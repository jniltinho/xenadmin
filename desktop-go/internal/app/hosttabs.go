package app

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// The Memory, Storage, Networking and NICs tabs of a server or pool, and the Storage tab of an SR (Details.Memory,
// Storage, Networking and Nics). shortcut: read-only tables; their buttons (New Network, bonds, Add disk...) and the
// SR-IOV and tunnel columns come later.

// hostTabData is what the tabs read of a connection.
type hostTabData struct {
	s     *xenapi.Session
	hosts map[xenapi.HostRef]xenapi.HostRecord
	pifs  map[xenapi.PIFRef]xenapi.PIFRecord
	pm    map[xenapi.PIFMetricsRef]xenapi.PIFMetricsRecord
	bonds map[xenapi.BondRef]xenapi.BondRecord
	// the pool's coordinator (Network.Name prefers its PIF)
	coordinator xenapi.HostRef
}

func readHostTabData(s *xenapi.Session) (*hostTabData, error) {
	d := &hostTabData{s: s}
	var err error
	if d.hosts, err = xenapi.Host.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	if d.pifs, err = xenapi.PIF.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	if d.pm, err = xenapi.PIFMetrics.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	if d.bonds, err = xenapi.Bond.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	if pools, err := xenapi.Pool.GetAllRecords(s); err == nil {
		for _, p := range pools {
			d.coordinator = p.Master
		}
	}
	return d, nil
}

// pifPhysical is PIF.IsPhysical: not a VLAN, tunnel access or SR-IOV logical PIF.
func pifPhysical(p xenapi.PIFRecord) bool {
	return p.VLAN == -1 && len(p.TunnelAccessPIFOf) == 0 && len(p.SriovLogicalPIFOf) == 0
}

// pifName is PIF.Name: "NIC 0" or "Bond 0+1" for a physical PIF (BondSuffix: the members' numbers, sorted), else the
// name of its physical PIF on the same device.
func (d *hostTabData) pifName(p xenapi.PIFRecord) string {
	if len(p.BondMasterOf) > 0 {
		var nums []string
		for _, b := range p.BondMasterOf {
			for _, sl := range d.bonds[b].Slaves {
				nums = append(nums, strings.TrimPrefix(d.pifs[sl].Device, "eth"))
			}
		}
		sort.Strings(nums) // PIF.BondSuffix
		return "Bond " + strings.Join(nums, "+")
	}
	if pifPhysical(p) {
		return "NIC " + strings.TrimPrefix(p.Device, "eth")
	}
	for _, o := range d.pifs { // a VLAN: the physical PIF of its server on the same device
		if o.Host == p.Host && o.Device == p.Device && pifPhysical(o) {
			return d.pifName(o)
		}
	}
	return p.Device
}

// linkState is PIF.LinkStatus for a physical or VLAN PIF: its metrics' carrier, unknown without metrics.
// shortcut: tunnel access and SR-IOV PIFs are not told apart.
func (d *hostTabData) linkState(p xenapi.PIFRecord) string {
	m, ok := d.pm[p.Metrics]
	switch {
	case !ok:
		return "Unknown"
	case m.Carrier:
		return "Connected"
	}
	return "Disconnected"
}

// HostTab is one of these tabs for a server ("hosts"), a pool ("pools") or an SR ("srs").
func HostTab(c *xapi.Connection, kind, uuid, tab string) (Table, error) {
	s := c.Session
	var host xenapi.HostRef // "" for a pool
	switch kind {
	case "hosts":
		r, err := xenapi.Host.GetByUUID(s, uuid)
		if err != nil {
			return Table{}, xapi.Sanitize(err)
		}
		host = r
	case "srs":
		if tab != "storage" {
			return Table{}, fmt.Errorf("an SR has no %s tab", tab)
		}
		return srDisks(s, uuid)
	case "pools":
	default:
		return Table{}, fmt.Errorf("no %s tab for %s", tab, kind)
	}
	d, err := readHostTabData(s)
	if err != nil {
		return Table{}, err
	}
	switch tab {
	case "memory":
		return d.memory(host)
	case "storage":
		return d.storage(host)
	case "networking":
		return d.networking(host)
	case "nics":
		if host == "" {
			return Table{}, fmt.Errorf("a pool has no NICs tab")
		}
		return d.nics(host), nil
	}
	return Table{}, fmt.Errorf("unknown tab %s", tab)
}

