package rules

import (
	"archive/tar"
	"bufio"
	"bytes"
	"cmp"
	"compress/gzip"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	"xcpng-center/internal/messages"
)

// OVF/OVA export (ExportApplianceAction): Export.cs makes one envelope per VM with XenOvf's OVF.CreateEnvelope,
// AddVirtualSystem, SetCPUs, SetMemory, AddNetwork, AddCDROM, AddDisk, AddOtherSystemSettingData and AddStartupSection;
// OVF.Merge joins them, OVF.AddEula adds the EULAs, OVF.FinalizeEnvelope makes a one-VM collection a VirtualSystem and
// OVF.SaveAs writes it. The element namespaces are those of XenOvf's Definitions (ovfEnvelope.cs), which the C#
// XmlSerializer needs to read the descriptor back. The texts are XenOvf's Content.resx.

const (
	nsOVF  = "http://schemas.dmtf.org/ovf/envelope/1"
	nsRASD = "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_ResourceAllocationSettingData"
	nsVSSD = "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_VirtualSystemSettingData"
	nsXen  = "http://schemas.citrix.com/ovf/envelope/1"
	// vhdFormat is OVF.WIN_FILE_FORMAT_URI, the DiskSection format of a VHD
	vhdFormat = "http://www.microsoft.com/technet/virtualserver/downloads/vhdspec.mspx"
)

// ExportDisk is a disk of an exported VM: its file in the package and the VBD it came from (OVF.AddDisk).
type ExportDisk struct {
	Href        string
	Name        string
	Description string
	Device      string // the VBD's userdevice (SetTargetDeviceInRASD)
	Bootable    bool
	Size        int64 // the file's size
	Capacity    int64 // the VDI's virtual size
	Compressed  bool  // gzip (OVF.CompressFiles)
}

// ExportNIC is an interface of an exported VM (OVF.AddNetwork): its network's uuid, name and description, and its MAC.
type ExportNIC struct {
	Network, NetworkName, NetworkDescription, MAC string
}

// ExportVM is what Export.cs reads of a VM.
type ExportVM struct {
	Name, Description string
	HVM               bool
	Arch              string // domarch
	OSName            string // the guest's os_version "name" ("" without one)
	OSMajor, OSMinor  string
	VCPUs, VCPUsMax   int64
	Memory            int64 // memory_dynamic_max
	Order             int64
	StartDelay        int64
	ShutdownDelay     int64
	NICs              []ExportNIC
	CDs               []string // the DVD drives' userdevices
	Disks             []ExportDisk
	// the Xen specifics (ADD XEN SPECIFICS)
	HVMBootParams, Platform, NVRAM                                  map[string]string
	HVMBootPolicy, PVArgs, PVBootloader, PVBootloaderArgs, PVKernel string
	PVLegacyArgs, PVRamdisk, Recommendations, ReferenceLabel        string
	HVMShadowMultiplier                                             float64
	HardwarePlatformVersion                                         int64
	HasVendorDevice                                                 bool
}

// newID is Guid.NewGuid().ToString(): a random (version 4) UUID.
func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6], b[8] = b[6]&0x0f|0x40, b[8]&0x3f|0x80
	h := hex.EncodeToString(b[:])
	return h[:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:]
}

type ovfWriter struct{ bytes.Buffer }

func (w *ovfWriter) esc(s string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

// el writes <name>text</name>.
func (w *ovfWriter) el(name, text string) {
	fmt.Fprintf(w, "<%s>%s</%s>", name, w.esc(text), name)
}

// rasd is one ResourceAllocationSettingData item; its elements are written in the schema's (alphabetical) order.
type rasd map[string]string

var rasdOrder = []string{"Address", "AddressOnParent", "AllocationUnits", "AutomaticAllocation", "Caption", "Connection",
	"ConsumerVisibility", "Description", "ElementName", "HostResource", "InstanceID", "Limit", "MappingBehavior",
	"ResourceSubType", "ResourceType", "VirtualQuantity", "Weight"}

func (w *ovfWriter) item(r rasd) {
	w.WriteString(`<ovf:Item ovf:required="false">`)
	for _, k := range rasdOrder {
		if v, ok := r[k]; ok {
			w.el("rasd:"+k, v)
		}
	}
	w.WriteString("</ovf:Item>")
}

// joinDict is string.Join(";", dict.Select(k=v)), by key (a Go map has no order).
func joinDict(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = k + "=" + m[k]
	}
	return strings.Join(parts, ";")
}

