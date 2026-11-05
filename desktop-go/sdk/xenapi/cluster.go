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

type ClusterRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// A list of the cluster_host objects associated with the Cluster
	ClusterHosts []ClusterHostRef `json:"cluster_hosts,omitempty"`
	// Internal field used by Host.destroy to store the IP of cluster members marked as permanently dead but not yet removed
	PendingForget []string `json:"pending_forget,omitempty"`
	// The secret key used by xapi-clusterd when it talks to itself on other hosts
	ClusterToken string `json:"cluster_token,omitempty"`
	// Simply the string &apos;corosync&apos;. No other cluster stacks are currently supported
	ClusterStack string `json:"cluster_stack,omitempty"`
	// Version of cluster stack, not writable via the API. Defaulting to 2 for backwards compatibility when upgrading from a cluster without this field, which means it is necessarily running version 2 of corosync, the only cluster stack supported so far.
	ClusterStackVersion int `json:"cluster_stack_version,omitempty"`
	// Whether the cluster stack thinks the cluster is quorate
	IsQuorate bool `json:"is_quorate,omitempty"`
	// Number of live hosts in order to be quorate
	Quorum int `json:"quorum,omitempty"`
	// Current number of live hosts, according to the cluster stack
	LiveHosts int `json:"live_hosts,omitempty"`
	// Total number of hosts expected by the cluster stack
	ExpectedHosts int `json:"expected_hosts,omitempty"`
	// list of the operations allowed in this state. This list is advisory only and the server state may have changed by the time this field is read by a client.
	AllowedOperations []ClusterOperation `json:"allowed_operations,omitempty"`
	// links each of the running tasks using this object (by reference) to a current_operation enum which describes the nature of the task.
	CurrentOperations map[string]ClusterOperation `json:"current_operations,omitempty"`
	// True if automatically joining new pool members to the cluster. This will be `true` in the first release
	PoolAutoJoin bool `json:"pool_auto_join,omitempty"`
	// The corosync token timeout in seconds
	TokenTimeout float64 `json:"token_timeout,omitempty"`
	// The corosync token timeout coefficient in seconds
	TokenTimeoutCoefficient float64 `json:"token_timeout_coefficient,omitempty"`
	// Contains read-only settings for the cluster, such as timeouts and other options. It can only be set at cluster create time
	ClusterConfig map[string]string `json:"cluster_config,omitempty"`
	// Additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
}

type ClusterRef string

// Cluster-wide Cluster metadata
type cluster struct{}

var Cluster cluster

