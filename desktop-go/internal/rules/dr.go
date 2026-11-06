package rules

import (
	"fmt"

	"xcpng-center/internal/messages"
)

// Disaster recovery: XenAdmin.Desktop.Core/DrRules.cs (DRConfigureDialog, DRFailoverWizard) and the XenModel actions
// under XenModel/Actions/DR.

// DRMaxSRs is DRConfigureDialog.MAX_SR_SELECTED.
const DRMaxSRs = 8

// The three runs of the DR wizard (DrType, DRFailoverWizardFirstPage).
const (
	DRFailover = "failover"
	DRFailback = "failback"
	DRTest     = "test"
)

// DRKnownType: one of DRFailover, DRFailback, DRTest.
func DRKnownType(t string) bool { return t == DRFailover || t == DRFailback || t == DRTest }

// What to do with the recovered VMs and vApps (DrRules.StartMode).
const (
	DRStartNone   = "none"
	DRStart       = "start"
	DRStartPaused = "paused"
)

// DRStartModes is DrWizard.ObjectsPage's "After recovery" choices: a test failover never boots the recovered guests
// (paused or not started).
func DRStartModes(t string) []string {
	if t == DRTest {
		return []string{DRStartPaused, DRStartNone}
	}
	return []string{DRStart, DRStartPaused, DRStartNone}
}

// DRConfigureWarning is DrRules.ConfigureWarning (DRConfigureDialog.ToggleWarningsVisibleState): the warning for this
// choice ("" for none) and whether OK is allowed.
func DRConfigureWarning(originallyEnabled bool, checked int) (string, bool) {
	switch {
	case checked > DRMaxSRs:
		return messages.Format(messages.Get("DR_CONFIGURE_TOO_MANY_SRS"), DRMaxSRs), false
	case originallyEnabled && checked == 0:
		return "Because no VM configuration data will be stored anywhere in the pool, it will not be possible to recover VMs or vApp from this pool in the event of a disaster.", true
	case !originallyEnabled && checked > 0:
		return "DR configuration is not complete until storage mirroring has been set up and the selected SRs are being replicated.", true
	}
	return "", true
}

// DRMetadataSR is GetMetadataVDIsAction's filter: metadata VDIs are looked for on shared SRs only, not on the tools SR
// nor an ISO library.
func DRMetadataSR(shared, tools bool, contentType string) bool {
	return shared && !tools && contentType != "iso"
}

// DRRealVM is VM.IsRealVm: what VdiLoadMetadataAction keeps from a pool's metadata (no template, snapshot or control
// domain).
func DRRealVM(template, snapshot, controlDomain bool) bool {
	return !template && !snapshot && !controlDomain
}

// DRCheck is one precheck of the DR wizard (DRFailoverWizardPrecheckPage): an error blocks the recovery, a warning
// does not.
type DRCheck struct {
	Text  string `json:"text"`
	Error bool   `json:"error"`
}

// DRHACheck is the test failover's HA check: HA must be off on this pool ("" when it can go on).
func DRHACheck(t string, haEnabled bool, pool string) string {
	if t == DRTest && haEnabled {
		return fmt.Sprintf("HA is enabled on '%s'. Disable HA before you test a failover on this pool.", pool)
	}
	return ""
}

// DRRunningCheck is RunningVmCheck: a VM still running in its source pool blocks a failover or failback.
func DRRunningCheck(vm, pool string) DRCheck {
	return DRCheck{fmt.Sprintf("The VM '%s' is still running in its pool '%s': shut it down there first.", vm, pool), true}
}

// DRExistingCheck is DrRules.Prechecks for a VM of the metadata that already exists in this pool (by uuid): it must be
// older than the metadata's and halted; then it is replaced (a warning). A test failover refuses it too: its clean-up
// destroys what it recovered, and must never take a VM that was here before (a safety the C# app lacks).
func DRExistingCheck(name string, existingVersion, metadataVersion int, halted, test bool) DRCheck {
	switch {
	case existingVersion >= metadataVersion:
		return DRCheck{fmt.Sprintf("The VM '%s' already exists in this pool and is not older than the one in the metadata: delete it first.", name), true}
	case !halted:
		return DRCheck{fmt.Sprintf("The VM '%s' already exists in this pool and is running: shut it down first.", name), true}
	case test:
		return DRCheck{fmt.Sprintf("The VM '%s' already exists in this pool: a test failover does not replace it.", name), true}
	}
	return DRCheck{fmt.Sprintf("The VM '%s' already exists in this pool and will be replaced.", name), false}
}

// DRTestCollision is a test failover's refusal of a vApp or snapshot (kind) of the metadata that this pool has already.
func DRTestCollision(kind, name string) DRCheck {
	return DRCheck{fmt.Sprintf("The %s '%s' already exists in this pool: a test failover does not replace it.", kind, name), true}
}

// DRBlocked: one of the checks is an error (DrWizard.Validate: fix the errors before you continue).
func DRBlocked(checks []DRCheck) bool {
	for _, c := range checks {
		if c.Error {
			return true
		}
	}
	return false
}
