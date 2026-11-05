package ops

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// the host other_config keys of XenModel's Host extension (MAINTENANCE_MODE, MAINTENANCE_MODE_EVACUATED_VMS_*)
const (
	maintenanceKey = "MAINTENANCE_MODE"
	evacMigrated   = "MAINTENANCE_MODE_EVACUATED_VMS_MIGRATED"
	evacHalted     = "MAINTENANCE_MODE_EVACUATED_VMS_HALTED"
	evacSuspended  = "MAINTENANCE_MODE_EVACUATED_VMS_SUSPENDED"
)

// host is a server and what the server commands read about it and its pool.
type host struct {
	ref  xenapi.HostRef
	rec  xenapi.HostRecord
	pool xenapi.PoolRef
	prec xenapi.PoolRecord
	info rules.HostInfo
	newC xenapi.HostRef // the member a coordinator hands over to
	vms  map[xenapi.VMRef]xenapi.VMRecord
}

func readHost(s *xenapi.Session, uuid string) (*host, error) {
	ref, err := xenapi.Host.GetByUUID(s, uuid)
	if err != nil {
		return nil, err
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return nil, err
	}
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return nil, err
	}
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return nil, err
	}
	h := &host{ref: ref, rec: hosts[ref], vms: vms}
	for pr, p := range pools {
		h.pool, h.prec = pr, p
	}
	h.info = rules.HostInfo{Name: h.rec.NameLabel, Coordinator: h.prec.Master == ref, PoolHosts: len(hosts), HAEnabled: h.prec.HaEnabled}
	// EvacuateHostDialog: the coordinator hands over to an enabled, live member (the first by name)
	var names []string
	byName := map[string]xenapi.HostRef{}
	for r, o := range hosts {
		if r == ref || !o.Enabled {
			continue
		}
		if live, err := xenapi.HostMetrics.GetLive(s, o.Metrics); err != nil || !live {
			continue
		}
		names, byName[o.NameLabel] = append(names, o.NameLabel), r
	}
	sort.Strings(names)
	if len(names) > 0 {
		h.info.NewCoordinator, h.newC = names[0], byName[names[0]]
	}
	for _, v := range h.rec.ResidentVMs {
		vm := vms[v]
		local, err := onLocalStorage(s, vm)
		if err != nil {
			return nil, err
		}
		h.info.VMs = append(h.info.VMs, rules.HostVM{Name: vm.NameLabel, Running: vm.PowerState == xenapi.VMPowerStateRunning,
			ControlDomain: vm.IsControlDomain, HCIWarn: vm.OtherConfig["hci-warn-before-shutdown"] != "", LocalStorage: local})
	}
	return h, nil
}

// onLocalStorage: a disk or CD of the VM is on an SR that is not shared (VM.IsOnSharedStorage).
func onLocalStorage(s *xenapi.Session, vm xenapi.VMRecord) (bool, error) {
	if vm.IsControlDomain {
		return false, nil
	}
	for _, b := range vm.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil {
			return false, err
		}
		if vbd.VDI == nullRef || vbd.Empty {
			continue
		}
		vdi, err := xenapi.VDI.GetRecord(s, vbd.VDI)
		if err != nil {
			return false, err
		}
		sr, err := xenapi.SR.GetRecord(s, vdi.SR)
		if err != nil {
			return false, err
		}
		if !sr.Shared {
			return true, nil
		}
	}
	return false, nil
}

// HostInfo is what the window shows before a server command (its confirmation texts).
func HostInfo(c *xapi.Connection, uuid string) (rules.HostInfo, error) {
	h, err := readHost(c.Session, uuid)
	if err != nil {
		return rules.HostInfo{}, xapi.Sanitize(err)
	}
	return h.info, nil
}

