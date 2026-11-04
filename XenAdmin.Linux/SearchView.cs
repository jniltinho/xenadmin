using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using Avalonia.Threading;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAdmin.XenSearch;
using XenAPI;

namespace XenAdmin.Linux
{
    /// <summary>
    /// SearchPage default "Overview": pool → hosts → VMs → custom templates, with the CPU, memory,
    /// disk, network, address and uptime columns computed by XenModel (PropertyAccessorHelper + MetricUpdater).
    /// </summary>
    public class SearchView : DockPanel, ILiveView
    {
        private static MetricUpdater _metrics;
        private readonly IReadOnlyList<IXenConnection> _connections;
        private readonly Grid _grid = new Grid();
        private readonly HashSet<string> _collapsed = new HashSet<string>();
        private static readonly string[] Columns = { "Name", Messages.OVERVIEW_CPU_USAGE, Messages.OVERVIEW_MEMORY_USAGE, "Disks\n(avg / max KBs)", "Network\n(avg / max KBs)", "Address", "Uptime" };
        private static readonly double[] Widths = { 250, 115, 125, 120, 120, 120, 170 };

        public static string Title(IXenConnection connection) =>
            connection == null ? "XCP-ng Center: Overview" : $"{Helpers.GetName(Helpers.GetPoolOfOne(connection))}: Overview";

        public SearchView(IReadOnlyList<IXenConnection> connections)
        {
            _connections = connections;
            Background = Brushes.White;
            if (_metrics == null)
            {
                _metrics = new MetricUpdater();
                PropertyAccessorHelper.MetricUpdater = _metrics;
            }

            var buttons = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            foreach (var t in new[] { "New Search", "Edit Search", "Saved Searches  ▾", "Export...", "Import..." })
                buttons.Children.Add(new Button { Content = t, IsEnabled = false });
            var top = new Border { BorderBrush = Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1, 1, 1, 0), Padding = new Thickness(10, 8), Child = buttons };
            var columns = new Button { Content = "Show Columns  ▾", IsEnabled = false, HorizontalAlignment = HorizontalAlignment.Right, Margin = new Thickness(0, 6, 4, 6) };
            var body = new DockPanel();
            DockPanel.SetDock(columns, Dock.Top);
            body.Children.Add(columns);
            body.Children.Add(new ScrollViewer { Content = _grid, HorizontalScrollBarVisibility = Avalonia.Controls.Primitives.ScrollBarVisibility.Auto });
            var frame = new Border { BorderBrush = Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1), Padding = new Thickness(4, 0, 4, 4), Child = body };
            var content = new DockPanel { Margin = new Thickness(12) };
            DockPanel.SetDock(top, Dock.Top);
            content.Children.Add(top);
            content.Children.Add(frame);
            Children.Add(content);

            foreach (var w in Widths)
                _grid.ColumnDefinitions.Add(new ColumnDefinition(w, GridUnitType.Pixel));

