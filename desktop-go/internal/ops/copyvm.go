package ops

import (
	"errors"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/xapi"
)

// CopyVM is CopyVMDialog's action on a halted VM (VmRules.CanCopy): VMCloneAction (a fast storage-level clone) or, with
// full, VMCopyAction (a full copy to the SR); then the copy's description. It returns the copy's uuid.
func CopyVM(c *xapi.Connection, vmUUID, name, description string, full bool, srUUID string) (string, error) {
	s := c.Session
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("enter a name for the copy")
	}
	vm, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	rec, err := xenapi.VM.GetRecord(s, vm)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	_, defaultTemplate := rec.OtherConfig["default_template"]
	_, internal := rec.OtherConfig["xensource_internal"]
	switch { // VmRules.CanCopy, or CanCopyTemplate for a user template
	case rec.IsASnapshot || rec.IsControlDomain || rec.IsATemplate && (defaultTemplate || internal):
		return "", errors.New("only a VM or a user template can be copied")
	case !rec.IsATemplate && rec.PowerState != xenapi.VMPowerStateHalted:
		return "", errors.New("only a shut down VM can be copied")
	case !contains(rec.AllowedOperations, xenapi.VMOperationsClone) && !contains(rec.AllowedOperations, xenapi.VMOperationsCopy):
		return "", errors.New("the VM cannot be copied now")
	}
	var ref string
	if full {
		if srUUID == "" {
			return "", errors.New("select the storage repository for the copy")
		}
		sr, err := xenapi.SR.GetByUUID(s, srUUID)
		if err != nil {
			return "", xapi.Sanitize(err)
		}
		ref, err = waitRef(s)(xenapi.VM.AsyncCopy(s, vm, name, sr))
		if err != nil {
			return "", xapi.Sanitize(err)
		}
	} else {
		if ref, err = waitRef(s)(xenapi.VM.AsyncClone(s, vm, name)); err != nil {
			return "", xapi.Sanitize(err)
		}
	}
	copied := xenapi.VMRef(ref)
	if err := xenapi.VM.SetNameDescription(s, copied, description); err != nil {
		return "", xapi.Sanitize(err)
	}
	uuid, err := xenapi.VM.GetUUID(s, copied)
	return uuid, xapi.Sanitize(err)
}
