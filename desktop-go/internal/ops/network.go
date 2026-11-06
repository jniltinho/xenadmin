package ops

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	xenapi "go/xenapi"

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

// vlansInUse is NetworkRules.VlansInUse: the VLANs already used on a device (SR-IOV VLANs apart).
func (d *netState) vlansInUse(device string) []int {
	out := []int{}
	for _, p := range d.pifs {
		if p.Device == device && !d.sriov(p) {
			out = append(out, p.VLAN)
		}
	}
	return out
}

// VLANsInUse are the VLANs already used on each device, for the New Network wizard.
func VLANsInUse(s *xenapi.Session) (map[string][]int, error) {
	d, err := readNet(s)
	if err != nil {
		return nil, err
	}
	out := map[string][]int{}
	for _, p := range d.pifs {
		out[p.Device] = d.vlansInUse(p.Device)
	}
	return out, nil
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
}

// CreateNetwork is NetworkAction's create (NewNetworkWizard.Finish): network.create, then for an external network
// pool.create_VLAN_from_PIF on the NIC.
func CreateNetwork(c *xapi.Connection, spec NewNetwork) error {
	s := c.Session
	d, err := readNet(s)
	if err != nil {
		return err
	}
	rec := xenapi.NetworkRecord{NameLabel: strings.TrimSpace(spec.Name), NameDescription: spec.Description, Managed: true, MTU: rules.MTUDefault,
		OtherConfig: map[string]string{"automatic": map[bool]string{true: "true", false: "false"}[spec.Auto]}}
	if rec.NameLabel == "" {
		return errors.New("enter a name for the new network")
	}
	calls := []string{"network.create"}
	var nic xenapi.PIFRef
	vlan := 0
	if spec.External {
		for r, p := range d.pifs { // NetworkRules.WizardNics: the server's shown physical NICs, not bond members
			if p.UUID == spec.NIC && d.hosts[p.Host].UUID == spec.Host && physical(p) && p.Managed && p.OtherConfig["HideFromXenCenter"] != "true" && !valid(p.BondSlaveOf) {
				nic = r
			}
		}
		if nic == "" {
			return errors.New("select a NIC")
		}
		var msg string
		var ok bool
		if vlan, msg, ok = rules.ValidateVLAN(spec.VLAN, rules.MinVLAN(d.vswitch()), d.vlansInUse(d.pifs[nic].Device)); !ok {
			return errors.New(msg)
		}
		if rec.MTU, msg, ok = rules.ValidateMTU(spec.MTU, rules.MTUMin, max(rules.MTUMin, rules.NICMaxMTU(d.pifs[nic].MTU))); !ok {
			return errors.New(msg)
		}
		calls = append(calls, "pool.create_VLAN_from_PIF")
	}
	if err := permitted(s, calls...); err != nil {
		return err
	}
	ref, err := xenapi.Network.Create(s, rec)
	if err != nil || !spec.External {
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
	calls := []string{"network.destroy"}
	for _, r := range n.PIFs {
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
	if err := permitted(s, calls...); err != nil {
		return err
	}
	for _, r := range n.PIFs { // NetworkAction.destroyPIFs
		p := d.pifs[r]
		switch {
		case len(p.TunnelAccessPIFOf) > 0:
			for _, t := range p.TunnelAccessPIFOf {
				err = xenapi.Tunnel.Destroy(s, t)
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
	return xapi.Sanitize(xenapi.Network.Destroy(s, ref))
}

// NetworkChanges are the network's Properties (GeneralEditPage's name and description, EditNetworkPage's auto-add and MTU).
type NetworkChanges struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Auto        bool   `json:"auto"`
	MTU         string `json:"mtu"`
}

// NetworkMTU is EditNetworkPage.SetMTUControlEnablement for a network: whether its MTU is shown and can be changed.
// noTools tells a running VM without I/O drivers (VM.GetVirtualizationStatus).
func NetworkMTU(s *xenapi.Session, uuid string, noTools func(xenapi.VMRecord) bool) (rules.MTUState, error) {
	d, err := readNet(s)
	if err != nil {
		return rules.MTUState{}, err
	}
	_, n, err := d.network(uuid)
	if err != nil {
		return rules.MTUState{}, err
	}
	return d.mtu(n, noTools), nil
}

func (d *netState) mtu(n xenapi.NetworkRecord, noTools func(xenapi.VMRecord) bool) rules.MTUState {
	var tunnel, mgmt bool
	var mgmtName string
	for _, r := range n.PIFs {
		p := d.pifs[r]
		tunnel = tunnel || len(p.TunnelAccessPIFOf) > 0
		// GetSelectedPIF looks at the coordinator's PIF only; every PIF counts here, as the new MTU replugs them all
		if d.management(p) && (!mgmt || p.Host == d.coordinator) {
			mgmt, mgmtName = true, managementName(p)
		}
	}
	vmsNoTools := false
	for _, v := range n.VIFs {
		vm, err := xenapi.VIF.GetVM(d.s, v)
		rec, err2 := xenapi.VM.GetRecord(d.s, vm)
		if err != nil || err2 != nil || rec.PowerState == xenapi.VMPowerStateRunning && noTools(rec) { // unknown: as without tools
			vmsNoTools = true
		}
	}
	return rules.NetworkMTU(n.MTU, tunnel, len(n.PIFs) == 0, mgmt, mgmtName, vmsNoTools)
}

// UpdateNetwork saves the network's Properties: name, description and auto-add (Network.SetAutoPlug); a new MTU is
// set and the network's PIFs and running VMs' VIFs are unplugged and plugged again (UnplugPlugNetworkAction).
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
	mtu := n.MTU
	if m := d.mtu(n, noTools); m.Shown && m.Enabled {
		var msg string
		var ok bool
		if mtu, msg, ok = rules.ValidateMTU(ch.MTU, m.Min, m.Max); !ok {
			return errors.New(msg)
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
		calls = append(calls, "network.set_MTU", "pif.unplug", "pif.plug", "vif.unplug", "vif.plug")
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
	if mtu == n.MTU {
		return nil
	}
	pifs, vifs, err := attached(s, n)
	if err != nil {
		return err
	}
	if err := xenapi.Network.SetMTU(s, ref, mtu); err != nil {
		return xapi.Sanitize(err)
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
