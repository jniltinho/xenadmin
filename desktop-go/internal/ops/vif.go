package ops

import (
	"errors"
	"strconv"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// VIFSettings are the Add Virtual Interface / Virtual Interface Properties dialog's values (VIFDialog).
type VIFSettings struct {
	Network string `json:"network"` // network uuid
	MAC     string `json:"mac"`     // "" to auto-generate
	QoS     bool   `json:"qos"`
	Limit   string `json:"limit"` // Kbytes/s
}

// vifRecord is NetworkRules.Settings: the new VIF, or a copy of the existing one with the dialog's values.
func vifRecord(base xenapi.VIFRecord, vm xenapi.VMRef, net xenapi.NetworkRef, device int, set VIFSettings) xenapi.VIFRecord {
	v := base
	v.VM, v.Network, v.MAC, v.Device = vm, net, set.MAC, strconv.Itoa(device)
	if v.OtherConfig == nil {
		v.OtherConfig = map[string]string{}
	}
	if set.QoS {
		v.QosAlgorithmType = "ratelimit"
	} else if base.QosAlgorithmType == "ratelimit" {
		v.QosAlgorithmType = ""
	}
	if set.Limit != "" {
		v.QosAlgorithmParams = map[string]string{"kbps": set.Limit}
	} else if v.QosAlgorithmParams == nil {
		v.QosAlgorithmParams = map[string]string{}
	}
	return v
}

// vifChanged is NetworkRules.Changed: an edit that changes nothing is a Cancel.
func vifChanged(old xenapi.VIFRecord, net xenapi.NetworkRef, set VIFSettings) bool {
	if old.Network != net || old.MAC != set.MAC {
		return true
	}
	if old.QosAlgorithmType == "ratelimit" {
		return !set.QoS || old.QosAlgorithmParams["kbps"] != set.Limit
	}
	return old.QosAlgorithmType == "" && set.QoS
}

// checkVIF validates the dialog's values (rules.ValidateVIF) and finds the network.
func checkVIF(s *xenapi.Session, set VIFSettings) (xenapi.NetworkRef, error) {
	set.Limit = strings.TrimSpace(set.Limit)
	if why, ok := rules.ValidateVIF(set.Network != "", set.MAC != "", set.MAC, set.QoS, set.Limit); !ok {
		if why == "" {
			why = "Complete the interface's settings."
		}
		return "", errors.New(why)
	}
	net, err := xenapi.Network.GetByUUID(s, set.Network)
	return net, xapi.Sanitize(err)
}

// CreateVIF is CreateVIFAction (after AddNetworkButton_Click's limit): VIF.create on the next free device, plugged at
// once into a running VM when allowed; rebootRequired when it waits for the next boot.
func CreateVIF(c *xapi.Connection, vmUUID string, set VIFSettings) (rebootRequired bool, err error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return false, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return false, xapi.Sanitize(err)
	}
	if len(vm.CurrentOperations) > 0 { // NetworkRules.CanAdd
		return false, errors.New("the VM is busy")
	}
	if int64(len(vm.VIFs)) >= rules.Restriction(vm.Recommendations, "number-of-vifs", "max", 7) { // NetworkRules.AtMaxInterfaces
		text, _ := messages.FailureText("VIFS_MAX_ALLOWED", nil)
		return false, errors.New(text)
	}
	net, err := checkVIF(s, set)
	if err != nil {
		return false, err
	}
	var used []string
	for _, v := range vm.VIFs {
		if d, err := xenapi.VIF.GetDevice(s, v); err == nil {
			used = append(used, d)
		}
	}
	rec := vifRecord(xenapi.VIFRecord{MTU: 1500, LockingMode: xenapi.VifLockingModeNetworkDefault}, ref, net, rules.VIFDevice(used), set)
	if err := permitted(s, vifCalls(vm.PowerState, false)...); err != nil {
		return false, err
	}
	return createVIF(s, vm.PowerState, rec)
}

// vifCalls are the API calls of CreateVIFAction (and, with del, of UpdateVIFAction's DeleteVIFAction first), as their
// ApiMethodsToRoleCheck.
func vifCalls(state xenapi.VMPowerState, del bool) []string {
	calls := []string{"vif.create"}
	if state == xenapi.VMPowerStateRunning {
		calls = append(calls, "vif.plug")
	}
	if del {
		calls = append(calls, "vif.destroy")
		if state == xenapi.VMPowerStateRunning {
			calls = append(calls, "vif.unplug")
		}
	}
	return calls
}

