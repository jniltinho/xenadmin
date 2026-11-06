package app

import (
	"testing"

	xenapi "go/xenapi"
)

// Two servers' parts of a view: value groups of the same name join, object groups (here two networks named "lan")
// stay apart; groups come before objects, both by name.
func TestMergeViews(t *testing.T) {
	a := &ViewNode{}
	add(a, []group{{"web", "tag", ""}}, Object{Kind: "vms", UUID: "1", Name: "b-vm"}, "s1")
	add(a, []group{{"db", "tag", ""}}, Object{Kind: "vms", UUID: "2", Name: "x"}, "s1")
	add(a, []group{{"lan", "net", "s1/n1"}}, Object{Kind: "vms", UUID: "1", Name: "b-vm"}, "s1")
	b := &ViewNode{}
	add(b, []group{{"web", "tag", ""}}, Object{Kind: "vms", UUID: "3", Name: "a-vm"}, "s2")
	add(b, []group{{"lan", "net", "s2/n9"}}, Object{Kind: "vms", UUID: "3", Name: "a-vm"}, "s2")
	add(b, nil, Object{Kind: "hosts", UUID: "4", Name: "h"}, "s2")
	got := MergeViews(a.Children, b.Children)
	if len(got) != 5 || got[0].Name != "db" || got[1].Name != "lan" || got[2].Name != "lan" || got[3].Name != "web" || got[4].Object == nil {
		t.Fatalf("%d nodes: %+v", len(got), got)
	}
	web := got[3].Children
	if len(web) != 2 || web[0].Name != "a-vm" || web[0].Server != "s2" || web[1].Server != "s1" {
		t.Fatalf("web: %+v %+v", web[0], web[1])
	}
}

func TestViewValues(t *testing.T) {
	if got := osName(xenapi.VMGuestMetricsRecord{OsVersion: map[string]string{"name": "Microsoft Windows Server 2022|C:\\Windows"}}); got != "Windows Server 2022" {
		t.Errorf("os %q", got)
	}
	if osName(xenapi.VMGuestMetricsRecord{}) != "Unknown" {
		t.Error("unknown os")
	}
	bad, _ := customFieldValue(customField{"d", "Date"}, "not a date")
	str, _ := customFieldValue(customField{"s", "String"}, "x")
	_, k1 := customFieldValue(customField{"d", "Date"}, "2026-01-02 03:04:05Z")
	_, k2 := customFieldValue(customField{"d", "Date"}, "2026-01-02 03:04:59Z")
	_, k3 := customFieldValue(customField{"d", "Date"}, "2026-01-02T03:04:05Z")
	if bad != "" || str != "x" || k1 == k2 || k1 != k3 {
		t.Error("custom field values")
	}
	if p := folderPath("/a//b/"); len(p) != 2 || p[1].name != "b" {
		t.Errorf("%v", p)
	}
}
