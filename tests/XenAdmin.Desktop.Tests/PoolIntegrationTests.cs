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
    /// Forms a pool from the two standalone lab hosts the way the New Pool dialog does, then ejects the member
    /// (XAPI resets and reboots it) so the lab is standalone again. Needs XCP_LAB_HOSTS=coordinator,member
    /// (tools/vagrant/xcp-ng-lab) and XCP_ALLOW_WRITE=1; never point it at real servers.
    /// </summary>
    [TestClass]
    public class PoolIntegrationTests
    {
        [ClassInitialize]
        public static void Init(TestContext ctx) => ConnectionIntegrationTests.Init(ctx);

        private static void Run(AsyncAction action, int minutes)
        {
            var done = new ManualResetEventSlim();
            action.Completed += _ => done.Set();
            action.RunAsync();
            Assert.IsTrue(done.Wait(TimeSpan.FromMinutes(minutes)), action.Title + " timed out");
            Assert.IsTrue(action.Succeeded, action.Exception?.ToString());
        }

        private static void Until(Func<bool> cond, int seconds, string what)
        {
            var deadline = DateTime.Now.AddSeconds(seconds);
            while (!cond() && DateTime.Now < deadline)
                Thread.Sleep(1000);
            Assert.IsTrue(cond(), what);
        }

        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void NewPoolThenEjectOnLab()
        {
            var lab = Environment.GetEnvironmentVariable("XCP_LAB_HOSTS")?.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);
            if (lab == null || lab.Length != 2 || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_LAB_HOSTS=coordinator,member and XCP_ALLOW_WRITE=1 required");
            var coordinatorConn = ConnectionIntegrationTests.Connect(lab[0]);
            var memberConn = ConnectionIntegrationTests.Connect(lab[1]);
            // XenConnection adds itself to ConnectionsManager once connected (HandleSuccessfulConnection)
            try
            {
                if (Helpers.GetPool(coordinatorConn) != null || Helpers.GetPool(memberConn) != null)
                    Assert.Inconclusive("the lab hosts are already a pool (the storage tests use it); eject the member to rerun");

                var candidates = PoolRules.Candidates(new IXenConnection[] { memberConn, coordinatorConn }, coordinatorConn);
                Assert.IsTrue(candidates[0].WillBeCoordinator && candidates[0].Connection == coordinatorConn, "coordinator sorts first");
                Assert.IsTrue(candidates[1].AllowedAsSupporter, "member: " + candidates[1].Description);
                Assert.AreEqual(Messages.POOL_NAME_EMPTY, PoolRules.Invalid(" ", coordinatorConn, 1));
                Assert.IsNull(PoolRules.Invalid("lab-pool", coordinatorConn, 1));

                var coordinator = Helpers.GetCoordinator(coordinatorConn);
                var member = Helpers.GetCoordinator(memberConn);
                var memberUuid = member.uuid;
                System.Console.WriteLine("warnings: " + string.Join(" | ", PoolRules.Warnings(new() { member }, coordinator)));
                Run(PoolRules.Create(coordinator, new() { member }, "lab-pool", "created by the desktop app test"), 10);

                bool Dropped() { lock (ConnectionsManager.ConnectionsLock) return !ConnectionsManager.XenConnections.Contains(memberConn); }
                Until(Dropped, 30, "member connection dropped after the join");
                Until(() => coordinatorConn.Cache.Hosts.Length == 2 && Helpers.GetPool(coordinatorConn)?.Name() == "lab-pool", 120, "pool of two in the cache");
                var pool = Helpers.GetPool(coordinatorConn);
                Assert.AreEqual("created by the desktop app test", pool.Description());

                var joined = coordinatorConn.Cache.Hosts.Single(h => h.uuid == memberUuid);
                Until(() => coordinatorConn.Resolve(new XenRef<Host>(joined.opaque_ref))?.IsLive() == true, 300, "member live in the pool");
                Run(new EjectHostFromPoolAction(pool, joined), 10);
                Until(() => coordinatorConn.Cache.Hosts.Length == 1, 120, "member gone from the pool");
                var single = coordinatorConn.Cache.Pools.Single();
                Run(new DelegatedAsyncAction(coordinatorConn, "reset lab pool", "", "", s =>
                {
                    Pool.set_name_label(s, single.opaque_ref, "");
                    Pool.set_name_description(s, single.opaque_ref, "");
                }, true), 2);
            }
            finally
            {
                // a failed join leaves the coordinator as a named pool of one, which counts as a pool
                if (coordinatorConn.IsConnected && coordinatorConn.Cache.Hosts.Length == 1 && coordinatorConn.Cache.Pools.FirstOrDefault() is Pool one && one.name_label != "")
                    Pool.set_name_label(coordinatorConn.DuplicateSession(), one.opaque_ref, "");
                lock (ConnectionsManager.ConnectionsLock)
                {
                    ConnectionsManager.XenConnections.Remove(coordinatorConn);
                    ConnectionsManager.XenConnections.Remove(memberConn);
                }
                coordinatorConn.EndConnect();
                memberConn.EndConnect();
            }
        }
    }
}
