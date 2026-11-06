package rules

import (
	"slices"
	"sort"
	"strconv"
	"strings"

	"xcpng-center/internal/messages"
)

// A VM's GPU and USB pages (XenAdmin/SettingsPanels/GpuEditPage.cs, Dialogs/AddVGPUDialog.cs,
// Controls/ComboBoxes/VgpuComboBox.cs, SettingsPanels/USBEditPage.cs, Dialogs/AttachUsbDialog.cs,
// Dialogs/PropertiesDialog.cs) and the SR-IOV networks (Wizards/NewNetworkWizard_Pages/NetWTypeSelect.cs,
// NetWSriovDetails.cs, TabPages/NICPage.cs).

// CanHaveGPU is VM.CanHaveGpu: an HVM guest whose allow-gpu-passthrough restriction is absent or not zero.
func CanHaveGPU(hvm bool, recommendations string) bool {
	return hvm && Restriction(recommendations, "allow-gpu-passthrough", "value", 1) != 0
}

// CanHaveVGPU is VM.CanHaveVGpu: CanHaveGPU and an allow-vgpu restriction absent or not zero.
func CanHaveVGPU(hvm bool, recommendations string) bool {
	return CanHaveGPU(hvm, recommendations) && Restriction(recommendations, "allow-vgpu", "value", 1) != 0
}

// GPUGroup is what the GPU page reads of a GPU group (GPU_group).
type GPUGroup struct {
	Ref       string
	Name      string // GPUGroupName
	PGPUs     int
	Supported []string // supported_VGPU_types
	Enabled   []string // enabled_VGPU_types
	HasVGPU   bool     // GPU_group.HasVGpu: one of its GPUs supports a type that is not pass-through
}

// GPUGroupName is GPU_group.Name: the name without its "Group of " prefix.
func GPUGroupName(label string) string { return strings.TrimPrefix(label, "Group of ") }

// GPUsAvailable is Helpers.GpusAvailable: a GPU group with GPUs and supported vGPU types.
func GPUsAvailable(groups []GPUGroup) bool {
	return slices.ContainsFunc(groups, func(g GPUGroup) bool { return g.PGPUs > 0 && len(g.Supported) > 0 })
}

// GPUPageShown is PropertiesDialog's test for a VM's GPU page: VM.CanHaveGpu and Helpers.GpusAvailable.
func GPUPageShown(hvm bool, recommendations string, groups []GPUGroup) bool {
	return CanHaveGPU(hvm, recommendations) && GPUsAvailable(groups)
}

// VGPUTypeText is VGPU_type.Description: "Pass-through whole GPU", or the model with its vGPUs per GPU (and whether
// several fit in one VM).
func VGPUTypeText(t VGPUType) string {
	if t.Passthrough() {
		return messages.Get("VGPU_PASSTHRU_TOSTRING")
	}
	if len(t.Compatible) > 0 {
		return messages.Format(messages.Get("VGPU_DESCRIPTION_MULTIPLE_VGPU_SUPPORT"), t.Model, t.Capacity)
	}
	return messages.Format(messages.Get("VGPU_DESCRIPTION"), t.Model, t.Capacity)
}

// VGPUChoice is an item of the Add virtual GPU dialog's list (VgpuComboBox's GpuTuple).
type VGPUChoice struct {
	Group   string `json:"group"`
	Type    string `json:"type"` // "" for a group's header
	Text    string `json:"text"`
	Header  bool   `json:"header"`  // the group's bold line, not selectable
	Sub     bool   `json:"sub"`     // a type under its group's header (indented)
	Enabled bool   `json:"enabled"` // a type the GPUs do not allow is grey and not selectable
}

// Selectable is VgpuComboBox.IsItemNonSelectable's opposite.
func (c VGPUChoice) Selectable() bool { return !c.Header && c.Enabled }

