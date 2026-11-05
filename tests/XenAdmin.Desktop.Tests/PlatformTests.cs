using System.Collections.Generic;
using System.IO;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using XenAdmin.Desktop.Platform;

namespace XenAdmin.Desktop.Tests
{
    [TestClass]
    public class PlatformTests
    {
        [AssemblyInitialize]
        public static void UseHostPlatform(TestContext _) => PlatformServices.Current = DesktopPlatform.Detect();

        private static System.Func<string, string> Env(Dictionary<string, string> vars) =>
            k => vars.TryGetValue(k, out var v) ? v : null;

        [TestMethod]
        public void UnixUsesXdgConfigHomeThenDotConfig()
        {
            var withXdg = new UnixPlatform("Linux", Env(new() { ["XDG_CONFIG_HOME"] = "/x/cfg" }), "/home/u");
            Assert.AreEqual(Path.Combine("/x/cfg", "xcp-ng-center"), withXdg.ConfigDirectory);
            var without = new UnixPlatform("FreeBSD", Env(new()), "/home/u");
            Assert.AreEqual(Path.Combine("/home/u", ".config", "xcp-ng-center"), without.ConfigDirectory);
            Assert.AreEqual("FreeBSD", without.Name);
        }

        [TestMethod]
        public void MacUsesApplicationSupport()
        {
            var mac = new MacPlatform(Env(new()), "/Users/u");
            Assert.AreEqual(Path.Combine("/Users/u", "Library", "Application Support", "xcp-ng-center"), mac.ConfigDirectory);
        }

        [TestMethod]
        public void WindowsUsesAppDataSeparateFromWinForms()
        {
            var win = new WindowsPlatform(Env(new() { ["APPDATA"] = "C:/Users/u/AppData/Roaming" }), "C:/Users/u");
            Assert.AreEqual(Path.Combine("C:/Users/u/AppData/Roaming", "xcp-ng-center"), win.ConfigDirectory);
            var fallback = new WindowsPlatform(Env(new()), "C:/Users/u");
            Assert.AreEqual(Path.Combine("C:/Users/u", "AppData", "Roaming", "xcp-ng-center"), fallback.ConfigDirectory);
        }

        [TestMethod]
        public void AppConfigFollowsCurrentPlatformAndNoSecretStoreStoresNothing()
        {
            var saved = PlatformServices.Current;
            try
            {
                PlatformServices.Current = new UnixPlatform("Linux", Env(new() { ["XDG_CONFIG_HOME"] = "/t" }), "/h");
                Assert.AreEqual(Path.Combine("/t", "xcp-ng-center"), AppConfig.DefaultDirectory());
            }
            finally { PlatformServices.Current = saved; }

            var none = new NoSecretStore();
            Assert.IsFalse(none.IsAvailable);
            Assert.IsNull(none.Load("k"));
        }

        [TestMethod]
        public void UnixLaunchersPickTheTerminalAndRdpClient()
        {
            System.Func<string, string> Have(params string[] names) => n => System.Array.IndexOf(names, n) >= 0 ? "/usr/bin/" + n : null;
            // $TERMINAL wins, then the Debian alternative, then known terminals; each with its "run this" flag
            var c = UnixPlatform.SshCommand(Have("ssh", "konsole", "xterm"), null, "root@10.0.0.5").Value;
            Assert.AreEqual("/usr/bin/konsole", c.Exe);
            CollectionAssert.AreEqual(new[] { "-e", "/usr/bin/ssh", "--", "root@10.0.0.5" }, c.Args);
            c = UnixPlatform.SshCommand(Have("ssh", "gnome-terminal", "myterm"), "myterm", "root@h").Value;
            Assert.AreEqual("/usr/bin/myterm", c.Exe);
            c = UnixPlatform.SshCommand(Have("ssh", "gnome-terminal"), null, "root@h").Value;
            CollectionAssert.AreEqual(new[] { "--", "/usr/bin/ssh", "--", "root@h" }, c.Args);
            c = UnixPlatform.SshCommand(Have("ssh", "gnome-terminal", "xterm"), "gnome-terminal", "root@h").Value;
            CollectionAssert.AreEqual(new[] { "--", "/usr/bin/ssh", "--", "root@h" }, c.Args, "$TERMINAL keeps its own run arguments");
            Assert.IsNull(UnixPlatform.SshCommand(Have("xterm"), null, "root@h"), "no ssh client");

            // FreeRDP runs in a terminal (its certificate question needs one); Remmina asks in its own window
            var r = UnixPlatform.RdpCommand(Have("remmina", "xfreerdp", "xterm"), null, "10.0.0.7").Value;
            Assert.AreEqual("/usr/bin/xterm", r.Exe);
            CollectionAssert.AreEqual(new[] { "-e", "/usr/bin/xfreerdp", "/v:10.0.0.7", "/dynamic-resolution", "+clipboard" }, r.Args);
            r = UnixPlatform.RdpCommand(Have("remmina", "xfreerdp"), null, "10.0.0.7").Value;
            CollectionAssert.AreEqual(new[] { "-c", "rdp://10.0.0.7" }, r.Args, "no terminal: Remmina");
            Assert.IsNull(UnixPlatform.RdpCommand(Have("xfreerdp"), null, "h"), "FreeRDP without a terminal cannot ask about the certificate");
        }

        [TestMethod]
        public void SshTargetAcceptsOnlyLiteralAddressesAndPlainUsers()
        {
            Assert.AreEqual("root@10.77.58.60", ConsoleRules.SshTarget("root", "10.77.58.60"));
            Assert.AreEqual("DOM\\ann@fe80::1", ConsoleRules.SshTarget("DOM\\ann", "fe80::1"));
            Assert.IsNull(ConsoleRules.SshTarget("root", "$(id)"), "guest data is never a command");
            Assert.IsNull(ConsoleRules.SshTarget("root", "127.0.0.1 -oProxyCommand=calc.exe"));
            Assert.IsNull(ConsoleRules.SshTarget("-oProxyCommand=x", "10.0.0.1"));
            Assert.IsNull(ConsoleRules.SshTarget("ann'; rm", "10.0.0.1"));
        }
    }
}
