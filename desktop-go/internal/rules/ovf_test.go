package rules

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// a VMware-style descriptor: two controllers, two disks, a DVD drive, an interface and a EULA by message id
const vmwareOVF = `<?xml version="1.0" encoding="UTF-8"?>
<Envelope vmw:buildId="build-1" xmlns="http://schemas.dmtf.org/ovf/envelope/1" xmlns:ovf="http://schemas.dmtf.org/ovf/envelope/1"
  xmlns:rasd="http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_ResourceAllocationSettingData"
  xmlns:vssd="http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_VirtualSystemSettingData" xmlns:vmw="http://www.vmware.com/schema/ovf">
  <References>
    <File ovf:href="web-disk1.vmdk" ovf:id="file1" ovf:size="1024"/>
    <File ovf:href="web-disk2.vmdk" ovf:id="file2" ovf:size="1024"/>
  </References>
  <DiskSection><Info>disks</Info>
    <Disk ovf:capacity="2" ovf:capacityAllocationUnits="byte * 2^30" ovf:diskId="vmdisk1" ovf:fileRef="file1" ovf:format="http://www.vmware.com/interfaces/specifications/vmdk.html#streamOptimized"/>
    <Disk ovf:capacity="512" ovf:capacityAllocationUnits="byte * 2^20" ovf:diskId="vmdisk2" ovf:fileRef="file2"/>
  </DiskSection>
  <NetworkSection><Info>nets</Info><Network ovf:name="VM Network"><Description>The VM Network network</Description></Network></NetworkSection>
  <VirtualSystem ovf:id="web">
    <Info>A virtual machine</Info>
    <Name>web</Name>
    <EulaSection><Info>License</Info><License ovf:msgid="eula">fallback</License></EulaSection>
    <VirtualHardwareSection>
      <Info>Virtual hardware requirements</Info>
      <System><vssd:ElementName>Virtual Hardware Family</vssd:ElementName><vssd:VirtualSystemIdentifier>web</vssd:VirtualSystemIdentifier><vssd:VirtualSystemType>vmx-07</vssd:VirtualSystemType></System>
      <Item><rasd:AllocationUnits>hertz * 10^6</rasd:AllocationUnits><rasd:ElementName>2 virtual CPU(s)</rasd:ElementName><rasd:InstanceID>1</rasd:InstanceID><rasd:ResourceType>3</rasd:ResourceType><rasd:VirtualQuantity>2</rasd:VirtualQuantity></Item>
      <Item><rasd:AllocationUnits>byte * 2^20</rasd:AllocationUnits><rasd:InstanceID>2</rasd:InstanceID><rasd:ResourceType>4</rasd:ResourceType><rasd:VirtualQuantity>1024</rasd:VirtualQuantity></Item>
      <Item><rasd:Address>0</rasd:Address><rasd:ElementName>SCSI controller 0</rasd:ElementName><rasd:InstanceID>3</rasd:InstanceID><rasd:ResourceSubType>lsilogic</rasd:ResourceSubType><rasd:ResourceType>6</rasd:ResourceType></Item>
      <Item><rasd:Address>1</rasd:Address><rasd:ElementName>IDE 1</rasd:ElementName><rasd:InstanceID>4</rasd:InstanceID><rasd:ResourceType>5</rasd:ResourceType></Item>
      <Item><rasd:Address>0</rasd:Address><rasd:ElementName>IDE 0</rasd:ElementName><rasd:InstanceID>5</rasd:InstanceID><rasd:ResourceType>5</rasd:ResourceType></Item>
      <Item ovf:required="false"><rasd:AddressOnParent>0</rasd:AddressOnParent><rasd:AutomaticAllocation>false</rasd:AutomaticAllocation><rasd:ElementName>CD/DVD drive 1</rasd:ElementName><rasd:InstanceID>7</rasd:InstanceID><rasd:Parent>4</rasd:Parent><rasd:ResourceType>15</rasd:ResourceType></Item>
      <Item><rasd:AddressOnParent>1</rasd:AddressOnParent><rasd:ElementName>Hard disk 2</rasd:ElementName><rasd:HostResource>ovf:/disk/vmdisk2</rasd:HostResource><rasd:InstanceID>9</rasd:InstanceID><rasd:Parent>3</rasd:Parent><rasd:ResourceType>17</rasd:ResourceType></Item>
      <Item><rasd:AddressOnParent>0</rasd:AddressOnParent><rasd:ElementName>Hard disk 1</rasd:ElementName><rasd:HostResource>ovf:/disk/vmdisk1</rasd:HostResource><rasd:InstanceID>8</rasd:InstanceID><rasd:Parent>3</rasd:Parent><rasd:ResourceType>17</rasd:ResourceType></Item>
      <Item><rasd:AddressOnParent>7</rasd:AddressOnParent><rasd:AutomaticAllocation>true</rasd:AutomaticAllocation><rasd:Connection>VM Network</rasd:Connection><rasd:ElementName>Network adapter 1</rasd:ElementName><rasd:InstanceID>10</rasd:InstanceID><rasd:ResourceSubType>E1000</rasd:ResourceSubType><rasd:ResourceType>10</rasd:ResourceType></Item>
      <Item><rasd:Address>0050569a1b2c</rasd:Address><rasd:InstanceID>11</rasd:InstanceID><rasd:ResourceType>10</rasd:ResourceType></Item>
    </VirtualHardwareSection>
  </VirtualSystem>
  <Strings xml:lang="en-US"><Msg ovf:msgid="eula">You agree.</Msg></Strings>
</Envelope>`

