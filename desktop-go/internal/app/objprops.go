package app

import (
	"fmt"
	"sort"
	"strings"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// ObjectProperties are the General and Custom Fields pages of an object's Properties (GeneralEditPage,
// CustomFieldsPage), with the folders and tags of its connection for the Change Folder and Edit Tags dialogs.
type ObjectProperties struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Folder      string       `json:"folder"`
	Tags        []string     `json:"tags"`
	AllTags     []string     `json:"all_tags"` // Tags.GetAllTags of the connection
	Folders     []string     `json:"folders"`  // the folder paths in use, with the pool's empty folders
	IQN         *string      `json:"iqn,omitempty"`
	Fields      []FieldValue `json:"fields"`
	// a server's Log Destination page (LogDestinationEditPage: the remote syslog server, "" for local only) and
	// Multipathing page (HostMultipathPage: editable in maintenance mode or when the server is down)
	LogDestination    *string `json:"log_destination,omitempty"`
	Multipath         *bool   `json:"multipath,omitempty"`
	MultipathEditable bool    `json:"multipath_editable,omitempty"`
}

// ReadObjectProperties reads an object ("vms", "templates", "hosts", "pools", "srs", "networks") for its Properties.
func ReadObjectProperties(c *xapi.Connection, kind, uuid string) (ObjectProperties, error) {
	s := c.Session
	out := ObjectProperties{Tags: []string{}, AllTags: []string{}, Folders: []string{}, Fields: []FieldValue{}}
	tags, folders := map[string]bool{}, map[string]bool{}
	note := func(t []string, other map[string]string) {
		for _, x := range t {
			tags[x] = true
		}
		if f := other["folder"]; f != "" {
			folders[f] = true
		}
	}
	var other map[string]string
	found := false
	pick := func(u string, name, desc string, t []string, o map[string]string) {
		if u == uuid {
			out.Name, out.Description, out.Tags, other, found = name, desc, append([]string{}, t...), o, true
		}
	}
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	var fieldDefs []rules.CustomField
	for _, p := range pools {
		note(p.Tags, p.OtherConfig)
		for _, f := range strings.Split(p.OtherConfig["EMPTY_FOLDERS"], ";") { // Folders.EMPTY_FOLDERS
			if f != "" {
				folders[f] = true
			}
		}
		fieldDefs = rules.ParseCustomFields(p.GuiConfig[rules.CustomFieldsKey])
		if kind == "pools" {
			pick(p.UUID, p.NameLabel, p.NameDescription, p.Tags, p.OtherConfig)
		}
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	for _, h := range hosts {
		note(h.Tags, h.OtherConfig)
		if kind == "hosts" && h.UUID == uuid {
			pick(h.UUID, h.NameLabel, h.NameDescription, h.Tags, h.OtherConfig)
			iqn := h.IscsiIqn // Host.GetIscsiIqn
			out.IQN = &iqn
			logTo, mp := h.Logging["syslog_destination"], h.Multipathing // Host.GetSysLogDestination, MultipathEnabled
			out.LogDestination, out.Multipath = &logTo, &mp
			live, err := xenapi.HostMetrics.GetLive(s, h.Metrics)
			out.MultipathEditable = h.OtherConfig["MAINTENANCE_MODE"] == "true" || err == nil && !live // unknown metrics: not down
		}
	}
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	for _, v := range vms {
		note(v.Tags, v.OtherConfig)
		if kind == "vms" || kind == "templates" {
			pick(v.UUID, v.NameLabel, v.NameDescription, v.Tags, v.OtherConfig)
		}
	}
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	for _, r := range srs {
		note(r.Tags, r.OtherConfig)
		if kind == "srs" {
			pick(r.UUID, r.NameLabel, r.NameDescription, r.Tags, r.OtherConfig)
		}
	}
	nets, err := xenapi.Network.GetAllRecords(s)
	if err != nil {
		return out, xapi.Sanitize(err)
	}
	for _, n := range nets {
		note(n.Tags, n.OtherConfig)
		if kind == "networks" {
			pick(n.UUID, n.NameLabel, n.NameDescription, n.Tags, n.OtherConfig)
		}
	}
	if !found {
		return out, fmt.Errorf("the object %s is gone", uuid)
	}
	out.Folder = other["folder"]
	for t := range tags {
		out.AllTags = append(out.AllTags, t)
	}
	for f := range folders {
		out.Folders = append(out.Folders, f)
	}
	sort.Strings(out.Tags)
	sort.Strings(out.AllTags)
	sort.Strings(out.Folders)
	for _, f := range fieldDefs {
		v := other[rules.CustomFieldsKey+"."+f.Name]
		if f.Type == "Date" {
			t, ok := rules.ParseFieldDate(v)
			v = ""
			if ok {
				v = t.UTC().Format(time.RFC3339)
			}
		}
		out.Fields = append(out.Fields, FieldValue{f, v})
	}
	sort.Slice(out.Fields, func(i, j int) bool { return out.Fields[i].Name < out.Fields[j].Name })
	return out, nil
}
