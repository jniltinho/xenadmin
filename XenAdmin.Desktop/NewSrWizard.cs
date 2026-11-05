using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using XenAdmin.Actions;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// NewSRWizard: Type, Name, Location (VHDoNFS / CifsFrontend / LVMoISCSI pages). Creates, reattaches or introduces the SR
    /// with the same actions as the WinForms wizard; the other storage types are listed disabled until ported.
    /// </summary>
    public class NewSrWizard : Window
    {
        private static readonly string[] Steps = { "Type", "Name", "Location" };

        private readonly IXenConnection _c;
        private readonly StackPanel _stepList = new StackPanel();
        private readonly TextBlock _title = new TextBlock { FontWeight = FontWeight.Bold, FontSize = 13, VerticalAlignment = VerticalAlignment.Center };
        private readonly ContentControl _page = new ContentControl();
        private readonly Button _prev = new Button { Content = "< Previous" };
        private readonly Button _next = new Button { Content = "Next >", IsDefault = true };
        private readonly Button _finish = new Button { Content = "Finish" };
        private readonly TextBlock _error = new TextBlock { Foreground = Brushes.Red, VerticalAlignment = VerticalAlignment.Center, TextWrapping = TextWrapping.Wrap };
        private int _step;

        // wizard state
        private SrRules.Kind _kind = SrRules.Kind.Nfs;
        private string _name;
        private SrRules.Kind? _nameFor;
        private string _description = "";
        private bool _autoDescription = true;
        private string _path = "", _options = "", _nfsVersion = "", _user = "", _password = "";
        private List<SR.SRInfo> _found = new List<SR.SRInfo>();
        private SR.SRInfo _existing;
        private bool _reattach;      // "Reattach an existing SR" chosen (needs a selection)
        private int _scanGeneration; // a Scan's result is dropped when the location changed meanwhile
        // iSCSI: target, CHAP, the scanned IQNs and LUNs and the picks
        private string _iscsiHost = "", _iscsiPort = SrRules.DefaultIscsiPort.ToString(), _chapUser = "", _chapPassword = "";
        private bool _chap;
        private List<IScsiIqnInfo> _iqns = new List<IScsiIqnInfo>();
        private List<ISCSIInfo> _luns = new List<ISCSIInfo>();
        private IScsiIqnInfo? _iqn;
        private ISCSIInfo? _lun;
        private bool _otherUser; // SMB ISO: "Use different user name"
        private string ChapUser => _chap ? _chapUser : null;
        private string ChapPassword => _chap ? _chapPassword : null;
        private bool _busy;
        private readonly Action<AsyncAction> _run;

        /// <summary>The wizard runs its action through run (the main window's task list) and closes only when it succeeds.</summary>
        public NewSrWizard(IXenConnection connection, Action<AsyncAction> run)
        {
            _c = connection;
            _run = run;
            Title = "New Storage Repository - " + XenAdmin.Core.Helpers.GetName(connection);
            Width = 810;
            Height = 545;
            MinWidth = 680;
            MinHeight = 480;
            WindowStartupLocation = WindowStartupLocation.CenterOwner;
            ShowInTaskbar = false;
            Background = Brush.Parse("#F0F0F0");

            var header = new DockPanel { Height = 60, Background = Brushes.White };
            var hicon = Views.Icon("000_NewStorage_h32bit_32", 32);
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
            _finish.Click += async (_, _) => { if (Validate()) await Finish(); };
            var cancel = new Button { Content = Messages.CANCEL, IsCancel = true };
            cancel.Click += (_, _) => Close();
            var buttons = new DockPanel { Margin = new Thickness(12, 10) };
            var right = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            right.Children.Add(_prev);
            right.Children.Add(_next);
            right.Children.Add(_finish);
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
            Go(0);
        }

        private void Go(int step)
        {
            _error.Text = "";
            _step = Math.Clamp(step, 0, Steps.Length - 1);
            if (_step > 0 && _nameFor != _kind)
            {
                _name = SrRules.DefaultName(_kind, _c);
                _nameFor = _kind;
            }
            _title.Text = _step switch
            {
                0 => "Choose the type of new storage",
                1 => "What do you want to call this Storage Repository?",
                _ => _kind switch
                {
                    SrRules.Kind.Nfs or SrRules.Kind.NfsIso => Messages.NEWSR_PATH_NFS, SrRules.Kind.Iscsi => Messages.NEWSR_PATH_ISCSI, _ => Messages.NEWSR_PATH_CIFS
                }
            };
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
            _finish.IsVisible = _step == Steps.Length - 1;
            _finish.IsDefault = _finish.IsVisible;
            _next.IsDefault = !_finish.IsVisible;
            _page.Content = _step switch { 0 => TypePage(), 1 => NamePage(), _ => _kind == SrRules.Kind.Iscsi ? IscsiPage() : LocationPage() };
            UpdateFinish();
        }

        private bool Validate()
        {
            string err = _step switch
            {
                1 when string.IsNullOrWhiteSpace(_name) => "Enter a name for the new SR.",
                2 when _kind == SrRules.Kind.Iscsi && (_iqn == null || _lun == null) => "Scan the target host and select the target IQN and LUN.",
                2 when _kind == SrRules.Kind.Iscsi && SrRules.LunInUse(_lun.Value) is SR used =>
                    string.Format(Messages.NEWSR_LUN_IN_USE_ON_POOL, used.Name(), XenAdmin.Core.Helpers.GetName(used.Connection)),
                2 when _kind == SrRules.Kind.Iscsi => null,
                2 when !SrRules.ValidLocation(_kind, _path) => SrRules.IsNfsPath(_kind) ? "Enter the share as server:/path." : @"Enter the share as \\server\share.",
                2 when SrRules.IsIso(_kind) => SrRules.IsoAlreadyAttached(_c, _kind, _path),
                2 when _reattach && _existing == null => "Select the SR to reattach.",
                2 when _existing != null && SrRules.InUse(_existing.UUID) != null => "This SR is already in use by " + XenAdmin.Core.Helpers.GetName(SrRules.InUse(_existing.UUID).Connection) + ".",
                _ => null
            };
            _error.Text = err ?? "";
            return err == null;
        }

        private void UpdateFinish()
        {
            _finish.IsEnabled = !_busy && (_kind == SrRules.Kind.Iscsi
                ? _iqn != null && _lun != null
                : SrRules.ValidLocation(_kind, _path) && (!_reattach || _existing != null));
            _next.IsEnabled = _prev.IsEnabled = !_busy;
            if (_step == 0) _prev.IsEnabled = false;
        }

        private static TextBlock Text(string s) => new TextBlock { Text = s, TextWrapping = TextWrapping.Wrap };

        // ---------- pages ----------

        private Control TypePage()
        {
            var typeName = new TextBlock { Text = SrRules.TypeName(_kind), FontStyle = FontStyle.Italic, Margin = new Thickness(0, 0, 0, 12) };
            var blurb = Text(SrRules.Blurb(_kind));
            RadioButton Option(string text, SrRules.Kind? kind)
            {
                var r = new RadioButton { Content = text, GroupName = "srtype", IsEnabled = kind != null, IsChecked = kind == _kind, Margin = new Thickness(24, 2, 0, 2) };
                if (kind != null)
                    r.IsCheckedChanged += (_, _) =>
                    {
                        if (r.IsChecked != true) return;
                        _kind = kind.Value;
                        typeName.Text = SrRules.TypeName(_kind);
                        blurb.Text = SrRules.Blurb(_kind);
                        _existing = null;
                        _found.Clear();
                    };
                return r;
            }
            TextBlock Group(string text, double left) => new TextBlock { Text = text, FontWeight = FontWeight.SemiBold, Margin = new Thickness(left, 10, 0, 4) };
            var list = new StackPanel();
            list.Children.Add(Group("Virtual disk storage", 0));
            list.Children.Add(Group("Block based storage", 12));
            list.Children.Add(Option(Messages.NEWSR_LVMOISCSI_TYPE_NAME, SrRules.Kind.Iscsi));
            list.Children.Add(Option(Messages.NEWSR_LVMOHBA_TYPE_NAME, null));
            list.Children.Add(Option(Messages.NEWSR_LVMOFCOE_TYPE_NAME, null));
            list.Children.Add(Group("File based storage", 12));
            list.Children.Add(Option(SrRules.TypeName(SrRules.Kind.Nfs), SrRules.Kind.Nfs));
            list.Children.Add(Option(SrRules.TypeName(SrRules.Kind.Smb), SrRules.Kind.Smb));
            list.Children.Add(Option("Local Storage", null));
            list.Children.Add(Group("ISO library", 0));
            list.Children.Add(Option(Messages.NEWSR_CIFS_ISO_TYPE_NAME, SrRules.Kind.SmbIso));
            list.Children.Add(Option(Messages.NEWSR_NFS_ISO_TYPE_NAME, SrRules.Kind.NfsIso));
            var right = new StackPanel { Margin = new Thickness(0, 10, 0, 0) };
            right.Children.Add(typeName);
            right.Children.Add(blurb);
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("236,*"), ColumnSpacing = 16 };
            grid.Children.Add(list);
            Grid.SetColumn(right, 1);
            grid.Children.Add(right);
            return grid;
        }

        private Control NamePage()
        {
            var name = new TextBox { Text = _name };
            var auto = new CheckBox { Content = Dialogs.M("&Autogenerate description based on SR settings (e.g., IP address, LUN etc.)"), IsChecked = _autoDescription, Margin = new Thickness(0, 16, 0, 0) };
            var desc = new TextBox { Text = _description, AcceptsReturn = true, TextWrapping = TextWrapping.Wrap, Height = 80, IsEnabled = !_autoDescription };
            name.TextChanged += (_, _) => _name = name.Text ?? "";
            desc.TextChanged += (_, _) => _description = desc.Text ?? "";
            var descLabel = new Label { Content = Dialogs.M("&Description:"), Target = desc, Padding = new Thickness(0, 4, 0, 0), IsEnabled = !_autoDescription, VerticalAlignment = VerticalAlignment.Top };
            auto.IsCheckedChanged += (_, _) => { _autoDescription = auto.IsChecked == true; desc.IsEnabled = descLabel.IsEnabled = !_autoDescription; };
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("78,*"), RowDefinitions = new RowDefinitions("Auto,Auto,Auto"), RowSpacing = 8 };
            var l = new Label { Content = Dialogs.M("&Name:"), Target = name, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
            Grid.SetColumn(name, 1);
            Grid.SetRow(auto, 1); Grid.SetColumnSpan(auto, 2);
            Grid.SetRow(descLabel, 2); Grid.SetRow(desc, 2); Grid.SetColumn(desc, 1);
            grid.Children.Add(l); grid.Children.Add(name); grid.Children.Add(auto); grid.Children.Add(descLabel); grid.Children.Add(desc);
            var p = new StackPanel { Spacing = 14 };
            p.Children.Add(Text("Provide a name and a description (optional) for your SR."));
            p.Children.Add(grid);
            Opened += (_, _) => { };
            name.AttachedToVisualTree += (_, _) => { name.Focus(); name.SelectAll(); };
            return p;
        }

        private Control LocationPage()
        {
            var nfs = SrRules.IsNfsPath(_kind);
            var iso = SrRules.IsIso(_kind); // ISO libraries are created only: no Scan, no reattach
            var path = new TextBox { Text = _path };
            var scan = new Button { Content = Dialogs.M("&Scan"), IsEnabled = SrRules.ValidLocation(_kind, _path) };
            var status = new TextBlock { Opacity = 0.7, VerticalAlignment = VerticalAlignment.Center };
            var createNew = new RadioButton { Content = Dialogs.M("&Create a new SR"), GroupName = "srmode", IsChecked = !_reattach };
            var reattach = new RadioButton { Content = Dialogs.M("&Reattach an existing SR:"), GroupName = "srmode", IsChecked = _reattach, IsEnabled = _found.Count > 0 };
            string Item(SR.SRInfo f) => SrRules.InUse(f.UUID) is SR used ? $"{f} (in use by {XenAdmin.Core.Helpers.GetName(used.Connection)})" : f.ToString();
            var existing = new ListBox { Height = 90, IsEnabled = _found.Count > 0, ItemsSource = _found.Select(Item).ToList() };
            existing.SelectedIndex = _existing == null ? -1 : _found.IndexOf(_existing);

            void Reset()
            {
                _scanGeneration++;
                _found.Clear();
                _existing = null;
                _reattach = false;
                existing.ItemsSource = null;
                existing.IsEnabled = reattach.IsEnabled = false;
                createNew.IsChecked = true;
                scan.IsEnabled = SrRules.ValidLocation(_kind, _path);
                UpdateFinish();
            }
            path.TextChanged += (_, _) => { _path = path.Text ?? ""; Reset(); };
            existing.SelectionChanged += (_, _) =>
            {
                if (existing.SelectedIndex < 0 || existing.SelectedIndex >= _found.Count) return;
                _existing = _found[existing.SelectedIndex];
                reattach.IsChecked = true;
                UpdateFinish();
            };
            createNew.IsCheckedChanged += (_, _) => { if (createNew.IsChecked == true) { _existing = null; existing.SelectedIndex = -1; } };
            reattach.IsCheckedChanged += (_, _) => { _reattach = reattach.IsChecked == true; UpdateFinish(); };

            var fields = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*,Auto"), RowSpacing = 8, ColumnSpacing = 8 };
            void Row(int r, string label, Control field, Control extra = null)
            {
                fields.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var top = field is StackPanel;
                var l = new Label { Content = Dialogs.M(label), Target = field, VerticalAlignment = top ? VerticalAlignment.Top : VerticalAlignment.Center, VerticalContentAlignment = top ? VerticalAlignment.Top : VerticalAlignment.Center, Padding = new Thickness(0, top ? 2 : 0, 0, 0) };
                Grid.SetRow(l, r); Grid.SetRow(field, r); Grid.SetColumn(field, 1);
                fields.Children.Add(l); fields.Children.Add(field);
                if (extra != null) { Grid.SetRow(extra, r); Grid.SetColumn(extra, 2); fields.Children.Add(extra); }
            }
            scan.IsVisible = createNew.IsVisible = reattach.IsVisible = !iso;
            Row(0, "&Share Name:", path, scan);
            var example = new TextBlock { Text = nfs ? "Example: server:/path" : @"Example: \\server\share", Margin = new Thickness(0, -4, 0, 0) };
            fields.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            Grid.SetRow(example, 1); Grid.SetColumn(example, 1);
            fields.Children.Add(example);
            var p = new StackPanel { Spacing = 10 };
            p.Children.Add(Text(_kind == SrRules.Kind.NfsIso ? "Provide the name of the share where your NFS ISO storage is located."
                : _kind == SrRules.Kind.SmbIso ? "Provide the name of the share where your SR is located. You can optionally specify alternative credentials by setting the server options."
                : nfs
                ? "Provide the name of the share where your SR is located, optionally specifying advanced options. Indicate whether you want to create a new SR or reattach an existing SR before proceeding."
                : "Provide the name of the share where your SR is located, optionally specifying the credentials for it. Indicate whether you want to create a new SR or reattach an existing SR before proceeding."));
            p.Children.Add(fields);
            if (nfs)
            {
                var options = new TextBox { Text = _options };
                options.TextChanged += (_, _) => { _options = options.Text ?? ""; Reset(); };
                var versions = new StackPanel { Spacing = 6 };
                foreach (var (text, value) in new[] { ("NFSv&3", ""), ("NFSv&4", "4"), ("NFSv4.&1", "4.1") })
                {
                    if (iso && value == "4.1") continue; // NFS_ISO offers 3 and 4
                    var r = new RadioButton { Content = Dialogs.M(text), GroupName = "nfsver", IsChecked = _nfsVersion == value };
                    r.IsCheckedChanged += (_, _) => { if (r.IsChecked == true) _nfsVersion = value; };
                    versions.Children.Add(r);
                }
                Row(2, "NFS Version:", versions);
                if (!iso)
                    Row(3, "Advanced &Options:", options);
            }
            else if (iso)
            {
                var other = new CheckBox { Content = Dialogs.M("Use &different user name"), IsChecked = _otherUser };
                var user = new TextBox { Text = _user, IsEnabled = _otherUser };
                var pass = new TextBox { Text = _password, PasswordChar = '•', IsEnabled = _otherUser };
                other.IsCheckedChanged += (_, _) => { _otherUser = other.IsChecked == true; user.IsEnabled = pass.IsEnabled = _otherUser; };
                user.TextChanged += (_, _) => _user = user.Text ?? "";
                pass.TextChanged += (_, _) => _password = pass.Text ?? "";
                fields.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                Grid.SetRow(other, 2); Grid.SetColumn(other, 1);
                fields.Children.Add(other);
                Row(3, "&Username:", user);
                Row(4, "Pass&word:", pass);
            }
            else
            {
                var user = new TextBox { Text = _user };
                var pass = new TextBox { Text = _password, PasswordChar = '•' };
                user.TextChanged += (_, _) => { _user = user.Text ?? ""; Reset(); };
                pass.TextChanged += (_, _) => { _password = pass.Text ?? ""; Reset(); };
                Row(2, "&Username:", user);
                Row(3, "&Password:", pass);
            }
            p.Children.Add(status);
            if (!iso)
            {
                p.Children.Add(createNew);
                p.Children.Add(reattach);
                p.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Margin = new Thickness(22, 0, 0, 0), Child = existing });
            }

            scan.Click += async (_, _) =>
            {
                scan.IsEnabled = false;
                status.Text = string.Format(Messages.ACTION_SCANNING_SR_FROM, XenAdmin.Core.Helpers.GetName(_c));
                var generation = ++_scanGeneration;
                var probe = SrRules.Probe(_c, _kind, SrRules.DeviceConfig(_kind, _path, _options, _nfsVersion, _user, _password, forProbe: true));
                await Run(probe);
                if (generation != _scanGeneration)
                    return; // the location changed while scanning: these SRs belong to the old one
                Reset();
                if (!probe.Succeeded)
                {
                    status.Text = probe.Exception?.Message ?? "Scan failed.";
                    return;
                }
                _found = probe.SRs ?? new List<SR.SRInfo>();
                status.Text = _found.Count == 0 ? Messages.NEWSR_NO_SRS_FOUND : "";
                existing.ItemsSource = _found.Select(Item).ToList();
                existing.IsEnabled = reattach.IsEnabled = _found.Count > 0;
            };
            return p;
        }

        /// <summary>LVMoISCSI: target host and port, optional CHAP, Scan Target Host, then the target IQN and LUN.</summary>
        private Control IscsiPage()
        {
            var host = new TextBox { Text = _iscsiHost };
            var port = new TextBox { Text = _iscsiPort, Width = 60 };
            var chap = new CheckBox { Content = Dialogs.M("Use &CHAP authentication"), IsChecked = _chap };
            var user = new TextBox { Text = _chapUser, IsEnabled = _chap };
            var pass = new TextBox { Text = _chapPassword, PasswordChar = '•', IsEnabled = _chap };
            var scan = new Button { Content = Dialogs.M("&Scan Target Host") };
            var iqns = new ComboBox { HorizontalAlignment = HorizontalAlignment.Stretch };
            var luns = new ComboBox { HorizontalAlignment = HorizontalAlignment.Stretch };
            var status = new TextBlock { Opacity = 0.8, TextWrapping = TextWrapping.Wrap };

            var filling = false; // ItemsSource changes raise SelectionChanged: not a user's pick
            void FillIqns()
            {
                filling = true;
                iqns.ItemsSource = _iqns.Count == 0 ? null : new[] { Messages.SELECT_TARGET_IQN }.Concat(_iqns.Select(SrRules.IqnLabel)).ToList();
                iqns.SelectedIndex = _iqn is IScsiIqnInfo i ? _iqns.IndexOf(i) + 1 : _iqns.Count > 0 ? 0 : -1;
                iqns.IsEnabled = _iqns.Count > 0 && !_busy;
                filling = false;
            }
            void FillLuns()
            {
                filling = true;
                luns.ItemsSource = _luns.Count == 0 ? null : new[] { Messages.SELECT_TARGET_LUN }.Concat(_luns.Select(SrRules.LunLabel)).ToList();
                luns.SelectedIndex = _lun is ISCSIInfo l ? _luns.FindIndex(x => x.ScsiID == l.ScsiID && x.LunID == l.LunID) + 1 : _luns.Count > 0 ? 0 : -1;
                luns.IsEnabled = _luns.Count > 0 && !_busy;
                filling = false;
            }
            // a changed target or CHAP setting invalidates what was scanned (ResetAll)
            void Reset()
            {
                _scanGeneration++;
                _iqns.Clear(); _luns.Clear(); _iqn = null; _lun = null;
                FillIqns(); FillLuns();
                scan.IsEnabled = !string.IsNullOrWhiteSpace(_iscsiHost) && ushort.TryParse(_iscsiPort, out _);
                status.Text = "";
                UpdateFinish();
            }
            host.TextChanged += (_, _) => { _iscsiHost = (host.Text ?? "").Trim(); Reset(); };
            port.TextChanged += (_, _) => { _iscsiPort = (port.Text ?? "").Trim(); Reset(); };
            chap.IsCheckedChanged += (_, _) => { _chap = chap.IsChecked == true; user.IsEnabled = pass.IsEnabled = _chap; Reset(); };
            user.TextChanged += (_, _) => { _chapUser = user.Text ?? ""; Reset(); };
            pass.TextChanged += (_, _) => { _chapPassword = pass.Text ?? ""; Reset(); };

            scan.Click += async (_, _) =>
            {
                if (!ushort.TryParse(_iscsiPort, out var p)) return;
                var generation = ++_scanGeneration;
                scan.IsEnabled = false;
                status.Text = string.Format(Messages.ACTION_ISCSI_IQN_SCANNING, _iscsiHost);
                var action = SrRules.ScanIqns(_c, _iscsiHost, p, ChapUser, ChapPassword);
                await Run(action);
                if (generation != _scanGeneration) return; // the target changed meanwhile
                scan.IsEnabled = true;
                _iqns = action.Succeeded ? action.IQNs.Where(i => !string.IsNullOrEmpty(i.TargetIQN)).ToList() : new List<IScsiIqnInfo>();
                _iqn = _iqns.Count == 1 ? _iqns[0] : null;
                status.Text = !action.Succeeded ? SrRules.IscsiError(action) : _iqns.Count == 0 ? Messages.NEWSR_NO_IQNS_FOUND : "";
                FillIqns();
                if (_iqn != null) await ScanLuns();
            };
            async System.Threading.Tasks.Task ScanLuns()
            {
                _luns.Clear(); _lun = null; FillLuns();
                if (!(_iqn is IScsiIqnInfo iqn)) { UpdateFinish(); return; }
                var generation = _scanGeneration;
                status.Text = string.Format(Messages.ACTION_ISCSI_LUN_SCANNING, iqn.IpAddress);
                var action = SrRules.ScanLuns(_c, iqn, ChapUser, ChapPassword);
                await Run(action);
                if (generation != _scanGeneration || !Equals(_iqn, iqn)) return;
                _luns = action.Succeeded ? action.LUNs.ToList() : new List<ISCSIInfo>();
                _lun = _luns.Count == 1 ? _luns[0] : null;
                status.Text = !action.Succeeded ? SrRules.IscsiError(action) : _luns.Count == 0 ? Messages.NEWSR_NO_LUNS_FOUND : "";
                FillLuns();
                UpdateFinish();
            }
            iqns.SelectionChanged += async (_, _) =>
            {
                if (filling) return;
                var i = iqns.SelectedIndex - 1;
                var pick = i >= 0 && i < _iqns.Count ? _iqns[i] : (IScsiIqnInfo?)null;
                if (Equals(pick, _iqn)) return;
                _iqn = pick;
                await ScanLuns();
            };
            luns.SelectionChanged += (_, _) =>
            {
                if (filling) return;
                var i = luns.SelectedIndex - 1;
                _lun = i >= 0 && i < _luns.Count ? _luns[i] : null;
                UpdateFinish();
            };

            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*,Auto,Auto"), RowSpacing = 8, ColumnSpacing = 8 };
            void Row(int r, string label, Control field, int span = 1)
            {
                while (grid.RowDefinitions.Count <= r) grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var l = new Label { Content = Dialogs.M(label), Target = field, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
                Grid.SetRow(l, r); Grid.SetRow(field, r); Grid.SetColumn(field, 1); Grid.SetColumnSpan(field, span);
                grid.Children.Add(l); grid.Children.Add(field);
            }
            Row(0, "Target &host name/IP address:", host);
            var colon = new TextBlock { Text = ":", VerticalAlignment = VerticalAlignment.Center };
            Grid.SetColumn(colon, 2); Grid.SetColumn(port, 3);
            grid.Children.Add(colon); grid.Children.Add(port);
            grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            Grid.SetRow(chap, 1); Grid.SetColumnSpan(chap, 4);
            grid.Children.Add(chap);
            Row(2, "CHAP &username:", user, 3);
            Row(3, "CHAP pass&word:", pass, 3);
            grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            scan.HorizontalAlignment = HorizontalAlignment.Right;
            Grid.SetRow(scan, 4); Grid.SetColumn(scan, 1); Grid.SetColumnSpan(scan, 3);
            grid.Children.Add(scan);
            Row(5, "Target &IQN:", iqns, 3);
            Row(6, "Target &LUN:", luns, 3);

            var group = new StackPanel { Spacing = 4 };
            group.Children.Add(new TextBlock { Text = "iSCSI target", FontWeight = FontWeight.SemiBold });
            group.Children.Add(new Border { BorderBrush = Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1), Padding = new Thickness(10), Child = grid });
            var p = new StackPanel { Spacing = 12 };
            p.Children.Add(Text("Provide a target host for your iSCSI storage, indicating your target IQN and your target LUN before proceeding."));
            p.Children.Add(group);
            p.Children.Add(status);
            scan.IsEnabled = !string.IsNullOrWhiteSpace(_iscsiHost);
            FillIqns(); FillLuns();
            host.AttachedToVisualTree += (_, _) => host.Focus();
            return p;
        }

        /// <summary>LVMoISCSI.PageLeave: probe the LUN, then reattach the SR found on it or format (LVMoIsciWarningDialog).</summary>
        private async System.Threading.Tasks.Task<(bool Go, string Uuid)> ProbeIscsiLun(Dictionary<string, string> dconf)
        {
            var probe = SrRules.Probe(_c, SrRules.Kind.Iscsi, dconf);
            await Run(probe);
            if (!probe.Succeeded)
            {
                _error.Text = SrRules.IscsiError(probe);
                return (false, null);
            }
            var type = SR.GetFriendlyTypeName(SR.SRTypes.lvmoiscsi);
            var warning = string.Format("To prevent data loss you must ensure that the LUN is not in use by any other system, including {0} hosts that are not connected to {1}.",
                XenAdmin.Core.BrandManager.ProductBrand, XenAdmin.Core.BrandManager.BrandConsole);
            var found = probe.SRs?.FirstOrDefault();
            if (found == null)
                return (await Dialogs.Confirm(this, Title, $"{Messages.LVMOHBA_WARNING_DIALOG_HEADER_NO_EXISTING_SRS}\n\n{string.Format(Messages.LVMOHBA_WARNING_DIALOG_FORMAT_INFO, type)}\n\n{warning}",
                    Dialogs.M("&Format"), Messages.CANCEL, cancelIsDefault: true), null);
            var choice = await Dialogs.Choose(this, Title,
                $"{string.Format(Messages.LVMOHBA_WARNING_DIALOG_HEADER_FOUND_EXISTING_SR, type)}\n\n{found}\n\n" +
                $"Reattach: {string.Format(Messages.LVMOHBA_WARNING_DIALOG_REATTACH_INFO, type)}\nFormat: {string.Format(Messages.LVMOHBA_WARNING_DIALOG_FORMAT_INFO, type)}\n\n{warning}",
                "&Reattach", "&Format");
            return choice switch { true => (true, found.UUID), false => (true, null), _ => (false, null) };
        }

        private static System.Threading.Tasks.Task Run(AsyncAction action)
        {
            var done = new TaskCompletionSource();
            action.Completed += _ => Avalonia.Threading.Dispatcher.UIThread.Post(() => done.TrySetResult());
            action.RunAsync();
            return done.Task;
        }

        private async System.Threading.Tasks.Task Finish()
        {
            if (_kind == SrRules.Kind.Iscsi)
            {
                var dconf = SrRules.IscsiDeviceConfig(_iqn.Value, _lun.Value, ChapUser, ChapPassword);
                _busy = true;
                UpdateFinish();
                var (go, existingUuid) = await ProbeIscsiLun(dconf);
                _busy = false;
                UpdateFinish();
                if (go)
                    await Submit(SrRules.Finish(_c, _kind, _name.Trim(), _autoDescription ? SrRules.IscsiDescription(_iqn.Value, _lun.Value) : _description, dconf, existingUuid), existingUuid);
                return;
            }
            if (SrRules.IsIso(_kind))
            {
                var dconf = SrRules.IsoDeviceConfig(_kind, _path, _nfsVersion == "4", _otherUser ? _user : null, _otherUser ? _password : null);
                await Submit(SrRules.Finish(_c, _kind, _name.Trim(), _autoDescription ? SrRules.AutoDescription(_kind, _path) : _description, dconf, null), null);
                return;
            }
            var uuid = _reattach ? _existing?.UUID : null;
            if (uuid != null && SrRules.InUse(uuid) != null)
                return;
            var warning = SrRules.AttachWarning(_kind, _c, uuid);
            if (warning != null && !await Dialogs.Confirm(this, Title, warning, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION, cancelIsDefault: true))
                return;
            var description = _autoDescription ? SrRules.AutoDescription(_kind, _path) : _description;
            await Submit(SrRules.Finish(_c, _kind, _name.Trim(), description, SrRules.DeviceConfig(_kind, _path, _options, _nfsVersion, _user, _password), uuid), uuid);
        }

        /// <summary>Runs the create/reattach/introduce action; the wizard closes only when it succeeds.</summary>
        private async System.Threading.Tasks.Task Submit(AsyncAction action, string uuid)
        {
            _busy = true;
            UpdateFinish();
            _error.Text = "";
            var done = new System.Threading.Tasks.TaskCompletionSource();
            action.Completed += _ => Avalonia.Threading.Dispatcher.UIThread.Post(() => done.TrySetResult());
            _run(action);
            await done.Task;
            if (action.Succeeded)
            {
                Close();
                return;
            }
            if (uuid != null && SrRules.RollbackFailedReattach(_c, uuid) is AsyncAction rollback)
                await Run(rollback);
            _busy = false;
            UpdateFinish();
            _error.Text = action.Exception?.Message ?? "The storage repository could not be added.";
        }
    }
}
