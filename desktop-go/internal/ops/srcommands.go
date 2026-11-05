package ops

import (
	"errors"
	"fmt"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// SRStates reads what the SR commands need for every SR of the pool (the C# app reads it from its cache).
func SRStates(s *xenapi.Session) (map[xenapi.SRRef]rules.SRState, map[xenapi.SRRef]xenapi.SRRecord, error) {
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return nil, nil, xapi.Sanitize(err)
	}
	pbds, err := xenapi.PBD.GetAllRecords(s)
	if err != nil {
		return nil, nil, xapi.Sanitize(err)
	}
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return nil, nil, xapi.Sanitize(err)
	}
	sms, err := xenapi.SM.GetAllRecords(s)
	if err != nil {
		return nil, nil, xapi.Sanitize(err)
	}
	hosts, err := xenapi.Host.GetAll(s)
	if err != nil {
		return nil, nil, xapi.Sanitize(err)
	}
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return nil, nil, xapi.Sanitize(err)
	}
	attached := map[xenapi.SRRef]bool{}
	unplugged := map[xenapi.SRRef]bool{}
	for _, p := range pbds {
		attached[p.SR] = attached[p.SR] || p.CurrentlyAttached
		unplugged[p.SR] = unplugged[p.SR] || !p.CurrentlyAttached
	}
	vdiCreate, multipath := map[string]bool{}, map[string]bool{}
	for _, sm := range sms {
		vdiCreate[sm.Type] = contains(sm.Capabilities, "VDI_CREATE")
		_, multipath[sm.Type] = sm.Features["SR_MULTIPATH"]
	}
	// !SR.MultipathAOK for a LUN per SR: a multipathed PBD whose "mpath-..." status has fewer paths than the maximum
	// shortcut: LUN-per-VDI SRs (sm_config "LUNperVDI"/"scsi-...") are not checked; only HBA/iSCSI raw LUNs have them
	degraded := map[xenapi.SRRef]bool{}
	for _, p := range pbds {
		if p.OtherConfig["multipathed"] != "true" {
			continue
		}
		for k, v := range p.OtherConfig {
			if strings.HasPrefix(k, "mpath") {
				degraded[p.SR] = degraded[p.SR] || rules.MultipathDown(v)
				break
			}
		}
	}
	vbds, err := xenapi.VBD.GetAllRecords(s)
	if err != nil {
		return nil, nil, xapi.Sanitize(err)
	}
	vdis, err := xenapi.VDI.GetAllRecords(s)
	if err != nil {
		return nil, nil, xapi.Sanitize(err)
	}
	busy := map[xenapi.SRRef]bool{} // SR.HasRunningVMs
	for _, b := range vbds {
		vm, running := vms[b.VM]
		vdi, ok := vdis[b.VDI]
		if !running || vm.PowerState != xenapi.VMPowerStateRunning || !ok || vdi.Type == xenapi.VdiTypeMetadata && vm.IsControlDomain { // PR-1223
			continue
		}
		busy[vdi.SR] = true
	}
	out := map[xenapi.SRRef]rules.SRState{}
	for ref, r := range srs {
		st := rules.SRState{Type: r.Type, Detached: !attached[ref], HasPBDs: len(r.PBDs) > 0, RunningVMs: busy[ref],
			Shared: r.Shared, Hosts: len(hosts), Broken: unplugged[ref], Hidden: r.OtherConfig["HideFromXenCenter"] == "true",
			// SR.SupportsVdiCreate: never an ISO library or a memory SR, whatever the driver says (CA-40119)
			VDICreate: vdiCreate[r.Type] && r.ContentType != "iso" && r.Type != "tmpfs",
			Degraded:  (multipath[r.Type] || r.SmConfig["multipathable"] == "true") && degraded[ref]}
		if r.Shared && len(r.PBDs) != len(hosts) || !r.Shared && len(r.PBDs) != 1 { // SR.IsBroken
			st.Broken = true
		}
		for _, op := range r.AllowedOperations {
			st.Allowed = append(st.Allowed, string(op))
		}
		for _, p := range pools {
			st.Default = st.Default || p.DefaultSR == ref
		}
		out[ref] = st
	}
	return out, srs, nil
}

