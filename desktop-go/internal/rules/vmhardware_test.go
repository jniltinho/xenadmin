package rules

import (
	"reflect"
	"testing"
)

func TestCanHaveGPU(t *testing.T) {
	no := `<restrictions><restriction field="allow-gpu-passthrough" value="0"/></restrictions>`
	noV := `<restrictions><restriction field="allow-vgpu" value="0"/></restrictions>`
	cases := []struct {
		hvm       bool
		recs      string
		gpu, vgpu bool
	}{
		{true, "", true, true}, // absent restrictions allow both
		{false, "", false, false},
		{true, no, false, false},
		{true, noV, true, false},
	}
	for _, c := range cases {
		if g, v := CanHaveGPU(c.hvm, c.recs), CanHaveVGPU(c.hvm, c.recs); g != c.gpu || v != c.vgpu {
			t.Errorf("%v %q: gpu %v vgpu %v, want %v %v", c.hvm, c.recs, g, v, c.gpu, c.vgpu)
		}
	}
}

func TestGPUPageShown(t *testing.T) {
	emulated := []GPUGroup{{Ref: "g", PGPUs: 1}} // the lab: a GPU with no supported vGPU type
	real := []GPUGroup{{Ref: "g", PGPUs: 1, Supported: []string{"pt"}}}
	empty := []GPUGroup{{Ref: "g", Supported: []string{"pt"}}}
	if GPUPageShown(true, "", emulated) || GPUPageShown(true, "", empty) || GPUPageShown(true, "", nil) {
		t.Error("shown without a GPU that supports a vGPU type")
	}
	if !GPUPageShown(true, "", real) || GPUPageShown(false, "", real) {
		t.Error("HVM with GPUs: shown; PV: hidden")
	}
	if GPUGroupName("Group of NVIDIA GRID K1") != "NVIDIA GRID K1" {
		t.Error("GPUGroupName")
	}
}

// a pool of a vGPU-capable group (K1: K100 and K140Q fit together, K180Q alone, pass-through) and a pass-through only
// group
func gpuPool() ([]GPUGroup, map[string]VGPUType) {
	types := map[string]VGPUType{
		"pt":    {Ref: "pt", Model: "passthrough"},
		"k100":  {Ref: "k100", Model: "GRID K100", Capacity: 8, Resolution: 1920 * 1200, Compatible: []string{"k100", "k140"}},
		"k140":  {Ref: "k140", Model: "GRID K140Q", Capacity: 4, Resolution: 2560 * 1600, Compatible: []string{"k100", "k140"}, Experimental: true},
		"k180":  {Ref: "k180", Model: "GRID K180Q", Capacity: 1, Resolution: 2560 * 1600},
		"pt2":   {Ref: "pt2", Model: "passthrough"},
		"other": {Ref: "other", Model: "X"},
	}
	groups := []GPUGroup{
		{Ref: "b", Name: "Radeon", PGPUs: 1, Supported: []string{"pt2"}, Enabled: []string{"pt2"}},
		{Ref: "a", Name: "GRID K1", PGPUs: 2, Supported: []string{"k100", "k140", "k180", "pt"}, Enabled: []string{"k100", "k180", "pt"}, HasVGPU: true},
		{Ref: "c", Name: "Empty", PGPUs: 0, Supported: []string{"other"}},
	}
	return groups, types
}

