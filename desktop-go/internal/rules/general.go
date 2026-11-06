package rules

import (
	"strconv"
	"strings"
	"time"

	"xcpng-center/internal/messages"
)

// The General tab's texts (XenAdmin.Desktop.Core/Details.cs General and the XenModel helpers it calls).

// vmTools is BrandManager.VmTools of the XCP-ng branding.
const vmTools = "XCP-ng VM Tools"

// YesNo is Helpers.BoolToString.
func YesNo(b bool) string {
	if b {
		return messages.Get("YES")
	}
	return messages.Get("NO")
}

// HostEnabledText is the server's Enabled row: not reachable, in maintenance mode, disabled or Yes.
func HostEnabledText(live, enabled, maintenance bool) string {
	switch {
	case !live:
		return messages.Get("HOST_NOT_LIVE")
	case !enabled && maintenance:
		return messages.Get("HOST_IN_MAINTENANCE_MODE")
	case !enabled:
		return messages.Get("DISABLED")
	}
	return messages.Get("YES")
}

// SRStateText is the SR's State row: Detached, Broken, Multipath Failure or OK.
func SRStateText(detached, broken, multipathOK bool) string {
	switch {
	case detached:
		return messages.Get("DETACHED")
	case broken:
		return messages.Get("GENERAL_SR_STATE_BROKEN")
	case !multipathOK:
		return messages.Get("GENERAL_MULTIPATH_FAILURE")
	}
	return messages.Get("GENERAL_STATE_OK")
}

// PBDStatusText is PBD.StatusString: unplugged (or no server), server not reachable, or connected.
func PBDStatusText(attached, hostKnown, hostLive bool) string {
	switch {
	case !attached || !hostKnown:
		return messages.Get("REPAIR_SR_DIALOG_UNPLUGGED")
	case !hostLive:
		return messages.Get("HOST_NOT_LIVE")
	}
	return messages.Get("CONNECTED")
}

// BootOrderText is the Boot order row: VM.GetBootOrder ("CD" when unset) with the missing devices after it
// (Union with D, C, N), one per line.
func BootOrderText(order string, set bool) string {
	if !set {
		order = "CD"
	}
	var lines []string
	seen := map[rune]bool{}
	for _, c := range strings.ToUpper(order) + "DCN" {
		if seen[c] {
			continue
		}
		seen[c] = true
		switch c {
		case 'C':
			lines = append(lines, messages.Get("BOOT_HARD_DISK"))
		case 'D':
			lines = append(lines, messages.Get("BOOT_DVD_DRIVE"))
		default:
			lines = append(lines, messages.Get("NETWORK"))
		}
	}
	return strings.Join(lines, "\n")
}

// BootModeText is the Boot mode row: BIOS, UEFI, or UEFI Secure Boot (with the missing certificates when the pool has
// none); firmware is HVM_boot_params "firmware", secureBoot the platform's "secureboot".
func BootModeText(firmware, secureBoot string, hasCerts bool) string {
	if strings.ToLower(strings.TrimSpace(firmware)) != "uefi" { // VM.IsDefaultBootModeUefi
		return messages.Get("BIOS_BOOT")
	}
	switch secure := strings.ToLower(strings.TrimSpace(secureBoot)); {
	case secure == "true" && !hasCerts:
		return messages.Format(messages.Get("CUSTOM_FIELD_NAME_AND_TYPE"), messages.Get("UEFI_SECURE_BOOT"), messages.Get("GUEFI_SECUREBOOT_MODE_MISSING_CERTIFICATES"))
	case secure == "true" || secure == "auto" && hasCerts:
		return messages.Get("UEFI_SECURE_BOOT")
	}
	return messages.Get("UEFI_BOOT")
}

// VMTopologyText is VM.Topology: the vCPUs' sockets and cores per socket (platform "cores-per-socket", 1 when unset or
// not a number), "invalid" when the cores do not divide the vCPUs or need more than MaxSockets.
func VMTopologyText(vcpusMax int64, coresPerSocket string) string {
	cores := int64(1) // VM.GetCoresPerSocket
	if n, err := strconv.ParseInt(coresPerSocket, 10, 64); err == nil {
		cores = n
	}
	var sockets int64 // VM.ValidVCPUConfiguration; 0 cores (a division by zero in C#) is invalid too
	if cores > 0 && vcpusMax%cores == 0 && vcpusMax/cores <= MaxSockets {
		sockets = vcpusMax / cores
	}
	return TopologyText(sockets, cores)
}

