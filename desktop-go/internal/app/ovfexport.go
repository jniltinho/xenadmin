package app

import (
	"slices"
	"sort"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// ExportPick is a VM of the Export wizard's list (ExportSelectVMsPage.GetDataGridViewRow).
type ExportPick struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        string `json:"size"` // its disks' virtual size (GetTotalVmSize)
	Vapp        string `json:"vapp"` // "" for none
	Selected    bool   `json:"selected"`
	VTPM        bool   `json:"vtpm"` // VTPM_EXPORT_UNSUPPORTED_FOR_OVF
}

// ExportChoices are the VMs the wizard offers (ExportSelectVMsPage.PopulatePage): the exportable ones of the pool, the
// selected VM or the vApp's VMs ticked, ticked ones first.
type ExportChoices struct {
	Pool   string       `json:"pool"`
	IsPool bool         `json:"is_pool"`
	VMs    []ExportPick `json:"vms"`
}

// ReadExportChoices reads them for the selected VM uuid or vApp uuid (either may be "").
func ReadExportChoices(c *xapi.Connection, vmUUID, vappUUID string) (ExportChoices, error) {
	s := c.Session
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return ExportChoices{}, xapi.Sanitize(err)
	}
	apps, err := xenapi.VMAppliance.GetAllRecords(s)
	if err != nil {
		return ExportChoices{}, xapi.Sanitize(err)
	}
	hosts, _ := xenapi.Host.GetAll(s)
	out := ExportChoices{Pool: PoolName(c), IsPool: len(hosts) > 1, VMs: []ExportPick{}}
	if pools, err := xenapi.Pool.GetAllRecords(s); err == nil {
		for _, p := range pools {
			out.IsPool = out.IsPool || p.NameLabel != ""
		}
	}
	for _, v := range vms {
		// IsVmExportable: a real, shown VM, halted or suspended, not locked, that allows export
		if v.IsATemplate || v.IsASnapshot || v.IsControlDomain || !vmShown(v) || len(v.CurrentOperations) > 0 ||
			v.PowerState != xenapi.VMPowerStateHalted && v.PowerState != xenapi.VMPowerStateSuspended ||
			!slices.Contains(v.AllowedOperations, xenapi.VMOperationsExport) {
			continue
		}
		var size int64
		for _, b := range v.VBDs {
			if vbd, err := xenapi.VBD.GetRecord(s, b); err == nil && vbd.Type != xenapi.VbdTypeCD && vbd.VDI != nullRef {
				if n, err := xenapi.VDI.GetVirtualSize(s, vbd.VDI); err == nil {
					size += int64(n)
				}
			}
		}
		p := ExportPick{UUID: v.UUID, Name: v.NameLabel, Description: v.NameDescription, Size: rules.DiskSizeText(size), VTPM: len(v.VTPMs) > 0,
			Selected: v.UUID == vmUUID}
		if a, ok := apps[v.Appliance]; ok {
			p.Vapp = a.NameLabel
			p.Selected = p.Selected || vappUUID != "" && a.UUID == vappUUID
		}
		out.VMs = append(out.VMs, p)
	}
	sort.Slice(out.VMs, func(i, j int) bool {
		if a, b := out.VMs[i], out.VMs[j]; a.Selected != b.Selected {
			return a.Selected
		}
		return strings.ToLower(out.VMs[i].Name) < strings.ToLower(out.VMs[j].Name)
	})
	return out, nil
}
