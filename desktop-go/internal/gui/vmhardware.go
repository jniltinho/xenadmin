package gui

import (
	"xcpng-center/internal/app"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// A VM's GPU and USB pages of its Properties, and the New Network wizard's SR-IOV network.

// VMGPU is a VM's GPU page.
func (w *Workspace) VMGPU(key, vm string) (out app.VMGPUPage, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadVMGPU(c, vm); return err })
	return out, err
}

// VMGPUState is the GPU page's buttons and warnings for its rows (GpuEditPage.ShowHideWarnings).
func (w *Workspace) VMGPUState(in rules.GPUPageInput) rules.GPUPageState { return rules.GPUPage(in) }

// VGPUChoices is the Add virtual GPU dialog's list beside the page's vGPUs (their type uuids).
func (w *Workspace) VGPUChoices(key, vm string, existing []string) (out app.VGPUPicks, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadVGPUChoices(c, vm, existing); return err })
	return out, err
}

// AssignVGPUs saves the GPU page: the vGPUs removed (uuids) and the ones added.
func (w *Workspace) AssignVGPUs(key, vm string, removed []string, add []ops.NewVGPU) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.AssignVGPUs(c, vm, removed, add) })
}

// VMUSB is a VM's USB page; haPriority is the Start Options page's choice.
func (w *Workspace) VMUSB(key, vm, haPriority string) (out app.VMUSBPage, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadVMUSB(c, vm, haPriority); return err })
	return out, err
}

// USBAttachChoices is the Attach USB dialog of a VM.
func (w *Workspace) USBAttachChoices(key, vm string) (out app.USBAttachChoices, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadUSBAttach(c, vm); return err })
	return out, err
}

// AttachUSB passes a USB device through to a VM.
func (w *Workspace) AttachUSB(key, vm, pusb string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.AttachUSB(c, vm, pusb) })
}

// DetachUSB takes a USB device from a VM (after the page's question).
func (w *Workspace) DetachUSB(key, vm, vusb string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.DetachUSB(c, vm, vusb) })
}

// CreateSriov makes an SR-IOV network (after the wizard's question).
func (w *Workspace) CreateSriov(key string, spec ops.NewSriov) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.CreateSriov(c, spec) })
}
