using System;
using System.Linq;
using System.Threading;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAPI;

namespace XenAdmin.Desktop.Tests
{
    /// <summary>
    /// Server commands on the lab pool: Enter Maintenance Mode evacuates a member (moving a VM on shared storage away),
    /// Exit Maintenance Mode restores it, Reboot brings the member back. Needs XCP_LAB_HOSTS (coordinator first),
    /// XCP_ALLOW_WRITE=1 and a pool of two; never point it at real servers.
    /// </summary>
    [TestClass]
    public class HostIntegrationTests
    {
        [ClassInitialize]
        public static void Init(TestContext ctx) => ConnectionIntegrationTests.Init(ctx);

        private static void Run(AsyncAction action, int minutes = 5)
        {
            var done = new ManualResetEventSlim();
            action.Completed += _ => done.Set();
            action.RunAsync();
            Assert.IsTrue(done.Wait(TimeSpan.FromMinutes(minutes)), action.Title + " timed out");
            Assert.IsTrue(action.Succeeded, action.Title + ": " + action.Exception);
        }

        private static T Until<T>(Func<T> get, Func<T, bool> ok, string what, int seconds = 180)
        {
            var deadline = DateTime.Now.AddSeconds(seconds);
            T v;
            while (!ok(v = get()) && DateTime.Now < deadline)
                Thread.Sleep(1000);
            Assert.IsTrue(ok(v), what);
            return v;
        }

        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void MaintenanceModeAndRebootOnLabMember()
        {
            var lab = Environment.GetEnvironmentVariable("XCP_LAB_HOSTS")?.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);
            if (lab == null || lab.Length == 0 || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_LAB_HOSTS and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(lab[0]);
            Host member = null;
            Host Member() => c.Resolve(new XenRef<Host>(member.opaque_ref));
            try
            {
                if (c.Cache.Hosts.Length < 2)
                    Assert.Inconclusive("the lab hosts are not a pool yet");
                member = c.Cache.Hosts.First(h => !h.IsCoordinator());
                // a member still booting (an earlier run's reboot) is not live and comes back enabled by itself
                Until(Member, h => h.IsLive(), "member live", 900);
                try { Until(Member, h => h.enabled, "member enabled", 120); } catch (AssertFailedException) { }
                member = Member();
                if (HostRules.InMaintenance(member))
                    Run(HostRules.ExitMaintenance(member, false)); // left over by an interrupted run
                Until(Member, h => !HostRules.InMaintenance(h), "member enabled before the test");
                Assert.IsTrue(HostRules.CanToggleMaintenance(member) && !HostRules.InMaintenance(member), "member live and enabled");
                Assert.AreEqual(Messages.ENTER_MAINTENANCE_MODE_CONTEXT_MENU, HostRules.MaintenanceText(member));
                Assert.IsNull(HostRules.CannotEnterMaintenance(member));

                // a running VM whose disks are all on shared storage goes to the member first, so evacuation has work to do
                var vm = c.Cache.VMs.FirstOrDefault(v => v.IsRealVm() && v.power_state == vm_power_state.Running && string.IsNullOrEmpty(v.IsOnSharedStorage()));
                if (vm == null)
                    Assert.Inconclusive("needs a running VM on shared storage (the migrate test's Alpine on the lab NFS SR)");
                if (vm.resident_on.opaque_ref != member.opaque_ref)
                {
                    Run(MigrateRules.Migrate(vm, member));
                    Until(() => c.Resolve(new XenRef<VM>(vm.opaque_ref)), v => v.resident_on.opaque_ref == member.opaque_ref, "VM on the member");
                }

                Run(HostRules.EnterMaintenance(Member()));
                Until(Member, h => HostRules.InMaintenance(h), "member disabled");
                if (vm != null)
                    Until(() => c.Resolve(new XenRef<VM>(vm.opaque_ref)), v => v.power_state == vm_power_state.Running && v.resident_on.opaque_ref != member.opaque_ref, "VM evacuated");
                Assert.AreEqual(Messages.EXIT_MAINTENANCE_MODE_CONTEXT_MENU, HostRules.MaintenanceText(Member()));
                if (vm != null)
                    Until(Member, h => HostRules.VmsToRestore(h).Any(v => v.opaque_ref == vm.opaque_ref), "evacuated VM offered for restore", 60);

                Run(HostRules.ExitMaintenance(Member(), true));
                Until(Member, h => !HostRules.InMaintenance(h), "member enabled again");
                if (vm != null)
                    Until(() => c.Resolve(new XenRef<VM>(vm.opaque_ref)), v => v.resident_on.opaque_ref == member.opaque_ref, "VM restored to the member");

                // Reboot: the confirmation names the server; the member comes back live and enabled
                var (title, text) = HostRules.Confirm(Member(), reboot: true);
                Assert.AreEqual(Messages.CONFIRM_REBOOT_SERVER_TITLE, title);
                StringAssert.Contains(text, member.Name());
                if (vm != null)
                    Run(MigrateRules.Migrate(c.Resolve(new XenRef<VM>(vm.opaque_ref)), c.Cache.Hosts.First(h => h.IsCoordinator())));
                var bootedAt = Member().BootTime();
                Run(HostRules.Reboot(Member()), 10);
                // a new boot time proves the reboot happened (live alone may still be the old heartbeat)
                Until(Member, h => h.BootTime() > bootedAt && h.IsLive() && h.enabled, "member rebooted and back", 900);
            }
            finally
            {
                // never leave the lab member in maintenance mode
                if (member != null && c.IsConnected && Member() is Host h && HostRules.InMaintenance(h) && h.IsLive())
                    try { Run(HostRules.ExitMaintenance(h, false)); } catch (AssertFailedException) { }
                c.EndConnect();
            }
        }
    }
}
