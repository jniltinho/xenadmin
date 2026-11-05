using System;

namespace XenAdmin.Desktop.Platform
{
    /// <summary>
    /// Everything the desktop app needs from the operating system. One implementation per OS;
    /// a port (macOS, FreeBSD, Windows without WinForms) only adds or adjusts one of these.
    /// </summary>
    public interface IPlatformServices
    {
        string Name { get; }
        /// <summary>Per-user settings directory (server list, certificate pins).</summary>
        string ConfigDirectory { get; }
        /// <summary>Password storage of the desktop (Secret Service, Keychain, DPAPI); may be unavailable.</summary>
        ISecretStore Secrets { get; }
        /// <summary>Opens the desktop's terminal running "ssh target" (Launch SSH Console); false when there is none.</summary>
        bool LaunchSsh(string target);
        /// <summary>Opens the desktop's Remote Desktop client on host (Switch to Remote Desktop); false when there is none.</summary>
        bool LaunchRdp(string host);
    }

    public interface ISecretStore
    {
        bool IsAvailable { get; }
        void Save(string key, string secret);
        string Load(string key);
        void Delete(string key);
    }

    /// <summary>No password storage: passwords are asked on every connection and never written to disk.</summary>
    public sealed class NoSecretStore : ISecretStore
    {
        public bool IsAvailable => false;
        public void Save(string key, string secret) => throw new NotSupportedException("No secret store on this platform");
        public string Load(string key) => null;
        public void Delete(string key) { }
    }

    public static class PlatformServices
    {
        private static IPlatformServices _current;

        /// <summary>Set once at startup by the front end (it knows the OS); tests may replace it.</summary>
        public static IPlatformServices Current
        {
            get => _current ?? throw new InvalidOperationException("PlatformServices.Current is not set");
            set => _current = value;
        }
    }
}
