using System;
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
                c.EndConnect();
            }
        }
    }
}
