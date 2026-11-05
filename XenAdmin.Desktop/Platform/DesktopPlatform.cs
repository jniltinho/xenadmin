using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;

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

    /// <summary>Starting the desktop's own SSH terminal and Remote Desktop client (XenCenter used PuTTY/OpenSSH and an embedded RDP control).</summary>
    public static class ExternalTools
    {
        /// <summary>The full path of a program on PATH, or null.</summary>
        public static string Which(string name) =>
            Path.IsPathRooted(name) ? (File.Exists(name) ? name : null)
            : (Environment.GetEnvironmentVariable("PATH") ?? "").Split(Path.PathSeparator, StringSplitOptions.RemoveEmptyEntries)
                .Select(d => Path.Combine(d, name)).FirstOrDefault(File.Exists);

        public static bool Start((string Exe, string[] Args)? command, bool shell = false)
        {
            if (command is not var (exe, args))
                return false;
            try
            {
                // ArgumentList keeps each argument whole (also with UseShellExecute on .NET 10): no re-parsing of guest data
                var psi = new ProcessStartInfo(exe) { UseShellExecute = shell };
                foreach (var a in args) psi.ArgumentList.Add(a);
                return Process.Start(psi) != null;
            }
            catch (Exception)
            {
                return false;
            }
        }
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

        public bool LaunchSsh(string target) => ExternalTools.Start(SshCommand(ExternalTools.Which, Environment.GetEnvironmentVariable("TERMINAL"), target));
        public bool LaunchRdp(string host) => ExternalTools.Start(RdpCommand(ExternalTools.Which, Environment.GetEnvironmentVariable("TERMINAL"), host));

        // terminal emulators and how each runs a command: $TERMINAL first, then Debian's alternative, then the common ones
        private static readonly (string Name, string[] Run)[] Terminals =
        {
            ("x-terminal-emulator", new[] { "-e" }), ("gnome-terminal", new[] { "--" }), ("kgx", new[] { "--" }), ("konsole", new[] { "-e" }),
            ("xfce4-terminal", new[] { "-x" }), ("mate-terminal", new[] { "-x" }), ("tilix", new[] { "-e" }), ("alacritty", new[] { "-e" }),
            ("kitty", new string[0]), ("foot", new string[0]), ("wezterm", new[] { "start", "--" }), ("xterm", new[] { "-e" })
        };

        /// <summary>The terminal to run a command in, with its "run this" arguments.</summary>
        private static (string Exe, string[] Run)? Terminal(Func<string, string> which, string terminalEnv)
        {
            // $TERMINAL keeps the "run this" arguments of the terminal it names (gnome-terminal needs "--", not "-e")
            var known = Terminals.FirstOrDefault(t => t.Name == Path.GetFileName(terminalEnv ?? "")).Run;
            var candidates = string.IsNullOrEmpty(terminalEnv) ? Terminals : Terminals.Prepend((terminalEnv, known ?? new[] { "-e" }));
            foreach (var (name, run) in candidates)
                if (which(name) is string exe)
                    return (exe, run);
            return null;
        }

        private static (string Exe, string[] Args)? InTerminal(Func<string, string> which, string terminalEnv, params string[] command) =>
            Terminal(which, terminalEnv) is var (exe, run) ? (exe, run.Concat(command).ToArray()) : null;

        public static (string Exe, string[] Args)? SshCommand(Func<string, string> which, string terminalEnv, string target) =>
            which("ssh") is string ssh ? InTerminal(which, terminalEnv, ssh, "--", target) : null;

        /// <summary>
        /// FreeRDP asks on its terminal whether to trust the server's certificate, so it runs in one (never /cert-ignore);
        /// Remmina asks in its own window.
        /// </summary>
        public static (string Exe, string[] Args)? RdpCommand(Func<string, string> which, string terminalEnv, string host)
        {
            foreach (var name in new[] { "xfreerdp3", "xfreerdp", "sdl-freerdp3", "wlfreerdp" })
                if (which(name) is string freerdp && InTerminal(which, terminalEnv, freerdp, "/v:" + host, "/dynamic-resolution", "+clipboard") is var command && command != null)
                    return command;
            return which("remmina") is string remmina ? (remmina, new[] { "-c", "rdp://" + host }) : null;
        }
    }

    /// <summary>macOS: ~/Library/Application Support.</summary>
    public sealed class MacPlatform : IPlatformServices
    {
        public MacPlatform(Func<string, string> env, string home) =>
            ConfigDirectory = Path.Combine(home, "Library", "Application Support", "xcp-ng-center");

        public string Name => "macOS";
        public string ConfigDirectory { get; }
        public ISecretStore Secrets { get; } = new NoSecretStore();

        // Terminal.app runs the command; the target travels as an argument and is shell-quoted by AppleScript
        public bool LaunchSsh(string target) =>
            ExternalTools.Start(("/usr/bin/osascript", new[] { "-e", "on run argv", "-e",
                "tell application \"Terminal\" to do script \"ssh -- \" & quoted form of (item 1 of argv)", "-e", "end run", target }));

        // the rdp: URL scheme of Microsoft's Windows App / Remote Desktop
        public bool LaunchRdp(string host) => ExternalTools.Start(("/usr/bin/open", new[] { "rdp://full%20address=s:" + host }));
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

        // Windows' own OpenSSH client in a new console window, and mstsc
        public bool LaunchSsh(string target) => ExternalTools.Which("ssh.exe") is string ssh && ExternalTools.Start((ssh, new[] { "--", target }), shell: true);
        public bool LaunchRdp(string host) => ExternalTools.Start(("mstsc.exe", new[] { "/v:" + host }), shell: true);
    }
}
