using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using Avalonia.Threading;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// AssignPriorities with HaNtolIndicator: each protectable VM's restart priority, start order and delay (a VM that is
    /// not agile can only be "restart if possible" or "do not restart"), and the server failures to tolerate, whose
    /// maximum the pool computes for the current priorities.
    /// </summary>
    public class HaPlanEditor : StackPanel
    {
        private sealed class Row
        {
            public VM Vm;
            public ComboBox Priority;
            public NumericUpDown Order, Delay;
            public bool Agile = true;
            public TextBlock Agility;
        }

        private readonly Pool _pool;
        private readonly List<Row> _rows = new List<Row>();
        private readonly List<(VM.HaRestartPriority Priority, string Text)> _priorities;
        private readonly NumericUpDown _ntol = new NumericUpDown { Minimum = 0, Increment = 1, FormatString = "0", Width = 110 };
        private readonly TextBlock _ntolInfo = new TextBlock { TextWrapping = TextWrapping.Wrap };
        private long _max = -1;
        private int _generation;
        private bool _ntolTouched;   // the user set the failures to tolerate: never change it for them
        private string _failure;     // a background query failed (connection lost): the plan cannot be applied
        private bool _agilityPending = true; // priorities stay locked until the agility check set them

        /// <summary>The plan changed (re-check Error).</summary>
        public event Action Changed;

        public HaPlanEditor(Pool pool)
        {
            _pool = pool;
            Spacing = 10;
            _priorities = HaRules.Priorities(pool.Connection);
            var vms = HaRules.ProtectableVms(pool.Connection);
            // AssignPriorities.ProtectVmsByDefault: first-time defaults only when HA is being turned on
            var firstTime = !pool.ha_enabled && HaRules.FirstTime(vms);

            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("*,Auto,Auto,Auto,Auto"), ColumnSpacing = 8, RowSpacing = 4 };
            grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            string[] headers = { "VM", "HA restart priority", "Start order", "Delay (s)", "Agile" };
            for (int c = 0; c < headers.Length; c++)
            {
                var h = new TextBlock { Text = headers[c], FontWeight = FontWeight.SemiBold };
                Grid.SetColumn(h, c);
                grid.Children.Add(h);
            }
            foreach (var vm in vms)
            {
                var row = new Row
                {
                    Vm = vm,
                    Priority = new ComboBox { ItemsSource = _priorities.Select(p => p.Text).ToList(), Width = 170 },
                    Order = new NumericUpDown { Minimum = 0, Maximum = 9999, Value = vm.order, FormatString = "0", Width = 100 },
                    Delay = new NumericUpDown { Minimum = 0, Maximum = 9999, Value = vm.start_delay, FormatString = "0", Width = 100 },
                    Agility = new TextBlock { Text = "Checking...", VerticalAlignment = VerticalAlignment.Center }
                };
                Select(row, HaRules.InitialPriority(vm, firstTime, true));
                row.Priority.IsEnabled = false; // AssignPriorities: until VM.assert_agile answered
                row.Priority.SelectionChanged += (_, _) => Recompute();
                _rows.Add(row);
                var r = grid.RowDefinitions.Count;
                grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                Control[] cells = { new TextBlock { Text = vm.Name(), VerticalAlignment = VerticalAlignment.Center }, row.Priority, row.Order, row.Delay, row.Agility };
                for (int c = 0; c < cells.Length; c++) { Grid.SetRow(cells[c], r); Grid.SetColumn(cells[c], c); grid.Children.Add(cells[c]); }
            }
            Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Padding = new Thickness(6), Background = Brushes.White,
                Child = new ScrollViewer { Content = vms.Count > 0 ? grid : new TextBlock { Text = "There are no VMs in this pool that HA can protect." }, MaxHeight = 220 } });

            var ntolRow = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            ntolRow.Children.Add(new Label { Content = Dialogs.M("&Failures tolerated:"), Target = _ntol, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) });
            ntolRow.Children.Add(_ntol);
            _ntol.Maximum = Math.Max(1, pool.Connection.Cache.Hosts.Length);
            _ntol.Value = pool.ha_enabled ? pool.ha_host_failures_to_tolerate : 0;
            _ntolTouched = pool.ha_enabled;
            _ntol.ValueChanged += (_, _) => { _ntolTouched = true; Changed?.Invoke(); };
            Children.Add(new TextBlock { Text = "You can specify the number of server failures you need to be able to tolerate in the pool.", TextWrapping = TextWrapping.Wrap });
            Children.Add(ntolRow);
            Children.Add(_ntolInfo);

            CheckAgility(firstTime);
            Recompute();
        }

        public long Ntol => (long)(_ntol.Value ?? -1);

        /// <summary>AssignPriorities.updateButtons: null when the plan can be applied.</summary>
        public string Error => _failure ?? (_agilityPending ? "Checking which VMs are agile..." : HaRules.PlanError(Ntol, _max));

        /// <summary>AssignPriorities.GetChangedStartupOptions: every VM's priority, order and delay.</summary>
        public Dictionary<VM, VMStartupOptions> Options() =>
            _rows.ToDictionary(r => r.Vm, r => new VMStartupOptions((long)(r.Order.Value ?? 0), (long)(r.Delay.Value ?? 0), PriorityOf(r)));

        public IReadOnlyDictionary<VM.HaRestartPriority, int> Counts() =>
            _rows.GroupBy(PriorityOf).ToDictionary(g => g.Key, g => g.Count());

        private VM.HaRestartPriority PriorityOf(Row r) =>
            r.Priority.SelectedIndex >= 0 ? _priorities[r.Priority.SelectedIndex].Priority : VM.HaRestartPriority.BestEffort;

        private void Select(Row r, VM.HaRestartPriority p)
        {
            var i = _priorities.FindIndex(x => x.Priority == p);
            r.Priority.SelectedIndex = i >= 0 ? i : _priorities.FindIndex(x => x.Priority == VM.HaRestartPriority.BestEffort);
        }

        /// <summary>AssignPriorities.UpdateVMsAgility: VM.assert_agile in the background; a VM that is not agile drops to best effort.</summary>
        private void CheckAgility(bool firstTime)
        {
            var vms = _rows.Select(r => r.Vm).ToList();
            System.Threading.ThreadPool.QueueUserWorkItem(_ =>
            {
                Dictionary<string, string> reasons;
                try
                {
                    var session = _pool.Connection.DuplicateSession();
                    reasons = vms.ToDictionary(v => v.opaque_ref, v => HaRules.NotAgileReason(session, v));
                }
                catch (Exception e)
                {
                    Dispatcher.UIThread.Post(() => { _failure = e.Message; Changed?.Invoke(); });
                    return;
                }
                Dispatcher.UIThread.Post(() =>
                {
                    foreach (var row in _rows)
                    {
                        var reason = reasons[row.Vm.opaque_ref];
                        row.Agile = reason == null;
                        row.Agility.Text = row.Agile ? Messages.YES : Messages.NO;
                        ToolTip.SetTip(row.Agility, reason);
                        if (firstTime && row.Agile)
                            Select(row, VM.HaHighestProtectionAvailable(_pool.Connection));
                        else if (!HaRules.Allowed(PriorityOf(row), row.Agile))
                            Select(row, VM.HaRestartPriority.BestEffort);
                        row.Priority.IsEnabled = true;
                    }
                    _agilityPending = false;
                    Recompute();
                });
            });
        }

        /// <summary>HaNtolControl.updateNtol: the maximum failures for these priorities, computed off the UI thread.</summary>
        private void Recompute()
        {
            foreach (var row in _rows.Where(r => !HaRules.Allowed(PriorityOf(r), r.Agile)))
                Select(row, VM.HaRestartPriority.BestEffort);
            var settings = _rows.ToDictionary(r => r.Vm, PriorityOf);
            var mine = ++_generation;
            _max = -1;
            _ntolInfo.Text = "Computing the number of server failures the pool can tolerate...";
            Changed?.Invoke();
            System.Threading.ThreadPool.QueueUserWorkItem(_ =>
            {
                long max;
                try { max = HaRules.MaxFailures(_pool.Connection.DuplicateSession(), settings); }
                catch { max = -2; }
                Dispatcher.UIThread.Post(() =>
                {
                    if (mine != _generation) return;
                    _max = max < 0 ? -1 : max;
                    // HaNtolControl: a new configuration starts at the maximum; a value the user chose stays (overcommit is
                    // reported and blocks the plan instead of lowering the protection)
                    if (max >= 0 && !_ntolTouched)
                    {
                        _ntol.Value = max;
                        _ntolTouched = false;
                    }
                    _ntolInfo.Text = max == -2 ? Messages.HA_NTOL_UNKNOWN
                        : max == 0 ? Messages.HA_NTOL_ZERO
                        : string.Format(Messages.HA_NTOL_BLURB, max);
                    Changed?.Invoke();
                });
            });
        }
    }
}
