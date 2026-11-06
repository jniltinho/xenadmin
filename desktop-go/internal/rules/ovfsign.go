package rules

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1" //nolint:gosec // a certificate's thumbprint, and the C# app's SHA-1 signatures
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"software.sslmate.com/src/go-pkcs12"

	"xcpng-center/internal/messages"
)

// Signed OVF packages as the C# app makes and checks them: ExportApplianceAction.Sign writes "<name>.cert", a
// FileDigest line with the signature of "<name>.mf" (StreamUtilities.ComputeSignedHash), then the signer's certificate
// in PEM; Package.VerifySignature checks the certificate (X509Certificate2.Verify) then the signature
// (StreamUtilities.VerifyAgainstDigest), and the import checks the manifest after it.
// The C# app's SHA-256 signature is RSACryptoServiceProvider.SignData over the manifest's SHA-256, so the manifest is
// hashed twice; ovftool and the OVF specification sign the manifest's hash once. Both are accepted, the C# one is
// written. The C# app puts X509ContentType.SerializedCert (a Windows certificate store element) in the PEM; this app
// writes the DER certificate (what the C# comment says it writes, which X509Certificate2 and other tools read) and
// reads both. go-pkcs12 reads the .pfx: golang.org/x/crypto/pkcs12 cannot open the AES (PBES2) files that OpenSSL 3
// and current Windows write.
// shortcut: RSA keys only (as the C# app's RSACryptoServiceProvider).

// CertificateInfo is what View Certificate shows of the signer (X509Certificate2UI's General page, and the thumbprint).
type CertificateInfo struct {
	Subject    string    `json:"subject"`
	Issuer     string    `json:"issuer"`
	ValidFrom  time.Time `json:"valid_from"`
	ValidTo    time.Time `json:"valid_to"`
	Thumbprint string    `json:"thumbprint"` // the certificate's SHA-1, as Windows shows it
}

func describeCertificate(c *x509.Certificate) *CertificateInfo {
	sum := sha1.Sum(c.Raw) //nolint:gosec // the thumbprint
	return &CertificateInfo{Subject: c.Subject.String(), Issuer: c.Issuer.String(), ValidFrom: c.NotBefore, ValidTo: c.NotAfter,
		Thumbprint: strings.ToUpper(hex.EncodeToString(sum[:]))}
}

const certPropCert = 32 // CERT_CERT_PROP_ID: the certificate itself in a serialized store element

// packageCertificate is the signer's certificate in a "<name>.cert": the PEM block after the signature line, DER or a
// serialized store element (properties as id, encoding, length and data).
func packageCertificate(raw []byte) (*x509.Certificate, error) {
	b, _ := pem.Decode(raw)
	if b == nil || b.Type != "CERTIFICATE" {
		return nil, errors.New("the certificate file has no certificate")
	}
	if c, err := x509.ParseCertificate(b.Bytes); err == nil {
		return c, nil
	}
	for d := b.Bytes; len(d) >= 12; {
		id, n := binary.LittleEndian.Uint32(d), binary.LittleEndian.Uint32(d[8:])
		if uint64(n) > uint64(len(d)-12) {
			break
		}
		if id == certPropCert {
			return x509.ParseCertificate(d[12 : 12+n])
		}
		d = d[12+n:]
	}
	return nil, errors.New("the certificate cannot be read")
}

// packageChain is the intermediate certificates this app writes after the signer's (the C# app reads the first
// certificate only and finds the others in the system's stores); they help build the chain but are not trusted.
func packageChain(raw []byte) []*x509.Certificate {
	var out []*x509.Certificate
	_, rest := pem.Decode(raw)
	for b, r := pem.Decode(rest); b != nil; b, r = pem.Decode(r) {
		if c, err := x509.ParseCertificate(b.Bytes); err == nil && b.Type == "CERTIFICATE" {
			out = append(out, c)
		}
	}
	return out
}

// trusted is X509Certificate2.Verify: a chain to a trusted root (roots, nil for the system's) valid now, for any use,
// whose certificates are not revoked (checkRevocation).
func trusted(c *x509.Certificate, chain []*x509.Certificate, roots *x509.CertPool) error {
	inter := x509.NewCertPool()
	for _, ca := range chain {
		inter.AddCert(ca)
	}
	chains, err := c.Verify(x509.VerifyOptions{Roots: roots, Intermediates: inter, KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny}})
	if err != nil {
		return err
	}
	for _, ch := range chains { // X509Chain builds one chain: any of these that passes
		if err = checkRevocation(ch, time.Now()); err == nil {
			return nil
		}
	}
	return err
}

// CRLFetcher gets the CRL at a certificate's CRL distribution point.
type CRLFetcher interface {
	FetchCRL(url string) ([]byte, error)
}

const crlMaxSize = 16 << 20 // a large public CA's CRL is a few MB

type httpCRLFetcher struct{ client *http.Client }

func (f httpCRLFetcher) FetchCRL(url string) ([]byte, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the CRL at %s cannot be read: %s", url, resp.Status)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, crlMaxSize+1))
	if err == nil && len(b) > crlMaxSize {
		err = fmt.Errorf("the CRL at %s is too large", url)
	}
	return b, err
}

