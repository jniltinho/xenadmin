package rules

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"xcpng-center/internal/messages"
)

// DesignedGraph is a graph of the Performance tab (CustomDataGraph.DesignedGraph): its name and the names of its data
// sources ("cpu0", "memory_internal_free"...), in order.
type DesignedGraph struct {
	Name    string   `json:"name"`
	Sources []string `json:"sources"`
}

// GraphLayoutKey, GraphNameKey and DataSourceColourKey are Palette.GetLayoutKey, GetGraphNameKey and GetColorKey: the
// pool's gui_config keys of a server's ("host") or VM's ("vm") graph i, its name and a data source's colour.
func GraphLayoutKey(i int, kind, uuid string) string {
	return fmt.Sprintf("XenCenter.GraphLayout.%d.%s.%s", i, kind, uuid)
}

func GraphNameKey(i int, kind, uuid string) string {
	return fmt.Sprintf("XenCenter.GraphName.%d.%s.%s", i, kind, uuid)
}

func DataSourceColourKey(kind, uuid, source string) string {
	return fmt.Sprintf("XenCenter.DataSource.%s.%s.%s", kind, uuid, source)
}

// graphKey is Palette.LayoutKey and Palette.GraphNameKey (their dots unescaped, as in C#).
var graphKey = regexp.MustCompile(`^XenCenter.(GraphLayout|GraphName).[0-9]+.(host|vm).[a-zA-Z0-9_\-]+$`)

// ParseGraphLayout is GraphList.LoadGraphs: graphs 0, 1... up to the first index without a layout key, each the
// comma-separated data source names and its name; nil when there is none (the default graphs are shown).
func ParseGraphLayout(gui map[string]string, kind, uuid string) []DesignedGraph {
	var out []DesignedGraph
	for i := 0; ; i++ {
		layout, ok := gui[GraphLayoutKey(i, kind, uuid)]
		if !ok {
			return out
		}
		g := DesignedGraph{Name: gui[GraphNameKey(i, kind, uuid)], Sources: []string{}}
		if layout != "" { // a graph with no sources (C# saves a default graph of no devices as "")
			g.Sources = strings.Split(layout, ",")
		}
		out = append(out, g)
	}
}

// GraphLayoutConfig is SaveDataSourceStateAction.Run: gui_config without the object's layout and name keys, then the
// graphs' (a name key only when the name is not empty). No graphs is the default graphs.
func GraphLayoutConfig(gui map[string]string, kind, uuid string, graphs []DesignedGraph) map[string]string {
	out := map[string]string{}
	for k, v := range gui {
		if graphKey.MatchString(k) && strings.Contains(k, uuid) {
			continue
		}
		out[k] = v
	}
	for i, g := range graphs {
		out[GraphLayoutKey(i, kind, uuid)] = strings.Join(g.Sources, ",")
		if g.Name != "" {
			out[GraphNameKey(i, kind, uuid)] = g.Name
		}
	}
	return out
}

// GraphColours is Palette.DefaultColors.
var GraphColours = []string{"#0066ff", "#339900", "#ff9900", "#660099", "#990033", "#995e00", "#ffff00", "#00ffb4", "#009982", "#bfa4a6"}

// ParseColour is Palette.LoadSetColor's value: Color.ToArgb as a decimal int32, as "#rrggbb" (the alpha dropped).
func ParseColour(v string) (string, bool) {
	argb, err := strconv.ParseInt(strings.TrimSpace(v), 10, 32)
	if err != nil {
		return "", false
	}
	return fmt.Sprintf("#%06x", uint32(argb)&0xffffff), true
}

var hexColour = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

// ColourValue is the gui_config value of "#rrggbb": Color.ToArgb of the opaque colour, a decimal int32.
func ColourValue(hex string) (string, error) {
	if !hexColour.MatchString(hex) {
		return "", fmt.Errorf("invalid colour %q", hex)
	}
	rgb, _ := strconv.ParseUint(hex[1:], 16, 32)
	return strconv.Itoa(int(int32(uint32(0xff000000) | uint32(rgb)))), nil
}

