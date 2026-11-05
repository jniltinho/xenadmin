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

type SubjectRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// the subject identifier, unique in the external directory service
	SubjectIdentifier string `json:"subject_identifier,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// the roles associated with this subject
	Roles []RoleRef `json:"roles,omitempty"`
}

type SubjectRef string

// A user or group that can log in xapi
type subject struct{}

var Subject subject

// GetAllRecords: Return a map of subject references to subject records for all subjects known to the system.
// Version: george
func (subject) GetAllRecords(session *Session) (retval map[SubjectRef]SubjectRecord, err error) {
	method := "subject.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeSubjectRefToSubjectRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of subject references to subject records for all subjects known to the system.
// Version: george
func (subject) GetAllRecords1(session *Session) (retval map[SubjectRef]SubjectRecord, err error) {
	method := "subject.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeSubjectRefToSubjectRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the subjects known to the system.
// Version: george
func (subject) GetAll(session *Session) (retval []SubjectRef, err error) {
	method := "subject.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeSubjectRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the subjects known to the system.
// Version: george
func (subject) GetAll1(session *Session) (retval []SubjectRef, err error) {
	method := "subject.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeSubjectRefSet(method+" -> ", result)
	return
}

// GetPermissionsNameLabel: This call returns a list of permission names given a subject
// Version: midnight-ride
func (subject) GetPermissionsNameLabel(session *Session, self SubjectRef) (retval []string, err error) {
	method := "subject.get_permissions_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPermissionsNameLabel2: This call returns a list of permission names given a subject
// Version: midnight-ride
func (subject) GetPermissionsNameLabel2(session *Session, self SubjectRef) (retval []string, err error) {
	method := "subject.get_permissions_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromRoles: This call removes a role from a subject
// Version: midnight-ride
func (subject) RemoveFromRoles(session *Session, self SubjectRef, role RoleRef) (err error) {
	method := "subject.remove_from_roles"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	roleArg, err := serializeRoleRef(fmt.Sprintf("%s(%s)", method, "role"), role)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, roleArg)
	return
}

// RemoveFromRoles3: This call removes a role from a subject
// Version: midnight-ride
func (subject) RemoveFromRoles3(session *Session, self SubjectRef, role RoleRef) (err error) {
	method := "subject.remove_from_roles"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	roleArg, err := serializeRoleRef(fmt.Sprintf("%s(%s)", method, "role"), role)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, roleArg)
	return
}

// AddToRoles: This call adds a new role to a subject
// Version: midnight-ride
func (subject) AddToRoles(session *Session, self SubjectRef, role RoleRef) (err error) {
	method := "subject.add_to_roles"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	roleArg, err := serializeRoleRef(fmt.Sprintf("%s(%s)", method, "role"), role)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, roleArg)
	return
}

// AddToRoles3: This call adds a new role to a subject
// Version: midnight-ride
func (subject) AddToRoles3(session *Session, self SubjectRef, role RoleRef) (err error) {
	method := "subject.add_to_roles"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	roleArg, err := serializeRoleRef(fmt.Sprintf("%s(%s)", method, "role"), role)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, roleArg)
	return
}

// GetRoles: Get the roles field of the given subject.
// Version: george
func (subject) GetRoles(session *Session, self SubjectRef) (retval []RoleRef, err error) {
	method := "subject.get_roles"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeRoleRefSet(method+" -> ", result)
	return
}

// GetRoles2: Get the roles field of the given subject.
// Version: george
func (subject) GetRoles2(session *Session, self SubjectRef) (retval []RoleRef, err error) {
	method := "subject.get_roles"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeRoleRefSet(method+" -> ", result)
	return
}

// GetOtherConfig: Get the other_config field of the given subject.
// Version: george
func (subject) GetOtherConfig(session *Session, self SubjectRef) (retval map[string]string, err error) {
	method := "subject.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetOtherConfig2: Get the other_config field of the given subject.
// Version: george
func (subject) GetOtherConfig2(session *Session, self SubjectRef) (retval map[string]string, err error) {
	method := "subject.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToStringMap(method+" -> ", result)
	return
}

// GetSubjectIdentifier: Get the subject_identifier field of the given subject.
// Version: george
func (subject) GetSubjectIdentifier(session *Session, self SubjectRef) (retval string, err error) {
	method := "subject.get_subject_identifier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSubjectIdentifier2: Get the subject_identifier field of the given subject.
// Version: george
func (subject) GetSubjectIdentifier2(session *Session, self SubjectRef) (retval string, err error) {
	method := "subject.get_subject_identifier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given subject.
// Version: george
func (subject) GetUUID(session *Session, self SubjectRef) (retval string, err error) {
	method := "subject.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given subject.
// Version: george
func (subject) GetUUID2(session *Session, self SubjectRef) (retval string, err error) {
	method := "subject.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy: Destroy the specified subject instance.
// Version: george
func (subject) Destroy(session *Session, self SubjectRef) (err error) {
	method := "subject.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy the specified subject instance.
// Version: george
func (subject) AsyncDestroy(session *Session, self SubjectRef) (retval TaskRef, err error) {
	method := "Async.subject.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy the specified subject instance.
// Version: george
func (subject) Destroy2(session *Session, self SubjectRef) (err error) {
	method := "subject.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy the specified subject instance.
// Version: george
func (subject) AsyncDestroy2(session *Session, self SubjectRef) (retval TaskRef, err error) {
	method := "Async.subject.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Create a new subject instance, and return its handle. The constructor args are: subject_identifier, other_config (* = non-optional).
// Version: george
func (subject) Create(session *Session, args SubjectRecord) (retval SubjectRef, err error) {
	method := "subject.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeSubjectRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeSubjectRef(method+" -> ", result)
	return
}

// AsyncCreate: Create a new subject instance, and return its handle. The constructor args are: subject_identifier, other_config (* = non-optional).
// Version: george
func (subject) AsyncCreate(session *Session, args SubjectRecord) (retval TaskRef, err error) {
	method := "Async.subject.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeSubjectRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
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

// Create2: Create a new subject instance, and return its handle. The constructor args are: subject_identifier, other_config (* = non-optional).
// Version: george
func (subject) Create2(session *Session, args SubjectRecord) (retval SubjectRef, err error) {
	method := "subject.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeSubjectRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeSubjectRef(method+" -> ", result)
	return
}

// AsyncCreate2: Create a new subject instance, and return its handle. The constructor args are: subject_identifier, other_config (* = non-optional).
// Version: george
func (subject) AsyncCreate2(session *Session, args SubjectRecord) (retval TaskRef, err error) {
	method := "Async.subject.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeSubjectRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
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

// GetByUUID: Get a reference to the subject instance with the specified UUID.
// Version: george
func (subject) GetByUUID(session *Session, uuid string) (retval SubjectRef, err error) {
	method := "subject.get_by_uuid"
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
	retval, err = deserializeSubjectRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the subject instance with the specified UUID.
// Version: george
func (subject) GetByUUID2(session *Session, uuid string) (retval SubjectRef, err error) {
	method := "subject.get_by_uuid"
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
	retval, err = deserializeSubjectRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given subject.
// Version: george
func (subject) GetRecord(session *Session, self SubjectRef) (retval SubjectRecord, err error) {
	method := "subject.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSubjectRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given subject.
// Version: george
func (subject) GetRecord2(session *Session, self SubjectRef) (retval SubjectRecord, err error) {
	method := "subject.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeSubjectRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSubjectRecord(method+" -> ", result)
	return
}
