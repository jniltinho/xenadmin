package app

import (
	"crypto/rand"
	"encoding/xml"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
)

// The Search tab's editor (XenAdmin.Desktop.Core SearchRules, SearchEditor): the object types, a full-text filter
// and one grouping, written as SearchMarshalling XML (major version 2) that the C# app loads. An edited or imported
// search keeps everything else it has: its filter (the text filter is AND-ed on top), the types not offered here, a
// nested or unlisted grouping while the choice is not changed, its columns and sorting.
// shortcut: no QueryElement tree of property conditions, as the C# Desktop editor; port QueryElement when users
// need field-by-field conditions.

// SearchChoice is an object type or a grouping the editor offers.
type SearchChoice struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SearchForm is the editor filled from a search.
type SearchForm struct {
	Name      string         `json:"name"`
	Types     []string       `json:"types"`      // the offered types it wants
	Grouping  int            `json:"grouping"`   // index in Groupings, 0 (None) when not listed
	HasFilter bool           `json:"has_filter"` // its conditions are kept
	TypeList  []SearchChoice `json:"type_list"`
	Groupings []SearchChoice `json:"groupings"`
}

// SearchEdit is what the editor asks for (SearchRules.Build).
type SearchEdit struct {
	Name            string   `json:"name"`
	Types           []string `json:"types"`
	Text            string   `json:"text"`
	Grouping        int      `json:"grouping"`
	GroupingChanged bool     `json:"grouping_changed"` // an unchanged choice keeps the edited search's own grouping
	Edited          string   `json:"edited"`           // the edited search's XML ("" for New Search)
	Saved           bool     `json:"saved"`            // the edited search is saved in a pool: it keeps its uuid
}

// scopeOrder is ObjectTypes in enum order, as QueryScope.ToXmlNode writes them.
var scopeOrder = []string{"Pool", "Server", "DisconnectedServer", "VM", "Snapshot", "UserTemplate", "DefaultTemplate",
	"RemoteSR", "LocalSR", "VDI", "Network", "Folder", "Appliance", "DockerContainer"}

// searchGroupings are SearchRules.Groupings: None, then one property (element and property attribute).
var searchGroupings = []struct{ name, element, property string }{
	{messages.Get("NONE"), "", ""},
	{"Type", "PropertyGrouping", "type"},
	{messages.Get("POOL"), "XenModelObjectPropertyGrouping", "pool"},
	{messages.Get("SERVER"), "XenModelObjectPropertyGrouping", "host"},
	{"Operating System", "PropertyGrouping", "os_name"},
	{"Power State", "PropertyGrouping", "power_state"},
	{messages.Get("TAGS"), "PropertyGrouping", "tags"},
	{messages.Get("FOLDERS"), "FolderGrouping", ""},
	{messages.Get("VM_APPLIANCE"), "XenModelObjectPropertyGrouping", "appliance"},
}

