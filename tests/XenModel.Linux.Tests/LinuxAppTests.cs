using System;
using System.Collections.Generic;
using System.IO;
using System.Net.Security;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using XenAdmin.Actions.VMActions;
using XenAdmin.Linux;
using XenAPI;

namespace XenModel.Linux.Tests
{
    [TestClass]
    public class LinuxAppTests
    {
        private static AppConfig TempConfig() =>
            AppConfig.Load(Path.Combine(Path.GetTempPath(), Guid.NewGuid() + "/config.json"));

        [TestMethod]
        public void TofuPolicy()
        {
            var trust = new CertificateTrust(TempConfig());
            Assert.AreEqual(TrustDecision.Trusted, trust.Evaluate("h:443", "AA", SslPolicyErrors.None));
            Assert.AreEqual(TrustDecision.Prompt, trust.Evaluate("h:443", "AA", SslPolicyErrors.RemoteCertificateChainErrors));
            trust.Pin("h:443", "AA");
            Assert.AreEqual(TrustDecision.Trusted, trust.Evaluate("H:443", "aa", SslPolicyErrors.RemoteCertificateChainErrors));
            // a stored pin wins even over a system-trusted chain
            Assert.AreEqual(TrustDecision.Mismatch, trust.Evaluate("h:443", "BB", SslPolicyErrors.None));
        }

        [TestMethod]
        public void ValidateNeverBlocksAndRecordsPending()
        {
            using var key = System.Security.Cryptography.RSA.Create(2048);
            var req = new System.Security.Cryptography.X509Certificates.CertificateRequest("CN=xcp", key, System.Security.Cryptography.HashAlgorithmName.SHA256, System.Security.Cryptography.RSASignaturePadding.Pkcs1);
            using var cert = req.CreateSelfSigned(DateTimeOffset.Now.AddDays(-1), DateTimeOffset.Now.AddDays(1));
#pragma warning disable SYSLIB0014
            var sender = System.Net.WebRequest.CreateHttp("https://xcp1:443/jsonrpc");
#pragma warning restore SYSLIB0014
            var trust = new CertificateTrust(TempConfig());

            Assert.IsFalse(trust.Validate(sender, cert, null, SslPolicyErrors.RemoteCertificateChainErrors));
            Assert.IsTrue(trust.Pending.TryGetValue("xcp1:443", out var p));
            Assert.AreEqual(CertificateTrust.Fingerprint(cert), p.Fingerprint);
            Assert.IsNull(p.Old);

            trust.Pin("xcp1:443", p.Fingerprint);
            Assert.IsTrue(trust.Validate(sender, cert, null, SslPolicyErrors.RemoteCertificateChainErrors));
        }

        [TestMethod]
        public void ConfigRoundTripHasNoPassword()
        {
            var cfg = TempConfig();
            cfg.AddOrUpdateServer("xcp1", 443, "root");
            cfg.Fingerprints["xcp1:443"] = "AB";
            cfg.Save();
            var json = File.ReadAllText(cfg.Path);
            StringAssert.DoesNotMatch(json, new System.Text.RegularExpressions.Regex("(?i)pass"));
            var again = AppConfig.Load(cfg.Path);
            Assert.AreEqual("xcp1", again.Servers[0].Host);
            Assert.AreEqual("AB", again.Fingerprints["XCP1:443"]);
        }

        [TestMethod]
        public void FilterKeepsAncestors()
        {
            var root = new TreeNode { Name = "pool" };
            var host = new TreeNode { Name = "host1" };
            host.Children.Add(new TreeNode { Name = "Web-01" });
            host.Children.Add(new TreeNode { Name = "db" });
            root.Children.Add(host);
            var f = ResourceTree.Filter(root, "web");
            Assert.AreEqual("host1", f.Children[0].Name);
            Assert.AreEqual(1, f.Children[0].Children.Count);
            Assert.IsNull(ResourceTree.Filter(root, "nothing"));
            Assert.AreSame(root, ResourceTree.Filter(root, " "));
        }

        [TestMethod]
        public void TabsMatchWindowsOrder()
        {
            CollectionAssert.AreEqual(new[] { "General", "Memory", "Storage", "Networking", "Console", "Performance", "Snapshots", "Search" }, Details.Tabs(new VM()));
            Assert.IsFalse(Details.NotPorted(new VM(), "Console"));
            Assert.IsTrue(Details.NotPorted(new Host(), "Console"));
            Assert.IsTrue(Details.NotPorted(new VM(), "Snapshots"));
        }

        [TestMethod]
        public void TreeExpansionSurvivesRebuild()
        {
            var a = new TreeNode { Name = "pool-x" };
            Assert.IsFalse(a.IsExpanded);
            a.IsExpanded = true;
            Assert.IsTrue(new TreeNode { Name = "pool-x" }.IsExpanded, "state is kept by key, not by node instance");
            a.IsExpanded = false;
            Assert.IsTrue(new TreeNode { IsRoot = true }.IsExpanded, "root starts expanded");
        }

        [TestMethod]
        public void OperationsFollowAllowedOperations()
        {
            var vm = new VM { allowed_operations = new List<vm_operations> { vm_operations.start } };
            Assert.IsTrue(VmOperations.IsAllowed(vm, VmOp.Start));
            Assert.IsFalse(VmOperations.IsAllowed(vm, VmOp.Shutdown));
            Assert.IsFalse(VmOperations.IsAllowed(null, VmOp.Start));
            Assert.IsTrue(VmOperations.NeedsConfirmation(VmOp.ForceShutdown));
            Assert.IsTrue(VmOperations.NeedsConfirmation(VmOp.ForceReboot));
            Assert.IsFalse(VmOperations.NeedsConfirmation(VmOp.Shutdown));
        }

        [TestMethod]
        public void StartActionsNeverUnprotectHa()
        {
            XenAdmin.XenAdminConfigManager.Provider ??= new LinuxConfigProvider();
            var vm = new VM { Connection = new XenAdmin.Network.XenConnection(), VBDs = new List<XenRef<VBD>>() };
            Assert.IsTrue(new VMStartAction(vm, null, null).AllowHaUnprotect, "Windows default must stay unchanged");
            Assert.IsFalse(((VMStartAbstractAction)VmOperations.Create(vm, VmOp.Start)).AllowHaUnprotect);
            Assert.IsFalse(((VMStartAbstractAction)VmOperations.Create(vm, VmOp.Resume)).AllowHaUnprotect);
        }
    }
}