func TestVGPUChoices(t *testing.T) {
	groups, types := gpuPool()
	got, picked := VGPUChoices(groups, types, nil, false)
	want := []VGPUChoice{
		{Group: "a", Text: "GRID K1", Header: true},
		{Group: "a", Type: "pt", Text: "Pass-through whole GPU", Sub: true, Enabled: true},
		{Group: "a", Type: "k180", Text: "GRID K180Q virtual GPU (1 per GPU)", Sub: true, Enabled: true},
		{Group: "a", Type: "k140", Text: "GRID K140Q virtual GPU (4 per GPU, multiple vGPU support)", Sub: true}, // not allowed on the GPUs
		{Group: "a", Type: "k100", Text: "GRID K100 virtual GPU (8 per GPU, multiple vGPU support)", Sub: true, Enabled: true},
		{Group: "b", Type: "pt2", Text: "Radeon", Enabled: true}, // a pass-through only group: no header, its name
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("choices\n got %+v\nwant %+v", got, want)
	}
	if picked != 5 { // the only type of Radeon (C#: the first item when a group has one allowed type)
		t.Errorf("picked %d", picked)
	}
	if got[0].Selectable() || got[3].Selectable() || !got[1].Selectable() {
		t.Error("headers and types not allowed are not selectable")
	}

	// a K100 already on the VM: only K100 and K140Q fit beside it, in the K1 group only
	got, _ = VGPUChoices(groups, types, []string{"k100"}, false)
	var refs []string
	for _, c := range got {
		refs = append(refs, c.Group+":"+c.Type)
	}
	if !reflect.DeepEqual(refs, []string{"a:", "a:k140", "a:k100"}) {
		t.Errorf("beside a K100: %v", refs)
	}
	// a K180Q (no multiple vGPU support): nothing fits
	if got, _ := VGPUChoices(groups, types, []string{"k180"}, false); len(got) != 0 {
		t.Errorf("beside a K180Q: %+v", got)
	}
	// vGPUs restricted (licence or template): the pass-through of each group that has one
	got, picked = VGPUChoices(groups, types, nil, true)
	want = []VGPUChoice{{Group: "a", Type: "pt", Text: "GRID K1", Enabled: true}, {Group: "b", Type: "pt2", Text: "Radeon", Enabled: true}}
	if !reflect.DeepEqual(got, want) || picked != -1 {
		t.Errorf("pass-through only: %+v %d", got, picked)
	}
	groups[0].Enabled = nil // Radeon's pass-through not allowed on its GPU
	if got, _ = VGPUChoices(groups, types, nil, true); got[1].Selectable() {
		t.Errorf("a pass-through type the GPUs do not allow: %+v", got[1])
	}
}

func TestGPUPage(t *testing.T) {
	k100 := VGPURow{Known: true, Multiple: true}
	k180 := VGPURow{Known: true}
	pt := VGPURow{Known: true, Passthrough: true}
	exp := VGPURow{Known: true, Multiple: true, Experimental: true}
	base := GPUPageInput{Halted: true, GPUsAvailable: true}

	st := GPUPage(base)
	if !st.CanAdd || st.CanDelete || st.Multiple != "" || st.NeedGPU || st.RDP {
		t.Errorf("no vGPU: %+v", st)
	}
	in := base
	in.Rows, in.Changed, in.Selected = []VGPURow{k100, exp}, true, true
	st = GPUPage(in)
	if !st.CanAdd || !st.CanDelete || !st.Experimental || !st.NeedGPU || st.RDP || st.Multiple != "The selected virtual GPU type supports multiple instances." {
		t.Errorf("two multi-vGPU types: %+v", st)
	}
	in.Rows = []VGPURow{k180}
	if st = GPUPage(in); st.CanAdd || st.Multiple != "The selected virtual GPU type does not support multiple instances." {
		t.Errorf("a single vGPU type: %+v", st)
	}
	in.Rows = []VGPURow{{}} // a type that is gone: pass-through, no multiple
	if st = GPUPage(in); st.CanAdd || !st.RDP {
		t.Errorf("unknown type: %+v", st)
	}
	in.Rows = []VGPURow{pt}
	if st = GPUPage(in); !st.RDP {
		t.Error("a new pass-through on a VM without GPU warns about Remote Desktop")
	}
	in.Original = 1
	if st = GPUPage(in); st.RDP {
		t.Error("no Remote Desktop warning when the VM had a GPU")
	}
	in = base
	in.GPUsAvailable = false
	if GPUPage(in).CanAdd {
		t.Error("no GPU available: no Add")
	}
	in = base
	in.Halted, in.Rows, in.Selected = false, []VGPURow{exp}, true
	if st = GPUPage(in); st.CanAdd || st.CanDelete || !st.StopVM || !st.Experimental || st.NeedGPU || st.Multiple != "" {
		t.Errorf("running: %+v", st)
	}
	in = base
	in.HAEnabled, in.HAPriority, in.Selected = true, "restart", true
	if st = GPUPage(in); st.CanAdd || st.CanDelete || !st.HA || st.StopVM {
		t.Errorf("HA restart: %+v", st)
	}
	in.HAPriority = "best-effort"
	if st = GPUPage(in); !st.CanAdd || st.HA {
		t.Errorf("HA best effort: %+v", st)
	}
}

