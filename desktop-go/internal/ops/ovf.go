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

// OVF/OVA appliance import: OvfRules.ImportAction, i.e. ImportApplianceAction (ImportApplianceAction.cs and Import.cs)
// with OvfRules.Mappings: every disk on one SR and every interface on one network.
// shortcut: as the C# OvfRules, no per-VM/disk/interface choice, no OS fixup and no encrypted or signed-only checks
// (the certificate is not verified); no vGPU, no PVS proxy, no InstallSection run, and the DVD drive stays empty.
// shortcut: no progress bar or Cancel yet (as the XVA import).

// ApplianceSpaceError is OvfRules.SpaceError: "" when the SR has room for every disk (thin-provisioned SRs not checked).
func ApplianceSpaceError(c *xapi.Connection, srUUID string, required int64) (string, error) {
	s := c.Session
	ref, err := xenapi.SR.GetByUUID(s, srUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	sr, err := xenapi.SR.GetRecord(s, ref)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	sms, err := xenapi.SM.GetAllRecords(s)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	for _, sm := range sms {
		if sm.Type == sr.Type && slices.Contains(sm.Capabilities, "THIN_PROVISIONING") {
			return "", nil
		}
	}
	if free := int64(sr.PhysicalSize - sr.PhysicalUtilisation); free < required {
		return fmt.Sprintf("'%s' has %s free; the appliance's disks need %s.", sr.NameLabel, rules.DiskSizeText(free), rules.DiskSizeText(required)), nil
	}
	return "", nil
}

// ImportAppliance imports the appliance file to the SR, its interfaces on the network (networkUUID may be "" when it
// has none), after verifying its manifest when asked, and starts it when asked. It returns the first VM's uuid.
func ImportAppliance(c *xapi.Connection, policy *trust.Policy, file, srUUID, networkUUID string, verify, start bool) (string, error) {
	p, err := rules.OpenPackage(file)
	if err != nil {
		return "", err
	}
	a := p.Appliance
	if verify { // ImportApplianceAction.RunCore: VERIFYING_MANIFEST
		if err := p.VerifyManifest(); err != nil {
			return "", errors.New(messages.Format(messages.Get("VERIFYING_MANIFEST_ERROR"), err.Error()))
		}
	}
	s := c.Session
	sr, err := xenapi.SR.GetByUUID(s, srUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	srRec, err := xenapi.SR.GetRecord(s, sr)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	if rules.SRFull(srRec.Type, srRec.PhysicalSize, srRec.PhysicalUtilisation) || srRec.ContentType == "iso" { // XvaRules.TargetError
		return "", errors.New("the storage repository is no longer available or is full")
	}
	network := xenapi.NetworkRef(nullRef)
	if a.Interfaces > 0 {
		if networkUUID == "" {
			return "", errors.New("Select a network.")
		}
		if network, err = xenapi.Network.GetByUUID(s, networkUUID); err != nil {
			return "", xapi.Sanitize(err)
		}
	}
	im := &applianceImport{c: c, policy: policy, pkg: p, sr: sr, srRec: srRec, network: network}
	first, err := im.run(a)
	if err != nil {
		im.cleanUp()
		return "", err
	}
	if start { // StartApplianceAction
		appUUID, _ := xenapi.VMAppliance.GetUUID(s, im.app)
		return first, VappPower(c, appUUID, "start")
	}
	return first, nil
}

// run is Process: the appliance (IMPORT_CREATING_APPLIANCE), then each VM; the first VM's uuid.
func (im *applianceImport) run(a *rules.Appliance) (string, error) {
	s := im.c.Session
	app, err := xenapi.VMAppliance.Create(s, xenapi.VMApplianceRecord{NameLabel: a.Name, AllowedOperations: []xenapi.VMApplianceOperation{},
		CurrentOperations: map[string]xenapi.VMApplianceOperation{}, VMs: []xenapi.VMRef{}})
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	im.app = app
	var first string
	for _, spec := range a.VMs {
		uuid, err := im.vm(spec)
		if err != nil {
			return "", err
		}
		first = cmp.Or(first, uuid)
	}
	return first, nil
}

// applianceImport is one import's state: what it made, removed when it fails.
type applianceImport struct {
	c       *xapi.Connection
	policy  *trust.Policy
	pkg     *rules.Package
	sr      xenapi.SRRef
	srRec   xenapi.SRRecord
	network xenapi.NetworkRef
	app     xenapi.VMApplianceRef
	vms     []xenapi.VMRef
	vdi     xenapi.VDIRef // a disk not attached yet
}

// cleanUp is Process' CleanUpVm for every VM made (and its disks), the loose disk and the appliance.
func (im *applianceImport) cleanUp() {
	s := im.c.Session
	for _, vm := range im.vms {
		destroyVM(s, vm)
	}
	if im.vdi != "" {
		destroyVDI(s, im.vdi)
	}
	if im.app != "" {
		_ = xenapi.VMAppliance.Destroy(s, im.app)
	}
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

// vm is CreateVm, AddResourceSettingData for each item, and the VM made visible.
func (im *applianceImport) vm(spec rules.ApplianceVM) (string, error) {
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
	for _, d := range spec.Disks {
		if err := im.disk(vm, d); err != nil {
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
		if _, err := xenapi.VIF.Create(s, xenapi.VIFRecord{Device: strconv.Itoa(i), Network: im.network, VM: vm, MAC: mac, MTU: 1500,
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

// disk is ImportFile and "connect VBD": a new VDI on the SR, the disk sent to import_raw_vdi, then attached to the VM.
func (im *applianceImport) disk(vm xenapi.VMRef, d rules.ApplianceDisk) error {
	s := im.c.Session
	stream, err := im.pkg.OpenDisk(d)
	if err != nil {
		return err
	}
	defer stream.Close()
	if free := int64(im.srRec.PhysicalSize - im.srRec.PhysicalUtilisation); free < stream.VirtualSize {
		return errors.New(messages.Format(messages.Get("SR_NOT_ENOUGH_SPACE"), im.srRec.UUID, rules.DiskSizeText(stream.VirtualSize), d.Href))
	}
	// shortcut: a VDI.create whose reply is lost leaves that disk behind (cleanUp only knows the disks it was told of);
	// tag the disks with the import and reconcile, as removeLostClone, if that shows up.
	vdi, err := xenapi.VDI.Create(s, xenapi.VDIRecord{NameLabel: d.Name, NameDescription: d.Description, SR: im.sr, VirtualSize: int(stream.VirtualSize),
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
	address, port := storageHost(im.c, im.srRec)
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
