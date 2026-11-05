using System;
using System.Collections.Generic;
using System.IO;
using Newtonsoft.Json;

namespace XenAdmin.Desktop
{
    public class SavedServer
    {
        public string Host { get; set; }
        public int Port { get; set; } = 443;
        public string User { get; set; } = "root";
        /// <summary>Pool (or host) name shown while disconnected, like XenConnection.FriendlyName in the WinForms settings.</summary>
        public string Name { get; set; }
    }

    /// <summary>
    /// Server list and accepted certificate fingerprints, stored as JSON in the platform's settings directory.
    /// Passwords are never written here.
    /// </summary>
    public class AppConfig
    {
        public List<SavedServer> Servers { get; set; } = new List<SavedServer>();

        /// <summary>host:port → SHA-256 fingerprint (hex, upper case)</summary>
        public Dictionary<string, string> Fingerprints { get; set; } = new Dictionary<string, string>(StringComparer.OrdinalIgnoreCase);

        /// <summary>Tools > Options.</summary>
        public AppOptions Options { get; set; } = new AppOptions();

        [JsonIgnore]
        public string Path { get; private set; }

        private static readonly object FileLock = new object();

        /// <summary>The platform's per-user settings directory (see Platform.PlatformServices).</summary>
        public static string DefaultDirectory() => Platform.PlatformServices.Current.ConfigDirectory;

        public static AppConfig Load(string path = null)
        {
            path ??= System.IO.Path.Combine(DefaultDirectory(), "config.json");
            AppConfig cfg = null;
            lock (FileLock)
            {
                if (File.Exists(path))
                    cfg = JsonConvert.DeserializeObject<AppConfig>(File.ReadAllText(path));
            }
            cfg ??= new AppConfig();
            cfg.Fingerprints = new Dictionary<string, string>(cfg.Fingerprints ?? new Dictionary<string, string>(), StringComparer.OrdinalIgnoreCase);
            cfg.Options ??= new AppOptions();
            cfg.Path = path;
            return cfg;
        }

        public void Save()
        {
            lock (FileLock)
            {
                Directory.CreateDirectory(System.IO.Path.GetDirectoryName(Path));
                var tmp = Path + ".tmp";
                File.WriteAllText(tmp, JsonConvert.SerializeObject(this, Formatting.Indented));
                File.Move(tmp, Path, true);
            }
        }

        public void AddOrUpdateServer(string host, int port, string user, string name = null)
        {
            var s = Servers.Find(x => string.Equals(x.Host, host, StringComparison.OrdinalIgnoreCase) && x.Port == port);
            if (s == null)
                Servers.Add(new SavedServer { Host = host, Port = port, User = user, Name = name });
            else
            {
                s.User = user;
                s.Name = name ?? s.Name;
            }
            Save();
        }

        public void RemoveServer(string host, int port)
        {
            Servers.RemoveAll(x => string.Equals(x.Host, host, StringComparison.OrdinalIgnoreCase) && x.Port == port);
            Save();
        }
    }
}
