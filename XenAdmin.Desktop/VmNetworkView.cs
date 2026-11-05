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
        private readonly Grid _rows = new Grid();
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
            var table = new Border { Focusable = true, BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Background = Brushes.White, Child = new ScrollViewer { Content = _rows, HorizontalScrollBarVisibility = Avalonia.Controls.Primitives.ScrollBarVisibility.Auto } };
            var grid = new Grid { Margin = new Thickness(12), RowDefinitions = new RowDefinitions("3*,Auto,*"), MaxWidth = 924, HorizontalAlignment = HorizontalAlignment.Left };
            Grid.SetRow(bar, 1);
            grid.Children.Add(table);
            grid.Children.Add(bar);
            Children.Add(grid);
            // keyboard selection, as in the DataGridView: Up/Down move, Enter opens Properties
            table.KeyDown += async (_, e) =>
            {
                var vifs = Vifs();
                var i = Array.FindIndex(vifs, v => v.device == _selected);
                if (e.Key == Avalonia.Input.Key.Down && i < vifs.Length - 1) { _selected = vifs[i + 1].device; Refresh(); e.Handled = true; }
                else if (e.Key == Avalonia.Input.Key.Up && i > 0) { _selected = vifs[i - 1].device; Refresh(); e.Handled = true; }
                else if (e.Key == Avalonia.Input.Key.Enter && NetworkRules.CanEdit(_vm, Selected)) { e.Handled = true; await Edit(); }
            };
            table.PointerPressed += (_, _) => table.Focus();
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
            var table = Details.Networking(_vm);
            _rows.Children.Clear();
            _rows.RowDefinitions.Clear();
            _rows.ColumnDefinitions = new ColumnDefinitions(string.Join(",", table.Headers.Select((_, i) => i == 4 ? "*" : "Auto")));
            _rows.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            for (int c = 0; c < table.Headers.Length; c++)
                Cell(0, c, table.Headers[c], null);
            // Details.Networking lists the VIFs in the same device order as Vifs()
            for (int r = 0; r < table.Rows.Count && r < vifs.Length; r++)
            {
                _rows.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                for (int c = 0; c < table.Headers.Length; c++)
                    Cell(r + 1, c, table.Rows[r][c], vifs[r]);
            }
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

        private void Cell(int row, int col, string text, VIF vif)
        {
            var header = vif == null;
            var selected = vif != null && vif.device == _selected;
            var b = new Border
            {
                Background = selected ? Brush.Parse("#0078D7") : Brushes.White,
                BorderBrush = Brush.Parse("#E3E3E3"),
                BorderThickness = header ? new Thickness(0, 0, 1, 1) : new Thickness(0),
                Padding = new Thickness(6, 3, 14, 3),
                Tag = vif?.device,
                Child = new TextBlock { Text = text ?? "", Foreground = selected ? Brushes.White : Brushes.Black, TextTrimming = TextTrimming.CharacterEllipsis }
            };
            if (vif != null)
            {
                b.PointerPressed += (_, _) => { _selected = vif.device; EditInFlight.Remove(_vm.opaque_ref); Views.PaintSelection(_rows, _selected); UpdateButtons(); };
                b.DoubleTapped += async (_, _) => { if (NetworkRules.CanEdit(_vm, Selected)) await Edit(); };
            }
            Grid.SetRow(b, row);
            Grid.SetColumn(b, col);
            _rows.Children.Add(b);
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
