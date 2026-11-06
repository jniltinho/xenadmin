package rules

import (
	"errors"
	"slices"
	"strconv"
	"strings"

	"xcpng-center/internal/messages"
)

// Bonded and host networks: BondRules and the host-network part of NetworkRules (XenAdmin.Desktop.Core/BondRules.cs,
// NetworkRules.cs), from BondDetails, DestroyBondCommand, NetWDetails and NetWName.

// The MTU range of a network (XenAPI.Network.MTU_MIN, MTU_DEFAULT, MTU_MAX).
const (
	MTUMin     = 1500
	MTUDefault = 1500
	MTUMax     = 9216
)

// Warning is BondRules.Warning: a question before a bond is created or deleted; Block shows Text as an error and stops.
type Warning struct {
	Block bool   `json:"block"`
	Title string `json:"title"`
	Text  string `json:"text"`
	OK    string `json:"ok"` // the button that goes on ("" when blocked)
}

// BondSizeLimit is Helpers.BondSizeLimit: 4 NICs per bond on the vSwitch backend, else 2.
func BondSizeLimit(vswitch bool) int {
	if vswitch {
		return 4
	}
	return 2
}

// BondValidSelection is BondDetails.SetValid: two NICs up to the limit of the pool.
func BondValidSelection(count, limit int) bool { return count >= 2 && count <= limit }

// BondMaxMTU is BondDetails.UpdateMtuControls: 1500 up to the smallest MTU of the members (at most 9216).
func BondMaxMTU(mtus []int) int {
	if len(mtus) == 0 {
		return MTUMax
	}
	return min(slices.Min(mtus), MTUMax)
}

// BondName is BondDetails.BondName (PIF.BondSuffix): "Bond 0+1", the members' device numbers sorted.
func BondName(devices []string) string {
	ids := make([]string, len(devices))
	for i, d := range devices {
		ids[i] = strings.ReplaceAll(d, "eth", "")
	}
	slices.Sort(ids)
	return messages.Format(messages.Get("PIF_BOND"), strings.Join(ids, "+"))
}

// BondCreateWarning is BondDetails.CanCreateBond: nil when nothing needs confirming. primary/secondary: a member is
// the primary or a secondary management interface; clustering: one is used by clustering.
func BondCreateWarning(primary, secondary, clustering, haEnabled bool, poolName string) *Warning {
	title := messages.Get("BOND_CREATE")
	switch {
	case primary && secondary:
		return &Warning{Block: true, Title: title, Text: messages.Get("BOND_CREATE_WILL_DISTURB_BOTH")}
	case primary && haEnabled:
		return &Warning{Block: true, Title: title, Text: messages.Format(messages.Get("BOND_CREATE_HA_ENABLED"), poolName)}
	case primary:
		return &Warning{Title: title, Text: messages.Format(messages.Get("BOND_CREATE_WILL_DISTURB_PRIMARY"), brandConsole),
			OK: strings.ReplaceAll(messages.Get("BOND_CREATE_CONTINUE"), "&", "")}
	case secondary && clustering:
		return &Warning{Title: title, Text: messages.Get("BOND_CREATE_WILL_DISTURB_CLUSTERING"), OK: "OK"}
	case secondary:
		return &Warning{Title: title, Text: messages.Get("BOND_CREATE_WILL_DISTURB_SECONDARY"), OK: "OK"}
	}
	return nil
}

