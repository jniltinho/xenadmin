package rules

import "testing"

func TestVMRulesMirrorTheCSharpTable(t *testing.T) {
	allowed := []string{"clean_shutdown", "hard_reboot"}
	for op, want := range map[VMOp]bool{VMShutdown: true, VMForceReboot: true, VMStart: false, VMReboot: false} {
		if VMAllowed(allowed, op) != want {
			t.Errorf("%s allowed = %v", op, !want)
		}
	}
	if NeedsConfirmation(VMStart) || NeedsConfirmation(VMResume) || !NeedsConfirmation(VMForceShutdown) {
		t.Error("only Start and Resume skip the confirmation")
	}
	if _, text := Confirmation(VMShutdown, true); text == "" || text[:24] != "The selected VM is curre" {
		t.Errorf("HA shutdown text = %q", text)
	}
	for _, op := range VMOps {
		if title, _ := Confirmation(op, false); NeedsConfirmation(op) != (title != "") {
			t.Errorf("%s: confirmation title %q", op, title)
		}
	}
	if HAStartBlocker(true, true, []bool{true, false}, "vm") == "" || HAStartBlocker(true, false, []bool{false}, "vm") != "" ||
		HAStartBlocker(false, true, []bool{false}, "vm") != "" || HAStartBlocker(true, true, []bool{true}, "vm") != "" {
		t.Error("HAStartBlocker")
	}
	if Known("delete") || !Known(VMResume) {
		t.Error("Known")
	}
}
