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

type PVSProxyRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// PVS site this proxy is part of
	Site PVSSiteRef `json:"site,omitempty"`
	// VIF of the VM using the proxy
	VIF VIFRef `json:"VIF,omitempty"`
	// true = VM is currently proxied
	CurrentlyAttached bool `json:"currently_attached,omitempty"`
	// The run-time status of the proxy
	Status PvsProxyStatus `json:"status,omitempty"`
}

type PVSProxyRef string

// a proxy connects a VM/VIF with a PVS site
type pvsProxy struct{}

var PVSProxy pvsProxy

// GetAllRecords: Return a map of PVS_proxy references to PVS_proxy records for all PVS_proxys known to the system.
// Version: ely
func (pvsProxy) GetAllRecords(session *Session) (retval map[PVSProxyRef]PVSProxyRecord, err error) {
	method := "PVS_proxy.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSProxyRefToPVSProxyRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of PVS_proxy references to PVS_proxy records for all PVS_proxys known to the system.
// Version: ely
func (pvsProxy) GetAllRecords1(session *Session) (retval map[PVSProxyRef]PVSProxyRecord, err error) {
	method := "PVS_proxy.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSProxyRefToPVSProxyRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the PVS_proxys known to the system.
// Version: ely
func (pvsProxy) GetAll(session *Session) (retval []PVSProxyRef, err error) {
	method := "PVS_proxy.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSProxyRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the PVS_proxys known to the system.
// Version: ely
func (pvsProxy) GetAll1(session *Session) (retval []PVSProxyRef, err error) {
	method := "PVS_proxy.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSProxyRefSet(method+" -> ", result)
	return
}

// Destroy: remove (or switch off) a PVS proxy for this VM
// Version: ely
func (pvsProxy) Destroy(session *Session, self PVSProxyRef) (err error) {
	method := "PVS_proxy.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: remove (or switch off) a PVS proxy for this VM
// Version: ely
func (pvsProxy) AsyncDestroy(session *Session, self PVSProxyRef) (retval TaskRef, err error) {
	method := "Async.PVS_proxy.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: remove (or switch off) a PVS proxy for this VM
// Version: ely
func (pvsProxy) Destroy2(session *Session, self PVSProxyRef) (err error) {
	method := "PVS_proxy.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: remove (or switch off) a PVS proxy for this VM
// Version: ely
func (pvsProxy) AsyncDestroy2(session *Session, self PVSProxyRef) (retval TaskRef, err error) {
	method := "Async.PVS_proxy.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Configure a VM/VIF to use a PVS proxy
// Version: ely
func (pvsProxy) Create(session *Session, site PVSSiteRef, vif VIFRef) (retval PVSProxyRef, err error) {
	method := "PVS_proxy.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	siteArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "site"), site)
	if err != nil {
		return
	}
	vifArg, err := serializeVIFRef(fmt.Sprintf("%s(%s)", method, "VIF"), vif)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, siteArg, vifArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSProxyRef(method+" -> ", result)
	return
}

// AsyncCreate: Configure a VM/VIF to use a PVS proxy
// Version: ely
func (pvsProxy) AsyncCreate(session *Session, site PVSSiteRef, vif VIFRef) (retval TaskRef, err error) {
	method := "Async.PVS_proxy.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	siteArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "site"), site)
	if err != nil {
		return
	}
	vifArg, err := serializeVIFRef(fmt.Sprintf("%s(%s)", method, "VIF"), vif)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, siteArg, vifArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create3: Configure a VM/VIF to use a PVS proxy
// Version: ely
func (pvsProxy) Create3(session *Session, site PVSSiteRef, vif VIFRef) (retval PVSProxyRef, err error) {
	method := "PVS_proxy.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	siteArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "site"), site)
	if err != nil {
		return
	}
	vifArg, err := serializeVIFRef(fmt.Sprintf("%s(%s)", method, "VIF"), vif)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, siteArg, vifArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSProxyRef(method+" -> ", result)
	return
}

// AsyncCreate3: Configure a VM/VIF to use a PVS proxy
// Version: ely
func (pvsProxy) AsyncCreate3(session *Session, site PVSSiteRef, vif VIFRef) (retval TaskRef, err error) {
	method := "Async.PVS_proxy.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	siteArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "site"), site)
	if err != nil {
		return
	}
	vifArg, err := serializeVIFRef(fmt.Sprintf("%s(%s)", method, "VIF"), vif)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, siteArg, vifArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetStatus: Get the status field of the given PVS_proxy.
// Version: ely
func (pvsProxy) GetStatus(session *Session, self PVSProxyRef) (retval PvsProxyStatus, err error) {
	method := "PVS_proxy.get_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPvsProxyStatus(method+" -> ", result)
	return
}

// GetStatus2: Get the status field of the given PVS_proxy.
// Version: ely
func (pvsProxy) GetStatus2(session *Session, self PVSProxyRef) (retval PvsProxyStatus, err error) {
	method := "PVS_proxy.get_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPvsProxyStatus(method+" -> ", result)
	return
}

// GetCurrentlyAttached: Get the currently_attached field of the given PVS_proxy.
// Version: ely
func (pvsProxy) GetCurrentlyAttached(session *Session, self PVSProxyRef) (retval bool, err error) {
	method := "PVS_proxy.get_currently_attached"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetCurrentlyAttached2: Get the currently_attached field of the given PVS_proxy.
// Version: ely
func (pvsProxy) GetCurrentlyAttached2(session *Session, self PVSProxyRef) (retval bool, err error) {
	method := "PVS_proxy.get_currently_attached"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBool(method+" -> ", result)
	return
}

// GetVIF: Get the VIF field of the given PVS_proxy.
// Version: ely
func (pvsProxy) GetVIF(session *Session, self PVSProxyRef) (retval VIFRef, err error) {
	method := "PVS_proxy.get_VIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVIFRef(method+" -> ", result)
	return
}

// GetVIF2: Get the VIF field of the given PVS_proxy.
// Version: ely
func (pvsProxy) GetVIF2(session *Session, self PVSProxyRef) (retval VIFRef, err error) {
	method := "PVS_proxy.get_VIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVIFRef(method+" -> ", result)
	return
}

// GetSite: Get the site field of the given PVS_proxy.
// Version: ely
func (pvsProxy) GetSite(session *Session, self PVSProxyRef) (retval PVSSiteRef, err error) {
	method := "PVS_proxy.get_site"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSite2: Get the site field of the given PVS_proxy.
// Version: ely
func (pvsProxy) GetSite2(session *Session, self PVSProxyRef) (retval PVSSiteRef, err error) {
	method := "PVS_proxy.get_site"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given PVS_proxy.
// Version: ely
func (pvsProxy) GetUUID(session *Session, self PVSProxyRef) (retval string, err error) {
	method := "PVS_proxy.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given PVS_proxy.
// Version: ely
func (pvsProxy) GetUUID2(session *Session, self PVSProxyRef) (retval string, err error) {
	method := "PVS_proxy.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the PVS_proxy instance with the specified UUID.
// Version: ely
func (pvsProxy) GetByUUID(session *Session, uuid string) (retval PVSProxyRef, err error) {
	method := "PVS_proxy.get_by_uuid"
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
	retval, err = deserializePVSProxyRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the PVS_proxy instance with the specified UUID.
// Version: ely
func (pvsProxy) GetByUUID2(session *Session, uuid string) (retval PVSProxyRef, err error) {
	method := "PVS_proxy.get_by_uuid"
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
	retval, err = deserializePVSProxyRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given PVS_proxy.
// Version: ely
func (pvsProxy) GetRecord(session *Session, self PVSProxyRef) (retval PVSProxyRecord, err error) {
	method := "PVS_proxy.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSProxyRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given PVS_proxy.
// Version: ely
func (pvsProxy) GetRecord2(session *Session, self PVSProxyRef) (retval PVSProxyRecord, err error) {
	method := "PVS_proxy.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSProxyRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSProxyRecord(method+" -> ", result)
	return
}
