package ops

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/trust"
	"xcpng-center/internal/xapi"
)

// OVF/OVA appliance export: ExportApplianceAction (ExportApplianceAction.cs, Export.cs) with the Export wizard's
// checks. The package is made in a temporary folder beside the target and renamed to "<dir>/<name>" once complete, so
// a failed export leaves nothing at the target.
// shortcut: no encryption (hidden in the C# wizard too), no progress bar or Cancel (as the XVA export).

// ApplianceExport is what the Export wizard asks for the OVF path (ExportApplianceAction's arguments).
type ApplianceExport struct {
	Dir      string   `json:"dir"`
	Name     string   `json:"name"`
	VMs      []string `json:"vms"`   // the VMs' uuids, in order
	Eulas    []string `json:"eulas"` // the EULA files
	Manifest bool     `json:"manifest"`
	OVA      bool     `json:"ova"`
	Compress bool     `json:"compress"`
	Verify   bool     `json:"verify"`
	// Sign the package with the certificate file (a .pfx) and its private key's password, typed in the dialog: used
	// to open the file and never kept nor logged
	Sign        bool   `json:"sign"`
	Certificate string `json:"certificate"`
	Password    string `json:"password"`
}

// diskMaxSize is SR.DISK_MAX_SIZE (2 TiB): ExportSelectVMsPage.CheckDiskSizeForTransfer refuses a bigger VM for OVF.
const diskMaxSize = 2 << 40

// ApplianceTargetError is ExportAppliancePage's checks (CheckPathValid, CheckDestinationFolderExists,
// CheckApplianceExists): "" when the package can be made as "<dir>/<name>". The C# check refuses a folder holding
// "<name>.ovf" or "<name>.ova"; here the folder must not exist or be empty, since the package replaces it whole. For an
// XVA (xva) it is the file "<dir>/<name>.xva" that must not exist (ApplianceExistsCheck.VerifyXva).
func ApplianceTargetError(dir, name string, xva bool) string {
	if why := rules.ApplianceNameError(name); why != "" {
		return why
	}
	if !filepath.IsAbs(dir) {
		return messages.Get("EXPORT_APPLIANCE_PAGE_ERROR_INVALID_DIR")
	}
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return messages.Get("ERROR_DESTINATION_DIR_NON_EXIST")
	}
	if xva {
		if !strings.HasSuffix(name, ".xva") {
			name += ".xva"
		}
		if _, err := os.Lstat(filepath.Join(dir, name)); !os.IsNotExist(err) {
			return messages.Get("EXPORT_APPLIANCE_PAGE_ERROR_APP_EXISTS")
		}
		return ""
	}
	if entries, err := os.ReadDir(filepath.Join(dir, name)); err == nil && len(entries) > 0 || err != nil && !os.IsNotExist(err) {
		return messages.Get("EXPORT_APPLIANCE_PAGE_ERROR_APP_EXISTS")
	}
	return ""
}

// exportDisk is a disk of a VM to export: its VDI and where its spec is.
type exportDisk struct {
	vm, disk int
	vdi      xenapi.VDIRecord
}

