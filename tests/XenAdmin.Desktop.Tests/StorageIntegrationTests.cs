using System;
using XenAdmin.Actions.DR;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop.Tests
{
    /// <summary>
    /// New Storage wizard flows on the lab pool (tools/vagrant/xcp-ng-lab: lab.sh storage-up): create a shared NFS, SMB or
    /// iSCSI SR, add and delete a disk, detach, find it again (Scan / LUN probe), reattach, destroy. Needs XCP_LAB_HOSTS (pool
    /// coordinator first), XCP_ALLOW_WRITE=1 and the lab storage on the lab network's .1; never point it at real servers.
    /// </summary>
    [TestClass]
    public class StorageIntegrationTests
    {
        [ClassInitialize]
        public static void Init(TestContext ctx) => ConnectionIntegrationTests.Init(ctx);

        private static void Run(AsyncAction action)
        {
            var done = new ManualResetEventSlim();
            action.Completed += _ => done.Set();
            action.RunAsync();
            Assert.IsTrue(done.Wait(TimeSpan.FromMinutes(5)), action.Title + " timed out");
            Assert.IsTrue(action.Succeeded, action.Title + ": " + action.Exception);
        }

        private static T Until<T>(Func<T> get, Func<T, bool> ok, string what, int seconds = 120)
        {
            var deadline = DateTime.Now.AddSeconds(seconds);
            T v;
            while (!ok(v = get()) && DateTime.Now < deadline)
                Thread.Sleep(500);
            Assert.IsTrue(ok(v), what);
            return v;
        }

        /// <summary>Destroy, as the WinForms DestroySRCommand: unplug the PBDs first, then SR.destroy.</summary>
        private static void Destroy(XenConnection c, SR sr)
        {
            if (c.ResolveAll(sr.PBDs).Any(p => p.currently_attached))
                Run(new DetachSrAction(sr));
            Run(new DestroySrAction(c.Resolve(new XenRef<SR>(sr.opaque_ref)) ?? sr));
        }

        /// <summary>ISO libraries are forgotten, not destroyed (ForgetSRCommand): the share keeps its ISOs.</summary>
        private static void Forget(XenConnection c, SR sr)
        {
            if (c.ResolveAll(sr.PBDs).Any(p => p.currently_attached))
                Run(new DetachSrAction(sr));
            Run(new ForgetSrAction(c.Resolve(new XenRef<SR>(sr.opaque_ref)) ?? sr));
        }

        private static XenConnection LabPool()
        {
            var lab = Environment.GetEnvironmentVariable("XCP_LAB_HOSTS")?.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);
            if (lab == null || lab.Length == 0 || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_LAB_HOSTS and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(lab[0]);
            if (c.Cache.Hosts.Length < 2)
            {
                c.EndConnect();
                Assert.Inconclusive("the lab hosts are not a pool yet");
            }
            return c;
        }

        /// <summary>lab.sh storage-up serves NFS and SMB on the lab network's gateway (x.y.z.1).</summary>
        private static string StorageHost() =>
            Environment.GetEnvironmentVariable("XCP_LAB_HOSTS")?.Split(',')[0].Trim() is string h && h.Count(ch => ch == '.') == 3
                ? h.Substring(0, h.LastIndexOf('.')) + ".1" : "192.168.56.1";

        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void NfsSharedSrLifecycleOnLabPool() =>
            Lifecycle(SrRules.Kind.Nfs, StorageHost() + ":/export", null, null);

        /// <summary>
        /// Move Virtual Disk (MoveVirtualDiskDialog): a 64 MB unattached disk on the coordinator's local SR moves to a
        /// shared NFS SR (MoveVirtualDiskAction, the disk is not attached so it is not migrated live), then both go away.
        /// </summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void MoveDiskToNfsOnLabPool()
        {
            var c = LabPool();
            const string name = "xcpng-center-test-move-nfs", disk = "xcpng-center-test-move-disk";
            SR nfs = null;
            try
            {
                foreach (var stale in c.Cache.SRs.Where(sr => sr.name_label == name).ToList())
                    Destroy(c, stale);
                var path = StorageHost() + ":/export";
                Run(SrRules.Finish(c, SrRules.Kind.Nfs, name, SrRules.AutoDescription(SrRules.Kind.Nfs, path), SrRules.DeviceConfig(SrRules.Kind.Nfs, path, "", "", null, null), null));
                nfs = Until(() => c.Cache.SRs.FirstOrDefault(s => s.name_label == name), s => s != null && s.HasPBDs(), "NFS SR attached");

                var coordinator = Helpers.GetCoordinator(c);
                var local = c.Cache.SRs.First(s => s.IsLocalSR() && s.SupportsVdiCreate() && s.Home()?.opaque_ref == coordinator.opaque_ref);
                Run(new CreateDiskAction(new VDI { Connection = c, SR = new XenRef<SR>(local), virtual_size = 64L << 20, name_label = disk, name_description = "", type = vdi_type.user }));
                var vdi = Until(() => c.Cache.VDIs.FirstOrDefault(v => v.name_label == disk), v => v != null, "disk on the local SR");

                Assert.AreEqual(false, StorageRules.UseMigrate(vdi), "an unattached disk is moved, not migrated live");
                var targets = StorageRules.MoveTargets(vdi, false);
                Assert.IsFalse(targets.Single(t => t.Sr.opaque_ref == local.opaque_ref).Enabled, "its own SR is the current location");
                Assert.AreEqual(Messages.CURRENT_LOCATION, targets.Single(t => t.Sr.opaque_ref == local.opaque_ref).Text);
                Assert.IsTrue(targets.Single(t => t.Sr.opaque_ref == nfs.opaque_ref).Enabled, "the NFS SR can take the disk");

                Run(new MoveVirtualDiskAction(c, vdi, nfs));
                var moved = Until(() => c.Cache.VDIs.FirstOrDefault(v => v.name_label == disk), v => v != null && v.SR.opaque_ref == nfs.opaque_ref, "disk on the NFS SR");
                Assert.IsFalse(c.Cache.VDIs.Any(v => v.name_label == disk && v.SR.opaque_ref == local.opaque_ref), "the source copy is gone");
                Run(new DestroyDiskAction(moved));
            }
            finally
            {
                foreach (var v in c.Cache.VDIs.Where(v => v.name_label == disk).ToList())
                    try { VDI.destroy(c.DuplicateSession(), v.opaque_ref); } catch { }
                if (nfs != null) Destroy(c, nfs);
                c.EndConnect();
            }
        }

        /// <summary>
        /// Pool > High Availability on the lab pool: an NFS SR as the heartbeat (GetHeartbeatSRsAction), the plan's
        /// maximum failures, EnableHAAction, then Disable and the SR goes away.
        /// </summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void HaEnableDisableOnLabPool()
        {
            var c = LabPool();
            const string name = "xcpng-center-test-ha-heartbeat";
            SR sr = null;
            try
            {
                foreach (var stale in c.Cache.SRs.Where(x => x.name_label == name).ToList())
                {
                    foreach (var vdi in c.ResolveAll(stale.VDIs))
                        try { VDI.destroy(c.DuplicateSession(), vdi.opaque_ref); } catch { }
                    Destroy(c, stale);
                }
                var pool = Helpers.GetPool(c);
                Assert.IsTrue(HaRules.CanConfigure(pool, false) && !HaRules.CanDisable(pool, false), "HA off");
                var path = StorageHost() + ":/export";
                Run(SrRules.Finish(c, SrRules.Kind.Nfs, name, SrRules.AutoDescription(SrRules.Kind.Nfs, path), SrRules.DeviceConfig(SrRules.Kind.Nfs, path, "", "", null, null), null));
                sr = Until(() => c.Cache.SRs.FirstOrDefault(x => x.name_label == name), x => x != null && x.HasPBDs(), "heartbeat SR attached");

                var scan = new GetHeartbeatSRsAction(pool);
                Run(scan);
                Assert.IsTrue(scan.SRs.Any(w => w.enabled && w.sr.opaque_ref == sr.opaque_ref), "the NFS SR can hold the HA statefile");

                var vms = HaRules.ProtectableVms(c);
                var settings = vms.ToDictionary(v => v, v => VM.HaRestartPriority.BestEffort);
                var max = HaRules.MaxFailures(c.DuplicateSession(), settings);
                Assert.IsTrue(max >= 0, "the pool computes the failures it can tolerate");
                Assert.IsNull(HaRules.PlanError(0, max));
                var options = vms.ToDictionary(v => v, v => new VMStartupOptions(v.order, v.start_delay, VM.HaRestartPriority.BestEffort));
                Run(HaRules.EnableAction(pool, options, sr, 0));
                Until(() => Helpers.GetPool(c), p => p.ha_enabled, "HA on");
                Assert.IsTrue(HaRules.CanDisable(Helpers.GetPool(c), false));

                Run(HaRules.DisableAction(Helpers.GetPool(c)));
                Until(() => Helpers.GetPool(c), p => !p.ha_enabled, "HA off again");
            }
            finally
            {
                var p = Helpers.GetPool(c);
                if (p != null && p.ha_enabled)
                    try { Pool.disable_ha(c.DuplicateSession()); } catch { }
                if (sr != null)
                {
                    // disabling HA leaves its statefile and redo log on the heartbeat SR
                    foreach (var vdi in c.ResolveAll(sr.VDIs).Where(v => v.type == vdi_type.ha_statefile || v.type == vdi_type.redo_log))
                        try { VDI.destroy(c.DuplicateSession(), vdi.opaque_ref); } catch { }
                    Destroy(c, c.Resolve(new XenRef<SR>(sr.opaque_ref)) ?? sr);
                }
                c.EndConnect();
            }
        }

        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void SmbSharedSrLifecycleOnLabPool() =>
            Lifecycle(SrRules.Kind.Smb, @"\\" + StorageHost() + @"\sr", "xcp", Environment.GetEnvironmentVariable("XCP_PASS"));

        /// <summary>lab.sh storage-up: the iSCSI target sits at .20 on the KVM lab network, on the gateway for VirtualBox.</summary>
        private static string IscsiHost() =>
            StorageHost() is var gw && gw.StartsWith("10.77.58.") ? "10.77.58.20" : StorageHost();

        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void IscsiSharedSrLifecycleOnLabPool()
        {
            var c = LabPool();
            const string name = "xcpng-center-test-iscsi";
            try
            {
                foreach (var stale in c.Cache.SRs.Where(sr => sr.name_label == name).ToList())
                    Destroy(c, stale);

                // Scan Target Host -> IQN -> LUN, as the LVMoISCSI page
                var scan = SrRules.ScanIqns(c, IscsiHost(), SrRules.DefaultIscsiPort, null, null);
                Run(scan);
                var iqn = scan.IQNs.Single(i => i.TargetIQN == "iqn.2026-11.org.xcp-ng.lab:sr");
                StringAssert.StartsWith(SrRules.IqnLabel(iqn), iqn.TargetIQN);
                var luns = SrRules.ScanLuns(c, iqn, null, null);
                Run(luns);
                var lun = luns.LUNs.Single();
                StringAssert.StartsWith(SrRules.LunLabel(lun), "LUN " + lun.LunID);
                Assert.IsNull(SrRules.LunInUse(lun), "free LUN");
                var dconf = SrRules.IscsiDeviceConfig(iqn, lun, null, null);
                Assert.AreEqual(lun.ScsiID, dconf["SCSIid"]);

                // Finish: the probe finds no SR -> Format (create)
                var probe = SrRules.Probe(c, SrRules.Kind.Iscsi, dconf);
                Run(probe);
                Assert.AreEqual(0, probe.SRs.Count, "a fresh LUN holds no SR");
                Run(SrRules.Finish(c, SrRules.Kind.Iscsi, name, SrRules.IscsiDescription(iqn, lun), dconf, null));
                var sr = Until(() => c.Cache.SRs.FirstOrDefault(s => s.name_label == name), s => s != null, "SR in the cache");
                Until(() => c.ResolveAll(sr.PBDs), pbds => pbds.Count == c.Cache.Hosts.Length && pbds.All(p => p.currently_attached), "plugged on every host");
                Assert.IsTrue(sr.shared && sr.GetSRType(false) == SR.SRTypes.lvmoiscsi, "shared LVM over iSCSI");
                Assert.AreSame(sr, SrRules.LunInUse(lun), "the wizard now refuses this LUN");

                // Detach, probe the LUN again -> Reattach the SR found on it
                var uuid = sr.uuid;
                Run(new DetachSrAction(sr, true));
                Until(() => c.Cache.SRs.FirstOrDefault(s => s.uuid == uuid), s => s != null && !s.HasPBDs(), "detached");
                probe = SrRules.Probe(c, SrRules.Kind.Iscsi, dconf);
                Run(probe);
                Assert.AreEqual(uuid, probe.SRs.Single().UUID, "the LUN probe finds the detached SR");
                Run(SrRules.Finish(c, SrRules.Kind.Iscsi, name, "", dconf, uuid));
                sr = Until(() => c.Cache.SRs.FirstOrDefault(s => s.uuid == uuid), s => s != null && s.HasPBDs()
                    && c.ResolveAll(s.PBDs).Count == c.Cache.Hosts.Length && c.ResolveAll(s.PBDs).All(p => p.currently_attached), "reattached on every host");

                Destroy(c, sr);
                Until(() => c.Cache.SRs.Any(s => s.uuid == uuid), any => !any, "destroyed");
            }
            finally
            {
                c.EndConnect();
            }
        }

        /// <summary>
        /// Disaster Recovery on the lab: the pool stores its metadata on an iSCSI SR holding a throwaway VM's disk (DR
        /// Configure); the SR leaves the pool, the member leaves too and attaches the SR as a standalone server, where the
        /// DR wizard's rules find the pool's metadata, run the prechecks, recover the VM (a test failover) and remove it.
        /// Everything is put back: the member rejoins and the SR and VM are destroyed. Needs XCP_LAB_HOSTS=coordinator,member.
        /// </summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void DisasterRecoveryTestFailoverOnLab()
        {
            var lab = Environment.GetEnvironmentVariable("XCP_LAB_HOSTS")?.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);
            if (lab == null || lab.Length != 2 || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_LAB_HOSTS=coordinator,member and XCP_ALLOW_WRITE=1 required");
            var c = LabPool();
            XenConnection member = null;
            const string name = "xcpng-center-test-dr";
            const string vmName = "xcpng-center-linux-test-dr";
            string uuid = null;
            Dictionary<string, string> dconf = null;
            try
            {
                foreach (var stale in c.Cache.SRs.Where(sr => sr.name_label == name).ToList())
                    DestroyDrSr(c, stale);
                var pool = Helpers.GetPool(c);
                var scan = SrRules.ScanIqns(c, IscsiHost(), SrRules.DefaultIscsiPort, null, null);
                Run(scan);
                var iqn = scan.IQNs.Single(i => i.TargetIQN == "iqn.2026-11.org.xcp-ng.lab:sr");
                var luns = SrRules.ScanLuns(c, iqn, null, null);
                Run(luns);
                dconf = SrRules.IscsiDeviceConfig(iqn, luns.LUNs.Single(), null, null);
                var probe = SrRules.Probe(c, SrRules.Kind.Iscsi, dconf);
                Run(probe);
                Run(SrRules.Finish(c, SrRules.Kind.Iscsi, name, "DR test", dconf, probe.SRs.Count > 0 ? probe.SRs[0].UUID : null));
                var sr = Until(() => c.Cache.SRs.FirstOrDefault(x => x.name_label == name && x.HasPBDs()), x => x != null, "SR in the cache");
                uuid = sr.uuid;

                // a throwaway VM with a 64 MB disk on that SR
                var session = c.DuplicateSession();
                var template = c.Cache.VMs.First(v => v.is_a_template && v.name_label == "Other install media");
                var vmRef = VM.clone(session, template.opaque_ref, vmName).opaque_ref;
                VM.set_is_a_template(session, vmRef, false);
                var vdiRef = VDI.create(session, new VDI { name_label = vmName + "-disk", SR = new XenRef<SR>(sr.opaque_ref), virtual_size = 64L << 20, type = vdi_type.user, sharable = false, read_only = false, other_config = new Dictionary<string, string>(), sm_config = new Dictionary<string, string>(), xenstore_data = new Dictionary<string, string>(), tags = new string[0] });
                VBD.create(session, new VBD { VM = new XenRef<VM>(vmRef), VDI = vdiRef, userdevice = "0", bootable = true, mode = vbd_mode.RW, type = vbd_type.Disk, empty = false, other_config = new Dictionary<string, string>(), qos_algorithm_params = new Dictionary<string, string>(), qos_algorithm_type = "", device = "", unpluggable = true });
                var vmUuid = VM.get_uuid(session, vmRef);

                // DR Configure: the pool's metadata on that SR
                Assert.IsTrue(DrRules.ReplicationSrs(c).Any(x => x.uuid == uuid), "iSCSI supports database replication");
                if (DrRules.IsDrEnabled(sr, pool))
                {
                    // a reattached SR from an earlier run still holds the pool's metadata: DR Configure turns it off
                    Assert.IsNotNull(DrRules.ConfigureWarning(true, 0).Text, "the disable warning");
                    Run(DrRules.ConfigureAction(pool, new SR[0], new[] { sr }));
                    Until(() => c.Cache.SRs.First(x => x.uuid == uuid), x => !DrRules.IsDrEnabled(x, pool), "DR disabled");
                }
                Assert.IsNotNull(DrRules.ConfigureWarning(false, 1).Text, "the enable warning");
                Run(DrRules.ConfigureAction(pool, new[] { sr }, new SR[0]));
                Until(() => c.Cache.SRs.First(x => x.uuid == uuid), x => DrRules.IsDrEnabled(x, pool), "DR enabled");
                Thread.Sleep(10000); // XAPI writes the database to the metadata VDI

                // the primary site goes away: the SR leaves the pool, then the member leaves the pool
                Run(new DetachSrAction(c.Cache.SRs.First(x => x.uuid == uuid), true));
                var coordinator = Helpers.GetCoordinator(c);
                var joined = c.Cache.Hosts.Single(h => h.opaque_ref != coordinator.opaque_ref);
                RunLong(new EjectHostFromPoolAction(pool, joined), 10);
                var deadline = DateTime.Now.AddMinutes(10);
                while (member == null && DateTime.Now < deadline)
                {
                    try { member = ConnectionIntegrationTests.Connect(lab[1]); }
                    catch (AssertFailedException) { Thread.Sleep(15000); }
                }
                Assert.IsNotNull(member, "the ejected server came back");

                // the DR site attaches the mirrored SR
                Run(SrRules.Finish(member, SrRules.Kind.Iscsi, name, "", dconf, uuid));
                var drSr = Until(() => member.Cache.SRs.FirstOrDefault(x => x.uuid == uuid), x => x != null && x.HasPBDs(), "SR attached on the DR site");

                // DR wizard: storage, metadata, prechecks, recovery (test failover), clean-up
                var choice = DrRules.StorageChoices(member).Single(x => x.Sr.uuid == uuid);
                Assert.IsTrue(choice.HasMetadata, "metadata detected");
                var metadata = DrRules.LoadMetadata(member, new[] { drSr });
                var md = metadata.Single(m => m.Pool.uuid == pool.uuid);
                Assert.IsFalse(DrRules.IsCurrentPool(md, member));
                var vm = DrRules.LooseVms(md).Single(v => v.uuid == vmUuid);
                var chosen = new List<(PoolMetadata, List<VM_appliance>, List<VM>)> { (md, new List<VM_appliance>(), new List<VM> { vm }) };
                var checks = DrRules.Prechecks(member, DrType.Test, chosen);
                Assert.IsFalse(checks.Any(x => x.Error), string.Join("; ", checks.Select(x => x.Text)));
                var recover = new DrRules.RecoverAction(member, chosen, DrRules.StartMode.None);
                RunLong(recover, 5);
                CollectionAssert.Contains(recover.RecoveredVmUuids, vmUuid);
                RunLong(DrRules.TestCleanupAction(member, recover), 5); // it waits for the recovered VM itself
                Until(() => member.Cache.VMs.Any(v => v.uuid == vmUuid), any => !any, "test failover cleaned up");
            }
            finally
            {
                // put the lab back: the member forgets the SR and rejoins; the pool destroys the VM and the SR
                if (member != null)
                {
                    try
                    {
                        if (member.Cache.SRs.FirstOrDefault(x => x.uuid == uuid) is SR drSr)
                            Forget(member, drSr);
                        RunLong(PoolRules.AddToPool(Helpers.GetCoordinator(member), Helpers.GetPool(c)), 10);
                        Until(() => c.Cache.Hosts.Length == 2 && c.Cache.Hosts.All(h => h.IsLive()), x => x, "pool of two again", 300);
                    }
                    catch (Exception e) { System.Console.WriteLine("DR test clean-up (member): " + e.Message); }
                }
                try
                {
                    var session = c.DuplicateSession();
                    foreach (var leftover in c.Cache.VMs.Where(v => v.name_label == vmName).ToList())
                        VM.destroy(session, leftover.opaque_ref);
                    if (uuid != null && c.Cache.SRs.FirstOrDefault(x => x.uuid == uuid) is SR sr)
                    {
                        if (!sr.HasPBDs() && dconf != null)
                            Run(SrRules.Finish(c, SrRules.Kind.Iscsi, name, "", dconf, uuid));
                        DestroyDrSr(c, Until(() => c.Cache.SRs.First(x => x.uuid == uuid), x => x.HasPBDs(), "SR reattached"));
                    }
                }
                catch (Exception e) { System.Console.WriteLine("DR test clean-up (pool): " + e.Message); }
                lock (ConnectionsManager.ConnectionsLock)
                {
                    ConnectionsManager.XenConnections.Remove(c);
                    if (member != null) ConnectionsManager.XenConnections.Remove(member);
                }
                c.EndConnect();
                member?.EndConnect();
            }
        }

        /// <summary>The DR test's SR: stop the replication, remove its disks (the VM's and the metadata), then destroy it.</summary>
        private static void DestroyDrSr(XenConnection c, SR sr)
        {
            var session = c.DuplicateSession();
            foreach (var pbd in c.ResolveAll(sr.PBDs).Where(p => !p.currently_attached))
                try { PBD.plug(session, pbd.opaque_ref); } catch { }
            try { SR.disable_database_replication(session, sr.opaque_ref); } catch { }
            foreach (var vdi in SR.get_VDIs(session, sr.opaque_ref))
                try { VDI.destroy(session, vdi.opaque_ref); } catch { }
            // straight to the server: the cache may not have seen the PBDs plugged above yet
            foreach (var pbd in SR.get_PBDs(session, sr.opaque_ref))
            {
                if (PBD.get_currently_attached(session, pbd.opaque_ref))
                    PBD.unplug(session, pbd.opaque_ref);
                PBD.destroy(session, pbd.opaque_ref);
            }
            SR.forget(session, sr.opaque_ref); // no PBD left to destroy it with; its disks are gone already
        }

        private static void RunLong(AsyncAction action, int minutes)
        {
            var done = new ManualResetEventSlim();
            action.Completed += _ => done.Set();
            action.RunAsync();
            Assert.IsTrue(done.Wait(TimeSpan.FromMinutes(minutes)), action.Title + " timed out");
            Assert.IsTrue(action.Succeeded, action.Exception?.ToString());
        }

        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void NfsIsoLibraryOnLabPool()
        {
            var c = LabPool();
            const string name = "xcpng-center-test-iso";
            var path = StorageHost() + ":/export/iso";
            try
            {
                foreach (var stale in c.Cache.SRs.Where(sr => sr.name_label == name).ToList())
                    Forget(c, stale);
                Until(() => c.Cache.SRs.Any(s => s.name_label == name), any => !any, "leftover library forgotten");
                Assert.IsTrue(SrRules.ValidLocation(SrRules.Kind.NfsIso, path));
                Assert.IsNull(SrRules.IsoAlreadyAttached(c, SrRules.Kind.NfsIso, path));
                Run(SrRules.Finish(c, SrRules.Kind.NfsIso, name, SrRules.AutoDescription(SrRules.Kind.NfsIso, path),
                    SrRules.IsoDeviceConfig(SrRules.Kind.NfsIso, path, false, null, null), null));
                var sr = Until(() => c.Cache.SRs.FirstOrDefault(s => s.name_label == name), s => s != null, "ISO SR in the cache");
                Assert.AreEqual(SR.Content_Type_ISO, sr.content_type);
                Assert.IsTrue(sr.shared, "the library is shared by the pool");
                Until(() => c.ResolveAll(sr.VDIs), vdis => vdis.Any(v => v.name_label.EndsWith(".iso")), "the share's ISOs are listed");
                Assert.IsNotNull(SrRules.IsoAlreadyAttached(c, SrRules.Kind.NfsIso, path), "the wizard refuses the same share twice");
                Forget(c, sr);
            }
            finally
            {
                c.EndConnect();
            }
        }

        private static void Lifecycle(SrRules.Kind kind, string path, string user, string password)
        {
            var c = LabPool();
            var name = "xcpng-center-test-" + kind.ToString().ToLowerInvariant();
            try
            {
                Assert.IsTrue(SrRules.ValidLocation(kind, path), path);
                foreach (var stale in c.Cache.SRs.Where(sr => sr.name_label == name).ToList())
                    Destroy(c, stale);

                // Location page: Scan, then Create a new SR
                var probe = SrRules.Probe(c, kind, SrRules.DeviceConfig(kind, path, "", "", user, password, forProbe: true));
                Run(probe);
                var before = probe.SRs?.Select(i => i.UUID).ToList() ?? new();
                Run(SrRules.Finish(c, kind, name, SrRules.AutoDescription(kind, path), SrRules.DeviceConfig(kind, path, "", "", user, password), null));
                var sr = Until(() => c.Cache.SRs.FirstOrDefault(s => s.name_label == name), s => s != null, "SR in the cache");
                Until(() => c.ResolveAll(sr.PBDs), pbds => pbds.Count == c.Cache.Hosts.Length && pbds.All(p => p.currently_attached), "plugged on every host");
                Assert.IsTrue(sr.shared, "shared SR");
                if (password != null)
                    Assert.IsTrue(c.ResolveAll(sr.PBDs).All(p => !p.device_config.ContainsKey("password") && p.device_config.ContainsKey("password_secret")),
                        "the SMB password is kept as a secret, never in the PBD");
                Assert.IsTrue(StorageRules.CanHoldNewDisk(sr, null), "offered for new disks");

                // a disk on it, from the coordinator and visible to the member
                var session = c.DuplicateSession(); // shares the connection's session: never log it out
                var vdi = VDI.create(session, new VDI
                {
                    name_label = name + "-disk", SR = new XenRef<SR>(sr.opaque_ref), virtual_size = 64L << 20, type = vdi_type.user,
                    sharable = false, read_only = false, other_config = new(), xenstore_data = new(), sm_config = new(), tags = new string[0]
                });
                Until(() => c.Resolve(vdi), v => v != null, "VDI in the cache");
                VDI.destroy(session, vdi.opaque_ref);

                // Detach, find it again with Scan and reattach it (Reattach an existing SR)
                var uuid = sr.uuid;
                Run(new DetachSrAction(sr, true));
                Until(() => c.Cache.SRs.FirstOrDefault(s => s.uuid == uuid), s => s != null && !s.HasPBDs(), "detached");
                probe = SrRules.Probe(c, kind, SrRules.DeviceConfig(kind, path, "", "", user, password, forProbe: true));
                Run(probe);
                Assert.IsTrue(probe.SRs.Any(i => i.UUID == uuid), "Scan finds the detached SR (found before: " + before.Count + ")");
                Assert.IsNull(SrRules.InUse(uuid), "not in use while detached");
                Run(SrRules.Finish(c, kind, name, "", SrRules.DeviceConfig(kind, path, "", "", user, password), uuid));
                sr = Until(() => c.Cache.SRs.FirstOrDefault(s => s.uuid == uuid), s => s != null && s.HasPBDs()
                    && c.ResolveAll(s.PBDs).Count == c.Cache.Hosts.Length && c.ResolveAll(s.PBDs).All(p => p.currently_attached), "reattached on every host");

                Destroy(c, sr);
                Until(() => c.Cache.SRs.Any(s => s.uuid == uuid), any => !any, "destroyed");
            }
            finally
            {
                if (password != null)
                    RemoveOrphanSecrets(c, password);
                c.EndConnect();
            }
        }

        /// <summary>
        /// XAPI turns a PBD's "password" into a secret of its own (the reattach creates one per host) and keeps it when
        /// the PBD is destroyed: remove the test's leftovers so the lab does not collect them.
        /// </summary>
        private static void RemoveOrphanSecrets(XenConnection c, string password)
        {
            var session = c.DuplicateSession(); // shares the connection's session: never log it out
            var used = PBD.get_all_records(session).Values.SelectMany(p => p.device_config.Where(kv => kv.Key.EndsWith("secret")).Select(kv => kv.Value)).ToHashSet();
            foreach (var (reference, secret) in Secret.get_all_records(session).Select(kv => (kv.Key, kv.Value)))
                if (!used.Contains(secret.uuid) && secret.value == password)
                    Secret.destroy(session, reference.opaque_ref);
        }
    }
}
