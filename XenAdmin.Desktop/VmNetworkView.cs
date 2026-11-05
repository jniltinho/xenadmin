using System;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using Avalonia.Threading;
using XenAdmin.Actions;
using XenAPI;
using Task = System.Threading.Tasks.Task;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// VM Networking tab (NetworkList for a VM): the interface list and Add Interface / Properties / Remove Interface /
    /// Activate-Deactivate, with the rules of NetworkRules and the XenModel VIF actions.
    /// </summary>
    public class VmNetworkView : DockPanel, ILiveView
    {
        private readonly VM _vm;
        private readonly Action<AsyncAction> _run;
        private readonly SelectableGrid _grid = new SelectableGrid();
        private readonly Button _add, _edit, _remove, _toggle;
        // the selected device per VM: edits recreate the VIF (new opaque_ref) and the tab is rebuilt after each action
        private static readonly System.Collections.Generic.Dictionary<string, string> SelectedDevice = new System.Collections.Generic.Dictionary<string, string>();

        // device being edited per VM: UpdateVIFAction destroys and recreates the VIF, so the device is missing until the
        // replacement reaches the cache (WinForms sets InBuildList during the edit); cleared then, or when the edit fails
        private static readonly System.Collections.Generic.Dictionary<string, (string Device, string OldRef)> EditInFlight = new System.Collections.Generic.Dictionary<string, (string, string)>();

        private string _selected // VIF device
        {
            get => SelectedDevice.TryGetValue(_vm.opaque_ref, out var d) ? d : null;
            set => SelectedDevice[_vm.opaque_ref] = value;
        }

        public VmNetworkView(VM vm, Action<AsyncAction> run)
        {
            _vm = vm;
            _run = run;
            Background = Brushes.White;
            Button B(string text) => new Button { Content = Dialogs.M(text), MinWidth = 94 };
            _add = B(Messages.VM_NETWORK_TAB_ADD_BUTTON_LABEL);
            _edit = B(Messages.VM_NETWORK_TAB_EDIT_BUTTON_LABEL);
            _remove = B(Messages.VM_NETWORK_TAB_REMOVE_BUTTON_LABEL);
            _toggle = B(Messages.VM_NETWORK_TAB_DEACTIVATE_BUTTON_LABEL);
            _add.Click += async (_, _) => await Add();
            _edit.Click += async (_, _) => await Edit();
            _remove.Click += async (_, _) => await Remove();
            _toggle.Click += (_, _) => Toggle();

            // WinForms order: Add, Properties, Remove, Activate/Deactivate (each disabled button keeps its tooltip)
            var bar = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(0, 12, 0, 0) };
            foreach (var b in new[] { _add, _edit, _remove, _toggle })
            {
                ToolTip.SetShowOnDisabled(b, true);
                bar.Children.Add(b);
            }
            var grid = new Grid { Margin = new Thickness(12), RowDefinitions = new RowDefinitions("3*,Auto,*"), MaxWidth = 924, HorizontalAlignment = HorizontalAlignment.Left };
            Grid.SetRow(bar, 1);
            grid.Children.Add(_grid);
            grid.Children.Add(bar);
            Children.Add(grid);
            _grid.SelectionChanged += () => { _selected = _grid.Selected; EditInFlight.Remove(_vm.opaque_ref); UpdateButtons(); };
            _grid.Activated += async () => { if (NetworkRules.CanEdit(_vm, Selected)) await Edit(); };
            Refresh();
        }

        private VIF[] Vifs() => _vm.Connection.ResolveAll(_vm.VIFs).OrderBy(v => v.device).ToArray(); // same order as Details.Networking

        private VIF Selected => Vifs().FirstOrDefault(v => v.device == _selected);

        private Window Owner => TopLevel.GetTopLevel(this) as Window;

        public void Refresh()
        {
            var vifs = Vifs();
            string pending = null;
            if (EditInFlight.TryGetValue(_vm.opaque_ref, out var edit))
            {
                if (NetworkRules.EditReplaced(vifs.Select(v => (v.device, v.opaque_ref)), edit.Device, edit.OldRef))
                    EditInFlight.Remove(_vm.opaque_ref);
                else
                    pending = edit.Device;
            }
            _selected = NetworkRules.Selection(vifs.Select(v => v.device).ToList(), _selected, pending);
            // Details.Networking lists the VIFs in the same device order as Vifs()
            _grid.Fill(Details.Networking(_vm), vifs.Select(v => v.device).ToList(), _selected, keepMissing: true);
            UpdateButtons();
        }

        private void UpdateButtons()
        {
            var vif = Selected;
            _add.IsEnabled = NetworkRules.CanAdd(_vm);
            _edit.IsEnabled = NetworkRules.CanEdit(_vm, vif);
            _remove.IsEnabled = NetworkRules.CanRemove(_vm, vif);
            _toggle.IsEnabled = NetworkRules.CanToggle(_vm, vif);
            _toggle.Content = Dialogs.M(NetworkRules.ToggleLabel(vif));
            var tips = NetworkRules.Tooltips(_vm, vif);
            ToolTip.SetTip(_edit, tips.Edit);
            ToolTip.SetTip(_remove, tips.Remove);
            ToolTip.SetTip(_toggle, tips.Toggle);
        }

        // ---------- commands (NetworkList button handlers) ----------

        private async Task Add()
        {
            if (NetworkRules.AtMaxInterfaces(_vm))
            {
                await Dialogs.Info(Owner, FriendlyErrorNames.VIFS_MAX_ALLOWED_TITLE, FriendlyErrorNames.VIFS_MAX_ALLOWED);
                return;
            }
            var owner = Owner; // the view may be gone (tab or VM switched) when the action completes
            var vif = await Dialogs.Vif(owner, _vm, null, VIF.GetDeviceId(_vm));
            if (vif == null) return;
            _selected = vif.device;
            var action = new CreateVIFAction(_vm, vif);
            action.Completed += _ => { if (action.RebootRequired) Dispatcher.UIThread.Post(() => ShowHotPlugError(owner)); };
            _run(action);
        }

        private async Task Edit()
        {
            var existing = Selected;
            if (existing == null || !int.TryParse(existing.device, out var device)) return;
            var owner = Owner;
            var vif = await Dialogs.Vif(owner, _vm, existing, device);
            if (vif == null) return;
            var action = new UpdateVIFAction(_vm, existing, vif);
            var vmRef = _vm.opaque_ref;
            EditInFlight[vmRef] = (existing.device, existing.opaque_ref);
            action.Completed += _ => Dispatcher.UIThread.Post(() =>
            {
                if (!action.Succeeded) EditInFlight.Remove(vmRef);
                if (action.RebootRequired) ShowHotPlugError(owner);
            });
            _run(action);
        }

        private async Task Remove()
        {
            var vif = Selected;
            if (vif == null) return;
            if (await Dialogs.Confirm(Owner, Messages.MESSAGEBOX_VIF_DELETE_TITLE, Messages.MESSAGEBOX_VIF_DELETE, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION))
                _run(new DeleteVIFAction(vif));
        }

        private void Toggle()
        {
            var vif = Selected;
            if (vif == null) return;
            _run(vif.currently_attached ? (AsyncAction)new UnplugVIFAction(vif) : new PlugVIFAction(vif));
        }

        private static void ShowHotPlugError(Window owner)
        {
            if (owner != null)
                _ = Dialogs.Info(owner, Messages.VIF_HOTPLUG_FAILED_TITLE, Messages.VIF_HOTPLUG_FAILED_MESSAGE);
        }
    }
}
