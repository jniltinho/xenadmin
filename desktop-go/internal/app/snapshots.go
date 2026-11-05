package app

import (
	"sort"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// Snapshot is one row of the Snapshots tab (SnapshotsPage).
type Snapshot struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Time        time.Time `json:"time"`
	Type        string    `json:"type"` // Disks only / Disks and memory
	Size        int64     `json:"size"` // SnapshotsPage.GetSnapshotSize: disks plus the suspend image
}

// VMSnapshots is the Snapshots tab of a VM: its snapshots (newest first) and the snapshot modes it allows now.
type VMSnapshots struct {
	Snapshots       []Snapshot          `json:"snapshots"`
	Modes           rules.SnapshotModes `json:"modes"`
	SnapshotAllowed bool                `json:"snapshot_allowed"` // RevertToSnapshotCommand's "snapshot first" option
}

// ListSnapshots reads a VM's snapshots and what the Take Snapshot dialog offers.
func ListSnapshots(c *xapi.Connection, vmUUID string) (VMSnapshots, error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return VMSnapshots{}, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return VMSnapshots{}, xapi.Sanitize(err)
	}
	allowed := make([]string, len(vm.AllowedOperations))
	for i, a := range vm.AllowedOperations {
		allowed[i] = string(a)
	}
	io := false
	if gm, err := xenapi.VMGuestMetrics.GetRecord(s, vm.GuestMetrics); err == nil {
		io = gm.PVDriversDetected
	}
	out := VMSnapshots{Snapshots: []Snapshot{}, Modes: rules.SnapshotModesFor(allowed, vm.PowerState == xenapi.VMPowerStateRunning, io)}
	out.SnapshotAllowed = out.Modes.Disk
	for _, sr := range vm.Snapshots {
		snap, err := xenapi.VM.GetRecord(s, sr)
		if err != nil {
			return VMSnapshots{}, xapi.Sanitize(err)
		}
		row := Snapshot{UUID: snap.UUID, Name: snap.NameLabel, Description: snap.NameDescription, Time: snap.SnapshotTime,
			Type: rules.SnapshotTypeText(string(snap.PowerState))}
		vdis := []xenapi.VDIRef{snap.SuspendVDI}
		for _, b := range snap.VBDs {
			if vbd, err := xenapi.VBD.GetRecord(s, b); err == nil && vbd.Type == xenapi.VbdTypeDisk {
				vdis = append(vdis, vbd.VDI)
			}
		}
		for _, v := range vdis {
			if v == "" || v == "OpaqueRef:NULL" {
				continue
			}
			if u, err := xenapi.VDI.GetPhysicalUtilisation(s, v); err == nil && u > 0 {
				row.Size += int64(u)
			}
		}
		out.Snapshots = append(out.Snapshots, row)
	}
	sort.Slice(out.Snapshots, func(i, j int) bool { return out.Snapshots[i].Time.After(out.Snapshots[j].Time) })
	return out, nil
}
