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

type ClusterHostRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// Reference to the Cluster object
	Cluster ClusterRef `json:"cluster,omitempty"`
	// Reference to the Host object
	Host HostRef `json:"host,omitempty"`
	// Whether the cluster host believes that clustering should be enabled on this host. This field can be altered by calling the enable/disable message on a cluster host. Only enabled members run the underlying cluster stack. Disabled members are still considered a member of the cluster (see joined), and can be re-enabled by the user.
	Enabled bool `json:"enabled,omitempty"`
	// Reference to the PIF object
	PIF PIFRef `json:"PIF,omitempty"`
	// Whether the cluster host has joined the cluster. Contrary to enabled, a host that is not joined is not considered a member of the cluster, and hence enable and disable operations cannot be performed on this host.
	Joined bool `json:"joined,omitempty"`
	// Whether the underlying cluster stack thinks we are live. This field is set automatically based on updates from the cluster stack and cannot be altered by the user.
	Live bool `json:"live,omitempty"`
	// Time when the live field was last updated based on information from the cluster stack
	LastUpdateLive time.Time `json:"last_update_live,omitempty"`
	// list of the operations allowed in this state. This list is advisory only and the server state may have changed by the time this field is read by a client.
	AllowedOperations []ClusterHostOperation `json:"allowed_operations,omitempty"`
	// links each of the running tasks using this object (by reference) to a current_operation enum which describes the nature of the task.
	CurrentOperations map[string]ClusterHostOperation `json:"current_operations,omitempty"`
	// Additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
}

type ClusterHostRef string

// Cluster member metadata
type clusterHost struct{}

var ClusterHost clusterHost

