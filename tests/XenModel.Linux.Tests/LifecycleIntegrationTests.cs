using System;
using System.Linq;
using System.Threading;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using XenAdmin.Actions;
using XenAdmin.Linux;
using XenAPI;

namespace XenModel.Linux.Tests
{
    /// <summary>
    /// Creates a throwaway diskless VM on XCP_HOST, drives it through the same actions the app uses, then destroys it.
    /// Runs only with XCP_ALLOW_WRITE=1.
    /// </summary>
    [TestClass]
    public class LifecycleIntegrationTests
    {
        private const string VmName = "xcpng-center-linux-test";

        [ClassInitialize]
        public static void Init(TestContext ctx) => ConnectionIntegrationTests.Init(ctx);

        private static void Run(AsyncAction action)
        {
            var done = new ManualResetEventSlim();
            action.Completed += _ => done.Set();
            action.RunAsync();
            Assert.IsTrue(done.Wait(TimeSpan.FromMinutes(3)), action.Title + " timed out");
        }

        private static VM WaitFor(XenAdmin.Network.XenConnection c, string vmRef, Func<VM, bool> cond)
        {
            var deadline = DateTime.Now.AddSeconds(60);
            VM vm;
            while (!((vm = c.Resolve(new XenRef<VM>(vmRef))) != null && cond(vm)) && DateTime.Now < deadline)
                Thread.Sleep(250);
            Assert.IsTrue(vm != null && cond(vm), "cache did not reflect the expected VM state");
            return vm;
        }

        /// <summary>
        /// Saves properties the way PropertiesDialog does (SaveChangesAction on a clone + page actions)
        /// and checks the result through the cache, then destroys the VM.
        /// </summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void PropertiesRoundTripOnThrowawayVm()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            string vmRef = null;
            try
            {
                vmRef = CreateThrowawayVm(c, session, VmName + "-props");
                var vm = WaitFor(c, vmRef, v => !v.is_a_template);

                var copy = (VM)vm.Clone();
                copy.name_label = VmName + "-renamed";
                copy.name_description = "edited by the Linux Properties dialog test";
                copy.SetAutoPowerOn(true);
                copy.SetBootOrder("DC");
                copy.SetVcpuWeight(1024);
                copy.SetCoresPerSocket(2);
                var actions = new System.Collections.Generic.List<AsyncAction>
                {
                    new SaveChangesAction(copy, true, vm),
                    new XenAdmin.Actions.GeneralEditPageAction(vm, copy, "/linux-test", new System.Collections.Generic.List<string> { "linux", "test" }, true),
                    new XenAdmin.Actions
.ChangeVCPUSettingsAction(copy, 2, 2)
                };
                var all = new MultipleAction(c, "Update properties", "Updating", "Updated", actions);
                vm.Locked = true; // as PropertiesDialog/ShowProperties do: SaveChanges requires the original locked
                Run(all);
                vm.Locked = false;
                Assert.IsTrue(all.Succeeded, all.Exception?.ToString());

                vm = WaitFor(c, vmRef, v => v.name_label == VmName + "-renamed" && v.VCPUs_max == 2 && v.VCPUs_at_startup == 2 && v.Path == "/linux-test");
                System.Console.WriteLine($"name={vm.name_label} vcpus={vm.VCPUs_max}/{vm.VCPUs_at_startup} path='{vm.Path}' folder='{(vm.other_config.TryGetValue("folder", out var f) ? f : null)}' tags={string.Join(",", vm.tags)} boot={vm.GetBootOrder()} weight={vm.GetVcpuWeight()} cores={vm.GetCoresPerSocket()} auto={vm.GetAutoPowerOn()}");
                Assert.AreEqual("edited by the Linux Properties dialog test", vm.name_description);
                Assert.IsTrue(vm.GetAutoPowerOn());
                Assert.AreEqual("DC", vm.GetBootOrder().ToUpper());
                Assert.AreEqual(1024, vm.GetVcpuWeight());
                Assert.AreEqual(2, vm.GetCoresPerSocket());
                CollectionAssert.AreEquivalent(new[] { "linux", "test" }, XenAdmin.Model.Tags.GetTags(vm));
                CollectionAssert.AreEquivalent(new[] { (2L, 1L), (1L, 2L) }, XenAdmin.Linux.CpuPage.ValidTopologies(2, vm.MaxCoresPerSocket()));
            }
            finally
            {
                if (vmRef != null)
                {
                    try { XenAdmin.Model.Folders.Unfolder(session, c.Resolve(new XenRef<VM>(vmRef))); } catch { }
                    VM.destroy(session, vmRef);
                }
                c.EndConnect();
            }
        }

