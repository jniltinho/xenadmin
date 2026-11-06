package app

import (
	"fmt"
	"sort"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// The GPU tab of a pool or server and the USB tab of a server (XenAdmin.Desktop/HardwareViews.cs, HardwareRules.cs).

// GPUTab is GpuPage: the placement policy (when the pool can run vGPUs) and a row per GPUs of the same settings.
type GPUTab struct {
	Policy      string   `json:"policy"`       // "Placement policy: ..."; empty when not shown
	PolicyValue string   `json:"policy_value"` // depth_first, breadth_first or unknown (mixed)
	NoGPU       string   `json:"no_gpu"`       // the text when no GPU does vGPU
	Groups      []GPURow `json:"groups"`
	Density     string   `json:"density"` // the placement choices (PoolGpuEditPage)
	Performance string   `json:"performance"`
}

// GPURow is GpuRow: the GPUs by server, their vGPU types, and the Edit button (GpuConfiguration).
type GPURow struct {
	Hosts   []GPUHost    `json:"hosts"`
	Types   []GPUTypeRow `json:"types"`
	PGPUs   []string     `json:"pgpus"` // uuids
	CanEdit bool         `json:"can_edit"`
	Edit    string       `json:"edit"`   // the button's text
	Title   string       `json:"title"`  // GpuConfiguration's heading
	Rubric  string       `json:"rubric"` // and its text
}

// GPUHost is a server's GPUs of a row.
type GPUHost struct {
	Label string   `json:"label"`
	GPUs  []string `json:"gpus"`
}

// GPUTypeRow is HardwareRules.VgpuTypeRow.
type GPUTypeRow struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Max     string `json:"max"`
	VRAM    string `json:"vram"`
	Enabled bool   `json:"enabled"`
	InUse   bool   `json:"in_use"`
}

