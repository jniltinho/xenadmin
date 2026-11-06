package app

import (
	"slices"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// VMCopy is CopyVMDialog: the copy's name ("Copy of ...", unique) and description, which of fast clone and full copy
// the VM allows (VmRules.CopyModes), and the SRs a full copy can go to.
type VMCopy struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Clone       bool    `json:"clone"`
	Copy        bool    `json:"copy"`
	SRs         []SRUse `json:"srs"`
	DefaultSR   string  `json:"default_sr"`
}

// ReadVMCopy fills the Copy VM dialog.
func ReadVMCopy(c *xapi.Connection, vmUUID string) (VMCopy, error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return VMCopy{}, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return VMCopy{}, xapi.Sanitize(err)
	}
	names, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return VMCopy{}, xapi.Sanitize(err)
	}
	taken := map[string]bool{}
	for _, v := range names {
		taken[v.NameLabel] = true
	}
	// shortcut: the SRs are those a new disk of this VM can use (its home server), not VM.GetStorageHost's
	targets, err := NewDiskOptions(c, vmUUID)
	if err != nil {
		return VMCopy{}, err
	}
	out := VMCopy{Name: rules.UniqueName(messages.Format(messages.Get("ACTION_TEMPLATE_CLONE_NEW_NAME"), vm.NameLabel), taken),
		Description: vm.NameDescription, SRs: targets.SRs, DefaultSR: targets.DefaultSR}
	hasDisk, fastClonable := false, false
	sms, _ := xenapi.SM.GetAllRecords(s)
	caps := map[string][]string{}
	for _, sm := range sms {
		caps[sm.Type] = sm.Capabilities
	}
	var disks []xenapi.VDIRecord
	for _, b := range vm.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil || vbd.Type != xenapi.VbdTypeDisk {
			continue
		}
		hasDisk = true
		vdi, err := xenapi.VDI.GetRecord(s, vbd.VDI)
		if err != nil {
			continue
		}
		disks = append(disks, vdi)
		if t, err := xenapi.SR.GetType(s, vdi.SR); err == nil && slices.Contains(caps[t], "VDI_CLONE") {
			fastClonable = true
		}
	}
	// SrPickerCopyItem: an SR that cannot hold the VM's disks is shown with the reason, and cannot be picked
	for i, t := range out.SRs {
		ref, err := xenapi.SR.GetByUUID(s, t.UUID)
		if err != nil {
			continue
		}
		if sr, err := xenapi.SR.GetRecord(s, ref); err == nil {
			out.SRs[i].Reason = cannotFitDisks(sr, caps[sr.Type], disks)
		}
	}
	if out.DefaultSR != "" {
		for _, t := range out.SRs {
			if t.UUID == out.DefaultSR && t.Reason != "" {
				out.DefaultSR = ""
			}
		}
	}
	allowCopy := !vm.IsATemplate || slices.Contains(vm.AllowedOperations, xenapi.VMOperationsCopy)
	out.Clone, out.Copy = !allowCopy || fastClonable || !hasDisk, allowCopy && hasDisk
	return out, nil
}

// diskMaxSize is SR.DISK_MAX_SIZE (2 TiB), the largest disk an SR without LARGE_VDI takes.
const diskMaxSize = 2 << 40

// cannotFitDisks is SR.CanFitDisks: "" when the SR can hold copies of the disks, else why not (too big for the SR's
// driver, the SR full, larger than the SR, or not enough free space; a thin SR counts what the disks use).
func cannotFitDisks(sr xenapi.SRRecord, caps []string, disks []xenapi.VDIRecord) string {
	var virtual, physical int
	for _, d := range disks {
		if d.VirtualSize > diskMaxSize && !slices.Contains(caps, "LARGE_VDI") {
			return messages.Format(messages.Get("SR_DISKSIZE_EXCEEDS_DISK_MAX_SIZE"), diskSize(diskMaxSize))
		}
		virtual, physical = virtual+d.VirtualSize, physical+d.PhysicalUtilisation
	}
	free := sr.PhysicalSize - sr.PhysicalUtilisation
	if rules.SRFull(sr.Type, sr.PhysicalSize, sr.PhysicalUtilisation) {
		return messages.Get("SRPICKER_SR_FULL")
	}
	need := virtual
	if slices.Contains(caps, "THIN_PROVISIONING") {
		need = physical
	}
	if need > sr.PhysicalSize {
		return messages.Format(messages.Get("SR_PICKER_DISK_TOO_BIG"), diskSize(need), diskSize(sr.PhysicalSize))
	}
	if need > free {
		return messages.Format(messages.Get("SR_PICKER_INSUFFICIENT_SPACE"), diskSize(need), diskSize(free))
	}
	return ""
}
