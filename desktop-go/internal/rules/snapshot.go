package rules

import (
	"strings"

	"xcpng-center/internal/messages"
)

// SnapshotModes is VmRules.SnapshotModes: whether the VM allows a disk snapshot and a disk and memory snapshot
// (checkpoint), and why the second is not offered. Quiesced snapshots are not shown (QuebecOrGreater: every
// XCP-ng release).
type SnapshotModes struct {
	Disk       bool   `json:"disk"`
	Memory     bool   `json:"memory"`
	MemoryInfo string `json:"memory_info"`
}

// SnapshotModesFor is VmRules.SnapshotModesFor; ioDrivers is the guest's PV I/O drivers being detected.
func SnapshotModesFor(allowed []string, running, ioDrivers bool) SnapshotModes {
	info := messages.Get("INFO_DISKMEMORY_MODE_MISC")
	switch {
	case !running:
		info = messages.Get("INFO_DISKMEMORY_MODE_POWER_STATE")
	case !ioDrivers:
		info = messages.Get("INFO_DISKMEMORY_MODE_NO_IO_DRIVERS")
	}
	has := func(op string) bool {
		for _, a := range allowed {
			if a == op {
				return true
			}
		}
		return false
	}
	return SnapshotModes{Disk: has("snapshot"), Memory: has("checkpoint"), MemoryInfo: info}
}

// SnapshotTypeText is VmRules.SnapshotTypeText: a suspended snapshot holds the memory too.
func SnapshotTypeText(powerState string) string {
	if powerState == "Suspended" {
		return messages.Get("DISKS_AND_MEMORY")
	}
	return messages.Get("DISKS_ONLY")
}

// RevertQuestion is SnapshotsView.Revert's question and its "snapshot first" option ("" when not allowed).
func RevertQuestion(snapshot string, snapshotAllowed bool) (text, option string) {
	text = messages.Format(messages.Get("SNAPSHOT_REVERT_BLURB"), snapshot)
	if snapshotAllowed {
		option = strings.ReplaceAll(messages.Get("SNAPSHOT_REVERT_NEW_SNAPSHOT"), "&", "")
	}
	return text, option
}

// DeleteSnapshotQuestion is SnapshotsView.Delete's question.
func DeleteSnapshotQuestion(snapshot string) string {
	return messages.Format(messages.Get("REMOVE_SNAPSHOT_TEXT"), messages.Format(messages.Get("THE_SNAPSHOT"), snapshot))
}