// xenSettings are the VirtualSystemOtherConfigurationData of ADD XEN SPECIFICS: name, value and description.
// shortcut: no "vgpu" nor "pvssite" (the app has no vGPU or PVS); add them with those features.
func xenSettings(vm ExportVM) [][3]string {
	hypervisor := "XCP-ng Specific Information" // OTHER_SYSTEM_SETTING_DESCRIPTION_1 with BrandManager.ProductBrand
	var out [][3]string
	add := func(name, value, info string) { out = append(out, [3]string{name, value, info}) }
	if len(vm.HVMBootParams) > 0 {
		add("HVM_boot_params", joinDict(vm.HVMBootParams), "HVM Boot Parameters")
	}
	if vm.HVMBootPolicy != "" {
		add("HVM_boot_policy", vm.HVMBootPolicy, "Boot Order: CD/DVD - Network - Disk")
	}
	if vm.HVMShadowMultiplier != 1.0 {
		add("HVM_shadow_multiplier", strconv.FormatFloat(vm.HVMShadowMultiplier, 'g', -1, 64), hypervisor)
	}
	if len(vm.Platform) > 0 {
		add("platform", joinDict(vm.Platform), "XCP-ng Platform Information")
	}
	if len(vm.NVRAM) > 0 {
		add("NVRAM", joinDict(vm.NVRAM), "Non-Volatile Random-Access Memory")
	}
	for _, s := range [][2]string{{"PV_args", vm.PVArgs}, {"PV_bootloader", vm.PVBootloader}, {"PV_bootloader_args", vm.PVBootloaderArgs},
		{"PV_kernel", vm.PVKernel}, {"PV_legacy_args", vm.PVLegacyArgs}, {"PV_ramdisk", vm.PVRamdisk}} {
		if s[1] != "" {
			add(s[0], s[1], hypervisor)
		}
	}
	if vm.HardwarePlatformVersion >= 0 {
		add("hardware_platform_version", strconv.FormatInt(vm.HardwarePlatformVersion, 10), hypervisor)
	}
	if vm.Recommendations != "" {
		add("recommendations", vm.Recommendations, hypervisor)
	}
	if vm.ReferenceLabel != "" {
		add("reference_label", vm.ReferenceLabel, hypervisor)
	}
	if vm.HasVendorDevice { // another name, so that it is not read back as a plain field of the wrong type
		add("VM_has_vendor_device", "True", hypervisor)
	}
	return out
}

func (w *ovfWriter) eulas(eulas []string) {
	for _, e := range eulas {
		fmt.Fprintf(w, `<ovf:EulaSection xenovf:Id="%s"><ovf:Info>License Agreement</ovf:Info>`, newID())
		w.el("ovf:License", e)
		w.WriteString("</ovf:EulaSection>")
	}
}