// HiddenDataSource is DataSet.Hide: the memory totals are drawn as the used memory, never on their own.
func HiddenDataSource(name string) bool { return name == "memory" || name == "memory_total_kib" }

var (
	cpuSrc          = regexp.MustCompile(`^cpu([0-9]+)$`)
	cpuAvgFreqSrc   = regexp.MustCompile(`^CPU([0-9]+)-avg-freq$`)
	cpuStateSrc     = regexp.MustCompile(`^cpu([0-9]+)-(C|P)([0-9]+)$`)
	cpuOtherSrc     = regexp.MustCompile(`^cpu_avg$`)
	vcpuSrc         = regexp.MustCompile(`^runstate_(blocked|concurrency_hazard|full_contention|fullrun|partial_contention|partial_run)$`)
	vifSrc          = regexp.MustCompile(`^vif_([0-9]+)_(tx|rx)((_errors)?)$`)
	pifEthSrc       = regexp.MustCompile(`^pif_eth([0-9]+)_(tx|rx)((_errors)?)$`)
	pifVlanSrc      = regexp.MustCompile(`^pif_eth([0-9]+).([0-9]+)_(tx|rx)((_errors)?)$`)
	pifBrSrc        = regexp.MustCompile(`^pif_xenbr([0-9]+)_(tx|rx)((_errors)?)$`)
	pifXapiSrc      = regexp.MustCompile(`^pif_xapi([0-9]+)_(tx|rx)((_errors)?)$`)
	pifTapSrc       = regexp.MustCompile(`^pif_tap([0-9]+)_(tx|rx)((_errors)?)$`)
	pifLoSrc        = regexp.MustCompile(`^pif_lo_(tx|rx)((_errors)?)$`)
	pifBondSrc      = regexp.MustCompile(`^pif_(bond[0-9]+)_(tx|rx)((_errors)?)$`)
	pifOtherSrc     = regexp.MustCompile(`^pif_aggr_(tx|rx)$`)
	diskSrc         = regexp.MustCompile(`^vbd_((xvd|hd)[a-z]+)(_(read|write))?(_latency)?$`)
	diskIopsSrc     = regexp.MustCompile(`^vbd_((xvd|hd)[a-z]+)_iops_(read|write|total)$`)
	diskThruSrc     = regexp.MustCompile(`^vbd_((xvd|hd)[a-z]+)_io_throughput_(read|write|total)$`)
	diskOtherSrc    = regexp.MustCompile(`^vbd_((xvd|hd)[a-z]+)_(avgqu_sz|inflight|iowait)$`)
	netLatencySrc   = regexp.MustCompile(`^network/latency$`)
	xapiLatencySrc  = regexp.MustCompile(`^xapi_healthcheck/latency$`)
	xapiMemorySrc   = regexp.MustCompile(`^xapi_(allocation|free_memory|live_memory|memory_usage)_kib$`)
	statefileSrc    = regexp.MustCompile(`^statefile/[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}/latency$`)
	srCacheSrc      = regexp.MustCompile(`^sr_[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}_cache_(size|hits|misses)`)
	srIOSrc         = regexp.MustCompile(`^(io_throughput|iops)_(read|write|total)_([a-f0-9]{8})$`)
	srOtherSrc      = regexp.MustCompile(`^(latency|avgqu_sz|inflight|iowait)_([a-f0-9]{8})$`)
	srReadWriteSrc  = regexp.MustCompile(`^((read|write)(_latency)?)_([a-f0-9]{8})$`)
	gpuSrc          = regexp.MustCompile(`^gpu_((memory_(free|used))|power_usage|temperature|(utilisation_(compute|memory_io)))_((([a-fA-F0-9]{4}/)|([a-fA-F0-9]{8}/))?[a-fA-F0-9]{2}/[0-1][a-fA-F0-9].[0-7])$`)
	oldSrThroughput = regexp.MustCompile(`^io_throughput_(read|write)_([a-f0-9]{8})$`)
	newSrReadWrite  = regexp.MustCompile(`^(read|write)_([a-f0-9]{8})$`)
)

