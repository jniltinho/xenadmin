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

type VDIRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// a human-readable name
	NameLabel string `json:"name_label,omitempty"`
	// a notes field containing human-readable description
	NameDescription string `json:"name_description,omitempty"`
	// list of the operations allowed in this state. This list is advisory only and the server state may have changed by the time this field is read by a client.
	AllowedOperations []VdiOperations `json:"allowed_operations,omitempty"`
	// links each of the running tasks using this object (by reference) to a current_operation enum which describes the nature of the task.
	CurrentOperations map[string]VdiOperations `json:"current_operations,omitempty"`
	// storage repository in which the VDI resides
	SR SRRef `json:"SR,omitempty"`
	// list of vbds that refer to this disk
	VBDs []VBDRef `json:"VBDs,omitempty"`
	// list of crash dumps that refer to this disk
	CrashDumps []CrashdumpRef `json:"crash_dumps,omitempty"`
	// size of disk as presented to the guest (in bytes). Note that, depending on storage backend type, requested size may not be respected exactly
	VirtualSize int `json:"virtual_size,omitempty"`
	// amount of physical space that the disk image is currently taking up on the storage repository (in bytes)
	PhysicalUtilisation int `json:"physical_utilisation,omitempty"`
	// type of the VDI
	Type VdiType `json:"type,omitempty"`
	// true if this disk may be shared
	Sharable bool `json:"sharable,omitempty"`
	// true if this disk may ONLY be mounted read-only
	ReadOnly bool `json:"read_only,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// true if this disk is locked at the storage level
	StorageLock bool `json:"storage_lock,omitempty"`
	// location information
	Location string `json:"location,omitempty"`
	//
	Managed bool `json:"managed,omitempty"`
	// true if SR scan operation reported this VDI as not present on disk
	Missing bool `json:"missing,omitempty"`
	// This field is always null. Deprecated
	Parent VDIRef `json:"parent,omitempty"`
	// data to be inserted into the xenstore tree (/local/domain/0/backend/vbd/&lt;domid&gt;/&lt;device-id&gt;/sm-data) after the VDI is attached. This is generally set by the SM backends on vdi_attach.
	XenstoreData map[string]string `json:"xenstore_data,omitempty"`
	// SM dependent data
	SmConfig map[string]string `json:"sm_config,omitempty"`
	// true if this is a snapshot.
	IsASnapshot bool `json:"is_a_snapshot,omitempty"`
	// Ref pointing to the VDI this snapshot is of.
	SnapshotOf VDIRef `json:"snapshot_of,omitempty"`
	// List pointing to all the VDIs snapshots.
	Snapshots []VDIRef `json:"snapshots,omitempty"`
	// Date/time when this snapshot was created.
	SnapshotTime time.Time `json:"snapshot_time,omitempty"`
	// user-specified tags for categorization purposes
	Tags []string `json:"tags,omitempty"`
	// true if this VDI is to be cached in the local cache SR
	AllowCaching bool `json:"allow_caching,omitempty"`
	// The behaviour of this VDI on a VM boot
	OnBoot OnBoot `json:"on_boot,omitempty"`
	// The pool whose metadata is contained in this VDI
	MetadataOfPool PoolRef `json:"metadata_of_pool,omitempty"`
	// Whether this VDI contains the latest known accessible metadata for the pool
	MetadataLatest bool `json:"metadata_latest,omitempty"`
	// Whether this VDI is a Tools ISO
	IsToolsIso bool `json:"is_tools_iso,omitempty"`
	// True if changed blocks are tracked for this VDI
	CbtEnabled bool `json:"cbt_enabled,omitempty"`
}

type VDIRef string

// A virtual disk image
type vdi struct{}

var VDI vdi

// GetAllRecords: Return a map of VDI references to VDI records for all VDIs known to the system.
// Version: rio
func (vdi) GetAllRecords(session *Session) (retval map[VDIRef]VDIRecord, err error) {
	method := "VDI.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRefToVDIRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of VDI references to VDI records for all VDIs known to the system.
// Version: rio
func (vdi) GetAllRecords1(session *Session) (retval map[VDIRef]VDIRecord, err error) {
	method := "VDI.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRefToVDIRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the VDIs known to the system.
// Version: rio
func (vdi) GetAll(session *Session) (retval []VDIRef, err error) {
	method := "VDI.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the VDIs known to the system.
// Version: rio
func (vdi) GetAll1(session *Session) (retval []VDIRef, err error) {
	method := "VDI.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRefSet(method+" -> ", result)
	return
}

// Revert: Copy the contents of a snapshot to the VDI it&apos;s related to. The original contents of the VDI are lost.
// Version: closed
//
// Errors:
// UNIMPLEMENTED_IN_SM_BACKEND - You have attempted a function which is not implemented
func (vdi) Revert(session *Session, snapshot VDIRef) (err error) {
	method := "VDI.revert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	snapshotArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "snapshot"), snapshot)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, snapshotArg)
	return
}

// AsyncRevert: Copy the contents of a snapshot to the VDI it&apos;s related to. The original contents of the VDI are lost.
// Version: closed
//
// Errors:
// UNIMPLEMENTED_IN_SM_BACKEND - You have attempted a function which is not implemented
func (vdi) AsyncRevert(session *Session, snapshot VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.revert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	snapshotArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "snapshot"), snapshot)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, snapshotArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Revert2: Copy the contents of a snapshot to the VDI it&apos;s related to. The original contents of the VDI are lost.
// Version: closed
//
// Errors:
// UNIMPLEMENTED_IN_SM_BACKEND - You have attempted a function which is not implemented
func (vdi) Revert2(session *Session, snapshot VDIRef) (err error) {
	method := "VDI.revert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	snapshotArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "snapshot"), snapshot)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, snapshotArg)
	return
}

// AsyncRevert2: Copy the contents of a snapshot to the VDI it&apos;s related to. The original contents of the VDI are lost.
// Version: closed
//
// Errors:
// UNIMPLEMENTED_IN_SM_BACKEND - You have attempted a function which is not implemented
func (vdi) AsyncRevert2(session *Session, snapshot VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.revert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	snapshotArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "snapshot"), snapshot)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, snapshotArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetNbdInfo: Get details specifying how to access this VDI via a Network Block Device server. For each of a set of NBD server addresses on which the VDI is available, the return value set contains a vdi_nbd_server_info object that contains an exportname to request once the NBD connection is established, and connection details for the address. An empty list is returned if there is no network that has a PIF on a host with access to the relevant SR, or if no such network has been assigned an NBD-related purpose in its purpose field. To access the given VDI, any of the vdi_nbd_server_info objects can be used to make a connection to a server, and then the VDI will be available by requesting the exportname.
// Version: inverness
//
// Errors:
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
func (vdi) GetNbdInfo(session *Session, self VDIRef) (retval []VdiNbdServerInfoRecord, err error) {
	method := "VDI.get_nbd_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVdiNbdServerInfoRecordSet(method+" -> ", result)
	return
}

// GetNbdInfo2: Get details specifying how to access this VDI via a Network Block Device server. For each of a set of NBD server addresses on which the VDI is available, the return value set contains a vdi_nbd_server_info object that contains an exportname to request once the NBD connection is established, and connection details for the address. An empty list is returned if there is no network that has a PIF on a host with access to the relevant SR, or if no such network has been assigned an NBD-related purpose in its purpose field. To access the given VDI, any of the vdi_nbd_server_info objects can be used to make a connection to a server, and then the VDI will be available by requesting the exportname.
// Version: inverness
//
// Errors:
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
func (vdi) GetNbdInfo2(session *Session, self VDIRef) (retval []VdiNbdServerInfoRecord, err error) {
	method := "VDI.get_nbd_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVdiNbdServerInfoRecordSet(method+" -> ", result)
	return
}

// ListChangedBlocks: Compare two VDIs in 64k block increments and report which blocks differ. This operation is not allowed when vdi_to is attached to a VM.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// VDI_IN_USE - This operation cannot be performed because this VDI is in use by some other operation
func (vdi) ListChangedBlocks(session *Session, vdiFrom VDIRef, vdiTo VDIRef) (retval string, err error) {
	method := "VDI.list_changed_blocks"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiFromArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi_from"), vdiFrom)
	if err != nil {
		return
	}
	vdiToArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi_to"), vdiTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiFromArg, vdiToArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncListChangedBlocks: Compare two VDIs in 64k block increments and report which blocks differ. This operation is not allowed when vdi_to is attached to a VM.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// VDI_IN_USE - This operation cannot be performed because this VDI is in use by some other operation
func (vdi) AsyncListChangedBlocks(session *Session, vdiFrom VDIRef, vdiTo VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.list_changed_blocks"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiFromArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi_from"), vdiFrom)
	if err != nil {
		return
	}
	vdiToArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi_to"), vdiTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiFromArg, vdiToArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ListChangedBlocks3: Compare two VDIs in 64k block increments and report which blocks differ. This operation is not allowed when vdi_to is attached to a VM.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// VDI_IN_USE - This operation cannot be performed because this VDI is in use by some other operation
func (vdi) ListChangedBlocks3(session *Session, vdiFrom VDIRef, vdiTo VDIRef) (retval string, err error) {
	method := "VDI.list_changed_blocks"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiFromArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi_from"), vdiFrom)
	if err != nil {
		return
	}
	vdiToArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi_to"), vdiTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiFromArg, vdiToArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncListChangedBlocks3: Compare two VDIs in 64k block increments and report which blocks differ. This operation is not allowed when vdi_to is attached to a VM.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// VDI_IN_USE - This operation cannot be performed because this VDI is in use by some other operation
func (vdi) AsyncListChangedBlocks3(session *Session, vdiFrom VDIRef, vdiTo VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.list_changed_blocks"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiFromArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi_from"), vdiFrom)
	if err != nil {
		return
	}
	vdiToArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi_to"), vdiTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiFromArg, vdiToArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DataDestroy: Delete the data of the snapshot VDI, but keep its changed block tracking metadata. When successful, this call changes the type of the VDI to cbt_metadata. This operation is idempotent: calling it on a VDI of type cbt_metadata results in a no-op, and no error will be thrown.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
// VDI_NO_CBT_METADATA - The requested operation is not allowed because the specified VDI does not have changed block tracking metadata.
// VDI_IN_USE - This operation cannot be performed because this VDI is in use by some other operation
// VDI_IS_A_PHYSICAL_DEVICE - The operation cannot be performed on physical device
func (vdi) DataDestroy(session *Session, self VDIRef) (err error) {
	method := "VDI.data_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDataDestroy: Delete the data of the snapshot VDI, but keep its changed block tracking metadata. When successful, this call changes the type of the VDI to cbt_metadata. This operation is idempotent: calling it on a VDI of type cbt_metadata results in a no-op, and no error will be thrown.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
// VDI_NO_CBT_METADATA - The requested operation is not allowed because the specified VDI does not have changed block tracking metadata.
// VDI_IN_USE - This operation cannot be performed because this VDI is in use by some other operation
// VDI_IS_A_PHYSICAL_DEVICE - The operation cannot be performed on physical device
func (vdi) AsyncDataDestroy(session *Session, self VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.data_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DataDestroy2: Delete the data of the snapshot VDI, but keep its changed block tracking metadata. When successful, this call changes the type of the VDI to cbt_metadata. This operation is idempotent: calling it on a VDI of type cbt_metadata results in a no-op, and no error will be thrown.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
// VDI_NO_CBT_METADATA - The requested operation is not allowed because the specified VDI does not have changed block tracking metadata.
// VDI_IN_USE - This operation cannot be performed because this VDI is in use by some other operation
// VDI_IS_A_PHYSICAL_DEVICE - The operation cannot be performed on physical device
func (vdi) DataDestroy2(session *Session, self VDIRef) (err error) {
	method := "VDI.data_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDataDestroy2: Delete the data of the snapshot VDI, but keep its changed block tracking metadata. When successful, this call changes the type of the VDI to cbt_metadata. This operation is idempotent: calling it on a VDI of type cbt_metadata results in a no-op, and no error will be thrown.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
// VDI_NO_CBT_METADATA - The requested operation is not allowed because the specified VDI does not have changed block tracking metadata.
// VDI_IN_USE - This operation cannot be performed because this VDI is in use by some other operation
// VDI_IS_A_PHYSICAL_DEVICE - The operation cannot be performed on physical device
func (vdi) AsyncDataDestroy2(session *Session, self VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.data_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableCbt: Disable changed block tracking for the VDI. This call is only allowed on VDIs that support enabling CBT. It is an idempotent operation - disabling CBT for a VDI for which CBT is not enabled results in a no-op, and no error will be thrown.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
// VDI_ON_BOOT_MODE_INCOMPATIBLE_WITH_OPERATION - This operation is not permitted on VDIs in the &apos;on-boot=reset&apos; mode, or on VMs having such VDIs.
func (vdi) DisableCbt(session *Session, self VDIRef) (err error) {
	method := "VDI.disable_cbt"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableCbt: Disable changed block tracking for the VDI. This call is only allowed on VDIs that support enabling CBT. It is an idempotent operation - disabling CBT for a VDI for which CBT is not enabled results in a no-op, and no error will be thrown.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
// VDI_ON_BOOT_MODE_INCOMPATIBLE_WITH_OPERATION - This operation is not permitted on VDIs in the &apos;on-boot=reset&apos; mode, or on VMs having such VDIs.
func (vdi) AsyncDisableCbt(session *Session, self VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.disable_cbt"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DisableCbt2: Disable changed block tracking for the VDI. This call is only allowed on VDIs that support enabling CBT. It is an idempotent operation - disabling CBT for a VDI for which CBT is not enabled results in a no-op, and no error will be thrown.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
// VDI_ON_BOOT_MODE_INCOMPATIBLE_WITH_OPERATION - This operation is not permitted on VDIs in the &apos;on-boot=reset&apos; mode, or on VMs having such VDIs.
func (vdi) DisableCbt2(session *Session, self VDIRef) (err error) {
	method := "VDI.disable_cbt"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisableCbt2: Disable changed block tracking for the VDI. This call is only allowed on VDIs that support enabling CBT. It is an idempotent operation - disabling CBT for a VDI for which CBT is not enabled results in a no-op, and no error will be thrown.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
// VDI_ON_BOOT_MODE_INCOMPATIBLE_WITH_OPERATION - This operation is not permitted on VDIs in the &apos;on-boot=reset&apos; mode, or on VMs having such VDIs.
func (vdi) AsyncDisableCbt2(session *Session, self VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.disable_cbt"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// EnableCbt: Enable changed block tracking for the VDI. This call is idempotent - enabling CBT for a VDI for which CBT is already enabled results in a no-op, and no error will be thrown.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
// VDI_ON_BOOT_MODE_INCOMPATIBLE_WITH_OPERATION - This operation is not permitted on VDIs in the &apos;on-boot=reset&apos; mode, or on VMs having such VDIs.
func (vdi) EnableCbt(session *Session, self VDIRef) (err error) {
	method := "VDI.enable_cbt"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncEnableCbt: Enable changed block tracking for the VDI. This call is idempotent - enabling CBT for a VDI for which CBT is already enabled results in a no-op, and no error will be thrown.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
// VDI_ON_BOOT_MODE_INCOMPATIBLE_WITH_OPERATION - This operation is not permitted on VDIs in the &apos;on-boot=reset&apos; mode, or on VMs having such VDIs.
func (vdi) AsyncEnableCbt(session *Session, self VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.enable_cbt"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// EnableCbt2: Enable changed block tracking for the VDI. This call is idempotent - enabling CBT for a VDI for which CBT is already enabled results in a no-op, and no error will be thrown.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
// VDI_ON_BOOT_MODE_INCOMPATIBLE_WITH_OPERATION - This operation is not permitted on VDIs in the &apos;on-boot=reset&apos; mode, or on VMs having such VDIs.
func (vdi) EnableCbt2(session *Session, self VDIRef) (err error) {
	method := "VDI.enable_cbt"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncEnableCbt2: Enable changed block tracking for the VDI. This call is idempotent - enabling CBT for a VDI for which CBT is already enabled results in a no-op, and no error will be thrown.
// Version: inverness
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
// VDI_MISSING - This operation cannot be performed because the specified VDI could not be found on the storage substrate
// SR_NOT_ATTACHED - The SR is not attached.
// SR_HAS_NO_PBDS - The SR has no attached PBDs
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VDI_INCOMPATIBLE_TYPE - This operation cannot be performed because the specified VDI is of an incompatible type (eg: an HA statefile cannot be attached to a guest)
// VDI_ON_BOOT_MODE_INCOMPATIBLE_WITH_OPERATION - This operation is not permitted on VDIs in the &apos;on-boot=reset&apos; mode, or on VMs having such VDIs.
func (vdi) AsyncEnableCbt2(session *Session, self VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.enable_cbt"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolMigrate: Migrate a VDI, which may be attached to a running guest, to a different SR. The destination SR must be visible to the guest.
// Version: tampa
func (vdi) PoolMigrate(session *Session, vdi VDIRef, sr SRRef, options map[string]string) (retval VDIRef, err error) {
	method := "VDI.pool_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, srArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncPoolMigrate: Migrate a VDI, which may be attached to a running guest, to a different SR. The destination SR must be visible to the guest.
// Version: tampa
func (vdi) AsyncPoolMigrate(session *Session, vdi VDIRef, sr SRRef, options map[string]string) (retval TaskRef, err error) {
	method := "Async.VDI.pool_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, srArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// PoolMigrate4: Migrate a VDI, which may be attached to a running guest, to a different SR. The destination SR must be visible to the guest.
// Version: tampa
func (vdi) PoolMigrate4(session *Session, vdi VDIRef, sr SRRef, options map[string]string) (retval VDIRef, err error) {
	method := "VDI.pool_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, srArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncPoolMigrate4: Migrate a VDI, which may be attached to a running guest, to a different SR. The destination SR must be visible to the guest.
// Version: tampa
func (vdi) AsyncPoolMigrate4(session *Session, vdi VDIRef, sr SRRef, options map[string]string) (retval TaskRef, err error) {
	method := "Async.VDI.pool_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, srArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ReadDatabasePoolUUID: Check the VDI cache for the pool UUID of the database on this VDI.
// Version: boston
func (vdi) ReadDatabasePoolUUID(session *Session, self VDIRef) (retval string, err error) {
	method := "VDI.read_database_pool_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncReadDatabasePoolUUID: Check the VDI cache for the pool UUID of the database on this VDI.
// Version: boston
func (vdi) AsyncReadDatabasePoolUUID(session *Session, self VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.read_database_pool_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// ReadDatabasePoolUUID2: Check the VDI cache for the pool UUID of the database on this VDI.
// Version: boston
func (vdi) ReadDatabasePoolUUID2(session *Session, self VDIRef) (retval string, err error) {
	method := "VDI.read_database_pool_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncReadDatabasePoolUUID2: Check the VDI cache for the pool UUID of the database on this VDI.
// Version: boston
func (vdi) AsyncReadDatabasePoolUUID2(session *Session, self VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.read_database_pool_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// OpenDatabase: Load the metadata found on the supplied VDI and return a session reference which can be used in API calls to query its contents.
// Version: boston
func (vdi) OpenDatabase(session *Session, self VDIRef) (retval SessionRef, err error) {
	method := "VDI.open_database"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	return
}

// AsyncOpenDatabase: Load the metadata found on the supplied VDI and return a session reference which can be used in API calls to query its contents.
// Version: boston
func (vdi) AsyncOpenDatabase(session *Session, self VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.open_database"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// OpenDatabase2: Load the metadata found on the supplied VDI and return a session reference which can be used in API calls to query its contents.
// Version: boston
func (vdi) OpenDatabase2(session *Session, self VDIRef) (retval SessionRef, err error) {
	method := "VDI.open_database"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSessionRef(method+" -> ", result)
	return
}

// AsyncOpenDatabase2: Load the metadata found on the supplied VDI and return a session reference which can be used in API calls to query its contents.
// Version: boston
func (vdi) AsyncOpenDatabase2(session *Session, self VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.open_database"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetAllowCaching: Set the value of the allow_caching parameter. This value can only be changed when the VDI is not attached to a running VM. The caching behaviour is only affected by this flag for VHD-based VDIs that have one parent and no child VHDs. Moreover, caching only takes place when the host running the VM containing this VDI has a nominated SR for local caching.
// Version: cowley
func (vdi) SetAllowCaching(session *Session, self VDIRef, value bool) (err error) {
	method := "VDI.set_allow_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetAllowCaching: Set the value of the allow_caching parameter. This value can only be changed when the VDI is not attached to a running VM. The caching behaviour is only affected by this flag for VHD-based VDIs that have one parent and no child VHDs. Moreover, caching only takes place when the host running the VM containing this VDI has a nominated SR for local caching.
// Version: cowley
func (vdi) AsyncSetAllowCaching(session *Session, self VDIRef, value bool) (retval TaskRef, err error) {
	method := "Async.VDI.set_allow_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetAllowCaching3: Set the value of the allow_caching parameter. This value can only be changed when the VDI is not attached to a running VM. The caching behaviour is only affected by this flag for VHD-based VDIs that have one parent and no child VHDs. Moreover, caching only takes place when the host running the VM containing this VDI has a nominated SR for local caching.
// Version: cowley
func (vdi) SetAllowCaching3(session *Session, self VDIRef, value bool) (err error) {
	method := "VDI.set_allow_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetAllowCaching3: Set the value of the allow_caching parameter. This value can only be changed when the VDI is not attached to a running VM. The caching behaviour is only affected by this flag for VHD-based VDIs that have one parent and no child VHDs. Moreover, caching only takes place when the host running the VM containing this VDI has a nominated SR for local caching.
// Version: cowley
func (vdi) AsyncSetAllowCaching3(session *Session, self VDIRef, value bool) (retval TaskRef, err error) {
	method := "Async.VDI.set_allow_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetOnBoot: Set the value of the on_boot parameter. This value can only be changed when the VDI is not attached to a running VM.
// Version: cowley
func (vdi) SetOnBoot(session *Session, self VDIRef, value OnBoot) (err error) {
	method := "VDI.set_on_boot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnBoot(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetOnBoot: Set the value of the on_boot parameter. This value can only be changed when the VDI is not attached to a running VM.
// Version: cowley
func (vdi) AsyncSetOnBoot(session *Session, self VDIRef, value OnBoot) (retval TaskRef, err error) {
	method := "Async.VDI.set_on_boot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnBoot(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetOnBoot3: Set the value of the on_boot parameter. This value can only be changed when the VDI is not attached to a running VM.
// Version: cowley
func (vdi) SetOnBoot3(session *Session, self VDIRef, value OnBoot) (err error) {
	method := "VDI.set_on_boot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnBoot(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetOnBoot3: Set the value of the on_boot parameter. This value can only be changed when the VDI is not attached to a running VM.
// Version: cowley
func (vdi) AsyncSetOnBoot3(session *Session, self VDIRef, value OnBoot) (retval TaskRef, err error) {
	method := "Async.VDI.set_on_boot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnBoot(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetNameDescription: Set the name description of the VDI. This can only happen when its SR is currently attached.
// Version: rio
func (vdi) SetNameDescription(session *Session, self VDIRef, value string) (err error) {
	method := "VDI.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetNameDescription: Set the name description of the VDI. This can only happen when its SR is currently attached.
// Version: rio
func (vdi) AsyncSetNameDescription(session *Session, self VDIRef, value string) (retval TaskRef, err error) {
	method := "Async.VDI.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetNameDescription3: Set the name description of the VDI. This can only happen when its SR is currently attached.
// Version: rio
func (vdi) SetNameDescription3(session *Session, self VDIRef, value string) (err error) {
	method := "VDI.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetNameDescription3: Set the name description of the VDI. This can only happen when its SR is currently attached.
// Version: rio
func (vdi) AsyncSetNameDescription3(session *Session, self VDIRef, value string) (retval TaskRef, err error) {
	method := "Async.VDI.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetNameLabel: Set the name label of the VDI. This can only happen when then its SR is currently attached.
// Version: rio
func (vdi) SetNameLabel(session *Session, self VDIRef, value string) (err error) {
	method := "VDI.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetNameLabel: Set the name label of the VDI. This can only happen when then its SR is currently attached.
// Version: rio
func (vdi) AsyncSetNameLabel(session *Session, self VDIRef, value string) (retval TaskRef, err error) {
	method := "Async.VDI.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetNameLabel3: Set the name label of the VDI. This can only happen when then its SR is currently attached.
// Version: rio
func (vdi) SetNameLabel3(session *Session, self VDIRef, value string) (err error) {
	method := "VDI.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetNameLabel3: Set the name label of the VDI. This can only happen when then its SR is currently attached.
// Version: rio
func (vdi) AsyncSetNameLabel3(session *Session, self VDIRef, value string) (retval TaskRef, err error) {
	method := "Async.VDI.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetReadOnly: Sets the VDI&apos;s read_only field
// Version: rio
func (vdi) SetReadOnly(session *Session, self VDIRef, value bool) (err error) {
	method := "VDI.set_read_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetReadOnly3: Sets the VDI&apos;s read_only field
// Version: rio
func (vdi) SetReadOnly3(session *Session, self VDIRef, value bool) (err error) {
	method := "VDI.set_read_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetSharable: Sets the VDI&apos;s sharable field
// Version: george
func (vdi) SetSharable(session *Session, self VDIRef, value bool) (err error) {
	method := "VDI.set_sharable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetSharable3: Sets the VDI&apos;s sharable field
// Version: george
func (vdi) SetSharable3(session *Session, self VDIRef, value bool) (err error) {
	method := "VDI.set_sharable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// Forget: Removes a VDI record from the database
// Version: rio
func (vdi) Forget(session *Session, vdi VDIRef) (err error) {
	method := "VDI.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vdiArg)
	return
}

// AsyncForget: Removes a VDI record from the database
// Version: rio
func (vdi) AsyncForget(session *Session, vdi VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.forget"
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

// Forget2: Removes a VDI record from the database
// Version: rio
func (vdi) Forget2(session *Session, vdi VDIRef) (err error) {
	method := "VDI.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vdiArg)
	return
}

// AsyncForget2: Removes a VDI record from the database
// Version: rio
func (vdi) AsyncForget2(session *Session, vdi VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.forget"
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

// Copy: Copy either a full VDI or the block differences between two VDIs into either a fresh VDI or an existing VDI.
// Version: clearwater-felton
//
// Errors:
// VDI_READONLY - The operation required write access but this VDI is read-only
// VDI_TOO_SMALL - The VDI is too small. Please resize it to at least the minimum size.
// VDI_NOT_SPARSE - The VDI is not stored using a sparse format. It is not possible to query and manipulate only the changed blocks (or &apos;block differences&apos; or &apos;disk deltas&apos;) between two VDIs. Please select a VDI which uses a sparse-aware technology such as VHD.
func (vdi) Copy(session *Session, vdi VDIRef, sr SRRef, baseVdi VDIRef, intoVdi VDIRef) (retval VDIRef, err error) {
	method := "VDI.copy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	baseVdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "base_vdi"), baseVdi)
	if err != nil {
		return
	}
	intoVdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "into_vdi"), intoVdi)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, srArg, baseVdiArg, intoVdiArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncCopy: Copy either a full VDI or the block differences between two VDIs into either a fresh VDI or an existing VDI.
// Version: clearwater-felton
//
// Errors:
// VDI_READONLY - The operation required write access but this VDI is read-only
// VDI_TOO_SMALL - The VDI is too small. Please resize it to at least the minimum size.
// VDI_NOT_SPARSE - The VDI is not stored using a sparse format. It is not possible to query and manipulate only the changed blocks (or &apos;block differences&apos; or &apos;disk deltas&apos;) between two VDIs. Please select a VDI which uses a sparse-aware technology such as VHD.
func (vdi) AsyncCopy(session *Session, vdi VDIRef, sr SRRef, baseVdi VDIRef, intoVdi VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.copy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	baseVdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "base_vdi"), baseVdi)
	if err != nil {
		return
	}
	intoVdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "into_vdi"), intoVdi)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, srArg, baseVdiArg, intoVdiArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Copy5: Copy either a full VDI or the block differences between two VDIs into either a fresh VDI or an existing VDI.
// Version: clearwater-felton
//
// Errors:
// VDI_READONLY - The operation required write access but this VDI is read-only
// VDI_TOO_SMALL - The VDI is too small. Please resize it to at least the minimum size.
// VDI_NOT_SPARSE - The VDI is not stored using a sparse format. It is not possible to query and manipulate only the changed blocks (or &apos;block differences&apos; or &apos;disk deltas&apos;) between two VDIs. Please select a VDI which uses a sparse-aware technology such as VHD.
func (vdi) Copy5(session *Session, vdi VDIRef, sr SRRef, baseVdi VDIRef, intoVdi VDIRef) (retval VDIRef, err error) {
	method := "VDI.copy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	baseVdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "base_vdi"), baseVdi)
	if err != nil {
		return
	}
	intoVdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "into_vdi"), intoVdi)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, srArg, baseVdiArg, intoVdiArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncCopy5: Copy either a full VDI or the block differences between two VDIs into either a fresh VDI or an existing VDI.
// Version: clearwater-felton
//
// Errors:
// VDI_READONLY - The operation required write access but this VDI is read-only
// VDI_TOO_SMALL - The VDI is too small. Please resize it to at least the minimum size.
// VDI_NOT_SPARSE - The VDI is not stored using a sparse format. It is not possible to query and manipulate only the changed blocks (or &apos;block differences&apos; or &apos;disk deltas&apos;) between two VDIs. Please select a VDI which uses a sparse-aware technology such as VHD.
func (vdi) AsyncCopy5(session *Session, vdi VDIRef, sr SRRef, baseVdi VDIRef, intoVdi VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.copy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	baseVdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "base_vdi"), baseVdi)
	if err != nil {
		return
	}
	intoVdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "into_vdi"), intoVdi)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, srArg, baseVdiArg, intoVdiArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Copy3: Copy either a full VDI or the block differences between two VDIs into either a fresh VDI or an existing VDI.
// Version: rio
//
// Errors:
// VDI_READONLY - The operation required write access but this VDI is read-only
// VDI_TOO_SMALL - The VDI is too small. Please resize it to at least the minimum size.
// VDI_NOT_SPARSE - The VDI is not stored using a sparse format. It is not possible to query and manipulate only the changed blocks (or &apos;block differences&apos; or &apos;disk deltas&apos;) between two VDIs. Please select a VDI which uses a sparse-aware technology such as VHD.
func (vdi) Copy3(session *Session, vdi VDIRef, sr SRRef) (retval VDIRef, err error) {
	method := "VDI.copy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, srArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncCopy3: Copy either a full VDI or the block differences between two VDIs into either a fresh VDI or an existing VDI.
// Version: rio
//
// Errors:
// VDI_READONLY - The operation required write access but this VDI is read-only
// VDI_TOO_SMALL - The VDI is too small. Please resize it to at least the minimum size.
// VDI_NOT_SPARSE - The VDI is not stored using a sparse format. It is not possible to query and manipulate only the changed blocks (or &apos;block differences&apos; or &apos;disk deltas&apos;) between two VDIs. Please select a VDI which uses a sparse-aware technology such as VHD.
func (vdi) AsyncCopy3(session *Session, vdi VDIRef, sr SRRef) (retval TaskRef, err error) {
	method := "Async.VDI.copy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, srArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Update: Ask the storage backend to refresh the fields in the VDI object
// Version: symc
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
func (vdi) Update(session *Session, vdi VDIRef) (err error) {
	method := "VDI.update"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vdiArg)
	return
}

// AsyncUpdate: Ask the storage backend to refresh the fields in the VDI object
// Version: symc
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
func (vdi) AsyncUpdate(session *Session, vdi VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.update"
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

// Update2: Ask the storage backend to refresh the fields in the VDI object
// Version: symc
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
func (vdi) Update2(session *Session, vdi VDIRef) (err error) {
	method := "VDI.update"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vdiArg)
	return
}

// AsyncUpdate2: Ask the storage backend to refresh the fields in the VDI object
// Version: symc
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
func (vdi) AsyncUpdate2(session *Session, vdi VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.update"
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

// Introduce: Create a new VDI record in the database only
// Version: tampa
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
func (vdi) Introduce(session *Session, uuid string, nameLabel string, nameDescription string, sr SRRef, typeKey VdiType, sharable bool, readOnly bool, otherConfig map[string]string, location string, xenstoreData map[string]string, smConfig map[string]string, managed bool, virtualSize int, physicalUtilisation int, metadataOfPool PoolRef, isASnapshot bool, snapshotTime time.Time, snapshotOf VDIRef) (retval VDIRef, err error) {
	method := "VDI.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "SR"), sr)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeEnumVdiType(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	sharableArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "sharable"), sharable)
	if err != nil {
		return
	}
	readOnlyArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "read_only"), readOnly)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	locationArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "location"), location)
	if err != nil {
		return
	}
	xenstoreDataArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "xenstore_data"), xenstoreData)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	virtualSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "virtual_size"), virtualSize)
	if err != nil {
		return
	}
	physicalUtilisationArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_utilisation"), physicalUtilisation)
	if err != nil {
		return
	}
	metadataOfPoolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "metadata_of_pool"), metadataOfPool)
	if err != nil {
		return
	}
	isASnapshotArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "is_a_snapshot"), isASnapshot)
	if err != nil {
		return
	}
	snapshotTimeArg, err := serializeTime(fmt.Sprintf("%s(%s)", method, "snapshot_time"), snapshotTime)
	if err != nil {
		return
	}
	snapshotOfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "snapshot_of"), snapshotOf)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, srArg, typeKeyArg, sharableArg, readOnlyArg, otherConfigArg, locationArg, xenstoreDataArg, smConfigArg, managedArg, virtualSizeArg, physicalUtilisationArg, metadataOfPoolArg, isASnapshotArg, snapshotTimeArg, snapshotOfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncIntroduce: Create a new VDI record in the database only
// Version: tampa
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
func (vdi) AsyncIntroduce(session *Session, uuid string, nameLabel string, nameDescription string, sr SRRef, typeKey VdiType, sharable bool, readOnly bool, otherConfig map[string]string, location string, xenstoreData map[string]string, smConfig map[string]string, managed bool, virtualSize int, physicalUtilisation int, metadataOfPool PoolRef, isASnapshot bool, snapshotTime time.Time, snapshotOf VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "SR"), sr)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeEnumVdiType(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	sharableArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "sharable"), sharable)
	if err != nil {
		return
	}
	readOnlyArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "read_only"), readOnly)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	locationArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "location"), location)
	if err != nil {
		return
	}
	xenstoreDataArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "xenstore_data"), xenstoreData)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	virtualSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "virtual_size"), virtualSize)
	if err != nil {
		return
	}
	physicalUtilisationArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_utilisation"), physicalUtilisation)
	if err != nil {
		return
	}
	metadataOfPoolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "metadata_of_pool"), metadataOfPool)
	if err != nil {
		return
	}
	isASnapshotArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "is_a_snapshot"), isASnapshot)
	if err != nil {
		return
	}
	snapshotTimeArg, err := serializeTime(fmt.Sprintf("%s(%s)", method, "snapshot_time"), snapshotTime)
	if err != nil {
		return
	}
	snapshotOfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "snapshot_of"), snapshotOf)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, srArg, typeKeyArg, sharableArg, readOnlyArg, otherConfigArg, locationArg, xenstoreDataArg, smConfigArg, managedArg, virtualSizeArg, physicalUtilisationArg, metadataOfPoolArg, isASnapshotArg, snapshotTimeArg, snapshotOfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce19: Create a new VDI record in the database only
// Version: tampa
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
func (vdi) Introduce19(session *Session, uuid string, nameLabel string, nameDescription string, sr SRRef, typeKey VdiType, sharable bool, readOnly bool, otherConfig map[string]string, location string, xenstoreData map[string]string, smConfig map[string]string, managed bool, virtualSize int, physicalUtilisation int, metadataOfPool PoolRef, isASnapshot bool, snapshotTime time.Time, snapshotOf VDIRef) (retval VDIRef, err error) {
	method := "VDI.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "SR"), sr)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeEnumVdiType(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	sharableArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "sharable"), sharable)
	if err != nil {
		return
	}
	readOnlyArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "read_only"), readOnly)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	locationArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "location"), location)
	if err != nil {
		return
	}
	xenstoreDataArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "xenstore_data"), xenstoreData)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	virtualSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "virtual_size"), virtualSize)
	if err != nil {
		return
	}
	physicalUtilisationArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_utilisation"), physicalUtilisation)
	if err != nil {
		return
	}
	metadataOfPoolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "metadata_of_pool"), metadataOfPool)
	if err != nil {
		return
	}
	isASnapshotArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "is_a_snapshot"), isASnapshot)
	if err != nil {
		return
	}
	snapshotTimeArg, err := serializeTime(fmt.Sprintf("%s(%s)", method, "snapshot_time"), snapshotTime)
	if err != nil {
		return
	}
	snapshotOfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "snapshot_of"), snapshotOf)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, srArg, typeKeyArg, sharableArg, readOnlyArg, otherConfigArg, locationArg, xenstoreDataArg, smConfigArg, managedArg, virtualSizeArg, physicalUtilisationArg, metadataOfPoolArg, isASnapshotArg, snapshotTimeArg, snapshotOfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncIntroduce19: Create a new VDI record in the database only
// Version: tampa
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
func (vdi) AsyncIntroduce19(session *Session, uuid string, nameLabel string, nameDescription string, sr SRRef, typeKey VdiType, sharable bool, readOnly bool, otherConfig map[string]string, location string, xenstoreData map[string]string, smConfig map[string]string, managed bool, virtualSize int, physicalUtilisation int, metadataOfPool PoolRef, isASnapshot bool, snapshotTime time.Time, snapshotOf VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "SR"), sr)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeEnumVdiType(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	sharableArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "sharable"), sharable)
	if err != nil {
		return
	}
	readOnlyArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "read_only"), readOnly)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	locationArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "location"), location)
	if err != nil {
		return
	}
	xenstoreDataArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "xenstore_data"), xenstoreData)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	virtualSizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "virtual_size"), virtualSize)
	if err != nil {
		return
	}
	physicalUtilisationArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "physical_utilisation"), physicalUtilisation)
	if err != nil {
		return
	}
	metadataOfPoolArg, err := serializePoolRef(fmt.Sprintf("%s(%s)", method, "metadata_of_pool"), metadataOfPool)
	if err != nil {
		return
	}
	isASnapshotArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "is_a_snapshot"), isASnapshot)
	if err != nil {
		return
	}
	snapshotTimeArg, err := serializeTime(fmt.Sprintf("%s(%s)", method, "snapshot_time"), snapshotTime)
	if err != nil {
		return
	}
	snapshotOfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "snapshot_of"), snapshotOf)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, srArg, typeKeyArg, sharableArg, readOnlyArg, otherConfigArg, locationArg, xenstoreDataArg, smConfigArg, managedArg, virtualSizeArg, physicalUtilisationArg, metadataOfPoolArg, isASnapshotArg, snapshotTimeArg, snapshotOfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce12: Create a new VDI record in the database only
// Version: miami
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
func (vdi) Introduce12(session *Session, uuid string, nameLabel string, nameDescription string, sr SRRef, typeKey VdiType, sharable bool, readOnly bool, otherConfig map[string]string, location string, xenstoreData map[string]string, smConfig map[string]string) (retval VDIRef, err error) {
	method := "VDI.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "SR"), sr)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeEnumVdiType(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	sharableArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "sharable"), sharable)
	if err != nil {
		return
	}
	readOnlyArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "read_only"), readOnly)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	locationArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "location"), location)
	if err != nil {
		return
	}
	xenstoreDataArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "xenstore_data"), xenstoreData)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, srArg, typeKeyArg, sharableArg, readOnlyArg, otherConfigArg, locationArg, xenstoreDataArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncIntroduce12: Create a new VDI record in the database only
// Version: miami
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
func (vdi) AsyncIntroduce12(session *Session, uuid string, nameLabel string, nameDescription string, sr SRRef, typeKey VdiType, sharable bool, readOnly bool, otherConfig map[string]string, location string, xenstoreData map[string]string, smConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.VDI.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "SR"), sr)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeEnumVdiType(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	sharableArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "sharable"), sharable)
	if err != nil {
		return
	}
	readOnlyArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "read_only"), readOnly)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	locationArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "location"), location)
	if err != nil {
		return
	}
	xenstoreDataArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "xenstore_data"), xenstoreData)
	if err != nil {
		return
	}
	smConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "sm_config"), smConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, srArg, typeKeyArg, sharableArg, readOnlyArg, otherConfigArg, locationArg, xenstoreDataArg, smConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce11: Create a new VDI record in the database only
// Version: rio
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
func (vdi) Introduce11(session *Session, uuid string, nameLabel string, nameDescription string, sr SRRef, typeKey VdiType, sharable bool, readOnly bool, otherConfig map[string]string, location string, xenstoreData map[string]string) (retval VDIRef, err error) {
	method := "VDI.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "SR"), sr)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeEnumVdiType(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	sharableArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "sharable"), sharable)
	if err != nil {
		return
	}
	readOnlyArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "read_only"), readOnly)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	locationArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "location"), location)
	if err != nil {
		return
	}
	xenstoreDataArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "xenstore_data"), xenstoreData)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, srArg, typeKeyArg, sharableArg, readOnlyArg, otherConfigArg, locationArg, xenstoreDataArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncIntroduce11: Create a new VDI record in the database only
// Version: rio
//
// Errors:
// SR_OPERATION_NOT_SUPPORTED - The SR backend does not support the operation (check the SR&apos;s allowed operations)
func (vdi) AsyncIntroduce11(session *Session, uuid string, nameLabel string, nameDescription string, sr SRRef, typeKey VdiType, sharable bool, readOnly bool, otherConfig map[string]string, location string, xenstoreData map[string]string) (retval TaskRef, err error) {
	method := "Async.VDI.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	uuidArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "uuid"), uuid)
	if err != nil {
		return
	}
	nameLabelArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_label"), nameLabel)
	if err != nil {
		return
	}
	nameDescriptionArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name_description"), nameDescription)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "SR"), sr)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeEnumVdiType(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	sharableArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "sharable"), sharable)
	if err != nil {
		return
	}
	readOnlyArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "read_only"), readOnly)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	locationArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "location"), location)
	if err != nil {
		return
	}
	xenstoreDataArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "xenstore_data"), xenstoreData)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, uuidArg, nameLabelArg, nameDescriptionArg, srArg, typeKeyArg, sharableArg, readOnlyArg, otherConfigArg, locationArg, xenstoreDataArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ResizeOnline: Resize the VDI which may or may not be attached to running guests.
// Version: rio
func (vdi) ResizeOnline(session *Session, vdi VDIRef, size int) (err error) {
	method := "VDI.resize_online"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	sizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "size"), size)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vdiArg, sizeArg)
	return
}

// AsyncResizeOnline: Resize the VDI which may or may not be attached to running guests.
// Version: rio
func (vdi) AsyncResizeOnline(session *Session, vdi VDIRef, size int) (retval TaskRef, err error) {
	method := "Async.VDI.resize_online"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	sizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "size"), size)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, sizeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ResizeOnline3: Resize the VDI which may or may not be attached to running guests.
// Version: rio
func (vdi) ResizeOnline3(session *Session, vdi VDIRef, size int) (err error) {
	method := "VDI.resize_online"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	sizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "size"), size)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vdiArg, sizeArg)
	return
}

// AsyncResizeOnline3: Resize the VDI which may or may not be attached to running guests.
// Version: rio
func (vdi) AsyncResizeOnline3(session *Session, vdi VDIRef, size int) (retval TaskRef, err error) {
	method := "Async.VDI.resize_online"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	sizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "size"), size)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, sizeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Resize: Resize the VDI.
// Version: rio
func (vdi) Resize(session *Session, vdi VDIRef, size int) (err error) {
	method := "VDI.resize"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	sizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "size"), size)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vdiArg, sizeArg)
	return
}

// AsyncResize: Resize the VDI.
// Version: rio
func (vdi) AsyncResize(session *Session, vdi VDIRef, size int) (retval TaskRef, err error) {
	method := "Async.VDI.resize"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	sizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "size"), size)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, sizeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Resize3: Resize the VDI.
// Version: rio
func (vdi) Resize3(session *Session, vdi VDIRef, size int) (err error) {
	method := "VDI.resize"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	sizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "size"), size)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vdiArg, sizeArg)
	return
}

// AsyncResize3: Resize the VDI.
// Version: rio
func (vdi) AsyncResize3(session *Session, vdi VDIRef, size int) (retval TaskRef, err error) {
	method := "Async.VDI.resize"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	sizeArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "size"), size)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, sizeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Clone: Take an exact copy of the VDI and return a reference to the new disk. If any driver_params are specified then these are passed through to the storage-specific substrate driver that implements the clone operation. NB the clone lives in the same Storage Repository as its parent.
// Version: miami
func (vdi) Clone(session *Session, vdi VDIRef, driverParams map[string]string) (retval VDIRef, err error) {
	method := "VDI.clone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	driverParamsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "driver_params"), driverParams)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, driverParamsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncClone: Take an exact copy of the VDI and return a reference to the new disk. If any driver_params are specified then these are passed through to the storage-specific substrate driver that implements the clone operation. NB the clone lives in the same Storage Repository as its parent.
// Version: miami
func (vdi) AsyncClone(session *Session, vdi VDIRef, driverParams map[string]string) (retval TaskRef, err error) {
	method := "Async.VDI.clone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	driverParamsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "driver_params"), driverParams)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, driverParamsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Clone3: Take an exact copy of the VDI and return a reference to the new disk. If any driver_params are specified then these are passed through to the storage-specific substrate driver that implements the clone operation. NB the clone lives in the same Storage Repository as its parent.
// Version: miami
func (vdi) Clone3(session *Session, vdi VDIRef, driverParams map[string]string) (retval VDIRef, err error) {
	method := "VDI.clone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	driverParamsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "driver_params"), driverParams)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, driverParamsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncClone3: Take an exact copy of the VDI and return a reference to the new disk. If any driver_params are specified then these are passed through to the storage-specific substrate driver that implements the clone operation. NB the clone lives in the same Storage Repository as its parent.
// Version: miami
func (vdi) AsyncClone3(session *Session, vdi VDIRef, driverParams map[string]string) (retval TaskRef, err error) {
	method := "Async.VDI.clone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	driverParamsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "driver_params"), driverParams)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, driverParamsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Clone2: Take an exact copy of the VDI and return a reference to the new disk. If any driver_params are specified then these are passed through to the storage-specific substrate driver that implements the clone operation. NB the clone lives in the same Storage Repository as its parent.
// Version: rio
func (vdi) Clone2(session *Session, vdi VDIRef) (retval VDIRef, err error) {
	method := "VDI.clone"
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
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncClone2: Take an exact copy of the VDI and return a reference to the new disk. If any driver_params are specified then these are passed through to the storage-specific substrate driver that implements the clone operation. NB the clone lives in the same Storage Repository as its parent.
// Version: rio
func (vdi) AsyncClone2(session *Session, vdi VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.clone"
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

// Snapshot: Take a read-only snapshot of the VDI, returning a reference to the snapshot. If any driver_params are specified then these are passed through to the storage-specific substrate driver that takes the snapshot. NB the snapshot lives in the same Storage Repository as its parent.
// Version: miami
func (vdi) Snapshot(session *Session, vdi VDIRef, driverParams map[string]string) (retval VDIRef, err error) {
	method := "VDI.snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	driverParamsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "driver_params"), driverParams)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, driverParamsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncSnapshot: Take a read-only snapshot of the VDI, returning a reference to the snapshot. If any driver_params are specified then these are passed through to the storage-specific substrate driver that takes the snapshot. NB the snapshot lives in the same Storage Repository as its parent.
// Version: miami
func (vdi) AsyncSnapshot(session *Session, vdi VDIRef, driverParams map[string]string) (retval TaskRef, err error) {
	method := "Async.VDI.snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	driverParamsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "driver_params"), driverParams)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, driverParamsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Snapshot3: Take a read-only snapshot of the VDI, returning a reference to the snapshot. If any driver_params are specified then these are passed through to the storage-specific substrate driver that takes the snapshot. NB the snapshot lives in the same Storage Repository as its parent.
// Version: miami
func (vdi) Snapshot3(session *Session, vdi VDIRef, driverParams map[string]string) (retval VDIRef, err error) {
	method := "VDI.snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	driverParamsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "driver_params"), driverParams)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, driverParamsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncSnapshot3: Take a read-only snapshot of the VDI, returning a reference to the snapshot. If any driver_params are specified then these are passed through to the storage-specific substrate driver that takes the snapshot. NB the snapshot lives in the same Storage Repository as its parent.
// Version: miami
func (vdi) AsyncSnapshot3(session *Session, vdi VDIRef, driverParams map[string]string) (retval TaskRef, err error) {
	method := "Async.VDI.snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vdiArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "vdi"), vdi)
	if err != nil {
		return
	}
	driverParamsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "driver_params"), driverParams)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vdiArg, driverParamsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Snapshot2: Take a read-only snapshot of the VDI, returning a reference to the snapshot. If any driver_params are specified then these are passed through to the storage-specific substrate driver that takes the snapshot. NB the snapshot lives in the same Storage Repository as its parent.
// Version: rio
func (vdi) Snapshot2(session *Session, vdi VDIRef) (retval VDIRef, err error) {
	method := "VDI.snapshot"
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
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncSnapshot2: Take a read-only snapshot of the VDI, returning a reference to the snapshot. If any driver_params are specified then these are passed through to the storage-specific substrate driver that takes the snapshot. NB the snapshot lives in the same Storage Repository as its parent.
// Version: rio
func (vdi) AsyncSnapshot2(session *Session, vdi VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.snapshot"
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

// RemoveTags: Remove the given value from the tags field of the given VDI.  If the value is not in that Set, then do nothing.
// Version: orlando
func (vdi) RemoveTags(session *Session, self VDIRef, value string) (err error) {
	method := "VDI.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// RemoveTags3: Remove the given value from the tags field of the given VDI.  If the value is not in that Set, then do nothing.
// Version: orlando
func (vdi) RemoveTags3(session *Session, self VDIRef, value string) (err error) {
	method := "VDI.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// RemoveTags2: Remove the given value from the tags field of the given VDI.  If the value is not in that Set, then do nothing.
// Version: rio
func (vdi) RemoveTags2(session *Session, self VDIRef) (err error) {
	method := "VDI.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddTags: Add the given value to the tags field of the given VDI.  If the value is already in that Set, then do nothing.
// Version: orlando
func (vdi) AddTags(session *Session, self VDIRef, value string) (err error) {
	method := "VDI.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AddTags3: Add the given value to the tags field of the given VDI.  If the value is already in that Set, then do nothing.
// Version: orlando
func (vdi) AddTags3(session *Session, self VDIRef, value string) (err error) {
	method := "VDI.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AddTags2: Add the given value to the tags field of the given VDI.  If the value is already in that Set, then do nothing.
// Version: rio
func (vdi) AddTags2(session *Session, self VDIRef) (err error) {
	method := "VDI.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetTags: Set the tags field of the given VDI.
// Version: orlando
func (vdi) SetTags(session *Session, self VDIRef, value []string) (err error) {
	method := "VDI.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetTags3: Set the tags field of the given VDI.
// Version: orlando
func (vdi) SetTags3(session *Session, self VDIRef, value []string) (err error) {
	method := "VDI.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetTags2: Set the tags field of the given VDI.
// Version: rio
func (vdi) SetTags2(session *Session, self VDIRef) (err error) {
	method := "VDI.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveFromSmConfig: Remove the given key and its corresponding value from the sm_config field of the given VDI.  If the key is not in that Map, then do nothing.
// Version: miami
func (vdi) RemoveFromSmConfig(session *Session, self VDIRef, key string) (err error) {
	method := "VDI.remove_from_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromSmConfig3: Remove the given key and its corresponding value from the sm_config field of the given VDI.  If the key is not in that Map, then do nothing.
// Version: miami
func (vdi) RemoveFromSmConfig3(session *Session, self VDIRef, key string) (err error) {
	method := "VDI.remove_from_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromSmConfig2: Remove the given key and its corresponding value from the sm_config field of the given VDI.  If the key is not in that Map, then do nothing.
// Version: rio
func (vdi) RemoveFromSmConfig2(session *Session, self VDIRef) (err error) {
	method := "VDI.remove_from_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToSmConfig: Add the given key-value pair to the sm_config field of the given VDI.
// Version: miami
func (vdi) AddToSmConfig(session *Session, self VDIRef, key string, value string) (err error) {
	method := "VDI.add_to_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToSmConfig4: Add the given key-value pair to the sm_config field of the given VDI.
// Version: miami
func (vdi) AddToSmConfig4(session *Session, self VDIRef, key string, value string) (err error) {
	method := "VDI.add_to_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToSmConfig2: Add the given key-value pair to the sm_config field of the given VDI.
// Version: rio
func (vdi) AddToSmConfig2(session *Session, self VDIRef) (err error) {
	method := "VDI.add_to_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetSmConfig: Set the sm_config field of the given VDI.
// Version: miami
func (vdi) SetSmConfig(session *Session, self VDIRef, value map[string]string) (err error) {
	method := "VDI.set_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetSmConfig3: Set the sm_config field of the given VDI.
// Version: miami
func (vdi) SetSmConfig3(session *Session, self VDIRef, value map[string]string) (err error) {
	method := "VDI.set_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetSmConfig2: Set the sm_config field of the given VDI.
// Version: rio
func (vdi) SetSmConfig2(session *Session, self VDIRef) (err error) {
	method := "VDI.set_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveFromXenstoreData: Remove the given key and its corresponding value from the xenstore_data field of the given VDI.  If the key is not in that Map, then do nothing.
// Version: miami
func (vdi) RemoveFromXenstoreData(session *Session, self VDIRef, key string) (err error) {
	method := "VDI.remove_from_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromXenstoreData3: Remove the given key and its corresponding value from the xenstore_data field of the given VDI.  If the key is not in that Map, then do nothing.
// Version: miami
func (vdi) RemoveFromXenstoreData3(session *Session, self VDIRef, key string) (err error) {
	method := "VDI.remove_from_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromXenstoreData2: Remove the given key and its corresponding value from the xenstore_data field of the given VDI.  If the key is not in that Map, then do nothing.
// Version: rio
func (vdi) RemoveFromXenstoreData2(session *Session, self VDIRef) (err error) {
	method := "VDI.remove_from_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToXenstoreData: Add the given key-value pair to the xenstore_data field of the given VDI.
// Version: miami
func (vdi) AddToXenstoreData(session *Session, self VDIRef, key string, value string) (err error) {
	method := "VDI.add_to_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToXenstoreData4: Add the given key-value pair to the xenstore_data field of the given VDI.
// Version: miami
func (vdi) AddToXenstoreData4(session *Session, self VDIRef, key string, value string) (err error) {
	method := "VDI.add_to_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToXenstoreData2: Add the given key-value pair to the xenstore_data field of the given VDI.
// Version: rio
func (vdi) AddToXenstoreData2(session *Session, self VDIRef) (err error) {
	method := "VDI.add_to_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetXenstoreData: Set the xenstore_data field of the given VDI.
// Version: miami
func (vdi) SetXenstoreData(session *Session, self VDIRef, value map[string]string) (err error) {
	method := "VDI.set_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetXenstoreData3: Set the xenstore_data field of the given VDI.
// Version: miami
func (vdi) SetXenstoreData3(session *Session, self VDIRef, value map[string]string) (err error) {
	method := "VDI.set_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetXenstoreData2: Set the xenstore_data field of the given VDI.
// Version: rio
func (vdi) SetXenstoreData2(session *Session, self VDIRef) (err error) {
	method := "VDI.set_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given VDI.  If the key is not in that Map, then do nothing.
// Version: rio
func (vdi) RemoveFromOtherConfig(session *Session, self VDIRef, key string) (err error) {
	method := "VDI.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given VDI.  If the key is not in that Map, then do nothing.
// Version: rio
func (vdi) RemoveFromOtherConfig3(session *Session, self VDIRef, key string) (err error) {
	method := "VDI.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given VDI.
// Version: rio
func (vdi) AddToOtherConfig(session *Session, self VDIRef, key string, value string) (err error) {
	method := "VDI.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given VDI.
// Version: rio
func (vdi) AddToOtherConfig4(session *Session, self VDIRef, key string, value string) (err error) {
	method := "VDI.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig: Set the other_config field of the given VDI.
// Version: rio
func (vdi) SetOtherConfig(session *Session, self VDIRef, value map[string]string) (err error) {
	method := "VDI.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given VDI.
// Version: rio
func (vdi) SetOtherConfig3(session *Session, self VDIRef, value map[string]string) (err error) {
	method := "VDI.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCbtEnabled: Get the cbt_enabled field of the given VDI.
// Version: rio
func (vdi) GetCbtEnabled(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_cbt_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCbtEnabled2: Get the cbt_enabled field of the given VDI.
// Version: rio
func (vdi) GetCbtEnabled2(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_cbt_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsToolsIso: Get the is_tools_iso field of the given VDI.
// Version: rio
func (vdi) GetIsToolsIso(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_is_tools_iso"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsToolsIso2: Get the is_tools_iso field of the given VDI.
// Version: rio
func (vdi) GetIsToolsIso2(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_is_tools_iso"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMetadataLatest: Get the metadata_latest field of the given VDI.
// Version: rio
func (vdi) GetMetadataLatest(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_metadata_latest"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMetadataLatest2: Get the metadata_latest field of the given VDI.
// Version: rio
func (vdi) GetMetadataLatest2(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_metadata_latest"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMetadataOfPool: Get the metadata_of_pool field of the given VDI.
// Version: rio
func (vdi) GetMetadataOfPool(session *Session, self VDIRef) (retval PoolRef, err error) {
	method := "VDI.get_metadata_of_pool"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolRef(method+" -> ", result)
	return
}

// GetMetadataOfPool2: Get the metadata_of_pool field of the given VDI.
// Version: rio
func (vdi) GetMetadataOfPool2(session *Session, self VDIRef) (retval PoolRef, err error) {
	method := "VDI.get_metadata_of_pool"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePoolRef(method+" -> ", result)
	return
}

// GetOnBoot: Get the on_boot field of the given VDI.
// Version: rio
func (vdi) GetOnBoot(session *Session, self VDIRef) (retval OnBoot, err error) {
	method := "VDI.get_on_boot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumOnBoot(method+" -> ", result)
	return
}

// GetOnBoot2: Get the on_boot field of the given VDI.
// Version: rio
func (vdi) GetOnBoot2(session *Session, self VDIRef) (retval OnBoot, err error) {
	method := "VDI.get_on_boot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumOnBoot(method+" -> ", result)
	return
}

// GetAllowCaching: Get the allow_caching field of the given VDI.
// Version: rio
func (vdi) GetAllowCaching(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_allow_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAllowCaching2: Get the allow_caching field of the given VDI.
// Version: rio
func (vdi) GetAllowCaching2(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_allow_caching"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTags: Get the tags field of the given VDI.
// Version: rio
func (vdi) GetTags(session *Session, self VDIRef) (retval []string, err error) {
	method := "VDI.get_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetTags2: Get the tags field of the given VDI.
// Version: rio
func (vdi) GetTags2(session *Session, self VDIRef) (retval []string, err error) {
	method := "VDI.get_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetSnapshotTime: Get the snapshot_time field of the given VDI.
// Version: rio
func (vdi) GetSnapshotTime(session *Session, self VDIRef) (retval time.Time, err error) {
	method := "VDI.get_snapshot_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSnapshotTime2: Get the snapshot_time field of the given VDI.
// Version: rio
func (vdi) GetSnapshotTime2(session *Session, self VDIRef) (retval time.Time, err error) {
	method := "VDI.get_snapshot_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSnapshots: Get the snapshots field of the given VDI.
// Version: rio
func (vdi) GetSnapshots(session *Session, self VDIRef) (retval []VDIRef, err error) {
	method := "VDI.get_snapshots"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRefSet(method+" -> ", result)
	return
}

// GetSnapshots2: Get the snapshots field of the given VDI.
// Version: rio
func (vdi) GetSnapshots2(session *Session, self VDIRef) (retval []VDIRef, err error) {
	method := "VDI.get_snapshots"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRefSet(method+" -> ", result)
	return
}

// GetSnapshotOf: Get the snapshot_of field of the given VDI.
// Version: rio
func (vdi) GetSnapshotOf(session *Session, self VDIRef) (retval VDIRef, err error) {
	method := "VDI.get_snapshot_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSnapshotOf2: Get the snapshot_of field of the given VDI.
// Version: rio
func (vdi) GetSnapshotOf2(session *Session, self VDIRef) (retval VDIRef, err error) {
	method := "VDI.get_snapshot_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsASnapshot: Get the is_a_snapshot field of the given VDI.
// Version: rio
func (vdi) GetIsASnapshot(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_is_a_snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsASnapshot2: Get the is_a_snapshot field of the given VDI.
// Version: rio
func (vdi) GetIsASnapshot2(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_is_a_snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSmConfig: Get the sm_config field of the given VDI.
// Version: rio
func (vdi) GetSmConfig(session *Session, self VDIRef) (retval map[string]string, err error) {
	method := "VDI.get_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSmConfig2: Get the sm_config field of the given VDI.
// Version: rio
func (vdi) GetSmConfig2(session *Session, self VDIRef) (retval map[string]string, err error) {
	method := "VDI.get_sm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetXenstoreData: Get the xenstore_data field of the given VDI.
// Version: rio
func (vdi) GetXenstoreData(session *Session, self VDIRef) (retval map[string]string, err error) {
	method := "VDI.get_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetXenstoreData2: Get the xenstore_data field of the given VDI.
// Version: rio
func (vdi) GetXenstoreData2(session *Session, self VDIRef) (retval map[string]string, err error) {
	method := "VDI.get_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetParent: Get the parent field of the given VDI.
// Version: rio
func (vdi) GetParent(session *Session, self VDIRef) (retval VDIRef, err error) {
	method := "VDI.get_parent"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetParent2: Get the parent field of the given VDI.
// Version: rio
func (vdi) GetParent2(session *Session, self VDIRef) (retval VDIRef, err error) {
	method := "VDI.get_parent"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMissing: Get the missing field of the given VDI.
// Version: rio
func (vdi) GetMissing(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_missing"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMissing2: Get the missing field of the given VDI.
// Version: rio
func (vdi) GetMissing2(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_missing"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetManaged: Get the managed field of the given VDI.
// Version: rio
func (vdi) GetManaged(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_managed"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetManaged2: Get the managed field of the given VDI.
// Version: rio
func (vdi) GetManaged2(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_managed"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLocation: Get the location field of the given VDI.
// Version: rio
func (vdi) GetLocation(session *Session, self VDIRef) (retval string, err error) {
	method := "VDI.get_location"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLocation2: Get the location field of the given VDI.
// Version: rio
func (vdi) GetLocation2(session *Session, self VDIRef) (retval string, err error) {
	method := "VDI.get_location"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetStorageLock: Get the storage_lock field of the given VDI.
// Version: rio
func (vdi) GetStorageLock(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_storage_lock"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetStorageLock2: Get the storage_lock field of the given VDI.
// Version: rio
func (vdi) GetStorageLock2(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_storage_lock"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given VDI.
// Version: rio
func (vdi) GetOtherConfig(session *Session, self VDIRef) (retval map[string]string, err error) {
	method := "VDI.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given VDI.
// Version: rio
func (vdi) GetOtherConfig2(session *Session, self VDIRef) (retval map[string]string, err error) {
	method := "VDI.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetReadOnly: Get the read_only field of the given VDI.
// Version: rio
func (vdi) GetReadOnly(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_read_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetReadOnly2: Get the read_only field of the given VDI.
// Version: rio
func (vdi) GetReadOnly2(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_read_only"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSharable: Get the sharable field of the given VDI.
// Version: rio
func (vdi) GetSharable(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_sharable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSharable2: Get the sharable field of the given VDI.
// Version: rio
func (vdi) GetSharable2(session *Session, self VDIRef) (retval bool, err error) {
	method := "VDI.get_sharable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetType: Get the type field of the given VDI.
// Version: rio
func (vdi) GetType(session *Session, self VDIRef) (retval VdiType, err error) {
	method := "VDI.get_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVdiType(method+" -> ", result)
	return
}

// GetType2: Get the type field of the given VDI.
// Version: rio
func (vdi) GetType2(session *Session, self VDIRef) (retval VdiType, err error) {
	method := "VDI.get_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVdiType(method+" -> ", result)
	return
}

// GetPhysicalUtilisation: Get the physical_utilisation field of the given VDI.
// Version: rio
func (vdi) GetPhysicalUtilisation(session *Session, self VDIRef) (retval int, err error) {
	method := "VDI.get_physical_utilisation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPhysicalUtilisation2: Get the physical_utilisation field of the given VDI.
// Version: rio
func (vdi) GetPhysicalUtilisation2(session *Session, self VDIRef) (retval int, err error) {
	method := "VDI.get_physical_utilisation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVirtualSize: Get the virtual_size field of the given VDI.
// Version: rio
func (vdi) GetVirtualSize(session *Session, self VDIRef) (retval int, err error) {
	method := "VDI.get_virtual_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVirtualSize2: Get the virtual_size field of the given VDI.
// Version: rio
func (vdi) GetVirtualSize2(session *Session, self VDIRef) (retval int, err error) {
	method := "VDI.get_virtual_size"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCrashDumps: Get the crash_dumps field of the given VDI.
// Version: rio
func (vdi) GetCrashDumps(session *Session, self VDIRef) (retval []CrashdumpRef, err error) {
	method := "VDI.get_crash_dumps"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeCrashdumpRefSet(method+" -> ", result)
	return
}

// GetCrashDumps2: Get the crash_dumps field of the given VDI.
// Version: rio
func (vdi) GetCrashDumps2(session *Session, self VDIRef) (retval []CrashdumpRef, err error) {
	method := "VDI.get_crash_dumps"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeCrashdumpRefSet(method+" -> ", result)
	return
}

// GetVBDs: Get the VBDs field of the given VDI.
// Version: rio
func (vdi) GetVBDs(session *Session, self VDIRef) (retval []VBDRef, err error) {
	method := "VDI.get_VBDs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVBDRefSet(method+" -> ", result)
	return
}

// GetVBDs2: Get the VBDs field of the given VDI.
// Version: rio
func (vdi) GetVBDs2(session *Session, self VDIRef) (retval []VBDRef, err error) {
	method := "VDI.get_VBDs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVBDRefSet(method+" -> ", result)
	return
}

// GetSR: Get the SR field of the given VDI.
// Version: rio
func (vdi) GetSR(session *Session, self VDIRef) (retval SRRef, err error) {
	method := "VDI.get_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// GetSR2: Get the SR field of the given VDI.
// Version: rio
func (vdi) GetSR2(session *Session, self VDIRef) (retval SRRef, err error) {
	method := "VDI.get_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRef(method+" -> ", result)
	return
}

// GetCurrentOperations: Get the current_operations field of the given VDI.
// Version: rio
func (vdi) GetCurrentOperations(session *Session, self VDIRef) (retval map[string]VdiOperations, err error) {
	method := "VDI.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumVdiOperationsMap(method+" -> ", result)
	return
}

// GetCurrentOperations2: Get the current_operations field of the given VDI.
// Version: rio
func (vdi) GetCurrentOperations2(session *Session, self VDIRef) (retval map[string]VdiOperations, err error) {
	method := "VDI.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumVdiOperationsMap(method+" -> ", result)
	return
}

// GetAllowedOperations: Get the allowed_operations field of the given VDI.
// Version: rio
func (vdi) GetAllowedOperations(session *Session, self VDIRef) (retval []VdiOperations, err error) {
	method := "VDI.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVdiOperationsSet(method+" -> ", result)
	return
}

// GetAllowedOperations2: Get the allowed_operations field of the given VDI.
// Version: rio
func (vdi) GetAllowedOperations2(session *Session, self VDIRef) (retval []VdiOperations, err error) {
	method := "VDI.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVdiOperationsSet(method+" -> ", result)
	return
}

// GetNameDescription: Get the name/description field of the given VDI.
// Version: rio
func (vdi) GetNameDescription(session *Session, self VDIRef) (retval string, err error) {
	method := "VDI.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription2: Get the name/description field of the given VDI.
// Version: rio
func (vdi) GetNameDescription2(session *Session, self VDIRef) (retval string, err error) {
	method := "VDI.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel: Get the name/label field of the given VDI.
// Version: rio
func (vdi) GetNameLabel(session *Session, self VDIRef) (retval string, err error) {
	method := "VDI.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel2: Get the name/label field of the given VDI.
// Version: rio
func (vdi) GetNameLabel2(session *Session, self VDIRef) (retval string, err error) {
	method := "VDI.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given VDI.
// Version: rio
func (vdi) GetUUID(session *Session, self VDIRef) (retval string, err error) {
	method := "VDI.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given VDI.
// Version: rio
func (vdi) GetUUID2(session *Session, self VDIRef) (retval string, err error) {
	method := "VDI.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByNameLabel: Get all the VDI instances with the given label.
// Version: rio
func (vdi) GetByNameLabel(session *Session, label string) (retval []VDIRef, err error) {
	method := "VDI.get_by_name_label"
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
	retval, err = deserializeVDIRefSet(method+" -> ", result)
	return
}

// GetByNameLabel2: Get all the VDI instances with the given label.
// Version: rio
func (vdi) GetByNameLabel2(session *Session, label string) (retval []VDIRef, err error) {
	method := "VDI.get_by_name_label"
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
	retval, err = deserializeVDIRefSet(method+" -> ", result)
	return
}

// Destroy: Destroy the specified VDI instance.
// Version: rio
func (vdi) Destroy(session *Session, self VDIRef) (err error) {
	method := "VDI.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy the specified VDI instance.
// Version: rio
func (vdi) AsyncDestroy(session *Session, self VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy the specified VDI instance.
// Version: rio
func (vdi) Destroy2(session *Session, self VDIRef) (err error) {
	method := "VDI.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy the specified VDI instance.
// Version: rio
func (vdi) AsyncDestroy2(session *Session, self VDIRef) (retval TaskRef, err error) {
	method := "Async.VDI.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Create a new VDI instance, and return its handle. The constructor args are: name_label, name_description, SR*, virtual_size*, type*, sharable*, read_only*, other_config*, xenstore_data, sm_config, tags (* = non-optional).
// Version: rio
func (vdi) Create(session *Session, args VDIRecord) (retval VDIRef, err error) {
	method := "VDI.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVDIRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncCreate: Create a new VDI instance, and return its handle. The constructor args are: name_label, name_description, SR*, virtual_size*, type*, sharable*, read_only*, other_config*, xenstore_data, sm_config, tags (* = non-optional).
// Version: rio
func (vdi) AsyncCreate(session *Session, args VDIRecord) (retval TaskRef, err error) {
	method := "Async.VDI.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVDIRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create2: Create a new VDI instance, and return its handle. The constructor args are: name_label, name_description, SR*, virtual_size*, type*, sharable*, read_only*, other_config*, xenstore_data, sm_config, tags (* = non-optional).
// Version: rio
func (vdi) Create2(session *Session, args VDIRecord) (retval VDIRef, err error) {
	method := "VDI.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVDIRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// AsyncCreate2: Create a new VDI instance, and return its handle. The constructor args are: name_label, name_description, SR*, virtual_size*, type*, sharable*, read_only*, other_config*, xenstore_data, sm_config, tags (* = non-optional).
// Version: rio
func (vdi) AsyncCreate2(session *Session, args VDIRecord) (retval TaskRef, err error) {
	method := "Async.VDI.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVDIRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetByUUID: Get a reference to the VDI instance with the specified UUID.
// Version: rio
func (vdi) GetByUUID(session *Session, uuid string) (retval VDIRef, err error) {
	method := "VDI.get_by_uuid"
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
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the VDI instance with the specified UUID.
// Version: rio
func (vdi) GetByUUID2(session *Session, uuid string) (retval VDIRef, err error) {
	method := "VDI.get_by_uuid"
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
	retval, err = deserializeVDIRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given VDI.
// Version: rio
func (vdi) GetRecord(session *Session, self VDIRef) (retval VDIRecord, err error) {
	method := "VDI.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given VDI.
// Version: rio
func (vdi) GetRecord2(session *Session, self VDIRef) (retval VDIRecord, err error) {
	method := "VDI.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVDIRecord(method+" -> ", result)
	return
}
