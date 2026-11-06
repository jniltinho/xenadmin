package rules

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// testVHD is a dynamic VHD of size bytes (2 MiB blocks) whose first block is allocated, as export_raw_vdi writes one.
func testVHD(size int64) []byte {
	be := binary.BigEndian
	const block = 2 << 20
	entries := (size + block - 1) / block
	foot := make([]byte, 512)
	copy(foot, "conectix")
	be.PutUint64(foot[16:], 512)
	be.PutUint64(foot[40:], uint64(size))
	be.PutUint64(foot[48:], uint64(size))
	be.PutUint32(foot[60:], 3)
	be.PutUint32(foot[64:], vhdChecksum(foot, 64))
	head := make([]byte, 1024)
	copy(head, "cxsparse")
	be.PutUint64(head[8:], ^uint64(0))
	be.PutUint64(head[16:], 1536)
	be.PutUint32(head[28:], uint32(entries))
	be.PutUint32(head[32:], block)
	be.PutUint32(head[36:], vhdChecksum(head, 36))
	tableLen := (entries*4 + 511) / 512 * 512
	table := bytes.Repeat([]byte{0xff}, int(tableLen))
	data := int64(1536) + tableLen
	be.PutUint32(table, uint32(data/512))
	out := append(append(append([]byte{}, foot...), head...), table...)
	out = append(out, bytes.Repeat([]byte{0xff}, 512)...)   // the block's bitmap
	out = append(out, bytes.Repeat([]byte{0xa5}, block)...) // its data
	return append(out, foot...)
}

func TestVerifyVHD(t *testing.T) {
	good := testVHD(8 << 20)
	if err := VerifyVHD(bytes.NewReader(good), int64(len(good))); err != nil {
		t.Fatal(err)
	}
	cut := good[:len(good)-1000]
	if VerifyVHD(bytes.NewReader(cut), int64(len(cut))) == nil {
		t.Error("a cut VHD passed")
	}
	bad := bytes.Clone(good)
	bad[600]++ // the header
	if VerifyVHD(bytes.NewReader(bad), int64(len(bad))) == nil {
		t.Error("a damaged header passed")
	}
	fixed := append(make([]byte, 1<<20), good[len(good)-512:]...)
	binary.BigEndian.PutUint64(fixed[len(fixed)-512+48:], 1<<20)
	binary.BigEndian.PutUint32(fixed[len(fixed)-512+60:], 2)
	f := fixed[len(fixed)-512:]
	binary.BigEndian.PutUint32(f[64:], vhdChecksum(f, 64))
	if err := VerifyVHD(bytes.NewReader(fixed), int64(len(fixed))); err != nil {
		t.Errorf("fixed: %v", err)
	}
}

func exportVM(name string) ExportVM {
	return ExportVM{Name: name, Description: "web & db <1>", HVM: true, VCPUs: 2, VCPUsMax: 4, Memory: 1 << 30, Order: 1, StartDelay: 20, ShutdownDelay: 5,
		OSName: "Alpine Linux|/dev/xvda", OSMajor: "3", OSMinor: "22",
		NICs:          []ExportNIC{{Network: "n1", NetworkName: "Pool-wide", MAC: "aa:bb:cc:00:00:01"}, {Network: "n2", NetworkName: "lan"}},
		CDs:           []string{"3"},
		Disks:         []ExportDisk{{Href: name + "-0.vhd", Name: "root", Device: "0", Bootable: true, Size: 4096, Capacity: 32 << 20}, {Href: name + "-1.vhd.gz", Name: "data", Device: "2", Size: 100, Capacity: 1 << 30, Compressed: true}},
		HVMBootPolicy: "BIOS order", HVMBootParams: map[string]string{"order": "cd", "firmware": "bios"}, Platform: map[string]string{"acpi": "1", "nx": "true"},
		NVRAM: map[string]string{"EFI-variables": "abc"}, HVMShadowMultiplier: 1.5, PVArgs: "quiet", HardwarePlatformVersion: 2, HasVendorDevice: true}
}

