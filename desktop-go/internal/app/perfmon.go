package app

import (
	"math"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// PerfmonPage is the Alerts page of a VM's, server's or SR's Properties (PerfmonAlertEditPage), in the page's units.
type PerfmonPage struct {
	Alerts   []PerfmonSetting `json:"alerts"`
	Interval float64          `json:"interval"` // the alert repeat interval, minutes
}

// PerfmonSetting is one alert of the page (AlertGroup).
type PerfmonSetting struct {
	Name    string  `json:"name"`
	Enabled bool    `json:"enabled"`
	Level   float64 `json:"level"`
	Period  float64 `json:"period"` // minutes
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Sub     string  `json:"sub"` // the summary's format: {0} the level, {1} the minutes
}

// ReadPerfmon reads the Alerts page of a VM, server or SR; a server's control domain memory alert is in its dom0's
// other_config (PerfmonAlertEditPage.Repopulate).
func ReadPerfmon(c *xapi.Connection, kind, uuid string) (PerfmonPage, error) {
	s := c.Session
	out := PerfmonPage{Alerts: []PerfmonSetting{}, Interval: 60}
	var other, dom0 map[string]string
	memMax := 0.0
	switch kind {
	case "vms", "templates":
		ref, err := xenapi.VM.GetByUUID(s, uuid)
		if err == nil {
			other, err = xenapi.VM.GetOtherConfig(s, ref)
		}
		if err != nil {
			return out, xapi.Sanitize(err)
		}
	case "hosts":
		ref, err := xenapi.Host.GetByUUID(s, uuid)
		if err != nil {
			return out, xapi.Sanitize(err)
		}
		h, err := xenapi.Host.GetRecord(s, ref)
		if err != nil {
			return out, xapi.Sanitize(err)
		}
		other = h.OtherConfig
		if d, err := xenapi.VM.GetOtherConfig(s, h.ControlDomain); err == nil {
			dom0 = d
		}
		if total, err := xenapi.HostMetrics.GetMemoryTotal(s, h.Metrics); err == nil {
			memMax = float64(total) / (1024 * 1024) // nudMemoryUsage.Maximum: the server's memory, MB
		}
	case "srs":
		ref, err := xenapi.SR.GetByUUID(s, uuid)
		if err == nil {
			other, err = xenapi.SR.GetOtherConfig(s, ref)
		}
		if err != nil {
			return out, xapi.Sanitize(err)
		}
	default:
		return out, nil
	}
	defs := map[string]rules.Perfmon{}
	var order []rules.Perfmon
	for _, d := range rules.ParsePerfmon(other[rules.PerfmonKey]) {
		if d.Name != "mem_usage" { // IsDom0MemoryUsage: read from the control domain only
			defs[d.Name], order = d, append(order, d)
		}
	}
	for _, d := range rules.ParsePerfmon(dom0[rules.PerfmonKey]) {
		if d.Name == "mem_usage" {
			defs[d.Name], order = d, append(order, d)
		}
	}
	for _, a := range rules.PerfmonAlerts(kind) {
		st := PerfmonSetting{Name: a.Name, Level: a.Default, Period: 1, Min: a.Min, Max: a.Max, Sub: messages.Get(a.Sub)}
		if a.Name == "memory_free_kib" && memMax > 0 {
			st.Max = math.Floor(memMax)
		}
		if d, ok := defs[a.Name]; ok { // AlertGroup.Populate
			st.Enabled = true
			st.Level = clamp(exact(d.Level*a.Factor), st.Min, st.Max) // kept as stored (no rounding), as the C# app
			st.Period = clamp(exact(d.Period/60), 1, 60)
		}
		out.Alerts = append(out.Alerts, st)
	}
	for _, d := range order { // every populated alert sets the shared interval; the last one wins
		for _, a := range rules.PerfmonAlerts(kind) {
			if a.Name == d.Name {
				if m := exact(d.Inhibit / 60); m >= 5 && m <= 86400 {
					out.Interval = m
				} else {
					out.Interval = 60 // AlertGroup.Populate: an interval the field refuses is 60 minutes
				}
			}
		}
	}
	return out, nil
}

func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }

// exact drops the float noise of a unit change (0.07 * 100 is 7.000000000000001), as decimal arithmetic would.
func exact(v float64) float64 { return math.Round(v*1e9) / 1e9 }
