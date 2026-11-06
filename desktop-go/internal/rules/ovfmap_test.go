package rules

import (
	"strings"
	"testing"
)

func TestCheckApplianceMappings(t *testing.T) {
	a := &Appliance{VMs: []ApplianceVM{
		{SysID: "vm1", Name: "one", Disks: []ApplianceDisk{{Name: "d0", Capacity: 4 << 30}, {Name: "d1", Capacity: 8 << 30}}, MACs: []string{"", ""}},
		{SysID: "vm2", Name: "two", Disks: []ApplianceDisk{{Name: "d0", Capacity: 1 << 30}}},
	}}
	st := &ApplianceState{
		Hosts: []string{"h1", "h2"},
		SRs: map[string]ApplianceSR{
			"nfs":    {Name: "nfs", Hosts: []string{"h1", "h2"}, Free: 100 << 30},
			"local1": {Name: "local1", Hosts: []string{"h1"}, Free: 100 << 30},
			"local2": {Name: "local2", Hosts: []string{"h2"}, Free: 10 << 30},
			"full":   {Name: "full", Hosts: []string{"h1", "h2"}, Full: true},
			"thin":   {Name: "thin", Hosts: []string{"h1", "h2"}, Thin: true},
		},
		Networks: map[string][]string{"priv": nil, "eth0": {"h1", "h2"}, "eth1-h2": {"h2"}},
		HostName: func(u string) string { return "host " + u },
	}
	ok := func() []ApplianceMapping { return DefaultApplianceMappings(a, "nfs", "eth0") }
	if err := CheckApplianceMappings(a, ok(), st); err != nil {
		t.Fatal(err)
	}
	// the lab's case: home h2, disk 0 shared, disk 1 on h2's local SR, interfaces on two networks
	m := ok()
	m[0].Host, m[0].SRs[1], m[0].Networks[1] = "h2", "local2", "eth1-h2"
	a.VMs[0].Disks[1].Capacity = 9 << 30
	if err := CheckApplianceMappings(a, m, st); err != nil {
		t.Fatal(err)
	}
	a.VMs[0].Disks[1].Capacity = 8 << 30
	// the pool (no server): local SRs only when one server sees all the VM's disks
	m = ok()
	m[0].SRs[1] = "local1"
	if err := CheckApplianceMappings(a, m, st); err != nil {
		t.Fatal(err)
	}
	m[0].SRs[0] = "local2"
	expect := func(m []ApplianceMapping, want string) {
		t.Helper()
		err := CheckApplianceMappings(a, m, st)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("got %v, want %q", err, want)
		}
	}
	expect(m, "No server sees all the storage")
	for _, c := range []struct {
		edit func(m []ApplianceMapping)
		want string
	}{
		{func(m []ApplianceMapping) { m[1].SysID = "other" }, "do not match"},
		{func(m []ApplianceMapping) { m[0].SRs = m[0].SRs[:1] }, "do not match"},
		{func(m []ApplianceMapping) { m[0].Networks = append(m[0].Networks, "eth0") }, "do not match"},
		{func(m []ApplianceMapping) { m[0].Host = "h9" }, "not in the pool"},
		{func(m []ApplianceMapping) { m[0].SRs[0] = "nope" }, "No SR available"},
		{func(m []ApplianceMapping) { m[0].Host, m[0].SRs[0] = "h2", "local1" }, "cannot be seen from host h2"},
		{func(m []ApplianceMapping) { m[1].SRs[0] = "full" }, "is full"},
		{func(m []ApplianceMapping) { m[0].Networks[0] = "nope" }, "No available networks"},
		{func(m []ApplianceMapping) { m[0].Networks[0] = "eth1-h2" }, "No available networks"},                  // not every server
		{func(m []ApplianceMapping) { m[0].Host, m[0].Networks[0] = "h1", "eth1-h2" }, "No available networks"}, // not h1
		{func(m []ApplianceMapping) { m[0].Host, m[0].SRs[0], m[0].SRs[1] = "h2", "local2", "local2" }, "not enough available space on SR local2"},
	} {
		m := ok()
		c.edit(m)
		expect(m, c.want)
	}
	expect(ok()[:1], "do not match")
	// a thin SR is not counted; a private network is seen by every server
	m = ok()
	a.VMs[0].Disks[1].Capacity = 1 << 50
	m[0].SRs[1], m[0].Networks[0] = "thin", "priv"
	expect(m, "cannot create a disk greater than")
	st.SRs["thin"] = ApplianceSR{Name: "thin", Hosts: []string{"h1", "h2"}, Thin: true, LargeVDI: true}
	if err := CheckApplianceMappings(a, m, st); err != nil {
		t.Fatal(err)
	}
}
