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
            Assert.IsFalse(Details.NotPorted(new VM(), "Performance"));
            Assert.IsFalse(Details.NotPorted(new Host(), "Performance"));
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

        [TestMethod]
        public void RrdUpdatesParseScaleAndUsedMemory()
        {
            const string xport = @"<xport><meta><start>1000</start><step>5</step><end>1010</end><rows>2</rows><columns>4</columns>
<legend><entry>AVERAGE:vm:u1:cpu0</entry><entry>AVERAGE:vm:u1:memory</entry><entry>AVERAGE:vm:u1:memory_internal_free</entry><entry>AVERAGE:vm:u1:vif_0_rx</entry></legend></meta>
<data><row><t>1010</t><v>0.25</v><v>536870912</v><v>262144</v><v>Infinity</v></row><row><t>1005</t><v>0.5</v><v>536870912</v><v>131072</v><v>2048</v></row></data></xport>";
            var raw = PerformanceRules.ParseUpdates(new System.IO.MemoryStream(System.Text.Encoding.UTF8.GetBytes(xport)));
            Assert.AreEqual(4, raw.Count);
            Assert.AreEqual(1005, new DateTimeOffset(raw["vm:u1:cpu0"][0].Time).ToUnixTimeSeconds(), "rows sorted by time");

            var units = new Dictionary<string, string> { ["cpu0"] = "(fraction)", ["memory"] = "B", ["memory_internal_free"] = "KiB", ["vif_0_rx"] = "B/s" };
            var data = PerformanceRules.Apply(raw, units);
            Assert.AreEqual(50, data["vm:u1:cpu0"][0].Value, 1e-9, "fraction shown as percent");
            // used = memory - free (KiB -> B): 512 MiB - 128 MiB
            Assert.AreEqual(384L << 20, data["vm:u1:memory_internal_free"][0].Value, 1);
            Assert.IsTrue(double.IsNaN(data["vm:u1:vif_0_rx"][1].Value), "infinities become gaps");

            // DataRange.RoundToNearestPowerOf10
            Assert.AreEqual(1, PerformanceRules.Scale(PerformanceRules.Unit.Percent, 0, 0.6).Max, "auto ranges start at 1");
            Assert.AreEqual(10, PerformanceRules.Scale(PerformanceRules.Unit.Percent, 0, 1.4).Max);
            Assert.AreEqual(100, PerformanceRules.Scale(PerformanceRules.Unit.Percent, 0, 42).Max);
            Assert.AreEqual(512L << 20, PerformanceRules.Scale(PerformanceRules.Unit.Bytes, 512L << 20, 1).Max, "memory keeps its fixed max");
            var axis = PerformanceRules.Scale(PerformanceRules.Unit.BytesPerSecond, 0, 49_000);
            Assert.AreEqual(100 * 1024, axis.Max, "10^n * 1024^k");
            Assert.AreEqual("kB/s", axis.Label);
            Assert.AreEqual(1024, PerformanceRules.Scale(PerformanceRules.Unit.BytesPerSecond, 0, 3).Max, "at least 1 kB/s");
            Assert.AreEqual(1024 * 1024, PerformanceRules.Scale(PerformanceRules.Unit.BytesPerSecond, 0, 1000 * 1024).Max, "1000 kB/s is shown as 1 MB/s");
            Assert.AreEqual("MB", PerformanceRules.AxisUnit(PerformanceRules.Unit.Bytes, 512L << 20).Label);
            Assert.AreEqual("kB/s", PerformanceRules.AxisUnit(PerformanceRules.Unit.BytesPerSecond, 1024).Label);
        }
    

        [TestMethod]
        public void ServerAddressParsing()
        {
            Assert.IsTrue(ServerAddress.TryParse(" xcp1.lan ", out var h, out var p)); Assert.AreEqual(("xcp1.lan", 443), (h, p));
            Assert.IsTrue(ServerAddress.TryParse("10.0.0.5:8443", out h, out p)); Assert.AreEqual(("10.0.0.5", 8443), (h, p));
            Assert.IsTrue(ServerAddress.TryParse("2001:db8::11", out h, out p)); Assert.AreEqual(("2001:db8::11", 443), (h, p), "bare IPv6 has no port");
            Assert.IsTrue(ServerAddress.TryParse("[2001:db8::11]:4443", out h, out p)); Assert.AreEqual(("2001:db8::11", 4443), (h, p));
            Assert.IsFalse(ServerAddress.TryParse("", out _, out _));
            Assert.IsFalse(ServerAddress.TryParse("host:99999", out _, out _));
            Assert.IsFalse(ServerAddress.TryParse("host:abc", out _, out _));
            Assert.IsTrue(ServerAddress.IsPlainHttp(80) && ServerAddress.IsPlainHttp(8080) && !ServerAddress.IsPlainHttp(443));

            // the parsed address reaches XenAPI as a valid URL (no network: the constructor only builds the client)
            ServerAddress.TryParse("[2001:db8::11]:4443", out h, out p);
            var url = new Uri(new Session(h, p).Url);
            Assert.AreEqual(("[2001:db8::11]", 4443, "https"), (url.Host, url.Port, url.Scheme));
            Assert.AreEqual("https://10.0.0.5:443", new Session("10.0.0.5", 443).Url);
        }

        [TestMethod]
        public void AlertPriorityAndIcon()
        {
            Assert.AreEqual(XenAdmin.Alerts.AlertPriority.Priority5, NotificationRules.PriorityOf(new Message { priority = 5 }));
            Assert.AreEqual(XenAdmin.Alerts.AlertPriority.Unknown, NotificationRules.PriorityOf(new Message { priority = 42 }), "out of range");
            Assert.AreEqual(XenAdmin.Alerts.AlertPriority.Unknown, NotificationRules.PriorityOf(new Message { priority = -1 }));
            Assert.AreEqual("alert1_16", NotificationRules.Icon(XenAdmin.Alerts.AlertPriority.Priority1));
            Assert.AreEqual("alert6_16", NotificationRules.Icon(XenAdmin.Alerts.AlertPriority.Unknown));
            Assert.IsFalse(NotificationRules.CanDismiss(new Message()), "no session, no dismiss");

            // FailedLoginAttemptAlert: totals, then one line per known source
            var text = NotificationRules.FailedLoginText("<body><total>5</total><unknown>2</unknown><known><number>3</number><ip>10.0.0.9</ip><username>root</username></known></body>");
            StringAssert.Contains(text, "5");
            StringAssert.Contains(text, "- ");
            StringAssert.Contains(text, "10.0.0.9");
            StringAssert.Contains(text, "root");
            Assert.IsNull(NotificationRules.FailedLoginText("not xml"));
        }
    }
}
