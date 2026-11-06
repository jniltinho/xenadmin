package rules

import (
	"errors"
	"fmt"
	"slices"
	"sort"

	"xcpng-center/internal/messages"
)

// The import wizard's per-VM choices for an appliance (ImportWizard's VmMapping, XenModel/Mappings/VmMapping.cs): the
// server of each VM (ImportSelectHostPage), the SR of each disk (ImportSelectStoragePage) and the network of each
// interface (ImportSelectNetworkPage), as ImportApplianceAction applies them (OVF.SetTargetSRInRASD,
// SetTargetNetworkInRASD).
// shortcut: no OS fixup (ImportOptionsPage: the fixup ISO "xenserver-linuxfixup-disk.iso" ships with the Windows app
// only, and running it is Import.cs' ConfigureForXenServer, a boot of up to 600 s per VM); port both with an ISO to test.

// ApplianceMapping is the VmMapping of one VM of the package, in the package's VM order.
type ApplianceMapping struct {
	SysID string `json:"sys_id"`
	// Host is the server whose local SRs and networks the VM may use ("" for the pool or standalone server). As in the
	// wizard ("If you want to place the disks of a VM on local storage, specify the server below") it is no affinity:
	// ImportApplianceAction never sets one.
	Host     string   `json:"host"`
	SRs      []string `json:"srs"`      // the SR of each disk, in the VM's disk order
	Networks []string `json:"networks"` // the network of each interface, in device order
}

// DefaultApplianceMappings is OvfRules.Mappings: every VM on the pool, every disk on sr, every interface on network.
func DefaultApplianceMappings(a *Appliance, sr, network string) []ApplianceMapping {
	out := make([]ApplianceMapping, len(a.VMs))
	for i, vm := range a.VMs {
		out[i] = ApplianceMapping{SysID: vm.SysID, SRs: make([]string, len(vm.Disks)), Networks: make([]string, len(vm.MACs))}
		for j := range vm.Disks {
			out[i].SRs[j] = sr
		}
		for j := range vm.MACs {
			out[i].Networks[j] = network
		}
	}
	return out
}

// ApplianceSR is an SR that can take new disks (VDI_CREATE, not ISO), as the storage page's pickers see it.
type ApplianceSR struct {
	Name     string
	Hosts    []string // the servers it is attached to
	Free     int64
	Full     bool // SRFull
	Thin     bool // THIN_PROVISIONING: free space not checked (OvfRules.SpaceError)
	LargeVDI bool // LARGE_VDI: disks over DiskMaxSize allowed
}

// ApplianceState is the server's state the mappings are checked against (only the networks the pickers offer: not
// hidden, no bond member nor SR-IOV).
type ApplianceState struct {
	Hosts    []string                 // the pool's servers
	SRs      map[string]ApplianceSR   // by uuid
	Networks map[string][]string      // by uuid: the servers that see it (Host.CanSeeNetwork; nil when it has no PIF: all)
	HostName func(uuid string) string // for the messages
}

// DiskMaxSize is SR.DISK_MAX_SIZE (2 TiB), the largest disk an SR without LARGE_VDI takes.
const DiskMaxSize = 2 << 40

// CheckApplianceMappings refuses mappings that do not match the package (one per VM, its sys id, one SR per disk and
// one network per interface) or the server: an unknown server, SR or network; an SR the VM's server does not see, or
// that no server sees with the VM's other disks; a full SR or one too small for the disks put on it
// (SelectVMStorageWithMultipleVirtualDisksPage.CheckStorageRequirements, SrPickerItem's reasons); a network the server
// does not see, or not every server when none is chosen (SelectMultipleVMNetworkPage.ShowNetwork).
func CheckApplianceMappings(a *Appliance, maps []ApplianceMapping, st *ApplianceState) error {
	mismatch := errors.New("the choices do not match the appliance; read it again")
	if len(maps) != len(a.VMs) {
		return mismatch
	}
	need := map[string]int64{}
	for i, vm := range a.VMs {
		m := maps[i]
		if m.SysID != vm.SysID || len(m.SRs) != len(vm.Disks) || len(m.Networks) != len(vm.MACs) {
			return mismatch
		}
		seers := st.Hosts // the servers that see every disk of the VM
		if m.Host != "" {
			if !slices.Contains(st.Hosts, m.Host) {
				return fmt.Errorf("The server chosen for %s is not in the pool.", vm.Name)
			}
			seers = []string{m.Host}
		}
		for j, d := range vm.Disks {
			sr, ok := st.SRs[m.SRs[j]]
			label := vm.Name + " - " + d.Name
			switch {
			case !ok:
				return fmt.Errorf("%s: %s", label, messages.Get("IMPORT_SELECT_STORAGE_PAGE_ERROR_NO_SR_AVAILABLE"))
			case m.Host != "" && !slices.Contains(sr.Hosts, m.Host):
				return fmt.Errorf("%s: %s", label, messages.Format(messages.Get("SR_CANNOT_BE_SEEN"), st.HostName(m.Host)))
			case sr.Full:
				return fmt.Errorf("%s: %s", label, messages.Get("SRPICKER_SR_FULL"))
			case d.Capacity > DiskMaxSize && !sr.LargeVDI:
				return fmt.Errorf("%s: %s", label, messages.Format(messages.Get("SR_DISKSIZE_EXCEEDS_DISK_MAX_SIZE"), DiskSizeText(DiskMaxSize)))
			}
			seers = slices.DeleteFunc(slices.Clone(seers), func(h string) bool { return !slices.Contains(sr.Hosts, h) })
			if len(seers) == 0 {
				return fmt.Errorf("No server sees all the storage chosen for the disks of %s.", vm.Name)
			}
			if need[m.SRs[j]] += d.Capacity; need[m.SRs[j]] < 0 {
				need[m.SRs[j]] = 1<<63 - 1
			}
		}
		for j := range vm.MACs {
			hosts, ok := st.Networks[m.Networks[j]]
			sees := func(h string) bool { return hosts == nil || slices.Contains(hosts, h) }
			if !ok || m.Host != "" && !sees(m.Host) || m.Host == "" && slices.ContainsFunc(st.Hosts, func(h string) bool { return !sees(h) }) {
				return fmt.Errorf("%s - interface %d: %s", vm.Name, j, messages.Get("IMPORT_SELECT_NETWORK_PAGE_NO_AVAIL_NETWORKS"))
			}
		}
	}
	srs := make([]string, 0, len(need))
	for u := range need {
		srs = append(srs, u)
	}
	sort.Strings(srs)
	for _, u := range srs { // FreeSpace() <= the disks put on it
		if sr := st.SRs[u]; !sr.Thin && need[u] > 0 && sr.Free <= need[u] {
			return errors.New(messages.Format(messages.Get("IMPORT_SELECT_STORAGE_PAGE_ERROR_SPACE_LIMITATION"), sr.Name))
		}
	}
	return nil
}
