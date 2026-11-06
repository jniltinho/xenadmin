package app

import (
	"encoding/xml"
	"strings"
	"testing"
)

// A search as SearchMarshalling writes it: running VMs tagged "web" or named "db*", grouped by power state then tag.
const savedSearch = `<Search xmlns="http://www.citrix.com/XenCenter/XenSearch/schema" uuid="u1" name="Web VMs" major_version="2" minor_version="0" show_expanded="yes">
<Query><QueryScope><VM /><Server /></QueryScope>
<GroupQuery type="And">
  <EnumPropertyQuery property="power_state" equals="True" query="Running" />
  <GroupQuery type="Or">
    <StringListContainsQuery property="tags" contains="True" query="web" />
    <StringPropertyQuery property="label" casesensitive="False" query="DB" type="startswith" />
  </GroupQuery>
</GroupQuery></Query>
<PropertyGrouping property="power_state"><PropertyGrouping property="tags" /></PropertyGrouping>
<columns><column name="name" width="200" /></columns>
</Search>`

func TestParseSearch(t *testing.T) {
	s, err := parseSearch(savedSearch)
	if err != nil || s.name != "Web VMs" || !s.types["VM"] || !s.types["Server"] || s.types["Pool"] || len(s.groupings) != 2 {
		t.Fatalf("%+v %v", s, err)
	}
	vm := func(name, power string, tags ...string) facts {
		return facts{values: map[string]string{"label": name, "power_state": power, "type": "VM"}, lists: map[string][]string{"tags": tags}}
	}
	for _, c := range []struct {
		f    facts
		want tri
	}{
		{vm("web1", "Running", "web"), yes}, {vm("db1", "Running"), yes}, {vm("web2", "Halted", "web"), no}, {vm("x", "Running", "other"), no},
		{facts{values: map[string]string{"label": "host", "type": "Server"}, lists: map[string][]string{"tags": nil}}, no}, // no power state: null
	} {
		if got := match(*s.filter, c.f); got != c.want {
			t.Errorf("%v: %v", c.f.values, got)
		}
	}
	paths := groupPaths(s.groupings, vm("web1", "Running", "web", "prod"))
	if len(paths) != 2 || paths[0][0].name != "Running" || paths[1][1].name != "prod" {
		t.Fatalf("%+v", paths)
	}
	for _, bad := range []string{strings.Replace(savedSearch, "EnumPropertyQuery", "IPAddressQuery", 2), strings.Replace(savedSearch, `property="tags" /`, `property="ip_address" /`, 1)} {
		if err := SearchSupported(bad); err == nil {
			t.Error("an unsupported search runs")
		}
	}
}

func TestSearchSemantics(t *testing.T) {
	tags := []xmlNode{{XMLName: xml.Name{Local: "PropertyGrouping"}, Attrs: []xml.Attr{{Name: xml.Name{Local: "property"}, Value: "tags"}}},
		{XMLName: xml.Name{Local: "PropertyGrouping"}, Attrs: []xml.Attr{{Name: xml.Name{Local: "property"}, Value: "power_state"}}}}
	untagged := facts{values: map[string]string{"power_state": "Running"}}
	if p := groupPaths(tags, untagged); len(p) != 1 || len(p[0]) != 1 || p[0][0].name != "Running" {
		t.Errorf("an empty level must not stop the next one: %+v", p)
	}
	byPool := []xmlNode{{XMLName: xml.Name{Local: "XenModelObjectPropertyGrouping"}, Attrs: []xml.Attr{{Name: xml.Name{Local: "property"}, Value: "pool"}}}}
	pool := facts{groups: map[string][]group{"pool": {{"lab", "", "s/p"}}}, self: "s/p"}
	if p := groupPaths(byPool, pool); len(p) != 1 || p[0][0].key != "s/p" {
		t.Errorf("%+v", p)
	}
	sr := facts{values: map[string]string{"type": "RemoteSR"}}
	q := xmlNode{XMLName: xml.Name{Local: "EnumPropertyQuery"}, Attrs: []xml.Attr{{Name: xml.Name{Local: "property"}, Value: "type"},
		{Name: xml.Name{Local: "query"}, Value: "SR"}, {Name: xml.Name{Local: "equals"}, Value: "yes"}}}
	if match(q, sr) != yes {
		t.Error("the SR alias")
	}
	empty := `<Search name="x" major_version="2"><Query><QueryScope /></Query></Search>`
	if s, err := parseSearch(empty); err != nil || len(s.types) != 0 {
		t.Errorf("an empty scope wants nothing: %+v %v", s, err)
	}
	for _, bad := range []string{`<Search name="x" major_version="2"><Query><QueryScope><Folder /></QueryScope></Query></Search>`,
		`<Search name="x" major_version="2"><Query><QueryScope><VM /></QueryScope><ListEmptyQuery property="networks" empty="no" /></Query></Search>`,
		`<Search name="x" major_version="2"><Query><QueryScope><VM /></QueryScope><GroupQuery type="Xor" /></Query></Search>`} {
		if SearchSupported(bad) == nil {
			t.Errorf("runs %s", bad)
		}
	}
}
