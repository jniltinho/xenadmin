package app

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/trust"
	"xcpng-center/internal/xapi"
)

// Zooms are PerformancePage's Zoom menu: the RRD step in seconds and the period shown (PerformanceRules.Span).
var Zooms = map[string]struct {
	Step   int
	Window time.Duration
}{
	"10m": {5, 10 * time.Minute}, "1h": {60, time.Hour}, "1d": {3600, 24 * time.Hour},
	"1w": {3600, 7 * 24 * time.Hour}, "1mo": {86400, 30 * 24 * time.Hour}, "1y": {86400, 365 * 24 * time.Hour},
}

// Graph is one of PerformancePage's graphs: its series and the axis (PerformanceRules.Scale).
type Graph struct {
	Title   string       `json:"title"`
	Max     float64      `json:"max"`     // the axis top, in the data's units
	Divisor float64      `json:"divisor"` // the labels are value / divisor
	Label   string       `json:"label"`   // "%", "MB", "kB/s"
	Series  []PerfSeries `json:"series"`
}

// PerfSeries is one line: its friendly name and values (null marks a gap) at the Performance.Times.
type PerfSeries struct {
	Name   string     `json:"name"`
	Values []*float64 `json:"values"`
}

// PerfEvent is a VM lifecycle event drawn on the graphs (PerformancePage.CheckMessage).
type PerfEvent struct {
	Time time.Time `json:"time"`
	Name string    `json:"name"` // VM_STARTED...
	Text string    `json:"text"`
}

// Performance is the Performance tab of a VM or a server for one zoom.
type Performance struct {
	Times  []int64     `json:"times"` // unix seconds
	Start  int64       `json:"start"`
	End    int64       `json:"end"`
	Graphs []Graph     `json:"graphs"`
	Events []PerfEvent `json:"events"`
}

type unitKind int

const (
	unitPercent unitKind = iota
	unitBytes
	unitBytesPerSecond
)

type graphDef struct {
	title    string
	unit     unitKind
	fixedMax float64
	series   []struct {
		id, name string
		hidden   bool
	}
}

