package app

import (
	"sort"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// The Memory tab's bars (MemoryView): a VM's used memory inside its dynamic range (MemoryRules.VmBar, VMShinyBar), and
// each server's memory split into Xen and the control domain, each VM and free memory (MemoryRules.HostSegments,
// HostShinyBar).

// MemorySegment is a part of a server's bar.
type MemorySegment struct {
	Name  string `json:"name"`
	Bytes int64  `json:"bytes"`
	Text  string `json:"text"`
}

// MemoryBars are the bars of a VM ("vms", "templates") or of the servers of a server or pool.
type MemoryBars struct {
	Used       float64      `json:"used"` // fractions of the static maximum
	DynamicMin float64      `json:"dynamic_min"`
	DynamicMax float64      `json:"dynamic_max"`
	UsedText   string       `json:"used_text"`
	Hosts      []HostMemBar `json:"hosts"`
}

// HostMemBar is one server's bar.
type HostMemBar struct {
	Name     string          `json:"name"`
	Segments []MemorySegment `json:"segments"`
}

// ReadMemoryBars reads them.
func ReadMemoryBars(c *xapi.Connection, kind, uuid string) (MemoryBars, error) {
	s := c.Session
	out := MemoryBars{Hosts: []HostMemBar{}}
	switch kind {
	case "vms", "templates":
		ref, err := xenapi.VM.GetByUUID(s, uuid)
		if err != nil {
			return out, xapi.Sanitize(err)
		}
		vm, err := xenapi.VM.GetRecord(s, ref)
		if err != nil {
			return out, xapi.Sanitize(err)
		}
		top := float64(max(1, vm.MemoryStaticMax))
		out.DynamicMin, out.DynamicMax = float64(vm.MemoryDynamicMin)/top, float64(vm.MemoryDynamicMax)/top
		if vm.PowerState == xenapi.VMPowerStateRunning || vm.PowerState == xenapi.VMPowerStatePaused {
			if m, err := xenapi.VMMetrics.GetRecord(s, vm.Metrics); err == nil {
				out.Used = min(1, float64(m.MemoryActual)/top)
				out.UsedText = "In use: " + rules.MemoryText(int64(m.MemoryActual))
			}
		}
		if out.UsedText == "" {
			out.UsedText = "The VM is not running"
		}
		return out, nil
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	for _, h := range hosts {
		if kind == "hosts" && h.UUID != uuid {
			continue
		}
		hm, err := xenapi.HostMetrics.GetRecord(s, h.Metrics)
		if err != nil {
			continue
		}
		xen := int64(h.MemoryOverhead)
		var vms []MemorySegment
		for _, r := range h.ResidentVMs {
			vm, err := xenapi.VM.GetRecord(s, r)
			if err != nil {
				continue
			}
			actual := int64(0)
			if m, err := xenapi.VMMetrics.GetRecord(s, vm.Metrics); err == nil {
				actual = int64(m.MemoryActual)
			}
			if vm.IsControlDomain {
				xen += actual + int64(vm.MemoryOverhead)
				continue
			}
			b := actual + int64(vm.MemoryOverhead)
			vms = append(vms, MemorySegment{Name: vm.NameLabel, Bytes: b, Text: vm.NameLabel + ": " + rules.MemoryText(b)})
		}
		sort.Slice(vms, func(i, j int) bool { return vms[i].Name < vms[j].Name })
		name := "XCP-ng (Xen and control domain)"
		free := max(0, int64(hm.MemoryFree))
		segs := append([]MemorySegment{{Name: name, Bytes: xen, Text: name + ": " + rules.MemoryText(xen)}}, vms...)
		segs = append(segs, MemorySegment{Name: "Free", Bytes: free, Text: "Free: " + rules.MemoryText(free)})
		out.Hosts = append(out.Hosts, HostMemBar{Name: h.NameLabel, Segments: segs})
	}
	sort.Slice(out.Hosts, func(i, j int) bool { return out.Hosts[i].Name < out.Hosts[j].Name })
	return out, nil
}
