package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	xenapi "go/xenapi"

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

// TestLabDiskLifecycle adds a 64 MB disk to the throwaway VM, detaches it and deletes it.
func TestLabDiskLifecycle(t *testing.T) {
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
		opts, err := NewDiskOptions(c, vm.UUID)
		if err != nil || opts.DefaultSR == "" || opts.Name == "" {
			t.Fatalf("options %+v %v", opts, err)
		}
		notice, err := ops.AddDisk(c, vm.UUID, "go-lab-disk", "lab test", opts.DefaultSR, 64<<20)
		if err != nil {
			t.Fatal(err)
		}
		st, _ := ListVMStorage(c, vm.UUID)
		var d Disk
		for _, x := range st.Disks {
			if x.Name == "go-lab-disk" {
				d = x
			}
		}
		if d.VBD == "" || (notice == "" && !d.Attached) || d.System {
			t.Fatalf("added disk %+v (notice %q)", d, notice)
		}
		if !d.CanDetach || d.CanDelete != !d.Attached {
			t.Fatalf("rules for the new disk: %+v", d)
		}
		if d.Attached {
			if err := ops.SetDiskActive(c, d.VBD, false); err != nil {
				t.Fatal(err)
			}
		}
		if err := ops.DeleteDisk(c, vm.UUID, d.VBD); err != nil {
			t.Fatal(err)
		}
		st, _ = ListVMStorage(c, vm.UUID)
		for _, x := range st.Disks {
			if x.Name == "go-lab-disk" {
				t.Fatal("the disk is still attached")
			}
		}
		// Detach keeps the disk: it leaves the VM and stays on the SR (removed here by the test)
		if _, err := ops.AddDisk(c, vm.UUID, "go-lab-detach", "", opts.DefaultSR, 64<<20); err != nil {
			t.Fatal(err)
		}
		st, _ = ListVMStorage(c, vm.UUID)
		for _, x := range st.Disks {
			if x.Name == "go-lab-detach" {
				d = x
			}
		}
		if err := ops.DetachDisk(c, d.VBD); err != nil {
			t.Fatal(err)
		}
		// the detached disk is offered by Attach Disk, and attaches again, then goes for good (read-only is refused by XAPI for HVM guests)
		choices, err := AttachableDisks(c, vm.UUID)
		found := false
		for _, ch := range choices {
			found = found || ch.VDI == d.VDI
		}
		if err != nil || !found {
			t.Fatalf("the detached disk is not offered: %v %v", choices, err)
		}
		if _, err := ops.AttachDisk(c, vm.UUID, d.VDI, false); err != nil {
			t.Fatal(err)
		}
		if _, err := ops.AttachDisk(c, vm.UUID, d.VDI, false); err == nil {
			t.Fatal("the same disk attached twice")
		}
		st, _ = ListVMStorage(c, vm.UUID)
		for _, x := range st.Disks {
			if x.VDI == d.VDI {
				d = x
			}
		}
		if d.Attached {
			_ = ops.SetDiskActive(c, d.VBD, false)
		}
		if err := ops.DetachDisk(c, d.VBD); err != nil {
			t.Fatal(err)
		}
		vdi, err := xenapi.VDI.GetByUUID(c.Session, d.VDI)
		if err != nil {
			t.Fatalf("the detached disk is gone: %v", err)
		}
		if vbds, _ := xenapi.VDI.GetVBDs(c.Session, vdi); len(vbds) != 0 {
			t.Fatal("the detached disk still has a VBD")
		}
		if err := xenapi.VDI.Destroy(c.Session, vdi); err != nil {
			t.Fatal(err)
		}
		return
	}
	t.Fatalf("no VM named %q", name)
}

