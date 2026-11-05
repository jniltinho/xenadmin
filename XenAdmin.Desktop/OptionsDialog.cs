using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using XenAPI;
using Task = System.Threading.Tasks.Task;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// Tools > Options (XenAdmin/Dialogs/OptionsDialog): the pages on the left (Security, Connection, Confirmations,
    /// Display), the chosen page with its title and description on the right, OK applies all of them.
    /// </summary>
    public static class OptionsDialog
    {
        private sealed class Page
        {
            public string Title, SubText, Icon;
            public Control Body;
        }

        /// <summary>Shows the dialog; true when OK changed <paramref name="options"/> (the caller saves them).</summary>
        public static async System.Threading.Tasks.Task<bool> Show(Window owner, AppOptions options)
        {
            static CheckBox Check(string text, bool value) => new CheckBox { Content = Dialogs.M(text), IsChecked = value };
            static TextBlock Wrap(string text) => new TextBlock { Text = text, TextWrapping = TextWrapping.Wrap };
            static TextBlock Group(string text) => new TextBlock { Text = text, FontWeight = FontWeight.SemiBold, Margin = new Thickness(0, 6, 0, 0) };
            static StackPanel Stack(params Control[] c) { var p = new StackPanel { Spacing = 8 }; p.Children.AddRange(c); return p; }

            // SecurityOptionsPage
            var newCert = Check("Warn me when a &new certificate is found", options.WarnUnrecognizedCertificate);
            var changedCert = Check("Warn me when a certificate ch&anges", options.WarnChangedCertificate);
            var remind = Check("&Remind me to change the server root password when I rotate the pool secret", options.RotateSecretReminder);
            var security = Stack(
                Wrap("You can configure XCP-ng Center to display a warning whenever a new or modified certificate is found on a managed server."),
                newCert, changedCert,
                Wrap("You can configure XCP-ng Center to display a reminder for changing the server root password whenever you rotate the pool secret."),
                remind);

            // ConnectionOptionsPage
            var direct = new RadioButton { Content = Dialogs.M("&Don't use a proxy server"), GroupName = "proxy", IsChecked = options.ProxyStyle == HTTPHelper.ProxyStyle.DirectConnection };
            var system = new RadioButton { Content = Dialogs.M("Use your &operating system's proxy server settings"), GroupName = "proxy", IsChecked = options.ProxyStyle == HTTPHelper.ProxyStyle.SystemProxy };
            var specified = new RadioButton { Content = Dialogs.M("Use &this proxy server:"), GroupName = "proxy", IsChecked = options.ProxyStyle == HTTPHelper.ProxyStyle.SpecifiedProxy };
            var address = new TextBox { Text = options.ProxyAddress, Width = 220 };
            var port = new TextBox { Text = options.ProxyPort.ToString(), Width = 70 };
            var bypass = Check("&Bypass proxy server for all server connections", options.BypassProxyForServers);
            var proxyRow = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6, Margin = new Thickness(22, 0, 0, 0) };
            proxyRow.Children.AddRange(new Control[]
            {
                new Label { Content = Dialogs.M("&Address:"), Target = address, VerticalAlignment = VerticalAlignment.Center }, address,
                new Label { Content = Dialogs.M("&Port:"), Target = port, VerticalAlignment = VerticalAlignment.Center }, port
            });
            bypass.Margin = new Thickness(22, 0, 0, 0);
            void ProxyEnabled() => proxyRow.IsEnabled = bypass.IsEnabled = specified.IsChecked == true;
            specified.IsCheckedChanged += (_, _) => ProxyEnabled();
            ProxyEnabled();
            var timeout = new NumericUpDown { Minimum = AppOptions.MinTimeout, Maximum = AppOptions.MaxTimeout, Value = options.ConnectionTimeoutSeconds, FormatString = "0", Increment = 1, Width = 130 };
            var timeoutRow = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
            timeoutRow.Children.AddRange(new Control[]
            {
                new Label { Content = Dialogs.M("Co&nnection timeout:"), Target = timeout, VerticalAlignment = VerticalAlignment.Center }, timeout,
                new TextBlock { Text = "seconds", VerticalAlignment = VerticalAlignment.Center }
            });
            var connection = Stack(
                Wrap("These settings allow you to configure a proxy server or to customize your server connections for high latency environments."),
                Group("Proxy server"), direct, system, specified, proxyRow, bypass,
                Group("Connection timeout"), Wrap("Specify how long to wait when the connection to a managed server is disrupted."), timeoutRow);

            // ConfirmationOptionsPage
            var dismissAlerts = Check("Don't show dismiss &alerts confirmation dialog", options.DoNotConfirmDismissAlerts);
            var confirmations = Stack(
                Wrap("You can configure XCP-ng Center to display a warning whenever you dismiss one or more alerts on the 'Notifications' view."),
                dismissAlerts);

            // DisplayOptionsPage
            var area = new RadioButton { Content = Dialogs.M("&Area graph"), GroupName = "graph", IsChecked = options.FillAreaUnderGraphs };
            var line = new RadioButton { Content = Dialogs.M("&Line graph"), GroupName = "graph", IsChecked = !options.FillAreaUnderGraphs };
            var rememberTab = Check("&Remember the last selected Tab for a resource", options.RememberLastSelectedTab);
            var display = Stack(
                Group("Graph type"), Wrap("Configure the appearance of the performance graphs displayed on the Performance tab."), area, line,
                Group("Tab options"), rememberTab);

            var pages = new List<Page>
            {
                new Page { Title = Messages.SECURITY, SubText = Messages.SECURITY_DESC, Icon = "padlock", Body = security },
                new Page { Title = Messages.CONNECTION, SubText = Messages.CONNECTION_DESC, Icon = "000_Network_h32bit_16", Body = connection },
                new Page { Title = Messages.CONFIRMATIONS, SubText = Messages.CONFIRMATIONS_DETAIL, Icon = "075_TickRound_h32bit_16", Body = confirmations },
                new Page { Title = Messages.DISPLAY, SubText = Messages.DISPLAY_DETAILS, Icon = "001_PerformanceGraph_h32bit_16", Body = display },
            };

            var list = new ListBox
            {
                Width = 170,
                ItemsSource = pages.Select(p =>
                {
                    var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
                    row.Children.Add(Views.Icon(p.Icon, 16));
                    row.Children.Add(new TextBlock { Text = p.Title, VerticalAlignment = VerticalAlignment.Center });
                    return row;
                }).ToList()
            };
            var heading = new TextBlock { FontWeight = FontWeight.SemiBold, FontSize = 14 };
            var sub = new TextBlock { Opacity = 0.75 };
            var body = new ContentControl();
            list.SelectionChanged += (_, _) =>
            {
                if (list.SelectedIndex < 0) return;
                var p = pages[list.SelectedIndex];
                heading.Text = p.Title;
                sub.Text = p.SubText;
                body.Content = p.Body;
            };

            var error = new TextBlock { Foreground = Brushes.Red, TextWrapping = TextWrapping.Wrap };
            var ok = new Button { Content = "OK", IsDefault = true, MinWidth = 75 };
            var cancel = new Button { Content = Dialogs.M(Messages.CANCEL), IsCancel = true, MinWidth = 75 };

            var right = new DockPanel { Margin = new Thickness(12, 0, 0, 0) };
            var top = Stack(heading, sub, new Border { Height = 1, Background = Brush.Parse("#DCDCDC") });
            DockPanel.SetDock(top, Dock.Top);
            right.Children.Add(top);
            right.Children.Add(body);
            var main = new DockPanel { Height = 400 };
            DockPanel.SetDock(list, Dock.Left);
            main.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = list, [DockPanel.DockProperty] = Dock.Left });
            main.Children.Add(right);
            var panel = new StackPanel { Spacing = 8 };
            panel.Children.Add(main);
            panel.Children.Add(error);
            panel.Children.Add(Dialogs.Buttons(ok, cancel));
            var w = Dialogs.Frame("Options", panel, 760);
            list.SelectedIndex = 0;

            var saved = false;
            ok.Click += (_, _) =>
            {
                var style = specified.IsChecked == true ? HTTPHelper.ProxyStyle.SpecifiedProxy
                    : system.IsChecked == true ? HTTPHelper.ProxyStyle.SystemProxy : HTTPHelper.ProxyStyle.DirectConnection;
                if (AppOptions.ProxyError(style, address.Text, port.Text) is string err)
                {
                    list.SelectedIndex = 1;
                    error.Text = err;
                    return;
                }
                options.WarnUnrecognizedCertificate = newCert.IsChecked == true;
                options.WarnChangedCertificate = changedCert.IsChecked == true;
                options.RotateSecretReminder = remind.IsChecked == true;
                options.ProxyStyle = style;
                if (style == HTTPHelper.ProxyStyle.SpecifiedProxy)
                {
                    options.ProxyAddress = address.Text.Trim();
                    options.ProxyPort = int.Parse(port.Text.Trim());
                }
                options.BypassProxyForServers = bypass.IsChecked == true;
                options.ConnectionTimeoutSeconds = (int)(timeout.Value ?? 20);
                options.DoNotConfirmDismissAlerts = dismissAlerts.IsChecked == true;
                options.FillAreaUnderGraphs = area.IsChecked == true;
                options.RememberLastSelectedTab = rememberTab.IsChecked == true;
                saved = true;
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return saved;
        }
    }
}
