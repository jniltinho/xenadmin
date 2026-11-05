using System;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using XenAdmin.Actions;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// Host NICs tab (NICPage): the physical interfaces and Create Bond / Delete Bond / Rescan (RescanPIFsAction).
    /// Bond creation and deletion are listed disabled until the bond dialogs are ported.
    /// </summary>
    public class NicsView : DockPanel, ILiveView
    {
        private readonly Host _host;
        private readonly SelectableGrid _grid = new SelectableGrid();
        private readonly Button _rescan;

        public NicsView(Host host, Action<AsyncAction> run)
        {
            _host = host;
            Background = Brushes.White;
            Button B(string text, bool enabled) => new Button { Content = Dialogs.M(text), MinWidth = 94, IsEnabled = enabled };
            var create = B("&Create Bond...", false);
            var delete = B("&Delete Bond", false);
            _rescan = B("&Rescan", true);
            foreach (var b in new[] { create, delete })
            {
                ToolTip.SetShowOnDisabled(b, true);
                ToolTip.SetTip(b, "Not available in this version yet.");
            }
            _rescan.Click += (_, _) => run(new RescanPIFsAction(_host));

            var bar = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(0, 12, 0, 0) };
            bar.Children.Add(create);
            bar.Children.Add(delete);
            bar.Children.Add(_rescan);
            var grid = new Grid { Margin = new Thickness(12), RowDefinitions = new RowDefinitions("3*,Auto,*"), MaxWidth = 1200, HorizontalAlignment = HorizontalAlignment.Left };
            Grid.SetRow(bar, 1);
            grid.Children.Add(_grid);
            grid.Children.Add(bar);
            Children.Add(grid);
            Refresh();
        }

        public void Refresh()
        {
            _grid.Fill(Details.Nics(_host), NetworkRules.HostNics(_host).Select(p => p.opaque_ref).ToList(), _grid.Selected);
            _rescan.IsEnabled = !_host.Locked;
        }
    }
}
