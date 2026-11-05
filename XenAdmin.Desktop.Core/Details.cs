using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Core;
using XenAdmin.Model;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    public class Section
    {
        public string Title;
        public bool Expanded;
        public List<(string Key, string Value)> Rows = new List<(string, string)>();
        public Section Add(string key, string value) { Rows.Add((key, string.IsNullOrEmpty(value) ? "" : value)); return this; }
    }

    public class Table
    {
        public string Title;
        public string[] Headers;
        /// <summary>Index of the column that takes the spare width (DataGridView AutoSizeMode=Fill).</summary>
        public int Fill = -1;
        public List<string[]> Rows = new List<string[]>();
    }

    /// <summary>
    /// Same sections, labels and values as the WinForms General/Storage/Networking/Memory tabs
    /// (GeneralTabPage, VMStoragePage, NetworkPage), built only from XenModel.
    /// </summary>
    public static class Details
    {
        private static string FN(string property) =>
            FriendlyNameManager.GetFriendlyName($"Label-{property}") ?? property;

        private static string Mem(long bytes) => Util.MemorySizeStringSuitableUnits(bytes, true);

        public static string TabTitle(IXenObject o, string tab)
        {
            var kind = o switch { VM _ => "VM", Host _ => "Server", Pool _ => "Pool", SR _ => "Storage", _ => "" };
            return tab == "General" ? $"{kind} General Properties" : tab;
        }

        public static string Heading(IXenObject o)
        {
            switch (o)
            {
                case VM vm:
                    var host = vm.Connection?.Resolve(vm.resident_on) ?? vm.Home();
                    return host == null ? vm.Name() : $"{vm.Name()} on '{host.Name()}'";
                case Host h when Helpers.GetPool(h.Connection) is Pool p:
                    return $"{h.Name()} in '{p.Name()}'";
                case SR sr:
                    return sr.Name();
                default:
                    return o.Name();
            }
        }

        /// <summary>Same tabs and order as the WinForms app; the ones in <see cref="NotPorted"/> are shown disabled.</summary>
        public static string[] Tabs(IXenObject o) => o switch
        {
            VM _ => new[] { "General", "Memory", "Storage", "Networking", "Console", "Performance", "Snapshots", "Search" },
            Host _ => new[] { "General", "Memory", "Storage", "Networking", "NICs", "GPU", "USB", "Console", "Performance", "Search" },
            Pool _ => new[] { "General", "Memory", "Storage", "Networking", "GPU", "Search" },
            SR _ => new[] { "General", "Storage", "Search" },
            _ => new[] { "General" }
        };

        public static bool NotPorted(IXenObject o, string tab) =>
            tab is "NICs" or "GPU" or "USB" || tab == "Console" && !(o is VM) || tab == "Search" && o is SR;

        private static Section Common(IXenObject o)
        {
            var s = new Section { Title = "General", Expanded = true };
            s.Add(FN("host.name_label"), Helpers.GetName(o));
            if (!(o is VM vm) || vm.DescriptionType() != VM.VmDescriptionType.None)
                s.Add(FN("host.name_description"), o.Description());
            var tags = Tags.GetTags(o);
            s.Add(Messages.TAGS, tags == null || tags.Length == 0 ? Messages.NONE : string.Join(", ", tags.OrderBy(t => t)));
            s.Add(Messages.FOLDER, string.IsNullOrEmpty(o.Path) ? Messages.NONE : o.Path);
            return s;
        }

        public static List<Section> General(IXenObject o)
        {
            var list = new List<Section>();
            var s = Common(o);
            list.Add(s);
            switch (o)
            {
                case Host host:
                    if (Helpers.GetPool(host.Connection) != null)
                        s.Add(Messages.POOL_COORDINATOR, host.IsCoordinator() ? Messages.YES : Messages.NO);
                    s.Add(FN("host.enabled"), !host.IsLive() ? Messages.HOST_NOT_LIVE
                        : !host.enabled ? (host.MaintenanceMode() ? Messages.HOST_IN_MAINTENANCE_MODE : Messages.DISABLED)
                        : Messages.YES);
                    s.Add("Autoboot of VMs enabled", Helpers.BoolToString(host.GetVmAutostartEnabled()));
                    var pool1 = Helpers.GetPoolOfOne(host.Connection);
                    if (pool1 != null && Helpers.CloudOrGreater(host) && Helpers.XapiEqualOrGreater_1_290_0(host))
                        s.Add(Messages.CERTIFICATE_VERIFICATION_KEY,
                            pool1.tls_verification_enabled && (!Helpers.XapiEqualOrGreater_1_313_0(host) || host.tls_verification_enabled)
                                ? Messages.ENABLED : Messages.DISABLED);
                    s.Add(FN("host.iscsi_iqn"), host.GetIscsiIqn());
                    var sysLog = host.GetSysLogDestination();
                    s.Add(FN("host.log_destination"), string.IsNullOrEmpty(sysLog)
                        ? Messages.HOST_LOG_DESTINATION_LOCAL
                        : string.Format(Messages.HOST_LOG_DESTINATION_LOCAL_AND_REMOTE, sysLog));
                    s.Add(FN("host.uptime"), host.Uptime()?.ToString());
                    s.Add(FN("host.agentUptime"), host.AgentUptime()?.ToString());
                    s.Add(FN("host.uuid"), host.uuid);

                    var version = new Section { Title = "Version Details" };
                    if (host.software_version != null)
                    {
                        if (host.software_version.TryGetValue("date", out var date)) version.Add(Messages.SOFTWARE_VERSION_DATE, date);
                        version.Add(Messages.SOFTWARE_VERSION_PRODUCT_VERSION, $"{host.ProductBrand()} {host.ProductVersionText()}");
                        if (host.software_version.TryGetValue("build_number", out var build)) version.Add("Build number", build);
                    }
                    list.Add(version);
                    list.Insert(1, Certificates(host));
                    list.Add(ManagementInterfaces(host.Connection, new[] { host }));
                    list.Add(HostMemory(host));
                    var cpus = new Section { Title = "CPUs" };
                    var allCpus = host.Connection.ResolveAll(host.host_CPUs).OrderBy(c => c.number).ToList();
                    if (allCpus.Count > 0)
                        cpus.Add(string.Format(Messages.GENERAL_DETAILS_CPU_NUMBER, allCpus.Count > 1 ? $"0 - {allCpus.Count - 1}" : "0"),
                            Helpers.GetCPUProperties(allCpus[0]));
                    list.Add(cpus);
                    break;

                case VM vm:
                    s.Add(FN("VM.OSName"), vm.GetOSName());
                    s.Add(FN("VM.OperatingMode"), vm.IsHVM() ? Messages.VM_OPERATING_MODE_HVM : Messages.VM_OPERATING_MODE_PV);
                    if (!vm.DefaultTemplate())
                        s.Add(Messages.BIOS_STRINGS_COPIED, vm.BiosStringsCopied() ? Messages.YES : Messages.NO);
                    if (!vm.is_a_template)
                        VirtualizationState(s, vm);
                    var running = vm.RunningTime();
                    if (running != null)
                        s.Add(FN("VM.uptime"), running.ToString());
                    // GeneralTabPage.VMCanChooseHomeServer: hidden under WLB or when the disks are not on shared storage
                    if (!vm.is_a_template && !Helpers.WlbEnabledAndConfigured(vm.Connection)
                        && (string.IsNullOrEmpty(vm.IsOnSharedStorage()) || vm.HasNoDisksAndNoLocalCD()))
                        s.Add(FN("VM.affinity"), vm.AffinityServerString());
                    s.Add(FN("host.uuid"), vm.uuid);

                    var boot = new Section { Title = "Boot Options" };
                    boot.Add(FN("VM.auto_boot"), Helpers.BoolToString(vm.GetAutoPowerOn()));
                    if (vm.IsHVM())
                    {
                        var order = vm.GetBootOrder().ToUpper().Union(new[] { 'D', 'C', 'N' });
                        boot.Add(FN("VM.BootOrder"), string.Join("\n", order.Select(c =>
                            c == 'C' ? Messages.BOOT_HARD_DISK : c == 'D' ? Messages.BOOT_DVD_DRIVE : Messages.NETWORK)));
                        if (Helpers.NaplesOrGreater(vm.Connection))
                        {
                            string mode = Messages.BIOS_BOOT;
                            if (vm.IsDefaultBootModeUefi())
                            {
                                var secure = vm.GetSecureBootMode();
                                var hasCerts = !string.IsNullOrEmpty(Helpers.GetPoolOfOne(vm.Connection)?.uefi_certificates);
                                mode = secure == "true" && !hasCerts
                                    ? string.Format(Messages.CUSTOM_FIELD_NAME_AND_TYPE, Messages.UEFI_SECURE_BOOT, Messages.GUEFI_SECUREBOOT_MODE_MISSING_CERTIFICATES)
                                    : secure == "true" || secure == "auto" && hasCerts ? Messages.UEFI_SECURE_BOOT : Messages.UEFI_BOOT;
                            }
                            boot.Add(FN("VM.BootMode"), mode);
                        }
                    }
                    else
                        boot.Add(FN("VM.PV_args"), vm.PV_args);
                    list.Add(boot);

                    var vcpus = new Section { Title = "CPUs" };
                    vcpus.Add(FN("VM.VCPUs"), vm.VCPUs_at_startup.ToString());
                    if (vm.VCPUs_at_startup != vm.VCPUs_max || vm.SupportsVcpuHotplug())
                        vcpus.Add(FN("VM.MaxVCPUs"), vm.VCPUs_max.ToString());
                    vcpus.Add(FN("VM.Topology"), vm.Topology());
                    list.Add(vcpus);
                    if (vm.IsRunning())
                    {
                        var rc = new Section { Title = "Read Caching" };
                        var pvsProxy = vm.PvsProxy();
                        if (pvsProxy != null)
                            rc.Add(FN("VM.pvs_read_caching_status"), pvs_proxy_status_extensions.ToFriendlyString(pvsProxy.status));
                        else if (vm.ReadCachingEnabled())
                        {
                            rc.Add(FN("VM.read_caching_status"), Messages.VM_READ_CACHING_ENABLED);
                            rc.Add(FN("VM.read_caching_disks"), string.Join("\n", vm.ReadCachingVDIs().Select(v => v.NameWithLocation())));
                        }
                        else
                        {
                            rc.Add(FN("VM.read_caching_status"), Messages.VM_READ_CACHING_DISABLED);
                            var reason = vm.ReadCachingDisabledReason();
                            if (reason != null)
                                rc.Add(FN("VM.read_caching_reason"), reason);
                        }
                        list.Add(rc);
                    }
                    break;

                case Pool pool:
                    s.Add(Messages.NUMBER_OF_SOCKETS, pool.CpuSockets().ToString());
                    s.Add("Autoboot of VMs enabled", Helpers.BoolToString(pool.GetVmAutostartEnabled()));
                    if (Helpers.CloudOrGreater(pool.Connection) && Helpers.XapiEqualOrGreater_1_290_0(pool.Connection))
                        s.Add(Messages.CERTIFICATE_VERIFICATION_KEY, pool.tls_verification_enabled ? Messages.ENABLED : Messages.DISABLED);
                    var coordinator = pool.Connection.Resolve(pool.master);
                    if (coordinator != null)
                        s.Add(Messages.SOFTWARE_VERSION_PRODUCT_VERSION, $"{coordinator.ProductBrand()} {coordinator.ProductVersionText()}");
                    s.Add(FN("host.uuid"), pool.uuid);
                    list.Add(ManagementInterfaces(pool.Connection, pool.Connection.Cache.Hosts));
                    break;

                case SR sr:
                    s.Add(Messages.TYPE, sr.FriendlyTypeName());
                    if (sr.content_type != SR.Content_Type_ISO && sr.GetSRType(false) != SR.SRTypes.udev)
                        s.Add(FN("SR.size"), sr.SizeString());
                    s.Add(FN("SR.scsiid"), sr.GetScsiID() ?? "");
                    s.Add(FN("host.uuid"), sr.uuid);
                    var status = new Section { Title = "Status" };
                    var multipath = new Section { Title = "Multipathing" };
                    if (!sr.MultipathCapable())
                        multipath.Add(Messages.MULTIPATH_CAPABLE, Messages.NO);
                    else
                        foreach (var h in sr.Connection.Cache.Hosts.OrderBy(h => h.Name()))
                        {
                            var pbd = sr.GetPBDFor(h);
                            multipath.Add(h.Name(), pbd != null && pbd.MultipathActive() ? Messages.MULTIPATH_ACTIVE : Messages.MULTIPATH_NOT_ACTIVE);
                        }
                    status.Add(FN("SR.state"), sr.IsDetached() ? Messages.DETACHED : sr.IsBroken() ? Messages.GENERAL_SR_STATE_BROKEN
                        : !sr.MultipathAOK() ? Messages.GENERAL_MULTIPATH_FAILURE : Messages.GENERAL_STATE_OK);
                    foreach (var pbd in sr.Connection.ResolveAll(sr.PBDs))
                        status.Add(Helpers.GetName(sr.Connection.Resolve(pbd.host)), pbd.StatusString());
                    list.Add(status);
                    list.Add(multipath);
                    break;
            }
            return list;
        }

        /// <summary>GeneralTabPage.GenerateVirtualisationStatusForGeneralBox.</summary>
        private static void VirtualizationState(Section s, VM vm)
        {
            if (vm.HasNewVirtualizationStates())
            {
                var status = vm.GetVirtualizationStatus(out var statusString);
                var lines = new List<string>();
                if (vm.power_state == vm_power_state.Running)
                {
                    if (status.HasFlag(VM.VirtualizationStatus.Unknown))
                        lines.Add(statusString);
                    else
                    {
                        lines.Add(status.HasFlag(VM.VirtualizationStatus.IoDriversInstalled) ? Messages.VIRTUALIZATION_STATE_VM_IO_OPTIMIZED : Messages.VIRTUALIZATION_STATE_VM_IO_NOT_OPTIMIZED);
                        lines.Add(status.HasFlag(VM.VirtualizationStatus.ManagementInstalled) ? Messages.VIRTUALIZATION_STATE_VM_MANAGEMENT_AGENT_INSTALLED : Messages.VIRTUALIZATION_STATE_VM_MANAGEMENT_AGENT_NOT_INSTALLED);
                    }
                }
                lines.Add(vm.has_vendor_device ? Messages.VIRTUALIZATION_STATE_VM_RECEIVING_UPDATES : Messages.VIRTUALIZATION_STATE_VM_NOT_RECEIVING_UPDATES);
                s.Add(FN("VM.VirtualizationState"), string.Join("\n", lines));
            }
            else if (vm.power_state == vm_power_state.Running)
            {
                vm.GetVirtualizationStatus(out var statusString);
                s.Add(FN("VM.VirtualizationState"), statusString);
            }
        }

        private static Section Certificates(Host host)
        {
            var s = new Section { Title = "Certificates" };
            if (!Helpers.StockholmOrGreater(host) || host.certificates == null)
                return s;
            foreach (var c in host.certificates.Select(r => host.Connection.Resolve(r)).Where(c => c != null).OrderBy(c => c.type))
            {
                var type = c.type switch
                {
                    certificate_type.ca => Messages.CERTIFICATE_TYPE_CA,
                    certificate_type.host => Messages.CERTIFICATE_TYPE_HOST,
                    _ => Messages.CERTIFICATE_TYPE_HOST_INTERNAL
                };
                var validity = string.Format(Messages.CERTIFICATE_VALIDITY_PERIOD_VALUE,
                    c.not_before.ToLocalTime().ToString(Messages.DATEFORMAT_DMY_HM), c.not_after.ToLocalTime().ToString(Messages.DATEFORMAT_DMY_HM));
                s.Add(type, $"{validity}\n{string.Format(Messages.CERTIFICATE_THUMBPRINT_VALUE, c.fingerprint)}");
            }
            return s;
        }

        private static Section ManagementInterfaces(IXenConnection c, IEnumerable<Host> hosts)
        {
            var s = new Section { Title = "Management Interfaces" };
            foreach (var host in hosts.OrderBy(h => h.Name()))
            {
                s.Add(FN("host.hostname"), host.hostname);
                foreach (var pif in c.ResolveAll(host.PIFs).Where(p => p.management))
                    s.Add(Messages.MANAGEMENT_INTERFACE, pif.FriendlyIPAddress());
            }
            return s;
        }

        private static Section HostMemory(Host host)
        {
            var s = new Section { Title = "Memory" };
            s.Add(FN("host.ServerMemory"), host.HostMemoryString());
            s.Add(FN("host.VMMemory"), host.ResidentVMMemoryUsageString());
            s.Add(FN("host.XenMemory"), host.XenMemoryString());
            return s;
        }

        public static List<Section> Memory(IXenObject o)
        {
            switch (o)
            {
                case Pool pool:
                    return pool.Connection.Cache.Hosts.OrderBy(h => h.Name()).Select(h =>
                    {
                        var m = HostMemory(h);
                        m.Title = h.Name();
                        m.Expanded = true;
                        return m;
                    }).ToList();
                case Host host:
                    var h = HostMemory(host);
                    h.Expanded = true;
                    return new List<Section> { h };
                case VM vm:
                    var s = new Section { Title = "Memory", Expanded = true };
                    s.Add("Dynamic minimum", Mem(vm.memory_dynamic_min));
                    s.Add("Dynamic maximum", Mem(vm.memory_dynamic_max));
                    s.Add("Static maximum", Mem(vm.memory_static_max));
                    var metrics = vm.Connection.Resolve(vm.metrics);
                    if (vm.power_state == vm_power_state.Running && metrics != null)
                        s.Add("Current memory usage", Mem(metrics.memory_actual));
                    return new List<Section> { s };
            }
            return new List<Section>();
        }

        public static Table Storage(IXenObject o)
        {
            switch (o)
            {
                case VM vm:
                {
                    var t = new Table { Title = "Virtual Disks", Fill = 2, Headers = new[] { "Position", "Name", "Description", "SR", "Size", "Read Only", "Priority", "Active", "Device Path" } };
                    foreach (var vbd in vm.Connection.ResolveAll(vm.VBDs).Where(b => b.type != vbd_type.CD).OrderBy(b => b.userdevice))
                    {
                        var vdi = vm.Connection.Resolve(vbd.VDI);
                        if (vdi == null) continue;
                        var sr = vm.Connection.Resolve(vdi.SR);
                        t.Rows.Add(new[] { vbd.userdevice, vdi.Name(), vdi.Description(), sr?.NameWithLocation() ?? "", vdi.SizeText(),
                            vbd.IsReadOnly() ? Messages.YES : Messages.NO,
                            vbd.GetIoNice() switch { 0 => Messages.STORAGE_WORST, 7 => Messages.STORAGE_BEST, var p => p.ToString() },
                            vbd.currently_attached ? Messages.YES : Messages.NO,
                            string.IsNullOrEmpty(vbd.device) ? "" : "/dev/" + vbd.device });
                    }
                    return t;
                }
                case SR sr:
                {
                    var t = new Table { Title = "Virtual Disks", Headers = new[] { "Name", "Description", "Size", "VM" } };
                    foreach (var vdi in sr.Connection.ResolveAll(sr.VDIs).Where(v => v.Show(false)).OrderBy(v => v.Name()))
                    {
                        var vms = vdi.Connection.ResolveAll(vdi.VBDs).Select(b => vdi.Connection.Resolve(b.VM)).Where(v => v != null).Select(v => v.Name());
                        t.Rows.Add(new[] { vdi.Name(), vdi.Description(), vdi.SizeText(), string.Join(", ", vms) });
                    }
                    return t;
                }
                case Host _:
                case Pool _:
                {
                    var t = new Table { Title = "Storage Repositories", Headers = new[] { "Name", "Description", "Type", "Shared", "Usage", "Size", "Virtual allocation" } };
                    var host = o as Host;
                    foreach (var sr in o.Connection.Cache.SRs.Where(s => s.Show(false) && !s.IsToolsSR()).OrderBy(s => s.Name()))
                    {
                        if (host != null && !sr.shared && sr.Home()?.opaque_ref != host.opaque_ref) continue;
                        var pct = sr.physical_size > 0 ? $"{sr.physical_utilisation * 100 / sr.physical_size}% ({Util.DiskSizeString(sr.physical_utilisation)} used)" : "";
                        t.Rows.Add(new[] { sr.Name(), sr.Description(), sr.FriendlyTypeName(), sr.shared ? Messages.YES : Messages.NO,
                            pct, Util.DiskSizeString(sr.physical_size), Util.DiskSizeString(sr.virtual_allocation) });
                    }
                    return t;
                }
            }
            return null;
        }

        public static Table Networking(IXenObject o)
        {
            if (o is VM vm)
            {
                var t = new Table { Title = "Virtual Network Interfaces", Headers = new[] { "Device", "MAC", "Network", "IP Address", "Active" } };
                var gm = vm.Connection.Resolve(vm.guest_metrics);
                foreach (var vif in vm.Connection.ResolveAll(vm.VIFs).OrderBy(v => v.device))
                {
                    var net = vm.Connection.Resolve(vif.network);
                    var ips = gm?.networks?.Where(kv => kv.Key.StartsWith(vif.device + "/ip")).Select(kv => kv.Value) ?? Enumerable.Empty<string>();
                    t.Rows.Add(new[] { vif.device, vif.MAC, net?.Name() ?? "", string.Join(", ", ips),
                        vif.currently_attached ? Messages.YES : Messages.NO });
                }
                return t;
            }
            if (o is Host || o is Pool)
            {
                var t = new Table { Title = "Networks", Headers = new[] { "Name", "Description", "NIC", "VLAN", "MAC", "MTU" } };
                var host = o as Host ?? Helpers.GetCoordinator(o.Connection);
                foreach (var net in o.Connection.Cache.Networks.Where(n => n.Show(false)).OrderBy(n => n.Name()))
                {
                    var pif = o.Connection.ResolveAll(net.PIFs).FirstOrDefault(p => host == null || p.host.opaque_ref == host.opaque_ref);
                    t.Rows.Add(new[] { net.Name(), net.Description(), pif?.Name() ?? "", pif == null || pif.VLAN < 0 ? "-" : pif.VLAN.ToString(),
                        pif?.MAC ?? "", net.MTU.ToString() });
                }
                return t;
            }
            return null;
        }
    }
}
