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

namespace XenAdmin.Desktop
{
    /// <summary>
    /// GPU tab (GpuPage): the placement policy of the pool with Edit (PoolGpuEditPage), then one row per group of GPUs
    /// with the same settings (GpuRow): the GPUs by server, the running vGPUs, the vGPU types each GPU allows, and
    /// Edit (GpuConfiguration) to choose them.
    /// </summary>
    public class GpuView : StackPanel, ILiveView
    {
        private readonly IXenObject _o;
        private readonly Action<AsyncAction> _run;
        private string _signature;

        public GpuView(IXenObject o, Action<AsyncAction> run)
        {
            _o = o;
            _run = run;
            Margin = new Thickness(12);
            Spacing = 14;
            MaxWidth = 1100;
            HorizontalAlignment = HorizontalAlignment.Left;
            Refresh();
        }

        private Window Owner => TopLevel.GetTopLevel(this) as Window;

        public void Refresh()
        {
            var c = _o.Connection;
            var groups = HardwareRules.Groups(_o);
            var signature = HardwareRules.PlacementPolicy(c) + "|" + HardwareRules.ShowPlacementPolicy(c) + "|" + string.Join(";", groups.SelectMany(g => g).Select(p =>
                $"{p.opaque_ref}:{string.Join(",", p.enabled_VGPU_types)}:{string.Join(",", p.supported_VGPU_types)}:{p.resident_VGPUs.Count}"))
                + "|" + string.Join(",", groups.Select(HardwareRules.CanEditTypes));
            if (signature == _signature) return;
            _signature = signature;
            Children.Clear();

            if (HardwareRules.ShowPlacementPolicy(c))
            {
                var row = new DockPanel();
                var edit = new Button { Content = Dialogs.M(Messages.GPU_EDIT_ALLOWED_TYPES_SINGLE), MinWidth = 80 };
                edit.Click += async (_, _) => { if (await PlacementDialog(Owner, c) is AsyncAction a) _run(a); };
                DockPanel.SetDock(edit, Dock.Right);
                row.Children.Add(edit);
                row.Children.Add(new TextBlock
                {
                    Text = string.Format(Messages.GPU_PLACEMENT_POLICY_DESCRIPTION, HardwareRules.PolicyText(HardwareRules.PlacementPolicy(c))),
                    VerticalAlignment = VerticalAlignment.Center, TextWrapping = TextWrapping.Wrap
                });
                Children.Add(Box(row));
            }
            foreach (var g in groups)
                Children.Add(Box(GroupRow(g)));
            if (groups.Count == 0)
                Children.Add(new TextBlock { Text = HardwareRules.NoGpuText(_o), TextWrapping = TextWrapping.Wrap });
        }

        private static Border Box(Control child) =>
            new Border { BorderBrush = Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1), Padding = new Thickness(10), Background = Brushes.White, Child = child };

        private Control GroupRow(List<PGPU> pgpus)
        {
            var panel = new StackPanel { Spacing = 6 };
            foreach (var byHost in pgpus.GroupBy(p => p.Connection.Resolve(p.host)))
            {
                panel.Children.Add(new TextBlock { Text = string.Format(Messages.GPU_ON_HOST_LABEL, byHost.Key?.Name()), FontWeight = FontWeight.SemiBold });
                foreach (var p in byHost)
                    panel.Children.Add(new TextBlock { Text = $"    {p.Name()}   ({p.resident_VGPUs.Count} vGPUs running)" });
            }
            var types = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*,Auto,Auto"), ColumnSpacing = 14, RowSpacing = 3, Margin = new Thickness(0, 6, 0, 0) };
            types.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            void Cell(Control ctl, int r, int col) { Grid.SetRow(ctl, r); Grid.SetColumn(ctl, col); types.Children.Add(ctl); }
            Cell(new TextBlock { Text = "Allowed vGPU types", FontWeight = FontWeight.SemiBold }, 0, 1);
            Cell(new TextBlock { Text = "Max per GPU", FontWeight = FontWeight.SemiBold }, 0, 2);
            Cell(new TextBlock { Text = "Video RAM", FontWeight = FontWeight.SemiBold }, 0, 3);
            foreach (var t in HardwareRules.TypeRows(pgpus))
            {
                var r = types.RowDefinitions.Count;
                types.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                Cell(Views.Icon(t.Enabled ? "000_Tick_h32bit_16" : "000_Abort_h32bit_16"), r, 0);
                Cell(new TextBlock { Text = t.Name }, r, 1);
                Cell(new TextBlock { Text = t.MaxPerGpu > 0 ? t.MaxPerGpu.ToString() : "" }, r, 2);
                Cell(new TextBlock { Text = t.VideoRam > 0 ? Util.MemorySizeStringSuitableUnits(t.VideoRam, true) : "" }, r, 3);
            }
            panel.Children.Add(types);
            if (HardwareRules.CanEditTypes(pgpus))
            {
                var edit = new Button
                {
                    Content = Dialogs.M(pgpus.Count > 1 ? Messages.GPU_EDIT_ALLOWED_TYPES_MULTIPLE : Messages.GPU_EDIT_ALLOWED_TYPES_SINGLE),
                    HorizontalAlignment = HorizontalAlignment.Left, Margin = new Thickness(0, 6, 0, 0)
                };
                edit.Click += async (_, _) => { if (await TypesDialog(Owner, pgpus) is AsyncAction a) _run(a); };
                panel.Children.Add(edit);
            }
            return panel;
        }

