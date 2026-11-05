using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions.VMActions;
using XenAdmin.Core;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>VM > Migrate to Server (MigrateVMToolStripMenuItem, VMOperationHostCommand) inside a pool.</summary>
    public static class MigrateRules
    {
        /// <summary>MigrateVmToServerCommand.CanRun.</summary>
        public static bool CanMigrate(VM vm) =>
            vm?.Connection != null && vm.Connection.IsConnected && !vm.is_a_template && !vm.Locked
            && !Helpers.FeatureForbidden(vm.Connection, Host.RestrictIntraPoolMigrate)
            && vm.allowed_operations != null && vm.allowed_operations.Contains(vm_operations.pool_migrate);

        /// <summary>
        /// The submenu's hosts with the reason each one cannot take the VM (null when it can), as VmCanBootOnHost.
        /// Calls VM.assert_can_boot_here per host, so run it off the UI thread.
        /// </summary>
        public static List<(Host Host, string Reason)> Targets(VM vm, Session session) =>
            vm.Connection.Cache.Hosts.OrderBy(h => h.Name(), StringComparer.CurrentCultureIgnoreCase)
                .Select(h => (h, CannotBootReason(vm, h, session))).ToList();

        public static string CannotBootReason(VM vm, Host host, Session session)
        {
            if (vm.power_state == vm_power_state.Running && vm.Connection.Resolve(vm.resident_on) is Host resident)
            {
                if (host.opaque_ref == resident.opaque_ref)
                    return Messages.HOST_MENU_CURRENT_SERVER;
                if (Helpers.ProductVersionCompare(Helpers.HostProductVersion(host), Helpers.HostProductVersion(resident)) < 0)
                    return Messages.OLDER_THAN_CURRENT_SERVER;
            }
            if ((vm.power_state == vm_power_state.Running || vm.power_state == vm_power_state.Suspended)
                && vm.last_boot_CPU_flags != null && vm.last_boot_CPU_flags.TryGetValue("vendor", out var vendor)
                && host.cpu_info != null && host.cpu_info.TryGetValue("vendor", out var hostVendor) && vendor != hostVendor)
                return FriendlyErrorNames.VM_INCOMPATIBLE_WITH_THIS_HOST;
            try
            {
                VM.assert_can_boot_here(session, vm.opaque_ref, host.opaque_ref);
                return null;
            }
            catch (Failure f)
            {
                if (f.ErrorDescription.Count > 2 && f.ErrorDescription[0] == Failure.VM_REQUIRES_SR
                    && host.Connection.Resolve(new XenRef<SR>(f.ErrorDescription[2]))?.content_type == SR.Content_Type_ISO)
                    return Messages.MIGRATE_PLEASE_EJECT_YOUR_CD;
                return f.ShortMessage;
            }
            catch (Exception)
            {
                return Messages.HOST_MENU_UNKNOWN_ERROR;
            }
        }

        /// <summary>VMOperationHomeServerCommand: the first item targets the VM's home server (affinity), if any.</summary>
        public static (Host Host, string Reason) HomeServerTarget(VM vm, Session session) =>
            vm.Connection.Resolve(vm.affinity) is Host home ? (home, CannotBootReason(vm, home, session)) : (null, Messages.NO_HOME_SERVER);

        public static string HomeServerMenuText(string reason) =>
            reason == null ? Messages.HOME_SERVER_MENU_ITEM : string.Format(Messages.MAINWINDOW_CONTEXT_REASON, Messages.HOME_SERVER_MENU_ITEM, reason);

        /// <summary>"xcp-lab-2 (Current server)" — MAINWINDOW_CONTEXT_REASON.</summary>
        public static string MenuText(Host host, string reason) =>
            reason == null ? host.Name() : string.Format(Messages.MAINWINDOW_CONTEXT_REASON, host.Name(), reason);

        public static VMMigrateAction Migrate(VM vm, Host host) => new VMMigrateAction(vm, host);
    }
}
