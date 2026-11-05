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

type VMGuestMetricsRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// version of the OS
	OsVersion map[string]string `json:"os_version,omitempty"`
	// The NETBIOS name of the machine
	NetbiosName map[string]string `json:"netbios_name,omitempty"`
	// version of the PV drivers
	PVDriversVersion map[string]string `json:"PV_drivers_version,omitempty"`
	// Logically equivalent to PV_drivers_detected
	PVDriversUpToDate bool `json:"PV_drivers_up_to_date,omitempty"`
	// This field exists but has no data. Use the memory and memory_internal_free RRD data-sources instead.
	Memory map[string]string `json:"memory,omitempty"`
	// This field exists but has no data.
	Disks map[string]string `json:"disks,omitempty"`
	// network configuration
	Networks map[string]string `json:"networks,omitempty"`
	// anything else
	Other map[string]string `json:"other,omitempty"`
	// Time at which this information was last updated
	LastUpdated time.Time `json:"last_updated,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// True if the guest is sending heartbeat messages via the guest agent
	Live bool `json:"live,omitempty"`
	// The guest&apos;s statement of whether it supports VBD hotplug, i.e. whether it is capable of responding immediately to instantiation of a new VBD by bringing online a new PV block device. If the guest states that it is not capable, then the VBD plug and unplug operations will not be allowed while the guest is running.
	CanUseHotplugVbd TristateType `json:"can_use_hotplug_vbd,omitempty"`
	// The guest&apos;s statement of whether it supports VIF hotplug, i.e. whether it is capable of responding immediately to instantiation of a new VIF by bringing online a new PV network device. If the guest states that it is not capable, then the VIF plug and unplug operations will not be allowed while the guest is running.
	CanUseHotplugVif TristateType `json:"can_use_hotplug_vif,omitempty"`
	// At least one of the guest&apos;s devices has successfully connected to the backend.
	PVDriversDetected bool `json:"PV_drivers_detected,omitempty"`
	// The guest&apos;s services data.
	Services map[string]string `json:"services,omitempty"`
}

type VMGuestMetricsRef string

// The metrics reported by the guest (as opposed to inferred from outside)
type vmGuestMetrics struct{}

var VMGuestMetrics vmGuestMetrics

// GetAllRecords: Return a map of VM_guest_metrics references to VM_guest_metrics records for all VM_guest_metrics instances known to the system.
// Version: rio
func (vmGuestMetrics) GetAllRecords(session *Session) (retval map[VMGuestMetricsRef]VMGuestMetricsRecord, err error) {
	method := "VM_guest_metrics.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMGuestMetricsRefToVMGuestMetricsRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of VM_guest_metrics references to VM_guest_metrics records for all VM_guest_metrics instances known to the system.
// Version: rio
func (vmGuestMetrics) GetAllRecords1(session *Session) (retval map[VMGuestMetricsRef]VMGuestMetricsRecord, err error) {
	method := "VM_guest_metrics.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMGuestMetricsRefToVMGuestMetricsRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the VM_guest_metrics instances known to the system.
// Version: rio
func (vmGuestMetrics) GetAll(session *Session) (retval []VMGuestMetricsRef, err error) {
	method := "VM_guest_metrics.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMGuestMetricsRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the VM_guest_metrics instances known to the system.
// Version: rio
func (vmGuestMetrics) GetAll1(session *Session) (retval []VMGuestMetricsRef, err error) {
	method := "VM_guest_metrics.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMGuestMetricsRefSet(method+" -> ", result)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given VM_guest_metrics.  If the key is not in that Map, then do nothing.
// Version: orlando
func (vmGuestMetrics) RemoveFromOtherConfig(session *Session, self VMGuestMetricsRef, key string) (err error) {
	method := "VM_guest_metrics.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given VM_guest_metrics.  If the key is not in that Map, then do nothing.
// Version: orlando
func (vmGuestMetrics) RemoveFromOtherConfig3(session *Session, self VMGuestMetricsRef, key string) (err error) {
	method := "VM_guest_metrics.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig2: Remove the given key and its corresponding value from the other_config field of the given VM_guest_metrics.  If the key is not in that Map, then do nothing.
// Version: rio
func (vmGuestMetrics) RemoveFromOtherConfig2(session *Session, self VMGuestMetricsRef) (err error) {
	method := "VM_guest_metrics.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given VM_guest_metrics.
// Version: orlando
func (vmGuestMetrics) AddToOtherConfig(session *Session, self VMGuestMetricsRef, key string, value string) (err error) {
	method := "VM_guest_metrics.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given VM_guest_metrics.
// Version: orlando
func (vmGuestMetrics) AddToOtherConfig4(session *Session, self VMGuestMetricsRef, key string, value string) (err error) {
	method := "VM_guest_metrics.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig2: Add the given key-value pair to the other_config field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) AddToOtherConfig2(session *Session, self VMGuestMetricsRef) (err error) {
	method := "VM_guest_metrics.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetOtherConfig: Set the other_config field of the given VM_guest_metrics.
// Version: orlando
func (vmGuestMetrics) SetOtherConfig(session *Session, self VMGuestMetricsRef, value map[string]string) (err error) {
	method := "VM_guest_metrics.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given VM_guest_metrics.
// Version: orlando
func (vmGuestMetrics) SetOtherConfig3(session *Session, self VMGuestMetricsRef, value map[string]string) (err error) {
	method := "VM_guest_metrics.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig2: Set the other_config field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) SetOtherConfig2(session *Session, self VMGuestMetricsRef) (err error) {
	method := "VM_guest_metrics.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// GetServices: Get the services field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetServices(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_services"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetServices2: Get the services field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetServices2(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_services"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVDriversDetected: Get the PV_drivers_detected field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetPVDriversDetected(session *Session, self VMGuestMetricsRef) (retval bool, err error) {
	method := "VM_guest_metrics.get_PV_drivers_detected"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVDriversDetected2: Get the PV_drivers_detected field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetPVDriversDetected2(session *Session, self VMGuestMetricsRef) (retval bool, err error) {
	method := "VM_guest_metrics.get_PV_drivers_detected"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCanUseHotplugVif: Get the can_use_hotplug_vif field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetCanUseHotplugVif(session *Session, self VMGuestMetricsRef) (retval TristateType, err error) {
	method := "VM_guest_metrics.get_can_use_hotplug_vif"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumTristateType(method+" -> ", result)
	return
}

// GetCanUseHotplugVif2: Get the can_use_hotplug_vif field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetCanUseHotplugVif2(session *Session, self VMGuestMetricsRef) (retval TristateType, err error) {
	method := "VM_guest_metrics.get_can_use_hotplug_vif"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumTristateType(method+" -> ", result)
	return
}

// GetCanUseHotplugVbd: Get the can_use_hotplug_vbd field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetCanUseHotplugVbd(session *Session, self VMGuestMetricsRef) (retval TristateType, err error) {
	method := "VM_guest_metrics.get_can_use_hotplug_vbd"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumTristateType(method+" -> ", result)
	return
}

// GetCanUseHotplugVbd2: Get the can_use_hotplug_vbd field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetCanUseHotplugVbd2(session *Session, self VMGuestMetricsRef) (retval TristateType, err error) {
	method := "VM_guest_metrics.get_can_use_hotplug_vbd"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumTristateType(method+" -> ", result)
	return
}

// GetLive: Get the live field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetLive(session *Session, self VMGuestMetricsRef) (retval bool, err error) {
	method := "VM_guest_metrics.get_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLive2: Get the live field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetLive2(session *Session, self VMGuestMetricsRef) (retval bool, err error) {
	method := "VM_guest_metrics.get_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetOtherConfig(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetOtherConfig2(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastUpdated: Get the last_updated field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetLastUpdated(session *Session, self VMGuestMetricsRef) (retval time.Time, err error) {
	method := "VM_guest_metrics.get_last_updated"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastUpdated2: Get the last_updated field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetLastUpdated2(session *Session, self VMGuestMetricsRef) (retval time.Time, err error) {
	method := "VM_guest_metrics.get_last_updated"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOther: Get the other field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetOther(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_other"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOther2: Get the other field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetOther2(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_other"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNetworks: Get the networks field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetNetworks(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_networks"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNetworks2: Get the networks field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetNetworks2(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_networks"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDisks: Get the disks field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetDisks(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_disks"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDisks2: Get the disks field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetDisks2(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_disks"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemory: Get the memory field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetMemory(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemory2: Get the memory field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetMemory2(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVDriversUpToDate: Get the PV_drivers_up_to_date field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetPVDriversUpToDate(session *Session, self VMGuestMetricsRef) (retval bool, err error) {
	method := "VM_guest_metrics.get_PV_drivers_up_to_date"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVDriversUpToDate2: Get the PV_drivers_up_to_date field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetPVDriversUpToDate2(session *Session, self VMGuestMetricsRef) (retval bool, err error) {
	method := "VM_guest_metrics.get_PV_drivers_up_to_date"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVDriversVersion: Get the PV_drivers_version field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetPVDriversVersion(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_PV_drivers_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVDriversVersion2: Get the PV_drivers_version field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetPVDriversVersion2(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_PV_drivers_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNetbiosName: Get the netbios_name field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetNetbiosName(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_netbios_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNetbiosName2: Get the netbios_name field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetNetbiosName2(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_netbios_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOsVersion: Get the os_version field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetOsVersion(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_os_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOsVersion2: Get the os_version field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetOsVersion2(session *Session, self VMGuestMetricsRef) (retval map[string]string, err error) {
	method := "VM_guest_metrics.get_os_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetUUID(session *Session, self VMGuestMetricsRef) (retval string, err error) {
	method := "VM_guest_metrics.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetUUID2(session *Session, self VMGuestMetricsRef) (retval string, err error) {
	method := "VM_guest_metrics.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the VM_guest_metrics instance with the specified UUID.
// Version: rio
func (vmGuestMetrics) GetByUUID(session *Session, uuid string) (retval VMGuestMetricsRef, err error) {
	method := "VM_guest_metrics.get_by_uuid"
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
	retval, err = deserializeVMGuestMetricsRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the VM_guest_metrics instance with the specified UUID.
// Version: rio
func (vmGuestMetrics) GetByUUID2(session *Session, uuid string) (retval VMGuestMetricsRef, err error) {
	method := "VM_guest_metrics.get_by_uuid"
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
	retval, err = deserializeVMGuestMetricsRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetRecord(session *Session, self VMGuestMetricsRef) (retval VMGuestMetricsRecord, err error) {
	method := "VM_guest_metrics.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMGuestMetricsRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given VM_guest_metrics.
// Version: rio
func (vmGuestMetrics) GetRecord2(session *Session, self VMGuestMetricsRef) (retval VMGuestMetricsRecord, err error) {
	method := "VM_guest_metrics.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMGuestMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMGuestMetricsRecord(method+" -> ", result)
	return
}
