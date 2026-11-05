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

type DriverVariantRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// Name identifying the driver variant within the driver
	Name string `json:"name,omitempty"`
	// Driver this variant is a part of
	Driver HostDriverRef `json:"driver,omitempty"`
	// Unique version of this driver variant
	Version string `json:"version,omitempty"`
	// True if the hardware for this variant is present on the host
	HardwarePresent bool `json:"hardware_present,omitempty"`
	// Priority; this needs an explanation how this is ordered
	Priority float64 `json:"priority,omitempty"`
	// Development and release status of this variant, like &apos;alpha&apos;
	Status string `json:"status,omitempty"`
}

type DriverVariantRef string

// UNSUPPORTED. Variant of a host driver
type driverVariant struct{}

var DriverVariant driverVariant

// GetAllRecords: Return a map of Driver_variant references to Driver_variant records for all Driver_variants known to the system.
// Version: closed
func (driverVariant) GetAllRecords(session *Session) (retval map[DriverVariantRef]DriverVariantRecord, err error) {
	method := "Driver_variant.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeDriverVariantRefToDriverVariantRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of Driver_variant references to Driver_variant records for all Driver_variants known to the system.
// Version: closed
func (driverVariant) GetAllRecords1(session *Session) (retval map[DriverVariantRef]DriverVariantRecord, err error) {
	method := "Driver_variant.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeDriverVariantRefToDriverVariantRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the Driver_variants known to the system.
// Version: closed
func (driverVariant) GetAll(session *Session) (retval []DriverVariantRef, err error) {
	method := "Driver_variant.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeDriverVariantRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the Driver_variants known to the system.
// Version: closed
func (driverVariant) GetAll1(session *Session) (retval []DriverVariantRef, err error) {
	method := "Driver_variant.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeDriverVariantRefSet(method+" -> ", result)
	return
}

// Select: UNSUPPORTED Select this variant of a driver to become active after reboot or immediately if currently no version is active
// Version: closed
func (driverVariant) Select(session *Session, self DriverVariantRef) (err error) {
	method := "Driver_variant.select"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncSelect: UNSUPPORTED Select this variant of a driver to become active after reboot or immediately if currently no version is active
// Version: closed
func (driverVariant) AsyncSelect(session *Session, self DriverVariantRef) (retval TaskRef, err error) {
	method := "Async.Driver_variant.select"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Select2: UNSUPPORTED Select this variant of a driver to become active after reboot or immediately if currently no version is active
// Version: closed
func (driverVariant) Select2(session *Session, self DriverVariantRef) (err error) {
	method := "Driver_variant.select"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncSelect2: UNSUPPORTED Select this variant of a driver to become active after reboot or immediately if currently no version is active
// Version: closed
func (driverVariant) AsyncSelect2(session *Session, self DriverVariantRef) (retval TaskRef, err error) {
	method := "Async.Driver_variant.select"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetStatus: Get the status field of the given Driver_variant.
// Version: closed
func (driverVariant) GetStatus(session *Session, self DriverVariantRef) (retval string, err error) {
	method := "Driver_variant.get_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetStatus2: Get the status field of the given Driver_variant.
// Version: closed
func (driverVariant) GetStatus2(session *Session, self DriverVariantRef) (retval string, err error) {
	method := "Driver_variant.get_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPriority: Get the priority field of the given Driver_variant.
// Version: closed
func (driverVariant) GetPriority(session *Session, self DriverVariantRef) (retval float64, err error) {
	method := "Driver_variant.get_priority"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeFloat(method+" -> ", result)
	return
}

// GetPriority2: Get the priority field of the given Driver_variant.
// Version: closed
func (driverVariant) GetPriority2(session *Session, self DriverVariantRef) (retval float64, err error) {
	method := "Driver_variant.get_priority"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeFloat(method+" -> ", result)
	return
}

// GetHardwarePresent: Get the hardware_present field of the given Driver_variant.
// Version: closed
func (driverVariant) GetHardwarePresent(session *Session, self DriverVariantRef) (retval bool, err error) {
	method := "Driver_variant.get_hardware_present"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHardwarePresent2: Get the hardware_present field of the given Driver_variant.
// Version: closed
func (driverVariant) GetHardwarePresent2(session *Session, self DriverVariantRef) (retval bool, err error) {
	method := "Driver_variant.get_hardware_present"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVersion: Get the version field of the given Driver_variant.
// Version: closed
func (driverVariant) GetVersion(session *Session, self DriverVariantRef) (retval string, err error) {
	method := "Driver_variant.get_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVersion2: Get the version field of the given Driver_variant.
// Version: closed
func (driverVariant) GetVersion2(session *Session, self DriverVariantRef) (retval string, err error) {
	method := "Driver_variant.get_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDriver: Get the driver field of the given Driver_variant.
// Version: closed
func (driverVariant) GetDriver(session *Session, self DriverVariantRef) (retval HostDriverRef, err error) {
	method := "Driver_variant.get_driver"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostDriverRef(method+" -> ", result)
	return
}

// GetDriver2: Get the driver field of the given Driver_variant.
// Version: closed
func (driverVariant) GetDriver2(session *Session, self DriverVariantRef) (retval HostDriverRef, err error) {
	method := "Driver_variant.get_driver"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostDriverRef(method+" -> ", result)
	return
}

// GetName: Get the name field of the given Driver_variant.
// Version: closed
func (driverVariant) GetName(session *Session, self DriverVariantRef) (retval string, err error) {
	method := "Driver_variant.get_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetName2: Get the name field of the given Driver_variant.
// Version: closed
func (driverVariant) GetName2(session *Session, self DriverVariantRef) (retval string, err error) {
	method := "Driver_variant.get_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given Driver_variant.
// Version: closed
func (driverVariant) GetUUID(session *Session, self DriverVariantRef) (retval string, err error) {
	method := "Driver_variant.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given Driver_variant.
// Version: closed
func (driverVariant) GetUUID2(session *Session, self DriverVariantRef) (retval string, err error) {
	method := "Driver_variant.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the Driver_variant instance with the specified UUID.
// Version: closed
func (driverVariant) GetByUUID(session *Session, uuid string) (retval DriverVariantRef, err error) {
	method := "Driver_variant.get_by_uuid"
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
	retval, err = deserializeDriverVariantRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the Driver_variant instance with the specified UUID.
// Version: closed
func (driverVariant) GetByUUID2(session *Session, uuid string) (retval DriverVariantRef, err error) {
	method := "Driver_variant.get_by_uuid"
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
	retval, err = deserializeDriverVariantRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given Driver_variant.
// Version: closed
func (driverVariant) GetRecord(session *Session, self DriverVariantRef) (retval DriverVariantRecord, err error) {
	method := "Driver_variant.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDriverVariantRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given Driver_variant.
// Version: closed
func (driverVariant) GetRecord2(session *Session, self DriverVariantRef) (retval DriverVariantRecord, err error) {
	method := "Driver_variant.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDriverVariantRecord(method+" -> ", result)
	return
}
