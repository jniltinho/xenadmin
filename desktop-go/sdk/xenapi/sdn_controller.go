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

type SDNControllerRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// Protocol to connect with SDN controller
	Protocol SdnControllerProtocol `json:"protocol,omitempty"`
	// IP address of the controller
	Address string `json:"address,omitempty"`
	// TCP port of the controller
	Port int `json:"port,omitempty"`
}

type SDNControllerRef string

// Describes the SDN controller that is to connect with the pool
type sdnController struct{}

var SDNController sdnController

// GetAllRecords: Return a map of SDN_controller references to SDN_controller records for all SDN_controllers known to the system.
// Version: falcon
func (sdnController) GetAllRecords(session *Session) (retval map[SDNControllerRef]SDNControllerRecord, err error) {
	method := "SDN_controller.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeSDNControllerRefToSDNControllerRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of SDN_controller references to SDN_controller records for all SDN_controllers known to the system.
// Version: falcon
func (sdnController) GetAllRecords1(session *Session) (retval map[SDNControllerRef]SDNControllerRecord, err error) {
	method := "SDN_controller.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeSDNControllerRefToSDNControllerRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the SDN_controllers known to the system.
// Version: falcon
func (sdnController) GetAll(session *Session) (retval []SDNControllerRef, err error) {
	method := "SDN_controller.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeSDNControllerRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the SDN_controllers known to the system.
// Version: falcon
func (sdnController) GetAll1(session *Session) (retval []SDNControllerRef, err error) {
	method := "SDN_controller.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeSDNControllerRefSet(method+" -> ", result)
	return
}

// Forget: Remove the OVS manager of the pool and destroy the db record.
// Version: falcon
func (sdnController) Forget(session *Session, self SDNControllerRef) (err error) {
	method := "SDN_controller.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncForget: Remove the OVS manager of the pool and destroy the db record.
// Version: falcon
func (sdnController) AsyncForget(session *Session, self SDNControllerRef) (retval TaskRef, err error) {
	method := "Async.SDN_controller.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Forget2: Remove the OVS manager of the pool and destroy the db record.
// Version: falcon
func (sdnController) Forget2(session *Session, self SDNControllerRef) (err error) {
	method := "SDN_controller.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncForget2: Remove the OVS manager of the pool and destroy the db record.
// Version: falcon
func (sdnController) AsyncForget2(session *Session, self SDNControllerRef) (retval TaskRef, err error) {
	method := "Async.SDN_controller.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Introduce: Introduce an SDN controller to the pool.
// Version: falcon
func (sdnController) Introduce(session *Session, protocol SdnControllerProtocol, address string, port int) (retval SDNControllerRef, err error) {
	method := "SDN_controller.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	protocolArg, err := serializeEnumSdnControllerProtocol(fmt.Sprintf("%s(%s)", method, "protocol"), protocol)
	if err != nil {
		return
	}
	addressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "address"), address)
	if err != nil {
		return
	}
	portArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "port"), port)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, protocolArg, addressArg, portArg)
	if err != nil {
		return
	}
	retval, err = deserializeSDNControllerRef(method+" -> ", result)
	return
}

// AsyncIntroduce: Introduce an SDN controller to the pool.
// Version: falcon
func (sdnController) AsyncIntroduce(session *Session, protocol SdnControllerProtocol, address string, port int) (retval TaskRef, err error) {
	method := "Async.SDN_controller.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	protocolArg, err := serializeEnumSdnControllerProtocol(fmt.Sprintf("%s(%s)", method, "protocol"), protocol)
	if err != nil {
		return
	}
	addressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "address"), address)
	if err != nil {
		return
	}
	portArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "port"), port)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, protocolArg, addressArg, portArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce4: Introduce an SDN controller to the pool.