// memory is Details.Memory (HostMemory): the server's free and total memory, its VMs' and Xen's; for a pool, each
// server in a section of its own.
func (d *hostTabData) memory(only xenapi.HostRef) (Table, error) {
	s := d.s
	t := Table{Title: "Memory", Headers: []string{"", ""}, Rows: [][]string{}, Keys: []string{}}
	var refs []xenapi.HostRef
	for r := range d.hosts {
		if only == "" || r == only {
			refs = append(refs, r)
		}
	}
	sort.Slice(refs, func(i, j int) bool { return d.hosts[refs[i]].NameLabel < d.hosts[refs[j]].NameLabel })
	for _, r := range refs {
		h := d.hosts[r]
		if only == "" {
			t.Rows = append(t.Rows, []string{h.NameLabel, ""})
		}
		hm, err := xenapi.HostMetrics.GetRecord(s, h.Metrics)
		if err != nil {
			t.Rows = append(t.Rows, []string{"Server", "Unknown"}, []string{"VMs", "Unknown"}, []string{"XCP-ng", "Unknown"})
			continue
		}
		xen, dynMin := int64(h.MemoryOverhead), int64(0) // Host.xen_memory_calc and tot_dyn_min
		var used []string
		for _, vr := range h.ResidentVMs {
			vm, err := xenapi.VM.GetRecord(s, vr)
			if err != nil {
				continue
			}
			xen += int64(vm.MemoryOverhead)
			m, merr := xenapi.VMMetrics.GetRecord(s, vm.Metrics)
			if vm.IsControlDomain {
				if merr == nil {
					xen += int64(m.MemoryActual)
				}
				continue
			}
			dynMin += int64(vm.MemoryStaticMax) // Host.tot_dyn_min: the static maximum of a VM without ballooning
			if gm, err := xenapi.VMGuestMetrics.GetRecord(s, vm.GuestMetrics); err == nil {
				if _, balloon := gm.Other["feature-balloon"]; balloon { // VM.SupportsBallooning
					dynMin += int64(vm.MemoryDynamicMin) - int64(vm.MemoryStaticMax)
				}
			}
			if merr == nil { // ResidentVMMemoryUsageString
				used = append(used, fmt.Sprintf("%s: using %s", vm.NameLabel, rules.MemoryText(int64(m.MemoryActual))))
			}
		}
		avail := max(0, int64(hm.MemoryTotal)-dynMin-xen) // Host.memory_available_calc
		t.Rows = append(t.Rows,
			[]string{"Server", fmt.Sprintf("%s RAM available (%s total)", rules.MemoryText(avail), rules.MemoryText(int64(hm.MemoryTotal)))},
			[]string{"VMs", strings.Join(used, "\n")},
			[]string{"XCP-ng", rules.MemoryText(xen)})
	}
	return t, nil
}

