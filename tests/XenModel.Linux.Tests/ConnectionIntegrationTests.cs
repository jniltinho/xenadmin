using System;
using System.Collections.Concurrent;
using System.ComponentModel;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using XenAdmin;
using XenAdmin.Linux;
using XenAdmin.Network;
using XenAPI;

namespace XenModel.Linux.Tests
{
    /// <summary>A dedicated "UI" thread so XenModel's InvokeHelper behaves as in the app.</summary>
    public sealed class ThreadSynchronizer : ISynchronizeInvoke
    {
        private readonly BlockingCollection<Action> _queue = new BlockingCollection<Action>();
        private readonly Thread _thread;

        public ThreadSynchronizer()
        {
            _thread = new Thread(() => { foreach (var a in _queue.GetConsumingEnumerable()) a(); }) { IsBackground = true, Name = "test-ui" };
            _thread.Start();
        }

        public bool InvokeRequired => Thread.CurrentThread != _thread;

        public IAsyncResult BeginInvoke(Delegate method, object[] args)
        {
            var tcs = new TaskCompletionSource<object>();
            _queue.Add(() => { try { tcs.SetResult(method.DynamicInvoke(args)); } catch (Exception e) { tcs.SetException(e); } });
            return tcs.Task;
        }

        public object EndInvoke(IAsyncResult result) => ((Task<object>)result).GetAwaiter().GetResult();
        public object Invoke(Delegate method, object[] args) => InvokeRequired ? EndInvoke(BeginInvoke(method, args)) : method.DynamicInvoke(args);
    }

    [TestClass]
    public class ConnectionIntegrationTests
    {
        internal static ThreadSynchronizer Sync;

        [ClassInitialize]
        public static void Init(TestContext _)
        {
            XenAdminConfigManager.Provider ??= new LinuxConfigProvider();
            Sync ??= new ThreadSynchronizer();
            InvokeHelper.Initialize(Sync);
#pragma warning disable SYSLIB0014
            // lab hosts use self-signed certificates; the app uses CertificateTrust instead
            System.Net.ServicePointManager.ServerCertificateValidationCallback = (_, _, _, _) => true;
#pragma warning restore SYSLIB0014
        }

        internal static XenConnection Connect(string host)
        {
            var c = new XenConnection
            {
                Hostname = host, Port = 443,
                Username = Environment.GetEnvironmentVariable("XCP_USER"),
                Password = Environment.GetEnvironmentVariable("XCP_PASS")
            };
            ConnectionResultEventArgs result = null;
            var done = new ManualResetEventSlim();
            c.ConnectionResult += (_, e) => { result = e; done.Set(); };
            Sync.Invoke(new Action(() => c.BeginConnect(false, (_, _) => false)), null);
            Assert.IsTrue(done.Wait(TimeSpan.FromSeconds(60)), "timeout connecting to " + host);
            Assert.IsTrue(result.Connected, $"{host}: {result.Reason} {result.Error}");
            var deadline = DateTime.Now.AddSeconds(60);
            while (!c.CacheIsPopulated && DateTime.Now < deadline)
                Thread.Sleep(200);
            Assert.IsTrue(c.CacheIsPopulated, "cache not populated");
            return c;
        }

        private static void ReadOnlyChecks(string envVar)
        {
            var host = Environment.GetEnvironmentVariable(envVar);
            if (string.IsNullOrEmpty(host))
                Assert.Inconclusive(envVar + " not set");
            var c = Connect(host);
            try
            {
                var tree = ResourceTree.Build(c);
                // the tools ISO SR is hidden, like in the WinForms tree
                Assert.IsFalse(Flatten(tree).Any(n => n.XenObject is SR sr && sr.IsToolsSR()));
                var hosts = c.Cache.Hosts;
                Assert.IsTrue(hosts.Length > 0);
                var vms = c.Cache.VMs.Where(ResourceTree.IsRealVm).ToList();
                var pool = XenAdmin.Core.Helpers.GetPool(c);
                System.Console.WriteLine($"{host}: pool={pool?.Name() ?? "-"} coordinator={c.Hostname} hosts={string.Join(",", hosts.Select(h => h.Name()))} vms={vms.Count} srs={c.Cache.SRs.Length}");
                // every real VM must be somewhere in the tree, exactly once
                var flat = Flatten(tree).Where(n => n.XenObject is VM v && ResourceTree.IsRealVm(v)).Select(n => n.XenObject.opaque_ref).ToList();
                CollectionAssert.AreEquivalent(vms.Select(v => v.opaque_ref).ToList(), flat);
                if (hosts.Length > 1)
                    Assert.IsNotNull(pool, "multi-host connection should show a pool");
                foreach (var vm in vms.Where(v => v.power_state == vm_power_state.Running))
                    Assert.IsNotNull(c.Resolve(vm.resident_on), $"running VM {vm.Name()} without resident host");
                foreach (var o in new IXenObject[] { hosts[0], vms.FirstOrDefault(), c.Cache.SRs.FirstOrDefault(), pool }.Where(o => o != null))
                {
                    Assert.IsTrue(Details.General(o).Any(sec => sec.Rows.Count > 0), "General tab empty for " + o.GetType().Name);
                    foreach (var tab in Details.Tabs(o).Where(t => !Details.NotPorted(o, t)))
                    {
                        if (tab == "Storage") Assert.IsNotNull(Details.Storage(o));
                        if (tab == "Networking") Assert.IsNotNull(Details.Networking(o));
                        if (tab == "Memory") Details.Memory(o);
                    }
                }
            }
            finally
            {
                c.EndConnect();
            }
        }

        private static System.Collections.Generic.IEnumerable<TreeNode> Flatten(TreeNode n) => new[] { n }.Concat(n.Children.SelectMany(Flatten));

        [TestMethod, TestCategory("Integration")]
        public void StandaloneHostReadOnly() => ReadOnlyChecks("XCP_HOST");

        [TestMethod, TestCategory("Integration")]
        public void PoolReadOnly() => ReadOnlyChecks("XCP_RO_HOST");
    }
}
