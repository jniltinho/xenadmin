package rules

import "xcpng-center/internal/messages"

// The SR Storage tab's rules (XenAdmin.Desktop.Core/StorageRules.cs: CanRescan, CanMove, CanMigrate, UseMigrate,
// MoveTargets, DeleteConfirmation).

// CanRescanSR is StorageRules.CanRescan: not while an operation runs on the SR, nor on a detached SR (its tooltip).
// shortcut: no "scan in progress" tooltip; the tab's buttons stay disabled while its own Rescan runs.
func CanRescanSR(locked, detached bool) (bool, string) {
	switch {
	case locked:
		return false, ""
	case detached:
		return false, messages.Get("SR_DETACHED")
	}
	return true, ""
}

// DiskMove is what MoveVirtualDiskCommand.CanBeMoved and MigrateVirtualDiskCommand.CanBeMigrated read about a disk.
type DiskMove struct {
	Snapshot, Locked, HAType, CBT bool
	LunPerVDI                     bool // its SR has one LUN per disk (SR.HBALunPerVDI)
	VBDs                          int
	VMsHalted                     bool // every VM of the disk is halted
	SRMigrates                    bool // its SR supports storage migration (VDI_SNAPSHOT and VDI_MIRROR)
	MigrateForbidden              bool // the license forbids storage live migration (restrict_storage_xen_motion)
}

// UseMigrate is StorageRules.UseMigrate: ok when the disk can go to another SR, migrate when live (VDI.pool_migrate),
// else an offline move (copy and delete).
func UseMigrate(d DiskMove) (migrate, ok bool) {
	if d.Snapshot || d.Locked || d.HAType || d.CBT || d.LunPerVDI {
		return false, false
	}
	if d.VBDs > 0 && d.SRMigrates && !d.MigrateForbidden {
		return true, true
	}
	return false, d.VMsHalted
}

// MoveTarget is what StorageRules.MoveTargets reads about an SR for a disk.
type MoveTarget struct {
	Current, Detached, Local, Migrates bool
	Homes                              int    // the disk's VMs that have a home server
	Unseen                             string // a home server that cannot see the SR ("" none)
	UnseenRunning                      bool   // that home's VM is running
	CannotFit                          string // SR.CanFitDisks' reason ("" fits)
}

// MoveTargetReason is why an SR cannot take the disk (SrPickerMoveItem, SrPickerMigrateItem), "" when it can.
func MoveTargetReason(t MoveTarget, migrate bool) string {
	switch {
	case t.Current:
		return messages.Get("CURRENT_LOCATION")
	case t.Detached:
		return messages.Get("SR_DETACHED")
	case migrate && t.Local && t.Homes == 0:
		return messages.Get("SR_IS_LOCAL")
	case migrate && t.Local && t.Unseen != "" && t.UnseenRunning:
		return messages.Get("SRPICKER_ERROR_LOCAL_SR_MUST_BE_RESIDENT_HOSTS")
	case migrate && t.Local && t.Unseen != "":
		return messages.Format(messages.Get("SR_CANNOT_BE_SEEN"), t.Unseen)
	case migrate && !t.Migrates:
		return messages.Get("UNSUPPORTED_SR_TYPE")
	}
	return t.CannotFit
}

// SRDiskQuestion is StorageRules.DeleteConfirmation (title, text): a snapshot disk deletes its whole snapshot, an ISO,
// a system disk, else a virtual disk.
func SRDiskQuestion(snapshot, iso, system bool, name string) (title, text string) {
	switch {
	case snapshot:
		return messages.Get("MESSAGEBOX_DELETE_SNAPSHOT_TITLE"), messages.Get("MESSAGEBOX_DELETE_SNAPSHOT")
	case iso:
		return messages.Get("MESSAGEBOX_DELETE_ISO_TITLE"), messages.Format(messages.Get("MESSAGEBOX_DELETE_ISO"), name)
	}
	return DiskQuestion(true, system)
}
