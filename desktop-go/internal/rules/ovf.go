package rules

import (
	"cmp"
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"

	"xcpng-center/internal/messages"
)

// The OVF descriptor of an appliance and the VMs it describes (XenOvf: OVF.cs, the Definitions types, Validation.cs,
// and ImportApplianceAction's Import.cs, which turns each VirtualSystem into a VM, its disks and its interfaces).
// The descriptor comes from the user's file: it is read with a size limit, Go's XML decoder expands no entities, and
// every file it names must be a plain name inside the package.

// IsAppliance is OvfRules.IsAppliance: an OVF descriptor or an OVA archive (optionally compressed), by its name.
func IsAppliance(file string) bool {
	name := strings.ToLower(strings.TrimSpace(file))
	for _, ext := range []string{".ovf", ".ova", ".ova.gz", ".ova.bz2"} {
		if strings.HasSuffix(name, ext) {
			return true
		}
	}
	return false
}

// ovfMaxDescriptor is the largest descriptor read (a real one is a few kB; this stops a hostile one filling memory).
const ovfMaxDescriptor = 16 << 20

type ovfMsg struct {
	MsgID string `xml:"msgid,attr"`
	Value string `xml:",chardata"`
}

type ovfFile struct {
	ID          string `xml:"id,attr"`
	Href        string `xml:"href,attr"`
	Size        string `xml:"size,attr"`
	Compression string `xml:"compression,attr"`
}

type ovfDisk struct {
	DiskID     string `xml:"diskId,attr"`
	FileRef    string `xml:"fileRef,attr"`
	Capacity   string `xml:"capacity,attr"`
	Units      string `xml:"capacityAllocationUnits,attr"`
	Format     string `xml:"format,attr"`
	IsBootable bool   `xml:"isBootable,attr"`
}

type ovfEula struct {
	License ovfMsg `xml:"License"`
}

type ovfStartup struct {
	Items []struct {
		ID         string `xml:"id,attr"`
		Order      int    `xml:"order,attr"`
		StartDelay int    `xml:"startDelay,attr"`
		StopDelay  int    `xml:"stopDelay,attr"`
	} `xml:"Item"`
}

// ovfRASD is a ResourceAllocationSettingData item of a VirtualHardwareSection.
type ovfRASD struct {
	ResourceType    string   `xml:"ResourceType"`
	ResourceSubType string   `xml:"ResourceSubType"`
	InstanceID      string   `xml:"InstanceID"`
	ElementName     string   `xml:"ElementName"`
	Caption         string   `xml:"Caption"`
	Description     string   `xml:"Description"`
	VirtualQuantity string   `xml:"VirtualQuantity"`
	AllocationUnits string   `xml:"AllocationUnits"`
	Limit           string   `xml:"Limit"`
	HostResource    []string `xml:"HostResource"`
	Connection      []string `xml:"Connection"`
	Address         string   `xml:"Address"`
	AddressOnParent string   `xml:"AddressOnParent"`
	Parent          string   `xml:"Parent"`
}

func (r *ovfRASD) kind() int { n, _ := strconv.Atoi(strings.TrimSpace(r.ResourceType)); return n }

type ovfHardware struct {
	System struct {
		VirtualSystemType       string `xml:"VirtualSystemType"`
		VirtualSystemIdentifier string `xml:"VirtualSystemIdentifier"`
		Description             string `xml:"Description"`
	} `xml:"System"`
	Items       []ovfRASD `xml:"Item"`
	OtherConfig []struct {
		Name  string `xml:"Name,attr"`
		Value string `xml:"Value"`
	} `xml:"VirtualSystemOtherConfigurationData"`
}

type ovfSystem struct {
	ID       string        `xml:"id,attr"`
	Name     string        `xml:"Name"`
	Hardware []ovfHardware `xml:"VirtualHardwareSection"`
	Eulas    []ovfEula     `xml:"EulaSection"`
	OS       []struct {
		Description string `xml:"Description"`
	} `xml:"OperatingSystemSection"`
	Install []struct{} `xml:"InstallSection"`
}

type ovfCollection struct {
	ID          string          `xml:"id,attr"`
	Systems     []ovfSystem     `xml:"VirtualSystem"`
	Collections []ovfCollection `xml:"VirtualSystemCollection"`
	Eulas       []ovfEula       `xml:"EulaSection"`
	Startup     []ovfStartup    `xml:"StartupSection"`
}

