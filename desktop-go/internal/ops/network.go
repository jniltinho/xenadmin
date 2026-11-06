package ops

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// Server networks and bonds: CreateBondAction, DestroyBondAction, RescanPIFsAction, NetworkAction (create, destroy),
// UnplugPlugNetworkAction and the Network Settings page (XenModel/Actions/Network, XenAdmin/SettingsPanels/
// EditNetworkPage.cs), with the rules of XenAdmin.Desktop.Core/BondRules.cs and NetworkRules.cs.

const (
	createInProgress = "XenCenterCreateInProgress" // Network.CREATE_IN_PROGRESS
	purposeKey       = "management_purpose"        // PIF.GetManagementPurpose
)

// netState is what these operations read of a connection.
type netState struct {
	s           *xenapi.Session
	hosts       map[xenapi.HostRef]xenapi.HostRecord
	pifs        map[xenapi.PIFRef]xenapi.PIFRecord
	bonds       map[xenapi.BondRef]xenapi.BondRecord
	nets        map[xenapi.NetworkRef]xenapi.NetworkRecord
	vlans       map[xenapi.VLANRef]xenapi.VLANRecord
	pool        xenapi.PoolRecord
	coordinator xenapi.HostRef
}

func readNet(s *xenapi.Session) (*netState, error) {
	d := &netState{s: s}
	var err error
	if d.hosts, err = xenapi.Host.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	if d.pifs, err = xenapi.PIF.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	if d.bonds, err = xenapi.Bond.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	if d.nets, err = xenapi.Network.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	if d.vlans, err = xenapi.VLAN.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	for _, p := range pools {
		d.pool, d.coordinator = p, p.Master
	}
	if _, ok := d.hosts[d.coordinator]; !ok {
		return nil, errors.New("the pool's coordinator is gone")
	}
	return d, nil
}

func valid[T ~string](ref T) bool { return ref != "" && ref != "OpaqueRef:NULL" }

// physical is PIF.IsPhysical: not a VLAN, tunnel access or SR-IOV logical PIF (a bond interface is one).
func physical(p xenapi.PIFRecord) bool {
	return p.VLAN == -1 && len(p.TunnelAccessPIFOf) == 0 && len(p.SriovLogicalPIFOf) == 0
}

// sriov is PIF.NetworkSriov != null: an SR-IOV logical PIF, or a VLAN on one.
func (d *netState) sriov(p xenapi.PIFRecord) bool {
	if len(p.SriovLogicalPIFOf) > 0 {
		return true
	}
	v, ok := d.vlans[p.VLANMasterOf]
	return p.VLAN >= 0 && ok && len(d.pifs[v.TaggedPIF].SriovLogicalPIFOf) > 0
}

// vlansInUse is NetWDetails.GetVLANList: the VLANs already used on a device, by its SR-IOV network (onSriov:
// checkBoxSriov ticked) or by the other networks.
func (d *netState) vlansInUse(device string, onSriov bool) []int {
	out := []int{}
	for _, p := range d.pifs {
		if p.Device == device && d.sriov(p) == onSriov {
			out = append(out, p.VLAN)
		}
	}
	return out
}

// sriovLogicalPIF is NewNetworkWizard.CreateNonBonded's PIF for a VLAN on the SR-IOV network: the logical PIF of the
// network_sriov the NIC is the physical PIF of (network_sriov.logical_PIF, read from its other side,
// PIF.sriov_logical_PIF_of); "" when the NIC has no SR-IOV network or its logical PIF is gone.
func (d *netState) sriovLogicalPIF(nic xenapi.PIFRecord) xenapi.PIFRef {
	if len(nic.SriovPhysicalPIFOf) == 0 {
		return ""
	}
	for r, p := range d.pifs {
		if slices.Contains(p.SriovLogicalPIFOf, nic.SriovPhysicalPIFOf[0]) {
			return r
		}
	}
	return ""
}

// VLANsInUse are the VLANs already used on each device, for the New Network wizard: by the other networks, and by
// the device's SR-IOV network.
func VLANsInUse(s *xenapi.Session) (others, sriov map[string][]int, err error) {
	d, err := readNet(s)
	if err != nil {
		return nil, nil, err
	}
	others, sriov = map[string][]int{}, map[string][]int{}
	for _, p := range d.pifs {
		others[p.Device], sriov[p.Device] = d.vlansInUse(p.Device, false), d.vlansInUse(p.Device, true)
	}
	return others, sriov, nil
}

// Vswitch is Host.vSwitchNetworkBackend of the coordinator (bonds of 4 NICs, LACP, VLAN 0).
func (d *netState) vswitch() bool {
	return d.hosts[d.coordinator].SoftwareVersion["network_backend"] == "openvswitch"
}

// shown is Network.Show(false): not the guest installer network, not hidden, no hidden or bond member PIF.
func (d *netState) shown(n xenapi.NetworkRecord) bool {
	if n.OtherConfig["is_guest_installer_network"] == "true" || n.OtherConfig["HideFromXenCenter"] == "true" {
		return false
	}
	for _, r := range n.PIFs {
		if p := d.pifs[r]; !p.Managed || p.OtherConfig["HideFromXenCenter"] == "true" || valid(p.BondSlaveOf) {
			return false
		}
	}
	return true
}

// secondary is PIF.IsSecondaryManagementInterface: not the management PIF, with an IP, on a shown network; with
// hidden (Network.Show(true), as NetworkingHelper.ContainsSecondaryManagement) any network of managed PIFs counts.
func (d *netState) secondary(p xenapi.PIFRecord, hidden bool) bool {
	if p.Management || p.IPConfigurationMode != xenapi.IPConfigurationModeDHCP && p.IPConfigurationMode != xenapi.IPConfigurationModeStatic {
		return false
	}
	n, ok := d.nets[p.Network]
	if !ok || !hidden {
		return ok && d.shown(n)
	}
	return !slices.ContainsFunc(n.PIFs, func(r xenapi.PIFRef) bool { return !d.pifs[r].Managed })
}

