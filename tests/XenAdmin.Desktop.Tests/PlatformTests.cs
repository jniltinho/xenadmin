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
    }
}
