package ops

import (
	"errors"
	"slices"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/xapi"
)

// vApps (VappRules): Manage vApps, the New vApp wizard, a vApp's Properties, Start and Shut down, and VM > Assign to
// vApp, with CreateVMApplianceAction, DestroyVMApplianceAction, StartApplianceAction, ShutDownApplianceAction,
// AssignVMsToVMApplianceAction and SetVMStartupOptionsAction.

// VappStart is a VM's start order and delay in a vApp (the VM Startup Sequence page).
type VappStart struct {
	Order int `json:"order"`
	Delay int `json:"delay"`
}

// VappSpec is the New vApp wizard's or the Properties dialog's values.
type VappSpec struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	VMs         []string             `json:"vms"`     // the members' uuids
	Startup     map[string]VappStart `json:"startup"` // by VM uuid, only the changed ones
}

func vmRefs(s *xenapi.Session, uuids []string) ([]xenapi.VMRef, error) {
	var out []xenapi.VMRef
	for _, u := range uuids {
		ref, err := xenapi.VM.GetByUUID(s, u)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		rec, err := xenapi.VM.GetRecord(s, ref)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		if rec.IsATemplate || rec.IsASnapshot || rec.IsControlDomain { // VappRules.Candidates: real VMs only
			return nil, errors.New("only VMs can be in a vApp")
		}
		out = append(out, ref)
	}
	return out, nil
}

// setVappStartup is SetVMStartupOptionsAction for the changed VMs, then pool.sync_database.
func setVappStartup(s *xenapi.Session, startup map[string]VappStart) error {
	if len(startup) == 0 {
		return nil
	}
	for _, st := range startup {
		if st.Order < 0 || st.Delay < 0 {
			return errors.New("the start order and the delay must be 0 or more")
		}
	}
	for u, st := range startup {
		ref, err := xenapi.VM.GetByUUID(s, u)
		if err != nil {
			return xapi.Sanitize(err)
		}
		if err := xenapi.VM.SetOrder(s, ref, st.Order); err != nil {
			return xapi.Sanitize(err)
		}
		if err := xenapi.VM.SetStartDelay(s, ref, st.Delay); err != nil {
			return xapi.Sanitize(err)
		}
	}
	return xapi.Sanitize(wait(s)(xenapi.Pool.AsyncSyncDatabase(s)))
}

// CreateVapp is the New vApp wizard's finish: CreateVMApplianceAction (the vApp, then each VM's appliance) and the
// changed startup options; the new vApp's uuid.
func CreateVapp(c *xapi.Connection, spec VappSpec) (string, error) {
	s := c.Session
	if strings.TrimSpace(spec.Name) == "" { // NewVMApplianceNamePage.EnableNext
		return "", errors.New("enter a name for the vApp")
	}
	vms, err := vmRefs(s, spec.VMs)
	if err != nil {
		return "", err
	}
	r, err := waitRef(s)(xenapi.VMAppliance.AsyncCreate(s, xenapi.VMApplianceRecord{NameLabel: strings.TrimSpace(spec.Name), NameDescription: spec.Description}))
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	app := xenapi.VMApplianceRef(r)
	for _, vm := range vms {
		if err := xenapi.VM.SetAppliance(s, vm, app); err != nil {
			return "", xapi.Sanitize(err)
		}
	}
	uuid, _ := xenapi.VMAppliance.GetUUID(s, app)
	return uuid, setVappStartup(s, spec.Startup)
}

