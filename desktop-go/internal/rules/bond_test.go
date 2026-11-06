package rules

import (
	"strings"
	"testing"
)

func TestBondRules(t *testing.T) {
	if BondSizeLimit(true) != 4 || BondSizeLimit(false) != 2 || !BondValidSelection(2, 2) || BondValidSelection(1, 4) || BondValidSelection(3, 2) {
		t.Error("bond size")
	}
	if BondMaxMTU(nil) != 9216 || BondMaxMTU([]int{9000, 1500}) != 1500 || BondMaxMTU([]int{65000}) != 9216 {
		t.Error("bond MTU")
	}
	if n := BondName([]string{"eth2", "eth1"}); n != "Bond 1+2" {
		t.Errorf("BondName %q", n)
	}
	if BondCreateWarning(false, false, false, true, "p") != nil {
		t.Error("no warning without management interfaces")
	}
	if w := BondCreateWarning(true, true, false, false, "p"); !w.Block {
		t.Error("primary and secondary block")
	}
	if w := BondCreateWarning(true, false, false, true, "p"); !w.Block || !strings.Contains(w.Text, "HA") {
		t.Errorf("HA blocks: %+v", w)
	}
	if w := BondCreateWarning(true, false, false, false, "p"); w.Block || w.OK != "Create bond anyway" || !strings.Contains(w.Text, "XCP-ng Center") {
		t.Errorf("primary: %+v", w)
	}
	if w := BondCreateWarning(false, true, true, false, "p"); w.Block || !strings.Contains(w.Text, "clustering") {
		t.Errorf("clustering: %+v", w)
	}
	if w := BondDeleteWarning("Bond 1+2", false, false, true, "p"); w.Block || w.Text != "Are you sure that you want to delete Bond 1+2?" || w.OK != "OK" {
		t.Errorf("plain delete: %+v", w)
	}
	if w := BondDeleteWarning("Bond 0+1", true, false, true, "p"); !w.Block {
		t.Error("HA blocks deleting the management bond")
	}
	if w := BondDeleteWarning("Bond 0+1", true, true, false, "p"); w.Block || !strings.Contains(w.Text, "secondary") || w.OK != "Delete bond anyway" {
		t.Errorf("both: %+v", w)
	}
}

func TestVLANAndMTU(t *testing.T) {
	if VLANFor(5, []int{-1, 3}) != 5 || VLANFor(0, []int{-1, 0, 1, 2}) != 3 {
		t.Error("VLANFor")
	}
	if MinVLAN(true) != 0 || MinVLAN(false) != 1 {
		t.Error("MinVLAN")
	}
	if v, _, ok := ValidateVLAN(" 7 ", 1, []int{-1}); !ok || v != 7 {
		t.Error("VLAN 7")
	}
	if _, msg, ok := ValidateVLAN("0", 0, nil); !ok || !strings.Contains(msg, "VLAN 0") {
		t.Error("VLAN 0 is allowed with a message")
	}
	for _, bad := range []string{"x", "0", "4095", "3"} {
		if _, _, ok := ValidateVLAN(bad, 1, []int{3}); ok {
			t.Errorf("VLAN %q", bad)
		}
	}
	if m, msg, ok := ValidateMTU("9000", 1500, 9216); !ok || m != 9000 || msg != "Allowed MTU range: 1500 to 9216" {
		t.Errorf("MTU 9000: %q", msg)
	}
	if _, msg, ok := ValidateMTU("1600", 1500, 1500); ok || msg != "Allowed MTU value: 1500" {
		t.Errorf("fixed MTU: %q", msg)
	}
	if NICMaxMTU(65535) != 9216 || NICMaxMTU(1500) != 1500 {
		t.Error("NICMaxMTU")
	}
}

