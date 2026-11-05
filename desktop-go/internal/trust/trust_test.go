package trust

import "testing"

type memStore map[string]string

func (m memStore) Fingerprint(k string) (string, bool) { v, ok := m[k]; return v, ok }
func (m memStore) Pin(k, fp string) error              { m[k] = fp; return nil }

func TestPolicy(t *testing.T) {
	s := memStore{"pinned:443": "AA"}
	p := &Policy{Store: s, OperatorPins: map[string]string{}}
	if p.Evaluate("pinned:443", "AA", false) != Trusted {
		t.Error("the pinned certificate is trusted")
	}
	if p.Evaluate("pinned:443", "BB", true) != Mismatch {
		t.Error("a changed certificate is a mismatch even if the system trusts it")
	}
	if p.Evaluate("new:443", "CC", false) != Prompt {
		t.Error("an unknown self-signed certificate needs the user")
	}
	if p.Evaluate("new:443", "CC", true) != Trusted {
		t.Error("a system-trusted certificate on first use")
	}
	p.RequirePin = true
	if p.Evaluate("new:443", "CC", true) != Prompt {
		t.Error("RequirePin ignores the system's verdict")
	}
	p.OperatorPins["*"] = "c:c"
	if p.Evaluate("new:443", "CC", false) != Trusted || p.Evaluate("new:443", "DD", false) != Mismatch {
		t.Error("the operator's fingerprint (normalized) decides for unpinned hosts")
	}
	if fp, ok := p.Expected("new:443"); !ok || fp != "CC" {
		t.Errorf("Expected = %q, %v: the operator's pin, so a mismatch is reported as a change, never offered", fp, ok)
	}
	if Normalize(" ab:cd ") != "ABCD" {
		t.Error("Normalize")
	}
}
