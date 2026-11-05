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

namespace XenAdmin.Desktop
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
        /// <summary>Show Columns: the columns the user hid (Name always shows), for every Search tab.</summary>
        private static readonly HashSet<int> HiddenColumns = new HashSet<int>();

        /// <summary>The main window shows a search in the tree (null: back to Infrastructure).</summary>
        public static Action<Search> ShowSearch;
        /// <summary>The main window runs and tracks an action.</summary>
        public static Action<XenAdmin.Actions.AsyncAction> Run;

        private readonly Search _search;


        private readonly Func<IEnumerable<(IXenObject Object, int Depth)>> _rows;

        /// <param name="rows">The objects to list (a group of an organization view or a saved search); by default the overview.</param>
        /// <param name="search">The search shown (Edit Search and Export start from it); the overview when null.</param>
        /// <param name="canEdit">False for a group of an organization view or search: Edit Search and Export need its definition.</param>
        public SearchView(IReadOnlyList<IXenConnection> connections, Func<IEnumerable<(IXenObject Object, int Depth)>> rows = null, Search search = null, bool canEdit = true)
        {
            _connections = connections;
            _rows = rows;
            _search = search;
            Background = Brushes.White;
            if (_metrics == null)
            {
                _metrics = new MetricUpdater();
                PropertyAccessorHelper.MetricUpdater = _metrics;
            }

            var buttons = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            Button B(string text, Action click)
            {
                var b = new Button { Content = text };
                b.Click += (_, _) => click();
                buttons.Children.Add(b);
                return b;
            }
            B("New Search", async () => Apply(await SearchEditor.Show(Owner, null, _connections)));
            B("Edit Search", async () => Apply(await SearchEditor.Show(Owner, _search ?? SearchRules.Overview(), _connections))).IsEnabled = canEdit;
            Button saved = null;
            saved = B("Saved Searches  ▾", () =>
            {
                var menu = new ContextMenu();
                foreach (var s in NavigationViews.SavedSearches())
                {
                    var item = new MenuItem { Header = s.Name.Replace("_", "__"), Icon = Views.Icon(s.DefaultSearch ? "000_defaultSpyglass_h32bit_16" : "000_Search_h32bit_16") };
                    item.Click += (_, _) => ShowSearch?.Invoke(s);
                    menu.Items.Add(item);
                }
                menu.Open(saved);
            });
            B("Export...", async () => await Export()).IsEnabled = canEdit;
            B("Import...", async () => await Import());
            var top = new Border { BorderBrush = Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1, 1, 1, 0), Padding = new Thickness(10, 8), Child = buttons };
            var columns = new Button { Content = "Show Columns  ▾", HorizontalAlignment = HorizontalAlignment.Right, Margin = new Thickness(0, 6, 4, 6) };
            columns.Click += (_, _) =>
            {
                var menu = new ContextMenu();
                for (var i = 1; i < Columns.Length; i++)
                {
                    var c = i;
                    var item = new MenuItem { Header = Columns[c].Replace("\n", " "), ToggleType = MenuItemToggleType.CheckBox, IsChecked = !HiddenColumns.Contains(c) };
                    item.Click += (_, _) => { if (!HiddenColumns.Remove(c)) HiddenColumns.Add(c); Refresh(); };
                    menu.Items.Add(item);
                }
                menu.Open(columns);
            };
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
            if (_rows != null)
            {
                foreach (var (o, depth) in _rows())
                    yield return new Row { Object = o, Depth = depth };
                yield break;
            }
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



        public void Refresh()
        {
            _grid.Children.Clear();
            _grid.RowDefinitions.Clear();
            _grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            for (int c = 0; c < Columns.Length; c++)
                _grid.ColumnDefinitions[c].Width = new GridLength(HiddenColumns.Contains(c) ? 0 : Widths[c], GridUnitType.Pixel);
            for (int c = 0; c < Columns.Length; c++)
            {
                if (HiddenColumns.Contains(c)) continue;
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
                name.Children.Add(Views.Icon(IconNames.For(r.Object)));
                var lines = new StackPanel();
                lines.Children.Add(new TextBlock { Text = r.Object.Name(), TextTrimming = TextTrimming.CharacterEllipsis });
                var desc = r.Object.Description();
                if (!string.IsNullOrEmpty(desc))
                    lines.Children.Add(new TextBlock { Text = desc, Opacity = 0.6, TextTrimming = TextTrimming.CharacterEllipsis, MaxWidth = Widths[0] - 52 - r.Depth * 25, Margin = new Thickness(0, -2, 0, 0) });
                name.Children.Add(lines);
                Grid.SetRow(name, row);
                _grid.Children.Add(name);
                var values = OverviewColumns.Values(r.Object);
                for (int c = 0; c < values.Length; c++)
                {
                    if (HiddenColumns.Contains(c + 1)) continue;
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

        private Window Owner => TopLevel.GetTopLevel(this) as Window;

        /// <summary>The editor's result: show the search, save it to a pool first, or delete it.</summary>
        private void Apply(SearchEditor.Result result)
        {
            if (result == null) return;
            if (result.Delete)
            {
                var c = result.Search.Connection ?? _connections.FirstOrDefault(x => x.IsConnected);
                if (c != null) Run?.Invoke(SearchRules.DeleteAction(result.Search, c));
                ShowSearch?.Invoke(null);
                return;
            }
            if (result.SaveTo != null)
                Run?.Invoke(SearchRules.SaveAction(result.Search, result.SaveTo));
            ShowSearch?.Invoke(result.Search);
        }

        private static readonly Avalonia.Platform.Storage.FilePickerFileType SearchFiles =
            new Avalonia.Platform.Storage.FilePickerFileType("XenCenter search (*.xensearch)") { Patterns = new[] { "*.xensearch" } };

        /// <summary>SearchPage export: the search shown (the overview by default) to a .xensearch file.</summary>
        private async System.Threading.Tasks.Task Export()
        {
            var search = _search ?? SearchRules.Overview();
            var f = await Owner.StorageProvider.SaveFilePickerAsync(new Avalonia.Platform.Storage.FilePickerSaveOptions
            {
                Title = "Export Search", SuggestedFileName = (string.IsNullOrEmpty(search.Name) ? "Overview" : search.Name) + ".xensearch",
                DefaultExtension = "xensearch", ShowOverwritePrompt = true, FileTypeChoices = new[] { SearchFiles }
            });
            if (f == null || !(Avalonia.Platform.Storage.StorageProviderExtensions.TryGetLocalPath(f) is string path)) return;
            try { SearchRules.Export(search, path); }
            catch (Exception e) { await Dialogs.Info(Owner, "Export Search", e.Message); }
        }

        /// <summary>SearchPage import: a .xensearch file's search is shown (Edit Search can then save it to a pool).</summary>
        private async System.Threading.Tasks.Task Import()
        {
            var files = await Owner.StorageProvider.OpenFilePickerAsync(new Avalonia.Platform.Storage.FilePickerOpenOptions
                { Title = "Import Search", AllowMultiple = false, FileTypeFilter = new[] { SearchFiles } });
            if (files.Count == 0 || !(Avalonia.Platform.Storage.StorageProviderExtensions.TryGetLocalPath(files[0]) is string path)) return;
            var searches = SearchRules.Import(path);
            if (searches == null || searches.Count == 0)
            {
                await Dialogs.Info(Owner, "Import Search", "The file is not a search exported by XCP-ng Center.");
                return;
            }
            ShowSearch?.Invoke(searches[0]);
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
