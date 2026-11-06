package gui

import (
	"path/filepath"
	"testing"

	"xcpng-center/internal/config"
	"xcpng-center/internal/trust"
)

// TestChangePasswordChecks: the dialog's checks come before any server call (PoolRules.ChangePasswordError).
func TestChangePasswordChecks(t *testing.T) {
	s, _ := config.Load(filepath.Join(t.TempDir(), "settings.toml"))
	w := NewWorkspace(s, &trust.Policy{Store: s})
	w.passwords["h:443"] = "old"
	if err := w.ChangePassword("h:443", "", "a", "a"); err == nil {
		t.Error("an empty field is accepted")
	}
	if err := w.ChangePassword("h:443", "wrong", "a", "a"); err == nil || err.Error() != "The password you typed was incorrect" {
		t.Errorf("wrong current: %v", err)
	}
	if err := w.ChangePassword("h:443", "old", "a", "b"); err == nil || err.Error() != "Passwords do not match" {
		t.Errorf("mismatch: %v", err)
	}
	if err := w.ChangePassword("h:443", "old", "a", "a"); err == nil { // no connection: refused, password kept
		t.Error("changed without a connection")
	}
	if w.passwords["h:443"] != "old" {
		t.Error("the kept password changed")
	}
}
