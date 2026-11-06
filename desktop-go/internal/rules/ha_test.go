package rules

import (
	"maps"
	"strings"
	"testing"
)

func TestHAPrioritiesMirrorTheCSharp(t *testing.T) {
	var names []string
	for _, p := range HAPriorities() {
		names = append(names, p.Value+"="+p.Text)
	}
	if got := strings.Join(names, ","); got != "restart=Restart,best-effort=Restart if possible,=Do not restart" {
		t.Errorf("priorities: %s", got)
	}
	for in, want := range map[string]string{"restart": "restart", "best-effort": "best-effort", "": "", "0": "restart", "1": "restart"} {
		if got := HANormalPriority(in); got != want {
			t.Errorf("HANormalPriority(%q) = %q, want %q", in, got, want)
		}
	}
	if HAValidPriority("0") || !HAValidPriority("") || !HAIsRestart("1") || HAIsRestart("best-effort") {
		t.Error("valid / restart")
	}
}

func TestHAPlanRules(t *testing.T) {
	if !HAProtectable(false, false, false, false, "web") || HAProtectable(true, false, false, false, "t") ||
		HAProtectable(false, true, false, false, "s") || HAProtectable(false, false, true, false, "dom0") ||
		HAProtectable(false, false, false, true, "hidden") || HAProtectable(false, false, false, false, "__gui__tmp") {
		t.Error("HaCanProtect")
	}
	if !HAFirstTime([]string{"", ""}) || HAFirstTime([]string{"", "best-effort"}) {
		t.Error("IsHaActivatedFirstTime")
	}
	// Allowed: only best effort or no restart for a VM that is not agile
	if !HAAllowed(HARestart, true) || HAAllowed(HARestart, false) || !HAAllowed(HABestEffort, false) || !HAAllowed(HADoNotRestart, false) {
		t.Error("Allowed")
	}
	cases := []struct {
		current          string
		firstTime, agile bool
		want             string
	}{
		{"", true, true, HARestart},             // the first time an agile VM is protected
		{"", true, false, HABestEffort},         // and one that is not agile is restarted if possible
		{"", false, true, HADoNotRestart},       // later each keeps its own
		{"restart", false, false, HABestEffort}, // no longer agile: dropped to best effort
		{"1", false, true, HARestart},
	}
	for _, c := range cases {
		if got := HAInitialPriority(c.current, c.firstTime, c.agile); got != c.want {
			t.Errorf("InitialPriority(%q, %v, %v) = %q, want %q", c.current, c.firstTime, c.agile, got, c.want)
		}
	}
	got := HAPlanForAPI(map[string]string{"a": "restart", "b": "best-effort", "c": "", "d": "0"})
	if !maps.Equal(got, map[string]string{"a": "restart", "d": "restart"}) {
		t.Errorf("GetVmHaRestartPrioritiesForApi: %v", got)
	}
}

func TestHANtolAndBlockers(t *testing.T) {
	if HAMaxTolerable(2, map[string]string{"hci-limit-fault-tolerance": "1"}) != 1 || HAMaxTolerable(2, map[string]string{"hci-limit-fault-tolerance": "3"}) != 2 ||
		HAMaxTolerable(2, map[string]string{"hci-limit-fault-tolerance": "x"}) != 2 || HAMaxTolerable(1, nil) != 1 {
		t.Error("GetMaximumTolerableHostFailures' HCI cap")
	}
	if HAPlanError(1, 1) != "" || HAPlanError(0, 2) != "" {
		t.Error("a plan that fits")
	}
	if !strings.HasPrefix(HAPlanError(2, 1), "HA cannot be guaranteed") {
		t.Errorf("overcommit: %q", HAPlanError(2, 1))
	}
	if HAPlanError(0, -1) == "" || HAPlanError(-1, 1) == "" {
		t.Error("unknown maximum / no value")
	}
	if !strings.Contains(HANtolInfo(2), "is 2.") || !strings.HasPrefix(HANtolInfo(0), "This pool cannot support") ||
		!strings.HasPrefix(HANtolInfo(-1), "The number of server failures that can be tolerated could not be determined") {
		t.Errorf("ntol info: %q / %q / %q", HANtolInfo(2), HANtolInfo(0), HANtolInfo(-1))
	}
	if HAConfigureBlocker(false, false, false, "p") != "" || HAConfigureBlocker(false, true, true, "p") != "" {
		t.Error("Configure can run")
	}
	if got := HAConfigureBlocker(false, true, false, "lab"); got != "Cannot configure HA for pool 'lab' because an HA Statefile VDI could not be found in the pool." {
		t.Errorf("no statefile: %q", got)
	}
	if HAConfigureBlocker(true, false, true, "p") == "" || HADisableBlocker(true, true) == "" {
		t.Error("an HA operation is running")
	}
	if HADisableBlocker(false, true) != "" || !strings.Contains(HADisableBlocker(false, false), "HA is not enabled") {
		t.Errorf("Disable: %q", HADisableBlocker(false, false))
	}
	if title, text := HADisableConfirmation("lab"); title != "Disable HA" || text != "Are you sure you want to disable HA for pool 'lab'?" {
		t.Errorf("Disable question: %q %q", title, text)
	}
	if HeartbeatUnsuitable("SR_HAS_NO_PBDS", "x") != "The SR is currently detached." ||
		HeartbeatUnsuitable("RBAC_PERMISSION_DENIED", "denied\ndetails") != "denied" || HeartbeatUnsuitable("OTHER", "why") != "why" {
		t.Error("heartbeat SR reasons")
	}
}
