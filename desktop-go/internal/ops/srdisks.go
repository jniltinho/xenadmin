package ops

import (
	"encoding/xml"
	"errors"
	"slices"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// The buttons of an SR's Storage tab (XenAdmin.Desktop/SrStorageView.cs: Add, Rescan, Move, Delete) and of a server's
// or pool's (PhysicalStorageView.cs: Reclaim freed space), with the XenModel actions they run: CreateDiskAction,
// SrRefreshAction, MigrateVirtualDiskAction, MoveVirtualDiskAction, DestroyDiskAction, VMSnapshotDeleteAction,
// SrTrimAction. Every change checks the role first (ApiMethodsToRoleCheck) and the rules again on the server's state.

// DiskEnv is what the SR disk rules read once per connection: the SM drivers' capabilities and features, and whether
// the license forbids storage live migration (Host.RestrictCrossPoolMigrate on any server).
type DiskEnv struct {
	caps             map[string][]string
	features         map[string]map[string]int
	migrateForbidden bool
}

// ReadDiskEnv reads the SM drivers and the servers' licenses.
func ReadDiskEnv(s *xenapi.Session) (DiskEnv, error) {
	sms, err := xenapi.SM.GetAllRecords(s)
	if err != nil {
		return DiskEnv{}, xapi.Sanitize(err)
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return DiskEnv{}, xapi.Sanitize(err)
	}
	e := DiskEnv{caps: map[string][]string{}, features: map[string]map[string]int{}, migrateForbidden: featureForbidden(hosts, "restrict_storage_xen_motion")}
	for _, sm := range sms {
		e.caps[sm.Type], e.features[sm.Type] = sm.Capabilities, sm.Features
	}
	return e, nil
}

// Caps are the capabilities of an SR type's driver (VDI_CREATE, THIN_PROVISIONING...).
func (e DiskEnv) Caps(srType string) []string { return e.caps[srType] }

// Migrates is SR.SupportsStorageMigration: not ISO nor memory, the driver snapshots and mirrors disks.
func (e DiskEnv) Migrates(sr xenapi.SRRecord) bool {
	c := e.caps[sr.Type]
	return sr.ContentType != "iso" && sr.Type != "tmpfs" && slices.Contains(c, "VDI_SNAPSHOT") && slices.Contains(c, "VDI_MIRROR")
}

// Trims is SR.SupportsTrim: the driver has the SR_TRIM feature.
func (e DiskEnv) Trims(sr xenapi.SRRecord) bool {
	_, ok := e.features[sr.Type]["SR_TRIM"]
	return ok
}

// Move reads what UseMigrate needs about a disk on its SR.
func (e DiskEnv) Move(s *xenapi.Session, ref xenapi.VDIRef, vdi xenapi.VDIRecord, sr xenapi.SRRecord) (rules.DiskMove, error) {
	d := rules.DiskMove{Snapshot: vdi.IsASnapshot, Locked: len(vdi.CurrentOperations) > 0, CBT: vdi.CbtEnabled,
		HAType: vdi.Type == xenapi.VdiTypeHaStatefile || vdi.Type == xenapi.VdiTypeRedoLog, LunPerVDI: sr.Type == "rawhba",
		VBDs: len(vdi.VBDs), SRMigrates: e.Migrates(sr), MigrateForbidden: e.migrateForbidden, VMsHalted: true}
	vms, err := DiskVMs(s, ref, vdi)
	if err != nil {
		return rules.DiskMove{}, err
	}
	for _, vm := range vms {
		state, err := xenapi.VM.GetPowerState(s, vm)
		if err != nil {
			return rules.DiskMove{}, xapi.Sanitize(err)
		}
		d.VMsHalted = d.VMsHalted && state == xenapi.VMPowerStateHalted
	}
	return d, nil
}

// DiskVMs is VDI.GetVMs: a crash dump's VM, the VMs suspended to the disk, else the VMs of its VBDs.
func DiskVMs(s *xenapi.Session, ref xenapi.VDIRef, vdi xenapi.VDIRecord) ([]xenapi.VMRef, error) {
	var out []xenapi.VMRef
	switch vdi.Type {
	case xenapi.VdiTypeCrashdump:
		for _, cd := range vdi.CrashDumps {
			vm, err := xenapi.Crashdump.GetVM(s, cd)
			if err != nil {
				return nil, xapi.Sanitize(err)
			}
			out = append(out, vm)
		}
	case xenapi.VdiTypeSuspend:
		vms, err := xenapi.VM.GetAllRecords(s)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		for r, vm := range vms {
			if vm.SuspendVDI == ref {
				out = append(out, r)
			}
		}
	default:
		for _, b := range vdi.VBDs {
			vm, err := xenapi.VBD.GetVM(s, b)
			if err != nil {
				return nil, xapi.Sanitize(err)
			}
			out = append(out, vm)
		}
	}
	return out, nil
}

// SRAttachedHost is SR.GetFirstAttachedStorageHost: the server of its first plugged PBD ("" when detached).
func SRAttachedHost(s *xenapi.Session, sr xenapi.SRRecord) xenapi.HostRef {
	for _, p := range sr.PBDs {
		if pbd, err := xenapi.PBD.GetRecord(s, p); err == nil && pbd.CurrentlyAttached {
			return pbd.Host
		}
	}
	return ""
}

// CreateSRDisk is NewDiskDialog(connection, sr) + CreateDiskAction from the SR Storage tab: a user disk on the SR,
// attached to no VM; it returns the disk's uuid.
func CreateSRDisk(c *xapi.Connection, srUUID, name, description string, size int) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || size <= 0 {
		return "", errors.New("enter a name and a size for the disk")
	}
	s := c.Session
	if err := permitted(s, "vdi.create"); err != nil {
		return "", err
	}
	sr, err := xenapi.SR.GetByUUID(s, srUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	ref, err := xenapi.VDI.Create(s, xenapi.VDIRecord{NameLabel: name, NameDescription: description, SR: sr, VirtualSize: size,
		Type: xenapi.VdiTypeUser, OtherConfig: map[string]string{}, XenstoreData: map[string]string{}, SmConfig: map[string]string{}})
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	uuid, err := xenapi.VDI.GetUUID(s, ref)
	return uuid, xapi.Sanitize(err)
}

// RescanSRs is SrRefreshAction (SR.scan) for each SR, as Rescan of the SR Storage tab and of the Move dialog
// (SrPicker.ScanSRs: the detached SRs are skipped).
func RescanSRs(c *xapi.Connection, srUUIDs ...string) error {
	s := c.Session
	if err := permitted(s, "sr.scan"); err != nil {
		return err
	}
	var errs []error
	for _, u := range srUUIDs {
		ref, err := xenapi.SR.GetByUUID(s, u)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		rec, err := xenapi.SR.GetRecord(s, ref)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if SRAttachedHost(s, rec) == "" {
			continue
		}
		if err := wait(s)(xenapi.SR.AsyncScan(s, ref)); err != nil {
			errs = append(errs, err)
		}
	}
	return xapi.Sanitize(errors.Join(errs...))
}

// diskOf reads a disk and its SR by the disk's uuid.
func diskOf(s *xenapi.Session, uuid string) (xenapi.VDIRef, xenapi.VDIRecord, xenapi.SRRecord, error) {
	ref, err := xenapi.VDI.GetByUUID(s, uuid)
	if err != nil {
		return "", xenapi.VDIRecord{}, xenapi.SRRecord{}, xapi.Sanitize(err)
	}
	vdi, err := xenapi.VDI.GetRecord(s, ref)
	if err != nil {
		return "", xenapi.VDIRecord{}, xenapi.SRRecord{}, xapi.Sanitize(err)
	}
	sr, err := xenapi.SR.GetRecord(s, vdi.SR)
	return ref, vdi, sr, xapi.Sanitize(err)
}

// MoveSRDisk is MoveVirtualDiskDialog's Move: a live storage migration (MigrateVirtualDiskAction: VDI.pool_migrate)
// when StorageRules.UseMigrate says so, else an offline move (MoveVirtualDiskAction: copy to the SR, the VBDs moved
// to the copy, the disk destroyed). It returns the uuid of the disk on the new SR.
// shortcut: the target's free space and its visibility from the VMs' servers are checked by the dialog
// (MoveDiskTargets); the server refuses what still does not fit.
func MoveSRDisk(c *xapi.Connection, vdiUUID, srUUID string) (string, error) {
	s := c.Session
	ref, vdi, sr, err := diskOf(s, vdiUUID)
	if err != nil {
		return "", err
	}
	env, err := ReadDiskEnv(s)
	if err != nil {
		return "", err
	}
	d, err := env.Move(s, ref, vdi, sr)
	if err != nil {
		return "", err
	}
	migrate, ok := rules.UseMigrate(d)
	if !ok {
		return "", errors.New("this disk cannot be moved now: refresh and try again")
	}
	target, err := xenapi.SR.GetByUUID(s, srUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	trec, err := xenapi.SR.GetRecord(s, target)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	switch {
	case target == vdi.SR:
		return "", errors.New(messages.Get("CURRENT_LOCATION"))
	case SRAttachedHost(s, trec) == "":
		return "", errors.New(messages.Get("SR_DETACHED"))
	case migrate && !env.Migrates(trec):
		return "", errors.New(messages.Get("UNSUPPORTED_SR_TYPE"))
	}
	var moved string
	if migrate {
		if err := permitted(s, "vdi.pool_migrate"); err != nil {
			return "", err
		}
		if moved, err = waitRef(s)(xenapi.VDI.AsyncPoolMigrate(s, ref, target, map[string]string{})); err != nil {
			return "", xapi.Sanitize(err)
		}
	} else if moved, err = moveDisk(s, ref, vdi, target); err != nil {
		return "", err
	}
	uuid, err := xenapi.VDI.GetUUID(s, xenapi.VDIRef(moved))
	return uuid, xapi.Sanitize(err)
}

// moveDisk is MoveVirtualDiskAction.Run for a disk of halted VMs: VDI.copy, the old VBDs destroyed, the disk
// destroyed, the VBDs made again on the copy (same position, mode, boot flag and owner).
func moveDisk(s *xenapi.Session, ref xenapi.VDIRef, vdi xenapi.VDIRecord, target xenapi.SRRef) (string, error) {
	if err := permitted(s, "vdi.copy", "vdi.destroy", "vbd.destroy", "vbd.create"); err != nil {
		return "", err
	}
	if _, err := unpluggedVBDs(s, ref, vdi.VBDs); err != nil { // CanMove: only the disks of halted VMs move offline
		return "", err
	}
	moved, err := waitRef(s)(xenapi.VDI.AsyncCopy(s, ref, target, nullRef, nullRef))
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	// read again after the copy (MoveVirtualDiskAction reads the VBDs then): a disk attached or plugged meanwhile is
	// not moved, the copy goes; the VBDs are made again with their current settings
	vbds, err := unpluggedVBDs(s, ref, vdi.VBDs)
	if err != nil {
		_ = xenapi.VDI.Destroy(s, xenapi.VDIRef(moved))
		return "", err
	}
	for _, b := range vdi.VBDs {
		if err := xenapi.VBD.Destroy(s, b); err != nil {
			return "", xapi.Sanitize(err)
		}
	}
	if err := wait(s)(xenapi.VDI.AsyncDestroy(s, ref)); err != nil {
		return "", xapi.Sanitize(err)
	}
	for _, b := range vbds {
		if _, err := xenapi.VBD.Create(s, xenapi.VBDRecord{VM: b.VM, VDI: xenapi.VDIRef(moved), Userdevice: b.Userdevice, Bootable: b.Bootable,
			Mode: b.Mode, Type: b.Type, Unpluggable: b.Unpluggable, OtherConfig: b.OtherConfig, QosAlgorithmParams: map[string]string{}}); err != nil {
			return "", xapi.Sanitize(err)
		}
	}
	return moved, nil
}

// unpluggedVBDs reads the disk's VBDs, which must still be exactly these, none plugged.
func unpluggedVBDs(s *xenapi.Session, ref xenapi.VDIRef, want []xenapi.VBDRef) ([]xenapi.VBDRecord, error) {
	changed := errors.New("this disk is in use, or was attached to a VM meanwhile: shut down its VMs and try again")
	now, err := xenapi.VDI.GetVBDs(s, ref)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	if len(now) != len(want) {
		return nil, changed
	}
	var out []xenapi.VBDRecord
	for _, b := range now {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		if !slices.Contains(want, b) || vbd.CurrentlyAttached {
			return nil, changed
		}
		out = append(out, vbd)
	}
	return out, nil
}

// DeleteSRDisk is the SR Storage tab's Delete (DeleteVirtualDiskCommand with AllowMultipleVBDDelete), its rules
// checked again on the server: a disk of VM snapshots deletes those snapshots (VMSnapshotDeleteAction), any other
// disk is detached from its VMs and destroyed (DestroyDiskAction).
func DeleteSRDisk(c *xapi.Connection, vdiUUID string) error {
	s := c.Session
	ref, vdi, sr, err := diskOf(s, vdiUUID)
	if err != nil {
		return err
	}
	info, err := DiskInfo(s, "", vdi, sr)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if !rules.CanDeleteSRDisk(info) {
		return errors.New("this disk cannot be deleted now (it is in use, busy, or its storage does not allow it)")
	}
	for _, b := range vdi.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil {
			return xapi.Sanitize(err)
		}
		if vbd.CurrentlyAttached {
			name, _ := xenapi.VM.GetNameLabel(s, vbd.VM)
			return errors.New(messages.Format(messages.Get("CANNOT_DELETE_VDI_ACTIVE_ON"), name))
		}
	}
	owners, err := DiskVMs(s, ref, vdi)
	if err != nil {
		return err
	}
	if vdi.IsASnapshot && len(owners) > 0 { // VDI.VDIType SNAPSHOT: StorageRules.DeleteActions deletes its VM snapshots
		var snapshots []string
		for _, vm := range owners {
			rec, err := xenapi.VM.GetRecord(s, vm)
			if err != nil {
				return xapi.Sanitize(err)
			}
			if rec.IsASnapshot {
				snapshots = append(snapshots, rec.UUID)
			}
		}
		if len(snapshots) == 0 { // the C# app runs no action then: the disk is never destroyed alone
			return errors.New("this snapshot disk belongs to no VM snapshot: it is not deleted")
		}
		if err := permitted(s, "vm.destroy", "vdi.destroy"); err != nil {
			return err
		}
		var errs []error
		for _, u := range snapshots {
			errs = append(errs, DeleteSnapshot(c, u))
		}
		return errors.Join(errs...)
	}
	if err := permitted(s, "vbd.unplug", "vbd.destroy", "vdi.destroy"); err != nil {
		return err
	}
	for _, b := range vdi.VBDs { // DetachVirtualDiskAction: the disk is attached nowhere (checked above)
		if err := wait(s)(xenapi.VBD.AsyncDestroy(s, b)); err != nil {
			return xapi.Sanitize(err)
		}
	}
	return xapi.Sanitize(wait(s)(xenapi.VDI.AsyncDestroy(s, ref)))
}

// TrimSR is SrTrimAction: the trim plugin's do_trim on the first server the SR is plugged into; a refusal is told in
// the words of its error code (FriendlyNames "Message.name-...") or its message.
func TrimSR(c *xapi.Connection, srUUID string) error {
	s := c.Session
	if err := permitted(s, "host.call_plugin"); err != nil {
		return err
	}
	ref, err := xenapi.SR.GetByUUID(s, srUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	sr, err := xenapi.SR.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	host := SRAttachedHost(s, sr)
	if host == "" {
		return errors.New(messages.Get("SR_TRIM_NO_STORAGE_HOST_ERROR"))
	}
	result, err := xenapi.Host.CallPlugin(s, host, "trim", "do_trim", map[string]string{"sr_uuid": sr.UUID})
	if err != nil {
		return xapi.Sanitize(err)
	}
	if strings.EqualFold(result, "true") {
		return nil
	}
	return errors.New(TrimError(result))
}

// TrimError is SrTrimAction.GetTrimError: the errcode (in FriendlyNames' words) or errmsg of a trim_response.
func TrimError(response string) string {
	var doc struct {
		Pairs []struct {
			Key   string `xml:"key"`
			Value string `xml:"value"`
		} `xml:"key_value_pair"`
	}
	code, msg := "", ""
	if xml.Unmarshal([]byte(response), &doc) == nil {
		for _, p := range doc.Pairs {
			switch p.Key {
			case "errcode":
				code = p.Value
			case "errmsg":
				msg = p.Value
			}
		}
	}
	if t := messages.Friendly("Message.name-" + code); code != "" && t != "" {
		return t
	}
	if msg != "" {
		return msg
	}
	return "unknown error"
}
