using System.Linq;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// Which 16px icon (file name in XenAdmin/Images, no extension) represents an object.
    /// Pure logic; loading the bitmaps is the front end's job.
    /// </summary>
    public static class IconNames
    {
        public const string Home = "infra_view_16_textured";

        // shortcut: port of XenAdmin/Images.cs GetIconFor (connection, pool, host, VM, SR only),
        // kept in step by hand; move it into XenModel when the WinForms app can be built on Linux CI.
        public static string For(IXenConnection connection)
        {
            var pool = Helpers.GetPool(connection);
            if (pool != null)
                return For(pool);
            var host = Helpers.GetCoordinator(connection);
            if (host != null)
                return For(host);
            return connection.InProgress ? "000_ServerInProgress_h32bit_16" : "000_ServerDisconnected_h32bit_16";
        }

        public static string For(IXenObject o) => o switch
        {
            VM vm => For(vm),
            Host h => For(h),
            Pool p => For(p),
            SR sr => For(sr),
            // Images.GetIconFor(Network / VDI / VM_appliance / Folder)
            XenAPI.Network _ => "000_Network_h32bit_16",
            VDI vdi => vdi.is_a_snapshot ? "virtualstorage_snapshot_16" : "000_VirtualStorage_h32bit_16",
            VM_appliance _ => "000_VirtualAppliance_h32bit_16",
            XenAdmin.Model.Folder f => f.Grey ? "folder_grey" : "000_Folder_open_h32bit_16",
            _ => "Logo"
        };

        public static string For(VM vm)
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

        public static string For(SR sr)
        {
            if (!sr.HasPBDs() || sr.IsHidden())
                return "000_StorageDisabled_h32bit_16";
            if (sr.IsDetached() || sr.IsBroken() || !sr.MultipathAOK())
                return "000_StorageBroken_h32bit_16";
            return SR.IsDefaultSr(sr) ? "000_StorageDefault_h32bit_16" : "000_Storage_h32bit_16";
        }

        public static string For(Host host)
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

        public static string For(Pool pool)
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
