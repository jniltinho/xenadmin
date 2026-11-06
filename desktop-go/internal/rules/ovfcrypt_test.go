package rules

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"testing/iotest"
	"unicode/utf16"

	"xcpng-center/internal/messages"
)

// csharpEncrypt is Security.cs' encryption (CryptoSetup with a version, CryptoStream, FlushFinalBlock) done by openssl
// alone: the key and IV the SHA-384 of the password (openssl dgst), AES-256-CBC with PKCS7 padding (openssl enc).
func csharpEncrypt(t *testing.T, password string, plain []byte) []byte {
	t.Helper()
	if _, err := exec.LookPath("openssl"); err != nil {
		t.Skip("no openssl")
	}
	dgst := exec.Command("openssl", "dgst", "-sha384", "-binary")
	dgst.Stdin = strings.NewReader(password)
	sum, err := dgst.Output()
	if err != nil || len(sum) != 48 {
		t.Fatalf("openssl dgst: %v", err)
	}
	enc := exec.Command("openssl", "enc", "-aes-256-cbc", "-K", hex.EncodeToString(sum[:32]), "-iv", hex.EncodeToString(sum[32:]))
	enc.Stdin = bytes.NewReader(plain)
	out, err := enc.Output()
	if err != nil {
		t.Fatalf("openssl enc: %v", err)
	}
	return out
}

// csharpCheckValue is CryptoElement(KnownEncrypt): the sentence (with its trailing spaces) in UTF-16, encrypted.
func csharpCheckValue(t *testing.T, password string) string {
	var b []byte
	for _, u := range utf16.Encode([]rune(ovfKnownEncrypt + strings.Repeat(" ", 46))) {
		b = append(b, byte(u), byte(u>>8))
	}
	return base64.StdEncoding.EncodeToString(csharpEncrypt(t, password, b))
}

// encryptedOVF is a descriptor as Security.cs' CryptoFileWrapper leaves it: a SecuritySection with one Security that
// lists the encrypted files' ids, its version and the password's check value.
func encryptedOVF(security string) string {
	return strings.Replace(`<?xml version="1.0"?>
<Envelope ovf:version="1.0" xmlns="http://schemas.dmtf.org/ovf/envelope/1" xmlns:ovf="http://schemas.dmtf.org/ovf/envelope/1"
  xmlns:rasd="http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_ResourceAllocationSettingData"
  xmlns:xenovf="http://schemas.citrix.com/ovf/envelope/1" xmlns:xenc="http://www.w3.org/2001/04/xmlenc#"
  xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
  xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
  <References><File ovf:id="xen_f1" ovf:href="disk.img"/><File ovf:id="xen_f2" ovf:href="plain.img"/>
    <File ovf:id="xen_f3" ovf:href="disk2.img.gz" ovf:compression="gzip"/></References>
  <SecuritySection>SECURITY</SecuritySection>
  <DiskSection><Disk ovf:diskId="d1" ovf:fileRef="xen_f1" ovf:capacity="4096"/><Disk ovf:diskId="d2" ovf:fileRef="xen_f2" ovf:capacity="4096"/>
    <Disk ovf:diskId="d3" ovf:fileRef="xen_f3" ovf:capacity="4096"/></DiskSection>
  <VirtualSystem ovf:id="tiny"><VirtualHardwareSection>
    <Item><rasd:ResourceType>17</rasd:ResourceType><rasd:InstanceID>9</rasd:InstanceID><rasd:HostResource>ovf:/disk/d1</rasd:HostResource></Item>
    <Item><rasd:ResourceType>17</rasd:ResourceType><rasd:InstanceID>10</rasd:InstanceID><rasd:HostResource>ovf:/disk/d2</rasd:HostResource></Item>
    <Item><rasd:ResourceType>17</rasd:ResourceType><rasd:InstanceID>11</rasd:InstanceID><rasd:HostResource>ovf:/disk/d3</rasd:HostResource></Item>
  </VirtualHardwareSection></VirtualSystem>
</Envelope>`, "SECURITY", security, 1)
}

func securityXML(version, check string) string {
	return `<Info>Encrypted Content Definition</Info><wsse:Security wsu:Id="xen_s" xenovf:version="` + version + `">
    <xenc:ReferenceList><xenc:DataReference ValueType="xen_f1"/><xenc:DataReference ValueType="xen_f3"/></xenc:ReferenceList>
    <xenc:EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#aes192-cbc"><xenc:KeySize>0</xenc:KeySize></xenc:EncryptionMethod>
    <xenc:EncryptedData><xenc:CipherData><xenc:CipherValue>` + check + `</xenc:CipherValue></xenc:CipherData></xenc:EncryptedData>
  </wsse:Security>`
}

