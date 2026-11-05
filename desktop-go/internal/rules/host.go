package rules

import (
	"sort"
	"strings"

	"xcpng-center/internal/messages"
)

// HostVM is what the server commands need about a VM resident on the host.
type HostVM struct {
	Name          string
	Running       bool
	ControlDomain bool
	HCIWarn       bool // VM.HciWarnBeforeShutdown: other_config "hci-warn-before-shutdown" (a storage appliance)
	LocalStorage  bool // a disk or CD on storage that is not shared (VM.IsOnSharedStorage is not empty)
}

// HostInfo is what HostRules reads about a server and its pool.
type HostInfo struct {
	Name           string
	Coordinator    bool
	PoolHosts      int
	HAEnabled      bool
	VMs            []HostVM
	NewCoordinator string // the enabled, live member a coordinator hands over to ("" when none)
}

// HostConfirmation is HostRules.Confirm: ShutDownHostCommand / RebootHostCommand for one server.
func HostConfirmation(h HostInfo, reboot bool) (title, text string) {
	verb := map[bool]string{true: "REBOOT", false: "SHUTDOWN"}[reboot]
	hci, running := false, false
	for _, v := range h.VMs {
		hci = hci || v.HCIWarn
		running = running || (v.Running && !v.ControlDomain)
	}
	switch {
	case hci:
		text = messages.Get("CONFIRM_" + verb + "_HCI_WARN_SERVER")
	case running:
		text = messages.Get("CONFIRM_" + verb + "_SERVER")
	default:
		text = messages.Get("CONFIRM_" + verb + "_SERVER_NO_VMS")
	}
	text = messages.Format(text, h.Name)
	if !reboot && h.Coordinator && h.PoolHosts > 1 {
		text += "\n\n" + messages.Format(messages.Get("SHUT_DOWN_POOL_COORDINATOR_SINGLE"), h.Name)
	}
	return messages.Get("CONFIRM_" + verb + "_SERVER_TITLE"), text
}

// CannotEnterMaintenance is HostRules.CannotEnterMaintenance: why the server cannot be evacuated ("" when it can).
func CannotEnterMaintenance(h HostInfo) string {
	if h.HAEnabled && h.Coordinator {
		return messages.Format(messages.Get("HA_CANNOT_EVACUATE_COORDINATOR"), h.Name)
	}
	var stuck []string
	for _, v := range h.VMs {
		if !v.ControlDomain && v.Running && v.LocalStorage {
			stuck = append(stuck, v.Name)
		}
	}
	if len(stuck) > 0 {
		return "These VMs cannot be migrated off the server (their disks or CDs are on its local storage). " +
			"Shut them down, suspend them or move their storage first:\n\n" + strings.Join(stuck, "\n")
	}
	if h.Coordinator && h.PoolHosts > 1 && h.NewCoordinator == "" {
		return "No other enabled server of the pool can become the coordinator."
	}
	return ""
}

// EvacuateConfirmation is HostRules.EvacuateConfirmation (EvacuateHostDialog's text and VM list).
func EvacuateConfirmation(h HostInfo) (title, text string) {
	text = "This operation will migrate all VMs running on this server and transition it into maintenance mode. If a VM cannot be migrated you need to suspend it or shut it down."
	var vms []string
	for _, v := range h.VMs {
		if !v.ControlDomain {
			vms = append(vms, v.Name)
		}
	}
	sort.Strings(vms)
	if len(vms) == 0 {
		text += "\n\nNo VMs are running on this server."
	} else {
		text += "\n\nVMs to migrate:\n" + strings.Join(vms, "\n")
	}
	if h.Coordinator {
		text += "\n\nThis server is the pool coordinator. Entering maintenance mode will nominate a new coordinator for the pool. XCP-ng Center will temporarily lose its connection to the pool."
	}
	return messages.Format(messages.Get("EVACUATE_HOST_DIALOG_TITLE"), h.Name), text
}

// ExitMaintenanceQuestion is MainWindow.ToggleMaintenance's question when evacuated VMs can go back ("" when none).
func ExitMaintenanceQuestion(name string, restore []string) string {
	if len(restore) == 0 {
		return ""
	}
	return "Would you like to restore your VMs to their previous location on '" + name + "'?\n\nVMs to restore:\n" + strings.Join(restore, "\n")
}
