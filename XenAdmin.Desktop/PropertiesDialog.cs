using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using XenAdmin.Actions;
using XenAdmin.Actions.VMActions;
using XenAdmin.Core;
using XenAdmin.Model;
using XenAPI;
using XenCenterLib;

namespace XenAdmin.Desktop
{
    /// <summary>One page of the Properties dialog (the WinForms IEditPage contract).</summary>
    public abstract class EditPage
    {
        public abstract string Title { get; }
        public abstract string Icon { get; }
        public virtual string SubText => "";
        public virtual bool HasChanged => false;
        public virtual string ValidationError => null;
        public Control View { get; protected set; }
        public event Action Changed;
        protected void OnChanged() => Changed?.Invoke();
        /// <summary>IsValidToSave with a question (the page may warn and let the user cancel the save).</summary>
        public virtual System.Threading.Tasks.Task<bool> ConfirmSave(Window owner) => System.Threading.Tasks.Task.FromResult(true);

        /// <summary>Applies edits to the cloned object; returns an extra action to run, or null.</summary>
        public virtual AsyncAction SaveSettings() => null;

        protected static TextBlock Rubric(string text) =>
            new TextBlock { Text = text, TextWrapping = TextWrapping.Wrap, Margin = new Thickness(0, 0, 0, 14) };
    }

    /// <summary>Page shown in the list but not ported yet (disabled, as with the main window tabs).</summary>
    public sealed class NotPortedPage : EditPage
    {
        private readonly string _title, _icon;
        public NotPortedPage(string title, string icon) { _title = title; _icon = icon; }
        public override string Title => _title;
        public override string Icon => _icon;
    }

    /// <summary>GeneralEditPage: name, description, folder, tags (and the iSCSI IQN for servers).</summary>
    public sealed class GeneralPage : EditPage
    {
        private readonly IXenObject _orig, _copy;
        private readonly TextBox _name = new TextBox();
        private readonly TextBox _description = new TextBox { AcceptsReturn = true, TextWrapping = TextWrapping.Wrap, Height = 50 };
        private readonly TextBlock _folderText = new TextBlock { VerticalAlignment = VerticalAlignment.Center };
        private readonly WrapPanel _tagPanel = new WrapPanel();
        private string _folder;
        private List<string> _tags;
        private readonly TextBox _iqn = new TextBox();
        private readonly bool _saveDescription;

        public override string Title => Messages.NAME_DESCRIPTION_TAGS;
        public override string Icon => "edit_16";
        public override string SubText => _name.Text;

        public GeneralPage(IXenObject orig, IXenObject copy)
        {
            _orig = orig;
            _copy = copy;
            _name.Text = Helpers.GetName(copy);
            _folder = copy.Path ?? "";
            _tags = OriginalTags;
            var vm = copy as VM;
            var descriptionType = vm?.DescriptionType() ?? VM.VmDescriptionType.ReadWrite;
            _saveDescription = descriptionType == VM.VmDescriptionType.ReadWrite;
            _description.Text = copy.Description();
            _description.IsReadOnly = !_saveDescription;

            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowSpacing = 10, ColumnSpacing = 18 };
            int row = 0;
            void Row(string label, Control c)
            {
                grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var l = new Label { Content = Dialogs.M(label), Target = c, VerticalAlignment = VerticalAlignment.Top, Padding = new Thickness(0, 6, 0, 0) };
                Grid.SetRow(l, row); Grid.SetRow(c, row); Grid.SetColumn(c, 1);
                grid.Children.Add(l); grid.Children.Add(c);
                row++;
            }
            Row("Na&me:", _name);
            if (descriptionType != VM.VmDescriptionType.None)
                Row("&Description:", _description);
            // GeneralEditPage: vApps have no folder or tags
            if (!(copy is VM_appliance))
            {
                Row("Folder:", FolderBox());
                Row("Tags:", TagsBox());
                RefreshFolder();
                RefreshTags();
            }
            if (copy is Host host)
            {
                _iqn.Text = host.GetIscsiIqn();
                Row("iSCSI &IQN:", _iqn);
            }
            var panel = new StackPanel();
            panel.Children.Add(Rubric("XCP-ng Center provides several methods for labelling and organizing your resources. By using features such as folders and tags you can quickly group and manage resources across your connected servers."));
            panel.Children.Add(grid);
            View = panel;
            foreach (var t in new[] { _name, _description, _iqn })
                t.TextChanged += (_, _) => OnChanged();
        }

