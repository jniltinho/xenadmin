package gui

import (
	"xcpng-center/internal/app"
	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/xapi"
)

// The buttons of the Storage tabs of an SR (SrStorageView: Rescan, Add..., Move..., Delete, Properties) and of a
// server or pool (PhysicalStorageView: New SR..., Reclaim freed space, Properties).

// SRStorage is an SR's Storage tab with its buttons.
func (w *Workspace) SRStorage(key, srUUID string) (out app.SRStorage, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.SRStorageTab(c, srUUID); return err })
	return out, err
}

// PhysicalStorage is the Storage tab of a server ("hosts") or pool ("pools") with its buttons.
func (w *Workspace) PhysicalStorage(key, kind, uuid string) (out app.PhysicalStorage, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.PhysicalStorageTab(c, kind, uuid); return err })
	return out, err
}

// NewSRDiskOptions are the Add Virtual Disk dialog's defaults for a disk on this SR.
func (w *Workspace) NewSRDiskOptions(key, srUUID string) (out app.NewDisk, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.NewSRDiskOptions(c, srUUID); return err })
	return out, err
}

// CreateSRDisk makes a disk attached to no VM; it returns its uuid.
func (w *Workspace) CreateSRDisk(key, srUUID, name, description string, size int) (out string, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		out, err = ops.CreateSRDisk(c, srUUID, name, description, size)
		return err
	})
	return out, err
}

// RescanSRs scans the SRs (those detached are skipped).
func (w *Workspace) RescanSRs(key string, srUUIDs []string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.RescanSRs(c, srUUIDs...) })
}

// MoveDiskTargets is the Move Virtual Disk dialog's SR list.
func (w *Workspace) MoveDiskTargets(key, vdiUUID string) (out app.DiskTargets, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.MoveDiskTargets(c, vdiUUID); return err })
	return out, err
}

// MoveSRDisk moves or migrates the disk to the SR; it returns the disk's uuid there.
func (w *Workspace) MoveSRDisk(key, vdiUUID, srUUID string) (out string, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = ops.MoveSRDisk(c, vdiUUID, srUUID); return err })
	return out, err
}

// DeleteSRDisk deletes the disk (a snapshot disk: its snapshots), after the row's question.
func (w *Workspace) DeleteSRDisk(key, vdiUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.DeleteSRDisk(c, vdiUUID) })
}

// TrimSR reclaims the SR's freed space, after the CONFIRM_TRIM_SR question.
func (w *Workspace) TrimSR(key, srUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.TrimSR(c, srUUID) })
}

// TrimQuestion is Reclaim freed space's question.
func (w *Workspace) TrimQuestion() Confirm {
	return Confirm{Title: messages.Get("CONFIRM_TRIM_SR_TITLE"), Text: messages.Get("CONFIRM_TRIM_SR")}
}
