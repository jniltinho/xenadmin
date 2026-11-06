package app

import (
	"errors"
	"slices"
	"sort"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// The buttons of the Storage tab of an SR (XenAdmin.Desktop/SrStorageView.cs: Rescan, Add..., Move..., Delete,
// Properties) and of a server or pool (PhysicalStorageView.cs: New SR..., Reclaim freed space, Properties), with the
// rules of StorageRules.cs and the Move dialog's SR list (Dialogs.MoveDisk, StorageRules.MoveTargets).

// SRStorage is an SR's Storage tab with its buttons.
type SRStorage struct {
	Table     Table    `json:"table"`
	Disks     []SRDisk `json:"disks"`   // the rows' disks, as Table.Keys
	CanAdd    bool     `json:"can_add"` // StorageRules.CanAddDisk: the SR is not busy
	CanRescan bool     `json:"can_rescan"`
	RescanTip string   `json:"rescan_tip,omitempty"`
}

// SRDisk is what the buttons need about a row: Properties (CanEditDisk), Delete (CanDelete with several VMs allowed,
// and its DeleteConfirmation) and Move (UseMigrate).
type SRDisk struct {
	UUID        string `json:"uuid"`
	CanEdit     bool   `json:"can_edit"`
	CanDelete   bool   `json:"can_delete"`
	CanMove     bool   `json:"can_move"`
	DeleteTitle string `json:"delete_title"`
	DeleteText  string `json:"delete_text"`
}

// SRStorageTab is the Storage tab of an SR with its buttons.
func SRStorageTab(c *xapi.Connection, srUUID string) (SRStorage, error) {
	s := c.Session
	t, err := srDisks(s, srUUID)
	if err != nil {
		return SRStorage{}, err
	}
	ref, err := xenapi.SR.GetByUUID(s, srUUID)
	if err != nil {
		return SRStorage{}, xapi.Sanitize(err)
	}
	sr, err := xenapi.SR.GetRecord(s, ref)
	if err != nil {
		return SRStorage{}, xapi.Sanitize(err)
	}
	env, err := ops.ReadDiskEnv(s)
	if err != nil {
		return SRStorage{}, err
	}
	locked := len(sr.CurrentOperations) > 0 // SR.Locked
	out := SRStorage{Table: t, Disks: []SRDisk{}, CanAdd: !locked}
	out.CanRescan, out.RescanTip = rules.CanRescanSR(locked, ops.SRAttachedHost(s, sr) == "")
	for _, k := range t.Keys {
		row := SRDisk{UUID: k}
		vref, err := xenapi.VDI.GetByUUID(s, k)
		if err != nil {
			return SRStorage{}, xapi.Sanitize(err)
		}
		v, err := xenapi.VDI.GetRecord(s, vref)
		if err != nil {
			return SRStorage{}, xapi.Sanitize(err)
		}
		info, err := ops.DiskInfo(s, "", v, sr)
		if err != nil {
			return SRStorage{}, xapi.Sanitize(err)
		}
		move, err := env.Move(s, vref, v, sr)
		if err != nil {
			return SRStorage{}, err
		}
		_, row.CanMove = rules.UseMigrate(move)
		row.CanDelete = rules.CanDeleteSRDisk(info)
		row.CanEdit = !locked && !v.IsASnapshot && len(v.CurrentOperations) == 0 // StorageRules.CanEditDisk
		owners, err := ops.DiskVMs(s, vref, v)
		if err != nil {
			return SRStorage{}, err
		}
		row.DeleteTitle, row.DeleteText = rules.SRDiskQuestion(v.IsASnapshot && len(owners) > 0, sr.ContentType == "iso",
			v.Type == xenapi.VdiTypeSystem, v.NameLabel)
		out.Disks = append(out.Disks, row)
	}
	return out, nil
}

// NewSRDiskOptions is NewDiskDialog(connection, sr): the defaults and the SRs a disk can go to, this SR preselected
// when it can hold one.
func NewSRDiskOptions(c *xapi.Connection, srUUID string) (NewDisk, error) {
	out, err := newDiskOptions(c.Session, "")
	if err != nil {
		return out, err
	}
	out.DefaultSR = ""
	if slices.ContainsFunc(out.SRs, func(u SRUse) bool { return u.UUID == srUUID }) {
		out.DefaultSR = srUUID
	}
	return out, nil
}

// DiskTargets is the Move Virtual Disk dialog: live migration or offline move, and every SR with the reason it cannot
// take the disk.
type DiskTargets struct {
	Migrate bool    `json:"migrate"`
	SRs     []SRUse `json:"srs"`
}

// MoveDiskTargets is StorageRules.MoveTargets for a disk (SrPickerMoveItem, SrPickerMigrateItem): the shown SRs that
// can hold disks, by name, each with its free space or why it cannot take the disk.
func MoveDiskTargets(c *xapi.Connection, vdiUUID string) (DiskTargets, error) {
	s := c.Session
	ref, err := xenapi.VDI.GetByUUID(s, vdiUUID)
	if err != nil {
		return DiskTargets{}, xapi.Sanitize(err)
	}
	vdi, err := xenapi.VDI.GetRecord(s, ref)
	if err != nil {
		return DiskTargets{}, xapi.Sanitize(err)
	}
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return DiskTargets{}, xapi.Sanitize(err)
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return DiskTargets{}, xapi.Sanitize(err)
	}
	env, err := ops.ReadDiskEnv(s)
	if err != nil {
		return DiskTargets{}, err
	}
	move, err := env.Move(s, ref, vdi, srs[vdi.SR])
	if err != nil {
		return DiskTargets{}, err
	}
	migrate, ok := rules.UseMigrate(move)
	if !ok {
		return DiskTargets{}, errors.New("this disk cannot be moved now")
	}
	type home struct {
		host    xenapi.HostRef
		running bool
	}
	vms, err := ops.DiskVMs(s, ref, vdi)
	if err != nil {
		return DiskTargets{}, err
	}
	var homes []home // the disk's VMs with a home server (VM.Home)
	for _, vmRef := range vms {
		vm, err := xenapi.VM.GetRecord(s, vmRef)
		if err != nil {
			return DiskTargets{}, xapi.Sanitize(err)
		}
		if h := vmHome(s, vm); h != "" {
			homes = append(homes, home{h, vm.PowerState == xenapi.VMPowerStateRunning})
		}
	}
	type row struct {
		ref  xenapi.SRRef
		rec  xenapi.SRRecord
		name string
	}
	var list []row
	for r, sr := range srs {
		if sr.Type == "rawhba" || sr.ContentType == "iso" || sr.Type == "tmpfs" || !slices.Contains(env.Caps(sr.Type), "VDI_CREATE") || !srShown(sr) {
			continue
		}
		list = append(list, row{r, sr, srI18N(s, hosts, sr, "name_label", sr.NameLabel)})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].name < list[j].name })
	out := DiskTargets{Migrate: migrate, SRs: []SRUse{}}
	for _, r := range list {
		t := rules.MoveTarget{Current: r.ref == vdi.SR, Detached: ops.SRAttachedHost(s, r.rec) == "", Local: localSRTypes[r.rec.Type],
			Migrates: env.Migrates(r.rec), Homes: len(homes), CannotFit: cannotFitDisks(r.rec, env.Caps(r.rec.Type), []xenapi.VDIRecord{vdi})}
		if t.Local && migrate { // only a migration to a local SR looks at the VMs' servers
			for _, h := range homes {
				if !seenFrom(s, r.rec, h.host) {
					t.Unseen, t.UnseenRunning = hosts[h.host].NameLabel, h.running
					break
				}
			}
		}
		reason := rules.MoveTargetReason(t, migrate)
		text := reason
		if text == "" {
			text = messages.Format(messages.Get("SRPICKER_DISK_FREE"), diskSize(r.rec.PhysicalSize-r.rec.PhysicalUtilisation), diskSize(r.rec.PhysicalSize))
		}
		out.SRs = append(out.SRs, SRUse{UUID: r.rec.UUID, Text: r.name + "   " + text, Reason: reason})
	}
	return out, nil
}

