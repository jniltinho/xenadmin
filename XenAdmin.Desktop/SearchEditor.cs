using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;
using XenAdmin.Network;
using XenAdmin.XenSearch;
using Task = System.Threading.Tasks.Task;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// New Search / Edit Search (SearchPage with QueryPanel): name, the object types, a text filter and the grouping,
    /// then Apply (show it), Save to a pool (it joins Saved Searches) or, for a saved search, Delete.
    /// </summary>
    public static class SearchEditor
    {
        public sealed class Result
        {
            public Search Search;
            /// <summary>The pool to save it to (null: show only).</summary>
            public IXenConnection SaveTo;
            public bool Delete;
        }

        public static async System.Threading.Tasks.Task<Result> Show(Window owner, Search edited, IEnumerable<IXenConnection> connections)
        {
            var name = new TextBox { Text = edited == null ? "New Search" : edited.Name };
            var types = edited?.Query?.QueryScope?.ObjectTypes ?? Search.DefaultObjectTypes();
            var typeBoxes = SearchRules.Types().Select(t => (t.Type, Box: new CheckBox { Content = t.Name, IsChecked = (types & t.Type) == t.Type && t.Type != 0 })).ToList();
            var text = new TextBox { Watermark = edited?.Query?.QueryFilter != null ? "Add a text filter to the current conditions" : "Name, description or IP address contains..." };
            var groupings = SearchRules.Groupings();
            var groupBy = new ComboBox { ItemsSource = groupings.Select(g => g.Name).ToList(), SelectedIndex = SearchRules.GroupingIndex(edited), Width = 220 };
            var groupingChanged = false; // an unchanged choice keeps the search's own grouping (it may be nested or not listed)
            groupBy.SelectionChanged += (_, _) => groupingChanged = true;
            var pools = connections.Where(c => c.IsConnected && c.Cache.Pools.Length > 0).ToList();
            var saveTo = new ComboBox { ItemsSource = pools.Select(c => XenAdmin.Core.Helpers.GetName(c)).ToList(), SelectedIndex = pools.Count > 0 ? Math.Max(0, pools.IndexOf(edited?.Connection)) : -1, Width = 220, IsEnabled = pools.Count > 0 };
            var error = new TextBlock { Foreground = Brushes.Red, TextWrapping = TextWrapping.Wrap };

            var typeGrid = new WrapPanel { Orientation = Orientation.Horizontal };
            foreach (var (_, box) in typeBoxes)
            {
                box.Margin = new Thickness(0, 0, 14, 4);
                box.MinWidth = 190;
                typeGrid.Children.Add(box);
            }
            Grid Row(string label, Control field)
            {
                var g = new Grid { ColumnDefinitions = new ColumnDefinitions("110,*") };
                var l = new Label { Content = Dialogs.M(label), Target = field, VerticalAlignment = VerticalAlignment.Center, Padding = new Thickness(0) };
                Grid.SetColumn(field, 1);
                g.Children.Add(l);
                g.Children.Add(field);
                return g;
            }

            var apply = new Button { Content = Dialogs.M("&Apply"), IsDefault = true, MinWidth = 75 };
            var save = new Button { Content = Dialogs.M("&Save"), MinWidth = 75, IsEnabled = pools.Count > 0 };
            var delete = new Button { Content = Dialogs.M("&Delete"), MinWidth = 75, IsVisible = SearchRules.CanDelete(edited) };
            var cancel = new Button { Content = Dialogs.M(Messages.CANCEL), IsCancel = true, MinWidth = 75 };
            var panel = new StackPanel { Spacing = 10 };
            panel.Children.Add(Row("&Name:", name));
            panel.Children.Add(new TextBlock { Text = "Search for:", FontWeight = FontWeight.SemiBold });
            panel.Children.Add(typeGrid);
            panel.Children.Add(Row("&Filter:", text));
            if (edited?.Query?.QueryFilter != null)
                panel.Children.Add(new TextBlock { Text = "The search's current conditions are kept.", Opacity = 0.7 });
            panel.Children.Add(Row("&Group by:", groupBy));
            panel.Children.Add(Row("Save &to pool:", saveTo));
            panel.Children.Add(error);
            var buttons = Dialogs.Buttons(apply, save, cancel);
            buttons.Children.Insert(0, delete);
            panel.Children.Add(buttons);
            var w = Dialogs.Frame(edited == null ? "New Search" : $"Edit Search '{edited.Name}'", panel, 680);

            Result result = null;
            Search Build()
            {
                var chosen = typeBoxes.Where(b => b.Box.IsChecked == true).Aggregate(ObjectTypes.None, (a, b) => a | b.Type);
                var err = SearchRules.Validate(name.Text, chosen);
                if (err != null)
                {
                    error.Text = err;
                    return null;
                }
                var grouping = groupingChanged || edited == null ? groupings[Math.Max(0, groupBy.SelectedIndex)].Make() : edited.Grouping;
                return SearchRules.Build(name.Text, chosen, text.Text, grouping, edited);
            }
            apply.Click += (_, _) => { if (Build() is Search s) { result = new Result { Search = s }; w.Close(); } };
            save.Click += (_, _) =>
            {
                if (saveTo.SelectedIndex < 0) { error.Text = "Choose the pool to save the search to."; return; }
                if (Build() is Search s) { result = new Result { Search = s, SaveTo = pools[saveTo.SelectedIndex] }; w.Close(); }
            };
            delete.Click += async (_, _) =>
            {
                if (await Dialogs.Confirm(w, "Delete Search", $"Delete the saved search '{edited.Name}'?", Messages.YES, Messages.NO))
                {
                    result = new Result { Search = edited, Delete = true };
                    w.Close();
                }
            };
            cancel.Click += (_, _) => w.Close();
            await w.ShowDialog(owner);
            return result;
        }
    }
}
