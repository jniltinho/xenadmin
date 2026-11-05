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

type TunnelRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// The interface through which the tunnel is accessed
	AccessPIF PIFRef `json:"access_PIF,omitempty"`
	// The interface used by the tunnel
	TransportPIF PIFRef `json:"transport_PIF,omitempty"`
	// Status information about the tunnel
	Status map[string]string `json:"status,omitempty"`
	// Additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// The protocol used for tunneling (either GRE or VxLAN)
	Protocol TunnelProtocol `json:"protocol,omitempty"`
}

type TunnelRef string

// A tunnel for network traffic
type tunnel struct{}

var Tunnel tunnel

// GetAllRecords: Return a map of tunnel references to tunnel records for all tunnels known to the system.
// Version: cowley
func (tunnel) GetAllRecords(session *Session) (retval map[TunnelRef]TunnelRecord, err error) {
	method := "tunnel.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRefToTunnelRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of tunnel references to tunnel records for all tunnels known to the system.
// Version: cowley
func (tunnel) GetAllRecords1(session *Session) (retval map[TunnelRef]TunnelRecord, err error) {
	method := "tunnel.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRefToTunnelRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the tunnels known to the system.
// Version: cowley
func (tunnel) GetAll(session *Session) (retval []TunnelRef, err error) {
	method := "tunnel.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the tunnels known to the system.
// Version: cowley
func (tunnel) GetAll1(session *Session) (retval []TunnelRef, err error) {
	method := "tunnel.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRefSet(method+" -> ", result)
	return
}

// Destroy: Destroy a tunnel
// Version: cowley
func (tunnel) Destroy(session *Session, self TunnelRef) (err error) {
	method := "tunnel.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy a tunnel
// Version: cowley
func (tunnel) AsyncDestroy(session *Session, self TunnelRef) (retval TaskRef, err error) {
	method := "Async.tunnel.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy a tunnel
// Version: cowley
func (tunnel) Destroy2(session *Session, self TunnelRef) (err error) {
	method := "tunnel.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy a tunnel
// Version: cowley
func (tunnel) AsyncDestroy2(session *Session, self TunnelRef) (retval TaskRef, err error) {
	method := "Async.tunnel.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Create a tunnel
// Version: 1.250.0
//
// Errors:
// OPENVSWITCH_NOT_ACTIVE - This operation needs the OpenVSwitch networking backend to be enabled on all hosts in the pool.
// TRANSPORT_PIF_NOT_CONFIGURED - The tunnel transport PIF has no IP configuration set.
// IS_TUNNEL_ACCESS_PIF - Cannot create a VLAN or tunnel on top of a tunnel access PIF - use the underlying transport PIF instead.
func (tunnel) Create(session *Session, transportPIF PIFRef, network NetworkRef, protocol TunnelProtocol) (retval TunnelRef, err error) {
	method := "tunnel.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	transportPIFArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "transport_PIF"), transportPIF)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	protocolArg, err := serializeEnumTunnelProtocol(fmt.Sprintf("%s(%s)", method, "protocol"), protocol)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, transportPIFArg, networkArg, protocolArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRef(method+" -> ", result)
	return
}

// AsyncCreate: Create a tunnel
// Version: 1.250.0
//
// Errors:
// OPENVSWITCH_NOT_ACTIVE - This operation needs the OpenVSwitch networking backend to be enabled on all hosts in the pool.
// TRANSPORT_PIF_NOT_CONFIGURED - The tunnel transport PIF has no IP configuration set.
// IS_TUNNEL_ACCESS_PIF - Cannot create a VLAN or tunnel on top of a tunnel access PIF - use the underlying transport PIF instead.
func (tunnel) AsyncCreate(session *Session, transportPIF PIFRef, network NetworkRef, protocol TunnelProtocol) (retval TaskRef, err error) {
	method := "Async.tunnel.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	transportPIFArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "transport_PIF"), transportPIF)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	protocolArg, err := serializeEnumTunnelProtocol(fmt.Sprintf("%s(%s)", method, "protocol"), protocol)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, transportPIFArg, networkArg, protocolArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create4: Create a tunnel
// Version: 1.250.0
//
// Errors:
// OPENVSWITCH_NOT_ACTIVE - This operation needs the OpenVSwitch networking backend to be enabled on all hosts in the pool.
// TRANSPORT_PIF_NOT_CONFIGURED - The tunnel transport PIF has no IP configuration set.
// IS_TUNNEL_ACCESS_PIF - Cannot create a VLAN or tunnel on top of a tunnel access PIF - use the underlying transport PIF instead.
func (tunnel) Create4(session *Session, transportPIF PIFRef, network NetworkRef, protocol TunnelProtocol) (retval TunnelRef, err error) {
	method := "tunnel.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	transportPIFArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "transport_PIF"), transportPIF)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	protocolArg, err := serializeEnumTunnelProtocol(fmt.Sprintf("%s(%s)", method, "protocol"), protocol)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, transportPIFArg, networkArg, protocolArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRef(method+" -> ", result)
	return
}

// AsyncCreate4: Create a tunnel
// Version: 1.250.0
//
// Errors:
// OPENVSWITCH_NOT_ACTIVE - This operation needs the OpenVSwitch networking backend to be enabled on all hosts in the pool.
// TRANSPORT_PIF_NOT_CONFIGURED - The tunnel transport PIF has no IP configuration set.
// IS_TUNNEL_ACCESS_PIF - Cannot create a VLAN or tunnel on top of a tunnel access PIF - use the underlying transport PIF instead.
func (tunnel) AsyncCreate4(session *Session, transportPIF PIFRef, network NetworkRef, protocol TunnelProtocol) (retval TaskRef, err error) {
	method := "Async.tunnel.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	transportPIFArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "transport_PIF"), transportPIF)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	protocolArg, err := serializeEnumTunnelProtocol(fmt.Sprintf("%s(%s)", method, "protocol"), protocol)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, transportPIFArg, networkArg, protocolArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Create3: Create a tunnel
// Version: dundee
//
// Errors:
// OPENVSWITCH_NOT_ACTIVE - This operation needs the OpenVSwitch networking backend to be enabled on all hosts in the pool.
// TRANSPORT_PIF_NOT_CONFIGURED - The tunnel transport PIF has no IP configuration set.
// IS_TUNNEL_ACCESS_PIF - Cannot create a VLAN or tunnel on top of a tunnel access PIF - use the underlying transport PIF instead.
func (tunnel) Create3(session *Session, transportPIF PIFRef, network NetworkRef) (retval TunnelRef, err error) {
	method := "tunnel.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	transportPIFArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "transport_PIF"), transportPIF)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, transportPIFArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRef(method+" -> ", result)
	return
}

// AsyncCreate3: Create a tunnel
// Version: dundee
//
// Errors:
// OPENVSWITCH_NOT_ACTIVE - This operation needs the OpenVSwitch networking backend to be enabled on all hosts in the pool.
// TRANSPORT_PIF_NOT_CONFIGURED - The tunnel transport PIF has no IP configuration set.
// IS_TUNNEL_ACCESS_PIF - Cannot create a VLAN or tunnel on top of a tunnel access PIF - use the underlying transport PIF instead.
func (tunnel) AsyncCreate3(session *Session, transportPIF PIFRef, network NetworkRef) (retval TaskRef, err error) {
	method := "Async.tunnel.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	transportPIFArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "transport_PIF"), transportPIF)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, transportPIFArg, networkArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetProtocol: Set the protocol field of the given tunnel.
// Version: 1.250.0
func (tunnel) SetProtocol(session *Session, self TunnelRef, value TunnelProtocol) (err error) {
	method := "tunnel.set_protocol"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumTunnelProtocol(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetProtocol3: Set the protocol field of the given tunnel.
// Version: 1.250.0
func (tunnel) SetProtocol3(session *Session, self TunnelRef, value TunnelProtocol) (err error) {
	method := "tunnel.set_protocol"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumTunnelProtocol(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetProtocol2: Set the protocol field of the given tunnel.
// Version: cowley
func (tunnel) SetProtocol2(session *Session, self TunnelRef) (err error) {
	method := "tunnel.set_protocol"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given tunnel.  If the key is not in that Map, then do nothing.
// Version: cowley
func (tunnel) RemoveFromOtherConfig(session *Session, self TunnelRef, key string) (err error) {
	method := "tunnel.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given tunnel.  If the key is not in that Map, then do nothing.
// Version: cowley
func (tunnel) RemoveFromOtherConfig3(session *Session, self TunnelRef, key string) (err error) {
	method := "tunnel.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given tunnel.
// Version: cowley
func (tunnel) AddToOtherConfig(session *Session, self TunnelRef, key string, value string) (err error) {
	method := "tunnel.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given tunnel.
// Version: cowley
func (tunnel) AddToOtherConfig4(session *Session, self TunnelRef, key string, value string) (err error) {
	method := "tunnel.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig: Set the other_config field of the given tunnel.
// Version: cowley
func (tunnel) SetOtherConfig(session *Session, self TunnelRef, value map[string]string) (err error) {
	method := "tunnel.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given tunnel.
// Version: cowley
func (tunnel) SetOtherConfig3(session *Session, self TunnelRef, value map[string]string) (err error) {
	method := "tunnel.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromStatus: Remove the given key and its corresponding value from the status field of the given tunnel.  If the key is not in that Map, then do nothing.
// Version: cowley
func (tunnel) RemoveFromStatus(session *Session, self TunnelRef, key string) (err error) {
	method := "tunnel.remove_from_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromStatus3: Remove the given key and its corresponding value from the status field of the given tunnel.  If the key is not in that Map, then do nothing.
// Version: cowley
func (tunnel) RemoveFromStatus3(session *Session, self TunnelRef, key string) (err error) {
	method := "tunnel.remove_from_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToStatus: Add the given key-value pair to the status field of the given tunnel.
// Version: cowley
func (tunnel) AddToStatus(session *Session, self TunnelRef, key string, value string) (err error) {
	method := "tunnel.add_to_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToStatus4: Add the given key-value pair to the status field of the given tunnel.
// Version: cowley
func (tunnel) AddToStatus4(session *Session, self TunnelRef, key string, value string) (err error) {
	method := "tunnel.add_to_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetStatus: Set the status field of the given tunnel.
// Version: cowley
func (tunnel) SetStatus(session *Session, self TunnelRef, value map[string]string) (err error) {
	method := "tunnel.set_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetStatus3: Set the status field of the given tunnel.
// Version: cowley
func (tunnel) SetStatus3(session *Session, self TunnelRef, value map[string]string) (err error) {
	method := "tunnel.set_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetProtocol: Get the protocol field of the given tunnel.
// Version: cowley
func (tunnel) GetProtocol(session *Session, self TunnelRef) (retval TunnelProtocol, err error) {
	method := "tunnel.get_protocol"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumTunnelProtocol(method+" -> ", result)
	return
}

// GetProtocol2: Get the protocol field of the given tunnel.
// Version: cowley
func (tunnel) GetProtocol2(session *Session, self TunnelRef) (retval TunnelProtocol, err error) {
	method := "tunnel.get_protocol"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumTunnelProtocol(method+" -> ", result)
	return
}

// GetOtherConfig: Get the other_config field of the given tunnel.
// Version: cowley
func (tunnel) GetOtherConfig(session *Session, self TunnelRef) (retval map[string]string, err error) {
	method := "tunnel.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given tunnel.
// Version: cowley
func (tunnel) GetOtherConfig2(session *Session, self TunnelRef) (retval map[string]string, err error) {
	method := "tunnel.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetStatus: Get the status field of the given tunnel.
// Version: cowley
func (tunnel) GetStatus(session *Session, self TunnelRef) (retval map[string]string, err error) {
	method := "tunnel.get_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetStatus2: Get the status field of the given tunnel.
// Version: cowley
func (tunnel) GetStatus2(session *Session, self TunnelRef) (retval map[string]string, err error) {
	method := "tunnel.get_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTransportPIF: Get the transport_PIF field of the given tunnel.
// Version: cowley
func (tunnel) GetTransportPIF(session *Session, self TunnelRef) (retval PIFRef, err error) {
	method := "tunnel.get_transport_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTransportPIF2: Get the transport_PIF field of the given tunnel.
// Version: cowley
func (tunnel) GetTransportPIF2(session *Session, self TunnelRef) (retval PIFRef, err error) {
	method := "tunnel.get_transport_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAccessPIF: Get the access_PIF field of the given tunnel.
// Version: cowley
func (tunnel) GetAccessPIF(session *Session, self TunnelRef) (retval PIFRef, err error) {
	method := "tunnel.get_access_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAccessPIF2: Get the access_PIF field of the given tunnel.
// Version: cowley
func (tunnel) GetAccessPIF2(session *Session, self TunnelRef) (retval PIFRef, err error) {
	method := "tunnel.get_access_PIF"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given tunnel.
// Version: cowley
func (tunnel) GetUUID(session *Session, self TunnelRef) (retval string, err error) {
	method := "tunnel.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given tunnel.
// Version: cowley
func (tunnel) GetUUID2(session *Session, self TunnelRef) (retval string, err error) {
	method := "tunnel.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the tunnel instance with the specified UUID.
// Version: cowley
func (tunnel) GetByUUID(session *Session, uuid string) (retval TunnelRef, err error) {
	method := "tunnel.get_by_uuid"
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
	retval, err = deserializeTunnelRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the tunnel instance with the specified UUID.
// Version: cowley
func (tunnel) GetByUUID2(session *Session, uuid string) (retval TunnelRef, err error) {
	method := "tunnel.get_by_uuid"
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
	retval, err = deserializeTunnelRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given tunnel.
// Version: cowley
func (tunnel) GetRecord(session *Session, self TunnelRef) (retval TunnelRecord, err error) {
	method := "tunnel.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given tunnel.
// Version: cowley
func (tunnel) GetRecord2(session *Session, self TunnelRef) (retval TunnelRecord, err error) {
	method := "tunnel.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeTunnelRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRecord(method+" -> ", result)
	return
}
