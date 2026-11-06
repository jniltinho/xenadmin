package ops

import (
	"cmp"
	"errors"
	"maps"
	"slices"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// A VM's GPU and USB pages' writes: GpuAssignAction (XenModel/Actions/VM/GpuAssignAction.cs), CreateVUSBAction and
// DeleteVUSBAction (XenModel/Actions/USB), with the rules of GpuEditPage, AddVGPUDialog, USBEditPage and
// AttachUsbDialog (rules/vmhardware.go). Each checks the server's current state again.

// GPUData is what the GPU page reads of a pool: its GPU groups, vGPU types and GPUs, by uuid for the rules.
type GPUData struct {
	Groups     map[xenapi.GPUGroupRef]xenapi.GPUGroupRecord
	Types      map[xenapi.VGPUTypeRef]xenapi.VGPUTypeRecord
	PGPUs      map[xenapi.PGPURef]xenapi.PGPURecord
	NoVGPU     bool // Host.RestrictVgpu on a server (Helpers.FeatureForbidden)
	ruleTypes  map[string]rules.VGPUType
	ruleGroups []rules.GPUGroup
}

// ReadGPUData reads the GPU groups, vGPU types and GPUs of the connection.
func ReadGPUData(s *xenapi.Session) (*GPUData, error) {
	d := &GPUData{ruleTypes: map[string]rules.VGPUType{}}
	var err error
	if d.Groups, err = xenapi.GPUGroup.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	if d.Types, err = xenapi.VGPUType.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	if d.PGPUs, err = xenapi.PGPU.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	d.NoVGPU = featureForbidden(hosts, "restrict_vgpu")
	for r := range d.Types {
		d.ruleTypes[d.Types[r].UUID] = d.Type(r)
	}
	for _, g := range d.Groups {
		rg := rules.GPUGroup{Ref: g.UUID, Name: rules.GPUGroupName(g.NameLabel), PGPUs: len(g.PGPUs)}
		for _, t := range g.SupportedVGPUTypes {
			rg.Supported = append(rg.Supported, d.Types[t].UUID)
		}
		for _, t := range g.EnabledVGPUTypes {
			rg.Enabled = append(rg.Enabled, d.Types[t].UUID)
		}
		for _, p := range g.PGPUs { // GPU_group.HasVGpu: PGPU.HasVGpu
			for _, t := range d.PGPUs[p].SupportedVGPUTypes {
				tt, ok := d.Types[t]
				rg.HasVGPU = rg.HasVGPU || ok && tt.ModelName != "passthrough"
			}
		}
		d.ruleGroups = append(d.ruleGroups, rg)
	}
	slices.SortFunc(d.ruleGroups, func(a, b rules.GPUGroup) int { return cmp.Compare(a.Ref, b.Ref) }) // map order is random
	return d, nil
}

// Type is a vGPU type for the rules (VGPU_type.Capacity: from the first GPU that supports it); an unknown ref gives
// a zero type.
func (d *GPUData) Type(r xenapi.VGPUTypeRef) rules.VGPUType {
	t, ok := d.Types[r]
	if !ok {
		return rules.VGPUType{}
	}
	out := rules.VGPUType{Ref: t.UUID, Model: t.ModelName, Resolution: int64(t.MaxResolutionX) * int64(t.MaxResolutionY),
		Framebuffer: int64(t.FramebufferSize), Experimental: t.Experimental}
	if len(t.SupportedOnPGPUs) > 0 {
		out.Capacity = int64(d.PGPUs[t.SupportedOnPGPUs[0]].SupportedVGPUMaxCapacities[r])
	}
	for _, c := range t.CompatibleTypesInVM {
		if ct, ok := d.Types[c]; ok {
			out.Compatible = append(out.Compatible, ct.UUID)
		}
	}
	return out
}

// RuleGroups are the GPU groups for the rules (by uuid).
func (d *GPUData) RuleGroups() []rules.GPUGroup { return d.ruleGroups }

// Choices is AddVGPUDialog's list for a VM whose vGPUs (on the page) have these type uuids.
func (d *GPUData) Choices(vm xenapi.VMRecord, existing []string) ([]rules.VGPUChoice, int) {
	passthroughOnly := d.NoVGPU || !rules.CanHaveVGPU(vm.HVMBootPolicy != "", vm.Recommendations)
	return rules.VGPUChoices(d.ruleGroups, d.ruleTypes, existing, passthroughOnly)
}

// NewVGPU is a vGPU added on the GPU page: its GPU group and vGPU type (uuids).
type NewVGPU struct {
	Group string `json:"group"`
	Type  string `json:"type"`
}

// AssignVGPUs is GpuAssignAction: the vGPUs the page removed (uuids of those it read; one added elsewhere since is
// kept, as the C# action works on the dialog's copy of the VM) are destroyed, then the new ones are created (device
// "0": the first free slot). The VM must be halted and not restarted by HA (GpuEditPage.ShowHideWarnings), and each
// new vGPU one the Add dialog offers beside the vGPUs that stay (AddVGPUDialog).
func AssignVGPUs(c *xapi.Connection, vmUUID string, removed []string, add []NewVGPU) error {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if err := vmEditable(s, vm, true); err != nil {
		return err
	}
	d, err := ReadGPUData(s)
	if err != nil {
		return err
	}
	var remove []xenapi.VGPURef // a removed vGPU already gone is left out
	var types []string          // the types beside which a new vGPU must fit
	for _, r := range vm.VGPUs {
		v, err := xenapi.VGPU.GetRecord(s, r)
		if err != nil {
			return xapi.Sanitize(err)
		}
		if slices.Contains(removed, v.UUID) {
			remove = append(remove, r)
		} else {
			types = append(types, d.Types[v.Type].UUID)
		}
	}
	if len(add) > 0 && !rules.CanHaveGPU(vm.HVMBootPolicy != "", vm.Recommendations) {
		return errors.New("this VM cannot have a GPU")
	}
	groupRef := map[string]xenapi.GPUGroupRef{}
	for r, g := range d.Groups {
		groupRef[g.UUID] = r
	}
	typeRef := map[string]xenapi.VGPUTypeRef{}
	for r, t := range d.Types {
		typeRef[t.UUID] = r
	}
	for _, a := range add {
		choices, _ := d.Choices(vm, types)
		if !slices.ContainsFunc(choices, func(ch rules.VGPUChoice) bool { return ch.Selectable() && ch.Group == a.Group && ch.Type == a.Type }) {
			return errors.New("this virtual GPU type cannot be added (any more) beside the VM's other virtual GPUs")
		}
		types = append(types, a.Type)
	}
	var calls []string // RbacMethod drops "async_"
	if len(remove) > 0 {
		calls = append(calls, "VGPU.destroy")
	}
	if len(add) > 0 {
		calls = append(calls, "VGPU.create")
	}
	if err := permitted(s, calls...); err != nil {
		return err
	}
	for _, r := range remove {
		if err := xenapi.VGPU.Destroy(s, r); err != nil {
			return xapi.Sanitize(err)
		}
	}
	for _, a := range add {
		if d.NoVGPU { // GpuAssignAction.AddGpu: no type when vGPUs are restricted
			err = wait(s)(xenapi.VGPU.AsyncCreate5(s, ref, groupRef[a.Group], "0", map[string]string{}))
		} else {
			err = wait(s)(xenapi.VGPU.AsyncCreate(s, ref, groupRef[a.Group], "0", map[string]string{}, typeRef[a.Type]))
		}
		if err != nil {
			return xapi.Sanitize(err)
		}
	}
	return nil
}

// vmEditable is the blocker of GpuEditPage.ShowHideWarnings (gpu, its labelStopVM and labelHA) and of
// USBEditPage.ShowHideWarnings: a VM that is not halted, or that HA restarts.
func vmEditable(s *xenapi.Session, vm xenapi.VMRecord, gpu bool) error {
	ha, err := PoolHAEnabled(s)
	if err != nil {
		return err
	}
	ok, why := rules.USBAttach(vm.PowerState == xenapi.VMPowerStateHalted, ha, vm.HaRestartPriority)
	switch {
	case ok:
		return nil
	case gpu && vm.PowerState != xenapi.VMPowerStateHalted:
		return errors.New("The virtual GPUs can only be changed when the VM is shut down.")
	case gpu:
		return errors.New("The virtual GPUs cannot be changed because the VM is currently protected by HA.")
	}
	return errors.New(why)
}

// PoolHAEnabled is Pool.ha_enabled.
func PoolHAEnabled(s *xenapi.Session) (bool, error) {
	pools, err := xenapi.Pool.GetAllRecords(s)
	for _, p := range pools {
		return p.HaEnabled, nil
	}
	return false, xapi.Sanitize(err)
}

// usingUpstreamQemu is VM.UsingUpstreamQemu.
func usingUpstreamQemu(vm xenapi.VMRecord) bool {
	return vm.Platform["device-model"] == "qemu-upstream-compat"
}

// AttachUSB is AttachUsbDialog's Attach then CreateVUSBAction: a VM not on the upstream QEMU is moved to it
// (platform device-model, irreversible), then the device's USB group is given to the VM. The device must still be
// free with pass-through on, the VM halted and not restarted by HA.
func AttachUSB(c *xapi.Connection, vmUUID, pusbUUID string) error {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	calls := []string{"VUSB.create"}
	if !usingUpstreamQemu(vm) {
		calls = append(calls, "VM.set_platform")
	}
	if err := permitted(s, calls...); err != nil {
		return err
	}
	if vm.IsATemplate || vm.HVMBootPolicy == "" {
		return errors.New("USB devices are passed through to HVM VMs only")
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if featureForbidden(hosts, "restrict_usb_passthrough") {
		return errors.New("USB pass-through is restricted by the licence")
	}
	if err := vmEditable(s, vm, false); err != nil {
		return err
	}
	pref, err := xenapi.PUSB.GetByUUID(s, pusbUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	p, err := xenapi.PUSB.GetRecord(s, pref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vusbs, err := xenapi.USBGroup.GetVUSBs(s, p.USBGroup)
	if err != nil {
		return xapi.Sanitize(err)
	}
	possible, err := xenapi.VM.GetPossibleHosts(s, ref) // AttachUsbDialog lists the devices of these servers only
	if err != nil {
		return xapi.Sanitize(err)
	}
	if !rules.USBAttachable(p.PassthroughEnabled, len(vusbs) > 0) || !slices.Contains(possible, p.Host) {
		return errors.New("this USB device is no longer available to the VM")
	}
	if !usingUpstreamQemu(vm) {
		platform := maps.Clone(vm.Platform)
		if platform == nil {
			platform = map[string]string{}
		}
		platform["device-model"] = "qemu-upstream-compat"
		if err := xenapi.VM.SetPlatform(s, ref, platform); err != nil {
			return xapi.Sanitize(err)
		}
	}
	_, err = xenapi.VUSB.Create(s, ref, p.USBGroup, map[string]string{})
	return xapi.Sanitize(err)
}

// DetachUSB is USBEditPage's Detach then DeleteVUSBAction: an attached device is unplugged first when it can be,
// then the VUSB is destroyed even if the unplug failed.
func DetachUSB(c *xapi.Connection, vmUUID, vusbUUID string) error {
	s := c.Session
	ref, err := xenapi.VUSB.GetByUUID(s, vusbUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	v, err := xenapi.VUSB.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if u, err := xenapi.VM.GetUUID(s, v.VM); err != nil || u != vmUUID {
		return errors.New("this USB device is not attached to the VM")
	}
	calls := []string{"VUSB.destroy"} // RbacMethod drops "async_"
	if v.CurrentlyAttached {
		calls = append(calls, "VUSB.get_allowed_operations", "VUSB.unplug")
	}
	if err := permitted(s, calls...); err != nil {
		return err
	}
	var unplugErr error
	if v.CurrentlyAttached {
		if allowed, err := xenapi.VUSB.GetAllowedOperations(s, ref); err != nil {
			unplugErr = err
		} else if slices.Contains(allowed, xenapi.VusbOperationsUnplug) {
			unplugErr = wait(s)(xenapi.VUSB.AsyncUnplug(s, ref))
		}
	}
	if err := wait(s)(xenapi.VUSB.AsyncDestroy(s, ref)); err != nil {
		return xapi.Sanitize(err)
	}
	return xapi.Sanitize(unplugErr)
}
