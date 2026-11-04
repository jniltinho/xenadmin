using System;
using System.IO;

namespace XenAdmin.Desktop.Platform
{
    /// <summary>The OS adapters of IPlatformServices; a port adds or adjusts one class here.</summary>
    public static class DesktopPlatform
    {
        public static IPlatformServices Detect() =>
            OperatingSystem.IsMacOS() ? new MacPlatform(Environment.GetEnvironmentVariable, Home())
            : OperatingSystem.IsWindows() ? new WindowsPlatform(Environment.GetEnvironmentVariable, Home())
            : new UnixPlatform(OperatingSystem.IsFreeBSD() ? "FreeBSD" : "Linux", Environment.GetEnvironmentVariable, Home());

        private static string Home() => Environment.GetFolderPath(Environment.SpecialFolder.UserProfile);
    }

    /// <summary>Linux and FreeBSD: XDG base directories.</summary>
    public sealed class UnixPlatform : IPlatformServices
    {
        public UnixPlatform(string name, Func<string, string> env, string home)
        {
            Name = name;
            var xdg = env("XDG_CONFIG_HOME");
            ConfigDirectory = Path.Combine(string.IsNullOrEmpty(xdg) ? Path.Combine(home, ".config") : xdg, "xcp-ng-center");
        }

        public string Name { get; }
        public string ConfigDirectory { get; }
        public ISecretStore Secrets { get; } = new NoSecretStore();
    }

    /// <summary>macOS: ~/Library/Application Support.</summary>
    public sealed class MacPlatform : IPlatformServices
    {
        public MacPlatform(Func<string, string> env, string home) =>
            ConfigDirectory = Path.Combine(home, "Library", "Application Support", "xcp-ng-center");

        public string Name => "macOS";
        public string ConfigDirectory { get; }
        public ISecretStore Secrets { get; } = new NoSecretStore();
    }

    /// <summary>Windows (Avalonia front end, no WinForms): %APPDATA%, kept apart from the WinForms app's settings.</summary>
    public sealed class WindowsPlatform : IPlatformServices
    {
        public WindowsPlatform(Func<string, string> env, string home)
        {
            var appData = env("APPDATA");
            ConfigDirectory = Path.Combine(string.IsNullOrEmpty(appData) ? Path.Combine(home, "AppData", "Roaming") : appData, "xcp-ng-center");
        }

        public string Name => "Windows";
        public string ConfigDirectory { get; }
        public ISecretStore Secrets { get; } = new NoSecretStore();
    }
}