// system writes the VirtualSystem of one VM (Export.cs). Its interfaces are in device order: the C# AddRasdToAllVHS
// puts each new item first, which reverses the interfaces' order on import.
func (w *ovfWriter) system(vm ExportVM, id string, diskIDs []string, eulas []string) {
	fmt.Fprintf(w, `<ovf:VirtualSystem ovf:id="%s"><ovf:Info>XenOVF Generated</ovf:Info>`, id)
	w.el("ovf:Name", vm.Name)
	// TRY TO ID OS. shortcut: the OS id is 1 (Other) with a guest agent: ValueMaps.OperatingSystem matches the CIM
	// names exactly, which a guest's name hardly ever is
	osID, desc, info := 0, "XenOVF Created Export", "Guest Operating System"
	if name, _, _ := strings.Cut(vm.OSName, "|"); vm.OSName != "" {
		osID, desc = 1, name
		if vm.OSMajor != "" && vm.OSMinor != "" {
			info = "Version " + vm.OSMajor + "." + vm.OSMinor // SECTION_OPERATINGSYSTEM_VERSION
		}
	}
	fmt.Fprintf(w, `<ovf:OperatingSystemSection ovf:id="%d" xenovf:Id="%s">`, osID, newID())
	w.el("ovf:Info", info)
	w.el("ovf:Description", desc)
	w.WriteString("</ovf:OperatingSystemSection>")
	// FinalizeEnvelope's roll-up
	mb := vm.Memory / (1 << 20)
	fmt.Fprintf(w, `<ovf:VirtualHardwareSection xenovf:Id="%s">`, newID())
	w.el("ovf:Info", fmt.Sprintf("Virtual Hardware Requirements: %d MB RAM; %d CPU(s), %d Disk(s), %d Network(s)", mb, vm.VCPUs, len(vm.Disks), len(vm.NICs)))
	pv := "xen"
	if vm.HVM {
		pv = "hvm"
	}
	w.WriteString("<ovf:System>")
	w.el("vssd:Caption", "Xen Virtual Machine")
	if vm.Description != "" {
		w.el("vssd:Description", vm.Description)
	}
	w.el("vssd:ElementName", vm.Name)
	w.el("vssd:InstanceID", newID())
	w.el("vssd:VirtualSystemIdentifier", newID())
	w.el("vssd:VirtualSystemType", pv+"-3.0-"+cmp.Or(vm.Arch, "unknown"))
	w.WriteString("</ovf:System>")
	w.item(rasd{"AllocationUnits": "Processor Cores", "AutomaticAllocation": "true", "Caption": "Processor", "ConsumerVisibility": "0",
		"Description": "Processor", "ElementName": "Processor", "InstanceID": newID(), "Limit": strconv.FormatInt(vm.VCPUsMax, 10),
		"MappingBehavior": "0", "ResourceType": "3", "VirtualQuantity": strconv.FormatInt(vm.VCPUs, 10), "Weight": "100"})
	w.item(rasd{"AllocationUnits": "byte * 2^20", "AutomaticAllocation": "true", "Caption": "Memory", "ConsumerVisibility": "2",
		"Description": "Memory Setting Data", "ElementName": "Memory", "InstanceID": newID(), "Limit": "128", "MappingBehavior": "0",
		"ResourceSubType": "Machine Memory", "ResourceType": "4", "VirtualQuantity": strconv.FormatInt(mb, 10), "Weight": "0"})
	for _, n := range vm.NICs {
		r := rasd{"AllocationUnits": "Ports", "AutomaticAllocation": "true", "Caption": "Network Adapter", "Connection": n.Network,
			"ConsumerVisibility": "3", "Description": cmp.Or(n.NetworkDescription, "Network Adapter"), "ElementName": n.NetworkName,
			"InstanceID": newID(), "Limit": "1", "MappingBehavior": "0", "ResourceType": "10", "VirtualQuantity": "1", "Weight": "0"}
		if n.MAC != "" {
			r["Address"] = n.MAC
		}
		w.item(r)
	}
	for _, dev := range vm.CDs {
		w.item(rasd{"AddressOnParent": dev, "AllocationUnits": "DVD Drives", "AutomaticAllocation": "true", "Caption": "CD/DVD Drive",
			"Connection": "device=" + dev, "ConsumerVisibility": "3", "Description": "Xen CD/DVD Drive Setting Data.", "ElementName": "CD/DVD Drive",
			"InstanceID": newID(), "Limit": "1", "MappingBehavior": "0", "ResourceType": "16", "VirtualQuantity": "1", "Weight": "0"})
	}
	for i, d := range vm.Disks {
		id := diskIDs[i]
		w.item(rasd{"AddressOnParent": d.Device, "AllocationUnits": "Disk", "AutomaticAllocation": "true", "Caption": "Hard Disk Image",
			"Connection": id + ",device=" + d.Device, "ConsumerVisibility": "3", "Description": d.Description, "ElementName": d.Name,
			"HostResource": "ovf:/disk/" + id, "InstanceID": id, "Limit": "1", "MappingBehavior": "0", "ResourceSubType": "Virtual Hard Disk Image",
			"ResourceType": "19", "VirtualQuantity": "1", "Weight": "100"})
	}
	for _, s := range xenSettings(vm) {
		fmt.Fprintf(w, `<xenovf:VirtualSystemOtherConfigurationData id="%s" Name="%s">`, newID(), w.esc(s[0]))
		w.el("xenovf:Info", s[2])
		w.el("xenovf:Value", s[1])
		w.WriteString("</xenovf:VirtualSystemOtherConfigurationData>")
	}
	w.WriteString("</ovf:VirtualHardwareSection>")
	w.eulas(eulas)
	w.WriteString("</ovf:VirtualSystem>")
}