// TestApplianceDescriptorRoundTrip reads the descriptor back with the import's parser.
func TestApplianceDescriptorRoundTrip(t *testing.T) {
	desc := ApplianceDescriptor("pkg", []ExportVM{exportVM("a")}, []string{"You agree & accept."})
	env, err := parseOVF(desc)
	if err != nil {
		t.Fatalf("%v\n%s", err, desc)
	}
	if env.System == nil || env.Collection != nil || env.Version != "1.0.0" {
		t.Fatalf("one VM is the envelope's VirtualSystem: %s", desc)
	}
	a, err := describeAppliance(env, "pkg", func(string) bool { return true }, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Warnings) != 0 || len(a.Eulas) != 1 || a.Eulas[0] != "You agree & accept." || a.Disks != 2 || a.Interfaces != 2 {
		t.Fatalf("%+v", a)
	}
	vm := a.VMs[0]
	if vm.Name != "a" || vm.Description != "web & db <1>" || vm.VCPUs != 2 || vm.VCPUsMax != 4 || vm.Memory != 1<<30 || vm.Order != 1 || vm.StartDelay != 20 || vm.ShutdownDelay != 5 {
		t.Fatalf("%+v", vm)
	}
	if vm.HVMBootPolicy != "BIOS order" || vm.HVMBootParams["order"] != "cd" || vm.HVMBootParams["firmware"] != "bios" || vm.Platform["acpi"] != "1" ||
		vm.NVRAM["EFI-variables"] != "abc" || !vm.HasVendorDevice || vm.Fields["PV_args"] != "quiet" || vm.Fields["HVM_shadow_multiplier"] != "1.5" ||
		vm.Fields["hardware_platform_version"] != "2" {
		t.Fatalf("settings %+v", vm)
	}
	if len(vm.MACs) != 2 || vm.MACs[0] != "aa:bb:cc:00:00:01" || vm.MACs[1] != "" || vm.DVD != "3" {
		t.Fatalf("interfaces %q, DVD %q", vm.MACs, vm.DVD)
	}
	d0, d1 := vm.Disks[0], vm.Disks[1]
	if d0.Href != "a-0.vhd" || d0.Device != "0" || !d0.Bootable || d0.Capacity != 32<<20 || d0.Name != "root" || d0.Compression != "" {
		t.Errorf("disk 0 %+v", d0)
	}
	if d1.Href != "a-1.vhd.gz" || d1.Device != "2" || d1.Bootable || d1.Capacity != 1<<30 || d1.Compression != "gzip" {
		t.Errorf("disk 1 %+v", d1)
	}
	if !strings.Contains(string(desc), `<ovf:OperatingSystemSection ovf:id="1"`) || !strings.Contains(string(desc), "<ovf:Description>Alpine Linux</ovf:Description>") ||
		!strings.Contains(string(desc), "<vssd:VirtualSystemType>hvm-3.0-unknown</vssd:VirtualSystemType>") {
		t.Errorf("OS or system type: %s", desc)
	}
}

func TestApplianceDescriptorCollection(t *testing.T) {
	b := exportVM("b")
	b.Disks, b.Order, b.HVM, b.Arch = nil, 2, false, "x64"
	desc := ApplianceDescriptor("app", []ExportVM{exportVM("a"), b}, []string{"EULA"})
	env, err := parseOVF(desc)
	if err != nil {
		t.Fatal(err)
	}
	a, err := describeAppliance(env, "app", func(string) bool { return true }, false)
	if err != nil {
		t.Fatal(err)
	}
	if env.Collection == nil || len(a.VMs) != 2 || a.VMs[0].Name != "a" || a.VMs[1].Name != "b" || a.VMs[1].Order != 2 || a.Disks != 2 || len(a.Eulas) != 1 {
		t.Fatalf("%+v", a)
	}
	if strings.Count(string(desc), "<ovf:Network ") != 2 || !strings.Contains(string(desc), "xen-3.0-x64") {
		t.Errorf("networks or type: %s", desc)
	}
}

