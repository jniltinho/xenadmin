package rules

import (
	"archive/tar"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1" //nolint:gosec // the C# app's SHA-1 signatures
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"software.sslmate.com/src/go-pkcs12"

	"xcpng-center/internal/messages"
)

// testSigner is a throwaway self-signed certificate and its key, as a .pfx with password "pw", and a pool trusting it.
func testSigner(t *testing.T, dir string, alg x509.SignatureAlgorithm) (string, *x509.Certificate, *rsa.PrivateKey, *x509.CertPool) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(7), Subject: pkix.Name{CommonName: "OVF test signer", Organization: []string{"Lab"}},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour), IsCA: true, BasicConstraintsValid: true, SignatureAlgorithm: alg,
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, _ := x509.ParseCertificate(der)
	pfx, err := pkcs12.Modern.Encode(key, cert, nil, "pw")
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(dir, "signer.pfx")
	if err := os.WriteFile(file, pfx, 0o600); err != nil {
		t.Fatal(err)
	}
	roots := x509.NewCertPool()
	roots.AddCert(cert)
	return file, cert, key, roots
}

func TestLoadSigner(t *testing.T) {
	dir := t.TempDir()
	file, _, _, roots := testSigner(t, dir, x509.SHA256WithRSA)
	if _, err := LoadSigner(file, "pw", roots); err != nil {
		t.Fatal(err)
	}
	other, _, _, _ := testSigner(t, t.TempDir(), x509.SHA256WithRSA) // a certificate the pool does not trust
	for _, c := range []struct{ file, password, want string }{
		{"", "pw", "EXPORT_SECURITY_PAGE_ERROR_INVALID_CERT"},
		{"signer.pfx", "pw", "EXPORT_SECURITY_PAGE_ERROR_INVALID_CERT"},
		{filepath.Join(dir, "none.pfx"), "pw", "EXPORT_SECURITY_PAGE_ERROR_CERT_NON_EXIST"},
		{dir, "pw", "EXPORT_SECURITY_PAGE_ERROR_CERT_NON_EXIST"},
		{file, "wrong", "EXPORT_SECURITY_PAGE_ERROR_SIGN_INVALID"},
		{filepath.Join(dir, "signer.pfx.txt"), "pw", "EXPORT_SECURITY_PAGE_ERROR_SIGN_INVALID"},
		{other, "pw", "EXPORT_SECURITY_PAGE_ERROR_CERTIFICATE_INVALID"},
	} {
		if strings.HasSuffix(c.file, ".txt") {
			_ = os.WriteFile(c.file, []byte("not a pfx"), 0o600)
		}
		if _, err := LoadSigner(c.file, c.password, roots); err == nil || err.Error() != messages.Get(c.want) {
			t.Errorf("%s: %v, want %s", c.file, err, c.want)
		}
	}
	// the system's roots do not trust a self-signed certificate (X509Certificate2.Verify)
	if _, err := LoadSigner(file, "pw", nil); err == nil {
		t.Error("a self-signed certificate trusted by the system")
	}
	// ComputeSignedHash refuses a certificate that is neither SHA-256 nor SHA-1
	f384, _, _, roots384 := testSigner(t, t.TempDir(), x509.SHA384WithRSA)
	if _, err := LoadSigner(f384, "pw", roots384); err == nil || !strings.Contains(err.Error(), "SHA384-RSA") {
		t.Errorf("a SHA-384 certificate: %v", err)
	}
}

