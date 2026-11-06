package rules

import (
	"archive/tar"
	"crypto/sha1" //nolint:gosec // the XVA format's block checksums (Export.checksum_sha1)
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"strings"
)

// VerifyXVA is CommandLib's Export.verify: each data block of the archive is followed by its ".checksum" (SHA-1, lower
// case hex) or ".xxhash" (XXH64, upper case hex) entry, and a "checksum.xml" table, when there is one, must agree.
func VerifyXVA(r io.Reader) error {
	tail := &tailReader{r: r}
	tr := tar.NewReader(tail)
	recomputed := map[string]string{}
	var original map[string]string
	for entries := 0; ; entries++ {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			// a whole archive ends with two zero blocks (Export.nextHeader's EndOfArchive); tar also says EOF for an
			// empty or cut stream
			if entries == 0 || !tail.endsWithZeroBlocks() {
				return errors.New("the export is truncated (no end of archive)")
			}
			break
		}
		if err != nil {
			return fmt.Errorf("the export is not a valid XVA archive (header checksum failed): %w", err)
		}
		data, err := io.ReadAll(tr)
		if err != nil {
			return err
		}
		switch {
		case h.Name == "ova.xml":
		case strings.HasSuffix(h.Name, "checksum.xml"):
			if original, err = parseChecksumTable(data); err != nil {
				return err
			}
		default:
			// a data block: its checksum comes right after it
			ch, err := tr.Next()
			if err != nil {
				return fmt.Errorf("the block %s has no checksum: %w", h.Name, err)
			}
			want, err := io.ReadAll(tr)
			if err != nil {
				return err
			}
			got := fmt.Sprintf("%x", sha1.Sum(data))
			if strings.HasSuffix(ch.Name, ".xxhash") {
				got = fmt.Sprintf("%016X", XXH64(data))
			}
			if got != string(want) {
				return fmt.Errorf("failed to verify the block checksum: block = %s; recomputed = %s; original = %s", h.Name, got, want)
			}
			recomputed[h.Name] = got
		}
	}
	for name, sum := range recomputed { // compare_tables
		if original != nil && original[name] != sum {
			return fmt.Errorf("failed to verify the block checksum: block = %s; recomputed = %s; original = %s", name, sum, original[name])
		}
	}
	return nil
}

// parseChecksumTable is Export.parse_checksum_table: the <member><name/><value/></member> pairs of an XML-RPC struct.
func parseChecksumTable(data []byte) (map[string]string, error) {
	var doc struct {
		Members []struct {
			Name  string `xml:"name"`
			Value string `xml:"value"`
		} `xml:"struct>member"`
	}
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("the checksum table of the export is not valid: %w", err)
	}
	out := map[string]string{}
	for _, m := range doc.Members {
		out[m.Name] = m.Value
	}
	return out, nil
}

const (
	xxPrime1 uint64 = 11400714785074694791
	xxPrime2 uint64 = 14029467366897019727
	xxPrime3 uint64 = 1609587929392839161
	xxPrime4 uint64 = 9650029242287828579
	xxPrime5 uint64 = 2870177450012600261
)

// XXH64 is CommandLib's XXHash64 with seed 0 (the ".xxhash" block checksums of newer XVA exports).
func XXH64(b []byte) uint64 {
	n := uint64(len(b))
	var h uint64
	round := func(acc, in uint64) uint64 { return bits.RotateLeft64(acc+in*xxPrime2, 31) * xxPrime1 }
	merge := func(acc, v uint64) uint64 { return (acc^round(0, v))*xxPrime1 + xxPrime4 }
	if len(b) >= 32 {
		p1, p2 := xxPrime1, xxPrime2 // variables: the sums wrap around, as the C# ulong arithmetic
		v1, v2, v3, v4 := p1+p2, p2, uint64(0), -p1
		for ; len(b) >= 32; b = b[32:] {
			v1 = round(v1, binary.LittleEndian.Uint64(b))
			v2 = round(v2, binary.LittleEndian.Uint64(b[8:]))
			v3 = round(v3, binary.LittleEndian.Uint64(b[16:]))
			v4 = round(v4, binary.LittleEndian.Uint64(b[24:]))
		}
		h = bits.RotateLeft64(v1, 1) + bits.RotateLeft64(v2, 7) + bits.RotateLeft64(v3, 12) + bits.RotateLeft64(v4, 18)
		h = merge(merge(merge(merge(h, v1), v2), v3), v4)
	} else {
		h = xxPrime5
	}
	h += n
	for ; len(b) >= 8; b = b[8:] {
		h ^= round(0, binary.LittleEndian.Uint64(b))
		h = bits.RotateLeft64(h, 27)*xxPrime1 + xxPrime4
	}
	if len(b) >= 4 {
		h ^= uint64(binary.LittleEndian.Uint32(b)) * xxPrime1
		h = bits.RotateLeft64(h, 23)*xxPrime2 + xxPrime3
		b = b[4:]
	}
	for _, c := range b {
		h ^= uint64(c) * xxPrime5
		h = bits.RotateLeft64(h, 11) * xxPrime1
	}
	h ^= h >> 33
	h *= xxPrime2
	h ^= h >> 29
	h *= xxPrime3
	h ^= h >> 32
	return h
}

// tailReader keeps the last 1024 bytes read: the two zero blocks that end a tar archive.
type tailReader struct {
	r    io.Reader
	last []byte
}

func (t *tailReader) Read(p []byte) (int, error) {
	n, err := t.r.Read(p)
	t.last = append(t.last, p[:n]...)
	if len(t.last) > 1024 {
		t.last = t.last[len(t.last)-1024:]
	}
	return n, err
}

func (t *tailReader) endsWithZeroBlocks() bool {
	if len(t.last) < 1024 {
		return false
	}
	for _, b := range t.last {
		if b != 0 {
			return false
		}
	}
	return true
}
