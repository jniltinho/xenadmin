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
	"time"
)

type PIFMetricsRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// Read bandwidth (KiB/s)
	IoReadKbs float64 `json:"io_read_kbs,omitempty"`
	// Write bandwidth (KiB/s)
	IoWriteKbs float64 `json:"io_write_kbs,omitempty"`
	// Report if the PIF got a carrier or not
	Carrier bool `json:"carrier,omitempty"`
	// Report vendor ID
	VendorID string `json:"vendor_id,omitempty"`
	// Report vendor name
	VendorName string `json:"vendor_name,omitempty"`
	// Report device ID
	DeviceID string `json:"device_id,omitempty"`
	// Report device name
	DeviceName string `json:"device_name,omitempty"`
	// Speed of the link in Mbit/s (if available)
	Speed int `json:"speed,omitempty"`
	// Full duplex capability of the link (if available)
	Duplex bool `json:"duplex,omitempty"`
	// PCI bus path of the pif (if available)
	PciBusPath string `json:"pci_bus_path,omitempty"`
	// Time at which this information was last updated
	LastUpdated time.Time `json:"last_updated,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
}

type PIFMetricsRef string

// The metrics associated with a physical network interface
type pifMetrics struct{}

var PIFMetrics pifMetrics

// GetAllRecords: Return a map of PIF_metrics references to PIF_metrics records for all PIF_metrics instances known to the system.
// Version: rio
func (pifMetrics) GetAllRecords(session *Session) (retval map[PIFMetricsRef]PIFMetricsRecord, err error) {
	method := "PIF_metrics.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFMetricsRefToPIFMetricsRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of PIF_metrics references to PIF_metrics records for all PIF_metrics instances known to the system.
// Version: rio
func (pifMetrics) GetAllRecords1(session *Session) (retval map[PIFMetricsRef]PIFMetricsRecord, err error) {
	method := "PIF_metrics.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFMetricsRefToPIFMetricsRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the PIF_metrics instances known to the system.
// Version: rio
func (pifMetrics) GetAll(session *Session) (retval []PIFMetricsRef, err error) {
	method := "PIF_metrics.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFMetricsRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the PIF_metrics instances known to the system.
// Version: rio
func (pifMetrics) GetAll1(session *Session) (retval []PIFMetricsRef, err error) {
	method := "PIF_metrics.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFMetricsRefSet(method+" -> ", result)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given PIF_metrics.  If the key is not in that Map, then do nothing.
// Version: orlando
func (pifMetrics) RemoveFromOtherConfig(session *Session, self PIFMetricsRef, key string) (err error) {
	method := "PIF_metrics.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given PIF_metrics.  If the key is not in that Map, then do nothing.
// Version: orlando
func (pifMetrics) RemoveFromOtherConfig3(session *Session, self PIFMetricsRef, key string) (err error) {
	method := "PIF_metrics.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig2: Remove the given key and its corresponding value from the other_config field of the given PIF_metrics.  If the key is not in that Map, then do nothing.
// Version: rio
func (pifMetrics) RemoveFromOtherConfig2(session *Session, self PIFMetricsRef) (err error) {
	method := "PIF_metrics.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given PIF_metrics.
// Version: orlando
func (pifMetrics) AddToOtherConfig(session *Session, self PIFMetricsRef, key string, value string) (err error) {
	method := "PIF_metrics.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given PIF_metrics.
// Version: orlando
func (pifMetrics) AddToOtherConfig4(session *Session, self PIFMetricsRef, key string, value string) (err error) {
	method := "PIF_metrics.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig2: Add the given key-value pair to the other_config field of the given PIF_metrics.
// Version: rio
func (pifMetrics) AddToOtherConfig2(session *Session, self PIFMetricsRef) (err error) {
	method := "PIF_metrics.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetOtherConfig: Set the other_config field of the given PIF_metrics.
// Version: orlando
func (pifMetrics) SetOtherConfig(session *Session, self PIFMetricsRef, value map[string]string) (err error) {
	method := "PIF_metrics.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given PIF_metrics.
// Version: orlando
func (pifMetrics) SetOtherConfig3(session *Session, self PIFMetricsRef, value map[string]string) (err error) {
	method := "PIF_metrics.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig2: Set the other_config field of the given PIF_metrics.
// Version: rio
func (pifMetrics) SetOtherConfig2(session *Session, self PIFMetricsRef) (err error) {
	method := "PIF_metrics.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// GetOtherConfig: Get the other_config field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetOtherConfig(session *Session, self PIFMetricsRef) (retval map[string]string, err error) {
	method := "PIF_metrics.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetOtherConfig2(session *Session, self PIFMetricsRef) (retval map[string]string, err error) {
	method := "PIF_metrics.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastUpdated: Get the last_updated field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetLastUpdated(session *Session, self PIFMetricsRef) (retval time.Time, err error) {
	method := "PIF_metrics.get_last_updated"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetLastUpdated2: Get the last_updated field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetLastUpdated2(session *Session, self PIFMetricsRef) (retval time.Time, err error) {
	method := "PIF_metrics.get_last_updated"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTime(method+" -> ", result)
	return
}

// GetPciBusPath: Get the pci_bus_path field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetPciBusPath(session *Session, self PIFMetricsRef) (retval string, err error) {
	method := "PIF_metrics.get_pci_bus_path"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPciBusPath2: Get the pci_bus_path field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetPciBusPath2(session *Session, self PIFMetricsRef) (retval string, err error) {
	method := "PIF_metrics.get_pci_bus_path"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDuplex: Get the duplex field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetDuplex(session *Session, self PIFMetricsRef) (retval bool, err error) {
	method := "PIF_metrics.get_duplex"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDuplex2: Get the duplex field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetDuplex2(session *Session, self PIFMetricsRef) (retval bool, err error) {
	method := "PIF_metrics.get_duplex"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSpeed: Get the speed field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetSpeed(session *Session, self PIFMetricsRef) (retval int, err error) {
	method := "PIF_metrics.get_speed"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSpeed2: Get the speed field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetSpeed2(session *Session, self PIFMetricsRef) (retval int, err error) {
	method := "PIF_metrics.get_speed"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDeviceName: Get the device_name field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetDeviceName(session *Session, self PIFMetricsRef) (retval string, err error) {
	method := "PIF_metrics.get_device_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDeviceName2: Get the device_name field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetDeviceName2(session *Session, self PIFMetricsRef) (retval string, err error) {
	method := "PIF_metrics.get_device_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDeviceID: Get the device_id field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetDeviceID(session *Session, self PIFMetricsRef) (retval string, err error) {
	method := "PIF_metrics.get_device_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDeviceID2: Get the device_id field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetDeviceID2(session *Session, self PIFMetricsRef) (retval string, err error) {
	method := "PIF_metrics.get_device_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVendorName: Get the vendor_name field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetVendorName(session *Session, self PIFMetricsRef) (retval string, err error) {
	method := "PIF_metrics.get_vendor_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVendorName2: Get the vendor_name field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetVendorName2(session *Session, self PIFMetricsRef) (retval string, err error) {
	method := "PIF_metrics.get_vendor_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVendorID: Get the vendor_id field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetVendorID(session *Session, self PIFMetricsRef) (retval string, err error) {
	method := "PIF_metrics.get_vendor_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVendorID2: Get the vendor_id field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetVendorID2(session *Session, self PIFMetricsRef) (retval string, err error) {
	method := "PIF_metrics.get_vendor_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCarrier: Get the carrier field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetCarrier(session *Session, self PIFMetricsRef) (retval bool, err error) {
	method := "PIF_metrics.get_carrier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCarrier2: Get the carrier field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetCarrier2(session *Session, self PIFMetricsRef) (retval bool, err error) {
	method := "PIF_metrics.get_carrier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIoWriteKbs: Get the io/write_kbs field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetIoWriteKbs(session *Session, self PIFMetricsRef) (retval float64, err error) {
	method := "PIF_metrics.get_io_write_kbs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIoWriteKbs2: Get the io/write_kbs field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetIoWriteKbs2(session *Session, self PIFMetricsRef) (retval float64, err error) {
	method := "PIF_metrics.get_io_write_kbs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIoReadKbs: Get the io/read_kbs field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetIoReadKbs(session *Session, self PIFMetricsRef) (retval float64, err error) {
	method := "PIF_metrics.get_io_read_kbs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIoReadKbs2: Get the io/read_kbs field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetIoReadKbs2(session *Session, self PIFMetricsRef) (retval float64, err error) {
	method := "PIF_metrics.get_io_read_kbs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetUUID(session *Session, self PIFMetricsRef) (retval string, err error) {
	method := "PIF_metrics.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetUUID2(session *Session, self PIFMetricsRef) (retval string, err error) {
	method := "PIF_metrics.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the PIF_metrics instance with the specified UUID.
// Version: rio
func (pifMetrics) GetByUUID(session *Session, uuid string) (retval PIFMetricsRef, err error) {
	method := "PIF_metrics.get_by_uuid"
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
	retval, err = deserializePIFMetricsRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the PIF_metrics instance with the specified UUID.
// Version: rio
func (pifMetrics) GetByUUID2(session *Session, uuid string) (retval PIFMetricsRef, err error) {
	method := "PIF_metrics.get_by_uuid"
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
	retval, err = deserializePIFMetricsRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetRecord(session *Session, self PIFMetricsRef) (retval PIFMetricsRecord, err error) {
	method := "PIF_metrics.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFMetricsRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given PIF_metrics.
// Version: rio
func (pifMetrics) GetRecord2(session *Session, self PIFMetricsRef) (retval PIFMetricsRecord, err error) {
	method := "PIF_metrics.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFMetricsRecord(method+" -> ", result)
	return
}
