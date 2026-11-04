using System;
using System.Net;
using Avalonia;
using Avalonia.Controls.ApplicationLifetimes;
using Avalonia.Themes.Fluent;
using XenAdmin;

namespace XenAdmin.Linux
{
    public class App : Application
    {
        public override void Initialize()
        {
            Styles.Add(new FluentTheme());
            // the WinForms app has a single light look
            RequestedThemeVariant = Avalonia.Styling.ThemeVariant.Light;
            Styles.Add(new Avalonia.Markup.Xaml.Styling.StyleInclude(new Uri("avares://xcp-ng-center/")) { Source = new Uri("avares://xcp-ng-center/Theme.axaml") });
        }

        public override void OnFrameworkInitializationCompleted()
        {
            if (ApplicationLifetime is IClassicDesktopStyleApplicationLifetime desktop)
            {
                var config = AppConfig.Load();
                var trust = new CertificateTrust(config);
#pragma warning disable SYSLIB0014 // XenAPI's JsonRpc falls back to this global callback, as on Windows
                ServicePointManager.ServerCertificateValidationCallback = trust.Validate;
#pragma warning restore SYSLIB0014
                InvokeHelper.Initialize(new DispatcherSynchronizer());
                // keeps IXenObject.Path in step with other_config["folder"], as WinForms MainWindow does
                XenAdmin.Model.Folders.InitFolders();
                ConsoleView.Trust = trust;
                desktop.MainWindow = new MainWindow(config, trust);
            }
            base.OnFrameworkInitializationCompleted();
        }
    }

    public static class Program
    {
        [STAThread]
        public static int Main(string[] args)
        {
            if (Array.IndexOf(args, "--version") >= 0)
            {
                Console.WriteLine($"xcp-ng-center {typeof(Program).Assembly.GetName().Version}");
                return 0;
            }
            XenAdminConfigManager.Provider = new LinuxConfigProvider();
            return BuildAvaloniaApp().StartWithClassicDesktopLifetime(args);
        }

        public static AppBuilder BuildAvaloniaApp() =>
            AppBuilder.Configure<App>().UsePlatformDetect().WithInterFont().LogToTrace();
    }
}