// SearchTypes are SearchRules.Types: the object types a search can return, by name.
// shortcut: only the types RunSearch reads (no Disconnected servers, Folders, vApps); an edited search keeps those.
func SearchTypes() []SearchChoice {
	out := make([]SearchChoice, 0, len(objectTypes))
	for t, key := range objectTypes {
		out = append(out, SearchChoice{t, messages.Format(messages.Get(key), "XCP-ng")})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func searchGroupingChoices() []SearchChoice {
	out := make([]SearchChoice, len(searchGroupings))
	for i, g := range searchGroupings {
		out[i] = SearchChoice{fmt.Sprint(i), g.name}
	}
	return out
}

func attr(name, value string) xml.Attr { return xml.Attr{Name: xml.Name{Local: name}, Value: value} }

func element(name string, attrs ...xml.Attr) xmlNode {
	return xmlNode{XMLName: xml.Name{Local: name}, Attrs: attrs}
}

func (n *xmlNode) setAttr(name, value string) {
	for i := range n.Attrs {
		if n.Attrs[i].Name.Local == name {
			n.Attrs[i].Value = value
			return
		}
	}
	n.Attrs = append(n.Attrs, attr(name, value))
}

// clean drops the namespaces of a search read from a file (the Searches element's xmlns) so it is written as the C#
// app writes it in gui_config: a plain <Search> element.
func clean(n *xmlNode) {
	n.XMLName.Space = ""
	n.Attrs = slices.DeleteFunc(n.Attrs, func(a xml.Attr) bool { return a.Name.Space == "xmlns" || a.Name.Local == "xmlns" })
	for i := range n.Children {
		clean(&n.Children[i])
	}
}

// readSearch parses a search's XML as SearchMarshalling.LoadSearch needs it: major version 2, a Query with a scope.
func readSearch(text string) (xmlNode, error) {
	var root xmlNode
	if err := xml.Unmarshal([]byte(text), &root); err != nil {
		return root, err
	}
	clean(&root)
	if err := checkSearch(root); err != nil {
		return root, err
	}
	return root, nil
}

// checkSearch is what SearchMarshalling.LoadSearch requires of a version 2 search.
// shortcut: version 0 and 1 searches (a filter without a Query) are refused, not converted; the filters and
// groupings are passed through unchecked (a hand-edited file the C# loader refuses stays refused there). Check each
// QueryFilter's attributes when the Go side ports them all.
func checkSearch(root xmlNode) error {
	if root.XMLName.Local != "Search" && root.XMLName.Local != "search" || root.attr("major_version") != "2" ||
		len(root.Children) == 0 || root.Children[0].XMLName.Local != "Query" ||
		!slices.ContainsFunc(root.Children[0].Children, func(c xmlNode) bool { return c.XMLName.Local == "QueryScope" }) {
		return errors.New("not a version 2 search")
	}
	return nil
}

func searchXML(root xmlNode) string {
	b, _ := xml.Marshal(root) // xmlNode always marshals
	return string(b)
}

// newUUID is Guid.NewGuid: a random (version 4) uuid.
func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6], b[8] = b[6]&0x0f|0x40, b[8]&0x3f|0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// OverviewSearch is SearchRules.Overview (Search.SearchForAllTypes): every type, no filter, no grouping, no name.
// shortcut: the types RunSearch reads, not AllExcFolders (no disconnected servers here).
func OverviewSearch() string {
	scope := element("QueryScope")
	for _, t := range scopeOrder {
		if _, ok := objectTypes[t]; ok {
			scope.Children = append(scope.Children, element(t))
		}
	}
	return searchXML(newSearch("", scope))
}

func newSearch(name string, scope xmlNode) xmlNode {
	root := element("Search", attr("uuid", newUUID()), attr("name", name), attr("major_version", "2"), attr("minor_version", "0"))
	root.Children = []xmlNode{{XMLName: xml.Name{Local: "Query"}, Children: []xmlNode{scope}}, element("sorting")}
	return root
}

// EditSearchForm fills the editor: from a search's XML, or New Search's defaults ("" text: Search.DefaultObjectTypes
// that are offered, no grouping).
func EditSearchForm(text string) (SearchForm, error) {
	f := SearchForm{Name: "New Search", Types: []string{"Server", "VM", "RemoteSR"}, TypeList: SearchTypes(), Groupings: searchGroupingChoices()}
	if text == "" {
		return f, nil
	}
	root, err := readSearch(text)
	if err != nil {
		return f, err
	}
	f.Name, f.Types = root.attr("name"), nil
	for _, c := range root.Children[0].Children {
		if c.XMLName.Local != "QueryScope" {
			f.HasFilter = true
			continue
		}
		for _, t := range c.Children {
			if _, ok := objectTypes[t.XMLName.Local]; ok && !slices.Contains(f.Types, t.XMLName.Local) {
				f.Types = append(f.Types, t.XMLName.Local)
			}
		}
	}
	// SearchRules.GroupingIndex: a listed grouping with no subgrouping, else None
	for _, g := range root.Children[1:] {
		if g.XMLName.Local == "columns" || g.XMLName.Local == "sorting" {
			continue
		}
		for i, choice := range searchGroupings[1:] {
			if g.XMLName.Local == choice.element && g.attr("property") == choice.property && len(g.Children) == 0 {
				f.Grouping = i + 1
			}
		}
		break
	}
	return f, nil
}

// fullQueryFor is Search.FullQueryFor: each word in the name or the description (contains, any case), any of them.
// shortcut: no IPAddressQuery for a word that is an address (RunSearch does not read addresses yet); port it with
// IPAddressQuery's matching.
func fullQueryFor(text string) xmlNode {
	q := element("GroupQuery", attr("type", "Or"))
	str := func(prop, word string) xmlNode {
		return element("StringPropertyQuery", attr("property", prop), attr("casesensitive", "no"), attr("query", word), attr("type", "contains"))
	}
	for _, w := range strings.Split(text, " ") {
		if w != "" {
			q.Children = append(q.Children, str("label", w), str("description", w))
		}
	}
	if len(q.Children) == 0 {
		q.Children = []xmlNode{str("label", "")}
	}
	return q
}

// BuildSearch is SearchRules.Build: the search the editor describes, as XML, or the validation text as the error.
// The edited search's filter stays (the text filter AND-ed on top), with its other parts; only a saved search keeps
// its uuid (saving replaces it), any other one gets a new one.
func BuildSearch(e SearchEdit) (string, error) {
	types := 0
	for _, t := range e.Types {
		if _, ok := objectTypes[t]; ok {
			types++
		}
	}
	if msg := rules.ValidateSearch(e.Name, types); msg != "" {
		return "", errors.New(msg)
	}
	if e.Grouping < 0 || e.Grouping >= len(searchGroupings) {
		return "", errors.New("unknown grouping")
	}
	root := newSearch("", element("QueryScope"))
	if e.Edited != "" {
		var err error
		if root, err = readSearch(e.Edited); err != nil {
			return "", err
		}
	}
	root.setAttr("name", strings.TrimSpace(e.Name))
	if !e.Saved || root.attr("uuid") == "" {
		root.setAttr("uuid", newUUID())
	}
	root.setAttr("major_version", "2")
	if root.attr("minor_version") == "" {
		root.setAttr("minor_version", "0")
	}

	query := &root.Children[0]
	var scope, filter *xmlNode
	var rest []xmlNode // anything else of the Query, kept
	for i := range query.Children {
		switch c := &query.Children[i]; {
		case c.XMLName.Local == "QueryScope" && scope == nil:
			scope = c
		case c.XMLName.Local != "QueryScope" && filter == nil:
			filter = c
		default:
			rest = append(rest, *c)
		}
	}
	// the scope: the chosen types, and those of the edited search the editor does not offer, in ObjectTypes order
	want := map[string]bool{}
	for _, t := range e.Types {
		if _, ok := objectTypes[t]; ok {
			want[t] = true
		}
	}
	var unknown []xmlNode
	for _, t := range scope.Children {
		if _, offered := objectTypes[t.XMLName.Local]; !offered && t.XMLName.Local != "NONE" && !strings.HasPrefix(t.XMLName.Local, "ALL") {
			if slices.Contains(scopeOrder, t.XMLName.Local) {
				want[t.XMLName.Local] = true
			} else {
				unknown = append(unknown, t)
			}
		}
	}
	newScope := element("QueryScope")
	for _, t := range scopeOrder {
		if want[t] {
			newScope.Children = append(newScope.Children, element(t))
		}
	}
	newScope.Children = append(newScope.Children, unknown...)
	children := []xmlNode{newScope}
	if strings.TrimSpace(e.Text) != "" {
		text := fullQueryFor(strings.TrimSpace(e.Text))
		if filter != nil {
			text = xmlNode{XMLName: xml.Name{Local: "GroupQuery"}, Attrs: []xml.Attr{attr("type", "And")}, Children: []xmlNode{*filter, text}}
		}
		filter = &text
	}
	if filter != nil {
		children = append(children, *filter)
	}
	query.Children = append(children, rest...)

	if e.GroupingChanged || e.Edited == "" {
		kept := []xmlNode{root.Children[0]}
		if g := searchGroupings[e.Grouping]; g.element != "" {
			n := element(g.element)
			if g.property != "" {
				n.Attrs = []xml.Attr{attr("property", g.property)}
			}
			kept = append(kept, n)
		}
		for _, c := range root.Children[1:] {
			if c.XMLName.Local == "columns" || c.XMLName.Local == "sorting" || !isGrouping(c.XMLName.Local) {
				kept = append(kept, c)
			}
		}
		root.Children = kept
	}
	return searchXML(root), nil
}

// isGrouping tells a Grouping element (GroupingTypes.cs) from the other parts of a search.
func isGrouping(name string) bool {
	switch name {
	case "PropertyGrouping", "BoolGrouping", "XenModelObjectPropertyGrouping", "FolderGrouping", "CustomFieldGrouping", "AllCustomFieldsGrouping":
		return true
	}
	return false
}

// SearchName and SearchUUID read a search's name and uuid ("" when the text is not a version 2 search).
func SearchName(text string) string { return searchAttr(text, "name") }
func SearchUUID(text string) string { return searchAttr(text, "uuid") }

func searchAttr(text, name string) string {
	r, err := readSearch(text)
	if err != nil {
		return ""
	}
	return r.attr(name)
}

// searchNamespace is SearchMarshalling's xmlns of the Searches element.
const searchNamespace = "http://www.citrix.com/XenCenter/XenSearch/schema"

// SearchFile is Search.Save(filename) (SearchMarshalling.SearchesToXML): the .xensearch file of one search, with the
// declaration, the doctype and the Searches element the C# app writes.
func SearchFile(text string) (string, error) {
	root, err := readSearch(text)
	if err != nil {
		return "", err
	}
	if root.attr("uuid") == "" {
		root.setAttr("uuid", newUUID())
	}
	doc := xmlNode{XMLName: xml.Name{Local: "Searches"}, Attrs: []xml.Attr{attr("xmlns", searchNamespace)}, Children: []xmlNode{root}}
	b, err := xml.MarshalIndent(doc, "", "    ")
	if err != nil {
		return "", err
	}
	return `<?xml version="1.0" encoding="utf-8"?>` + "\n" +
		`<!DOCTYPE xensearch PUBLIC "-//XENSEARCH//DTD XENSEARCH 1//EN" "xensearch-1.dtd">` + "\n" + string(b) + "\n", nil
}

// ImportSearches is SearchRules.Import (Search.LoadFile): the searches of a .xensearch file, each with a new
// identity (never the original's); an error when the file is not one, or a search of it cannot be loaded.
func ImportSearches(data []byte) ([]string, error) {
	var root xmlNode
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	clean(&root)
	var found []xmlNode
	var walk func(n xmlNode)
	walk = func(n xmlNode) { // SearchMarshalling.LoadSearches: every Search (or old search) element
		if n.XMLName.Local == "Search" || n.XMLName.Local == "search" {
			found = append(found, n)
			return
		}
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(root)
	var out []string
	for _, s := range found {
		if err := checkSearch(s); err != nil {
			return nil, err
		}
		s.XMLName.Local = "Search"
		s.setAttr("uuid", newUUID())
		out = append(out, searchXML(s))
	}
	if len(out) == 0 {
		return nil, errors.New("no search")
	}
	return out, nil
}
