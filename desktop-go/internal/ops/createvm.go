package ops

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/xapi"
)

// NewVM is what the New VM wizard asks for (CreateVMAction's arguments).
type NewVM struct {
	Template    string `json:"template"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// Install is the installation method (CreateVMAction's InstallMethod): "cd" (the DVD drive), "network" (boot order
	// "n" first, or a PV installer's URL) or "none" (a custom template without a DVD drive)
	Install    string    `json:"install"`
	URL        string    `json:"url"`       // a PV (eliloader) template's network install repository
	ISO        string    `json:"iso"`       // ISO VDI uuid, "" for an empty drive
	BootMode   string    `json:"boot_mode"` // bios, uefi, secureuefi (HVM)
	HomeServer string    `json:"home_server"`
	VCPUs      int       `json:"vcpus"`
	Cores      int64     `json:"cores"`
	Memory     int       `json:"memory"`
	Disks      []NewDisk `json:"disks"`
	Networks   []NewVIF  `json:"networks"` // one interface each
	Start      bool      `json:"start"`
}

// NewVIF is an interface of the new VM: its network, and the device and bandwidth limit a custom template's VIF had.
type NewVIF struct {
	Network   string            `json:"network"` // network uuid
	Device    string            `json:"device,omitempty"`
	QosType   string            `json:"qos_type,omitempty"`
	QosParams map[string]string `json:"qos_params,omitempty"`
}

// Created is CreateVM's outcome: the VM exists once UUID is set, even when starting it failed (StartError).
type Created struct {
	UUID       string `json:"uuid"`
	StartError string `json:"start_error,omitempty"`
}

// NewDisk is a disk of the new VM: a new one on an SR at a position (a custom template's disk at that position is
// kept, or moved to that SR).
type NewDisk struct {
	Position    string `json:"position"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int    `json:"size"`
	SR          string `json:"sr"`
}

// CreateVM is CreateVMAction.Run: clone the template under a hidden name, apply the vendor-device recommendation,
// provision it (without the template's disk list: the wizard's disks are created here), set name, vCPUs, cores,
// home server and memory, the boot parameters, the DVD drive with the ISO, the disks (the first one's SR becomes the
// suspend SR) and the network interfaces, then start it when asked. A failure before the VM is complete destroys it.
// shortcut: no full copy, BIOS strings, vTPM, vGPU or cloud-config drive; disks of a custom template are kept on their
// SR only (not moved).
func CreateVM(c *xapi.Connection, spec NewVM) (Created, error) {
	if strings.TrimSpace(spec.Name) == "" {
		return Created{}, errors.New("enter a name for the VM")
	}
	s := c.Session
	tref, err := xenapi.VM.GetByUUID(s, spec.Template)
	if err != nil {
		return Created{}, xapi.Sanitize(err)
	}
	t, err := xenapi.VM.GetRecord(s, tref)
	if err != nil {
		return Created{}, xapi.Sanitize(err)
	}
	if !t.IsATemplate {
		return Created{}, errors.New("not a template")
	}
	// Helpers.MakeHiddenName, unique to this attempt: a lost clone can be found by it and removed, nothing else
	nonce := make([]byte, 8)
	if _, err := rand.Read(nonce); err != nil {
		return Created{}, err
	}
	hidden := "__gui__" + spec.Name + "-" + hex.EncodeToString(nonce)
	result, err := waitResult(s)(xenapi.VM.AsyncClone(s, tref, hidden))
	vm := xenapi.VMRef(resultRef(result))
	if vm == "" {
		// only a failure XAPI reported for the task proves there is no clone; anything else (a lost answer, an
		// unfollowed task) may still produce one later
		var api *xenapi.APIError
		removeLostClone(s, hidden, !errors.As(err, &api))
		if err == nil {
			err = errors.New("the template was cloned but XAPI did not say which VM it is")
		}
		return Created{}, xapi.Sanitize(err)
	}
	complete := false
	defer func() {
		if !complete { // CreateVMAction.CleanOnError: the half-made VM goes, with the disks it owns (retried)
			for i := 0; i < 30 && !destroyVM(s, vm); i++ {
				time.Sleep(2 * time.Second)
			}
		}
	}()
	if err := create(s, vm, t, spec); err != nil {
		return Created{}, xapi.Sanitize(err)
	}
	complete = true
	out := Created{}
	out.UUID, _ = xenapi.VM.GetUUID(s, vm)
	if spec.Start {
		if err := wait(s)(xenapi.VM.AsyncStart(s, vm, false, false)); err != nil {
			out.StartError = xapi.Sanitize(err).Error() // the VM stays: only starting it failed
		}
	}
	return out, nil
}