// SRCommand runs an SR command of the Storage menu, re-checked on the server first: repair (SrRepairAction),
// set_default (SetSrAsDefaultAction),
// detach (DetachSrAction: unplug the PBDs, the coordinator's last), forget (ForgetSrAction) or destroy (DestroySrAction).
func SRCommand(c *xapi.Connection, uuid, command string) error {
	s := c.Session
	sr, err := xenapi.SR.GetByUUID(s, uuid)
	if err != nil {
		return xapi.Sanitize(err)
	}
	states, _, err := SRStates(s)
	if err != nil {
		return err
	}
	if !contains(rules.SRCommands(states[sr]), command) {
		return fmt.Errorf("%s cannot run on this SR now", command)
	}
	switch command {
	case "repair":
		return repairSR(s, sr)
	case "set_default":
		return setDefaultSR(s, sr)
	case "detach":
		return unplugPBDs(s, sr)
	case "forget":
		return xapi.Sanitize(wait(s)(xenapi.SR.AsyncForget(s, sr)))
	case "destroy":
		return xapi.Sanitize(wait(s)(xenapi.SR.AsyncDestroy(s, sr)))
	}
	return errors.New("unknown SR command")
}

// unplugPBDs is DetachSrAction.UnplugPBDs: the coordinator's PBD last (CA-176935, CA-173497).
func unplugPBDs(s *xenapi.Session, sr xenapi.SRRef) error {
	pbds, err := xenapi.SR.GetPBDs(s, sr)
	if err != nil {
		return xapi.Sanitize(err)
	}
	coordinator, err := coordinatorOf(s)
	if err != nil {
		return err
	}
	var last []xenapi.PBDRef
	for _, p := range pbds {
		if h, err := xenapi.PBD.GetHost(s, p); err == nil && h == coordinator {
			last = append(last, p)
			continue
		}
		if err := wait(s)(xenapi.PBD.AsyncUnplug(s, p)); err != nil {
			return xapi.Sanitize(err)
		}
	}
	for _, p := range last {
		if err := wait(s)(xenapi.PBD.AsyncUnplug(s, p)); err != nil {
			return xapi.Sanitize(err)
		}
	}
	return nil
}

// repairSR is SrRepairAction: a GFS2 SR's cluster resynced first; then for the coordinator first, then every other host (a local SR only on its own host), a
// shared SR gets the PBD a host lacks (a copy of its first PBD's device_config), and an unplugged PBD is plugged; every
// host is tried, and the last failure is returned.
func repairSR(s *xenapi.Session, sr xenapi.SRRef) error {
	rec, err := xenapi.SR.GetRecord(s, sr)
	if err != nil {
		return xapi.Sanitize(err)
	}
	coordinator, err := coordinatorOf(s)
	if err != nil {
		return err
	}
	all, err := xenapi.Host.GetAll(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	pbdOf := map[xenapi.HostRef]xenapi.PBDRecord{}
	refOf := map[xenapi.HostRef]xenapi.PBDRef{}
	for _, p := range rec.PBDs {
		r, err := xenapi.PBD.GetRecord(s, p)
		if err != nil {
			return xapi.Sanitize(err)
		}
		pbdOf[r.Host], refOf[r.Host] = r, p
	}
	if rec.Type == "gfs2" { // the cluster first, best effort, as the C# app
		if clusters, err := xenapi.Cluster.GetAll(s); err == nil && len(clusters) > 0 {
			_ = xenapi.Cluster.PoolResync(s, clusters[0])
		}
	}
	hosts := []xenapi.HostRef{coordinator}
	for _, h := range all {
		if h != coordinator {
			hosts = append(hosts, h)
		}
	}
	var failure error
	for _, h := range hosts {
		if _, ok := pbdOf[h]; !ok && !rec.Shared {
			continue // a local SR is repaired on its own host only
		}
		if _, ok := pbdOf[h]; !ok && len(rec.PBDs) > 0 {
			first, err := xenapi.PBD.GetRecord(s, rec.PBDs[0])
			if err != nil {
				failure = xapi.Sanitize(err)
				continue
			}
			ref, err := waitRef(s)(xenapi.PBD.AsyncCreate(s, xenapi.PBDRecord{Host: h, SR: sr, DeviceConfig: first.DeviceConfig, OtherConfig: map[string]string{}}))
			if err != nil {
				failure = xapi.Sanitize(err)
				continue
			}
			refOf[h], pbdOf[h] = xenapi.PBDRef(ref), xenapi.PBDRecord{}
		}
		if pbd, ok := refOf[h]; ok && !pbdOf[h].CurrentlyAttached {
			if err := wait(s)(xenapi.PBD.AsyncPlug(s, pbd)); err != nil {
				failure = xapi.Sanitize(err)
			}
		}
	}
	return failure
}
