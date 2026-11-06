package app

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
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
	Sources []string     `json:"sources"` // its data source names, for the graph editor (DesignedGraph)
}

// PerfSeries is one line: its friendly name, colour and values (null marks a gap) at the Performance.Times.
type PerfSeries struct {
	Name   string     `json:"name"`
	Colour string     `json:"colour"` // "#rrggbb"
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
	// the graph editor: GraphList.ShowingDefaultGraphs, PerformancePage.UpdateActionButtons' tooltips of New and Edit
	// Graph on a VM that is not running (they are disabled then), GraphDetailsDialog's suggested name
	Default     bool   `json:"default"`
	NewBlocker  string `json:"new_blocker,omitempty"`
	EditBlocker string `json:"edit_blocker,omitempty"`
	NewName     string `json:"new_name"`
}

type unitKind int

const (
	unitPercent unitKind = iota
	unitBytes
	unitBytesPerSecond
	unitOther // a count, seconds...: powers of ten, labelled with the data source's units
)

type graphDef struct {
	title    string
	unit     unitKind
	label    string // unitOther's
	fixedMax float64
	sources  []string
	series   []struct {
		id, name string
		hidden   bool
	}
}

// graphUnit is DataSet's CustomYRange by the data source's units: the axis of a designed graph is its first shown
// source's. shortcut: one axis per graph (C# scales each data set to its own range, so mixed units differ) and s/ms/μs
// in their own units (C# converts them to ns); give each series its own range if users mix units.
func graphUnit(units string, name string) (unitKind, string) {
	switch units {
	case "(fraction)", "percent":
		return unitPercent, ""
	case "B", "KiB":
		return unitBytes, ""
	case "B/s", "MiB/s":
		return unitBytesPerSecond, ""
	case "":
		if strings.HasPrefix(name, "cpu") { // a halted VM has no units: cpuN is a fraction
			return unitPercent, ""
		}
	}
	return unitOther, units
}

// sourceFinder is what Helpers.GetFriendlyDataSourceName looks up for rules.FriendlyDataSource: a VM's VIFs and VBDs,
// the PIFs of the server (a VM's: where it runs), the SRs; each read once, when first needed.
func sourceFinder(s *xenapi.Session, host xenapi.HostRef, vm *xenapi.VMRecord) func(what string, ids ...string) string {
	var vifs []xenapi.VIFRecord
	var vbds []xenapi.VBDRecord
	var pifs []xenapi.PIFRecord
	var srs []xenapi.SRRecord
	loaded := map[string]bool{}
	load := func(what string) {
		if loaded[what] {
			return
		}
		loaded[what] = true
		switch what {
		case "vif":
			for _, r := range vm.VIFs {
				if v, err := xenapi.VIF.GetRecord(s, r); err == nil {
					vifs = append(vifs, v)
				}
			}
		case "vbd":
			for _, r := range vm.VBDs {
				if v, err := xenapi.VBD.GetRecord(s, r); err == nil {
					vbds = append(vbds, v)
				}
			}
		case "pif":
			refs, _ := xenapi.Host.GetPIFs(s, host)
			for _, r := range refs {
				if p, err := xenapi.PIF.GetRecord(s, r); err == nil {
					pifs = append(pifs, p)
				}
			}
		case "sr":
			all, _ := xenapi.SR.GetAllRecords(s)
			for _, r := range all {
				srs = append(srs, r)
			}
		}
	}
	network := func(n xenapi.NetworkRef) string { name, _ := xenapi.Network.GetNameLabel(s, n); return name }
	return func(what string, ids ...string) string {
		switch {
		case (what == "vif" || what == "vbd") && vm == nil, (what == "pif" || what == "vlan") && (host == "" || host == nullRef):
			return ""
		}
		switch what {
		case "vif": // FindNetworkOfVIF
			load("vif")
			for _, v := range vifs {
				if v.Device == ids[0] {
					return network(v.Network)
				}
			}
		case "vbd": // FindVBD
			load("vbd")
			for _, v := range vbds {
				if v.Device == ids[0] {
					return v.Userdevice
				}
			}
		case "pif", "vlan": // FindNetworkOfPIF (a physical NIC), FindVlan
			load("pif")
			for _, p := range pifs {
				if p.Device == ids[0] && (what == "pif" && p.Physical || what == "vlan" && len(ids) > 1 && strconv.Itoa(p.VLAN) == ids[1]) {
					return network(p.Network)
				}
			}
		case "sr": // FindSr
			load("sr")
			for _, r := range srs {
				if strings.HasPrefix(r.UUID, ids[0]) {
					return r.NameLabel
				}
			}
		}
		return ""
	}
}

