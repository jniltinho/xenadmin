using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using XenAdmin.Actions;
using XenAPI;
using Task = System.Threading.Tasks.Task;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// Host and pool Networking tab (NetworkList for a host or pool): the networks and Add Network / Properties /
    /// Remove Network, with the rules of NetworkRules, the New Network wizard and NetworkAction.
    /// </summary>
    public class HostNetworkView : DockPanel, ILiveView
    {
        // the selected network per host or pool: the tab is rebuilt after each action
        private static readonly Dictionary<string, string> SelectedNetwork = new Dictionary<string, string>();

        private readonly IXenObject _owner;
        private readonly Action<AsyncAction> _run;
        private readonly Action<IXenObject> _properties;
        private readonly SelectableGrid _grid = new SelectableGrid();
        private readonly Button _add, _edit, _remove;

        public HostNetworkView(IXenObject hostOrPool, Action<AsyncAction> run, Action<IXenObject> properties)
        {
            _owner = hostOrPool;
            _run = run;
            _properties = properties;
            Background = Brushes.White;
            Button B(string text) => new Button { Content = Dialogs.M(text), MinWidth = 94 };
            _add = B(Messages.HOST_NETWORK_TAB_ADD_BUTTON_LABEL);
            _edit = B(Messages.HOST_NETWORK_TAB_EDIT_BUTTON_LABEL);
            _remove = B(Messages.HOST_NETWORK_TAB_REMOVE_BUTTON_LABEL);
            ToolTip.SetShowOnDisabled(_remove, true);
            _add.Click += async (_, _) => await Add();
            _edit.Click += (_, _) => Edit();
            _remove.Click += async (_, _) => await Remove();
            _grid.SelectionChanged += () => { SelectedNetwork[_owner.opaque_ref] = _grid.Selected; UpdateButtons(); };
            _grid.Activated += Edit;

            var bar = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(0, 12, 0, 0) };
            bar.Children.Add(_add);
            bar.Children.Add(_edit);
            bar.Children.Add(_remove);
            var grid = new Grid { Margin = new Thickness(12), RowDefinitions = new RowDefinitions("3*,Auto,*"), MaxWidth = 1100, HorizontalAlignment = HorizontalAlignment.Left };
            Grid.SetRow(bar, 1);
            grid.Children.Add(_grid);
            grid.Children.Add(bar);
            Children.Add(grid);
            Refresh();
        }

        private XenAPI.Network Selected => NetworkRules.HostNetworks(_owner.Connection).FirstOrDefault(n => n.opaque_ref == _grid.Selected);

        private Window Owner => TopLevel.GetTopLevel(this) as Window;

        public void Refresh()
        {
            SelectedNetwork.TryGetValue(_owner.opaque_ref, out var selected);
            // Details.Networking lists the same networks in the same order as NetworkRules.HostNetworks
            _grid.Fill(Details.Networking(_owner), NetworkRules.HostNetworks(_owner.Connection).Select(n => n.opaque_ref).ToList(), selected);
            SelectedNetwork[_owner.opaque_ref] = _grid.Selected;
            UpdateButtons();
        }

        private void UpdateButtons()
        {
            var n = Selected;
            _add.IsEnabled = !_owner.Locked && NetworkRules.WizardHost(_owner) != null;
            _edit.IsEnabled = NetworkRules.CanEditNetwork(_owner, n);
            var remove = NetworkRules.CanRemoveNetwork(_owner, n);
            _remove.IsEnabled = remove.Enabled;
            ToolTip.SetTip(_remove, remove.Tooltip);
        }

        private async Task Add() => await new NewNetworkWizard(_owner, _run).ShowDialog(Owner);

        private void Edit()
        {
            if (NetworkRules.CanEditNetwork(_owner, Selected)) _properties(Selected);
        }

        private async Task Remove()
        {
            var n = Selected;
            if (!NetworkRules.CanRemoveNetwork(_owner, n).Enabled) return;
            if (await Dialogs.Confirm(Owner, Messages.MESSAGEBOX_NETWORK_DELETE_TITLE, Messages.MESSAGEBOX_NETWORK_DELETE, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION))
                _run(new NetworkAction(n.Connection, n, false));
        }
    }
}
