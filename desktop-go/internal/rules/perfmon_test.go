package rules

import (
	"reflect"
	"testing"
)

func TestPerfmonXML(t *testing.T) {
	defs := []Perfmon{{"cpu_usage", 0.5, 60, 3600}, {"network_usage", 102400, 120, 3600}}
	text := PerfmonXML(defs)
	want := `<config><variable><name value="cpu_usage" /><alarm_trigger_level value="0.5" /><alarm_trigger_period value="60" />` +
		`<alarm_auto_inhibit_period value="3600" /></variable><variable><name value="network_usage" /><alarm_trigger_level value="102400" />` +
		`<alarm_trigger_period value="120" /><alarm_auto_inhibit_period value="3600" /></variable></config>`
	if text != want {
		t.Fatal(text)
	}
	if got := ParsePerfmon(text); !reflect.DeepEqual(got, defs) {
		t.Fatal(got)
	}
	// a variable with no value that parses is dropped, as an unreadable blob
	if got := ParsePerfmon(`<config><variable><name value="x"/><alarm_trigger_level value="abc"/></variable></config>`); len(got) != 0 {
		t.Fatal(got)
	}
	if got := ParsePerfmon("<config"); got != nil {
		t.Fatal(got)
	}
}

func TestPerfmonAlerts(t *testing.T) {
	if len(PerfmonAlerts("hosts")) != 4 || len(PerfmonAlerts("srs")) != 1 || PerfmonAlerts("pools") != nil {
		t.Fatal("alerts per kind")
	}
}
