package app

import (
	"encoding/xml"
	"sort"
	"strings"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// ViewNode is a node of an Organization View or a saved search (MainWindowTreeBuilder.CreateNewRootNode): a group
// with its icon, or an object of a server (its icon is the window's, by kind and state), with children. Key is a
// group's identity: its value, or the server and uuid of the object it is (a network, a vApp, a pool, a server).
type ViewNode struct {
	Name     string      `json:"name"`
	Key      string      `json:"key,omitempty"`
	Icon     string      `json:"icon,omitempty"`
	Server   string      `json:"server,omitempty"`
	Object   *Object     `json:"object,omitempty"`
	Children []*ViewNode `json:"children"`
}

// ViewInfo names a view: the Organization Views (OrganizationalView) and the default searches (Search.InitDefaultSearches).
type ViewInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Supported bool   `json:"supported"` // a saved search this port can run (app.SearchSupported); the others are listed only
}

// OrganizationViews are the views of the Organization Views button, in the C# app's order.
func OrganizationViews() []ViewInfo {
	return []ViewInfo{
		{ID: "tags", Name: messages.Get("TAGS"), Icon: "000_Tag_h32bit_16", Supported: true},
		{ID: "folders", Name: messages.Get("FOLDERS"), Icon: "000_Folder_open_h32bit_16", Supported: true},
		{ID: "custom_fields", Name: messages.Get("CUSTOM_FIELDS"), Icon: "000_Fields_h32bit_16", Supported: true},
		{ID: "vapps", Name: messages.Get("VM_APPLIANCES"), Icon: "000_VirtualAppliance_h32bit_16", Supported: true},
	}
}

