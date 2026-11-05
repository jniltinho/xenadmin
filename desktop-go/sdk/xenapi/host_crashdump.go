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

type HostCrashdumpRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// Host the crashdump relates to
	Host HostRef `json:"host,omitempty"`
	// Time the crash happened
	Timestamp time.Time `json:"timestamp,omitempty"`
	// Size of the crashdump
	Size int `json:"size,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
}

type HostCrashdumpRef string

// Represents a host crash dump
type hostCrashdump struct{}

var HostCrashdump hostCrashdump

// GetAllRecords: Return a map of host_crashdump references to host_crashdump records for all host_crashdumps known to the system.
// Version: rio
func (hostCrashdump) GetAllRecords(session *Session) (retval map[HostCrashdumpRef]HostCrashdumpRecord, err error) {
	method := "host_crashdump.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostCrashdumpRefToHostCrashdumpRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of host_crashdump references to host_crashdump records for all host_crashdumps known to the system.
// Version: rio
func (hostCrashdump) GetAllRecords1(session *Session) (retval map[HostCrashdumpRef]HostCrashdumpRecord, err error) {
	method := "host_crashdump.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostCrashdumpRefToHostCrashdumpRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the host_crashdumps known to the system.
// Version: rio
func (hostCrashdump) GetAll(session *Session) (retval []HostCrashdumpRef, err error) {
	method := "host_crashdump.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostCrashdumpRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the host_crashdumps known to the system.
// Version: rio
func (hostCrashdump) GetAll1(session *Session) (retval []HostCrashdumpRef, err error) {
	method := "host_crashdump.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostCrashdumpRefSet(method+" -> ", result)
	return
}

// Upload: Upload the specified host crash dump to a specified URL
// Version: rio
func (hostCrashdump) Upload(session *Session, self HostCrashdumpRef, url string, options map[string]string) (err error) {
	method := "host_crashdump.upload"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, urlArg, optionsArg)
	return
}

// AsyncUpload: Upload the specified host crash dump to a specified URL
// Version: rio
func (hostCrashdump) AsyncUpload(session *Session, self HostCrashdumpRef, url string, options map[string]string) (retval TaskRef, err error) {
	method := "Async.host_crashdump.upload"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, urlArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Upload4: Upload the specified host crash dump to a specified URL
// Version: rio
func (hostCrashdump) Upload4(session *Session, self HostCrashdumpRef, url string, options map[string]string) (err error) {
	method := "host_crashdump.upload"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, urlArg, optionsArg)
	return
}

// AsyncUpload4: Upload the specified host crash dump to a specified URL
// Version: rio
func (hostCrashdump) AsyncUpload4(session *Session, self HostCrashdumpRef, url string, options map[string]string) (retval TaskRef, err error) {
	method := "Async.host_crashdump.upload"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, urlArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Destroy: Destroy specified host crash dump, removing it from the disk.
// Version: rio
func (hostCrashdump) Destroy(session *Session, self HostCrashdumpRef) (err error) {
	method := "host_crashdump.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy specified host crash dump, removing it from the disk.
// Version: rio
func (hostCrashdump) AsyncDestroy(session *Session, self HostCrashdumpRef) (retval TaskRef, err error) {
	method := "Async.host_crashdump.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy specified host crash dump, removing it from the disk.
// Version: rio
func (hostCrashdump) Destroy2(session *Session, self HostCrashdumpRef) (err error) {
	method := "host_crashdump.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy specified host crash dump, removing it from the disk.
// Version: rio
func (hostCrashdump) AsyncDestroy2(session *Session, self HostCrashdumpRef) (retval TaskRef, err error) {
	method := "Async.host_crashdump.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given host_crashdump.  If the key is not in that Map, then do nothing.
// Version: miami
func (hostCrashdump) RemoveFromOtherConfig(session *Session, self HostCrashdumpRef, key string) (err error) {
	method := "host_crashdump.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given host_crashdump.  If the key is not in that Map, then do nothing.
// Version: miami
func (hostCrashdump) RemoveFromOtherConfig3(session *Session, self HostCrashdumpRef, key string) (err error) {
	method := "host_crashdump.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig2: Remove the given key and its corresponding value from the other_config field of the given host_crashdump.  If the key is not in that Map, then do nothing.
// Version: rio
func (hostCrashdump) RemoveFromOtherConfig2(session *Session, self HostCrashdumpRef) (err error) {
	method := "host_crashdump.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given host_crashdump.
// Version: miami
func (hostCrashdump) AddToOtherConfig(session *Session, self HostCrashdumpRef, key string, value string) (err error) {
	method := "host_crashdump.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given host_crashdump.
// Version: miami
func (hostCrashdump) AddToOtherConfig4(session *Session, self HostCrashdumpRef, key string, value string) (err error) {
	method := "host_crashdump.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig2: Add the given key-value pair to the other_config field of the given host_crashdump.
// Version: rio
func (hostCrashdump) AddToOtherConfig2(session *Session, self HostCrashdumpRef) (err error) {
	method := "host_crashdump.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetOtherConfig: Set the other_config field of the given host_crashdump.
// Version: miami
func (hostCrashdump) SetOtherConfig(session *Session, self HostCrashdumpRef, value map[string]string) (err error) {
	method := "host_crashdump.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given host_crashdump.
// Version: miami
func (hostCrashdump) SetOtherConfig3(session *Session, self HostCrashdumpRef, value map[string]string) (err error) {
	method := "host_crashdump.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig2: Set the other_config field of the given host_crashdump.
// Version: rio
func (hostCrashdump) SetOtherConfig2(session *Session, self HostCrashdumpRef) (err error) {
	method := "host_crashdump.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// GetOtherConfig: Get the other_config field of the given host_crashdump.
// Version: rio
func (hostCrashdump) GetOtherConfig(session *Session, self HostCrashdumpRef) (retval map[string]string, err error) {
	method := "host_crashdump.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given host_crashdump.
// Version: rio
func (hostCrashdump) GetOtherConfig2(session *Session, self HostCrashdumpRef) (retval map[string]string, err error) {
	method := "host_crashdump.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSize: Get the size field of the given host_crashdump.
// Version: rio
func (hostCrashdump) GetSize(session *Session, self HostCrashdumpRef) (retval int, err error) {
	method := "host_crashdump.get_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSize2: Get the size field of the given host_crashdump.
// Version: rio
func (hostCrashdump) GetSize2(session *Session, self HostCrashdumpRef) (retval int, err error) {
	method := "host_crashdump.get_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTimestamp: Get the timestamp field of the given host_crashdump.
// Version: rio
func (hostCrashdump) GetTimestamp(session *Session, self HostCrashdumpRef) (retval time.Time, err error) {
	method := "host_crashdump.get_timestamp"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTimestamp2: Get the timestamp field of the given host_crashdump.
// Version: rio
func (hostCrashdump) GetTimestamp2(session *Session, self HostCrashdumpRef) (retval time.Time, err error) {
	method := "host_crashdump.get_timestamp"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHost: Get the host field of the given host_crashdump.
// Version: rio
func (hostCrashdump) GetHost(session *Session, self HostCrashdumpRef) (retval HostRef, err error) {
	method := "host_crashdump.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHost2: Get the host field of the given host_crashdump.
// Version: rio
func (hostCrashdump) GetHost2(session *Session, self HostCrashdumpRef) (retval HostRef, err error) {
	method := "host_crashdump.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given host_crashdump.
// Version: rio
func (hostCrashdump) GetUUID(session *Session, self HostCrashdumpRef) (retval string, err error) {
	method := "host_crashdump.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given host_crashdump.
// Version: rio
func (hostCrashdump) GetUUID2(session *Session, self HostCrashdumpRef) (retval string, err error) {
	method := "host_crashdump.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the host_crashdump instance with the specified UUID.
// Version: rio
func (hostCrashdump) GetByUUID(session *Session, uuid string) (retval HostCrashdumpRef, err error) {
	method := "host_crashdump.get_by_uuid"
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
	retval, err = deserializeHostCrashdumpRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the host_crashdump instance with the specified UUID.
// Version: rio
func (hostCrashdump) GetByUUID2(session *Session, uuid string) (retval HostCrashdumpRef, err error) {
	method := "host_crashdump.get_by_uuid"
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
	retval, err = deserializeHostCrashdumpRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given host_crashdump.
// Version: rio
func (hostCrashdump) GetRecord(session *Session, self HostCrashdumpRef) (retval HostCrashdumpRecord, err error) {
	method := "host_crashdump.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostCrashdumpRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given host_crashdump.
// Version: rio
func (hostCrashdump) GetRecord2(session *Session, self HostCrashdumpRef) (retval HostCrashdumpRecord, err error) {
	method := "host_crashdump.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeHostCrashdumpRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostCrashdumpRecord(method+" -> ", result)
	return
}
