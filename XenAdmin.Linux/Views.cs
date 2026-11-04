using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Input;
using Avalonia.Layout;
using Avalonia.Media;

namespace XenAdmin.Linux
{
    /// <summary>Building blocks that mimic the WinForms look (title band, GradientPanel, PDSection, DataGridView).</summary>
    public static class Views
    {
        public static readonly IBrush Link = Brush.Parse("#0000EE");
        private static readonly IBrush Selection = Brush.Parse("#0078D7");
        private static readonly IBrush SectionBorder = Brush.Parse("#DCDCDC");   // Gainsboro, PDSection
        private static readonly IBrush SectionHeader = Brush.Parse("#F5F5F5");   // WhiteSmoke, PDSection
        private static readonly IBrush GridLine = Brush.Parse("#E3E3E3");
        private static readonly Cursor Hand = new Cursor(StandardCursorType.Hand);

        public static IBrush Gradient(string from, string to, bool horizontal = false) => new LinearGradientBrush
        {
            StartPoint = new RelativePoint(0, 0, RelativeUnit.Relative),
            EndPoint = horizontal ? new RelativePoint(1, 0, RelativeUnit.Relative) : new RelativePoint(0, 1, RelativeUnit.Relative),
            GradientStops = { new GradientStop(Color.Parse(from), 0), new GradientStop(Color.Parse(to), 1) }
        };

        public static Image Icon(string file, double size = 16) =>
            new Image { Source = Images.Get(file), Width = size, Height = size, VerticalAlignment = VerticalAlignment.Center };

        /// <summary>Grey band above the tabs: icon + object heading, "Logged in as" on the right.</summary>
        public static Control TitleBand(string icon, string title, string right)
        {
            var dock = new DockPanel { Height = 24, Background = Gradient("#F4F4F4", "#C9C9C9") };
            var r = new TextBlock { Text = right ?? "", Margin = new Thickness(0, 0, 8, 0), VerticalAlignment = VerticalAlignment.Center };
            DockPanel.SetDock(r, Dock.Right);
            dock.Children.Add(r);
            var left = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 6, Margin = new Thickness(4, 0) };
            left.Children.Add(Icon(icon));
            left.Children.Add(new TextBlock { Text = title, FontWeight = FontWeight.Bold, FontSize = 13, VerticalAlignment = VerticalAlignment.Center });
            dock.Children.Add(left);
            return new Border { BorderBrush = Brush.Parse("#A0A0A0"), BorderThickness = new Thickness(1), Child = dock };
        }

        /// <summary>The teal GradientPanel at the top of each tab (#396D8C → #3F8B89).</summary>
        public static Control TabHeader(string text) => new Border
        {
            Height = 36,
            Background = Gradient("#396D8C", "#3F8B89", true),
            Child = new TextBlock { Text = text, Foreground = Brushes.White, FontWeight = FontWeight.Bold, FontSize = 13, Margin = new Thickness(10, 0), VerticalAlignment = VerticalAlignment.Center }
        };

        public static TextBlock LinkText(string text, Action onClick)
        {
            var t = new TextBlock { Text = text, Foreground = onClick == null ? Brushes.Gray : Link, Cursor = onClick == null ? null : Hand };
            if (onClick != null)
                t.PointerPressed += (_, _) => onClick();
            return t;
        }

        /// <summary>Content limited to the WinForms MaximumWidth (900px), left-aligned, stretching below it.</summary>
        private static Control MaxWidth(Control content, double max = 900)
        {
            var host = new Grid { ColumnDefinitions = { new ColumnDefinition(GridLength.Star) { MaxWidth = max }, new ColumnDefinition(GridLength.Auto) } };
            host.Children.Add(content);
            return host;
        }

