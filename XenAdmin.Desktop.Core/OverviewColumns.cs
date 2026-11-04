using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Actions.VMActions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;
using XenAdmin.XenSearch;

namespace XenAdmin.Desktop
{
    /// <summary>Values of the Search overview columns, from XenModel's PropertyAccessorHelper and search properties.</summary>
    public static class OverviewColumns
    {
        public static string Title(IXenConnection connection) =>
            connection == null ? "XCP-ng Center: Overview" : $"{Helpers.GetName(Helpers.GetPoolOfOne(connection))}: Overview";

        /// <summary>(text, bar 0..100 or -1) per column, as QueryPanel shows them.</summary>
        public static (string Text, int Bar)[] Values(IXenObject o)
        {
            const string none = "-";
            switch (o)
            {
                case Host h:
                    return new[]
                    {
                        (PropertyAccessorHelper.hostCpuUsageString(h), PropertyAccessorHelper.hostCpuUsageRank(h)),
                        (PropertyAccessorHelper.hostMemoryUsageString(h), PropertyAccessorHelper.hostMemoryUsageRank(h)),
                        (none, -1),
                        (PropertyAccessorHelper.hostNetworkUsageString(h), -1),
                        (Address(h), -1),
                        (h.Uptime()?.ToString() ?? none, -1)
                    };
                case VM vm when !vm.is_a_template:
                    var running = vm.power_state == vm_power_state.Running;
                    return new[]
                    {
                        (running ? PropertyAccessorHelper.vmCpuUsageString(vm) : none, running ? PropertyAccessorHelper.vmCpuUsageRank(vm) : -1),
                        (running ? PropertyAccessorHelper.vmMemoryUsageString(vm) : none, running ? PropertyAccessorHelper.vmMemoryUsageRank(vm) : -1),
                        (running ? PropertyAccessorHelper.vmDiskUsageString(vm) : none, -1),
                        (running ? PropertyAccessorHelper.vmNetworkUsageString(vm) : none, -1),
                        (Address(vm), -1),
                        (vm.RunningTime()?.ToString() ?? none, -1)
                    };
                default:
                    return Enumerable.Repeat((none, -1), 4).Append((o is Pool ? none : "", -1)).Append((o is Pool ? none : "", -1)).ToArray();
            }
        }

        /// <summary>The search "ip_address" property (distinct, VIF/PIF addresses), as the WinForms Address column.</summary>
        public static string Address(IXenObject o)
        {
            var text = PropertyAccessors.Get(PropertyNames.ip_address)(o)?.ToString();
            return string.IsNullOrEmpty(text) ? "-" : text;
        }

    }
}
