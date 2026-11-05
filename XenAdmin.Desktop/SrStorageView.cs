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
        private readonly Grid _rows = new Grid();
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
            var table = new Border { Focusable = true, BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Background = Brushes.White, Child = new ScrollViewer { Content = _rows, HorizontalScrollBarVisibility = Avalonia.Controls.Primitives.ScrollBarVisibility.Auto } };
            var grid = new Grid { Margin = new Thickness(12), RowDefinitions = new RowDefinitions("3*,Auto,*"), MaxWidth = 924, HorizontalAlignment = HorizontalAlignment.Left };
            Grid.SetRow(bar, 1);
            grid.Children.Add(table);
            grid.Children.Add(bar);
            Children.Add(grid);
            table.KeyDown += (_, e) =>
            {
                var disks = StorageRules.SrDisks(_sr);
                var i = disks.FindIndex(v => v.opaque_ref == SelectedRef);
                if (e.Key == Avalonia.Input.Key.Down && i < disks.Count - 1) { SelectedRef = disks[i + 1].opaque_ref; Refresh(); e.Handled = true; }
                else if (e.Key == Avalonia.Input.Key.Up && i > 0) { SelectedRef = disks[i - 1].opaque_ref; Refresh(); e.Handled = true; }
            };
            table.PointerPressed += (_, _) => table.Focus();
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
            var disks = StorageRules.SrDisks(_sr);
            if (disks.All(v => v.opaque_ref != SelectedRef))
                SelectedRef = disks.FirstOrDefault()?.opaque_ref;
            var table = Details.Storage(_sr);
            _rows.Children.Clear();
            _rows.RowDefinitions.Clear();
            _rows.ColumnDefinitions = new ColumnDefinitions(string.Join(",", table.Headers.Select((_, i) => i == table.Fill ? "*" : "Auto")));
            _rows.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            for (int c = 0; c < table.Headers.Length; c++)
                Cell(0, c, table.Headers[c], null);
            // Details.Storage lists the same disks in the same order as StorageRules.SrDisks
            for (int r = 0; r < table.Rows.Count && r < disks.Count; r++)
            {
                _rows.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                for (int c = 0; c < table.Headers.Length; c++)
                    Cell(r + 1, c, table.Rows[r][c], disks[r]);
            }
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

        private void Cell(int row, int col, string text, VDI vdi)
        {
            var header = vdi == null;
            var selected = vdi != null && vdi.opaque_ref == SelectedRef;
            var b = new Border
            {
                Background = selected ? Brush.Parse("#0078D7") : Brushes.White,
                BorderBrush = Brush.Parse("#E3E3E3"),
                BorderThickness = header ? new Thickness(0, 0, 1, 1) : new Thickness(0),
                Padding = new Thickness(6, 3, 14, 3),
                Tag = vdi?.opaque_ref,
                Child = new TextBlock { Text = text ?? "", Foreground = selected ? Brushes.White : Brushes.Black, TextTrimming = TextTrimming.CharacterEllipsis }
            };
            if (vdi != null)
            {
                b.PointerPressed += (_, _) => { SelectedRef = vdi.opaque_ref; Views.PaintSelection(_rows, SelectedRef); UpdateButtons(); };
                b.DoubleTapped += (_, _) => { if (StorageRules.CanEditDisk(_sr, Selected)) _properties(Selected); };
            }
            Grid.SetRow(b, row);
            Grid.SetColumn(b, col);
            _rows.Children.Add(b);
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
