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

type PoolPatchRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// a human-readable name
	NameLabel string `json:"name_label,omitempty"`
	// a notes field containing human-readable description
	NameDescription string `json:"name_description,omitempty"`
	// Patch version number
	Version string `json:"version,omitempty"`
	// Size of the patch
	Size int `json:"size,omitempty"`
	// This patch should be applied across the entire pool
	PoolApplied bool `json:"pool_applied,omitempty"`
	// This hosts this patch is applied to.
	HostPatches []HostPatchRef `json:"host_patches,omitempty"`
	// What the client should do after this patch has been applied.
	AfterApplyGuidance []AfterApplyGuidance `json:"after_apply_guidance,omitempty"`
	// A reference to the associated pool_update object
	PoolUpdate PoolUpdateRef `json:"pool_update,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
}

type PoolPatchRef string

// Pool-wide patches
type poolPatch struct{}

var PoolPatch poolPatch

// GetAllRecords: Return a map of pool_patch references to pool_patch records for all pool_patchs known to the system.
// Version: miami
func (poolPatch) GetAllRecords(session *Session) (retval map[PoolPatchRef]PoolPatchRecord, err error) {
	method := "pool_patch.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolPatchRefToPoolPatchRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of pool_patch references to pool_patch records for all pool_patchs known to the system.
// Version: miami
func (poolPatch) GetAllRecords1(session *Session) (retval map[PoolPatchRef]PoolPatchRecord, err error) {
	method := "pool_patch.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolPatchRefToPoolPatchRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the pool_patchs known to the system.
// Version: miami
func (poolPatch) GetAll(session *Session) (retval []PoolPatchRef, err error) {
	method := "pool_patch.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolPatchRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the pool_patchs known to the system.
// Version: miami
func (poolPatch) GetAll1(session *Session) (retval []PoolPatchRef, err error) {
	method := "pool_patch.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolPatchRefSet(method+" -> ", result)
	return
}

// CleanOnHost: Removes the patch&apos;s files from the specified host
// Version: tampa
func (poolPatch) CleanOnHost(session *Session, self PoolPatchRef, host HostRef) (err error) {
	method := "pool_patch.clean_on_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncCleanOnHost: Removes the patch&apos;s files from the specified host
// Version: tampa
func (poolPatch) AsyncCleanOnHost(session *Session, self PoolPatchRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.clean_on_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// CleanOnHost3: Removes the patch&apos;s files from the specified host
// Version: tampa
func (poolPatch) CleanOnHost3(session *Session, self PoolPatchRef, host HostRef) (err error) {
	method := "pool_patch.clean_on_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncCleanOnHost3: Removes the patch&apos;s files from the specified host
// Version: tampa
func (poolPatch) AsyncCleanOnHost3(session *Session, self PoolPatchRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.clean_on_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy: Removes the patch&apos;s files from all hosts in the pool, and removes the database entries.  Only works on unapplied patches.
// Version: miami
func (poolPatch) Destroy(session *Session, self PoolPatchRef) (err error) {
	method := "pool_patch.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Removes the patch&apos;s files from all hosts in the pool, and removes the database entries.  Only works on unapplied patches.
// Version: miami
func (poolPatch) AsyncDestroy(session *Session, self PoolPatchRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Removes the patch&apos;s files from all hosts in the pool, and removes the database entries.  Only works on unapplied patches.
// Version: miami
func (poolPatch) Destroy2(session *Session, self PoolPatchRef) (err error) {
	method := "pool_patch.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Removes the patch&apos;s files from all hosts in the pool, and removes the database entries.  Only works on unapplied patches.
// Version: miami
func (poolPatch) AsyncDestroy2(session *Session, self PoolPatchRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolClean: Removes the patch&apos;s files from all hosts in the pool, but does not remove the database entries
// Version: tampa
func (poolPatch) PoolClean(session *Session, self PoolPatchRef) (err error) {
	method := "pool_patch.pool_clean"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolClean: Removes the patch&apos;s files from all hosts in the pool, but does not remove the database entries
// Version: tampa
func (poolPatch) AsyncPoolClean(session *Session, self PoolPatchRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.pool_clean"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolClean2: Removes the patch&apos;s files from all hosts in the pool, but does not remove the database entries
// Version: tampa
func (poolPatch) PoolClean2(session *Session, self PoolPatchRef) (err error) {
	method := "pool_patch.pool_clean"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolClean2: Removes the patch&apos;s files from all hosts in the pool, but does not remove the database entries
// Version: tampa
func (poolPatch) AsyncPoolClean2(session *Session, self PoolPatchRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.pool_clean"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Clean: Removes the patch&apos;s files from the server
// Version: miami
func (poolPatch) Clean(session *Session, self PoolPatchRef) (err error) {
	method := "pool_patch.clean"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncClean: Removes the patch&apos;s files from the server
// Version: miami
func (poolPatch) AsyncClean(session *Session, self PoolPatchRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.clean"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Clean2: Removes the patch&apos;s files from the server
// Version: miami
func (poolPatch) Clean2(session *Session, self PoolPatchRef) (err error) {
	method := "pool_patch.clean"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncClean2: Removes the patch&apos;s files from the server
// Version: miami
func (poolPatch) AsyncClean2(session *Session, self PoolPatchRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.clean"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Precheck: Execute the precheck stage of the selected patch on a host and return its output
// Version: miami
func (poolPatch) Precheck(session *Session, self PoolPatchRef, host HostRef) (retval string, err error) {
	method := "pool_patch.precheck"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncPrecheck: Execute the precheck stage of the selected patch on a host and return its output
// Version: miami
func (poolPatch) AsyncPrecheck(session *Session, self PoolPatchRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.precheck"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Precheck3: Execute the precheck stage of the selected patch on a host and return its output
// Version: miami
func (poolPatch) Precheck3(session *Session, self PoolPatchRef, host HostRef) (retval string, err error) {
	method := "pool_patch.precheck"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncPrecheck3: Execute the precheck stage of the selected patch on a host and return its output
// Version: miami
func (poolPatch) AsyncPrecheck3(session *Session, self PoolPatchRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.precheck"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolApply: Apply the selected patch to all hosts in the pool and return a map of host_ref -&gt; patch output
// Version: miami
func (poolPatch) PoolApply(session *Session, self PoolPatchRef) (err error) {
	method := "pool_patch.pool_apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolApply: Apply the selected patch to all hosts in the pool and return a map of host_ref -&gt; patch output
// Version: miami
func (poolPatch) AsyncPoolApply(session *Session, self PoolPatchRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.pool_apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolApply2: Apply the selected patch to all hosts in the pool and return a map of host_ref -&gt; patch output
// Version: miami
func (poolPatch) PoolApply2(session *Session, self PoolPatchRef) (err error) {
	method := "pool_patch.pool_apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolApply2: Apply the selected patch to all hosts in the pool and return a map of host_ref -&gt; patch output
// Version: miami
func (poolPatch) AsyncPoolApply2(session *Session, self PoolPatchRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.pool_apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Apply: Apply the selected patch to a host and return its output
// Version: miami
func (poolPatch) Apply(session *Session, self PoolPatchRef, host HostRef) (retval string, err error) {
	method := "pool_patch.apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncApply: Apply the selected patch to a host and return its output
// Version: miami
func (poolPatch) AsyncApply(session *Session, self PoolPatchRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Apply3: Apply the selected patch to a host and return its output
// Version: miami
func (poolPatch) Apply3(session *Session, self PoolPatchRef, host HostRef) (retval string, err error) {
	method := "pool_patch.apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncApply3: Apply the selected patch to a host and return its output
// Version: miami
func (poolPatch) AsyncApply3(session *Session, self PoolPatchRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.pool_patch.apply"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given pool_patch.  If the key is not in that Map, then do nothing.
// Version: miami
func (poolPatch) RemoveFromOtherConfig(session *Session, self PoolPatchRef, key string) (err error) {
	method := "pool_patch.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given pool_patch.  If the key is not in that Map, then do nothing.
// Version: miami
func (poolPatch) RemoveFromOtherConfig3(session *Session, self PoolPatchRef, key string) (err error) {
	method := "pool_patch.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given pool_patch.
// Version: miami
func (poolPatch) AddToOtherConfig(session *Session, self PoolPatchRef, key string, value string) (err error) {
	method := "pool_patch.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given pool_patch.
// Version: miami
func (poolPatch) AddToOtherConfig4(session *Session, self PoolPatchRef, key string, value string) (err error) {
	method := "pool_patch.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig: Set the other_config field of the given pool_patch.
// Version: miami
func (poolPatch) SetOtherConfig(session *Session, self PoolPatchRef, value map[string]string) (err error) {
	method := "pool_patch.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given pool_patch.
// Version: miami
func (poolPatch) SetOtherConfig3(session *Session, self PoolPatchRef, value map[string]string) (err error) {
	method := "pool_patch.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given pool_patch.
// Version: miami
func (poolPatch) GetOtherConfig(session *Session, self PoolPatchRef) (retval map[string]string, err error) {
	method := "pool_patch.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given pool_patch.
// Version: miami
func (poolPatch) GetOtherConfig2(session *Session, self PoolPatchRef) (retval map[string]string, err error) {
	method := "pool_patch.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPoolUpdate: Get the pool_update field of the given pool_patch.
// Version: miami
func (poolPatch) GetPoolUpdate(session *Session, self PoolPatchRef) (retval PoolUpdateRef, err error) {
	method := "pool_patch.get_pool_update"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRef(method+" -> ", result)
	return
}

// GetPoolUpdate2: Get the pool_update field of the given pool_patch.
// Version: miami
func (poolPatch) GetPoolUpdate2(session *Session, self PoolPatchRef) (retval PoolUpdateRef, err error) {
	method := "pool_patch.get_pool_update"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolUpdateRef(method+" -> ", result)
	return
}

// GetAfterApplyGuidance: Get the after_apply_guidance field of the given pool_patch.
// Version: miami
func (poolPatch) GetAfterApplyGuidance(session *Session, self PoolPatchRef) (retval []AfterApplyGuidance, err error) {
	method := "pool_patch.get_after_apply_guidance"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumAfterApplyGuidanceSet(method+" -> ", result)
	return
}

// GetAfterApplyGuidance2: Get the after_apply_guidance field of the given pool_patch.
// Version: miami
func (poolPatch) GetAfterApplyGuidance2(session *Session, self PoolPatchRef) (retval []AfterApplyGuidance, err error) {
	method := "pool_patch.get_after_apply_guidance"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumAfterApplyGuidanceSet(method+" -> ", result)
	return
}

// GetHostPatches: Get the host_patches field of the given pool_patch.
// Version: miami
func (poolPatch) GetHostPatches(session *Session, self PoolPatchRef) (retval []HostPatchRef, err error) {
	method := "pool_patch.get_host_patches"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostPatchRefSet(method+" -> ", result)
	return
}

// GetHostPatches2: Get the host_patches field of the given pool_patch.
// Version: miami
func (poolPatch) GetHostPatches2(session *Session, self PoolPatchRef) (retval []HostPatchRef, err error) {
	method := "pool_patch.get_host_patches"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostPatchRefSet(method+" -> ", result)
	return
}

// GetPoolApplied: Get the pool_applied field of the given pool_patch.
// Version: miami
func (poolPatch) GetPoolApplied(session *Session, self PoolPatchRef) (retval bool, err error) {
	method := "pool_patch.get_pool_applied"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPoolApplied2: Get the pool_applied field of the given pool_patch.
// Version: miami
func (poolPatch) GetPoolApplied2(session *Session, self PoolPatchRef) (retval bool, err error) {
	method := "pool_patch.get_pool_applied"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSize: Get the size field of the given pool_patch.
// Version: miami
func (poolPatch) GetSize(session *Session, self PoolPatchRef) (retval int, err error) {
	method := "pool_patch.get_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSize2: Get the size field of the given pool_patch.
// Version: miami
func (poolPatch) GetSize2(session *Session, self PoolPatchRef) (retval int, err error) {
	method := "pool_patch.get_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVersion: Get the version field of the given pool_patch.
// Version: miami
func (poolPatch) GetVersion(session *Session, self PoolPatchRef) (retval string, err error) {
	method := "pool_patch.get_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVersion2: Get the version field of the given pool_patch.
// Version: miami
func (poolPatch) GetVersion2(session *Session, self PoolPatchRef) (retval string, err error) {
	method := "pool_patch.get_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription: Get the name/description field of the given pool_patch.
// Version: miami
func (poolPatch) GetNameDescription(session *Session, self PoolPatchRef) (retval string, err error) {
	method := "pool_patch.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription2: Get the name/description field of the given pool_patch.
// Version: miami
func (poolPatch) GetNameDescription2(session *Session, self PoolPatchRef) (retval string, err error) {
	method := "pool_patch.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel: Get the name/label field of the given pool_patch.
// Version: miami
func (poolPatch) GetNameLabel(session *Session, self PoolPatchRef) (retval string, err error) {
	method := "pool_patch.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel2: Get the name/label field of the given pool_patch.
// Version: miami
func (poolPatch) GetNameLabel2(session *Session, self PoolPatchRef) (retval string, err error) {
	method := "pool_patch.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given pool_patch.
// Version: miami
func (poolPatch) GetUUID(session *Session, self PoolPatchRef) (retval string, err error) {
	method := "pool_patch.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given pool_patch.
// Version: miami
func (poolPatch) GetUUID2(session *Session, self PoolPatchRef) (retval string, err error) {
	method := "pool_patch.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByNameLabel: Get all the pool_patch instances with the given label.
// Version: miami
func (poolPatch) GetByNameLabel(session *Session, label string) (retval []PoolPatchRef, err error) {
	method := "pool_patch.get_by_name_label"
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
	retval, err = deserializePoolPatchRefSet(method+" -> ", result)
	return
}

// GetByNameLabel2: Get all the pool_patch instances with the given label.
// Version: miami
func (poolPatch) GetByNameLabel2(session *Session, label string) (retval []PoolPatchRef, err error) {
	method := "pool_patch.get_by_name_label"
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
	retval, err = deserializePoolPatchRefSet(method+" -> ", result)
	return
}

// GetByUUID: Get a reference to the pool_patch instance with the specified UUID.
// Version: miami
func (poolPatch) GetByUUID(session *Session, uuid string) (retval PoolPatchRef, err error) {
	method := "pool_patch.get_by_uuid"
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
	retval, err = deserializePoolPatchRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the pool_patch instance with the specified UUID.
// Version: miami
func (poolPatch) GetByUUID2(session *Session, uuid string) (retval PoolPatchRef, err error) {
	method := "pool_patch.get_by_uuid"
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
	retval, err = deserializePoolPatchRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given pool_patch.
// Version: miami
func (poolPatch) GetRecord(session *Session, self PoolPatchRef) (retval PoolPatchRecord, err error) {
	method := "pool_patch.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolPatchRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given pool_patch.
// Version: miami
func (poolPatch) GetRecord2(session *Session, self PoolPatchRef) (retval PoolPatchRecord, err error) {
	method := "pool_patch.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePoolPatchRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolPatchRecord(method+" -> ", result)
	return
}