// ApplianceDescriptor is the .ovf of the package name for the VMs (in order) and the EULAs' texts: Export.cs for each
// VM, then OVF.Merge, OVF.AddEula and OVF.FinalizeEnvelope.
func ApplianceDescriptor(name string, vms []ExportVM, eulas []string) []byte {
	var w ovfWriter
	w.WriteString(xml.Header)
	fmt.Fprintf(&w, `<ovf:Envelope ovf:version="1.0.0" xml:lang="en-US" xenovf:Name="%s" xenovf:id="%s" xmlns:ovf="%s" xmlns:rasd="%s" xmlns:vssd="%s" xmlns:xenovf="%s">`,
		w.esc(name), newID(), nsOVF, nsRASD, nsVSSD, nsXen)
	sysIDs := make([]string, len(vms))
	diskIDs := make([][]string, len(vms))
	w.WriteString("<ovf:References>")
	for i, vm := range vms {
		sysIDs[i] = newID()
		for _, d := range vm.Disks {
			id := newID()
			diskIDs[i] = append(diskIDs[i], id)
			fmt.Fprintf(&w, `<ovf:File ovf:id="%s" ovf:href="%s" ovf:size="%d"`, id, w.esc(d.Href), d.Size)
			if d.Compressed {
				w.WriteString(` ovf:compression="gzip"`)
			}
			w.WriteString("/>")
		}
	}
	w.WriteString("</ovf:References>")
	if slices.ContainsFunc(vms, func(v ExportVM) bool { return len(v.Disks) > 0 }) {
		fmt.Fprintf(&w, `<ovf:DiskSection xenovf:Id="%s"><ovf:Info>Available Disks</ovf:Info>`, newID())
		for i, vm := range vms {
			for j, d := range vm.Disks {
				fmt.Fprintf(&w, `<ovf:Disk ovf:diskId="%s" ovf:fileRef="%s" ovf:capacity="%d" ovf:format="%s" xenovf:isBootable="%t"/>`,
					diskIDs[i][j], diskIDs[i][j], d.Capacity, vhdFormat, d.Bootable)
			}
		}
		w.WriteString("</ovf:DiskSection>")
	}
	seen := map[string]bool{} // Merge: one network per uuid
	var nets bytes.Buffer
	for _, vm := range vms {
		for _, n := range vm.NICs {
			if !seen[n.Network] {
				seen[n.Network] = true
				fmt.Fprintf(&nets, `<ovf:Network ovf:name="%s">`, w.esc(n.Network))
				nets.WriteString("<ovf:Description>" + w.esc(cmp.Or(n.NetworkDescription, "Network Adapter")) + "</ovf:Description></ovf:Network>")
			}
		}
	}
	if nets.Len() > 0 {
		fmt.Fprintf(&w, `<ovf:NetworkSection xenovf:Id="%s"><ovf:Info>Available Networks</ovf:Info>`, newID())
		w.Write(nets.Bytes())
		w.WriteString("</ovf:NetworkSection>")
	}
	fmt.Fprintf(&w, `<ovf:StartupSection xenovf:Id="%s"><ovf:Info>Startup Options</ovf:Info>`, newID())
	for i, vm := range vms {
		fmt.Fprintf(&w, `<ovf:Item ovf:id="%s" ovf:order="%d" ovf:startDelay="%d" ovf:waitingForGuest="false" ovf:stopDelay="%d" ovf:startAction="powerOn" ovf:stopAction="powerOff"/>`,
			sysIDs[i], vm.Order, vm.StartDelay, vm.ShutdownDelay)
	}
	w.WriteString("</ovf:StartupSection>")
	if len(vms) == 1 { // FinalizeEnvelope: one VM is the envelope's VirtualSystem, with the EULAs
		w.system(vms[0], sysIDs[0], diskIDs[0], eulas)
	} else {
		fmt.Fprintf(&w, `<ovf:VirtualSystemCollection ovf:id="%s"><ovf:Info>OVF/OVA Packages</ovf:Info>`, newID())
		w.eulas(eulas)
		for i, vm := range vms {
			w.system(vm, sysIDs[i], diskIDs[i], nil)
		}
		w.WriteString("</ovf:VirtualSystemCollection>")
	}
	w.WriteString("</ovf:Envelope>\n")
	return w.Bytes()
}

