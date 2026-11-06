package rules

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestVMDKSplit reads the split VMDKs qemu-img writes (twoGbMaxExtentSparse and twoGbMaxExtentFlat: a descriptor and
// extents of up to 2 GB) of a disk over 2 GB, mostly unallocated, back as the raw disk: the data across the extents'
// boundary and the zeros around it; then descriptors that do not match their extents.
func TestVMDKSplit(t *testing.T) {
	qemu, err := exec.LookPath("qemu-img")
	if err != nil {
		t.Skip("no qemu-img")
	}
	dir := t.TempDir()
	const size = 2<<30 + 64<<20
	f, err := os.Create(filepath.Join(dir, "d.raw"))
	if err != nil {
		t.Fatal(err)
	}
	for _, at := range []int64{0, 2<<30 - 3000, 2<<30 + 5<<20, size - 4096} {
		if _, err := f.WriteAt([]byte(strings.Repeat(fmt.Sprintf("at %d ", at), 1000)[:4096]), at); err != nil {
			t.Fatal(err)
		}
	}
	if err := f.Truncate(size); err != nil {
		t.Fatal(err)
	}
	f.Close()
	for n, kind := range map[string]string{"s": "Sparse", "f": "Flat"} {
		if out, err := exec.Command(qemu, "convert", "-f", "raw", "-O", "vmdk", "-o", "subformat=twoGbMaxExtent"+kind,
			filepath.Join(dir, "d.raw"), filepath.Join(dir, n+".vmdk")).CombinedOutput(); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
	}
	file := func(n string) string {
		b, err := os.ReadFile(filepath.Join(dir, n))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	if !strings.Contains(file("s.vmdk"), `RW 131072 SPARSE "s-s002.vmdk"`) || !strings.Contains(file("f.vmdk"), `RW 131072 FLAT "f-f002.vmdk" 0`) {
		t.Fatalf("qemu's descriptors: %q %q", file("s.vmdk"), file("f.vmdk"))
	}
	same := func(t *testing.T, p *Package, disk string) {
		t.Helper()
		s, err := p.OpenDisk(ApplianceDisk{Href: disk})
		if err != nil {
			t.Fatal(err)
		}
		defer s.Close()
		if s.Length != size || s.VirtualSize != size || s.Format != "raw" {
			t.Fatalf("%s: %d %d %s", disk, s.Length, s.VirtualSize, s.Format)
		}
		raw, err := os.Open(filepath.Join(dir, "d.raw"))
		if err != nil {
			t.Fatal(err)
		}
		defer raw.Close()
		a, b := make([]byte, 4<<20), make([]byte, 4<<20)
		for at := int64(0); ; at += int64(len(a)) { // the stream, never the disk, in memory
			n, err := io.ReadFull(s, a)
			m, _ := io.ReadFull(raw, b)
			if n != m || !bytes.Equal(a[:n], b[:m]) {
				t.Fatalf("%s differs from the raw disk in the 4 MB at %d (%d bytes, %d)", disk, at, n, m)
			}
			if err != nil {
				if at+int64(n) != size || !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Fatalf("%s: %v at %d", disk, err, at+int64(n))
				}
				return
			}
		}
	}
	for _, disk := range []string{"s.vmdk", "f.vmdk"} {
		desc := strings.ReplaceAll(rawOVF, "disk.img", disk)
		if err := os.WriteFile(filepath.Join(dir, "tiny.ovf"), []byte(desc), 0o600); err != nil {
			t.Fatal(err)
		}
		p, err := OpenPackage(filepath.Join(dir, "tiny.ovf"))
		if err != nil {
			t.Fatal(err)
		}
		same(t, p, disk)
	}
	// the sparse extents in an OVA (the flat ones would be 2 GB members), after a member of an odd size
	desc := strings.ReplaceAll(rawOVF, "disk.img", "s.vmdk")
	ova(t, filepath.Join(dir, "t.ova"), false, [2]string{"tiny.ovf", desc}, [2]string{"odd", "x"}, [2]string{"s.vmdk", file("s.vmdk")},
		[2]string{"s-s001.vmdk", file("s-s001.vmdk")}, [2]string{"s-s002.vmdk", file("s-s002.vmdk")})
	p, err := OpenPackage(filepath.Join(dir, "t.ova"))
	if err != nil {
		t.Fatal(err)
	}
	same(t, p, "s.vmdk")
	ova(t, filepath.Join(dir, "m.ova"), false, [2]string{"tiny.ovf", desc}, [2]string{"s.vmdk", file("s.vmdk")}, [2]string{"s-s001.vmdk", file("s-s001.vmdk")})
	if p, err := OpenPackage(filepath.Join(dir, "m.ova")); err != nil {
		t.Fatal(err)
	} else if _, err := p.OpenDisk(ApplianceDisk{Href: "s.vmdk"}); err == nil {
		t.Error("an OVA without its second extent")
	}
	// descriptors that do not match their extents, or name files they must not
	for name, d := range map[string]string{
		"sparse size":   strings.Replace(file("s.vmdk"), "RW 131072 SPARSE", "RW 131000 SPARSE", 1),
		"flat too long": strings.Replace(file("f.vmdk"), "RW 131072 FLAT", "RW 131073 FLAT", 1),
		"flat offset":   strings.Replace(file("f.vmdk"), `"f-f002.vmdk" 0`, `"f-f002.vmdk" 1`, 1),
		"missing":       strings.Replace(file("s.vmdk"), "s-s002.vmdk", "s-s003.vmdk", 1),
		"outside":       strings.Replace(file("f.vmdk"), `"f-f002.vmdk"`, `"../f-f002.vmdk"`, 1),
		"itself":        strings.Replace(file("f.vmdk"), `"f-f002.vmdk"`, `"g.vmdk"`, 1),
		"descriptor":    strings.Replace(file("f.vmdk"), `"f-f002.vmdk"`, `"tiny.ovf"`, 1),
		"twice":         strings.Replace(file("f.vmdk"), `"f-f002.vmdk"`, `"f-f001.vmdk"`, 1),
		"a snapshot":    strings.Replace(file("s.vmdk"), "parentCID=ffffffff", "parentCID=0badf00d", 1),
	} {
		if err := os.WriteFile(filepath.Join(dir, "g.vmdk"), []byte(d), 0o600); err != nil {
			t.Fatal(err)
		}
		p := &Package{file: filepath.Join(dir, "x.ovf"), dir: dir}
		if s, err := p.OpenDisk(ApplianceDisk{Href: "g.vmdk"}); err == nil {
			s.Close()
			t.Errorf("%s: accepted", name)
		}
	}
	// a manifest that lists the first extent only
	p = &Package{file: filepath.Join(dir, "x.ovf"), dir: dir, manifest: []byte("SHA256(s.vmdk)= 00\nSHA256(s-s001.vmdk)= 00\n")}
	if _, err := p.OpenDisk(ApplianceDisk{Href: "s.vmdk"}); err == nil || !strings.Contains(err.Error(), "s-s002.vmdk") {
		t.Errorf("an extent the manifest does not list: %v", err)
	}
}
