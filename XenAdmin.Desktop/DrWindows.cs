using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using Avalonia.Threading;
using XenAdmin.Actions;
using XenAdmin.Actions.DR;
using XenAdmin.Network;
using XenAPI;
using Task = System.Threading.Tasks.Task;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// DRConfigureDialog: the SRs that can hold the pool's VM metadata (found in the background), ticked where it is
    /// stored already; OK enables or disables database replication where the choice changed.
    /// </summary>
    public static class DrConfigureDialog
    {
        public static async System.Threading.Tasks.Task<AsyncAction> Show(Window owner, Pool pool)
        {
            var list = new StackPanel { Spacing = 4 };
            var warning = new TextBlock { TextWrapping = TextWrapping.Wrap, Foreground = Brushes.DarkGoldenrod };
            var ok = new Button { Content = "OK", IsDefault = true, MinWidth = 75, IsEnabled = false };
            var cancel = new Button { Content = Dialogs.M(Messages.CANCEL), IsCancel = true, MinWidth = 75 };
            var panel = new StackPanel { Spacing = 10 };
            panel.Children.Add(new TextBlock { Text = "The configuration data needed to recover VMs and vApps in the event of a disaster must be stored on one or more SRs in this pool.", TextWrapping = TextWrapping.Wrap });
            panel.Children.Add(new TextBlock { Text = "Select the SRs where the VM configuration data will be stored:" });
            panel.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Padding = new Thickness(6), MinHeight = 160, Background = Brushes.White,
                Child = new ScrollViewer { Content = list, MaxHeight = 260 } });
            panel.Children.Add(warning);
            panel.Children.Add(Dialogs.Buttons(ok, cancel));
            var w = Dialogs.Frame(string.Format(Messages.DR_CONFIGURE_TITLE, pool.Name()), panel, 560);

            var rows = new List<(SR Sr, bool Enabled, CheckBox Box)>();
            var originallyEnabled = false;
            void Update()
            {
                var (text, canSave) = DrRules.ConfigureWarning(originallyEnabled, rows.Count(r => r.Box.IsChecked == true));
                warning.Text = text ?? "";
                ok.IsEnabled = canSave;
            }
            list.Children.Add(new TextBlock { Text = "Looking for the SRs that can store the VM configuration data...", Opacity = 0.7 });
            w.Opened += async (_, _) =>
            {
                List<SR> srs;
                try { srs = await Task.Run(() => DrRules.ReplicationSrs(pool.Connection)); }
                catch (Exception e) { list.Children.Clear(); list.Children.Add(new TextBlock { Text = e.Message, Foreground = Brushes.Red }); return; }
                list.Children.Clear();
                foreach (var sr in srs)
                {
                    var enabled = DrRules.IsDrEnabled(sr, pool);
                    originallyEnabled |= enabled;
                    var box = new CheckBox { Content = $"{sr.Name()}   ({sr.GetSRType(false)}, {XenAdmin.Util.DiskSizeString(sr.FreeSpace())} free)", IsChecked = enabled };
                    box.IsCheckedChanged += (_, _) => Update();
                    rows.Add((sr, enabled, box));
                    list.Children.Add(box);
                }
                if (srs.Count == 0)
                    list.Children.Add(new TextBlock { Text = "There are no SRs in this pool that can store the VM configuration data.", TextWrapping = TextWrapping.Wrap });
                Update();
            };
            AsyncAction result = null;
            ok.Click += (_, _) =>
            {
                result = DrRules.ConfigureAction(pool,
                    rows.Where(r => !r.Enabled && r.Box.IsChecked == true).Select(r => r.Sr),
                    rows.Where(r => r.Enabled && r.Box.IsChecked != true).Select(r => r.Sr));
                w.Close();
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }
    }

    /// <summary>
    /// DRFailoverWizard: Before You Start (failover, failback or test failover), Storage (the SRs with metadata),
    /// VMs and vApps (from each pool's metadata, and what to do with them after recovery), Prechecks, Summary; Finish
    /// recovers them, and a test failover then offers to remove what it recovered.
    /// </summary>
    public class DrWizard : WizardWindow
    {
        private readonly IXenConnection _c;
        private readonly Action<AsyncAction> _run;
        private DrType _type = DrType.Failover;
        private List<(SR Sr, bool HasMetadata)> _srs;
        private readonly HashSet<string> _chosenSrs = new HashSet<string>();
        private List<PoolMetadata> _metadata;
        private string _metadataFor;
        private readonly HashSet<string> _chosenObjects = new HashSet<string>();
        private DrRules.StartMode _start = DrRules.StartMode.Start;
        private List<(string Text, bool Error)> _checks;
        private int _metadataGeneration, _checksGeneration; // a newer request makes an older answer stale

        public DrWizard(IXenConnection c, Action<AsyncAction> run)
            : base(string.Format("Disaster Recovery - {0}", XenAdmin.Core.Helpers.GetName(XenAdmin.Core.Helpers.GetPoolOfOne(c))), "000_Storage_h32bit_32",
                "Before You Start", "Storage", "VMs and vApps", "Prechecks", "Summary")
        {
            _c = c;
            _run = run;
        }

        private string Verb => _type switch { DrType.Failback => "failback", DrType.Test => "test failover", _ => "failover" };

        protected override string PageTitle(int step) => step switch
        {
            0 => "Disaster Recovery Wizard",
            1 => "Locate the mirrored storage",
            2 => "Select the vApps and VMs to recover",
            3 => "Check that the selected vApps and VMs can be recovered",
            _ => $"Ready to start the {Verb}"
        };

        private static TextBlock Text(string s, double opacity = 1) => new TextBlock { Text = s, TextWrapping = TextWrapping.Wrap, Opacity = opacity };

        protected override Control Page(int step)
        {
            switch (step)
            {
                case 0: return TypePage();
                case 1: return StoragePage();
                case 2: return ObjectsPage();
                case 3: return ChecksPage();
                default: return SummaryPage();
            }
        }

        protected override string Validate(int step) => step switch
        {
            1 when _srs == null => "Looking for the SRs...",
            1 when _chosenSrs.Count == 0 => "Select the SRs holding the metadata of the pool to recover.",
            2 when _metadata == null => "Loading the metadata...",
            2 when Chosen().Sum(x => x.Vapps.Count + x.Vms.Count) == 0 => "Select the vApps and VMs to recover.",
            3 when _checks == null => "Running the prechecks...",
            3 when _checks.Any(c => c.Error) => "Fix the errors above before you continue.",
            _ => null
        };

        // ---------- pages ----------

        private Control TypePage()
        {
            var p = new StackPanel { Spacing = 8 };
            p.Children.Add(Text("This wizard will take you through the steps needed to fail over your critical vApps and VMs in the event of a site/datacenter failure at your primary data site, to fail them back, or to test a failover."));
            p.Children.Add(Text("Configuration information for the vApps and VMs to be recovered will be retrieved from remote storage mirrors. Before the recovery begins, the wizard will run a number of checks to ensure that the selected vApps and VMs can be brought back on this pool.", 0.85));
            RadioButton Option(DrType t, string label, string desc)
            {
                var rb = new RadioButton { Content = Dialogs.M(label), GroupName = "drtype", IsChecked = _type == t, Margin = new Thickness(0, 8, 0, 0) };
                rb.IsCheckedChanged += (_, _) => { if (rb.IsChecked == true) _type = t; };
                p.Children.Add(rb);
                p.Children.Add(new TextBlock { Text = desc, TextWrapping = TextWrapping.Wrap, Margin = new Thickness(22, 0, 0, 0), Opacity = 0.8 });
                return rb;
            }
            Option(DrType.Failover, "&Failover", "Recover the vApps and VMs of a failed primary site on this (secondary) pool.");
            Option(DrType.Failback, "Fail&back", "Bring the vApps and VMs back to the primary site once it has been restored.");
            Option(DrType.Test, "&Test Failover", "Recover the vApps and VMs here to check the failover works; they are removed when the wizard ends.");
            p.Children.Add(new TextBlock { Text = "IMPORTANT", FontWeight = FontWeight.Bold, Margin = new Thickness(0, 10, 0, 0) });
            p.Children.Add(Text("Before failover can take place, storage mirroring - remote replication of your VM data - must be halted. The mirrored SRs must be attached to this pool."));
            return p;
        }

        private Control StoragePage()
        {
            var list = new StackPanel { Spacing = 4 };
            var p = new StackPanel { Spacing = 8 };
            p.Children.Add(Text("Select the SRs holding the VM configuration data (metadata) of the pool you want to recover:"));
            p.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Padding = new Thickness(6), MinHeight = 180, Background = Brushes.White,
                Child = new ScrollViewer { Content = list, MaxHeight = 280 } });
            void Fill()
            {
                list.Children.Clear();
                foreach (var (sr, has) in _srs)
                {
                    var box = new CheckBox { Content = $"{sr.Name()}   ({sr.GetSRType(false)}){(has ? "   - metadata detected" : "")}", IsChecked = _chosenSrs.Contains(sr.uuid) };
                    box.IsCheckedChanged += (_, _) => { if (box.IsChecked == true) _chosenSrs.Add(sr.uuid); else _chosenSrs.Remove(sr.uuid); };
                    list.Children.Add(box);
                }
                if (_srs.Count == 0)
                    list.Children.Add(Text("No SR attached to this pool can hold VM metadata. Attach the mirrored SR with New SR first."));
            }
            if (_srs != null) Fill();
            else
            {
                list.Children.Add(Text("Looking for the SRs...", 0.7));
                Task.Run(() => DrRules.StorageChoices(_c)).ContinueWith(t => Dispatcher.UIThread.Post(() =>
                {
                    _srs = t.IsFaulted ? new List<(SR, bool)>() : t.Result;
                    foreach (var (sr, has) in _srs.Where(x => x.HasMetadata)) _chosenSrs.Add(sr.uuid);
                    ShowError(t.IsFaulted ? t.Exception?.InnerException?.Message : null);
                    Fill();
                }));
            }
            return p;
        }

        private Control ObjectsPage()
        {
            var key = string.Join(",", _chosenSrs.OrderBy(x => x));
            var list = new StackPanel { Spacing = 3 };
            var p = new StackPanel { Spacing = 8 };
            p.Children.Add(Text("Select the vApps and VMs to recover. A vApp brings its own VMs."));
            p.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Padding = new Thickness(6), MinHeight = 200, Background = Brushes.White,
                Child = new ScrollViewer { Content = list, MaxHeight = 280 } });
            // DRFailoverWizardAppliancesPage: a test failover never boots the recovered guests (paused or not started)
            var modes = _type == DrType.Test
                ? new List<(DrRules.StartMode Mode, string Text)> { (DrRules.StartMode.StartPaused, "Start them paused"), (DrRules.StartMode.None, "Do not start them") }
                : new List<(DrRules.StartMode Mode, string Text)> { (DrRules.StartMode.Start, "Start the recovered vApps and VMs"), (DrRules.StartMode.StartPaused, "Start them paused"), (DrRules.StartMode.None, "Do not start them") };
            if (modes.All(m => m.Mode != _start)) _start = modes[0].Mode;
            var power = new ComboBox { ItemsSource = modes.Select(m => m.Text).ToList(), SelectedIndex = modes.FindIndex(m => m.Mode == _start), Width = 280 };
            power.SelectionChanged += (_, _) => { if (power.SelectedIndex >= 0) _start = modes[power.SelectedIndex].Mode; };
            var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            row.Children.Add(new TextBlock { Text = "After recovery:", VerticalAlignment = VerticalAlignment.Center });
            row.Children.Add(power);
            p.Children.Add(row);
            void Fill()
            {
                list.Children.Clear();
                foreach (var md in _metadata)
                {
                    var current = DrRules.IsCurrentPool(md, _c);
                    list.Children.Add(new TextBlock { Text = md.Pool.Name() + (current ? $"   ({Messages.DR_WIZARD_APPLIANCESPAGE_CURRENT_POOL})" : ""), FontWeight = FontWeight.SemiBold, Margin = new Thickness(0, 6, 0, 0) });
                    void Item(IXenObject o, string text)
                    {
                        var box = new CheckBox { Content = text, IsEnabled = !current, IsChecked = !current && _chosenObjects.Contains(o.opaque_ref), Margin = new Thickness(18, 0, 0, 0) };
                        box.IsCheckedChanged += (_, _) => { if (box.IsChecked == true) _chosenObjects.Add(o.opaque_ref); else _chosenObjects.Remove(o.opaque_ref); };
                        list.Children.Add(box);
                    }
                    foreach (var app in md.VmAppliances.Values.OrderBy(a => a.Name()))
                        Item(app, $"vApp '{app.Name()}' ({app.VMs.Count} VMs)");
                    foreach (var vm in DrRules.LooseVms(md))
                        Item(vm, vm.Name());
                }
                if (_metadata.Count == 0)
                    list.Children.Add(Text("No pool metadata was found on the selected SRs."));
            }
            if (_metadata != null && _metadataFor == key) Fill();
            else
            {
                _metadata = null;
                _chosenObjects.Clear();
                list.Children.Add(Text("Loading the metadata...", 0.7));
                var srs = _srs.Where(x => _chosenSrs.Contains(x.Sr.uuid)).Select(x => x.Sr).ToList();
                var generation = ++_metadataGeneration;
                Task.Run(() => DrRules.LoadMetadata(_c, srs)).ContinueWith(t => Dispatcher.UIThread.Post(() =>
                {
                    if (generation != _metadataGeneration) return;
                    _metadata = t.IsFaulted ? new List<PoolMetadata>() : t.Result;
                    _metadataFor = key;
                    ShowError(t.IsFaulted ? t.Exception?.InnerException?.Message : null);
                    Fill();
                }));
            }
            return p;
        }

        private List<(PoolMetadata Metadata, List<VM_appliance> Vapps, List<VM> Vms)> Chosen() =>
            (_metadata ?? new List<PoolMetadata>()).Where(md => !DrRules.IsCurrentPool(md, _c))
                .Select(md => (md, md.VmAppliances.Values.Where(a => _chosenObjects.Contains(a.opaque_ref)).ToList(),
                               DrRules.LooseVms(md).Where(v => _chosenObjects.Contains(v.opaque_ref)).ToList()))
                .Where(x => x.Item2.Count + x.Item3.Count > 0).ToList();

        private Control ChecksPage()
        {
            var list = new StackPanel { Spacing = 4 };
            var p = new StackPanel { Spacing = 8 };
            p.Children.Add(Text($"Checking that the {Verb} can take place:"));
            p.Children.Add(new Border { BorderBrush = Brush.Parse("#828790"), BorderThickness = new Thickness(1), Padding = new Thickness(6), MinHeight = 160, Background = Brushes.White,
                Child = new ScrollViewer { Content = list, MaxHeight = 280 } });
            void Fill()
            {
                list.Children.Clear();
                foreach (var (text, error) in _checks)
                {
                    var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
                    row.Children.Add(Views.Icon(error ? "000_Abort_h32bit_16" : "000_Alert2_h32bit_16"));
                    row.Children.Add(new TextBlock { Text = text, TextWrapping = TextWrapping.Wrap, MaxWidth = 520 });
                    list.Children.Add(row);
                }
                if (_checks.Count == 0)
                {
                    var row = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6 };
                    row.Children.Add(Views.Icon("000_Tick_h32bit_16"));
                    row.Children.Add(new TextBlock { Text = "All checks passed." });
                    list.Children.Add(row);
                }
                ShowError(_checks.Any(c => c.Error) ? "Fix the errors above before you continue." : null);
            }
            // the checks run again each time the page is shown: the user may have fixed what failed
            _checks = null;
            list.Children.Add(Text("Running the prechecks...", 0.7));
            var chosen = Chosen();
            var type = _type;
            var generation = ++_checksGeneration;
            Task.Run(() => DrRules.Prechecks(_c, type, chosen)).ContinueWith(t => Dispatcher.UIThread.Post(() =>
            {
                if (generation != _checksGeneration) return;
                _checks = t.IsFaulted ? new List<(string, bool)> { (t.Exception?.InnerException?.Message ?? "The prechecks failed.", true) } : t.Result;
                Fill();
            }));
            return p;
        }

        private Control SummaryPage()
        {
            var chosen = Chosen();
            var p = new StackPanel { Spacing = 8 };
            p.Children.Add(Text($"Press Finish to start the {Verb}. These will be recovered on '{XenAdmin.Core.Helpers.GetName(_c)}':"));
            foreach (var (md, vapps, vms) in chosen)
            {
                p.Children.Add(new TextBlock { Text = "From " + md.Pool.Name(), FontWeight = FontWeight.SemiBold });
                foreach (var a in vapps) p.Children.Add(new TextBlock { Text = $"    vApp '{a.Name()}'" });
                foreach (var v in vms) p.Children.Add(new TextBlock { Text = "    " + v.Name() });
            }
            p.Children.Add(Text(_start switch
            {
                DrRules.StartMode.Start => "They will be started after recovery.",
                DrRules.StartMode.StartPaused => "They will be started paused after recovery.",
                _ => "They will not be started."
            }, 0.8));
            if (_type == DrType.Test)
                p.Children.Add(Text("This is a test: when it is done, you will be asked to remove what it recovered.", 0.8));
            return p;
        }

        protected override Task Finish()
        {
            var chosen = Chosen();
            if (chosen.Count == 0) return Task.CompletedTask;
            var recover = new DrRules.RecoverAction(_c, chosen, _start);
            if (_type == DrType.Test)
            {
                var owner = Owner as Window;
                recover.Completed += _ => Dispatcher.UIThread.Post(async () =>
                {
                    // DRFailoverWizard.FinishWizard (test failover): remove what the test recovered
                    if (DrRules.TestCleanupAction(_c, recover) is AsyncAction cleanup
                        && await Dialogs.Confirm(owner, "Test Failover", "The test failover is done. Remove the vApps and VMs it recovered?", Messages.YES, Messages.NO))
                        _run(cleanup);
                });
            }
            _run(recover);
            Close();
            return Task.CompletedTask;
        }
    }
}
