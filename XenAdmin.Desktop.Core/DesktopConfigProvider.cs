using System;
using System.Collections.Generic;
using System.IO;
using System.Net;
using XenAdmin.Actions;
using XenAdmin.Network;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>Linux counterpart of WinformsXenAdminConfigProvider, with the same defaults.</summary>
    public class DesktopConfigProvider : IXenAdminConfigProvider
    {
        private readonly HashSet<string> _hidden = new HashSet<string>();

        public Func<List<Role>, IXenConnection, string, AsyncAction.SudoElevationResult> ElevatedSessionDelegate => (roles, connection, title) => null;
        public int ConnectionTimeout => 20000;
        public Session CreateActionSession(Session session, IXenConnection connection) => new Session(session, connection) { Timeout = ConnectionTimeout };
        public bool Exiting { get; set; }
        public bool ForcedExiting => false;
        public string XenCenterUUID { get; } = Guid.NewGuid().ToString();
        public bool DontSudo => true;
        public int GetProxyTimeout(bool timeout) => timeout ? 40000 : 0;
        public void ShowObject(string opaqueRef) { lock (_hidden) _hidden.Remove(opaqueRef); }
        public void HideObject(string opaqueRef) { lock (_hidden) _hidden.Add(opaqueRef); }
        public bool ObjectIsHidden(string opaqueRef) { lock (_hidden) return _hidden.Contains(opaqueRef); }
        public string GetLogFile() => Path.Combine(AppConfig.DefaultDirectory(), "xcp-ng-center.log");
        public void UpdateServerHistory(string hostnameWithPort) { }
        public void SaveSettingsIfRequired() { }
        public bool ShowHiddenVMs => false;
        public string GetXenCenterMetadata() => "{}";
        public string GetCustomClientUpdatesXmlLocation() => null;
        public string GetCustomCfuLocation() => null;
        public string GetClientUpdatesQueryParam() => string.Empty;
        public string GetCustomFileServicePrefix() => null;
        public string FileServiceUsername => null;
        public string FileServiceClientId => null;
        public string GetCustomTokenUrl() => null;
        public IWebProxy GetProxyFromSettings(IXenConnection connection) => null;
        public IWebProxy GetProxyFromSettings(IXenConnection connection, bool isForXenServer) => null;
    }
}