// an XCP-ng Center export: a collection, the Xen settings and the disk's Connection
const xenOVF = `<?xml version="1.0"?>
<ovf:Envelope ovf:version="1.0.0" xmlns:ovf="http://schemas.dmtf.org/ovf/envelope/1" xmlns:xenovf="http://schemas.citrix.com/ovf/envelope/1"
  xmlns:rasd="http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_ResourceAllocationSettingData" xmlns:vssd="http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_VirtualSystemSettingData">
  <ovf:References><ovf:File ovf:id="f1" ovf:href="a.vhd"/></ovf:References>
  <ovf:DiskSection><ovf:Disk ovf:diskId="d1" ovf:fileRef="f1" ovf:capacity="1073741824"/></ovf:DiskSection>
  <ovf:VirtualSystemCollection ovf:id="app">
    <ovf:StartupSection><ovf:Item ovf:id="vm1" ovf:order="2" ovf:startDelay="30" ovf:stopDelay="10"/></ovf:StartupSection>
    <ovf:VirtualSystem ovf:id="vm1">
      <ovf:VirtualHardwareSection>
        <ovf:System><vssd:Description>my vm</vssd:Description><vssd:VirtualSystemType>hvm-3.0-unknown</vssd:VirtualSystemType></ovf:System>
        <ovf:Item><rasd:ResourceType>4</rasd:ResourceType><rasd:AllocationUnits>byte * 2^30</rasd:AllocationUnits><rasd:VirtualQuantity>2</rasd:VirtualQuantity></ovf:Item>
        <ovf:Item><rasd:ResourceType>3</rasd:ResourceType><rasd:VirtualQuantity>1</rasd:VirtualQuantity><rasd:Limit>4</rasd:Limit></ovf:Item>
        <ovf:Item><rasd:ResourceType>19</rasd:ResourceType><rasd:InstanceID>d1</rasd:InstanceID><rasd:ElementName>sys</rasd:ElementName><rasd:Connection>device=1,bootable=false,mode=r</rasd:Connection></ovf:Item>
        <xenovf:VirtualSystemOtherConfigurationData xenovf:Name="HVM_boot_policy"><xenovf:Value>BIOS order</xenovf:Value></xenovf:VirtualSystemOtherConfigurationData>
        <xenovf:VirtualSystemOtherConfigurationData xenovf:Name="HVM_boot_params"><xenovf:Value>cd</xenovf:Value></xenovf:VirtualSystemOtherConfigurationData>
        <xenovf:VirtualSystemOtherConfigurationData xenovf:Name="platform"><xenovf:Value>acpi=1;apic=true</xenovf:Value></xenovf:VirtualSystemOtherConfigurationData>
        <xenovf:VirtualSystemOtherConfigurationData xenovf:Name="PV_args"><xenovf:Value>quiet</xenovf:Value></xenovf:VirtualSystemOtherConfigurationData>
        <xenovf:VirtualSystemOtherConfigurationData xenovf:Name="VM_has_vendor_device"><xenovf:Value>True</xenovf:Value></xenovf:VirtualSystemOtherConfigurationData>
      </ovf:VirtualHardwareSection>
    </ovf:VirtualSystem>
  </ovf:VirtualSystemCollection>
</ovf:Envelope>`