// TestSignedPackage signs exported packages (folder and OVA) and checks them back: the C# file names, order and line,
// the signature hashed twice, and a changed manifest or disk refused.
func TestSignedPackage(t *testing.T) {
	keys := t.TempDir()
	pfx, cert, key, roots := testSigner(t, keys, x509.SHA256WithRSA)
	signer, err := LoadSigner(pfx, "pw", roots)
	if err != nil {
		t.Fatal(err)
	}
	for _, isOVA := range []bool{false, true} {
		dir := t.TempDir()
		vhd := testVHD(4 << 20)
		if err := os.WriteFile(filepath.Join(dir, "d.vhd"), vhd, 0o644); err != nil {
			t.Fatal(err)
		}
		vm := ExportVM{Name: "x", VCPUs: 1, VCPUsMax: 1, Memory: 256 << 20, HardwarePlatformVersion: -1,
			Disks: []ExportDisk{{Href: "d.vhd", Name: "x", Device: "0", Size: int64(len(vhd)), Capacity: 4 << 20}}}
		// no manifest asked: signing makes it (ManifestAndSign)
		file, err := WritePackage(dir, "x", []ExportVM{vm}, nil, false, isOVA, false, signer)
		if err != nil {
			t.Fatal(err)
		}
		if isOVA { // ConvertOVFtoOVA's order: the descriptor, the manifest, the certificate, the disks
			f, _ := os.Open(filepath.Join(dir, file))
			tr := tar.NewReader(f)
			var names []string
			for h, err := tr.Next(); err == nil; h, err = tr.Next() {
				names = append(names, h.Name)
			}
			f.Close()
			if strings.Join(names, " ") != "x.ovf x.mf x.cert d.vhd" {
				t.Fatalf("OVA order %v", names)
			}
		} else {
			raw, err := os.ReadFile(filepath.Join(dir, "x.cert"))
			if err != nil {
				t.Fatal(err)
			}
			mf, _ := os.ReadFile(filepath.Join(dir, "x.mf"))
			first, rest, _ := strings.Cut(string(raw), "\n")
			m := regexp.MustCompile(`^SHA256\(x\.mf\)= ([0-9a-f]+)$`).FindStringSubmatch(first)
			if m == nil || !strings.HasPrefix(rest, "-----BEGIN CERTIFICATE-----\n") {
				t.Fatalf("certificate file %q", raw)
			}
			var sig []byte
			fmt.Sscanf(m[1], "%x", &sig)
			once := sha256.Sum256(mf)
			if twice := sha256.Sum256(once[:]); rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, twice[:], sig) != nil {
				t.Fatal("not the C# app's signature (SignData over the manifest's SHA-256)")
			}
		}
		p, err := OpenPackage(filepath.Join(dir, file))
		if err != nil {
			t.Fatal(err)
		}
		a := p.Appliance
		if !a.HasSignature || !a.HasManifest || a.Certificate == nil || a.CertificateError != "" ||
			a.Certificate.Subject != "CN=OVF test signer,O=Lab" || len(a.Certificate.Thumbprint) != 40 || !a.Certificate.ValidTo.Equal(cert.NotAfter) {
			t.Fatalf("%+v %+v", a, a.Certificate)
		}
		if err := p.VerifySignature(roots); err != nil {
			t.Fatal(err)
		}
		if err := p.VerifyManifest(); err != nil {
			t.Fatal(err)
		}
		if err := p.VerifySignature(nil); err == nil || err.Error() != "Certificate is invalid." {
			t.Errorf("untrusted signer: %v", err)
		}
		if isOVA {
			continue
		}
		// a changed manifest fails the signature; a changed disk with its manifest kept fails the manifest
		mf, _ := os.ReadFile(filepath.Join(dir, "x.mf"))
		_ = os.WriteFile(filepath.Join(dir, "x.mf"), append(mf, "SHA256(y)= 00\n"...), 0o644)
		if p, _ := OpenPackage(filepath.Join(dir, file)); p.VerifySignature(roots) == nil {
			t.Error("a changed manifest")
		}
		_ = os.WriteFile(filepath.Join(dir, "x.mf"), mf, 0o644)
		vhd[100] ^= 1
		_ = os.WriteFile(filepath.Join(dir, "d.vhd"), vhd, 0o644)
		if p, _ := OpenPackage(filepath.Join(dir, file)); p.VerifySignature(roots) != nil || p.VerifyManifest() == nil {
			t.Error("a changed disk")
		}
	}
}

