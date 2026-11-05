using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;

using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using Avalonia.Threading;
using XenAdmin.Core;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// PerformancePage: the default graphs (GraphList) over the RRD archives, the Zoom menu and the lifecycle events.
    /// shortcut: no graph designer (New/Edit/Delete Graph, Move Up/Down) and no DataPlotNav overview; add with custom graphs.
    /// </summary>
    public class PerformanceView : DockPanel, ILiveView
    {
        // Palette.DefaultColors
        private static readonly IBrush[] Colors = new[] { "#0066FF", "#339900", "#FF9900", "#660099", "#990033", "#995E00", "#E6E600", "#00FFB4", "#009982", "#BFA4A6" }
            .Select(c => (IBrush)Brush.Parse(c)).ToArray();

        private readonly IXenObject _o;
        private readonly StackPanel _graphs = new StackPanel { Spacing = 14, Margin = new Thickness(8) };
        private readonly StackPanel _events = new StackPanel { Spacing = 2, Margin = new Thickness(4) };
        private readonly TextBlock _status = new TextBlock { Opacity = 0.7, Margin = new Thickness(8, 4) };
        private readonly DispatcherTimer _timer = new DispatcherTimer();
        private PerformanceRules.Zoom _zoom = PerformanceRules.Zoom.TenMinutes;
        private Dictionary<string, string> _units;
        private bool _busy;
        private int _generation;

        public PerformanceView(IXenObject o)
        {
            _o = o;
            Background = Brushes.White;

            var zoom = new Button { Content = "Zoom  ▾" };
            var menu = new ContextMenu();
            foreach (var z in Enum.GetValues<PerformanceRules.Zoom>())
            {
                var item = new MenuItem { Header = Dialogs.M(PerformanceRules.ZoomText(z)) };
                item.Click += (_, _) => { _zoom = z; Reload(); };
                menu.Items.Add(item);
            }
            zoom.Click += (_, _) => menu.Open(zoom);
            var toolbar = new DockPanel { Margin = new Thickness(0, 0, 0, 6) };
            DockPanel.SetDock(zoom, Dock.Right);
            toolbar.Children.Add(zoom);
            var left = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
            left.Children.Add(new Button { Content = Dialogs.M("Move &Up"), IsEnabled = false });
            left.Children.Add(new Button { Content = Dialogs.M("Move &Down"), IsEnabled = false });
            left.Children.Add(new Button { Content = Dialogs.M("&Actions") + "  ▾", IsEnabled = false });
            toolbar.Children.Add(left);

            var graphsFrame = new Border { BorderBrush = Brush.Parse("#A0A0A0"), BorderThickness = new Thickness(1), Child = new ScrollViewer { Content = _graphs } };
            var main = new DockPanel { Margin = new Thickness(10) };
            DockPanel.SetDock(toolbar, Dock.Top);
            DockPanel.SetDock(_status, Dock.Bottom);
            main.Children.Add(toolbar);
            main.Children.Add(_status);
            main.Children.Add(graphsFrame);

            var eventsTitle = new Border
            {
                Background = Views.Gradient("#F4F4F4", "#DCDCDC", false), BorderBrush = Brush.Parse("#A0A0A0"), BorderThickness = new Thickness(0, 0, 0, 1),
                Child = new TextBlock { Text = "VM Lifecycle Events", FontWeight = FontWeight.SemiBold, Margin = new Thickness(8, 4) }
            };
            var events = new DockPanel();
            DockPanel.SetDock(eventsTitle, Dock.Top);
            events.Children.Add(eventsTitle);
            events.Children.Add(new ScrollViewer { Content = _events });
            var eventsFrame = new Border { Width = 160, Margin = new Thickness(0, 10, 10, 10), BorderBrush = Brush.Parse("#A0A0A0"), BorderThickness = new Thickness(1), Child = events };
            DockPanel.SetDock(eventsFrame, Dock.Right);
            Children.Add(eventsFrame);
            Children.Add(main);

            _timer.Tick += (_, _) => Load();
            AttachedToVisualTree += (_, _) => Reload();
            DetachedFromVisualTree += (_, _) => { _timer.Stop(); _generation++; };
            Refresh();
        }

        /// <summary>The cache changed: the lifecycle events may have.</summary>
        public void Refresh()
        {
            _events.Children.Clear();
            foreach (var m in PerformanceRules.LifecycleEvents(_o))
            {
                var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 4 };
                row.Children.Add(Views.Icon(PerformanceRules.EventIcon(m), 16));
                row.Children.Add(new TextBlock { Text = m.TimestampLocal().ToString("MMM d, yyyy h:mm tt", CultureInfo.CurrentCulture), FontSize = 11, VerticalAlignment = VerticalAlignment.Center });
                ToolTip.SetTip(row, m.Name());
                _events.Children.Add(row);
            }
        }

        private void Reload()
        {
            var (step, _) = PerformanceRules.Span(_zoom);
            _timer.Interval = TimeSpan.FromSeconds(Math.Max(5, Math.Min(step, 60)));
            _timer.Start();
            Load();
        }

        private async void Load()
        {
            if (_busy || !_o.Connection.IsConnected) return;
            _busy = true;
            var generation = _generation;
            var zoom = _zoom;
            try
            {
                var data = await System.Threading.Tasks.Task.Run(() =>
                {
                    // shared session: never log it out
                    _units ??= PerformanceRules.Units(_o, _o.Connection.DuplicateSession());
                    return PerformanceRules.Fetch(_o, zoom, _units);
                });
                if (generation != _generation || zoom != _zoom) return;
                Draw(data);
                _status.Text = "";
            }
            catch (Exception e)
            {
                _status.Text = "Unable to load the performance data: " + e.Message;
            }
            finally
            {
                _busy = false;
            }
        }

        private void Draw(Dictionary<string, List<PerformanceRules.Point>> data)
        {
            var (_, window) = PerformanceRules.Span(_zoom);
            var end = DateTime.UtcNow - _o.Connection.ServerTimeOffset;
            var start = end - window;
            _graphs.Children.Clear();
            var colour = 0;
            foreach (var g in PerformanceRules.DefaultGraphs(_o).Where(g => g.Series.Any(s => !s.Hidden)))
            {
                var series = g.Series.Select(s => (s.Name, s.Hidden, Brush: Colors[colour++ % Colors.Length],
                    Points: data.TryGetValue(s.Id, out var p) ? p.Where(x => x.Time >= start).ToList() : new List<PerformanceRules.Point>()))
                    .Where(s => !s.Hidden).ToList();
                var dataMax = series.SelectMany(s => s.Points).Select(p => p.Value).Where(double.IsFinite).DefaultIfEmpty(0).Max();
                var (max, divisor, label) = PerformanceRules.Scale(g.Unit, g.FixedMax, dataMax);

                var legend = new StackPanel { Margin = new Thickness(6, 4) };
                foreach (var s in series)
                {
                    var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
                    row.Children.Add(new Border { Width = 10, Height = 2, Background = s.Brush, VerticalAlignment = VerticalAlignment.Center });
                    row.Children.Add(new TextBlock { Text = s.Name, FontSize = 11 });
                    legend.Children.Add(row);
                }
                var plot = new Plot(series.Select(s => (s.Brush, s.Points)).ToList(), start, end, max, divisor, label, _o.Connection.ServerTimeOffset) { Height = 130 };
                var line = new Grid { ColumnDefinitions = new ColumnDefinitions("*,150"), ColumnSpacing = 10 };
                // DataKey is a scrolling list: hosts may have dozens of CPUs
                var legendBox = new Border
                {
                    BorderBrush = Brush.Parse("#A0A0A0"), BorderThickness = new Thickness(1), Height = 110, VerticalAlignment = VerticalAlignment.Top, Margin = new Thickness(0, 10, 0, 0),
                    Child = new ScrollViewer { Content = legend, HorizontalScrollBarVisibility = Avalonia.Controls.Primitives.ScrollBarVisibility.Auto }
                };
                Grid.SetColumn(legendBox, 1);
                line.Children.Add(plot);
                line.Children.Add(legendBox);
                var panel = new StackPanel();
                panel.Children.Add(new TextBlock { Text = g.Title, Margin = new Thickness(26, 0, 0, 2) });
                panel.Children.Add(line);
                _graphs.Children.Add(panel);
            }
        }

        /// <summary>DataPlot: grid, one line per series (gaps at NaN), value axis on the right, time axis below.</summary>
        private sealed class Plot : Control
        {
            private static readonly IPen Grid = new Pen(new SolidColorBrush(Color.FromArgb(64, 109, 109, 109)), 1);
            private static readonly IPen Frame = new Pen(Brushes.Gray, 1);
            private readonly List<(IBrush Brush, List<PerformanceRules.Point> Points)> _series;
            private readonly DateTime _start, _end;
            private readonly double _max, _divisor;
            private readonly string _unit;
            private readonly TimeSpan _serverOffset;

            public Plot(List<(IBrush, List<PerformanceRules.Point>)> series, DateTime start, DateTime end, double max, double divisor, string unit, TimeSpan serverOffset)
            {
                _series = series; _start = start; _end = end; _max = max <= 0 ? 1 : max; _divisor = divisor; _unit = unit; _serverOffset = serverOffset;
            }

            private static FormattedText Text(string s) =>
                new FormattedText(s, CultureInfo.CurrentCulture, FlowDirection.LeftToRight, new Typeface("Segoe UI, Noto Sans, DejaVu Sans"), 11, Brushes.Black);

            public override void Render(DrawingContext dc)
            {
                var area = new Rect(26, 6, Math.Max(10, Bounds.Width - 26 - 44), Math.Max(10, Bounds.Height - 6 - 22));
                dc.FillRectangle(Brushes.White, area);
                for (var i = 1; i < 10; i++)
                {
                    var x = area.X + area.Width * i / 10;
                    dc.DrawLine(Grid, new Point(x, area.Top), new Point(x, area.Bottom));
                }
                for (var i = 1; i < 5; i++)
                {
                    var y = area.Y + area.Height * i / 5;
                    dc.DrawLine(Grid, new Point(area.Left, y), new Point(area.Right, y));
                }
                dc.DrawRectangle(null, Frame, area);

                var span = (_end - _start).TotalSeconds;
                foreach (var (brush, points) in _series)
                {
                    var pen = new Pen(brush, 1.2);
                    // DisplayOptionsPage "Area graph": each unbroken run of points is filled down to the axis
                    var fill = DesktopConfigProvider.CurrentOptions.FillAreaUnderGraphs && brush is ISolidColorBrush solid
                        ? new SolidColorBrush(solid.Color, 0.35) : null;
                    var run = new List<Point>();
                    void Close()
                    {
                        if (fill != null && run.Count > 1)
                        {
                            var g = new StreamGeometry();
                            using (var ctx = g.Open())
                            {
                                ctx.BeginFigure(new Point(run[0].X, area.Bottom), true);
                                run.ForEach(p => ctx.LineTo(p));
                                ctx.LineTo(new Point(run[^1].X, area.Bottom));
                                ctx.EndFigure(true);
                            }
                            dc.DrawGeometry(fill, null, g);
                        }
                        for (var i = 1; i < run.Count; i++) dc.DrawLine(pen, run[i - 1], run[i]);
                        run.Clear();
                    }
                    foreach (var p in points)
                    {
                        if (!double.IsFinite(p.Value) || p.Value < 0) { Close(); continue; }
                        var x = area.Left + area.Width * (p.Time - _start).TotalSeconds / span;
                        var y = area.Bottom - area.Height * Math.Min(1, p.Value / _max);
                        run.Add(new Point(x, y));
                    }
                    Close();
                }

                // value axis on the right: max, unit, 0
                var top = (_max / _divisor).ToString(_max / _divisor >= 10 ? "0" : "0.#", CultureInfo.CurrentCulture);
                dc.DrawText(Text(top), new Point(area.Right + 4, area.Top - 7));
                var unit = Text(_unit);
                unit.SetFontWeight(FontWeight.Bold);
                dc.DrawText(unit, new Point(area.Right + 4, area.Top + area.Height / 2 - 7));
                dc.DrawText(Text("0"), new Point(area.Right + 4, area.Bottom - 7));

                // time axis: five local times
                var fmt = span <= 86400 ? "t" : "MMM d";
                for (var i = 0; i <= 4; i++)
                {
                    // the points are in the server's time; the labels show this machine's clock (DataGrid adds ClientServerOffset)
                    var t = (_start.AddSeconds(span * (i + 0.5) / 5) + _serverOffset).ToLocalTime();
                    var label = Text(t.ToString(fmt, CultureInfo.CurrentCulture));
                    var x = area.Left + area.Width * (i + 0.5) / 5 - label.Width / 2;
                    dc.DrawText(label, new Point(x, area.Bottom + 3));
                }
            }
        }
    }
}
