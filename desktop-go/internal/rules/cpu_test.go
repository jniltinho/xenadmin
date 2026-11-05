package rules

import "testing"

func TestCPURulesMirrorTheCSharp(t *testing.T) {
	got := ValidTopologies(4, 2)
	if len(got) != 2 || got[0].Sockets != 4 || got[1].Cores != 2 || got[1].Text != "2 sockets with 2 cores per socket" {
		t.Errorf("%+v", got)
	}
	if len(ValidTopologies(32, 1)) != 0 { // 32 sockets of 1 core is over the Xen limit
		t.Error("socket limit")
	}
	if ValidVCPUConfiguration(6, 4) == "" || ValidVCPUConfiguration(8, 4) != "" {
		t.Error("ValidVCPUConfiguration")
	}
	if Priority(256) != 4 || Weight(4) != 256 || Priority(65536) != 8 || Priority(1) != 0 {
		t.Error("priority")
	}
	if TopologyText(1, 1) != "1 socket with 1 core per socket" {
		t.Error(TopologyText(1, 1))
	}
}
