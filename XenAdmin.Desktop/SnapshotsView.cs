using System;
using System.Collections.Generic;
using System.Linq;
using Task = System.Threading.Tasks.Task;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Input;
using Avalonia.Layout;
using Avalonia.Media;
using XenAdmin.Actions;
using XenAdmin.Actions.VMActions;
using XenAdmin.Core;
using XenAdmin.Model;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>A tab body that refreshes itself instead of being rebuilt on every cache event.</summary>
    public interface ILiveView
    {
        void Refresh();
    }

    /// <summary>SnapshotsPage (list view): Take Snapshot, Revert To, Actions, Delete, details panel.</summary>
    public class SnapshotsView : DockPanel, ILiveView
    {
        private readonly VM _vm;
        private readonly Action<AsyncAction> _run;
        private readonly Action<IXenObject> _properties;
        private readonly Grid _rows = new Grid();
        private readonly TextBlock _empty = new TextBlock { Text = Messages.SNAPSHOTS_EMPTY.Replace(". ", ".\n"), FontWeight = FontWeight.Bold, TextAlignment = TextAlignment.Center, HorizontalAlignment = HorizontalAlignment.Center, VerticalAlignment = VerticalAlignment.Center };
        private readonly Button _take, _revert, _delete;
        private readonly Button _actions = new Button { Content = "Actions  ▾" };
        private readonly Button _view = new Button { Content = "View  ▾" };
        private readonly StackPanel _details = new StackPanel { Width = 262, Margin = new Thickness(6, 0, 0, 0) };
        private readonly TextBlock _created = new TextBlock { Opacity = 0.7 };
        private readonly TextBlock _name = new TextBlock { TextWrapping = TextWrapping.Wrap, FontWeight = FontWeight.SemiBold };
        private readonly TextBlock _size = new TextBlock();
        private readonly TextBlock _desc = new TextBlock { TextWrapping = TextWrapping.Wrap };
        private readonly Image _scheduleIcon = Views.Icon("000_Info3_h32bit_16");
        private readonly TextBlock _scheduleText = new TextBlock { TextWrapping = TextWrapping.Wrap, MaxWidth = 230 };
        private readonly StackPanel _schedule = new StackPanel { Margin = new Thickness(6, 0, 0, 10), Spacing = 2, VerticalAlignment = VerticalAlignment.Bottom };
        private ContentControl _rightHost;
        private readonly TextBlock _type = new TextBlock();
        private readonly TextBlock _tags = new TextBlock { TextWrapping = TextWrapping.Wrap };
        private readonly TextBlock _folder = new TextBlock();
        private readonly Button _props = new Button { Content = Dialogs.M("P&roperties"), HorizontalAlignment = HorizontalAlignment.Right };
        private string _selected;

        public SnapshotsView(VM vm, Action<AsyncAction> run, Action<IXenObject> properties)
        {
            _vm = vm;
            _run = run;
            _properties = properties;
            Background = Brushes.White;

            _take = new Button { Content = Dialogs.M("Ta&ke Snapshot...") };
            _revert = new Button { Content = Dialogs.M("R&evert To...") };
            _delete = new Button { Content = Dialogs.M("&Delete") };
            _take.Click += async (_, _) => await TakeSnapshot();
            _revert.Click += async (_, _) => await Revert();
            _delete.Click += async (_, _) => await Delete();
            _props.Click += (_, _) => { if (Selected is VM s) _properties(s); };
            _actions.Click += (_, _) => ActionsMenu().Open(_actions);
            _view.Click += (_, _) => new ContextMenu
            {
                ItemsSource = new object[]
                {
                    new MenuItem { Header = Dialogs.M("&Tree View"), IsEnabled = false },
                    new MenuItem { Header = Dialogs.M("&List View"), Icon = new TextBlock { Text = "✓" } },
                    new Separator(),
                    new MenuItem { Header = Dialogs.M("&Scheduled Snapshots"), IsEnabled = false }
                }
            }.Open(_view);

            var toolbar = new DockPanel { Margin = new Thickness(12, 12, 12, 8) };
            var left = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            foreach (var b in new[] { _take, _revert, _actions, _delete, _view })
            {
                b.Classes.Add("tall"); // this toolbar uses taller buttons than dialogs, as in SnapshotsPage
                left.Children.Add(b);
            }
            var hide = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6, Cursor = new Cursor(StandardCursorType.Hand), VerticalAlignment = VerticalAlignment.Center };
            var chevron = Views.Icon("PDChevronRight", 18);
            hide.Children.Add(chevron);
            var hideText = new TextBlock { Text = "Hide Details", VerticalAlignment = VerticalAlignment.Center };
            hide.Children.Add(hideText);
            hide.PointerPressed += (_, _) =>
            {
                // collapse the whole right column so the list takes the space
                _rightHost.IsVisible = !_rightHost.IsVisible;
                hideText.Text = _rightHost.IsVisible ? "Hide Details" : "Show Details";
                chevron.Source = Images.Get(_rightHost.IsVisible ? "PDChevronRight" : "PDChevronLeft");
            };
            DockPanel.SetDock(hide, Dock.Right);
            toolbar.Children.Add(hide);
            toolbar.Children.Add(left);

            // details panel: "Created on ...", screenshot box, description, type, tags, folder, Properties
            var screenshot = new Border { Width = 100, Height = 75, Background = Brushes.Black, HorizontalAlignment = HorizontalAlignment.Center, Margin = new Thickness(0, 8, 0, 60), Child = new TextBlock { Text = Messages.NO_SCREENSHOT, Foreground = Brushes.White, HorizontalAlignment = HorizontalAlignment.Center, VerticalAlignment = VerticalAlignment.Center } };
            var fields = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowSpacing = 12, ColumnSpacing = 8 };
            int r = 0;
            void Field(string label, Control value)
            {
                fields.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var l = new TextBlock { Text = label, Opacity = 0.6 };
                Grid.SetRow(l, r); Grid.SetRow(value, r); Grid.SetColumn(value, 1);
                fields.Children.Add(l); fields.Children.Add(value);
                r++;
            }
            Field("Name:", _name);
            Field("Description:", _desc);
            Field("Type:", _type);
            Field("Size:", _size);
            Field("Tags:", _tags);
            Field("Folder:", _folder);
            var box = new StackPanel { Margin = new Thickness(8, 4, 8, 10) };
            box.Children.Add(_created);
            box.Children.Add(screenshot);
            box.Children.Add(fields);
            box.Children.Add(new Border { Height = 18 });
            box.Children.Add(_props);
            _details.Children.Add(new Border { BorderBrush = Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1), Child = box });
            var info = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
            info.Children.Add(_scheduleIcon);
            info.Children.Add(_scheduleText);
            _schedule.Children.Add(info);
            var link = Views.LinkText(Messages.VIEW_VMSS_POLICIES, null);
            link.TextDecorations = TextDecorations.Underline;
            link.Margin = new Thickness(22, 0, 0, 0);
            _schedule.Children.Add(link);
            var right = new DockPanel();
            DockPanel.SetDock(_schedule, Dock.Bottom);
            right.Children.Add(_schedule);
            right.Children.Add(_details);
            var rightHost = _rightHost = new ContentControl { Content = right, Margin = new Thickness(0, 0, 12, 12) };

            var listBox = new Border
            {
                BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Background = Brushes.White, Margin = new Thickness(12, 0, 0, 12),
                Child = new Panel { Children = { new ScrollViewer { Content = _rows, HorizontalScrollBarVisibility = Avalonia.Controls.Primitives.ScrollBarVisibility.Auto }, _empty } }
            };
            DockPanel.SetDock(toolbar, Dock.Top);
            DockPanel.SetDock(rightHost, Dock.Right);
            Children.Add(toolbar);
            Children.Add(rightHost);
            Children.Add(listBox);
            Refresh();
        }

        private List<VM> Snapshots() =>
            _vm.Connection.ResolveAll(_vm.snapshots).Where(s => s != null).OrderByDescending(s => s.snapshot_time).ToList();

        private VM Selected => Snapshots().FirstOrDefault(s => s.opaque_ref == _selected);



        public void Refresh()
        {
            var snapshots = Snapshots();
            if (snapshots.All(s => s.opaque_ref != _selected))
                _selected = snapshots.FirstOrDefault()?.opaque_ref;
            _empty.IsVisible = snapshots.Count == 0;
            _rows.Children.Clear();
            _rows.RowDefinitions.Clear();
            _rows.ColumnDefinitions = new ColumnDefinitions("Auto,Auto,Auto,Auto,Auto,*");
            _rows.IsVisible = snapshots.Count > 0;
            var headers = new[] { "Type", "Name", "Created On", "Size", "Tags", "Description" };
            _rows.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            for (int c = 0; c < headers.Length; c++)
                Cell(0, c, headers[c], true, null);
            for (int i = 0; i < snapshots.Count; i++)
            {
                var s = snapshots[i];
                _rows.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var values = new[]
                {
                    VmRules.SnapshotTypeText(s), s.name_label, (s.snapshot_time.ToLocalTime() + s.Connection.ServerTimeOffset).ToString(Messages.DATEFORMAT_DMY_HMS),
                    Util.DiskSizeString(VmRules.SnapshotSize(s)), string.Join(", ", Tags.GetTags(s) ?? Array.Empty<string>()),
                    string.IsNullOrEmpty(s.Description()) ? Messages.NONE : s.Description()
                };
                for (int c = 0; c < values.Length; c++)
                    Cell(i + 1, c, values[c], false, s);
            }
            UpdateDetails();
        }

        private void Cell(int row, int col, string text, bool header, VM snapshot)
        {
            var selected = snapshot != null && snapshot.opaque_ref == _selected;
            var b = new Border
            {
                Background = selected ? Brush.Parse("#0078D7") : Brushes.White,
                BorderBrush = Brush.Parse("#E3E3E3"),
                BorderThickness = header ? new Thickness(0, 0, 1, 1) : new Thickness(0),
                Padding = new Thickness(6, 3, 14, 3),
                Child = new TextBlock { Text = text, Foreground = selected ? Brushes.White : Brushes.Black, TextTrimming = TextTrimming.CharacterEllipsis }
            };
            if (snapshot != null)
            {
                b.PointerPressed += (_, _) => { _selected = snapshot.opaque_ref; Refresh(); };
                b.DoubleTapped += (_, _) => _properties(snapshot);
            }
            Grid.SetRow(b, row);
            Grid.SetColumn(b, col);
            _rows.Children.Add(b);
        }

        private void UpdateDetails()
        {
            var s = Selected;
            _take.IsEnabled = VmRules.CanTakeSnapshot(_vm);
            _revert.IsEnabled = s != null && s.allowed_operations.Contains(vm_operations.revert);
            _delete.IsEnabled = s != null && s.allowed_operations.Contains(vm_operations.destroy);
            _actions.IsEnabled = s != null;
            _props.IsEnabled = s != null;
            _created.Text = s == null ? "Created on" : string.Format("Created on {0}", (s.snapshot_time.ToLocalTime() + s.Connection.ServerTimeOffset).ToString(Messages.DATEFORMAT_DMY_HMS));
            _name.Text = s?.name_label ?? "";
            _size.Text = s == null ? "" : Util.DiskSizeString(VmRules.SnapshotSize(s));
            RefreshSchedule();
            _desc.Text = s == null ? "" : string.IsNullOrEmpty(s.Description()) ? Messages.NONE : s.Description();
            _type.Text = s == null ? "" : VmRules.SnapshotTypeText(s);
            _tags.Text = s == null ? "" : (Tags.GetTags(s) is { Length: > 0 } t ? string.Join(", ", t) : Messages.NONE);
            _folder.Text = s == null ? "" : string.IsNullOrEmpty(s.Path) ? Messages.NONE : s.Path;
        }

        /// <summary>SnapshotsPage.RefreshVMSSPanel.</summary>
        private void RefreshSchedule()
        {
            _schedule.IsVisible = Helpers.FalconOrGreater(_vm.Connection);
            var forbidden = Helpers.FeatureForbidden(_vm.Connection, Host.RestrictVMSnapshotSchedule);
            var vmss = _vm.Connection.Resolve(_vm.snapshot_schedule);
            var inSchedule = vmss != null && !forbidden;
            _scheduleText.Text = inSchedule ? string.Format(Messages.THIS_VM_IS_IN_VMSS, vmss.Name()) : Messages.THIS_VM_IS_NOT_IN_VMSS;
            _scheduleIcon.Source = Images.Get(inSchedule ? "000_Tick_h32bit_16" : "000_Info3_h32bit_16");
        }

        private ContextMenu ActionsMenu()
        {
            var s = Selected;
            MenuItem Item(string text, bool enabled, Func<Task> run = null)
            {
                var m = new MenuItem { Header = Dialogs.M(text), IsEnabled = enabled && run != null };
                if (run != null) m.Click += async (_, _) => await run();
                return m;
            }
            return new ContextMenu
            {
                ItemsSource = new object[]
                {
                    Item("Ta&ke Snapshot...", _take.IsEnabled, TakeSnapshot),
                    Item("R&evert To Snapshot...", _revert.IsEnabled, Revert),
                    new Separator(),
                    Item("&New VM from Snapshot...", false),
                    Item("New Template from &Snapshot...", false),
                    Item("E&xport Snapshot as Template...", false),
                    new Separator(),
                    Item("&Delete", _delete.IsEnabled, Delete),
                    Item("P&roperties", s != null, () => { _properties(s); return Task.CompletedTask; })
                }
            };
        }

        private Window Owner => TopLevel.GetTopLevel(this) as Window;

        private async Task TakeSnapshot()
        {
            var action = await Dialogs.TakeSnapshot(Owner, _vm);
            if (action != null)
                _run(action);
        }

        /// <summary>RevertToSnapshotCommand: optional "snapshot current state first", then revert.</summary>
        private async Task Revert()
        {
            var s = Selected;
            if (s == null) return;
            var snapshotAllowed = _vm.allowed_operations.Contains(vm_operations.snapshot);
            var (yes, snapshotFirst) = await Dialogs.ConfirmWithCheckbox(Owner, "XCP-ng Center",
                string.Format(Messages.SNAPSHOT_REVERT_BLURB, s.Name()), snapshotAllowed ? Messages.SNAPSHOT_REVERT_NEW_SNAPSHOT : null, snapshotAllowed);
            if (!yes) return;
            if (snapshotFirst)
            {
                var take = await Dialogs.TakeSnapshot(Owner, _vm);
                if (take == null) return;
                take.Completed += a => { if (a.Succeeded) Avalonia.Threading.Dispatcher.UIThread.Post(() => _run(new VMSnapshotRevertAction(s))); };
                _run(take);
            }
            else
                _run(new VMSnapshotRevertAction(s));
        }

        private async Task Delete()
        {
            var s = Selected;
            if (s == null) return;
            if (await Dialogs.Confirm(Owner, "XCP-ng Center", string.Format(Messages.REMOVE_SNAPSHOT_TEXT, string.Format(Messages.THE_SNAPSHOT, s.Name())), Messages.YES, Messages.NO))
                _run(new VMSnapshotDeleteAction(s));
        }
    }
}
