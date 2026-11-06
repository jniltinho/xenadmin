package ops

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// Disaster recovery (XenAdmin.Desktop.Core/DrRules.cs): DRConfigureCommand's database replication, the pool metadata
// of XenModel/Actions/DR (GetMetadataVDIsAction, VdiOpenDatabaseAction, VdiLoadMetadataAction), DrRules.RecoverAction
// (DrRecoverAction, StartVMsAndAppliancesAction) and the clean-up of a test failover (ShutdownAndDestroyVMsAction,
// ShutdownAndDestroyVmAppliancesAction).

// PoolOf is the pool of the connection (Helpers.GetPoolOfOne), with its number of servers.
func PoolOf(s *xenapi.Session) (xenapi.PoolRef, xenapi.PoolRecord, int, error) { return poolOf(s) }

// ReplicationSR is an SR that supports database replication, and whether it holds this pool's metadata (DR on) or
// any pool's.
type ReplicationSR struct {
	Ref         xenapi.SRRef
	Record      xenapi.SRRecord
	Enabled     bool // SrRow.IsDrEnabled: a metadata VDI of this pool
	HasMetadata bool // DRFailoverWizardStoragePage: a metadata VDI of any pool
}

// ReplicationSRs is DrRules.ReplicationSrs and StorageChoices: the SRs where SR.assert_supports_database_replication
// passes, by name.
func ReplicationSRs(s *xenapi.Session) ([]ReplicationSR, error) {
	poolRef, _, _, err := poolOf(s)
	if err != nil {
		return nil, err
	}
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	vdis, err := xenapi.VDI.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	out := []ReplicationSR{}
	for ref, sr := range srs {
		if err := xenapi.SR.AssertSupportsDatabaseReplication(s, ref); err != nil {
			var apiErr *xenapi.APIError
			if errors.As(err, &apiErr) { // SR.SupportsDatabaseReplication: a failure is a no
				continue
			}
			return nil, xapi.Sanitize(err)
		}
		r := ReplicationSR{Ref: ref, Record: sr}
		for _, v := range vdis {
			if v.SR == ref && v.Type == xenapi.VdiTypeMetadata {
				r.HasMetadata = true
				r.Enabled = r.Enabled || v.MetadataOfPool == poolRef
			}
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := strings.ToLower(out[i].Record.NameLabel), strings.ToLower(out[j].Record.NameLabel)
		return a < b || a == b && out[i].Record.UUID < out[j].Record.UUID
	})
	return out, nil
}

// ConfigureDR is DRConfigureCommand.RunCore: database replication off on the SRs left out of want (uuids), then on
// for the new ones (each a DelegatedAsyncAction, all run: the failures are given together).
func ConfigureDR(c *xapi.Connection, want []string) error {
	s := c.Session
	if text, ok := rules.DRConfigureWarning(false, len(want)); !ok {
		return errors.New(text)
	}
	srs, err := ReplicationSRs(s)
	if err != nil {
		return err
	}
	var enable, disable []ReplicationSR
	for _, uuid := range want {
		i := slices.IndexFunc(srs, func(r ReplicationSR) bool { return r.Record.UUID == uuid })
		if i < 0 {
			return fmt.Errorf("SR %s cannot store the VM configuration data", uuid)
		}
		if !srs[i].Enabled && !slices.ContainsFunc(enable, func(r ReplicationSR) bool { return r.Record.UUID == uuid }) {
			enable = append(enable, srs[i])
		}
	}
	for _, r := range srs {
		if r.Enabled && !slices.Contains(want, r.Record.UUID) {
			disable = append(disable, r)
		}
	}
	var calls []string
	if len(disable) > 0 {
		calls = append(calls, "sr.disable_database_replication")
	}
	if len(enable) > 0 {
		calls = append(calls, "sr.enable_database_replication")
	}
	if err := permitted(s, calls...); err != nil {
		return err
	}
	var failures []string
	for _, r := range disable {
		if err := wait(s)(xenapi.SR.AsyncDisableDatabaseReplication(s, r.Ref)); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", messages.Format(messages.Get("ACTION_DR_DISABLING_ON"), r.Record.NameLabel), xapi.Sanitize(err)))
		}
	}
	for _, r := range enable {
		if err := wait(s)(xenapi.SR.AsyncEnableDatabaseReplication(s, r.Ref)); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", messages.Format(messages.Get("ACTION_DR_ENABLING_ON"), r.Record.NameLabel), xapi.Sanitize(err)))
		}
	}
	if len(failures) > 0 {
		return errors.New(strings.Join(failures, "\n"))
	}
	return nil
}

