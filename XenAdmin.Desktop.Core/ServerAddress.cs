using System.Net;
using System.Net.Sockets;

namespace XenAdmin.Desktop
{
    /// <summary>The "Server:" field of Add New Server / Connect: host, IPv4, IPv6 (bare or [bracketed]) and an optional port.</summary>
    public static class ServerAddress
    {
        public const int DefaultPort = 443;

        /// <summary>Splits "host[:port]", "[v6]:port" or a bare IPv6 address; false when it is empty or the port is invalid.</summary>
        public static bool TryParse(string text, out string host, out int port)
        {
            host = (text ?? "").Trim();
            port = DefaultPort;
            if (host.StartsWith("["))
            {
                var end = host.IndexOf(']');
                if (end < 0) return false;
                var rest = host[(end + 1)..];
                host = host[1..end];
                if (rest.Length > 0 && !(rest[0] == ':' && int.TryParse(rest[1..], out port)))
                    return false;
            }
            else if (IPAddress.TryParse(host, out var ip) && ip.AddressFamily == AddressFamily.InterNetworkV6)
            {
                // a bare IPv6 address has no port: "2001:db8::11" is not host "2001:db8:" port 11
            }
            else
            {
                var idx = host.LastIndexOf(':');
                if (idx >= 0)
                {
                    if (!int.TryParse(host[(idx + 1)..], out port)) return false;
                    host = host[..idx];
                }
            }
            return host.Length > 0 && port >= 1 && port <= 65535;
        }

        /// <summary>Session.GetUrl talks plain HTTP on 80 and 8080: the password would cross the network unencrypted.</summary>
        public static bool IsPlainHttp(int port) => port == 80 || port == 8080;
    }
}
