package gui

import (
	"path/filepath"
	"testing"

	"xcpng-center/internal/config"
	"xcpng-center/internal/trust"
)

func workspace(t *testing.T) *Workspace {
	s, err := config.Load(filepath.Join(t.TempDir(), "settings.toml"))
	if err != nil {
		t.Fatal(err)
	}
	return NewWorkspace(s, &trust.Policy{Store: s, OperatorPins: map[string]string{}})
}

func TestTrustNeverReplacesAChangedCertificate(t *testing.T) {
	w := workspace(t)
	if err := w.Trust("h:443", "AA"); err != nil {
		t.Fatal(err)
	}
	if err := w.Trust("h:443", "BB"); err == nil {
		t.Fatal("a pinned certificate was replaced from the window")
	}
	if fp, _ := w.settings.Fingerprint("h:443"); fp != "AA" {
		t.Fatalf("pin = %q", fp)
	}
	// a pin another process saved after this window loaded its settings is not replaced either
	other, _ := config.Load(w.settings.Path())
	if err := other.Pin("k:443", "CC"); err != nil {
		t.Fatal(err)
	}
	if err := w.Trust("k:443", "DD"); err == nil {
		t.Fatal("another process's pin was replaced")
	}
	if fp, _ := w.settings.Fingerprint("k:443"); fp != "CC" {
		t.Fatalf("pin = %q", fp)
	}
}

func TestTrustRefusesAnOperatorPinMismatch(t *testing.T) {
	w := workspace(t)
	w.policy.OperatorPins["*"] = "aa:bb"
	if err := w.Trust("h:443", "CC"); err == nil {
		t.Fatal("a certificate other than --trust-fingerprint was pinned")
	}
	if err := w.Trust("h:443", "AABB"); err != nil {
		t.Fatal(err)
	}
}

func TestConnectAndObjectsErrors(t *testing.T) {
	w := workspace(t)
	if r := w.Connect("h:80", "root", "secret"); r.Error == "" || r.Connected {
		t.Fatalf("plain HTTP port accepted: %+v", r)
	}
	if _, err := w.Objects("h:443", "vms"); err == nil {
		t.Fatal("objects of a server that is not connected")
	}
	if rows := w.Servers(); len(rows) != 0 {
		t.Fatalf("servers = %+v", rows)
	}
}