// BiosStringsCopied is VM.BiosStringsCopied: a VM's BIOS strings are not Xen's own.
func BiosStringsCopied(defaultTemplate bool, bios map[string]string) bool {
	if defaultTemplate || len(bios) == 0 {
		return false
	}
	return bios["bios-vendor"] != "Xen" || bios["system-manufacturer"] != "Xen"
}

// templateOS are VM.TemplateType's names of the default templates (each has its own description text in the C# app,
// so the General tab shows none).
var templateOS = []string{"citrix", "debian", "gooroom", "rocky", "centos", "linx", "red hat", "oracle", "suse", "scientific",
	"windows", "ubuntu", "yinhe", "kylin", "asianux", "turbo", "solaris", "coreos"}

// VMDescription is VM.DescriptionType and VM.Description: shown is false for a default template of a known OS; a
// default template of another OS shows its FriendlyNames description; a P2V import's description is empty.
func VMDescription(template, snapshot, defaultTemplate bool, name, description string, other map[string]string) (string, bool) {
	if template && !snapshot && defaultTemplate {
		os := strings.ToLower(name)
		for _, t := range templateOS {
			if strings.Contains(os, t) {
				return "", false
			}
		}
		if d := messages.Friendly("VM.TemplateDescription-" + name); d != "" {
			return d, true
		}
		return description, true
	}
	_, src := other["p2v_source_machine"] // VM.IsP2V
	_, date := other["p2v_import_date"]
	if src && date && strings.HasPrefix(description, "VM imported from physical machine") {
		return "", true
	}
	return description, true
}

// GuestTools is what the Virtualization state row reads of a VM, its metrics and guest metrics.
type GuestTools struct {
	Windows      bool // VM.HasNewVirtualizationStates
	Running      bool // power state Running
	Metrics      bool // the VM_metrics resolve
	GuestMetrics bool // the VM_guest_metrics resolve
	IODrivers    bool // PV_drivers_detected
	Management   string
	Early        bool              // started less than two minutes ago (GetBodgeStartupTime)
	PVVersion    map[string]string // PV_drivers_version
	UpToDate     bool              // PV_drivers_up_to_date
	VendorDevice bool              // has_vendor_device
}

// VirtualizationText is GeneralTabPage.GenerateVirtualisationStatusForGeneralBox over VM.GetVirtualizationStatus:
// a Windows guest's I/O drivers, management agent and Windows Update lines, another running guest's tools state; ""
// when the row is not shown (a guest other than Windows that is not running).
func VirtualizationText(t GuestTools) string {
	text, unknown, io, mgmt := virtualizationStatus(t)
	if t.Windows {
		var lines []string
		if t.Running {
			if unknown {
				lines = append(lines, text)
			} else {
				lines = append(lines, pick(io, "VIRTUALIZATION_STATE_VM_IO_OPTIMIZED", "VIRTUALIZATION_STATE_VM_IO_NOT_OPTIMIZED"),
					pick(mgmt, "VIRTUALIZATION_STATE_VM_MANAGEMENT_AGENT_INSTALLED", "VIRTUALIZATION_STATE_VM_MANAGEMENT_AGENT_NOT_INSTALLED"))
			}
		}
		lines = append(lines, pick(t.VendorDevice, "VIRTUALIZATION_STATE_VM_RECEIVING_UPDATES", "VIRTUALIZATION_STATE_VM_NOT_RECEIVING_UPDATES"))
		return strings.Join(lines, "\n")
	}
	if t.Running {
		return text
	}
	return ""
}

func pick(b bool, yes, no string) string {
	if b {
		return messages.Get(yes)
	}
	return messages.Get(no)
}

