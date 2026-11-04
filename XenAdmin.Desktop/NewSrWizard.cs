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
    /// NewSRWizard: Type, Name, Location (VHDoNFS / CifsFrontend pages). Creates, reattaches or introduces the SR
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
                _ => _kind == SrRules.Kind.Nfs ? Messages.NEWSR_PATH_NFS : Messages.NEWSR_PATH_CIFS
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
            _page.Content = _step switch { 0 => TypePage(), 1 => NamePage(), _ => LocationPage() };
            UpdateFinish();
        }

        private bool Validate()
        {
            string err = _step switch
            {
                1 when string.IsNullOrWhiteSpace(_name) => "Enter a name for the new SR.",
                2 when !SrRules.ValidLocation(_kind, _path) => _kind == SrRules.Kind.Nfs ? "Enter the share as server:/path." : @"Enter the share as \\server\share.",
                2 when _reattach && _existing == null => "Select the SR to reattach.",
                2 when _existing != null && SrRules.InUse(_existing.UUID) != null => "This SR is already in use by " + XenAdmin.Core.Helpers.GetName(SrRules.InUse(_existing.UUID).Connection) + ".",
                _ => null
            };
            _error.Text = err ?? "";
            return err == null;
        }

        private void UpdateFinish()
        {
            _finish.IsEnabled = !_busy && SrRules.ValidLocation(_kind, _path) && (!_reattach || _existing != null);
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
            list.Children.Add(Option(Messages.NEWSR_LVMOISCSI_TYPE_NAME, null));
            list.Children.Add(Option(Messages.NEWSR_LVMOHBA_TYPE_NAME, null));
            list.Children.Add(Option(Messages.NEWSR_LVMOFCOE_TYPE_NAME, null));
            list.Children.Add(Group("File based storage", 12));
            list.Children.Add(Option(SrRules.TypeName(SrRules.Kind.Nfs), SrRules.Kind.Nfs));
            list.Children.Add(Option(SrRules.TypeName(SrRules.Kind.Smb), SrRules.Kind.Smb));
            list.Children.Add(Option("Local Storage", null));
            list.Children.Add(Group("ISO library", 0));
            list.Children.Add(Option(Messages.NEWSR_CIFS_ISO_TYPE_NAME, null));
            list.Children.Add(Option(Messages.NEWSR_NFS_ISO_TYPE_NAME, null));
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
            var nfs = _kind == SrRules.Kind.Nfs;
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
                var l = new Label { Content = Dialogs.M(label), Target = field, VerticalAlignment = top ? VerticalAlignment.Top : VerticalAlignment.Center, Padding = new Thickness(0, top ? 2 : 0, 0, 0) };
                Grid.SetRow(l, r); Grid.SetRow(field, r); Grid.SetColumn(field, 1);
                fields.Children.Add(l); fields.Children.Add(field);
                if (extra != null) { Grid.SetRow(extra, r); Grid.SetColumn(extra, 2); fields.Children.Add(extra); }
            }
            Row(0, "&Share Name:", path, scan);
            var example = new TextBlock { Text = nfs ? "Example: server:/path" : @"Example: \\server\share", Margin = new Thickness(0, -4, 0, 0) };
            fields.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            Grid.SetRow(example, 1); Grid.SetColumn(example, 1);
            fields.Children.Add(example);
            var p = new StackPanel { Spacing = 10 };
            p.Children.Add(Text(nfs
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
                    var r = new RadioButton { Content = Dialogs.M(text), GroupName = "nfsver", IsChecked = _nfsVersion == value };
                    r.IsCheckedChanged += (_, _) => { if (r.IsChecked == true) _nfsVersion = value; };
                    versions.Children.Add(r);
                }
                Row(2, "NFS Version:", versions);
                Row(3, "Advanced &Options:", options);
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
            p.Children.Add(createNew);
            p.Children.Add(reattach);
            p.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Margin = new Thickness(22, 0, 0, 0), Child = existing });

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

        private static System.Threading.Tasks.Task Run(AsyncAction action)
        {
            var done = new TaskCompletionSource();
            action.Completed += _ => Avalonia.Threading.Dispatcher.UIThread.Post(() => done.TrySetResult());
            action.RunAsync();
            return done.Task;
        }

        private async System.Threading.Tasks.Task Finish()
        {
            var uuid = _reattach ? _existing?.UUID : null;
            if (uuid != null && SrRules.InUse(uuid) != null)
                return;
            var warning = SrRules.AttachWarning(_kind, _c, uuid);
            if (warning != null && !await Dialogs.Confirm(this, Title, warning, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION, cancelIsDefault: true))
                return;
            var description = _autoDescription ? SrRules.AutoDescription(_kind, _path) : _description;
            var action = SrRules.Finish(_c, _kind, _name.Trim(), description, SrRules.DeviceConfig(_kind, _path, _options, _nfsVersion, _user, _password), uuid);
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
