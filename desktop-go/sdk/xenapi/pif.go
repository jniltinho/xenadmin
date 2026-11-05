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

type PIFRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// machine-readable name of the interface (for example, eth0)
	Device string `json:"device,omitempty"`
	// virtual network to which this pif is connected
	Network NetworkRef `json:"network,omitempty"`
	// physical machine to which this pif is connected
	Host HostRef `json:"host,omitempty"`
	// ethernet MAC address of physical interface
	MAC string `json:"MAC,omitempty"`
	// MTU in octets
	MTU int `json:"MTU,omitempty"`
	// VLAN tag for all traffic passing through this interface
	VLAN int `json:"VLAN,omitempty"`
	// metrics associated with this PIF
	Metrics PIFMetricsRef `json:"metrics,omitempty"`
	// true if this represents a physical network interface
	Physical bool `json:"physical,omitempty"`
	// true if this interface is online
	CurrentlyAttached bool `json:"currently_attached,omitempty"`
	// Sets if and how this interface gets an IP address
	IPConfigurationMode IPConfigurationMode `json:"ip_configuration_mode,omitempty"`
	// IP address
	IP string `json:"IP,omitempty"`
	// IP netmask
	Netmask string `json:"netmask,omitempty"`
	// IP gateway
	Gateway string `json:"gateway,omitempty"`
	// Comma separated list of the IP addresses of the DNS servers to use
	DNS string `json:"DNS,omitempty"`
	// Indicates which bond this interface is part of
	BondSlaveOf BondRef `json:"bond_slave_of,omitempty"`
	// Indicates this PIF represents the results of a bond
	BondMasterOf []BondRef `json:"bond_master_of,omitempty"`
	// Indicates which VLAN this interface receives untagged traffic from
	VLANMasterOf VLANRef `json:"VLAN_master_of,omitempty"`
	// Indicates which VLANs this interface transmits tagged traffic to
	VLANSlaveOf []VLANRef `json:"VLAN_slave_of,omitempty"`
	// Indicates whether the control software is listening for connections on this interface
	Management bool `json:"management,omitempty"`
	// Additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// Prevent this PIF from being unplugged; set this to notify the management tool-stack that the PIF has a special use and should not be unplugged under any circumstances (for example, because you&apos;re running storage traffic over it)
	DisallowUnplug bool `json:"disallow_unplug,omitempty"`
	// Indicates to which tunnel this PIF gives access
	TunnelAccessPIFOf []TunnelRef `json:"tunnel_access_PIF_of,omitempty"`
	// Indicates to which tunnel this PIF provides transport
	TunnelTransportPIFOf []TunnelRef `json:"tunnel_transport_PIF_of,omitempty"`
	// Sets if and how this interface gets an IPv6 address
	Ipv6ConfigurationMode Ipv6ConfigurationMode `json:"ipv6_configuration_mode,omitempty"`
	// IPv6 address
	IPv6 []string `json:"IPv6,omitempty"`
	// IPv6 gateway
	Ipv6Gateway string `json:"ipv6_gateway,omitempty"`
	// Which protocol should define the primary address of this interface
	PrimaryAddressType PrimaryAddressType `json:"primary_address_type,omitempty"`
	// Indicates whether the interface is managed by xapi. If it is not, then xapi will not configure the interface, the commands PIF.plug/unplug/reconfigure_ip(v6) cannot be used, nor can the interface be bonded or have VLANs based on top through xapi.
	Managed bool `json:"managed,omitempty"`
	// Additional configuration properties for the interface.
	Properties map[string]string `json:"properties,omitempty"`
	// Additional capabilities on the interface.
	Capabilities []string `json:"capabilities,omitempty"`
	// The IGMP snooping status of the corresponding network bridge
	IgmpSnoopingStatus PifIgmpStatus `json:"igmp_snooping_status,omitempty"`
	// Indicates which network_sriov this interface is physical of
	SriovPhysicalPIFOf []NetworkSriovRef `json:"sriov_physical_PIF_of,omitempty"`
	// Indicates which network_sriov this interface is logical of
	SriovLogicalPIFOf []NetworkSriovRef `json:"sriov_logical_PIF_of,omitempty"`
	// Link to underlying PCI device
	PCI PCIRef `json:"PCI,omitempty"`
}

type PIFRef string

// A physical network interface (note separate VLANs are represented as several PIFs)
type pif struct{}

var PIF pif

// GetAllRecords: Return a map of PIF references to PIF records for all PIFs known to the system.
// Version: rio
func (pif) GetAllRecords(session *Session) (retval map[PIFRef]PIFRecord, err error) {
	method := "PIF.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRefToPIFRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of PIF references to PIF records for all PIFs known to the system.
// Version: rio
func (pif) GetAllRecords1(session *Session) (retval map[PIFRef]PIFRecord, err error) {
	method := "PIF.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRefToPIFRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the PIFs known to the system.
// Version: rio
func (pif) GetAll(session *Session) (retval []PIFRef, err error) {
	method := "PIF.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the PIFs known to the system.
// Version: rio
func (pif) GetAll1(session *Session) (retval []PIFRef, err error) {
	method := "PIF.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRefSet(method+" -> ", result)
	return
}

// SetProperty: Set the value of a property of the PIF
// Version: creedence
func (pif) SetProperty(session *Session, self PIFRef, name string, value string) (err error) {
	method := "PIF.set_property"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetProperty: Set the value of a property of the PIF
// Version: creedence
func (pif) AsyncSetProperty(session *Session, self PIFRef, name string, value string) (retval TaskRef, err error) {
	method := "Async.PIF.set_property"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetProperty4: Set the value of a property of the PIF
// Version: creedence
func (pif) SetProperty4(session *Session, self PIFRef, name string, value string) (err error) {
	method := "PIF.set_property"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetProperty4: Set the value of a property of the PIF
// Version: creedence
func (pif) AsyncSetProperty4(session *Session, self PIFRef, name string, value string) (retval TaskRef, err error) {
	method := "Async.PIF.set_property"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DBForget: Destroy a PIF database record.
// Version: orlando
func (pif) DBForget(session *Session, self PIFRef) (err error) {
	method := "PIF.db_forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDBForget: Destroy a PIF database record.
// Version: orlando
func (pif) AsyncDBForget(session *Session, self PIFRef) (retval TaskRef, err error) {
	method := "Async.PIF.db_forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DBForget2: Destroy a PIF database record.
// Version: orlando
func (pif) DBForget2(session *Session, self PIFRef) (err error) {
	method := "PIF.db_forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDBForget2: Destroy a PIF database record.
// Version: orlando
func (pif) AsyncDBForget2(session *Session, self PIFRef) (retval TaskRef, err error) {
	method := "Async.PIF.db_forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// DBIntroduce: Create a new PIF record in the database only
// Version: creedence
func (pif) DBIntroduce(session *Session, device string, network NetworkRef, host HostRef, mac string, mtu int, vlan int, physical bool, ipConfigurationMode IPConfigurationMode, ip string, netmask string, gateway string, dns string, bondSlaveOf BondRef, vlanMasterOf VLANRef, management bool, otherConfig map[string]string, disallowUnplug bool, ipv6ConfigurationMode Ipv6ConfigurationMode, ipv6 []string, ipv6Gateway string, primaryAddressType PrimaryAddressType, managed bool, properties map[string]string) (retval PIFRef, err error) {
	method := "PIF.db_introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	mtuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "MTU"), mtu)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	physicalArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "physical"), physical)
	if err != nil {
		return
	}
	ipConfigurationModeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "ip_configuration_mode"), ipConfigurationMode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	bondSlaveOfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "bond_slave_of"), bondSlaveOf)
	if err != nil {
		return
	}
	vlanMasterOfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "VLAN_master_of"), vlanMasterOf)
	if err != nil {
		return
	}
	managementArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "management"), management)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	disallowUnplugArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "disallow_unplug"), disallowUnplug)
	if err != nil {
		return
	}
	ipv6ConfigurationModeArg, err := serializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s(%s)", method, "ipv6_configuration_mode"), ipv6ConfigurationMode)
	if err != nil {
		return
	}
	ipv6Arg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "IPv6"), ipv6)
	if err != nil {
		return
	}
	ipv6GatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "ipv6_gateway"), ipv6Gateway)
	if err != nil {
		return
	}
	primaryAddressTypeArg, err := serializeEnumPrimaryAddressType(fmt.Sprintf("%s(%s)", method, "primary_address_type"), primaryAddressType)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	propertiesArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "properties"), properties)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, macArg, mtuArg, vlanArg, physicalArg, ipConfigurationModeArg, ipArg, netmaskArg, gatewayArg, dnsArg, bondSlaveOfArg, vlanMasterOfArg, managementArg, otherConfigArg, disallowUnplugArg, ipv6ConfigurationModeArg, ipv6Arg, ipv6GatewayArg, primaryAddressTypeArg, managedArg, propertiesArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// AsyncDBIntroduce: Create a new PIF record in the database only
