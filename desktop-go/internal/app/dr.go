package app

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// What Pool > Disaster Recovery reads (XenAdmin.Desktop.Core/DrRules.cs, XenAdmin.Desktop/DrWindows.cs): the SRs that
// can hold the pool's metadata, the pools, vApps and VMs found in the metadata, and the wizard's prechecks.
// shortcut: the wizard reads metadata from SRs already attached to this pool, as the C# Desktop app; it does not probe
// iSCSI/HBA devices for mirrored SRs (DRFailoverWizardStoragePage "Find Storage Repositories"): attach them with New
// Storage first. Upgrade when the C# app gets that page.

// DRSR is an SR that can hold pool metadata: DR on (this pool's metadata on it) and any pool's metadata on it.
type DRSR struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Free        string `json:"free"`
	Enabled     bool   `json:"enabled"`
	HasMetadata bool   `json:"has_metadata"`
}

// DRSRs is DRConfigureDialog's and DRFailoverWizardStoragePage's list, with the dialog's title and its action's.
type DRSRs struct {
	Pool   string `json:"pool"`
	Title  string `json:"title"`
	Action string `json:"action"`
	MaxSRs int    `json:"max_srs"`
	SRs    []DRSR `json:"srs"`
}

// ReadDRSRs is DrRules.ReplicationSrs with IsDrEnabled and StorageChoices' metadata (one server call per SR).
func ReadDRSRs(c *xapi.Connection) (DRSRs, error) {
	s := c.Session
	pool := ops.PoolName(c)
	srs, err := ops.ReplicationSRs(s)
	if err != nil {
		return DRSRs{}, err
	}
	out := DRSRs{Pool: pool, Title: messages.Format(messages.Get("DR_CONFIGURE_TITLE"), pool), Action: messages.Get("ACTION_DR_CONFIGURING"), MaxSRs: rules.DRMaxSRs, SRs: []DRSR{}}
	for _, r := range srs {
		out.SRs = append(out.SRs, DRSR{UUID: r.Record.UUID, Name: r.Record.NameLabel, Type: srTypeName(r.Record.Type),
			Free: rules.DiskSizeText(int64(r.Record.PhysicalSize - r.Record.PhysicalUtilisation)), Enabled: r.Enabled, HasMetadata: r.HasMetadata})
	}
	return out, nil
}

// DRObject is a vApp (with its VM count) or a loose VM of a pool's metadata.
type DRObject struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	VMs  int    `json:"vms,omitempty"`
}

// DRPool is one pool's metadata found on the chosen SRs: its vApps and the VMs outside them (LooseVms); Current is
// IsCurrentPool (its own metadata cannot be recovered onto it), with Note saying so.
type DRPool struct {
	UUID    string     `json:"uuid"`
	Name    string     `json:"name"`
	VDI     string     `json:"vdi"`
	Current bool       `json:"current"`
	Note    string     `json:"note,omitempty"`
	VApps   []DRObject `json:"vapps"`
	VMs     []DRObject `json:"vms"`
}

