package rules

import (
	"archive/tar"
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"crypto/sha1" //nolint:gosec // the manifests of older OVF packages are SHA-1
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"xcpng-center/internal/messages"
)

// An OVF package (XenOvf's Package.cs): a folder with the .ovf descriptor beside its files (FolderPackage), or an OVA
// tar archive, gzip or bzip2 compressed or not (TarPackage). Nothing is extracted: the files are read from the archive
// as streams, so an entry's name is only ever compared, never used as a path.

// Package is an opened appliance package.
type Package struct {
	file        string
	dir         string           // a folder package's folder ("" for an archive)
	descriptor  string           // the .ovf's name
	entries     map[string]int64 // an archive's files and their sizes
	manifest    []byte
	certificate []byte   // "<name>.cert": the manifest's signature and the signer's certificate
	hrefs       []string // the descriptor and the files it names (a folder's manifest must list them)
	security    []ovfSecuritySection
	password    string          // an encrypted package's, once Unlock has checked it
	packed      map[string]bool // the files the descriptor says are encrypted (IsThisIdEncrypted) or compressed
	Appliance   *Appliance
}

const (
	ovfMaxManifest = 1 << 20
	ovfMaxEntries  = 10_000  // a real OVA has a descriptor, a manifest and its disks
	ovfMaxNames    = 1 << 20 // the bytes of all the archive's names
)

func (p *Package) archive() bool { return p.dir == "" }

