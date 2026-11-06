package rules

import (
	"strings"
	"testing"
)

func TestUseMigrate(t *testing.T) {
	for _, c := range []struct {
		name        string
		d           DiskMove
		migrate, ok bool
	}{
		{"unattached disk: offline move", DiskMove{VMsHalted: true, SRMigrates: true}, false, true},
		{"attached, SR migrates: live", DiskMove{VBDs: 1, SRMigrates: true}, true, true},
		{"attached, licence forbids, halted: move", DiskMove{VBDs: 1, SRMigrates: true, MigrateForbidden: true, VMsHalted: true}, false, true},
		{"attached to a running VM, SR cannot migrate", DiskMove{VBDs: 1}, false, false},
		{"snapshot", DiskMove{Snapshot: true, VMsHalted: true}, false, false},
		{"CBT", DiskMove{CBT: true, VMsHalted: true}, false, false},
		{"HA disk", DiskMove{HAType: true, VMsHalted: true}, false, false},
		{"LUN per VDI", DiskMove{LunPerVDI: true, VMsHalted: true, VBDs: 1, SRMigrates: true}, false, false},
		{"busy", DiskMove{Locked: true, VMsHalted: true}, false, false},
	} {
		if m, ok := UseMigrate(c.d); m != c.migrate || ok != c.ok {
			t.Errorf("%s: got %v %v, want %v %v", c.name, m, ok, c.migrate, c.ok)
		}
	}
}

func TestMoveTargetReason(t *testing.T) {
	if r := MoveTargetReason(MoveTarget{}, false); r != "" {
		t.Errorf("free SR: %q", r)
	}
	if r := MoveTargetReason(MoveTarget{CannotFit: "full"}, false); r != "full" {
		t.Errorf("full SR: %q", r)
	}
	if r := MoveTargetReason(MoveTarget{Current: true}, false); r != "Current location" {
		t.Errorf("current: %q", r)
	}
	if r := MoveTargetReason(MoveTarget{Local: true}, false); r != "" {
		t.Errorf("a move to a local SR needs no home: %q", r)
	}
	if r := MoveTargetReason(MoveTarget{Local: true, Migrates: true}, true); !strings.Contains(r, "home server") {
		t.Errorf("migrate to local, no home: %q", r)
	}
	if r := MoveTargetReason(MoveTarget{Local: true, Migrates: true, Homes: 1, Unseen: "h2"}, true); !strings.Contains(r, "h2") {
		t.Errorf("migrate to a local SR not seen from h2: %q", r)
	}
	if r := MoveTargetReason(MoveTarget{Homes: 1}, true); r != "Unsupported SR type" {
		t.Errorf("migrate to an SR without mirroring: %q", r)
	}
}

func TestCanDeleteSRDisk(t *testing.T) {
	two := DiskInfo{DestroyAllowed: true, VBDs: []DiskVBD{{}, {}}}
	if CanDeleteDisk(two) || !CanDeleteSRDisk(two) {
		t.Error("a disk of two halted VMs: only the SR tab deletes it")
	}
	if CanDeleteSRDisk(DiskInfo{DestroyAllowed: true, Locked: true}) {
		t.Error("a busy disk is not deleted")
	}
	if CanDeleteSRDisk(DiskInfo{DestroyAllowed: true, VBDs: []DiskVBD{{}, {Attached: true}}}) {
		t.Error("an attached disk is not deleted")
	}
}

func TestSRDiskQuestion(t *testing.T) {
	if title, _ := SRDiskQuestion(true, false, false, "d"); title != "Delete entire snapshot" {
		t.Errorf("snapshot: %q", title)
	}
	if _, text := SRDiskQuestion(false, true, false, "x.iso"); !strings.Contains(text, "x.iso") {
		t.Errorf("iso: %q", text)
	}
	if a, _ := SRDiskQuestion(false, false, false, "d"); a == "" {
		t.Error("disk: no title")
	}
}