func describe(t *testing.T, desc string, has func(string) bool) *Appliance {
	t.Helper()
	env, err := parseOVF([]byte(desc))
	if err != nil {
		t.Fatal(err)
	}
	if has == nil {
		has = func(string) bool { return true }
	}
	a, err := describeAppliance(env, "pkg", has, false)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func TestIsAppliance(t *testing.T) {
	for name, want := range map[string]bool{"a.ovf": true, "A.OVA": true, "a.ova.gz": true, "a.ova.bz2": true, "a.xva": false, "a.ova.zip": false, "": false} {
		if IsAppliance(name) != want {
			t.Errorf("%q: %v", name, !want)
		}
	}
}

func TestApplianceVMware(t *testing.T) {
	a := describe(t, vmwareOVF, nil)
	if len(a.VMs) != 1 || a.Disks != 2 || a.Interfaces != 2 {
		t.Fatalf("%+v", a)
	}
	vm := a.VMs[0]
	if vm.Name != "web" || vm.VCPUs != 2 || vm.VCPUsMax != 2 || vm.Memory != 1<<30 || vm.HVMBootPolicy != "BIOS order" || vm.Platform["stdvga"] != "0" {
		t.Fatalf("%+v", vm)
	}
	if vm.Description != "Imported from an OVF/OVA package." {
		t.Fatalf("description %q", vm.Description)
	}
	// SCSI (type 6) after the two IDE controllers: the DVD drive at 2, the disks at 0 and 1 by AddressOnParent
	if vm.DVD != "2" {
		t.Errorf("DVD at %q", vm.DVD)
	}
	d1, d2 := vm.Disks[0], vm.Disks[1]
	if d1.Href != "web-disk2.vmdk" || d1.Device != "1" || d1.Bootable || d1.Capacity != 512<<20 {
		t.Errorf("disk 2: %+v", d1)
	}
	if d2.Href != "web-disk1.vmdk" || d2.Device != "0" || !d2.Bootable || d2.Capacity != 2<<30 || d2.Name != "Hard disk 1" {
		t.Errorf("disk 1: %+v", d2)
	}
	if vm.MACs[0] != "" || vm.MACs[1] != "00:50:56:9a:1b:2c" {
		t.Errorf("MACs %q", vm.MACs)
	}
	if len(a.Eulas) != 1 || a.Eulas[0] != "You agree." {
		t.Errorf("EULAs %q", a.Eulas)
	}
	if a.RequiredSpace != 2<<30+512<<20 {
		t.Errorf("space %d", a.RequiredSpace)
	}
	if !strings.Contains(strings.Join(a.Warnings, "\n"), "version is not set") {
		t.Errorf("warnings %q", a.Warnings)
	}
}

func TestApplianceXen(t *testing.T) {
	a := describe(t, xenOVF, nil)
	vm := a.VMs[0]
	if vm.Name != "vm1" || vm.Description != "my vm" || vm.Memory != 2<<30 || vm.VCPUs != 1 || vm.VCPUsMax != 4 {
		t.Fatalf("%+v", vm)
	}
	if vm.HVMBootPolicy != "BIOS order" || vm.HVMBootParams["order"] != "cd" || vm.Platform["acpi"] != "1" || vm.Platform["nx"] != "true" ||
		vm.Fields["PV_args"] != "quiet" || !vm.HasVendorDevice {
		t.Fatalf("%+v", vm)
	}
	if vm.Order != 2 || vm.StartDelay != 30 || vm.ShutdownDelay != 10 {
		t.Errorf("startup %+v", vm)
	}
	if d := vm.Disks[0]; d.Device != "1" || d.Bootable || !d.ReadOnly || d.Capacity != 1<<30 || d.Name != "sys" {
		t.Errorf("disk %+v", d)
	}
	if len(a.Warnings) != 0 || len(a.Eulas) != 0 || vm.DVD != "" {
		t.Errorf("%+v", a)
	}
}

func TestApplianceHostile(t *testing.T) {
	for name, desc := range map[string]string{
		"traversal": strings.Replace(xenOVF, `ovf:href="a.vhd"`, `ovf:href="../../etc/passwd"`, 1),
		"absolute":  strings.Replace(xenOVF, `ovf:href="a.vhd"`, `ovf:href="/etc/passwd"`, 1),
		"url":       strings.Replace(xenOVF, `ovf:href="a.vhd"`, `ovf:href="http://example.com/a.vhd"`, 1),
		"doctype":   strings.Replace(xenOVF, `<?xml version="1.0"?>`, `<?xml version="1.0"?><!DOCTYPE x [<!ENTITY a "aaaa">]>`, 1),
		"entity":    strings.Replace(xenOVF, `my vm`, `&a;`, 1),
		"encrypted": strings.Replace(xenOVF, `<ovf:References>`, `<ovf:SecuritySection/><ovf:References>`, 1),
		"not ovf":   `<html/>`,
		"empty":     `<Envelope/>`,
	} {
		env, err := parseOVF([]byte(desc))
		if err == nil {
			_, err = describeAppliance(env, "pkg", func(string) bool { return true }, false)
		}
		if err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
	env, _ := parseOVF([]byte(xenOVF))
	if _, err := describeAppliance(env, "pkg", func(string) bool { return false }, false); err == nil || !strings.Contains(err.Error(), "missing") {
		t.Errorf("a missing file: %v", err)
	}
	if _, err := parseOVF(bytes.Repeat([]byte(" "), ovfMaxDescriptor+1)); err == nil {
		t.Error("a huge descriptor")
	}
	for href, want := range map[string]bool{"a.vmdk": true, "dir/a.vmdk": true, "../a": false, "a/../../b": false, "/a": false, `a\b`: false, "c:a": false, "": false, "./a": false} {
		if SafeHref(href) != want {
			t.Errorf("SafeHref(%q) = %v", href, !want)
		}
	}
}

// streamVMDK writes a stream-optimized VMDK of size bytes with the given grains (128 sectors), in order, as QEMU does.
func streamVMDK(t *testing.T, size int64, grains map[int64][]byte, order []int64) []byte {
	t.Helper()
	var b bytes.Buffer
	h := make([]byte, sector)
	le := binary.LittleEndian
	le.PutUint32(h, vmdkMagic)
	le.PutUint32(h[4:], 3)
	le.PutUint32(h[8:], vmdkFlagCompressed|vmdkFlagMarkers|3)
	le.PutUint64(h[12:], uint64(size/sector))
	le.PutUint64(h[20:], 128)
	le.PutUint64(h[64:], 2) // a sector of embedded descriptor, then the grains
	le.PutUint16(h[77:], 1)
	b.Write(h)
	b.Write(make([]byte, sector))
	pad := func() { b.Write(make([]byte, (sector-b.Len()%sector)%sector)) }
	meta := func(sectors uint64, kind uint32) {
		m := make([]byte, sector)
		le.PutUint64(m, sectors)
		le.PutUint32(m[12:], kind)
		b.Write(m)
		b.Write(make([]byte, sectors*sector))
	}
	for _, g := range order {
		var z bytes.Buffer
		w := zlib.NewWriter(&z)
		w.Write(grains[g])
		w.Close()
		m := make([]byte, 12)
		le.PutUint64(m, uint64(g*128))
		le.PutUint32(m[8:], uint32(z.Len()))
		b.Write(m)
		b.Write(z.Bytes())
		pad()
	}
	meta(1, 1) // a grain table, the directory and the footer, skipped
	meta(1, 2)
	meta(1, 3)
	meta(0, 0)
	return b.Bytes()
}

func TestVMDKReader(t *testing.T) {
	grain := 128 * sector
	g1, g3 := bytes.Repeat([]byte{1}, grain), bytes.Repeat([]byte("xcp-"), grain/4)
	size := int64(4*grain + 3*sector) // the last grain past the disk's end
	file := streamVMDK(t, size, map[int64][]byte{1: g1, 3: g3, 4: bytes.Repeat([]byte{9}, grain)}, []int64{1, 3, 4})
	v, err := NewVMDKReader(bytes.NewReader(file))
	if err != nil {
		t.Fatal(err)
	}
	got, err := io.ReadAll(v)
	if err != nil {
		t.Fatal(err)
	}
	want := make([]byte, size)
	copy(want[grain:], g1)
	copy(want[3*grain:], g3)
	copy(want[4*grain:], bytes.Repeat([]byte{9}, 3*sector))
	if v.Size() != size || !bytes.Equal(got, want) {
		t.Fatalf("read %d bytes, want %d; equal %v", len(got), size, bytes.Equal(got, want))
	}
	// grains out of order, a truncated file, a VMDK that is not stream-optimized
	if _, err := io.ReadAll(must(NewVMDKReader(bytes.NewReader(streamVMDK(t, size, map[int64][]byte{1: g1, 3: g3}, []int64{3, 1}))))); err == nil {
		t.Error("out of order")
	}
	if _, err := io.ReadAll(must(NewVMDKReader(bytes.NewReader(file[:len(file)-600])))); err == nil {
		t.Error("truncated")
	}
	// a grain inside the disk that is short, and one with more than a grain of data
	for _, g := range [][]byte{g1[:sector], append(g1, 1)} {
		if _, err := io.ReadAll(must(NewVMDKReader(bytes.NewReader(streamVMDK(t, size, map[int64][]byte{1: g}, []int64{1}))))); err == nil {
			t.Errorf("a grain of %d bytes", len(g))
		}
	}
	flat := append([]byte{}, file...)
	binary.LittleEndian.PutUint32(flat[8:], 3)
	if _, err := NewVMDKReader(bytes.NewReader(flat)); err == nil {
		t.Error("not stream-optimized")
	}
}

func must(v *vmdkReader, err error) io.Reader {
	if err != nil {
		return strings.NewReader("")
	}
	return v
}

// TestVMDKQemu reads a VMDK written by qemu-img (when it is installed) back as the raw disk it came from.
func TestVMDKQemu(t *testing.T) {
	qemu, err := exec.LookPath("qemu-img")
	if err != nil {
		t.Skip("no qemu-img")
	}
	dir := t.TempDir()
	raw := make([]byte, 3<<20+5*sector)
	for i := 0; i < len(raw); i += 4096 * 7 {
		copy(raw[i:], fmt.Sprintf("sector %d", i/sector))
	}
	if err := os.WriteFile(filepath.Join(dir, "d.raw"), raw, 0o600); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command(qemu, "convert", "-f", "raw", "-O", "vmdk", "-o", "subformat=streamOptimized",
		filepath.Join(dir, "d.raw"), filepath.Join(dir, "d.vmdk")).CombinedOutput(); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	f, err := os.Open(filepath.Join(dir, "d.vmdk"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	v, err := NewVMDKReader(f)
	if err != nil {
		t.Fatal(err)
	}
	got, err := io.ReadAll(v)
	if err != nil {
		t.Fatal(err)
	}
	// qemu rounds the capacity up to a sector; the data is the same
	if int64(len(got)) != v.Size() || !bytes.Equal(got[:len(raw)], raw) {
		t.Fatalf("read %d bytes of %d", len(got), len(raw))
	}
}

// ova writes a tar of the files (in order), optionally gzipped.
func ova(t *testing.T, file string, gz bool, files ...[2]string) {
	t.Helper()
	f, err := os.Create(file)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var w io.Writer = f
	if gz {
		z := gzip.NewWriter(f)
		defer z.Close()
		w = z
	}
	tw := tar.NewWriter(w)
	for _, nf := range files {
		_ = tw.WriteHeader(&tar.Header{Name: nf[0], Mode: 0o644, Size: int64(len(nf[1])), Typeflag: tar.TypeReg})
		_, _ = tw.Write([]byte(nf[1]))
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
}

const rawOVF = `<?xml version="1.0"?>
<Envelope ovf:version="1.0" xmlns="http://schemas.dmtf.org/ovf/envelope/1" xmlns:ovf="http://schemas.dmtf.org/ovf/envelope/1"
  xmlns:rasd="http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_ResourceAllocationSettingData">
  <References><File ovf:id="f1" ovf:href="disk.img"/></References>
  <DiskSection><Disk ovf:diskId="d1" ovf:fileRef="f1" ovf:capacity="4096"/></DiskSection>
  <VirtualSystem ovf:id="tiny"><VirtualHardwareSection>
    <Item><rasd:ResourceType>17</rasd:ResourceType><rasd:InstanceID>9</rasd:InstanceID><rasd:HostResource>ovf:/disk/d1</rasd:HostResource></Item>
  </VirtualHardwareSection></VirtualSystem>
</Envelope>`

func TestPackage(t *testing.T) {
	dir := t.TempDir()
	disk := strings.Repeat("d", 4096)
	sum := func(s string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(s))) }
	mf := fmt.Sprintf("SHA256(tiny.ovf)= %s\nSHA256(disk.img)= %s\n", sum(rawOVF), sum(disk))
	for _, gz := range []bool{false, true} {
		file := filepath.Join(dir, fmt.Sprintf("tiny%v.ova", gz))
		if gz {
			file += ".gz"
		}
		ova(t, file, gz, [2]string{"tiny.ovf", rawOVF}, [2]string{"tiny.mf", mf}, [2]string{"./disk.img", disk})
		p, err := OpenPackage(file)
		if err != nil {
			t.Fatal(err)
		}
		if p.Appliance.Name != "tiny" || !p.Appliance.HasManifest || p.Appliance.VMs[0].Memory != defaultImportMemory {
			t.Fatalf("%+v", p.Appliance)
		}
		if err := p.VerifyManifest(); err != nil {
			t.Fatal(err)
		}
		s, err := p.OpenDisk(p.Appliance.VMs[0].Disks[0])
		if err != nil {
			t.Fatal(err)
		}
		got, _ := io.ReadAll(s)
		s.Close()
		if s.Format != "raw" || s.Length != 4096 || s.VirtualSize != 4096 || string(got) != disk {
			t.Fatalf("%+v %d", s, len(got))
		}
	}
	// a changed disk, a file the manifest does not list
	bad := filepath.Join(dir, "bad.ova")
	ova(t, bad, false, [2]string{"tiny.ovf", rawOVF}, [2]string{"tiny.mf", mf}, [2]string{"disk.img", disk[1:] + "x"})
	if p, err := OpenPackage(bad); err != nil || p.VerifyManifest() == nil {
		t.Errorf("a changed disk: %v", err)
	}
	ova(t, bad, false, [2]string{"tiny.ovf", rawOVF}, [2]string{"tiny.mf", mf}, [2]string{"disk.img", disk}, [2]string{"extra", "x"})
	if p, err := OpenPackage(bad); err != nil || p.VerifyManifest() == nil {
		t.Errorf("an extra file: %v", err)
	}
	ova(t, bad, false, [2]string{"tiny.ovf", rawOVF})
	if _, err := OpenPackage(bad); err == nil {
		t.Error("a missing disk")
	}
	// a folder package: the descriptor beside its disk, a manifest naming a file outside it
	if err := os.WriteFile(filepath.Join(dir, "tiny.ovf"), []byte(rawOVF), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "disk.img"), []byte(disk), 0o600); err != nil {
		t.Fatal(err)
	}
	for manifest, want := range map[string]string{mf: "", mf + "SHA256(../x)= 00\n": "outside", fmt.Sprintf("SHA256(tiny.ovf)= %s\n", sum(rawOVF)): "not listed"} {
		if err := os.WriteFile(filepath.Join(dir, "tiny.mf"), []byte(manifest), 0o600); err != nil {
			t.Fatal(err)
		}
		p, err := OpenPackage(filepath.Join(dir, "tiny.ovf"))
		if err != nil {
			t.Fatal(err)
		}
		if err := p.VerifyManifest(); want == "" && err != nil || want != "" && (err == nil || !strings.Contains(err.Error(), want)) {
			t.Errorf("manifest %q: %v", manifest, err)
		}
		if s, err := p.OpenDisk(p.Appliance.VMs[0].Disks[0]); err != nil || s.Length != 4096 {
			t.Errorf("folder disk: %v", err)
		} else {
			s.Close()
		}
	}
	// a disk that is a link out of the folder
	outside := filepath.Join(t.TempDir(), "secret")
	if err := os.WriteFile(outside, []byte(disk), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dir, "disk.img")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "disk.img")); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenPackage(filepath.Join(dir, "tiny.ovf")); err == nil {
		t.Error("a disk linked out of the package")
	}
}

