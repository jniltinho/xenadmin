using System;
using System.Linq;
using System.Net;
using System.Net.Http;
using System.Net.Security;
using System.Net.Sockets;
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;

namespace XenAdmin.Desktop
{
    public enum TrustDecision { Trusted, Prompt, Mismatch }

    /// <summary>
    /// Trust-on-first-use: a stored SHA-256 pin always wins over the system chain;
    /// without a pin, a system-trusted chain is accepted and pinned; otherwise the user is asked.
    /// </summary>
    public class CertificateTrust
    {
        private readonly AppConfig _config;

        /// <summary>
        /// Certificates rejected during a handshake, waiting for the user's decision (hostKey → (new, pinned or null)).
        /// The TLS callback must never block on UI: XAPI drops the request (408) if the handshake stalls.
        /// </summary>
        public System.Collections.Concurrent.ConcurrentDictionary<string, (string Fingerprint, string Old)> Pending { get; } =
            new System.Collections.Concurrent.ConcurrentDictionary<string, (string, string)>(StringComparer.OrdinalIgnoreCase);

        /// <summary>
        /// Fingerprints that an already-authenticated XAPI session reports for a "host:port"
        /// (pool members' host certificates). Used only when nothing is pinned yet.
        /// </summary>
        public Func<string, System.Collections.Generic.IEnumerable<string>> VouchedFingerprints { get; set; }

        /// <summary>Hex without separators, upper case (XAPI uses "AA:BB:..").</summary>
        public static string Normalize(string fingerprint) => (fingerprint ?? "").Replace(":", "").Trim().ToUpperInvariant();

        public CertificateTrust(AppConfig config)
        {
            _config = config;
        }

        public static string Fingerprint(X509Certificate cert)
        {
            return Convert.ToHexString(SHA256.HashData(cert.GetRawCertData()));
        }

        public TrustDecision Evaluate(string hostKey, string fingerprint, SslPolicyErrors errors)
        {
            lock (_config)
            {
                if (_config.Fingerprints.TryGetValue(hostKey, out var pinned))
                    return string.Equals(pinned, fingerprint, StringComparison.OrdinalIgnoreCase) ? TrustDecision.Trusted : TrustDecision.Mismatch;
            }
            return errors == SslPolicyErrors.None ? TrustDecision.Trusted : TrustDecision.Prompt;
        }

        public void Pin(string hostKey, string fingerprint)
        {
            lock (_config)
            {
                _config.Fingerprints[hostKey] = fingerprint;
                _config.Save();
            }
        }

        /// <summary>Callback for ServicePointManager/HttpWebRequest (sender is the request).</summary>
        public bool Validate(object sender, X509Certificate cert, X509Chain chain, SslPolicyErrors errors)
        {
            var hostKey = HostKey(sender);
            return hostKey == null ? cert != null && errors == SslPolicyErrors.None : ValidateHost(hostKey, cert, errors);
        }

        /// <summary>Same policy for raw TLS streams (console tunnels), where there is no request object.</summary>
        public bool ValidateHost(string hostKey, X509Certificate cert, SslPolicyErrors errors)
        {
            if (cert == null)
                return false;
            var fp = Fingerprint(cert);
            string old;
            lock (_config)
                _config.Fingerprints.TryGetValue(hostKey, out old);
            var decision = Evaluate(hostKey, fp, errors);
            if (decision == TrustDecision.Prompt && IsVouched(hostKey, cert))
                decision = TrustDecision.Trusted;
            if (decision == TrustDecision.Trusted)
            {
                if (old == null)
                    Pin(hostKey, fp);
                return true;
            }
            Pending[hostKey] = (fp, old);
            return false;
        }

        private bool IsVouched(string hostKey, X509Certificate cert)
        {
            var known = VouchedFingerprints?.Invoke(hostKey)?.Select(Normalize).ToList();
            if (known == null || known.Count == 0)
                return false;
            var raw = cert.GetRawCertData();
            return known.Contains(Convert.ToHexString(SHA256.HashData(raw))) || known.Contains(Convert.ToHexString(SHA1.HashData(raw)));
        }

        private static string HostKey(object sender)
        {
            Uri uri = sender switch
            {
                HttpWebRequest r => r.RequestUri,
                HttpRequestMessage m => m.RequestUri,
                _ => null
            };
            return uri == null ? null : $"{uri.Host}:{uri.Port}";
        }
    }
}