// GetAllRecords: Return a map of Cluster_host references to Cluster_host records for all Cluster_hosts known to the system.
// Version: lima
func (clusterHost) GetAllRecords(session *Session) (retval map[ClusterHostRef]ClusterHostRecord, err error) {
	method := "Cluster_host.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterHostRefToClusterHostRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of Cluster_host references to Cluster_host records for all Cluster_hosts known to the system.
// Version: lima
func (clusterHost) GetAllRecords1(session *Session) (retval map[ClusterHostRef]ClusterHostRecord, err error) {
	method := "Cluster_host.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterHostRefToClusterHostRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the Cluster_hosts known to the system.
// Version: lima
func (clusterHost) GetAll(session *Session) (retval []ClusterHostRef, err error) {
	method := "Cluster_host.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterHostRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the Cluster_hosts known to the system.
// Version: lima
func (clusterHost) GetAll1(session *Session) (retval []ClusterHostRef, err error) {
	method := "Cluster_host.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterHostRefSet(method+" -> ", result)
	return
}

// Disable: Disable cluster membership for an enabled cluster host.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
func (clusterHost) Disable(session *Session, self ClusterHostRef) (err error) {
	method := "Cluster_host.disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisable: Disable cluster membership for an enabled cluster host.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
func (clusterHost) AsyncDisable(session *Session, self ClusterHostRef) (retval TaskRef, err error) {
	method := "Async.Cluster_host.disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Disable2: Disable cluster membership for an enabled cluster host.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
func (clusterHost) Disable2(session *Session, self ClusterHostRef) (err error) {
	method := "Cluster_host.disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDisable2: Disable cluster membership for an enabled cluster host.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
func (clusterHost) AsyncDisable2(session *Session, self ClusterHostRef) (retval TaskRef, err error) {
	method := "Async.Cluster_host.disable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// ForceDestroy: Remove a host from an existing cluster forcefully.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
func (clusterHost) ForceDestroy(session *Session, self ClusterHostRef) (err error) {
	method := "Cluster_host.force_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncForceDestroy: Remove a host from an existing cluster forcefully.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
func (clusterHost) AsyncForceDestroy(session *Session, self ClusterHostRef) (retval TaskRef, err error) {
	method := "Async.Cluster_host.force_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// ForceDestroy2: Remove a host from an existing cluster forcefully.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
func (clusterHost) ForceDestroy2(session *Session, self ClusterHostRef) (err error) {
	method := "Cluster_host.force_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncForceDestroy2: Remove a host from an existing cluster forcefully.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
func (clusterHost) AsyncForceDestroy2(session *Session, self ClusterHostRef) (retval TaskRef, err error) {
	method := "Async.Cluster_host.force_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Enable: Enable cluster membership for a disabled cluster host.
// Version: lima
//
// Errors:
// PIF_ALLOWS_UNPLUG - The operation you requested cannot be performed because the specified PIF allows unplug.
// REQUIRED_PIF_IS_UNPLUGGED - The operation you requested cannot be performed because the specified PIF is currently unplugged.
func (clusterHost) Enable(session *Session, self ClusterHostRef) (err error) {
	method := "Cluster_host.enable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncEnable: Enable cluster membership for a disabled cluster host.
// Version: lima
//
// Errors:
// PIF_ALLOWS_UNPLUG - The operation you requested cannot be performed because the specified PIF allows unplug.
// REQUIRED_PIF_IS_UNPLUGGED - The operation you requested cannot be performed because the specified PIF is currently unplugged.
func (clusterHost) AsyncEnable(session *Session, self ClusterHostRef) (retval TaskRef, err error) {
	method := "Async.Cluster_host.enable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Enable2: Enable cluster membership for a disabled cluster host.
// Version: lima
//
// Errors:
// PIF_ALLOWS_UNPLUG - The operation you requested cannot be performed because the specified PIF allows unplug.
// REQUIRED_PIF_IS_UNPLUGGED - The operation you requested cannot be performed because the specified PIF is currently unplugged.
func (clusterHost) Enable2(session *Session, self ClusterHostRef) (err error) {
	method := "Cluster_host.enable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncEnable2: Enable cluster membership for a disabled cluster host.
// Version: lima
//
// Errors:
// PIF_ALLOWS_UNPLUG - The operation you requested cannot be performed because the specified PIF allows unplug.
// REQUIRED_PIF_IS_UNPLUGGED - The operation you requested cannot be performed because the specified PIF is currently unplugged.
func (clusterHost) AsyncEnable2(session *Session, self ClusterHostRef) (retval TaskRef, err error) {
	method := "Async.Cluster_host.enable"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy: Remove the host from an existing cluster. This operation is allowed even if a cluster host is not enabled.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
// CLUSTERING_DISABLED - An operation was attempted while clustering was disabled on the cluster_host.
// CLUSTER_HOST_IS_LAST - The last cluster host cannot be destroyed. Destroy the cluster instead
func (clusterHost) Destroy(session *Session, self ClusterHostRef) (err error) {
	method := "Cluster_host.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Remove the host from an existing cluster. This operation is allowed even if a cluster host is not enabled.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
// CLUSTERING_DISABLED - An operation was attempted while clustering was disabled on the cluster_host.
// CLUSTER_HOST_IS_LAST - The last cluster host cannot be destroyed. Destroy the cluster instead
func (clusterHost) AsyncDestroy(session *Session, self ClusterHostRef) (retval TaskRef, err error) {
	method := "Async.Cluster_host.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Remove the host from an existing cluster. This operation is allowed even if a cluster host is not enabled.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
// CLUSTERING_DISABLED - An operation was attempted while clustering was disabled on the cluster_host.
// CLUSTER_HOST_IS_LAST - The last cluster host cannot be destroyed. Destroy the cluster instead
func (clusterHost) Destroy2(session *Session, self ClusterHostRef) (err error) {
	method := "Cluster_host.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Remove the host from an existing cluster. This operation is allowed even if a cluster host is not enabled.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
// CLUSTERING_DISABLED - An operation was attempted while clustering was disabled on the cluster_host.
// CLUSTER_HOST_IS_LAST - The last cluster host cannot be destroyed. Destroy the cluster instead
func (clusterHost) AsyncDestroy2(session *Session, self ClusterHostRef) (retval TaskRef, err error) {
	method := "Async.Cluster_host.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Add a new host to an existing cluster.
// Version: lima
//
// Errors:
// PIF_NOT_ATTACHED_TO_HOST - Cluster_host creation failed as the PIF provided is not attached to the host.
// NO_CLUSTER_HOSTS_REACHABLE - No other cluster host was reachable when joining
func (clusterHost) Create(session *Session, cluster ClusterRef, host HostRef, pif PIFRef) (retval ClusterHostRef, err error) {
	method := "Cluster_host.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	clusterArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "cluster"), cluster)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, clusterArg, hostArg, pifArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterHostRef(method+" -> ", result)
	return
}

// AsyncCreate: Add a new host to an existing cluster.
// Version: lima
//
// Errors:
// PIF_NOT_ATTACHED_TO_HOST - Cluster_host creation failed as the PIF provided is not attached to the host.
// NO_CLUSTER_HOSTS_REACHABLE - No other cluster host was reachable when joining
func (clusterHost) AsyncCreate(session *Session, cluster ClusterRef, host HostRef, pif PIFRef) (retval TaskRef, err error) {
	method := "Async.Cluster_host.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	clusterArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "cluster"), cluster)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, clusterArg, hostArg, pifArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create4: Add a new host to an existing cluster.
// Version: lima
//
// Errors:
// PIF_NOT_ATTACHED_TO_HOST - Cluster_host creation failed as the PIF provided is not attached to the host.
// NO_CLUSTER_HOSTS_REACHABLE - No other cluster host was reachable when joining
func (clusterHost) Create4(session *Session, cluster ClusterRef, host HostRef, pif PIFRef) (retval ClusterHostRef, err error) {
	method := "Cluster_host.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	clusterArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "cluster"), cluster)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, clusterArg, hostArg, pifArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterHostRef(method+" -> ", result)
	return
}

// AsyncCreate4: Add a new host to an existing cluster.
// Version: lima
//
// Errors:
// PIF_NOT_ATTACHED_TO_HOST - Cluster_host creation failed as the PIF provided is not attached to the host.
// NO_CLUSTER_HOSTS_REACHABLE - No other cluster host was reachable when joining
func (clusterHost) AsyncCreate4(session *Session, cluster ClusterRef, host HostRef, pif PIFRef) (retval TaskRef, err error) {
	method := "Async.Cluster_host.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	clusterArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "cluster"), cluster)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "pif"), pif)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, clusterArg, hostArg, pifArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetOtherConfig: Get the other_config field of the given Cluster_host.
// Version: lima
func (clusterHost) GetOtherConfig(session *Session, self ClusterHostRef) (retval map[string]string, err error) {
	method := "Cluster_host.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given Cluster_host.
// Version: lima
func (clusterHost) GetOtherConfig2(session *Session, self ClusterHostRef) (retval map[string]string, err error) {
	method := "Cluster_host.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCurrentOperations: Get the current_operations field of the given Cluster_host.
// Version: lima
func (clusterHost) GetCurrentOperations(session *Session, self ClusterHostRef) (retval map[string]ClusterHostOperation, err error) {
	method := "Cluster_host.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumClusterHostOperationMap(method+" -> ", result)
	return
}

// GetCurrentOperations2: Get the current_operations field of the given Cluster_host.
// Version: lima
func (clusterHost) GetCurrentOperations2(session *Session, self ClusterHostRef) (retval map[string]ClusterHostOperation, err error) {
	method := "Cluster_host.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumClusterHostOperationMap(method+" -> ", result)
	return
}

// GetAllowedOperations: Get the allowed_operations field of the given Cluster_host.
// Version: lima
func (clusterHost) GetAllowedOperations(session *Session, self ClusterHostRef) (retval []ClusterHostOperation, err error) {
	method := "Cluster_host.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumClusterHostOperationSet(method+" -> ", result)
	return
}

// GetAllowedOperations2: Get the allowed_operations field of the given Cluster_host.
// Version: lima
func (clusterHost) GetAllowedOperations2(session *Session, self ClusterHostRef) (retval []ClusterHostOperation, err error) {
	method := "Cluster_host.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumClusterHostOperationSet(method+" -> ", result)
	return
}

// GetLastUpdateLive: Get the last_update_live field of the given Cluster_host.
// Version: lima
func (clusterHost) GetLastUpdateLive(session *Session, self ClusterHostRef) (retval time.Time, err error) {
	method := "Cluster_host.get_last_update_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastUpdateLive2: Get the last_update_live field of the given Cluster_host.
// Version: lima
func (clusterHost) GetLastUpdateLive2(session *Session, self ClusterHostRef) (retval time.Time, err error) {
	method := "Cluster_host.get_last_update_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLive: Get the live field of the given Cluster_host.
// Version: lima
func (clusterHost) GetLive(session *Session, self ClusterHostRef) (retval bool, err error) {
	method := "Cluster_host.get_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLive2: Get the live field of the given Cluster_host.
// Version: lima
func (clusterHost) GetLive2(session *Session, self ClusterHostRef) (retval bool, err error) {
	method := "Cluster_host.get_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetJoined: Get the joined field of the given Cluster_host.
// Version: lima
func (clusterHost) GetJoined(session *Session, self ClusterHostRef) (retval bool, err error) {
	method := "Cluster_host.get_joined"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetJoined2: Get the joined field of the given Cluster_host.
// Version: lima
func (clusterHost) GetJoined2(session *Session, self ClusterHostRef) (retval bool, err error) {
	method := "Cluster_host.get_joined"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPIF: Get the PIF field of the given Cluster_host.
// Version: lima
func (clusterHost) GetPIF(session *Session, self ClusterHostRef) (retval PIFRef, err error) {
	method := "Cluster_host.get_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPIF2: Get the PIF field of the given Cluster_host.
// Version: lima
func (clusterHost) GetPIF2(session *Session, self ClusterHostRef) (retval PIFRef, err error) {
	method := "Cluster_host.get_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetEnabled: Get the enabled field of the given Cluster_host.
// Version: lima
func (clusterHost) GetEnabled(session *Session, self ClusterHostRef) (retval bool, err error) {
	method := "Cluster_host.get_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetEnabled2: Get the enabled field of the given Cluster_host.
// Version: lima
func (clusterHost) GetEnabled2(session *Session, self ClusterHostRef) (retval bool, err error) {
	method := "Cluster_host.get_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHost: Get the host field of the given Cluster_host.
// Version: lima
func (clusterHost) GetHost(session *Session, self ClusterHostRef) (retval HostRef, err error) {
	method := "Cluster_host.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHost2: Get the host field of the given Cluster_host.
// Version: lima
func (clusterHost) GetHost2(session *Session, self ClusterHostRef) (retval HostRef, err error) {
	method := "Cluster_host.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCluster: Get the cluster field of the given Cluster_host.
// Version: lima
func (clusterHost) GetCluster(session *Session, self ClusterHostRef) (retval ClusterRef, err error) {
	method := "Cluster_host.get_cluster"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterRef(method+" -> ", result)
	return
}

// GetCluster2: Get the cluster field of the given Cluster_host.
// Version: lima
func (clusterHost) GetCluster2(session *Session, self ClusterHostRef) (retval ClusterRef, err error) {
	method := "Cluster_host.get_cluster"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterRef(method+" -> ", result)
	return
}

// GetUUID: Get the uuid field of the given Cluster_host.
// Version: lima
func (clusterHost) GetUUID(session *Session, self ClusterHostRef) (retval string, err error) {
	method := "Cluster_host.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given Cluster_host.
// Version: lima
func (clusterHost) GetUUID2(session *Session, self ClusterHostRef) (retval string, err error) {
	method := "Cluster_host.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the Cluster_host instance with the specified UUID.
// Version: lima
func (clusterHost) GetByUUID(session *Session, uuid string) (retval ClusterHostRef, err error) {
	method := "Cluster_host.get_by_uuid"
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
	retval, err = deserializeClusterHostRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the Cluster_host instance with the specified UUID.
// Version: lima
func (clusterHost) GetByUUID2(session *Session, uuid string) (retval ClusterHostRef, err error) {
	method := "Cluster_host.get_by_uuid"
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
	retval, err = deserializeClusterHostRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given Cluster_host.
// Version: lima
func (clusterHost) GetRecord(session *Session, self ClusterHostRef) (retval ClusterHostRecord, err error) {
	method := "Cluster_host.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterHostRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given Cluster_host.
// Version: lima
func (clusterHost) GetRecord2(session *Session, self ClusterHostRef) (retval ClusterHostRecord, err error) {
	method := "Cluster_host.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterHostRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterHostRecord(method+" -> ", result)
	return
}
