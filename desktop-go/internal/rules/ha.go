package rules

import (
	"strconv"
	"strings"

	"xcpng-center/internal/messages"
)

// High availability: XenAdmin.Desktop.Core/HaRules.cs (HAConfigureCommand, HADisableCommand, the HA wizard's
// AssignPriorities and HaNtolControl) and the restart priorities of XenModel/XenAPI-Extensions/VM.cs.

// The restart priorities XAPI stores in VM.ha_restart_priority (VM.PriorityToString); the pre-Boston "0" and "1" are
// read as Restart (VM.StringToPriority gives AlwaysRestart*, which HA restarts).
const (
	HARestart      = "restart"
	HABestEffort   = "best-effort"
	HADoNotRestart = ""
)

// HAPriority is one restart priority with its name.
type HAPriority struct {
	Value string `json:"value"`
	Text  string `json:"text"`
}

// HAPriorities is VM.GetAvailableRestartPriorities with Helpers.RestartPriorityI18n: Restart, Restart if possible,
// Do not restart.
func HAPriorities() []HAPriority {
	return []HAPriority{
		{HARestart, messages.Friendly("Label-VM.ha_restart_priority.Restart")},
		{HABestEffort, messages.Friendly("Label-VM.ha_restart_priority.BestEffort")},
		{HADoNotRestart, messages.Friendly("Label-VM.ha_restart_priority.DoNotRestart")},
	}
}

// HANormalPriority is VM.HARestartPriority as one of HAPriorities' values ("0" and "1" of old pools are Restart).
func HANormalPriority(p string) string {
	switch p {
	case HABestEffort, HADoNotRestart:
		return p
	}
	return HARestart
}

// HAValidPriority: p is one of HAPriorities' values (what the app ever sets).
func HAValidPriority(p string) bool {
	return p == HARestart || p == HABestEffort || p == HADoNotRestart
}

// HAIsRestart is VM.HaPriorityIsRestart: the priority HA guarantees (and plans for).
func HAIsRestart(p string) bool { return HANormalPriority(p) == HARestart }

// HAProtectable is VM.HaCanProtect(false): a real VM (no template, snapshot or control domain) that the tree shows
// (VM.Show: no internal template, no "__gui__" temporary object, not hidden).
func HAProtectable(template, snapshot, controlDomain, hidden bool, name string) bool {
	return !template && !snapshot && !controlDomain && !hidden && !strings.HasPrefix(name, "__gui__")
}

// HAFirstTime is AssignPriorities.IsHaActivatedFirstTime: no VM has a restart priority yet.
func HAFirstTime(priorities []string) bool {
	for _, p := range priorities {
		if p != "" {
			return false
		}
	}
	return true
}

// HAAllowed is HaRules.Allowed: a VM that is not agile cannot be guaranteed a restart, only best effort or none.
func HAAllowed(p string, agile bool) bool { return agile || !HAIsRestart(p) }

// HAInitialPriority is HaRules.InitialPriority with the agility check of AssignPriorities.UpdateVMsAgility: the
// first time, Restart for an agile VM and best effort otherwise; later the VM's own, dropped to best effort when the
// VM is not agile.
func HAInitialPriority(current string, firstTime, agile bool) string {
	p := HANormalPriority(current)
	if firstTime {
		p = HABestEffort
		if agile {
			p = HARestart // VM.HaHighestProtectionAvailable
		}
	}
	if !HAAllowed(p, agile) {
		return HABestEffort
	}
	return p
}

// HAPlanForAPI is Helpers.GetVmHaRestartPrioritiesForApi: the VMs to restart, for
// pool.ha_compute_hypothetical_max_host_failures_to_tolerate (best effort and do not restart are left out: the
// server refuses them).
func HAPlanForAPI(priorities map[string]string) map[string]string {
	out := map[string]string{}
	for vm, p := range priorities {
		if HAIsRestart(p) {
			out[vm] = HARestart
		}
	}
	return out
}

// HAMaxTolerable is Pool.GetMaximumTolerableHostFailures' cap: the pool's other_config hci-limit-fault-tolerance
// (FAULT_TOLERANCE_LIMIT_FOR_HCI) when it is a number no greater than the computed maximum.
func HAMaxTolerable(max int, otherConfig map[string]string) int {
	if v, err := strconv.Atoi(otherConfig["hci-limit-fault-tolerance"]); err == nil && v <= max {
		return v
	}
	return max
}

// HAPlanError is HaRules.PlanError (AssignPriorities.updateButtons): "" when the failures to tolerate fit the plan.
func HAPlanError(ntol, max int) string {
	switch {
	case max < 0:
		return "The number of server failures that can be tolerated is still being computed."
	case ntol < 0:
		return "Enter the number of server failures to tolerate."
	case ntol > max:
		return messages.Get("HA_OVERCOMMIT_MESSAGE")
	}
	return ""
}

// HANtolInfo is HaNtolControl's text under the failures to tolerate (max < 0: unknown).
func HANtolInfo(max int) string {
	switch {
	case max < 0:
		return messages.Get("HA_NTOL_UNKNOWN")
	case max == 0:
		return messages.Get("HA_NTOL_ZERO")
	}
	return messages.Format(messages.Get("HA_NTOL_BLURB"), max)
}

// HAConfigureBlocker is HaRules.CanConfigure and ConfigureBlocker: why Configure cannot run ("" when it can): an HA
// operation running on the pool, or HA on without a statefile VDI to be found. pool is the name, ellipsised (30).
func HAConfigureBlocker(busy, haEnabled, statefileFound bool, pool string) string {
	switch {
	case busy:
		return "An HA operation is running on this pool: wait for it to end."
	case haEnabled && !statefileFound:
		return messages.Format(messages.Get("HA_CONFIGURE_NO_STATEFILE"), pool)
	}
	return ""
}

// HADisableBlocker is HaRules.CanDisable: HA on and no HA operation running ("" when Disable can run).
func HADisableBlocker(busy, haEnabled bool) string {
	switch {
	case busy:
		return "An HA operation is running on this pool: wait for it to end."
	case !haEnabled:
		text, _ := messages.FailureText("HA_NOT_ENABLED", nil)
		return text
	}
	return ""
}

// HADisableConfirmation is HADisableCommand's question (pool ellipsised to 30).
func HADisableConfirmation(pool string) (title, text string) {
	return messages.Get("DISABLE_HA"), messages.Format(messages.Get("HA_DISABLE_QUERY"), pool)
}

// HeartbeatUnsuitable is GetHeartbeatSRsAction's reason an SR cannot hold the heartbeat, from the failure of
// SR.assert_can_host_ha_statefile (its code and friendly text).
func HeartbeatUnsuitable(code, text string) string {
	switch code {
	case "SR_HAS_NO_PBDS":
		return messages.Get("SR_DETACHED")
	case "RBAC_PERMISSION_DENIED":
		first, _, _ := strings.Cut(text, "\n") // the first line, not the details
		return first
	}
	return text
}
