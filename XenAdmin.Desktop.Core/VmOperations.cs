using System.Collections.Generic;
using XenAdmin.Actions;
using XenAdmin.Actions.VMActions;
using XenAdmin.Core;
using XenAPI;

namespace XenAdmin.Desktop
{
    public enum VmOp { Start, Shutdown, Reboot, Suspend, Resume, ForceShutdown, ForceReboot }

    /// <summary>Maps UI operations to XAPI allowed_operations and to the existing XenModel actions.</summary>
    public static class VmOperations
    {
        public static readonly IReadOnlyDictionary<VmOp, vm_operations> Required = new Dictionary<VmOp, vm_operations>
        {
            [VmOp.Start] = vm_operations.start,
            [VmOp.Shutdown] = vm_operations.clean_shutdown,
            [VmOp.Reboot] = vm_operations.clean_reboot,
            [VmOp.Suspend] = vm_operations.suspend,
            [VmOp.Resume] = vm_operations.resume,
            [VmOp.ForceShutdown] = vm_operations.hard_shutdown,
            [VmOp.ForceReboot] = vm_operations.hard_reboot,
        };

        public static bool IsAllowed(VM vm, VmOp op) =>
            vm != null && vm.allowed_operations != null && vm.allowed_operations.Contains(Required[op]);

        /// <summary>ShutDownVMCommand, RebootVMCommand, SuspendVMCommand and the forced ones always confirm; Start and Resume do not.</summary>
        public static bool NeedsConfirmation(VmOp op) => op != VmOp.Start && op != VmOp.Resume;

        /// <summary>The commands' ConfirmationDialogTitle / ConfirmationDialogText (HA-protected VMs get the HA warning).</summary>
        public static (string Title, string Text) Confirmation(VM vm, VmOp op)
        {
            var ha = vm.HAIsProtected();
            return op switch
            {
                VmOp.Shutdown => (Messages.CONFIRM_SHUTDOWN_VM_TITLE, ha ? Messages.HA_CONFIRM_SHUTDOWN_VM : Messages.CONFIRM_SHUTDOWN_VM),
                VmOp.Reboot => (Messages.CONFIRM_REBOOT_VM_TITLE, Messages.CONFIRM_REBOOT_VM),
                VmOp.Suspend => (Messages.CONFIRM_SUSPEND_VM_TITLE, ha ? Messages.HA_CONFIRM_SUSPEND_VM : Messages.CONFIRM_SUSPEND_VM),
                VmOp.ForceShutdown => (Messages.CONFIRM_FORCESHUTDOWN_VM_TITLE, Messages.CONFIRM_FORCESHUTDOWN_VM_NO_CANCEL_TASKS),
                VmOp.ForceReboot => (Messages.CONFIRM_FORCEREBOOT_VM_TITLE, Messages.CONFIRM_FORCEREBOOT_VM_NO_CANCEL_TASKS),
                _ => (null, null)
            };
        }

        public static AsyncAction Create(VM vm, VmOp op)
        {
            switch (op)
            {
                case VmOp.Start:
                    // never unprotect an HA-protected VM behind the user's back
                    return new VMStartAction(vm, null, null) { AllowHaUnprotect = false };
                case VmOp.Resume:
                    return new VMResumeAction(vm, null, null) { AllowHaUnprotect = false };
                case VmOp.Shutdown: return new VMCleanShutdown(vm);
                case VmOp.Reboot: return new VMCleanReboot(vm);
                case VmOp.Suspend: return new VMSuspendAction(vm);
                case VmOp.ForceShutdown: return new VMHardShutdown(vm);
                case VmOp.ForceReboot: return new VMHardReboot(vm);
            }
            return null;
        }
    }
}
