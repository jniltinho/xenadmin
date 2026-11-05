using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Actions.VMActions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>Enablement rules of the VM storage commands (as the WinForms commands' CanRun).</summary>
    public static class StorageRules
    {
        /// <summary>SrPicker (new disk / copy / new VM): visible, writable, plugged, not ISO, reachable from the home host.</summary>
        public static bool CanHoldNewDisk(SR sr, Host home) =>
            sr != null && sr.SupportsVdiCreate() && sr.HasPBDs() && !sr.IsBroken() && sr.content_type != SR.Content_Type_ISO
            && sr.Show(false) && (home == null || sr.CanBeSeenFrom(home));

        public static List<SR> NewDiskTargets(IXenConnection c, Host home) =>
            c.Cache.SRs.Where(sr => CanHoldNewDisk(sr, home)).OrderBy(sr => sr.Name()).ToList();

        public static bool CanActivate(VM vm, VBD vbd) =>
            vm.IsRealVm() && vm.power_state == vm_power_state.Running && !vbd.Locked && !vbd.currently_attached
            && vbd.allowed_operations.Contains(vbd_operations.plug);

        public static bool CanDeactivate(VM vm, VBD vbd) =>
            vm.IsRealVm() && vm.power_state == vm_power_state.Running && !vbd.Locked && vbd.currently_attached
            && vbd.allowed_operations.Contains(vbd_operations.unplug);

        /// <summary>DetachVirtualDiskCommand.CanRun: every attached VBD of this VM must be unpluggable.</summary>
        public static bool CanDetach(VM vm, VDI vdi) =>
            vdi != null && !vdi.Locked && vdi.Connection.ResolveAll(vdi.VBDs)
                .Where(b => b.VM.opaque_ref == vm.opaque_ref && b.currently_attached).All(b => CanDeactivate(vm, b));

        /// <summary>DeleteVirtualDiskCommand.CanRun (AllowRunningVMDelete = false).</summary>
        public static bool CanDelete(VDI vdi)
        {
            var sr = vdi?.Connection.Resolve(vdi.SR);
            if (sr == null || vdi.Locked || sr.Physical() || sr.IsToolsSR() || vdi.IsUsedByHA())
                return false;
            var vbds = vdi.Connection.ResolveAll(vdi.VBDs);
            if (vbds.Count > 1)
                return false;
            foreach (var vbd in vbds)
            {
                var vm = vdi.Connection.Resolve(vbd.VM);
                if (vdi.type == vdi_type.system && vm?.power_state == vm_power_state.Running)
                    return false;
                if (vbd.Locked || vbd.currently_attached)
                    return false;
            }
            return true;
        }


        /// <summary>One item of the DVD drive's ISO list (CDChanger): the ISO, or null for "&lt;empty&gt;".</summary>
        public sealed class IsoChoice
        {
            public VDI Vdi;
            public string Text;
            public override string ToString() => Text;
        }

        /// <summary>"&lt;empty&gt;" plus ISOs of ISO SRs visible from the VM's home server.</summary>
        public static List<IsoChoice> IsoChoices(VM vm)
        {
            var list = new List<IsoChoice> { new IsoChoice { Text = "<empty>" } };
            var home = vm.Home();
            // ISODropDownBox: the tools SR is hidden from Stockholm (XCP-ng 8.x) on, the guest agent comes from the distro
            foreach (var sr in vm.Connection.Cache.SRs.Where(s => s.content_type == SR.Content_Type_ISO && s.Show(false)
                         && !(s.IsToolsSR() && XenAdmin.Core.Helpers.StockholmOrGreater(vm.Connection))
                         && (home == null || s.CanBeSeenFrom(home))).OrderBy(s => s.Name()))
                foreach (var vdi in vm.Connection.ResolveAll(sr.VDIs).Where(v => v.Show(false)).OrderBy(v => v.Name()))
                    list.Add(new IsoChoice { Vdi = vdi, Text = $"{sr.Name()}: {vdi.Name()}" });
            return list;
        }
    }
}
