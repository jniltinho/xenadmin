using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Actions.GPU;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// The GPU tab (GpuPage with GpuPlacementPolicyPanel, GpuRow and GpuConfiguration, PoolGpuEditPage) and the USB tab
    /// (UsbPage) of servers and pools.
    /// </summary>
    public static class HardwareRules
    {
        /// <summary>MainWindow.ChangeToNewTabs: GPU for a pool or server unless GPUs are restricted by the license.</summary>
        public static bool ShowGpuTab(IXenObject o) =>
            (o is Pool || o is Host host && host.IsLive()) && o.Connection != null && !Helpers.FeatureForbidden(o.Connection, Host.RestrictGpu);

        /// <summary>MainWindow.ChangeToNewTabs: USB for a server with USB devices, unless pass-through is restricted.</summary>
        public static bool ShowUsbTab(IXenObject o) =>
            o is Host h && h.Connection != null && h.IsLive() && h.PUSBs != null && h.PUSBs.Count > 0 && !Helpers.FeatureForbidden(h.Connection, Host.RestrictUsbPassthrough);

        // ---------- GPU ----------

        /// <summary>GpuPage.Rebuild: the GPUs that support vGPU types, of the pool or of this server, by server then name.</summary>
        public static List<PGPU> Pgpus(IXenObject o) =>
            o.Connection.Cache.PGPUs
                .Select(p => (Pgpu: p, Host: o.Connection.Resolve(p.host)))
                .Where(x => x.Pgpu.supported_VGPU_types.Count > 0 && (o is Pool || x.Host != null && x.Host.opaque_ref == o.opaque_ref))
                .OrderBy(x => x.Host).ThenBy(x => x.Pgpu.Name())
                .Select(x => x.Pgpu).ToList();

        /// <summary>GpuPage.Rebuild: one row per distinct settings (name, enabled and supported vGPU types).</summary>
        public static List<List<PGPU>> Groups(IXenObject o)
        {
            string Key(PGPU p) => p.Name() + "|" + string.Join(",", p.enabled_VGPU_types.Select(r => r.opaque_ref).OrderBy(r => r))
                                  + "|" + string.Join(",", p.supported_VGPU_types.Select(r => r.opaque_ref).OrderBy(r => r));
            var groups = new List<List<PGPU>>();
            foreach (var g in Pgpus(o).GroupBy(Key))
                groups.Add(g.ToList());
            return groups;
        }

        public static string NoGpuText(IXenObject o) => o is Pool ? Messages.NO_GPU_IN_POOL : Messages.NO_GPU_ON_HOST;

        /// <summary>GpuRow.vGpuCapability: vGPU types can be chosen (licensed and the GPUs do vGPU).</summary>
        public static bool CanEditTypes(List<PGPU> pgpus) =>
            pgpus.Count > 0 && !Helpers.FeatureForbidden(pgpus[0].Connection, Host.RestrictVgpu) && pgpus[0].HasVGpu();

        /// <summary>GpuPage: the placement policy panel is shown when the pool can run vGPUs.</summary>
        public static bool ShowPlacementPolicy(IXenConnection c) => Helpers.VGpuCapability(c);

        /// <summary>GpuPlacementPolicyPanel.PopulatePage: the GPU groups' common algorithm (unknown when mixed).</summary>
        public static allocation_algorithm PlacementPolicy(IXenConnection c)
        {
            var algorithms = c.Cache.GPU_groups.Select(g => g.allocation_algorithm).Distinct().ToList();
            return algorithms.Count == 1 ? algorithms[0] : allocation_algorithm.unknown;
        }

        /// <summary>GpuPlacementPolicyPanel.ToString.</summary>
        public static string PolicyText(allocation_algorithm a) => a switch
        {
            allocation_algorithm.breadth_first => Messages.GPU_PLACEMENT_POLICY_MAX_PERFORMANCE_DESCRIPTION,
            allocation_algorithm.depth_first => Messages.GPU_PLACEMENT_POLICY_MAX_DENSITY_DESCRIPTION,
            _ => Messages.GPU_PLACEMENT_POLICY_MIXED_DESCRIPTION
        };

        /// <summary>PoolGpuEditPage: SetGpuPlacementPolicyAction (null when unchanged).</summary>
        public static AsyncAction PlacementAction(IXenConnection c, allocation_algorithm chosen) =>
            chosen == PlacementPolicy(c) || chosen == allocation_algorithm.unknown ? null
                : new SetGpuPlacementPolicyAction(Helpers.GetPoolOfOne(c), chosen);

        public sealed record VgpuTypeRow(XenRef<VGPU_type> Ref, string Name, long MaxPerGpu, long VideoRam, bool Enabled, bool InUse);

        /// <summary>GpuConfiguration.PopulateGridView: each supported type, enabled or not, and in use (cannot be disallowed).</summary>
        public static List<VgpuTypeRow> TypeRows(List<PGPU> pgpus)
        {
            var first = pgpus[0];
            return first.supported_VGPU_types
                .Select(r => (Ref: r, Type: first.Connection.Resolve(r)))
                .Where(x => x.Type != null)
                .OrderByDescending(x => x.Type)
                .Select(x => new VgpuTypeRow(x.Ref,
                    x.Type.IsPassthrough() ? Messages.VGPU_PASSTHRU_TOSTRING : x.Type.model_name,
                    first.supported_VGPU_max_capacities != null && first.supported_VGPU_max_capacities.TryGetValue(x.Ref, out var max) ? max : 0,
                    x.Type.framebuffer_size,
                    first.enabled_VGPU_types.Any(e => e.opaque_ref == x.Ref.opaque_ref),
                    pgpus.Any(p => p.Connection.ResolveAll(p.resident_VGPUs).Any(v => v.type.opaque_ref == x.Ref.opaque_ref))))
                .ToList();
        }

        /// <summary>The vGPUs running on these GPUs.</summary>
        public static int RunningVgpus(List<PGPU> pgpus) => pgpus.Sum(p => p.resident_VGPUs.Count);

        /// <summary>
        /// GpuConfiguration.SaveChanges: only the boxes the user changed (type ref, now allowed) are applied to each GPU's
        /// current allowed types, so changes made elsewhere meanwhile stay; null when nothing changed.
        /// </summary>
        public static AsyncAction ConfigureAction(List<PGPU> pgpus, IDictionary<string, bool> changes)
        {
            if (changes.Count == 0)
                return null;
            var updated = new Dictionary<PGPU, List<XenRef<VGPU_type>>>();
            foreach (var p in pgpus)
            {
                var list = new List<XenRef<VGPU_type>>(p.enabled_VGPU_types);
                foreach (var (typeRef, allowed) in changes)
                {
                    var has = list.Any(r => r.opaque_ref == typeRef);
                    if (has && !allowed)
                        list.RemoveAll(r => r.opaque_ref == typeRef);
                    else if (!has && allowed)
                        list.Add(new XenRef<VGPU_type>(typeRef));
                }
                updated[p] = list;
            }
            return new VgpuConfigurationAction(updated, pgpus[0].Connection);
        }

        // ---------- USB ----------

        /// <summary>UsbPage: the device is attached to a VM (its USB group has a VUSB).</summary>
        public static bool IsAttached(PUSB pusb)
        {
            var group = pusb.Connection.Resolve(pusb.USB_group);
            return group?.VUSBs != null && group.VUSBs.Count > 0;
        }

        /// <summary>UsbPage.HostUsbRow: the VM a pass-through device is attached to (null if none).</summary>
        public static VM AttachedVm(PUSB pusb)
        {
            var group = pusb.Connection.Resolve(pusb.USB_group);
            var vusb = group?.VUSBs != null && group.VUSBs.Count > 0 ? pusb.Connection.Resolve(group.VUSBs[0]) : null;
            return vusb == null ? null : pusb.Connection.Resolve(vusb.VM);
        }

        public static List<PUSB> Pusbs(Host host) => host.Connection.ResolveAll(host.PUSBs).OrderBy(p => p.path).ToList();

        /// <summary>UsbPage.dataGridViewUsbList_SelectionChanged: the button text and whether it can be used.</summary>
        public static (string Text, bool Enabled) PassthroughButton(PUSB pusb)
        {
            if (pusb == null)
                return (Messages.USBLIST_ENABLE_PASSTHROUGH_HOTKEY, false);
            if (pusb.passthrough_enabled)
                return (Messages.USBLIST_DISABLE_PASSTHROUGH_HOTKEY, !IsAttached(pusb));
            return (Messages.USBLIST_ENABLE_PASSTHROUGH_HOTKEY, true);
        }

        /// <summary>UsbPage.buttonPassthrough_Click: the question and its Yes button.</summary>
        public static (string Text, string Yes) PassthroughConfirmation(PUSB pusb) => pusb.passthrough_enabled
            ? (Messages.DIALOG_USB_USAGE_DISABLE_PASSTHROUGH, Messages.DIALOG_USB_USAGE_OKBUTTON_DISABLE)
            : (Messages.DIALOG_USB_USAGE_ENABLE_PASSTHROUGH, Messages.DIALOG_USB_USAGE_OKBUTTON_ENABLE);

        /// <summary>SetUsbPassthroughAction to the state the user confirmed.</summary>
        public static AsyncAction PassthroughAction(PUSB pusb, bool enable) => new SetUsbPassthroughAction(pusb, enable);
    }
}
