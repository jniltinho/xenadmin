package ops

import (
	"cmp"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/trust"
	"xcpng-center/internal/xapi"
)

// OVF/OVA appliance import: ImportApplianceAction (ImportApplianceAction.cs and Import.cs) with the import wizard's
// VmMappings: a server, the SR of each disk and the network of each interface per VM (rules.ApplianceMapping).
// shortcut: no OS fixup (see rules/ovfmap.go) and no encrypted packages (XenOvf's Security.cs ciphers); no vGPU, no
// PVS proxy, no InstallSection run, and the DVD drive stays empty.
// shortcut: no progress bar or Cancel yet (as the XVA import).

// applianceState reads what rules.CheckApplianceMappings checks against: the servers, the SRs that can take new disks
// (VDI_CREATE, not ISO) with the servers they are attached to, and the networks with the servers that see them.
func applianceState(s *xenapi.Session) (*rules.ApplianceState, error) {
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	srs, err := xenapi.SR.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	pbds, err := xenapi.PBD.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	sms, err := xenapi.SM.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	nets, err := xenapi.Network.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	pifs, err := xenapi.PIF.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	caps := map[string][]string{}
	for _, sm := range sms {
		caps[sm.Type] = append(caps[sm.Type], sm.Capabilities...)
	}
	names := map[string]string{}
	st := &rules.ApplianceState{SRs: map[string]rules.ApplianceSR{}, Networks: map[string][]string{}, HostName: func(u string) string { return names[u] }}
	for _, h := range hosts {
		st.Hosts = append(st.Hosts, h.UUID)
		names[h.UUID] = h.NameLabel
	}
	for _, sr := range srs {
		if sr.ContentType == "iso" || !slices.Contains(caps[sr.Type], "VDI_CREATE") {
			continue
		}
		out := rules.ApplianceSR{Name: sr.NameLabel, Free: int64(sr.PhysicalSize - sr.PhysicalUtilisation),
			Full: rules.SRFull(sr.Type, sr.PhysicalSize, sr.PhysicalUtilisation),
			Thin: slices.Contains(caps[sr.Type], "THIN_PROVISIONING"), LargeVDI: slices.Contains(caps[sr.Type], "LARGE_VDI")}
		for _, p := range sr.PBDs {
			if pbd, ok := pbds[p]; ok && pbd.CurrentlyAttached {
				out.Hosts = append(out.Hosts, hosts[pbd.Host].UUID)
			}
		}
		st.SRs[sr.UUID] = out
	}
nets:
	for _, n := range nets {
		if n.OtherConfig["is_guest_installer_network"] == "true" || n.OtherConfig["HideFromXenCenter"] == "true" { // Network.Show
			continue
		}
		var seers []string // nil: no PIF, every server sees it
		if len(n.PIFs) > 0 {
			seers = []string{}
		}
		for _, p := range n.PIFs {
			if pif := pifs[p]; !pif.Managed || pif.OtherConfig["HideFromXenCenter"] == "true" || // PIF.Show, IsMember, IsSriov
				pif.BondSlaveOf != "" && pif.BondSlaveOf != nullRef || len(pif.SriovLogicalPIFOf) > 0 {
				continue nets
			}
			seers = append(seers, hosts[pifs[p].Host].UUID)
		}
		st.Networks[n.UUID] = seers
	}
	return st, nil
}

// ApplianceMappingsError is why the mappings cannot be used for the appliance now ("" when they can).
func ApplianceMappingsError(c *xapi.Connection, a *rules.Appliance, maps []rules.ApplianceMapping) (string, error) {
	st, err := applianceState(c.Session)
	if err != nil {
		return "", err
	}
	if err := rules.CheckApplianceMappings(a, maps, st); err != nil {
		return err.Error(), nil
	}
	return "", nil
}