// OfferedDataSources is DataSourceItemList.BuildList's filter of the data source names: no memory totals, and no old
// SR throughput sources when the server has their replacements.
func OfferedDataSources(names []string) []string {
	replaced := slices.ContainsFunc(names, newSrReadWrite.MatchString)
	var out []string
	for _, n := range names {
		if !HiddenDataSource(n) && !(replaced && oldSrThroughput.MatchString(n)) {
			out = append(out, n)
		}
	}
	return out
}

// DataSourceCategory is Helpers.GetDataSourceCategory + ToStringI18N: the Type column of Graph Details.
func DataSourceCategory(name string) string {
	key := "DATATYPE_CUSTOM"
	switch {
	case cpuSrc.MatchString(name) || cpuAvgFreqSrc.MatchString(name) || cpuStateSrc.MatchString(name) || cpuOtherSrc.MatchString(name) || vcpuSrc.MatchString(name):
		key = "DATATYPE_CPU"
	case vifSrc.MatchString(name) || pifEthSrc.MatchString(name) || pifVlanSrc.MatchString(name) || pifBrSrc.MatchString(name) ||
		pifXapiSrc.MatchString(name) || pifBondSrc.MatchString(name) || pifLoSrc.MatchString(name) || pifTapSrc.MatchString(name) || pifOtherSrc.MatchString(name):
		key = "DATATYPE_NETWORK"
	case diskSrc.MatchString(name) || diskIopsSrc.MatchString(name) || diskThruSrc.MatchString(name) || diskOtherSrc.MatchString(name):
		key = "DATATYPE_DISK"
	case srCacheSrc.MatchString(name) || srIOSrc.MatchString(name) || srOtherSrc.MatchString(name) || srReadWriteSrc.MatchString(name):
		key = "DATATYPE_STORAGE"
	case gpuSrc.MatchString(name):
		key = "DATATYPE_GPU"
	case netLatencySrc.MatchString(name) || xapiLatencySrc.MatchString(name) || statefileSrc.MatchString(name):
		key = "DATATYPE_LATENCY"
	case strings.Contains(name, "loadavg") || strings.Contains(name, "hostload"):
		key = "DATATYPE_LOADAVERAGE"
	case strings.HasPrefix(name, "pvsaccelerator"):
		key = "DATATYPE_PVS"
	case xapiMemorySrc.MatchString(name) || strings.HasPrefix(name, "memory"):
		key = "DATATYPE_MEMORY"
	}
	return messages.Get(key)
}

