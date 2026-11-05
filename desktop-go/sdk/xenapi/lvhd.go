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

type LVHDRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
}

type LVHDRef string

// LVHD SR specific operations
type lvhd struct{}

var LVHD lvhd

// EnableThinProvisioning: Upgrades an LVHD SR to enable thin-provisioning. Future VDIs created in this SR will be thinly-provisioned, although existing VDIs will be left alone. Note that the SR must be attached to the SRmaster for upgrade to work.
// Version: dundee
func (lvhd) EnableThinProvisioning(session *Session, host HostRef, sr SRRef, initialAllocation int, allocationQuantum int) (retval string, err error) {
	method := "LVHD.enable_thin_provisioning"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "SR"), sr)
	if err != nil {
		return
	}
	initialAllocationArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "initial_allocation"), initialAllocation)
	if err != nil {
		return
	}
	allocationQuantumArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "allocation_quantum"), allocationQuantum)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, srArg, initialAllocationArg, allocationQuantumArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncEnableThinProvisioning: Upgrades an LVHD SR to enable thin-provisioning. Future VDIs created in this SR will be thinly-provisioned, although existing VDIs will be left alone. Note that the SR must be attached to the SRmaster for upgrade to work.
// Version: dundee
func (lvhd) AsyncEnableThinProvisioning(session *Session, host HostRef, sr SRRef, initialAllocation int, allocationQuantum int) (retval TaskRef, err error) {
	method := "Async.LVHD.enable_thin_provisioning"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "SR"), sr)
	if err != nil {
		return
	}
	initialAllocationArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "initial_allocation"), initialAllocation)
	if err != nil {
		return
	}
	allocationQuantumArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "allocation_quantum"), allocationQuantum)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, srArg, initialAllocationArg, allocationQuantumArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// EnableThinProvisioning5: Upgrades an LVHD SR to enable thin-provisioning. Future VDIs created in this SR will be thinly-provisioned, although existing VDIs will be left alone. Note that the SR must be attached to the SRmaster for upgrade to work.
// Version: dundee
func (lvhd) EnableThinProvisioning5(session *Session, host HostRef, sr SRRef, initialAllocation int, allocationQuantum int) (retval string, err error) {
	method := "LVHD.enable_thin_provisioning"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "SR"), sr)
	if err != nil {
		return
	}
	initialAllocationArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "initial_allocation"), initialAllocation)
	if err != nil {
		return
	}
	allocationQuantumArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "allocation_quantum"), allocationQuantum)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, srArg, initialAllocationArg, allocationQuantumArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncEnableThinProvisioning5: Upgrades an LVHD SR to enable thin-provisioning. Future VDIs created in this SR will be thinly-provisioned, although existing VDIs will be left alone. Note that the SR must be attached to the SRmaster for upgrade to work.
// Version: dundee
func (lvhd) AsyncEnableThinProvisioning5(session *Session, host HostRef, sr SRRef, initialAllocation int, allocationQuantum int) (retval TaskRef, err error) {
	method := "Async.LVHD.enable_thin_provisioning"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "SR"), sr)
	if err != nil {
		return
	}
	initialAllocationArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "initial_allocation"), initialAllocation)
	if err != nil {
		return
	}
	allocationQuantumArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "allocation_quantum"), allocationQuantum)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, srArg, initialAllocationArg, allocationQuantumArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetUUID: Get the uuid field of the given LVHD.
// Version: dundee
func (lvhd) GetUUID(session *Session, self LVHDRef) (retval string, err error) {
	method := "LVHD.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeLVHDRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given LVHD.
// Version: dundee
func (lvhd) GetUUID2(session *Session, self LVHDRef) (retval string, err error) {
	method := "LVHD.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeLVHDRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the LVHD instance with the specified UUID.
// Version: dundee
func (lvhd) GetByUUID(session *Session, uuid string) (retval LVHDRef, err error) {
	method := "LVHD.get_by_uuid"
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
	retval, err = deserializeLVHDRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the LVHD instance with the specified UUID.
// Version: dundee
func (lvhd) GetByUUID2(session *Session, uuid string) (retval LVHDRef, err error) {
	method := "LVHD.get_by_uuid"
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
	retval, err = deserializeLVHDRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given LVHD.
// Version: dundee
func (lvhd) GetRecord(session *Session, self LVHDRef) (retval LVHDRecord, err error) {
	method := "LVHD.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeLVHDRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeLVHDRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given LVHD.
// Version: dundee
func (lvhd) GetRecord2(session *Session, self LVHDRef) (retval LVHDRecord, err error) {
	method := "LVHD.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeLVHDRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeLVHDRecord(method+" -> ", result)
	return
}
