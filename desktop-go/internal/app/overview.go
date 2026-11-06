package app

import (
	"fmt"
	"math"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/trust"
	"xcpng-center/internal/xapi"
)

// OverviewRow is one object's live columns of the Search overview (OverviewColumns.Values): CPU, memory, disks,
// network, address and uptime; Bar is 0..100 for the CPU and memory bars, -1 for none.
type OverviewRow struct {
	CPU       string `json:"cpu"`
	CPUBar    int    `json:"cpu_bar"`
	Memory    string `json:"memory"`
	MemoryBar int    `json:"memory_bar"`
	Disks     string `json:"disks"`
	Network   string `json:"network"`
	Address   string `json:"address"`
	Uptime    string `json:"uptime"`
}

// Overview reads the live columns of the connection's servers and VMs (MetricUpdater: rrd_updates of the last ten
// seconds from each live server, its VMs with it), keyed by uuid. A server whose figures cannot be read shows "-".
func Overview(c *xapi.Connection, policy *trust.Policy) (map[string]OverviewRow, error) {
	s := c.Session
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	vms, err := xenapi.VM.GetAllRecords(s)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	pifs, _ := xenapi.PIF.GetAllRecords(s)
	// the server's clock (Connection.ServerTimeOffset): the samples asked for and the uptimes are in its time
	now := time.Now()
	if pools, err := xenapi.Pool.GetAllRecords(s); err == nil {
		for _, p := range pools {
			if t, err := xenapi.Host.GetServertime(s, p.Master); err == nil && !t.IsZero() {
				now = t
			}
		}
	}
	out := map[string]OverviewRow{}
	values := map[string]map[string]float64{} // uuid -> name -> the newest value
	for _, h := range hosts {
		live, _ := xenapi.HostMetrics.GetLive(s, h.Metrics)
		if !live {
			continue
		}
		query := fmt.Sprintf("start=%d&cf=AVERAGE&interval=5&host=true", now.Add(-10*time.Second).Unix())
		u, err := fetchUpdates(s, policy, h.Address, c.Port, query)
		if err != nil {
			continue
		}
		for key, vals := range u.series {
			parts := strings.SplitN(key, ":", 3) // "vm:uuid:name" or "host:uuid:name"
			if len(parts) != 3 {
				continue
			}
			v := math.NaN()
			for i := len(vals) - 1; i >= 0 && math.IsNaN(v); i-- {
				v = vals[i]
			}
			if values[parts[1]] == nil {
				values[parts[1]] = map[string]float64{}
			}
			values[parts[1]][parts[2]] = v
		}
	}
	// MetricUpdater.GetValue: 0 for a figure the server did not send
	get := func(uuid, name string) float64 { return values[uuid][name] }
	for _, h := range hosts {
		r := OverviewRow{CPUBar: -1, MemoryBar: -1, Disks: "-", Address: "-", Uptime: "-"}
		r.CPU, r.CPUBar = cpuText(len(h.HostCPUs), func(i int) float64 { return get(h.UUID, fmt.Sprintf("cpu%d", i)) })
		free, total := get(h.UUID, "memory_free_kib"), get(h.UUID, "memory_total_kib") // hostMemoryUsageString
		r.Memory, r.MemoryBar = "-", -1
		if total != 0 && !math.IsNaN(total) && !math.IsNaN(free) {
			r.Memory = messages.Format(messages.Get("QUERY_MEMORY_USAGE"), rules.MemorySizeText((total-free)*1024, false), rules.MemorySizeText(total*1024, false))
			r.MemoryBar = int(math.Round((total - free) * 100 / total))
		}
		var nics []string
		var addrs []string
		for _, p := range h.PIFs {
			pif := pifs[p]
			if pif.Physical {
				nics = append(nics, pif.Device)
			}
			if ip := net.ParseIP(pif.IP); ip != nil { // the search's ip_address of a server: its PIFs'
				addrs = append(addrs, ip.String())
			}
		}
		// hostNetworkUsageString: each physical NIC's rx plus "vbd_<device>_tx", as the C# reads it
		r.Network = dataAvgMax(nics, func(d string) float64 { return get(h.UUID, "pif_"+d+"_rx") + get(h.UUID, "vbd_"+d+"_tx") })
		if len(addrs) > 0 {
			r.Address = strings.Join(addrs, ", ")
		}
		if boot, err := strconv.ParseFloat(h.OtherConfig["boot_time"], 64); err == nil && boot != 0 { // Host.Uptime
			r.Uptime = prettyTimeSpan(now.Sub(time.Unix(int64(boot), 0)))
		}
		out[h.UUID] = r
	}
	gms, _ := xenapi.VMGuestMetrics.GetAllRecords(s)
	metrics, _ := xenapi.VMMetrics.GetAllRecords(s)
	vbds, _ := xenapi.VBD.GetAllRecords(s)
	vifs, _ := xenapi.VIF.GetAllRecords(s)
	for _, vm := range vms {
		if vm.IsATemplate || vm.IsASnapshot || vm.IsControlDomain {
			continue
		}
		r := OverviewRow{CPU: "-", CPUBar: -1, Memory: "-", MemoryBar: -1, Disks: "-", Network: "-", Address: "-", Uptime: "-"}
		m := metrics[vm.Metrics]
		if vm.PowerState == xenapi.VMPowerStateRunning {
			r.CPU, r.CPUBar = cpuText(m.VCPUsNumber, func(i int) float64 { return get(vm.UUID, fmt.Sprintf("cpu%d", i)) })
			free, total := get(vm.UUID, "memory_internal_free"), get(vm.UUID, "memory") // vmMemoryUsageString
			if total != 0 && !math.IsNaN(total) && !math.IsNaN(free) && total >= free*1024 {
				r.Memory = messages.Format(messages.Get("QUERY_MEMORY_USAGE"), rules.MemorySizeText(total-free*1024, false), rules.MemorySizeText(total, false))
				r.MemoryBar = int(math.Round((total - free*1024) * 100 / total))
			}
			var devs []string
			for _, b := range vm.VBDs {
				devs = append(devs, vbds[b].Device)
			}
			r.Disks = dataAvgMax(devs, func(d string) float64 { return get(vm.UUID, "vbd_"+d+"_read") + get(vm.UUID, "vbd_"+d+"_write") })
			devs = nil
			for _, v := range vm.VIFs {
				devs = append(devs, vifs[v].Device)
			}
			// vmNetworkUsageString: each VIF's rx plus "vbd_<device>_tx", as the C# reads it
			r.Network = dataAvgMax(devs, func(d string) float64 { return get(vm.UUID, "vif_"+d+"_rx") + get(vm.UUID, "vbd_"+d+"_tx") })
		}
		// the search's ip_address of a VM: the guest's addresses of its VIFs, distinct
		if gm, ok := gms[vm.GuestMetrics]; ok {
			seen := map[string]bool{}
			var addrs []string
			for _, v := range vm.VIFs {
				dev := vifs[v].Device
				keys := make([]string, 0, len(gm.Networks))
				for k := range gm.Networks {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys { // Helpers.FindIpAddresses: a value may hold several, split on "\n" or "%n"
					if !strings.HasPrefix(k, dev+"/ip") {
						continue
					}
					for _, a := range strings.FieldsFunc(strings.ReplaceAll(gm.Networks[k], "%n", "\n"), func(r rune) bool { return r == '\n' }) {
						if ip := net.ParseIP(strings.TrimSpace(a)); ip != nil && !seen[ip.String()] {
							seen[ip.String()], addrs = true, append(addrs, ip.String())
						}
					}
				}
			}
			if len(addrs) > 0 {
				r.Address = strings.Join(addrs, ", ")
			}
		}
		running := vm.PowerState == xenapi.VMPowerStateRunning || vm.PowerState == xenapi.VMPowerStatePaused || vm.PowerState == xenapi.VMPowerStateSuspended
		if running && !m.StartTime.IsZero() && m.StartTime.Unix() > 0 { // VM.RunningTime
			r.Uptime = prettyTimeSpan(now.Sub(m.StartTime))
		}
		out[vm.UUID] = r
	}
	return out, nil
}

// cpuText is the CPU column: "N% of 1 CPU" or "N% of M CPUs" (the mean), and its bar; "-" without CPUs.
func cpuText(n int, value func(int) float64) (string, int) {
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += value(i)
	}
	if n == 0 || math.IsNaN(sum) {
		return "-", -1
	}
	if n == 1 {
		return messages.Format(messages.Get("QUERY_PERCENT_OF_CPU"), fmt.Sprintf("%.0f", sum*100)), int(math.Round(sum * 100))
	}
	return messages.Format(messages.Get("QUERY_PERCENT_OF_CPUS"), fmt.Sprintf("%.0f", sum*100/float64(n)), n), int(math.Round(sum * 100 / float64(n)))
}

