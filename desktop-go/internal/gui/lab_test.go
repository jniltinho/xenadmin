package gui

import (
	"os"
	"path/filepath"
	"testing"

	xenapi "go/xenapi"

	"xcpng-center/internal/config"
	"xcpng-center/internal/trust"
	"xcpng-center/internal/xapi"
)

// TestLabRestartCoordinatorToolstack restarts the toolstack of the lab's coordinator from the window's workspace: the
// command returns once the workspace has logged in again (RestartToolstackAction), and the connection works
// (XCP_LAB_HOST, XCP_PASS, XCP_TRUST_FINGERPRINT, XCP_ALLOW_WRITE=1).
func TestLabRestartCoordinatorToolstack(t *testing.T) {
	host, pass := os.Getenv("XCP_LAB_HOST"), os.Getenv("XCP_PASS")
	if host == "" || pass == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_HOST, XCP_PASS and XCP_ALLOW_WRITE=1 not set")
	}
	s, err := config.Load(filepath.Join(t.TempDir(), "settings.toml"))
	if err != nil {
		t.Fatal(err)
	}
	h, p, err := xapi.ParseAddress(host)
	if err != nil {
		t.Fatal(err)
	}
	w := NewWorkspace(s, &trust.Policy{Store: s, OperatorPins: map[string]string{xapi.HostKey(h, p): os.Getenv("XCP_TRUST_FINGERPRINT")}})
	defer w.Close()
	r := w.Connect(host, "root", pass)
	if !r.Connected {
		t.Fatalf("%+v", r)
	}
	var coordinator string
	if err := w.with(r.Key, func(c *xapi.Connection) error {
		pools, err := xenapi.Pool.GetAllRecords(c.Session)
		for _, p := range pools {
			coordinator, _ = xenapi.Host.GetUUID(c.Session, p.Master)
		}
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if err := w.HostCommand(r.Key, coordinator, "restart_toolstack", false); err != nil {
		t.Fatal(err)
	}
	if objs, err := w.Objects(r.Key, "hosts"); err != nil || len(objs) == 0 {
		t.Fatalf("not connected again: %v %v", objs, err)
	}
}

// TestLabChangePassword runs Change Server Password with the same password (so the lab keeps its password): the
// server takes it and the connection goes on (XCP_ALLOW_WRITE=1).
func TestLabChangePassword(t *testing.T) {
	host, pass := os.Getenv("XCP_LAB_HOST"), os.Getenv("XCP_PASS")
	if host == "" || pass == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_HOST, XCP_PASS and XCP_ALLOW_WRITE=1 not set")
	}
	s, _ := config.Load(filepath.Join(t.TempDir(), "settings.toml"))
	h, p, _ := xapi.ParseAddress(host)
	w := NewWorkspace(s, &trust.Policy{Store: s, OperatorPins: map[string]string{xapi.HostKey(h, p): os.Getenv("XCP_TRUST_FINGERPRINT")}})
	defer w.Close()
	r := w.Connect(host, "root", pass)
	if !r.Connected {
		t.Fatalf("%+v", r)
	}
	if ok, err := w.CanChangePassword(r.Key); err != nil || !ok {
		t.Fatal("root may not change the password", err)
	}
	if err := w.ChangePassword(r.Key, pass, pass, pass); err != nil {
		t.Fatal(err)
	}
	if _, err := w.Objects(r.Key, "hosts"); err != nil {
		t.Fatal(err)
	}
}
