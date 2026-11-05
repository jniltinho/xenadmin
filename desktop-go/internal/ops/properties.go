package ops

import (
	"errors"
	"strconv"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// VMChanges are the Properties dialog's edits; nil fields are unchanged.
type VMChanges struct {
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	VCPUsMax       *int    `json:"vcpus_max,omitempty"`
	VCPUsAtStartup *int    `json:"vcpus_at_startup,omitempty"`
	CoresPerSocket *int64  `json:"cores_per_socket,omitempty"`
	Priority       *int    `json:"priority,omitempty"`
	HomeServer     *string `json:"home_server,omitempty"` // host uuid, "" for none
}

// SaveVMProperties is the pages' SaveSettings: name and description, the vCPU weight and cores per socket (only on a
// halted VM, as the page), ChangeVCPUSettingsAction (a running VM only gains vCPUs live) and VM.set_affinity.
func SaveVMProperties(c *xapi.Connection, vmUUID string, ch VMChanges) error {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if ch.Name != nil {
		if *ch.Name == "" {
			return errors.New("enter a name for the VM")
		}
		if err := xenapi.VM.SetNameLabel(s, ref, *ch.Name); err != nil {
			return xapi.Sanitize(err)
		}
	}
	if ch.Description != nil {
		if err := xenapi.VM.SetNameDescription(s, ref, *ch.Description); err != nil {
			return xapi.Sanitize(err)
		}
	}
	halted := vm.PowerState == xenapi.VMPowerStateHalted
	if (ch.Priority != nil || ch.CoresPerSocket != nil) && !halted {
		return errors.New(messages.Get("VCPU_ONLY_WHEN_HALTED"))
	}
	if ch.Priority != nil {
		_ = xenapi.VM.RemoveFromVCPUsParams(s, ref, "weight")
		if err := xenapi.VM.AddToVCPUsParams(s, ref, "weight", strconv.Itoa(rules.Weight(*ch.Priority))); err != nil {
			return xapi.Sanitize(err)
		}
	}
	if ch.CoresPerSocket != nil {
		vcpus := int64(vm.VCPUsMax)
		if ch.VCPUsMax != nil {
			vcpus = int64(*ch.VCPUsMax)
		}
		if why := rules.ValidVCPUConfiguration(vcpus, *ch.CoresPerSocket); why != "" {
			return errors.New(why)
		}
		_ = xenapi.VM.RemoveFromPlatform(s, ref, "cores-per-socket")
		if err := xenapi.VM.AddToPlatform(s, ref, "cores-per-socket", strconv.FormatInt(*ch.CoresPerSocket, 10)); err != nil {
			return xapi.Sanitize(err)
		}
	}
	if ch.VCPUsMax != nil || ch.VCPUsAtStartup != nil {
		if err := changeVCPUs(s, ref, vm, ch); err != nil {
			return err
		}
	}
	if ch.HomeServer != nil {
		host := xenapi.HostRef(nullRef)
		if *ch.HomeServer != "" {
			if host, err = xenapi.Host.GetByUUID(s, *ch.HomeServer); err != nil {
				return xapi.Sanitize(err)
			}
		}
		if err := xenapi.VM.SetAffinity(s, ref, host); err != nil {
			return xapi.Sanitize(err)
		}
	}
	return nil
}

// VCPUHotplugRestricted is Helpers.FeatureForbidden(Host.RestrictVcpuHotplug): a server's licence forbids changing
// the vCPUs of a running VM (VM.SupportsVcpuHotplug is then false). The key restricts unless it says "false"
// (BoolKeyPreferTrue).
func VCPUHotplugRestricted(hosts map[xenapi.HostRef]xenapi.HostRecord) bool {
	for _, h := range hosts {
		if v, ok := h.LicenseParams["restrict_set_vcpus_number_live"]; !ok || v != "false" {
			return true
		}
	}
	return false
}

// changeVCPUs is ChangeVCPUSettingsAction: a running VM only takes more vCPUs live; a halted one changes at_startup
// and max in the order that keeps at_startup <= max.
func changeVCPUs(s *xenapi.Session, ref xenapi.VMRef, vm xenapi.VMRecord, ch VMChanges) error {
	maxV, start := vm.VCPUsMax, vm.VCPUsAtStartup
	if ch.VCPUsMax != nil {
		maxV = *ch.VCPUsMax
	}
	if ch.VCPUsAtStartup != nil {
		start = *ch.VCPUsAtStartup
	} else if ch.VCPUsMax != nil && start > maxV {
		start = maxV
	}
	if start < 1 || start > maxV {
		return errors.New("the initial number of vCPUs must be between 1 and the maximum")
	}
	if vm.PowerState == xenapi.VMPowerStateRunning {
		if hosts, err := xenapi.Host.GetAllRecords(s); err != nil {
			return xapi.Sanitize(err)
		} else if VCPUHotplugRestricted(hosts) {
			return errors.New(messages.Get("VCPU_ONLY_WHEN_HALTED"))
		}
		if vm.VCPUsAtStartup > start {
			return errors.New(messages.Format(messages.Get("VM_VCPU_CANNOT_UNPLUG_LIVE"), vm.VCPUsAtStartup))
		}
		return xapi.Sanitize(xenapi.VM.SetVCPUsNumberLive(s, ref, start))
	}
	if vm.PowerState != xenapi.VMPowerStateHalted {
		return errors.New(messages.Get("VCPU_ONLY_WHEN_HALTED"))
	}
	if vm.VCPUsAtStartup > start {
		if err := xenapi.VM.SetVCPUsAtStartup(s, ref, start); err != nil {
			return xapi.Sanitize(err)
		}
		return xapi.Sanitize(xenapi.VM.SetVCPUsMax(s, ref, maxV))
	}
	if err := xenapi.VM.SetVCPUsMax(s, ref, maxV); err != nil {
		return xapi.Sanitize(err)
	}
	return xapi.Sanitize(xenapi.VM.SetVCPUsAtStartup(s, ref, start))
}
