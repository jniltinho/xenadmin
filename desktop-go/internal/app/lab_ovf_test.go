package app

import (
	"archive/tar"
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	xenapi "go/xenapi"

	"xcpng-center/internal/ops"
)

// labVMDK is a stream-optimized VMDK of size bytes whose first grain (64 KiB) is data.
func labVMDK(size int64) []byte {
	le := binary.LittleEndian
	h := make([]byte, 512)
	le.PutUint32(h, 0x564d444b)
	le.PutUint32(h[4:], 3)
	le.PutUint32(h[8:], 1<<16|1<<17|3)
	le.PutUint64(h[12:], uint64(size/512))
	le.PutUint64(h[20:], 128)
	le.PutUint64(h[64:], 1)
	le.PutUint16(h[77:], 1)
	var z bytes.Buffer
	w := zlib.NewWriter(&z)
	w.Write(bytes.Repeat([]byte("ovf!"), 128*512/4))
	w.Close()
	m := make([]byte, 12)
	le.PutUint32(m[8:], uint32(z.Len()))
	out := append(append(append(h, m...), z.Bytes()...), make([]byte, (512-(12+z.Len())%512)%512)...)
	return append(out, make([]byte, 512)...) // the end-of-stream marker
}

// labOVA writes an OVA of one VM (2 vCPUs, 512 MiB, an interface with a MAC) with a raw disk, a VMDK and, when qemu-img
// is installed, a dynamic VHD; and its SHA-256 manifest. broken cuts the VMDK short.
func labOVA(t *testing.T, name string, broken bool) (string, []int64) {
	dir := t.TempDir()
	type file struct {
		name string
		data []byte
	}
	files := []file{{"disk0.img", bytes.Repeat([]byte{0xa5}, 4<<20)}, {"disk1.vmdk", labVMDK(8 << 20)}}
	sizes := []int64{4 << 20, 8 << 20}
	if broken {
		files[1].data = files[1].data[:700]
	}
	if qemu, err := exec.LookPath("qemu-img"); err == nil {
		vhd := filepath.Join(dir, "d.vhd")
		if out, err := exec.Command(qemu, "create", "-f", "vpc", "-o", "subformat=dynamic,force_size=on", vhd, "6M").CombinedOutput(); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
		data, err := os.ReadFile(vhd)
		if err != nil {
			t.Fatal(err)
		}
		files, sizes = append(files, file{"disk2.vhd", data}), append(sizes, 6<<20)
	}
	var refs, disks, items strings.Builder
	for i, f := range files {
		fmt.Fprintf(&refs, `<File ovf:id="f%d" ovf:href="%s"/>`, i, f.name)
		fmt.Fprintf(&disks, `<Disk ovf:diskId="d%d" ovf:fileRef="f%d" ovf:capacity="%d"/>`, i, i, sizes[i])
		fmt.Fprintf(&items, `<Item><rasd:ResourceType>17</rasd:ResourceType><rasd:InstanceID>%d</rasd:InstanceID><rasd:ElementName>%s disk %d</rasd:ElementName>`+
			`<rasd:HostResource>ovf:/disk/d%d</rasd:HostResource><rasd:Connection>device=%d</rasd:Connection></Item>`, 10+i, name, i, i, i)
	}
	desc := fmt.Sprintf(`<?xml version="1.0"?>
<Envelope ovf:version="1.0" xmlns="http://schemas.dmtf.org/ovf/envelope/1" xmlns:ovf="http://schemas.dmtf.org/ovf/envelope/1"
  xmlns:rasd="http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_ResourceAllocationSettingData">
  <References>%s</References><DiskSection>%s</DiskSection>
  <VirtualSystem ovf:id="%s"><Name>%s</Name><VirtualHardwareSection>
    <Item><rasd:ResourceType>3</rasd:ResourceType><rasd:VirtualQuantity>2</rasd:VirtualQuantity></Item>
    <Item><rasd:ResourceType>4</rasd:ResourceType><rasd:AllocationUnits>byte * 2^20</rasd:AllocationUnits><rasd:VirtualQuantity>512</rasd:VirtualQuantity></Item>
    <Item><rasd:ResourceType>10</rasd:ResourceType><rasd:InstanceID>5</rasd:InstanceID><rasd:Address>02:16:3e:00:0a:01</rasd:Address></Item>
    %s
  </VirtualHardwareSection></VirtualSystem>
</Envelope>`, refs.String(), disks.String(), name, name, items.String())
	files = append([]file{{name + ".ovf", []byte(desc)}}, files...)
	var mf strings.Builder
	for _, f := range files {
		fmt.Fprintf(&mf, "SHA256(%s)= %x\n", f.name, sha256.Sum256(f.data))
	}
	files = append(files[:1], append([]file{{name + ".mf", []byte(mf.String())}}, files[1:]...)...)
	out := filepath.Join(dir, name+".ova")
	f, err := os.Create(out)
	if err != nil {
		t.Fatal(err)
	}
	tw := tar.NewWriter(f)
	for _, nf := range files {
		_ = tw.WriteHeader(&tar.Header{Name: nf.name, Mode: 0o644, Size: int64(len(nf.data)), Typeflag: tar.TypeReg})
		_, _ = tw.Write(nf.data)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return out, sizes
}

// TestLabImportAppliance imports a small OVA (raw, VMDK and VHD disks, an interface) to lab-nfs and checks the VM, then
// removes it and its vApp; and checks that a failed import leaves nothing behind (XCP_ALLOW_WRITE=1).
func TestLabImportAppliance(t *testing.T) {
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
	opts, err := ReadImportOptions(c, "")
	if err != nil || len(opts.Networks) == 0 {
		t.Fatalf("no network: %v", err)
	}
	name := fmt.Sprintf("xcpng-center-go-test-ova-%d", os.Getpid())
	leftovers := func() (vms []xenapi.VMRef, vdis []xenapi.VDIRef, apps []xenapi.VMApplianceRef) {
		vms, _ = xenapi.VM.GetByNameLabel(s, name)
		for i := range 3 {
			v, _ := xenapi.VDI.GetByNameLabel(s, fmt.Sprintf("%s disk %d", name, i))
			vdis = append(vdis, v...)
		}
		apps, _ = xenapi.VMAppliance.GetByNameLabel(s, name)
		return
	}
	clean := func() { // whatever an import left: the VMs, their disks and the vApp
		vms, vdis, apps := leftovers()
		for _, vm := range vms {
			_ = xenapi.VM.HardShutdown(s, vm)
			if err := xenapi.VM.Destroy(s, vm); err != nil {
				t.Errorf("VM not destroyed: %v", err)
			}
		}
		for _, v := range vdis {
			if err := xenapi.VDI.Destroy(s, v); err != nil {
				t.Errorf("disk not destroyed: %v", err)
			}
		}
		for _, a := range apps {
			if err := xenapi.VMAppliance.Destroy(s, a); err != nil {
				t.Errorf("vApp not destroyed: %v", err)
			}
		}
	}
	t.Cleanup(clean)

	// a broken VMDK: the import fails and removes what it made
	bad, _ := labOVA(t, name, true)
	if _, err := ops.ImportAppliance(c, labPolicy(t), bad, srUUID, opts.Networks[0].UUID, false, false); err == nil {
		t.Fatal("the broken appliance was imported")
	} else {
		t.Logf("broken: %v", err)
	}
	if vms, vdis, apps := leftovers(); len(vms)+len(vdis)+len(apps) != 0 {
		t.Fatalf("a failed import left %d VMs, %d disks, %d vApps", len(vms), len(vdis), len(apps))
	}

	file, sizes := labOVA(t, name, false)
	uuid, err := ops.ImportAppliance(c, labPolicy(t), file, srUUID, opts.Networks[0].UUID, true, false)
	if err != nil {
		t.Fatal(err)
	}
	vm, err := xenapi.VM.GetByUUID(s, uuid)
	if err != nil {
		t.Fatal(err)
	}
	rec, err := xenapi.VM.GetRecord(s, vm)
	if err != nil {
		t.Fatal(err)
	}
	if rec.NameLabel != name || rec.VCPUsAtStartup != 2 || rec.VCPUsMax != 2 || rec.MemoryStaticMax != 512<<20 || rec.HVMBootPolicy != "BIOS order" ||
		rec.OtherConfig["HideFromXenCenter"] != "" {
		t.Fatalf("the VM: %s, %d/%d vCPUs, %d memory, %q", rec.NameLabel, rec.VCPUsAtStartup, rec.VCPUsMax, rec.MemoryStaticMax, rec.HVMBootPolicy)
	}
	if app, _ := xenapi.VMAppliance.GetNameLabel(s, rec.Appliance); app != name {
		t.Errorf("the vApp is %q", app)
	}
	var got []int64
	for _, b := range rec.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil || vbd.Type != xenapi.VbdTypeDisk {
			continue
		}
		size, _ := xenapi.VDI.GetVirtualSize(s, vbd.VDI)
		got = append(got, int64(size))
		t.Logf("disk %s: %d bytes, bootable %v", vbd.Userdevice, size, vbd.Bootable)
	}
	slices.Sort(got)
	want := slices.Clone(sizes)
	slices.Sort(want)
	if len(got) != len(want) {
		t.Fatalf("disks %v, want %v", got, want)
	}
	for i := range got {
		if got[i] < want[i] || got[i] > want[i]+4<<20 { // the SR may round a size up
			t.Errorf("disks %v, want %v", got, want)
		}
	}
	if len(rec.VIFs) != 1 {
		t.Fatalf("%d interfaces", len(rec.VIFs))
	}
	if mac, _ := xenapi.VIF.GetMAC(s, rec.VIFs[0]); mac != "02:16:3e:00:0a:01" {
		t.Errorf("MAC %q", mac)
	}
}
