package ops

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// graphObject is a server ("hosts") or VM ("vms") of the Performance tab: the gui_config kind ("host", "vm"), its data
// sources (none for a VM that is not running) and record_data_source.
func graphObject(s *xenapi.Session, kind, uuid string) (short string, sources []xenapi.DataSourceRecord, record func(string) error, err error) {
	switch kind {
	case "hosts":
		ref, err := xenapi.Host.GetByUUID(s, uuid)
		if err != nil {
			return "", nil, nil, xapi.Sanitize(err)
		}
		sources, _ = xenapi.Host.GetDataSources(s, ref)
		return "host", sources, func(n string) error { return xenapi.Host.RecordDataSource(s, ref, n) }, nil
	case "vms":
		ref, err := xenapi.VM.GetByUUID(s, uuid)
		if err != nil {
			return "", nil, nil, xapi.Sanitize(err)
		}
		sources, _ = xenapi.VM.GetDataSources(s, ref)
		return "vm", sources, func(n string) error { return xenapi.VM.RecordDataSource(s, ref, n) }, nil
	}
	return "", nil, nil, fmt.Errorf("no performance graphs for %s", kind)
}

// SaveGraphs is GraphList.SaveGraphs + SaveDataSourceStateAction: the graphs of a server or VM in the pool's
// gui_config, where the Windows app keeps them ("XenCenter.GraphLayout.<i>.<vm|host>.<uuid>" the comma-separated data
// sources, "XenCenter.GraphName.<i>..." the name; no graphs is the default ones), the colours given ("#rrggbb" by data
// source, kept as "XenCenter.DataSource.<vm|host>.<uuid>.<source>"), and recording on for the non-standard data
// sources a graph shows. base is the layout the change was made on (none: the defaults): a layout changed meanwhile
// (by either app) is not overwritten. Each key is written on its own, so other keys changed meanwhile are kept.
// shortcut: C# also forgets the archives of the non-standard sources on no graph (forget_data_source_archives); that
// deletes history, so it is left to the Windows app.
func SaveGraphs(c *xapi.Connection, kind, uuid string, base, graphs []rules.DesignedGraph, colours map[string]string) error {
	for _, g := range graphs { // Graph Details asks for a source; a graph of none (C#'s default of no devices) is kept
		for _, src := range g.Sources {
			if src == "" || strings.Contains(src, ",") {
				return fmt.Errorf("invalid data source %q", src)
			}
		}
	}
	values := map[string]string{}
	for src, hex := range colours {
		v, err := rules.ColourValue(hex)
		if err != nil {
			return err
		}
		values[src] = v
	}
	s := c.Session
	short, sources, record, err := graphObject(s, kind, uuid)
	if err != nil {
		return err
	}
	shown := map[string]bool{}
	for _, g := range graphs {
		for _, src := range g.Sources {
			shown[src] = true
		}
	}
	var enable []string // UpdateDataSources: a non-standard source on a graph is recorded
	for _, d := range sources {
		if shown[d.NameLabel] && !d.Standard && !d.Enabled {
			enable = append(enable, d.NameLabel)
		}
	}
	calls := []string{"pool.add_to_gui_config", "pool.remove_from_gui_config"}
	if len(enable) > 0 {
		calls = append(calls, strings.TrimSuffix(kind, "s")+".record_data_source")
	}
	if err := permitted(s, calls...); err != nil {
		return err
	}
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err)
	}
	var pool xenapi.PoolRef
	old := map[string]string{}
	for r, p := range pools {
		pool, old = r, p.GuiConfig
	}
	if pool == "" {
		return errors.New("the server has no pool")
	}
	if !slices.EqualFunc(rules.ParseGraphLayout(old, short, uuid), base, func(a, b rules.DesignedGraph) bool {
		return a.Name == b.Name && slices.Equal(a.Sources, b.Sources)
	}) {
		return errors.New("the graphs have changed meanwhile: look at them again before changing them")
	}
	for _, n := range enable {
		if err := record(n); err != nil {
			return xapi.Sanitize(err)
		}
	}
	want := rules.GraphLayoutConfig(old, short, uuid, graphs)
	for src, v := range values {
		want[rules.DataSourceColourKey(short, uuid, src)] = v
	}
	// every key this save changed is put back as it was when a write fails: a lost layout key would hide the graphs
	// after it in both apps (C# writes the whole gui_config at once); a key another client wrote is left alone
	var touched []string
	rollback := func(err error) error {
		for i := len(touched) - 1; i >= 0; i-- {
			k := touched[i]
			_ = xenapi.Pool.RemoveFromGuiConfig(s, pool, k)
			if v, ok := old[k]; ok {
				_ = xenapi.Pool.AddToGuiConfig(s, pool, k, v)
			}
		}
		return xapi.Sanitize(err)
	}
	for k := range old {
		if _, ok := want[k]; !ok {
			if err := xenapi.Pool.RemoveFromGuiConfig(s, pool, k); err != nil {
				return rollback(err)
			}
			touched = append(touched, k)
		}
	}
	for k, v := range want {
		ov, ok := old[k]
		if ok && ov == v {
			continue
		}
		if ok {
			if err := xenapi.Pool.RemoveFromGuiConfig(s, pool, k); err != nil {
				return rollback(err)
			}
			touched = append(touched, k)
		}
		if err := xenapi.Pool.AddToGuiConfig(s, pool, k, v); err != nil {
			return rollback(err)
		}
		if !ok {
			touched = append(touched, k)
		}
	}
	return nil
}

// EnableDataSource is EnableDataSourceAction (Graph Details' Enable): recording on for a data source; text is its
// friendly name for the error.
func EnableDataSource(c *xapi.Connection, kind, uuid, name, text string) error {
	s := c.Session
	if err := permitted(s, strings.TrimSuffix(kind, "s")+".record_data_source"); err != nil {
		return err
	}
	_, _, record, err := graphObject(s, kind, uuid)
	if err != nil {
		return err
	}
	if err := record(name); err != nil {
		return fmt.Errorf("%s: %w", messages.Format(messages.Get("ACTION_ENABLING_DATASOURCE_ERROR"), text), xapi.Sanitize(err))
	}
	return nil
}
