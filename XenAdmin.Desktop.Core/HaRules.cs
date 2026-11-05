using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// High availability: the rules of HAConfigureCommand / HADisableCommand (XenAdmin/Commands/HACommand.cs), the HA
    /// wizard (ChooseSR, AssignPriorities with HaNtolControl) and EditVmHaPrioritiesDialog.
    /// </summary>
    public static class HaRules
    {
        /// <summary>HACommand.CanRunHACommand: a visible pool with its coordinator, not busy (busy: an HA action runs).</summary>
        private static bool CanRun(Pool pool, bool busy) =>
            pool != null && pool.IsVisible() && !pool.Locked && pool.Connection != null && pool.Connection.IsConnected
            && Helpers.GetCoordinator(pool.Connection) != null && !busy;

        public static bool CanConfigure(Pool pool, bool busy) => CanRun(pool, busy);

        public static bool CanDisable(Pool pool, bool busy) => CanRun(pool, busy) && pool.ha_enabled;

        /// <summary>HAConfigureCommand with HA on: the statefile must be resolvable before the priorities can be edited.</summary>
        public static string ConfigureBlocker(Pool pool) =>
            pool.ha_enabled && pool.ha_statefiles.All(sf => pool.Connection.Resolve(new XenRef<VDI>(sf)) == null)
                ? string.Format(Messages.HA_CONFIGURE_NO_STATEFILE, Helpers.GetName(pool).Ellipsise(30))
                : null;

        /// <summary>AssignPriorities: the VMs HA can protect.</summary>
        public static List<VM> ProtectableVms(IXenConnection c) =>
            c.Cache.VMs.Where(v => v.HaCanProtect(false)).OrderBy(v => v.Name()).ToList();

        /// <summary>The restart priorities this pool offers, with their names.</summary>
        public static List<(VM.HaRestartPriority Priority, string Text)> Priorities(IXenConnection c) =>
            VM.GetAvailableRestartPriorities(c).Select(p => (p, Helpers.RestartPriorityI18n(p))).ToList();

        /// <summary>AssignPriorities.IsHaActivatedFirstTime: no VM has a restart priority yet.</summary>
        public static bool FirstTime(IEnumerable<VM> vms) => vms.All(v => string.IsNullOrEmpty(v.ha_restart_priority));

        /// <summary>The priority a VM starts with: its own, or (first time) the highest when agile, else best effort.</summary>
        public static VM.HaRestartPriority InitialPriority(VM vm, bool firstTime, bool agile) =>
            !firstTime ? vm.HARestartPriority()
            : agile ? VM.HaHighestProtectionAvailable(vm.Connection) : VM.HaRestartPriority.BestEffort;

        /// <summary>AssignPriorities.bgWorker_DoWork: null when the VM is agile, else the reason (VM.assert_agile).</summary>
        public static string NotAgileReason(Session session, VM vm)
        {
            try
            {
                VM.assert_agile(session, vm.opaque_ref);
                return null;
            }
            catch (Failure f)
            {
                return f.Message;
            }
        }

        /// <summary>A VM that is not agile cannot be protected (restarted): only best effort or do not restart.</summary>
        public static bool Allowed(VM.HaRestartPriority p, bool agile) =>
            agile || p == VM.HaRestartPriority.BestEffort || p == VM.HaRestartPriority.DoNotRestart;

        /// <summary>HaNtolControl: the most server failures the plan tolerates (Pool.ha_compute_hypothetical_max...).</summary>
        public static long MaxFailures(Session session, Dictionary<VM, VM.HaRestartPriority> settings) =>
            Pool.GetMaximumTolerableHostFailures(session, Helpers.GetVmHaRestartPrioritiesForApi(settings));

        /// <summary>AssignPriorities.updateButtons: a plan can be applied when the failures to tolerate fit (not overcommitted).</summary>
        public static string PlanError(long ntol, long max)
        {
            if (max < 0) return "The number of server failures that can be tolerated is still being computed.";
            if (ntol < 0) return "Enter the number of server failures to tolerate.";
            return ntol > max ? Messages.HA_OVERCOMMIT_MESSAGE : null;
        }

        /// <summary>HAWizard.FinishWizard: EnableHAAction with the heartbeat SR, the priorities and the failures to tolerate.</summary>
        public static AsyncAction EnableAction(Pool pool, Dictionary<VM, VMStartupOptions> options, SR heartbeat, long ntol) =>
            new EnableHAAction(pool, options, new List<SR> { heartbeat }, ntol);

        /// <summary>EditVmHaPrioritiesDialog: SetHaPrioritiesAction.</summary>
        public static AsyncAction SetPrioritiesAction(IXenConnection c, Dictionary<VM, VMStartupOptions> options, long ntol) =>
            new SetHaPrioritiesAction(c, options, ntol, false);

        /// <summary>HADisableCommand's confirmation.</summary>
        public static (string Title, string Text) DisableConfirmation(Pool pool) =>
            (Messages.DISABLE_HA, string.Format(Messages.HA_DISABLE_QUERY, Helpers.GetName(pool).Ellipsise(30)));

        public static AsyncAction DisableAction(Pool pool) => new DisableHAAction(pool);
    }
}
