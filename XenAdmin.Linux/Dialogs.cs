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

        public static async Task<bool> Confirm(Window owner, string title, string message, string okText = "OK")
        {
            var ok = new Button { Content = okText, IsDefault = true };
            var cancel = new Button { Content = "Cancelar", IsCancel = true };
            var panel = new StackPanel();
            panel.Children.Add(new TextBlock { Text = message, TextWrapping = Avalonia.Media.TextWrapping.Wrap });
            panel.Children.Add(Buttons(cancel, ok));
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

        public static async Task<ServerInput> AddServer(Window owner, string host = "", string user = "root")
        {
            var hostBox = new TextBox { Text = host, Watermark = "servidor ou servidor:porta" };
            var userBox = new TextBox { Text = user };
            var passBox = new TextBox { PasswordChar = '•' };
            var error = new TextBlock { Foreground = Avalonia.Media.Brushes.IndianRed, TextWrapping = Avalonia.Media.TextWrapping.Wrap };
            var ok = new Button { Content = "Conectar", IsDefault = true };
            var cancel = new Button { Content = "Cancelar", IsCancel = true };
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowDefinitions = new RowDefinitions("Auto,Auto,Auto,Auto,Auto"), RowSpacing = 8, ColumnSpacing = 8 };
            void Row(int r, string label, Control c)
            {
                var l = new Label { Content = label, Target = c, VerticalAlignment = VerticalAlignment.Center };
                Grid.SetRow(l, r); Grid.SetRow(c, r); Grid.SetColumn(c, 1);
                grid.Children.Add(l); grid.Children.Add(c);
            }
            Row(0, "_Servidor:", hostBox);
            Row(1, "_Usuário:", userBox);
            Row(2, "_Senha:", passBox);
            var note = new TextBlock { Text = "A senha não é gravada em disco.", Opacity = 0.7 };
            Grid.SetRow(note, 3); Grid.SetColumnSpan(note, 2); grid.Children.Add(note);
            Grid.SetRow(error, 4); Grid.SetColumnSpan(error, 2); grid.Children.Add(error);
            var panel = new StackPanel();
            panel.Children.Add(grid);
            panel.Children.Add(Buttons(cancel, ok));
            var w = Frame("Adicionar servidor", panel);
            ServerInput result = null;
            ok.Click += (_, _) =>
            {
                var text = (hostBox.Text ?? "").Trim();
                var port = 443;
                var idx = text.LastIndexOf(':');
                if (idx > 0 && int.TryParse(text[(idx + 1)..], out var p)) { port = p; text = text[..idx]; }
                if (text.Length == 0 || port < 1 || port > 65535) { error.Text = "Informe um servidor válido."; return; }
                if (string.IsNullOrWhiteSpace(userBox.Text)) { error.Text = "Informe o usuário."; return; }
                result = new ServerInput { Host = text, Port = port, User = userBox.Text.Trim(), Password = passBox.Text ?? "" };
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            w.Opened += (_, _) => (string.IsNullOrEmpty(host) ? hostBox : passBox).Focus();
            await w.ShowDialog(owner);
            return result;
        }

        public static Task<bool> TrustCertificate(Window owner, string hostKey, string fingerprint, string oldFingerprint)
        {
            var msg = oldFingerprint == null
                ? $"O certificado de {hostKey} não é confiável pelo sistema.\n\nSHA-256:\n{Pretty(fingerprint)}\n\nConfie somente se a impressão digital conferir com a do servidor (xe host-param-get param-name=certificates… ou console do host)."
                : $"ATENÇÃO: o certificado de {hostKey} MUDOU.\n\nAnterior:\n{Pretty(oldFingerprint)}\n\nNovo:\n{Pretty(fingerprint)}\n\nIsso pode indicar um ataque. Aceite apenas se o certificado foi trocado de propósito.";
            return Confirm(owner, oldFingerprint == null ? "Certificado não confiável" : "Certificado alterado", msg, "Confiar e conectar");
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
