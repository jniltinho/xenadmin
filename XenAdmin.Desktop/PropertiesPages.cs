using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAdmin.CustomFields;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// VMHAEditPage: the VM's start order and the delay before the next VM starts, and (HA on) its restart priority,
    /// saved with SetVMStartupOptionsAction or SetHaPrioritiesAction. shortcut: with HA on, the failures to tolerate stay
    /// the pool's current value (the WinForms page recomputes it live); an overcommitted plan is refused by the server.
    /// </summary>
    public sealed class StartOptionsPage : EditPage
    {
        private readonly VM _vm;
        private readonly Pool _pool;
        private readonly NumericUpDown _order, _delay;
        private readonly ComboBox _priority;
        private readonly List<(VM.HaRestartPriority Priority, string Text)> _priorities;

        public override string Title => Messages.START_UP_OPTIONS;
        public override string Icon => "000_RebootVM_h32bit_16";
        public override string SubText => _pool?.ha_enabled == true && _priority.SelectedIndex >= 0
            ? _priorities[_priority.SelectedIndex].Text
            : $"Start order: {(long)(_order.Value ?? 0)}, delay: {(long)(_delay.Value ?? 0)} s";

        public StartOptionsPage(VM vm)
        {
            _vm = vm;
            _pool = Helpers.GetPoolOfOne(vm.Connection);
            _order = new NumericUpDown { Minimum = 0, Maximum = long.MaxValue, Value = vm.order, FormatString = "0", Width = 140 };
            _delay = new NumericUpDown { Minimum = 0, Maximum = long.MaxValue, Value = vm.start_delay, FormatString = "0", Width = 140 };
            // VMHAEditPage: a VM that cannot be agile (a disk on local storage) is only offered the priorities that do
            // not restart it, besides the one it already has
            var agile = VmOperations.DisksShared(vm).All(shared => shared);
            _priorities = HaRules.Priorities(vm.Connection)
                .Where(p => HaRules.Allowed(p.Priority, agile) || p.Priority == vm.HARestartPriority()).ToList();
            _priority = new ComboBox { ItemsSource = _priorities.Select(p => p.Text).ToList(), Width = 220,
                SelectedIndex = Math.Max(0, _priorities.FindIndex(p => p.Priority == vm.HARestartPriority())) };
            foreach (var c in new Control[] { _order, _delay })
                ((NumericUpDown)c).ValueChanged += (_, _) => OnChanged();
            _priority.SelectionChanged += (_, _) => OnChanged();

            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,Auto,Auto"), RowDefinitions = new RowDefinitions("Auto,Auto,Auto"), RowSpacing = 10, ColumnSpacing = 10 };
            void Row(int r, string label, Control c, string unit = null)
            {
                var l = new Label { Content = Dialogs.M(label), Target = c, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
                Grid.SetRow(l, r); Grid.SetRow(c, r); Grid.SetColumn(c, 1);
                grid.Children.Add(l); grid.Children.Add(c);
                if (unit != null) { var u = new TextBlock { Text = unit, VerticalAlignment = VerticalAlignment.Center }; Grid.SetRow(u, r); Grid.SetColumn(u, 2); grid.Children.Add(u); }
            }
            Row(0, "&Start order:", _order);
            Row(1, "&Attempt to start next VM after:", _delay, "seconds");
            var panel = new StackPanel();
            if (_pool?.ha_enabled == true)
            {
                Row(2, "HA restart &priority:", _priority);
                panel.Children.Add(Rubric("HA is enabled on this pool: set this VM's restart priority, its start order and the delay before the next VM starts."));
            }
            else
                panel.Children.Add(Rubric("Set the order in which this VM starts in a vApp or when HA restarts it, and the delay before the next VM starts."));
            panel.Children.Add(grid);
            View = panel;
        }

        private long Order => (long)(_order.Value ?? 0);
        private long Delay => (long)(_delay.Value ?? 0);
        private VM.HaRestartPriority Priority => _priority.SelectedIndex >= 0 ? _priorities[_priority.SelectedIndex].Priority : _vm.HARestartPriority();

        public override bool HasChanged => Order != _vm.order || Delay != _vm.start_delay || _pool?.ha_enabled == true && Priority != _vm.HARestartPriority();

        public override AsyncAction SaveSettings()
        {
            var settings = new Dictionary<VM, VMStartupOptions> { [_vm] = new VMStartupOptions(Order, Delay, Priority) };
            return _pool?.ha_enabled == true
                ? new SetHaPrioritiesAction(_vm.Connection, settings, _pool.ha_host_failures_to_tolerate, true)
                : new SetVMStartupOptionsAction(_vm.Connection, settings, true);
        }
    }

    /// <summary>VMAdvancedEditPage: the shadow memory multiplier (general use 1, Virtual Apps 4, or manual), live when running.</summary>
    public sealed class AdvancedPage : EditPage
    {
        private const double General = 1.0, VirtualApps = 4.0;
        private readonly VM _vm;
        private readonly RadioButton _general = new RadioButton { Content = Dialogs.M("Optimize for &general use"), GroupName = "shadow" };
        private readonly RadioButton _apps = new RadioButton { Content = Dialogs.M("Optimize for &Virtual Apps"), GroupName = "shadow" };
        private readonly RadioButton _manual = new RadioButton { Content = Dialogs.M("Optimize manually (&advanced use only)"), GroupName = "shadow" };
        private readonly TextBox _value = new TextBox { Width = 90 };

        public override string Title => Messages.ADVANCED_OPTIONS;
        public override string Icon => "002_Configure_h32bit_16";
        public override string SubText => _general.IsChecked == true ? "Optimize for general use"
            : _apps.IsChecked == true ? "Optimize for Virtual Apps" : string.Format(Messages.SHADOW_MEMORY_MULTIPLIER, _value.Text);

        public AdvancedPage(VM vm)
        {
            _vm = vm;
            var mul = vm.HVM_shadow_multiplier;
            _value.Text = mul.ToString(CultureInfo.InvariantCulture);
            (mul == General ? _general : mul == VirtualApps ? _apps : _manual).IsChecked = true;
            _general.IsCheckedChanged += (_, _) => { if (_general.IsChecked == true) _value.Text = "1"; };
            _apps.IsCheckedChanged += (_, _) => { if (_apps.IsChecked == true) _value.Text = "4"; };
            _value.TextChanged += (_, _) => OnChanged();
            _value.IsEnabled = _manual.IsChecked == true;
            _manual.IsCheckedChanged += (_, _) => _value.IsEnabled = _manual.IsChecked == true;
            var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(22, 0, 0, 0) };
            row.Children.Add(new Label { Content = Dialogs.M("Shadow memory &multiplier:"), Target = _value, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) });
            row.Children.Add(_value);
            var panel = new StackPanel { Spacing = 8 };
            panel.Children.Add(Rubric("Use these options to optimize this VM for increased performance in different environments."));
            panel.Children.AddRange(new Control[] { _general, _apps, _manual, row });
            // VMAdvancedEditPage: not while suspended or paused
            if (vm.power_state == vm_power_state.Suspended || vm.power_state == vm_power_state.Paused)
            {
                foreach (var c in new Control[] { _general, _apps, _manual, _value }) c.IsEnabled = false;
                panel.Children.Add(new TextBlock { Text = "These settings cannot be adjusted when the VM is suspended", Foreground = Brushes.DarkGoldenrod });
            }
            View = panel;
        }

        private double Value => double.TryParse(_value.Text, NumberStyles.Any, CultureInfo.InvariantCulture, out var v) ? v : -1;

        public override bool HasChanged => Value != _vm.HVM_shadow_multiplier;

        public override string ValidationError => Value >= 1.0 ? null : Messages.SHADOW_MEMORY_MULTIPLIER_VALUE;

        public override AsyncAction SaveSettings()
        {
            var value = Value;
            if (_vm.power_state == vm_power_state.Running)
                return new DelegatedAsyncAction(_vm.Connection, Messages.ACTION_CHANGE_SHADOW_MULTIPLIER,
                    string.Format(Messages.ACTION_CHANGING_SHADOW_MULTIPLIER_FOR, _vm), null,
                    s => VM.set_shadow_multiplier_live(s, _vm.opaque_ref, value), true, "vm.set_shadow_multiplier_live");
            _vm.HVM_shadow_multiplier = value; // the clone: SaveChangesAction writes it
            return null;
        }
    }

    /// <summary>
    /// CustomFieldsDisplayPage: a value for each custom field of the pool (text, or a date), and Edit Custom Fields
    /// (add or delete a field for every object of the pool), saved with SaveCustomFieldsAction.
    /// </summary>
    public sealed class CustomFieldsPage : EditPage
    {
        private readonly IXenObject _o;
        private readonly StackPanel _rows = new StackPanel { Spacing = 6 };
        private readonly List<(CustomFieldDefinition Def, Func<object> Value)> _fields = new List<(CustomFieldDefinition, Func<object>)>();

        public override string Title => Messages.CUSTOM_FIELDS;
        public override string Icon => "000_Fields_h32bit_16";
        public override string SubText => _fields.Count == 0 ? Messages.NONE : string.Join(", ", _fields.Select(f => f.Def.Name));

        public CustomFieldsPage(IXenObject o)
        {
            _o = o;
            var edit = new Button { Content = Dialogs.M("&Edit Custom Fields..."), HorizontalAlignment = HorizontalAlignment.Left };
            edit.Click += async (_, _) => await EditFields();
            var panel = new StackPanel { Spacing = 10 };
            panel.Children.Add(Rubric("Custom fields hold your own information about this object, such as its owner or its expiry date."));
            panel.Children.Add(_rows);
            panel.Children.Add(edit);
            View = panel;
            Fill();
            // follow the pool's field definitions while the page is shown
            Action changed = () => Avalonia.Threading.Dispatcher.UIThread.Post(Fill);
            View.AttachedToVisualTree += (_, _) => { CustomFieldsManager.CustomFieldsChanged += changed; Fill(); }; // changes while hidden
            View.DetachedFromVisualTree += (_, _) => CustomFieldsManager.CustomFieldsChanged -= changed;
        }

        private void Fill()
        {
            var previous = _fields.ToDictionary(f => f.Def.Name, f => f.Value());
            _fields.Clear();
            _rows.Children.Clear();
            foreach (var def in CustomFieldsManager.GetCustomFields(_o.Connection).OrderBy(d => d.Name))
            {
                var current = previous.TryGetValue(def.Name, out var kept) ? kept : CustomFieldsManager.GetCustomFieldValue(_o, def);
                var row = new Grid { ColumnDefinitions = new ColumnDefinitions("160,*") };
                row.Children.Add(new TextBlock { Text = def.Name + ":", VerticalAlignment = VerticalAlignment.Center, TextTrimming = TextTrimming.CharacterEllipsis });
                Control input;
                Func<object> value;
                if (def.Type == CustomFieldDefinition.Types.Date)
                {
                    // the date and its time of day, as the WinForms DateTimePicker pair
                    var was = current as DateTime?;
                    var picker = new CalendarDatePicker { SelectedDate = was?.Date, Width = 160 };
                    var time = new TimePicker { SelectedTime = was?.TimeOfDay, ClockIdentifier = "24HourClock" };
                    picker.SelectedDateChanged += (_, _) => OnChanged();
                    time.SelectedTimeChanged += (_, _) => OnChanged();
                    var both = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
                    both.Children.Add(picker);
                    both.Children.Add(time);
                    input = both;
                    value = () => picker.SelectedDate is DateTime d ? d.Date + (time.SelectedTime ?? TimeSpan.Zero) : (DateTime?)null;
                }
                else
                {
                    var box = new TextBox { Text = current as string ?? "" };
                    box.TextChanged += (_, _) => OnChanged();
                    input = box;
                    value = () => string.IsNullOrEmpty(box.Text) ? null : box.Text;
                }
                Grid.SetColumn(input, 1);
                row.Children.Add(input);
                _rows.Children.Add(row);
                _fields.Add((def, value));
            }
            if (_fields.Count == 0)
                _rows.Children.Add(new TextBlock { Text = "There are no custom fields. Use Edit Custom Fields to add one.", Opacity = 0.7 });
            OnChanged();
        }

        private List<CustomField> Changed() =>
            _fields.Select(f => (f.Def, Now: f.Value(), Was: CustomFieldsManager.GetCustomFieldValue(_o, f.Def)))
                .Where(x => !(x.Now == null && x.Was == null) && !(x.Now != null && x.Now.Equals(x.Was)))
                .Select(x => new CustomField(x.Def, x.Now)).ToList();

        public override bool HasChanged => Changed().Count > 0;

        public override AsyncAction SaveSettings() => new SaveCustomFieldsAction(_o, Changed(), true);

        /// <summary>CustomFieldsDialog: add a field (name and type) or delete one, for the whole pool.</summary>
        private async System.Threading.Tasks.Task EditFields()
        {
            var owner = TopLevel.GetTopLevel(View) as Window;
            var name = new TextBox { Watermark = "Name of the new field", Width = 220 };
            var type = new ComboBox { ItemsSource = new[] { "Text", "Date and time" }, SelectedIndex = 0, Width = 140 };
            var add = new Button { Content = Dialogs.M("&Add") };
            var list = new StackPanel { Spacing = 4 };
            var close = new Button { Content = "Close", IsCancel = true, MinWidth = 75 };
            void FillList()
            {
                list.Children.Clear();
                foreach (var def in CustomFieldsManager.GetCustomFields(_o.Connection).OrderBy(d => d.Name))
                {
                    var row = new DockPanel();
                    var delete = new Button { Content = Dialogs.M("&Delete"), Margin = new Thickness(8, 0, 0, 0) };
                    delete.Click += async (_, _) =>
                    {
                        if (await Dialogs.Confirm(owner, "Delete Custom Field", $"Delete the custom field '{def.Name}' and its values on every object?", Messages.YES, Messages.NO))
                            new RemoveCustomFieldAction(_o.Connection, def).RunAsync();
                    };
                    DockPanel.SetDock(delete, Dock.Right);
                    row.Children.Add(delete);
                    row.Children.Add(new TextBlock { Text = $"{def.Name}   ({(def.Type == CustomFieldDefinition.Types.Date ? "Date and time" : "Text")})", VerticalAlignment = VerticalAlignment.Center });
                    list.Children.Add(row);
                }
            }
            FillList();
            Action changed = () => Avalonia.Threading.Dispatcher.UIThread.Post(FillList);
            CustomFieldsManager.CustomFieldsChanged += changed;
            add.Click += (_, _) =>
            {
                var n = name.Text?.Trim();
                if (string.IsNullOrEmpty(n) || CustomFieldsManager.GetCustomFields(_o.Connection).Any(d => d.Name == n)) return;
                new AddCustomFieldAction(_o.Connection, new CustomFieldDefinition(n, type.SelectedIndex == 1 ? CustomFieldDefinition.Types.Date : CustomFieldDefinition.Types.String)).RunAsync();
                name.Text = "";
            };
            var addRow = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            addRow.Children.AddRange(new Control[] { name, type, add });
            var panel = new StackPanel { Spacing = 10 };
            panel.Children.Add(new ScrollViewer { Content = list, MaxHeight = 260 });
            panel.Children.Add(addRow);
            panel.Children.Add(Dialogs.Buttons(close));
            var w = Dialogs.Frame("Edit Custom Fields", panel, 520);
            close.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            CustomFieldsManager.CustomFieldsChanged -= changed;
        }
    }
}