// GetAllRecords: Return a map of Cluster references to Cluster records for all Clusters known to the system.
// Version: lima
func (cluster) GetAllRecords(session *Session) (retval map[ClusterRef]ClusterRecord, err error) {
	method := "Cluster.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterRefToClusterRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of Cluster references to Cluster records for all Clusters known to the system.
// Version: lima
func (cluster) GetAllRecords1(session *Session) (retval map[ClusterRef]ClusterRecord, err error) {
	method := "Cluster.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterRefToClusterRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the Clusters known to the system.
// Version: lima
func (cluster) GetAll(session *Session) (retval []ClusterRef, err error) {
	method := "Cluster.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the Clusters known to the system.
// Version: lima
func (cluster) GetAll1(session *Session) (retval []ClusterRef, err error) {
	method := "Cluster.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterRefSet(method+" -> ", result)
	return
}

// PoolResync: Resynchronise the cluster_host objects across the pool. Creates them where they need creating and then plugs them
// Version: lima
func (cluster) PoolResync(session *Session, self ClusterRef) (err error) {
	method := "Cluster.pool_resync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolResync: Resynchronise the cluster_host objects across the pool. Creates them where they need creating and then plugs them
// Version: lima
func (cluster) AsyncPoolResync(session *Session, self ClusterRef) (retval TaskRef, err error) {
	method := "Async.Cluster.pool_resync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolResync2: Resynchronise the cluster_host objects across the pool. Creates them where they need creating and then plugs them
// Version: lima
func (cluster) PoolResync2(session *Session, self ClusterRef) (err error) {
	method := "Cluster.pool_resync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolResync2: Resynchronise the cluster_host objects across the pool. Creates them where they need creating and then plugs them
// Version: lima
func (cluster) AsyncPoolResync2(session *Session, self ClusterRef) (retval TaskRef, err error) {
	method := "Async.Cluster.pool_resync"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolDestroy: Attempt to destroy the Cluster_host objects for all hosts in the pool and then destroy the Cluster.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
// CLUSTERING_DISABLED - An operation was attempted while clustering was disabled on the cluster_host.
// CLUSTER_HOST_IS_LAST - The last cluster host cannot be destroyed. Destroy the cluster instead
func (cluster) PoolDestroy(session *Session, self ClusterRef) (err error) {
	method := "Cluster.pool_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolDestroy: Attempt to destroy the Cluster_host objects for all hosts in the pool and then destroy the Cluster.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
// CLUSTERING_DISABLED - An operation was attempted while clustering was disabled on the cluster_host.
// CLUSTER_HOST_IS_LAST - The last cluster host cannot be destroyed. Destroy the cluster instead
func (cluster) AsyncPoolDestroy(session *Session, self ClusterRef) (retval TaskRef, err error) {
	method := "Async.Cluster.pool_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolDestroy2: Attempt to destroy the Cluster_host objects for all hosts in the pool and then destroy the Cluster.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
// CLUSTERING_DISABLED - An operation was attempted while clustering was disabled on the cluster_host.
// CLUSTER_HOST_IS_LAST - The last cluster host cannot be destroyed. Destroy the cluster instead
func (cluster) PoolDestroy2(session *Session, self ClusterRef) (err error) {
	method := "Cluster.pool_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolDestroy2: Attempt to destroy the Cluster_host objects for all hosts in the pool and then destroy the Cluster.
// Version: lima
//
// Errors:
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
// CLUSTERING_DISABLED - An operation was attempted while clustering was disabled on the cluster_host.
// CLUSTER_HOST_IS_LAST - The last cluster host cannot be destroyed. Destroy the cluster instead
func (cluster) AsyncPoolDestroy2(session *Session, self ClusterRef) (retval TaskRef, err error) {
	method := "Async.Cluster.pool_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolForceDestroy: Attempt to force destroy the Cluster_host objects, and then destroy the Cluster.
// Version: lima
//
// Errors:
// CLUSTER_FORCE_DESTROY_FAILED - Force destroy failed on a Cluster_host while force destroying the cluster.
func (cluster) PoolForceDestroy(session *Session, self ClusterRef) (err error) {
	method := "Cluster.pool_force_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolForceDestroy: Attempt to force destroy the Cluster_host objects, and then destroy the Cluster.
// Version: lima
//
// Errors:
// CLUSTER_FORCE_DESTROY_FAILED - Force destroy failed on a Cluster_host while force destroying the cluster.
func (cluster) AsyncPoolForceDestroy(session *Session, self ClusterRef) (retval TaskRef, err error) {
	method := "Async.Cluster.pool_force_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolForceDestroy2: Attempt to force destroy the Cluster_host objects, and then destroy the Cluster.
// Version: lima
//
// Errors:
// CLUSTER_FORCE_DESTROY_FAILED - Force destroy failed on a Cluster_host while force destroying the cluster.
func (cluster) PoolForceDestroy2(session *Session, self ClusterRef) (err error) {
	method := "Cluster.pool_force_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPoolForceDestroy2: Attempt to force destroy the Cluster_host objects, and then destroy the Cluster.
// Version: lima
//
// Errors:
// CLUSTER_FORCE_DESTROY_FAILED - Force destroy failed on a Cluster_host while force destroying the cluster.
func (cluster) AsyncPoolForceDestroy2(session *Session, self ClusterRef) (retval TaskRef, err error) {
	method := "Async.Cluster.pool_force_destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// PoolCreate: Attempt to create a Cluster from the entire pool
// Version: kolkata
func (cluster) PoolCreate(session *Session, network NetworkRef, clusterStack string, tokenTimeout float64, tokenTimeoutCoefficient float64) (retval ClusterRef, err error) {
	method := "Cluster.pool_create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	clusterStackArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cluster_stack"), clusterStack)
	if err != nil {
		return
	}
	tokenTimeoutArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout"), tokenTimeout)
	if err != nil {
		return
	}
	tokenTimeoutCoefficientArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout_coefficient"), tokenTimeoutCoefficient)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg, clusterStackArg, tokenTimeoutArg, tokenTimeoutCoefficientArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterRef(method+" -> ", result)
	return
}

// AsyncPoolCreate: Attempt to create a Cluster from the entire pool
// Version: kolkata
func (cluster) AsyncPoolCreate(session *Session, network NetworkRef, clusterStack string, tokenTimeout float64, tokenTimeoutCoefficient float64) (retval TaskRef, err error) {
	method := "Async.Cluster.pool_create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	clusterStackArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cluster_stack"), clusterStack)
	if err != nil {
		return
	}
	tokenTimeoutArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout"), tokenTimeout)
	if err != nil {
		return
	}
	tokenTimeoutCoefficientArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout_coefficient"), tokenTimeoutCoefficient)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg, clusterStackArg, tokenTimeoutArg, tokenTimeoutCoefficientArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// PoolCreate5: Attempt to create a Cluster from the entire pool
// Version: kolkata
func (cluster) PoolCreate5(session *Session, network NetworkRef, clusterStack string, tokenTimeout float64, tokenTimeoutCoefficient float64) (retval ClusterRef, err error) {
	method := "Cluster.pool_create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	clusterStackArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cluster_stack"), clusterStack)
	if err != nil {
		return
	}
	tokenTimeoutArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout"), tokenTimeout)
	if err != nil {
		return
	}
	tokenTimeoutCoefficientArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout_coefficient"), tokenTimeoutCoefficient)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg, clusterStackArg, tokenTimeoutArg, tokenTimeoutCoefficientArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterRef(method+" -> ", result)
	return
}

// AsyncPoolCreate5: Attempt to create a Cluster from the entire pool
// Version: kolkata
func (cluster) AsyncPoolCreate5(session *Session, network NetworkRef, clusterStack string, tokenTimeout float64, tokenTimeoutCoefficient float64) (retval TaskRef, err error) {
	method := "Async.Cluster.pool_create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	clusterStackArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cluster_stack"), clusterStack)
	if err != nil {
		return
	}
	tokenTimeoutArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout"), tokenTimeout)
	if err != nil {
		return
	}
	tokenTimeoutCoefficientArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout_coefficient"), tokenTimeoutCoefficient)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, networkArg, clusterStackArg, tokenTimeoutArg, tokenTimeoutCoefficientArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetNetwork: Returns the network used by the cluster for inter-host communication, i.e. the network shared by all cluster host PIFs
// Version: lima
func (cluster) GetNetwork(session *Session, self ClusterRef) (retval NetworkRef, err error) {
	method := "Cluster.get_network"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkRef(method+" -> ", result)
	return
}

// AsyncGetNetwork: Returns the network used by the cluster for inter-host communication, i.e. the network shared by all cluster host PIFs
// Version: lima
func (cluster) AsyncGetNetwork(session *Session, self ClusterRef) (retval TaskRef, err error) {
	method := "Async.Cluster.get_network"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNetwork2: Returns the network used by the cluster for inter-host communication, i.e. the network shared by all cluster host PIFs
// Version: lima
func (cluster) GetNetwork2(session *Session, self ClusterRef) (retval NetworkRef, err error) {
	method := "Cluster.get_network"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkRef(method+" -> ", result)
	return
}

// AsyncGetNetwork2: Returns the network used by the cluster for inter-host communication, i.e. the network shared by all cluster host PIFs
// Version: lima
func (cluster) AsyncGetNetwork2(session *Session, self ClusterRef) (retval TaskRef, err error) {
	method := "Async.Cluster.get_network"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy: Destroys a Cluster object and the one remaining Cluster_host member
// Version: lima
//
// Errors:
// CLUSTER_DOES_NOT_HAVE_ONE_NODE - An operation failed as it expected the cluster to have only one node but found multiple cluster_hosts.
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
func (cluster) Destroy(session *Session, self ClusterRef) (err error) {
	method := "Cluster.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroys a Cluster object and the one remaining Cluster_host member
// Version: lima
//
// Errors:
// CLUSTER_DOES_NOT_HAVE_ONE_NODE - An operation failed as it expected the cluster to have only one node but found multiple cluster_hosts.
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
func (cluster) AsyncDestroy(session *Session, self ClusterRef) (retval TaskRef, err error) {
	method := "Async.Cluster.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroys a Cluster object and the one remaining Cluster_host member
// Version: lima
//
// Errors:
// CLUSTER_DOES_NOT_HAVE_ONE_NODE - An operation failed as it expected the cluster to have only one node but found multiple cluster_hosts.
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
func (cluster) Destroy2(session *Session, self ClusterRef) (err error) {
	method := "Cluster.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroys a Cluster object and the one remaining Cluster_host member
// Version: lima
//
// Errors:
// CLUSTER_DOES_NOT_HAVE_ONE_NODE - An operation failed as it expected the cluster to have only one node but found multiple cluster_hosts.
// CLUSTER_STACK_IN_USE - The cluster stack is still in use by at least one plugged PBD.
func (cluster) AsyncDestroy2(session *Session, self ClusterRef) (retval TaskRef, err error) {
	method := "Async.Cluster.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Creates a Cluster object and one Cluster_host object as its first member
// Version: kolkata
//
// Errors:
// INVALID_CLUSTER_STACK - The cluster stack provided is not supported.
// INVALID_VALUE - The value given is invalid
// PIF_ALLOWS_UNPLUG - The operation you requested cannot be performed because the specified PIF allows unplug.
// REQUIRED_PIF_IS_UNPLUGGED - The operation you requested cannot be performed because the specified PIF is currently unplugged.
func (cluster) Create(session *Session, pif PIFRef, clusterStack string, poolAutoJoin bool, tokenTimeout float64, tokenTimeoutCoefficient float64) (retval ClusterRef, err error) {
	method := "Cluster.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "PIF"), pif)
	if err != nil {
		return
	}
	clusterStackArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cluster_stack"), clusterStack)
	if err != nil {
		return
	}
	poolAutoJoinArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "pool_auto_join"), poolAutoJoin)
	if err != nil {
		return
	}
	tokenTimeoutArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout"), tokenTimeout)
	if err != nil {
		return
	}
	tokenTimeoutCoefficientArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout_coefficient"), tokenTimeoutCoefficient)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg, clusterStackArg, poolAutoJoinArg, tokenTimeoutArg, tokenTimeoutCoefficientArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterRef(method+" -> ", result)
	return
}

// AsyncCreate: Creates a Cluster object and one Cluster_host object as its first member
// Version: kolkata
//
// Errors:
// INVALID_CLUSTER_STACK - The cluster stack provided is not supported.
// INVALID_VALUE - The value given is invalid
// PIF_ALLOWS_UNPLUG - The operation you requested cannot be performed because the specified PIF allows unplug.
// REQUIRED_PIF_IS_UNPLUGGED - The operation you requested cannot be performed because the specified PIF is currently unplugged.
func (cluster) AsyncCreate(session *Session, pif PIFRef, clusterStack string, poolAutoJoin bool, tokenTimeout float64, tokenTimeoutCoefficient float64) (retval TaskRef, err error) {
	method := "Async.Cluster.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "PIF"), pif)
	if err != nil {
		return
	}
	clusterStackArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cluster_stack"), clusterStack)
	if err != nil {
		return
	}
	poolAutoJoinArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "pool_auto_join"), poolAutoJoin)
	if err != nil {
		return
	}
	tokenTimeoutArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout"), tokenTimeout)
	if err != nil {
		return
	}
	tokenTimeoutCoefficientArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout_coefficient"), tokenTimeoutCoefficient)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg, clusterStackArg, poolAutoJoinArg, tokenTimeoutArg, tokenTimeoutCoefficientArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create6: Creates a Cluster object and one Cluster_host object as its first member
// Version: kolkata
//
// Errors:
// INVALID_CLUSTER_STACK - The cluster stack provided is not supported.
// INVALID_VALUE - The value given is invalid
// PIF_ALLOWS_UNPLUG - The operation you requested cannot be performed because the specified PIF allows unplug.
// REQUIRED_PIF_IS_UNPLUGGED - The operation you requested cannot be performed because the specified PIF is currently unplugged.
func (cluster) Create6(session *Session, pif PIFRef, clusterStack string, poolAutoJoin bool, tokenTimeout float64, tokenTimeoutCoefficient float64) (retval ClusterRef, err error) {
	method := "Cluster.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "PIF"), pif)
	if err != nil {
		return
	}
	clusterStackArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cluster_stack"), clusterStack)
	if err != nil {
		return
	}
	poolAutoJoinArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "pool_auto_join"), poolAutoJoin)
	if err != nil {
		return
	}
	tokenTimeoutArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout"), tokenTimeout)
	if err != nil {
		return
	}
	tokenTimeoutCoefficientArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout_coefficient"), tokenTimeoutCoefficient)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg, clusterStackArg, poolAutoJoinArg, tokenTimeoutArg, tokenTimeoutCoefficientArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterRef(method+" -> ", result)
	return
}

// AsyncCreate6: Creates a Cluster object and one Cluster_host object as its first member
// Version: kolkata
//
// Errors:
// INVALID_CLUSTER_STACK - The cluster stack provided is not supported.
// INVALID_VALUE - The value given is invalid
// PIF_ALLOWS_UNPLUG - The operation you requested cannot be performed because the specified PIF allows unplug.
// REQUIRED_PIF_IS_UNPLUGGED - The operation you requested cannot be performed because the specified PIF is currently unplugged.
func (cluster) AsyncCreate6(session *Session, pif PIFRef, clusterStack string, poolAutoJoin bool, tokenTimeout float64, tokenTimeoutCoefficient float64) (retval TaskRef, err error) {
	method := "Async.Cluster.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	pifArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "PIF"), pif)
	if err != nil {
		return
	}
	clusterStackArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "cluster_stack"), clusterStack)
	if err != nil {
		return
	}
	poolAutoJoinArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "pool_auto_join"), poolAutoJoin)
	if err != nil {
		return
	}
	tokenTimeoutArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout"), tokenTimeout)
	if err != nil {
		return
	}
	tokenTimeoutCoefficientArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "token_timeout_coefficient"), tokenTimeoutCoefficient)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, pifArg, clusterStackArg, poolAutoJoinArg, tokenTimeoutArg, tokenTimeoutCoefficientArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given Cluster.  If the key is not in that Map, then do nothing.
// Version: lima
func (cluster) RemoveFromOtherConfig(session *Session, self ClusterRef, key string) (err error) {
	method := "Cluster.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given Cluster.  If the key is not in that Map, then do nothing.
// Version: lima
func (cluster) RemoveFromOtherConfig3(session *Session, self ClusterRef, key string) (err error) {
	method := "Cluster.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given Cluster.
// Version: lima
func (cluster) AddToOtherConfig(session *Session, self ClusterRef, key string, value string) (err error) {
	method := "Cluster.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given Cluster.
// Version: lima
func (cluster) AddToOtherConfig4(session *Session, self ClusterRef, key string, value string) (err error) {
	method := "Cluster.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig: Set the other_config field of the given Cluster.
// Version: lima
func (cluster) SetOtherConfig(session *Session, self ClusterRef, value map[string]string) (err error) {
	method := "Cluster.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given Cluster.
// Version: lima
func (cluster) SetOtherConfig3(session *Session, self ClusterRef, value map[string]string) (err error) {
	method := "Cluster.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given Cluster.
// Version: lima
func (cluster) GetOtherConfig(session *Session, self ClusterRef) (retval map[string]string, err error) {
	method := "Cluster.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given Cluster.
// Version: lima
func (cluster) GetOtherConfig2(session *Session, self ClusterRef) (retval map[string]string, err error) {
	method := "Cluster.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClusterConfig: Get the cluster_config field of the given Cluster.
// Version: lima
func (cluster) GetClusterConfig(session *Session, self ClusterRef) (retval map[string]string, err error) {
	method := "Cluster.get_cluster_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClusterConfig2: Get the cluster_config field of the given Cluster.
// Version: lima
func (cluster) GetClusterConfig2(session *Session, self ClusterRef) (retval map[string]string, err error) {
	method := "Cluster.get_cluster_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTokenTimeoutCoefficient: Get the token_timeout_coefficient field of the given Cluster.
// Version: lima
func (cluster) GetTokenTimeoutCoefficient(session *Session, self ClusterRef) (retval float64, err error) {
	method := "Cluster.get_token_timeout_coefficient"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTokenTimeoutCoefficient2: Get the token_timeout_coefficient field of the given Cluster.
// Version: lima
func (cluster) GetTokenTimeoutCoefficient2(session *Session, self ClusterRef) (retval float64, err error) {
	method := "Cluster.get_token_timeout_coefficient"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTokenTimeout: Get the token_timeout field of the given Cluster.
// Version: lima
func (cluster) GetTokenTimeout(session *Session, self ClusterRef) (retval float64, err error) {
	method := "Cluster.get_token_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTokenTimeout2: Get the token_timeout field of the given Cluster.
// Version: lima
func (cluster) GetTokenTimeout2(session *Session, self ClusterRef) (retval float64, err error) {
	method := "Cluster.get_token_timeout"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPoolAutoJoin: Get the pool_auto_join field of the given Cluster.
// Version: lima
func (cluster) GetPoolAutoJoin(session *Session, self ClusterRef) (retval bool, err error) {
	method := "Cluster.get_pool_auto_join"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPoolAutoJoin2: Get the pool_auto_join field of the given Cluster.
// Version: lima
func (cluster) GetPoolAutoJoin2(session *Session, self ClusterRef) (retval bool, err error) {
	method := "Cluster.get_pool_auto_join"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCurrentOperations: Get the current_operations field of the given Cluster.
// Version: lima
func (cluster) GetCurrentOperations(session *Session, self ClusterRef) (retval map[string]ClusterOperation, err error) {
	method := "Cluster.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumClusterOperationMap(method+" -> ", result)
	return
}

// GetCurrentOperations2: Get the current_operations field of the given Cluster.
// Version: lima
func (cluster) GetCurrentOperations2(session *Session, self ClusterRef) (retval map[string]ClusterOperation, err error) {
	method := "Cluster.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumClusterOperationMap(method+" -> ", result)
	return
}

// GetAllowedOperations: Get the allowed_operations field of the given Cluster.
// Version: lima
func (cluster) GetAllowedOperations(session *Session, self ClusterRef) (retval []ClusterOperation, err error) {
	method := "Cluster.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumClusterOperationSet(method+" -> ", result)
	return
}

// GetAllowedOperations2: Get the allowed_operations field of the given Cluster.
// Version: lima
func (cluster) GetAllowedOperations2(session *Session, self ClusterRef) (retval []ClusterOperation, err error) {
	method := "Cluster.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumClusterOperationSet(method+" -> ", result)
	return
}

// GetExpectedHosts: Get the expected_hosts field of the given Cluster.
// Version: lima
func (cluster) GetExpectedHosts(session *Session, self ClusterRef) (retval int, err error) {
	method := "Cluster.get_expected_hosts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetExpectedHosts2: Get the expected_hosts field of the given Cluster.
// Version: lima
func (cluster) GetExpectedHosts2(session *Session, self ClusterRef) (retval int, err error) {
	method := "Cluster.get_expected_hosts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLiveHosts: Get the live_hosts field of the given Cluster.
// Version: lima
func (cluster) GetLiveHosts(session *Session, self ClusterRef) (retval int, err error) {
	method := "Cluster.get_live_hosts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLiveHosts2: Get the live_hosts field of the given Cluster.
// Version: lima
func (cluster) GetLiveHosts2(session *Session, self ClusterRef) (retval int, err error) {
	method := "Cluster.get_live_hosts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetQuorum: Get the quorum field of the given Cluster.
// Version: lima
func (cluster) GetQuorum(session *Session, self ClusterRef) (retval int, err error) {
	method := "Cluster.get_quorum"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetQuorum2: Get the quorum field of the given Cluster.
// Version: lima
func (cluster) GetQuorum2(session *Session, self ClusterRef) (retval int, err error) {
	method := "Cluster.get_quorum"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsQuorate: Get the is_quorate field of the given Cluster.
// Version: lima
func (cluster) GetIsQuorate(session *Session, self ClusterRef) (retval bool, err error) {
	method := "Cluster.get_is_quorate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsQuorate2: Get the is_quorate field of the given Cluster.
// Version: lima
func (cluster) GetIsQuorate2(session *Session, self ClusterRef) (retval bool, err error) {
	method := "Cluster.get_is_quorate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClusterStackVersion: Get the cluster_stack_version field of the given Cluster.
// Version: lima
func (cluster) GetClusterStackVersion(session *Session, self ClusterRef) (retval int, err error) {
	method := "Cluster.get_cluster_stack_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClusterStackVersion2: Get the cluster_stack_version field of the given Cluster.
// Version: lima
func (cluster) GetClusterStackVersion2(session *Session, self ClusterRef) (retval int, err error) {
	method := "Cluster.get_cluster_stack_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClusterStack: Get the cluster_stack field of the given Cluster.
// Version: lima
func (cluster) GetClusterStack(session *Session, self ClusterRef) (retval string, err error) {
	method := "Cluster.get_cluster_stack"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClusterStack2: Get the cluster_stack field of the given Cluster.
// Version: lima
func (cluster) GetClusterStack2(session *Session, self ClusterRef) (retval string, err error) {
	method := "Cluster.get_cluster_stack"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClusterToken: Get the cluster_token field of the given Cluster.
// Version: lima
func (cluster) GetClusterToken(session *Session, self ClusterRef) (retval string, err error) {
	method := "Cluster.get_cluster_token"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClusterToken2: Get the cluster_token field of the given Cluster.
// Version: lima
func (cluster) GetClusterToken2(session *Session, self ClusterRef) (retval string, err error) {
	method := "Cluster.get_cluster_token"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPendingForget: Get the pending_forget field of the given Cluster.
// Version: lima
func (cluster) GetPendingForget(session *Session, self ClusterRef) (retval []string, err error) {
	method := "Cluster.get_pending_forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPendingForget2: Get the pending_forget field of the given Cluster.
// Version: lima
func (cluster) GetPendingForget2(session *Session, self ClusterRef) (retval []string, err error) {
	method := "Cluster.get_pending_forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetClusterHosts: Get the cluster_hosts field of the given Cluster.
// Version: lima
func (cluster) GetClusterHosts(session *Session, self ClusterRef) (retval []ClusterHostRef, err error) {
	method := "Cluster.get_cluster_hosts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterHostRefSet(method+" -> ", result)
	return
}

// GetClusterHosts2: Get the cluster_hosts field of the given Cluster.
// Version: lima
func (cluster) GetClusterHosts2(session *Session, self ClusterRef) (retval []ClusterHostRef, err error) {
	method := "Cluster.get_cluster_hosts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterHostRefSet(method+" -> ", result)
	return
}

// GetUUID: Get the uuid field of the given Cluster.
// Version: lima
func (cluster) GetUUID(session *Session, self ClusterRef) (retval string, err error) {
	method := "Cluster.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given Cluster.
// Version: lima
func (cluster) GetUUID2(session *Session, self ClusterRef) (retval string, err error) {
	method := "Cluster.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the Cluster instance with the specified UUID.
// Version: lima
func (cluster) GetByUUID(session *Session, uuid string) (retval ClusterRef, err error) {
	method := "Cluster.get_by_uuid"
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
	retval, err = deserializeClusterRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the Cluster instance with the specified UUID.
// Version: lima
func (cluster) GetByUUID2(session *Session, uuid string) (retval ClusterRef, err error) {
	method := "Cluster.get_by_uuid"
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
	retval, err = deserializeClusterRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given Cluster.
// Version: lima
func (cluster) GetRecord(session *Session, self ClusterRef) (retval ClusterRecord, err error) {
	method := "Cluster.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given Cluster.
// Version: lima
func (cluster) GetRecord2(session *Session, self ClusterRef) (retval ClusterRecord, err error) {
	method := "Cluster.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeClusterRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeClusterRecord(method+" -> ", result)
	return
}
