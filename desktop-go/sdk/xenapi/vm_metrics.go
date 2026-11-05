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

type VMMetricsRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// Guest&apos;s actual memory (bytes)
	MemoryActual int `json:"memory_actual,omitempty"`
	// Current number of VCPUs
	VCPUsNumber int `json:"VCPUs_number,omitempty"`
	// Utilisation for all of guest&apos;s current VCPUs
	VCPUsUtilisation map[int]float64 `json:"VCPUs_utilisation,omitempty"`
	// VCPU to PCPU map
	VCPUsCPU map[int]int `json:"VCPUs_CPU,omitempty"`
	// The live equivalent to VM.VCPUs_params
	VCPUsParams map[string]string `json:"VCPUs_params,omitempty"`
	// CPU flags (blocked,online,running)
	VCPUsFlags map[int][]string `json:"VCPUs_flags,omitempty"`
	// The state of the guest, eg blocked, dying etc
	State []string `json:"state,omitempty"`
	// Time at which this VM was last booted
	StartTime time.Time `json:"start_time,omitempty"`
	// Time at which the VM was installed
	InstallTime time.Time `json:"install_time,omitempty"`
	// Time at which this information was last updated
	LastUpdated time.Time `json:"last_updated,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// hardware virtual machine
	Hvm bool `json:"hvm,omitempty"`
	// VM supports nested virtualisation
	NestedVirt bool `json:"nested_virt,omitempty"`
	// VM is immobile and can&apos;t migrate between hosts
	Nomigrate bool `json:"nomigrate,omitempty"`
	// The current domain type of the VM (for running,suspended, or paused VMs). The last-known domain type for halted VMs.
	CurrentDomainType DomainType `json:"current_domain_type,omitempty"`
	// If the VM is optimised for NUMA
	NumaOptimised bool `json:"numa_optimised,omitempty"`
	// number of NUMA nodes of the host the VM is using
	NumaNodes int `json:"numa_nodes,omitempty"`
	// mapping a NUMA node (int) to an amount of memory (bytes) in that node.
	NumaNodeMemory map[int]int `json:"numa_node_memory,omitempty"`
}

type VMMetricsRef string

// The metrics associated with a VM
type vmMetrics struct{}

var VMMetrics vmMetrics

// GetAllRecords: Return a map of VM_metrics references to VM_metrics records for all VM_metrics instances known to the system.
// Version: rio
func (vmMetrics) GetAllRecords(session *Session) (retval map[VMMetricsRef]VMMetricsRecord, err error) {
	method := "VM_metrics.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMMetricsRefToVMMetricsRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of VM_metrics references to VM_metrics records for all VM_metrics instances known to the system.
// Version: rio
func (vmMetrics) GetAllRecords1(session *Session) (retval map[VMMetricsRef]VMMetricsRecord, err error) {
	method := "VM_metrics.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMMetricsRefToVMMetricsRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the VM_metrics instances known to the system.
// Version: rio
func (vmMetrics) GetAll(session *Session) (retval []VMMetricsRef, err error) {
	method := "VM_metrics.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMMetricsRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the VM_metrics instances known to the system.
// Version: rio
func (vmMetrics) GetAll1(session *Session) (retval []VMMetricsRef, err error) {
	method := "VM_metrics.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMMetricsRefSet(method+" -> ", result)
	return
}

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given VM_metrics.  If the key is not in that Map, then do nothing.
// Version: orlando
func (vmMetrics) RemoveFromOtherConfig(session *Session, self VMMetricsRef, key string) (err error) {
	method := "VM_metrics.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given VM_metrics.  If the key is not in that Map, then do nothing.
// Version: orlando
func (vmMetrics) RemoveFromOtherConfig3(session *Session, self VMMetricsRef, key string) (err error) {
	method := "VM_metrics.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig2: Remove the given key and its corresponding value from the other_config field of the given VM_metrics.  If the key is not in that Map, then do nothing.
// Version: rio
func (vmMetrics) RemoveFromOtherConfig2(session *Session, self VMMetricsRef) (err error) {
	method := "VM_metrics.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given VM_metrics.
// Version: orlando
func (vmMetrics) AddToOtherConfig(session *Session, self VMMetricsRef, key string, value string) (err error) {
	method := "VM_metrics.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given VM_metrics.
// Version: orlando
func (vmMetrics) AddToOtherConfig4(session *Session, self VMMetricsRef, key string, value string) (err error) {
	method := "VM_metrics.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig2: Add the given key-value pair to the other_config field of the given VM_metrics.
// Version: rio
func (vmMetrics) AddToOtherConfig2(session *Session, self VMMetricsRef) (err error) {
	method := "VM_metrics.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetOtherConfig: Set the other_config field of the given VM_metrics.
// Version: orlando
func (vmMetrics) SetOtherConfig(session *Session, self VMMetricsRef, value map[string]string) (err error) {
	method := "VM_metrics.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given VM_metrics.
// Version: orlando
func (vmMetrics) SetOtherConfig3(session *Session, self VMMetricsRef, value map[string]string) (err error) {
	method := "VM_metrics.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig2: Set the other_config field of the given VM_metrics.
// Version: rio
func (vmMetrics) SetOtherConfig2(session *Session, self VMMetricsRef) (err error) {
	method := "VM_metrics.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// GetNumaNodeMemory: Get the numa_node_memory field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetNumaNodeMemory(session *Session, self VMMetricsRef) (retval map[int]int, err error) {
	method := "VM_metrics.get_numa_node_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeIntToIntMap(method+" -> ", result)
	return
}

// GetNumaNodeMemory2: Get the numa_node_memory field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetNumaNodeMemory2(session *Session, self VMMetricsRef) (retval map[int]int, err error) {
	method := "VM_metrics.get_numa_node_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeIntToIntMap(method+" -> ", result)
	return
}

// GetNumaNodes: Get the numa_nodes field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetNumaNodes(session *Session, self VMMetricsRef) (retval int, err error) {
	method := "VM_metrics.get_numa_nodes"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNumaNodes2: Get the numa_nodes field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetNumaNodes2(session *Session, self VMMetricsRef) (retval int, err error) {
	method := "VM_metrics.get_numa_nodes"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNumaOptimised: Get the numa_optimised field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetNumaOptimised(session *Session, self VMMetricsRef) (retval bool, err error) {
	method := "VM_metrics.get_numa_optimised"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNumaOptimised2: Get the numa_optimised field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetNumaOptimised2(session *Session, self VMMetricsRef) (retval bool, err error) {
	method := "VM_metrics.get_numa_optimised"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCurrentDomainType: Get the current_domain_type field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetCurrentDomainType(session *Session, self VMMetricsRef) (retval DomainType, err error) {
	method := "VM_metrics.get_current_domain_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumDomainType(method+" -> ", result)
	return
}

// GetCurrentDomainType2: Get the current_domain_type field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetCurrentDomainType2(session *Session, self VMMetricsRef) (retval DomainType, err error) {
	method := "VM_metrics.get_current_domain_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumDomainType(method+" -> ", result)
	return
}

// GetNomigrate: Get the nomigrate field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetNomigrate(session *Session, self VMMetricsRef) (retval bool, err error) {
	method := "VM_metrics.get_nomigrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNomigrate2: Get the nomigrate field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetNomigrate2(session *Session, self VMMetricsRef) (retval bool, err error) {
	method := "VM_metrics.get_nomigrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNestedVirt: Get the nested_virt field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetNestedVirt(session *Session, self VMMetricsRef) (retval bool, err error) {
	method := "VM_metrics.get_nested_virt"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNestedVirt2: Get the nested_virt field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetNestedVirt2(session *Session, self VMMetricsRef) (retval bool, err error) {
	method := "VM_metrics.get_nested_virt"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHvm: Get the hvm field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetHvm(session *Session, self VMMetricsRef) (retval bool, err error) {
	method := "VM_metrics.get_hvm"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHvm2: Get the hvm field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetHvm2(session *Session, self VMMetricsRef) (retval bool, err error) {
	method := "VM_metrics.get_hvm"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetOtherConfig(session *Session, self VMMetricsRef) (retval map[string]string, err error) {
	method := "VM_metrics.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetOtherConfig2(session *Session, self VMMetricsRef) (retval map[string]string, err error) {
	method := "VM_metrics.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastUpdated: Get the last_updated field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetLastUpdated(session *Session, self VMMetricsRef) (retval time.Time, err error) {
	method := "VM_metrics.get_last_updated"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastUpdated2: Get the last_updated field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetLastUpdated2(session *Session, self VMMetricsRef) (retval time.Time, err error) {
	method := "VM_metrics.get_last_updated"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetInstallTime: Get the install_time field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetInstallTime(session *Session, self VMMetricsRef) (retval time.Time, err error) {
	method := "VM_metrics.get_install_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetInstallTime2: Get the install_time field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetInstallTime2(session *Session, self VMMetricsRef) (retval time.Time, err error) {
	method := "VM_metrics.get_install_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetStartTime: Get the start_time field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetStartTime(session *Session, self VMMetricsRef) (retval time.Time, err error) {
	method := "VM_metrics.get_start_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetStartTime2: Get the start_time field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetStartTime2(session *Session, self VMMetricsRef) (retval time.Time, err error) {
	method := "VM_metrics.get_start_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetState: Get the state field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetState(session *Session, self VMMetricsRef) (retval []string, err error) {
	method := "VM_metrics.get_state"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetState2: Get the state field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetState2(session *Session, self VMMetricsRef) (retval []string, err error) {
	method := "VM_metrics.get_state"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVCPUsFlags: Get the VCPUs/flags field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetVCPUsFlags(session *Session, self VMMetricsRef) (retval map[int][]string, err error) {
	method := "VM_metrics.get_VCPUs_flags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeIntToStringSetMap(method+" -> ", result)
	return
}

// GetVCPUsFlags2: Get the VCPUs/flags field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetVCPUsFlags2(session *Session, self VMMetricsRef) (retval map[int][]string, err error) {
	method := "VM_metrics.get_VCPUs_flags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeIntToStringSetMap(method+" -> ", result)
	return
}

// GetVCPUsParams: Get the VCPUs/params field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetVCPUsParams(session *Session, self VMMetricsRef) (retval map[string]string, err error) {
	method := "VM_metrics.get_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVCPUsParams2: Get the VCPUs/params field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetVCPUsParams2(session *Session, self VMMetricsRef) (retval map[string]string, err error) {
	method := "VM_metrics.get_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVCPUsCPU: Get the VCPUs/CPU field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetVCPUsCPU(session *Session, self VMMetricsRef) (retval map[int]int, err error) {
	method := "VM_metrics.get_VCPUs_CPU"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeIntToIntMap(method+" -> ", result)
	return
}

// GetVCPUsCPU2: Get the VCPUs/CPU field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetVCPUsCPU2(session *Session, self VMMetricsRef) (retval map[int]int, err error) {
	method := "VM_metrics.get_VCPUs_CPU"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeIntToIntMap(method+" -> ", result)
	return
}

// GetVCPUsUtilisation: Get the VCPUs/utilisation field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetVCPUsUtilisation(session *Session, self VMMetricsRef) (retval map[int]float64, err error) {
	method := "VM_metrics.get_VCPUs_utilisation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeIntToFloatMap(method+" -> ", result)
	return
}

// GetVCPUsUtilisation2: Get the VCPUs/utilisation field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetVCPUsUtilisation2(session *Session, self VMMetricsRef) (retval map[int]float64, err error) {
	method := "VM_metrics.get_VCPUs_utilisation"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeIntToFloatMap(method+" -> ", result)
	return
}

// GetVCPUsNumber: Get the VCPUs/number field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetVCPUsNumber(session *Session, self VMMetricsRef) (retval int, err error) {
	method := "VM_metrics.get_VCPUs_number"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVCPUsNumber2: Get the VCPUs/number field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetVCPUsNumber2(session *Session, self VMMetricsRef) (retval int, err error) {
	method := "VM_metrics.get_VCPUs_number"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryActual: Get the memory/actual field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetMemoryActual(session *Session, self VMMetricsRef) (retval int, err error) {
	method := "VM_metrics.get_memory_actual"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryActual2: Get the memory/actual field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetMemoryActual2(session *Session, self VMMetricsRef) (retval int, err error) {
	method := "VM_metrics.get_memory_actual"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetUUID(session *Session, self VMMetricsRef) (retval string, err error) {
	method := "VM_metrics.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given VM_metrics.
// Version: rio
func (vmMetrics) GetUUID2(session *Session, self VMMetricsRef) (retval string, err error) {
	method := "VM_metrics.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByUUID: Get a reference to the VM_metrics instance with the specified UUID.
// Version: rio
func (vmMetrics) GetByUUID(session *Session, uuid string) (retval VMMetricsRef, err error) {
	method := "VM_metrics.get_by_uuid"
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
	retval, err = deserializeVMMetricsRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the VM_metrics instance with the specified UUID.
// Version: rio
func (vmMetrics) GetByUUID2(session *Session, uuid string) (retval VMMetricsRef, err error) {
	method := "VM_metrics.get_by_uuid"
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
	retval, err = deserializeVMMetricsRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given VM_metrics.
// Version: rio
func (vmMetrics) GetRecord(session *Session, self VMMetricsRef) (retval VMMetricsRecord, err error) {
	method := "VM_metrics.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMMetricsRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given VM_metrics.
// Version: rio
func (vmMetrics) GetRecord2(session *Session, self VMMetricsRef) (retval VMMetricsRecord, err error) {
	method := "VM_metrics.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMMetricsRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMMetricsRecord(method+" -> ", result)
	return
}
