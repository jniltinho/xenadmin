package rules

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// testCert is a throwaway certificate and its key: a root when parent is nil, else issued by parent.
type testCert struct {
	cert *x509.Certificate
	key  *rsa.PrivateKey
}

func newTestCert(t *testing.T, cn string, ca bool, parent *testCert, cdps ...string) *testCert {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(time.Now().UnixNano()), Subject: pkix.Name{CommonName: cn},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour), IsCA: ca, BasicConstraintsValid: true,
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign, CRLDistributionPoints: cdps}
	signer, signerKey := tmpl, key
	if parent != nil {
		signer, signerKey = parent.cert, parent.key
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, signer, &key.PublicKey, signerKey)
	if err != nil {
		t.Fatal(err)
	}
	c, _ := x509.ParseCertificate(der)
	return &testCert{c, key}
}

// crl is the CA's CRL (DER) listing the revoked certificates, current until next.
func (ca *testCert) crl(t *testing.T, next time.Time, revoked ...*x509.Certificate) []byte {
	t.Helper()
	var entries []x509.RevocationListEntry
	for _, c := range revoked {
		entries = append(entries, x509.RevocationListEntry{SerialNumber: c.SerialNumber, RevocationTime: time.Now().Add(-time.Minute)})
	}
	der, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{Number: big.NewInt(1), ThisUpdate: time.Now().Add(-time.Hour),
		NextUpdate: next, RevokedCertificateEntries: entries}, ca.cert, ca.key)
	if err != nil {
		t.Fatal(err)
	}
	return der
}

// crlServer serves the CRLs put in it by path, through the app's own HTTP fetcher.
type crlServer struct {
	*httptest.Server
	mu    sync.Mutex
	files map[string][]byte
}

func newCRLServer(t *testing.T) *crlServer {
	s := &crlServer{files: map[string][]byte{}}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		b, ok := s.files[r.URL.Path]
		s.mu.Unlock()
		if !ok {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write(b)
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *crlServer) put(path string, b []byte) {
	s.mu.Lock()
	s.files[path] = b
	s.mu.Unlock()
}

// TestRevocation: X509Certificate2.Verify's online revocation check of a certificate against its issuer's CRLs.
func TestRevocation(t *testing.T) {
	srv := newCRLServer(t)
	gone := httptest.NewServer(http.NotFoundHandler())
	unreachable := gone.URL + "/ca.crl"
	gone.Close()
	ca := newTestCert(t, "CA", true, nil)
	impostor := newTestCert(t, "CA", true, nil) // the same name, another key
	now, later := time.Now(), time.Now().Add(time.Hour)
	revokedLeaf := newTestCert(t, "revoked", false, ca, srv.URL+"/revoked.crl")
	srv.put("/ok.crl", ca.crl(t, later))
	srv.put("/revoked.crl", ca.crl(t, later, revokedLeaf.cert))
	srv.put("/pem.crl", pem.EncodeToMemory(&pem.Block{Type: "X509 CRL", Bytes: ca.crl(t, later)}))
	srv.put("/stale.crl", ca.crl(t, now.Add(-time.Minute)))
	srv.put("/impostor.crl", impostor.crl(t, later))
	srv.put("/garbage.crl", []byte("not a CRL"))
	srv.put("/big.crl", bytes.Repeat([]byte{0x30}, crlMaxSize+1))
	delta, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{Number: big.NewInt(2), ThisUpdate: now.Add(-time.Hour), NextUpdate: later,
		ExtraExtensions: []pkix.Extension{{Id: asn1.ObjectIdentifier{2, 5, 29, 27}, Critical: true, Value: []byte{2, 1, 1}}}}, ca.cert, ca.key)
	if err != nil {
		t.Fatal(err)
	}
	srv.put("/delta.crl", delta)
	freshest, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{Number: big.NewInt(3), ThisUpdate: now.Add(-time.Hour), NextUpdate: later,
		ExtraExtensions: []pkix.Extension{{Id: asn1.ObjectIdentifier{2, 5, 29, 46}, Value: []byte{0x30, 0}}}}, ca.cert, ca.key)
	if err != nil {
		t.Fatal(err)
	}
	srv.put("/freshest.crl", freshest)
	for name, c := range map[string]struct {
		cdps []string
		want string // "" when not revoked
	}{
		"not revoked":    {[]string{srv.URL + "/ok.crl"}, ""},
		"PEM":            {[]string{srv.URL + "/pem.crl"}, ""},
		"second point":   {[]string{unreachable, "ldap://ca/crl", srv.URL + "/ok.crl"}, ""},
		"unreachable":    {[]string{unreachable}, "connection refused"},
		"no point":       {nil, "no HTTP CRL distribution point"},
		"LDAP only":      {[]string{"ldap://ca/crl"}, "no HTTP CRL distribution point"},
		"not found":      {[]string{srv.URL + "/none.crl"}, "404"},
		"too large":      {[]string{srv.URL + "/big.crl"}, "too large"},
		"not a CRL":      {[]string{srv.URL + "/garbage.crl"}, "not valid"},
		"wrong issuer":   {[]string{srv.URL + "/impostor.crl"}, "not signed by the issuer"},
		"not current":    {[]string{srv.URL + "/stale.crl"}, "not current"},
		"delta":          {[]string{srv.URL + "/delta.crl"}, "delta"},
		"names deltas":   {[]string{srv.URL + "/freshest.crl"}, "delta"},
		"others revoked": {[]string{srv.URL + "/revoked.crl"}, ""},
	} {
		leaf := newTestCert(t, name, false, ca, c.cdps...)
		if err := notRevoked(leaf.cert, ca.cert, now); c.want == "" && err != nil || c.want != "" && (err == nil || !strings.Contains(err.Error(), c.want)) {
			t.Errorf("%s: %v, want %q", name, err, c.want)
		}
	}
	if err := notRevoked(revokedLeaf.cert, ca.cert, now); err == nil || !strings.Contains(err.Error(), "is revoked") {
		t.Errorf("revoked: %v", err)
	}
	// a certificate that names delta CRLs itself
	deltaTmpl := &x509.Certificate{SerialNumber: big.NewInt(99), Subject: pkix.Name{CommonName: "deltas"}, NotBefore: now.Add(-time.Hour),
		NotAfter: later, CRLDistributionPoints: []string{srv.URL + "/ok.crl"},
		ExtraExtensions: []pkix.Extension{{Id: asn1.ObjectIdentifier{2, 5, 29, 46}, Value: []byte{0x30, 0}}}}
	der, err := x509.CreateCertificate(rand.Reader, deltaTmpl, ca.cert, &ca.key.PublicKey, ca.key)
	if err != nil {
		t.Fatal(err)
	}
	if dc, _ := x509.ParseCertificate(der); notRevoked(dc, ca.cert, now) == nil {
		t.Error("a certificate naming delta CRLs")
	}
	// the root is not checked (ExcludeRoot), every certificate under it is
	leaf := newTestCert(t, "leaf", false, ca, srv.URL+"/ok.crl")
	if err := checkRevocation([]*x509.Certificate{leaf.cert, ca.cert}, now); err != nil {
		t.Errorf("a root without a distribution point: %v", err)
	}
	inter := newTestCert(t, "inter", true, ca)
	under := newTestCert(t, "under", false, inter, srv.URL+"/inter.crl")
	srv.put("/inter.crl", inter.crl(t, later))
	if err := checkRevocation([]*x509.Certificate{under.cert, inter.cert, ca.cert}, now); err == nil {
		t.Error("an intermediate without a distribution point")
	}
}
