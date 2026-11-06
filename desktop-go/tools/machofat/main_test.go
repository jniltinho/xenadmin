package main

import (
	"debug/macho"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestJoin builds this tool for darwin/amd64 and darwin/arm64 and joins them; debug/macho reads the universal file
// back with both slices intact.
func TestJoin(t *testing.T) {
	dir := t.TempDir()
	var thin []string
	for _, arch := range []string{"arm64", "amd64"} {
		p := filepath.Join(dir, arch)
		cmd := exec.Command("go", "build", "-o", p, ".")
		cmd.Env = append(os.Environ(), "GOOS=darwin", "GOARCH="+arch, "CGO_ENABLED=0")
		if b, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s: %v %s", arch, err, b)
		}
		thin = append(thin, p)
	}
	out := filepath.Join(dir, "universal")
	if err := join(out, thin); err != nil {
		t.Fatal(err)
	}
	ff, err := macho.OpenFat(out)
	if err != nil {
		t.Fatal(err)
	}
	defer ff.Close()
	if len(ff.Arches) != 2 || ff.Arches[0].Cpu != macho.CpuAmd64 || ff.Arches[1].Cpu != macho.CpuArm64 {
		t.Fatalf("%+v", ff.Arches)
	}
	for i, a := range ff.Arches {
		if a.Offset%(1<<alignShift) != 0 || a.Type != macho.TypeExec {
			t.Fatalf("slice %d: %+v", i, a.FatArchHeader)
		}
		want, _ := os.ReadFile(thin[1-i]) // amd64 first
		got := make([]byte, a.Size)
		f, _ := os.Open(out)
		f.ReadAt(got, int64(a.Offset))
		f.Close()
		if string(got) != string(want) {
			t.Fatalf("slice %d differs from its thin file", i)
		}
	}
	if err := join(out, []string{thin[0], thin[0]}); err == nil {
		t.Fatal("two slices for one CPU were taken")
	}
	if err := join(thin[0], thin); err == nil {
		t.Fatal("an input was taken as the output")
	}
	if b, _ := os.ReadFile(thin[0]); len(b) == 0 {
		t.Fatal("the input was truncated")
	}
	if err := join(out, []string{thin[0], "main.go"}); err == nil {
		t.Fatal("a file that is not Mach-O was taken")
	}
}
