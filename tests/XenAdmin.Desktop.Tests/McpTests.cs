using System;
using System.IO;
using System.Linq;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Newtonsoft.Json.Linq;
using XenAdmin.Desktop.Mcp;

namespace XenAdmin.Desktop.Tests
{
    /// <summary>The MCP server's protocol and guards, without servers (the lab test drives the real process).</summary>
    [TestClass]
    public class McpTests
    {
        private static McpServer Server(bool allowWrite) => new McpServer("test", "1.0", new[]
        {
            new McpTool { Name = "read", Description = "r", InputSchema = new JObject { ["type"] = "object" }, Run = a => new { echo = (string)a["x"] } },
            new McpTool { Name = "fail", Description = "f", InputSchema = new JObject { ["type"] = "object" }, Run = _ => throw new McpToolException("no: s3cret") },
            new McpTool { Name = "write", Description = "w", Writes = true, InputSchema = new JObject { ["type"] = "object" }, Run = _ => "done" },
        }, allowWrite, s => s.Replace("s3cret", "***"));

        private static JObject Ask(McpServer s, string json) => JObject.Parse(s.Handle(json));

        [TestMethod]
        public void HandshakeNegotiatesTheVersion()
        {
            var s = Server(false);
            var old = Ask(s, "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"initialize\",\"params\":{\"protocolVersion\":\"2024-11-05\",\"capabilities\":{},\"clientInfo\":{\"name\":\"c\",\"version\":\"1\"}}}");
            Assert.AreEqual("2024-11-05", (string)old["result"]["protocolVersion"]);
            Assert.IsNotNull(old["result"]["capabilities"]["tools"]);
            var other = Ask(s, "{\"jsonrpc\":\"2.0\",\"id\":\"a\",\"method\":\"initialize\",\"params\":{\"protocolVersion\":\"2025-03-26\"}}");
            Assert.AreEqual("2025-06-18", (string)other["result"]["protocolVersion"], "unsupported versions get the latest");
            Assert.AreEqual("a", (string)other["id"], "string ids come back as they were");
            Assert.IsNull(s.Handle("{\"jsonrpc\":\"2.0\",\"method\":\"notifications/initialized\"}"), "notifications get no response");
            Assert.IsNull(s.Handle("{\"jsonrpc\":\"2.0\",\"method\":\"notifications/whatever\"}"));
            Assert.AreEqual(0, ((JObject)Ask(s, "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"ping\"}")["result"]).Count);
        }

        [TestMethod]
        public void ProtocolErrors()
        {
            var s = Server(false);
            Assert.AreEqual(-32700, (int)Ask(s, "{nope")["error"]["code"]);
            Assert.AreEqual(-32600, (int)Ask(s, "[{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"ping\"}]")["error"]["code"], "no batches");
            Assert.AreEqual(-32600, (int)Ask(s, "{\"id\":1,\"method\":\"ping\"}")["error"]["code"], "jsonrpc 2.0 is required");
            Assert.AreEqual(-32601, (int)Ask(s, "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"resources/list\"}")["error"]["code"]);
            Assert.AreEqual(-32602, (int)Ask(s, "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"nope\"}}")["error"]["code"]);
            Assert.AreEqual(-32600, (int)Ask(s, "{\"jsonrpc\":{},\"id\":1,\"method\":\"ping\"}")["error"]["code"], "a malformed envelope still gets an answer");
            Assert.AreEqual(-32600, (int)Ask(s, "{\"jsonrpc\":\"2.0\",\"id\":{\"a\":1},\"method\":\"ping\"}")["error"]["code"], "ids are strings or integers");
            Assert.AreEqual(-32600, (int)Ask(s, "{\"jsonrpc\":\"2.0\",\"id\":null,\"method\":\"ping\"}")["error"]["code"]);
        }

        [TestMethod]
        public void ReadOnlyByDefaultAndWritesNeedConfirm()
        {
            var ro = Server(false);
            var listed = Ask(ro, "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}")["result"]["tools"].Select(t => (string)t["name"]).ToList();
            CollectionAssert.AreEquivalent(new[] { "read", "fail" }, listed);
            Assert.AreEqual(-32602, (int)Ask(ro, "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"write\",\"arguments\":{\"confirm\":true}}}")["error"]["code"]);

            var rw = Server(true);
            var noConfirm = Ask(rw, "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"write\",\"arguments\":{\"confirm\":\"true\"}}}")["result"];
            Assert.IsTrue((bool)noConfirm["isError"], "confirm must be the boolean true");
            var ok = Ask(rw, "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"write\",\"arguments\":{\"confirm\":true}}}")["result"];
            Assert.IsFalse((bool)ok["isError"]);
            Assert.AreEqual("done", (string)ok["content"][0]["text"]);
        }

