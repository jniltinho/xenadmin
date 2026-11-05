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

type PVSServerRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// IPv4 addresses of this server
	Addresses []string `json:"addresses,omitempty"`
	// First UDP port accepted by this server
	FirstPort int `json:"first_port,omitempty"`
	// Last UDP port accepted by this server
	LastPort int `json:"last_port,omitempty"`
	// PVS site this server is part of
	Site PVSSiteRef `json:"site,omitempty"`
}

type PVSServerRef string

// individual machine serving provisioning (block) data
type pvsServer struct{}

var PVSServer pvsServer

// GetAllRecords: Return a map of PVS_server references to PVS_server records for all PVS_servers known to the system.
// Version: ely
func (pvsServer) GetAllRecords(session *Session) (retval map[PVSServerRef]PVSServerRecord, err error) {
	method := "PVS_server.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSServerRefToPVSServerRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of PVS_server references to PVS_server records for all PVS_servers known to the system.
// Version: ely
func (pvsServer) GetAllRecords1(session *Session) (retval map[PVSServerRef]PVSServerRecord, err error) {
	method := "PVS_server.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSServerRefToPVSServerRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the PVS_servers known to the system.
// Version: ely
func (pvsServer) GetAll(session *Session) (retval []PVSServerRef, err error) {
	method := "PVS_server.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSServerRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the PVS_servers known to the system.
// Version: ely
func (pvsServer) GetAll1(session *Session) (retval []PVSServerRef, err error) {
	method := "PVS_server.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSServerRefSet(method+" -> ", result)
	return
}

// Forget: forget a PVS server
// Version: ely
func (pvsServer) Forget(session *Session, self PVSServerRef) (err error) {
	method := "PVS_server.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncForget: forget a PVS server
// Version: ely
func (pvsServer) AsyncForget(session *Session, self PVSServerRef) (retval TaskRef, err error) {
	method := "Async.PVS_server.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Forget2: forget a PVS server
// Version: ely
func (pvsServer) Forget2(session *Session, self PVSServerRef) (err error) {
	method := "PVS_server.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncForget2: forget a PVS server
// Version: ely
func (pvsServer) AsyncForget2(session *Session, self PVSServerRef) (retval TaskRef, err error) {
	method := "Async.PVS_server.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Introduce: introduce new PVS server
// Version: ely
func (pvsServer) Introduce(session *Session, addresses []string, firstPort int, lastPort int, site PVSSiteRef) (retval PVSServerRef, err error) {
	method := "PVS_server.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	addressesArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "addresses"), addresses)
	if err != nil {
		return
	}
	firstPortArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "first_port"), firstPort)
	if err != nil {
		return
	}
	lastPortArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "last_port"), lastPort)
	if err != nil {
		return
	}
	siteArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "site"), site)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, addressesArg, firstPortArg, lastPortArg, siteArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSServerRef(method+" -> ", result)
	return
}

// AsyncIntroduce: introduce new PVS server
// Version: ely
func (pvsServer) AsyncIntroduce(session *Session, addresses []string, firstPort int, lastPort int, site PVSSiteRef) (retval TaskRef, err error) {
	method := "Async.PVS_server.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	addressesArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "addresses"), addresses)
	if err != nil {
		return
	}
	firstPortArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "first_port"), firstPort)
	if err != nil {
		return
	}
	lastPortArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "last_port"), lastPort)
	if err != nil {
		return
	}
	siteArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "site"), site)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, addressesArg, firstPortArg, lastPortArg, siteArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce5: introduce new PVS server
// Version: ely
func (pvsServer) Introduce5(session *Session, addresses []string, firstPort int, lastPort int, site PVSSiteRef) (retval PVSServerRef, err error) {
	method := "PVS_server.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	addressesArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "addresses"), addresses)
	if err != nil {
		return
	}
	firstPortArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "first_port"), firstPort)
	if err != nil {
		return
	}
	lastPortArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "last_port"), lastPort)
	if err != nil {
		return
	}
	siteArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "site"), site)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, addressesArg, firstPortArg, lastPortArg, siteArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSServerRef(method+" -> ", result)
	return
}

