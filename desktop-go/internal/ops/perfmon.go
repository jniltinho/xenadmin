package ops

import (
	"errors"
	"fmt"
	"math"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// PerfmonChange is an alert left on in the Alerts page, in the page's units (minutes for the period).
type PerfmonChange struct {
	Name   string  `json:"name"`
	Level  float64 `json:"level"`
	Period float64 `json:"period"`
}

// SavePerfmon is PerfmonAlertEditPage.SaveSettings and PerfmonDefinitionAction: the alerts left on (every other one
// off) with the repeat interval in minutes, then perfmon's refresh on the servers concerned. Definitions this page
// does not show (fs_usage, physical_utilisation...) are kept.
func SavePerfmon(c *xapi.Connection, kind, uuid string, on []PerfmonChange, interval float64) error {
	s := c.Session
	if interval != math.Trunc(interval) || interval < 5 || interval > 86400 || int(interval)%5 != 0 {
		return errors.New(messages.Get("PERFORM_ALERT_EDIT_INTERVAL_WRONG_MULTIPLE"))
	}
	shown := map[string]rules.PerfmonAlert{}
	for _, a := range rules.PerfmonAlerts(kind) {
		shown[a.Name] = a
	}
	var page []rules.Perfmon
	var dom0Def *rules.Perfmon
	for _, ch := range on {
		a, ok := shown[ch.Name]
		if !ok {
			return fmt.Errorf("no %s alert for this object", ch.Name)
		}
		if ch.Level < a.Min || ch.Level > a.Max || ch.Period < 1 || ch.Period > 60 {
			return fmt.Errorf("the %s alert's values are out of range", ch.Name)
		}
		d := rules.Perfmon{Name: ch.Name, Level: ch.Level / a.Factor, Period: ch.Period * 60, Inhibit: interval * 60}
		if ch.Name == "mem_usage" {
			dom0Def = &d
			continue
		}
		page = append(page, d)
	}

	var (
		typ      string
		getOther func() (map[string]string, error)
		setOther func(xml string) error
		hosts    []xenapi.HostRef // perfmon's refresh (PerfmonDefinitionAction)
		dom0     xenapi.VMRef
	)
	switch kind {
	case "vms", "templates":
		ref, err := xenapi.VM.GetByUUID(s, uuid)
		if err != nil {
			return xapi.Sanitize(err)
		}
		vm, err := xenapi.VM.GetRecord(s, ref)
		if err != nil {
			return xapi.Sanitize(err)
		}
		typ = "vm"
		getOther = func() (map[string]string, error) { return xenapi.VM.GetOtherConfig(s, ref) }
		setOther = func(x string) error { return setVMOther(s, ref, x) }
		// shortcut: VM.Home is the running server or the affinity here (no local-storage server); the refresh is only
		// a hurry-up, perfmon reads the change within half an hour anyway
		for _, h := range []xenapi.HostRef{vm.ResidentOn, vm.Affinity} {
			if h != "" && h != "OpaqueRef:NULL" {
				hosts = append(hosts, h)
				break
			}
		}
	case "hosts":
		ref, err := xenapi.Host.GetByUUID(s, uuid)
		if err != nil {
			return xapi.Sanitize(err)
		}
		if dom0, err = xenapi.Host.GetControlDomain(s, ref); err != nil {
			return xapi.Sanitize(err)
		}
		typ, hosts = "host", []xenapi.HostRef{ref}
		getOther = func() (map[string]string, error) { return xenapi.Host.GetOtherConfig(s, ref) }
		setOther = func(x string) error {
			if err := xenapi.Host.RemoveFromOtherConfig(s, ref, rules.PerfmonKey); err != nil || x == "" {
				return err
			}
			return xenapi.Host.AddToOtherConfig(s, ref, rules.PerfmonKey, x)
		}
	case "srs":
		ref, err := xenapi.SR.GetByUUID(s, uuid)
		if err != nil {
			return xapi.Sanitize(err)
		}
		pbds, err := xenapi.SR.GetPBDs(s, ref)
		if err != nil {
			return xapi.Sanitize(err)
		}
		for _, p := range pbds {
			if h, err := xenapi.PBD.GetHost(s, p); err == nil {
				hosts = append(hosts, h)
			}
		}
		typ = "sr"
		getOther = func() (map[string]string, error) { return xenapi.SR.GetOtherConfig(s, ref) }
		setOther = func(x string) error {
			if err := xenapi.SR.RemoveFromOtherConfig(s, ref, rules.PerfmonKey); err != nil || x == "" {
				return err
			}
			return xenapi.SR.AddToOtherConfig(s, ref, rules.PerfmonKey, x)
		}
	default:
		return errors.New("this object has no alerts")
	}
	key := "/key:" + rules.PerfmonKey
	calls := []string{typ + ".add_to_other_config" + key, typ + ".remove_from_other_config" + key}
	if dom0 != "" {
		calls = append(calls, "vm.add_to_other_config"+key, "vm.remove_from_other_config"+key)
	}
	if err := permitted(s, calls...); err != nil {
		return err
	}

	if dom0 != "" { // the control domain memory alert lives in the dom0's other_config
		other, err := xenapi.VM.GetOtherConfig(s, dom0)
		if err != nil {
			return xapi.Sanitize(err)
		}
		defs, had := keepOthers(rules.ParsePerfmon(other[rules.PerfmonKey]), map[string]bool{"mem_usage": true})
		if dom0Def != nil {
			defs = append(defs, *dom0Def)
		}
		if dom0Def != nil || had {
			if err := setVMOther(s, dom0, rules.PerfmonXML(defs)); err != nil {
				return xapi.Sanitize(err)
			}
		}
	}
	other, err := getOther()
	if err != nil {
		return xapi.Sanitize(err)
	}
	names := map[string]bool{}
	for n := range shown {
		names[n] = true
	}
	defs, _ := keepOthers(rules.ParsePerfmon(other[rules.PerfmonKey]), names)
	defs = append(defs, page...)
	text := ""
	if len(defs) > 0 {
		text = rules.PerfmonXML(defs)
	}
	if err := setOther(text); err != nil {
		return xapi.Sanitize(err)
	}
	for _, h := range hosts {
		refreshPerfmon(s, h)
	}
	return nil
}

// keepOthers drops the definitions named in drop; had tells whether one was dropped.
func keepOthers(defs []rules.Perfmon, drop map[string]bool) (kept []rules.Perfmon, had bool) {
	for _, d := range defs {
		if drop[d.Name] {
			had = true
			continue
		}
		kept = append(kept, d)
	}
	return kept, had
}

// setVMOther is Helpers.SetOtherConfig for a VM's perfmon key ("" removes it).
func setVMOther(s *xenapi.Session, ref xenapi.VMRef, x string) error {
	if err := xenapi.VM.RemoveFromOtherConfig(s, ref, rules.PerfmonKey); err != nil || x == "" {
		return err
	}
	return xenapi.VM.AddToOtherConfig(s, ref, rules.PerfmonKey, x)
}

// refreshPerfmon is PerfmonDefinitionAction's refresh: perfmon re-reads its configuration now, started first when it
// is not running (ERROR_111); a failure only delays the alerts.
func refreshPerfmon(s *xenapi.Session, h xenapi.HostRef) {
	_, err := xenapi.Host.CallPlugin(s, h, "perfmon", "refresh", map[string]string{})
	if err != nil && strings.Contains(err.Error(), "ERROR_111") {
		if _, err = xenapi.Host.CallPlugin(s, h, "perfmon", "start", map[string]string{}); err == nil {
			_, _ = xenapi.Host.CallPlugin(s, h, "perfmon", "refresh", map[string]string{})
		}
	}
}