// ApplianceNameError is ExportAppliancePage.CheckPathValid's name check (PathValidator.IsFileNameValid): "" when the
// package can be called name.
func ApplianceNameError(name string) string {
	if strings.TrimSpace(name) == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\:*?"<>|`) ||
		strings.ContainsFunc(name, func(r rune) bool { return r < 32 }) {
		return messages.Get("EXPORT_APPLIANCE_PAGE_ERROR_INALID_APP")
	}
	return ""
}

// WriteManifest is ExportApplianceAction.Manifest: "<name>.mf" in dir with the SHA-256 of each file (FileDigest's
// "SHA256(file)= digest" lines), the descriptor first.
func WriteManifest(dir, name string, files []string) error {
	var mf strings.Builder
	for _, f := range files {
		in, err := os.Open(filepath.Join(dir, f))
		if err != nil {
			return err
		}
		h := sha256.New()
		_, err = io.Copy(h, in)
		in.Close()
		if err != nil {
			return err
		}
		fmt.Fprintf(&mf, "SHA256(%s)= %x\n", f, h.Sum(nil))
	}
	return writeSynced(filepath.Join(dir, name+".mf"), func(w io.Writer) error { _, err := io.WriteString(w, mf.String()); return err })
}

// writeSynced makes file with fill's output, flushed to the disk.
func writeSynced(file string, fill func(io.Writer) error) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	bw := bufio.NewWriterSize(f, 1<<20)
	err = fill(bw)
	if err == nil {
		err = bw.Flush()
	}
	if err == nil {
		err = f.Sync()
	}
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	return err
}

// WriteOVA is OVF.ConvertOVFtoOVA: "<name>.ova" (".ova.gz" when compressed) in dir, a tar of the files in their order
// (the descriptor, the manifest, the disks), each file removed once it is in the archive (deleteOriginalFiles). It
// returns the archive's name.
func WriteOVA(dir, name string, files []string, compress bool) (string, error) {
	ova := name + ".ova"
	if compress {
		ova += ".gz"
	}
	err := writeSynced(filepath.Join(dir, ova), func(out io.Writer) error {
		var z *gzip.Writer
		if compress {
			z = gzip.NewWriter(out)
			out = z
		}
		tw := tar.NewWriter(out)
		for _, f := range files {
			if err := addToTar(tw, filepath.Join(dir, f), f); err != nil {
				return err
			}
			_ = os.Remove(filepath.Join(dir, f))
		}
		if err := tw.Close(); err != nil {
			return err
		}
		if z != nil {
			return z.Close()
		}
		return nil
	})
	return ova, err
}