func TestEncryptedPackage(t *testing.T) {
	const password = "s3cret pass"
	dir := t.TempDir()
	disk := bytes.Repeat([]byte("encrypted disk "), 4000)                                          // not a multiple of the block
	disk2, plain := bytes.Repeat([]byte("0123456789abcdef"), 512), bytes.Repeat([]byte("p"), 4096) // a multiple: a whole block of padding
	var gz bytes.Buffer
	z := gzip.NewWriter(&gz)
	z.Write(disk2)
	z.Close()
	write := func(name string, b []byte) {
		if err := os.WriteFile(filepath.Join(dir, name), b, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write("disk.img", csharpEncrypt(t, password, disk))
	write("plain.img", plain)
	write("disk2.img.gz", csharpEncrypt(t, password, gz.Bytes())) // compressed, then encrypted
	write("tiny.ovf", []byte(encryptedOVF(securityXML("1.3.1", csharpCheckValue(t, password)))))
	p, err := OpenPackage(filepath.Join(dir, "tiny.ovf"))
	if err != nil {
		t.Fatal(err)
	}
	a := p.Appliance
	if !a.Encrypted || !a.VMs[0].Disks[0].Encrypted || a.VMs[0].Disks[1].Encrypted || !a.VMs[0].Disks[2].Encrypted {
		t.Fatalf("%+v", a.VMs[0].Disks)
	}
	if _, err := p.OpenDisk(a.VMs[0].Disks[0]); err == nil || err.Error() != messages.Get("ERROR_NO_PASSWORD") {
		t.Errorf("no password: %v", err)
	}
	for _, pw := range []string{"", "wrong", "s3cret pas", "s3cret pass "} {
		if err := p.Unlock(pw); err == nil {
			t.Errorf("password %q accepted", pw)
		}
	}
	if err := p.Unlock(password); err != nil {
		t.Fatal(err)
	}
	for i, want := range [][]byte{disk, plain, disk2} {
		s, err := p.OpenDisk(a.VMs[0].Disks[i])
		if err != nil {
			t.Fatal(err)
		}
		got, err := io.ReadAll(s)
		s.Close()
		if err != nil || !bytes.Equal(got, want) || s.Length != int64(len(want)) {
			t.Errorf("disk %d: %v, %d bytes of %d, length %d", i, err, len(got), len(want), s.Length)
		}
	}
	// an encrypted stream-optimized VMDK, then one byte past its ciphertext: the file is read to its end
	grain := bytes.Repeat([]byte("g"), 128*sector)
	vmdk := streamVMDK(t, 2*int64(len(grain)), map[int64][]byte{1: grain}, []int64{1})
	for _, extra := range []string{"", "x"} {
		write("disk.img", append(csharpEncrypt(t, password, vmdk), extra...))
		s, err := p.OpenDisk(a.VMs[0].Disks[0])
		if err != nil {
			t.Fatal(err)
		}
		got, err := io.ReadAll(s)
		s.Close()
		if extra == "" && (err != nil || !bytes.Equal(got[len(grain):], grain)) || extra != "" && err == nil {
			t.Errorf("VMDK %q: %v", extra, err)
		}
	}
	// an extent the package lists as encrypted (disk.img) or compressed (disk2.img.gz), and one it does not (plain.img)
	for extent, ok := range map[string]bool{"disk.img": false, "disk2.img.gz": false, "plain.img": true} {
		x, err := p.openExtents("d.vmdk", []vmdkExtent{{file: extent, size: 4096}})
		if (err == nil) != ok {
			t.Errorf("extent %s: %v", extent, err)
		} else if x != nil {
			x.Close()
		}
	}
	// a disk encrypted with another password than the check value's
	write("disk.img", csharpEncrypt(t, "other", disk))
	if _, err := p.OpenDisk(a.VMs[0].Disks[0]); err == nil {
		t.Error("a disk of another password")
	}
	// the descriptors that cannot be decrypted
	check := csharpCheckValue(t, password)
	for name, sec := range map[string]string{
		"no version":     securityXML("", check),
		"bad version":    securityXML("1.x", check),
		"no check value": strings.ReplaceAll(securityXML("1.3.1", check), "xenc:EncryptedData", "xenc:Other"),
		"empty":          "",
		"two sections":   securityXML("1.3.1", check) + `</SecuritySection><SecuritySection>` + securityXML("1.3.1", check),
	} {
		write("bad.ovf", []byte(encryptedOVF(sec)))
		if _, err := OpenPackage(filepath.Join(dir, "bad.ovf")); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
	// CheckPassword's other spelling of EncryptedData, and a check value of another password
	for pw, sec := range map[string]string{
		password: strings.ReplaceAll(securityXML("1.3", check), "xenc:EncryptedData", "EncrypteData"),
		"x":      securityXML("1.3.1", csharpCheckValue(t, "y")),
	} {
		write("other.ovf", []byte(encryptedOVF(sec)))
		p, err := OpenPackage(filepath.Join(dir, "other.ovf"))
		if err != nil {
			t.Fatal(err)
		}
		if err := p.Unlock(pw); (err == nil) != (pw == password) {
			t.Errorf("%q: %v", pw, err)
		}
	}
}

func TestOVFDecrypter(t *testing.T) {
	data := bytes.Repeat([]byte("xcp-ng center "), 300)
	enc := csharpEncrypt(t, "pw", data)
	// byte by byte, and in reads of half a buffer
	for _, r := range []io.Reader{iotest.OneByteReader(bytes.NewReader(enc)), iotest.HalfReader(bytes.NewReader(enc))} {
		if got, err := io.ReadAll(newOVFDecrypter(r, "pw")); err != nil || !bytes.Equal(got, data) {
			t.Errorf("%v, %d bytes", err, len(got))
		}
	}
	bad := append([]byte{}, enc...)
	bad[len(bad)-1] ^= 0xff // the padding
	for name, c := range map[string][]byte{"truncated": enc[:len(enc)-3], "padding": bad, "empty": nil} {
		if _, err := io.ReadAll(newOVFDecrypter(bytes.NewReader(c), "pw")); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
	// Encoding.ASCII: a UTF-16 unit past ASCII is "?" (two for a character past the BMP)
	if k1, _ := ovfKey("pé😀"); !bytes.Equal(k1, func() []byte { k, _ := ovfKey("p???"); return k }()) {
		t.Error("ovfKey of non-ASCII characters")
	}
}
