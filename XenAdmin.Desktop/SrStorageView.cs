using System;
using System.Collections.Generic;
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
    /// SR Storage tab (SrStoragePage): the SR's virtual disks and Add / Rescan / Move / Delete / Properties, with the rules
    /// of StorageRules and the XenModel actions (SrRefreshAction, CreateDiskAction, DestroyDiskAction, VMSnapshotDeleteAction).
    /// </summary>
    public class SrStorageView : DockPanel, ILiveView
    {
        // SRs being rescanned (HelpersGUI.BeingScanned looks at the action history, which this app does not keep)
        private static readonly HashSet<string> Scanning = new HashSet<string>();
        private static readonly Dictionary<string, string> SelectedDisk = new Dictionary<string, string>();

        private readonly SR _sr;
        private readonly Action<AsyncAction> _run;
        private readonly Action<IXenObject> _properties;
        private readonly SelectableGrid _grid = new SelectableGrid();
        private readonly Button _add, _rescan, _move, _delete, _props;

        public SrStorageView(SR sr, Action<AsyncAction> run, Action<IXenObject> properties)
        {
            _sr = sr;
            _run = run;
            _properties = properties;
            Background = Brushes.White;
            Button B(string text) => new Button { Content = Dialogs.M(text), MinWidth = 94 };
            _rescan = B("Res&can");
            _add = B("&Add...");
            _move = B("Mov&e...");
            _delete = B("&Delete");
            _props = B("P&roperties");
            _move.IsEnabled = false; // MoveVirtualDiskDialog not ported yet (as in the VM Storage tab)
            _rescan.Click += (_, _) => Rescan();
            _add.Click += async (_, _) => await Add();
            _delete.Click += async (_, _) => await Delete();
            _props.Click += (_, _) => { if (Selected is VDI v) _properties(v); };

            var bar = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(0, 12, 0, 0) };
            foreach (var b in new[] { _rescan, _add, _move, _delete, _props })
            {
                ToolTip.SetShowOnDisabled(b, true);
                bar.Children.Add(b);
            }
            var grid = new Grid { Margin = new Thickness(12), RowDefinitions = new RowDefinitions("3*,Auto,*"), MaxWidth = 924, HorizontalAlignment = HorizontalAlignment.Left };
            Grid.SetRow(bar, 1);
            grid.Children.Add(_grid);
            grid.Children.Add(bar);
            Children.Add(grid);
            _grid.SelectionChanged += () => { SelectedRef = _grid.Selected; UpdateButtons(); };
            _grid.Activated += () => { if (StorageRules.CanEditDisk(_sr, Selected)) _properties(Selected); };
            Refresh();
        }

        private string SelectedRef
        {
            get => SelectedDisk.TryGetValue(_sr.opaque_ref, out var r) ? r : null;
            set => SelectedDisk[_sr.opaque_ref] = value;
        }

        private VDI Selected => StorageRules.SrDisks(_sr).FirstOrDefault(v => v.opaque_ref == SelectedRef);

        private Window Owner => TopLevel.GetTopLevel(this) as Window;

        public void Refresh()
        {
            // Details.Storage lists the same disks in the same order as StorageRules.SrDisks
            _grid.Fill(Details.Storage(_sr), StorageRules.SrDisks(_sr).Select(v => v.opaque_ref).ToList(), SelectedRef);
            SelectedRef = _grid.Selected;
            UpdateButtons();
        }

        private void UpdateButtons()
        {
            var vdi = Selected;
            var rescan = StorageRules.CanRescan(_sr, Scanning.Contains(_sr.opaque_ref));
            _rescan.IsEnabled = rescan.Enabled;
            ToolTip.SetTip(_rescan, rescan.Tooltip);
            _add.IsEnabled = StorageRules.CanAddDisk(_sr);
            _delete.IsEnabled = vdi != null && StorageRules.CanDelete(vdi, allowMultipleVbds: true);
            _props.IsEnabled = StorageRules.CanEditDisk(_sr, vdi);
        }

        // ---------- commands (SrStoragePage handlers) ----------

        private void Rescan()
        {
            var srRef = _sr.opaque_ref;
            var action = new SrRefreshAction(_sr);
            Scanning.Add(srRef);
            action.Completed += _ => Dispatcher.UIThread.Post(() => Scanning.Remove(srRef));
            _run(action);
            Refresh();
        }

        private async Task Add()
        {
            var action = await Dialogs.NewDisk(Owner, _sr);
            if (action != null) _run(action);
        }

        private async Task Delete()
        {
            var vdi = Selected;
            if (vdi == null) return;
            var (title, text) = StorageRules.DeleteConfirmation(vdi);
            if (!await Dialogs.Confirm(Owner, title, text, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION))
                return;
            var actions = StorageRules.DeleteActions(vdi);
            if (actions.Count == 1)
                _run(actions[0]);
            else if (actions.Count > 1)
                _run(new MultipleAction(_sr.Connection, Messages.ACTION_DELETING_MULTIPLE_STORAGE_ITEMS_TITLE,
                    Messages.ACTION_DELETING_MULTIPLE_STORAGE_ITEMS_STATUS, Messages.COMPLETED, actions));
        }
    }
}