        [TestMethod]
        public void ToolResultsAndRedaction()
        {
            var s = Server(false);
            var r = Ask(s, "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"read\",\"arguments\":{\"x\":\"hi\"}}}")["result"];
            Assert.AreEqual("text", (string)r["content"][0]["type"]);
            Assert.AreEqual("hi", (string)JObject.Parse((string)r["content"][0]["text"])["echo"]);
            var f = Ask(s, "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"fail\"}}")["result"];
            Assert.IsTrue((bool)f["isError"]);
            Assert.AreEqual("no: ***", (string)f["content"][0]["text"], "secrets never go back to the agent");
        }

        [TestMethod]
        public void RunAnswersEveryRequestOnItsOwnLine()
        {
            var input = new StringReader("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"ping\"}\n\n{\"jsonrpc\":\"2.0\",\"method\":\"notifications/initialized\"}\n{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/list\"}\n");
            var output = new StringWriter();
            Server(false).Run(input, output);
            var lines = output.ToString().Split('\n', StringSplitOptions.RemoveEmptyEntries).Select(JObject.Parse).ToList();
            CollectionAssert.AreEquivalent(new[] { 1, 2 }, lines.Select(l => (int)l["id"]).ToList());
        }

        [TestMethod]
        public void ToolsOnlyReachTheOperatorsServers()
        {
            var env = new System.Collections.Generic.Dictionary<string, string> { ["XCP_HOST"] = "10.0.0.5", ["XCP_USER"] = "root", ["XCP_PASS"] = "pw-123", ["XCP_TRUST_FINGERPRINT"] = "aa:bb" };
            var trust = new CertificateTrust(AppConfig.Load(Path.Combine(Path.GetTempPath(), Guid.NewGuid() + "/config.json"))) { RequirePin = true };
            var tools = new XcpTools(k => env.TryGetValue(k, out var v) ? v : null, trust, new ThreadSynchronizer());
            Assert.AreEqual("AABB", trust.OperatorPins["10.0.0.5:443"], "a single fingerprint pins the single configured host");
            Assert.AreEqual(TrustDecision.Prompt, trust.Evaluate("10.0.0.9:443", "CC", System.Net.Security.SslPolicyErrors.None), "no trust-on-first-use, even for a valid chain");
            Assert.AreEqual(TrustDecision.Mismatch, trust.Evaluate("10.0.0.5:443", "CC", System.Net.Security.SslPolicyErrors.None));
            var server = new McpServer("t", "1", tools.All(), true, tools.Redact);
            var r = JObject.Parse(server.Handle("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"connect\",\"arguments\":{\"host\":\"evil.example\"}}}"))["result"];
            Assert.IsTrue((bool)r["isError"]);
            StringAssert.Contains((string)r["content"][0]["text"], "not one of the servers");
            Assert.AreEqual("x *** y", tools.Redact("x pw-123 y"));
            env["XCP_PASS"] = "a\"b\\c";
            var tricky = new McpServer("t", "1", new[] { new McpTool { Name = "r", InputSchema = new JObject(), Run = _ => new { error = "bad a\"b\\c", total = 1, operation_id = "1db84bb163e2", result = "c12420ea-5343-8d19-63fc-2326556c6a68" } } }, false, tools.Redact);
            var text = (string)JObject.Parse(tricky.Handle("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"r\"}}"))["result"]["content"][0]["text"];
            var parsed = JObject.Parse(text);
            Assert.AreEqual("bad ***", (string)parsed["error"], "secrets with quotes or backslashes are removed before serializing");
            Assert.AreEqual(1, (int)parsed["total"], "numbers stay intact");
            env["XCP_PASS"] = "1";
            parsed = JObject.Parse((string)JObject.Parse(tricky.Handle("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"r\"}}"))["result"]["content"][0]["text"]);
            Assert.AreEqual("1db84bb163e2", (string)parsed["operation_id"], "generated ids are not redacted");
            Assert.AreEqual("c12420ea-5343-8d19-63fc-2326556c6a68", (string)parsed["result"]);
            env["XCP_PASS"] = "x";
            Assert.AreEqual("***y", tools.Redact("xy"), "short secrets are removed too");
            env["XCP_PASS"] = "pw-123";
            var rw = new McpServer("t", "1", tools.All(), true, tools.Redact);
            var restart = JObject.Parse(rw.Handle("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"host_power\",\"arguments\":{\"host\":\"h\",\"operation\":\"restart\",\"confirm\":true}}}"));
            Assert.AreEqual(-32602, (int)restart["error"]["code"], "an unknown operation never becomes a shutdown");
            var noEnter = JObject.Parse(rw.Handle("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"host_maintenance\",\"arguments\":{\"host\":\"h\",\"confirm\":true}}}"));
            Assert.AreEqual(-32602, (int)noEnter["error"]["code"], "enter is required");
            var names = JObject.Parse(new McpServer("t", "1", tools.All(), false).Handle("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}"))["result"]["tools"].Select(t => (string)t["name"]).ToList();
            Assert.IsFalse(names.Contains("vm_power") || names.Contains("sr_create") || names.Contains("snapshot_delete"), "read-only by default");
            Assert.IsTrue(names.Contains("list_objects") && names.Contains("operation_status"));
        }
    }
}
