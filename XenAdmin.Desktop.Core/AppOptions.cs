using System;
using System.Net;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>
    /// Tools > Options (XenAdmin/Dialogs/OptionsDialog): the settings of the Security, Connection, Confirmations and
    /// Display pages that apply to this app, saved in config.json with the server list. Save and Restore is not
    /// ported (this app never stores passwords), nor Consoles, External Tools and Plugins.
    /// </summary>
    public class AppOptions
    {
        // SecurityOptionsPage. WinForms does not warn about a new certificate by default; this app always asked
        // (the safer behaviour) and keeps that default.
        public bool WarnUnrecognizedCertificate { get; set; } = true;
        public bool WarnChangedCertificate { get; set; } = true;
        public bool RotateSecretReminder { get; set; } = true; // RemindChangePassword

        // ConnectionOptionsPage
        public int ConnectionTimeoutSeconds { get; set; } = 20;
        public HTTPHelper.ProxyStyle ProxyStyle { get; set; } = HTTPHelper.ProxyStyle.DirectConnection;
        public string ProxyAddress { get; set; } = "";
        public int ProxyPort { get; set; } = 80;
        public bool BypassProxyForServers { get; set; }

        // ConfirmationOptionsPage
        public bool DoNotConfirmDismissAlerts { get; set; }

        // DisplayOptionsPage
        public bool FillAreaUnderGraphs { get; set; }
        public bool RememberLastSelectedTab { get; set; } = true;

        /// <summary>ConnectionOptionsPage's NumericUpDown range.</summary>
        public const int MinTimeout = 1, MaxTimeout = 2147483;

        /// <summary>ConnectionOptionsPage.IsValidToSave: a host name or IPv4 address and a valid port for a specified proxy.</summary>
        public static string ProxyError(HTTPHelper.ProxyStyle style, string address, string port)
        {
            if (style != HTTPHelper.ProxyStyle.SpecifiedProxy)
                return null;
            var type = Uri.CheckHostName(address?.Trim());
            if (type == UriHostNameType.Unknown || type == UriHostNameType.IPv6)
                return "Enter the proxy server's host name or IPv4 address.";
            return Util.IsValidPort(port?.Trim()) ? null : "Enter a proxy port between 1 and 65535.";
        }

        /// <summary>
        /// WinformsXenAdminConfigProvider.GetProxyFromSettings. shortcut: no proxy authentication (its credentials would
        /// have to be stored); add it with the platform secret store.
        /// </summary>
        public IWebProxy Proxy(bool isForXenServer)
        {
            switch (ProxyStyle)
            {
                case HTTPHelper.ProxyStyle.SpecifiedProxy:
                    return isForXenServer && BypassProxyForServers ? null : new WebProxy($"http://{ProxyAddress}:{ProxyPort}", false);
                case HTTPHelper.ProxyStyle.SystemProxy:
                    return WebRequest.GetSystemWebProxy();
                default:
                    return null;
            }
        }
    }
}
