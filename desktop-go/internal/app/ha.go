package app

import (
	"errors"
	"sort"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// What Pool > High Availability reads (XenAdmin.Desktop.Core/HaRules.cs, XenAdmin.Desktop/HaPlanEditor.cs): the
// pool's HA state, and the plan of the HA wizard and of the restart priorities dialog.

// HAConfig is the pool's HA state and, when Configure can run, its plan (AssignPriorities).
type HAConfig struct {
	Pool       string             `json:"pool"`
	HAEnabled  bool               `json:"ha_enabled"`
	Blocker    string             `json:"blocker,omitempty"` // Configure cannot run (HaRules.ConfigureBlocker): say why
	Ntol       int                `json:"ntol"`              // ha_host_failures_to_tolerate (HA on)
	Hosts      int                `json:"hosts"`             // the failures to tolerate go up to this (at least 1)
	Priorities []rules.HAPriority `json:"priorities"`
	VMs        []HAPlanVM         `json:"vms"`
	// the plan's texts: too many failures to tolerate (HaRules.PlanError), and EditVmHaPrioritiesDialog's question
	// when none is tolerated
	Overcommit    string `json:"overcommit"`
	NtolZeroQuery string `json:"ntol_zero_query"`
}

// HAPlanVM is a VM HA can protect: its priority to start with (HaRules.InitialPriority), start order and delay, and
// why it is not agile ("" when it is: VM.assert_agile).
type HAPlanVM struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Priority   string `json:"priority"`
	Order      int    `json:"order"`
	StartDelay int    `json:"start_delay"`
	NotAgile   string `json:"not_agile,omitempty"`
}

// ReadHAConfig reads the pool's HA state; with withPlan and no blocker, also the protectable VMs, each checked for
// agility (AssignPriorities.UpdateVMsAgility; a VM that is not agile can only be best effort or not restarted).
func ReadHAConfig(c *xapi.Connection, withPlan bool) (HAConfig, error) {
	s := c.Session
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return HAConfig{}, xapi.Sanitize(err)
	}
	var cfg HAConfig
	var p xenapi.PoolRecord
	found := false
	for _, r := range pools {
		p, found = r, true
	}
	if !found {
		return HAConfig{}, errors.New("the server has no pool")
	}
	hosts, err := xenapi.Host.GetAll(s)
	if err != nil {
		return HAConfig{}, xapi.Sanitize(err)
	}
	cfg.Pool, cfg.HAEnabled, cfg.Ntol, cfg.Hosts = p.NameLabel, p.HaEnabled, p.HaHostFailuresToTolerate, max(1, len(hosts))
	statefile := false
	for _, sf := range p.HaStatefiles { // HaRules.ConfigureBlocker: the statefile VDI must resolve
		if _, err := xenapi.VDI.GetUUID(s, xenapi.VDIRef(sf)); err == nil {
			statefile = true
		}
	}
	cfg.Blocker = rules.HAConfigureBlocker(ops.HABusy(p), p.HaEnabled, statefile, ops.Ellipsise(p.NameLabel, 30))
	if !withPlan || cfg.Blocker != "" {
		return cfg, nil
	}
	cfg.Priorities = rules.HAPriorities()
	cfg.Overcommit, cfg.NtolZeroQuery = messages.Get("HA_OVERCOMMIT_MESSAGE"), messages.Get("HA_NTOL_ZERO_QUERY")
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return HAConfig{}, xapi.Sanitize(err)
	}
	var current []string
	refs := map[string]xenapi.VMRef{}
	cfg.VMs = []HAPlanVM{}
	for ref, v := range vms {
		if !rules.HAProtectable(v.IsATemplate, v.IsASnapshot, v.IsControlDomain, v.OtherConfig["HideFromXenCenter"] == "true", v.NameLabel) {
			continue
		}
		current = append(current, v.HaRestartPriority)
		refs[v.UUID] = ref
		cfg.VMs = append(cfg.VMs, HAPlanVM{UUID: v.UUID, Name: v.NameLabel, Priority: v.HaRestartPriority, Order: v.Order, StartDelay: v.StartDelay})
	}
	// AssignPriorities.ProtectVmsByDefault: first-time defaults only when HA is being turned on
	firstTime := !p.HaEnabled && rules.HAFirstTime(current)
	for i := range cfg.VMs {
		v := &cfg.VMs[i]
		why, err := ops.NotAgileReason(s, refs[v.UUID])
		if err != nil {
			return HAConfig{}, err
		}
		v.NotAgile = why
		v.Priority = rules.HAInitialPriority(v.Priority, firstTime, why == "")
	}
	sort.Slice(cfg.VMs, func(i, j int) bool { // AssignPriorities: by name
		a, b := strings.ToLower(cfg.VMs[i].Name), strings.ToLower(cfg.VMs[j].Name)
		return a < b || a == b && cfg.VMs[i].UUID < cfg.VMs[j].UUID
	})
	return cfg, nil
}
