package rules

import (
	"strings"
	"testing"
)

func TestHostRulesMirrorTheCSharp(t *testing.T) {
	h := HostInfo{Name: "h1", Coordinator: true, PoolHosts: 2, VMs: []HostVM{{Name: "dom0", ControlDomain: true, Running: true}}}
	title, text := HostConfirmation(h, false)
	if title != "Shut Down Server" || !strings.HasPrefix(text, "Are you sure you want to shut down 'h1'?") || !strings.Contains(text, "pool coordinator") {
		t.Errorf("shutdown: %q %q", title, text)
	}
	if _, text := HostConfirmation(HostInfo{Name: "h1", VMs: []HostVM{{HCIWarn: true}}}, true); !strings.Contains(text, "storage virtual appliance") {
		t.Errorf("HCI reboot: %q", text)
	}
	if CannotEnterMaintenance(h) == "" {
		t.Error("a coordinator with no member to hand over to")
	}
	h.NewCoordinator = "h2"
	if why := CannotEnterMaintenance(h); why != "" {
		t.Errorf("coordinator with a member: %q", why)
	}
	h.HAEnabled = true
	if !strings.Contains(CannotEnterMaintenance(h), "HA-enabled pool") {
		t.Error("HA coordinator")
	}
	h = HostInfo{Name: "h2", VMs: []HostVM{{Name: "local", Running: true, LocalStorage: true}, {Name: "b"}}}
	if !strings.Contains(CannotEnterMaintenance(h), "local") {
		t.Error("VM on local storage")
	}
	title, text = EvacuateConfirmation(h)
	if title != "Enter Maintenance Mode - h2" || !strings.Contains(text, "VMs to migrate:\nb\nlocal") {
		t.Errorf("evacuate: %q %q", title, text)
	}
	if ExitMaintenanceQuestion("h2", nil) != "" || !strings.Contains(ExitMaintenanceQuestion("h2", []string{"b"}), "VMs to restore:\nb") {
		t.Error("exit question")
	}
}
