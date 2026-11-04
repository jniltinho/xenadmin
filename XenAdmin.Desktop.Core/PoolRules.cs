using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>NewPoolDialog / ConnectionWrapperWithMoreStuff: who can be coordinator or member, validation and warnings.</summary>
    public static class PoolRules
    {
        public sealed record Candidate(IXenConnection Connection, PoolJoinRules.Reason Reason)
        {
            public string Name => Helpers.GetName(Helpers.GetCoordinator(Connection)) is var n && n != "" ? n : Connection.Hostname;
            public string Description => PoolJoinRules.ReasonMessage(Reason);
            public bool WillBeCoordinator => Reason == PoolJoinRules.Reason.WillBeCoordinator;
            public bool AllowedAsSupporter => Reason == PoolJoinRules.Reason.Allowed;
            public bool CanBeCoordinator => Reason != PoolJoinRules.Reason.Connecting && Reason != PoolJoinRules.Reason.NotConnected
                                            && Reason != PoolJoinRules.Reason.LicenseRestriction && Reason != PoolJoinRules.Reason.IsAPool;
        }

        // shortcut: no AD prompt yet, so a server whose AD config differs is listed as not allowed; upgrade when AD login is ported.
        private const bool AllowSupporterAdConfig = false;

        /// <summary>Connected servers with their join reason against the proposed coordinator, sorted like the WinForms list.</summary>
        public static List<Candidate> Candidates(IEnumerable<IXenConnection> connections, IXenConnection coordinator) =>
            connections.Where(c => c != null && c.IsConnected)
                .Select(c => new Candidate(c, PoolJoinRules.CanJoinPool(c, coordinator, true, AllowSupporterAdConfig)))
                .OrderBy(x => (int)x.Reason).ThenBy(x => x.Name, StringComparer.CurrentCultureIgnoreCase).ToList();

        /// <summary>NewPoolDialog.validToClose; null when Create Pool can be pressed.</summary>
        public static string Invalid(string name, IXenConnection coordinator, int supporters)
        {
            if (string.IsNullOrWhiteSpace(name)) return Messages.POOL_NAME_EMPTY;
            if (coordinator == null || Helpers.GetCoordinator(coordinator) == null) return Messages.NO_ELIGIBLE_COORDINATOR;
            if (PoolJoinRules.WillExceedPoolMaxSize(coordinator, supporters))
                return string.Format(Messages.NEWPOOL_WILL_EXCEED_POOL_MAX_SIZE, InvisibleMessages.LICENSE_BUY_URL);
            return null;
        }

        /// <summary>The Yes/No warnings NewPoolDialog.createPool and GetPermissionForCpuFeatureLevelling ask, in order.</summary>
        public static List<string> Warnings(List<Host> supporters, Host coordinator)
        {
            string Names(IEnumerable<Host> hosts) => string.Join("\n", hosts.Select(h => h.Name()));
            string Pick(List<Host> hosts, string one, string many) => hosts.Count == 0 ? null : string.Format(hosts.Count == 1 ? one : many, Names(hosts));
            var result = new List<string>();
            var packs = PoolJoinRules.HomogeneousSuppPacksDiffering(supporters, coordinator);
            if (packs.Count > 0)
                result.Add(string.Format(packs.Count == 1 ? Messages.NEW_POOL_SUPP_PACK : Messages.NEW_POOL_SUPP_PACKS, string.Join("\n", packs)));
            result.Add(Pick(supporters.FindAll(h => PoolJoinRules.FreeHostPaidCoordinator(h, coordinator, false)), Messages.NEW_POOL_LICENSE_MESSAGE, Messages.NEW_POOL_LICENSE_MESSAGE_MULTIPLE));
            var cpus = supporters.FindAll(h => !PoolJoinRules.CompatibleCPUs(h, coordinator));
            if (cpus.Count > 0)
                result.Add(string.Format(cpus.Count == 1 ? Messages.NEW_POOL_CPU_MASKING_MESSAGE : Messages.NEW_POOL_CPU_MASKING_MESSAGE_MULTIPLE, Names(cpus), BrandManager.ProductBrand));
            var pool = Helpers.GetPoolOfOne(coordinator.Connection);
            if (pool != null)
            {
                var fewer = supporters.Where(h => PoolJoinRules.HostHasFewerFeatures(h, pool)).ToList();
                var more = supporters.Where(h => PoolJoinRules.HostHasMoreFeatures(h, pool)).ToList();
                result.Add(fewer.Count > 0 && more.Count > 0
                    ? Pick(fewer.Union(more).ToList(), Messages.ADD_HOST_TO_POOL_CPU_DOWN_LEVEL_POOL_AND_HOST_MESSAGE, Messages.ADD_HOST_TO_POOL_CPU_DOWN_LEVEL_POOL_AND_HOST_MESSAGE_MULTIPLE)
                    : fewer.Count > 0 ? Pick(fewer, Messages.ADD_HOST_TO_POOL_CPU_DOWN_LEVEL_POOL_MESSAGE, Messages.ADD_HOST_TO_POOL_CPU_DOWN_LEVEL_POOL_MESSAGE_MULTIPLE)
                    : Pick(more, Messages.ADD_HOST_TO_POOL_CPU_DOWN_LEVEL_HOST_MESSAGE, Messages.ADD_HOST_TO_POOL_CPU_DOWN_LEVEL_HOST_MESSAGE_MULTIPLE));
            }
            return result.Where(m => m != null).ToList();
        }

        /// <summary>CreatePoolAction as NewPoolDialog runs it; licence failures are reported through the action's error.</summary>
        public static CreatePoolAction Create(Host coordinator, List<Host> supporters, string name, string description) =>
            new CreatePoolAction(coordinator, supporters, name.Trim(), description ?? "", _ => null, (_, _) => { });
    }
}
