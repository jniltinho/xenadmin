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
)

type PVSSiteRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// a human-readable name
	NameLabel string `json:"name_label,omitempty"`
	// a notes field containing human-readable description
	NameDescription string `json:"name_description,omitempty"`
	// Unique identifier of the PVS site, as configured in PVS
	PVSUUID string `json:"PVS_uuid,omitempty"`
	// The SR used by PVS proxy for the cache
	CacheStorage []PVSCacheStorageRef `json:"cache_storage,omitempty"`
	// The set of PVS servers in the site
	Servers []PVSServerRef `json:"servers,omitempty"`
	// The set of proxies associated with the site
	Proxies []PVSProxyRef `json:"proxies,omitempty"`
}

type PVSSiteRef string

// machines serving blocks of data for provisioning VMs
type pvsSite struct{}

var PVSSite pvsSite

// GetAllRecords: Return a map of PVS_site references to PVS_site records for all PVS_sites known to the system.
// Version: ely
func (pvsSite) GetAllRecords(session *Session) (retval map[PVSSiteRef]PVSSiteRecord, err error) {
	method := "PVS_site.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSSiteRefToPVSSiteRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of PVS_site references to PVS_site records for all PVS_sites known to the system.
// Version: ely
func (pvsSite) GetAllRecords1(session *Session) (retval map[PVSSiteRef]PVSSiteRecord, err error) {
	method := "PVS_site.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSSiteRefToPVSSiteRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the PVS_sites known to the system.
// Version: ely
func (pvsSite) GetAll(session *Session) (retval []PVSSiteRef, err error) {
	method := "PVS_site.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSSiteRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the PVS_sites known to the system.
// Version: ely
func (pvsSite) GetAll1(session *Session) (retval []PVSSiteRef, err error) {
	method := "PVS_site.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSSiteRefSet(method+" -> ", result)
	return
}

// SetPVSUUID: Update the PVS UUID of the PVS site
// Version: ely
func (pvsSite) SetPVSUUID(session *Session, self PVSSiteRef, value string) (err error) {
	method := "PVS_site.set_PVS_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetPVSUUID: Update the PVS UUID of the PVS site
// Version: ely
func (pvsSite) AsyncSetPVSUUID(session *Session, self PVSSiteRef, value string) (retval TaskRef, err error) {
	method := "Async.PVS_site.set_PVS_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetPVSUUID3: Update the PVS UUID of the PVS site
// Version: ely
func (pvsSite) SetPVSUUID3(session *Session, self PVSSiteRef, value string) (err error) {
	method := "PVS_site.set_PVS_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetPVSUUID3: Update the PVS UUID of the PVS site
// Version: ely
func (pvsSite) AsyncSetPVSUUID3(session *Session, self PVSSiteRef, value string) (retval TaskRef, err error) {
	method := "Async.PVS_site.set_PVS_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Forget: Remove a site&apos;s meta data
// Version: ely
//
// Errors:
// PVS_SITE_CONTAINS_RUNNING_PROXIES - The PVS site contains running proxies.
// PVS_SITE_CONTAINS_SERVERS - The PVS site contains servers and cannot be forgotten.
func (pvsSite) Forget(session *Session, self PVSSiteRef) (err error) {
	method := "PVS_site.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncForget: Remove a site&apos;s meta data
// Version: ely
//
// Errors:
// PVS_SITE_CONTAINS_RUNNING_PROXIES - The PVS site contains running proxies.
// PVS_SITE_CONTAINS_SERVERS - The PVS site contains servers and cannot be forgotten.
func (pvsSite) AsyncForget(session *Session, self PVSSiteRef) (retval TaskRef, err error) {
	method := "Async.PVS_site.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Forget2: Remove a site&apos;s meta data
// Version: ely
//
// Errors:
// PVS_SITE_CONTAINS_RUNNING_PROXIES - The PVS site contains running proxies.
// PVS_SITE_CONTAINS_SERVERS - The PVS site contains servers and cannot be forgotten.
func (pvsSite) Forget2(session *Session, self PVSSiteRef) (err error) {
	method := "PVS_site.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncForget2: Remove a site&apos;s meta data
// Version: ely
//
// Errors:
// PVS_SITE_CONTAINS_RUNNING_PROXIES - The PVS site contains running proxies.
// PVS_SITE_CONTAINS_SERVERS - The PVS site contains servers and cannot be forgotten.
func (pvsSite) AsyncForget2(session *Session, self PVSSiteRef) (retval TaskRef, err error) {
	method := "Async.PVS_site.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce: Introduce new PVS site
// Version: ely
func (pvsSite) Introduce(session *Session, nameLabel string, nameDescription string, pvsUUID string) (retval PVSSiteRef, err error) {
	method := "PVS_site.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	pvsUUIDArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "PVS_uuid"), pvsUUID)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameLabelArg, nameDescriptionArg, pvsUUIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSSiteRef(method+" -> ", result)
	return
}

// AsyncIntroduce: Introduce new PVS site
// Version: ely
func (pvsSite) AsyncIntroduce(session *Session, nameLabel string, nameDescription string, pvsUUID string) (retval TaskRef, err error) {
	method := "Async.PVS_site.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	pvsUUIDArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "PVS_uuid"), pvsUUID)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameLabelArg, nameDescriptionArg, pvsUUIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce4: Introduce new PVS site
// Version: ely
func (pvsSite) Introduce4(session *Session, nameLabel string, nameDescription string, pvsUUID string) (retval PVSSiteRef, err error) {
	method := "PVS_site.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	pvsUUIDArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "PVS_uuid"), pvsUUID)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameLabelArg, nameDescriptionArg, pvsUUIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSSiteRef(method+" -> ", result)
	return
}

// AsyncIntroduce4: Introduce new PVS site
// Version: ely
func (pvsSite) AsyncIntroduce4(session *Session, nameLabel string, nameDescription string, pvsUUID string) (retval TaskRef, err error) {
	method := "Async.PVS_site.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	pvsUUIDArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "PVS_uuid"), pvsUUID)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, nameLabelArg, nameDescriptionArg, pvsUUIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetNameDescription: Set the name/description field of the given PVS_site.
// Version: ely
func (pvsSite) SetNameDescription(session *Session, self PVSSiteRef, value string) (err error) {
	method := "PVS_site.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetNameDescription3: Set the name/description field of the given PVS_site.
// Version: ely
func (pvsSite) SetNameDescription3(session *Session, self PVSSiteRef, value string) (err error) {
	method := "PVS_site.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetNameLabel: Set the name/label field of the given PVS_site.
// Version: ely
func (pvsSite) SetNameLabel(session *Session, self PVSSiteRef, value string) (err error) {
	method := "PVS_site.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetNameLabel3: Set the name/label field of the given PVS_site.
// Version: ely
func (pvsSite) SetNameLabel3(session *Session, self PVSSiteRef, value string) (err error) {
	method := "PVS_site.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// GetProxies: Get the proxies field of the given PVS_site.
// Version: ely
func (pvsSite) GetProxies(session *Session, self PVSSiteRef) (retval []PVSProxyRef, err error) {
	method := "PVS_site.get_proxies"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSProxyRefSet(method+" -> ", result)
	return
}

// GetProxies2: Get the proxies field of the given PVS_site.
// Version: ely
func (pvsSite) GetProxies2(session *Session, self PVSSiteRef) (retval []PVSProxyRef, err error) {
	method := "PVS_site.get_proxies"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSProxyRefSet(method+" -> ", result)
	return
}

// GetServers: Get the servers field of the given PVS_site.
// Version: ely
func (pvsSite) GetServers(session *Session, self PVSSiteRef) (retval []PVSServerRef, err error) {
	method := "PVS_site.get_servers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSServerRefSet(method+" -> ", result)
	return
}

// GetServers2: Get the servers field of the given PVS_site.
// Version: ely
func (pvsSite) GetServers2(session *Session, self PVSSiteRef) (retval []PVSServerRef, err error) {
	method := "PVS_site.get_servers"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSServerRefSet(method+" -> ", result)
	return
}

// GetCacheStorage: Get the cache_storage field of the given PVS_site.
// Version: ely
func (pvsSite) GetCacheStorage(session *Session, self PVSSiteRef) (retval []PVSCacheStorageRef, err error) {
	method := "PVS_site.get_cache_storage"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSCacheStorageRefSet(method+" -> ", result)
	return
}

// GetCacheStorage2: Get the cache_storage field of the given PVS_site.
// Version: ely
func (pvsSite) GetCacheStorage2(session *Session, self PVSSiteRef) (retval []PVSCacheStorageRef, err error) {
	method := "PVS_site.get_cache_storage"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSCacheStorageRefSet(method+" -> ", result)
	return
}

// GetPVSUUID: Get the PVS_uuid field of the given PVS_site.
// Version: ely
func (pvsSite) GetPVSUUID(session *Session, self PVSSiteRef) (retval string, err error) {
	method := "PVS_site.get_PVS_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVSUUID2: Get the PVS_uuid field of the given PVS_site.
// Version: ely
func (pvsSite) GetPVSUUID2(session *Session, self PVSSiteRef) (retval string, err error) {
	method := "PVS_site.get_PVS_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription: Get the name/description field of the given PVS_site.
// Version: ely
func (pvsSite) GetNameDescription(session *Session, self PVSSiteRef) (retval string, err error) {
	method := "PVS_site.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription2: Get the name/description field of the given PVS_site.
// Version: ely
func (pvsSite) GetNameDescription2(session *Session, self PVSSiteRef) (retval string, err error) {
	method := "PVS_site.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel: Get the name/label field of the given PVS_site.
// Version: ely
func (pvsSite) GetNameLabel(session *Session, self PVSSiteRef) (retval string, err error) {
	method := "PVS_site.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel2: Get the name/label field of the given PVS_site.
// Version: ely
func (pvsSite) GetNameLabel2(session *Session, self PVSSiteRef) (retval string, err error) {
	method := "PVS_site.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given PVS_site.
// Version: ely
func (pvsSite) GetUUID(session *Session, self PVSSiteRef) (retval string, err error) {
	method := "PVS_site.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given PVS_site.
// Version: ely
func (pvsSite) GetUUID2(session *Session, self PVSSiteRef) (retval string, err error) {
	method := "PVS_site.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByNameLabel: Get all the PVS_site instances with the given label.
// Version: ely
func (pvsSite) GetByNameLabel(session *Session, label string) (retval []PVSSiteRef, err error) {
	method := "PVS_site.get_by_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	labelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "label"), label)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, labelArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSSiteRefSet(method+" -> ", result)
	return
}

// GetByNameLabel2: Get all the PVS_site instances with the given label.
// Version: ely
func (pvsSite) GetByNameLabel2(session *Session, label string) (retval []PVSSiteRef, err error) {
	method := "PVS_site.get_by_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	labelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "label"), label)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, labelArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSSiteRefSet(method+" -> ", result)
	return
}

// GetByUUID: Get a reference to the PVS_site instance with the specified UUID.
// Version: ely
func (pvsSite) GetByUUID(session *Session, uuid string) (retval PVSSiteRef, err error) {
	method := "PVS_site.get_by_uuid"
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
	retval, err = deserializePVSSiteRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the PVS_site instance with the specified UUID.
// Version: ely
func (pvsSite) GetByUUID2(session *Session, uuid string) (retval PVSSiteRef, err error) {
	method := "PVS_site.get_by_uuid"
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
	retval, err = deserializePVSSiteRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given PVS_site.
// Version: ely
func (pvsSite) GetRecord(session *Session, self PVSSiteRef) (retval PVSSiteRecord, err error) {
	method := "PVS_site.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSSiteRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given PVS_site.
// Version: ely
func (pvsSite) GetRecord2(session *Session, self PVSSiteRef) (retval PVSSiteRecord, err error) {
	method := "PVS_site.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSSiteRecord(method+" -> ", result)
	return
}
