package app

import (
	"strings"
	"testing"
)

func TestBuildNewSearch(t *testing.T) {
	if _, err := BuildSearch(SearchEdit{Name: " ", Types: []string{"VM"}}); err == nil || err.Error() != "Enter a name for the search." {
		t.Fatal(err)
	}
	if _, err := BuildSearch(SearchEdit{Name: "x", Types: []string{"Folder"}}); err == nil || !strings.HasPrefix(err.Error(), "Select at least one type") {
		t.Fatal(err)
	}
	text, err := BuildSearch(SearchEdit{Name: " Web ", Types: []string{"VM", "Server"}, Text: "web  db", Grouping: 5})
	if err != nil {
		t.Fatal(err)
	}
	// SearchMarshalling's order: the scope in ObjectTypes order, the filter, the grouping, sorting
	want := `<Search uuid="` + SearchUUID(text) + `" name="Web" major_version="2" minor_version="0"><Query><QueryScope><Server></Server><VM></VM></QueryScope>` +
		`<GroupQuery type="Or"><StringPropertyQuery property="label" casesensitive="no" query="web" type="contains"></StringPropertyQuery>` +
		`<StringPropertyQuery property="description" casesensitive="no" query="web" type="contains"></StringPropertyQuery>` +
		`<StringPropertyQuery property="label" casesensitive="no" query="db" type="contains"></StringPropertyQuery>` +
		`<StringPropertyQuery property="description" casesensitive="no" query="db" type="contains"></StringPropertyQuery></GroupQuery></Query>` +
		`<PropertyGrouping property="power_state"></PropertyGrouping><sorting></sorting></Search>`
	if text != want || len(SearchUUID(text)) != 36 {
		t.Fatalf("%s", text)
	}
	s, err := parseSearch(text)
	if err != nil || s.name != "Web" || !s.types["VM"] || !s.types["Server"] || len(s.types) != 2 || len(s.groupings) != 1 {
		t.Fatalf("%+v %v", s, err)
	}
	vm := func(name, desc string) facts {
		return facts{values: map[string]string{"label": name, "description": desc, "type": "VM"}}
	}
	if match(*s.filter, vm("my-WEB-1", "")) != yes || match(*s.filter, vm("x", "the DB")) != yes || match(*s.filter, vm("x", "y")) != no {
		t.Error("the text filter")
	}
	f, err := EditSearchForm(text)
	if err != nil || f.Name != "Web" || f.Grouping != 5 || !f.HasFilter || strings.Join(f.Types, ",") != "Server,VM" {
		t.Fatalf("%+v %v", f, err)
	}
}

