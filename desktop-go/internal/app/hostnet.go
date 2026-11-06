package app

import (
	"fmt"
	"sort"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// The buttons of the Networking and NICs tabs of a server or pool (HostNetworkView, NicsView), the BondDetails
// control (BondEditor), the New Network wizard's choices (NewNetworkWizard) and a network's Properties (EditNetworkPage).

// HostNetworking is the Networking tab of a server or pool with its buttons.
type HostNetworking struct {
	Table    Table         `json:"table"`
	Networks []HostNetwork `json:"networks"` // the rows' networks, as Table.Keys
	CanAdd   bool          `json:"can_add"`  // Add Network...: the server or pool is not busy
}

// HostNetwork is what Properties and Remove Network need about a row (NetworkRules.CanEditNetwork, CanRemoveNetwork).
type HostNetwork struct {
	UUID      string `json:"uuid"`
	CanEdit   bool   `json:"can_edit"`
	CanRemove bool   `json:"can_remove"`
	RemoveTip string `json:"remove_tip,omitempty"`
	Bond      bool   `json:"bond"` // Remove Network deletes the bond (DestroyBondCommand)
}

// HostNics is the NICs tab of a server with its buttons.
type HostNics struct {
	Table   Table     `json:"table"`
	NICs    []HostNic `json:"nics"`     // the rows' PIFs, as Table.Keys
	CanBond bool      `json:"can_bond"` // Create Bond... and Rescan: the server is not busy
}

// HostNic is a NICs row: a bond interface names the network Delete Bond removes ("" when it cannot: BondRules.CanDelete).
type HostNic struct {
	UUID        string `json:"uuid"`
	BondNetwork string `json:"bond_network,omitempty"`
}

// ownerBusy reads the server or pool of a tab: its ref ("" for a pool) and whether an operation runs on it
// (IXenObject.Locked).
func ownerBusy(s *xenapi.Session, kind, uuid string) (xenapi.HostRef, bool, error) {
	switch kind {
	case "hosts":
		r, err := xenapi.Host.GetByUUID(s, uuid)
		if err != nil {
			return "", false, xapi.Sanitize(err)
		}
		cur, err := xenapi.Host.GetCurrentOperations(s, r)
		return r, len(cur) > 0, xapi.Sanitize(err)
	case "pools":
		r, err := xenapi.Pool.GetByUUID(s, uuid)
		if err != nil {
			return "", false, xapi.Sanitize(err)
		}
		cur, err := xenapi.Pool.GetCurrentOperations(s, r)
		return "", len(cur) > 0, xapi.Sanitize(err)
	}
	return "", false, fmt.Errorf("no networks for %s", kind)
}

// HostNetworkingTab is the Networking tab of a server ("hosts") or pool ("pools") with its buttons.
func HostNetworkingTab(c *xapi.Connection, kind, uuid string) (HostNetworking, error) {
	s := c.Session
	host, busy, err := ownerBusy(s, kind, uuid)
	if err != nil {
		return HostNetworking{}, err
	}
	d, err := readHostTabData(s)
	if err != nil {
		return HostNetworking{}, err
	}
	t, err := d.networking(host)
	if err != nil {
		return HostNetworking{}, err
	}
	infos, err := ops.NetworkInfos(s, busy)
	if err != nil {
		return HostNetworking{}, err
	}
	out := HostNetworking{Table: t, Networks: []HostNetwork{}, CanAdd: !busy}
	for _, k := range t.Keys {
		i := infos[k]
		ok, tip := rules.CanRemoveNetwork(i)
		out.Networks = append(out.Networks, HostNetwork{UUID: k, CanEdit: rules.CanEditNetwork(i), CanRemove: ok, RemoveTip: tip, Bond: i.Bond})
	}
	return out, nil
}

// HostNicsTab is the NICs tab of a server with its buttons.
func HostNicsTab(c *xapi.Connection, uuid string) (HostNics, error) {
	s := c.Session
	host, busy, err := ownerBusy(s, "hosts", uuid)
	if err != nil {
		return HostNics{}, err
	}
	d, err := readHostTabData(s)
	if err != nil {
		return HostNics{}, err
	}
	infos, err := ops.NetworkInfos(s, busy)
	if err != nil {
		return HostNics{}, err
	}
	nets, err := xenapi.Network.GetAllRecords(s)
	if err != nil {
		return HostNics{}, xapi.Sanitize(err)
	}
	t := d.nics(host)
	out := HostNics{Table: t, NICs: []HostNic{}, CanBond: !busy}
	for _, k := range t.Keys {
		row := HostNic{UUID: k}
		for _, p := range d.pifs { // NicsView.BondNetwork: the network of a bond interface
			if p.UUID == k && len(p.BondMasterOf) > 0 {
				n := nets[p.Network].UUID
				if ok, _ := rules.CanRemoveNetwork(infos[n]); ok && infos[n].Bond {
					row.BondNetwork = n
				}
			}
		}
		out.NICs = append(out.NICs, row)
	}
	return out, nil
}

// BondChoices are the BondDetails control's NICs and limits.
type BondChoices struct {
	NICs  []BondNIC `json:"nics"`
	Limit int       `json:"limit"` // Helpers.BondSizeLimit
	LACP  bool      `json:"lacp"`  // Helpers.SupportsLinkAggregationBond
}

// BondNIC is a row of BondDetails: a NIC already in a bond cannot be ticked.
type BondNIC struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	InBond bool   `json:"in_bond"`
	MAC    string `json:"mac"`
	Link   string `json:"link"`
	Speed  string `json:"speed"`
	Vendor string `json:"vendor"`
	MTU    int    `json:"mtu"`
}