// management is PIF.IsManagementInterface(false): the primary management PIF (not on the guest installer network)
// or a secondary one.
func (d *netState) management(p xenapi.PIFRecord) bool {
	return p.Management && d.nets[p.Network].OtherConfig["is_guest_installer_network"] != "true" || d.secondary(p, false)
}

// managementName is PIF.ManagementInterfaceNameOrUnknown.
func managementName(p xenapi.PIFRecord) string {
	if p.Management {
		return "Management" // Messages.MANAGEMENT
	}
	if v := p.OtherConfig[purposeKey]; v != "" {
		return v
	}
	return "unknown" // Messages.NETWORKING_PROPERTIES_PURPOSE_UNKNOWN
}

// sortByName sorts PIFs as NetworkingHelper.Sort (PIF.Name: by device number).
func (d *netState) sortByName(refs []xenapi.PIFRef) {
	slices.SortFunc(refs, func(a, b xenapi.PIFRef) int { return strings.Compare(d.pifs[a].Device, d.pifs[b].Device) })
}

// poolName is Pool.Name: the pool's name, else its coordinator's.
func (d *netState) poolName() string {
	if d.pool.NameLabel != "" {
		return d.pool.NameLabel
	}
	return d.hosts[d.coordinator].NameLabel
}

// NetworkInfos are the Networking tab's facts of each network (rules.NetworkInfo), by uuid; busy is the server's or
// pool's own state, added to each.
func NetworkInfos(s *xenapi.Session, busy bool) (map[string]rules.NetworkInfo, error) {
	d, err := readNet(s)
	if err != nil {
		return nil, err
	}
	out := map[string]rules.NetworkInfo{}
	for _, n := range d.nets {
		out[n.UUID] = d.info(n, busy)
	}
	return out, nil
}

func (d *netState) info(n xenapi.NetworkRecord, busy bool) rules.NetworkInfo {
	i := rules.NetworkInfo{Busy: busy || len(n.CurrentOperations) > 0, Creating: n.OtherConfig[createInProgress] == "true",
		GuestInstal: n.OtherConfig["is_guest_installer_network"] == "true"}
	for _, r := range n.PIFs {
		p := d.pifs[r]
		i.BondMember = i.BondMember || valid(p.BondSlaveOf)
		i.Bond = i.Bond || len(p.BondMasterOf) > 0
		i.Physical = i.Physical || physical(p) && len(p.BondMasterOf) == 0
	}
	return i
}

// ---------- bonds ----------

// BondSpec is the BondDetails control's values (Create Bond, the New Network wizard's Bonded Network).
type BondSpec struct {
	Members []string `json:"members"` // the ticked NICs (PIF uuids of one server)
	Mode    string   `json:"mode"`    // balance-slb, active-backup or lacp
	Hashing string   `json:"hashing"` // lacp: src_mac or tcpudp_ports
	MTU     int      `json:"mtu"`
	Auto    bool     `json:"auto"` // add this network to new VMs
}

// members checks the ticked NICs (BondDetails: physical, not bond interfaces or members, one per device, up to the
// limit); a pool's NICs may come from different servers (GetAllPhysicalPIFs(pool)), only their devices count.
func (d *netState) members(uuids []string) ([]xenapi.PIFRef, error) {
	var refs []xenapi.PIFRef
	devices := map[string]bool{}
	for r, p := range d.pifs {
		if !slices.Contains(uuids, p.UUID) {
			continue
		}
		if !physical(p) || len(p.BondMasterOf) > 0 || valid(p.BondSlaveOf) || devices[p.Device] {
			return nil, fmt.Errorf("%s cannot be bonded", p.Device)
		}
		devices[p.Device] = true
		refs = append(refs, r)
	}
	if len(refs) != len(uuids) || !rules.BondValidSelection(len(refs), rules.BondSizeLimit(d.vswitch())) {
		return nil, errors.New("select two or more NICs that are not already bonded") // NewNetworkWizard.Validate
	}
	d.sortByName(refs)
	return refs, nil
}

// createWarning is BondRules.CreateWarning for these members.
func (d *netState) createWarning(refs []xenapi.PIFRef) *rules.Warning {
	clustering := map[xenapi.PIFRef]bool{} // PIF.IsUsedByClustering
	if chs, err := xenapi.ClusterHost.GetAllRecords(d.s); err == nil {
		for _, ch := range chs {
			clustering[ch.PIF] = true
		}
	}
	var primary, secondary, cluster bool
	for _, r := range refs {
		primary = primary || d.pifs[r].Management
		secondary = secondary || d.secondary(d.pifs[r], true)
		cluster = cluster || clustering[r]
	}
	return rules.BondCreateWarning(primary, secondary, cluster, d.pool.HaEnabled, d.poolName())
}

// BondCreateWarning is the question Create Bond asks first (nil: none).
func BondCreateWarning(c *xapi.Connection, members []string) (*rules.Warning, error) {
	d, err := readNet(c.Session)
	if err != nil {
		return nil, err
	}
	refs, err := d.members(members)
	if err != nil {
		return nil, err
	}
	return d.createWarning(refs), nil
}

