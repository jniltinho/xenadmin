using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Media.Imaging;
using Avalonia.Platform;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Linux
{
    /// <summary>The WinForms app's own PNGs (XenAdmin/Images), embedded as Avalonia resources.</summary>
    public static class Images
    {
        private static readonly Dictionary<string, Bitmap> Cache = new Dictionary<string, Bitmap>();

        public static Bitmap Get(string file)
        {
            lock (Cache)
            {
                if (!Cache.TryGetValue(file, out var bmp))
                {
                    var uri = file == "Logo" ? "avares://xcp-ng-center/Assets/xcp-ng-center.png" : $"avares://xcp-ng-center/Images/{file}.png";
                    try { bmp = new Bitmap(AssetLoader.Open(new Uri(uri))); }
                    catch { bmp = null; }
                    Cache[file] = bmp;
                }
                return bmp;
            }
        }

        private static readonly Dictionary<string, Bitmap> GrayCache = new Dictionary<string, Bitmap>();

        /// <summary>Greyscale copy, as ToolStrip renders disabled buttons.</summary>
        public static Bitmap Gray(string file)
        {
            lock (GrayCache)
            {
                if (GrayCache.TryGetValue(file, out var gray))
                    return gray;
                var src = Get(file);
                if (src == null)
                    return null;
                var size = src.PixelSize;
                var wb = new WriteableBitmap(size, src.Dpi, Avalonia.Platform.PixelFormat.Bgra8888, AlphaFormat.Unpremul);
                using (var fb = wb.Lock())
                {
                    src.CopyPixels(new PixelRect(size), fb.Address, fb.RowBytes * size.Height, fb.RowBytes);
                    var px = new byte[fb.RowBytes * size.Height];
                    System.Runtime.InteropServices.Marshal.Copy(fb.Address, px, 0, px.Length);
                    for (int i = 0; i < px.Length; i += 4)
                    {
                        // lighter grey, like ControlPaint.DrawImageDisabled
                        var l = (byte)Math.Min(255, 70 + (px[i] * 11 + px[i + 1] * 59 + px[i + 2] * 30) / 100 * 2 / 3);
                        px[i] = px[i + 1] = px[i + 2] = l;
                    }
                    System.Runtime.InteropServices.Marshal.Copy(px, 0, fb.Address, px.Length);
                }
                return GrayCache[file] = wb;
            }
        }

        public const string Home = "infra_view_16_textured";

        // shortcut: port of XenAdmin/Images.cs GetIconFor (connection, pool, host, VM, SR only),
        // kept in step by hand; move it into XenModel when the WinForms app can be built on Linux CI.
        public static string IconFor(IXenConnection connection)
        {
            var pool = Helpers.GetPool(connection);
            if (pool != null)
                return IconFor(pool);
            var host = Helpers.GetCoordinator(connection);
            if (host != null)
                return IconFor(host);
            return connection.InProgress ? "000_ServerInProgress_h32bit_16" : "000_ServerDisconnected_h32bit_16";
        }

        public static string IconFor(IXenObject o) => o switch
        {
            VM vm => IconFor(vm),
            Host h => IconFor(h),
            Pool p => IconFor(p),
            SR sr => IconFor(sr),
            _ => "Logo"
        };

        public static string IconFor(VM vm)
        {
            bool disabled = vm.IsHidden();
            if (vm.is_a_template)
                return disabled ? "000_TemplateDisabled_h32bit_16" : vm.DefaultTemplate() ? "000_VMTemplate_h32bit_16" : "000_UserTemplate_h32bit_16";
            if (!vm.ExistsOnServer())
                return disabled ? "000_VMStoppedDisabled_h32bit_16" : "000_StoppedVM_h32bit_16";
            if (vm.current_operations.ContainsValue(vm_operations.migrate_send))
                return vm.power_state switch
                {
                    vm_power_state.Halted => "000_MigrateStoppedVM_h32bit_16",
                    vm_power_state.Suspended or vm_power_state.Paused => "000_MigrateSuspendedVM_h32bit_16",
                    _ => "000_MigrateVM_h32bit_16"
                };
            if (vm.current_operations.Values.Any(VM.is_lifecycle_operation))
                return disabled ? "000_VMStartingDisabled_h32bit_16" : "000_VMStarting_h32bit_16";
            return vm.power_state switch
            {
                vm_power_state.Suspended => disabled ? "000_VMPausedDisabled_h32bit_16" : "000_SuspendVM_h32bit_16",
                vm_power_state.Running => disabled ? "000_VMDisabled_h32bit_16" : "000_StartVM_h32bit_16",
                vm_power_state.Paused => disabled ? "000_VMPausedDisabled_h32bit_16_green" : "000_SuspendVM_h32bit_16_green",
                _ => disabled ? "000_VMStoppedDisabled_h32bit_16" : "000_StoppedVM_h32bit_16"
            };
        }

        public static string IconFor(SR sr)
        {
            if (!sr.HasPBDs() || sr.IsHidden())
                return "000_StorageDisabled_h32bit_16";
            if (sr.IsDetached() || sr.IsBroken() || !sr.MultipathAOK())
                return "000_StorageBroken_h32bit_16";
            return SR.IsDefaultSr(sr) ? "000_StorageDefault_h32bit_16" : "000_Storage_h32bit_16";
        }

        public static string IconFor(Host host)
        {
            var metrics = host.Connection.Resolve(host.metrics);
            if (metrics != null && metrics.live)
            {
                if (host.IsExpired())
                    return "000_Server_h32bit_16-w-alert";
                if (host.HasCrashDumps())
                    return "000_ServerErrorFile_h32bit_16";
                if (host.current_operations.ContainsValue(host_allowed_operations.evacuate) || !host.enabled)
                    return "000_ServerMaintenance_h32bit_16";
                if (Helpers.IsOlderThanCoordinator(host))
                    return "server_up_16";
                if (Helpers.CloudOrGreater(host) &&
                    host.Connection.Cache.Hosts.Select(h => h.latest_synced_updates_applied).Distinct().Count() > 1 &&
                    host.latest_synced_updates_applied != latest_synced_updates_applied_state.yes)
                    return "000_HostUnpatched_h32bit_16";
                return "000_TreeConnected_h32bit_16";
            }
            return host.Connection.InProgress && !host.Connection.IsConnected ? "000_ServerInProgress_h32bit_16" : "000_ServerDisconnected_h32bit_16";
        }

        public static string IconFor(Pool pool)
        {
            if (!pool.Connection.IsConnected)
                return "000_ServerDisconnected_h32bit_16";
            if (pool.Connection.Cache.Hosts.Any(h => h.IsExpired()))
                return "000_Pool_h32bit_16-w-alert";
            if (!pool.IsPoolFullyUpgraded())
                return "pool_up_16";
            if (Helpers.CloudOrGreater(pool.Connection) &&
                pool.Connection.Cache.Hosts.Select(h => h.latest_synced_updates_applied).Distinct().Count() > 1)
                return "pool_unpatched";
            return "000_PoolConnected_h32bit_16";
        }
    }
}
