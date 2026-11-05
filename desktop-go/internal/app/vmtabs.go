package app

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

const nullRef = "OpaqueRef:NULL"

// Table is one of Details' tables (XenAdmin.Desktop.Core/Details.cs): a title, headers and rows; Keys identify the rows
// for the tab's buttons (the VBD or VIF uuid).
type Table struct {
	Title   string     `json:"title"`
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
	Keys    []string   `json:"keys"`
}

// Disk is what the Storage tab's buttons need about a row.
type Disk struct {
	VBD           string `json:"vbd"`
	Attached      bool   `json:"attached"`
	CanActivate   bool   `json:"can_activate"`   // StorageRules.CanActivate
	CanDeactivate bool   `json:"can_deactivate"` // StorageRules.CanDeactivate
	VDI           string `json:"vdi"`
	Name          string `json:"name"`
	System        bool   `json:"system"`     // bootable or a system disk: the Detach/Delete questions warn
	CanDetach     bool   `json:"can_detach"` // StorageRules.CanDetach
	CanDelete     bool   `json:"can_delete"` // StorageRules.CanDelete
}

// ISO is one choice of the DVD drive (StorageRules.IsoChoices); an empty VDI is "<empty>".
type ISO struct {
	VDI  string `json:"vdi"`
	Text string `json:"text"`
}

// VMStorage is the Storage tab of a VM (VmStorageView): the disks table, the disks' state and the DVD drive.
type VMStorage struct {
	Table     Table  `json:"table"`
	Disks     []Disk `json:"disks"`
	DVD       string `json:"dvd,omitempty"`      // the DVD drive's VBD uuid ("" when the VM has none)
	CanAddDVD bool   `json:"can_add_dvd"`        // a real VM without a DVD drive (the "create a DVD drive" link)
	DVDName   string `json:"dvd_name,omitempty"` // "DVD Drive 3"
	ISO       string `json:"iso"`                // the inserted ISO's VDI uuid ("" when empty)
	ISOs      []ISO  `json:"isos"`
}

func yesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// diskSize is Util.DiskSizeString: binary units, one decimal.
func diskSize(bytes int) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	v, i := float64(bytes), 0
	for v >= 1024 && i < len(units)-1 {
		v, i = v/1024, i+1
	}
	if i == 0 {
		return fmt.Sprintf("%d B", bytes)
	}
	return fmt.Sprintf("%.1f %s", v, units[i])
}

// location is XenObject.LocationString for the connection: "in 'pool'" or "on 'server'".
func location(s *xenapi.Session) string {
	pools, _ := xenapi.Pool.GetAllRecords(s)
	hosts, _ := xenapi.Host.GetAllRecords(s)
	for _, p := range pools {
		if p.NameLabel != "" || len(hosts) > 1 {
			name := p.NameLabel
			if name == "" {
				name = hosts[p.Master].NameLabel
			}
			return messages.Format(messages.Get("IN_POOL"), name)
		}
		return messages.Format(messages.Get("ON_SERVER"), hosts[p.Master].NameLabel)
	}
	return ""
}