// LoadDRMetadata is DrRules.LoadMetadata: the pools whose metadata is on these SRs (uuids), once each, by name.
func LoadDRMetadata(c *xapi.Connection, srUUIDs []string) ([]DRPool, error) {
	s := c.Session
	_, here, _, err := ops.PoolOf(s)
	if err != nil {
		return nil, err
	}
	vdis, err := ops.MetadataVDIs(s, srUUIDs)
	if err != nil {
		return nil, err
	}
	out := []DRPool{}
	for _, vdi := range vdis {
		_, err := ops.WithMetadata(s, vdi, func(_ *xenapi.Session, md ops.DRMetadata) error {
			for _, p := range out {
				if p.UUID == md.Pool.UUID {
					return nil
				}
			}
			out = append(out, drPool(md, here.UUID))
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	return out, nil
}

// drPool is a pool's metadata as DRFailoverWizardAppliancesPage lists it: its vApps and loose VMs, by name, and
// whether it is this pool (here: its uuid).
func drPool(md ops.DRMetadata, here string) DRPool {
	p := DRPool{UUID: md.Pool.UUID, Name: md.Pool.NameLabel, VDI: md.VDIUUID, Current: md.Pool.UUID == here, VApps: []DRObject{}, VMs: []DRObject{}}
	if p.Current {
		p.Note = messages.Get("DR_WIZARD_APPLIANCESPAGE_CURRENT_POOL")
	}
	for _, a := range md.VApps {
		p.VApps = append(p.VApps, DRObject{UUID: a.UUID, Name: a.NameLabel, VMs: len(a.VMs)})
	}
	for _, v := range md.VMs {
		if ops.Loose(v) {
			p.VMs = append(p.VMs, DRObject{UUID: v.UUID, Name: v.NameLabel})
		}
	}
	sortObjects(p.VApps)
	sortObjects(p.VMs)
	return p
}

func sortObjects(list []DRObject) {
	sort.Slice(list, func(i, j int) bool {
		a, b := strings.ToLower(list[i].Name), strings.ToLower(list[j].Name)
		return a < b || a == b && list[i].UUID < list[j].UUID
	})
}

// DRSource is a pool connected here (RunningVmCheck's source pool): its VMs that are not halted, uuid to name.
type DRSource map[string]string

// ReadDRSource is the pool's uuid and its VMs that are not halted (for the prechecks of another pool's wizard).
func ReadDRSource(c *xapi.Connection) (string, DRSource, error) {
	s := c.Session
	_, p, _, err := ops.PoolOf(s)
	if err != nil {
		return "", nil, err
	}
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return "", nil, xapi.Sanitize(err)
	}
	running := DRSource{}
	for _, v := range vms {
		if rules.DRRealVM(v.IsATemplate, v.IsASnapshot, v.IsControlDomain) && v.PowerState != xenapi.VMPowerStateHalted {
			running[v.UUID] = v.NameLabel
		}
	}
	return p.UUID, running, nil
}

// DRPrechecks is DrRules.Prechecks (DRFailoverWizardPrecheckPage.GenerateChecks): for a test failover HA must be off
// here; for a failover or failback a VM still running in its source pool (sources: the pools connected here, by
// uuid) blocks the recovery; a VM (or a vApp's VM) that already exists here must be older and halted; every chosen
// object must pass VM / VM_appliance.assert_can_be_recovered.
func DRPrechecks(c *xapi.Connection, typ string, choices []ops.DRChoice, sources map[string]DRSource) ([]rules.DRCheck, error) {
	if !rules.DRKnownType(typ) {
		return nil, errors.New("unknown recovery type")
	}
	s := c.Session
	_, pool, _, err := ops.PoolOf(s)
	if err != nil {
		return nil, err
	}
	out := []rules.DRCheck{}
	if text := rules.DRHACheck(typ, pool.HaEnabled, pool.NameLabel); text != "" {
		out = append(out, rules.DRCheck{Text: text, Error: true})
	}
	local, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	localApps, err := xenapi.VMAppliance.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	for _, ch := range choices {
		vdi, err := ops.MetadataVDI(s, ch.VDI)
		if err != nil {
			return nil, err
		}
		ok, err := ops.WithMetadata(s, vdi, func(meta *xenapi.Session, md ops.DRMetadata) error {
			apps, vms, err := md.Resolve(ch)
			if err != nil {
				return err
			}
			all := append([]xenapi.VMRef{}, vms...)
			for _, a := range apps {
				for _, v := range md.VApps[a].VMs {
					if _, ok := md.VMs[v]; ok {
						all = append(all, v)
					}
				}
			}
			if src, connected := sources[md.Pool.UUID]; typ != rules.DRTest && connected {
				for _, v := range all {
					if name, running := src[md.VMs[v].UUID]; running {
						out = append(out, rules.DRRunningCheck(name, md.Pool.NameLabel))
					}
				}
			}
			for _, v := range all {
				for _, e := range local {
					if e.UUID == md.VMs[v].UUID {
						out = append(out, rules.DRExistingCheck(e.NameLabel, e.Version, md.VMs[v].Version, e.PowerState == xenapi.VMPowerStateHalted, typ == rules.DRTest))
					}
				}
			}
			if typ == rules.DRTest { // nor the vApps and snapshots it would bring in (ops.RecoverDR refuses them too)
				objs := md.Objects(apps, vms)
				for _, e := range local {
					if e.IsASnapshot && slices.Contains(objs.Snapshots, e.UUID) {
						out = append(out, rules.DRTestCollision("snapshot", e.NameLabel))
					}
				}
				for _, e := range localApps {
					if slices.Contains(objs.VApps, e.UUID) {
						out = append(out, rules.DRTestCollision("vApp", e.NameLabel))
					}
				}
			}
			for _, a := range apps {
				if err := xenapi.VMAppliance.AssertCanBeRecovered(meta, a, s.Ref()); err != nil {
					out = append(out, rules.DRCheck{Text: fmt.Sprintf("vApp '%s': %v", md.VApps[a].NameLabel, xapi.Sanitize(err)), Error: true})
				}
			}
			for _, v := range vms {
				if err := xenapi.VM.AssertCanBeRecovered(meta, v, s.Ref()); err != nil {
					out = append(out, rules.DRCheck{Text: fmt.Sprintf("VM '%s': %v", md.VMs[v].NameLabel, xapi.Sanitize(err)), Error: true})
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		if !ok {
			out = append(out, rules.DRCheck{Text: fmt.Sprintf("The metadata on VDI %s could not be opened.", ch.VDI), Error: true})
		}
	}
	return out, nil
}
