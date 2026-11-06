package app

import (
	"sort"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// NetworkPick is a network the Add Virtual Interface dialog offers.
type NetworkPick struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// VIFNetworks is NetworkRules.Networks: the visible networks that are not bond members, SR-IOV ones only when the VM
// recommends them. shortcut: sorted by name_label, not Network.Name()'s "Network n" for the default names.
func VIFNetworks(c *xapi.Connection, vmUUID string) ([]NetworkPick, error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	recs, err := xenapi.VM.GetRecommendations(s, ref)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	sriov := rules.Restriction(recs, "allow-network-sriov", "value", 0) != 0 // VM.HasSriovRecommendation
	nets, err := xenapi.Network.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	pifs, err := xenapi.PIF.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	out := []NetworkPick{}
	for _, n := range nets {
		isSriov := false // Network.IsSriov
		for _, p := range n.PIFs {
			isSriov = isSriov || len(pifs[p].SriovLogicalPIFOf) > 0
		}
		if networkShown(n, pifs) && (!isSriov || sriov) {
			out = append(out, NetworkPick{UUID: n.UUID, Name: n.NameLabel})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// DuplicateMAC is NetworkRules.DuplicateMac on one connection: the name and location of a real VM, other than the
// interface vifUUID's ("" for a new one), that already has this MAC ("" when none).
func DuplicateMAC(c *xapi.Connection, vifUUID, mac string) (string, error) {
	if mac == "" {
		return "", nil
	}
	s := c.Session
	vifs, err := xenapi.VIF.GetAllRecords(s)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	for _, v := range vifs {
		if v.UUID == vifUUID && v.MAC == mac {
			return "", nil // the interface keeps its own MAC
		}
	}
	for _, v := range vifs {
		if v.UUID == vifUUID || v.MAC != mac {
			continue
		}
		if vm, err := xenapi.VM.GetRecord(s, v.VM); err == nil && !vm.IsATemplate && !vm.IsASnapshot && !vm.IsControlDomain {
			return vm.NameLabel + " " + location(s), nil // VM.NameWithLocation
		}
	}
	return "", nil
}