// ReadPerformance is PerformanceRules.DefaultGraphs + Fetch: the default graphs (CPU, memory, network and, for VMs,
// disk) from the rrd_updates of the server that holds the RRDs, the lifecycle events of the period.
// shortcut: default graphs only (no graph editor); a halted VM shows no history, as the C# port.
func ReadPerformance(c *xapi.Connection, policy *trust.Policy, kind, uuid, zoom string) (Performance, error) {
	z, ok := Zooms[zoom]
	if !ok {
		return Performance{}, fmt.Errorf("unknown zoom %q", zoom)
	}
	s := c.Session
	var defs []graphDef
	var rrdHost xenapi.HostRef
	var filter string
	units := map[string]string{}
	events := map[string]bool{} // VM uuids whose lifecycle messages are drawn
	add := func(d *graphDef, id, name string, hidden bool) {
		d.series = append(d.series, struct {
			id, name string
			hidden   bool
		}{id, name, hidden})
	}
	switch kind {
	case "hosts":
		ref, err := xenapi.Host.GetByUUID(s, uuid)
		if err != nil {
			return Performance{}, xapi.Sanitize(err)
		}
		h, err := xenapi.Host.GetRecord(s, ref)
		if err != nil {
			return Performance{}, xapi.Sanitize(err)
		}
		rrdHost, filter = ref, "host=true&vm_uuid=none"
		p := "host:" + uuid
		cpu := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_CPU"), unit: unitPercent}
		var nums []int
		for _, cr := range h.HostCPUs {
			if n, err := xenapi.HostCPU.GetNumber(s, cr); err == nil {
				nums = append(nums, n)
			}
		}
		sort.Ints(nums)
		for _, n := range nums {
			add(&cpu, fmt.Sprintf("%s:cpu%d", p, n), friendly(s, fmt.Sprintf("cpu%d", n), nil), false)
		}
		total, _ := xenapi.HostMetrics.GetMemoryTotal(s, h.Metrics)
		mem := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_MEMORY"), unit: unitBytes, fixedMax: float64(total)}
		add(&mem, p+":memory_total_kib", "", true)
		add(&mem, p+":memory_free_kib", messages.Friendly("Label-performance.memory_used_kib"), false)
		net := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_NETWORK"), unit: unitBytesPerSecond}
		var devices []string
		for _, pr := range h.PIFs {
			if d, err := xenapi.PIF.GetDevice(s, pr); err == nil {
				devices = append(devices, d)
			}
		}
		sort.Strings(devices)
		for _, d := range devices {
			for _, dir := range []string{"tx", "rx"} {
				name := fmt.Sprintf("pif_%s_%s", d, dir)
				add(&net, p+":"+name, friendly(s, name, nil), false)
			}
		}
		defs = []graphDef{cpu, mem, net}
		if sources, err := xenapi.Host.GetDataSources(s, ref); err == nil {
			for _, d := range sources {
				units[d.NameLabel] = d.Units
			}
		}
		for _, v := range h.ResidentVMs {
			if u, err := xenapi.VM.GetUUID(s, v); err == nil {
				events[u] = true
			}
		}
	case "vms":
		ref, err := xenapi.VM.GetByUUID(s, uuid)
		if err != nil {
			return Performance{}, xapi.Sanitize(err)
		}
		vm, err := xenapi.VM.GetRecord(s, ref)
		if err != nil {
			return Performance{}, xapi.Sanitize(err)
		}
		rrdHost, filter = vm.ResidentOn, "vm_uuid="+uuid
		if rrdHost == "" || rrdHost == nullRef { // a halted VM: the coordinator
			pools, _ := xenapi.Pool.GetAllRecords(s)
			for _, pl := range pools {
				rrdHost = pl.Master
			}
		}
		p := "vm:" + uuid
		cpu := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_CPU"), unit: unitPercent}
		for i := 0; i < max(1, vm.VCPUsAtStartup); i++ {
			add(&cpu, fmt.Sprintf("%s:cpu%d", p, i), friendly(s, fmt.Sprintf("cpu%d", i), nil), false)
		}
		memMax := float64(vm.MemoryDynamicMax)
		if actual, err := xenapi.VMMetrics.GetMemoryActual(s, vm.Metrics); err == nil && actual > 0 {
			memMax = float64(actual)
		}
		mem := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_MEMORY"), unit: unitBytes, fixedMax: memMax}
		add(&mem, p+":memory", "", true)
		add(&mem, p+":memory_internal_free", messages.Friendly("Label-performance.memory_internal_used"), false)
		net := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_NETWORK"), unit: unitBytesPerSecond}
		vifs := map[string]xenapi.NetworkRef{}
		var vifDevs []string
		for _, v := range vm.VIFs {
			if r, err := xenapi.VIF.GetRecord(s, v); err == nil {
				vifs[r.Device], vifDevs = r.Network, append(vifDevs, r.Device)
			}
		}
		sort.Strings(vifDevs)
		for _, d := range vifDevs {
			for _, dir := range []string{"tx", "rx"} {
				name := fmt.Sprintf("vif_%s_%s", d, dir)
				add(&net, p+":"+name, friendly(s, name, vifs), false)
			}
		}
		disk := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_DISK"), unit: unitBytesPerSecond}
		vbds := map[string]string{} // device -> userdevice
		var devs []string
		for _, b := range vm.VBDs {
			if r, err := xenapi.VBD.GetRecord(s, b); err == nil && r.Device != "" {
				vbds[r.Device], devs = r.Userdevice, append(devs, r.Device)
			}
		}
		sort.Strings(devs)
		for _, d := range devs {
			for _, dir := range []string{"read", "write"} {
				add(&disk, fmt.Sprintf("%s:vbd_%s_%s", p, d, dir), messages.Format(messages.Friendly("Label-performance.vbd_"+dir), vbds[d]), false)
			}
		}
		defs = []graphDef{cpu, mem, net, disk}
		if vm.PowerState == xenapi.VMPowerStateRunning {
			if sources, err := xenapi.VM.GetDataSources(s, ref); err == nil {
				for _, d := range sources {
					units[d.NameLabel] = d.Units
				}
			}
		} else {
			units["memory"], units["memory_internal_free"] = "B", "KiB"
		}
		events[uuid] = true
	default:
		return Performance{}, fmt.Errorf("no performance data for %s", kind)
	}

	address, err := xenapi.Host.GetAddress(s, rrdHost)
	if err != nil {
		return Performance{}, xapi.Sanitize(err)
	}
	// the server's clock, as ServerTimeOffset: a client clock ahead or behind would ask for the wrong period
	now, err := xenapi.Host.GetServertime(s, rrdHost)
	if err != nil {
		now = time.Now()
	}
	start := now.Add(-z.Window)
	raw, err := fetchUpdates(s, policy, address, c.Port, fmt.Sprintf("start=%d&cf=AVERAGE&interval=%d&%s", start.Unix(), z.Step, filter))
	if err != nil {
		return Performance{}, xapi.Sanitize(err)
	}
	out := Performance{Times: raw.times, Start: start.Unix(), End: now.Unix(), Events: lifecycle(s, events, start)}
	for _, d := range defs {
		g := Graph{Title: d.title, Series: []PerfSeries{}}
		dataMax := math.NaN()
		values := map[string][]float64{}
		for _, sr := range d.series {
			values[sr.id] = scaled(raw, sr.id, units)
		}
		// DataSet: used memory is total - free
		for _, sr := range d.series {
			if strings.HasSuffix(sr.id, ":memory_free_kib") || strings.HasSuffix(sr.id, ":memory_internal_free") {
				total := values[strings.TrimSuffix(strings.TrimSuffix(sr.id, "memory_free_kib"), "memory_internal_free")+
					map[bool]string{true: "memory_total_kib", false: "memory"}[strings.HasSuffix(sr.id, "_kib")]]
				free := values[sr.id]
				for i := range free {
					if i < len(total) {
						free[i] = total[i] - free[i]
					} else {
						free[i] = math.NaN()
					}
				}
			}
		}
		for _, sr := range d.series {
			if sr.hidden {
				continue
			}
			ps := PerfSeries{Name: sr.name, Values: make([]*float64, len(raw.times))}
			for i, v := range values[sr.id] {
				if !math.IsNaN(v) && !math.IsInf(v, 0) {
					v := v
					ps.Values[i] = &v
					if math.IsNaN(dataMax) || v > dataMax {
						dataMax = v
					}
				}
			}
			g.Series = append(g.Series, ps)
		}
		g.Max, g.Divisor, g.Label = axis(d.unit, d.fixedMax, dataMax)
		out.Graphs = append(out.Graphs, g)
	}
	return out, nil
}

