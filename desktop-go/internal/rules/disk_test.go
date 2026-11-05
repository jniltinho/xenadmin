package rules

import (
	"strings"
	"testing"
)

func TestDiskRulesMirrorTheCSharp(t *testing.T) {
	if title, text := DiskQuestion(true, true); title == "" || !strings.Contains(strings.ToLower(text+title), "delete") {
		t.Errorf("delete system: %q %q", title, text)
	}
	if _, text := DiskQuestion(false, false); !strings.HasPrefix(text, "Detaching this virtual disk will remove the association") {
		t.Errorf("detach: %q", text)
	}
	if UniqueName("New virtual disk", map[string]bool{"New virtual disk": true, "New virtual disk (1)": true}) != "New virtual disk (2)" {
		t.Error("UniqueName")
	}
	free := DiskInfo{DestroyAllowed: true, VBDs: []DiskVBD{{ThisVM: true}}}
	if !CanDeleteDisk(free) || !CanDetachDisk(free) {
		t.Error("an unattached disk can be deleted and detached")
	}
	attached := DiskInfo{DestroyAllowed: true, VBDs: []DiskVBD{{ThisVM: true, Attached: true, VMRunning: true, CanUnplug: false}}}
	if CanDeleteDisk(attached) || CanDetachDisk(attached) {
		t.Error("an attached disk that cannot be unplugged")
	}
	attached.VBDs[0].CanUnplug = true
	if !CanDetachDisk(attached) || CanDeleteDisk(attached) {
		t.Error("an unpluggable attached disk detaches but is not deleted while attached")
	}
	if CanDeleteDisk(DiskInfo{DestroyAllowed: true, VBDs: []DiskVBD{{}, {}}}) || CanDeleteDisk(DiskInfo{DestroyAllowed: true, SRTools: true}) {
		t.Error("shared or tools disks")
	}
}