// ReadPerformance is PerformanceRules.DefaultGraphs + Fetch: the graphs saved in the pool's gui_config by either app
// (GraphList.LoadGraphs), else the default ones (CPU, memory, network and, for VMs, disk), from the rrd_updates of the
// server that holds the RRDs, the lifecycle events of the period.
// shortcut: a halted VM shows no history, as the C# port.
func ReadPerformance(c *xapi.Connection, policy *trust.Policy, kind, uuid, zoom string) (Performance, error) {
	z, ok := Zooms[zoom]
	if !ok {
		return Performance{}, fmt.Errorf("unknown zoom %q", zoom)
	}
	s := c.Session
	var defs []graphDef
	var rrdHost xenapi.HostRef
	var filter, short, editBlocker, newBlocker string
	var find func(what string, ids ...string) string
	var memMax float64 // the memory graph's top (DataSet.GetMemoryMax)
	units := map[string]string{}
	events := map[string]bool{} // VM uuids whose lifecycle messages are drawn
	add := func(d *graphDef, source string) {
		d.sources = append(d.sources, source)
		d.series = append(d.series, struct {
			id, name string
			hidden   bool
		}{short + ":" + uuid + ":" + source, rules.FriendlyDataSource(source, find), rules.HiddenDataSource(source)})
	}
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return Performance{}, xapi.Sanitize(err)
	}
	gui := map[string]string{}
	var coordinator xenapi.HostRef
	for _, pl := range pools {
		gui, coordinator = pl.GuiConfig, pl.Master
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
		rrdHost, filter, short, find = ref, "host=true&vm_uuid=none", "host", sourceFinder(s, ref, nil)
		cpu := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_CPU"), unit: unitPercent}
		var nums []int
		for _, cr := range h.HostCPUs {
			if n, err := xenapi.HostCPU.GetNumber(s, cr); err == nil {
				nums = append(nums, n)
			}
		}
		sort.Ints(nums)
		for _, n := range nums {
			add(&cpu, fmt.Sprintf("cpu%d", n))
		}
		total, _ := xenapi.HostMetrics.GetMemoryTotal(s, h.Metrics)
		memMax = float64(total)
		mem := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_MEMORY"), unit: unitBytes, fixedMax: memMax}
		add(&mem, "memory_total_kib")
		add(&mem, "memory_free_kib")
		net := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_NETWORK"), unit: unitBytesPerSecond}
		var devices []string
		for _, pr := range h.PIFs {
			if d, err := xenapi.PIF.GetDevice(s, pr); err == nil {
				devices = append(devices, d)
			}
		}
		sort.Strings(devices)
		for _, d := range devices {
			add(&net, fmt.Sprintf("pif_%s_tx", d))
			add(&net, fmt.Sprintf("pif_%s_rx", d))
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
		rrdHost, filter, short, find = vm.ResidentOn, "vm_uuid="+uuid, "vm", sourceFinder(s, vm.ResidentOn, &vm)
		if rrdHost == "" || rrdHost == nullRef { // a halted VM: the coordinator
			rrdHost = coordinator
		}
		if len(vm.CurrentOperations) > 0 || vm.PowerState != xenapi.VMPowerStateRunning {
			newBlocker = messages.Format(messages.Get("GRAPHS_CANNOT_ADD_VM_HALTED"), vm.NameLabel)
			editBlocker = messages.Format(messages.Get("GRAPHS_CANNOT_EDIT_VM_HALTED"), vm.NameLabel)
		}
		cpu := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_CPU"), unit: unitPercent}
		for i := 0; i < max(1, vm.VCPUsAtStartup); i++ {
			add(&cpu, fmt.Sprintf("cpu%d", i))
		}
		memMax = float64(vm.MemoryDynamicMax)
		if actual, err := xenapi.VMMetrics.GetMemoryActual(s, vm.Metrics); err == nil && actual > 0 {
			memMax = float64(actual)
		}
		mem := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_MEMORY"), unit: unitBytes, fixedMax: memMax}
		add(&mem, "memory")
		add(&mem, "memory_internal_free")
		net := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_NETWORK"), unit: unitBytesPerSecond}
		var vifDevs []string
		for _, v := range vm.VIFs {
			if d, err := xenapi.VIF.GetDevice(s, v); err == nil {
				vifDevs = append(vifDevs, d)
			}
		}
		sort.Strings(vifDevs)
		for _, d := range vifDevs {
			add(&net, fmt.Sprintf("vif_%s_tx", d))
			add(&net, fmt.Sprintf("vif_%s_rx", d))
		}
		disk := graphDef{title: messages.Get("GRAPHS_DEFAULT_NAME_DISK"), unit: unitBytesPerSecond}
		var devs []string
		for _, b := range vm.VBDs {
			if d, err := xenapi.VBD.GetDevice(s, b); err == nil && d != "" {
				devs = append(devs, d)
			}
		}
		sort.Strings(devs)
		for _, d := range devs {
			add(&disk, fmt.Sprintf("vbd_%s_read", d))
			add(&disk, fmt.Sprintf("vbd_%s_write", d))
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

	out := Performance{Default: true, NewBlocker: newBlocker, EditBlocker: editBlocker}
	if layout := rules.ParseGraphLayout(gui, short, uuid); layout != nil {
		out.Default, defs = false, nil
		for _, g := range layout {
			d := graphDef{title: g.Name, unit: unitOther}
			for _, src := range g.Sources {
				add(&d, src)
			}
			for _, sr := range d.series {
				if !sr.hidden {
					name := sr.id[strings.LastIndex(sr.id, ":")+1:]
					d.unit, d.label = graphUnit(units[name], name)
					if name == "memory_free_kib" || name == "memory_internal_free" { // DataSet: the used memory
						d.unit, d.fixedMax = unitBytes, memMax
					}
					break
				}
			}
			defs = append(defs, d)
		}
	}
	taken := map[string]bool{}
	for _, d := range defs {
		taken[d.title] = true
	}
	out.NewName = rules.UniqueName(messages.Get("GRAPH_NAME"), taken)

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
	out.Times, out.Start, out.End, out.Events = raw.times, start.Unix(), now.Unix(), lifecycle(s, events, start)
	prefix := short + ":" + uuid + ":"
	for _, d := range defs {
		g := Graph{Title: d.title, Series: []PerfSeries{}, Sources: append([]string{}, d.sources...)}
		dataMax := math.NaN()
		for _, sr := range d.series {
			if sr.hidden {
				continue
			}
			values := scaled(raw, sr.id, units)
			// DataSet.AddPoint: used memory is total - free
			if total, ok := map[string]string{prefix + "memory_free_kib": "memory_total_kib", prefix + "memory_internal_free": "memory"}[sr.id]; ok {
				totals := scaled(raw, prefix+total, units)
				for i := range values {
					values[i] = totals[i] - values[i]
				}
			}
			source := sr.id[len(prefix):]
			colour, ok := rules.ParseColour(gui[rules.DataSourceColourKey(short, uuid, source)])
			if !ok {
				colour = rules.GraphColours[len(g.Series)%len(rules.GraphColours)]
			}
			ps := PerfSeries{Name: sr.name, Colour: colour, Values: make([]*float64, len(raw.times))}
			for i, v := range values {
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
		if d.unit == unitOther {
			g.Label = d.label
		}
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
	if unit == unitPercent || unit == unitOther {
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
	if unit == unitPercent || unit == unitOther {
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

// DataSource is a row of GraphDetailsDialog (DataSourceItem): the data source, its friendly name, Type column,
// recording state, whether it is hidden (no units) and the colour saved in gui_config ("" when none).
type DataSource struct {
	Name        string `json:"name"`
	Text        string `json:"text"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Enabled     bool   `json:"enabled"`
	Standard    bool   `json:"standard"`
	Hidden      bool   `json:"hidden"`
	Colour      string `json:"colour,omitempty"`
}

// GraphDataSources is GetDataSourcesAction + DataSourceItemList.BuildList for a server ("hosts") or a VM ("vms").
func GraphDataSources(c *xapi.Connection, kind, uuid string) ([]DataSource, error) {
	s := c.Session
	var sources []xenapi.DataSourceRecord
	var find func(what string, ids ...string) string
	var short string
	switch kind {
	case "hosts":
		ref, err := xenapi.Host.GetByUUID(s, uuid)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		if sources, err = xenapi.Host.GetDataSources(s, ref); err != nil {
			return nil, xapi.Sanitize(err)
		}
		short, find = "host", sourceFinder(s, ref, nil)
	case "vms":
		ref, err := xenapi.VM.GetByUUID(s, uuid)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		vm, err := xenapi.VM.GetRecord(s, ref)
		if err != nil {
			return nil, xapi.Sanitize(err)
		}
		if sources, err = xenapi.VM.GetDataSources(s, ref); err != nil {
			return nil, xapi.Sanitize(err)
		}
		short, find = "vm", sourceFinder(s, vm.ResidentOn, &vm)
	default:
		return nil, fmt.Errorf("no performance data for %s", kind)
	}
	gui := map[string]string{}
	if pools, err := xenapi.Pool.GetAllRecords(s); err == nil {
		for _, p := range pools {
			gui = p.GuiConfig
		}
	}
	byName := map[string]xenapi.DataSourceRecord{}
	var names []string
	for _, d := range sources {
		byName[d.NameLabel] = d
		names = append(names, d.NameLabel)
	}
	out := []DataSource{}
	for _, n := range rules.OfferedDataSources(names) {
		d := byName[n]
		colour, _ := rules.ParseColour(gui[rules.DataSourceColourKey(short, uuid, n)])
		out = append(out, DataSource{Name: n, Text: rules.FriendlyDataSource(n, find), Description: d.NameDescription, Type: rules.DataSourceCategory(n),
			Enabled: d.Enabled, Standard: d.Standard, Hidden: d.Units == "" || d.Units == "unknown", Colour: colour})
	}
	return out, nil
}