// Version: creedence
func (pif) AsyncDBIntroduce(session *Session, device string, network NetworkRef, host HostRef, mac string, mtu int, vlan int, physical bool, ipConfigurationMode IPConfigurationMode, ip string, netmask string, gateway string, dns string, bondSlaveOf BondRef, vlanMasterOf VLANRef, management bool, otherConfig map[string]string, disallowUnplug bool, ipv6ConfigurationMode Ipv6ConfigurationMode, ipv6 []string, ipv6Gateway string, primaryAddressType PrimaryAddressType, managed bool, properties map[string]string) (retval TaskRef, err error) {
	method := "Async.PIF.db_introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	mtuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "MTU"), mtu)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	physicalArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "physical"), physical)
	if err != nil {
		return
	}
	ipConfigurationModeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "ip_configuration_mode"), ipConfigurationMode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	bondSlaveOfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "bond_slave_of"), bondSlaveOf)
	if err != nil {
		return
	}
	vlanMasterOfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "VLAN_master_of"), vlanMasterOf)
	if err != nil {
		return
	}
	managementArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "management"), management)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	disallowUnplugArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "disallow_unplug"), disallowUnplug)
	if err != nil {
		return
	}
	ipv6ConfigurationModeArg, err := serializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s(%s)", method, "ipv6_configuration_mode"), ipv6ConfigurationMode)
	if err != nil {
		return
	}
	ipv6Arg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "IPv6"), ipv6)
	if err != nil {
		return
	}
	ipv6GatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "ipv6_gateway"), ipv6Gateway)
	if err != nil {
		return
	}
	primaryAddressTypeArg, err := serializeEnumPrimaryAddressType(fmt.Sprintf("%s(%s)", method, "primary_address_type"), primaryAddressType)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	propertiesArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "properties"), properties)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, macArg, mtuArg, vlanArg, physicalArg, ipConfigurationModeArg, ipArg, netmaskArg, gatewayArg, dnsArg, bondSlaveOfArg, vlanMasterOfArg, managementArg, otherConfigArg, disallowUnplugArg, ipv6ConfigurationModeArg, ipv6Arg, ipv6GatewayArg, primaryAddressTypeArg, managedArg, propertiesArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DBIntroduce24: Create a new PIF record in the database only
// Version: creedence
func (pif) DBIntroduce24(session *Session, device string, network NetworkRef, host HostRef, mac string, mtu int, vlan int, physical bool, ipConfigurationMode IPConfigurationMode, ip string, netmask string, gateway string, dns string, bondSlaveOf BondRef, vlanMasterOf VLANRef, management bool, otherConfig map[string]string, disallowUnplug bool, ipv6ConfigurationMode Ipv6ConfigurationMode, ipv6 []string, ipv6Gateway string, primaryAddressType PrimaryAddressType, managed bool, properties map[string]string) (retval PIFRef, err error) {
	method := "PIF.db_introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	mtuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "MTU"), mtu)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	physicalArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "physical"), physical)
	if err != nil {
		return
	}
	ipConfigurationModeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "ip_configuration_mode"), ipConfigurationMode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	bondSlaveOfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "bond_slave_of"), bondSlaveOf)
	if err != nil {
		return
	}
	vlanMasterOfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "VLAN_master_of"), vlanMasterOf)
	if err != nil {
		return
	}
	managementArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "management"), management)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	disallowUnplugArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "disallow_unplug"), disallowUnplug)
	if err != nil {
		return
	}
	ipv6ConfigurationModeArg, err := serializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s(%s)", method, "ipv6_configuration_mode"), ipv6ConfigurationMode)
	if err != nil {
		return
	}
	ipv6Arg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "IPv6"), ipv6)
	if err != nil {
		return
	}
	ipv6GatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "ipv6_gateway"), ipv6Gateway)
	if err != nil {
		return
	}
	primaryAddressTypeArg, err := serializeEnumPrimaryAddressType(fmt.Sprintf("%s(%s)", method, "primary_address_type"), primaryAddressType)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	propertiesArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "properties"), properties)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, macArg, mtuArg, vlanArg, physicalArg, ipConfigurationModeArg, ipArg, netmaskArg, gatewayArg, dnsArg, bondSlaveOfArg, vlanMasterOfArg, managementArg, otherConfigArg, disallowUnplugArg, ipv6ConfigurationModeArg, ipv6Arg, ipv6GatewayArg, primaryAddressTypeArg, managedArg, propertiesArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// AsyncDBIntroduce24: Create a new PIF record in the database only
// Version: creedence
func (pif) AsyncDBIntroduce24(session *Session, device string, network NetworkRef, host HostRef, mac string, mtu int, vlan int, physical bool, ipConfigurationMode IPConfigurationMode, ip string, netmask string, gateway string, dns string, bondSlaveOf BondRef, vlanMasterOf VLANRef, management bool, otherConfig map[string]string, disallowUnplug bool, ipv6ConfigurationMode Ipv6ConfigurationMode, ipv6 []string, ipv6Gateway string, primaryAddressType PrimaryAddressType, managed bool, properties map[string]string) (retval TaskRef, err error) {
	method := "Async.PIF.db_introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	mtuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "MTU"), mtu)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	physicalArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "physical"), physical)
	if err != nil {
		return
	}
	ipConfigurationModeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "ip_configuration_mode"), ipConfigurationMode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	bondSlaveOfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "bond_slave_of"), bondSlaveOf)
	if err != nil {
		return
	}
	vlanMasterOfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "VLAN_master_of"), vlanMasterOf)
	if err != nil {
		return
	}
	managementArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "management"), management)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	disallowUnplugArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "disallow_unplug"), disallowUnplug)
	if err != nil {
		return
	}
	ipv6ConfigurationModeArg, err := serializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s(%s)", method, "ipv6_configuration_mode"), ipv6ConfigurationMode)
	if err != nil {
		return
	}
	ipv6Arg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "IPv6"), ipv6)
	if err != nil {
		return
	}
	ipv6GatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "ipv6_gateway"), ipv6Gateway)
	if err != nil {
		return
	}
	primaryAddressTypeArg, err := serializeEnumPrimaryAddressType(fmt.Sprintf("%s(%s)", method, "primary_address_type"), primaryAddressType)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	propertiesArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "properties"), properties)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, macArg, mtuArg, vlanArg, physicalArg, ipConfigurationModeArg, ipArg, netmaskArg, gatewayArg, dnsArg, bondSlaveOfArg, vlanMasterOfArg, managementArg, otherConfigArg, disallowUnplugArg, ipv6ConfigurationModeArg, ipv6Arg, ipv6GatewayArg, primaryAddressTypeArg, managedArg, propertiesArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DBIntroduce23: Create a new PIF record in the database only
