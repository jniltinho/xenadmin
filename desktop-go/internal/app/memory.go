package app

import (
	"errors"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// MemoryEdit is the Memory Settings dialog (BallooningDialog) of a VM: its values, its mode and its limits.
type MemoryEdit struct {
	Name        string `json:"name"`
	CanEdit     bool   `json:"can_edit"`    // MemoryRules.CanEdit
	Advanced    bool   `json:"advanced"`    // dynamic minimum, maximum and static maximum
	Dynamic     bool   `json:"dynamic"`     // basic mode may set a range
	Unavailable string `json:"unavailable"` // why it cannot ("" when it can)
	Rubric      string `json:"rubric"`
	InMB        bool   `json:"in_mb"` // MemoryRules.InMegabytes
	StaticMin   int64  `json:"static_min"`
	StaticMax   int64  `json:"static_max"`
	DynamicMin  int64  `json:"dynamic_min"`
	DynamicMax  int64  `json:"dynamic_max"`
	Max         int64  `json:"max"` // MemoryRules.MemoryMax
}

// memoryVM reads rules.MemoryVM of a VM, with its virtualization state.
func memoryVM(s *xenapi.Session, ref xenapi.VMRef, vm xenapi.VMRecord) (rules.MemoryVM, bool, string, error) {
	io, win, status := vmTools(s, vm)
	m, restricted, err := ops.MemoryFacts(s, ref, vm, io, win)
	return m, restricted, status, err
}

// vmTools is VM.GetVirtualizationStatus's IoDriversInstalled flag, VM.HasNewVirtualizationStates and the state's name.
func vmTools(s *xenapi.Session, vm xenapi.VMRecord) (bool, bool, string) {
	gm, _ := xenapi.VMGuestMetrics.GetRecord(s, vm.GuestMetrics)
	m, _ := xenapi.VMMetrics.GetRecord(s, vm.Metrics)
	_, status, _ := virtualizationStatus(s, vm, gm, m)
	return strings.Contains(status, "IoDriversInstalled"), isWindows(s, vm), status
}

// ReadMemoryEdit reads a VM for its Memory Settings dialog.
func ReadMemoryEdit(c *xapi.Connection, vmUUID string) (MemoryEdit, error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return MemoryEdit{}, xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return MemoryEdit{}, xapi.Sanitize(err)
	}
	m, restricted, status, err := memoryVM(s, ref, vm)
	if err != nil {
		return MemoryEdit{}, err
	}
	e := MemoryEdit{Name: vm.NameLabel, Advanced: m.UsesBallooning(), Dynamic: m.Ballooning, Rubric: m.Rubric(),
		InMB: m.StaticMax <= 1<<30, StaticMin: m.StaticMin, StaticMax: m.StaticMax, DynamicMin: m.DynamicMin,
		DynamicMax: m.DynamicMax, Max: max(m.StaticMax, m.MaxAllowed)}
	e.CanEdit = len(vm.CurrentOperations) == 0 && (m.Halted || m.Running && !strings.Contains(status, "Unknown"))
	if restricted { // MemoryRules.DynamicUnavailable: the licence comes first
		e.Unavailable = messages.Get("DMC_UNAVAILABLE_LICENSE_RESTRICTION")
	} else {
		e.Unavailable = m.DynamicUnavailable()
	}
	return e, nil
}

// MemoryQuestion is MemoryRules.Validate and Confirmation for the dialog's values: an error when they cannot be set,
// else the question to ask first ("" for none).
func MemoryQuestion(c *xapi.Connection, vmUUID string, dynMin, dynMax, staticMax int64) (string, error) {
	s := c.Session
	ref, err := xenapi.VM.GetByUUID(s, vmUUID)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	vm, err := xenapi.VM.GetRecord(s, ref)
	if err != nil {
		return "", xapi.Sanitize(err)
	}
	m, _, _, err := memoryVM(s, ref, vm)
	if err != nil {
		return "", err
	}
	if !m.Ballooning { // as ops.ChangeMemory: one value without a range
		dynMin, dynMax = staticMax, staticMax
	}
	if why := rules.ValidateMemory(m, dynMin, dynMax, staticMax); why != "" {
		return "", errors.New(why)
	}
	_, _, staticMax = rules.MemorySettings(m, dynMin, dynMax, staticMax)
	return rules.MemoryConfirmation(m, staticMax), nil
}

// ChangeMemory applies the dialog's values (ops.ChangeMemory).
func ChangeMemory(c *xapi.Connection, vmUUID string, dynMin, dynMax, staticMax int64) error {
	return ops.ChangeMemory(c, vmUUID, dynMin, dynMax, staticMax, func(vm xenapi.VMRecord) (bool, bool) {
		io, win, _ := vmTools(c.Session, vm)
		return io, win
	})
}