// VGPUChoices is AddVGPUDialog.PopulateComboBox: per GPU group (by name), the types it supports that every vGPU the VM
// already has accepts beside it (compatible_types_in_vm), under the group's header; with vGPUs restricted by the
// licence or by the VM's template only the groups offering whole-GPU pass-through. existing are the types of the VM's
// vGPUs (on the page, unsaved ones included). picked is the item selected first (-1: none).
//
// The C# pass-through item has no type and its Add does nothing (AddVGPUDialog.buttonAdd_Click); here it carries the
// group's pass-through type, so it can be added (GpuAssignAction then creates it without a type when vGPUs are
// restricted).
func VGPUChoices(groups []GPUGroup, types map[string]VGPUType, existing []string, passthroughOnly bool) (out []VGPUChoice, picked int) {
	out, picked = []VGPUChoice{}, -1
	var list []GPUGroup
	for _, g := range groups {
		if g.PGPUs > 0 && len(g.Supported) > 0 {
			list = append(list, g)
		}
	}
	sort.SliceStable(list, func(i, j int) bool { return HostLess(list[i].Name, list[i].Ref, list[j].Name, list[j].Ref) })
	for _, g := range list {
		if passthroughOnly {
			for _, r := range g.Supported { // GPU_group.HasPassthrough
				if t, ok := types[r]; ok && t.Passthrough() { // not offered when the GPUs do not allow it (C#: always shown)
					out = append(out, VGPUChoice{Group: g.Ref, Type: r, Text: g.Name, Enabled: slices.Contains(g.Enabled, r)})
					break
				}
			}
			continue
		}
		var common []VGPUType
		for _, r := range g.Supported {
			t, ok := types[r]
			if !ok {
				continue
			}
			fits := true
			for _, e := range existing {
				if et, ok := types[e]; ok { // an unknown type does not narrow the choice
					fits = fits && slices.Contains(et.Compatible, r)
				}
			}
			if fits {
				common = append(common, t)
			}
		}
		SortVGPUTypes(common)
		if g.HasVGPU && len(common) > 0 {
			out = append(out, VGPUChoice{Group: g.Ref, Text: g.Name, Header: true})
		}
		for _, t := range common {
			c := VGPUChoice{Group: g.Ref, Type: t.Ref, Text: g.Name, Sub: g.HasVGPU, Enabled: slices.Contains(g.Enabled, t.Ref)}
			if c.Sub {
				c.Text = VGPUTypeText(t)
			}
			out = append(out, c)
		}
		if len(common) == 1 && slices.Contains(g.Enabled, common[0].Ref) && picked < 0 {
			picked = len(out) - 1
		}
	}
	return out, picked
}

// VGPURow is what the GPU page's state needs of a row's vGPU type (VGPU.IsPassthrough, IsExperimental and
// compatible_types_in_vm); Known is false when the type is gone, which counts as pass-through.
type VGPURow struct {
	Known        bool `json:"known"`
	Passthrough  bool `json:"passthrough"`
	Multiple     bool `json:"multiple"` // compatible_types_in_vm is not empty
	Experimental bool `json:"experimental"`
}

// GPUPageInput is GpuEditPage.ShowHideWarnings' facts.
type GPUPageInput struct {
	Halted        bool      `json:"halted"`
	HAEnabled     bool      `json:"ha_enabled"`
	HAPriority    string    `json:"ha_priority"` // the Start Options page's choice (PropertiesDialog: SelectedPriority)
	GPUsAvailable bool      `json:"gpus_available"`
	Rows          []VGPURow `json:"rows"`
	Original      int       `json:"original"` // the VM's vGPUs when the page opened
	Changed       bool      `json:"changed"`  // HasChanged
	Selected      bool      `json:"selected"`
}

// GPUPageState is the GPU page's buttons and warnings (the warnings table's lines, top to bottom).
type GPUPageState struct {
	CanAdd       bool   `json:"can_add"`
	CanDelete    bool   `json:"can_delete"`
	Experimental bool   `json:"experimental"`
	RDP          bool   `json:"rdp"`
	NeedGPU      bool   `json:"need_gpu"` // the driver and GPU-at-boot lines
	StopVM       bool   `json:"stop_vm"`
	HA           bool   `json:"ha"`
	Multiple     string `json:"multiple"` // NEWVMWIZARD_VGPUPAGE_(MULTIPLE|SINGLE)_VGPU_INFO, "" when hidden
}

// GPUPage is GpuEditPage.ShowHideWarnings: a VM that is not halted, or that HA restarts, cannot change its vGPUs; a
// vGPU is added when GPUs are available and every vGPU of the VM accepts another.
func GPUPage(in GPUPageInput) GPUPageState {
	var st GPUPageState
	n := len(in.Rows)
	st.Experimental = slices.ContainsFunc(in.Rows, func(r VGPURow) bool { return r.Known && r.Experimental })
	if !in.Halted {
		st.StopVM = true
		return st
	}
	if in.HAEnabled && HAIsRestart(in.HAPriority) {
		st.HA = true
		return st
	}
	multiple := !slices.ContainsFunc(in.Rows, func(r VGPURow) bool { return !r.Known || !r.Multiple })
	st.CanAdd, st.CanDelete = in.GPUsAvailable && multiple, in.Selected
	if n > 0 {
		st.Multiple = messages.Get("NEWVMWIZARD_VGPUPAGE_SINGLE_VGPU_INFO")
		if multiple {
			st.Multiple = messages.Get("NEWVMWIZARD_VGPUPAGE_MULTIPLE_VGPU_INFO")
		}
	}
	st.RDP = in.Changed && in.Original == 0 && slices.ContainsFunc(in.Rows, func(r VGPURow) bool { return !r.Known || r.Passthrough })
	st.NeedGPU = n > 0
	return st
}

