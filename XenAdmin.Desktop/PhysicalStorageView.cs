using System;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using XenAdmin.Actions;
using XenAPI;
using Task = System.Threading.Tasks.Task;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// Host and pool Storage tab (PhysicalStoragePage): the storage repositories and New SR / Reclaim freed space /
    /// Properties (NewSRCommand, TrimSRCommand with its confirmation, PropertiesDialog).
    /// </summary>
    public class PhysicalStorageView : DockPanel, ILiveView
    {
        // the selected SR per host or pool: the tab is rebuilt after each action
        private static readonly System.Collections.Generic.Dictionary<string, string> SelectedSr = new System.Collections.Generic.Dictionary<string, string>();

        private readonly IXenObject _owner;
        private readonly Action<AsyncAction> _run;
        private readonly Action<IXenObject> _properties;
        private readonly SelectableGrid _grid = new SelectableGrid();
        private readonly Button _newSr, _trim, _props;

        public PhysicalStorageView(IXenObject hostOrPool, Action<AsyncAction> run, Action<IXenObject> properties, Action newSr)
        {
            _owner = hostOrPool;
            _run = run;
            _properties = properties;
            Background = Brushes.White;
            Button B(string text) => new Button { Content = Dialogs.M(text), MinWidth = 94 };
            _newSr = B("&New SR...");
            _trim = B("Re&claim freed space");
            _props = B("P&roperties");
            ToolTip.SetShowOnDisabled(_trim, true);
            _newSr.Click += (_, _) => newSr();
            _trim.Click += async (_, _) => await Trim();
            _props.Click += (_, _) => OpenProperties();
            _grid.SelectionChanged += () => { SelectedSr[_owner.opaque_ref] = _grid.Selected; UpdateButtons(); };
            _grid.Activated += OpenProperties;

            var bar = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(0, 12, 0, 0) };
            bar.Children.Add(_newSr);
            bar.Children.Add(_trim);
            bar.Children.Add(_props);
            var grid = new Grid { Margin = new Thickness(12), RowDefinitions = new RowDefinitions("3*,Auto,*"), MaxWidth = 924, HorizontalAlignment = HorizontalAlignment.Left };
            Grid.SetRow(bar, 1);
            grid.Children.Add(_grid);
            grid.Children.Add(bar);
            Children.Add(grid);
            Refresh();
        }

        private SR Selected => StorageRules.PhysicalSrs(_owner).FirstOrDefault(s => s.opaque_ref == _grid.Selected);

        public void Refresh()
        {
            SelectedSr.TryGetValue(_owner.opaque_ref, out var selected);
            _grid.Fill(Details.Storage(_owner), StorageRules.PhysicalSrs(_owner).Select(s => s.opaque_ref).ToList(), selected);
            SelectedSr[_owner.opaque_ref] = _grid.Selected;
            UpdateButtons();
        }

        private void UpdateButtons()
        {
            var sr = Selected;
            _newSr.IsEnabled = _owner.Connection.IsConnected && !_owner.Locked;
            _trim.IsEnabled = StorageRules.CanTrim(sr);
            ToolTip.SetTip(_trim, sr != null && !sr.SupportsTrim() ? Messages.TOOLTIP_SR_TRIM_UNSUPPORTED : null);
            _props.IsEnabled = sr != null;
        }

        private void OpenProperties()
        {
            if (Selected is SR sr) _properties(sr);
        }

        private async Task Trim()
        {
            var sr = Selected;
            if (!StorageRules.CanTrim(sr)) return;
            if (await Dialogs.Confirm(TopLevel.GetTopLevel(this) as Window, Messages.CONFIRM_TRIM_SR_TITLE, Messages.CONFIRM_TRIM_SR, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION))
                _run(new SrTrimAction(sr.Connection, sr));
        }
    }
}