// TestPackageRoundTrip writes packages as the export does (WritePackage: a folder with its disks compressed or not, an
// OVA compressed or not, with a manifest) and opens them with the import's package reader, manifest checked.
func TestPackageRoundTrip(t *testing.T) {
	for _, c := range []struct{ ova, compress bool }{{false, false}, {false, true}, {true, false}, {true, true}} {
		dir := t.TempDir()
		vhd := testVHD(4 << 20)
		if err := os.WriteFile(filepath.Join(dir, "d.vhd"), vhd, 0o644); err != nil {
			t.Fatal(err)
		}
		disk := ExportDisk{Href: "d.vhd", Name: "x", Device: "0", Size: int64(len(vhd)), Capacity: 4 << 20}
		if c.compress && !c.ova { // CompressFiles: each disk on its own
			n, err := GzipFile(filepath.Join(dir, "d.vhd"))
			if err != nil {
				t.Fatal(err)
			}
			disk.Href, disk.Size, disk.Compressed = "d.vhd.gz", n, true
		}
		vm := ExportVM{Name: "x", VCPUs: 1, VCPUsMax: 1, Memory: 256 << 20, HardwarePlatformVersion: -1, Disks: []ExportDisk{disk}}
		file, err := WritePackage(dir, "x", []ExportVM{vm}, []string{"eula"}, true, c.ova, c.compress, nil)
		if err != nil {
			t.Fatal(err)
		}
		want := map[[2]bool]string{{false, false}: "x.ovf", {false, true}: "x.ovf", {true, false}: "x.ova", {true, true}: "x.ova.gz"}[[2]bool{c.ova, c.compress}]
		if left, _ := os.ReadDir(dir); file != want || c.ova && len(left) != 1 || !c.ova && len(left) != 3 {
			t.Fatalf("%+v: %s, left %v", c, file, left)
		}
		p, err := OpenPackage(filepath.Join(dir, file))
		if err != nil {
			t.Fatal(err)
		}
		if err := p.VerifyManifest(); err != nil {
			t.Fatalf("%+v: %v", c, err)
		}
		s, err := p.OpenDisk(p.Appliance.VMs[0].Disks[0])
		if err != nil {
			t.Fatal(err)
		}
		s.Close()
		if s.Format != "vhd" || s.VirtualSize != 4<<20 || s.Length != int64(len(vhd)) || len(p.Appliance.Eulas) != 1 {
			t.Errorf("%+v: disk %+v", c, s)
		}
	}
}

// TestLongNameOVA: a package name longer than ustar's 100 bytes, and not ASCII, goes as a GNU long name (which the C#
// app's SharpZipLib reads; it ignores PAX names) and opens again.
func TestLongNameOVA(t *testing.T) {
	dir := t.TempDir()
	name := strings.Repeat("Máquina ", 14)
	file, err := WritePackage(dir, name, []ExportVM{{Name: "x", VCPUs: 1, VCPUsMax: 1, Memory: 256 << 20}}, nil, true, true, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(dir, file))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(raw, []byte("././@LongLink")) || bytes.Contains(raw, []byte("PaxHeaders")) {
		t.Error("not a GNU long name")
	}
	p, err := OpenPackage(filepath.Join(dir, file))
	if err != nil {
		t.Fatal(err)
	}
	if p.Appliance.Name != name {
		t.Errorf("the package is %q", p.Appliance.Name)
	}
	if err := p.VerifyManifest(); err != nil {
		t.Fatal(err)
	}
}

func TestGzipFile(t *testing.T) {
	f := filepath.Join(t.TempDir(), "d.vhd")
	if err := os.WriteFile(f, bytes.Repeat([]byte("x"), 1<<20), 0o644); err != nil {
		t.Fatal(err)
	}
	n, err := GzipFile(f)
	if err != nil || n <= 0 || n > 1<<16 {
		t.Fatalf("%d %v", n, err)
	}
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Error("the plain file is left")
	}
}

func TestReadEulas(t *testing.T) {
	dir := t.TempDir()
	write := func(name string, b []byte) string {
		f := filepath.Join(dir, name)
		if err := os.WriteFile(f, b, 0o644); err != nil {
			t.Fatal(err)
		}
		return f
	}
	utf16 := []byte{0xff, 0xfe, 'O', 0, 'k', 0, 0xe9, 0} // "Oké" as UTF-16LE with its mark
	got, err := ReadEulas([]string{write("a.txt", utf16), write("b.txt", []byte("\xef\xbb\xbfplain"))})
	if err != nil || len(got) != 2 || got[0] != "Oké" || got[1] != "plain" {
		t.Fatalf("%q %v", got, err)
	}
	if _, err := ReadEulas([]string{write("empty.txt", []byte(" \n"))}); err == nil {
		t.Error("an empty EULA was read")
	}
	if _, err := ReadEulas([]string{dir}); err == nil {
		t.Error("a folder was read as a EULA")
	}
	big := strings.Repeat("<", 1<<20-1) // each "<" is "&lt;" in the descriptor
	eulas := []string{big, big, big, big}
	if _, err := WritePackage(t.TempDir(), "x", []ExportVM{{Name: "x"}}, eulas, false, false, false, nil); err == nil {
		t.Error("a descriptor larger than the import reads was written")
	}
}

func TestApplianceNameError(t *testing.T) {
	for name, ok := range map[string]bool{"web": true, "my vApp 2": true, "": false, " ": false, "a/b": false, "a:b": false, "..": false} {
		if (ApplianceNameError(name) == "") != ok {
			t.Errorf("%q", name)
		}
	}
}
