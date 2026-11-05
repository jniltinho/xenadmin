// Package trust decides whether a server's TLS certificate is accepted (port of
// XenAdmin.Desktop.Core/CertificateTrust.cs): trust on first use, a pinned certificate must not change, and operator
// pins (XCP_TRUST_FINGERPRINT) for unattended use.
package trust

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
)

// Decision is the outcome for one certificate.
type Decision int

const (
	// Trusted: the pinned or vouched certificate, or (first use, no pin required) a certificate the system trusts.
	Trusted Decision = iota
	// Prompt: unknown certificate; the user must accept its fingerprint.
	Prompt
	// Mismatch: the certificate differs from the pinned one.
	Mismatch
)

// Store keeps the accepted fingerprints ("host:port" -> SHA-256 hex, upper case).
type Store interface {
	Fingerprint(hostKey string) (string, bool)
	Pin(hostKey, fingerprint string) error
}

// Fingerprint is the SHA-256 of the DER certificate, upper-case hex (as the C# app).
func Fingerprint(der []byte) string {
	sum := sha256.Sum256(der)
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

// Normalize accepts "AB:CD..." or "abcd..." forms.
func Normalize(fp string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(fp), ":", ""))
}

// Policy is CertificateTrust: the store, the operator's pins and whether a pin is always required (MCP, CLI scripts).
type Policy struct {
	Store        Store
	OperatorPins map[string]string // hostKey or "*" -> fingerprint
	RequirePin   bool
	mu           sync.Mutex
}

// Evaluate is CertificateTrust.Evaluate: a stored pin wins, then an operator pin, then the system's verdict.
func (p *Policy) Evaluate(hostKey, fp string, systemTrusted bool) Decision {
	p.mu.Lock()
	defer p.mu.Unlock()
	if pinned, ok := p.Store.Fingerprint(hostKey); ok {
		if strings.EqualFold(pinned, fp) {
			return Trusted
		}
		return Mismatch
	}
	if given, ok := p.operatorPin(hostKey); ok {
		if strings.EqualFold(Normalize(given), fp) {
			return Trusted
		}
		return Mismatch
	}
	if systemTrusted && !p.RequirePin {
		return Trusted
	}
	return Prompt
}

// Expected is the fingerprint a host must present: its stored pin, else the operator's pin (none: trust on first use).
func (p *Policy) Expected(hostKey string) (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if pinned, ok := p.Store.Fingerprint(hostKey); ok {
		return pinned, true
	}
	if given, ok := p.operatorPin(hostKey); ok {
		return Normalize(given), true
	}
	return "", false
}

func (p *Policy) operatorPin(hostKey string) (string, bool) {
	if v, ok := p.OperatorPins[hostKey]; ok {
		return v, true
	}
	v, ok := p.OperatorPins["*"]
	return v, ok
}

// UntrustedError is returned by a TLS handshake the policy did not accept: the UI or CLI shows the fingerprint and,
// once the user accepts it, pins it and connects again.
type UntrustedError struct {
	HostKey     string
	Fingerprint string
	Previous    string // the pinned fingerprint, when the certificate changed
}

func (e *UntrustedError) Error() string {
	if e.Previous != "" {
		return fmt.Sprintf("the certificate of %s changed (now %s, was %s): verify it before trusting it", e.HostKey, e.Fingerprint, e.Previous)
	}
	return fmt.Sprintf("the certificate of %s is not trusted yet (SHA-256 %s)", e.HostKey, e.Fingerprint)
}
