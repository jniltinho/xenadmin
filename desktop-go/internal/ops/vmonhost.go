package ops

import (
	"errors"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/xapi"
)

// hostOps are the VM operations that take a server (VMOperationHostCommand), with the allowed operation each needs.
var hostOps = map[string]xenapi.VMOperations{"migrate": xenapi.VMOperationsPoolMigrate, "start_on": xenapi.VMOperationsStartOn, "resume_on": xenapi.VMOperationsResumeOn}

// VMOnHost runs Migrate to Server (VMMigrateAction: a live VM.pool_migrate), Start on Server (VM.start_on) or Resume on
// Server (VM.resume_on), if the VM allows it now.
func VMOnHost(c *xapi.Connection, vmUUID, hostUUID, op string) error {
	s := c.Session
	need, ok := hostOps[op]
	if !ok {
		return errors.New("unknown operation")
	}
	vm, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	rec, err := xenapi.VM.GetRecord(s, vm)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if rec.IsATemplate || rec.IsASnapshot || rec.IsControlDomain || !contains(rec.AllowedOperations, need) {
		return errors.New("the VM does not allow this now")
	}
	host, err := xenapi.Host.GetByUUID(s, hostUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	switch op {
	case "migrate":
		err = wait(s)(xenapi.VM.AsyncPoolMigrate(s, vm, host, map[string]string{"live": "true"}))
		var api *xenapi.APIError
		if errors.As(err, &api) && api.Message == "VM_MIGRATE_FAILED" && len(api.Data) >= 4 && strings.Contains(api.Data[3], "VDI_MISSING") {
			return errors.New(messages.Get("MIGRATE_EJECT_TOOLS_ON_UPGRADE"))
		}
	case "start_on":
		err = wait(s)(xenapi.VM.AsyncStartOn(s, vm, host, false, false))
	case "resume_on":
		err = wait(s)(xenapi.VM.AsyncResumeOn(s, vm, host, false, false))
	}
	return xapi.Sanitize(err)
}
