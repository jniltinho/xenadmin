package rules

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

// A stream-optimized VMDK (the disks of an OVA from VMware or VirtualBox; DiscUtils.Vmdk in the C# app) read as the raw
// disk, as it comes: each compressed grain is preceded by a marker with its sector, so the disk is written in order
// with zeros between the grains, one grain in memory at a time.
// shortcut: stream-optimized VMDKs only (the hosted sparse and flat kinds need random access); the unallocated grains
// are sent as zeros (import_raw_vdi takes raw or VHD), convert to VHD when large sparse disks are too slow.

const (
	vmdkMagic          = 0x564d444b // "KDMV"
	vmdkFlagCompressed = 1 << 16
	vmdkFlagMarkers    = 1 << 17
	sector             = 512
)

// IsVMDK tells whether a disk starts with a VMDK sparse extent header.
func IsVMDK(head []byte) bool { return len(head) >= 4 && binary.LittleEndian.Uint32(head) == vmdkMagic }

type vmdkReader struct {
	in       io.Reader
	pos      int64 // bytes read from in
	size     int64 // the raw disk's size
	grain    int64 // bytes in a grain
	out      int64 // bytes given out
	zeros    int64 // zeros to give before buf
	buf      []byte
	done     bool
	compBuf  []byte
	grainBuf []byte
}

// NewVMDKReader reads a stream-optimized VMDK from r (at its first byte) as its raw disk; Size is the disk's size.
func NewVMDKReader(r io.Reader) (*vmdkReader, error) {
	h := make([]byte, sector)
	if _, err := io.ReadFull(r, h); err != nil {
		return nil, fmt.Errorf("the VMDK disk is truncated: %w", err)
	}
	le := binary.LittleEndian
	if le.Uint32(h) != vmdkMagic {
		return nil, errors.New("the disk is not a VMDK sparse extent")
	}
	flags := le.Uint32(h[8:])
	capacity, grain, overhead := le.Uint64(h[12:]), le.Uint64(h[20:]), le.Uint64(h[64:])
	if flags&vmdkFlagCompressed == 0 || flags&vmdkFlagMarkers == 0 || le.Uint16(h[77:]) != 1 {
		return nil, errors.New("only stream-optimized VMDK disks can be imported in this version")
	}
	if grain == 0 || grain > 2048 || capacity == 0 || capacity > math.MaxInt64/sector || overhead > capacity+1<<20 {
		return nil, errors.New("the VMDK disk's header is not valid")
	}
	v := &vmdkReader{in: r, pos: sector, size: int64(capacity) * sector, grain: int64(grain) * sector}
	if err := v.skip(int64(overhead)*sector - v.pos); err != nil {
		return nil, err
	}
	return v, nil
}

// Size is the raw disk's size in bytes.
func (v *vmdkReader) Size() int64 { return v.size }

func (v *vmdkReader) Read(p []byte) (int, error) {
	for {
		switch {
		case len(p) == 0:
			return 0, nil
		case v.zeros > 0:
			n := int(min(int64(len(p)), v.zeros))
			clear(p[:n])
			v.zeros -= int64(n)
			v.out += int64(n)
			return n, nil
		case len(v.buf) > 0:
			n := copy(p, v.buf)
			v.buf = v.buf[n:]
			v.out += int64(n)
			return n, nil
		case v.done:
			return 0, io.EOF
		}
		if err := v.next(); err != nil {
			return 0, err
		}
	}
}

func (v *vmdkReader) read(b []byte) error {
	n, err := io.ReadFull(v.in, b)
	v.pos += int64(n)
	if err != nil {
		return fmt.Errorf("the VMDK disk is truncated: %w", err)
	}
	return nil
}

func (v *vmdkReader) skip(n int64) error {
	if n < 0 {
		return errors.New("the VMDK disk is not valid (its data goes backwards)")
	}
	got, err := io.CopyN(io.Discard, v.in, n)
	v.pos += got
	if err != nil {
		return fmt.Errorf("the VMDK disk is truncated: %w", err)
	}
	return nil
}

// next reads the next marker: a grain becomes the zeros up to its sector and its data; the end of stream the zeros
// up to the disk's end; the grain tables, directory and footer are skipped.
func (v *vmdkReader) next() error {
	m := make([]byte, 12)
	if err := v.read(m); err != nil {
		return err
	}
	val, size := binary.LittleEndian.Uint64(m), int64(binary.LittleEndian.Uint32(m[8:]))
	if size == 0 { // a metadata marker: its type, the rest of its sector, then val sectors
		t := make([]byte, 4)
		if err := v.read(t); err != nil {
			return err
		}
		if err := v.skip(sector - 16); err != nil {
			return err
		}
		if binary.LittleEndian.Uint32(t) == 0 { // end of stream
			v.done, v.zeros = true, v.size-v.out
			return nil
		}
		if val > math.MaxInt64/sector {
			return errors.New("the VMDK disk is not valid")
		}
		return v.skip(int64(val) * sector)
	}
	if size > 2*v.grain+sector {
		return errors.New("the VMDK disk is not valid (a grain is too large)")
	}
	if int64(cap(v.compBuf)) < size {
		v.compBuf = make([]byte, size)
	}
	comp := v.compBuf[:size]
	if err := v.read(comp); err != nil {
		return err
	}
	if pad := (sector - (12+size)%sector) % sector; pad > 0 {
		if err := v.skip(pad); err != nil {
			return err
		}
	}
	if val > math.MaxInt64/sector {
		return errors.New("the VMDK disk is not valid")
	}
	at := int64(val) * sector
	if at%v.grain != 0 || at >= v.size {
		return fmt.Errorf("the VMDK disk is not valid (a grain at sector %d)", val)
	}
	if at < v.out {
		return errors.New("the VMDK disk's grains are not in order; it cannot be read as a stream")
	}
	z, err := zlib.NewReader(bytes.NewReader(comp))
	if err != nil {
		return fmt.Errorf("the VMDK disk has a damaged grain: %w", err)
	}
	if v.grainBuf == nil {
		v.grainBuf = make([]byte, v.grain)
	}
	// a whole grain (the last one at least up to the disk's end, it may go past it), then the end of the zlib stream,
	// which checks its checksum
	n, err := io.ReadFull(z, v.grainBuf)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) || int64(n) < min(v.grain, v.size-at) {
		return fmt.Errorf("the VMDK disk has a damaged grain at sector %d", val)
	}
	if extra, err := z.Read(make([]byte, 1)); extra != 0 || !errors.Is(err, io.EOF) {
		return fmt.Errorf("the VMDK disk has a damaged grain at sector %d", val)
	}
	v.zeros = at - v.out
	v.buf = v.grainBuf[:min(int64(n), v.size-at)]
	return nil
}
