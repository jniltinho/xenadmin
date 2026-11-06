package gui

import (
	"xcpng-center/internal/app"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// Vapps are a connection's vApps and the VMs that may join one (Manage vApps).
func (w *Workspace) Vapps(key string) (out app.Vapps, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadVapps(c); return err })
	return out, err
}

// VappDeleteQuestion is Delete's question for a vApp with vms VMs.
func (w *Workspace) VappDeleteQuestion(name string, vms int) string {
	return rules.VappDeleteQuestion(name, vms)
}

// VappChangeWarning is the question before VMs already in another vApp move to target ("" for none).
func (w *Workspace) VappChangeWarning(moves []rules.VappMove, target string) string {
	return rules.VappChangeWarning(moves, target)
}

// VappSummary is the New vApp wizard's Finish page.
func (w *Workspace) VappSummary(name string, vms []string) string {
	return rules.VappSummary(name, vms)
}

// CreateVapp creates a vApp with its VMs; its uuid.
func (w *Workspace) CreateVapp(key string, spec ops.VappSpec) (out string, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = ops.CreateVapp(c, spec); return err })
	w.changed(key)
	return out, err
}

// SaveVapp applies a vApp's Properties.
func (w *Workspace) SaveVapp(key, uuid string, spec ops.VappSpec) error {
	err := w.with(key, func(c *xapi.Connection) error { return ops.SaveVapp(c, uuid, spec) })
	w.changed(key)
	return err
}

// DeleteVapp deletes a vApp (its VMs stay).
func (w *Workspace) DeleteVapp(key, uuid string) error {
	err := w.with(key, func(c *xapi.Connection) error { return ops.DeleteVapp(c, uuid) })
	w.changed(key)
	return err
}

// VappPower starts ("start") or shuts down ("shutdown") a vApp.
func (w *Workspace) VappPower(key, uuid, op string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.VappPower(c, uuid, op) })
}

// AssignVapp puts a VM into a vApp, or takes it out when it is already there.
func (w *Workspace) AssignVapp(key, vappUUID, vmUUID string) error {
	err := w.with(key, func(c *xapi.Connection) error { return ops.AssignVapp(c, vappUUID, vmUUID) })
	w.changed(key)
	return err
}
