package app

import (
	"cmp"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// The General tab (XenAdmin.Desktop.Core/Details.cs General, Common, HostMemory, Certificates and
// ManagementInterfaces, after GeneralTabPage): the sections of a VM or template, a server, a pool, an SR or a network.
// shortcut: no Custom Fields section (Details.General has none; the Properties dialog edits them), no PVS-Accelerator
// status (PVS_proxy, a Citrix feature) and no LUN-per-VDI multipath check; add them when XCP-ng servers have them.

// Section is one of Details' property sections: a title, its rows (name, value; a value may hold several lines) and
// whether it starts expanded (Views.Sections).
type Section struct {
	Title    string      `json:"title"`
	Rows     [][2]string `json:"rows"`
	Expanded bool        `json:"expanded"`
}

func (s *Section) add(key, value string) { s.Rows = append(s.Rows, [2]string{key, value}) }

// fn is Details.FN: the FriendlyNames label of a property, else the property.
func fn(property string) string {
	if n := messages.Friendly("Label-" + property); n != "" {
		return n
	}
	return property
}

// atLeast is Helpers.ProductVersionCompare(v, min) >= 0 for a version that is known (the Helpers.*OrGreater checks).
func atLeast(v, min string) bool { return v != "" && rules.CompareVersions(v, min) >= 0 }

// xapiVersion is Host.GetXapiVersion.
func xapiVersion(h xenapi.HostRecord) string {
	if v := h.SoftwareVersion["xapi_build"]; v != "" {
		return v
	}
	return h.SoftwareVersion["xapi"]
}

// productVersion is "{ProductBrand} {ProductVersionText}" (product_version_text, else product_version).
func productVersion(h xenapi.HostRecord) string {
	v := h.SoftwareVersion["product_version_text"]
	if v == "" {
		v = h.SoftwareVersion["product_version"]
	}
	return h.SoftwareVersion["product_brand"] + " " + v
}

// sinceUnix is Host.Uptime and AgentUptime: the time from an other_config Unix time to now (the server's), "" when
// unset.
func sinceUnix(now time.Time, v string) string {
	t, err := strconv.ParseFloat(v, 64)
	if err != nil || t == 0 {
		return ""
	}
	return prettyTimeSpan(now.Sub(time.Unix(int64(t), 0)))
}

// common is Details.Common: name, description (shown unless a default template's), tags and folder.
func common(name, description string, showDescription bool, tags []string, other map[string]string) Section {
	s := Section{Title: "General", Expanded: true}
	s.add(fn("host.name_label"), name)
	if showDescription {
		s.add(fn("host.name_description"), description)
	}
	t := messages.Get("NONE")
	if len(tags) > 0 {
		sorted := append([]string{}, tags...)
		sort.Strings(sorted)
		t = strings.Join(sorted, ", ")
	}
	s.add(messages.Get("TAGS"), t)
	folder := other["folder"]
	if folder == "" {
		folder = messages.Get("NONE")
	}
	s.add(messages.Get("FOLDER"), folder)
	return s
}

// general is what every kind reads of the connection.
type general struct {
	*hostTabData
	pool    xenapi.PoolRecord // Helpers.GetPoolOfOne
	visible bool              // Helpers.GetPool != null (Pool.IsVisible)
	metrics map[xenapi.HostMetricsRef]xenapi.HostMetricsRecord
	now     time.Time // the server's clock (Connection.ServerTimeOffset), as Overview
}

func (g *general) live(h xenapi.HostRecord) bool { return g.metrics[h.Metrics].Live } // Host.IsLive

func (g *general) coordinator() xenapi.HostRecord { return g.hosts[g.pool.Master] }

// sortedHosts are the servers by name (OrderBy(h => h.Name())).
func (g *general) sortedHosts() []xenapi.HostRef {
	refs := make([]xenapi.HostRef, 0, len(g.hosts))
	for r := range g.hosts {
		refs = append(refs, r)
	}
	sort.Slice(refs, func(i, j int) bool { return g.hosts[refs[i]].NameLabel < g.hosts[refs[j]].NameLabel })
	return refs
}

// ReadGeneral reads the General tab of an object ("vms", "templates", "hosts", "pools", "srs", "networks"); the
// sections without rows are left out, as Views.Sections does.
func ReadGeneral(c *xapi.Connection, kind, uuid string) ([]Section, error) {
	s := c.Session
	d, err := readHostTabData(s)
	if err != nil {
		return nil, err
	}
	g := &general{hostTabData: d}
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	for _, p := range pools {
		g.pool = p
	}
	g.visible = len(pools) > 0 && (g.pool.NameLabel != "" || len(d.hosts) > 1)
	g.now = time.Now()
	if t, err := xenapi.Host.GetServertime(s, g.pool.Master); err == nil && !t.IsZero() {
		g.now = t
	}
	if g.metrics, err = xenapi.HostMetrics.GetAllRecords(s); err != nil {
		return nil, xapi.Sanitize(err)
	}
	var list []Section
	switch kind {
	case "hosts":
		list, err = g.host(uuid)
	case "vms", "templates":
		list, err = g.vm(uuid)
	case "pools":
		list, err = g.poolSections(uuid)
	case "srs":
		list, err = g.sr(uuid)
	case "networks":
		list, err = g.network(uuid)
	default:
		err = fmt.Errorf("no General tab for %s", kind)
	}
	if err != nil {
		return nil, err
	}
	out := []Section{}
	for _, sec := range list {
		if len(sec.Rows) > 0 {
			out = append(out, sec)
		}
	}
	return out, nil
}

func (g *general) host(uuid string) ([]Section, error) {
	s := g.s
	ref, err := xenapi.Host.GetByUUID(s, uuid)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	h, ok := g.hosts[ref]
	if !ok {
		return nil, fmt.Errorf("the server %s is gone", uuid)
	}
	desc := h.NameDescription // Host.Description: the installer's description with the server's brand
	if desc == "Default install of XenServer" || desc == "Default install" {
		brand, ok := h.SoftwareVersion["product_brand"]
		if !ok {
			brand = "XCP-ng"
		}
		desc = messages.Format(messages.Get("DEFAULT_INSTALL_OF_XENSERVER"), brand)
	}
	gen := common(h.NameLabel, desc, true, h.Tags, h.OtherConfig)
	if g.visible {
		gen.add(messages.Get("POOL_COORDINATOR"), rules.YesNo(ref == g.pool.Master))
	}
	gen.add(fn("host.enabled"), rules.HostEnabledText(g.live(h), h.Enabled, h.OtherConfig["MAINTENANCE_MODE"] == "true"))
	gen.add("Autoboot of VMs enabled", rules.YesNo(g.pool.OtherConfig["auto_poweron"] == "true"))       // Host.GetVmAutostartEnabled
	if atLeast(h.SoftwareVersion["platform_version"], "3.2.50") && atLeast(xapiVersion(h), "1.290.0") { // CloudOrGreater
		on := g.pool.TLSVerificationEnabled && (!atLeast(xapiVersion(h), "1.313.0") || h.TLSVerificationEnabled)
		gen.add(messages.Get("CERTIFICATE_VERIFICATION_KEY"), map[bool]string{true: messages.Get("ENABLED"), false: messages.Get("DISABLED")}[on])
	}
	gen.add(fn("host.iscsi_iqn"), h.IscsiIqn) // Host.GetIscsiIqn from Kolkata (every XCP-ng 8.x)
	gen.add(fn("host.log_destination"), rules.LogDestinationText(h.Logging["syslog_destination"]))
	gen.add(fn("host.uptime"), sinceUnix(g.now, h.OtherConfig["boot_time"]))
	gen.add(fn("host.agentUptime"), sinceUnix(g.now, h.OtherConfig["agent_start_time"]))
	gen.add(fn("host.uuid"), h.UUID)

	version := Section{Title: "Version Details"}
	if len(h.SoftwareVersion) > 0 {
		if date, ok := h.SoftwareVersion["date"]; ok {
			version.add(messages.Get("SOFTWARE_VERSION_DATE"), date)
		}
		version.add(messages.Get("SOFTWARE_VERSION_PRODUCT_VERSION"), productVersion(h))
		if build, ok := h.SoftwareVersion["build_number"]; ok {
			version.add("Build number", build)
		}
	}
	memory := Section{Title: "Memory"}
	if t, err := g.memory(ref); err == nil {
		for _, r := range t.Rows {
			memory.add(r[0], r[1])
		}
	}
	cpus := Section{Title: "CPUs"}
	var all []xenapi.HostCPURecord
	for _, r := range h.HostCPUs {
		if cpu, err := xenapi.HostCPU.GetRecord(s, r); err == nil {
			all = append(all, cpu)
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Number < all[j].Number })
	if len(all) > 0 {
		n := "0"
		if len(all) > 1 {
			n = fmt.Sprintf("0 - %d", len(all)-1)
		}
		cpus.add(messages.Format(messages.Get("GENERAL_DETAILS_CPU_NUMBER"), n), rules.CPUText(all[0].Vendor, all[0].Modelname, all[0].Speed))
	}
	return []Section{gen, g.certificates(h), version, g.managementInterfaces([]xenapi.HostRef{ref}), memory, cpus}, nil
}

// certificates is Details.Certificates: the server's certificates by type, from Stockholm.
func (g *general) certificates(h xenapi.HostRecord) Section {
	sec := Section{Title: "Certificates"}
	if !atLeast(h.SoftwareVersion["platform_version"], "3.1.50") {
		return sec
	}
	var certs []xenapi.CertificateRecord
	for _, r := range h.Certificates {
		if c, err := xenapi.Certificate.GetRecord(g.s, r); err == nil {
			certs = append(certs, c)
		}
	}
	sort.SliceStable(certs, func(i, j int) bool { return certs[i].Type < certs[j].Type }) // ca, host, host_internal
	for _, c := range certs {
		fp := c.Fingerprint
		if fp == "" { // "fingerprint" is deprecated for fingerprint_sha256
			fp = c.FingerprintSha256
		}
		sec.add(rules.CertificateTypeText(string(c.Type)), rules.CertificateText(c.NotBefore, c.NotAfter, fp))
	}
	return sec
}

// managementInterfaces is Details.ManagementInterfaces: each server's DNS hostname and management PIFs' addresses
// (PIF.FriendlyIPAddress).
func (g *general) managementInterfaces(hosts []xenapi.HostRef) Section {
	sec := Section{Title: "Management Interfaces"}
	for _, r := range hosts {
		h := g.hosts[r]
		sec.add(fn("host.hostname"), h.Hostname)
		for _, pr := range h.PIFs {
			p, ok := g.pifs[pr]
			if !ok || !p.Management {
				continue
			}
			ip := p.IP
			if ip == "" {
				ip = messages.Get("PIF_UNKNOWN")
				if p.IPConfigurationMode == xenapi.IPConfigurationModeDHCP {
					ip = messages.Get("PIF_DHCP")
				}
			}
			sec.add(messages.Get("MANAGEMENT_INTERFACE"), ip)
		}
	}
	return sec
}

func (g *general) poolSections(uuid string) ([]Section, error) {
	p := g.pool
	if p.UUID != uuid {
		return nil, fmt.Errorf("the pool %s is gone", uuid)
	}
	coord, hasCoord := g.hosts[p.Master]
	name := p.NameLabel // Pool.Name
	if name == "" {
		name = coord.NameLabel
	}
	gen := common(name, p.NameDescription, true, p.Tags, p.OtherConfig)
	sockets := 0 // Pool.CpuSockets: the servers' socket_count (0 when unknown)
	for _, h := range g.hosts {
		n, _ := strconv.Atoi(h.CPUInfo["socket_count"])
		sockets += n
	}
	gen.add(messages.Get("NUMBER_OF_SOCKETS"), strconv.Itoa(sockets))
	gen.add("Autoboot of VMs enabled", rules.YesNo(p.OtherConfig["auto_poweron"] == "true"))
	// CloudOrGreater and XapiEqualOrGreater_1_290_0 of the connection: true without a coordinator
	if !hasCoord || atLeast(coord.SoftwareVersion["platform_version"], "3.2.50") && atLeast(xapiVersion(coord), "1.290.0") {
		gen.add(messages.Get("CERTIFICATE_VERIFICATION_KEY"), map[bool]string{true: messages.Get("ENABLED"), false: messages.Get("DISABLED")}[p.TLSVerificationEnabled])
	}
	if hasCoord {
		gen.add(messages.Get("SOFTWARE_VERSION_PRODUCT_VERSION"), productVersion(coord))
	}
	gen.add(fn("host.uuid"), p.UUID)
	return []Section{gen, g.managementInterfaces(g.sortedHosts())}, nil
}

func (g *general) vm(uuid string) ([]Section, error) {
	s := g.s
	ref, err := xenapi.VM.GetByUUID(s, uuid)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	m, merr := xenapi.VMMetrics.GetRecord(s, vm.Metrics)
	gm, gerr := xenapi.VMGuestMetrics.GetRecord(s, vm.GuestMetrics)
	def := defaultTemplate(vm)
	desc, showDesc := rules.VMDescription(vm.IsATemplate, vm.IsASnapshot || vm.IsSnapshotFromVmpp, def, vm.NameLabel, vm.NameDescription, vm.OtherConfig)
	gen := common(vm.NameLabel, desc, showDesc, vm.Tags, vm.OtherConfig)
	hvm := vm.HVMBootPolicy != "" // VM.IsHVM
	gen.add(fn("VM.OSName"), osName(gm))
	gen.add(fn("VM.OperatingMode"), map[bool]string{true: messages.Get("VM_OPERATING_MODE_HVM"), false: messages.Get("VM_OPERATING_MODE_PV")}[hvm])
	if !def {
		gen.add(messages.Get("BIOS_STRINGS_COPIED"), rules.YesNo(rules.BiosStringsCopied(def, vm.BiosStrings)))
	}
	windows := isWindows(s, vm)
	running := vm.PowerState == xenapi.VMPowerStateRunning
	if !vm.IsATemplate {
		state := rules.VirtualizationText(rules.GuestTools{Windows: windows, Running: running, Metrics: merr == nil, GuestMetrics: gerr == nil,
			IODrivers: gm.PVDriversDetected, Management: gm.Other["feature-static-ip-setting"],
			Early:     merr == nil && !m.StartTime.IsZero() && g.now.Sub(m.StartTime) < 2*time.Minute, // GetBodgeStartupTime
			PVVersion: gm.PVDriversVersion, UpToDate: gm.PVDriversUpToDate, VendorDevice: vm.HasVendorDevice})
		if state != "" {
			gen.add(fn("VM.VirtualizationState"), state)
		}
	}
	up := vm.PowerState == xenapi.VMPowerStateRunning || vm.PowerState == xenapi.VMPowerStatePaused || vm.PowerState == xenapi.VMPowerStateSuspended
	if up && merr == nil && m.StartTime.Unix() > 0 { // VM.RunningTime: not at the epoch nor unset
		gen.add(fn("VM.uptime"), prettyTimeSpan(g.now.Sub(m.StartTime)))
	}
	vbds := make([]xenapi.VBDRecord, 0, len(vm.VBDs))
	for _, r := range vm.VBDs {
		if b, err := xenapi.VBD.GetRecord(s, r); err == nil {
			vbds = append(vbds, b)
		}
	}
	// GeneralTabPage.VMCanChooseHomeServer: not under WLB, and the disks on shared storage (VM.IsOnSharedStorage) or no
	// disk nor local CD (VM.HasNoDisksAndNoLocalCD), which is the same set of VMs: none of its VDIs on a local SR
	if !vm.IsATemplate && !(g.pool.WlbEnabled && g.pool.WlbURL != "") && !onLocalSR(s, vbds) {
		home := messages.Get("NONE") // VM.AffinityServerString
		if h, ok := g.hosts[vm.Affinity]; ok {
			home = h.NameLabel
		}
		gen.add(fn("VM.affinity"), home)
	}
	gen.add(fn("host.uuid"), vm.UUID)

	boot := Section{Title: "Boot Options"}
	boot.add(fn("VM.auto_boot"), rules.YesNo(vm.OtherConfig["auto_poweron"] == "true"))
	if hvm {
		order, set := vm.HVMBootParams["order"]
		boot.add(fn("VM.BootOrder"), rules.BootOrderText(order, set))
		if atLeast(g.coordinator().SoftwareVersion["platform_version"], "2.9.50") { // NaplesOrGreater
			boot.add(fn("VM.BootMode"), rules.BootModeText(vm.HVMBootParams["firmware"], vm.Platform["secureboot"], g.pool.UefiCertificates != ""))
		}
	} else {
		boot.add(fn("VM.PV_args"), vm.PVArgs)
	}

	cpus := Section{Title: "CPUs"}
	cpus.add(fn("VM.VCPUs"), strconv.Itoa(vm.VCPUsAtStartup))
	if vm.VCPUsAtStartup != vm.VCPUsMax || !windows && !ops.VCPUHotplugRestricted(g.hosts) { // VM.SupportsVcpuHotplug
		cpus.add(fn("VM.MaxVCPUs"), strconv.Itoa(vm.VCPUsMax))
	}
	cpus.add(fn("VM.Topology"), rules.VMTopologyText(int64(vm.VCPUsMax), vm.Platform["cores-per-socket"]))
	out := []Section{gen, boot, cpus}
	if running {
		out = append(out, g.readCaching(vm, vbds))
	}
	return out, nil
}

// onLocalSR is !VM.IsOnSharedStorage: one of the VM's disks or CDs is on an SR that is not shared.
func onLocalSR(s *xenapi.Session, vbds []xenapi.VBDRecord) bool {
	for _, b := range vbds {
		sr, err := xenapi.VDI.GetSR(s, b.VDI)
		if err != nil {
			continue
		}
		if shared, err := xenapi.SR.GetShared(s, sr); err == nil && !shared {
			return true
		}
	}
	return false
}

// readCaching is the Read Caching section of a running VM (VM.ReadCachingVDIs and ReadCachingDisabledReason): the
// attached disks cached on its server, else why none is.
func (g *general) readCaching(vm xenapi.VMRecord, vbds []xenapi.VBDRecord) Section {
	sec := Section{Title: "Read Caching"}
	resident, ok := g.hosts[vm.ResidentOn]
	var cached, reasons []string
	where := ""
	for _, b := range vbds {
		if !b.CurrentlyAttached || !ok {
			continue
		}
		vdi, err := xenapi.VDI.GetRecord(g.s, b.VDI)
		if err != nil {
			continue
		}
		if vdi.SmConfig["read-caching-enabled-on-"+resident.UUID] == "true" { // VDI.ReadCachingEnabled
			if where == "" {
				where = location(g.s)
			}
			// VDI.NameWithLocation: "disk on 'SR' in 'pool'"; a local SR of a standalone server has no location
			// (SR.LocationString through Host.LocationString)
			name := strings.TrimSpace(vdi.NameLabel + " " + where)
			if sr, err := xenapi.SR.GetRecord(g.s, vdi.SR); err == nil {
				at := where
				if !g.visible && !sr.Shared {
					at = ""
				}
				name = strings.TrimSpace(messages.Format(messages.Get("VDI_ON_SR_TITLE"), vdi.NameLabel, srI18N(g.s, g.hosts, sr, "name_label", sr.NameLabel), at))
			}
			cached = append(cached, name)
		} else {
			reasons = append(reasons, vdi.SmConfig["read-caching-reason-"+resident.UUID])
		}
	}
	if len(cached) > 0 {
		sec.add(fn("VM.read_caching_status"), messages.Get("VM_READ_CACHING_ENABLED"))
		sec.add(fn("VM.read_caching_disks"), strings.Join(cached, "\n"))
		return sec
	}
	sec.add(fn("VM.read_caching_status"), messages.Get("VM_READ_CACHING_DISABLED"))
	if r := rules.ReadCachingReason(reasons, ops.ReadCachingRestricted(g.hosts)); r != "" {
		sec.add(fn("VM.read_caching_reason"), r)
	}
	return sec
}

func (g *general) sr(uuid string) ([]Section, error) {
	s := g.s
	ref, err := xenapi.SR.GetByUUID(s, uuid)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	sr, err := xenapi.SR.GetRecord(s, ref)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	gen := common(srI18N(s, g.hosts, sr, "name_label", sr.NameLabel), srI18N(s, g.hosts, sr, "name_description", sr.NameDescription), true, sr.Tags, sr.OtherConfig)
	sms, _ := xenapi.SM.GetAllRecords(s)
	typeName := srTypeName(sr.Type) // SR.FriendlyTypeName: an unknown type is its SMAPIv3 plugin's name, else Unknown
	if messages.Friendly("Label-SR.SRTypes-"+sr.Type) == "" {
		typeName = messages.Friendly("Label-SR.SRTypes-unknown")
		for _, sm := range sms {
			if sm.Type == sr.Type && rules.CompareVersions(sm.RequiredAPIVersion, "3.0") >= 0 {
				typeName = cmp.Or(sm.NameLabel, sr.Type)
			}
		}
	}
	gen.add(messages.Get("TYPE"), typeName)
	if sr.ContentType != "iso" && sr.Type != "udev" { // SR.SizeString
		gen.add(fn("SR.size"), messages.Format(messages.Get("SR_SIZE_USED"), diskSize(sr.PhysicalUtilisation), diskSize(sr.PhysicalSize), diskSize(sr.VirtualAllocation)))
	}
	var pbds []xenapi.PBDRecord
	for _, r := range sr.PBDs {
		if p, err := xenapi.PBD.GetRecord(s, r); err == nil {
			pbds = append(pbds, p)
		}
	}
	gen.add(fn("SR.scsiid"), scsiID(sr, pbds))
	gen.add(fn("host.uuid"), sr.UUID)

	capable := sr.SmConfig["multipathable"] == "true" // SR.MultipathCapable: sm_config, or the SM's SR_MULTIPATH feature
	for _, sm := range sms {
		if _, ok := sm.Features["SR_MULTIPATH"]; ok && sm.Type == sr.Type {
			capable = true
		}
	}
	detached, multipathOK := true, true
	for _, p := range pbds {
		detached = detached && !p.CurrentlyAttached                           // SR.IsDetached
		if capable && p.OtherConfig["multipathed"] == "true" && !pathsUp(p) { // SR.MultipathAOK for a LUN per SR
			multipathOK = false
		}
	}
	status := Section{Title: "Status"}
	status.add(fn("SR.state"), rules.SRStateText(detached, broken(s, sr, len(g.hosts)), multipathOK))
	for _, p := range pbds {
		h, known := g.hosts[p.Host]
		status.add(h.NameLabel, rules.PBDStatusText(p.CurrentlyAttached, known, g.live(h)))
	}
	multipath := Section{Title: "Multipathing"}
	if !capable {
		multipath.add(messages.Get("MULTIPATH_CAPABLE"), messages.Get("NO"))
	} else {
		for _, r := range g.sortedHosts() {
			active := false // SR.GetPBDFor(h).MultipathActive
			for _, p := range pbds {
				if p.Host == r && p.OtherConfig["multipathed"] == "true" {
					active = true
				}
			}
			multipath.add(g.hosts[r].NameLabel, map[bool]string{true: messages.Get("MULTIPATH_ACTIVE"), false: messages.Get("MULTIPATH_NOT_ACTIVE")}[active])
		}
	}
	return []Section{gen, status, multipath}, nil
}

// pathsUp is SR.GetMultiPathStatusLunPerSR + CheckMultipathString for a PBD: its first "mpath..." status has every
// path up (or cannot be read).
func pathsUp(p xenapi.PBDRecord) bool {
	keys := make([]string, 0, len(p.OtherConfig))
	for k := range p.OtherConfig {
		if strings.HasPrefix(k, "mpath") {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return len(keys) == 0 || !rules.MultipathDown(p.OtherConfig[keys[0]])
}

// scsiID is SR.GetScsiID: a PBD's SCSIid, else the sm_config devserial without "scsi-" and trailing commas.
func scsiID(sr xenapi.SRRecord, pbds []xenapi.PBDRecord) string {
	for _, p := range pbds {
		if id, ok := p.DeviceConfig["SCSIid"]; ok {
			return id
		}
	}
	id, ok := sr.SmConfig["devserial"]
	if !ok {
		return ""
	}
	return strings.TrimRight(strings.TrimPrefix(id, "scsi-"), ",")
}

func (g *general) network(uuid string) ([]Section, error) {
	ref, err := xenapi.Network.GetByUUID(g.s, uuid)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	n, err := xenapi.Network.GetRecord(g.s, ref)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	return []Section{common(g.networkName(n), n.NameDescription, true, n.Tags, n.OtherConfig)}, nil
}
