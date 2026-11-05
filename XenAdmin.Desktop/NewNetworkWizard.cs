using System;
using System.Collections.Generic;
using System.Linq;
using Task = System.Threading.Tasks.Task;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using XenAdmin.Actions;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// NewNetworkWizard: Select Type, Details (NetWName), Interface (NetWDetails) for External and Single-Server Private
    /// networks (NetworkAction), Select Type, Details, Interface (NetWSriovDetails) for an SR-IOV network
    /// (CreateSriovAction), or Select Type, Bond Members (NetWBondDetails) for a bond (CreateBondAction), as the WinForms
    /// wizard does. Cross-Server Private networks are gone since Stockholm (XenServer 8 removed the vSwitch Controller).
    /// shortcut: the External page has no "create the VLAN on the SR-IOV network" option yet; add it when needed.
    /// </summary>
    public class NewNetworkWizard : WizardWindow
    {
        private readonly IXenConnection _c;
        private readonly Host _host;
        private readonly Action<AsyncAction> _run;
        private readonly IXenObject _owner;

        private bool _external = true, _bonded, _sriov;
        private BondEditor _bond;
        private string _nameFor;
        private string _name = "", _description = "";
        private PIF _nic;
        private string _vlan, _mtu = XenAPI.Network.MTU_DEFAULT.ToString();
        private bool _autoPlug = true;

        public NewNetworkWizard(IXenObject hostOrPool, Action<AsyncAction> run)
            : base(string.Format(Messages.NEWNETWORKWIZARD_TITLE, XenAdmin.Core.Helpers.GetName(hostOrPool.Connection)), "000_NewNetwork_h32bit_32",
                "Select Type", "Details", "Interface")
        {
            _owner = hostOrPool;
            _c = hostOrPool.Connection;
            _host = NetworkRules.WizardHost(hostOrPool);
            _run = run;
        }

        protected override string PageTitle(int step) => step switch
        {
            0 => Messages.NETW_TYPE_SELECT_TITLE,
            1 when _bonded => Messages.NETW_BOND_DETAILS_TITLE,
            2 when _sriov => Messages.NETW_INTERNAL_DETAILS_TITLE,
            1 => Messages.NETW_NAME_TITLE,
            _ => _external ? Messages.NETW_EXTERNAL_DETAILS_TITLE : Messages.NETW_INTERNAL_DETAILS_TITLE
        };

        protected override Control Page(int step)
        {
            var kind = _sriov ? "sriov" : _external ? "external" : "private";
            if (step > 0 && _nameFor != kind)
            {
                _name = _sriov ? NetworkRules.DefaultSriovName(_c) : NetworkRules.DefaultNetworkName(_c, _external);
                _nameFor = kind;
            }
            if (step == 1 && _bonded)
            {
                // NetWBondDetails: the BondDetails control is the whole page
                _bond = new BondEditor(_owner);
                _bond.Changed += () => FinishButton.IsEnabled = _bond.IsValid;
                FinishButton.IsEnabled = _bond.IsValid;
                return _bond;
            }
            return step switch { 0 => TypePage(), 1 => NamePage(), _ => _sriov ? SriovPage() : InterfacePage() };
        }

        protected override string Validate(int step)
        {
            if (step == 1 && _bonded)
                return _bond != null && _bond.IsValid ? null : "Select two or more NICs that are not already bonded.";
            if (step == 1 && string.IsNullOrWhiteSpace(_name))
                return "Enter a name for the new network.";
            if (step == 2 && _sriov && _nic == null)
                return "Select a NIC to enable SR-IOV.";
            if (step == 2 && _external && !_sriov)
            {
                if (_nic == null)
                    return "Select a NIC.";
                var vlan = NetworkRules.ValidateVlan(_vlan, NetworkRules.MinVlan(_c), NetworkRules.VlansInUse(_nic));
                if (vlan.Error)
                    return vlan.Message;
                var mtu = NetworkRules.ValidateMtu(_mtu, NetworkRules.MaxMtu(_nic));
                if (mtu.Error)
                    return mtu.Message;
            }
            return null;
        }

        protected override async Task Finish()
        {
            if (_bonded)
            {
                if (await _bond.CreateAction(this) is AsyncAction bond)
                {
                    _run(bond);
                    Close();
                }
                return;
            }
            if (_sriov)
            {
                // NetWSriovDetails.PageLeaveCore: SR-IOV disturbs the pool's connections
                var (warning, button) = NetworkRules.SriovWarning();
                if (_nic == null || !await Dialogs.Confirm(this, Messages.NETWORK_SRIOV, warning, button, Messages.CANCEL))
                    return;
                if (NetworkRules.SriovAction(_c, NetworkRules.NewNetwork(_name, _description, _autoPlug, null), _nic) is AsyncAction sriov)
                    _run(sriov);
                Close();
                return;
            }
            int? mtu = _external && int.TryParse(_mtu.Trim(), out var m) ? m : null;
            var network = NetworkRules.NewNetwork(_name, _description, _autoPlug, mtu);
            _run(_external
                ? new NetworkAction(_c, network, _nic, int.Parse(_vlan.Trim()))
                : new NetworkAction(_c, network, true));
            Close();
        }

        // ---------- pages ----------

        private static TextBlock Text(string text, double opacity = 1) =>
            new TextBlock { Text = text, TextWrapping = TextWrapping.Wrap, Opacity = opacity };

        private Control TypePage()
        {
            var panel = new StackPanel { Spacing = 6 };
            panel.Children.Add(Text("Select the type of new network you would like to create:"));
            RadioButton Option(string label, string description, bool enabled, bool selected, Action pick, string reason = null)
            {
                var rb = new RadioButton { Content = Dialogs.M(label), GroupName = "nettype", IsEnabled = enabled, IsChecked = selected, Margin = new Thickness(0, 8, 0, 0) };
                rb.IsCheckedChanged += (_, _) => { if (rb.IsChecked == true) pick(); };
                panel.Children.Add(rb);
                panel.Children.Add(new TextBlock { Text = description, TextWrapping = TextWrapping.Wrap, Margin = new Thickness(22, 0, 0, 0), Opacity = enabled ? 0.8 : 0.5 });
                if (reason != null)
                    panel.Children.Add(new TextBlock { Text = reason, TextWrapping = TextWrapping.Wrap, Margin = new Thickness(22, 2, 0, 0), Foreground = Brushes.DarkGoldenrod });
                return rb;
            }
            // UpdateWizardContent: a bond has one more page (Bond Members) instead of Details and Interface
            void Pick(bool external, bool bonded, bool sriov = false)
            {
                _external = external;
                _bonded = bonded;
                if (_sriov != sriov) _nic = null; // the NIC lists differ
                _sriov = sriov;
                SetSteps(bonded ? new[] { "Select Type", Messages.NETW_BOND_DETAILS_TEXT } : new[] { "Select Type", "Details", "Interface" });
            }
            Option("&External Network", "Create a network that passes traffic over one of your VLANs.", true, _external && !_bonded && !_sriov, () => Pick(true, false));
            Option("&Single-Server Private Network", "Create a network that does not leave each server. This can be used as a private connection between VMs on the same host.", true, !_external && !_bonded && !_sriov, () => Pick(false, false));
            Option("&Bonded Network", "Create a network that bonds together two or more of your NICs. This will create a single higher performing channel.", true, _bonded, () => Pick(false, true));
            var (sriovVisible, sriovBlocker) = NetworkRules.SriovOption(_c);
            if (sriovVisible)
                Option("S&R-IOV Network", "Enable SR-IOV on a NIC and create an SR-IOV network on that NIC.", sriovBlocker == null, _sriov,
                    () => Pick(false, false, true), sriovBlocker);
            return panel;
        }

        private Control NamePage()
        {
            var name = new TextBox { Text = _name };
            var description = new TextBox { Text = _description, AcceptsReturn = true, Height = 70, TextWrapping = TextWrapping.Wrap };
            name.TextChanged += (_, _) => _name = name.Text ?? "";
            description.TextChanged += (_, _) => _description = description.Text ?? "";
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowDefinitions = new RowDefinitions("Auto,Auto"), RowSpacing = 8, ColumnSpacing = 10 };
            var l1 = new Label { Content = Dialogs.M("Na&me:"), Target = name, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
            var l2 = new Label { Content = Dialogs.M("&Description:"), Target = description, Padding = new Thickness(0, 4, 0, 0) };
            Grid.SetColumn(name, 1); Grid.SetRow(l2, 1); Grid.SetRow(description, 1); Grid.SetColumn(description, 1);
            grid.Children.Add(l1); grid.Children.Add(name); grid.Children.Add(l2); grid.Children.Add(description);
            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(Text("Provide a name and optional description for the new network."));
            panel.Children.Add(grid);
            name.AttachedToVisualTree += (_, _) => name.Focus();
            return panel;
        }

        /// <summary>NetWSriovDetails: the NIC to enable SR-IOV on, and whether new VMs get this network.</summary>
        private Control SriovPage()
        {
            var panel = new StackPanel { Spacing = 10 };
            var nics = _host == null ? new List<PIF>() : NetworkRules.SriovNics(_host);
            if (_nic == null || nics.All(p => p.opaque_ref != _nic.opaque_ref))
                _nic = nics.FirstOrDefault();
            var nic = new ComboBox { ItemsSource = nics.Select(p => p.Name()).ToList(), SelectedIndex = nics.FindIndex(p => p.opaque_ref == _nic?.opaque_ref), Width = 220 };
            nic.SelectionChanged += (_, _) => { _nic = nic.SelectedIndex >= 0 ? nics[nic.SelectedIndex] : null; FinishButton.IsEnabled = _nic != null; };
            var auto = new CheckBox { Content = Dialogs.M("&Automatically add this network to new virtual machines."), IsChecked = _autoPlug };
            auto.IsCheckedChanged += (_, _) => _autoPlug = auto.IsChecked == true;
            var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            row.Children.Add(new Label { Content = Dialogs.M("&NIC:"), Target = nic, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) });
            row.Children.Add(nic);
            panel.Children.Add(Text("Select a NIC to enable SR-IOV."));
            panel.Children.Add(row);
            panel.Children.Add(auto);
            FinishButton.IsEnabled = _nic != null;
            return panel;
        }

        private Control InterfacePage()
        {
            var panel = new StackPanel { Spacing = 10 };
            var auto = new CheckBox { Content = Dialogs.M("&Automatically add this network to new virtual machines"), IsChecked = _autoPlug };
            auto.IsCheckedChanged += (_, _) => _autoPlug = auto.IsChecked == true;
            FinishButton.IsEnabled = true; // an earlier External page may have left it disabled
            if (!_external)
            {
                panel.Children.Add(Text("Select whether you would like to automatically add this network to new VMs."));
                panel.Children.Add(auto);
                return panel;
            }
            panel.Children.Add(Text("Your new network will be mapped to an existing physical network interface and assigned a VLAN number to use on that interface. You can select the physical interface you would like to use below."));
            var nics = _host == null ? new List<PIF>() : NetworkRules.WizardNics(_host);
            _nic ??= nics.FirstOrDefault();
            var nic = new ComboBox { ItemsSource = nics.Select(p => p.Name()).ToList(), SelectedIndex = Math.Max(0, nics.FindIndex(p => p.opaque_ref == _nic?.opaque_ref)), Width = 220 };
            var vlan = new TextBox { Text = _vlan, Width = 90 };
            var mtu = new TextBox { Text = _mtu, Width = 90 };
            var vlanInfo = new TextBlock { TextWrapping = TextWrapping.Wrap, VerticalAlignment = VerticalAlignment.Center };
            var mtuInfo = new TextBlock { TextWrapping = TextWrapping.Wrap, VerticalAlignment = VerticalAlignment.Center, Opacity = 0.8 };
            void PickVlan()
            {
                // the NumericUpDown starts at its minimum; a VLAN already used on the NIC moves to the first free one
                if (_nic == null) return;
                var current = int.TryParse(_vlan, out var v) ? v : NetworkRules.MinVlan(_c);
                _vlan = NetworkRules.VlanFor(current, NetworkRules.VlansInUse(_nic)).ToString();
                vlan.Text = _vlan;
            }
            void Update()
            {
                if (_nic == null) return;
                var v = NetworkRules.ValidateVlan(_vlan, NetworkRules.MinVlan(_c), NetworkRules.VlansInUse(_nic));
                vlanInfo.Text = v.Message ?? "";
                vlanInfo.Foreground = v.Error ? Brushes.Red : Brushes.Black;
                var max = NetworkRules.MaxMtu(_nic);
                mtu.IsEnabled = max > XenAPI.Network.MTU_MIN;
                var m = NetworkRules.ValidateMtu(_mtu, max);
                mtuInfo.Text = m.Message ?? "";
                mtuInfo.Foreground = m.Error ? Brushes.Red : Brushes.Black;
                FinishButton.IsEnabled = !v.Error && !m.Error;
            }
            void ClampMtu()
            {
                // the NumericUpDown keeps its value inside the NIC's range; a fixed range resets it to the only value
                if (_nic == null) return;
                var max = NetworkRules.MaxMtu(_nic);
                if (!int.TryParse(_mtu, out var m) || m > max || m < XenAPI.Network.MTU_MIN)
                {
                    _mtu = Math.Clamp(int.TryParse(_mtu, out m) ? m : (int)XenAPI.Network.MTU_DEFAULT, (int)XenAPI.Network.MTU_MIN, max).ToString();
                    mtu.Text = _mtu;
                }
            }
            nic.SelectionChanged += (_, _) => { _nic = nic.SelectedIndex >= 0 && nic.SelectedIndex < nics.Count ? nics[nic.SelectedIndex] : null; PickVlan(); ClampMtu(); Update(); };
            vlan.TextChanged += (_, _) => { _vlan = vlan.Text ?? ""; Update(); };
            mtu.TextChanged += (_, _) => { _mtu = mtu.Text ?? ""; Update(); };

            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,Auto,*"), RowDefinitions = new RowDefinitions("Auto,Auto,Auto"), RowSpacing = 8, ColumnSpacing = 10 };
            void Row(int r, string label, Control field, Control info)
            {
                var l = new Label { Content = Dialogs.M(label), Target = field, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
                Grid.SetRow(l, r); Grid.SetRow(field, r); Grid.SetColumn(field, 1);
                grid.Children.Add(l); grid.Children.Add(field);
                if (info != null) { Grid.SetRow(info, r); Grid.SetColumn(info, 2); grid.Children.Add(info); }
            }
            Row(0, "NI&C:", nic, null);
            Row(1, "&VLAN:", vlan, vlanInfo);
            Row(2, "&MTU:", mtu, mtuInfo);
            panel.Children.Add(grid);
            panel.Children.Add(auto);
            PickVlan();
            ClampMtu();
            Update();
            return panel;
        }
    }
}
