package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"xcpng-center/internal/config"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/trust"
	"xcpng-center/internal/xapi"
)

// labConnection logs in to XCP_LAB_HOST (a throwaway lab server) whose certificate is XCP_TRUST_FINGERPRINT.
// Without XCP_LAB_HOST and XCP_PASS the lab tests are skipped.
func labConnection(t *testing.T) *xapi.Connection {
	host, pass := os.Getenv("XCP_LAB_HOST"), os.Getenv("XCP_PASS")
	if host == "" || pass == "" {
		t.Skip("XCP_LAB_HOST and XCP_PASS not set")
	}
	s, err := config.Load(filepath.Join(t.TempDir(), "settings.toml"))
	if err != nil {
		t.Fatal(err)
	}
	h, p, err := xapi.ParseAddress(host)
	if err != nil {
		t.Fatal(err)
	}
	c, err := xapi.Connect(context.Background(), &trust.Policy{Store: s, OperatorPins: map[string]string{"*": os.Getenv("XCP_TRUST_FINGERPRINT")}}, h, p, "root", pass)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	return c
}

// TestLabVMTabs reads the Storage, Networking and Memory tabs of every VM of the lab (read-only).
func TestLabVMTabs(t *testing.T) {
	c := labConnection(t)
	vms, err := ListObjects(c, "vms")
	if err != nil {
		t.Fatal(err)
	}
	for _, vm := range vms {
		st, err := ListVMStorage(c, vm.UUID)
		if err != nil {
			t.Fatalf("%s storage: %v", vm.Name, err)
		}
		t.Logf("%s: disks %v, dvd %q iso %q, isos %v", vm.Name, st.Table.Rows, st.DVDName, st.ISO, st.ISOs)
		if _, err := ListVMNetworking(c, vm.UUID); err != nil {
			t.Fatalf("%s networking: %v", vm.Name, err)
		}
		if _, err := ListVMMemory(c, vm.UUID); err != nil {
			t.Fatalf("%s memory: %v", vm.Name, err)
		}
	}
}

// TestLabChangeISO inserts and ejects an ISO in the DVD drive of the lab's throwaway VM (XCP_LAB_VM, writes only
// with XCP_ALLOW_WRITE=1).
func TestLabChangeISO(t *testing.T) {
	c := labConnection(t)
	name := os.Getenv("XCP_LAB_VM")
	if name == "" || os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_LAB_VM and XCP_ALLOW_WRITE=1 not set")
	}
	vms, _ := ListObjects(c, "vms")
	for _, vm := range vms {
		if vm.Name != name {
			continue
		}
		st, err := ListVMStorage(c, vm.UUID)
		if err != nil || st.DVD == "" || len(st.ISOs) < 2 {
			t.Fatalf("no DVD drive or ISO: %+v %v", st, err)
		}
		iso := st.ISOs[len(st.ISOs)-1]
		if err := ops.ChangeISO(c, st.DVD, iso.VDI); err != nil {
			t.Fatal(err)
		}
		if st, _ = ListVMStorage(c, vm.UUID); st.ISO != iso.VDI {
			t.Fatalf("inserted %q, drive holds %q", iso.VDI, st.ISO)
		}
		if err := ops.ChangeISO(c, st.DVD, ""); err != nil {
			t.Fatal(err)
		}
		if st, _ = ListVMStorage(c, vm.UUID); st.ISO != "" {
			t.Fatalf("still holds %q after eject", st.ISO)
		}
		return
	}
	t.Fatalf("no VM named %q", name)
}