type ovfEnvelope struct {
	XMLName    xml.Name       `xml:"Envelope"`
	Version    string         `xml:"version,attr"`
	Files      []ovfFile      `xml:"References>File"`
	Disks      []ovfDisk      `xml:"DiskSection>Disk"`
	Networks   []struct{}     `xml:"NetworkSection>Network"`
	System     *ovfSystem     `xml:"VirtualSystem"`
	Collection *ovfCollection `xml:"VirtualSystemCollection"`
	Eulas      []ovfEula      `xml:"EulaSection"`
	Startup    []ovfStartup   `xml:"StartupSection"`
	Encryption []struct{}     `xml:"SecuritySection"`
	Strings    []struct {
		Lang string   `xml:"lang,attr"`
		Msgs []ovfMsg `xml:"Msg"`
	} `xml:"Strings"`
}

// systems is OVF.FindSystemIds' order: the VirtualSystem, or every one of the collection (nested ones too).
func (e *ovfEnvelope) systems() []ovfSystem {
	if e.System != nil {
		return []ovfSystem{*e.System}
	}
	var out []ovfSystem
	var walk func(c *ovfCollection)
	walk = func(c *ovfCollection) {
		out = append(out, c.Systems...)
		for i := range c.Collections {
			walk(&c.Collections[i])
		}
	}
	if e.Collection != nil {
		walk(e.Collection)
	}
	return out
}

// message is OVF.FindStringsMessage: a message's text from the Strings bundles by its id, else its own text.
func (e *ovfEnvelope) message(m ovfMsg) string {
	if m.MsgID != "" {
		for _, b := range e.Strings {
			for _, s := range b.Msgs {
				if s.MsgID == m.MsgID && s.Value != "" && s.Value != "Empty" {
					return strings.TrimSpace(s.Value)
				}
			}
		}
	}
	return strings.TrimSpace(m.Value)
}

// parseOVF reads a descriptor. A DOCTYPE is refused: no OVF has one, and it is how entity tricks start.
func parseOVF(data []byte) (*ovfEnvelope, error) {
	if len(data) > ovfMaxDescriptor {
		return nil, errors.New("the OVF descriptor is too large")
	}
	d := xml.NewDecoder(strings.NewReader(string(data)))
	d.Strict = true
	var env ovfEnvelope
	for {
		tok, err := d.Token()
		if err != nil {
			return nil, fmt.Errorf("the OVF descriptor is not valid: %w", err)
		}
		if dir, ok := tok.(xml.Directive); ok && strings.HasPrefix(strings.ToUpper(strings.TrimSpace(string(dir))), "DOCTYPE") {
			return nil, errors.New("the OVF descriptor is not valid: it has a DOCTYPE")
		}
		if start, ok := tok.(xml.StartElement); ok {
			if start.Name.Local != "Envelope" {
				return nil, errors.New("the OVF descriptor is not valid: it is not an OVF envelope")
			}
			if err := d.DecodeElement(&env, &start); err != nil {
				return nil, fmt.Errorf("the OVF descriptor is not valid: %w", err)
			}
			return &env, nil
		}
	}
}

// SafeHref is a file name of the package: relative, inside it (no "..", no absolute path, no URL).
func SafeHref(href string) bool {
	return href != "" && !strings.Contains(href, "\\") && !strings.Contains(href, ":") && !path.IsAbs(href) &&
		path.Clean(href) == href && !strings.HasPrefix(href, "../") && href != ".."
}

// ApplianceDisk is a disk of an appliance VM (AddResourceSettingData, resource types 17, 19 and 21).
type ApplianceDisk struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Href        string `json:"href"`        // the file inside the package
	Compression string `json:"compression"` // "" or "gzip" (the File's ovf:compression)
	Capacity    int64  `json:"capacity"`    // DiskSection capacity in bytes (0 when unknown)
	Device      string `json:"device"`      // the wanted userdevice ("" for the first free one)
	Bootable    bool   `json:"bootable"`
	ReadOnly    bool   `json:"read_only"`
}

