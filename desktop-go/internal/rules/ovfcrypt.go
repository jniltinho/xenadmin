package rules

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode/utf16"

	"xcpng-center/internal/messages"
)

// Encrypted OVF packages (XenOvf's Security.cs), as older XenCenters made them: a SecuritySection whose Security
// elements list the encrypted files (DataReference ValueType: their ids), a version, and an EncryptedData holding a
// known sentence encrypted with the password, which CheckPassword decrypts to check it. The files are decrypted as they
// are read, before they are uncompressed (Import.ImportFile: DecryptToTempFile, then UncompressFile).
// The scheme, as the C# code runs it for a package with a version: CheckSecurityVersion compares the version with
// itself (it reads both sides from its first argument), so every version is ">= 1.3.1", and CryptoStreamWrapper uses
// RijndaelManaged whatever EncryptionMethod's Algorithm says: AES-256-CBC (RijndaelManaged's default 256-bit key and
// 128-bit block) with PKCS7 padding, the key and IV the SHA-384 of the password's ASCII bytes (GenerateKey: the first
// 32 bytes, then the next 16). So the TripleDES and RC2 classes GetAlgorithmClass maps are never used for such a
// package. The Windows app itself refuses encrypted packages (ImportSourcePage) and never makes them (the export
// wizard has no encryption): this reads those of older XenCenters.
// shortcut: a package without a version (the deprecated scheme: DeprecatedCryptoFile drops the last block and pads
// with zeros, DeprecatedCheckPassword decodes UTF-16 as UTF-8) is refused; port it if such a package shows up.

// ovfKnownEncrypt is Security.cs' KnownEncrypt, trimmed.
const ovfKnownEncrypt = "Nihil tam munitum quod non expugnari pecunia possit."

// ovfSecurity is a Security element of a SecuritySection (Security_Type, with XENC's ReferenceList and EncryptedData).
type ovfSecurity struct {
	Version string `xml:"version,attr"`
	Refs    []struct {
		ValueType string `xml:"ValueType,attr"`
	} `xml:"ReferenceList>DataReference"`
	Cipher []string `xml:"EncryptedData>CipherData>CipherValue"`
	Typo   []string `xml:"EncrypteData>CipherData>CipherValue"` // CheckPassword reads this spelling too
}

type ovfSecuritySection struct {
	Security []ovfSecurity `xml:"Security"`
}

// encrypted is IsThisIdEncrypted: a DataReference of a Security names the file's id (a substring, as there).
func (e *ovfEnvelope) encrypted(id string) bool {
	for _, s := range e.Security {
		for _, sec := range s.Security {
			for _, r := range sec.Refs {
				if id != "" && strings.Contains(r.ValueType, id) {
					return true
				}
			}
		}
	}
	return false
}

// checkEncryption is what the import can decrypt: one SecuritySection (CheckPassword reads no other), each Security
// with a version (CheckSecurityVersion's Convert.ToInt32 on its first two parts) and an EncryptedData to check the
// password against.
func (e *ovfEnvelope) checkEncryption() error {
	invalid := errors.New("Security section is invalid.") // XenOvf's SECURITY_SECTION_INVALID
	if len(e.Security) != 1 || len(e.Security[0].Security) == 0 {
		return invalid
	}
	for _, s := range e.Security[0].Security {
		if strings.TrimSpace(s.Version) == "" {
			return errors.New("The appliance is encrypted with the scheme of an older version, which cannot be imported in this version.")
		}
		parts := strings.Split(s.Version, ".")
		for _, part := range parts[:min(2, len(parts))] {
			if _, err := strconv.ParseInt(strings.TrimSpace(part), 10, 32); err != nil {
				return invalid
			}
		}
		if len(s.Cipher)+len(s.Typo) == 0 {
			return invalid
		}
	}
	return nil
}

// ovfKey is GenerateKey for RijndaelManaged: the SHA-384 of the password's ASCII bytes (Encoding.ASCII: "?" for any
// other UTF-16 unit, so two for a character past the BMP) is the key (32 bytes) then the IV (16).
func ovfKey(password string) (key, iv []byte) {
	b := make([]byte, 0, len(password))
	for _, r := range password {
		switch {
		case r > 0xffff:
			b = append(b, '?', '?')
		case r > 0x7f:
			b = append(b, '?')
		default:
			b = append(b, byte(r))
		}
	}
	sum := sha512.Sum384(b)
	return sum[:32], sum[32:48]
}