func createVIF(s *xenapi.Session, state xenapi.VMPowerState, rec xenapi.VIFRecord) (bool, error) {
	r, err := waitRef(s)(xenapi.VIF.AsyncCreate(s, rec))
	if err != nil {
		return false, xapi.Sanitize(err)
	}
	vif := xenapi.VIFRef(r)
	if state != xenapi.VMPowerStateRunning {
		return false, nil
	}
	allowed, err := xenapi.VIF.GetAllowedOperations(s, vif)
	if err != nil || !contains(allowed, xenapi.VifOperationsPlug) {
		return true, nil
	}
	return false, xapi.Sanitize(xenapi.VIF.Plug(s, vif)) // try hot-plugging
}

// UpdateVIF is UpdateVIFAction: the VIF is deleted and created again on its device with the new values (nothing is done
// when nothing changed).
func UpdateVIF(c *xapi.Connection, vifUUID string, set VIFSettings) (rebootRequired bool, err error) {
	s := c.Session
	ref, err := xenapi.VIF.GetByUUID(s, vifUUID)
	if err != nil {
		return false, xapi.Sanitize(err)
	}
	old, err := xenapi.VIF.GetRecord(s, ref)
	if err != nil {
		return false, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, old.VM)
	if err != nil {
		return false, xapi.Sanitize(err)
	}
	if !canEditVIF(vm, old) {
		return false, errors.New("the interface cannot be changed now")
	}
	net, err := checkVIF(s, set)
	if err != nil || !vifChanged(old, net, set) {
		return false, err
	}
	device, err := strconv.Atoi(old.Device)
	if err != nil {
		return false, errors.New("the interface has no device number")
	}
	if err := permitted(s, vifCalls(vm.PowerState, true)...); err != nil { // never delete what cannot be created again
		return false, err
	}
	if err := deleteVIF(s, vm.PowerState, ref); err != nil {
		return false, err
	}
	return createVIF(s, vm.PowerState, vifRecord(old, old.VM, net, device, set))
}

// canEditVIF is NetworkRules.CanEdit (and CanRemove): the VM is not busy and the VIF can be unplugged or is not plugged.
func canEditVIF(vm xenapi.VMRecord, vif xenapi.VIFRecord) bool {
	return len(vm.CurrentOperations) == 0 && (contains(vif.AllowedOperations, xenapi.VifOperationsUnplug) || !vif.CurrentlyAttached)
}

// DeleteVIF is DeleteVIFAction (after the MESSAGEBOX_VIF_DELETE question).
func DeleteVIF(c *xapi.Connection, vifUUID string) error {
	s := c.Session
	ref, err := xenapi.VIF.GetByUUID(s, vifUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vif, err := xenapi.VIF.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, vif.VM)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if !canEditVIF(vm, vif) {
		return errors.New("the interface cannot be removed now")
	}
	calls := []string{"vif.destroy"}
	if vm.PowerState == xenapi.VMPowerStateRunning {
		calls = append(calls, "vif.unplug")
	}
	if err := permitted(s, calls...); err != nil {
		return err
	}
	return deleteVIF(s, vm.PowerState, ref)
}

// deleteVIF unplugs the VIF of a running VM when allowed (already detached is fine) and destroys it.
func deleteVIF(s *xenapi.Session, state xenapi.VMPowerState, ref xenapi.VIFRef) error {
	if state == xenapi.VMPowerStateRunning {
		if allowed, err := xenapi.VIF.GetAllowedOperations(s, ref); err == nil && contains(allowed, xenapi.VifOperationsUnplug) {
			var api *xenapi.APIError
			if err := xenapi.VIF.Unplug(s, ref); err != nil && !(errors.As(err, &api) && api.Message == "DEVICE_ALREADY_DETACHED") {
				return xapi.Sanitize(err)
			}
		}
	}
	return xapi.Sanitize(xenapi.VIF.Destroy(s, ref))
}

// SetVIFActive is PlugVIFAction / UnplugVIFAction: only a running VM's interface, when the server allows it.
func SetVIFActive(c *xapi.Connection, vifUUID string, active bool) error {
	s := c.Session
	ref, err := xenapi.VIF.GetByUUID(s, vifUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vif, err := xenapi.VIF.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	state, err := xenapi.VM.GetPowerState(s, vif.VM)
	if err != nil {
		return xapi.Sanitize(err)
	}
	op := map[bool]xenapi.VifOperations{true: xenapi.VifOperationsPlug, false: xenapi.VifOperationsUnplug}[active]
	if state != xenapi.VMPowerStateRunning || !contains(vif.AllowedOperations, op) {
		return nil
	}
	if active {
		return xapi.Sanitize(xenapi.VIF.Plug(s, ref))
	}
	return xapi.Sanitize(xenapi.VIF.Unplug(s, ref))
}
