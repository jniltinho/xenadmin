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

type NetworkSriovRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// The PIF that has SR-IOV enabled
	PhysicalPIF PIFRef `json:"physical_PIF,omitempty"`
	// The logical PIF to connect to the SR-IOV network after enable SR-IOV on the physical PIF
	LogicalPIF PIFRef `json:"logical_PIF,omitempty"`
	// Indicates whether the host need to be rebooted before SR-IOV is enabled on the physical PIF
	RequiresReboot bool `json:"requires_reboot,omitempty"`
	// The mode for configure network sriov
	ConfigurationMode SriovConfigurationMode `json:"configuration_mode,omitempty"`
}

type NetworkSriovRef string

// network-sriov which connects logical pif and physical pif
type networkSriov struct{}

var NetworkSriov networkSriov

// GetAllRecords: Return a map of network_sriov references to network_sriov records for all network_sriovs known to the system.
// Version: kolkata
func (networkSriov) GetAllRecords(session *Session) (retval map[NetworkSriovRef]NetworkSriovRecord, err error) {
	method := "network_sriov.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkSriovRefToNetworkSriovRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of network_sriov references to network_sriov records for all network_sriovs known to the system.
// Version: kolkata
func (networkSriov) GetAllRecords1(session *Session) (retval map[NetworkSriovRef]NetworkSriovRecord, err error) {
	method := "network_sriov.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkSriovRefToNetworkSriovRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the network_sriovs known to the system.
// Version: kolkata
func (networkSriov) GetAll(session *Session) (retval []NetworkSriovRef, err error) {
	method := "network_sriov.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkSriovRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the network_sriovs known to the system.
// Version: kolkata
func (networkSriov) GetAll1(session *Session) (retval []NetworkSriovRef, err error) {
	method := "network_sriov.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkSriovRefSet(method+" -> ", result)
	return
}

// GetRemainingCapacity: Get the number of free SR-IOV VFs on the associated PIF
// Version: kolkata
func (networkSriov) GetRemainingCapacity(session *Session, self NetworkSriovRef) (retval int, err error) {
	method := "network_sriov.get_remaining_capacity"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncGetRemainingCapacity: Get the number of free SR-IOV VFs on the associated PIF
// Version: kolkata
func (networkSriov) AsyncGetRemainingCapacity(session *Session, self NetworkSriovRef) (retval TaskRef, err error) {
	method := "Async.network_sriov.get_remaining_capacity"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRemainingCapacity2: Get the number of free SR-IOV VFs on the associated PIF
// Version: kolkata
func (networkSriov) GetRemainingCapacity2(session *Session, self NetworkSriovRef) (retval int, err error) {
	method := "network_sriov.get_remaining_capacity"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncGetRemainingCapacity2: Get the number of free SR-IOV VFs on the associated PIF
// Version: kolkata
func (networkSriov) AsyncGetRemainingCapacity2(session *Session, self NetworkSriovRef) (retval TaskRef, err error) {
	method := "Async.network_sriov.get_remaining_capacity"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy: Disable SR-IOV on the specific PIF. It will destroy the network-sriov and the logical PIF accordingly.
// Version: kolkata
func (networkSriov) Destroy(session *Session, self NetworkSriovRef) (err error) {
	method := "network_sriov.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Disable SR-IOV on the specific PIF. It will destroy the network-sriov and the logical PIF accordingly.
// Version: kolkata
func (networkSriov) AsyncDestroy(session *Session, self NetworkSriovRef) (retval TaskRef, err error) {
	method := "Async.network_sriov.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Disable SR-IOV on the specific PIF. It will destroy the network-sriov and the logical PIF accordingly.
// Version: kolkata
func (networkSriov) Destroy2(session *Session, self NetworkSriovRef) (err error) {
	method := "network_sriov.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Disable SR-IOV on the specific PIF. It will destroy the network-sriov and the logical PIF accordingly.
// Version: kolkata
func (networkSriov) AsyncDestroy2(session *Session, self NetworkSriovRef) (retval TaskRef, err error) {
	method := "Async.network_sriov.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Enable SR-IOV on the specific PIF. It will create a network-sriov based on the specific PIF and automatically create a logical PIF to connect the specific network.
// Version: kolkata
func (networkSriov) Create(session *Session, pif PIFRef, network NetworkRef) (retval NetworkSriovRef, err error) {
	method := "network_sriov.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkSriovRef(method+" -> ", result)
	return
}

// AsyncCreate: Enable SR-IOV on the specific PIF. It will create a network-sriov based on the specific PIF and automatically create a logical PIF to connect the specific network.
// Version: kolkata
func (networkSriov) AsyncCreate(session *Session, pif PIFRef, network NetworkRef) (retval TaskRef, err error) {
	method := "Async.network_sriov.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create3: Enable SR-IOV on the specific PIF. It will create a network-sriov based on the specific PIF and automatically create a logical PIF to connect the specific network.
// Version: kolkata
func (networkSriov) Create3(session *Session, pif PIFRef, network NetworkRef) (retval NetworkSriovRef, err error) {
	method := "network_sriov.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkSriovRef(method+" -> ", result)
	return
}

// AsyncCreate3: Enable SR-IOV on the specific PIF. It will create a network-sriov based on the specific PIF and automatically create a logical PIF to connect the specific network.
// Version: kolkata
func (networkSriov) AsyncCreate3(session *Session, pif PIFRef, network NetworkRef) (retval TaskRef, err error) {
	method := "Async.network_sriov.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetConfigurationMode: Get the configuration_mode field of the given network_sriov.
// Version: kolkata
func (networkSriov) GetConfigurationMode(session *Session, self NetworkSriovRef) (retval SriovConfigurationMode, err error) {
	method := "network_sriov.get_configuration_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumSriovConfigurationMode(method+" -> ", result)
	return
}

// GetConfigurationMode2: Get the configuration_mode field of the given network_sriov.
// Version: kolkata
func (networkSriov) GetConfigurationMode2(session *Session, self NetworkSriovRef) (retval SriovConfigurationMode, err error) {
	method := "network_sriov.get_configuration_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumSriovConfigurationMode(method+" -> ", result)
	return
}

// GetRequiresReboot: Get the requires_reboot field of the given network_sriov.
// Version: kolkata
func (networkSriov) GetRequiresReboot(session *Session, self NetworkSriovRef) (retval bool, err error) {
	method := "network_sriov.get_requires_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRequiresReboot2: Get the requires_reboot field of the given network_sriov.
// Version: kolkata
func (networkSriov) GetRequiresReboot2(session *Session, self NetworkSriovRef) (retval bool, err error) {
	method := "network_sriov.get_requires_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLogicalPIF: Get the logical_PIF field of the given network_sriov.
// Version: kolkata
func (networkSriov) GetLogicalPIF(session *Session, self NetworkSriovRef) (retval PIFRef, err error) {
	method := "network_sriov.get_logical_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// GetLogicalPIF2: Get the logical_PIF field of the given network_sriov.
// Version: kolkata
func (networkSriov) GetLogicalPIF2(session *Session, self NetworkSriovRef) (retval PIFRef, err error) {
	method := "network_sriov.get_logical_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// GetPhysicalPIF: Get the physical_PIF field of the given network_sriov.
// Version: kolkata
func (networkSriov) GetPhysicalPIF(session *Session, self NetworkSriovRef) (retval PIFRef, err error) {
	method := "network_sriov.get_physical_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// GetPhysicalPIF2: Get the physical_PIF field of the given network_sriov.
// Version: kolkata
func (networkSriov) GetPhysicalPIF2(session *Session, self NetworkSriovRef) (retval PIFRef, err error) {
	method := "network_sriov.get_physical_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// GetUUID: Get the uuid field of the given network_sriov.
// Version: kolkata
func (networkSriov) GetUUID(session *Session, self NetworkSriovRef) (retval string, err error) {
	method := "network_sriov.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given network_sriov.
// Version: kolkata
func (networkSriov) GetUUID2(session *Session, self NetworkSriovRef) (retval string, err error) {
	method := "network_sriov.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the network_sriov instance with the specified UUID.
// Version: kolkata
func (networkSriov) GetByUUID(session *Session, uuid string) (retval NetworkSriovRef, err error) {
	method := "network_sriov.get_by_uuid"
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
	retval, err = deserializeNetworkSriovRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the network_sriov instance with the specified UUID.
// Version: kolkata
func (networkSriov) GetByUUID2(session *Session, uuid string) (retval NetworkSriovRef, err error) {
	method := "network_sriov.get_by_uuid"
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
	retval, err = deserializeNetworkSriovRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given network_sriov.
// Version: kolkata
func (networkSriov) GetRecord(session *Session, self NetworkSriovRef) (retval NetworkSriovRecord, err error) {
	method := "network_sriov.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkSriovRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given network_sriov.
// Version: kolkata
func (networkSriov) GetRecord2(session *Session, self NetworkSriovRef) (retval NetworkSriovRecord, err error) {
	method := "network_sriov.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeNetworkSriovRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkSriovRecord(method+" -> ", result)
	return
}
