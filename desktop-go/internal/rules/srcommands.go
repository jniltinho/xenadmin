package rules

import (
	"regexp"
	"strconv"
	"strings"

	"xcpng-center/internal/messages"
)

// SRState is what the SR commands look at (SrRules.CanDetach / CanForget / CanDestroy / CanSetDefault).
type SRState struct {
	Type       string
	Detached   bool // no PBD currently attached (SR.IsDetached)
	HasPBDs    bool
	RunningVMs bool     // a running VM uses one of its disks (SR.HasRunningVMs; dom0 on metadata disks ignored)
	Allowed    []string // the SR's allowed_operations
	Default    bool
	VDICreate  bool // its SM driver can create disks (SR.SupportsVdiCreate)
	Shared     bool
	Hosts      int
	Broken     bool // SR.IsBroken: a PBD missing for a host, or not attached
	Hidden     bool // other_config HideFromXenCenter (SR.IsHidden)
	Degraded   bool // a multipathed PBD has fewer paths than it should (!SR.MultipathAOK)
}

// SRIcon is Images.GetIconFor(SR): no PBDs or hidden, broken (detached, or a multipath path down), the pool's
// default, or a working SR.
func SRIcon(sr SRState) string {
	switch {
	case !sr.HasPBDs || sr.Hidden:
		return "000_StorageDisabled_h32bit_16"
	case sr.Detached || sr.Broken || sr.Degraded:
		return "000_StorageBroken_h32bit_16"
	case sr.Default:
		return "000_StorageDefault_h32bit_16"
	}
	return "000_Storage_h32bit_16"
}

// creatable is SR.CanCreateWithXenCenter: the types the New Storage wizard knows.
var creatable = map[string]bool{"iso": true, "lvmoiscsi": true, "nfs": true, "lvmohba": true, "smb": true, "lvmofcoe": true, "gfs2": true}

// SRCommands are the SR commands that can run now, in the Storage menu's order: repair, set_default, detach, forget,
// destroy.
func SRCommands(sr SRState) []string {
	var out []string
	if sr.HasPBDs && (sr.Detached || sr.Broken || sr.Degraded) { // RepairSRCommand
		out = append(out, "repair")
	}
	if sr.HasPBDs && !sr.Default && sr.VDICreate && (sr.Shared || sr.Hosts <= 1) {
		out = append(out, "set_default")
	}
	if !sr.Detached && !sr.RunningVMs {
		out = append(out, "detach")
	}
	if !sr.RunningVMs && creatable[sr.Type] && contains(sr.Allowed, "forget") {
		out = append(out, "forget")
	}
	if !sr.RunningVMs && creatable[sr.Type] && contains(sr.Allowed, "destroy") {
		out = append(out, "destroy")
	}
	return out
}

// SRConfirmation is SrRules.Confirmation: the title, text and yes button of a command's question ("" for Set as
// Default, which asks nothing).
func SRConfirmation(command, name string) (title, text, yes string) {
	switch command {
	case "detach":
		return messages.Get("MESSAGEBOX_DETACH_SR_CONTINUE_TITLE"), messages.Get("MESSAGEBOX_DETACH_SR_CONTINUE"), "Yes"
	case "forget":
		return messages.Get("MESSAGEBOX_FORGET_SR_CONTINUE_TITLE"), messages.Format(messages.Get("MESSAGEBOX_FORGET_SR_CONTINUE"), name),
			strings.ReplaceAll(messages.Get("MESSAGEBOX_FORGET_SR_CONTINUE_YES_BUTTON_LABEL"), "&", "")
	case "destroy":
		return messages.Get("MESSAGEBOX_DESTROY_SR_CONTINUE_TITLE"), messages.Format(messages.Get("MESSAGEBOX_DESTROY_SR_CONTINUE"), name),
			strings.ReplaceAll(messages.Get("MESSAGEBOX_DESTROY_SR_YES_BUTTON_LABEL"), "&", "")
	}
	return "", "", ""
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

var pathCounts = regexp.MustCompile(`\[(\d+)L?,\s(\d+)L?`)

// MultipathDown is PBD.ParsePathCounts + SR.CheckMultipathString on a PBD's "mpath-..." status: fewer paths than the
// maximum.
func MultipathDown(status string) bool {
	m := pathCounts.FindStringSubmatch(status)
	if m == nil {
		return false
	}
	current, err1 := strconv.Atoi(m[1])
	most, err2 := strconv.Atoi(m[2])
	return err1 == nil && err2 == nil && current < most
}