// OpenPackage is OvfRules.Open / Load: the package read (descriptor and manifest), validated and described.
func OpenPackage(file string) (*Package, error) {
	if st, err := os.Stat(file); err != nil || st.IsDir() {
		return nil, errors.New("The file does not exist.")
	}
	if !IsAppliance(file) {
		return nil, errors.New("The file is not an OVF or OVA appliance.")
	}
	p := &Package{file: file}
	var desc []byte
	if strings.HasSuffix(strings.ToLower(file), ".ovf") {
		p.dir, p.descriptor = filepath.Dir(file), filepath.Base(file)
		var err error
		if desc, err = readLimited(file, ovfMaxDescriptor); err != nil {
			return nil, err
		}
		for _, m := range []struct {
			ext string
			to  *[]byte
		}{{".mf", &p.manifest}, {".cert", &p.certificate}} {
			if f, st, err := p.openInDir(p.name() + m.ext); err == nil {
				*m.to, err = io.ReadAll(io.LimitReader(f, ovfMaxManifest+1))
				f.Close()
				if err != nil || st.Size() > ovfMaxManifest {
					return nil, fmt.Errorf("the file %s of the appliance cannot be read or is too large", p.name()+m.ext)
				}
			}
		}
	} else {
		p.entries = map[string]int64{}
		// TarPackage.Load: the first .ovf, then its manifest; the other files are only listed
		err := p.walk(func(name string, h *tar.Header, r io.Reader) error {
			p.entries[name] = h.Size
			switch ext := strings.ToLower(path.Ext(name)); {
			case ext == ".ovf" && p.descriptor == "":
				p.descriptor = name
				b, err := io.ReadAll(io.LimitReader(r, ovfMaxDescriptor+1))
				if err != nil {
					return err
				}
				desc = b
			case (ext == ".mf" || ext == ".cert") && p.descriptor != "" && strings.EqualFold(name, p.name()+ext):
				b, err := io.ReadAll(io.LimitReader(r, ovfMaxManifest+1))
				if err != nil {
					return err
				}
				if len(b) > ovfMaxManifest {
					break
				}
				if ext == ".mf" {
					p.manifest = b
				} else {
					p.certificate = b
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		if p.descriptor == "" {
			return nil, errors.New("The appliance has no OVF descriptor.")
		}
	}
	env, err := parseOVF(desc)
	if err != nil {
		return nil, err
	}
	if p.Appliance, err = describeAppliance(env, p.name(), p.has, p.manifest != nil); err != nil {
		return nil, err
	}
	if p.certificate != nil { // ImportSecurityPage: HasSignature, and View Certificate's check
		p.Appliance.HasSignature = true
		if c, err := packageCertificate(p.certificate); err != nil {
			p.Appliance.CertificateError = messages.Get("IMPORT_SECURITY_PAGE_ERROR_MISSING_CERTIFICATE")
		} else {
			p.Appliance.Certificate = describeCertificate(c)
		}
	}
	p.security, p.packed = env.Security, map[string]bool{}
	for _, f := range env.Files {
		c := strings.ToLower(strings.TrimSpace(f.Compression))
		p.packed[f.Href] = p.packed[f.Href] || env.encrypted(f.ID) || c != "" && c != "identity"
	}
	p.hrefs = []string{p.descriptor}
	for _, f := range env.Files {
		if ext := strings.ToLower(path.Ext(f.Href)); ext != ".mf" && ext != ".cert" {
			p.hrefs = append(p.hrefs, f.Href)
		}
	}
	return p, nil
}

// name is Package.Name: the descriptor's name without its extension.
func (p *Package) name() string {
	b := path.Base(p.descriptor)
	return strings.TrimSuffix(b, path.Ext(b))
}

func readLimited(file string, limit int64) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("%s is too large", filepath.Base(file))
	}
	return b, nil
}

// entryName is an archive member's name as the descriptor names it ("./disk.vmdk" is "disk.vmdk").
func entryName(n string) string {
	return strings.TrimPrefix(path.Clean(strings.TrimPrefix(n, "./")), "/")
}

// open is the archive's tar stream, decompressed by the file's extension.
func (p *Package) open() (*tar.Reader, io.Closer, error) {
	f, err := os.Open(p.file)
	if err != nil {
		return nil, nil, err
	}
	var r io.Reader = f // a plain tar: tar skips the members it is not asked for by seeking
	switch name := strings.ToLower(p.file); {
	case strings.HasSuffix(name, ".gz"):
		z, err := gzip.NewReader(bufio.NewReaderSize(f, 1<<20))
		if err != nil {
			f.Close()
			return nil, nil, fmt.Errorf("the appliance is not a valid gzip file: %w", err)
		}
		r = z
	case strings.HasSuffix(name, ".bz2"):
		r = bzip2.NewReader(bufio.NewReaderSize(f, 1<<20))
	}
	return tar.NewReader(r), f, nil
}

// walk calls fn for each regular file of the archive, in its order.
func (p *Package) walk(fn func(name string, h *tar.Header, r io.Reader) error) error {
	tr, c, err := p.open()
	if err != nil {
		return err
	}
	defer c.Close()
	for entries, names := 1, 0; ; entries++ {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("the appliance is not a valid OVA archive: %w", err)
		}
		if names += len(h.Name) + len(h.Linkname); entries > ovfMaxEntries || names > ovfMaxNames { // every header counts
			return errors.New("the appliance has too many files")
		}
		if h.Typeflag != tar.TypeReg && h.Typeflag != tar.TypeRegA {
			continue
		}
		if err := fn(entryName(h.Name), h, tr); err != nil {
			return err
		}
	}
}

func (p *Package) has(href string) bool {
	if p.archive() {
		_, ok := p.entries[href]
		return ok
	}
	f, _, err := p.openInDir(href)
	if err == nil {
		f.Close()
	}
	return err == nil
}

// openInDir opens a regular file of a folder package; os.Root keeps a symbolic link from leading out of the folder.
func (p *Package) openInDir(href string) (*os.File, os.FileInfo, error) {
	root, err := os.OpenRoot(p.dir)
	if err != nil {
		return nil, nil, err
	}
	defer root.Close()
	// a FIFO or device would block the open: only a regular file is opened (and checked again once open)
	if st, err := root.Stat(filepath.FromSlash(href)); err != nil || !st.Mode().IsRegular() {
		return nil, nil, fmt.Errorf("%s is not a file of the package", href)
	}
	f, err := root.Open(filepath.FromSlash(href))
	if err != nil {
		return nil, nil, err
	}
	st, err := f.Stat()
	if err != nil || !st.Mode().IsRegular() {
		f.Close()
		return nil, nil, fmt.Errorf("%s is not a file of the package", href)
	}
	return f, st, nil
}

type readCloser struct {
	io.Reader
	io.Closer
}

// Open is a file of the package as a stream, and its size.
func (p *Package) Open(href string) (io.ReadCloser, int64, error) {
	if !SafeHref(href) {
		return nil, 0, fmt.Errorf("the appliance names a file outside the package: %q", href)
	}
	if !p.archive() {
		f, st, err := p.openInDir(href)
		if err != nil {
			return nil, 0, errors.New(messages.Format(messages.Get("ERROR_FILE_NOT_FOUND"), href))
		}
		return f, st.Size(), nil
	}
	tr, c, err := p.open()
	if err != nil {
		return nil, 0, err
	}
	for {
		h, err := tr.Next()
		if err != nil {
			c.Close()
			if errors.Is(err, io.EOF) {
				return nil, 0, errors.New(messages.Format(messages.Get("ERROR_FILE_NOT_FOUND"), href))
			}
			return nil, 0, fmt.Errorf("the appliance is not a valid OVA archive: %w", err)
		}
		if (h.Typeflag == tar.TypeReg || h.Typeflag == tar.TypeRegA) && entryName(h.Name) == href {
			return readCloser{tr, c}, h.Size, nil
		}
	}
}

// openAt is a file of the package with random access, which a hosted sparse VMDK needs: a folder's file, or a member
// of an uncompressed OVA as a section of the archive (tar reads exactly up to a member's data).
// shortcut: not from a compressed OVA (.ova.gz, .ova.bz2); the C# app extracts the archive first, extract it by hand.
func (p *Package) openAt(href string) (io.ReaderAt, int64, io.Closer, error) {
	if !SafeHref(href) {
		return nil, 0, nil, fmt.Errorf("the appliance names a file outside the package: %q", href)
	}
	if !p.archive() {
		f, st, err := p.openInDir(href)
		if err != nil {
			return nil, 0, nil, errors.New(messages.Format(messages.Get("ERROR_FILE_NOT_FOUND"), href))
		}
		return f, st.Size(), f, nil
	}
	if name := strings.ToLower(p.file); strings.HasSuffix(name, ".gz") || strings.HasSuffix(name, ".bz2") {
		return nil, 0, nil, fmt.Errorf("the disk %s is a hosted sparse VMDK, which cannot be read from a compressed OVA; extract the OVA and import its .ovf", href)
	}
	f, err := os.Open(p.file)
	if err != nil {
		return nil, 0, nil, err
	}
	fail := func(err error) (io.ReaderAt, int64, io.Closer, error) { f.Close(); return nil, 0, nil, err }
	tr := tar.NewReader(f)
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return fail(errors.New(messages.Format(messages.Get("ERROR_FILE_NOT_FOUND"), href)))
		}
		if err != nil {
			return fail(fmt.Errorf("the appliance is not a valid OVA archive: %w", err))
		}
		if h.Typeflag != tar.TypeReg && h.Typeflag != tar.TypeRegA || entryName(h.Name) != href {
			continue
		}
		for k := range h.PAXRecords { // a PAX sparse file's data is not the member's bytes as they lie
			if strings.HasPrefix(k, "GNU.sparse.") {
				return fail(fmt.Errorf("the file %s of the appliance is a sparse tar member", href))
			}
		}
		at, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			return fail(err)
		}
		st, err := f.Stat()
		if err != nil || h.Size > st.Size()-at {
			return fail(fmt.Errorf("the file %s of the appliance is truncated", href))
		}
		return io.NewSectionReader(f, at, h.Size), h.Size, f, nil
	}
}

