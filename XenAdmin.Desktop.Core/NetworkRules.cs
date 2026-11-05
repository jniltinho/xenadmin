using System;
using System.Collections.Generic;
using System.Linq;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// VM network interfaces: the button rules of NetworkList.UpdateEnablement and the validation and result of
    /// VIFDialog (XenAdmin/Controls/NetworkingTab/NetworkList.cs, XenAdmin/Dialogs/VIFDialog.cs).
    /// </summary>
    public static class NetworkRules
    {
        public static bool CanAdd(VM vm) => !vm.Locked;

        /// <summary>Properties and Remove Interface: the VIF can be unplugged, or is not plugged.</summary>
        public static bool CanEdit(VM vm, VIF vif) =>
            vif != null && !vm.Locked && (vif.allowed_operations.Contains(vif_operations.unplug) || !vif.currently_attached);

        public static bool CanRemove(VM vm, VIF vif) => CanEdit(vm, vif);

        /// <summary>Activate/Deactivate: plug a detached VIF or unplug an attached one, when the server allows it.</summary>
        public static bool CanToggle(VM vm, VIF vif) =>
            vif != null && !vm.Locked
            && (vif.currently_attached ? vif.allowed_operations.Contains(vif_operations.unplug) : vif.allowed_operations.Contains(vif_operations.plug));

        public static string ToggleLabel(VIF vif) =>
            vif != null && !vif.currently_attached ? Messages.VM_NETWORK_TAB_ACTIVATE_BUTTON_LABEL : Messages.VM_NETWORK_TAB_DEACTIVATE_BUTTON_LABEL;

        /// <summary>The tooltips of Properties, Remove and Activate/Deactivate (suspended VM, or running without I/O drivers); null when none.</summary>
        public static (string Edit, string Remove, string Toggle) Tooltips(VM vm, VIF vif)
        {
            if (vif == null)
                return (null, null, null);
            var edit = vm.HasNewVirtualizationStates() ? Messages.TOOLTIP_EDIT_NETWORK_IO_DRIVERS : string.Format(Messages.TOOLTIP_EDIT_NETWORK_TOOLS, BrandManager.VmTools);
            if (vm.power_state == vm_power_state.Suspended)
                return (edit, Messages.TOOLTIP_REMOVE_NETWORK_SUSPENDED,
                    vif.currently_attached ? Messages.TOOLTIP_DEACTIVATE_VIF_SUSPENDED : Messages.TOOLTIP_ACTIVATE_VIF_SUSPENDED);
            if (vm.power_state == vm_power_state.Running && !vm.GetVirtualizationStatus(out _).HasFlag(VM.VirtualizationStatus.IoDriversInstalled))
                return (edit,
                    vm.HasNewVirtualizationStates() ? Messages.TOOLTIP_REMOVE_NETWORK_IO_DRIVERS : string.Format(Messages.TOOLTIP_REMOVE_NETWORK_TOOLS, BrandManager.VmTools),
                    vif.currently_attached ? Messages.TOOLTIP_DEACTIVATE_VIF_TOOLS : Messages.TOOLTIP_ACTIVATE_VIF_TOOLS);
            return (null, null, null);
        }

        /// <summary>VIFDialog.LoadDetails: before Stockholm, a pool with a vSwitch Controller manages QoS itself.</summary>
        public static bool QosManagedByController(IXenConnection c)
        {
            var pool = Helpers.GetPoolOfOne(c);
            return !Helpers.StockholmOrGreater(c) && pool != null && pool.vSwitchController();
        }

        /// <summary>AddNetworkButton_Click: the VM already has the most interfaces it can take.</summary>
        public static bool AtMaxInterfaces(VM vm) => vm.VIFs.Count >= vm.MaxVIFsAllowed();

        /// <summary>VIFDialog.LoadNetworks: visible networks that are not bond members (SR-IOV only when the VM recommends it).</summary>
        public static List<XenAPI.Network> Networks(IXenConnection c, VM vm)
        {
            var sriov = vm.HasSriovRecommendation();
            var list = c.Cache.Networks.Where(n => n.Show(false) && !n.IsMember() && (!n.IsSriov() || sriov)).ToList();
            list.Sort();
            return list;
        }

        /// <summary>VIFDialog.ValidateInput: null when OK can be pressed; "" when it cannot but there is nothing to say.</summary>
        public static string Validate(XenAPI.Network network, bool manualMac, string mac, bool qos, string qosLimit)
        {
            if (network == null)
                return Messages.SELECT_NETWORK_TOOLTIP;
            if (manualMac)
            {
                if (string.IsNullOrWhiteSpace(mac))
                    return "";
                if (!Helpers.IsValidMAC(mac))
                    return Messages.MAC_INVALID;
            }
            if (qos)
            {
                if (string.IsNullOrWhiteSpace(qosLimit))
                    return "";
                if (!int.TryParse(qosLimit, out var kbps) || kbps <= 0)
                    return Messages.ENTER_VALID_QOS;
            }
            return null;
        }

        /// <summary>VIFDialog.GetNewSettings: the new VIF, or a copy of the existing one with the dialog's values.</summary>
        public static VIF Settings(VM vm, VIF existing, XenAPI.Network network, string mac, int device, bool qos, string qosLimit)
        {
            var vif = new VIF();
            if (existing != null)
                vif.UpdateFrom(existing);
            vif.Connection = vm.Connection;
            vif.VM = new XenRef<VM>(vm.opaque_ref);
            vif.network = new XenRef<XenAPI.Network>(network.opaque_ref);
            vif.MAC = mac ?? "";
            vif.device = device.ToString();
            if (qos)
                vif.qos_algorithm_type = VIF.RATE_LIMIT_QOS_VALUE;
            else if (existing != null && existing.qos_algorithm_type == VIF.RATE_LIMIT_QOS_VALUE)
                vif.qos_algorithm_type = "";
            if (!string.IsNullOrEmpty(qosLimit))
                vif.qos_algorithm_params = new Dictionary<string, string> { { VIF.KBPS_QOS_PARAMS_KEY, qosLimit } };
            return vif;
        }

        /// <summary>VIFDialog.ChangesHaveBeenMade: an edit that changes nothing is a Cancel.</summary>
        public static bool Changed(VIF existing, XenAPI.Network network, string mac, int device, bool qos, string qosLimit)
        {
            if (existing == null)
                return true;
            if (existing.network.opaque_ref != network.opaque_ref || existing.MAC != (mac ?? "") || existing.device != device.ToString())
                return true;
            if (existing.qos_algorithm_type == VIF.RATE_LIMIT_QOS_VALUE)
                return !qos || !existing.qos_algorithm_params.TryGetValue(VIF.KBPS_QOS_PARAMS_KEY, out var old) || old != qosLimit;
            return string.IsNullOrEmpty(existing.qos_algorithm_type) && qos;
        }

        /// <summary>VIFDialog_FormClosing: a real VM, other than this interface's, that already has this MAC.</summary>
        public static VM DuplicateMac(IEnumerable<IXenConnection> connections, VIF existing, string mac)
        {
            if (string.IsNullOrEmpty(mac) || existing != null && existing.MAC == mac)
                return null;
            foreach (var c in connections.Where(c => c.IsConnected))
                foreach (var vif in c.Cache.VIFs)
                {
                    var vm = c.Resolve(vif.VM);
                    if (vif.opaque_ref != existing?.opaque_ref && vif.MAC == mac && vm != null && vm.IsRealVm())
                        return vm;
                }
            return null;
        }

        /// <summary>An edit of the VIF oldRef on this device is over once a different VIF holds the device.</summary>
        public static bool EditReplaced(IEnumerable<(string Device, string Ref)> vifs, string device, string oldRef) =>
            vifs.Any(v => v.Device == device && v.Ref != oldRef);

        /// <summary>
        /// The row to select after a refresh: the selected device if it is listed; the device being edited (pending) even
        /// while its VIF is recreated; otherwise the first device.
        /// </summary>
        public static string Selection(IList<string> devices, string selected, string pending)
        {
            if (selected != null && devices.Contains(selected))
                return selected;
            if (pending != null)
                return pending;
            return devices.FirstOrDefault();
        }

        // ---------- host and pool networks (NetworkList for a host or pool, NewNetworkWizard) ----------

        /// <summary>NetworkList.BuildList: the visible networks, sorted.</summary>
        public static List<XenAPI.Network> HostNetworks(IXenConnection c)
        {
            var list = c.Cache.Networks.Where(n => n.Show(false)).ToList();
            list.Sort();
            return list;
        }

        /// <summary>NetworkRow.NetworkName: bond members and hidden networks are marked.</summary>
        public static string NetworkName(XenAPI.Network n) =>
            n.IsMember() ? string.Format(Messages.NIC_BONDED_MEMBER, n.Name())
            : !n.Show(false) ? string.Format(Messages.NIC_HIDDEN, n.Name())
            : n.Name();

        /// <summary>UpdateEnablement for a network: Properties (and Remove, unless it has a physical NIC).</summary>
        public static bool CanEditNetwork(IXenObject owner, XenAPI.Network n) =>
            n != null && !owner.Locked && !n.Locked && !n.IsMember() && !n.CreateInProgress() && !n.IsGuestInstallerNetwork();

        /// <summary>
        /// Remove Network, with its tooltip when disabled: a network on a physical NIC cannot be removed; deleting a bond
        /// (DestroyBondCommand and its management-interface checks) is not ported yet.
        /// </summary>
        public static (bool Enabled, string Tooltip) CanRemoveNetwork(IXenObject owner, XenAPI.Network n)
        {
            if (n != null && n.Connection.ResolveAll(n.PIFs).Any(p => p.IsPhysical() && !p.IsBondNIC()))
                return (false, Messages.TOOLTIP_REMOVE_PIF);
            if (n != null && n.IsBond())
                return (false, "Deleting a bond is not available in this version yet.");
            return (CanEditNetwork(owner, n), null);
        }

        /// <summary>The host whose NICs the New Network wizard offers: the host, or the pool coordinator.</summary>
        public static Host WizardHost(IXenObject owner) => owner as Host ?? Helpers.GetCoordinator(owner.Connection);

        /// <summary>NetWName.PopulatePage: "New Network" for external networks, "New Private Network" otherwise, made unique.</summary>
        public static string DefaultNetworkName(IXenConnection c, bool external) =>
            Helpers.MakeUniqueName(external ? Messages.NEWNETWORK_NAME : Messages.NEWNETWORK_VNAME, c.Cache.Networks.Select(n => n.Name()).ToList());

        /// <summary>NetWDetails.PopulatePage: the host's physical NICs that are not bond members.</summary>
        public static List<PIF> WizardNics(Host host) =>
            host.Connection.Cache.PIFs.Where(p => p.host.opaque_ref == host.opaque_ref && p.IsPhysical() && p.Show(false) && !p.IsBondMember())
                .OrderBy(p => p.Name()).ToList();

        /// <summary>NetWDetails.GetVLANList: the VLANs already used on this NIC's device (SR-IOV VLANs apart).</summary>
        public static List<int> VlansInUse(PIF nic) =>
            nic.Connection.Cache.PIFs.Where(p => p.device == nic.device && p.NetworkSriov() == null).Select(p => (int)p.VLAN).ToList();

        /// <summary>NetWDetails (CA-72484, CA-19111): keep the VLAN if free on the NIC, else take the first free one from 1.</summary>
        public static int VlanFor(int current, ICollection<int> inUse)
        {
            if (!inUse.Contains(current))
                return current;
            for (int i = 1; i <= 4094; i++)
                if (!inUse.Contains(i))
                    return i;
            return current;
        }

        public static int MinVlan(IXenConnection c) => Helpers.VLAN0Allowed(c) ? 0 : 1;

        /// <summary>NetWDetails.ValidateVLANValue: (error, message); VLAN 0 is allowed with an information message.</summary>
        public static (bool Error, string Message) ValidateVlan(string text, int min, ICollection<int> inUse)
        {
            if (!int.TryParse((text ?? "").Trim(), out var vlan))
                return (true, Messages.INVALID_NUMBER);
            if (vlan < min || vlan > 4094)
                return (true, string.Format(Messages.NETW_DETAILS_VLAN_RANGE, min, 4094));
            if (inUse.Contains(vlan))
                return (true, Messages.NETW_DETAILS_VLAN_NUMBER_IN_USE);
            return (false, vlan == 0 ? Messages.NETW_VLAN_ZERO : null);
        }

        /// <summary>NetWDetails: the MTU goes from 1500 to the NIC's MTU (at most 9216).</summary>
        public static int MaxMtu(PIF nic) => (int)Math.Min(nic?.MTU ?? XenAPI.Network.MTU_MAX, XenAPI.Network.MTU_MAX);

        /// <summary>NetWDetails.ValidateMtuValue: (error, message); the allowed range is always shown.</summary>
        public static (bool Error, string Message) ValidateMtu(string text, int max)
        {
            var min = (int)XenAPI.Network.MTU_MIN;
            var range = min == max ? string.Format(Messages.ALLOWED_MTU_VALUE, min) : string.Format(Messages.ALLOWED_MTU_RANGE, min, max);
            if (!int.TryParse((text ?? "").Trim(), out var mtu))
                return (true, Messages.INVALID_NUMBER);
            return (mtu < min || mtu > max, range);
        }

        /// <summary>NewNetworkWizard.CreateNetworkObject (external and single-server private networks).</summary>
        public static XenAPI.Network NewNetwork(string name, string description, bool autoPlug, int? mtu)
        {
            var n = new XenAPI.Network { name_label = name.Trim(), name_description = description ?? "", managed = true };
            n.SetAutoPlug(autoPlug);
            if (mtu.HasValue)
                n.MTU = mtu.Value;
            return n;
        }

        /// <summary>The Limit column of the VM Networking tab (VifRow).</summary>
        public static string Limit(VIF vif) => string.IsNullOrEmpty(vif.qos_algorithm_type) ? "" : vif.LimitString();
    }
}