        // FolderEditor: "<None>  Change..." in a field-like box
        private Control FolderBox()
        {
            var link = Views.LinkText("Change...", async () =>
            {
                var owner = TopLevel.GetTopLevel(View) as Window;
                var chosen = await Dialogs.ChangeFolder(owner, _copy.Connection, _folder);
                if (chosen == null) return;
                _folder = chosen;
                RefreshFolder();
                OnChanged();
            });
            link.TextDecorations = TextDecorations.Underline;
            link.VerticalAlignment = VerticalAlignment.Center;
            var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 28 };
            row.Children.Add(_folderText);
            row.Children.Add(link);
            return new Border { BorderBrush = Brush.Parse("#7A7A7A"), BorderThickness = new Thickness(1), Background = Brushes.White, Padding = new Thickness(4, 2), MinHeight = 23, Child = row };
        }

        // TagsEditor: tag list plus "Edit tags..." link in a tall box
        private Control TagsBox()
        {
            var link = Views.LinkText(Messages.NEW_TAG_LINK, async () =>
            {
                var owner = TopLevel.GetTopLevel(View) as Window;
                var chosen = await Dialogs.EditTags(owner, _copy.Connection, _tags);
                if (chosen == null) return;
                _tags = chosen.OrderBy(t => t).ToList();
                RefreshTags();
                OnChanged();
            });
            link.TextDecorations = TextDecorations.Underline;
            var linkRow = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 4, Margin = new Thickness(0, 8, 0, 0) };
            linkRow.Children.Add(Views.Icon("000_Tag_h32bit_16"));
            linkRow.Children.Add(link);
            var content = new StackPanel();
            content.Children.Add(_tagPanel);
            content.Children.Add(linkRow);
            return new Border { BorderBrush = Brush.Parse("#7A7A7A"), BorderThickness = new Thickness(1), Background = Brushes.White, Padding = new Thickness(4, 3), Height = 262, Child = content };
        }

        private void RefreshFolder() => _folderText.Text = string.IsNullOrEmpty(_folder) ? Messages.NONE : _folder;

        private void RefreshTags()
        {
            _tagPanel.Children.Clear();
            if (_tags.Count == 0)
                _tagPanel.Children.Add(new TextBlock { Text = Messages.NONE });
            foreach (var t in _tags)
                _tagPanel.Children.Add(new Border
                {
                    Background = Brush.Parse("#F0F0F0"), BorderBrush = Brush.Parse("#BFBFBF"), BorderThickness = new Thickness(1),
                    Padding = new Thickness(4, 0), Margin = new Thickness(0, 0, 4, 4), Child = new TextBlock { Text = t }
                });
        }

        private List<string> OriginalTags => (Tags.GetTags(_copy) ?? Array.Empty<string>()).OrderBy(t => t).ToList();
        private List<string> NewTags => _tags;
        // folders are paths like "/Production/Web"; empty means no folder
        private string NewFolder => string.IsNullOrWhiteSpace(_folder) ? "" : "/" + _folder.Trim().Trim('/');
        private bool FolderChanged => NewFolder != (_copy.Path ?? "");
        private bool TagsChanged => !OriginalTags.SequenceEqual(NewTags);
        private bool IqnChanged => _copy is Host h && h.GetIscsiIqn() != _iqn.Text.Trim();

        public override bool HasChanged =>
            Helpers.GetName(_copy) != _name.Text.Trim() || _saveDescription && _copy.Description() != _description.Text
            || FolderChanged || TagsChanged || IqnChanged;

        public override string ValidationError =>
            _name.Text.Trim().Length == 0 ? Messages.GENERAL_EDIT_INVALID_NAME
            : IqnChanged && !Helpers.ValidateIscsiIQN(_iqn.Text.Trim()) ? Messages.GENERAL_EDIT_INVALID_IQN
            : null;

        public override AsyncAction SaveSettings()
        {
            var name = _name.Text.Trim();
            if (name != "" && name != Helpers.GetName(_copy))
                _copy.Set("name_label", name);
            if (_saveDescription && _description.Text != _copy.Description())
                _copy.Set("name_description", _description.Text);
            if (_copy is Host host && IqnChanged)
                host.SetIscsiIqn(_iqn.Text.Trim());
            return FolderChanged || TagsChanged ? new GeneralEditPageAction(_orig, _copy, NewFolder, NewTags, true) : null;
        }
    }

    /// <summary>CpuMemoryEditPage: vCPUs, topology and vCPU priority.</summary>
    public sealed class CpuPage : EditPage
    {
        private readonly VM _vm;
        private readonly bool _hotplug;
        private readonly long _origVcpus, _origAtStartup;
        private readonly int _origPriority;
        private readonly ComboBox _vcpus = new ComboBox { MinWidth = 60 };
        private readonly ComboBox _initial = new ComboBox { MinWidth = 60 };
        private readonly ComboBox _topology = new ComboBox { MinWidth = 260 };
        private readonly Slider _priority = new Slider { Minimum = 0, Maximum = 8, TickFrequency = 1, IsSnapToTickEnabled = true, TickPlacement = Avalonia.Controls.TickPlacement.BottomRight, Width = 520 };
        private readonly TextBlock _warning = new TextBlock { Foreground = Brushes.DarkRed, TextWrapping = TextWrapping.Wrap };

        public override string Title => Messages.CPU;
        public override string Icon => "000_CPU_h32bit_16";
        public override string SubText => string.Format(Messages.CPU_SUB, SelectedAtStartup);

        public CpuPage(VM vm)
        {
            _vm = vm;
            _hotplug = vm.SupportsVcpuHotplug();
            _origVcpus = _hotplug ? Math.Max(vm.VCPUs_max, 1) : Math.Max(vm.VCPUs_at_startup, 1);
            _origAtStartup = Math.Max(vm.VCPUs_at_startup, 1);
            _origPriority = (int)Math.Round(Math.Log(vm.GetVcpuWeight()) / Math.Log(4.0));
            var halted = vm.power_state == vm_power_state.Halted;

            var maxAllowed = Math.Max(vm.MaxVCPUsAllowed(), _origVcpus);
            _vcpus.ItemsSource = Enumerable.Range(1, (int)maxAllowed).Select(i => (long)i)
                .Where(i => i == _origVcpus || Topologies(i).Count > 0).ToList();
            _vcpus.SelectedItem = _origVcpus;
            _vcpus.IsEnabled = _topology.IsEnabled = _priority.IsEnabled = halted;
            _initial.IsEnabled = _hotplug && (halted || vm.power_state == vm_power_state.Running);
            _priority.Value = _origPriority;
            FillInitial();
            FillTopology(vm.GetCoresPerSocket());

            var panel = new StackPanel();
            panel.Children.Add(Rubric(Messages.VM_CPUMEMPAGE_RUBRIC + (_hotplug ? Messages.VM_CPUMEMPAGE_RUBRIC_HOTPLUG : "")));
            if (!halted)
            {
                var info = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6, Margin = new Thickness(0, 0, 0, 14) };
                info.Children.Add(Views.Icon("000_Info3_h32bit_16"));
                info.Children.Add(new TextBlock { Text = _hotplug ? Messages.VM_CPUMEMPAGE_MAX_VCPUS_READONLY : Messages.VCPU_ONLY_WHEN_HALTED, TextWrapping = TextWrapping.Wrap, MaxWidth = 560 });
                panel.Children.Add(info);
            }
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,Auto"), RowSpacing = 8, ColumnSpacing = 12 };
            int row = 0;
            void Row(string label, Control c)
            {
                grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var l = new Label { Content = Dialogs.M(label), Target = c, VerticalAlignment = VerticalAlignment.Center };
                Grid.SetRow(l, row); Grid.SetRow(c, row); Grid.SetColumn(c, 1);
                grid.Children.Add(l); grid.Children.Add(c);
                row++;
            }
            Row(_hotplug ? Messages.VM_CPUMEMPAGE_MAX_VCPUS_LABEL : Messages.VM_CPUMEMPAGE_VCPUS_LABEL, _vcpus);
            Row("&Topology:", _topology);
            if (_hotplug)
                Row(halted ? Messages.VM_CPUMEMPAGE_INITIAL_VCPUS_LABEL : Messages.VM_CPUMEMPAGE_CURRENT_VCPUS_LABEL, _initial);
            panel.Children.Add(grid);
            panel.Children.Add(new TextBlock { Text = "vCPU priority for this virtual machine:", Margin = new Thickness(0, 20, 0, 4) });
            panel.Children.Add(_priority);
            var scale = new Grid { ColumnDefinitions = new ColumnDefinitions("*,*,*"), Width = 520, HorizontalAlignment = HorizontalAlignment.Left };
            void Scale(string t, int col, HorizontalAlignment a) { var x = new TextBlock { Text = t, HorizontalAlignment = a, Opacity = 0.7 }; Grid.SetColumn(x, col); scale.Children.Add(x); }
            Scale("Lowest", 0, HorizontalAlignment.Left);
            Scale("Normal", 1, HorizontalAlignment.Center);
            Scale("Highest", 2, HorizontalAlignment.Right);
            panel.Children.Add(scale);
            panel.Children.Add(_warning);
            View = panel;

            _vcpus.SelectionChanged += (_, _) =>
            {
                FillTopology((_topology.SelectedItem as Topology)?.Cores ?? _vm.GetCoresPerSocket());
                FillInitial();
                Validate();
                OnChanged();
            };
            _initial.SelectionChanged += (_, _) => { Validate(); OnChanged(); };
            _topology.SelectionChanged += (_, _) => { Validate(); OnChanged(); };
            _priority.PropertyChanged += (_, e) => { if (e.Property == Avalonia.Controls.Primitives.RangeBase.ValueProperty) OnChanged(); };
            Validate();
        }

        private sealed class Topology
        {
            public long Sockets, Cores;
            public override string ToString() => VM.GetTopology(Sockets, Cores);
        }


        private List<Topology> Topologies(long vcpus) =>
            VmRules.ValidTopologies(vcpus, _vm.MaxCoresPerSocket()).Select(t => new Topology { Sockets = t.Sockets, Cores = t.Cores }).ToList();

        private long SelectedVcpus => _vcpus.SelectedItem is long l ? l : _origVcpus;
        private long SelectedAtStartup => _hotplug && _initial.SelectedItem is long l ? l : SelectedVcpus;
        private long SelectedCores => (_topology.SelectedItem as Topology)?.Cores ?? _vm.GetCoresPerSocket();

        private void FillTopology(long keepCores)
        {
            var list = Topologies(SelectedVcpus);
            if (list.All(t => t.Cores != keepCores) && SelectedVcpus == _origVcpus)
                list.Add(new Topology { Sockets = 0, Cores = keepCores }); // invalid current value stays selectable, as in WinForms
            _topology.ItemsSource = list;
            _topology.SelectedItem = list.FirstOrDefault(t => t.Cores == keepCores) ?? list.FirstOrDefault();
        }

        private void FillInitial()
        {
            if (!_hotplug) return;
            var current = _initial.SelectedItem is long l ? Math.Min(l, SelectedVcpus) : _origAtStartup;
            var min = _vm.power_state == vm_power_state.Halted ? 1 : _origAtStartup;
            _initial.ItemsSource = Enumerable.Range((int)min, (int)Math.Max(0, SelectedVcpus - min + 1)).Select(i => (long)i).ToList();
            _initial.SelectedItem = current;
        }

        private void Validate()
        {
            var warnings = new List<string>();
            var minVcpus = _vm.MinVCPUs();
            if (SelectedVcpus < minVcpus || SelectedAtStartup < minVcpus)
                warnings.Add(string.Format(Messages.VM_CPUMEMPAGE_VCPU_MIN_WARNING, minVcpus));
            var maxPhysical = _vm.Connection.Cache.Hosts.Select(h => h.host_CPUs.Count).DefaultIfEmpty(0).Max();
            if (maxPhysical > 0 && maxPhysical < SelectedVcpus)
                warnings.Add(Messages.VM_CPUMEMPAGE_VCPU_WARNING);
            var topologyWarning = VM.ValidVCPUConfiguration(SelectedVcpus, SelectedCores);
            if (!string.IsNullOrEmpty(topologyWarning))
                warnings.Add(topologyWarning + ".");
            _warning.Text = string.Join("\n\n", warnings);
            _warning.Margin = new Thickness(0, warnings.Count > 0 ? 16 : 0, 0, 0);
        }

        private bool VcpusChanged => _origVcpus != SelectedVcpus || _hotplug && _origAtStartup != SelectedAtStartup;
        private bool TopologyChanged => _vm.GetCoresPerSocket() != SelectedCores;
        private bool PriorityChanged => (int)_priority.Value != _origPriority;

        public override bool HasChanged => VcpusChanged || TopologyChanged || PriorityChanged;

        public override AsyncAction SaveSettings()
        {
            if (PriorityChanged)
                _vm.SetVcpuWeight((int)Math.Pow(4.0, _priority.Value));
            if (TopologyChanged)
                _vm.SetCoresPerSocket(SelectedCores);
            return VcpusChanged ? new ChangeVCPUSettingsAction(_vm, SelectedVcpus, SelectedAtStartup) : null;
        }
    }

    /// <summary>BootOptionsEditPage: autostart, HVM boot order, PV boot device and parameters.</summary>
    /// <summary>HomeServerEditPage with its AffinityPicker: no home server (shared storage) or one of the live servers.</summary>
    public sealed class HomeServerPage : EditPage
    {
        private readonly VM _vm;
        private readonly Host _orig;
        private readonly RadioButton _none = new RadioButton { GroupName = "home" };
        private readonly RadioButton _server = new RadioButton { Content = Dialogs.M("Place the VM on this &server:"), GroupName = "home" };
        private readonly ListBox _hosts = new ListBox { Height = 150, BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Margin = new Thickness(20, 4, 0, 0) };

        public override string Title => Messages.HOME_SERVER;
        public override string Icon => "000_ServerHome_h32bit_16";

        public HomeServerPage(VM vm)
        {
            _vm = vm;
            _orig = vm.Connection.Resolve(vm.affinity);
            // AutoHeightRadioButton: the long text wraps
            _none.Content = new TextBlock { Text = (Helpers.HasFullyConnectedSharedStorage(vm.Connection)
                ? Messages.AFFINITY_PICKER_DYNAMIC_SHARED_SR : Messages.AFFINITY_PICKER_DYNAMIC_NOT_SHARED_SR).Replace("&", ""), TextWrapping = TextWrapping.Wrap };
            _none.IsEnabled = VmRules.CanHaveNoHomeServer(vm);

            var rows = new List<Control>();
            foreach (var host in vm.Connection.Cache.Hosts.OrderBy(h => h))
            {
                var live = host.IsLive();
                var row = new Grid { ColumnDefinitions = new ColumnDefinitions("22,*,Auto"), Margin = new Thickness(0, 1) };
                row.Children.Add(Views.Icon(IconNames.For(host), 16));
                var name = new TextBlock { Text = host.Name(), Margin = new Thickness(4, 0) };
                var reason = new TextBlock { Text = live ? host.HostMemoryString() : Messages.HOMESERVER_PICKER_HOST_NOT_LIVE, Foreground = Brushes.Gray };
                Grid.SetColumn(name, 1);
                Grid.SetColumn(reason, 2);
                row.Children.Add(name);
                row.Children.Add(reason);
                // the item itself is disabled, so neither the mouse nor the keyboard can pick a server that is down
                rows.Add(new ListBoxItem { Content = row, Tag = host, IsEnabled = live });
            }
            _hosts.ItemsSource = rows;
            // SelectRadioButtons: the home server when it can be selected, else "no home server" when allowed
            var current = rows.FirstOrDefault(r => r.IsEnabled && ((Host)r.Tag).opaque_ref == _orig?.opaque_ref);
            if (current == null && _none.IsEnabled)
                _none.IsChecked = true;
            else
            {
                _server.IsChecked = true;
                _hosts.SelectedItem = current;
            }
            _hosts.IsEnabled = _server.IsChecked == true;

            _server.IsCheckedChanged += (_, _) =>
            {
                _hosts.IsEnabled = _server.IsChecked == true;
                if (_hosts.IsEnabled && _hosts.SelectedItem == null)
                    _hosts.SelectedItem = rows.FirstOrDefault(r => r.IsEnabled);
                OnChanged();
            };
            _hosts.SelectionChanged += (_, _) => OnChanged();

            var panel = new StackPanel { Spacing = 6 };
            panel.Children.Add(Rubric("When you nominate a home server for a virtual machine, the virtual machine will always be started up on that server if it is available. If this is not possible, then an alternate server within the same pool will be selected automatically."));
            panel.Children.Add(_none);
            panel.Children.Add(_server);
            panel.Children.Add(_hosts);
            View = panel;
        }

        // _server, not _none: the group unchecks _none only after _server's IsCheckedChanged has run
        private Host Selected => _server.IsChecked == true ? (_hosts.SelectedItem as Control)?.Tag as Host : null;

        public override string SubText => Selected?.Name() ?? (_server.IsChecked == true ? _orig?.Name() : null) ?? Messages.NONE_DEFINED;
        // a home server that is down stays as it is until the user picks another one
        public override bool HasChanged => !(_server.IsChecked == true && Selected == null) && Selected?.opaque_ref != _orig?.opaque_ref;
        public override AsyncAction SaveSettings() => VmRules.SetHomeServer(_vm, Selected);
    }

    public sealed class BootPage : EditPage
    {
        private readonly VM _vm;
        private readonly CheckBox _autoBoot = new CheckBox { Content = "Autostart on server boot" };
        private readonly ListBox _order = new ListBox { Width = 145, Height = 95, BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Classes = { "checked" } };
        private readonly List<(char Code, CheckBox Box)> _devices = new List<(char, CheckBox)>();
        private readonly ComboBox _pvDevice = new ComboBox { MinWidth = 160 };
        private readonly TextBox _pvArgs = new TextBox { MinWidth = 320 };
        private readonly bool _origBootFromCd;

        public override string Title => Messages.GENERAL_HEADING_BOOT_OPTIONS;
        public override string Icon => "001_PowerOn_h32bit_16";

        public static string DeviceName(char c) =>
            c == 'C' ? Messages.BOOT_HARD_DISK : c == 'D' ? Messages.BOOT_DVD_DRIVE : Messages.NETWORK;

        public BootPage(VM vm)
        {
            _vm = vm;
            _autoBoot.IsChecked = vm.GetAutoPowerOn();
            _origBootFromCd = vm.Connection.ResolveAll(vm.VBDs).Any(b => b.IsCDROM() && b.bootable);

            var panel = new StackPanel();
            panel.Children.Add(Rubric(vm.IsHVM()
                ? "Choose the devices from which you would like to boot the VM’s operating system, and configure the boot device order. The first item in the list will be scanned first when booting; if boot media is not found, the second device will then be scanned, and so on down the list of selected devices."
                : "Choose the device from which you would like to boot the VM’s operating system, and, optionally, specify any additional boot parameters."));
            var autostart = new Border { BorderBrush = Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1), Padding = new Thickness(6, 4, 6, 18), Margin = new Thickness(0, 0, 0, 16) };
            var autostartPanel = new StackPanel();
            autostartPanel.Children.Add(new TextBlock { Text = "Autostart", Margin = new Thickness(0, 0, 0, 4) });
            autostartPanel.Children.Add(_autoBoot);
            autostart.Child = autostartPanel;
            panel.Children.Add(autostart);

            if (vm.IsHVM())
            {
                var order = vm.GetBootOrder().ToUpper();
                foreach (var c in order.Concat("CDN".Where(x => !order.Contains(x))))
                {
                    var box = new CheckBox { Content = DeviceName(c), IsChecked = order.Contains(c), MinHeight = 18, Padding = new Thickness(2, 0) };
                    box.IsCheckedChanged += (_, _) => OnChanged();
                    // the check box takes the click: select its row too, so Move Up/Down apply to it
                    box.AddHandler(Avalonia.Input.InputElement.PointerPressedEvent, (_, _) => _order.SelectedItem = box, Avalonia.Interactivity.RoutingStrategies.Tunnel);
                    box.GotFocus += (_, _) => _order.SelectedItem = box;
                    _devices.Add((c, box));
                }
                RefreshOrder();
                var up = new Button { Content = Dialogs.M("Move &Up"), MinWidth = 82, IsEnabled = false };
                var down = new Button { Content = Dialogs.M("Move &Down"), MinWidth = 82, IsEnabled = false };
                void Toggle() { up.IsEnabled = _order.SelectedIndex > 0; down.IsEnabled = _order.SelectedIndex >= 0 && _order.SelectedIndex < _devices.Count - 1; }
                _order.SelectionChanged += (_, _) => Toggle();
                up.Click += (_, _) => Move(-1);
                down.Click += (_, _) => Move(1);
                var buttons = new StackPanel { Spacing = 6 };
                buttons.Children.Add(up);
                buttons.Children.Add(down);
                var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
                row.Children.Add(new TextBlock { Text = "Boot order:", Width = 130 });
                row.Children.Add(_order);
                row.Children.Add(buttons);
                panel.Children.Add(row);
                // Convert to PV is shown but not ported yet
                panel.Children.Add(new Button { Content = "Convert to PV", IsEnabled = false, Width = 132, HorizontalAlignment = HorizontalAlignment.Left, Margin = new Thickness(0, 30, 0, 0) });
            }
            else
            {
                var items = new List<string> { Messages.BOOT_HARD_DISK };
                if (vm.HasCD())
                    items.Add(Messages.DVD_DRIVE);
                _pvDevice.ItemsSource = items;
                _pvDevice.SelectedItem = _origBootFromCd && vm.HasCD() ? Messages.DVD_DRIVE : Messages.BOOT_HARD_DISK;
                _pvArgs.Text = vm.PV_args;
                var g = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,Auto"), RowDefinitions = new RowDefinitions("Auto,Auto"), RowSpacing = 8, ColumnSpacing = 12 };
                var l1 = new Label { Content = Dialogs.M("&Boot device:"), Target = _pvDevice };
                var l2 = new Label { Content = Dialogs.M("&OS boot parameters:"), Target = _pvArgs };
                Grid.SetColumn(_pvDevice, 1); Grid.SetRow(l2, 1); Grid.SetRow(_pvArgs, 1); Grid.SetColumn(_pvArgs, 1);
                g.Children.Add(l1); g.Children.Add(_pvDevice); g.Children.Add(l2); g.Children.Add(_pvArgs);
                panel.Children.Add(g);
                _pvDevice.SelectionChanged += (_, _) => OnChanged();
                _pvArgs.TextChanged += (_, _) => OnChanged();
            }
            View = panel;
            _autoBoot.IsCheckedChanged += (_, _) => OnChanged();
        }

        private void RefreshOrder()
        {
            var selected = _order.SelectedIndex;
            _order.ItemsSource = _devices.Select(d => d.Box).ToList();
            _order.SelectedIndex = selected;
        }

        private void Move(int delta)
        {
            var i = _order.SelectedIndex;
            var j = i + delta;
            if (i < 0 || j < 0 || j >= _devices.Count) return;
            (_devices[i], _devices[j]) = (_devices[j], _devices[i]);
            _order.ItemsSource = null;
            RefreshOrder();
            _order.SelectedIndex = j;
            OnChanged();
        }

        private string NewOrder => new string(_devices.Where(d => d.Box.IsChecked == true).Select(d => d.Code).ToArray());
        private bool BootFromCd => (string)_pvDevice.SelectedItem == Messages.DVD_DRIVE;

        public override string SubText
        {
            get
            {
                if (!_vm.IsHVM())
                    return Messages.NONE_DEFINED;
                var order = string.Join(", ", _devices.Select(d => DeviceName(d.Code)));
                return string.Format(_autoBoot.IsChecked == true ? Messages.BOOTORDER_AUTOSTART : Messages.BOOTORDER, order);
            }
        }

        public override bool HasChanged =>
            _autoBoot.IsChecked != _vm.GetAutoPowerOn()
            || _vm.IsHVM() && NewOrder != _vm.GetBootOrder().ToUpper()
            || !_vm.IsHVM() && (_pvArgs.Text != _vm.PV_args || BootFromCd != _origBootFromCd);

        public override AsyncAction SaveSettings()
        {
            if (_vm.IsHVM())
                _vm.SetBootOrder(NewOrder);
            else
                _vm.PV_args = _pvArgs.Text;
            _vm.SetAutoPowerOn(_autoBoot.IsChecked == true);
            if (_vm.IsHVM())
                return null;
            return VmRules.PvBootableAction(_vm, BootFromCd);
        }
    }

    /// <summary>PropertiesDialog: vertical page list, teal page header, OK/Cancel; saves via SaveChangesAction + page actions.</summary>
    public class PropertiesDialog : Window
    {
        private readonly IXenObject _orig, _copy;
        private readonly List<EditPage> _pages;
        private readonly ContentControl _pageHost = new ContentControl();
        private readonly ContentControl _pageHeader = new ContentControl();
        private readonly ListBox _list = new ListBox();
        private readonly TextBlock _error = new TextBlock { Foreground = Brushes.Red, VerticalAlignment = VerticalAlignment.Center, TextWrapping = TextWrapping.Wrap };

        public AsyncAction Action { get; private set; }

        public static List<EditPage> PagesFor(IXenObject orig, IXenObject copy)
        {
            var pages = new List<EditPage> { new GeneralPage(orig, copy), new NotPortedPage(Messages.CUSTOM_FIELDS, "000_Fields_h32bit_16") };
            switch (copy)
            {
                case VM vm when !vm.is_a_snapshot:
                    pages.Add(new CpuPage(vm));
                    pages.Add(new BootPage(vm));
                    pages.Add(new NotPortedPage(Messages.START_UP_OPTIONS, "000_RebootVM_h32bit_16"));
                    pages.Add(new NotPortedPage(Messages.ALERTS, "000_Alert2_h32bit_16"));
                    if (!Helpers.WlbEnabledAndConfigured(vm.Connection)) // WLB places the VMs itself
                        pages.Add(new HomeServerPage(vm));
                    pages.Add(new NotPortedPage("USB", "usb_16"));
                    pages.Add(new NotPortedPage(Messages.ADVANCED_OPTIONS, "002_Configure_h32bit_16"));
                    break;
                case Host _:
                    pages.Add(new NotPortedPage(Messages.ALERTS, "000_Alert2_h32bit_16"));
                    pages.Add(new NotPortedPage(Messages.MULTIPATHING, "000_Storage_h32bit_16"));
                    pages.Add(new NotPortedPage(Messages.LOG_DESTINATION, "log_destination_16"));
                    break;
                case Pool _:
                case SR _:
                    pages.Add(new NotPortedPage(Messages.ALERTS, "000_Alert2_h32bit_16"));
                    break;
                case XenAPI.Network _:
                    pages.Add(new NotPortedPage("Network Settings", "000_Network_h32bit_16"));
                    break;
                case VM_appliance vapp:
                    // the vApp's VMs, and the startup sequence of the VMs chosen there
                    var vms = new VappVmsPage(vapp.Connection, vapp);
                    var order = new VappOrderPage(vapp.Connection, () => vms.Selected);
                    vms.Changed += order.Refresh;
                    pages.Add(vms);
                    pages.Add(order);
                    break;
            }
            return pages;
        }

        public PropertiesDialog(IXenObject xenObject)
        {
            _orig = xenObject;
            _copy = xenObject.Clone();
            _pages = PagesFor(_orig, _copy);
            Title = $"'{Helpers.GetName(xenObject)}' Properties";
            Width = 880;
            Height = 590;
            MinWidth = 700;
            MinHeight = 480;
            WindowStartupLocation = WindowStartupLocation.CenterOwner;
            ShowInTaskbar = false;
            Background = Brush.Parse("#F0F0F0");

            _list.ItemTemplate = new Avalonia.Controls.Templates.FuncDataTemplate<EditPage>((p, _) =>
            {
                var row = new Grid { ColumnDefinitions = new ColumnDefinitions("26,*"), Margin = new Thickness(0, 2), IsEnabled = !(p is NotPortedPage) };
                var icon = Views.Icon(p.Icon);
                icon.VerticalAlignment = VerticalAlignment.Top;
                icon.Margin = new Thickness(0, 2, 0, 0);
                var text = new StackPanel();
                text.Children.Add(new TextBlock { Text = p.Title });
                var sub = new TextBlock { Text = p.SubText, Opacity = 0.65, Margin = new Thickness(8, 0, 0, 0), TextTrimming = TextTrimming.CharacterEllipsis };
                p.Changed += () => sub.Text = p.SubText;
                text.Children.Add(sub);
                Grid.SetColumn(text, 1);
                row.Children.Add(icon);
                row.Children.Add(text);
                return row;
            });
            _list.ItemsSource = _pages;
            _list.Classes.Add("props");
            _list.SelectionChanged += (_, _) =>
            {
                if (_list.SelectedItem is EditPage p && p is NotPortedPage)
                {
                    _list.SelectedItem = _pages[0];
                    return;
                }
                ShowPage(_list.SelectedItem as EditPage);
            };

            var ok = new Button { Content = Messages.OK, IsDefault = true, MinWidth = 75 };
            var cancel = new Button { Content = Messages.CANCEL, IsCancel = true, MinWidth = 75 };
            ok.Click += async (_, _) => await Ok();
            cancel.Click += (_, _) => Close();
            // Esc closes from any field (a focused TextBox would otherwise keep the key)
            AddHandler(KeyDownEvent, (_, e) => { if (e.Key == Avalonia.Input.Key.Escape) { e.Handled = true; Close(); } }, Avalonia.Interactivity.RoutingStrategies.Tunnel);
            var buttons = new DockPanel { Margin = new Thickness(0, 10, 0, 0) };
            var right = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, HorizontalAlignment = HorizontalAlignment.Right };
            right.Children.Add(ok);
            right.Children.Add(cancel);
            DockPanel.SetDock(right, Dock.Right);
            buttons.Children.Add(right);
            buttons.Children.Add(_error);

            var page = new DockPanel { Background = Brushes.White };
            DockPanel.SetDock(_pageHeader, Dock.Top);
            page.Children.Add(_pageHeader);
            page.Children.Add(new ScrollViewer { Content = new Border { Padding = new Thickness(16, 14), Child = _pageHost } });

            var body = new Grid { ColumnDefinitions = new ColumnDefinitions("230,6,*") };
            var listBorder = new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Background = Brushes.White, Child = _list };
            var pageBorder = new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = page };
            Grid.SetColumn(pageBorder, 2);
            body.Children.Add(listBorder);
            body.Children.Add(pageBorder);

            var root = new DockPanel { Margin = new Thickness(12) };
            DockPanel.SetDock(buttons, Dock.Bottom);
            root.Children.Add(buttons);
            root.Children.Add(body);
            Content = root;
            _list.SelectedIndex = 0;
            foreach (var p in _pages)
                p.Changed += () => _error.Text = "";
        }

        private void ShowPage(EditPage p)
        {
            if (p == null) return;
            var head = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 10, Margin = new Thickness(10, 0), VerticalAlignment = VerticalAlignment.Center };
            head.Children.Add(Views.Icon(p.Icon));
            head.Children.Add(new TextBlock { Text = p.Title, Foreground = Brushes.White, FontWeight = FontWeight.Bold, FontSize = 14, VerticalAlignment = VerticalAlignment.Center });
            _pageHeader.Content = new Border { Height = 38, Background = Views.Gradient("#396D8C", "#3F8B89", true), Child = head };
            _pageHost.Content = p.View;
        }

        /// <summary>PropertiesDialog.okButton_Click: validate, then one MultipleAction with SaveChangesAction first.</summary>
        private async System.Threading.Tasks.Task Ok()
        {
            foreach (var p in _pages)
            {
                var error = p.ValidationError;
                if (error != null)
                {
                    _list.SelectedItem = p;
                    _error.Text = error;
                    return;
                }
            }
            var changed = _pages.Where(p => p.HasChanged).ToList();
            if (changed.Count == 0)
            {
                Close();
                return;
            }
            foreach (var p in changed)
                if (!await p.ConfirmSave(this))
                {
                    _list.SelectedItem = p;
                    return;
                }
            var actions = changed.Select(p => p.SaveSettings()).Where(a => a != null).ToList();
            actions.Insert(0, new SaveChangesAction(_copy, true, _orig));
            var name = Helpers.GetName(_orig).Ellipsise(50);
            Action = new MultipleAction(_orig.Connection, string.Format(Messages.UPDATE_PROPERTIES, name),
                Messages.UPDATING_PROPERTIES, string.Format(Messages.UPDATED_PROPERTIES, name), actions);
            Action.SetObject(_copy);
            Close();
        }
    }
}
