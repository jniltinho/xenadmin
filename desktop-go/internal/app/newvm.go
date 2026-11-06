package app

import (
	"encoding/xml"
	"sort"
	"strconv"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// NewVMTemplate is one row of the New VM wizard's template page (Page_Template).
// shortcut: no OS category column or logo (VmTemplateType); the custom templates come first, as the C# order.
type NewVMTemplate struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Custom bool   `json:"custom"`
}

// NewVMDisk is one disk of the storage page (DiskDescription): a new disk on an SR at a position.
type NewVMDisk struct {
	Position    string `json:"position"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int    `json:"size"`
	SR          string `json:"sr"` // SR uuid
	Bootable    bool   `json:"bootable"`
	Inherited   bool   `json:"inherited"` // a custom template's disk: kept (or moved to another SR), never resized nor left out
}

// NewVMDefaults are NewVmDefaults.FromTemplate, DisksFor and VifsFor for one template and home server.
type NewVMDefaults struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	HVM         bool         `json:"hvm"`
	BootMode    string       `json:"boot_mode"` // bios, uefi
	Install     string       `json:"install"`   // cd, network or none (CreateVMAction's InstallMethod)
	VCPUs       int          `json:"vcpus"`
	Cores       int64        `json:"cores"`
	Memory      int          `json:"memory"`
	Disks       []NewVMDisk  `json:"disks"`
	Networks    []ops.NewVIF `json:"networks"`
	NoDVD       bool         `json:"no_dvd"`    // a custom template without a DVD drive (Helpers.CustomWithNoDVD): no media page choice
	Eliloader   bool         `json:"eliloader"` // a PV template installed from a URL or the DVD drive
}

// NewVMOptions is everything the wizard's pages offer: templates, ISOs, servers, networks and SRs for new disks.
type NewVMOptions struct {
	Templates []NewVMTemplate `json:"templates"`
	ISOs      []ISO           `json:"isos"`
	Servers   []ServerPick    `json:"servers"`
	Networks  []ISO           `json:"networks"` // uuid (in VDI) and name
	SRs       []SRUse         `json:"srs"`
	DefaultSR string          `json:"default_sr"`
}

func defaultTemplate(t xenapi.VMRecord) bool { _, ok := t.OtherConfig["default_template"]; return ok }

// ReadNewVMOptions reads the wizard's choices (templates in Page_Template's order: custom ones first, then the newest
// version of each first).
func ReadNewVMOptions(c *xapi.Connection) (NewVMOptions, error) {
	s := c.Session
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return NewVMOptions{}, xapi.Sanitize(err)
	}
	out := NewVMOptions{Templates: []NewVMTemplate{}, ISOs: []ISO{}, Servers: []ServerPick{}, Networks: []ISO{}, SRs: []SRUse{}}
	for _, v := range vms {
		if !v.IsATemplate || v.IsASnapshot || v.OtherConfig["HideFromXenCenter"] == "true" {
			continue
		}
		if _, internal := v.OtherConfig["xensource_internal"]; internal {
			continue
		}
		out.Templates = append(out.Templates, NewVMTemplate{UUID: v.UUID, Name: v.NameLabel, Custom: !defaultTemplate(v)})
	}
	sort.SliceStable(out.Templates, func(i, j int) bool {
		a, b := out.Templates[i], out.Templates[j]
		if a.Custom != b.Custom {
			return a.Custom
		}
		return a.Name > b.Name // newest version first
	})
	// ISOs of every ISO SR the pool shows (the media page filters by the chosen server)
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return NewVMOptions{}, xapi.Sanitize(err)
	}
	for _, sr := range srs {
		if sr.ContentType != "iso" || sr.IsToolsSr || !srShown(sr) {
			continue
		}
		for _, vr := range sr.VDIs {
			if v, err := xenapi.VDI.GetRecord(s, vr); err == nil && v.Managed && !v.Missing && !v.IsASnapshot {
				out.ISOs = append(out.ISOs, ISO{VDI: v.UUID, Text: sr.NameLabel + ": " + v.NameLabel})
			}
		}
	}
	sort.Slice(out.ISOs, func(i, j int) bool { return out.ISOs[i].Text < out.ISOs[j].Text })
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return NewVMOptions{}, xapi.Sanitize(err)
	}
	for _, h := range hosts {
		live, _ := xenapi.HostMetrics.GetLive(s, h.Metrics)
		out.Servers = append(out.Servers, ServerPick{UUID: h.UUID, Name: h.NameLabel, Live: live})
	}
	sort.Slice(out.Servers, func(i, j int) bool { return out.Servers[i].Name < out.Servers[j].Name })
	nets, err := xenapi.Network.GetAllRecords(s)
	if err != nil {
		return NewVMOptions{}, xapi.Sanitize(err)
	}
	for _, n := range nets {
		if n.OtherConfig["HideFromXenCenter"] == "true" || n.OtherConfig["is_guest_installer_network"] == "true" {
			continue
		}
		out.Networks = append(out.Networks, ISO{VDI: n.UUID, Text: n.NameLabel})
	}
	sort.Slice(out.Networks, func(i, j int) bool { return out.Networks[i].Text < out.Networks[j].Text })
	if len(vms) > 0 { // the disk SRs: those every server can use (NewDiskOptions of any VM gives the same list)
		opts, err := newDiskSRs(s, "")
		if err != nil {
			return NewVMOptions{}, err
		}
		out.SRs, out.DefaultSR = opts.SRs, opts.DefaultSR
	}
	return out, nil
}

// newDiskSRs is StorageRules.NewDiskTargets for a home server ("" for any): the SRs that can hold a new disk.
func newDiskSRs(s *xenapi.Session, home xenapi.HostRef) (NewDisk, error) {
	out := NewDisk{SRs: []SRUse{}}
	sms, _ := xenapi.SM.GetAllRecords(s)
	canCreate := map[string]bool{}
	for _, sm := range sms {
		for _, cp := range sm.Capabilities {
			canCreate[sm.Type] = canCreate[sm.Type] || cp == "VDI_CREATE"
		}
	}
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
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
		if !canCreate[sr.Type] || sr.ContentType == "iso" || !srShown(sr) || broken(s, sr, len(hosts)) || (home != "" && !seenFrom(s, sr, home)) {
			continue
		}
		list = append(list, row{r, sr})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].rec.NameLabel < list[j].rec.NameLabel })
	hostRecs, _ := xenapi.Host.GetAllRecords(s)
	for _, r := range list {
		name := srI18N(s, hostRecs, r.rec, "name_label", r.rec.NameLabel) // SR.Name: "Local storage on host-a"
		if r.rec.Shared && loc != "" {
			name += " " + loc
		}
		out.SRs = append(out.SRs, SRUse{UUID: r.rec.UUID, Text: name + "   " + diskSize(r.rec.PhysicalSize-r.rec.PhysicalUtilisation) + " free of " + diskSize(r.rec.PhysicalSize)})
		if r.ref == def {
			out.DefaultSR = r.rec.UUID
		}
	}
	return out, nil
}

// ReadNewVMDefaults is NewVmDefaults.FromTemplate + DisksFor + VifsFor for a template and an optional home server.
func ReadNewVMDefaults(c *xapi.Connection, templateUUID, homeUUID string) (NewVMDefaults, error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, templateUUID)
	if err != nil {
		return NewVMDefaults{}, xapi.Sanitize(err)
	}
	t, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return NewVMDefaults{}, xapi.Sanitize(err)
	}
	var home xenapi.HostRef
	if homeUUID != "" {
		if home, err = xenapi.Host.GetByUUID(s, homeUUID); err != nil {
			return NewVMDefaults{}, xapi.Sanitize(err)
		}
	}
	vms, _ := xenapi.VM.GetAllRecords(s)
	taken := map[string]bool{}
	for _, v := range vms {
		if !v.IsATemplate {
			taken[v.NameLabel] = true
		}
	}
	hvm := t.HVMBootPolicy != ""
	cores := int64(1)
	if n, err := strconv.ParseInt(t.Platform["cores-per-socket"], 10, 64); err == nil && n > 0 {
		cores = n
	}
	d := NewVMDefaults{Name: rules.UniqueName(t.NameLabel, taken), HVM: hvm, BootMode: "bios", Install: map[bool]string{true: "cd", false: "network"}[hvm],
		VCPUs: max(1, t.VCPUsAtStartup), Cores: cores, Memory: t.MemoryStaticMax, Disks: []NewVMDisk{}, Networks: []ops.NewVIF{},
		Eliloader: !hvm && t.PVBootloader == "eliloader"}
	if !defaultTemplate(t) {
		d.NoDVD = true
		for _, b := range t.VBDs {
			if typ, err := xenapi.VBD.GetType(s, b); err == nil && typ == xenapi.VbdTypeCD {
				d.NoDVD = false
			}
		}
	}
	switch {
	case d.NoDVD:
		d.Install = "none" // the media page is skipped: nothing to install from, the boot order stays
	case d.Eliloader:
		d.Install = "cd" // a PV installer reads the DVD drive by default
	}
	if !defaultTemplate(t) {
		d.Description = t.NameDescription
	}
	if hvm && strings.EqualFold(strings.TrimSpace(t.HVMBootParams["firmware"]), "uefi") {
		d.BootMode = "uefi"
	}
	targets, err := newDiskSRs(s, home)
	if err != nil {
		return NewVMDefaults{}, err
	}
	// Page_Storage.GetBestDiskStorage: the suggested SR if usable, else the default, else the first usable (shared first)
	usable := map[string]bool{}
	for _, sr := range targets.SRs {
		usable[sr.UUID] = true
	}
	best := func(suggested string) string {
		if usable[suggested] {
			return suggested
		}
		if usable[targets.DefaultSR] {
			return targets.DefaultSR
		}
		if len(targets.SRs) > 0 {
			return targets.SRs[0].UUID
		}
		return ""
	}
	// Page_Storage.LoadDisks: the provision XML of a default template, else the template's own disks
	var prov struct {
		Disks []struct {
			Device   string `xml:"device,attr"`
			Size     string `xml:"size,attr"`
			SR       string `xml:"sr,attr"`
			Bootable string `xml:"bootable,attr"`
		} `xml:"disk"`
	}
	if x := t.OtherConfig["disks"]; x != "" && xml.Unmarshal([]byte(x), &prov) == nil && len(prov.Disks) > 0 {
		for _, p := range prov.Disks {
			size, _ := strconv.Atoi(p.Size)
			d.Disks = append(d.Disks, NewVMDisk{Position: p.Device, Name: d.Name + " " + p.Device,
				Description: messages.Get("NEWVMWIZARD_STORAGEPAGE_DISK_DESCRIPTION"), Size: size, SR: best(p.SR), Bootable: p.Bootable == "true"})
		}
	} else {
		for _, b := range t.VBDs {
			vbd, err := xenapi.VBD.GetRecord(s, b)
			if err != nil || vbd.Type != xenapi.VbdTypeDisk {
				continue
			}
			vdi, err := xenapi.VDI.GetRecord(s, vbd.VDI)
			if err != nil {
				continue
			}
			srUUID, _ := xenapi.SR.GetUUID(s, vdi.SR)
			d.Disks = append(d.Disks, NewVMDisk{Position: vbd.Userdevice, Name: vdi.NameLabel, Description: vdi.NameDescription,
				Size: vdi.VirtualSize, SR: best(srUUID), Bootable: vbd.Bootable, Inherited: true})
		}
	}
	// Page_Networking.LoadNetworks: up to 4 auto-plug networks for a default template, else the template's VIFs
	if defaultTemplate(t) {
		nets, _ := xenapi.Network.GetAllRecords(s)
		var names []string
		byName := map[string]string{}
		for _, n := range nets {
			if n.OtherConfig["is_guest_installer_network"] == "true" || n.OtherConfig["automatic"] == "false" || n.OtherConfig["HideFromXenCenter"] == "true" {
				continue
			}
			if bonded(s, n) {
				continue
			}
			names, byName[n.NameLabel+"\x00"+n.UUID] = append(names, n.NameLabel+"\x00"+n.UUID), n.UUID
		}
		sort.Strings(names)
		for _, n := range names {
			if len(d.Networks) == 4 {
				break
			}
			d.Networks = append(d.Networks, ops.NewVIF{Network: byName[n]})
		}
	} else {
		var vifs []xenapi.VIFRecord
		for _, v := range t.VIFs {
			if r, err := xenapi.VIF.GetRecord(s, v); err == nil {
				vifs = append(vifs, r)
			}
		}
		sort.Slice(vifs, func(i, j int) bool { return vifs[i].Device < vifs[j].Device })
		for _, r := range vifs { // device and bandwidth limit kept (VifsFor)
			if u, err := xenapi.Network.GetUUID(s, r.Network); err == nil {
				d.Networks = append(d.Networks, ops.NewVIF{Network: u, Device: r.Device, QosType: r.QosAlgorithmType, QosParams: r.QosAlgorithmParams})
			}
		}
	}
	return d, nil
}

// bonded is Network.IsMember: the network of a NIC that is part of a bond (not offered to VMs).
func bonded(s *xenapi.Session, n xenapi.NetworkRecord) bool {
	for _, p := range n.PIFs {
		if bm, err := xenapi.PIF.GetBondSlaveOf(s, p); err == nil && bm != "" && bm != nullRef {
			return true
		}
	}
	return false
}
