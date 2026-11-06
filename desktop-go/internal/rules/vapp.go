package rules

import (
	"strings"

	"xcpng-center/internal/messages"
)

// Ellipsise is StringExtensions.Ellipsise: a string of n characters or more cut to n-3 and "...".
func Ellipsise(s string, n int) string {
	r := []rune(s)
	if len(r) < n {
		return s
	}
	return string(r[:n-3]) + "..."
}

// VappDeleteQuestion is VappRules.DeleteConfirmation for one vApp with vms VMs.
func VappDeleteQuestion(name string, vms int) string {
	switch vms {
	case 0:
		return messages.Format(messages.Get("CONFIRM_DELETE_VM_APPLIANCE_0"), Ellipsise(name, 120))
	case 1:
		return messages.Format(messages.Get("CONFIRM_DELETE_VM_APPLIANCE_1"), Ellipsise(name, 120), vms)
	}
	return messages.Format(messages.Get("CONFIRM_DELETE_VM_APPLIANCE"), Ellipsise(name, 120), vms)
}

// VappMove is a VM going into a vApp: its name and the vApp it is in now ("" for none).
type VappMove struct {
	VM   string `json:"vm"`
	From string `json:"from"`
}

// VappChangeWarning is VappRules.ChangeWarning: the question when VMs already belong to another vApp ("" for none).
// The callers pass only the VMs whose vApp is not the target (compared by reference: two vApps may share a name).
func VappChangeWarning(moves []VappMove, target string) string {
	var moving []VappMove
	for _, m := range moves {
		if m.From != "" {
			moving = append(moving, m)
		}
	}
	switch len(moving) {
	case 0:
		return ""
	case 1:
		return messages.Format(messages.Get("CHANGE_VMS_APPLIANCE_WARNING"), Ellipsise(moving[0].VM, 250), Ellipsise(moving[0].From, 250), Ellipsise(target, 250))
	}
	return messages.Format(messages.Get("CHANGE_VMS_APPLIANCES_WARNING"), Ellipsise(target, 250))
}

// VappSummary is NewVMApplianceWizard.GetSummary.
func VappSummary(name string, vms []string) string {
	text := strings.NewReplacer(`\n`, "\n", `\r`, "\r").Replace(messages.Get("VM_APPLIANCE_SUMMARY"))
	return messages.Format(text, name, strings.Join(vms, ", "))
}
