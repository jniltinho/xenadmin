using Avalonia.Input;
using System.Collections.Generic;
using System;
using System.Linq;
using System.Threading.Tasks;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;

namespace XenAdmin.Desktop
{
    public static class Dialogs
    {
        private static Window Frame(string title, Control content, double width = 460)
        {
            var w = new Window
            {
                Title = title,
                Width = width,
                SizeToContent = SizeToContent.Height,
                CanResize = false,
                WindowStartupLocation = WindowStartupLocation.CenterOwner,
                Content = new Border { Padding = new Thickness(14, 12), Child = content }
            };
            // Esc cancels from any field, as in WinForms (a focused TextBox would otherwise keep the key)
            w.AddHandler(InputElement.KeyDownEvent, (_, e) => { if (e.Key == Avalonia.Input.Key.Escape) { e.Handled = true; w.Close(); } }, Avalonia.Interactivity.RoutingStrategies.Tunnel);
            return w;
        }

        private static StackPanel Buttons(params Button[] buttons)
        {
            var p = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, HorizontalAlignment = HorizontalAlignment.Right, Margin = new Thickness(0, 4, 0, 0) };
            p.Children.AddRange(buttons);
            return p;
        }

        public static async Task<bool> Confirm(Window owner, string title, string message, string okText = "OK", string cancelText = "Cancel", bool cancelIsDefault = false)
        {
            var ok = new Button { Content = M(okText), IsDefault = !cancelIsDefault, MinWidth = 75 };
            var cancel = new Button { Content = M(cancelText), IsCancel = true, IsDefault = cancelIsDefault, MinWidth = 75 };
            var panel = new StackPanel();
            panel.Children.Add(new TextBlock { Text = message, TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame(title, panel);
            if (cancelIsDefault) w.Opened += (_, _) => cancel.Focus();
            ok.Click += (_, _) => w.Close(true);
            cancel.Click += (_, _) => w.Close(false);
            return await w.ShowDialog<bool>(owner);
        }

        /// <summary>Three answers (ExitMaintenanceModeDialog: Restore VMs / Skip / Cancel): true, false, or null when cancelled.</summary>
        public static async Task<bool?> Choose(Window owner, string title, string message, string yesText, string noText)
        {
            var yes = new Button { Content = M(yesText), IsDefault = true, MinWidth = 75 };
            var no = new Button { Content = M(noText), MinWidth = 75 };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true, MinWidth = 75 };
            var panel = new StackPanel();
            panel.Children.Add(new TextBlock { Text = message, TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            panel.Children.Add(Buttons(yes, no, cancel));
            var w = Frame(title, panel);
            bool? result = null;
            yes.Click += (_, _) => { result = true; w.Close(); };
            no.Click += (_, _) => { result = false; w.Close(); };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        public static Task Info(Window owner, string title, string message)
        {
            var ok = new Button { Content = "OK", IsDefault = true, IsCancel = true };
            var panel = new StackPanel();
            panel.Children.Add(new SelectableTextBlock { Text = message, TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            panel.Children.Add(Buttons(ok));
            var w = Frame(title, panel);
            ok.Click += (_, _) => w.Close();
            return w.ShowDialog(owner);
        }

        public class ServerInput
        {
            public string Host;
            public int Port;
            public string User;
            public string Password;
        }

        /// <summary>WinForms access keys use '&amp;', Avalonia uses '_'.</summary>
        public static string M(string text) => (text ?? "").Replace("&", "_");

        /// <summary>AddServerDialog: "Add New Server", or "Connect to Server" for a known, disconnected server.</summary>
        public static async Task<ServerInput> AddServer(Window owner, string host = "", string user = "root", bool connectMode = false)
        {
            // Server: an editable field when adding; a disabled combo-like box for a known server (as in WinForms)
            var hostBox = new TextBox { Text = host };
            Control serverField = hostBox;
            if (connectMode)
                serverField = new ComboBox { ItemsSource = new[] { host }, SelectedIndex = 0, IsEnabled = false, HorizontalAlignment = HorizontalAlignment.Stretch };
            var userBox = new TextBox { Text = user };
            var passBox = new TextBox { PasswordChar = '•' };
            var error = new TextBlock { Foreground = Avalonia.Media.Brushes.Red, TextWrapping = Avalonia.Media.TextWrapping.Wrap, IsVisible = false };
            error.PropertyChanged += (_, e) => { if (e.Property == TextBlock.TextProperty) error.IsVisible = !string.IsNullOrEmpty(error.Text); };
            var ok = new Button { Content = M(connectMode ? Messages.CONNECT : Messages.ADD), IsDefault = true };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };

            Grid Fields(params (string Label, Control Field)[] rows)
            {
                var g = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowSpacing = 6, ColumnSpacing = 6 };
                for (int i = 0; i < rows.Length; i++)
                {
                    g.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                    var l = new Label { Content = M(rows[i].Label), Target = rows[i].Field, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
                    Grid.SetRow(l, i); Grid.SetRow(rows[i].Field, i); Grid.SetColumn(rows[i].Field, 1);
                    g.Children.Add(l); g.Children.Add(rows[i].Field);
                }
                return g;
            }
            var credentials = new StackPanel { Spacing = 6 };
            credentials.Children.Add(new TextBlock { Text = "User login credentials" });
            credentials.Children.Add(Fields(("&Username:", userBox), ("&Password:", passBox)));
            var group = new Border { BorderBrush = Avalonia.Media.Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1), Padding = new Thickness(6, 4, 6, 8), Child = credentials };

            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(new TextBlock
            {
                Text = connectMode ? Messages.CONNECT_TO_SERVER_BLURB
                    : "Enter the host name or IP address of the server you want to add and your user login credentials for that server.",
                TextWrapping = Avalonia.Media.TextWrapping.Wrap
            });
            panel.Children.Add(Fields(("&Server:", serverField)));
            panel.Children.Add(group);
            panel.Children.Add(error);
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame(connectMode ? Messages.CONNECT_TO_SERVER : "Add New Server", panel, 400);
            ServerInput result = null;
            ok.Click += (_, _) =>
            {
                if (!ServerAddress.TryParse(hostBox.Text, out var text, out var port)) { error.Text = "Enter a valid server name or address."; return; }
                if (ServerAddress.IsPlainHttp(port)) { error.Text = "Ports 80 and 8080 are plain HTTP; connect over HTTPS (443)."; return; }
                if (string.IsNullOrWhiteSpace(userBox.Text)) { error.Text = "Enter a user name."; return; }
                result = new ServerInput { Host = text, Port = port, User = userBox.Text.Trim(), Password = passBox.Text ?? "" };
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            w.Opened += (_, _) => (connectMode || !string.IsNullOrEmpty(host) ? passBox : hostBox).Focus();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>ConnectingToServerDialog: shown while a connection is in progress; Cancel ends it.</summary>
        public static Window Connecting(Window owner, string host, System.Action cancel)
        {
            var c = new Button { Content = M(Messages.CANCEL), IsCancel = true, MinWidth = 75 };
            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(new TextBlock { Text = string.Format(Messages.LABEL_ATTEMPT, host) });
            panel.Children.Add(Buttons(c));
            var w = Frame(string.Format(Messages.CONNECTING_NOTICE_TITLE, host), panel, 380);
            c.Click += (_, _) => { cancel(); w.Close(); };
            w.Show(owner);
            return w;
        }

        /// <summary>NewTagDialog ("Edit Tags"): every tag in use on the connection, checked if applied; "New Tag" + Create.</summary>
        public static async Task<System.Collections.Generic.List<string>> EditTags(Window owner, XenAdmin.Network.IXenConnection connection, System.Collections.Generic.IList<string> current)
        {
            var all = new System.Collections.Generic.SortedSet<string>(current, System.StringComparer.CurrentCulture);
            foreach (var o in connection.Cache.XenSearchableObjects)
                foreach (var t in XenAdmin.Model.Tags.GetTags(o) ?? System.Array.Empty<string>())
                    all.Add(t);
            var boxes = new System.Collections.Generic.List<CheckBox>();
            var list = new StackPanel();
            void AddTag(string t, bool on)
            {
                var existing = boxes.Find(b => (string)b.Content == t);
                if (existing != null) { existing.IsChecked = true; return; }
                var b = new CheckBox { Content = t, IsChecked = on };
                boxes.Add(b);
                list.Children.Add(b);
            }
            foreach (var t in all) AddTag(t, current.Contains(t));
            var newTag = new TextBox { Width = 230 };
            var create = new Button { Content = M("C&reate") };
            create.Click += (_, _) => { var t = (newTag.Text ?? "").Trim(); if (t.Length > 0) { AddTag(t, true); newTag.Text = ""; } };
            var top = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            top.Children.Add(new Label { Content = M("&New Tag:"), Target = newTag, VerticalAlignment = VerticalAlignment.Center });
            top.Children.Add(newTag);
            top.Children.Add(create);
            var ok = new Button { Content = M(Messages.OK), IsDefault = true };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            var panel = new StackPanel { Spacing = 10 };
            panel.Children.Add(top);
            panel.Children.Add(new Border { BorderBrush = Avalonia.Media.Brush.Parse("#828790"), BorderThickness = new Thickness(1), Background = Avalonia.Media.Brushes.White, Height = 220, Padding = new Thickness(4), Child = new ScrollViewer { Content = list } });
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame("Edit Tags", panel, 420);
            System.Collections.Generic.List<string> result = null;
            ok.Click += (_, _) => { result = boxes.FindAll(b => b.IsChecked == true).ConvertAll(b => (string)b.Content); w.Close(); };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>FolderChangeDialog ("Change Folder"): no folder, or an existing/new folder path.</summary>
        public static async Task<string> ChangeFolder(Window owner, XenAdmin.Network.IXenConnection connection, string current)
        {
            var folders = new System.Collections.Generic.SortedSet<string>(System.StringComparer.CurrentCulture);
            foreach (var o in connection.Cache.XenSearchableObjects)
                if (!string.IsNullOrEmpty(o.Path)) folders.Add(o.Path);
            if (!string.IsNullOrEmpty(current)) folders.Add(current);
            var none = new RadioButton { Content = M("Do n&ot place this resource in any folder"), GroupName = "f", IsChecked = string.IsNullOrEmpty(current) };
            var choose = new RadioButton { Content = M("Place this resource &in this folder:"), GroupName = "f", IsChecked = !string.IsNullOrEmpty(current) };
            var list = new ListBox { ItemsSource = folders, Height = 180, SelectedItem = string.IsNullOrEmpty(current) ? null : current };
            var path = new TextBox { Text = current ?? "", Watermark = "/Folder/Subfolder" };
            list.SelectionChanged += (_, _) => { if (list.SelectedItem is string f) { path.Text = f; choose.IsChecked = true; } };
            path.TextChanged += (_, _) => { if (!string.IsNullOrWhiteSpace(path.Text)) choose.IsChecked = true; };
            var ok = new Button { Content = M(Messages.OK), IsDefault = true };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            var panel = new StackPanel { Spacing = 8 };
            panel.Children.Add(none);
            panel.Children.Add(choose);
            panel.Children.Add(new Border { BorderBrush = Avalonia.Media.Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = list, Margin = new Thickness(20, 0, 0, 0) });
            var newRow = new DockPanel { Margin = new Thickness(20, 0, 0, 0) };
            var l = new Label { Content = M("&New..."), Target = path, VerticalAlignment = VerticalAlignment.Center, Width = 60 };
            DockPanel.SetDock(l, Dock.Left);
            newRow.Children.Add(l);
            newRow.Children.Add(path);
            panel.Children.Add(newRow);
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame("Change Folder", panel, 440);
            string result = null;
            ok.Click += (_, _) => { result = none.IsChecked == true ? "" : (path.Text ?? "").Trim(); w.Close(); };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>WarningDialog with an optional check box (RevertToSnapshotCommand).</summary>
        public static async Task<(bool Yes, bool Checked)> ConfirmWithCheckbox(Window owner, string title, string message, string checkboxText, bool isChecked)
        {
            var yes = new Button { Content = M(Messages.YES), IsDefault = true };
            var no = new Button { Content = M(Messages.NO), IsCancel = true };
            var box = new CheckBox { Content = M(checkboxText ?? ""), IsChecked = isChecked, IsVisible = checkboxText != null };
            var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 10 };
            row.Children.Add(Views.Icon("000_WarningAlert_h32bit_32", 32));
            row.Children.Add(new TextBlock { Text = message, TextWrapping = Avalonia.Media.TextWrapping.Wrap, MaxWidth = 360, VerticalAlignment = VerticalAlignment.Center });
            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(row);
            panel.Children.Add(box);
            panel.Children.Add(Buttons(yes, no));
            var w = Frame(title, panel, 460);
            var result = false;
            yes.Click += (_, _) => { result = true; w.Close(); };
            no.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return (result, box.IsChecked == true);
        }

        /// <summary>VmSnapshotDialog ("Take Snapshot"): name, description and snapshot mode; returns the action to run.</summary>
        public static async Task<XenAdmin.Actions.VMSnapshotCreateAction> TakeSnapshot(Window owner, XenAPI.VM vm)
        {
            var name = new TextBox();
            var description = new TextBox { AcceptsReturn = true, TextWrapping = Avalonia.Media.TextWrapping.Wrap, Height = 60 };
            var modes = VmRules.SnapshotModesFor(vm);
            var disk = new RadioButton { Content = M("&Snapshot the virtual machine's disks"), GroupName = "mode", IsChecked = true, IsEnabled = modes.Disk };
            var quiesce = new CheckBox
            {
                Content = M("&Quiesce the VM before taking the snapshot (Windows only)"), Margin = new Thickness(22, 0, 0, 0), IsVisible = modes.QuiesceVisible, IsEnabled = modes.Quiesce
            };
            var memory = new RadioButton { Content = M("Sn&apshot the virtual machine's disks and memory"), GroupName = "mode", IsEnabled = modes.Memory };
            if (!disk.IsEnabled) { disk.IsChecked = false; memory.IsChecked = memory.IsEnabled; }
            var memoryInfo = new TextBlock { Text = modes.MemoryInfo, TextWrapping = Avalonia.Media.TextWrapping.Wrap, Opacity = 0.75, Margin = new Thickness(22, 0, 0, 0), IsVisible = !memory.IsEnabled };
            var ok = new Button { Content = M("&Take Snapshot"), IsDefault = true };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            void Update() => ok.IsEnabled = !string.IsNullOrWhiteSpace(name.Text) && (disk.IsEnabled || quiesce.IsEnabled || memory.IsEnabled);
            name.TextChanged += (_, _) => Update();
            Update();

            var fields = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowDefinitions = new RowDefinitions("Auto,Auto"), RowSpacing = 8, ColumnSpacing = 10 };
            var ln = new Label { Content = M("&Name:"), Target = name, Padding = new Thickness(0, 4, 0, 0) };
            var ld = new Label { Content = M("&Description:"), Target = description, Padding = new Thickness(0, 4, 0, 0) };
            Grid.SetColumn(name, 1); Grid.SetRow(ld, 1); Grid.SetRow(description, 1); Grid.SetColumn(description, 1);
            fields.Children.Add(ln); fields.Children.Add(name); fields.Children.Add(ld); fields.Children.Add(description);
            var mode = new StackPanel { Spacing = 6 };
            mode.Children.Add(new TextBlock { Text = "Snapshot mode" });
            mode.Children.Add(disk);
            mode.Children.Add(quiesce);
            mode.Children.Add(memory);
            mode.Children.Add(memoryInfo);
            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(fields);
            panel.Children.Add(new Border { BorderBrush = Avalonia.Media.Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1), Padding = new Thickness(8, 4, 8, 8), Child = mode });
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame("Take Snapshot", panel, 480);
            XenAdmin.Actions.VMSnapshotCreateAction result = null;
            ok.Click += (_, _) =>
            {
                var type = quiesce.IsChecked == true && quiesce.IsEnabled ? XenAdmin.Actions.SnapshotType.QUIESCED_DISK
                    : disk.IsChecked == true ? XenAdmin.Actions.SnapshotType.DISK : XenAdmin.Actions.SnapshotType.DISK_AND_MEMORY;
                // no console screenshot provider yet: the snapshot is saved without the VNC thumbnail
                result = new XenAdmin.Actions.VMSnapshotCreateAction(vm, name.Text.Trim(), description.Text ?? "", type, null);
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            w.Opened += (_, _) => name.Focus();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>
        /// VIFDialog ("Add Virtual Interface" / "Virtual Interface Properties"): network, MAC and QoS limit. Returns the
        /// new settings, or null when cancelled or nothing changed (NetworkRules has the validation).
        /// </summary>
        public static async Task<XenAPI.VIF> Vif(Window owner, XenAPI.VM vm, XenAPI.VIF existing, int device)
        {
            var networks = NetworkRules.Networks(vm.Connection, vm);
            var network = new ComboBox { ItemsSource = networks.Count > 0 ? networks.Select(n => n.Name()).ToList() : new List<string> { Messages.NONE }, SelectedIndex = 0, HorizontalAlignment = HorizontalAlignment.Stretch };
            var auto = new RadioButton { Content = M("&Auto-generate a MAC address"), GroupName = "mac", IsChecked = true };
            var manual = new RadioButton { Content = M("&Use this MAC address:"), GroupName = "mac" };
            var mac = new TextBox { Watermark = "aa:bb:cc:dd:ee:ff", Width = 200 };
            var qos = new CheckBox { Content = M("Enable a &QoS limit of:") };
            var limit = new TextBox { Width = 90 };
            var error = new TextBlock { Foreground = Avalonia.Media.Brushes.Red, TextWrapping = Avalonia.Media.TextWrapping.Wrap, IsVisible = false };
            var ok = new Button { Content = M(existing == null ? "&Add" : Messages.OK), IsDefault = true };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            if (existing != null)
            {
                network.SelectedIndex = Math.Max(0, networks.FindIndex(n => n.opaque_ref == existing.network.opaque_ref));
                mac.Text = existing.MAC;
                manual.IsChecked = !string.IsNullOrEmpty(existing.MAC);
                auto.IsChecked = !manual.IsChecked;
                limit.Text = existing.LimitString();
                qos.IsChecked = existing.qos_algorithm_type == XenAPI.VIF.RATE_LIMIT_QOS_VALUE;
            }
            // VIFDialog.LoadDetails: QoS is set on the vSwitch Controller (pre-Stockholm pools that have one)
            var controller = NetworkRules.QosManagedByController(vm.Connection);
            if (controller)
            {
                qos.IsChecked = false;
                qos.IsEnabled = limit.IsEnabled = false;
            }
            XenAPI.Network Selected() => network.SelectedIndex >= 0 && network.SelectedIndex < networks.Count ? networks[network.SelectedIndex] : null;
            string Mac() => manual.IsChecked == true ? mac.Text?.Trim() ?? "" : "";
            void Validate()
            {
                var e = NetworkRules.Validate(Selected(), manual.IsChecked == true, mac.Text, qos.IsChecked == true, limit.Text);
                ok.IsEnabled = e == null;
                error.Text = e;
                error.IsVisible = !string.IsNullOrEmpty(e);
            }
            network.SelectionChanged += (_, _) => Validate();
            auto.IsCheckedChanged += (_, _) => Validate();
            manual.IsCheckedChanged += (_, _) => Validate();
            mac.GotFocus += (_, _) => manual.IsChecked = true;
            mac.TextChanged += (_, _) => Validate();
            qos.IsCheckedChanged += (_, _) => Validate();
            limit.GotFocus += (_, _) => qos.IsChecked = true;
            limit.TextChanged += (_, _) => Validate();

            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowDefinitions = new RowDefinitions("Auto,Auto,Auto"), RowSpacing = 10, ColumnSpacing = 10 };
            void Row(int r, Control label, Control field)
            {
                Grid.SetRow(label, r); Grid.SetRow(field, r); Grid.SetColumn(field, 1);
                grid.Children.Add(label); grid.Children.Add(field);
            }
            var macPanel = new StackPanel { Spacing = 4 };
            macPanel.Children.Add(auto);
            var manualRow = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
            manualRow.Children.Add(manual);
            manualRow.Children.Add(mac);
            macPanel.Children.Add(manualRow);
            var qosRow = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
            qosRow.Children.Add(qos);
            qosRow.Children.Add(limit);
            qosRow.Children.Add(new TextBlock { Text = "Kbytes/s", VerticalAlignment = VerticalAlignment.Center });
            Row(0, new Label { Content = M("&Network:"), Target = network, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) }, network);
            Row(1, new TextBlock { Text = "MAC address:", Margin = new Thickness(0, 6, 0, 0) }, macPanel);
            Row(2, new TextBlock { Text = "QoS settings:", VerticalAlignment = VerticalAlignment.Center }, qosRow);
            if (controller)
            {
                grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var info = new TextBlock { Text = "QoS settings should be configured using the vSwitch Controller", TextWrapping = Avalonia.Media.TextWrapping.Wrap };
                Grid.SetRow(info, 3); Grid.SetColumn(info, 1);
                grid.Children.Add(info);
            }
            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(new TextBlock { Text = "Select your network and MAC address for this virtual interface and, optionally, specify a QoS limit.", TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            panel.Children.Add(grid);
            panel.Children.Add(error);
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame(existing == null ? "Add Virtual Interface" : Messages.VIRTUAL_INTERFACE_PROPERTIES, panel, 520);
            Validate();
            XenAPI.VIF result = null;
            ok.Click += async (_, _) =>
            {
                var net = Selected();
                if (net == null) return;
                var dup = NetworkRules.DuplicateMac(XenAdmin.ConnectionsManager.XenConnectionsCopy, existing, Mac());
                if (dup != null && !await Confirm(w, Messages.PROBLEM_MAC_ADDRESS_IS_DUPLICATE_TITLE,
                        string.Format(Messages.PROBLEM_MAC_ADDRESS_IS_DUPLICATE, Mac(), dup.NameWithLocation()), Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION, cancelIsDefault: true))
                    return;
                var limitText = limit.Text?.Trim() ?? "";
                if (NetworkRules.Changed(existing, net, Mac(), device, qos.IsChecked == true, limitText))
                    result = NetworkRules.Settings(vm, existing, net, Mac(), device, qos.IsChecked == true, limitText);
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>
        /// MoveVirtualDiskDialog / MigrateVirtualDiskDialog: the SR list with the reason a target cannot take the disk,
        /// Rescan, and Move (MigrateVirtualDiskAction live, MoveVirtualDiskAction offline). Returns the action, or null.
        /// </summary>
        public static async Task<XenAdmin.Actions.AsyncAction> MoveDisk(Window owner, XenAPI.VDI vdi)
        {
            var migrate = StorageRules.UseMigrate(vdi);
            if (migrate == null) return null;
            var list = new ListBox { Height = 220 };
            var move = new Button { Content = M("&Move"), IsDefault = true };
            var rescan = new Button { Content = M("&Rescan") };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            List<StorageRules.SrTarget> targets = null;
            void Fill()
            {
                targets = StorageRules.MoveTargets(vdi, migrate.Value);
                list.ItemsSource = targets.Select(t => new TextBlock { Text = $"{t.Sr.Name()}   {t.Text}", Opacity = t.Enabled ? 1 : 0.5 }).ToList();
                list.SelectedIndex = targets.FindIndex(t => t.Enabled);
                move.IsEnabled = list.SelectedIndex >= 0;
            }
            list.SelectionChanged += (_, _) =>
            {
                // disabled targets stay visible with their reason but cannot be picked (SrPicker)
                if (list.SelectedIndex >= 0 && !targets[list.SelectedIndex].Enabled)
                    list.SelectedIndex = targets.FindIndex(t => t.Enabled);
                move.IsEnabled = list.SelectedIndex >= 0 && targets[list.SelectedIndex].Enabled;
            };
            var panel = new StackPanel { Spacing = 10 };
            panel.Children.Add(new TextBlock { Text = "Select the storage repository that you would like to move the virtual disk to:", TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            panel.Children.Add(new Border { BorderBrush = Avalonia.Media.Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = list });
            var buttons = new DockPanel();
            DockPanel.SetDock(rescan, Dock.Left);
            buttons.Children.Add(rescan);
            buttons.Children.Add(Buttons(move, cancel));
            panel.Children.Add(buttons);
            var w = Frame("Move Virtual Disk", panel, 560);
            Fill();
            XenAdmin.Actions.AsyncAction result = null;
            rescan.Click += (_, _) =>
            {
                // SrPicker.ScanSRs: rescan the attached SRs, then list them again
                rescan.IsEnabled = false;
                var scans = targets.Where(t => !t.Sr.IsDetached()).Select(t => (XenAdmin.Actions.AsyncAction)new XenAdmin.Actions.SrRefreshAction(t.Sr)).ToList();
                var all = new XenAdmin.Actions.MultipleAction(vdi.Connection, Messages.SR_REFRESH_ACTION_TITLE_GENERIC, "", "", scans, true);
                all.Completed += _ => Avalonia.Threading.Dispatcher.UIThread.Post(async () =>
                {
                    Fill();
                    rescan.IsEnabled = true;
                    // the dialog may be closed by now: report against whichever window is still open
                    if (!all.Succeeded && all.Exception != null && (w.IsVisible ? w : owner) is Window shown && shown.IsVisible)
                        await Info(shown, Messages.SR_REFRESH_ACTION_TITLE_GENERIC, all.Exception.Message);
                });
                all.RunAsync();
            };
            move.Click += (_, _) =>
            {
                if (list.SelectedIndex < 0 || !targets[list.SelectedIndex].Enabled) return;
                var sr = targets[list.SelectedIndex].Sr;
                result = migrate.Value
                    ? new XenAdmin.Actions.MigrateVirtualDiskAction(vdi.Connection, vdi, sr)
                    : new XenAdmin.Actions.MoveVirtualDiskAction(vdi.Connection, vdi, sr);
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>
        /// BallooningDialog ("Memory Settings"): one memory size (VMMemoryControlsBasic) or dynamic minimum, dynamic
        /// maximum and static maximum (VMMemoryControlsAdvanced), in MB up to 1 GB and GB above; the restart warning for a
        /// running VM whose static maximum changes. Returns ChangeMemorySettingsAction, or null.
        /// </summary>
        public static async Task<XenAdmin.Actions.AsyncAction> MemorySettings(Window owner, XenAPI.VM vm)
        {
            if (!MemoryRules.CanEdit(vm)) return null;
            var advanced = MemoryRules.Advanced(vm);
            var mb = MemoryRules.InMegabytes(vm);
            var unit = mb ? XenAdmin.Util.BINARY_MEGA : XenAdmin.Util.BINARY_GIGA;
            NumericUpDown Spinner(long bytes) => new NumericUpDown
            {
                Value = (decimal)Math.Round((double)bytes / unit, mb ? 0 : 1), Increment = mb ? 64 : 0.1m, Minimum = 0, FormatString = mb ? "0" : "0.0", Width = 130
            };
            var dmin = Spinner(vm.memory_dynamic_min);
            var dmax = Spinner(vm.memory_dynamic_max);
            var smax = Spinner(vm.memory_static_max);
            // VMMemoryControlsBasic: a fixed size, or a range when ballooning is allowed
            var fixedRadio = new RadioButton { Content = M("&Set a fixed memory of"), GroupName = "mem", IsChecked = vm.memory_dynamic_min == vm.memory_static_max || !MemoryRules.DynamicAllowed(vm) };
            var rangeRadio = new RadioButton { Content = M("&Automatically allocate memory within this range:"), GroupName = "mem", IsChecked = fixedRadio.IsChecked != true, IsEnabled = MemoryRules.DynamicAllowed(vm) };
            var error = new TextBlock { Foreground = Avalonia.Media.Brushes.Red, TextWrapping = Avalonia.Media.TextWrapping.Wrap };
            var ok = new Button { Content = "OK", IsDefault = true };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            // a field left as shown keeps the server's exact bytes (the display rounds, e.g. 1280 MiB shows as 1.3 GB)
            var original = new Dictionary<NumericUpDown, (decimal? Shown, long Bytes)>
            {
                [dmin] = (dmin.Value, vm.memory_dynamic_min), [dmax] = (dmax.Value, vm.memory_dynamic_max), [smax] = (smax.Value, vm.memory_static_max)
            };
            long Bytes(NumericUpDown n) => n.Value == original[n].Shown ? original[n].Bytes : (long)((double)(n.Value ?? 0) * unit);
            (long, long, long) Values() => advanced ? (Bytes(dmin), Bytes(dmax), Bytes(smax))
                : rangeRadio.IsChecked == true ? (Bytes(dmin), Bytes(dmax), Bytes(dmax)) : (Bytes(smax), Bytes(smax), Bytes(smax));
            void Validate()
            {
                var (a, b, c) = Values();
                error.Text = MemoryRules.Validate(vm, a, b, c) ?? "";
                ok.IsEnabled = error.Text == "";
            }
            foreach (var n in new[] { dmin, dmax, smax })
                n.ValueChanged += (_, _) => Validate();
            fixedRadio.IsCheckedChanged += (_, _) => Validate();
            dmin.GotFocus += (_, _) => { if (!advanced && rangeRadio.IsEnabled) rangeRadio.IsChecked = true; };
            dmax.GotFocus += (_, _) => { if (!advanced && rangeRadio.IsEnabled) rangeRadio.IsChecked = true; };
            smax.GotFocus += (_, _) => { if (!advanced) fixedRadio.IsChecked = true; };
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,Auto,Auto"), RowSpacing = 8, ColumnSpacing = 10 };
            void Row(string label, NumericUpDown field)
            {
                var r = grid.RowDefinitions.Count;
                grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var l = new TextBlock { Text = label, VerticalAlignment = VerticalAlignment.Center };
                var u = new TextBlock { Text = mb ? "MB" : "GB", VerticalAlignment = VerticalAlignment.Center };
                Grid.SetRow(l, r); Grid.SetRow(field, r); Grid.SetRow(u, r); Grid.SetColumn(field, 1); Grid.SetColumn(u, 2);
                grid.Children.Add(l); grid.Children.Add(field); grid.Children.Add(u);
            }
            if (advanced)
            {
                Row("Minimum memory:", dmin);
                Row("Maximum memory:", dmax);
                Row("Static maximum:", smax);
            }
            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(new TextBlock { Text = MemoryRules.Rubric(vm), TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            if (advanced)
                panel.Children.Add(grid);
            else
            {
                var fixedRow = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
                fixedRow.Children.Add(fixedRadio);
                fixedRow.Children.Add(smax);
                fixedRow.Children.Add(new TextBlock { Text = mb ? "MB" : "GB", VerticalAlignment = VerticalAlignment.Center });
                panel.Children.Add(fixedRow);
                panel.Children.Add(rangeRadio);
                Row("Minimum memory:", dmin);
                Row("Maximum memory:", dmax);
                grid.Margin = new Thickness(22, 0, 0, 0);
                grid.IsEnabled = rangeRadio.IsEnabled;
                panel.Children.Add(grid);
                if (MemoryRules.DynamicUnavailable(vm) is string why)
                    panel.Children.Add(new TextBlock { Text = why, TextWrapping = Avalonia.Media.TextWrapping.Wrap, Opacity = 0.8 });
            }
            panel.Children.Add(error);
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame("Memory Settings", panel, 480);
            Validate();
            XenAdmin.Actions.AsyncAction result = null;
            ok.Click += async (_, _) =>
            {
                var (a, b, c) = MemoryRules.Settings(vm, Values().Item1, Values().Item2, Values().Item3);
                if (MemoryRules.Confirmation(vm, c) is string warning && !await Confirm(w, "Memory Settings", warning, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION))
                    return;
                result = MemoryRules.ChangeAction(vm, a, b, c);
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>
        /// DateFilterDialog ("Select Date Filter"): a start and an end (not set, 24 hours, 3, 7 or 30 days ago, or a custom
        /// date and time); the end must be after the start. Returns the range, or null when cancelled.
        /// </summary>
        public static async Task<(DateTime? From, DateTime? To)?> DateFilter(Window owner, DateTime? from, DateTime? to)
        {
            (ComboBox Choice, DatePicker Date, TimePicker Time) Bound(DateTime? value)
            {
                var combo = new ComboBox { ItemsSource = NotificationRules.DateBoundChoices, SelectedIndex = value == null ? 0 : NotificationRules.CustomBound, Width = 140 };
                var date = new DatePicker { SelectedDate = value ?? DateTime.Now, MinWidth = 230 };
                var time = new TimePicker { SelectedTime = (value ?? DateTime.Now).TimeOfDay, ClockIdentifier = "24HourClock" };
                return (combo, date, time);
            }
            var start = Bound(from);
            var end = Bound(to);
            var warning = new TextBlock { Text = "End date must be after the start date.", Foreground = Avalonia.Media.Brushes.Red, IsVisible = false };
            var apply = new Button { Content = M("&Apply Filter"), IsDefault = true };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            DateTime? Value((ComboBox Choice, DatePicker Date, TimePicker Time) b, bool isStart, DateTime now)
            {
                var custom = (b.Date.SelectedDate?.Date ?? DateTime.Today) + (b.Time.SelectedTime ?? TimeSpan.Zero);
                return NotificationRules.DateBound(b.Choice.SelectedIndex, isStart, now, custom);
            }
            // both bounds from one clock reading, so "24 hours ago" to "24 hours ago" is never a valid range
            (DateTime? From, DateTime? To) Range() { var now = DateTime.Now; return (Value(start, true, now), Value(end, false, now)); }
            void Validate()
            {
                foreach (var b in new[] { start, end })
                    b.Date.IsEnabled = b.Time.IsEnabled = b.Choice.SelectedIndex == NotificationRules.CustomBound;
                var (f, t) = Range();
                var ok = NotificationRules.ValidRange(f, t);
                warning.IsVisible = !ok;
                apply.IsEnabled = ok;
            }
            foreach (var b in new[] { start, end })
            {
                b.Choice.SelectionChanged += (_, _) => Validate();
                b.Date.SelectedDateChanged += (_, _) => Validate();
                b.Time.SelectedTimeChanged += (_, _) => Validate();
            }
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,Auto,Auto,Auto"), RowDefinitions = new RowDefinitions("Auto,Auto"), RowSpacing = 8, ColumnSpacing = 8 };
            void Row(int r, string label, (ComboBox Choice, DatePicker Date, TimePicker Time) b)
            {
                var l = new Label { Content = M(label), Target = b.Choice, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
                Control[] cells = { l, b.Choice, b.Date, b.Time };
                for (int c = 0; c < cells.Length; c++) { Grid.SetRow(cells[c], r); Grid.SetColumn(cells[c], c); grid.Children.Add(cells[c]); }
            }
            Row(0, "&Start date:", start);
            Row(1, "&End date:", end);
            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(new TextBlock { Text = "Set the filter to show results from a specific time period by altering the start and end dates below", TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            panel.Children.Add(grid);
            panel.Children.Add(warning);
            panel.Children.Add(Buttons(apply, cancel));
            var w = Frame("Select Date Filter", panel, 720);
            Validate();
            (DateTime?, DateTime?)? result = null;
            apply.Click += (_, _) =>
            {
                // validate the exact pair that is applied (time has passed since the last check)
                var range = Range();
                if (!NotificationRules.ValidRange(range.From, range.To)) { Validate(); return; }
                result = range;
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>
        /// ChangeServerPasswordDialog: the current root password (checked against the connection's), the new one twice,
        /// and the pool secret advice. Returns ChangeHostPasswordAction, or null. The passwords stay in memory only.
        /// </summary>
        public static async Task<XenAdmin.Actions.AsyncAction> ChangeServerPassword(Window owner, XenAPI.IXenObject hostOrPool)
        {
            var c = hostOrPool.Connection;
            TextBox Pass() => new TextBox { PasswordChar = '•' };
            var current = Pass();
            var next = Pass();
            var confirm = Pass();
            var error = new TextBlock { Foreground = Avalonia.Media.Brushes.Red, TextWrapping = Avalonia.Media.TextWrapping.Wrap };
            var ok = new Button { Content = "OK", IsDefault = true, IsEnabled = false };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            foreach (var box in new[] { current, next, confirm })
                box.TextChanged += (_, _) =>
                {
                    error.Text = "";
                    ok.IsEnabled = PoolRules.ChangePasswordError(c, current.Text, next.Text, confirm.Text) != "";
                };
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowDefinitions = new RowDefinitions("Auto,Auto,Auto"), RowSpacing = 8, ColumnSpacing = 10 };
            void Row(int r, string label, TextBox box)
            {
                var l = new Label { Content = M(label), Target = box, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
                Grid.SetRow(l, r); Grid.SetRow(box, r); Grid.SetColumn(box, 1);
                grid.Children.Add(l); grid.Children.Add(box);
            }
            Row(0, "&Current root password:", current);
            Row(1, "&New root password:", next);
            Row(2, "&Re-enter new password:", confirm);
            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(grid);
            panel.Children.Add(new TextBlock { Text = "Passwords are case sensitive. If you lose or forget the password, it cannot be recovered.", TextWrapping = Avalonia.Media.TextWrapping.Wrap, Opacity = 0.8 });
            if (PoolRules.CanRotateSecret(c))
                panel.Children.Add(new TextBlock { Text = "To reduce the risk of unauthorized access to your system, it is recommended that you also rotate the pool secret.", TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            panel.Children.Add(error);
            panel.Children.Add(Buttons(ok, cancel));
            var name = XenAdmin.Core.Helpers.GetPoolOfOne(c)?.Name() ?? XenAdmin.Core.Helpers.GetName(c);
            var w = Frame(string.Format(Messages.CHANGEPASS_DIALOG_TITLE, name.Ellipsise(30)), panel, 480);
            XenAdmin.Actions.AsyncAction result = null;
            ok.Click += (_, _) =>
            {
                var err = PoolRules.ChangePasswordError(c, current.Text, next.Text, confirm.Text);
                if (err != null)
                {
                    error.Text = err;
                    (err == Messages.PASSWORD_INCORRECT ? current : next).Focus();
                    return;
                }
                result = PoolRules.ChangePasswordAction(c, current.Text, next.Text);
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>DestroyBondCommand: the warning for the management interfaces (or the HA error), then DestroyBondAction.</summary>
        public static async Task<XenAdmin.Actions.AsyncAction> DeleteBond(Window owner, XenAPI.Network network)
        {
            if (!BondRules.CanDelete(network)) return null;
            var w = BondRules.DeleteWarning(network);
            if (w.Block)
            {
                await Info(owner, w.Title, w.Text);
                return null;
            }
            if (!await Confirm(owner, w.Title, w.Text, w.OkText, Messages.CANCEL))
                return null;
            // DestroyBondCommand: the bond may have gone while the confirmation was open
            return BondRules.BondOf(network) is XenAPI.Bond bond ? new XenAdmin.Actions.DestroyBondAction(bond) : null;
        }

        /// <summary>BondProperties (Create Bond from the NICs tab): the BondDetails control in a dialog.</summary>
        public static async Task<XenAdmin.Actions.AsyncAction> CreateBond(Window owner, XenAPI.Host host)
        {
            var editor = new BondEditor(host);
            var create = new Button { Content = M("&Create"), IsDefault = true, IsEnabled = editor.IsValid };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            editor.Changed += () => create.IsEnabled = editor.IsValid;
            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(editor);
            panel.Children.Add(Buttons(create, cancel));
            var w = Frame(Messages.BOND_CREATE, panel, 620);
            XenAdmin.Actions.AsyncAction result = null;
            create.Click += async (_, _) =>
            {
                result = await editor.CreateAction(w);
                if (result != null) w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>NewDiskDialog ("Add Virtual Disk"): CreateDiskAction + VbdCreateAndPlugAction, as the original.</summary>
        public static async Task<XenAdmin.Actions.AsyncAction> NewDisk(Window owner, XenAPI.VM vm) => await NewDisk(owner, vm, null);

        /// <summary>NewDiskDialog(connection, sr) from the SR Storage tab: a disk on that SR, not attached to any VM.</summary>
        public static async Task<XenAdmin.Actions.AsyncAction> NewDisk(Window owner, XenAPI.SR sr) => await NewDisk(owner, null, sr);

        private static async Task<XenAdmin.Actions.AsyncAction> NewDisk(Window owner, XenAPI.VM vm, XenAPI.SR target)
        {
            var c = vm?.Connection ?? target.Connection;
            var name = new TextBox { Text = XenAdmin.Core.Helpers.MakeUniqueName(Messages.DEFAULT_VDI_NAME, c.Cache.VDIs.Select(v => v.Name()).ToList()) };
            var description = new TextBox();
            var size = new NumericUpDown { Value = 10, Minimum = 0.01m, Maximum = 2040, Increment = 1, FormatString = "0.##", Width = 110 };
            var unit = new ComboBox { ItemsSource = new[] { "GB", "MB" }, SelectedIndex = 0, Width = 70 };
            var home = vm?.Home();
            var srs = StorageRules.NewDiskTargets(c, home);  // the SR picker preselects target only when it can hold a disk
            var srList = new ListBox { ItemsSource = srs.Select(sr => $"{sr.NameWithLocation()}   {Util.DiskSizeString(sr.FreeSpace())} free of {Util.DiskSizeString(sr.physical_size)}").ToList(), Height = 140 };
            var def = target ?? c.Resolve(XenAdmin.Core.Helpers.GetPoolOfOne(c)?.default_SR);
            srList.SelectedIndex = Math.Max(0, srs.FindIndex(sr => def != null && sr.opaque_ref == def.opaque_ref));
            var error = new TextBlock { Foreground = Avalonia.Media.Brushes.Red, TextWrapping = Avalonia.Media.TextWrapping.Wrap };
            var ok = new Button { Content = M("&Add"), IsDefault = true };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowDefinitions = new RowDefinitions("Auto,Auto,Auto,Auto"), RowSpacing = 8, ColumnSpacing = 10 };
            void Row(int r, string label, Control field)
            {
                var l = new Label { Content = M(label), Target = field, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
                Grid.SetRow(l, r); Grid.SetRow(field, r); Grid.SetColumn(field, 1);
                grid.Children.Add(l); grid.Children.Add(field);
            }
            var sizeRow = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
            sizeRow.Children.Add(size);
            sizeRow.Children.Add(unit);
            Row(0, "&Name:", name);
            Row(1, "&Description:", description);
            Row(2, "&Size:", sizeRow);
            Row(3, "&Location:", new Border { BorderBrush = Avalonia.Media.Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = srList });
            var panel = new StackPanel { Spacing = 10 };
            panel.Children.Add(new TextBlock { Text = "Enter a name, description and size for the new virtual disk, and select the storage repository where it will be stored.", TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            panel.Children.Add(grid);
            panel.Children.Add(error);
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame(Messages.DISK_ADD, panel, 560);
            XenAdmin.Actions.AsyncAction result = null;
            ok.Click += (_, _) =>
            {
                var bytes = (long)((size.Value ?? 0) * (unit.SelectedIndex == 0 ? Util.BINARY_GIGA : Util.BINARY_MEGA));
                var sr = srList.SelectedIndex >= 0 ? srs[srList.SelectedIndex] : null;
                if (string.IsNullOrWhiteSpace(name.Text)) { error.Text = "Enter a name for the disk."; return; }
                if (sr == null) { error.Text = "Select a storage repository."; return; }
                if (bytes <= 0) { error.Text = "Enter a disk size."; return; }
                var vdi = new XenAPI.VDI
                {
                    Connection = c, read_only = false, SR = new XenAPI.XenRef<XenAPI.SR>(sr), virtual_size = bytes,
                    name_label = name.Text.Trim(), name_description = description.Text ?? "", sharable = false, type = XenAPI.vdi_type.user
                };
                if (vm == null)
                {
                    result = new XenAdmin.Actions.CreateDiskAction(vdi);
                    w.Close();
                    return;
                }
                vdi.SetVmHint(vm.uuid);
                var vbd = new XenAPI.VBD { Connection = c, device = "", empty = false, type = XenAPI.vbd_type.Disk, mode = XenAPI.vbd_mode.RW, unpluggable = true };
                vbd.SetIsOwner(true);
                var actions = new System.Collections.Generic.List<XenAdmin.Actions.AsyncAction>
                {
                    new XenAdmin.Actions.CreateDiskAction(vdi, vbd, vm),   // fills in the VBD
                    new XenAdmin.Actions.VbdCreateAndPlugAction(vm, vbd, vdi.Name(), false)
                };
                result = new XenAdmin.Actions.MultipleAction(c, $"Add virtual disk '{vdi.Name()}' to '{vm.Name()}'", "", "", actions, false, false, true);
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            w.Opened += (_, _) => name.Focus();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>AttachDiskDialog: an existing, unattached disk; VbdCreateAndPlugAction on a free userdevice.</summary>
        public static async Task<XenAdmin.Actions.AsyncAction> AttachDisk(Window owner, XenAPI.VM vm)
        {
            var c = vm.Connection;
            var home = vm.Home();
            var attached = c.ResolveAll(vm.VBDs).Select(b => b.VDI.opaque_ref).ToHashSet();
            var vdis = c.Cache.SRs.Where(sr => sr.Show(false) && sr.content_type != XenAPI.SR.Content_Type_ISO && !sr.IsToolsSR() && (home == null || sr.CanBeSeenFrom(home)))
                .SelectMany(sr => c.ResolveAll(sr.VDIs).Where(v => v.Show(false) && !v.is_a_snapshot && !attached.Contains(v.opaque_ref) && (v.sharable || v.VBDs.Count == 0)))
                .OrderBy(v => v.Name()).ToList();
            var list = new ListBox
            {
                ItemsSource = vdis.Select(v => $"{v.Name()}  ({v.SizeText()})  —  {c.Resolve(v.SR)?.Name()}").ToList(), Height = 220
            };
            var ro = new CheckBox { Content = M("Attach as &read-only") };
            var ok = new Button { Content = M("&Attach"), IsDefault = true, IsEnabled = false };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            list.SelectionChanged += (_, _) => ok.IsEnabled = list.SelectedIndex >= 0;
            var panel = new StackPanel { Spacing = 10 };
            panel.Children.Add(new TextBlock { Text = vdis.Count == 0 ? "There are no unattached virtual disks available to this VM." : "Select the virtual disk you want to attach to this VM.", TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            panel.Children.Add(new Border { BorderBrush = Avalonia.Media.Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = list });
            panel.Children.Add(ro);
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame(Messages.DISK_ATTACH, panel, 520);
            XenAPI.VDI chosen = null;
            ok.Click += (_, _) => { chosen = vdis[list.SelectedIndex]; w.Close(); };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            if (chosen == null)
                return null;
            var readOnly = ro.IsChecked == true;
            // a free device is a server call: done off the UI thread, as the original
            var uds = await System.Threading.Tasks.Task.Run(() => XenAPI.VM.get_allowed_VBD_devices(c.DuplicateSession(), vm.opaque_ref));
            if (uds.Length == 0)
            {
                await Info(owner, Messages.DISK_ATTACH, "This VM already has the maximum number of virtual disks.");
                return null;
            }
            var ud = uds[0];
            var vbd = new XenAPI.VBD
            {
                VDI = new XenAPI.XenRef<XenAPI.VDI>(chosen), VM = new XenAPI.XenRef<XenAPI.VM>(vm), bootable = ud == "0", device = "", empty = false,
                userdevice = ud, type = XenAPI.vbd_type.Disk, mode = readOnly ? XenAPI.vbd_mode.RO : XenAPI.vbd_mode.RW, unpluggable = true
            };
            vbd.SetIsOwner(chosen.VBDs.Count == 0);
            return new XenAdmin.Actions.VbdCreateAndPlugAction(vm, vbd, chosen.Name(), false);
        }

        /// <summary>CopyVMDialog: fast clone (VMCloneAction) or full copy to an SR (VMCopyAction).</summary>
        public static async Task<XenAdmin.Actions.AsyncAction> CopyVm(Window owner, XenAPI.VM vm)
        {
            var c = vm.Connection;
            var name = new TextBox { Text = XenAdmin.Core.Helpers.MakeUniqueName(string.Format(Messages.ACTION_TEMPLATE_CLONE_NEW_NAME, vm.Name()), c.Cache.VMs.Select(v => v.Name()).ToList()) };
            var description = new TextBox { Text = vm.Description(), AcceptsReturn = true, TextWrapping = Avalonia.Media.TextWrapping.Wrap, Height = 50 };
            var modes = VmRules.CopyModes(vm);
            var clone = new RadioButton { Content = M("&Fast clone"), GroupName = "copy", IsEnabled = modes.Clone };
            var copy = new RadioButton { Content = M("F&ull copy"), GroupName = "copy", IsEnabled = modes.Copy };
            clone.IsChecked = clone.IsEnabled;
            copy.IsChecked = !clone.IsEnabled;
            var host = vm.GetStorageHost(false);
            var srs = StorageRules.NewDiskTargets(c, host);
            var srList = new ListBox { ItemsSource = srs.Select(sr => $"{sr.NameWithLocation()}   {Util.DiskSizeString(sr.FreeSpace())} free").ToList(), Height = 120, Margin = new Thickness(22, 0, 0, 0) };
            srList.SelectedIndex = srs.Count > 0 ? 0 : -1;
            void Update() => srList.IsEnabled = copy.IsChecked == true;
            copy.IsCheckedChanged += (_, _) => Update();
            Update();
            var ok = new Button { Content = M("C&opy"), IsDefault = true };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowDefinitions = new RowDefinitions("Auto,Auto"), RowSpacing = 8, ColumnSpacing = 10 };
            var l1 = new Label { Content = M("&Name:"), Target = name, Padding = new Thickness(0, 4, 0, 0) };
            var l2 = new Label { Content = M("&Description:"), Target = description, Padding = new Thickness(0, 4, 0, 0) };
            Grid.SetColumn(name, 1); Grid.SetRow(l2, 1); Grid.SetRow(description, 1); Grid.SetColumn(description, 1);
            grid.Children.Add(l1); grid.Children.Add(name); grid.Children.Add(l2); grid.Children.Add(description);
            var mode = new StackPanel { Spacing = 6 };
            mode.Children.Add(new TextBlock { Text = "Copy mode" });
            mode.Children.Add(clone);
            mode.Children.Add(new TextBlock { Text = Messages.COPY_VM_FAST_CLONE_DESCRIPTION, Opacity = 0.7, Margin = new Thickness(22, 0, 0, 0), TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            mode.Children.Add(copy);
            mode.Children.Add(new TextBlock { Text = vm.is_a_template ? Messages.COPY_TEMPLATE_SELECT_SR : Messages.COPY_VM_SELECT_SR, Opacity = 0.7, Margin = new Thickness(22, 0, 0, 0), TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            mode.Children.Add(new Border { BorderBrush = Avalonia.Media.Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = srList, Margin = new Thickness(22, 0, 0, 0) });
            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(grid);
            panel.Children.Add(new Border { BorderBrush = Avalonia.Media.Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1), Padding = new Thickness(8, 4, 8, 8), Child = mode });
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame(vm.is_a_template ? Messages.COPY_TEMPLATE : Messages.COPY_VM, panel, 520);
            XenAdmin.Actions.AsyncAction result = null;
            ok.Click += (_, _) =>
            {
                result = VmRules.CopyAction(vm, copy.IsChecked == true, srList.SelectedIndex >= 0 ? srs[srList.SelectedIndex] : null, name.Text, description.Text);
                if (result != null) w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>NewPoolDialog: name, description, coordinator and checked members (PoolJoinRules reasons); CreatePoolAction.</summary>
        public static async Task<XenAdmin.Actions.CreatePoolAction> NewPool(Window owner, Func<IEnumerable<XenAdmin.Network.IXenConnection>> connections,
            XenAdmin.Network.IXenConnection preselected, Func<Task> addServer)
        {
            var name = new TextBox();
            var description = new TextBox();
            var coordinator = new ComboBox { HorizontalAlignment = HorizontalAlignment.Stretch };
            var members = new StackPanel { Spacing = 2, Margin = new Thickness(3) };
            var checkedMembers = new HashSet<XenAdmin.Network.IXenConnection>();
            var create = new Button { Content = M("&Create Pool"), IsDefault = true, Width = 95, Height = 24 };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true, Width = 95, Height = 24 };
            var addNew = new Button { Content = M("Add New &Server"), Width = 120, Height = 26 };
            List<PoolRules.Candidate> coordinators = new();
            List<PoolRules.Candidate> candidates = new();
            XenAdmin.Network.IXenConnection chosen = preselected;

            XenAdmin.Network.IXenConnection Coordinator() => coordinator.SelectedIndex >= 0 ? coordinators[coordinator.SelectedIndex].Connection : null;
            List<XenAPI.Host> Supporters() => candidates.Where(x => x.AllowedAsSupporter && checkedMembers.Contains(x.Connection))
                .Select(x => XenAdmin.Core.Helpers.GetCoordinator(x.Connection)).Where(h => h != null).ToList();
            void UpdateButtons()
            {
                var why = PoolRules.Invalid(name.Text, Coordinator(), Supporters().Count);
                create.IsEnabled = why == null;
                ToolTip.SetTip(create, why);
            }
            void Members()
            {
                candidates = PoolRules.Candidates(connections(), Coordinator());
                members.Children.Clear();
                foreach (var x in candidates)
                {
                    var box = new CheckBox
                    {
                        Content = x.Name, IsEnabled = x.AllowedAsSupporter,
                        IsChecked = x.WillBeCoordinator || x.AllowedAsSupporter && checkedMembers.Contains(x.Connection)
                    };
                    box.IsCheckedChanged += (_, _) => { if (box.IsChecked == true) checkedMembers.Add(x.Connection); else checkedMembers.Remove(x.Connection); UpdateButtons(); };
                    var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
                    row.Children.Add(box);
                    row.Children.Add(new TextBlock { Text = x.Description, Opacity = 0.6, VerticalAlignment = VerticalAlignment.Center });
                    members.Children.Add(row);
                }
                UpdateButtons();
            }
            void Coordinators()
            {
                // the proposed coordinator is the one the reasons are computed against
                coordinators = PoolRules.Candidates(connections(), chosen).Where(x => x.CanBeCoordinator).ToList();
                coordinator.ItemsSource = coordinators.Select(x => x.Name).ToList();
                var i = coordinators.FindIndex(x => x.Connection == chosen);
                coordinator.SelectedIndex = i >= 0 ? i : coordinators.Count > 0 ? 0 : -1;
                Members();
            }
            coordinator.SelectionChanged += (_, _) =>
            {
                var c = Coordinator();
                if (c == null || c == chosen) return;
                chosen = c;
                checkedMembers.Remove(c);
                Members();
            };
            name.TextChanged += (_, _) => UpdateButtons();
            Coordinators();

            var fields = new Grid { ColumnDefinitions = new ColumnDefinitions("70,*,Auto"), RowDefinitions = new RowDefinitions("Auto,Auto"), RowSpacing = 5 };
            var ln = new Label { Content = M("&Name:"), Target = name, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
            var ld = new Label { Content = M("&Description:"), Target = description, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
            var optional = new TextBlock { Text = "(optional)", VerticalAlignment = VerticalAlignment.Center, Margin = new Thickness(4, 0, 0, 0) };
            Grid.SetColumn(name, 1); Grid.SetRow(ld, 1); Grid.SetRow(description, 1); Grid.SetColumn(description, 1); Grid.SetRow(optional, 1); Grid.SetColumn(optional, 2);
            fields.Children.Add(ln); fields.Children.Add(name); fields.Children.Add(ld); fields.Children.Add(description); fields.Children.Add(optional);

            var coordRow = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), ColumnSpacing = 4 };
            var lc = new Label { Content = M("C&oordinator:"), Target = coordinator, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
            Grid.SetColumn(coordinator, 1);
            coordRow.Children.Add(lc); coordRow.Children.Add(coordinator);
            var servers = new StackPanel { Spacing = 6 };
            servers.Children.Add(new TextBlock { Text = "Servers" });
            servers.Children.Add(coordRow);
            servers.Children.Add(new TextBlock { Text = "Additional members:", Margin = new Thickness(0, 8, 0, 0) });
            servers.Children.Add(new Border
            {
                BorderBrush = Avalonia.Media.Brush.Parse("#828790"), BorderThickness = new Thickness(1), Background = Avalonia.Media.Brushes.White, Height = 114,
                Child = new ScrollViewer { Content = members }
            });
            servers.Children.Add(addNew);

            var panel = new StackPanel { Spacing = 10 };
            panel.Children.Add(new TextBlock { Text = "To create your pool provide a name and select which servers you would like to be added to the pool.", TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            panel.Children.Add(fields);
            panel.Children.Add(new Border { BorderBrush = Avalonia.Media.Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1), Padding = new Thickness(8, 4, 8, 8), Child = servers });
            panel.Children.Add(Buttons(create, cancel));
            var w = Frame("Create New Pool", panel, 427);
            addNew.Click += async (_, _) => await addServer();
            // the dialog follows servers connecting or disconnecting while it is open (WinForms listens to ConnectionsManager)
            string Connected() => string.Join(",", connections().Where(c => c.IsConnected).Select(c => c.Hostname));
            var seen = Connected();
            var watch = new Avalonia.Threading.DispatcherTimer { Interval = TimeSpan.FromSeconds(1) };
            watch.Tick += (_, _) => { var now = Connected(); if (now != seen) { seen = now; Coordinators(); } };
            watch.Start();
            w.Closed += (_, _) => watch.Stop();
            XenAdmin.Actions.CreatePoolAction result = null;
            create.Click += async (_, _) =>
            {
                var coord = XenAdmin.Core.Helpers.GetCoordinator(Coordinator());
                var supporters = Supporters();
                if (coord == null || PoolRules.Invalid(name.Text, Coordinator(), supporters.Count) != null) return;
                foreach (var warning in PoolRules.Warnings(supporters, coord))
                    if (!await Confirm(w, "XCP-ng Center", warning, "Yes", "No"))
                        return;
                result = PoolRules.Create(coord, supporters, name.Text, description.Text);
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>ConfirmVMDeleteDialog: disks owned by the VM checked, snapshots unchecked; VMDestroyAction.</summary>
        public static async Task<XenAdmin.Actions.AsyncAction> DeleteVm(Window owner, XenAPI.VM vm)
        {
            var c = vm.Connection;
            var disks = c.ResolveAll(vm.VBDs).Where(b => b.type == XenAPI.vbd_type.Disk && c.Resolve(b.VDI) != null).ToList();
            var snaps = c.ResolveAll(vm.snapshots).ToList();
            var items = new List<(CheckBox Box, object Item)>();
            var list = new StackPanel { Spacing = 2 };
            foreach (var b in disks)
            {
                var vdi = c.Resolve(b.VDI);
                var box = new CheckBox { Content = $"{vdi.Name()}   ({vdi.SizeText()})", IsChecked = b.GetIsOwner() };
                items.Add((box, b)); list.Children.Add(box);
            }
            foreach (var sn in snaps)
            {
                var box = new CheckBox { Content = $"{sn.Name()}   (snapshot)", IsChecked = false };
                items.Add((box, sn)); list.Children.Add(box);
            }
            var all = new Button { Content = M("&Select All") };
            var none = new Button { Content = M("&Clear All") };
            all.Click += (_, _) => items.ForEach(i => i.Box.IsChecked = true);
            none.Click += (_, _) => items.ForEach(i => i.Box.IsChecked = false);
            var ok = new Button { Content = M("&Delete"), IsDefault = true };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true };
            var panel = new StackPanel { Spacing = 10 };
            var warn = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 10 };
            warn.Children.Add(Views.Icon("000_WarningAlert_h32bit_32", 32));
            var texts = new StackPanel();
            texts.Children.Add(new TextBlock { Text = string.Format(Messages.CONFIRM_DELETE_VM, vm.Name(), XenAdmin.Core.Helpers.GetName(c)), TextWrapping = Avalonia.Media.TextWrapping.Wrap, MaxWidth = 400 });
            texts.Children.Add(new TextBlock { Text = "Deleted VMs and Templates cannot be recovered.", FontWeight = Avalonia.Media.FontWeight.SemiBold });
            warn.Children.Add(texts);
            panel.Children.Add(warn);
            if (items.Count > 0)
            {
                panel.Children.Add(new TextBlock { Text = "The selected virtual disks and snapshots below will also be deleted:" });
                panel.Children.Add(new Border { BorderBrush = Avalonia.Media.Brush.Parse("#828790"), BorderThickness = new Thickness(1), Background = Avalonia.Media.Brushes.White, Padding = new Thickness(4), Height = 140, Child = new ScrollViewer { Content = list } });
                var sel = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
                sel.Children.Add(all); sel.Children.Add(none);
                panel.Children.Add(sel);
            }
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame("Delete VM", panel, 480);
            XenAdmin.Actions.AsyncAction result = null;
            ok.Click += (_, _) =>
            {
                result = new XenAdmin.Actions.VMActions.VMDestroyAction(vm,
                    items.Where(i => i.Box.IsChecked == true && i.Item is XenAPI.VBD).Select(i => (XenAPI.VBD)i.Item).ToList(),
                    items.Where(i => i.Box.IsChecked == true && i.Item is XenAPI.VM).Select(i => (XenAPI.VM)i.Item).ToList());
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        public static Task<bool> TrustCertificate(Window owner, string hostKey, string fingerprint, string oldFingerprint)
        {
            var msg = oldFingerprint == null
                ? $"The certificate presented by {hostKey} is not trusted.\n\nSHA-256 fingerprint:\n{Pretty(fingerprint)}\n\nOnly continue if this matches the fingerprint of the server's certificate."
                : $"WARNING: the certificate presented by {hostKey} has CHANGED.\n\nPrevious:\n{Pretty(oldFingerprint)}\n\nNew:\n{Pretty(fingerprint)}\n\nThis could indicate an attack. Only continue if the certificate was replaced on purpose.";
            return Confirm(owner, oldFingerprint == null ? "Unknown Certificate" : "Certificate Changed", msg, "Trust and Connect");
        }

        private static string Pretty(string hex)
        {
            var parts = new System.Collections.Generic.List<string>();
            for (int i = 0; i < hex.Length; i += 2)
                parts.Add(hex.Substring(i, 2));
            return string.Join(":", parts);
        }
    }
}
