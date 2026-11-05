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
        public void OptionsCertificateWarningsProxyAndRoundTrip()
        {
            using var rsa = System.Security.Cryptography.RSA.Create(2048);
            var req = new System.Security.Cryptography.X509Certificates.CertificateRequest("CN=opt", rsa, System.Security.Cryptography.HashAlgorithmName.SHA256, System.Security.Cryptography.RSASignaturePadding.Pkcs1);
            using var cert = req.CreateSelfSigned(DateTimeOffset.Now.AddDays(-1), DateTimeOffset.Now.AddDays(1));
            var cfg = TempConfig();
            var trust = new CertificateTrust(cfg);
            Assert.IsFalse(trust.ValidateHost("new:443", cert, SslPolicyErrors.RemoteCertificateChainErrors), "new certificates warn by default");

            cfg.Options.WarnUnrecognizedCertificate = false;
            Assert.IsTrue(trust.ValidateHost("new:443", cert, SslPolicyErrors.RemoteCertificateChainErrors), "no warning: trusted and pinned");
            Assert.AreEqual(CertificateTrust.Fingerprint(cert), cfg.Fingerprints["new:443"]);

            trust.Pin("changed:443", "00");
            Assert.IsFalse(trust.ValidateHost("changed:443", cert, SslPolicyErrors.None), "a changed certificate warns by default");
            cfg.Options.WarnChangedCertificate = false;
            Assert.IsTrue(trust.ValidateHost("changed:443", cert, SslPolicyErrors.None), "no warning: the new certificate replaces the pin");
            Assert.AreEqual(CertificateTrust.Fingerprint(cert), cfg.Fingerprints["changed:443"]);

            trust.OperatorPins["op:443"] = "00";
            Assert.IsFalse(trust.ValidateHost("op:443", cert, SslPolicyErrors.None), "an operator pin is never overridden");
            var strict = new CertificateTrust(cfg) { RequirePin = true };
            Assert.IsFalse(strict.ValidateHost("strict:443", cert, SslPolicyErrors.RemoteCertificateChainErrors), "RequirePin (MCP) ignores the options");

            Assert.IsNull(AppOptions.ProxyError(XenAPI.HTTPHelper.ProxyStyle.DirectConnection, "", ""));
            Assert.IsNull(AppOptions.ProxyError(XenAPI.HTTPHelper.ProxyStyle.SpecifiedProxy, "proxy.lan", "3128"));
            Assert.IsNotNull(AppOptions.ProxyError(XenAPI.HTTPHelper.ProxyStyle.SpecifiedProxy, "proxy lan", "3128"));
            Assert.IsNotNull(AppOptions.ProxyError(XenAPI.HTTPHelper.ProxyStyle.SpecifiedProxy, "proxy.lan", "70000"));
            var o = new AppOptions { ProxyStyle = XenAPI.HTTPHelper.ProxyStyle.SpecifiedProxy, ProxyAddress = "proxy.lan", ProxyPort = 3128, BypassProxyForServers = true };
            Assert.IsNull(o.Proxy(true), "bypassed for servers");
            Assert.AreEqual(new Uri("http://proxy.lan:3128/"), ((System.Net.WebProxy)o.Proxy(false)).Address);
            Assert.IsNull(new AppOptions().Proxy(true), "direct by default");

            cfg.Options.ConnectionTimeoutSeconds = 45;
            cfg.Options.FillAreaUnderGraphs = true;
            cfg.Save();
            var again = AppConfig.Load(cfg.Path);
            Assert.AreEqual(45, again.Options.ConnectionTimeoutSeconds);
            Assert.IsTrue(again.Options.FillAreaUnderGraphs);
            Assert.IsFalse(again.Options.WarnChangedCertificate);
            Assert.AreEqual(45000, new DesktopConfigProvider { Options = again.Options }.ConnectionTimeout);
        }

        /// <summary>An appliance like the VHD import path builds: one VM, an 8 MB dynamic VHD (disk1.vhd) and its envelope.</summary>
        internal static string WriteSampleAppliance(string dir, string name)
        {
            Directory.CreateDirectory(dir);
            var vhd = Path.Combine(dir, "disk1.vhd");
            const long capacity = 8L << 20;
            using (var stream = File.Create(vhd))
            using (DiscUtils.Vhd.Disk.InitializeDynamic(stream, DiscUtils.Streams.Ownership.None, capacity)) { }
            var env = XenOvf.OVF.CreateOvfEnvelope(Guid.NewGuid().ToString(), name, 1, 256UL << 20, "", "",
                (ulong)capacity, "disk1.vhd", (ulong)new FileInfo(vhd).Length, XenAdmin.Core.BrandManager.ProductBrand);
            var ovf = Path.Combine(dir, name + ".ovf");
            XenOvf.OVF.SaveAs(env, ovf);
            return ovf;
        }

        [TestMethod]
        public void ArchiveEntriesCannotEscapeTheExtractionDirectory()
        {
            var dir = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString());
            var absolute = Path.Combine(dir, "absolute-escape.txt");
            // a relative climb and an absolute path (Path.Combine keeps the absolute one)
            foreach (var evilName in new[] { "../escape.txt", absolute })
            {
                var target = Path.Combine(dir, "extract");
                Directory.CreateDirectory(target);
                var tar = Path.Combine(dir, "evil.ova");
                using (var writer = new System.Formats.Tar.TarWriter(File.Create(tar), System.Formats.Tar.TarEntryFormat.Ustar, false))
                {
                    writer.WriteEntry(new System.Formats.Tar.UstarTarEntry(System.Formats.Tar.TarEntryType.RegularFile, "inside.txt") { DataStream = new MemoryStream(new byte[] { 1 }) });
                    writer.WriteEntry(new System.Formats.Tar.UstarTarEntry(System.Formats.Tar.TarEntryType.RegularFile, evilName) { DataStream = new MemoryStream(new byte[] { 2 }) });
                }
                using (var stream = File.OpenRead(tar))
                using (var iterator = XenCenterLib.Archive.ArchiveFactory.Reader(XenCenterLib.Archive.ArchiveFactory.Type.Tar, stream))
                    Assert.ThrowsException<IOException>(() => iterator.ExtractAllContents(target), evilName);
                Assert.IsTrue(File.Exists(Path.Combine(target, "inside.txt")), "entries inside are extracted");
                Assert.IsFalse(File.Exists(Path.Combine(dir, "escape.txt")) || File.Exists(absolute), "nothing is written outside: " + evilName);
            }
            Directory.Delete(dir, true);
        }

        [TestMethod]
        public void OvfApplianceOpensWithoutServer()
        {
            // XCP_SAMPLE_OVF_DIR keeps the appliance for a manual check of the import dialog
            var keep = Environment.GetEnvironmentVariable("XCP_SAMPLE_OVF_DIR");
            var dir = string.IsNullOrEmpty(keep) ? Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString()) : keep;
            try
            {
                var ovf = WriteSampleAppliance(dir, "sample-appliance");
                Assert.IsTrue(OvfRules.IsAppliance(ovf));
                Assert.IsFalse(OvfRules.IsAppliance("vm.xva"));
                var (package, error, _) = OvfRules.Open(ovf);
                Assert.IsNull(error, error);
                CollectionAssert.AreEqual(new[] { "sample-appliance" }, OvfRules.Vms(package).Select(v => v.Name).ToArray());
                Assert.AreEqual((1, 1), OvfRules.Counts(package), "one disk and the interface CreateOvfEnvelope adds");
                Assert.IsNotNull(OvfRules.Open(Path.Combine(dir, "missing.ova")).Error);
            }
            finally
            {
                if (string.IsNullOrEmpty(keep)) Directory.Delete(dir, true);
            }
        }

        [TestMethod]
        public void SearchEditorBuildExportImport()
        {
            var types = XenAdmin.XenSearch.ObjectTypes.VM | XenAdmin.XenSearch.ObjectTypes.Server;
            Assert.IsNotNull(SearchRules.Validate(" ", types));
            Assert.IsNotNull(SearchRules.Validate("x", XenAdmin.XenSearch.ObjectTypes.None));
            var groupings = SearchRules.Groupings();
            var os = groupings.FindIndex(g => g.Name == "Operating System");
            var search = SearchRules.Build("My VMs", types, "web", groupings[os].Make(), null);
            Assert.AreEqual("My VMs", search.Name);
            Assert.AreEqual(types, search.Query.QueryScope.ObjectTypes);
            Assert.IsNotNull(search.Query.QueryFilter, "the text filter");
            Assert.AreEqual(os, SearchRules.GroupingIndex(search), "the editor shows the grouping again");

            var edited = SearchRules.Build("My VMs 2", types, "", null, search);
            Assert.AreSame(search.Query.QueryFilter, edited.Query.QueryFilter, "editing keeps the conditions");
            Assert.AreEqual(0, SearchRules.GroupingIndex(edited));

            var file = Path.Combine(Path.GetTempPath(), Guid.NewGuid() + ".xensearch");
            try
            {
                SearchRules.Export(search, file);
                var back = SearchRules.Import(file);
                Assert.AreEqual(1, back.Count);
                Assert.AreEqual("My VMs", back[0].Name);
                Assert.AreEqual(types, back[0].Query.QueryScope.ObjectTypes);
                File.WriteAllText(file, "not a search");
                Assert.IsTrue(SearchRules.Import(file) == null || SearchRules.Import(file).Count == 0, "a file that is not a search");
            }
            finally
            {
                File.Delete(file);
            }
        }

        [TestMethod]
        public void ConfigRoundTripHasNoPassword()
        {
            var cfg = TempConfig();
            cfg.AddOrUpdateServer("xcp1", 443, "root");
            cfg.Fingerprints["xcp1:443"] = "AB";
            cfg.Save();
            var json = File.ReadAllText(cfg.Path);
            StringAssert.DoesNotMatch(json, new System.Text.RegularExpressions.Regex("(?i)passw"));
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
            Assert.IsFalse(Details.NotPorted(new Host(), "Console"), "the server console is its control domain's");
            Assert.IsFalse(Details.NotPorted(new VM(), "Snapshots"));
            Assert.IsFalse(Details.NotPorted(new VM(), "Search"));
            Assert.IsTrue(Details.NotPorted(new SR(), "Search"));
            Assert.IsFalse(Details.NotPorted(new VM(), "Performance"));
            Assert.IsFalse(Details.NotPorted(new Host(), "Performance"));
            Assert.IsFalse(Details.NotPorted(new Host(), "NICs"));
            Assert.IsFalse(Details.NotPorted(new Host(), "GPU"));
            Assert.IsFalse(Details.NotPorted(new Host(), "USB"));
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
            Assert.IsTrue(VmOperations.NeedsConfirmation(VmOp.Shutdown), "ShutDownVMCommand always confirms");
            Assert.IsFalse(VmOperations.NeedsConfirmation(VmOp.Start));
        }

        [TestMethod]
        public void HaBlocksStartOfANonAgileProtectedVm()
        {
            Assert.IsNotNull(VmOperations.HaStartBlocker(true, true, new[] { true, false }, "vm"), "a disk on local storage");
            Assert.IsNull(VmOperations.HaStartBlocker(true, true, new[] { true }, "vm"), "all disks shared: the server decides");
            Assert.IsNull(VmOperations.HaStartBlocker(false, true, new[] { false }, "vm"), "HA off");
            Assert.IsNull(VmOperations.HaStartBlocker(true, false, new[] { false }, "vm"), "best effort / do not restart");
            Assert.IsNull(VmOperations.HaStartBlocker(null), "nothing selected");
        }

        [TestMethod]
        public void FailuresShowReadableMessages()
        {
            // the tracked actions report Exception.Message: XAPI codes become the Windows app's texts
            var noMemory = new Failure(Failure.HOST_NOT_ENOUGH_FREE_MEMORY, "1073741824", "536870912");
            Assert.IsFalse(noMemory.Message.Contains(Failure.HOST_NOT_ENOUGH_FREE_MEMORY), noMemory.Message);
            Assert.IsFalse(string.IsNullOrWhiteSpace(noMemory.Message));
            var rbac = new Failure(Failure.RBAC_PERMISSION_DENIED, "VM.start", "No permission");
            Assert.IsFalse(string.IsNullOrWhiteSpace(rbac.Message));
            Assert.IsFalse(rbac.Message.StartsWith(Failure.RBAC_PERMISSION_DENIED), rbac.Message);
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
        public void StorageDeviceConfigs()
        {
            // CIFS_ISO: a sub-folder of the share goes to iso_path; credentials only with "Use different user name"
            var smb = SrRules.IsoDeviceConfig(SrRules.Kind.SmbIso, @"\\nas\isos\linux\alpine", false, null, null);
            Assert.AreEqual("//nas/isos", smb["location"]);
            Assert.AreEqual("/linux/alpine", smb["iso_path"]);
            Assert.IsFalse(smb.ContainsKey("username"));
            Assert.AreEqual("u", SrRules.IsoDeviceConfig(SrRules.Kind.SmbIso, @"\\nas\isos", false, "u", "p")["username"]);
            var nfs = SrRules.IsoDeviceConfig(SrRules.Kind.NfsIso, " nas:/isos ", true, null, null);
            Assert.AreEqual(("nas:/isos", "nfs_iso", "4"), (nfs["location"], nfs["type"], nfs["nfsversion"]));
            Assert.IsTrue(SrRules.ValidLocation(SrRules.Kind.NfsIso, "nas:/isos") && !SrRules.ValidLocation(SrRules.Kind.SmbIso, "nas:/isos"));
            Assert.AreEqual(SR.Content_Type_ISO, SrRules.ContentType(SrRules.Kind.NfsIso));
            Assert.AreEqual("", SrRules.ContentType(SrRules.Kind.Iscsi));

            // LVMoISCSI: the IQN's own address, the LUN by SCSI id, or by serial when the target gives no LUN id
            var iqn = new XenAdmin.Actions.IScsiIqnInfo(0, "iqn.x:sr", "10.0.0.20", 3260);
            var lun = new XenAdmin.Actions.ISCSIInfo { LunID = 1, ScsiID = "3600", Serial = "beaf11", Vendor = "IET", Size = 4L << 30 };
            var d = SrRules.IscsiDeviceConfig(iqn, lun, "chap", "secret");
            Assert.AreEqual(("10.0.0.20", "3260", "iqn.x:sr", "3600", "chap"), (d["target"], d["port"], d["targetIQN"], d["SCSIid"], d["chapuser"]));
            lun.LunID = -1;
            Assert.AreEqual("beaf11", SrRules.IscsiDeviceConfig(iqn, lun, null, null)["LUNSerial"]);
            StringAssert.StartsWith(SrRules.LunLabel(new XenAdmin.Actions.ISCSIInfo { LunID = 1, Serial = "beaf11", Vendor = "IET", Size = 4L << 30 }), "LUN 1: beaf11: 4 GB (IET)");
        }

        [TestMethod]
        public void VmNetworkInterfaceRules()
        {
            var vm = new VM { opaque_ref = "OpaqueRef:vm" };
            var plugged = new VIF { currently_attached = true, allowed_operations = new List<vif_operations> { vif_operations.unplug } };
            var stuck = new VIF { currently_attached = true, allowed_operations = new List<vif_operations>() };
            var detached = new VIF { currently_attached = false, allowed_operations = new List<vif_operations> { vif_operations.plug } };
            Assert.IsTrue(NetworkRules.CanEdit(vm, plugged) && NetworkRules.CanToggle(vm, plugged));
            Assert.IsFalse(NetworkRules.CanEdit(vm, stuck) || NetworkRules.CanRemove(vm, stuck) || NetworkRules.CanToggle(vm, stuck));
            Assert.IsTrue(NetworkRules.CanRemove(vm, detached) && NetworkRules.CanToggle(vm, detached));
            Assert.AreEqual(Messages.VM_NETWORK_TAB_ACTIVATE_BUTTON_LABEL, NetworkRules.ToggleLabel(detached));
            Assert.AreEqual(Messages.VM_NETWORK_TAB_DEACTIVATE_BUTTON_LABEL, NetworkRules.ToggleLabel(plugged));
            Assert.IsFalse(NetworkRules.CanEdit(vm, null) || NetworkRules.CanToggle(vm, null));

            // VIFDialog.ValidateInput: a network is required; MAC and QoS are checked only when chosen
            var net = new XenAPI.Network { opaque_ref = "OpaqueRef:net" };
            Assert.AreEqual(Messages.SELECT_NETWORK_TOOLTIP, NetworkRules.Validate(null, false, null, false, null));
            Assert.IsNull(NetworkRules.Validate(net, false, "junk", false, "junk"));
            Assert.AreEqual("", NetworkRules.Validate(net, true, " ", false, null));
            Assert.AreEqual(Messages.MAC_INVALID, NetworkRules.Validate(net, true, "aa:bb", false, null));
            Assert.IsNull(NetworkRules.Validate(net, true, "aa:bb:cc:dd:ee:ff", false, null));
            Assert.AreEqual(Messages.ENTER_VALID_QOS, NetworkRules.Validate(net, false, null, true, "0"));
            Assert.IsNull(NetworkRules.Validate(net, false, null, true, "100"));

            // GetNewSettings / ChangesHaveBeenMade: turning QoS off clears the algorithm; an unchanged edit is a cancel
            var existing = new VIF
            {
                opaque_ref = "OpaqueRef:vif", device = "0", MAC = "aa:bb:cc:dd:ee:ff", network = new XenRef<XenAPI.Network>("OpaqueRef:net"),
                qos_algorithm_type = VIF.RATE_LIMIT_QOS_VALUE, qos_algorithm_params = new Dictionary<string, string> { { VIF.KBPS_QOS_PARAMS_KEY, "100" } }
            };
            Assert.IsFalse(NetworkRules.Changed(existing, net, "aa:bb:cc:dd:ee:ff", 0, true, "100"));
            Assert.IsTrue(NetworkRules.Changed(existing, net, "aa:bb:cc:dd:ee:ff", 0, true, "200"));
            Assert.IsTrue(NetworkRules.Changed(existing, net, "aa:bb:cc:dd:ee:ff", 0, false, "100"));
            Assert.IsTrue(NetworkRules.Changed(null, net, "", 1, false, ""));
            var off = NetworkRules.Settings(vm, existing, net, "", 0, false, "");
            Assert.AreEqual(("", "", "0", "OpaqueRef:vm"), (off.qos_algorithm_type, off.MAC, off.device, off.VM.opaque_ref));
            var added = NetworkRules.Settings(vm, null, net, "", 2, true, "500");
            Assert.AreEqual((VIF.RATE_LIMIT_QOS_VALUE, "500", "2"), (added.qos_algorithm_type, added.qos_algorithm_params[VIF.KBPS_QOS_PARAMS_KEY], added.device));
            Assert.AreEqual("", NetworkRules.Limit(new VIF { qos_algorithm_type = "" }));
            Assert.AreEqual("100", NetworkRules.Limit(existing));

            // selection while UpdateVIFAction recreates device 1: it stays on 1, even after the action completed and
            // before the replacement reached the cache; with no edit pending it falls back to the first device
            Assert.AreEqual("1", NetworkRules.Selection(new[] { "0" }, "1", "1"));
            Assert.AreEqual("1", NetworkRules.Selection(new[] { "0", "1" }, "1", null));
            Assert.AreEqual("0", NetworkRules.Selection(new[] { "0" }, "1", null));
            Assert.IsNull(NetworkRules.Selection(new string[0], null, null));

            // the edit guard across refreshes: original present -> absent -> replacement present, device 1 selected throughout
            var before = new[] { ("0", "r0"), ("1", "old") };
            var gap = new[] { ("0", "r0") };
            var after = new[] { ("0", "r0"), ("1", "new") };
            Assert.IsFalse(NetworkRules.EditReplaced(before, "1", "old"));
            Assert.AreEqual("1", NetworkRules.Selection(new[] { "0", "1" }, "1", "1"));
            Assert.IsFalse(NetworkRules.EditReplaced(gap, "1", "old"));
            Assert.AreEqual("1", NetworkRules.Selection(new[] { "0" }, "1", "1"));
            Assert.IsTrue(NetworkRules.EditReplaced(after, "1", "old"));
            Assert.AreEqual("1", NetworkRules.Selection(new[] { "0", "1" }, "1", null));
        }

        [TestMethod]
        public void HostNetworkRules()
        {
            // NetWDetails: VLAN range, in use, VLAN 0 information; MTU range shown always
            Assert.AreEqual((true, Messages.INVALID_NUMBER), NetworkRules.ValidateVlan("x", 1, new int[0]));
            Assert.AreEqual((true, string.Format(Messages.NETW_DETAILS_VLAN_RANGE, 1, 4094)), NetworkRules.ValidateVlan("0", 1, new int[0]));
            Assert.AreEqual((true, Messages.NETW_DETAILS_VLAN_NUMBER_IN_USE), NetworkRules.ValidateVlan("5", 1, new[] { -1, 5 }));
            Assert.AreEqual((false, Messages.NETW_VLAN_ZERO), NetworkRules.ValidateVlan("0", 0, new[] { -1 }));
            Assert.AreEqual((false, (string)null), NetworkRules.ValidateVlan("12", 0, new[] { -1 }));
            Assert.AreEqual(0, NetworkRules.VlanFor(0, new[] { -1 }), "a free VLAN is kept");
            Assert.AreEqual(3, NetworkRules.VlanFor(1, new[] { -1, 1, 2 }), "a used VLAN moves to the first free one");
            var range = string.Format(Messages.ALLOWED_MTU_RANGE, 1500, 9000);
            Assert.AreEqual((false, range), NetworkRules.ValidateMtu("9000", 9000));
            Assert.AreEqual((true, range), NetworkRules.ValidateMtu("9216", 9000));
            Assert.AreEqual((true, string.Format(Messages.ALLOWED_MTU_VALUE, 1500)), NetworkRules.ValidateMtu("1400", 1500));
            Assert.AreEqual(9216, NetworkRules.MaxMtu(new PIF { MTU = 65000 }));

            // CreateNetworkObject
            var n = NetworkRules.NewNetwork(" lab ", "d", false, 9000);
            Assert.AreEqual(("lab", "d", true, 9000L, false), (n.name_label, n.name_description, n.managed, n.MTU, n.GetAutoPlug()));
        }

        [TestMethod]
        public void MemorySettingsKeepExactValues()
        {
            // ChangeMemorySettings: a static maximum unchanged to the MB keeps the server's bytes (no resize by rounding)
            var vm = new VM { memory_static_max = 536870912L + 4096, memory_dynamic_min = 268435456L, memory_dynamic_max = 536870912L + 4096 };
            var same = MemoryRules.Settings(vm, 268435456L, 536870912L, 536870912L);
            Assert.AreEqual((268435456L, 536870912L + 4096, 536870912L + 4096), same);
            var bigger = MemoryRules.Settings(vm, 268435456L, 1073741824L, 1073741824L);
            Assert.AreEqual((268435456L, 1073741824L, 1073741824L), bigger);
            Assert.IsTrue(MemoryRules.InMegabytes(vm));
        }

        [TestMethod]
        public void TemplateIconsExist()
        {
            // every template type maps to a PNG the app ships (XenAdmin/Images, linked as Avalonia resources)
            var dir = AppContext.BaseDirectory;
            while (dir != null && !Directory.Exists(Path.Combine(dir, "XenAdmin", "Images")))
                dir = Path.GetDirectoryName(dir);
            Assert.IsNotNull(dir, "repository root");
            foreach (var type in Enum.GetValues<VM.VmTemplateType>())
            {
                var icon = NewVmDefaults.TemplateIcon(type);
                var file = icon == "Logo" ? Path.Combine(dir, "XenAdmin.Desktop", "packaging", "xcp-ng-center.png") : Path.Combine(dir, "XenAdmin", "Images", icon + ".png");
                Assert.IsTrue(File.Exists(file), type + " icon");
            }
            Assert.AreEqual("windows_h32bit_16", NewVmDefaults.TemplateIcon(VM.VmTemplateType.WindowsServer));
            Assert.AreEqual("ubuntu_16x", NewVmDefaults.TemplateIcon(VM.VmTemplateType.Ubuntu));
        }

        [TestMethod]
        public void BondRulesWithoutServer()
        {
            Assert.IsFalse(BondRules.ValidSelection(1, 4), "a bond needs two NICs");
            Assert.IsTrue(BondRules.ValidSelection(2, 4));
            Assert.IsFalse(BondRules.ValidSelection(5, 4), "no more than the pool's limit");
            Assert.AreEqual(9000, BondRules.MaxMtu(new[] { new PIF { MTU = 9000 }, new PIF { MTU = 9216 } }), "the smallest member MTU");
            Assert.AreEqual(9216, BondRules.MaxMtu(new PIF[0]));
        }

        [TestMethod]
        public void AlertDateFilterAndCsv()
        {
            var now = new DateTime(2026, 11, 5, 14, 30, 0);
            Assert.IsNull(NotificationRules.DateFrom(NotificationRules.DateFilter.All, now));
            Assert.AreEqual(new DateTime(2026, 11, 4, 14, 30, 0), NotificationRules.DateFrom(NotificationRules.DateFilter.Last24Hours, now));
            Assert.AreEqual(new DateTime(2026, 11, 2), NotificationRules.DateFrom(NotificationRules.DateFilter.Last3Days, now), "from the beginning of the day");
            Assert.AreEqual(new DateTime(2026, 10, 6), NotificationRules.DateFrom(NotificationRules.DateFilter.Last30Days, now));

            var row = new NotificationRules.AlertRow(null, XenAdmin.Alerts.AlertPriority.Priority2, "Disk \"sda\" full", "a, b", "host1", now);
            var line = NotificationRules.CsvLine(row);
            StringAssert.StartsWith(line, "\"Disk \"\"sda\"\" full\",");
            StringAssert.Contains(line, "\"a, b\",\"host1\",");
            Assert.AreEqual(5, NotificationRules.CsvHeader().Split(',').Length);

            // DateFilterDialog: start at the beginning and end at the end of the day; unset bounds never hide
            Assert.AreEqual(new DateTime(2026, 11, 2), NotificationRules.DateBound(2, true, now, now));
            Assert.AreEqual(new DateTime(2026, 11, 2, 23, 59, 59).AddTicks(9999999), NotificationRules.DateBound(2, false, now, now));
            Assert.IsNull(NotificationRules.DateBound(0, true, now, now));
            var custom = new DateTime(2026, 1, 2, 3, 4, 5);
            Assert.AreEqual(custom, NotificationRules.DateBound(NotificationRules.CustomBound, true, now, custom));
            Assert.IsFalse(NotificationRules.ValidRange(now, now.AddDays(-1)), "the end must be after the start");
            Assert.IsTrue(NotificationRules.ValidRange(null, now));
            Assert.IsTrue(NotificationRules.HideByDate(now, now.AddHours(1), null));
            Assert.IsTrue(NotificationRules.HideByDate(now, null, now.AddHours(-1)));
            Assert.IsFalse(NotificationRules.HideByDate(now, null, null));
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
