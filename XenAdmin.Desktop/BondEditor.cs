using System;
using System.Collections.Generic;
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
    /// BondDetails (XenAdmin/Controls/BondDetails.cs): the NICs to bond (ticked up to the pool's limit; NICs already in a
    /// bond are dimmed), the bond mode, the MTU and "add to new VMs". Used by the New Network wizard and Create Bond.
    /// </summary>
    public class BondEditor : StackPanel
    {
        private readonly IXenObject _owner;
        private readonly List<PIF> _nics;
        private readonly List<CheckBox> _boxes = new List<CheckBox>();
        private readonly int _limit;
        private readonly RadioButton _activeActive, _activePassive, _lacpMac, _lacpPorts;
        private readonly TextBox _mtu = new TextBox { Width = 90, Text = XenAPI.Network.MTU_DEFAULT.ToString() };
        private readonly TextBlock _mtuInfo = new TextBlock { VerticalAlignment = VerticalAlignment.Center, Opacity = 0.8 };
        private readonly CheckBox _auto = new CheckBox { Content = Dialogs.M("A&utomatically add this network to new virtual machines"), IsChecked = true };

        /// <summary>The selection, mode or MTU changed (re-check IsValid).</summary>
        public event Action Changed;

        public BondEditor(IXenObject hostOrPool)
        {
            _owner = hostOrPool;
            _nics = BondRules.Candidates(hostOrPool);
            _limit = BondRules.SizeLimit(hostOrPool.Connection);
            Spacing = 10;

            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,Auto,Auto,Auto,Auto,*") };
            void Cell(int r, int c, Control content)
            {
                var b = new Border { Padding = new Thickness(6, 3, 12, 3), Child = content, BorderBrush = Brush.Parse("#E3E3E3"), BorderThickness = r == 0 ? new Thickness(0, 0, 1, 1) : new Thickness(0) };
                Grid.SetRow(b, r); Grid.SetColumn(b, c);
                grid.Children.Add(b);
            }
            grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            var headers = new[] { "", "NIC", "MAC", "Link Status", "Speed", "Vendor" };
            for (int c = 0; c < headers.Length; c++)
                Cell(0, c, new TextBlock { Text = headers[c] });
            for (int i = 0; i < _nics.Count; i++)
            {
                var pif = _nics[i];
                var box = new CheckBox { MinHeight = 0, Padding = new Thickness(0) };
                box.IsCheckedChanged += (_, _) => Update();
                _boxes.Add(box);
                grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var inBond = BondRules.InBond(pif);
                Cell(i + 1, 0, box);
                Cell(i + 1, 1, new TextBlock { Text = (pif.Name() + " " + inBond).Trim() });
                Cell(i + 1, 2, new TextBlock { Text = pif.MAC });
                var net = pif.Connection.Resolve(pif.network);
                Cell(i + 1, 3, new TextBlock { Text = net != null && net.PIFs.Count > 1 ? net.LinkStatusString() : pif.LinkStatusString() });
                Cell(i + 1, 4, new TextBlock { Text = pif.Carrier() ? pif.Speed() : Messages.HYPHEN });
                Cell(i + 1, 5, new TextBlock { Text = pif.PIFMetrics()?.vendor_name ?? "" });
            }
            Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Background = Brushes.White, Child = new ScrollViewer { Content = grid, MaxHeight = 150, HorizontalScrollBarVisibility = Avalonia.Controls.Primitives.ScrollBarVisibility.Auto } });

            var modes = new StackPanel { Spacing = 4 };
            modes.Children.Add(new TextBlock { Text = "Bond mode", FontWeight = FontWeight.SemiBold });
            RadioButton Mode(string text, bool on)
            {
                var r = new RadioButton { Content = Dialogs.M(text), GroupName = "bondmode", IsChecked = on };
                r.IsCheckedChanged += (_, _) => Update();
                modes.Children.Add(r);
                return r;
            }
            _activeActive = Mode("&Active-active", true);
            _activePassive = Mode("Active-pa&ssive", false);
            _lacpMac = Mode("&LACP with load balancing based on source MAC address", false);
            _lacpPorts = Mode("LA&CP with load balancing based on IP and port of source and destination", false);
            _lacpMac.IsVisible = _lacpPorts.IsVisible = BondRules.SupportsLacp(hostOrPool.Connection);
            var lacpWarning = new TextBlock { Text = "LACP must also be configured on the switch ports", Foreground = Brushes.DarkGoldenrod, IsVisible = false };
            _lacpMac.IsCheckedChanged += (_, _) => lacpWarning.IsVisible = _lacpMac.IsChecked == true || _lacpPorts.IsChecked == true;
            _lacpPorts.IsCheckedChanged += (_, _) => lacpWarning.IsVisible = _lacpMac.IsChecked == true || _lacpPorts.IsChecked == true;
            modes.Children.Add(lacpWarning);
            Children.Add(modes);

            var mtuRow = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            mtuRow.Children.Add(new Label { Content = Dialogs.M("&MTU:"), Target = _mtu, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) });
            mtuRow.Children.Add(_mtu);
            mtuRow.Children.Add(_mtuInfo);
            _mtu.TextChanged += (_, _) => Update();
            Children.Add(mtuRow);
            Children.Add(_auto);
            Update();
        }

        public List<PIF> Members => _nics.Where((_, i) => _boxes[i].IsChecked == true).ToList();

        public bond_mode Mode => _activePassive.IsChecked == true ? bond_mode.active_backup : _activeActive.IsChecked == true ? bond_mode.balance_slb : bond_mode.lacp;

        public Bond.hashing_algoritm Hashing => _lacpMac.IsChecked == true ? Bond.hashing_algoritm.src_mac
            : _lacpPorts.IsChecked == true ? Bond.hashing_algoritm.tcpudp_ports : Bond.hashing_algoritm.unknown;

        public bool IsValid => BondRules.ValidSelection(Members.Count, _limit) && !NetworkRules.ValidateMtu(_mtu.Text, BondRules.MaxMtu(Members)).Error;

        /// <summary>BondDetails.CanCreateBond, then CreateBondAction (null when blocked or cancelled).</summary>
        public async System.Threading.Tasks.Task<AsyncAction> CreateAction(Window owner)
        {
            var members = Members;
            if (BondRules.CreateWarning(_owner.Connection, members) is BondRules.Warning w)
            {
                if (w.Block)
                {
                    await Dialogs.Info(owner, w.Title, w.Text);
                    return null;
                }
                if (!await Dialogs.Confirm(owner, w.Title, w.Text, w.OkText, Messages.CANCEL))
                    return null;
            }
            return new CreateBondAction(_owner.Connection, BondRules.BondName(members), members, _auto.IsChecked == true,
                int.Parse(_mtu.Text.Trim()), Mode, Hashing);
        }

        private void Update()
        {
            var members = Members;
            // BondDetails.ToggleNICRowCheckable: ticked rows stay; the others only while the limit allows
            for (int i = 0; i < _nics.Count; i++)
                if (_boxes[i].IsChecked != true)
                {
                    var can = BondRules.CanTick(_nics[i], members.Count, _limit);
                    _boxes[i].IsEnabled = can;
                }
            var max = BondRules.MaxMtu(members);
            _mtu.IsEnabled = max > XenAPI.Network.MTU_MIN;
            // the NumericUpDown keeps its value inside 1500..max; a fixed range resets it to the only value
            var current = int.TryParse(_mtu.Text, out var m) ? m : (int)XenAPI.Network.MTU_DEFAULT;
            var clamped = Math.Clamp(current, (int)XenAPI.Network.MTU_MIN, max);
            if (!_mtu.IsEnabled || clamped != current)
                _mtu.Text = (_mtu.IsEnabled ? clamped : max).ToString();
            var check = NetworkRules.ValidateMtu(_mtu.Text, max);
            _mtuInfo.Text = check.Message;
            _mtuInfo.Foreground = check.Error ? Brushes.Red : Brushes.Black;
            Changed?.Invoke();
        }
    }
}