// VGPUSlot is VGpuDetailRow's PCI slot: the device number shifted by 11 (xapi keeps 0-20 for vGPUs, guest slots are
// 11-31; QuebecOrGreater, every XCP-ng 8.x), "-" for none.
func VGPUSlot(device string) string {
	if n, err := strconv.Atoi(device); err == nil {
		return strconv.Itoa(n + 11)
	}
	return messages.Get("HYPHEN")
}

// USBPageShown is PropertiesDialog's test for a VM's USB page: an HVM VM (not a template), USB pass-through not
// restricted by the licence, and a USB device on a server of the pool.
func USBPageShown(hvm, template, forbidden, anyDevice bool) bool {
	return hvm && !template && !forbidden && anyDevice
}

// USBAttach is USBEditPage.ShowHideWarnings: a device is attached to a halted VM that HA does not restart; else why not.
func USBAttach(halted, haEnabled bool, haPriority string) (ok bool, why string) {
	switch {
	case !halted:
		return false, messages.Get("USB_ATTACH_NOT_ALLOWED_NOT_HALTED")
	case haEnabled && HAIsRestart(haPriority):
		return false, messages.Get("USB_ATTACH_NOT_ALLOWED_HA")
	}
	return true, ""
}

// USBAttachable is AttachUsbDialog.BuildList's test for a device: pass-through on, no VM using it.
func USBAttachable(passthrough, attached bool) bool { return passthrough && !attached }

// USBDescription is PUSB.Description: its description, else its vendor and product ids.
func USBDescription(description, vendor, product string) string {
	if description != "" {
		return description
	}
	return messages.Format(messages.Get("USB_DEVICE_VENDOR_PRODUCT"), vendor, product)
}

// SriovCapable is PIF.SriovCapable.
func SriovCapable(capabilities []string) bool { return slices.Contains(capabilities, "sriov") }

// SriovCandidate is NetWSriovDetails.PopulateHostNicList's test for a NIC: physical, not a bond interface, SR-IOV
// capable, SR-IOV not enabled on it yet (not an SR-IOV physical PIF).
func SriovCandidate(physical, bond, capable, enabled bool) bool {
	return physical && !bond && capable && !enabled
}

// SriovOption is NetWTypeSelect.PopulatePage's SR-IOV Network choice (always shown: every XCP-ng is Kolkata or
// greater): enabled, else the warning under it. experimental is Host.SriovNetworkDisabled on a server, forbidden
// Host.RestrictSriovNetwork; anyCapable is Pool.HasSriovNic, anyFree a physical capable NIC without SR-IOV yet.
func SriovOption(experimental, forbidden, anyCapable, anyFree bool) (ok bool, why string) {
	name := messages.Get("NETWORK_SRIOV")
	switch {
	case experimental:
		return false, messages.Format(messages.Get("FEATURE_EXPERIMENTAL"), name)
	case forbidden:
		return false, messages.Format(messages.Get("FEATURE_DISABLED"), name)
	case anyCapable && !anyFree:
		return false, messages.Get("NICS_ARE_SRIOV_ENABLED")
	case !anyCapable:
		return false, messages.Get("SRIOV_NEED_NICSUPPORT")
	}
	return true, ""
}

// SriovNIC is NICPage's SR-IOV Capable cell: No; Yes but no SR-IOV network yet; the server needs a reboot; the
// logical PIF is unplugged; else the virtual functions left (remaining < 0: unreadable, plain Yes).
func SriovNIC(capable, enabled, needsReboot, logicalAttached bool, remaining int) string {
	switch {
	case !capable:
		return messages.Get("NO")
	case !enabled:
		return messages.Get("SRIOV_NETWORK_SHOULD_BE_CREATED")
	case needsReboot:
		return messages.Get("HOST_NEEDS_REBOOT_ENABLE_SRIOV")
	case !logicalAttached:
		return messages.Get("SRIOV_LOGICAL_PIF_UNPLUGGED")
	case remaining < 0:
		return messages.Get("YES")
	}
	return messages.Format(messages.Get("REMAINING_VFS"), remaining)
}