func (d *hostTabData) vswitch() bool {
	return d.hosts[d.coordinator].SoftwareVersion["network_backend"] == "openvswitch"
}

// bondChoices is BondRules.Candidates (NetworkingHelper.GetAllPhysicalPIFs): the physical NICs of the server, or of
// every server of the pool once per device (the coordinator's first), by name.
func (d *hostTabData) bondChoices(host xenapi.HostRef, nets map[xenapi.NetworkRef]xenapi.NetworkRecord) BondChoices {
	out := BondChoices{NICs: []BondNIC{}, Limit: rules.BondSizeLimit(d.vswitch()), LACP: d.vswitch()}
	var all, list []xenapi.PIFRecord
	for _, p := range d.pifs {
		if (host == "" || p.Host == host) && pifPhysical(p) && len(p.BondMasterOf) == 0 {
			all = append(all, p)
		}
	}
	sort.SliceStable(all, func(i, j int) bool { return all[i].Host == d.coordinator && all[j].Host != d.coordinator })
	seen := map[string]bool{}
	for _, p := range all {
		if !seen[p.Device] {
			seen[p.Device] = true
			list = append(list, p)
		}
	}
	sort.Slice(list, func(i, j int) bool { return d.pifName(list[i]) < d.pifName(list[j]) })
	for _, p := range list {
		n := BondNIC{UUID: p.UUID, Name: d.pifName(p), MAC: p.MAC, Link: d.linkState(p), Speed: "-", MTU: p.MTU}
		if b, ok := d.bonds[p.BondSlaveOf]; ok { // BondRules.InBond
			n.InBond, n.Name = true, n.Name+" "+messages.Format(messages.Get("ALREADY_IN_BOND"), d.pifName(d.pifs[b.Master]))
		}
		if net, ok := nets[p.Network]; ok && len(net.PIFs) > 1 {
			n.Link = d.networkLink(net)
		}
		if m, ok := d.pm[p.Metrics]; ok {
			n.Vendor = m.VendorName
			if m.Carrier {
				n.Speed = fmt.Sprintf("%d Mbit/s", m.Speed)
			}
		}
		out.NICs = append(out.NICs, n)
	}
	return out
}

// NewNetworkChoices are the New Network wizard's defaults and NICs for a server or pool.
type NewNetworkChoices struct {
	Title        string      `json:"title"`
	Host         string      `json:"host"` // NetworkRules.WizardHost: the server, or the coordinator
	ExternalName string      `json:"external_name"`
	PrivateName  string      `json:"private_name"`
	NICs         []WizardNIC `json:"nics"`
	MinVLAN      int         `json:"min_vlan"`
	Bond         BondChoices `json:"bond"`
}