// MetadataVDIs is GetMetadataVDIsAction on these SRs (uuids) followed by DrRules.LoadMetadata's filter: the latest
// metadata VDIs on shared SRs (no tools SR, no ISO library) that no VBD has attached.
func MetadataVDIs(s *xenapi.Session, srUUIDs []string) ([]xenapi.VDIRef, error) {
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	var out []xenapi.VDIRef
	for _, sr := range srs {
		if !rules.DRMetadataSR(sr.Shared, sr.IsToolsSr || sr.NameLabel == "XenServer Tools", sr.ContentType) || !slices.Contains(srUUIDs, sr.UUID) {
			continue
		}
		for _, ref := range sr.VDIs {
			v, err := xenapi.VDI.GetRecord(s, ref)
			if err != nil {
				return nil, xapi.Sanitize(err)
			}
			if v.Type != xenapi.VdiTypeMetadata || !v.MetadataLatest {
				continue
			}
			inUse := false
			for _, b := range v.VBDs {
				if attached, err := xenapi.VBD.GetCurrentlyAttached(s, b); err == nil && attached {
					inUse = true
				}
			}
			if !inUse {
				out = append(out, ref)
			}
		}
	}
	return out, nil
}

// DRMetadata is VdiLoadMetadataAction's PoolMetadata: the pool, vApps and real VMs of a metadata VDI's database.
type DRMetadata struct {
	VDI     xenapi.VDIRef
	VDIUUID string
	Pool    xenapi.PoolRecord
	VApps   map[xenapi.VMApplianceRef]xenapi.VMApplianceRecord
	VMs     map[xenapi.VMRef]xenapi.VMRecord
	all     map[xenapi.VMRef]xenapi.VMRecord // every VM record, snapshots included
}

// DRObjects are the uuids of what a recovery brings into the pool: the vApps, their VMs, the VMs chosen on their own,
// and the snapshots of all those VMs (XAPI recovers them with their VM).
type DRObjects struct {
	VApps, Members, VMs, Snapshots []string
}

// All is every uuid of o.
func (o DRObjects) All() []string { return slices.Concat(o.VApps, o.Members, o.VMs, o.Snapshots) }

// Objects are the DRObjects of recovering these vApps and VMs of the metadata.
func (md DRMetadata) Objects(apps []xenapi.VMApplianceRef, vms []xenapi.VMRef) DRObjects {
	var o DRObjects
	snaps := func(v xenapi.VMRecord) {
		for _, s := range v.Snapshots {
			if r, ok := md.all[s]; ok {
				o.Snapshots = append(o.Snapshots, r.UUID)
			}
		}
	}
	for _, a := range apps {
		o.VApps = append(o.VApps, md.VApps[a].UUID)
		for _, m := range md.VApps[a].VMs {
			if v, ok := md.VMs[m]; ok {
				o.Members = append(o.Members, v.UUID)
				snaps(v)
			}
		}
	}
	for _, v := range vms {
		o.VMs = append(o.VMs, md.VMs[v].UUID)
		snaps(md.VMs[v])
	}
	return o
}

