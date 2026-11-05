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
	"math"
	"reflect"
	"strconv"
	"time"
)

//nolint:unparam
func serializeFloat(context string, value float64) (interface{}, error) {
	_ = context
	if math.IsInf(value, 0) {
		if math.IsInf(value, 1) {
			return "+Inf", nil
		}
		return "-Inf", nil
	} else if math.IsNaN(value) {
		return "NaN", nil
	}
	return value, nil
}

func serializeBlobRef(context string, ref BlobRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeBondRef(context string, ref BondRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeBool(context string, value bool) (bool, error) {
	_ = context
	return value, nil
}

func serializeCertificateRef(context string, ref CertificateRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeClusterHostRef(context string, ref ClusterHostRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeClusterRef(context string, ref ClusterRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeConsoleRecord(context string, record ConsoleRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["protocol"], err = serializeEnumConsoleProtocol(fmt.Sprintf("%s.%s", context, "protocol"), record.Protocol)
	if err != nil {
		return
	}
	rpcStruct["location"], err = serializeString(fmt.Sprintf("%s.%s", context, "location"), record.Location)
	if err != nil {
		return
	}
	rpcStruct["VM"], err = serializeVMRef(fmt.Sprintf("%s.%s", context, "VM"), record.VM)
	if err != nil {
		return
	}
	rpcStruct["other_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), record.OtherConfig)
	if err != nil {
		return
	}
	return
}

func serializeConsoleRef(context string, ref ConsoleRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeConsoleRefSet(context string, slice []ConsoleRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeConsoleRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeCrashdumpRef(context string, ref CrashdumpRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeCrashdumpRefSet(context string, slice []CrashdumpRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeCrashdumpRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeDRTaskRef(context string, ref DRTaskRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeDriverVariantRef(context string, ref DriverVariantRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeEnumAllocationAlgorithm(context string, value AllocationAlgorithm) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumBondMode(context string, value BondMode) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumCertificatePurpose(context string, value CertificatePurpose) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumCertificatePurposeSet(context string, slice []CertificatePurpose) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeEnumCertificatePurpose(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeEnumCls(context string, value Cls) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumConsoleProtocol(context string, value ConsoleProtocol) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumDomainType(context string, value DomainType) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumHostDisplay(context string, value HostDisplay) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumHostNtpMode(context string, value HostNtpMode) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumHostNumaAffinityPolicy(context string, value HostNumaAffinityPolicy) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumHostSchedGran(context string, value HostSchedGran) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumIPConfigurationMode(context string, value IPConfigurationMode) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumIpv6ConfigurationMode(context string, value Ipv6ConfigurationMode) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumNetworkDefaultLockingMode(context string, value NetworkDefaultLockingMode) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumNetworkOperations(context string, value NetworkOperations) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumNetworkOperationsSet(context string, slice []NetworkOperations) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeEnumNetworkOperations(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeEnumNetworkPurpose(context string, value NetworkPurpose) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumNetworkPurposeSet(context string, slice []NetworkPurpose) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeEnumNetworkPurpose(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeEnumOnBoot(context string, value OnBoot) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumOnCrashBehaviour(context string, value OnCrashBehaviour) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumOnNormalExit(context string, value OnNormalExit) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumOnSoftrebootBehavior(context string, value OnSoftrebootBehavior) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumPlacementPolicy(context string, value PlacementPolicy) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumPrimaryAddressType(context string, value PrimaryAddressType) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumSdnControllerProtocol(context string, value SdnControllerProtocol) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumTaskStatusType(context string, value TaskStatusType) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumTunnelProtocol(context string, value TunnelProtocol) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumUpdateGuidances(context string, value UpdateGuidances) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumUpdateGuidancesSet(context string, slice []UpdateGuidances) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeEnumUpdateGuidances(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeEnumUpdateSyncFrequency(context string, value UpdateSyncFrequency) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVMApplianceOperation(context string, value VMApplianceOperation) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVMApplianceOperationSet(context string, slice []VMApplianceOperation) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeEnumVMApplianceOperation(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeEnumVMOperations(context string, value VMOperations) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVMOperationsSet(context string, slice []VMOperations) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeEnumVMOperations(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeEnumVMOperationsToStringMap(context string, goMap map[VMOperations]string) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeEnumVMOperations(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeString(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeEnumVMPowerState(context string, value VMPowerState) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVMSecurebootCertificatesState(context string, value VMSecurebootCertificatesState) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVMUefiMode(context string, value VMUefiMode) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVbdMode(context string, value VbdMode) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVbdOperations(context string, value VbdOperations) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVbdOperationsSet(context string, slice []VbdOperations) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeEnumVbdOperations(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeEnumVbdType(context string, value VbdType) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVdiOperations(context string, value VdiOperations) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVdiOperationsSet(context string, slice []VdiOperations) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeEnumVdiOperations(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeEnumVdiType(context string, value VdiType) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVifIpv4ConfigurationMode(context string, value VifIpv4ConfigurationMode) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVifIpv6ConfigurationMode(context string, value VifIpv6ConfigurationMode) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVifLockingMode(context string, value VifLockingMode) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVifOperations(context string, value VifOperations) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVifOperationsSet(context string, slice []VifOperations) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeEnumVifOperations(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeEnumVmppArchiveFrequency(context string, value VmppArchiveFrequency) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVmppArchiveTargetType(context string, value VmppArchiveTargetType) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVmppBackupFrequency(context string, value VmppBackupFrequency) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVmppBackupType(context string, value VmppBackupType) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVmssFrequency(context string, value VmssFrequency) (string, error) {
	_ = context
	return string(value), nil
}

func serializeEnumVmssType(context string, value VmssType) (string, error) {
	_ = context
	return string(value), nil
}

func serializeFeatureRef(context string, ref FeatureRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeGPUGroupRef(context string, ref GPUGroupRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeHostCPURef(context string, ref HostCPURef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeHostCrashdumpRef(context string, ref HostCrashdumpRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeHostDriverRef(context string, ref HostDriverRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeHostMetricsRef(context string, ref HostMetricsRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeHostPatchRef(context string, ref HostPatchRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeHostRef(context string, ref HostRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeHostRefSet(context string, slice []HostRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeHostRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeInt(context string, value int) (int, error) {
	_ = context
	return value, nil
}

func serializeLVHDRef(context string, ref LVHDRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeMessageRef(context string, ref MessageRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeMessageRefSet(context string, slice []MessageRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeMessageRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeNetworkRecord(context string, record NetworkRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["name_label"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_label"), record.NameLabel)
	if err != nil {
		return
	}
	rpcStruct["name_description"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_description"), record.NameDescription)
	if err != nil {
		return
	}
	rpcStruct["allowed_operations"], err = serializeEnumNetworkOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), record.AllowedOperations)
	if err != nil {
		return
	}
	rpcStruct["current_operations"], err = serializeStringToEnumNetworkOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), record.CurrentOperations)
	if err != nil {
		return
	}
	rpcStruct["VIFs"], err = serializeVIFRefSet(fmt.Sprintf("%s.%s", context, "VIFs"), record.VIFs)
	if err != nil {
		return
	}
	rpcStruct["PIFs"], err = serializePIFRefSet(fmt.Sprintf("%s.%s", context, "PIFs"), record.PIFs)
	if err != nil {
		return
	}
	rpcStruct["MTU"], err = serializeInt(fmt.Sprintf("%s.%s", context, "MTU"), record.MTU)
	if err != nil {
		return
	}
	rpcStruct["other_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), record.OtherConfig)
	if err != nil {
		return
	}
	rpcStruct["bridge"], err = serializeString(fmt.Sprintf("%s.%s", context, "bridge"), record.Bridge)
	if err != nil {
		return
	}
	rpcStruct["managed"], err = serializeBool(fmt.Sprintf("%s.%s", context, "managed"), record.Managed)
	if err != nil {
		return
	}
	rpcStruct["blobs"], err = serializeStringToBlobRefMap(fmt.Sprintf("%s.%s", context, "blobs"), record.Blobs)
	if err != nil {
		return
	}
	rpcStruct["tags"], err = serializeStringSet(fmt.Sprintf("%s.%s", context, "tags"), record.Tags)
	if err != nil {
		return
	}
	rpcStruct["default_locking_mode"], err = serializeEnumNetworkDefaultLockingMode(fmt.Sprintf("%s.%s", context, "default_locking_mode"), record.DefaultLockingMode)
	if err != nil {
		return
	}
	rpcStruct["assigned_ips"], err = serializeVIFRefToStringMap(fmt.Sprintf("%s.%s", context, "assigned_ips"), record.AssignedIps)
	if err != nil {
		return
	}
	rpcStruct["purpose"], err = serializeEnumNetworkPurposeSet(fmt.Sprintf("%s.%s", context, "purpose"), record.Purpose)
	if err != nil {
		return
	}
	return
}

func serializeNetworkRef(context string, ref NetworkRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeNetworkSriovRef(context string, ref NetworkSriovRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeObserverRecord(context string, record ObserverRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["name_label"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_label"), record.NameLabel)
	if err != nil {
		return
	}
	rpcStruct["name_description"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_description"), record.NameDescription)
	if err != nil {
		return
	}
	rpcStruct["hosts"], err = serializeHostRefSet(fmt.Sprintf("%s.%s", context, "hosts"), record.Hosts)
	if err != nil {
		return
	}
	rpcStruct["attributes"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "attributes"), record.Attributes)
	if err != nil {
		return
	}
	rpcStruct["endpoints"], err = serializeStringSet(fmt.Sprintf("%s.%s", context, "endpoints"), record.Endpoints)
	if err != nil {
		return
	}
	rpcStruct["components"], err = serializeStringSet(fmt.Sprintf("%s.%s", context, "components"), record.Components)
	if err != nil {
		return
	}
	rpcStruct["enabled"], err = serializeBool(fmt.Sprintf("%s.%s", context, "enabled"), record.Enabled)
	if err != nil {
		return
	}
	return
}

func serializeObserverRef(context string, ref ObserverRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePBDRecord(context string, record PBDRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["host"], err = serializeHostRef(fmt.Sprintf("%s.%s", context, "host"), record.Host)
	if err != nil {
		return
	}
	rpcStruct["SR"], err = serializeSRRef(fmt.Sprintf("%s.%s", context, "SR"), record.SR)
	if err != nil {
		return
	}
	rpcStruct["device_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "device_config"), record.DeviceConfig)
	if err != nil {
		return
	}
	rpcStruct["currently_attached"], err = serializeBool(fmt.Sprintf("%s.%s", context, "currently_attached"), record.CurrentlyAttached)
	if err != nil {
		return
	}
	rpcStruct["other_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), record.OtherConfig)
	if err != nil {
		return
	}
	return
}

func serializePBDRef(context string, ref PBDRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePCIRef(context string, ref PCIRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePCIRefSet(context string, slice []PCIRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializePCIRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializePGPURef(context string, ref PGPURef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePIFMetricsRef(context string, ref PIFMetricsRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePIFRef(context string, ref PIFRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePIFRefSet(context string, slice []PIFRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializePIFRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializePUSBRef(context string, ref PUSBRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePVSCacheStorageRecord(context string, record PVSCacheStorageRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["host"], err = serializeHostRef(fmt.Sprintf("%s.%s", context, "host"), record.Host)
	if err != nil {
		return
	}
	rpcStruct["SR"], err = serializeSRRef(fmt.Sprintf("%s.%s", context, "SR"), record.SR)
	if err != nil {
		return
	}
	rpcStruct["site"], err = serializePVSSiteRef(fmt.Sprintf("%s.%s", context, "site"), record.Site)
	if err != nil {
		return
	}
	rpcStruct["size"], err = serializeInt(fmt.Sprintf("%s.%s", context, "size"), record.Size)
	if err != nil {
		return
	}
	rpcStruct["VDI"], err = serializeVDIRef(fmt.Sprintf("%s.%s", context, "VDI"), record.VDI)
	if err != nil {
		return
	}
	return
}

func serializePVSCacheStorageRef(context string, ref PVSCacheStorageRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePVSProxyRef(context string, ref PVSProxyRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePVSServerRef(context string, ref PVSServerRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePVSSiteRef(context string, ref PVSSiteRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePoolPatchRef(context string, ref PoolPatchRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePoolRef(context string, ref PoolRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializePoolUpdateRef(context string, ref PoolUpdateRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeRepositoryRef(context string, ref RepositoryRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeRepositoryRefSet(context string, slice []RepositoryRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeRepositoryRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeRoleRef(context string, ref RoleRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeRoleRefSet(context string, slice []RoleRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeRoleRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeSDNControllerRef(context string, ref SDNControllerRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeSMRef(context string, ref SMRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeSRRef(context string, ref SRRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeSRRefSet(context string, slice []SRRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeSRRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeSecretRecord(context string, record SecretRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["value"], err = serializeString(fmt.Sprintf("%s.%s", context, "value"), record.Value)
	if err != nil {
		return
	}
	rpcStruct["other_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), record.OtherConfig)
	if err != nil {
		return
	}
	return
}

func serializeSecretRef(context string, ref SecretRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeSessionRef(context string, ref SessionRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeString(context string, value string) (string, error) {
	_ = context
	return value, nil
}

func serializeStringSet(context string, slice []string) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeString(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeStringToBlobRefMap(context string, goMap map[string]BlobRef) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeString(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeBlobRef(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeStringToEnumNetworkOperationsMap(context string, goMap map[string]NetworkOperations) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeString(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeEnumNetworkOperations(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeStringToEnumVMApplianceOperationMap(context string, goMap map[string]VMApplianceOperation) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeString(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeEnumVMApplianceOperation(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeStringToEnumVMOperationsMap(context string, goMap map[string]VMOperations) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeString(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeEnumVMOperations(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeStringToEnumVbdOperationsMap(context string, goMap map[string]VbdOperations) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeString(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeEnumVbdOperations(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeStringToEnumVdiOperationsMap(context string, goMap map[string]VdiOperations) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeString(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeEnumVdiOperations(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeStringToEnumVifOperationsMap(context string, goMap map[string]VifOperations) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeString(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeEnumVifOperations(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeStringToStringMap(context string, goMap map[string]string) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeString(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeString(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeSubjectRecord(context string, record SubjectRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["subject_identifier"], err = serializeString(fmt.Sprintf("%s.%s", context, "subject_identifier"), record.SubjectIdentifier)
	if err != nil {
		return
	}
	rpcStruct["other_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), record.OtherConfig)
	if err != nil {
		return
	}
	rpcStruct["roles"], err = serializeRoleRefSet(fmt.Sprintf("%s.%s", context, "roles"), record.Roles)
	if err != nil {
		return
	}
	return
}

func serializeSubjectRef(context string, ref SubjectRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeTaskRef(context string, ref TaskRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeTunnelRef(context string, ref TunnelRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeUSBGroupRef(context string, ref USBGroupRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeUserRecord(context string, record UserRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["short_name"], err = serializeString(fmt.Sprintf("%s.%s", context, "short_name"), record.ShortName)
	if err != nil {
		return
	}
	rpcStruct["fullname"], err = serializeString(fmt.Sprintf("%s.%s", context, "fullname"), record.Fullname)
	if err != nil {
		return
	}
	rpcStruct["other_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), record.OtherConfig)
	if err != nil {
		return
	}
	return
}

func serializeUserRef(context string, ref UserRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVBDMetricsRef(context string, ref VBDMetricsRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVBDRecord(context string, record VBDRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["allowed_operations"], err = serializeEnumVbdOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), record.AllowedOperations)
	if err != nil {
		return
	}
	rpcStruct["current_operations"], err = serializeStringToEnumVbdOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), record.CurrentOperations)
	if err != nil {
		return
	}
	rpcStruct["VM"], err = serializeVMRef(fmt.Sprintf("%s.%s", context, "VM"), record.VM)
	if err != nil {
		return
	}
	rpcStruct["VDI"], err = serializeVDIRef(fmt.Sprintf("%s.%s", context, "VDI"), record.VDI)
	if err != nil {
		return
	}
	rpcStruct["device"], err = serializeString(fmt.Sprintf("%s.%s", context, "device"), record.Device)
	if err != nil {
		return
	}
	rpcStruct["userdevice"], err = serializeString(fmt.Sprintf("%s.%s", context, "userdevice"), record.Userdevice)
	if err != nil {
		return
	}
	rpcStruct["bootable"], err = serializeBool(fmt.Sprintf("%s.%s", context, "bootable"), record.Bootable)
	if err != nil {
		return
	}
	rpcStruct["mode"], err = serializeEnumVbdMode(fmt.Sprintf("%s.%s", context, "mode"), record.Mode)
	if err != nil {
		return
	}
	rpcStruct["type"], err = serializeEnumVbdType(fmt.Sprintf("%s.%s", context, "type"), record.Type)
	if err != nil {
		return
	}
	rpcStruct["unpluggable"], err = serializeBool(fmt.Sprintf("%s.%s", context, "unpluggable"), record.Unpluggable)
	if err != nil {
		return
	}
	rpcStruct["storage_lock"], err = serializeBool(fmt.Sprintf("%s.%s", context, "storage_lock"), record.StorageLock)
	if err != nil {
		return
	}
	rpcStruct["empty"], err = serializeBool(fmt.Sprintf("%s.%s", context, "empty"), record.Empty)
	if err != nil {
		return
	}
	rpcStruct["other_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), record.OtherConfig)
	if err != nil {
		return
	}
	rpcStruct["currently_attached"], err = serializeBool(fmt.Sprintf("%s.%s", context, "currently_attached"), record.CurrentlyAttached)
	if err != nil {
		return
	}
	rpcStruct["status_code"], err = serializeInt(fmt.Sprintf("%s.%s", context, "status_code"), record.StatusCode)
	if err != nil {
		return
	}
	rpcStruct["status_detail"], err = serializeString(fmt.Sprintf("%s.%s", context, "status_detail"), record.StatusDetail)
	if err != nil {
		return
	}
	rpcStruct["runtime_properties"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "runtime_properties"), record.RuntimeProperties)
	if err != nil {
		return
	}
	rpcStruct["qos_algorithm_type"], err = serializeString(fmt.Sprintf("%s.%s", context, "qos_algorithm_type"), record.QosAlgorithmType)
	if err != nil {
		return
	}
	rpcStruct["qos_algorithm_params"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "qos_algorithm_params"), record.QosAlgorithmParams)
	if err != nil {
		return
	}
	rpcStruct["qos_supported_algorithms"], err = serializeStringSet(fmt.Sprintf("%s.%s", context, "qos_supported_algorithms"), record.QosSupportedAlgorithms)
	if err != nil {
		return
	}
	rpcStruct["metrics"], err = serializeVBDMetricsRef(fmt.Sprintf("%s.%s", context, "metrics"), record.Metrics)
	if err != nil {
		return
	}
	return
}

func serializeVBDRef(context string, ref VBDRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVBDRefSet(context string, slice []VBDRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeVBDRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeVDIRecord(context string, record VDIRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["name_label"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_label"), record.NameLabel)
	if err != nil {
		return
	}
	rpcStruct["name_description"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_description"), record.NameDescription)
	if err != nil {
		return
	}
	rpcStruct["allowed_operations"], err = serializeEnumVdiOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), record.AllowedOperations)
	if err != nil {
		return
	}
	rpcStruct["current_operations"], err = serializeStringToEnumVdiOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), record.CurrentOperations)
	if err != nil {
		return
	}
	rpcStruct["SR"], err = serializeSRRef(fmt.Sprintf("%s.%s", context, "SR"), record.SR)
	if err != nil {
		return
	}
	rpcStruct["VBDs"], err = serializeVBDRefSet(fmt.Sprintf("%s.%s", context, "VBDs"), record.VBDs)
	if err != nil {
		return
	}
	rpcStruct["crash_dumps"], err = serializeCrashdumpRefSet(fmt.Sprintf("%s.%s", context, "crash_dumps"), record.CrashDumps)
	if err != nil {
		return
	}
	rpcStruct["virtual_size"], err = serializeInt(fmt.Sprintf("%s.%s", context, "virtual_size"), record.VirtualSize)
	if err != nil {
		return
	}
	rpcStruct["physical_utilisation"], err = serializeInt(fmt.Sprintf("%s.%s", context, "physical_utilisation"), record.PhysicalUtilisation)
	if err != nil {
		return
	}
	rpcStruct["type"], err = serializeEnumVdiType(fmt.Sprintf("%s.%s", context, "type"), record.Type)
	if err != nil {
		return
	}
	rpcStruct["sharable"], err = serializeBool(fmt.Sprintf("%s.%s", context, "sharable"), record.Sharable)
	if err != nil {
		return
	}
	rpcStruct["read_only"], err = serializeBool(fmt.Sprintf("%s.%s", context, "read_only"), record.ReadOnly)
	if err != nil {
		return
	}
	rpcStruct["other_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), record.OtherConfig)
	if err != nil {
		return
	}
	rpcStruct["storage_lock"], err = serializeBool(fmt.Sprintf("%s.%s", context, "storage_lock"), record.StorageLock)
	if err != nil {
		return
	}
	rpcStruct["location"], err = serializeString(fmt.Sprintf("%s.%s", context, "location"), record.Location)
	if err != nil {
		return
	}
	rpcStruct["managed"], err = serializeBool(fmt.Sprintf("%s.%s", context, "managed"), record.Managed)
	if err != nil {
		return
	}
	rpcStruct["missing"], err = serializeBool(fmt.Sprintf("%s.%s", context, "missing"), record.Missing)
	if err != nil {
		return
	}
	rpcStruct["parent"], err = serializeVDIRef(fmt.Sprintf("%s.%s", context, "parent"), record.Parent)
	if err != nil {
		return
	}
	rpcStruct["xenstore_data"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "xenstore_data"), record.XenstoreData)
	if err != nil {
		return
	}
	rpcStruct["sm_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "sm_config"), record.SmConfig)
	if err != nil {
		return
	}
	rpcStruct["is_a_snapshot"], err = serializeBool(fmt.Sprintf("%s.%s", context, "is_a_snapshot"), record.IsASnapshot)
	if err != nil {
		return
	}
	rpcStruct["snapshot_of"], err = serializeVDIRef(fmt.Sprintf("%s.%s", context, "snapshot_of"), record.SnapshotOf)
	if err != nil {
		return
	}
	rpcStruct["snapshots"], err = serializeVDIRefSet(fmt.Sprintf("%s.%s", context, "snapshots"), record.Snapshots)
	if err != nil {
		return
	}
	rpcStruct["snapshot_time"], err = serializeTime(fmt.Sprintf("%s.%s", context, "snapshot_time"), record.SnapshotTime)
	if err != nil {
		return
	}
	rpcStruct["tags"], err = serializeStringSet(fmt.Sprintf("%s.%s", context, "tags"), record.Tags)
	if err != nil {
		return
	}
	rpcStruct["allow_caching"], err = serializeBool(fmt.Sprintf("%s.%s", context, "allow_caching"), record.AllowCaching)
	if err != nil {
		return
	}
	rpcStruct["on_boot"], err = serializeEnumOnBoot(fmt.Sprintf("%s.%s", context, "on_boot"), record.OnBoot)
	if err != nil {
		return
	}
	rpcStruct["metadata_of_pool"], err = serializePoolRef(fmt.Sprintf("%s.%s", context, "metadata_of_pool"), record.MetadataOfPool)
	if err != nil {
		return
	}
	rpcStruct["metadata_latest"], err = serializeBool(fmt.Sprintf("%s.%s", context, "metadata_latest"), record.MetadataLatest)
	if err != nil {
		return
	}
	rpcStruct["is_tools_iso"], err = serializeBool(fmt.Sprintf("%s.%s", context, "is_tools_iso"), record.IsToolsIso)
	if err != nil {
		return
	}
	rpcStruct["cbt_enabled"], err = serializeBool(fmt.Sprintf("%s.%s", context, "cbt_enabled"), record.CbtEnabled)
	if err != nil {
		return
	}
	return
}

func serializeVDIRef(context string, ref VDIRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVDIRefSet(context string, slice []VDIRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeVDIRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeVDIRefToSRRefMap(context string, goMap map[VDIRef]SRRef) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeVDIRef(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeSRRef(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeVDIRefToStringMap(context string, goMap map[VDIRef]string) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeVDIRef(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeString(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeVGPURef(context string, ref VGPURef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVGPURefSet(context string, slice []VGPURef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeVGPURef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeVGPURefToGPUGroupRefMap(context string, goMap map[VGPURef]GPUGroupRef) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeVGPURef(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeGPUGroupRef(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeVGPUTypeRef(context string, ref VGPUTypeRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVGPUTypeRefSet(context string, slice []VGPUTypeRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeVGPUTypeRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeVIFMetricsRef(context string, ref VIFMetricsRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVIFRecord(context string, record VIFRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["allowed_operations"], err = serializeEnumVifOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), record.AllowedOperations)
	if err != nil {
		return
	}
	rpcStruct["current_operations"], err = serializeStringToEnumVifOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), record.CurrentOperations)
	if err != nil {
		return
	}
	rpcStruct["device"], err = serializeString(fmt.Sprintf("%s.%s", context, "device"), record.Device)
	if err != nil {
		return
	}
	rpcStruct["network"], err = serializeNetworkRef(fmt.Sprintf("%s.%s", context, "network"), record.Network)
	if err != nil {
		return
	}
	rpcStruct["VM"], err = serializeVMRef(fmt.Sprintf("%s.%s", context, "VM"), record.VM)
	if err != nil {
		return
	}
	rpcStruct["MAC"], err = serializeString(fmt.Sprintf("%s.%s", context, "MAC"), record.MAC)
	if err != nil {
		return
	}
	rpcStruct["MTU"], err = serializeInt(fmt.Sprintf("%s.%s", context, "MTU"), record.MTU)
	if err != nil {
		return
	}
	rpcStruct["other_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), record.OtherConfig)
	if err != nil {
		return
	}
	rpcStruct["currently_attached"], err = serializeBool(fmt.Sprintf("%s.%s", context, "currently_attached"), record.CurrentlyAttached)
	if err != nil {
		return
	}
	rpcStruct["status_code"], err = serializeInt(fmt.Sprintf("%s.%s", context, "status_code"), record.StatusCode)
	if err != nil {
		return
	}
	rpcStruct["status_detail"], err = serializeString(fmt.Sprintf("%s.%s", context, "status_detail"), record.StatusDetail)
	if err != nil {
		return
	}
	rpcStruct["runtime_properties"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "runtime_properties"), record.RuntimeProperties)
	if err != nil {
		return
	}
	rpcStruct["qos_algorithm_type"], err = serializeString(fmt.Sprintf("%s.%s", context, "qos_algorithm_type"), record.QosAlgorithmType)
	if err != nil {
		return
	}
	rpcStruct["qos_algorithm_params"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "qos_algorithm_params"), record.QosAlgorithmParams)
	if err != nil {
		return
	}
	rpcStruct["qos_supported_algorithms"], err = serializeStringSet(fmt.Sprintf("%s.%s", context, "qos_supported_algorithms"), record.QosSupportedAlgorithms)
	if err != nil {
		return
	}
	rpcStruct["metrics"], err = serializeVIFMetricsRef(fmt.Sprintf("%s.%s", context, "metrics"), record.Metrics)
	if err != nil {
		return
	}
	rpcStruct["MAC_autogenerated"], err = serializeBool(fmt.Sprintf("%s.%s", context, "MAC_autogenerated"), record.MACAutogenerated)
	if err != nil {
		return
	}
	rpcStruct["locking_mode"], err = serializeEnumVifLockingMode(fmt.Sprintf("%s.%s", context, "locking_mode"), record.LockingMode)
	if err != nil {
		return
	}
	rpcStruct["ipv4_allowed"], err = serializeStringSet(fmt.Sprintf("%s.%s", context, "ipv4_allowed"), record.Ipv4Allowed)
	if err != nil {
		return
	}
	rpcStruct["ipv6_allowed"], err = serializeStringSet(fmt.Sprintf("%s.%s", context, "ipv6_allowed"), record.Ipv6Allowed)
	if err != nil {
		return
	}
	rpcStruct["ipv4_configuration_mode"], err = serializeEnumVifIpv4ConfigurationMode(fmt.Sprintf("%s.%s", context, "ipv4_configuration_mode"), record.Ipv4ConfigurationMode)
	if err != nil {
		return
	}
	rpcStruct["ipv4_addresses"], err = serializeStringSet(fmt.Sprintf("%s.%s", context, "ipv4_addresses"), record.Ipv4Addresses)
	if err != nil {
		return
	}
	rpcStruct["ipv4_gateway"], err = serializeString(fmt.Sprintf("%s.%s", context, "ipv4_gateway"), record.Ipv4Gateway)
	if err != nil {
		return
	}
	rpcStruct["ipv6_configuration_mode"], err = serializeEnumVifIpv6ConfigurationMode(fmt.Sprintf("%s.%s", context, "ipv6_configuration_mode"), record.Ipv6ConfigurationMode)
	if err != nil {
		return
	}
	rpcStruct["ipv6_addresses"], err = serializeStringSet(fmt.Sprintf("%s.%s", context, "ipv6_addresses"), record.Ipv6Addresses)
	if err != nil {
		return
	}
	rpcStruct["ipv6_gateway"], err = serializeString(fmt.Sprintf("%s.%s", context, "ipv6_gateway"), record.Ipv6Gateway)
	if err != nil {
		return
	}
	return
}

func serializeVIFRef(context string, ref VIFRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVIFRefSet(context string, slice []VIFRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeVIFRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeVIFRefToNetworkRefMap(context string, goMap map[VIFRef]NetworkRef) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeVIFRef(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeNetworkRef(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeVIFRefToStringMap(context string, goMap map[VIFRef]string) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeVIFRef(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeString(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeVLANRef(context string, ref VLANRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVMApplianceRecord(context string, record VMApplianceRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["name_label"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_label"), record.NameLabel)
	if err != nil {
		return
	}
	rpcStruct["name_description"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_description"), record.NameDescription)
	if err != nil {
		return
	}
	rpcStruct["allowed_operations"], err = serializeEnumVMApplianceOperationSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), record.AllowedOperations)
	if err != nil {
		return
	}
	rpcStruct["current_operations"], err = serializeStringToEnumVMApplianceOperationMap(fmt.Sprintf("%s.%s", context, "current_operations"), record.CurrentOperations)
	if err != nil {
		return
	}
	rpcStruct["VMs"], err = serializeVMRefSet(fmt.Sprintf("%s.%s", context, "VMs"), record.VMs)
	if err != nil {
		return
	}
	return
}

func serializeVMApplianceRef(context string, ref VMApplianceRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVMGroupRecord(context string, record VMGroupRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["name_label"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_label"), record.NameLabel)
	if err != nil {
		return
	}
	rpcStruct["name_description"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_description"), record.NameDescription)
	if err != nil {
		return
	}
	rpcStruct["placement"], err = serializeEnumPlacementPolicy(fmt.Sprintf("%s.%s", context, "placement"), record.Placement)
	if err != nil {
		return
	}
	rpcStruct["VMs"], err = serializeVMRefSet(fmt.Sprintf("%s.%s", context, "VMs"), record.VMs)
	if err != nil {
		return
	}
	return
}

func serializeVMGroupRef(context string, ref VMGroupRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVMGroupRefSet(context string, slice []VMGroupRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeVMGroupRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeVMGuestMetricsRef(context string, ref VMGuestMetricsRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVMMetricsRef(context string, ref VMMetricsRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVMPPRecord(context string, record VMPPRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["name_label"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_label"), record.NameLabel)
	if err != nil {
		return
	}
	rpcStruct["name_description"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_description"), record.NameDescription)
	if err != nil {
		return
	}
	rpcStruct["is_policy_enabled"], err = serializeBool(fmt.Sprintf("%s.%s", context, "is_policy_enabled"), record.IsPolicyEnabled)
	if err != nil {
		return
	}
	rpcStruct["backup_type"], err = serializeEnumVmppBackupType(fmt.Sprintf("%s.%s", context, "backup_type"), record.BackupType)
	if err != nil {
		return
	}
	rpcStruct["backup_retention_value"], err = serializeInt(fmt.Sprintf("%s.%s", context, "backup_retention_value"), record.BackupRetentionValue)
	if err != nil {
		return
	}
	rpcStruct["backup_frequency"], err = serializeEnumVmppBackupFrequency(fmt.Sprintf("%s.%s", context, "backup_frequency"), record.BackupFrequency)
	if err != nil {
		return
	}
	rpcStruct["backup_schedule"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "backup_schedule"), record.BackupSchedule)
	if err != nil {
		return
	}
	rpcStruct["is_backup_running"], err = serializeBool(fmt.Sprintf("%s.%s", context, "is_backup_running"), record.IsBackupRunning)
	if err != nil {
		return
	}
	rpcStruct["backup_last_run_time"], err = serializeTime(fmt.Sprintf("%s.%s", context, "backup_last_run_time"), record.BackupLastRunTime)
	if err != nil {
		return
	}
	rpcStruct["archive_target_type"], err = serializeEnumVmppArchiveTargetType(fmt.Sprintf("%s.%s", context, "archive_target_type"), record.ArchiveTargetType)
	if err != nil {
		return
	}
	rpcStruct["archive_target_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "archive_target_config"), record.ArchiveTargetConfig)
	if err != nil {
		return
	}
	rpcStruct["archive_frequency"], err = serializeEnumVmppArchiveFrequency(fmt.Sprintf("%s.%s", context, "archive_frequency"), record.ArchiveFrequency)
	if err != nil {
		return
	}
	rpcStruct["archive_schedule"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "archive_schedule"), record.ArchiveSchedule)
	if err != nil {
		return
	}
	rpcStruct["is_archive_running"], err = serializeBool(fmt.Sprintf("%s.%s", context, "is_archive_running"), record.IsArchiveRunning)
	if err != nil {
		return
	}
	rpcStruct["archive_last_run_time"], err = serializeTime(fmt.Sprintf("%s.%s", context, "archive_last_run_time"), record.ArchiveLastRunTime)
	if err != nil {
		return
	}
	rpcStruct["VMs"], err = serializeVMRefSet(fmt.Sprintf("%s.%s", context, "VMs"), record.VMs)
	if err != nil {
		return
	}
	rpcStruct["is_alarm_enabled"], err = serializeBool(fmt.Sprintf("%s.%s", context, "is_alarm_enabled"), record.IsAlarmEnabled)
	if err != nil {
		return
	}
	rpcStruct["alarm_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "alarm_config"), record.AlarmConfig)
	if err != nil {
		return
	}
	rpcStruct["recent_alerts"], err = serializeStringSet(fmt.Sprintf("%s.%s", context, "recent_alerts"), record.RecentAlerts)
	if err != nil {
		return
	}
	return
}

func serializeVMPPRef(context string, ref VMPPRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVMRecord(context string, record VMRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["allowed_operations"], err = serializeEnumVMOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), record.AllowedOperations)
	if err != nil {
		return
	}
	rpcStruct["current_operations"], err = serializeStringToEnumVMOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), record.CurrentOperations)
	if err != nil {
		return
	}
	rpcStruct["name_label"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_label"), record.NameLabel)
	if err != nil {
		return
	}
	rpcStruct["name_description"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_description"), record.NameDescription)
	if err != nil {
		return
	}
	rpcStruct["power_state"], err = serializeEnumVMPowerState(fmt.Sprintf("%s.%s", context, "power_state"), record.PowerState)
	if err != nil {
		return
	}
	rpcStruct["user_version"], err = serializeInt(fmt.Sprintf("%s.%s", context, "user_version"), record.UserVersion)
	if err != nil {
		return
	}
	rpcStruct["is_a_template"], err = serializeBool(fmt.Sprintf("%s.%s", context, "is_a_template"), record.IsATemplate)
	if err != nil {
		return
	}
	rpcStruct["is_default_template"], err = serializeBool(fmt.Sprintf("%s.%s", context, "is_default_template"), record.IsDefaultTemplate)
	if err != nil {
		return
	}
	rpcStruct["suspend_VDI"], err = serializeVDIRef(fmt.Sprintf("%s.%s", context, "suspend_VDI"), record.SuspendVDI)
	if err != nil {
		return
	}
	rpcStruct["resident_on"], err = serializeHostRef(fmt.Sprintf("%s.%s", context, "resident_on"), record.ResidentOn)
	if err != nil {
		return
	}
	rpcStruct["scheduled_to_be_resident_on"], err = serializeHostRef(fmt.Sprintf("%s.%s", context, "scheduled_to_be_resident_on"), record.ScheduledToBeResidentOn)
	if err != nil {
		return
	}
	rpcStruct["affinity"], err = serializeHostRef(fmt.Sprintf("%s.%s", context, "affinity"), record.Affinity)
	if err != nil {
		return
	}
	rpcStruct["memory_overhead"], err = serializeInt(fmt.Sprintf("%s.%s", context, "memory_overhead"), record.MemoryOverhead)
	if err != nil {
		return
	}
	rpcStruct["memory_target"], err = serializeInt(fmt.Sprintf("%s.%s", context, "memory_target"), record.MemoryTarget)
	if err != nil {
		return
	}
	rpcStruct["memory_static_max"], err = serializeInt(fmt.Sprintf("%s.%s", context, "memory_static_max"), record.MemoryStaticMax)
	if err != nil {
		return
	}
	rpcStruct["memory_dynamic_max"], err = serializeInt(fmt.Sprintf("%s.%s", context, "memory_dynamic_max"), record.MemoryDynamicMax)
	if err != nil {
		return
	}
	rpcStruct["memory_dynamic_min"], err = serializeInt(fmt.Sprintf("%s.%s", context, "memory_dynamic_min"), record.MemoryDynamicMin)
	if err != nil {
		return
	}
	rpcStruct["memory_static_min"], err = serializeInt(fmt.Sprintf("%s.%s", context, "memory_static_min"), record.MemoryStaticMin)
	if err != nil {
		return
	}
	rpcStruct["VCPUs_params"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "VCPUs_params"), record.VCPUsParams)
	if err != nil {
		return
	}
	rpcStruct["VCPUs_max"], err = serializeInt(fmt.Sprintf("%s.%s", context, "VCPUs_max"), record.VCPUsMax)
	if err != nil {
		return
	}
	rpcStruct["VCPUs_at_startup"], err = serializeInt(fmt.Sprintf("%s.%s", context, "VCPUs_at_startup"), record.VCPUsAtStartup)
	if err != nil {
		return
	}
	rpcStruct["actions_after_softreboot"], err = serializeEnumOnSoftrebootBehavior(fmt.Sprintf("%s.%s", context, "actions_after_softreboot"), record.ActionsAfterSoftreboot)
	if err != nil {
		return
	}
	rpcStruct["actions_after_shutdown"], err = serializeEnumOnNormalExit(fmt.Sprintf("%s.%s", context, "actions_after_shutdown"), record.ActionsAfterShutdown)
	if err != nil {
		return
	}
	rpcStruct["actions_after_reboot"], err = serializeEnumOnNormalExit(fmt.Sprintf("%s.%s", context, "actions_after_reboot"), record.ActionsAfterReboot)
	if err != nil {
		return
	}
	rpcStruct["actions_after_crash"], err = serializeEnumOnCrashBehaviour(fmt.Sprintf("%s.%s", context, "actions_after_crash"), record.ActionsAfterCrash)
	if err != nil {
		return
	}
	rpcStruct["consoles"], err = serializeConsoleRefSet(fmt.Sprintf("%s.%s", context, "consoles"), record.Consoles)
	if err != nil {
		return
	}
	rpcStruct["VIFs"], err = serializeVIFRefSet(fmt.Sprintf("%s.%s", context, "VIFs"), record.VIFs)
	if err != nil {
		return
	}
	rpcStruct["VBDs"], err = serializeVBDRefSet(fmt.Sprintf("%s.%s", context, "VBDs"), record.VBDs)
	if err != nil {
		return
	}
	rpcStruct["VUSBs"], err = serializeVUSBRefSet(fmt.Sprintf("%s.%s", context, "VUSBs"), record.VUSBs)
	if err != nil {
		return
	}
	rpcStruct["crash_dumps"], err = serializeCrashdumpRefSet(fmt.Sprintf("%s.%s", context, "crash_dumps"), record.CrashDumps)
	if err != nil {
		return
	}
	rpcStruct["VTPMs"], err = serializeVTPMRefSet(fmt.Sprintf("%s.%s", context, "VTPMs"), record.VTPMs)
	if err != nil {
		return
	}
	rpcStruct["PV_bootloader"], err = serializeString(fmt.Sprintf("%s.%s", context, "PV_bootloader"), record.PVBootloader)
	if err != nil {
		return
	}
	rpcStruct["PV_kernel"], err = serializeString(fmt.Sprintf("%s.%s", context, "PV_kernel"), record.PVKernel)
	if err != nil {
		return
	}
	rpcStruct["PV_ramdisk"], err = serializeString(fmt.Sprintf("%s.%s", context, "PV_ramdisk"), record.PVRamdisk)
	if err != nil {
		return
	}
	rpcStruct["PV_args"], err = serializeString(fmt.Sprintf("%s.%s", context, "PV_args"), record.PVArgs)
	if err != nil {
		return
	}
	rpcStruct["PV_bootloader_args"], err = serializeString(fmt.Sprintf("%s.%s", context, "PV_bootloader_args"), record.PVBootloaderArgs)
	if err != nil {
		return
	}
	rpcStruct["PV_legacy_args"], err = serializeString(fmt.Sprintf("%s.%s", context, "PV_legacy_args"), record.PVLegacyArgs)
	if err != nil {
		return
	}
	rpcStruct["HVM_boot_policy"], err = serializeString(fmt.Sprintf("%s.%s", context, "HVM_boot_policy"), record.HVMBootPolicy)
	if err != nil {
		return
	}
	rpcStruct["HVM_boot_params"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "HVM_boot_params"), record.HVMBootParams)
	if err != nil {
		return
	}
	rpcStruct["HVM_shadow_multiplier"], err = serializeFloat(fmt.Sprintf("%s.%s", context, "HVM_shadow_multiplier"), record.HVMShadowMultiplier)
	if err != nil {
		return
	}
	rpcStruct["platform"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "platform"), record.Platform)
	if err != nil {
		return
	}
	rpcStruct["PCI_bus"], err = serializeString(fmt.Sprintf("%s.%s", context, "PCI_bus"), record.PCIBus)
	if err != nil {
		return
	}
	rpcStruct["other_config"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), record.OtherConfig)
	if err != nil {
		return
	}
	rpcStruct["domid"], err = serializeInt(fmt.Sprintf("%s.%s", context, "domid"), record.Domid)
	if err != nil {
		return
	}
	rpcStruct["domarch"], err = serializeString(fmt.Sprintf("%s.%s", context, "domarch"), record.Domarch)
	if err != nil {
		return
	}
	rpcStruct["last_boot_CPU_flags"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "last_boot_CPU_flags"), record.LastBootCPUFlags)
	if err != nil {
		return
	}
	rpcStruct["is_control_domain"], err = serializeBool(fmt.Sprintf("%s.%s", context, "is_control_domain"), record.IsControlDomain)
	if err != nil {
		return
	}
	rpcStruct["metrics"], err = serializeVMMetricsRef(fmt.Sprintf("%s.%s", context, "metrics"), record.Metrics)
	if err != nil {
		return
	}
	rpcStruct["guest_metrics"], err = serializeVMGuestMetricsRef(fmt.Sprintf("%s.%s", context, "guest_metrics"), record.GuestMetrics)
	if err != nil {
		return
	}
	rpcStruct["last_booted_record"], err = serializeString(fmt.Sprintf("%s.%s", context, "last_booted_record"), record.LastBootedRecord)
	if err != nil {
		return
	}
	rpcStruct["recommendations"], err = serializeString(fmt.Sprintf("%s.%s", context, "recommendations"), record.Recommendations)
	if err != nil {
		return
	}
	rpcStruct["xenstore_data"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "xenstore_data"), record.XenstoreData)
	if err != nil {
		return
	}
	rpcStruct["ha_always_run"], err = serializeBool(fmt.Sprintf("%s.%s", context, "ha_always_run"), record.HaAlwaysRun)
	if err != nil {
		return
	}
	rpcStruct["ha_restart_priority"], err = serializeString(fmt.Sprintf("%s.%s", context, "ha_restart_priority"), record.HaRestartPriority)
	if err != nil {
		return
	}
	rpcStruct["is_a_snapshot"], err = serializeBool(fmt.Sprintf("%s.%s", context, "is_a_snapshot"), record.IsASnapshot)
	if err != nil {
		return
	}
	rpcStruct["snapshot_of"], err = serializeVMRef(fmt.Sprintf("%s.%s", context, "snapshot_of"), record.SnapshotOf)
	if err != nil {
		return
	}
	rpcStruct["snapshots"], err = serializeVMRefSet(fmt.Sprintf("%s.%s", context, "snapshots"), record.Snapshots)
	if err != nil {
		return
	}
	rpcStruct["snapshot_time"], err = serializeTime(fmt.Sprintf("%s.%s", context, "snapshot_time"), record.SnapshotTime)
	if err != nil {
		return
	}
	rpcStruct["transportable_snapshot_id"], err = serializeString(fmt.Sprintf("%s.%s", context, "transportable_snapshot_id"), record.TransportableSnapshotID)
	if err != nil {
		return
	}
	rpcStruct["blobs"], err = serializeStringToBlobRefMap(fmt.Sprintf("%s.%s", context, "blobs"), record.Blobs)
	if err != nil {
		return
	}
	rpcStruct["tags"], err = serializeStringSet(fmt.Sprintf("%s.%s", context, "tags"), record.Tags)
	if err != nil {
		return
	}
	rpcStruct["blocked_operations"], err = serializeEnumVMOperationsToStringMap(fmt.Sprintf("%s.%s", context, "blocked_operations"), record.BlockedOperations)
	if err != nil {
		return
	}
	rpcStruct["snapshot_info"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "snapshot_info"), record.SnapshotInfo)
	if err != nil {
		return
	}
	rpcStruct["snapshot_metadata"], err = serializeString(fmt.Sprintf("%s.%s", context, "snapshot_metadata"), record.SnapshotMetadata)
	if err != nil {
		return
	}
	rpcStruct["parent"], err = serializeVMRef(fmt.Sprintf("%s.%s", context, "parent"), record.Parent)
	if err != nil {
		return
	}
	rpcStruct["children"], err = serializeVMRefSet(fmt.Sprintf("%s.%s", context, "children"), record.Children)
	if err != nil {
		return
	}
	rpcStruct["bios_strings"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "bios_strings"), record.BiosStrings)
	if err != nil {
		return
	}
	rpcStruct["protection_policy"], err = serializeVMPPRef(fmt.Sprintf("%s.%s", context, "protection_policy"), record.ProtectionPolicy)
	if err != nil {
		return
	}
	rpcStruct["is_snapshot_from_vmpp"], err = serializeBool(fmt.Sprintf("%s.%s", context, "is_snapshot_from_vmpp"), record.IsSnapshotFromVmpp)
	if err != nil {
		return
	}
	rpcStruct["snapshot_schedule"], err = serializeVMSSRef(fmt.Sprintf("%s.%s", context, "snapshot_schedule"), record.SnapshotSchedule)
	if err != nil {
		return
	}
	rpcStruct["is_vmss_snapshot"], err = serializeBool(fmt.Sprintf("%s.%s", context, "is_vmss_snapshot"), record.IsVmssSnapshot)
	if err != nil {
		return
	}
	rpcStruct["appliance"], err = serializeVMApplianceRef(fmt.Sprintf("%s.%s", context, "appliance"), record.Appliance)
	if err != nil {
		return
	}
	rpcStruct["start_delay"], err = serializeInt(fmt.Sprintf("%s.%s", context, "start_delay"), record.StartDelay)
	if err != nil {
		return
	}
	rpcStruct["shutdown_delay"], err = serializeInt(fmt.Sprintf("%s.%s", context, "shutdown_delay"), record.ShutdownDelay)
	if err != nil {
		return
	}
	rpcStruct["order"], err = serializeInt(fmt.Sprintf("%s.%s", context, "order"), record.Order)
	if err != nil {
		return
	}
	rpcStruct["VGPUs"], err = serializeVGPURefSet(fmt.Sprintf("%s.%s", context, "VGPUs"), record.VGPUs)
	if err != nil {
		return
	}
	rpcStruct["attached_PCIs"], err = serializePCIRefSet(fmt.Sprintf("%s.%s", context, "attached_PCIs"), record.AttachedPCIs)
	if err != nil {
		return
	}
	rpcStruct["suspend_SR"], err = serializeSRRef(fmt.Sprintf("%s.%s", context, "suspend_SR"), record.SuspendSR)
	if err != nil {
		return
	}
	rpcStruct["version"], err = serializeInt(fmt.Sprintf("%s.%s", context, "version"), record.Version)
	if err != nil {
		return
	}
	rpcStruct["generation_id"], err = serializeString(fmt.Sprintf("%s.%s", context, "generation_id"), record.GenerationID)
	if err != nil {
		return
	}
	rpcStruct["hardware_platform_version"], err = serializeInt(fmt.Sprintf("%s.%s", context, "hardware_platform_version"), record.HardwarePlatformVersion)
	if err != nil {
		return
	}
	rpcStruct["has_vendor_device"], err = serializeBool(fmt.Sprintf("%s.%s", context, "has_vendor_device"), record.HasVendorDevice)
	if err != nil {
		return
	}
	rpcStruct["requires_reboot"], err = serializeBool(fmt.Sprintf("%s.%s", context, "requires_reboot"), record.RequiresReboot)
	if err != nil {
		return
	}
	rpcStruct["reference_label"], err = serializeString(fmt.Sprintf("%s.%s", context, "reference_label"), record.ReferenceLabel)
	if err != nil {
		return
	}
	rpcStruct["domain_type"], err = serializeEnumDomainType(fmt.Sprintf("%s.%s", context, "domain_type"), record.DomainType)
	if err != nil {
		return
	}
	rpcStruct["NVRAM"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "NVRAM"), record.NVRAM)
	if err != nil {
		return
	}
	rpcStruct["pending_guidances"], err = serializeEnumUpdateGuidancesSet(fmt.Sprintf("%s.%s", context, "pending_guidances"), record.PendingGuidances)
	if err != nil {
		return
	}
	rpcStruct["pending_guidances_recommended"], err = serializeEnumUpdateGuidancesSet(fmt.Sprintf("%s.%s", context, "pending_guidances_recommended"), record.PendingGuidancesRecommended)
	if err != nil {
		return
	}
	rpcStruct["pending_guidances_full"], err = serializeEnumUpdateGuidancesSet(fmt.Sprintf("%s.%s", context, "pending_guidances_full"), record.PendingGuidancesFull)
	if err != nil {
		return
	}
	rpcStruct["groups"], err = serializeVMGroupRefSet(fmt.Sprintf("%s.%s", context, "groups"), record.Groups)
	if err != nil {
		return
	}
	rpcStruct["secureboot_certificates_state"], err = serializeEnumVMSecurebootCertificatesState(fmt.Sprintf("%s.%s", context, "secureboot_certificates_state"), record.SecurebootCertificatesState)
	if err != nil {
		return
	}
	return
}

func serializeVMRef(context string, ref VMRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVMRefSet(context string, slice []VMRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeVMRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeVMRefToStringMap(context string, goMap map[VMRef]string) (xenMap map[string]interface{}, err error) {
	xenMap = make(map[string]interface{})
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := serializeVMRef(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := serializeString(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}

func serializeVMSSRecord(context string, record VMSSRecord) (rpcStruct map[string]interface{}, err error) {
	rpcStruct = map[string]interface{}{}
	rpcStruct["uuid"], err = serializeString(fmt.Sprintf("%s.%s", context, "uuid"), record.UUID)
	if err != nil {
		return
	}
	rpcStruct["name_label"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_label"), record.NameLabel)
	if err != nil {
		return
	}
	rpcStruct["name_description"], err = serializeString(fmt.Sprintf("%s.%s", context, "name_description"), record.NameDescription)
	if err != nil {
		return
	}
	rpcStruct["enabled"], err = serializeBool(fmt.Sprintf("%s.%s", context, "enabled"), record.Enabled)
	if err != nil {
		return
	}
	rpcStruct["type"], err = serializeEnumVmssType(fmt.Sprintf("%s.%s", context, "type"), record.Type)
	if err != nil {
		return
	}
	rpcStruct["retained_snapshots"], err = serializeInt(fmt.Sprintf("%s.%s", context, "retained_snapshots"), record.RetainedSnapshots)
	if err != nil {
		return
	}
	rpcStruct["frequency"], err = serializeEnumVmssFrequency(fmt.Sprintf("%s.%s", context, "frequency"), record.Frequency)
	if err != nil {
		return
	}
	rpcStruct["schedule"], err = serializeStringToStringMap(fmt.Sprintf("%s.%s", context, "schedule"), record.Schedule)
	if err != nil {
		return
	}
	rpcStruct["last_run_time"], err = serializeTime(fmt.Sprintf("%s.%s", context, "last_run_time"), record.LastRunTime)
	if err != nil {
		return
	}
	rpcStruct["VMs"], err = serializeVMRefSet(fmt.Sprintf("%s.%s", context, "VMs"), record.VMs)
	if err != nil {
		return
	}
	return
}

func serializeVMSSRef(context string, ref VMSSRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVTPMRef(context string, ref VTPMRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVTPMRefSet(context string, slice []VTPMRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeVTPMRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

func serializeVUSBRef(context string, ref VUSBRef) (string, error) {
	_ = context
	return string(ref), nil
}

func serializeVUSBRefSet(context string, slice []VUSBRef) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := serializeVUSBRef(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}

var timeFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",

	// no dashes, no colons
	"20060102T15:04:05Z",
	"20060102T15:04:05",
	"20060102T150405.999999999Z0700",
	"20060102T150405",
	"20060102T150405Z07",
	"20060102T150405Z07:00",

	// no dashes, with colons
	"20060102T15:04:05Z07",
	"20060102T15:04:05Z0700",
	"20060102T15:04:05Z07:00",
	"20060102T15:04:05.999999999Z07",
	"20060102T15:04:05.999999999Z07:00",
	"20060102T15:04:05.999999999Z07",

	// dashes and colon patterns not covered by `time.RFC3339`
	"2006-01-02T15:04:05Z07",
	"2006-01-02T15:04:05Z0700",
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05.999999999Z07",
	"2006-01-02T15:04:05.999999999Z07:00",
	"2006-01-02T15:04:05.999999999Z07",
}

//nolint:unparam
func serializeTime(context string, value time.Time) (string, error) {
	_ = context
	return value.Format(time.RFC3339), nil
}

func deserializeBlobRecord(context string, input interface{}) (record BlobRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	sizeValue, ok := rpcStruct["size"]
	if ok && sizeValue != nil {
		record.Size, err = deserializeInt(fmt.Sprintf("%s.%s", context, "size"), sizeValue)
		if err != nil {
			return
		}
	}
	publicValue, ok := rpcStruct["public"]
	if ok && publicValue != nil {
		record.Public, err = deserializeBool(fmt.Sprintf("%s.%s", context, "public"), publicValue)
		if err != nil {
			return
		}
	}
	lastUpdatedValue, ok := rpcStruct["last_updated"]
	if ok && lastUpdatedValue != nil {
		record.LastUpdated, err = deserializeTime(fmt.Sprintf("%s.%s", context, "last_updated"), lastUpdatedValue)
		if err != nil {
			return
		}
	}
	mimeTypeValue, ok := rpcStruct["mime_type"]
	if ok && mimeTypeValue != nil {
		record.MimeType, err = deserializeString(fmt.Sprintf("%s.%s", context, "mime_type"), mimeTypeValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeBlobRef(context string, input interface{}) (BlobRef, error) {
	var ref BlobRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return BlobRef(value), nil
}

func deserializeBlobRefSet(context string, input interface{}) (slice []BlobRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]BlobRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeBlobRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeBlobRefToBlobRecordMap(context string, input interface{}) (goMap map[BlobRef]BlobRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[BlobRef]BlobRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeBlobRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeBlobRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeBondRecord(context string, input interface{}) (record BondRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	masterValue, ok := rpcStruct["master"]
	if ok && masterValue != nil {
		record.Master, err = deserializePIFRef(fmt.Sprintf("%s.%s", context, "master"), masterValue)
		if err != nil {
			return
		}
	}
	slavesValue, ok := rpcStruct["slaves"]
	if ok && slavesValue != nil {
		record.Slaves, err = deserializePIFRefSet(fmt.Sprintf("%s.%s", context, "slaves"), slavesValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	primarySlaveValue, ok := rpcStruct["primary_slave"]
	if ok && primarySlaveValue != nil {
		record.PrimarySlave, err = deserializePIFRef(fmt.Sprintf("%s.%s", context, "primary_slave"), primarySlaveValue)
		if err != nil {
			return
		}
	}
	modeValue, ok := rpcStruct["mode"]
	if ok && modeValue != nil {
		record.Mode, err = deserializeEnumBondMode(fmt.Sprintf("%s.%s", context, "mode"), modeValue)
		if err != nil {
			return
		}
	}
	propertiesValue, ok := rpcStruct["properties"]
	if ok && propertiesValue != nil {
		record.Properties, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "properties"), propertiesValue)
		if err != nil {
			return
		}
	}
	linksUpValue, ok := rpcStruct["links_up"]
	if ok && linksUpValue != nil {
		record.LinksUp, err = deserializeInt(fmt.Sprintf("%s.%s", context, "links_up"), linksUpValue)
		if err != nil {
			return
		}
	}
	autoUpdateMacValue, ok := rpcStruct["auto_update_mac"]
	if ok && autoUpdateMacValue != nil {
		record.AutoUpdateMac, err = deserializeBool(fmt.Sprintf("%s.%s", context, "auto_update_mac"), autoUpdateMacValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeBondRef(context string, input interface{}) (BondRef, error) {
	var ref BondRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return BondRef(value), nil
}

func deserializeBondRefSet(context string, input interface{}) (slice []BondRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]BondRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeBondRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeBondRefToBondRecordMap(context string, input interface{}) (goMap map[BondRef]BondRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[BondRef]BondRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeBondRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeBondRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeBool(context string, input interface{}) (value bool, err error) {
	if input == nil {
		return
	}
	value, ok := input.(bool)
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "bool", context, reflect.TypeOf(input), input)
	}
	return
}

func deserializeCertificateRecord(context string, input interface{}) (record CertificateRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameValue, ok := rpcStruct["name"]
	if ok && nameValue != nil {
		record.Name, err = deserializeString(fmt.Sprintf("%s.%s", context, "name"), nameValue)
		if err != nil {
			return
		}
	}
	typeValue, ok := rpcStruct["type"]
	if ok && typeValue != nil {
		record.Type, err = deserializeEnumCertificateType(fmt.Sprintf("%s.%s", context, "type"), typeValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	notBeforeValue, ok := rpcStruct["not_before"]
	if ok && notBeforeValue != nil {
		record.NotBefore, err = deserializeTime(fmt.Sprintf("%s.%s", context, "not_before"), notBeforeValue)
		if err != nil {
			return
		}
	}
	notAfterValue, ok := rpcStruct["not_after"]
	if ok && notAfterValue != nil {
		record.NotAfter, err = deserializeTime(fmt.Sprintf("%s.%s", context, "not_after"), notAfterValue)
		if err != nil {
			return
		}
	}
	fingerprintValue, ok := rpcStruct["fingerprint"]
	if ok && fingerprintValue != nil {
		record.Fingerprint, err = deserializeString(fmt.Sprintf("%s.%s", context, "fingerprint"), fingerprintValue)
		if err != nil {
			return
		}
	}
	fingerprintSha256Value, ok := rpcStruct["fingerprint_sha256"]
	if ok && fingerprintSha256Value != nil {
		record.FingerprintSha256, err = deserializeString(fmt.Sprintf("%s.%s", context, "fingerprint_sha256"), fingerprintSha256Value)
		if err != nil {
			return
		}
	}
	fingerprintSha1Value, ok := rpcStruct["fingerprint_sha1"]
	if ok && fingerprintSha1Value != nil {
		record.FingerprintSha1, err = deserializeString(fmt.Sprintf("%s.%s", context, "fingerprint_sha1"), fingerprintSha1Value)
		if err != nil {
			return
		}
	}
	purposeValue, ok := rpcStruct["purpose"]
	if ok && purposeValue != nil {
		record.Purpose, err = deserializeEnumCertificatePurposeSet(fmt.Sprintf("%s.%s", context, "purpose"), purposeValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeCertificateRef(context string, input interface{}) (CertificateRef, error) {
	var ref CertificateRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return CertificateRef(value), nil
}

func deserializeCertificateRefSet(context string, input interface{}) (slice []CertificateRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]CertificateRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeCertificateRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeCertificateRefToCertificateRecordMap(context string, input interface{}) (goMap map[CertificateRef]CertificateRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[CertificateRef]CertificateRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeCertificateRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeCertificateRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeClusterHostRecord(context string, input interface{}) (record ClusterHostRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	clusterValue, ok := rpcStruct["cluster"]
	if ok && clusterValue != nil {
		record.Cluster, err = deserializeClusterRef(fmt.Sprintf("%s.%s", context, "cluster"), clusterValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	enabledValue, ok := rpcStruct["enabled"]
	if ok && enabledValue != nil {
		record.Enabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "enabled"), enabledValue)
		if err != nil {
			return
		}
	}
	pIFValue, ok := rpcStruct["PIF"]
	if ok && pIFValue != nil {
		record.PIF, err = deserializePIFRef(fmt.Sprintf("%s.%s", context, "PIF"), pIFValue)
		if err != nil {
			return
		}
	}
	joinedValue, ok := rpcStruct["joined"]
	if ok && joinedValue != nil {
		record.Joined, err = deserializeBool(fmt.Sprintf("%s.%s", context, "joined"), joinedValue)
		if err != nil {
			return
		}
	}
	liveValue, ok := rpcStruct["live"]
	if ok && liveValue != nil {
		record.Live, err = deserializeBool(fmt.Sprintf("%s.%s", context, "live"), liveValue)
		if err != nil {
			return
		}
	}
	lastUpdateLiveValue, ok := rpcStruct["last_update_live"]
	if ok && lastUpdateLiveValue != nil {
		record.LastUpdateLive, err = deserializeTime(fmt.Sprintf("%s.%s", context, "last_update_live"), lastUpdateLiveValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumClusterHostOperationSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumClusterHostOperationMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeClusterHostRef(context string, input interface{}) (ClusterHostRef, error) {
	var ref ClusterHostRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return ClusterHostRef(value), nil
}

func deserializeClusterHostRefSet(context string, input interface{}) (slice []ClusterHostRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]ClusterHostRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeClusterHostRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeClusterHostRefToClusterHostRecordMap(context string, input interface{}) (goMap map[ClusterHostRef]ClusterHostRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[ClusterHostRef]ClusterHostRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeClusterHostRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeClusterHostRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeClusterRecord(context string, input interface{}) (record ClusterRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	clusterHostsValue, ok := rpcStruct["cluster_hosts"]
	if ok && clusterHostsValue != nil {
		record.ClusterHosts, err = deserializeClusterHostRefSet(fmt.Sprintf("%s.%s", context, "cluster_hosts"), clusterHostsValue)
		if err != nil {
			return
		}
	}
	pendingForgetValue, ok := rpcStruct["pending_forget"]
	if ok && pendingForgetValue != nil {
		record.PendingForget, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "pending_forget"), pendingForgetValue)
		if err != nil {
			return
		}
	}
	clusterTokenValue, ok := rpcStruct["cluster_token"]
	if ok && clusterTokenValue != nil {
		record.ClusterToken, err = deserializeString(fmt.Sprintf("%s.%s", context, "cluster_token"), clusterTokenValue)
		if err != nil {
			return
		}
	}
	clusterStackValue, ok := rpcStruct["cluster_stack"]
	if ok && clusterStackValue != nil {
		record.ClusterStack, err = deserializeString(fmt.Sprintf("%s.%s", context, "cluster_stack"), clusterStackValue)
		if err != nil {
			return
		}
	}
	clusterStackVersionValue, ok := rpcStruct["cluster_stack_version"]
	if ok && clusterStackVersionValue != nil {
		record.ClusterStackVersion, err = deserializeInt(fmt.Sprintf("%s.%s", context, "cluster_stack_version"), clusterStackVersionValue)
		if err != nil {
			return
		}
	}
	isQuorateValue, ok := rpcStruct["is_quorate"]
	if ok && isQuorateValue != nil {
		record.IsQuorate, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_quorate"), isQuorateValue)
		if err != nil {
			return
		}
	}
	quorumValue, ok := rpcStruct["quorum"]
	if ok && quorumValue != nil {
		record.Quorum, err = deserializeInt(fmt.Sprintf("%s.%s", context, "quorum"), quorumValue)
		if err != nil {
			return
		}
	}
	liveHostsValue, ok := rpcStruct["live_hosts"]
	if ok && liveHostsValue != nil {
		record.LiveHosts, err = deserializeInt(fmt.Sprintf("%s.%s", context, "live_hosts"), liveHostsValue)
		if err != nil {
			return
		}
	}
	expectedHostsValue, ok := rpcStruct["expected_hosts"]
	if ok && expectedHostsValue != nil {
		record.ExpectedHosts, err = deserializeInt(fmt.Sprintf("%s.%s", context, "expected_hosts"), expectedHostsValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumClusterOperationSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumClusterOperationMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	poolAutoJoinValue, ok := rpcStruct["pool_auto_join"]
	if ok && poolAutoJoinValue != nil {
		record.PoolAutoJoin, err = deserializeBool(fmt.Sprintf("%s.%s", context, "pool_auto_join"), poolAutoJoinValue)
		if err != nil {
			return
		}
	}
	tokenTimeoutValue, ok := rpcStruct["token_timeout"]
	if ok && tokenTimeoutValue != nil {
		record.TokenTimeout, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "token_timeout"), tokenTimeoutValue)
		if err != nil {
			return
		}
	}
	tokenTimeoutCoefficientValue, ok := rpcStruct["token_timeout_coefficient"]
	if ok && tokenTimeoutCoefficientValue != nil {
		record.TokenTimeoutCoefficient, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "token_timeout_coefficient"), tokenTimeoutCoefficientValue)
		if err != nil {
			return
		}
	}
	clusterConfigValue, ok := rpcStruct["cluster_config"]
	if ok && clusterConfigValue != nil {
		record.ClusterConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "cluster_config"), clusterConfigValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeClusterRef(context string, input interface{}) (ClusterRef, error) {
	var ref ClusterRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return ClusterRef(value), nil
}

func deserializeClusterRefSet(context string, input interface{}) (slice []ClusterRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]ClusterRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeClusterRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeClusterRefToClusterRecordMap(context string, input interface{}) (goMap map[ClusterRef]ClusterRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[ClusterRef]ClusterRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeClusterRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeClusterRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeConsoleRecord(context string, input interface{}) (record ConsoleRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	protocolValue, ok := rpcStruct["protocol"]
	if ok && protocolValue != nil {
		record.Protocol, err = deserializeEnumConsoleProtocol(fmt.Sprintf("%s.%s", context, "protocol"), protocolValue)
		if err != nil {
			return
		}
	}
	locationValue, ok := rpcStruct["location"]
	if ok && locationValue != nil {
		record.Location, err = deserializeString(fmt.Sprintf("%s.%s", context, "location"), locationValue)
		if err != nil {
			return
		}
	}
	vMValue, ok := rpcStruct["VM"]
	if ok && vMValue != nil {
		record.VM, err = deserializeVMRef(fmt.Sprintf("%s.%s", context, "VM"), vMValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeConsoleRef(context string, input interface{}) (ConsoleRef, error) {
	var ref ConsoleRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return ConsoleRef(value), nil
}

func deserializeConsoleRefSet(context string, input interface{}) (slice []ConsoleRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]ConsoleRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeConsoleRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeConsoleRefToConsoleRecordMap(context string, input interface{}) (goMap map[ConsoleRef]ConsoleRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[ConsoleRef]ConsoleRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeConsoleRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeConsoleRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeCrashdumpRecord(context string, input interface{}) (record CrashdumpRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	vMValue, ok := rpcStruct["VM"]
	if ok && vMValue != nil {
		record.VM, err = deserializeVMRef(fmt.Sprintf("%s.%s", context, "VM"), vMValue)
		if err != nil {
			return
		}
	}
	vDIValue, ok := rpcStruct["VDI"]
	if ok && vDIValue != nil {
		record.VDI, err = deserializeVDIRef(fmt.Sprintf("%s.%s", context, "VDI"), vDIValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeCrashdumpRef(context string, input interface{}) (CrashdumpRef, error) {
	var ref CrashdumpRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return CrashdumpRef(value), nil
}

func deserializeCrashdumpRefSet(context string, input interface{}) (slice []CrashdumpRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]CrashdumpRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeCrashdumpRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeCrashdumpRefToCrashdumpRecordMap(context string, input interface{}) (goMap map[CrashdumpRef]CrashdumpRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[CrashdumpRef]CrashdumpRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeCrashdumpRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeCrashdumpRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeDRTaskRecord(context string, input interface{}) (record DRTaskRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	introducedSRsValue, ok := rpcStruct["introduced_SRs"]
	if ok && introducedSRsValue != nil {
		record.IntroducedSRs, err = deserializeSRRefSet(fmt.Sprintf("%s.%s", context, "introduced_SRs"), introducedSRsValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeDRTaskRef(context string, input interface{}) (DRTaskRef, error) {
	var ref DRTaskRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return DRTaskRef(value), nil
}

func deserializeDRTaskRefSet(context string, input interface{}) (slice []DRTaskRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]DRTaskRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeDRTaskRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeDRTaskRefToDRTaskRecordMap(context string, input interface{}) (goMap map[DRTaskRef]DRTaskRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[DRTaskRef]DRTaskRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeDRTaskRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeDRTaskRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeDataSourceRecord(context string, input interface{}) (record DataSourceRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	enabledValue, ok := rpcStruct["enabled"]
	if ok && enabledValue != nil {
		record.Enabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "enabled"), enabledValue)
		if err != nil {
			return
		}
	}
	standardValue, ok := rpcStruct["standard"]
	if ok && standardValue != nil {
		record.Standard, err = deserializeBool(fmt.Sprintf("%s.%s", context, "standard"), standardValue)
		if err != nil {
			return
		}
	}
	unitsValue, ok := rpcStruct["units"]
	if ok && unitsValue != nil {
		record.Units, err = deserializeString(fmt.Sprintf("%s.%s", context, "units"), unitsValue)
		if err != nil {
			return
		}
	}
	minValue, ok := rpcStruct["min"]
	if ok && minValue != nil {
		record.Min, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "min"), minValue)
		if err != nil {
			return
		}
	}
	maxValue, ok := rpcStruct["max"]
	if ok && maxValue != nil {
		record.Max, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "max"), maxValue)
		if err != nil {
			return
		}
	}
	valueValue, ok := rpcStruct["value"]
	if ok && valueValue != nil {
		record.Value, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "value"), valueValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeDataSourceRecordSet(context string, input interface{}) (slice []DataSourceRecord, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]DataSourceRecord, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeDataSourceRecord(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeDriverVariantRecord(context string, input interface{}) (record DriverVariantRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameValue, ok := rpcStruct["name"]
	if ok && nameValue != nil {
		record.Name, err = deserializeString(fmt.Sprintf("%s.%s", context, "name"), nameValue)
		if err != nil {
			return
		}
	}
	driverValue, ok := rpcStruct["driver"]
	if ok && driverValue != nil {
		record.Driver, err = deserializeHostDriverRef(fmt.Sprintf("%s.%s", context, "driver"), driverValue)
		if err != nil {
			return
		}
	}
	versionValue, ok := rpcStruct["version"]
	if ok && versionValue != nil {
		record.Version, err = deserializeString(fmt.Sprintf("%s.%s", context, "version"), versionValue)
		if err != nil {
			return
		}
	}
	hardwarePresentValue, ok := rpcStruct["hardware_present"]
	if ok && hardwarePresentValue != nil {
		record.HardwarePresent, err = deserializeBool(fmt.Sprintf("%s.%s", context, "hardware_present"), hardwarePresentValue)
		if err != nil {
			return
		}
	}
	priorityValue, ok := rpcStruct["priority"]
	if ok && priorityValue != nil {
		record.Priority, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "priority"), priorityValue)
		if err != nil {
			return
		}
	}
	statusValue, ok := rpcStruct["status"]
	if ok && statusValue != nil {
		record.Status, err = deserializeString(fmt.Sprintf("%s.%s", context, "status"), statusValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeDriverVariantRef(context string, input interface{}) (DriverVariantRef, error) {
	var ref DriverVariantRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return DriverVariantRef(value), nil
}

func deserializeDriverVariantRefSet(context string, input interface{}) (slice []DriverVariantRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]DriverVariantRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeDriverVariantRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeDriverVariantRefToDriverVariantRecordMap(context string, input interface{}) (goMap map[DriverVariantRef]DriverVariantRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[DriverVariantRef]DriverVariantRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeDriverVariantRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeDriverVariantRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeEnumAfterApplyGuidance(context string, input interface{}) (value AfterApplyGuidance, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "restartHVM":
		value = AfterApplyGuidanceRestartHVM
	case "restartPV":
		value = AfterApplyGuidanceRestartPV
	case "restartHost":
		value = AfterApplyGuidanceRestartHost
	case "restartXAPI":
		value = AfterApplyGuidanceRestartXAPI
	default:
		value = AfterApplyGuidanceUnrecognized
	}
	return
}

func deserializeEnumAfterApplyGuidanceSet(context string, input interface{}) (slice []AfterApplyGuidance, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]AfterApplyGuidance, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumAfterApplyGuidance(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumAllocationAlgorithm(context string, input interface{}) (value AllocationAlgorithm, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "breadth_first":
		value = AllocationAlgorithmBreadthFirst
	case "depth_first":
		value = AllocationAlgorithmDepthFirst
	default:
		value = AllocationAlgorithmUnrecognized
	}
	return
}

func deserializeEnumBondMode(context string, input interface{}) (value BondMode, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "balance-slb":
		value = BondModeBalanceSlb
	case "active-backup":
		value = BondModeActiveBackup
	case "lacp":
		value = BondModeLacp
	default:
		value = BondModeUnrecognized
	}
	return
}

func deserializeEnumCertificatePurpose(context string, input interface{}) (value CertificatePurpose, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "licensing":
		value = CertificatePurposeLicensing
	case "ldaps":
		value = CertificatePurposeLdaps
	default:
		value = CertificatePurposeUnrecognized
	}
	return
}

func deserializeEnumCertificatePurposeSet(context string, input interface{}) (slice []CertificatePurpose, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]CertificatePurpose, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumCertificatePurpose(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumCertificateType(context string, input interface{}) (value CertificateType, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "ca":
		value = CertificateTypeCa
	case "host":
		value = CertificateTypeHost
	case "host_internal":
		value = CertificateTypeHostInternal
	case "pinned":
		value = CertificateTypePinned
	default:
		value = CertificateTypeUnrecognized
	}
	return
}

func deserializeEnumCls(context string, input interface{}) (value Cls, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "VM":
		value = ClsVM
	case "Host":
		value = ClsHost
	case "SR":
		value = ClsSR
	case "Pool":
		value = ClsPool
	case "VMPP":
		value = ClsVMPP
	case "VMSS":
		value = ClsVMSS
	case "PVS_proxy":
		value = ClsPVSProxy
	case "VDI":
		value = ClsVDI
	case "Certificate":
		value = ClsCertificate
	default:
		value = ClsUnrecognized
	}
	return
}

func deserializeEnumClusterHostOperation(context string, input interface{}) (value ClusterHostOperation, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "enable":
		value = ClusterHostOperationEnable
	case "disable":
		value = ClusterHostOperationDisable
	case "destroy":
		value = ClusterHostOperationDestroy
	default:
		value = ClusterHostOperationUnrecognized
	}
	return
}

func deserializeEnumClusterHostOperationSet(context string, input interface{}) (slice []ClusterHostOperation, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]ClusterHostOperation, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumClusterHostOperation(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumClusterOperation(context string, input interface{}) (value ClusterOperation, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "add":
		value = ClusterOperationAdd
	case "remove":
		value = ClusterOperationRemove
	case "enable":
		value = ClusterOperationEnable
	case "disable":
		value = ClusterOperationDisable
	case "destroy":
		value = ClusterOperationDestroy
	default:
		value = ClusterOperationUnrecognized
	}
	return
}

func deserializeEnumClusterOperationSet(context string, input interface{}) (slice []ClusterOperation, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]ClusterOperation, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumClusterOperation(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumConsoleProtocol(context string, input interface{}) (value ConsoleProtocol, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "vt100":
		value = ConsoleProtocolVt100
	case "rfb":
		value = ConsoleProtocolRfb
	case "rdp":
		value = ConsoleProtocolRdp
	default:
		value = ConsoleProtocolUnrecognized
	}
	return
}

func deserializeEnumDomainType(context string, input interface{}) (value DomainType, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "hvm":
		value = DomainTypeHvm
	case "pv":
		value = DomainTypePv
	case "pv_in_pvh":
		value = DomainTypePvInPvh
	case "pvh":
		value = DomainTypePvh
	case "unspecified":
		value = DomainTypeUnspecified
	default:
		value = DomainTypeUnrecognized
	}
	return
}

func deserializeEnumEventOperation(context string, input interface{}) (value EventOperation, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "add":
		value = EventOperationAdd
	case "del":
		value = EventOperationDel
	case "mod":
		value = EventOperationMod
	default:
		value = EventOperationUnrecognized
	}
	return
}

func deserializeEnumHostAllowedOperations(context string, input interface{}) (value HostAllowedOperations, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "provision":
		value = HostAllowedOperationsProvision
	case "evacuate":
		value = HostAllowedOperationsEvacuate
	case "shutdown":
		value = HostAllowedOperationsShutdown
	case "reboot":
		value = HostAllowedOperationsReboot
	case "power_on":
		value = HostAllowedOperationsPowerOn
	case "vm_start":
		value = HostAllowedOperationsVMStart
	case "vm_resume":
		value = HostAllowedOperationsVMResume
	case "vm_migrate":
		value = HostAllowedOperationsVMMigrate
	case "apply_updates":
		value = HostAllowedOperationsApplyUpdates
	case "enable":
		value = HostAllowedOperationsEnable
	default:
		value = HostAllowedOperationsUnrecognized
	}
	return
}

func deserializeEnumHostAllowedOperationsSet(context string, input interface{}) (slice []HostAllowedOperations, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]HostAllowedOperations, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumHostAllowedOperations(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumHostDisplay(context string, input interface{}) (value HostDisplay, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "enabled":
		value = HostDisplayEnabled
	case "disable_on_reboot":
		value = HostDisplayDisableOnReboot
	case "disabled":
		value = HostDisplayDisabled
	case "enable_on_reboot":
		value = HostDisplayEnableOnReboot
	default:
		value = HostDisplayUnrecognized
	}
	return
}

func deserializeEnumHostNtpMode(context string, input interface{}) (value HostNtpMode, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "DHCP":
		value = HostNtpModeDHCP
	case "Custom":
		value = HostNtpModeCustom
	case "Factory":
		value = HostNtpModeFactory
	case "Disabled":
		value = HostNtpModeDisabled
	default:
		value = HostNtpModeUnrecognized
	}
	return
}

func deserializeEnumHostNumaAffinityPolicy(context string, input interface{}) (value HostNumaAffinityPolicy, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "any":
		value = HostNumaAffinityPolicyAny
	case "best_effort":
		value = HostNumaAffinityPolicyBestEffort
	case "default_policy":
		value = HostNumaAffinityPolicyDefaultPolicy
	default:
		value = HostNumaAffinityPolicyUnrecognized
	}
	return
}

func deserializeEnumHostSchedGran(context string, input interface{}) (value HostSchedGran, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "core":
		value = HostSchedGranCore
	case "cpu":
		value = HostSchedGranCPU
	case "socket":
		value = HostSchedGranSocket
	default:
		value = HostSchedGranUnrecognized
	}
	return
}

func deserializeEnumIPConfigurationMode(context string, input interface{}) (value IPConfigurationMode, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "None":
		value = IPConfigurationModeNone
	case "DHCP":
		value = IPConfigurationModeDHCP
	case "Static":
		value = IPConfigurationModeStatic
	default:
		value = IPConfigurationModeUnrecognized
	}
	return
}

func deserializeEnumImageFormatType(context string, input interface{}) (value ImageFormatType, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "raw":
		value = ImageFormatTypeRaw
	case "vhd":
		value = ImageFormatTypeVhd
	case "qcow2":
		value = ImageFormatTypeQcow2
	default:
		value = ImageFormatTypeUnrecognized
	}
	return
}

func deserializeEnumImageFormatTypeSet(context string, input interface{}) (slice []ImageFormatType, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]ImageFormatType, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumImageFormatType(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumIpv6ConfigurationMode(context string, input interface{}) (value Ipv6ConfigurationMode, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "None":
		value = Ipv6ConfigurationModeNone
	case "DHCP":
		value = Ipv6ConfigurationModeDHCP
	case "Static":
		value = Ipv6ConfigurationModeStatic
	case "Autoconf":
		value = Ipv6ConfigurationModeAutoconf
	default:
		value = Ipv6ConfigurationModeUnrecognized
	}
	return
}

func deserializeEnumLatestSyncedUpdatesAppliedState(context string, input interface{}) (value LatestSyncedUpdatesAppliedState, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "yes":
		value = LatestSyncedUpdatesAppliedStateYes
	case "no":
		value = LatestSyncedUpdatesAppliedStateNo
	case "unknown":
		value = LatestSyncedUpdatesAppliedStateUnknown
	default:
		value = LatestSyncedUpdatesAppliedStateUnrecognized
	}
	return
}

func deserializeEnumLivepatchStatus(context string, input interface{}) (value LivepatchStatus, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "ok_livepatch_complete":
		value = LivepatchStatusOkLivepatchComplete
	case "ok_livepatch_incomplete":
		value = LivepatchStatusOkLivepatchIncomplete
	case "ok":
		value = LivepatchStatusOk
	default:
		value = LivepatchStatusUnrecognized
	}
	return
}

func deserializeEnumNetworkDefaultLockingMode(context string, input interface{}) (value NetworkDefaultLockingMode, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "unlocked":
		value = NetworkDefaultLockingModeUnlocked
	case "disabled":
		value = NetworkDefaultLockingModeDisabled
	default:
		value = NetworkDefaultLockingModeUnrecognized
	}
	return
}

func deserializeEnumNetworkOperations(context string, input interface{}) (value NetworkOperations, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "attaching":
		value = NetworkOperationsAttaching
	default:
		value = NetworkOperationsUnrecognized
	}
	return
}

func deserializeEnumNetworkOperationsSet(context string, input interface{}) (slice []NetworkOperations, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]NetworkOperations, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumNetworkOperations(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumNetworkPurpose(context string, input interface{}) (value NetworkPurpose, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "nbd":
		value = NetworkPurposeNbd
	case "insecure_nbd":
		value = NetworkPurposeInsecureNbd
	default:
		value = NetworkPurposeUnrecognized
	}
	return
}

func deserializeEnumNetworkPurposeSet(context string, input interface{}) (slice []NetworkPurpose, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]NetworkPurpose, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumNetworkPurpose(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumOnBoot(context string, input interface{}) (value OnBoot, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "reset":
		value = OnBootReset
	case "persist":
		value = OnBootPersist
	default:
		value = OnBootUnrecognized
	}
	return
}

func deserializeEnumOnCrashBehaviour(context string, input interface{}) (value OnCrashBehaviour, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "destroy":
		value = OnCrashBehaviourDestroy
	case "coredump_and_destroy":
		value = OnCrashBehaviourCoredumpAndDestroy
	case "restart":
		value = OnCrashBehaviourRestart
	case "coredump_and_restart":
		value = OnCrashBehaviourCoredumpAndRestart
	case "preserve":
		value = OnCrashBehaviourPreserve
	case "rename_restart":
		value = OnCrashBehaviourRenameRestart
	default:
		value = OnCrashBehaviourUnrecognized
	}
	return
}

func deserializeEnumOnNormalExit(context string, input interface{}) (value OnNormalExit, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "destroy":
		value = OnNormalExitDestroy
	case "restart":
		value = OnNormalExitRestart
	default:
		value = OnNormalExitUnrecognized
	}
	return
}

func deserializeEnumOnSoftrebootBehavior(context string, input interface{}) (value OnSoftrebootBehavior, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "soft_reboot":
		value = OnSoftrebootBehaviorSoftReboot
	case "destroy":
		value = OnSoftrebootBehaviorDestroy
	case "restart":
		value = OnSoftrebootBehaviorRestart
	case "preserve":
		value = OnSoftrebootBehaviorPreserve
	default:
		value = OnSoftrebootBehaviorUnrecognized
	}
	return
}

func deserializeEnumOrigin(context string, input interface{}) (value Origin, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "remote":
		value = OriginRemote
	case "bundle":
		value = OriginBundle
	case "remote_pool":
		value = OriginRemotePool
	default:
		value = OriginUnrecognized
	}
	return
}

func deserializeEnumPciDom0Access(context string, input interface{}) (value PciDom0Access, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "enabled":
		value = PciDom0AccessEnabled
	case "disable_on_reboot":
		value = PciDom0AccessDisableOnReboot
	case "disabled":
		value = PciDom0AccessDisabled
	case "enable_on_reboot":
		value = PciDom0AccessEnableOnReboot
	default:
		value = PciDom0AccessUnrecognized
	}
	return
}

func deserializeEnumPersistenceBackend(context string, input interface{}) (value PersistenceBackend, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "xapi":
		value = PersistenceBackendXapi
	default:
		value = PersistenceBackendUnrecognized
	}
	return
}

func deserializeEnumPifIgmpStatus(context string, input interface{}) (value PifIgmpStatus, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "enabled":
		value = PifIgmpStatusEnabled
	case "disabled":
		value = PifIgmpStatusDisabled
	case "unknown":
		value = PifIgmpStatusUnknown
	default:
		value = PifIgmpStatusUnrecognized
	}
	return
}

func deserializeEnumPlacementPolicy(context string, input interface{}) (value PlacementPolicy, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "anti_affinity":
		value = PlacementPolicyAntiAffinity
	case "normal":
		value = PlacementPolicyNormal
	default:
		value = PlacementPolicyUnrecognized
	}
	return
}

func deserializeEnumPoolAllowedOperations(context string, input interface{}) (value PoolAllowedOperations, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "ha_enable":
		value = PoolAllowedOperationsHaEnable
	case "ha_disable":
		value = PoolAllowedOperationsHaDisable
	case "cluster_create":
		value = PoolAllowedOperationsClusterCreate
	case "designate_new_master":
		value = PoolAllowedOperationsDesignateNewMaster
	case "configure_repositories":
		value = PoolAllowedOperationsConfigureRepositories
	case "sync_updates":
		value = PoolAllowedOperationsSyncUpdates
	case "sync_bundle":
		value = PoolAllowedOperationsSyncBundle
	case "get_updates":
		value = PoolAllowedOperationsGetUpdates
	case "apply_updates":
		value = PoolAllowedOperationsApplyUpdates
	case "tls_verification_enable":
		value = PoolAllowedOperationsTLSVerificationEnable
	case "cert_refresh":
		value = PoolAllowedOperationsCertRefresh
	case "exchange_certificates_on_join":
		value = PoolAllowedOperationsExchangeCertificatesOnJoin
	case "exchange_ca_certificates_on_join":
		value = PoolAllowedOperationsExchangeCaCertificatesOnJoin
	case "copy_primary_host_certs":
		value = PoolAllowedOperationsCopyPrimaryHostCerts
	case "eject":
		value = PoolAllowedOperationsEject
	case "exchange_crls_on_join":
		value = PoolAllowedOperationsExchangeCrlsOnJoin
	default:
		value = PoolAllowedOperationsUnrecognized
	}
	return
}

func deserializeEnumPoolAllowedOperationsSet(context string, input interface{}) (slice []PoolAllowedOperations, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PoolAllowedOperations, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumPoolAllowedOperations(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumPoolGuestSecurebootReadiness(context string, input interface{}) (value PoolGuestSecurebootReadiness, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "ready":
		value = PoolGuestSecurebootReadinessReady
	case "ready_no_dbx":
		value = PoolGuestSecurebootReadinessReadyNoDbx
	case "not_ready":
		value = PoolGuestSecurebootReadinessNotReady
	default:
		value = PoolGuestSecurebootReadinessUnrecognized
	}
	return
}

func deserializeEnumPrimaryAddressType(context string, input interface{}) (value PrimaryAddressType, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "IPv4":
		value = PrimaryAddressTypeIPv4
	case "IPv6":
		value = PrimaryAddressTypeIPv6
	default:
		value = PrimaryAddressTypeUnrecognized
	}
	return
}

func deserializeEnumPvsProxyStatus(context string, input interface{}) (value PvsProxyStatus, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "stopped":
		value = PvsProxyStatusStopped
	case "initialised":
		value = PvsProxyStatusInitialised
	case "caching":
		value = PvsProxyStatusCaching
	case "incompatible_write_cache_mode":
		value = PvsProxyStatusIncompatibleWriteCacheMode
	case "incompatible_protocol_version":
		value = PvsProxyStatusIncompatibleProtocolVersion
	default:
		value = PvsProxyStatusUnrecognized
	}
	return
}

func deserializeEnumSdnControllerProtocol(context string, input interface{}) (value SdnControllerProtocol, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "ssl":
		value = SdnControllerProtocolSsl
	case "pssl":
		value = SdnControllerProtocolPssl
	default:
		value = SdnControllerProtocolUnrecognized
	}
	return
}

func deserializeEnumSrHealth(context string, input interface{}) (value SrHealth, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "healthy":
		value = SrHealthHealthy
	case "recovering":
		value = SrHealthRecovering
	case "unreachable":
		value = SrHealthUnreachable
	case "unavailable":
		value = SrHealthUnavailable
	default:
		value = SrHealthUnrecognized
	}
	return
}

func deserializeEnumSriovConfigurationMode(context string, input interface{}) (value SriovConfigurationMode, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "sysfs":
		value = SriovConfigurationModeSysfs
	case "modprobe":
		value = SriovConfigurationModeModprobe
	case "manual":
		value = SriovConfigurationModeManual
	case "unknown":
		value = SriovConfigurationModeUnknown
	default:
		value = SriovConfigurationModeUnrecognized
	}
	return
}

func deserializeEnumStorageOperations(context string, input interface{}) (value StorageOperations, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "scan":
		value = StorageOperationsScan
	case "destroy":
		value = StorageOperationsDestroy
	case "forget":
		value = StorageOperationsForget
	case "plug":
		value = StorageOperationsPlug
	case "unplug":
		value = StorageOperationsUnplug
	case "update":
		value = StorageOperationsUpdate
	case "vdi_create":
		value = StorageOperationsVdiCreate
	case "vdi_introduce":
		value = StorageOperationsVdiIntroduce
	case "vdi_destroy":
		value = StorageOperationsVdiDestroy
	case "vdi_resize":
		value = StorageOperationsVdiResize
	case "vdi_clone":
		value = StorageOperationsVdiClone
	case "vdi_snapshot":
		value = StorageOperationsVdiSnapshot
	case "vdi_mirror":
		value = StorageOperationsVdiMirror
	case "vdi_enable_cbt":
		value = StorageOperationsVdiEnableCbt
	case "vdi_disable_cbt":
		value = StorageOperationsVdiDisableCbt
	case "vdi_data_destroy":
		value = StorageOperationsVdiDataDestroy
	case "vdi_list_changed_blocks":
		value = StorageOperationsVdiListChangedBlocks
	case "vdi_set_on_boot":
		value = StorageOperationsVdiSetOnBoot
	case "vdi_blocked":
		value = StorageOperationsVdiBlocked
	case "vdi_copy":
		value = StorageOperationsVdiCopy
	case "vdi_force_unlock":
		value = StorageOperationsVdiForceUnlock
	case "vdi_forget":
		value = StorageOperationsVdiForget
	case "vdi_generate_config":
		value = StorageOperationsVdiGenerateConfig
	case "vdi_resize_online":
		value = StorageOperationsVdiResizeOnline
	case "vdi_update":
		value = StorageOperationsVdiUpdate
	case "vdi_revert":
		value = StorageOperationsVdiRevert
	case "pbd_create":
		value = StorageOperationsPbdCreate
	case "pbd_destroy":
		value = StorageOperationsPbdDestroy
	default:
		value = StorageOperationsUnrecognized
	}
	return
}

func deserializeEnumStorageOperationsSet(context string, input interface{}) (slice []StorageOperations, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]StorageOperations, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumStorageOperations(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumTaskAllowedOperations(context string, input interface{}) (value TaskAllowedOperations, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "cancel":
		value = TaskAllowedOperationsCancel
	case "destroy":
		value = TaskAllowedOperationsDestroy
	default:
		value = TaskAllowedOperationsUnrecognized
	}
	return
}

func deserializeEnumTaskAllowedOperationsSet(context string, input interface{}) (slice []TaskAllowedOperations, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]TaskAllowedOperations, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumTaskAllowedOperations(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumTaskStatusType(context string, input interface{}) (value TaskStatusType, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "pending":
		value = TaskStatusTypePending
	case "success":
		value = TaskStatusTypeSuccess
	case "failure":
		value = TaskStatusTypeFailure
	case "cancelling":
		value = TaskStatusTypeCancelling
	case "cancelled":
		value = TaskStatusTypeCancelled
	default:
		value = TaskStatusTypeUnrecognized
	}
	return
}

func deserializeEnumTelemetryFrequency(context string, input interface{}) (value TelemetryFrequency, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "daily":
		value = TelemetryFrequencyDaily
	case "weekly":
		value = TelemetryFrequencyWeekly
	case "monthly":
		value = TelemetryFrequencyMonthly
	default:
		value = TelemetryFrequencyUnrecognized
	}
	return
}

func deserializeEnumTristateType(context string, input interface{}) (value TristateType, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "yes":
		value = TristateTypeYes
	case "no":
		value = TristateTypeNo
	case "unspecified":
		value = TristateTypeUnspecified
	default:
		value = TristateTypeUnrecognized
	}
	return
}

func deserializeEnumTunnelProtocol(context string, input interface{}) (value TunnelProtocol, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "gre":
		value = TunnelProtocolGre
	case "vxlan":
		value = TunnelProtocolVxlan
	default:
		value = TunnelProtocolUnrecognized
	}
	return
}

func deserializeEnumUpdateAfterApplyGuidance(context string, input interface{}) (value UpdateAfterApplyGuidance, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "restartHVM":
		value = UpdateAfterApplyGuidanceRestartHVM
	case "restartPV":
		value = UpdateAfterApplyGuidanceRestartPV
	case "restartHost":
		value = UpdateAfterApplyGuidanceRestartHost
	case "restartXAPI":
		value = UpdateAfterApplyGuidanceRestartXAPI
	default:
		value = UpdateAfterApplyGuidanceUnrecognized
	}
	return
}

func deserializeEnumUpdateAfterApplyGuidanceSet(context string, input interface{}) (slice []UpdateAfterApplyGuidance, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]UpdateAfterApplyGuidance, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumUpdateAfterApplyGuidance(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumUpdateGuidances(context string, input interface{}) (value UpdateGuidances, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "reboot_host":
		value = UpdateGuidancesRebootHost
	case "reboot_host_on_livepatch_failure":
		value = UpdateGuidancesRebootHostOnLivepatchFailure
	case "reboot_host_on_kernel_livepatch_failure":
		value = UpdateGuidancesRebootHostOnKernelLivepatchFailure
	case "reboot_host_on_xen_livepatch_failure":
		value = UpdateGuidancesRebootHostOnXenLivepatchFailure
	case "restart_toolstack":
		value = UpdateGuidancesRestartToolstack
	case "restart_device_model":
		value = UpdateGuidancesRestartDeviceModel
	case "restart_vm":
		value = UpdateGuidancesRestartVM
	default:
		value = UpdateGuidancesUnrecognized
	}
	return
}

func deserializeEnumUpdateGuidancesSet(context string, input interface{}) (slice []UpdateGuidances, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]UpdateGuidances, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumUpdateGuidances(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumUpdateSyncFrequency(context string, input interface{}) (value UpdateSyncFrequency, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "daily":
		value = UpdateSyncFrequencyDaily
	case "weekly":
		value = UpdateSyncFrequencyWeekly
	default:
		value = UpdateSyncFrequencyUnrecognized
	}
	return
}

func deserializeEnumVMApplianceOperation(context string, input interface{}) (value VMApplianceOperation, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "start":
		value = VMApplianceOperationStart
	case "clean_shutdown":
		value = VMApplianceOperationCleanShutdown
	case "hard_shutdown":
		value = VMApplianceOperationHardShutdown
	case "shutdown":
		value = VMApplianceOperationShutdown
	default:
		value = VMApplianceOperationUnrecognized
	}
	return
}

func deserializeEnumVMApplianceOperationSet(context string, input interface{}) (slice []VMApplianceOperation, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VMApplianceOperation, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumVMApplianceOperation(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumVMOperations(context string, input interface{}) (value VMOperations, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "snapshot":
		value = VMOperationsSnapshot
	case "clone":
		value = VMOperationsClone
	case "copy":
		value = VMOperationsCopy
	case "create_template":
		value = VMOperationsCreateTemplate
	case "revert":
		value = VMOperationsRevert
	case "checkpoint":
		value = VMOperationsCheckpoint
	case "snapshot_with_quiesce":
		value = VMOperationsSnapshotWithQuiesce
	case "provision":
		value = VMOperationsProvision
	case "start":
		value = VMOperationsStart
	case "start_on":
		value = VMOperationsStartOn
	case "pause":
		value = VMOperationsPause
	case "unpause":
		value = VMOperationsUnpause
	case "clean_shutdown":
		value = VMOperationsCleanShutdown
	case "clean_reboot":
		value = VMOperationsCleanReboot
	case "hard_shutdown":
		value = VMOperationsHardShutdown
	case "power_state_reset":
		value = VMOperationsPowerStateReset
	case "hard_reboot":
		value = VMOperationsHardReboot
	case "suspend":
		value = VMOperationsSuspend
	case "csvm":
		value = VMOperationsCsvm
	case "resume":
		value = VMOperationsResume
	case "resume_on":
		value = VMOperationsResumeOn
	case "pool_migrate":
		value = VMOperationsPoolMigrate
	case "migrate_send":
		value = VMOperationsMigrateSend
	case "get_boot_record":
		value = VMOperationsGetBootRecord
	case "send_sysrq":
		value = VMOperationsSendSysrq
	case "send_trigger":
		value = VMOperationsSendTrigger
	case "query_services":
		value = VMOperationsQueryServices
	case "shutdown":
		value = VMOperationsShutdown
	case "call_plugin":
		value = VMOperationsCallPlugin
	case "changing_memory_live":
		value = VMOperationsChangingMemoryLive
	case "awaiting_memory_live":
		value = VMOperationsAwaitingMemoryLive
	case "changing_dynamic_range":
		value = VMOperationsChangingDynamicRange
	case "changing_static_range":
		value = VMOperationsChangingStaticRange
	case "changing_memory_limits":
		value = VMOperationsChangingMemoryLimits
	case "changing_shadow_memory":
		value = VMOperationsChangingShadowMemory
	case "changing_shadow_memory_live":
		value = VMOperationsChangingShadowMemoryLive
	case "changing_VCPUs":
		value = VMOperationsChangingVCPUs
	case "changing_VCPUs_live":
		value = VMOperationsChangingVCPUsLive
	case "changing_NVRAM":
		value = VMOperationsChangingNVRAM
	case "assert_operation_valid":
		value = VMOperationsAssertOperationValid
	case "data_source_op":
		value = VMOperationsDataSourceOp
	case "update_allowed_operations":
		value = VMOperationsUpdateAllowedOperations
	case "make_into_template":
		value = VMOperationsMakeIntoTemplate
	case "import":
		value = VMOperationsImport
	case "export":
		value = VMOperationsExport
	case "metadata_export":
		value = VMOperationsMetadataExport
	case "reverting":
		value = VMOperationsReverting
	case "destroy":
		value = VMOperationsDestroy
	case "create_vtpm":
		value = VMOperationsCreateVtpm
	case "sysprep":
		value = VMOperationsSysprep
	default:
		value = VMOperationsUnrecognized
	}
	return
}

func deserializeEnumVMOperationsSet(context string, input interface{}) (slice []VMOperations, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VMOperations, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumVMOperations(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumVMOperationsToStringMap(context string, input interface{}) (goMap map[VMOperations]string, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VMOperations]string, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeEnumVMOperations(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeString(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeEnumVMPowerState(context string, input interface{}) (value VMPowerState, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "Halted":
		value = VMPowerStateHalted
	case "Paused":
		value = VMPowerStatePaused
	case "Running":
		value = VMPowerStateRunning
	case "Suspended":
		value = VMPowerStateSuspended
	default:
		value = VMPowerStateUnrecognized
	}
	return
}

func deserializeEnumVMSecurebootCertificatesState(context string, input interface{}) (value VMSecurebootCertificatesState, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "ok":
		value = VMSecurebootCertificatesStateOk
	case "update_available":
		value = VMSecurebootCertificatesStateUpdateAvailable
	case "update_on_boot":
		value = VMSecurebootCertificatesStateUpdateOnBoot
	default:
		value = VMSecurebootCertificatesStateUnrecognized
	}
	return
}

func deserializeEnumVMSecurebootReadiness(context string, input interface{}) (value VMSecurebootReadiness, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "not_supported":
		value = VMSecurebootReadinessNotSupported
	case "disabled":
		value = VMSecurebootReadinessDisabled
	case "first_boot":
		value = VMSecurebootReadinessFirstBoot
	case "ready":
		value = VMSecurebootReadinessReady
	case "ready_no_dbx":
		value = VMSecurebootReadinessReadyNoDbx
	case "setup_mode":
		value = VMSecurebootReadinessSetupMode
	case "certs_incomplete":
		value = VMSecurebootReadinessCertsIncomplete
	default:
		value = VMSecurebootReadinessUnrecognized
	}
	return
}

func deserializeEnumVbdMode(context string, input interface{}) (value VbdMode, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "RO":
		value = VbdModeRO
	case "RW":
		value = VbdModeRW
	default:
		value = VbdModeUnrecognized
	}
	return
}

func deserializeEnumVbdOperations(context string, input interface{}) (value VbdOperations, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "attach":
		value = VbdOperationsAttach
	case "eject":
		value = VbdOperationsEject
	case "insert":
		value = VbdOperationsInsert
	case "plug":
		value = VbdOperationsPlug
	case "unplug":
		value = VbdOperationsUnplug
	case "unplug_force":
		value = VbdOperationsUnplugForce
	case "pause":
		value = VbdOperationsPause
	case "unpause":
		value = VbdOperationsUnpause
	default:
		value = VbdOperationsUnrecognized
	}
	return
}

func deserializeEnumVbdOperationsSet(context string, input interface{}) (slice []VbdOperations, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VbdOperations, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumVbdOperations(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumVbdType(context string, input interface{}) (value VbdType, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "CD":
		value = VbdTypeCD
	case "Disk":
		value = VbdTypeDisk
	case "Floppy":
		value = VbdTypeFloppy
	default:
		value = VbdTypeUnrecognized
	}
	return
}

func deserializeEnumVdiOperations(context string, input interface{}) (value VdiOperations, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "clone":
		value = VdiOperationsClone
	case "copy":
		value = VdiOperationsCopy
	case "resize":
		value = VdiOperationsResize
	case "resize_online":
		value = VdiOperationsResizeOnline
	case "snapshot":
		value = VdiOperationsSnapshot
	case "mirror":
		value = VdiOperationsMirror
	case "destroy":
		value = VdiOperationsDestroy
	case "forget":
		value = VdiOperationsForget
	case "update":
		value = VdiOperationsUpdate
	case "force_unlock":
		value = VdiOperationsForceUnlock
	case "generate_config":
		value = VdiOperationsGenerateConfig
	case "enable_cbt":
		value = VdiOperationsEnableCbt
	case "disable_cbt":
		value = VdiOperationsDisableCbt
	case "data_destroy":
		value = VdiOperationsDataDestroy
	case "list_changed_blocks":
		value = VdiOperationsListChangedBlocks
	case "set_on_boot":
		value = VdiOperationsSetOnBoot
	case "blocked":
		value = VdiOperationsBlocked
	case "revert_to":
		value = VdiOperationsRevertTo
	case "revert_from":
		value = VdiOperationsRevertFrom
	default:
		value = VdiOperationsUnrecognized
	}
	return
}

func deserializeEnumVdiOperationsSet(context string, input interface{}) (slice []VdiOperations, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VdiOperations, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumVdiOperations(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumVdiType(context string, input interface{}) (value VdiType, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "system":
		value = VdiTypeSystem
	case "user":
		value = VdiTypeUser
	case "ephemeral":
		value = VdiTypeEphemeral
	case "suspend":
		value = VdiTypeSuspend
	case "crashdump":
		value = VdiTypeCrashdump
	case "ha_statefile":
		value = VdiTypeHaStatefile
	case "metadata":
		value = VdiTypeMetadata
	case "redo_log":
		value = VdiTypeRedoLog
	case "rrd":
		value = VdiTypeRrd
	case "pvs_cache":
		value = VdiTypePvsCache
	case "cbt_metadata":
		value = VdiTypeCbtMetadata
	default:
		value = VdiTypeUnrecognized
	}
	return
}

func deserializeEnumVgpuTypeImplementation(context string, input interface{}) (value VgpuTypeImplementation, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "passthrough":
		value = VgpuTypeImplementationPassthrough
	case "nvidia":
		value = VgpuTypeImplementationNvidia
	case "nvidia_sriov":
		value = VgpuTypeImplementationNvidiaSriov
	case "gvt_g":
		value = VgpuTypeImplementationGvtG
	case "mxgpu":
		value = VgpuTypeImplementationMxgpu
	default:
		value = VgpuTypeImplementationUnrecognized
	}
	return
}

func deserializeEnumVifIpv4ConfigurationMode(context string, input interface{}) (value VifIpv4ConfigurationMode, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "None":
		value = VifIpv4ConfigurationModeNone
	case "Static":
		value = VifIpv4ConfigurationModeStatic
	case "DHCP":
		value = VifIpv4ConfigurationModeDHCP
	default:
		value = VifIpv4ConfigurationModeUnrecognized
	}
	return
}

func deserializeEnumVifIpv6ConfigurationMode(context string, input interface{}) (value VifIpv6ConfigurationMode, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "None":
		value = VifIpv6ConfigurationModeNone
	case "Static":
		value = VifIpv6ConfigurationModeStatic
	case "Autoconf":
		value = VifIpv6ConfigurationModeAutoconf
	default:
		value = VifIpv6ConfigurationModeUnrecognized
	}
	return
}

func deserializeEnumVifLockingMode(context string, input interface{}) (value VifLockingMode, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "network_default":
		value = VifLockingModeNetworkDefault
	case "locked":
		value = VifLockingModeLocked
	case "unlocked":
		value = VifLockingModeUnlocked
	case "disabled":
		value = VifLockingModeDisabled
	default:
		value = VifLockingModeUnrecognized
	}
	return
}

func deserializeEnumVifOperations(context string, input interface{}) (value VifOperations, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "attach":
		value = VifOperationsAttach
	case "plug":
		value = VifOperationsPlug
	case "unplug":
		value = VifOperationsUnplug
	default:
		value = VifOperationsUnrecognized
	}
	return
}

func deserializeEnumVifOperationsSet(context string, input interface{}) (slice []VifOperations, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VifOperations, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumVifOperations(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumVmppArchiveFrequency(context string, input interface{}) (value VmppArchiveFrequency, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "never":
		value = VmppArchiveFrequencyNever
	case "always_after_backup":
		value = VmppArchiveFrequencyAlwaysAfterBackup
	case "daily":
		value = VmppArchiveFrequencyDaily
	case "weekly":
		value = VmppArchiveFrequencyWeekly
	default:
		value = VmppArchiveFrequencyUnrecognized
	}
	return
}

func deserializeEnumVmppArchiveTargetType(context string, input interface{}) (value VmppArchiveTargetType, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "none":
		value = VmppArchiveTargetTypeNone
	case "cifs":
		value = VmppArchiveTargetTypeCifs
	case "nfs":
		value = VmppArchiveTargetTypeNfs
	default:
		value = VmppArchiveTargetTypeUnrecognized
	}
	return
}

func deserializeEnumVmppBackupFrequency(context string, input interface{}) (value VmppBackupFrequency, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "hourly":
		value = VmppBackupFrequencyHourly
	case "daily":
		value = VmppBackupFrequencyDaily
	case "weekly":
		value = VmppBackupFrequencyWeekly
	default:
		value = VmppBackupFrequencyUnrecognized
	}
	return
}

func deserializeEnumVmppBackupType(context string, input interface{}) (value VmppBackupType, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "snapshot":
		value = VmppBackupTypeSnapshot
	case "checkpoint":
		value = VmppBackupTypeCheckpoint
	default:
		value = VmppBackupTypeUnrecognized
	}
	return
}

func deserializeEnumVmssFrequency(context string, input interface{}) (value VmssFrequency, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "hourly":
		value = VmssFrequencyHourly
	case "daily":
		value = VmssFrequencyDaily
	case "weekly":
		value = VmssFrequencyWeekly
	default:
		value = VmssFrequencyUnrecognized
	}
	return
}

func deserializeEnumVmssType(context string, input interface{}) (value VmssType, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "snapshot":
		value = VmssTypeSnapshot
	case "checkpoint":
		value = VmssTypeCheckpoint
	case "snapshot_with_quiesce":
		value = VmssTypeSnapshotWithQuiesce
	default:
		value = VmssTypeUnrecognized
	}
	return
}

func deserializeEnumVtpmOperations(context string, input interface{}) (value VtpmOperations, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "destroy":
		value = VtpmOperationsDestroy
	default:
		value = VtpmOperationsUnrecognized
	}
	return
}

func deserializeEnumVtpmOperationsSet(context string, input interface{}) (slice []VtpmOperations, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VtpmOperations, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumVtpmOperations(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEnumVusbOperations(context string, input interface{}) (value VusbOperations, err error) {
	strValue, err := deserializeString(context, input)
	if err != nil {
		return
	}
	switch strValue {
	case "attach":
		value = VusbOperationsAttach
	case "plug":
		value = VusbOperationsPlug
	case "unplug":
		value = VusbOperationsUnplug
	default:
		value = VusbOperationsUnrecognized
	}
	return
}

func deserializeEnumVusbOperationsSet(context string, input interface{}) (slice []VusbOperations, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VusbOperations, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEnumVusbOperations(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEventRecord(context string, input interface{}) (record EventRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	iDValue, ok := rpcStruct["id"]
	if ok && iDValue != nil {
		record.ID, err = deserializeInt(fmt.Sprintf("%s.%s", context, "id"), iDValue)
		if err != nil {
			return
		}
	}
	timestampValue, ok := rpcStruct["timestamp"]
	if ok && timestampValue != nil {
		record.Timestamp, err = deserializeTime(fmt.Sprintf("%s.%s", context, "timestamp"), timestampValue)
		if err != nil {
			return
		}
	}
	classValue, ok := rpcStruct["class"]
	if ok && classValue != nil {
		record.Class, err = deserializeString(fmt.Sprintf("%s.%s", context, "class"), classValue)
		if err != nil {
			return
		}
	}
	operationValue, ok := rpcStruct["operation"]
	if ok && operationValue != nil {
		record.Operation, err = deserializeEnumEventOperation(fmt.Sprintf("%s.%s", context, "operation"), operationValue)
		if err != nil {
			return
		}
	}
	refValue, ok := rpcStruct["ref"]
	if ok && refValue != nil {
		record.Ref, err = deserializeString(fmt.Sprintf("%s.%s", context, "ref"), refValue)
		if err != nil {
			return
		}
	}
	objUUIDValue, ok := rpcStruct["obj_uuid"]
	if ok && objUUIDValue != nil {
		record.ObjUUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "obj_uuid"), objUUIDValue)
		if err != nil {
			return
		}
	}
	snapshotValue, ok := rpcStruct["snapshot"]
	if ok && snapshotValue != nil {
		record.Snapshot, err = deserializeRecordInterface(fmt.Sprintf("%s.%s", context, "snapshot"), snapshotValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeEventRecordSet(context string, input interface{}) (slice []EventRecord, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]EventRecord, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeEventRecord(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeFeatureRecord(context string, input interface{}) (record FeatureRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	enabledValue, ok := rpcStruct["enabled"]
	if ok && enabledValue != nil {
		record.Enabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "enabled"), enabledValue)
		if err != nil {
			return
		}
	}
	experimentalValue, ok := rpcStruct["experimental"]
	if ok && experimentalValue != nil {
		record.Experimental, err = deserializeBool(fmt.Sprintf("%s.%s", context, "experimental"), experimentalValue)
		if err != nil {
			return
		}
	}
	versionValue, ok := rpcStruct["version"]
	if ok && versionValue != nil {
		record.Version, err = deserializeString(fmt.Sprintf("%s.%s", context, "version"), versionValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeFeatureRef(context string, input interface{}) (FeatureRef, error) {
	var ref FeatureRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return FeatureRef(value), nil
}

func deserializeFeatureRefSet(context string, input interface{}) (slice []FeatureRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]FeatureRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeFeatureRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeFeatureRefToFeatureRecordMap(context string, input interface{}) (goMap map[FeatureRef]FeatureRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[FeatureRef]FeatureRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeFeatureRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeFeatureRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeFloat(context string, input interface{}) (value float64, err error) {
	_ = context
	if input == nil {
		return
	}
	strValue := fmt.Sprintf("%v", input)
	value, err = strconv.ParseFloat(strValue, 64)
	if err != nil {
		switch strValue {
		case "+Inf":
			return math.Inf(1), nil
		case "-Inf":
			return math.Inf(-1), nil
		case "NaN":
			return math.NaN(), nil
		}
	}
	return
}

func deserializeGPUGroupRecord(context string, input interface{}) (record GPUGroupRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	pGPUsValue, ok := rpcStruct["PGPUs"]
	if ok && pGPUsValue != nil {
		record.PGPUs, err = deserializePGPURefSet(fmt.Sprintf("%s.%s", context, "PGPUs"), pGPUsValue)
		if err != nil {
			return
		}
	}
	vGPUsValue, ok := rpcStruct["VGPUs"]
	if ok && vGPUsValue != nil {
		record.VGPUs, err = deserializeVGPURefSet(fmt.Sprintf("%s.%s", context, "VGPUs"), vGPUsValue)
		if err != nil {
			return
		}
	}
	gPUTypesValue, ok := rpcStruct["GPU_types"]
	if ok && gPUTypesValue != nil {
		record.GPUTypes, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "GPU_types"), gPUTypesValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	allocationAlgorithmValue, ok := rpcStruct["allocation_algorithm"]
	if ok && allocationAlgorithmValue != nil {
		record.AllocationAlgorithm, err = deserializeEnumAllocationAlgorithm(fmt.Sprintf("%s.%s", context, "allocation_algorithm"), allocationAlgorithmValue)
		if err != nil {
			return
		}
	}
	supportedVGPUTypesValue, ok := rpcStruct["supported_VGPU_types"]
	if ok && supportedVGPUTypesValue != nil {
		record.SupportedVGPUTypes, err = deserializeVGPUTypeRefSet(fmt.Sprintf("%s.%s", context, "supported_VGPU_types"), supportedVGPUTypesValue)
		if err != nil {
			return
		}
	}
	enabledVGPUTypesValue, ok := rpcStruct["enabled_VGPU_types"]
	if ok && enabledVGPUTypesValue != nil {
		record.EnabledVGPUTypes, err = deserializeVGPUTypeRefSet(fmt.Sprintf("%s.%s", context, "enabled_VGPU_types"), enabledVGPUTypesValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeGPUGroupRef(context string, input interface{}) (GPUGroupRef, error) {
	var ref GPUGroupRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return GPUGroupRef(value), nil
}

func deserializeGPUGroupRefSet(context string, input interface{}) (slice []GPUGroupRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]GPUGroupRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeGPUGroupRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeGPUGroupRefToGPUGroupRecordMap(context string, input interface{}) (goMap map[GPUGroupRef]GPUGroupRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[GPUGroupRef]GPUGroupRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeGPUGroupRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeGPUGroupRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeHostCPURecord(context string, input interface{}) (record HostCPURecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	numberValue, ok := rpcStruct["number"]
	if ok && numberValue != nil {
		record.Number, err = deserializeInt(fmt.Sprintf("%s.%s", context, "number"), numberValue)
		if err != nil {
			return
		}
	}
	vendorValue, ok := rpcStruct["vendor"]
	if ok && vendorValue != nil {
		record.Vendor, err = deserializeString(fmt.Sprintf("%s.%s", context, "vendor"), vendorValue)
		if err != nil {
			return
		}
	}
	speedValue, ok := rpcStruct["speed"]
	if ok && speedValue != nil {
		record.Speed, err = deserializeInt(fmt.Sprintf("%s.%s", context, "speed"), speedValue)
		if err != nil {
			return
		}
	}
	modelnameValue, ok := rpcStruct["modelname"]
	if ok && modelnameValue != nil {
		record.Modelname, err = deserializeString(fmt.Sprintf("%s.%s", context, "modelname"), modelnameValue)
		if err != nil {
			return
		}
	}
	familyValue, ok := rpcStruct["family"]
	if ok && familyValue != nil {
		record.Family, err = deserializeInt(fmt.Sprintf("%s.%s", context, "family"), familyValue)
		if err != nil {
			return
		}
	}
	modelValue, ok := rpcStruct["model"]
	if ok && modelValue != nil {
		record.Model, err = deserializeInt(fmt.Sprintf("%s.%s", context, "model"), modelValue)
		if err != nil {
			return
		}
	}
	steppingValue, ok := rpcStruct["stepping"]
	if ok && steppingValue != nil {
		record.Stepping, err = deserializeString(fmt.Sprintf("%s.%s", context, "stepping"), steppingValue)
		if err != nil {
			return
		}
	}
	flagsValue, ok := rpcStruct["flags"]
	if ok && flagsValue != nil {
		record.Flags, err = deserializeString(fmt.Sprintf("%s.%s", context, "flags"), flagsValue)
		if err != nil {
			return
		}
	}
	featuresValue, ok := rpcStruct["features"]
	if ok && featuresValue != nil {
		record.Features, err = deserializeString(fmt.Sprintf("%s.%s", context, "features"), featuresValue)
		if err != nil {
			return
		}
	}
	utilisationValue, ok := rpcStruct["utilisation"]
	if ok && utilisationValue != nil {
		record.Utilisation, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "utilisation"), utilisationValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeHostCPURef(context string, input interface{}) (HostCPURef, error) {
	var ref HostCPURef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return HostCPURef(value), nil
}

func deserializeHostCPURefSet(context string, input interface{}) (slice []HostCPURef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]HostCPURef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeHostCPURef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeHostCPURefToHostCPURecordMap(context string, input interface{}) (goMap map[HostCPURef]HostCPURecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[HostCPURef]HostCPURecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeHostCPURef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeHostCPURecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeHostCrashdumpRecord(context string, input interface{}) (record HostCrashdumpRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	timestampValue, ok := rpcStruct["timestamp"]
	if ok && timestampValue != nil {
		record.Timestamp, err = deserializeTime(fmt.Sprintf("%s.%s", context, "timestamp"), timestampValue)
		if err != nil {
			return
		}
	}
	sizeValue, ok := rpcStruct["size"]
	if ok && sizeValue != nil {
		record.Size, err = deserializeInt(fmt.Sprintf("%s.%s", context, "size"), sizeValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeHostCrashdumpRef(context string, input interface{}) (HostCrashdumpRef, error) {
	var ref HostCrashdumpRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return HostCrashdumpRef(value), nil
}

func deserializeHostCrashdumpRefSet(context string, input interface{}) (slice []HostCrashdumpRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]HostCrashdumpRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeHostCrashdumpRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeHostCrashdumpRefToHostCrashdumpRecordMap(context string, input interface{}) (goMap map[HostCrashdumpRef]HostCrashdumpRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[HostCrashdumpRef]HostCrashdumpRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeHostCrashdumpRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeHostCrashdumpRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeHostDriverRecord(context string, input interface{}) (record HostDriverRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	nameValue, ok := rpcStruct["name"]
	if ok && nameValue != nil {
		record.Name, err = deserializeString(fmt.Sprintf("%s.%s", context, "name"), nameValue)
		if err != nil {
			return
		}
	}
	friendlyNameValue, ok := rpcStruct["friendly_name"]
	if ok && friendlyNameValue != nil {
		record.FriendlyName, err = deserializeString(fmt.Sprintf("%s.%s", context, "friendly_name"), friendlyNameValue)
		if err != nil {
			return
		}
	}
	variantsValue, ok := rpcStruct["variants"]
	if ok && variantsValue != nil {
		record.Variants, err = deserializeDriverVariantRefSet(fmt.Sprintf("%s.%s", context, "variants"), variantsValue)
		if err != nil {
			return
		}
	}
	activeVariantValue, ok := rpcStruct["active_variant"]
	if ok && activeVariantValue != nil {
		record.ActiveVariant, err = deserializeDriverVariantRef(fmt.Sprintf("%s.%s", context, "active_variant"), activeVariantValue)
		if err != nil {
			return
		}
	}
	selectedVariantValue, ok := rpcStruct["selected_variant"]
	if ok && selectedVariantValue != nil {
		record.SelectedVariant, err = deserializeDriverVariantRef(fmt.Sprintf("%s.%s", context, "selected_variant"), selectedVariantValue)
		if err != nil {
			return
		}
	}
	typeValue, ok := rpcStruct["type"]
	if ok && typeValue != nil {
		record.Type, err = deserializeString(fmt.Sprintf("%s.%s", context, "type"), typeValue)
		if err != nil {
			return
		}
	}
	descriptionValue, ok := rpcStruct["description"]
	if ok && descriptionValue != nil {
		record.Description, err = deserializeString(fmt.Sprintf("%s.%s", context, "description"), descriptionValue)
		if err != nil {
			return
		}
	}
	infoValue, ok := rpcStruct["info"]
	if ok && infoValue != nil {
		record.Info, err = deserializeString(fmt.Sprintf("%s.%s", context, "info"), infoValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeHostDriverRef(context string, input interface{}) (HostDriverRef, error) {
	var ref HostDriverRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return HostDriverRef(value), nil
}

func deserializeHostDriverRefSet(context string, input interface{}) (slice []HostDriverRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]HostDriverRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeHostDriverRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeHostDriverRefToHostDriverRecordMap(context string, input interface{}) (goMap map[HostDriverRef]HostDriverRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[HostDriverRef]HostDriverRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeHostDriverRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeHostDriverRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeHostMetricsRecord(context string, input interface{}) (record HostMetricsRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	memoryTotalValue, ok := rpcStruct["memory_total"]
	if ok && memoryTotalValue != nil {
		record.MemoryTotal, err = deserializeInt(fmt.Sprintf("%s.%s", context, "memory_total"), memoryTotalValue)
		if err != nil {
			return
		}
	}
	memoryFreeValue, ok := rpcStruct["memory_free"]
	if ok && memoryFreeValue != nil {
		record.MemoryFree, err = deserializeInt(fmt.Sprintf("%s.%s", context, "memory_free"), memoryFreeValue)
		if err != nil {
			return
		}
	}
	liveValue, ok := rpcStruct["live"]
	if ok && liveValue != nil {
		record.Live, err = deserializeBool(fmt.Sprintf("%s.%s", context, "live"), liveValue)
		if err != nil {
			return
		}
	}
	lastUpdatedValue, ok := rpcStruct["last_updated"]
	if ok && lastUpdatedValue != nil {
		record.LastUpdated, err = deserializeTime(fmt.Sprintf("%s.%s", context, "last_updated"), lastUpdatedValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeHostMetricsRef(context string, input interface{}) (HostMetricsRef, error) {
	var ref HostMetricsRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return HostMetricsRef(value), nil
}

func deserializeHostMetricsRefSet(context string, input interface{}) (slice []HostMetricsRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]HostMetricsRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeHostMetricsRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeHostMetricsRefToHostMetricsRecordMap(context string, input interface{}) (goMap map[HostMetricsRef]HostMetricsRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[HostMetricsRef]HostMetricsRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeHostMetricsRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeHostMetricsRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeHostPatchRecord(context string, input interface{}) (record HostPatchRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	versionValue, ok := rpcStruct["version"]
	if ok && versionValue != nil {
		record.Version, err = deserializeString(fmt.Sprintf("%s.%s", context, "version"), versionValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	appliedValue, ok := rpcStruct["applied"]
	if ok && appliedValue != nil {
		record.Applied, err = deserializeBool(fmt.Sprintf("%s.%s", context, "applied"), appliedValue)
		if err != nil {
			return
		}
	}
	timestampAppliedValue, ok := rpcStruct["timestamp_applied"]
	if ok && timestampAppliedValue != nil {
		record.TimestampApplied, err = deserializeTime(fmt.Sprintf("%s.%s", context, "timestamp_applied"), timestampAppliedValue)
		if err != nil {
			return
		}
	}
	sizeValue, ok := rpcStruct["size"]
	if ok && sizeValue != nil {
		record.Size, err = deserializeInt(fmt.Sprintf("%s.%s", context, "size"), sizeValue)
		if err != nil {
			return
		}
	}
	poolPatchValue, ok := rpcStruct["pool_patch"]
	if ok && poolPatchValue != nil {
		record.PoolPatch, err = deserializePoolPatchRef(fmt.Sprintf("%s.%s", context, "pool_patch"), poolPatchValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeHostPatchRef(context string, input interface{}) (HostPatchRef, error) {
	var ref HostPatchRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return HostPatchRef(value), nil
}

func deserializeHostPatchRefSet(context string, input interface{}) (slice []HostPatchRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]HostPatchRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeHostPatchRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeHostPatchRefToHostPatchRecordMap(context string, input interface{}) (goMap map[HostPatchRef]HostPatchRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[HostPatchRef]HostPatchRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeHostPatchRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeHostPatchRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeHostRecord(context string, input interface{}) (record HostRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	memoryOverheadValue, ok := rpcStruct["memory_overhead"]
	if ok && memoryOverheadValue != nil {
		record.MemoryOverhead, err = deserializeInt(fmt.Sprintf("%s.%s", context, "memory_overhead"), memoryOverheadValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumHostAllowedOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumHostAllowedOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	aPIVersionMajorValue, ok := rpcStruct["API_version_major"]
	if ok && aPIVersionMajorValue != nil {
		record.APIVersionMajor, err = deserializeInt(fmt.Sprintf("%s.%s", context, "API_version_major"), aPIVersionMajorValue)
		if err != nil {
			return
		}
	}
	aPIVersionMinorValue, ok := rpcStruct["API_version_minor"]
	if ok && aPIVersionMinorValue != nil {
		record.APIVersionMinor, err = deserializeInt(fmt.Sprintf("%s.%s", context, "API_version_minor"), aPIVersionMinorValue)
		if err != nil {
			return
		}
	}
	aPIVersionVendorValue, ok := rpcStruct["API_version_vendor"]
	if ok && aPIVersionVendorValue != nil {
		record.APIVersionVendor, err = deserializeString(fmt.Sprintf("%s.%s", context, "API_version_vendor"), aPIVersionVendorValue)
		if err != nil {
			return
		}
	}
	aPIVersionVendorImplementationValue, ok := rpcStruct["API_version_vendor_implementation"]
	if ok && aPIVersionVendorImplementationValue != nil {
		record.APIVersionVendorImplementation, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "API_version_vendor_implementation"), aPIVersionVendorImplementationValue)
		if err != nil {
			return
		}
	}
	enabledValue, ok := rpcStruct["enabled"]
	if ok && enabledValue != nil {
		record.Enabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "enabled"), enabledValue)
		if err != nil {
			return
		}
	}
	softwareVersionValue, ok := rpcStruct["software_version"]
	if ok && softwareVersionValue != nil {
		record.SoftwareVersion, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "software_version"), softwareVersionValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	capabilitiesValue, ok := rpcStruct["capabilities"]
	if ok && capabilitiesValue != nil {
		record.Capabilities, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "capabilities"), capabilitiesValue)
		if err != nil {
			return
		}
	}
	cPUConfigurationValue, ok := rpcStruct["cpu_configuration"]
	if ok && cPUConfigurationValue != nil {
		record.CPUConfiguration, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "cpu_configuration"), cPUConfigurationValue)
		if err != nil {
			return
		}
	}
	schedPolicyValue, ok := rpcStruct["sched_policy"]
	if ok && schedPolicyValue != nil {
		record.SchedPolicy, err = deserializeString(fmt.Sprintf("%s.%s", context, "sched_policy"), schedPolicyValue)
		if err != nil {
			return
		}
	}
	supportedBootloadersValue, ok := rpcStruct["supported_bootloaders"]
	if ok && supportedBootloadersValue != nil {
		record.SupportedBootloaders, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "supported_bootloaders"), supportedBootloadersValue)
		if err != nil {
			return
		}
	}
	residentVMsValue, ok := rpcStruct["resident_VMs"]
	if ok && residentVMsValue != nil {
		record.ResidentVMs, err = deserializeVMRefSet(fmt.Sprintf("%s.%s", context, "resident_VMs"), residentVMsValue)
		if err != nil {
			return
		}
	}
	loggingValue, ok := rpcStruct["logging"]
	if ok && loggingValue != nil {
		record.Logging, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "logging"), loggingValue)
		if err != nil {
			return
		}
	}
	pIFsValue, ok := rpcStruct["PIFs"]
	if ok && pIFsValue != nil {
		record.PIFs, err = deserializePIFRefSet(fmt.Sprintf("%s.%s", context, "PIFs"), pIFsValue)
		if err != nil {
			return
		}
	}
	suspendImageSrValue, ok := rpcStruct["suspend_image_sr"]
	if ok && suspendImageSrValue != nil {
		record.SuspendImageSr, err = deserializeSRRef(fmt.Sprintf("%s.%s", context, "suspend_image_sr"), suspendImageSrValue)
		if err != nil {
			return
		}
	}
	crashDumpSrValue, ok := rpcStruct["crash_dump_sr"]
	if ok && crashDumpSrValue != nil {
		record.CrashDumpSr, err = deserializeSRRef(fmt.Sprintf("%s.%s", context, "crash_dump_sr"), crashDumpSrValue)
		if err != nil {
			return
		}
	}
	crashdumpsValue, ok := rpcStruct["crashdumps"]
	if ok && crashdumpsValue != nil {
		record.Crashdumps, err = deserializeHostCrashdumpRefSet(fmt.Sprintf("%s.%s", context, "crashdumps"), crashdumpsValue)
		if err != nil {
			return
		}
	}
	patchesValue, ok := rpcStruct["patches"]
	if ok && patchesValue != nil {
		record.Patches, err = deserializeHostPatchRefSet(fmt.Sprintf("%s.%s", context, "patches"), patchesValue)
		if err != nil {
			return
		}
	}
	updatesValue, ok := rpcStruct["updates"]
	if ok && updatesValue != nil {
		record.Updates, err = deserializePoolUpdateRefSet(fmt.Sprintf("%s.%s", context, "updates"), updatesValue)
		if err != nil {
			return
		}
	}
	pBDsValue, ok := rpcStruct["PBDs"]
	if ok && pBDsValue != nil {
		record.PBDs, err = deserializePBDRefSet(fmt.Sprintf("%s.%s", context, "PBDs"), pBDsValue)
		if err != nil {
			return
		}
	}
	hostCPUsValue, ok := rpcStruct["host_CPUs"]
	if ok && hostCPUsValue != nil {
		record.HostCPUs, err = deserializeHostCPURefSet(fmt.Sprintf("%s.%s", context, "host_CPUs"), hostCPUsValue)
		if err != nil {
			return
		}
	}
	cPUInfoValue, ok := rpcStruct["cpu_info"]
	if ok && cPUInfoValue != nil {
		record.CPUInfo, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "cpu_info"), cPUInfoValue)
		if err != nil {
			return
		}
	}
	hostnameValue, ok := rpcStruct["hostname"]
	if ok && hostnameValue != nil {
		record.Hostname, err = deserializeString(fmt.Sprintf("%s.%s", context, "hostname"), hostnameValue)
		if err != nil {
			return
		}
	}
	addressValue, ok := rpcStruct["address"]
	if ok && addressValue != nil {
		record.Address, err = deserializeString(fmt.Sprintf("%s.%s", context, "address"), addressValue)
		if err != nil {
			return
		}
	}
	metricsValue, ok := rpcStruct["metrics"]
	if ok && metricsValue != nil {
		record.Metrics, err = deserializeHostMetricsRef(fmt.Sprintf("%s.%s", context, "metrics"), metricsValue)
		if err != nil {
			return
		}
	}
	licenseParamsValue, ok := rpcStruct["license_params"]
	if ok && licenseParamsValue != nil {
		record.LicenseParams, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "license_params"), licenseParamsValue)
		if err != nil {
			return
		}
	}
	haStatefilesValue, ok := rpcStruct["ha_statefiles"]
	if ok && haStatefilesValue != nil {
		record.HaStatefiles, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "ha_statefiles"), haStatefilesValue)
		if err != nil {
			return
		}
	}
	haNetworkPeersValue, ok := rpcStruct["ha_network_peers"]
	if ok && haNetworkPeersValue != nil {
		record.HaNetworkPeers, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "ha_network_peers"), haNetworkPeersValue)
		if err != nil {
			return
		}
	}
	blobsValue, ok := rpcStruct["blobs"]
	if ok && blobsValue != nil {
		record.Blobs, err = deserializeStringToBlobRefMap(fmt.Sprintf("%s.%s", context, "blobs"), blobsValue)
		if err != nil {
			return
		}
	}
	tagsValue, ok := rpcStruct["tags"]
	if ok && tagsValue != nil {
		record.Tags, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "tags"), tagsValue)
		if err != nil {
			return
		}
	}
	externalAuthTypeValue, ok := rpcStruct["external_auth_type"]
	if ok && externalAuthTypeValue != nil {
		record.ExternalAuthType, err = deserializeString(fmt.Sprintf("%s.%s", context, "external_auth_type"), externalAuthTypeValue)
		if err != nil {
			return
		}
	}
	externalAuthServiceNameValue, ok := rpcStruct["external_auth_service_name"]
	if ok && externalAuthServiceNameValue != nil {
		record.ExternalAuthServiceName, err = deserializeString(fmt.Sprintf("%s.%s", context, "external_auth_service_name"), externalAuthServiceNameValue)
		if err != nil {
			return
		}
	}
	externalAuthConfigurationValue, ok := rpcStruct["external_auth_configuration"]
	if ok && externalAuthConfigurationValue != nil {
		record.ExternalAuthConfiguration, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "external_auth_configuration"), externalAuthConfigurationValue)
		if err != nil {
			return
		}
	}
	editionValue, ok := rpcStruct["edition"]
	if ok && editionValue != nil {
		record.Edition, err = deserializeString(fmt.Sprintf("%s.%s", context, "edition"), editionValue)
		if err != nil {
			return
		}
	}
	licenseServerValue, ok := rpcStruct["license_server"]
	if ok && licenseServerValue != nil {
		record.LicenseServer, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "license_server"), licenseServerValue)
		if err != nil {
			return
		}
	}
	biosStringsValue, ok := rpcStruct["bios_strings"]
	if ok && biosStringsValue != nil {
		record.BiosStrings, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "bios_strings"), biosStringsValue)
		if err != nil {
			return
		}
	}
	powerOnModeValue, ok := rpcStruct["power_on_mode"]
	if ok && powerOnModeValue != nil {
		record.PowerOnMode, err = deserializeString(fmt.Sprintf("%s.%s", context, "power_on_mode"), powerOnModeValue)
		if err != nil {
			return
		}
	}
	powerOnConfigValue, ok := rpcStruct["power_on_config"]
	if ok && powerOnConfigValue != nil {
		record.PowerOnConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "power_on_config"), powerOnConfigValue)
		if err != nil {
			return
		}
	}
	localCacheSrValue, ok := rpcStruct["local_cache_sr"]
	if ok && localCacheSrValue != nil {
		record.LocalCacheSr, err = deserializeSRRef(fmt.Sprintf("%s.%s", context, "local_cache_sr"), localCacheSrValue)
		if err != nil {
			return
		}
	}
	chipsetInfoValue, ok := rpcStruct["chipset_info"]
	if ok && chipsetInfoValue != nil {
		record.ChipsetInfo, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "chipset_info"), chipsetInfoValue)
		if err != nil {
			return
		}
	}
	pCIsValue, ok := rpcStruct["PCIs"]
	if ok && pCIsValue != nil {
		record.PCIs, err = deserializePCIRefSet(fmt.Sprintf("%s.%s", context, "PCIs"), pCIsValue)
		if err != nil {
			return
		}
	}
	pGPUsValue, ok := rpcStruct["PGPUs"]
	if ok && pGPUsValue != nil {
		record.PGPUs, err = deserializePGPURefSet(fmt.Sprintf("%s.%s", context, "PGPUs"), pGPUsValue)
		if err != nil {
			return
		}
	}
	pUSBsValue, ok := rpcStruct["PUSBs"]
	if ok && pUSBsValue != nil {
		record.PUSBs, err = deserializePUSBRefSet(fmt.Sprintf("%s.%s", context, "PUSBs"), pUSBsValue)
		if err != nil {
			return
		}
	}
	sslLegacyValue, ok := rpcStruct["ssl_legacy"]
	if ok && sslLegacyValue != nil {
		record.SslLegacy, err = deserializeBool(fmt.Sprintf("%s.%s", context, "ssl_legacy"), sslLegacyValue)
		if err != nil {
			return
		}
	}
	guestVCPUsParamsValue, ok := rpcStruct["guest_VCPUs_params"]
	if ok && guestVCPUsParamsValue != nil {
		record.GuestVCPUsParams, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "guest_VCPUs_params"), guestVCPUsParamsValue)
		if err != nil {
			return
		}
	}
	displayValue, ok := rpcStruct["display"]
	if ok && displayValue != nil {
		record.Display, err = deserializeEnumHostDisplay(fmt.Sprintf("%s.%s", context, "display"), displayValue)
		if err != nil {
			return
		}
	}
	virtualHardwarePlatformVersionsValue, ok := rpcStruct["virtual_hardware_platform_versions"]
	if ok && virtualHardwarePlatformVersionsValue != nil {
		record.VirtualHardwarePlatformVersions, err = deserializeIntSet(fmt.Sprintf("%s.%s", context, "virtual_hardware_platform_versions"), virtualHardwarePlatformVersionsValue)
		if err != nil {
			return
		}
	}
	controlDomainValue, ok := rpcStruct["control_domain"]
	if ok && controlDomainValue != nil {
		record.ControlDomain, err = deserializeVMRef(fmt.Sprintf("%s.%s", context, "control_domain"), controlDomainValue)
		if err != nil {
			return
		}
	}
	updatesRequiringRebootValue, ok := rpcStruct["updates_requiring_reboot"]
	if ok && updatesRequiringRebootValue != nil {
		record.UpdatesRequiringReboot, err = deserializePoolUpdateRefSet(fmt.Sprintf("%s.%s", context, "updates_requiring_reboot"), updatesRequiringRebootValue)
		if err != nil {
			return
		}
	}
	featuresValue, ok := rpcStruct["features"]
	if ok && featuresValue != nil {
		record.Features, err = deserializeFeatureRefSet(fmt.Sprintf("%s.%s", context, "features"), featuresValue)
		if err != nil {
			return
		}
	}
	iscsiIqnValue, ok := rpcStruct["iscsi_iqn"]
	if ok && iscsiIqnValue != nil {
		record.IscsiIqn, err = deserializeString(fmt.Sprintf("%s.%s", context, "iscsi_iqn"), iscsiIqnValue)
		if err != nil {
			return
		}
	}
	multipathingValue, ok := rpcStruct["multipathing"]
	if ok && multipathingValue != nil {
		record.Multipathing, err = deserializeBool(fmt.Sprintf("%s.%s", context, "multipathing"), multipathingValue)
		if err != nil {
			return
		}
	}
	uefiCertificatesValue, ok := rpcStruct["uefi_certificates"]
	if ok && uefiCertificatesValue != nil {
		record.UefiCertificates, err = deserializeString(fmt.Sprintf("%s.%s", context, "uefi_certificates"), uefiCertificatesValue)
		if err != nil {
			return
		}
	}
	certificatesValue, ok := rpcStruct["certificates"]
	if ok && certificatesValue != nil {
		record.Certificates, err = deserializeCertificateRefSet(fmt.Sprintf("%s.%s", context, "certificates"), certificatesValue)
		if err != nil {
			return
		}
	}
	editionsValue, ok := rpcStruct["editions"]
	if ok && editionsValue != nil {
		record.Editions, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "editions"), editionsValue)
		if err != nil {
			return
		}
	}
	pendingGuidancesValue, ok := rpcStruct["pending_guidances"]
	if ok && pendingGuidancesValue != nil {
		record.PendingGuidances, err = deserializeEnumUpdateGuidancesSet(fmt.Sprintf("%s.%s", context, "pending_guidances"), pendingGuidancesValue)
		if err != nil {
			return
		}
	}
	tLSVerificationEnabledValue, ok := rpcStruct["tls_verification_enabled"]
	if ok && tLSVerificationEnabledValue != nil {
		record.TLSVerificationEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "tls_verification_enabled"), tLSVerificationEnabledValue)
		if err != nil {
			return
		}
	}
	lastSoftwareUpdateValue, ok := rpcStruct["last_software_update"]
	if ok && lastSoftwareUpdateValue != nil {
		record.LastSoftwareUpdate, err = deserializeTime(fmt.Sprintf("%s.%s", context, "last_software_update"), lastSoftwareUpdateValue)
		if err != nil {
			return
		}
	}
	hTTPSOnlyValue, ok := rpcStruct["https_only"]
	if ok && hTTPSOnlyValue != nil {
		record.HTTPSOnly, err = deserializeBool(fmt.Sprintf("%s.%s", context, "https_only"), hTTPSOnlyValue)
		if err != nil {
			return
		}
	}
	latestSyncedUpdatesAppliedValue, ok := rpcStruct["latest_synced_updates_applied"]
	if ok && latestSyncedUpdatesAppliedValue != nil {
		record.LatestSyncedUpdatesApplied, err = deserializeEnumLatestSyncedUpdatesAppliedState(fmt.Sprintf("%s.%s", context, "latest_synced_updates_applied"), latestSyncedUpdatesAppliedValue)
		if err != nil {
			return
		}
	}
	numaAffinityPolicyValue, ok := rpcStruct["numa_affinity_policy"]
	if ok && numaAffinityPolicyValue != nil {
		record.NumaAffinityPolicy, err = deserializeEnumHostNumaAffinityPolicy(fmt.Sprintf("%s.%s", context, "numa_affinity_policy"), numaAffinityPolicyValue)
		if err != nil {
			return
		}
	}
	pendingGuidancesRecommendedValue, ok := rpcStruct["pending_guidances_recommended"]
	if ok && pendingGuidancesRecommendedValue != nil {
		record.PendingGuidancesRecommended, err = deserializeEnumUpdateGuidancesSet(fmt.Sprintf("%s.%s", context, "pending_guidances_recommended"), pendingGuidancesRecommendedValue)
		if err != nil {
			return
		}
	}
	pendingGuidancesFullValue, ok := rpcStruct["pending_guidances_full"]
	if ok && pendingGuidancesFullValue != nil {
		record.PendingGuidancesFull, err = deserializeEnumUpdateGuidancesSet(fmt.Sprintf("%s.%s", context, "pending_guidances_full"), pendingGuidancesFullValue)
		if err != nil {
			return
		}
	}
	lastUpdateHashValue, ok := rpcStruct["last_update_hash"]
	if ok && lastUpdateHashValue != nil {
		record.LastUpdateHash, err = deserializeString(fmt.Sprintf("%s.%s", context, "last_update_hash"), lastUpdateHashValue)
		if err != nil {
			return
		}
	}
	sSHEnabledValue, ok := rpcStruct["ssh_enabled"]
	if ok && sSHEnabledValue != nil {
		record.SSHEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "ssh_enabled"), sSHEnabledValue)
		if err != nil {
			return
		}
	}
	sSHEnabledTimeoutValue, ok := rpcStruct["ssh_enabled_timeout"]
	if ok && sSHEnabledTimeoutValue != nil {
		record.SSHEnabledTimeout, err = deserializeInt(fmt.Sprintf("%s.%s", context, "ssh_enabled_timeout"), sSHEnabledTimeoutValue)
		if err != nil {
			return
		}
	}
	sSHExpiryValue, ok := rpcStruct["ssh_expiry"]
	if ok && sSHExpiryValue != nil {
		record.SSHExpiry, err = deserializeTime(fmt.Sprintf("%s.%s", context, "ssh_expiry"), sSHExpiryValue)
		if err != nil {
			return
		}
	}
	consoleIdleTimeoutValue, ok := rpcStruct["console_idle_timeout"]
	if ok && consoleIdleTimeoutValue != nil {
		record.ConsoleIdleTimeout, err = deserializeInt(fmt.Sprintf("%s.%s", context, "console_idle_timeout"), consoleIdleTimeoutValue)
		if err != nil {
			return
		}
	}
	sSHAutoModeValue, ok := rpcStruct["ssh_auto_mode"]
	if ok && sSHAutoModeValue != nil {
		record.SSHAutoMode, err = deserializeBool(fmt.Sprintf("%s.%s", context, "ssh_auto_mode"), sSHAutoModeValue)
		if err != nil {
			return
		}
	}
	maxCstateValue, ok := rpcStruct["max_cstate"]
	if ok && maxCstateValue != nil {
		record.MaxCstate, err = deserializeString(fmt.Sprintf("%s.%s", context, "max_cstate"), maxCstateValue)
		if err != nil {
			return
		}
	}
	secureBootValue, ok := rpcStruct["secure_boot"]
	if ok && secureBootValue != nil {
		record.SecureBoot, err = deserializeBool(fmt.Sprintf("%s.%s", context, "secure_boot"), secureBootValue)
		if err != nil {
			return
		}
	}
	ntpModeValue, ok := rpcStruct["ntp_mode"]
	if ok && ntpModeValue != nil {
		record.NtpMode, err = deserializeEnumHostNtpMode(fmt.Sprintf("%s.%s", context, "ntp_mode"), ntpModeValue)
		if err != nil {
			return
		}
	}
	ntpCustomServersValue, ok := rpcStruct["ntp_custom_servers"]
	if ok && ntpCustomServersValue != nil {
		record.NtpCustomServers, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "ntp_custom_servers"), ntpCustomServersValue)
		if err != nil {
			return
		}
	}
	timezoneValue, ok := rpcStruct["timezone"]
	if ok && timezoneValue != nil {
		record.Timezone, err = deserializeString(fmt.Sprintf("%s.%s", context, "timezone"), timezoneValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeHostRef(context string, input interface{}) (HostRef, error) {
	var ref HostRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return HostRef(value), nil
}

func deserializeHostRefSet(context string, input interface{}) (slice []HostRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]HostRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeHostRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeHostRefToHostRecordMap(context string, input interface{}) (goMap map[HostRef]HostRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[HostRef]HostRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeHostRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeHostRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeHostRefToStringSetMap(context string, input interface{}) (goMap map[HostRef][]string, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[HostRef][]string, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeHostRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeStringSet(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeInt(context string, input interface{}) (value int, err error) {
	_ = context
	if input == nil {
		return
	}
	strValue := fmt.Sprintf("%v", input)
	value, err = strconv.Atoi(strValue)
	if err != nil {
		floatValue, err1 := strconv.ParseFloat(strValue, 64)
		if err1 == nil {
			return int(floatValue), nil
		}
	}
	return
}

func deserializeIntSet(context string, input interface{}) (slice []int, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]int, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeInt(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeIntToFloatMap(context string, input interface{}) (goMap map[int]float64, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[int]float64, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeInt(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeFloat(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeIntToIntMap(context string, input interface{}) (goMap map[int]int, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[int]int, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeInt(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeInt(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeIntToStringSetMap(context string, input interface{}) (goMap map[int][]string, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[int][]string, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeInt(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeStringSet(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeLVHDRecord(context string, input interface{}) (record LVHDRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeLVHDRef(context string, input interface{}) (LVHDRef, error) {
	var ref LVHDRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return LVHDRef(value), nil
}

func deserializeMessageRecord(context string, input interface{}) (record MessageRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameValue, ok := rpcStruct["name"]
	if ok && nameValue != nil {
		record.Name, err = deserializeString(fmt.Sprintf("%s.%s", context, "name"), nameValue)
		if err != nil {
			return
		}
	}
	priorityValue, ok := rpcStruct["priority"]
	if ok && priorityValue != nil {
		record.Priority, err = deserializeInt(fmt.Sprintf("%s.%s", context, "priority"), priorityValue)
		if err != nil {
			return
		}
	}
	clsValue, ok := rpcStruct["cls"]
	if ok && clsValue != nil {
		record.Cls, err = deserializeEnumCls(fmt.Sprintf("%s.%s", context, "cls"), clsValue)
		if err != nil {
			return
		}
	}
	objUUIDValue, ok := rpcStruct["obj_uuid"]
	if ok && objUUIDValue != nil {
		record.ObjUUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "obj_uuid"), objUUIDValue)
		if err != nil {
			return
		}
	}
	timestampValue, ok := rpcStruct["timestamp"]
	if ok && timestampValue != nil {
		record.Timestamp, err = deserializeTime(fmt.Sprintf("%s.%s", context, "timestamp"), timestampValue)
		if err != nil {
			return
		}
	}
	bodyValue, ok := rpcStruct["body"]
	if ok && bodyValue != nil {
		record.Body, err = deserializeString(fmt.Sprintf("%s.%s", context, "body"), bodyValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeMessageRef(context string, input interface{}) (MessageRef, error) {
	var ref MessageRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return MessageRef(value), nil
}

func deserializeMessageRefSet(context string, input interface{}) (slice []MessageRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]MessageRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeMessageRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeMessageRefToMessageRecordMap(context string, input interface{}) (goMap map[MessageRef]MessageRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[MessageRef]MessageRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeMessageRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeMessageRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeNetworkRecord(context string, input interface{}) (record NetworkRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumNetworkOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumNetworkOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	vIFsValue, ok := rpcStruct["VIFs"]
	if ok && vIFsValue != nil {
		record.VIFs, err = deserializeVIFRefSet(fmt.Sprintf("%s.%s", context, "VIFs"), vIFsValue)
		if err != nil {
			return
		}
	}
	pIFsValue, ok := rpcStruct["PIFs"]
	if ok && pIFsValue != nil {
		record.PIFs, err = deserializePIFRefSet(fmt.Sprintf("%s.%s", context, "PIFs"), pIFsValue)
		if err != nil {
			return
		}
	}
	mTUValue, ok := rpcStruct["MTU"]
	if ok && mTUValue != nil {
		record.MTU, err = deserializeInt(fmt.Sprintf("%s.%s", context, "MTU"), mTUValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	bridgeValue, ok := rpcStruct["bridge"]
	if ok && bridgeValue != nil {
		record.Bridge, err = deserializeString(fmt.Sprintf("%s.%s", context, "bridge"), bridgeValue)
		if err != nil {
			return
		}
	}
	managedValue, ok := rpcStruct["managed"]
	if ok && managedValue != nil {
		record.Managed, err = deserializeBool(fmt.Sprintf("%s.%s", context, "managed"), managedValue)
		if err != nil {
			return
		}
	}
	blobsValue, ok := rpcStruct["blobs"]
	if ok && blobsValue != nil {
		record.Blobs, err = deserializeStringToBlobRefMap(fmt.Sprintf("%s.%s", context, "blobs"), blobsValue)
		if err != nil {
			return
		}
	}
	tagsValue, ok := rpcStruct["tags"]
	if ok && tagsValue != nil {
		record.Tags, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "tags"), tagsValue)
		if err != nil {
			return
		}
	}
	defaultLockingModeValue, ok := rpcStruct["default_locking_mode"]
	if ok && defaultLockingModeValue != nil {
		record.DefaultLockingMode, err = deserializeEnumNetworkDefaultLockingMode(fmt.Sprintf("%s.%s", context, "default_locking_mode"), defaultLockingModeValue)
		if err != nil {
			return
		}
	}
	assignedIpsValue, ok := rpcStruct["assigned_ips"]
	if ok && assignedIpsValue != nil {
		record.AssignedIps, err = deserializeVIFRefToStringMap(fmt.Sprintf("%s.%s", context, "assigned_ips"), assignedIpsValue)
		if err != nil {
			return
		}
	}
	purposeValue, ok := rpcStruct["purpose"]
	if ok && purposeValue != nil {
		record.Purpose, err = deserializeEnumNetworkPurposeSet(fmt.Sprintf("%s.%s", context, "purpose"), purposeValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeNetworkRef(context string, input interface{}) (NetworkRef, error) {
	var ref NetworkRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return NetworkRef(value), nil
}

func deserializeNetworkRefSet(context string, input interface{}) (slice []NetworkRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]NetworkRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeNetworkRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeNetworkRefToNetworkRecordMap(context string, input interface{}) (goMap map[NetworkRef]NetworkRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[NetworkRef]NetworkRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeNetworkRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeNetworkRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeNetworkSriovRecord(context string, input interface{}) (record NetworkSriovRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	physicalPIFValue, ok := rpcStruct["physical_PIF"]
	if ok && physicalPIFValue != nil {
		record.PhysicalPIF, err = deserializePIFRef(fmt.Sprintf("%s.%s", context, "physical_PIF"), physicalPIFValue)
		if err != nil {
			return
		}
	}
	logicalPIFValue, ok := rpcStruct["logical_PIF"]
	if ok && logicalPIFValue != nil {
		record.LogicalPIF, err = deserializePIFRef(fmt.Sprintf("%s.%s", context, "logical_PIF"), logicalPIFValue)
		if err != nil {
			return
		}
	}
	requiresRebootValue, ok := rpcStruct["requires_reboot"]
	if ok && requiresRebootValue != nil {
		record.RequiresReboot, err = deserializeBool(fmt.Sprintf("%s.%s", context, "requires_reboot"), requiresRebootValue)
		if err != nil {
			return
		}
	}
	configurationModeValue, ok := rpcStruct["configuration_mode"]
	if ok && configurationModeValue != nil {
		record.ConfigurationMode, err = deserializeEnumSriovConfigurationMode(fmt.Sprintf("%s.%s", context, "configuration_mode"), configurationModeValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeNetworkSriovRef(context string, input interface{}) (NetworkSriovRef, error) {
	var ref NetworkSriovRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return NetworkSriovRef(value), nil
}

func deserializeNetworkSriovRefSet(context string, input interface{}) (slice []NetworkSriovRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]NetworkSriovRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeNetworkSriovRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeNetworkSriovRefToNetworkSriovRecordMap(context string, input interface{}) (goMap map[NetworkSriovRef]NetworkSriovRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[NetworkSriovRef]NetworkSriovRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeNetworkSriovRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeNetworkSriovRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeObserverRecord(context string, input interface{}) (record ObserverRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	hostsValue, ok := rpcStruct["hosts"]
	if ok && hostsValue != nil {
		record.Hosts, err = deserializeHostRefSet(fmt.Sprintf("%s.%s", context, "hosts"), hostsValue)
		if err != nil {
			return
		}
	}
	attributesValue, ok := rpcStruct["attributes"]
	if ok && attributesValue != nil {
		record.Attributes, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "attributes"), attributesValue)
		if err != nil {
			return
		}
	}
	endpointsValue, ok := rpcStruct["endpoints"]
	if ok && endpointsValue != nil {
		record.Endpoints, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "endpoints"), endpointsValue)
		if err != nil {
			return
		}
	}
	componentsValue, ok := rpcStruct["components"]
	if ok && componentsValue != nil {
		record.Components, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "components"), componentsValue)
		if err != nil {
			return
		}
	}
	enabledValue, ok := rpcStruct["enabled"]
	if ok && enabledValue != nil {
		record.Enabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "enabled"), enabledValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeObserverRef(context string, input interface{}) (ObserverRef, error) {
	var ref ObserverRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return ObserverRef(value), nil
}

func deserializeObserverRefSet(context string, input interface{}) (slice []ObserverRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]ObserverRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeObserverRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeObserverRefToObserverRecordMap(context string, input interface{}) (goMap map[ObserverRef]ObserverRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[ObserverRef]ObserverRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeObserverRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeObserverRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeOptionSrStatRecord(context string, input interface{}) (option OptionSrStatRecord, err error) {
	if input == nil {
		return
	}
	value, err := deserializeSrStatRecord(context, input)
	if err != nil {
		return
	}
	option = &value
	return
}

func deserializeOptionString(context string, input interface{}) (option OptionString, err error) {
	if input == nil {
		return
	}
	value, err := deserializeString(context, input)
	if err != nil {
		return
	}
	option = &value
	return
}

func deserializePBDRecord(context string, input interface{}) (record PBDRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	sRValue, ok := rpcStruct["SR"]
	if ok && sRValue != nil {
		record.SR, err = deserializeSRRef(fmt.Sprintf("%s.%s", context, "SR"), sRValue)
		if err != nil {
			return
		}
	}
	deviceConfigValue, ok := rpcStruct["device_config"]
	if ok && deviceConfigValue != nil {
		record.DeviceConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "device_config"), deviceConfigValue)
		if err != nil {
			return
		}
	}
	currentlyAttachedValue, ok := rpcStruct["currently_attached"]
	if ok && currentlyAttachedValue != nil {
		record.CurrentlyAttached, err = deserializeBool(fmt.Sprintf("%s.%s", context, "currently_attached"), currentlyAttachedValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePBDRef(context string, input interface{}) (PBDRef, error) {
	var ref PBDRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PBDRef(value), nil
}

func deserializePBDRefSet(context string, input interface{}) (slice []PBDRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PBDRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePBDRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePBDRefToPBDRecordMap(context string, input interface{}) (goMap map[PBDRef]PBDRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PBDRef]PBDRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePBDRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePBDRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializePCIRecord(context string, input interface{}) (record PCIRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	classIDValue, ok := rpcStruct["class_id"]
	if ok && classIDValue != nil {
		record.ClassID, err = deserializeString(fmt.Sprintf("%s.%s", context, "class_id"), classIDValue)
		if err != nil {
			return
		}
	}
	classNameValue, ok := rpcStruct["class_name"]
	if ok && classNameValue != nil {
		record.ClassName, err = deserializeString(fmt.Sprintf("%s.%s", context, "class_name"), classNameValue)
		if err != nil {
			return
		}
	}
	vendorIDValue, ok := rpcStruct["vendor_id"]
	if ok && vendorIDValue != nil {
		record.VendorID, err = deserializeString(fmt.Sprintf("%s.%s", context, "vendor_id"), vendorIDValue)
		if err != nil {
			return
		}
	}
	vendorNameValue, ok := rpcStruct["vendor_name"]
	if ok && vendorNameValue != nil {
		record.VendorName, err = deserializeString(fmt.Sprintf("%s.%s", context, "vendor_name"), vendorNameValue)
		if err != nil {
			return
		}
	}
	deviceIDValue, ok := rpcStruct["device_id"]
	if ok && deviceIDValue != nil {
		record.DeviceID, err = deserializeString(fmt.Sprintf("%s.%s", context, "device_id"), deviceIDValue)
		if err != nil {
			return
		}
	}
	deviceNameValue, ok := rpcStruct["device_name"]
	if ok && deviceNameValue != nil {
		record.DeviceName, err = deserializeString(fmt.Sprintf("%s.%s", context, "device_name"), deviceNameValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	pciIDValue, ok := rpcStruct["pci_id"]
	if ok && pciIDValue != nil {
		record.PciID, err = deserializeString(fmt.Sprintf("%s.%s", context, "pci_id"), pciIDValue)
		if err != nil {
			return
		}
	}
	dependenciesValue, ok := rpcStruct["dependencies"]
	if ok && dependenciesValue != nil {
		record.Dependencies, err = deserializePCIRefSet(fmt.Sprintf("%s.%s", context, "dependencies"), dependenciesValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	subsystemVendorIDValue, ok := rpcStruct["subsystem_vendor_id"]
	if ok && subsystemVendorIDValue != nil {
		record.SubsystemVendorID, err = deserializeString(fmt.Sprintf("%s.%s", context, "subsystem_vendor_id"), subsystemVendorIDValue)
		if err != nil {
			return
		}
	}
	subsystemVendorNameValue, ok := rpcStruct["subsystem_vendor_name"]
	if ok && subsystemVendorNameValue != nil {
		record.SubsystemVendorName, err = deserializeString(fmt.Sprintf("%s.%s", context, "subsystem_vendor_name"), subsystemVendorNameValue)
		if err != nil {
			return
		}
	}
	subsystemDeviceIDValue, ok := rpcStruct["subsystem_device_id"]
	if ok && subsystemDeviceIDValue != nil {
		record.SubsystemDeviceID, err = deserializeString(fmt.Sprintf("%s.%s", context, "subsystem_device_id"), subsystemDeviceIDValue)
		if err != nil {
			return
		}
	}
	subsystemDeviceNameValue, ok := rpcStruct["subsystem_device_name"]
	if ok && subsystemDeviceNameValue != nil {
		record.SubsystemDeviceName, err = deserializeString(fmt.Sprintf("%s.%s", context, "subsystem_device_name"), subsystemDeviceNameValue)
		if err != nil {
			return
		}
	}
	driverNameValue, ok := rpcStruct["driver_name"]
	if ok && driverNameValue != nil {
		record.DriverName, err = deserializeString(fmt.Sprintf("%s.%s", context, "driver_name"), driverNameValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePCIRef(context string, input interface{}) (PCIRef, error) {
	var ref PCIRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PCIRef(value), nil
}

func deserializePCIRefSet(context string, input interface{}) (slice []PCIRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PCIRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePCIRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePCIRefToPCIRecordMap(context string, input interface{}) (goMap map[PCIRef]PCIRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PCIRef]PCIRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePCIRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePCIRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializePGPURecord(context string, input interface{}) (record PGPURecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	pCIValue, ok := rpcStruct["PCI"]
	if ok && pCIValue != nil {
		record.PCI, err = deserializePCIRef(fmt.Sprintf("%s.%s", context, "PCI"), pCIValue)
		if err != nil {
			return
		}
	}
	gPUGroupValue, ok := rpcStruct["GPU_group"]
	if ok && gPUGroupValue != nil {
		record.GPUGroup, err = deserializeGPUGroupRef(fmt.Sprintf("%s.%s", context, "GPU_group"), gPUGroupValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	supportedVGPUTypesValue, ok := rpcStruct["supported_VGPU_types"]
	if ok && supportedVGPUTypesValue != nil {
		record.SupportedVGPUTypes, err = deserializeVGPUTypeRefSet(fmt.Sprintf("%s.%s", context, "supported_VGPU_types"), supportedVGPUTypesValue)
		if err != nil {
			return
		}
	}
	enabledVGPUTypesValue, ok := rpcStruct["enabled_VGPU_types"]
	if ok && enabledVGPUTypesValue != nil {
		record.EnabledVGPUTypes, err = deserializeVGPUTypeRefSet(fmt.Sprintf("%s.%s", context, "enabled_VGPU_types"), enabledVGPUTypesValue)
		if err != nil {
			return
		}
	}
	residentVGPUsValue, ok := rpcStruct["resident_VGPUs"]
	if ok && residentVGPUsValue != nil {
		record.ResidentVGPUs, err = deserializeVGPURefSet(fmt.Sprintf("%s.%s", context, "resident_VGPUs"), residentVGPUsValue)
		if err != nil {
			return
		}
	}
	supportedVGPUMaxCapacitiesValue, ok := rpcStruct["supported_VGPU_max_capacities"]
	if ok && supportedVGPUMaxCapacitiesValue != nil {
		record.SupportedVGPUMaxCapacities, err = deserializeVGPUTypeRefToIntMap(fmt.Sprintf("%s.%s", context, "supported_VGPU_max_capacities"), supportedVGPUMaxCapacitiesValue)
		if err != nil {
			return
		}
	}
	dom0AccessValue, ok := rpcStruct["dom0_access"]
	if ok && dom0AccessValue != nil {
		record.Dom0Access, err = deserializeEnumPciDom0Access(fmt.Sprintf("%s.%s", context, "dom0_access"), dom0AccessValue)
		if err != nil {
			return
		}
	}
	isSystemDisplayDeviceValue, ok := rpcStruct["is_system_display_device"]
	if ok && isSystemDisplayDeviceValue != nil {
		record.IsSystemDisplayDevice, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_system_display_device"), isSystemDisplayDeviceValue)
		if err != nil {
			return
		}
	}
	compatibilityMetadataValue, ok := rpcStruct["compatibility_metadata"]
	if ok && compatibilityMetadataValue != nil {
		record.CompatibilityMetadata, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "compatibility_metadata"), compatibilityMetadataValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePGPURef(context string, input interface{}) (PGPURef, error) {
	var ref PGPURef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PGPURef(value), nil
}

func deserializePGPURefSet(context string, input interface{}) (slice []PGPURef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PGPURef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePGPURef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePGPURefToPGPURecordMap(context string, input interface{}) (goMap map[PGPURef]PGPURecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PGPURef]PGPURecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePGPURef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePGPURecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializePIFMetricsRecord(context string, input interface{}) (record PIFMetricsRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	ioReadKbsValue, ok := rpcStruct["io_read_kbs"]
	if ok && ioReadKbsValue != nil {
		record.IoReadKbs, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "io_read_kbs"), ioReadKbsValue)
		if err != nil {
			return
		}
	}
	ioWriteKbsValue, ok := rpcStruct["io_write_kbs"]
	if ok && ioWriteKbsValue != nil {
		record.IoWriteKbs, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "io_write_kbs"), ioWriteKbsValue)
		if err != nil {
			return
		}
	}
	carrierValue, ok := rpcStruct["carrier"]
	if ok && carrierValue != nil {
		record.Carrier, err = deserializeBool(fmt.Sprintf("%s.%s", context, "carrier"), carrierValue)
		if err != nil {
			return
		}
	}
	vendorIDValue, ok := rpcStruct["vendor_id"]
	if ok && vendorIDValue != nil {
		record.VendorID, err = deserializeString(fmt.Sprintf("%s.%s", context, "vendor_id"), vendorIDValue)
		if err != nil {
			return
		}
	}
	vendorNameValue, ok := rpcStruct["vendor_name"]
	if ok && vendorNameValue != nil {
		record.VendorName, err = deserializeString(fmt.Sprintf("%s.%s", context, "vendor_name"), vendorNameValue)
		if err != nil {
			return
		}
	}
	deviceIDValue, ok := rpcStruct["device_id"]
	if ok && deviceIDValue != nil {
		record.DeviceID, err = deserializeString(fmt.Sprintf("%s.%s", context, "device_id"), deviceIDValue)
		if err != nil {
			return
		}
	}
	deviceNameValue, ok := rpcStruct["device_name"]
	if ok && deviceNameValue != nil {
		record.DeviceName, err = deserializeString(fmt.Sprintf("%s.%s", context, "device_name"), deviceNameValue)
		if err != nil {
			return
		}
	}
	speedValue, ok := rpcStruct["speed"]
	if ok && speedValue != nil {
		record.Speed, err = deserializeInt(fmt.Sprintf("%s.%s", context, "speed"), speedValue)
		if err != nil {
			return
		}
	}
	duplexValue, ok := rpcStruct["duplex"]
	if ok && duplexValue != nil {
		record.Duplex, err = deserializeBool(fmt.Sprintf("%s.%s", context, "duplex"), duplexValue)
		if err != nil {
			return
		}
	}
	pciBusPathValue, ok := rpcStruct["pci_bus_path"]
	if ok && pciBusPathValue != nil {
		record.PciBusPath, err = deserializeString(fmt.Sprintf("%s.%s", context, "pci_bus_path"), pciBusPathValue)
		if err != nil {
			return
		}
	}
	lastUpdatedValue, ok := rpcStruct["last_updated"]
	if ok && lastUpdatedValue != nil {
		record.LastUpdated, err = deserializeTime(fmt.Sprintf("%s.%s", context, "last_updated"), lastUpdatedValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePIFMetricsRef(context string, input interface{}) (PIFMetricsRef, error) {
	var ref PIFMetricsRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PIFMetricsRef(value), nil
}

func deserializePIFMetricsRefSet(context string, input interface{}) (slice []PIFMetricsRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PIFMetricsRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePIFMetricsRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePIFMetricsRefToPIFMetricsRecordMap(context string, input interface{}) (goMap map[PIFMetricsRef]PIFMetricsRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PIFMetricsRef]PIFMetricsRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePIFMetricsRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePIFMetricsRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializePIFRecord(context string, input interface{}) (record PIFRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	deviceValue, ok := rpcStruct["device"]
	if ok && deviceValue != nil {
		record.Device, err = deserializeString(fmt.Sprintf("%s.%s", context, "device"), deviceValue)
		if err != nil {
			return
		}
	}
	networkValue, ok := rpcStruct["network"]
	if ok && networkValue != nil {
		record.Network, err = deserializeNetworkRef(fmt.Sprintf("%s.%s", context, "network"), networkValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	mACValue, ok := rpcStruct["MAC"]
	if ok && mACValue != nil {
		record.MAC, err = deserializeString(fmt.Sprintf("%s.%s", context, "MAC"), mACValue)
		if err != nil {
			return
		}
	}
	mTUValue, ok := rpcStruct["MTU"]
	if ok && mTUValue != nil {
		record.MTU, err = deserializeInt(fmt.Sprintf("%s.%s", context, "MTU"), mTUValue)
		if err != nil {
			return
		}
	}
	vLANValue, ok := rpcStruct["VLAN"]
	if ok && vLANValue != nil {
		record.VLAN, err = deserializeInt(fmt.Sprintf("%s.%s", context, "VLAN"), vLANValue)
		if err != nil {
			return
		}
	}
	metricsValue, ok := rpcStruct["metrics"]
	if ok && metricsValue != nil {
		record.Metrics, err = deserializePIFMetricsRef(fmt.Sprintf("%s.%s", context, "metrics"), metricsValue)
		if err != nil {
			return
		}
	}
	physicalValue, ok := rpcStruct["physical"]
	if ok && physicalValue != nil {
		record.Physical, err = deserializeBool(fmt.Sprintf("%s.%s", context, "physical"), physicalValue)
		if err != nil {
			return
		}
	}
	currentlyAttachedValue, ok := rpcStruct["currently_attached"]
	if ok && currentlyAttachedValue != nil {
		record.CurrentlyAttached, err = deserializeBool(fmt.Sprintf("%s.%s", context, "currently_attached"), currentlyAttachedValue)
		if err != nil {
			return
		}
	}
	iPConfigurationModeValue, ok := rpcStruct["ip_configuration_mode"]
	if ok && iPConfigurationModeValue != nil {
		record.IPConfigurationMode, err = deserializeEnumIPConfigurationMode(fmt.Sprintf("%s.%s", context, "ip_configuration_mode"), iPConfigurationModeValue)
		if err != nil {
			return
		}
	}
	iPValue, ok := rpcStruct["IP"]
	if ok && iPValue != nil {
		record.IP, err = deserializeString(fmt.Sprintf("%s.%s", context, "IP"), iPValue)
		if err != nil {
			return
		}
	}
	netmaskValue, ok := rpcStruct["netmask"]
	if ok && netmaskValue != nil {
		record.Netmask, err = deserializeString(fmt.Sprintf("%s.%s", context, "netmask"), netmaskValue)
		if err != nil {
			return
		}
	}
	gatewayValue, ok := rpcStruct["gateway"]
	if ok && gatewayValue != nil {
		record.Gateway, err = deserializeString(fmt.Sprintf("%s.%s", context, "gateway"), gatewayValue)
		if err != nil {
			return
		}
	}
	dNSValue, ok := rpcStruct["DNS"]
	if ok && dNSValue != nil {
		record.DNS, err = deserializeString(fmt.Sprintf("%s.%s", context, "DNS"), dNSValue)
		if err != nil {
			return
		}
	}
	bondSlaveOfValue, ok := rpcStruct["bond_slave_of"]
	if ok && bondSlaveOfValue != nil {
		record.BondSlaveOf, err = deserializeBondRef(fmt.Sprintf("%s.%s", context, "bond_slave_of"), bondSlaveOfValue)
		if err != nil {
			return
		}
	}
	bondMasterOfValue, ok := rpcStruct["bond_master_of"]
	if ok && bondMasterOfValue != nil {
		record.BondMasterOf, err = deserializeBondRefSet(fmt.Sprintf("%s.%s", context, "bond_master_of"), bondMasterOfValue)
		if err != nil {
			return
		}
	}
	vLANMasterOfValue, ok := rpcStruct["VLAN_master_of"]
	if ok && vLANMasterOfValue != nil {
		record.VLANMasterOf, err = deserializeVLANRef(fmt.Sprintf("%s.%s", context, "VLAN_master_of"), vLANMasterOfValue)
		if err != nil {
			return
		}
	}
	vLANSlaveOfValue, ok := rpcStruct["VLAN_slave_of"]
	if ok && vLANSlaveOfValue != nil {
		record.VLANSlaveOf, err = deserializeVLANRefSet(fmt.Sprintf("%s.%s", context, "VLAN_slave_of"), vLANSlaveOfValue)
		if err != nil {
			return
		}
	}
	managementValue, ok := rpcStruct["management"]
	if ok && managementValue != nil {
		record.Management, err = deserializeBool(fmt.Sprintf("%s.%s", context, "management"), managementValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	disallowUnplugValue, ok := rpcStruct["disallow_unplug"]
	if ok && disallowUnplugValue != nil {
		record.DisallowUnplug, err = deserializeBool(fmt.Sprintf("%s.%s", context, "disallow_unplug"), disallowUnplugValue)
		if err != nil {
			return
		}
	}
	tunnelAccessPIFOfValue, ok := rpcStruct["tunnel_access_PIF_of"]
	if ok && tunnelAccessPIFOfValue != nil {
		record.TunnelAccessPIFOf, err = deserializeTunnelRefSet(fmt.Sprintf("%s.%s", context, "tunnel_access_PIF_of"), tunnelAccessPIFOfValue)
		if err != nil {
			return
		}
	}
	tunnelTransportPIFOfValue, ok := rpcStruct["tunnel_transport_PIF_of"]
	if ok && tunnelTransportPIFOfValue != nil {
		record.TunnelTransportPIFOf, err = deserializeTunnelRefSet(fmt.Sprintf("%s.%s", context, "tunnel_transport_PIF_of"), tunnelTransportPIFOfValue)
		if err != nil {
			return
		}
	}
	ipv6ConfigurationModeValue, ok := rpcStruct["ipv6_configuration_mode"]
	if ok && ipv6ConfigurationModeValue != nil {
		record.Ipv6ConfigurationMode, err = deserializeEnumIpv6ConfigurationMode(fmt.Sprintf("%s.%s", context, "ipv6_configuration_mode"), ipv6ConfigurationModeValue)
		if err != nil {
			return
		}
	}
	iPv6Value, ok := rpcStruct["IPv6"]
	if ok && iPv6Value != nil {
		record.IPv6, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "IPv6"), iPv6Value)
		if err != nil {
			return
		}
	}
	ipv6GatewayValue, ok := rpcStruct["ipv6_gateway"]
	if ok && ipv6GatewayValue != nil {
		record.Ipv6Gateway, err = deserializeString(fmt.Sprintf("%s.%s", context, "ipv6_gateway"), ipv6GatewayValue)
		if err != nil {
			return
		}
	}
	primaryAddressTypeValue, ok := rpcStruct["primary_address_type"]
	if ok && primaryAddressTypeValue != nil {
		record.PrimaryAddressType, err = deserializeEnumPrimaryAddressType(fmt.Sprintf("%s.%s", context, "primary_address_type"), primaryAddressTypeValue)
		if err != nil {
			return
		}
	}
	managedValue, ok := rpcStruct["managed"]
	if ok && managedValue != nil {
		record.Managed, err = deserializeBool(fmt.Sprintf("%s.%s", context, "managed"), managedValue)
		if err != nil {
			return
		}
	}
	propertiesValue, ok := rpcStruct["properties"]
	if ok && propertiesValue != nil {
		record.Properties, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "properties"), propertiesValue)
		if err != nil {
			return
		}
	}
	capabilitiesValue, ok := rpcStruct["capabilities"]
	if ok && capabilitiesValue != nil {
		record.Capabilities, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "capabilities"), capabilitiesValue)
		if err != nil {
			return
		}
	}
	igmpSnoopingStatusValue, ok := rpcStruct["igmp_snooping_status"]
	if ok && igmpSnoopingStatusValue != nil {
		record.IgmpSnoopingStatus, err = deserializeEnumPifIgmpStatus(fmt.Sprintf("%s.%s", context, "igmp_snooping_status"), igmpSnoopingStatusValue)
		if err != nil {
			return
		}
	}
	sriovPhysicalPIFOfValue, ok := rpcStruct["sriov_physical_PIF_of"]
	if ok && sriovPhysicalPIFOfValue != nil {
		record.SriovPhysicalPIFOf, err = deserializeNetworkSriovRefSet(fmt.Sprintf("%s.%s", context, "sriov_physical_PIF_of"), sriovPhysicalPIFOfValue)
		if err != nil {
			return
		}
	}
	sriovLogicalPIFOfValue, ok := rpcStruct["sriov_logical_PIF_of"]
	if ok && sriovLogicalPIFOfValue != nil {
		record.SriovLogicalPIFOf, err = deserializeNetworkSriovRefSet(fmt.Sprintf("%s.%s", context, "sriov_logical_PIF_of"), sriovLogicalPIFOfValue)
		if err != nil {
			return
		}
	}
	pCIValue, ok := rpcStruct["PCI"]
	if ok && pCIValue != nil {
		record.PCI, err = deserializePCIRef(fmt.Sprintf("%s.%s", context, "PCI"), pCIValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePIFRef(context string, input interface{}) (PIFRef, error) {
	var ref PIFRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PIFRef(value), nil
}

func deserializePIFRefSet(context string, input interface{}) (slice []PIFRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PIFRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePIFRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePIFRefToPIFRecordMap(context string, input interface{}) (goMap map[PIFRef]PIFRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PIFRef]PIFRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePIFRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePIFRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializePUSBRecord(context string, input interface{}) (record PUSBRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	uSBGroupValue, ok := rpcStruct["USB_group"]
	if ok && uSBGroupValue != nil {
		record.USBGroup, err = deserializeUSBGroupRef(fmt.Sprintf("%s.%s", context, "USB_group"), uSBGroupValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	pathValue, ok := rpcStruct["path"]
	if ok && pathValue != nil {
		record.Path, err = deserializeString(fmt.Sprintf("%s.%s", context, "path"), pathValue)
		if err != nil {
			return
		}
	}
	vendorIDValue, ok := rpcStruct["vendor_id"]
	if ok && vendorIDValue != nil {
		record.VendorID, err = deserializeString(fmt.Sprintf("%s.%s", context, "vendor_id"), vendorIDValue)
		if err != nil {
			return
		}
	}
	vendorDescValue, ok := rpcStruct["vendor_desc"]
	if ok && vendorDescValue != nil {
		record.VendorDesc, err = deserializeString(fmt.Sprintf("%s.%s", context, "vendor_desc"), vendorDescValue)
		if err != nil {
			return
		}
	}
	productIDValue, ok := rpcStruct["product_id"]
	if ok && productIDValue != nil {
		record.ProductID, err = deserializeString(fmt.Sprintf("%s.%s", context, "product_id"), productIDValue)
		if err != nil {
			return
		}
	}
	productDescValue, ok := rpcStruct["product_desc"]
	if ok && productDescValue != nil {
		record.ProductDesc, err = deserializeString(fmt.Sprintf("%s.%s", context, "product_desc"), productDescValue)
		if err != nil {
			return
		}
	}
	serialValue, ok := rpcStruct["serial"]
	if ok && serialValue != nil {
		record.Serial, err = deserializeString(fmt.Sprintf("%s.%s", context, "serial"), serialValue)
		if err != nil {
			return
		}
	}
	versionValue, ok := rpcStruct["version"]
	if ok && versionValue != nil {
		record.Version, err = deserializeString(fmt.Sprintf("%s.%s", context, "version"), versionValue)
		if err != nil {
			return
		}
	}
	descriptionValue, ok := rpcStruct["description"]
	if ok && descriptionValue != nil {
		record.Description, err = deserializeString(fmt.Sprintf("%s.%s", context, "description"), descriptionValue)
		if err != nil {
			return
		}
	}
	passthroughEnabledValue, ok := rpcStruct["passthrough_enabled"]
	if ok && passthroughEnabledValue != nil {
		record.PassthroughEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "passthrough_enabled"), passthroughEnabledValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	speedValue, ok := rpcStruct["speed"]
	if ok && speedValue != nil {
		record.Speed, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "speed"), speedValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePUSBRef(context string, input interface{}) (PUSBRef, error) {
	var ref PUSBRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PUSBRef(value), nil
}

func deserializePUSBRefSet(context string, input interface{}) (slice []PUSBRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PUSBRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePUSBRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePUSBRefToPUSBRecordMap(context string, input interface{}) (goMap map[PUSBRef]PUSBRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PUSBRef]PUSBRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePUSBRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePUSBRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializePVSCacheStorageRecord(context string, input interface{}) (record PVSCacheStorageRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	hostValue, ok := rpcStruct["host"]
	if ok && hostValue != nil {
		record.Host, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "host"), hostValue)
		if err != nil {
			return
		}
	}
	sRValue, ok := rpcStruct["SR"]
	if ok && sRValue != nil {
		record.SR, err = deserializeSRRef(fmt.Sprintf("%s.%s", context, "SR"), sRValue)
		if err != nil {
			return
		}
	}
	siteValue, ok := rpcStruct["site"]
	if ok && siteValue != nil {
		record.Site, err = deserializePVSSiteRef(fmt.Sprintf("%s.%s", context, "site"), siteValue)
		if err != nil {
			return
		}
	}
	sizeValue, ok := rpcStruct["size"]
	if ok && sizeValue != nil {
		record.Size, err = deserializeInt(fmt.Sprintf("%s.%s", context, "size"), sizeValue)
		if err != nil {
			return
		}
	}
	vDIValue, ok := rpcStruct["VDI"]
	if ok && vDIValue != nil {
		record.VDI, err = deserializeVDIRef(fmt.Sprintf("%s.%s", context, "VDI"), vDIValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePVSCacheStorageRef(context string, input interface{}) (PVSCacheStorageRef, error) {
	var ref PVSCacheStorageRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PVSCacheStorageRef(value), nil
}

func deserializePVSCacheStorageRefSet(context string, input interface{}) (slice []PVSCacheStorageRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PVSCacheStorageRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePVSCacheStorageRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePVSCacheStorageRefToPVSCacheStorageRecordMap(context string, input interface{}) (goMap map[PVSCacheStorageRef]PVSCacheStorageRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PVSCacheStorageRef]PVSCacheStorageRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePVSCacheStorageRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePVSCacheStorageRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializePVSProxyRecord(context string, input interface{}) (record PVSProxyRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	siteValue, ok := rpcStruct["site"]
	if ok && siteValue != nil {
		record.Site, err = deserializePVSSiteRef(fmt.Sprintf("%s.%s", context, "site"), siteValue)
		if err != nil {
			return
		}
	}
	vIFValue, ok := rpcStruct["VIF"]
	if ok && vIFValue != nil {
		record.VIF, err = deserializeVIFRef(fmt.Sprintf("%s.%s", context, "VIF"), vIFValue)
		if err != nil {
			return
		}
	}
	currentlyAttachedValue, ok := rpcStruct["currently_attached"]
	if ok && currentlyAttachedValue != nil {
		record.CurrentlyAttached, err = deserializeBool(fmt.Sprintf("%s.%s", context, "currently_attached"), currentlyAttachedValue)
		if err != nil {
			return
		}
	}
	statusValue, ok := rpcStruct["status"]
	if ok && statusValue != nil {
		record.Status, err = deserializeEnumPvsProxyStatus(fmt.Sprintf("%s.%s", context, "status"), statusValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePVSProxyRef(context string, input interface{}) (PVSProxyRef, error) {
	var ref PVSProxyRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PVSProxyRef(value), nil
}

func deserializePVSProxyRefSet(context string, input interface{}) (slice []PVSProxyRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PVSProxyRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePVSProxyRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePVSProxyRefToPVSProxyRecordMap(context string, input interface{}) (goMap map[PVSProxyRef]PVSProxyRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PVSProxyRef]PVSProxyRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePVSProxyRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePVSProxyRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializePVSServerRecord(context string, input interface{}) (record PVSServerRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	addressesValue, ok := rpcStruct["addresses"]
	if ok && addressesValue != nil {
		record.Addresses, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "addresses"), addressesValue)
		if err != nil {
			return
		}
	}
	firstPortValue, ok := rpcStruct["first_port"]
	if ok && firstPortValue != nil {
		record.FirstPort, err = deserializeInt(fmt.Sprintf("%s.%s", context, "first_port"), firstPortValue)
		if err != nil {
			return
		}
	}
	lastPortValue, ok := rpcStruct["last_port"]
	if ok && lastPortValue != nil {
		record.LastPort, err = deserializeInt(fmt.Sprintf("%s.%s", context, "last_port"), lastPortValue)
		if err != nil {
			return
		}
	}
	siteValue, ok := rpcStruct["site"]
	if ok && siteValue != nil {
		record.Site, err = deserializePVSSiteRef(fmt.Sprintf("%s.%s", context, "site"), siteValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePVSServerRef(context string, input interface{}) (PVSServerRef, error) {
	var ref PVSServerRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PVSServerRef(value), nil
}

func deserializePVSServerRefSet(context string, input interface{}) (slice []PVSServerRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PVSServerRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePVSServerRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePVSServerRefToPVSServerRecordMap(context string, input interface{}) (goMap map[PVSServerRef]PVSServerRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PVSServerRef]PVSServerRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePVSServerRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePVSServerRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializePVSSiteRecord(context string, input interface{}) (record PVSSiteRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	pVSUUIDValue, ok := rpcStruct["PVS_uuid"]
	if ok && pVSUUIDValue != nil {
		record.PVSUUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "PVS_uuid"), pVSUUIDValue)
		if err != nil {
			return
		}
	}
	cacheStorageValue, ok := rpcStruct["cache_storage"]
	if ok && cacheStorageValue != nil {
		record.CacheStorage, err = deserializePVSCacheStorageRefSet(fmt.Sprintf("%s.%s", context, "cache_storage"), cacheStorageValue)
		if err != nil {
			return
		}
	}
	serversValue, ok := rpcStruct["servers"]
	if ok && serversValue != nil {
		record.Servers, err = deserializePVSServerRefSet(fmt.Sprintf("%s.%s", context, "servers"), serversValue)
		if err != nil {
			return
		}
	}
	proxiesValue, ok := rpcStruct["proxies"]
	if ok && proxiesValue != nil {
		record.Proxies, err = deserializePVSProxyRefSet(fmt.Sprintf("%s.%s", context, "proxies"), proxiesValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePVSSiteRef(context string, input interface{}) (PVSSiteRef, error) {
	var ref PVSSiteRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PVSSiteRef(value), nil
}

func deserializePVSSiteRefSet(context string, input interface{}) (slice []PVSSiteRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PVSSiteRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePVSSiteRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePVSSiteRefToPVSSiteRecordMap(context string, input interface{}) (goMap map[PVSSiteRef]PVSSiteRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PVSSiteRef]PVSSiteRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePVSSiteRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePVSSiteRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializePoolPatchRecord(context string, input interface{}) (record PoolPatchRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	versionValue, ok := rpcStruct["version"]
	if ok && versionValue != nil {
		record.Version, err = deserializeString(fmt.Sprintf("%s.%s", context, "version"), versionValue)
		if err != nil {
			return
		}
	}
	sizeValue, ok := rpcStruct["size"]
	if ok && sizeValue != nil {
		record.Size, err = deserializeInt(fmt.Sprintf("%s.%s", context, "size"), sizeValue)
		if err != nil {
			return
		}
	}
	poolAppliedValue, ok := rpcStruct["pool_applied"]
	if ok && poolAppliedValue != nil {
		record.PoolApplied, err = deserializeBool(fmt.Sprintf("%s.%s", context, "pool_applied"), poolAppliedValue)
		if err != nil {
			return
		}
	}
	hostPatchesValue, ok := rpcStruct["host_patches"]
	if ok && hostPatchesValue != nil {
		record.HostPatches, err = deserializeHostPatchRefSet(fmt.Sprintf("%s.%s", context, "host_patches"), hostPatchesValue)
		if err != nil {
			return
		}
	}
	afterApplyGuidanceValue, ok := rpcStruct["after_apply_guidance"]
	if ok && afterApplyGuidanceValue != nil {
		record.AfterApplyGuidance, err = deserializeEnumAfterApplyGuidanceSet(fmt.Sprintf("%s.%s", context, "after_apply_guidance"), afterApplyGuidanceValue)
		if err != nil {
			return
		}
	}
	poolUpdateValue, ok := rpcStruct["pool_update"]
	if ok && poolUpdateValue != nil {
		record.PoolUpdate, err = deserializePoolUpdateRef(fmt.Sprintf("%s.%s", context, "pool_update"), poolUpdateValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePoolPatchRef(context string, input interface{}) (PoolPatchRef, error) {
	var ref PoolPatchRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PoolPatchRef(value), nil
}

func deserializePoolPatchRefSet(context string, input interface{}) (slice []PoolPatchRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PoolPatchRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePoolPatchRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePoolPatchRefToPoolPatchRecordMap(context string, input interface{}) (goMap map[PoolPatchRef]PoolPatchRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PoolPatchRef]PoolPatchRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePoolPatchRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePoolPatchRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializePoolRecord(context string, input interface{}) (record PoolRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	masterValue, ok := rpcStruct["master"]
	if ok && masterValue != nil {
		record.Master, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "master"), masterValue)
		if err != nil {
			return
		}
	}
	defaultSRValue, ok := rpcStruct["default_SR"]
	if ok && defaultSRValue != nil {
		record.DefaultSR, err = deserializeSRRef(fmt.Sprintf("%s.%s", context, "default_SR"), defaultSRValue)
		if err != nil {
			return
		}
	}
	suspendImageSRValue, ok := rpcStruct["suspend_image_SR"]
	if ok && suspendImageSRValue != nil {
		record.SuspendImageSR, err = deserializeSRRef(fmt.Sprintf("%s.%s", context, "suspend_image_SR"), suspendImageSRValue)
		if err != nil {
			return
		}
	}
	crashDumpSRValue, ok := rpcStruct["crash_dump_SR"]
	if ok && crashDumpSRValue != nil {
		record.CrashDumpSR, err = deserializeSRRef(fmt.Sprintf("%s.%s", context, "crash_dump_SR"), crashDumpSRValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	haEnabledValue, ok := rpcStruct["ha_enabled"]
	if ok && haEnabledValue != nil {
		record.HaEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "ha_enabled"), haEnabledValue)
		if err != nil {
			return
		}
	}
	haConfigurationValue, ok := rpcStruct["ha_configuration"]
	if ok && haConfigurationValue != nil {
		record.HaConfiguration, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "ha_configuration"), haConfigurationValue)
		if err != nil {
			return
		}
	}
	haStatefilesValue, ok := rpcStruct["ha_statefiles"]
	if ok && haStatefilesValue != nil {
		record.HaStatefiles, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "ha_statefiles"), haStatefilesValue)
		if err != nil {
			return
		}
	}
	haHostFailuresToTolerateValue, ok := rpcStruct["ha_host_failures_to_tolerate"]
	if ok && haHostFailuresToTolerateValue != nil {
		record.HaHostFailuresToTolerate, err = deserializeInt(fmt.Sprintf("%s.%s", context, "ha_host_failures_to_tolerate"), haHostFailuresToTolerateValue)
		if err != nil {
			return
		}
	}
	haPlanExistsForValue, ok := rpcStruct["ha_plan_exists_for"]
	if ok && haPlanExistsForValue != nil {
		record.HaPlanExistsFor, err = deserializeInt(fmt.Sprintf("%s.%s", context, "ha_plan_exists_for"), haPlanExistsForValue)
		if err != nil {
			return
		}
	}
	haAllowOvercommitValue, ok := rpcStruct["ha_allow_overcommit"]
	if ok && haAllowOvercommitValue != nil {
		record.HaAllowOvercommit, err = deserializeBool(fmt.Sprintf("%s.%s", context, "ha_allow_overcommit"), haAllowOvercommitValue)
		if err != nil {
			return
		}
	}
	haOvercommittedValue, ok := rpcStruct["ha_overcommitted"]
	if ok && haOvercommittedValue != nil {
		record.HaOvercommitted, err = deserializeBool(fmt.Sprintf("%s.%s", context, "ha_overcommitted"), haOvercommittedValue)
		if err != nil {
			return
		}
	}
	blobsValue, ok := rpcStruct["blobs"]
	if ok && blobsValue != nil {
		record.Blobs, err = deserializeStringToBlobRefMap(fmt.Sprintf("%s.%s", context, "blobs"), blobsValue)
		if err != nil {
			return
		}
	}
	tagsValue, ok := rpcStruct["tags"]
	if ok && tagsValue != nil {
		record.Tags, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "tags"), tagsValue)
		if err != nil {
			return
		}
	}
	guiConfigValue, ok := rpcStruct["gui_config"]
	if ok && guiConfigValue != nil {
		record.GuiConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "gui_config"), guiConfigValue)
		if err != nil {
			return
		}
	}
	healthCheckConfigValue, ok := rpcStruct["health_check_config"]
	if ok && healthCheckConfigValue != nil {
		record.HealthCheckConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "health_check_config"), healthCheckConfigValue)
		if err != nil {
			return
		}
	}
	wlbURLValue, ok := rpcStruct["wlb_url"]
	if ok && wlbURLValue != nil {
		record.WlbURL, err = deserializeString(fmt.Sprintf("%s.%s", context, "wlb_url"), wlbURLValue)
		if err != nil {
			return
		}
	}
	wlbUsernameValue, ok := rpcStruct["wlb_username"]
	if ok && wlbUsernameValue != nil {
		record.WlbUsername, err = deserializeString(fmt.Sprintf("%s.%s", context, "wlb_username"), wlbUsernameValue)
		if err != nil {
			return
		}
	}
	wlbEnabledValue, ok := rpcStruct["wlb_enabled"]
	if ok && wlbEnabledValue != nil {
		record.WlbEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "wlb_enabled"), wlbEnabledValue)
		if err != nil {
			return
		}
	}
	wlbVerifyCertValue, ok := rpcStruct["wlb_verify_cert"]
	if ok && wlbVerifyCertValue != nil {
		record.WlbVerifyCert, err = deserializeBool(fmt.Sprintf("%s.%s", context, "wlb_verify_cert"), wlbVerifyCertValue)
		if err != nil {
			return
		}
	}
	redoLogEnabledValue, ok := rpcStruct["redo_log_enabled"]
	if ok && redoLogEnabledValue != nil {
		record.RedoLogEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "redo_log_enabled"), redoLogEnabledValue)
		if err != nil {
			return
		}
	}
	redoLogVdiValue, ok := rpcStruct["redo_log_vdi"]
	if ok && redoLogVdiValue != nil {
		record.RedoLogVdi, err = deserializeVDIRef(fmt.Sprintf("%s.%s", context, "redo_log_vdi"), redoLogVdiValue)
		if err != nil {
			return
		}
	}
	vswitchControllerValue, ok := rpcStruct["vswitch_controller"]
	if ok && vswitchControllerValue != nil {
		record.VswitchController, err = deserializeString(fmt.Sprintf("%s.%s", context, "vswitch_controller"), vswitchControllerValue)
		if err != nil {
			return
		}
	}
	restrictionsValue, ok := rpcStruct["restrictions"]
	if ok && restrictionsValue != nil {
		record.Restrictions, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "restrictions"), restrictionsValue)
		if err != nil {
			return
		}
	}
	metadataVDIsValue, ok := rpcStruct["metadata_VDIs"]
	if ok && metadataVDIsValue != nil {
		record.MetadataVDIs, err = deserializeVDIRefSet(fmt.Sprintf("%s.%s", context, "metadata_VDIs"), metadataVDIsValue)
		if err != nil {
			return
		}
	}
	haClusterStackValue, ok := rpcStruct["ha_cluster_stack"]
	if ok && haClusterStackValue != nil {
		record.HaClusterStack, err = deserializeString(fmt.Sprintf("%s.%s", context, "ha_cluster_stack"), haClusterStackValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumPoolAllowedOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumPoolAllowedOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	guestAgentConfigValue, ok := rpcStruct["guest_agent_config"]
	if ok && guestAgentConfigValue != nil {
		record.GuestAgentConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "guest_agent_config"), guestAgentConfigValue)
		if err != nil {
			return
		}
	}
	cPUInfoValue, ok := rpcStruct["cpu_info"]
	if ok && cPUInfoValue != nil {
		record.CPUInfo, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "cpu_info"), cPUInfoValue)
		if err != nil {
			return
		}
	}
	policyNoVendorDeviceValue, ok := rpcStruct["policy_no_vendor_device"]
	if ok && policyNoVendorDeviceValue != nil {
		record.PolicyNoVendorDevice, err = deserializeBool(fmt.Sprintf("%s.%s", context, "policy_no_vendor_device"), policyNoVendorDeviceValue)
		if err != nil {
			return
		}
	}
	livePatchingDisabledValue, ok := rpcStruct["live_patching_disabled"]
	if ok && livePatchingDisabledValue != nil {
		record.LivePatchingDisabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "live_patching_disabled"), livePatchingDisabledValue)
		if err != nil {
			return
		}
	}
	igmpSnoopingEnabledValue, ok := rpcStruct["igmp_snooping_enabled"]
	if ok && igmpSnoopingEnabledValue != nil {
		record.IgmpSnoopingEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "igmp_snooping_enabled"), igmpSnoopingEnabledValue)
		if err != nil {
			return
		}
	}
	uefiCertificatesValue, ok := rpcStruct["uefi_certificates"]
	if ok && uefiCertificatesValue != nil {
		record.UefiCertificates, err = deserializeString(fmt.Sprintf("%s.%s", context, "uefi_certificates"), uefiCertificatesValue)
		if err != nil {
			return
		}
	}
	customUefiCertificatesValue, ok := rpcStruct["custom_uefi_certificates"]
	if ok && customUefiCertificatesValue != nil {
		record.CustomUefiCertificates, err = deserializeString(fmt.Sprintf("%s.%s", context, "custom_uefi_certificates"), customUefiCertificatesValue)
		if err != nil {
			return
		}
	}
	isPsrPendingValue, ok := rpcStruct["is_psr_pending"]
	if ok && isPsrPendingValue != nil {
		record.IsPsrPending, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_psr_pending"), isPsrPendingValue)
		if err != nil {
			return
		}
	}
	tLSVerificationEnabledValue, ok := rpcStruct["tls_verification_enabled"]
	if ok && tLSVerificationEnabledValue != nil {
		record.TLSVerificationEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "tls_verification_enabled"), tLSVerificationEnabledValue)
		if err != nil {
			return
		}
	}
	repositoriesValue, ok := rpcStruct["repositories"]
	if ok && repositoriesValue != nil {
		record.Repositories, err = deserializeRepositoryRefSet(fmt.Sprintf("%s.%s", context, "repositories"), repositoriesValue)
		if err != nil {
			return
		}
	}
	clientCertificateAuthEnabledValue, ok := rpcStruct["client_certificate_auth_enabled"]
	if ok && clientCertificateAuthEnabledValue != nil {
		record.ClientCertificateAuthEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "client_certificate_auth_enabled"), clientCertificateAuthEnabledValue)
		if err != nil {
			return
		}
	}
	clientCertificateAuthNameValue, ok := rpcStruct["client_certificate_auth_name"]
	if ok && clientCertificateAuthNameValue != nil {
		record.ClientCertificateAuthName, err = deserializeString(fmt.Sprintf("%s.%s", context, "client_certificate_auth_name"), clientCertificateAuthNameValue)
		if err != nil {
			return
		}
	}
	repositoryProxyURLValue, ok := rpcStruct["repository_proxy_url"]
	if ok && repositoryProxyURLValue != nil {
		record.RepositoryProxyURL, err = deserializeString(fmt.Sprintf("%s.%s", context, "repository_proxy_url"), repositoryProxyURLValue)
		if err != nil {
			return
		}
	}
	repositoryProxyUsernameValue, ok := rpcStruct["repository_proxy_username"]
	if ok && repositoryProxyUsernameValue != nil {
		record.RepositoryProxyUsername, err = deserializeString(fmt.Sprintf("%s.%s", context, "repository_proxy_username"), repositoryProxyUsernameValue)
		if err != nil {
			return
		}
	}
	repositoryProxyPasswordValue, ok := rpcStruct["repository_proxy_password"]
	if ok && repositoryProxyPasswordValue != nil {
		record.RepositoryProxyPassword, err = deserializeSecretRef(fmt.Sprintf("%s.%s", context, "repository_proxy_password"), repositoryProxyPasswordValue)
		if err != nil {
			return
		}
	}
	migrationCompressionValue, ok := rpcStruct["migration_compression"]
	if ok && migrationCompressionValue != nil {
		record.MigrationCompression, err = deserializeBool(fmt.Sprintf("%s.%s", context, "migration_compression"), migrationCompressionValue)
		if err != nil {
			return
		}
	}
	coordinatorBiasValue, ok := rpcStruct["coordinator_bias"]
	if ok && coordinatorBiasValue != nil {
		record.CoordinatorBias, err = deserializeBool(fmt.Sprintf("%s.%s", context, "coordinator_bias"), coordinatorBiasValue)
		if err != nil {
			return
		}
	}
	localAuthMaxThreadsValue, ok := rpcStruct["local_auth_max_threads"]
	if ok && localAuthMaxThreadsValue != nil {
		record.LocalAuthMaxThreads, err = deserializeInt(fmt.Sprintf("%s.%s", context, "local_auth_max_threads"), localAuthMaxThreadsValue)
		if err != nil {
			return
		}
	}
	extAuthMaxThreadsValue, ok := rpcStruct["ext_auth_max_threads"]
	if ok && extAuthMaxThreadsValue != nil {
		record.ExtAuthMaxThreads, err = deserializeInt(fmt.Sprintf("%s.%s", context, "ext_auth_max_threads"), extAuthMaxThreadsValue)
		if err != nil {
			return
		}
	}
	extAuthCacheEnabledValue, ok := rpcStruct["ext_auth_cache_enabled"]
	if ok && extAuthCacheEnabledValue != nil {
		record.ExtAuthCacheEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "ext_auth_cache_enabled"), extAuthCacheEnabledValue)
		if err != nil {
			return
		}
	}
	extAuthCacheSizeValue, ok := rpcStruct["ext_auth_cache_size"]
	if ok && extAuthCacheSizeValue != nil {
		record.ExtAuthCacheSize, err = deserializeInt(fmt.Sprintf("%s.%s", context, "ext_auth_cache_size"), extAuthCacheSizeValue)
		if err != nil {
			return
		}
	}
	extAuthCacheExpiryValue, ok := rpcStruct["ext_auth_cache_expiry"]
	if ok && extAuthCacheExpiryValue != nil {
		record.ExtAuthCacheExpiry, err = deserializeInt(fmt.Sprintf("%s.%s", context, "ext_auth_cache_expiry"), extAuthCacheExpiryValue)
		if err != nil {
			return
		}
	}
	telemetryUUIDValue, ok := rpcStruct["telemetry_uuid"]
	if ok && telemetryUUIDValue != nil {
		record.TelemetryUUID, err = deserializeSecretRef(fmt.Sprintf("%s.%s", context, "telemetry_uuid"), telemetryUUIDValue)
		if err != nil {
			return
		}
	}
	telemetryFrequencyValue, ok := rpcStruct["telemetry_frequency"]
	if ok && telemetryFrequencyValue != nil {
		record.TelemetryFrequency, err = deserializeEnumTelemetryFrequency(fmt.Sprintf("%s.%s", context, "telemetry_frequency"), telemetryFrequencyValue)
		if err != nil {
			return
		}
	}
	telemetryNextCollectionValue, ok := rpcStruct["telemetry_next_collection"]
	if ok && telemetryNextCollectionValue != nil {
		record.TelemetryNextCollection, err = deserializeTime(fmt.Sprintf("%s.%s", context, "telemetry_next_collection"), telemetryNextCollectionValue)
		if err != nil {
			return
		}
	}
	lastUpdateSyncValue, ok := rpcStruct["last_update_sync"]
	if ok && lastUpdateSyncValue != nil {
		record.LastUpdateSync, err = deserializeTime(fmt.Sprintf("%s.%s", context, "last_update_sync"), lastUpdateSyncValue)
		if err != nil {
			return
		}
	}
	updateSyncFrequencyValue, ok := rpcStruct["update_sync_frequency"]
	if ok && updateSyncFrequencyValue != nil {
		record.UpdateSyncFrequency, err = deserializeEnumUpdateSyncFrequency(fmt.Sprintf("%s.%s", context, "update_sync_frequency"), updateSyncFrequencyValue)
		if err != nil {
			return
		}
	}
	updateSyncDayValue, ok := rpcStruct["update_sync_day"]
	if ok && updateSyncDayValue != nil {
		record.UpdateSyncDay, err = deserializeInt(fmt.Sprintf("%s.%s", context, "update_sync_day"), updateSyncDayValue)
		if err != nil {
			return
		}
	}
	updateSyncEnabledValue, ok := rpcStruct["update_sync_enabled"]
	if ok && updateSyncEnabledValue != nil {
		record.UpdateSyncEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "update_sync_enabled"), updateSyncEnabledValue)
		if err != nil {
			return
		}
	}
	recommendationsValue, ok := rpcStruct["recommendations"]
	if ok && recommendationsValue != nil {
		record.Recommendations, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "recommendations"), recommendationsValue)
		if err != nil {
			return
		}
	}
	licenseServerValue, ok := rpcStruct["license_server"]
	if ok && licenseServerValue != nil {
		record.LicenseServer, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "license_server"), licenseServerValue)
		if err != nil {
			return
		}
	}
	haRebootVMOnInternalShutdownValue, ok := rpcStruct["ha_reboot_vm_on_internal_shutdown"]
	if ok && haRebootVMOnInternalShutdownValue != nil {
		record.HaRebootVMOnInternalShutdown, err = deserializeBool(fmt.Sprintf("%s.%s", context, "ha_reboot_vm_on_internal_shutdown"), haRebootVMOnInternalShutdownValue)
		if err != nil {
			return
		}
	}
	limitConsoleSessionsValue, ok := rpcStruct["limit_console_sessions"]
	if ok && limitConsoleSessionsValue != nil {
		record.LimitConsoleSessions, err = deserializeBool(fmt.Sprintf("%s.%s", context, "limit_console_sessions"), limitConsoleSessionsValue)
		if err != nil {
			return
		}
	}
	vMConsoleIdleTimeoutValue, ok := rpcStruct["vm_console_idle_timeout"]
	if ok && vMConsoleIdleTimeoutValue != nil {
		record.VMConsoleIdleTimeout, err = deserializeInt(fmt.Sprintf("%s.%s", context, "vm_console_idle_timeout"), vMConsoleIdleTimeoutValue)
		if err != nil {
			return
		}
	}
	autoUpdateVMSecurebootCertificatesValue, ok := rpcStruct["auto_update_vm_secureboot_certificates"]
	if ok && autoUpdateVMSecurebootCertificatesValue != nil {
		record.AutoUpdateVMSecurebootCertificates, err = deserializeBool(fmt.Sprintf("%s.%s", context, "auto_update_vm_secureboot_certificates"), autoUpdateVMSecurebootCertificatesValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePoolRef(context string, input interface{}) (PoolRef, error) {
	var ref PoolRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PoolRef(value), nil
}

func deserializePoolRefSet(context string, input interface{}) (slice []PoolRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PoolRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePoolRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePoolRefToPoolRecordMap(context string, input interface{}) (goMap map[PoolRef]PoolRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PoolRef]PoolRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePoolRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePoolRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializePoolUpdateRecord(context string, input interface{}) (record PoolUpdateRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	versionValue, ok := rpcStruct["version"]
	if ok && versionValue != nil {
		record.Version, err = deserializeString(fmt.Sprintf("%s.%s", context, "version"), versionValue)
		if err != nil {
			return
		}
	}
	installationSizeValue, ok := rpcStruct["installation_size"]
	if ok && installationSizeValue != nil {
		record.InstallationSize, err = deserializeInt(fmt.Sprintf("%s.%s", context, "installation_size"), installationSizeValue)
		if err != nil {
			return
		}
	}
	keyValue, ok := rpcStruct["key"]
	if ok && keyValue != nil {
		record.Key, err = deserializeString(fmt.Sprintf("%s.%s", context, "key"), keyValue)
		if err != nil {
			return
		}
	}
	afterApplyGuidanceValue, ok := rpcStruct["after_apply_guidance"]
	if ok && afterApplyGuidanceValue != nil {
		record.AfterApplyGuidance, err = deserializeEnumUpdateAfterApplyGuidanceSet(fmt.Sprintf("%s.%s", context, "after_apply_guidance"), afterApplyGuidanceValue)
		if err != nil {
			return
		}
	}
	vdiValue, ok := rpcStruct["vdi"]
	if ok && vdiValue != nil {
		record.Vdi, err = deserializeVDIRef(fmt.Sprintf("%s.%s", context, "vdi"), vdiValue)
		if err != nil {
			return
		}
	}
	hostsValue, ok := rpcStruct["hosts"]
	if ok && hostsValue != nil {
		record.Hosts, err = deserializeHostRefSet(fmt.Sprintf("%s.%s", context, "hosts"), hostsValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	enforceHomogeneityValue, ok := rpcStruct["enforce_homogeneity"]
	if ok && enforceHomogeneityValue != nil {
		record.EnforceHomogeneity, err = deserializeBool(fmt.Sprintf("%s.%s", context, "enforce_homogeneity"), enforceHomogeneityValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializePoolUpdateRef(context string, input interface{}) (PoolUpdateRef, error) {
	var ref PoolUpdateRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return PoolUpdateRef(value), nil
}

func deserializePoolUpdateRefSet(context string, input interface{}) (slice []PoolUpdateRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]PoolUpdateRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializePoolUpdateRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializePoolUpdateRefToPoolUpdateRecordMap(context string, input interface{}) (goMap map[PoolUpdateRef]PoolUpdateRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[PoolUpdateRef]PoolUpdateRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializePoolUpdateRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializePoolUpdateRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeProbeResultRecord(context string, input interface{}) (record ProbeResultRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	configurationValue, ok := rpcStruct["configuration"]
	if ok && configurationValue != nil {
		record.Configuration, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "configuration"), configurationValue)
		if err != nil {
			return
		}
	}
	completeValue, ok := rpcStruct["complete"]
	if ok && completeValue != nil {
		record.Complete, err = deserializeBool(fmt.Sprintf("%s.%s", context, "complete"), completeValue)
		if err != nil {
			return
		}
	}
	record.Sr, err = deserializeOptionSrStatRecord(fmt.Sprintf("%s.%s", context, "sr"), rpcStruct["sr"])
	if err != nil {
		return
	}
	extraInfoValue, ok := rpcStruct["extra_info"]
	if ok && extraInfoValue != nil {
		record.ExtraInfo, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "extra_info"), extraInfoValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeProbeResultRecordSet(context string, input interface{}) (slice []ProbeResultRecord, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]ProbeResultRecord, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeProbeResultRecord(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeRepositoryRecord(context string, input interface{}) (record RepositoryRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	binaryURLValue, ok := rpcStruct["binary_url"]
	if ok && binaryURLValue != nil {
		record.BinaryURL, err = deserializeString(fmt.Sprintf("%s.%s", context, "binary_url"), binaryURLValue)
		if err != nil {
			return
		}
	}
	sourceURLValue, ok := rpcStruct["source_url"]
	if ok && sourceURLValue != nil {
		record.SourceURL, err = deserializeString(fmt.Sprintf("%s.%s", context, "source_url"), sourceURLValue)
		if err != nil {
			return
		}
	}
	updateValue, ok := rpcStruct["update"]
	if ok && updateValue != nil {
		record.Update, err = deserializeBool(fmt.Sprintf("%s.%s", context, "update"), updateValue)
		if err != nil {
			return
		}
	}
	hashValue, ok := rpcStruct["hash"]
	if ok && hashValue != nil {
		record.Hash, err = deserializeString(fmt.Sprintf("%s.%s", context, "hash"), hashValue)
		if err != nil {
			return
		}
	}
	upToDateValue, ok := rpcStruct["up_to_date"]
	if ok && upToDateValue != nil {
		record.UpToDate, err = deserializeBool(fmt.Sprintf("%s.%s", context, "up_to_date"), upToDateValue)
		if err != nil {
			return
		}
	}
	gpgkeyPathValue, ok := rpcStruct["gpgkey_path"]
	if ok && gpgkeyPathValue != nil {
		record.GpgkeyPath, err = deserializeString(fmt.Sprintf("%s.%s", context, "gpgkey_path"), gpgkeyPathValue)
		if err != nil {
			return
		}
	}
	originValue, ok := rpcStruct["origin"]
	if ok && originValue != nil {
		record.Origin, err = deserializeEnumOrigin(fmt.Sprintf("%s.%s", context, "origin"), originValue)
		if err != nil {
			return
		}
	}
	certificateValue, ok := rpcStruct["certificate"]
	if ok && certificateValue != nil {
		record.Certificate, err = deserializeString(fmt.Sprintf("%s.%s", context, "certificate"), certificateValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeRepositoryRef(context string, input interface{}) (RepositoryRef, error) {
	var ref RepositoryRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return RepositoryRef(value), nil
}

func deserializeRepositoryRefSet(context string, input interface{}) (slice []RepositoryRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]RepositoryRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeRepositoryRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeRepositoryRefToRepositoryRecordMap(context string, input interface{}) (goMap map[RepositoryRef]RepositoryRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[RepositoryRef]RepositoryRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeRepositoryRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeRepositoryRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeRoleRecord(context string, input interface{}) (record RoleRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	subrolesValue, ok := rpcStruct["subroles"]
	if ok && subrolesValue != nil {
		record.Subroles, err = deserializeRoleRefSet(fmt.Sprintf("%s.%s", context, "subroles"), subrolesValue)
		if err != nil {
			return
		}
	}
	isInternalValue, ok := rpcStruct["is_internal"]
	if ok && isInternalValue != nil {
		record.IsInternal, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_internal"), isInternalValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeRoleRef(context string, input interface{}) (RoleRef, error) {
	var ref RoleRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return RoleRef(value), nil
}

func deserializeRoleRefSet(context string, input interface{}) (slice []RoleRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]RoleRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeRoleRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeRoleRefToRoleRecordMap(context string, input interface{}) (goMap map[RoleRef]RoleRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[RoleRef]RoleRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeRoleRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeRoleRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeSDNControllerRecord(context string, input interface{}) (record SDNControllerRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	protocolValue, ok := rpcStruct["protocol"]
	if ok && protocolValue != nil {
		record.Protocol, err = deserializeEnumSdnControllerProtocol(fmt.Sprintf("%s.%s", context, "protocol"), protocolValue)
		if err != nil {
			return
		}
	}
	addressValue, ok := rpcStruct["address"]
	if ok && addressValue != nil {
		record.Address, err = deserializeString(fmt.Sprintf("%s.%s", context, "address"), addressValue)
		if err != nil {
			return
		}
	}
	portValue, ok := rpcStruct["port"]
	if ok && portValue != nil {
		record.Port, err = deserializeInt(fmt.Sprintf("%s.%s", context, "port"), portValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeSDNControllerRef(context string, input interface{}) (SDNControllerRef, error) {
	var ref SDNControllerRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return SDNControllerRef(value), nil
}

func deserializeSDNControllerRefSet(context string, input interface{}) (slice []SDNControllerRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]SDNControllerRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeSDNControllerRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeSDNControllerRefToSDNControllerRecordMap(context string, input interface{}) (goMap map[SDNControllerRef]SDNControllerRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[SDNControllerRef]SDNControllerRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeSDNControllerRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeSDNControllerRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeSMRecord(context string, input interface{}) (record SMRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	typeValue, ok := rpcStruct["type"]
	if ok && typeValue != nil {
		record.Type, err = deserializeString(fmt.Sprintf("%s.%s", context, "type"), typeValue)
		if err != nil {
			return
		}
	}
	vendorValue, ok := rpcStruct["vendor"]
	if ok && vendorValue != nil {
		record.Vendor, err = deserializeString(fmt.Sprintf("%s.%s", context, "vendor"), vendorValue)
		if err != nil {
			return
		}
	}
	copyrightValue, ok := rpcStruct["copyright"]
	if ok && copyrightValue != nil {
		record.Copyright, err = deserializeString(fmt.Sprintf("%s.%s", context, "copyright"), copyrightValue)
		if err != nil {
			return
		}
	}
	versionValue, ok := rpcStruct["version"]
	if ok && versionValue != nil {
		record.Version, err = deserializeString(fmt.Sprintf("%s.%s", context, "version"), versionValue)
		if err != nil {
			return
		}
	}
	requiredAPIVersionValue, ok := rpcStruct["required_api_version"]
	if ok && requiredAPIVersionValue != nil {
		record.RequiredAPIVersion, err = deserializeString(fmt.Sprintf("%s.%s", context, "required_api_version"), requiredAPIVersionValue)
		if err != nil {
			return
		}
	}
	configurationValue, ok := rpcStruct["configuration"]
	if ok && configurationValue != nil {
		record.Configuration, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "configuration"), configurationValue)
		if err != nil {
			return
		}
	}
	capabilitiesValue, ok := rpcStruct["capabilities"]
	if ok && capabilitiesValue != nil {
		record.Capabilities, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "capabilities"), capabilitiesValue)
		if err != nil {
			return
		}
	}
	featuresValue, ok := rpcStruct["features"]
	if ok && featuresValue != nil {
		record.Features, err = deserializeStringToIntMap(fmt.Sprintf("%s.%s", context, "features"), featuresValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	driverFilenameValue, ok := rpcStruct["driver_filename"]
	if ok && driverFilenameValue != nil {
		record.DriverFilename, err = deserializeString(fmt.Sprintf("%s.%s", context, "driver_filename"), driverFilenameValue)
		if err != nil {
			return
		}
	}
	requiredClusterStackValue, ok := rpcStruct["required_cluster_stack"]
	if ok && requiredClusterStackValue != nil {
		record.RequiredClusterStack, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "required_cluster_stack"), requiredClusterStackValue)
		if err != nil {
			return
		}
	}
	supportedImageFormatsValue, ok := rpcStruct["supported_image_formats"]
	if ok && supportedImageFormatsValue != nil {
		record.SupportedImageFormats, err = deserializeEnumImageFormatTypeSet(fmt.Sprintf("%s.%s", context, "supported_image_formats"), supportedImageFormatsValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeSMRef(context string, input interface{}) (SMRef, error) {
	var ref SMRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return SMRef(value), nil
}

func deserializeSMRefSet(context string, input interface{}) (slice []SMRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]SMRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeSMRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeSMRefToSMRecordMap(context string, input interface{}) (goMap map[SMRef]SMRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[SMRef]SMRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeSMRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeSMRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeSRRecord(context string, input interface{}) (record SRRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumStorageOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumStorageOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	vDIsValue, ok := rpcStruct["VDIs"]
	if ok && vDIsValue != nil {
		record.VDIs, err = deserializeVDIRefSet(fmt.Sprintf("%s.%s", context, "VDIs"), vDIsValue)
		if err != nil {
			return
		}
	}
	pBDsValue, ok := rpcStruct["PBDs"]
	if ok && pBDsValue != nil {
		record.PBDs, err = deserializePBDRefSet(fmt.Sprintf("%s.%s", context, "PBDs"), pBDsValue)
		if err != nil {
			return
		}
	}
	virtualAllocationValue, ok := rpcStruct["virtual_allocation"]
	if ok && virtualAllocationValue != nil {
		record.VirtualAllocation, err = deserializeInt(fmt.Sprintf("%s.%s", context, "virtual_allocation"), virtualAllocationValue)
		if err != nil {
			return
		}
	}
	physicalUtilisationValue, ok := rpcStruct["physical_utilisation"]
	if ok && physicalUtilisationValue != nil {
		record.PhysicalUtilisation, err = deserializeInt(fmt.Sprintf("%s.%s", context, "physical_utilisation"), physicalUtilisationValue)
		if err != nil {
			return
		}
	}
	physicalSizeValue, ok := rpcStruct["physical_size"]
	if ok && physicalSizeValue != nil {
		record.PhysicalSize, err = deserializeInt(fmt.Sprintf("%s.%s", context, "physical_size"), physicalSizeValue)
		if err != nil {
			return
		}
	}
	typeValue, ok := rpcStruct["type"]
	if ok && typeValue != nil {
		record.Type, err = deserializeString(fmt.Sprintf("%s.%s", context, "type"), typeValue)
		if err != nil {
			return
		}
	}
	contentTypeValue, ok := rpcStruct["content_type"]
	if ok && contentTypeValue != nil {
		record.ContentType, err = deserializeString(fmt.Sprintf("%s.%s", context, "content_type"), contentTypeValue)
		if err != nil {
			return
		}
	}
	sharedValue, ok := rpcStruct["shared"]
	if ok && sharedValue != nil {
		record.Shared, err = deserializeBool(fmt.Sprintf("%s.%s", context, "shared"), sharedValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	tagsValue, ok := rpcStruct["tags"]
	if ok && tagsValue != nil {
		record.Tags, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "tags"), tagsValue)
		if err != nil {
			return
		}
	}
	smConfigValue, ok := rpcStruct["sm_config"]
	if ok && smConfigValue != nil {
		record.SmConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "sm_config"), smConfigValue)
		if err != nil {
			return
		}
	}
	blobsValue, ok := rpcStruct["blobs"]
	if ok && blobsValue != nil {
		record.Blobs, err = deserializeStringToBlobRefMap(fmt.Sprintf("%s.%s", context, "blobs"), blobsValue)
		if err != nil {
			return
		}
	}
	localCacheEnabledValue, ok := rpcStruct["local_cache_enabled"]
	if ok && localCacheEnabledValue != nil {
		record.LocalCacheEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "local_cache_enabled"), localCacheEnabledValue)
		if err != nil {
			return
		}
	}
	introducedByValue, ok := rpcStruct["introduced_by"]
	if ok && introducedByValue != nil {
		record.IntroducedBy, err = deserializeDRTaskRef(fmt.Sprintf("%s.%s", context, "introduced_by"), introducedByValue)
		if err != nil {
			return
		}
	}
	clusteredValue, ok := rpcStruct["clustered"]
	if ok && clusteredValue != nil {
		record.Clustered, err = deserializeBool(fmt.Sprintf("%s.%s", context, "clustered"), clusteredValue)
		if err != nil {
			return
		}
	}
	isToolsSrValue, ok := rpcStruct["is_tools_sr"]
	if ok && isToolsSrValue != nil {
		record.IsToolsSr, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_tools_sr"), isToolsSrValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeSRRef(context string, input interface{}) (SRRef, error) {
	var ref SRRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return SRRef(value), nil
}

func deserializeSRRefSet(context string, input interface{}) (slice []SRRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]SRRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeSRRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeSRRefToSRRecordMap(context string, input interface{}) (goMap map[SRRef]SRRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[SRRef]SRRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeSRRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeSRRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeSecretRecord(context string, input interface{}) (record SecretRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	valueValue, ok := rpcStruct["value"]
	if ok && valueValue != nil {
		record.Value, err = deserializeString(fmt.Sprintf("%s.%s", context, "value"), valueValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeSecretRef(context string, input interface{}) (SecretRef, error) {
	var ref SecretRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return SecretRef(value), nil
}

func deserializeSecretRefSet(context string, input interface{}) (slice []SecretRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]SecretRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeSecretRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeSecretRefToSecretRecordMap(context string, input interface{}) (goMap map[SecretRef]SecretRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[SecretRef]SecretRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeSecretRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeSecretRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeSessionRecord(context string, input interface{}) (record SessionRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	thisHostValue, ok := rpcStruct["this_host"]
	if ok && thisHostValue != nil {
		record.ThisHost, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "this_host"), thisHostValue)
		if err != nil {
			return
		}
	}
	thisUserValue, ok := rpcStruct["this_user"]
	if ok && thisUserValue != nil {
		record.ThisUser, err = deserializeUserRef(fmt.Sprintf("%s.%s", context, "this_user"), thisUserValue)
		if err != nil {
			return
		}
	}
	lastActiveValue, ok := rpcStruct["last_active"]
	if ok && lastActiveValue != nil {
		record.LastActive, err = deserializeTime(fmt.Sprintf("%s.%s", context, "last_active"), lastActiveValue)
		if err != nil {
			return
		}
	}
	poolValue, ok := rpcStruct["pool"]
	if ok && poolValue != nil {
		record.Pool, err = deserializeBool(fmt.Sprintf("%s.%s", context, "pool"), poolValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	isLocalSuperuserValue, ok := rpcStruct["is_local_superuser"]
	if ok && isLocalSuperuserValue != nil {
		record.IsLocalSuperuser, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_local_superuser"), isLocalSuperuserValue)
		if err != nil {
			return
		}
	}
	subjectValue, ok := rpcStruct["subject"]
	if ok && subjectValue != nil {
		record.Subject, err = deserializeSubjectRef(fmt.Sprintf("%s.%s", context, "subject"), subjectValue)
		if err != nil {
			return
		}
	}
	validationTimeValue, ok := rpcStruct["validation_time"]
	if ok && validationTimeValue != nil {
		record.ValidationTime, err = deserializeTime(fmt.Sprintf("%s.%s", context, "validation_time"), validationTimeValue)
		if err != nil {
			return
		}
	}
	authUserSidValue, ok := rpcStruct["auth_user_sid"]
	if ok && authUserSidValue != nil {
		record.AuthUserSid, err = deserializeString(fmt.Sprintf("%s.%s", context, "auth_user_sid"), authUserSidValue)
		if err != nil {
			return
		}
	}
	authUserNameValue, ok := rpcStruct["auth_user_name"]
	if ok && authUserNameValue != nil {
		record.AuthUserName, err = deserializeString(fmt.Sprintf("%s.%s", context, "auth_user_name"), authUserNameValue)
		if err != nil {
			return
		}
	}
	rbacPermissionsValue, ok := rpcStruct["rbac_permissions"]
	if ok && rbacPermissionsValue != nil {
		record.RbacPermissions, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "rbac_permissions"), rbacPermissionsValue)
		if err != nil {
			return
		}
	}
	tasksValue, ok := rpcStruct["tasks"]
	if ok && tasksValue != nil {
		record.Tasks, err = deserializeTaskRefSet(fmt.Sprintf("%s.%s", context, "tasks"), tasksValue)
		if err != nil {
			return
		}
	}
	parentValue, ok := rpcStruct["parent"]
	if ok && parentValue != nil {
		record.Parent, err = deserializeSessionRef(fmt.Sprintf("%s.%s", context, "parent"), parentValue)
		if err != nil {
			return
		}
	}
	originatorValue, ok := rpcStruct["originator"]
	if ok && originatorValue != nil {
		record.Originator, err = deserializeString(fmt.Sprintf("%s.%s", context, "originator"), originatorValue)
		if err != nil {
			return
		}
	}
	clientCertificateValue, ok := rpcStruct["client_certificate"]
	if ok && clientCertificateValue != nil {
		record.ClientCertificate, err = deserializeBool(fmt.Sprintf("%s.%s", context, "client_certificate"), clientCertificateValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeSessionRef(context string, input interface{}) (SessionRef, error) {
	var ref SessionRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return SessionRef(value), nil
}

func deserializeSrStatRecord(context string, input interface{}) (record SrStatRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	record.UUID, err = deserializeOptionString(fmt.Sprintf("%s.%s", context, "uuid"), rpcStruct["uuid"])
	if err != nil {
		return
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	freeSpaceValue, ok := rpcStruct["free_space"]
	if ok && freeSpaceValue != nil {
		record.FreeSpace, err = deserializeInt(fmt.Sprintf("%s.%s", context, "free_space"), freeSpaceValue)
		if err != nil {
			return
		}
	}
	totalSpaceValue, ok := rpcStruct["total_space"]
	if ok && totalSpaceValue != nil {
		record.TotalSpace, err = deserializeInt(fmt.Sprintf("%s.%s", context, "total_space"), totalSpaceValue)
		if err != nil {
			return
		}
	}
	clusteredValue, ok := rpcStruct["clustered"]
	if ok && clusteredValue != nil {
		record.Clustered, err = deserializeBool(fmt.Sprintf("%s.%s", context, "clustered"), clusteredValue)
		if err != nil {
			return
		}
	}
	healthValue, ok := rpcStruct["health"]
	if ok && healthValue != nil {
		record.Health, err = deserializeEnumSrHealth(fmt.Sprintf("%s.%s", context, "health"), healthValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeString(context string, input interface{}) (value string, err error) {
	if input == nil {
		return
	}
	value, ok := input.(string)
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return
}

func deserializeStringSet(context string, input interface{}) (slice []string, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]string, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeString(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeStringSetSet(context string, input interface{}) (slice [][]string, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([][]string, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeStringSet(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeStringToBlobRefMap(context string, input interface{}) (goMap map[string]BlobRef, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]BlobRef, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeBlobRef(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumClusterHostOperationMap(context string, input interface{}) (goMap map[string]ClusterHostOperation, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]ClusterHostOperation, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumClusterHostOperation(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumClusterOperationMap(context string, input interface{}) (goMap map[string]ClusterOperation, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]ClusterOperation, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumClusterOperation(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumHostAllowedOperationsMap(context string, input interface{}) (goMap map[string]HostAllowedOperations, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]HostAllowedOperations, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumHostAllowedOperations(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumNetworkOperationsMap(context string, input interface{}) (goMap map[string]NetworkOperations, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]NetworkOperations, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumNetworkOperations(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumPoolAllowedOperationsMap(context string, input interface{}) (goMap map[string]PoolAllowedOperations, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]PoolAllowedOperations, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumPoolAllowedOperations(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumStorageOperationsMap(context string, input interface{}) (goMap map[string]StorageOperations, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]StorageOperations, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumStorageOperations(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumTaskAllowedOperationsMap(context string, input interface{}) (goMap map[string]TaskAllowedOperations, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]TaskAllowedOperations, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumTaskAllowedOperations(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumVMApplianceOperationMap(context string, input interface{}) (goMap map[string]VMApplianceOperation, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]VMApplianceOperation, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumVMApplianceOperation(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumVMOperationsMap(context string, input interface{}) (goMap map[string]VMOperations, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]VMOperations, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumVMOperations(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumVbdOperationsMap(context string, input interface{}) (goMap map[string]VbdOperations, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]VbdOperations, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumVbdOperations(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumVdiOperationsMap(context string, input interface{}) (goMap map[string]VdiOperations, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]VdiOperations, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumVdiOperations(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumVifOperationsMap(context string, input interface{}) (goMap map[string]VifOperations, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]VifOperations, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumVifOperations(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumVtpmOperationsMap(context string, input interface{}) (goMap map[string]VtpmOperations, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]VtpmOperations, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumVtpmOperations(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToEnumVusbOperationsMap(context string, input interface{}) (goMap map[string]VusbOperations, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]VusbOperations, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeEnumVusbOperations(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToIntMap(context string, input interface{}) (goMap map[string]int, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]int, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeInt(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeStringToStringMap(context string, input interface{}) (goMap map[string]string, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[string]string, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeString(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeString(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeSubjectRecord(context string, input interface{}) (record SubjectRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	subjectIdentifierValue, ok := rpcStruct["subject_identifier"]
	if ok && subjectIdentifierValue != nil {
		record.SubjectIdentifier, err = deserializeString(fmt.Sprintf("%s.%s", context, "subject_identifier"), subjectIdentifierValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	rolesValue, ok := rpcStruct["roles"]
	if ok && rolesValue != nil {
		record.Roles, err = deserializeRoleRefSet(fmt.Sprintf("%s.%s", context, "roles"), rolesValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeSubjectRef(context string, input interface{}) (SubjectRef, error) {
	var ref SubjectRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return SubjectRef(value), nil
}

func deserializeSubjectRefSet(context string, input interface{}) (slice []SubjectRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]SubjectRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeSubjectRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeSubjectRefToSubjectRecordMap(context string, input interface{}) (goMap map[SubjectRef]SubjectRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[SubjectRef]SubjectRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeSubjectRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeSubjectRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeTaskRecord(context string, input interface{}) (record TaskRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumTaskAllowedOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumTaskAllowedOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	createdValue, ok := rpcStruct["created"]
	if ok && createdValue != nil {
		record.Created, err = deserializeTime(fmt.Sprintf("%s.%s", context, "created"), createdValue)
		if err != nil {
			return
		}
	}
	finishedValue, ok := rpcStruct["finished"]
	if ok && finishedValue != nil {
		record.Finished, err = deserializeTime(fmt.Sprintf("%s.%s", context, "finished"), finishedValue)
		if err != nil {
			return
		}
	}
	statusValue, ok := rpcStruct["status"]
	if ok && statusValue != nil {
		record.Status, err = deserializeEnumTaskStatusType(fmt.Sprintf("%s.%s", context, "status"), statusValue)
		if err != nil {
			return
		}
	}
	residentOnValue, ok := rpcStruct["resident_on"]
	if ok && residentOnValue != nil {
		record.ResidentOn, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "resident_on"), residentOnValue)
		if err != nil {
			return
		}
	}
	progressValue, ok := rpcStruct["progress"]
	if ok && progressValue != nil {
		record.Progress, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "progress"), progressValue)
		if err != nil {
			return
		}
	}
	typeValue, ok := rpcStruct["type"]
	if ok && typeValue != nil {
		record.Type, err = deserializeString(fmt.Sprintf("%s.%s", context, "type"), typeValue)
		if err != nil {
			return
		}
	}
	resultValue, ok := rpcStruct["result"]
	if ok && resultValue != nil {
		record.Result, err = deserializeString(fmt.Sprintf("%s.%s", context, "result"), resultValue)
		if err != nil {
			return
		}
	}
	errorInfoValue, ok := rpcStruct["error_info"]
	if ok && errorInfoValue != nil {
		record.ErrorInfo, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "error_info"), errorInfoValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	subtaskOfValue, ok := rpcStruct["subtask_of"]
	if ok && subtaskOfValue != nil {
		record.SubtaskOf, err = deserializeTaskRef(fmt.Sprintf("%s.%s", context, "subtask_of"), subtaskOfValue)
		if err != nil {
			return
		}
	}
	subtasksValue, ok := rpcStruct["subtasks"]
	if ok && subtasksValue != nil {
		record.Subtasks, err = deserializeTaskRefSet(fmt.Sprintf("%s.%s", context, "subtasks"), subtasksValue)
		if err != nil {
			return
		}
	}
	backtraceValue, ok := rpcStruct["backtrace"]
	if ok && backtraceValue != nil {
		record.Backtrace, err = deserializeString(fmt.Sprintf("%s.%s", context, "backtrace"), backtraceValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeTaskRef(context string, input interface{}) (TaskRef, error) {
	var ref TaskRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return TaskRef(value), nil
}

func deserializeTaskRefSet(context string, input interface{}) (slice []TaskRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]TaskRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeTaskRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeTaskRefToTaskRecordMap(context string, input interface{}) (goMap map[TaskRef]TaskRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[TaskRef]TaskRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeTaskRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeTaskRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeTime(context string, input interface{}) (value time.Time, err error) {
	_ = context
	if input == nil {
		return
	}
	strValue := fmt.Sprintf("%v", input)
	floatValue, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		for _, timeFormat := range timeFormats {
			value, err = time.Parse(timeFormat, strValue)
			if err == nil {
				return value, nil
			}
		}
		return
	}
	unixTimestamp, err := strconv.ParseInt(strconv.Itoa(int(floatValue)), 10, 64)
	value = time.Unix(unixTimestamp, 0).UTC()

	return
}

func deserializeTunnelRecord(context string, input interface{}) (record TunnelRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	accessPIFValue, ok := rpcStruct["access_PIF"]
	if ok && accessPIFValue != nil {
		record.AccessPIF, err = deserializePIFRef(fmt.Sprintf("%s.%s", context, "access_PIF"), accessPIFValue)
		if err != nil {
			return
		}
	}
	transportPIFValue, ok := rpcStruct["transport_PIF"]
	if ok && transportPIFValue != nil {
		record.TransportPIF, err = deserializePIFRef(fmt.Sprintf("%s.%s", context, "transport_PIF"), transportPIFValue)
		if err != nil {
			return
		}
	}
	statusValue, ok := rpcStruct["status"]
	if ok && statusValue != nil {
		record.Status, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "status"), statusValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	protocolValue, ok := rpcStruct["protocol"]
	if ok && protocolValue != nil {
		record.Protocol, err = deserializeEnumTunnelProtocol(fmt.Sprintf("%s.%s", context, "protocol"), protocolValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeTunnelRef(context string, input interface{}) (TunnelRef, error) {
	var ref TunnelRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return TunnelRef(value), nil
}

func deserializeTunnelRefSet(context string, input interface{}) (slice []TunnelRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]TunnelRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeTunnelRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeTunnelRefToTunnelRecordMap(context string, input interface{}) (goMap map[TunnelRef]TunnelRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[TunnelRef]TunnelRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeTunnelRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeTunnelRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeUSBGroupRecord(context string, input interface{}) (record USBGroupRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	pUSBsValue, ok := rpcStruct["PUSBs"]
	if ok && pUSBsValue != nil {
		record.PUSBs, err = deserializePUSBRefSet(fmt.Sprintf("%s.%s", context, "PUSBs"), pUSBsValue)
		if err != nil {
			return
		}
	}
	vUSBsValue, ok := rpcStruct["VUSBs"]
	if ok && vUSBsValue != nil {
		record.VUSBs, err = deserializeVUSBRefSet(fmt.Sprintf("%s.%s", context, "VUSBs"), vUSBsValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeUSBGroupRef(context string, input interface{}) (USBGroupRef, error) {
	var ref USBGroupRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return USBGroupRef(value), nil
}

func deserializeUSBGroupRefSet(context string, input interface{}) (slice []USBGroupRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]USBGroupRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeUSBGroupRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeUSBGroupRefToUSBGroupRecordMap(context string, input interface{}) (goMap map[USBGroupRef]USBGroupRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[USBGroupRef]USBGroupRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeUSBGroupRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeUSBGroupRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeUserRecord(context string, input interface{}) (record UserRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	shortNameValue, ok := rpcStruct["short_name"]
	if ok && shortNameValue != nil {
		record.ShortName, err = deserializeString(fmt.Sprintf("%s.%s", context, "short_name"), shortNameValue)
		if err != nil {
			return
		}
	}
	fullnameValue, ok := rpcStruct["fullname"]
	if ok && fullnameValue != nil {
		record.Fullname, err = deserializeString(fmt.Sprintf("%s.%s", context, "fullname"), fullnameValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeUserRef(context string, input interface{}) (UserRef, error) {
	var ref UserRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return UserRef(value), nil
}

func deserializeVBDMetricsRecord(context string, input interface{}) (record VBDMetricsRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	ioReadKbsValue, ok := rpcStruct["io_read_kbs"]
	if ok && ioReadKbsValue != nil {
		record.IoReadKbs, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "io_read_kbs"), ioReadKbsValue)
		if err != nil {
			return
		}
	}
	ioWriteKbsValue, ok := rpcStruct["io_write_kbs"]
	if ok && ioWriteKbsValue != nil {
		record.IoWriteKbs, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "io_write_kbs"), ioWriteKbsValue)
		if err != nil {
			return
		}
	}
	lastUpdatedValue, ok := rpcStruct["last_updated"]
	if ok && lastUpdatedValue != nil {
		record.LastUpdated, err = deserializeTime(fmt.Sprintf("%s.%s", context, "last_updated"), lastUpdatedValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVBDMetricsRef(context string, input interface{}) (VBDMetricsRef, error) {
	var ref VBDMetricsRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VBDMetricsRef(value), nil
}

func deserializeVBDMetricsRefSet(context string, input interface{}) (slice []VBDMetricsRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VBDMetricsRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVBDMetricsRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVBDMetricsRefToVBDMetricsRecordMap(context string, input interface{}) (goMap map[VBDMetricsRef]VBDMetricsRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VBDMetricsRef]VBDMetricsRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVBDMetricsRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVBDMetricsRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVBDRecord(context string, input interface{}) (record VBDRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumVbdOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumVbdOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	vMValue, ok := rpcStruct["VM"]
	if ok && vMValue != nil {
		record.VM, err = deserializeVMRef(fmt.Sprintf("%s.%s", context, "VM"), vMValue)
		if err != nil {
			return
		}
	}
	vDIValue, ok := rpcStruct["VDI"]
	if ok && vDIValue != nil {
		record.VDI, err = deserializeVDIRef(fmt.Sprintf("%s.%s", context, "VDI"), vDIValue)
		if err != nil {
			return
		}
	}
	deviceValue, ok := rpcStruct["device"]
	if ok && deviceValue != nil {
		record.Device, err = deserializeString(fmt.Sprintf("%s.%s", context, "device"), deviceValue)
		if err != nil {
			return
		}
	}
	userdeviceValue, ok := rpcStruct["userdevice"]
	if ok && userdeviceValue != nil {
		record.Userdevice, err = deserializeString(fmt.Sprintf("%s.%s", context, "userdevice"), userdeviceValue)
		if err != nil {
			return
		}
	}
	bootableValue, ok := rpcStruct["bootable"]
	if ok && bootableValue != nil {
		record.Bootable, err = deserializeBool(fmt.Sprintf("%s.%s", context, "bootable"), bootableValue)
		if err != nil {
			return
		}
	}
	modeValue, ok := rpcStruct["mode"]
	if ok && modeValue != nil {
		record.Mode, err = deserializeEnumVbdMode(fmt.Sprintf("%s.%s", context, "mode"), modeValue)
		if err != nil {
			return
		}
	}
	typeValue, ok := rpcStruct["type"]
	if ok && typeValue != nil {
		record.Type, err = deserializeEnumVbdType(fmt.Sprintf("%s.%s", context, "type"), typeValue)
		if err != nil {
			return
		}
	}
	unpluggableValue, ok := rpcStruct["unpluggable"]
	if ok && unpluggableValue != nil {
		record.Unpluggable, err = deserializeBool(fmt.Sprintf("%s.%s", context, "unpluggable"), unpluggableValue)
		if err != nil {
			return
		}
	}
	storageLockValue, ok := rpcStruct["storage_lock"]
	if ok && storageLockValue != nil {
		record.StorageLock, err = deserializeBool(fmt.Sprintf("%s.%s", context, "storage_lock"), storageLockValue)
		if err != nil {
			return
		}
	}
	emptyValue, ok := rpcStruct["empty"]
	if ok && emptyValue != nil {
		record.Empty, err = deserializeBool(fmt.Sprintf("%s.%s", context, "empty"), emptyValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	currentlyAttachedValue, ok := rpcStruct["currently_attached"]
	if ok && currentlyAttachedValue != nil {
		record.CurrentlyAttached, err = deserializeBool(fmt.Sprintf("%s.%s", context, "currently_attached"), currentlyAttachedValue)
		if err != nil {
			return
		}
	}
	statusCodeValue, ok := rpcStruct["status_code"]
	if ok && statusCodeValue != nil {
		record.StatusCode, err = deserializeInt(fmt.Sprintf("%s.%s", context, "status_code"), statusCodeValue)
		if err != nil {
			return
		}
	}
	statusDetailValue, ok := rpcStruct["status_detail"]
	if ok && statusDetailValue != nil {
		record.StatusDetail, err = deserializeString(fmt.Sprintf("%s.%s", context, "status_detail"), statusDetailValue)
		if err != nil {
			return
		}
	}
	runtimePropertiesValue, ok := rpcStruct["runtime_properties"]
	if ok && runtimePropertiesValue != nil {
		record.RuntimeProperties, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "runtime_properties"), runtimePropertiesValue)
		if err != nil {
			return
		}
	}
	qosAlgorithmTypeValue, ok := rpcStruct["qos_algorithm_type"]
	if ok && qosAlgorithmTypeValue != nil {
		record.QosAlgorithmType, err = deserializeString(fmt.Sprintf("%s.%s", context, "qos_algorithm_type"), qosAlgorithmTypeValue)
		if err != nil {
			return
		}
	}
	qosAlgorithmParamsValue, ok := rpcStruct["qos_algorithm_params"]
	if ok && qosAlgorithmParamsValue != nil {
		record.QosAlgorithmParams, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "qos_algorithm_params"), qosAlgorithmParamsValue)
		if err != nil {
			return
		}
	}
	qosSupportedAlgorithmsValue, ok := rpcStruct["qos_supported_algorithms"]
	if ok && qosSupportedAlgorithmsValue != nil {
		record.QosSupportedAlgorithms, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "qos_supported_algorithms"), qosSupportedAlgorithmsValue)
		if err != nil {
			return
		}
	}
	metricsValue, ok := rpcStruct["metrics"]
	if ok && metricsValue != nil {
		record.Metrics, err = deserializeVBDMetricsRef(fmt.Sprintf("%s.%s", context, "metrics"), metricsValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVBDRef(context string, input interface{}) (VBDRef, error) {
	var ref VBDRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VBDRef(value), nil
}

func deserializeVBDRefSet(context string, input interface{}) (slice []VBDRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VBDRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVBDRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVBDRefToVBDRecordMap(context string, input interface{}) (goMap map[VBDRef]VBDRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VBDRef]VBDRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVBDRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVBDRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVDIRecord(context string, input interface{}) (record VDIRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumVdiOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumVdiOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	sRValue, ok := rpcStruct["SR"]
	if ok && sRValue != nil {
		record.SR, err = deserializeSRRef(fmt.Sprintf("%s.%s", context, "SR"), sRValue)
		if err != nil {
			return
		}
	}
	vBDsValue, ok := rpcStruct["VBDs"]
	if ok && vBDsValue != nil {
		record.VBDs, err = deserializeVBDRefSet(fmt.Sprintf("%s.%s", context, "VBDs"), vBDsValue)
		if err != nil {
			return
		}
	}
	crashDumpsValue, ok := rpcStruct["crash_dumps"]
	if ok && crashDumpsValue != nil {
		record.CrashDumps, err = deserializeCrashdumpRefSet(fmt.Sprintf("%s.%s", context, "crash_dumps"), crashDumpsValue)
		if err != nil {
			return
		}
	}
	virtualSizeValue, ok := rpcStruct["virtual_size"]
	if ok && virtualSizeValue != nil {
		record.VirtualSize, err = deserializeInt(fmt.Sprintf("%s.%s", context, "virtual_size"), virtualSizeValue)
		if err != nil {
			return
		}
	}
	physicalUtilisationValue, ok := rpcStruct["physical_utilisation"]
	if ok && physicalUtilisationValue != nil {
		record.PhysicalUtilisation, err = deserializeInt(fmt.Sprintf("%s.%s", context, "physical_utilisation"), physicalUtilisationValue)
		if err != nil {
			return
		}
	}
	typeValue, ok := rpcStruct["type"]
	if ok && typeValue != nil {
		record.Type, err = deserializeEnumVdiType(fmt.Sprintf("%s.%s", context, "type"), typeValue)
		if err != nil {
			return
		}
	}
	sharableValue, ok := rpcStruct["sharable"]
	if ok && sharableValue != nil {
		record.Sharable, err = deserializeBool(fmt.Sprintf("%s.%s", context, "sharable"), sharableValue)
		if err != nil {
			return
		}
	}
	readOnlyValue, ok := rpcStruct["read_only"]
	if ok && readOnlyValue != nil {
		record.ReadOnly, err = deserializeBool(fmt.Sprintf("%s.%s", context, "read_only"), readOnlyValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	storageLockValue, ok := rpcStruct["storage_lock"]
	if ok && storageLockValue != nil {
		record.StorageLock, err = deserializeBool(fmt.Sprintf("%s.%s", context, "storage_lock"), storageLockValue)
		if err != nil {
			return
		}
	}
	locationValue, ok := rpcStruct["location"]
	if ok && locationValue != nil {
		record.Location, err = deserializeString(fmt.Sprintf("%s.%s", context, "location"), locationValue)
		if err != nil {
			return
		}
	}
	managedValue, ok := rpcStruct["managed"]
	if ok && managedValue != nil {
		record.Managed, err = deserializeBool(fmt.Sprintf("%s.%s", context, "managed"), managedValue)
		if err != nil {
			return
		}
	}
	missingValue, ok := rpcStruct["missing"]
	if ok && missingValue != nil {
		record.Missing, err = deserializeBool(fmt.Sprintf("%s.%s", context, "missing"), missingValue)
		if err != nil {
			return
		}
	}
	parentValue, ok := rpcStruct["parent"]
	if ok && parentValue != nil {
		record.Parent, err = deserializeVDIRef(fmt.Sprintf("%s.%s", context, "parent"), parentValue)
		if err != nil {
			return
		}
	}
	xenstoreDataValue, ok := rpcStruct["xenstore_data"]
	if ok && xenstoreDataValue != nil {
		record.XenstoreData, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "xenstore_data"), xenstoreDataValue)
		if err != nil {
			return
		}
	}
	smConfigValue, ok := rpcStruct["sm_config"]
	if ok && smConfigValue != nil {
		record.SmConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "sm_config"), smConfigValue)
		if err != nil {
			return
		}
	}
	isASnapshotValue, ok := rpcStruct["is_a_snapshot"]
	if ok && isASnapshotValue != nil {
		record.IsASnapshot, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_a_snapshot"), isASnapshotValue)
		if err != nil {
			return
		}
	}
	snapshotOfValue, ok := rpcStruct["snapshot_of"]
	if ok && snapshotOfValue != nil {
		record.SnapshotOf, err = deserializeVDIRef(fmt.Sprintf("%s.%s", context, "snapshot_of"), snapshotOfValue)
		if err != nil {
			return
		}
	}
	snapshotsValue, ok := rpcStruct["snapshots"]
	if ok && snapshotsValue != nil {
		record.Snapshots, err = deserializeVDIRefSet(fmt.Sprintf("%s.%s", context, "snapshots"), snapshotsValue)
		if err != nil {
			return
		}
	}
	snapshotTimeValue, ok := rpcStruct["snapshot_time"]
	if ok && snapshotTimeValue != nil {
		record.SnapshotTime, err = deserializeTime(fmt.Sprintf("%s.%s", context, "snapshot_time"), snapshotTimeValue)
		if err != nil {
			return
		}
	}
	tagsValue, ok := rpcStruct["tags"]
	if ok && tagsValue != nil {
		record.Tags, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "tags"), tagsValue)
		if err != nil {
			return
		}
	}
	allowCachingValue, ok := rpcStruct["allow_caching"]
	if ok && allowCachingValue != nil {
		record.AllowCaching, err = deserializeBool(fmt.Sprintf("%s.%s", context, "allow_caching"), allowCachingValue)
		if err != nil {
			return
		}
	}
	onBootValue, ok := rpcStruct["on_boot"]
	if ok && onBootValue != nil {
		record.OnBoot, err = deserializeEnumOnBoot(fmt.Sprintf("%s.%s", context, "on_boot"), onBootValue)
		if err != nil {
			return
		}
	}
	metadataOfPoolValue, ok := rpcStruct["metadata_of_pool"]
	if ok && metadataOfPoolValue != nil {
		record.MetadataOfPool, err = deserializePoolRef(fmt.Sprintf("%s.%s", context, "metadata_of_pool"), metadataOfPoolValue)
		if err != nil {
			return
		}
	}
	metadataLatestValue, ok := rpcStruct["metadata_latest"]
	if ok && metadataLatestValue != nil {
		record.MetadataLatest, err = deserializeBool(fmt.Sprintf("%s.%s", context, "metadata_latest"), metadataLatestValue)
		if err != nil {
			return
		}
	}
	isToolsIsoValue, ok := rpcStruct["is_tools_iso"]
	if ok && isToolsIsoValue != nil {
		record.IsToolsIso, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_tools_iso"), isToolsIsoValue)
		if err != nil {
			return
		}
	}
	cbtEnabledValue, ok := rpcStruct["cbt_enabled"]
	if ok && cbtEnabledValue != nil {
		record.CbtEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "cbt_enabled"), cbtEnabledValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVDIRef(context string, input interface{}) (VDIRef, error) {
	var ref VDIRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VDIRef(value), nil
}

func deserializeVDIRefSet(context string, input interface{}) (slice []VDIRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VDIRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVDIRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVDIRefToVDIRecordMap(context string, input interface{}) (goMap map[VDIRef]VDIRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VDIRef]VDIRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVDIRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVDIRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVGPURecord(context string, input interface{}) (record VGPURecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	vMValue, ok := rpcStruct["VM"]
	if ok && vMValue != nil {
		record.VM, err = deserializeVMRef(fmt.Sprintf("%s.%s", context, "VM"), vMValue)
		if err != nil {
			return
		}
	}
	gPUGroupValue, ok := rpcStruct["GPU_group"]
	if ok && gPUGroupValue != nil {
		record.GPUGroup, err = deserializeGPUGroupRef(fmt.Sprintf("%s.%s", context, "GPU_group"), gPUGroupValue)
		if err != nil {
			return
		}
	}
	deviceValue, ok := rpcStruct["device"]
	if ok && deviceValue != nil {
		record.Device, err = deserializeString(fmt.Sprintf("%s.%s", context, "device"), deviceValue)
		if err != nil {
			return
		}
	}
	currentlyAttachedValue, ok := rpcStruct["currently_attached"]
	if ok && currentlyAttachedValue != nil {
		record.CurrentlyAttached, err = deserializeBool(fmt.Sprintf("%s.%s", context, "currently_attached"), currentlyAttachedValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	typeValue, ok := rpcStruct["type"]
	if ok && typeValue != nil {
		record.Type, err = deserializeVGPUTypeRef(fmt.Sprintf("%s.%s", context, "type"), typeValue)
		if err != nil {
			return
		}
	}
	residentOnValue, ok := rpcStruct["resident_on"]
	if ok && residentOnValue != nil {
		record.ResidentOn, err = deserializePGPURef(fmt.Sprintf("%s.%s", context, "resident_on"), residentOnValue)
		if err != nil {
			return
		}
	}
	scheduledToBeResidentOnValue, ok := rpcStruct["scheduled_to_be_resident_on"]
	if ok && scheduledToBeResidentOnValue != nil {
		record.ScheduledToBeResidentOn, err = deserializePGPURef(fmt.Sprintf("%s.%s", context, "scheduled_to_be_resident_on"), scheduledToBeResidentOnValue)
		if err != nil {
			return
		}
	}
	compatibilityMetadataValue, ok := rpcStruct["compatibility_metadata"]
	if ok && compatibilityMetadataValue != nil {
		record.CompatibilityMetadata, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "compatibility_metadata"), compatibilityMetadataValue)
		if err != nil {
			return
		}
	}
	extraArgsValue, ok := rpcStruct["extra_args"]
	if ok && extraArgsValue != nil {
		record.ExtraArgs, err = deserializeString(fmt.Sprintf("%s.%s", context, "extra_args"), extraArgsValue)
		if err != nil {
			return
		}
	}
	pCIValue, ok := rpcStruct["PCI"]
	if ok && pCIValue != nil {
		record.PCI, err = deserializePCIRef(fmt.Sprintf("%s.%s", context, "PCI"), pCIValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVGPURef(context string, input interface{}) (VGPURef, error) {
	var ref VGPURef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VGPURef(value), nil
}

func deserializeVGPURefSet(context string, input interface{}) (slice []VGPURef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VGPURef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVGPURef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVGPURefToVGPURecordMap(context string, input interface{}) (goMap map[VGPURef]VGPURecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VGPURef]VGPURecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVGPURef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVGPURecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVGPUTypeRecord(context string, input interface{}) (record VGPUTypeRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	vendorNameValue, ok := rpcStruct["vendor_name"]
	if ok && vendorNameValue != nil {
		record.VendorName, err = deserializeString(fmt.Sprintf("%s.%s", context, "vendor_name"), vendorNameValue)
		if err != nil {
			return
		}
	}
	modelNameValue, ok := rpcStruct["model_name"]
	if ok && modelNameValue != nil {
		record.ModelName, err = deserializeString(fmt.Sprintf("%s.%s", context, "model_name"), modelNameValue)
		if err != nil {
			return
		}
	}
	framebufferSizeValue, ok := rpcStruct["framebuffer_size"]
	if ok && framebufferSizeValue != nil {
		record.FramebufferSize, err = deserializeInt(fmt.Sprintf("%s.%s", context, "framebuffer_size"), framebufferSizeValue)
		if err != nil {
			return
		}
	}
	maxHeadsValue, ok := rpcStruct["max_heads"]
	if ok && maxHeadsValue != nil {
		record.MaxHeads, err = deserializeInt(fmt.Sprintf("%s.%s", context, "max_heads"), maxHeadsValue)
		if err != nil {
			return
		}
	}
	maxResolutionXValue, ok := rpcStruct["max_resolution_x"]
	if ok && maxResolutionXValue != nil {
		record.MaxResolutionX, err = deserializeInt(fmt.Sprintf("%s.%s", context, "max_resolution_x"), maxResolutionXValue)
		if err != nil {
			return
		}
	}
	maxResolutionYValue, ok := rpcStruct["max_resolution_y"]
	if ok && maxResolutionYValue != nil {
		record.MaxResolutionY, err = deserializeInt(fmt.Sprintf("%s.%s", context, "max_resolution_y"), maxResolutionYValue)
		if err != nil {
			return
		}
	}
	supportedOnPGPUsValue, ok := rpcStruct["supported_on_PGPUs"]
	if ok && supportedOnPGPUsValue != nil {
		record.SupportedOnPGPUs, err = deserializePGPURefSet(fmt.Sprintf("%s.%s", context, "supported_on_PGPUs"), supportedOnPGPUsValue)
		if err != nil {
			return
		}
	}
	enabledOnPGPUsValue, ok := rpcStruct["enabled_on_PGPUs"]
	if ok && enabledOnPGPUsValue != nil {
		record.EnabledOnPGPUs, err = deserializePGPURefSet(fmt.Sprintf("%s.%s", context, "enabled_on_PGPUs"), enabledOnPGPUsValue)
		if err != nil {
			return
		}
	}
	vGPUsValue, ok := rpcStruct["VGPUs"]
	if ok && vGPUsValue != nil {
		record.VGPUs, err = deserializeVGPURefSet(fmt.Sprintf("%s.%s", context, "VGPUs"), vGPUsValue)
		if err != nil {
			return
		}
	}
	supportedOnGPUGroupsValue, ok := rpcStruct["supported_on_GPU_groups"]
	if ok && supportedOnGPUGroupsValue != nil {
		record.SupportedOnGPUGroups, err = deserializeGPUGroupRefSet(fmt.Sprintf("%s.%s", context, "supported_on_GPU_groups"), supportedOnGPUGroupsValue)
		if err != nil {
			return
		}
	}
	enabledOnGPUGroupsValue, ok := rpcStruct["enabled_on_GPU_groups"]
	if ok && enabledOnGPUGroupsValue != nil {
		record.EnabledOnGPUGroups, err = deserializeGPUGroupRefSet(fmt.Sprintf("%s.%s", context, "enabled_on_GPU_groups"), enabledOnGPUGroupsValue)
		if err != nil {
			return
		}
	}
	implementationValue, ok := rpcStruct["implementation"]
	if ok && implementationValue != nil {
		record.Implementation, err = deserializeEnumVgpuTypeImplementation(fmt.Sprintf("%s.%s", context, "implementation"), implementationValue)
		if err != nil {
			return
		}
	}
	identifierValue, ok := rpcStruct["identifier"]
	if ok && identifierValue != nil {
		record.Identifier, err = deserializeString(fmt.Sprintf("%s.%s", context, "identifier"), identifierValue)
		if err != nil {
			return
		}
	}
	experimentalValue, ok := rpcStruct["experimental"]
	if ok && experimentalValue != nil {
		record.Experimental, err = deserializeBool(fmt.Sprintf("%s.%s", context, "experimental"), experimentalValue)
		if err != nil {
			return
		}
	}
	compatibleTypesInVMValue, ok := rpcStruct["compatible_types_in_vm"]
	if ok && compatibleTypesInVMValue != nil {
		record.CompatibleTypesInVM, err = deserializeVGPUTypeRefSet(fmt.Sprintf("%s.%s", context, "compatible_types_in_vm"), compatibleTypesInVMValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVGPUTypeRef(context string, input interface{}) (VGPUTypeRef, error) {
	var ref VGPUTypeRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VGPUTypeRef(value), nil
}

func deserializeVGPUTypeRefSet(context string, input interface{}) (slice []VGPUTypeRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VGPUTypeRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVGPUTypeRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVGPUTypeRefToIntMap(context string, input interface{}) (goMap map[VGPUTypeRef]int, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VGPUTypeRef]int, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVGPUTypeRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeInt(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVGPUTypeRefToVGPUTypeRecordMap(context string, input interface{}) (goMap map[VGPUTypeRef]VGPUTypeRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VGPUTypeRef]VGPUTypeRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVGPUTypeRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVGPUTypeRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVIFMetricsRecord(context string, input interface{}) (record VIFMetricsRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	ioReadKbsValue, ok := rpcStruct["io_read_kbs"]
	if ok && ioReadKbsValue != nil {
		record.IoReadKbs, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "io_read_kbs"), ioReadKbsValue)
		if err != nil {
			return
		}
	}
	ioWriteKbsValue, ok := rpcStruct["io_write_kbs"]
	if ok && ioWriteKbsValue != nil {
		record.IoWriteKbs, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "io_write_kbs"), ioWriteKbsValue)
		if err != nil {
			return
		}
	}
	lastUpdatedValue, ok := rpcStruct["last_updated"]
	if ok && lastUpdatedValue != nil {
		record.LastUpdated, err = deserializeTime(fmt.Sprintf("%s.%s", context, "last_updated"), lastUpdatedValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVIFMetricsRef(context string, input interface{}) (VIFMetricsRef, error) {
	var ref VIFMetricsRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VIFMetricsRef(value), nil
}

func deserializeVIFMetricsRefSet(context string, input interface{}) (slice []VIFMetricsRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VIFMetricsRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVIFMetricsRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVIFMetricsRefToVIFMetricsRecordMap(context string, input interface{}) (goMap map[VIFMetricsRef]VIFMetricsRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VIFMetricsRef]VIFMetricsRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVIFMetricsRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVIFMetricsRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVIFRecord(context string, input interface{}) (record VIFRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumVifOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumVifOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	deviceValue, ok := rpcStruct["device"]
	if ok && deviceValue != nil {
		record.Device, err = deserializeString(fmt.Sprintf("%s.%s", context, "device"), deviceValue)
		if err != nil {
			return
		}
	}
	networkValue, ok := rpcStruct["network"]
	if ok && networkValue != nil {
		record.Network, err = deserializeNetworkRef(fmt.Sprintf("%s.%s", context, "network"), networkValue)
		if err != nil {
			return
		}
	}
	vMValue, ok := rpcStruct["VM"]
	if ok && vMValue != nil {
		record.VM, err = deserializeVMRef(fmt.Sprintf("%s.%s", context, "VM"), vMValue)
		if err != nil {
			return
		}
	}
	mACValue, ok := rpcStruct["MAC"]
	if ok && mACValue != nil {
		record.MAC, err = deserializeString(fmt.Sprintf("%s.%s", context, "MAC"), mACValue)
		if err != nil {
			return
		}
	}
	mTUValue, ok := rpcStruct["MTU"]
	if ok && mTUValue != nil {
		record.MTU, err = deserializeInt(fmt.Sprintf("%s.%s", context, "MTU"), mTUValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	currentlyAttachedValue, ok := rpcStruct["currently_attached"]
	if ok && currentlyAttachedValue != nil {
		record.CurrentlyAttached, err = deserializeBool(fmt.Sprintf("%s.%s", context, "currently_attached"), currentlyAttachedValue)
		if err != nil {
			return
		}
	}
	statusCodeValue, ok := rpcStruct["status_code"]
	if ok && statusCodeValue != nil {
		record.StatusCode, err = deserializeInt(fmt.Sprintf("%s.%s", context, "status_code"), statusCodeValue)
		if err != nil {
			return
		}
	}
	statusDetailValue, ok := rpcStruct["status_detail"]
	if ok && statusDetailValue != nil {
		record.StatusDetail, err = deserializeString(fmt.Sprintf("%s.%s", context, "status_detail"), statusDetailValue)
		if err != nil {
			return
		}
	}
	runtimePropertiesValue, ok := rpcStruct["runtime_properties"]
	if ok && runtimePropertiesValue != nil {
		record.RuntimeProperties, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "runtime_properties"), runtimePropertiesValue)
		if err != nil {
			return
		}
	}
	qosAlgorithmTypeValue, ok := rpcStruct["qos_algorithm_type"]
	if ok && qosAlgorithmTypeValue != nil {
		record.QosAlgorithmType, err = deserializeString(fmt.Sprintf("%s.%s", context, "qos_algorithm_type"), qosAlgorithmTypeValue)
		if err != nil {
			return
		}
	}
	qosAlgorithmParamsValue, ok := rpcStruct["qos_algorithm_params"]
	if ok && qosAlgorithmParamsValue != nil {
		record.QosAlgorithmParams, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "qos_algorithm_params"), qosAlgorithmParamsValue)
		if err != nil {
			return
		}
	}
	qosSupportedAlgorithmsValue, ok := rpcStruct["qos_supported_algorithms"]
	if ok && qosSupportedAlgorithmsValue != nil {
		record.QosSupportedAlgorithms, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "qos_supported_algorithms"), qosSupportedAlgorithmsValue)
		if err != nil {
			return
		}
	}
	metricsValue, ok := rpcStruct["metrics"]
	if ok && metricsValue != nil {
		record.Metrics, err = deserializeVIFMetricsRef(fmt.Sprintf("%s.%s", context, "metrics"), metricsValue)
		if err != nil {
			return
		}
	}
	mACAutogeneratedValue, ok := rpcStruct["MAC_autogenerated"]
	if ok && mACAutogeneratedValue != nil {
		record.MACAutogenerated, err = deserializeBool(fmt.Sprintf("%s.%s", context, "MAC_autogenerated"), mACAutogeneratedValue)
		if err != nil {
			return
		}
	}
	lockingModeValue, ok := rpcStruct["locking_mode"]
	if ok && lockingModeValue != nil {
		record.LockingMode, err = deserializeEnumVifLockingMode(fmt.Sprintf("%s.%s", context, "locking_mode"), lockingModeValue)
		if err != nil {
			return
		}
	}
	ipv4AllowedValue, ok := rpcStruct["ipv4_allowed"]
	if ok && ipv4AllowedValue != nil {
		record.Ipv4Allowed, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "ipv4_allowed"), ipv4AllowedValue)
		if err != nil {
			return
		}
	}
	ipv6AllowedValue, ok := rpcStruct["ipv6_allowed"]
	if ok && ipv6AllowedValue != nil {
		record.Ipv6Allowed, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "ipv6_allowed"), ipv6AllowedValue)
		if err != nil {
			return
		}
	}
	ipv4ConfigurationModeValue, ok := rpcStruct["ipv4_configuration_mode"]
	if ok && ipv4ConfigurationModeValue != nil {
		record.Ipv4ConfigurationMode, err = deserializeEnumVifIpv4ConfigurationMode(fmt.Sprintf("%s.%s", context, "ipv4_configuration_mode"), ipv4ConfigurationModeValue)
		if err != nil {
			return
		}
	}
	ipv4AddressesValue, ok := rpcStruct["ipv4_addresses"]
	if ok && ipv4AddressesValue != nil {
		record.Ipv4Addresses, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "ipv4_addresses"), ipv4AddressesValue)
		if err != nil {
			return
		}
	}
	ipv4GatewayValue, ok := rpcStruct["ipv4_gateway"]
	if ok && ipv4GatewayValue != nil {
		record.Ipv4Gateway, err = deserializeString(fmt.Sprintf("%s.%s", context, "ipv4_gateway"), ipv4GatewayValue)
		if err != nil {
			return
		}
	}
	ipv6ConfigurationModeValue, ok := rpcStruct["ipv6_configuration_mode"]
	if ok && ipv6ConfigurationModeValue != nil {
		record.Ipv6ConfigurationMode, err = deserializeEnumVifIpv6ConfigurationMode(fmt.Sprintf("%s.%s", context, "ipv6_configuration_mode"), ipv6ConfigurationModeValue)
		if err != nil {
			return
		}
	}
	ipv6AddressesValue, ok := rpcStruct["ipv6_addresses"]
	if ok && ipv6AddressesValue != nil {
		record.Ipv6Addresses, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "ipv6_addresses"), ipv6AddressesValue)
		if err != nil {
			return
		}
	}
	ipv6GatewayValue, ok := rpcStruct["ipv6_gateway"]
	if ok && ipv6GatewayValue != nil {
		record.Ipv6Gateway, err = deserializeString(fmt.Sprintf("%s.%s", context, "ipv6_gateway"), ipv6GatewayValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVIFRef(context string, input interface{}) (VIFRef, error) {
	var ref VIFRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VIFRef(value), nil
}

func deserializeVIFRefSet(context string, input interface{}) (slice []VIFRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VIFRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVIFRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVIFRefToStringMap(context string, input interface{}) (goMap map[VIFRef]string, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VIFRef]string, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVIFRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeString(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVIFRefToVIFRecordMap(context string, input interface{}) (goMap map[VIFRef]VIFRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VIFRef]VIFRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVIFRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVIFRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVLANRecord(context string, input interface{}) (record VLANRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	taggedPIFValue, ok := rpcStruct["tagged_PIF"]
	if ok && taggedPIFValue != nil {
		record.TaggedPIF, err = deserializePIFRef(fmt.Sprintf("%s.%s", context, "tagged_PIF"), taggedPIFValue)
		if err != nil {
			return
		}
	}
	untaggedPIFValue, ok := rpcStruct["untagged_PIF"]
	if ok && untaggedPIFValue != nil {
		record.UntaggedPIF, err = deserializePIFRef(fmt.Sprintf("%s.%s", context, "untagged_PIF"), untaggedPIFValue)
		if err != nil {
			return
		}
	}
	tagValue, ok := rpcStruct["tag"]
	if ok && tagValue != nil {
		record.Tag, err = deserializeInt(fmt.Sprintf("%s.%s", context, "tag"), tagValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVLANRef(context string, input interface{}) (VLANRef, error) {
	var ref VLANRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VLANRef(value), nil
}

func deserializeVLANRefSet(context string, input interface{}) (slice []VLANRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VLANRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVLANRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVLANRefToVLANRecordMap(context string, input interface{}) (goMap map[VLANRef]VLANRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VLANRef]VLANRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVLANRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVLANRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVMApplianceRecord(context string, input interface{}) (record VMApplianceRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumVMApplianceOperationSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumVMApplianceOperationMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	vMsValue, ok := rpcStruct["VMs"]
	if ok && vMsValue != nil {
		record.VMs, err = deserializeVMRefSet(fmt.Sprintf("%s.%s", context, "VMs"), vMsValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVMApplianceRef(context string, input interface{}) (VMApplianceRef, error) {
	var ref VMApplianceRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VMApplianceRef(value), nil
}

func deserializeVMApplianceRefSet(context string, input interface{}) (slice []VMApplianceRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VMApplianceRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVMApplianceRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVMApplianceRefToVMApplianceRecordMap(context string, input interface{}) (goMap map[VMApplianceRef]VMApplianceRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VMApplianceRef]VMApplianceRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVMApplianceRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVMApplianceRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVMGroupRecord(context string, input interface{}) (record VMGroupRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	placementValue, ok := rpcStruct["placement"]
	if ok && placementValue != nil {
		record.Placement, err = deserializeEnumPlacementPolicy(fmt.Sprintf("%s.%s", context, "placement"), placementValue)
		if err != nil {
			return
		}
	}
	vMsValue, ok := rpcStruct["VMs"]
	if ok && vMsValue != nil {
		record.VMs, err = deserializeVMRefSet(fmt.Sprintf("%s.%s", context, "VMs"), vMsValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVMGroupRef(context string, input interface{}) (VMGroupRef, error) {
	var ref VMGroupRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VMGroupRef(value), nil
}

func deserializeVMGroupRefSet(context string, input interface{}) (slice []VMGroupRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VMGroupRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVMGroupRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVMGroupRefToVMGroupRecordMap(context string, input interface{}) (goMap map[VMGroupRef]VMGroupRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VMGroupRef]VMGroupRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVMGroupRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVMGroupRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVMGuestMetricsRecord(context string, input interface{}) (record VMGuestMetricsRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	osVersionValue, ok := rpcStruct["os_version"]
	if ok && osVersionValue != nil {
		record.OsVersion, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "os_version"), osVersionValue)
		if err != nil {
			return
		}
	}
	netbiosNameValue, ok := rpcStruct["netbios_name"]
	if ok && netbiosNameValue != nil {
		record.NetbiosName, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "netbios_name"), netbiosNameValue)
		if err != nil {
			return
		}
	}
	pVDriversVersionValue, ok := rpcStruct["PV_drivers_version"]
	if ok && pVDriversVersionValue != nil {
		record.PVDriversVersion, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "PV_drivers_version"), pVDriversVersionValue)
		if err != nil {
			return
		}
	}
	pVDriversUpToDateValue, ok := rpcStruct["PV_drivers_up_to_date"]
	if ok && pVDriversUpToDateValue != nil {
		record.PVDriversUpToDate, err = deserializeBool(fmt.Sprintf("%s.%s", context, "PV_drivers_up_to_date"), pVDriversUpToDateValue)
		if err != nil {
			return
		}
	}
	memoryValue, ok := rpcStruct["memory"]
	if ok && memoryValue != nil {
		record.Memory, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "memory"), memoryValue)
		if err != nil {
			return
		}
	}
	disksValue, ok := rpcStruct["disks"]
	if ok && disksValue != nil {
		record.Disks, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "disks"), disksValue)
		if err != nil {
			return
		}
	}
	networksValue, ok := rpcStruct["networks"]
	if ok && networksValue != nil {
		record.Networks, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "networks"), networksValue)
		if err != nil {
			return
		}
	}
	otherValue, ok := rpcStruct["other"]
	if ok && otherValue != nil {
		record.Other, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other"), otherValue)
		if err != nil {
			return
		}
	}
	lastUpdatedValue, ok := rpcStruct["last_updated"]
	if ok && lastUpdatedValue != nil {
		record.LastUpdated, err = deserializeTime(fmt.Sprintf("%s.%s", context, "last_updated"), lastUpdatedValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	liveValue, ok := rpcStruct["live"]
	if ok && liveValue != nil {
		record.Live, err = deserializeBool(fmt.Sprintf("%s.%s", context, "live"), liveValue)
		if err != nil {
			return
		}
	}
	canUseHotplugVbdValue, ok := rpcStruct["can_use_hotplug_vbd"]
	if ok && canUseHotplugVbdValue != nil {
		record.CanUseHotplugVbd, err = deserializeEnumTristateType(fmt.Sprintf("%s.%s", context, "can_use_hotplug_vbd"), canUseHotplugVbdValue)
		if err != nil {
			return
		}
	}
	canUseHotplugVifValue, ok := rpcStruct["can_use_hotplug_vif"]
	if ok && canUseHotplugVifValue != nil {
		record.CanUseHotplugVif, err = deserializeEnumTristateType(fmt.Sprintf("%s.%s", context, "can_use_hotplug_vif"), canUseHotplugVifValue)
		if err != nil {
			return
		}
	}
	pVDriversDetectedValue, ok := rpcStruct["PV_drivers_detected"]
	if ok && pVDriversDetectedValue != nil {
		record.PVDriversDetected, err = deserializeBool(fmt.Sprintf("%s.%s", context, "PV_drivers_detected"), pVDriversDetectedValue)
		if err != nil {
			return
		}
	}
	servicesValue, ok := rpcStruct["services"]
	if ok && servicesValue != nil {
		record.Services, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "services"), servicesValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVMGuestMetricsRef(context string, input interface{}) (VMGuestMetricsRef, error) {
	var ref VMGuestMetricsRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VMGuestMetricsRef(value), nil
}

func deserializeVMGuestMetricsRefSet(context string, input interface{}) (slice []VMGuestMetricsRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VMGuestMetricsRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVMGuestMetricsRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVMGuestMetricsRefToVMGuestMetricsRecordMap(context string, input interface{}) (goMap map[VMGuestMetricsRef]VMGuestMetricsRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VMGuestMetricsRef]VMGuestMetricsRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVMGuestMetricsRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVMGuestMetricsRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVMMetricsRecord(context string, input interface{}) (record VMMetricsRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	memoryActualValue, ok := rpcStruct["memory_actual"]
	if ok && memoryActualValue != nil {
		record.MemoryActual, err = deserializeInt(fmt.Sprintf("%s.%s", context, "memory_actual"), memoryActualValue)
		if err != nil {
			return
		}
	}
	vCPUsNumberValue, ok := rpcStruct["VCPUs_number"]
	if ok && vCPUsNumberValue != nil {
		record.VCPUsNumber, err = deserializeInt(fmt.Sprintf("%s.%s", context, "VCPUs_number"), vCPUsNumberValue)
		if err != nil {
			return
		}
	}
	vCPUsUtilisationValue, ok := rpcStruct["VCPUs_utilisation"]
	if ok && vCPUsUtilisationValue != nil {
		record.VCPUsUtilisation, err = deserializeIntToFloatMap(fmt.Sprintf("%s.%s", context, "VCPUs_utilisation"), vCPUsUtilisationValue)
		if err != nil {
			return
		}
	}
	vCPUsCPUValue, ok := rpcStruct["VCPUs_CPU"]
	if ok && vCPUsCPUValue != nil {
		record.VCPUsCPU, err = deserializeIntToIntMap(fmt.Sprintf("%s.%s", context, "VCPUs_CPU"), vCPUsCPUValue)
		if err != nil {
			return
		}
	}
	vCPUsParamsValue, ok := rpcStruct["VCPUs_params"]
	if ok && vCPUsParamsValue != nil {
		record.VCPUsParams, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "VCPUs_params"), vCPUsParamsValue)
		if err != nil {
			return
		}
	}
	vCPUsFlagsValue, ok := rpcStruct["VCPUs_flags"]
	if ok && vCPUsFlagsValue != nil {
		record.VCPUsFlags, err = deserializeIntToStringSetMap(fmt.Sprintf("%s.%s", context, "VCPUs_flags"), vCPUsFlagsValue)
		if err != nil {
			return
		}
	}
	stateValue, ok := rpcStruct["state"]
	if ok && stateValue != nil {
		record.State, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "state"), stateValue)
		if err != nil {
			return
		}
	}
	startTimeValue, ok := rpcStruct["start_time"]
	if ok && startTimeValue != nil {
		record.StartTime, err = deserializeTime(fmt.Sprintf("%s.%s", context, "start_time"), startTimeValue)
		if err != nil {
			return
		}
	}
	installTimeValue, ok := rpcStruct["install_time"]
	if ok && installTimeValue != nil {
		record.InstallTime, err = deserializeTime(fmt.Sprintf("%s.%s", context, "install_time"), installTimeValue)
		if err != nil {
			return
		}
	}
	lastUpdatedValue, ok := rpcStruct["last_updated"]
	if ok && lastUpdatedValue != nil {
		record.LastUpdated, err = deserializeTime(fmt.Sprintf("%s.%s", context, "last_updated"), lastUpdatedValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	hvmValue, ok := rpcStruct["hvm"]
	if ok && hvmValue != nil {
		record.Hvm, err = deserializeBool(fmt.Sprintf("%s.%s", context, "hvm"), hvmValue)
		if err != nil {
			return
		}
	}
	nestedVirtValue, ok := rpcStruct["nested_virt"]
	if ok && nestedVirtValue != nil {
		record.NestedVirt, err = deserializeBool(fmt.Sprintf("%s.%s", context, "nested_virt"), nestedVirtValue)
		if err != nil {
			return
		}
	}
	nomigrateValue, ok := rpcStruct["nomigrate"]
	if ok && nomigrateValue != nil {
		record.Nomigrate, err = deserializeBool(fmt.Sprintf("%s.%s", context, "nomigrate"), nomigrateValue)
		if err != nil {
			return
		}
	}
	currentDomainTypeValue, ok := rpcStruct["current_domain_type"]
	if ok && currentDomainTypeValue != nil {
		record.CurrentDomainType, err = deserializeEnumDomainType(fmt.Sprintf("%s.%s", context, "current_domain_type"), currentDomainTypeValue)
		if err != nil {
			return
		}
	}
	numaOptimisedValue, ok := rpcStruct["numa_optimised"]
	if ok && numaOptimisedValue != nil {
		record.NumaOptimised, err = deserializeBool(fmt.Sprintf("%s.%s", context, "numa_optimised"), numaOptimisedValue)
		if err != nil {
			return
		}
	}
	numaNodesValue, ok := rpcStruct["numa_nodes"]
	if ok && numaNodesValue != nil {
		record.NumaNodes, err = deserializeInt(fmt.Sprintf("%s.%s", context, "numa_nodes"), numaNodesValue)
		if err != nil {
			return
		}
	}
	numaNodeMemoryValue, ok := rpcStruct["numa_node_memory"]
	if ok && numaNodeMemoryValue != nil {
		record.NumaNodeMemory, err = deserializeIntToIntMap(fmt.Sprintf("%s.%s", context, "numa_node_memory"), numaNodeMemoryValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVMMetricsRef(context string, input interface{}) (VMMetricsRef, error) {
	var ref VMMetricsRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VMMetricsRef(value), nil
}

func deserializeVMMetricsRefSet(context string, input interface{}) (slice []VMMetricsRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VMMetricsRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVMMetricsRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVMMetricsRefToVMMetricsRecordMap(context string, input interface{}) (goMap map[VMMetricsRef]VMMetricsRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VMMetricsRef]VMMetricsRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVMMetricsRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVMMetricsRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVMPPRecord(context string, input interface{}) (record VMPPRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	isPolicyEnabledValue, ok := rpcStruct["is_policy_enabled"]
	if ok && isPolicyEnabledValue != nil {
		record.IsPolicyEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_policy_enabled"), isPolicyEnabledValue)
		if err != nil {
			return
		}
	}
	backupTypeValue, ok := rpcStruct["backup_type"]
	if ok && backupTypeValue != nil {
		record.BackupType, err = deserializeEnumVmppBackupType(fmt.Sprintf("%s.%s", context, "backup_type"), backupTypeValue)
		if err != nil {
			return
		}
	}
	backupRetentionValueValue, ok := rpcStruct["backup_retention_value"]
	if ok && backupRetentionValueValue != nil {
		record.BackupRetentionValue, err = deserializeInt(fmt.Sprintf("%s.%s", context, "backup_retention_value"), backupRetentionValueValue)
		if err != nil {
			return
		}
	}
	backupFrequencyValue, ok := rpcStruct["backup_frequency"]
	if ok && backupFrequencyValue != nil {
		record.BackupFrequency, err = deserializeEnumVmppBackupFrequency(fmt.Sprintf("%s.%s", context, "backup_frequency"), backupFrequencyValue)
		if err != nil {
			return
		}
	}
	backupScheduleValue, ok := rpcStruct["backup_schedule"]
	if ok && backupScheduleValue != nil {
		record.BackupSchedule, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "backup_schedule"), backupScheduleValue)
		if err != nil {
			return
		}
	}
	isBackupRunningValue, ok := rpcStruct["is_backup_running"]
	if ok && isBackupRunningValue != nil {
		record.IsBackupRunning, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_backup_running"), isBackupRunningValue)
		if err != nil {
			return
		}
	}
	backupLastRunTimeValue, ok := rpcStruct["backup_last_run_time"]
	if ok && backupLastRunTimeValue != nil {
		record.BackupLastRunTime, err = deserializeTime(fmt.Sprintf("%s.%s", context, "backup_last_run_time"), backupLastRunTimeValue)
		if err != nil {
			return
		}
	}
	archiveTargetTypeValue, ok := rpcStruct["archive_target_type"]
	if ok && archiveTargetTypeValue != nil {
		record.ArchiveTargetType, err = deserializeEnumVmppArchiveTargetType(fmt.Sprintf("%s.%s", context, "archive_target_type"), archiveTargetTypeValue)
		if err != nil {
			return
		}
	}
	archiveTargetConfigValue, ok := rpcStruct["archive_target_config"]
	if ok && archiveTargetConfigValue != nil {
		record.ArchiveTargetConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "archive_target_config"), archiveTargetConfigValue)
		if err != nil {
			return
		}
	}
	archiveFrequencyValue, ok := rpcStruct["archive_frequency"]
	if ok && archiveFrequencyValue != nil {
		record.ArchiveFrequency, err = deserializeEnumVmppArchiveFrequency(fmt.Sprintf("%s.%s", context, "archive_frequency"), archiveFrequencyValue)
		if err != nil {
			return
		}
	}
	archiveScheduleValue, ok := rpcStruct["archive_schedule"]
	if ok && archiveScheduleValue != nil {
		record.ArchiveSchedule, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "archive_schedule"), archiveScheduleValue)
		if err != nil {
			return
		}
	}
	isArchiveRunningValue, ok := rpcStruct["is_archive_running"]
	if ok && isArchiveRunningValue != nil {
		record.IsArchiveRunning, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_archive_running"), isArchiveRunningValue)
		if err != nil {
			return
		}
	}
	archiveLastRunTimeValue, ok := rpcStruct["archive_last_run_time"]
	if ok && archiveLastRunTimeValue != nil {
		record.ArchiveLastRunTime, err = deserializeTime(fmt.Sprintf("%s.%s", context, "archive_last_run_time"), archiveLastRunTimeValue)
		if err != nil {
			return
		}
	}
	vMsValue, ok := rpcStruct["VMs"]
	if ok && vMsValue != nil {
		record.VMs, err = deserializeVMRefSet(fmt.Sprintf("%s.%s", context, "VMs"), vMsValue)
		if err != nil {
			return
		}
	}
	isAlarmEnabledValue, ok := rpcStruct["is_alarm_enabled"]
	if ok && isAlarmEnabledValue != nil {
		record.IsAlarmEnabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_alarm_enabled"), isAlarmEnabledValue)
		if err != nil {
			return
		}
	}
	alarmConfigValue, ok := rpcStruct["alarm_config"]
	if ok && alarmConfigValue != nil {
		record.AlarmConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "alarm_config"), alarmConfigValue)
		if err != nil {
			return
		}
	}
	recentAlertsValue, ok := rpcStruct["recent_alerts"]
	if ok && recentAlertsValue != nil {
		record.RecentAlerts, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "recent_alerts"), recentAlertsValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVMPPRef(context string, input interface{}) (VMPPRef, error) {
	var ref VMPPRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VMPPRef(value), nil
}

func deserializeVMPPRefSet(context string, input interface{}) (slice []VMPPRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VMPPRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVMPPRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVMPPRefToVMPPRecordMap(context string, input interface{}) (goMap map[VMPPRef]VMPPRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VMPPRef]VMPPRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVMPPRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVMPPRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVMRecord(context string, input interface{}) (record VMRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumVMOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumVMOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	powerStateValue, ok := rpcStruct["power_state"]
	if ok && powerStateValue != nil {
		record.PowerState, err = deserializeEnumVMPowerState(fmt.Sprintf("%s.%s", context, "power_state"), powerStateValue)
		if err != nil {
			return
		}
	}
	userVersionValue, ok := rpcStruct["user_version"]
	if ok && userVersionValue != nil {
		record.UserVersion, err = deserializeInt(fmt.Sprintf("%s.%s", context, "user_version"), userVersionValue)
		if err != nil {
			return
		}
	}
	isATemplateValue, ok := rpcStruct["is_a_template"]
	if ok && isATemplateValue != nil {
		record.IsATemplate, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_a_template"), isATemplateValue)
		if err != nil {
			return
		}
	}
	isDefaultTemplateValue, ok := rpcStruct["is_default_template"]
	if ok && isDefaultTemplateValue != nil {
		record.IsDefaultTemplate, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_default_template"), isDefaultTemplateValue)
		if err != nil {
			return
		}
	}
	suspendVDIValue, ok := rpcStruct["suspend_VDI"]
	if ok && suspendVDIValue != nil {
		record.SuspendVDI, err = deserializeVDIRef(fmt.Sprintf("%s.%s", context, "suspend_VDI"), suspendVDIValue)
		if err != nil {
			return
		}
	}
	residentOnValue, ok := rpcStruct["resident_on"]
	if ok && residentOnValue != nil {
		record.ResidentOn, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "resident_on"), residentOnValue)
		if err != nil {
			return
		}
	}
	scheduledToBeResidentOnValue, ok := rpcStruct["scheduled_to_be_resident_on"]
	if ok && scheduledToBeResidentOnValue != nil {
		record.ScheduledToBeResidentOn, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "scheduled_to_be_resident_on"), scheduledToBeResidentOnValue)
		if err != nil {
			return
		}
	}
	affinityValue, ok := rpcStruct["affinity"]
	if ok && affinityValue != nil {
		record.Affinity, err = deserializeHostRef(fmt.Sprintf("%s.%s", context, "affinity"), affinityValue)
		if err != nil {
			return
		}
	}
	memoryOverheadValue, ok := rpcStruct["memory_overhead"]
	if ok && memoryOverheadValue != nil {
		record.MemoryOverhead, err = deserializeInt(fmt.Sprintf("%s.%s", context, "memory_overhead"), memoryOverheadValue)
		if err != nil {
			return
		}
	}
	memoryTargetValue, ok := rpcStruct["memory_target"]
	if ok && memoryTargetValue != nil {
		record.MemoryTarget, err = deserializeInt(fmt.Sprintf("%s.%s", context, "memory_target"), memoryTargetValue)
		if err != nil {
			return
		}
	}
	memoryStaticMaxValue, ok := rpcStruct["memory_static_max"]
	if ok && memoryStaticMaxValue != nil {
		record.MemoryStaticMax, err = deserializeInt(fmt.Sprintf("%s.%s", context, "memory_static_max"), memoryStaticMaxValue)
		if err != nil {
			return
		}
	}
	memoryDynamicMaxValue, ok := rpcStruct["memory_dynamic_max"]
	if ok && memoryDynamicMaxValue != nil {
		record.MemoryDynamicMax, err = deserializeInt(fmt.Sprintf("%s.%s", context, "memory_dynamic_max"), memoryDynamicMaxValue)
		if err != nil {
			return
		}
	}
	memoryDynamicMinValue, ok := rpcStruct["memory_dynamic_min"]
	if ok && memoryDynamicMinValue != nil {
		record.MemoryDynamicMin, err = deserializeInt(fmt.Sprintf("%s.%s", context, "memory_dynamic_min"), memoryDynamicMinValue)
		if err != nil {
			return
		}
	}
	memoryStaticMinValue, ok := rpcStruct["memory_static_min"]
	if ok && memoryStaticMinValue != nil {
		record.MemoryStaticMin, err = deserializeInt(fmt.Sprintf("%s.%s", context, "memory_static_min"), memoryStaticMinValue)
		if err != nil {
			return
		}
	}
	vCPUsParamsValue, ok := rpcStruct["VCPUs_params"]
	if ok && vCPUsParamsValue != nil {
		record.VCPUsParams, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "VCPUs_params"), vCPUsParamsValue)
		if err != nil {
			return
		}
	}
	vCPUsMaxValue, ok := rpcStruct["VCPUs_max"]
	if ok && vCPUsMaxValue != nil {
		record.VCPUsMax, err = deserializeInt(fmt.Sprintf("%s.%s", context, "VCPUs_max"), vCPUsMaxValue)
		if err != nil {
			return
		}
	}
	vCPUsAtStartupValue, ok := rpcStruct["VCPUs_at_startup"]
	if ok && vCPUsAtStartupValue != nil {
		record.VCPUsAtStartup, err = deserializeInt(fmt.Sprintf("%s.%s", context, "VCPUs_at_startup"), vCPUsAtStartupValue)
		if err != nil {
			return
		}
	}
	actionsAfterSoftrebootValue, ok := rpcStruct["actions_after_softreboot"]
	if ok && actionsAfterSoftrebootValue != nil {
		record.ActionsAfterSoftreboot, err = deserializeEnumOnSoftrebootBehavior(fmt.Sprintf("%s.%s", context, "actions_after_softreboot"), actionsAfterSoftrebootValue)
		if err != nil {
			return
		}
	}
	actionsAfterShutdownValue, ok := rpcStruct["actions_after_shutdown"]
	if ok && actionsAfterShutdownValue != nil {
		record.ActionsAfterShutdown, err = deserializeEnumOnNormalExit(fmt.Sprintf("%s.%s", context, "actions_after_shutdown"), actionsAfterShutdownValue)
		if err != nil {
			return
		}
	}
	actionsAfterRebootValue, ok := rpcStruct["actions_after_reboot"]
	if ok && actionsAfterRebootValue != nil {
		record.ActionsAfterReboot, err = deserializeEnumOnNormalExit(fmt.Sprintf("%s.%s", context, "actions_after_reboot"), actionsAfterRebootValue)
		if err != nil {
			return
		}
	}
	actionsAfterCrashValue, ok := rpcStruct["actions_after_crash"]
	if ok && actionsAfterCrashValue != nil {
		record.ActionsAfterCrash, err = deserializeEnumOnCrashBehaviour(fmt.Sprintf("%s.%s", context, "actions_after_crash"), actionsAfterCrashValue)
		if err != nil {
			return
		}
	}
	consolesValue, ok := rpcStruct["consoles"]
	if ok && consolesValue != nil {
		record.Consoles, err = deserializeConsoleRefSet(fmt.Sprintf("%s.%s", context, "consoles"), consolesValue)
		if err != nil {
			return
		}
	}
	vIFsValue, ok := rpcStruct["VIFs"]
	if ok && vIFsValue != nil {
		record.VIFs, err = deserializeVIFRefSet(fmt.Sprintf("%s.%s", context, "VIFs"), vIFsValue)
		if err != nil {
			return
		}
	}
	vBDsValue, ok := rpcStruct["VBDs"]
	if ok && vBDsValue != nil {
		record.VBDs, err = deserializeVBDRefSet(fmt.Sprintf("%s.%s", context, "VBDs"), vBDsValue)
		if err != nil {
			return
		}
	}
	vUSBsValue, ok := rpcStruct["VUSBs"]
	if ok && vUSBsValue != nil {
		record.VUSBs, err = deserializeVUSBRefSet(fmt.Sprintf("%s.%s", context, "VUSBs"), vUSBsValue)
		if err != nil {
			return
		}
	}
	crashDumpsValue, ok := rpcStruct["crash_dumps"]
	if ok && crashDumpsValue != nil {
		record.CrashDumps, err = deserializeCrashdumpRefSet(fmt.Sprintf("%s.%s", context, "crash_dumps"), crashDumpsValue)
		if err != nil {
			return
		}
	}
	vTPMsValue, ok := rpcStruct["VTPMs"]
	if ok && vTPMsValue != nil {
		record.VTPMs, err = deserializeVTPMRefSet(fmt.Sprintf("%s.%s", context, "VTPMs"), vTPMsValue)
		if err != nil {
			return
		}
	}
	pVBootloaderValue, ok := rpcStruct["PV_bootloader"]
	if ok && pVBootloaderValue != nil {
		record.PVBootloader, err = deserializeString(fmt.Sprintf("%s.%s", context, "PV_bootloader"), pVBootloaderValue)
		if err != nil {
			return
		}
	}
	pVKernelValue, ok := rpcStruct["PV_kernel"]
	if ok && pVKernelValue != nil {
		record.PVKernel, err = deserializeString(fmt.Sprintf("%s.%s", context, "PV_kernel"), pVKernelValue)
		if err != nil {
			return
		}
	}
	pVRamdiskValue, ok := rpcStruct["PV_ramdisk"]
	if ok && pVRamdiskValue != nil {
		record.PVRamdisk, err = deserializeString(fmt.Sprintf("%s.%s", context, "PV_ramdisk"), pVRamdiskValue)
		if err != nil {
			return
		}
	}
	pVArgsValue, ok := rpcStruct["PV_args"]
	if ok && pVArgsValue != nil {
		record.PVArgs, err = deserializeString(fmt.Sprintf("%s.%s", context, "PV_args"), pVArgsValue)
		if err != nil {
			return
		}
	}
	pVBootloaderArgsValue, ok := rpcStruct["PV_bootloader_args"]
	if ok && pVBootloaderArgsValue != nil {
		record.PVBootloaderArgs, err = deserializeString(fmt.Sprintf("%s.%s", context, "PV_bootloader_args"), pVBootloaderArgsValue)
		if err != nil {
			return
		}
	}
	pVLegacyArgsValue, ok := rpcStruct["PV_legacy_args"]
	if ok && pVLegacyArgsValue != nil {
		record.PVLegacyArgs, err = deserializeString(fmt.Sprintf("%s.%s", context, "PV_legacy_args"), pVLegacyArgsValue)
		if err != nil {
			return
		}
	}
	hVMBootPolicyValue, ok := rpcStruct["HVM_boot_policy"]
	if ok && hVMBootPolicyValue != nil {
		record.HVMBootPolicy, err = deserializeString(fmt.Sprintf("%s.%s", context, "HVM_boot_policy"), hVMBootPolicyValue)
		if err != nil {
			return
		}
	}
	hVMBootParamsValue, ok := rpcStruct["HVM_boot_params"]
	if ok && hVMBootParamsValue != nil {
		record.HVMBootParams, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "HVM_boot_params"), hVMBootParamsValue)
		if err != nil {
			return
		}
	}
	hVMShadowMultiplierValue, ok := rpcStruct["HVM_shadow_multiplier"]
	if ok && hVMShadowMultiplierValue != nil {
		record.HVMShadowMultiplier, err = deserializeFloat(fmt.Sprintf("%s.%s", context, "HVM_shadow_multiplier"), hVMShadowMultiplierValue)
		if err != nil {
			return
		}
	}
	platformValue, ok := rpcStruct["platform"]
	if ok && platformValue != nil {
		record.Platform, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "platform"), platformValue)
		if err != nil {
			return
		}
	}
	pCIBusValue, ok := rpcStruct["PCI_bus"]
	if ok && pCIBusValue != nil {
		record.PCIBus, err = deserializeString(fmt.Sprintf("%s.%s", context, "PCI_bus"), pCIBusValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	domidValue, ok := rpcStruct["domid"]
	if ok && domidValue != nil {
		record.Domid, err = deserializeInt(fmt.Sprintf("%s.%s", context, "domid"), domidValue)
		if err != nil {
			return
		}
	}
	domarchValue, ok := rpcStruct["domarch"]
	if ok && domarchValue != nil {
		record.Domarch, err = deserializeString(fmt.Sprintf("%s.%s", context, "domarch"), domarchValue)
		if err != nil {
			return
		}
	}
	lastBootCPUFlagsValue, ok := rpcStruct["last_boot_CPU_flags"]
	if ok && lastBootCPUFlagsValue != nil {
		record.LastBootCPUFlags, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "last_boot_CPU_flags"), lastBootCPUFlagsValue)
		if err != nil {
			return
		}
	}
	isControlDomainValue, ok := rpcStruct["is_control_domain"]
	if ok && isControlDomainValue != nil {
		record.IsControlDomain, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_control_domain"), isControlDomainValue)
		if err != nil {
			return
		}
	}
	metricsValue, ok := rpcStruct["metrics"]
	if ok && metricsValue != nil {
		record.Metrics, err = deserializeVMMetricsRef(fmt.Sprintf("%s.%s", context, "metrics"), metricsValue)
		if err != nil {
			return
		}
	}
	guestMetricsValue, ok := rpcStruct["guest_metrics"]
	if ok && guestMetricsValue != nil {
		record.GuestMetrics, err = deserializeVMGuestMetricsRef(fmt.Sprintf("%s.%s", context, "guest_metrics"), guestMetricsValue)
		if err != nil {
			return
		}
	}
	lastBootedRecordValue, ok := rpcStruct["last_booted_record"]
	if ok && lastBootedRecordValue != nil {
		record.LastBootedRecord, err = deserializeString(fmt.Sprintf("%s.%s", context, "last_booted_record"), lastBootedRecordValue)
		if err != nil {
			return
		}
	}
	recommendationsValue, ok := rpcStruct["recommendations"]
	if ok && recommendationsValue != nil {
		record.Recommendations, err = deserializeString(fmt.Sprintf("%s.%s", context, "recommendations"), recommendationsValue)
		if err != nil {
			return
		}
	}
	xenstoreDataValue, ok := rpcStruct["xenstore_data"]
	if ok && xenstoreDataValue != nil {
		record.XenstoreData, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "xenstore_data"), xenstoreDataValue)
		if err != nil {
			return
		}
	}
	haAlwaysRunValue, ok := rpcStruct["ha_always_run"]
	if ok && haAlwaysRunValue != nil {
		record.HaAlwaysRun, err = deserializeBool(fmt.Sprintf("%s.%s", context, "ha_always_run"), haAlwaysRunValue)
		if err != nil {
			return
		}
	}
	haRestartPriorityValue, ok := rpcStruct["ha_restart_priority"]
	if ok && haRestartPriorityValue != nil {
		record.HaRestartPriority, err = deserializeString(fmt.Sprintf("%s.%s", context, "ha_restart_priority"), haRestartPriorityValue)
		if err != nil {
			return
		}
	}
	isASnapshotValue, ok := rpcStruct["is_a_snapshot"]
	if ok && isASnapshotValue != nil {
		record.IsASnapshot, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_a_snapshot"), isASnapshotValue)
		if err != nil {
			return
		}
	}
	snapshotOfValue, ok := rpcStruct["snapshot_of"]
	if ok && snapshotOfValue != nil {
		record.SnapshotOf, err = deserializeVMRef(fmt.Sprintf("%s.%s", context, "snapshot_of"), snapshotOfValue)
		if err != nil {
			return
		}
	}
	snapshotsValue, ok := rpcStruct["snapshots"]
	if ok && snapshotsValue != nil {
		record.Snapshots, err = deserializeVMRefSet(fmt.Sprintf("%s.%s", context, "snapshots"), snapshotsValue)
		if err != nil {
			return
		}
	}
	snapshotTimeValue, ok := rpcStruct["snapshot_time"]
	if ok && snapshotTimeValue != nil {
		record.SnapshotTime, err = deserializeTime(fmt.Sprintf("%s.%s", context, "snapshot_time"), snapshotTimeValue)
		if err != nil {
			return
		}
	}
	transportableSnapshotIDValue, ok := rpcStruct["transportable_snapshot_id"]
	if ok && transportableSnapshotIDValue != nil {
		record.TransportableSnapshotID, err = deserializeString(fmt.Sprintf("%s.%s", context, "transportable_snapshot_id"), transportableSnapshotIDValue)
		if err != nil {
			return
		}
	}
	blobsValue, ok := rpcStruct["blobs"]
	if ok && blobsValue != nil {
		record.Blobs, err = deserializeStringToBlobRefMap(fmt.Sprintf("%s.%s", context, "blobs"), blobsValue)
		if err != nil {
			return
		}
	}
	tagsValue, ok := rpcStruct["tags"]
	if ok && tagsValue != nil {
		record.Tags, err = deserializeStringSet(fmt.Sprintf("%s.%s", context, "tags"), tagsValue)
		if err != nil {
			return
		}
	}
	blockedOperationsValue, ok := rpcStruct["blocked_operations"]
	if ok && blockedOperationsValue != nil {
		record.BlockedOperations, err = deserializeEnumVMOperationsToStringMap(fmt.Sprintf("%s.%s", context, "blocked_operations"), blockedOperationsValue)
		if err != nil {
			return
		}
	}
	snapshotInfoValue, ok := rpcStruct["snapshot_info"]
	if ok && snapshotInfoValue != nil {
		record.SnapshotInfo, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "snapshot_info"), snapshotInfoValue)
		if err != nil {
			return
		}
	}
	snapshotMetadataValue, ok := rpcStruct["snapshot_metadata"]
	if ok && snapshotMetadataValue != nil {
		record.SnapshotMetadata, err = deserializeString(fmt.Sprintf("%s.%s", context, "snapshot_metadata"), snapshotMetadataValue)
		if err != nil {
			return
		}
	}
	parentValue, ok := rpcStruct["parent"]
	if ok && parentValue != nil {
		record.Parent, err = deserializeVMRef(fmt.Sprintf("%s.%s", context, "parent"), parentValue)
		if err != nil {
			return
		}
	}
	childrenValue, ok := rpcStruct["children"]
	if ok && childrenValue != nil {
		record.Children, err = deserializeVMRefSet(fmt.Sprintf("%s.%s", context, "children"), childrenValue)
		if err != nil {
			return
		}
	}
	biosStringsValue, ok := rpcStruct["bios_strings"]
	if ok && biosStringsValue != nil {
		record.BiosStrings, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "bios_strings"), biosStringsValue)
		if err != nil {
			return
		}
	}
	protectionPolicyValue, ok := rpcStruct["protection_policy"]
	if ok && protectionPolicyValue != nil {
		record.ProtectionPolicy, err = deserializeVMPPRef(fmt.Sprintf("%s.%s", context, "protection_policy"), protectionPolicyValue)
		if err != nil {
			return
		}
	}
	isSnapshotFromVmppValue, ok := rpcStruct["is_snapshot_from_vmpp"]
	if ok && isSnapshotFromVmppValue != nil {
		record.IsSnapshotFromVmpp, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_snapshot_from_vmpp"), isSnapshotFromVmppValue)
		if err != nil {
			return
		}
	}
	snapshotScheduleValue, ok := rpcStruct["snapshot_schedule"]
	if ok && snapshotScheduleValue != nil {
		record.SnapshotSchedule, err = deserializeVMSSRef(fmt.Sprintf("%s.%s", context, "snapshot_schedule"), snapshotScheduleValue)
		if err != nil {
			return
		}
	}
	isVmssSnapshotValue, ok := rpcStruct["is_vmss_snapshot"]
	if ok && isVmssSnapshotValue != nil {
		record.IsVmssSnapshot, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_vmss_snapshot"), isVmssSnapshotValue)
		if err != nil {
			return
		}
	}
	applianceValue, ok := rpcStruct["appliance"]
	if ok && applianceValue != nil {
		record.Appliance, err = deserializeVMApplianceRef(fmt.Sprintf("%s.%s", context, "appliance"), applianceValue)
		if err != nil {
			return
		}
	}
	startDelayValue, ok := rpcStruct["start_delay"]
	if ok && startDelayValue != nil {
		record.StartDelay, err = deserializeInt(fmt.Sprintf("%s.%s", context, "start_delay"), startDelayValue)
		if err != nil {
			return
		}
	}
	shutdownDelayValue, ok := rpcStruct["shutdown_delay"]
	if ok && shutdownDelayValue != nil {
		record.ShutdownDelay, err = deserializeInt(fmt.Sprintf("%s.%s", context, "shutdown_delay"), shutdownDelayValue)
		if err != nil {
			return
		}
	}
	orderValue, ok := rpcStruct["order"]
	if ok && orderValue != nil {
		record.Order, err = deserializeInt(fmt.Sprintf("%s.%s", context, "order"), orderValue)
		if err != nil {
			return
		}
	}
	vGPUsValue, ok := rpcStruct["VGPUs"]
	if ok && vGPUsValue != nil {
		record.VGPUs, err = deserializeVGPURefSet(fmt.Sprintf("%s.%s", context, "VGPUs"), vGPUsValue)
		if err != nil {
			return
		}
	}
	attachedPCIsValue, ok := rpcStruct["attached_PCIs"]
	if ok && attachedPCIsValue != nil {
		record.AttachedPCIs, err = deserializePCIRefSet(fmt.Sprintf("%s.%s", context, "attached_PCIs"), attachedPCIsValue)
		if err != nil {
			return
		}
	}
	suspendSRValue, ok := rpcStruct["suspend_SR"]
	if ok && suspendSRValue != nil {
		record.SuspendSR, err = deserializeSRRef(fmt.Sprintf("%s.%s", context, "suspend_SR"), suspendSRValue)
		if err != nil {
			return
		}
	}
	versionValue, ok := rpcStruct["version"]
	if ok && versionValue != nil {
		record.Version, err = deserializeInt(fmt.Sprintf("%s.%s", context, "version"), versionValue)
		if err != nil {
			return
		}
	}
	generationIDValue, ok := rpcStruct["generation_id"]
	if ok && generationIDValue != nil {
		record.GenerationID, err = deserializeString(fmt.Sprintf("%s.%s", context, "generation_id"), generationIDValue)
		if err != nil {
			return
		}
	}
	hardwarePlatformVersionValue, ok := rpcStruct["hardware_platform_version"]
	if ok && hardwarePlatformVersionValue != nil {
		record.HardwarePlatformVersion, err = deserializeInt(fmt.Sprintf("%s.%s", context, "hardware_platform_version"), hardwarePlatformVersionValue)
		if err != nil {
			return
		}
	}
	hasVendorDeviceValue, ok := rpcStruct["has_vendor_device"]
	if ok && hasVendorDeviceValue != nil {
		record.HasVendorDevice, err = deserializeBool(fmt.Sprintf("%s.%s", context, "has_vendor_device"), hasVendorDeviceValue)
		if err != nil {
			return
		}
	}
	requiresRebootValue, ok := rpcStruct["requires_reboot"]
	if ok && requiresRebootValue != nil {
		record.RequiresReboot, err = deserializeBool(fmt.Sprintf("%s.%s", context, "requires_reboot"), requiresRebootValue)
		if err != nil {
			return
		}
	}
	referenceLabelValue, ok := rpcStruct["reference_label"]
	if ok && referenceLabelValue != nil {
		record.ReferenceLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "reference_label"), referenceLabelValue)
		if err != nil {
			return
		}
	}
	domainTypeValue, ok := rpcStruct["domain_type"]
	if ok && domainTypeValue != nil {
		record.DomainType, err = deserializeEnumDomainType(fmt.Sprintf("%s.%s", context, "domain_type"), domainTypeValue)
		if err != nil {
			return
		}
	}
	nVRAMValue, ok := rpcStruct["NVRAM"]
	if ok && nVRAMValue != nil {
		record.NVRAM, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "NVRAM"), nVRAMValue)
		if err != nil {
			return
		}
	}
	pendingGuidancesValue, ok := rpcStruct["pending_guidances"]
	if ok && pendingGuidancesValue != nil {
		record.PendingGuidances, err = deserializeEnumUpdateGuidancesSet(fmt.Sprintf("%s.%s", context, "pending_guidances"), pendingGuidancesValue)
		if err != nil {
			return
		}
	}
	pendingGuidancesRecommendedValue, ok := rpcStruct["pending_guidances_recommended"]
	if ok && pendingGuidancesRecommendedValue != nil {
		record.PendingGuidancesRecommended, err = deserializeEnumUpdateGuidancesSet(fmt.Sprintf("%s.%s", context, "pending_guidances_recommended"), pendingGuidancesRecommendedValue)
		if err != nil {
			return
		}
	}
	pendingGuidancesFullValue, ok := rpcStruct["pending_guidances_full"]
	if ok && pendingGuidancesFullValue != nil {
		record.PendingGuidancesFull, err = deserializeEnumUpdateGuidancesSet(fmt.Sprintf("%s.%s", context, "pending_guidances_full"), pendingGuidancesFullValue)
		if err != nil {
			return
		}
	}
	groupsValue, ok := rpcStruct["groups"]
	if ok && groupsValue != nil {
		record.Groups, err = deserializeVMGroupRefSet(fmt.Sprintf("%s.%s", context, "groups"), groupsValue)
		if err != nil {
			return
		}
	}
	securebootCertificatesStateValue, ok := rpcStruct["secureboot_certificates_state"]
	if ok && securebootCertificatesStateValue != nil {
		record.SecurebootCertificatesState, err = deserializeEnumVMSecurebootCertificatesState(fmt.Sprintf("%s.%s", context, "secureboot_certificates_state"), securebootCertificatesStateValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVMRef(context string, input interface{}) (VMRef, error) {
	var ref VMRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VMRef(value), nil
}

func deserializeVMRefSet(context string, input interface{}) (slice []VMRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VMRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVMRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVMRefToStringSetMap(context string, input interface{}) (goMap map[VMRef][]string, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VMRef][]string, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVMRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeStringSet(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVMRefToStringToStringMapMap(context string, input interface{}) (goMap map[VMRef]map[string]string, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VMRef]map[string]string, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVMRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeStringToStringMap(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVMRefToVMRecordMap(context string, input interface{}) (goMap map[VMRef]VMRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VMRef]VMRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVMRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVMRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVMSSRecord(context string, input interface{}) (record VMSSRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	nameLabelValue, ok := rpcStruct["name_label"]
	if ok && nameLabelValue != nil {
		record.NameLabel, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_label"), nameLabelValue)
		if err != nil {
			return
		}
	}
	nameDescriptionValue, ok := rpcStruct["name_description"]
	if ok && nameDescriptionValue != nil {
		record.NameDescription, err = deserializeString(fmt.Sprintf("%s.%s", context, "name_description"), nameDescriptionValue)
		if err != nil {
			return
		}
	}
	enabledValue, ok := rpcStruct["enabled"]
	if ok && enabledValue != nil {
		record.Enabled, err = deserializeBool(fmt.Sprintf("%s.%s", context, "enabled"), enabledValue)
		if err != nil {
			return
		}
	}
	typeValue, ok := rpcStruct["type"]
	if ok && typeValue != nil {
		record.Type, err = deserializeEnumVmssType(fmt.Sprintf("%s.%s", context, "type"), typeValue)
		if err != nil {
			return
		}
	}
	retainedSnapshotsValue, ok := rpcStruct["retained_snapshots"]
	if ok && retainedSnapshotsValue != nil {
		record.RetainedSnapshots, err = deserializeInt(fmt.Sprintf("%s.%s", context, "retained_snapshots"), retainedSnapshotsValue)
		if err != nil {
			return
		}
	}
	frequencyValue, ok := rpcStruct["frequency"]
	if ok && frequencyValue != nil {
		record.Frequency, err = deserializeEnumVmssFrequency(fmt.Sprintf("%s.%s", context, "frequency"), frequencyValue)
		if err != nil {
			return
		}
	}
	scheduleValue, ok := rpcStruct["schedule"]
	if ok && scheduleValue != nil {
		record.Schedule, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "schedule"), scheduleValue)
		if err != nil {
			return
		}
	}
	lastRunTimeValue, ok := rpcStruct["last_run_time"]
	if ok && lastRunTimeValue != nil {
		record.LastRunTime, err = deserializeTime(fmt.Sprintf("%s.%s", context, "last_run_time"), lastRunTimeValue)
		if err != nil {
			return
		}
	}
	vMsValue, ok := rpcStruct["VMs"]
	if ok && vMsValue != nil {
		record.VMs, err = deserializeVMRefSet(fmt.Sprintf("%s.%s", context, "VMs"), vMsValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVMSSRef(context string, input interface{}) (VMSSRef, error) {
	var ref VMSSRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VMSSRef(value), nil
}

func deserializeVMSSRefSet(context string, input interface{}) (slice []VMSSRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VMSSRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVMSSRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVMSSRefToVMSSRecordMap(context string, input interface{}) (goMap map[VMSSRef]VMSSRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VMSSRef]VMSSRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVMSSRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVMSSRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVTPMRecord(context string, input interface{}) (record VTPMRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumVtpmOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumVtpmOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	vMValue, ok := rpcStruct["VM"]
	if ok && vMValue != nil {
		record.VM, err = deserializeVMRef(fmt.Sprintf("%s.%s", context, "VM"), vMValue)
		if err != nil {
			return
		}
	}
	backendValue, ok := rpcStruct["backend"]
	if ok && backendValue != nil {
		record.Backend, err = deserializeVMRef(fmt.Sprintf("%s.%s", context, "backend"), backendValue)
		if err != nil {
			return
		}
	}
	persistenceBackendValue, ok := rpcStruct["persistence_backend"]
	if ok && persistenceBackendValue != nil {
		record.PersistenceBackend, err = deserializeEnumPersistenceBackend(fmt.Sprintf("%s.%s", context, "persistence_backend"), persistenceBackendValue)
		if err != nil {
			return
		}
	}
	isUniqueValue, ok := rpcStruct["is_unique"]
	if ok && isUniqueValue != nil {
		record.IsUnique, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_unique"), isUniqueValue)
		if err != nil {
			return
		}
	}
	isProtectedValue, ok := rpcStruct["is_protected"]
	if ok && isProtectedValue != nil {
		record.IsProtected, err = deserializeBool(fmt.Sprintf("%s.%s", context, "is_protected"), isProtectedValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVTPMRef(context string, input interface{}) (VTPMRef, error) {
	var ref VTPMRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VTPMRef(value), nil
}

func deserializeVTPMRefSet(context string, input interface{}) (slice []VTPMRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VTPMRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVTPMRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVTPMRefToVTPMRecordMap(context string, input interface{}) (goMap map[VTPMRef]VTPMRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VTPMRef]VTPMRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVTPMRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVTPMRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVUSBRecord(context string, input interface{}) (record VUSBRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	uUIDValue, ok := rpcStruct["uuid"]
	if ok && uUIDValue != nil {
		record.UUID, err = deserializeString(fmt.Sprintf("%s.%s", context, "uuid"), uUIDValue)
		if err != nil {
			return
		}
	}
	allowedOperationsValue, ok := rpcStruct["allowed_operations"]
	if ok && allowedOperationsValue != nil {
		record.AllowedOperations, err = deserializeEnumVusbOperationsSet(fmt.Sprintf("%s.%s", context, "allowed_operations"), allowedOperationsValue)
		if err != nil {
			return
		}
	}
	currentOperationsValue, ok := rpcStruct["current_operations"]
	if ok && currentOperationsValue != nil {
		record.CurrentOperations, err = deserializeStringToEnumVusbOperationsMap(fmt.Sprintf("%s.%s", context, "current_operations"), currentOperationsValue)
		if err != nil {
			return
		}
	}
	vMValue, ok := rpcStruct["VM"]
	if ok && vMValue != nil {
		record.VM, err = deserializeVMRef(fmt.Sprintf("%s.%s", context, "VM"), vMValue)
		if err != nil {
			return
		}
	}
	uSBGroupValue, ok := rpcStruct["USB_group"]
	if ok && uSBGroupValue != nil {
		record.USBGroup, err = deserializeUSBGroupRef(fmt.Sprintf("%s.%s", context, "USB_group"), uSBGroupValue)
		if err != nil {
			return
		}
	}
	otherConfigValue, ok := rpcStruct["other_config"]
	if ok && otherConfigValue != nil {
		record.OtherConfig, err = deserializeStringToStringMap(fmt.Sprintf("%s.%s", context, "other_config"), otherConfigValue)
		if err != nil {
			return
		}
	}
	currentlyAttachedValue, ok := rpcStruct["currently_attached"]
	if ok && currentlyAttachedValue != nil {
		record.CurrentlyAttached, err = deserializeBool(fmt.Sprintf("%s.%s", context, "currently_attached"), currentlyAttachedValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVUSBRef(context string, input interface{}) (VUSBRef, error) {
	var ref VUSBRef
	value, ok := input.(string)
	if !ok {
		return ref, fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	}
	return VUSBRef(value), nil
}

func deserializeVUSBRefSet(context string, input interface{}) (slice []VUSBRef, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VUSBRef, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVUSBRef(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeVUSBRefToVUSBRecordMap(context string, input interface{}) (goMap map[VUSBRef]VUSBRecord, err error) {
	xenMap, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make(map[VUSBRef]VUSBRecord, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := deserializeVUSBRef(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := deserializeVUSBRecord(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}

func deserializeVdiNbdServerInfoRecord(context string, input interface{}) (record VdiNbdServerInfoRecord, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	exportnameValue, ok := rpcStruct["exportname"]
	if ok && exportnameValue != nil {
		record.Exportname, err = deserializeString(fmt.Sprintf("%s.%s", context, "exportname"), exportnameValue)
		if err != nil {
			return
		}
	}
	addressValue, ok := rpcStruct["address"]
	if ok && addressValue != nil {
		record.Address, err = deserializeString(fmt.Sprintf("%s.%s", context, "address"), addressValue)
		if err != nil {
			return
		}
	}
	portValue, ok := rpcStruct["port"]
	if ok && portValue != nil {
		record.Port, err = deserializeInt(fmt.Sprintf("%s.%s", context, "port"), portValue)
		if err != nil {
			return
		}
	}
	certValue, ok := rpcStruct["cert"]
	if ok && certValue != nil {
		record.Cert, err = deserializeString(fmt.Sprintf("%s.%s", context, "cert"), certValue)
		if err != nil {
			return
		}
	}
	subjectValue, ok := rpcStruct["subject"]
	if ok && subjectValue != nil {
		record.Subject, err = deserializeString(fmt.Sprintf("%s.%s", context, "subject"), subjectValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeVdiNbdServerInfoRecordSet(context string, input interface{}) (slice []VdiNbdServerInfoRecord, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make([]VdiNbdServerInfoRecord, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := deserializeVdiNbdServerInfoRecord(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}

func deserializeEventBatch(context string, input interface{}) (batch EventBatch, err error) {
	rpcStruct, ok := input.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "map[string]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	tokenValue, ok := rpcStruct["token"]
	if ok && tokenValue != nil {
		batch.Token, err = deserializeString(fmt.Sprintf("%s.%s", context, "token"), tokenValue)
		if err != nil {
			return
		}
	}
	validRefCountsValue, ok := rpcStruct["validRefCounts"]
	if ok && validRefCountsValue != nil {
		batch.ValidRefCounts, err = deserializeStringToIntMap(fmt.Sprintf("%s.%s", context, "validRefCounts"), validRefCountsValue)
		if err != nil {
			return
		}
	}
	eventsValue, ok := rpcStruct["events"]
	if ok && eventsValue != nil {
		batch.Events, err = deserializeEventRecordSet(fmt.Sprintf("%s.%s", context, "events"), eventsValue)
		if err != nil {
			return
		}
	}
	return
}

func deserializeRecordInterface(context string, input interface{}) (inter RecordInterface, err error) {
	_ = context
	inter = input
	return
}