// AsyncIntroduce5: introduce new PVS server
// Version: ely
func (pvsServer) AsyncIntroduce5(session *Session, addresses []string, firstPort int, lastPort int, site PVSSiteRef) (retval TaskRef, err error) {
	method := "Async.PVS_server.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	addressesArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "addresses"), addresses)
	if err != nil {
		return
	}
	firstPortArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "first_port"), firstPort)
	if err != nil {
		return
	}
	lastPortArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "last_port"), lastPort)
	if err != nil {
		return
	}
	siteArg, err := serializePVSSiteRef(fmt.Sprintf("%s(%s)", method, "site"), site)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, addressesArg, firstPortArg, lastPortArg, siteArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetSite: Get the site field of the given PVS_server.
// Version: ely
func (pvsServer) GetSite(session *Session, self PVSServerRef) (retval PVSSiteRef, err error) {
	method := "PVS_server.get_site"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSite2: Get the site field of the given PVS_server.
// Version: ely
func (pvsServer) GetSite2(session *Session, self PVSServerRef) (retval PVSSiteRef, err error) {
	method := "PVS_server.get_site"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastPort: Get the last_port field of the given PVS_server.
// Version: ely
func (pvsServer) GetLastPort(session *Session, self PVSServerRef) (retval int, err error) {
	method := "PVS_server.get_last_port"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastPort2: Get the last_port field of the given PVS_server.
// Version: ely
func (pvsServer) GetLastPort2(session *Session, self PVSServerRef) (retval int, err error) {
	method := "PVS_server.get_last_port"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetFirstPort: Get the first_port field of the given PVS_server.
// Version: ely
func (pvsServer) GetFirstPort(session *Session, self PVSServerRef) (retval int, err error) {
	method := "PVS_server.get_first_port"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetFirstPort2: Get the first_port field of the given PVS_server.
// Version: ely
func (pvsServer) GetFirstPort2(session *Session, self PVSServerRef) (retval int, err error) {
	method := "PVS_server.get_first_port"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAddresses: Get the addresses field of the given PVS_server.
// Version: ely
func (pvsServer) GetAddresses(session *Session, self PVSServerRef) (retval []string, err error) {
	method := "PVS_server.get_addresses"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetAddresses2: Get the addresses field of the given PVS_server.
// Version: ely
func (pvsServer) GetAddresses2(session *Session, self PVSServerRef) (retval []string, err error) {
	method := "PVS_server.get_addresses"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetUUID: Get the uuid field of the given PVS_server.
// Version: ely
func (pvsServer) GetUUID(session *Session, self PVSServerRef) (retval string, err error) {
	method := "PVS_server.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given PVS_server.
// Version: ely
func (pvsServer) GetUUID2(session *Session, self PVSServerRef) (retval string, err error) {
	method := "PVS_server.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the PVS_server instance with the specified UUID.
// Version: ely
func (pvsServer) GetByUUID(session *Session, uuid string) (retval PVSServerRef, err error) {
	method := "PVS_server.get_by_uuid"
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
	retval, err = deserializePVSServerRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the PVS_server instance with the specified UUID.
// Version: ely
func (pvsServer) GetByUUID2(session *Session, uuid string) (retval PVSServerRef, err error) {
	method := "PVS_server.get_by_uuid"
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
	retval, err = deserializePVSServerRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given PVS_server.
// Version: ely
func (pvsServer) GetRecord(session *Session, self PVSServerRef) (retval PVSServerRecord, err error) {
	method := "PVS_server.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSServerRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given PVS_server.
// Version: ely
func (pvsServer) GetRecord2(session *Session, self PVSServerRef) (retval PVSServerRecord, err error) {
	method := "PVS_server.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePVSServerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePVSServerRecord(method+" -> ", result)
	return
}
