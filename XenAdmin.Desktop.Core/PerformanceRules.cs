using System;
using System.Collections.Generic;
using System.Globalization;
using System.IO;
using System.Linq;
using System.Xml;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// The Performance tab's data (GraphList.LoadDefaultGraphs, ArchiveMaintainer, DataSet): default graphs, the
    /// rrd_updates download, units and the derived "used memory" series.
    /// </summary>
    public static class PerformanceRules
    {
        public enum Zoom { TenMinutes, OneHour, OneDay, OneWeek, OneMonth, OneYear }
        public enum Unit { None, Percent, Bytes, BytesPerSecond, MegaHertz }

        /// <summary>Hidden series are not drawn but take a colour, as Palette.GetColour does for the memory totals.</summary>
        public sealed record Series(string Id, string Name, bool Hidden = false);
        public sealed record Graph(string Title, List<Series> Series, Unit Unit, double FixedMax);
        public sealed record Point(DateTime Time, double Value);

        /// <summary>The Zoom menu of PerformancePage.</summary>
        public static string ZoomText(Zoom z) => z switch
        {
            Zoom.TenMinutes => "10 &Minutes", Zoom.OneHour => "1 &Hour", Zoom.OneDay => "1 &Day",
            Zoom.OneWeek => "1 &Week", Zoom.OneMonth => "1 M&onth", _ => "1 &Year"
        };

        /// <summary>The RRD archive behind each zoom (5 s, 1 min, 1 h, 1 day steps) and the period shown.</summary>
        public static (int Step, TimeSpan Window) Span(Zoom z) => z switch
        {
            Zoom.TenMinutes => (5, TimeSpan.FromMinutes(10)),
            Zoom.OneHour => (60, TimeSpan.FromHours(1)),
            Zoom.OneDay => (3600, TimeSpan.FromDays(1)),
            Zoom.OneWeek => (3600, TimeSpan.FromDays(7)),
            Zoom.OneMonth => (86400, TimeSpan.FromDays(30)),
            _ => (86400, TimeSpan.FromDays(365))
        };

        /// <summary>PerformancePage.CheckMessage: lifecycle messages of the VM, or of the VMs resident on the host; newest first.</summary>
        public static List<Message> LifecycleEvents(IXenObject o)
        {
            var uuids = o is VM vm ? new HashSet<string> { vm.uuid }
                : o is Host h ? h.Connection.ResolveAll(h.resident_VMs).Select(v => v.uuid).ToHashSet() : new HashSet<string>();
            return o.Connection.Cache.Messages.Where(m => m.ShowOnGraphs() && m.cls == cls.VM && uuids.Contains(m.obj_uuid))
                .OrderByDescending(m => m.timestamp).ToList();
        }

        /// <summary>Images.GetIconFor(Message.MessageType) for the lifecycle messages.</summary>
        public static string EventIcon(Message m) => m.Type switch
        {
            Message.MessageType.VM_STARTED => "001_PowerOn_h32bit_16", Message.MessageType.VM_SHUTDOWN => "001_ShutDown_h32bit_16",
            Message.MessageType.VM_REBOOTED => "001_Reboot_h32bit_16", Message.MessageType.VM_SUSPENDED => "000_paused_h32bit_16",
            Message.MessageType.VM_RESUMED => "000_Resumed_h32bit_16", Message.MessageType.VM_CLONED => "clonevm_16",
            _ => "alert6_16"
        };

        /// <summary>GraphList.LoadDefaultGraphs: CPU, memory and network (hosts); plus disk for VMs.</summary>
        public static List<Graph> DefaultGraphs(IXenObject o)
        {
            var graphs = new List<Graph>();
            Series S(string prefix, string name) => new Series($"{prefix}:{name}", Helpers.GetFriendlyDataSourceName(name, o));
            if (o is Host host)
            {
                var p = $"host:{host.uuid}";
                graphs.Add(new Graph(Messages.GRAPHS_DEFAULT_NAME_CPU,
                    host.Connection.ResolveAll(host.host_CPUs).OrderBy(c => c.number).Select(c => S(p, $"cpu{c.number}")).ToList(), Unit.Percent, 0));
                // memory_free_kib is drawn as used memory (DataSet: total - free)
                graphs.Add(new Graph(Messages.GRAPHS_DEFAULT_NAME_MEMORY,
                    new List<Series> { new Series($"{p}:memory_total_kib", "", true), new Series($"{p}:memory_free_kib", Helpers.GetFriendlyDataSourceName("memory_used_kib", o)) },
                    Unit.Bytes, host.Connection.Resolve(host.metrics)?.memory_total ?? 0));
                graphs.Add(new Graph(Messages.GRAPHS_DEFAULT_NAME_NETWORK,
                    host.Connection.ResolveAll(host.PIFs).OrderBy(x => x.device)
                        .SelectMany(x => new[] { S(p, $"pif_{x.device}_tx"), S(p, $"pif_{x.device}_rx") }).ToList(), Unit.BytesPerSecond, 0));
            }
            else if (o is VM vm)
            {
                var p = $"vm:{vm.uuid}";
                graphs.Add(new Graph(Messages.GRAPHS_DEFAULT_NAME_CPU,
                    Enumerable.Range(0, (int)Math.Max(1, vm.VCPUs_at_startup)).Select(i => S(p, $"cpu{i}")).ToList(), Unit.Percent, 0));
                graphs.Add(new Graph(Messages.GRAPHS_DEFAULT_NAME_MEMORY,
                    new List<Series> { new Series($"{p}:memory", "", true), new Series($"{p}:memory_internal_free", Helpers.GetFriendlyDataSourceName("memory_internal_used", o)) },
                    Unit.Bytes, vm.Connection.Resolve(vm.metrics)?.memory_actual ?? vm.memory_dynamic_max));
                graphs.Add(new Graph(Messages.GRAPHS_DEFAULT_NAME_NETWORK,
                    vm.Connection.ResolveAll(vm.VIFs).OrderBy(x => x.device)
                        .SelectMany(x => new[] { S(p, $"vif_{x.device}_tx"), S(p, $"vif_{x.device}_rx") }).ToList(), Unit.BytesPerSecond, 0));
                graphs.Add(new Graph(Messages.GRAPHS_DEFAULT_NAME_DISK,
                    vm.Connection.ResolveAll(vm.VBDs).Where(b => !string.IsNullOrEmpty(b.device)).OrderBy(b => b.device)
                        .SelectMany(b => new[] { S(p, $"vbd_{b.device}_read"), S(p, $"vbd_{b.device}_write") }).ToList(), Unit.BytesPerSecond, 0));
            }
            return graphs;
        }

        /// <summary>The host that serves the object's RRDs (a VM's resident host, else the coordinator).</summary>
        private static Host RrdHost(IXenObject o) =>
            o as Host ?? (o is VM vm ? vm.Connection.Resolve(vm.resident_on) ?? Helpers.GetCoordinator(vm.Connection) : null);

        /// <summary>ArchiveMaintainer.UpdateUri for one zoom level.</summary>
        public static Uri UpdatesUri(IXenObject o, Zoom zoom)
        {
            var session = o.Connection?.Session?.opaque_ref;
            var host = RrdHost(o);
            if (session == null || host == null)
                return null;
            var (step, window) = Span(zoom);
            var serverNow = DateTime.UtcNow - o.Connection.ServerTimeOffset;
            var start = (long)(serverNow - window - DateTime.UnixEpoch).TotalSeconds;
            // host=true alone also returns every resident VM's series
            var filter = o is VM vm ? $"vm_uuid={vm.uuid}" : "host=true&vm_uuid=none";
            return new UriBuilder
            {
                Scheme = o.Connection.UriScheme, Host = host.address, Port = o.Connection.Port, Path = "rrd_updates",
                Query = $"session_id={Uri.EscapeDataString(session)}&start={start}&cf=AVERAGE&interval={step}&{filter}"
            }.Uri;
        }

        /// <summary>Parses an rrd_updates xport: legend "AVERAGE:vm:uuid:name" becomes series "vm:uuid:name", times are UTC.</summary>
        public static Dictionary<string, List<Point>> ParseUpdates(Stream xml)
        {
            var legend = new List<string>();
            var rows = new List<(long T, List<double> V)>();
            using (var reader = XmlReader.Create(xml, new XmlReaderSettings { DtdProcessing = DtdProcessing.Prohibit }))
            {
                (long T, List<double> V)? row = null;
                // ReadElementContentAsString already moves to the next node: only Read() when nothing was consumed
                while (!reader.EOF)
                {
                    if (reader.NodeType != XmlNodeType.Element)
                    {
                        reader.Read();
                        continue;
                    }
                    switch (reader.Name)
                    {
                        case "entry":
                            var e = reader.ReadElementContentAsString();
                            var i = e.IndexOf(':');
                            legend.Add(i >= 0 ? e.Substring(i + 1) : e);
                            break;
                        case "row":
                            if (row != null) rows.Add(row.Value);
                            row = (0, new List<double>());
                            reader.Read();
                            break;
                        case "t" when row != null:
                            row = (long.Parse(reader.ReadElementContentAsString(), CultureInfo.InvariantCulture), row.Value.V);
                            break;
                        case "v" when row != null:
                            row.Value.V.Add(double.TryParse(reader.ReadElementContentAsString(), NumberStyles.Float, CultureInfo.InvariantCulture, out var d) ? d : double.NaN);
                            break;
                        default:
                            reader.Read();
                            break;
                    }
                }
                if (row != null) rows.Add(row.Value);
            }
            var result = legend.Distinct().ToDictionary(l => l, _ => new List<Point>());
            foreach (var (t, values) in rows.OrderBy(r => r.T))
                for (var k = 0; k < Math.Min(values.Count, legend.Count); k++)
                    result[legend[k]].Add(new Point(DateTimeOffset.FromUnixTimeSeconds(t).UtcDateTime, values[k]));
            return result;
        }

        /// <summary>DataSet's unit factors by the data source's units ("(fraction)" is shown as a percentage, KiB as bytes...).</summary>
        public static double Factor(string units) => units switch
        {
            "(fraction)" => 100, "KiB" => Util.BINARY_KILO, "MiB/s" => Util.BINARY_MEGA,
            "s" => Util.DEC_GIGA, "ms" => Util.DEC_MEGA, "μs" => Util.DEC_KILO, _ => 1
        };

        /// <summary>Scales by the units and derives used memory (total - free) as DataSet does; NaN marks a gap.</summary>
        public static Dictionary<string, List<Point>> Apply(Dictionary<string, List<Point>> raw, IReadOnlyDictionary<string, string> units)
        {
            double F(string id)
            {
                var name = id.Substring(id.LastIndexOf(':') + 1);
                return units.TryGetValue(name, out var u) ? Factor(u) : name.StartsWith("cpu") ? 100 : 1; // cpuN is a fraction
            }
            // DataSet: NaN and infinities are gaps
            var scaled = raw.ToDictionary(kv => kv.Key, kv => kv.Value
                .Select(p => p with { Value = double.IsFinite(p.Value) ? p.Value * F(kv.Key) : double.NaN }).ToList());
            void Used(string totalName, string freeName)
            {
                foreach (var freeId in scaled.Keys.Where(k => k.EndsWith(":" + freeName)).ToList())
                {
                    var totalId = freeId.Substring(0, freeId.Length - freeName.Length) + totalName;
                    if (!scaled.TryGetValue(totalId, out var total)) continue;
                    var byTime = total.ToDictionary(p => p.Time, p => p.Value);
                    scaled[freeId] = scaled[freeId].Select(p => p with { Value = byTime.TryGetValue(p.Time, out var t) ? t - p.Value : double.NaN }).ToList();
                }
            }
            Used("memory_total_kib", "memory_free_kib");
            Used("memory", "memory_internal_free");
            return scaled;
        }

        /// <summary>Downloads one zoom level for the object and applies the units (call off the UI thread).</summary>
        public static Dictionary<string, List<Point>> Fetch(IXenObject o, Zoom zoom, IReadOnlyDictionary<string, string> units)
        {
            var uri = UpdatesUri(o, zoom);
            if (uri == null)
                return new Dictionary<string, List<Point>>();
            using var stream = HTTPHelper.GET(uri, o.Connection, true);
            return Apply(ParseUpdates(stream), units);
        }

        /// <summary>
        /// Data source name to units ("cpu0" to "(fraction)"), from get_data_sources; that call needs a running VM, so
        /// halted ones get the units of the default graphs' memory sources (cpuN is treated as a fraction anyway).
        /// shortcut: a halted VM shows no history (ArchiveMaintainer loads vm_rrds for it); add with the full RRD parser.
        /// </summary>
        public static Dictionary<string, string> Units(IXenObject o, Session session)
        {
            if (o is VM v && v.power_state != vm_power_state.Running)
                return new Dictionary<string, string> { ["memory"] = "B", ["memory_internal_free"] = "KiB" };
            var sources = o is VM vm ? VM.get_data_sources(session, vm.opaque_ref)
                : o is Host h ? Host.get_data_sources(session, h.opaque_ref) : new List<Data_source>();
            return sources.GroupBy(d => d.name_label).ToDictionary(g => g.Key, g => g.First().units);
        }

        /// <summary>The axis of a graph: its top value, the divisor for the labels and the unit label ("MB", "kB/s").</summary>
        public static (double Max, double Divisor, string Label) Scale(Unit unit, double fixedMax, double dataMax)
        {
            // DataPlot: auto ranges start at 1 and are rounded by DataRange.RoundToNearestPowerOf10; memory keeps its fixed max
            var max = fixedMax > 0 ? fixedMax : RoundToNearestPowerOf10(unit, double.IsNaN(dataMax) || dataMax < 1 ? 1 : dataMax);
            var (divisor, label) = AxisUnit(unit, max);
            return (max, divisor, label);
        }

        /// <summary>DataRange.RoundToNearestPowerOf10: powers of 10, and 10^n * 1024^k for bytes.</summary>
        public static double RoundToNearestPowerOf10(Unit unit, double value)
        {
            if (unit != Unit.Bytes && unit != Unit.BytesPerSecond)
            {
                var pow = 0;
                var m = value;
                if (m > 1)
                    while (m > 1) { m /= 10; pow++; }
                else if (m > 0 && m < 1)
                    while (m < 1) { m *= 10; pow--; }
                return Math.Pow(10, pow);
            }
            int pows = 1, tens = 0;
            var x = value / 1024;
            while (x > 1024) { x /= 1024; pows++; }
            while (x > 1) { x /= 10; tens++; }
            if (tens == 3) { tens = 0; pows++; }
            return Math.Pow(10, tens) * Math.Pow(1024, pows);
        }

        /// <summary>The axis label for a maximum: "%", or the byte unit that keeps it readable ("MB", "kB/s").</summary>
        public static (double Divisor, string Label) AxisUnit(Unit unit, double max)
        {
            if (unit == Unit.Percent) return (1, "%");
            if (unit == Unit.MegaHertz) return (1, "MHz");
            if (unit != Unit.Bytes && unit != Unit.BytesPerSecond) return (1, "");
            var suffix = unit == Unit.BytesPerSecond ? "/s" : "";
            foreach (var (d, l) in new[] { (Util.BINARY_GIGA, "GB"), (Util.BINARY_MEGA, "MB"), (Util.BINARY_KILO, "kB") })
                if (max >= d) return (d, l + suffix);
            return (1, "B" + suffix);
        }
    }
}