func addToTar(tw *tar.Writer, file, name string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return err
	}
	// ustar, or for a file of 8 GiB or more a PAX size record (as GNU tar): the C# app's SharpZipLib 0.85 neither writes
	// nor reads such an entry (it cuts the size to 11 octal digits), so it cannot open that OVA; this app and others can.
	// A name longer than ustar's 100 bytes, or not ASCII, is a GNU long name, which SharpZipLib writes and reads too.
	h := &tar.Header{Name: name, Mode: 0o644, Size: st.Size(), ModTime: st.ModTime().Truncate(time.Second), Typeflag: tar.TypeReg}
	if len(name) > 100 || strings.ContainsFunc(name, func(r rune) bool { return r >= 0x80 }) {
		h.Format = tar.FormatGNU
	}
	if err := tw.WriteHeader(h); err != nil {
		return err
	}
	_, err = io.Copy(tw, f)
	return err
}

// GzipFile is CompressionFactory.CompressFile for OVF.CompressFiles: src as "<src>.gz", src removed; the new size.
func GzipFile(src string) (int64, error) {
	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()
	err = writeSynced(src+".gz", func(w io.Writer) error {
		z := gzip.NewWriter(w)
		if _, err := io.Copy(z, bufio.NewReaderSize(in, 1<<20)); err != nil {
			return err
		}
		return z.Close()
	})
	if err != nil {
		_ = os.Remove(src + ".gz")
		return 0, err
	}
	in.Close()
	_ = os.Remove(src)
	st, err := os.Stat(src + ".gz")
	if err != nil {
		return 0, err
	}
	return st.Size(), nil
}

// vhdChecksum is the VHD footer's and header's: the one's complement of the sum of the bytes but the checksum's own.
func vhdChecksum(b []byte, at int) uint32 {
	var sum uint32
	for i, c := range b {
		if i < at || i >= at+4 {
			sum += uint32(c)
		}
	}
	return ^sum
}

// VerifyVHD is the "Verify export on completion" of a disk (DiscUtils' Vhd FileChecker, in short): the footer and its
// checksum, and for a dynamic disk the copy of the footer, the header and the block table, every block inside the file.
func VerifyVHD(r io.ReaderAt, size int64) error {
	bad := func(why string) error { return fmt.Errorf("the disk is not a valid VHD: %s", why) }
	if size < sector {
		return bad("it is too short")
	}
	foot := make([]byte, sector)
	if _, err := r.ReadAt(foot, size-sector); err != nil {
		return err
	}
	be := binary.BigEndian
	if string(foot[:8]) != "conectix" || be.Uint32(foot[64:]) != vhdChecksum(foot, 64) {
		return bad("its footer is damaged")
	}
	virtual := int64(be.Uint64(foot[48:]))
	switch be.Uint32(foot[60:]) {
	case 2: // fixed: the disk, then the footer
		if size-sector != virtual {
			return bad("its size does not match its footer")
		}
		return nil
	case 3:
	default:
		return bad("it is neither a fixed nor a dynamic disk")
	}
	first := make([]byte, sector)
	if _, err := r.ReadAt(first, 0); err != nil {
		return err
	}
	if !bytes.Equal(first, foot) {
		return bad("the copy of its footer does not match")
	}
	headAt := int64(be.Uint64(foot[16:]))
	head := make([]byte, 1024)
	if headAt < 0 || headAt+1024 > size-sector {
		return bad("its header is outside the file")
	}
	if _, err := r.ReadAt(head, headAt); err != nil {
		return err
	}
	if string(head[:8]) != "cxsparse" || be.Uint32(head[36:]) != vhdChecksum(head, 36) {
		return bad("its header is damaged")
	}
	tableAt, entries, block := int64(be.Uint64(head[16:])), int64(be.Uint32(head[28:])), int64(be.Uint32(head[32:]))
	if block <= 0 || block%sector != 0 || entries*block < virtual || tableAt < 0 || tableAt+entries*4 > size-sector {
		return bad("its block table does not match its size")
	}
	bitmap := (block/sector/8 + sector - 1) / sector * sector
	table := make([]byte, entries*4)
	if _, err := r.ReadAt(table, tableAt); err != nil {
		return err
	}
	for i := int64(0); i < entries; i++ {
		if at := be.Uint32(table[i*4:]); at != 0xffffffff && int64(at)*sector+bitmap+block > size-sector {
			return bad("a block is outside the file")
		}
	}
	return nil
}