// ListVMStorage is Details.Storage for a VM plus VmStorageView's DVD drive.
// shortcut: no Add, Attach, Move, Detach, Delete or Properties of disks yet (StorageRules has them); add with 6.1.
func ListVMStorage(c *xapi.Connection, vmUUID string) (VMStorage, error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return VMStorage{}, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return VMStorage{}, xapi.Sanitize(err)
	}
	out := VMStorage{Table: Table{Title: "Virtual Disks", Headers: []string{"Position", "Name", "Description", "SR", "Size", "Read Only", "Priority", "Active", "Device Path"},
		Rows: [][]string{}, Keys: []string{}}, Disks: []Disk{}, ISOs: []ISO{{Text: "<empty>"}}}
	real := !vm.IsATemplate && !vm.IsASnapshot && !vm.IsControlDomain
	running := vm.PowerState == xenapi.VMPowerStateRunning
	loc := location(s)
	type row struct {
		vbd xenapi.VBDRecord
		vdi xenapi.VDIRecord
	}
	var rows []row
	for _, b := range vm.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil {
			return VMStorage{}, xapi.Sanitize(err)
		}
		if vbd.Type == xenapi.VbdTypeCD {
			if out.DVD == "" {
				out.DVD, out.DVDName = vbd.UUID, "DVD Drive "+vbd.Userdevice
				if !vbd.Empty && vbd.VDI != nullRef {
					out.ISO, _ = xenapi.VDI.GetUUID(s, vbd.VDI)
				}
			}
			continue
		}
		if vbd.Type == xenapi.VbdTypeFloppy || vbd.VDI == nullRef {
			continue
		}
		vdi, err := xenapi.VDI.GetRecord(s, vbd.VDI)
		if err != nil {
			continue // the disk is gone (VmDisks skips it)
		}
		rows = append(rows, row{vbd, vdi})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].vbd.Userdevice < rows[j].vbd.Userdevice }) // ordinal, as VmDisks
	for _, r := range rows {
		sr, srErr := xenapi.SR.GetRecord(s, r.vdi.SR)
		srName := sr.NameLabel
		if sr.Shared && loc != "" { // SR.NameWithLocation: local SRs show their name only
			srName += " " + loc
		}
		desc := r.vdi.NameDescription
		if desc == "" {
			desc = "No description"
		}
		prio := strconv.Itoa(ioNice(r.vbd))
		switch prio {
		case "0":
			prio = messages.Get("STORAGE_WORST")
		case "7":
			prio = messages.Get("STORAGE_BEST")
		}
		dev := ""
		if r.vbd.Device != "" {
			dev = "/dev/" + r.vbd.Device
		}
		out.Table.Rows = append(out.Table.Rows, []string{r.vbd.Userdevice, r.vdi.NameLabel, desc, srName, diskSize(r.vdi.VirtualSize),
			yesNo(r.vbd.Mode == xenapi.VbdModeRO), prio, yesNo(r.vbd.CurrentlyAttached), dev})
		out.Table.Keys = append(out.Table.Keys, r.vbd.UUID)
		has := func(op xenapi.VbdOperations) bool {
			for _, a := range r.vbd.AllowedOperations {
				if a == op {
					return true
				}
			}
			return false
		}
		info, infoErr := ops.DiskInfo(s, ref, r.vdi, sr) // on a failed read the disk can be neither detached nor deleted
		if srErr != nil {
			infoErr = srErr
		}
		out.Disks = append(out.Disks, Disk{VBD: r.vbd.UUID, Attached: r.vbd.CurrentlyAttached,
			CanActivate:   real && running && !r.vbd.CurrentlyAttached && has(xenapi.VbdOperationsPlug),
			CanDeactivate: real && running && r.vbd.CurrentlyAttached && has(xenapi.VbdOperationsUnplug),
			VDI:           r.vdi.UUID, Name: r.vdi.NameLabel, System: r.vbd.Bootable || r.vdi.Type == xenapi.VdiTypeSystem,
			CanDetach: infoErr == nil && rules.CanDetachDisk(info), CanDelete: infoErr == nil && rules.CanDeleteDisk(info)})
	}
	out.CanAddDVD = out.DVD == "" && real
	if out.DVD != "" {
		isos, err := isoChoices(s, vm)
		if err != nil {
			return VMStorage{}, xapi.Sanitize(err)
		}
		out.ISOs = append(out.ISOs, isos...)
	}
	return out, nil
}

// NewDisk is what the Add Virtual Disk dialog offers (NewDiskDialog): a unique default name, the SRs that can hold the
// disk with their free space, and the pool's default SR.
type NewDisk struct {
	Name      string  `json:"name"`
	SRs       []SRUse `json:"srs"`
	DefaultSR string  `json:"default_sr"`
}

// SRUse is one SR of the picker.
type SRUse struct {
	UUID string `json:"uuid"`
	Text string `json:"text"` // "name   10.0 GB free of 20.0 GB"
}

// NewDiskOptions is NewDiskDialog's defaults and StorageRules.NewDiskTargets for the VM's home server.
func NewDiskOptions(c *xapi.Connection, vmUUID string) (NewDisk, error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return NewDisk{}, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return NewDisk{}, xapi.Sanitize(err)
	}
	vdis, err := xenapi.VDI.GetAllRecords(s)
	if err != nil {
		return NewDisk{}, xapi.Sanitize(err)
	}
	taken := map[string]bool{}
	for _, v := range vdis {
		taken[v.NameLabel] = true
	}
	out := NewDisk{Name: rules.UniqueName(messages.Get("DEFAULT_VDI_NAME"), taken), SRs: []SRUse{}}
	sms, _ := xenapi.SM.GetAllRecords(s)
	canCreate := map[string]bool{}
	for _, sm := range sms {
		for _, c := range sm.Capabilities {
			canCreate[sm.Type] = canCreate[sm.Type] || c == "VDI_CREATE"
		}
	}
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return NewDisk{}, xapi.Sanitize(err)
	}
	home := vmHome(s, vm)
	hosts, _ := xenapi.Host.GetAll(s)
	pools, _ := xenapi.Pool.GetAllRecords(s)
	var def xenapi.SRRef
	for _, p := range pools {
		def = p.DefaultSR
	}
	loc := location(s)
	type row struct {
		ref xenapi.SRRef
		rec xenapi.SRRecord
	}
	var list []row
	for r, sr := range srs {
		// StorageRules.CanHoldNewDisk: VDI_CREATE, not ISO nor memory, plugged everywhere, shown, seen from home
		if !canCreate[sr.Type] || sr.ContentType == "iso" || !srShown(sr) {
			continue
		}
		if broken(s, sr, len(hosts)) || (home != "" && !seenFrom(s, sr, home)) {
			continue
		}
		list = append(list, row{r, sr})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].rec.NameLabel < list[j].rec.NameLabel })
	for _, r := range list {
		name := r.rec.NameLabel
		if r.rec.Shared && loc != "" {
			name += " " + loc
		}
		out.SRs = append(out.SRs, SRUse{UUID: r.rec.UUID, Text: fmt.Sprintf("%s   %s free of %s", name,
			diskSize(r.rec.PhysicalSize-r.rec.PhysicalUtilisation), diskSize(r.rec.PhysicalSize))})
		if r.ref == def {
			out.DefaultSR = r.rec.UUID
		}
	}
	return out, nil
}

