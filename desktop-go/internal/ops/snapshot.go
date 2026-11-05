package ops

import (
	"errors"

	xenapi "go/xenapi"

	"xcpng-center/internal/xapi"
)

// TakeSnapshot is VMSnapshotCreateAction: a disk snapshot, or a disk and memory one (checkpoint), then its
// description. shortcut: no console screenshot saved with a memory snapshot (SaveImageInBlob).
func TakeSnapshot(c *xapi.Connection, vmUUID, name, description string, memory bool) error {
	if name == "" {
		return errors.New("enter a name for the snapshot")
	}
	s := c.Session
	vm, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	var result string
	if memory {
		result, err = waitResult(s)(xenapi.VM.AsyncCheckpoint(s, vm, name))
	} else {
		result, err = waitResult(s)(xenapi.VM.AsyncSnapshot(s, vm, name, nil))
	}
	if err != nil {
		return xapi.Sanitize(err)
	}
	// always set, also when empty: a snapshot would otherwise keep the VM's description
	if ref := resultRef(result); ref != "" {
		return xapi.Sanitize(xenapi.VM.SetNameDescription(s, xenapi.VMRef(ref), description))
	}
	return errors.New("the snapshot was taken but XAPI did not say which it is; its description was not set")
}

// snapshotOf finds a snapshot by uuid and refuses anything else (deleting a VM through here must be impossible).
func snapshotOf(s *xenapi.Session, uuid string) (xenapi.VMRef, xenapi.VMRecord, error) {
	ref, err := xenapi.VM.GetByUUID(s, uuid)
	if err != nil {
		return "", xenapi.VMRecord{}, err
	}
	rec, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return "", xenapi.VMRecord{}, err
	}
	if !rec.IsASnapshot {
		return "", xenapi.VMRecord{}, errors.New("this uuid is not a snapshot: refusing")
	}
	return ref, rec, nil
}

// RevertSnapshot is VMSnapshotRevertAction: VM.revert, then, when the VM was running at the snapshot, start or resume
// it, on its previous server when it can boot there.
func RevertSnapshot(c *xapi.Connection, snapshotUUID string) error {
	s := c.Session
	snap, rec, err := snapshotOf(s, snapshotUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	previous, _ := xenapi.VM.GetResidentOn(s, rec.SnapshotOf)
	if err := wait(s)(xenapi.VM.AsyncRevert(s, snap)); err != nil {
		return xapi.Sanitize(err)
	}
	if rec.SnapshotInfo["power-state-at-snapshot"] != "Running" {
		return nil
	}
	vm := rec.SnapshotOf
	state, err := xenapi.VM.GetPowerState(s, vm)
	if err != nil {
		return xapi.Sanitize(err)
	}
	onPrevious := previous != "" && previous != nullRef && xenapi.VM.AssertCanBootHere(s, vm, previous) == nil
	switch state {
	case xenapi.VMPowerStateHalted:
		if onPrevious {
			err = wait(s)(xenapi.VM.AsyncStartOn(s, vm, previous, false, false))
		} else {
			err = wait(s)(xenapi.VM.AsyncStart(s, vm, false, false))
		}
	case xenapi.VMPowerStateSuspended:
		if onPrevious {
			err = wait(s)(xenapi.VM.AsyncResumeOn(s, vm, previous, false, false))
		} else {
			err = wait(s)(xenapi.VM.AsyncResume(s, vm, false, false))
		}
	}
	return xapi.Sanitize(err)
}

// DeleteSnapshot is VMSnapshotDeleteAction: a suspended snapshot is shut down first; the snapshot goes with the disks
// it owns (VBD other_config "owner") and its suspend image; a disk XAPI already removed is not an error.
func DeleteSnapshot(c *xapi.Connection, snapshotUUID string) error {
	s := c.Session
	snap, rec, err := snapshotOf(s, snapshotUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if rec.PowerState == xenapi.VMPowerStateSuspended {
		if err := xenapi.VM.HardShutdown(s, snap); err != nil {
			return xapi.Sanitize(err)
		}
	}
	var vdis []xenapi.VDIRef
	for _, b := range rec.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil {
			return xapi.Sanitize(err)
		}
		if _, owner := vbd.OtherConfig["owner"]; owner && vbd.VDI != nullRef {
			vdis = append(vdis, vbd.VDI)
		}
	}
	if rec.SuspendVDI != "" && rec.SuspendVDI != nullRef {
		vdis = append(vdis, rec.SuspendVDI)
	}
	if err := xenapi.VM.Destroy(s, snap); err != nil {
		return xapi.Sanitize(err)
	}
	var errs []error
	for _, v := range vdis {
		err := xenapi.VDI.Destroy(s, v)
		var api *xenapi.APIError
		if errors.As(err, &api) && api.Message == "HANDLE_INVALID" {
			continue // CA-115249: XAPI deleted it already
		}
		errs = append(errs, err)
	}
	return xapi.Sanitize(errors.Join(errs...))
}
