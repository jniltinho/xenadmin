using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>Server commands (ShutDownHostCommand, RebootHostCommand, HostMaintenanceModeCommand).</summary>
    public static class HostRules
    {
        /// <summary>ShutDownHostCommand.CanRun; busy is HelpersGUI.HasActiveHostAction (a server action of ours is running).</summary>
        public static bool CanShutDown(Host host, bool busy = false) => host != null && host.Connection.IsConnected && host.IsLive() && !host.Locked && !busy;
        public static bool CanReboot(Host host, bool busy = false) => CanShutDown(host, busy);

        /// <summary>HostMaintenanceModeCommand.CanRunCore: a live host of a connected pool.</summary>
        public static bool CanToggleMaintenance(Host host, bool busy = false) =>
            host != null && host.Connection.IsConnected && host.IsLive() && !host.Locked && !busy && Helpers.GetCoordinator(host.Connection) != null;

        public static bool InMaintenance(Host host) => host != null && !host.enabled;

        public static string MaintenanceText(Host host) =>
            InMaintenance(host) ? Messages.EXIT_MAINTENANCE_MODE_CONTEXT_MENU : Messages.ENTER_MAINTENANCE_MODE_CONTEXT_MENU;

        private static bool HasRunningVms(Host host) =>
            host.Connection.ResolveAll(host.resident_VMs).Any(v => !v.is_control_domain && v.power_state == vm_power_state.Running);

        /// <summary>The confirmation of ShutDownHostCommand / RebootHostCommand for one server (the HCI warning first).</summary>
        public static (string Title, string Text) Confirm(Host host, bool reboot)
        {
            var hci = host.Connection.ResolveAll(host.resident_VMs).Exists(v => v.HciWarnBeforeShutdown());
            var text = hci
                ? string.Format(reboot ? Messages.CONFIRM_REBOOT_HCI_WARN_SERVER : Messages.CONFIRM_SHUTDOWN_HCI_WARN_SERVER, host.Name())
                : HasRunningVms(host)
                ? string.Format(reboot ? Messages.CONFIRM_REBOOT_SERVER : Messages.CONFIRM_SHUTDOWN_SERVER, host.Name())
                : string.Format(reboot ? Messages.CONFIRM_REBOOT_SERVER_NO_VMS : Messages.CONFIRM_SHUTDOWN_SERVER_NO_VMS, host.Name());
            if (!reboot && host.IsCoordinator() && host.Connection.Cache.Hosts.Length > 1)
                text += "\n\n" + string.Format(Messages.SHUT_DOWN_POOL_COORDINATOR_SINGLE, host.Name());
            return (reboot ? Messages.CONFIRM_REBOOT_SERVER_TITLE : Messages.CONFIRM_SHUTDOWN_SERVER_TITLE, text);
        }

        // shortcut: no failures-to-tolerate dialog (AddHostToPoolCommand.NtolDialog) yet. True cancels the action rather
        // than lowering the pool's HA protection without consent; on enable, false leaves the tolerance as it is.
        private static bool KeepNtol(HostAbstractAction a, Pool p, long current, long target) => true;
        private static bool KeepNtolOnEnable(Pool p, Host h, long current, long target) => false;

        public static AsyncAction ShutDown(Host host) => new ShutdownHostAction(host, KeepNtol);
        public static AsyncAction Reboot(Host host) => new RebootHostAction(host, KeepNtol);

        /// <summary>
        /// EnterMaintenanceMode: refused for the HA coordinator; evacuates with XAPI's own placement.
        /// shortcut: the EvacuateHostDialog (choose targets, new coordinator) is not ported; host.evacuate places the VMs.
        /// </summary>
        public static string CannotEnterMaintenance(Host host)
        {
            var pool = Helpers.GetPool(host.Connection);
            if (pool != null && pool.ha_enabled && host.IsCoordinator())
                return string.Format(Messages.HA_CANNOT_EVACUATE_COORDINATOR, host.Name());
            // shortcut: EvacuateHostDialog offers to eject CDs, suspend or shut these down; here the user does it first
            var stuck = host.Connection.ResolveAll(host.resident_VMs)
                .Where(v => !v.is_control_domain && v.power_state == vm_power_state.Running && !string.IsNullOrEmpty(v.IsOnSharedStorage()))
                .Select(v => v.Name()).ToList();
            if (stuck.Count > 0)
                return "These VMs cannot be migrated off the server (their disks or CDs are on its local storage). " +
                       "Shut them down, suspend them or move their storage first:\n\n" + string.Join("\n", stuck);
            if (host.IsCoordinator() && host.Connection.Cache.Hosts.Length > 1 && NewCoordinator(host) == null)
                return "No other enabled server of the pool can become the coordinator.";
            return null;
        }

        /// <summary>EvacuateHostDialog: the coordinator hands over to an enabled, live member.</summary>
        private static Host NewCoordinator(Host host) =>
            host.Connection.Cache.Hosts.Where(h => h.opaque_ref != host.opaque_ref && h.enabled && h.IsLive()).OrderBy(h => h.Name()).FirstOrDefault();

        /// <summary>EvacuateHostDialog: what entering maintenance mode does, the VMs it migrates and the coordinator note.</summary>
        public static (string Title, string Text) EvacuateConfirmation(Host host)
        {
            var text = "This operation will migrate all VMs running on this server and transition it into maintenance mode. If a VM cannot be migrated you need to suspend it or shut it down.";
            var vms = host.Connection.ResolveAll(host.resident_VMs).Where(v => !v.is_control_domain).Select(v => v.Name()).OrderBy(n => n).ToList();
            text += vms.Count == 0 ? "\n\nNo VMs are running on this server." : "\n\nVMs to migrate:\n" + string.Join("\n", vms);
            if (host.IsCoordinator())
                text += "\n\n" + string.Format("This server is the pool coordinator. Entering maintenance mode will nominate a new coordinator for the pool. {0} will temporarily lose its connection to the pool.", BrandManager.BrandConsole);
            return (string.Format(Messages.EVACUATE_HOST_DIALOG_TITLE, host.Name()), text);
        }

        public static AsyncAction EnterMaintenance(Host host)
        {
            // a coordinator with other members hands over to one of them (EvacuateHostDialog's default choice)
            var newCoordinator = host.IsCoordinator() ? NewCoordinator(host) : null;
            // EvacuateHostDialog.Scan saves where the VMs are first, so Exit Maintenance Mode can offer to restore them
            var save = new DelegatedAsyncAction(host.Connection, Messages.SAVING_VM_PROPERTIES_ACTION_TITLE, Messages.SAVING_VM_PROPERTIES_ACTION_DESC,
                Messages.COMPLETED, session => host.SaveEvacuatedVMs(session), true, "host.remove_from_other_config", "host.add_to_other_config");
            var evacuate = new EvacuateHostAction(host, newCoordinator, new Dictionary<XenRef<VM>, string[]>(), KeepNtol, KeepNtolOnEnable);
            return new MultipleAction(host.Connection, string.Format(Messages.EVACUATE_HOST_DIALOG_TITLE, host.Name()), Messages.HOST_EVACUATE,
                Messages.COMPLETED, new List<AsyncAction> { save, evacuate });
        }

        /// <summary>ExitMaintenanceMode: the evacuated VMs that are not back on the host.</summary>
        public static List<VM> VmsToRestore(Host host) =>
            host.GetHaltedEvacuatedVMs().Concat(host.GetMigratedEvacuatedVMs()).Concat(host.GetSuspendedEvacuatedVMs())
                .Where(v => v.resident_on.opaque_ref != host.opaque_ref).ToList();

        public static AsyncAction ExitMaintenance(Host host, bool restoreVms) => new EnableHostAction(host, restoreVms, KeepNtolOnEnable);
    }
}