// readExportVM is Export.cs' read of the VM (before any disk is exported): its spec and its disks' VDIs.
func readExportVM(s *xenapi.Session, uuid string) (rules.ExportVM, []xenapi.VDIRecord, error) {
	var spec rules.ExportVM
	ref, err := xenapi.VM.GetByUUID(s, uuid)
	if err != nil {
		return spec, nil, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return spec, nil, xapi.Sanitize(err)
	}
	switch {
	case vm.IsATemplate || vm.IsASnapshot || vm.IsControlDomain:
		return spec, nil, errors.New("Only VMs can be exported.")
	case vm.PowerState != xenapi.VMPowerStateHalted && vm.PowerState != xenapi.VMPowerStateSuspended:
		return spec, nil, errors.New(messages.Format(messages.Get("ERROR_VM_NOT_HALTED"), vm.NameLabel))
	case len(vm.CurrentOperations) > 0 || !contains(vm.AllowedOperations, xenapi.VMOperationsExport): // IsVmExportable
		return spec, nil, fmt.Errorf("The VM '%s' cannot be exported now.", vm.NameLabel)
	}
	spec = rules.ExportVM{Name: vm.NameLabel, Description: vm.NameDescription, HVM: vm.HVMBootPolicy != "", Arch: vm.Domarch,
		VCPUs: int64(vm.VCPUsAtStartup), VCPUsMax: int64(vm.VCPUsMax), Memory: int64(vm.MemoryDynamicMax),
		Order: int64(vm.Order), StartDelay: int64(vm.StartDelay), ShutdownDelay: int64(vm.ShutdownDelay),
		HVMBootParams: vm.HVMBootParams, Platform: vm.Platform, NVRAM: vm.NVRAM, HVMBootPolicy: vm.HVMBootPolicy,
		PVArgs: vm.PVArgs, PVBootloader: vm.PVBootloader, PVBootloaderArgs: vm.PVBootloaderArgs, PVKernel: vm.PVKernel,
		PVLegacyArgs: vm.PVLegacyArgs, PVRamdisk: vm.PVRamdisk, Recommendations: vm.Recommendations, ReferenceLabel: vm.ReferenceLabel,
		HVMShadowMultiplier: vm.HVMShadowMultiplier, HardwarePlatformVersion: int64(vm.HardwarePlatformVersion), HasVendorDevice: vm.HasVendorDevice}
	if vm.GuestMetrics != nullRef { // TRY TO ID OS
		if v, err := xenapi.VMGuestMetrics.GetOsVersion(s, vm.GuestMetrics); err == nil {
			spec.OSName, spec.OSMajor, spec.OSMinor = v["name"], v["major"], v["minor"]
		}
	}
	// ADD NETWORKS, in device order
	type nic struct {
		device int
		rules.ExportNIC
	}
	// a device that cannot be read fails the export: a package without it would look complete
	var nics []nic
	for _, v := range vm.VIFs {
		vif, err := xenapi.VIF.GetRecord(s, v)
		if err != nil {
			return spec, nil, xapi.Sanitize(err)
		}
		net, err := xenapi.Network.GetRecord(s, vif.Network)
		if err != nil {
			return spec, nil, xapi.Sanitize(err)
		}
		d, _ := strconv.Atoi(vif.Device)
		nics = append(nics, nic{d, rules.ExportNIC{Network: net.UUID, NetworkName: net.NameLabel, NetworkDescription: net.NameDescription, MAC: vif.MAC}})
	}
	slices.SortStableFunc(nics, func(a, b nic) int { return a.device - b.device })
	for _, n := range nics {
		spec.NICs = append(spec.NICs, n.ExportNIC)
	}
	// EXPORT DISKS, in device order
	var vbds []xenapi.VBDRecord
	for _, b := range vm.VBDs {
		vbd, err := xenapi.VBD.GetRecord(s, b)
		if err != nil {
			return spec, nil, xapi.Sanitize(err)
		}
		vbds = append(vbds, vbd)
	}
	slices.SortStableFunc(vbds, func(a, b xenapi.VBDRecord) int {
		x, _ := strconv.Atoi(a.Userdevice)
		y, _ := strconv.Atoi(b.Userdevice)
		return x - y
	})
	var vdis []xenapi.VDIRecord
	var size int64
	for _, vbd := range vbds {
		if n, err := strconv.Atoi(vbd.Userdevice); err != nil || n < 0 || n > 15 { // SetTargetDeviceInRASD
			return spec, nil, fmt.Errorf("Device number %s is less than zero or greater than 15.", vbd.Userdevice)
		}
		if vbd.Type == xenapi.VbdTypeCD {
			spec.CDs = append(spec.CDs, vbd.Userdevice)
			continue
		}
		if vbd.VDI == nullRef {
			continue
		}
		vdi, err := xenapi.VDI.GetRecord(s, vbd.VDI)
		if err != nil {
			return spec, nil, xapi.Sanitize(err)
		}
		name := vdi.NameLabel
		if name == "" {
			name = fmt.Sprintf("Hard Disk Image %d", len(vdis)) // RASD_19_CAPTION and the disk's index
		}
		spec.Disks = append(spec.Disks, rules.ExportDisk{Name: name, Description: vdi.NameDescription, Device: vbd.Userdevice,
			Bootable: vbd.Bootable, Capacity: int64(vdi.VirtualSize)})
		vdis = append(vdis, vdi)
		size += int64(vdi.VirtualSize)
	}
	if size > diskMaxSize {
		return spec, nil, errors.New(messages.Format(messages.Get("EXPORT_ERROR_EXCEEDS_MAX_SIZE_VDI_OVA_OVF"), "2 TB"))
	}
	return spec, vdis, nil
}

