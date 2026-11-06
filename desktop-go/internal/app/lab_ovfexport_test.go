package app

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	xenapi "go/xenapi"

	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
)

// TestLabExportAppliance makes a halted VM (2 vCPUs, 256 MiB, a 32 MiB disk on lab-nfs, two interfaces), exports it as
// an OVA (manifest, verified) and as a compressed OVF folder, imports the OVA back with the appliance import and checks
// the copy; every VM, disk and vApp it made is removed, also on failure (XCP_ALLOW_WRITE=1).
func TestLabExportAppliance(t *testing.T) {
	c := labConnection(t)
	if os.Getenv("XCP_ALLOW_WRITE") != "1" {
		t.Skip("XCP_ALLOW_WRITE=1 not set")
	}
	s := c.Session
	srs, err := xenapi.SR.GetByNameLabel(s, "lab-nfs")
	if err != nil || len(srs) == 0 {
		t.Fatalf("no lab-nfs: %v", err)
	}
	srUUID, _ := xenapi.SR.GetUUID(s, srs[0])
	imp, err := ReadImportOptions(c, "")
	if err != nil || len(imp.Networks) == 0 {
		t.Fatalf("no network: %v", err)
	}
	name := fmt.Sprintf("xcpng-center-go-test-ovfexp-%d", os.Getpid())
	t.Cleanup(func() { // whatever the test made: the VMs of that name, their disks and the vApp
		vms, _ := xenapi.VM.GetByNameLabel(s, name)
		for _, vm := range vms {
			_ = xenapi.VM.HardShutdown(s, vm)
			vbds, _ := xenapi.VM.GetVBDs(s, vm)
			var vdis []xenapi.VDIRef
			for _, b := range vbds {
				if v, err := xenapi.VBD.GetVDI(s, b); err == nil && v != "OpaqueRef:NULL" {
					vdis = append(vdis, v)
				}
			}
			if err := xenapi.VM.Destroy(s, vm); err != nil {
				t.Errorf("VM not destroyed: %v", err)
			}
			for _, v := range vdis {
				if err := xenapi.VDI.Destroy(s, v); err != nil {
					t.Errorf("disk not destroyed: %v", err)
				}
			}
		}
		left, _ := xenapi.VDI.GetByNameLabel(s, name+" 0")
		for _, v := range left {
			if err := xenapi.VDI.Destroy(s, v); err != nil {
				t.Errorf("disk not destroyed: %v", err)
			}
		}
		apps, _ := xenapi.VMAppliance.GetByNameLabel(s, name)
		for _, a := range apps {
			if err := xenapi.VMAppliance.Destroy(s, a); err != nil {
				t.Errorf("vApp not destroyed: %v", err)
			}
		}
	})

	opts, err := ReadNewVMOptions(c)
	if err != nil {
		t.Fatal(err)
	}
	var tpl string
	for _, x := range opts.Templates {
		if x.Name == "Other install media" {
			tpl = x.UUID
		}
	}
	d, err := ReadNewVMDefaults(c, tpl, "")
	if err != nil {
		t.Fatal(err)
	}
	net := imp.Networks[0].UUID
	created, err := ops.CreateVM(c, ops.NewVM{Template: tpl, Name: name, Install: "cd", BootMode: d.BootMode, VCPUs: 2, Cores: 1, Memory: 256 << 20,
		Disks:    []ops.NewDisk{{Position: "0", Name: name + " 0", Size: 32 << 20, SR: srUUID}},
		Networks: []ops.NewVIF{{Network: net}, {Network: net}}})
	if err != nil {
		t.Fatal(err)
	}
	vm, err := xenapi.VM.GetByUUID(s, created.UUID)
	if err != nil {
		t.Fatal(err)
	}
	orig, err := xenapi.VM.GetRecord(s, vm)
	if err != nil {
		t.Fatal(err)
	}
	var macs []string
	for _, v := range orig.VIFs {
		r, _ := xenapi.VIF.GetRecord(s, v)
		macs = append(macs, r.Device+"="+r.MAC)
	}
	slices.Sort(macs)

	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "busy"), 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(dir, "busy", "x"), nil, 0o644)
	if _, err := ops.ExportAppliance(c, labPolicy(t), ops.ApplianceExport{Dir: dir, Name: "busy", VMs: []string{created.UUID}}); err == nil {
		t.Fatal("exported into a folder that is not empty")
	}
	ova, err := ops.ExportAppliance(c, labPolicy(t), ops.ApplianceExport{Dir: dir, Name: name, VMs: []string{created.UUID}, Manifest: true, OVA: true, Verify: true})
	if err != nil {
		t.Fatal(err)
	}
	if ova != filepath.Join(dir, name, name+".ova") {
		t.Fatalf("the package is %s", ova)
	}
	ovf, err := ops.ExportAppliance(c, labPolicy(t), ops.ApplianceExport{Dir: dir, Name: "folder", VMs: []string{created.UUID}, Manifest: true, Compress: true, Verify: true})
	if err != nil {
		t.Fatal(err)
	}
	if files, _ := os.ReadDir(filepath.Dir(ovf)); len(files) != 3 { // the descriptor, the manifest and the .vhd.gz
		t.Errorf("the OVF folder holds %v", files)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 3 { // busy, the OVA's folder, the OVF's: no temporary folder left
		t.Errorf("the target holds %v", entries)
	}

	pkg, err := rules.OpenPackage(ova)
	if err != nil {
		t.Fatal(err)
	}
	uuid, err := ops.ImportAppliance(c, labPolicy(t), ova, rules.DefaultApplianceMappings(pkg.Appliance, srUUID, net), true, false)
	if err != nil {
		t.Fatal(err)
	}
	back, err := xenapi.VM.GetByUUID(s, uuid)
	if err != nil {
		t.Fatal(err)
	}
	rec, err := xenapi.VM.GetRecord(s, back)
	if err != nil {
		t.Fatal(err)
	}
	if rec.NameLabel != name || rec.VCPUsAtStartup != orig.VCPUsAtStartup || rec.VCPUsMax != orig.VCPUsMax || rec.MemoryStaticMax != orig.MemoryDynamicMax ||
		rec.HVMBootPolicy != orig.HVMBootPolicy || rec.Platform["device-model"] != orig.Platform["device-model"] {
		t.Fatalf("the copy: %s, %d/%d vCPUs, %d memory, %q; the VM: %d/%d, %d, %q", rec.NameLabel, rec.VCPUsAtStartup, rec.VCPUsMax, rec.MemoryStaticMax,
			rec.HVMBootPolicy, orig.VCPUsAtStartup, orig.VCPUsMax, orig.MemoryDynamicMax, orig.HVMBootPolicy)
	}
	var sizes []int
	cds := 0
	for _, b := range rec.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil {
			continue
		}
		if vbd.Type == xenapi.VbdTypeCD {
			cds++
			continue
		}
		size, _ := xenapi.VDI.GetVirtualSize(s, vbd.VDI)
		sizes = append(sizes, size)
		t.Logf("disk %s: %d bytes, bootable %v", vbd.Userdevice, size, vbd.Bootable)
	}
	if len(sizes) != 1 || sizes[0] < 32<<20 || sizes[0] > 36<<20 || cds != 1 {
		t.Errorf("disks %v, %d DVD drives", sizes, cds)
	}
	var got []string
	for _, v := range rec.VIFs {
		r, _ := xenapi.VIF.GetRecord(s, v)
		got = append(got, r.Device+"="+r.MAC)
	}
	slices.Sort(got)
	if !slices.Equal(got, macs) {
		t.Errorf("interfaces %v, want %v", got, macs)
	}
}