// WizardNIC is a NIC of the Interface page: its MTU range's top and the VLANs already used on its device.
type WizardNIC struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	MaxMTU int    `json:"max_mtu"`
	VLANs  []int  `json:"vlans"`
}

// NewNetworkOptions reads the New Network wizard's choices (NetWName, NetWDetails, NetWBondDetails).
func NewNetworkOptions(c *xapi.Connection, kind, uuid string) (NewNetworkChoices, error) {
	s := c.Session
	owner, _, err := ownerBusy(s, kind, uuid)
	if err != nil {
		return NewNetworkChoices{}, err
	}
	d, err := readHostTabData(s)
	if err != nil {
		return NewNetworkChoices{}, err
	}
	nets, err := xenapi.Network.GetAllRecords(s)
	if err != nil {
		return NewNetworkChoices{}, xapi.Sanitize(err)
	}
	host := owner
	if host == "" {
		host = d.coordinator
	}
	taken := map[string]bool{}
	for _, n := range nets {
		taken[d.networkName(n)] = true
	}
	out := NewNetworkChoices{Title: messages.Format(messages.Get("NEWNETWORKWIZARD_TITLE"), PoolName(c)), Host: d.hosts[host].UUID,
		ExternalName: rules.UniqueName(messages.Get("NEWNETWORK_NAME"), taken), PrivateName: rules.UniqueName(messages.Get("NEWNETWORK_VNAME"), taken),
		NICs: []WizardNIC{}, MinVLAN: rules.MinVLAN(d.vswitch()), Bond: d.bondChoices(owner, nets)}
	var list []xenapi.PIFRecord // NetworkRules.WizardNics
	for _, p := range d.pifs {
		if p.Host == host && pifPhysical(p) && p.Managed && p.OtherConfig["HideFromXenCenter"] != "true" && (p.BondSlaveOf == "" || p.BondSlaveOf == nullRef) {
			list = append(list, p)
		}
	}
	sort.Slice(list, func(i, j int) bool { return d.pifName(list[i]) < d.pifName(list[j]) })
	inUse, err := ops.VLANsInUse(s)
	if err != nil {
		return NewNetworkChoices{}, err
	}
	for _, p := range list {
		out.NICs = append(out.NICs, WizardNIC{UUID: p.UUID, Name: d.pifName(p), MaxMTU: max(rules.MTUMin, rules.NICMaxMTU(p.MTU)), VLANs: inUse[p.Device]})
	}
	return out, nil
}

// NetworkEdit is a network's Properties: General (name, description) and Network Settings (auto-add, MTU).
type NetworkEdit struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Auto        bool           `json:"auto"`
	AutoEnabled bool           `json:"auto_enabled"` // not the guest installer network
	MTU         int            `json:"mtu"`
	MTUState    rules.MTUState `json:"mtu_state"`
}

// noTools is a running VM without I/O drivers (VM.GetVirtualizationStatus), for the MTU rules.
func noTools(s *xenapi.Session) func(xenapi.VMRecord) bool {
	return func(vm xenapi.VMRecord) bool { io, _, _ := vmTools(s, vm); return !io }
}

// ReadNetworkEdit reads a network for its Properties.
func ReadNetworkEdit(c *xapi.Connection, uuid string) (NetworkEdit, error) {
	s := c.Session
	ref, err := xenapi.Network.GetByUUID(s, uuid)
	if err != nil {
		return NetworkEdit{}, xapi.Sanitize(err)
	}
	n, err := xenapi.Network.GetRecord(s, ref)
	if err != nil {
		return NetworkEdit{}, xapi.Sanitize(err)
	}
	m, err := ops.NetworkMTU(s, uuid, noTools(s))
	return NetworkEdit{Name: n.NameLabel, Description: n.NameDescription, Auto: n.OtherConfig["automatic"] != "false",
		AutoEnabled: n.OtherConfig["is_guest_installer_network"] != "true", MTU: n.MTU, MTUState: m}, err
}

// SaveNetwork saves a network's Properties (ops.UpdateNetwork).
func SaveNetwork(c *xapi.Connection, uuid string, ch ops.NetworkChanges) error {
	return ops.UpdateNetwork(c, uuid, ch, noTools(c.Session))
}