// FriendlyDataSource is Helpers.GetFriendlyDataSourceName (with DataSet's memory_free_kib / memory_internal_free as the
// used memory). find names an object of the server or VM: "vif" (device) and "pif" (a physical NIC's device) give its
// network, "vlan" (device, tag) the VLAN's network, "vbd" (device) its userdevice, "sr" (uuid prefix) the SR; "" when
// there is none, which gives the raw name as in C#.
// shortcut: bonds and GPUs show the raw name (C# names the bond's PIF and the GPU); port PIF.Name / PGPU when asked.
func FriendlyDataSource(name string, find func(what string, ids ...string) string) string {
	switch name {
	case "memory_free_kib":
		name = "memory_used_kib"
	case "memory_internal_free":
		name = "memory_internal_used"
	}
	f := func(key string, args ...any) string { return messages.Format(messages.Friendly(key), args...) }
	found := func(what string, ids ...string) (string, bool) { v := find(what, ids...); return v, v != "" }
	text := ""
	if m := cpuSrc.FindStringSubmatch(name); m != nil {
		text = f("Label-performance.cpu", m[1])
	} else if m := cpuAvgFreqSrc.FindStringSubmatch(name); m != nil {
		text = f("Label-performance.cpu-avg-freq", m[1])
	} else if m := vifSrc.FindStringSubmatch(name); m != nil {
		if n, ok := found("vif", m[1]); ok {
			text = f("Label-performance.vif_"+m[2]+m[3], n)
		}
	} else if m := pifEthSrc.FindStringSubmatch(name); m != nil {
		text = f("Label-performance.nic_"+m[2]+m[3], m[1])
	} else if m := pifVlanSrc.FindStringSubmatch(name); m != nil {
		if n, ok := found("vlan", "eth"+m[1], m[2]); ok {
			text = f("Label-performance.vlan_"+m[3]+m[4], n)
		}
	} else if m := pifBrSrc.FindStringSubmatch(name); m != nil {
		if n, ok := found("pif", "eth"+m[1]); ok {
			text = f("Label-performance.xenbr_"+m[2]+m[3], n)
		}
	} else if m := pifXapiSrc.FindStringSubmatch(name); m != nil {
		text = f("Label-performance.xapi_"+m[2]+m[3], m[1])
	} else if pifBondSrc.MatchString(name) {
		// shortcut above
	} else if m := pifLoSrc.FindStringSubmatch(name); m != nil {
		text = f("Label-performance.lo_" + m[1] + m[2])
	} else if m := pifTapSrc.FindStringSubmatch(name); m != nil {
		text = f("Label-performance.tap_"+m[2]+m[3], m[1])
	} else if m := diskSrc.FindStringSubmatch(name); m != nil {
		if u, ok := found("vbd", m[1]); ok {
			text = f("Label-performance.vbd"+m[3]+m[5], u)
		}
	} else if m := diskIopsSrc.FindStringSubmatch(name); m != nil {
		if u, ok := found("vbd", m[1]); ok {
			text = f("Label-performance.vbd_iops_"+m[3], u)
		}
	} else if m := diskThruSrc.FindStringSubmatch(name); m != nil {
		if u, ok := found("vbd", m[1]); ok {
			text = f("Label-performance.vbd_io_throughput_"+m[3], u)
		}
	} else if m := diskOtherSrc.FindStringSubmatch(name); m != nil {
		if u, ok := found("vbd", m[1]); ok {
			text = f("Label-performance.vbd_"+m[3], u)
		}
	} else if m := srCacheSrc.FindStringSubmatch(name); m != nil {
		text = f("Label-performance.sr_cache_" + m[1])
	} else if m := cpuStateSrc.FindStringSubmatch(name); m != nil {
		text = f("Label-performance.cpu-state", m[1], m[2], m[3])
	} else if m := srIOSrc.FindStringSubmatch(name); m != nil {
		if n, ok := found("sr", m[3]); ok {
			text = f("Label-performance.sr_"+m[1]+"_"+m[2], Ellipsise(n, 30))
		}
	} else if m := srOtherSrc.FindStringSubmatch(name); m != nil {
		if n, ok := found("sr", m[2]); ok {
			text = f("Label-performance.sr_"+m[1], Ellipsise(n, 30))
		}
	} else if m := srReadWriteSrc.FindStringSubmatch(name); m != nil {
		if n, ok := found("sr", m[4]); ok {
			text = f("Label-performance.sr_rw_"+m[1], Ellipsise(n, 30))
		}
	} else if gpuSrc.MatchString(name) {
		// shortcut above
	} else if netLatencySrc.MatchString(name) {
		text = messages.Friendly("Label-performance.network_latency")
	} else if xapiLatencySrc.MatchString(name) {
		text = messages.Friendly("Label-performance.xapi_latency")
	} else if statefileSrc.MatchString(name) {
		text = messages.Friendly("Label-performance.statefile_latency")
	} else if strings.Contains(name, "loadavg") {
		text = messages.Friendly("Label-performance.loadavg")
	} else if strings.Contains(name, "hostload") {
		text = messages.Friendly("OPERATIONAL_METRICS_HOSTLOAD")
	} else {
		text = messages.Friendly("Label-performance." + name)
	}
	if text == "" {
		return name
	}
	return text
}
