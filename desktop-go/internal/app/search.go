package app

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// A search saved in a pool (XenModel/XenSearch: SearchMarshalling XML, major version 2) is a query (the object
// types of QueryScope and a filter) and a grouping, each element named after its C# class. This is the part of
// XenSearch the saved searches use most; a search with any other filter, property or grouping is listed, not run.
// shortcut: no XenModelObject*Query, Recursive*Query or CustomField*Query yet; IPAddressQuery
// reads VM and server addresses (not an SR's target);
// port them from QueryTypes.cs when a saved search needs one.

// xmlNode is any element of a search, with its attributes and children.
type xmlNode struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Children []xmlNode  `xml:",any"`
}

func (n xmlNode) attr(name string) string {
	for _, a := range n.Attrs {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}

// search is a parsed saved search: its name, the object types it wants, its filter and its groupings (outermost first).
type search struct {
	name      string
	types     map[string]bool
	filter    *xmlNode
	groupings []xmlNode
}

// objectTypes are ObjectTypes as QueryScope writes them, with the group names of a type grouping (ObjectTypes_i18n).
// DisconnectedServer is left out: the window does not make objects of the disconnected servers (GroupAlg) yet.
var objectTypes = map[string]string{"Pool": "POOLS", "Server": "SERVERS", "VM": "VMS",
	"Snapshot": "SNAPSHOTS", "UserTemplate": "CUSTOM_TEMPLATES", "DefaultTemplate": "XENSERVER_TEMPLATES",
	"RemoteSR": "REMOTE_SRS", "LocalSR": "LOCAL_SRS", "VDI": "VIRTUAL_DISKS", "Network": "NETWORKS"}

// filterProperties are the properties each filter reads here, of the type the C# PropertyAccessors give them; any
// other pair is refused, never run with a wrong meaning.
var filterProperties = map[string]map[string]bool{
	"StringPropertyQuery":     {"label": true, "uuid": true, "description": true, "os_name": true},
	"EnumPropertyQuery":       {"type": true, "power_state": true, "virtualisation_status": true},
	"BooleanQuery":            {"shared": true, "has_custom_fields": true, "in_any_appliance": true, "ha_enabled": true},
	"ListEmptyQuery":          {"tags": true},
	"StringListContainsQuery": {"tags": true},
	"IntPropertyQuery":        {"memory": true, "size": true},
	"NullQuery":               {"pool": true, "folder": true, "appliance": true},
	"IPAddressQuery":          {"ip_address": true},
	"DatePropertyQuery":       {"start_time": true},
}

// groupingProperties are the properties each grouping reads here.
var groupingProperties = map[string]map[string]bool{
	"PropertyGrouping":               {"type": true, "power_state": true, "os_name": true, "tags": true, "virtualisation_status": true},
	"BoolGrouping":                   {"shared": true, "has_custom_fields": true, "in_any_appliance": true, "ha_enabled": true},
	"XenModelObjectPropertyGrouping": {"pool": true, "host": true, "networks": true, "appliance": true},
}

// parseSearch reads a saved search's XML; an error names what this port cannot run.
func parseSearch(text string) (*search, error) {
	var root xmlNode
	if err := xml.Unmarshal([]byte(text), &root); err != nil {
		return nil, err
	}
	if v, _ := strconv.Atoi(root.attr("major_version")); v != 2 || len(root.Children) == 0 || root.Children[0].XMLName.Local != "Query" {
		return nil, fmt.Errorf("search version %q", root.attr("major_version"))
	}
	s := &search{name: root.attr("name"), types: map[string]bool{}}
	scoped := false
	for _, c := range root.Children[0].Children {
		if c.XMLName.Local == "QueryScope" { // an empty scope wants nothing, as the C# QueryScope
			scoped = true
			for _, t := range c.Children {
				name := t.XMLName.Local
				if name == "NONE" || strings.HasPrefix(name, "ALL") {
					continue
				}
				if _, ok := objectTypes[name]; !ok {
					return nil, fmt.Errorf("object type %s", name)
				}
				s.types[name] = true
			}
			continue
		}
		f := c
		if err := checkFilter(f); err != nil {
			return nil, err
		}
		s.filter = &f
	}
	if !scoped { // Query without a scope: every type but folders
		for t := range objectTypes {
			s.types[t] = true
		}
	}
	for _, c := range root.Children[1:] {
		if c.XMLName.Local == "columns" || c.XMLName.Local == "sorting" {
			continue
		}
		for g := &c; g != nil; { // a grouping holds its subgrouping
			if err := checkGrouping(*g); err != nil {
				return nil, err
			}
			s.groupings = append(s.groupings, xmlNode{XMLName: g.XMLName, Attrs: g.Attrs})
			var next *xmlNode
			for i := range g.Children {
				if g.Children[i].XMLName.Local != "CustomFieldDefinition" {
					next = &g.Children[i]
					break
				}
			}
			g = next
		}
	}
	return s, nil
}

func checkFilter(f xmlNode) error {
	switch name := f.XMLName.Local; name {
	case "DummyQuery":
		return nil
	case "GroupQuery":
		if t := f.attr("type"); t != "And" && t != "Or" && t != "Nor" {
			return fmt.Errorf("group query type %q", t)
		}
		for _, c := range f.Children {
			if err := checkFilter(c); err != nil {
				return err
			}
		}
		return nil
	default:
		if props, ok := filterProperties[name]; !ok || !props[f.attr("property")] {
			return fmt.Errorf("%s on %q", name, f.attr("property"))
		}
		if t := f.attr("type"); name == "DatePropertyQuery" && !strings.Contains(" today yesterday thisweek lastweek before after exact ", " "+t+" ") {
			return fmt.Errorf("date query type %q", t) // Enum.Parse refuses it
		}
		if _, err := rules.ParseDateQuery(f.attr("query")); name == "DatePropertyQuery" && err != nil {
			return fmt.Errorf("date %q", f.attr("query"))
		}
		return nil
	}
}

func checkGrouping(g xmlNode) error {
	switch name := g.XMLName.Local; name {
	case "FolderGrouping", "AllCustomFieldsGrouping":
		return nil
	default:
		if props, ok := groupingProperties[name]; !ok || !props[g.attr("property")] {
			return fmt.Errorf("%s by %q", name, g.attr("property"))
		}
		return nil
	}
}

// SearchSupported tells whether this port runs a saved search (and why not).
func SearchSupported(text string) error { _, err := parseSearch(text); return err }

// facts are the properties of an object a search reads; a property an object does not have is absent (nil in the
// C# accessors: such a filter does not match, as PropertyQuery's nullProtect).
type facts struct {
	values map[string]string   // string, enum and boolean properties ("true"/"false"), as the XML writes them
	labels map[string]string   // the group name of an enum value (virtualisation_status)
	ints   map[string]int64    // memory, size
	lists  map[string][]string // tags
	groups map[string][]group  // pool, host, networks, appliance: the objects they are, for a grouping
	self   string              // the object's group key when it is a pool, a server or a network (BelongsAsGroupNotMember)
	folder []group
	fields []customFieldGroup
	ips    []string             // IPAddressProperty: a VM's guest addresses, a server's PIF addresses; none is null
	noIPs  bool                 // ...except a template's or snapshot's, an empty list (IsRealVm is false)
	dates  map[string]time.Time // start_time of a real VM (VM_metrics.start_time; the zero time without metrics)
}

type customFieldGroup struct{ name, value, key string }

// RunSearch runs a saved search over one server: the objects of its types its filter matches, grouped (Query.Match,
// Search.PopulateAdapters).
func RunSearch(c *xapi.Connection, server, text string) ([]*ViewNode, error) {
	sr, err := parseSearch(text)
	if err != nil {
		return nil, err
	}
	s := c.Session
	inv, err := readInventory(c, sr.types["VDI"])
	if err != nil {
		return nil, err
	}
	gms, err := xenapi.VMGuestMetrics.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	metrics, err := xenapi.VMMetrics.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	apps, err := xenapi.VMAppliance.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	pifs, err := xenapi.PIF.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	vifs, err := xenapi.VIF.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	nets, err := xenapi.Network.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	pbds, err := xenapi.PBD.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	home, err := homes(s, inv)
	if err != nil {
		return nil, err
	}
	fields := customFields(inv.pools)
	vmByUUID := map[string]xenapi.VMRecord{}
	for _, r := range inv.vms {
		vmByUUID[r.UUID] = r
	}
	srByUUID := map[string]xenapi.SRRecord{}
	for _, r := range inv.srs {
		srByUUID[r.UUID] = r
	}
	vdiByUUID := map[string]xenapi.VDIRecord{}
	for _, r := range inv.vdis {
		vdiByUUID[r.UUID] = r
	}
	netByUUID := map[string]xenapi.NetworkRecord{}
	for _, r := range nets {
		netByUUID[r.UUID] = r
	}
	hostGroup := func(h xenapi.HostRecord) group {
		return group{h.NameLabel, "000_TreeConnected_h32bit_16", server + "/" + h.UUID}
	}
	var allHosts []group
	for _, h := range inv.hosts {
		allHosts = append(allHosts, hostGroup(h))
	}
	// Helpers.GetPool: the pool, when it is one (named or of several servers); a standalone server has none
	var pool *group
	poolUUID := ""
	for _, p := range inv.pools {
		if p.NameLabel != "" || len(inv.hosts) > 1 {
			name := p.NameLabel
			if name == "" {
				name = inv.hosts[p.Master].NameLabel
			}
			pool, poolUUID = &group{name, "000_PoolConnected_h32bit_16", server + "/" + p.UUID}, p.UUID
		}
	}
	// SR.Home: a local SR's server (one PBD)
	srHome := func(r xenapi.SRRecord) []group {
		if r.Shared || len(r.PBDs) != 1 {
			return nil
		}
		if h, ok := inv.hosts[pbds[r.PBDs[0]].Host]; ok {
			return []group{hostGroup(h)}
		}
		return nil
	}
	root := &ViewNode{}
	for _, it := range inv.items {
		f := facts{values: map[string]string{"label": it.obj.Name, "uuid": it.obj.UUID, "description": it.obj.Description},
			labels: map[string]string{}, ints: map[string]int64{}, lists: map[string][]string{"tags": it.tags}, groups: map[string][]group{},
			folder: folderPath(it.other["folder"])}
		if pool != nil {
			f.values["pool"], f.groups["pool"] = poolUUID, []group{*pool}
		}
		if len(f.folder) > 0 {
			f.values["folder"] = it.other["folder"]
		}
		hasField := false
		for _, cf := range fields {
			raw := it.other["XenCenter.CustomFields."+cf.name]
			hasField = hasField || raw != "" // CustomFieldsManager.HasCustomFields: a value, valid or not
			if v, key := customFieldValue(cf, raw); v != "" {
				if key == "" {
					key = "cf:" + cf.name + "=" + v
				}
				f.fields = append(f.fields, customFieldGroup{cf.name, v, key})
			}
		}
		f.values["has_custom_fields"] = strconv.FormatBool(hasField)
		var hosts []group // HostProperty
		switch it.obj.Kind {
		case "pools":
			if pool == nil { // Cache.XenSearchableObjects: a standalone server's pool of one is no object
				continue
			}
			f.values["type"] = "Pool"
			for _, p := range inv.pools {
				f.values["ha_enabled"] = strconv.FormatBool(p.HaEnabled)
			}
			if pool != nil {
				f.self = pool.key
			}
		case "hosts":
			f.values["type"] = "Server"
			for ref, h := range inv.hosts {
				if h.UUID == it.obj.UUID {
					hosts = []group{hostGroup(h)}
					f.self = hostGroup(h).key
					for _, p := range pifs {
						if p.Host == ref && p.IP != "" {
							f.ips = append(f.ips, p.IP)
						}
					}
				}
			}
		case "vms", "templates", "snapshots":
			r := vmByUUID[it.obj.UUID]
			_, defaultTemplate := r.OtherConfig["default_template"] // VM.DefaultTemplate: the key is there
			switch {
			case r.IsASnapshot:
				f.values["type"] = "Snapshot"
				f.noIPs = true
				if src, ok := inv.all[r.SnapshotOf]; ok { // VM.Home of a snapshot: its VM's, hidden or not
					if h, ok := inv.hosts[home(src)]; ok {
						hosts = []group{hostGroup(h)}
					}
				}
			case r.IsATemplate && defaultTemplate:
				f.values["type"], f.noIPs = "DefaultTemplate", true
			case r.IsATemplate:
				f.values["type"], f.noIPs = "UserTemplate", true
			default: // IsRealVm: these properties are a real VM's only
				f.values["type"] = "VM"
				f.values["power_state"] = string(r.PowerState)
				f.values["os_name"] = osName(gms[r.GuestMetrics])
				// VM.GetVirtualizationStatus: Unknown unless running with metrics
				f.values["virtualisation_status"], f.labels["virtualisation_status"] = "Unknown", messages.Get("UNKNOWN")
				if _, ok := metrics[r.Metrics]; ok && r.PowerState == xenapi.VMPowerStateRunning {
					label, value, _ := virtualizationStatus(s, r, gms[r.GuestMetrics], metrics[r.Metrics])
					f.values["virtualisation_status"], f.labels["virtualisation_status"] = value, label
				}
				f.dates = map[string]time.Time{"start_time": {}}
				if m, ok := metrics[r.Metrics]; ok {
					f.ints["memory"] = int64(m.MemoryActual)
					f.dates["start_time"] = m.StartTime
				}
				if gm, ok := gms[r.GuestMetrics]; ok { // Helpers.FindIpAddresses for each VIF
					for _, v := range r.VIFs {
						for k, addr := range gm.Networks {
							if strings.HasPrefix(k, vifs[v].Device+"/ip") {
								f.ips = append(f.ips, strings.FieldsFunc(strings.ReplaceAll(addr, "%n", "\n"), func(c rune) bool { return c == '\n' })...)
							}
						}
					}
				}
				if h, ok := inv.hosts[home(r)]; ok {
					hosts = []group{hostGroup(h)}
				}
				done := map[xenapi.NetworkRef]bool{}
				for _, v := range r.VIFs {
					if n, ok := nets[vifs[v].Network]; ok && !done[vifs[v].Network] {
						done[vifs[v].Network] = true
						f.groups["networks"] = append(f.groups["networks"], group{n.NameLabel, "000_Network_h32bit_16", server + "/" + n.UUID})
					}
				}
				a, inApp := apps[r.Appliance]
				f.values["in_any_appliance"] = strconv.FormatBool(inApp)
				if inApp {
					f.values["appliance"] = a.UUID
					f.groups["appliance"] = []group{{a.NameLabel, "000_VirtualAppliance_h32bit_16", server + "/" + a.UUID}}
				}
			}
		case "srs":
			r := srByUUID[it.obj.UUID]
			f.values["type"] = "RemoteSR"
			if localSRTypes[r.Type] {
				f.values["type"] = "LocalSR"
			}
			f.values["shared"] = strconv.FormatBool(r.Shared)
			hosts = srHome(r)
		case "networks":
			f.values["type"] = "Network"
			if n := netByUUID[it.obj.UUID]; len(n.PIFs) == 0 { // a private network: every server
				hosts = allHosts
			}
			f.self = server + "/" + it.obj.UUID
			f.groups["networks"] = []group{{it.obj.Name, "000_Network_h32bit_16", f.self}} // NetworksProperty: itself
		case "vdis":
			r := vdiByUUID[it.obj.UUID]
			f.values["type"] = "VDI"
			f.ints["size"] = int64(r.VirtualSize)
			f.values["shared"] = strconv.FormatBool(len(r.VBDs) >= 2) // SharedProperty: on two VMs or more
			hosts = srHome(inv.srs[r.SR])
		}
		if pool == nil { // not in a pool: every object under the server
			hosts = allHosts
		}
		f.groups["host"] = hosts
		if !sr.types[f.values["type"]] || sr.filter != nil && match(*sr.filter, f) == no {
			continue
		}
		for _, path := range groupPaths(sr.groupings, f) {
			if len(path) > 0 && f.self != "" && path[len(path)-1].key == f.self {
				addGroups(root, path) // XenModelObjectPropertyGrouping.BelongsAsGroupNotMember: the object is its group
				continue
			}
			add(root, path, it.obj, server)
		}
	}
	return root.Children, nil
}

// localSRTypes are SR.IsLocalSR's types.
var localSRTypes = map[string]bool{"local": true, "file": true, "ext": true, "lvm": true, "udev": true, "egeneracd": true,
	"dummy": true, "ext4": true, "xfs": true}

// tri is the C# bool? of QueryFilter.Match: yes, no, or no opinion (a GroupQuery of none).
type tri int

const (
	none tri = iota
	yes
	no
)

func of(b bool) tri {
	if b {
		return yes
	}
	return no
}

// match is QueryFilter.Match for the filters this port reads (QueryTypes.cs).
func match(q xmlNode, f facts) tri {
	prop := q.attr("property")
	value, has := f.values[prop]
	switch q.XMLName.Local {
	case "DummyQuery":
		return yes
	case "GroupQuery":
		result := none
		for _, c := range q.Children {
			b := match(c, f)
			switch q.attr("type") {
			case "And":
				if b == no {
					return no
				}
				if b == yes {
					result = yes
				}
			case "Or":
				if b == yes {
					return yes
				}
				if b == no {
					result = no
				}
			case "Nor":
				if b == yes {
					return no
				}
				if b == no {
					result = yes
				}
			}
		}
		return result
	case "StringPropertyQuery":
		if !has {
			return no
		}
		return of(matchString(value, q.attr("query"), q.attr("type"), parseBool(q.attr("casesensitive"))))
	case "EnumPropertyQuery":
		if !has {
			return no
		}
		want := q.attr("query")
		if prop == "type" && want == "SR" { // EnumPropertyQuery<ObjectTypes>: an old search's "SR" is RemoteSR
			want = "RemoteSR"
		}
		return of((value == want) == parseBool(q.attr("equals")))
	case "BooleanQuery":
		if !has {
			return no
		}
		return of(value == strconv.FormatBool(parseBool(q.attr("query"))))
	case "NullQuery": // a missing object is null: XenModelObjects here are the folder and the appliance
		return of((value == "") == parseBool(q.attr("query")))
	case "ListEmptyQuery":
		list, ok := f.lists[prop]
		if !ok {
			return no
		}
		return of((len(list) == 0) == parseBool(q.attr("empty")))
	case "DatePropertyQuery":
		v, ok := f.dates[prop]
		if !ok {
			return no
		}
		day, _ := rules.ParseDateQuery(q.attr("query")) // checked by checkFilter
		return of(rules.MatchDate(v, day, q.attr("type"), time.Now()))
	case "IPAddressQuery": // ListContainsQuery<ComparableAddress>: a null list matches nothing
		want, ok := rules.ParseAddress(q.attr("address"), true, true)
		if !ok || len(f.ips) == 0 && !f.noIPs {
			return no
		}
		found := false
		for _, v := range f.ips {
			if a, ok := rules.ParseAddress(v, false, true); ok && want.Matches(a) {
				found = true
			}
		}
		return of(found == parseBool(q.attr("contains")))
	case "StringListContainsQuery":
		list, ok := f.lists[prop]
		if !ok {
			return no
		}
		found := false
		for _, v := range list {
			found = found || v == q.attr("query")
		}
		return of(found == parseBool(q.attr("contains")))
	case "IntPropertyQuery":
		n, ok := f.ints[prop]
		want, err := strconv.ParseInt(q.attr("query"), 10, 64)
		if !ok || err != nil {
			return no
		}
		switch q.attr("type") {
		case "gt":
			return of(n > want)
		case "lt":
			return of(n < want)
		}
		return of(n == want)
	}
	return no
}

// parseBool is SearchMarshalling.ParseBool: "yes", or "True" from older searches.
func parseBool(v string) bool { return v == "yes" || v == "True" }

// matchString is StringPropertyQuery.MatchString: an empty query matches everything.
func matchString(candidate, query, kind string, caseSensitive bool) bool {
	if query == "" {
		return true
	}
	if !caseSensitive {
		candidate, query = strings.ToLower(candidate), strings.ToLower(query)
	}
	switch kind {
	case "contains":
		return strings.Contains(candidate, query)
	case "startswith":
		return strings.HasPrefix(candidate, query)
	case "endswith":
		return strings.HasSuffix(candidate, query)
	case "exactmatch":
		return candidate == query
	case "notcontain":
		return !strings.Contains(candidate, query)
	}
	return false
}

// groupPaths are the paths an object goes under for the groupings (each level may give several groups, as tags and
// networks). A level with no group for the object is skipped and the next one applies (NodeGroup.AddUngrouped),
// except a folder grouping, under which an object in no folder stays (FolderGroup); an object that is the group of
// a level (a pool, a server, a network) stops there.
func groupPaths(groupings []xmlNode, f facts) [][]group {
	paths := [][]group{nil}
	for _, g := range groupings {
		var level [][]group
		switch g.XMLName.Local {
		case "FolderGrouping": // FolderGroup: its members go in a leaf, no subgrouping applies under it
			if len(f.folder) == 0 {
				return paths
			}
			var out [][]group
			for _, p := range paths {
				if len(p) > 0 && f.self != "" && p[len(p)-1].key == f.self { // the object is already its group
					out = append(out, p)
					continue
				}
				out = append(out, append(append([]group{}, p...), f.folder...))
			}
			return out
		case "AllCustomFieldsGrouping":
			for _, cf := range f.fields {
				level = append(level, []group{{cf.name, "000_Fields_h32bit_16", "cf:" + cf.name}, {cf.value, "000_Fields_h32bit_16", cf.key}})
			}
		default: // PropertyGrouping, BoolGrouping, XenModelObjectPropertyGrouping
			prop := g.attr("property")
			switch {
			case g.XMLName.Local == "XenModelObjectPropertyGrouping":
				for _, x := range f.groups[prop] {
					level = append(level, []group{x})
				}
			case prop == "tags": // GroupKey: the grouping's identity with the value, never a value of another property
				for _, t := range f.lists["tags"] {
					level = append(level, []group{{t, "000_Tag_h32bit_16", "tags:" + t}})
				}
			case prop == "type" && f.values["type"] != "":
				name := messages.Format(messages.Get(objectTypes[f.values["type"]]), "XCP-ng")
				level = append(level, []group{{name, "000_Folder_open_h32bit_16", "type:" + f.values["type"]}})
			case f.labels[prop] != "":
				level = append(level, []group{{f.labels[prop], "000_Folder_open_h32bit_16", prop + ":" + f.values[prop]}})
			case f.values[prop] != "":
				level = append(level, []group{{f.values[prop], "000_Folder_open_h32bit_16", prop + ":" + f.values[prop]}})
			}
		}
		if len(level) == 0 {
			continue
		}
		var next [][]group
		for _, p := range paths {
			if len(p) > 0 && f.self != "" && p[len(p)-1].key == f.self { // the object is already its group
				next = append(next, p)
				continue
			}
			for _, l := range level {
				next = append(next, append(append([]group{}, p...), l...))
			}
		}
		paths = next
	}
	return paths
}
