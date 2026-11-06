package ops

import (
	"errors"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// The GPU and USB tabs' writes (HardwareRules: SetGpuPlacementPolicyAction, VgpuConfigurationAction,
// SetUsbPassthroughAction).

// FeatureForbidden is Helpers.FeatureForbidden: a licence restriction on any server of the connection.
func FeatureForbidden(hosts map[xenapi.HostRef]xenapi.HostRecord, key string) bool {
	return featureForbidden(hosts, key)
}

// SetGPUPlacement is SetGpuPlacementPolicyAction: every GPU group of the pool gets the algorithm.
func SetGPUPlacement(c *xapi.Connection, algorithm string) error {
	s := c.Session
	a := xenapi.AllocationAlgorithm(algorithm)
	if a != xenapi.AllocationAlgorithmBreadthFirst && a != xenapi.AllocationAlgorithmDepthFirst {
		return errors.New("choose a placement policy")
	}
	if err := permitted(s, "GPU_group.set_allocation_algorithm"); err != nil {
		return err
	}
	groups, err := xenapi.GPUGroup.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	for ref := range groups {
		if err := xenapi.GPUGroup.SetAllocationAlgorithm(s, ref, a); err != nil {
			return xapi.Sanitize(err)
		}
	}
	return nil
}

// SetVGPUTypes is VgpuConfigurationAction: the changed types (vGPU type uuid -> now allowed) apply to each GPU's
// current allowed types (HardwareRules.ConfigureAction).
func SetVGPUTypes(c *xapi.Connection, pgpuUUIDs []string, changes map[string]bool) error {
	s := c.Session
	if len(changes) == 0 {
		return nil
	}
	if err := permitted(s, "PGPU.set_enabled_VGPU_types"); err != nil {
		return err
	}
	byRef := map[string]bool{}
	for u, allowed := range changes {
		ref, err := xenapi.VGPUType.GetByUUID(s, u)
		if err != nil {
			return xapi.Sanitize(err)
		}
		byRef[string(ref)] = allowed
	}
	for _, u := range pgpuUUIDs {
		ref, err := xenapi.PGPU.GetByUUID(s, u)
		if err != nil {
			return xapi.Sanitize(err)
		}
		enabled, err := xenapi.PGPU.GetEnabledVGPUTypes(s, ref)
		if err != nil {
			return xapi.Sanitize(err)
		}
		current := make([]string, len(enabled))
		for i, r := range enabled {
			current[i] = string(r)
		}
		var list []xenapi.VGPUTypeRef
		for _, r := range rules.ApplyVGPUChanges(current, byRef) {
			list = append(list, xenapi.VGPUTypeRef(r))
		}
		if err := xenapi.PGPU.SetEnabledVGPUTypes(s, ref, list); err != nil {
			return xapi.Sanitize(err)
		}
	}
	return nil
}

// USBAttached is HardwareRules.IsAttached: the device's USB group has a VUSB.
func USBAttached(s *xenapi.Session, p xenapi.PUSBRecord) bool {
	vusbs, err := xenapi.USBGroup.GetVUSBs(s, p.USBGroup)
	return err == nil && len(vusbs) > 0
}

// USBAttachedVM is HardwareRules.AttachedVm: the name of the VM a device is passed through to ("" if none).
func USBAttachedVM(s *xenapi.Session, p xenapi.PUSBRecord) string {
	vusbs, err := xenapi.USBGroup.GetVUSBs(s, p.USBGroup)
	if err != nil || len(vusbs) == 0 {
		return ""
	}
	vm, err := xenapi.VUSB.GetVM(s, vusbs[0])
	if err != nil {
		return ""
	}
	name, _ := xenapi.VM.GetNameLabel(s, vm)
	return name
}

// SetUSBPassthrough is UsbPage.buttonPassthrough_Click then SetUsbPassthroughAction: nothing when the device changed
// since the user was asked (the user confirmed something else), or when it is now attached to a VM.
func SetUSBPassthrough(c *xapi.Connection, pusbUUID string, wasEnabled bool) error {
	s := c.Session
	if err := permitted(s, "PUSB.set_passthrough_enabled"); err != nil {
		return err
	}
	ref, err := xenapi.PUSB.GetByUUID(s, pusbUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	p, err := xenapi.PUSB.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if p.PassthroughEnabled != wasEnabled || wasEnabled && USBAttached(s, p) {
		return nil
	}
	return xapi.Sanitize(xenapi.PUSB.SetPassthroughEnabled(s, ref, !wasEnabled))
}