// fileDigest is a manifest line (FileDigest): "SHA256(disk.vmdk)= 0123...".
type fileDigest struct{ algorithm, name, digest string }

var manifestLine = regexp.MustCompile(`^\s*([A-Za-z0-9-]+)\s*\((.*)\)\s*=\s*([0-9A-Fa-f]+)\s*$`)

func (p *Package) digests() ([]fileDigest, error) {
	var out []fileDigest
	for _, line := range strings.Split(strings.ReplaceAll(string(p.manifest), "\r", ""), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		m := manifestLine.FindStringSubmatch(line)
		if m == nil {
			return nil, fmt.Errorf("the manifest line %q is not valid", line)
		}
		out = append(out, fileDigest{strings.ToUpper(m[1]), strings.TrimSpace(m[2]), strings.ToLower(m[3])})
	}
	return out, nil
}

func newHash(algorithm string) (hash.Hash, error) {
	switch strings.ReplaceAll(algorithm, "-", "") {
	case "SHA1":
		return sha1.New(), nil //nolint:gosec // the manifest's own algorithm
	case "SHA256":
		return sha256.New(), nil
	case "SHA512":
		return sha512.New(), nil
	}
	return nil, fmt.Errorf("the manifest's algorithm %s is not supported", algorithm)
}

// VerifyManifest is Package.VerifyManifest: every file of the package (but the manifest and the certificate) is
// listed in the manifest with its digest, and every file the manifest lists is in the package.
func (p *Package) VerifyManifest() error {
	if p.manifest == nil {
		return errors.New("the appliance has no manifest")
	}
	want, err := p.digests()
	if err != nil {
		return err
	}
	check := func(d fileDigest, r io.Reader) error {
		h, err := newHash(d.algorithm)
		if err != nil {
			return err
		}
		if _, err := io.Copy(h, r); err != nil {
			return err
		}
		if hex.EncodeToString(h.Sum(nil)) != d.digest {
			return fmt.Errorf("the file %s does not match its digest in the manifest", d.name)
		}
		return nil
	}
	if !p.archive() { // FolderPackage: by the manifest's order, which must list the descriptor and every file it names
		for _, h := range p.hrefs {
			if !slices.ContainsFunc(want, func(d fileDigest) bool { return d.name == h }) {
				return fmt.Errorf("the file %s of the appliance is not listed in its manifest", h)
			}
		}
		for _, d := range want {
			r, _, err := p.Open(d.name)
			if err != nil {
				return err
			}
			err = check(d, r)
			r.Close()
			if err != nil {
				return err
			}
		}
		return nil
	}
	left := append([]fileDigest{}, want...)
	err = p.walk(func(name string, _ *tar.Header, r io.Reader) error { // TarPackage: by the archive's order
		if ext := strings.ToLower(path.Ext(name)); ext == ".mf" || ext == ".cert" {
			return nil
		}
		for i, d := range left {
			if strings.EqualFold(d.name, name) {
				left = append(left[:i], left[i+1:]...)
				return check(d, r)
			}
		}
		return fmt.Errorf("the file %s of the appliance is not listed in its manifest", name)
	})
	if err != nil {
		return err
	}
	if len(left) > 0 {
		return fmt.Errorf("the file %s listed in the manifest is missing from the appliance", left[0].name)
	}
	return nil
}

