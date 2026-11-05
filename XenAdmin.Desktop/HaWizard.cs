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
    /// HAWizard: Prerequisites (Intro), Heartbeat SR (ChooseSR, scanned with GetHeartbeatSRsAction), HA Plan
    /// (AssignPriorities) and Finish (HAFinishPage), then EnableHAAction.
    /// </summary>
    public class HaWizard : WizardWindow
    {
        private readonly Pool _pool;
        private readonly Action<AsyncAction> _run;
        private List<SRWrapper> _srs;
        private SR _heartbeat;
        private HaPlanEditor _plan;

        public HaWizard(Pool pool, Action<AsyncAction> run)
            : base(string.Format("Configure HA on '{0}'", XenAdmin.Core.Helpers.GetName(pool)), "000_HAServer_h32bit_32",
                Messages.HA_INTRO_TEXT, Messages.HA_CHOOSESR_PAGE_TEXT, "HA Plan", Messages.FINISH_PAGE_TEXT)
        {
            _pool = pool;
            _run = run;
        }

        protected override string PageTitle(int step) => step switch
        {
            0 => Messages.HA_INTRO_PAGE_TITLE,
            1 => Messages.HA_CHOOSESR_PAGE_PAGETITLE,
            2 => "Configure the HA restart priority, restart order and delay interval for the VMs in this pool",
            _ => Messages.HA_WIZARD_FINISH_PAGE_TITLE
        };

        protected override Control Page(int step) => step switch
        {
            0 => new TextBlock
            {
                TextWrapping = TextWrapping.Wrap,
                Text = "This wizard helps you configure HA (high availability) on a resource pool.\n\nBefore you begin, ensure that the following requirements are satisfied for all servers and virtual machines in the pool:\n\n•    Shared storage must be available, including at least one iSCSI, NFS or Fibre Channel LUN. This LUN will be used for the heartbeat SR.\n\n•    All the virtual machines you want to protect with HA must be agile.\n\nWe strongly recommend the use of a bonded management interface on the servers in the pool if HA is enabled, and multipathed storage for the heartbeat SR.\n\n\nTo continue, click Next."
            },
            1 => SrPage(),
            2 => _plan ??= Plan(),
            _ => Summary()
        };

        private HaPlanEditor Plan()
        {
            var plan = new HaPlanEditor(_pool);
            plan.Changed += () => ShowError(plan.Error ?? "");
            return plan;
        }

        protected override string Validate(int step) => step switch
        {
            1 when _heartbeat == null => "Select a heartbeat SR.",
            2 when _plan?.Error is string e => e,
            _ => null
        };

        protected override Task Finish()
        {
            if (_plan?.Error != null || _heartbeat == null) return Task.CompletedTask;
            _run(HaRules.EnableAction(_pool, _plan.Options(), _heartbeat, _plan.Ntol));
            Close();
            return Task.CompletedTask;
        }

        /// <summary>ChooseSR: the shared SRs, each usable for the heartbeat or with the reason it is not.</summary>
        private Control SrPage()
        {
            var panel = new StackPanel { Spacing = 10 };
            panel.Children.Add(new TextBlock { Text = "Select a heartbeat SR that will be used to monitor the availability and health of servers in the pool. Heartbeat SRs must use shared storage.", TextWrapping = TextWrapping.Wrap });
            var list = new ListBox { Height = 240 };
            var rescan = new Button { Content = Dialogs.M("&Rescan") };
            panel.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Child = list });
            panel.Children.Add(rescan);
            void Fill()
            {
                if (_srs == null)
                {
                    list.ItemsSource = new[] { Messages.HA_SCANNING_SRS.Replace("{0}", XenAdmin.Core.Helpers.GetName(_pool)) };
                    return;
                }
                list.ItemsSource = _srs.Select(w => new TextBlock { Text = w.enabled ? w.sr.Name() : $"{w.sr.Name()}   ({w.ReasonUnsuitable})", Opacity = w.enabled ? 1 : 0.5, Tag = w }).ToList();
                var current = _srs.FindIndex(w => w.enabled && w.sr.opaque_ref == _heartbeat?.opaque_ref);
                list.SelectedIndex = current >= 0 ? current : _srs.FindIndex(w => w.enabled);
            }
            list.SelectionChanged += (_, _) =>
            {
                if (list.SelectedItem is TextBlock { Tag: SRWrapper w })
                {
                    if (w.enabled) _heartbeat = w.sr;
                    else list.SelectedIndex = _srs.FindIndex(x => x.enabled && x.sr.opaque_ref == _heartbeat?.opaque_ref);
                }
            };
            void Scan()
            {
                // a new scan decides again which SRs are usable: nothing stays selected from the previous one
                _srs = null;
                _heartbeat = null;
                rescan.IsEnabled = false;
                Fill();
                var action = new GetHeartbeatSRsAction(_pool);
                action.Completed += _ => Dispatcher.UIThread.Post(() =>
                {
                    _srs = action.SRs.OrderBy(w => !w.enabled).ThenBy(w => w.sr.Name()).ToList();
                    rescan.IsEnabled = true;
                    Fill();
                    if (_srs.Count == 0) ShowError("There are no shared storage repositories in this pool.");
                });
                action.RunAsync();
            }
            rescan.Click += (_, _) => Scan();
            if (_srs == null) Scan(); else Fill();
            return panel;
        }

        /// <summary>HAFinishPage: the configuration summary.</summary>
        private Control Summary()
        {
            var counts = _plan.Counts();
            int Count(VM.HaRestartPriority p) => counts.TryGetValue(p, out var n) ? n : 0;
            var grid = new Grid { ColumnDefinitions = new ColumnDefinitions("Auto,*"), RowSpacing = 6, ColumnSpacing = 16 };
            void Row(string label, string value)
            {
                var r = grid.RowDefinitions.Count;
                grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                var l = new TextBlock { Text = label };
                var v = new TextBlock { Text = value };
                Grid.SetRow(l, r); Grid.SetRow(v, r); Grid.SetColumn(v, 1);
                grid.Children.Add(l); grid.Children.Add(v);
            }
            Row("Heartbeat SR:", _heartbeat?.Name() ?? "");
            Row("Failure Tolerance:", _plan.Ntol.ToString());
            Row("Restart:", $"{Count(VM.HaRestartPriority.Restart) + Count(VM.HaRestartPriority.AlwaysRestart) + Count(VM.HaRestartPriority.AlwaysRestartHighPriority)} VMs");
            Row("Restart if possible:", $"{Count(VM.HaRestartPriority.BestEffort)} VMs");
            Row("Do not restart:", $"{Count(VM.HaRestartPriority.DoNotRestart)} VMs");
            var panel = new StackPanel { Spacing = 12 };
            panel.Children.Add(new TextBlock { Text = "Review the final HA configuration below. You can press Finish to enable HA with this configuration, or Previous to alter the configuration.", TextWrapping = TextWrapping.Wrap });
            panel.Children.Add(new Border { BorderBrush = Brush.Parse("#DCDCDC"), BorderThickness = new Thickness(1), Padding = new Thickness(10), Child = grid });
            if (_plan.Ntol == 0)
                panel.Children.Add(new TextBlock { Text = "Under this configuration, the failover tolerance of your pool would be 0 (zero) and it may not be possible to restart your virtual machines in the event of a server failure.", TextWrapping = TextWrapping.Wrap, Foreground = Brushes.DarkGoldenrod });
            return panel;
        }
    }
}