// WithMetadata is VdiOpenDatabaseAction + VdiLoadMetadataAction: it opens the VDI's database (VDI.open_database),
// reads its pool, vApps and real VMs, runs fn with the metadata session and logs that session out. ok is false (fn
// not run) when the database holds no pool, or no vApp nor VM (PoolMetadata is null then).
func WithMetadata(s *xenapi.Session, vdi xenapi.VDIRef, fn func(meta *xenapi.Session, md DRMetadata) error) (ok bool, err error) {
	if err := permitted(s, "vdi.open_database"); err != nil {
		return false, err
	}
	uuid, err := xenapi.VDI.GetUUID(s, vdi)
	if err != nil {
		return false, xapi.Sanitize(err)
	}
	ref, err := xenapi.VDI.OpenDatabase(s, vdi)
	if err != nil {
		return false, xapi.Sanitize(err)
	}
	meta := s.WithRef(ref)
	defer meta.Logout() //nolint:errcheck // best effort, as the C# app
	md := DRMetadata{VDI: vdi, VDIUUID: uuid}
	pools, err := xenapi.Pool.GetAllRecords(meta)
	if err != nil {
		return false, nil // VdiLoadMetadataAction ignores what it cannot read
	}
	found := false
	for _, p := range pools {
		md.Pool, found = p, true
		break
	}
	if !found {
		return false, nil
	}
	if md.Pool.NameLabel == "" && md.Pool.Master != "" && md.Pool.Master != nullRef {
		md.Pool.NameLabel, _ = xenapi.Host.GetNameLabel(meta, md.Pool.Master)
	}
	if md.VApps, err = xenapi.VMAppliance.GetAllRecords(meta); err != nil {
		return false, nil
	}
	vms, err := xenapi.VM.GetAllRecords(meta)
	if err != nil {
		return false, nil
	}
	md.VMs, md.all = map[xenapi.VMRef]xenapi.VMRecord{}, vms
	for r, v := range vms {
		if rules.DRRealVM(v.IsATemplate, v.IsASnapshot, v.IsControlDomain) {
			md.VMs[r] = v
		}
	}
	if len(md.VApps) == 0 && len(md.VMs) == 0 {
		return false, nil
	}
	return true, fn(meta, md)
}

// Loose is DrRules.LooseVms' test: a VM outside a vApp (a vApp brings its own VMs).
func Loose(v xenapi.VMRecord) bool { return v.Appliance == "" || v.Appliance == nullRef }

// DRChoice is what the wizard recovers from one pool's metadata: the metadata VDI and the vApps and loose VMs chosen
// (uuids).
type DRChoice struct {
	VDI   string   `json:"vdi"`
	VApps []string `json:"vapps"`
	VMs   []string `json:"vms"`
}

// Resolve checks a choice against the metadata (trust boundary): every vApp is in it, every VM is one of its loose
// real VMs, once each.
func (md DRMetadata) Resolve(ch DRChoice) ([]xenapi.VMApplianceRef, []xenapi.VMRef, error) {
	var apps []xenapi.VMApplianceRef
	var vms []xenapi.VMRef
	for _, uuid := range ch.VApps {
		ref, ok := findRef(md.VApps, func(a xenapi.VMApplianceRecord) bool { return a.UUID == uuid })
		if !ok || slices.Contains(apps, ref) {
			return nil, nil, fmt.Errorf("vApp %s is not in the metadata of '%s'", uuid, md.Pool.NameLabel)
		}
		apps = append(apps, ref)
	}
	for _, uuid := range ch.VMs {
		ref, ok := findRef(md.VMs, func(v xenapi.VMRecord) bool { return v.UUID == uuid && Loose(v) })
		if !ok || slices.Contains(vms, ref) {
			return nil, nil, fmt.Errorf("VM %s is not in the metadata of '%s'", uuid, md.Pool.NameLabel)
		}
		vms = append(vms, ref)
	}
	return apps, vms, nil
}

func findRef[K comparable, V any](m map[K]V, match func(V) bool) (K, bool) {
	for k, v := range m {
		if match(v) {
			return k, true
		}
	}
	var zero K
	return zero, false
}

// MetadataVDI is the metadata VDI of a choice (its uuid), checked to be one.
func MetadataVDI(s *xenapi.Session, uuid string) (xenapi.VDIRef, error) {
	ref, err := xenapi.VDI.GetByUUID(s, uuid)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	if t, err := xenapi.VDI.GetType(s, ref); err != nil {
		return "", xapi.Sanitize(err)
	} else if t != xenapi.VdiTypeMetadata {
		return "", fmt.Errorf("VDI %s holds no pool metadata", uuid)
	}
	return ref, nil
}

// DRRecovered is what a recovery brought into the pool, for the start and the clean-up of a test failover.
type DRRecovered = DRObjects

// Empty: nothing was recovered.
func (o DRObjects) Empty() bool { return len(o.All()) == 0 }

// localUUIDs are the uuids of this pool's VMs (snapshots included) and vApps.
func localUUIDs(s *xenapi.Session) (map[string]bool, error) {
	out := map[string]bool{}
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	for _, v := range vms {
		out[v.UUID] = true
	}
	apps, err := xenapi.VMAppliance.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	for _, a := range apps {
		out[a.UUID] = true
	}
	return out, nil
}

