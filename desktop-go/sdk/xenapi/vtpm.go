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

type VTPMRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// list of the operations allowed in this state. This list is advisory only and the server state may have changed by the time this field is read by a client.
	AllowedOperations []VtpmOperations `json:"allowed_operations,omitempty"`
	// links each of the running tasks using this object (by reference) to a current_operation enum which describes the nature of the task.
	CurrentOperations map[string]VtpmOperations `json:"current_operations,omitempty"`
	// The virtual machine the TPM is attached to
	VM VMRef `json:"VM,omitempty"`
	// The domain where the backend is located (unused)
	Backend VMRef `json:"backend,omitempty"`
	// The backend where the vTPM is persisted
	PersistenceBackend PersistenceBackend `json:"persistence_backend,omitempty"`
	// Whether the contents are never copied, satisfying the TPM spec
	IsUnique bool `json:"is_unique,omitempty"`
	// Whether the contents of the VTPM are secured according to the TPM spec
	IsProtected bool `json:"is_protected,omitempty"`
}

type VTPMRef string

// A virtual TPM device
type vtpm struct{}

var VTPM vtpm

// GetAllRecords: Return a map of VTPM references to VTPM records for all VTPMs known to the system.
// Version: closed
func (vtpm) GetAllRecords(session *Session) (retval map[VTPMRef]VTPMRecord, err error) {
	method := "VTPM.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVTPMRefToVTPMRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of VTPM references to VTPM records for all VTPMs known to the system.
// Version: closed
func (vtpm) GetAllRecords1(session *Session) (retval map[VTPMRef]VTPMRecord, err error) {
	method := "VTPM.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVTPMRefToVTPMRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the VTPMs known to the system.
// Version: closed
func (vtpm) GetAll(session *Session) (retval []VTPMRef, err error) {
	method := "VTPM.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVTPMRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the VTPMs known to the system.
// Version: closed
func (vtpm) GetAll1(session *Session) (retval []VTPMRef, err error) {
	method := "VTPM.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVTPMRefSet(method+" -> ", result)
	return
}

// Destroy: Destroy the specified VTPM instance, along with its state.
// Version: closed
func (vtpm) Destroy(session *Session, self VTPMRef) (err error) {
	method := "VTPM.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy the specified VTPM instance, along with its state.
// Version: closed
func (vtpm) AsyncDestroy(session *Session, self VTPMRef) (retval TaskRef, err error) {
	method := "Async.VTPM.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy the specified VTPM instance, along with its state.
// Version: closed
func (vtpm) Destroy2(session *Session, self VTPMRef) (err error) {
	method := "VTPM.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy the specified VTPM instance, along with its state.
// Version: closed
func (vtpm) AsyncDestroy2(session *Session, self VTPMRef) (retval TaskRef, err error) {
	method := "Async.VTPM.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Create a new VTPM instance, and return its handle.
// Version: closed
func (vtpm) Create(session *Session, vm VMRef, isUnique bool) (retval VTPMRef, err error) {
	method := "VTPM.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vM"), vm)
	if err != nil {
		return
	}
	isUniqueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "is_unique"), isUnique)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, isUniqueArg)
	if err != nil {
		return
	}
	retval, err = deserializeVTPMRef(method+" -> ", result)
	return
}

// AsyncCreate: Create a new VTPM instance, and return its handle.
// Version: closed
func (vtpm) AsyncCreate(session *Session, vm VMRef, isUnique bool) (retval TaskRef, err error) {
	method := "Async.VTPM.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vM"), vm)
	if err != nil {
		return
	}
	isUniqueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "is_unique"), isUnique)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, isUniqueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create3: Create a new VTPM instance, and return its handle.
// Version: closed
func (vtpm) Create3(session *Session, vm VMRef, isUnique bool) (retval VTPMRef, err error) {
	method := "VTPM.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vM"), vm)
	if err != nil {
		return
	}
	isUniqueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "is_unique"), isUnique)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, isUniqueArg)
	if err != nil {
		return
	}
	retval, err = deserializeVTPMRef(method+" -> ", result)
	return
}

// AsyncCreate3: Create a new VTPM instance, and return its handle.
// Version: closed
func (vtpm) AsyncCreate3(session *Session, vm VMRef, isUnique bool) (retval TaskRef, err error) {
	method := "Async.VTPM.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vM"), vm)
	if err != nil {
		return
	}
	isUniqueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "is_unique"), isUnique)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, isUniqueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetIsProtected: Get the is_protected field of the given VTPM.
