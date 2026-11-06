// Package rules holds the UI-free decisions of the C# app (XenAdmin.Desktop.Core): what an operation needs, when it
// asks for confirmation and with which words.
package rules

import (
	"fmt"

	"xcpng-center/internal/messages"
)

// VMOp is a power operation of the VM menu and toolbar (VmOperations.VmOp; the names are the MCP vm_power ones).
type VMOp string

const (
	VMStart         VMOp = "start"
	VMShutdown      VMOp = "shutdown"
	VMReboot        VMOp = "reboot"
	VMSuspend       VMOp = "suspend"
	VMResume        VMOp = "resume"
	VMForceShutdown VMOp = "force_shutdown"
	VMForceReboot   VMOp = "force_reboot"
	VMPause         VMOp = "pause"   // the window's Pause button (PauseVMCommand); not an MCP operation
	VMUnpause       VMOp = "unpause" // UnPauseVMCommand
)

// VMOps in the order of the C# toolbar.
var VMOps = []VMOp{VMStart, VMShutdown, VMReboot, VMSuspend, VMResume, VMForceShutdown, VMForceReboot}

// required is VmOperations.Required: the XAPI allowed_operations entry each operation needs.
var required = map[VMOp]string{
	VMStart: "start", VMShutdown: "clean_shutdown", VMReboot: "clean_reboot", VMSuspend: "suspend",
	VMResume: "resume", VMForceShutdown: "hard_shutdown", VMForceReboot: "hard_reboot", VMPause: "pause", VMUnpause: "unpause",
}

// Known reports whether op is one of VMOps or Pause/Unpause.
func Known(op VMOp) bool { _, ok := required[op]; return ok }

// VMAllowed is VmOperations.IsAllowed: the server lists the operation in the VM's allowed_operations.
func VMAllowed(allowed []string, op VMOp) bool {
	for _, a := range allowed {
		if a == required[op] {
			return true
		}
	}
	return false
}

// NeedsConfirmation is VmOperations.NeedsConfirmation: Start and Resume do not ask.
func NeedsConfirmation(op VMOp) bool {
	return op != VMStart && op != VMResume && op != VMPause && op != VMUnpause // PauseVMCommand.ConfirmationRequired: false
}

// Confirmation is VmOperations.Confirmation (the commands' ConfirmationDialogTitle/Text from Messages.resx);
// haProtected is VM.HAIsProtected (pool HA on and a restart priority set).
func Confirmation(op VMOp, haProtected bool) (title, text string) {
	key := map[VMOp]string{VMShutdown: "SHUTDOWN_VM", VMReboot: "REBOOT_VM", VMSuspend: "SUSPEND_VM",
		VMForceShutdown: "FORCESHUTDOWN_VM", VMForceReboot: "FORCEREBOOT_VM"}[op]
	if key == "" {
		return "", ""
	}
	text = messages.Get("CONFIRM_" + key)
	switch {
	case haProtected && (op == VMShutdown || op == VMSuspend):
		text = messages.Get("HA_CONFIRM_" + key)
	case op == VMForceShutdown || op == VMForceReboot:
		text = messages.Get("CONFIRM_" + key + "_NO_CANCEL_TASKS")
	}
	return messages.Get("CONFIRM_" + key + "_TITLE"), text
}

// HAStartBlocker is VmOperations.HaStartBlocker: with HA on, a VM whose restart priority protects it must be agile; a
// disk on storage that is not shared makes it not agile, so Start/Resume would fail (the app never unprotects it).
func HAStartBlocker(haEnabled, restartPriority bool, disksShared []bool, name string) string {
	if !haEnabled || !restartPriority {
		return ""
	}
	for _, shared := range disksShared {
		if !shared {
			return fmt.Sprintf("'%s' is protected by HA but has disks on local storage, so it is not agile. Set its HA restart priority to \"Restart if possible\" or \"Do not restart\" (Pool > High Availability) before starting it.", name)
		}
	}
	return ""
}
