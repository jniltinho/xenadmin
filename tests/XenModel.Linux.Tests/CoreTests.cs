using System;
using System.Linq;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Newtonsoft.Json;
using XenAdmin.Actions.Updates;
using XenAPI;

namespace XenModel.Linux.Tests
{
    [TestClass]
    public class CoreTests
    {
        [TestMethod]
        public void EmbeddedResourcesLoad()
        {
            Assert.IsFalse(string.IsNullOrEmpty(XenAdmin.Messages.ResourceManager.GetString("OK")));
        }

        [TestMethod]
        public void FileServiceTokenDeserializes()
        {
            var t = JsonConvert.DeserializeObject<FileServiceToken>("{\"token\":\"abc\",\"token_type\":\"Bearer\"}");
            Assert.AreEqual("abc", t.token);
            Assert.AreEqual("Bearer", t.token_type);
        }
    }

    /// <summary>Runs only when XCP_HOST/XCP_USER/XCP_PASS are set (see .env.example).</summary>
    [TestClass]
    public class XcpIntegrationTests
    {
        [TestMethod]
        [TestCategory("Integration")]
        public void LoginAndListVms()
        {
            var host = Environment.GetEnvironmentVariable("XCP_HOST");
            if (string.IsNullOrEmpty(host))
                Assert.Inconclusive("XCP_HOST not set");

            // lab host has a self-signed certificate
            var session = new Session($"https://{host}") { ServerCertificateValidationCallback = (_, _, _, _) => true };
            session.login_with_password(Environment.GetEnvironmentVariable("XCP_USER"), Environment.GetEnvironmentVariable("XCP_PASS"));
            try
            {
                var hosts = Host.get_all_records(session);
                var vms = VM.get_all_records(session).Values.Where(v => !v.is_a_template && !v.is_control_domain).ToList();
                Assert.IsTrue(hosts.Count > 0);
                System.Console.WriteLine($"API {session.APIVersion}; hosts: {string.Join(", ", hosts.Values.Select(h => h.name_label))}");
                foreach (var vm in vms)
                    System.Console.WriteLine($"VM {vm.name_label} [{vm.power_state}]");
            }
            finally
            {
                session.logout();
            }
        }
    }
}