// crls fetches the CRLs (a test may put another fetcher here).
var crls CRLFetcher = httpCRLFetcher{&http.Client{Timeout: 10 * time.Second}}

// checkRevocation is X509Certificate2.Verify's revocation check. Its chain policy is X509RevocationMode.Online,
// X509RevocationFlag.ExcludeRoot and X509VerificationFlags.NoFlag, so every certificate of the chain but the root must
// be found not revoked in a current CRL of one of its distribution points, signed by its issuer: a revoked certificate
// fails, and so does one whose revocation cannot be told (RevocationStatusUnknown or OfflineRevocation: no
// distribution point, none reachable, no valid CRL there). The CRLs are untrusted: HTTP(S) only, with a timeout and a
// size cap, and the issuer's signature checked.
// shortcut: CRLs only, no OCSP (Windows also asks the OCSP responder in a certificate's AIA, so a certificate with an
// OCSP responder and no CRL distribution point passes there and fails here); add OCSP if such a signer shows up.
// shortcut: complete CRLs only: a delta CRL, a partitioned one (an issuing distribution point), one that names delta
// CRLs (Freshest CRL) or one with another critical extension is not used, so a CA that publishes delta CRLs fails
// here (Windows reads them); add the delta CRLs if such a signer shows up.
func checkRevocation(chain []*x509.Certificate, now time.Time) error {
	for i := 0; i+1 < len(chain); i++ {
		if err := notRevoked(chain[i], chain[i+1], now); err != nil {
			return err
		}
	}
	return nil
}

// partialCRL tells whether a CRL's extension makes it list only some revocations: the delta CRL indicator, the issuing
// distribution point, the Freshest CRL (delta CRLs to read with it), or any critical extension Go does not read (it
// reads the authority key id and CRL number).
func partialCRL(e pkix.Extension) bool {
	id := e.Id.String()
	return id == "2.5.29.27" || id == "2.5.29.28" || id == "2.5.29.46" || e.Critical && id != "2.5.29.35" && id != "2.5.29.20"
}

func notRevoked(c, issuer *x509.Certificate, now time.Time) error {
	for _, e := range c.Extensions {
		if e.Id.String() == "2.5.29.46" { // Freshest CRL: its delta CRLs would be needed too (see checkRevocation)
			return fmt.Errorf("the revocation of the certificate %s cannot be checked: it names delta CRLs", c.Subject)
		}
	}
	err := fmt.Errorf("the revocation of the certificate %s cannot be checked: it has no HTTP CRL distribution point", c.Subject)
	for _, url := range c.CRLDistributionPoints {
		if u := strings.ToLower(url); !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
			continue
		}
		der, ferr := crls.FetchCRL(url)
		if ferr != nil {
			err = ferr
			continue
		}
		if b, _ := pem.Decode(der); b != nil && b.Type == "X509 CRL" {
			der = b.Bytes
		}
		crl, perr := x509.ParseRevocationList(der)
		switch {
		case perr != nil:
			err = fmt.Errorf("the CRL at %s is not valid: %w", url, perr)
		case slices.ContainsFunc(crl.Extensions, partialCRL):
			err = fmt.Errorf("the CRL at %s is a delta or partitioned CRL, which cannot be checked", url)
		case !bytes.Equal(crl.RawIssuer, c.RawIssuer) || crl.CheckSignatureFrom(issuer) != nil:
			err = fmt.Errorf("the CRL at %s is not signed by the issuer of %s", url, c.Subject)
		case now.Before(crl.ThisUpdate) || !crl.NextUpdate.IsZero() && now.After(crl.NextUpdate):
			err = fmt.Errorf("the CRL at %s is not current", url)
		default:
			for _, r := range crl.RevokedCertificateEntries {
				if r.SerialNumber != nil && r.SerialNumber.Cmp(c.SerialNumber) == 0 {
					return fmt.Errorf("the certificate %s is revoked", c.Subject)
				}
			}
			return nil
		}
	}
	return err
}

func signatureHash(algorithm string) (crypto.Hash, error) {
	switch strings.ReplaceAll(strings.ToUpper(algorithm), "-", "") {
	case "SHA1":
		return crypto.SHA1, nil
	case "SHA256":
		return crypto.SHA256, nil
	case "SHA512":
		return crypto.SHA512, nil
	}
	return 0, fmt.Errorf("the signature's algorithm %s is not supported", algorithm)
}