// ImportAppliance imports the appliance file with a mapping per VM (its server, the SR of each disk and the network of
// each interface), after verifying its signature (a signed package) and its manifest when asked, and starts it when
// asked. It returns the first VM's uuid.
func ImportAppliance(c *xapi.Connection, policy *trust.Policy, file string, maps []rules.ApplianceMapping, verify, start bool) (string, error) {
	s := c.Session
	if err := permitted(s, "vm_appliance.create", "vm_appliance.destroy", "vm.create", "vm.destroy", "vm.remove_from_other_config",
		"vdi.create", "vdi.destroy", "vbd.create", "vif.create"); err != nil { // ApplianceAction.StaticRBACDependencies
		return "", err
	}
	p, err := rules.OpenPackage(file)
	if err != nil {
		return "", err
	}
	a := p.Appliance
	if why, err := ApplianceMappingsError(c, a, maps); err != nil {
		return "", err
	} else if why != "" {
		return "", errors.New(why)
	}
	if verify && a.HasSignature { // ImportApplianceAction.RunCore: VERIFYING_SIGNATURE, then the manifest it signs
		if err := p.VerifySignature(nil); err != nil {
			return "", errors.New(messages.Format(messages.Get("VERIFYING_SIGNATURE_ERROR"), err.Error()))
		}
	}
	if verify { // ImportApplianceAction.RunCore: VERIFYING_MANIFEST
		if err := p.VerifyManifest(); err != nil {
			return "", errors.New(messages.Format(messages.Get("VERIFYING_MANIFEST_ERROR"), err.Error()))
		}
	}
	im := &applianceImport{c: c, policy: policy, pkg: p}
	first, err := im.run(a, maps)
	if err != nil {
		return "", errors.Join(err, im.cleanUp())
	}
	if start { // StartApplianceAction
		appUUID, _ := xenapi.VMAppliance.GetUUID(s, im.app)
		return first, VappPower(c, appUUID, "start")
	}
	return first, nil
}

// run is Process: the appliance (IMPORT_CREATING_APPLIANCE), then each VM; the first VM's uuid.
func (im *applianceImport) run(a *rules.Appliance, maps []rules.ApplianceMapping) (string, error) {
	s := im.c.Session
	app, err := xenapi.VMAppliance.Create(s, xenapi.VMApplianceRecord{NameLabel: a.Name, AllowedOperations: []xenapi.VMApplianceOperation{},
		CurrentOperations: map[string]xenapi.VMApplianceOperation{}, VMs: []xenapi.VMRef{}})
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	im.app = app
	var first string
	for i, spec := range a.VMs {
		uuid, err := im.vm(spec, maps[i])
		if err != nil {
			return "", err
		}
		first = cmp.Or(first, uuid)
	}
	return first, nil
}

// applianceImport is one import's state: what it made, removed when it fails.
type applianceImport struct {
	c      *xapi.Connection
	policy *trust.Policy
	pkg    *rules.Package
	app    xenapi.VMApplianceRef
	vms    []xenapi.VMRef
	vdi    xenapi.VDIRef // a disk not attached yet
}

// cleanUp is Process' CleanUpVm for every VM made (and its disks), the loose disk and the appliance; a VM that cannot
// be removed after a few tries (and its vApp) is reported, never silently left behind.
func (im *applianceImport) cleanUp() error {
	s := im.c.Session
	var left []string
	var disks []xenapi.VDIRef // the VMs' disks, checked after: destroyVM gives up on one quietly
	for _, vm := range im.vms {
		if vbds, err := xenapi.VM.GetVBDs(s, vm); err == nil {
			for _, b := range vbds {
				if v, err := xenapi.VBD.GetVDI(s, b); err == nil && v != nullRef {
					disks = append(disks, v)
				}
			}
		}
		gone := false
		for i := 0; i < 3 && !gone; i++ {
			if gone = destroyVM(s, vm); !gone {
				time.Sleep(time.Second)
			}
		}
		if !gone {
			name, _ := xenapi.VM.GetNameLabel(s, vm)
			left = append(left, "VM '"+name+"'")
		}
	}
	if im.vdi != "" {
		disks = append(disks, im.vdi)
		destroyVDI(s, im.vdi)
	}
	for _, v := range disks {
		if name, err := xenapi.VDI.GetNameLabel(s, v); err == nil { // still there
			left = append(left, "disk '"+name+"'")
		}
	}
	if im.app != "" && len(left) == 0 {
		if err := xenapi.VMAppliance.Destroy(s, im.app); err != nil {
			left = append(left, "the vApp")
		}
	}
	if len(left) > 0 {
		return errors.New("not removed after the failed import: " + strings.Join(left, ", "))
	}
	return nil
}