// DiskStream is a disk of the package as import_raw_vdi takes it.
type DiskStream struct {
	io.Reader
	io.Closer
	Length      int64  // the bytes to send
	VirtualSize int64  // the disk's size
	Format      string // "raw" or "vhd"
}

// OpenDisk is ImportApplianceAction.ImportFile's "open disk" part: a VHD goes as VHD (its size from its footer), a
// VMDK (stream-optimized, hosted sparse, or a descriptor of flat or sparse extents) as its raw disk, a .raw or .img
// file as raw (DiscUtils' disk classes there); an encrypted file is decrypted and a gzip-compressed one uncompressed
// as it is read (shortcut: its size is found by reading it once more, where the C# app writes a temporary file).
func (p *Package) OpenDisk(d ApplianceDisk) (*DiskStream, error) {
	ext := strings.ToLower(d.Href)
	if d.Compression == "gzip" {
		ext = strings.TrimSuffix(ext, ".gz")
	}
	ext = path.Ext(ext)
	openData := func() (io.Reader, io.Closer, int64, error) {
		rc, n, err := p.Open(d.Href)
		if err != nil {
			return nil, nil, 0, err
		}
		var r io.Reader = rc
		if d.Encrypted { // DecryptToTempFile, before the uncompression
			if p.password == "" {
				rc.Close()
				return nil, nil, 0, errors.New(messages.Get("ERROR_NO_PASSWORD"))
			}
			r, n = newOVFDecrypter(rc, p.password), -1 // its size without the padding is known at its end
		}
		if d.Compression != "gzip" {
			return r, rc, n, nil
		}
		z, err := gzip.NewReader(bufio.NewReaderSize(r, 1<<20))
		if err != nil {
			rc.Close()
			return nil, nil, 0, fmt.Errorf("the disk %s is not a valid gzip file: %w", d.Href, err)
		}
		return z, rc, -1, nil
	}
	r, c, n, err := openData()
	if err != nil {
		return nil, err
	}
	br := bufio.NewReaderSize(r, 1<<20)
	head, _ := br.Peek(sector)
	fail := func(err error) (*DiskStream, error) { c.Close(); return nil, err }
	if IsVMDK(head) && !isStreamVMDK(head) { // a hosted sparse VMDK: read again with random access
		c.Close()
		// shortcut: not decrypted either (CBC can be read at random, a block and the one before it; the C# app decrypts
		// to a temporary file); decrypt such a disk by hand if one shows up
		if d.Compression != "" || d.Encrypted {
			return nil, fmt.Errorf("the disk %s is a hosted sparse VMDK, which cannot be read compressed or encrypted", d.Href)
		}
		ra, size, rc, err := p.openAt(d.Href)
		if err != nil {
			return nil, err
		}
		v, err := NewSparseVMDKReader(ra, size)
		if err != nil {
			rc.Close()
			return nil, err
		}
		return &DiskStream{Reader: v, Closer: rc, Length: v.Size(), VirtualSize: v.Size(), Format: "raw"}, nil
	}
	if IsVMDK(head) {
		v, err := NewVMDKReader(br)
		if err != nil {
			return fail(err)
		}
		// the file read to its end after the disk's end marker, so that its decryption checks the padding (and gzip its
		// checksum)
		return &DiskStream{Reader: drainAtEOF{v, br}, Closer: c, Length: v.Size(), VirtualSize: v.Size(), Format: "raw"}, nil
	}
	switch ext {
	case ".vhd", ".raw", ".img":
	case ".vmdk": // a descriptor: its extents, other files of the package, are the raw disk
		desc, err := io.ReadAll(io.LimitReader(br, vmdkMaxDescriptor+1))
		c.Close()
		if err != nil {
			return nil, err
		}
		if len(desc) > vmdkMaxDescriptor {
			desc = nil // not a descriptor: vmdkExtents says which kinds are read
		}
		extents, err := vmdkExtents(desc)
		if err != nil {
			return nil, err
		}
		x, err := p.openExtents(d.Href, extents)
		if err != nil {
			return nil, err
		}
		return &DiskStream{Reader: x, Closer: x, Length: x.size, VirtualSize: x.size, Format: "raw"}, nil
	default:
		return fail(errors.New(messages.Format(messages.Get("UNSUPPORTED_FILE_TYPE"), ext)))
	}
	if n < 0 { // a compressed disk: its size, read once, then the stream again
		if n, err = io.Copy(io.Discard, br); err != nil {
			return fail(err)
		}
		c.Close()
		if r, c, _, err = openData(); err != nil {
			return nil, err
		}
		br = bufio.NewReaderSize(r, 1<<20)
	}
	s := &DiskStream{Reader: br, Closer: c, Length: n, VirtualSize: n, Format: "raw"}
	if ext == ".vhd" {
		s.Format = "vhd"
		if len(head) >= sector && string(head[:8]) == "conectix" { // a dynamic VHD starts with a copy of its footer
			s.VirtualSize = int64(binary.BigEndian.Uint64(head[48:]))
		} else { // a fixed VHD is the disk then its footer
			s.VirtualSize = n - sector
		}
		if s.VirtualSize <= 0 {
			return fail(fmt.Errorf("the disk %s is not a valid VHD", d.Href))
		}
	}
	if s.VirtualSize <= 0 {
		return fail(fmt.Errorf("the disk %s is empty", d.Href))
	}
	return s, nil
}