// VerifySignature is Package.VerifySignature: the signer's certificate is trusted (roots, nil for the system's) and the
// signature on the certificate file's first line is the manifest's.
func (p *Package) VerifySignature(roots *x509.CertPool) error {
	if p.certificate == nil {
		return errors.New("the appliance is not signed")
	}
	c, err := packageCertificate(p.certificate)
	if err != nil {
		return errors.New(messages.Get("IMPORT_SECURITY_PAGE_ERROR_MISSING_CERTIFICATE"))
	}
	if trusted(c, packageChain(p.certificate), roots) != nil {
		return errors.New("Certificate is invalid.") // XenOvf's CERTIFICATE_IS_INVALID
	}
	if p.manifest == nil {
		return errors.New("the appliance has no manifest")
	}
	first, _, _ := strings.Cut(strings.ReplaceAll(string(p.certificate), "\r", ""), "\n")
	m := manifestLine.FindStringSubmatch(first)
	if m == nil {
		return fmt.Errorf("the signature line %q is not valid", first)
	}
	name := strings.TrimSpace(m[2])
	h, err := signatureHash(m[1])
	if err != nil {
		return err
	}
	if len(m[3])%2 == 1 { // FileDigest.ToArray: an odd digest has a leading zero
		m[3] = "0" + m[3]
	}
	sig, _ := hex.DecodeString(m[3])
	pub, ok := c.PublicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("the certificate's key is not an RSA key")
	}
	hh := h.New()
	hh.Write(p.manifest)
	sum := hh.Sum(nil)
	if rsa.VerifyPKCS1v15(pub, h, sum, sig) == nil {
		return nil
	}
	if h == crypto.SHA256 { // the C# app's: RSACryptoServiceProvider.VerifyData hashes the manifest's hash again
		if twice := sha256.Sum256(sum); rsa.VerifyPKCS1v15(pub, h, twice[:], sig) == nil {
			return nil
		}
	}
	return fmt.Errorf("The file %s does not match the original.", name) // XenOvf's SECURITY_SIGNATURE_FAILED
}

// Signer is the Export wizard's certificate (ExportOptionsPage.Certificate) with its private key.
type Signer struct {
	cert  *x509.Certificate
	key   *rsa.PrivateKey
	hash  crypto.Hash
	alg   string              // the hash's name in the FileDigest line
	chain []*x509.Certificate // the .pfx's other certificates, written after the signer's
}

// LoadSigner is ExportOptionsPage.CheckSignature: the file is a PKCS#12 (.pfx) that the password opens, with an RSA
// key of its certificate, which is trusted (roots, nil for the system's) and valid now. The password is used only here.
func LoadSigner(file, password string, roots *x509.CertPool) (*Signer, error) {
	if strings.TrimSpace(file) == "" || !filepath.IsAbs(file) { // CheckCertificatePathValid
		return nil, errors.New(messages.Get("EXPORT_SECURITY_PAGE_ERROR_INVALID_CERT"))
	}
	if st, err := os.Stat(file); err != nil || !st.Mode().IsRegular() { // CheckCertificatePathExists
		return nil, errors.New(messages.Get("EXPORT_SECURITY_PAGE_ERROR_CERT_NON_EXIST"))
	}
	invalid := errors.New(messages.Get("EXPORT_SECURITY_PAGE_ERROR_SIGN_INVALID")) // CheckSignPasswordValid
	data, err := readLimited(file, ovfMaxManifest)
	if err != nil {
		return nil, invalid
	}
	key, cert, chain, err := pkcs12.DecodeChain(data, password)
	if err != nil {
		return nil, invalid
	}
	k, ok := key.(*rsa.PrivateKey)
	if !ok || !k.PublicKey.Equal(cert.PublicKey) {
		return nil, invalid
	}
	if trusted(cert, chain, roots) != nil { // CheckCertificateValid: Certificate.Verify(), revocation included
		return nil, errors.New(messages.Get("EXPORT_SECURITY_PAGE_ERROR_CERTIFICATE_INVALID"))
	}
	s := &Signer{cert: cert, key: k, chain: chain}
	// ComputeSignedHash: the certificate's own hash, SHA-256 or SHA-1 (refused there at signing, here before the export)
	switch cert.SignatureAlgorithm {
	case x509.SHA256WithRSA:
		s.hash, s.alg = crypto.SHA256, "SHA256"
	case x509.SHA1WithRSA:
		s.hash, s.alg = crypto.SHA1, "SHA1"
	default:
		return nil, fmt.Errorf("The certificate's signature algorithm %s is not supported: only SHA-256 and SHA-1 with RSA are.", cert.SignatureAlgorithm)
	}
	return s, nil
}

// Sign is ExportApplianceAction.Sign: "<name>.cert" in dir, the FileDigest line of the signed hash of "<name>.mf"
// (ComputeSignedHash: SHA-256 hashed twice, or SHA-1), then the certificate in PEM.
func (s *Signer) Sign(dir, name string) error {
	mf, err := readLimited(filepath.Join(dir, name+".mf"), ovfMaxManifest)
	if err != nil {
		return err
	}
	h, alg := s.hash, s.alg
	hh := h.New()
	hh.Write(mf)
	sum := hh.Sum(nil)
	if h == crypto.SHA256 {
		twice := sha256.Sum256(sum)
		sum = twice[:]
	}
	sig, err := rsa.SignPKCS1v15(rand.Reader, s.key, h, sum)
	if err != nil {
		return err
	}
	out := fmt.Sprintf("%s(%s.mf)= %x\n", alg, name, sig)
	for _, c := range append([]*x509.Certificate{s.cert}, s.chain...) {
		out += string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.Raw}))
	}
	return writeSynced(filepath.Join(dir, name+".cert"), func(w io.Writer) error { _, err := io.WriteString(w, out); return err })
}
