using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using XenAdmin.Actions.OvfActions;
using XenAdmin.Core;
using XenAdmin.Mappings;
using XenAdmin.Network;
using XenAPI;
using XenOvf;
using XenOvf.Definitions;
using XenOvf.Utilities;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// Import of an OVF/OVA appliance: the import wizard's appliance path (ImportSourcePage.LoadAppliance, ImportEulaPage,
    /// ImportSelectHostPage, the storage and network pages, ImportSecurityPage) ending in ImportApplianceAction.
    /// shortcut: one target, one SR for every disk and one network for every interface (the WinForms pages allow a
    /// choice per VM, disk and interface), no OS fixup and no encrypted packages; add the per-item grids when needed.
    /// </summary>
    public static class OvfRules
    {
        /// <summary>An OVF descriptor or an OVA archive (optionally compressed), by its name.</summary>
        public static bool IsAppliance(string path)
        {
            var name = (path ?? "").ToLowerInvariant();
            return name.EndsWith(".ovf") || name.EndsWith(".ova") || name.EndsWith(".ova.gz") || name.EndsWith(".ova.bz2");
        }

        private static readonly object RegisterLock = new object();
        private static bool _registered;

        /// <summary>Program.Main (WinForms): DiscUtils opens VHD, VHDX, VMDK and raw disks only once they are registered.</summary>
        public static void RegisterDiskFormats()
        {
            lock (RegisterLock)
            {
                if (_registered) return;
                DiscUtils.Setup.SetupHelper.RegisterAssembly(typeof(DiscUtils.Vhd.Disk).Assembly);
                DiscUtils.Setup.SetupHelper.RegisterAssembly(typeof(DiscUtils.Vhdx.Disk).Assembly);
                DiscUtils.Setup.SetupHelper.RegisterAssembly(typeof(DiscUtils.Vmdk.Disk).Assembly);
                DiscUtils.Setup.SetupHelper.RegisterAssembly(typeof(DiscUtils.Raw.Disk).Assembly);
                _registered = true;
            }
        }

        /// <summary>ImportSourcePage.LoadAppliance: the package, or the reason it cannot be imported.</summary>
        public static (Package Package, string Error, List<string> Warnings) Open(string path)
        {
            RegisterDiskFormats();
            if (string.IsNullOrWhiteSpace(path) || !File.Exists(path))
                return (null, "The file does not exist.", null);
            try
            {
                var package = Package.Create(path);
                if (!OVF.Validate(package, out var warnings))
                    return (null, warnings.LastOrDefault() ?? "The appliance is not valid.", warnings);
                if (package.HasEncryption())
                    return (null, "Encrypted appliances cannot be imported in this version.", warnings);
                return (package, null, warnings ?? new List<string>());
            }
            catch (Exception e)
            {
                return (null, e.Message, null);
            }
        }

        /// <summary>The VMs of the appliance: their system ids and names (ImportWizard.FindVMName, first choices).</summary>
        public static List<(string SysId, string Name)> Vms(Package package)
        {
            var env = package.OvfEnvelope;
            return OVF.FindSystemIds(env).Select(id =>
            {
                var vs = OVF.FindVirtualSystemById(env, id);
                var name = Tools.ValidateProperty("Name", vs) && !string.IsNullOrEmpty(vs.Name[0].Value) ? vs.Name[0].Value
                    : !string.IsNullOrEmpty(vs.id) ? vs.id : package.Name;
                return (id, name);
            }).ToList();
        }

        /// <summary>ImportEulaPage: the license texts the user must accept (empty when there are none).</summary>
        public static List<string> Eulas(Package package) =>
            OVF.HasEula(package.OvfEnvelope)
                ? OVF.FindSections<EulaSection_Type>(package.OvfEnvelope).Select(e => OVF.FindStringsMessage(package.OvfEnvelope, e.License)).ToList()
                : new List<string>();

        /// <summary>ImportSecurityPage: the manifest can be verified.</summary>
        public static bool HasManifest(Package package) => package.HasManifest();

        /// <summary>The number of disks and network interfaces, for the summary.</summary>
        public static (int Disks, int Interfaces) Counts(Package package) =>
            (OVF.FindSystemIds(package.OvfEnvelope).Sum(id => OVF.FindDiskRasds(package.OvfEnvelope, id).Length),
             OVF.FindSystemIds(package.OvfEnvelope).Sum(id => OVF.FindRasdByType(package.OvfEnvelope, id, 10).Length));

        /// <summary>
        /// SelectMultipleVMNetworkPage.ShowNetwork: the networks a VM interface can use on the target (a server sees it,
        /// or every server does for a pool), visible and not bond members, auto-add networks first.
        /// </summary>
        public static List<XenAPI.Network> Networks(IXenConnection c, Host target) =>
            c.Cache.Networks.Where(n => n.Show(false) && !n.IsMember() && !n.IsSriov()
                                        && (target != null ? target.CanSeeNetwork(n) : n.AllHostsCanSeeNetwork()))
                .OrderBy(n => !n.GetAutoPlug()).ThenBy(n => n.Name()).ToList();

        /// <summary>The space the appliance's disks need (their capacities, as the storage page adds them up).</summary>
        public static ulong RequiredSpace(Package package)
        {
            var env = package.OvfEnvelope;
            ulong total = 0;
            foreach (var id in OVF.FindSystemIds(env))
                foreach (var rasd in OVF.FindDiskRasds(env, id))
                    if (OVF.FindDiskReference(env, rasd) is VirtualDiskDesc_Type disk && !string.IsNullOrEmpty(disk.capacity))
                        total += OVF.ComputeCapacity(Convert.ToInt64(disk.capacity), disk.capacityAllocationUnits);
            return total;
        }

        /// <summary>null when the SR has room for every disk of the appliance (thin-provisioned SRs are not checked).</summary>
        public static string SpaceError(SR sr, ulong required)
        {
            var sm = sr.GetSM();
            var thin = sm != null && Array.IndexOf(sm.capabilities, "THIN_PROVISIONING") != -1;
            return !thin && sr.FreeSpace() < (long)required
                ? string.Format("'{0}' has {1} free; the appliance's disks need {2}.", sr.Name(), Util.DiskSizeString(sr.FreeSpace()), Util.DiskSizeString((long)required))
                : null;
        }

        /// <summary>What the import dialog shows about a package, read once (off the UI thread: an OVA may be large).</summary>
        public sealed record Info(Package Package, string Error, List<string> Warnings, List<(string SysId, string Name)> Vms,
            int Disks, int Interfaces, List<string> Eulas, bool HasManifest, ulong RequiredSpace);

        public static Info Load(string path)
        {
            var (package, error, warnings) = Open(path);
            if (package == null)
                return new Info(null, error, warnings, null, 0, 0, null, false, 0);
            var (disks, vifs) = Counts(package);
            return new Info(package, null, warnings, Vms(package), disks, vifs, Eulas(package), HasManifest(package), RequiredSpace(package));
        }

        /// <summary>
        /// The mappings of ImportSelectHostPage, the storage page and the network page: every VM on the target, every disk
        /// on that SR and every interface on that network.
        /// </summary>
        public static Dictionary<string, VmMapping> Mappings(Package package, IXenObject target, SR sr, XenAPI.Network network)
        {
            var env = package.OvfEnvelope;
            var result = new Dictionary<string, VmMapping>();
            foreach (var (sysId, name) in Vms(package))
            {
                var mapping = new VmMapping(sysId)
                {
                    VmNameLabel = name,
                    XenRef = target is Pool ? (object)new XenRef<Pool>(target.opaque_ref) : new XenRef<Host>(target.opaque_ref),
                    TargetName = Helpers.GetName(target)
                };
                foreach (var rasd in OVF.FindDiskRasds(env, sysId))
                    if (rasd.InstanceID?.Value is string rasdId)
                        mapping.Storage[rasdId] = sr;
                foreach (var rasd in OVF.FindRasdByType(env, sysId, 10))
                {
                    if (rasd.InstanceID?.Value is string netId)
                        mapping.Networks[netId] = network;
                    // SelectMultipleVMNetworkPage: interfaces without a MAC are keyed by the "autogenerated" text
                    var mac = rasd.Address == null || string.IsNullOrEmpty(rasd.Address.Value) ? Messages.NEWVMWIZARD_NETWORKINGPAGE_AUTOGEN : rasd.Address.Value;
                    mapping.VIFs[mac] = network;
                }
                result[sysId] = mapping;
            }
            return result;
        }

        /// <summary>ImportWizard.FinishWizard (OVF): ImportApplianceAction.</summary>
        public static XenAdmin.Actions.AsyncAction ImportAction(IXenConnection c, Package package, Dictionary<string, VmMapping> mappings,
            bool verifyManifest, bool start) =>
            new ImportApplianceAction(c, package, mappings, verifyManifest, false, null, false, null, start);
    }
}
