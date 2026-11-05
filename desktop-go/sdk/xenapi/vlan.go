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

type VLANRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// interface on which traffic is tagged
	TaggedPIF PIFRef `json:"tagged_PIF,omitempty"`
	// interface on which traffic is untagged
	UntaggedPIF PIFRef `json:"untagged_PIF,omitempty"`
	// VLAN tag in use
	Tag int `json:"tag,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
}

type VLANRef string

// A VLAN mux/demux
type vlan struct{}

var VLAN vlan

// GetAllRecords: Return a map of VLAN references to VLAN records for all VLANs known to the system.
// Version: miami
func (vlan) GetAllRecords(session *Session) (retval map[VLANRef]VLANRecord, err error) {
	method := "VLAN.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVLANRefToVLANRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of VLAN references to VLAN records for all VLANs known to the system.
// Version: miami
func (vlan) GetAllRecords1(session *Session) (retval map[VLANRef]VLANRecord, err error) {
	method := "VLAN.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVLANRefToVLANRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the VLANs known to the system.
// Version: miami
func (vlan) GetAll(session *Session) (retval []VLANRef, err error) {
	method := "VLAN.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVLANRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the VLANs known to the system.
// Version: miami
func (vlan) GetAll1(session *Session) (retval []VLANRef, err error) {
	method := "VLAN.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVLANRefSet(method+" -> ", result)
	return
}

// Destroy: Destroy a VLAN mux/demuxer
// Version: miami
func (vlan) Destroy(session *Session, self VLANRef) (err error) {
	method := "VLAN.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy a VLAN mux/demuxer
// Version: miami
func (vlan) AsyncDestroy(session *Session, self VLANRef) (retval TaskRef, err error) {
	method := "Async.VLAN.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy a VLAN mux/demuxer
// Version: miami
func (vlan) Destroy2(session *Session, self VLANRef) (err error) {
	method := "VLAN.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy a VLAN mux/demuxer
// Version: miami
func (vlan) AsyncDestroy2(session *Session, self VLANRef) (retval TaskRef, err error) {
	method := "Async.VLAN.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Create a VLAN mux/demuxer
// Version: miami
func (vlan) Create(session *Session, taggedPIF PIFRef, tag int, network NetworkRef) (retval VLANRef, err error) {
	method := "VLAN.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	taggedPIFArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "tagged_PIF"), taggedPIF)
	if err != nil {
		return
	}
	tagArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "tag"), tag)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, taggedPIFArg, tagArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeVLANRef(method+" -> ", result)
	return
}

// AsyncCreate: Create a VLAN mux/demuxer
// Version: miami
func (vlan) AsyncCreate(session *Session, taggedPIF PIFRef, tag int, network NetworkRef) (retval TaskRef, err error) {
	method := "Async.VLAN.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	taggedPIFArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "tagged_PIF"), taggedPIF)
	if err != nil {
		return
	}
	tagArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "tag"), tag)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, taggedPIFArg, tagArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create4: Create a VLAN mux/demuxer
// Version: miami
func (vlan) Create4(session *Session, taggedPIF PIFRef, tag int, network NetworkRef) (retval VLANRef, err error) {
	method := "VLAN.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	taggedPIFArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "tagged_PIF"), taggedPIF)
	if err != nil {
		return
	}
	tagArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "tag"), tag)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, taggedPIFArg, tagArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeVLANRef(method+" -> ", result)
	return
}

// AsyncCreate4: Create a VLAN mux/demuxer
// Version: miami
func (vlan) AsyncCreate4(session *Session, taggedPIF PIFRef, tag int, network NetworkRef) (retval TaskRef, err error) {
	method := "Async.VLAN.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	taggedPIFArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "tagged_PIF"), taggedPIF)
	if err != nil {
		return
	}
	tagArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "tag"), tag)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, taggedPIFArg, tagArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given VLAN.  If the key is not in that Map, then do nothing.
// Version: miami
func (vlan) RemoveFromOtherConfig(session *Session, self VLANRef, key string) (err error) {
	method := "VLAN.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given VLAN.  If the key is not in that Map, then do nothing.
// Version: miami
func (vlan) RemoveFromOtherConfig3(session *Session, self VLANRef, key string) (err error) {
	method := "VLAN.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given VLAN.
// Version: miami
func (vlan) AddToOtherConfig(session *Session, self VLANRef, key string, value string) (err error) {
	method := "VLAN.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given VLAN.
// Version: miami
func (vlan) AddToOtherConfig4(session *Session, self VLANRef, key string, value string) (err error) {
	method := "VLAN.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig: Set the other_config field of the given VLAN.
// Version: miami
func (vlan) SetOtherConfig(session *Session, self VLANRef, value map[string]string) (err error) {
	method := "VLAN.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given VLAN.
// Version: miami
func (vlan) SetOtherConfig3(session *Session, self VLANRef, value map[string]string) (err error) {
	method := "VLAN.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given VLAN.
// Version: miami
func (vlan) GetOtherConfig(session *Session, self VLANRef) (retval map[string]string, err error) {
	method := "VLAN.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given VLAN.
// Version: miami
func (vlan) GetOtherConfig2(session *Session, self VLANRef) (retval map[string]string, err error) {
	method := "VLAN.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTag: Get the tag field of the given VLAN.
// Version: miami
func (vlan) GetTag(session *Session, self VLANRef) (retval int, err error) {
	method := "VLAN.get_tag"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTag2: Get the tag field of the given VLAN.
// Version: miami
func (vlan) GetTag2(session *Session, self VLANRef) (retval int, err error) {
	method := "VLAN.get_tag"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUntaggedPIF: Get the untagged_PIF field of the given VLAN.
// Version: miami
func (vlan) GetUntaggedPIF(session *Session, self VLANRef) (retval PIFRef, err error) {
	method := "VLAN.get_untagged_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUntaggedPIF2: Get the untagged_PIF field of the given VLAN.
// Version: miami
func (vlan) GetUntaggedPIF2(session *Session, self VLANRef) (retval PIFRef, err error) {
	method := "VLAN.get_untagged_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTaggedPIF: Get the tagged_PIF field of the given VLAN.
// Version: miami
func (vlan) GetTaggedPIF(session *Session, self VLANRef) (retval PIFRef, err error) {
	method := "VLAN.get_tagged_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTaggedPIF2: Get the tagged_PIF field of the given VLAN.
// Version: miami
func (vlan) GetTaggedPIF2(session *Session, self VLANRef) (retval PIFRef, err error) {
	method := "VLAN.get_tagged_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given VLAN.
// Version: miami
func (vlan) GetUUID(session *Session, self VLANRef) (retval string, err error) {
	method := "VLAN.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given VLAN.
// Version: miami
func (vlan) GetUUID2(session *Session, self VLANRef) (retval string, err error) {
	method := "VLAN.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the VLAN instance with the specified UUID.
// Version: miami
func (vlan) GetByUUID(session *Session, uuid string) (retval VLANRef, err error) {
	method := "VLAN.get_by_uuid"
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
	retval, err = deserializeVLANRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the VLAN instance with the specified UUID.
// Version: miami
func (vlan) GetByUUID2(session *Session, uuid string) (retval VLANRef, err error) {
	method := "VLAN.get_by_uuid"
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
	retval, err = deserializeVLANRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given VLAN.
// Version: miami
func (vlan) GetRecord(session *Session, self VLANRef) (retval VLANRecord, err error) {
	method := "VLAN.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVLANRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given VLAN.
// Version: miami
func (vlan) GetRecord2(session *Session, self VLANRef) (retval VLANRecord, err error) {
	method := "VLAN.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVLANRecord(method+" -> ", result)
	return
}