// RecoverDR is DrRules.RecoverAction: per pool, open its metadata database and recover each chosen vApp and VM
// (DrRecoverAction: VM_appliance / VM.async_recover with force), then start them if asked (StartVMsAndAppliancesAction;
// an HA-protected VM that is not agile is never unprotected to start it). The failures are given together; what was
// recovered is returned with them. The current pool's own metadata is refused (DRFailoverWizardAppliancesPage).
// A test failover refuses to bring in any object (vApp, VM or snapshot) this pool has already, and records what it
// finds in the pool afterwards, a failed recovery's leftovers included, so that its clean-up only destroys what the
// test made; with an outcome unknown (ErrUncertain) nothing is recorded, as a recovery may still be running.
func RecoverDR(c *xapi.Connection, typ string, choices []DRChoice, start string) (DRRecovered, error) {
	var rec DRRecovered
	s := c.Session
	if !rules.DRKnownType(typ) || !slices.Contains(rules.DRStartModes(typ), start) {
		return rec, errors.New("unknown recovery type or start mode")
	}
	_, here, _, err := poolOf(s)
	if err != nil {
		return rec, err
	}
	calls := []string{"vdi.open_database", "vm.recover", "vm_appliance.recover"}
	if start != rules.DRStartNone {
		calls = append(calls, "vm.start", "vm_appliance.start")
	}
	if err := permitted(s, calls...); err != nil {
		return rec, err
	}
	type plan struct {
		vdi  xenapi.VDIRef
		ch   DRChoice
		objs DRObjects
	}
	var plans []plan
	for _, ch := range choices { // everything is checked before anything is recovered
		vdi, err := MetadataVDI(s, ch.VDI)
		if err != nil {
			return rec, err
		}
		var objs DRObjects
		ok, err := WithMetadata(s, vdi, func(_ *xenapi.Session, md DRMetadata) error {
			if md.Pool.UUID == here.UUID {
				return fmt.Errorf("'%s' is this pool: its own metadata cannot be recovered onto it", md.Pool.NameLabel)
			}
			apps, vms, err := md.Resolve(ch)
			objs = md.Objects(apps, vms)
			return err
		})
		if err != nil {
			return rec, err
		}
		if !ok {
			return rec, fmt.Errorf("VDI %s holds no pool metadata", ch.VDI)
		}
		plans = append(plans, plan{vdi, ch, objs})
	}
	if typ == rules.DRTest {
		existed, err := localUUIDs(s)
		if err != nil {
			return rec, err
		}
		for _, p := range plans {
			for _, u := range p.objs.All() {
				if existed[u] {
					return rec, fmt.Errorf("%s is in this pool already: a test failover does not replace it", u)
				}
			}
		}
	}
	var failures []string
	uncertain := false
	note := func(title string, err error) {
		var apiErr *xenapi.APIError // only a XAPI failure is a known outcome: a lost reply may hide a recovery still running
		uncertain = uncertain || errors.Is(err, ErrUncertain) || !errors.As(err, &apiErr)
		failures = append(failures, fmt.Sprintf("%s: %v", title, xapi.Sanitize(err)))
	}
	for _, p := range plans {
		ok, err := WithMetadata(s, p.vdi, func(meta *xenapi.Session, md DRMetadata) error {
			apps, vms, err := md.Resolve(p.ch)
			if err != nil {
				return err
			}
			for _, a := range apps {
				if err := wait(s)(xenapi.VMAppliance.AsyncRecover(meta, a, s.Ref(), true)); err != nil {
					note(messages.Format(messages.Get("ACTION_DR_RECOVER_APPLIANCE_TITLE"), md.VApps[a].NameLabel), err)
				} else {
					rec.VApps = append(rec.VApps, md.VApps[a].UUID)
				}
			}
			for _, v := range vms {
				if err := wait(s)(xenapi.VM.AsyncRecover(meta, v, s.Ref(), true)); err != nil {
					note(messages.Format(messages.Get("ACTION_DR_RECOVER_VM_TITLE"), md.VMs[v].NameLabel), err)
				} else {
					rec.VMs = append(rec.VMs, md.VMs[v].UUID)
				}
			}
			return nil
		})
		if err != nil {
			failures = append(failures, err.Error())
		} else if !ok {
			failures = append(failures, "The metadata of a pool could not be read again: nothing was recovered from it.")
		}
	}
	if typ == rules.DRTest { // what the test brought in, as the pool has it now (none of it was here before)
		var made DRObjects
		if now, err := localUUIDs(s); err != nil || uncertain {
			failures = append(failures, "What the test failover recovered could not be listed: remove it by hand.")
		} else {
			in := func(list []string) (out []string) {
				for _, u := range list {
					if now[u] {
						out = append(out, u)
					}
				}
				return out
			}
			for _, p := range plans {
				made.VApps = append(made.VApps, in(p.objs.VApps)...)
				made.Members = append(made.Members, in(p.objs.Members)...)
				made.VMs = append(made.VMs, in(p.objs.VMs)...)
				made.Snapshots = append(made.Snapshots, in(p.objs.Snapshots)...)
			}
		}
		started := DRObjects{VApps: slices.DeleteFunc(slices.Clone(rec.VApps), func(u string) bool { return !slices.Contains(made.VApps, u) }),
			VMs: slices.DeleteFunc(slices.Clone(rec.VMs), func(u string) bool { return !slices.Contains(made.VMs, u) })}
		rec = made
		if start != rules.DRStartNone && !started.Empty() {
			failures = append(failures, startRecovered(s, started, start == rules.DRStartPaused)...)
		}
	} else if start != rules.DRStartNone && !rec.Empty() {
		failures = append(failures, startRecovered(s, rec, start == rules.DRStartPaused)...)
	}
	if len(failures) > 0 {
		return rec, errors.New(strings.Join(failures, "\n"))
	}
	return rec, nil
}

