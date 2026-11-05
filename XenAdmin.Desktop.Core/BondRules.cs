using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// Bonded networks: the rules of BondDetails (New Network wizard, XenAdmin/Controls/BondDetails.cs) and
    /// DestroyBondCommand (XenAdmin/Commands/DestroyBondCommand.cs).
    /// </summary>
    public static class BondRules
    {
        /// <summary>BondDetails.SetHost/SetPool: the physical NICs of the host, or of the pool coordinator.</summary>
        public static List<PIF> Candidates(IXenObject owner) => owner switch
        {
            Host h => NetworkingHelper.GetAllPhysicalPIFs(h),
            Pool p => NetworkingHelper.GetAllPhysicalPIFs(p),
            _ => new List<PIF>()
        };

        /// <summary>BondDetails.PIFDescription: a NIC already in a bond cannot be picked ("" when free).</summary>
        public static string InBond(PIF pif) => pif.BondMemberOf() is Bond b ? string.Format(Messages.ALREADY_IN_BOND, b.Name()) : "";

        public static int SizeLimit(IXenConnection c) => Helpers.BondSizeLimit(c);

        /// <summary>BondDetails.SetValid: two NICs up to the limit of the pool.</summary>
        public static bool ValidSelection(int count, int limit) => count >= 2 && count <= limit;

        /// <summary>BondDetails.ToggleNICRowCheckable: a free NIC can be ticked while the limit is not reached.</summary>
        public static bool CanTick(PIF pif, int ticked, int limit) => InBond(pif) == "" && ticked < limit;

        /// <summary>LACP modes need Tampa or newer with the vSwitch (Helpers.SupportsLinkAggregationBond).</summary>
        public static bool SupportsLacp(IXenConnection c) => Helpers.SupportsLinkAggregationBond(c);

        /// <summary>BondDetails.UpdateMtuControls: the MTU goes from 1500 to the smallest MTU of the members (at most 9216).</summary>
        public static int MaxMtu(IReadOnlyCollection<PIF> members) =>
            (int)Math.Min(members.Select(p => p.MTU).DefaultIfEmpty(XenAPI.Network.MTU_MAX).Min(), XenAPI.Network.MTU_MAX);

        /// <summary>BondDetails.BondName: "Bond 0+1".</summary>
        public static string BondName(List<PIF> members) => string.Format(Messages.PIF_BOND, PIF.BondSuffix(members));

        /// <summary>A warning before a bond is created or deleted: Block shows the text as an error and stops.</summary>
        public sealed record Warning(bool Block, string Title, string Text, string OkText);

        /// <summary>BondDetails.CanCreateBond: null when nothing needs confirming.</summary>
        public static Warning CreateWarning(IXenConnection c, List<PIF> members)
        {
            var primary = NetworkingHelper.ContainsPrimaryManagement(members);
            var secondary = NetworkingHelper.ContainsSecondaryManagement(members);
            if (primary && secondary)
                return new Warning(true, Messages.BOND_CREATE, Messages.BOND_CREATE_WILL_DISTURB_BOTH, null);
            if (primary)
            {
                var pool = Helpers.GetPool(c);
                if (pool != null && pool.ha_enabled)
                    return new Warning(true, Messages.BOND_CREATE, string.Format(Messages.BOND_CREATE_HA_ENABLED, pool.Name()), null);
                return new Warning(false, Messages.BOND_CREATE, string.Format(Messages.BOND_CREATE_WILL_DISTURB_PRIMARY, BrandManager.BrandConsole), Messages.BOND_CREATE_CONTINUE);
            }
            if (secondary)
                return new Warning(false, Messages.BOND_CREATE,
                    NetworkingHelper.ContainsClusteringPif(members) ? Messages.BOND_CREATE_WILL_DISTURB_CLUSTERING : Messages.BOND_CREATE_WILL_DISTURB_SECONDARY, Messages.OK);
            return null;
        }

        /// <summary>The bond behind a bond network (from its PIF on the coordinator), or null.</summary>
        public static Bond BondOf(XenAPI.Network n) =>
            n?.Connection.ResolveAll(n.PIFs).Select(p => p.BondInterfaceOf()).FirstOrDefault(b => b != null);

        /// <summary>NICPage / NetworkList: Delete Bond for a bond network that is not busy.</summary>
        public static bool CanDelete(XenAPI.Network n) => BondOf(n) is Bond b && !b.Locked && !n.Locked;

        /// <summary>DestroyBondCommand.RunCore: the confirmation (always asked), or a blocking error under HA.</summary>
        public static Warning DeleteWarning(XenAPI.Network n)
        {
            var pifs = n.Connection.ResolveAll(n.PIFs);
            var msg = string.Format(Messages.DELETE_BOND_MESSAGE, pifs[0].Name());
            var primary = NetworkingHelper.ContainsPrimaryManagement(pifs);
            var secondary = NetworkingHelper.ContainsSecondaryManagement(pifs);
            if (primary)
            {
                var pool = Helpers.GetPool(n.Connection);
                if (pool != null && pool.ha_enabled)
                    return new Warning(true, Messages.DELETE_BOND, string.Format(Messages.BOND_DELETE_HA_ENABLED, pifs[0].Name(), pool.Name()), null);
                return new Warning(false, Messages.DELETE_BOND,
                    string.Format(secondary ? Messages.BOND_DELETE_WILL_DISTURB_BOTH : Messages.BOND_DELETE_WILL_DISTURB_PRIMARY, msg, BrandManager.BrandConsole), Messages.BOND_DELETE_CONTINUE);
            }
            if (secondary)
                return new Warning(false, Messages.DELETE_BOND, string.Format(Messages.BOND_DELETE_WILL_DISTURB_SECONDARY, msg), Messages.OK);
            return new Warning(false, Messages.DELETE_BOND, msg, Messages.OK);
        }
    }
}
