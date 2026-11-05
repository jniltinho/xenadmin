// Package rules holds the UI-free decisions of the C# app (XenAdmin.Desktop.Core): what an operation needs, when it
// asks for confirmation and with which words.
package rules

import "fmt"

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
)

// VMOps in the order of the C# toolbar.
var VMOps = []VMOp{VMStart, VMShutdown, VMReboot, VMSuspend, VMResume, VMForceShutdown, VMForceReboot}

// required is VmOperations.Required: the XAPI allowed_operations entry each operation needs.
var required = map[VMOp]string{
	VMStart: "start", VMShutdown: "clean_shutdown", VMReboot: "clean_reboot", VMSuspend: "suspend",
	VMResume: "resume", VMForceShutdown: "hard_shutdown", VMForceReboot: "hard_reboot",
}

// Known reports whether op is one of VMOps.
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
func NeedsConfirmation(op VMOp) bool { return op != VMStart && op != VMResume }

// Confirmation is VmOperations.Confirmation (the commands' ConfirmationDialogTitle/Text, XenModel/Messages.resx);
// haProtected is VM.HAIsProtected (pool HA on and a restart priority set).
func Confirmation(op VMOp, haProtected bool) (title, text string) {
	switch op {
	case VMShutdown:
		if haProtected {
			return "Shut Down VM", "The selected VM is currently protected by HA. Are you sure you want to shut it down?"
		}
		return "Shut Down VM", "Are you sure you want to shut down the selected VM?"
	case VMReboot:
		return "Reboot VM", "Are you sure you want to reboot the selected VM?"
	case VMSuspend:
		if haProtected {
			return "Suspend VM", "The selected VM is currently protected by HA. Are you sure you want to suspend it?"
		}
		return "Suspend VM", "Are you sure you want to suspend the selected VM?"
	case VMForceShutdown:
		return "Force Shut Down VM", "Are you sure you want to forcibly shut down the selected VM? This operation can result in data loss."
	case VMForceReboot:
		return "Force VM Reboot", "Are you sure you want to forcibly reboot the selected VM? This operation can result in data loss."
	}
	return "", ""
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