// Version: closed
func (vtpm) GetIsProtected(session *Session, self VTPMRef) (retval bool, err error) {
	method := "VTPM.get_is_protected"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsProtected2: Get the is_protected field of the given VTPM.
// Version: closed
func (vtpm) GetIsProtected2(session *Session, self VTPMRef) (retval bool, err error) {
	method := "VTPM.get_is_protected"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsUnique: Get the is_unique field of the given VTPM.
// Version: closed
func (vtpm) GetIsUnique(session *Session, self VTPMRef) (retval bool, err error) {
	method := "VTPM.get_is_unique"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsUnique2: Get the is_unique field of the given VTPM.
// Version: closed
func (vtpm) GetIsUnique2(session *Session, self VTPMRef) (retval bool, err error) {
	method := "VTPM.get_is_unique"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPersistenceBackend: Get the persistence_backend field of the given VTPM.
// Version: closed
func (vtpm) GetPersistenceBackend(session *Session, self VTPMRef) (retval PersistenceBackend, err error) {
	method := "VTPM.get_persistence_backend"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPersistenceBackend(method+" -> ", result)
	return
}

// GetPersistenceBackend2: Get the persistence_backend field of the given VTPM.
// Version: closed
func (vtpm) GetPersistenceBackend2(session *Session, self VTPMRef) (retval PersistenceBackend, err error) {
	method := "VTPM.get_persistence_backend"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPersistenceBackend(method+" -> ", result)
	return
}

// GetBackend: Get the backend field of the given VTPM.
// Version: closed
func (vtpm) GetBackend(session *Session, self VTPMRef) (retval VMRef, err error) {
	method := "VTPM.get_backend"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// GetBackend2: Get the backend field of the given VTPM.
// Version: closed
func (vtpm) GetBackend2(session *Session, self VTPMRef) (retval VMRef, err error) {
	method := "VTPM.get_backend"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// GetVM: Get the VM field of the given VTPM.
// Version: closed
func (vtpm) GetVM(session *Session, self VTPMRef) (retval VMRef, err error) {
	method := "VTPM.get_VM"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// GetVM2: Get the VM field of the given VTPM.
// Version: closed
func (vtpm) GetVM2(session *Session, self VTPMRef) (retval VMRef, err error) {
	method := "VTPM.get_VM"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// GetCurrentOperations: Get the current_operations field of the given VTPM.
// Version: closed
func (vtpm) GetCurrentOperations(session *Session, self VTPMRef) (retval map[string]VtpmOperations, err error) {
	method := "VTPM.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumVtpmOperationsMap(method+" -> ", result)
	return
}

// GetCurrentOperations2: Get the current_operations field of the given VTPM.
// Version: closed
func (vtpm) GetCurrentOperations2(session *Session, self VTPMRef) (retval map[string]VtpmOperations, err error) {
	method := "VTPM.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumVtpmOperationsMap(method+" -> ", result)
	return
}

// GetAllowedOperations: Get the allowed_operations field of the given VTPM.
// Version: closed
func (vtpm) GetAllowedOperations(session *Session, self VTPMRef) (retval []VtpmOperations, err error) {
	method := "VTPM.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVtpmOperationsSet(method+" -> ", result)
	return
}

// GetAllowedOperations2: Get the allowed_operations field of the given VTPM.
// Version: closed
func (vtpm) GetAllowedOperations2(session *Session, self VTPMRef) (retval []VtpmOperations, err error) {
	method := "VTPM.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVtpmOperationsSet(method+" -> ", result)
	return
}

// GetUUID: Get the uuid field of the given VTPM.
// Version: closed
func (vtpm) GetUUID(session *Session, self VTPMRef) (retval string, err error) {
	method := "VTPM.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given VTPM.
// Version: closed
func (vtpm) GetUUID2(session *Session, self VTPMRef) (retval string, err error) {
	method := "VTPM.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the VTPM instance with the specified UUID.
// Version: closed
func (vtpm) GetByUUID(session *Session, uuid string) (retval VTPMRef, err error) {
	method := "VTPM.get_by_uuid"
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
	retval, err = deserializeVTPMRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the VTPM instance with the specified UUID.
// Version: closed
func (vtpm) GetByUUID2(session *Session, uuid string) (retval VTPMRef, err error) {
	method := "VTPM.get_by_uuid"
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
	retval, err = deserializeVTPMRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given VTPM.
// Version: closed
func (vtpm) GetRecord(session *Session, self VTPMRef) (retval VTPMRecord, err error) {
	method := "VTPM.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVTPMRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given VTPM.
// Version: closed
func (vtpm) GetRecord2(session *Session, self VTPMRef) (retval VTPMRecord, err error) {
	method := "VTPM.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVTPMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVTPMRecord(method+" -> ", result)
	return
}