// An edited search keeps what the editor does not show: its filter (AND-ed with the text), the types not offered,
// a nested grouping while unchanged, its columns; a saved one keeps its uuid.
func TestEditSearchKeepsTheRest(t *testing.T) {
	edited := strings.Replace(savedSearch, "<VM /><Server />", "<VM /><Server /><Folder /><DockerContainer />", 1)
	f, err := EditSearchForm(edited)
	if err != nil || f.Name != "Web VMs" || f.Grouping != 0 || !f.HasFilter || strings.Join(f.Types, ",") != "VM,Server" {
		t.Fatalf("%+v %v", f, err)
	}
	text, err := BuildSearch(SearchEdit{Name: "Web 2", Types: []string{"VM"}, Text: "prod", Edited: edited, Saved: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{`uuid="u1"`, `<QueryScope><VM></VM><Folder></Folder><DockerContainer></DockerContainer></QueryScope><GroupQuery type="And"><GroupQuery type="And"><EnumPropertyQuery`,
		`query="prod" type="contains"></StringPropertyQuery></GroupQuery></GroupQuery></Query><PropertyGrouping property="power_state"><PropertyGrouping property="tags">`,
		`<columns><column name="name" width="200"></column></columns>`, `show_expanded="yes"`} {
		if !strings.Contains(text, part) {
			t.Errorf("lost %s in %s", part, text)
		}
	}
	if strings.Contains(text, "xmlns") || strings.Contains(text, "<Server>") {
		t.Errorf("%s", text)
	}
	// a new grouping replaces the nested one; an unsaved search gets a new uuid
	text, err = BuildSearch(SearchEdit{Name: "Web 2", Types: []string{"VM"}, Grouping: 7, GroupingChanged: true, Edited: edited})
	if err != nil || SearchUUID(text) == "u1" || !strings.Contains(text, `</Query><FolderGrouping></FolderGrouping><columns>`) || strings.Contains(text, "PropertyGrouping") {
		t.Fatalf("%s %v", text, err)
	}
	text, _ = BuildSearch(SearchEdit{Name: "Web 2", Types: []string{"VM"}, Grouping: 0, GroupingChanged: true, Edited: edited})
	if strings.Contains(text, "Grouping") {
		t.Errorf("None keeps a grouping: %s", text)
	}
	if _, err := EditSearchForm(`<Search name="old" major_version="1"><DummyQuery /></Search>`); err == nil {
		t.Error("a version 1 search is edited")
	}
	if SearchUUID(`<Bad uuid="u1"/>`) != "" { // SaveSearch never writes what is not a search
		t.Error("the uuid of a non-search")
	}
}

func TestSearchFileRoundTrip(t *testing.T) {
	text, _ := BuildSearch(SearchEdit{Name: "Tagged & <odd>", Types: []string{"Pool", "VDI"}, Grouping: 6})
	file, err := SearchFile(text)
	if err != nil || !strings.HasPrefix(file, `<?xml version="1.0" encoding="utf-8"?>`+"\n"+`<!DOCTYPE xensearch PUBLIC "-//XENSEARCH//DTD XENSEARCH 1//EN" "xensearch-1.dtd">`+"\n"+
		`<Searches xmlns="http://www.citrix.com/XenCenter/XenSearch/schema">`+"\n"+`    <Search uuid=`) {
		t.Fatalf("%s %v", file, err)
	}
	got, err := ImportSearches([]byte(file))
	if err != nil || len(got) != 1 || SearchUUID(got[0]) == SearchUUID(text) || SearchName(got[0]) != "Tagged & <odd>" || strings.Contains(got[0], "xmlns") {
		t.Fatalf("%v %v", got, err)
	}
	s, err := parseSearch(got[0])
	if err != nil || !s.types["Pool"] || !s.types["VDI"] || len(s.groupings) != 1 || s.filter != nil {
		t.Fatalf("%+v %v", s, err)
	}
	// a file the C# app wrote: the Search element in the Searches namespace (or with xmlns="")
	cs := `<?xml version="1.0" encoding="utf-8"?><!DOCTYPE xensearch PUBLIC "-//XENSEARCH//DTD XENSEARCH 1//EN" "xensearch-1.dtd">
<Searches xmlns="http://www.citrix.com/XenCenter/XenSearch/schema">
    <Search uuid="a" name="C#" major_version="2" minor_version="0" xmlns="">
        <Query><QueryScope><VM /></QueryScope></Query>
        <sorting />
    </Search>
</Searches>`
	if got, err := ImportSearches([]byte(cs)); err != nil || len(got) != 1 || SearchSupported(got[0]) != nil || SearchName(got[0]) != "C#" {
		t.Fatalf("%v %v", got, err)
	}
	for _, bad := range []string{"not xml", `<Searches xmlns="x"></Searches>`, `<Searches><Search name="v1" major_version="1"><DummyQuery /></Search></Searches>`} {
		if _, err := ImportSearches([]byte(bad)); err == nil {
			t.Errorf("imported %q", bad)
		}
	}
	if SearchSupported(OverviewSearch()) != nil || SearchName(OverviewSearch()) != "" {
		t.Error("the overview")
	}
	if f, _ := EditSearchForm(OverviewSearch()); len(f.Types) != len(objectTypes) {
		t.Errorf("%+v", f)
	}
}
