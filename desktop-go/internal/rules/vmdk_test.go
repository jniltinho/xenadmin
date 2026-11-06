package rules

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// hostedVMDK is a hosted sparse VMDK of capacity sectors with 8-sector grains and 4 entries a table, the grains given
// by index (a nil grain is a zeroed-grain entry), each table in its own sector before the grains.
func hostedVMDK(capacity int64, grains map[int64][]byte) []byte {
	const grain, gtes = 8, 4
	tables := (capacity + grain*gtes - 1) / (grain * gtes)
	gt0 := int64(2) // the directory in sector 1
	overhead := gt0 + tables
	f := make([]byte, overhead*sector)
	le := binary.LittleEndian
	le.PutUint32(f, vmdkMagic)
	le.PutUint32(f[4:], 1)
	le.PutUint32(f[8:], 1|vmdkFlagZeroGrain)
	le.PutUint64(f[12:], uint64(capacity))
	le.PutUint64(f[20:], grain)
	le.PutUint32(f[44:], gtes)
	le.PutUint64(f[56:], 1)
	le.PutUint64(f[64:], uint64(overhead))
	for t := int64(0); t < tables; t++ {
		le.PutUint32(f[sector+t*4:], uint32(gt0+t))
	}
	for g := int64(0); g < tables*gtes; g++ {
		data, ok := grains[g]
		if !ok {
			continue
		}
		at := (gt0+g/gtes)*sector + g%gtes*4
		if data == nil {
			le.PutUint32(f[at:], 1)
			continue
		}
		le.PutUint32(f[at:], uint32(len(f)/sector))
		f = append(f, data...)
		f = append(f, make([]byte, grain*sector-len(data))...)
	}
	return f
}

func readSparse(b []byte) ([]byte, error) {
	v, err := NewSparseVMDKReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, err
	}
	out, err := io.ReadAll(v)
	if err == nil && int64(len(out)) != v.Size() {
		err = fmt.Errorf("read %d bytes of %d", len(out), v.Size())
	}
	return out, err
}

