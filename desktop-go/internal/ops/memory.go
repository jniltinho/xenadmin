package ops

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// MemoryFacts reads what MemoryRules needs of a VM (rules.MemoryVM); restricted is the DMC licence restriction.
func MemoryFacts(s *xenapi.Session, ref xenapi.VMRef, vm xenapi.VMRecord, ioDrivers, windows bool) (m rules.MemoryVM, restricted bool, err error) {
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return m, false, xapi.Sanitize(err)
	}
	restricted = featureForbidden(hosts, "restrict_dmc")
	supports := vm.MemoryDynamicMin != vm.MemoryStaticMax // VM.SupportsBallooning
	if !vm.IsATemplate {
		gm, err := xenapi.VMGuestMetrics.GetRecord(s, vm.GuestMetrics)
		_, balloon := gm.Other["feature-balloon"]
		supports = err == nil && balloon
	}
	m = rules.MemoryVM{StaticMin: int64(vm.MemoryStaticMin), StaticMax: int64(vm.MemoryStaticMax),
		DynamicMin: int64(vm.MemoryDynamicMin), DynamicMax: int64(vm.MemoryDynamicMax),
		Halted: vm.PowerState == xenapi.VMPowerStateHalted, Running: vm.PowerState == xenapi.VMPowerStateRunning,
		Template: vm.IsATemplate, Ballooning: supports && !restricted, IODrivers: ioDrivers, Windows: windows,
		Ratio: 0.25, MaxDynamicMin: -1, MaxAllowed: rules.Restriction(vm.Recommendations, "memory-static-max", "max", 1<<40)}
	key := "memory-ratio-pv" // MemoryRules.Ratio
	if vm.HVMBootPolicy != "" {
		key = "memory-ratio-hvm"
	}
	if pools, err := xenapi.Pool.GetAllRecords(s); err == nil {
		for _, p := range pools {
			if f, err := strconv.ParseFloat(p.OtherConfig[key], 64); err == nil && f >= 0 && f <= 1 {
				m.Ratio = f
			}
		}
	}
	// MemoryRules.MaxDynMin: what the VM's server can still give it while it runs
	host, ok := hosts[vm.ResidentOn]
	if (vm.PowerState == xenapi.VMPowerStateRunning || vm.PowerState == xenapi.VMPowerStatePaused) && ok {
		if hm, err := xenapi.HostMetrics.GetRecord(s, host.Metrics); err == nil {
			used := int64(host.MemoryOverhead)
			for _, r := range host.ResidentVMs {
				other, err := xenapi.VM.GetRecord(s, r)
				if err != nil {
					continue
				}
				used += int64(other.MemoryOverhead)
				if other.IsControlDomain {
					if om, err := xenapi.VMMetrics.GetRecord(s, other.Metrics); err == nil {
						used += int64(om.MemoryActual)
					}
				} else if r != ref {
					used += int64(other.MemoryDynamicMin)
				}
			}
			m.MaxDynamicMin = max(0, int64(hm.MemoryTotal)-used)
		}
	}
	return m, restricted, nil
}