        /// <summary>Collapsible property sections with "Expand all / Collapse all" and a selectable row, like PDSection.</summary>
        public static Control Sections(IReadOnlyList<Section> sections, Action properties = null)
        {
            var panel = new StackPanel { Margin = new Thickness(12, 10, 12, 12), Spacing = 10 };
            var expanders = new List<Action<bool>>();
            var top = new DockPanel();
            var links = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 18, VerticalAlignment = VerticalAlignment.Center };
            links.Children.Add(LinkText("Expand all", () => expanders.ForEach(e => e(true))));
            links.Children.Add(LinkText("Collapse all", () => expanders.ForEach(e => e(false))));
            DockPanel.SetDock(links, Dock.Right);
            top.Children.Add(links);
            var props = new Button { Content = "Properties", IsEnabled = properties != null, Padding = new Thickness(8, 1), HorizontalAlignment = HorizontalAlignment.Left };
            if (properties != null)
                props.Click += (_, _) => properties();
            top.Children.Add(props);
            panel.Children.Add(top);

            Border selected = null;
            void Select(Border row)
            {
                if (selected != null) Paint(selected, false);
                selected = row;
                Paint(row, true);
            }
            static void Paint(Border row, bool on)
            {
                row.Background = on ? Selection : Brushes.Transparent;
                foreach (var t in ((Grid)row.Child).Children.OfType<TextBlock>())
                    t.Foreground = on ? Brushes.White : Brushes.Black;
            }

