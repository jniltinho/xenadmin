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
	file       string
	dir        string           // a folder package's folder ("" for an archive)
	descriptor string           // the .ovf's name
	entries    map[string]int64 // an archive's files and their sizes
	manifest   []byte
	hrefs      []string // the descriptor and the files it names (a folder's manifest must list them)
	Appliance  *Appliance
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
		if f, st, err := p.openInDir(p.name() + ".mf"); err == nil {
			p.manifest, err = io.ReadAll(io.LimitReader(f, ovfMaxManifest+1))
			f.Close()
			if err != nil || st.Size() > ovfMaxManifest {
				return nil, errors.New("the manifest of the appliance cannot be read or is too large")
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
			case ext == ".mf" && p.descriptor != "" && strings.EqualFold(name, p.name()+".mf"):
				b, err := io.ReadAll(io.LimitReader(r, ovfMaxManifest+1))
				if err != nil {
					return err
				}
				if len(b) <= ovfMaxManifest {
					p.manifest = b
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
// stream-optimized VMDK as its raw disk, a .raw or .img file as raw; a gzip-compressed file is uncompressed as it is
// read (shortcut: its size is found by reading it once more, where the C# app uncompresses it to a temporary file).
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
		if d.Compression != "gzip" {
			return rc, rc, n, nil
		}
		z, err := gzip.NewReader(bufio.NewReaderSize(rc, 1<<20))
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
	if IsVMDK(head) {
		v, err := NewVMDKReader(br)
		if err != nil {
			return fail(err)
		}
		return &DiskStream{Reader: v, Closer: c, Length: v.Size(), VirtualSize: v.Size(), Format: "raw"}, nil
	}
	switch ext {
	case ".vhd", ".raw", ".img":
	case ".vmdk":
		return fail(errors.New("only stream-optimized VMDK disks can be imported in this version"))
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
