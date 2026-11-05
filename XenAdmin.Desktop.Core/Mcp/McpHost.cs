using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.ComponentModel;
using System.IO;
using System.Linq;
using System.Net;
using System.Reflection;
using System.Threading;
using System.Threading.Tasks;
using XenAdmin.Actions;
using XenAdmin.Network;

namespace XenAdmin.Desktop.Mcp
{
    /// <summary>
    /// <c>xcp-ng-center mcp [--allow-write]</c>: the headless bootstrap (the parts of App.OnFrameworkInitializationCompleted
    /// that XenModel needs, without a window) and the MCP server on stdin/stdout. Logs go to stderr only.
    /// </summary>
    public static class McpHost
    {
        public static int Run(string[] args, TextReader input, TextWriter output, TextWriter log)
        {
            var allowWrite = args.Contains("--allow-write");
            var env = new Func<string, string>(Environment.GetEnvironmentVariable);

            // XenModel raises its events through InvokeHelper: give it a thread of its own, as the UI thread does in the app
            var sync = new ThreadSynchronizer();
            InvokeHelper.Initialize(sync);
            sync.Invoke(new Action(XenAdmin.Model.Folders.InitFolders), null);

            // the app's TLS policy, without trust-on-first-use (McpTools explains it to the agent)
            var trust = new CertificateTrust(AppConfig.Load()) { RequirePin = true, VouchedFingerprints = CertificateTrust.PoolMemberFingerprints };
            var tools = new XcpTools(env, trust, sync);
#pragma warning disable SYSLIB0014 // XenAPI's JsonRpc falls back to this global callback, as on Windows
            ServicePointManager.ServerCertificateValidationCallback = trust.Validate;
#pragma warning restore SYSLIB0014
            XenAPI.HTTP.ServerCertificateValidator = trust.ValidateHost;

            var version = Assembly.GetEntryAssembly()?.GetName().Version?.ToString() ?? "0";
            var server = new McpServer("xcp-ng-center", version, tools.All(), allowWrite, tools.Redact);
            log.WriteLine($"xcp-ng-center MCP server {version}: {(allowWrite ? "read-write" : "read-only")}, servers: {string.Join(", ", tools.ConfiguredHosts)}");
            // in the background: an unreachable server must not delay the client's initialize
            _ = Task.Run(() => tools.AutoConnect(log));
            server.Run(input, output);
            tools.DisconnectAll();
            return 0;
        }
    }

    /// <summary>Runs delegates on one background thread (the "UI thread" of a headless XenModel).</summary>
    public sealed class ThreadSynchronizer : ISynchronizeInvoke
    {
        private readonly BlockingCollection<Action> _queue = new BlockingCollection<Action>();
        private readonly Thread _thread;

        public ThreadSynchronizer()
        {
            _thread = new Thread(() => { foreach (var a in _queue.GetConsumingEnumerable()) a(); }) { IsBackground = true, Name = "xenmodel-events" };
            _thread.Start();
        }

        public bool InvokeRequired => Thread.CurrentThread != _thread;

        public IAsyncResult BeginInvoke(Delegate method, object[] args)
        {
            var tcs = new TaskCompletionSource<object>();
            // DynamicInvoke wraps the delegate's exception: pass on the real one (tools report its message)
            _queue.Add(() => { try { tcs.SetResult(method.DynamicInvoke(args)); } catch (System.Reflection.TargetInvocationException e) when (e.InnerException != null) { tcs.SetException(e.InnerException); } catch (Exception e) { tcs.SetException(e); } });
            return tcs.Task;
        }

        public object EndInvoke(IAsyncResult result) => ((Task<object>)result).GetAwaiter().GetResult();
        public object Invoke(Delegate method, object[] args) => InvokeRequired ? EndInvoke(BeginInvoke(method, args)) : method.DynamicInvoke(args);
    }

    /// <summary>The write operations started by the tools, kept for operation_status.</summary>
    public sealed class Operations
    {
        private readonly ConcurrentDictionary<string, AsyncAction> _running = new ConcurrentDictionary<string, AsyncAction>();

        public string Start(AsyncAction action)
        {
            var id = Guid.NewGuid().ToString("N").Substring(0, 12);
            _running[id] = action;
            action.RunAsync();
            return id;
        }

        public AsyncAction Get(string id) => id != null && _running.TryGetValue(id, out var a) ? a : null;

        /// <summary>Waits up to 50 s (below the clients' usual 60 s request timeout).</summary>
        public static void Wait(AsyncAction action, int seconds)
        {
            var deadline = DateTime.UtcNow.AddSeconds(Math.Clamp(seconds, 0, 50));
            while (!action.IsCompleted && DateTime.UtcNow < deadline)
                Thread.Sleep(200);
        }

        public static object Status(string id, AsyncAction a) => new
        {
            operation_id = id,
            title = a.Title,
            state = !a.IsCompleted ? "running" : a.Succeeded ? "succeeded" : a.Cancelled ? "cancelled" : "failed",
            percent = a.PercentComplete,
            result = a.Succeeded ? Uuid(a) ?? a.Result : null,
            error = a.IsCompleted && !a.Succeeded ? a.Exception?.Message ?? a.Description : null
        };

        /// <summary>The uuid of the object an action created (a snapshot, an SR), when its result is a reference.</summary>
        private static string Uuid(AsyncAction a)
        {
            if (a.Result == null || !a.Result.StartsWith("OpaqueRef:") || a.Connection?.Cache == null)
                return null;
            var c = a.Connection;
            return c.Cache.VMs.FirstOrDefault(o => o.opaque_ref == a.Result)?.uuid ?? c.Cache.SRs.FirstOrDefault(o => o.opaque_ref == a.Result)?.uuid;
        }

        /// <summary>Runs a short action (scans, probes) to the end, failing with its error.</summary>
        public static void RunToEnd(AsyncAction action, int seconds = 120)
        {
            action.RunAsync();
            var deadline = DateTime.UtcNow.AddSeconds(seconds);
            while (!action.IsCompleted && DateTime.UtcNow < deadline)
                Thread.Sleep(200);
            if (!action.IsCompleted)
                throw new McpToolException(action.Title + " did not finish in " + seconds + " s");
        }
    }
}
