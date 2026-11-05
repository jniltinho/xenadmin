package app

import (
	"strings"
	"testing"

	xenapi "go/xenapi"
)

func TestCannotFitDisks(t *testing.T) {
	gb := 1 << 30
	sr := xenapi.SRRecord{Type: "nfs", PhysicalSize: 100 * gb, PhysicalUtilisation: 90 * gb}
	disk := []xenapi.VDIRecord{{VirtualSize: 20 * gb, PhysicalUtilisation: 2 * gb}}
	if r := cannotFitDisks(sr, nil, disk); !strings.Contains(r, "required when only") {
		t.Errorf("thick: %q", r)
	}
	if r := cannotFitDisks(sr, []string{"THIN_PROVISIONING"}, disk); r != "" {
		t.Errorf("thin: %q", r)
	}
	if r := cannotFitDisks(sr, nil, []xenapi.VDIRecord{{VirtualSize: 3 << 40}}); !strings.Contains(r, "greater than") {
		t.Errorf("max size: %q", r)
	}
	if r := cannotFitDisks(xenapi.SRRecord{Type: "nfs", PhysicalSize: gb, PhysicalUtilisation: gb}, nil, nil); r != "This storage repository is full" {
		t.Errorf("full: %q", r)
	}
}
