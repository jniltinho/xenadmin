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
                Content = new Border { Padding = new Thickness(16), Child = content }
            };
        }

        private static StackPanel Buttons(params Button[] buttons)
        {
            var p = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, HorizontalAlignment = HorizontalAlignment.Right, Margin = new Thickness(0, 12, 0, 0) };
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
            var hostBox = new TextBox { Text = host, IsReadOnly = connectMode };
            var userBox = new TextBox { Text = user };
            var passBox = new TextBox { PasswordChar = '•' };
            var error = new TextBlock { Foreground = Avalonia.Media.Brushes.Red, TextWrapping = Avalonia.Media.TextWrapping.Wrap };
            var ok = new Button { Content = M(connectMode ? Messages.CONNECT : Messages.ADD), IsDefault = true, MinWidth = 75 };
            var cancel = new Button { Content = M(Messages.CANCEL), IsCancel = true, MinWidth = 75 };
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowDefinitions = new RowDefinitions("Auto,Auto,Auto,Auto"), RowSpacing = 8, ColumnSpacing = 8 };
            void Row(int r, string label, Control c)
            {
                var l = new Label { Content = M(label), Target = c, VerticalAlignment = VerticalAlignment.Center };
                Grid.SetRow(l, r); Grid.SetRow(c, r); Grid.SetColumn(c, 1);
                grid.Children.Add(l); grid.Children.Add(c);
            }
            Row(0, "&Server:", hostBox);
            Row(1, "&Username:", userBox);
            Row(2, "&Password:", passBox);
            Grid.SetRow(error, 3); Grid.SetColumnSpan(error, 2); grid.Children.Add(error);
            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(new TextBlock
            {
                Text = connectMode ? Messages.CONNECT_TO_SERVER_BLURB
                    : "Enter the host name or IP address of the server you want to add and your user login credentials for that server.",
                TextWrapping = Avalonia.Media.TextWrapping.Wrap
            });
            panel.Children.Add(grid);
            panel.Children.Add(Buttons(ok, cancel));
            var w = Frame(connectMode ? Messages.CONNECT_TO_SERVER : "Add New Server", panel);
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