func create(s *xenapi.Session, vm xenapi.VMRef, t xenapi.VMRecord, spec NewVM) error {
	// ApplyRecommendationsForVendorDevice
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return err
	}
	noVendor := true
	for _, p := range pools {
		noVendor = p.PolicyNoVendorDevice
	}
	if hasRecommendation(t.Recommendations, "has-vendor-device") && !noVendor && !vendorDeviceRestricted(s) {
		if err := xenapi.VM.SetHasVendorDevice(s, vm, true); err != nil {
			return err
		}
	}
	// ProvisionVM: the wizard creates every disk itself
	if _, ok := t.OtherConfig["disks"]; ok {
		if err := xenapi.VM.RemoveFromOtherConfig(s, vm, "disks"); err != nil {
			return err
		}
	}
	if err := wait(s)(xenapi.VM.AsyncProvision(s, vm)); err != nil {
		return err
	}
	rec, err := xenapi.VM.GetRecord(s, vm)
	if err != nil {
		return err
	}
	// SetVMParams
	if err := xenapi.VM.SetNameLabel(s, vm, strings.TrimSpace(spec.Name)); err != nil {
		return err
	}
	if err := xenapi.VM.SetNameDescription(s, vm, spec.Description); err != nil {
		return err
	}
	if spec.VCPUs < 1 {
		spec.VCPUs = 1
	}
	if err := changeVCPUs(s, vm, rec, VMChanges{VCPUsMax: &spec.VCPUs, VCPUsAtStartup: &spec.VCPUs}); err != nil {
		return err
	}
	platform := map[string]string{}
	for k, v := range rec.Platform {
		platform[k] = v
	}
	if spec.Cores > 0 {
		platform["cores-per-socket"] = itoa64(spec.Cores)
	}
	home := xenapi.HostRef(nullRef)
	if spec.HomeServer != "" {
		if home, err = xenapi.Host.GetByUUID(s, spec.HomeServer); err != nil {
			return err
		}
	}
	if home != t.Affinity {
		if err := xenapi.VM.SetAffinity(s, vm, home); err != nil {
			return err
		}
	}
	if spec.Memory > 0 && (t.MemoryDynamicMin != spec.Memory || t.MemoryDynamicMax != spec.Memory || t.MemoryStaticMax != spec.Memory) {
		if err := xenapi.VM.SetMemoryLimits(s, vm, t.MemoryStaticMin, spec.Memory, spec.Memory, spec.Memory); err != nil {
			return err
		}
	}
	// SetVMBootParams
	hvm := t.HVMBootPolicy != ""
	if hvm {
		params := map[string]string{}
		for k, v := range rec.HVMBootParams {
			params[k] = v
		}
		if len(spec.Disks) == 0 || spec.Install == "network" { // CA-46213: network boot first
			params["order"] = networkFirst(params["order"])
		}
		params["firmware"] = map[bool]string{true: "bios", false: "uefi"}[spec.BootMode == "" || spec.BootMode == "bios"]
		if err := xenapi.VM.SetHVMBootParams(s, vm, params); err != nil {
			return err
		}
		platform["secureboot"] = map[bool]string{true: "true", false: "false"}[spec.BootMode == "secureuefi"]
	} else if err := xenapi.VM.SetPVArgs(s, vm, t.PVArgs); err != nil {
		return err
	}
	if !hvm && t.PVBootloader == "eliloader" && spec.Install != "none" { // the installer's repository: the URL, or the DVD drive
		repo := "cdrom"
		if spec.Install == "network" {
			repo = spec.URL
			if d := rec.OtherConfig["install-distro"]; d == "rhel41" || d == "rhel44" || d == "rhlike" {
				repo = rhelRepo(repo)
			}
		}
		_ = xenapi.VM.RemoveFromOtherConfig(s, vm, "install-repository")
		if err := xenapi.VM.AddToOtherConfig(s, vm, "install-repository", repo); err != nil {
			return err
		}
	}
	if err := xenapi.VM.SetPlatform(s, vm, platform); err != nil {
		return err
	}
	if !customWithNoDVD(s, t) { // a custom template without a DVD drive gets none (Helpers.CustomWithNoDVD)
		if err := addCD(s, vm, rec, spec); err != nil {
			return err
		}
	}
	if err := addDisks(s, vm, spec); err != nil {
		return err
	}
	return addNetworks(s, vm, rec, spec)
}

