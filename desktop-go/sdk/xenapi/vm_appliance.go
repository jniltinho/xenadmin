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

type VMApplianceRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// a human-readable name
	NameLabel string `json:"name_label,omitempty"`
	// a notes field containing human-readable description
	NameDescription string `json:"name_description,omitempty"`
	// list of the operations allowed in this state. This list is advisory only and the server state may have changed by the time this field is read by a client.
	AllowedOperations []VMApplianceOperation `json:"allowed_operations,omitempty"`
	// links each of the running tasks using this object (by reference) to a current_operation enum which describes the nature of the task.
	CurrentOperations map[string]VMApplianceOperation `json:"current_operations,omitempty"`
	// all VMs in this appliance
	VMs []VMRef `json:"VMs,omitempty"`
}

type VMApplianceRef string

// VM appliance
type vmAppliance struct{}

var VMAppliance vmAppliance

// GetAllRecords: Return a map of VM_appliance references to VM_appliance records for all VM_appliances known to the system.
// Version: boston
func (vmAppliance) GetAllRecords(session *Session) (retval map[VMApplianceRef]VMApplianceRecord, err error) {
	method := "VM_appliance.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMApplianceRefToVMApplianceRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of VM_appliance references to VM_appliance records for all VM_appliances known to the system.
// Version: boston
func (vmAppliance) GetAllRecords1(session *Session) (retval map[VMApplianceRef]VMApplianceRecord, err error) {
	method := "VM_appliance.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMApplianceRefToVMApplianceRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the VM_appliances known to the system.
// Version: boston
func (vmAppliance) GetAll(session *Session) (retval []VMApplianceRef, err error) {
	method := "VM_appliance.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMApplianceRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the VM_appliances known to the system.
// Version: boston
func (vmAppliance) GetAll1(session *Session) (retval []VMApplianceRef, err error) {
	method := "VM_appliance.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMApplianceRefSet(method+" -> ", result)
	return
}

// Recover: Recover the VM appliance
// Version: boston
//
// Errors:
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
func (vmAppliance) Recover(session *Session, self VMApplianceRef, sessionTo SessionRef, force bool) (err error) {
	method := "VM_appliance.recover"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg, forceArg)
	return
}

// AsyncRecover: Recover the VM appliance
// Version: boston
//
// Errors:
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
func (vmAppliance) AsyncRecover(session *Session, self VMApplianceRef, sessionTo SessionRef, force bool) (retval TaskRef, err error) {
	method := "Async.VM_appliance.recover"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Recover4: Recover the VM appliance
// Version: boston
//
// Errors:
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
func (vmAppliance) Recover4(session *Session, self VMApplianceRef, sessionTo SessionRef, force bool) (err error) {
	method := "VM_appliance.recover"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg, forceArg)
	return
}

// AsyncRecover4: Recover the VM appliance
// Version: boston
//
// Errors:
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
func (vmAppliance) AsyncRecover4(session *Session, self VMApplianceRef, sessionTo SessionRef, force bool) (retval TaskRef, err error) {
	method := "Async.VM_appliance.recover"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetSRsRequiredForRecovery: Get the list of SRs required by the VM appliance to recover.
// Version: creedence
func (vmAppliance) GetSRsRequiredForRecovery(session *Session, self VMApplianceRef, sessionTo SessionRef) (retval []SRRef, err error) {
	method := "VM_appliance.get_SRs_required_for_recovery"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRefSet(method+" -> ", result)
	return
}

// AsyncGetSRsRequiredForRecovery: Get the list of SRs required by the VM appliance to recover.
// Version: creedence
func (vmAppliance) AsyncGetSRsRequiredForRecovery(session *Session, self VMApplianceRef, sessionTo SessionRef) (retval TaskRef, err error) {
	method := "Async.VM_appliance.get_SRs_required_for_recovery"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetSRsRequiredForRecovery3: Get the list of SRs required by the VM appliance to recover.
// Version: creedence
func (vmAppliance) GetSRsRequiredForRecovery3(session *Session, self VMApplianceRef, sessionTo SessionRef) (retval []SRRef, err error) {
	method := "VM_appliance.get_SRs_required_for_recovery"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRefSet(method+" -> ", result)
	return
}

// AsyncGetSRsRequiredForRecovery3: Get the list of SRs required by the VM appliance to recover.
// Version: creedence
func (vmAppliance) AsyncGetSRsRequiredForRecovery3(session *Session, self VMApplianceRef, sessionTo SessionRef) (retval TaskRef, err error) {
	method := "Async.VM_appliance.get_SRs_required_for_recovery"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AssertCanBeRecovered: Assert whether all SRs required to recover this VM appliance are available.
// Version: boston
//
// Errors:
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
func (vmAppliance) AssertCanBeRecovered(session *Session, self VMApplianceRef, sessionTo SessionRef) (err error) {
	method := "VM_appliance.assert_can_be_recovered"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	return
}

// AsyncAssertCanBeRecovered: Assert whether all SRs required to recover this VM appliance are available.
// Version: boston
//
// Errors:
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
func (vmAppliance) AsyncAssertCanBeRecovered(session *Session, self VMApplianceRef, sessionTo SessionRef) (retval TaskRef, err error) {
	method := "Async.VM_appliance.assert_can_be_recovered"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AssertCanBeRecovered3: Assert whether all SRs required to recover this VM appliance are available.
// Version: boston
//
// Errors:
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
func (vmAppliance) AssertCanBeRecovered3(session *Session, self VMApplianceRef, sessionTo SessionRef) (err error) {
	method := "VM_appliance.assert_can_be_recovered"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	return
}

// AsyncAssertCanBeRecovered3: Assert whether all SRs required to recover this VM appliance are available.
// Version: boston
//
// Errors:
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
func (vmAppliance) AsyncAssertCanBeRecovered3(session *Session, self VMApplianceRef, sessionTo SessionRef) (retval TaskRef, err error) {
	method := "Async.VM_appliance.assert_can_be_recovered"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Shutdown: For each VM in the appliance, try to shut it down cleanly. If this fails, perform a hard shutdown of the VM.
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) Shutdown(session *Session, self VMApplianceRef) (err error) {
	method := "VM_appliance.shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncShutdown: For each VM in the appliance, try to shut it down cleanly. If this fails, perform a hard shutdown of the VM.
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) AsyncShutdown(session *Session, self VMApplianceRef) (retval TaskRef, err error) {
	method := "Async.VM_appliance.shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Shutdown2: For each VM in the appliance, try to shut it down cleanly. If this fails, perform a hard shutdown of the VM.
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) Shutdown2(session *Session, self VMApplianceRef) (err error) {
	method := "VM_appliance.shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncShutdown2: For each VM in the appliance, try to shut it down cleanly. If this fails, perform a hard shutdown of the VM.
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) AsyncShutdown2(session *Session, self VMApplianceRef) (retval TaskRef, err error) {
	method := "Async.VM_appliance.shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// HardShutdown: Perform a hard shutdown of all the VMs in the appliance
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) HardShutdown(session *Session, self VMApplianceRef) (err error) {
	method := "VM_appliance.hard_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncHardShutdown: Perform a hard shutdown of all the VMs in the appliance
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) AsyncHardShutdown(session *Session, self VMApplianceRef) (retval TaskRef, err error) {
	method := "Async.VM_appliance.hard_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// HardShutdown2: Perform a hard shutdown of all the VMs in the appliance
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) HardShutdown2(session *Session, self VMApplianceRef) (err error) {
	method := "VM_appliance.hard_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncHardShutdown2: Perform a hard shutdown of all the VMs in the appliance
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) AsyncHardShutdown2(session *Session, self VMApplianceRef) (retval TaskRef, err error) {
	method := "Async.VM_appliance.hard_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// CleanShutdown: Perform a clean shutdown of all the VMs in the appliance
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) CleanShutdown(session *Session, self VMApplianceRef) (err error) {
	method := "VM_appliance.clean_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncCleanShutdown: Perform a clean shutdown of all the VMs in the appliance
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) AsyncCleanShutdown(session *Session, self VMApplianceRef) (retval TaskRef, err error) {
	method := "Async.VM_appliance.clean_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// CleanShutdown2: Perform a clean shutdown of all the VMs in the appliance
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) CleanShutdown2(session *Session, self VMApplianceRef) (err error) {
	method := "VM_appliance.clean_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncCleanShutdown2: Perform a clean shutdown of all the VMs in the appliance
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) AsyncCleanShutdown2(session *Session, self VMApplianceRef) (retval TaskRef, err error) {
	method := "Async.VM_appliance.clean_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Start: Start all VMs in the appliance
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) Start(session *Session, self VMApplianceRef, paused bool) (err error) {
	method := "VM_appliance.start"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	pausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "paused"), paused)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, pausedArg)
	return
}

// AsyncStart: Start all VMs in the appliance
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) AsyncStart(session *Session, self VMApplianceRef, paused bool) (retval TaskRef, err error) {
	method := "Async.VM_appliance.start"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	pausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "paused"), paused)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, pausedArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Start3: Start all VMs in the appliance
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) Start3(session *Session, self VMApplianceRef, paused bool) (err error) {
	method := "VM_appliance.start"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	pausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "paused"), paused)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, pausedArg)
	return
}

