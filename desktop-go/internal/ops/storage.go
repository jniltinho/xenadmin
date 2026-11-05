package ops

import (
	"errors"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/xapi"
)

// SetDiskActive is VmStorageView.ToggleActive: VBD.plug or VBD.unplug of a VM's disk, if the VBD allows it now.
func SetDiskActive(c *xapi.Connection, vbdUUID string, active bool) error {
	s := c.Session
	ref, err := xenapi.VBD.GetByUUID(s, vbdUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if active {
		return xapi.Sanitize(wait(s)(xenapi.VBD.AsyncPlug(s, ref)))
	}
	return xapi.Sanitize(wait(s)(xenapi.VBD.AsyncUnplug(s, ref)))
}

// ChangeISO is ChangeVMISOAction: eject what the DVD drive holds, then insert the ISO ("" leaves it empty).
func ChangeISO(c *xapi.Connection, dvdUUID, isoUUID string) error {
	s := c.Session
	ref, err := xenapi.VBD.GetByUUID(s, dvdUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vbd, err := xenapi.VBD.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if vbd.Type != xenapi.VbdTypeCD {
		return errors.New("this is not a DVD drive")
	}
	if !vbd.Empty {
		if err := wait(s)(xenapi.VBD.AsyncEject(s, ref)); err != nil {
			return xapi.Sanitize(err)
		}
	}
	if isoUUID == "" {
		return nil
	}
	vdi, err := xenapi.VDI.GetByUUID(s, isoUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	return xapi.Sanitize(wait(s)(xenapi.VBD.AsyncInsert(s, ref, vdi)))
}

// CreateDVDDrive is CreateCdDriveAction + VbdCreateAndPlugAction: an empty, read-only DVD drive on position 3 when
// free (else the first allowed one), plugged at once when the VM allows it; otherwise, on a VM that is not halted, the
// notice says the drive comes with the next reboot (ShowUserInstruction).
func CreateDVDDrive(c *xapi.Connection, vmUUID string) (notice string, err error) {
	s := c.Session
	vm, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	vbds, err := xenapi.VM.GetVBDs(s, vm)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	for _, b := range vbds { // CreateCdDriveAction: nothing to do when the VM has one already (a stale link)
		if t, err := xenapi.VBD.GetType(s, b); err == nil && t == xenapi.VbdTypeCD {
			return "", nil
		}
	}
	devices, err := xenapi.VM.GetAllowedVBDDevices(s, vm)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	if len(devices) == 0 {
		return "", errors.New("the VM has the maximum number of virtual disks and drives")
	}
	position := devices[0]
	for _, d := range devices {
		if d == "3" {
			position = d
		}
	}
	ref, err := xenapi.VBD.Create(s, xenapi.VBDRecord{VM: vm, VDI: nullRef, Userdevice: position, Bootable: false, Mode: xenapi.VbdModeRO,
		Type: xenapi.VbdTypeCD, Empty: true, Unpluggable: true, OtherConfig: map[string]string{}, QosAlgorithmParams: map[string]string{}})
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	if allowed, _ := xenapi.VBD.GetAllowedOperations(s, ref); contains(allowed, xenapi.VbdOperationsPlug) {
		return "", xapi.Sanitize(wait(s)(xenapi.VBD.AsyncPlug(s, ref)))
	}
	if state, _ := xenapi.VM.GetPowerState(s, vm); state != xenapi.VMPowerStateHalted {
		return messages.Get("NEW_DVD_DRIVE_REBOOT"), nil
	}
	return "", nil
}

func contains[T comparable](list []T, v T) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
