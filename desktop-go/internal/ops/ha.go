package ops

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// High availability (XenAdmin.Desktop.Core/HaRules.cs): EnableHAAction, SetHaPrioritiesAction and DisableHAAction
// (XenModel/Actions/Pool, XenModel/Actions/VM) and GetHeartbeatSRsAction's scan.

// HAStartup is one VM's VMStartupOptions: its restart priority (rules.HARestart, HABestEffort or HADoNotRestart),
// start order and delay before the next VM starts.
type HAStartup struct {
	VM         string `json:"vm"` // uuid
	Priority   string `json:"priority"`
	Order      int    `json:"order"`
	StartDelay int    `json:"start_delay"`
}

// HeartbeatSR is an SR of GetHeartbeatSRsAction: usable for the heartbeat, or the reason it is not.
type HeartbeatSR struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Reason string `json:"reason,omitempty"`
}

// HABusy: an HA operation runs on the pool (MainWindow's _haBusy and FindActiveHaAction, here the pool's current
// operations, which every client sees).
func HABusy(p xenapi.PoolRecord) bool {
	for _, op := range p.CurrentOperations {
		if op == xenapi.PoolAllowedOperationsHaEnable || op == xenapi.PoolAllowedOperationsHaDisable {
			return true
		}
	}
	return false
}

// NotAgileReason is HaRules.NotAgileReason: "" when the VM is agile, else the reason VM.assert_agile gives; an error
// that is not a XAPI failure (the connection) is returned as an error.
func NotAgileReason(s *xenapi.Session, ref xenapi.VMRef) (string, error) {
	err := xenapi.VM.AssertAgile(s, ref)
	var apiErr *xenapi.APIError
	switch {
	case err == nil:
		return "", nil
	case errors.As(err, &apiErr):
		return xapi.Sanitize(err).Error(), nil
	}
	return "", xapi.Sanitize(err)
}

