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
    /// Memory tab (BallooningPage): a bar per VM (used memory inside the dynamic range, VMShinyBar) with Edit for the
    /// Memory Settings dialog, or per server (Xen and control domain, each VM, free memory, HostShinyBar), above the
    /// figures of Details.Memory.
    /// </summary>
    public class MemoryView : StackPanel, ILiveView
    {
        private static readonly IBrush Used = Brush.Parse("#4E9A06"), Range = Brush.Parse("#C9E7B0"), Empty = Brush.Parse("#F2F2F2"), Edge = Brush.Parse("#828790");
        private static readonly IBrush[] VmColours = { Brush.Parse("#3C78B4"), Brush.Parse("#5A9BD5"), Brush.Parse("#7FB3E0"), Brush.Parse("#A5C9EA") };

        private readonly IXenObject _o;
        private readonly Action<AsyncAction> _run;

        public MemoryView(IXenObject o, Action<AsyncAction> run)
        {
            _o = o;
            _run = run;
            Margin = new Thickness(12);
            Spacing = 12;
            Refresh();
        }

        public void Refresh()
        {
            Children.Clear();
            switch (_o)
            {
                case VM vm:
                    Children.Add(VmRow(vm));
                    break;
                case Host host:
                    Children.Add(HostBar(host));
                    break;
                case Pool pool:
                    foreach (var h in pool.Connection.Cache.Hosts.OrderBy(h => h.Name()))
                    {
                        Children.Add(new TextBlock { Text = h.Name(), FontWeight = FontWeight.SemiBold });
                        Children.Add(HostBar(h));
                    }
                    break;
            }
            Children.Add(Views.Sections(Details.Memory(_o)));
        }

        private Control VmRow(VM vm)
        {
            var (used, dmin, dmax) = MemoryRules.VmBar(vm);
            var bar = new Grid { Height = 26, Width = 600, HorizontalAlignment = HorizontalAlignment.Left, ColumnDefinitions = new ColumnDefinitions("*") };
            bar.Children.Add(new Border { Background = Empty, BorderBrush = Edge, BorderThickness = new Thickness(1) });
            // the dynamic range, then the memory in use on top (VMShinyBar)
            bar.Children.Add(new Border { Background = Range, Width = Math.Max(0, (dmax - dmin) * 600), Margin = new Thickness(dmin * 600, 1, 0, 1), HorizontalAlignment = HorizontalAlignment.Left });
            bar.Children.Add(new Border { Background = Used, Width = used * 600, Margin = new Thickness(0, 1), HorizontalAlignment = HorizontalAlignment.Left });
            ToolTip.SetTip(bar, used > 0 ? $"In use: {XenAdmin.Util.MemorySizeStringSuitableUnits((long)(used * vm.memory_static_max), true)}" : "The VM is not running");
            var edit = new Button { Content = Dialogs.M("&Edit..."), IsEnabled = MemoryRules.CanEdit(vm), VerticalAlignment = VerticalAlignment.Center };
            edit.Click += async (_, _) =>
            {
                if (await Dialogs.MemorySettings(TopLevel.GetTopLevel(this) as Window, vm) is AsyncAction a) _run(a);
            };
            var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 12 };
            row.Children.Add(bar);
            row.Children.Add(edit);
            return row;
        }

        private static Control HostBar(Host host)
        {
            var segments = MemoryRules.HostSegments(host);
            double total = Math.Max(1, segments.Sum(s => s.Bytes));
            var bar = new StackPanel { Orientation = Orientation.Horizontal, Height = 26, HorizontalAlignment = HorizontalAlignment.Left };
            for (int i = 0; i < segments.Count; i++)
            {
                var (name, bytes) = segments[i];
                var free = i == segments.Count - 1;
                var cell = new Border
                {
                    Width = Math.Max(1, bytes / total * 600), Background = i == 0 ? Brush.Parse("#2E5E8A") : free ? Empty : VmColours[(i - 1) % VmColours.Length],
                    BorderBrush = Edge, BorderThickness = new Thickness(i == 0 ? 1 : 0, 1, free ? 1 : 0, 1)
                };
                ToolTip.SetTip(cell, $"{name}: {XenAdmin.Util.MemorySizeStringSuitableUnits(bytes, true)}");
                bar.Children.Add(cell);
            }
            return bar;
        }
    }
}
