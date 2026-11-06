package rules

import (
	"archive/tar"
	"bytes"
	"crypto/sha1" //nolint:gosec // the XVA format
	"fmt"
	"strings"
	"testing"
)

func TestXXH64(t *testing.T) {
	// the reference vectors of xxHash (seed 0)
	for in, want := range map[string]uint64{"": 0xEF46DB3751D8E999, "a": 0xD24EC4F1A98C6E5B, "abc": 0x44BC2CF5AD770999,
		"Nobody inspects the spammish repetition": 0xFBCEA83C8A378BF1} {
		if got := XXH64([]byte(in)); got != want {
			t.Errorf("%q: %016X, want %016X", in, got, want)
		}
	}
}

func TestVerifyXVA(t *testing.T) {
	block := bytes.Repeat([]byte("xcp-ng"), 2000)
	build := func(sum string, kind string) []byte {
		var buf bytes.Buffer
		w := tar.NewWriter(&buf)
		add := func(name string, data []byte) {
			_ = w.WriteHeader(&tar.Header{Name: name, Mode: 0o644, Size: int64(len(data))})
			_, _ = w.Write(data)
		}
		add("ova.xml", []byte("<value/>"))
		add("Ref:1/00000000", block)
		add("Ref:1/00000000"+kind, []byte(sum))
		_ = w.Close()
		return buf.Bytes()
	}
	good := fmt.Sprintf("%x", sha1.Sum(block))
	if err := VerifyXVA(bytes.NewReader(build(good, ".checksum"))); err != nil {
		t.Fatal(err)
	}
	if err := VerifyXVA(bytes.NewReader(build(fmt.Sprintf("%016X", XXH64(block)), ".xxhash"))); err != nil {
		t.Fatal(err)
	}
	if err := VerifyXVA(bytes.NewReader(build(strings.Repeat("0", 40), ".checksum"))); err == nil {
		t.Fatal("a wrong checksum passes")
	}
	if err := VerifyXVA(strings.NewReader("not a tar archive at all, but long enough to be read as a header block")); err == nil {
		t.Fatal("garbage passes")
	}
	whole := build(good, ".checksum")
	if err := VerifyXVA(bytes.NewReader(whole[:len(whole)-1024])); err == nil { // cut at an entry boundary
		t.Fatal("a truncated archive passes")
	}
	if err := VerifyXVA(bytes.NewReader(nil)); err == nil {
		t.Fatal("an empty file passes")
	}
}