        /// <summary>PoolGpuEditPage: maximum density or maximum performance.</summary>
        private static async System.Threading.Tasks.Task<AsyncAction> PlacementDialog(Window owner, XenAdmin.Network.IXenConnection c)
        {
            var current = HardwareRules.PlacementPolicy(c);
            var density = new RadioButton { Content = Messages.GPU_PLACEMENT_POLICY_MAX_DENSITY_DESCRIPTION, GroupName = "gpu", IsChecked = current == allocation_algorithm.depth_first };
            var performance = new RadioButton { Content = Messages.GPU_PLACEMENT_POLICY_MAX_PERFORMANCE_DESCRIPTION, GroupName = "gpu", IsChecked = current == allocation_algorithm.breadth_first };
            var ok = new Button { Content = "OK", IsDefault = true, MinWidth = 75 };
            var cancel = new Button { Content = Dialogs.M(Messages.CANCEL), IsCancel = true, MinWidth = 75 };
            var panel = new StackPanel { Spacing = 8 };
            panel.Children.Add(new TextBlock { Text = "Choose how VMs with vGPUs are placed on the GPUs of this pool:", TextWrapping = TextWrapping.Wrap });
            panel.Children.Add(density);
            panel.Children.Add(performance);
            panel.Children.Add(Dialogs.Buttons(ok, cancel));
            var w = Dialogs.Frame("GPU Placement Policy", panel, 520);
            AsyncAction result = null;
            ok.Click += (_, _) =>
            {
                var chosen = density.IsChecked == true ? allocation_algorithm.depth_first
                    : performance.IsChecked == true ? allocation_algorithm.breadth_first : allocation_algorithm.unknown;
                result = HardwareRules.PlacementAction(c, chosen);
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }

        /// <summary>GpuConfiguration: tick the vGPU types the GPUs allow; types in use cannot be cleared.</summary>
        private static async System.Threading.Tasks.Task<AsyncAction> TypesDialog(Window owner, List<PGPU> pgpus)
        {
            var rows = HardwareRules.TypeRows(pgpus);
            var boxes = rows.Select(t => (t, Box: new CheckBox
            {
                Content = t.InUse ? $"{t.Name}  (in use)" : t.Name,
                IsChecked = t.Enabled, IsEnabled = !t.InUse
            })).ToList();
            var ok = new Button { Content = "OK", IsDefault = true, MinWidth = 75 };
            var cancel = new Button { Content = Dialogs.M(Messages.CANCEL), IsCancel = true, MinWidth = 75 };
            var list = new StackPanel { Spacing = 4 };
            boxes.ForEach(b => list.Children.Add(b.Box));
            var name = pgpus[0].Name();
            var panel = new StackPanel { Spacing = 8 };
            panel.Children.Add(new TextBlock
            {
                Text = pgpus.Count == 1 ? string.Format(Messages.GPU_GROUP_NAME_AND_NO_OF_GPUS_ONE, name) : string.Format(Messages.GPU_GROUP_NAME_AND_NO_OF_GPUS, name, pgpus.Count),
                FontWeight = FontWeight.SemiBold
            });
            panel.Children.Add(new TextBlock { Text = pgpus.Count == 1 ? Messages.GPU_RUBRIC_PLEASE_SELECT_WHICH_GPU_ONE : Messages.GPU_RUBRIC_PLEASE_SELECT_WHICH_GPU_MULTIPLE, TextWrapping = TextWrapping.Wrap });
            panel.Children.Add(new ScrollViewer { Content = list, MaxHeight = 300 });
            panel.Children.Add(Dialogs.Buttons(ok, cancel));
            var w = Dialogs.Frame("GPU Configuration", panel, 560);
            AsyncAction result = null;
            ok.Click += (_, _) =>
            {
                var changes = boxes.Where(b => b.Box.IsEnabled && (b.Box.IsChecked == true) != b.t.Enabled)
                    .ToDictionary(b => b.t.Ref.opaque_ref, b => b.Box.IsChecked == true);
                result = HardwareRules.ConfigureAction(pgpus, changes);
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }
    }

    /// <summary>USB tab (UsbPage): the server's USB devices with Enable / Disable pass-through.</summary>
    public class UsbView : DockPanel, ILiveView
    {
        private readonly Host _host;
        private readonly Action<AsyncAction> _run;
        private readonly SelectableGrid _grid = new SelectableGrid();
        private readonly Button _passthrough = new Button { MinWidth = 140 };

        public UsbView(Host host, Action<AsyncAction> run)
        {
            _host = host;
            _run = run;
            Background = Brushes.White;
            _grid.SelectionChanged += UpdateButton;
            _passthrough.Click += async (_, _) => await Toggle();
            var title = new TextBlock { Text = Messages.USB_DEVICES, FontWeight = FontWeight.SemiBold, Margin = new Thickness(0, 0, 0, 6) };
            var bar = new StackPanel { Orientation = Orientation.Horizontal, Margin = new Thickness(0, 10, 0, 0) };
            bar.Children.Add(_passthrough);
            var grid = new Grid { Margin = new Thickness(12), RowDefinitions = new RowDefinitions("Auto,3*,Auto,*"), MaxWidth = 1200, HorizontalAlignment = HorizontalAlignment.Left };
            Grid.SetRow(_grid, 1);
            Grid.SetRow(bar, 2);
            grid.Children.Add(title);
            grid.Children.Add(_grid);
            grid.Children.Add(bar);
            Children.Add(grid);
            Refresh();
        }

        private PUSB Selected => HardwareRules.Pusbs(_host).FirstOrDefault(p => p.opaque_ref == _grid.Selected);

        public void Refresh()
        {
            var pusbs = HardwareRules.Pusbs(_host);
            var table = new Table { Headers = new[] { "Location", "Description", "Pass-through", "Virtual machine" }, Fill = 1 };
            foreach (var p in pusbs)
                table.Rows.Add(new[] { p.path, p.Description(), p.passthrough_enabled ? Messages.ENABLED : Messages.DISABLED, HardwareRules.AttachedVm(p)?.name_label ?? "" });
            _grid.Fill(table, pusbs.Select(p => p.opaque_ref).ToList(), _grid.Selected);
            UpdateButton();
        }

        private void UpdateButton()
        {
            var (text, enabled) = HardwareRules.PassthroughButton(Selected);
            _passthrough.Content = Dialogs.M(text);
            _passthrough.IsEnabled = enabled && !_host.Locked;
        }

        private async System.Threading.Tasks.Task Toggle()
        {
            if (!(Selected is PUSB pusb)) return;
            var wasEnabled = pusb.passthrough_enabled; // cached objects change in place while the question is open
            var (text, yes) = HardwareRules.PassthroughConfirmation(pusb);
            if (!await Dialogs.Confirm(TopLevel.GetTopLevel(this) as Window, Messages.USB, text, yes, Messages.CANCEL))
                return;
            var now = _host.Connection.Resolve(new XenRef<PUSB>(pusb.opaque_ref));
            if (now == null || now.passthrough_enabled != wasEnabled || wasEnabled && HardwareRules.IsAttached(now))
                return; // changed meanwhile: the user confirmed something else
            _run(HardwareRules.PassthroughAction(now, !wasEnabled));
        }
    }
}
