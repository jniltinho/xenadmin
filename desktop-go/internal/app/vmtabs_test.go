package app

import (
	"testing"

	xenapi "go/xenapi"
)

func TestDetailsFormats(t *testing.T) {
	for in, want := range map[int]string{512: "512 B", 10 << 30: "10.0 GB", 1536 << 20: "1.5 GB"} {
		if got := diskSize(in); got != want {
			t.Errorf("diskSize(%d) = %q", in, got)
		}
	}
	if ioNice(xenapi.VBDRecord{}) != 0 || ioNice(xenapi.VBDRecord{QosAlgorithmParams: map[string]string{"class": "7"}}) != 7 {
		t.Error("ioNice")
	}
}