// Version: vgpu-productisation
func (pif) DBIntroduce23(session *Session, device string, network NetworkRef, host HostRef, mac string, mtu int, vlan int, physical bool, ipConfigurationMode IPConfigurationMode, ip string, netmask string, gateway string, dns string, bondSlaveOf BondRef, vlanMasterOf VLANRef, management bool, otherConfig map[string]string, disallowUnplug bool, ipv6ConfigurationMode Ipv6ConfigurationMode, ipv6 []string, ipv6Gateway string, primaryAddressType PrimaryAddressType, managed bool) (retval PIFRef, err error) {
	method := "PIF.db_introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	mtuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "MTU"), mtu)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	physicalArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "physical"), physical)
	if err != nil {
		return
	}
	ipConfigurationModeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "ip_configuration_mode"), ipConfigurationMode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	bondSlaveOfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "bond_slave_of"), bondSlaveOf)
	if err != nil {
		return
	}
	vlanMasterOfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "VLAN_master_of"), vlanMasterOf)
	if err != nil {
		return
	}
	managementArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "management"), management)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	disallowUnplugArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "disallow_unplug"), disallowUnplug)
	if err != nil {
		return
	}
	ipv6ConfigurationModeArg, err := serializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s(%s)", method, "ipv6_configuration_mode"), ipv6ConfigurationMode)
	if err != nil {
		return
	}
	ipv6Arg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "IPv6"), ipv6)
	if err != nil {
		return
	}
	ipv6GatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "ipv6_gateway"), ipv6Gateway)
	if err != nil {
		return
	}
	primaryAddressTypeArg, err := serializeEnumPrimaryAddressType(fmt.Sprintf("%s(%s)", method, "primary_address_type"), primaryAddressType)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, macArg, mtuArg, vlanArg, physicalArg, ipConfigurationModeArg, ipArg, netmaskArg, gatewayArg, dnsArg, bondSlaveOfArg, vlanMasterOfArg, managementArg, otherConfigArg, disallowUnplugArg, ipv6ConfigurationModeArg, ipv6Arg, ipv6GatewayArg, primaryAddressTypeArg, managedArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// AsyncDBIntroduce23: Create a new PIF record in the database only
// Version: vgpu-productisation
func (pif) AsyncDBIntroduce23(session *Session, device string, network NetworkRef, host HostRef, mac string, mtu int, vlan int, physical bool, ipConfigurationMode IPConfigurationMode, ip string, netmask string, gateway string, dns string, bondSlaveOf BondRef, vlanMasterOf VLANRef, management bool, otherConfig map[string]string, disallowUnplug bool, ipv6ConfigurationMode Ipv6ConfigurationMode, ipv6 []string, ipv6Gateway string, primaryAddressType PrimaryAddressType, managed bool) (retval TaskRef, err error) {
	method := "Async.PIF.db_introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	mtuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "MTU"), mtu)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	physicalArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "physical"), physical)
	if err != nil {
		return
	}
	ipConfigurationModeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "ip_configuration_mode"), ipConfigurationMode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	bondSlaveOfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "bond_slave_of"), bondSlaveOf)
	if err != nil {
		return
	}
	vlanMasterOfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "VLAN_master_of"), vlanMasterOf)
	if err != nil {
		return
	}
	managementArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "management"), management)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	disallowUnplugArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "disallow_unplug"), disallowUnplug)
	if err != nil {
		return
	}
	ipv6ConfigurationModeArg, err := serializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s(%s)", method, "ipv6_configuration_mode"), ipv6ConfigurationMode)
	if err != nil {
		return
	}
	ipv6Arg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "IPv6"), ipv6)
	if err != nil {
		return
	}
	ipv6GatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "ipv6_gateway"), ipv6Gateway)
	if err != nil {
		return
	}
	primaryAddressTypeArg, err := serializeEnumPrimaryAddressType(fmt.Sprintf("%s(%s)", method, "primary_address_type"), primaryAddressType)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, macArg, mtuArg, vlanArg, physicalArg, ipConfigurationModeArg, ipArg, netmaskArg, gatewayArg, dnsArg, bondSlaveOfArg, vlanMasterOfArg, managementArg, otherConfigArg, disallowUnplugArg, ipv6ConfigurationModeArg, ipv6Arg, ipv6GatewayArg, primaryAddressTypeArg, managedArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DBIntroduce22: Create a new PIF record in the database only
// Version: boston
func (pif) DBIntroduce22(session *Session, device string, network NetworkRef, host HostRef, mac string, mtu int, vlan int, physical bool, ipConfigurationMode IPConfigurationMode, ip string, netmask string, gateway string, dns string, bondSlaveOf BondRef, vlanMasterOf VLANRef, management bool, otherConfig map[string]string, disallowUnplug bool, ipv6ConfigurationMode Ipv6ConfigurationMode, ipv6 []string, ipv6Gateway string, primaryAddressType PrimaryAddressType) (retval PIFRef, err error) {
	method := "PIF.db_introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	mtuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "MTU"), mtu)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	physicalArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "physical"), physical)
	if err != nil {
		return
	}
	ipConfigurationModeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "ip_configuration_mode"), ipConfigurationMode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	bondSlaveOfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "bond_slave_of"), bondSlaveOf)
	if err != nil {
		return
	}
	vlanMasterOfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "VLAN_master_of"), vlanMasterOf)
	if err != nil {
		return
	}
	managementArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "management"), management)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	disallowUnplugArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "disallow_unplug"), disallowUnplug)
	if err != nil {
		return
	}
	ipv6ConfigurationModeArg, err := serializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s(%s)", method, "ipv6_configuration_mode"), ipv6ConfigurationMode)
	if err != nil {
		return
	}
	ipv6Arg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "IPv6"), ipv6)
	if err != nil {
		return
	}
	ipv6GatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "ipv6_gateway"), ipv6Gateway)
	if err != nil {
		return
	}
	primaryAddressTypeArg, err := serializeEnumPrimaryAddressType(fmt.Sprintf("%s(%s)", method, "primary_address_type"), primaryAddressType)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, macArg, mtuArg, vlanArg, physicalArg, ipConfigurationModeArg, ipArg, netmaskArg, gatewayArg, dnsArg, bondSlaveOfArg, vlanMasterOfArg, managementArg, otherConfigArg, disallowUnplugArg, ipv6ConfigurationModeArg, ipv6Arg, ipv6GatewayArg, primaryAddressTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// AsyncDBIntroduce22: Create a new PIF record in the database only
// Version: boston
func (pif) AsyncDBIntroduce22(session *Session, device string, network NetworkRef, host HostRef, mac string, mtu int, vlan int, physical bool, ipConfigurationMode IPConfigurationMode, ip string, netmask string, gateway string, dns string, bondSlaveOf BondRef, vlanMasterOf VLANRef, management bool, otherConfig map[string]string, disallowUnplug bool, ipv6ConfigurationMode Ipv6ConfigurationMode, ipv6 []string, ipv6Gateway string, primaryAddressType PrimaryAddressType) (retval TaskRef, err error) {
	method := "Async.PIF.db_introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	mtuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "MTU"), mtu)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	physicalArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "physical"), physical)
	if err != nil {
		return
	}
	ipConfigurationModeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "ip_configuration_mode"), ipConfigurationMode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	bondSlaveOfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "bond_slave_of"), bondSlaveOf)
	if err != nil {
		return
	}
	vlanMasterOfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "VLAN_master_of"), vlanMasterOf)
	if err != nil {
		return
	}
	managementArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "management"), management)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	disallowUnplugArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "disallow_unplug"), disallowUnplug)
	if err != nil {
		return
	}
	ipv6ConfigurationModeArg, err := serializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s(%s)", method, "ipv6_configuration_mode"), ipv6ConfigurationMode)
	if err != nil {
		return
	}
	ipv6Arg, err := serializeStringSet(fmt.Sprintf("%s(%s)", method, "IPv6"), ipv6)
	if err != nil {
		return
	}
	ipv6GatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "ipv6_gateway"), ipv6Gateway)
	if err != nil {
		return
	}
	primaryAddressTypeArg, err := serializeEnumPrimaryAddressType(fmt.Sprintf("%s(%s)", method, "primary_address_type"), primaryAddressType)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, macArg, mtuArg, vlanArg, physicalArg, ipConfigurationModeArg, ipArg, netmaskArg, gatewayArg, dnsArg, bondSlaveOfArg, vlanMasterOfArg, managementArg, otherConfigArg, disallowUnplugArg, ipv6ConfigurationModeArg, ipv6Arg, ipv6GatewayArg, primaryAddressTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// DBIntroduce18: Create a new PIF record in the database only