// addCD is AddCdDrive: the template's DVD drive at position 0-3, else a new one; emptied, then the ISO inserted.
func addCD(s *xenapi.Session, vm xenapi.VMRef, rec xenapi.VMRecord, spec NewVM) error {
	var cd xenapi.VBDRef
	empty := true
	for _, b := range rec.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err == nil && vbd.Type == xenapi.VbdTypeCD && strings.Contains("0123", vbd.Userdevice) && len(vbd.Userdevice) == 1 {
			cd, empty = b, vbd.Empty
			break
		}
	}
	if cd == "" {
		devices, err := xenapi.VM.GetAllowedVBDDevices(s, vm)
		if err != nil {
			return err
		}
		if len(devices) == 0 {
			return errors.New(messages.Get("NO_MORE_USERDEVICES"))
		}
		position := devices[0]
		if contains(devices, "3") {
			position = "3"
		}
		if cd, err = xenapi.VBD.Create(s, xenapi.VBDRecord{VM: vm, VDI: nullRef, Userdevice: position, Bootable: (spec.Install == "cd"), Mode: xenapi.VbdModeRO,
			Type: xenapi.VbdTypeCD, Empty: true, Unpluggable: true, OtherConfig: map[string]string{}, QosAlgorithmParams: map[string]string{}}); err != nil {
			return err
		}
	}
	if !empty {
		if err := wait(s)(xenapi.VBD.AsyncEject(s, cd)); err != nil {
			return err
		}
	}
	if (spec.Install == "cd") && spec.ISO != "" {
		iso, err := xenapi.VDI.GetByUUID(s, spec.ISO)
		if err != nil {
			return err
		}
		return wait(s)(xenapi.VBD.AsyncInsert(s, cd, iso))
	}
	return nil
}

// addDisks is AddDisks: for each disk, the VM's disk at that position (cloned from a custom template) is kept when it
// is on the chosen SR, else copied there (MoveDisk); a position without one gets a new disk (CreateDisk). New and moved
// disks get the VM as placement hint and a VBD owned by the VM, at that position when free (else the first free one),
// bootable at position 0 unless installing from the DVD; the first disk's SR is the suspend SR.
func addDisks(s *xenapi.Session, vm xenapi.VMRef, spec NewVM) error {
	suspend := xenapi.SRRef(nullRef)
	existing := map[string]xenapi.VBDRecord{}
	refs := map[string]xenapi.VBDRef{}
	if vbds, err := xenapi.VM.GetVBDs(s, vm); err == nil {
		for _, b := range vbds {
			if r, err := xenapi.VBD.GetRecord(s, b); err == nil && r.Type == xenapi.VbdTypeDisk {
				existing[r.Userdevice], refs[r.Userdevice] = r, b
			}
		}
	}
	for i, d := range spec.Disks {
		sr, err := xenapi.SR.GetByUUID(s, d.SR)
		if err != nil {
			return err
		}
		var vdi xenapi.VDIRef
		if old, ok := existing[d.Position]; ok && old.VDI != nullRef {
			oldSR, err := xenapi.VDI.GetSR(s, old.VDI)
			if err != nil {
				return err
			}
			if oldSR == sr { // DiskOk: kept as cloned
				vdi = old.VDI
			} else { // MoveDisk
				res, err := waitResult(s)(xenapi.VDI.AsyncCopy(s, old.VDI, sr, nullRef, nullRef))
				if err != nil {
					return err
				}
				vdi = xenapi.VDIRef(resultRef(res))
				if vdi == "" {
					return errors.New("the disk was copied but XAPI did not say which disk it is")
				}
				// until the copy is attached, a failure deletes it (destroyVM only finds attached disks)
				attach := func() error {
					if err := addVMHint(s, vdi, vm); err != nil {
						return err
					}
					if err := wait(s)(xenapi.VBD.AsyncDestroy(s, refs[d.Position])); err != nil {
						return err
					}
					if err := wait(s)(xenapi.VDI.AsyncDestroy(s, old.VDI)); err != nil {
						return err
					}
					return createVBD(s, vm, vdi, d.Position, d.Position == "0")
				}
				if err := attach(); err != nil {
					_ = xenapi.VDI.Destroy(s, vdi)
					return err
				}
			}
		} else {
			vdi, err = xenapi.VDI.Create(s, xenapi.VDIRecord{NameLabel: d.Name, NameDescription: d.Description, SR: sr, VirtualSize: d.Size,
				Type: xenapi.VdiTypeUser, OtherConfig: map[string]string{}, XenstoreData: map[string]string{}, SmConfig: map[string]string{"vmhint": string(vm)}})
			if err != nil {
				return err
			}
			if err := createVBD(s, vm, vdi, d.Position, d.Position == "0" && spec.Install != "cd"); err != nil {
				_ = xenapi.VDI.Destroy(s, vdi)
				return err
			}
		}
		if err := xenapi.VDI.SetNameLabel(s, vdi, d.Name); err != nil {
			return err
		}
		if err := xenapi.VDI.SetNameDescription(s, vdi, d.Description); err != nil {
			return err
		}
		if i == 0 {
			if t, err := xenapi.SR.GetType(s, sr); err == nil && t != "rawhba" { // SR.HBALunPerVDI has no suspend image
				suspend = sr
			}
		}
	}
	return xenapi.VM.SetSuspendSR(s, vm, suspend)
}

