using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>VM > Manage vTPMs: VtpmCommand and VtpmManagementDialog / VtpmManagementPage.</summary>
    public static class VtpmRules
    {
        /// <summary>VtpmCommand.CanRunCore: a real HVM VM booting with UEFI, where vTPMs are supported and licensed.</summary>
        public static bool CanManage(VM vm) =>
            vm != null && vm.IsRealVm() && vm.Connection.IsConnected && !Helpers.FeatureForbidden(vm, Host.RestrictVtpm)
            && Helpers.XapiEqualOrGreater_22_26_0(vm.Connection) && vm.IsHVM() && vm.IsDefaultBootModeUefi();

        /// <summary>VtpmCommand.GetCantRunReasonCore: supported here, but the VM does not boot with UEFI.</summary>
        public static string Blocker(VM vm) =>
            vm != null && vm.IsRealVm() && !Helpers.FeatureForbidden(vm, Host.RestrictVtpm) && Helpers.XapiEqualOrGreater_22_26_0(vm.Connection)
            && !(vm.IsHVM() && vm.IsDefaultBootModeUefi())
                ? Messages.COMMAND_VTPM_DISABLED_NON_UEFI : null;

        public static List<VTPM> Vtpms(VM vm) => vm.Connection.ResolveAll(vm.VTPMs);

        /// <summary>VM.CanAddVtpm: null when a vTPM can be added, else the reason.</summary>
        public static string AddBlocker(VM vm) => vm.CanAddVtpm(out var reason) ? null : reason;

        /// <summary>VtpmManagementDialog.CanRemoveVtpm: null when this vTPM can be removed, else the reason.</summary>
        public static string RemoveBlocker(VM vm, VTPM vtpm)
        {
            if (Helpers.XapiEqualOrGreater_23_11_0(vtpm.Connection))
                return vtpm.allowed_operations.Contains(vtpm_operations.destroy) ? null : Messages.VTPM_OPERATION_DISALLOWED_REMOVE;
            return vm.CanRemoveVtpm(out var reason) ? null : reason;
        }

        public static string RemoveWarning => Messages.VTPM_REMOVE_WARNING;

        public static AsyncAction AddAction(VM vm) => new NewVtpmAction(vm.Connection, vm);

        public static AsyncAction RemoveAction(VM vm, VTPM vtpm) => new RemoveVtpmAction(vtpm.Connection, vtpm, vm);

        /// <summary>VtpmManagementPage: Protected and Unique.</summary>
        public static (string Protected, string Unique) Properties(VTPM vtpm) =>
            (vtpm.is_protected ? Messages.YES : Messages.NO, vtpm.is_unique ? Messages.YES : Messages.NO);
    }
}