// recovered is DrRules.RecoveredVms / RecoveredVapps: the recovered objects in this pool, waited for (up to 10 s) as
// they reach the database just after the calls return.
func recovered(s *xenapi.Session, rec DRRecovered) (map[string]xenapi.VMRef, map[string]xenapi.VMApplianceRef) {
	vms, apps := map[string]xenapi.VMRef{}, map[string]xenapi.VMApplianceRef{}
	for try := 0; try < 40 && len(vms)+len(apps) < len(rec.VMs)+len(rec.VApps); try++ {
		if try > 0 {
			time.Sleep(250 * time.Millisecond)
		}
		for _, u := range rec.VMs {
			if ref, err := xenapi.VM.GetByUUID(s, u); err == nil {
				vms[u] = ref
			}
		}
		for _, u := range rec.VApps {
			if ref, err := xenapi.VMAppliance.GetByUUID(s, u); err == nil {
				apps[u] = ref
			}
		}
	}
	return vms, apps
}

// startRecovered is StartVMsAndAppliancesAction: the halted recovered VMs, then the vApps (their VMs in order).
func startRecovered(s *xenapi.Session, rec DRRecovered, paused bool) []string {
	var failures []string
	vms, apps := recovered(s, rec)
	for _, ref := range vms {
		vm, err := xenapi.VM.GetRecord(s, ref)
		if err != nil {
			failures = append(failures, xapi.Sanitize(err).Error())
			continue
		}
		if vm.PowerState != xenapi.VMPowerStateHalted {
			continue
		}
		if why, err := haBlocker(s, vm); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", vm.NameLabel, err))
			continue
		} else if why != "" {
			failures = append(failures, fmt.Sprintf("%s: %s", vm.NameLabel, why))
			continue
		}
		if err := wait(s)(xenapi.VM.AsyncStart(s, ref, paused, false)); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", vm.NameLabel, xapi.Sanitize(err)))
		}
	}
	for _, ref := range apps {
		if err := wait(s)(xenapi.VMAppliance.AsyncStart(s, ref, paused)); err != nil {
			name, _ := xenapi.VMAppliance.GetNameLabel(s, ref)
			failures = append(failures, fmt.Sprintf("%s: %v", name, xapi.Sanitize(err)))
		}
	}
	return failures
}