// createVBD is CreateVbd: owned by the VM, at the position when free, else the first free one.
func createVBD(s *xenapi.Session, vm xenapi.VMRef, vdi xenapi.VDIRef, position string, bootable bool) error {
	devices, err := xenapi.VM.GetAllowedVBDDevices(s, vm)
	if err != nil || len(devices) == 0 {
		return errors.New(messages.Get("NO_MORE_USERDEVICES"))
	}
	if !contains(devices, position) {
		position = devices[0]
	}
	_, err = xenapi.VBD.Create(s, xenapi.VBDRecord{VM: vm, VDI: vdi, Userdevice: position, Bootable: bootable, Mode: xenapi.VbdModeRW,
		Type: xenapi.VbdTypeDisk, Unpluggable: true, OtherConfig: map[string]string{"owner": "true"}, QosAlgorithmParams: map[string]string{}})
	return err
}

// addVMHint is AddVMHint: the VM's reference as the disk's placement hint.
func addVMHint(s *xenapi.Session, vdi xenapi.VDIRef, vm xenapi.VMRef) error {
	_ = xenapi.VDI.RemoveFromSmConfig(s, vdi, "vmhint")
	return xenapi.VDI.AddToSmConfig(s, vdi, "vmhint", string(vm))
}

// addNetworks is AddNetworks: the template's interfaces go, one per chosen network is created on a free device.
func addNetworks(s *xenapi.Session, vm xenapi.VMRef, rec xenapi.VMRecord, spec NewVM) error {
	for _, v := range rec.VIFs {
		if err := xenapi.VIF.Destroy(s, v); err != nil {
			return err
		}
	}
	for i, n := range spec.Networks {
		net, err := xenapi.Network.GetByUUID(s, n.Network)
		if err != nil {
			return err
		}
		devices, err := xenapi.VM.GetAllowedVIFDevices(s, vm)
		if err != nil {
			return err
		}
		if len(devices) == 0 {
			return nil // more interfaces than the server allows: the rest are not created, as the C# app
		}
		want := n.Device
		if want == "" {
			want = itoa64(int64(i))
		}
		device := devices[0]
		if contains(devices, want) {
			device = want
		}
		params := n.QosParams
		if params == nil {
			params = map[string]string{}
		}
		if _, err := xenapi.VIF.Create(s, xenapi.VIFRecord{Device: device, Network: net, VM: vm, MTU: 1500, LockingMode: xenapi.VifLockingModeNetworkDefault,
			OtherConfig: map[string]string{}, QosAlgorithmType: n.QosType, QosAlgorithmParams: params}); err != nil {
			return err
		}
	}
	return nil
}

