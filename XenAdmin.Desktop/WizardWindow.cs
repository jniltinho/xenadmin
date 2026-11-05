using System;
using System.Threading.Tasks;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Layout;
using Avalonia.Media;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// The XenWizardBase frame: header with icon and page title, the step list with the logo, the page, and
    /// Previous / Next / Finish / Cancel with an error line. Subclasses give the pages, their validation and Finish.
    /// </summary>
    public abstract class WizardWindow : Window
    {
        private string[] _steps;
        private readonly StackPanel _stepList = new StackPanel();
        private readonly TextBlock _title = new TextBlock { FontWeight = FontWeight.Bold, FontSize = 13, VerticalAlignment = VerticalAlignment.Center };
        private readonly ContentControl _page = new ContentControl();
        private readonly Button _prev = new Button { Content = "< Previous" };
        private readonly Button _next = new Button { Content = "Next >", IsDefault = true };
        protected readonly Button FinishButton = new Button { Content = "Finish" };
        private readonly TextBlock _error = new TextBlock { Foreground = Brushes.Red, VerticalAlignment = VerticalAlignment.Center, TextWrapping = TextWrapping.Wrap };

        protected int Step { get; private set; }

        protected WizardWindow(string title, string icon32, params string[] steps)
        {
            _steps = steps;
            Title = title;
            Width = 810;
            Height = 545;
            MinWidth = 680;
            MinHeight = 480;
            WindowStartupLocation = WindowStartupLocation.CenterOwner;
            ShowInTaskbar = false;
            Background = Brush.Parse("#F0F0F0");

            var header = new DockPanel { Height = 60, Background = Brushes.White };
            var hicon = Views.Icon(icon32, 32);
            hicon.Margin = new Thickness(12, 0, 16, 0);
            header.Children.Add(hicon);
            header.Children.Add(_title);

            var logo = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 10, Margin = new Thickness(12, 0, 0, 14), VerticalAlignment = VerticalAlignment.Bottom };
            try { logo.Children.Add(new Image { Source = new Avalonia.Media.Imaging.Bitmap(Avalonia.Platform.AssetLoader.Open(new Uri("avares://xcp-ng-center/Assets/xcp-ng-center.png"))), Width = 40, Height = 40 }); } catch { }
            logo.Children.Add(new TextBlock { Text = "XCP-ng", FontSize = 20, VerticalAlignment = VerticalAlignment.Center });
            var side = new DockPanel { Width = 190, Background = Brushes.White };
            DockPanel.SetDock(logo, Dock.Bottom);
            side.Children.Add(logo);
            side.Children.Add(_stepList);

            _prev.Click += (_, _) => Go(Step - 1);
            _next.Click += (_, _) => { if (Check()) Go(Step + 1); };
            FinishButton.Click += async (_, _) => { if (Check()) await Finish(); };
            var cancel = new Button { Content = Messages.CANCEL, IsCancel = true };
            cancel.Click += (_, _) => Close();
            var buttons = new DockPanel { Margin = new Thickness(12, 10) };
            var right = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 8 };
            right.Children.Add(_prev);
            right.Children.Add(_next);
            right.Children.Add(FinishButton);
            right.Children.Add(cancel);
            DockPanel.SetDock(right, Dock.Right);
            buttons.Children.Add(right);
            buttons.Children.Add(_error);

            var body = new DockPanel();
            DockPanel.SetDock(side, Dock.Left);
            body.Children.Add(side);
            body.Children.Add(new Border { Padding = new Thickness(16, 14), Child = _page });

            var root = new DockPanel();
            DockPanel.SetDock(header, Dock.Top);
            DockPanel.SetDock(buttons, Dock.Bottom);
            root.Children.Add(header);
            var line = new Border { Height = 1, Background = Brush.Parse("#C8C8C8") };
            DockPanel.SetDock(line, Dock.Top);
            root.Children.Add(line);
            root.Children.Add(buttons);
            root.Children.Add(body);
            Content = root;
            Opened += (_, _) => Go(0);
        }

        /// <summary>The header title of a step.</summary>
        protected abstract string PageTitle(int step);

        /// <summary>The content of a step (built each time the step is shown).</summary>
        protected abstract Control Page(int step);

        /// <summary>Null when the step can be left, else the error to show.</summary>
        protected abstract string Validate(int step);

        protected abstract Task Finish();

        /// <summary>Shows an error under the page (e.g. from Finish).</summary>
        protected void ShowError(string text) => _error.Text = text ?? "";

        /// <summary>Changes the step list (a choice on one page adds or removes later pages, as UpdateWizardContent).</summary>
        protected void SetSteps(params string[] steps) => _steps = steps;

        protected void Go(int step)
        {
            _error.Text = "";
            Step = Math.Clamp(step, 0, _steps.Length - 1);
            _title.Text = PageTitle(Step);
            _stepList.Children.Clear();
            for (int i = 0; i < _steps.Length; i++)
                _stepList.Children.Add(new Border
                {
                    Background = i == Step ? Brush.Parse("#0078D7") : Brushes.Transparent,
                    Padding = new Thickness(12, 4),
                    Child = new TextBlock { Text = _steps[i], FontWeight = i == Step ? FontWeight.Bold : FontWeight.Normal, Foreground = i == Step ? Brushes.White : Brushes.Black }
                });
            _prev.IsEnabled = Step > 0;
            _next.IsVisible = Step < _steps.Length - 1;
            FinishButton.IsVisible = Step == _steps.Length - 1;
            FinishButton.IsDefault = FinishButton.IsVisible;
            _next.IsDefault = !FinishButton.IsVisible;
            _page.Content = Page(Step);
        }

        private bool Check()
        {
            var err = Validate(Step);
            _error.Text = err ?? "";
            return err == null;
        }
    }
}
