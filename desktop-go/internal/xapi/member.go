package xapi

import (
	"crypto/tls"
	"errors"

	"go/xenapi"

	"xcpng-center/internal/trust"
)

// VouchedFingerprints are the host certificates the authenticated coordinator reports for a server address
// (CertificateTrust.PoolMemberFingerprints): a pool member is trusted through them when nothing is pinned for it.
func VouchedFingerprints(s *xenapi.Session, address string) []string {
	var out []string
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return nil
	}
	for _, h := range hosts {
		if h.Address != address {
			continue
		}
		for _, cr := range h.Certificates {
			if cert, err := xenapi.Certificate.GetRecord(s, cr); err == nil && cert.Type == xenapi.CertificateTypeHost {
				out = append(out, trust.Normalize(cert.FingerprintSha256), trust.Normalize(cert.Fingerprint))
			}
		}
	}
	return out
}

// MemberTLS is the TLS configuration for a direct stream to a pool server (a console, RRD updates): the server is
// accepted when the policy trusts it (no prompt), or, with nothing pinned for it, when the coordinator vouches for
// its certificate.
func MemberTLS(policy *trust.Policy, hostKey string, vouched []string) *tls.Config {
	return &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true, // checked below
		VerifyConnection: func(cs tls.ConnectionState) error {
			if len(cs.PeerCertificates) == 0 {
				return errors.New("the server sent no certificate")
			}
			fp := trust.Fingerprint(cs.PeerCertificates[0].Raw)
			if policy.Evaluate(hostKey, fp, false) == trust.Trusted {
				return nil
			}
			if _, pinned := policy.Expected(hostKey); !pinned {
				for _, v := range vouched {
					if v != "" && v == fp {
						return nil
					}
				}
			}
			return &trust.UntrustedError{HostKey: hostKey, Fingerprint: fp}
		}}
}
