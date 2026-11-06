package rules

import (
	"slices"
	"strconv"
	"strings"

	"xcpng-center/internal/messages"
)

// PoolFacts is what PoolJoinRules reads of a connected server (its pool of one or its pool).
type PoolFacts struct {
	Key             string // the connection
	Name            string // the coordinator's name
	Hosts           int
	PoolName        string // a named pool, or one of more hosts, is a pool (Pool.IsVisible)
	CPUVendor       string
	APIMajor        int
	APIMinor        int
	DBSchema        string // software_version db_schema
	PlatformVersion string
	NetworkBackend  string
	SharedStorage   bool // an SR shared, other than the tools SR
	RunningVMs      bool
	HA              bool
	Clustered       bool
	TLSVerification bool
	CPU             map[string]string // a member's host cpu_info
	PoolCPU         map[string]string // a coordinator's pool cpu_info (the features common to its hosts)
	ReadError       string            // the server could not be read: it is offered for nothing
}

// JoinReason is PoolJoinRules.Reason, in its order (the first reason found is the one shown).
type JoinReason int

const (
	JoinWillBeCoordinator JoinReason = iota
	JoinAllowed
	JoinIsAPool
	JoinDifferentCPUs
	JoinDifferentServerVersion
	JoinHasSharedStorage
	JoinHasRunningVMs
	JoinDifferentNetworkBackends
	JoinCoordinatorHasHA
	JoinHasClusteringEnabled
	JoinTLSOnlyOnPool
	JoinTLSOnlyOnPoolJoiner
	JoinTLSOnlyOnCoordinator
	JoinTLSOnlyOnCoordinatorJoiner
	JoinNotRead // the server could not be read (not a C# reason: its cache would be empty, "Connecting...")
)

// CanJoinPool is PoolJoinRules.CanJoinPool for XCP-ng: the licensing, AD, RBAC, homogeneous update and management
// interface checks of the C# app are left to pool.join, whose refusal the window shows in the C# words.
// shortcut: those reasons are not shown before the user picks the server; port them if pool.join refusals confuse
func CanJoinPool(member, coordinator PoolFacts) JoinReason {
	switch {
	case member.ReadError != "" || coordinator.ReadError != "":
		return JoinNotRead
	case member.Hosts > 1 || member.PoolName != "":
		return JoinIsAPool
	case member.Key == coordinator.Key:
		return JoinWillBeCoordinator
	case member.CPUVendor != "" && coordinator.CPUVendor != "" && member.CPUVendor != coordinator.CPUVendor:
		return JoinDifferentCPUs // CompatibleCPUs: feature levelling makes one vendor's CPUs compatible
	case member.APIMajor != coordinator.APIMajor || member.APIMinor != coordinator.APIMinor ||
		member.DBSchema != coordinator.DBSchema || member.PlatformVersion != coordinator.PlatformVersion:
		return JoinDifferentServerVersion
	case member.SharedStorage:
		return JoinHasSharedStorage
	case member.RunningVMs:
		return JoinHasRunningVMs
	case member.NetworkBackend != "" && coordinator.NetworkBackend != "" && member.NetworkBackend != coordinator.NetworkBackend:
		return JoinDifferentNetworkBackends
	case coordinator.HA:
		return JoinCoordinatorHasHA
	case member.Clustered:
		return JoinHasClusteringEnabled
	case coordinator.TLSVerification && !member.TLSVerification:
		if coordinator.Hosts > 1 || coordinator.PoolName != "" {
			return JoinTLSOnlyOnPool
		}
		return JoinTLSOnlyOnCoordinator
	case !coordinator.TLSVerification && member.TLSVerification:
		if coordinator.Hosts > 1 || coordinator.PoolName != "" {
			return JoinTLSOnlyOnPoolJoiner
		}
		return JoinTLSOnlyOnCoordinatorJoiner
	}
	return JoinAllowed
}

