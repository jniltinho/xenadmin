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

type DRTaskRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// All SRs introduced by this appliance
	IntroducedSRs []SRRef `json:"introduced_SRs,omitempty"`
}

type DRTaskRef string

// DR task
type drTask struct{}

var DRTask drTask

// GetAllRecords: Return a map of DR_task references to DR_task records for all DR_tasks known to the system.
// Version: boston
func (drTask) GetAllRecords(session *Session) (retval map[DRTaskRef]DRTaskRecord, err error) {
	method := "DR_task.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeDRTaskRefToDRTaskRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of DR_task references to DR_task records for all DR_tasks known to the system.
// Version: boston
func (drTask) GetAllRecords1(session *Session) (retval map[DRTaskRef]DRTaskRecord, err error) {
	method := "DR_task.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeDRTaskRefToDRTaskRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the DR_tasks known to the system.
// Version: boston
func (drTask) GetAll(session *Session) (retval []DRTaskRef, err error) {
	method := "DR_task.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeDRTaskRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the DR_tasks known to the system.
// Version: boston
func (drTask) GetAll1(session *Session) (retval []DRTaskRef, err error) {
	method := "DR_task.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeDRTaskRefSet(method+" -> ", result)
	return
}

// Destroy: Destroy the disaster recovery task, detaching and forgetting any SRs introduced which are no longer required
// Version: boston
func (drTask) Destroy(session *Session, self DRTaskRef) (err error) {
	method := "DR_task.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDRTaskRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy the disaster recovery task, detaching and forgetting any SRs introduced which are no longer required
// Version: boston
func (drTask) AsyncDestroy(session *Session, self DRTaskRef) (retval TaskRef, err error) {
	method := "Async.DR_task.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDRTaskRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy the disaster recovery task, detaching and forgetting any SRs introduced which are no longer required
// Version: boston
func (drTask) Destroy2(session *Session, self DRTaskRef) (err error) {
	method := "DR_task.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDRTaskRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy the disaster recovery task, detaching and forgetting any SRs introduced which are no longer required
// Version: boston
func (drTask) AsyncDestroy2(session *Session, self DRTaskRef) (retval TaskRef, err error) {
	method := "Async.DR_task.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDRTaskRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Create a disaster recovery task which will query the supplied list of devices
// Version: boston
func (drTask) Create(session *Session, typeKey string, deviceConfig map[string]string, whitelist []string) (retval DRTaskRef, err error) {
	method := "DR_task.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	whitelistArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "whitelist"), whitelist)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, typeKeyArg, deviceConfigArg, whitelistArg)
	if err != nil {
		return
	}
	retval, err = deserializeDRTaskRef(method+" -> ", result)
	return
}

// AsyncCreate: Create a disaster recovery task which will query the supplied list of devices
// Version: boston
func (drTask) AsyncCreate(session *Session, typeKey string, deviceConfig map[string]string, whitelist []string) (retval TaskRef, err error) {
	method := "Async.DR_task.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	whitelistArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "whitelist"), whitelist)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, typeKeyArg, deviceConfigArg, whitelistArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create4: Create a disaster recovery task which will query the supplied list of devices
// Version: boston
func (drTask) Create4(session *Session, typeKey string, deviceConfig map[string]string, whitelist []string) (retval DRTaskRef, err error) {
	method := "DR_task.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	whitelistArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "whitelist"), whitelist)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, typeKeyArg, deviceConfigArg, whitelistArg)
	if err != nil {
		return
	}
	retval, err = deserializeDRTaskRef(method+" -> ", result)
	return
}

// AsyncCreate4: Create a disaster recovery task which will query the supplied list of devices
// Version: boston
func (drTask) AsyncCreate4(session *Session, typeKey string, deviceConfig map[string]string, whitelist []string) (retval TaskRef, err error) {
	method := "Async.DR_task.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	deviceConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "device_config"), deviceConfig)
	if err != nil {
		return
	}
	whitelistArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "whitelist"), whitelist)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, typeKeyArg, deviceConfigArg, whitelistArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetIntroducedSRs: Get the introduced_SRs field of the given DR_task.
// Version: boston
func (drTask) GetIntroducedSRs(session *Session, self DRTaskRef) (retval []SRRef, err error) {
	method := "DR_task.get_introduced_SRs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDRTaskRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRefSet(method+" -> ", result)
	return
}

// GetIntroducedSRs2: Get the introduced_SRs field of the given DR_task.
// Version: boston
func (drTask) GetIntroducedSRs2(session *Session, self DRTaskRef) (retval []SRRef, err error) {
	method := "DR_task.get_introduced_SRs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDRTaskRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRefSet(method+" -> ", result)
	return
}

// GetUUID: Get the uuid field of the given DR_task.
// Version: boston
func (drTask) GetUUID(session *Session, self DRTaskRef) (retval string, err error) {
	method := "DR_task.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDRTaskRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given DR_task.
// Version: boston
func (drTask) GetUUID2(session *Session, self DRTaskRef) (retval string, err error) {
	method := "DR_task.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDRTaskRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the DR_task instance with the specified UUID.
// Version: boston
func (drTask) GetByUUID(session *Session, uuid string) (retval DRTaskRef, err error) {
	method := "DR_task.get_by_uuid"
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
	retval, err = deserializeDRTaskRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the DR_task instance with the specified UUID.
// Version: boston
func (drTask) GetByUUID2(session *Session, uuid string) (retval DRTaskRef, err error) {
	method := "DR_task.get_by_uuid"
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
	retval, err = deserializeDRTaskRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given DR_task.
// Version: boston
func (drTask) GetRecord(session *Session, self DRTaskRef) (retval DRTaskRecord, err error) {
	method := "DR_task.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDRTaskRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDRTaskRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given DR_task.
// Version: boston
func (drTask) GetRecord2(session *Session, self DRTaskRef) (retval DRTaskRecord, err error) {
	method := "DR_task.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDRTaskRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDRTaskRecord(method+" -> ", result)
	return
}
