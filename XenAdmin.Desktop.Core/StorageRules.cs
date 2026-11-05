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

        /// <summary>
        /// DeleteVirtualDiskCommand.CanRun (AllowRunningVMDelete = false); the SR Storage tab sets AllowMultipleVBDDelete
        /// (a disk shown with several VMs can be deleted there).
        /// </summary>
        public static bool CanDelete(VDI vdi, bool allowMultipleVbds = false)
        {
            var sr = vdi?.Connection.Resolve(vdi.SR);
            if (sr == null || vdi.Locked || sr.Physical() || sr.IsToolsSR() || vdi.IsUsedByHA())
                return false;
            var vbds = vdi.Connection.ResolveAll(vdi.VBDs);
            if (vbds.Count > 1 && !allowMultipleVbds)
                return false;
            foreach (var vbd in vbds)
            {
                var vm = vdi.Connection.Resolve(vbd.VM);
                if (vdi.type == vdi_type.system && vm?.power_state == vm_power_state.Running)
                    return false;
                if (vbd.Locked || vbd.currently_attached)
                    return false;
            }
            return sr.HBALunPerVDI() || vdi.allowed_operations.Contains(vdi_operations.destroy);
        }

        /// <summary>SrStoragePage.RefreshButtons: Rescan, and its tooltip when disabled (null when none).</summary>
        public static (bool Enabled, string Tooltip) CanRescan(SR sr, bool beingScanned)
        {
            if (sr == null || sr.Locked)
                return (false, null);
            if (sr.IsDetached())
                return (false, Messages.SR_DETACHED);
            if (beingScanned)
                return (false, Messages.SCAN_IN_PROGRESS_TOOLTIP);
            return (true, null);
        }

        /// <summary>VMStoragePage: the VM's disks (no DVD or floppy drives, nor VBDs whose disk is gone), by position.</summary>
        public static List<(VBD Vbd, VDI Vdi)> VmDisks(VM vm) =>
            vm.Connection.ResolveAll(vm.VBDs).Where(b => b.type != vbd_type.CD && b.type != vbd_type.Floppy)
                .OrderBy(b => b.userdevice, StringComparer.Ordinal)
                .Select(b => (b, vm.Connection.Resolve(b.VDI))).Where(d => d.Item2 != null).ToList();

        /// <summary>PhysicalStoragePage: the SRs of a pool, or the shared ones and its own of a host.</summary>
        public static List<SR> PhysicalSrs(IXenObject o)
        {
            var host = o as Host;
            return o.Connection.Cache.SRs.Where(s => s.Show(false) && !s.IsToolsSR())
                .Where(s => host == null || s.shared || s.Home()?.opaque_ref == host.opaque_ref).OrderBy(s => s.Name()).ToList();
        }

        /// <summary>TrimSRCommand.CanRun: Reclaim freed space needs an SR that supports trim and is attached to a host.</summary>
        public static bool CanTrim(SR sr) => sr != null && sr.SupportsTrim() && sr.GetFirstAttachedStorageHost() != null;

        /// <summary>SrStoragePage: Add needs an unlocked SR; Properties one unlocked disk that is not a snapshot.</summary>
        public static bool CanAddDisk(SR sr) => sr != null && !sr.Locked;

        public static bool CanEditDisk(SR sr, VDI vdi) => sr != null && !sr.Locked && vdi != null && !vdi.is_a_snapshot && !vdi.Locked;

        /// <summary>The disks of the SR Storage tab (SrStoragePage.GetCurrentData).</summary>
        public static List<VDI> SrDisks(SR sr) =>
            sr.Connection.ResolveAll(sr.VDIs).Where(v => v.Show(false) && !v.IsAnIntermediateStorageMotionSnapshot()).OrderBy(v => v.Name()).ToList();

        /// <summary>DeleteVirtualDiskCommand.ConfirmationDialogText for one disk (title, text).</summary>
        public static (string Title, string Text) DeleteConfirmation(VDI vdi) => vdi.VDIType() switch
        {
            VDI.FriendlyType.SNAPSHOT => (Messages.MESSAGEBOX_DELETE_SNAPSHOT_TITLE, Messages.MESSAGEBOX_DELETE_SNAPSHOT),
            VDI.FriendlyType.ISO => (Messages.MESSAGEBOX_DELETE_ISO_TITLE, string.Format(Messages.MESSAGEBOX_DELETE_ISO, vdi.Name())),
            VDI.FriendlyType.SYSTEM_DISK => (Messages.MESSAGEBOX_DELETE_SYS_DISK_TITLE, Messages.MESSAGEBOX_DELETE_SYS_DISK),
            _ => (Messages.MESSAGEBOX_DELETE_VD_TITLE, Messages.MESSAGEBOX_DELETE_VD)
        };

        /// <summary>DeleteVirtualDiskCommand.getDestroyVDIAction: a snapshot disk deletes its VM snapshots, others the disk.</summary>
        public static List<AsyncAction> DeleteActions(VDI vdi)
        {
            if (vdi.is_a_snapshot && vdi.GetVMs().Count >= 1)
                return vdi.GetVMs().Where(vm => vm.is_a_snapshot).Select(vm => (AsyncAction)new VMSnapshotDeleteAction(vm)).ToList();
            return new List<AsyncAction> { new DestroyDiskAction(vdi) };
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
