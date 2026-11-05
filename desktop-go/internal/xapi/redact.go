package xapi

import (
	"regexp"
	"strings"
)

var opaqueRef = regexp.MustCompile(`OpaqueRef:[0-9a-fA-F-]{8,}`)

// Redact hides session and object references (a session reference is a bearer token) and the given secrets (the
// password), wherever a server or library echoed them (XcpTools.Redact in the C# app).
func Redact(s string, secrets ...string) string {
	for _, secret := range secrets {
		if secret != "" { // any length: a short password is still a password
			s = strings.ReplaceAll(s, secret, "***")
		}
	}
	return opaqueRef.ReplaceAllString(s, "OpaqueRef:…")
}
