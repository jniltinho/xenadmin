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

type PCIRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// PCI class ID
	ClassID string `json:"class_id,omitempty"`
	// PCI class name
	ClassName string `json:"class_name,omitempty"`
	// Vendor ID
	VendorID string `json:"vendor_id,omitempty"`
	// Vendor name
	VendorName string `json:"vendor_name,omitempty"`
	// Device ID
	DeviceID string `json:"device_id,omitempty"`
	// Device name
	DeviceName string `json:"device_name,omitempty"`
	// Physical machine that owns the PCI device
	Host HostRef `json:"host,omitempty"`
	// PCI ID of the physical device
	PciID string `json:"pci_id,omitempty"`
	// List of dependent PCI devices
	Dependencies []PCIRef `json:"dependencies,omitempty"`
	// Additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// Subsystem vendor ID
	SubsystemVendorID string `json:"subsystem_vendor_id,omitempty"`
	// Subsystem vendor name
	SubsystemVendorName string `json:"subsystem_vendor_name,omitempty"`
	// Subsystem device ID
	SubsystemDeviceID string `json:"subsystem_device_id,omitempty"`
	// Subsystem device name
	SubsystemDeviceName string `json:"subsystem_device_name,omitempty"`
	// Driver name
	DriverName string `json:"driver_name,omitempty"`
}

type PCIRef string

// A PCI device
type pci struct{}

var PCI pci

