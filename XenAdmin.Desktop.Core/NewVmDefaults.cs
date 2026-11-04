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
        public static string Category(VM template)
        {
            var t = template.TemplateType();
            switch (t)
            {
                case VM.VmTemplateType.Custom: return Messages.NEWVMWIZARD_TEMPLATEPAGE_CUSTOM;
                case VM.VmTemplateType.Windows:
                case VM.VmTemplateType.WindowsServer:
                case VM.VmTemplateType.LegacyWindows: return Messages.NEWVMWIZARD_TEMPLATEPAGE_WINDOWS;
                case VM.VmTemplateType.Centos: return Messages.NEWVMWIZARD_TEMPLATEPAGE_CENTOS;
                case VM.VmTemplateType.Debian: return Messages.NEWVMWIZARD_TEMPLATEPAGE_DEBIAN;
                case VM.VmTemplateType.Ubuntu: return Messages.NEW_VM_WIZARD_TEMPLATEPAGE_UBUNTU;
                case VM.VmTemplateType.RedHat: return Messages.NEWVMWIZARD_TEMPLATEPAGE_REDHAT;
                case VM.VmTemplateType.Rocky: return Messages.NEWVMWIZARD_TEMPLATEPAGE_ROCKY;
                case VM.VmTemplateType.Oracle: return Messages.NEWVMWIZARD_TEMPLATEPAGE_ORACLE;
                case VM.VmTemplateType.Suse: return Messages.NEWVMWIZARD_TEMPLATEPAGE_SUSE;
                default: return Messages.NEWVMWIZARD_TEMPLATEPAGE_MISC;
            }
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