// moveManagementName is NetworkingActionHelpers.MoveManagementInterfaceName: a secondary management interface's name
// (other_config management_purpose) goes from one PIF to the other.
func moveManagementName(s *xenapi.Session, from, to xenapi.PIFRef) error {
	other, err := xenapi.PIF.GetOtherConfig(s, from)
	if err != nil || other[purposeKey] == "" {
		return xapi.Sanitize(err)
	}
	_ = xenapi.PIF.RemoveFromOtherConfig(s, to, purposeKey) // add fails on an existing key
	if err := xenapi.PIF.AddToOtherConfig(s, to, purposeKey, other[purposeKey]); err != nil {
		return xapi.Sanitize(err)
	}
	return xapi.Sanitize(xenapi.PIF.RemoveFromOtherConfig(s, from, purposeKey))
}

// managementNameCalls are the calls of moveManagementName (PIF.SaveChanges of other_config).
var managementNameCalls = []string{"pif.add_to_other_config", "pif.remove_from_other_config"}

// CreateBond is CreateBondAction (after BondDetails.CanCreateBond's question): a network named after the bond, the
// bond of these NICs' devices on every server (the coordinator last), then the secondary management interfaces'
// names moved onto the bonds; on a failure the names go back and the bonds and the network are destroyed.
func CreateBond(c *xapi.Connection, spec BondSpec) error {
	s := c.Session
	d, err := readNet(s)
	if err != nil {
		return err
	}
	refs, err := d.members(spec.Members)
	if err != nil {
		return err
	}
	mode := xenapi.BondMode(spec.Mode)
	switch {
	case mode != xenapi.BondModeBalanceSlb && mode != xenapi.BondModeActiveBackup && mode != xenapi.BondModeLacp:
		return fmt.Errorf("unknown bond mode %q", spec.Mode)
	case mode == xenapi.BondModeLacp && (!d.vswitch() || spec.Hashing != "src_mac" && spec.Hashing != "tcpudp_ports"):
		return errors.New("LACP needs the vSwitch and a load balancing method") // Helpers.SupportsLinkAggregationBond
	}
	var devices []string
	var minMTU []int
	for _, r := range refs {
		devices = append(devices, d.pifs[r].Device)
		minMTU = append(minMTU, d.pifs[r].MTU)
	}
	if _, msg, ok := rules.ValidateMTU(strconv.Itoa(spec.MTU), rules.MTUMin, max(rules.MTUMin, rules.BondMaxMTU(minMTU))); !ok {
		return errors.New(msg)
	}
	if w := d.createWarning(refs); w != nil && w.Block {
		return errors.New(w.Text)
	}
	if err := permitted(s, append([]string{"host.management_reconfigure", "network.create", "network.destroy", "network.remove_from_other_config",
		"pif.reconfigure_ip", "pif.plug", "bond.create", "bond.destroy"}, managementNameCalls...)...); err != nil {
		return err
	}
	// NetworkingHelper.PIFsOnAllHosts: the physical PIFs of these devices on each server, by name
	onHost := map[xenapi.HostRef][]xenapi.PIFRef{}
	for r, p := range d.pifs {
		if p.Physical && slices.Contains(devices, p.Device) {
			onHost[p.Host] = append(onHost[p.Host], r)
		}
	}
	hosts := []xenapi.HostRef{} // GetHostsCoordinatorLast
	for h := range d.hosts {
		if h != d.coordinator {
			hosts = append(hosts, h)
		}
	}
	hosts = append(hosts, d.coordinator)

	name := rules.BondName(devices)
	auto := map[bool]string{true: "true", false: "false"}[spec.Auto] // Network.SetAutoPlug
	netRef, err := waitRef(s)(xenapi.Network.AsyncCreate(s, xenapi.NetworkRecord{NameLabel: name, MTU: spec.MTU, Managed: true,
		OtherConfig: map[string]string{"automatic": auto, createInProgress: "true"}}))
	if err != nil {
		return xapi.Sanitize(err)
	}
	network := xenapi.NetworkRef(netRef)
	type newBond struct {
		bond    xenapi.BondRef
		master  xenapi.PIFRef
		members []xenapi.PIFRef
	}
	var made []newBond
	err = func() error {
		if err := xenapi.Network.RemoveFromOtherConfig(s, network, createInProgress); err != nil {
			return err
		}
		props := map[string]string{}
		if mode == xenapi.BondModeLacp {
			props["hashing_algorithm"] = spec.Hashing
		}
		for _, h := range hosts {
			pifs := onHost[h]
			d.sortByName(pifs)
			r, err := waitRef(s)(xenapi.Bond.AsyncCreate(s, network, pifs, "", mode, props))
			if api := (*xenapi.APIError)(nil); err != nil && !errors.As(err, &api) && !errors.Is(err, ErrUncertain) {
				return fmt.Errorf("%w (%v)", ErrUncertain, err) // a lost reply: the bond may be made, so nothing is undone
			}
			if err != nil {
				return err
			}
			b, err := xenapi.Bond.GetRecord(s, xenapi.BondRef(r))
			if err != nil {
				return fmt.Errorf("%w (the bond was made, its interface is unreadable: %v)", ErrUncertain, err)
			}
			made = append(made, newBond{xenapi.BondRef(r), b.Master, pifs})
		}
		for _, b := range made { // ReconfigureManagementInterfaces
			for _, m := range b.members {
				if err := moveManagementName(s, m, b.master); err != nil {
					return err
				}
			}
		}
		return nil
	}()
	if err == nil || errors.Is(err, ErrUncertain) { // an uncertain outcome is never undone
		return xapi.Sanitize(err)
	}
	for _, b := range made { // RevertManagementInterfaces, DestroyBond, DestroyNetwork: best effort, as the C# app
		back := b.members[0]
		if p, err := xenapi.Bond.GetPrimarySlave(s, b.bond); err == nil && valid(p) {
			back = p
		}
		_ = moveManagementName(s, b.master, back)
		_ = xenapi.Bond.Destroy(s, b.bond)
	}
	_ = xenapi.Network.Destroy(s, network)
	return xapi.Sanitize(err)
}

