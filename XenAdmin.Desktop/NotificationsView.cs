using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using Avalonia.Platform.Storage;
using Avalonia.Threading;
using XenAdmin.Actions;
using XenAdmin.Alerts;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// AlertSummaryPage: the alerts grid (severity, message, server/pool, date, Dismiss), severity and server filters,
    /// Refresh and Dismiss All / Dismiss Selected.
    /// shortcut: no column sorting or keyboard expansion, and the dismiss confirmations have no "do not show again"
    /// setting; add them with the rest of AlertSummaryPage.
    /// </summary>
    public class NotificationsView : DockPanel, ILiveView
    {
        private readonly Func<IReadOnlyList<IXenConnection>> _connections;
        private readonly Action<AsyncAction> _run;
        private readonly Grid _grid = new Grid { ColumnDefinitions = new ColumnDefinitions("22,26,*,150,170,80") };
        private readonly HashSet<AlertPriority> _hiddenSeverities = new HashSet<AlertPriority>();
        private readonly HashSet<string> _hiddenServers = new HashSet<string>(); // host uuids (MessageAlert.HostUuid)
        private readonly HashSet<string> _expanded = new HashSet<string>();
        private readonly HashSet<string> _selected = new HashSet<string>();
        private readonly HashSet<string> _dismissing = new HashSet<string>();
        private List<NotificationRules.AlertRow> _rows = new List<NotificationRules.AlertRow>();
        private string _signature;
        private NotificationRules.DateFilter _dateFilter = NotificationRules.DateFilter.All;
        private DateTime? _customFrom, _customTo; // DateFilterDialog's range when "Custom..." is chosen
        private bool _custom;

        /// <summary>Raised when the visible alerts or the filters change, so the title band and the left list follow.</summary>
        public event Action Changed;

        /// <summary>All alerts, filtered or not, for the title band and the left list (MainWindow.UpdateAlertCount).</summary>
        public int Total { get; private set; }
        public string FiltersText => FilterIsOn ? "Filters are ON" : "Filters are OFF";

        public NotificationsView(Func<IReadOnlyList<IXenConnection>> connections, Action<AsyncAction> run)
        {
            _connections = connections;
            _run = run;
            Background = Brushes.White;

            var severity = MenuButton("Filter by Severit&y", () => Enum.GetValues<AlertPriority>().OrderBy(p => p == AlertPriority.Unknown ? int.MaxValue : (int)p).Select(p =>
                Check(p == AlertPriority.Unknown ? "Unknown" : $"Priority {p.GetString()}", !_hiddenSeverities.Contains(p), on => Toggle(_hiddenSeverities, p, on))));
            var server = MenuButton("Filter by S&erver", () => _connections().Where(c => c.IsConnected).SelectMany(c => c.Cache.Hosts).OrderBy(h => h).Select(h =>
                Check(h.Name(), !_hiddenServers.Contains(h.uuid), on => Toggle(_hiddenServers, h.uuid, on))));
            // FilterDatesToolStripDropDownButton order: the periods, then Show All
            var date = MenuButton("Filter by &Date", () => Enum.GetValues<NotificationRules.DateFilter>()
                .Where(f => f != NotificationRules.DateFilter.All).Select(f =>
                    Check(Dialogs.M(NotificationRules.DateFilterText(f)), !_custom && _dateFilter == f, on => { _custom = false; _dateFilter = on ? f : NotificationRules.DateFilter.All; }))
                .Append(Check(Dialogs.M(Messages.DATE_FILTER_CUSTOM), _custom, _ => CustomDates()))
                .Append(Check(Dialogs.M(NotificationRules.DateFilterText(NotificationRules.DateFilter.All)), !_custom && _dateFilter == NotificationRules.DateFilter.All,
                    _ => { _custom = false; _dateFilter = NotificationRules.DateFilter.All; })));
            var export = new Button { Content = Dialogs.M("E&xport") };
            export.Click += async (_, _) => await Export();
            var refresh = new Button { Content = Dialogs.M("&Refresh") };
            refresh.Click += (_, _) => { _signature = null; Refresh(); };
            var dismiss = MenuButton("Dismiss A&ll", () => new[]
            {
                Action("Dismiss &All", DismissAll),
                Action("Dismiss Sele&cted", async () =>
                {
                    var selected = Shown.Where(r => _selected.Contains(r.Message.opaque_ref)).Select(r => r.Message).ToList();
                    if (selected.Count > 0 && await Dialogs.Confirm(Owner, Title, Messages.ALERT_DISMISS_SELECTED_CONFIRM, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION))
                        DismissMany(selected);
                })
            });
            var toolbar = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 4 };
            toolbar.Children.Add(severity);
            toolbar.Children.Add(server);
            toolbar.Children.Add(date);
            toolbar.Children.Add(refresh);
            toolbar.Children.Add(export);
            toolbar.Children.Add(dismiss);
            var bar = new Border { Background = Brush.Parse("#F5F5F5"), BorderBrush = Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1), Padding = new Thickness(4), Child = toolbar, Margin = new Thickness(0, 0, 0, 8) };

            var frame = new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = new ScrollViewer { Content = _grid } };
            var content = new DockPanel { Margin = new Thickness(12) };
            DockPanel.SetDock(bar, Dock.Top);
            content.Children.Add(bar);
            content.Children.Add(frame);
            Children.Add(content);
            Refresh();
        }

        private static void Toggle<T>(HashSet<T> hidden, T item, bool shown) { if (shown) hidden.Remove(item); else hidden.Add(item); }

        private static MenuItem Check(string text, bool on, Action<bool> changed)
        {
            var m = new MenuItem { Header = text, ToggleType = MenuItemToggleType.CheckBox, IsChecked = on };
            m.Click += (_, _) => changed(m.IsChecked);
            return m;
        }

        private static MenuItem Action(string text, Action a)
        {
            var m = new MenuItem { Header = Dialogs.M(text) };
            m.Click += (_, _) => a();
            return m;
        }

        private Button MenuButton(string text, Func<IEnumerable<MenuItem>> items)
        {
            var b = new Button { Content = Dialogs.M(text) + "  ▾" };
            b.Click += (_, _) =>
            {
                var menu = new ContextMenu { ItemsSource = items().ToList() };
                menu.Closed += (_, _) => { _signature = null; Refresh(); };
                menu.Open(b);
            };
            return b;
        }

        private const int MaxRows = 1000; // AlertSummaryPage shows the first 1000
        private const string Title = "Dismiss Alerts";
        private IEnumerable<NotificationRules.AlertRow> Shown => _rows;
        private Window Owner => TopLevel.GetTopLevel(this) as Window;
        private bool FilterIsOn => _hiddenSeverities.Count + _hiddenServers.Count > 0
                                   || (_custom ? _customFrom != null || _customTo != null : _dateFilter != NotificationRules.DateFilter.All);

        /// <summary>FilterDatesToolStripDropDownButton "Custom...": the DateFilterDialog range.</summary>
        private async void CustomDates()
        {
            var range = await Dialogs.DateFilter(Owner, _customFrom, _customTo);
            if (range == null) return;
            (_customFrom, _customTo) = range.Value;
            _custom = true;
            _signature = null;
            Refresh();
        }

        /// <summary>
        /// toolStripButtonExportAll_Click: with filters on, all alerts or only the shown ones; then a .csv file chosen by
        /// the user, written by an action so it shows in Events.
        /// </summary>
        private async System.Threading.Tasks.Task Export()
        {
            var all = true;
            if (FilterIsOn)
            {
                var choice = await Dialogs.Choose(Owner, "Export", Messages.ALERT_EXPORT_ALL_OR_FILTERED, Messages.EXPORT_ALL_BUTTON, Messages.EXPORT_FILTERED_BUTTON);
                if (choice == null) return;
                all = choice.Value;
            }
            var file = await Owner.StorageProvider.SaveFilePickerAsync(new Avalonia.Platform.Storage.FilePickerSaveOptions
            {
                Title = Messages.EXPORT_ALL,
                DefaultExtension = "csv",
                SuggestedFileName = "alerts.csv",
                ShowOverwritePrompt = true,
                FileTypeChoices = new[] { new Avalonia.Platform.Storage.FilePickerFileType(Messages.CSV_DESCRIPTION) { Patterns = new[] { "*.csv" } } }
            });
            var path = file?.TryGetLocalPath();
            if (path == null) return;
            var rows = all
                ? NotificationRules.AlertMessages(_connections()).Where(m => !_dismissing.Contains(m.opaque_ref)).OrderByDescending(m => m.timestamp).Select(NotificationRules.Row).ToList()
                : Shown.ToList();
            _run(new DelegatedAsyncAction(null, string.Format(Messages.EXPORT_SYSTEM_ALERTS, path), string.Format(Messages.EXPORTING_SYSTEM_ALERTS, path),
                string.Format(Messages.EXPORTED_SYSTEM_ALERTS, path), _ =>
                {
                    using var stream = new System.IO.StreamWriter(path, false, System.Text.Encoding.UTF8);
                    stream.WriteLine(NotificationRules.CsvHeader());
                    foreach (var r in rows)
                        stream.WriteLine(NotificationRules.CsvLine(r));
                }));
        }

        /// <summary>tsmiDismissAll_Click: with filters on, all alerts or only the shown ones; without, a confirmation.</summary>
        private async void DismissAll()
        {
            bool? all = true;
            if (FilterIsOn)
                all = await Dialogs.Choose(Owner, Title, Messages.ALERT_DISMISS_ALL_CONTINUE, Messages.DISMISS_ALL_CONFIRM_BUTTON, Messages.DISMISS_FILTERED_CONFIRM_BUTTON);
            else if (!await Dialogs.Confirm(Owner, Title, Messages.ALERT_DISMISS_ALL_NO_FILTER_CONTINUE, Messages.DISMISS_ALL_YES_CONFIRM_BUTTON, Messages.CANCEL))
                all = null;
            if (all == null)
                return;
            DismissMany(all == true ? NotificationRules.AlertMessages(_connections()).ToList() : Shown.Select(r => r.Message).ToList());
        }

        private async void DismissOne(Message m)
        {
            if (await Dialogs.Confirm(Owner, Title, Messages.ALERT_DISMISS_CONFIRM, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION))
                DismissMany(new[] { m });
        }

        /// <summary>Hides the alerts while Message.destroy runs; failed ones come back when the action ends.</summary>
        private void DismissMany(IEnumerable<Message> messages)
        {
            foreach (var group in messages.Where(NotificationRules.CanDismiss).GroupBy(m => m.Connection))
            {
                var refs = group.Select(m => m.opaque_ref).ToList();
                refs.ForEach(r => _dismissing.Add(r));
                var action = NotificationRules.Dismiss(group.Key, group.ToList());
                action.Completed += _ => Dispatcher.UIThread.Post(() =>
                {
                    refs.ForEach(r => _dismissing.Remove(r));
                    _signature = null;
                    Refresh();
                });
                _run(action);
            }
            _signature = null;
            Refresh();
        }

        /// <summary>Rebuilds the grid when the alerts, the filters or the expanded rows change.</summary>
        public void Refresh()
        {
            var all = NotificationRules.AlertMessages(_connections()).Where(m => !_dismissing.Contains(m.opaque_ref)).ToList();
            Total = all.Count;
            // filter and cut first, then format only the rows that change (MessageAlert texts resolve objects by uuid)
            var from = _custom ? _customFrom : NotificationRules.DateFrom(_dateFilter, DateTime.Now);
            var to = _custom ? _customTo : null;
            var shown = all.Where(m => !_hiddenSeverities.Contains(NotificationRules.PriorityOf(m))
                    && (_hiddenServers.Count == 0 || !_hiddenServers.Contains(NotificationRules.HostUuid(m)))
                    && !NotificationRules.HideByDate(m.TimestampLocal(), from, to))
                .OrderByDescending(m => m.timestamp).Take(MaxRows).ToList();
            var signature = string.Join("|", shown.Select(m => m.opaque_ref)) + "#" + string.Join(",", _expanded) + "#" + string.Join(",", _selected) + "#" + Total;
            if (signature == _signature) return;
            _signature = signature;
            _rows = shown.Select(NotificationRules.Row).ToList();
            if (Parent != null) Changed?.Invoke();

            _grid.Children.Clear();
            _grid.RowDefinitions.Clear();
            void Cell(Control c, int row, int col)
            {
                Grid.SetRow(c, row);
                Grid.SetColumn(c, col);
                _grid.Children.Add(c);
            }
            _grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            foreach (var (text, col) in new[] { ("Message", 2), ("Server / Pool", 3), ("Date", 4), ("Actions", 5) })
                Cell(new Border { Background = Views.Gradient("#FFFFFF", "#E8E8E8"), BorderBrush = Brush.Parse("#D5D5D5"), BorderThickness = new Thickness(0, 0, 1, 1), Padding = new Thickness(6, 3), Child = new TextBlock { Text = text } }, 0, col);
            Cell(new Border { Background = Views.Gradient("#FFFFFF", "#E8E8E8"), BorderBrush = Brush.Parse("#D5D5D5"), BorderThickness = new Thickness(0, 0, 1, 1) }, 0, 0);
            Cell(new Border { Background = Views.Gradient("#FFFFFF", "#E8E8E8"), BorderBrush = Brush.Parse("#D5D5D5"), BorderThickness = new Thickness(0, 0, 1, 1) }, 0, 1);

            var r = 1;
            foreach (var a in Shown)
            {
                var id = a.Message.opaque_ref;
                var selected = _selected.Contains(id);
                var expanded = _expanded.Contains(id);
                _grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var background = new Border { Background = selected ? Brush.Parse("#0078D7") : Brushes.Transparent };
                Grid.SetColumnSpan(background, 6);
                Cell(background, r, 0);
                var fg = selected ? Brushes.White : Brushes.Black;

                var expander = new TextBlock { Text = expanded ? "▾" : "▸", Foreground = fg, HorizontalAlignment = HorizontalAlignment.Center, Margin = new Thickness(0, 3), Cursor = new Avalonia.Input.Cursor(Avalonia.Input.StandardCursorType.Hand) };
                expander.PointerPressed += (_, e) => { if (!_expanded.Remove(id)) _expanded.Add(id); e.Handled = true; Refresh(); };
                Cell(expander, r, 0);
                var icon = Views.Icon(NotificationRules.Icon(a.Priority), 16);
                icon.Margin = new Thickness(4, 3);
                icon.VerticalAlignment = VerticalAlignment.Top;
                ToolTip.SetTip(icon, a.Priority == AlertPriority.Unknown ? Messages.UNKNOWN : $"Priority {a.Priority.GetString()}");
                Cell(icon, r, 1);
                var message = new StackPanel { Margin = new Thickness(6, 3) };
                message.Children.Add(new TextBlock { Text = a.Title, Foreground = fg, TextWrapping = TextWrapping.Wrap });
                if (expanded)
                    message.Children.Add(new SelectableTextBlock { Text = a.Description, Foreground = fg, TextWrapping = TextWrapping.Wrap, Margin = new Thickness(0, 4, 0, 2) });
                Cell(message, r, 2);
                Cell(new TextBlock { Text = a.AppliesTo, Foreground = fg, Margin = new Thickness(6, 3) }, r, 3);
                Cell(new TextBlock { Text = a.Time.ToString("MMM d, yyyy h:mm tt", CultureInfo.CurrentCulture), Foreground = fg, Margin = new Thickness(6, 3) }, r, 4);
                var dismiss = new Button { Content = "Dismiss", Padding = new Thickness(4, 0), Margin = new Thickness(2), IsEnabled = NotificationRules.CanDismiss(a.Message), MinWidth = 0, Height = 20, VerticalAlignment = VerticalAlignment.Top };
                dismiss.Click += (_, _) => DismissOne(a.Message);
                Cell(dismiss, r, 5);

                background.PointerPressed += (_, e) =>
                {
                    var ctrl = e.KeyModifiers.HasFlag(Avalonia.Input.KeyModifiers.Control);
                    if (!ctrl) _selected.Clear();
                    if (!_selected.Remove(id) || !ctrl) _selected.Add(id);
                    Refresh();
                };
                message.PointerPressed += (_, e) => background.RaiseEvent(e);
                r++;
            }
        }
    }
}