// BondDeleteWarning is BondRules.DeleteWarning (DestroyBondCommand.RunCore): always a question, or an error under HA.
// nic is the name of the bond's interface ("Bond 1+2").
func BondDeleteWarning(nic string, primary, secondary, haEnabled bool, poolName string) Warning {
	title, msg := messages.Get("DELETE_BOND"), messages.Format(messages.Get("DELETE_BOND_MESSAGE"), nic)
	goOn := strings.ReplaceAll(messages.Get("BOND_DELETE_CONTINUE"), "&", "")
	switch {
	case primary && haEnabled:
		return Warning{Block: true, Title: title, Text: messages.Format(messages.Get("BOND_DELETE_HA_ENABLED"), nic, poolName)}
	case primary && secondary:
		return Warning{Title: title, Text: messages.Format(messages.Get("BOND_DELETE_WILL_DISTURB_BOTH"), msg, brandConsole), OK: goOn}
	case primary:
		return Warning{Title: title, Text: messages.Format(messages.Get("BOND_DELETE_WILL_DISTURB_PRIMARY"), msg, brandConsole), OK: goOn}
	case secondary:
		return Warning{Title: title, Text: messages.Format(messages.Get("BOND_DELETE_WILL_DISTURB_SECONDARY"), msg), OK: "OK"}
	}
	return Warning{Title: title, Text: msg, OK: "OK"}
}

// brandConsole is BrandManager.BrandConsole.
const brandConsole = "XCP-ng Center"

// VLANFor is NetworkRules.VlanFor (NetWDetails, CA-72484): the VLAN if free on the NIC, else the first free one from 1.
func VLANFor(current int, inUse []int) int {
	if !slices.Contains(inUse, current) {
		return current
	}
	for i := 1; i <= 4094; i++ {
		if !slices.Contains(inUse, i) {
			return i
		}
	}
	return current
}

// MinVLAN is NetworkRules.MinVlan (Helpers.VLAN0Allowed): VLAN 0 on the vSwitch backend.
func MinVLAN(vswitch bool) int {
	if vswitch {
		return 0
	}
	return 1
}

// ValidateVLAN is NetworkRules.ValidateVlan (NetWDetails.ValidateVLANValue): the VLAN, or an error; VLAN 0 is allowed
// with an information message.
func ValidateVLAN(text string, min int, inUse []int) (vlan int, msg string, ok bool) {
	vlan, err := strconv.Atoi(strings.TrimSpace(text))
	switch {
	case err != nil:
		return 0, messages.Get("INVALID_NUMBER"), false
	case vlan < min || vlan > 4094:
		return 0, messages.Format(messages.Get("NETW_DETAILS_VLAN_RANGE"), min, 4094), false
	case slices.Contains(inUse, vlan):
		return 0, messages.Get("NETW_DETAILS_VLAN_NUMBER_IN_USE"), false
	case vlan == 0:
		return 0, messages.Get("NETW_VLAN_ZERO"), true
	}
	return vlan, "", true
}

// NICMaxMTU is NetworkRules.MaxMtu: up to the NIC's MTU, at most 9216.
func NICMaxMTU(mtu int) int { return min(mtu, MTUMax) }

// ValidateMTU is NetworkRules.ValidateMtu (NetWDetails.ValidateMtuValue): the MTU, the allowed range (always shown)
// and whether it is inside it. lo is the range's minimum (1500, or a lower MTU set on the CLI).
func ValidateMTU(text string, lo, hi int) (mtu int, msg string, ok bool) {
	msg = messages.Format(messages.Get("ALLOWED_MTU_RANGE"), lo, hi)
	if lo == hi {
		msg = messages.Format(messages.Get("ALLOWED_MTU_VALUE"), lo)
	}
	mtu, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		return 0, messages.Get("INVALID_NUMBER"), false
	}
	return mtu, msg, mtu >= lo && mtu <= hi
}

// NetworkInfo is what the Networking tab's rules read of a network.
type NetworkInfo struct {
	Busy        bool // the network, or the server or pool, has an operation running (IXenObject.Locked)
	BondMember  bool // Network.IsMember: a PIF of it is a bond member
	Bond        bool // Network.IsBond: a PIF of it is a bond interface
	Physical    bool // a PIF of it is physical and not a bond interface
	Creating    bool // other_config XenCenterCreateInProgress (a bond being made)
	GuestInstal bool // other_config is_guest_installer_network
}

// CanEditNetwork is NetworkRules.CanEditNetwork: Properties for a network that is not busy, not a bond member, not
// being made and not the guest installer network.
func CanEditNetwork(n NetworkInfo) bool {
	return !n.Busy && !n.BondMember && !n.Creating && !n.GuestInstal
}