// DRTestCleanup is DrRules.TestCleanupAction: the vApps a test failover recovered are destroyed with their VMs
// (ShutdownAndDestroyVmAppliancesAction), then the VMs it recovered on their own (ShutdownAndDestroyVMsAction). Only
// what the test recorded is touched: a vApp holding any other VM, or whose VMs could not all be removed, is kept (its
// VMs are shut down one by one, not with VM_appliance.hard_shutdown). The disks stay on the mirrored SR, as in the C#
// app.
func DRTestCleanup(c *xapi.Connection, rec DRRecovered) error {
	s := c.Session
	if err := permitted(s, "vm.hard_shutdown", "vm.destroy", "vm_appliance.destroy"); err != nil {
		return err
	}
	all := DRRecovered{VMs: slices.Concat(rec.VMs, rec.Members), VApps: rec.VApps} // a vApp's VMs left without it are removed too
	vms, apps := recovered(s, all)
	var failures []string
	for _, ref := range apps {
		app, err := xenapi.VMAppliance.GetRecord(s, ref)
		if err != nil {
			failures = append(failures, xapi.Sanitize(err).Error())
			continue
		}
		ours := true // every VM of the vApp, checked before any is touched
		for _, v := range app.VMs {
			if uuid, err := xenapi.VM.GetUUID(s, v); err != nil && !gone(err) || err == nil && !slices.Contains(rec.Members, uuid) {
				ours = false
			}
		}
		if !ours {
			failures = append(failures, fmt.Sprintf("vApp '%s' holds a VM the test failover did not recover (or one that cannot be read): it is kept", app.NameLabel))
			continue
		}
		removed := true
		for _, v := range app.VMs {
			if err := destroyRecoveredVM(s, v, rec.Snapshots); err != nil {
				failures = append(failures, err.Error())
				removed = false
			}
		}
		if !removed { // the vApp stays, so what is left of it can still be found
			continue
		}
		if err := xenapi.VMAppliance.Destroy(s, ref); err != nil && !gone(err) {
			failures = append(failures, fmt.Sprintf("%s: %v", app.NameLabel, xapi.Sanitize(err)))
		}
	}
	for _, ref := range vms {
		if err := destroyRecoveredVM(s, ref, rec.Snapshots); err != nil {
			failures = append(failures, err.Error())
		}
	}
	if len(vms) < len(all.VMs) || len(apps) < len(all.VApps) {
		failures = append(failures, "Some recovered vApps or VMs could not be found to remove.")
	}
	if len(failures) > 0 {
		return errors.New(strings.Join(failures, "\n"))
	}
	return nil
}

// destroyRecoveredVM is ShutdownAndDestroyVMsAction.DestroyVM after its hard shutdown: the snapshots, then the VM
// (VM.destroy takes its VBDs and VIFs; the disks stay). A VM with a snapshot the test did not recover is kept.
func destroyRecoveredVM(s *xenapi.Session, ref xenapi.VMRef, snapshots []string) error {
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		if gone(err) {
			return nil
		}
		return xapi.Sanitize(err)
	}
	for _, snap := range vm.Snapshots {
		if uuid, err := xenapi.VM.GetUUID(s, snap); err != nil && !gone(err) || err == nil && !slices.Contains(snapshots, uuid) {
			return fmt.Errorf("%s has a snapshot the test failover did not recover: it is kept", vm.NameLabel)
		}
	}
	if vm.PowerState == xenapi.VMPowerStateRunning || vm.PowerState == xenapi.VMPowerStatePaused || vm.PowerState == xenapi.VMPowerStateSuspended {
		if err := wait(s)(xenapi.VM.AsyncHardShutdown(s, ref)); err != nil {
			return fmt.Errorf("%s: %w", vm.NameLabel, xapi.Sanitize(err))
		}
	}
	for _, snap := range vm.Snapshots {
		if err := xenapi.VM.Destroy(s, snap); err != nil && !gone(err) {
			return fmt.Errorf("%s: %w", vm.NameLabel, xapi.Sanitize(err))
		}
	}
	if err := xenapi.VM.Destroy(s, ref); err != nil && !gone(err) {
		return fmt.Errorf("%s: %w", vm.NameLabel, xapi.Sanitize(err))
	}
	return nil
}

// gone: the object was deleted meanwhile (HANDLE_INVALID), which the C# app ignores.
func gone(err error) bool {
	var apiErr *xenapi.APIError
	return errors.As(err, &apiErr) && apiErr.Message == "HANDLE_INVALID"
}
