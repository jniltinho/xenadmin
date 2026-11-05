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

type HostDriverRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// Host where this driver is installed
	Host HostRef `json:"host,omitempty"`
	// Name identifying the driver uniquely
	Name string `json:"name,omitempty"`
	// Descriptive name, not used for identification
	FriendlyName string `json:"friendly_name,omitempty"`
	// Variants of this driver available for selection
	Variants []DriverVariantRef `json:"variants,omitempty"`
	// Currently active variant of this driver, if any, or Null.
	ActiveVariant DriverVariantRef `json:"active_variant,omitempty"`
	// Variant (if any) selected to become active after reboot. Or Null
	SelectedVariant DriverVariantRef `json:"selected_variant,omitempty"`
	// Device type this driver supports, like network or storage
	Type string `json:"type,omitempty"`
	// Description of the driver
	Description string `json:"description,omitempty"`
	// Information about the driver
	Info string `json:"info,omitempty"`
}

type HostDriverRef string

// UNSUPPORTED. A multi-version driver on a host
type hostDriver struct{}

var HostDriver hostDriver

// GetAllRecords: Return a map of Host_driver references to Host_driver records for all Host_drivers known to the system.
// Version: closed
func (hostDriver) GetAllRecords(session *Session) (retval map[HostDriverRef]HostDriverRecord, err error) {
	method := "Host_driver.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostDriverRefToHostDriverRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of Host_driver references to Host_driver records for all Host_drivers known to the system.
// Version: closed
func (hostDriver) GetAllRecords1(session *Session) (retval map[HostDriverRef]HostDriverRecord, err error) {
	method := "Host_driver.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostDriverRefToHostDriverRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the Host_drivers known to the system.
// Version: closed
func (hostDriver) GetAll(session *Session) (retval []HostDriverRef, err error) {
	method := "Host_driver.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostDriverRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the Host_drivers known to the system.
// Version: closed
func (hostDriver) GetAll1(session *Session) (retval []HostDriverRef, err error) {
	method := "Host_driver.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostDriverRefSet(method+" -> ", result)
	return
}

// Rescan: UNSUPPORTED. Re-scan a host&apos;s drivers and update information about them. This is mostly  for trouble shooting.
// Version: closed
func (hostDriver) Rescan(session *Session, host HostRef) (err error) {
	method := "Host_driver.rescan"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncRescan: UNSUPPORTED. Re-scan a host&apos;s drivers and update information about them. This is mostly  for trouble shooting.
// Version: closed
func (hostDriver) AsyncRescan(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.Host_driver.rescan"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Rescan2: UNSUPPORTED. Re-scan a host&apos;s drivers and update information about them. This is mostly  for trouble shooting.
// Version: closed
func (hostDriver) Rescan2(session *Session, host HostRef) (err error) {
	method := "Host_driver.rescan"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, hostArg)
	return
}

// AsyncRescan2: UNSUPPORTED. Re-scan a host&apos;s drivers and update information about them. This is mostly  for trouble shooting.
// Version: closed
func (hostDriver) AsyncRescan2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.Host_driver.rescan"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Deselect: UNSUPPORTED. Deselect the currently active variant of this driver after reboot. No action will be taken if no variant is currently active.
// Version: closed
func (hostDriver) Deselect(session *Session, self HostDriverRef) (err error) {
	method := "Host_driver.deselect"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDeselect: UNSUPPORTED. Deselect the currently active variant of this driver after reboot. No action will be taken if no variant is currently active.
// Version: closed
func (hostDriver) AsyncDeselect(session *Session, self HostDriverRef) (retval TaskRef, err error) {
	method := "Async.Host_driver.deselect"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Deselect2: UNSUPPORTED. Deselect the currently active variant of this driver after reboot. No action will be taken if no variant is currently active.
// Version: closed
func (hostDriver) Deselect2(session *Session, self HostDriverRef) (err error) {
	method := "Host_driver.deselect"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDeselect2: UNSUPPORTED. Deselect the currently active variant of this driver after reboot. No action will be taken if no variant is currently active.
// Version: closed
func (hostDriver) AsyncDeselect2(session *Session, self HostDriverRef) (retval TaskRef, err error) {
	method := "Async.Host_driver.deselect"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Select: UNSUPPORTED. Select a variant of the driver to become active after reboot or immediately if currently no version is active
// Version: closed
func (hostDriver) Select(session *Session, self HostDriverRef, variant DriverVariantRef) (err error) {
	method := "Host_driver.select"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	variantArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "variant"), variant)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, variantArg)
	return
}

// AsyncSelect: UNSUPPORTED. Select a variant of the driver to become active after reboot or immediately if currently no version is active
// Version: closed
func (hostDriver) AsyncSelect(session *Session, self HostDriverRef, variant DriverVariantRef) (retval TaskRef, err error) {
	method := "Async.Host_driver.select"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	variantArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "variant"), variant)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, variantArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Select3: UNSUPPORTED. Select a variant of the driver to become active after reboot or immediately if currently no version is active
// Version: closed
func (hostDriver) Select3(session *Session, self HostDriverRef, variant DriverVariantRef) (err error) {
	method := "Host_driver.select"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	variantArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "variant"), variant)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, variantArg)
	return
}

// AsyncSelect3: UNSUPPORTED. Select a variant of the driver to become active after reboot or immediately if currently no version is active
// Version: closed
func (hostDriver) AsyncSelect3(session *Session, self HostDriverRef, variant DriverVariantRef) (retval TaskRef, err error) {
	method := "Async.Host_driver.select"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	variantArg, err := serializeDriverVariantRef(fmt.Sprintf("%s(%s)", method, "variant"), variant)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, variantArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetInfo: Get the info field of the given Host_driver.
// Version: closed
func (hostDriver) GetInfo(session *Session, self HostDriverRef) (retval string, err error) {
	method := "Host_driver.get_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetInfo2: Get the info field of the given Host_driver.
// Version: closed
func (hostDriver) GetInfo2(session *Session, self HostDriverRef) (retval string, err error) {
	method := "Host_driver.get_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDescription: Get the description field of the given Host_driver.
// Version: closed
func (hostDriver) GetDescription(session *Session, self HostDriverRef) (retval string, err error) {
	method := "Host_driver.get_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDescription2: Get the description field of the given Host_driver.
// Version: closed
func (hostDriver) GetDescription2(session *Session, self HostDriverRef) (retval string, err error) {
	method := "Host_driver.get_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetType: Get the type field of the given Host_driver.
// Version: closed
func (hostDriver) GetType(session *Session, self HostDriverRef) (retval string, err error) {
	method := "Host_driver.get_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetType2: Get the type field of the given Host_driver.
// Version: closed
func (hostDriver) GetType2(session *Session, self HostDriverRef) (retval string, err error) {
	method := "Host_driver.get_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSelectedVariant: Get the selected_variant field of the given Host_driver.
// Version: closed
func (hostDriver) GetSelectedVariant(session *Session, self HostDriverRef) (retval DriverVariantRef, err error) {
	method := "Host_driver.get_selected_variant"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDriverVariantRef(method+" -> ", result)
	return
}

// GetSelectedVariant2: Get the selected_variant field of the given Host_driver.
// Version: closed
func (hostDriver) GetSelectedVariant2(session *Session, self HostDriverRef) (retval DriverVariantRef, err error) {
	method := "Host_driver.get_selected_variant"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDriverVariantRef(method+" -> ", result)
	return
}

// GetActiveVariant: Get the active_variant field of the given Host_driver.
// Version: closed
func (hostDriver) GetActiveVariant(session *Session, self HostDriverRef) (retval DriverVariantRef, err error) {
	method := "Host_driver.get_active_variant"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDriverVariantRef(method+" -> ", result)
	return
}

// GetActiveVariant2: Get the active_variant field of the given Host_driver.
// Version: closed
func (hostDriver) GetActiveVariant2(session *Session, self HostDriverRef) (retval DriverVariantRef, err error) {
	method := "Host_driver.get_active_variant"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDriverVariantRef(method+" -> ", result)
	return
}

// GetVariants: Get the variants field of the given Host_driver.
// Version: closed
func (hostDriver) GetVariants(session *Session, self HostDriverRef) (retval []DriverVariantRef, err error) {
	method := "Host_driver.get_variants"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDriverVariantRefSet(method+" -> ", result)
	return
}

// GetVariants2: Get the variants field of the given Host_driver.
// Version: closed
func (hostDriver) GetVariants2(session *Session, self HostDriverRef) (retval []DriverVariantRef, err error) {
	method := "Host_driver.get_variants"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDriverVariantRefSet(method+" -> ", result)
	return
}

// GetFriendlyName: Get the friendly_name field of the given Host_driver.
// Version: closed
func (hostDriver) GetFriendlyName(session *Session, self HostDriverRef) (retval string, err error) {
	method := "Host_driver.get_friendly_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetFriendlyName2: Get the friendly_name field of the given Host_driver.
// Version: closed
func (hostDriver) GetFriendlyName2(session *Session, self HostDriverRef) (retval string, err error) {
	method := "Host_driver.get_friendly_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetName: Get the name field of the given Host_driver.
// Version: closed
func (hostDriver) GetName(session *Session, self HostDriverRef) (retval string, err error) {
	method := "Host_driver.get_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetName2: Get the name field of the given Host_driver.
// Version: closed
func (hostDriver) GetName2(session *Session, self HostDriverRef) (retval string, err error) {
	method := "Host_driver.get_name"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHost: Get the host field of the given Host_driver.
// Version: closed
func (hostDriver) GetHost(session *Session, self HostDriverRef) (retval HostRef, err error) {
	method := "Host_driver.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHost2: Get the host field of the given Host_driver.
// Version: closed
func (hostDriver) GetHost2(session *Session, self HostDriverRef) (retval HostRef, err error) {
	method := "Host_driver.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given Host_driver.
// Version: closed
func (hostDriver) GetUUID(session *Session, self HostDriverRef) (retval string, err error) {
	method := "Host_driver.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given Host_driver.
// Version: closed
func (hostDriver) GetUUID2(session *Session, self HostDriverRef) (retval string, err error) {
	method := "Host_driver.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the Host_driver instance with the specified UUID.
// Version: closed
func (hostDriver) GetByUUID(session *Session, uuid string) (retval HostDriverRef, err error) {
	method := "Host_driver.get_by_uuid"
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
	retval, err = deserializeHostDriverRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the Host_driver instance with the specified UUID.
// Version: closed
func (hostDriver) GetByUUID2(session *Session, uuid string) (retval HostDriverRef, err error) {
	method := "Host_driver.get_by_uuid"
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
	retval, err = deserializeHostDriverRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given Host_driver.
// Version: closed
func (hostDriver) GetRecord(session *Session, self HostDriverRef) (retval HostDriverRecord, err error) {
	method := "Host_driver.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostDriverRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given Host_driver.
// Version: closed
func (hostDriver) GetRecord2(session *Session, self HostDriverRef) (retval HostDriverRecord, err error) {
	method := "Host_driver.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostDriverRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostDriverRecord(method+" -> ", result)
	return
}
