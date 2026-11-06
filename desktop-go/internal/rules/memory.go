package rules

import (
	"fmt"

	"xcpng-center/internal/messages"
)

// MemoryVM is what MemoryRules reads of a VM: its memory limits, whether it balloons and the limits of its server.
type MemoryVM struct {
	StaticMin, StaticMax, DynamicMin, DynamicMax int64
	Halted, Running, Template                    bool
	Ballooning                                   bool    // VM.SupportsBallooning: the guest reports feature-balloon (a template: dynamic min below static max)
	IODrivers, Windows                           bool    // for DynamicUnavailable's reason
	Ratio                                        float64 // MemoryRules.Ratio: the pool's memory-ratio-hvm/pv, else 0.25
	MaxDynamicMin                                int64   // MemoryRules.MaxDynMin: the free memory of its server for it, -1 when not running
	MaxAllowed                                   int64   // VM.MaxMemAllowed: the memory-static-max restriction, else 1 TiB
}

const mib = 1 << 20

// UsesBallooning is VM.UsesBallooning: the advanced controls (dynamic range and static maximum).
func (m MemoryVM) UsesBallooning() bool { return m.DynamicMax != m.StaticMax && m.Ballooning }

// DynamicUnavailable is MemoryRules.DynamicUnavailable: why a memory range cannot be set ("" when it can).
func (m MemoryVM) DynamicUnavailable() string {
	switch {
	case m.Ballooning:
		return ""
	case m.Template:
		return messages.Get("DMC_UNAVAILABLE_TEMPLATE")
	case m.IODrivers:
		return messages.Get("DMC_UNAVAILABLE_NOTSUPPORTED")
	case m.Windows:
		return messages.Get("DMC_UNAVAILABLE_NO_IO_NO_MGMNT")
	}
	return messages.Format(messages.Get("DMC_UNAVAILABLE_NOTOOLS"), "XCP-ng VM Tools")
}

// Rubric is MemoryRules.Rubric.
func (m MemoryVM) Rubric() string {
	key := "BALLOONING_RUBRIC"
	if m.UsesBallooning() {
		key += "_ADVANCED"
	}
	if m.Template {
		key += "_TEMPLATE"
	}
	return messages.Get(key)
}

func (m MemoryVM) memoryMax() int64 { return max(m.StaticMax, m.MaxAllowed) }

// ValidateMemory is MemoryRules.Validate: "" when the dynamic minimum, dynamic maximum and static maximum can be set.
func ValidateMemory(m MemoryVM, dynMin, dynMax, staticMax int64) string {
	lowest := max(m.StaticMin, int64(float64(staticMax)*m.Ratio))
	switch {
	case staticMax < mib || staticMax > m.memoryMax():
		return fmt.Sprintf("The memory must be between 1 MB and %s.", MemoryText(m.memoryMax()))
	case m.MaxDynamicMin >= 0 && float64(staticMax)*m.Ratio > float64(m.MaxDynamicMin) && staticMax > m.StaticMax:
		return "The server does not have enough free memory for this maximum."
	case dynMax > staticMax || dynMin > dynMax:
		return "The dynamic minimum must not be above the dynamic maximum, nor the dynamic maximum above the static maximum."
	case dynMin < lowest && (staticMax != m.StaticMax || dynMin < m.DynamicMin):
		return fmt.Sprintf("The dynamic minimum must be at least %s.", MemoryText(lowest))
	case m.MaxDynamicMin >= 0 && dynMin > m.MaxDynamicMin && dynMin > m.DynamicMin:
		return "The server does not have enough free memory for this dynamic minimum."
	}
	return ""
}

// MemorySettings is MemoryRules.Settings: a static maximum unchanged to the MB keeps the server's exact bytes (and so
// do the dynamic values equal to it).
func MemorySettings(m MemoryVM, dynMin, dynMax, staticMax int64) (int64, int64, int64) {
	if m.StaticMax/mib == staticMax/mib {
		if dynMin == staticMax {
			dynMin = m.StaticMax
		}
		if dynMax == staticMax {
			dynMax = m.StaticMax
		}
		staticMax = m.StaticMax
	}
	return dynMin, dynMax, staticMax
}

// MemoryConfirmation is MemoryRules.Confirmation: the question before a change that reboots the VM ("" for none).
func MemoryConfirmation(m MemoryVM, staticMax int64) string {
	switch {
	case m.Halted || m.StaticMax == staticMax:
		return ""
	case m.UsesBallooning():
		return messages.Get("CONFIRM_CHANGE_STATIC_MAX_SINGULAR")
	case m.Ballooning:
		return messages.Get("CONFIRM_CHANGE_MEMORY_MAX_SINGULAR")
	}
	return messages.Get("CONFIRM_CHANGE_MEMORY_SINGULAR")
}

// MemoryText is Util.MemorySizeStringSuitableUnits: MB up to 1 GB, GB above (one decimal when needed).
func MemoryText(bytes int64) string {
	if bytes < 1<<30 {
		return fmt.Sprintf("%d MB", bytes/mib)
	}
	gb := float64(bytes) / (1 << 30)
	if gb == float64(int64(gb)) {
		return fmt.Sprintf("%d GB", int64(gb))
	}
	return fmt.Sprintf("%.1f GB", gb)
}