// bondOf is BondRules.BondOf: the bond behind a bond network (from its PIF on the coordinator first).
func (d *netState) bondOf(n xenapi.NetworkRecord) (xenapi.BondRef, bool) {
	var found xenapi.BondRef
	for _, r := range n.PIFs {
		p := d.pifs[r]
		if len(p.BondMasterOf) == 0 {
			continue
		}
		if _, ok := d.bonds[p.BondMasterOf[0]]; ok && (found == "" || p.Host == d.coordinator) {
			found = p.BondMasterOf[0]
		}
	}
	return found, found != ""
}

func (d *netState) network(uuid string) (xenapi.NetworkRef, xenapi.NetworkRecord, error) {
	for r, n := range d.nets {
		if n.UUID == uuid {
			return r, n, nil
		}
	}
	return "", xenapi.NetworkRecord{}, errors.New("the network is gone")
}

// deleteWarning is BondRules.DeleteWarning for a bond network.
func (d *netState) deleteWarning(n xenapi.NetworkRecord) (rules.Warning, error) {
	bond, ok := d.bondOf(n)
	if !ok {
		return rules.Warning{}, errors.New("this network is not a bond")
	}
	var devices []string
	for _, sl := range d.bonds[bond].Slaves {
		devices = append(devices, d.pifs[sl].Device)
	}
	var primary, secondary bool
	for _, r := range n.PIFs {
		primary = primary || d.pifs[r].Management
		secondary = secondary || d.secondary(d.pifs[r], true)
	}
	return rules.BondDeleteWarning(rules.BondName(devices), primary, secondary, d.pool.HaEnabled, d.poolName()), nil
}

// BondDeleteWarning is the question Delete Bond asks first (DestroyBondCommand).
func BondDeleteWarning(c *xapi.Connection, networkUUID string) (rules.Warning, error) {
	d, err := readNet(c.Session)
	if err != nil {
		return rules.Warning{}, err
	}
	_, n, err := d.network(networkUUID)
	if err != nil {
		return rules.Warning{}, err
	}
	return d.deleteWarning(n)
}