// VerifyVHDError is EXPORTING_VDI_VERIFICATION_FAILURE for the disk file, with the reason (the C# app logs it).
func VerifyVHDError(file string, err error) error {
	return errors.New(messages.Format(messages.Get("EXPORTING_VDI_VERIFICATION_FAILURE"), file) + " (" + err.Error() + ")")
}

// maxEulas is ExportEulaPage.MAX_EULA_DOCS.
const maxEulas = 25

// ReadEulas is OVF.AddEula's read of each EULA file (ExportEulaPage: at most 25, none empty, CA-55247).
func ReadEulas(files []string) ([]string, error) {
	if len(files) > maxEulas {
		return nil, fmt.Errorf("at most %d EULAs can be included in the appliance", maxEulas)
	}
	out := []string{}
	for _, f := range files {
		if st, err := os.Stat(f); err != nil || !st.Mode().IsRegular() { // a FIFO would block the open
			return nil, errors.New(messages.Format(messages.Get("EXPORT_EULA_PAGE_INVALID_FILE"), f))
		}
		b, err := readLimited(f, ovfMaxManifest)
		if err != nil {
			return nil, fmt.Errorf("Export Halted: Cannot read EULA: %w", err)
		}
		text := eulaText(b)
		if strings.TrimSpace(text) == "" {
			return nil, errors.New(messages.Format(messages.Get("EXPORT_EULA_PAGE_INVALID_FILE"), f))
		}
		out = append(out, text)
	}
	return out, nil
}

// eulaText is StreamReader.ReadToEnd: UTF-16 by its byte order mark, else UTF-8 (its mark dropped, a bad byte as U+FFFD).
func eulaText(b []byte) string {
	var order binary.ByteOrder
	switch {
	case bytes.HasPrefix(b, []byte{0xff, 0xfe}):
		order = binary.LittleEndian
	case bytes.HasPrefix(b, []byte{0xfe, 0xff}):
		order = binary.BigEndian
	default:
		return strings.ToValidUTF8(string(bytes.TrimPrefix(b, []byte{0xef, 0xbb, 0xbf})), "�")
	}
	u := make([]uint16, (len(b)-2)/2)
	for i := range u {
		u[i] = order.Uint16(b[2+2*i:])
	}
	return string(utf16.Decode(u))
}

// WritePackage is the end of ExportApplianceAction.RunCore, once the disks are in dir: the descriptor "<name>.ovf"
// (OVF.SaveAs), the manifest when asked (ManifestAndSign: after CompressFiles compressed the disks one by one), then for
// an OVA the archive (ConvertOVFtoOVA, compressed as a whole when asked). It returns the package's file in dir.
func WritePackage(dir, name string, vms []ExportVM, eulas []string, manifest, ova, compress bool) (string, error) {
	ovf := name + ".ovf"
	desc := ApplianceDescriptor(name, vms, eulas)
	if len(desc) > ovfMaxDescriptor { // the import would refuse it (the EULAs make it big)
		return "", errors.New("the EULAs are too large: the OVF descriptor would be larger than 16 MB")
	}
	if err := writeSynced(filepath.Join(dir, ovf), func(w io.Writer) error { _, err := w.Write(desc); return err }); err != nil {
		return "", err
	}
	files := []string{ovf}
	for _, vm := range vms {
		for _, d := range vm.Disks {
			files = append(files, d.Href)
		}
	}
	if manifest {
		if err := WriteManifest(dir, name, files); err != nil {
			return "", err
		}
	}
	if !ova {
		return ovf, nil
	}
	if manifest { // the archive's order: the descriptor, the manifest, the disks
		files = slices.Insert(files, 1, name+".mf")
	}
	return WriteOVA(dir, name, files, compress)
}
