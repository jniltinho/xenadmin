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

type BondRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// The bonded interface
	Master PIFRef `json:"master,omitempty"`
	// The interfaces which are part of this bond
	Slaves []PIFRef `json:"slaves,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// The PIF of which the IP configuration and MAC were copied to the bond, and which will receive all configuration/VLANs/VIFs on the bond if the bond is destroyed
	PrimarySlave PIFRef `json:"primary_slave,omitempty"`
	// The algorithm used to distribute traffic among the bonded NICs
	Mode BondMode `json:"mode,omitempty"`
	// Additional configuration properties specific to the bond mode.
	Properties map[string]string `json:"properties,omitempty"`
	// Number of links up in this bond
	LinksUp int `json:"links_up,omitempty"`
	// true if the MAC was taken from the primary slave when the bond was created, and false if the client specified the MAC
	AutoUpdateMac bool `json:"auto_update_mac,omitempty"`
}

type BondRef string

// A Network bond that combines physical network interfaces, also known as link aggregation
type bond struct{}

var Bond bond

// GetAllRecords: Return a map of Bond references to Bond records for all Bonds known to the system.
// Version: miami
func (bond) GetAllRecords(session *Session) (retval map[BondRef]BondRecord, err error) {
	method := "Bond.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRefToBondRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of Bond references to Bond records for all Bonds known to the system.
// Version: miami
func (bond) GetAllRecords1(session *Session) (retval map[BondRef]BondRecord, err error) {
	method := "Bond.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRefToBondRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the Bonds known to the system.
// Version: miami
func (bond) GetAll(session *Session) (retval []BondRef, err error) {
	method := "Bond.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the Bonds known to the system.
// Version: miami
func (bond) GetAll1(session *Session) (retval []BondRef, err error) {
	method := "Bond.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRefSet(method+" -> ", result)
	return
}

// SetProperty: Set the value of a property of the bond
// Version: tampa
func (bond) SetProperty(session *Session, self BondRef, name string, value string) (err error) {
	method := "Bond.set_property"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, nameArg, valueArg)
	return
}

// AsyncSetProperty: Set the value of a property of the bond
// Version: tampa
func (bond) AsyncSetProperty(session *Session, self BondRef, name string, value string) (retval TaskRef, err error) {
	method := "Async.Bond.set_property"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, nameArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetProperty4: Set the value of a property of the bond
// Version: tampa
func (bond) SetProperty4(session *Session, self BondRef, name string, value string) (err error) {
	method := "Bond.set_property"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, nameArg, valueArg)
	return
}

// AsyncSetProperty4: Set the value of a property of the bond
// Version: tampa
func (bond) AsyncSetProperty4(session *Session, self BondRef, name string, value string) (retval TaskRef, err error) {
	method := "Async.Bond.set_property"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, nameArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMode: Change the bond mode
// Version: boston
func (bond) SetMode(session *Session, self BondRef, value BondMode) (err error) {
	method := "Bond.set_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumBondMode(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetMode: Change the bond mode
// Version: boston
func (bond) AsyncSetMode(session *Session, self BondRef, value BondMode) (retval TaskRef, err error) {
	method := "Async.Bond.set_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumBondMode(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMode3: Change the bond mode
// Version: boston
func (bond) SetMode3(session *Session, self BondRef, value BondMode) (err error) {
	method := "Bond.set_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumBondMode(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetMode3: Change the bond mode
// Version: boston
func (bond) AsyncSetMode3(session *Session, self BondRef, value BondMode) (retval TaskRef, err error) {
	method := "Async.Bond.set_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumBondMode(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Destroy: Destroy an interface bond
// Version: miami
func (bond) Destroy(session *Session, self BondRef) (err error) {
	method := "Bond.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy an interface bond
// Version: miami
func (bond) AsyncDestroy(session *Session, self BondRef) (retval TaskRef, err error) {
	method := "Async.Bond.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy an interface bond
// Version: miami
func (bond) Destroy2(session *Session, self BondRef) (err error) {
	method := "Bond.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy an interface bond
// Version: miami
func (bond) AsyncDestroy2(session *Session, self BondRef) (retval TaskRef, err error) {
	method := "Async.Bond.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Create an interface bond
// Version: tampa
func (bond) Create(session *Session, network NetworkRef, members []PIFRef, mac string, mode BondMode, properties map[string]string) (retval BondRef, err error) {
	method := "Bond.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	membersArg, err := serializePIFRefSet(fmt.Sprintf("%s(%s)", method, "members"), members)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumBondMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	propertiesArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "properties"), properties)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg, membersArg, macArg, modeArg, propertiesArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRef(method+" -> ", result)
	return
}

// AsyncCreate: Create an interface bond
// Version: tampa
func (bond) AsyncCreate(session *Session, network NetworkRef, members []PIFRef, mac string, mode BondMode, properties map[string]string) (retval TaskRef, err error) {
	method := "Async.Bond.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	membersArg, err := serializePIFRefSet(fmt.Sprintf("%s(%s)", method, "members"), members)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumBondMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	propertiesArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "properties"), properties)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg, membersArg, macArg, modeArg, propertiesArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create6: Create an interface bond
// Version: tampa
func (bond) Create6(session *Session, network NetworkRef, members []PIFRef, mac string, mode BondMode, properties map[string]string) (retval BondRef, err error) {
	method := "Bond.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	membersArg, err := serializePIFRefSet(fmt.Sprintf("%s(%s)", method, "members"), members)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumBondMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	propertiesArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "properties"), properties)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg, membersArg, macArg, modeArg, propertiesArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRef(method+" -> ", result)
	return
}

// AsyncCreate6: Create an interface bond
// Version: tampa
func (bond) AsyncCreate6(session *Session, network NetworkRef, members []PIFRef, mac string, mode BondMode, properties map[string]string) (retval TaskRef, err error) {
	method := "Async.Bond.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	membersArg, err := serializePIFRefSet(fmt.Sprintf("%s(%s)", method, "members"), members)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumBondMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	propertiesArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "properties"), properties)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg, membersArg, macArg, modeArg, propertiesArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create5: Create an interface bond
// Version: boston
func (bond) Create5(session *Session, network NetworkRef, members []PIFRef, mac string, mode BondMode) (retval BondRef, err error) {
	method := "Bond.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	membersArg, err := serializePIFRefSet(fmt.Sprintf("%s(%s)", method, "members"), members)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumBondMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg, membersArg, macArg, modeArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRef(method+" -> ", result)
	return
}

// AsyncCreate5: Create an interface bond
// Version: boston
func (bond) AsyncCreate5(session *Session, network NetworkRef, members []PIFRef, mac string, mode BondMode) (retval TaskRef, err error) {
	method := "Async.Bond.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	membersArg, err := serializePIFRefSet(fmt.Sprintf("%s(%s)", method, "members"), members)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumBondMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg, membersArg, macArg, modeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create4: Create an interface bond
// Version: miami
func (bond) Create4(session *Session, network NetworkRef, members []PIFRef, mac string) (retval BondRef, err error) {
	method := "Bond.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	membersArg, err := serializePIFRefSet(fmt.Sprintf("%s(%s)", method, "members"), members)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg, membersArg, macArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRef(method+" -> ", result)
	return
}

// AsyncCreate4: Create an interface bond
// Version: miami
func (bond) AsyncCreate4(session *Session, network NetworkRef, members []PIFRef, mac string) (retval TaskRef, err error) {
	method := "Async.Bond.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	membersArg, err := serializePIFRefSet(fmt.Sprintf("%s(%s)", method, "members"), members)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg, membersArg, macArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given Bond.  If the key is not in that Map, then do nothing.
// Version: miami
func (bond) RemoveFromOtherConfig(session *Session, self BondRef, key string) (err error) {
	method := "Bond.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given Bond.  If the key is not in that Map, then do nothing.
// Version: miami
func (bond) RemoveFromOtherConfig3(session *Session, self BondRef, key string) (err error) {
	method := "Bond.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given Bond.
// Version: miami
func (bond) AddToOtherConfig(session *Session, self BondRef, key string, value string) (err error) {
	method := "Bond.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given Bond.
// Version: miami
func (bond) AddToOtherConfig4(session *Session, self BondRef, key string, value string) (err error) {
	method := "Bond.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig: Set the other_config field of the given Bond.
// Version: miami
func (bond) SetOtherConfig(session *Session, self BondRef, value map[string]string) (err error) {
	method := "Bond.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given Bond.
// Version: miami
func (bond) SetOtherConfig3(session *Session, self BondRef, value map[string]string) (err error) {
	method := "Bond.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAutoUpdateMac: Get the auto_update_mac field of the given Bond.
// Version: miami
func (bond) GetAutoUpdateMac(session *Session, self BondRef) (retval bool, err error) {
	method := "Bond.get_auto_update_mac"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAutoUpdateMac2: Get the auto_update_mac field of the given Bond.
// Version: miami
func (bond) GetAutoUpdateMac2(session *Session, self BondRef) (retval bool, err error) {
	method := "Bond.get_auto_update_mac"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLinksUp: Get the links_up field of the given Bond.
// Version: miami
func (bond) GetLinksUp(session *Session, self BondRef) (retval int, err error) {
	method := "Bond.get_links_up"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLinksUp2: Get the links_up field of the given Bond.
// Version: miami
func (bond) GetLinksUp2(session *Session, self BondRef) (retval int, err error) {
	method := "Bond.get_links_up"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetProperties: Get the properties field of the given Bond.
// Version: miami
func (bond) GetProperties(session *Session, self BondRef) (retval map[string]string, err error) {
	method := "Bond.get_properties"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetProperties2: Get the properties field of the given Bond.
// Version: miami
func (bond) GetProperties2(session *Session, self BondRef) (retval map[string]string, err error) {
	method := "Bond.get_properties"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMode: Get the mode field of the given Bond.
// Version: miami
func (bond) GetMode(session *Session, self BondRef) (retval BondMode, err error) {
	method := "Bond.get_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumBondMode(method+" -> ", result)
	return
}

// GetMode2: Get the mode field of the given Bond.
// Version: miami
func (bond) GetMode2(session *Session, self BondRef) (retval BondMode, err error) {
	method := "Bond.get_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumBondMode(method+" -> ", result)
	return
}

// GetPrimarySlave: Get the primary_slave field of the given Bond.
// Version: miami
func (bond) GetPrimarySlave(session *Session, self BondRef) (retval PIFRef, err error) {
	method := "Bond.get_primary_slave"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPrimarySlave2: Get the primary_slave field of the given Bond.
// Version: miami
func (bond) GetPrimarySlave2(session *Session, self BondRef) (retval PIFRef, err error) {
	method := "Bond.get_primary_slave"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given Bond.
// Version: miami
func (bond) GetOtherConfig(session *Session, self BondRef) (retval map[string]string, err error) {
	method := "Bond.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given Bond.
// Version: miami
func (bond) GetOtherConfig2(session *Session, self BondRef) (retval map[string]string, err error) {
	method := "Bond.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSlaves: Get the slaves field of the given Bond.
// Version: miami
func (bond) GetSlaves(session *Session, self BondRef) (retval []PIFRef, err error) {
	method := "Bond.get_slaves"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRefSet(method+" -> ", result)
	return
}

// GetSlaves2: Get the slaves field of the given Bond.
// Version: miami
func (bond) GetSlaves2(session *Session, self BondRef) (retval []PIFRef, err error) {
	method := "Bond.get_slaves"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRefSet(method+" -> ", result)
	return
}

// GetMaster: Get the master field of the given Bond.
// Version: miami
func (bond) GetMaster(session *Session, self BondRef) (retval PIFRef, err error) {
	method := "Bond.get_master"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMaster2: Get the master field of the given Bond.
// Version: miami
func (bond) GetMaster2(session *Session, self BondRef) (retval PIFRef, err error) {
	method := "Bond.get_master"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given Bond.
// Version: miami
func (bond) GetUUID(session *Session, self BondRef) (retval string, err error) {
	method := "Bond.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given Bond.
// Version: miami
func (bond) GetUUID2(session *Session, self BondRef) (retval string, err error) {
	method := "Bond.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the Bond instance with the specified UUID.
// Version: miami
func (bond) GetByUUID(session *Session, uuid string) (retval BondRef, err error) {
	method := "Bond.get_by_uuid"
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
	retval, err = deserializeBondRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the Bond instance with the specified UUID.
// Version: miami
func (bond) GetByUUID2(session *Session, uuid string) (retval BondRef, err error) {
	method := "Bond.get_by_uuid"
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
	retval, err = deserializeBondRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given Bond.
// Version: miami
func (bond) GetRecord(session *Session, self BondRef) (retval BondRecord, err error) {
	method := "Bond.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given Bond.
// Version: miami
func (bond) GetRecord2(session *Session, self BondRef) (retval BondRecord, err error) {
	method := "Bond.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRecord(method+" -> ", result)
	return
}
