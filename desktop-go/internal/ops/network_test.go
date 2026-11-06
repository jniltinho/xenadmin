package ops

import (
	"slices"
	"strings"
	"testing"

	xenapi "go/xenapi"
)

// sriovNet: eth0 has an SR-IOV network (logical PIF L, VLAN 5 on it) and VLAN 7 of another network; eth1 has neither.
func sriovNet() *netState {
	return &netState{coordinator: "H",
		hosts: map[xenapi.HostRef]xenapi.HostRecord{"H": {UUID: "h", SoftwareVersion: map[string]string{"network_backend": "openvswitch"}}},
		pifs: map[xenapi.PIFRef]xenapi.PIFRecord{
			"P0": {UUID: "p0", Host: "H", Device: "eth0", VLAN: -1, Managed: true, MTU: 9000, SriovPhysicalPIFOf: []xenapi.NetworkSriovRef{"NS"}},
			"L":  {UUID: "l", Host: "H", Device: "eth0", VLAN: -1, Managed: true, SriovLogicalPIFOf: []xenapi.NetworkSriovRef{"NS"}},
			"V5": {UUID: "v5", Host: "H", Device: "eth0", VLAN: 5, Managed: true, VLANMasterOf: "VL5"},
			"V7": {UUID: "v7", Host: "H", Device: "eth0", VLAN: 7, Managed: true, VLANMasterOf: "VL7"},
			"P1": {UUID: "p1", Host: "H", Device: "eth1", VLAN: -1, Managed: true, MTU: 1500},
		},
		vlans: map[xenapi.VLANRef]xenapi.VLANRecord{"VL5": {TaggedPIF: "L"}, "VL7": {TaggedPIF: "P0"}},
	}
}

func TestVLANsInUseOnSriovNetwork(t *testing.T) { // NetWDetails.GetVLANList
	d := sriovNet()
	others, sriov := d.vlansInUse("eth0", false), d.vlansInUse("eth0", true)
	slices.Sort(others)
	slices.Sort(sriov)
	if !slices.Equal(others, []int{-1, 7}) || !slices.Equal(sriov, []int{-1, 5}) {
		t.Errorf("others %v, sriov %v", others, sriov)
	}
	if got := d.vlansInUse("eth1", true); len(got) != 0 {
		t.Errorf("eth1 on SR-IOV: %v", got)
	}
}

func TestSriovLogicalPIF(t *testing.T) {
	d := sriovNet()
	if got := d.sriovLogicalPIF(d.pifs["P0"]); got != "L" {
		t.Errorf("P0: %q", got)
	}
	if got := d.sriovLogicalPIF(d.pifs["P1"]); got != "" { // no SR-IOV network
		t.Errorf("P1: %q", got)
	}
	delete(d.pifs, "L")
	if got := d.sriovLogicalPIF(d.pifs["P0"]); got != "" { // its logical PIF is gone
		t.Errorf("P0 without L: %q", got)
	}
}

func TestNewNetworkOnSriov(t *testing.T) { // NetWDetails' checks, NewNetworkWizard.CreateNonBonded's PIF
	d := sriovNet()
	spec := func(nic, vlan string, sriov bool) NewNetwork {
		return NewNetwork{External: true, Name: "n", Host: "h", NIC: nic, VLAN: vlan, MTU: "1500", Sriov: sriov}
	}
	for _, c := range []struct {
		spec NewNetwork
		nic  xenapi.PIFRef
		vlan int
		err  string
	}{
		{spec("p0", "5", true), "", 0, "already in use"},  // used on the SR-IOV network
		{spec("p0", "7", true), "L", 7, ""},               // used by another network only: free on the SR-IOV one
		{spec("p0", "5", false), "P0", 5, ""},             // used on the SR-IOV network only
		{spec("p0", "7", false), "", 0, "already in use"}, // used by another network
		{spec("p0", "0", true), "L", 0, ""},               // VLAN 0 on vSwitch
		{spec("p0", "4095", true), "", 0, "range"},
		{spec("p0", "x", true), "", 0, "Invalid"},
		{spec("p1", "5", false), "P1", 5, ""},               // no SR-IOV network: the NIC itself
		{spec("p1", "5", true), "", 0, "no SR-IOV network"}, // the box is not offered there: a stale choice stops
	} {
		_, nic, vlan, calls, err := d.newNetwork(c.spec)
		switch {
		case c.err != "":
			if err == nil || !strings.Contains(err.Error(), c.err) {
				t.Errorf("%+v: err %v, want %q", c.spec, err, c.err)
			}
		case err != nil || nic != c.nic || vlan != c.vlan || !slices.Equal(calls, []string{"network.create", "pool.create_VLAN_from_PIF"}):
			t.Errorf("%+v: nic %q vlan %d calls %v err %v", c.spec, nic, vlan, calls, err)
		}
	}
	d.hosts["H"] = xenapi.HostRecord{UUID: "h"} // Linux bridge: no VLAN 0
	if _, _, _, _, err := d.newNetwork(spec("p0", "0", true)); err == nil || !strings.Contains(err.Error(), "1 to 4094") {
		t.Errorf("VLAN 0 on a bridge: %v", err)
	}
	delete(d.pifs, "L")
	if _, _, _, _, err := d.newNetwork(spec("p0", "8", true)); err == nil {
		t.Error("SR-IOV network gone: no error")
	}
	if _, nic, _, calls, err := d.newNetwork(NewNetwork{Name: "p", Sriov: true}); err != nil || nic != "" || !slices.Equal(calls, []string{"network.create"}) {
		t.Errorf("private: nic %q calls %v err %v", nic, calls, err)
	}
}
