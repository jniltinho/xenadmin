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
