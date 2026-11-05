using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using Avalonia.Threading;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;
using Task = System.Threading.Tasks.Task;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// NewVMGroupVMsPage&lt;VM_appliance&gt;: the pool's VMs with a checkbox each and their current vApp, Select All /
    /// Clear All. Used by the New vApp wizard and as the vApp's Properties page.
    /// </summary>
    public sealed class VappVmsPage : EditPage
    {
        private readonly VM_appliance _vapp;
        private readonly List<(VM Vm, CheckBox Box)> _rows = new List<(VM, CheckBox)>();
        private readonly List<VM> _unlisted = new List<VM>(); // members the list does not show (hidden VMs): they stay
        private readonly TextBlock _count = new TextBlock { VerticalAlignment = VerticalAlignment.Center };

        public VappVmsPage(IXenConnection c, VM_appliance vapp, IEnumerable<VM> preselected = null)
        {
            _vapp = vapp;
            var chosen = new HashSet<string>((preselected ?? Enumerable.Empty<VM>()).Select(v => v.opaque_ref));
            if (vapp != null)
                chosen.UnionWith(vapp.VMs.Select(r => r.opaque_ref));

            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*,Auto"), ColumnSpacing = 10, RowSpacing = 4 };
            grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            void Cell(Control ctl, int row, int col) { Grid.SetRow(ctl, row); Grid.SetColumn(ctl, col); grid.Children.Add(ctl); }
            Cell(new TextBlock { Text = "VM", FontWeight = FontWeight.SemiBold }, 0, 1);
            Cell(new TextBlock { Text = Messages.CURRENT_VAPP, FontWeight = FontWeight.SemiBold }, 0, 2);
            var candidates = VappRules.Candidates(c);
            if (vapp != null)
                _unlisted.AddRange(c.ResolveAll(vapp.VMs).Where(v => candidates.All(x => x.opaque_ref != v.opaque_ref)));
            foreach (var vm in candidates)
            {
                var r = grid.RowDefinitions.Count;
                grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var box = new CheckBox { IsChecked = chosen.Contains(vm.opaque_ref), MinHeight = 0, Padding = new Thickness(0) };
                box.IsCheckedChanged += (_, _) => Update();
                _rows.Add((vm, box));
                Cell(box, r, 0);
                var name = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
                name.Children.Add(Views.Icon(IconNames.For(vm)));
                name.Children.Add(new TextBlock { Text = vm.Name(), VerticalAlignment = VerticalAlignment.Center });
                Cell(name, r, 1);
                Cell(new TextBlock { Text = VappRules.Of(vm)?.Name() ?? Messages.NONE, VerticalAlignment = VerticalAlignment.Center }, r, 2);
            }

            var all = new Button { Content = Dialogs.M("&Select All") };
            var none = new Button { Content = Dialogs.M("C&lear All") };
            all.Click += (_, _) => _rows.ForEach(x => x.Box.IsChecked = true);
            none.Click += (_, _) => _rows.ForEach(x => x.Box.IsChecked = false);
            var buttons = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            buttons.Children.Add(all);
            buttons.Children.Add(none);
            buttons.Children.Add(_count);

            var panel = new DockPanel();
            var top = new StackPanel { Spacing = 6 };
            top.Children.Add(Rubric(Messages.CHOOSE_VMS_VAPP_RUBRIC));
            top.Children.Add(new TextBlock { Text = string.Format(Helpers.IsPool(c) ? Messages.VMS_IN_POOL : Messages.VMS_IN_SERVER, Helpers.GetName(Helpers.GetPoolOfOne(c)).Ellipsise(60)) });
            DockPanel.SetDock(top, Dock.Top);
            DockPanel.SetDock(buttons, Dock.Bottom);
            buttons.Margin = new Thickness(0, 8, 0, 0);
            panel.Children.Add(top);
            panel.Children.Add(buttons);
            panel.Children.Add(new Border
            {
                BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Background = Brushes.White, Padding = new Thickness(6),
                MinHeight = 200, Child = new ScrollViewer { Content = grid, MaxHeight = 300 }
            });
            View = panel;
            Update();
        }

        public override string Title => Messages.NEWVMAPPLIANCE_VMSPAGE_TEXT;
        public override string Icon => "000_VM_h32bit_16";
        public override string SubText => Selected.Count == 1 ? Messages.ONE_VM_SELECTED : string.Format(Messages.MOREONE_VM_SELECTED, Selected.Count);

        public List<VM> Selected => _rows.Where(x => x.Box.IsChecked == true).Select(x => x.Vm).ToList();

        /// <summary>The vApp's members after the edit: the ticked VMs and the members the list does not show.</summary>
        private List<VM> Members => Selected.Concat(_unlisted).ToList();

        public override bool HasChanged =>
            _vapp != null && !new HashSet<string>(_vapp.VMs.Select(r => r.opaque_ref)).SetEquals(Members.Select(v => v.opaque_ref));

        /// <summary>NewVMGroupVMsPage.userAcceptsWarning: VMs taken from another vApp need the user's OK.</summary>
        public override async System.Threading.Tasks.Task<bool> ConfirmSave(Window owner) => await Confirm(owner, _vapp?.Name());

        public async System.Threading.Tasks.Task<bool> Confirm(Window owner, string vappName) =>
            !(VappRules.ChangeWarning(Selected, _vapp, vappName ?? "") is string warning)
            || await Dialogs.Confirm(owner, Messages.CHANGE_VMS_APPLIANCE, warning, Messages.YES, Messages.NO);

        public override AsyncAction SaveSettings() => VappRules.SetMembersAction(_vapp, Members);

        private void Update()
        {
            _count.Text = SubText;
            OnChanged();
        }
    }

    /// <summary>NewVMApplianceVMOrderAndDelaysPage: each VM's start order and the delay before the next VM starts.</summary>
    public sealed class VappOrderPage : EditPage
    {
        private readonly IXenConnection _connection;
        private readonly Func<IEnumerable<VM>> _vms;
        private readonly Dictionary<string, (long Order, long Delay)> _values = new Dictionary<string, (long, long)>();
        private readonly ContentControl _host = new ContentControl();

        public VappOrderPage(IXenConnection c, Func<IEnumerable<VM>> vms)
        {
            _connection = c;
            _vms = vms;
            var panel = new StackPanel { Spacing = 6 };
            panel.Children.Add(Rubric("Specify the order in which the VMs in the vApp are started, and the delay before the next VM is started."));
            panel.Children.Add(_host);
            View = panel;
            Refresh();
        }

        public override string Title => Messages.NEWVMAPPLIANCE_VMORDERANDDELAYSPAGE_TEXT;
        public override string Icon => "000_RebootVM_h32bit_16";

        /// <summary>Rebuilds the rows for the VMs chosen now, keeping the values already entered.</summary>
        public void Refresh()
        {
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("*,Auto,Auto"), ColumnSpacing = 10, RowSpacing = 4 };
            grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            string[] headers = { "VM", "Start order", "Delay interval (seconds)" };
            for (var i = 0; i < headers.Length; i++)
            {
                var h = new TextBlock { Text = headers[i], FontWeight = FontWeight.SemiBold };
                Grid.SetColumn(h, i);
                grid.Children.Add(h);
            }
            foreach (var vm in _vms().OrderBy(v => v))
            {
                var key = vm.opaque_ref;
                var (order, delay) = _values.TryGetValue(key, out var v) ? v : (vm.order, vm.start_delay);
                var orderBox = new NumericUpDown { Minimum = 0, Maximum = long.MaxValue, Value = order, FormatString = "0", Width = 130 };
                var delayBox = new NumericUpDown { Minimum = 0, Maximum = long.MaxValue, Value = delay, FormatString = "0", Width = 130 };
                void Save() { _values[key] = ((long)(orderBox.Value ?? 0), (long)(delayBox.Value ?? 0)); OnChanged(); }
                orderBox.ValueChanged += (_, _) => Save();
                delayBox.ValueChanged += (_, _) => Save();
                var r = grid.RowDefinitions.Count;
                grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                Control[] cells = { new TextBlock { Text = vm.Name(), VerticalAlignment = VerticalAlignment.Center }, orderBox, delayBox };
                for (var c = 0; c < cells.Length; c++) { Grid.SetRow(cells[c], r); Grid.SetColumn(cells[c], c); grid.Children.Add(cells[c]); }
            }
            _host.Content = grid.RowDefinitions.Count > 1 ? new ScrollViewer { Content = grid, MaxHeight = 320 } : new TextBlock { Text = "No VMs are selected." };
        }

        private Dictionary<VM, VMStartupOptions> Changed() =>
            VappRules.ChangedStartup(_vms().Select(vm => _values.TryGetValue(vm.opaque_ref, out var v) ? (vm, v.Order, v.Delay) : (vm, vm.order, vm.start_delay)));

        public override bool HasChanged => Changed().Count > 0;

        public override AsyncAction SaveSettings() => VappRules.StartupAction(_connection, Changed());
    }

    /// <summary>NewVMApplianceWizard: Name, Virtual Machines, VM Startup Sequence, Finish.</summary>
    public class VappWizard : WizardWindow
    {
        private readonly IXenConnection _connection;
        private readonly Action<AsyncAction> _run;
        private readonly TextBox _name = new TextBox();
        private readonly TextBox _description = new TextBox { AcceptsReturn = true, Height = 80, TextWrapping = TextWrapping.Wrap };
        private readonly VappVmsPage _vms;
        private readonly VappOrderPage _order;
        private Control _namePage;

        public VappWizard(IXenConnection c, Action<AsyncAction> run, IEnumerable<VM> preselected = null)
            : base(string.Format("New vApp - {0}", Helpers.GetName(Helpers.GetPoolOfOne(c))), "000_NewVirtualAppliance_h32bit_32",
                Messages.NEWVMAPPLIANCE_NAMEPAGE_TEXT, Messages.NEWVMAPPLIANCE_VMSPAGE_TEXT, Messages.NEWVMAPPLIANCE_VMORDERANDDELAYSPAGE_TEXT, Messages.NEWVMAPPLIANCE_FINISHPAGE_TEXT)
        {
            _connection = c;
            _run = run;
            _vms = new VappVmsPage(c, null, preselected);
            _order = new VappOrderPage(c, () => _vms.Selected);
        }

        protected override string PageTitle(int step) => step switch
        {
            0 => Messages.NEWVMAPPLIANCE_NAMEPAGE_TITLE,
            1 => Messages.NEWVMAPPLIANCE_VMSPAGE_TITLE,
            2 => Messages.NEWVMAPPLIANCE_VMORDERANDDELAYSPAGE_TITLE,
            _ => Messages.NEWVMAPPLIANCE_FINISHPAGE_TITLE
        };

        protected override Control Page(int step)
        {
            switch (step)
            {
                case 0:
                    if (_namePage != null) return _namePage;
                    var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowDefinitions = new RowDefinitions("Auto,Auto"), RowSpacing = 8, ColumnSpacing = 10 };
                    var nameLabel = new Label { Content = Dialogs.M("Na&me:"), Target = _name, VerticalAlignment = VerticalAlignment.Center };
                    var descLabel = new Label { Content = Dialogs.M("&Description:"), Target = _description };
                    Grid.SetColumn(_name, 1);
                    Grid.SetRow(descLabel, 1);
                    Grid.SetRow(_description, 1);
                    Grid.SetColumn(_description, 1);
                    grid.Children.AddRange(new Control[] { nameLabel, _name, descLabel, _description });
                    var p = new StackPanel { Spacing = 10 };
                    p.Children.Add(new TextBlock { Text = "Enter a name and an optional description for the new vApp.", TextWrapping = TextWrapping.Wrap });
                    p.Children.Add(grid);
                    Dispatcher.UIThread.Post(() => _name.Focus());
                    return _namePage = p;
                case 1:
                    return _vms.View;
                case 2:
                    _order.Refresh();
                    return _order.View;
                default:
                    var summary = new TextBox { Text = VappRules.Summary(_name.Text?.Trim(), _vms.Selected), IsReadOnly = true, TextWrapping = TextWrapping.Wrap, Height = 200 };
                    var f = new StackPanel { Spacing = 10 };
                    f.Children.Add(new TextBlock { Text = "The new vApp will be created with the following settings:" });
                    f.Children.Add(summary);
                    return f;
            }
        }

        protected override string Validate(int step) => step == 0 ? VappRules.NameError(_name.Text) : null;

        protected override async Task Finish()
        {
            if (VappRules.NameError(_name.Text) != null || !await _vms.Confirm(this, _name.Text.Trim())) return;
            var create = VappRules.CreateAction(_connection, _name.Text, _description.Text, _vms.Selected);
            var startup = _order.SaveSettings();
            if (startup != null)
                create.Completed += a => { if (a.Succeeded) Dispatcher.UIThread.Post(() => _run(startup)); };
            _run(create);
            Close();
        }
    }

    /// <summary>
    /// VMAppliancesDialog (Manage vApps): the pool's vApps with New vApp, Delete, Properties, Start and Shut down, and
    /// the VMs of the selected one. Import and Export (OVF) are not ported. shortcut: one vApp is selected at a time
    /// (WinForms deletes several at once); allow a multiple selection if it is missed.
    /// </summary>
    public class VappsDialog : Window
    {
        private readonly IXenConnection _connection;
        private readonly Action<AsyncAction> _run;
        private readonly Func<IXenObject, System.Threading.Tasks.Task> _properties;
        private readonly SelectableGrid _grid = new SelectableGrid { MinHeight = 220 };
        private readonly TextBlock _members = new TextBlock { TextWrapping = TextWrapping.Wrap };
        private readonly Button _delete = new Button { Content = Dialogs.M("&Delete") };
        private readonly Button _edit = new Button { Content = Dialogs.M("&Properties") };
        private readonly Button _start = new Button { Content = Dialogs.M("&Start") };
        private readonly Button _shutdown = new Button { Content = Dialogs.M("S&hut down") };
        private readonly DispatcherTimer _timer = new DispatcherTimer { Interval = TimeSpan.FromSeconds(1) };
        private string _signature;

        public VappsDialog(IXenConnection c, Action<AsyncAction> run, Func<IXenObject, System.Threading.Tasks.Task> properties)
        {
            _connection = c;
            _run = run;
            _properties = properties;
            var pool = Helpers.GetPoolOfOne(c);
            Title = string.Format(Messages.VM_APPLIANCES_TITLE, pool.Name());
            Width = 760;
            Height = 520;
            WindowStartupLocation = WindowStartupLocation.CenterOwner;
            ShowInTaskbar = false;

            var add = new Button { Content = Dialogs.M("&New vApp...") };
            add.Click += (_, _) => new VappWizard(_connection, _run).Show(this);
            _delete.Click += async (_, _) => await Delete();
            _edit.Click += async (_, _) => { if (Current is VM_appliance a) await _properties(a); };
            _start.Click += (_, _) => { if (VappRules.CanStart(Current)) _run(VappRules.StartAction(Current)); };
            _shutdown.Click += async (_, _) =>
            {
                var target = Current;
                if (!VappRules.CanShutDown(target)
                    || !await Dialogs.Confirm(this, Messages.VM_APPLIANCE_SHUT_DOWN, Messages.CONFIRM_SHUT_DOWN_APPLIANCE, Messages.YES, Messages.NO))
                    return;
                // the vApp chosen before the question, as it is now (it may have gone meanwhile)
                var now = _connection.Resolve(new XenRef<VM_appliance>(target.opaque_ref));
                if (VappRules.CanShutDown(now))
                    _run(VappRules.ShutDownAction(now));
            };
            var importButton = new Button { Content = Dialogs.M("&Import"), IsEnabled = false };
            var exportButton = new Button { Content = Dialogs.M("&Export"), IsEnabled = false };
            ToolTip.SetTip(importButton, "Importing a vApp (OVF/OVA) is not available yet.");
            ToolTip.SetTip(exportButton, "Exporting a vApp (OVF/OVA) is not available yet.");
            var toolbar = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
            toolbar.Children.AddRange(new Control[] { add, _delete, _edit, _start, _shutdown, importButton, exportButton });
            _grid.SelectionChanged += Buttons;
            _grid.Activated += async () => { if (Current is VM_appliance a) await _properties(a); };

            var close = new Button { Content = Messages.CLOSE, IsCancel = true, MinWidth = 75, HorizontalAlignment = HorizontalAlignment.Right };
            close.Click += (_, _) => Close();
            var membersBox = new StackPanel { Spacing = 4 };
            membersBox.Children.Add(new TextBlock { Text = "VMs in this vApp:", FontWeight = FontWeight.SemiBold });
            membersBox.Children.Add(_members);

            var root = new DockPanel { Margin = new Thickness(12) };
            var top = new StackPanel { Spacing = 8 };
            top.Children.Add(new TextBlock { Text = string.Format(Helpers.IsPool(c) ? Messages.VM_APPLIANCES_DEFINED_FOR_POOL : Messages.VM_APPLIANCES_DEFINED_FOR_SERVER, pool.Name()) });
            top.Children.Add(toolbar);
            DockPanel.SetDock(top, Dock.Top);
            DockPanel.SetDock(close, Dock.Bottom);
            DockPanel.SetDock(membersBox, Dock.Bottom);
            membersBox.Margin = new Thickness(0, 10);
            _grid.Margin = new Thickness(0, 8, 0, 0);
            root.Children.AddRange(new Control[] { top, close, membersBox, _grid });
            Content = root;

            _timer.Tick += (_, _) => Fill();
            Opened += (_, _) => { Fill(); _timer.Start(); };
            Closed += (_, _) => _timer.Stop();
        }

        private VM_appliance Current => _grid.Selected == null ? null : _connection.Resolve(new XenRef<VM_appliance>(_grid.Selected));

        /// <summary>VMAppliancesDialog refreshes on cache changes; this polls the cache once a second.</summary>
        private void Fill()
        {
            var apps = VappRules.Appliances(_connection);
            var signature = string.Join("|", apps.Select(a => $"{a.opaque_ref}:{a.Name()}:{a.Description()}:{a.VMs.Count}:{string.Join(",", a.allowed_operations)}"))
                + "|" + string.Join(",", apps.SelectMany(a => _connection.ResolveAll(a.VMs)).Select(v => v.Name() + v.power_state));
            if (signature != _signature)
            {
                _signature = signature;
                var table = new Table { Headers = new[] { "Name", "Description", "#VMs" }, Fill = 1 };
                foreach (var a in apps)
                    table.Rows.Add(new[] { a.Name(), a.Description(), a.VMs.Count.ToString() });
                _grid.Fill(table, apps.Select(a => a.opaque_ref).ToList(), _grid.Selected);
            }
            Buttons();
        }

        private void Buttons()
        {
            var a = Current;
            _edit.IsEnabled = _delete.IsEnabled = a != null;
            _start.IsEnabled = VappRules.CanStart(a);
            _shutdown.IsEnabled = VappRules.CanShutDown(a);
            _members.Text = a == null ? "" : a.VMs.Count == 0 ? Messages.NONE
                : string.Join(", ", _connection.ResolveAll(a.VMs).OrderBy(v => v).Select(v => $"{v.Name()} ({v.power_state})"));
        }

        private async Task Delete()
        {
            if (!(Current is VM_appliance a)) return;
            if (!await Dialogs.Confirm(this, Messages.DELETE_VM_APPLIANCE_TITLE, VappRules.DeleteConfirmation(new[] { a }), Messages.YES, Messages.NO))
                return;
            if (_connection.Resolve(new XenRef<VM_appliance>(a.opaque_ref)) is VM_appliance now)
                _run(VappRules.DeleteAction(_connection, new List<VM_appliance> { now }));
        }
    }
}