// AsyncStart3: Start all VMs in the appliance
// Version: boston
//
// Errors:
// OPERATION_PARTIALLY_FAILED - Some VMs belonging to the appliance threw an exception while carrying out the specified operation
func (vmAppliance) AsyncStart3(session *Session, self VMApplianceRef, paused bool) (retval TaskRef, err error) {
	method := "Async.VM_appliance.start"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	pausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "paused"), paused)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, pausedArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetNameDescription: Set the name/description field of the given VM_appliance.
// Version: boston
func (vmAppliance) SetNameDescription(session *Session, self VMApplianceRef, value string) (err error) {
	method := "VM_appliance.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameDescription3: Set the name/description field of the given VM_appliance.
// Version: boston
func (vmAppliance) SetNameDescription3(session *Session, self VMApplianceRef, value string) (err error) {
	method := "VM_appliance.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameDescription2: Set the name/description field of the given VM_appliance.
// Version: rio
func (vmAppliance) SetNameDescription2(session *Session, value string) (err error) {
	method := "VM_appliance.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, valueArg)
	return
}

// SetNameLabel: Set the name/label field of the given VM_appliance.
// Version: boston
func (vmAppliance) SetNameLabel(session *Session, self VMApplianceRef, value string) (err error) {
	method := "VM_appliance.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameLabel3: Set the name/label field of the given VM_appliance.
// Version: boston
func (vmAppliance) SetNameLabel3(session *Session, self VMApplianceRef, value string) (err error) {
	method := "VM_appliance.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameLabel2: Set the name/label field of the given VM_appliance.
// Version: rio
func (vmAppliance) SetNameLabel2(session *Session, value string) (err error) {
	method := "VM_appliance.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, valueArg)
	return
}

// GetVMs: Get the VMs field of the given VM_appliance.
// Version: boston
func (vmAppliance) GetVMs(session *Session, self VMApplianceRef) (retval []VMRef, err error) {
	method := "VM_appliance.get_VMs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// GetVMs2: Get the VMs field of the given VM_appliance.
// Version: boston
func (vmAppliance) GetVMs2(session *Session, self VMApplianceRef) (retval []VMRef, err error) {
	method := "VM_appliance.get_VMs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// GetCurrentOperations: Get the current_operations field of the given VM_appliance.
// Version: boston
func (vmAppliance) GetCurrentOperations(session *Session, self VMApplianceRef) (retval map[string]VMApplianceOperation, err error) {
	method := "VM_appliance.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumVMApplianceOperationMap(method+" -> ", result)
	return
}

// GetCurrentOperations2: Get the current_operations field of the given VM_appliance.
// Version: boston
func (vmAppliance) GetCurrentOperations2(session *Session, self VMApplianceRef) (retval map[string]VMApplianceOperation, err error) {
	method := "VM_appliance.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumVMApplianceOperationMap(method+" -> ", result)
	return
}

// GetAllowedOperations: Get the allowed_operations field of the given VM_appliance.
// Version: boston
func (vmAppliance) GetAllowedOperations(session *Session, self VMApplianceRef) (retval []VMApplianceOperation, err error) {
	method := "VM_appliance.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVMApplianceOperationSet(method+" -> ", result)
	return
}

// GetAllowedOperations2: Get the allowed_operations field of the given VM_appliance.
// Version: boston
func (vmAppliance) GetAllowedOperations2(session *Session, self VMApplianceRef) (retval []VMApplianceOperation, err error) {
	method := "VM_appliance.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVMApplianceOperationSet(method+" -> ", result)
	return
}

// GetNameDescription: Get the name/description field of the given VM_appliance.
// Version: boston
func (vmAppliance) GetNameDescription(session *Session, self VMApplianceRef) (retval string, err error) {
	method := "VM_appliance.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription2: Get the name/description field of the given VM_appliance.
// Version: boston
func (vmAppliance) GetNameDescription2(session *Session, self VMApplianceRef) (retval string, err error) {
	method := "VM_appliance.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel: Get the name/label field of the given VM_appliance.
// Version: boston
func (vmAppliance) GetNameLabel(session *Session, self VMApplianceRef) (retval string, err error) {
	method := "VM_appliance.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel2: Get the name/label field of the given VM_appliance.
// Version: boston
func (vmAppliance) GetNameLabel2(session *Session, self VMApplianceRef) (retval string, err error) {
	method := "VM_appliance.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given VM_appliance.
// Version: boston
func (vmAppliance) GetUUID(session *Session, self VMApplianceRef) (retval string, err error) {
	method := "VM_appliance.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given VM_appliance.
// Version: boston
func (vmAppliance) GetUUID2(session *Session, self VMApplianceRef) (retval string, err error) {
	method := "VM_appliance.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByNameLabel: Get all the VM_appliance instances with the given label.
// Version: boston
func (vmAppliance) GetByNameLabel(session *Session, label string) (retval []VMApplianceRef, err error) {
	method := "VM_appliance.get_by_name_label"
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
	retval, err = deserializeVMApplianceRefSet(method+" -> ", result)
	return
}

// GetByNameLabel2: Get all the VM_appliance instances with the given label.
// Version: boston
func (vmAppliance) GetByNameLabel2(session *Session, label string) (retval []VMApplianceRef, err error) {
	method := "VM_appliance.get_by_name_label"
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
	retval, err = deserializeVMApplianceRefSet(method+" -> ", result)
	return
}

// Destroy: Destroy the specified VM_appliance instance.
// Version: boston
func (vmAppliance) Destroy(session *Session, self VMApplianceRef) (err error) {
	method := "VM_appliance.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy the specified VM_appliance instance.
// Version: boston
func (vmAppliance) AsyncDestroy(session *Session, self VMApplianceRef) (retval TaskRef, err error) {
	method := "Async.VM_appliance.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy the specified VM_appliance instance.
// Version: boston
func (vmAppliance) Destroy2(session *Session, self VMApplianceRef) (err error) {
	method := "VM_appliance.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy the specified VM_appliance instance.
// Version: boston
func (vmAppliance) AsyncDestroy2(session *Session, self VMApplianceRef) (retval TaskRef, err error) {
	method := "Async.VM_appliance.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Create a new VM_appliance instance, and return its handle. The constructor args are: name_label, name_description (* = non-optional).
// Version: boston
func (vmAppliance) Create(session *Session, args VMApplianceRecord) (retval VMApplianceRef, err error) {
	method := "VM_appliance.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVMApplianceRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMApplianceRef(method+" -> ", result)
	return
}

// AsyncCreate: Create a new VM_appliance instance, and return its handle. The constructor args are: name_label, name_description (* = non-optional).
// Version: boston
func (vmAppliance) AsyncCreate(session *Session, args VMApplianceRecord) (retval TaskRef, err error) {
	method := "Async.VM_appliance.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVMApplianceRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
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

// Create2: Create a new VM_appliance instance, and return its handle. The constructor args are: name_label, name_description (* = non-optional).
// Version: boston
func (vmAppliance) Create2(session *Session, args VMApplianceRecord) (retval VMApplianceRef, err error) {
	method := "VM_appliance.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVMApplianceRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMApplianceRef(method+" -> ", result)
	return
}

// AsyncCreate2: Create a new VM_appliance instance, and return its handle. The constructor args are: name_label, name_description (* = non-optional).
// Version: boston
func (vmAppliance) AsyncCreate2(session *Session, args VMApplianceRecord) (retval TaskRef, err error) {
	method := "Async.VM_appliance.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVMApplianceRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
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

// GetByUUID: Get a reference to the VM_appliance instance with the specified UUID.
// Version: boston
func (vmAppliance) GetByUUID(session *Session, uuid string) (retval VMApplianceRef, err error) {
	method := "VM_appliance.get_by_uuid"
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
	retval, err = deserializeVMApplianceRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the VM_appliance instance with the specified UUID.
// Version: boston
func (vmAppliance) GetByUUID2(session *Session, uuid string) (retval VMApplianceRef, err error) {
	method := "VM_appliance.get_by_uuid"
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
	retval, err = deserializeVMApplianceRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given VM_appliance.
// Version: boston
func (vmAppliance) GetRecord(session *Session, self VMApplianceRef) (retval VMApplianceRecord, err error) {
	method := "VM_appliance.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMApplianceRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given VM_appliance.
// Version: boston
func (vmAppliance) GetRecord2(session *Session, self VMApplianceRef) (retval VMApplianceRecord, err error) {
	method := "VM_appliance.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMApplianceRecord(method+" -> ", result)
	return
}
