package app

import (
	"slices"
	"sort"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/xapi"
)

// VappVM is a VM of a vApp, or one that may join it, with its start order and delay.
type VappVM struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	PowerState string `json:"power_state"`
	Vapp       string `json:"vapp"` // the vApp it is in ("" for none)
	Order      int    `json:"order"`
	Delay      int    `json:"delay"`
}

// Vapp is a vApp of Manage vApps.
type Vapp struct {
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	VMs         []VappVM `json:"vms"`
	CanStart    bool     `json:"can_start"`    // VappRules.CanStart
	CanShutDown bool     `json:"can_shutdown"` // VappRules.CanShutDown
}

// Vapps are the vApps of a connection (VappRules.Appliances, by name) and the VMs that may be in one
// (VappRules.Candidates: the shown real VMs, by name).
type Vapps struct {
	Pool       string   `json:"pool"`
	IsPool     bool     `json:"is_pool"`
	Vapps      []Vapp   `json:"vapps"`
	Candidates []VappVM `json:"candidates"`
}

// ReadVapps reads them.
func ReadVapps(c *xapi.Connection) (Vapps, error) {
	s := c.Session
	apps, err := xenapi.VMAppliance.GetAllRecords(s)
	if err != nil {
		return Vapps{}, xapi.Sanitize(err)
	}
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return Vapps{}, xapi.Sanitize(err)
	}
	hosts, _ := xenapi.Host.GetAll(s)
	pools, _ := xenapi.Pool.GetAllRecords(s)
	out := Vapps{Pool: PoolName(c), Vapps: []Vapp{}, Candidates: []VappVM{}}
	for _, p := range pools {
		out.IsPool = p.NameLabel != "" || len(hosts) > 1
	}
	uuidOf := map[xenapi.VMApplianceRef]string{}
	for r, a := range apps {
		uuidOf[r] = a.UUID
	}
	row := func(v xenapi.VMRecord) VappVM {
		return VappVM{UUID: v.UUID, Name: v.NameLabel, PowerState: string(v.PowerState), Vapp: uuidOf[v.Appliance], Order: v.Order, Delay: v.StartDelay}
	}
	byName := func(l []VappVM) {
		sort.Slice(l, func(i, j int) bool { return strings.ToLower(l[i].Name) < strings.ToLower(l[j].Name) })
	}
	for _, v := range vms {
		if !v.IsATemplate && !v.IsASnapshot && !v.IsControlDomain && vmShown(v) {
			out.Candidates = append(out.Candidates, row(v))
		}
	}
	byName(out.Candidates)
	for _, a := range apps {
		va := Vapp{UUID: a.UUID, Name: a.NameLabel, Description: a.NameDescription, VMs: []VappVM{},
			CanStart:    slices.Contains(a.AllowedOperations, xenapi.VMApplianceOperationStart),
			CanShutDown: slices.Contains(a.AllowedOperations, xenapi.VMApplianceOperationCleanShutdown) || slices.Contains(a.AllowedOperations, xenapi.VMApplianceOperationHardShutdown)}
		for _, r := range a.VMs {
			if v, ok := vms[r]; ok {
				va.VMs = append(va.VMs, row(v))
			}
		}
		byName(va.VMs)
		out.Vapps = append(out.Vapps, va)
	}
	sort.Slice(out.Vapps, func(i, j int) bool { return strings.ToLower(out.Vapps[i].Name) < strings.ToLower(out.Vapps[j].Name) })
	return out, nil
}

// vmShown is VM.Show(false): not internal, not a GUI temporary object, not hidden.
func vmShown(v xenapi.VMRecord) bool {
	return vmShownHidden(v) && v.OtherConfig["HideFromXenCenter"] != "true"
}
