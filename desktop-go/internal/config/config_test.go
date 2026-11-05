package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSaveLoadAndPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "settings.toml")
	s, err := Load(path)
	if err != nil || len(s.Servers) != 0 {
		t.Fatal("a missing file is an empty configuration", err)
	}
	if err := s.AddServer(Server{Host: "xcp1", Port: 443, User: "root"}); err != nil {
		t.Fatal(err)
	}
	if err := s.AddServer(Server{Host: "XCP1", Port: 443, User: "admin", Name: "pool"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Pin("XCP1:443", "AB"); err != nil {
		t.Fatal(err)
	}
	again, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(again.Servers) != 1 || again.Servers[0].User != "admin" || again.Servers[0].Name != "pool" {
		t.Errorf("servers: %+v", again.Servers)
	}
	if fp, ok := again.Fingerprint("xcp1:443"); !ok || fp != "AB" {
		t.Error("fingerprints are kept, host keys case-insensitive")
	}
	data, _ := os.ReadFile(path)
	if strings.Contains(strings.ToLower(string(data)), "pass") {
		t.Error("no password in the settings file")
	}
	if runtime.GOOS != "windows" {
		if st, _ := os.Stat(path); st.Mode().Perm() != 0o600 {
			t.Errorf("mode %v, want 0600", st.Mode().Perm())
		}
	}
}

func TestTwoProcessesKeepEachOthersPins(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.toml")
	a, _ := Load(path)
	b, _ := Load(path) // a second process with its own copy
	if err := a.Pin("one:443", "AA"); err != nil {
		t.Fatal(err)
	}
	if err := b.Pin("two:443", "BB"); err != nil {
		t.Fatal(err)
	}
	if err := b.AddServer(Server{Host: "two", Port: 443, User: "root"}); err != nil {
		t.Fatal(err)
	}
	final, _ := Load(path)
	if _, ok := final.Fingerprint("one:443"); !ok {
		t.Error("the first process's pin was lost")
	}
	if _, ok := final.Fingerprint("two:443"); !ok || len(final.Servers) != 1 {
		t.Error("the second process's changes are kept")
	}
	if _, err := os.Stat(path + ".lock"); err == nil {
		t.Error("the lock file is released")
	}
}
