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

        /// <summary>MoveVirtualDiskCommand.CanBeMoved: a disk of halted VMs only, not a snapshot, HA or CBT disk.</summary>
        public static bool CanMove(VDI vdi)
        {
            if (vdi == null || vdi.is_a_snapshot || vdi.Locked || vdi.IsHaType() || vdi.cbt_enabled)
                return false;
            var sr = vdi.Connection.Resolve(vdi.SR);
            return sr != null && !sr.HBALunPerVDI() && vdi.GetVMs().All(vm => vm.power_state == vm_power_state.Halted);
        }

        /// <summary>MigrateVirtualDiskCommand.CanBeMigrated: an attached disk on an SR that supports storage migration.</summary>
        public static bool CanMigrate(VDI vdi)
        {
            if (vdi == null || vdi.is_a_snapshot || vdi.Locked || vdi.IsHaType() || vdi.cbt_enabled)
                return false;
            if (vdi.Connection.ResolveAll(vdi.VBDs).Count == 0)
                return false;
            var sr = vdi.Connection.Resolve(vdi.SR);
            return sr != null && !sr.HBALunPerVDI() && sr.SupportsStorageMigration();
        }

        /// <summary>
        /// MoveVirtualDiskDialog.MoveMigrateCommand: live storage migration when the disk allows it and the license does not
        /// forbid it, else an offline move; null when neither is possible.
        /// </summary>
        public static bool? UseMigrate(VDI vdi) =>
            CanMigrate(vdi) && !Helpers.FeatureForbidden(vdi.Connection, Host.RestrictCrossPoolMigrate) ? true
            : CanMove(vdi) ? false
            : (bool?)null;

        public sealed record SrTarget(SR Sr, bool Enabled, string Text);

        /// <summary>
        /// SrPicker for Move / Migrate (SrPickerMoveItem, SrPickerMigrateItem): the visible SRs that can hold disks, each
        /// enabled or with the reason it cannot take the disk.
        /// </summary>
        public static List<SrTarget> MoveTargets(VDI vdi, bool migrate)
        {
            var vdis = new[] { vdi };
            var result = new List<SrTarget>();
            foreach (var sr in vdi.Connection.Cache.SRs.Where(x => !x.HBALunPerVDI() && x.SupportsVdiCreate() && x.Show(false)).OrderBy(x => x.Name()))
            {
                string reason = null;
                if (vdi.SR.opaque_ref == sr.opaque_ref)
                    reason = Messages.CURRENT_LOCATION;
                else if (sr.IsDetached())
                    reason = Messages.SR_DETACHED;
                else if (migrate && sr.IsLocalSR())
                {
                    var homes = vdi.GetVMs().Select(vm => (vm, home: vm.Home())).Where(x => x.home != null).ToList();
                    var unseen = homes.FirstOrDefault(x => !sr.CanBeSeenFrom(x.home));
                    if (homes.Count == 0)
                        reason = Messages.SR_IS_LOCAL;
                    else if (unseen.vm != null)
                        reason = unseen.vm.power_state == vm_power_state.Running
                            ? Messages.SRPICKER_ERROR_LOCAL_SR_MUST_BE_RESIDENT_HOSTS
                            : string.Format(Messages.SR_CANNOT_BE_SEEN, Helpers.GetName(unseen.home));
                }
                if (reason == null && migrate && !sr.SupportsStorageMigration())
                    reason = Messages.UNSUPPORTED_SR_TYPE;
                if (reason == null && !sr.CanFitDisks(out var cannotFit, vdis))
                    reason = cannotFit;
                result.Add(new SrTarget(sr, reason == null, reason ?? string.Format(Messages.SRPICKER_DISK_FREE,
                    Util.DiskSizeString(sr.FreeSpace(), 2), Util.DiskSizeString(sr.physical_size, 2))));
            }
            return result;
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
