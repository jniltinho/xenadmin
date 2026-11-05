using System;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using Avalonia.Threading;
using XenAdmin.Actions;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// VtpmManagementDialog: the VM's vTPMs on the left, the selected one's properties (Protected, Unique) with Remove
    /// vTPM on the right, Add vTPM and Close. It follows the VM while open (power state, vTPMs).
    /// </summary>
    public class VtpmDialog : Window
    {
        private readonly VM _vm;
        private readonly Action<AsyncAction> _run;
        private readonly ListBox _list = new ListBox { Width = 180 };
        private readonly ContentControl _page = new ContentControl();
        private readonly Button _add = new Button { Content = Dialogs.M("A&dd vTPM") };
        private readonly TextBlock _addReason = new TextBlock { Opacity = 0.7, TextWrapping = TextWrapping.Wrap, VerticalAlignment = VerticalAlignment.Center };
        private readonly DispatcherTimer _timer = new DispatcherTimer { Interval = TimeSpan.FromSeconds(1) };
        private string _signature;
        private bool _busy; // an add or remove is running (Windows shows it in a modal progress dialog)

        public VtpmDialog(VM vm, Action<AsyncAction> run)
        {
            _vm = vm;
            _run = run;
            Title = $"Trusted Platform Module (TPM) Manager - {vm.Name().Ellipsise(50)}";
            Width = 640;
            Height = 400;
            WindowStartupLocation = WindowStartupLocation.CenterOwner;
            ShowInTaskbar = false;

            _add.Click += (_, _) => { if (!_busy && VtpmRules.AddBlocker(_vm) == null) Start(VtpmRules.AddAction(_vm)); };
            _list.SelectionChanged += (_, _) => ShowSelected();
            var close = new Button { Content = "Close", IsCancel = true, MinWidth = 75 };
            close.Click += (_, _) => Close();

            var bottom = new DockPanel { Margin = new Thickness(0, 10, 0, 0) };
            DockPanel.SetDock(close, Dock.Right);
            var left = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            left.Children.Add(_add);
            left.Children.Add(_addReason);
            bottom.Children.Add(close);
            bottom.Children.Add(left);

            var body = new DockPanel();
            var listBorder = new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = _list, Margin = new Thickness(0, 0, 10, 0) };
            DockPanel.SetDock(listBorder, Dock.Left);
            body.Children.Add(listBorder);
            body.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Padding = new Thickness(12), Background = Brushes.White, Child = _page });

            var root = new DockPanel { Margin = new Thickness(12) };
            var blurb = new TextBlock { Text = $"Attach, remove, and configure vTPMs on '{vm.Name().Ellipsise(80)}'.", Margin = new Thickness(0, 0, 0, 10) };
            DockPanel.SetDock(blurb, Dock.Top);
            DockPanel.SetDock(bottom, Dock.Bottom);
            root.Children.Add(blurb);
            root.Children.Add(bottom);
            root.Children.Add(body);
            Content = root;

            _timer.Tick += (_, _) => Fill();
            Opened += (_, _) => { Fill(); _timer.Start(); };
            Closed += (_, _) => _timer.Stop();
        }

        private VTPM Selected => _list.SelectedItem is ListBoxItem { Tag: string r } ? _vm.Connection.Resolve(new XenRef<VTPM>(r)) : null;

        private void Start(AsyncAction action)
        {
            _busy = true;
            _signature = null;
            action.Completed += _ => Dispatcher.UIThread.Post(() => { _busy = false; _signature = null; Fill(); });
            _run(action);
            Fill();
        }

        private void Fill()
        {
            // the connection was lost or the VM is gone (a reconnect builds a new cache): this window is stale
            if (!_vm.Connection.IsConnected || !ReferenceEquals(_vm.Connection.Resolve(new XenRef<VM>(_vm.opaque_ref)), _vm))
            {
                Close();
                return;
            }
            var vtpms = VtpmRules.Vtpms(_vm);
            var signature = $"{_vm.power_state}|{string.Join(",", _vm.allowed_operations)}|" +
                            string.Join(",", vtpms.Select(t => $"{t.opaque_ref}:{t.is_protected}:{t.is_unique}:{string.Join("+", t.allowed_operations)}"));
            if (signature == _signature) return;
            _signature = signature;
            var selected = (_list.SelectedItem as ListBoxItem)?.Tag as string;
            _list.ItemsSource = vtpms.Select((t, i) => new ListBoxItem { Content = vtpms.Count == 1 ? "vTPM" : $"vTPM {i + 1}", Tag = t.opaque_ref }).ToList();
            var items = (System.Collections.Generic.List<ListBoxItem>)_list.ItemsSource;
            _list.SelectedItem = items.FirstOrDefault(i => (string)i.Tag == selected) ?? items.FirstOrDefault();
            var reason = VtpmRules.AddBlocker(_vm);
            _add.IsEnabled = reason == null && !_busy;
            _addReason.Text = reason ?? "";
            ShowSelected();
        }

        private void ShowSelected()
        {
            if (!(Selected is VTPM vtpm))
            {
                _page.Content = new TextBlock { Text = "There are no vTPMs attached to this VM.", VerticalAlignment = VerticalAlignment.Center, HorizontalAlignment = HorizontalAlignment.Center };
                return;
            }
            var (prot, unique) = VtpmRules.Properties(vtpm);
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowDefinitions = new RowDefinitions("Auto,Auto,Auto,Auto"), RowSpacing = 6, ColumnSpacing = 12 };
            void Row(int r, string label, string value, string hint)
            {
                var l = new TextBlock { Text = label, FontWeight = FontWeight.SemiBold };
                var v = new TextBlock { Text = value };
                var h = new TextBlock { Text = hint, Opacity = 0.7, TextWrapping = TextWrapping.Wrap, Margin = new Thickness(0, 0, 0, 6) };
                Grid.SetRow(l, r * 2); Grid.SetRow(v, r * 2); Grid.SetColumn(v, 1);
                Grid.SetRow(h, r * 2 + 1); Grid.SetColumnSpan(h, 2);
                grid.Children.AddRange(new Control[] { l, v, h });
            }
            Row(0, "Protected:", prot, "The contents of a protected vTPM are shielded from unauthorized access and undetected tampering when saved to persistent storage or sent via network.");
            Row(1, "Unique:", unique, "The contents of a unique vTPM cannot be duplicated through any VM operation such as copy, clone, export, or creation from a snapshot.");
            var remove = new Button { Content = Dialogs.M("&Remove vTPM"), HorizontalAlignment = HorizontalAlignment.Left };
            var reason = VtpmRules.RemoveBlocker(_vm, vtpm);
            remove.IsEnabled = reason == null && !_busy;
            ToolTip.SetTip(remove, reason);
            remove.Click += async (_, _) =>
            {
                var target = vtpm.opaque_ref;
                if (!await Dialogs.Confirm(this, "Remove vTPM", VtpmRules.RemoveWarning, Messages.YES, Messages.NO)) return;
                if (!_busy && _vm.Connection.Resolve(new XenRef<VTPM>(target)) is VTPM now && VtpmRules.RemoveBlocker(_vm, now) == null)
                    Start(VtpmRules.RemoveAction(_vm, now));
            };
            var panel = new StackPanel { Spacing = 10 };
            panel.Children.Add(new TextBlock { Text = "Properties", FontWeight = FontWeight.Bold });
            panel.Children.Add(grid);
            panel.Children.Add(remove);
            if (reason != null)
                panel.Children.Add(new TextBlock { Text = reason, Opacity = 0.7, TextWrapping = TextWrapping.Wrap });
            _page.Content = panel;
        }
    }
}