// GetAllRecords: Return a map of PCI references to PCI records for all PCIs known to the system.
// Version: boston
func (pci) GetAllRecords(session *Session) (retval map[PCIRef]PCIRecord, err error) {
	method := "PCI.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRefToPCIRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of PCI references to PCI records for all PCIs known to the system.
// Version: boston
func (pci) GetAllRecords1(session *Session) (retval map[PCIRef]PCIRecord, err error) {
	method := "PCI.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRefToPCIRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the PCIs known to the system.
// Version: boston
func (pci) GetAll(session *Session) (retval []PCIRef, err error) {
	method := "PCI.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the PCIs known to the system.
// Version: boston
func (pci) GetAll1(session *Session) (retval []PCIRef, err error) {
	method := "PCI.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRefSet(method+" -> ", result)
	return
}

// GetDom0AccessStatus: Return a PCI device dom0 access status.
// Version: closed
func (pci) GetDom0AccessStatus(session *Session, self PCIRef) (retval PciDom0Access, err error) {
	method := "PCI.get_dom0_access_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPciDom0Access(method+" -> ", result)
	return
}

// AsyncGetDom0AccessStatus: Return a PCI device dom0 access status.
// Version: closed
func (pci) AsyncGetDom0AccessStatus(session *Session, self PCIRef) (retval TaskRef, err error) {
	method := "Async.PCI.get_dom0_access_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDom0AccessStatus2: Return a PCI device dom0 access status.
// Version: closed
func (pci) GetDom0AccessStatus2(session *Session, self PCIRef) (retval PciDom0Access, err error) {
	method := "PCI.get_dom0_access_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPciDom0Access(method+" -> ", result)
	return
}

// AsyncGetDom0AccessStatus2: Return a PCI device dom0 access status.
// Version: closed
func (pci) AsyncGetDom0AccessStatus2(session *Session, self PCIRef) (retval TaskRef, err error) {
	method := "Async.PCI.get_dom0_access_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// EnableDom0Access: Unhide a PCI device from the dom0 kernel. (Takes affect after next boot.)
// Version: closed
func (pci) EnableDom0Access(session *Session, self PCIRef) (retval PciDom0Access, err error) {
	method := "PCI.enable_dom0_access"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPciDom0Access(method+" -> ", result)
	return
}

// AsyncEnableDom0Access: Unhide a PCI device from the dom0 kernel. (Takes affect after next boot.)
// Version: closed
func (pci) AsyncEnableDom0Access(session *Session, self PCIRef) (retval TaskRef, err error) {
	method := "Async.PCI.enable_dom0_access"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// EnableDom0Access2: Unhide a PCI device from the dom0 kernel. (Takes affect after next boot.)
// Version: closed
func (pci) EnableDom0Access2(session *Session, self PCIRef) (retval PciDom0Access, err error) {
	method := "PCI.enable_dom0_access"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPciDom0Access(method+" -> ", result)
	return
}

// AsyncEnableDom0Access2: Unhide a PCI device from the dom0 kernel. (Takes affect after next boot.)
// Version: closed
func (pci) AsyncEnableDom0Access2(session *Session, self PCIRef) (retval TaskRef, err error) {
	method := "Async.PCI.enable_dom0_access"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableDom0Access: Hide a PCI device from the dom0 kernel. (Takes affect after next boot.)
// Version: closed
func (pci) DisableDom0Access(session *Session, self PCIRef) (retval PciDom0Access, err error) {
	method := "PCI.disable_dom0_access"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPciDom0Access(method+" -> ", result)
	return
}

// AsyncDisableDom0Access: Hide a PCI device from the dom0 kernel. (Takes affect after next boot.)
// Version: closed
func (pci) AsyncDisableDom0Access(session *Session, self PCIRef) (retval TaskRef, err error) {
	method := "Async.PCI.disable_dom0_access"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableDom0Access2: Hide a PCI device from the dom0 kernel. (Takes affect after next boot.)
// Version: closed
func (pci) DisableDom0Access2(session *Session, self PCIRef) (retval PciDom0Access, err error) {
	method := "PCI.disable_dom0_access"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPciDom0Access(method+" -> ", result)
	return
}

// AsyncDisableDom0Access2: Hide a PCI device from the dom0 kernel. (Takes affect after next boot.)
// Version: closed
func (pci) AsyncDisableDom0Access2(session *Session, self PCIRef) (retval TaskRef, err error) {
	method := "Async.PCI.disable_dom0_access"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given PCI.  If the key is not in that Map, then do nothing.
// Version: boston
func (pci) RemoveFromOtherConfig(session *Session, self PCIRef, key string) (err error) {
	method := "PCI.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given PCI.  If the key is not in that Map, then do nothing.
// Version: boston
func (pci) RemoveFromOtherConfig3(session *Session, self PCIRef, key string) (err error) {
	method := "PCI.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given PCI.
// Version: boston
func (pci) AddToOtherConfig(session *Session, self PCIRef, key string, value string) (err error) {
	method := "PCI.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	return
}

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given PCI.
// Version: boston
func (pci) AddToOtherConfig4(session *Session, self PCIRef, key string, value string) (err error) {
	method := "PCI.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	return
}

// SetOtherConfig: Set the other_config field of the given PCI.
// Version: boston
func (pci) SetOtherConfig(session *Session, self PCIRef, value map[string]string) (err error) {
	method := "PCI.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetOtherConfig3: Set the other_config field of the given PCI.
// Version: boston
func (pci) SetOtherConfig3(session *Session, self PCIRef, value map[string]string) (err error) {
	method := "PCI.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// GetDriverName: Get the driver_name field of the given PCI.
// Version: boston
func (pci) GetDriverName(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_driver_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDriverName2: Get the driver_name field of the given PCI.
// Version: boston
func (pci) GetDriverName2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_driver_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSubsystemDeviceName: Get the subsystem_device_name field of the given PCI.
// Version: boston
func (pci) GetSubsystemDeviceName(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_subsystem_device_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSubsystemDeviceName2: Get the subsystem_device_name field of the given PCI.
// Version: boston
func (pci) GetSubsystemDeviceName2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_subsystem_device_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSubsystemDeviceID: Get the subsystem_device_id field of the given PCI.
// Version: boston
func (pci) GetSubsystemDeviceID(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_subsystem_device_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSubsystemDeviceID2: Get the subsystem_device_id field of the given PCI.
// Version: boston
func (pci) GetSubsystemDeviceID2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_subsystem_device_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSubsystemVendorName: Get the subsystem_vendor_name field of the given PCI.
// Version: boston
func (pci) GetSubsystemVendorName(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_subsystem_vendor_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSubsystemVendorName2: Get the subsystem_vendor_name field of the given PCI.
// Version: boston
func (pci) GetSubsystemVendorName2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_subsystem_vendor_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSubsystemVendorID: Get the subsystem_vendor_id field of the given PCI.
// Version: boston
func (pci) GetSubsystemVendorID(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_subsystem_vendor_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSubsystemVendorID2: Get the subsystem_vendor_id field of the given PCI.
// Version: boston
func (pci) GetSubsystemVendorID2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_subsystem_vendor_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given PCI.
// Version: boston
func (pci) GetOtherConfig(session *Session, self PCIRef) (retval map[string]string, err error) {
	method := "PCI.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given PCI.
// Version: boston
func (pci) GetOtherConfig2(session *Session, self PCIRef) (retval map[string]string, err error) {
	method := "PCI.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDependencies: Get the dependencies field of the given PCI.
// Version: boston
func (pci) GetDependencies(session *Session, self PCIRef) (retval []PCIRef, err error) {
	method := "PCI.get_dependencies"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRefSet(method+" -> ", result)
	return
}

// GetDependencies2: Get the dependencies field of the given PCI.
// Version: boston
func (pci) GetDependencies2(session *Session, self PCIRef) (retval []PCIRef, err error) {
	method := "PCI.get_dependencies"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRefSet(method+" -> ", result)
	return
}

// GetPciID: Get the pci_id field of the given PCI.
// Version: boston
func (pci) GetPciID(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_pci_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPciID2: Get the pci_id field of the given PCI.
// Version: boston
func (pci) GetPciID2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_pci_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHost: Get the host field of the given PCI.
// Version: boston
func (pci) GetHost(session *Session, self PCIRef) (retval HostRef, err error) {
	method := "PCI.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRef(method+" -> ", result)
	return
}

// GetHost2: Get the host field of the given PCI.
// Version: boston
func (pci) GetHost2(session *Session, self PCIRef) (retval HostRef, err error) {
	method := "PCI.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRef(method+" -> ", result)
	return
}

// GetDeviceName: Get the device_name field of the given PCI.
// Version: boston
func (pci) GetDeviceName(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_device_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDeviceName2: Get the device_name field of the given PCI.
// Version: boston
func (pci) GetDeviceName2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_device_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDeviceID: Get the device_id field of the given PCI.
// Version: boston
func (pci) GetDeviceID(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_device_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDeviceID2: Get the device_id field of the given PCI.
// Version: boston
func (pci) GetDeviceID2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_device_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVendorName: Get the vendor_name field of the given PCI.
// Version: boston
func (pci) GetVendorName(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_vendor_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVendorName2: Get the vendor_name field of the given PCI.
// Version: boston
func (pci) GetVendorName2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_vendor_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVendorID: Get the vendor_id field of the given PCI.
// Version: boston
func (pci) GetVendorID(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_vendor_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVendorID2: Get the vendor_id field of the given PCI.
// Version: boston
func (pci) GetVendorID2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_vendor_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClassName: Get the class_name field of the given PCI.
// Version: boston
func (pci) GetClassName(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_class_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClassName2: Get the class_name field of the given PCI.
// Version: boston
func (pci) GetClassName2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_class_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClassID: Get the class_id field of the given PCI.
// Version: boston
func (pci) GetClassID(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_class_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClassID2: Get the class_id field of the given PCI.
// Version: boston
func (pci) GetClassID2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_class_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given PCI.
// Version: boston
func (pci) GetUUID(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given PCI.
// Version: boston
func (pci) GetUUID2(session *Session, self PCIRef) (retval string, err error) {
	method := "PCI.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the PCI instance with the specified UUID.
// Version: boston
func (pci) GetByUUID(session *Session, uuid string) (retval PCIRef, err error) {
	method := "PCI.get_by_uuid"
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
	retval, err = deserializePCIRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the PCI instance with the specified UUID.
// Version: boston
func (pci) GetByUUID2(session *Session, uuid string) (retval PCIRef, err error) {
	method := "PCI.get_by_uuid"
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
	retval, err = deserializePCIRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given PCI.
// Version: boston
func (pci) GetRecord(session *Session, self PCIRef) (retval PCIRecord, err error) {
	method := "PCI.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given PCI.
// Version: boston
func (pci) GetRecord2(session *Session, self PCIRef) (retval PCIRecord, err error) {
	method := "PCI.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePCIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRecord(method+" -> ", result)
	return
}
