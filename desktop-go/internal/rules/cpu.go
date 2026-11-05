package rules

import (
	"math"

	"xcpng-center/internal/messages"
)

// MaxSockets is VM.MAX_SOCKETS (the Xen limit, CA-198276).
const MaxSockets = 16

// Topology is a split of the vCPUs into sockets and cores per socket.
type Topology struct {
	Sockets int64  `json:"sockets"`
	Cores   int64  `json:"cores"`
	Text    string `json:"text"`
}

// TopologyText is VM.GetTopology ("2 sockets with 2 cores per socket"; sockets 0 for an invalid value).
func TopologyText(sockets, cores int64) string {
	switch {
	case sockets == 0 && cores == 1:
		return messages.Get("CPU_TOPOLOGY_STRING_INVALID_VALUE_1")
	case sockets == 0:
		return messages.Format(messages.Get("CPU_TOPOLOGY_STRING_INVALID_VALUE"), cores)
	case sockets == 1 && cores == 1:
		return messages.Get("CPU_TOPOLOGY_STRING_1_SOCKET_1_CORE")
	case sockets == 1:
		return messages.Format(messages.Get("CPU_TOPOLOGY_STRING_1_SOCKET_N_CORE"), cores)
	case cores == 1:
		return messages.Format(messages.Get("CPU_TOPOLOGY_STRING_N_SOCKET_1_CORE"), sockets)
	}
	return messages.Format(messages.Get("CPU_TOPOLOGY_STRING_N_SOCKET_N_CORE"), sockets, cores)
}

// ValidTopologies is VmRules.ValidTopologies: the (sockets, cores) splits of the vCPUs within the cores per socket
// of the server and the socket limit.
func ValidTopologies(vcpus, maxCoresPerSocket int64) []Topology {
	limit := vcpus
	if maxCoresPerSocket > 0 && maxCoresPerSocket < limit {
		limit = maxCoresPerSocket
	}
	out := []Topology{}
	for cores := int64(1); cores <= limit; cores++ {
		if vcpus%cores == 0 && vcpus/cores <= MaxSockets {
			out = append(out, Topology{Sockets: vcpus / cores, Cores: cores, Text: TopologyText(vcpus/cores, cores)})
		}
	}
	return out
}

// ValidVCPUConfiguration is VM.ValidVCPUConfiguration: "" when the vCPUs fit the cores per socket.
func ValidVCPUConfiguration(vcpus, cores int64) string {
	if cores > 0 {
		if vcpus%cores != 0 {
			return messages.Get("CPU_TOPOLOGY_INVALID_REASON_MULTIPLE")
		}
		if vcpus/cores > MaxSockets {
			return messages.Get("CPU_TOPOLOGY_INVALID_REASON_SOCKETS")
		}
	}
	return ""
}

// Priority is CpuPage's slider position (0-8) for a vCPU weight: log4(weight), rounded.
func Priority(weight int) int { return int(math.Round(math.Log(float64(weight)) / math.Log(4))) }

// Weight is the vCPU weight of a slider position: 4^priority.
func Weight(priority int) int { return int(math.Pow(4, float64(priority))) }
