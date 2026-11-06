package rules

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// A stream-optimized VMDK (the disks of an OVA from VMware or VirtualBox; DiscUtils.Vmdk in the C# app) read as the raw
// disk, as it comes: each compressed grain is preceded by a marker with its sector, so the disk is written in order
// with zeros between the grains, one grain in memory at a time. The hosted sparse VMDK (HostedSparseExtentFile: the
// monolithicSparse disk qemu-img and VMware Workstation write) is read through random access (sparseVMDK), and a
// monolithic flat or split (twoGbMaxExtent) VMDK is its descriptor's extents one after the other (vmdkExtents).
// shortcut: the unallocated grains are sent as zeros (import_raw_vdi takes raw or VHD), convert to VHD when large
// sparse disks are too slow.

const (
	vmdkMagic          = 0x564d444b // "KDMV"
	vmdkFlagCompressed = 1 << 16
	vmdkFlagMarkers    = 1 << 17
	sector             = 512
)

// IsVMDK tells whether a disk starts with a VMDK sparse extent header.
func IsVMDK(head []byte) bool { return len(head) >= 4 && binary.LittleEndian.Uint32(head) == vmdkMagic }

// isStreamVMDK tells whether a sparse extent header is stream-optimized (compressed grains, markers), which
// NewVMDKReader reads as a stream; the others are hosted sparse extents (NewSparseVMDKReader).
func isStreamVMDK(head []byte) bool {
	return len(head) >= 12 && binary.LittleEndian.Uint32(head[8:])&(vmdkFlagCompressed|vmdkFlagMarkers) != 0
}

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

const vmdkFlagZeroGrain = 1 << 2 // a grain table entry of 1 is a grain of zeros

var vmdkParentCID = regexp.MustCompile(`(?m)^\s*parentCID\s*=\s*"?([0-9A-Za-z]+)`)

// sparseVMDK is a hosted sparse VMDK (DiscUtils.Vmdk's HostedSparseExtentFile) read as its raw disk in order through
// random access, one grain table in memory at a time. The file is untrusted: the grain directory and tables must lie
// in the metadata before overHead and every grain after it and inside the file, so a table or a grain never aliases
// the header or the tables. Two entries naming the same grain only repeat its data: the disk's size stays capacity.
type sparseVMDK struct {
	r       io.ReaderAt
	file    int64 // the file's size
	size    int64 // the raw disk's size
	grain   int64 // bytes in a grain
	gtes    int64 // entries in a grain table
	gd      int64 // the grain directory's offset
	data    int64 // where the grains start (overHead)
	zeroGTE bool
	out     int64 // bytes given out
	table   int64 // the grain table in gt (-1: none read yet)
	gt      []byte
	noGT    bool // that table is not allocated: its grains are zeros
}

// NewSparseVMDKReader reads a hosted sparse VMDK of size bytes from r as its raw disk; Size is the disk's size.
func NewSparseVMDKReader(r io.ReaderAt, size int64) (*sparseVMDK, error) {
	h := make([]byte, sector)
	if n, _ := r.ReadAt(h, 0); n < sector {
		return nil, errors.New("the VMDK disk is truncated")
	}
	le := binary.LittleEndian
	if le.Uint32(h) != vmdkMagic || isStreamVMDK(h) {
		return nil, errors.New("the disk is not a hosted sparse VMDK")
	}
	version, flags := le.Uint32(h[4:]), le.Uint32(h[8:])
	capacity, grain, gtes := le.Uint64(h[12:]), le.Uint64(h[20:]), uint64(le.Uint32(h[44:]))
	gd, overhead := le.Uint64(h[56:]), le.Uint64(h[64:])
	// qemu's vmdk_open_vmdk4 limits: a power-of-two grain, at most 512 entries in a table
	if version < 1 || version > 3 || grain < 8 || grain > 2048 || grain&(grain-1) != 0 || gtes == 0 || gtes > 512 ||
		capacity == 0 || capacity > math.MaxInt64/sector || size < 0 || overhead > uint64(size)/sector || gd == 0 || gd >= overhead {
		return nil, errors.New("the VMDK disk's header is not valid")
	}
	if tables := (capacity + grain*gtes - 1) / (grain * gtes); tables > (overhead-gd)*sector/4 {
		return nil, errors.New("the VMDK disk's grain directory is not valid")
	}
	// the embedded descriptor: a snapshot's delta (a parent) holds only its changes, not the disk
	if descAt, descLen := le.Uint64(h[28:]), le.Uint64(h[36:]); descAt != 0 {
		if descLen > vmdkMaxDescriptor/sector || descAt > overhead || descLen > overhead-descAt {
			return nil, errors.New("the VMDK disk's descriptor is not valid")
		}
		desc := make([]byte, descLen*sector)
		if n, _ := r.ReadAt(desc, int64(descAt)*sector); n < len(desc) {
			return nil, errors.New("the VMDK disk is truncated")
		}
		if m := vmdkParentCID.FindSubmatch(desc); m != nil && !strings.EqualFold(string(m[1]), "ffffffff") {
			return nil, errors.New("the VMDK disk is a snapshot of another disk, which cannot be imported on its own")
		}
	}
	return &sparseVMDK{r: r, file: size, size: int64(capacity) * sector, grain: int64(grain) * sector, gtes: int64(gtes),
		gd: int64(gd) * sector, data: int64(overhead) * sector, zeroGTE: flags&vmdkFlagZeroGrain != 0, table: -1}, nil
}

// Size is the raw disk's size in bytes.
func (v *sparseVMDK) Size() int64 { return v.size }

