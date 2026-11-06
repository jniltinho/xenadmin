package gui

import (
	"xcpng-center/internal/app"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/xapi"
)

// GPU is the GPU tab of a pool or server.
func (w *Workspace) GPU(key, kind, uuid string) (out app.GPUTab, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadGPU(c, kind, uuid); return err })
	return out, err
}

// SetGPUPlacement sets the pool's GPU placement policy.
func (w *Workspace) SetGPUPlacement(key, algorithm string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.SetGPUPlacement(c, algorithm) })
}

// SetVGPUTypes applies the changed allowed vGPU types to GPUs.
func (w *Workspace) SetVGPUTypes(key string, pgpus []string, changes map[string]bool) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.SetVGPUTypes(c, pgpus, changes) })
}

// USB is the USB tab of a server.
func (w *Workspace) USB(key, uuid string) (out app.USBTab, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadUSB(c, uuid); return err })
	return out, err
}

// SetUSBPassthrough turns a device's pass-through over, if it is still as the user saw it.
func (w *Workspace) SetUSBPassthrough(key, pusb string, wasEnabled bool) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.SetUSBPassthrough(c, pusb, wasEnabled) })
}