// PhysicalStorage is the Storage tab of a server or pool with its buttons.
type PhysicalStorage struct {
	Table  Table        `json:"table"`
	SRs    []PhysicalSR `json:"srs"`     // the rows' SRs, as Table.Keys
	CanNew bool         `json:"can_new"` // New SR...: the server or pool is not busy
	Pool   string       `json:"pool"`    // the New Storage wizard's title
}

// PhysicalSR is Reclaim freed space for a row (StorageRules.CanTrim and its tooltip).
type PhysicalSR struct {
	UUID    string `json:"uuid"`
	CanTrim bool   `json:"can_trim"`
	TrimTip string `json:"trim_tip,omitempty"`
}

// PhysicalStorageTab is the Storage tab of a server ("hosts") or pool ("pools") with its buttons.
func PhysicalStorageTab(c *xapi.Connection, kind, uuid string) (PhysicalStorage, error) {
	s := c.Session
	host, busy, err := ownerBusy(s, kind, uuid)
	if err != nil {
		return PhysicalStorage{}, err
	}
	d, err := readHostTabData(s)
	if err != nil {
		return PhysicalStorage{}, err
	}
	t, err := d.storage(host)
	if err != nil {
		return PhysicalStorage{}, err
	}
	env, err := ops.ReadDiskEnv(s)
	if err != nil {
		return PhysicalStorage{}, err
	}
	out := PhysicalStorage{Table: t, SRs: []PhysicalSR{}, CanNew: !busy, Pool: PoolName(c)}
	for _, k := range t.Keys {
		row := PhysicalSR{UUID: k}
		if ref, err := xenapi.SR.GetByUUID(s, k); err == nil {
			if sr, err := xenapi.SR.GetRecord(s, ref); err == nil {
				trims := env.Trims(sr)
				row.CanTrim = trims && ops.SRAttachedHost(s, sr) != ""
				if !trims {
					row.TrimTip = messages.Get("TOOLTIP_SR_TRIM_UNSUPPORTED")
				}
			}
		}
		out.SRs = append(out.SRs, row)
	}
	return out, nil
}