// HeartbeatSRs is GetHeartbeatSRsAction: the shared SRs (not the tools SR), each asked whether it can hold the HA
// statefile (SR.assert_can_host_ha_statefile); the usable ones first, then by name.
func HeartbeatSRs(c *xapi.Connection) ([]HeartbeatSR, error) {
	s := c.Session
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	out := []HeartbeatSR{}
	for ref, sr := range srs {
		if !sr.Shared || sr.IsToolsSr || sr.NameLabel == "XenServer Tools" { // SR.IsToolsSR: unsuitable, not listed
			continue
		}
		h := HeartbeatSR{UUID: sr.UUID, Name: sr.NameLabel}
		if err := wait(s)(xenapi.SR.AsyncAssertCanHostHaStatefile(s, ref)); err != nil {
			var apiErr *xenapi.APIError
			if !errors.As(err, &apiErr) {
				return nil, xapi.Sanitize(err)
			}
			h.Reason = rules.HeartbeatUnsuitable(apiErr.Message, xapi.Sanitize(err).Error())
		}
		out = append(out, h)
	}
	sort.Slice(out, func(i, j int) bool {
		if (out[i].Reason == "") != (out[j].Reason == "") {
			return out[i].Reason == ""
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

// haEntry is a checked HAStartup.
type haEntry struct {
	HAStartup
	ref  xenapi.VMRef
	name string
}

// HAMaxFailures is HaRules.MaxFailures (HaNtolControl): the most server failures the pool tolerates with these
// priorities (VM uuid to priority).
func HAMaxFailures(c *xapi.Connection, priorities map[string]string) (int, error) {
	s := c.Session
	_, p, _, err := poolOf(s)
	if err != nil {
		return 0, err
	}
	plan := map[xenapi.VMRef]string{}
	for uuid, prio := range rules.HAPlanForAPI(priorities) {
		ref, err := xenapi.VM.GetByUUID(s, uuid)
		if err != nil {
			return 0, xapi.Sanitize(err)
		}
		plan[ref] = prio
	}
	return maxFailures(s, p, plan)
}

// maxFailures is Pool.GetMaximumTolerableHostFailures: pool.ha_compute_hypothetical_max_host_failures_to_tolerate for
// the VMs to restart, capped by the pool's HCI fault tolerance limit.
func maxFailures(s *xenapi.Session, p xenapi.PoolRecord, plan map[xenapi.VMRef]string) (int, error) {
	n, err := xenapi.Pool.HaComputeHypotheticalMaxHostFailuresToTolerate(s, plan)
	if err != nil {
		return 0, xapi.Sanitize(err)
	}
	return rules.HAMaxTolerable(n, p.OtherConfig), nil
}

// checkPlan validates a plan at the trust boundary: known protectable VMs, once each, known priorities, an order and
// a delay of 0 or more. When the plan protects more (HA being turned on, a VM newly restarted, more failures to
// tolerate) it also checks that the newly restarted VMs are agile (HaRules.Allowed) and that the pool's plan (these
// priorities over the other VMs' own) holds the failures to tolerate (HaRules.PlanError: no overcommit); an order or
// delay change on a degraded plan still goes through, as XAPI allows.
func checkPlan(s *xenapi.Session, p xenapi.PoolRecord, plan []HAStartup, ntol int) ([]haEntry, error) {
	if ntol < 0 {
		return nil, errors.New("enter the number of server failures to tolerate (0 or more)")
	}
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	byUUID := map[string]xenapi.VMRef{}
	priorities := map[string]string{}
	for ref, v := range vms {
		if rules.HAProtectable(v.IsATemplate, v.IsASnapshot, v.IsControlDomain, v.OtherConfig["HideFromXenCenter"] == "true", v.NameLabel) {
			byUUID[v.UUID] = ref
			priorities[v.UUID] = v.HaRestartPriority
		}
	}
	var out []haEntry
	seen := map[string]bool{}
	raise := !p.HaEnabled || ntol > p.HaHostFailuresToTolerate
	for _, e := range plan {
		ref, ok := byUUID[e.VM]
		switch {
		case !ok:
			return nil, fmt.Errorf("VM %s is not one HA can protect", e.VM)
		case seen[e.VM]:
			return nil, fmt.Errorf("VM %s is given twice", e.VM)
		case !rules.HAValidPriority(e.Priority):
			return nil, fmt.Errorf("unknown HA restart priority %q", e.Priority)
		case e.Order < 0 || e.StartDelay < 0:
			return nil, errors.New("the start order and the delay must be 0 or more")
		}
		seen[e.VM] = true
		name := vms[ref].NameLabel
		if rules.HAIsRestart(e.Priority) && (!p.HaEnabled || !rules.HAIsRestart(vms[ref].HaRestartPriority)) {
			raise = true
			why, err := NotAgileReason(s, ref)
			if err != nil {
				return nil, err
			}
			if why != "" {
				return nil, fmt.Errorf("'%s' is not agile, so HA cannot guarantee to restart it: %s", name, why)
			}
		}
		priorities[e.VM] = e.Priority
		out = append(out, haEntry{HAStartup: e, ref: ref, name: name})
	}
	if !raise {
		return out, nil
	}
	api := map[xenapi.VMRef]string{}
	for uuid, prio := range rules.HAPlanForAPI(priorities) {
		api[byUUID[uuid]] = prio
	}
	max, err := maxFailures(s, p, api)
	if err != nil {
		return nil, err
	}
	if why := rules.HAPlanError(ntol, max); why != "" {
		return nil, errors.New(why)
	}
	return out, nil
}

// setStartup is the VM.SetHaRestartPriority, VM.set_order and VM.set_start_delay of the HA actions.
func setStartup(s *xenapi.Session, e haEntry) error {
	if err := xenapi.VM.SetHaRestartPriority(s, e.ref, e.Priority); err != nil {
		return fmt.Errorf("%s: %w", e.name, xapi.Sanitize(err))
	}
	if err := xenapi.VM.SetOrder(s, e.ref, e.Order); err != nil {
		return fmt.Errorf("%s: %w", e.name, xapi.Sanitize(err))
	}
	if err := xenapi.VM.SetStartDelay(s, e.ref, e.StartDelay); err != nil {
		return fmt.Errorf("%s: %w", e.name, xapi.Sanitize(err))
	}
	return nil
}

// EnableHA is EnableHAAction: the VMs' priorities, start order and delay, the failures to tolerate, then
// pool.enable_ha with the heartbeat SR (heartbeatSR is its uuid).
func EnableHA(c *xapi.Connection, heartbeatSR string, plan []HAStartup, ntol int) error {
	s := c.Session
	poolRef, p, _, err := poolOf(s)
	if err != nil {
		return err
	}
	switch {
	case HABusy(p):
		return errors.New(rules.HAConfigureBlocker(true, p.HaEnabled, true, ""))
	case p.HaEnabled:
		text, _ := messages.FailureText("HA_IS_ENABLED", nil)
		return errors.New(text)
	}
	sr, err := xenapi.SR.GetByUUID(s, heartbeatSR)
	if err != nil {
		return xapi.Sanitize(err)
	}
	entries, err := checkPlan(s, p, plan, ntol)
	if err != nil {
		return err
	}
	// EnableHAAction.ApiMethodsToRoleCheck: nothing is changed when the role cannot finish
	if err := permitted(s, "vm.set_ha_restart_priority", "vm.set_order", "vm.set_start_delay", "pool.set_ha_host_failures_to_tolerate", "pool.enable_ha"); err != nil {
		return err
	}
	for _, e := range entries {
		if err := setStartup(s, e); err != nil {
			return err
		}
	}
	if err := xenapi.Pool.SetHaHostFailuresToTolerate(s, poolRef, ntol); err != nil {
		return xapi.Sanitize(err)
	}
	err = wait(s)(xenapi.Pool.AsyncEnableHa(s, []xenapi.SRRef{sr}, map[string]string{}))
	var apiErr *xenapi.APIError
	if errors.As(err, &apiErr) && apiErr.Message == "VDI_NOT_AVAILABLE" && len(apiErr.Data) > 0 { // the VDI's uuid, not its ref
		if uuid, uerr := xenapi.VDI.GetUUID(s, xenapi.VDIRef(apiErr.Data[0])); uerr == nil {
			text, _ := messages.FailureText("VDI_NOT_AVAILABLE", []string{uuid})
			return errors.New(text)
		}
	}
	return xapi.Sanitize(err)
}

// SetHAPriorities is SetHaPrioritiesAction on a pool with HA on: the VMs that lose their guaranteed restart first,
// then the failures to tolerate, then the VMs that gain it (so the plan holds at every step), and pool.sync_database.
func SetHAPriorities(c *xapi.Connection, plan []HAStartup, ntol int) error {
	return setHAPriorities(c.Session, plan, ntol, false)
}

// setHAPriorities with keepNtol (the Start Options page) leaves the pool's failures to tolerate as they are now, never
// writing back a value read earlier.
func setHAPriorities(s *xenapi.Session, plan []HAStartup, ntol int, keepNtol bool) error {
	poolRef, p, _, err := poolOf(s)
	if err != nil {
		return err
	}
	if why := rules.HADisableBlocker(HABusy(p), p.HaEnabled); why != "" { // HA on, no HA operation running
		return errors.New(why)
	}
	if keepNtol {
		ntol = p.HaHostFailuresToTolerate
	}
	entries, err := checkPlan(s, p, plan, ntol)
	if err != nil {
		return err
	}
	if err := permitted(s, "vm.set_ha_restart_priority", "vm.set_order", "vm.set_start_delay", "pool.set_ha_host_failures_to_tolerate", "pool.sync_database"); err != nil {
		return err
	}
	for _, restart := range []bool{false, true} {
		if restart && !keepNtol {
			if err := xenapi.Pool.SetHaHostFailuresToTolerate(s, poolRef, ntol); err != nil {
				return xapi.Sanitize(err)
			}
		}
		for _, e := range entries {
			if rules.HAIsRestart(e.Priority) == restart {
				if err := setStartup(s, e); err != nil {
					return err
				}
			}
		}
	}
	return xapi.Sanitize(wait(s)(xenapi.Pool.AsyncSyncDatabase(s)))
}

// DisableHABlocker is why Disable cannot run ("" when it can): HA off, an HA operation running, or a role without
// pool.disable_ha (DisableHAAction.ApiMethodsToRoleCheck).
func DisableHABlocker(c *xapi.Connection) string {
	_, p, _, err := poolOf(c.Session)
	if err != nil {
		return err.Error()
	}
	if why := rules.HADisableBlocker(HABusy(p), p.HaEnabled); why != "" {
		return why
	}
	if err := permitted(c.Session, "pool.disable_ha"); err != nil {
		return err.Error()
	}
	return ""
}

// DisableHA is DisableHAAction: pool.disable_ha.
func DisableHA(c *xapi.Connection) error {
	if why := DisableHABlocker(c); why != "" {
		return errors.New(why)
	}
	s := c.Session
	return xapi.Sanitize(wait(s)(xenapi.Pool.AsyncDisableHa(s)))
}