func TestNetworkButtons(t *testing.T) {
	if !CanEditNetwork(NetworkInfo{}) || CanEditNetwork(NetworkInfo{BondMember: true}) || CanEditNetwork(NetworkInfo{Busy: true}) {
		t.Error("CanEditNetwork")
	}
	if ok, tip := CanRemoveNetwork(NetworkInfo{Physical: true}); ok || tip == "" {
		t.Error("a physical NIC's network cannot be removed")
	}
	if ok, _ := CanRemoveNetwork(NetworkInfo{Bond: true}); !ok {
		t.Error("a bond network is removed by deleting the bond")
	}
	if ok, _ := CanRemoveNetwork(NetworkInfo{GuestInstal: true}); ok {
		t.Error("guest installer network")
	}
}

func TestNetworkMTU(t *testing.T) {
	if m := NetworkMTU(1500, true, false, false, "", false); m.Shown {
		t.Error("no MTU on a tunnel network")
	}
	if m := NetworkMTU(1500, false, true, false, "", false); !m.Shown || m.Enabled {
		t.Error("a private network's MTU is fixed")
	}
	if m := NetworkMTU(1500, false, false, true, "Management", true); m.Enabled || !strings.Contains(m.Note, "'Management'") {
		t.Errorf("management first: %+v", m)
	}
	if m := NetworkMTU(1500, false, false, false, "", true); m.Enabled || !strings.Contains(m.Note, "without tools") {
		t.Errorf("VMs without tools: %+v", m)
	}
	if m := NetworkMTU(1400, false, false, false, "", false); !m.Enabled || m.Min != 1400 || m.Max != 9216 {
		t.Errorf("a CLI value widens the range: %+v", m)
	}
}

func TestNetworkNICs(t *testing.T) {
	if n := NetworkNICs(true, false, false, "", true); !n.Shown || !n.Enabled || n.Note != "" {
		t.Errorf("a private network moves to a NIC even with VMs on it: %+v", n)
	}
	if n := NetworkNICs(false, false, false, "", false); n.Shown || n.Enabled {
		t.Errorf("a NIC, bond or tunnel network: %+v", n)
	}
	if n := NetworkNICs(false, true, true, "Storage", true); n.Enabled || !strings.Contains(n.Note, "management interface 'Storage'") {
		t.Errorf("management first: %+v", n)
	}
	if n := NetworkNICs(false, true, false, "", true); n.Enabled || !strings.Contains(n.Note, "active VMs attached") {
		t.Errorf("VMs attached: %+v", n)
	}
	if n := NetworkNICs(false, true, false, "", false); !n.Shown || !n.Enabled {
		t.Errorf("a VLAN: %+v", n)
	}
}

func TestBondChanges(t *testing.T) {
	for _, c := range []struct {
		cur, curHash, mode, hash string
		lacp, setMode, setHash   bool
		err                      bool
	}{
		{"balance-slb", "", "balance-slb", "src_mac", true, false, false, false},  // nothing: the hashing only goes with LACP
		{"balance-slb", "", "active-backup", "", true, true, false, false},        // a new mode
		{"balance-slb", "", "lacp", "tcpudp_ports", true, true, true, false},      // to LACP, with its hashing
		{"lacp", "src_mac", "lacp", "tcpudp_ports", true, false, true, false},     // the hashing only
		{"lacp", "", "lacp", "src_mac", true, false, true, false},                 // an unknown hashing is set (as the radio shows src_mac)
		{"lacp", "src_mac", "active-backup", "src_mac", true, true, false, false}, // from LACP
		{"lacp", "", "lacp", "", false, false, false, false},                      // no vSwitch: LACP stays, hidden
		{"balance-slb", "", "lacp", "src_mac", false, false, false, true},         // no vSwitch: no LACP
		{"balance-slb", "", "lacp", "", true, false, false, true},                 // LACP needs a hashing
		{"balance-slb", "", "802.3ad", "", true, false, false, true},
	} {
		m, h, err := BondChanges(c.cur, c.curHash, c.mode, c.hash, c.lacp)
		if m != c.setMode || h != c.setHash || (err != nil) != c.err {
			t.Errorf("%+v: %v %v %v", c, m, h, err)
		}
	}
}
