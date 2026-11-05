package app

import (
	"encoding/xml"
	"sort"
	"strconv"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// VMProperties is what the VM's Properties dialog shows (PropertiesDialog: GeneralPage, CpuPage, HomeServerPage).
// shortcut: General has the name and description only (no tags, folder or custom fields yet); no Boot Options page.
type VMProperties struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Halted      bool   `json:"halted"`
	Running     bool   `json:"running"`

	// CPU (CpuPage)
	Hotplug        bool         `json:"hotplug"` // VM.SupportsVcpuHotplug: not a Windows guest
	VCPUsMax       int          `json:"vcpus_max"`
	VCPUsAtStartup int          `json:"vcpus_at_startup"`
	CoresPerSocket int64        `json:"cores_per_socket"`
	Priority       int          `json:"priority"` // 0-8, log4 of the vCPU weight
	MinVCPUs       int          `json:"min_vcpus"`
	VCPUChoices    []VCPUChoice `json:"vcpu_choices"` // vCPU counts with at least one valid topology
	MaxPhysical    int          `json:"max_physical"` // the most CPUs of a server (CpuPage's warning)

	// Home Server (HomeServerPage / AffinityPicker)
	HomeServer  string       `json:"home_server"` // host uuid, "" for none
	Servers     []ServerPick `json:"servers"`
	CanHaveNone bool         `json:"can_have_none"` // VmRules.CanHaveNoHomeServer
}

// VCPUChoice is one vCPU count and its topologies.
type VCPUChoice struct {
	VCPUs      int64            `json:"vcpus"`
	Topologies []rules.Topology `json:"topologies"`
}

// ServerPick is a server of the Home Server page.
type ServerPick struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Live bool   `json:"live"`
}

// restriction reads a "vcpus-max"/"vcpus-min" restriction of the VM's recommendations (VM.GetMaxRestrictionValue).
// shortcut: the value of a matching template with the same reference_label (CP-44766) is not looked up.
func restriction(recommendations, field, attr string, def int) int {
	var doc struct {
		R []struct {
			Field string `xml:"field,attr"`
			Max   string `xml:"max,attr"`
			Min   string `xml:"min,attr"`
			Value string `xml:"value,attr"`
		} `xml:"restriction"`
	}
	if xml.Unmarshal([]byte(recommendations), &doc) != nil {
		return def
	}
	for _, r := range doc.R {
		if r.Field != field {
			continue
		}
		v := map[string]string{"max": r.Max, "min": r.Min, "value": r.Value}[attr]
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// linuxDistros are VM.LinuxDistros: a guest agent reporting one of them is not Windows, whatever the platform flags.
var linuxDistros = map[string]bool{"debian": true, "rhel": true, "fedora": true, "centos": true, "scientific": true, "oracle": true,
	"sles": true, "lsb": true, "boot2docker": true, "freebsd": true, "ubuntu": true, "neokylin": true, "gooroom": true, "rocky": true, "alpine": true}

// isWindows is VM.IsWindows: not when the guest agent reports a Linux distribution (or NetScaler), else an HVM guest
// with viridian on. shortcut: the guest metrics of the last booted record are not read.
func isWindows(s *xenapi.Session, vm xenapi.VMRecord) bool {
	if gm, err := xenapi.VMGuestMetrics.GetRecord(s, vm.GuestMetrics); err == nil {
		if linuxDistros[strings.ToLower(gm.OsVersion["distro"])] || strings.Contains(strings.ToLower(gm.OsVersion["uname"]), "netscaler") {
			return false
		}
	}
	return vm.HVMBootPolicy != "" && vm.Platform["viridian"] == "true"
}

// coresPerSocket is Host.CoresPerSocket: cpu_count / socket_count (0 when unknown).
func coresPerSocket(h xenapi.HostRecord) int64 {
	cpus, _ := strconv.ParseInt(h.CPUInfo["cpu_count"], 10, 64)
	sockets, _ := strconv.ParseInt(h.CPUInfo["socket_count"], 10, 64)
	if cpus > 0 && sockets > 0 {
		return cpus / sockets
	}
	return 0
}

// ReadVMProperties reads the VM for the Properties dialog.
func ReadVMProperties(c *xapi.Connection, vmUUID string) (VMProperties, error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return VMProperties{}, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return VMProperties{}, xapi.Sanitize(err)
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return VMProperties{}, xapi.Sanitize(err)
	}
	weight := 256 // VM.GetVcpuWeight
	if w, ok := vm.VCPUsParams["weight"]; ok {
		if n, err := strconv.Atoi(w); err != nil {
			weight = 65536
		} else if n > 0 {
			weight = n
		} else {
			weight = 1
		}
	}
	cores := int64(1) // VM.GetCoresPerSocket
	if n, err := strconv.ParseInt(vm.Platform["cores-per-socket"], 10, 64); err == nil {
		cores = n
	}
	p := VMProperties{Name: vm.NameLabel, Description: vm.NameDescription, Halted: vm.PowerState == xenapi.VMPowerStateHalted,
		Running: vm.PowerState == xenapi.VMPowerStateRunning, Hotplug: !isWindows(s, vm) && !ops.VCPUHotplugRestricted(hosts),
		VCPUsMax: max(vm.VCPUsMax, 1), VCPUsAtStartup: max(vm.VCPUsAtStartup, 1), CoresPerSocket: cores, Priority: rules.Priority(weight),
		MinVCPUs: restriction(vm.Recommendations, "vcpus-min", "min", 1), Servers: []ServerPick{}}
	// VM.MaxCoresPerSocket: the home server's, else the largest of the pool
	home := vmHome(s, vm)
	var maxCores int64
	for r, h := range hosts {
		if c := coresPerSocket(h); (home != "" && r == home) || (home == "" && c > maxCores) {
			maxCores = c
		}
		n, _ := strconv.Atoi(h.CPUInfo["cpu_count"])
		p.MaxPhysical = max(p.MaxPhysical, n)
		live := false
		if m, err := xenapi.HostMetrics.GetLive(s, h.Metrics); err == nil {
			live = m
		}
		p.Servers = append(p.Servers, ServerPick{UUID: h.UUID, Name: h.NameLabel, Live: live})
		if r == vm.Affinity {
			p.HomeServer = h.UUID
		}
	}
	sort.Slice(p.Servers, func(i, j int) bool { return p.Servers[i].Name < p.Servers[j].Name })
	current := int64(p.VCPUsMax)
	if !p.Hotplug {
		current = int64(p.VCPUsAtStartup)
	}
	limit := max(int64(restriction(vm.Recommendations, "vcpus-max", "max", 16)), current) // VM.MaxVCPUsAllowed
	for n := int64(1); n <= limit; n++ {
		if t := rules.ValidTopologies(n, maxCores); len(t) > 0 || n == current {
			p.VCPUChoices = append(p.VCPUChoices, VCPUChoice{VCPUs: n, Topologies: t})
		}
	}
	// VmRules.CanHaveNoHomeServer: shared storage that can hold disks and reaches every server, or no home server now
	p.CanHaveNone = p.HomeServer == ""
	if srs, err := xenapi.SR.GetAllRecords(s); err == nil {
		for _, sr := range srs {
			if sr.Shared && sr.ContentType != "iso" && !broken(s, sr, len(hosts)) {
				p.CanHaveNone = true
			}
		}
	}
	return p, nil
}