// keepNtol is HostRules.KeepNtol: when HA would have to tolerate fewer failures for this server to go, the operation
// is refused instead of lowering the pool's protection without consent (MaybeReduceNtolBeforeOp).
// shortcut: no failures-to-tolerate dialog; the user lowers it in Pool > High Availability first.
func keepNtol(s *xenapi.Session, h *host) error {
	if !h.prec.HaEnabled {
		return nil
	}
	max, err := xenapi.Pool.HaComputeMaxHostFailuresToTolerate(s)
	if err != nil {
		return err
	}
	if target := maxInt(0, max-1); h.prec.HaHostFailuresToTolerate > target {
		return fmt.Errorf("HA would have to tolerate %d server failures instead of %d for '%s' to go; lower it in Pool > High Availability first",
			target, h.prec.HaHostFailuresToTolerate, h.info.Name)
	}
	return nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// HostPower is ShutdownHostAction / RebootHostAction: disable the server, shut its running VMs down (clean when
// allowed), then shut it down or reboot it; on a failure the server is enabled again.
// shortcut: no progress per VM; a coordinator's connection ends with it (the window reconnects).
func HostPower(c *xapi.Connection, uuid string, reboot bool) error {
	s := c.Session
	h, err := readHost(s, uuid)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if err := keepNtol(s, h); err != nil {
		return xapi.Sanitize(err)
	}
	if err := wait(s)(xenapi.Host.AsyncDisable2(s, h.ref)); err != nil {
		return xapi.Sanitize(err)
	}
	fail := func(err error) error {
		if h.rec.Enabled && !errors.Is(err, ErrUncertain) {
			_ = wait(s)(xenapi.Host.AsyncEnable(s, h.ref)) // let the user shut the VMs down by hand
		}
		return xapi.Sanitize(err)
	}
	// the VMs resident now: one may have moved away before the server was disabled
	residents, err := xenapi.Host.GetResidentVMs(s, h.ref)
	if err != nil {
		return fail(err)
	}
	for _, v := range residents {
		vm, err := xenapi.VM.GetRecord(s, v)
		if err != nil {
			return fail(err)
		}
		if vm.IsControlDomain || vm.PowerState != xenapi.VMPowerStateRunning || vm.ResidentOn != h.ref {
			continue
		}
		op := xenapi.VM.AsyncHardShutdown
		for _, a := range vm.AllowedOperations {
			if a == xenapi.VMOperationsCleanShutdown {
				op = xenapi.VM.AsyncCleanShutdown
			}
		}
		if err := wait(s)(op(s, v)); err != nil {
			return fail(err)
		}
	}
	if reboot {
		err = wait(s)(xenapi.Host.AsyncReboot(s, h.ref))
	} else {
		err = wait(s)(xenapi.Host.AsyncShutdown(s, h.ref))
	}
	if err != nil {
		return fail(err)
	}
	return nil
}

// EnterMaintenance is HostRules.EnterMaintenance: save where the VMs are (Host.SaveEvacuatedVMs), disable the server
// with MAINTENANCE_MODE=true, host.evacuate with XAPI's placement, and hand a coordinator's role to another member.
// Every step is a task polled to its end, so a failure (never a slow evacuation) enables the server again.
func EnterMaintenance(c *xapi.Connection, uuid string) error {
	s := c.Session
	h, err := readHost(s, uuid)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if !h.rec.Enabled {
		return errors.New("the server is already in maintenance mode")
	}
	if why := rules.CannotEnterMaintenance(h.info); why != "" {
		return errors.New(why)
	}
	if err := keepNtol(s, h); err != nil {
		return xapi.Sanitize(err)
	}
	if err := saveEvacuated(s, h); err != nil {
		return xapi.Sanitize(err)
	}
	err = func() error {
		if err := wait(s)(xenapi.Host.AsyncDisable2(s, h.ref)); err != nil {
			return err
		}
		_ = xenapi.Host.RemoveFromOtherConfig(s, h.ref, maintenanceKey)
		if err := xenapi.Host.AddToOtherConfig(s, h.ref, maintenanceKey, "true"); err != nil {
			return err
		}
		if err := wait(s)(xenapi.Host.AsyncEvacuate2(s, h.ref)); err != nil {
			return err
		}
		if h.info.Coordinator && h.newC != "" {
			return wait(s)(xenapi.Pool.AsyncDesignateNewMaster(s, h.newC))
		}
		return nil
	}()
	if err != nil {
		if !errors.Is(err, ErrUncertain) { // an evacuation that may still run is never undone under it
			_ = enable(s, h.ref) // EvacuateHostAction: back out of maintenance mode
		}
		return xapi.Sanitize(err)
	}
	return nil
}

// saveEvacuated is Host.SaveEvacuatedVMs: the first time every VM resident on the server is saved as migrated; once
// lists exist (a second try after a partial evacuation), the saved VMs are sorted again by their power state, so
// the ones already moved away stay on the lists.
func saveEvacuated(s *xenapi.Session, h *host) error {
	saved := 0
	byUUID := map[string]xenapi.VMRecord{}
	for _, v := range h.vms {
		byUUID[v.UUID] = v
	}
	var all []xenapi.VMRecord
	seen := map[string]bool{}
	for _, key := range []string{evacMigrated, evacSuspended, evacHalted} {
		for _, u := range strings.Split(h.rec.OtherConfig[key], ",") {
			if vm, ok := byUUID[u]; ok && !seen[u] {
				seen[u], saved = true, saved+1
				all = append(all, vm)
			}
		}
	}
	lists := map[string][]string{}
	if saved < len(h.rec.ResidentVMs)-1 { // the control domain is resident too
		for _, v := range h.rec.ResidentVMs {
			if vm := h.vms[v]; !vm.IsControlDomain {
				lists[evacMigrated] = append(lists[evacMigrated], vm.UUID)
			}
		}
		return saveList(s, h.ref, evacMigrated, lists[evacMigrated])
	}
	for _, vm := range all {
		switch vm.PowerState {
		case xenapi.VMPowerStateHalted:
			lists[evacHalted] = append(lists[evacHalted], vm.UUID)
		case xenapi.VMPowerStateRunning:
			lists[evacMigrated] = append(lists[evacMigrated], vm.UUID)
		case xenapi.VMPowerStateSuspended:
			lists[evacSuspended] = append(lists[evacSuspended], vm.UUID)
		}
	}
	for _, key := range []string{evacMigrated, evacHalted, evacSuspended} {
		if err := saveList(s, h.ref, key, lists[key]); err != nil {
			return err
		}
	}
	return nil
}

// saveList is Host.SaveVMList: the VM uuids, comma-separated, in the host's other_config.
func saveList(s *xenapi.Session, ref xenapi.HostRef, key string, uuids []string) error {
	_ = xenapi.Host.RemoveFromOtherConfig(s, ref, key)
	return xenapi.Host.AddToOtherConfig(s, ref, key, strings.Join(uuids, ","))
}

func enable(s *xenapi.Session, ref xenapi.HostRef) error {
	_ = xenapi.Host.RemoveFromOtherConfig(s, ref, maintenanceKey)
	return wait(s)(xenapi.Host.AsyncEnable(s, ref))
}

// evacuated is Host.GetMigrated/Halted/SuspendedEvacuatedVMs: the saved VMs of each list now in that power state and
// not on the server (HostRules.VmsToRestore).
func evacuated(h *host) (migrated, halted, suspended []xenapi.VMRef) {
	byUUID := map[string]xenapi.VMRef{}
	for r, v := range h.vms {
		byUUID[v.UUID] = r
	}
	var seen = map[xenapi.VMRef]bool{}
	for key, state := range map[string]xenapi.VMPowerState{evacMigrated: xenapi.VMPowerStateRunning, evacHalted: xenapi.VMPowerStateHalted, evacSuspended: xenapi.VMPowerStateSuspended} {
		for _, u := range strings.Split(h.rec.OtherConfig[key], ",") {
			r, ok := byUUID[u]
			if !ok || seen[r] {
				continue
			}
			vm := h.vms[r]
			if vm.PowerState != state || vm.ResidentOn == h.ref {
				continue
			}
			seen[r] = true
			switch state {
			case xenapi.VMPowerStateRunning:
				migrated = append(migrated, r)
			case xenapi.VMPowerStateHalted:
				halted = append(halted, r)
			default:
				suspended = append(suspended, r)
			}
		}
	}
	return
}

// VMsToRestore names the server and the evacuated VMs that Exit Maintenance Mode can put back.
func VMsToRestore(c *xapi.Connection, uuid string) (string, []string, error) {
	h, err := readHost(c.Session, uuid)
	if err != nil {
		return "", nil, xapi.Sanitize(err)
	}
	m, ha, su := evacuated(h)
	var names []string
	for _, r := range append(append(m, ha...), su...) {
		names = append(names, h.vms[r].NameLabel)
	}
	sort.Strings(names)
	return h.info.Name, names, nil
}

// ExitMaintenance is EnableHostAction: enable the server and, when asked, put the evacuated VMs back (live migrate
// the running ones, start the halted ones and resume the suspended ones on it).
func ExitMaintenance(c *xapi.Connection, uuid string, restore bool) error {
	s := c.Session
	h, err := readHost(s, uuid)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if err := enable(s, h.ref); err != nil {
		return xapi.Sanitize(err)
	}
	if !restore {
		return nil
	}
	migrated, halted, suspended := evacuated(h)
	for _, k := range []string{evacMigrated, evacHalted, evacSuspended} {
		_ = xenapi.Host.RemoveFromOtherConfig(s, h.ref, k)
	}
	var errs []error
	for _, v := range migrated {
		errs = append(errs, wait(s)(xenapi.VM.AsyncPoolMigrate(s, v, h.ref, map[string]string{"live": "true"})))
	}
	for _, v := range halted {
		errs = append(errs, wait(s)(xenapi.VM.AsyncStartOn(s, v, h.ref, false, false)))
	}
	for _, v := range suspended {
		errs = append(errs, wait(s)(xenapi.VM.AsyncResumeOn(s, v, h.ref, false, false)))
	}
	return xapi.Sanitize(errors.Join(errs...))
}
