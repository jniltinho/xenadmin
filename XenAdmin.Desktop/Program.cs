using System;
using System.Linq;
using System.Net;
using Avalonia;
using Avalonia.Controls.ApplicationLifetimes;
using Avalonia.Themes.Fluent;
using XenAdmin;

namespace XenAdmin.Desktop
{
    public class App : Application
    {
        public override void Initialize()
        {
            Name = "XCP-ng Center"; // the macOS menu bar shows it instead of "Avalonia Application"
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
                var provider = (DesktopConfigProvider)XenAdminConfigManager.Provider;
                provider.Options = config.Options;
                provider.ApplyProxy();
                var trust = new CertificateTrust(config);
#pragma warning disable SYSLIB0014 // XenAPI's JsonRpc falls back to this global callback, as on Windows
                ServicePointManager.ServerCertificateValidationCallback = trust.Validate;
                // console tunnels and RRD downloads (XenAPI.HTTP) use the same policy
                XenAPI.HTTP.ServerCertificateValidator = trust.ValidateHost;
                // pool members: trust the host certificate the authenticated coordinator reports
                trust.VouchedFingerprints = CertificateTrust.PoolMemberFingerprints;
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
            Platform.PlatformServices.Current = Platform.DesktopPlatform.Detect();
            XenAdminConfigManager.Provider = new DesktopConfigProvider();
            if (args.Length > 0 && args[0] == "mcp")
            {
                // MCP server on stdio, no window: stdout carries only the protocol (UTF-8 without BOM)
                var utf8 = new System.Text.UTF8Encoding(false);
                Console.InputEncoding = utf8;
                Console.OutputEncoding = utf8;
                return Mcp.McpHost.Run(args.Skip(1).ToArray(), Console.In, Console.Out, Console.Error);
            }
            return BuildAvaloniaApp().StartWithClassicDesktopLifetime(args);
        }

        public static AppBuilder BuildAvaloniaApp() =>
            AppBuilder.Configure<App>().UsePlatformDetect().WithInterFont().LogToTrace();
    }
}
