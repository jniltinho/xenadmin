package ops

import (
	"testing"

	xenapi "go/xenapi"
)

func TestEvacuatedVMsByListAndPowerState(t *testing.T) {
	h := &host{ref: "H", rec: xenapi.HostRecord{OtherConfig: map[string]string{
		evacMigrated: "a,b,c,gone", evacHalted: "d",
	}}, vms: map[xenapi.VMRef]xenapi.VMRecord{
		"A": {UUID: "a", PowerState: xenapi.VMPowerStateRunning, ResidentOn: "H2"}, // migrated away: goes back
		"B": {UUID: "b", PowerState: xenapi.VMPowerStateRunning, ResidentOn: "H"},  // already back
		"C": {UUID: "c", PowerState: xenapi.VMPowerStateHalted},                    // saved as migrated, now halted: not this list
		"D": {UUID: "d", PowerState: xenapi.VMPowerStateHalted},                    // halted: started on the server
	}}
	m, ha, su := evacuated(h)
	if len(m) != 1 || m[0] != "A" || len(ha) != 1 || ha[0] != "D" || len(su) != 0 {
		t.Errorf("migrated %v halted %v suspended %v", m, ha, su)
	}
}

func TestResultRef(t *testing.T) {
	for in, want := range map[string]string{"<value>OpaqueRef:ab-1</value>": "OpaqueRef:ab-1", "OpaqueRef:x": "OpaqueRef:x", "": ""} {
		if got := resultRef(in); got != want {
			t.Errorf("resultRef(%q) = %q", in, got)
		}
	}
}
