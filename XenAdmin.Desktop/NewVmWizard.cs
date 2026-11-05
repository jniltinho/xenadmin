using System;
using System.Collections.Generic;
using System.Linq;
using System.Xml;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using XenAdmin.Actions;
using XenAdmin.Actions.VMActions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// NewVMWizard: Template, Name, Installation Media, Home Server, CPU &amp; Memory, Storage, Networking, Finish.
    /// Builds the same CreateVMAction as the WinForms wizard.
    /// </summary>
    public class NewVmWizard : Window
    {
        private static readonly string[] Steps = { "Template", "Name", "Installation Media", "Home Server", "CPU & Memory", "Storage", "Networking", "Finish" };
        private static readonly string[] Titles =
        {
            "Select a VM template", "Name the new virtual machine", "Locate the operating system installation media", "Select a home server",
            "Allocate processor and memory resources", "Configure storage for the new VM", "Configure networking on the new VM", "Ready to create the new virtual machine"
        };

        private readonly IXenConnection _c;
        private readonly StackPanel _stepList = new StackPanel();
        private readonly TextBlock _title = new TextBlock { FontWeight = FontWeight.Bold, FontSize = 13, VerticalAlignment = VerticalAlignment.Center };
        private readonly ContentControl _page = new ContentControl();
        private readonly Button _prev = new Button { Content = "< Previous" };
        private readonly Button _next = new Button { Content = "Next >", IsDefault = true };
        private readonly Button _create = new Button { Content = "Create Now" };
        private readonly TextBlock _error = new TextBlock { Foreground = Brushes.Red, VerticalAlignment = VerticalAlignment.Center, TextWrapping = TextWrapping.Wrap };
        private int _step;

        // wizard state
        private VM _template;
        private string _templateRefForDefaults;
        private string _name = "", _description = "";
        private InstallMethod _install = InstallMethod.CD;
        private VDI _iso;
        private string _url = "";
        private VmBootMode _bootMode = VmBootMode.Bios;
        private Host _home;
        private long _vcpus = 1, _cores = 1;
        private long _memory;
        private readonly List<DiskDescription> _disks = new List<DiskDescription>();
        private readonly List<VIF> _vifs = new List<VIF>();
        private bool _start = true;

        public CreateVMAction Action { get; private set; }

        public NewVmWizard(IXenConnection connection)
        {
            _c = connection;
            Title = "New VM";
            Width = 810;
            Height = 570;
            MinWidth = 700;
            MinHeight = 500;
            WindowStartupLocation = WindowStartupLocation.CenterOwner;
            ShowInTaskbar = false;
            Background = Brush.Parse("#F0F0F0");

            var header = new DockPanel { Height = 60, Background = Brushes.White };
            var hicon = Views.Icon("000_CreateVM_h32bit_32", 32);
            hicon.Margin = new Thickness(12, 0, 16, 0);
            header.Children.Add(hicon);
            header.Children.Add(_title);

            var logo = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 10, Margin = new Thickness(12, 0, 0, 14), VerticalAlignment = VerticalAlignment.Bottom };
            try { logo.Children.Add(new Image { Source = new Avalonia.Media.Imaging.Bitmap(Avalonia.Platform.AssetLoader.Open(new Uri("avares://xcp-ng-center/Assets/xcp-ng-center.png"))), Width = 40, Height = 40 }); } catch { }
            logo.Children.Add(new TextBlock { Text = "XCP-ng", FontSize = 20, VerticalAlignment = VerticalAlignment.Center });
            var side = new DockPanel { Width = 190, Background = Brushes.White };
            DockPanel.SetDock(logo, Dock.Bottom);
            side.Children.Add(logo);
            side.Children.Add(_stepList);

            _prev.Click += (_, _) => Go(_step - 1);
            _next.Click += (_, _) => { if (Validate()) Go(_step + 1); };
            _create.Click += (_, _) => { if (Validate()) Finish(); };
            var cancel = new Button { Content = Messages.CANCEL, IsCancel = true };
            cancel.Click += (_, _) => Close();
            var buttons = new DockPanel { Margin = new Thickness(12, 10) };
            var right = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            right.Children.Add(_prev);
            right.Children.Add(_next);
            right.Children.Add(_create);
            right.Children.Add(cancel);
            DockPanel.SetDock(right, Dock.Right);
            buttons.Children.Add(right);
            buttons.Children.Add(_error);

            var body = new DockPanel();
            DockPanel.SetDock(side, Dock.Left);
            body.Children.Add(side);
            body.Children.Add(new Border { Padding = new Thickness(16, 14), Child = _page });

            var root = new DockPanel();
            DockPanel.SetDock(header, Dock.Top);
            DockPanel.SetDock(buttons, Dock.Bottom);
            root.Children.Add(header);
            root.Children.Add(new Border { Height = 1, Background = Brush.Parse("#C8C8C8") });
            DockPanel.SetDock(root.Children[1], Dock.Top);
            root.Children.Add(buttons);
            root.Children.Add(body);
            Content = root;

            _template = Templates().FirstOrDefault();
            Go(0);
        }

        // ---------- navigation ----------

        private void Go(int step)
        {
            _error.Text = "";
            _step = Math.Clamp(step, 0, Steps.Length - 1);
            if (_template != null && _template.opaque_ref != _templateRefForDefaults && _step > 0)
                ApplyTemplateDefaults();
            _title.Text = Titles[_step];
            _stepList.Children.Clear();
            for (int i = 0; i < Steps.Length; i++)
                _stepList.Children.Add(new Border
                {
                    Background = i == _step ? Brush.Parse("#0078D7") : Brushes.Transparent,
                    Padding = new Thickness(12, 4),
                    Child = new TextBlock { Text = Steps[i], FontWeight = i == _step ? FontWeight.Bold : FontWeight.Normal, Foreground = i == _step ? Brushes.White : Brushes.Black }
                });
            _prev.IsEnabled = _step > 0;
            _next.IsVisible = _step < Steps.Length - 1;
            _create.IsVisible = _step == Steps.Length - 1;
            _create.IsDefault = _create.IsVisible;
            _next.IsDefault = !_create.IsVisible;
            _page.Content = _step switch
            {
                0 => TemplatePage(), 1 => NamePage(), 2 => MediaPage(), 3 => HomePage(),
                4 => CpuMemPage(), 5 => StoragePage(), 6 => NetworkPage(), _ => FinishPage()
            };
        }

        private bool Validate()
        {
            string err = _step switch
            {
                0 when _template == null => "Select a template.",
                1 when string.IsNullOrWhiteSpace(_name) => "Enter a name for the new VM.",
                5 when _disks.Any(d => d.Disk.SR == null || d.Disk.SR.opaque_ref == Helper.NullOpaqueRef) => "Select a storage repository for every disk.",
                _ => null
            };
            _error.Text = err ?? "";
            return err == null;
        }

        private IEnumerable<VM> Templates() =>
            _c.Cache.VMs.Where(v => v.is_a_template && !v.is_a_snapshot && v.Show(false))
                .OrderBy(v => v.DefaultTemplate() ? 1 : 0).ThenBy(v => NewVmDefaults.Category(v)).ThenBy(v => v.Name(), StringComparer.CurrentCultureIgnoreCase);


        // ---------- defaults from the template (Page_Name, Page_CpuMem, Page_Storage, Page_Networking) ----------

        private void ApplyTemplateDefaults()
        {
            var t = _template;
            _templateRefForDefaults = t.opaque_ref;
            (_name, _description, _bootMode, _install, _vcpus, _cores, _memory) = NewVmDefaults.FromTemplate(t);
            _home = null;
            LoadDisks();
            LoadVifs();
        }

        public void LoadDisks()
        {
            _disks.Clear();
            _disks.AddRange(NewVmDefaults.DisksFor(_c, _template, _name, _home));
        }



        private void LoadVifs()
        {
            _vifs.Clear();
            _vifs.AddRange(NewVmDefaults.VifsFor(_c, _template));
        }


        // ---------- pages ----------

        private static TextBlock Text(string t) => new TextBlock { Text = t, TextWrapping = TextWrapping.Wrap, Margin = new Thickness(0, 0, 0, 12) };

        private Control TemplatePage()
        {
            var search = new TextBox { Watermark = "Search..." };
            var list = new ListBox { Height = 330, Classes = { "props" } };
            var details = new TextBlock { TextWrapping = TextWrapping.Wrap };
            void Fill()
            {
                var all = Templates().Where(t => string.IsNullOrWhiteSpace(search.Text) || t.Name().IndexOf(search.Text, StringComparison.CurrentCultureIgnoreCase) >= 0).ToList();
                list.ItemsSource = all.Select(t =>
                {
                    var row = new Grid { ColumnDefinitions = new ColumnDefinitions("20,*,90"), Tag = t };
                    row.Children.Add(Views.Icon(IconNames.For(t)));
                    var n = new TextBlock { Text = t.Name() }; Grid.SetColumn(n, 1); row.Children.Add(n);
                    var c = new TextBlock { Text = NewVmDefaults.Category(t) }; Grid.SetColumn(c, 2); row.Children.Add(c);
                    return row;
                }).ToList();
                list.SelectedItem = list.Items.OfType<Grid>().FirstOrDefault(g => ((VM)g.Tag).opaque_ref == _template?.opaque_ref) ?? list.Items.OfType<Grid>().FirstOrDefault();
            }
            list.SelectionChanged += (_, _) =>
            {
                if (list.SelectedItem is Grid g)
                {
                    _template = (VM)g.Tag;
                    details.Text = _template.DescriptionType() == VM.VmDescriptionType.None ? _template.Name() : _template.Description();
                }
            };
            search.TextChanged += (_, _) => Fill();
            Fill();
            var left = new DockPanel();
            DockPanel.SetDock(search, Dock.Top);
            left.Children.Add(search);
            var header = new Grid { ColumnDefinitions = new ColumnDefinitions("20,*,90"), Margin = new Thickness(6, 6, 6, 2) };
            var h1 = new TextBlock { Text = "Name" }; Grid.SetColumn(h1, 1); header.Children.Add(h1);
            var h2 = new TextBlock { Text = "Category" }; Grid.SetColumn(h2, 2); header.Children.Add(h2);
            DockPanel.SetDock(header, Dock.Top);
            left.Children.Add(header);
            left.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = list });
            var detailBox = new Border { BorderBrush = Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1), Padding = new Thickness(8), Width = 200, Margin = new Thickness(8, 0, 0, 0) };
            var dp = new StackPanel();
            dp.Children.Add(new TextBlock { Text = "Template details", Margin = new Thickness(0, 0, 0, 8) });
            dp.Children.Add(details);
            detailBox.Child = dp;
            var grid = new DockPanel();
            DockPanel.SetDock(detailBox, Dock.Right);
            grid.Children.Add(detailBox);
            grid.Children.Add(left);
            return grid;
        }

        private Control NamePage()
        {
            var name = new TextBox { Text = _name };
            var desc = new TextBox { Text = _description, AcceptsReturn = true, TextWrapping = TextWrapping.Wrap, Height = 80 };
            name.TextChanged += (_, _) => _name = name.Text ?? "";
            desc.TextChanged += (_, _) => _description = desc.Text ?? "";
            var p = new StackPanel { Spacing = 8 };
            p.Children.Add(Text("Enter a name that will help you to identify the virtual machine later. This could be a name that describes its software and hardware, such as Windows Server 2022 Web Server."));
            p.Children.Add(new Label { Content = Dialogs.M("&Name:"), Target = name, Padding = new Thickness(0) });
            p.Children.Add(name);
            p.Children.Add(new Label { Content = Dialogs.M("&Description:"), Target = desc, Padding = new Thickness(0) });
            p.Children.Add(desc);
            Opened += (_, _) => { };
            return p;
        }

        private Control MediaPage()
        {
            var p = new StackPanel { Spacing = 8 };
            if (_template.IsHVM())
            {
                p.Children.Add(Text("Select the installation method for the operating system software you want to install on the new VM."));
                var isoRadio = new RadioButton { Content = "Install from ISO library or DVD drive:", GroupName = "m", IsChecked = _install == InstallMethod.CD };
                var home = _home;
                var isos = new List<VDI> { null };
                isos.AddRange(_c.Cache.SRs.Where(s => s.content_type == SR.Content_Type_ISO && s.Show(false) && (home == null || s.CanBeSeenFrom(home)))
                    .SelectMany(s => _c.ResolveAll(s.VDIs).Where(v => v.Show(false))).OrderBy(v => v.Name()));
                var combo = new ComboBox { ItemsSource = isos.Select(v => v == null ? "<empty>" : $"{_c.Resolve(v.SR)?.Name()}: {v.Name()}").ToList(), Width = 420, Margin = new Thickness(22, 0, 0, 0) };
                combo.SelectedIndex = Math.Max(0, isos.FindIndex(v => v?.opaque_ref == _iso?.opaque_ref));
                combo.SelectionChanged += (_, _) => { _iso = combo.SelectedIndex > 0 ? isos[combo.SelectedIndex] : null; isoRadio.IsChecked = true; };
                var netRadio = new RadioButton { Content = "Boot from network", GroupName = "m", IsChecked = _install == InstallMethod.Network };
                isoRadio.IsCheckedChanged += (_, _) => { if (isoRadio.IsChecked == true) _install = InstallMethod.CD; };
                netRadio.IsCheckedChanged += (_, _) => { if (netRadio.IsChecked == true) _install = InstallMethod.Network; };
                p.Children.Add(isoRadio);
                p.Children.Add(combo);
                p.Children.Add(netRadio);
                if (_template.SupportsUefiBoot())
                {
                    p.Children.Add(new TextBlock { Text = "Boot Mode", Margin = new Thickness(0, 16, 0, 0) });
                    var bios = new RadioButton { Content = "BIOS Boot", GroupName = "b", IsChecked = _bootMode == VmBootMode.Bios };
                    var uefi = new RadioButton { Content = "UEFI Boot", GroupName = "b", IsChecked = _bootMode == VmBootMode.Uefi };
                    var secure = new RadioButton { Content = "UEFI Secure Boot", GroupName = "b", IsChecked = _bootMode == VmBootMode.SecureUefi, IsEnabled = _template.SupportsSecureUefiBoot() };
                    bios.IsCheckedChanged += (_, _) => { if (bios.IsChecked == true) _bootMode = VmBootMode.Bios; };
                    uefi.IsCheckedChanged += (_, _) => { if (uefi.IsChecked == true) _bootMode = VmBootMode.Uefi; };
                    secure.IsCheckedChanged += (_, _) => { if (secure.IsChecked == true) _bootMode = VmBootMode.SecureUefi; };
                    p.Children.Add(bios); p.Children.Add(uefi); p.Children.Add(secure);
                }
            }
            else
            {
                p.Children.Add(Text("Enter the URL of the network installation repository for the operating system."));
                var url = new TextBox { Text = _url, Watermark = "http://, https://, ftp:// or nfs://" };
                url.TextChanged += (_, _) => _url = url.Text ?? "";
                p.Children.Add(new Label { Content = Dialogs.M("&URL:"), Target = url, Padding = new Thickness(0) });
                p.Children.Add(url);
            }
            return p;
        }

        private Control HomePage()
        {
            var p = new StackPanel { Spacing = 8 };
            p.Children.Add(Text("Nominating a home server for a virtual machine ensures that the virtual machine will be started on that server if it is available."));
            var none = new RadioButton { Content = "Don't assign this VM a home server. The VM will be started on any server with the necessary resources.", GroupName = "h", IsChecked = _home == null };
            var place = new RadioButton { Content = "Place the VM on this server:", GroupName = "h", IsChecked = _home != null };
            var hosts = _c.Cache.Hosts.OrderBy(h => h.Name()).ToList();
            var list = new ListBox { ItemsSource = hosts.Select(h => h.Name()).ToList(), Height = 200, Margin = new Thickness(22, 0, 0, 0), SelectedIndex = _home == null ? -1 : hosts.FindIndex(h => h.opaque_ref == _home.opaque_ref) };
            list.SelectionChanged += (_, _) => { if (list.SelectedIndex >= 0) { _home = hosts[list.SelectedIndex]; place.IsChecked = true; LoadDisks(); } };
            none.IsCheckedChanged += (_, _) => { if (none.IsChecked == true) { _home = null; list.SelectedIndex = -1; LoadDisks(); } };
            p.Children.Add(none);
            p.Children.Add(place);
            p.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = list, Margin = new Thickness(22, 0, 0, 0) });
            return p;
        }

        private Control CpuMemPage()
        {
            var max = Math.Max(1, _template.MaxVCPUsAllowed());
            var vcpus = new ComboBox { ItemsSource = Enumerable.Range(1, max).Select(i => (long)i).ToList(), SelectedItem = Math.Min(_vcpus, max), Width = 70 };
            var topology = new ComboBox { Width = 260 };
            void FillTopology()
            {
                var opts = VmRules.ValidTopologies(_vcpus, _template.MaxCoresPerSocket());
                topology.ItemsSource = opts.Select(o => VM.GetTopology(o.Sockets, o.Cores)).ToList();
                var i = opts.FindIndex(o => o.Cores == _cores);
                topology.SelectedIndex = i >= 0 ? i : 0;
                _cores = opts.Count > 0 ? opts[Math.Max(0, i)].Cores : 1;
            }
            FillTopology();
            vcpus.SelectionChanged += (_, _) => { if (vcpus.SelectedItem is long v) { _vcpus = v; FillTopology(); } };
            topology.SelectionChanged += (_, _) =>
            {
                var opts = VmRules.ValidTopologies(_vcpus, _template.MaxCoresPerSocket());
                if (topology.SelectedIndex >= 0 && topology.SelectedIndex < opts.Count) _cores = opts[topology.SelectedIndex].Cores;
            };
            var memory = new NumericUpDown { Minimum = Math.Max(0.125m, (decimal)_template.memory_static_min / Util.BINARY_GIGA), Maximum = 1536, Increment = 0.5m, FormatString = "0.###", Width = 120, Value = (decimal)_memory / Util.BINARY_GIGA };
            memory.ValueChanged += (_, _) => _memory = (long)((memory.Value ?? 1) * Util.BINARY_GIGA);
            var g = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,Auto,Auto"), RowDefinitions = new RowDefinitions("Auto,Auto,Auto"), RowSpacing = 10, ColumnSpacing = 10 };
            void Row(int r, string label, Control c, string suffix = null)
            {
                var l = new TextBlock { Text = label, VerticalAlignment = VerticalAlignment.Center };
                Grid.SetRow(l, r); Grid.SetRow(c, r); Grid.SetColumn(c, 1);
                g.Children.Add(l); g.Children.Add(c);
                if (suffix != null) { var s = new TextBlock { Text = suffix, VerticalAlignment = VerticalAlignment.Center }; Grid.SetRow(s, r); Grid.SetColumn(s, 2); g.Children.Add(s); }
            }
            Row(0, "Number of vCPUs:", vcpus);
            Row(1, "Topology:", topology);
            Row(2, "Memory:", memory, "GB");
            var p = new StackPanel();
            p.Children.Add(Text("Specify the number of virtual CPUs, their topology, and the amount of memory that will be allocated to the new virtual machine."));
            p.Children.Add(g);
            return p;
        }

        private Control StoragePage()
        {
            var srs = StorageRules.NewDiskTargets(_c, _home);
            var rows = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*,Auto,Auto"), RowSpacing = 6, ColumnSpacing = 10 };
            void H(int col, string t) { var x = new TextBlock { Text = t, FontWeight = FontWeight.SemiBold }; Grid.SetColumn(x, col); rows.Children.Add(x); }
            rows.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            H(0, "Position"); H(1, "Name"); H(2, "Size"); H(3, "Location");
            int r = 1;
            DiskDescription selected = null;
            var remove = new Button { Content = Dialogs.M("&Delete"), IsEnabled = false, MinWidth = 75 };
            foreach (var d in _disks)
            {
                rows.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var pos = new RadioButton { Content = d.Device.userdevice, GroupName = "disk", VerticalAlignment = VerticalAlignment.Center };
                pos.IsCheckedChanged += (_, _) => { if (pos.IsChecked == true) { selected = d; remove.IsEnabled = true; } };
                var name = new TextBox { Text = d.Disk.name_label, MinWidth = 160 };
                name.TextChanged += (_, _) => d.Disk.name_label = name.Text ?? "";
                var size = new NumericUpDown { Value = (decimal)d.Disk.virtual_size / Util.BINARY_GIGA, Minimum = 0.0625m, Maximum = 2040, Increment = 1, FormatString = "0.###", Width = 110 };
                size.ValueChanged += (_, _) => d.Disk.virtual_size = (long)((size.Value ?? 1) * Util.BINARY_GIGA);
                var sr = new ComboBox { ItemsSource = srs.Select(s => $"{s.NameWithLocation()} ({Util.DiskSizeString(s.FreeSpace())} free)").ToList(), Width = 300,
                    SelectedIndex = srs.FindIndex(s => s.opaque_ref == d.Disk.SR?.opaque_ref) };
                sr.SelectionChanged += (_, _) => { if (sr.SelectedIndex >= 0) d.Disk.SR = new XenRef<SR>(srs[sr.SelectedIndex]); };
                Grid.SetRow(pos, r); Grid.SetRow(name, r); Grid.SetRow(size, r); Grid.SetRow(sr, r);
                Grid.SetColumn(name, 1); Grid.SetColumn(size, 2); Grid.SetColumn(sr, 3);
                rows.Children.Add(pos); rows.Children.Add(name); rows.Children.Add(size); rows.Children.Add(sr);
                r++;
            }
            // Page_Storage: Add a new disk, Delete the selected one (the radio in the Position column)
            var add = new Button { Content = Dialogs.M("&Add..."), MinWidth = 75, IsEnabled = srs.Count > 0 };
            add.Click += (_, _) => { _disks.Add(NewVmDefaults.NewDisk(_c, _home, _disks)); _page.Content = StoragePage(); };
            remove.Click += (_, _) => { if (selected != null) { _disks.Remove(selected); _page.Content = StoragePage(); } };
            var buttons = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(0, 10, 0, 0) };
            buttons.Children.Add(add);
            buttons.Children.Add(remove);
            var p = new StackPanel();
            p.Children.Add(Text(_disks.Count == 0
                ? "This template has no virtual disks. Add one, or the new VM will be created without disks and can boot from the network or a DVD."
                : "The virtual machine template you selected provides the virtual disks listed below. You can change their size and location, add disks or delete them."));
            p.Children.Add(rows);
            p.Children.Add(buttons);
            return p;
        }

        private Control NetworkPage()
        {
            var list = new ListBox { Height = 200 };
            void Fill() => list.ItemsSource = _vifs.Select(v => $"{v.device}   {(string.IsNullOrEmpty(v.MAC) ? "<auto-generated>" : v.MAC)}   {_c.Resolve(v.network)?.Name()}").ToList();
            Fill();
            var add = new Button { Content = Dialogs.M("&Add...") };
            var nets = _c.Cache.Networks.Where(n => n.Show(false) && !n.IsGuestInstallerNetwork() && !n.IsMember()).OrderBy(n => n.Name()).ToList();
            var netCombo = new ComboBox { ItemsSource = nets.Select(n => n.Name()).ToList(), SelectedIndex = nets.Count > 0 ? 0 : -1, Width = 260 };
            add.IsEnabled = nets.Count > 0;
            add.Click += (_, _) =>
            {
                if (netCombo.SelectedIndex < 0 || _vifs.Count >= _template.MaxVIFsAllowed()) return;
                var device = Enumerable.Range(0, 64).First(i => _vifs.All(v => v.device != i.ToString())).ToString();
                _vifs.Add(new VIF { Connection = _c, device = device, network = new XenRef<XenAPI.Network>(nets[netCombo.SelectedIndex]) });
                Fill();
            };
            var del = new Button { Content = Dialogs.M("&Delete") };
            del.Click += (_, _) => { if (list.SelectedIndex >= 0) { _vifs.RemoveAt(list.SelectedIndex); Fill(); } };
            var bar = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(0, 8, 0, 0) };
            bar.Children.Add(netCombo); bar.Children.Add(add); bar.Children.Add(del);
            var p = new StackPanel();
            p.Children.Add(Text("The new VM will be connected to the networks listed below. You can add or remove virtual network interfaces."));
            p.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = list });
            p.Children.Add(bar);
            return p;
        }

        private Control FinishPage()
        {
            var start = new CheckBox { Content = "Start the new VM automatically", IsChecked = _start };
            start.IsCheckedChanged += (_, _) => _start = start.IsChecked == true;
            var summary = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowSpacing = 4, ColumnSpacing = 16 };
            var items = new List<(string, string)>
            {
                ("Template", _template.Name()), ("Name", _name),
                ("Install method", _install == InstallMethod.CD ? (_iso == null ? "DVD drive (empty)" : _iso.Name()) : _template.IsHVM() ? "Network boot" : _url),
                ("Home server", _home?.Name() ?? "None"), ("vCPUs", _vcpus.ToString()), ("Topology", VM.GetTopology(_vcpus / Math.Max(1, _cores), _cores)),
                ("Memory", Util.MemorySizeStringSuitableUnits(_memory, true)),
            };
            items.AddRange(_disks.Select(d => ("Disk", $"{d.Disk.name_label} ({Util.DiskSizeString(d.Disk.virtual_size)}) on {_c.Resolve(d.Disk.SR)?.Name()}")));
            items.AddRange(_vifs.Select(v => ("Network interface", _c.Resolve(v.network)?.Name())));
            for (int i = 0; i < items.Count; i++)
            {
                summary.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var k = new TextBlock { Text = items[i].Item1 + ":", Opacity = 0.7 };
                var v = new TextBlock { Text = items[i].Item2 ?? "", TextWrapping = TextWrapping.Wrap };
                Grid.SetRow(k, i); Grid.SetRow(v, i); Grid.SetColumn(v, 1);
                summary.Children.Add(k); summary.Children.Add(v);
            }
            var p = new StackPanel { Spacing = 10 };
            p.Children.Add(Text("All the necessary information has been collected and the wizard is ready to provision the new virtual machine using the settings shown below."));
            p.Children.Add(new ScrollViewer { Content = summary, MaxHeight = 280 });
            p.Children.Add(start);
            return p;
        }

        private void Finish()
        {
            if (_templateRefForDefaults != _template.opaque_ref)
                ApplyTemplateDefaults();
            Action = NewVmDefaults.Build(_c, _template, _name, _description, _install, _iso, _url, _bootMode, _home, _vcpus, _cores, _memory, _disks, _vifs, _start);
            Close();
        }

    }
}