func TestSparseVMDK(t *testing.T) {
	g := 8 * sector
	a, b, c := bytes.Repeat([]byte{1}, g), bytes.Repeat([]byte("xcp-"), g/4), bytes.Repeat([]byte{9}, g)
	// 13 grains, the last one 4 sectors; grain 5 a zeroed grain, the third table (grains 8 to 11) not allocated
	file := hostedVMDK(100, map[int64][]byte{0: a, 2: b, 5: nil, 12: c})
	le := binary.LittleEndian
	le.PutUint32(file[sector+2*4:], 0)
	got, err := readSparse(file)
	if err != nil {
		t.Fatal(err)
	}
	want := make([]byte, 100*sector)
	copy(want, a)
	copy(want[2*g:], b)
	copy(want[12*g:], c[:4*sector])
	if !bytes.Equal(got, want) {
		t.Fatal("not the raw disk")
	}
	// the last grain may end at the disk's end
	if _, err := readSparse(file[:len(file)-4*sector]); err != nil {
		t.Errorf("a last grain cut at the disk's end: %v", err)
	}
	gte0 := 2 * sector // grain 0's entry in the first table
	for name, mutate := range map[string]func(f []byte) []byte{
		"truncated":               func(f []byte) []byte { return f[:len(f)-5*sector] },
		"a grain past the end":    func(f []byte) []byte { le.PutUint32(f[gte0:], math.MaxUint32); return f },
		"a grain in the tables":   func(f []byte) []byte { le.PutUint32(f[gte0:], 2); return f },
		"a grain on the header":   func(f []byte) []byte { le.PutUint32(f[8:], 1); return f }, // grain 5's entry 1, no zero flag
		"a table past the end":    func(f []byte) []byte { le.PutUint32(f[sector:], math.MaxUint32); return f },
		"a table on the grains":   func(f []byte) []byte { le.PutUint32(f[sector:], uint32(len(f)/sector-1)); return f },
		"a huge capacity":         func(f []byte) []byte { le.PutUint64(f[12:], math.MaxInt64/sector); return f },
		"a capacity overflowing":  func(f []byte) []byte { le.PutUint64(f[12:], math.MaxUint64); return f },
		"a grain of 3 sectors":    func(f []byte) []byte { le.PutUint64(f[20:], 3); return f },
		"no table entries":        func(f []byte) []byte { le.PutUint32(f[44:], 0); return f },
		"metadata past the end":   func(f []byte) []byte { le.PutUint64(f[64:], 1<<40); return f },
		"a directory in the data": func(f []byte) []byte { le.PutUint64(f[56:], 6); return f },
		"stream-optimized":        func(f []byte) []byte { le.PutUint32(f[8:], vmdkFlagCompressed|vmdkFlagMarkers); return f },
	} {
		f := mutate(append([]byte{}, file...))
		if _, err := readSparse(f); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

func TestVMDKExtents(t *testing.T) {
	for desc, want := range map[string]string{
		"# Disk DescriptorFile\ncreateType=\"monolithicFlat\"\nRW 5860 FLAT \"d-flat.vmdk\" 0\n\x00\x00": "[{d-flat.vmdk false 0 3000320}]",
		"RW 8 VMFS \"d.raw\"\r\n":                                              "[{d.raw false 0 4096}]",
		"RDONLY 8 FLAT \"d.raw\" 2\n":                                          "[{d.raw false 1024 4096}]",
		"parentCID=ffffffff\nRW 8 FLAT \"a\" 0\nRW 4 FLAT \"b\" 0\n":           "[{a false 0 4096} {b false 0 2048}]",
		"RW 8 SPARSE \"d-s001.vmdk\"\nRW 8 SPARSE \"d-s002.vmdk\"\n":           "[{d-s001.vmdk true 0 4096} {d-s002.vmdk true 0 4096}]",
		"RW 8 SPARSE \"d-s001.vmdk\" 2\n":                                      "", // a sparse extent has no offset
		"RW 8 FLAT \"a\" 0\nRW 8 FLAT \"./a\" 0\n":                             "", // one file twice
		"parentCID=1234abcd\nRW 8 FLAT \"a\" 0\n":                              "", // a snapshot
		"NOACCESS 8 FLAT \"d.raw\" 0\n":                                        "",
		"RW 8 FLAT \"a\" 0\nRW 8 FLAT \"b\" -1\nRW 8 FLAT \"c\" 0\n":           "", // an extent that cannot be read
		"RW 8 FLAT \"a\" 0\nRW\n":                                              "",
		"RW 8 ZERO\n":                                                          "",
		"RW 8 VMFSSPARSE \"d.vmdk\"\n":                                         "",
		"RW 99999999999999999 FLAT \"d.raw\" 0\n":                              "",
		"RW 9000000000000000 FLAT \"a\" 0\nRW 9100000000000000 FLAT \"b\" 0\n": "", // the sum overflows
		"RW 8 FLAT \"d.raw\" 99999999999999999999999\n":                        "",
		"# Disk DescriptorFile\n":                                              "",
	} {
		got, err := vmdkExtents([]byte(desc))
		if s := fmt.Sprint(got); want == "" && err == nil || want != "" && (err != nil || s != want) {
			t.Errorf("%q: %s, %v", desc, s, err)
		}
	}
	var b strings.Builder // one extent more than vmdkMaxExtents
	for i := 0; i <= vmdkMaxExtents; i++ {
		fmt.Fprintf(&b, "RW 8 FLAT \"x%d\" 0\n", i)
	}
	if _, err := vmdkExtents([]byte(b.String())); err == nil {
		t.Error("too many extents accepted")
	}
}

// TestVMDKQemuKinds reads the hosted sparse and monolithic flat VMDKs qemu-img writes (when it is installed), in a
// folder package and in an OVA, back as the raw disk they came from.
func TestVMDKQemuKinds(t *testing.T) {
	qemu, err := exec.LookPath("qemu-img")
	if err != nil {
		t.Skip("no qemu-img")
	}
	dir := t.TempDir()
	raw := make([]byte, 3<<20+5*sector) // zeros between the patterns: unallocated grains
	for i := 0; i < len(raw); i += 300 << 10 {
		copy(raw[i:], strings.Repeat(fmt.Sprintf("sector %d ", i/sector), 50))
	}
	if err := os.WriteFile(filepath.Join(dir, "d.raw"), raw, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"s.vmdk"}, {"f.vmdk", "-o", "subformat=monolithicFlat"}} {
		cmd := append([]string{"convert", "-f", "raw", "-O", "vmdk"}, args[1:]...)
		if out, err := exec.Command(qemu, append(cmd, filepath.Join(dir, "d.raw"), filepath.Join(dir, args[0]))...).CombinedOutput(); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
	}
	read := func(t *testing.T, ovf, disk string) {
		t.Helper()
		p, err := OpenPackage(ovf)
		if err != nil {
			t.Fatal(err)
		}
		s, err := p.OpenDisk(ApplianceDisk{Href: disk})
		if err != nil {
			t.Fatal(err)
		}
		defer s.Close()
		got, err := io.ReadAll(s)
		if err != nil {
			t.Fatal(err)
		}
		// qemu rounds the capacity up to a sector; the data is the same
		if s.Format != "raw" || int64(len(got)) != s.Length || s.VirtualSize != s.Length || !bytes.Equal(got[:len(raw)], raw) ||
			bytes.ContainsFunc(got[len(raw):], func(r rune) bool { return r != 0 }) {
			t.Fatalf("%s: %d bytes of %d", disk, len(got), len(raw))
		}
	}
	file := func(n string) string {
		b, err := os.ReadFile(filepath.Join(dir, n))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	for _, disk := range []string{"s.vmdk", "f.vmdk"} {
		desc := strings.ReplaceAll(rawOVF, "disk.img", disk)
		if err := os.WriteFile(filepath.Join(dir, "tiny.ovf"), []byte(desc), 0o600); err != nil {
			t.Fatal(err)
		}
		read(t, filepath.Join(dir, "tiny.ovf"), disk)
		// an OVA: the disk is a section of the archive, after a member of an odd size
		members := [][2]string{{"tiny.ovf", desc}, {"odd", "x"}, {disk, file(disk)}}
		if disk == "f.vmdk" {
			members = append(members, [2]string{"f-flat.vmdk", file("f-flat.vmdk")})
		}
		ova(t, filepath.Join(dir, "t.ova"), false, members...)
		read(t, filepath.Join(dir, "t.ova"), disk)
	}
	// a hosted sparse disk in a compressed OVA cannot be read at random
	desc := strings.ReplaceAll(rawOVF, "disk.img", "s.vmdk")
	ova(t, filepath.Join(dir, "t.ova.gz"), true, [2]string{"tiny.ovf", desc}, [2]string{"s.vmdk", file("s.vmdk")})
	if p, err := OpenPackage(filepath.Join(dir, "t.ova.gz")); err != nil {
		t.Fatal(err)
	} else if _, err := p.OpenDisk(ApplianceDisk{Href: "s.vmdk"}); err == nil {
		t.Error("a hosted sparse disk in a compressed OVA")
	}
	// a flat extent shorter than its descriptor says, and one outside the package
	for _, d := range []string{regexp.MustCompile(`RW \d+`).ReplaceAllString(file("f.vmdk"), "RW 99999"), strings.Replace(file("f.vmdk"), `"f-flat.vmdk"`, `"../f-flat.vmdk"`, 1)} {
		if err := os.WriteFile(filepath.Join(dir, "g.vmdk"), []byte(d), 0o600); err != nil {
			t.Fatal(err)
		}
		p := &Package{file: filepath.Join(dir, "x.ovf"), dir: dir}
		if _, err := p.OpenDisk(ApplianceDisk{Href: "g.vmdk"}); err == nil {
			t.Errorf("descriptor %q accepted", d[:80])
		}
	}
	// a folder's manifest that does not list the flat extent: its check would not cover the disk's data
	p := &Package{file: filepath.Join(dir, "x.ovf"), dir: dir, manifest: []byte("SHA256(f.vmdk)= 00\n")}
	if _, err := p.OpenDisk(ApplianceDisk{Href: "f.vmdk"}); err == nil || !strings.Contains(err.Error(), "manifest") {
		t.Errorf("an extent the manifest does not list: %v", err)
	}
	p.manifest = append(p.manifest, "SHA256(F-flat.vmdk)= 00\n"...) // another file on a case-sensitive system
	if _, err := p.OpenDisk(ApplianceDisk{Href: "f.vmdk"}); err == nil {
		t.Error("an extent the manifest lists in another case")
	}
	p.manifest = append(p.manifest, "SHA256(f-flat.vmdk)= 00\n"...)
	if s, err := p.OpenDisk(ApplianceDisk{Href: "f.vmdk"}); err != nil {
		t.Errorf("an extent the manifest lists: %v", err)
	} else {
		s.Close()
	}
	// in an OVA too, and an extent the manifest check would skip
	fdesc := strings.ReplaceAll(rawOVF, "disk.img", "f.vmdk")
	for extent, mf := range map[string]string{"f-flat.vmdk": "SHA256(tiny.ovf)= 00\nSHA256(f.vmdk)= 00\n", "tiny.cert": "SHA256(tiny.cert)= 00\n"} {
		ova(t, filepath.Join(dir, "m.ova"), false, [2]string{"tiny.ovf", fdesc}, [2]string{"tiny.mf", mf},
			[2]string{"f.vmdk", strings.Replace(file("f.vmdk"), `"f-flat.vmdk"`, `"`+extent+`"`, 1)}, [2]string{extent, file("f-flat.vmdk")})
		if p, err := OpenPackage(filepath.Join(dir, "m.ova")); err != nil {
			t.Fatal(err)
		} else if _, err := p.OpenDisk(ApplianceDisk{Href: "f.vmdk"}); err == nil {
			t.Errorf("OVA extent %s accepted", extent)
		}
	}
	// a snapshot's delta (its embedded descriptor names a parent) holds only its changes
	delta := strings.Replace(file("s.vmdk"), "parentCID=ffffffff", "parentCID=1234abcd", 1)
	if delta == file("s.vmdk") {
		t.Fatal("no parentCID in qemu's descriptor")
	}
	if _, err := NewSparseVMDKReader(strings.NewReader(delta), int64(len(delta))); err == nil || !strings.Contains(err.Error(), "snapshot") {
		t.Errorf("a delta disk: %v", err)
	}
}
