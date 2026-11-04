using System.Threading.Tasks;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;

namespace XenAdmin.Linux
{
    public static class Dialogs
    {
        private static Window Frame(string title, Control content, double width = 460)
        {
            return new Window
            {
                Title = title,
                Width = width,
                SizeToContent = SizeToContent.Height,
                CanResize = false,
                WindowStartupLocation = WindowStartupLocation.CenterOwner,
                Content = new Border { Padding = new Thickness(14, 12), Child = content }
            };
        }

        private static StackPanel Buttons(params Button[] buttons)
        {
            var p = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, HorizontalAlignment = HorizontalAlignment.Right, Margin = new Thickness(0, 4, 0, 0) };
            p.Children.AddRange(buttons);
            return p;
        }

        public static async Task<bool> Confirm(Window owner, string title, string message, string okText = "OK", string cancelText = "Cancel")
        {
            var ok = new Button { Content = M(okText), IsDefault = true, MinWidth = 75 };
            var cancel = new Button { Content = M(cancelText), IsCancel = true, MinWidth = 75 };
            var panel = new StackPanel();
            panel.Children.Add(new TextBlock { Text = message, TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame(title, panel);
            ok.Click += (_, _) => w.Close(true);
            cancel.Click += (_, _) => w.Close(false);
            return await w.ShowDialog<bool>(owner);
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
                var text = (hostBox.Text ?? "").Trim();
                var port = 443;
                var idx = text.LastIndexOf(':');
                if (idx > 0 && int.TryParse(text[(idx + 1)..], out var p)) { port = p; text = text[..idx]; }
                if (text.Length == 0 || port < 1 || port > 65535) { error.Text = "Enter a valid server name or address."; return; }
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
