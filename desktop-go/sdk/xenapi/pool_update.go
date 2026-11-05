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

type PoolUpdateRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// a human-readable name
	NameLabel string `json:"name_label,omitempty"`
	// a notes field containing human-readable description
	NameDescription string `json:"name_description,omitempty"`
	// Update version number
	Version string `json:"version,omitempty"`
	// Size of the update in bytes
	InstallationSize int `json:"installation_size,omitempty"`
	// GPG key of the update
	Key string `json:"key,omitempty"`
	// What the client should do after this update has been applied.
	AfterApplyGuidance []UpdateAfterApplyGuidance `json:"after_apply_guidance,omitempty"`
	// VDI the update was uploaded to
	Vdi VDIRef `json:"vdi,omitempty"`
	// The hosts that have applied this update.
	Hosts []HostRef `json:"hosts,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// Flag - if true, all hosts in a pool must apply this update
	EnforceHomogeneity bool `json:"enforce_homogeneity,omitempty"`
}

type PoolUpdateRef string

// Pool-wide updates to the host software
type poolUpdate struct{}

var PoolUpdate poolUpdate

// GetAllRecords: Return a map of pool_update references to pool_update records for all pool_updates known to the system.
// Version: ely
func (poolUpdate) GetAllRecords(session *Session) (retval map[PoolUpdateRef]PoolUpdateRecord, err error) {
	method := "pool_update.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRefToPoolUpdateRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of pool_update references to pool_update records for all pool_updates known to the system.
// Version: ely
func (poolUpdate) GetAllRecords1(session *Session) (retval map[PoolUpdateRef]PoolUpdateRecord, err error) {
	method := "pool_update.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRefToPoolUpdateRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the pool_updates known to the system.
// Version: ely
func (poolUpdate) GetAll(session *Session) (retval []PoolUpdateRef, err error) {
	method := "pool_update.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the pool_updates known to the system.
// Version: ely
func (poolUpdate) GetAll1(session *Session) (retval []PoolUpdateRef, err error) {
	method := "pool_update.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRefSet(method+" -> ", result)
	return
}

// Destroy: Removes the database entry. Only works on unapplied update.
// Version: ely
func (poolUpdate) Destroy(session *Session, self PoolUpdateRef) (err error) {
	method := "pool_update.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Removes the database entry. Only works on unapplied update.
// Version: ely
func (poolUpdate) AsyncDestroy(session *Session, self PoolUpdateRef) (retval TaskRef, err error) {
	method := "Async.pool_update.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Removes the database entry. Only works on unapplied update.
// Version: ely
func (poolUpdate) Destroy2(session *Session, self PoolUpdateRef) (err error) {
	method := "pool_update.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Removes the database entry. Only works on unapplied update.
// Version: ely
func (poolUpdate) AsyncDestroy2(session *Session, self PoolUpdateRef) (retval TaskRef, err error) {
	method := "Async.pool_update.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolClean: Removes the update&apos;s files from all hosts in the pool, but does not revert the update
// Version: ely
func (poolUpdate) PoolClean(session *Session, self PoolUpdateRef) (err error) {
	method := "pool_update.pool_clean"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolClean: Removes the update&apos;s files from all hosts in the pool, but does not revert the update
// Version: ely
func (poolUpdate) AsyncPoolClean(session *Session, self PoolUpdateRef) (retval TaskRef, err error) {
	method := "Async.pool_update.pool_clean"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolClean2: Removes the update&apos;s files from all hosts in the pool, but does not revert the update
// Version: ely
func (poolUpdate) PoolClean2(session *Session, self PoolUpdateRef) (err error) {
	method := "pool_update.pool_clean"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolClean2: Removes the update&apos;s files from all hosts in the pool, but does not revert the update
// Version: ely
func (poolUpdate) AsyncPoolClean2(session *Session, self PoolUpdateRef) (retval TaskRef, err error) {
	method := "Async.pool_update.pool_clean"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolApply: Apply the selected update to all hosts in the pool
// Version: ely
func (poolUpdate) PoolApply(session *Session, self PoolUpdateRef) (err error) {
	method := "pool_update.pool_apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolApply: Apply the selected update to all hosts in the pool
// Version: ely
func (poolUpdate) AsyncPoolApply(session *Session, self PoolUpdateRef) (retval TaskRef, err error) {
	method := "Async.pool_update.pool_apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolApply2: Apply the selected update to all hosts in the pool
// Version: ely
func (poolUpdate) PoolApply2(session *Session, self PoolUpdateRef) (err error) {
	method := "pool_update.pool_apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolApply2: Apply the selected update to all hosts in the pool
// Version: ely
func (poolUpdate) AsyncPoolApply2(session *Session, self PoolUpdateRef) (retval TaskRef, err error) {
	method := "Async.pool_update.pool_apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Apply: Apply the selected update to a host
// Version: ely
func (poolUpdate) Apply(session *Session, self PoolUpdateRef, host HostRef) (err error) {
	method := "pool_update.apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, hostArg)
	return
}

// AsyncApply: Apply the selected update to a host
// Version: ely
func (poolUpdate) AsyncApply(session *Session, self PoolUpdateRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool_update.apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Apply3: Apply the selected update to a host
// Version: ely
func (poolUpdate) Apply3(session *Session, self PoolUpdateRef, host HostRef) (err error) {
	method := "pool_update.apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, hostArg)
	return
}

// AsyncApply3: Apply the selected update to a host
// Version: ely
func (poolUpdate) AsyncApply3(session *Session, self PoolUpdateRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool_update.apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Precheck: Execute the precheck stage of the selected update on a host
// Version: ely
func (poolUpdate) Precheck(session *Session, self PoolUpdateRef, host HostRef) (retval LivepatchStatus, err error) {
	method := "pool_update.precheck"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumLivepatchStatus(method+" -> ", result)
	return
}

// AsyncPrecheck: Execute the precheck stage of the selected update on a host
// Version: ely
func (poolUpdate) AsyncPrecheck(session *Session, self PoolUpdateRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool_update.precheck"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Precheck3: Execute the precheck stage of the selected update on a host
// Version: ely
func (poolUpdate) Precheck3(session *Session, self PoolUpdateRef, host HostRef) (retval LivepatchStatus, err error) {
	method := "pool_update.precheck"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumLivepatchStatus(method+" -> ", result)
	return
}

// AsyncPrecheck3: Execute the precheck stage of the selected update on a host
// Version: ely
func (poolUpdate) AsyncPrecheck3(session *Session, self PoolUpdateRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool_update.precheck"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce: Introduce update VDI
// Version: ely
func (poolUpdate) Introduce(session *Session, vdi VDIRef) (retval PoolUpdateRef, err error) {
	method := "pool_update.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRef(method+" -> ", result)
	return
}

// AsyncIntroduce: Introduce update VDI
// Version: ely
func (poolUpdate) AsyncIntroduce(session *Session, vdi VDIRef) (retval TaskRef, err error) {
	method := "Async.pool_update.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce2: Introduce update VDI
// Version: ely
func (poolUpdate) Introduce2(session *Session, vdi VDIRef) (retval PoolUpdateRef, err error) {
	method := "pool_update.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRef(method+" -> ", result)
	return
}

// AsyncIntroduce2: Introduce update VDI
// Version: ely
func (poolUpdate) AsyncIntroduce2(session *Session, vdi VDIRef) (retval TaskRef, err error) {
	method := "Async.pool_update.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given pool_update.  If the key is not in that Map, then do nothing.
// Version: inverness
func (poolUpdate) RemoveFromOtherConfig(session *Session, self PoolUpdateRef, key string) (err error) {
	method := "pool_update.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given pool_update.  If the key is not in that Map, then do nothing.
// Version: inverness
func (poolUpdate) RemoveFromOtherConfig3(session *Session, self PoolUpdateRef, key string) (err error) {
	method := "pool_update.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig2: Remove the given key and its corresponding value from the other_config field of the given pool_update.  If the key is not in that Map, then do nothing.
// Version: ely
func (poolUpdate) RemoveFromOtherConfig2(session *Session, self PoolUpdateRef) (err error) {
	method := "pool_update.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given pool_update.
// Version: inverness
func (poolUpdate) AddToOtherConfig(session *Session, self PoolUpdateRef, key string, value string) (err error) {
	method := "pool_update.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given pool_update.
// Version: inverness
func (poolUpdate) AddToOtherConfig4(session *Session, self PoolUpdateRef, key string, value string) (err error) {
	method := "pool_update.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig2: Add the given key-value pair to the other_config field of the given pool_update.
// Version: ely
func (poolUpdate) AddToOtherConfig2(session *Session, self PoolUpdateRef) (err error) {
	method := "pool_update.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetOtherConfig: Set the other_config field of the given pool_update.
// Version: inverness
func (poolUpdate) SetOtherConfig(session *Session, self PoolUpdateRef, value map[string]string) (err error) {
	method := "pool_update.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given pool_update.
// Version: inverness
func (poolUpdate) SetOtherConfig3(session *Session, self PoolUpdateRef, value map[string]string) (err error) {
	method := "pool_update.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig2: Set the other_config field of the given pool_update.
// Version: ely
func (poolUpdate) SetOtherConfig2(session *Session, self PoolUpdateRef) (err error) {
	method := "pool_update.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// GetEnforceHomogeneity: Get the enforce_homogeneity field of the given pool_update.
// Version: ely
func (poolUpdate) GetEnforceHomogeneity(session *Session, self PoolUpdateRef) (retval bool, err error) {
	method := "pool_update.get_enforce_homogeneity"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetEnforceHomogeneity2: Get the enforce_homogeneity field of the given pool_update.
// Version: ely
func (poolUpdate) GetEnforceHomogeneity2(session *Session, self PoolUpdateRef) (retval bool, err error) {
	method := "pool_update.get_enforce_homogeneity"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given pool_update.
// Version: ely
func (poolUpdate) GetOtherConfig(session *Session, self PoolUpdateRef) (retval map[string]string, err error) {
	method := "pool_update.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given pool_update.
// Version: ely
func (poolUpdate) GetOtherConfig2(session *Session, self PoolUpdateRef) (retval map[string]string, err error) {
	method := "pool_update.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHosts: Get the hosts field of the given pool_update.
// Version: ely
func (poolUpdate) GetHosts(session *Session, self PoolUpdateRef) (retval []HostRef, err error) {
	method := "pool_update.get_hosts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRefSet(method+" -> ", result)
	return
}

// GetHosts2: Get the hosts field of the given pool_update.
// Version: ely
func (poolUpdate) GetHosts2(session *Session, self PoolUpdateRef) (retval []HostRef, err error) {
	method := "pool_update.get_hosts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRefSet(method+" -> ", result)
	return
}

// GetVdi: Get the vdi field of the given pool_update.
// Version: ely
func (poolUpdate) GetVdi(session *Session, self PoolUpdateRef) (retval VDIRef, err error) {
	method := "pool_update.get_vdi"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// GetVdi2: Get the vdi field of the given pool_update.
// Version: ely
func (poolUpdate) GetVdi2(session *Session, self PoolUpdateRef) (retval VDIRef, err error) {
	method := "pool_update.get_vdi"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// GetAfterApplyGuidance: Get the after_apply_guidance field of the given pool_update.
// Version: ely
func (poolUpdate) GetAfterApplyGuidance(session *Session, self PoolUpdateRef) (retval []UpdateAfterApplyGuidance, err error) {
	method := "pool_update.get_after_apply_guidance"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateAfterApplyGuidanceSet(method+" -> ", result)
	return
}

// GetAfterApplyGuidance2: Get the after_apply_guidance field of the given pool_update.
// Version: ely
func (poolUpdate) GetAfterApplyGuidance2(session *Session, self PoolUpdateRef) (retval []UpdateAfterApplyGuidance, err error) {
	method := "pool_update.get_after_apply_guidance"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateAfterApplyGuidanceSet(method+" -> ", result)
	return
}

// GetKey: Get the key field of the given pool_update.
// Version: ely
func (poolUpdate) GetKey(session *Session, self PoolUpdateRef) (retval string, err error) {
	method := "pool_update.get_key"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetKey2: Get the key field of the given pool_update.
// Version: ely
func (poolUpdate) GetKey2(session *Session, self PoolUpdateRef) (retval string, err error) {
	method := "pool_update.get_key"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetInstallationSize: Get the installation_size field of the given pool_update.
// Version: ely
func (poolUpdate) GetInstallationSize(session *Session, self PoolUpdateRef) (retval int, err error) {
	method := "pool_update.get_installation_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetInstallationSize2: Get the installation_size field of the given pool_update.
// Version: ely
func (poolUpdate) GetInstallationSize2(session *Session, self PoolUpdateRef) (retval int, err error) {
	method := "pool_update.get_installation_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVersion: Get the version field of the given pool_update.
// Version: ely
func (poolUpdate) GetVersion(session *Session, self PoolUpdateRef) (retval string, err error) {
	method := "pool_update.get_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVersion2: Get the version field of the given pool_update.
// Version: ely
func (poolUpdate) GetVersion2(session *Session, self PoolUpdateRef) (retval string, err error) {
	method := "pool_update.get_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription: Get the name/description field of the given pool_update.
// Version: ely
func (poolUpdate) GetNameDescription(session *Session, self PoolUpdateRef) (retval string, err error) {
	method := "pool_update.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription2: Get the name/description field of the given pool_update.
// Version: ely
func (poolUpdate) GetNameDescription2(session *Session, self PoolUpdateRef) (retval string, err error) {
	method := "pool_update.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel: Get the name/label field of the given pool_update.
// Version: ely
func (poolUpdate) GetNameLabel(session *Session, self PoolUpdateRef) (retval string, err error) {
	method := "pool_update.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel2: Get the name/label field of the given pool_update.
// Version: ely
func (poolUpdate) GetNameLabel2(session *Session, self PoolUpdateRef) (retval string, err error) {
	method := "pool_update.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given pool_update.
// Version: ely
func (poolUpdate) GetUUID(session *Session, self PoolUpdateRef) (retval string, err error) {
	method := "pool_update.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given pool_update.
// Version: ely
func (poolUpdate) GetUUID2(session *Session, self PoolUpdateRef) (retval string, err error) {
	method := "pool_update.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByNameLabel: Get all the pool_update instances with the given label.
// Version: ely
func (poolUpdate) GetByNameLabel(session *Session, label string) (retval []PoolUpdateRef, err error) {
	method := "pool_update.get_by_name_label"
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
	retval, err = deserializePoolUpdateRefSet(method+" -> ", result)
	return
}

// GetByNameLabel2: Get all the pool_update instances with the given label.
// Version: ely
func (poolUpdate) GetByNameLabel2(session *Session, label string) (retval []PoolUpdateRef, err error) {
	method := "pool_update.get_by_name_label"
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
	retval, err = deserializePoolUpdateRefSet(method+" -> ", result)
	return
}

// GetByUUID: Get a reference to the pool_update instance with the specified UUID.
// Version: ely
func (poolUpdate) GetByUUID(session *Session, uuid string) (retval PoolUpdateRef, err error) {
	method := "pool_update.get_by_uuid"
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
	retval, err = deserializePoolUpdateRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the pool_update instance with the specified UUID.
// Version: ely
func (poolUpdate) GetByUUID2(session *Session, uuid string) (retval PoolUpdateRef, err error) {
	method := "pool_update.get_by_uuid"
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
	retval, err = deserializePoolUpdateRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given pool_update.
// Version: ely
func (poolUpdate) GetRecord(session *Session, self PoolUpdateRef) (retval PoolUpdateRecord, err error) {
	method := "pool_update.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given pool_update.
// Version: ely
func (poolUpdate) GetRecord2(session *Session, self PoolUpdateRef) (retval PoolUpdateRecord, err error) {
	method := "pool_update.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolUpdateRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRecord(method+" -> ", result)
	return
}
