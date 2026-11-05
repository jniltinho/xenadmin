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

        /// <summary>The Limit column of the VM Networking tab (VifRow).</summary>
        public static string Limit(VIF vif) => string.IsNullOrEmpty(vif.qos_algorithm_type) ? "" : vif.LimitString();
    }
}
