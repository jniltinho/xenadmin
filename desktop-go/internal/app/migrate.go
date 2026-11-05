package app

import (
	"errors"
	"sort"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// Target is a server of the Migrate to Server, Start on Server and Resume on Server submenus: its text with the reason
// it cannot take the VM, if any (VmCanBootOnHost); Home marks the first item, the VM's home server.
type Target struct {
	UUID    string `json:"uuid"`
	Text    string `json:"text"`
	Enabled bool   `json:"enabled"`
	Home    bool   `json:"home,omitempty"`
}

// ReadTargets is MigrateRules.Targets with HomeServerTarget first: "Home Server" (or why not), then every server by
// name. It asks VM.assert_can_boot_here of each server, as the C# app.
func ReadTargets(c *xapi.Connection, vmUUID string) ([]Target, error) {
	s := c.Session
	vm, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	rec, err := xenapi.VM.GetRecord(s, vm)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	// the reasons come from server failures: no session reference ever reaches the window
	reason := func(ref xenapi.HostRef) string { return xapi.Redact(cannotBootReason(s, vm, rec, hosts, ref)) }
	home := Target{Text: rules.TargetText(strings.ReplaceAll(messages.Get("HOME_SERVER_MENU_ITEM"), "&", ""), messages.Get("NO_HOME_SERVER")), Home: true}
	if h, ok := hosts[rec.Affinity]; ok {
		why := reason(rec.Affinity)
		home = Target{UUID: h.UUID, Text: rules.TargetText(strings.ReplaceAll(messages.Get("HOME_SERVER_MENU_ITEM"), "&", ""), why), Enabled: why == "", Home: true}
	}
	out := []Target{home}
	refs := make([]xenapi.HostRef, 0, len(hosts))
	for ref := range hosts {
		refs = append(refs, ref)
	}
	sort.Slice(refs, func(i, j int) bool {
		return strings.ToLower(hosts[refs[i]].NameLabel) < strings.ToLower(hosts[refs[j]].NameLabel)
	})
	for _, ref := range refs {
		why := reason(ref)
		out = append(out, Target{UUID: hosts[ref].UUID, Text: rules.TargetText(hosts[ref].NameLabel, why), Enabled: why == ""})
	}
	return out, nil
}

// cannotBootReason is MigrateRules.CannotBootReason ("" when the server can take the VM).
func cannotBootReason(s *xenapi.Session, vm xenapi.VMRef, rec xenapi.VMRecord, hosts map[xenapi.HostRef]xenapi.HostRecord, ref xenapi.HostRef) string {
	host := hosts[ref]
	if resident, ok := hosts[rec.ResidentOn]; ok && rec.PowerState == xenapi.VMPowerStateRunning {
		if ref == rec.ResidentOn {
			return messages.Get("HOST_MENU_CURRENT_SERVER")
		}
		if rules.CompareVersions(host.SoftwareVersion["product_version"], resident.SoftwareVersion["product_version"]) < 0 {
			return messages.Get("OLDER_THAN_CURRENT_SERVER")
		}
	}
	if rec.PowerState == xenapi.VMPowerStateRunning || rec.PowerState == xenapi.VMPowerStateSuspended {
		if v, ok := rec.LastBootCPUFlags["vendor"]; ok && host.CPUInfo["vendor"] != "" && v != host.CPUInfo["vendor"] {
			text, _ := messages.FailureText("VM_INCOMPATIBLE_WITH_THIS_HOST", nil)
			return text
		}
	}
	err := xenapi.VM.AssertCanBootHere(s, vm, ref)
	var api *xenapi.APIError
	switch {
	case err == nil:
		return ""
	case errors.As(err, &api):
		if api.Message == "VM_REQUIRES_SR" && len(api.Data) > 1 {
			if ct, err := xenapi.SR.GetContentType(s, xenapi.SRRef(api.Data[1])); err == nil && ct == "iso" {
				return messages.Get("MIGRATE_PLEASE_EJECT_YOUR_CD")
			}
		}
		_, short := messages.FailureText(api.Message, api.Data)
		return short
	}
	return messages.Get("HOST_MENU_UNKNOWN_ERROR")
}
