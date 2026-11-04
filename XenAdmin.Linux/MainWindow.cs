using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Controls.Templates;
using Avalonia.Layout;
using Avalonia.Styling;
using Avalonia.Threading;
using XenAdmin.Actions;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Linux
{
    public class MainWindow : Window
    {
        private readonly AppConfig _config;
        private readonly CertificateTrust _trust;
        private readonly List<XenConnection> _connections = new List<XenConnection>();
        private readonly TreeView _tree = new TreeView();
        private readonly TextBox _search = new TextBox { Watermark = "Buscar…", Width = 220 };
        private readonly Grid _props = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), ColumnSpacing = 12, RowSpacing = 4 };
        private readonly ObservableCollection<string> _tasks = new ObservableCollection<string>();
        private readonly TextBlock _status = new TextBlock { Margin = new Thickness(8, 4) };
        private readonly Dictionary<VmOp, Button> _opButtons = new Dictionary<VmOp, Button>();
        private readonly Button _connect = new Button { Content = "Conectar" };
        private readonly Button _disconnect = new Button { Content = "Desconectar" };
        private readonly Button _remove = new Button { Content = "Remover" };
        private readonly DispatcherTimer _refreshTimer;
        private string _selectedKey;
        private string _treeSignature;
        private List<TreeNode> _roots = new List<TreeNode>();
        private bool _refreshing;

        private static readonly (VmOp op, string text)[] OpTexts =
        {
            (VmOp.Start, "Iniciar"), (VmOp.Shutdown, "Desligar"), (VmOp.Reboot, "Reiniciar"),
            (VmOp.Suspend, "Suspender"), (VmOp.Resume, "Retomar"),
            (VmOp.ForceShutdown, "Forçar desligamento"), (VmOp.ForceReboot, "Forçar reinício"),
        };

        public MainWindow(AppConfig config, CertificateTrust trust)
        {
            _config = config;
            _trust = trust;
            Title = "XCP-ng Center";
            Width = 1100;
            Height = 700;
            try { Icon = new WindowIcon(Avalonia.Platform.AssetLoader.Open(new Uri("avares://xcp-ng-center/Assets/AppIcon.ico"))); } catch { }

            var add = new Button { Content = "Adicionar servidor" };
            add.Click += async (_, _) => await AddServer();
            _connect.Click += async (_, _) => { if (SelectedConnection() is XenConnection c) await Connect(c); };
            _disconnect.Click += (_, _) => { if (SelectedConnection() is XenConnection c) Disconnect(c); };
            _remove.Click += (_, _) => RemoveSelected();

            var toolbar = new WrapPanel { Margin = new Thickness(6) };
            foreach (var b in new[] { add, _connect, _disconnect, _remove })
                toolbar.Children.Add(Spaced(b));
            toolbar.Children.Add(new Border { Width = 12 });
            foreach (var (op, text) in OpTexts)
            {
                var b = new Button { Content = text, IsEnabled = false };
                var captured = op;
                b.Click += async (_, _) => await RunVmOp(captured);
                _opButtons[op] = b;
                toolbar.Children.Add(Spaced(b));
            }
            toolbar.Children.Add(Spaced(_search));
            _search.TextChanged += (_, _) => RefreshTree();

            _tree.ItemTemplate = new FuncTreeDataTemplate<TreeNode>((n, _) =>
            {
                var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
                row.Children.Add(new Avalonia.Controls.Shapes.Ellipse { Width = 10, Height = 10, Fill = Avalonia.Media.Brush.Parse(n.Color), VerticalAlignment = VerticalAlignment.Center });
                row.Children.Add(new TextBlock { Text = n.Name });
                row.Children.Add(new TextBlock { Text = n.Kind, Opacity = 0.6 });
                Avalonia.Automation.AutomationProperties.SetName(row, n.Label);
                return row;
            }, n => n.Children);
            _tree.Styles.Add(new Style(x => x.OfType<TreeViewItem>()) { Setters = { new Setter(TreeViewItem.IsExpandedProperty, true) } });
            _tree.SelectionChanged += (_, _) =>
            {
                // resetting ItemsSource raises a null selection; keep the user's choice
                if (_refreshing || !(_tree.SelectedItem is TreeNode n)) return;
                _selectedKey = n.Key;
                UpdateSelection(n);
            };

            var tasksList = new ListBox { ItemsSource = _tasks, Height = 140 };
            var right = new DockPanel { Margin = new Thickness(8) };
            var tasksHeader = new TextBlock { Text = "Tarefas", FontWeight = Avalonia.Media.FontWeight.Bold, Margin = new Thickness(0, 8, 0, 4) };
            DockPanel.SetDock(tasksList, Dock.Bottom);
            DockPanel.SetDock(tasksHeader, Dock.Bottom);
            right.Children.Add(tasksList);
            right.Children.Add(tasksHeader);
            right.Children.Add(new ScrollViewer { Content = _props });

            var split = new Grid { ColumnDefinitions = new ColumnDefinitions("350,4,*") };
            Grid.SetColumn(_tree, 0);
            var splitter = new GridSplitter { ResizeDirection = GridResizeDirection.Columns };
            Grid.SetColumn(splitter, 1);
            Grid.SetColumn(right, 2);
            split.Children.Add(_tree);
            split.Children.Add(splitter);
            split.Children.Add(right);

            var root = new DockPanel();
            DockPanel.SetDock(toolbar, Dock.Top);
            DockPanel.SetDock(_status, Dock.Bottom);
            root.Children.Add(toolbar);
            root.Children.Add(_status);
            root.Children.Add(split);
            Content = root;

            // coalesce bursts of cache events into one tree rebuild
            _refreshTimer = new DispatcherTimer { Interval = TimeSpan.FromMilliseconds(500) };
            _refreshTimer.Tick += (_, _) => { _refreshTimer.Stop(); RefreshTree(); };

            foreach (var s in _config.Servers)
                Track(new XenConnection { Hostname = s.Host, Port = s.Port, Username = s.User });
            RefreshTree();
            Closing += (_, _) =>
            {
                ((LinuxConfigProvider)XenAdminConfigManager.Provider).Exiting = true;
                foreach (var c in _connections.Where(c => c.IsConnected))
                    c.EndConnect(true, true);
            };
        }

        private static Control Spaced(Control c) { c.Margin = new Thickness(2); return c; }

        private void Track(XenConnection c)
        {
            _connections.Add(c);
            lock (ConnectionsManager.ConnectionsLock)
                ConnectionsManager.XenConnections.Add(c);
            void Schedule() => Dispatcher.UIThread.Post(() => { if (!_refreshTimer.IsEnabled) _refreshTimer.Start(); });
            c.CachePopulated += _ => Schedule();
            c.ConnectionStateChanged += _ => Schedule();
            c.ConnectionLost += _ => Dispatcher.UIThread.Post(() => SetStatus($"Conexão com {c.Hostname} perdida; reconectando…"));
            c.XenObjectsUpdated += (_, _) => Schedule();
            c.ConnectionMessageChanged += (_, m) => Dispatcher.UIThread.Post(() => SetStatus($"{c.Hostname}: {m}"));
            c.ConnectionResult += (_, e) => Dispatcher.UIThread.Post(async () =>
            {
                Schedule();
                if (e.Connected)
                {
                    SetStatus($"Conectado a {c.Hostname} (API {c.Session?.APIVersion})");
                    _config.AddOrUpdateServer(c.Hostname, c.Port, c.Username);
                }
                else if (await ResolvePendingCertificates())
                {
                    SetStatus($"Certificado aceito; reconectando a {c.Hostname}…");
                    c.BeginConnect(false, (_, _) => false);
                }
                else
                {
                    c.Password = null; // ask again next time
                    SetStatus($"Falha ao conectar em {c.Hostname}");
                    await Dialogs.Info(this, "Falha na conexão", e.Reason ?? e.Error?.Message ?? "Erro desconhecido");
                }
            });
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
            SetStatus($"Conectando a {c.Hostname}…");
            c.BeginConnect(false, (_, _) => false);
            RefreshTree();
        }

        private async System.Threading.Tasks.Task Connect(XenConnection c)
        {
            if (c.IsConnected || c.InProgress) return;
            if (string.IsNullOrEmpty(c.Password))
            {
                var input = await Dialogs.AddServer(this, c.Port == 443 ? c.Hostname : $"{c.Hostname}:{c.Port}", c.Username ?? "root");
                if (input == null) return;
                c.Username = input.User;
                c.Password = input.Password;
            }
            SetStatus($"Conectando a {c.Hostname}…");
            c.BeginConnect(false, (_, _) => false);
        }

        private void Disconnect(XenConnection c)
        {
            if (c.IsConnected || c.InProgress)
                c.EndConnect();
            SetStatus($"Desconectado de {c.Hostname}");
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
            RefreshTree();
        }

        private XenConnection SelectedConnection() => (_tree.SelectedItem as TreeNode)?.Connection as XenConnection;

        private void RefreshTree()
        {
            var roots = _connections.Select(ResourceTree.Build)
                .Select(n => ResourceTree.Filter(n, _search.Text)).Where(n => n != null).ToList();
            // cache objects are updated in place, so only rebuild the tree when its shape or labels change
            var signature = string.Join("\n", Flatten(roots).Select(n => $"{n.Key}|{n.Label}|{n.Color}|{n.Children.Count}"));
            if (signature != _treeSignature)
            {
                _treeSignature = signature;
                _roots = roots;
                _refreshing = true;
                try { _tree.ItemsSource = roots; }
                finally { _refreshing = false; }
                var match = Flatten(roots).FirstOrDefault(n => n.Key == _selectedKey);
                if (match != null)
                    Dispatcher.UIThread.Post(() => { _refreshing = true; _tree.SelectedItem = match; _refreshing = false; }, DispatcherPriority.Loaded);
            }
            UpdateSelection(Flatten(_roots).FirstOrDefault(n => n.Key == _selectedKey));
        }

        private static IEnumerable<TreeNode> Flatten(IEnumerable<TreeNode> nodes) =>
            nodes.SelectMany(n => new[] { n }.Concat(Flatten(n.Children)));

        private void UpdateSelection(TreeNode node)
        {
            var vm = node?.XenObject as VM;
            foreach (var (op, b) in _opButtons)
                b.IsEnabled = VmOperations.IsAllowed(vm, op);
            var c = node?.Connection;
            _connect.IsEnabled = c != null && !c.IsConnected && !c.InProgress;
            _disconnect.IsEnabled = c != null && (c.IsConnected || c.InProgress);
            _remove.IsEnabled = c != null && node.XenObject == null;

            _props.Children.Clear();
            _props.RowDefinitions.Clear();
            var rows = node?.XenObject != null
                ? ResourceTree.Properties(node.XenObject).ToList()
                : c != null ? new List<(string, string)> { ("Servidor", c.HostnameWithPort), ("Usuário", c.Username), ("Estado", c.IsConnected ? "conectado" : c.InProgress ? "conectando" : "desconectado") }
                : new List<(string, string)>();
            for (int i = 0; i < rows.Count; i++)
            {
                _props.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var k = new TextBlock { Text = rows[i].Item1, FontWeight = Avalonia.Media.FontWeight.SemiBold };
                var v = new SelectableTextBlock { Text = rows[i].Item2 ?? "", TextWrapping = Avalonia.Media.TextWrapping.Wrap };
                Grid.SetRow(k, i); Grid.SetRow(v, i); Grid.SetColumn(v, 1);
                _props.Children.Add(k);
                _props.Children.Add(v);
            }
        }

        private async System.Threading.Tasks.Task RunVmOp(VmOp op)
        {
            if (!((_tree.SelectedItem as TreeNode)?.XenObject is VM vm) || !VmOperations.IsAllowed(vm, op)) return;
            if (VmOperations.NeedsConfirmation(op) &&
                !await Dialogs.Confirm(this, "Confirmar", $"{OpTexts.First(o => o.op == op).text} da VM \"{vm.Name()}\"? Dados não salvos podem ser perdidos.", "Confirmar"))
                return;

            var action = VmOperations.Create(vm, op);
            var index = _tasks.Count;
            _tasks.Add($"{action.Title}: iniciando…");
            void Update(ActionBase a) => Dispatcher.UIThread.Post(() =>
            {
                _tasks[index] = a.IsCompleted
                    ? $"{a.Title}: {(a.Succeeded ? "concluído" : "falhou — " + a.Exception?.Message)}"
                    : $"{a.Title}: {a.PercentComplete}% {a.Description}";
            });
            action.Changed += Update;
            action.Completed += a =>
            {
                Update(a);
                if (!a.Succeeded && a.Exception != null)
                    Dispatcher.UIThread.Post(async () => await Dialogs.Info(this, "Falha na operação", $"{a.Title}\n\n{a.Exception.Message}"));
            };
            action.RunAsync();
        }

        private void SetStatus(string text) => _status.Text = text;

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
