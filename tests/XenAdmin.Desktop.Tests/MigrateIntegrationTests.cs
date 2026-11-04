using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAPI;

namespace XenAdmin.Desktop.Tests
{
    /// <summary>
    /// Migrate to Server on a pool with a diskless PV Alpine VM (lab.sh guest-kernel puts its kernel in every dom0).
    /// Needs XCP_LAB_HOSTS (coordinator first) and XCP_ALLOW_WRITE=1, and hosts that can run guests: the VirtualBox
    /// lab on AMD cannot (no nested NPT, and XCP-ng 8.3's Xen has no shadow paging), so there it is inconclusive.
    /// </summary>
    [TestClass]
    public class MigrateIntegrationTests
    {
        private const string VmName = "xcpng-center-test-migrate";

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

        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void MigrateRunningPvVmBetweenLabHosts()
        {
            var lab = Environment.GetEnvironmentVariable("XCP_LAB_HOSTS")?.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);
            if (lab == null || lab.Length == 0 || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_LAB_HOSTS and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(lab[0]);
            var session = c.DuplicateSession(); // shares the connection's session: never log it out
            string vmRef = null;
            try
            {
                if (c.Cache.Hosts.Length < 2)
                    Assert.Inconclusive("the lab hosts are not a pool yet");
                foreach (var stale in c.Cache.VMs.Where(v => v.name_label == VmName && !v.is_a_template))
                {
                    if (stale.power_state != vm_power_state.Halted) VM.hard_shutdown(session, stale.opaque_ref);
                    VM.destroy(session, stale.opaque_ref);
                }

                var template = c.Cache.VMs.Single(v => v.is_a_template && v.name_label == "Other install media");
                vmRef = VM.clone(session, template.opaque_ref, VmName).opaque_ref;
                VM.provision(session, vmRef);
                VM.set_is_a_template(session, vmRef, false);
                VM.set_memory_limits(session, vmRef, 256L << 20, 256L << 20, 256L << 20, 256L << 20);
                VM.set_HVM_boot_policy(session, vmRef, "");
                VM.set_domain_type(session, vmRef, domain_type.pv);
                VM.set_PV_bootloader(session, vmRef, "");
                VM.set_PV_kernel(session, vmRef, "/var/lib/xcp/guest/alpine-vmlinuz-virt");
                VM.set_PV_ramdisk(session, vmRef, "/var/lib/xcp/guest/alpine-initramfs-virt");
                VM.set_PV_args(session, vmRef, "console=hvc0");

                var hosts = c.Cache.Hosts.OrderBy(h => h.name_label).ToList();
                try { VM.start_on(session, vmRef, hosts[0].opaque_ref, false, true); }
                catch (Failure f) when (f.Message.Contains("shadow") || f.Message.Contains("paging mode"))
                {
                    // nested Xen without NPT (VirtualBox on AMD) and XCP-ng 8.3 without shadow paging: no guest can run
                    Assert.Inconclusive("the lab hosts cannot run guests: " + f.Message);
                }
                var vm = Until(() => c.Resolve(new XenRef<VM>(vmRef)), v => v?.power_state == vm_power_state.Running
                    && v.resident_on.opaque_ref == hosts[0].opaque_ref && v.allowed_operations.Contains(vm_operations.pool_migrate), "running on " + hosts[0].name_label);

                // the submenu: the current server disabled with its reason, the other one offered
                Assert.IsTrue(MigrateRules.CanMigrate(vm));
                var targets = MigrateRules.Targets(vm, session);
                Assert.AreEqual(Messages.HOST_MENU_CURRENT_SERVER, targets.Single(t => t.Host.opaque_ref == hosts[0].opaque_ref).Reason);
                var target = targets.Single(t => t.Host.opaque_ref == hosts[1].opaque_ref);
                Assert.IsNull(target.Reason, target.Reason);
                Assert.AreEqual($"{hosts[0].Name()} ({Messages.HOST_MENU_CURRENT_SERVER})", MigrateRules.MenuText(hosts[0], Messages.HOST_MENU_CURRENT_SERVER));

                Run(MigrateRules.Migrate(vm, hosts[1]));
                Until(() => c.Resolve(new XenRef<VM>(vmRef)), v => v.power_state == vm_power_state.Running && v.resident_on.opaque_ref == hosts[1].opaque_ref,
                    "running on " + hosts[1].name_label + " after the migration");
            }
            finally
            {
                if (vmRef != null)
                {
                    try { VM.hard_shutdown(session, vmRef); } catch (Failure) { }
                    VM.destroy(session, vmRef);
                }
                c.EndConnect();
            }
        }
    }
}
