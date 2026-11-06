package app

import (
	"fmt"
	"slices"
	"sort"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// A VM's GPU and USB pages of its Properties (XenAdmin/SettingsPanels/GpuEditPage.cs, Dialogs/AddVGPUDialog.cs,
// SettingsPanels/USBEditPage.cs, Dialogs/AttachUsbDialog.cs), shown as PropertiesDialog.cs decides.

// VMGPUPage is the GPU page: its rows and GpuEditPage.ShowHideWarnings' facts (rules.GPUPage; the page keeps Rows,
// Changed and Selected up to date).
type VMGPUPage struct {
	Shown bool               `json:"shown"`
	Rows  []VMGPURow         `json:"rows"`
	Input rules.GPUPageInput `json:"input"`
}

// VMGPURow is a row (VGpuDetailRow): an existing vGPU (its uuid) or one added on the page ("").
type VMGPURow struct {
	UUID   string        `json:"uuid"`
	Group  string        `json:"group"` // GPU group uuid
	Type   string        `json:"type"`  // vGPU type uuid
	Slot   string        `json:"slot"`  // PCI slot number
	Name   string        `json:"name"`  // Virtual GPU type
	PerGPU string        `json:"per_gpu"`
	VRAM   string        `json:"vram"`
	Facts  rules.VGPURow `json:"facts"`
}

func vmByUUID(s *xenapi.Session, uuid string) (xenapi.VMRef, xenapi.VMRecord, error) {
	ref, err := xenapi.VM.GetByUUID(s, uuid)
	if err != nil {
		return "", xenapi.VMRecord{}, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	return ref, vm, xapi.Sanitize(err)
}

// gpuRow is VGpuDetailRow.SetCells.
func gpuRow(d *ops.GPUData, uuid string, g xenapi.GPUGroupRef, t xenapi.VGPUTypeRef, device string) VMGPURow {
	rt := d.Type(t)
	_, known := d.Types[t]
	group := rules.GPUGroup{Name: rules.GPUGroupName(d.Groups[g].NameLabel)}
	if i := slices.IndexFunc(d.RuleGroups(), func(x rules.GPUGroup) bool { return x.Ref == d.Groups[g].UUID }); i >= 0 {
		group = d.RuleGroups()[i]
	}
	row := VMGPURow{UUID: uuid, Group: d.Groups[g].UUID, Type: rt.Ref, Slot: rules.VGPUSlot(device), Name: rt.Model,
		Facts: rules.VGPURow{Known: known, Passthrough: !known || rt.Passthrough(), Multiple: len(rt.Compatible) > 0, Experimental: rt.Experimental}}
	switch {
	case row.Facts.Passthrough && group.HasVGPU:
		row.Name = messages.Get("VGPU_PASSTHRU_TOSTRING")
	case row.Facts.Passthrough:
		row.Name = group.Name
	default:
		row.PerGPU = fmt.Sprint(rt.Capacity)
	}
	if rt.Framebuffer > 0 {
		row.VRAM = rules.MemorySizeText(float64(rt.Framebuffer), true)
	}
	return row
}

// ReadVMGPU reads a VM's GPU page.
func ReadVMGPU(c *xapi.Connection, vmUUID string) (VMGPUPage, error) {
	s := c.Session
	_, vm, err := vmByUUID(s, vmUUID)
	if err != nil {
		return VMGPUPage{}, err
	}
	d, err := ops.ReadGPUData(s)
	if err != nil {
		return VMGPUPage{}, err
	}
	ha, err := ops.PoolHAEnabled(s)
	if err != nil {
		return VMGPUPage{}, err
	}
	out := VMGPUPage{Rows: []VMGPURow{}, Shown: rules.GPUPageShown(vm.HVMBootPolicy != "", vm.Recommendations, d.RuleGroups()),
		Input: rules.GPUPageInput{Halted: vm.PowerState == xenapi.VMPowerStateHalted, HAEnabled: ha,
			HAPriority: vm.HaRestartPriority, GPUsAvailable: rules.GPUsAvailable(d.RuleGroups()), Rows: []rules.VGPURow{}}}
	for _, r := range vm.VGPUs {
		v, err := xenapi.VGPU.GetRecord(s, r)
		if err != nil { // a page that missed a vGPU would destroy it on OK
			return VMGPUPage{}, xapi.Sanitize(err)
		}
		row := gpuRow(d, v.UUID, v.GPUGroup, v.Type, v.Device)
		out.Rows = append(out.Rows, row)
		out.Input.Rows = append(out.Input.Rows, row.Facts)
	}
	out.Input.Original = len(out.Rows)
	return out, nil
}

// VGPUPick is an item of the Add virtual GPU dialog, with the row it adds.
type VGPUPick struct {
	rules.VGPUChoice
	Row VMGPURow `json:"row"`
}

// VGPUPicks is the Add virtual GPU dialog's list and the item selected first (-1: none).
type VGPUPicks struct {
	Items  []VGPUPick `json:"items"`
	Picked int        `json:"picked"`
}

// ReadVGPUChoices is AddVGPUDialog.PopulateComboBox for a VM whose page has vGPUs of these type uuids.
func ReadVGPUChoices(c *xapi.Connection, vmUUID string, existing []string) (VGPUPicks, error) {
	s := c.Session
	_, vm, err := vmByUUID(s, vmUUID)
	if err != nil {
		return VGPUPicks{}, err
	}
	d, err := ops.ReadGPUData(s)
	if err != nil {
		return VGPUPicks{}, err
	}
	choices, picked := d.Choices(vm, existing)
	out := VGPUPicks{Items: []VGPUPick{}, Picked: picked}
	for _, ch := range choices {
		p := VGPUPick{VGPUChoice: ch}
		if ch.Type != "" {
			g, _ := xenapi.GPUGroup.GetByUUID(s, ch.Group)
			t, _ := xenapi.VGPUType.GetByUUID(s, ch.Type)
			p.Row = gpuRow(d, "", g, t, "")
		}
		out.Items = append(out.Items, p)
	}
	return out, nil
}

// VMUSBPage is the USB page (USBEditPage): the VM's USB devices, and whether one can be attached (else why not).
type VMUSBPage struct {
	Shown     bool   `json:"shown"`
	Title     string `json:"title"`
	Table     Table  `json:"table"` // keys: VUSB uuids
	CanAttach bool   `json:"can_attach"`
	Warning   string `json:"warning"`
	Confirm   string `json:"confirm"` // Detach's question and title
	ConfirmT  string `json:"confirm_title"`
}

// ReadVMUSB reads a VM's USB page; haPriority is the Start Options page's choice (USBEditPage.SelectedPriority).
func ReadVMUSB(c *xapi.Connection, vmUUID, haPriority string) (VMUSBPage, error) {
	s := c.Session
	_, vm, err := vmByUUID(s, vmUUID)
	if err != nil {
		return VMUSBPage{}, err
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return VMUSBPage{}, xapi.Sanitize(err)
	}
	pusbs, err := xenapi.PUSB.GetAllRecords(s)
	if err != nil {
		return VMUSBPage{}, xapi.Sanitize(err)
	}
	out := VMUSBPage{Title: messages.Get("USB_DEVICES"), Confirm: messages.Get("ACTION_VUSB_DETACH_CONFIRM"), ConfirmT: messages.Get("ACTION_VUSB_DETACH"),
		Table: Table{Headers: []string{"Location", "Description", "Attached"}, Rows: [][]string{}, Keys: []string{}},
		Shown: rules.USBPageShown(vm.HVMBootPolicy != "", vm.IsATemplate, ops.FeatureForbidden(hosts, "restrict_usb_passthrough"), len(pusbs) > 0)}
	ha, err := ops.PoolHAEnabled(s)
	if err != nil {
		return VMUSBPage{}, err
	}
	out.CanAttach, out.Warning = rules.USBAttach(vm.PowerState == xenapi.VMPowerStateHalted, ha, haPriority)
	for _, r := range vm.VUSBs {
		v, err := xenapi.VUSB.GetRecord(s, r)
		if err != nil {
			continue
		}
		loc, desc := "", ""
		if ps, err := xenapi.USBGroup.GetPUSBs(s, v.USBGroup); err == nil && len(ps) > 0 { // VMUsbRow: the group's first device
			p := pusbs[ps[0]]
			loc, desc = p.Path, rules.USBDescription(p.Description, p.VendorID, p.ProductID)
		}
		out.Table.Rows = append(out.Table.Rows, []string{loc, desc, map[bool]string{true: messages.Get("YES"), false: messages.Get("NO")}[v.CurrentlyAttached]})
		out.Table.Keys = append(out.Table.Keys, v.UUID)
	}
	return out, nil
}

// USBAttachChoices is the Attach USB dialog (AttachUsbDialog.BuildList): the free devices of each server the VM can
// start on, and whether the VM's first attach moves it to the upstream QEMU (labelWarningLine3).
type USBAttachChoices struct {
	Hosts   []USBAttachHost `json:"hosts"`
	Empty   string          `json:"empty"`
	Upgrade bool            `json:"upgrade"`
}

// USBAttachHost is a server of the dialog and its devices (PUSB uuid, "path description").
type USBAttachHost struct {
	Name    string        `json:"name"`
	Devices []NetworkPick `json:"devices"`
}

// ReadUSBAttach reads the Attach USB dialog of a VM.
func ReadUSBAttach(c *xapi.Connection, vmUUID string) (USBAttachChoices, error) {
	s := c.Session
	ref, vm, err := vmByUUID(s, vmUUID)
	if err != nil {
		return USBAttachChoices{}, err
	}
	out := USBAttachChoices{Hosts: []USBAttachHost{}, Empty: messages.Get("DIALOG_ATTACH_USB_NO_DEVICES_AVAILABLE"),
		Upgrade: vm.Platform["device-model"] != "qemu-upstream-compat"}
	possible, err := xenapi.VM.GetPossibleHosts(s, ref)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	pusbs, err := xenapi.PUSB.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	type host struct {
		ref  xenapi.HostRef
		name string
	}
	var hs []host
	for _, h := range possible {
		name, err := xenapi.Host.GetNameLabel(s, h)
		if err == nil {
			hs = append(hs, host{h, name})
		}
	}
	sort.Slice(hs, func(i, j int) bool {
		return rules.HostLess(hs[i].name, string(hs[i].ref), hs[j].name, string(hs[j].ref))
	})
	for _, h := range hs {
		var list []xenapi.PUSBRecord
		for _, p := range pusbs {
			if p.Host == h.ref && rules.USBAttachable(p.PassthroughEnabled, ops.USBAttached(s, p)) {
				list = append(list, p)
			}
		}
		if len(list) == 0 { // a server is shown only with devices
			continue
		}
		sort.Slice(list, func(i, j int) bool { return list[i].Path < list[j].Path })
		row := USBAttachHost{Name: h.name}
		for _, p := range list {
			row.Devices = append(row.Devices, NetworkPick{UUID: p.UUID,
				Name: messages.Format(messages.Get("STRING_SPACE_STRING"), p.Path, rules.USBDescription(p.Description, p.VendorID, p.ProductID))})
		}
		out.Hosts = append(out.Hosts, row)
	}
	return out, nil
}