// ReadGPU reads the GPU tab of a pool ("pools") or a server.
func ReadGPU(c *xapi.Connection, kind, uuid string) (GPUTab, error) {
	s := c.Session
	out := GPUTab{Groups: []GPURow{}, Density: messages.Get("GPU_PLACEMENT_POLICY_MAX_DENSITY_DESCRIPTION"),
		Performance: messages.Get("GPU_PLACEMENT_POLICY_MAX_PERFORMANCE_DESCRIPTION")}
	out.NoGPU = messages.Get("NO_GPU_ON_HOST")
	if kind == "pools" {
		out.NoGPU = messages.Get("NO_GPU_IN_POOL")
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	pgpus, err := xenapi.PGPU.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	types, err := xenapi.VGPUType.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	vgpus, err := xenapi.VGPU.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	groups, err := xenapi.GPUGroup.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	noVGPU := ops.FeatureForbidden(hosts, "restrict_vgpu")
	name := func(p xenapi.PGPURecord) string { // PGPU.Name: the PCI device's name
		if d, err := xenapi.PCI.GetDeviceName(s, p.PCI); err == nil {
			return d
		}
		return p.UUID
	}
	hasVGPU := func(p xenapi.PGPURecord) bool { // PGPU.HasVGpu
		for _, t := range p.SupportedVGPUTypes {
			if tt, ok := types[t]; ok && tt.ModelName != "passthrough" {
				return true
			}
		}
		return false
	}

	// GpuPlacementPolicyPanel, shown when the pool can run vGPUs (Helpers.VGpuCapability)
	capable := false
	for _, p := range pgpus {
		capable = capable || hasVGPU(p)
	}
	if capable && !noVGPU {
		algos := map[xenapi.AllocationAlgorithm]bool{}
		for _, g := range groups {
			algos[g.AllocationAlgorithm] = true
		}
		out.PolicyValue, out.Policy = "unknown", messages.Get("GPU_PLACEMENT_POLICY_MIXED_DESCRIPTION")
		if len(algos) == 1 {
			for a := range algos {
				out.PolicyValue = string(a)
			}
			switch out.PolicyValue {
			case "breadth_first":
				out.Policy = out.Performance
			case "depth_first":
				out.Policy = out.Density
			}
		}
		out.Policy = messages.Format(messages.Get("GPU_PLACEMENT_POLICY_DESCRIPTION"), out.Policy)
	}

	// HardwareRules.Pgpus: the GPUs that support vGPU types, of the pool or of this server, by server then name
	type gpu struct {
		ref        xenapi.PGPURef
		rec        xenapi.PGPURecord
		host, name string
	}
	var list []gpu
	for ref, p := range pgpus {
		h, ok := hosts[p.Host]
		if len(p.SupportedVGPUTypes) == 0 || kind != "pools" && (!ok || h.UUID != uuid) {
			continue
		}
		list = append(list, gpu{ref, p, h.NameLabel, name(p)})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].rec.Host != list[j].rec.Host {
			return rules.HostLess(list[i].host, string(list[i].rec.Host), list[j].host, string(list[j].rec.Host))
		}
		return list[i].name < list[j].name
	})
	key := func(g gpu) string { // HardwareRules.Groups: one row per name, enabled and supported types
		refs := func(rs []xenapi.VGPUTypeRef) string {
			ss := make([]string, len(rs))
			for i, r := range rs {
				ss[i] = string(r)
			}
			sort.Strings(ss)
			return strings.Join(ss, ",")
		}
		return g.name + "|" + refs(g.rec.EnabledVGPUTypes) + "|" + refs(g.rec.SupportedVGPUTypes)
	}
	var order []string
	byKey := map[string][]gpu{}
	for _, g := range list {
		k := key(g)
		if _, ok := byKey[k]; !ok {
			order = append(order, k)
		}
		byKey[k] = append(byKey[k], g)
	}
	for _, k := range order {
		gs := byKey[k]
		first := gs[0].rec
		row := GPURow{Hosts: []GPUHost{}, Types: []GPUTypeRow{}, CanEdit: !noVGPU && hasVGPU(first),
			Edit: strings.ReplaceAll(messages.Get("GPU_EDIT_ALLOWED_TYPES_SINGLE"), "&", "")}
		var last xenapi.HostRef
		for i, g := range gs { // by server (the list is sorted by server)
			if i == 0 || g.rec.Host != last {
				row.Hosts, last = append(row.Hosts, GPUHost{Label: messages.Format(messages.Get("GPU_ON_HOST_LABEL"), g.host)}), g.rec.Host
			}
			h := &row.Hosts[len(row.Hosts)-1]
			h.GPUs = append(h.GPUs, fmt.Sprintf("%s   (%d vGPUs running)", g.name, len(g.rec.ResidentVGPUs)))
			row.PGPUs = append(row.PGPUs, g.rec.UUID)
		}
		row.Title, row.Rubric = messages.Format(messages.Get("GPU_GROUP_NAME_AND_NO_OF_GPUS_ONE"), gs[0].name), messages.Get("GPU_RUBRIC_PLEASE_SELECT_WHICH_GPU_ONE")
		if len(gs) > 1 {
			row.Edit = strings.ReplaceAll(messages.Get("GPU_EDIT_ALLOWED_TYPES_MULTIPLE"), "&", "")
			row.Title = messages.Format(messages.Get("GPU_GROUP_NAME_AND_NO_OF_GPUS"), gs[0].name, len(gs))
			row.Rubric = messages.Get("GPU_RUBRIC_PLEASE_SELECT_WHICH_GPU_MULTIPLE")
		}
		// GpuConfiguration.PopulateGridView: each supported type, enabled or not, in use on these GPUs
		var sorted []rules.VGPUType
		for _, r := range first.SupportedVGPUTypes {
			t, ok := types[r]
			if !ok {
				continue
			}
			capacity := int64(0) // VGPU_type.Capacity: from the first GPU that supports it
			if len(t.SupportedOnPGPUs) > 0 {
				capacity = int64(pgpus[t.SupportedOnPGPUs[0]].SupportedVGPUMaxCapacities[r])
			}
			sorted = append(sorted, rules.VGPUType{Ref: string(r), Model: t.ModelName, Capacity: capacity,
				Resolution: int64(t.MaxResolutionX) * int64(t.MaxResolutionY), Framebuffer: int64(t.FramebufferSize)})
		}
		rules.SortVGPUTypes(sorted)
		for _, t := range sorted {
			r := xenapi.VGPUTypeRef(t.Ref)
			tr := GPUTypeRow{UUID: types[r].UUID, Name: t.Model}
			if t.Passthrough() {
				tr.Name = messages.Get("VGPU_PASSTHRU_TOSTRING")
			}
			if m := first.SupportedVGPUMaxCapacities[r]; m > 0 {
				tr.Max = fmt.Sprint(m)
			}
			if t.Framebuffer > 0 {
				tr.VRAM = rules.MemorySizeText(float64(t.Framebuffer), true)
			}
			for _, e := range first.EnabledVGPUTypes {
				tr.Enabled = tr.Enabled || e == r
			}
			for _, g := range gs {
				for _, v := range g.rec.ResidentVGPUs {
					tr.InUse = tr.InUse || vgpus[v].Type == r
				}
			}
			row.Types = append(row.Types, tr)
		}
		out.Groups = append(out.Groups, row)
	}
	return out, nil
}