// storage is Details.Storage for a server or pool (StorageRules.PhysicalSrs): the shown SRs, not the tools SR, shared
// or local to the server.
func (d *hostTabData) storage(host xenapi.HostRef) (Table, error) {
	s := d.s
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return Table{}, xapi.Sanitize(err)
	}
	t := Table{Title: "Storage Repositories", Headers: []string{"Name", "Description", "Type", "Shared", "Usage", "Size", "Virtual allocation"},
		Rows: [][]string{}, Keys: []string{}}
	var list []xenapi.SRRecord
	for _, sr := range srs {
		if !srShown(sr) || sr.IsToolsSr || sr.NameLabel == "XenServer Tools" { // SR.IsToolsSR
			continue
		}
		if host != "" && !sr.Shared { // SR.Home: the server of its one PBD
			if len(sr.PBDs) != 1 {
				continue
			}
			if h, err := xenapi.PBD.GetHost(s, sr.PBDs[0]); err != nil || h != host {
				continue
			}
		}
		list = append(list, sr)
	}
	names := map[string][2]string{}
	for _, sr := range list {
		names[sr.UUID] = [2]string{srI18N(s, d.hosts, sr, "name_label", sr.NameLabel), srI18N(s, d.hosts, sr, "name_description", sr.NameDescription)}
	}
	sort.Slice(list, func(i, j int) bool { return names[list[i].UUID][0] < names[list[j].UUID][0] })
	for _, sr := range list {
		usage := ""
		if sr.PhysicalSize > 0 {
			usage = fmt.Sprintf("%d%% (%s used)", int64(sr.PhysicalUtilisation)*100/int64(sr.PhysicalSize), diskSize(sr.PhysicalUtilisation))
		}
		t.Rows = append(t.Rows, []string{names[sr.UUID][0], names[sr.UUID][1], srTypeName(sr.Type), yesNo(sr.Shared), usage,
			diskSize(sr.PhysicalSize), diskSize(sr.VirtualAllocation)})
		t.Keys = append(t.Keys, sr.UUID)
	}
	return t, nil
}

// networking is Details.Networking for a server or pool: the shown networks with the server's PIF on each (a pool
// shows the network's link status over all its PIFs).
func (d *hostTabData) networking(host xenapi.HostRef) (Table, error) {
	nets, err := xenapi.Network.GetAllRecords(d.s)
	if err != nil {
		return Table{}, xapi.Sanitize(err)
	}
	t := Table{Title: "Networks", Headers: []string{"Name", "Description", "NIC", "VLAN", "Auto", "Link Status", "MAC", "MTU"}, Rows: [][]string{}, Keys: []string{}}
	var list []xenapi.NetworkRecord
	for _, n := range nets {
		if networkShown(n, d.pifs) {
			list = append(list, n)
		}
	}
	sort.Slice(list, func(i, j int) bool { return d.networkName(list[i]) < d.networkName(list[j]) })
	for _, n := range list {
		var pif *xenapi.PIFRecord // Helpers.FindPIF
		for _, r := range n.PIFs {
			if p, ok := d.pifs[r]; ok && (host == "" || p.Host == host) {
				pif = &p
				break
			}
		}
		nic, vlan, link, mac := "", " - ", "None", " - "
		if pif != nil {
			nic = d.pifName(*pif)
			if pif.VLAN != -1 {
				vlan = strconv.Itoa(pif.VLAN)
			}
			link = d.linkState(*pif)
			if pifPhysical(*pif) {
				mac = pif.MAC
			}
		}
		if host == "" { // Network.LinkStatusString
			link = d.networkLink(n)
		}
		mtu := strconv.Itoa(n.MTU)
		if slices.ContainsFunc(n.PIFs, func(r xenapi.PIFRef) bool { return len(d.pifs[r].TunnelAccessPIFOf) > 0 }) { // CanUseJumboFrames
			mtu = " - "
		}
		t.Rows = append(t.Rows, []string{d.networkName(n), n.NameDescription, nic, vlan, yesNo(n.OtherConfig["automatic"] != "false"), link, mac, mtu})
		t.Keys = append(t.Keys, n.UUID)
	}
	return t, nil
}

// networkLink is Network.LinkStatusString: Connected, Partially connected or Disconnected over the network's PIFs.
func (d *hostTabData) networkLink(n xenapi.NetworkRecord) string {
	if len(n.PIFs) == 0 {
		return "None"
	}
	states := map[string]bool{}
	for _, r := range n.PIFs {
		states[d.linkState(d.pifs[r])] = true
	}
	switch {
	case states["Connected"] && (states["Disconnected"] || states["Unknown"]):
		return "Partially connected"
	case states["Connected"]:
		return "Connected"
	case states["Disconnected"] && !states["Unknown"]:
		return "Disconnected"
	}
	return "Unknown"
}

