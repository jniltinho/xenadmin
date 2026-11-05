package rules

import (
	"strings"
	"testing"
)

func TestSnapshotRulesMirrorTheCSharp(t *testing.T) {
	m := SnapshotModesFor([]string{"snapshot", "clean_shutdown"}, true, false)
	if !m.Disk || m.Memory || m.MemoryInfo != "Disk and memory snapshots need I/O drivers to be installed" {
		t.Errorf("%+v", m)
	}
	if m := SnapshotModesFor(nil, false, true); m.MemoryInfo != "Disk and memory snapshots need the VM to be running" {
		t.Errorf("%+v", m)
	}
	if SnapshotTypeText("Suspended") != "Disks and memory" || SnapshotTypeText("Halted") != "Disks only" {
		t.Error("type text")
	}
	text, option := RevertQuestion("s1", true)
	if !strings.HasPrefix(text, "Are you sure you want to revert to 's1'?\n\nReverting to this snapshot") || option != "Take a snapshot of the VM's current state and then revert" {
		t.Errorf("%q %q", text, option)
	}
	if _, option := RevertQuestion("s1", false); option != "" {
		t.Error("snapshot-first offered when not allowed")
	}
	if DeleteSnapshotQuestion("s1") != "Are you sure you want to delete the snapshot 's1'?  This operation cannot be undone." {
		t.Error(DeleteSnapshotQuestion("s1"))
	}
}