// drainAtEOF is a disk read from a file that goes on after it: at the disk's end, the rest of the file is read.
type drainAtEOF struct {
	io.Reader
	rest io.Reader
}

func (d drainAtEOF) Read(p []byte) (int, error) {
	n, err := d.Reader.Read(p)
	if errors.Is(err, io.EOF) {
		if _, derr := io.Copy(io.Discard, d.rest); derr != nil {
			return n, derr
		}
	}
	return n, err
}

// fileSize is the size of a file of the package, without reading it.
func (p *Package) fileSize(href string) (int64, error) {
	if SafeHref(href) {
		if n, ok := p.entries[href]; ok && p.archive() {
			return n, nil
		}
		if f, st, err := p.openInDir(href); err == nil && !p.archive() {
			f.Close()
			return st.Size(), nil
		}
	}
	return 0, errors.New(messages.Format(messages.Get("ERROR_FILE_NOT_FOUND"), href))
}

// extentsReader is a VMDK descriptor's extents read one after the other as the raw disk, each opened when the one
// before it is done: one file open at a time, one grain table or read buffer in memory.
type extentsReader struct {
	p       *Package
	extents []vmdkExtent // those not read yet, their files named as in the package
	r       io.Reader
	c       io.Closer
	left    int64 // bytes the current extent has still to give
	size    int64 // the disk's
}