func TestVGPUSlot(t *testing.T) {
	for in, want := range map[string]string{"0": "11", "3": "14", "": "-", "x": "-"} {
		if got := VGPUSlot(in); got != want {
			t.Errorf("%q: %q", in, got)
		}
	}
}

func TestUSBRules(t *testing.T) {
	if !USBPageShown(true, false, false, true) {
		t.Error("an HVM VM with a USB device in the pool shows the page")
	}
	for _, c := range [][4]bool{{false, false, false, true}, {true, true, false, true}, {true, false, true, true}, {true, false, false, false}} {
		if USBPageShown(c[0], c[1], c[2], c[3]) {
			t.Errorf("shown for %v", c)
		}
	}
	if ok, why := USBAttach(true, false, "restart"); !ok || why != "" {
		t.Error("halted, HA off: attach")
	}
	if ok, why := USBAttach(false, false, ""); ok || why != "USB devices can only be attached when the VM is shut down." {
		t.Error(why)
	}
	if ok, why := USBAttach(true, true, "restart"); ok || why != "Cannot attach USB devices because the VM is currently protected by HA." {
		t.Error(why)
	}
	if ok, _ := USBAttach(true, true, ""); !ok {
		t.Error("HA on, not restarted: attach")
	}
	if !USBAttachable(true, false) || USBAttachable(true, true) || USBAttachable(false, false) {
		t.Error("USBAttachable")
	}
	if USBDescription("", "0781", "5567") != "Vendor: 0781; Product: 5567" {
		t.Error("USBDescription")
	}
	if USBDescription("Flash", "a", "b") != "Flash" {
		t.Error("USBDescription keeps the description")
	}
}

func TestSriovRules(t *testing.T) {
	if !SriovCapable([]string{"fcoe", "sriov"}) || SriovCapable([]string{"fcoe"}) || SriovCapable(nil) {
		t.Error("SriovCapable")
	}
	if !SriovCandidate(true, false, true, false) || SriovCandidate(false, false, true, false) || SriovCandidate(true, true, true, false) ||
		SriovCandidate(true, false, false, false) || SriovCandidate(true, false, true, true) {
		t.Error("SriovCandidate")
	}
	cases := []struct {
		exp, forbidden, capable, free bool
		ok                            bool
		why                           string
	}{
		{false, false, true, true, true, ""},
		{true, false, true, true, false, "Feature 'SR-IOV' must be enabled on all hosts in the pool before use. See release notes."},
		{false, true, true, true, false, "'SR-IOV' feature is disabled due to license restrictions on the server."},
		{false, false, true, false, false, "SR-IOV is already enabled on all the SR-IOV capable NICs."},
		{false, false, false, false, false, "SR-IOV networks require an SR-IOV capable NIC."}, // the lab
	}
	for _, c := range cases {
		if ok, why := SriovOption(c.exp, c.forbidden, c.capable, c.free); ok != c.ok || why != c.why {
			t.Errorf("%+v: %v %q", c, ok, why)
		}
	}
	nic := [][2]string{
		{SriovNIC(false, false, false, false, 0), "No"},
		{SriovNIC(true, false, false, false, 0), "Yes (To enable SR-IOV create an SR-IOV network on this NIC)"},
		{SriovNIC(true, true, true, true, 3), "Yes (Host needs reboot)"},
		{SriovNIC(true, true, false, false, 3), "Yes (SR-IOV logical PIF is unplugged)"},
		{SriovNIC(true, true, false, true, 3), "Yes (3 VFs remaining)"},
		{SriovNIC(true, true, false, true, -1), "Yes"},
	}
	for _, c := range nic {
		if c[0] != c[1] {
			t.Errorf("SriovNIC %q, want %q", c[0], c[1])
		}
	}
}
