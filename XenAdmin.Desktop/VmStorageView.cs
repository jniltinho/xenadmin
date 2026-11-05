using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAPI;
using Task = System.Threading.Tasks.Task;

namespace XenAdmin.Desktop
{
    /// <summary>VMStoragePage: DVD drive, virtual disk list and Add/Attach/Activate/Detach/Delete/Properties.</summary>
    public class VmStorageView : DockPanel, ILiveView
    {
        private readonly VM _vm;
        private readonly Action<AsyncAction> _run;
        private readonly Action<IXenObject> _properties;
        private readonly ContentControl _dvd = new ContentControl();
        private readonly Grid _rows = new Grid();
        private readonly Button _add, _attach, _activate, _move, _detach, _delete, _props;
        private string _selected; // VBD opaque_ref
        private bool _updatingDvd;

        public VmStorageView(VM vm, Action<AsyncAction> run, Action<IXenObject> properties)
        {
            _vm = vm;
            _run = run;
            _properties = properties;
            Background = Brushes.White;
            Button B(string text) => new Button { Content = Dialogs.M(text), MinWidth = 94 };
            _add = B("&Add...");
            _attach = B("Attach Disk...");
            _activate = B("Deactivate");
            _move = B("Move...");
            _detach = B("Detach");
            _delete = B("&Delete");
            _props = B("P&roperties");
            _move.IsEnabled = false; // MoveVirtualDiskDialog not ported yet
            _add.Click += async (_, _) => await AddDisk();
            _attach.Click += async (_, _) => await AttachDisk();
            _activate.Click += (_, _) => ToggleActive();
            _detach.Click += async (_, _) => await Detach();
            _delete.Click += async (_, _) => await Delete();
            _props.Click += (_, _) => { if (SelectedVdi is VDI v) _properties(v); };

            var bar = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(0, 12, 0, 0) };
            foreach (var b in new[] { _add, _attach, _activate, _move, _detach, _delete, _props })
                bar.Children.Add(b);
            var table = new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Background = Brushes.White, Child = new ScrollViewer { Content = _rows, HorizontalScrollBarVisibility = Avalonia.Controls.Primitives.ScrollBarVisibility.Auto } };
            var grid = new Grid { Margin = new Thickness(12), RowDefinitions = new RowDefinitions("Auto,3*,Auto,*"), MaxWidth = 924, HorizontalAlignment = HorizontalAlignment.Left };
            Grid.SetRow(table, 1);
            Grid.SetRow(bar, 2);
            grid.Children.Add(_dvd);
            grid.Children.Add(table);
            grid.Children.Add(bar);
            Children.Add(grid);
            Refresh();
        }

        // ---------- rules (as the WinForms commands' CanRun) ----------





        private List<VBD> DiskVbds() =>
            _vm.Connection.ResolveAll(_vm.VBDs).Where(b => b.type != vbd_type.CD && b.type != vbd_type.Floppy).OrderBy(b => b.userdevice, StringComparer.Ordinal).ToList();

        private VBD SelectedVbd => DiskVbds().FirstOrDefault(b => b.opaque_ref == _selected);
        private VDI SelectedVdi => SelectedVbd is VBD b ? _vm.Connection.Resolve(b.VDI) : null;

        public void Refresh()
        {
            RefreshDvd();
            var vbds = DiskVbds();
            if (vbds.All(b => b.opaque_ref != _selected))
                _selected = vbds.FirstOrDefault()?.opaque_ref;
            _rows.Children.Clear();
            _rows.RowDefinitions.Clear();
            _rows.ColumnDefinitions = new ColumnDefinitions("Auto,Auto,*,Auto,Auto,Auto,Auto,Auto,Auto");
            var headers = new[] { "Position", "Name", "Description", "SR", "Size", "Read Only", "Priority", "Active", "Device Path" };
            _rows.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            for (int c = 0; c < headers.Length; c++)
                Cell(0, c, headers[c], null);
            int r = 1;
            foreach (var vbd in vbds)
            {
                var vdi = _vm.Connection.Resolve(vbd.VDI);
                if (vdi == null) continue;
                var sr = _vm.Connection.Resolve(vdi.SR);
                _rows.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var values = new[]
                {
                    vbd.userdevice, vdi.Name(), string.IsNullOrEmpty(vdi.Description()) ? "No description" : vdi.Description(),
                    sr?.NameWithLocation() ?? "", vdi.SizeText(), vbd.IsReadOnly() ? Messages.YES : Messages.NO,
                    vbd.GetIoNice() switch { 0 => Messages.STORAGE_WORST, 7 => Messages.STORAGE_BEST, var p => p.ToString() },
                    vbd.currently_attached ? Messages.YES : Messages.NO, string.IsNullOrEmpty(vbd.device) ? "" : "/dev/" + vbd.device
                };
                for (int c = 0; c < values.Length; c++)
                    Cell(r, c, values[c], vbd);
                r++;
            }

            var sel = SelectedVbd;
            var selVdi = SelectedVdi;
            var room = _vm.VBDs.Count < _vm.MaxVBDsAllowed() && !_vm.Locked;
            _add.IsEnabled = room;
            _attach.IsEnabled = room;
            _activate.Content = sel != null && !sel.currently_attached ? "Activate" : "Deactivate";
            _activate.IsEnabled = sel != null && (sel.currently_attached ? StorageRules.CanDeactivate(_vm, sel) : StorageRules.CanActivate(_vm, sel));
            _detach.IsEnabled = selVdi != null && StorageRules.CanDetach(_vm, selVdi);
            _delete.IsEnabled = selVdi != null && StorageRules.CanDelete(selVdi);
            _props.IsEnabled = selVdi != null && !selVdi.Locked;
        }

        private void Cell(int row, int col, string text, VBD vbd)
        {
            var header = vbd == null;
            var selected = vbd != null && vbd.opaque_ref == _selected;
            var b = new Border
            {
                Background = selected ? Brush.Parse("#0078D7") : Brushes.White,
                BorderBrush = Brush.Parse("#E3E3E3"),
                BorderThickness = header ? new Thickness(0, 0, 1, 1) : new Thickness(0),
                Padding = new Thickness(6, 3, 14, 3),
                Child = new TextBlock { Text = text ?? "", Foreground = selected ? Brushes.White : Brushes.Black, TextTrimming = TextTrimming.CharacterEllipsis }
            };
            if (vbd != null)
            {
                b.PointerPressed += (_, _) => { _selected = vbd.opaque_ref; Refresh(); };
                b.DoubleTapped += (_, _) => { if (SelectedVdi is VDI v) _properties(v); };
            }
            Grid.SetRow(b, row);
            Grid.SetColumn(b, col);
            _rows.Children.Add(b);
        }

        // ---------- DVD drive (CD changer) ----------

        private void RefreshDvd()
        {
            var cd = _vm.Connection.ResolveAll(_vm.VBDs).FirstOrDefault(b => b.IsCDROM());
            if (cd == null)
            {
                var link = Views.LinkText("Click here to create a DVD drive", _vm.IsRealVm() && !_vm.Locked ? () => _run(new CreateCdDriveAction(_vm)) : null);
                link.TextDecorations = TextDecorations.Underline;
                link.FontWeight = FontWeight.Bold;
                link.HorizontalAlignment = HorizontalAlignment.Center;
                link.Margin = new Thickness(0, 12, 0, 18);
                _dvd.Content = link;
                return;
            }
            var isos = IsoChoices();
            if (_dvd.Content is StackPanel existing && existing.Tag as string == cd.opaque_ref && existing.Children[1] is ComboBox combo)
            {
                // keep the open combo; just follow the current ISO
                _updatingDvd = true;
                combo.ItemsSource = isos;
                combo.SelectedItem = isos.FirstOrDefault(i => i.Vdi?.opaque_ref == (cd.empty ? null : cd.VDI.opaque_ref)) ?? isos[0];
                combo.IsEnabled = !cd.Locked;
                if (existing.Children[2] is Button ejectButton)
                    ejectButton.IsEnabled = !cd.empty;
                _updatingDvd = false;
                return;
            }
            var box = new ComboBox { ItemsSource = isos, Width = 360, IsEnabled = !cd.Locked };
            box.SelectedItem = isos.FirstOrDefault(i => i.Vdi?.opaque_ref == (cd.empty ? null : cd.VDI.opaque_ref)) ?? isos[0];
            box.SelectionChanged += (_, _) =>
            {
                if (_updatingDvd || !(box.SelectedItem is IsoChoice choice)) return;
                var current = _vm.Connection.Resolve(cd.VDI);
                if ((choice.Vdi?.opaque_ref) == (cd.empty ? null : current?.opaque_ref)) return;
                _run(new ChangeVMISOAction(_vm.Connection, _vm, choice.Vdi, cd));
            };
            var eject = new Button { Content = "Eject", IsEnabled = !cd.empty };
            eject.Click += (_, _) => _run(new ChangeVMISOAction(_vm.Connection, _vm, null, cd));
            var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(0, 10, 0, 14), Tag = cd.opaque_ref };
            row.Children.Add(new TextBlock { Text = $"DVD Drive {cd.userdevice}:", VerticalAlignment = VerticalAlignment.Center });
            row.Children.Add(box);
            row.Children.Add(eject);
            _dvd.Content = row;
        }

        public sealed class IsoChoice
        {
            public VDI Vdi;
            public string Text;
            public override string ToString() => Text;
        }

        /// <summary>"&lt;empty&gt;" plus ISOs of ISO SRs visible from the VM's home server.</summary>
        private List<IsoChoice> IsoChoices()
        {
            var list = new List<IsoChoice> { new IsoChoice { Text = "<empty>" } };
            var home = _vm.Home();
            foreach (var sr in _vm.Connection.Cache.SRs.Where(s => s.content_type == SR.Content_Type_ISO && s.Show(false)
                         && (home == null || s.CanBeSeenFrom(home))).OrderBy(s => s.Name()))
                foreach (var vdi in _vm.Connection.ResolveAll(sr.VDIs).Where(v => v.Show(false)).OrderBy(v => v.Name()))
                    list.Add(new IsoChoice { Vdi = vdi, Text = $"{sr.Name()}: {vdi.Name()}" });
            return list;
        }

        // ---------- commands ----------

        private Window Owner => TopLevel.GetTopLevel(this) as Window;

        private void ToggleActive()
        {
            var vbd = SelectedVbd;
            if (vbd == null) return;
            var plug = !vbd.currently_attached;
            var name = SelectedVdi?.Name() ?? vbd.userdevice;
            _run(new DelegatedAsyncAction(vbd.Connection,
                plug ? $"Activating '{name}'" : $"Deactivating '{name}'", "", "",
                session => { if (plug) VBD.plug(session, vbd.opaque_ref); else VBD.unplug(session, vbd.opaque_ref); },
                plug ? "vbd.plug" : "vbd.unplug"));
        }

        private async Task Detach()
        {
            var vdi = SelectedVdi;
            var vbd = SelectedVbd;
            if (vdi == null) return;
            var system = vbd?.bootable == true || vdi.type == vdi_type.system;
            if (await Dialogs.Confirm(Owner, system ? Messages.MESSAGEBOX_DETACH_SYSTEMVD_TITLE : Messages.MESSAGEBOX_DETACH_VD_TITLE,
                    system ? Messages.MESSAGEBOX_DETACH_SYSTEMVD : Messages.MESSAGEBOX_DETACH_VD, Messages.YES, Messages.NO))
                _run(new DetachVirtualDiskAction(vdi, _vm, true));
        }

        private async Task Delete()
        {
            var vdi = SelectedVdi;
            if (vdi == null) return;
            var system = vdi.VDIType() == VDI.FriendlyType.SYSTEM_DISK;
            if (await Dialogs.Confirm(Owner, system ? Messages.MESSAGEBOX_DELETE_SYS_DISK_TITLE : Messages.MESSAGEBOX_DELETE_VD_TITLE,
                    system ? Messages.MESSAGEBOX_DELETE_SYS_DISK : Messages.MESSAGEBOX_DELETE_VD, Messages.YES, Messages.NO))
                _run(new DestroyDiskAction(vdi));
        }

        private async Task AddDisk()
        {
            var action = await Dialogs.NewDisk(Owner, _vm);
            if (action != null) _run(action);
        }

        private async Task AttachDisk()
        {
            var action = await Dialogs.AttachDisk(Owner, _vm);
            if (action != null) _run(action);
        }
    }
}
