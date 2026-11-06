package rules

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// Performance alerts (XenModel/Alerts/PerfmonDefinition.cs, XenAdmin/SettingsPanels/PerfmonAlertEditPage.cs): the
// "perfmon" other_config key of a VM, server or SR holds <config><variable>...</variable></config>.

// PerfmonKey is PerfmonDefinition.PERFMON_KEY_NAME.
const PerfmonKey = "perfmon"

// Perfmon is a PerfmonDefinition, in the server's units.
type Perfmon struct {
	Name                   string
	Level, Period, Inhibit float64 // alarm_trigger_level, alarm_trigger_period (s), alarm_auto_inhibit_period (s)
}

// PerfmonAlert is one alert of the Alerts page: its definition name and how the page shows it (AlertGroup).
type PerfmonAlert struct {
	Name     string
	Factor   float64 // the page's value is the server's times Factor (XapiToGuiTriggerLevel)
	Min, Max float64 // the page's value range (the designer's NumericUpDown)
	Default  float64
	Sub      string // the Messages key of the summary (SubTextFormat)
}

const int64Max = 9223372036854775807 // the designer's {-1, 2147483647, 0, 0}: decimal Int64.MaxValue

// PerfmonAlerts are the alerts of the page per object kind (PerfmonAlertEditPage.SetXenObjects); mem_usage is kept in
// the control domain's other_config (PerfmonDefinitionAction).
func PerfmonAlerts(kind string) []PerfmonAlert {
	cpu := PerfmonAlert{"cpu_usage", 100, 5, 100, 50, "ALERT_CPUS_SUB_TEXT"}
	net := PerfmonAlert{"network_usage", 1.0 / 1024, 1, int64Max, 100, "ALERT_NET_SUB_TEXT"}
	switch kind {
	case "vms", "templates":
		return []PerfmonAlert{cpu, net, {"disk_usage", 1.0 / 1024, 1, int64Max, 1000, "ALERT_DISK_SUB_TEXT"}}
	case "hosts":
		return []PerfmonAlert{cpu, net, {"memory_free_kib", 1.0 / 1024, 1, int64Max, 1000, "ALERT_MEMORY_SUB_TEXT"},
			{"mem_usage", 100, 5, 100, 95, "ALERT_DOM0_MEMORY_SUB_TEXT"}}
	case "srs":
		return []PerfmonAlert{{"sr_io_throughput_total_per_host", 1024, 1, int64Max, 1000, "ALERT_SR_SUB_TEXT"}}
	}
	return nil
}

type perfmonValue struct {
	Value string `xml:"value,attr"`
}

type perfmonVariable struct {
	Name    *perfmonValue `xml:"name"`
	Level   *perfmonValue `xml:"alarm_trigger_level"`
	Period  *perfmonValue `xml:"alarm_trigger_period"`
	Inhibit *perfmonValue `xml:"alarm_auto_inhibit_period"`
}

// ParsePerfmon is PerfmonDefinition.GetPerfmonDefinitions: the variables with a name and at least one value that
// parses; anything unreadable gives none, as the C# app.
func ParsePerfmon(text string) []Perfmon {
	var doc struct {
		Variables []perfmonVariable `xml:"variable"`
	}
	if strings.TrimSpace(text) == "" || xml.Unmarshal([]byte(text), &doc) != nil {
		return nil
	}
	var out []Perfmon
	for _, v := range doc.Variables {
		if v.Name == nil {
			continue
		}
		p, set := Perfmon{Name: v.Name.Value}, false
		for _, f := range []struct {
			from *perfmonValue
			to   *float64
		}{{v.Level, &p.Level}, {v.Period, &p.Period}, {v.Inhibit, &p.Inhibit}} {
			if f.from == nil {
				continue
			}
			if n, err := strconv.ParseFloat(strings.TrimSpace(f.from.Value), 64); err == nil {
				*f.to, set = n, true
			}
		}
		if set {
			out = append(out, p)
		}
	}
	return out
}

// PerfmonXML is PerfmonDefinition.GetPerfmonDefinitionXML (XmlDocument.OuterXml).
func PerfmonXML(defs []Perfmon) string {
	var b strings.Builder
	b.WriteString("<config>")
	num := func(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }
	for _, d := range defs {
		var name strings.Builder
		_ = xml.EscapeText(&name, []byte(d.Name))
		b.WriteString(`<variable><name value="` + strings.ReplaceAll(name.String(), `"`, "&quot;") + `" /><alarm_trigger_level value="` +
			num(d.Level) + `" /><alarm_trigger_period value="` + num(d.Period) + `" /><alarm_auto_inhibit_period value="` +
			num(d.Inhibit) + `" /></variable>`)
	}
	b.WriteString("</config>")
	return b.String()
}