// Version: falcon
func (sdnController) Introduce4(session *Session, protocol SdnControllerProtocol, address string, port int) (retval SDNControllerRef, err error) {
	method := "SDN_controller.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	protocolArg, err := serializeEnumSdnControllerProtocol(fmt.Sprintf("%s(%s)", method, "protocol"), protocol)
	if err != nil {
		return
	}
	addressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "address"), address)
	if err != nil {
		return
	}
	portArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "port"), port)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, protocolArg, addressArg, portArg)
	if err != nil {
		return
	}
	retval, err = deserializeSDNControllerRef(method+" -> ", result)
	return
}

// AsyncIntroduce4: Introduce an SDN controller to the pool.
// Version: falcon
func (sdnController) AsyncIntroduce4(session *Session, protocol SdnControllerProtocol, address string, port int) (retval TaskRef, err error) {
	method := "Async.SDN_controller.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	protocolArg, err := serializeEnumSdnControllerProtocol(fmt.Sprintf("%s(%s)", method, "protocol"), protocol)
	if err != nil {
		return
	}
	addressArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "address"), address)
	if err != nil {
		return
	}
	portArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "port"), port)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, protocolArg, addressArg, portArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetPort: Get the port field of the given SDN_controller.
// Version: falcon
func (sdnController) GetPort(session *Session, self SDNControllerRef) (retval int, err error) {
	method := "SDN_controller.get_port"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPort2: Get the port field of the given SDN_controller.
// Version: falcon
func (sdnController) GetPort2(session *Session, self SDNControllerRef) (retval int, err error) {
	method := "SDN_controller.get_port"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAddress: Get the address field of the given SDN_controller.
// Version: falcon
func (sdnController) GetAddress(session *Session, self SDNControllerRef) (retval string, err error) {
	method := "SDN_controller.get_address"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAddress2: Get the address field of the given SDN_controller.
// Version: falcon
func (sdnController) GetAddress2(session *Session, self SDNControllerRef) (retval string, err error) {
	method := "SDN_controller.get_address"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetProtocol: Get the protocol field of the given SDN_controller.
// Version: falcon
func (sdnController) GetProtocol(session *Session, self SDNControllerRef) (retval SdnControllerProtocol, err error) {
	method := "SDN_controller.get_protocol"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumSdnControllerProtocol(method+" -> ", result)
	return
}

// GetProtocol2: Get the protocol field of the given SDN_controller.
// Version: falcon
func (sdnController) GetProtocol2(session *Session, self SDNControllerRef) (retval SdnControllerProtocol, err error) {
	method := "SDN_controller.get_protocol"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumSdnControllerProtocol(method+" -> ", result)
	return
}

// GetUUID: Get the uuid field of the given SDN_controller.
// Version: falcon
func (sdnController) GetUUID(session *Session, self SDNControllerRef) (retval string, err error) {
	method := "SDN_controller.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given SDN_controller.
// Version: falcon
func (sdnController) GetUUID2(session *Session, self SDNControllerRef) (retval string, err error) {
	method := "SDN_controller.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the SDN_controller instance with the specified UUID.
// Version: falcon
func (sdnController) GetByUUID(session *Session, uuid string) (retval SDNControllerRef, err error) {
	method := "SDN_controller.get_by_uuid"
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
	retval, err = deserializeSDNControllerRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the SDN_controller instance with the specified UUID.
// Version: falcon
func (sdnController) GetByUUID2(session *Session, uuid string) (retval SDNControllerRef, err error) {
	method := "SDN_controller.get_by_uuid"
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
	retval, err = deserializeSDNControllerRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given SDN_controller.
// Version: falcon
func (sdnController) GetRecord(session *Session, self SDNControllerRef) (retval SDNControllerRecord, err error) {
	method := "SDN_controller.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSDNControllerRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given SDN_controller.
// Version: falcon
func (sdnController) GetRecord2(session *Session, self SDNControllerRef) (retval SDNControllerRecord, err error) {
	method := "SDN_controller.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSDNControllerRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSDNControllerRecord(method+" -> ", result)
	return
}
