package ops

import (
	"errors"
	"slices"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// vTPMs (VtpmRules): VM > Manage vTPMs, with NewVtpmAction and RemoveVtpmAction. The server's allowed operations
// decide (XAPI 23.11 on, as XCP-ng 8.3).

// Vtpm is a vTPM of the VM.
type Vtpm struct {
	UUID          string `json:"uuid"`
	Protected     bool   `json:"protected"`
	Unique        bool   `json:"unique"`
	RemoveBlocker string `json:"remove_blocker"` // VtpmRules.RemoveBlocker ("" when it can be removed)
}

// Vtpms is the Manage vTPMs dialog of a VM.
type Vtpms struct {
	Name       string `json:"name"`
	Blocker    string `json:"blocker"`     // VtpmRules.CanManage / Blocker: why vTPMs cannot be managed ("" when they can)
	AddBlocker string `json:"add_blocker"` // VM.CanAddVtpm
	Vtpms      []Vtpm `json:"vtpms"`
}

// xapiAtLeast is Helpers.XapiEqualOrGreater_*: the coordinator's xapi_build (else xapi) version (true when unknown).
func xapiAtLeast(s *xenapi.Session, hosts map[xenapi.HostRef]xenapi.HostRecord, v string) bool {
	_, p, _, err := poolOf(s)
	if err != nil {
		return true
	}
	h := hosts[p.Master]
	have := h.SoftwareVersion["xapi_build"]
	if have == "" {
		have = h.SoftwareVersion["xapi"]
	}
	return have == "" || rules.CompareVersions(have, v) >= 0
}

// vtpmVM reads the VM and why its vTPMs cannot be managed: a real HVM VM booting with UEFI, where vTPMs are licensed.
func vtpmVM(s *xenapi.Session, uuid string) (xenapi.VMRef, xenapi.VMRecord, string, error) {
	vm, ref, blocker, _, err := vtpmVMNew(s, uuid)
	return ref, vm, blocker, err
}

// vtpmVMNew is vtpmVM with whether the server decides by its allowed operations (XAPI 23.11 on).
func vtpmVMNew(s *xenapi.Session, uuid string) (xenapi.VMRecord, xenapi.VMRef, string, bool, error) {
	ref, err := xenapi.VM.GetByUUID(s, uuid)
	if err != nil {
		return xenapi.VMRecord{}, "", "", false, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return xenapi.VMRecord{}, "", "", false, xapi.Sanitize(err)
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return xenapi.VMRecord{}, "", "", false, xapi.Sanitize(err)
	}
	modern := xapiAtLeast(s, hosts, "23.11.0")
	switch {
	case vm.IsATemplate || vm.IsASnapshot || vm.IsControlDomain:
		return vm, ref, "Only VMs have vTPMs.", modern, nil
	case featureForbidden(hosts, "restrict_vtpm"): // Host.RestrictVtpm (BoolKeyPreferTrue)
		return vm, ref, "vTPMs are not available due to license restrictions on the server.", modern, nil
	case !xapiAtLeast(s, hosts, "22.26.0"):
		return vm, ref, "vTPMs are not supported by this version of the server.", modern, nil
	case vm.HVMBootPolicy == "" || strings.ToLower(strings.TrimSpace(vm.HVMBootParams["firmware"])) != "uefi": // IsHVM, IsDefaultBootModeUefi
		return vm, ref, messages.Get("COMMAND_VTPM_DISABLED_NON_UEFI"), modern, nil
	}
	return vm, ref, "", modern, nil
}

// ReadVtpms reads the VM's vTPMs.
func ReadVtpms(c *xapi.Connection, vmUUID string) (Vtpms, error) {
	s := c.Session
	vm, _, blocker, modern, err := vtpmVMNew(s, vmUUID)
	if err != nil {
		return Vtpms{}, err
	}
	out := Vtpms{Name: vm.NameLabel, Blocker: blocker, AddBlocker: addBlocker(vm, modern), Vtpms: []Vtpm{}}
	for _, r := range vm.VTPMs {
		t, err := xenapi.VTPM.GetRecord(s, r)
		if err != nil {
			continue
		}
		v := Vtpm{UUID: t.UUID, Protected: t.IsProtected, Unique: t.IsUnique, RemoveBlocker: removeBlocker(vm, t, modern)}
		out.Vtpms = append(out.Vtpms, v)
	}
	return out, nil
}

// AddVtpm is NewVtpmAction: VTPM.create, not unique.
func AddVtpm(c *xapi.Connection, vmUUID string) error {
	s := c.Session
	vm, ref, blocker, modern, err := vtpmVMNew(s, vmUUID)
	switch {
	case err != nil:
		return err
	case blocker != "":
		return errors.New(blocker)
	}
	if why := addBlocker(vm, modern); why != "" {
		return errors.New(why)
	}
	_, err = xenapi.VTPM.Create(s, ref, false)
	return xapi.Sanitize(err)
}

// RemoveVtpm is RemoveVtpmAction (after VTPM_REMOVE_WARNING): VTPM.destroy of a vTPM of this VM.
func RemoveVtpm(c *xapi.Connection, vmUUID, vtpmUUID string) error {
	s := c.Session
	ref, err := xenapi.VTPM.GetByUUID(s, vtpmUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	t, err := xenapi.VTPM.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if owner, err := xenapi.VM.GetUUID(s, t.VM); err != nil || owner != vmUUID {
		return errors.New("the vTPM is not this VM's")
	}
	vm, _, _, modern, err := vtpmVMNew(s, vmUUID)
	if err != nil {
		return err
	}
	if why := removeBlocker(vm, t, modern); why != "" {
		return errors.New(why)
	}
	return xapi.Sanitize(xenapi.VTPM.Destroy(s, ref))
}

// addBlocker is VM.CanAddVtpm: the allowed operations, else (before XAPI 23.11) one vTPM at most, on a halted VM.
func addBlocker(vm xenapi.VMRecord, modern bool) string {
	switch {
	case modern && !slices.Contains(vm.AllowedOperations, xenapi.VMOperationsCreateVtpm):
		return messages.Get("VTPM_OPERATION_DISALLOWED_ADD")
	case modern:
		return ""
	case len(vm.VTPMs) >= 1: // VM.MAX_ALLOWED_VTPMS
		return messages.Format(messages.Get("VTPM_MAX_REACHED"), 1)
	case vm.PowerState != xenapi.VMPowerStateHalted:
		return messages.Get("VTPM_POWER_STATE_WRONG_ATTACH")
	}
	return ""
}

// removeBlocker is VtpmRules.RemoveBlocker: the vTPM's allowed operations, else (before XAPI 23.11) a halted VM.
func removeBlocker(vm xenapi.VMRecord, t xenapi.VTPMRecord, modern bool) string {
	switch {
	case modern && !slices.Contains(t.AllowedOperations, xenapi.VtpmOperationsDestroy):
		return messages.Get("VTPM_OPERATION_DISALLOWED_REMOVE")
	case !modern && vm.PowerState != xenapi.VMPowerStateHalted:
		return messages.Get("VTPM_POWER_STATE_WRONG_REMOVE")
	}
	return ""
}
