// Package app holds the use cases shared by the desktop app, the CLI and the MCP server (no UI code).
package app

import (
	"fmt"
	"sort"
	"strings"

	"go/xenapi"

	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// Kinds that ListObjects accepts (the MCP list_objects kinds of the C# app).
var Kinds = []string{"pools", "hosts", "vms", "templates", "srs", "networks"}

// Object is one row of a listing.
type Object struct {
	Kind        string   `json:"kind"`
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	PowerState  string   `json:"power_state,omitempty"`
	Host        string   `json:"host,omitempty"`
	HostUUID    string   `json:"host_uuid,omitempty"` // the resident host of a VM (names may repeat)
	Allowed     []string `json:"allowed,omitempty"`   // a VM's allowed_operations, an SR's commands (what the menus enable)
	HAProtected bool     `json:"ha_protected,omitempty"`
	HAEnabled   bool     `json:"ha_enabled,omitempty"`  // a pool with HA on (the Pool > High Availability items)
	Maintenance bool     `json:"maintenance,omitempty"` // a host that is disabled (HostRules.InMaintenance)
	Coordinator bool     `json:"coordinator,omitempty"`
	Live        bool     `json:"live,omitempty"` // a host that answers
	Icon        string   `json:"icon,omitempty"` // an SR's tree icon (Images.GetIconFor)
	// ControlDomain is a host's dom0 VM (its Console tab shows the server's console)
	ControlDomain string `json:"control_domain,omitempty"`
	// Resident is the number of VMs a host runs, its control domain included (HostRules.CanRemoveFromPool: at most 1)
	Resident int `json:"resident,omitempty"`
	// UserTemplate is a template made by a user (ResourceTree.IsUserTemplate: not a default one, shown in the tree)
	UserTemplate bool `json:"user_template,omitempty"`
	// GPUTab, USBTab: the pool or server shows a GPU tab, the server a USB tab (HardwareRules.ShowGpuTab, ShowUsbTab)
	GPUTab bool `json:"gpu_tab,omitempty"`
	USBTab bool `json:"usb_tab,omitempty"`
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
		hosts, err := xenapi.Host.GetAllRecords(s)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		for _, r := range recs {
			out = append(out, Object{Kind: kind, UUID: r.UUID, Name: r.NameLabel, Description: r.NameDescription, HAEnabled: r.HaEnabled,
				GPUTab: !ops.FeatureForbidden(hosts, "restrict_gpu")})
		}
	case "hosts":
		recs, err := xenapi.Host.GetAllRecords(s)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		pools, err := xenapi.Pool.GetAllRecords(s)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		for ref, r := range recs {
			live, err := xenapi.HostMetrics.GetLive(s, r.Metrics)
			if err != nil {
				live = false // no metrics yet: the server is starting
			}
			o := Object{Kind: kind, UUID: r.UUID, Name: r.NameLabel, Description: r.NameDescription, Maintenance: !r.Enabled, Live: live,
				Resident: len(r.ResidentVMs)}
			if dom0, err := xenapi.VM.GetUUID(s, r.ControlDomain); err == nil {
				o.ControlDomain = dom0
			}
			for _, p := range pools {
				o.Coordinator = o.Coordinator || p.Master == ref
			}
			o.GPUTab = live && !ops.FeatureForbidden(recs, "restrict_gpu")
			o.USBTab = live && len(r.PUSBs) > 0 && !ops.FeatureForbidden(recs, "restrict_usb_passthrough")
			out = append(out, o)
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
		pools, err := xenapi.Pool.GetAllRecords(s)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		out = vmObjects(kind, vms, hosts, pools)
	case "srs":
		states, recs, err := ops.SRStates(s)
		if err != nil {
			return nil, err
		}
		for ref, r := range recs {
			out = append(out, Object{Kind: kind, UUID: r.UUID, Name: r.NameLabel, Description: r.NameDescription,
				Allowed: rules.SRCommands(states[ref]), Icon: rules.SRIcon(states[ref])})
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

// vmObjects are the real VMs ("vms": no templates, snapshots or control domains) or the templates of the records,
// with what the window shows of a VM (state, resident server, allowed operations, HA protection).
func vmObjects(kind string, vms map[xenapi.VMRef]xenapi.VMRecord, hosts map[xenapi.HostRef]xenapi.HostRecord, pools map[xenapi.PoolRef]xenapi.PoolRecord) []Object {
	ha := false
	for _, p := range pools {
		ha = p.HaEnabled
	}
	var out []Object
	for _, r := range vms {
		if r.IsControlDomain || r.IsASnapshot || r.IsATemplate != (kind == "templates") {
			continue
		}
		o := Object{Kind: kind, UUID: r.UUID, Name: r.NameLabel, Description: r.NameDescription}
		if kind == "templates" {
			_, def := r.OtherConfig["default_template"]
			_, internal := r.OtherConfig["xensource_internal"]
			o.UserTemplate = !def && !internal && r.OtherConfig["HideFromXenCenter"] != "true" && !strings.HasPrefix(r.NameLabel, "__gui__")
			for _, a := range r.AllowedOperations {
				o.Allowed = append(o.Allowed, string(a))
			}
		}
		if kind == "vms" {
			o.PowerState = string(r.PowerState)
			if h, ok := hosts[r.ResidentOn]; ok {
				o.Host, o.HostUUID = h.NameLabel, h.UUID
			}
			for _, a := range r.AllowedOperations {
				o.Allowed = append(o.Allowed, string(a))
			}
			o.HAProtected = ha && r.HaRestartPriority != "" // VM.HAIsProtected
		}
		out = append(out, o)
	}
	return out
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