// srShown is SR.Show(false): not a temporary "__gui__" SR, not the old local CD type nor memory, not introduced by a
// disaster recovery task, not hidden.
func srShown(sr xenapi.SRRecord) bool {
	return !strings.HasPrefix(sr.NameLabel, "__gui__") && sr.Type != "local" && sr.Type != "tmpfs" &&
		(sr.IntroducedBy == "" || sr.IntroducedBy == nullRef) && sr.OtherConfig["HideFromXenCenter"] != "true"
}

// broken is SR.IsBroken: no PBD, a local SR with more than one, a shared SR without one per server, or one unplugged.
func broken(s *xenapi.Session, sr xenapi.SRRecord, hosts int) bool {
	if len(sr.PBDs) == 0 || (!sr.Shared && len(sr.PBDs) != 1) || (sr.Shared && len(sr.PBDs) != hosts) {
		return true
	}
	for _, p := range sr.PBDs {
		if attached, err := xenapi.PBD.GetCurrentlyAttached(s, p); err != nil || !attached {
			return true
		}
	}
	return false
}

// ioNice is VBD.GetIoNice: the "class" of the I/O scheduling parameters, 0 when unset.
func ioNice(vbd xenapi.VBDRecord) int {
	n, _ := strconv.Atoi(vbd.QosAlgorithmParams["class"])
	return n
}

// isoChoices is StorageRules.IsoChoices: the ISOs of ISO SRs shown in the app and visible from the VM's home server,
// without the tools SR (XCP-ng guests get their agent from the distribution).
func isoChoices(s *xenapi.Session, vm xenapi.VMRecord) ([]ISO, error) {
	home := vmHome(s, vm)
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return nil, err
	}
	type srRow struct {
		ref xenapi.SRRef
		rec xenapi.SRRecord
	}
	var list []srRow
	for r, sr := range srs {
		if sr.ContentType != "iso" || sr.IsToolsSr || !srShown(sr) {
			continue
		}
		if home != "" && home != nullRef && !seenFrom(s, sr, home) { // SR.CanBeSeenFrom, shared SRs too
			continue
		}
		list = append(list, srRow{r, sr})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].rec.NameLabel < list[j].rec.NameLabel })
	var out []ISO
	for _, sr := range list {
		var vdis []xenapi.VDIRecord
		for _, v := range sr.rec.VDIs {
			// VDI.Show(false): managed, present, not hidden, not a snapshot
			if vdi, err := xenapi.VDI.GetRecord(s, v); err == nil && vdi.Managed && !vdi.Missing && vdi.OtherConfig["HideFromXenCenter"] != "true" && !vdi.IsASnapshot {
				vdis = append(vdis, vdi)
			}
		}
		sort.Slice(vdis, func(i, j int) bool { return vdis[i].NameLabel < vdis[j].NameLabel })
		for _, v := range vdis {
			out = append(out, ISO{VDI: v.UUID, Text: sr.rec.NameLabel + ": " + v.NameLabel})
		}
	}
	return out, nil
}