// SaveVapp is a vApp's Properties: name and description, the members (VappRules.SetMembersAction: the others leave)
// and the changed startup options.
func SaveVapp(c *xapi.Connection, uuid string, spec VappSpec) error {
	s := c.Session
	if strings.TrimSpace(spec.Name) == "" {
		return errors.New("enter a name for the vApp")
	}
	ref, err := xenapi.VMAppliance.GetByUUID(s, uuid)
	if err != nil {
		return xapi.Sanitize(err)
	}
	app, err := xenapi.VMAppliance.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vms, err := vmRefs(s, spec.VMs)
	if err != nil {
		return err
	}
	if name := strings.TrimSpace(spec.Name); name != app.NameLabel {
		if err := xenapi.VMAppliance.SetNameLabel(s, ref, name); err != nil {
			return xapi.Sanitize(err)
		}
	}
	if spec.Description != app.NameDescription {
		if err := xenapi.VMAppliance.SetNameDescription(s, ref, spec.Description); err != nil {
			return xapi.Sanitize(err)
		}
	}
	if err := setMembers(s, ref, app.VMs, vms); err != nil {
		return err
	}
	return setVappStartup(s, spec.Startup)
}

// setMembers is AssignVMsToVMApplianceAction: the members not kept leave, the new ones join.
func setMembers(s *xenapi.Session, app xenapi.VMApplianceRef, current, want []xenapi.VMRef) error {
	for _, vm := range current {
		if !slices.Contains(want, vm) {
			if err := xenapi.VM.SetAppliance(s, vm, nullRef); err != nil {
				return xapi.Sanitize(err)
			}
		}
	}
	for _, vm := range want {
		if a, err := xenapi.VM.GetAppliance(s, vm); err != nil || a != app {
			if err := xenapi.VM.SetAppliance(s, vm, app); err != nil {
				return xapi.Sanitize(err)
			}
		}
	}
	return nil
}

// DeleteVapp is DestroyVMApplianceAction: the VMs leave the vApp (they stay), then it is destroyed.
func DeleteVapp(c *xapi.Connection, uuid string) error {
	s := c.Session
	ref, err := xenapi.VMAppliance.GetByUUID(s, uuid)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vms, err := xenapi.VMAppliance.GetVMs(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	for _, vm := range vms {
		if err := xenapi.VM.SetAppliance(s, vm, nullRef); err != nil {
			return xapi.Sanitize(err)
		}
	}
	return xapi.Sanitize(xenapi.VMAppliance.Destroy(s, ref))
}

// VappPower is StartApplianceAction ("start": its VMs in the startup order) or ShutDownApplianceAction ("shutdown"),
// when the vApp allows it now.
func VappPower(c *xapi.Connection, uuid, op string) error {
	s := c.Session
	ref, err := xenapi.VMAppliance.GetByUUID(s, uuid)
	if err != nil {
		return xapi.Sanitize(err)
	}
	app, err := xenapi.VMAppliance.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	switch op {
	case "start":
		if len(app.VMs) == 0 || !slices.Contains(app.AllowedOperations, xenapi.VMApplianceOperationStart) {
			return errors.New("the vApp cannot be started now")
		}
		return xapi.Sanitize(wait(s)(xenapi.VMAppliance.AsyncStart(s, ref, false)))
	case "shutdown":
		if !slices.Contains(app.AllowedOperations, xenapi.VMApplianceOperationShutdown) {
			return errors.New("the vApp cannot be shut down now")
		}
		return xapi.Sanitize(wait(s)(xenapi.VMAppliance.AsyncShutdown(s, ref)))
	}
	return errors.New("unknown vApp operation")
}

// AssignVapp is AssignGroupToVMCommand for one VM: a VM already in that vApp leaves it (the checked menu item),
// otherwise it joins it.
func AssignVapp(c *xapi.Connection, vappUUID, vmUUID string) error {
	s := c.Session
	app, err := xenapi.VMAppliance.GetByUUID(s, vappUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vms, err := vmRefs(s, []string{vmUUID})
	if err != nil {
		return err
	}
	current, err := xenapi.VM.GetAppliance(s, vms[0])
	if err != nil {
		return xapi.Sanitize(err)
	}
	if current == app {
		return xapi.Sanitize(xenapi.VM.SetAppliance(s, vms[0], nullRef))
	}
	return xapi.Sanitize(xenapi.VM.SetAppliance(s, vms[0], app))
}