            EventHandler onMetrics = (_, _) => Dispatcher.UIThread.Post(Refresh);
            AttachedToVisualTree += (_, _) =>
            {
                _metrics.MetricsUpdated += onMetrics;
                _metrics.SetXenObjects(Rows().Select(r => r.Object).Where(o => o is Host || o is VM).ToArray());
                _metrics.Start();
                _metrics.Prod();
            };
            DetachedFromVisualTree += (_, _) =>
            {
                _metrics.MetricsUpdated -= onMetrics;
                _metrics.Pause();
            };
            Refresh();
        }

        private sealed class Row
        {
            public IXenObject Object;
            public int Depth;
        }

        private IEnumerable<Row> Rows()
        {
            foreach (var c in _connections.Where(c => c.IsConnected && c.CacheIsPopulated))
            {
                var pool = Helpers.GetPool(c);
                var depth = 0;
                if (pool != null)
                    yield return new Row { Object = pool, Depth = depth++ };
                var vms = c.Cache.VMs.Where(ResourceTree.IsRealVm).OrderBy(v => v.Name(), StringComparer.CurrentCultureIgnoreCase).ToList();
                foreach (var host in c.Cache.Hosts.OrderBy(h => h.Name(), StringComparer.CurrentCultureIgnoreCase))
                {
                    yield return new Row { Object = host, Depth = depth };
                    foreach (var vm in vms.Where(v => v.Home()?.opaque_ref == host.opaque_ref))
                        yield return new Row { Object = vm, Depth = depth + 1 };
                }
                foreach (var vm in vms.Where(v => v.Home() == null))
                    yield return new Row { Object = vm, Depth = depth };
                foreach (var t in c.Cache.VMs.Where(ResourceTree.IsUserTemplate).OrderBy(v => v.Name(), StringComparer.CurrentCultureIgnoreCase))
                    yield return new Row { Object = t, Depth = depth };
            }
        }

        /// <summary>(text, bar 0..100 or -1) per column, as QueryPanel shows them.</summary>
        public static (string Text, int Bar)[] Values(IXenObject o)
        {
            const string none = "-";
            switch (o)
            {
                case Host h:
                    return new[]
                    {
                        (PropertyAccessorHelper.hostCpuUsageString(h), PropertyAccessorHelper.hostCpuUsageRank(h)),
                        (PropertyAccessorHelper.hostMemoryUsageString(h), PropertyAccessorHelper.hostMemoryUsageRank(h)),
                        (none, -1),
                        (PropertyAccessorHelper.hostNetworkUsageString(h), -1),
                        (Address(h), -1),
                        (h.Uptime()?.ToString() ?? none, -1)
                    };
                case VM vm when !vm.is_a_template:
                    var running = vm.power_state == vm_power_state.Running;
                    return new[]
                    {
                        (running ? PropertyAccessorHelper.vmCpuUsageString(vm) : none, running ? PropertyAccessorHelper.vmCpuUsageRank(vm) : -1),
                        (running ? PropertyAccessorHelper.vmMemoryUsageString(vm) : none, running ? PropertyAccessorHelper.vmMemoryUsageRank(vm) : -1),
                        (running ? PropertyAccessorHelper.vmDiskUsageString(vm) : none, -1),
                        (running ? PropertyAccessorHelper.vmNetworkUsageString(vm) : none, -1),
                        (Address(vm), -1),
                        (vm.RunningTime()?.ToString() ?? none, -1)
                    };
                default:
                    return Enumerable.Repeat((none, -1), 4).Append((o is Pool ? none : "", -1)).Append((o is Pool ? none : "", -1)).ToArray();
            }
        }

        /// <summary>The search "ip_address" property (distinct, VIF/PIF addresses), as the WinForms Address column.</summary>
        private static string Address(IXenObject o)
        {
            var text = PropertyAccessors.Get(PropertyNames.ip_address)(o)?.ToString();
            return string.IsNullOrEmpty(text) ? "-" : text;
        }

        public void Refresh()
        {
            _grid.Children.Clear();
            _grid.RowDefinitions.Clear();
            _grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            for (int c = 0; c < Columns.Length; c++)
            {
                var h = new Border
                {
                    BorderBrush = Brush.Parse("#E3E3E3"), BorderThickness = new Thickness(0, 0, 1, 1), Padding = new Thickness(4, 0), Height = 38,
                    Child = new TextBlock { Text = Columns[c], FontWeight = FontWeight.SemiBold, TextAlignment = TextAlignment.Center, HorizontalAlignment = HorizontalAlignment.Center, VerticalAlignment = VerticalAlignment.Center, FontSize = 12 }
                };
                Grid.SetColumn(h, c);
                _grid.Children.Add(h);
            }
            var row = 1;
            var all = Rows().ToList();
            int skipDeeperThan = -1;
            for (int i = 0; i < all.Count; i++)
            {
                var r = all[i];
                if (skipDeeperThan >= 0 && r.Depth > skipDeeperThan)
                    continue; // inside a collapsed node
                skipDeeperThan = -1;
                var hasChildren = i + 1 < all.Count && all[i + 1].Depth > r.Depth;
                var key = r.Object.opaque_ref;
                var collapsed = _collapsed.Contains(key);
                if (hasChildren && collapsed)
                    skipDeeperThan = r.Depth;
                _grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var name = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 4, Margin = new Thickness(4 + r.Depth * 25, 0, 0, 0) };
                // QueryPanel +/- expander, 12px before the icon
                var expander = new Image { Width = 9, Height = 9, Margin = new Thickness(0, 0, 3, 0), VerticalAlignment = VerticalAlignment.Center,
                    Source = hasChildren ? Images.Get(collapsed ? "tree_plus" : "tree_minus") : null, Cursor = hasChildren ? new Avalonia.Input.Cursor(Avalonia.Input.StandardCursorType.Hand) : null };
                if (hasChildren)
                    expander.PointerPressed += (_, _) => { if (!_collapsed.Remove(key)) _collapsed.Add(key); Refresh(); };
                name.Children.Add(expander);
                name.Children.Add(Views.Icon(Images.IconFor(r.Object)));
                var lines = new StackPanel();
                lines.Children.Add(new TextBlock { Text = r.Object.Name(), TextTrimming = TextTrimming.CharacterEllipsis });
                var desc = r.Object.Description();
                if (!string.IsNullOrEmpty(desc))
                    lines.Children.Add(new TextBlock { Text = desc, Opacity = 0.6, TextTrimming = TextTrimming.CharacterEllipsis, MaxWidth = Widths[0] - 52 - r.Depth * 25, Margin = new Thickness(0, -2, 0, 0) });
                name.Children.Add(lines);
                Grid.SetRow(name, row);
                _grid.Children.Add(name);
                var values = Values(r.Object);
                for (int c = 0; c < values.Length; c++)
                {
                    var cell = new StackPanel { VerticalAlignment = VerticalAlignment.Center, HorizontalAlignment = HorizontalAlignment.Center, Margin = new Thickness(4, 0) };
                    if (values[c].Bar >= 0)
                        cell.Children.Add(Bar(values[c].Bar));
                    cell.Children.Add(new TextBlock { Text = values[c].Text ?? "", HorizontalAlignment = HorizontalAlignment.Center, TextTrimming = TextTrimming.CharacterEllipsis, MaxWidth = Widths[c + 1] - 8 });
                    Grid.SetRow(cell, row);
                    Grid.SetColumn(cell, c + 1);
                    _grid.Children.Add(cell);
                }
                row++;
            }
        }

        /// <summary>The usage bar drawn above CPU and memory values.</summary>
        private static Control Bar(int percent)
        {
            var p = Math.Clamp(percent, 0, 100);
            var bar = new Grid { Width = 70, Height = 6, Margin = new Thickness(0, 0, 0, 2) };
            bar.Children.Add(new Border { Background = Brush.Parse("#3A3A3A") });
            bar.Children.Add(new Border { Background = Brush.Parse("#2F8FE0"), Width = 70 * p / 100.0, HorizontalAlignment = HorizontalAlignment.Left });
            return bar;
        }
    }
}
