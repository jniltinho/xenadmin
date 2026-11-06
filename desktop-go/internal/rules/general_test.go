package rules

import (
	"testing"
	"time"
)

func TestGeneralTexts(t *testing.T) {
	if HostEnabledText(false, true, false) != "Server is not reachable" || HostEnabledText(true, false, true) != "In maintenance mode" ||
		HostEnabledText(true, false, false) != "Disabled" || HostEnabledText(true, true, true) != "Yes" {
		t.Error("host enabled")
	}
	if SRStateText(true, true, false) != "Detached" || SRStateText(false, true, false) != "Broken" ||
		SRStateText(false, false, false) != "Multipath Failure" || SRStateText(false, false, true) != "OK" {
		t.Error("SR state")
	}
	if PBDStatusText(false, true, true) != "Unplugged" || PBDStatusText(true, false, true) != "Unplugged" ||
		PBDStatusText(true, true, false) != "Server is not reachable" || PBDStatusText(true, true, true) != "Connected" {
		t.Error("PBD status")
	}
	if got := BootOrderText("", false); got != "Hard Disk\nDVD-Drive\nNetwork" {
		t.Errorf("default boot order: %q", got)
	}
	if got := BootOrderText("nd", true); got != "Network\nDVD-Drive\nHard Disk" {
		t.Errorf("boot order: %q", got)
	}
	if BootModeText("bios", "true", false) != "BIOS Boot" || BootModeText(" UEFI ", "false", true) != "UEFI Boot" ||
		BootModeText("uefi", "true", true) != "UEFI Secure Boot" || BootModeText("uefi", "auto", true) != "UEFI Secure Boot" ||
		BootModeText("uefi", "auto", false) != "UEFI Boot" ||
		BootModeText("uefi", "true", false) != "UEFI Secure Boot (UEFI certificates are missing from your system)" {
		t.Error("boot mode")
	}
	if VMTopologyText(4, "2") != "2 sockets with 2 cores per socket" || VMTopologyText(2, "") != "2 sockets with 1 core per socket" ||
		VMTopologyText(3, "2") != TopologyText(0, 2) || VMTopologyText(2, "0") != TopologyText(0, 0) {
		t.Error("topology", VMTopologyText(4, "2"), VMTopologyText(2, ""))
	}
	xen := map[string]string{"bios-vendor": "Xen", "system-manufacturer": "Xen"}
	if BiosStringsCopied(false, xen) || BiosStringsCopied(true, map[string]string{"bios-vendor": "Dell"}) || BiosStringsCopied(false, nil) ||
		!BiosStringsCopied(false, map[string]string{"bios-vendor": "Dell", "system-manufacturer": "Xen"}) {
		t.Error("BIOS strings")
	}
	if d, ok := VMDescription(true, false, true, "Debian Bookworm 12", "x", nil); ok || d != "" {
		t.Error("a known default template has no description row")
	}
	if d, ok := VMDescription(true, false, true, "Other install media", "raw", nil); !ok || d == "raw" || d == "" {
		t.Errorf("a Misc default template shows its FriendlyNames text: %q", d)
	}
	if d, ok := VMDescription(false, false, false, "vm", "VM imported from physical machine x",
		map[string]string{"p2v_source_machine": "a", "p2v_import_date": "b"}); !ok || d != "" {
		t.Error("P2V")
	}
	if d, ok := VMDescription(true, false, false, "Debian custom", "mine", nil); !ok || d != "mine" {
		t.Error("a user template")
	}
	if CertificateTypeText("ca") != "Trusted CA certificate" || CertificateTypeText("host_internal") != "Server internal identity certificate" {
		t.Error("certificate type")
	}
	from := time.Date(2025, 3, 4, 13, 5, 0, 0, time.Local)
	if got := CertificateText(from, from.AddDate(1, 0, 0), "AB:CD"); got != "Valid from Mar 4, 2025 1:05 PM to Mar 4, 2026 1:05 PM\nThumbprint: AB:CD" {
		t.Errorf("certificate: %q", got)
	}
	if CPUText("AuthenticAMD", "EPYC", 3000) != "Vendor: AuthenticAMD\nModel: EPYC\nSpeed: 3000 MHz" {
		t.Error("CPU")
	}
	if LogDestinationText("") != "Local" || LogDestinationText("10.0.0.1") != "Local and Remote (10.0.0.1)" {
		t.Error("log destination")
	}
}

func TestVirtualizationText(t *testing.T) {
	linux := GuestTools{Running: true, Metrics: true, GuestMetrics: true, PVVersion: map[string]string{"major": "8", "minor": "4"}, UpToDate: true}
	if got := VirtualizationText(linux); got != "Optimized (version 8.4 installed)" {
		t.Errorf("Linux with tools: %q", got)
	}
	old := linux
	old.UpToDate = false
	if got := VirtualizationText(old); got != "Tools out of date (version 8.4 installed)" {
		t.Errorf("old tools: %q", got)
	}
	none := GuestTools{Running: true, Metrics: true}
	if got := VirtualizationText(none); got != "XCP-ng VM Tools not installed" {
		t.Errorf("no tools: %q", got)
	}
	none.Early = true
	if got := VirtualizationText(none); got != "Unknown (please wait)" {
		t.Errorf("just started: %q", got)
	}
	if VirtualizationText(GuestTools{}) != "" {
		t.Error("a halted Linux VM has no row")
	}
	win := GuestTools{Windows: true, Running: true, Metrics: true, GuestMetrics: true, IODrivers: true, Management: "1"}
	if got := VirtualizationText(win); got != "I/O optimized\nManagement Agent installed\nNot able to receive updates from Windows Update" {
		t.Errorf("Windows with tools: %q", got)
	}
	if got := VirtualizationText(GuestTools{Windows: true, VendorDevice: true}); got != "Able to receive updates from Windows Update" {
		t.Errorf("a halted Windows VM: %q", got)
	}
	if got := VirtualizationText(GuestTools{Windows: true, Running: true, Metrics: true, Early: true}); got != "Unknown (please wait)\nNot able to receive updates from Windows Update" {
		t.Errorf("a Windows VM just started: %q", got)
	}
}

func TestReadCachingReason(t *testing.T) {
	if ReadCachingReason(nil, true) != "" || ReadCachingReason([]string{"bogus"}, true) != "" {
		t.Error("no known reason")
	}
	if ReadCachingReason([]string{"LICENSE_RESTRICTION", "NO_RO_IMAGE"}, true) != ReadCachingReason([]string{"NO_RO_IMAGE"}, false) {
		t.Error("the last reason wins")
	}
	if ReadCachingReason([]string{"LICENSE_RESTRICTION"}, true) == ReadCachingReason([]string{"LICENSE_RESTRICTION"}, false) {
		t.Error("licence now or then")
	}
}