// JoinReasonText is PoolJoinRules.ReasonMessage.
func JoinReasonText(r JoinReason) string {
	return messages.Get(map[JoinReason]string{
		JoinWillBeCoordinator: "COORDINATOR", JoinIsAPool: "NEWPOOL_IS_A_POOL", JoinDifferentCPUs: "NEWPOOL_DIFF_HARDWARE",
		JoinDifferentServerVersion: "NEWPOOL_DIFF_SERVER", JoinHasSharedStorage: "NEWPOOL_HAS_SHARED_STORAGE",
		JoinHasRunningVMs: "NEWPOOL_HAS_RUNNING_VMS", JoinDifferentNetworkBackends: "NEWPOOL_DIFFERENT_NETWORK_BACKENDS",
		JoinCoordinatorHasHA: "POOL_JOIN_FORBIDDEN_BY_HA", JoinHasClusteringEnabled: "NEW_POOL_CLUSTERING_ENABLED",
		JoinTLSOnlyOnPool: "POOL_JOIN_CERTIFICATE_CHECKING_ONLY_ON_POOL", JoinTLSOnlyOnPoolJoiner: "POOL_JOIN_CERTIFICATE_CHECKING_ONLY_ON_POOL_JOINER",
		JoinTLSOnlyOnCoordinator: "POOL_JOIN_CERTIFICATE_CHECKING_ONLY_ON_COORDINATOR", JoinTLSOnlyOnCoordinatorJoiner: "POOL_JOIN_CERTIFICATE_CHECKING_ONLY_ON_COORDINATOR_JOINER",
		JoinNotRead: "CONNECTING",
	}[r])
}

// CanBeCoordinator is Candidate.CanBeCoordinator: any connected server, read, that is not a pool already.
func CanBeCoordinator(r JoinReason) bool { return r != JoinIsAPool && r != JoinNotRead }

// CPUWarning is PoolRules.Warnings' CPU part: the members with fewer CPU features than the pool (the pool levels
// down), with more (they level down), or both; "" when none. The window asks before creating the pool.
func CPUWarning(members []PoolFacts, coordinator PoolFacts) string {
	var fewer, more []string
	for _, m := range members {
		if fewerFeatures(m.CPU, coordinator.PoolCPU) {
			fewer = append(fewer, m.Name)
		}
		if fewerFeatures(coordinator.PoolCPU, m.CPU) {
			more = append(more, m.Name)
		}
	}
	pick := func(names []string, key string) string {
		if len(names) == 1 {
			return messages.Format(messages.Get(key), names[0])
		}
		return messages.Format(messages.Get(key+"_MULTIPLE"), strings.Join(names, "\n"))
	}
	switch {
	case len(fewer) > 0 && len(more) > 0:
		return pick(union(fewer, more), "ADD_HOST_TO_POOL_CPU_DOWN_LEVEL_POOL_AND_HOST_MESSAGE")
	case len(fewer) > 0:
		return pick(fewer, "ADD_HOST_TO_POOL_CPU_DOWN_LEVEL_POOL_MESSAGE")
	case len(more) > 0:
		return pick(more, "ADD_HOST_TO_POOL_CPU_DOWN_LEVEL_HOST_MESSAGE")
	}
	return ""
}

func union(a, b []string) []string {
	out := append([]string{}, a...)
	for _, x := range b {
		if !slices.Contains(out, x) {
			out = append(out, x)
		}
	}
	return out
}

// fewerFeatures is PoolJoinRules.FewerFeatures: a has fewer CPU features than b in features_hvm or features_pv.
func fewerFeatures(a, b map[string]string) bool {
	return fewerFeatureBits(a["features_hvm"], b["features_hvm"]) || fewerFeatureBits(a["features_pv"], b["features_pv"])
}

// fewerFeatureBits compares two hex feature masks ("178bfbff-f6f83203-..."): a lacks a bit b has, in some 32-bit word.
func fewerFeatureBits(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	clean := strings.NewReplacer(" ", "", "-", "")
	a, b = clean.Replace(a), clean.Replace(b)
	for len(a) < len(b) {
		a += "0"
	}
	for len(b) < len(a) {
		b += "0"
	}
	for i := 0; i+8 <= len(a); i += 8 {
		x, errA := strconv.ParseUint(a[i:i+8], 16, 32)
		y, errB := strconv.ParseUint(b[i:i+8], 16, 32)
		if errA == nil && errB == nil && x&y != y {
			return true
		}
	}
	return false
}