// ExportAppliance is ExportApplianceAction: the VMs' disks exported as VHD (export_raw_vdi) and verified when asked,
// compressed one by one for a compressed OVF, then the descriptor, the manifest and for an OVA the archive. It returns
// the package's file (the .ovf, or the .ova / .ova.gz).
func ExportAppliance(c *xapi.Connection, policy *trust.Policy, e ApplianceExport) (string, error) {
	s := c.Session
	if err := permitted(s, "task.create", "http/get_export_raw_vdi"); err != nil {
		return "", err
	}
	if why := ApplianceTargetError(e.Dir, e.Name, false); why != "" {
		return "", errors.New(why)
	}
	if len(e.VMs) == 0 {
		return "", errors.New("Select the VMs to export.")
	}
	eulas, err := rules.ReadEulas(e.Eulas)
	if err != nil {
		return "", err
	}
	var signer *rules.Signer // ExportOptionsPage's checks, before any disk is exported
	if e.Sign {
		if signer, err = rules.LoadSigner(e.Certificate, e.Password, nil); err != nil {
			return "", err
		}
	}
	specs := make([]rules.ExportVM, len(e.VMs))
	var disks []exportDisk
	for i, uuid := range e.VMs {
		spec, vdis, err := readExportVM(s, uuid)
		if err != nil {
			return "", err
		}
		specs[i] = spec
		for j, vdi := range vdis {
			disks = append(disks, exportDisk{i, j, vdi})
		}
	}
	tmp, err := os.MkdirTemp(e.Dir, "."+e.Name+".*.tmp")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp) // gone once renamed to the target; else whatever a failed export wrote
	used := map[string]bool{}
	for _, d := range disks {
		file := d.vdi.UUID + ".vhd"
		for n := 1; used[file]; n++ { // the "VHD Name collision" of a disk in two VMs
			file = fmt.Sprintf("%s_%d.vhd", d.vdi.UUID, n)
		}
		used[file] = true
		size, err := exportVDI(c, policy, d.vdi, filepath.Join(tmp, file), e.Verify)
		if err != nil {
			return "", err
		}
		disk := &specs[d.vm].Disks[d.disk]
		disk.Href, disk.Size = file, size
		if e.Compress && !e.OVA { // OVF.CompressFiles
			if disk.Size, err = rules.GzipFile(filepath.Join(tmp, file)); err != nil {
				return "", err
			}
			disk.Href, disk.Compressed = file+".gz", true
		}
	}
	file, err := rules.WritePackage(tmp, e.Name, specs, eulas, e.Manifest, e.OVA, e.Compress, signer)
	if err != nil {
		return "", err
	}
	target := filepath.Join(e.Dir, e.Name)
	if why := ApplianceTargetError(e.Dir, e.Name, false); why != "" { // made meanwhile
		return "", errors.New(why)
	}
	_ = os.Remove(target) // an empty folder of that name (Windows renames no folder onto another)
	if err := os.Rename(tmp, target); err != nil {
		return "", err
	}
	return filepath.Join(target, file), nil
}

// exportVDI is Export.cs' EXPORT DISKS for one disk: GET /export_raw_vdi?format=vhd with a task into file, flushed to the
// disk, then checked when asked (EXPORTING_VDI_VERIFICATION). It returns the file's size.
func exportVDI(c *xapi.Connection, policy *trust.Policy, vdi xenapi.VDIRecord, file string, verify bool) (int64, error) {
	s := c.Session
	address, port := c.Host, c.Port
	if sr, err := xenapi.SR.GetRecord(s, vdi.SR); err == nil {
		address, port = storageHost(c, sr) // a local SR's disks are on its own server
	}
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644) // before the task, which nothing would end
	if err != nil {
		return 0, err
	}
	task, err := xenapi.Task.Create(s, "export_raw_vdi_task", fmt.Sprintf("Exporting VDI %s to %s", vdi.UUID, filepath.Base(file)))
	if err != nil {
		f.Close()
		return 0, xapi.Sanitize(err)
	}
	err = func() error {
		body, err := get(s, policy, address, port, "/export_raw_vdi?"+url.Values{"task_id": {string(task)}, "vdi": {vdi.UUID}, "format": {"vhd"}}.Encode())
		if err != nil {
			return err
		}
		defer body.Close()
		w := bufio.NewWriterSize(f, 1<<20)
		if _, err := io.Copy(w, body); err != nil {
			return err
		}
		if err := w.Flush(); err != nil {
			return err
		}
		return f.Sync()
	}()
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err = finishTask(s, task, err, nil); err != nil {
		return 0, err
	}
	st, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	if verify {
		r, err := os.Open(file)
		if err != nil {
			return 0, err
		}
		defer r.Close()
		if err := rules.VerifyVHD(r, st.Size()); err != nil {
			return 0, rules.VerifyVHDError(filepath.Base(file), err)
		}
	}
	return st.Size(), nil
}
