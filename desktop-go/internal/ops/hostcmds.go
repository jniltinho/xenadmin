package ops

import (
	"errors"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// liveHost reads a server and whether it is live (Host.IsLive: its metrics say so).
func liveHost(s *xenapi.Session, uuid string) (*host, bool, error) {
	h, err := readHost(s, uuid)
	if err != nil {
		return nil, false, xapi.Sanitize(err)
	}
	live, err := xenapi.HostMetrics.GetLive(s, h.rec.Metrics)
	return h, err == nil && live, nil
}

// RestartToolstackQuestion is HostRules.CanRestartToolstack and RestartToolstackConfirmation: title and text, or why not.
func RestartToolstackQuestion(c *xapi.Connection, uuid string) (title, text string, err error) {
	h, live, err := liveHost(c.Session, uuid)
	if err != nil {
		return "", "", err
	}
	if !live {
		return "", "", errors.New("the server is not live")
	}
	return messages.Get("CONFIRM_RESTART_TOOLSTACK_TITLE"),
		messages.Format(messages.Get("CONFIRM_RESTART_TOOLSTACK_ONE_SERVER"), Ellipsise(h.rec.NameLabel, 30), "XCP-ng Center"), nil
}

// RestartToolstack is RestartToolstackAction: host.restart_agent (the connection to that server drops for a moment).
func RestartToolstack(c *xapi.Connection, uuid string) error {
	s := c.Session
	h, live, err := liveHost(s, uuid)
	if err != nil {
		return err
	}
	if !live {
		return errors.New("the server is not live")
	}
	if len(h.rec.CurrentOperations) > 0 { // an evacuation or another server operation would be cut short
		return errors.New("the server is busy: wait for its current operation to finish")
	}
	return xapi.Sanitize(wait(s)(xenapi.Host.AsyncRestartAgent(s, h.ref)))
}

// removable is HostRules.CanRemoveFromPool: a live member (not the coordinator) of a pool, running at most its
// control domain.
func removable(h *host, live bool) string {
	switch {
	case h.info.PoolHosts < 2 || h.info.Coordinator:
		return "Only a member of a pool can be removed from it, not its coordinator."
	case !live:
		return "The server is not live."
	case len(h.rec.ResidentVMs) > 1:
		return "Shut down or migrate the VMs of the server first."
	}
	return ""
}

// RemoveFromPoolQuestion is RemoveFromPoolConfirmation: title, text and the Yes label (No is the default), or why not.
func RemoveFromPoolQuestion(c *xapi.Connection, uuid string) (title, text, yes string, err error) {
	h, live, err := liveHost(c.Session, uuid)
	if err != nil {
		return "", "", "", err
	}
	if why := removable(h, live); why != "" {
		return "", "", "", errors.New(why)
	}
	pool := h.prec.NameLabel // Helpers.GetName(pool): its name, else its coordinator's
	if pool == "" {
		pool, _ = xenapi.Host.GetNameLabel(c.Session, h.prec.Master)
	}
	name := Ellipsise(h.rec.NameLabel, 500)
	return messages.Get("MAINWINDOW_CONFIRM_REMOVE_FROM_POOL_TITLE"),
		messages.Format(messages.Get("MAINWINDOW_CONFIRM_REMOVE_FROM_POOL"), name, Ellipsise(pool, 500), name),
		strings.ReplaceAll(messages.Get("MAINWINDOW_CONFIRM_REMOVE_FROM_POOL_YES_BUTTON_LABEL"), "&", ""), nil
}

// RemoveFromPool is EjectHostFromPoolAction: pool.eject (the server reinitializes its local storage and restarts as a
// standalone server).
func RemoveFromPool(c *xapi.Connection, uuid string) error {
	s := c.Session
	h, live, err := liveHost(s, uuid)
	if err != nil {
		return err
	}
	if why := removable(h, live); why != "" {
		return errors.New(why)
	}
	return xapi.Sanitize(wait(s)(xenapi.Pool.AsyncEject(s, h.ref)))
}

// ConvertToTemplate is ConvertVMToTemplateCommand (VmRules.ConvertToTemplateAction): other_config "instant" (fast
// clones), then VM.set_is_a_template, for a VM that allows make_into_template.
func ConvertToTemplate(c *xapi.Connection, vmUUID string) error {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if vm.IsATemplate || len(vm.CurrentOperations) > 0 || !contains(vm.AllowedOperations, xenapi.VMOperationsMakeIntoTemplate) {
		return errors.New("the VM cannot be converted to a template now")
	}
	_ = xenapi.VM.RemoveFromOtherConfig(s, ref, "instant") // SetVMOtherConfigAction: remove, then add
	if err := xenapi.VM.AddToOtherConfig(s, ref, "instant", "true"); err != nil {
		return xapi.Sanitize(err)
	}
	return xapi.Sanitize(xenapi.VM.SetIsATemplate(s, ref, true))
}

// Ellipsise is StringExtensions.Ellipsise (rules.Ellipsise).
func Ellipsise(s string, n int) string { return rules.Ellipsise(s, n) }