// destroyVDI destroys a disk the import made, waiting for the import task to let it go.
func destroyVDI(s *xenapi.Session, vdi xenapi.VDIRef) {
	for i := 0; i < 10; i++ {
		err := xenapi.VDI.Destroy(s, vdi)
		var api *xenapi.APIError
		if err == nil || errors.As(err, &api) && api.Message == "HANDLE_INVALID" {
			return
		}
		time.Sleep(time.Second)
	}
}

// vm is CreateVm, AddResourceSettingData for each item on the mapping's SRs and networks, and the VM made visible.
func (im *applianceImport) vm(spec rules.ApplianceVM, m rules.ApplianceMapping) (string, error) {
	s := im.c.Session
	rec := xenapi.VMRecord{
		NameLabel: cmp.Or(spec.Name, messages.Get("UNDEFINED_NAME_LABEL")), NameDescription: spec.Description, UserVersion: 1,
		MemoryTarget: int(spec.Memory), MemoryStaticMax: int(spec.Memory), MemoryDynamicMax: int(spec.Memory),
		MemoryDynamicMin: int(spec.Memory), MemoryStaticMin: int(spec.Memory),
		VCPUsMax: int(spec.VCPUsMax), VCPUsAtStartup: int(spec.VCPUs),
		ActionsAfterShutdown: xenapi.OnNormalExitDestroy, ActionsAfterReboot: xenapi.OnNormalExitRestart,
		ActionsAfterCrash: xenapi.OnCrashBehaviourRestart, ActionsAfterSoftreboot: xenapi.OnSoftrebootBehaviorSoftReboot,
		HVMShadowMultiplier: 1.0, PowerState: xenapi.VMPowerStateHalted, DomainType: xenapi.DomainTypeUnspecified,
		HVMBootPolicy: spec.HVMBootPolicy, HVMBootParams: spec.HVMBootParams, Platform: spec.Platform, NVRAM: spec.NVRAM,
		HasVendorDevice: spec.HasVendorDevice, Appliance: im.app, Order: spec.Order, StartDelay: spec.StartDelay, ShutdownDelay: spec.ShutdownDelay,
		OtherConfig: map[string]string{"HideFromXenCenter": "true"}, // hidden until its devices are made
		Affinity:    nullRef, SuspendVDI: nullRef, SuspendSR: nullRef, ResidentOn: nullRef, ScheduledToBeResidentOn: nullRef,
		Metrics: nullRef, GuestMetrics: nullRef, SnapshotOf: nullRef, Parent: nullRef, ProtectionPolicy: nullRef, SnapshotSchedule: nullRef,
	}
	for k, v := range spec.Fields { // VM.UpdateFrom: the other settings of an XCP-ng Center export
		switch k {
		case "PV_args":
			rec.PVArgs = v
		case "PV_bootloader":
			rec.PVBootloader = v
		case "PV_bootloader_args":
			rec.PVBootloaderArgs = v
		case "PV_kernel":
			rec.PVKernel = v
		case "PV_ramdisk":
			rec.PVRamdisk = v
		case "PV_legacy_args":
			rec.PVLegacyArgs = v
		case "HVM_shadow_multiplier":
			if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
				rec.HVMShadowMultiplier = f
			}
		case "hardware_platform_version":
			rec.HardwarePlatformVersion, _ = strconv.Atoi(v)
		case "recommendations":
			rec.Recommendations = v
		case "reference_label":
			rec.ReferenceLabel = v
		}
	}
	vm, err := xenapi.VM.Create(s, rec)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	im.vms = append(im.vms, vm)
	for i, d := range spec.Disks {
		if err := im.disk(vm, d, m.SRs[i]); err != nil {
			return "", err
		}
	}
	if spec.DVD != "" { // CD/DVD drive: empty, at its device when free
		if _, err := xenapi.VBD.Create(s, xenapi.VBDRecord{VM: vm, VDI: nullRef, Userdevice: freeDevice(s, vm, spec.DVD), Mode: xenapi.VbdModeRO,
			Type: xenapi.VbdTypeCD, Empty: true, OtherConfig: map[string]string{"owner": "true"}, QosAlgorithmParams: map[string]string{}}); err != nil {
			return "", xapi.Sanitize(err)
		}
	}
	for i, mac := range spec.MACs { // Network: the device index in order, the MAC when the OVF has one
		network, err := xenapi.Network.GetByUUID(s, m.Networks[i])
		if err != nil {
			return "", xapi.Sanitize(err)
		}
		if _, err := xenapi.VIF.Create(s, xenapi.VIFRecord{Device: strconv.Itoa(i), Network: network, VM: vm, MAC: mac, MTU: 1500,
			LockingMode: xenapi.VifLockingModeNetworkDefault, OtherConfig: map[string]string{}, QosAlgorithmParams: map[string]string{}}); err != nil {
			return "", xapi.Sanitize(err)
		}
	}
	if err := xenapi.VM.RemoveFromOtherConfig(s, vm, "HideFromXenCenter"); err != nil {
		return "", xapi.Sanitize(err)
	}
	uuid, err := xenapi.VM.GetUUID(s, vm)
	return uuid, xapi.Sanitize(err)
}