// destroyVM is VMDestroyAction.DestroyVM(vm, true): the VM and the disks it owns; false when the VM is still there.
func destroyVM(s *xenapi.Session, vm xenapi.VMRef) bool {
	rec, err := xenapi.VM.GetRecord(s, vm)
	if err != nil {
		var api *xenapi.APIError
		return errors.As(err, &api) && api.Message == "HANDLE_INVALID" // already gone
	}
	var vdis []xenapi.VDIRef
	for _, b := range rec.VBDs {
		if vbd, err := xenapi.VBD.GetRecord(s, b); err == nil && vbd.OtherConfig["owner"] != "" && vbd.VDI != nullRef {
			vdis = append(vdis, vbd.VDI)
		}
	}
	if rec.PowerState != xenapi.VMPowerStateHalted {
		_ = xenapi.VM.HardShutdown(s, vm)
	}
	if err := xenapi.VM.Destroy(s, vm); err != nil {
		return false
	}
	// the VM is gone: its disks are only reachable from here, so each is retried before giving up
	for _, v := range vdis {
		for i := 0; i < 5; i++ {
			err := xenapi.VDI.Destroy(s, v)
			var api *xenapi.APIError
			if err == nil || (errors.As(err, &api) && api.Message == "HANDLE_INVALID") {
				break
			}
			time.Sleep(2 * time.Second)
		}
	}
	return true
}

// removeLostClone removes the clone of a failed attempt by its unique hidden name; when the clone's outcome is unknown
// (the task could not be followed) it waits up to two minutes for the clone to appear.
func removeLostClone(s *xenapi.Session, hidden string, uncertain bool) {
	for i := 0; ; i++ {
		found, err := xenapi.VM.GetByNameLabel(s, hidden)
		if err == nil && len(found) > 0 {
			removed := true
			for _, v := range found {
				removed = destroyVM(s, v) && removed
			}
			if removed {
				return
			}
			uncertain = true // a clone that could not be removed yet: try again until the deadline
		}
		if !uncertain || i >= 60 {
			return
		}
		time.Sleep(2 * time.Second)
	}
}

// customWithNoDVD is Helpers.CustomWithNoDVD: a template the user made that has no DVD drive.
func customWithNoDVD(s *xenapi.Session, t xenapi.VMRecord) bool {
	if _, def := t.OtherConfig["default_template"]; def {
		return false
	}
	for _, b := range t.VBDs {
		if typ, err := xenapi.VBD.GetType(s, b); err == nil && typ == xenapi.VbdTypeCD {
			return false
		}
	}
	return true
}

// rhelRepo is NormalizeRepoUrlForRHEL: nfs://host/path becomes nfs:host:/path.
func rhelRepo(url string) string {
	if rest, ok := strings.CutPrefix(url, "nfs://"); ok {
		if i := strings.Index(rest, "/"); i > 0 {
			return "nfs:" + rest[:i] + ":" + rest[i:]
		}
	}
	return url
}

// networkFirst is GetBootOrderNetworkFirst: "n" moved (or added) to the front of the boot order.
func networkFirst(order string) string {
	order = strings.ToLower(order)
	if order == "" {
		return "ncd"
	}
	return "n" + strings.Replace(order, "n", "", 1)
}

// hasRecommendation reads a boolean restriction of the template's recommendations (VM.GetBoolRestrictionValue).
func hasRecommendation(recommendations, field string) bool {
	var doc struct {
		R []struct {
			Field string `xml:"field,attr"`
			Value string `xml:"value,attr"`
		} `xml:"restriction"`
	}
	if xml.Unmarshal([]byte(recommendations), &doc) != nil {
		return false
	}
	for _, r := range doc.R {
		if r.Field == field {
			return r.Value == "true"
		}
	}
	return false
}

// vendorDeviceRestricted is FeatureForbidden(Host.RestrictVendorDevice): restricted unless the licence says "false".
func vendorDeviceRestricted(s *xenapi.Session) bool {
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return true
	}
	for _, h := range hosts {
		if v, ok := h.LicenseParams["restrict_pci_device_for_auto_update"]; !ok || v != "false" {
			return true
		}
	}
	return false
}

func itoa64(n int64) string { return strconv.FormatInt(n, 10) }