// cbcReader decrypts AES-256-CBC as it reads (CryptoStream over a RijndaelManaged decryptor), holding the last block
// back to remove its PKCS7 padding at the end.
type cbcReader struct {
	r    io.Reader
	mode cipher.BlockMode
	in   []byte // ciphertext not decrypted yet (the last block, or less)
	out  []byte // plaintext ready
	buf  []byte
	err  error
}

func newOVFDecrypter(r io.Reader, password string) io.Reader {
	key, iv := ovfKey(password)
	block, _ := aes.NewCipher(key) // a 32-byte key: no error
	return &cbcReader{r: r, mode: cipher.NewCBCDecrypter(block, iv), buf: make([]byte, 1<<20+aes.BlockSize)}
}

var errOVFDecrypt = errors.New("the file cannot be decrypted: it is damaged or not encrypted with this password")

func (d *cbcReader) Read(p []byte) (int, error) {
	for len(d.out) == 0 {
		if d.err != nil {
			return 0, d.err
		}
		d.fill()
	}
	n := copy(p, d.out)
	d.out = d.out[n:]
	return n, nil
}

func (d *cbcReader) fill() {
	have := copy(d.buf, d.in)
	n, err := d.r.Read(d.buf[have:])
	have += n
	if errors.Is(err, io.EOF) {
		if have == 0 || have%aes.BlockSize != 0 {
			d.err = errOVFDecrypt
			return
		}
		d.mode.CryptBlocks(d.buf[:have], d.buf[:have])
		pad := int(d.buf[have-1])
		if pad == 0 || pad > aes.BlockSize {
			d.err = errOVFDecrypt
			return
		}
		for _, b := range d.buf[have-pad : have] {
			if int(b) != pad {
				d.err = errOVFDecrypt
				return
			}
		}
		d.out, d.in, d.err = d.buf[:have-pad], nil, io.EOF
		return
	}
	if err != nil {
		d.err = err
		return
	}
	// every whole block but the last (which may be the padding), and at least one byte held back
	ready := 0
	if have > 0 {
		ready = (have - 1) / aes.BlockSize * aes.BlockSize
	}
	d.mode.CryptBlocks(d.buf[:ready], d.buf[:ready])
	d.out = d.buf[:ready]
	d.in = append(d.in[:0], d.buf[ready:have]...)
	// out aliases buf: the next fill copies in to buf's start only once out is used up
}

// checkOVFPassword is InternalCheckPassword: the EncryptedData decrypts to KnownEncrypt (UTF-16, trimmed of spaces and
// NULs).
func checkOVFPassword(cipherValue, password string) bool {
	data, err := base64.StdEncoding.DecodeString(strings.Join(strings.Fields(cipherValue), ""))
	if err != nil {
		return false
	}
	plain, err := io.ReadAll(newOVFDecrypter(strings.NewReader(string(data)), password))
	if err != nil {
		return false
	}
	u := make([]uint16, len(plain)/2)
	for i := range u {
		u[i] = uint16(plain[2*i]) | uint16(plain[2*i+1])<<8
	}
	return strings.Trim(string(utf16.Decode(u)), " \x00") == ovfKnownEncrypt
}

// Unlock is ImportSecurityPage's password for an encrypted package, checked (CheckPassword) and kept in the package
// to decrypt its disks; it is never stored nor logged. Each Security's EncryptedData must take it (CheckPassword
// keeps the last one's answer: they all have the same password).
func (p *Package) Unlock(password string) error {
	if !p.Appliance.Encrypted {
		return nil
	}
	if password == "" {
		return errors.New(messages.Get("ERROR_NO_PASSWORD"))
	}
	for _, s := range p.security[0].Security {
		if !checkOVFPassword(append(s.Cipher, s.Typo...)[0], password) {
			return errors.New("The password is invalid.") // XenOvf's PASSWORD_FAILED
		}
	}
	p.password = password
	return nil
}