// TestLabProperties reads the throwaway VM's Properties and changes its description and home server, then restores
// them.
func TestLabProperties(t *testing.T) {
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
		p, err := ReadVMProperties(c, vm.UUID)
		if err != nil || len(p.VCPUChoices) == 0 || len(p.Servers) == 0 {
			t.Fatalf("%+v %v", p, err)
		}
		desc, home := "go lab test", p.Servers[0].UUID
		if err := ops.SaveVMProperties(c, vm.UUID, ops.VMChanges{Description: &desc, HomeServer: &home}); err != nil {
			t.Fatal(err)
		}
		q, _ := ReadVMProperties(c, vm.UUID)
		if q.Description != desc || q.HomeServer != home {
			t.Fatalf("saved %+v", q)
		}
		// a running VM keeps its vCPU weight and topology (only a halted one changes them)
		prio := q.Priority + 1
		if err := ops.SaveVMProperties(c, vm.UUID, ops.VMChanges{Priority: &prio}); err == nil && q.Running {
			t.Fatal("vCPU priority changed on a running VM")
		}
		if err := ops.SaveVMProperties(c, vm.UUID, ops.VMChanges{Description: &p.Description, HomeServer: &p.HomeServer}); err != nil {
			t.Fatal(err)
		}
		return
	}
	t.Fatalf("no VM named %q", name)
}

// TestLabCreateVM creates a VM from "Other install media" with the lab's Alpine ISO and a 1 GB disk, starts it, then
// destroys it (XCP_ALLOW_WRITE=1; the VM is named xcpng-center-go-newvm).
func TestLabCreateVM(t *testing.T) {
	c := labConnection(t)
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	opts, err := ReadNewVMOptions(c)
	if err != nil {
		t.Fatal(err)
	}
	var tpl, iso string
	for _, x := range opts.Templates {
		if x.Name == "Other install media" {
			tpl = x.UUID
		}
	}
	for _, x := range opts.ISOs {
		if strings.Contains(x.Text, "alpine") {
			iso = x.VDI
		}
	}
	if tpl == "" || iso == "" || opts.DefaultSR == "" {
		t.Fatalf("template %q iso %q sr %q", tpl, iso, opts.DefaultSR)
	}
	d, err := ReadNewVMDefaults(c, tpl, "")
	if err != nil || !d.HVM || d.Install != "cd" {
		t.Fatalf("%+v %v", d, err)
	}
	spec := ops.NewVM{Template: tpl, Name: "xcpng-center-go-newvm", Install: "cd", ISO: iso, BootMode: d.BootMode, VCPUs: 1, Cores: 1,
		Memory: 256 << 20, Networks: d.Networks, Start: true,
		Disks: []ops.NewDisk{{Position: "0", Name: "xcpng-center-go-newvm 0", Size: 1 << 30, SR: opts.DefaultSR}}}
	created, err := ops.CreateVM(c, spec)
	uuid := created.UUID
	if uuid != "" {
		defer func() {
			ref, _ := xenapi.VM.GetByUUID(c.Session, uuid)
			rec, _ := xenapi.VM.GetRecord(c.Session, ref)
			if rec.PowerState != xenapi.VMPowerStateHalted {
				_ = xenapi.VM.HardShutdown(c.Session, ref)
			}
			for _, b := range rec.VBDs {
				if vbd, err := xenapi.VBD.GetRecord(c.Session, b); err == nil && vbd.OtherConfig["owner"] != "" {
					_ = xenapi.VDI.Destroy(c.Session, vbd.VDI)
				}
			}
			_ = xenapi.VM.Destroy(c.Session, ref)
		}()
	}
	if err != nil || created.StartError != "" {
		t.Fatal(err, created.StartError)
	}
	vms, _ := ListObjects(c, "vms")
	for _, v := range vms {
		if v.UUID == uuid {
			if v.Name != spec.Name || v.PowerState != "Running" {
				t.Fatalf("created %+v", v)
			}
			st, _ := ListVMStorage(c, uuid)
			if len(st.Disks) != 1 || st.ISO != iso {
				t.Fatalf("storage %+v", st)
			}
			return
		}
	}
	t.Fatal("the new VM is not listed")
}
