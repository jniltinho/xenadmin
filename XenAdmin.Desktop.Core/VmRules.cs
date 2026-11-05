using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Actions.VMActions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;
using XenCenterLib;

namespace XenAdmin.Desktop
{
    /// <summary>Rules of the VM commands and pages (CPU topology, snapshots, Copy/Delete VM).</summary>
    public static class VmRules
    {
        /// <summary>PropertiesCommand.CanRun for the objects that have a Properties dialog here.</summary>
        public static bool CanEditProperties(IXenObject o) =>
            o != null && o.Connection != null && o.Connection.IsConnected && !o.Locked && (o is VM || o is Host || o is Pool || o is SR || o is VDI || o is XenAPI.Network);

        public sealed record SnapshotModes(bool Disk, bool QuiesceVisible, bool Quiesce, bool Memory, string MemoryInfo);

        /// <summary>VmSnapshotDialog: which snapshot modes the VM allows, and why disk+memory is not offered.</summary>
        public static SnapshotModes SnapshotModesFor(VM vm)
        {
            var quiesceVisible = !Helpers.QuebecOrGreater(vm.Connection);
            var checkpointForbidden = Helpers.FeatureForbidden(vm, Host.RestrictCheckpoint);
            var info = checkpointForbidden ? Messages.FIELD_DISABLED
                : vm.power_state != vm_power_state.Running ? Messages.INFO_DISKMEMORY_MODE_POWER_STATE
                : !vm.GetVirtualizationStatus(out _).HasFlag(VM.VirtualizationStatus.IoDriversInstalled)
                    ? (vm.HasNewVirtualizationStates() ? Messages.INFO_DISKMEMORY_MODE_NO_IO_DRIVERS : string.Format(Messages.INFO_DISKMEMORY_MODE_NO_TOOLS, BrandManager.VmTools))
                : Messages.INFO_DISKMEMORY_MODE_MISC;
            return new SnapshotModes(
                vm.allowed_operations.Contains(vm_operations.snapshot),
                quiesceVisible,
                quiesceVisible && vm.allowed_operations.Contains(vm_operations.snapshot_with_quiesce) && !Helpers.FeatureForbidden(vm, Host.RestrictVss),
                vm.allowed_operations.Contains(vm_operations.checkpoint) && !checkpointForbidden,
                info);
        }

        /// <summary>CopyVMDialog.OnLoad: which of fast clone / full copy are available.</summary>
        public static (bool Clone, bool Copy) CopyModes(VM vm)
        {
            var allowCopy = !vm.is_a_template || vm.allowed_operations.Contains(vm_operations.copy);
            var hasDisk = vm.HasAtLeastOneDisk();
            return (!allowCopy || vm.AnyDiskFastClonable() || !hasDisk, allowCopy && hasDisk);
        }

        /// <summary>CopyVMDialog.MoveButton_Click; null when a full copy has no target SR (never falls back to a clone).</summary>
        public static AsyncAction CopyAction(VM vm, bool fullCopy, SR sr, string name, string description)
        {
            if (string.IsNullOrWhiteSpace(name)) return null;
            if (!fullCopy) return new VMCloneAction(vm, name.Trim(), description ?? "");
            return sr == null ? null : new VMCopyAction(vm, vm.GetStorageHost(false), sr, name.Trim(), description ?? "");
        }

        /// <summary>AffinityPicker: "no home server" needs shared storage, unless the VM has no home server already.</summary>
        public static bool CanHaveNoHomeServer(VM vm) =>
            Helpers.HasFullyConnectedSharedStorage(vm.Connection) || vm.Connection.Resolve(vm.affinity) == null;

        /// <summary>HomeServerEditPage.SaveSettings: VM.set_affinity, the null reference for no home server.</summary>
        public static AsyncAction SetHomeServer(VM vm, Host host) =>
            new DelegatedAsyncAction(vm.Connection, Messages.ACTION_CHANGE_HOME_SERVER, string.Format(Messages.ACTION_CHANGING_HOME_SERVER_FOR, vm), null,
                session => VM.set_affinity(session, vm.opaque_ref, host?.opaque_ref ?? Helper.NullOpaqueRef), true, "vm.set_affinity");

        /// <summary>BootOptionsEditPage (PV): the CD is bootable, or only the lowest numbered disk (CA-47457).</summary>
        public static AsyncAction PvBootableAction(VM vm, bool fromCd) =>
            new DelegatedAsyncAction(vm.Connection, "Change VBDs bootable", "Change VBDs bootable", null, session =>
            {
                var vbds = vm.Connection.ResolveAll(vm.VBDs);
                if (fromCd)
                {
                    foreach (var vbd in vbds)
                        VBD.set_bootable(session, vbd.opaque_ref, vbd.IsCDROM());
                    return;
                }
                vbds.Sort((a, b) => a.userdevice == "xvda" ? -1 : b.userdevice == "xvda" ? 1 : StringUtility.NaturalCompare(a.userdevice, b.userdevice));
                var found = false;
                foreach (var vbd in vbds)
                {
                    var bootable = !found && vbd.type == vbd_type.Disk;
                    found |= bootable;
                    VBD.set_bootable(session, vbd.opaque_ref, bootable);
                }
            }, true, "VBD.set_bootable");

        /// <summary>CPUTopologyComboBox.GetTopologies: (sockets, cores) splits of the vCPUs.</summary>
        public static List<(long Sockets, long Cores)> ValidTopologies(long vcpus, long maxCoresPerSocket)
        {
            var limit = maxCoresPerSocket > 0 ? Math.Min(vcpus, maxCoresPerSocket) : vcpus;
            var result = new List<(long, long)>();
            for (long cores = 1; cores <= limit; cores++)
                if (vcpus % cores == 0 && vcpus / cores <= VM.MAX_SOCKETS)
                    result.Add((vcpus / cores, cores));
            return result;
        }

        /// <summary>SnapshotsPage.GetSnapshotSize: disks plus the suspend VDI.</summary>
        public static long SnapshotSize(VM snapshot)
        {
            long size = 0;
            foreach (var vbd in snapshot.Connection.ResolveAll(snapshot.VBDs).Where(b => b.type == vbd_type.Disk))
            {
                var vdi = snapshot.Connection.Resolve(vbd.VDI);
                if (vdi != null && vdi.physical_utilisation > 0)
                    size += vdi.physical_utilisation;
            }
            var suspend = snapshot.Connection.Resolve(snapshot.suspend_VDI);
            if (suspend != null && suspend.physical_utilisation > 0)
                size += suspend.physical_utilisation;
            return size;
        }

        public static string SnapshotTypeText(VM snapshot) =>
            snapshot.power_state == vm_power_state.Suspended ? Messages.DISKS_AND_MEMORY : Messages.DISKS_ONLY;

        /// <summary>CopyVMCommand: a halted, unlocked VM that can be cloned or copied.</summary>
        public static bool CanCopy(VM vm) =>
            vm != null && !vm.is_a_snapshot && !vm.is_a_template && !vm.Locked && vm.power_state == vm_power_state.Halted
            && (vm.allowed_operations.Contains(vm_operations.clone) || vm.allowed_operations.Contains(vm_operations.copy));

        /// <summary>DeleteVMCommand: destroy allowed (the VM is shut down).</summary>
        public static bool CanDeleteVm(VM vm) =>
            vm != null && !vm.is_a_snapshot && !vm.Locked && vm.allowed_operations.Contains(vm_operations.destroy);
    }
}