// ChangeMemory is ChangeMemorySettingsAction, after MemoryRules.Validate on the server's current values: a new static
// maximum shuts a running VM down (clean if allowed, else hard), sets the limits and starts it again on the same server
// (even if the change failed); else only the dynamic range changes. checkVM gives the virtualization facts of the VM.
func ChangeMemory(c *xapi.Connection, vmUUID string, dynMin, dynMax, staticMax int64, checkVM func(xenapi.VMRecord) (ioDrivers, windows bool)) error {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if vm.IsASnapshot || vm.IsControlDomain || len(vm.CurrentOperations) > 0 {
		return errors.New("the VM cannot change its memory now")
	}
	io, win := checkVM(vm)
	m, _, err := MemoryFacts(s, ref, vm, io, win)
	if err != nil {
		return err
	}
	if !m.Ballooning { // basic mode without a range: one value
		dynMin, dynMax = staticMax, staticMax
	}
	if why := rules.ValidateMemory(m, dynMin, dynMax, staticMax); why != "" {
		return errors.New(why)
	}
	dynMin, dynMax, staticMax = rules.MemorySettings(m, dynMin, dynMax, staticMax)
	staticChanged := staticMax != m.StaticMax
	reboot := staticChanged && !m.Halted
	// ApiMethodsToRoleCheck: never shut the VM down when the role cannot change the memory or start it again
	calls := []string{"vm.set_memory_dynamic_range"}
	if staticChanged {
		calls = []string{"vm.set_memory_limits"}
	}
	if reboot {
		shutdown, start := "vm.hard_shutdown", "vm.start"
		if contains(vm.AllowedOperations, xenapi.VMOperationsCleanShutdown) {
			shutdown = "vm.clean_shutdown"
		}
		if vm.ResidentOn != nullRef {
			start = "vm.start_on"
		}
		calls = append(calls, shutdown, start)
	}
	if err := permitted(s, calls...); err != nil {
		return err
	}
	if reboot {
		if why, err := haBlocker(s, vm); err != nil { // MemoryRules.ChangeAction: an HA restart is never unprotected silently
			return err
		} else if why != "" {
			return errors.New(why)
		}
		if err := shutDown(s, ref, vm); err != nil {
			return err
		}
	}
	if staticChanged {
		err = xenapi.VM.SetMemoryLimits(s, ref, int(m.StaticMin), int(staticMax), int(dynMin), int(dynMax))
	} else {
		err = xenapi.VM.SetMemoryDynamicRange(s, ref, int(dynMin), int(dynMax))
	}
	err = xapi.Sanitize(err)
	if reboot { // boot the VM even if the memory settings failed after the shutdown
		var start error
		if vm.ResidentOn != nullRef {
			start = wait(s)(xenapi.VM.AsyncStartOn(s, ref, vm.ResidentOn, false, false))
		} else {
			start = wait(s)(xenapi.VM.AsyncStart(s, ref, false, false))
		}
		err = errors.Join(err, xapi.Sanitize(start))
	}
	return err
}

// shutDown is the restart's first half: a clean shutdown if allowed, else a hard one; a failure is ignored when the VM
// is halted (or another shutdown runs), and then it waits for Halted.
func shutDown(s *xenapi.Session, ref xenapi.VMRef, vm xenapi.VMRecord) error {
	var err error
	if contains(vm.AllowedOperations, xenapi.VMOperationsCleanShutdown) {
		err = wait(s)(xenapi.VM.AsyncCleanShutdown(s, ref))
	} else {
		err = wait(s)(xenapi.VM.AsyncHardShutdown(s, ref))
	}
	for deadline := time.Now().Add(10 * time.Minute); ; time.Sleep(pollEvery) {
		state, perr := xenapi.VM.GetPowerState(s, ref)
		if perr == nil && state == xenapi.VMPowerStateHalted {
			return nil
		}
		if err != nil {
			ops, _ := xenapi.VM.GetCurrentOperations(s, ref)
			shutting := false
			for _, op := range ops {
				shutting = shutting || op == xenapi.VMOperationsCleanShutdown || op == xenapi.VMOperationsHardShutdown || op == xenapi.VMOperationsShutdown
			}
			if !shutting {
				return xapi.Sanitize(err)
			}
		}
		if time.Now().After(deadline) {
			return errors.New("the VM did not shut down")
		}
	}
}

// permitted is the actions' RBAC check (ApiMethodsToRoleCheck, Role.CanPerform): a local superuser may call anything,
// another session only the API calls its role grants.
func permitted(s *xenapi.Session, calls ...string) error {
	if su, err := s.GetIsLocalSuperuser(s.Ref()); err == nil && su {
		return nil
	}
	perms, err := s.GetRbacPermissions(s.Ref())
	if err != nil {
		return xapi.Sanitize(err)
	}
	for _, c := range calls {
		if !slices.ContainsFunc(perms, func(p string) bool { return strings.EqualFold(p, c) }) {
			return errors.New(rbacDenied(c))
		}
	}
	return nil
}

func rbacDenied(call string) string {
	text, _ := messages.FailureText("RBAC_PERMISSION_DENIED", []string{call})
	return text
}
