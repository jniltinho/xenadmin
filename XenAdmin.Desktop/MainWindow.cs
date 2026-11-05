using System.Globalization;
using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Controls.Templates;
using Avalonia.Layout;
using Avalonia.Media;
using Avalonia.Styling;
using Avalonia.Threading;
using XenAdmin.Actions;
using XenAdmin.Core;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>Main window laid out like the WinForms one: menus, toolbar, navigation pane, title band and tabs.</summary>
    public class MainWindow : Window
    {
        private readonly AppConfig _config;
        private readonly CertificateTrust _trust;
        private readonly List<XenConnection> _connections = new List<XenConnection>();
        private readonly TreeView _tree = new TreeView();
        // NavigationPane: the view the tree shows (and the saved search of Saved Searches)
        private NavView _navView = NavView.Infrastructure;
        private XenAdmin.XenSearch.Search _savedSearch;
        private readonly List<Button> _navButtons = new List<Button>();
        private Button _savedButton, _infraButton;
        private readonly TextBox _search = new TextBox { Watermark = "Search...", BorderThickness = new Thickness(0) };
        private readonly ContentControl _titleBand = new ContentControl();
        private readonly TabControl _tabs = new TabControl();
        private readonly ObservableCollection<string> _events = new ObservableCollection<string>();
        private readonly TextBlock _status = new TextBlock { Margin = new Thickness(6, 2) };
        private readonly DispatcherTimer _refreshTimer;
        private readonly List<string> _history = new List<string>();
        private int _historyIndex = -1;
        private bool _showEvents;
        private bool _showAlerts = true; // Notifications: Alerts or Events
        private Border _leftPane;
        private readonly StackPanel _statusAlerts = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 4 };
        private NotificationsView _alertsView;
        private ListBox _notificationList;
        private readonly Border _alertBadge = new Border
        {
            Background = Brush.Parse("#C00000"), CornerRadius = new CornerRadius(9), MinWidth = 18, Height = 18, Padding = new Thickness(5, 0), IsVisible = false,
            VerticalAlignment = VerticalAlignment.Center, Child = new TextBlock { Foreground = Brushes.White, FontWeight = FontWeight.Bold, FontSize = 11, HorizontalAlignment = HorizontalAlignment.Center, VerticalAlignment = VerticalAlignment.Center }
        };
        private string _selectedKey = "<root>";
        private string _treeSignature;
        private string _tabsSignature;
        private TreeNode _rootNode;
        // last tab per object, as MainWindow.selectedTabs (RememberLastSelectedTab)
        private readonly Dictionary<string, string> _lastTab = new Dictionary<string, string>();
        private string _lastTabAny; // selectedOverviewTab: the last tab of any object, when tabs are not remembered per object
        private bool _refreshing;

        // toolbar / menu items whose state follows the selection
        private readonly Button _back, _forward, _addServer, _newPool, _newStorage, _newVm, _pause;
        private readonly Dictionary<VmOp, Button> _power = new Dictionary<VmOp, Button>();
        private readonly Dictionary<VmOp, MenuItem> _vmMenu = new Dictionary<VmOp, MenuItem>();
        private MenuItem _connectMenu, _disconnectMenu, _removeMenu, _poolProps, _hostProps, _hostReboot, _hostShutdown, _hostMaintenance, _vmProps, _newVmMenu, _newPoolMenu, _newSrMenu, _migrateMenu, _copyVmMenu, _deleteVmMenu,
            _restartToolstack, _removeFromPool, _snapshotMenu, _toTemplateMenu, _fromTemplateMenu, _copyTemplateMenu, _deleteTemplateMenu, _templateProps,
            _poolDisconnect, _poolReconnectAs, _rotateSecret, _poolPassword, _hostPassword, _poolAddServer, _poolRemoveServer, _makeStandalone, _haConfigure, _haDisable, _exportVmMenu, _vappsMenu, _assignVappMenu, _vtpmMenu, _drMenu;
        // connections with an Enable/Disable HA action of ours running (HelpersGUI.FindActiveHaAction)
        private readonly HashSet<IXenConnection> _haBusy = new HashSet<IXenConnection>();

        public MainWindow(AppConfig config, CertificateTrust trust)
        {
            // the Search tab's New/Edit Search, Saved Searches and Import show the search in the tree
            SearchView.ShowSearch = search =>
            {
                if (search == null) ShowView(NavView.Infrastructure, null, _infraButton);
                else ShowView(NavView.SavedSearch, search, _savedButton);
            };
            SearchView.Run = RunTracked;
            _config = config;
            _trust = trust;
            ConsoleView.RunAction = RunTracked;
            ConsoleView.Changed = () => { _tabsSignature = null; UpdateContent(); };
            Title = "XCP-ng Center";
            Width = 1280;
            Height = 800;
            WindowStartupLocation = WindowStartupLocation.CenterScreen;
            // a screen smaller than the default size (1280x800 VMs, small laptops) gets the window maximized
            Opened += (_, _) =>
            {
                if (Screens.ScreenFromWindow(this)?.WorkingArea is PixelRect area
                    && (area.Width / RenderScaling < Width || area.Height / RenderScaling < Height))
                    WindowState = WindowState.Maximized;
            };
            FontSize = 12;
            Background = Brush.Parse("#F0F0F0");
            try { Icon = new WindowIcon(Avalonia.Platform.AssetLoader.Open(new Uri("avares://xcp-ng-center/Assets/xcp-ng-center.png"))); } catch { }

            _back = Tool("001_Back_h32bit_24", "Back", () => Navigate(-1));
            _forward = Tool("001_Forward_h32bit_24", "Forward", () => Navigate(1));
            _addServer = Tool("000_AddApplicationServer_h32bit_24", "Add New Server", async () => await AddServer());
            _newPool = Tool("000_PoolNew_h32bit_24", "New Pool", () => ShowNewPool());
            _newStorage = Tool("000_NewStorage_h32bit_24", "New Storage", () => ShowNewSr());
            _newVm = Tool("000_CreateVM_h32bit_24", "New VM", () => ShowNewVm());
            _power[VmOp.Start] = Tool("001_PowerOn_h32bit_24", "Start", () => RunVmOp(VmOp.Start));
            _power[VmOp.Shutdown] = Tool("001_ShutDown_h32bit_24", "Shut Down", () => { if (SelectedNode()?.XenObject is Host h) HostPower(h, false); else RunVmOp(VmOp.Shutdown); });
            _power[VmOp.ForceShutdown] = Tool("001_ForceShutDown_h32bit_24", "Force Shutdown", () => RunVmOp(VmOp.ForceShutdown));
            _power[VmOp.Reboot] = Tool("001_Reboot_h32bit_24", "Reboot", () => { if (SelectedNode()?.XenObject is Host h) HostPower(h, true); else RunVmOp(VmOp.Reboot); });
            _power[VmOp.ForceReboot] = Tool("001_ForceReboot_h32bit_24", "Force Reboot", () => RunVmOp(VmOp.ForceReboot));
            _power[VmOp.Suspend] = Tool("000_Paused_h32bit_24", "Suspend", () => RunVmOp(VmOp.Suspend));
            _power[VmOp.Resume] = Tool("000_Resumed_h32bit_24", "Resume", () => RunVmOp(VmOp.Resume));
            _pause = Tool("000_Paused_h32bit_24_green", "Pause", null);

            var toolbar = new StackPanel { Orientation = Orientation.Horizontal, Height = 30, Margin = new Thickness(2, 0) };
            void Sep() => toolbar.Children.Add(new Border { Width = 1, Background = Brush.Parse("#C5C5C5"), Margin = new Thickness(4, 5) });
            toolbar.Children.Add(_back);
            toolbar.Children.Add(_forward);
            Sep();
            toolbar.Children.Add(_addServer);
            Sep();
            foreach (var b in new[] { _newPool, _newStorage, _newVm }) toolbar.Children.Add(b);
            Sep();
            foreach (var op in new[] { VmOp.Start, VmOp.Shutdown, VmOp.ForceShutdown, VmOp.Reboot, VmOp.ForceReboot, VmOp.Suspend, VmOp.Resume })
                toolbar.Children.Add(_power[op]);
            toolbar.Children.Add(_pause);

            var top = new StackPanel();
            top.Children.Add(BuildMenu());
            top.Children.Add(toolbar);

            // navigation pane: search box, tree, view buttons
            _tree.ItemTemplate = new FuncTreeDataTemplate<TreeNode>((n, _) =>
            {
                var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 4 };
                row.Children.Add(Views.Icon(n.Icon));
                row.Children.Add(new TextBlock { Text = n.Name, VerticalAlignment = VerticalAlignment.Center });
                return row;
            }, n => n.Children);
            _tree.Classes.Add("nav");
            _tree.Styles.Add(new Avalonia.Styling.Style(x => x.OfType<TreeViewItem>())
            {
                Setters = { new Avalonia.Styling.Setter(TreeViewItem.IsExpandedProperty, new Avalonia.Data.Binding(nameof(TreeNode.IsExpanded)) { Mode = Avalonia.Data.BindingMode.TwoWay }) }
            });
            _tree.SelectionChanged += (_, _) =>
            {
                // resetting ItemsSource raises a null selection; keep the user's choice
                if (_refreshing || !(_tree.SelectedItem is TreeNode n)) return;
                Select(n.Key, true);
            };
            _search.TextChanged += (_, _) => RefreshTree();
            _tree.DoubleTapped += async (_, e) =>
            {
                if ((e.Source as Control)?.DataContext is TreeNode n)
                    await DoubleClick(n);
            };
            _tree.ContextRequested += (_, e) =>
            {
                if (!((e.Source as Control)?.DataContext is TreeNode n)) return;
                if (n.Key != _selectedKey) { _refreshing = true; _tree.SelectedItem = n; _refreshing = false; Select(n.Key, true); }
                var items = ContextMenuFor(n);
                if (items.Count == 0) return;
                new ContextMenu { ItemsSource = items }.Open((Control)e.Source);
                e.Handled = true;
            };

            var searchRow = new DockPanel { Background = Brushes.White };
            var magnifier = Views.Icon("000_Search_h32bit_16");
            magnifier.Margin = new Thickness(4, 0);
            DockPanel.SetDock(magnifier, Dock.Right);
            searchRow.Children.Add(magnifier);
            searchRow.Children.Add(_search);

            var navButtons = new StackPanel();
            Button infra = null, objects = null, orgViews = null, saved = null;
            infra = (Button)NavButton("infra_view_24", "Infrastructure", true, () => ShowView(NavView.Infrastructure, null, infra));
            infra.Classes.Add("active");
            navButtons.Children.Add(infra);
            objects = (Button)NavButton("objects_24", "Objects", true, () => ShowView(NavView.Objects, null, objects));
            navButtons.Children.Add(objects);
            // NavigationPane: Organization Views and Saved Searches open a menu of views / searches
            orgViews = (Button)NavButton("org_view_24", "Organization Views", true, null);
            orgViews.Click += (_, _) =>
            {
                var menu = new ContextMenu();
                foreach (var (view, text) in new[] { (NavView.Tags, Messages.TAGS), (NavView.Folders, Messages.FOLDERS), (NavView.CustomFields, Messages.CUSTOM_FIELDS), (NavView.Vapps, Messages.VM_APPLIANCES) })
                {
                    var item = new MenuItem { Header = text, Icon = Views.Icon(NavigationViews.RootOf(view).Icon) };
                    item.Click += (_, _) => ShowView(view, null, orgViews);
                    menu.Items.Add(item);
                }
                menu.Open(orgViews);
            };
            navButtons.Children.Add(orgViews);
            saved = (Button)NavButton("saved_searches_24", "Saved Searches", true, null);
            _savedButton = saved;
            _infraButton = infra;
            saved.Click += (_, _) =>
            {
                var menu = new ContextMenu();
                foreach (var s in NavigationViews.SavedSearches())
                {
                    var item = new MenuItem { Header = s.Name.Replace("_", "__"), Icon = Views.Icon(s.DefaultSearch ? "000_defaultSpyglass_h32bit_16" : "000_Search_h32bit_16") };
                    item.Click += (_, _) => ShowView(NavView.SavedSearch, s, saved);
                    menu.Items.Add(item);
                }
                menu.Open(saved);
            };
            navButtons.Children.Add(saved);
            _navButtons.AddRange(new[] { infra, objects, orgViews, saved });
            // NavigationPane: the Notifications button carries the alert count in a red badge
            var notifications = NavButton("notif_none_24", "Notifications", true, () => { _showEvents = true; UpdateContent(); });
            ((StackPanel)((Button)notifications).Content).Children.Add(_alertBadge);
            navButtons.Children.Add(notifications);

            var left = new DockPanel();
            DockPanel.SetDock(searchRow, Dock.Top);
            DockPanel.SetDock(navButtons, Dock.Bottom);
            left.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = searchRow, Margin = new Thickness(2, 2, 2, 2) });
            DockPanel.SetDock(left.Children[0], Dock.Top);
            left.Children.Add(navButtons);
            _leftPane = new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Background = Brushes.White, Margin = new Thickness(2, 0, 2, 4), Child = _tree };
            left.Children.Add(_leftPane);

            var right = new DockPanel { Margin = new Thickness(2, 2, 2, 2) };
            DockPanel.SetDock(_titleBand, Dock.Top);
            right.Children.Add(_titleBand);
            right.Children.Add(_tabs);
            _tabs.SelectionChanged += (_, _) =>
            {
                if (_tabs.SelectedItem is TabItem t && SelectedNode()?.XenObject is IXenObject x && _tabsSignature?.StartsWith(x.opaque_ref + "|") == true)
                {
                    _lastTab[x.opaque_ref] = (string)t.Header;
                    _lastTabAny = (string)t.Header;
                }
            };

            var split = new Grid { ColumnDefinitions = new ColumnDefinitions("245,3,*") };
            Grid.SetColumn(left, 0);
            var splitter = new GridSplitter { ResizeDirection = GridResizeDirection.Columns, Background = Brushes.Transparent };
            Grid.SetColumn(splitter, 1);
            Grid.SetColumn(right, 2);
            split.Children.Add(left);
            split.Children.Add(splitter);
            split.Children.Add(right);

            var root = new DockPanel();
            DockPanel.SetDock(top, Dock.Top);
            // statusButtonAlerts: alert icon and count at the right, opens the alerts
            var alertsButton = new Button { Content = _statusAlerts, Background = Brushes.Transparent, BorderThickness = new Thickness(0), Padding = new Thickness(6, 0), IsVisible = false };
            alertsButton.Click += (_, _) => { _showEvents = true; _showAlerts = true; UpdateContent(); };
            _statusAlerts.Children.Add(Views.Icon("000_Alert2_h32bit_16", 16));
            _statusAlerts.Children.Add(new TextBlock { VerticalAlignment = VerticalAlignment.Center });
            var statusRow = new DockPanel();
            DockPanel.SetDock(alertsButton, Dock.Right);
            statusRow.Children.Add(alertsButton);
            statusRow.Children.Add(_status);
            var statusBar = new Border { BorderBrush = Brush.Parse("#D7D7D7"), BorderThickness = new Thickness(0, 1, 0, 0), Child = statusRow };
            DockPanel.SetDock(statusBar, Dock.Bottom);
            root.Children.Add(top);
            root.Children.Add(statusBar);
            root.Children.Add(split);
            Content = root;

            // coalesce bursts of cache events into one refresh
            _refreshTimer = new DispatcherTimer { Interval = TimeSpan.FromMilliseconds(500) };
            _refreshTimer.Tick += (_, _) => { _refreshTimer.Stop(); RefreshTree(); };

            ConsoleView.CertificatePending += async console =>
            {
                if (await ResolvePendingCertificates())
                    console.Retry();
            };

            foreach (var s in _config.Servers)
                Track(new XenConnection { Hostname = s.Host, Port = s.Port, Username = s.User, FriendlyName = s.Name });
            Select("<root>", true);
            RefreshTree();
            Closing += (_, _) =>
            {
                ((DesktopConfigProvider)XenAdminConfigManager.Provider).Exiting = true;
                ConsoleView.CloseAll();
                foreach (var c in _connections.Where(c => c.IsConnected))
                    c.EndConnect(true, true);
            };
        }

        private static Button Tool(string icon, string text, Action onClick)
        {
            var content = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 4 };
            var image = Views.Icon(icon, 22);
            content.Children.Add(image);
            content.Children.Add(new TextBlock { Text = text, VerticalAlignment = VerticalAlignment.Center });
            var b = new Button { Content = content, IsEnabled = onClick != null };
            b.Classes.Add("tool");
            // ToolStrip draws disabled images in grey
            void Paint() => image.Source = b.IsEnabled ? Images.Get(icon) : Images.Gray(icon);
            b.PropertyChanged += (_, e) => { if (e.Property == IsEnabledProperty) Paint(); };
            Paint();
            if (onClick != null)
                b.Click += (_, _) => onClick();
            return b;
        }

        private static Control NavButton(string icon, string text, bool enabled, Action onClick)
        {
            var content = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
            content.Children.Add(Views.Icon(icon, 22));
            content.Children.Add(new TextBlock { Text = text, VerticalAlignment = VerticalAlignment.Center });
            var b = new Button { Content = content, IsEnabled = enabled, HorizontalAlignment = HorizontalAlignment.Stretch, HorizontalContentAlignment = HorizontalAlignment.Left };
            b.Classes.Add("nav");
            if (onClick != null)
                b.Click += (_, _) => onClick();
            return b;
        }

        private Menu BuildMenu()
        {
            MenuItem Item(string header, Action onClick = null, string gesture = null, string icon = null)
            {
                var m = new MenuItem { Header = header, IsEnabled = onClick != null };
                if (icon != null) m.Icon = Views.Icon(icon, 16); // the WinForms command's MenuImage
                if (gesture != null) m.InputGesture = Avalonia.Input.KeyGesture.Parse(gesture);
                if (onClick != null) m.Click += (_, _) => onClick();
                return m;
            }
            MenuItem Top(string header, params object[] items)
            {
                var m = new MenuItem { Header = header };
                foreach (var i in items) m.Items.Add(i is string ? new Separator() : i);
                return m;
            }
            _connectMenu = Item("_Connect", async () => { if (SelectedConnection() is XenConnection c) await Connect(c); });
            _disconnectMenu = Item("_Disconnect", () => { if (SelectedConnection() is XenConnection c) Disconnect(c); });
            _removeMenu = Item("_Remove from XCP-ng Center", RemoveSelected);
            foreach (var (op, text, icon) in new[] { (VmOp.Start, "S_tart", "001_PowerOn_h32bit_16"), (VmOp.Shutdown, "_Shut Down", "001_ShutDown_h32bit_16"),
                (VmOp.Reboot, "_Reboot", "001_Reboot_h32bit_16"), (VmOp.Suspend, "S_uspend", "000_paused_h32bit_16"), (VmOp.Resume, "Resu_me", "000_Resumed_h32bit_16"),
                (VmOp.ForceShutdown, "_Force Shutdown", "001_ForceShutDown_h32bit_16"), (VmOp.ForceReboot, "Force Re_boot", "001_ForceReboot_h32bit_16") })
            {
                var captured = op;
                _vmMenu[op] = Item(text, () => RunVmOp(captured), icon: icon);
            }
            return new Menu
            {
                Background = Brush.Parse("#F5F5F5"),
                ItemsSource = new[]
                {
                    Top("_File", Item("_Add...", async () => await AddServer(), icon: "000_AddApplicationServer_h32bit_16"), "-",
                        Item("_Import...", ImportVm, icon: "000_ImportVirtualAppliance_h32bit_16"), "-",
                        Item("E_xit", Close, "Alt+F4")),
                    Top("_View", Item("_Infrastructure", () => { _showEvents = false; UpdateContent(); }),
                        Item("_Notifications", () => { _showEvents = true; UpdateContent(); })),
                    Top("_Pool", _newPoolMenu = Item("New Pool...", () => ShowNewPool(), icon: "000_PoolNew_h32bit_16"), "-",
                        _poolAddServer = PoolServersMenu("_Add Server", true), _poolRemoveServer = PoolServersMenu("Re_move Server", false),
                        _poolReconnectAs = Item("Reconnect _As...", async () => { if (SelectedConnection() is XenConnection c) await ReconnectAs(c); }),
                        _poolDisconnect = Item("_Disconnect", () => { if (SelectedConnection() is XenConnection c) Disconnect(c); }), "-",
                        _poolPassword = Item("Change Server Pass_word...", () => ChangePassword(SelectedNode()?.XenObject as Pool)),
                        _rotateSecret = Item("Rotate Pool _Secret", () => RotateSecret(SelectedNode()?.Connection)), "-",
                        Top(Dialogs.M(Messages.MAINWINDOW_HIGH_AVAILABILITY),
                            _haConfigure = Item("_Configure...", () => ConfigureHa(SelectedNode()?.Connection), icon: "ha_16"),
                            _haDisable = Item("_Disable", () => DisableHa(SelectedNode()?.Connection))), "-",
                        _drMenu = Top(Dialogs.M(Messages.DISASTER_RECOVERY_AMP),
                            Item("_Configure...", () => ConfigureDr(SelectedNode()?.Connection)),
                            Item("_Disaster Recovery Wizard...", () => { if (SelectedNode()?.Connection is IXenConnection drc && DrRules.CanRun(drc)) new DrWizard(drc, RunTracked).Show(this); })),
                        _vappsMenu = Item(Dialogs.M(Messages.VM_APPLIANCES_MENU), () => ManageVapps(SelectedNode()?.Connection), icon: "000_VirtualAppliance_h32bit_16"), "-",
                        _makeStandalone = Item(Dialogs.M(Messages.MAINWINDOW_DELETE_POOL), () => MakeStandalone(SelectedNode()?.Connection)), "-",
                        _poolProps = Item("P_roperties", () => ShowProperties(SelectedNode()?.XenObject as Pool), icon: "edit_16")),
                    Top("_Server", Item("_Add...", async () => await AddServer(), icon: "000_AddApplicationServer_h32bit_16"), "-", _connectMenu, _disconnectMenu, _removeMenu, "-",
                        _hostReboot = Item("_Reboot", () => { if (SelectedNode()?.XenObject is Host h) HostPower(h, true); }, icon: "001_Reboot_h32bit_16"),
                        _hostShutdown = Item("S_hut Down", () => { if (SelectedNode()?.XenObject is Host h) HostPower(h, false); }, icon: "001_ShutDown_h32bit_16"),
                        _restartToolstack = Item(Dialogs.M(Messages.MAINWINDOW_RESTART_TOOLSTACK), () => { if (SelectedNode()?.XenObject is Host h) RestartToolstack(h); }),
                        "-", _removeFromPool = Item("Remove From _Pool", () => { if (SelectedNode()?.XenObject is Host h) RemoveFromPool(h); }),
                        _hostMaintenance = Item(Dialogs.M(Messages.ENTER_MAINTENANCE_MODE_CONTEXT_MENU), () => { if (SelectedNode()?.XenObject is Host h) ToggleMaintenance(h); }, icon: "000_ServerMaintenance_h32bit_16"), "-",
                        _hostPassword = Item("Change Server Pass_word...", () => ChangePassword(SelectedNode()?.XenObject as Host)), "-",
                        _hostProps = Item("P_roperties", () => ShowProperties(SelectedNode()?.XenObject as Host), icon: "edit_16")),
                    Top("V_M", _newVmMenu = Item("_New VM...", () => ShowNewVm(), icon: "001_CreateVM_h32bit_16"), "-", _vmMenu[VmOp.Start], _vmMenu[VmOp.Shutdown], _vmMenu[VmOp.Reboot],
                        _vmMenu[VmOp.Suspend], _vmMenu[VmOp.Resume], "-", _vmMenu[VmOp.ForceShutdown], _vmMenu[VmOp.ForceReboot], "-",
                        _migrateMenu = MigrateMenu(() => SelectedNode()?.XenObject as VM), "-",
                        _copyVmMenu = Item("Cop_y VM...", () => CopyVm(SelectedNode()?.XenObject as VM)),
                        _snapshotMenu = Item("Ta_ke a Snapshot...", () => TakeSnapshot(SelectedNode()?.XenObject as VM), icon: "000_VMSnapShotDiskOnly_h32bit_16"),
                        _toTemplateMenu = Item("Con_vert to Template...", () => ConvertToTemplate(SelectedNode()?.XenObject as VM)),
                        _assignVappMenu = AssignVappMenu(),
                        _vtpmMenu = Item(Dialogs.M(Messages.COMMAND_VTPM_MENU), () => { if (SelectedNode()?.XenObject is VM v && VtpmRules.CanManage(v)) new VtpmDialog(v, RunTracked).Show(this); }), "-",
                        _exportVmMenu = Item("_Export...", async () => { if (SelectedNode()?.XenObject is VM v && await Dialogs.ExportVm(this, v) is AsyncAction a) RunTracked(a); }, icon: "000_ExportVirtualAppliance_h32bit_16"),
                        _deleteVmMenu = Item("_Delete VM...", () => DeleteVm(SelectedNode()?.XenObject as VM)), "-",
                        _vmProps = Item("P_roperties", () => ShowProperties(SelectedNode()?.XenObject as VM), icon: "edit_16")),
                    Top("S_torage", _newSrMenu = Item("New SR...", () => ShowNewSr(), icon: "000_NewStorage_h32bit_16")),
                    Top("_Templates", _fromTemplateMenu = Item("Create _VM From Selection...", () => { if (SelectedNode()?.XenObject is VM t) ShowNewVm(t); }, icon: "000_CreateVM_h32bit_16"), "-",
                        new MenuItem { Header = "_Export to File...", IsEnabled = false }, // ExportTemplateCommand: export is not ported yet
                        _copyTemplateMenu = Item("_Copy...", () => CopyVm(SelectedNode()?.XenObject as VM)), "-",
                        _deleteTemplateMenu = Item("_Delete Template...", () => DeleteVm(SelectedNode()?.XenObject as VM)), "-",
                        _templateProps = Item("P_roperties", () => ShowProperties(SelectedNode()?.XenObject as VM), icon: "edit_16")),
                    Top("T_ools", Item("Server Status Report..."), "-",
                        Item("_Options...", async () => { if (await OptionsDialog.Show(this, _config.Options)) { _config.Save(); ((DesktopConfigProvider)XenAdminConfigManager.Provider).ApplyProxy(); _tabsSignature = null; UpdateContent(); } })),
                    Top("_Help", Item("_About XCP-ng Center", async () => await Dialogs.Info(this, "About XCP-ng Center",
                        $"XCP-ng Center for {Platform.PlatformServices.Current.Name} {typeof(MainWindow).Assembly.GetName().Version}\nPreview build based on the XCP-ng Center core (XenModel).")))
                }
            };
        }

        // ---------- connections ----------

        private void Track(XenConnection c)
        {
            _connections.Add(c);
            lock (ConnectionsManager.ConnectionsLock)
                ConnectionsManager.XenConnections.Add(c);
            void Schedule() => Dispatcher.UIThread.Post(() => { if (!_refreshTimer.IsEnabled) _refreshTimer.Start(); });
            c.CachePopulated += _ => Schedule();
            c.ConnectionStateChanged += _ => Schedule();
            c.ConnectionLost += _ => Dispatcher.UIThread.Post(() => AddEvent($"Connection to {c.Hostname} lost; reconnecting..."));
            c.XenObjectsUpdated += (_, _) => Schedule();
            c.ConnectionMessageChanged += (_, m) => Dispatcher.UIThread.Post(() => SetStatus($"{c.Hostname}: {m}"));
            c.ConnectionResult += (_, e) => Dispatcher.UIThread.Post(async () =>
            {
                Schedule();
                EndConnecting(c);
                if (e.Connected)
                {
                    SetStatus($"Connected to {c.Hostname}");
                    _config.AddOrUpdateServer(c.Hostname, c.Port, c.Username, c.Name);
                }
                else if (await ResolvePendingCertificates())
                {
                    StartConnect(c);
                }
                else
                {
                    c.Password = null; // ask again next time
                    AddEvent($"Failed to connect to {c.Hostname}: {e.Reason ?? e.Error?.Message}");
                    await Dialogs.Info(this, "Connection failed", e.Reason ?? e.Error?.Message ?? "Unknown error");
                }
            });
        }

        private readonly Dictionary<XenConnection, Window> _connecting = new Dictionary<XenConnection, Window>();

        /// <summary>Starts a connection with the "Connecting to..." dialog, like XenConnectionUI.BeginConnect.</summary>
        private void StartConnect(XenConnection c)
        {
            if (_connecting.TryGetValue(c, out var existing)) { existing.Activate(); return; }
            // saved servers too: never send the password over plain HTTP
            if (ServerAddress.IsPlainHttp(c.Port)) { SetStatus($"{c.Hostname}:{c.Port} is plain HTTP; connect over HTTPS (443)."); return; }
            SetStatus(string.Format(Messages.LABEL_ATTEMPT, c.Hostname));
            _connecting[c] = Dialogs.Connecting(this, c.Hostname, () => { _connecting.Remove(c); c.EndConnect(); RefreshTree(); });
            c.BeginConnect(false, PromptForNewPassword);
            RefreshTree();
        }

        /// <summary>
        /// XenConnectionUI.PromptForNewPassword: a login failed, so ask again (the "incorrect password" variant of the
        /// connect dialog). Runs on the connection's worker thread and waits for the user. Cancel ends the attempt; the
        /// connection then closes without a ConnectionResult, so the "Connecting to..." dialog is closed here.
        /// </summary>
        private bool PromptForNewPassword(IXenConnection connection, string oldPassword)
        {
            if (Dispatcher.UIThread.CheckAccess() || !(connection is XenConnection c))
                return false; // never block the UI thread waiting for itself
            var answer = new System.Threading.Tasks.TaskCompletionSource<Dialogs.ServerInput>();
            Dispatcher.UIThread.Post(async () =>
            {
                try
                {
                    // the user may have cancelled "Connecting to..." meanwhile: then there is nothing to ask for
                    if (!_connecting.TryGetValue(c, out var owner))
                    {
                        answer.TrySetResult(null);
                        return;
                    }
                    var input = await Dialogs.AddServer(owner, c.Port == 443 ? c.Hostname : $"{c.Hostname}:{c.Port}", c.Username ?? "root", true, true);
                    answer.TrySetResult(_connecting.ContainsKey(c) ? input : null);
                }
                catch (Exception e)
                {
                    answer.TrySetException(e);
                }
            });
            Dialogs.ServerInput input;
            try { input = answer.Task.GetAwaiter().GetResult(); }
            catch { input = null; }
            if (input != null)
            {
                c.Username = input.User;
                c.Password = input.Password;
                return true;
            }
            Dispatcher.UIThread.Post(() =>
            {
                EndConnecting(c);
                c.Password = null; // ask again next time
                AddEvent($"Failed to connect to {c.Hostname}: {Messages.ADD_NEW_INCORRECT}");
                SetStatus(string.Format(Messages.CONNECTION_FAILED_TITLE, c.HostnameWithPort));
                RefreshTree();
            });
            return false;
        }

        private void EndConnecting(XenConnection c)
        {
            if (_connecting.Remove(c, out var w)) w.Close();
        }

        private async System.Threading.Tasks.Task AddServer()
        {
            var input = await Dialogs.AddServer(this);
            if (input == null) return;
            var c = _connections.FirstOrDefault(x => string.Equals(x.Hostname, input.Host, StringComparison.OrdinalIgnoreCase) && x.Port == input.Port);
            if (c == null)
            {
                c = new XenConnection { Hostname = input.Host, Port = input.Port };
                Track(c);
            }
            if (c.IsConnected) return;
            c.Username = input.User;
            c.Password = input.Password;
            StartConnect(c);
        }

        private async System.Threading.Tasks.Task Connect(XenConnection c)
        {
            if (c.IsConnected || c.InProgress) return;
            if (string.IsNullOrEmpty(c.Password))
            {
                var input = await Dialogs.AddServer(this, c.Port == 443 ? c.Hostname : $"{c.Hostname}:{c.Port}", c.Username ?? "root", true);
                if (input == null) return;
                c.Username = input.User;
                c.Password = input.Password;
            }
            StartConnect(c);
        }

        private void Disconnect(XenConnection c)
        {
            if (c.IsConnected || c.InProgress)
                c.EndConnect();
            SetStatus($"Disconnected from {c.Hostname}");
            RefreshTree();
        }

        private void RemoveSelected()
        {
            if (!(SelectedConnection() is XenConnection c)) return;
            Disconnect(c);
            _connections.Remove(c);
            lock (ConnectionsManager.ConnectionsLock)
                ConnectionsManager.XenConnections.Remove(c);
            _config.RemoveServer(c.Hostname, c.Port);
            Select("<root>", true);
            RefreshTree();
        }

        /// <summary>NavigationPane.SwitchToNavigationMode: the tree shows this view; the selection starts at its root.</summary>
        private void ShowView(NavView view, XenAdmin.XenSearch.Search saved, Button button)
        {
            _navView = view;
            _savedSearch = saved;
            _showEvents = false;
            foreach (var b in _navButtons) b.Classes.Remove("active");
            button.Classes.Add("active");
            _treeSignature = null;
            _selectedKey = "<root>";
            _tabsSignature = null;
            // node keys belong to one view: Back and Forward start again in the new one
            _history.Clear();
            _historyIndex = -1;
            RefreshTree();
            Select("<root>", true);
        }

        private TreeNode SelectedNode() =>
            _rootNode == null ? null : ResourceTree.Flatten(_rootNode).FirstOrDefault(n => n.Key == _selectedKey);

        private XenConnection SelectedConnection() => SelectedNode()?.Connection as XenConnection;

        // ---------- selection, history, refresh ----------

        private void Select(string key, bool record)
        {
            _selectedKey = key;
            if (record && (_historyIndex < 0 || _history[_historyIndex] != key))
            {
                _history.RemoveRange(_historyIndex + 1, _history.Count - _historyIndex - 1);
                _history.Add(key);
                _historyIndex = _history.Count - 1;
            }
            _showEvents = false;
            _tabsSignature = null;
            UpdateContent();
        }

        private void Navigate(int delta)
        {
            var i = _historyIndex + delta;
            if (i < 0 || i >= _history.Count) return;
            _historyIndex = i;
            Select(_history[i], false);
            var match = SelectedNode();
            if (match != null) { _refreshing = true; _tree.SelectedItem = match; _refreshing = false; }
        }

        private void RefreshTree()
        {
            PruneDroppedConnections();
            var root = ResourceTree.Filter(_navView == NavView.Infrastructure
                ? ResourceTree.Root(_connections)
                : NavigationViews.Build(_navView, _savedSearch), _search.Text);
            // cache objects are updated in place, so only rebuild the tree when its shape, names or icons change
            var signature = string.Join("\n", ResourceTree.Flatten(root).Select(n => $"{n.Key}|{n.Name}|{n.Icon}|{n.Children.Count}"));
            if (signature != _treeSignature)
            {
                _treeSignature = signature;
                _rootNode = root;
                _refreshing = true;
                try { _tree.ItemsSource = new[] { root }; }
                finally { _refreshing = false; }
                var match = SelectedNode();
                if (match != null)
                    Dispatcher.UIThread.Post(() => { _refreshing = true; _tree.SelectedItem = match; _refreshing = false; }, DispatcherPriority.Loaded);
            }
            UpdateContent();
        }

        private void UpdateContent()
        {
            var node = SelectedNode();
            var o = node?.XenObject;
            var vm = o as VM;

            _back.IsEnabled = _historyIndex > 0;
            _forward.IsEnabled = _historyIndex < _history.Count - 1;
            // same rule as the WinForms toolbar: clean operation if allowed, otherwise the forced one
            bool Allowed(VmOp op) => VmOperations.IsAllowed(vm, op);
            _power[VmOp.Start].IsVisible = vm != null && vm.power_state == vm_power_state.Halted;
            _power[VmOp.Shutdown].IsVisible = !_power[VmOp.Start].IsVisible && (vm == null || Allowed(VmOp.Shutdown) || !Allowed(VmOp.ForceShutdown));
            _power[VmOp.ForceShutdown].IsVisible = !_power[VmOp.Start].IsVisible && !_power[VmOp.Shutdown].IsVisible;
            _power[VmOp.ForceReboot].IsVisible = vm != null && !Allowed(VmOp.Reboot) && Allowed(VmOp.ForceReboot);
            _power[VmOp.Reboot].IsVisible = !_power[VmOp.ForceReboot].IsVisible;
            _power[VmOp.Resume].IsVisible = vm != null && vm.power_state == vm_power_state.Suspended;
            _power[VmOp.Suspend].IsVisible = !_power[VmOp.Resume].IsVisible;
            foreach (var (op, b) in _power)
                b.IsEnabled = Allowed(op);
            // menus and buttons disabled with a reason show it (Avalonia hides tooltips of disabled controls by default)
            foreach (var reasoned in new Control[] { _power[VmOp.Start], _power[VmOp.Resume], _vmMenu[VmOp.Start], _vmMenu[VmOp.Resume], _exportVmMenu, _vtpmMenu })
                ToolTip.SetShowOnDisabled(reasoned, true);
            // 4.1c: Start/Resume of an HA-protected VM that cannot be agile is disabled early, with the reason
            var haBlocker = VmOperations.HaStartBlocker(vm);
            foreach (var op in new[] { VmOp.Start, VmOp.Resume })
            {
                if (haBlocker != null) _power[op].IsEnabled = false;
                ToolTip.SetTip(_power[op], haBlocker);
            }
            // a selected server: the toolbar's Shut Down and Reboot are the server commands
            if (o is Host selectedHost)
            {
                _power[VmOp.Shutdown].IsVisible = _power[VmOp.Reboot].IsVisible = true;
                _power[VmOp.ForceShutdown].IsVisible = _power[VmOp.ForceReboot].IsVisible = false;
                _power[VmOp.Shutdown].IsEnabled = HostRules.CanShutDown(selectedHost, Busy(selectedHost));
                _power[VmOp.Reboot].IsEnabled = HostRules.CanReboot(selectedHost, Busy(selectedHost));
            }
            var serverHost = o as Host;
            _hostReboot.IsEnabled = HostRules.CanReboot(serverHost, Busy(serverHost));
            _hostShutdown.IsEnabled = HostRules.CanShutDown(serverHost, Busy(serverHost));
            _hostMaintenance.IsEnabled = HostRules.CanToggleMaintenance(serverHost, Busy(serverHost));
            _restartToolstack.IsEnabled = HostRules.CanRestartToolstack(serverHost) && !Busy(serverHost);
            _removeFromPool.IsEnabled = HostRules.CanRemoveFromPool(serverHost) && !Busy(serverHost);
            _hostMaintenance.Header = Dialogs.M(serverHost != null ? HostRules.MaintenanceText(serverHost) : Messages.ENTER_MAINTENANCE_MODE_CONTEXT_MENU);
            foreach (var (op, m) in _vmMenu)
                m.IsEnabled = Allowed(op) && !(haBlocker != null && (op == VmOp.Start || op == VmOp.Resume));
            ToolTip.SetTip(_vmMenu[VmOp.Start], haBlocker);
            ToolTip.SetTip(_vmMenu[VmOp.Resume], haBlocker);
            var c = node?.Connection;
            _connectMenu.IsEnabled = c != null && !c.IsConnected && !c.InProgress;
            _disconnectMenu.IsEnabled = c != null && (c.IsConnected || c.InProgress);
            _removeMenu.IsEnabled = c != null;
            _poolProps.IsEnabled = o is Pool && VmRules.CanEditProperties(o);
            // PoolReconnectAsCommand / DisconnectPoolCommand: any object of a connection that is a pool (PoolAncestor)
            var poolConnection = c != null && XenAdmin.Core.Helpers.GetPool(c) != null ? c : null;
            _poolReconnectAs.IsEnabled = poolConnection != null;
            _poolDisconnect.IsEnabled = poolConnection != null && (poolConnection.IsConnected || poolConnection.InProgress);
            _poolPassword.IsEnabled = o is Pool && PoolRules.CanChangePassword(o);
            _hostPassword.IsEnabled = o is Host && PoolRules.CanChangePassword(o);
            _rotateSecret.IsEnabled = (o is Pool || o is Host) && PoolRules.CanRotateSecret(c);
            var selectedPool = poolConnection == null ? null : XenAdmin.Core.Helpers.GetPool(poolConnection);
            _poolAddServer.IsEnabled = selectedPool != null && poolConnection.IsConnected;
            _poolRemoveServer.IsEnabled = selectedPool != null && poolConnection.IsConnected;
            _makeStandalone.IsEnabled = PoolRules.CanMakeStandalone(selectedPool);
            var haBusy = poolConnection != null && _haBusy.Contains(poolConnection);
            _haConfigure.IsEnabled = HaRules.CanConfigure(selectedPool, haBusy);
            _haDisable.IsEnabled = HaRules.CanDisable(selectedPool, haBusy);
            _hostProps.IsEnabled = o is Host && VmRules.CanEditProperties(o);
            _vmProps.IsEnabled = o is VM && VmRules.CanEditProperties(o);
            _newVm.IsEnabled = _newVmMenu.IsEnabled = c != null && c.IsConnected;
            _newStorage.IsEnabled = _newSrMenu.IsEnabled = c != null && c.IsConnected;
            _newPool.IsEnabled = _newPoolMenu.IsEnabled = _connections.Any(x => x.IsConnected);
            _copyVmMenu.IsEnabled = vm != null && VmRules.CanCopy(vm);
            _snapshotMenu.IsEnabled = VmRules.CanTakeSnapshot(vm);
            _toTemplateMenu.IsEnabled = VmRules.CanConvertToTemplate(vm);
            _exportVmMenu.IsEnabled = XvaRules.CanExport(vm);
            _vappsMenu.IsEnabled = VappRules.CanManage(c);
            _drMenu.IsEnabled = DrRules.CanRun(c);
            _assignVappMenu.IsEnabled = VappRules.CanAssign(vm);
            _vtpmMenu.IsEnabled = VtpmRules.CanManage(vm);
            ToolTip.SetTip(_vtpmMenu, VtpmRules.Blocker(vm));
            ToolTip.SetTip(_exportVmMenu, XvaRules.ExportBlocker(vm));
            var template = o as VM;
            _fromTemplateMenu.IsEnabled = VmRules.IsUserTemplate(template) && template.Connection.IsConnected;
            _copyTemplateMenu.IsEnabled = VmRules.CanCopyTemplate(template);
            _deleteTemplateMenu.IsEnabled = VmRules.CanDeleteTemplate(template);
            _templateProps.IsEnabled = template != null && template.is_a_template && VmRules.CanEditProperties(template);
            _migrateMenu.IsEnabled = vm != null && MigrateRules.CanMigrate(vm);
            _deleteVmMenu.IsEnabled = vm != null && VmRules.CanDeleteVm(vm);

            UpdateAlertCount();
            if (_showEvents)
            {
                ShowNotifications();
                return;
            }
            if (_leftPane.Child != _tree)
                _leftPane.Child = _tree;

            var user = c?.Session == null ? null : $"Logged in as:  {(c.Session.IsLocalSuperuser ? "Local root account" : c.Session.UserFriendlyName())}";
            // a group of an organization view or a saved search (or its root): the Search tab lists its objects
            // (a vApp, or a network of "VMs by network", groups VMs too: its page lists them as in WinForms)
            if (node != null && _navView != NavView.Infrastructure
                && (node.IsRoot || o == null || o is XenAdmin.Model.Folder || o is VM_appliance || node.Children.Count > 0 && !(o is VM || o is Host || o is Pool || o is SR)))
            {
                _titleBand.Content = Views.TitleBand(node.Icon, node.Name, null);
                var groupKey = "group:" + _navView + ":" + node.Key;
                var shown = _navView == NavView.SavedSearch && node.IsRoot ? _savedSearch : null;
                SetTabs(groupKey, () => new[] { Tab("Search", Body(node.Name, new SearchView(_connections, () => NavigationViews.Objects(SelectedNode() ?? node), shown, shown != null))) });
                if (_tabs.Items.OfType<TabItem>().FirstOrDefault()?.Content is Panel groupPanel)
                    groupPanel.Children.OfType<ILiveView>().FirstOrDefault()?.Refresh();
                return;
            }
            if (node == null || node.IsRoot)
            {
                _titleBand.Content = Views.TitleBand("Logo", ResourceTree.RootName, null);
                SetTabs("home", () => new[] { Tab("Home", Views.HomePage(async () => await AddServer())),
                    Tab("Search", Body(OverviewColumns.Title(null), new SearchView(_connections))) });
                if (_tabs.Items.OfType<TabItem>().FirstOrDefault(t => (string)t.Header == "Search")?.Content is Panel home)
                    home.Children.OfType<ILiveView>().FirstOrDefault()?.Refresh();
                return;
            }
            if (o == null)
            {
                _titleBand.Content = Views.TitleBand(node.Icon, node.Name, null);
                // GeneralTabPage.GenerateDisconnectedHostBox
                SetTabs("conn:" + c?.HostnameWithPort, () => new[] { Tab("General", Body("Connection General Properties",
                    Views.Sections(new[] { new Section { Title = "General", Expanded = true }
                        .Add(FriendlyNameManager.GetFriendlyName("Label-host.name_label"), (c as XenConnection)?.Name ?? c?.Hostname)
                        .Add(FriendlyNameManager.GetFriendlyName("Label-host.address"), c?.Hostname) },
                        null))),
                    new TabItem { Header = "Search", IsEnabled = false } });
                return;
            }

            _titleBand.Content = Views.TitleBand(node.Icon, Details.Heading(o), user);
            // rebuild the tab bodies on every refresh so values stay live, but keep the selected tab
            // the object and its tab list: a tab that comes or goes (USB plugged in, host down) rebuilds the tabs
            var sig = o.opaque_ref + "|" + string.Join(",", Details.Tabs(o));
            string current = _lastTabAny;
            if (_config.Options.RememberLastSelectedTab)
                _lastTab.TryGetValue(o.opaque_ref, out current);
            if (_tabsSignature == sig)
            {
                foreach (var existing in _tabs.Items.OfType<TabItem>().Where(i => i.IsEnabled && (string)i.Header != "Console"))
                {
                    // stateful tabs (selection, live metrics) refresh themselves; the rest are rebuilt
                    if (existing.Content is Panel p && p.Children.OfType<ILiveView>().FirstOrDefault() is ILiveView live)
                        live.Refresh();
                    else
                        existing.Content = TabBody(o, (string)existing.Header);
                }
                return;
            }
            var items = Details.Tabs(o).Select(t => Details.NotPorted(o, t)
                ? new TabItem { Header = t, IsEnabled = false }
                : Tab(t, TabBody(o, t))).ToList();
            _tabsSignature = sig;
            _tabs.ItemsSource = items;
            _tabs.SelectedItem = items.FirstOrDefault(i => (string)i.Header == current && i.IsEnabled) ?? items.First();
        }

        private void SetTabs(string signature, Func<TabItem[]> build)
        {
            if (_tabsSignature == signature) return;
            _tabsSignature = signature;
            var items = build();
            _tabs.ItemsSource = items;
            _tabs.SelectedItem = items[0];
        }

        private static TabItem Tab(string header, Control content) => new TabItem { Header = header, Content = content };

        private static Control Body(string title, Control content)
        {
            var d = new DockPanel { Background = Brushes.White };
            var h = Views.TabHeader(title);
            DockPanel.SetDock(h, Dock.Top);
            d.Children.Add(h);
            d.Children.Add(content);
            return d;
        }

        private Control TabBody(IXenObject o, string tab)
        {
            switch (tab)
            {
                case "General": return Body(Details.TabTitle(o, tab), Views.Sections(Details.General(o), VmRules.CanEditProperties(o) ? () => ShowProperties(o) : null));
                case "Memory": return Body("Memory", new MemoryView(o, RunTracked));
                case "Storage":
                    var st = Details.Storage(o);
                    return Body(st.Title, o switch
                    {
                        VM vmStorage => new VmStorageView(vmStorage, RunTracked, ShowProperties),
                        SR srStorage => new SrStorageView(srStorage, RunTracked, ShowProperties),
                        _ => new PhysicalStorageView(o, RunTracked, ShowProperties, ShowNewSr)
                    });
                case "Networking":
                    var nt = Details.Networking(o);
                    return Body(nt.Title, o is VM vmNetworking
                        ? new VmNetworkView(vmNetworking, RunTracked)
                        : new HostNetworkView(o, RunTracked, ShowProperties));
                case "NICs":
                    return Body("Network Interface Cards", new NicsView((Host)o, RunTracked));
                case "GPU":
                    return Body("GPU", new ScrollViewer { Content = new GpuView(o, RunTracked) });
                case "USB":
                    return Body(Messages.USB, new UsbView((Host)o, RunTracked));
                case "Console":
                    // TabPageConsole of a server: the console of its control domain
                    return ConsoleView.For(o is Host consoleHost ? consoleHost.ControlDomainZero() : o as VM);
                case "Performance":
                    return Body("Performance Graphs", new PerformanceView(o));
                case "Snapshots":
                    return Body("Virtual Machine Snapshots", new SnapshotsView((VM)o, RunTracked, ShowProperties));
                case "Search":
                    return Body(OverviewColumns.Title(o.Connection), new SearchView(new[] { o.Connection }));
            }
            return new TextBlock();
        }

        // ---------- server and VM commands of the menus ----------

        /// <summary>RestartToolstackCommand: confirm, then RestartToolstackAction (the connection reconnects by itself).</summary>
        private async void RestartToolstack(Host host)
        {
            if (!HostRules.CanRestartToolstack(host)) return;
            var (title, text) = HostRules.RestartToolstackConfirmation(host);
            if (await Dialogs.Confirm(this, title, text, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION))
                RunTracked(new RestartToolstackAction(host));
        }

        /// <summary>
        /// RemoveHostFromPoolCommand: confirm (No is the default), then EjectHostFromPoolAction; XAPI resets and reboots
        /// the server. shortcut: the ejected server is not added back to the tree after its reboot (WinForms WaitForReboot);
        /// add it with Add Server for now.
        /// </summary>
        private async void RemoveFromPool(Host host)
        {
            if (!HostRules.CanRemoveFromPool(host)) return;
            var (title, text, yes) = HostRules.RemoveFromPoolConfirmation(host);
            if (await Dialogs.Confirm(this, title, text, yes, Messages.NO_BUTTON_CAPTION, cancelIsDefault: true))
                RunTracked(new EjectHostFromPoolAction(Helpers.GetPool(host.Connection), host));
        }

        /// <summary>PoolReconnectAsCommand: disconnect, then connect again with other credentials.</summary>
        private async System.Threading.Tasks.Task ReconnectAs(XenConnection c)
        {
            var input = await Dialogs.AddServer(this, c.Port == 443 ? c.Hostname : $"{c.Hostname}:{c.Port}", c.Username ?? "root", true);
            if (input == null) return;
            Disconnect(c);
            c.Username = input.User;
            c.Password = input.Password;
            StartConnect(c);
        }

        /// <summary>ChangeHostPasswordCommand.</summary>
        private async void ChangePassword(IXenObject hostOrPool)
        {
            if (!PoolRules.CanChangePassword(hostOrPool)) return;
            if (await Dialogs.ChangeServerPassword(this, hostOrPool) is AsyncAction a) RunTracked(a);
        }

        /// <summary>RotatePoolSecretCommand: the Pool Admin restriction, the password reminder, then pool.rotate_secret.</summary>
        private async void RotateSecret(IXenConnection c)
        {
            if (!PoolRules.CanRotateSecret(c)) return;
            if (PoolRules.RotateSecretRestriction(c) is string restriction)
            {
                await Dialogs.Info(this, Messages.ROTATE_POOL_SECRET_TITLE, restriction);
                return;
            }
            if (_config.Options.RotateSecretReminder)
                await Dialogs.Info(this, Messages.ROTATE_POOL_SECRET_TITLE, Messages.ROTATE_POOL_SECRET_REMIND_CHANGE_PASSWORD);
            RunTracked(PoolRules.RotateSecretAction(c));
        }

        /// <summary>ConvertVMToTemplateCommand: the Convert / Cancel warning, then the two actions.</summary>
        private async void ConvertToTemplate(VM vm)
        {
            if (!VmRules.CanConvertToTemplate(vm)) return;
            if (await Dialogs.Confirm(this, Messages.CONVERT_TO_TEMPLATE, string.Format(Messages.CONVERT_TEMPLATE_DIALOG_TEXT, vm.Name().Ellipsise(25)), Messages.CONVERT, Messages.CANCEL))
                RunTracked(VmRules.ConvertToTemplateAction(vm));
        }

        /// <summary>TakeSnapshotCommand: the same Take Snapshot dialog as the Snapshots tab.</summary>
        private async void TakeSnapshot(VM vm)
        {
            if (!VmRules.CanTakeSnapshot(vm)) return;
            var action = await Dialogs.TakeSnapshot(this, vm);
            if (action != null) RunTracked(action);
        }

        // ---------- tree events (as NavigationView / ContextMenuBuilder) ----------

        /// <summary>Disconnected server: connect; halted VM: start; suspended VM: resume.</summary>
        private async System.Threading.Tasks.Task DoubleClick(TreeNode n)
        {
            if (n.Connection is XenConnection c && !c.IsConnected)
            {
                await Connect(c);
                return;
            }
            if (n.XenObject is VM vm && !vm.is_a_template)
            {
                if (vm.power_state == vm_power_state.Halted && vm.allowed_operations.Contains(vm_operations.start))
                    RunVmOp(VmOp.Start);
                else if (vm.power_state == vm_power_state.Suspended && vm.allowed_operations.Contains(vm_operations.resume))
                    RunVmOp(VmOp.Resume);
            }
        }

        /// <summary>Only the commands that apply, in the WinForms order (it uses AddIfEnabled).</summary>
        private List<object> ContextMenuFor(TreeNode n)
        {
            var items = new List<object>();
            MenuItem Item(string text, Action a, bool bold = false, string icon = null)
            {
                var m = new MenuItem { Header = Dialogs.M(text) };
                if (icon != null) m.Icon = Views.Icon(icon, 16);
                if (bold) m.FontWeight = FontWeight.Bold;
                m.Click += (_, _) => a();
                return m;
            }
            var c = n.Connection as XenConnection;
            // the server's part of ContextMenuBuilder: maintenance mode, reboot, shut down
            void AddHostCommands(Host host)
            {
                if (HostRules.CanToggleMaintenance(host, Busy(host)))
                    items.Add(Item(HostRules.MaintenanceText(host), () => ToggleMaintenance(host), icon: "000_ServerMaintenance_h32bit_16"));
                if (HostRules.CanReboot(host, Busy(host)))
                    items.Add(Item(Messages.MAINWINDOW_REBOOT, () => HostPower(host, true), icon: "001_Reboot_h32bit_16"));
                if (HostRules.CanShutDown(host, Busy(host)))
                    items.Add(Item(Messages.MAINWINDOW_SHUTDOWN, () => HostPower(host, false), icon: "001_ShutDown_h32bit_16"));
                if (HostRules.CanRestartToolstack(host) && !Busy(host))
                    items.Add(Item(Messages.MAINWINDOW_RESTART_TOOLSTACK, () => RestartToolstack(host)));
                if (HostRules.CanRemoveFromPool(host) && !Busy(host))
                    items.Add(Item(Messages.REMOVE_SERVER_FROM_POOL_CONTEXT_MENU_ITEM_TEXT, () => RemoveFromPool(host)));
                items.Add(new Separator());
            }
            if (n.XenObject is SR sr && sr.Connection.IsConnected)
            {
                // the SR part of ContextMenuBuilder (shortcut: Repair, Reattach and Reclaim space are not ported yet)
                async void Run(string command, Func<AsyncAction> action)
                {
                    var (title, text, yes) = SrRules.Confirmation(sr, command);
                    if (title == null || await Dialogs.Confirm(this, title, text, yes, Messages.CANCEL, cancelIsDefault: command == "destroy"))
                        RunTracked(action());
                }
                if (SrRules.CanSetDefault(sr))
                    items.Add(Item(Messages.MAINWINDOW_SET_AS_DEFAULT_SR, () => Run("default", () => new SetSrAsDefaultAction(sr)), icon: "000_StorageDefault_h32bit_16"));
                if (SrRules.CanDetach(sr))
                    items.Add(Item(Messages.MAINWINDOW_DETACH_SR, () => Run("detach", () => new DetachSrAction(sr))));
                if (SrRules.CanForget(sr))
                    items.Add(Item(Messages.MAINWINDOW_FORGET_SR, () => Run("forget", () => new ForgetSrAction(sr))));
                if (SrRules.CanDestroy(sr))
                    items.Add(Item(Messages.MAINWINDOW_DESTROY_SR, () => Run("destroy", () => new DestroySrAction(sr))));
            }
            if (n.XenObject is VM vm)
            {
                // the WinForms commands' MenuImage
                string IconOf(VmOp op) => op switch
                {
                    VmOp.Start => "001_PowerOn_h32bit_16", VmOp.Shutdown => "001_ShutDown_h32bit_16", VmOp.Reboot => "001_Reboot_h32bit_16",
                    VmOp.Suspend => "000_paused_h32bit_16", VmOp.Resume => "000_Resumed_h32bit_16",
                    VmOp.ForceShutdown => "001_ForceShutDown_h32bit_16", VmOp.ForceReboot => "001_ForceReboot_h32bit_16", _ => null
                };
                void Add(VmOp op, string text, bool bold = false) { if (VmOperations.IsAllowed(vm, op)) items.Add(Item(text, () => RunVmOp(op), bold, IconOf(op))); }
                Add(VmOp.Shutdown, Messages.MAINWINDOW_SHUTDOWN);
                Add(VmOp.Start, Messages.MAINWINDOW_START, vm.power_state == vm_power_state.Halted);
                Add(VmOp.Resume, Messages.MAINWINDOW_RESUME);
                Add(VmOp.Suspend, Messages.MAINWINDOW_SUSPEND);
                Add(VmOp.Reboot, Messages.MAINWINDOW_REBOOT);
                var before = items.Count;
                Add(VmOp.ForceShutdown, Messages.MAINWINDOW_FORCE_SHUTDOWN);
                Add(VmOp.ForceReboot, Messages.MAINWINDOW_FORCE_REBOOT);
                if (before > 0 && items.Count > before) items.Insert(before, new Separator());
                if (MigrateRules.CanMigrate(vm)) { items.Add(new Separator()); items.Add(MigrateMenu(() => vm)); }
                // shortcut: Assign to Snapshot Schedule / vApp are not ported yet
                if (VmRules.SnapshotModesFor(vm) is var modes && (modes.Disk || modes.Memory))
                {
                    items.Add(new Separator());
                    items.Add(Item("Take a S&napshot...", async () => { if (await Dialogs.TakeSnapshot(this, vm) is AsyncAction a) RunTracked(a); }));
                }
                if (VmRules.CanCopy(vm)) { items.Add(new Separator()); items.Add(Item("Cop&y VM...", () => CopyVm(vm))); }
                if (VmRules.CanDeleteVm(vm)) { items.Add(new Separator()); items.Add(Item("&Delete VM...", () => DeleteVm(vm))); }
            }
            else if (c != null && n.XenObject == null)
            {
                if (c.InProgress)
                    items.Add(Item(Messages.CANCEL_CONNECT, () => { EndConnecting(c); c.EndConnect(); RefreshTree(); }));
                else
                {
                    items.Add(Item(Messages.MAINWINDOW_CONNECTHOST, async () => await Connect(c), true));
                    if (!string.IsNullOrEmpty(c.Password))
                        items.Add(Item(Messages.MAINWINDOW_FORGET_SAVED_PASSWORD, () => c.Password = null));
                    items.Add(Item(string.Format(Messages.MAINWINDOW_REMOVE_HOST, "XCP-ng Center"), RemoveSelected));
                }
            }
            else if (c != null && c.IsConnected && (n.XenObject is Pool || n.XenObject is Host h && Helpers.GetPool(c) == null))
            {
                items.Add(Item("&New VM...", () => ShowNewVm(), icon: "001_CreateVM_h32bit_16"));
                items.Add(new Separator());
                if (n.XenObject is Host standalone) AddHostCommands(standalone);
                items.Add(Item(n.XenObject is Pool ? Messages.MAINWINDOW_DISCONNECT_CONTEXT_MENU_POOL : Messages.MAINWINDOW_DISCONNECT_CONTEXT_MENU, () => Disconnect(c)));
            }
            if (n.XenObject is Host member && Helpers.GetPool(member.Connection) != null && member.Connection.IsConnected)
            {
                items.Add(Item("&New VM...", () => ShowNewVm(), icon: "001_CreateVM_h32bit_16"));
                items.Add(new Separator());
                AddHostCommands(member);
            }
            if (VmRules.CanEditProperties(n.XenObject))
            {
                if (items.Count > 0) items.Add(new Separator());
                items.Add(Item("P&roperties", () => ShowProperties(n.XenObject), icon: "edit_16"));
            }
            // no leading, trailing or doubled separators (a group can turn out empty)
            for (var i = items.Count - 1; i >= 0; i--)
                if (items[i] is Separator && (i == 0 || i == items.Count - 1 || items[i - 1] is Separator))
                    items.RemoveAt(i);
            return items;
        }

        // ---------- New VM / Copy VM / Delete VM ----------



        /// <summary>
        /// MigrateVMToolStripMenuItem: one item per host, disabled with the reason when the VM cannot go there.
        /// The reasons need VM.assert_can_boot_here, so they load when the submenu opens.
        /// </summary>
        /// <summary>
        /// Pool > Add Server (the other connected standalone servers, each with the reason it cannot join) and Remove
        /// Server (the members; RemoveHostFromPoolCommand), filled when the submenu opens.
        /// </summary>
        /// <summary>DRConfigureCommand: the SRs holding the pool's VM metadata.</summary>
        private async void ConfigureDr(IXenConnection c)
        {
            if (!DrRules.CanRun(c) || !(XenAdmin.Core.Helpers.GetPoolOfOne(c) is Pool pool)) return;
            if (await DrConfigureDialog.Show(this, pool) is AsyncAction a) RunTracked(a);
        }

        /// <summary>VMGroupCommand&lt;VM_appliance&gt;: Manage vApps for the selected pool or server.</summary>
        private void ManageVapps(IXenConnection c)
        {
            if (VappRules.CanManage(c))
                new VappsDialog(c, RunTracked, o => ShowProperties(_vappsDialogOwner ?? this, o)).Show(this);
        }

        private Window _vappsDialogOwner => OwnedWindows.OfType<VappsDialog>().FirstOrDefault();

        /// <summary>
        /// AssignGroupToolStripMenuItem&lt;VM_appliance&gt;: New vApp... (with this VM), then each vApp; the VM's own vApp is
        /// checked and choosing it takes the VM out.
        /// </summary>
        private MenuItem AssignVappMenu()
        {
            var menu = new MenuItem { Header = Dialogs.M(Messages.ASSIGN_VM_APPLIANCE) };
            menu.Items.Add(new MenuItem { Header = Messages.DELAY_LOADING_COMBO_BOX_WAITING, IsEnabled = false });
            menu.SubmenuOpened += (_, e) =>
            {
                if (e.Source != menu || !(SelectedNode()?.XenObject is VM vm) || !VappRules.CanAssign(vm)) return;
                menu.Items.Clear();
                var add = new MenuItem { Header = Dialogs.M(Messages.NEW_VM_APPLIANCE), Icon = Views.Icon("000_NewVirtualAppliance_h32bit_16") };
                add.Click += (_, _) => new VappWizard(vm.Connection, RunTracked, new[] { vm }).Show(this);
                menu.Items.Add(add);
                var apps = VappRules.Appliances(vm.Connection);
                if (apps.Count > 0) menu.Items.Add(new Separator());
                for (var i = 0; i < apps.Count; i++)
                {
                    var app = apps[i];
                    var text = (i < 9 ? $"_{i + 1} " : "") + app.Name().Ellipsise(30).Replace("_", "__");
                    var item = new MenuItem { Header = text, ToggleType = MenuItemToggleType.CheckBox, IsChecked = VappRules.Of(vm)?.opaque_ref == app.opaque_ref };
                    item.Click += async (_, _) =>
                    {
                        var vms = new List<VM> { vm };
                        if (VappRules.ChangeWarning(vms, app, app.Name()) is string warning
                            && !await Dialogs.Confirm(this, Messages.CHANGE_VMS_APPLIANCE, warning, Messages.YES, Messages.NO))
                            return;
                        RunTracked(VappRules.AssignAction(app, vms));
                    };
                    menu.Items.Add(item);
                }
            };
            return menu;
        }

        private MenuItem PoolServersMenu(string header, bool add)
        {
            var menu = new MenuItem { Header = header };
            menu.Items.Add(new MenuItem { Header = Messages.DELAY_LOADING_COMBO_BOX_WAITING, IsEnabled = false });
            menu.SubmenuOpened += (_, e) =>
            {
                if (e.Source != menu || SelectedNode()?.Connection is not IXenConnection c || XenAdmin.Core.Helpers.GetPool(c) is not Pool pool) return;
                menu.Items.Clear();
                if (add)
                {
                    foreach (var candidate in PoolRules.JoinCandidates(_connections, pool))
                    {
                        var host = XenAdmin.Core.Helpers.GetCoordinator(candidate.Connection);
                        var item = new MenuItem { Header = candidate.AllowedAsSupporter ? candidate.Name : $"{candidate.Name} ({candidate.Description})", IsEnabled = candidate.AllowedAsSupporter && host != null };
                        item.Click += (_, _) => AddToPool(host, pool);
                        menu.Items.Add(item);
                    }
                    if (menu.Items.Count > 0) menu.Items.Add(new Separator());
                    var other = new MenuItem { Header = Dialogs.M(Messages.ADD_POOL_OR_SERVER) };
                    other.Click += async (_, _) => await AddServer();
                    menu.Items.Add(other);
                }
                else
                {
                    foreach (var host in c.Cache.Hosts.OrderBy(h => h.Name()))
                    {
                        var item = new MenuItem { Header = host.Name(), IsEnabled = HostRules.CanRemoveFromPool(host) && !Busy(host) };
                        item.Click += (_, _) => RemoveFromPool(host);
                        menu.Items.Add(item);
                    }
                }
                if (menu.Items.Count == 0)
                    menu.Items.Add(new MenuItem { Header = Messages.NONE, IsEnabled = false });
            };
            return menu;
        }

        /// <summary>AddHostToPoolCommand: the confirmation and the pool warnings, then AddHostToPoolAction.</summary>
        private async void AddToPool(Host host, Pool pool)
        {
            foreach (var message in PoolRules.AddToPoolMessages(host, pool))
                if (!await Dialogs.Confirm(this, Messages.POOLCREATE_ADDING, message, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION))
                    return;
            RunTracked(PoolRules.AddToPool(host, pool));
        }

        /// <summary>
        /// HAConfigureCommand: the HA wizard when HA is off; with HA on, the restart priorities (EditVmHaPrioritiesDialog
        /// with SetHaPrioritiesAction), unless the statefile cannot be found.
        /// </summary>
        private async void ConfigureHa(IXenConnection c)
        {
            var pool = c == null ? null : XenAdmin.Core.Helpers.GetPool(c);
            if (!HaRules.CanConfigure(pool, _haBusy.Contains(c))) return;
            if (HaRules.ConfigureBlocker(pool) is string blocker)
            {
                await Dialogs.Info(this, Messages.CONFIGURE_HA, blocker);
                return;
            }
            if (!pool.ha_enabled)
            {
                await new HaWizard(pool, a => RunHa(c, a)).ShowDialog(this);
                return;
            }
            var plan = new HaPlanEditor(pool);
            var ok = new Button { Content = "OK", IsDefault = true };
            var cancel = new Button { Content = Dialogs.M(Messages.CANCEL), IsCancel = true };
            var error = new TextBlock { Foreground = Brushes.Red, TextWrapping = TextWrapping.Wrap };
            plan.Changed += () => { error.Text = plan.Error ?? ""; ok.IsEnabled = plan.Error == null; };
            var buttons = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, HorizontalAlignment = HorizontalAlignment.Right };
            buttons.Children.Add(ok);
            buttons.Children.Add(cancel);
            var panel = new StackPanel { Spacing = 10, Margin = new Thickness(14, 12) };
            panel.Children.Add(plan);
            panel.Children.Add(error);
            panel.Children.Add(buttons);
            var w = new Window { Title = "HA restart priorities", Width = 760, SizeToContent = SizeToContent.Height, Content = panel, WindowStartupLocation = WindowStartupLocation.CenterOwner, ShowInTaskbar = false };
            ok.Click += async (_, _) =>
            {
                if (plan.Error != null) return;
                // EditVmHaPrioritiesDialog: tolerating no failure is practically useless, so ask (No is the default)
                if (plan.Ntol == 0 && !await Dialogs.Confirm(w, Messages.HIGH_AVAILABILITY, Messages.HA_NTOL_ZERO_QUERY, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION, cancelIsDefault: true))
                    return;
                RunHa(c, HaRules.SetPrioritiesAction(c, plan.Options(), plan.Ntol));
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(this);
        }

        /// <summary>HADisableCommand: Yes/No, then DisableHAAction.</summary>
        private async void DisableHa(IXenConnection c)
        {
            var pool = c == null ? null : XenAdmin.Core.Helpers.GetPool(c);
            if (!HaRules.CanDisable(pool, _haBusy.Contains(c))) return;
            var (title, text) = HaRules.DisableConfirmation(pool);
            if (await Dialogs.Confirm(this, title, text, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION))
                RunHa(c, HaRules.DisableAction(pool));
        }

        /// <summary>An HA action keeps the HA commands disabled until it ends, as FindActiveHaAction does.</summary>
        private void RunHa(IXenConnection c, AsyncAction action)
        {
            _haBusy.Add(c);
            action.Completed += _ => Dispatcher.UIThread.Post(() => { _haBusy.Remove(c); UpdateContent(); });
            RunTracked(action);
        }

        /// <summary>DeletePoolCommand ("Make into standalone server"): DestroyPoolAction on a pool of one server.</summary>
        private async void MakeStandalone(IXenConnection c)
        {
            var pool = c == null ? null : XenAdmin.Core.Helpers.GetPool(c);
            if (!PoolRules.CanMakeStandalone(pool)) return;
            if (PoolRules.MakeStandaloneBlocker(pool) is string blocker)
            {
                await Dialogs.Info(this, Messages.MAINWINDOW_DELETE_POOL.Replace("&", ""), blocker);
                return;
            }
            RunTracked(new DestroyPoolAction(pool));
        }

        private MenuItem MigrateMenu(Func<VM> current)
        {
            var menu = new MenuItem { Header = Dialogs.M(Messages.MAINWINDOW_MIGRATE_TO_SERVER) };
            menu.Items.Add(new MenuItem { Header = Messages.DELAY_LOADING_COMBO_BOX_WAITING, IsEnabled = false });
            var generation = 0;
            menu.SubmenuOpened += (_, e) =>
            {
                if (e.Source != menu || current() is not VM vm) return;
                // drop the previous VM's hosts at once; a slower, older query must not publish its items
                var mine = ++generation;
                menu.Items.Clear();
                menu.Items.Add(new MenuItem { Header = Messages.DELAY_LOADING_COMBO_BOX_WAITING, IsEnabled = false });
                System.Threading.ThreadPool.QueueUserWorkItem(_ =>
                {
                    List<(Host Host, string Reason)> targets;
                    (Host Host, string Reason) home;
                    try
                    {
                        var session = vm.Connection.DuplicateSession(); // shared session: never log it out
                        home = MigrateRules.HomeServerTarget(vm, session);
                        targets = MigrateRules.Targets(vm, session);
                    }
                    catch { home = (null, Messages.HOST_MENU_UNKNOWN_ERROR); targets = new List<(Host, string)>(); }
                    Dispatcher.UIThread.Post(() =>
                    {
                        if (mine != generation) return;
                        menu.Items.Clear();
                        var homeItem = new MenuItem { Header = Dialogs.M(MigrateRules.HomeServerMenuText(home.Reason).Replace("_", "__")), IsEnabled = home.Reason == null, Icon = Views.Icon("000_ServerHome_h32bit_16", 16) };
                        homeItem.Click += (_, _) => { if (mine == generation && home.Host != null && MigrateRules.CanMigrate(vm)) RunTracked(MigrateRules.Migrate(vm, home.Host)); };
                        menu.Items.Add(homeItem);
                        foreach (var (host, reason) in targets)
                        {
                            var item = new MenuItem { Header = MigrateRules.MenuText(host, reason).Replace("_", "__"), IsEnabled = reason == null, Icon = Views.Icon("000_TreeConnected_h32bit_16", 16) };
                            item.Click += (_, _) => { if (mine == generation && MigrateRules.CanMigrate(vm)) RunTracked(MigrateRules.Migrate(vm, host)); };
                            menu.Items.Add(item);
                        }
                        menu.Items.Add(new Separator());
                        menu.Items.Add(new MenuItem { Header = Dialogs.M("&Migrate VM wizard..."), IsEnabled = false, Icon = Views.Icon("000_MigrateVM_h32bit_16", 16) });
                    });
                });
            };
            return menu;
        }

        /// <summary>ShutDownHostCommand / RebootHostCommand: always confirmed, with the coordinator warning.</summary>
        private async void HostPower(Host host, bool reboot)
        {
            if (!(reboot ? HostRules.CanReboot(host, Busy(host)) : HostRules.CanShutDown(host, Busy(host)))) return;
            var (title, text) = HostRules.Confirm(host, reboot);
            if (await Dialogs.Confirm(this, title, text, reboot ? Messages.CONFIRM_REBOOT_SERVER_YES_BUTTON_LABEL : Messages.CONFIRM_SHUTDOWN_SERVER_YES_BUTTON_LABEL, Messages.CANCEL, cancelIsDefault: true))
                RunHostAction(host, reboot ? HostRules.Reboot(host) : HostRules.ShutDown(host));
        }

        /// <summary>HostMaintenanceModeCommand: enter (evacuate) or exit (enable, optionally restoring the VMs).</summary>
        private async void ToggleMaintenance(Host host)
        {
            if (!HostRules.CanToggleMaintenance(host, Busy(host))) return;
            if (!HostRules.InMaintenance(host))
            {
                if (HostRules.CannotEnterMaintenance(host) is string why)
                {
                    await Dialogs.Info(this, "XCP-ng Center", why);
                    return;
                }
                var (title, text) = HostRules.EvacuateConfirmation(host);
                if (await Dialogs.Confirm(this, title, text, "&Enter Maintenance Mode", Messages.CANCEL))
                    RunHostAction(host, HostRules.EnterMaintenance(host));
                return;
            }
            var restore = HostRules.VmsToRestore(host);
            var answer = restore.Count == 0 ? false : await Dialogs.Choose(this, "Exit Maintenance Mode",
                $"Would you like to restore your VMs to their previous location on '{host.Name()}'?\n\nVMs to restore:\n" + string.Join("\n", restore.Select(v => v.Name())),
                "&Restore VMs", "&Skip");
            if (answer is bool restoreVms) // Cancel leaves the server in maintenance mode
                RunHostAction(host, HostRules.ExitMaintenance(host, restoreVms));
        }

        // HelpersGUI.HasActiveHostAction: one server action at a time per host
        private readonly HashSet<string> _busyHosts = new HashSet<string>();
        private bool Busy(Host host) => host != null && _busyHosts.Contains(host.opaque_ref);

        private void RunHostAction(Host host, AsyncAction action)
        {
            _busyHosts.Add(host.opaque_ref);
            action.Completed += _ => Dispatcher.UIThread.Post(() => { _busyHosts.Remove(host.opaque_ref); UpdateContent(); });
            RunTracked(action);
            UpdateContent();
        }

        /// <summary>NewSRCommand: the wizard for the selected (or first connected) pool or server.</summary>
        private async void ShowNewSr()
        {
            var c = SelectedNode()?.Connection ?? _connections.FirstOrDefault(x => x.IsConnected);
            if (c == null || !c.IsConnected) return;
            await new NewSrWizard(c, RunTracked).ShowDialog(this);
        }

        private async void ShowNewVm(VM template = null)
        {
            var c = template?.Connection ?? SelectedNode()?.Connection ?? _connections.FirstOrDefault(x => x.IsConnected);
            if (c == null || !c.IsConnected) return;
            var wizard = new NewVmWizard(c, template);
            await wizard.ShowDialog(this);
            if (wizard.Action != null)
                RunTracked(wizard.Action);
        }

        /// <summary>Servers that joined a pool are removed from ConnectionsManager by XenModel; forget them here too.</summary>
        private void PruneDroppedConnections()
        {
            List<XenConnection> gone;
            lock (ConnectionsManager.ConnectionsLock)
                gone = _connections.Where(x => !ConnectionsManager.XenConnections.Contains(x)).ToList();
            foreach (var x in gone)
            {
                _connections.Remove(x);
                _config.RemoveServer(x.Hostname, x.Port);
            }
        }

        /// <summary>NewPoolCommand.</summary>
        private async void ShowNewPool()
        {
            var pre = SelectedNode()?.Connection;
            var action = await Dialogs.NewPool(this, () => _connections, pre != null && XenAdmin.Core.Helpers.GetPool(pre) == null ? pre : null, AddServer);
            if (action == null) return;
            // XenModel drops the members' connections on a background thread after the join
            action.Completed += _ => Dispatcher.UIThread.Post(() => DispatcherTimer.RunOnce(RefreshTree, TimeSpan.FromSeconds(2)));
            RunTracked(action);
        }

        private async void CopyVm(VM vm)
        {
            if (!VmRules.CanCopy(vm) && !VmRules.CanCopyTemplate(vm)) return;
            var action = await Dialogs.CopyVm(this, vm);
            if (action != null) RunTracked(action);
        }

        private async void DeleteVm(VM vm)
        {
            if (!VmRules.CanDeleteVm(vm)) return;
            var action = await Dialogs.DeleteVm(this, vm);
            if (action != null) RunTracked(action);
        }

        // ---------- Properties (PropertiesCommand) ----------

        private async void ShowProperties(IXenObject o) => await ShowProperties(this, o);

        private async System.Threading.Tasks.Task ShowProperties(Window owner, IXenObject o)
        {
            if (!VmRules.CanEditProperties(o)) return;
            var dialog = new PropertiesDialog(o);
            await dialog.ShowDialog(owner);
            if (dialog.Action == null) return;
            o.Locked = true; // SaveChangesAction requires the original locked
            RunTracked(dialog.Action, o);
        }

        /// <summary>Runs an action, mirroring its progress in Notifications/status and reporting failures.</summary>
        private void RunTracked(AsyncAction action) => RunTracked(action, null);

        private void RunTracked(AsyncAction action, IXenObject unlockWhenDone)
        {
            var index = _events.Count;
            _events.Add($"{action.Title}: starting...");
            void Update(ActionBase a) => Dispatcher.UIThread.Post(() =>
            {
                _events[index] = a.IsCompleted
                    ? $"{a.Title}: {(a.Succeeded ? "completed" : "failed - " + a.Exception?.Message)}"
                    : $"{a.Title}: {a.PercentComplete}% {a.Description}";
                SetStatus(_events[index]);
            });
            action.Changed += Update;
            action.Completed += a =>
            {
                if (unlockWhenDone != null) unlockWhenDone.Locked = false;
                Update(a);
                Dispatcher.UIThread.Post(() => { _tabsSignature = null; UpdateContent(); });
                if (!a.Succeeded && a.Exception != null)
                    Dispatcher.UIThread.Post(async () => await Dialogs.Info(this, "XCP-ng Center", $"{a.Title}\n\n{a.Exception.Message}"));
            };
            action.RunAsync();
        }

        /// <summary>
        /// File > Import: an XVA import, then "Start VM after import" as its own tracked action; or an OVF/OVA appliance
        /// (ImportApplianceAction starts its VMs itself when asked).
        /// </summary>
        private async void ImportVm()
        {
            var (import, start, appliance) = await Dialogs.ImportVm(this, _connections);
            if (appliance != null)
            {
                if (await Dialogs.ImportAppliance(this, _connections, appliance) is AsyncAction a) RunTracked(a);
                return;
            }
            if (import == null) return;
            if (start)
                import.Completed += _ => Dispatcher.UIThread.Post(() =>
                {
                    if (XvaRules.StartAfterImport(import) is AsyncAction a) RunTracked(a);
                });
            RunTracked(import);
        }

        // ---------- VM operations ----------

        private async void RunVmOp(VmOp op)
        {
            if (!(SelectedNode()?.XenObject is VM vm) || !VmOperations.IsAllowed(vm, op)) return;
            if (VmOperations.NeedsConfirmation(op))
            {
                var (title, text) = VmOperations.Confirmation(vm, op);
                if (!await Dialogs.Confirm(this, title, text, Messages.YES, Messages.NO))
                    return;
            }

            var action = VmOperations.Create(vm, op);
            var index = _events.Count;
            _events.Add($"{action.Title}: starting...");
            void Update(ActionBase a) => Dispatcher.UIThread.Post(() =>
            {
                _events[index] = a.IsCompleted
                    ? $"{a.Title}: {(a.Succeeded ? "completed" : "failed - " + a.Exception?.Message)}"
                    : $"{a.Title}: {a.PercentComplete}% {a.Description}";
                SetStatus(_events[index]);
            });
            action.Changed += Update;
            action.Completed += a =>
            {
                Update(a);
                if (!a.Succeeded && a.Exception != null)
                    Dispatcher.UIThread.Post(async () => await Dialogs.Info(this, "XCP-ng Center", $"{a.Title}\n\n{a.Exception.Message}"));
            };
            action.RunAsync();
        }

        private void AddEvent(string text)
        {
            _events.Add(text);
            SetStatus(text);
        }

        private void SetStatus(string text) => _status.Text = text;

        /// <summary>The red badge of the Notifications button (MainWindow.UpdateAlertCount).</summary>
        private void UpdateAlertCount()
        {
            var count = NotificationRules.Count(_connections);
            ((TextBlock)_alertBadge.Child).Text = count.ToString(CultureInfo.CurrentCulture);
            _alertBadge.IsVisible = count > 0;
            ((TextBlock)_statusAlerts.Children[1]).Text = ((TextBlock)_alertBadge.Child).Text;
            ((Control)_statusAlerts.Parent).IsVisible = count > 0;
        }

        /// <summary>NotificationsView: the left pane lists Alerts and Events; the right shows the chosen one.</summary>
        private void ShowNotifications()
        {
            if (_alertsView == null)
            {
                _alertsView = new NotificationsView(() => _connections, RunTracked);
                _alertsView.Changed += () => { if (_showEvents && _showAlerts) UpdateContent(); };
            }
            // the grid refreshes only while it is shown
            if (_showAlerts)
                _alertsView.Refresh();
            var alertsItem = $"Alerts ({(_showAlerts ? _alertsView.Total : NotificationRules.Count(_connections))})";
            if (_notificationList == null)
            {
                _notificationList = new ListBox { Margin = new Thickness(2) };
                _notificationList.SelectionChanged += (_, _) =>
                {
                    if (_notificationList.SelectedIndex < 0) return;
                    _showAlerts = _notificationList.SelectedIndex == 0;
                    UpdateContent();
                };
            }
            var items = new List<Control>();
            foreach (var (icon, text) in new[] { ("000_Alert2_h32bit_16", alertsItem), ("notif_events_16", "Events") })
            {
                var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(6, 6) };
                row.Children.Add(Views.Icon(icon, 16));
                row.Children.Add(new TextBlock { Text = text });
                items.Add(row);
            }
            var index = _showAlerts ? 0 : 1;
            if (_notificationList.ItemCount != 2 || !(_notificationList.Items[0] is StackPanel first && ((TextBlock)first.Children[1]).Text == alertsItem))
                _notificationList.ItemsSource = items;
            if (_notificationList.SelectedIndex != index)
                _notificationList.SelectedIndex = index;
            if (_leftPane.Child != _notificationList)
                _leftPane.Child = _notificationList;

            if (_showAlerts)
            {
                _titleBand.Content = Views.TitleBand("000_Alert2_h32bit_16", alertsItem, _alertsView.FiltersText);
                SetTabs("alerts", () => new[] { Tab("Alerts", _alertsView) });
            }
            else
            {
                _titleBand.Content = Views.TitleBand("notif_events_16", "Events", null);
                SetTabs("events", () =>
                {
                    var header = Views.TabHeader("Events");
                    DockPanel.SetDock(header, Dock.Top);
                    return new[] { Tab("Events", new DockPanel { Children = { header, new ListBox { ItemsSource = _events } } }) };
                });
            }
        }

        /// <summary>Asks about certificates rejected during the last handshake; true if the user trusted at least one.</summary>
        private async System.Threading.Tasks.Task<bool> ResolvePendingCertificates()
        {
            var accepted = false;
            foreach (var hostKey in _trust.Pending.Keys.ToList())
            {
                if (!_trust.Pending.TryRemove(hostKey, out var p)) continue;
                if (await Dialogs.TrustCertificate(this, hostKey, p.Fingerprint, p.Old))
                {
                    _trust.Pin(hostKey, p.Fingerprint);
                    accepted = true;
                }
            }
            return accepted;
        }
    }
}