// nics is Details.Nics (NICPage): the server's physical NICs, by name.
// shortcut: no SR-IOV Capable column.
func (d *hostTabData) nics(host xenapi.HostRef) Table {
	t := Table{Title: "Network Interface Cards", Headers: []string{"NIC", "MAC", "Link Status", "Speed", "Duplex", "Vendor", "Device", "PCI Bus Path", "FCoE Capable"},
		Rows: [][]string{}, Keys: []string{}}
	var list []xenapi.PIFRecord
	for _, p := range d.pifs {
		if p.Host == host && pifPhysical(p) {
			list = append(list, p)
		}
	}
	sort.Slice(list, func(i, j int) bool { return d.pifName(list[i]) < d.pifName(list[j]) })
	for _, p := range list {
		m, ok := d.pm[p.Metrics]
		up := ok && m.Carrier
		link, speed, duplex, vendor, device, bus := "Disconnected", "-", "-", "-", "-", "-"
		if up {
			link, speed, duplex = "Connected", fmt.Sprintf("%d Mbit/s", m.Speed), map[bool]string{true: "Full", false: "Half"}[m.Duplex]
		}
		if ok {
			vendor, device, bus = m.VendorName, m.DeviceName, m.PciBusPath
		}
		t.Rows = append(t.Rows, []string{d.pifName(p), p.MAC, link, speed, duplex, vendor, device, bus, yesNo(slices.Contains(p.Capabilities, "fcoe"))})
		t.Keys = append(t.Keys, p.UUID)
	}
	return t
}

// srDisks is Details.Storage for an SR (StorageRules.SrDisks): its shown disks, by name, with the VMs using them.
// shortcut: no Storage Volume or Changed Block Tracking columns.
func srDisks(s *xenapi.Session, uuid string) (Table, error) {
	ref, err := xenapi.SR.GetByUUID(s, uuid)
	if err != nil {
		return Table{}, xapi.Sanitize(err)
	}
	vdiRefs, err := xenapi.SR.GetVDIs(s, ref)
	if err != nil {
		return Table{}, xapi.Sanitize(err)
	}
	t := Table{Title: "Virtual Disks", Headers: []string{"Name", "Description", "Size", "Virtual Machine"}, Rows: [][]string{}, Keys: []string{}}
	var list []xenapi.VDIRecord
	for _, r := range vdiRefs {
		v, err := xenapi.VDI.GetRecord(s, r)
		if err != nil || v.Missing || !v.Managed || v.OtherConfig["HideFromXenCenter"] == "true" { // VDI.Show
			continue
		}
		if _, mirror := v.SmConfig["base_mirror"]; mirror { // IsAnIntermediateStorageMotionSnapshot
			continue
		}
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].NameLabel < list[j].NameLabel })
	for _, v := range list {
		size := diskSize(v.VirtualSize) // VDI.SizeText: none for a snapshot
		if v.IsASnapshot {
			size = ""
		}
		t.Rows = append(t.Rows, []string{v.NameLabel, v.NameDescription, size, vmsOfVDI(s, v)})
		t.Keys = append(t.Keys, v.UUID)
	}
	return t, nil
}

// srTypeName is SR.FriendlyTypeName: the FriendlyNames label of the type, else the type itself.
func srTypeName(t string) string {
	if n := messages.Friendly("Label-SR.SRTypes-" + t); n != "" {
		return n
	}
	return t
}

