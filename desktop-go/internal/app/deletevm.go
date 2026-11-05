package app

import (
	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/xapi"
)

// VMDelete is ConfirmVMDeleteDialog for one VM: its title and question, the disks only it uses (checked when the VM
// owns them) and its snapshots (unchecked).
type VMDelete struct {
	Title     string       `json:"title"`
	Text      string       `json:"text"`
	Disks     []DeleteItem `json:"disks"`
	Snapshots []DeleteItem `json:"snapshots"`
}

// DeleteItem is a row of the dialog: a disk's VBD or a snapshot.
type DeleteItem struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Size    int    `json:"size,omitempty"`
	Checked bool   `json:"checked"`
}

// ReadVMDelete fills the Delete VM dialog.
func ReadVMDelete(c *xapi.Connection, vmUUID string) (VMDelete, error) {
	s := c.Session
	vm, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return VMDelete{}, xapi.Sanitize(err)
	}
	rec, err := xenapi.VM.GetRecord(s, vm)
	if err != nil {
		return VMDelete{}, xapi.Sanitize(err)
	}
	out := VMDelete{Title: messages.Format(messages.Get("CONFIRM_DELETE_TITLE"), "VM"),
		Text: messages.Format(messages.Get("CONFIRM_DELETE_VM"), rec.NameLabel, PoolName(c)), Disks: []DeleteItem{}, Snapshots: []DeleteItem{}}
	for _, b := range rec.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil || vbd.Type == xenapi.VbdTypeCD || vbd.VDI == "" || vbd.VDI == "OpaqueRef:NULL" {
			continue
		}
		vdi, err := xenapi.VDI.GetRecord(s, vbd.VDI)
		if err != nil || len(vdi.VBDs) > 1 { // a disk other VMs use too is never offered for one VM
			continue
		}
		_, owner := vbd.OtherConfig["owner"]
		out.Disks = append(out.Disks, DeleteItem{UUID: vbd.UUID, Name: vdi.NameLabel, Size: vdi.VirtualSize, Checked: owner})
	}
	for _, sn := range rec.Snapshots {
		if name, err := xenapi.VM.GetNameLabel(s, sn); err == nil {
			uuid, _ := xenapi.VM.GetUUID(s, sn)
			out.Snapshots = append(out.Snapshots, DeleteItem{UUID: uuid, Name: name})
		}
	}
	return out, nil
}