        /// <summary>Take snapshot, revert, delete — the Snapshots tab actions — on a throwaway VM.</summary>
        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void SnapshotTakeRevertDelete()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");
            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            string vmRef = null;
            try
            {
                vmRef = CreateThrowawayVm(c, session, VmName + "-snap");
                var vm = WaitFor(c, vmRef, v => !v.is_a_template);
                Assert.IsTrue(vm.allowed_operations.Contains(vm_operations.snapshot));

                var take = new VMSnapshotCreateAction(vm, "linux-test-snapshot", "taken by the Linux Snapshots tab test", SnapshotType.DISK, null);
                Run(take);
                Assert.IsTrue(take.Succeeded, take.Exception?.ToString());
                // the description is set after the snapshot task completes, so wait for it too
                vm = WaitFor(c, vmRef, v => v.snapshots.Count == 1 && c.Resolve(v.snapshots[0])?.name_description == "taken by the Linux Snapshots tab test");
                var snap = c.Resolve(vm.snapshots[0]);
                Assert.AreEqual("linux-test-snapshot", snap.name_label);
                Assert.AreEqual("taken by the Linux Snapshots tab test", snap.name_description);
                Assert.AreEqual(XenAdmin.Linux.SnapshotsView.TypeText(snap), XenAdmin.Messages.DISKS_ONLY);
                Assert.IsTrue(XenAdmin.Linux.SnapshotsView.Size(snap) >= 0);
                System.Console.WriteLine($"snapshot {snap.name_label} size={XenAdmin.Linux.SnapshotsView.Size(snap)} allowed={string.Join(",", snap.allowed_operations)}");

                var revert = new VMSnapshotRevertAction(snap);
                Run(revert);
                Assert.IsTrue(revert.Succeeded, revert.Exception?.ToString());

                var delete = new VMSnapshotDeleteAction(snap);
                Run(delete);
                Assert.IsTrue(delete.Succeeded, delete.Exception?.ToString());
                WaitFor(c, vmRef, v => v.snapshots.Count == 0);
            }
            finally
            {
                if (vmRef != null)
                {
                    foreach (var s in VM.get_snapshots(session, vmRef))
                        try { VM.destroy(session, s.opaque_ref); } catch { }
                    VM.destroy(session, vmRef);
                }
                c.EndConnect();
            }
        }

        private static string CreateThrowawayVm(XenAdmin.Network.XenConnection c, Session session, string name)
        {
            foreach (var old in c.Cache.VMs.Where(v => v.name_label.StartsWith(VmName) && !v.is_a_template))
            {
                if (old.power_state != vm_power_state.Halted) VM.hard_shutdown(session, old.opaque_ref);
                VM.destroy(session, old.opaque_ref);
            }
            var template = c.Cache.VMs.First(v => v.is_a_template && v.name_label == "Other install media");
            var vmRef = VM.clone(session, template.opaque_ref, name).opaque_ref;
            VM.set_is_a_template(session, vmRef, false);
            const long mem = 256L * 1024 * 1024;
            VM.set_memory_limits(session, vmRef, mem, mem, mem, mem);
            VM.set_VCPUs_max(session, vmRef, 1);
            VM.set_VCPUs_at_startup(session, vmRef, 1);
            return vmRef;
        }

        [TestMethod, TestCategory("Integration"), TestCategory("Write")]
        public void PowerCycleThrowawayVm()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host) || Environment.GetEnvironmentVariable("XCP_ALLOW_WRITE") != "1")
                Assert.Inconclusive("XCP_HOST and XCP_ALLOW_WRITE=1 required");

            var c = ConnectionIntegrationTests.Connect(host);
            var session = c.DuplicateSession();
            string vmRef = null;
            try
            {
                // leftovers from an aborted run
                foreach (var old in c.Cache.VMs.Where(v => v.name_label == VmName && !v.is_a_template))
                {
                    if (old.power_state != vm_power_state.Halted) VM.hard_shutdown(session, old.opaque_ref);
                    VM.destroy(session, old.opaque_ref);
                }

                var template = c.Cache.VMs.First(v => v.is_a_template && v.name_label == "Other install media");
                vmRef = VM.clone(session, template.opaque_ref, VmName).opaque_ref;
                VM.set_is_a_template(session, vmRef, false);
                const long mem = 256L * 1024 * 1024;
                VM.set_memory_limits(session, vmRef, mem, mem, mem, mem);
                VM.set_VCPUs_max(session, vmRef, 1);
                VM.set_VCPUs_at_startup(session, vmRef, 1);

                var vm = WaitFor(c, vmRef, v => v.power_state == vm_power_state.Halted && !v.is_a_template);
                Assert.IsTrue(VmOperations.IsAllowed(vm, VmOp.Start));
                Assert.IsFalse(VmOperations.IsAllowed(vm, VmOp.ForceShutdown), "halted VM cannot be force-shutdown");

                Run(VmOperations.Create(vm, VmOp.Start));
                vm = WaitFor(c, vmRef, v => v.power_state == vm_power_state.Running);
                Assert.IsNotNull(c.Resolve(vm.resident_on), "running VM must have a resident host");
                System.Console.WriteLine("allowed while running: " + string.Join(",", vm.allowed_operations));

                // expected failure: a diskless VM without guest tools cannot be suspended
                if (!VmOperations.IsAllowed(vm, VmOp.Suspend))
                {
                    var suspend = VmOperations.Create(vm, VmOp.Suspend);
                    Run(suspend);
                    Assert.IsFalse(suspend.Succeeded, "suspend should fail without guest tools");
                    Assert.IsInstanceOfType(suspend.Exception, typeof(Failure));
                    System.Console.WriteLine("expected suspend failure: " + suspend.Exception.Message);
                }

                Assert.IsTrue(VmOperations.IsAllowed(vm, VmOp.ForceReboot));
                var reboot = VmOperations.Create(vm, VmOp.ForceReboot);
                Run(reboot);
                Assert.IsTrue(reboot.Succeeded, reboot.Exception?.Message);
                vm = WaitFor(c, vmRef, v => v.power_state == vm_power_state.Running);

                var shutdown = VmOperations.Create(vm, VmOp.ForceShutdown);
                Run(shutdown);
                Assert.IsTrue(shutdown.Succeeded, shutdown.Exception?.Message);
                WaitFor(c, vmRef, v => v.power_state == vm_power_state.Halted);
            }
            finally
            {
                if (vmRef != null)
                {
                    try { if (VM.get_power_state(session, vmRef) != vm_power_state.Halted) VM.hard_shutdown(session, vmRef); } catch { }
                    VM.destroy(session, vmRef);
                }
                c.EndConnect();
            }
        }
    }
}