// Version: orlando
func (pif) DBIntroduce18(session *Session, device string, network NetworkRef, host HostRef, mac string, mtu int, vlan int, physical bool, ipConfigurationMode IPConfigurationMode, ip string, netmask string, gateway string, dns string, bondSlaveOf BondRef, vlanMasterOf VLANRef, management bool, otherConfig map[string]string, disallowUnplug bool) (retval PIFRef, err error) {
	method := "PIF.db_introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	mtuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "MTU"), mtu)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	physicalArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "physical"), physical)
	if err != nil {
		return
	}
	ipConfigurationModeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "ip_configuration_mode"), ipConfigurationMode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	bondSlaveOfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "bond_slave_of"), bondSlaveOf)
	if err != nil {
		return
	}
	vlanMasterOfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "VLAN_master_of"), vlanMasterOf)
	if err != nil {
		return
	}
	managementArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "management"), management)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	disallowUnplugArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "disallow_unplug"), disallowUnplug)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, macArg, mtuArg, vlanArg, physicalArg, ipConfigurationModeArg, ipArg, netmaskArg, gatewayArg, dnsArg, bondSlaveOfArg, vlanMasterOfArg, managementArg, otherConfigArg, disallowUnplugArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// AsyncDBIntroduce18: Create a new PIF record in the database only
// Version: orlando
func (pif) AsyncDBIntroduce18(session *Session, device string, network NetworkRef, host HostRef, mac string, mtu int, vlan int, physical bool, ipConfigurationMode IPConfigurationMode, ip string, netmask string, gateway string, dns string, bondSlaveOf BondRef, vlanMasterOf VLANRef, management bool, otherConfig map[string]string, disallowUnplug bool) (retval TaskRef, err error) {
	method := "Async.PIF.db_introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	mtuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "MTU"), mtu)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	physicalArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "physical"), physical)
	if err != nil {
		return
	}
	ipConfigurationModeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "ip_configuration_mode"), ipConfigurationMode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	bondSlaveOfArg, err := serializeBondRef(fmt.Sprintf("%s(%s)", method, "bond_slave_of"), bondSlaveOf)
	if err != nil {
		return
	}
	vlanMasterOfArg, err := serializeVLANRef(fmt.Sprintf("%s(%s)", method, "VLAN_master_of"), vlanMasterOf)
	if err != nil {
		return
	}
	managementArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "management"), management)
	if err != nil {
		return
	}
	otherConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "other_config"), otherConfig)
	if err != nil {
		return
	}
	disallowUnplugArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "disallow_unplug"), disallowUnplug)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, macArg, mtuArg, vlanArg, physicalArg, ipConfigurationModeArg, ipArg, netmaskArg, gatewayArg, dnsArg, bondSlaveOfArg, vlanMasterOfArg, managementArg, otherConfigArg, disallowUnplugArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Plug: Attempt to bring up a physical interface
