package rules

import (
	"strings"

	"xcpng-center/internal/messages"
)

// DiskQuestion is VmStorageView's Detach and Delete question (title, text); system is a bootable or system disk.
func DiskQuestion(del, system bool) (title, text string) {
	key := "MESSAGEBOX_DETACH_VD"
	switch {
	case del && system:
		key = "MESSAGEBOX_DELETE_SYS_DISK"
	case del:
		key = "MESSAGEBOX_DELETE_VD"
	case system:
		key = "MESSAGEBOX_DETACH_SYSTEMVD"
	}
	return messages.Get(key + "_TITLE"), strings.TrimSpace(messages.Get(key))
}

// UniqueName is Helpers.MakeUniqueName: the stub, or "stub (1)", "stub (2)"... when taken.
func UniqueName(stub string, taken map[string]bool) string {
	name := stub
	for i := 1; taken[name]; i++ {
		name = messages.Format(messages.Get("NEWVM_DEFAULTNAME"), stub, i)
	}
	return name
}

// DiskInfo is what DeleteVirtualDiskCommand.CanRun and DetachVirtualDiskCommand.CanRun read about a VM's disk.
type DiskInfo struct {
	SRPhysical, SRTools, UsedByHA bool // the disk's SR is a local/CD one or the tools SR; an HA metadata disk on HA
	VBDs                          []DiskVBD
	DestroyAllowed                bool // VDI allowed_operations has destroy (or the SR has one LUN per disk)
	System                        bool // VDI type system
	Locked                        bool // an operation runs on the disk (VDI.Locked: current_operations)
}

// DiskVBD is one VBD of the disk.
type DiskVBD struct {
	VMRunning, Attached, CanUnplug, ThisVM bool
}

// CanDeleteDisk is StorageRules.CanDelete (one VBD at most, none attached, not a running VM's system disk).
func CanDeleteDisk(d DiskInfo) bool { return canDelete(d, false) }

// CanDeleteSRDisk is StorageRules.CanDelete(vdi, allowMultipleVbds: true): the SR Storage tab also deletes a disk
// of several VMs.
func CanDeleteSRDisk(d DiskInfo) bool { return canDelete(d, true) }

func canDelete(d DiskInfo, allowMultipleVBDs bool) bool {
	if d.Locked || d.SRPhysical || d.SRTools || d.UsedByHA || (len(d.VBDs) > 1 && !allowMultipleVBDs) {
		return false
	}
	for _, b := range d.VBDs {
		if (d.System && b.VMRunning) || b.Attached {
			return false
		}
	}
	return d.DestroyAllowed
}

// CanDetachDisk is StorageRules.CanDetach: every attached VBD of this VM can be unplugged.
func CanDetachDisk(d DiskInfo) bool {
	for _, b := range d.VBDs {
		if b.ThisVM && b.Attached && !(b.VMRunning && b.CanUnplug) {
			return false
		}
	}
	return true
}