// TestSignedChain: a signer under an intermediate CA (in its .pfx) verifies against the root alone, the intermediate
// written after the signer's certificate.
func TestSignedChain(t *testing.T) {
	dir := t.TempDir()
	issue := func(cn string, ca bool, parent *x509.Certificate, parentKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey) {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		tmpl := &x509.Certificate{SerialNumber: big.NewInt(time.Now().UnixNano()), Subject: pkix.Name{CommonName: cn},
			NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour), IsCA: ca, BasicConstraintsValid: true,
			KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign}
		if parent == nil {
			parent, parentKey = tmpl, key
		}
		der, err := x509.CreateCertificate(rand.Reader, tmpl, parent, &key.PublicKey, parentKey)
		if err != nil {
			t.Fatal(err)
		}
		c, _ := x509.ParseCertificate(der)
		return c, key
	}
	root, rootKey := issue("root", true, nil, nil)
	inter, interKey := issue("intermediate", true, root, rootKey)
	leaf, leafKey := issue("signer", false, inter, interKey)
	pfx, err := pkcs12.Modern.Encode(leafKey, leaf, []*x509.Certificate{inter}, "pw")
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(dir, "s.pfx"), pfx, 0o600)
	roots := x509.NewCertPool()
	roots.AddCert(root)
	s, err := LoadSigner(filepath.Join(dir, "s.pfx"), "pw", roots)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(dir, "tiny.ovf"), []byte(rawOVF), 0o600)
	_ = os.WriteFile(filepath.Join(dir, "disk.img"), []byte("d"), 0o600)
	if err := WriteManifest(dir, "tiny", []string{"tiny.ovf", "disk.img"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Sign(dir, "tiny"); err != nil {
		t.Fatal(err)
	}
	p, err := OpenPackage(filepath.Join(dir, "tiny.ovf"))
	if err != nil {
		t.Fatal(err)
	}
	if err := p.VerifySignature(roots); err != nil || p.Appliance.Certificate.Subject != "CN=signer" {
		t.Fatalf("%v %+v", err, p.Appliance.Certificate)
	}
	if err := p.VerifySignature(x509.NewCertPool()); err == nil { // the intermediate in the file is not a root
		t.Error("no trusted root")
	}
}

// TestSignatureFormats: a manifest signed once (ovftool, the OVF specification), a C# SHA-1 signature, the C# app's
// serialized certificate store element in the PEM block, and certificate files that cannot be read.
func TestSignatureFormats(t *testing.T) {
	dir := t.TempDir()
	_, cert, key, roots := testSigner(t, t.TempDir(), x509.SHA256WithRSA)
	mf := "SHA256(tiny.ovf)= 00\n"
	sign := func(h crypto.Hash, sum []byte) string {
		sig, err := rsa.SignPKCS1v15(rand.Reader, key, h, sum)
		if err != nil {
			t.Fatal(err)
		}
		return fmt.Sprintf("%x", sig)
	}
	s256, s1 := sha256.Sum256([]byte(mf)), sha1.Sum([]byte(mf)) //nolint:gosec // the C# SHA-1 signature
	derPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}))
	// X509ContentType.SerializedCert: a SHA-1 hash property (3) then the certificate (32), on one base64 line
	var ser []byte
	for _, prop := range []struct {
		id   uint32
		data []byte
	}{{3, make([]byte, 20)}, {certPropCert, cert.Raw}} {
		ser = binary.LittleEndian.AppendUint32(ser, prop.id)
		ser = binary.LittleEndian.AppendUint32(ser, 1)
		ser = binary.LittleEndian.AppendUint32(ser, uint32(len(prop.data)))
		ser = append(ser, prop.data...)
	}
	serPEM := "-----BEGIN CERTIFICATE-----\r\n" + base64.StdEncoding.EncodeToString(ser) + "\r\n-----END CERTIFICATE-----\r\n"
	for name, c := range map[string]struct{ cert, want string }{
		"once":       {"SHA256(tiny.mf)= " + sign(crypto.SHA256, s256[:]) + "\n" + derPEM, ""},
		"sha1":       {"SHA1(tiny.mf)= " + sign(crypto.SHA1, s1[:]) + "\r\n" + serPEM, ""},
		"wrong hash": {"SHA1(tiny.mf)= " + sign(crypto.SHA256, s256[:]) + "\n" + derPEM, "does not match"},
		"no line":    {derPEM, "not valid"},
		"no cert":    {"SHA256(tiny.mf)= 00\n", messages.Get("IMPORT_SECURITY_PAGE_ERROR_MISSING_CERTIFICATE")},
		"bad cert":   {"SHA256(tiny.mf)= 00\n-----BEGIN CERTIFICATE-----\nAAAA\n-----END CERTIFICATE-----\n", messages.Get("IMPORT_SECURITY_PAGE_ERROR_MISSING_CERTIFICATE")},
	} {
		file := filepath.Join(dir, name+".ova")
		ova(t, file, false, [2]string{"tiny.ovf", rawOVF}, [2]string{"tiny.mf", mf}, [2]string{"tiny.cert", c.cert}, [2]string{"disk.img", "d"})
		p, err := OpenPackage(file)
		if err != nil {
			t.Fatal(err)
		}
		if !p.Appliance.HasSignature || (p.Appliance.CertificateError != "") != strings.HasPrefix(c.want, "Could not") {
			t.Errorf("%s: %+v", name, p.Appliance)
		}
		if err := p.VerifySignature(roots); c.want == "" && err != nil || c.want != "" && (err == nil || !strings.Contains(err.Error(), c.want)) {
			t.Errorf("%s: %v, want %q", name, err, c.want)
		}
	}
	// an unsigned package
	file := filepath.Join(dir, "plain.ova")
	ova(t, file, false, [2]string{"tiny.ovf", rawOVF}, [2]string{"disk.img", "d"})
	if p, err := OpenPackage(file); err != nil || p.Appliance.HasSignature || p.VerifySignature(roots) == nil {
		t.Errorf("unsigned: %v", err)
	}
}
