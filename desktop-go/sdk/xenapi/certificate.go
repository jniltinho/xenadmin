/*
 * Copyright (c) Cloud Software Group, Inc.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 *   1) Redistributions of source code must retain the above copyright
 *      notice, this list of conditions and the following disclaimer.
 *
 *   2) Redistributions in binary form must reproduce the above
 *      copyright notice, this list of conditions and the following
 *      disclaimer in the documentation and/or other materials
 *      provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
 * COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
 * OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package xenapi

import (
	"fmt"
	"time"
)

type CertificateRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// The name of the certificate, only present on certificates of type &apos;ca&apos;
	Name string `json:"name,omitempty"`
	// The type of the certificate, either &apos;ca&apos;, &apos;host&apos; or &apos;host_internal&apos;
	Type CertificateType `json:"type,omitempty"`
	// The host where the certificate is installed
	Host HostRef `json:"host,omitempty"`
	// Date after which the certificate is valid
	NotBefore time.Time `json:"not_before,omitempty"`
	// Date before which the certificate is valid
	NotAfter time.Time `json:"not_after,omitempty"`
	// Use fingerprint_sha256 instead
	Fingerprint string `json:"fingerprint,omitempty"`
	// The certificate&apos;s SHA256 fingerprint / hash
	FingerprintSha256 string `json:"fingerprint_sha256,omitempty"`
	// The certificate&apos;s SHA1 fingerprint / hash
	FingerprintSha1 string `json:"fingerprint_sha1,omitempty"`
	// The purposes of the certificate
	Purpose []CertificatePurpose `json:"purpose,omitempty"`
}

type CertificateRef string

// An X509 certificate used for TLS connections
type certificate struct{}

var Certificate certificate

// GetAllRecords: Return a map of Certificate references to Certificate records for all Certificates known to the system.
// Version: stockholm
func (certificate) GetAllRecords(session *Session) (retval map[CertificateRef]CertificateRecord, err error) {
	method := "Certificate.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeCertificateRefToCertificateRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of Certificate references to Certificate records for all Certificates known to the system.
// Version: stockholm
func (certificate) GetAllRecords1(session *Session) (retval map[CertificateRef]CertificateRecord, err error) {
	method := "Certificate.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeCertificateRefToCertificateRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the Certificates known to the system.
// Version: stockholm
func (certificate) GetAll(session *Session) (retval []CertificateRef, err error) {
	method := "Certificate.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeCertificateRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the Certificates known to the system.
// Version: stockholm
func (certificate) GetAll1(session *Session) (retval []CertificateRef, err error) {
	method := "Certificate.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeCertificateRefSet(method+" -> ", result)
	return
}

// GetPurpose: Get the purpose field of the given Certificate.
// Version: stockholm
func (certificate) GetPurpose(session *Session, self CertificateRef) (retval []CertificatePurpose, err error) {
	method := "Certificate.get_purpose"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumCertificatePurposeSet(method+" -> ", result)
	return
}

// GetPurpose2: Get the purpose field of the given Certificate.
// Version: stockholm
func (certificate) GetPurpose2(session *Session, self CertificateRef) (retval []CertificatePurpose, err error) {
	method := "Certificate.get_purpose"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumCertificatePurposeSet(method+" -> ", result)
	return
}

// GetFingerprintSha1: Get the fingerprint_sha1 field of the given Certificate.
// Version: stockholm
func (certificate) GetFingerprintSha1(session *Session, self CertificateRef) (retval string, err error) {
	method := "Certificate.get_fingerprint_sha1"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetFingerprintSha12: Get the fingerprint_sha1 field of the given Certificate.
// Version: stockholm
func (certificate) GetFingerprintSha12(session *Session, self CertificateRef) (retval string, err error) {
	method := "Certificate.get_fingerprint_sha1"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetFingerprintSha256: Get the fingerprint_sha256 field of the given Certificate.
// Version: stockholm
func (certificate) GetFingerprintSha256(session *Session, self CertificateRef) (retval string, err error) {
	method := "Certificate.get_fingerprint_sha256"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetFingerprintSha2562: Get the fingerprint_sha256 field of the given Certificate.
// Version: stockholm
func (certificate) GetFingerprintSha2562(session *Session, self CertificateRef) (retval string, err error) {
	method := "Certificate.get_fingerprint_sha256"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetFingerprint: Get the fingerprint field of the given Certificate.
// Version: stockholm
func (certificate) GetFingerprint(session *Session, self CertificateRef) (retval string, err error) {
	method := "Certificate.get_fingerprint"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetFingerprint2: Get the fingerprint field of the given Certificate.
// Version: stockholm
func (certificate) GetFingerprint2(session *Session, self CertificateRef) (retval string, err error) {
	method := "Certificate.get_fingerprint"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetNotAfter: Get the not_after field of the given Certificate.
// Version: stockholm
func (certificate) GetNotAfter(session *Session, self CertificateRef) (retval time.Time, err error) {
	method := "Certificate.get_not_after"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetNotAfter2: Get the not_after field of the given Certificate.
// Version: stockholm
func (certificate) GetNotAfter2(session *Session, self CertificateRef) (retval time.Time, err error) {
	method := "Certificate.get_not_after"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetNotBefore: Get the not_before field of the given Certificate.
// Version: stockholm
func (certificate) GetNotBefore(session *Session, self CertificateRef) (retval time.Time, err error) {
	method := "Certificate.get_not_before"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetNotBefore2: Get the not_before field of the given Certificate.
// Version: stockholm
func (certificate) GetNotBefore2(session *Session, self CertificateRef) (retval time.Time, err error) {
	method := "Certificate.get_not_before"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetHost: Get the host field of the given Certificate.
// Version: stockholm
func (certificate) GetHost(session *Session, self CertificateRef) (retval HostRef, err error) {
	method := "Certificate.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRef(method+" -> ", result)
	return
}

// GetHost2: Get the host field of the given Certificate.
// Version: stockholm
func (certificate) GetHost2(session *Session, self CertificateRef) (retval HostRef, err error) {
	method := "Certificate.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRef(method+" -> ", result)
	return
}

// GetType: Get the type field of the given Certificate.
// Version: stockholm
func (certificate) GetType(session *Session, self CertificateRef) (retval CertificateType, err error) {
	method := "Certificate.get_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumCertificateType(method+" -> ", result)
	return
}

// GetType2: Get the type field of the given Certificate.
// Version: stockholm
func (certificate) GetType2(session *Session, self CertificateRef) (retval CertificateType, err error) {
	method := "Certificate.get_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumCertificateType(method+" -> ", result)
	return
}

// GetName: Get the name field of the given Certificate.
// Version: stockholm
func (certificate) GetName(session *Session, self CertificateRef) (retval string, err error) {
	method := "Certificate.get_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetName2: Get the name field of the given Certificate.
// Version: stockholm
func (certificate) GetName2(session *Session, self CertificateRef) (retval string, err error) {
	method := "Certificate.get_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetUUID: Get the uuid field of the given Certificate.
// Version: stockholm
func (certificate) GetUUID(session *Session, self CertificateRef) (retval string, err error) {
	method := "Certificate.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetUUID2: Get the uuid field of the given Certificate.
// Version: stockholm
func (certificate) GetUUID2(session *Session, self CertificateRef) (retval string, err error) {
	method := "Certificate.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// GetByUUID: Get a reference to the Certificate instance with the specified UUID.
// Version: stockholm
func (certificate) GetByUUID(session *Session, uuid string) (retval CertificateRef, err error) {
	method := "Certificate.get_by_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg)
	if err != nil {
		return
	}
	retval, err = deserializeCertificateRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the Certificate instance with the specified UUID.
// Version: stockholm
func (certificate) GetByUUID2(session *Session, uuid string) (retval CertificateRef, err error) {
	method := "Certificate.get_by_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg)
	if err != nil {
		return
	}
	retval, err = deserializeCertificateRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given Certificate.
// Version: stockholm
func (certificate) GetRecord(session *Session, self CertificateRef) (retval CertificateRecord, err error) {
	method := "Certificate.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeCertificateRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given Certificate.
// Version: stockholm
func (certificate) GetRecord2(session *Session, self CertificateRef) (retval CertificateRecord, err error) {
	method := "Certificate.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeCertificateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeCertificateRecord(method+" -> ", result)
	return
}