// openExtents checks a descriptor's extents before the disk is read: each one a disk file of the package (not the
// descriptor, nor a file the manifest check skips), listed in the manifest when there is one, as large as the
// descriptor says (a sparse extent's header capacity is its size).
func (p *Package) openExtents(disk string, extents []vmdkExtent) (*extentsReader, error) {
	var ds []fileDigest
	if p.manifest != nil {
		var err error
		if ds, err = p.digests(); err != nil {
			return nil, err
		}
	}
	x := &extentsReader{p: p}
	for _, e := range extents {
		e.file = path.Join(path.Dir(disk), e.file)
		// the OVF does not name the extents: the manifest must, so that its check covers the disk's data (and the check
		// skips a .mf or .cert file)
		if ext := strings.ToLower(path.Ext(e.file)); ext == ".mf" || ext == ".cert" || ext == ".ovf" || e.file == disk {
			return nil, fmt.Errorf("the extent %s of the disk %s is not a disk file", e.file, disk)
		}
		// shortcut: an extent the package lists as encrypted or compressed is refused, not decrypted nor uncompressed (no
		// XenCenter nor VMware wrote one)
		if p.packed[e.file] {
			return nil, fmt.Errorf("the extent %s of the disk %s is encrypted or compressed, which cannot be read in this version", e.file, disk)
		}
		// a folder's check opens the manifest's name exactly; an OVA's matches its members' names in any case
		same := func(d fileDigest) bool { return d.name == e.file || p.archive() && strings.EqualFold(d.name, e.file) }
		if p.manifest != nil && !slices.ContainsFunc(ds, same) {
			return nil, fmt.Errorf("the extent %s of the disk %s is not listed in the appliance's manifest", e.file, disk)
		}
		if e.sparse { // its header, read now (an OVA member is found at random, without reading the archive through)
			_, c, err := x.open(e)
			if err != nil {
				return nil, err
			}
			c.Close()
		} else if n, err := p.fileSize(e.file); err != nil {
			return nil, err
		} else if e.offset > n || e.size > n-e.offset {
			return nil, fmt.Errorf("the extent %s of the disk %s is truncated", e.file, disk)
		}
		x.size += e.size // vmdkExtents keeps the sum an int64
		x.extents = append(x.extents, e)
	}
	return x, nil
}

// open is an extent's data and the file to close after it.
func (x *extentsReader) open(e vmdkExtent) (io.Reader, io.Closer, error) {
	if e.sparse {
		ra, n, c, err := x.p.openAt(e.file)
		if err != nil {
			return nil, nil, err
		}
		v, err := NewSparseVMDKReader(ra, n)
		if err == nil && v.Size() != e.size {
			err = fmt.Errorf("the extent %s is not the size its VMDK descriptor says", e.file)
		}
		if err != nil {
			c.Close()
			return nil, nil, err
		}
		return v, c, nil
	}
	rc, n, err := x.p.Open(e.file)
	if err != nil {
		return nil, nil, err
	}
	if e.offset > n || e.size > n-e.offset {
		rc.Close()
		return nil, nil, fmt.Errorf("the extent %s is truncated", e.file)
	}
	if _, err := io.CopyN(io.Discard, rc, e.offset); err != nil {
		rc.Close()
		return nil, nil, err
	}
	return io.LimitReader(rc, e.size), rc, nil
}

func (x *extentsReader) Read(b []byte) (int, error) {
	for {
		if x.r == nil {
			if len(x.extents) == 0 {
				return 0, io.EOF
			}
			r, c, err := x.open(x.extents[0])
			if err != nil {
				return 0, err
			}
			x.r, x.c, x.left, x.extents = r, c, x.extents[0].size, x.extents[1:]
		}
		n, err := x.r.Read(b)
		x.left -= int64(n)
		if !errors.Is(err, io.EOF) {
			return n, err
		}
		x.c.Close()
		x.r, x.c = nil, nil
		if x.left != 0 {
			return n, errors.New("an extent of the VMDK disk is truncated")
		}
		if n > 0 {
			return n, nil
		}
	}
}

func (x *extentsReader) Close() error {
	x.extents = nil
	if x.c == nil {
		return nil
	}
	c := x.c
	x.r, x.c = nil, nil
	return c.Close()
}
