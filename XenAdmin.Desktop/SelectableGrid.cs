using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Input;
using Avalonia.Media;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// The WinForms DataGridView of the tab pages: column headers, one selectable row per key, Up/Down to move,
    /// Enter or double click to activate. Selecting repaints the rows instead of rebuilding them (a double click needs
    /// both clicks on the same control).
    /// </summary>
    public class SelectableGrid : Border
    {
        private readonly Grid _rows = new Grid();
        private IList<string> _keys = Array.Empty<string>();

        public string Selected { get; private set; }

        /// <summary>The user picked another row.</summary>
        public event Action SelectionChanged;

        /// <summary>Double click or Enter on the selected row.</summary>
        public event Action Activated;

        public SelectableGrid()
        {
            Focusable = true;
            BorderBrush = Brush.Parse("#828790");
            BorderThickness = new Thickness(1);
            Background = Brushes.White;
            Child = new ScrollViewer { Content = _rows, HorizontalScrollBarVisibility = Avalonia.Controls.Primitives.ScrollBarVisibility.Auto };
            PointerPressed += (_, _) => Focus();
            KeyDown += (_, e) =>
            {
                var i = _keys.IndexOf(Selected);
                if (e.Key == Key.Down && i < _keys.Count - 1) Select(_keys[i + 1]);
                else if (e.Key == Key.Up && i > 0) Select(_keys[i - 1]);
                else if (e.Key == Key.Enter && Selected != null) Activated?.Invoke();
                else return;
                e.Handled = true;
            };
        }

        /// <summary>
        /// Shows the table; keys[i] identifies row i. The selection stays on its key, else goes to the first row; with
        /// keepMissing it stays even on a key not listed (a row being recreated).
        /// </summary>
        public void Fill(Table table, IList<string> keys, string selected, bool keepMissing = false)
        {
            _keys = keys;
            Selected = selected != null && (keepMissing || keys.Contains(selected)) ? selected : keys.FirstOrDefault();
            _rows.Children.Clear();
            _rows.RowDefinitions.Clear();
            var fill = table.Fill >= 0 ? table.Fill : table.Headers.Length - 1;
            _rows.ColumnDefinitions = new ColumnDefinitions(string.Join(",", table.Headers.Select((_, i) => i == fill ? "*" : "Auto")));
            _rows.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            for (int c = 0; c < table.Headers.Length; c++)
                Cell(0, c, table.Headers[c], null);
            for (int r = 0; r < table.Rows.Count && r < keys.Count; r++)
            {
                _rows.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                for (int c = 0; c < table.Headers.Length; c++)
                    Cell(r + 1, c, c < table.Rows[r].Length ? table.Rows[r][c] : "", keys[r]);
            }
            Views.PaintSelection(_rows, Selected);
        }

        private void Select(string key)
        {
            Selected = key;
            Views.PaintSelection(_rows, key);
            _rows.Children.OfType<Border>().FirstOrDefault(b => (string)b.Tag == key)?.BringIntoView();
            SelectionChanged?.Invoke();
        }

        private void Cell(int row, int col, string text, string key)
        {
            var b = new Border
            {
                Background = Brushes.White,
                BorderBrush = Brush.Parse("#E3E3E3"),
                BorderThickness = key == null ? new Thickness(0, 0, 1, 1) : new Thickness(0),
                Padding = new Thickness(6, 3, 14, 3),
                Tag = key,
                Child = new TextBlock { Text = text ?? "", TextTrimming = TextTrimming.CharacterEllipsis }
            };
            if (key != null)
            {
                b.PointerPressed += (_, _) => Select(key);
                b.DoubleTapped += (_, _) => Activated?.Invoke();
            }
            Grid.SetRow(b, row);
            Grid.SetColumn(b, col);
            _rows.Children.Add(b);
        }
    }
}