// TestApplianceDevices: a disk id that is another's prefix, and a disk under a Microsoft synthetic drive (its Address).
func TestApplianceDevices(t *testing.T) {
	desc := strings.NewReplacer(`<Disk ovf:diskId="d1" ovf:fileRef="f1" ovf:capacity="4096"/>`,
		`<Disk ovf:diskId="d1" ovf:fileRef="f1" ovf:capacity="4096"/><Disk ovf:diskId="d10" ovf:fileRef="f2" ovf:capacity="8192"/>`,
		`<File ovf:id="f1" ovf:href="disk.img"/>`, `<File ovf:id="f1" ovf:href="disk.img"/><File ovf:id="f2" ovf:href="ten.img"/>`,
		`<Item><rasd:ResourceType>17</rasd:ResourceType><rasd:InstanceID>9</rasd:InstanceID><rasd:HostResource>ovf:/disk/d1</rasd:HostResource></Item>`,
		`<Item><rasd:ResourceType>5</rasd:ResourceType><rasd:InstanceID>ide0</rasd:InstanceID><rasd:Address>0</rasd:Address></Item>
		<Item><rasd:ResourceType>22</rasd:ResourceType><rasd:ResourceSubType>Microsoft Synthetic Disk Drive</rasd:ResourceSubType><rasd:InstanceID>syn</rasd:InstanceID><rasd:Parent>ide0</rasd:Parent><rasd:Address>0</rasd:Address></Item>
		<Item><rasd:ResourceType>21</rasd:ResourceType><rasd:Caption>Hard Disk Image</rasd:Caption><rasd:InstanceID>9</rasd:InstanceID><rasd:Parent>syn</rasd:Parent><rasd:HostResource>ovf:/disk/d10</rasd:HostResource></Item>`).Replace(rawOVF)
	vm := describe(t, desc, nil).VMs[0]
	if len(vm.Disks) != 1 || vm.Disks[0].Href != "ten.img" || vm.Disks[0].Capacity != 8192 || vm.Disks[0].Device != "0" || !vm.Disks[0].Bootable {
		t.Fatalf("%+v", vm.Disks)
	}
}