// CanRemoveNetwork is NetworkRules.CanRemoveNetwork: a network on a physical NIC cannot be removed (with its tooltip);
// a bond network is removed by deleting the bond.
func CanRemoveNetwork(n NetworkInfo) (ok bool, tip string) {
	switch {
	case n.Physical:
		return false, messages.Get("TOOLTIP_REMOVE_PIF")
	case n.Bond:
		return !n.Busy, ""
	}
	return CanEditNetwork(n), ""
}

// MTUState is the MTU field of a network's Properties.
type MTUState struct {
	Shown   bool   `json:"shown"`
	Enabled bool   `json:"enabled"`
	Min     int    `json:"min"`
	Max     int    `json:"max"`
	Note    string `json:"note"` // why it cannot be changed
}

// NetworkMTU is EditNetworkPage.SetMTUControlEnablement: hidden on a tunnel network (Network.CanUseJumboFrames), fixed
// on a private network, on a management interface or while running VMs without tools use it; the range keeps a value
// set on the CLI.
func NetworkMTU(mtu int, tunnel, private, management bool, managementName string, vmsWithoutTools bool) MTUState {
	m := MTUState{Shown: !tunnel, Min: min(mtu, MTUMin), Max: max(mtu, MTUMax)}
	switch {
	case tunnel, private:
	case management:
		m.Note = messages.Format(messages.Get("CANNOT_CONFIGURE_JUMBO_DISTURB_MANAGEMENT"), managementName)
	case vmsWithoutTools:
		m.Note = messages.Get("CANNOT_CONFIGURE_JUMBO_VM_NO_TOOLS")
	default:
		m.Enabled = true
	}
	return m
}

// NICState is the NIC and VLAN fields of a network's Properties.
type NICState struct {
	Shown   bool   `json:"shown"`   // EditNetworkPage.Editable: a private network or a VLAN (not a NIC, bond, tunnel or SR-IOV network)
	Enabled bool   `json:"enabled"` // HostPNICList.Enabled, and the VLAN's unless <internal> is picked
	Note    string `json:"note"`    // warningText: why they cannot be changed
}

// NetworkNICs is EditNetworkPage.Editable and SetNetSettingsEnablement: a private network can always move to a NIC;
// a VLAN cannot change while it is a management interface or VMs are attached to it (its PIFs would not unplug).
func NetworkNICs(private, editable, management bool, managementName string, vmsAttached bool) NICState {
	switch {
	case private:
		return NICState{Shown: true, Enabled: true}
	case !editable:
		return NICState{}
	case management:
		return NICState{Shown: true, Note: messages.Format(messages.Get("CANNOT_CONFIGURE_NET_DISTURB_MANAGEMENT"), managementName)}
	case vmsAttached:
		return NICState{Shown: true, Note: messages.Get("CANNOT_CONFIGURE_NET_VMS_ATTACHED")}
	}
	return NICState{Shown: true, Enabled: true}
}

// BondChanges is EditNetworkPage.BondModeHasChanged and HashingAlgorithmHasChanged: what to set on a bond network's
// bonds for the picked mode and hashing algorithm. cur* are the bond's ("" hashing: none or unknown); the hashing
// algorithm is set only with LACP on the vSwitch (the LACP choices are hidden otherwise).
func BondChanges(curMode, curHashing, mode, hashing string, lacp bool) (setMode, setHashing bool, err error) {
	switch {
	case mode != "balance-slb" && mode != "active-backup" && mode != "lacp":
		return false, false, errors.New("unknown bond mode " + strconv.Quote(mode))
	case mode == "lacp" && !lacp && curMode != "lacp":
		return false, false, errors.New("LACP needs the vSwitch") // Helpers.SupportsLinkAggregationBond
	case mode == "lacp" && lacp && hashing != "src_mac" && hashing != "tcpudp_ports":
		return false, false, errors.New("pick a load balancing method for LACP")
	}
	if mode != "lacp" || !lacp {
		hashing = "" // Bond.hashing_algoritm.unknown: never set
	}
	return mode != curMode, hashing != "" && hashing != curHashing, nil
}