// USBTab is UsbPage: the server's USB devices.
type USBTab struct {
	Table   Table       `json:"table"`
	Devices []USBDevice `json:"devices"`
	Title   string      `json:"title"`
}

// USBDevice is a row's facts: the pass-through button (HardwareRules.PassthroughButton) and its question.
type USBDevice struct {
	UUID        string `json:"uuid"`
	Passthrough bool   `json:"passthrough"`
	Button      string `json:"button"`
	CanToggle   bool   `json:"can_toggle"`
	Question    string `json:"question"`
	Yes         string `json:"yes"`
}

// ReadUSB reads the USB tab of a server.
func ReadUSB(c *xapi.Connection, hostUUID string) (USBTab, error) {
	s := c.Session
	out := USBTab{Table: Table{Headers: []string{"Location", "Description", "Pass-through", "Virtual machine"}, Rows: [][]string{}, Keys: []string{}},
		Devices: []USBDevice{},
		Title:   messages.Get("USB_DEVICES")}
	host, err := xenapi.Host.GetByUUID(s, hostUUID)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	pusbs, err := xenapi.PUSB.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	var list []xenapi.PUSBRecord
	for _, p := range pusbs {
		if p.Host == host {
			list = append(list, p)
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Path < list[j].Path })
	for _, p := range list {
		vm := ops.USBAttachedVM(s, p)
		desc := p.Description // PUSB.Description
		if desc == "" {
			desc = messages.Format(messages.Get("USB_DEVICE_VENDOR_PRODUCT"), p.VendorID, p.ProductID)
		}
		state := messages.Get("DISABLED")
		d := USBDevice{UUID: p.UUID, Passthrough: p.PassthroughEnabled, CanToggle: true,
			Button:   strings.ReplaceAll(messages.Get("USBLIST_ENABLE_PASSTHROUGH_HOTKEY"), "&", ""),
			Question: messages.Get("DIALOG_USB_USAGE_ENABLE_PASSTHROUGH"),
			Yes:      strings.ReplaceAll(messages.Get("DIALOG_USB_USAGE_OKBUTTON_ENABLE"), "&", "")}
		if p.PassthroughEnabled {
			state = messages.Get("ENABLED")
			d.Button = strings.ReplaceAll(messages.Get("USBLIST_DISABLE_PASSTHROUGH_HOTKEY"), "&", "")
			d.CanToggle = vm == "" && !ops.USBAttached(s, p)
			d.Question = messages.Get("DIALOG_USB_USAGE_DISABLE_PASSTHROUGH")
			d.Yes = strings.ReplaceAll(messages.Get("DIALOG_USB_USAGE_OKBUTTON_DISABLE"), "&", "")
		}
		out.Table.Rows = append(out.Table.Rows, []string{p.Path, desc, state, vm})
		out.Table.Keys = append(out.Table.Keys, p.UUID)
		out.Devices = append(out.Devices, d)
	}
	return out, nil
}
