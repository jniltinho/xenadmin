package gui

import (
	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/xapi"
)

// Vtpms is the Manage vTPMs dialog of a VM.
func (w *Workspace) Vtpms(key, vmUUID string) (out ops.Vtpms, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = ops.ReadVtpms(c, vmUUID); return err })
	return out, err
}

// AddVtpm attaches a new vTPM to the VM.
func (w *Workspace) AddVtpm(key, vmUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.AddVtpm(c, vmUUID) })
}

// VtpmRemoveWarning is Remove vTPM's question.
func (w *Workspace) VtpmRemoveWarning() string { return messages.Get("VTPM_REMOVE_WARNING") }

// RemoveVtpm removes a vTPM of the VM.
func (w *Workspace) RemoveVtpm(key, vmUUID, vtpmUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.RemoveVtpm(c, vmUUID, vtpmUUID) })
}
