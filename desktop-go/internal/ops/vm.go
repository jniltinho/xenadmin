// Package ops runs the operations on a server (the XenModel actions the C# app reuses).
package ops

import (
	"errors"
	"fmt"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

const nullRef = "OpaqueRef:NULL"

// VMPower is VmOperations.Create + the action: it refuses what is not a real VM or what the VM does not allow now,
// and Start/Resume of an HA-protected VM that is not agile (never unprotected behind the user's back).
// The call runs as a task polled to its end (wait). shortcut: no progress or cancel shown yet (4.2).
func VMPower(c *xapi.Connection, uuid string, op rules.VMOp) error {
	if !rules.Known(op) {
		return fmt.Errorf("unknown operation %q", op)
	}
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, uuid)
	if err != nil {
		return xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if vm.IsATemplate || vm.IsASnapshot || vm.IsControlDomain {
		return errors.New("not a VM (template, snapshot or control domain)")
	}
	allowed := make([]string, len(vm.AllowedOperations))
	for i, a := range vm.AllowedOperations {
		allowed[i] = string(a)
	}
	if !rules.VMAllowed(allowed, op) {
		return fmt.Errorf("'%s' is not allowed for this VM now (power state %s)", op, vm.PowerState)
	}
	if op == rules.VMStart || op == rules.VMResume {
		if why, err := haBlocker(s, vm); err != nil {
			return err
		} else if why != "" {
			return errors.New(why)
		}
	}
	var task xenapi.TaskRef
	switch op {
	case rules.VMStart:
		task, err = xenapi.VM.AsyncStart(s, ref, false, false)
	case rules.VMResume:
		task, err = xenapi.VM.AsyncResume(s, ref, false, false)
	case rules.VMShutdown:
		task, err = xenapi.VM.AsyncCleanShutdown(s, ref)
	case rules.VMReboot:
		task, err = xenapi.VM.AsyncCleanReboot(s, ref)
	case rules.VMSuspend:
		task, err = xenapi.VM.AsyncSuspend(s, ref)
	case rules.VMForceShutdown:
		task, err = xenapi.VM.AsyncHardShutdown(s, ref)
	case rules.VMForceReboot:
		task, err = xenapi.VM.AsyncHardReboot(s, ref)
	}
	return xapi.Sanitize(wait(s)(task, err))
}

// HasVM reports whether the connection's pool has a VM (or template) with this UUID.
func HasVM(c *xapi.Connection, uuid string) (bool, error) {
	_, err := xenapi.VM.GetByUUID(c.Session, uuid)
	return found(err)
}

// HasHost reports whether the connection's pool has a server with this UUID.
func HasHost(c *xapi.Connection, uuid string) (bool, error) {
	_, err := xenapi.Host.GetByUUID(c.Session, uuid)
	return found(err)
}

func found(err error) (bool, error) {
	var api *xenapi.APIError
	if errors.As(err, &api) && api.Message == "UUID_INVALID" { // JSON-RPC: the XAPI error code is the message
		return false, nil
	}
	return err == nil, xapi.Sanitize(err)
}

// haBlocker gathers what rules.HAStartBlocker needs: pool HA, the restart priority, the disks' SRs being shared.
func haBlocker(s *xenapi.Session, vm xenapi.VMRecord) (string, error) {
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	ha := false
	for _, p := range pools {
		ha = p.HaEnabled
	}
	if !ha || vm.HaRestartPriority != "restart" { // VM.HaPriorityIsRestart
		return "", nil
	}
	var shared []bool
	for _, b := range vm.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil {
			return "", xapi.Sanitize(err)
		}
		if vbd.Type != xenapi.VbdTypeDisk || vbd.VDI == nullRef {
			continue
		}
		vdi, err := xenapi.VDI.GetRecord(s, vbd.VDI)
		if err != nil {
			return "", xapi.Sanitize(err)
		}
		sr, err := xenapi.SR.GetRecord(s, vdi.SR)
		if err != nil {
			return "", xapi.Sanitize(err)
		}
		shared = append(shared, sr.Shared)
	}
	return rules.HAStartBlocker(true, true, shared, vm.NameLabel), nil
}