// DestroyBond is DestroyBondAction (after its question): on every server the bond matching this one (the same
// devices) gives its management interface name back to its primary member, then the bonds and the network go.
func DestroyBond(c *xapi.Connection, networkUUID string) error {
	s := c.Session
	d, err := readNet(s)
	if err != nil {
		return err
	}
	netRef, n, err := d.network(networkUUID)
	if err != nil {
		return err
	}
	w, err := d.deleteWarning(n)
	if err != nil {
		return err
	}
	if ok, _ := rules.CanRemoveNetwork(d.info(n, false)); !ok || w.Block {
		return errors.New(map[bool]string{true: w.Text, false: "the bond is busy"}[w.Block])
	}
	if err := permitted(s, append([]string{"host.management_reconfigure", "network.destroy", "vif.plug", "vif.unplug",
		"pif.reconfigure_ip", "pif.plug", "bond.destroy"}, managementNameCalls...)...); err != nil {
		return err
	}
	bond, _ := d.bondOf(n)
	devices := func(b xenapi.BondRecord) []string {
		var out []string
		for _, sl := range b.Slaves {
			out = append(out, d.pifs[sl].Device)
		}
		return out
	}
	want := devices(d.bonds[bond])
	matches := map[xenapi.HostRef][]xenapi.BondRef{}
	for r, b := range d.bonds { // NetworkingHelper.FindBond: on each server, a bond whose members are a subset of the other's
		got := devices(b)
		small, big := got, want
		if len(small) > len(big) {
			small, big = big, small
		}
		if valid(b.Master) && len(small) > 0 && all(small, func(dev string) bool { return slices.Contains(big, dev) }) {
			matches[d.pifs[b.Master].Host] = append(matches[d.pifs[b.Master].Host], r)
		}
	}
	perHost := map[xenapi.HostRef]xenapi.BondRef{}
	for h, rs := range matches {
		switch {
		case slices.Contains(rs, bond):
			perHost[h] = bond
		case len(rs) > 1: // FindBond takes the first: never guess which bond to destroy
			return fmt.Errorf("%s has more than one bond like this one: delete them from its own NICs tab", d.hosts[h].NameLabel)
		default:
			perHost[h] = rs[0]
		}
	}
	var bonds []xenapi.BondRef
	var firstErr error
	for _, r := range perHost {
		b := d.bonds[r]
		bonds = append(bonds, r)
		first := b.PrimarySlave // the first member when the primary one is unknown
		if !valid(first) && len(b.Slaves) > 0 {
			sl := slices.Clone(b.Slaves)
			d.sortByName(sl)
			first = sl[0]
		}
		if valid(first) {
			if err := moveManagementName(s, b.Master, first); err != nil {
				return err
			}
		}
	}
	for _, b := range bonds {
		if err := wait(s)(xenapi.Bond.AsyncDestroy(s, b)); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if err := xenapi.Network.Destroy(s, netRef); err != nil && firstErr == nil {
		firstErr = err
	}
	return xapi.Sanitize(firstErr)
}

func all[T any](list []T, f func(T) bool) bool {
	for _, v := range list {
		if !f(v) {
			return false
		}
	}
	return true
}

// RescanPIFs is RescanPIFsAction: PIF.scan on the server.
func RescanPIFs(c *xapi.Connection, hostUUID string) error {
	s := c.Session
	h, err := xenapi.Host.GetByUUID(s, hostUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if err := permitted(s, "pif.scan"); err != nil {
		return err
	}
	return xapi.Sanitize(xenapi.PIF.Scan(s, h))
}

// ---------- networks ----------

// NewNetwork is the New Network wizard's External or Single-Server Private network (NetWName, NetWDetails).
type NewNetwork struct {
	External    bool   `json:"external"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Auto        bool   `json:"auto"`
	Host        string `json:"host"` // the wizard's server (NetworkRules.WizardHost: the server, or the coordinator)
	NIC         string `json:"nic"`  // external: the PIF uuid of the NIC on that server
	VLAN        string `json:"vlan"`
	MTU         string `json:"mtu"`
	Sriov       bool   `json:"sriov"` // external: on the NIC's SR-IOV network (NetWDetails.CreateVlanOnSriovNetwork)
}

// newNetwork is NetWDetails' checks and NewNetworkWizard.CreateNonBonded's choice on this state: the network, the PIF
// its VLAN is made on (the NIC, or the logical PIF of its SR-IOV network; "" for a private network) and the calls.
func (d *netState) newNetwork(spec NewNetwork) (rec xenapi.NetworkRecord, nic xenapi.PIFRef, vlan int, calls []string, err error) {
	rec = xenapi.NetworkRecord{NameLabel: strings.TrimSpace(spec.Name), NameDescription: spec.Description, Managed: true, MTU: rules.MTUDefault,
		OtherConfig: map[string]string{"automatic": map[bool]string{true: "true", false: "false"}[spec.Auto]}}
	if rec.NameLabel == "" {
		return rec, "", 0, nil, errors.New("enter a name for the new network")
	}
	if !spec.External {
		return rec, "", 0, []string{"network.create"}, nil
	}
	for r, p := range d.pifs { // NetworkRules.WizardNics: the server's shown physical NICs, not bond members
		if p.UUID == spec.NIC && d.hosts[p.Host].UUID == spec.Host && physical(p) && p.Managed && p.OtherConfig["HideFromXenCenter"] != "true" && !valid(p.BondSlaveOf) {
			nic = r
		}
	}
	if nic == "" {
		return rec, "", 0, nil, errors.New("select a NIC")
	}
	p := d.pifs[nic]
	// the wizard sends Sriov only when checkBoxSriov is shown (PIF.IsSriovPhysicalPIF) and ticked; CreateNonBonded
	// stops when the NIC has no SR-IOV network any more
	if spec.Sriov && len(p.SriovPhysicalPIFOf) == 0 {
		return rec, "", 0, nil, errors.New("this NIC has no SR-IOV network")
	}
	var msg string
	var ok bool
	if vlan, msg, ok = rules.ValidateVLAN(spec.VLAN, rules.MinVLAN(d.vswitch()), d.vlansInUse(p.Device, spec.Sriov)); !ok {
		return rec, "", 0, nil, errors.New(msg)
	}
	if rec.MTU, msg, ok = rules.ValidateMTU(spec.MTU, rules.MTUMin, max(rules.MTUMin, rules.NICMaxMTU(p.MTU))); !ok {
		return rec, "", 0, nil, errors.New(msg)
	}
	if spec.Sriov {
		if nic = d.sriovLogicalPIF(p); nic == "" {
			return rec, "", 0, nil, errors.New("the SR-IOV network of this NIC is gone")
		}
	}
	return rec, nic, vlan, []string{"network.create", "pool.create_VLAN_from_PIF"}, nil
}

// CreateNetwork is NetworkAction's create (NewNetworkWizard.Finish): network.create, then for an external network
// pool.create_VLAN_from_PIF on the NIC, or on the logical PIF of its SR-IOV network.
func CreateNetwork(c *xapi.Connection, spec NewNetwork) error {
	s := c.Session
	d, err := readNet(s)
	if err != nil {
		return err
	}
	rec, nic, vlan, calls, err := d.newNetwork(spec)
	if err != nil {
		return err
	}
	if err := permitted(s, calls...); err != nil {
		return err
	}
	ref, err := xenapi.Network.Create(s, rec)
	if err != nil || nic == "" {
		return xapi.Sanitize(err)
	}
	if _, err := xenapi.Pool.CreateVLANFromPIF(s, nic, ref, vlan); err != nil {
		_ = xenapi.Network.Destroy(s, ref) // not in NetworkAction: an empty network is no use to anyone
		return xapi.Sanitize(err)
	}
	return nil
}

// RemoveNetwork is NetworkAction's destroy (after MESSAGEBOX_NETWORK_DELETE): its VLAN, tunnel or SR-IOV PIFs go
// first. A bond network is removed with DestroyBond.
func RemoveNetwork(c *xapi.Connection, uuid string) error {
	s := c.Session
	d, err := readNet(s)
	if err != nil {
		return err
	}
	ref, n, err := d.network(uuid)
	if err != nil {
		return err
	}
	info := d.info(n, false)
	switch ok, tip := rules.CanRemoveNetwork(info); {
	case info.Bond:
		return errors.New("a bond network is removed by deleting its bond")
	case tip != "":
		return errors.New(tip)
	case !ok:
		return errors.New("the network cannot be removed now")
	}
	if err := permitted(s, append([]string{"network.destroy"}, d.destroyCalls(n.PIFs)...)...); err != nil {
		return err
	}
	if err := d.destroyPIFs(n.PIFs); err != nil {
		return err
	}
	return xapi.Sanitize(xenapi.Network.Destroy(s, ref))
}

// destroyCalls are the calls of destroyPIFs for these PIFs.
func (d *netState) destroyCalls(pifs []xenapi.PIFRef) []string {
	var calls []string
	for _, r := range pifs {
		switch p := d.pifs[r]; {
		case len(p.TunnelAccessPIFOf) > 0:
			calls = append(calls, "tunnel.destroy")
		case p.Physical:
			calls = append(calls, "pif.forget")
		case p.VLAN != -1:
			calls = append(calls, "vlan.destroy")
		case len(p.SriovLogicalPIFOf) > 0:
			calls = append(calls, "network_sriov.destroy")
		}
	}
	return calls
}

// destroyPIFs is NetworkAction.destroyPIFs: a tunnel access PIF goes with its tunnel, a VLAN with its VLAN, an SR-IOV
// PIF with its SR-IOV network, a physical one is forgotten.
func (d *netState) destroyPIFs(pifs []xenapi.PIFRef) error {
	s := d.s
	for _, r := range pifs {
		var err error
		switch p := d.pifs[r]; {
		case len(p.TunnelAccessPIFOf) > 0:
			for _, t := range p.TunnelAccessPIFOf {
				if err = xenapi.Tunnel.Destroy(s, t); err != nil {
					break
				}
			}
		case p.Physical:
			err = xenapi.PIF.Forget(s, r)
		case p.VLAN != -1:
			err = xenapi.VLAN.Destroy(s, p.VLANMasterOf)
		case len(p.SriovLogicalPIFOf) > 0:
			err = xenapi.NetworkSriov.Destroy(s, p.SriovLogicalPIFOf[0])
		}
		if err != nil {
			return xapi.Sanitize(err)
		}
	}
	return nil
}

// NetworkChanges are the network's Properties (GeneralEditPage's name and description; EditNetworkPage's auto-add,
// MTU, NIC and VLAN, bond mode and hashing algorithm).
type NetworkChanges struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Auto        bool   `json:"auto"`
	MTU         string `json:"mtu"`
	NIC         string `json:"nic"`     // the PIF uuid of a NIC of the coordinator; "" <internal> (a private network)
	VLAN        string `json:"vlan"`    // with a NIC
	Mode        string `json:"mode"`    // a bond network's: balance-slb, active-backup or lacp
	Hashing     string `json:"hashing"` // lacp: src_mac or tcpudp_ports
}

// SettingsNIC is a NIC of the NIC list (populateHostNicList) with the VLANs other networks use on its device.
type SettingsNIC struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	VLANs []int  `json:"vlans"`
}

// NetworkSettings is the Network Settings page of a network (EditNetworkPage) but its auto-add.
type NetworkSettings struct {
	MTUState rules.MTUState `json:"mtu_state"`
	VLANMTU  rules.MTUState `json:"vlan_mtu_state"` // the MTU of a new VLAN (GetSelectedPIF null): fixed by running VMs without tools
	NICState rules.NICState `json:"nic_state"`      // the NIC and VLAN fields
	NICs     []SettingsNIC  `json:"nics"`
	NIC      string         `json:"nic"` // FindAssociatedPhysicalPIF: the NIC under the network; "" <internal>
	VLAN     int            `json:"vlan"`
	MinVLAN  int            `json:"min_vlan"` // Helpers.VLAN0Allowed
	Bond     bool           `json:"bond"`     // groupBoxBondMode: a bond network
	LACP     bool           `json:"lacp"`     // the LACP choices (Helpers.SupportsLinkAggregationBond)
	Mode     string         `json:"mode"`     // the bond's (NetworkBondMode)
	Hashing  string         `json:"hashing"`  // the bond's (HashingAlgorithm; "" unknown)
}

// ReadNetworkSettings reads a network's Network Settings page. noTools tells a running VM without I/O drivers
// (VM.GetVirtualizationStatus).
func ReadNetworkSettings(s *xenapi.Session, uuid string, noTools func(xenapi.VMRecord) bool) (NetworkSettings, error) {
	d, err := readNet(s)
	if err != nil {
		return NetworkSettings{}, err
	}
	_, n, err := d.network(uuid)
	if err != nil {
		return NetworkSettings{}, err
	}
	st, _ := d.settings(n, noTools)
	return st, nil
}

// managementOf is the network's management interface (blockDueToManagement), its PIF on the coordinator first.
func (d *netState) managementOf(n xenapi.NetworkRecord) (mgmt bool, name string) {
	for _, r := range n.PIFs {
		if p := d.pifs[r]; d.management(p) && (!mgmt || p.Host == d.coordinator) {
			mgmt, name = true, managementName(p)
		}
	}
	return mgmt, name
}

// settings is EditNetworkPage.Repopulate and its enablement, with the coordinator's PIFs by uuid for the NIC list.
func (d *netState) settings(n xenapi.NetworkRecord, noTools func(xenapi.VMRecord) bool) (NetworkSettings, map[string]xenapi.PIFRef) {
	st := NetworkSettings{NICs: []SettingsNIC{}, VLAN: 1, MinVLAN: rules.MinVLAN(d.vswitch()), LACP: d.vswitch()}
	var tunnel, attached, vmsNoTools bool
	var pif *xenapi.PIFRecord // GetNetworksPIF: the network's PIF on the coordinator
	for _, r := range n.PIFs {
		p := d.pifs[r]
		tunnel = tunnel || len(p.TunnelAccessPIFOf) > 0
		if p.Host == d.coordinator {
			pif = &p
		}
	}
	for _, v := range n.VIFs {
		vif, err := xenapi.VIF.GetRecord(d.s, v)
		var rec xenapi.VMRecord
		if err == nil {
			rec, err = xenapi.VM.GetRecord(d.s, vif.VM)
		}
		attached = attached || err != nil || vif.CurrentlyAttached                                            // unknown: as attached
		vmsNoTools = vmsNoTools || err != nil || rec.PowerState == xenapi.VMPowerStateRunning && noTools(rec) // and as without tools
	}
	mgmt, mgmtName := d.managementOf(n)
	// GetSelectedPIF looks at the coordinator's PIF only; every PIF counts here, as the new MTU replugs them all
	st.MTUState = rules.NetworkMTU(n.MTU, tunnel, len(n.PIFs) == 0, mgmt, mgmtName, vmsNoTools)
	st.VLANMTU = rules.NetworkMTU(n.MTU, false, false, false, "", vmsNoTools)

	editable := pif != nil && !physical(*pif) && len(pif.TunnelAccessPIFOf) == 0 && len(pif.SriovLogicalPIFOf) == 0 // Editable(pif)
	nics := map[string]xenapi.PIFRef{}
	for r, p := range d.pifs { // populateHostNicList: the coordinator's shown physical PIFs, not bond members
		if p.Host != d.coordinator || !physical(p) || !p.Managed || p.OtherConfig["HideFromXenCenter"] == "true" || valid(p.BondSlaveOf) {
			continue
		}
		nic := SettingsNIC{UUID: p.UUID, Name: d.nicName(p), VLANs: []int{}}
		for _, o := range d.pifs { // NetworkAction's PIF_VLAN_EXISTS check, this network's own VLANs too
			if o.Device == p.Device && o.VLAN >= 0 {
				nic.VLANs = append(nic.VLANs, o.VLAN)
			}
		}
		if editable && p.Device == pif.Device { // FindAssociatedPhysicalPIF
			st.NIC, st.VLAN = p.UUID, pif.VLAN
		}
		nics[p.UUID] = r
		st.NICs = append(st.NICs, nic)
	}
	slices.SortFunc(st.NICs, func(a, b SettingsNIC) int { return strings.Compare(a.Name, b.Name) })
	// a network with PIFs on other servers only, or a VLAN whose NIC is gone, shows no NIC: <internal> would be
	// picked for it, and OK would make it private (EditNetworkPage does)
	st.NICState = rules.NetworkNICs(len(n.PIFs) == 0, editable && st.NIC != "", mgmt, mgmtName, attached)

	if bond, ok := d.bondOf(n); ok {
		b := d.bonds[bond]
		st.Bond, st.Mode = true, string(b.Mode)
		if h := b.Properties["hashing_algorithm"]; h == "src_mac" || h == "tcpudp_ports" { // Bond.HashingAlgoritm
			st.Hashing = h
		}
	}
	return st, nics
}

// nicName is PIF.Name of a physical PIF: "NIC 1", or "Bond 1+2" for a bond interface.
func (d *netState) nicName(p xenapi.PIFRecord) string {
	if len(p.BondMasterOf) == 0 {
		return messages.Format(messages.Get("PIF_NIC"), strings.TrimPrefix(p.Device, "eth"))
	}
	var devices []string
	for _, sl := range d.bonds[p.BondMasterOf[0]].Slaves {
		devices = append(devices, d.pifs[sl].Device)
	}
	return rules.BondName(devices)
}

// UpdateNetwork is EditNetworkPage.SaveSettings (and GeneralEditPage's name and description): name, description and
// auto-add (Network.SetAutoPlug); a new MTU; a bond network's mode and hashing algorithm on each of its bonds
// (SetBondModeActions, SetBondPropertiesActions); a new NIC or VLAN, or <internal>, destroys the network's PIFs and
// makes the VLAN again on the pool (NetworkAction's update); without one a new MTU unplugs and plugs the network's PIFs
// and running VMs' VIFs (UnplugPlugNetworkAction). Everything is checked before the first write.
func UpdateNetwork(c *xapi.Connection, uuid string, ch NetworkChanges, noTools func(xenapi.VMRecord) bool) error {
	s := c.Session
	d, err := readNet(s)
	if err != nil {
		return err
	}
	ref, n, err := d.network(uuid)
	if err != nil {
		return err
	}
	if !rules.CanEditNetwork(d.info(n, false)) {
		return errors.New("the network cannot be changed now")
	}
	name := strings.TrimSpace(ch.Name)
	if name == "" {
		return errors.New("enter a name for the network")
	}
	st, nics := d.settings(n, noTools)

	// InternalToExternal, ExternalToInternal, ExternalChangedDeviceOrVlan
	var nic xenapi.PIFRef
	vlan, newPIFs := -1, false
	if st.NICState.Shown {
		if ch.NIC != "" {
			i := slices.IndexFunc(st.NICs, func(x SettingsNIC) bool { return x.UUID == ch.NIC })
			if i < 0 {
				return errors.New("select a NIC")
			}
			if v, err := strconv.Atoi(strings.TrimSpace(ch.VLAN)); err == nil && ch.NIC == st.NIC && v == st.VLAN {
				vlan = v // unchanged: not checked, as NetworkAction checks only a new VLAN
			} else {
				var msg string
				var ok bool
				if vlan, msg, ok = rules.ValidateVLAN(ch.VLAN, st.MinVLAN, st.NICs[i].VLANs); !ok {
					return errors.New(msg)
				}
			}
			nic = nics[ch.NIC]
		}
		newPIFs = ch.NIC != st.NIC || ch.NIC != "" && vlan != st.VLAN
		if newPIFs && !st.NICState.Enabled {
			return errors.New(strings.Trim(st.NICState.Note, "()"))
		}
	}

	mtu, m := n.MTU, st.MTUState
	if newPIFs { // SetMTUControlEnablement for the new choice: fixed on <internal>, a new VLAN is no management interface
		m = st.VLANMTU
		if nic == "" {
			m = rules.NetworkMTU(n.MTU, false, true, false, "", false)
		}
	}
	if m.Shown && m.Enabled {
		var msg string
		var ok bool
		if mtu, msg, ok = rules.ValidateMTU(ch.MTU, m.Min, m.Max); !ok {
			return errors.New(msg)
		}
	}

	var setMode, setHashing bool
	var bonds []xenapi.BondRef // Network.TheBonds
	if st.Bond {
		if setMode, setHashing, err = rules.BondChanges(st.Mode, st.Hashing, ch.Mode, ch.Hashing, st.LACP); err != nil {
			return err
		}
		for _, r := range n.PIFs {
			if p := d.pifs[r]; len(p.BondMasterOf) > 0 {
				bonds = append(bonds, p.BondMasterOf[0])
			}
		}
	}

	auto := n.OtherConfig["automatic"] != "false" // Network.GetAutoPlug
	var calls []string
	if name != n.NameLabel {
		calls = append(calls, "network.set_name_label")
	}
	if ch.Description != n.NameDescription {
		calls = append(calls, "network.set_name_description")
	}
	if ch.Auto != auto {
		calls = append(calls, "network.add_to_other_config", "network.remove_from_other_config")
	}
	if mtu != n.MTU {
		calls = append(calls, "network.set_MTU")
		if !newPIFs {
			calls = append(calls, "pif.unplug", "pif.plug", "vif.unplug", "vif.plug")
		}
	}
	if setMode {
		calls = append(calls, "bond.set_mode")
	}
	if setHashing {
		calls = append(calls, "bond.set_property")
	}
	if newPIFs {
		calls = append(calls, d.destroyCalls(n.PIFs)...)
		if nic != "" {
			calls = append(calls, "pool.create_VLAN_from_PIF")
		}
	}
	if len(calls) == 0 {
		return nil
	}
	if err := permitted(s, calls...); err != nil {
		return err
	}
	if name != n.NameLabel {
		if err := xenapi.Network.SetNameLabel(s, ref, name); err != nil {
			return xapi.Sanitize(err)
		}
	}
	if ch.Description != n.NameDescription {
		if err := xenapi.Network.SetNameDescription(s, ref, ch.Description); err != nil {
			return xapi.Sanitize(err)
		}
	}
	if ch.Auto != auto {
		_ = xenapi.Network.RemoveFromOtherConfig(s, ref, "automatic")
		if err := xenapi.Network.AddToOtherConfig(s, ref, "automatic", strconv.FormatBool(ch.Auto)); err != nil {
			return xapi.Sanitize(err)
		}
	}
	var pifs []xenapi.PIFRef
	var vifs []xenapi.VIFRef
	if mtu != n.MTU && !newPIFs {
		if pifs, vifs, err = attached(s, n); err != nil {
			return err
		}
	}
	if mtu != n.MTU {
		if err := xenapi.Network.SetMTU(s, ref, mtu); err != nil {
			return xapi.Sanitize(err)
		}
	}
	for _, b := range bonds {
		if setMode {
			if err := xenapi.Bond.SetMode(s, b, xenapi.BondMode(ch.Mode)); err != nil {
				return xapi.Sanitize(err)
			}
		}
	}
	for _, b := range bonds {
		if setHashing {
			if err := xenapi.Bond.SetProperty(s, b, "hashing_algorithm", ch.Hashing); err != nil {
				return xapi.Sanitize(err)
			}
		}
	}
	if newPIFs { // the new PIFs come with the network's MTU, nothing to replug
		if err := d.destroyPIFs(n.PIFs); err != nil {
			return err
		}
		if nic != "" {
			_, err := xenapi.Pool.CreateVLANFromPIF(s, nic, ref, vlan)
			return xapi.Sanitize(err)
		}
		return nil
	}
	if mtu == n.MTU {
		return nil
	}
	return unplugPlug(s, pifs, vifs)
}

// attached are the network's attached PIFs and its running VMs' attached VIFs (UnplugPlugNetworkAction's lists).
func attached(s *xenapi.Session, n xenapi.NetworkRecord) (pifs []xenapi.PIFRef, vifs []xenapi.VIFRef, err error) {
	for _, r := range n.PIFs {
		a, err := xenapi.PIF.GetCurrentlyAttached(s, r)
		if err != nil {
			return nil, nil, xapi.Sanitize(err)
		}
		if a {
			pifs = append(pifs, r)
		}
	}
	for _, r := range n.VIFs {
		v, err := xenapi.VIF.GetRecord(s, r)
		if err != nil {
			return nil, nil, xapi.Sanitize(err)
		}
		if !v.CurrentlyAttached {
			continue
		}
		st, err := xenapi.VM.GetPowerState(s, v.VM)
		if err != nil {
			return nil, nil, xapi.Sanitize(err)
		}
		if st == xenapi.VMPowerStateRunning {
			vifs = append(vifs, r)
		}
	}
	return pifs, vifs, nil
}

// unplugPlug is UnplugPlugNetworkAction: the VIFs and PIFs are unplugged and plugged again, carrying on past failures
// (to leave as much plugged as it can) and returning the first one.
func unplugPlug(s *xenapi.Session, pifs []xenapi.PIFRef, vifs []xenapi.VIFRef) error {
	var first error
	keep := func(err error) {
		if err != nil && first == nil {
			first = err
		}
	}
	for _, v := range vifs {
		keep(xenapi.VIF.Unplug(s, v))
	}
	for _, p := range pifs {
		keep(xenapi.PIF.Unplug(s, p))
	}
	for _, p := range pifs {
		keep(xenapi.PIF.Plug(s, p))
	}
	for _, v := range vifs {
		keep(xenapi.VIF.Plug(s, v))
	}
	return xapi.Sanitize(first)
}
