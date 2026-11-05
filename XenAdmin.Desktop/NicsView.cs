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
    /// Host NICs tab (NICPage): the physical interfaces and Create Bond (BondProperties) / Delete Bond
    /// (DestroyBondCommand) / Rescan (RescanPIFsAction).
    /// </summary>
    public class NicsView : DockPanel, ILiveView
    {
        private readonly Host _host;
        private readonly SelectableGrid _grid = new SelectableGrid();
        private readonly Button _rescan, _create, _delete;

        public NicsView(Host host, Action<AsyncAction> run)
        {
            _host = host;
            Background = Brushes.White;
            Button B(string text, bool enabled) => new Button { Content = Dialogs.M(text), MinWidth = 94, IsEnabled = enabled };
            _create = B("&Create Bond...", true);
            _delete = B("&Delete Bond", false);
            _rescan = B("&Rescan", true);
            _rescan.Click += (_, _) => run(new RescanPIFsAction(_host));
            _create.Click += async (_, _) => { if (await Dialogs.CreateBond(Owner, _host) is AsyncAction a) run(a); };
            _delete.Click += async (_, _) => { if (BondNetwork() is XenAPI.Network n && await Dialogs.DeleteBond(Owner, n) is AsyncAction a) run(a); };
            _grid.SelectionChanged += UpdateButtons;

            var bar = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(0, 12, 0, 0) };
            bar.Children.Add(_create);
            bar.Children.Add(_delete);
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
            UpdateButtons();
        }

        private Window Owner => TopLevel.GetTopLevel(this) as Window;

        /// <summary>NICPage.DeleteBondButton_Click: the network of the selected bond interface.</summary>
        private XenAPI.Network BondNetwork()
        {
            var pif = NetworkRules.HostNics(_host).FirstOrDefault(p => p.opaque_ref == _grid.Selected);
            return pif != null && pif.IsBondNIC() ? pif.Connection.Resolve(pif.network) : null;
        }

        private void UpdateButtons()
        {
            _rescan.IsEnabled = !_host.Locked;
            _create.IsEnabled = !_host.Locked;
            _delete.IsEnabled = BondRules.CanDelete(BondNetwork());
        }
    }
}
