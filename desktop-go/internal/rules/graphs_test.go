package rules

import (
	"maps"
	"reflect"
	"testing"
)

const graphVM = "4a1c2f0e-1111-2222-3333-444455556666"

// a pool's gui_config as the Windows app writes it (SaveDataSourceStateAction): two graphs of a VM, a graph of another
// VM, a colour and the custom fields; index 3 is past the gap at 2, so LoadGraphs never reads it
var csharpGui = map[string]string{
	"XenCenter.GraphLayout.0.vm." + graphVM:        "cpu0,cpu1",
	"XenCenter.GraphName.0.vm." + graphVM:          "CPU Performance",
	"XenCenter.GraphLayout.1.vm." + graphVM:        "memory,memory_internal_free",
	"XenCenter.GraphLayout.3.vm." + graphVM:        "vif_0_tx",
	"XenCenter.GraphLayout.0.vm.other-vm":          "cpu0",
	"XenCenter.GraphName.0.host." + graphVM:        "a server graph that happens to carry the uuid",
	"XenCenter.DataSource.vm." + graphVM + ".cpu0": "-16750849",
	"XenCenter.CustomFields":                       "<CustomFieldDefinitions />",
}

func TestParseGraphLayout(t *testing.T) {
	got := ParseGraphLayout(csharpGui, "vm", graphVM)
	want := []DesignedGraph{{Name: "CPU Performance", Sources: []string{"cpu0", "cpu1"}}, {Sources: []string{"memory", "memory_internal_free"}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v", got)
	}
	if ParseGraphLayout(csharpGui, "host", "nothing") != nil {
		t.Fatal("no layout is the defaults")
	}
	// C# saves a default graph of no devices (a VM without VIFs) as ""
	empty := ParseGraphLayout(map[string]string{"XenCenter.GraphLayout.0.vm.x": ""}, "vm", "x")
	if v, ok := GraphLayoutConfig(nil, "vm", "x", empty)["XenCenter.GraphLayout.0.vm.x"]; len(empty) != 1 || len(empty[0].Sources) != 0 || !ok || v != "" {
		t.Fatalf("%+v", empty)
	}
}

func TestGraphLayoutConfig(t *testing.T) {
	got := GraphLayoutConfig(csharpGui, "vm", graphVM, []DesignedGraph{{Name: "Disks", Sources: []string{"vbd_xvda_read", "vbd_xvda_write"}}, {Sources: []string{"cpu0"}}})
	want := map[string]string{
		"XenCenter.GraphLayout.0.vm." + graphVM:        "vbd_xvda_read,vbd_xvda_write",
		"XenCenter.GraphName.0.vm." + graphVM:          "Disks",
		"XenCenter.GraphLayout.1.vm." + graphVM:        "cpu0",
		"XenCenter.GraphLayout.0.vm.other-vm":          "cpu0",
		"XenCenter.DataSource.vm." + graphVM + ".cpu0": "-16750849",
		"XenCenter.CustomFields":                       "<CustomFieldDefinitions />",
	}
	if !maps.Equal(got, want) {
		t.Fatalf("got %v", got)
	}
	// Restore Default Graphs: no keys of the object left
	if d := GraphLayoutConfig(csharpGui, "vm", graphVM, nil); len(d) != 3 {
		t.Fatalf("defaults: %v", d)
	}
}

func TestGraphColours(t *testing.T) {
	if v, err := ColourValue("#0066ff"); err != nil || v != "-16750849" { // Color.FromArgb(0,102,255).ToArgb()
		t.Fatalf("%q %v", v, err)
	}
	if c, ok := ParseColour("-16750849"); !ok || c != "#0066ff" {
		t.Fatalf("%q", c)
	}
	if _, err := ColourValue("red"); err == nil {
		t.Fatal("not a colour")
	}
	if _, ok := ParseColour("x"); ok {
		t.Fatal("not a colour")
	}
}

func TestFriendlyDataSource(t *testing.T) {
	find := func(what string, ids ...string) string {
		switch what + ":" + ids[0] {
		case "vif:0":
			return "Pool-wide network"
		case "vbd:xvda":
			return "0"
		case "sr:abcd1234":
			return "Local storage"
		}
		return ""
	}
	for name, want := range map[string]string{
		"cpu3": "CPU 3", "memory_internal_free": "Used Memory", "vif_0_tx": "Pool-wide network Send", "vif_9_tx": "vif_9_tx",
		"pif_eth1_rx": "NIC 1 Receive", "vbd_xvda_read": "Disk 0 Read", "read_abcd1234": "'Local storage' Read Throughput",
		"hostload": "Host CPU Load", "no_such_source": "no_such_source",
	} {
		if got := FriendlyDataSource(name, find); got != want {
			t.Errorf("%s: %q, want %q", name, got, want)
		}
	}
	if DataSourceCategory("cpu0") != "CPU" {
		t.Error(DataSourceCategory("cpu0"))
	}
	if got := OfferedDataSources([]string{"memory", "io_throughput_read_abcd1234", "read_abcd1234", "cpu0"}); !reflect.DeepEqual(got, []string{"read_abcd1234", "cpu0"}) {
		t.Error(got)
	}
}
