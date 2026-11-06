package app

import (
	"sort"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// ImportOptions are the Import VM dialog's choices on a connection: the pool (or standalone server) itself and, in a
// pool, each of its servers as the home server; and the SRs that can take the VM's disks on the chosen one
// (XvaRules.Targets: StorageRules.NewDiskTargets).
type ImportOptions struct {
	Servers   []ServerPick  `json:"servers"` // the first one ("" uuid) is the pool or standalone server itself
	SRs       []SRUse       `json:"srs"`
	DefaultSR string        `json:"default_sr"`
	Networks  []NetworkPick `json:"networks"` // for an appliance's interfaces (OvfRules.Networks)
}

// applianceNetworks is OvfRules.Networks: the visible networks that are not bond members nor SR-IOV, that the home
// server sees (every server when none), auto-plug ones first, then by name.
func applianceNetworks(s *xenapi.Session, home xenapi.HostRef, hosts map[xenapi.HostRef]xenapi.HostRecord) ([]NetworkPick, error) {
	nets, err := xenapi.Network.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	pifs, err := xenapi.PIF.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	sees := func(n xenapi.NetworkRecord, h xenapi.HostRef) bool { // Host.CanSeeNetwork
		if len(n.PIFs) == 0 {
			return true
		}
		for _, p := range n.PIFs {
			if pifs[p].Host == h {
				return true
			}
		}
		return false
	}
	type pick struct {
		NetworkPick
		auto bool
	}
	var out []pick
	for _, n := range nets {
		ok := networkShown(n, pifs)
		for _, p := range n.PIFs {
			ok = ok && len(pifs[p].SriovLogicalPIFOf) == 0 // Network.IsSriov
		}
		if home != "" {
			ok = ok && sees(n, home)
		} else {
			for h := range hosts {
				ok = ok && sees(n, h)
			}
		}
		if ok {
			out = append(out, pick{NetworkPick{UUID: n.UUID, Name: n.NameLabel}, n.OtherConfig["automatic"] != "false"})
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].auto != out[j].auto {
			return out[i].auto
		}
		return out[i].Name < out[j].Name
	})
	picks := make([]NetworkPick, len(out))
	for i, p := range out {
		picks[i] = p.NetworkPick
	}
	return picks, nil
}

// ReadImportOptions reads them for the home server hostUUID ("" for none).
func ReadImportOptions(c *xapi.Connection, hostUUID string) (ImportOptions, error) {
	s := c.Session
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return ImportOptions{}, xapi.Sanitize(err)
	}
	out := ImportOptions{Servers: []ServerPick{{UUID: "", Name: PoolName(c), Live: true}}}
	var home xenapi.HostRef
	if len(hosts) > 1 { // Helpers.GetPool: the servers of a pool
		var picks []ServerPick
		for r, h := range hosts {
			live, _ := xenapi.HostMetrics.GetLive(s, h.Metrics)
			picks = append(picks, ServerPick{UUID: h.UUID, Name: h.NameLabel, Live: live})
			if h.UUID == hostUUID {
				home = r
			}
		}
		sort.Slice(picks, func(i, j int) bool { return picks[i].Name < picks[j].Name })
		out.Servers = append(out.Servers, picks...)
	}
	disks, err := newDiskSRs(s, home)
	if err != nil {
		return ImportOptions{}, err
	}
	if out.Networks, err = applianceNetworks(s, home, hosts); err != nil {
		return ImportOptions{}, err
	}
	out.SRs = []SRUse{} // XvaRules.Targets: SR.CanFitDisks with no disks, i.e. not full
	for _, u := range disks.SRs {
		if ref, err := xenapi.SR.GetByUUID(s, u.UUID); err == nil {
			if sr, err := xenapi.SR.GetRecord(s, ref); err == nil && !rules.SRFull(sr.Type, sr.PhysicalSize, sr.PhysicalUtilisation) {
				out.SRs = append(out.SRs, u)
				if u.UUID == disks.DefaultSR {
					out.DefaultSR = u.UUID
				}
			}
		}
	}
	return out, nil
}
