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

        /// <summary>
        /// New Network wizard, External network: a VLAN on a standalone lab host, then the two hosts form a pool and a VLAN
        /// is created on the pool (pool.create_VLAN_from_PIF puts it on every host). Both VLANs are removed; the pool is
        /// kept for the storage tests. Needs the standalone lab (XCP_LAB_HOSTS=coordinator,member).
        /// </summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void ExternalVlanOnHostThenPoolOnLab()
        {
            var lab = Environment.GetEnvironmentVariable("XCP_LAB_HOSTS")?.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);
            if (lab == null || lab.Length != 2 || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_LAB_HOSTS=coordinator,member and XCP_ALLOW_WRITE=1 required");
            var coordinatorConn = ConnectionIntegrationTests.Connect(lab[0]);
            // a rerun on the lab pool: the member answers HOST_IS_SLAVE, so it is connected only while standalone
            var memberConn = Helpers.GetPool(coordinatorConn) == null ? ConnectionIntegrationTests.Connect(lab[1]) : null;
            try
            {
                // a VLAN on a standalone host, as the wizard does from the host's Networking tab, then the New Pool dialog;
                // a rerun on the lab pool goes straight to the pool part
                var host = NetworkRules.WizardHost(Helpers.GetCoordinator(coordinatorConn));
                if (memberConn != null && Helpers.GetPool(memberConn) == null)
                {
                    CreateAndRemoveVlan(coordinatorConn, host, 101, 1);
                    Run(PoolRules.Create(host, new() { Helpers.GetCoordinator(memberConn) }, "lab-pool", "created by the desktop app test"), 10);
                }
                // the join restarts XAPI on the member and the coordinator's connection reconnects: wait until it is back
                Until(() => coordinatorConn.IsConnected && coordinatorConn.CacheIsPopulated && coordinatorConn.Cache.Hosts.Length == 2
                            && coordinatorConn.Cache.Hosts.All(h => h.IsLive()) && Helpers.GetPool(coordinatorConn)?.Name() == "lab-pool", 300, "pool of two, connected");
                Thread.Sleep(TimeSpan.FromSeconds(20));
                Until(() => coordinatorConn.IsConnected, 120, "coordinator connection stable");
                var pool = Helpers.GetPool(coordinatorConn);
                Assert.AreEqual(Helpers.GetCoordinator(coordinatorConn).opaque_ref, NetworkRules.WizardHost(pool).opaque_ref, "the pool's wizard offers the coordinator's NICs");
                CreateAndRemoveVlan(coordinatorConn, pool, 102, 2);
            }
            finally
            {
                lock (ConnectionsManager.ConnectionsLock)
                {
                    ConnectionsManager.XenConnections.Remove(coordinatorConn);
                    if (memberConn != null) ConnectionsManager.XenConnections.Remove(memberConn);
                }
                coordinatorConn.EndConnect();
                memberConn?.EndConnect();
            }
        }

        private static void CreateAndRemoveVlan(IXenConnection c, IXenObject owner, int vlan, int hosts)
        {
            var nic = NetworkRules.WizardNics(NetworkRules.WizardHost(owner)).First();
            var inUse = NetworkRules.VlansInUse(nic);
            Assert.AreEqual(vlan, NetworkRules.VlanFor(vlan, inUse), "VLAN free on " + nic.Name());
            Assert.IsFalse(NetworkRules.ValidateVlan(vlan.ToString(), NetworkRules.MinVlan(c), inUse).Error);
            var name = "xcpng-center-linux-test-vlan" + vlan;
            Run(new NetworkAction(c, NetworkRules.NewNetwork(name, "created by the Linux New Network wizard test", false, 1500), nic, vlan), 3);
            Until(() => NetworkRules.HostNetworks(c).Any(n => n.name_label == name && c.ResolveAll(n.PIFs).Count == hosts), 60, $"VLAN {vlan} on {hosts} host(s)");
            var net = NetworkRules.HostNetworks(c).Single(n => n.name_label == name);
            Assert.IsTrue(c.ResolveAll(net.PIFs).All(p => p.VLAN == vlan && !p.IsPhysical()));
            var row = Details.Networking(owner).Rows.Single(r => r[0] == name);
            Assert.AreEqual(vlan.ToString(), row[3], "VLAN column");
            Assert.AreEqual(Messages.NO, row[4], "Auto column");
            Assert.IsTrue(NetworkRules.CanRemoveNetwork(owner, net).Enabled, "a VLAN network can be removed");
            System.Console.WriteLine($"VLAN {vlan} on {owner.GetType().Name} '{owner.Name()}': {string.Join(" | ", row)}");
            Run(new NetworkAction(c, net, false), 3);
            Until(() => NetworkRules.HostNetworks(c).All(n => n.name_label != name), 60, $"VLAN {vlan} removed");
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
