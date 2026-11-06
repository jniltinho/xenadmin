package rules

import (
	"slices"
	"strings"
	"testing"
)

func TestDRConfigureWarningMirrorsTheCSharp(t *testing.T) {
	if text, ok := DRConfigureWarning(false, 9); ok || text != "You can only select up to 8 SRs to store the VM configuration data." {
		t.Errorf("too many: %q %v", text, ok)
	}
	if text, ok := DRConfigureWarning(true, 0); !ok || !strings.HasPrefix(text, "Because no VM configuration data") {
		t.Errorf("all cleared: %q %v", text, ok)
	}
	if text, ok := DRConfigureWarning(false, 1); !ok || !strings.HasPrefix(text, "DR configuration is not complete") {
		t.Errorf("first SR: %q %v", text, ok)
	}
	for _, c := range [][2]int{{0, 0}, {1, 2}, {1, 8}} {
		if text, ok := DRConfigureWarning(c[0] == 1, c[1]); !ok || text != "" {
			t.Errorf("%v: %q %v", c, text, ok)
		}
	}
}

func TestDRWizardRules(t *testing.T) {
	if !slices.Equal(DRStartModes(DRTest), []string{DRStartPaused, DRStartNone}) || !slices.Equal(DRStartModes(DRFailover), []string{DRStart, DRStartPaused, DRStartNone}) {
		t.Error("start modes")
	}
	if !DRKnownType(DRFailback) || DRKnownType("") || DRKnownType("Test") {
		t.Error("known types")
	}
	if !DRMetadataSR(true, false, "user") || DRMetadataSR(false, false, "user") || DRMetadataSR(true, true, "") || DRMetadataSR(true, false, "iso") {
		t.Error("GetMetadataVDIsAction's SRs")
	}
	if !DRRealVM(false, false, false) || DRRealVM(true, false, false) || DRRealVM(false, true, false) || DRRealVM(false, false, true) {
		t.Error("IsRealVm")
	}
	if DRHACheck(DRTest, true, "p") == "" || DRHACheck(DRFailover, true, "p") != "" || DRHACheck(DRTest, false, "p") != "" {
		t.Error("HA check")
	}
	// an existing VM: not older is an error, running is an error, else a warning (it is replaced)
	if c := DRExistingCheck("vm", 3, 3, true, false); !c.Error || !strings.Contains(c.Text, "not older") {
		t.Errorf("same version: %+v", c)
	}
	if c := DRExistingCheck("vm", 2, 3, false, false); !c.Error || !strings.Contains(c.Text, "is running") {
		t.Errorf("running: %+v", c)
	}
	if c := DRExistingCheck("vm", 2, 3, true, false); c.Error || !strings.Contains(c.Text, "will be replaced") {
		t.Errorf("replaced: %+v", c)
	}
	if c := DRExistingCheck("vm", 2, 3, true, true); !c.Error || !strings.Contains(c.Text, "a test failover does not replace it") {
		t.Errorf("test failover: %+v", c)
	}
	if c := DRTestCollision("snapshot", "s1"); !c.Error || c.Text != "The snapshot 's1' already exists in this pool: a test failover does not replace it." {
		t.Errorf("collision: %+v", c)
	}
	if !DRBlocked([]DRCheck{{"w", false}, {"e", true}}) || DRBlocked([]DRCheck{{"w", false}}) || DRBlocked(nil) {
		t.Error("blocked")
	}
}
