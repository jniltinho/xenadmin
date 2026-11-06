package rules

import (
	"regexp"
	"strconv"
	"strings"

	"xcpng-center/internal/messages"
)

var macRegex = regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`) // Helpers.IsValidMAC

// ValidateVIF is NetworkRules.Validate (VIFDialog.ValidateInput): ok when the dialog's OK can be pressed, else the
// reason ("" when there is nothing to say).
func ValidateVIF(hasNetwork, manualMAC bool, mac string, qos bool, limit string) (why string, ok bool) {
	if !hasNetwork {
		return messages.Get("SELECT_NETWORK_TOOLTIP"), false
	}
	if manualMAC {
		if strings.TrimSpace(mac) == "" {
			return "", false
		}
		if !macRegex.MatchString(mac) {
			return messages.Get("MAC_INVALID"), false
		}
	}
	if qos {
		if strings.TrimSpace(limit) == "" {
			return "", false
		}
		if kbps, err := strconv.Atoi(limit); err != nil || kbps <= 0 {
			return messages.Get("ENTER_VALID_QOS"), false
		}
	}
	return "", true
}

// VIFDevice is VIF.GetDeviceId: the lowest device number no interface of the VM uses.
func VIFDevice(used []string) int {
	taken := map[string]bool{}
	for _, d := range used {
		taken[d] = true
	}
	i := 0
	for taken[strconv.Itoa(i)] {
		i++
	}
	return i
}

// VIFTooltips is NetworkRules.Tooltips: why Properties, Remove and Activate/Deactivate may fail (a suspended VM, or a
// running one without I/O drivers); empty when none.
func VIFTooltips(suspended, running, ioDrivers, windows, attached bool) (edit, remove, toggle string) {
	tools := func(key string) string { return messages.Format(messages.Get(key), "XCP-ng VM Tools") }
	edit = tools("TOOLTIP_EDIT_NETWORK_TOOLS")
	if windows { // VM.HasNewVirtualizationStates
		edit = messages.Get("TOOLTIP_EDIT_NETWORK_IO_DRIVERS")
	}
	state := map[bool]string{true: "TOOLTIP_DEACTIVATE_VIF_", false: "TOOLTIP_ACTIVATE_VIF_"}[attached]
	switch {
	case suspended:
		return edit, messages.Get("TOOLTIP_REMOVE_NETWORK_SUSPENDED"), messages.Get(state + "SUSPENDED")
	case running && !ioDrivers:
		remove = tools("TOOLTIP_REMOVE_NETWORK_TOOLS")
		if windows {
			remove = messages.Get("TOOLTIP_REMOVE_NETWORK_IO_DRIVERS")
		}
		return edit, remove, messages.Get(state + "TOOLS")
	}
	return "", "", ""
}