// dataAvgMax is the disk and network columns: "avg/max" in KB/s over the devices ("-" for none).
func dataAvgMax(devices []string, value func(string) float64) string {
	sum, top := 0.0, 0.0
	for _, d := range devices {
		v := value(d)
		sum += v
		if v > top {
			top = v
		}
	}
	if len(devices) == 0 || math.IsNaN(sum) {
		return "-"
	}
	return messages.Format(messages.Get("QUERY_DATA_AVG_MAX"), fmt.Sprintf("%.0f", sum/(1024*float64(len(devices)))), fmt.Sprintf("%.0f", top/1024))
}

// prettyTimeSpan is PrettyTimeSpan.ToString: "2 days 3 hours 4 minutes", the days and hours only when reached.
func prettyTimeSpan(d time.Duration) string {
	plural := func(n int, one, many string) string {
		if n == 1 {
			return messages.Format(messages.Get(one), n)
		}
		return messages.Format(messages.Get(many), n)
	}
	days, hours, minutes := int(d.Hours())/24, int(d.Hours())%24, int(d.Minutes())%60
	out := plural(minutes, "TIME_MINUTE", "TIME_MINUTES")
	if d >= time.Hour {
		out = plural(hours, "TIME_HOUR", "TIME_HOURS") + " " + out
	}
	if d >= 24*time.Hour {
		out = plural(days, "TIME_DAY", "TIME_DAYS") + " " + out
	}
	return out
}