// virtualizationStatus is VM.GetVirtualizationStatus: its friendly text and flags (Unknown, IoDriversInstalled,
// ManagementInstalled).
func virtualizationStatus(t GuestTools) (text string, unknown, io, mgmt bool) {
	text = messages.Get("VIRTUALIZATION_UNKNOWN")
	if !t.Running || !t.Metrics {
		return text, true, false, false
	}
	notInstalled := messages.Format(messages.Get("PV_DRIVERS_NOT_INSTALLED"), vmTools)
	if t.Windows {
		io = t.GuestMetrics && t.IODrivers
		n, _ := strconv.Atoi(t.Management) // IntKey(other, "feature-static-ip-setting", 0)
		mgmt = t.GuestMetrics && n != 0
		switch {
		case io && mgmt:
			text = messages.Get("VIRTUALIZATION_STATE_VM_IO_DRIVERS_AND_MANAGEMENT_AGENT_INSTALLED")
		case io:
			text = messages.Get("VIRTUALIZATION_STATE_VM_MANAGEMENT_AGENT_NOT_INSTALLED")
		case t.Early:
			return text, true, false, false
		default:
			text = notInstalled
		}
		return text, false, io, mgmt
	}
	major, hasMajor := t.PVVersion["major"]
	minor, hasMinor := t.PVVersion["minor"]
	switch {
	case (!t.GuestMetrics || !hasMajor || !hasMinor) && t.Early:
		return text, true, false, false
	case !t.GuestMetrics || !hasMajor || !hasMinor:
		return notInstalled, false, false, false
	case !t.UpToDate:
		return messages.Format(messages.Get("VIRTUALIZATION_OUT_OF_DATE"), major+"."+minor), false, false, false
	}
	return messages.Format(messages.Get("VIRTUALIZATION_OPTIMIZED"), major+"."+minor), false, true, true
}

// CertificateTypeText is the Certificates section's row name of a certificate type.
func CertificateTypeText(kind string) string {
	switch kind {
	case "ca":
		return messages.Get("CERTIFICATE_TYPE_CA")
	case "host":
		return messages.Get("CERTIFICATE_TYPE_HOST")
	}
	return messages.Get("CERTIFICATE_TYPE_HOST_INTERNAL")
}

// CertificateText is a certificate's row: its validity in local time (DATEFORMAT_DMY_HM, "MMM d, yyyy h:mm tt") and
// its thumbprint.
func CertificateText(notBefore, notAfter time.Time, fingerprint string) string {
	const dmyhm = "Jan 2, 2006 3:04 PM"
	return messages.Format(messages.Get("CERTIFICATE_VALIDITY_PERIOD_VALUE"), notBefore.Local().Format(dmyhm), notAfter.Local().Format(dmyhm)) +
		"\n" + messages.Format(messages.Get("CERTIFICATE_THUMBPRINT_VALUE"), fingerprint)
}

// CPUText is Helpers.GetCPUProperties: the vendor, model and speed of a server's CPU, one per line.
func CPUText(vendor, model string, mhz int) string {
	return messages.Format(messages.Get("GENERAL_CPU_VENDOR"), vendor) + "\n" + messages.Format(messages.Get("GENERAL_CPU_MODEL"), model) +
		"\n" + messages.Format(messages.Get("GENERAL_CPU_SPEED"), mhz)
}

// LogDestinationText is the Log destination row of a server's syslog destination ("" for local only).
func LogDestinationText(syslog string) string {
	if syslog == "" {
		return messages.Get("HOST_LOG_DESTINATION_LOCAL")
	}
	return messages.Format(messages.Get("HOST_LOG_DESTINATION_LOCAL_AND_REMOTE"), syslog)
}

// readCachingReasons are VDI.ReadCachingDisabledReasonCode in order (UNKNOWN first).
var readCachingReasons = map[string]int{"LICENSE_RESTRICTION": 1, "SR_NOT_SUPPORTED": 2, "NO_RO_IMAGE": 3, "SR_OVERRIDE": 4}

// ReadCachingReason is VM.ReadCachingDisabledReason over its disks' "read-caching-reason-<host>" codes: the last
// reason of the list (the disk nearest to being cached), "" when none is known; forbidden is the licence's
// restriction now.
func ReadCachingReason(codes []string, forbidden bool) string {
	top := 0
	for _, c := range codes {
		top = max(top, readCachingReasons[c])
	}
	switch top {
	case 1:
		if forbidden {
			return messages.Get("VM_READ_CACHING_DISABLED_REASON_LICENSE")
		}
		return messages.Get("VM_READ_CACHING_DISABLED_REASON_PREV_LICENSE")
	case 2:
		return messages.Get("VM_READ_CACHING_DISABLED_REASON_SR_TYPE")
	case 3:
		return messages.Get("VM_READ_CACHING_DISABLED_REASON_NO_RO_IMAGE")
	case 4:
		return messages.Get("VM_READ_CACHING_DISABLED_REASON_TURNED_OFF")
	}
	return ""
}
