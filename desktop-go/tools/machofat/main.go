// machofat joins thin 64-bit Mach-O executables (one per CPU) into a universal ("fat") one, as Apple's
// `lipo -create` does, so the macOS builds made on Linux are universal (Apple Silicon and Intel). Usage:
//
//	go run ./tools/machofat -o out thin-amd64 thin-arm64
package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
)

const (
	fatMagic   = 0xcafebabe
	thinMagic  = 0xfeedfacf // MH_MAGIC_64, little-endian as written by amd64 and arm64 builds
	alignShift = 14         // 16 KiB: the page size of arm64 macOS, what lipo uses for arm64 slices
)

// slice is one thin Mach-O of the universal file.
type slice struct {
	path         string
	cpu, subtype uint32
	size, offset uint32
}

// readSlice checks that path is a thin 64-bit Mach-O and reads its CPU type and subtype.
func readSlice(path string) (slice, error) {
	f, err := os.Open(path)
	if err != nil {
		return slice{}, err
	}
	defer f.Close()
	var h [12]byte
	if _, err := io.ReadFull(f, h[:]); err != nil {
		return slice{}, fmt.Errorf("%s: %w", path, err)
	}
	if binary.LittleEndian.Uint32(h[0:]) != thinMagic {
		return slice{}, fmt.Errorf("%s: not a thin 64-bit Mach-O", path)
	}
	st, err := f.Stat()
	if err != nil {
		return slice{}, err
	}
	if st.Size() > 1<<32-1 {
		return slice{}, fmt.Errorf("%s: too large for a fat header", path)
	}
	return slice{path: path, cpu: binary.LittleEndian.Uint32(h[4:]), subtype: binary.LittleEndian.Uint32(h[8:]),
		size: uint32(st.Size())}, nil
}

// join writes the universal file: the fat header (big-endian), then each slice at an offset aligned to 2^14.
func join(out string, paths []string) error {
	if len(paths) < 2 {
		return errors.New("give two thin Mach-O files or more")
	}
	var slices []slice
	for _, p := range paths {
		if same(p, out) {
			return fmt.Errorf("%s is both an input and the output", p)
		}
		s, err := readSlice(p)
		if err != nil {
			return err
		}
		for _, o := range slices {
			if o.cpu == s.cpu {
				return fmt.Errorf("%s and %s are for the same CPU", o.path, s.path)
			}
		}
		slices = append(slices, s)
	}
	sort.Slice(slices, func(i, j int) bool { return slices[i].cpu < slices[j].cpu }) // x86_64 before arm64, as lipo
	align := uint64(1) << alignShift
	off := uint64(8 + 20*len(slices))
	for i := range slices {
		off = (off + align - 1) &^ (align - 1)
		if off+uint64(slices[i].size) > 1<<32-1 {
			return errors.New("the universal file would be over 4 GiB")
		}
		slices[i].offset = uint32(off)
		off += uint64(slices[i].size)
	}
	// written beside the output and renamed at the end: a failure leaves no half universal file
	f, err := os.CreateTemp(filepath.Dir(out), ".machofat-*")
	if err != nil {
		return err
	}
	done := false
	defer func() {
		if !done {
			f.Close()
			os.Remove(f.Name())
		}
	}()
	head := binary.BigEndian.AppendUint32(nil, fatMagic)
	head = binary.BigEndian.AppendUint32(head, uint32(len(slices)))
	for _, s := range slices {
		for _, v := range []uint32{s.cpu, s.subtype, s.offset, s.size, alignShift} {
			head = binary.BigEndian.AppendUint32(head, v)
		}
	}
	if _, err := f.Write(head); err != nil {
		return err
	}
	for _, s := range slices {
		in, err := os.Open(s.path)
		if err == nil {
			_, err = f.Seek(int64(s.offset), io.SeekStart)
		}
		if err == nil {
			_, err = io.Copy(f, in)
			in.Close()
		}
		if err != nil {
			return err
		}
	}
	if err := f.Chmod(0o755); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	done = true
	if err := os.Rename(f.Name(), out); err != nil {
		os.Remove(f.Name())
		return err
	}
	return nil
}

// same tells whether two paths are one file (an output that is also an input would be truncated before it is read).
func same(a, b string) bool {
	sa, err1 := os.Stat(a)
	sb, err2 := os.Stat(b)
	return err1 == nil && err2 == nil && os.SameFile(sa, sb)
}

func main() {
	out := flag.String("o", "", "the universal file to write")
	flag.Parse()
	if *out == "" {
		log.Fatal("usage: machofat -o out thin1 thin2...")
	}
	if err := join(*out, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
