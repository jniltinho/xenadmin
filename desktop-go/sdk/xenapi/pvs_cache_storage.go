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

type PVSCacheStorageRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// The host on which this object defines PVS cache storage
	Host HostRef `json:"host,omitempty"`
	// SR providing storage for the PVS cache
	SR SRRef `json:"SR,omitempty"`
	// The PVS_site for which this object defines the storage
	Site PVSSiteRef `json:"site,omitempty"`
	// The size of the cache VDI (in bytes)
	Size int `json:"size,omitempty"`
	// The VDI used for caching
	VDI VDIRef `json:"VDI,omitempty"`
}

type PVSCacheStorageRef string

// Describes the storage that is available to a PVS site for caching purposes
type pvsCacheStorage struct{}

var PVSCacheStorage pvsCacheStorage

// GetAllRecords: Return a map of PVS_cache_storage references to PVS_cache_storage records for all PVS_cache_storages known to the system.
// Version: ely
func (pvsCacheStorage) GetAllRecords(session *Session) (retval map[PVSCacheStorageRef]PVSCacheStorageRecord, err error) {
	method := "PVS_cache_storage.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSCacheStorageRefToPVSCacheStorageRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of PVS_cache_storage references to PVS_cache_storage records for all PVS_cache_storages known to the system.
// Version: ely
func (pvsCacheStorage) GetAllRecords1(session *Session) (retval map[PVSCacheStorageRef]PVSCacheStorageRecord, err error) {
	method := "PVS_cache_storage.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSCacheStorageRefToPVSCacheStorageRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the PVS_cache_storages known to the system.
// Version: ely
func (pvsCacheStorage) GetAll(session *Session) (retval []PVSCacheStorageRef, err error) {
	method := "PVS_cache_storage.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSCacheStorageRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the PVS_cache_storages known to the system.
// Version: ely
func (pvsCacheStorage) GetAll1(session *Session) (retval []PVSCacheStorageRef, err error) {
	method := "PVS_cache_storage.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSCacheStorageRefSet(method+" -> ", result)
	return
}

// GetVDI: Get the VDI field of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetVDI(session *Session, self PVSCacheStorageRef) (retval VDIRef, err error) {
	method := "PVS_cache_storage.get_VDI"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// GetVDI2: Get the VDI field of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetVDI2(session *Session, self PVSCacheStorageRef) (retval VDIRef, err error) {
	method := "PVS_cache_storage.get_VDI"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// GetSize: Get the size field of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetSize(session *Session, self PVSCacheStorageRef) (retval int, err error) {
	method := "PVS_cache_storage.get_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// GetSize2: Get the size field of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetSize2(session *Session, self PVSCacheStorageRef) (retval int, err error) {
	method := "PVS_cache_storage.get_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// GetSite: Get the site field of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetSite(session *Session, self PVSCacheStorageRef) (retval PVSSiteRef, err error) {
	method := "PVS_cache_storage.get_site"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSSiteRef(method+" -> ", result)
	return
}

// GetSite2: Get the site field of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetSite2(session *Session, self PVSCacheStorageRef) (retval PVSSiteRef, err error) {
	method := "PVS_cache_storage.get_site"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSSiteRef(method+" -> ", result)
	return
}

// GetSR: Get the SR field of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetSR(session *Session, self PVSCacheStorageRef) (retval SRRef, err error) {
	method := "PVS_cache_storage.get_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// GetSR2: Get the SR field of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetSR2(session *Session, self PVSCacheStorageRef) (retval SRRef, err error) {
	method := "PVS_cache_storage.get_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// GetHost: Get the host field of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetHost(session *Session, self PVSCacheStorageRef) (retval HostRef, err error) {
	method := "PVS_cache_storage.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHost2: Get the host field of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetHost2(session *Session, self PVSCacheStorageRef) (retval HostRef, err error) {
	method := "PVS_cache_storage.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetUUID(session *Session, self PVSCacheStorageRef) (retval string, err error) {
	method := "PVS_cache_storage.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetUUID2(session *Session, self PVSCacheStorageRef) (retval string, err error) {
	method := "PVS_cache_storage.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy: Destroy the specified PVS_cache_storage instance.
// Version: ely
func (pvsCacheStorage) Destroy(session *Session, self PVSCacheStorageRef) (err error) {
	method := "PVS_cache_storage.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy the specified PVS_cache_storage instance.
// Version: ely
func (pvsCacheStorage) AsyncDestroy(session *Session, self PVSCacheStorageRef) (retval TaskRef, err error) {
	method := "Async.PVS_cache_storage.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy the specified PVS_cache_storage instance.
// Version: ely
func (pvsCacheStorage) Destroy2(session *Session, self PVSCacheStorageRef) (err error) {
	method := "PVS_cache_storage.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy the specified PVS_cache_storage instance.
// Version: ely
func (pvsCacheStorage) AsyncDestroy2(session *Session, self PVSCacheStorageRef) (retval TaskRef, err error) {
	method := "Async.PVS_cache_storage.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Create a new PVS_cache_storage instance, and return its handle. The constructor args are: host, SR, site, size (* = non-optional).
// Version: ely
func (pvsCacheStorage) Create(session *Session, args PVSCacheStorageRecord) (retval PVSCacheStorageRef, err error) {
	method := "PVS_cache_storage.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializePVSCacheStorageRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSCacheStorageRef(method+" -> ", result)
	return
}

// AsyncCreate: Create a new PVS_cache_storage instance, and return its handle. The constructor args are: host, SR, site, size (* = non-optional).
// Version: ely
func (pvsCacheStorage) AsyncCreate(session *Session, args PVSCacheStorageRecord) (retval TaskRef, err error) {
	method := "Async.PVS_cache_storage.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializePVSCacheStorageRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create2: Create a new PVS_cache_storage instance, and return its handle. The constructor args are: host, SR, site, size (* = non-optional).
// Version: ely
func (pvsCacheStorage) Create2(session *Session, args PVSCacheStorageRecord) (retval PVSCacheStorageRef, err error) {
	method := "PVS_cache_storage.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializePVSCacheStorageRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSCacheStorageRef(method+" -> ", result)
	return
}

// AsyncCreate2: Create a new PVS_cache_storage instance, and return its handle. The constructor args are: host, SR, site, size (* = non-optional).
// Version: ely
func (pvsCacheStorage) AsyncCreate2(session *Session, args PVSCacheStorageRecord) (retval TaskRef, err error) {
	method := "Async.PVS_cache_storage.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializePVSCacheStorageRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetByUUID: Get a reference to the PVS_cache_storage instance with the specified UUID.
// Version: ely
func (pvsCacheStorage) GetByUUID(session *Session, uuid string) (retval PVSCacheStorageRef, err error) {
	method := "PVS_cache_storage.get_by_uuid"
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
	retval, err = deserializePVSCacheStorageRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the PVS_cache_storage instance with the specified UUID.
// Version: ely
func (pvsCacheStorage) GetByUUID2(session *Session, uuid string) (retval PVSCacheStorageRef, err error) {
	method := "PVS_cache_storage.get_by_uuid"
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
	retval, err = deserializePVSCacheStorageRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetRecord(session *Session, self PVSCacheStorageRef) (retval PVSCacheStorageRecord, err error) {
	method := "PVS_cache_storage.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSCacheStorageRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given PVS_cache_storage.
// Version: ely
func (pvsCacheStorage) GetRecord2(session *Session, self PVSCacheStorageRef) (retval PVSCacheStorageRecord, err error) {
	method := "PVS_cache_storage.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSCacheStorageRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSCacheStorageRecord(method+" -> ", result)
	return
}