// DefaultSearches are Search.InitDefaultSearches, by name as the C# list (NavigationViews.SavedSearches).
func DefaultSearches() []ViewInfo {
	out := []ViewInfo{
		{ID: "vms_by_os", Name: messages.Get("DEFAULT_SEARCH_VMS_BY_OS")},
		{ID: "vms_without_tools", Name: messages.Format(messages.Get("DEFAULT_SEARCH_VMS_WO_XS_TOOLS"), "XCP-ng VM Tools")},
		{ID: "vms_by_power_state", Name: messages.Get("DEFAULT_SEARCH_VMS_BY_POWERSTATE")},
		{ID: "vms_by_network", Name: messages.Get("DEFAULT_SEARCH_VMS_BY_NETWORK")},
		{ID: "vms_by_vapp", Name: messages.Get("DEFAULT_SEARCH_VMS_BY_APPLIANCE")},
		{ID: "objects_by_tag", Name: messages.Get("DEFAULT_SEARCH_OBJECTS_BY_TAG")},
		{ID: "snapshots_by_vm", Name: messages.Get("DEFAULT_SEARCH_SNAPSHOTS_BY_VM")},
	}
	for i := range out {
		out[i].Icon, out[i].Supported = "000_defaultSpyglass_h32bit_16", true
	}
	sort.SliceStable(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	return out
}

// searchPrefix is Search.SearchPrefix (rules.SearchKeyPrefix).
const searchPrefix = rules.SearchKeyPrefix

// PoolSearches are the searches saved in the pool's gui_config (SearchMarshalling XML), by their name.
func PoolSearches(c *xapi.Connection) ([]ViewInfo, error) {
	pools, err := xenapi.Pool.GetAllRecords(c.Session)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	var out []ViewInfo
	for _, p := range pools {
		for k, v := range p.GuiConfig {
			if !strings.HasPrefix(k, searchPrefix) {
				continue
			}
			var s struct {
				Name string `xml:"name,attr"`
			}
			if xml.Unmarshal([]byte(v), &s) != nil || s.Name == "" {
				continue
			}
			out = append(out, ViewInfo{ID: "saved:" + strings.TrimPrefix(k, searchPrefix), Name: s.Name, Icon: "000_Search_h32bit_16",
				Supported: SearchSupported(v) == nil})
		}
	}
	return out, nil
}

// PoolSearch is the XML of a search saved in this server's pool ("" when it is not there).
func PoolSearch(c *xapi.Connection, uuid string) (string, error) {
	pools, err := xenapi.Pool.GetAllRecords(c.Session)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	for _, p := range pools {
		if v, ok := p.GuiConfig[searchPrefix+uuid]; ok {
			return v, nil
		}
	}
	return "", nil
}

// group is one level of a view's path: a group node's name, icon and identity ("" groups by name).
type group struct{ name, icon, key string }

// item is an object of the view with the fields the groupings read.
type item struct {
	obj   Object
	tags  []string
	other map[string]string
}

// inventory is what the views read of a server, each record set once.
type inventory struct {
	pools map[xenapi.PoolRef]xenapi.PoolRecord
	hosts map[xenapi.HostRef]xenapi.HostRecord
	vms   map[xenapi.VMRef]xenapi.VMRecord // the VMs shown (VM.Show)
	all   map[xenapi.VMRef]xenapi.VMRecord // every VM record, to follow references (a snapshot's hidden VM)
	srs   map[xenapi.SRRef]xenapi.SRRecord
	vdis  map[xenapi.VDIRef]xenapi.VDIRecord
	items []item
	real  map[string]Object // the real VMs (IsRealVm) by uuid, with the window's state
}

// ReadView is one server's part of a view: the objects it holds grouped as the C# groupings do (tags, folder path,
// custom field and value, vApp; the default searches' properties).
func ReadView(c *xapi.Connection, server, view string) ([]*ViewNode, error) {
	inv, err := readInventory(c, false)
	if err != nil {
		return nil, err
	}
	root := &ViewNode{}
	switch view {
	case "tags", "objects_by_tag":
		for _, it := range inv.items {
			for _, t := range it.tags {
				add(root, []group{{t, "000_Tag_h32bit_16", ""}}, it.obj, server)
			}
		}
	case "folders":
		for _, p := range inv.pools { // Folders.EMPTY_FOLDERS: folders that hold nothing yet
			for _, f := range strings.Split(p.OtherConfig["EMPTY_FOLDERS"], ";") {
				addGroups(root, folderPath(f))
			}
		}
		for _, it := range inv.items {
			if path := folderPath(it.other["folder"]); len(path) > 0 {
				add(root, path, it.obj, server)
			}
		}
	case "custom_fields":
		for _, f := range customFields(inv.pools) {
			for _, it := range inv.items {
				if v, key := customFieldValue(f, it.other["XenCenter.CustomFields."+f.name]); v != "" {
					if key == "" {
						key = "cf:" + f.name + "=" + v
					}
					add(root, []group{{f.name, "000_Fields_h32bit_16", "cf:" + f.name}, {v, "000_Fields_h32bit_16", key}}, it.obj, server)
				}
			}
		}
	default:
		if err := vmView(c, inv, root, server, view); err != nil {
			return nil, err
		}
	}
	return root.Children, nil
}

// readInventory reads the pools, servers, VMs, SRs, networks and disks of a server once, and the window's objects;
// the disks are those a view can show (organized ones), or with allVDIs every disk VDI.Show allows (a search).
func readInventory(c *xapi.Connection, allVDIs bool) (*inventory, error) {
	s := c.Session
	inv := &inventory{real: map[string]Object{}}
	var err error
	if inv.pools, err = xenapi.Pool.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	if inv.hosts, err = xenapi.Host.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	if inv.vms, err = xenapi.VM.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	inv.all = make(map[xenapi.VMRef]xenapi.VMRecord, len(inv.vms))
	for ref, r := range inv.vms {
		inv.all[ref] = r
	}
	for ref, r := range inv.vms { // VM.Show: no internal template, no "__gui__" object, nothing hidden from XCP-ng Center
		_, internal := r.OtherConfig["xensource_internal"]
		if internal || strings.HasPrefix(r.NameLabel, "__gui__") || r.OtherConfig["HideFromXenCenter"] == "true" {
			delete(inv.vms, ref)
		}
	}
	objs := map[string]Object{}
	// the VMs and templates come from the records the groupings read (a VM that starts meanwhile is never grouped as
	// running with a stopped VM's state); servers and SRs from the window's own listing (their state reads more)
	for _, kind := range []string{"vms", "templates"} {
		for _, o := range vmObjects(kind, inv.vms, inv.hosts, inv.pools) {
			objs[o.UUID] = o
			if kind == "vms" {
				inv.real[o.UUID] = o
			}
		}
	}
	for _, kind := range []string{"hosts", "srs"} {
		list, err := ListObjects(c, kind)
		if err != nil {
			return nil, err
		}
		for _, o := range list {
			objs[o.UUID] = o
		}
	}
	keep := func(uuid string, tags []string, other map[string]string) {
		if o, ok := objs[uuid]; ok {
			inv.items = append(inv.items, item{obj: o, tags: tags, other: other})
		}
	}
	for _, r := range inv.pools {
		name := r.NameLabel
		if name == "" { // Pool.Name: an unnamed pool is called after its coordinator
			name = inv.hosts[r.Master].NameLabel
		}
		objs[r.UUID] = Object{Kind: "pools", UUID: r.UUID, Name: name, Description: r.NameDescription}
		keep(r.UUID, r.Tags, r.OtherConfig)
	}
	for _, r := range inv.hosts {
		keep(r.UUID, r.Tags, r.OtherConfig)
	}
	for _, r := range inv.vms {
		if r.IsASnapshot { // ObjectTypes.Snapshot: in the views too
			objs[r.UUID] = Object{Kind: "snapshots", UUID: r.UUID, Name: r.NameLabel, Description: r.NameDescription}
		}
		keep(r.UUID, r.Tags, r.OtherConfig)
	}
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	inv.srs = srs
	for _, r := range srs {
		if !srShown(r) { // SR.Show
			continue
		}
		keep(r.UUID, r.Tags, r.OtherConfig)
	}
	nets, err := xenapi.Network.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	pifs, err := xenapi.PIF.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	for _, r := range nets {
		if !networkShown(r, pifs) {
			continue
		}
		objs[r.UUID] = Object{Kind: "networks", UUID: r.UUID, Name: r.NameLabel, Description: r.NameDescription}
		keep(r.UUID, r.Tags, r.OtherConfig)
	}
	vdis, err := xenapi.VDI.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	inv.vdis = vdis
	for _, r := range vdis { // ObjectTypes.VDI: those VDI.Show allows (not missing, managed, not hidden) and organized
		if r.Missing || !r.Managed || r.OtherConfig["HideFromXenCenter"] == "true" {
			continue
		}
		if allVDIs || len(r.Tags) > 0 || r.OtherConfig["folder"] != "" || hasCustomField(r.OtherConfig) {
			name := r.NameLabel
			if name == "" {
				name = messages.Get("NO_NAME")
			}
			objs[r.UUID] = Object{Kind: "vdis", UUID: r.UUID, Name: name, Description: r.NameDescription}
			keep(r.UUID, r.Tags, r.OtherConfig)
		}
	}
	return inv, nil
}

// networkShown is Network.Show without the hidden objects: not the guest installer network, not hidden, not a bond
// member's network, every PIF shown.
func networkShown(n xenapi.NetworkRecord, pifs map[xenapi.PIFRef]xenapi.PIFRecord) bool {
	if n.OtherConfig["is_guest_installer_network"] == "true" || n.OtherConfig["HideFromXenCenter"] == "true" {
		return false
	}
	for _, p := range n.PIFs {
		pif := pifs[p]
		if !pif.Managed || pif.OtherConfig["HideFromXenCenter"] == "true" || pif.BondSlaveOf != "" && pif.BondSlaveOf != nullRef { // PIF.Show
			return false
		}
	}
	return true
}

func hasCustomField(other map[string]string) bool {
	for k, v := range other {
		if strings.HasPrefix(k, "XenCenter.CustomFields.") && v != "" {
			return true
		}
	}
	return false
}

// folderPath is Folders.PointToPath: "/a/b" as its groups, each keyed by its whole path (never a custom field or a
// tag of the same name).
func folderPath(path string) []group {
	var out []group
	at := ""
	for _, f := range strings.Split(path, "/") {
		if f != "" {
			at += "/" + f
			out = append(out, group{f, "000_Folder_open_h32bit_16", "folder:" + at})
		}
	}
	return out
}

// customField is a CustomFieldDefinition: its name and type (String or Date).
type customField struct{ name, kind string }

// customFields are the pool's custom fields (gui_config "XenCenter.CustomFields", rules.ParseCustomFields).
func customFields(pools map[xenapi.PoolRef]xenapi.PoolRecord) []customField {
	var out []customField
	for _, p := range pools {
		for _, f := range rules.ParseCustomFields(p.GuiConfig[rules.CustomFieldsKey]) {
			out = append(out, customField{f.Name, f.Type})
		}
	}
	return out
}

// customFieldValue is CustomFieldsManager.ParseValue as the grouping shows it: a Date field's valid date in local
// time, grouped by the exact moment (key, whatever its text), a String as it is; "" for none or an invalid date.
func customFieldValue(f customField, v string) (label, key string) {
	if v == "" || f.kind != "Date" {
		return v, ""
	}
	if t, ok := rules.ParseFieldDate(v); ok {
		return t.Local().Format("2006-01-02 15:04:05"), "date:" + t.UTC().Format(time.RFC3339Nano)
	}
	return "", ""
}

// vmView groups the VMs as the vApps view and the default searches do.
func vmView(c *xapi.Connection, inv *inventory, root *ViewNode, server, view string) error {
	s := c.Session
	switch view {
	case "vapps", "vms_by_vapp", "vms_by_os", "vms_without_tools", "vms_by_power_state", "vms_by_network", "snapshots_by_vm":
	default:
		return nil
	}
	gms, err := xenapi.VMGuestMetrics.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	metrics, err := xenapi.VMMetrics.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	apps, err := xenapi.VMAppliance.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vifs, err := xenapi.VIF.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	nets, err := xenapi.Network.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	poolGroup := func() []group {
		for ref, p := range inv.pools {
			name := p.NameLabel
			if name == "" {
				name = inv.hosts[p.Master].NameLabel
			}
			return []group{{name, "000_PoolConnected_h32bit_16", server + "/" + string(ref)}}
		}
		return nil
	}
	vmNodes := map[xenapi.VMRef]*ViewNode{}
	var home func(xenapi.VMRecord) xenapi.HostRef
	for ref, r := range inv.vms {
		o, ok := inv.real[r.UUID]
		if !ok {
			continue
		}
		switch view {
		case "vapps", "vms_by_vapp":
			if a, ok := apps[r.Appliance]; ok {
				add(root, []group{{a.NameLabel, "000_VirtualAppliance_h32bit_16", server + "/" + a.UUID}}, o, server)
			}
		case "vms_by_os":
			add(root, []group{{osName(gms[r.GuestMetrics]), "000_Folder_open_h32bit_16", ""}}, o, server)
		case "vms_without_tools":
			if status, optimized := toolsStatus(s, r, gms[r.GuestMetrics], metrics[r.Metrics]); r.PowerState == xenapi.VMPowerStateRunning && !optimized {
				add(root, []group{{status, "000_Folder_open_h32bit_16", ""}}, o, server)
			}
		case "vms_by_power_state":
			add(root, []group{{string(r.PowerState), "000_Folder_open_h32bit_16", ""}}, o, server)
		case "vms_by_network": // a VM in each of its networks; one without any stays ungrouped
			done := map[xenapi.NetworkRef]bool{} // two VIFs on one network: the VM once
			for _, v := range r.VIFs {
				if n, ok := nets[vifs[v].Network]; ok && !done[vifs[v].Network] {
					done[vifs[v].Network] = true
					add(root, []group{{n.NameLabel, "000_Network_h32bit_16", server + "/" + n.UUID}}, o, server)
				}
			}
			if len(done) == 0 {
				add(root, nil, o, server)
			}
		case "snapshots_by_vm": // pool, the VM's home server (VM.Home), the VM; its snapshots below
			if home == nil {
				if home, err = homes(s, inv); err != nil {
					return err
				}
			}
			path := poolGroup()
			if h, ok := inv.hosts[home(r)]; ok {
				path = append(path, group{h.NameLabel, "000_TreeConnected_h32bit_16", server + "/" + h.UUID})
			}
			vmNodes[ref] = add(root, path, o, server)
		}
	}
	if view == "snapshots_by_vm" { // every snapshot, also one whose VM is gone (then under its pool)
		for _, r := range inv.vms {
			if !r.IsASnapshot {
				continue
			}
			node := &ViewNode{Name: r.NameLabel, Server: server, Object: &Object{Kind: "snapshots", UUID: r.UUID, Name: r.NameLabel, Description: r.NameDescription}, Children: []*ViewNode{}}
			if parent, ok := vmNodes[r.SnapshotOf]; ok {
				parent.Children = append(parent.Children, node)
			} else {
				pool := addGroups(root, poolGroup())
				pool.Children = append(pool.Children, node)
			}
		}
	}
	return nil
}

// homes is vmHome (VM.Home) over records read once: the server a VM runs on, else the one holding its local disks or
// CD, else its affinity when that server is live.
func homes(s *xenapi.Session, inv *inventory) (func(xenapi.VMRecord) xenapi.HostRef, error) {
	vbds, err := xenapi.VBD.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	vdis, err := xenapi.VDI.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	pbds, err := xenapi.PBD.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	hostMetrics, err := xenapi.HostMetrics.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	srHost := map[xenapi.SRRef]xenapi.HostRef{}
	for _, p := range pbds {
		if _, ok := srHost[p.SR]; !ok {
			srHost[p.SR] = p.Host
		}
	}
	return func(vm xenapi.VMRecord) xenapi.HostRef {
		if vm.ResidentOn != "" && vm.ResidentOn != nullRef {
			return vm.ResidentOn
		}
		for _, b := range vm.VBDs {
			vbd, ok := vbds[b]
			if !ok || vbd.Empty || vbd.VDI == nullRef {
				continue
			}
			if sr, ok := srs[vdis[vbd.VDI].SR]; ok && !sr.Shared {
				if h, ok := srHost[vdis[vbd.VDI].SR]; ok {
					return h
				}
			}
		}
		if h, ok := inv.hosts[vm.Affinity]; ok && hostMetrics[h.Metrics].Live {
			return vm.Affinity
		}
		return ""
	}, nil
}

// osName is VM.GetOSName: the guest's os_version name (before a "|", without "Microsoft", CA-9631), else Unknown.
func osName(gm xenapi.VMGuestMetricsRecord) string {
	name := gm.OsVersion["name"]
	if i := strings.Index(name, "|"); i >= 1 {
		name = name[:i]
	}
	for _, ms := range []string{"Microsoft®", "Microsoft"} {
		if strings.HasPrefix(name, ms) && name != ms {
			name = strings.TrimSpace(name[len(ms):])
			break
		}
	}
	if name == "" {
		return messages.Get("UNKNOWN")
	}
	return name
}

// toolsStatus is VM.GetVirtualizationStatus named as the search's grouping (VirtualisationStatus_i18n): Windows guests
// by their I/O drivers and management agent (feature-static-ip-setting), the others by their PV drivers' version;
// Unknown in a VM's first two minutes. optimized: nothing to install.
func toolsStatus(s *xenapi.Session, vm xenapi.VMRecord, gm xenapi.VMGuestMetricsRecord, m xenapi.VMMetricsRecord) (string, bool) {
	label, _, optimized := virtualizationStatus(s, vm, gm, m)
	return label, optimized
}

// virtualizationStatus is toolsStatus with the VM.VirtualizationStatus name a search compares ("NotInstalled",
// "IoDriversInstalled, ManagementInstalled", "PvDriversOutOfDate", "Unknown"...).
func virtualizationStatus(s *xenapi.Session, vm xenapi.VMRecord, gm xenapi.VMGuestMetricsRecord, m xenapi.VMMetricsRecord) (string, string, bool) {
	early := !m.StartTime.IsZero() && time.Since(m.StartTime) < 2*time.Minute // GetBodgeStartupTime: the guest may not have reported yet
	if isWindows(s, vm) {
		io := gm.PVDriversDetected
		agent := gm.Other["feature-static-ip-setting"] != "" && gm.Other["feature-static-ip-setting"] != "0"
		switch {
		case io && agent:
			return messages.Get("VIRTUALIZATION_STATE_VM_OPTIMIZED"), "IoDriversInstalled, ManagementInstalled", true
		case io:
			return messages.Get("VIRTUALIZATION_STATE_VM_IO_OPTIMIZED_ONLY"), "IoDriversInstalled", false
		case early:
			return messages.Get("UNKNOWN"), "Unknown", false
		case agent:
			return messages.Get("VIRTUALIZATION_STATE_VM_MANAGEMENT_AGENT_INSTALLED_ONLY"), "ManagementInstalled", false
		}
		return messages.Get("VIRTUALIZATION_STATE_VM_NOT_OPTIMIZED"), "NotInstalled", false
	}
	_, major := gm.PVDriversVersion["major"]
	_, minor := gm.PVDriversVersion["minor"]
	switch {
	case (!major || !minor) && early: // VM_guest_metrics.PV_drivers_installed
		return messages.Get("UNKNOWN"), "Unknown", false
	case !major || !minor:
		return messages.Get("VIRTUALIZATION_STATE_VM_NOT_OPTIMIZED"), "NotInstalled", false
	case !gm.PVDriversUpToDate:
		return messages.Get("OUT_OF_DATE"), "PvDriversOutOfDate", false
	}
	return messages.Get("VIRTUALIZATION_STATE_VM_OPTIMIZED"), "IoDriversInstalled, ManagementInstalled", true
}

// addGroups makes the groups of a path that are missing and returns the last one (root for none).
func addGroups(root *ViewNode, path []group) *ViewNode {
	n := root
	for _, g := range path {
		var next *ViewNode
		for _, c := range n.Children {
			if c.Object == nil && sameGroup(c, g.name, g.key) {
				next = c
				break
			}
		}
		if next == nil {
			next = &ViewNode{Name: g.name, Key: g.key, Icon: g.icon, Children: []*ViewNode{}}
			n.Children = append(n.Children, next)
		}
		n = next
	}
	return n
}

// sameGroup: groups that are objects match by identity (two networks may share a name), the others by value.
func sameGroup(n *ViewNode, name, key string) bool {
	if n.Key != "" || key != "" {
		return n.Key == key
	}
	return n.Name == name
}

// add puts an object under the group path (making the groups it lacks) and returns the object's node.
func add(root *ViewNode, path []group, o Object, server string) *ViewNode {
	n := addGroups(root, path)
	obj := o
	node := &ViewNode{Name: o.Name, Server: server, Object: &obj, Children: []*ViewNode{}}
	n.Children = append(n.Children, node)
	return node
}

// MergeViews joins the servers' parts of a view: value groups of the same name become one (the C# view spans every
// connection), object groups stay their server's; groups come first by name, then objects by name.
func MergeViews(parts ...[]*ViewNode) []*ViewNode {
	root := &ViewNode{}
	for _, p := range parts {
		merge(root, p)
	}
	sortView(root)
	return root.Children
}

func merge(into *ViewNode, nodes []*ViewNode) {
	for _, n := range nodes {
		if n.Object != nil {
			into.Children = append(into.Children, n)
			continue
		}
		var same *ViewNode
		for _, c := range into.Children {
			if c.Object == nil && sameGroup(c, n.Name, n.Key) {
				same = c
				break
			}
		}
		if same == nil {
			same = &ViewNode{Name: n.Name, Key: n.Key, Icon: n.Icon, Children: []*ViewNode{}}
			into.Children = append(into.Children, same)
		}
		merge(same, n.Children)
	}
}

func sortView(n *ViewNode) {
	sort.SliceStable(n.Children, func(i, j int) bool {
		a, b := n.Children[i], n.Children[j]
		if (a.Object == nil) != (b.Object == nil) {
			return a.Object == nil
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
	for _, c := range n.Children {
		sortView(c)
	}
}