// vmsOfVDI is VDI.VMsOfVDI: the VMs of the disk (VDI.GetVMs: a crash dump's VM, the VM suspended to it, else the VMs of
// its VBDs), sorted, the shown ones, a snapshot as "name (snapshot)".
func vmsOfVDI(s *xenapi.Session, v xenapi.VDIRecord) string {
	var refs []xenapi.VMRef
	switch v.Type {
	case xenapi.VdiTypeCrashdump:
		for _, cd := range v.CrashDumps {
			if vm, err := xenapi.Crashdump.GetVM(s, cd); err == nil {
				refs = append(refs, vm)
			}
		}
	case xenapi.VdiTypeSuspend:
		if all, err := xenapi.VM.GetAllRecords(s); err == nil {
			for r, vm := range all {
				if suspend, err := xenapi.VDI.GetUUID(s, vm.SuspendVDI); err == nil && suspend == v.UUID {
					refs = append(refs, r)
				}
			}
		}
	default:
		for _, b := range v.VBDs {
			if vm, err := xenapi.VBD.GetVM(s, b); err == nil {
				refs = append(refs, vm)
			}
		}
	}
	var names []string
	for _, r := range refs {
		vm, err := xenapi.VM.GetRecord(s, r)
		if err != nil || !vmShownHidden(vm) {
			continue
		}
		if vm.IsASnapshot {
			names = append(names, messages.Format(messages.Get("SNAPSHOT_BRACKETS"), vm.NameLabel))
		} else {
			names = append(names, vm.NameLabel)
		}
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

// vmShownHidden is VM.Show(true): not an internal template (xensource_internal) nor a GUI temporary object.
func vmShownHidden(vm xenapi.VMRecord) bool {
	_, internal := vm.OtherConfig["xensource_internal"]
	return !internal && !strings.HasPrefix(vm.NameLabel, "__gui__")
}

// srI18N is SR.I18N for SR.Name and SR.Description: an SR the server made and nobody renamed shows the translated
// name, with its server's name for a local SR ("Local storage on host-a").
func srI18N(s *xenapi.Session, hosts map[xenapi.HostRef]xenapi.HostRecord, sr xenapi.SRRecord, field, value string) string {
	key, ok := sr.OtherConfig["i18n-key"]
	if !ok || sr.OtherConfig["i18n-original-value-"+field] != value {
		return value
	}
	host := "" // SR.GetStorageHost: the server of its one PBD, for an SR that is not shared
	if !sr.Shared && len(sr.PBDs) == 1 {
		if h, err := xenapi.PBD.GetHost(s, sr.PBDs[0]); err == nil {
			host = hosts[h].NameLabel
		}
	}
	if host == "" {
		if pattern := messages.Friendly("SR." + field + "-" + key); pattern != "" {
			return messages.Format(pattern, "XCP-ng")
		}
		return value
	}
	if pattern := messages.Friendly("SR." + field + "-" + key + "-host"); pattern != "" {
		return messages.Format(pattern, host)
	}
	return value
}

// networkName is Network.Name: a default name ("Network associated with eth0") as "Network 0" or "Network 0 on VLAN 5"
// after the coordinator's PIF (else the first one), the host internal management network translated.
func (d *hostTabData) networkName(n xenapi.NetworkRecord) string {
	if n.NameLabel == "Host internal management network" {
		return messages.Get("HOST_INTERNAL_MANAGEMENT_NETWORK")
	}
	if !strings.HasPrefix(n.NameLabel, "Network associated with ") && !strings.HasPrefix(n.NameLabel, "Pool-wide network associated with ") {
		return n.NameLabel
	}
	var first, coord *xenapi.PIFRecord
	for _, r := range n.PIFs {
		p, ok := d.pifs[r]
		if !ok {
			continue
		}
		if first == nil {
			first = &p
		}
		if d.coordinator != "" && p.Host == d.coordinator && coord == nil {
			coord = &p
		}
	}
	if coord == nil {
		coord = first
	}
	if coord == nil {
		return n.NameLabel
	}
	id := d.nicIdentifier(*coord)
	if coord.VLAN >= 0 {
		return messages.Format(messages.Get("VLAN_NETWORK_NAME"), id, coord.VLAN)
	}
	return messages.Format(messages.Get("NETWORK_NAME"), id)
}

// nicIdentifier is PIF.NICIdentifier: the device number, or the bond's members.
func (d *hostTabData) nicIdentifier(p xenapi.PIFRecord) string {
	if len(p.BondMasterOf) == 0 {
		return strings.TrimPrefix(p.Device, "eth")
	}
	return strings.TrimPrefix(d.pifName(p), "Bond ")
}