            foreach (var s in sections.Where(s => s.Rows.Count > 0))
            {
                var rows = new StackPanel();
                foreach (var (key, value) in s.Rows)
                {
                    var g = new Grid { ColumnDefinitions = new ColumnDefinitions("150,*") };
                    var k = new TextBlock { Text = key + ":", TextWrapping = TextWrapping.Wrap, Margin = new Thickness(8, 6, 8, 6) };
                    var v = new TextBlock { Text = value, TextWrapping = TextWrapping.Wrap, Margin = new Thickness(0, 6, 8, 6) };
                    Grid.SetColumn(v, 1);
                    g.Children.Add(k);
                    g.Children.Add(v);
                    var row = new Border { Child = g, Background = Brushes.Transparent };
                    row.PointerPressed += (_, _) => Select(row);
                    rows.Children.Add(row);
                    if (selected == null)
                        Select(row); // the first row starts selected, as in the WinForms tab
                }
                var chevron = new Image { Width = 20, Height = 20, VerticalAlignment = VerticalAlignment.Center };
                var header = new DockPanel { Height = 32, Background = SectionHeader, Cursor = Hand };
                var right = new StackPanel { Orientation = Orientation.Horizontal, Margin = new Thickness(0, 0, 6, 0) };
                right.Children.Add(new Border { Width = 1, Background = SectionBorder, Margin = new Thickness(0, 6, 8, 6) });
                right.Children.Add(chevron);
                DockPanel.SetDock(right, Dock.Right);
                header.Children.Add(right);
                header.Children.Add(new TextBlock { Text = s.Title, FontSize = 15, VerticalAlignment = VerticalAlignment.Center, Margin = new Thickness(8, 0) });
                void Set(bool open)
                {
                    rows.IsVisible = open;
                    chevron.Source = Images.Get(open ? "PDChevronUp" : "PDChevronDown");
                }
                Set(s.Expanded);
                header.PointerPressed += (_, _) => Set(!rows.IsVisible);
                expanders.Add(Set);
                var body = new StackPanel();
                body.Children.Add(header);
                body.Children.Add(rows);
                panel.Children.Add(new Border { BorderBrush = SectionBorder, BorderThickness = new Thickness(1), Background = Brushes.White, Child = body });
            }
            return new ScrollViewer { Content = MaxWidth(panel, 924) };
        }

        /// <summary>Read-only list with column headers, sized like the WinForms DataGridView (one Fill column).</summary>
        public static Control TableView(Table t, string empty = "None")
        {
            var grid = new Grid();
            var fill = t.Fill >= 0 ? t.Fill : t.Headers.Length - 1;
            for (int c = 0; c < t.Headers.Length; c++)
                grid.ColumnDefinitions.Add(new ColumnDefinition(c == fill ? GridLength.Star : GridLength.Auto));
            void Cell(int r, int c, string text, bool head)
            {
                var b = new Border
                {
                    BorderBrush = GridLine,
                    BorderThickness = head ? new Thickness(0, 0, 1, 1) : new Thickness(0),
                    Background = Brushes.White,
                    Padding = new Thickness(6, 3, 14, 3),
                    Child = new TextBlock { Text = text ?? "", TextWrapping = TextWrapping.NoWrap, TextTrimming = TextTrimming.CharacterEllipsis }
                };
                Grid.SetRow(b, r);
                Grid.SetColumn(b, c);
                grid.Children.Add(b);
            }
            grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
            for (int c = 0; c < t.Headers.Length; c++)
                Cell(0, c, t.Headers[c], true);
            for (int r = 0; r < t.Rows.Count; r++)
            {
                grid.RowDefinitions.Add(new RowDefinition(GridLength.Auto));
                for (int c = 0; c < t.Headers.Length; c++)
                    Cell(r + 1, c, c < t.Rows[r].Length ? t.Rows[r][c] : "", false);
            }
            var content = new StackPanel();
            content.Children.Add(grid);
            if (t.Rows.Count == 0)
                content.Children.Add(new TextBlock { Text = empty, Opacity = 0.6, Margin = new Thickness(8) });
            return new Border
            {
                BorderBrush = Brush.Parse("#828790"),
                BorderThickness = new Thickness(1),
                Background = Brushes.White,
                Child = new ScrollViewer { HorizontalScrollBarVisibility = Avalonia.Controls.Primitives.ScrollBarVisibility.Auto, Content = content }
            };
        }

        /// <summary>A grid page: optional link on top, the table, and the WinForms action buttons below.</summary>
        public static Control GridPage(Table t, string topLink = null, params string[] buttons)
        {
            // WinForms layout: link strip, then the table in 75% and the button row on top of the last 25%
            var dock = new Grid { Margin = new Thickness(12), RowDefinitions = new RowDefinitions("Auto,3*,Auto,*") };
            if (topLink != null)
            {
                var link = LinkText(topLink, null);
                link.HorizontalAlignment = HorizontalAlignment.Center;
                link.TextDecorations = TextDecorations.Underline;
                link.FontWeight = FontWeight.Bold;
                link.Margin = new Thickness(0, 12, 0, 18);
                dock.Children.Add(link);
            }
            if (buttons.Length > 0)
            {
                var bar = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8, Margin = new Thickness(0, 12, 0, 0) };
                foreach (var b in buttons)
                    bar.Children.Add(new Button { Content = b, MinWidth = 94, Height = 23, Padding = new Thickness(6, 0), IsEnabled = false, HorizontalContentAlignment = HorizontalAlignment.Center, VerticalContentAlignment = VerticalAlignment.Center });
                Grid.SetRow(bar, 2);
                dock.Children.Add(bar);
            }
            var table = TableView(t);
            Grid.SetRow(table, 1);
            dock.Children.Add(table);
            return MaxWidth(dock, 924);
        }

        public static Control HomePage(Action addServer)
        {
            var panel = new StackPanel { HorizontalAlignment = HorizontalAlignment.Center, VerticalAlignment = VerticalAlignment.Top, Spacing = 46, Margin = new Thickness(0, 50, 0, 0) };
            try
            {
                var bmp = new Avalonia.Media.Imaging.Bitmap(Avalonia.Platform.AssetLoader.Open(new Uri("avares://xcp-ng-center/Assets/home-logo.png")));
                panel.Children.Add(new Image { Source = bmp, Width = 300, Height = 300 });
            }
            catch { }
            var link = LinkText("Add a Server", addServer);
            link.FontFamily = new FontFamily("serif");
            link.FontSize = 15;
            link.TextDecorations = TextDecorations.Underline;
            link.HorizontalAlignment = HorizontalAlignment.Center;
            panel.Children.Add(link);
            return new ScrollViewer { Background = Brushes.White, VerticalScrollBarVisibility = Avalonia.Controls.Primitives.ScrollBarVisibility.Visible, Content = panel };
        }
    }
}