// vmHome is VM.Home: the server it runs on, else the one holding its local disks or CD (VM.GetStorageHost), else
// its affinity when that server is live; "" when none (every server then).
func vmHome(s *xenapi.Session, vm xenapi.VMRecord) xenapi.HostRef {
	if vm.ResidentOn != "" && vm.ResidentOn != nullRef {
		return vm.ResidentOn
	}
	for _, b := range vm.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil || vbd.Empty || vbd.VDI == nullRef { // disks and an inserted CD (GetStorageHost(false))
			continue
		}
		vdi, err := xenapi.VDI.GetRecord(s, vbd.VDI)
		if err != nil {
			continue
		}
		sr, err := xenapi.SR.GetRecord(s, vdi.SR)
		if err != nil || sr.Shared {
			continue
		}
		for _, p := range sr.PBDs {
			if pbd, err := xenapi.PBD.GetRecord(s, p); err == nil {
				return pbd.Host
			}
		}
	}
	if vm.Affinity != "" && vm.Affinity != nullRef { // a home server that is down does not count
		if metrics, err := xenapi.Host.GetMetrics(s, vm.Affinity); err == nil {
			if live, err := xenapi.HostMetrics.GetLive(s, metrics); err == nil && live {
				return vm.Affinity
			}
		}
	}
	return ""
}

// seenFrom is SR.CanBeSeenFrom: one of its PBDs is the host's and is plugged.
func seenFrom(s *xenapi.Session, sr xenapi.SRRecord, host xenapi.HostRef) bool {
	for _, p := range sr.PBDs {
		if pbd, err := xenapi.PBD.GetRecord(s, p); err == nil && pbd.Host == host && pbd.CurrentlyAttached {
			return true
		}
	}
	return false
}

// ListVMNetworking is Details.Networking for a VM (VmNetworkView's table).
// shortcut: no Add, Properties, Remove or Activate of interfaces yet (NetworkRules has them).
func ListVMNetworking(c *xapi.Connection, vmUUID string) (Table, error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return Table{}, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return Table{}, xapi.Sanitize(err)
	}
	t := Table{Title: "Virtual Network Interfaces", Headers: []string{"Device", "MAC", "Limit", "Network", "IP Address", "Active"}, Rows: [][]string{}, Keys: []string{}}
	gm, _ := xenapi.VMGuestMetrics.GetRecord(s, vm.GuestMetrics)
	var vifs []xenapi.VIFRecord
	for _, v := range vm.VIFs {
		if vif, err := xenapi.VIF.GetRecord(s, v); err == nil {
			vifs = append(vifs, vif)
		}
	}
	sort.Slice(vifs, func(i, j int) bool { return vifs[i].Device < vifs[j].Device })
	for _, vif := range vifs {
		net := "-"
		if n, err := xenapi.Network.GetNameLabel(s, vif.Network); err == nil {
			net = n
		}
		// "0/ip" and "0/ipv4/0" hold the same address: list it once
		var ips []string
		seen := map[string]bool{}
		keys := make([]string, 0, len(gm.Networks))
		for k := range gm.Networks {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if v := gm.Networks[k]; strings.HasPrefix(k, vif.Device+"/ip") && !seen[v] {
				seen[v], ips = true, append(ips, v)
			}
		}
		limit := ""
		if vif.QosAlgorithmType != "" { // NetworkRules.Limit / VIF.LimitString
			limit = vif.QosAlgorithmParams["kbps"]
		}
		t.Rows = append(t.Rows, []string{vif.Device, vif.MAC, limit, net, strings.Join(ips, ", "), yesNo(vif.CurrentlyAttached)})
		t.Keys = append(t.Keys, vif.UUID)
	}
	return t, nil
}

// ListVMMemory is the Memory tab's figures (MemoryView): the static and dynamic limits and what the VM uses now.
// shortcut: figures only; editing (MemoryRules.Validate, Settings, ChangeAction) and the host bar come next.
func ListVMMemory(c *xapi.Connection, vmUUID string) (Table, error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return Table{}, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return Table{}, xapi.Sanitize(err)
	}
	mb := func(b int) string { return fmt.Sprintf("%d MB", b>>20) }
	t := Table{Title: "Memory", Headers: []string{"", ""}, Keys: []string{}}
	if vm.PowerState == xenapi.VMPowerStateRunning {
		if actual, err := xenapi.VMMetrics.GetMemoryActual(s, vm.Metrics); err == nil {
			t.Rows = append(t.Rows, []string{"Current memory usage", mb(actual)})
		}
	}
	if vm.MemoryDynamicMin == vm.MemoryStaticMax { // fixed memory: one figure, as the C# rubric
		t.Rows = append(t.Rows, []string{"Memory", mb(vm.MemoryStaticMax)})
	} else {
		t.Rows = append(t.Rows, []string{"Minimum memory (dynamic)", mb(vm.MemoryDynamicMin)}, []string{"Maximum memory (dynamic)", mb(vm.MemoryDynamicMax)},
			[]string{"Maximum memory (static)", mb(vm.MemoryStaticMax)})
	}
	return t, nil
}
