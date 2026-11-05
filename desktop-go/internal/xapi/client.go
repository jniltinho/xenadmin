package xapi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"go/xenapi"

	"xcpng-center/internal/trust"
)

// HTTPClient is the only HTTP client that talks to servers: TLS 1.2+, the certificate decided by the trust policy
// (the system's verdict counts only for a first use without a pin), and no redirects (a redirect could send the
// session or password elsewhere).
// host is the destination the system's verdict is checked against (a name or an IP, never the SNI alone).
func HTTPClient(policy *trust.Policy, host, hostKey string, timeout time.Duration) *http.Client {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true, // verified below by VerifyConnection with the app's policy
		VerifyConnection: func(cs tls.ConnectionState) error {
			if len(cs.PeerCertificates) == 0 {
				return errors.New("the server sent no certificate")
			}
			leaf := cs.PeerCertificates[0]
			fp := trust.Fingerprint(leaf.Raw)
			system := verifiesWithSystem(cs, leaf, host)
			switch policy.Evaluate(hostKey, fp, system) {
			case trust.Trusted:
				if _, pinned := policy.Store.Fingerprint(hostKey); !pinned && !policy.RequirePin {
					// trust on first use, as CertificateTrust.ValidateHost; an unsaved pin would lose the protection
					// against a changed certificate, so the connection fails instead
					if err := policy.Store.Pin(hostKey, fp); err != nil {
						return fmt.Errorf("could not save the certificate of %s: %w", hostKey, err)
					}
				}
				return nil
			case trust.Mismatch:
				prev, _ := policy.Store.Fingerprint(hostKey)
				return &trust.UntrustedError{HostKey: hostKey, Fingerprint: fp, Previous: prev}
			default:
				return &trust.UntrustedError{HostKey: hostKey, Fingerprint: fp}
			}
		},
	}
	return &http.Client{
		Timeout:       timeout,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		Transport: &http.Transport{
			TLSClientConfig:     tlsConfig,
			Proxy:               http.ProxyFromEnvironment,
			DialContext:         (&net.Dialer{Timeout: 20 * time.Second}).DialContext,
			TLSHandshakeTimeout: 20 * time.Second,
		},
	}
}

// verifiesWithSystem: the system roots vouch for the certificate and it names the destination (DNS name or IP SAN).
func verifiesWithSystem(cs tls.ConnectionState, leaf *x509.Certificate, host string) bool {
	inter := x509.NewCertPool()
	for _, c := range cs.PeerCertificates[1:] {
		inter.AddCert(c)
	}
	_, err := leaf.Verify(x509.VerifyOptions{DNSName: host, Intermediates: inter})
	return err == nil
}

// Connection is a logged-in session on a pool coordinator (or standalone server).
type Connection struct {
	Host    string
	Port    int
	User    string
	Session *xenapi.Session
}

// Connect logs in and follows HOST_IS_SLAVE to the coordinator (XenConnection's redirect), whose certificate goes
// through the same policy. The password is never logged nor part of any error.
func Connect(ctx context.Context, policy *trust.Policy, host string, port int, user, password string) (*Connection, error) {
	// shortcut: the SDK's calls do not take a context, so a cancelled login finishes (bounded by the client's
	// timeout) and is then logged out instead of being used; thread the context through the SDK when it matters
	for hops := 0; hops < 3; hops++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		key := HostKey(host, port)
		session, err := xenapi.NewSessionE(&xenapi.ClientOpts{URL: "https://" + key, HTTPClient: HTTPClient(policy, host, key, 60*time.Second)})
		if err != nil {
			return nil, err
		}
		_, err = session.LoginWithPassword(user, password, "2.21", "xcp-ng-center-go")
		var apiErr *xenapi.APIError
		if errors.As(err, &apiErr) && apiErr.Message == "HOST_IS_SLAVE" && len(apiErr.Data) > 0 {
			if session.Ref() != "" { // logged in before a follow-up call answered HOST_IS_SLAVE
				_ = session.Logout()
			}
			host = apiErr.Data[0] // the coordinator's address
			continue
		}
		if err != nil {
			// the login may have succeeded before the SDK's follow-up calls failed: never leave that session behind
			if session.Ref() != "" {
				_ = session.Logout()
			}
			return nil, Sanitize(err, password)
		}
		if err := ctx.Err(); err != nil {
			_ = session.Logout()
			return nil, err
		}
		return &Connection{Host: host, Port: port, User: user, Session: session}, nil
	}
	return nil, fmt.Errorf("too many coordinator redirects")
}

// Close logs out.
func (c *Connection) Close() error { return c.Session.Logout() }

// Sanitize keeps errors free of session references (OpaqueRef values) before they reach logs, the CLI or MCP; the
// wrapped error stays reachable with errors.As (an untrusted certificate, a XAPI failure).
func Sanitize(err error, secrets ...string) error {
	if err == nil {
		return nil
	}
	return &sanitized{msg: Redact(err.Error(), secrets...), err: err}
}

type sanitized struct {
	msg string
	err error
}

func (e *sanitized) Error() string { return e.msg }
func (e *sanitized) Unwrap() error { return e.err }
