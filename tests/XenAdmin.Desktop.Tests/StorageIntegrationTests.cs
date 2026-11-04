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
    /// New Storage wizard flows on the lab pool (tools/vagrant/xcp-ng-lab: lab.sh storage-up): create a shared NFS or SMB
    /// SR, add and delete a disk, detach, find it again with Scan, reattach, destroy. Needs XCP_LAB_HOSTS (pool
    /// coordinator first), XCP_ALLOW_WRITE=1 and the lab storage at 192.168.56.1; never point it at real servers.
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

        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void NfsSharedSrLifecycleOnLabPool() =>
            Lifecycle(SrRules.Kind.Nfs, "192.168.56.1:/export", null, null);

        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void SmbSharedSrLifecycleOnLabPool() =>
            Lifecycle(SrRules.Kind.Smb, @"\\192.168.56.1\sr", "xcp", Environment.GetEnvironmentVariable("XCP_PASS"));

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
