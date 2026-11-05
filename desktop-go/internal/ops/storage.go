package ops

import (
	"errors"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
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

// AddDisk is NewDiskDialog + CreateDiskAction + VbdCreateAndPlugAction: a new disk on the SR (with the VM as its
// placement hint), attached on the first free position, bootable only as position 0 of a VM with no bootable disk,
// owned by the VM, plugged when the VM allows it; the notice says when the VM needs a restart first.
func AddDisk(c *xapi.Connection, vmUUID, name, description, srUUID string, size int) (notice string, err error) {
	if name == "" || size <= 0 {
		return "", errors.New("enter a name and a size for the disk")
	}
	s := c.Session
	vm, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	sr, err := xenapi.SR.GetByUUID(s, srUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	devices, err := xenapi.VM.GetAllowedVBDDevices(s, vm)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	if len(devices) == 0 {
		return "", errors.New("the VM has the maximum number of virtual disks and drives")
	}
	bootable := false
	vbds, _ := xenapi.VM.GetVBDs(s, vm)
	for _, b := range vbds { // CA-44959: bootable only when no other disk is
		if boot, err := xenapi.VBD.GetBootable(s, b); err == nil && boot {
			bootable = true
		}
	}
	vdi, err := xenapi.VDI.Create(s, xenapi.VDIRecord{NameLabel: name, NameDescription: description, SR: sr, VirtualSize: size,
		Type: xenapi.VdiTypeUser, OtherConfig: map[string]string{}, XenstoreData: map[string]string{}, SmConfig: map[string]string{"vmhint": vmUUID}})
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	ref, err := xenapi.VBD.Create(s, xenapi.VBDRecord{VM: vm, VDI: vdi, Userdevice: devices[0], Bootable: devices[0] == "0" && !bootable,
		Mode: xenapi.VbdModeRW, Type: xenapi.VbdTypeDisk, Unpluggable: true, OtherConfig: map[string]string{"owner": "true"},
		QosAlgorithmParams: map[string]string{}})
	if err != nil {
		_ = xenapi.VDI.Destroy(s, vdi) // the disk would be left alone on the SR
		return "", xapi.Sanitize(err)
	}
	if allowed, _ := xenapi.VBD.GetAllowedOperations(s, ref); contains(allowed, xenapi.VbdOperationsPlug) {
		return "", xapi.Sanitize(wait(s)(xenapi.VBD.AsyncPlug(s, ref)))
	}
	if state, _ := xenapi.VM.GetPowerState(s, vm); state != xenapi.VMPowerStateHalted {
		return messages.Get("NEWDISKWIZARD_MESSAGE"), nil
	}
	return "", nil
}

// DetachDisk is DetachVirtualDiskAction: unplug the VBD when attached and allowed, then destroy it; the disk stays.
func DetachDisk(c *xapi.Connection, vbdUUID string) error {
	s := c.Session
	ref, err := xenapi.VBD.GetByUUID(s, vbdUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vbd, err := xenapi.VBD.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if vbd.Type != xenapi.VbdTypeDisk {
		return errors.New("this is not a virtual disk")
	}
	var unplugErr error
	if vbd.CurrentlyAttached && contains(vbd.AllowedOperations, xenapi.VbdOperationsUnplug) {
		unplugErr = wait(s)(xenapi.VBD.AsyncUnplug(s, ref))
	}
	if unplugErr != nil {
		return xapi.Sanitize(unplugErr)
	}
	return xapi.Sanitize(wait(s)(xenapi.VBD.AsyncDestroy(s, ref)))
}

// DiskInfo gathers what the disk rules read: the SR's kind, the HA use and every VBD of the disk. Any failed read is
// an error: the rules must never decide on a partial picture.
func DiskInfo(s *xenapi.Session, vm xenapi.VMRef, vdi xenapi.VDIRecord, sr xenapi.SRRecord) (rules.DiskInfo, error) {
	ha := false
	if vdi.Type == xenapi.VdiTypeHaStatefile || vdi.Type == xenapi.VdiTypeRedoLog { // VDI.IsUsedByHA
		pools, err := xenapi.Pool.GetAllRecords(s)
		if err != nil {
			return rules.DiskInfo{}, err
		}
		for _, p := range pools {
			ha = p.HaEnabled
		}
	}
	d := rules.DiskInfo{SRPhysical: sr.Type == "local" || (sr.Type == "udev" && sr.SmConfig["type"] == "cd"), SRTools: sr.IsToolsSr,
		UsedByHA: ha, System: vdi.Type == xenapi.VdiTypeSystem}
	d.DestroyAllowed = sr.Type == "rawhba" // SR.HBALunPerVDI: one LUN per disk, deleted without the destroy operation
	for _, op := range vdi.AllowedOperations {
		d.DestroyAllowed = d.DestroyAllowed || op == xenapi.VdiOperationsDestroy
	}
	for _, b := range vdi.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil {
			return rules.DiskInfo{}, err
		}
		state, err := xenapi.VM.GetPowerState(s, vbd.VM)
		if err != nil {
			return rules.DiskInfo{}, err
		}
		unplug := false
		for _, op := range vbd.AllowedOperations {
			unplug = unplug || op == xenapi.VbdOperationsUnplug
		}
		d.VBDs = append(d.VBDs, rules.DiskVBD{VMRunning: state == xenapi.VMPowerStateRunning, Attached: vbd.CurrentlyAttached,
			CanUnplug: unplug, ThisVM: vbd.VM == vm})
	}
	return d, nil
}

// DeleteDisk is DestroyDiskAction for the disk of this VM's VBD: the VBD must still be the VM's and the rules
// (StorageRules.CanDelete) are checked again on the server's current state, so a disk shared or reassigned since
// the tab was drawn is never deleted; a disk active on a VM is refused (CANNOT_DELETE_VDI_ACTIVE_ON).
func DeleteDisk(c *xapi.Connection, vmUUID, vbdUUID string) error {
	s := c.Session
	vm, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vbdRef, err := xenapi.VBD.GetByUUID(s, vbdUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vbd, err := xenapi.VBD.GetRecord(s, vbdRef)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if vbd.VM != vm || vbd.Type != xenapi.VbdTypeDisk || vbd.VDI == nullRef {
		return errors.New("this disk is no longer attached to this VM: refresh and try again")
	}
	ref := vbd.VDI
	vdi, err := xenapi.VDI.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	sr, err := xenapi.SR.GetRecord(s, vdi.SR)
	if err != nil {
		return xapi.Sanitize(err)
	}
	for _, b := range vdi.VBDs {
		if attached, err := xenapi.VBD.GetCurrentlyAttached(s, b); err == nil && attached {
			vmRef, _ := xenapi.VBD.GetVM(s, b)
			name, _ := xenapi.VM.GetNameLabel(s, vmRef)
			return errors.New(messages.Format(messages.Get("CANNOT_DELETE_VDI_ACTIVE_ON"), name))
		}
	}
	if len(vdi.VBDs) != 1 || vdi.VBDs[0] != vbdRef {
		return errors.New("this disk is also attached to another VM: it is not deleted from here")
	}
	info, err := DiskInfo(s, vm, vdi, sr)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if !rules.CanDeleteDisk(info) {
		return errors.New("this disk cannot be deleted now (it is used by another VM, is a system disk of a running VM, or its storage does not allow it)")
	}
	for _, b := range vdi.VBDs {
		if err := wait(s)(xenapi.VBD.AsyncDestroy(s, b)); err != nil {
			return xapi.Sanitize(err)
		}
	}
	return xapi.Sanitize(wait(s)(xenapi.VDI.AsyncDestroy(s, ref)))
}

// AttachDisk is AttachDiskDialog + VbdCreateAndPlugAction: an existing disk on the first free position (bootable only
// at position 0), read-only when asked, owned by the VM when no other VM has it; refused when the disk is now attached
// to this VM, or to another one and not sharable.
func AttachDisk(c *xapi.Connection, vmUUID, vdiUUID string, readOnly bool) (notice string, err error) {
	s := c.Session
	vm, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	vdi, err := xenapi.VDI.GetByUUID(s, vdiUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	rec, err := xenapi.VDI.GetRecord(s, vdi)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	for _, b := range rec.VBDs {
		owner, err := xenapi.VBD.GetVM(s, b)
		if err != nil {
			return "", xapi.Sanitize(err)
		}
		if owner == vm || !rec.Sharable {
			return "", errors.New("this disk is already attached: refresh and choose another")
		}
	}
	devices, err := xenapi.VM.GetAllowedVBDDevices(s, vm)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	if len(devices) == 0 {
		return "", errors.New("this VM already has the maximum number of virtual disks")
	}
	mode := xenapi.VbdModeRW
	if readOnly {
		mode = xenapi.VbdModeRO
	}
	other := map[string]string{}
	if len(rec.VBDs) == 0 {
		other["owner"] = "true"
	}
	ref, err := xenapi.VBD.Create(s, xenapi.VBDRecord{VM: vm, VDI: vdi, Userdevice: devices[0], Bootable: devices[0] == "0", Mode: mode,
		Type: xenapi.VbdTypeDisk, Unpluggable: true, OtherConfig: other, QosAlgorithmParams: map[string]string{}})
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	if allowed, _ := xenapi.VBD.GetAllowedOperations(s, ref); contains(allowed, xenapi.VbdOperationsPlug) {
		return "", xapi.Sanitize(wait(s)(xenapi.VBD.AsyncPlug(s, ref)))
	}
	if state, _ := xenapi.VM.GetPowerState(s, vm); state != xenapi.VMPowerStateHalted {
		return messages.Get("NEWDISKWIZARD_MESSAGE"), nil
	}
	return "", nil
}