func (v *sparseVMDK) readAt(b []byte, at int64) error {
	if n, err := v.r.ReadAt(b, at); n < len(b) {
		return fmt.Errorf("the VMDK disk is truncated: %v", err)
	}
	return nil
}

// locate is where grain g is in the file, -1 for a grain of zeros.
func (v *sparseVMDK) locate(g int64) (int64, error) {
	if t := g / v.gtes; t != v.table {
		b := make([]byte, 4)
		if err := v.readAt(b, v.gd+t*4); err != nil {
			return 0, err
		}
		at := int64(binary.LittleEndian.Uint32(b)) * sector
		v.table, v.noGT = t, at == 0
		if !v.noGT {
			if at < sector || at+v.gtes*4 > v.data {
				return 0, fmt.Errorf("the VMDK disk is not valid (grain table %d)", t)
			}
			if v.gt == nil {
				v.gt = make([]byte, v.gtes*4)
			}
			if err := v.readAt(v.gt, at); err != nil {
				return 0, err
			}
		}
	}
	if v.noGT {
		return -1, nil
	}
	gte := int64(binary.LittleEndian.Uint32(v.gt[g%v.gtes*4:]))
	if gte == 0 || gte == 1 && v.zeroGTE {
		return -1, nil
	}
	if at := gte * sector; at >= v.data && at+min(v.grain, v.size-g*v.grain) <= v.file {
		return at, nil
	}
	return 0, fmt.Errorf("the VMDK disk is not valid (a grain at sector %d)", gte)
}

func (v *sparseVMDK) Read(p []byte) (int, error) {
	if v.out >= v.size {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return 0, nil
	}
	in := v.out % v.grain
	n := min(int64(len(p)), v.grain-in, v.size-v.out) // within one grain
	at, err := v.locate(v.out / v.grain)
	if err != nil {
		return 0, err
	}
	if at < 0 {
		clear(p[:n])
	} else if err := v.readAt(p[:n], at+in); err != nil {
		return 0, err
	}
	v.out += n
	return int(n), nil
}

const vmdkMaxDescriptor = 64 << 10

var vmdkExtentStart = regexp.MustCompile(`^(RW|RDONLY|NOACCESS)(\s|$)`)

var vmdkExtentLine = regexp.MustCompile(`^(RW|RDONLY|NOACCESS)\s+(\d+)\s+(\w+)(?:\s+"([^"]*)"(?:\s+(\d+))?)?\s*$`)

// vmdkExtent is an extent of a VMDK descriptor: its file (relative to the descriptor), a flat file of raw data (where
// the data starts in it) or a hosted sparse extent, and its size in bytes.
type vmdkExtent struct {
	file         string
	sparse       bool
	offset, size int64
}

const vmdkMaxExtents = 1 << 14 // of 2 GB each: 32 TB, past the largest VDI

// vmdkExtents is a VMDK descriptor (DiscUtils.Vmdk's DescriptorFile) read for its extents, the disk being their data
// one after the other: the monolithic flat kind (one FLAT or ESX VMFS extent) and the split twoGbMaxExtentFlat and
// twoGbMaxExtentSparse kinds (FLAT or SPARSE extents of up to 2 GB each), as DiscUtils.Vmdk.DiskImageFile opens them.
// The descriptor is untrusted: each extent names its own file, and a snapshot's descriptor (a parent) is refused.
// shortcut: no ZERO or NOACCESS extents nor the other ESX kinds (VMFSSPARSE, SESPARSE), which neither qemu-img nor
// VMware's OVF export writes; add them when such packages show up.
func vmdkExtents(desc []byte) ([]vmdkExtent, error) {
	unsupported := errors.New("only stream-optimized, hosted sparse, monolithic flat and split (twoGbMaxExtent) VMDK disks can be imported in this version")
	var out []vmdkExtent
	var total int64 // sectors
	seen := map[string]bool{}
	for _, line := range strings.Split(strings.ReplaceAll(string(desc), "\r", ""), "\n") {
		line = strings.TrimSpace(strings.Trim(line, "\x00"))
		m := vmdkExtentLine.FindStringSubmatch(line)
		if m == nil {
			if vmdkExtentStart.MatchString(line) { // an extent that cannot be read: the others would be out of place
				return nil, errors.New("the VMDK descriptor's extent is not valid")
			}
			continue
		}
		sparse := m[3] == "SPARSE"
		if m[1] == "NOACCESS" || m[4] == "" || !sparse && m[3] != "FLAT" && m[3] != "VMFS" || len(out) == vmdkMaxExtents {
			return nil, unsupported
		}
		sectors, err1 := strconv.ParseInt(m[2], 10, 64)
		start, err2 := int64(0), error(nil)
		if m[5] != "" {
			start, err2 = strconv.ParseInt(m[5], 10, 64)
		}
		if err1 != nil || err2 != nil || sectors <= 0 || sectors > math.MaxInt64/sector-total || start > math.MaxInt64/sector-sectors ||
			sparse && start != 0 { // a sparse extent has its own header, not an offset
			return nil, errors.New("the VMDK descriptor's extent is not valid")
		}
		name := path.Clean(m[4])
		if seen[name] {
			return nil, fmt.Errorf("the VMDK descriptor names the extent %s twice", m[4])
		}
		seen[name] = true
		total += sectors
		out = append(out, vmdkExtent{file: m[4], sparse: sparse, offset: start * sector, size: sectors * sector})
	}
	if len(out) == 0 {
		return nil, unsupported
	}
	if m := vmdkParentCID.FindSubmatch(desc); m != nil && !strings.EqualFold(string(m[1]), "ffffffff") {
		return nil, errors.New("the VMDK disk is a snapshot of another disk, which cannot be imported on its own")
	}
	return out, nil
}
