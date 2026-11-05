using System;
using System.Collections.Generic;
using System.Linq;
using System.Xml;
using XenAdmin.Actions;
using XenAdmin.Actions.VMActions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>Defaults and the final action of the New VM wizard (Page_Template, Page_Storage, Page_Networking, FinishWizard).</summary>
    public static class NewVmDefaults
    {
        /// <summary>Page_Template category column (VmTemplateType display names).</summary>
        /// <summary>Page_Template's Category column (VmTemplateTypeExtensions.ToDisplayString).</summary>
        public static string Category(VM template) => template.TemplateType() switch
        {
            VM.VmTemplateType.Custom => Messages.NEWVMWIZARD_TEMPLATEPAGE_CUSTOM,
            VM.VmTemplateType.Windows or VM.VmTemplateType.WindowsServer or VM.VmTemplateType.LegacyWindows => Messages.NEWVMWIZARD_TEMPLATEPAGE_WINDOWS,
            VM.VmTemplateType.Centos => Messages.NEWVMWIZARD_TEMPLATEPAGE_CENTOS,
            VM.VmTemplateType.Debian => Messages.NEWVMWIZARD_TEMPLATEPAGE_DEBIAN,
            VM.VmTemplateType.Gooroom => Messages.NEWVMWIZARD_TEMPLATEPAGE_GOOROOM,
            VM.VmTemplateType.Rocky => Messages.NEWVMWIZARD_TEMPLATEPAGE_ROCKY,
            VM.VmTemplateType.Linx => Messages.NEWVMWIZARD_TEMPLATEPAGE_LINX,
            VM.VmTemplateType.Oracle => Messages.NEWVMWIZARD_TEMPLATEPAGE_ORACLE,
            VM.VmTemplateType.RedHat => Messages.NEWVMWIZARD_TEMPLATEPAGE_REDHAT,
            VM.VmTemplateType.SciLinux => Messages.NEWVMWIZARD_TEMPLATEPAGE_SCILINUX,
            VM.VmTemplateType.Suse => Messages.NEWVMWIZARD_TEMPLATEPAGE_SUSE,
            VM.VmTemplateType.Ubuntu => Messages.NEW_VM_WIZARD_TEMPLATEPAGE_UBUNTU,
            VM.VmTemplateType.YinheKylin => Messages.NEW_VM_WIZARD_TEMPLATEPAGE_YINHEKYLIN,
            VM.VmTemplateType.NeoKylin => Messages.NEW_VM_WIZARD_TEMPLATEPAGE_NEOKYLIN,
            VM.VmTemplateType.Asianux => Messages.NEW_VM_WIZARD_TEMPLATEPAGE_ASIANUX,
            VM.VmTemplateType.Turbo => Messages.NEW_VM_WIZARD_TEMPLATEPAGE_TURBO,
            VM.VmTemplateType.Citrix => BrandManager.CompanyNameShort,
            VM.VmTemplateType.Solaris or VM.VmTemplateType.Misc => Messages.NEWVMWIZARD_TEMPLATEPAGE_MISC,
            VM.VmTemplateType.CoreOS => Messages.NEWVMWIZARD_TEMPLATEPAGE_COREOS,
            _ => Messages.NEWVMWIZARD_TEMPLATEPAGE_SNAPSHOTS
        };

        /// <summary>Page_Template's icon column (VmTemplateTypeExtensions.ToBitmap): the OS logo of the template.</summary>
        public static string TemplateIcon(VM template) => TemplateIcon(template.TemplateType());

        public static string TemplateIcon(VM.VmTemplateType type) => type switch
        {
            VM.VmTemplateType.Custom => "000_UserTemplate_h32bit_16",
            VM.VmTemplateType.Windows or VM.VmTemplateType.WindowsServer or VM.VmTemplateType.LegacyWindows => "windows_h32bit_16",
            VM.VmTemplateType.Centos => "centos_16x",
            VM.VmTemplateType.Debian => "debian_16x",
            VM.VmTemplateType.Gooroom => "gooroom_16x",
            VM.VmTemplateType.Rocky => "rocky_16x",
            VM.VmTemplateType.Linx => "linx_16x",
            VM.VmTemplateType.Oracle => "oracle_16x",
            VM.VmTemplateType.RedHat => "redhat_16x",
            VM.VmTemplateType.SciLinux => "scilinux_16x",
            VM.VmTemplateType.Suse => "suse_16x",
            VM.VmTemplateType.Ubuntu => "ubuntu_16x",
            VM.VmTemplateType.YinheKylin => "yinhekylin_16x",
            VM.VmTemplateType.NeoKylin => "neokylin_16x",
            VM.VmTemplateType.Asianux => "asianux_16x",
            VM.VmTemplateType.Turbo => "turbo_16x",
            VM.VmTemplateType.Citrix => "Logo", // the product logo (Images.Get maps it to the app icon)
            VM.VmTemplateType.Solaris or VM.VmTemplateType.Misc => "000_VMTemplate_h32bit_16",
            VM.VmTemplateType.CoreOS => "coreos-globe-icon",
            _ => "000_VMSession_h32bit_16"
        };

        /// <summary>
        /// Page_Template's TemplatesGridViewItem order: by template type (custom templates and Windows first, hidden ones
        /// last), then newest version first (reverse natural order of the names).
        /// </summary>
        public static List<VM> SortTemplates(IEnumerable<VM> templates) =>
            templates.OrderBy(t => (int)t.TemplateType() + (t.IsHidden() ? (int)VM.VmTemplateType.Count : 0))
                .ThenBy(t => t.Name(), Comparer<string>.Create((a, b) => -XenCenterLib.StringUtility.NaturalCompare(a, b)))
                .ThenBy(t => t.opaque_ref, StringComparer.Ordinal).ToList();

        /// <summary>
        /// Page_Storage Add (NewDiskDialog): "New virtual disk" with a unique name, 1 GB, on the best SR, at the first free
        /// position (CreateVMAction moves it to an allowed device if needed).
        /// </summary>
        public static DiskDescription NewDisk(IXenConnection c, Host home, IList<DiskDescription> disks)
        {
            var names = c.Cache.VDIs.Select(v => v.Name()).Concat(disks.Select(d => d.Disk.name_label)).ToList();
            var used = disks.Select(d => d.Device.userdevice).ToHashSet();
            var position = Enumerable.Range(0, 16).Select(i => i.ToString()).First(i => !used.Contains(i) && i != "3"); // 3: the DVD drive
            var sr = BestSr(c, home, null);
            var disk = new VDI
            {
                name_label = Helpers.MakeUniqueName(Messages.DEFAULT_VDI_NAME, names), name_description = "", virtual_size = Util.BINARY_GIGA,
                type = vdi_type.user, read_only = false, SR = new XenRef<SR>(sr != null ? sr.opaque_ref : Helper.NullOpaqueRef)
            };
            return new DiskDescription(disk, new VBD { userdevice = position, bootable = position == "0", mode = vbd_mode.RW });
        }

        /// <summary>Page_Storage.LoadDisks: provisioned disks of a default template, or the custom template's disks.</summary>
        public static List<DiskDescription> DisksFor(IXenConnection c, VM t, string vmName, Host home)
        {
            var disks = new List<DiskDescription>();
            XmlNode provision = t.ProvisionXml();
            if (provision != null)
            {
                foreach (XmlNode node in provision.ChildNodes)
                {
                    var device = new VBD { userdevice = node.Attributes["device"].Value, bootable = node.Attributes["bootable"].Value == "true", mode = vbd_mode.RW };
                    var disk = new VDI
                    {
                        name_label = string.Format(Messages.STRING_SPACE_STRING, vmName, device.userdevice),
                        name_description = Messages.NEWVMWIZARD_STORAGEPAGE_DISK_DESCRIPTION,
                        virtual_size = long.Parse(node.Attributes["size"].Value),
                        type = (vdi_type)Enum.Parse(typeof(vdi_type), node.Attributes["type"].Value),
                        read_only = false
                    };
                    var sr = BestSr(c, home, c.Cache.Find_By_Uuid<SR>(node.Attributes["sr"].Value));
                    disk.SR = new XenRef<SR>(sr != null ? sr.opaque_ref : Helper.NullOpaqueRef);
                    disks.Add(new DiskDescription(disk, device));
                }
                return disks;
            }
            foreach (var vbd in c.ResolveAll(t.VBDs).Where(b => b.type == vbd_type.Disk))
            {
                var vdi = c.Resolve(vbd.VDI);
                if (vdi == null) continue;
                var sr = BestSr(c, home, c.Resolve(vdi.SR));
                var disk = new VDI
                {
                    name_label = vdi.name_label, name_description = vdi.name_description, virtual_size = vdi.virtual_size, type = vdi.type,
                    read_only = vdi.read_only, sm_config = vdi.sm_config, SR = new XenRef<SR>(sr != null ? sr.opaque_ref : Helper.NullOpaqueRef)
                };
                disks.Add(new DiskDescription(disk, new VBD { userdevice = vbd.userdevice, bootable = vbd.bootable, mode = vbd.mode }));
            }
            return disks;
        }

        public sealed record TemplateDefaults(string Name, string Description, VmBootMode BootMode, InstallMethod Install, long Vcpus, long Cores, long Memory);

        /// <summary>Page_Name, Page_BootMode, Page_InstallationMedia and Page_CpuMem defaults for a template.</summary>
        public static TemplateDefaults FromTemplate(VM t) => new TemplateDefaults(
            Helpers.MakeUniqueName(t.Name(), t.Connection.Cache.VMs.Where(v => !v.is_a_template).Select(v => v.Name()).ToList()),
            t.DefaultTemplate() ? "" : t.Description(),
            t.IsHVM() && t.IsDefaultBootModeUefi() ? VmBootMode.Uefi : VmBootMode.Bios,
            t.IsHVM() ? InstallMethod.CD : InstallMethod.Network,
            Math.Max(1, t.VCPUs_at_startup),
            Math.Max(1, t.GetCoresPerSocket()),
            t.memory_static_max);

        /// <summary>Page_Storage.GetBestDiskStorage: suggested SR if visible, else the pool default, else any usable SR.</summary>
        public static SR BestSr(IXenConnection c, Host home, SR suggested)
        {
            bool Usable(SR sr) => StorageRules.CanHoldNewDisk(sr, home);
            if (Usable(suggested))
                return suggested;
            var def = c.Resolve(Helpers.GetPoolOfOne(c)?.default_SR);
            return Usable(def) ? def : c.Cache.SRs.Where(Usable).OrderByDescending(s => s.shared).ThenBy(s => s.Name()).FirstOrDefault();
        }

        /// <summary>Page_Networking.LoadNetworks: up to 4 auto-plug networks for default templates, else the template's VIFs.</summary>
        public static List<VIF> VifsFor(IXenConnection c, VM template)
        {
            var vifs = new List<VIF>();
            if (template.DefaultTemplate())
            {
                foreach (var n in c.Cache.Networks.OrderBy(n => n))
                {
                    if (n.IsGuestInstallerNetwork() || !n.GetAutoPlug() || !n.Show(false) || n.IsMember() || n.IsSriov())
                        continue;
                    if (vifs.Count >= 4)
                        break;
                    vifs.Add(new VIF { Connection = c, device = vifs.Count.ToString(), network = new XenRef<XenAPI.Network>(n.opaque_ref) });
                }
            }
            else
            {
                foreach (var v in c.ResolveAll(template.VIFs).OrderBy(v => v))
                    vifs.Add(new VIF { Connection = c, device = v.device, MAC = "", network = v.network, qos_algorithm_type = v.qos_algorithm_type, qos_algorithm_params = v.qos_algorithm_params });
            }
            return vifs;
        }

        /// <summary>The CreateVMAction exactly as NewVMWizard.FinishWizard builds it (no BIOS copy, vGPU, vTPM or cloud-config).</summary>
        public static CreateVMAction Build(IXenConnection c, VM template, string name, string description, InstallMethod install, VDI iso, string url,
            VmBootMode bootMode, Host home, long vcpus, long cores, long memory, List<DiskDescription> disks, List<VIF> vifs, bool start)
        {
            return new CreateVMAction(c, template, null, name, description, install, template.PV_args, iso, url, bootMode, home,
                vcpus, vcpus, memory, memory, memory, disks, null, vifs, start, false, null, null, null, false, cores, null);
        }

    }
}
