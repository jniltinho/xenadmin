package ops

import (
	"errors"
	"fmt"

	xenapi "go/xenapi"

	"xcpng-center/internal/xapi"
)

// DeleteVM is DeleteVMCommand's VMDestroyAction: the chosen snapshots (with their own disks), then the VM, then the
// chosen disks and the VM's suspend image. Everything is checked on the server first: a VM (not a template, snapshot
// or control domain) that allows destroy, disks attached to it alone, and snapshots of it; anything else is refused.
// The disks and snapshots that cannot be deleted do not stop the others; the first error is returned.
func DeleteVM(c *xapi.Connection, vmUUID string, vbdUUIDs, snapshotUUIDs []string) error {
	s := c.Session
	vm, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	rec, err := xenapi.VM.GetRecord(s, vm)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if rec.IsATemplate || rec.IsASnapshot || rec.IsControlDomain || !contains(rec.AllowedOperations, xenapi.VMOperationsDestroy) {
		return errors.New("this VM cannot be deleted now")
	}
	var vdis []xenapi.VDIRef
	for _, u := range vbdUUIDs {
		vdi, err := ownDisk(s, vm, rec, u)
		if err != nil {
			return err
		}
		vdis = append(vdis, vdi)
	}
	var snaps []string
	for _, u := range snapshotUUIDs {
		snap, err := xenapi.VM.GetByUUID(s, u)
		if err != nil {
			return xapi.Sanitize(err)
		}
		if !contains(rec.Snapshots, snap) {
			return fmt.Errorf("%s is not a snapshot of this VM", u)
		}
		snaps = append(snaps, u)
	}
	var caught []error
	for _, u := range snaps {
		if err := DeleteSnapshot(c, u); err != nil {
			caught = append(caught, err)
		}
	}
	if rec.SuspendVDI != "" && rec.SuspendVDI != nullRef { // CA-91072
		vdis = append(vdis, rec.SuspendVDI)
	}
	if err := xenapi.VM.Destroy(s, vm); err != nil {
		return xapi.Sanitize(errors.Join(append([]error{err}, caught...)...))
	}
	for _, v := range vdis {
		err := xenapi.VDI.Destroy(s, v)
		var api *xenapi.APIError
		if err != nil && !(errors.As(err, &api) && api.Message == "HANDLE_INVALID") { // CA-115249: gone already
			caught = append(caught, err)
		}
	}
	if len(caught) > 0 {
		return xapi.Sanitize(caught[0])
	}
	return nil
}

// ownDisk is a disk of the VM (by its VBD's uuid) that no other VM uses: what ConfirmVMDeleteDialog lists.
func ownDisk(s *xenapi.Session, vm xenapi.VMRef, rec xenapi.VMRecord, vbdUUID string) (xenapi.VDIRef, error) {
	ref, err := xenapi.VBD.GetByUUID(s, vbdUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	vbd, err := xenapi.VBD.GetRecord(s, ref)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	if vbd.VM != vm || vbd.Type == xenapi.VbdTypeCD || vbd.VDI == "" || vbd.VDI == nullRef {
		return "", fmt.Errorf("%s is not a virtual disk of this VM", vbdUUID)
	}
	users, err := xenapi.VDI.GetVBDs(s, vbd.VDI)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	for _, b := range users {
		if other, err := xenapi.VBD.GetVM(s, b); err != nil || other != vm {
			return "", fmt.Errorf("the disk of %s is used by another VM too", vbdUUID)
		}
	}
	return vbd.VDI, nil
}