type updates struct {
	times  []int64
	series map[string][]float64 // "vm:uuid:name" -> values at times
}

// fetchUpdates is ArchiveMaintainer's GET of rrd_updates on the server that holds the RRDs; the session goes in a
// cookie (never in the URL), the server's certificate through the member trust.
func fetchUpdates(s *xenapi.Session, policy *trust.Policy, address string, port int, query string) (updates, error) {
	key := xapi.HostKey(address, port)
	client := &http.Client{Timeout: 60 * time.Second, Transport: &http.Transport{TLSClientConfig: xapi.MemberTLS(policy, key, xapi.VouchedFingerprints(s, address)),
		Proxy: xapi.Proxy, DialContext: xapi.Dialer().DialContext, TLSHandshakeTimeout: xapi.HandshakeTimeout()},
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	req, err := http.NewRequest(http.MethodGet, "https://"+key+"/rrd_updates?"+query, nil)
	if err != nil {
		return updates{}, err
	}
	req.AddCookie(&http.Cookie{Name: "session_id", Value: string(s.Ref())})
	resp, err := client.Do(req)
	if err != nil {
		return updates{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return updates{}, fmt.Errorf("rrd_updates: %s", resp.Status)
	}
	return parseUpdates(io.LimitReader(resp.Body, 64<<20))
}

// parseUpdates is PerformanceRules.ParseUpdates: the legend "AVERAGE:vm:uuid:name" names the columns, rows by time.
func parseUpdates(r io.Reader) (updates, error) {
	var doc struct {
		Legend []string `xml:"meta>legend>entry"`
		Rows   []struct {
			T string   `xml:"t"`
			V []string `xml:"v"`
		} `xml:"data>row"`
	}
	if err := xml.NewDecoder(r).Decode(&doc); err != nil {
		return updates{}, err
	}
	u := updates{series: map[string][]float64{}}
	names := make([]string, len(doc.Legend))
	for i, e := range doc.Legend {
		if j := strings.Index(e, ":"); j >= 0 {
			e = e[j+1:]
		}
		names[i] = e
	}
	type row struct {
		t int64
		v []string
	}
	rows := make([]row, 0, len(doc.Rows))
	for _, r := range doc.Rows {
		t, err := strconv.ParseInt(r.T, 10, 64)
		if err == nil {
			rows = append(rows, row{t, r.V})
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].t < rows[j].t })
	for _, r := range rows {
		u.times = append(u.times, r.t)
	}
	for k, name := range names {
		vals := make([]float64, len(rows))
		for i, r := range rows {
			vals[i] = math.NaN()
			if k < len(r.v) {
				if f, err := strconv.ParseFloat(r.v[k], 64); err == nil {
					vals[i] = f
				}
			}
		}
		u.series[name] = vals
	}
	return u, nil
}

// scaled is PerformanceRules.Apply for one series: the units' factor ("(fraction)" as %, KiB as bytes; cpuN is a
// fraction), NaN for gaps and for a series the server does not have.
func scaled(u updates, id string, units map[string]string) []float64 {
	src, ok := u.series[id]
	out := make([]float64, len(u.times))
	name := id[strings.LastIndex(id, ":")+1:]
	f := 1.0
	switch unit, known := units[name]; {
	case known && unit == "(fraction)":
		f = 100
	case known && unit == "KiB":
		f = 1024
	case known && unit == "MiB/s":
		f = 1 << 20
	case !known && strings.HasPrefix(name, "cpu"):
		f = 100
	}
	for i := range out {
		out[i] = math.NaN()
		if ok && i < len(src) && !math.IsNaN(src[i]) && !math.IsInf(src[i], 0) {
			out[i] = src[i] * f
		}
	}
	return out
}

// axis is PerformanceRules.Scale: a fixed top (memory) or the data's, rounded by RoundToNearestPowerOf10; the unit
// label that keeps it readable.
func axis(unit unitKind, fixedMax, dataMax float64) (top, divisor float64, label string) {
	top = fixedMax
	if top <= 0 {
		v := dataMax
		if math.IsNaN(v) || v < 1 {
			v = 1
		}
		top = roundPow10(unit, v)
	}
	if unit == unitPercent {
		return top, 1, "%"
	}
	suffix := map[bool]string{true: "/s", false: ""}[unit == unitBytesPerSecond]
	for _, u := range []struct {
		d float64
		l string
	}{{1 << 30, "GB"}, {1 << 20, "MB"}, {1 << 10, "kB"}} {
		if top >= u.d {
			return top, u.d, u.l + suffix
		}
	}
	return top, 1, "B" + suffix
}

// roundPow10 is DataRange.RoundToNearestPowerOf10: powers of 10, and 10^n * 1024^k for bytes.
func roundPow10(unit unitKind, v float64) float64 {
	if unit == unitPercent {
		pow := 0
		for m := v; m > 1; m /= 10 {
			pow++
		}
		return math.Pow(10, float64(pow))
	}
	pows, tens := 1, 0
	x := v / 1024
	for x > 1024 {
		x /= 1024
		pows++
	}
	for x > 1 {
		x /= 10
		tens++
	}
	if tens == 3 {
		tens, pows = 0, pows+1
	}
	return math.Pow(10, float64(tens)) * math.Pow(1024, float64(pows))
}

var (
	cpuRe = regexp.MustCompile(`^cpu([0-9]+)$`)
	vifRe = regexp.MustCompile(`^vif_([0-9]+)_(tx|rx)$`)
	pifRe = regexp.MustCompile(`^pif_eth([0-9]+)_(tx|rx)$`)
)

// friendly is Helpers.GetFriendlyDataSourceName for the default graphs' sources (CPU, VIF by its network, NIC).
func friendly(s *xenapi.Session, name string, vifNetworks map[string]xenapi.NetworkRef) string {
	if m := cpuRe.FindStringSubmatch(name); m != nil {
		return messages.Format(messages.Friendly("Label-performance.cpu"), m[1])
	}
	if m := vifRe.FindStringSubmatch(name); m != nil {
		if n, err := xenapi.Network.GetNameLabel(s, vifNetworks[m[1]]); err == nil {
			return messages.Format(messages.Friendly("Label-performance.vif_"+m[2]), n)
		}
		return name
	}
	if m := pifRe.FindStringSubmatch(name); m != nil {
		return messages.Format(messages.Friendly("Label-performance.nic_"+m[2]), m[1])
	}
	return name
}

// lifecycle is PerformanceRules.LifecycleEvents: the lifecycle messages (Message.ShowOnGraphs) of these VMs since
// start, newest first.
func lifecycle(s *xenapi.Session, vms map[string]bool, start time.Time) []PerfEvent {
	out := []PerfEvent{}
	msgs, err := xenapi.Message.GetAllRecords(s)
	if err != nil {
		return out
	}
	for _, m := range msgs {
		if m.Cls != xenapi.ClsVM || !vms[m.ObjUUID] || !graphsOnly[m.Name] || m.Name == "HA_POOL_OVERCOMMITTED" || m.Timestamp.Before(start) {
			continue
		}
		out = append(out, PerfEvent{Time: m.Timestamp, Name: m.Name, Text: messages.Friendly("Message.name-" + strings.ToLower(m.Name))})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Time.After(out[j].Time) })
	return out
}