// ApplianceVM is the VM ImportApplianceAction.CreateVm makes for a VirtualSystem, with its disks, interfaces and DVD drive.
type ApplianceVM struct {
	SysID           string            `json:"sys_id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Memory          int64             `json:"memory"`
	VCPUs           int64             `json:"vcpus"`
	VCPUsMax        int64             `json:"vcpus_max"`
	HVMBootPolicy   string            `json:"hvm_boot_policy"`
	HVMBootParams   map[string]string `json:"hvm_boot_params"`
	Platform        map[string]string `json:"platform"`
	NVRAM           map[string]string `json:"nvram"`
	Fields          map[string]string `json:"fields"` // the other VirtualSystemOtherConfigurationData (VM.UpdateFrom)
	HasVendorDevice bool              `json:"has_vendor_device"`
	Order           int               `json:"order"`
	StartDelay      int               `json:"start_delay"`
	ShutdownDelay   int               `json:"shutdown_delay"`
	Disks           []ApplianceDisk   `json:"disks"`
	MACs            []string          `json:"macs"`    // one per interface ("" for a generated MAC), in device order
	DVD             string            `json:"dvd"`     // the DVD drive's userdevice, "" for none
	Install         bool              `json:"install"` // the VM has an InstallSection (not run: see ImportAppliance)
}

// defaultImportMemory is VM.DEFAULT_MEM_MIN_IMG_IMPORT: a VM without a memory item gets 256 MiB.
const defaultImportMemory = 256 << 20

// ovfUnits is the multiplier of "byte * 2^20" (OVF.ComputeCapacity, CreateVm's memory units); ok is false for any other
// form.
func ovfUnits(units string) (float64, bool) {
	f := strings.FieldsFunc(units, func(r rune) bool { return r == '*' || r == '^' })
	if len(f) != 3 || !strings.EqualFold(strings.TrimSpace(f[0]), "byte") {
		return 1, false
	}
	b, err1 := strconv.ParseFloat(strings.TrimSpace(f[1]), 64)
	p, err2 := strconv.ParseFloat(strings.TrimSpace(f[2]), 64)
	if err1 != nil || err2 != nil {
		return 1, false
	}
	return math.Pow(b, p), true
}

// clampInt64 is a float product as an int64, at most MaxInt64 (CreateVm's "> long.MaxValue" cases).
func clampInt64(v float64) int64 {
	if v >= math.MaxInt64 || math.IsNaN(v) {
		return math.MaxInt64
	}
	if v < 0 {
		return 0
	}
	return int64(v)
}

func ovfUint(s string) uint64 { n, _ := strconv.ParseUint(strings.TrimSpace(s), 10, 64); return n }

// splitDict is StringExtensions.SplitToDictionary: "k1=v1<sep>k2=v2"; a part without "=" is its own key and value.
func splitDict(s string, sep string) map[string]string {
	out := map[string]string{}
	for _, part := range strings.Split(s, sep) {
		if part == "" {
			continue
		}
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			v = k
		}
		out[k] = v
	}
	return out
}

// hardware is OVF.FindVirtualHardwareSectionByAffinity(..., "xen"): a xen section first, then xen PV, then HVM, then
// any, then one without a type.
func hardware(vs ovfSystem) *ovfHardware {
	best, rank := -1, 99
	for i, h := range vs.Hardware {
		t := strings.ToLower(h.System.VirtualSystemType)
		r := 3
		switch {
		case t == "":
			r = 5
		case strings.HasPrefix(t, "xen"):
			r = 0
		case strings.HasSuffix(t, "xen:pv"):
			r = 1
		case strings.HasPrefix(t, "hvm") || strings.HasSuffix(t, "xen:hvm"):
			r = 2
		}
		if r < rank || best < 0 {
			best, rank = i, r
		}
	}
	if best < 0 {
		return &ovfHardware{}
	}
	return &vs.Hardware[best]
}

// systemName is OVF.FindSystemName (ImportWizard.FindVMName): the Name, else the id, the VirtualSystemIdentifier, the
// OS description, else the package's name.
func systemName(vs ovfSystem, h *ovfHardware, pkg string) string {
	for _, n := range []string{vs.Name, vs.ID, h.System.VirtualSystemIdentifier} {
		if n = strings.TrimSpace(n); n != "" {
			return n
		}
	}
	if len(vs.OS) > 0 && strings.TrimSpace(vs.OS[0].Description) != "" {
		return strings.TrimSpace(vs.OS[0].Description)
	}
	return pkg
}

func isDiskRASD(r *ovfRASD) bool { k := r.kind(); return k == 17 || k == 19 || k == 21 }

// notADisk is AddResourceSettingData's skip of the COM, floppy and ISO "disks".
func notADisk(r *ovfRASD) bool {
	c := strings.ToUpper(r.Caption)
	return strings.Contains(c, "COM") || strings.Contains(c, "FLOPPY") || strings.Contains(c, "ISO")
}

// diskOf is OVF.FindDiskReference: the DiskSection disk the item names (HostResource "ovf:/disk/<id>", else InstanceID).
// The whole id is compared: the C# Contains takes "vmdisk1" for "ovf:/disk/vmdisk10".
func (e *ovfEnvelope) diskOf(r *ovfRASD) *ovfDisk {
	for _, ref := range []string{first(r.HostResource), strings.TrimSpace(r.InstanceID)} {
		if i := strings.LastIndex(ref, "/"); i >= 0 {
			ref = ref[i+1:]
		}
		if ref == "" {
			continue
		}
		for i := range e.Disks {
			if e.Disks[i].DiskID == ref {
				return &e.Disks[i]
			}
		}
	}
	return nil
}

func (e *ovfEnvelope) file(id string) *ovfFile {
	for i := range e.Files {
		if e.Files[i].ID == id {
			return &e.Files[i]
		}
	}
	return nil
}

func first(s []string) string {
	if len(s) > 0 {
		return strings.TrimSpace(s[0])
	}
	return ""
}

// deviceConnections is ImportApplianceAction.SetDeviceConnections: the disks and DVD drives on each controller (IDE,
// SCSI, ...) in order get the next free userdevice (a DVD drive from 2 up), unless their Connection names one. It
// returns each item's Connection.
func deviceConnections(items []ovfRASD) []string {
	conns := make([]string, len(items))
	for i := range items {
		conns[i] = first(items[i].Connection)
	}
	controllers := []int{}
	for i := range items {
		if k := items[i].kind(); k >= 5 && k <= 9 {
			controllers = append(controllers, i)
		}
	}
	addr := func(r *ovfRASD) int {
		n, err := strconv.Atoi(strings.TrimSpace(r.Address))
		if err != nil {
			return -1
		}
		return n
	}
	sort.SliceStable(controllers, func(a, b int) bool { // OVF.CompareControllerRasd
		ra, rb := &items[controllers[a]], &items[controllers[b]]
		if aa, ab := addr(ra), addr(rb); aa >= 0 && ab >= 0 {
			return ra.kind()*10+aa < rb.kind()*10+ab
		}
		return ra.kind() < rb.kind()
	})
	used := make([]bool, 16)
	next := func(offset int) int { // FindNextAvailable
		for i := offset; i < len(used); i++ {
			if !used[i] {
				used[i] = true
				return i
			}
		}
		return 0
	}
	// OVF.FindConnectedItems: the DVD drives and hard disks whose Parent is the item, through Microsoft synthetic drives
	// (type 22, whose Address the disks under them take); each synthetic drive is followed once, so a hostile loop of
	// parents ends
	followed := map[int]bool{}
	var connectedTo func(instance string, synthetic *string) []int
	connectedTo = func(instance string, synthetic *string) []int {
		instance = strings.ReplaceAll(strings.TrimSpace(instance), `\`, "")
		var out []int
		if instance == "" {
			return out
		}
		for i := range items {
			r := &items[i]
			parent := strings.ReplaceAll(strings.TrimSpace(r.Parent), `\`, "")
			if parent == "" || !strings.Contains(parent, instance) {
				continue
			}
			name, caption := strings.ToLower(r.ElementName), strings.ToLower(r.Caption)
			switch r.kind() {
			case 15, 16:
				out = append(out, i)
			case 22:
				if strings.Contains(strings.ToLower(r.ResourceSubType), "synthetic") && !followed[i] {
					followed[i] = true
					out = append(out, connectedTo(r.InstanceID, &r.Address)...)
				}
			case 17, 19, 21:
				if strings.Contains(name, "hard disk") || strings.Contains(caption, "hard disk") || strings.HasPrefix(caption, "disk") {
					if synthetic != nil {
						r.Address = *synthetic
					}
					out = append(out, i)
				}
			}
		}
		slices.Sort(out)
		return slices.Compact(out)
	}
	for _, c := range controllers {
		connected := connectedTo(items[c].InstanceID, nil)
		sort.SliceStable(connected, func(a, b int) bool { // OVF.CompareConnectedDisks
			ra, rb := &items[connected[a]], &items[connected[b]]
			if ra.AddressOnParent != "" && rb.AddressOnParent != "" {
				return ra.AddressOnParent < rb.AddressOnParent
			}
			return addr(ra) < addr(rb)
		})
		for _, i := range connected {
			offset := 0
			if k := items[i].kind(); k == 15 || k == 16 {
				offset = 2
			}
			if conns[i] == "" {
				conns[i] = "device=" + strconv.Itoa(next(offset))
			} else if !strings.Contains(strings.ToLower(conns[i]), "device=") {
				conns[i] += ",device=" + strconv.Itoa(next(offset))
			}
		}
	}
	return conns
}

// bootable is ImportApplianceAction.SetIfDeviceIsBootable: the Connection's "bootable", else the disk's isBootable, else
// the first device of its controller, else userdevice 0.
func (e *ovfEnvelope) bootable(r *ovfRASD, conn string, dict map[string]string) bool {
	if strings.Contains(conn, "bootable") {
		b, _ := strconv.ParseBool(dict["bootable"])
		return b
	}
	if d := e.diskOf(r); d != nil && d.IsBootable {
		return true
	}
	k := r.kind()
	return (k == 21 || k == 5) && strings.TrimSpace(r.Address) == "0" ||
		(k == 17 || k == 19) && strings.TrimSpace(r.AddressOnParent) == "0" ||
		strings.Contains(conn, "device=0")
}

// vmSpec is ImportApplianceAction.CreateVm and AddResourceSettingData for one VirtualSystem.
func (e *ovfEnvelope) vmSpec(vs ovfSystem, pkg string) (ApplianceVM, error) {
	h := hardware(vs)
	vm := ApplianceVM{SysID: vs.ID, Name: systemName(vs, h, pkg), Description: strings.TrimSpace(h.System.Description), Install: len(vs.Install) > 0}
	if vm.Description == "" {
		vm.Description = messages.Get("DEFAULT_IMPORT_DESCRIPTION")
	}
	// MEMORY: MB unless the units say "byte * base^power"
	var memory float64
	var cpus, maxCPUs uint64
	for i := range h.Items {
		r := &h.Items[i]
		switch r.kind() {
		case 4:
			mult := float64(1 << 20)
			if strings.HasPrefix(strings.ToLower(strings.TrimSpace(r.AllocationUnits)), "byte") {
				if m, ok := ovfUnits(r.AllocationUnits); ok {
					mult = m
				}
			}
			memory += float64(ovfUint(r.VirtualQuantity)) * mult
		case 3: // CPU COUNT: VirtualQuantity cores; Limit the maximum (CA-361078: 100,000 or more means none)
			q := ovfUint(r.VirtualQuantity)
			cpus += q
			limit := uint64(100_000)
			if strings.TrimSpace(r.Limit) != "" {
				limit = ovfUint(r.Limit)
			}
			if limit >= 100_000 {
				maxCPUs += q
			} else {
				maxCPUs += limit
			}
		}
	}
	vm.Memory = clampInt64(memory)
	if vm.Memory == 0 {
		vm.Memory = defaultImportMemory
	}
	vm.VCPUs, vm.VCPUsMax = int64(min(max(cpus, 1), math.MaxInt64)), int64(min(max(maxCPUs, 1), math.MaxInt64))
	// XEN SPECIFIC CONFIGURATION INFORMATION
	if len(h.OtherConfig) == 0 {
		vm.HVMBootPolicy = "BIOS order" // DEFAULT should work for all of HVM type or 301
		vm.HVMBootParams = map[string]string{"order": "dc"}
		vm.Platform = map[string]string{"nx": "true", "acpi": "true", "apic": "true", "pae": "true", "stdvga": "0"}
	} else {
		vm.Fields = map[string]string{}
		for _, d := range h.OtherConfig {
			key := strings.ReplaceAll(d.Name, "-", "_")
			switch strings.ToLower(key) {
			case "hvm_boot_params": // "k1=v1;k2=v2", or an old OVF's plain boot order
				if strings.Contains(d.Value, "=") {
					vm.HVMBootParams = splitDict(d.Value, ";")
				} else {
					vm.HVMBootParams = map[string]string{"order": d.Value}
				}
			case "platform":
				vm.Platform = splitDict(d.Value, ";")
				if _, ok := vm.Platform["nx"]; !ok {
					vm.Platform["nx"] = "true"
				}
			case "nvram":
				vm.NVRAM = splitDict(d.Value, ";")
			case "vgpu", "pvssite": // shortcut: no vGPU and no PVS proxy (CreateVgpus, CreatePvsProxy); add them with those features
			case "vm_has_vendor_device":
				vm.HasVendorDevice, _ = strconv.ParseBool(strings.TrimSpace(d.Value))
			case "hvm_boot_policy":
				vm.HVMBootPolicy = d.Value
			default:
				vm.Fields[key] = d.Value
			}
		}
	}
	conns := deviceConnections(h.Items)
	for i := range h.Items {
		r := &h.Items[i]
		dict := splitDict(conns[i], ",")
		switch r.kind() {
		case 10: // Network: the MAC as 00:00:00:00:00:00 when it has no colons
			mac := strings.TrimSpace(r.Address)
			if mac != "" && !strings.Contains(mac, ":") {
				var parts []string
				for j := 0; j < len(mac); j += 2 {
					parts = append(parts, mac[j:min(j+2, len(mac))])
				}
				mac = strings.Join(parts, ":")
			}
			vm.MACs = append(vm.MACs, mac)
		case 15, 16: // CD/DVD drive: one per VM, empty (shortcut: no ISO of the package is imported; the C# app needs an ISO SR mapping it never sets here)
			if vm.DVD == "" {
				vm.DVD = cmp.Or(dict["device"], "3")
			}
		case 17, 19, 21:
			if notADisk(r) {
				continue
			}
			d := e.diskOf(r)
			if d == nil {
				continue
			}
			f := e.file(d.FileRef)
			if f == nil {
				continue
			}
			if !SafeHref(f.Href) {
				return vm, fmt.Errorf("the appliance names a file outside the package: %q", f.Href)
			}
			if ext := strings.ToLower(path.Ext(f.Href)); ext == ".mf" || ext == ".cert" || ext == ".ovf" { // never checked as data by the manifest
				return vm, fmt.Errorf("the disk %s of the appliance is not a disk image", f.Href)
			}
			comp := strings.ToLower(strings.TrimSpace(f.Compression))
			if comp == "identity" {
				comp = ""
			}
			if comp != "" && comp != "gzip" {
				return vm, fmt.Errorf("the compression method %s is not supported", f.Compression)
			}
			disk := ApplianceDisk{Name: strings.TrimSpace(r.ElementName), Description: strings.TrimSpace(r.Description), Href: f.Href,
				Compression: comp, Device: dict["device"], Bootable: e.bootable(r, conns[i], dict), ReadOnly: dict["mode"] == "r"}
			if c := ovfUint(d.Capacity); c > 0 {
				m, _ := ovfUnits(d.Units)
				disk.Capacity = clampInt64(float64(c) * m)
			}
			if disk.Name == "" {
				disk.Name = vm.Name + "_" + disk.Device
			}
			vm.Disks = append(vm.Disks, disk)
		}
	}
	return vm, nil
}

// Appliance is what the import dialog shows about a package (OvfRules.Info) and what the import makes of it.
type Appliance struct {
	Name          string        `json:"name"`
	Warnings      []string      `json:"warnings"`
	VMs           []ApplianceVM `json:"vms"`
	Disks         int           `json:"disks"`
	Interfaces    int           `json:"interfaces"`
	Eulas         []string      `json:"eulas"`
	HasManifest   bool          `json:"has_manifest"`
	RequiredSpace int64         `json:"required_space"` // the disks' capacities (OvfRules.RequiredSpace)
	// a signed package (a "<name>.cert"): its signer's certificate, or why it cannot be read
	HasSignature     bool             `json:"has_signature"`
	Certificate      *CertificateInfo `json:"certificate,omitempty"`
	CertificateError string           `json:"certificate_error,omitempty"`
}

var knownOVFVersions = []string{"0.9", "1.0", "1.0.0", "1.0.0a", "1.0.0.b", "1.0.0c", "1.0.0d", "1.0.1", "1.0.0e"}

var knownSystemTypes = []string{"xen-3.0-unknown", "xen-3.0-x32", "xen-3.0-x86", "xen-3.0-x64", "hvm-3.0-unknown", "hvm-3.0-x32",
	"hvm-3.0-x86", "hvm-3.0-x64", "hvm-3.0-hvm", "301", "vmx-4", "vmx-04", "vmx-6", "vmx-06", "vmx-7", "vmx-07", "DMTF:xen:pv",
	"DMTF:xen:hvm", "virtualbox-2.2"}

var knownFileExts = []string{".vhd", ".pvp", ".vmdk", ".mf", ".cert", ".xva", ".ovf", ".wim", ".vdi", ".sdi", ".iso", ".gz"}

// describeAppliance is OvfRules.Load past the package reading: OVF.Validate (shortcut: no XML schema check), the
// encryption refusal, the VMs, the EULAs and the space the disks need. has tells whether the package holds a file.
func describeAppliance(env *ovfEnvelope, name string, has func(string) bool, hasManifest bool) (*Appliance, error) {
	a := &Appliance{Name: name, HasManifest: hasManifest, Warnings: []string{}, Eulas: []string{}}
	switch {
	case env.Version == "":
		a.Warnings = append(a.Warnings, "The OVF version is not set; 1.0.0 is assumed.")
	case !slices.Contains(knownOVFVersions, env.Version):
		a.Warnings = append(a.Warnings, fmt.Sprintf("The OVF version %s is not known.", env.Version))
	}
	for _, f := range env.Files {
		ext := strings.ToLower(path.Ext(f.Href))
		if ext == ".mf" || ext == ".cert" {
			continue
		}
		if !SafeHref(f.Href) {
			return nil, fmt.Errorf("the appliance names a file outside the package: %q", f.Href)
		}
		if !has(f.Href) {
			return nil, fmt.Errorf("the file %s of the appliance is missing", f.Href)
		}
		if !slices.Contains(knownFileExts, ext) {
			a.Warnings = append(a.Warnings, fmt.Sprintf("The file %s has an extension that is not supported.", f.Href))
		}
	}
	if len(env.Encryption) > 0 {
		return nil, errors.New("Encrypted appliances cannot be imported in this version.")
	}
	systems := env.systems()
	if len(systems) == 0 {
		return nil, errors.New("The appliance has no virtual systems.")
	}
	eulas := append([]ovfEula{}, env.Eulas...)
	if env.Collection != nil {
		eulas = append(eulas, env.Collection.Eulas...)
	}
	for _, vs := range systems {
		if t := hardware(vs).System.VirtualSystemType; t != "" && !slices.Contains(knownSystemTypes, t) {
			a.Warnings = append(a.Warnings, fmt.Sprintf("The virtual hardware type %s is not known; the VM may not start.", t))
		}
		eulas = append(eulas, vs.Eulas...)
		vm, err := env.vmSpec(vs, name)
		if err != nil {
			return nil, err
		}
		if vm.Install {
			a.Warnings = append(a.Warnings, fmt.Sprintf("The post-install steps of %s are not run.", vm.Name))
		}
		a.Disks += len(vm.Disks)
		a.Interfaces += len(vm.MACs)
		for _, d := range vm.Disks {
			if a.RequiredSpace += d.Capacity; a.RequiredSpace < 0 {
				a.RequiredSpace = math.MaxInt64
			}
		}
		a.VMs = append(a.VMs, vm)
	}
	for _, e := range eulas {
		if t := env.message(e.License); t != "" {
			a.Eulas = append(a.Eulas, t)
		}
	}
	// the StartupSection: each VM's order and delays (ImportApplianceAction.Process)
	startups := append([]ovfStartup{}, env.Startup...)
	if env.Collection != nil {
		startups = append(startups, env.Collection.Startup...)
	}
	for _, st := range startups {
		for _, it := range st.Items {
			for i := range a.VMs {
				if a.VMs[i].SysID == it.ID {
					a.VMs[i].Order, a.VMs[i].StartDelay, a.VMs[i].ShutdownDelay = it.Order, it.StartDelay, it.StopDelay
				}
			}
		}
	}
	return a, nil
}
