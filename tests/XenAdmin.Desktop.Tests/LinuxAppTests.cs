using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net.Security;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using XenAdmin.Actions.VMActions;
using XenAdmin.Desktop;
using XenAPI;

namespace XenAdmin.Desktop.Tests
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
        public void PoolMemberCertificateVouchedByXapi()
        {
            using var key = System.Security.Cryptography.RSA.Create(2048);
            var req = new System.Security.Cryptography.X509Certificates.CertificateRequest("CN=member", key, System.Security.Cryptography.HashAlgorithmName.SHA256, System.Security.Cryptography.RSASignaturePadding.Pkcs1);
            using var cert = req.CreateSelfSigned(DateTimeOffset.Now.AddDays(-1), DateTimeOffset.Now.AddDays(1));
            var fp = CertificateTrust.Fingerprint(cert);
            var xapiStyle = string.Join(":", Enumerable.Range(0, fp.Length / 2).Select(i => fp.Substring(i * 2, 2)));
            var trust = new CertificateTrust(TempConfig()) { VouchedFingerprints = h => h == "member:443" ? new[] { xapiStyle } : null };

            Assert.IsTrue(trust.ValidateHost("member:443", cert, SslPolicyErrors.RemoteCertificateChainErrors), "XAPI-reported fingerprint is trusted");
            Assert.IsFalse(trust.ValidateHost("other:443", cert, SslPolicyErrors.RemoteCertificateChainErrors), "unknown host still needs the user");
            trust.Pin("pinned:443", "00");
            trust.VouchedFingerprints = _ => new[] { xapiStyle };
            Assert.IsFalse(trust.ValidateHost("pinned:443", cert, SslPolicyErrors.None), "a stored pin always wins over the vouch");
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
        public void TopologiesMatchWinForms()
        {
            CollectionAssert.AreEqual(new[] { (4L, 1L), (2L, 2L), (1L, 4L) }, VmRules.ValidTopologies(4, 0));
            CollectionAssert.AreEqual(new[] { (6L, 1L), (3L, 2L) }, VmRules.ValidTopologies(6, 2), "limited by cores per socket");
            Assert.AreEqual(0, VmRules.ValidTopologies(7, 1).Count(t => t.Sockets > VM.MAX_SOCKETS), "never more sockets than allowed");
            Assert.AreEqual("Hard Disk", BootPage.DeviceName('C'));
        }

        [TestMethod]
        public void TabsMatchWindowsOrder()
        {
            CollectionAssert.AreEqual(new[] { "General", "Memory", "Storage", "Networking", "Console", "Performance", "Snapshots", "Search" }, Details.Tabs(new VM()));
            Assert.IsFalse(Details.NotPorted(new VM(), "Console"));
            Assert.IsTrue(Details.NotPorted(new Host(), "Console"));
            Assert.IsFalse(Details.NotPorted(new VM(), "Snapshots"));
            Assert.IsFalse(Details.NotPorted(new VM(), "Search"));
            Assert.IsTrue(Details.NotPorted(new SR(), "Search"));
            Assert.IsTrue(Details.NotPorted(new VM(), "Performance"));
            CollectionAssert.AreEqual(new[] { "-", "-", "-", "-", "-", "-" }, OverviewColumns.Values(new Pool()).Select(v => v.Text).ToArray());
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
            XenAdmin.XenAdminConfigManager.Provider ??= new DesktopConfigProvider();
            var vm = new VM { Connection = new XenAdmin.Network.XenConnection(), VBDs = new List<XenRef<VBD>>() };
            Assert.IsTrue(new VMStartAction(vm, null, null).AllowHaUnprotect, "Windows default must stay unchanged");
            Assert.IsFalse(((VMStartAbstractAction)VmOperations.Create(vm, VmOp.Start)).AllowHaUnprotect);
            Assert.IsFalse(((VMStartAbstractAction)VmOperations.Create(vm, VmOp.Resume)).AllowHaUnprotect);
        }

        [TestMethod]
        public void ConsoleKeyReleaseRepeatsThePressedKeysym()
        {
            var pressed = new System.Collections.Generic.Dictionary<(Avalonia.Input.PhysicalKey, Avalonia.Input.Key), int>();
            Avalonia.Input.KeyEventArgs Ev(string symbol) => new Avalonia.Input.KeyEventArgs { Key = Avalonia.Input.Key.OemComma, PhysicalKey = Avalonia.Input.PhysicalKey.Comma, KeySymbol = symbol };
            Assert.AreEqual('<', KeySyms.Track(pressed, Ev("<"), true));
            Assert.AreEqual('<', KeySyms.Track(pressed, Ev(","), false), "Shift released first");
            Assert.AreEqual('%', KeySyms.Track(pressed, Ev("%"), true));
            Assert.AreEqual('%', KeySyms.Track(pressed, Ev(null), false), "X11 reports no symbol on key up");
            Assert.AreEqual(0, pressed.Count);
        }

        [TestMethod]
        public void ConsoleKeysWithoutSymbolUseTheUsCharacterAndHvmGetsScancodes()
        {
            // X11 gives no symbol for Shift+5: the keysym path must still send '%'
            var shift5 = new Avalonia.Input.KeyEventArgs { Key = Avalonia.Input.Key.D5, PhysicalKey = Avalonia.Input.PhysicalKey.Digit5, KeyModifiers = Avalonia.Input.KeyModifiers.Shift };
            Assert.AreEqual('%', KeySyms.From(shift5));
            Assert.AreEqual(0x06, KeySyms.Scancode(Avalonia.Input.PhysicalKey.Digit5));
            Assert.AreEqual(0x2a, KeySyms.Scancode(Avalonia.Input.PhysicalKey.ShiftLeft));
            Assert.AreEqual(0x33, KeySyms.Scancode(Avalonia.Input.PhysicalKey.Comma));
            Assert.AreEqual(0x9d, KeySyms.Scancode(Avalonia.Input.PhysicalKey.ControlRight), "0xE0 0x1D");
            Assert.AreEqual(0x73, KeySyms.Scancode(Avalonia.Input.PhysicalKey.IntlRo), "ABNT2 /? key");
        }
    }
}