// freeDevice is VerifyUserDevice: the wanted userdevice when the VM allows it, else ("" or "99...") its first free one.
// The C# app leaves a disk whose device is taken unattached ("+"); here it goes to the first free one.
func freeDevice(s *xenapi.Session, vm xenapi.VMRef, want string) string {
	devices, err := xenapi.VM.GetAllowedVBDDevices(s, vm)
	if err != nil || len(devices) == 0 || want != "" && !strings.HasPrefix(want, "99") && slices.Contains(devices, want) {
		return want
	}
	return devices[0]
}

// disk is ImportFile and "connect VBD": a new VDI on its SR, the disk sent to import_raw_vdi, then attached to the VM.
func (im *applianceImport) disk(vm xenapi.VMRef, d rules.ApplianceDisk, srUUID string) error {
	s := im.c.Session
	sr, err := xenapi.SR.GetByUUID(s, srUUID)
	if err != nil {
		return xapi.Sanitize(err)
	}
	srRec, err := xenapi.SR.GetRecord(s, sr)
	if err != nil {
		return xapi.Sanitize(err)
	}
	stream, err := im.pkg.OpenDisk(d)
	if err != nil {
		return err
	}
	defer stream.Close()
	if free := int64(srRec.PhysicalSize - srRec.PhysicalUtilisation); free < stream.VirtualSize {
		return errors.New(messages.Format(messages.Get("SR_NOT_ENOUGH_SPACE"), srRec.UUID, rules.DiskSizeText(stream.VirtualSize), d.Href))
	}
	// shortcut: a VDI.create whose reply is lost leaves that disk behind (cleanUp only knows the disks it was told of);
	// tag the disks with the import and reconcile, as removeLostClone, if that shows up.
	vdi, err := xenapi.VDI.Create(s, xenapi.VDIRecord{NameLabel: d.Name, NameDescription: d.Description, SR: sr, VirtualSize: int(stream.VirtualSize),
		Type: xenapi.VdiTypeUser, OtherConfig: map[string]string{}, XenstoreData: map[string]string{}, SmConfig: map[string]string{}})
	if err != nil {
		return xapi.Sanitize(err)
	}
	im.vdi = vdi
	uuid, err := xenapi.VDI.GetUUID(s, vdi)
	if err != nil {
		return xapi.Sanitize(err)
	}
	task, err := xenapi.Task.Create(s, "import_raw_vdi_task", fmt.Sprintf("Importing disk %s to VDI %s", d.Href, uuid))
	if err != nil {
		return xapi.Sanitize(err)
	}
	address, port := storageHost(im.c, srRec)
	query := url.Values{"task_id": {string(task)}, "vdi": {uuid}, "format": {stream.Format}}
	err = put(s, im.policy, address, port, "/import_raw_vdi?"+query.Encode(), stream, stream.Length)
	if err = finishTask(s, task, err, nil); err != nil {
		return err
	}
	mode := xenapi.VbdModeRW
	if d.ReadOnly {
		mode = xenapi.VbdModeRO
	}
	if _, err := xenapi.VBD.Create(s, xenapi.VBDRecord{VM: vm, VDI: vdi, Userdevice: freeDevice(s, vm, d.Device), Bootable: d.Bootable, Mode: mode,
		Type: xenapi.VbdTypeDisk, OtherConfig: map[string]string{"owner": "true"}, QosAlgorithmParams: map[string]string{}}); err != nil {
		return xapi.Sanitize(err)
	}
	im.vdi = "" // the VM owns it now
	return nil
}
