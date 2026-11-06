package rules

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1" //nolint:gosec // a certificate's thumbprint, and the C# app's SHA-1 signatures
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
// shortcut: no revocation check (X509Certificate2.Verify asks the CRLs online) and RSA keys only (as the C# app's
// RSACryptoServiceProvider); add OCSP/CRL checks if a signer is ever revoked.

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

// trusted is X509Certificate2.Verify: a chain to a trusted root (roots, nil for the system's) valid now, for any use.
func trusted(c *x509.Certificate, chain []*x509.Certificate, roots *x509.CertPool) bool {
	inter := x509.NewCertPool()
	for _, ca := range chain {
		inter.AddCert(ca)
	}
	_, err := c.Verify(x509.VerifyOptions{Roots: roots, Intermediates: inter, KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny}})
	return err == nil
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
	if !trusted(c, packageChain(p.certificate), roots) {
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
	if !trusted(cert, chain, roots) { // CheckCertificateValid
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
