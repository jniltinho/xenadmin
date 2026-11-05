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
    /// Migrate to Server on a pool: an HVM Alpine VM whose disk is the hybrid Alpine ISO on a shared SR, started on
    /// the first host and live-migrated to the second. Needs XCP_LAB_HOSTS (coordinator first), XCP_ALLOW_WRITE=1,
    /// XCP_LAB_GUEST_ISO (path to alpine-virt-*.iso) and a shared NFS/SMB SR on the pool; hosts must run guests
    /// (the KVM lab does; the VirtualBox lab on AMD cannot, and the test is then inconclusive).
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
        public void LiveMigrateVmBetweenPoolHosts()
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
                foreach (var stale in c.Cache.VDIs.Where(v => v.name_label == VmName))
                    VDI.destroy(session, stale.opaque_ref);

                var iso = Environment.GetEnvironmentVariable("XCP_LAB_GUEST_ISO");
                if (string.IsNullOrEmpty(iso) || !System.IO.File.Exists(iso))
                    Assert.Inconclusive("XCP_LAB_GUEST_ISO must point to the Alpine ISO");
                var shared = c.Cache.SRs.FirstOrDefault(sr => sr.shared && (sr.type == "nfs" || sr.type == "smb")
                    && c.ResolveAll(sr.PBDs).Count(p => p.currently_attached) == c.Cache.Hosts.Length);
                if (shared == null)
                    Assert.Inconclusive("no shared NFS/SMB SR on the pool (create one with the New Storage wizard)");

                // the ISO is hybrid: written to a disk it boots like one
                var size = (new System.IO.FileInfo(iso).Length + (1 << 20) - 1) >> 20 << 20;
                var vdi = VDI.create(session, new VDI
                {
                    name_label = VmName, SR = new XenRef<SR>(shared.opaque_ref), virtual_size = size, type = vdi_type.user,
                    sharable = false, read_only = false, other_config = new(), xenstore_data = new(), sm_config = new(), tags = new string[0]
                });
                XenAPI.HTTP_actions.put_import_raw_vdi(_ => { }, () => false, 600_000, c.Hostname, null, iso,
                    session_id: session.opaque_ref, vdi: vdi.opaque_ref, format: "raw");

                var template = c.Cache.VMs.Single(v => v.is_a_template && v.name_label == "Other install media");
                vmRef = VM.clone(session, template.opaque_ref, VmName).opaque_ref;
                VM.provision(session, vmRef);
                VM.set_is_a_template(session, vmRef, false);
                VM.set_memory_limits(session, vmRef, 512L << 20, 512L << 20, 512L << 20, 512L << 20);
                VM.set_HVM_boot_params(session, vmRef, new Dictionary<string, string> { ["order"] = "c" });
                VBD.create(session, new VBD
                {
                    VM = new XenRef<VM>(vmRef), VDI = new XenRef<VDI>(vdi.opaque_ref), userdevice = "0", bootable = true, mode = vbd_mode.RW,
                    type = vbd_type.Disk, empty = false, other_config = new(), qos_algorithm_params = new(), qos_algorithm_type = "", unpluggable = false
                });

                var hosts = c.Cache.Hosts.OrderBy(h => h.name_label).ToList();
                try { VM.start_on(session, vmRef, hosts[0].opaque_ref, false, true); }
                catch (Failure f) when (f.Message.Contains("shadow") || f.Message.Contains("paging mode") || f.Message.Contains("Shadow"))
                {
                    // nested Xen without NPT (VirtualBox on AMD) and XCP-ng 8.3 without shadow paging: no guest can run
                    Assert.Inconclusive("the lab hosts cannot run guests: " + f.Message);
                }
                var vm = Until(() => c.Resolve(new XenRef<VM>(vmRef)), v => v?.power_state == vm_power_state.Running
                    && v.resident_on.opaque_ref == hosts[0].opaque_ref && v.allowed_operations.Contains(vm_operations.pool_migrate), "running on " + hosts[0].name_label + " with PV drivers", 240);

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
                    foreach (var vbd in VM.get_VBDs(session, vmRef))
                    {
                        var disk = VBD.get_VDI(session, vbd.opaque_ref);
                        if (disk.opaque_ref != Helper.NullOpaqueRef && VBD.get_type(session, vbd.opaque_ref) == vbd_type.Disk)
                            VDI.destroy(session, disk.opaque_ref);
                    }
                    VM.destroy(session, vmRef);
                }
                c.EndConnect();
            }
        }
    }
}
