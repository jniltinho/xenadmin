using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Actions.DR;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>The three runs of the DR wizard (DRFailoverWizardFirstPage).</summary>
    public enum DrType { Failover, Failback, Test }

    /// <summary>
    /// Pool > Disaster Recovery: Configure (DRConfigureCommand / DRConfigureDialog: the SRs holding the pool's VM
    /// metadata) and the Disaster Recovery wizard (DRFailoverWizard: storage, pools and VMs from the metadata, prechecks,
    /// recovery, and the clean-up of a test failover).
    /// shortcut: the wizard reads metadata from SRs already attached to this pool; it does not probe iSCSI/HBA devices
    /// for mirrored SRs (DRFailoverWizardStoragePage "Find Storage Repositories"); attach them with New SR first.
    /// </summary>
    public static class DrRules
    {
        /// <summary>DRConfigureDialog.MAX_SR_SELECTED.</summary>
        public const int MaxSrs = 8;

        /// <summary>DRConfigureCommand / DisasterRecoveryCommand: a connected pool or standalone server.</summary>
        public static bool CanRun(IXenConnection c) => c != null && c.IsConnected && Helpers.GetPoolOfOne(c) != null;

        /// <summary>DRConfigureDialog._worker_DoWork: the SRs that support database replication (one server call each).</summary>
        public static List<SR> ReplicationSrs(IXenConnection c) =>
            c.Cache.SRs.Where(sr => SR.SupportsDatabaseReplication(c, sr)).OrderBy(sr => sr.Name()).ToList();

        /// <summary>SrRow.IsDrEnabled: the SR holds this pool's metadata VDI.</summary>
        public static bool IsDrEnabled(SR sr, Pool pool) =>
            sr.Connection.ResolveAll(sr.VDIs).Any(v => v.type == vdi_type.metadata && v.metadata_of_pool?.opaque_ref == pool.opaque_ref);

        /// <summary>DRConfigureDialog.ToggleWarningsVisibleState: the warning for this choice, and whether OK is allowed.</summary>
        public static (string Text, bool CanSave) ConfigureWarning(bool originallyEnabled, int checkedSrs)
        {
            if (checkedSrs > MaxSrs)
                return (string.Format(Messages.DR_CONFIGURE_TOO_MANY_SRS, MaxSrs), false);
            if (originallyEnabled && checkedSrs == 0)
                return ("Because no VM configuration data will be stored anywhere in the pool, it will not be possible to recover VMs or vApp from this pool in the event of a disaster.", true);
            if (!originallyEnabled && checkedSrs > 0)
                return ("DR configuration is not complete until storage mirroring has been set up and the selected SRs are being replicated.", true);
            return (null, true);
        }

        /// <summary>DRConfigureCommand.RunCore: disable on the cleared SRs, enable on the ticked ones (null when nothing changed).</summary>
        public static AsyncAction ConfigureAction(Pool pool, IEnumerable<SR> enable, IEnumerable<SR> disable)
        {
            var actions = new List<AsyncAction>();
            foreach (var sr in disable)
            {
                var r = sr.opaque_ref;
                actions.Add(new DelegatedAsyncAction(pool.Connection, string.Format(Messages.ACTION_DR_DISABLING_ON, sr.Name()),
                    Messages.ACTION_DR_DISABLING, Messages.ACTION_DR_DISABLED, s => SR.disable_database_replication(s, r)) { Pool = pool });
            }
            foreach (var sr in enable)
            {
                var r = sr.opaque_ref;
                actions.Add(new DelegatedAsyncAction(pool.Connection, string.Format(Messages.ACTION_DR_ENABLING_ON, sr.Name()),
                    Messages.ACTION_DR_ENABLING, Messages.ACTION_DR_ENABLED, s => SR.enable_database_replication(s, r)) { Pool = pool });
            }
            return actions.Count == 0 ? null : new MultipleAction(pool.Connection, Messages.ACTION_DR_CONFIGURING, "", "", actions, true);
        }

        // ---------- the wizard ----------

        /// <summary>DRFailoverWizardStoragePage: an SR that may hold metadata, and whether a metadata VDI is on it.</summary>
        public static List<(SR Sr, bool HasMetadata)> StorageChoices(IXenConnection c) =>
            ReplicationSrs(c).Select(sr => (sr, c.ResolveAll(sr.VDIs).Any(v => v.type == vdi_type.metadata))).ToList();

        /// <summary>
        /// DRFailoverWizardStoragePage.LoadMetadata: the pools whose metadata is on these SRs (metadata VDIs not in use),
        /// with their VMs and vApps. Runs server calls: call it off the UI thread.
        /// </summary>
        public static List<PoolMetadata> LoadMetadata(IXenConnection c, IEnumerable<SR> srs)
        {
            var uuids = srs.Select(s => s.uuid).ToList();
            var find = new GetMetadataVDIsAction(c, uuids);
            find.RunSync(c.DuplicateSession());
            var result = new List<PoolMetadata>();
            foreach (var vdi in find.VDIs ?? new List<VDI>())
            {
                if (c.ResolveAll(vdi.VBDs).Any(vbd => vbd.currently_attached))
                    continue; // in use: DRFailoverWizardStoragePage does not load it
                var load = new VdiLoadMetadataAction(c, vdi);
                try
                {
                    load.RunSync(c.DuplicateSession());
                    if (load.PoolMetadata?.Pool != null && result.All(m => m.Pool.uuid != load.PoolMetadata.Pool.uuid))
                        result.Add(load.PoolMetadata);
                }
                finally
                {
                    try { load.MetadataSession?.logout(); } catch { }
                }
            }
            return result.OrderBy(m => m.Pool.Name()).ToList();
        }

        /// <summary>DRFailoverWizardAppliancesPage: the current pool's own metadata cannot be recovered onto it.</summary>
        public static bool IsCurrentPool(PoolMetadata metadata, IXenConnection c) => metadata.Pool.uuid == Helpers.GetPoolOfOne(c)?.uuid;

        /// <summary>The VMs listed for a pool: those outside a vApp (a vApp brings its own VMs).</summary>
        public static List<VM> LooseVms(PoolMetadata metadata) =>
            metadata.Vms.Values.Where(v => !v.IsAssignedToVapp()).OrderBy(v => v.Name()).ToList();

        /// <summary>
        /// DRFailoverWizardPrecheckPage.GenerateChecks: for a test failover HA must be off on this pool; for a failover or
        /// failback, a VM still running in its source pool (when that pool is connected here) blocks the recovery; a VM
        /// (or a vApp's VM) that already exists here must be older and halted; every chosen object must pass VM /
        /// VM_appliance.assert_can_be_recovered. Errors block the recovery; warnings (an existing VM to be replaced) do
        /// not. Runs server calls.
        /// </summary>
        public static List<(string Text, bool Error)> Prechecks(IXenConnection c, DrType type, IEnumerable<(PoolMetadata Metadata, List<VM_appliance> Vapps, List<VM> Vms)> chosen)
        {
            var result = new List<(string, bool)>();
            var pool = Helpers.GetPoolOfOne(c);
            if (type == DrType.Test && pool.ha_enabled)
                result.Add((string.Format("HA is enabled on '{0}'. Disable HA before you test a failover on this pool.", pool.Name()), true));
            foreach (var (metadata, vapps, vms) in chosen)
            {
                var allVms = vms.Concat(vapps.SelectMany(a => a.VMs.Select(r => metadata.Vms.TryGetValue(r, out var vm) ? vm : null)).Where(v => v != null)).ToList();
                // RunningVmCheck / RunningVmApplianceCheck: the source pool, if this app is connected to it
                var source = ConnectionsManager.XenConnectionsCopy.FirstOrDefault(x => x.IsConnected && Helpers.GetPoolOfOne(x)?.uuid == metadata.Pool.uuid);
                if (type != DrType.Test && source != null)
                    foreach (var running in allVms.Select(v => source.Cache.VMs.FirstOrDefault(x => x.uuid == v.uuid))
                                 .Where(x => x != null && x.power_state != vm_power_state.Halted))
                        result.Add((string.Format("The VM '{0}' is still running in its pool '{1}': shut it down there first.", running.Name(), metadata.Pool.Name()), true));
                foreach (var vm in allVms)
                {
                    var existing = c.Cache.VMs.FirstOrDefault(v => v.uuid == vm.uuid);
                    if (existing == null) continue;
                    if (existing.version >= vm.version)
                        result.Add((string.Format("The VM '{0}' already exists in this pool and is not older than the one in the metadata: delete it first.", existing.Name()), true));
                    else if (existing.power_state != vm_power_state.Halted)
                        result.Add((string.Format("The VM '{0}' already exists in this pool and is running: shut it down first.", existing.Name()), true));
                    else
                        result.Add((string.Format("The VM '{0}' already exists in this pool and will be replaced.", existing.Name()), false));
                }
                var open = new VdiOpenDatabaseAction(c, metadata.Vdi);
                try
                {
                    var session = c.DuplicateSession();
                    open.RunSync(session);
                    if (open.MetadataSession == null)
                    {
                        result.Add((string.Format("The metadata of '{0}' could not be opened.", metadata.Pool.Name()), true));
                        continue;
                    }
                    foreach (var app in vapps)
                        try { VM_appliance.assert_can_be_recovered(open.MetadataSession, app.opaque_ref, session.opaque_ref); }
                        catch (Failure f) { result.Add((string.Format("vApp '{0}': {1}", app.Name(), f.Message), true)); }
                    foreach (var vm in vms)
                        try { VM.assert_can_be_recovered(open.MetadataSession, vm.opaque_ref, session.opaque_ref); }
                        catch (Failure f) { result.Add((string.Format("VM '{0}': {1}", vm.Name(), f.Message), true)); }
                }
                finally
                {
                    try { open.MetadataSession?.logout(); } catch { }
                }
            }
            return result;
        }

        /// <summary>What to do with the recovered VMs and vApps (DRFailoverWizardAppliancesPage power state).</summary>
        public enum StartMode { None, Start, StartPaused }

        /// <summary>
        /// DRFailoverWizardRecoverPage: per pool, open its metadata database (VdiOpenDatabaseAction) and recover each
        /// chosen vApp and VM (DrRecoverAction), then start them if asked (StartVMsAndAppliancesAction). The UUIDs of what
        /// was recovered are kept for the clean-up of a test failover.
        /// </summary>
        public sealed class RecoverAction : AsyncAction
        {
            private readonly List<(PoolMetadata Metadata, List<VM_appliance> Vapps, List<VM> Vms)> _chosen;
            private readonly StartMode _start;

            public readonly List<string> RecoveredVmUuids = new List<string>();
            public readonly List<string> RecoveredVappUuids = new List<string>();
            public readonly List<string> Failures = new List<string>();

            public RecoverAction(IXenConnection c, List<(PoolMetadata, List<VM_appliance>, List<VM>)> chosen, StartMode start)
                : base(c, "Disaster recovery")
            {
                _chosen = chosen;
                _start = start;
                Pool = Helpers.GetPoolOfOne(c);
            }

            protected override void Run()
            {
                var total = Math.Max(1, _chosen.Sum(x => x.Vapps.Count + x.Vms.Count));
                var done = 0;
                foreach (var (metadata, vapps, vms) in _chosen)
                {
                    Description = string.Format(Messages.DR_WIZARD_RECOVERPAGE_RECOVERING_FROM, metadata.Vdi.Name());
                    var open = new VdiOpenDatabaseAction(Connection, metadata.Vdi);
                    open.RunSync(Session);
                    if (open.MetadataSession == null)
                    {
                        Failures.Add(string.Format("The metadata of '{0}' could not be opened.", metadata.Pool.Name()));
                        continue;
                    }
                    try
                    {
                        foreach (var obj in vapps.Cast<IXenObject>().Concat(vms))
                        {
                            var recover = new DrRecoverAction(Connection, obj) { MetadataSession = open.MetadataSession };
                            try
                            {
                                recover.RunSync(Session);
                                if (obj is VM vm) RecoveredVmUuids.Add(vm.uuid);
                                else if (obj is VM_appliance app) RecoveredVappUuids.Add(app.uuid);
                            }
                            catch (Exception e)
                            {
                                Failures.Add(string.Format("{0}: {1}", obj.Name(), e.Message));
                            }
                            PercentComplete = ++done * 90 / total;
                        }
                    }
                    finally
                    {
                        try { open.MetadataSession.logout(); } catch { }
                    }
                }
                if (_start != StartMode.None && (RecoveredVmUuids.Count > 0 || RecoveredVappUuids.Count > 0))
                {
                    Description = Messages.ACTION_START_VMS_AND_APPLIANCES_TITLE;
                    var vmsToStart = RecoveredVms(Connection, RecoveredVmUuids).Where(v => v.power_state == vm_power_state.Halted).ToList();
                    var appsToStart = RecoveredVapps(Connection, RecoveredVappUuids);
                    // never remove a VM's HA protection behind the user's back (as VmOperations' AllowHaUnprotect = false)
                    new StartVMsAndAppliancesAction(Connection, vmsToStart, appsToStart, (_, _) => throw new CancelledException(), (_, _) => { }, _start == StartMode.StartPaused).RunSync(Session);
                }
                if (Failures.Count > 0)
                    throw new Exception(string.Join("\n", Failures));
                Description = Messages.COMPLETED;
            }
        }

        private static List<VM> RecoveredVms(IXenConnection c, ICollection<string> uuids)
        {
            // the cache catches up with the recovered objects shortly after the calls return
            for (var i = 0; i < 40 && c.Cache.VMs.Count(v => uuids.Contains(v.uuid)) < uuids.Count; i++)
                System.Threading.Thread.Sleep(250);
            return c.Cache.VMs.Where(v => uuids.Contains(v.uuid)).ToList();
        }

        private static List<VM_appliance> RecoveredVapps(IXenConnection c, ICollection<string> uuids)
        {
            for (var i = 0; i < 40 && c.Cache.VM_appliances.Count(a => uuids.Contains(a.uuid)) < uuids.Count; i++)
                System.Threading.Thread.Sleep(250);
            return c.Cache.VM_appliances.Where(a => uuids.Contains(a.uuid)).ToList();
        }

        /// <summary>DRFailoverWizard (test failover) on Finish: shut down and destroy what the test recovered (null when nothing).</summary>
        public static AsyncAction TestCleanupAction(IXenConnection c, RecoverAction recovered) =>
            recovered.RecoveredVmUuids.Count + recovered.RecoveredVappUuids.Count == 0 ? null : new TestCleanup(c, recovered);

        /// <summary>Waits until every recovered object is in the cache, then ShutdownAndDestroy them (vApps first).</summary>
        private sealed class TestCleanup : AsyncAction
        {
            private readonly RecoverAction _recovered;

            public TestCleanup(IXenConnection c, RecoverAction recovered) : base(c, "Clean up the test failover")
            {
                _recovered = recovered;
                Pool = Helpers.GetPoolOfOne(c);
            }

            protected override void Run()
            {
                // the VMs recovered on their own (a vApp's VMs go with the vApp)
                var apps = RecoveredVapps(Connection, _recovered.RecoveredVappUuids);
                // the vApps' member VMs reach the cache too: ShutdownAndDestroyVmAppliancesAction skips unresolved ones
                for (var i = 0; i < 40 && apps.Any(a => Connection.ResolveAll(a.VMs).Count < a.VMs.Count); i++)
                    System.Threading.Thread.Sleep(250);
                var vms = RecoveredVms(Connection, _recovered.RecoveredVmUuids);
                if (apps.Count > 0) new ShutdownAndDestroyVmAppliancesAction(Connection, apps).RunSync(Session);
                if (vms.Count > 0) new ShutdownAndDestroyVMsAction(Connection, vms).RunSync(Session);
                if (apps.Count < _recovered.RecoveredVappUuids.Count || vms.Count < _recovered.RecoveredVmUuids.Count)
                    throw new Exception("Some recovered vApps or VMs could not be found to remove.");
            }
        }
    }
}
