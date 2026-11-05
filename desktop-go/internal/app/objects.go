// Package app holds the use cases shared by the desktop app, the CLI and the MCP server (no UI code).
package app

import (
	"fmt"
	"sort"
	"strings"

	"go/xenapi"

	"xcpng-center/internal/xapi"
)

// Kinds that ListObjects accepts (the MCP list_objects kinds of the C# app).
var Kinds = []string{"pools", "hosts", "vms", "templates", "srs", "networks"}

// Object is one row of a listing.
type Object struct {
	Kind        string `json:"kind"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	PowerState  string `json:"power_state,omitempty"`
	Host        string `json:"host,omitempty"`
}

// ListObjects lists one kind of object of a connection, sorted by name (as the C# tree and list_objects): VMs are the
// real VMs (no templates, snapshots or control domains), templates the user and default templates.
func ListObjects(c *xapi.Connection, kind string) ([]Object, error) {
	s := c.Session
	var out []Object
	switch kind {
	case "pools":
		recs, err := xenapi.Pool.GetAllRecords(s)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		for _, r := range recs {
			out = append(out, Object{Kind: kind, UUID: r.UUID, Name: r.NameLabel, Description: r.NameDescription})
		}
	case "hosts":
		recs, err := xenapi.Host.GetAllRecords(s)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		for _, r := range recs {
			out = append(out, Object{Kind: kind, UUID: r.UUID, Name: r.NameLabel, Description: r.NameDescription})
		}
	case "vms", "templates":
		vms, err := xenapi.VM.GetAllRecords(s)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		hosts, err := xenapi.Host.GetAllRecords(s)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		for _, r := range vms {
			if r.IsControlDomain || r.IsASnapshot || r.IsATemplate != (kind == "templates") {
				continue
			}
			o := Object{Kind: kind, UUID: r.UUID, Name: r.NameLabel, Description: r.NameDescription}
			if kind == "vms" {
				o.PowerState = string(r.PowerState)
				if h, ok := hosts[r.ResidentOn]; ok {
					o.Host = h.NameLabel
				}
			}
			out = append(out, o)
		}
	case "srs":
		recs, err := xenapi.SR.GetAllRecords(s)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		for _, r := range recs {
			out = append(out, Object{Kind: kind, UUID: r.UUID, Name: r.NameLabel, Description: r.NameDescription})
		}
	case "networks":
		recs, err := xenapi.Network.GetAllRecords(s)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		for _, r := range recs {
			out = append(out, Object{Kind: kind, UUID: r.UUID, Name: r.NameLabel, Description: r.NameDescription})
		}
	default:
		return nil, fmt.Errorf("unknown kind %q (one of %s)", kind, strings.Join(Kinds, ", "))
	}
	SortObjects(out)
	return out, nil
}

// SortObjects orders by name (case-insensitive), then UUID, so pages never repeat or skip equal names.
func SortObjects(objs []Object) {
	sort.SliceStable(objs, func(i, j int) bool {
		a, b := strings.ToLower(objs[i].Name), strings.ToLower(objs[j].Name)
		if a != b {
			return a < b
		}
		return objs[i].UUID < objs[j].UUID
	})
}

// KindFor maps the C# MCP types (pool, host, vm, template, sr, network) to the listing kinds.
func KindFor(t string) (string, bool) {
	for _, k := range Kinds {
		if k == t || k == t+"s" {
			return k, true
		}
	}
	return "", false
}