func TestOpenDiskFormats(t *testing.T) {
	dir := t.TempDir()
	vhd := make([]byte, 3*sector) // a dynamic VHD's first sector is a copy of its footer
	copy(vhd, "conectix")
	binary.BigEndian.PutUint64(vhd[48:], 8<<20)
	fixed := append(make([]byte, 4096), vhd[:sector]...)
	grain := 128 * sector
	vmdk := streamVMDK(t, int64(2*grain), map[int64][]byte{1: bytes.Repeat([]byte{7}, grain)}, []int64{1})
	var gz bytes.Buffer
	z := gzip.NewWriter(&gz)
	z.Write(fixed)
	z.Close()
	for name, c := range map[string]struct {
		data, comp    string
		format        string
		length, vsize int64
	}{
		"dyn.vhd":   {string(vhd), "", "vhd", 3 * sector, 8 << 20},
		"fixed.vhd": {string(fixed), "", "vhd", 4096 + sector, 4096},
		"d.vmdk":    {string(vmdk), "", "raw", int64(2 * grain), int64(2 * grain)},
		"f.vhd.gz":  {gz.String(), "gzip", "vhd", 4096 + sector, 4096},
		"bad.vmdk":  {"# Disk DescriptorFile", "", "", 0, 0},
		"x.qcow2":   {"QFI", "", "", 0, 0},
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(c.data), 0o600); err != nil {
			t.Fatal(err)
		}
		p := &Package{file: filepath.Join(dir, "x.ovf"), dir: dir}
		s, err := p.OpenDisk(ApplianceDisk{Href: name, Compression: c.comp})
		if c.format == "" {
			if err == nil {
				t.Errorf("%s: accepted", name)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		n, _ := io.Copy(io.Discard, s)
		s.Close()
		if s.Format != c.format || s.Length != c.length || s.VirtualSize != c.vsize || n != c.length {
			t.Errorf("%s: %+v, read %d", name, s, n)
		}
	}
}

// TestApplianceLoops: synthetic drives that are their own parents end; a disk named like a manifest is refused.
func TestApplianceLoops(t *testing.T) {
	loop := strings.Repeat(`<Item><rasd:ResourceType>22</rasd:ResourceType><rasd:ResourceSubType>synthetic</rasd:ResourceSubType><rasd:InstanceID>s</rasd:InstanceID><rasd:Parent>s</rasd:Parent></Item>`, 20)
	desc := strings.Replace(rawOVF, `<Item><rasd:ResourceType>17`, `<Item><rasd:ResourceType>5</rasd:ResourceType><rasd:InstanceID>s</rasd:InstanceID></Item>`+loop+`<Item><rasd:ResourceType>17`, 1)
	if a := describe(t, desc, nil); len(a.VMs[0].Disks) != 1 {
		t.Fatalf("%+v", a.VMs[0])
	}
	env, _ := parseOVF([]byte(strings.Replace(rawOVF, `disk.img`, `disk.cert`, 1)))
	if _, err := describeAppliance(env, "pkg", func(string) bool { return true }, false); err == nil {
		t.Error("a .cert disk")
	}
}