// Version: miami
//
// Errors:
// TRANSPORT_PIF_NOT_CONFIGURED - The tunnel transport PIF has no IP configuration set.
func (pif) Plug(session *Session, self PIFRef) (err error) {
	method := "PIF.plug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPlug: Attempt to bring up a physical interface
// Version: miami
//
// Errors:
// TRANSPORT_PIF_NOT_CONFIGURED - The tunnel transport PIF has no IP configuration set.
func (pif) AsyncPlug(session *Session, self PIFRef) (retval TaskRef, err error) {
	method := "Async.PIF.plug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Plug2: Attempt to bring up a physical interface
// Version: miami
//
// Errors:
// TRANSPORT_PIF_NOT_CONFIGURED - The tunnel transport PIF has no IP configuration set.
func (pif) Plug2(session *Session, self PIFRef) (err error) {
	method := "PIF.plug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncPlug2: Attempt to bring up a physical interface
// Version: miami
//
// Errors:
// TRANSPORT_PIF_NOT_CONFIGURED - The tunnel transport PIF has no IP configuration set.
func (pif) AsyncPlug2(session *Session, self PIFRef) (retval TaskRef, err error) {
	method := "Async.PIF.plug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetDisallowUnplug: Set whether unplugging the PIF is allowed
// Version: orlando
//
// Errors:
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) SetDisallowUnplug(session *Session, self PIFRef, value bool) (err error) {
	method := "PIF.set_disallow_unplug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetDisallowUnplug: Set whether unplugging the PIF is allowed
// Version: orlando
//
// Errors:
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) AsyncSetDisallowUnplug(session *Session, self PIFRef, value bool) (retval TaskRef, err error) {
	method := "Async.PIF.set_disallow_unplug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetDisallowUnplug3: Set whether unplugging the PIF is allowed
// Version: orlando
//
// Errors:
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) SetDisallowUnplug3(session *Session, self PIFRef, value bool) (err error) {
	method := "PIF.set_disallow_unplug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetDisallowUnplug3: Set whether unplugging the PIF is allowed
// Version: orlando
//
// Errors:
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) AsyncSetDisallowUnplug3(session *Session, self PIFRef, value bool) (retval TaskRef, err error) {
	method := "Async.PIF.set_disallow_unplug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Unplug: Attempt to bring down a physical interface
// Version: miami
//
// Errors:
// HA_OPERATION_WOULD_BREAK_FAILOVER_PLAN - This operation cannot be performed because it would invalidate VM failover planning such that the system would be unable to guarantee to restart protected VMs after a Host failure.
// VIF_IN_USE - Network has active VIFs
// PIF_DOES_NOT_ALLOW_UNPLUG - The operation you requested cannot be performed because the specified PIF does not allow unplug.
// PIF_HAS_FCOE_SR_IN_USE - The operation you requested cannot be performed because the specified PIF has FCoE SR in use.
func (pif) Unplug(session *Session, self PIFRef) (err error) {
	method := "PIF.unplug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncUnplug: Attempt to bring down a physical interface
// Version: miami
//
// Errors:
// HA_OPERATION_WOULD_BREAK_FAILOVER_PLAN - This operation cannot be performed because it would invalidate VM failover planning such that the system would be unable to guarantee to restart protected VMs after a Host failure.
// VIF_IN_USE - Network has active VIFs
// PIF_DOES_NOT_ALLOW_UNPLUG - The operation you requested cannot be performed because the specified PIF does not allow unplug.
// PIF_HAS_FCOE_SR_IN_USE - The operation you requested cannot be performed because the specified PIF has FCoE SR in use.
func (pif) AsyncUnplug(session *Session, self PIFRef) (retval TaskRef, err error) {
	method := "Async.PIF.unplug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Unplug2: Attempt to bring down a physical interface
// Version: miami
//
// Errors:
// HA_OPERATION_WOULD_BREAK_FAILOVER_PLAN - This operation cannot be performed because it would invalidate VM failover planning such that the system would be unable to guarantee to restart protected VMs after a Host failure.
// VIF_IN_USE - Network has active VIFs
// PIF_DOES_NOT_ALLOW_UNPLUG - The operation you requested cannot be performed because the specified PIF does not allow unplug.
// PIF_HAS_FCOE_SR_IN_USE - The operation you requested cannot be performed because the specified PIF has FCoE SR in use.
func (pif) Unplug2(session *Session, self PIFRef) (err error) {
	method := "PIF.unplug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncUnplug2: Attempt to bring down a physical interface
// Version: miami
//
// Errors:
// HA_OPERATION_WOULD_BREAK_FAILOVER_PLAN - This operation cannot be performed because it would invalidate VM failover planning such that the system would be unable to guarantee to restart protected VMs after a Host failure.
// VIF_IN_USE - Network has active VIFs
// PIF_DOES_NOT_ALLOW_UNPLUG - The operation you requested cannot be performed because the specified PIF does not allow unplug.
// PIF_HAS_FCOE_SR_IN_USE - The operation you requested cannot be performed because the specified PIF has FCoE SR in use.
func (pif) AsyncUnplug2(session *Session, self PIFRef) (retval TaskRef, err error) {
	method := "Async.PIF.unplug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Forget: Destroy the PIF object matching a particular network interface
// Version: miami
//
// Errors:
// PIF_TUNNEL_STILL_EXISTS - Operation cannot proceed while a tunnel exists on this interface.
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) Forget(session *Session, self PIFRef) (err error) {
	method := "PIF.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncForget: Destroy the PIF object matching a particular network interface
// Version: miami
//
// Errors:
// PIF_TUNNEL_STILL_EXISTS - Operation cannot proceed while a tunnel exists on this interface.
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) AsyncForget(session *Session, self PIFRef) (retval TaskRef, err error) {
	method := "Async.PIF.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Forget2: Destroy the PIF object matching a particular network interface
// Version: miami
//
// Errors:
// PIF_TUNNEL_STILL_EXISTS - Operation cannot proceed while a tunnel exists on this interface.
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) Forget2(session *Session, self PIFRef) (err error) {
	method := "PIF.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncForget2: Destroy the PIF object matching a particular network interface
// Version: miami
//
// Errors:
// PIF_TUNNEL_STILL_EXISTS - Operation cannot proceed while a tunnel exists on this interface.
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) AsyncForget2(session *Session, self PIFRef) (retval TaskRef, err error) {
	method := "Async.PIF.forget"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Introduce: Create a PIF object matching a particular network interface
// Version: vgpu-productisation
func (pif) Introduce(session *Session, host HostRef, mac string, device string, managed bool) (retval PIFRef, err error) {
	method := "PIF.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, macArg, deviceArg, managedArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// AsyncIntroduce: Create a PIF object matching a particular network interface
// Version: vgpu-productisation
func (pif) AsyncIntroduce(session *Session, host HostRef, mac string, device string, managed bool) (retval TaskRef, err error) {
	method := "Async.PIF.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, macArg, deviceArg, managedArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce5: Create a PIF object matching a particular network interface
// Version: vgpu-productisation
func (pif) Introduce5(session *Session, host HostRef, mac string, device string, managed bool) (retval PIFRef, err error) {
	method := "PIF.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, macArg, deviceArg, managedArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// AsyncIntroduce5: Create a PIF object matching a particular network interface
// Version: vgpu-productisation
func (pif) AsyncIntroduce5(session *Session, host HostRef, mac string, device string, managed bool) (retval TaskRef, err error) {
	method := "Async.PIF.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	managedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "managed"), managed)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, macArg, deviceArg, managedArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Introduce4: Create a PIF object matching a particular network interface
// Version: miami
func (pif) Introduce4(session *Session, host HostRef, mac string, device string) (retval PIFRef, err error) {
	method := "PIF.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, macArg, deviceArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// AsyncIntroduce4: Create a PIF object matching a particular network interface
// Version: miami
func (pif) AsyncIntroduce4(session *Session, host HostRef, mac string, device string) (retval TaskRef, err error) {
	method := "Async.PIF.introduce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	macArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "MAC"), mac)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, hostArg, macArg, deviceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Scan: Scan for physical interfaces on a host and create PIF objects to represent them
// Version: miami
func (pif) Scan(session *Session, host HostRef) (err error) {
	method := "PIF.scan"
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

// AsyncScan: Scan for physical interfaces on a host and create PIF objects to represent them
// Version: miami
func (pif) AsyncScan(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.PIF.scan"
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

// Scan2: Scan for physical interfaces on a host and create PIF objects to represent them
// Version: miami
func (pif) Scan2(session *Session, host HostRef) (err error) {
	method := "PIF.scan"
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

// AsyncScan2: Scan for physical interfaces on a host and create PIF objects to represent them
// Version: miami
func (pif) AsyncScan2(session *Session, host HostRef) (retval TaskRef, err error) {
	method := "Async.PIF.scan"
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

// SetPrimaryAddressType: Change the primary address type used by this PIF
// Version: tampa
func (pif) SetPrimaryAddressType(session *Session, self PIFRef, primaryAddressType PrimaryAddressType) (err error) {
	method := "PIF.set_primary_address_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	primaryAddressTypeArg, err := serializeEnumPrimaryAddressType(fmt.Sprintf("%s(%s)", method, "primary_address_type"), primaryAddressType)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, primaryAddressTypeArg)
	return
}

// AsyncSetPrimaryAddressType: Change the primary address type used by this PIF
// Version: tampa
func (pif) AsyncSetPrimaryAddressType(session *Session, self PIFRef, primaryAddressType PrimaryAddressType) (retval TaskRef, err error) {
	method := "Async.PIF.set_primary_address_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	primaryAddressTypeArg, err := serializeEnumPrimaryAddressType(fmt.Sprintf("%s(%s)", method, "primary_address_type"), primaryAddressType)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, primaryAddressTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetPrimaryAddressType3: Change the primary address type used by this PIF
// Version: tampa
func (pif) SetPrimaryAddressType3(session *Session, self PIFRef, primaryAddressType PrimaryAddressType) (err error) {
	method := "PIF.set_primary_address_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	primaryAddressTypeArg, err := serializeEnumPrimaryAddressType(fmt.Sprintf("%s(%s)", method, "primary_address_type"), primaryAddressType)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, primaryAddressTypeArg)
	return
}

// AsyncSetPrimaryAddressType3: Change the primary address type used by this PIF
// Version: tampa
func (pif) AsyncSetPrimaryAddressType3(session *Session, self PIFRef, primaryAddressType PrimaryAddressType) (retval TaskRef, err error) {
	method := "Async.PIF.set_primary_address_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	primaryAddressTypeArg, err := serializeEnumPrimaryAddressType(fmt.Sprintf("%s(%s)", method, "primary_address_type"), primaryAddressType)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, primaryAddressTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ReconfigureIpv6: Reconfigure the IPv6 address settings for this interface
// Version: tampa
//
// Errors:
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) ReconfigureIpv6(session *Session, self PIFRef, mode Ipv6ConfigurationMode, ipv6 string, gateway string, dns string) (err error) {
	method := "PIF.reconfigure_ipv6"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	ipv6Arg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IPv6"), ipv6)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, modeArg, ipv6Arg, gatewayArg, dnsArg)
	return
}

// AsyncReconfigureIpv6: Reconfigure the IPv6 address settings for this interface
// Version: tampa
//
// Errors:
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) AsyncReconfigureIpv6(session *Session, self PIFRef, mode Ipv6ConfigurationMode, ipv6 string, gateway string, dns string) (retval TaskRef, err error) {
	method := "Async.PIF.reconfigure_ipv6"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	ipv6Arg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IPv6"), ipv6)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, modeArg, ipv6Arg, gatewayArg, dnsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ReconfigureIpv66: Reconfigure the IPv6 address settings for this interface
// Version: tampa
//
// Errors:
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) ReconfigureIpv66(session *Session, self PIFRef, mode Ipv6ConfigurationMode, ipv6 string, gateway string, dns string) (err error) {
	method := "PIF.reconfigure_ipv6"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	ipv6Arg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IPv6"), ipv6)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, modeArg, ipv6Arg, gatewayArg, dnsArg)
	return
}

// AsyncReconfigureIpv66: Reconfigure the IPv6 address settings for this interface
// Version: tampa
//
// Errors:
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) AsyncReconfigureIpv66(session *Session, self PIFRef, mode Ipv6ConfigurationMode, ipv6 string, gateway string, dns string) (retval TaskRef, err error) {
	method := "Async.PIF.reconfigure_ipv6"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	ipv6Arg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IPv6"), ipv6)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, modeArg, ipv6Arg, gatewayArg, dnsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ReconfigureIP: Reconfigure the IP address settings for this interface
// Version: miami
//
// Errors:
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) ReconfigureIP(session *Session, self PIFRef, mode IPConfigurationMode, ip string, netmask string, gateway string, dns string) (err error) {
	method := "PIF.reconfigure_ip"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, modeArg, ipArg, netmaskArg, gatewayArg, dnsArg)
	return
}

// AsyncReconfigureIP: Reconfigure the IP address settings for this interface
// Version: miami
//
// Errors:
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) AsyncReconfigureIP(session *Session, self PIFRef, mode IPConfigurationMode, ip string, netmask string, gateway string, dns string) (retval TaskRef, err error) {
	method := "Async.PIF.reconfigure_ip"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, modeArg, ipArg, netmaskArg, gatewayArg, dnsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ReconfigureIP7: Reconfigure the IP address settings for this interface
// Version: miami
//
// Errors:
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) ReconfigureIP7(session *Session, self PIFRef, mode IPConfigurationMode, ip string, netmask string, gateway string, dns string) (err error) {
	method := "PIF.reconfigure_ip"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, modeArg, ipArg, netmaskArg, gatewayArg, dnsArg)
	return
}

// AsyncReconfigureIP7: Reconfigure the IP address settings for this interface
// Version: miami
//
// Errors:
// CLUSTERING_ENABLED - An operation was attempted while clustering was enabled on the cluster_host.
func (pif) AsyncReconfigureIP7(session *Session, self PIFRef, mode IPConfigurationMode, ip string, netmask string, gateway string, dns string) (retval TaskRef, err error) {
	method := "Async.PIF.reconfigure_ip"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumIPConfigurationMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	ipArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "IP"), ip)
	if err != nil {
		return
	}
	netmaskArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "netmask"), netmask)
	if err != nil {
		return
	}
	gatewayArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "gateway"), gateway)
	if err != nil {
		return
	}
	dnsArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "DNS"), dns)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, modeArg, ipArg, netmaskArg, gatewayArg, dnsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Destroy: Destroy the PIF object (provided it is a VLAN interface). This call is deprecated: use VLAN.destroy or Bond.destroy instead
// Version: rio
//
// Errors:
// PIF_IS_PHYSICAL - You tried to destroy a PIF, but it represents an aspect of the physical host configuration, and so cannot be destroyed. The parameter echoes the PIF handle you gave.
func (pif) Destroy(session *Session, self PIFRef) (err error) {
	method := "PIF.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy the PIF object (provided it is a VLAN interface). This call is deprecated: use VLAN.destroy or Bond.destroy instead
// Version: rio
//
// Errors:
// PIF_IS_PHYSICAL - You tried to destroy a PIF, but it represents an aspect of the physical host configuration, and so cannot be destroyed. The parameter echoes the PIF handle you gave.
func (pif) AsyncDestroy(session *Session, self PIFRef) (retval TaskRef, err error) {
	method := "Async.PIF.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy the PIF object (provided it is a VLAN interface). This call is deprecated: use VLAN.destroy or Bond.destroy instead
// Version: rio
//
// Errors:
// PIF_IS_PHYSICAL - You tried to destroy a PIF, but it represents an aspect of the physical host configuration, and so cannot be destroyed. The parameter echoes the PIF handle you gave.
func (pif) Destroy2(session *Session, self PIFRef) (err error) {
	method := "PIF.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy the PIF object (provided it is a VLAN interface). This call is deprecated: use VLAN.destroy or Bond.destroy instead
// Version: rio
//
// Errors:
// PIF_IS_PHYSICAL - You tried to destroy a PIF, but it represents an aspect of the physical host configuration, and so cannot be destroyed. The parameter echoes the PIF handle you gave.
func (pif) AsyncDestroy2(session *Session, self PIFRef) (retval TaskRef, err error) {
	method := "Async.PIF.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// CreateVLAN: Create a VLAN interface from an existing physical interface. This call is deprecated: use VLAN.create instead
// Version: rio
//
// Errors:
// VLAN_TAG_INVALID - You tried to create a VLAN, but the tag you gave was invalid -- it must be between 0 and 4094. The parameter echoes the VLAN tag you gave.
func (pif) CreateVLAN(session *Session, device string, network NetworkRef, host HostRef, vlan int) (retval PIFRef, err error) {
	method := "PIF.create_VLAN"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, vlanArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// AsyncCreateVLAN: Create a VLAN interface from an existing physical interface. This call is deprecated: use VLAN.create instead
// Version: rio
//
// Errors:
// VLAN_TAG_INVALID - You tried to create a VLAN, but the tag you gave was invalid -- it must be between 0 and 4094. The parameter echoes the VLAN tag you gave.
func (pif) AsyncCreateVLAN(session *Session, device string, network NetworkRef, host HostRef, vlan int) (retval TaskRef, err error) {
	method := "Async.PIF.create_VLAN"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, vlanArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateVLAN5: Create a VLAN interface from an existing physical interface. This call is deprecated: use VLAN.create instead
// Version: rio
//
// Errors:
// VLAN_TAG_INVALID - You tried to create a VLAN, but the tag you gave was invalid -- it must be between 0 and 4094. The parameter echoes the VLAN tag you gave.
func (pif) CreateVLAN5(session *Session, device string, network NetworkRef, host HostRef, vlan int) (retval PIFRef, err error) {
	method := "PIF.create_VLAN"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, vlanArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// AsyncCreateVLAN5: Create a VLAN interface from an existing physical interface. This call is deprecated: use VLAN.create instead
// Version: rio
//
// Errors:
// VLAN_TAG_INVALID - You tried to create a VLAN, but the tag you gave was invalid -- it must be between 0 and 4094. The parameter echoes the VLAN tag you gave.
func (pif) AsyncCreateVLAN5(session *Session, device string, network NetworkRef, host HostRef, vlan int) (retval TaskRef, err error) {
	method := "Async.PIF.create_VLAN"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	deviceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "device"), device)
	if err != nil {
		return
	}
	networkArg, err := serializeNetworkRef(fmt.Sprintf("%s(%s)", method, "network"), network)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	vlanArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "VLAN"), vlan)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, deviceArg, networkArg, hostArg, vlanArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given PIF.  If the key is not in that Map, then do nothing.
// Version: miami
func (pif) RemoveFromOtherConfig(session *Session, self PIFRef, key string) (err error) {
	method := "PIF.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given PIF.  If the key is not in that Map, then do nothing.
// Version: miami
func (pif) RemoveFromOtherConfig3(session *Session, self PIFRef, key string) (err error) {
	method := "PIF.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig2: Remove the given key and its corresponding value from the other_config field of the given PIF.  If the key is not in that Map, then do nothing.
// Version: rio
func (pif) RemoveFromOtherConfig2(session *Session, self PIFRef) (err error) {
	method := "PIF.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given PIF.
// Version: miami
func (pif) AddToOtherConfig(session *Session, self PIFRef, key string, value string) (err error) {
	method := "PIF.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given PIF.
// Version: miami
func (pif) AddToOtherConfig4(session *Session, self PIFRef, key string, value string) (err error) {
	method := "PIF.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig2: Add the given key-value pair to the other_config field of the given PIF.
// Version: rio
func (pif) AddToOtherConfig2(session *Session, self PIFRef) (err error) {
	method := "PIF.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetOtherConfig: Set the other_config field of the given PIF.
// Version: miami
func (pif) SetOtherConfig(session *Session, self PIFRef, value map[string]string) (err error) {
	method := "PIF.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given PIF.
// Version: miami
func (pif) SetOtherConfig3(session *Session, self PIFRef, value map[string]string) (err error) {
	method := "PIF.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig2: Set the other_config field of the given PIF.
// Version: rio
func (pif) SetOtherConfig2(session *Session, self PIFRef) (err error) {
	method := "PIF.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// GetPCI: Get the PCI field of the given PIF.
// Version: rio
func (pif) GetPCI(session *Session, self PIFRef) (retval PCIRef, err error) {
	method := "PIF.get_PCI"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRef(method+" -> ", result)
	return
}

// GetPCI2: Get the PCI field of the given PIF.
// Version: rio
func (pif) GetPCI2(session *Session, self PIFRef) (retval PCIRef, err error) {
	method := "PIF.get_PCI"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRef(method+" -> ", result)
	return
}

// GetSriovLogicalPIFOf: Get the sriov_logical_PIF_of field of the given PIF.
// Version: rio
func (pif) GetSriovLogicalPIFOf(session *Session, self PIFRef) (retval []NetworkSriovRef, err error) {
	method := "PIF.get_sriov_logical_PIF_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkSriovRefSet(method+" -> ", result)
	return
}

// GetSriovLogicalPIFOf2: Get the sriov_logical_PIF_of field of the given PIF.
// Version: rio
func (pif) GetSriovLogicalPIFOf2(session *Session, self PIFRef) (retval []NetworkSriovRef, err error) {
	method := "PIF.get_sriov_logical_PIF_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkSriovRefSet(method+" -> ", result)
	return
}

// GetSriovPhysicalPIFOf: Get the sriov_physical_PIF_of field of the given PIF.
// Version: rio
func (pif) GetSriovPhysicalPIFOf(session *Session, self PIFRef) (retval []NetworkSriovRef, err error) {
	method := "PIF.get_sriov_physical_PIF_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkSriovRefSet(method+" -> ", result)
	return
}

// GetSriovPhysicalPIFOf2: Get the sriov_physical_PIF_of field of the given PIF.
// Version: rio
func (pif) GetSriovPhysicalPIFOf2(session *Session, self PIFRef) (retval []NetworkSriovRef, err error) {
	method := "PIF.get_sriov_physical_PIF_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeNetworkSriovRefSet(method+" -> ", result)
	return
}

// GetIgmpSnoopingStatus: Get the igmp_snooping_status field of the given PIF.
// Version: rio
func (pif) GetIgmpSnoopingStatus(session *Session, self PIFRef) (retval PifIgmpStatus, err error) {
	method := "PIF.get_igmp_snooping_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPifIgmpStatus(method+" -> ", result)
	return
}

// GetIgmpSnoopingStatus2: Get the igmp_snooping_status field of the given PIF.
// Version: rio
func (pif) GetIgmpSnoopingStatus2(session *Session, self PIFRef) (retval PifIgmpStatus, err error) {
	method := "PIF.get_igmp_snooping_status"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPifIgmpStatus(method+" -> ", result)
	return
}

// GetCapabilities: Get the capabilities field of the given PIF.
// Version: rio
func (pif) GetCapabilities(session *Session, self PIFRef) (retval []string, err error) {
	method := "PIF.get_capabilities"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCapabilities2: Get the capabilities field of the given PIF.
// Version: rio
func (pif) GetCapabilities2(session *Session, self PIFRef) (retval []string, err error) {
	method := "PIF.get_capabilities"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetProperties: Get the properties field of the given PIF.
// Version: rio
func (pif) GetProperties(session *Session, self PIFRef) (retval map[string]string, err error) {
	method := "PIF.get_properties"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetProperties2: Get the properties field of the given PIF.
// Version: rio
func (pif) GetProperties2(session *Session, self PIFRef) (retval map[string]string, err error) {
	method := "PIF.get_properties"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetManaged: Get the managed field of the given PIF.
// Version: rio
func (pif) GetManaged(session *Session, self PIFRef) (retval bool, err error) {
	method := "PIF.get_managed"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetManaged2: Get the managed field of the given PIF.
// Version: rio
func (pif) GetManaged2(session *Session, self PIFRef) (retval bool, err error) {
	method := "PIF.get_managed"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPrimaryAddressType: Get the primary_address_type field of the given PIF.
// Version: rio
func (pif) GetPrimaryAddressType(session *Session, self PIFRef) (retval PrimaryAddressType, err error) {
	method := "PIF.get_primary_address_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPrimaryAddressType(method+" -> ", result)
	return
}

// GetPrimaryAddressType2: Get the primary_address_type field of the given PIF.
// Version: rio
func (pif) GetPrimaryAddressType2(session *Session, self PIFRef) (retval PrimaryAddressType, err error) {
	method := "PIF.get_primary_address_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumPrimaryAddressType(method+" -> ", result)
	return
}

// GetIpv6Gateway: Get the ipv6_gateway field of the given PIF.
// Version: rio
func (pif) GetIpv6Gateway(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_ipv6_gateway"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIpv6Gateway2: Get the ipv6_gateway field of the given PIF.
// Version: rio
func (pif) GetIpv6Gateway2(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_ipv6_gateway"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIPv6: Get the IPv6 field of the given PIF.
// Version: rio
func (pif) GetIPv6(session *Session, self PIFRef) (retval []string, err error) {
	method := "PIF.get_IPv6"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIPv62: Get the IPv6 field of the given PIF.
// Version: rio
func (pif) GetIPv62(session *Session, self PIFRef) (retval []string, err error) {
	method := "PIF.get_IPv6"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIpv6ConfigurationMode: Get the ipv6_configuration_mode field of the given PIF.
// Version: rio
func (pif) GetIpv6ConfigurationMode(session *Session, self PIFRef) (retval Ipv6ConfigurationMode, err error) {
	method := "PIF.get_ipv6_configuration_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumIpv6ConfigurationMode(method+" -> ", result)
	return
}

// GetIpv6ConfigurationMode2: Get the ipv6_configuration_mode field of the given PIF.
// Version: rio
func (pif) GetIpv6ConfigurationMode2(session *Session, self PIFRef) (retval Ipv6ConfigurationMode, err error) {
	method := "PIF.get_ipv6_configuration_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumIpv6ConfigurationMode(method+" -> ", result)
	return
}

// GetTunnelTransportPIFOf: Get the tunnel_transport_PIF_of field of the given PIF.
// Version: rio
func (pif) GetTunnelTransportPIFOf(session *Session, self PIFRef) (retval []TunnelRef, err error) {
	method := "PIF.get_tunnel_transport_PIF_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRefSet(method+" -> ", result)
	return
}

// GetTunnelTransportPIFOf2: Get the tunnel_transport_PIF_of field of the given PIF.
// Version: rio
func (pif) GetTunnelTransportPIFOf2(session *Session, self PIFRef) (retval []TunnelRef, err error) {
	method := "PIF.get_tunnel_transport_PIF_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRefSet(method+" -> ", result)
	return
}

// GetTunnelAccessPIFOf: Get the tunnel_access_PIF_of field of the given PIF.
// Version: rio
func (pif) GetTunnelAccessPIFOf(session *Session, self PIFRef) (retval []TunnelRef, err error) {
	method := "PIF.get_tunnel_access_PIF_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRefSet(method+" -> ", result)
	return
}

// GetTunnelAccessPIFOf2: Get the tunnel_access_PIF_of field of the given PIF.
// Version: rio
func (pif) GetTunnelAccessPIFOf2(session *Session, self PIFRef) (retval []TunnelRef, err error) {
	method := "PIF.get_tunnel_access_PIF_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeTunnelRefSet(method+" -> ", result)
	return
}

// GetDisallowUnplug: Get the disallow_unplug field of the given PIF.
// Version: rio
func (pif) GetDisallowUnplug(session *Session, self PIFRef) (retval bool, err error) {
	method := "PIF.get_disallow_unplug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDisallowUnplug2: Get the disallow_unplug field of the given PIF.
// Version: rio
func (pif) GetDisallowUnplug2(session *Session, self PIFRef) (retval bool, err error) {
	method := "PIF.get_disallow_unplug"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given PIF.
// Version: rio
func (pif) GetOtherConfig(session *Session, self PIFRef) (retval map[string]string, err error) {
	method := "PIF.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given PIF.
// Version: rio
func (pif) GetOtherConfig2(session *Session, self PIFRef) (retval map[string]string, err error) {
	method := "PIF.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetManagement: Get the management field of the given PIF.
// Version: rio
func (pif) GetManagement(session *Session, self PIFRef) (retval bool, err error) {
	method := "PIF.get_management"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetManagement2: Get the management field of the given PIF.
// Version: rio
func (pif) GetManagement2(session *Session, self PIFRef) (retval bool, err error) {
	method := "PIF.get_management"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVLANSlaveOf: Get the VLAN_slave_of field of the given PIF.
// Version: rio
func (pif) GetVLANSlaveOf(session *Session, self PIFRef) (retval []VLANRef, err error) {
	method := "PIF.get_VLAN_slave_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVLANRefSet(method+" -> ", result)
	return
}

// GetVLANSlaveOf2: Get the VLAN_slave_of field of the given PIF.
// Version: rio
func (pif) GetVLANSlaveOf2(session *Session, self PIFRef) (retval []VLANRef, err error) {
	method := "PIF.get_VLAN_slave_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVLANRefSet(method+" -> ", result)
	return
}

// GetVLANMasterOf: Get the VLAN_master_of field of the given PIF.
// Version: rio
func (pif) GetVLANMasterOf(session *Session, self PIFRef) (retval VLANRef, err error) {
	method := "PIF.get_VLAN_master_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVLANRef(method+" -> ", result)
	return
}

// GetVLANMasterOf2: Get the VLAN_master_of field of the given PIF.
// Version: rio
func (pif) GetVLANMasterOf2(session *Session, self PIFRef) (retval VLANRef, err error) {
	method := "PIF.get_VLAN_master_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVLANRef(method+" -> ", result)
	return
}

// GetBondMasterOf: Get the bond_master_of field of the given PIF.
// Version: rio
func (pif) GetBondMasterOf(session *Session, self PIFRef) (retval []BondRef, err error) {
	method := "PIF.get_bond_master_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRefSet(method+" -> ", result)
	return
}

// GetBondMasterOf2: Get the bond_master_of field of the given PIF.
// Version: rio
func (pif) GetBondMasterOf2(session *Session, self PIFRef) (retval []BondRef, err error) {
	method := "PIF.get_bond_master_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRefSet(method+" -> ", result)
	return
}

// GetBondSlaveOf: Get the bond_slave_of field of the given PIF.
// Version: rio
func (pif) GetBondSlaveOf(session *Session, self PIFRef) (retval BondRef, err error) {
	method := "PIF.get_bond_slave_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRef(method+" -> ", result)
	return
}

// GetBondSlaveOf2: Get the bond_slave_of field of the given PIF.
// Version: rio
func (pif) GetBondSlaveOf2(session *Session, self PIFRef) (retval BondRef, err error) {
	method := "PIF.get_bond_slave_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeBondRef(method+" -> ", result)
	return
}

// GetDNS: Get the DNS field of the given PIF.
// Version: rio
func (pif) GetDNS(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_DNS"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDNS2: Get the DNS field of the given PIF.
// Version: rio
func (pif) GetDNS2(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_DNS"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetGateway: Get the gateway field of the given PIF.
// Version: rio
func (pif) GetGateway(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_gateway"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetGateway2: Get the gateway field of the given PIF.
// Version: rio
func (pif) GetGateway2(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_gateway"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNetmask: Get the netmask field of the given PIF.
// Version: rio
func (pif) GetNetmask(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_netmask"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNetmask2: Get the netmask field of the given PIF.
// Version: rio
func (pif) GetNetmask2(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_netmask"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIP: Get the IP field of the given PIF.
// Version: rio
func (pif) GetIP(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_IP"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIP2: Get the IP field of the given PIF.
// Version: rio
func (pif) GetIP2(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_IP"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIPConfigurationMode: Get the ip_configuration_mode field of the given PIF.
// Version: rio
func (pif) GetIPConfigurationMode(session *Session, self PIFRef) (retval IPConfigurationMode, err error) {
	method := "PIF.get_ip_configuration_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumIPConfigurationMode(method+" -> ", result)
	return
}

// GetIPConfigurationMode2: Get the ip_configuration_mode field of the given PIF.
// Version: rio
func (pif) GetIPConfigurationMode2(session *Session, self PIFRef) (retval IPConfigurationMode, err error) {
	method := "PIF.get_ip_configuration_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumIPConfigurationMode(method+" -> ", result)
	return
}

// GetCurrentlyAttached: Get the currently_attached field of the given PIF.
// Version: rio
func (pif) GetCurrentlyAttached(session *Session, self PIFRef) (retval bool, err error) {
	method := "PIF.get_currently_attached"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCurrentlyAttached2: Get the currently_attached field of the given PIF.
// Version: rio
func (pif) GetCurrentlyAttached2(session *Session, self PIFRef) (retval bool, err error) {
	method := "PIF.get_currently_attached"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPhysical: Get the physical field of the given PIF.
// Version: rio
func (pif) GetPhysical(session *Session, self PIFRef) (retval bool, err error) {
	method := "PIF.get_physical"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPhysical2: Get the physical field of the given PIF.
// Version: rio
func (pif) GetPhysical2(session *Session, self PIFRef) (retval bool, err error) {
	method := "PIF.get_physical"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMetrics: Get the metrics field of the given PIF.
// Version: rio
func (pif) GetMetrics(session *Session, self PIFRef) (retval PIFMetricsRef, err error) {
	method := "PIF.get_metrics"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFMetricsRef(method+" -> ", result)
	return
}

// GetMetrics2: Get the metrics field of the given PIF.
// Version: rio
func (pif) GetMetrics2(session *Session, self PIFRef) (retval PIFMetricsRef, err error) {
	method := "PIF.get_metrics"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFMetricsRef(method+" -> ", result)
	return
}

// GetVLAN: Get the VLAN field of the given PIF.
// Version: rio
func (pif) GetVLAN(session *Session, self PIFRef) (retval int, err error) {
	method := "PIF.get_VLAN"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVLAN2: Get the VLAN field of the given PIF.
// Version: rio
func (pif) GetVLAN2(session *Session, self PIFRef) (retval int, err error) {
	method := "PIF.get_VLAN"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMTU: Get the MTU field of the given PIF.
// Version: rio
func (pif) GetMTU(session *Session, self PIFRef) (retval int, err error) {
	method := "PIF.get_MTU"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMTU2: Get the MTU field of the given PIF.
// Version: rio
func (pif) GetMTU2(session *Session, self PIFRef) (retval int, err error) {
	method := "PIF.get_MTU"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMAC: Get the MAC field of the given PIF.
// Version: rio
func (pif) GetMAC(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_MAC"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMAC2: Get the MAC field of the given PIF.
// Version: rio
func (pif) GetMAC2(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_MAC"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHost: Get the host field of the given PIF.
// Version: rio
func (pif) GetHost(session *Session, self PIFRef) (retval HostRef, err error) {
	method := "PIF.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHost2: Get the host field of the given PIF.
// Version: rio
func (pif) GetHost2(session *Session, self PIFRef) (retval HostRef, err error) {
	method := "PIF.get_host"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNetwork: Get the network field of the given PIF.
// Version: rio
func (pif) GetNetwork(session *Session, self PIFRef) (retval NetworkRef, err error) {
	method := "PIF.get_network"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNetwork2: Get the network field of the given PIF.
// Version: rio
func (pif) GetNetwork2(session *Session, self PIFRef) (retval NetworkRef, err error) {
	method := "PIF.get_network"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDevice: Get the device field of the given PIF.
// Version: rio
func (pif) GetDevice(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDevice2: Get the device field of the given PIF.
// Version: rio
func (pif) GetDevice2(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given PIF.
// Version: rio
func (pif) GetUUID(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given PIF.
// Version: rio
func (pif) GetUUID2(session *Session, self PIFRef) (retval string, err error) {
	method := "PIF.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the PIF instance with the specified UUID.
// Version: rio
func (pif) GetByUUID(session *Session, uuid string) (retval PIFRef, err error) {
	method := "PIF.get_by_uuid"
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
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the PIF instance with the specified UUID.
// Version: rio
func (pif) GetByUUID2(session *Session, uuid string) (retval PIFRef, err error) {
	method := "PIF.get_by_uuid"
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
	retval, err = deserializePIFRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given PIF.
// Version: rio
func (pif) GetRecord(session *Session, self PIFRef) (retval PIFRecord, err error) {
	method := "PIF.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given PIF.
// Version: rio
func (pif) GetRecord2(session *Session, self PIFRef) (retval PIFRecord, err error) {
	method := "PIF.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializePIFRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePIFRecord(method+" -> ", result)
	return
}
