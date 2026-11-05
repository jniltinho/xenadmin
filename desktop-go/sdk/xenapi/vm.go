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

type VMRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// list of the operations allowed in this state. This list is advisory only and the server state may have changed by the time this field is read by a client.
	AllowedOperations []VMOperations `json:"allowed_operations,omitempty"`
	// links each of the running tasks using this object (by reference) to a current_operation enum which describes the nature of the task.
	CurrentOperations map[string]VMOperations `json:"current_operations,omitempty"`
	// a human-readable name
	NameLabel string `json:"name_label,omitempty"`
	// a notes field containing human-readable description
	NameDescription string `json:"name_description,omitempty"`
	// Current power state of the machine
	PowerState VMPowerState `json:"power_state,omitempty"`
	// Creators of VMs and templates may store version information here.
	UserVersion int `json:"user_version,omitempty"`
	// true if this is a template. Template VMs can never be started, they are used only for cloning other VMs
	IsATemplate bool `json:"is_a_template,omitempty"`
	// true if this is a default template. Default template VMs can never be started or migrated, they are used only for cloning other VMs
	IsDefaultTemplate bool `json:"is_default_template,omitempty"`
	// The VDI that a suspend image is stored on. (Only has meaning if VM is currently suspended)
	SuspendVDI VDIRef `json:"suspend_VDI,omitempty"`
	// the host the VM is currently resident on
	ResidentOn HostRef `json:"resident_on,omitempty"`
	// the host on which the VM is due to be started/resumed/migrated. This acts as a memory reservation indicator
	ScheduledToBeResidentOn HostRef `json:"scheduled_to_be_resident_on,omitempty"`
	// A host which the VM has some affinity for (or NULL). This is used as a hint to the start call when it decides where to run the VM. Resource constraints may cause the VM to be started elsewhere.
	Affinity HostRef `json:"affinity,omitempty"`
	// Virtualization memory overhead (bytes).
	MemoryOverhead int `json:"memory_overhead,omitempty"`
	// Dynamically-set memory target (bytes). The value of this field indicates the current target for memory available to this VM.
	MemoryTarget int `json:"memory_target,omitempty"`
	// Statically-set (i.e. absolute) maximum (bytes). The value of this field at VM start time acts as a hard limit of the amount of memory a guest can use. New values only take effect on reboot.
	MemoryStaticMax int `json:"memory_static_max,omitempty"`
	// Dynamic maximum (bytes)
	MemoryDynamicMax int `json:"memory_dynamic_max,omitempty"`
	// Dynamic minimum (bytes)
	MemoryDynamicMin int `json:"memory_dynamic_min,omitempty"`
	// Statically-set (i.e. absolute) mininum (bytes). The value of this field indicates the least amount of memory this VM can boot with without crashing.
	MemoryStaticMin int `json:"memory_static_min,omitempty"`
	// configuration parameters for the selected VCPU policy
	VCPUsParams map[string]string `json:"VCPUs_params,omitempty"`
	// Max number of VCPUs
	VCPUsMax int `json:"VCPUs_max,omitempty"`
	// Boot number of VCPUs
	VCPUsAtStartup int `json:"VCPUs_at_startup,omitempty"`
	// action to take after soft reboot
	ActionsAfterSoftreboot OnSoftrebootBehavior `json:"actions_after_softreboot,omitempty"`
	// action to take after the guest has shutdown itself
	ActionsAfterShutdown OnNormalExit `json:"actions_after_shutdown,omitempty"`
	// action to take after the guest has rebooted itself
	ActionsAfterReboot OnNormalExit `json:"actions_after_reboot,omitempty"`
	// action to take if the guest crashes
	ActionsAfterCrash OnCrashBehaviour `json:"actions_after_crash,omitempty"`
	// virtual console devices
	Consoles []ConsoleRef `json:"consoles,omitempty"`
	// virtual network interfaces
	VIFs []VIFRef `json:"VIFs,omitempty"`
	// virtual block devices
	VBDs []VBDRef `json:"VBDs,omitempty"`
	// virtual usb devices
	VUSBs []VUSBRef `json:"VUSBs,omitempty"`
	// crash dumps associated with this VM
	CrashDumps []CrashdumpRef `json:"crash_dumps,omitempty"`
	// virtual TPMs
	VTPMs []VTPMRef `json:"VTPMs,omitempty"`
	// name of or path to bootloader
	PVBootloader string `json:"PV_bootloader,omitempty"`
	// path to the kernel
	PVKernel string `json:"PV_kernel,omitempty"`
	// path to the initrd
	PVRamdisk string `json:"PV_ramdisk,omitempty"`
	// kernel command-line arguments
	PVArgs string `json:"PV_args,omitempty"`
	// miscellaneous arguments for the bootloader
	PVBootloaderArgs string `json:"PV_bootloader_args,omitempty"`
	// to make Zurich guests boot
	PVLegacyArgs string `json:"PV_legacy_args,omitempty"`
	// HVM boot policy
	HVMBootPolicy string `json:"HVM_boot_policy,omitempty"`
	// HVM boot params
	HVMBootParams map[string]string `json:"HVM_boot_params,omitempty"`
	// multiplier applied to the amount of shadow that will be made available to the guest
	HVMShadowMultiplier float64 `json:"HVM_shadow_multiplier,omitempty"`
	// platform-specific configuration
	Platform map[string]string `json:"platform,omitempty"`
	// PCI bus path for pass-through devices
	PCIBus string `json:"PCI_bus,omitempty"`
	// additional configuration
	OtherConfig map[string]string `json:"other_config,omitempty"`
	// domain ID (if available, -1 otherwise)
	Domid int `json:"domid,omitempty"`
	// Domain architecture (if available, null string otherwise)
	Domarch string `json:"domarch,omitempty"`
	// describes the CPU flags on which the VM was last booted
	LastBootCPUFlags map[string]string `json:"last_boot_CPU_flags,omitempty"`
	// true if this is a control domain (domain 0 or a driver domain)
	IsControlDomain bool `json:"is_control_domain,omitempty"`
	// metrics associated with this VM
	Metrics VMMetricsRef `json:"metrics,omitempty"`
	// metrics associated with the running guest
	GuestMetrics VMGuestMetricsRef `json:"guest_metrics,omitempty"`
	// marshalled value containing VM record at time of last boot
	LastBootedRecord string `json:"last_booted_record,omitempty"`
	// An XML specification of recommended values and ranges for properties of this VM
	Recommendations string `json:"recommendations,omitempty"`
	// data to be inserted into the xenstore tree (/local/domain/&lt;domid&gt;/vm-data) after the VM is created.
	XenstoreData map[string]string `json:"xenstore_data,omitempty"`
	// if true then the system will attempt to keep the VM running as much as possible.
	HaAlwaysRun bool `json:"ha_always_run,omitempty"`
	// has possible values: &quot;best-effort&quot; meaning &quot;try to restart this VM if possible but don&apos;t consider the Pool to be overcommitted if this is not possible&quot;; &quot;restart&quot; meaning &quot;this VM should be restarted&quot;; &quot;&quot; meaning &quot;do not try to restart this VM&quot;
	HaRestartPriority string `json:"ha_restart_priority,omitempty"`
	// true if this is a snapshot. Snapshotted VMs can never be started, they are used only for cloning other VMs
	IsASnapshot bool `json:"is_a_snapshot,omitempty"`
	// Ref pointing to the VM this snapshot is of.
	SnapshotOf VMRef `json:"snapshot_of,omitempty"`
	// List pointing to all the VM snapshots.
	Snapshots []VMRef `json:"snapshots,omitempty"`
	// Date/time when this snapshot was created.
	SnapshotTime time.Time `json:"snapshot_time,omitempty"`
	// Transportable ID of the snapshot VM
	TransportableSnapshotID string `json:"transportable_snapshot_id,omitempty"`
	// Binary blobs associated with this VM
	Blobs map[string]BlobRef `json:"blobs,omitempty"`
	// user-specified tags for categorization purposes
	Tags []string `json:"tags,omitempty"`
	// List of operations which have been explicitly blocked and an error code
	BlockedOperations map[VMOperations]string `json:"blocked_operations,omitempty"`
	// Human-readable information concerning this snapshot
	SnapshotInfo map[string]string `json:"snapshot_info,omitempty"`
	// Encoded information about the VM&apos;s metadata this is a snapshot of
	SnapshotMetadata string `json:"snapshot_metadata,omitempty"`
	// Ref pointing to the parent of this VM
	Parent VMRef `json:"parent,omitempty"`
	// List pointing to all the children of this VM
	Children []VMRef `json:"children,omitempty"`
	// BIOS strings
	BiosStrings map[string]string `json:"bios_strings,omitempty"`
	// Ref pointing to a protection policy for this VM
	ProtectionPolicy VMPPRef `json:"protection_policy,omitempty"`
	// true if this snapshot was created by the protection policy
	IsSnapshotFromVmpp bool `json:"is_snapshot_from_vmpp,omitempty"`
	// Ref pointing to a snapshot schedule for this VM
	SnapshotSchedule VMSSRef `json:"snapshot_schedule,omitempty"`
	// true if this snapshot was created by the snapshot schedule
	IsVmssSnapshot bool `json:"is_vmss_snapshot,omitempty"`
	// the appliance to which this VM belongs
	Appliance VMApplianceRef `json:"appliance,omitempty"`
	// The delay to wait before proceeding to the next order in the startup sequence (seconds)
	StartDelay int `json:"start_delay,omitempty"`
	// The delay to wait before proceeding to the next order in the shutdown sequence (seconds)
	ShutdownDelay int `json:"shutdown_delay,omitempty"`
	// The point in the startup or shutdown sequence at which this VM will be started
	Order int `json:"order,omitempty"`
	// Virtual GPUs
	VGPUs []VGPURef `json:"VGPUs,omitempty"`
	// Currently passed-through PCI devices
	AttachedPCIs []PCIRef `json:"attached_PCIs,omitempty"`
	// The SR on which a suspend image is stored
	SuspendSR SRRef `json:"suspend_SR,omitempty"`
	// The number of times this VM has been recovered
	Version int `json:"version,omitempty"`
	// Generation ID of the VM
	GenerationID string `json:"generation_id,omitempty"`
	// The host virtual hardware platform version the VM can run on
	HardwarePlatformVersion int `json:"hardware_platform_version,omitempty"`
	// When an HVM guest starts, this controls the presence of the emulated C000 PCI device which triggers Windows Update to fetch or update PV drivers.
	HasVendorDevice bool `json:"has_vendor_device,omitempty"`
	// Indicates whether a VM requires a reboot in order to update its configuration, for example, its memory allocation.
	RequiresReboot bool `json:"requires_reboot,omitempty"`
	// Textual reference to the template used to create a VM. This can be used by clients in need of an immutable reference to the template since the latter&apos;s uuid and name_label may change, for example, after a package installation or upgrade.
	ReferenceLabel string `json:"reference_label,omitempty"`
	// The type of domain that will be created when the VM is started
	DomainType DomainType `json:"domain_type,omitempty"`
	// initial value for guest NVRAM (containing UEFI variables, etc). Cannot be changed while the VM is running
	NVRAM map[string]string `json:"NVRAM,omitempty"`
	// The set of pending mandatory guidances after applying updates, which must be applied, otherwise there may be, for example, VM failures
	PendingGuidances []UpdateGuidances `json:"pending_guidances,omitempty"`
	// The set of pending recommended guidances after applying updates, which most users should follow to make the updates effective, but if not followed, will not cause a failure
	PendingGuidancesRecommended []UpdateGuidances `json:"pending_guidances_recommended,omitempty"`
	// The set of pending full guidances after applying updates, which a user should follow to make some updates, for example, specific hardware drivers or CPU features, fully effective, but the &apos;average user&apos; doesn&apos;t need to
	PendingGuidancesFull []UpdateGuidances `json:"pending_guidances_full,omitempty"`
	// VM groups associated with the VM
	Groups []VMGroupRef `json:"groups,omitempty"`
	// The state of the Secure Boot certificates, showing whether an update is available, already scheduled, or not needed.
	SecurebootCertificatesState VMSecurebootCertificatesState `json:"secureboot_certificates_state,omitempty"`
}

type VMRef string

// A virtual machine (or &apos;guest&apos;).
type vm struct{}

var VM vm

// GetAllRecords: Return a map of VM references to VM records for all VMs known to the system.
// Version: rio
func (vm) GetAllRecords(session *Session) (retval map[VMRef]VMRecord, err error) {
	method := "VM.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefToVMRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of VM references to VM records for all VMs known to the system.
// Version: rio
func (vm) GetAllRecords1(session *Session) (retval map[VMRef]VMRecord, err error) {
	method := "VM.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefToVMRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the VMs known to the system.
// Version: rio
func (vm) GetAll(session *Session) (retval []VMRef, err error) {
	method := "VM.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the VMs known to the system.
// Version: rio
func (vm) GetAll1(session *Session) (retval []VMRef, err error) {
	method := "VM.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// SetPlatform: Set the platform field of the given VM.
// Version: rio
func (vm) SetPlatform(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_platform"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPlatform3: Set the platform field of the given VM.
// Version: rio
func (vm) SetPlatform3(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_platform"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromPlatform: Remove the given key and its corresponding value from the platform field of the given VM. If the key is not in that Map, then do nothing.
// Version: rio
func (vm) RemoveFromPlatform(session *Session, self VMRef, key string) (err error) {
	method := "VM.remove_from_platform"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromPlatform3: Remove the given key and its corresponding value from the platform field of the given VM. If the key is not in that Map, then do nothing.
// Version: rio
func (vm) RemoveFromPlatform3(session *Session, self VMRef, key string) (err error) {
	method := "VM.remove_from_platform"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToPlatform: Add the given key-value pair to the platform field of the given VM.
// Version: rio
func (vm) AddToPlatform(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_platform"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToPlatform4: Add the given key-value pair to the platform field of the given VM.
// Version: rio
func (vm) AddToPlatform4(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_platform"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig: Set the other_config field of the given VM.
// Version: rio
func (vm) SetOtherConfig(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetOtherConfig3: Set the other_config field of the given VM.
// Version: rio
func (vm) SetOtherConfig3(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig: Remove the given key and its corresponding value from the other_config field of the given VM. If the key is not in that Map, then do nothing.
// Version: rio
func (vm) RemoveFromOtherConfig(session *Session, self VMRef, key string) (err error) {
	method := "VM.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromOtherConfig3: Remove the given key and its corresponding value from the other_config field of the given VM. If the key is not in that Map, then do nothing.
// Version: rio
func (vm) RemoveFromOtherConfig3(session *Session, self VMRef, key string) (err error) {
	method := "VM.remove_from_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig: Add the given key-value pair to the other_config field of the given VM.
// Version: rio
func (vm) AddToOtherConfig(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToOtherConfig4: Add the given key-value pair to the other_config field of the given VM.
// Version: rio
func (vm) AddToOtherConfig4(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Sysprep: Pass unattend.xml to Windows sysprep and wait for the VM to shut down as part of a reboot.
// Version: closed
func (vm) Sysprep(session *Session, self VMRef, unattend string, timeout float64) (err error) {
	method := "VM.sysprep"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	unattendArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "unattend"), unattend)
	if err != nil {
		return
	}
	timeoutArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "timeout"), timeout)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, unattendArg, timeoutArg)
	return
}

// AsyncSysprep: Pass unattend.xml to Windows sysprep and wait for the VM to shut down as part of a reboot.
// Version: closed
func (vm) AsyncSysprep(session *Session, self VMRef, unattend string, timeout float64) (retval TaskRef, err error) {
	method := "Async.VM.sysprep"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	unattendArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "unattend"), unattend)
	if err != nil {
		return
	}
	timeoutArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "timeout"), timeout)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, unattendArg, timeoutArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Sysprep4: Pass unattend.xml to Windows sysprep and wait for the VM to shut down as part of a reboot.
// Version: closed
func (vm) Sysprep4(session *Session, self VMRef, unattend string, timeout float64) (err error) {
	method := "VM.sysprep"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	unattendArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "unattend"), unattend)
	if err != nil {
		return
	}
	timeoutArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "timeout"), timeout)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, unattendArg, timeoutArg)
	return
}

// AsyncSysprep4: Pass unattend.xml to Windows sysprep and wait for the VM to shut down as part of a reboot.
// Version: closed
func (vm) AsyncSysprep4(session *Session, self VMRef, unattend string, timeout float64) (retval TaskRef, err error) {
	method := "Async.VM.sysprep"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	unattendArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "unattend"), unattend)
	if err != nil {
		return
	}
	timeoutArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "timeout"), timeout)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, unattendArg, timeoutArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RemoveFromBlockedOperations: Update list of operations which have been explicitly blocked and an error code
// Version: orlando
func (vm) RemoveFromBlockedOperations(session *Session, self VMRef, key VMOperations) (err error) {
	method := "VM.remove_from_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeEnumVMOperations(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// AsyncRemoveFromBlockedOperations: Update list of operations which have been explicitly blocked and an error code
// Version: orlando
func (vm) AsyncRemoveFromBlockedOperations(session *Session, self VMRef, key VMOperations) (retval TaskRef, err error) {
	method := "Async.VM.remove_from_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeEnumVMOperations(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RemoveFromBlockedOperations3: Update list of operations which have been explicitly blocked and an error code
// Version: orlando
func (vm) RemoveFromBlockedOperations3(session *Session, self VMRef, key VMOperations) (err error) {
	method := "VM.remove_from_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeEnumVMOperations(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	return
}

// AsyncRemoveFromBlockedOperations3: Update list of operations which have been explicitly blocked and an error code
// Version: orlando
func (vm) AsyncRemoveFromBlockedOperations3(session *Session, self VMRef, key VMOperations) (retval TaskRef, err error) {
	method := "Async.VM.remove_from_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeEnumVMOperations(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, keyArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AddToBlockedOperations: Update list of operations which have been explicitly blocked and an error code
// Version: orlando
func (vm) AddToBlockedOperations(session *Session, self VMRef, key VMOperations, value string) (err error) {
	method := "VM.add_to_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeEnumVMOperations(fmt.Sprintf("%s(%s)", method, "key"), key)
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

// AsyncAddToBlockedOperations: Update list of operations which have been explicitly blocked and an error code
// Version: orlando
func (vm) AsyncAddToBlockedOperations(session *Session, self VMRef, key VMOperations, value string) (retval TaskRef, err error) {
	method := "Async.VM.add_to_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeEnumVMOperations(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AddToBlockedOperations4: Update list of operations which have been explicitly blocked and an error code
// Version: orlando
func (vm) AddToBlockedOperations4(session *Session, self VMRef, key VMOperations, value string) (err error) {
	method := "VM.add_to_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeEnumVMOperations(fmt.Sprintf("%s(%s)", method, "key"), key)
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

// AsyncAddToBlockedOperations4: Update list of operations which have been explicitly blocked and an error code
// Version: orlando
func (vm) AsyncAddToBlockedOperations4(session *Session, self VMRef, key VMOperations, value string) (retval TaskRef, err error) {
	method := "Async.VM.add_to_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	keyArg, err := serializeEnumVMOperations(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	valueArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetBlockedOperations: Update list of operations which have been explicitly blocked and an error code
// Version: orlando
func (vm) SetBlockedOperations(session *Session, self VMRef, value map[VMOperations]string) (err error) {
	method := "VM.set_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumVMOperationsToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetBlockedOperations: Update list of operations which have been explicitly blocked and an error code
// Version: orlando
func (vm) AsyncSetBlockedOperations(session *Session, self VMRef, value map[VMOperations]string) (retval TaskRef, err error) {
	method := "Async.VM.set_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumVMOperationsToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetBlockedOperations3: Update list of operations which have been explicitly blocked and an error code
// Version: orlando
func (vm) SetBlockedOperations3(session *Session, self VMRef, value map[VMOperations]string) (err error) {
	method := "VM.set_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumVMOperationsToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetBlockedOperations3: Update list of operations which have been explicitly blocked and an error code
// Version: orlando
func (vm) AsyncSetBlockedOperations3(session *Session, self VMRef, value map[VMOperations]string) (retval TaskRef, err error) {
	method := "Async.VM.set_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumVMOperationsToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// UpdateSecurebootCertificatesOnBoot: Mark or unmark secure boot certificate update on VM boot
// Version: closed
func (vm) UpdateSecurebootCertificatesOnBoot(session *Session, self VMRef, mark bool) (err error) {
	method := "VM.update_secureboot_certificates_on_boot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	markArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "mark"), mark)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, markArg)
	return
}

// AsyncUpdateSecurebootCertificatesOnBoot: Mark or unmark secure boot certificate update on VM boot
// Version: closed
func (vm) AsyncUpdateSecurebootCertificatesOnBoot(session *Session, self VMRef, mark bool) (retval TaskRef, err error) {
	method := "Async.VM.update_secureboot_certificates_on_boot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	markArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "mark"), mark)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, markArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// UpdateSecurebootCertificatesOnBoot3: Mark or unmark secure boot certificate update on VM boot
// Version: closed
func (vm) UpdateSecurebootCertificatesOnBoot3(session *Session, self VMRef, mark bool) (err error) {
	method := "VM.update_secureboot_certificates_on_boot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	markArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "mark"), mark)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, markArg)
	return
}

// AsyncUpdateSecurebootCertificatesOnBoot3: Mark or unmark secure boot certificate update on VM boot
// Version: closed
func (vm) AsyncUpdateSecurebootCertificatesOnBoot3(session *Session, self VMRef, mark bool) (retval TaskRef, err error) {
	method := "Async.VM.update_secureboot_certificates_on_boot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	markArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "mark"), mark)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, markArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetSecurebootReadiness: Return the secureboot readiness of the VM
// Version: closed
func (vm) GetSecurebootReadiness(session *Session, self VMRef) (retval VMSecurebootReadiness, err error) {
	method := "VM.get_secureboot_readiness"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVMSecurebootReadiness(method+" -> ", result)
	return
}

// AsyncGetSecurebootReadiness: Return the secureboot readiness of the VM
// Version: closed
func (vm) AsyncGetSecurebootReadiness(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.get_secureboot_readiness"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSecurebootReadiness2: Return the secureboot readiness of the VM
// Version: closed
func (vm) GetSecurebootReadiness2(session *Session, self VMRef) (retval VMSecurebootReadiness, err error) {
	method := "VM.get_secureboot_readiness"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVMSecurebootReadiness(method+" -> ", result)
	return
}

// AsyncGetSecurebootReadiness2: Return the secureboot readiness of the VM
// Version: closed
func (vm) AsyncGetSecurebootReadiness2(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.get_secureboot_readiness"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetUefiMode: Set the UEFI mode of a VM
// Version: closed
func (vm) SetUefiMode(session *Session, self VMRef, mode VMUefiMode) (retval string, err error) {
	method := "VM.set_uefi_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumVMUefiMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, modeArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncSetUefiMode: Set the UEFI mode of a VM
// Version: closed
func (vm) AsyncSetUefiMode(session *Session, self VMRef, mode VMUefiMode) (retval TaskRef, err error) {
	method := "Async.VM.set_uefi_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumVMUefiMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, modeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetUefiMode3: Set the UEFI mode of a VM
// Version: closed
func (vm) SetUefiMode3(session *Session, self VMRef, mode VMUefiMode) (retval string, err error) {
	method := "VM.set_uefi_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumVMUefiMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, modeArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncSetUefiMode3: Set the UEFI mode of a VM
// Version: closed
func (vm) AsyncSetUefiMode3(session *Session, self VMRef, mode VMUefiMode) (retval TaskRef, err error) {
	method := "Async.VM.set_uefi_mode"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	modeArg, err := serializeEnumVMUefiMode(fmt.Sprintf("%s(%s)", method, "mode"), mode)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, modeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RestartDeviceModels: Restart device models of the VM
// Version: closed
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) RestartDeviceModels(session *Session, self VMRef) (err error) {
	method := "VM.restart_device_models"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncRestartDeviceModels: Restart device models of the VM
// Version: closed
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) AsyncRestartDeviceModels(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.restart_device_models"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RestartDeviceModels2: Restart device models of the VM
// Version: closed
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) RestartDeviceModels2(session *Session, self VMRef) (err error) {
	method := "VM.restart_device_models"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncRestartDeviceModels2: Restart device models of the VM
// Version: closed
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) AsyncRestartDeviceModels2(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.restart_device_models"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHVMBootPolicy: Set the VM.HVM_boot_policy field of the given VM, which will take effect when it is next started
// Version: rio
func (vm) SetHVMBootPolicy(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_HVM_boot_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHVMBootPolicy3: Set the VM.HVM_boot_policy field of the given VM, which will take effect when it is next started
// Version: rio
func (vm) SetHVMBootPolicy3(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_HVM_boot_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetDomainType: Set the VM.domain_type field of the given VM, which will take effect when it is next started
// Version: kolkata
func (vm) SetDomainType(session *Session, self VMRef, value DomainType) (err error) {
	method := "VM.set_domain_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumDomainType(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetDomainType3: Set the VM.domain_type field of the given VM, which will take effect when it is next started
// Version: kolkata
func (vm) SetDomainType3(session *Session, self VMRef, value DomainType) (err error) {
	method := "VM.set_domain_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumDomainType(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetActionsAfterCrash: Sets the actions_after_crash parameter
// Version: rio
func (vm) SetActionsAfterCrash(session *Session, self VMRef, value OnCrashBehaviour) (err error) {
	method := "VM.set_actions_after_crash"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnCrashBehaviour(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetActionsAfterCrash: Sets the actions_after_crash parameter
// Version: rio
func (vm) AsyncSetActionsAfterCrash(session *Session, self VMRef, value OnCrashBehaviour) (retval TaskRef, err error) {
	method := "Async.VM.set_actions_after_crash"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnCrashBehaviour(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetActionsAfterCrash3: Sets the actions_after_crash parameter
// Version: rio
func (vm) SetActionsAfterCrash3(session *Session, self VMRef, value OnCrashBehaviour) (err error) {
	method := "VM.set_actions_after_crash"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnCrashBehaviour(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetActionsAfterCrash3: Sets the actions_after_crash parameter
// Version: rio
func (vm) AsyncSetActionsAfterCrash3(session *Session, self VMRef, value OnCrashBehaviour) (retval TaskRef, err error) {
	method := "Async.VM.set_actions_after_crash"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnCrashBehaviour(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// Import: Import an XVA from a URI
// Version: dundee
func (vm) Import(session *Session, url string, sr SRRef, fullRestore bool, force bool) (retval []VMRef, err error) {
	method := "VM.import"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	fullRestoreArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "full_restore"), fullRestore)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, urlArg, srArg, fullRestoreArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// AsyncImport: Import an XVA from a URI
// Version: dundee
func (vm) AsyncImport(session *Session, url string, sr SRRef, fullRestore bool, force bool) (retval TaskRef, err error) {
	method := "Async.VM.import"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	fullRestoreArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "full_restore"), fullRestore)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, urlArg, srArg, fullRestoreArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Import5: Import an XVA from a URI
// Version: dundee
func (vm) Import5(session *Session, url string, sr SRRef, fullRestore bool, force bool) (retval []VMRef, err error) {
	method := "VM.import"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	fullRestoreArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "full_restore"), fullRestore)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, urlArg, srArg, fullRestoreArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// AsyncImport5: Import an XVA from a URI
// Version: dundee
func (vm) AsyncImport5(session *Session, url string, sr SRRef, fullRestore bool, force bool) (retval TaskRef, err error) {
	method := "Async.VM.import"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	urlArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "url"), url)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	fullRestoreArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "full_restore"), fullRestore)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, urlArg, srArg, fullRestoreArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetHasVendorDevice: Controls whether, when the VM starts in HVM mode, its virtual hardware will include the emulated PCI device for which drivers may be available through Windows Update. Usually this should never be changed on a VM on which Windows has been installed: changing it on such a VM is likely to lead to a crash on next start.
// Version: dundee
func (vm) SetHasVendorDevice(session *Session, self VMRef, value bool) (err error) {
	method := "VM.set_has_vendor_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetHasVendorDevice: Controls whether, when the VM starts in HVM mode, its virtual hardware will include the emulated PCI device for which drivers may be available through Windows Update. Usually this should never be changed on a VM on which Windows has been installed: changing it on such a VM is likely to lead to a crash on next start.
// Version: dundee
func (vm) AsyncSetHasVendorDevice(session *Session, self VMRef, value bool) (retval TaskRef, err error) {
	method := "Async.VM.set_has_vendor_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHasVendorDevice3: Controls whether, when the VM starts in HVM mode, its virtual hardware will include the emulated PCI device for which drivers may be available through Windows Update. Usually this should never be changed on a VM on which Windows has been installed: changing it on such a VM is likely to lead to a crash on next start.
// Version: dundee
func (vm) SetHasVendorDevice3(session *Session, self VMRef, value bool) (err error) {
	method := "VM.set_has_vendor_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetHasVendorDevice3: Controls whether, when the VM starts in HVM mode, its virtual hardware will include the emulated PCI device for which drivers may be available through Windows Update. Usually this should never be changed on a VM on which Windows has been installed: changing it on such a VM is likely to lead to a crash on next start.
// Version: dundee
func (vm) AsyncSetHasVendorDevice3(session *Session, self VMRef, value bool) (retval TaskRef, err error) {
	method := "Async.VM.set_has_vendor_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// CallHostPlugin: Call an API plugin on the host where this vm resides
// Version: closed
func (vm) CallHostPlugin(session *Session, vm VMRef, plugin string, fn string, args map[string]string) (retval string, err error) {
	method := "VM.call_host_plugin"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	pluginArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "plugin"), plugin)
	if err != nil {
		return
	}
	fnArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "fn"), fn)
	if err != nil {
		return
	}
	argsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, pluginArg, fnArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncCallHostPlugin: Call an API plugin on the host where this vm resides
// Version: closed
func (vm) AsyncCallHostPlugin(session *Session, vm VMRef, plugin string, fn string, args map[string]string) (retval TaskRef, err error) {
	method := "Async.VM.call_host_plugin"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	pluginArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "plugin"), plugin)
	if err != nil {
		return
	}
	fnArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "fn"), fn)
	if err != nil {
		return
	}
	argsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, pluginArg, fnArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CallHostPlugin5: Call an API plugin on the host where this vm resides
// Version: closed
func (vm) CallHostPlugin5(session *Session, vm VMRef, plugin string, fn string, args map[string]string) (retval string, err error) {
	method := "VM.call_host_plugin"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	pluginArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "plugin"), plugin)
	if err != nil {
		return
	}
	fnArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "fn"), fn)
	if err != nil {
		return
	}
	argsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, pluginArg, fnArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncCallHostPlugin5: Call an API plugin on the host where this vm resides
// Version: closed
func (vm) AsyncCallHostPlugin5(session *Session, vm VMRef, plugin string, fn string, args map[string]string) (retval TaskRef, err error) {
	method := "Async.VM.call_host_plugin"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	pluginArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "plugin"), plugin)
	if err != nil {
		return
	}
	fnArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "fn"), fn)
	if err != nil {
		return
	}
	argsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, pluginArg, fnArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CallPlugin: Call an API plugin on this vm
// Version: cream
func (vm) CallPlugin(session *Session, vm VMRef, plugin string, fn string, args map[string]string) (retval string, err error) {
	method := "VM.call_plugin"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	pluginArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "plugin"), plugin)
	if err != nil {
		return
	}
	fnArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "fn"), fn)
	if err != nil {
		return
	}
	argsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, pluginArg, fnArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncCallPlugin: Call an API plugin on this vm
// Version: cream
func (vm) AsyncCallPlugin(session *Session, vm VMRef, plugin string, fn string, args map[string]string) (retval TaskRef, err error) {
	method := "Async.VM.call_plugin"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	pluginArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "plugin"), plugin)
	if err != nil {
		return
	}
	fnArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "fn"), fn)
	if err != nil {
		return
	}
	argsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, pluginArg, fnArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CallPlugin5: Call an API plugin on this vm
// Version: cream
func (vm) CallPlugin5(session *Session, vm VMRef, plugin string, fn string, args map[string]string) (retval string, err error) {
	method := "VM.call_plugin"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	pluginArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "plugin"), plugin)
	if err != nil {
		return
	}
	fnArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "fn"), fn)
	if err != nil {
		return
	}
	argsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, pluginArg, fnArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// AsyncCallPlugin5: Call an API plugin on this vm
// Version: cream
func (vm) AsyncCallPlugin5(session *Session, vm VMRef, plugin string, fn string, args map[string]string) (retval TaskRef, err error) {
	method := "Async.VM.call_plugin"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	pluginArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "plugin"), plugin)
	if err != nil {
		return
	}
	fnArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "fn"), fn)
	if err != nil {
		return
	}
	argsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, pluginArg, fnArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// QueryServices: Query the system services advertised by this VM and register them. This can only be applied to a system domain.
// Version: tampa
func (vm) QueryServices(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.query_services"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncQueryServices: Query the system services advertised by this VM and register them. This can only be applied to a system domain.
// Version: tampa
func (vm) AsyncQueryServices(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.query_services"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// QueryServices2: Query the system services advertised by this VM and register them. This can only be applied to a system domain.
// Version: tampa
func (vm) QueryServices2(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.query_services"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncQueryServices2: Query the system services advertised by this VM and register them. This can only be applied to a system domain.
// Version: tampa
func (vm) AsyncQueryServices2(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.query_services"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetGroups: Associate this VM with VM groups.
// Version: closed
//
// Errors:
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
func (vm) SetGroups(session *Session, self VMRef, value []VMGroupRef) (err error) {
	method := "VM.set_groups"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVMGroupRefSet(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetGroups: Associate this VM with VM groups.
// Version: closed
//
// Errors:
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
func (vm) AsyncSetGroups(session *Session, self VMRef, value []VMGroupRef) (retval TaskRef, err error) {
	method := "Async.VM.set_groups"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVMGroupRefSet(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetGroups3: Associate this VM with VM groups.
// Version: closed
//
// Errors:
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
func (vm) SetGroups3(session *Session, self VMRef, value []VMGroupRef) (err error) {
	method := "VM.set_groups"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVMGroupRefSet(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetGroups3: Associate this VM with VM groups.
// Version: closed
//
// Errors:
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
func (vm) AsyncSetGroups3(session *Session, self VMRef, value []VMGroupRef) (retval TaskRef, err error) {
	method := "Async.VM.set_groups"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVMGroupRefSet(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetAppliance: Assign this VM to an appliance.
// Version: boston
func (vm) SetAppliance(session *Session, self VMRef, value VMApplianceRef) (err error) {
	method := "VM.set_appliance"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetAppliance: Assign this VM to an appliance.
// Version: boston
func (vm) AsyncSetAppliance(session *Session, self VMRef, value VMApplianceRef) (retval TaskRef, err error) {
	method := "Async.VM.set_appliance"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetAppliance3: Assign this VM to an appliance.
// Version: boston
func (vm) SetAppliance3(session *Session, self VMRef, value VMApplianceRef) (err error) {
	method := "VM.set_appliance"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetAppliance3: Assign this VM to an appliance.
// Version: boston
func (vm) AsyncSetAppliance3(session *Session, self VMRef, value VMApplianceRef) (retval TaskRef, err error) {
	method := "Async.VM.set_appliance"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVMApplianceRef(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// ImportConvert: Import using a conversion service.
// Version: tampa
func (vm) ImportConvert(session *Session, typeKey string, username string, password string, sr SRRef, remoteConfig map[string]string) (err error) {
	method := "VM.import_convert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	usernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "username"), username)
	if err != nil {
		return
	}
	passwordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "password"), password)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	remoteConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "remote_config"), remoteConfig)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, typeKeyArg, usernameArg, passwordArg, srArg, remoteConfigArg)
	return
}

// AsyncImportConvert: Import using a conversion service.
// Version: tampa
func (vm) AsyncImportConvert(session *Session, typeKey string, username string, password string, sr SRRef, remoteConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.VM.import_convert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	usernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "username"), username)
	if err != nil {
		return
	}
	passwordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "password"), password)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	remoteConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "remote_config"), remoteConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, typeKeyArg, usernameArg, passwordArg, srArg, remoteConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ImportConvert6: Import using a conversion service.
// Version: tampa
func (vm) ImportConvert6(session *Session, typeKey string, username string, password string, sr SRRef, remoteConfig map[string]string) (err error) {
	method := "VM.import_convert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	usernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "username"), username)
	if err != nil {
		return
	}
	passwordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "password"), password)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	remoteConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "remote_config"), remoteConfig)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, typeKeyArg, usernameArg, passwordArg, srArg, remoteConfigArg)
	return
}

// AsyncImportConvert6: Import using a conversion service.
// Version: tampa
func (vm) AsyncImportConvert6(session *Session, typeKey string, username string, password string, sr SRRef, remoteConfig map[string]string) (retval TaskRef, err error) {
	method := "Async.VM.import_convert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	typeKeyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "type"), typeKey)
	if err != nil {
		return
	}
	usernameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "username"), username)
	if err != nil {
		return
	}
	passwordArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "password"), password)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	remoteConfigArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "remote_config"), remoteConfig)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, typeKeyArg, usernameArg, passwordArg, srArg, remoteConfigArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Recover: Recover the VM
// Version: boston
func (vm) Recover(session *Session, self VMRef, sessionTo SessionRef, force bool) (err error) {
	method := "VM.recover"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg, forceArg)
	return
}

// AsyncRecover: Recover the VM
// Version: boston
func (vm) AsyncRecover(session *Session, self VMRef, sessionTo SessionRef, force bool) (retval TaskRef, err error) {
	method := "Async.VM.recover"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Recover4: Recover the VM
// Version: boston
func (vm) Recover4(session *Session, self VMRef, sessionTo SessionRef, force bool) (err error) {
	method := "VM.recover"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg, forceArg)
	return
}

// AsyncRecover4: Recover the VM
// Version: boston
func (vm) AsyncRecover4(session *Session, self VMRef, sessionTo SessionRef, force bool) (retval TaskRef, err error) {
	method := "Async.VM.recover"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetSRsRequiredForRecovery: List all the SR&apos;s that are required for the VM to be recovered
// Version: creedence
func (vm) GetSRsRequiredForRecovery(session *Session, self VMRef, sessionTo SessionRef) (retval []SRRef, err error) {
	method := "VM.get_SRs_required_for_recovery"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRefSet(method+" -> ", result)
	return
}

// AsyncGetSRsRequiredForRecovery: List all the SR&apos;s that are required for the VM to be recovered
// Version: creedence
func (vm) AsyncGetSRsRequiredForRecovery(session *Session, self VMRef, sessionTo SessionRef) (retval TaskRef, err error) {
	method := "Async.VM.get_SRs_required_for_recovery"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetSRsRequiredForRecovery3: List all the SR&apos;s that are required for the VM to be recovered
// Version: creedence
func (vm) GetSRsRequiredForRecovery3(session *Session, self VMRef, sessionTo SessionRef) (retval []SRRef, err error) {
	method := "VM.get_SRs_required_for_recovery"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	if err != nil {
		return
	}
	retval, err = deserializeSRRefSet(method+" -> ", result)
	return
}

// AsyncGetSRsRequiredForRecovery3: List all the SR&apos;s that are required for the VM to be recovered
// Version: creedence
func (vm) AsyncGetSRsRequiredForRecovery3(session *Session, self VMRef, sessionTo SessionRef) (retval TaskRef, err error) {
	method := "Async.VM.get_SRs_required_for_recovery"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AssertCanBeRecovered: Assert whether all SRs required to recover this VM are available.
// Version: boston
//
// Errors:
// VM_IS_PART_OF_AN_APPLIANCE - This operation is not allowed as the VM is part of an appliance.
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
func (vm) AssertCanBeRecovered(session *Session, self VMRef, sessionTo SessionRef) (err error) {
	method := "VM.assert_can_be_recovered"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	return
}

// AsyncAssertCanBeRecovered: Assert whether all SRs required to recover this VM are available.
// Version: boston
//
// Errors:
// VM_IS_PART_OF_AN_APPLIANCE - This operation is not allowed as the VM is part of an appliance.
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
func (vm) AsyncAssertCanBeRecovered(session *Session, self VMRef, sessionTo SessionRef) (retval TaskRef, err error) {
	method := "Async.VM.assert_can_be_recovered"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AssertCanBeRecovered3: Assert whether all SRs required to recover this VM are available.
// Version: boston
//
// Errors:
// VM_IS_PART_OF_AN_APPLIANCE - This operation is not allowed as the VM is part of an appliance.
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
func (vm) AssertCanBeRecovered3(session *Session, self VMRef, sessionTo SessionRef) (err error) {
	method := "VM.assert_can_be_recovered"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	return
}

// AsyncAssertCanBeRecovered3: Assert whether all SRs required to recover this VM are available.
// Version: boston
//
// Errors:
// VM_IS_PART_OF_AN_APPLIANCE - This operation is not allowed as the VM is part of an appliance.
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
func (vm) AsyncAssertCanBeRecovered3(session *Session, self VMRef, sessionTo SessionRef) (retval TaskRef, err error) {
	method := "Async.VM.assert_can_be_recovered"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	sessionToArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_to"), sessionTo)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, sessionToArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetSuspendVDI: Set this VM&apos;s suspend VDI, which must be indentical to its current one
// Version: boston
func (vm) SetSuspendVDI(session *Session, self VMRef, value VDIRef) (err error) {
	method := "VM.set_suspend_VDI"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetSuspendVDI: Set this VM&apos;s suspend VDI, which must be indentical to its current one
// Version: boston
func (vm) AsyncSetSuspendVDI(session *Session, self VMRef, value VDIRef) (retval TaskRef, err error) {
	method := "Async.VM.set_suspend_VDI"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetSuspendVDI3: Set this VM&apos;s suspend VDI, which must be indentical to its current one
// Version: boston
func (vm) SetSuspendVDI3(session *Session, self VMRef, value VDIRef) (err error) {
	method := "VM.set_suspend_VDI"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetSuspendVDI3: Set this VM&apos;s suspend VDI, which must be indentical to its current one
// Version: boston
func (vm) AsyncSetSuspendVDI3(session *Session, self VMRef, value VDIRef) (retval TaskRef, err error) {
	method := "Async.VM.set_suspend_VDI"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVDIRef(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetOrder: Set this VM&apos;s boot order
// Version: boston
func (vm) SetOrder(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_order"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetOrder: Set this VM&apos;s boot order
// Version: boston
func (vm) AsyncSetOrder(session *Session, self VMRef, value int) (retval TaskRef, err error) {
	method := "Async.VM.set_order"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetOrder3: Set this VM&apos;s boot order
// Version: boston
func (vm) SetOrder3(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_order"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetOrder3: Set this VM&apos;s boot order
// Version: boston
func (vm) AsyncSetOrder3(session *Session, self VMRef, value int) (retval TaskRef, err error) {
	method := "Async.VM.set_order"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetShutdownDelay: Set this VM&apos;s shutdown delay in seconds
// Version: boston
func (vm) SetShutdownDelay(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_shutdown_delay"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetShutdownDelay: Set this VM&apos;s shutdown delay in seconds
// Version: boston
func (vm) AsyncSetShutdownDelay(session *Session, self VMRef, value int) (retval TaskRef, err error) {
	method := "Async.VM.set_shutdown_delay"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetShutdownDelay3: Set this VM&apos;s shutdown delay in seconds
// Version: boston
func (vm) SetShutdownDelay3(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_shutdown_delay"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetShutdownDelay3: Set this VM&apos;s shutdown delay in seconds
// Version: boston
func (vm) AsyncSetShutdownDelay3(session *Session, self VMRef, value int) (retval TaskRef, err error) {
	method := "Async.VM.set_shutdown_delay"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetStartDelay: Set this VM&apos;s start delay in seconds
// Version: boston
func (vm) SetStartDelay(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_start_delay"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetStartDelay: Set this VM&apos;s start delay in seconds
// Version: boston
func (vm) AsyncSetStartDelay(session *Session, self VMRef, value int) (retval TaskRef, err error) {
	method := "Async.VM.set_start_delay"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetStartDelay3: Set this VM&apos;s start delay in seconds
// Version: boston
func (vm) SetStartDelay3(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_start_delay"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetStartDelay3: Set this VM&apos;s start delay in seconds
// Version: boston
func (vm) AsyncSetStartDelay3(session *Session, self VMRef, value int) (retval TaskRef, err error) {
	method := "Async.VM.set_start_delay"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetSnapshotSchedule: Set the value of the snapshot schedule field
// Version: falcon
func (vm) SetSnapshotSchedule(session *Session, self VMRef, value VMSSRef) (err error) {
	method := "VM.set_snapshot_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVMSSRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetSnapshotSchedule3: Set the value of the snapshot schedule field
// Version: falcon
func (vm) SetSnapshotSchedule3(session *Session, self VMRef, value VMSSRef) (err error) {
	method := "VM.set_snapshot_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVMSSRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetProtectionPolicy: Set the value of the protection_policy field
// Version: cowley
func (vm) SetProtectionPolicy(session *Session, self VMRef, value VMPPRef) (err error) {
	method := "VM.set_protection_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetProtectionPolicy3: Set the value of the protection_policy field
// Version: cowley
func (vm) SetProtectionPolicy3(session *Session, self VMRef, value VMPPRef) (err error) {
	method := "VM.set_protection_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// CopyBiosStrings: Copy the BIOS strings from the given host to this VM
// Version: midnight-ride
func (vm) CopyBiosStrings(session *Session, vm VMRef, host HostRef) (err error) {
	method := "VM.copy_bios_strings"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, hostArg)
	return
}

// AsyncCopyBiosStrings: Copy the BIOS strings from the given host to this VM
// Version: midnight-ride
func (vm) AsyncCopyBiosStrings(session *Session, vm VMRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.VM.copy_bios_strings"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CopyBiosStrings3: Copy the BIOS strings from the given host to this VM
// Version: midnight-ride
func (vm) CopyBiosStrings3(session *Session, vm VMRef, host HostRef) (err error) {
	method := "VM.copy_bios_strings"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, hostArg)
	return
}

// AsyncCopyBiosStrings3: Copy the BIOS strings from the given host to this VM
// Version: midnight-ride
func (vm) AsyncCopyBiosStrings3(session *Session, vm VMRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.VM.copy_bios_strings"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, hostArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetBiosStrings: Set custom BIOS strings to this VM. VM will be given a default set of BIOS strings, only some of which can be overridden by the supplied values. Allowed keys are: &apos;bios-vendor&apos;, &apos;bios-version&apos;, &apos;system-manufacturer&apos;, &apos;system-product-name&apos;, &apos;system-version&apos;, &apos;system-serial-number&apos;, &apos;enclosure-asset-tag&apos;, &apos;baseboard-manufacturer&apos;, &apos;baseboard-product-name&apos;, &apos;baseboard-version&apos;, &apos;baseboard-serial-number&apos;, &apos;baseboard-asset-tag&apos;, &apos;baseboard-location-in-chassis&apos;, &apos;enclosure-asset-tag&apos;
// Version: inverness
//
// Errors:
// VM_BIOS_STRINGS_ALREADY_SET - The BIOS strings for this VM have already been set and cannot be changed.
// INVALID_VALUE - The value given is invalid
func (vm) SetBiosStrings(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_bios_strings"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetBiosStrings: Set custom BIOS strings to this VM. VM will be given a default set of BIOS strings, only some of which can be overridden by the supplied values. Allowed keys are: &apos;bios-vendor&apos;, &apos;bios-version&apos;, &apos;system-manufacturer&apos;, &apos;system-product-name&apos;, &apos;system-version&apos;, &apos;system-serial-number&apos;, &apos;enclosure-asset-tag&apos;, &apos;baseboard-manufacturer&apos;, &apos;baseboard-product-name&apos;, &apos;baseboard-version&apos;, &apos;baseboard-serial-number&apos;, &apos;baseboard-asset-tag&apos;, &apos;baseboard-location-in-chassis&apos;, &apos;enclosure-asset-tag&apos;
// Version: inverness
//
// Errors:
// VM_BIOS_STRINGS_ALREADY_SET - The BIOS strings for this VM have already been set and cannot be changed.
// INVALID_VALUE - The value given is invalid
func (vm) AsyncSetBiosStrings(session *Session, self VMRef, value map[string]string) (retval TaskRef, err error) {
	method := "Async.VM.set_bios_strings"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetBiosStrings3: Set custom BIOS strings to this VM. VM will be given a default set of BIOS strings, only some of which can be overridden by the supplied values. Allowed keys are: &apos;bios-vendor&apos;, &apos;bios-version&apos;, &apos;system-manufacturer&apos;, &apos;system-product-name&apos;, &apos;system-version&apos;, &apos;system-serial-number&apos;, &apos;enclosure-asset-tag&apos;, &apos;baseboard-manufacturer&apos;, &apos;baseboard-product-name&apos;, &apos;baseboard-version&apos;, &apos;baseboard-serial-number&apos;, &apos;baseboard-asset-tag&apos;, &apos;baseboard-location-in-chassis&apos;, &apos;enclosure-asset-tag&apos;
// Version: inverness
//
// Errors:
// VM_BIOS_STRINGS_ALREADY_SET - The BIOS strings for this VM have already been set and cannot be changed.
// INVALID_VALUE - The value given is invalid
func (vm) SetBiosStrings3(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_bios_strings"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncSetBiosStrings3: Set custom BIOS strings to this VM. VM will be given a default set of BIOS strings, only some of which can be overridden by the supplied values. Allowed keys are: &apos;bios-vendor&apos;, &apos;bios-version&apos;, &apos;system-manufacturer&apos;, &apos;system-product-name&apos;, &apos;system-version&apos;, &apos;system-serial-number&apos;, &apos;enclosure-asset-tag&apos;, &apos;baseboard-manufacturer&apos;, &apos;baseboard-product-name&apos;, &apos;baseboard-version&apos;, &apos;baseboard-serial-number&apos;, &apos;baseboard-asset-tag&apos;, &apos;baseboard-location-in-chassis&apos;, &apos;enclosure-asset-tag&apos;
// Version: inverness
//
// Errors:
// VM_BIOS_STRINGS_ALREADY_SET - The BIOS strings for this VM have already been set and cannot be changed.
// INVALID_VALUE - The value given is invalid
func (vm) AsyncSetBiosStrings3(session *Session, self VMRef, value map[string]string) (retval TaskRef, err error) {
	method := "Async.VM.set_bios_strings"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// RetrieveWlbRecommendations: Returns mapping of hosts to ratings, indicating the suitability of starting the VM at that location according to wlb. Rating is replaced with an error if the VM cannot boot there.
// Version: george
func (vm) RetrieveWlbRecommendations(session *Session, vm VMRef) (retval map[HostRef][]string, err error) {
	method := "VM.retrieve_wlb_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRefToStringSetMap(method+" -> ", result)
	return
}

// AsyncRetrieveWlbRecommendations: Returns mapping of hosts to ratings, indicating the suitability of starting the VM at that location according to wlb. Rating is replaced with an error if the VM cannot boot there.
// Version: george
func (vm) AsyncRetrieveWlbRecommendations(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.retrieve_wlb_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// RetrieveWlbRecommendations2: Returns mapping of hosts to ratings, indicating the suitability of starting the VM at that location according to wlb. Rating is replaced with an error if the VM cannot boot there.
// Version: george
func (vm) RetrieveWlbRecommendations2(session *Session, vm VMRef) (retval map[HostRef][]string, err error) {
	method := "VM.retrieve_wlb_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRefToStringSetMap(method+" -> ", result)
	return
}

// AsyncRetrieveWlbRecommendations2: Returns mapping of hosts to ratings, indicating the suitability of starting the VM at that location according to wlb. Rating is replaced with an error if the VM cannot boot there.
// Version: george
func (vm) AsyncRetrieveWlbRecommendations2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.retrieve_wlb_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AssertAgile: Returns an error if the VM is not considered agile, for example, because it is tied to a resource local to a host
// Version: orlando
func (vm) AssertAgile(session *Session, self VMRef) (err error) {
	method := "VM.assert_agile"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncAssertAgile: Returns an error if the VM is not considered agile, for example, because it is tied to a resource local to a host
// Version: orlando
func (vm) AsyncAssertAgile(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.assert_agile"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AssertAgile2: Returns an error if the VM is not considered agile, for example, because it is tied to a resource local to a host
// Version: orlando
func (vm) AssertAgile2(session *Session, self VMRef) (err error) {
	method := "VM.assert_agile"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncAssertAgile2: Returns an error if the VM is not considered agile, for example, because it is tied to a resource local to a host
// Version: orlando
func (vm) AsyncAssertAgile2(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.assert_agile"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// CreateNewBlob: Create a placeholder for a named binary blob of data that is associated with this VM
// Version: tampa
func (vm) CreateNewBlob(session *Session, vm VMRef, name string, mimeType string, public bool) (retval BlobRef, err error) {
	method := "VM.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	mimeTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "mime_type"), mimeType)
	if err != nil {
		return
	}
	publicArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "public"), public)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeBlobRef(method+" -> ", result)
	return
}

// AsyncCreateNewBlob: Create a placeholder for a named binary blob of data that is associated with this VM
// Version: tampa
func (vm) AsyncCreateNewBlob(session *Session, vm VMRef, name string, mimeType string, public bool) (retval TaskRef, err error) {
	method := "Async.VM.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	mimeTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "mime_type"), mimeType)
	if err != nil {
		return
	}
	publicArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "public"), public)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateNewBlob5: Create a placeholder for a named binary blob of data that is associated with this VM
// Version: tampa
func (vm) CreateNewBlob5(session *Session, vm VMRef, name string, mimeType string, public bool) (retval BlobRef, err error) {
	method := "VM.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	mimeTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "mime_type"), mimeType)
	if err != nil {
		return
	}
	publicArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "public"), public)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeBlobRef(method+" -> ", result)
	return
}

// AsyncCreateNewBlob5: Create a placeholder for a named binary blob of data that is associated with this VM
// Version: tampa
func (vm) AsyncCreateNewBlob5(session *Session, vm VMRef, name string, mimeType string, public bool) (retval TaskRef, err error) {
	method := "Async.VM.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	mimeTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "mime_type"), mimeType)
	if err != nil {
		return
	}
	publicArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "public"), public)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, nameArg, mimeTypeArg, publicArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CreateNewBlob4: Create a placeholder for a named binary blob of data that is associated with this VM
// Version: orlando
func (vm) CreateNewBlob4(session *Session, vm VMRef, name string, mimeType string) (retval BlobRef, err error) {
	method := "VM.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	mimeTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "mime_type"), mimeType)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, nameArg, mimeTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeBlobRef(method+" -> ", result)
	return
}

// AsyncCreateNewBlob4: Create a placeholder for a named binary blob of data that is associated with this VM
// Version: orlando
func (vm) AsyncCreateNewBlob4(session *Session, vm VMRef, name string, mimeType string) (retval TaskRef, err error) {
	method := "Async.VM.create_new_blob"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	nameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "name"), name)
	if err != nil {
		return
	}
	mimeTypeArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "mime_type"), mimeType)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, nameArg, mimeTypeArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AssertCanBootHere: Returns an error if the VM could not boot on this host for some reason
// Version: rio
//
// Errors:
// HOST_NOT_ENOUGH_FREE_MEMORY - Not enough server memory is available to perform this operation.
// HOST_NOT_ENOUGH_PCPUS - The host does not have enough pCPUs to run the VM. It needs at least as many as the VM has vCPUs.
// NETWORK_SRIOV_INSUFFICIENT_CAPACITY - There is insufficient capacity for VF reservation
// HOST_NOT_LIVE - This operation cannot be completed as the server is not live.
// HOST_DISABLED - The specified server is disabled.
// HOST_CANNOT_ATTACH_NETWORK - Server cannot attach network (in the case of NIC bonding, this may be because attaching the network on this server would require other networks - that are currently active - to be taken down).
// VM_HVM_REQUIRED - HVM is required for this operation
// VM_REQUIRES_GPU - You attempted to run a VM on a host which doesn&apos;t have a pGPU available in the GPU group needed by the VM. The VM has a vGPU attached to this GPU group.
// VM_REQUIRES_IOMMU - You attempted to run a VM on a host which doesn&apos;t have I/O virtualization (IOMMU/VT-d) enabled, which is needed by the VM.
// VM_REQUIRES_NETWORK - You attempted to run a VM on a host which doesn&apos;t have a PIF on a Network needed by the VM. The VM has at least one VIF attached to the Network.
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
// VM_REQUIRES_VGPU - You attempted to run a VM on a host on which the vGPU required by the VM cannot be allocated on any pGPUs in the GPU_group needed by the VM.
// VM_HOST_INCOMPATIBLE_VERSION - This VM operation cannot be performed on an older-versioned host during an upgrade.
// VM_HOST_INCOMPATIBLE_VIRTUAL_HARDWARE_PLATFORM_VERSION - You attempted to run a VM on a host that cannot provide the VM&apos;s required Virtual Hardware Platform version.
// INVALID_VALUE - The value given is invalid
// MEMORY_CONSTRAINT_VIOLATION - The dynamic memory range does not satisfy the following constraint.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VALUE_NOT_SUPPORTED - You attempted to set a value that is not supported by this implementation. The fully-qualified field name and the value that you tried to set are returned. Also returned is a developer-only diagnostic reason.
// VM_INCOMPATIBLE_WITH_THIS_HOST - The VM is incompatible with the CPU features of this host.
func (vm) AssertCanBootHere(session *Session, self VMRef, host HostRef) (err error) {
	method := "VM.assert_can_boot_here"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncAssertCanBootHere: Returns an error if the VM could not boot on this host for some reason
// Version: rio
//
// Errors:
// HOST_NOT_ENOUGH_FREE_MEMORY - Not enough server memory is available to perform this operation.
// HOST_NOT_ENOUGH_PCPUS - The host does not have enough pCPUs to run the VM. It needs at least as many as the VM has vCPUs.
// NETWORK_SRIOV_INSUFFICIENT_CAPACITY - There is insufficient capacity for VF reservation
// HOST_NOT_LIVE - This operation cannot be completed as the server is not live.
// HOST_DISABLED - The specified server is disabled.
// HOST_CANNOT_ATTACH_NETWORK - Server cannot attach network (in the case of NIC bonding, this may be because attaching the network on this server would require other networks - that are currently active - to be taken down).
// VM_HVM_REQUIRED - HVM is required for this operation
// VM_REQUIRES_GPU - You attempted to run a VM on a host which doesn&apos;t have a pGPU available in the GPU group needed by the VM. The VM has a vGPU attached to this GPU group.
// VM_REQUIRES_IOMMU - You attempted to run a VM on a host which doesn&apos;t have I/O virtualization (IOMMU/VT-d) enabled, which is needed by the VM.
// VM_REQUIRES_NETWORK - You attempted to run a VM on a host which doesn&apos;t have a PIF on a Network needed by the VM. The VM has at least one VIF attached to the Network.
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
// VM_REQUIRES_VGPU - You attempted to run a VM on a host on which the vGPU required by the VM cannot be allocated on any pGPUs in the GPU_group needed by the VM.
// VM_HOST_INCOMPATIBLE_VERSION - This VM operation cannot be performed on an older-versioned host during an upgrade.
// VM_HOST_INCOMPATIBLE_VIRTUAL_HARDWARE_PLATFORM_VERSION - You attempted to run a VM on a host that cannot provide the VM&apos;s required Virtual Hardware Platform version.
// INVALID_VALUE - The value given is invalid
// MEMORY_CONSTRAINT_VIOLATION - The dynamic memory range does not satisfy the following constraint.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VALUE_NOT_SUPPORTED - You attempted to set a value that is not supported by this implementation. The fully-qualified field name and the value that you tried to set are returned. Also returned is a developer-only diagnostic reason.
// VM_INCOMPATIBLE_WITH_THIS_HOST - The VM is incompatible with the CPU features of this host.
func (vm) AsyncAssertCanBootHere(session *Session, self VMRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.VM.assert_can_boot_here"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AssertCanBootHere3: Returns an error if the VM could not boot on this host for some reason
// Version: rio
//
// Errors:
// HOST_NOT_ENOUGH_FREE_MEMORY - Not enough server memory is available to perform this operation.
// HOST_NOT_ENOUGH_PCPUS - The host does not have enough pCPUs to run the VM. It needs at least as many as the VM has vCPUs.
// NETWORK_SRIOV_INSUFFICIENT_CAPACITY - There is insufficient capacity for VF reservation
// HOST_NOT_LIVE - This operation cannot be completed as the server is not live.
// HOST_DISABLED - The specified server is disabled.
// HOST_CANNOT_ATTACH_NETWORK - Server cannot attach network (in the case of NIC bonding, this may be because attaching the network on this server would require other networks - that are currently active - to be taken down).
// VM_HVM_REQUIRED - HVM is required for this operation
// VM_REQUIRES_GPU - You attempted to run a VM on a host which doesn&apos;t have a pGPU available in the GPU group needed by the VM. The VM has a vGPU attached to this GPU group.
// VM_REQUIRES_IOMMU - You attempted to run a VM on a host which doesn&apos;t have I/O virtualization (IOMMU/VT-d) enabled, which is needed by the VM.
// VM_REQUIRES_NETWORK - You attempted to run a VM on a host which doesn&apos;t have a PIF on a Network needed by the VM. The VM has at least one VIF attached to the Network.
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
// VM_REQUIRES_VGPU - You attempted to run a VM on a host on which the vGPU required by the VM cannot be allocated on any pGPUs in the GPU_group needed by the VM.
// VM_HOST_INCOMPATIBLE_VERSION - This VM operation cannot be performed on an older-versioned host during an upgrade.
// VM_HOST_INCOMPATIBLE_VIRTUAL_HARDWARE_PLATFORM_VERSION - You attempted to run a VM on a host that cannot provide the VM&apos;s required Virtual Hardware Platform version.
// INVALID_VALUE - The value given is invalid
// MEMORY_CONSTRAINT_VIOLATION - The dynamic memory range does not satisfy the following constraint.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VALUE_NOT_SUPPORTED - You attempted to set a value that is not supported by this implementation. The fully-qualified field name and the value that you tried to set are returned. Also returned is a developer-only diagnostic reason.
// VM_INCOMPATIBLE_WITH_THIS_HOST - The VM is incompatible with the CPU features of this host.
func (vm) AssertCanBootHere3(session *Session, self VMRef, host HostRef) (err error) {
	method := "VM.assert_can_boot_here"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncAssertCanBootHere3: Returns an error if the VM could not boot on this host for some reason
// Version: rio
//
// Errors:
// HOST_NOT_ENOUGH_FREE_MEMORY - Not enough server memory is available to perform this operation.
// HOST_NOT_ENOUGH_PCPUS - The host does not have enough pCPUs to run the VM. It needs at least as many as the VM has vCPUs.
// NETWORK_SRIOV_INSUFFICIENT_CAPACITY - There is insufficient capacity for VF reservation
// HOST_NOT_LIVE - This operation cannot be completed as the server is not live.
// HOST_DISABLED - The specified server is disabled.
// HOST_CANNOT_ATTACH_NETWORK - Server cannot attach network (in the case of NIC bonding, this may be because attaching the network on this server would require other networks - that are currently active - to be taken down).
// VM_HVM_REQUIRED - HVM is required for this operation
// VM_REQUIRES_GPU - You attempted to run a VM on a host which doesn&apos;t have a pGPU available in the GPU group needed by the VM. The VM has a vGPU attached to this GPU group.
// VM_REQUIRES_IOMMU - You attempted to run a VM on a host which doesn&apos;t have I/O virtualization (IOMMU/VT-d) enabled, which is needed by the VM.
// VM_REQUIRES_NETWORK - You attempted to run a VM on a host which doesn&apos;t have a PIF on a Network needed by the VM. The VM has at least one VIF attached to the Network.
// VM_REQUIRES_SR - You attempted to run a VM on a host which doesn&apos;t have access to an SR needed by the VM. The VM has at least one VBD attached to a VDI in the SR.
// VM_REQUIRES_VGPU - You attempted to run a VM on a host on which the vGPU required by the VM cannot be allocated on any pGPUs in the GPU_group needed by the VM.
// VM_HOST_INCOMPATIBLE_VERSION - This VM operation cannot be performed on an older-versioned host during an upgrade.
// VM_HOST_INCOMPATIBLE_VIRTUAL_HARDWARE_PLATFORM_VERSION - You attempted to run a VM on a host that cannot provide the VM&apos;s required Virtual Hardware Platform version.
// INVALID_VALUE - The value given is invalid
// MEMORY_CONSTRAINT_VIOLATION - The dynamic memory range does not satisfy the following constraint.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VALUE_NOT_SUPPORTED - You attempted to set a value that is not supported by this implementation. The fully-qualified field name and the value that you tried to set are returned. Also returned is a developer-only diagnostic reason.
// VM_INCOMPATIBLE_WITH_THIS_HOST - The VM is incompatible with the CPU features of this host.
func (vm) AsyncAssertCanBootHere3(session *Session, self VMRef, host HostRef) (retval TaskRef, err error) {
	method := "Async.VM.assert_can_boot_here"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPossibleHosts: Return the list of hosts on which this VM may run.
// Version: rio
func (vm) GetPossibleHosts(session *Session, vm VMRef) (retval []HostRef, err error) {
	method := "VM.get_possible_hosts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRefSet(method+" -> ", result)
	return
}

// AsyncGetPossibleHosts: Return the list of hosts on which this VM may run.
// Version: rio
func (vm) AsyncGetPossibleHosts(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.get_possible_hosts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetPossibleHosts2: Return the list of hosts on which this VM may run.
// Version: rio
func (vm) GetPossibleHosts2(session *Session, vm VMRef) (retval []HostRef, err error) {
	method := "VM.get_possible_hosts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeHostRefSet(method+" -> ", result)
	return
}

// AsyncGetPossibleHosts2: Return the list of hosts on which this VM may run.
// Version: rio
func (vm) AsyncGetPossibleHosts2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.get_possible_hosts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// GetAllowedVIFDevices: Returns a list of the allowed values that a VIF device field can take
// Version: rio
func (vm) GetAllowedVIFDevices(session *Session, vm VMRef) (retval []string, err error) {
	method := "VM.get_allowed_VIF_devices"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetAllowedVIFDevices2: Returns a list of the allowed values that a VIF device field can take
// Version: rio
func (vm) GetAllowedVIFDevices2(session *Session, vm VMRef) (retval []string, err error) {
	method := "VM.get_allowed_VIF_devices"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetAllowedVBDDevices: Returns a list of the allowed values that a VBD device field can take
// Version: rio
func (vm) GetAllowedVBDDevices(session *Session, vm VMRef) (retval []string, err error) {
	method := "VM.get_allowed_VBD_devices"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetAllowedVBDDevices2: Returns a list of the allowed values that a VBD device field can take
// Version: rio
func (vm) GetAllowedVBDDevices2(session *Session, vm VMRef) (retval []string, err error) {
	method := "VM.get_allowed_VBD_devices"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// UpdateAllowedOperations: Recomputes the list of acceptable operations
// Version: rio
func (vm) UpdateAllowedOperations(session *Session, self VMRef) (err error) {
	method := "VM.update_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncUpdateAllowedOperations: Recomputes the list of acceptable operations
// Version: rio
func (vm) AsyncUpdateAllowedOperations(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.update_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// UpdateAllowedOperations2: Recomputes the list of acceptable operations
// Version: rio
func (vm) UpdateAllowedOperations2(session *Session, self VMRef) (err error) {
	method := "VM.update_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncUpdateAllowedOperations2: Recomputes the list of acceptable operations
// Version: rio
func (vm) AsyncUpdateAllowedOperations2(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.update_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AssertOperationValid: Check to see whether this operation is acceptable in the current state of the system, raising an error if the operation is invalid for some reason
// Version: rio
func (vm) AssertOperationValid(session *Session, self VMRef, op VMOperations) (err error) {
	method := "VM.assert_operation_valid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	opArg, err := serializeEnumVMOperations(fmt.Sprintf("%s(%s)", method, "op"), op)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, opArg)
	return
}

// AsyncAssertOperationValid: Check to see whether this operation is acceptable in the current state of the system, raising an error if the operation is invalid for some reason
// Version: rio
func (vm) AsyncAssertOperationValid(session *Session, self VMRef, op VMOperations) (retval TaskRef, err error) {
	method := "Async.VM.assert_operation_valid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	opArg, err := serializeEnumVMOperations(fmt.Sprintf("%s(%s)", method, "op"), op)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, opArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AssertOperationValid3: Check to see whether this operation is acceptable in the current state of the system, raising an error if the operation is invalid for some reason
// Version: rio
func (vm) AssertOperationValid3(session *Session, self VMRef, op VMOperations) (err error) {
	method := "VM.assert_operation_valid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	opArg, err := serializeEnumVMOperations(fmt.Sprintf("%s(%s)", method, "op"), op)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, opArg)
	return
}

// AsyncAssertOperationValid3: Check to see whether this operation is acceptable in the current state of the system, raising an error if the operation is invalid for some reason
// Version: rio
func (vm) AsyncAssertOperationValid3(session *Session, self VMRef, op VMOperations) (retval TaskRef, err error) {
	method := "Async.VM.assert_operation_valid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	opArg, err := serializeEnumVMOperations(fmt.Sprintf("%s(%s)", method, "op"), op)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, opArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ForgetDataSourceArchives: Forget the recorded statistics related to the specified data source
// Version: orlando
func (vm) ForgetDataSourceArchives(session *Session, self VMRef, dataSource string) (err error) {
	method := "VM.forget_data_source_archives"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, dataSourceArg)
	return
}

// ForgetDataSourceArchives3: Forget the recorded statistics related to the specified data source
// Version: orlando
func (vm) ForgetDataSourceArchives3(session *Session, self VMRef, dataSource string) (err error) {
	method := "VM.forget_data_source_archives"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, dataSourceArg)
	return
}

// QueryDataSource: Query the latest value of the specified data source
// Version: orlando
func (vm) QueryDataSource(session *Session, self VMRef, dataSource string) (retval float64, err error) {
	method := "VM.query_data_source"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, dataSourceArg)
	if err != nil {
		return
	}
	retval, err = deserializeFloat(method+" -> ", result)
	return
}

// QueryDataSource3: Query the latest value of the specified data source
// Version: orlando
func (vm) QueryDataSource3(session *Session, self VMRef, dataSource string) (retval float64, err error) {
	method := "VM.query_data_source"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, dataSourceArg)
	if err != nil {
		return
	}
	retval, err = deserializeFloat(method+" -> ", result)
	return
}

// RecordDataSource: Start recording the specified data source
// Version: orlando
func (vm) RecordDataSource(session *Session, self VMRef, dataSource string) (err error) {
	method := "VM.record_data_source"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, dataSourceArg)
	return
}

// RecordDataSource3: Start recording the specified data source
// Version: orlando
func (vm) RecordDataSource3(session *Session, self VMRef, dataSource string) (err error) {
	method := "VM.record_data_source"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	dataSourceArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "data_source"), dataSource)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, dataSourceArg)
	return
}

// GetDataSources:
// Version: orlando
func (vm) GetDataSources(session *Session, self VMRef) (retval []DataSourceRecord, err error) {
	method := "VM.get_data_sources"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDataSourceRecordSet(method+" -> ", result)
	return
}

// GetDataSources2:
// Version: orlando
func (vm) GetDataSources2(session *Session, self VMRef) (retval []DataSourceRecord, err error) {
	method := "VM.get_data_sources"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeDataSourceRecordSet(method+" -> ", result)
	return
}

// GetBootRecord: Returns a record describing the VM&apos;s dynamic state, initialised when the VM boots and updated to reflect runtime configuration changes, for example CPU hotplug
// Version: rio
func (vm) GetBootRecord(session *Session, self VMRef) (retval VMRecord, err error) {
	method := "VM.get_boot_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRecord(method+" -> ", result)
	return
}

// GetBootRecord2: Returns a record describing the VM&apos;s dynamic state, initialised when the VM boots and updated to reflect runtime configuration changes, for example CPU hotplug
// Version: rio
func (vm) GetBootRecord2(session *Session, self VMRef) (retval VMRecord, err error) {
	method := "VM.get_boot_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRecord(method+" -> ", result)
	return
}

// AssertCanMigrate: Assert whether a VM can be migrated to the specified destination.
// Version: inverness
//
// Errors:
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AssertCanMigrate(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string, vgpuMap map[VGPURef]GPUGroupRef) (err error) {
	method := "VM.assert_can_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	vgpuMapArg, err := serializeVGPURefToGPUGroupRefMap(fmt.Sprintf("%s(%s)", method, "vgpu_map"), vgpuMap)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg, vgpuMapArg)
	return
}

// AsyncAssertCanMigrate: Assert whether a VM can be migrated to the specified destination.
// Version: inverness
//
// Errors:
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncAssertCanMigrate(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string, vgpuMap map[VGPURef]GPUGroupRef) (retval TaskRef, err error) {
	method := "Async.VM.assert_can_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	vgpuMapArg, err := serializeVGPURefToGPUGroupRefMap(fmt.Sprintf("%s(%s)", method, "vgpu_map"), vgpuMap)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg, vgpuMapArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AssertCanMigrate8: Assert whether a VM can be migrated to the specified destination.
// Version: inverness
//
// Errors:
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AssertCanMigrate8(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string, vgpuMap map[VGPURef]GPUGroupRef) (err error) {
	method := "VM.assert_can_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	vgpuMapArg, err := serializeVGPURefToGPUGroupRefMap(fmt.Sprintf("%s(%s)", method, "vgpu_map"), vgpuMap)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg, vgpuMapArg)
	return
}

// AsyncAssertCanMigrate8: Assert whether a VM can be migrated to the specified destination.
// Version: inverness
//
// Errors:
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncAssertCanMigrate8(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string, vgpuMap map[VGPURef]GPUGroupRef) (retval TaskRef, err error) {
	method := "Async.VM.assert_can_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	vgpuMapArg, err := serializeVGPURefToGPUGroupRefMap(fmt.Sprintf("%s(%s)", method, "vgpu_map"), vgpuMap)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg, vgpuMapArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AssertCanMigrate7: Assert whether a VM can be migrated to the specified destination.
// Version: tampa
//
// Errors:
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AssertCanMigrate7(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string) (err error) {
	method := "VM.assert_can_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg)
	return
}

// AsyncAssertCanMigrate7: Assert whether a VM can be migrated to the specified destination.
// Version: tampa
//
// Errors:
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncAssertCanMigrate7(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string) (retval TaskRef, err error) {
	method := "Async.VM.assert_can_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// MigrateSend: Migrate the VM to another host.  This can only be called when the specified VM is in the Running state.
// Version: 26.14.0-next
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) MigrateSend(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string, vgpuMap map[VGPURef]GPUGroupRef, vdiFormatMap map[VDIRef]string) (retval VMRef, err error) {
	method := "VM.migrate_send"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	vgpuMapArg, err := serializeVGPURefToGPUGroupRefMap(fmt.Sprintf("%s(%s)", method, "vgpu_map"), vgpuMap)
	if err != nil {
		return
	}
	vdiFormatMapArg, err := serializeVDIRefToStringMap(fmt.Sprintf("%s(%s)", method, "vdi_format_map"), vdiFormatMap)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg, vgpuMapArg, vdiFormatMapArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncMigrateSend: Migrate the VM to another host.  This can only be called when the specified VM is in the Running state.
// Version: 26.14.0-next
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncMigrateSend(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string, vgpuMap map[VGPURef]GPUGroupRef, vdiFormatMap map[VDIRef]string) (retval TaskRef, err error) {
	method := "Async.VM.migrate_send"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	vgpuMapArg, err := serializeVGPURefToGPUGroupRefMap(fmt.Sprintf("%s(%s)", method, "vgpu_map"), vgpuMap)
	if err != nil {
		return
	}
	vdiFormatMapArg, err := serializeVDIRefToStringMap(fmt.Sprintf("%s(%s)", method, "vdi_format_map"), vdiFormatMap)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg, vgpuMapArg, vdiFormatMapArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// MigrateSend9: Migrate the VM to another host.  This can only be called when the specified VM is in the Running state.
// Version: 26.14.0-next
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) MigrateSend9(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string, vgpuMap map[VGPURef]GPUGroupRef, vdiFormatMap map[VDIRef]string) (retval VMRef, err error) {
	method := "VM.migrate_send"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	vgpuMapArg, err := serializeVGPURefToGPUGroupRefMap(fmt.Sprintf("%s(%s)", method, "vgpu_map"), vgpuMap)
	if err != nil {
		return
	}
	vdiFormatMapArg, err := serializeVDIRefToStringMap(fmt.Sprintf("%s(%s)", method, "vdi_format_map"), vdiFormatMap)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg, vgpuMapArg, vdiFormatMapArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncMigrateSend9: Migrate the VM to another host.  This can only be called when the specified VM is in the Running state.
// Version: 26.14.0-next
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncMigrateSend9(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string, vgpuMap map[VGPURef]GPUGroupRef, vdiFormatMap map[VDIRef]string) (retval TaskRef, err error) {
	method := "Async.VM.migrate_send"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	vgpuMapArg, err := serializeVGPURefToGPUGroupRefMap(fmt.Sprintf("%s(%s)", method, "vgpu_map"), vgpuMap)
	if err != nil {
		return
	}
	vdiFormatMapArg, err := serializeVDIRefToStringMap(fmt.Sprintf("%s(%s)", method, "vdi_format_map"), vdiFormatMap)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg, vgpuMapArg, vdiFormatMapArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// MigrateSend8: Migrate the VM to another host.  This can only be called when the specified VM is in the Running state.
// Version: inverness
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) MigrateSend8(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string, vgpuMap map[VGPURef]GPUGroupRef) (retval VMRef, err error) {
	method := "VM.migrate_send"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	vgpuMapArg, err := serializeVGPURefToGPUGroupRefMap(fmt.Sprintf("%s(%s)", method, "vgpu_map"), vgpuMap)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg, vgpuMapArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncMigrateSend8: Migrate the VM to another host.  This can only be called when the specified VM is in the Running state.
// Version: inverness
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncMigrateSend8(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string, vgpuMap map[VGPURef]GPUGroupRef) (retval TaskRef, err error) {
	method := "Async.VM.migrate_send"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	vgpuMapArg, err := serializeVGPURefToGPUGroupRefMap(fmt.Sprintf("%s(%s)", method, "vgpu_map"), vgpuMap)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg, vgpuMapArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// MigrateSend7: Migrate the VM to another host.  This can only be called when the specified VM is in the Running state.
// Version: tampa
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) MigrateSend7(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string) (retval VMRef, err error) {
	method := "VM.migrate_send"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncMigrateSend7: Migrate the VM to another host.  This can only be called when the specified VM is in the Running state.
// Version: tampa
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncMigrateSend7(session *Session, vm VMRef, dest map[string]string, live bool, vdiMap map[VDIRef]SRRef, vifMap map[VIFRef]NetworkRef, options map[string]string) (retval TaskRef, err error) {
	method := "Async.VM.migrate_send"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	destArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "dest"), dest)
	if err != nil {
		return
	}
	liveArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "live"), live)
	if err != nil {
		return
	}
	vdiMapArg, err := serializeVDIRefToSRRefMap(fmt.Sprintf("%s(%s)", method, "vdi_map"), vdiMap)
	if err != nil {
		return
	}
	vifMapArg, err := serializeVIFRefToNetworkRefMap(fmt.Sprintf("%s(%s)", method, "vif_map"), vifMap)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, destArg, liveArg, vdiMapArg, vifMapArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// MaximiseMemory: Returns the maximum amount of guest memory which will fit, together with overheads, in the supplied amount of physical memory. If &apos;exact&apos; is true then an exact calculation is performed using the VM&apos;s current settings. If &apos;exact&apos; is false then a more conservative approximation is used
// Version: miami
func (vm) MaximiseMemory(session *Session, self VMRef, total int, approximate bool) (retval int, err error) {
	method := "VM.maximise_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	totalArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "total"), total)
	if err != nil {
		return
	}
	approximateArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "approximate"), approximate)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, totalArg, approximateArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// AsyncMaximiseMemory: Returns the maximum amount of guest memory which will fit, together with overheads, in the supplied amount of physical memory. If &apos;exact&apos; is true then an exact calculation is performed using the VM&apos;s current settings. If &apos;exact&apos; is false then a more conservative approximation is used
// Version: miami
func (vm) AsyncMaximiseMemory(session *Session, self VMRef, total int, approximate bool) (retval TaskRef, err error) {
	method := "Async.VM.maximise_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	totalArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "total"), total)
	if err != nil {
		return
	}
	approximateArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "approximate"), approximate)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, totalArg, approximateArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// MaximiseMemory4: Returns the maximum amount of guest memory which will fit, together with overheads, in the supplied amount of physical memory. If &apos;exact&apos; is true then an exact calculation is performed using the VM&apos;s current settings. If &apos;exact&apos; is false then a more conservative approximation is used
// Version: miami
func (vm) MaximiseMemory4(session *Session, self VMRef, total int, approximate bool) (retval int, err error) {
	method := "VM.maximise_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	totalArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "total"), total)
	if err != nil {
		return
	}
	approximateArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "approximate"), approximate)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, totalArg, approximateArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// AsyncMaximiseMemory4: Returns the maximum amount of guest memory which will fit, together with overheads, in the supplied amount of physical memory. If &apos;exact&apos; is true then an exact calculation is performed using the VM&apos;s current settings. If &apos;exact&apos; is false then a more conservative approximation is used
// Version: miami
func (vm) AsyncMaximiseMemory4(session *Session, self VMRef, total int, approximate bool) (retval TaskRef, err error) {
	method := "Async.VM.maximise_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	totalArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "total"), total)
	if err != nil {
		return
	}
	approximateArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "approximate"), approximate)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, totalArg, approximateArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SendTrigger: Send the named trigger to this VM.  This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) SendTrigger(session *Session, vm VMRef, trigger string) (err error) {
	method := "VM.send_trigger"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	triggerArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "trigger"), trigger)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, triggerArg)
	return
}

// AsyncSendTrigger: Send the named trigger to this VM.  This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) AsyncSendTrigger(session *Session, vm VMRef, trigger string) (retval TaskRef, err error) {
	method := "Async.VM.send_trigger"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	triggerArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "trigger"), trigger)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, triggerArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SendTrigger3: Send the named trigger to this VM.  This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) SendTrigger3(session *Session, vm VMRef, trigger string) (err error) {
	method := "VM.send_trigger"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	triggerArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "trigger"), trigger)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, triggerArg)
	return
}

// AsyncSendTrigger3: Send the named trigger to this VM.  This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) AsyncSendTrigger3(session *Session, vm VMRef, trigger string) (retval TaskRef, err error) {
	method := "Async.VM.send_trigger"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	triggerArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "trigger"), trigger)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, triggerArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SendSysrq: Send the given key as a sysrq to this VM.  The key is specified as a single character (a String of length 1).  This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) SendSysrq(session *Session, vm VMRef, key string) (err error) {
	method := "VM.send_sysrq"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, keyArg)
	return
}

// AsyncSendSysrq: Send the given key as a sysrq to this VM.  The key is specified as a single character (a String of length 1).  This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) AsyncSendSysrq(session *Session, vm VMRef, key string) (retval TaskRef, err error) {
	method := "Async.VM.send_sysrq"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, keyArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SendSysrq3: Send the given key as a sysrq to this VM.  The key is specified as a single character (a String of length 1).  This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) SendSysrq3(session *Session, vm VMRef, key string) (err error) {
	method := "VM.send_sysrq"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, keyArg)
	return
}

// AsyncSendSysrq3: Send the given key as a sysrq to this VM.  The key is specified as a single character (a String of length 1).  This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) AsyncSendSysrq3(session *Session, vm VMRef, key string) (retval TaskRef, err error) {
	method := "Async.VM.send_sysrq"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	keyArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "key"), key)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, keyArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetVCPUsAtStartup: Set the number of startup VCPUs for a halted VM
// Version: midnight-ride
func (vm) SetVCPUsAtStartup(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_VCPUs_at_startup"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetVCPUsAtStartup3: Set the number of startup VCPUs for a halted VM
// Version: midnight-ride
func (vm) SetVCPUsAtStartup3(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_VCPUs_at_startup"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetVCPUsMax: Set the maximum number of VCPUs for a halted VM
// Version: midnight-ride
func (vm) SetVCPUsMax(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_VCPUs_max"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetVCPUsMax3: Set the maximum number of VCPUs for a halted VM
// Version: midnight-ride
func (vm) SetVCPUsMax3(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_VCPUs_max"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetShadowMultiplierLive: Set the shadow memory multiplier on a running VM
// Version: rio
func (vm) SetShadowMultiplierLive(session *Session, self VMRef, multiplier float64) (err error) {
	method := "VM.set_shadow_multiplier_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	multiplierArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "multiplier"), multiplier)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, multiplierArg)
	return
}

// AsyncSetShadowMultiplierLive: Set the shadow memory multiplier on a running VM
// Version: rio
func (vm) AsyncSetShadowMultiplierLive(session *Session, self VMRef, multiplier float64) (retval TaskRef, err error) {
	method := "Async.VM.set_shadow_multiplier_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	multiplierArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "multiplier"), multiplier)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, multiplierArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetShadowMultiplierLive3: Set the shadow memory multiplier on a running VM
// Version: rio
func (vm) SetShadowMultiplierLive3(session *Session, self VMRef, multiplier float64) (err error) {
	method := "VM.set_shadow_multiplier_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	multiplierArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "multiplier"), multiplier)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, multiplierArg)
	return
}

// AsyncSetShadowMultiplierLive3: Set the shadow memory multiplier on a running VM
// Version: rio
func (vm) AsyncSetShadowMultiplierLive3(session *Session, self VMRef, multiplier float64) (retval TaskRef, err error) {
	method := "Async.VM.set_shadow_multiplier_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	multiplierArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "multiplier"), multiplier)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, multiplierArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetHVMShadowMultiplier: Set the shadow memory multiplier on a halted VM
// Version: midnight-ride
func (vm) SetHVMShadowMultiplier(session *Session, self VMRef, value float64) (err error) {
	method := "VM.set_HVM_shadow_multiplier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetHVMShadowMultiplier3: Set the shadow memory multiplier on a halted VM
// Version: midnight-ride
func (vm) SetHVMShadowMultiplier3(session *Session, self VMRef, value float64) (err error) {
	method := "VM.set_HVM_shadow_multiplier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeFloat(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// GetCooperative: Return true if the VM is currently &apos;co-operative&apos; i.e. is expected to reach a balloon target and actually has done
// Version: midnight-ride
func (vm) GetCooperative(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_cooperative"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncGetCooperative: Return true if the VM is currently &apos;co-operative&apos; i.e. is expected to reach a balloon target and actually has done
// Version: midnight-ride
func (vm) AsyncGetCooperative(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.get_cooperative"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCooperative2: Return true if the VM is currently &apos;co-operative&apos; i.e. is expected to reach a balloon target and actually has done
// Version: midnight-ride
func (vm) GetCooperative2(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_cooperative"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncGetCooperative2: Return true if the VM is currently &apos;co-operative&apos; i.e. is expected to reach a balloon target and actually has done
// Version: midnight-ride
func (vm) AsyncGetCooperative2(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.get_cooperative"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// WaitMemoryTargetLive: Wait for a running VM to reach its current memory target
// Version: orlando
func (vm) WaitMemoryTargetLive(session *Session, self VMRef) (err error) {
	method := "VM.wait_memory_target_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncWaitMemoryTargetLive: Wait for a running VM to reach its current memory target
// Version: orlando
func (vm) AsyncWaitMemoryTargetLive(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.wait_memory_target_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// WaitMemoryTargetLive2: Wait for a running VM to reach its current memory target
// Version: orlando
func (vm) WaitMemoryTargetLive2(session *Session, self VMRef) (err error) {
	method := "VM.wait_memory_target_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncWaitMemoryTargetLive2: Wait for a running VM to reach its current memory target
// Version: orlando
func (vm) AsyncWaitMemoryTargetLive2(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.wait_memory_target_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetMemoryTargetLive: Set the memory target for a running VM
// Version: rio
func (vm) SetMemoryTargetLive(session *Session, self VMRef, target int) (err error) {
	method := "VM.set_memory_target_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	targetArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "target"), target)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, targetArg)
	return
}

// AsyncSetMemoryTargetLive: Set the memory target for a running VM
// Version: rio
func (vm) AsyncSetMemoryTargetLive(session *Session, self VMRef, target int) (retval TaskRef, err error) {
	method := "Async.VM.set_memory_target_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	targetArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "target"), target)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, targetArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMemoryTargetLive3: Set the memory target for a running VM
// Version: rio
func (vm) SetMemoryTargetLive3(session *Session, self VMRef, target int) (err error) {
	method := "VM.set_memory_target_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	targetArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "target"), target)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, targetArg)
	return
}

// AsyncSetMemoryTargetLive3: Set the memory target for a running VM
// Version: rio
func (vm) AsyncSetMemoryTargetLive3(session *Session, self VMRef, target int) (retval TaskRef, err error) {
	method := "Async.VM.set_memory_target_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	targetArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "target"), target)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, targetArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMemory: Set the memory allocation of this VM. Sets all of memory_static_max, memory_dynamic_min, and memory_dynamic_max to the given value, and leaves memory_static_min untouched.
// Version: ely
func (vm) SetMemory(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetMemory: Set the memory allocation of this VM. Sets all of memory_static_max, memory_dynamic_min, and memory_dynamic_max to the given value, and leaves memory_static_min untouched.
// Version: ely
func (vm) AsyncSetMemory(session *Session, self VMRef, value int) (retval TaskRef, err error) {
	method := "Async.VM.set_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetMemory3: Set the memory allocation of this VM. Sets all of memory_static_max, memory_dynamic_min, and memory_dynamic_max to the given value, and leaves memory_static_min untouched.
// Version: ely
func (vm) SetMemory3(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// AsyncSetMemory3: Set the memory allocation of this VM. Sets all of memory_static_max, memory_dynamic_min, and memory_dynamic_max to the given value, and leaves memory_static_min untouched.
// Version: ely
func (vm) AsyncSetMemory3(session *Session, self VMRef, value int) (retval TaskRef, err error) {
	method := "Async.VM.set_memory"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
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

// SetMemoryLimits: Set the memory limits of this VM.
// Version: midnight-ride
func (vm) SetMemoryLimits(session *Session, self VMRef, staticMin int, staticMax int, dynamicMin int, dynamicMax int) (err error) {
	method := "VM.set_memory_limits"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	staticMinArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "static_min"), staticMin)
	if err != nil {
		return
	}
	staticMaxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "static_max"), staticMax)
	if err != nil {
		return
	}
	dynamicMinArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "dynamic_min"), dynamicMin)
	if err != nil {
		return
	}
	dynamicMaxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "dynamic_max"), dynamicMax)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, staticMinArg, staticMaxArg, dynamicMinArg, dynamicMaxArg)
	return
}

// AsyncSetMemoryLimits: Set the memory limits of this VM.
// Version: midnight-ride
func (vm) AsyncSetMemoryLimits(session *Session, self VMRef, staticMin int, staticMax int, dynamicMin int, dynamicMax int) (retval TaskRef, err error) {
	method := "Async.VM.set_memory_limits"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	staticMinArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "static_min"), staticMin)
	if err != nil {
		return
	}
	staticMaxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "static_max"), staticMax)
	if err != nil {
		return
	}
	dynamicMinArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "dynamic_min"), dynamicMin)
	if err != nil {
		return
	}
	dynamicMaxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "dynamic_max"), dynamicMax)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, staticMinArg, staticMaxArg, dynamicMinArg, dynamicMaxArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMemoryLimits6: Set the memory limits of this VM.
// Version: midnight-ride
func (vm) SetMemoryLimits6(session *Session, self VMRef, staticMin int, staticMax int, dynamicMin int, dynamicMax int) (err error) {
	method := "VM.set_memory_limits"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	staticMinArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "static_min"), staticMin)
	if err != nil {
		return
	}
	staticMaxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "static_max"), staticMax)
	if err != nil {
		return
	}
	dynamicMinArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "dynamic_min"), dynamicMin)
	if err != nil {
		return
	}
	dynamicMaxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "dynamic_max"), dynamicMax)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, staticMinArg, staticMaxArg, dynamicMinArg, dynamicMaxArg)
	return
}

// AsyncSetMemoryLimits6: Set the memory limits of this VM.
// Version: midnight-ride
func (vm) AsyncSetMemoryLimits6(session *Session, self VMRef, staticMin int, staticMax int, dynamicMin int, dynamicMax int) (retval TaskRef, err error) {
	method := "Async.VM.set_memory_limits"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	staticMinArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "static_min"), staticMin)
	if err != nil {
		return
	}
	staticMaxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "static_max"), staticMax)
	if err != nil {
		return
	}
	dynamicMinArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "dynamic_min"), dynamicMin)
	if err != nil {
		return
	}
	dynamicMaxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "dynamic_max"), dynamicMax)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, staticMinArg, staticMaxArg, dynamicMinArg, dynamicMaxArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMemoryStaticRange: Set the static (ie boot-time) range of virtual memory that the VM is allowed to use.
// Version: midnight-ride
func (vm) SetMemoryStaticRange(session *Session, self VMRef, min int, max int) (err error) {
	method := "VM.set_memory_static_range"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	minArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "min"), min)
	if err != nil {
		return
	}
	maxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "max"), max)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, minArg, maxArg)
	return
}

// AsyncSetMemoryStaticRange: Set the static (ie boot-time) range of virtual memory that the VM is allowed to use.
// Version: midnight-ride
func (vm) AsyncSetMemoryStaticRange(session *Session, self VMRef, min int, max int) (retval TaskRef, err error) {
	method := "Async.VM.set_memory_static_range"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	minArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "min"), min)
	if err != nil {
		return
	}
	maxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "max"), max)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, minArg, maxArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMemoryStaticRange4: Set the static (ie boot-time) range of virtual memory that the VM is allowed to use.
// Version: midnight-ride
func (vm) SetMemoryStaticRange4(session *Session, self VMRef, min int, max int) (err error) {
	method := "VM.set_memory_static_range"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	minArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "min"), min)
	if err != nil {
		return
	}
	maxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "max"), max)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, minArg, maxArg)
	return
}

// AsyncSetMemoryStaticRange4: Set the static (ie boot-time) range of virtual memory that the VM is allowed to use.
// Version: midnight-ride
func (vm) AsyncSetMemoryStaticRange4(session *Session, self VMRef, min int, max int) (retval TaskRef, err error) {
	method := "Async.VM.set_memory_static_range"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	minArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "min"), min)
	if err != nil {
		return
	}
	maxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "max"), max)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, minArg, maxArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMemoryStaticMin: Set the value of the memory_static_min field
// Version: midnight-ride
func (vm) SetMemoryStaticMin(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_memory_static_min"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetMemoryStaticMin3: Set the value of the memory_static_min field
// Version: midnight-ride
func (vm) SetMemoryStaticMin3(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_memory_static_min"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetMemoryStaticMax: Set the value of the memory_static_max field
// Version: orlando
//
// Errors:
// HA_OPERATION_WOULD_BREAK_FAILOVER_PLAN - This operation cannot be performed because it would invalidate VM failover planning such that the system would be unable to guarantee to restart protected VMs after a Host failure.
func (vm) SetMemoryStaticMax(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_memory_static_max"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetMemoryStaticMax3: Set the value of the memory_static_max field
// Version: orlando
//
// Errors:
// HA_OPERATION_WOULD_BREAK_FAILOVER_PLAN - This operation cannot be performed because it would invalidate VM failover planning such that the system would be unable to guarantee to restart protected VMs after a Host failure.
func (vm) SetMemoryStaticMax3(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_memory_static_max"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetMemoryDynamicRange: Set the minimum and maximum amounts of physical memory the VM is allowed to use.
// Version: midnight-ride
func (vm) SetMemoryDynamicRange(session *Session, self VMRef, min int, max int) (err error) {
	method := "VM.set_memory_dynamic_range"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	minArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "min"), min)
	if err != nil {
		return
	}
	maxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "max"), max)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, minArg, maxArg)
	return
}

// AsyncSetMemoryDynamicRange: Set the minimum and maximum amounts of physical memory the VM is allowed to use.
// Version: midnight-ride
func (vm) AsyncSetMemoryDynamicRange(session *Session, self VMRef, min int, max int) (retval TaskRef, err error) {
	method := "Async.VM.set_memory_dynamic_range"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	minArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "min"), min)
	if err != nil {
		return
	}
	maxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "max"), max)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, minArg, maxArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMemoryDynamicRange4: Set the minimum and maximum amounts of physical memory the VM is allowed to use.
// Version: midnight-ride
func (vm) SetMemoryDynamicRange4(session *Session, self VMRef, min int, max int) (err error) {
	method := "VM.set_memory_dynamic_range"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	minArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "min"), min)
	if err != nil {
		return
	}
	maxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "max"), max)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, minArg, maxArg)
	return
}

// AsyncSetMemoryDynamicRange4: Set the minimum and maximum amounts of physical memory the VM is allowed to use.
// Version: midnight-ride
func (vm) AsyncSetMemoryDynamicRange4(session *Session, self VMRef, min int, max int) (retval TaskRef, err error) {
	method := "Async.VM.set_memory_dynamic_range"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	minArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "min"), min)
	if err != nil {
		return
	}
	maxArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "max"), max)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, minArg, maxArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetMemoryDynamicMin: Set the value of the memory_dynamic_min field
// Version: midnight-ride
func (vm) SetMemoryDynamicMin(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_memory_dynamic_min"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetMemoryDynamicMin3: Set the value of the memory_dynamic_min field
// Version: midnight-ride
func (vm) SetMemoryDynamicMin3(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_memory_dynamic_min"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetMemoryDynamicMax: Set the value of the memory_dynamic_max field
// Version: midnight-ride
func (vm) SetMemoryDynamicMax(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_memory_dynamic_max"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetMemoryDynamicMax3: Set the value of the memory_dynamic_max field
// Version: midnight-ride
func (vm) SetMemoryDynamicMax3(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_memory_dynamic_max"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// ComputeMemoryOverhead: Computes the virtualization memory overhead of a VM.
// Version: midnight-ride
func (vm) ComputeMemoryOverhead(session *Session, vm VMRef) (retval int, err error) {
	method := "VM.compute_memory_overhead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// AsyncComputeMemoryOverhead: Computes the virtualization memory overhead of a VM.
// Version: midnight-ride
func (vm) AsyncComputeMemoryOverhead(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.compute_memory_overhead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ComputeMemoryOverhead2: Computes the virtualization memory overhead of a VM.
// Version: midnight-ride
func (vm) ComputeMemoryOverhead2(session *Session, vm VMRef) (retval int, err error) {
	method := "VM.compute_memory_overhead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeInt(method+" -> ", result)
	return
}

// AsyncComputeMemoryOverhead2: Computes the virtualization memory overhead of a VM.
// Version: midnight-ride
func (vm) AsyncComputeMemoryOverhead2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.compute_memory_overhead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetHaAlwaysRun: Set the value of the ha_always_run
// Version: orlando
func (vm) SetHaAlwaysRun(session *Session, self VMRef, value bool) (err error) {
	method := "VM.set_ha_always_run"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHaAlwaysRun3: Set the value of the ha_always_run
// Version: orlando
func (vm) SetHaAlwaysRun3(session *Session, self VMRef, value bool) (err error) {
	method := "VM.set_ha_always_run"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHaRestartPriority: Set the value of the ha_restart_priority field
// Version: orlando
func (vm) SetHaRestartPriority(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_ha_restart_priority"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHaRestartPriority3: Set the value of the ha_restart_priority field
// Version: orlando
func (vm) SetHaRestartPriority3(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_ha_restart_priority"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromNVRAM:
// Version: naples
func (vm) RemoveFromNVRAM(session *Session, self VMRef, key string) (err error) {
	method := "VM.remove_from_NVRAM"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromNVRAM3:
// Version: naples
func (vm) RemoveFromNVRAM3(session *Session, self VMRef, key string) (err error) {
	method := "VM.remove_from_NVRAM"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToNVRAM:
// Version: naples
func (vm) AddToNVRAM(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_NVRAM"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToNVRAM4:
// Version: naples
func (vm) AddToNVRAM4(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_NVRAM"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNVRAM:
// Version: naples
func (vm) SetNVRAM(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_NVRAM"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNVRAM3:
// Version: naples
func (vm) SetNVRAM3(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_NVRAM"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToVCPUsParamsLive: Add the given key-value pair to VM.VCPUs_params, and apply that value on the running VM
// Version: rio
func (vm) AddToVCPUsParamsLive(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_VCPUs_params_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncAddToVCPUsParamsLive: Add the given key-value pair to VM.VCPUs_params, and apply that value on the running VM
// Version: rio
func (vm) AsyncAddToVCPUsParamsLive(session *Session, self VMRef, key string, value string) (retval TaskRef, err error) {
	method := "Async.VM.add_to_VCPUs_params_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// AddToVCPUsParamsLive4: Add the given key-value pair to VM.VCPUs_params, and apply that value on the running VM
// Version: rio
func (vm) AddToVCPUsParamsLive4(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_VCPUs_params_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AsyncAddToVCPUsParamsLive4: Add the given key-value pair to VM.VCPUs_params, and apply that value on the running VM
// Version: rio
func (vm) AsyncAddToVCPUsParamsLive4(session *Session, self VMRef, key string, value string) (retval TaskRef, err error) {
	method := "Async.VM.add_to_VCPUs_params_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, keyArg, valueArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetVCPUsNumberLive: Set the number of VCPUs for a running VM
// Version: rio
//
// Errors:
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) SetVCPUsNumberLive(session *Session, self VMRef, nvcpu int) (err error) {
	method := "VM.set_VCPUs_number_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nvcpuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "nvcpu"), nvcpu)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, nvcpuArg)
	return
}

// AsyncSetVCPUsNumberLive: Set the number of VCPUs for a running VM
// Version: rio
//
// Errors:
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncSetVCPUsNumberLive(session *Session, self VMRef, nvcpu int) (retval TaskRef, err error) {
	method := "Async.VM.set_VCPUs_number_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nvcpuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "nvcpu"), nvcpu)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, nvcpuArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetVCPUsNumberLive3: Set the number of VCPUs for a running VM
// Version: rio
//
// Errors:
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) SetVCPUsNumberLive3(session *Session, self VMRef, nvcpu int) (err error) {
	method := "VM.set_VCPUs_number_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nvcpuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "nvcpu"), nvcpu)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, nvcpuArg)
	return
}

// AsyncSetVCPUsNumberLive3: Set the number of VCPUs for a running VM
// Version: rio
//
// Errors:
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncSetVCPUsNumberLive3(session *Session, self VMRef, nvcpu int) (retval TaskRef, err error) {
	method := "Async.VM.set_VCPUs_number_live"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	nvcpuArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "nvcpu"), nvcpu)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg, nvcpuArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// PoolMigrate: Migrate a VM to another Host.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) PoolMigrate(session *Session, vm VMRef, host HostRef, options map[string]string) (err error) {
	method := "VM.pool_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, hostArg, optionsArg)
	return
}

// AsyncPoolMigrate: Migrate a VM to another Host.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) AsyncPoolMigrate(session *Session, vm VMRef, host HostRef, options map[string]string) (retval TaskRef, err error) {
	method := "Async.VM.pool_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, hostArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// PoolMigrate4: Migrate a VM to another Host.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) PoolMigrate4(session *Session, vm VMRef, host HostRef, options map[string]string) (err error) {
	method := "VM.pool_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, hostArg, optionsArg)
	return
}

// AsyncPoolMigrate4: Migrate a VM to another Host.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
func (vm) AsyncPoolMigrate4(session *Session, vm VMRef, host HostRef, options map[string]string) (retval TaskRef, err error) {
	method := "Async.VM.pool_migrate"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	optionsArg, err := serializeStringToStringMap(fmt.Sprintf("%s(%s)", method, "options"), options)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, hostArg, optionsArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ResumeOn: Awaken the specified VM and resume it on a particular Host.  This can only be called when the specified VM is in the Suspended state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) ResumeOn(session *Session, vm VMRef, host HostRef, startPaused bool, force bool) (err error) {
	method := "VM.resume_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, hostArg, startPausedArg, forceArg)
	return
}

// AsyncResumeOn: Awaken the specified VM and resume it on a particular Host.  This can only be called when the specified VM is in the Suspended state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncResumeOn(session *Session, vm VMRef, host HostRef, startPaused bool, force bool) (retval TaskRef, err error) {
	method := "Async.VM.resume_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, hostArg, startPausedArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// ResumeOn5: Awaken the specified VM and resume it on a particular Host.  This can only be called when the specified VM is in the Suspended state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) ResumeOn5(session *Session, vm VMRef, host HostRef, startPaused bool, force bool) (err error) {
	method := "VM.resume_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, hostArg, startPausedArg, forceArg)
	return
}

// AsyncResumeOn5: Awaken the specified VM and resume it on a particular Host.  This can only be called when the specified VM is in the Suspended state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncResumeOn5(session *Session, vm VMRef, host HostRef, startPaused bool, force bool) (retval TaskRef, err error) {
	method := "Async.VM.resume_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, hostArg, startPausedArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Resume: Awaken the specified VM and resume it.  This can only be called when the specified VM is in the Suspended state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) Resume(session *Session, vm VMRef, startPaused bool, force bool) (err error) {
	method := "VM.resume"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, startPausedArg, forceArg)
	return
}

// AsyncResume: Awaken the specified VM and resume it.  This can only be called when the specified VM is in the Suspended state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncResume(session *Session, vm VMRef, startPaused bool, force bool) (retval TaskRef, err error) {
	method := "Async.VM.resume"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, startPausedArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Resume4: Awaken the specified VM and resume it.  This can only be called when the specified VM is in the Suspended state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) Resume4(session *Session, vm VMRef, startPaused bool, force bool) (err error) {
	method := "VM.resume"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, startPausedArg, forceArg)
	return
}

// AsyncResume4: Awaken the specified VM and resume it.  This can only be called when the specified VM is in the Suspended state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncResume4(session *Session, vm VMRef, startPaused bool, force bool) (retval TaskRef, err error) {
	method := "Async.VM.resume"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, startPausedArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Suspend: Suspend the specified VM to disk.  This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) Suspend(session *Session, vm VMRef) (err error) {
	method := "VM.suspend"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncSuspend: Suspend the specified VM to disk.  This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncSuspend(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.suspend"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Suspend2: Suspend the specified VM to disk.  This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) Suspend2(session *Session, vm VMRef) (err error) {
	method := "VM.suspend"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncSuspend2: Suspend the specified VM to disk.  This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncSuspend2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.suspend"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// HardReboot: Stop executing the specified VM without attempting a clean shutdown and immediately restart the VM.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) HardReboot(session *Session, vm VMRef) (err error) {
	method := "VM.hard_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncHardReboot: Stop executing the specified VM without attempting a clean shutdown and immediately restart the VM.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncHardReboot(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.hard_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// HardReboot2: Stop executing the specified VM without attempting a clean shutdown and immediately restart the VM.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) HardReboot2(session *Session, vm VMRef) (err error) {
	method := "VM.hard_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncHardReboot2: Stop executing the specified VM without attempting a clean shutdown and immediately restart the VM.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncHardReboot2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.hard_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// PowerStateReset: Reset the power-state of the VM to halted in the database only. (Used to recover from slave failures in pooling scenarios by resetting the power-states of VMs running on dead slaves to halted.) This is a potentially dangerous operation; use with care.
// Version: rio
func (vm) PowerStateReset(session *Session, vm VMRef) (err error) {
	method := "VM.power_state_reset"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncPowerStateReset: Reset the power-state of the VM to halted in the database only. (Used to recover from slave failures in pooling scenarios by resetting the power-states of VMs running on dead slaves to halted.) This is a potentially dangerous operation; use with care.
// Version: rio
func (vm) AsyncPowerStateReset(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.power_state_reset"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// PowerStateReset2: Reset the power-state of the VM to halted in the database only. (Used to recover from slave failures in pooling scenarios by resetting the power-states of VMs running on dead slaves to halted.) This is a potentially dangerous operation; use with care.
// Version: rio
func (vm) PowerStateReset2(session *Session, vm VMRef) (err error) {
	method := "VM.power_state_reset"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncPowerStateReset2: Reset the power-state of the VM to halted in the database only. (Used to recover from slave failures in pooling scenarios by resetting the power-states of VMs running on dead slaves to halted.) This is a potentially dangerous operation; use with care.
// Version: rio
func (vm) AsyncPowerStateReset2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.power_state_reset"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// HardShutdown: Stop executing the specified VM without attempting a clean shutdown.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) HardShutdown(session *Session, vm VMRef) (err error) {
	method := "VM.hard_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncHardShutdown: Stop executing the specified VM without attempting a clean shutdown.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncHardShutdown(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.hard_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// HardShutdown2: Stop executing the specified VM without attempting a clean shutdown.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) HardShutdown2(session *Session, vm VMRef) (err error) {
	method := "VM.hard_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncHardShutdown2: Stop executing the specified VM without attempting a clean shutdown.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncHardShutdown2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.hard_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CleanReboot: Attempt to cleanly shutdown the specified VM (Note: this may not be supported, for example, if a guest agent is not installed). This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) CleanReboot(session *Session, vm VMRef) (err error) {
	method := "VM.clean_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncCleanReboot: Attempt to cleanly shutdown the specified VM (Note: this may not be supported, for example, if a guest agent is not installed). This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncCleanReboot(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.clean_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CleanReboot2: Attempt to cleanly shutdown the specified VM (Note: this may not be supported, for example, if a guest agent is not installed). This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) CleanReboot2(session *Session, vm VMRef) (err error) {
	method := "VM.clean_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncCleanReboot2: Attempt to cleanly shutdown the specified VM (Note: this may not be supported, for example, if a guest agent is not installed). This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncCleanReboot2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.clean_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Shutdown: Attempts to first clean shutdown a VM and if it should fail then perform a hard shutdown on it.
// Version: clearwater
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) Shutdown(session *Session, vm VMRef) (err error) {
	method := "VM.shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncShutdown: Attempts to first clean shutdown a VM and if it should fail then perform a hard shutdown on it.
// Version: clearwater
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncShutdown(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Shutdown2: Attempts to first clean shutdown a VM and if it should fail then perform a hard shutdown on it.
// Version: clearwater
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) Shutdown2(session *Session, vm VMRef) (err error) {
	method := "VM.shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncShutdown2: Attempts to first clean shutdown a VM and if it should fail then perform a hard shutdown on it.
// Version: clearwater
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncShutdown2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CleanShutdown: Attempt to cleanly shutdown the specified VM. (Note: this may not be supported, for example, if a guest agent is not installed). This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) CleanShutdown(session *Session, vm VMRef) (err error) {
	method := "VM.clean_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncCleanShutdown: Attempt to cleanly shutdown the specified VM. (Note: this may not be supported, for example, if a guest agent is not installed). This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncCleanShutdown(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.clean_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// CleanShutdown2: Attempt to cleanly shutdown the specified VM. (Note: this may not be supported, for example, if a guest agent is not installed). This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) CleanShutdown2(session *Session, vm VMRef) (err error) {
	method := "VM.clean_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncCleanShutdown2: Attempt to cleanly shutdown the specified VM. (Note: this may not be supported, for example, if a guest agent is not installed). This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncCleanShutdown2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.clean_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Unpause: Resume the specified VM. This can only be called when the specified VM is in the Paused state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) Unpause(session *Session, vm VMRef) (err error) {
	method := "VM.unpause"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncUnpause: Resume the specified VM. This can only be called when the specified VM is in the Paused state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncUnpause(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.unpause"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Unpause2: Resume the specified VM. This can only be called when the specified VM is in the Paused state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) Unpause2(session *Session, vm VMRef) (err error) {
	method := "VM.unpause"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncUnpause2: Resume the specified VM. This can only be called when the specified VM is in the Paused state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncUnpause2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.unpause"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Pause: Pause the specified VM. This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) Pause(session *Session, vm VMRef) (err error) {
	method := "VM.pause"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncPause: Pause the specified VM. This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncPause(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.pause"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Pause2: Pause the specified VM. This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) Pause2(session *Session, vm VMRef) (err error) {
	method := "VM.pause"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncPause2: Pause the specified VM. This can only be called when the specified VM is in the Running state.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
func (vm) AsyncPause2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.pause"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// StartOn: Start the specified VM on a particular host.  This function can only be called with the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// BOOTLOADER_FAILED - The bootloader returned an error
// UNKNOWN_BOOTLOADER - The requested bootloader is unknown
func (vm) StartOn(session *Session, vm VMRef, host HostRef, startPaused bool, force bool) (err error) {
	method := "VM.start_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, hostArg, startPausedArg, forceArg)
	return
}

// AsyncStartOn: Start the specified VM on a particular host.  This function can only be called with the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// BOOTLOADER_FAILED - The bootloader returned an error
// UNKNOWN_BOOTLOADER - The requested bootloader is unknown
func (vm) AsyncStartOn(session *Session, vm VMRef, host HostRef, startPaused bool, force bool) (retval TaskRef, err error) {
	method := "Async.VM.start_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, hostArg, startPausedArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// StartOn5: Start the specified VM on a particular host.  This function can only be called with the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// BOOTLOADER_FAILED - The bootloader returned an error
// UNKNOWN_BOOTLOADER - The requested bootloader is unknown
func (vm) StartOn5(session *Session, vm VMRef, host HostRef, startPaused bool, force bool) (err error) {
	method := "VM.start_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, hostArg, startPausedArg, forceArg)
	return
}

// AsyncStartOn5: Start the specified VM on a particular host.  This function can only be called with the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// BOOTLOADER_FAILED - The bootloader returned an error
// UNKNOWN_BOOTLOADER - The requested bootloader is unknown
func (vm) AsyncStartOn5(session *Session, vm VMRef, host HostRef, startPaused bool, force bool) (retval TaskRef, err error) {
	method := "Async.VM.start_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	hostArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "host"), host)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, hostArg, startPausedArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Start: Start the specified VM.  This function can only be called with the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// VM_HVM_REQUIRED - HVM is required for this operation
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// BOOTLOADER_FAILED - The bootloader returned an error
// UNKNOWN_BOOTLOADER - The requested bootloader is unknown
// NO_HOSTS_AVAILABLE - There were no servers available to complete the specified operation.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) Start(session *Session, vm VMRef, startPaused bool, force bool) (err error) {
	method := "VM.start"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, startPausedArg, forceArg)
	return
}

// AsyncStart: Start the specified VM.  This function can only be called with the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// VM_HVM_REQUIRED - HVM is required for this operation
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// BOOTLOADER_FAILED - The bootloader returned an error
// UNKNOWN_BOOTLOADER - The requested bootloader is unknown
// NO_HOSTS_AVAILABLE - There were no servers available to complete the specified operation.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncStart(session *Session, vm VMRef, startPaused bool, force bool) (retval TaskRef, err error) {
	method := "Async.VM.start"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, startPausedArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Start4: Start the specified VM.  This function can only be called with the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// VM_HVM_REQUIRED - HVM is required for this operation
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// BOOTLOADER_FAILED - The bootloader returned an error
// UNKNOWN_BOOTLOADER - The requested bootloader is unknown
// NO_HOSTS_AVAILABLE - There were no servers available to complete the specified operation.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) Start4(session *Session, vm VMRef, startPaused bool, force bool) (err error) {
	method := "VM.start"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg, startPausedArg, forceArg)
	return
}

// AsyncStart4: Start the specified VM.  This function can only be called with the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// VM_HVM_REQUIRED - HVM is required for this operation
// VM_IS_TEMPLATE - The operation attempted is not valid for a template VM
// OTHER_OPERATION_IN_PROGRESS - Another operation involving the object is currently in progress
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// BOOTLOADER_FAILED - The bootloader returned an error
// UNKNOWN_BOOTLOADER - The requested bootloader is unknown
// NO_HOSTS_AVAILABLE - There were no servers available to complete the specified operation.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncStart4(session *Session, vm VMRef, startPaused bool, force bool) (retval TaskRef, err error) {
	method := "Async.VM.start"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	startPausedArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "start_paused"), startPaused)
	if err != nil {
		return
	}
	forceArg, err := serializeBool(fmt.Sprintf("%s(%s)", method, "force"), force)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, startPausedArg, forceArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Provision: Inspects the disk configuration contained within the VM&apos;s other_config, creates VDIs and VBDs and then executes any applicable post-install script.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) Provision(session *Session, vm VMRef) (err error) {
	method := "VM.provision"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncProvision: Inspects the disk configuration contained within the VM&apos;s other_config, creates VDIs and VBDs and then executes any applicable post-install script.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncProvision(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.provision"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Provision2: Inspects the disk configuration contained within the VM&apos;s other_config, creates VDIs and VBDs and then executes any applicable post-install script.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) Provision2(session *Session, vm VMRef) (err error) {
	method := "VM.provision"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, vmArg)
	return
}

// AsyncProvision2: Inspects the disk configuration contained within the VM&apos;s other_config, creates VDIs and VBDs and then executes any applicable post-install script.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncProvision2(session *Session, vm VMRef) (retval TaskRef, err error) {
	method := "Async.VM.provision"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Checkpoint: Checkpoints the specified VM, making a new VM. Checkpoint automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write) and saves the memory image as well.
// Version: midnight-ride
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_CHECKPOINT_SUSPEND_FAILED - An error occured while saving the memory image of the specified virtual machine
// VM_CHECKPOINT_RESUME_FAILED - An error occured while restoring the memory image of the specified virtual machine
func (vm) Checkpoint(session *Session, vm VMRef, newName string) (retval VMRef, err error) {
	method := "VM.checkpoint"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncCheckpoint: Checkpoints the specified VM, making a new VM. Checkpoint automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write) and saves the memory image as well.
// Version: midnight-ride
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_CHECKPOINT_SUSPEND_FAILED - An error occured while saving the memory image of the specified virtual machine
// VM_CHECKPOINT_RESUME_FAILED - An error occured while restoring the memory image of the specified virtual machine
func (vm) AsyncCheckpoint(session *Session, vm VMRef, newName string) (retval TaskRef, err error) {
	method := "Async.VM.checkpoint"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Checkpoint3: Checkpoints the specified VM, making a new VM. Checkpoint automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write) and saves the memory image as well.
// Version: midnight-ride
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_CHECKPOINT_SUSPEND_FAILED - An error occured while saving the memory image of the specified virtual machine
// VM_CHECKPOINT_RESUME_FAILED - An error occured while restoring the memory image of the specified virtual machine
func (vm) Checkpoint3(session *Session, vm VMRef, newName string) (retval VMRef, err error) {
	method := "VM.checkpoint"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncCheckpoint3: Checkpoints the specified VM, making a new VM. Checkpoint automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write) and saves the memory image as well.
// Version: midnight-ride
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_CHECKPOINT_SUSPEND_FAILED - An error occured while saving the memory image of the specified virtual machine
// VM_CHECKPOINT_RESUME_FAILED - An error occured while restoring the memory image of the specified virtual machine
func (vm) AsyncCheckpoint3(session *Session, vm VMRef, newName string) (retval TaskRef, err error) {
	method := "Async.VM.checkpoint"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Revert: Reverts the specified VM to a previous state.
// Version: midnight-ride
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// VM_REVERT_FAILED - An error occured while reverting the specified virtual machine to the specified snapshot
func (vm) Revert(session *Session, snapshot VMRef) (err error) {
	method := "VM.revert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	snapshotArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "snapshot"), snapshot)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, snapshotArg)
	return
}

// AsyncRevert: Reverts the specified VM to a previous state.
// Version: midnight-ride
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// VM_REVERT_FAILED - An error occured while reverting the specified virtual machine to the specified snapshot
func (vm) AsyncRevert(session *Session, snapshot VMRef) (retval TaskRef, err error) {
	method := "Async.VM.revert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	snapshotArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "snapshot"), snapshot)
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

// Revert2: Reverts the specified VM to a previous state.
// Version: midnight-ride
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// VM_REVERT_FAILED - An error occured while reverting the specified virtual machine to the specified snapshot
func (vm) Revert2(session *Session, snapshot VMRef) (err error) {
	method := "VM.revert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	snapshotArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "snapshot"), snapshot)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, snapshotArg)
	return
}

// AsyncRevert2: Reverts the specified VM to a previous state.
// Version: midnight-ride
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// VM_REVERT_FAILED - An error occured while reverting the specified virtual machine to the specified snapshot
func (vm) AsyncRevert2(session *Session, snapshot VMRef) (retval TaskRef, err error) {
	method := "Async.VM.revert"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	snapshotArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "snapshot"), snapshot)
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

// Copy: Copied the specified VM, making a new VM. Unlike clone, copy does not exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored. Instead, copy guarantees that the disk images of the newly created VM will be &apos;full disks&apos; - i.e. not part of a CoW chain.  This function can only be called when the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) Copy(session *Session, vm VMRef, newName string, sr SRRef) (retval VMRef, err error) {
	method := "VM.copy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg, srArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncCopy: Copied the specified VM, making a new VM. Unlike clone, copy does not exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored. Instead, copy guarantees that the disk images of the newly created VM will be &apos;full disks&apos; - i.e. not part of a CoW chain.  This function can only be called when the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncCopy(session *Session, vm VMRef, newName string, sr SRRef) (retval TaskRef, err error) {
	method := "Async.VM.copy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg, srArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Copy4: Copied the specified VM, making a new VM. Unlike clone, copy does not exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored. Instead, copy guarantees that the disk images of the newly created VM will be &apos;full disks&apos; - i.e. not part of a CoW chain.  This function can only be called when the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) Copy4(session *Session, vm VMRef, newName string, sr SRRef) (retval VMRef, err error) {
	method := "VM.copy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg, srArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncCopy4: Copied the specified VM, making a new VM. Unlike clone, copy does not exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored. Instead, copy guarantees that the disk images of the newly created VM will be &apos;full disks&apos; - i.e. not part of a CoW chain.  This function can only be called when the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncCopy4(session *Session, vm VMRef, newName string, sr SRRef) (retval TaskRef, err error) {
	method := "Async.VM.copy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	srArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "sr"), sr)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg, srArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Clone: Clones the specified VM, making a new VM. Clone automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write). This function can only be called when the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) Clone(session *Session, vm VMRef, newName string) (retval VMRef, err error) {
	method := "VM.clone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncClone: Clones the specified VM, making a new VM. Clone automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write). This function can only be called when the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncClone(session *Session, vm VMRef, newName string) (retval TaskRef, err error) {
	method := "Async.VM.clone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Clone3: Clones the specified VM, making a new VM. Clone automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write). This function can only be called when the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) Clone3(session *Session, vm VMRef, newName string) (retval VMRef, err error) {
	method := "VM.clone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncClone3: Clones the specified VM, making a new VM. Clone automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write). This function can only be called when the VM is in the Halted State.
// Version: rio
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// LICENCE_RESTRICTION - This operation is not allowed because your license lacks a needed feature. Please contact your support representative.
func (vm) AsyncClone3(session *Session, vm VMRef, newName string) (retval TaskRef, err error) {
	method := "Async.VM.clone"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SnapshotWithQuiesce: Snapshots the specified VM with quiesce, making a new VM. Snapshot automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write).
// Version: orlando
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_SNAPSHOT_WITH_QUIESCE_FAILED - The quiesced-snapshot operation failed for an unexpected reason
// VM_SNAPSHOT_WITH_QUIESCE_TIMEOUT - The VSS plug-in has timed out
// VM_SNAPSHOT_WITH_QUIESCE_PLUGIN_DEOS_NOT_RESPOND - The VSS plug-in cannot be contacted
// VM_SNAPSHOT_WITH_QUIESCE_NOT_SUPPORTED - The VSS plug-in is not installed on this virtual machine
func (vm) SnapshotWithQuiesce(session *Session, vm VMRef, newName string) (retval VMRef, err error) {
	method := "VM.snapshot_with_quiesce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncSnapshotWithQuiesce: Snapshots the specified VM with quiesce, making a new VM. Snapshot automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write).
// Version: orlando
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_SNAPSHOT_WITH_QUIESCE_FAILED - The quiesced-snapshot operation failed for an unexpected reason
// VM_SNAPSHOT_WITH_QUIESCE_TIMEOUT - The VSS plug-in has timed out
// VM_SNAPSHOT_WITH_QUIESCE_PLUGIN_DEOS_NOT_RESPOND - The VSS plug-in cannot be contacted
// VM_SNAPSHOT_WITH_QUIESCE_NOT_SUPPORTED - The VSS plug-in is not installed on this virtual machine
func (vm) AsyncSnapshotWithQuiesce(session *Session, vm VMRef, newName string) (retval TaskRef, err error) {
	method := "Async.VM.snapshot_with_quiesce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SnapshotWithQuiesce3: Snapshots the specified VM with quiesce, making a new VM. Snapshot automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write).
// Version: orlando
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_SNAPSHOT_WITH_QUIESCE_FAILED - The quiesced-snapshot operation failed for an unexpected reason
// VM_SNAPSHOT_WITH_QUIESCE_TIMEOUT - The VSS plug-in has timed out
// VM_SNAPSHOT_WITH_QUIESCE_PLUGIN_DEOS_NOT_RESPOND - The VSS plug-in cannot be contacted
// VM_SNAPSHOT_WITH_QUIESCE_NOT_SUPPORTED - The VSS plug-in is not installed on this virtual machine
func (vm) SnapshotWithQuiesce3(session *Session, vm VMRef, newName string) (retval VMRef, err error) {
	method := "VM.snapshot_with_quiesce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncSnapshotWithQuiesce3: Snapshots the specified VM with quiesce, making a new VM. Snapshot automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write).
// Version: orlando
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
// VM_SNAPSHOT_WITH_QUIESCE_FAILED - The quiesced-snapshot operation failed for an unexpected reason
// VM_SNAPSHOT_WITH_QUIESCE_TIMEOUT - The VSS plug-in has timed out
// VM_SNAPSHOT_WITH_QUIESCE_PLUGIN_DEOS_NOT_RESPOND - The VSS plug-in cannot be contacted
// VM_SNAPSHOT_WITH_QUIESCE_NOT_SUPPORTED - The VSS plug-in is not installed on this virtual machine
func (vm) AsyncSnapshotWithQuiesce3(session *Session, vm VMRef, newName string) (retval TaskRef, err error) {
	method := "Async.VM.snapshot_with_quiesce"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Snapshot: Snapshots the specified VM, making a new VM. Snapshot automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write).
// Version: 21.4.0
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
func (vm) Snapshot(session *Session, vm VMRef, newName string, ignoreVdis []VDIRef) (retval VMRef, err error) {
	method := "VM.snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	ignoreVdisArg, err := serializeVDIRefSet(fmt.Sprintf("%s(%s)", method, "ignore_vdis"), ignoreVdis)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg, ignoreVdisArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncSnapshot: Snapshots the specified VM, making a new VM. Snapshot automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write).
// Version: 21.4.0
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
func (vm) AsyncSnapshot(session *Session, vm VMRef, newName string, ignoreVdis []VDIRef) (retval TaskRef, err error) {
	method := "Async.VM.snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	ignoreVdisArg, err := serializeVDIRefSet(fmt.Sprintf("%s(%s)", method, "ignore_vdis"), ignoreVdis)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg, ignoreVdisArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Snapshot4: Snapshots the specified VM, making a new VM. Snapshot automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write).
// Version: 21.4.0
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
func (vm) Snapshot4(session *Session, vm VMRef, newName string, ignoreVdis []VDIRef) (retval VMRef, err error) {
	method := "VM.snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	ignoreVdisArg, err := serializeVDIRefSet(fmt.Sprintf("%s(%s)", method, "ignore_vdis"), ignoreVdis)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg, ignoreVdisArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncSnapshot4: Snapshots the specified VM, making a new VM. Snapshot automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write).
// Version: 21.4.0
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
func (vm) AsyncSnapshot4(session *Session, vm VMRef, newName string, ignoreVdis []VDIRef) (retval TaskRef, err error) {
	method := "Async.VM.snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	ignoreVdisArg, err := serializeVDIRefSet(fmt.Sprintf("%s(%s)", method, "ignore_vdis"), ignoreVdis)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg, ignoreVdisArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// Snapshot3: Snapshots the specified VM, making a new VM. Snapshot automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write).
// Version: orlando
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
func (vm) Snapshot3(session *Session, vm VMRef, newName string) (retval VMRef, err error) {
	method := "VM.snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncSnapshot3: Snapshots the specified VM, making a new VM. Snapshot automatically exploits the capabilities of the underlying storage repository in which the VM&apos;s disk images are stored (for example, Copy on Write).
// Version: orlando
//
// Errors:
// VM_BAD_POWER_STATE - You attempted an operation on a VM that was not in an appropriate power state at the time; for example, you attempted to start a VM that was already running. The parameters returned are the VM&apos;s handle, and the expected and actual VM state at the time of the call.
// SR_FULL - The SR is full. Requested new size exceeds the maximum size
// OPERATION_NOT_ALLOWED - You attempted an operation that was not allowed.
func (vm) AsyncSnapshot3(session *Session, vm VMRef, newName string) (retval TaskRef, err error) {
	method := "Async.VM.snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "vm"), vm)
	if err != nil {
		return
	}
	newNameArg, err := serializeString(fmt.Sprintf("%s(%s)", method, "new_name"), newName)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmArg, newNameArg)
	if err != nil {
		return
	}
	retval, err = deserializeTaskRef(method+" -> ", result)
	return
}

// SetHardwarePlatformVersion: Set the hardware_platform_version field of the given VM.
// Version: cream
func (vm) SetHardwarePlatformVersion(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_hardware_platform_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetHardwarePlatformVersion3: Set the hardware_platform_version field of the given VM.
// Version: cream
func (vm) SetHardwarePlatformVersion3(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_hardware_platform_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetHardwarePlatformVersion2: Set the hardware_platform_version field of the given VM.
// Version: rio
func (vm) SetHardwarePlatformVersion2(session *Session, self VMRef) (err error) {
	method := "VM.set_hardware_platform_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetSuspendSR: Set the suspend_SR field of the given VM.
// Version: boston
func (vm) SetSuspendSR(session *Session, self VMRef, value SRRef) (err error) {
	method := "VM.set_suspend_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetSuspendSR3: Set the suspend_SR field of the given VM.
// Version: boston
func (vm) SetSuspendSR3(session *Session, self VMRef, value SRRef) (err error) {
	method := "VM.set_suspend_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeSRRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetSuspendSR2: Set the suspend_SR field of the given VM.
// Version: rio
func (vm) SetSuspendSR2(session *Session, self VMRef) (err error) {
	method := "VM.set_suspend_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveTags: Remove the given value from the tags field of the given VM.  If the value is not in that Set, then do nothing.
// Version: orlando
func (vm) RemoveTags(session *Session, self VMRef, value string) (err error) {
	method := "VM.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveTags3: Remove the given value from the tags field of the given VM.  If the value is not in that Set, then do nothing.
// Version: orlando
func (vm) RemoveTags3(session *Session, self VMRef, value string) (err error) {
	method := "VM.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveTags2: Remove the given value from the tags field of the given VM.  If the value is not in that Set, then do nothing.
// Version: rio
func (vm) RemoveTags2(session *Session, self VMRef) (err error) {
	method := "VM.remove_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddTags: Add the given value to the tags field of the given VM.  If the value is already in that Set, then do nothing.
// Version: orlando
func (vm) AddTags(session *Session, self VMRef, value string) (err error) {
	method := "VM.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddTags3: Add the given value to the tags field of the given VM.  If the value is already in that Set, then do nothing.
// Version: orlando
func (vm) AddTags3(session *Session, self VMRef, value string) (err error) {
	method := "VM.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddTags2: Add the given value to the tags field of the given VM.  If the value is already in that Set, then do nothing.
// Version: rio
func (vm) AddTags2(session *Session, self VMRef) (err error) {
	method := "VM.add_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetTags: Set the tags field of the given VM.
// Version: orlando
func (vm) SetTags(session *Session, self VMRef, value []string) (err error) {
	method := "VM.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetTags3: Set the tags field of the given VM.
// Version: orlando
func (vm) SetTags3(session *Session, self VMRef, value []string) (err error) {
	method := "VM.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetTags2: Set the tags field of the given VM.
// Version: rio
func (vm) SetTags2(session *Session, self VMRef) (err error) {
	method := "VM.set_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// RemoveFromXenstoreData: Remove the given key and its corresponding value from the xenstore_data field of the given VM.  If the key is not in that Map, then do nothing.
// Version: miami
func (vm) RemoveFromXenstoreData(session *Session, self VMRef, key string) (err error) {
	method := "VM.remove_from_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromXenstoreData3: Remove the given key and its corresponding value from the xenstore_data field of the given VM.  If the key is not in that Map, then do nothing.
// Version: miami
func (vm) RemoveFromXenstoreData3(session *Session, self VMRef, key string) (err error) {
	method := "VM.remove_from_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromXenstoreData2: Remove the given key and its corresponding value from the xenstore_data field of the given VM.  If the key is not in that Map, then do nothing.
// Version: rio
func (vm) RemoveFromXenstoreData2(session *Session, self VMRef) (err error) {
	method := "VM.remove_from_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AddToXenstoreData: Add the given key-value pair to the xenstore_data field of the given VM.
// Version: miami
func (vm) AddToXenstoreData(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToXenstoreData4: Add the given key-value pair to the xenstore_data field of the given VM.
// Version: miami
func (vm) AddToXenstoreData4(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToXenstoreData2: Add the given key-value pair to the xenstore_data field of the given VM.
// Version: rio
func (vm) AddToXenstoreData2(session *Session, self VMRef) (err error) {
	method := "VM.add_to_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetXenstoreData: Set the xenstore_data field of the given VM.
// Version: miami
func (vm) SetXenstoreData(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetXenstoreData3: Set the xenstore_data field of the given VM.
// Version: miami
func (vm) SetXenstoreData3(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetXenstoreData2: Set the xenstore_data field of the given VM.
// Version: rio
func (vm) SetXenstoreData2(session *Session, self VMRef) (err error) {
	method := "VM.set_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// SetRecommendations: Set the recommendations field of the given VM.
// Version: rio
func (vm) SetRecommendations(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetRecommendations3: Set the recommendations field of the given VM.
// Version: rio
func (vm) SetRecommendations3(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPCIBus: Set the PCI_bus field of the given VM.
// Version: rio
func (vm) SetPCIBus(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PCI_bus"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPCIBus3: Set the PCI_bus field of the given VM.
// Version: rio
func (vm) SetPCIBus3(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PCI_bus"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromHVMBootParams: Remove the given key and its corresponding value from the HVM/boot_params field of the given VM.  If the key is not in that Map, then do nothing.
// Version: rio
func (vm) RemoveFromHVMBootParams(session *Session, self VMRef, key string) (err error) {
	method := "VM.remove_from_HVM_boot_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromHVMBootParams3: Remove the given key and its corresponding value from the HVM/boot_params field of the given VM.  If the key is not in that Map, then do nothing.
// Version: rio
func (vm) RemoveFromHVMBootParams3(session *Session, self VMRef, key string) (err error) {
	method := "VM.remove_from_HVM_boot_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToHVMBootParams: Add the given key-value pair to the HVM/boot_params field of the given VM.
// Version: rio
func (vm) AddToHVMBootParams(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_HVM_boot_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToHVMBootParams4: Add the given key-value pair to the HVM/boot_params field of the given VM.
// Version: rio
func (vm) AddToHVMBootParams4(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_HVM_boot_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHVMBootParams: Set the HVM/boot_params field of the given VM.
// Version: rio
func (vm) SetHVMBootParams(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_HVM_boot_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetHVMBootParams3: Set the HVM/boot_params field of the given VM.
// Version: rio
func (vm) SetHVMBootParams3(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_HVM_boot_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPVLegacyArgs: Set the PV/legacy_args field of the given VM.
// Version: rio
func (vm) SetPVLegacyArgs(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PV_legacy_args"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPVLegacyArgs3: Set the PV/legacy_args field of the given VM.
// Version: rio
func (vm) SetPVLegacyArgs3(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PV_legacy_args"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPVBootloaderArgs: Set the PV/bootloader_args field of the given VM.
// Version: rio
func (vm) SetPVBootloaderArgs(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PV_bootloader_args"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPVBootloaderArgs3: Set the PV/bootloader_args field of the given VM.
// Version: rio
func (vm) SetPVBootloaderArgs3(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PV_bootloader_args"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPVArgs: Set the PV/args field of the given VM.
// Version: rio
func (vm) SetPVArgs(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PV_args"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPVArgs3: Set the PV/args field of the given VM.
// Version: rio
func (vm) SetPVArgs3(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PV_args"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPVRamdisk: Set the PV/ramdisk field of the given VM.
// Version: rio
func (vm) SetPVRamdisk(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PV_ramdisk"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPVRamdisk3: Set the PV/ramdisk field of the given VM.
// Version: rio
func (vm) SetPVRamdisk3(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PV_ramdisk"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPVKernel: Set the PV/kernel field of the given VM.
// Version: rio
func (vm) SetPVKernel(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PV_kernel"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPVKernel3: Set the PV/kernel field of the given VM.
// Version: rio
func (vm) SetPVKernel3(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PV_kernel"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPVBootloader: Set the PV/bootloader field of the given VM.
// Version: rio
func (vm) SetPVBootloader(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PV_bootloader"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetPVBootloader3: Set the PV/bootloader field of the given VM.
// Version: rio
func (vm) SetPVBootloader3(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_PV_bootloader"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetActionsAfterReboot: Set the actions/after_reboot field of the given VM.
// Version: rio
func (vm) SetActionsAfterReboot(session *Session, self VMRef, value OnNormalExit) (err error) {
	method := "VM.set_actions_after_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnNormalExit(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetActionsAfterReboot3: Set the actions/after_reboot field of the given VM.
// Version: rio
func (vm) SetActionsAfterReboot3(session *Session, self VMRef, value OnNormalExit) (err error) {
	method := "VM.set_actions_after_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnNormalExit(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetActionsAfterShutdown: Set the actions/after_shutdown field of the given VM.
// Version: rio
func (vm) SetActionsAfterShutdown(session *Session, self VMRef, value OnNormalExit) (err error) {
	method := "VM.set_actions_after_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnNormalExit(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetActionsAfterShutdown3: Set the actions/after_shutdown field of the given VM.
// Version: rio
func (vm) SetActionsAfterShutdown3(session *Session, self VMRef, value OnNormalExit) (err error) {
	method := "VM.set_actions_after_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnNormalExit(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetActionsAfterSoftreboot: Set the actions/after_softreboot field of the given VM.
// Version: rio
func (vm) SetActionsAfterSoftreboot(session *Session, self VMRef, value OnSoftrebootBehavior) (err error) {
	method := "VM.set_actions_after_softreboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnSoftrebootBehavior(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetActionsAfterSoftreboot3: Set the actions/after_softreboot field of the given VM.
// Version: rio
func (vm) SetActionsAfterSoftreboot3(session *Session, self VMRef, value OnSoftrebootBehavior) (err error) {
	method := "VM.set_actions_after_softreboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnSoftrebootBehavior(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetActionsAfterSoftreboot2: Set the actions/after_softreboot field of the given VM.
// Version: closed
func (vm) SetActionsAfterSoftreboot2(session *Session, value OnSoftrebootBehavior) (err error) {
	method := "VM.set_actions_after_softreboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumOnSoftrebootBehavior(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, valueArg)
	return
}

// RemoveFromVCPUsParams: Remove the given key and its corresponding value from the VCPUs/params field of the given VM.  If the key is not in that Map, then do nothing.
// Version: rio
func (vm) RemoveFromVCPUsParams(session *Session, self VMRef, key string) (err error) {
	method := "VM.remove_from_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromVCPUsParams3: Remove the given key and its corresponding value from the VCPUs/params field of the given VM.  If the key is not in that Map, then do nothing.
// Version: rio
func (vm) RemoveFromVCPUsParams3(session *Session, self VMRef, key string) (err error) {
	method := "VM.remove_from_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToVCPUsParams: Add the given key-value pair to the VCPUs/params field of the given VM.
// Version: rio
func (vm) AddToVCPUsParams(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToVCPUsParams4: Add the given key-value pair to the VCPUs/params field of the given VM.
// Version: rio
func (vm) AddToVCPUsParams4(session *Session, self VMRef, key string, value string) (err error) {
	method := "VM.add_to_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetVCPUsParams: Set the VCPUs/params field of the given VM.
// Version: rio
func (vm) SetVCPUsParams(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetVCPUsParams3: Set the VCPUs/params field of the given VM.
// Version: rio
func (vm) SetVCPUsParams3(session *Session, self VMRef, value map[string]string) (err error) {
	method := "VM.set_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetAffinity: Set the affinity field of the given VM.
// Version: rio
func (vm) SetAffinity(session *Session, self VMRef, value HostRef) (err error) {
	method := "VM.set_affinity"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetAffinity3: Set the affinity field of the given VM.
// Version: rio
func (vm) SetAffinity3(session *Session, self VMRef, value HostRef) (err error) {
	method := "VM.set_affinity"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeHostRef(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetIsATemplate: Set the is_a_template field of the given VM.
// Version: rio
func (vm) SetIsATemplate(session *Session, self VMRef, value bool) (err error) {
	method := "VM.set_is_a_template"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetIsATemplate3: Set the is_a_template field of the given VM.
// Version: rio
func (vm) SetIsATemplate3(session *Session, self VMRef, value bool) (err error) {
	method := "VM.set_is_a_template"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetUserVersion: Set the user_version field of the given VM.
// Version: rio
func (vm) SetUserVersion(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_user_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetUserVersion3: Set the user_version field of the given VM.
// Version: rio
func (vm) SetUserVersion3(session *Session, self VMRef, value int) (err error) {
	method := "VM.set_user_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetNameDescription: Set the name/description field of the given VM.
// Version: rio
func (vm) SetNameDescription(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameDescription3: Set the name/description field of the given VM.
// Version: rio
func (vm) SetNameDescription3(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameLabel: Set the name/label field of the given VM.
// Version: rio
func (vm) SetNameLabel(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameLabel3: Set the name/label field of the given VM.
// Version: rio
func (vm) SetNameLabel3(session *Session, self VMRef, value string) (err error) {
	method := "VM.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSecurebootCertificatesState: Get the secureboot_certificates_state field of the given VM.
// Version: rio
func (vm) GetSecurebootCertificatesState(session *Session, self VMRef) (retval VMSecurebootCertificatesState, err error) {
	method := "VM.get_secureboot_certificates_state"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVMSecurebootCertificatesState(method+" -> ", result)
	return
}

// GetSecurebootCertificatesState2: Get the secureboot_certificates_state field of the given VM.
// Version: rio
func (vm) GetSecurebootCertificatesState2(session *Session, self VMRef) (retval VMSecurebootCertificatesState, err error) {
	method := "VM.get_secureboot_certificates_state"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVMSecurebootCertificatesState(method+" -> ", result)
	return
}

// GetGroups: Get the groups field of the given VM.
// Version: rio
func (vm) GetGroups(session *Session, self VMRef) (retval []VMGroupRef, err error) {
	method := "VM.get_groups"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMGroupRefSet(method+" -> ", result)
	return
}

// GetGroups2: Get the groups field of the given VM.
// Version: rio
func (vm) GetGroups2(session *Session, self VMRef) (retval []VMGroupRef, err error) {
	method := "VM.get_groups"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMGroupRefSet(method+" -> ", result)
	return
}

// GetPendingGuidancesFull: Get the pending_guidances_full field of the given VM.
// Version: rio
func (vm) GetPendingGuidancesFull(session *Session, self VMRef) (retval []UpdateGuidances, err error) {
	method := "VM.get_pending_guidances_full"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateGuidancesSet(method+" -> ", result)
	return
}

// GetPendingGuidancesFull2: Get the pending_guidances_full field of the given VM.
// Version: rio
func (vm) GetPendingGuidancesFull2(session *Session, self VMRef) (retval []UpdateGuidances, err error) {
	method := "VM.get_pending_guidances_full"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateGuidancesSet(method+" -> ", result)
	return
}

// GetPendingGuidancesRecommended: Get the pending_guidances_recommended field of the given VM.
// Version: rio
func (vm) GetPendingGuidancesRecommended(session *Session, self VMRef) (retval []UpdateGuidances, err error) {
	method := "VM.get_pending_guidances_recommended"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateGuidancesSet(method+" -> ", result)
	return
}

// GetPendingGuidancesRecommended2: Get the pending_guidances_recommended field of the given VM.
// Version: rio
func (vm) GetPendingGuidancesRecommended2(session *Session, self VMRef) (retval []UpdateGuidances, err error) {
	method := "VM.get_pending_guidances_recommended"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateGuidancesSet(method+" -> ", result)
	return
}

// GetPendingGuidances: Get the pending_guidances field of the given VM.
// Version: rio
func (vm) GetPendingGuidances(session *Session, self VMRef) (retval []UpdateGuidances, err error) {
	method := "VM.get_pending_guidances"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateGuidancesSet(method+" -> ", result)
	return
}

// GetPendingGuidances2: Get the pending_guidances field of the given VM.
// Version: rio
func (vm) GetPendingGuidances2(session *Session, self VMRef) (retval []UpdateGuidances, err error) {
	method := "VM.get_pending_guidances"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumUpdateGuidancesSet(method+" -> ", result)
	return
}

// GetNVRAM: Get the NVRAM field of the given VM.
// Version: rio
func (vm) GetNVRAM(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_NVRAM"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNVRAM2: Get the NVRAM field of the given VM.
// Version: rio
func (vm) GetNVRAM2(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_NVRAM"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDomainType: Get the domain_type field of the given VM.
// Version: rio
func (vm) GetDomainType(session *Session, self VMRef) (retval DomainType, err error) {
	method := "VM.get_domain_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDomainType2: Get the domain_type field of the given VM.
// Version: rio
func (vm) GetDomainType2(session *Session, self VMRef) (retval DomainType, err error) {
	method := "VM.get_domain_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetReferenceLabel: Get the reference_label field of the given VM.
// Version: rio
func (vm) GetReferenceLabel(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_reference_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetReferenceLabel2: Get the reference_label field of the given VM.
// Version: rio
func (vm) GetReferenceLabel2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_reference_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRequiresReboot: Get the requires_reboot field of the given VM.
// Version: rio
func (vm) GetRequiresReboot(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_requires_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRequiresReboot2: Get the requires_reboot field of the given VM.
// Version: rio
func (vm) GetRequiresReboot2(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_requires_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHasVendorDevice: Get the has_vendor_device field of the given VM.
// Version: rio
func (vm) GetHasVendorDevice(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_has_vendor_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHasVendorDevice2: Get the has_vendor_device field of the given VM.
// Version: rio
func (vm) GetHasVendorDevice2(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_has_vendor_device"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHardwarePlatformVersion: Get the hardware_platform_version field of the given VM.
// Version: rio
func (vm) GetHardwarePlatformVersion(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_hardware_platform_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHardwarePlatformVersion2: Get the hardware_platform_version field of the given VM.
// Version: rio
func (vm) GetHardwarePlatformVersion2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_hardware_platform_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetGenerationID: Get the generation_id field of the given VM.
// Version: rio
func (vm) GetGenerationID(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_generation_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetGenerationID2: Get the generation_id field of the given VM.
// Version: rio
func (vm) GetGenerationID2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_generation_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVersion: Get the version field of the given VM.
// Version: rio
func (vm) GetVersion(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVersion2: Get the version field of the given VM.
// Version: rio
func (vm) GetVersion2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSuspendSR: Get the suspend_SR field of the given VM.
// Version: rio
func (vm) GetSuspendSR(session *Session, self VMRef) (retval SRRef, err error) {
	method := "VM.get_suspend_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSuspendSR2: Get the suspend_SR field of the given VM.
// Version: rio
func (vm) GetSuspendSR2(session *Session, self VMRef) (retval SRRef, err error) {
	method := "VM.get_suspend_SR"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAttachedPCIs: Get the attached_PCIs field of the given VM.
// Version: rio
func (vm) GetAttachedPCIs(session *Session, self VMRef) (retval []PCIRef, err error) {
	method := "VM.get_attached_PCIs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRefSet(method+" -> ", result)
	return
}

// GetAttachedPCIs2: Get the attached_PCIs field of the given VM.
// Version: rio
func (vm) GetAttachedPCIs2(session *Session, self VMRef) (retval []PCIRef, err error) {
	method := "VM.get_attached_PCIs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializePCIRefSet(method+" -> ", result)
	return
}

// GetVGPUs: Get the VGPUs field of the given VM.
// Version: rio
func (vm) GetVGPUs(session *Session, self VMRef) (retval []VGPURef, err error) {
	method := "VM.get_VGPUs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVGPURefSet(method+" -> ", result)
	return
}

// GetVGPUs2: Get the VGPUs field of the given VM.
// Version: rio
func (vm) GetVGPUs2(session *Session, self VMRef) (retval []VGPURef, err error) {
	method := "VM.get_VGPUs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVGPURefSet(method+" -> ", result)
	return
}

// GetOrder: Get the order field of the given VM.
// Version: rio
func (vm) GetOrder(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_order"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOrder2: Get the order field of the given VM.
// Version: rio
func (vm) GetOrder2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_order"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetShutdownDelay: Get the shutdown_delay field of the given VM.
// Version: rio
func (vm) GetShutdownDelay(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_shutdown_delay"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetShutdownDelay2: Get the shutdown_delay field of the given VM.
// Version: rio
func (vm) GetShutdownDelay2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_shutdown_delay"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetStartDelay: Get the start_delay field of the given VM.
// Version: rio
func (vm) GetStartDelay(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_start_delay"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetStartDelay2: Get the start_delay field of the given VM.
// Version: rio
func (vm) GetStartDelay2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_start_delay"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAppliance: Get the appliance field of the given VM.
// Version: rio
func (vm) GetAppliance(session *Session, self VMRef) (retval VMApplianceRef, err error) {
	method := "VM.get_appliance"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMApplianceRef(method+" -> ", result)
	return
}

// GetAppliance2: Get the appliance field of the given VM.
// Version: rio
func (vm) GetAppliance2(session *Session, self VMRef) (retval VMApplianceRef, err error) {
	method := "VM.get_appliance"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMApplianceRef(method+" -> ", result)
	return
}

// GetIsVmssSnapshot: Get the is_vmss_snapshot field of the given VM.
// Version: rio
func (vm) GetIsVmssSnapshot(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_is_vmss_snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsVmssSnapshot2: Get the is_vmss_snapshot field of the given VM.
// Version: rio
func (vm) GetIsVmssSnapshot2(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_is_vmss_snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSnapshotSchedule: Get the snapshot_schedule field of the given VM.
// Version: rio
func (vm) GetSnapshotSchedule(session *Session, self VMRef) (retval VMSSRef, err error) {
	method := "VM.get_snapshot_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMSSRef(method+" -> ", result)
	return
}

// GetSnapshotSchedule2: Get the snapshot_schedule field of the given VM.
// Version: rio
func (vm) GetSnapshotSchedule2(session *Session, self VMRef) (retval VMSSRef, err error) {
	method := "VM.get_snapshot_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMSSRef(method+" -> ", result)
	return
}

// GetIsSnapshotFromVmpp: Get the is_snapshot_from_vmpp field of the given VM.
// Version: rio
func (vm) GetIsSnapshotFromVmpp(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_is_snapshot_from_vmpp"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsSnapshotFromVmpp2: Get the is_snapshot_from_vmpp field of the given VM.
// Version: rio
func (vm) GetIsSnapshotFromVmpp2(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_is_snapshot_from_vmpp"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetProtectionPolicy: Get the protection_policy field of the given VM.
// Version: rio
func (vm) GetProtectionPolicy(session *Session, self VMRef) (retval VMPPRef, err error) {
	method := "VM.get_protection_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMPPRef(method+" -> ", result)
	return
}

// GetProtectionPolicy2: Get the protection_policy field of the given VM.
// Version: rio
func (vm) GetProtectionPolicy2(session *Session, self VMRef) (retval VMPPRef, err error) {
	method := "VM.get_protection_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMPPRef(method+" -> ", result)
	return
}

// GetBiosStrings: Get the bios_strings field of the given VM.
// Version: rio
func (vm) GetBiosStrings(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_bios_strings"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBiosStrings2: Get the bios_strings field of the given VM.
// Version: rio
func (vm) GetBiosStrings2(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_bios_strings"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetChildren: Get the children field of the given VM.
// Version: rio
func (vm) GetChildren(session *Session, self VMRef) (retval []VMRef, err error) {
	method := "VM.get_children"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// GetChildren2: Get the children field of the given VM.
// Version: rio
func (vm) GetChildren2(session *Session, self VMRef) (retval []VMRef, err error) {
	method := "VM.get_children"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// GetParent: Get the parent field of the given VM.
// Version: rio
func (vm) GetParent(session *Session, self VMRef) (retval VMRef, err error) {
	method := "VM.get_parent"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// GetParent2: Get the parent field of the given VM.
// Version: rio
func (vm) GetParent2(session *Session, self VMRef) (retval VMRef, err error) {
	method := "VM.get_parent"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// GetSnapshotMetadata: Get the snapshot_metadata field of the given VM.
// Version: rio
func (vm) GetSnapshotMetadata(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_snapshot_metadata"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSnapshotMetadata2: Get the snapshot_metadata field of the given VM.
// Version: rio
func (vm) GetSnapshotMetadata2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_snapshot_metadata"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSnapshotInfo: Get the snapshot_info field of the given VM.
// Version: rio
func (vm) GetSnapshotInfo(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_snapshot_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSnapshotInfo2: Get the snapshot_info field of the given VM.
// Version: rio
func (vm) GetSnapshotInfo2(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_snapshot_info"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBlockedOperations: Get the blocked_operations field of the given VM.
// Version: rio
func (vm) GetBlockedOperations(session *Session, self VMRef) (retval map[VMOperations]string, err error) {
	method := "VM.get_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVMOperationsToStringMap(method+" -> ", result)
	return
}

// GetBlockedOperations2: Get the blocked_operations field of the given VM.
// Version: rio
func (vm) GetBlockedOperations2(session *Session, self VMRef) (retval map[VMOperations]string, err error) {
	method := "VM.get_blocked_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVMOperationsToStringMap(method+" -> ", result)
	return
}

// GetTags: Get the tags field of the given VM.
// Version: rio
func (vm) GetTags(session *Session, self VMRef) (retval []string, err error) {
	method := "VM.get_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTags2: Get the tags field of the given VM.
// Version: rio
func (vm) GetTags2(session *Session, self VMRef) (retval []string, err error) {
	method := "VM.get_tags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBlobs: Get the blobs field of the given VM.
// Version: rio
func (vm) GetBlobs(session *Session, self VMRef) (retval map[string]BlobRef, err error) {
	method := "VM.get_blobs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToBlobRefMap(method+" -> ", result)
	return
}

// GetBlobs2: Get the blobs field of the given VM.
// Version: rio
func (vm) GetBlobs2(session *Session, self VMRef) (retval map[string]BlobRef, err error) {
	method := "VM.get_blobs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToBlobRefMap(method+" -> ", result)
	return
}

// GetTransportableSnapshotID: Get the transportable_snapshot_id field of the given VM.
// Version: rio
func (vm) GetTransportableSnapshotID(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_transportable_snapshot_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetTransportableSnapshotID2: Get the transportable_snapshot_id field of the given VM.
// Version: rio
func (vm) GetTransportableSnapshotID2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_transportable_snapshot_id"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSnapshotTime: Get the snapshot_time field of the given VM.
// Version: rio
func (vm) GetSnapshotTime(session *Session, self VMRef) (retval time.Time, err error) {
	method := "VM.get_snapshot_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSnapshotTime2: Get the snapshot_time field of the given VM.
// Version: rio
func (vm) GetSnapshotTime2(session *Session, self VMRef) (retval time.Time, err error) {
	method := "VM.get_snapshot_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSnapshots: Get the snapshots field of the given VM.
// Version: rio
func (vm) GetSnapshots(session *Session, self VMRef) (retval []VMRef, err error) {
	method := "VM.get_snapshots"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// GetSnapshots2: Get the snapshots field of the given VM.
// Version: rio
func (vm) GetSnapshots2(session *Session, self VMRef) (retval []VMRef, err error) {
	method := "VM.get_snapshots"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// GetSnapshotOf: Get the snapshot_of field of the given VM.
// Version: rio
func (vm) GetSnapshotOf(session *Session, self VMRef) (retval VMRef, err error) {
	method := "VM.get_snapshot_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// GetSnapshotOf2: Get the snapshot_of field of the given VM.
// Version: rio
func (vm) GetSnapshotOf2(session *Session, self VMRef) (retval VMRef, err error) {
	method := "VM.get_snapshot_of"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// GetIsASnapshot: Get the is_a_snapshot field of the given VM.
// Version: rio
func (vm) GetIsASnapshot(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_is_a_snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsASnapshot2: Get the is_a_snapshot field of the given VM.
// Version: rio
func (vm) GetIsASnapshot2(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_is_a_snapshot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaRestartPriority: Get the ha_restart_priority field of the given VM.
// Version: rio
func (vm) GetHaRestartPriority(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_ha_restart_priority"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaRestartPriority2: Get the ha_restart_priority field of the given VM.
// Version: rio
func (vm) GetHaRestartPriority2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_ha_restart_priority"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaAlwaysRun: Get the ha_always_run field of the given VM.
// Version: rio
func (vm) GetHaAlwaysRun(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_ha_always_run"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHaAlwaysRun2: Get the ha_always_run field of the given VM.
// Version: rio
func (vm) GetHaAlwaysRun2(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_ha_always_run"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetXenstoreData: Get the xenstore_data field of the given VM.
// Version: rio
func (vm) GetXenstoreData(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetXenstoreData2: Get the xenstore_data field of the given VM.
// Version: rio
func (vm) GetXenstoreData2(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_xenstore_data"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRecommendations: Get the recommendations field of the given VM.
// Version: rio
func (vm) GetRecommendations(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRecommendations2: Get the recommendations field of the given VM.
// Version: rio
func (vm) GetRecommendations2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_recommendations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastBootedRecord: Get the last_booted_record field of the given VM.
// Version: rio
func (vm) GetLastBootedRecord(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_last_booted_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastBootedRecord2: Get the last_booted_record field of the given VM.
// Version: rio
func (vm) GetLastBootedRecord2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_last_booted_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetGuestMetrics: Get the guest_metrics field of the given VM.
// Version: rio
func (vm) GetGuestMetrics(session *Session, self VMRef) (retval VMGuestMetricsRef, err error) {
	method := "VM.get_guest_metrics"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMGuestMetricsRef(method+" -> ", result)
	return
}

// GetGuestMetrics2: Get the guest_metrics field of the given VM.
// Version: rio
func (vm) GetGuestMetrics2(session *Session, self VMRef) (retval VMGuestMetricsRef, err error) {
	method := "VM.get_guest_metrics"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMGuestMetricsRef(method+" -> ", result)
	return
}

// GetMetrics: Get the metrics field of the given VM.
// Version: rio
func (vm) GetMetrics(session *Session, self VMRef) (retval VMMetricsRef, err error) {
	method := "VM.get_metrics"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMMetricsRef(method+" -> ", result)
	return
}

// GetMetrics2: Get the metrics field of the given VM.
// Version: rio
func (vm) GetMetrics2(session *Session, self VMRef) (retval VMMetricsRef, err error) {
	method := "VM.get_metrics"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMMetricsRef(method+" -> ", result)
	return
}

// GetIsControlDomain: Get the is_control_domain field of the given VM.
// Version: rio
func (vm) GetIsControlDomain(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_is_control_domain"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsControlDomain2: Get the is_control_domain field of the given VM.
// Version: rio
func (vm) GetIsControlDomain2(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_is_control_domain"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastBootCPUFlags: Get the last_boot_CPU_flags field of the given VM.
// Version: rio
func (vm) GetLastBootCPUFlags(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_last_boot_CPU_flags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetLastBootCPUFlags2: Get the last_boot_CPU_flags field of the given VM.
// Version: rio
func (vm) GetLastBootCPUFlags2(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_last_boot_CPU_flags"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDomarch: Get the domarch field of the given VM.
// Version: rio
func (vm) GetDomarch(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_domarch"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDomarch2: Get the domarch field of the given VM.
// Version: rio
func (vm) GetDomarch2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_domarch"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDomid: Get the domid field of the given VM.
// Version: rio
func (vm) GetDomid(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_domid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetDomid2: Get the domid field of the given VM.
// Version: rio
func (vm) GetDomid2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_domid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig: Get the other_config field of the given VM.
// Version: rio
func (vm) GetOtherConfig(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetOtherConfig2: Get the other_config field of the given VM.
// Version: rio
func (vm) GetOtherConfig2(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_other_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPCIBus: Get the PCI_bus field of the given VM.
// Version: rio
func (vm) GetPCIBus(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PCI_bus"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPCIBus2: Get the PCI_bus field of the given VM.
// Version: rio
func (vm) GetPCIBus2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PCI_bus"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPlatform: Get the platform field of the given VM.
// Version: rio
func (vm) GetPlatform(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_platform"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPlatform2: Get the platform field of the given VM.
// Version: rio
func (vm) GetPlatform2(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_platform"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHVMShadowMultiplier: Get the HVM/shadow_multiplier field of the given VM.
// Version: rio
func (vm) GetHVMShadowMultiplier(session *Session, self VMRef) (retval float64, err error) {
	method := "VM.get_HVM_shadow_multiplier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHVMShadowMultiplier2: Get the HVM/shadow_multiplier field of the given VM.
// Version: rio
func (vm) GetHVMShadowMultiplier2(session *Session, self VMRef) (retval float64, err error) {
	method := "VM.get_HVM_shadow_multiplier"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHVMBootParams: Get the HVM/boot_params field of the given VM.
// Version: rio
func (vm) GetHVMBootParams(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_HVM_boot_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHVMBootParams2: Get the HVM/boot_params field of the given VM.
// Version: rio
func (vm) GetHVMBootParams2(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_HVM_boot_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHVMBootPolicy: Get the HVM/boot_policy field of the given VM.
// Version: rio
func (vm) GetHVMBootPolicy(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_HVM_boot_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetHVMBootPolicy2: Get the HVM/boot_policy field of the given VM.
// Version: rio
func (vm) GetHVMBootPolicy2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_HVM_boot_policy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVLegacyArgs: Get the PV/legacy_args field of the given VM.
// Version: rio
func (vm) GetPVLegacyArgs(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PV_legacy_args"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVLegacyArgs2: Get the PV/legacy_args field of the given VM.
// Version: rio
func (vm) GetPVLegacyArgs2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PV_legacy_args"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVBootloaderArgs: Get the PV/bootloader_args field of the given VM.
// Version: rio
func (vm) GetPVBootloaderArgs(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PV_bootloader_args"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVBootloaderArgs2: Get the PV/bootloader_args field of the given VM.
// Version: rio
func (vm) GetPVBootloaderArgs2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PV_bootloader_args"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVArgs: Get the PV/args field of the given VM.
// Version: rio
func (vm) GetPVArgs(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PV_args"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVArgs2: Get the PV/args field of the given VM.
// Version: rio
func (vm) GetPVArgs2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PV_args"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVRamdisk: Get the PV/ramdisk field of the given VM.
// Version: rio
func (vm) GetPVRamdisk(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PV_ramdisk"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVRamdisk2: Get the PV/ramdisk field of the given VM.
// Version: rio
func (vm) GetPVRamdisk2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PV_ramdisk"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVKernel: Get the PV/kernel field of the given VM.
// Version: rio
func (vm) GetPVKernel(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PV_kernel"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVKernel2: Get the PV/kernel field of the given VM.
// Version: rio
func (vm) GetPVKernel2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PV_kernel"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVBootloader: Get the PV/bootloader field of the given VM.
// Version: rio
func (vm) GetPVBootloader(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PV_bootloader"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPVBootloader2: Get the PV/bootloader field of the given VM.
// Version: rio
func (vm) GetPVBootloader2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_PV_bootloader"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVTPMs: Get the VTPMs field of the given VM.
// Version: rio
func (vm) GetVTPMs(session *Session, self VMRef) (retval []VTPMRef, err error) {
	method := "VM.get_VTPMs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVTPMRefSet(method+" -> ", result)
	return
}

// GetVTPMs2: Get the VTPMs field of the given VM.
// Version: rio
func (vm) GetVTPMs2(session *Session, self VMRef) (retval []VTPMRef, err error) {
	method := "VM.get_VTPMs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVTPMRefSet(method+" -> ", result)
	return
}

// GetCrashDumps: Get the crash_dumps field of the given VM.
// Version: rio
func (vm) GetCrashDumps(session *Session, self VMRef) (retval []CrashdumpRef, err error) {
	method := "VM.get_crash_dumps"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCrashDumps2: Get the crash_dumps field of the given VM.
// Version: rio
func (vm) GetCrashDumps2(session *Session, self VMRef) (retval []CrashdumpRef, err error) {
	method := "VM.get_crash_dumps"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVUSBs: Get the VUSBs field of the given VM.
// Version: rio
func (vm) GetVUSBs(session *Session, self VMRef) (retval []VUSBRef, err error) {
	method := "VM.get_VUSBs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVUSBRefSet(method+" -> ", result)
	return
}

// GetVUSBs2: Get the VUSBs field of the given VM.
// Version: rio
func (vm) GetVUSBs2(session *Session, self VMRef) (retval []VUSBRef, err error) {
	method := "VM.get_VUSBs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVUSBRefSet(method+" -> ", result)
	return
}

// GetVBDs: Get the VBDs field of the given VM.
// Version: rio
func (vm) GetVBDs(session *Session, self VMRef) (retval []VBDRef, err error) {
	method := "VM.get_VBDs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVBDs2: Get the VBDs field of the given VM.
// Version: rio
func (vm) GetVBDs2(session *Session, self VMRef) (retval []VBDRef, err error) {
	method := "VM.get_VBDs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVIFs: Get the VIFs field of the given VM.
// Version: rio
func (vm) GetVIFs(session *Session, self VMRef) (retval []VIFRef, err error) {
	method := "VM.get_VIFs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVIFRefSet(method+" -> ", result)
	return
}

// GetVIFs2: Get the VIFs field of the given VM.
// Version: rio
func (vm) GetVIFs2(session *Session, self VMRef) (retval []VIFRef, err error) {
	method := "VM.get_VIFs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVIFRefSet(method+" -> ", result)
	return
}

// GetConsoles: Get the consoles field of the given VM.
// Version: rio
func (vm) GetConsoles(session *Session, self VMRef) (retval []ConsoleRef, err error) {
	method := "VM.get_consoles"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeConsoleRefSet(method+" -> ", result)
	return
}

// GetConsoles2: Get the consoles field of the given VM.
// Version: rio
func (vm) GetConsoles2(session *Session, self VMRef) (retval []ConsoleRef, err error) {
	method := "VM.get_consoles"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeConsoleRefSet(method+" -> ", result)
	return
}

// GetActionsAfterCrash: Get the actions/after_crash field of the given VM.
// Version: rio
func (vm) GetActionsAfterCrash(session *Session, self VMRef) (retval OnCrashBehaviour, err error) {
	method := "VM.get_actions_after_crash"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumOnCrashBehaviour(method+" -> ", result)
	return
}

// GetActionsAfterCrash2: Get the actions/after_crash field of the given VM.
// Version: rio
func (vm) GetActionsAfterCrash2(session *Session, self VMRef) (retval OnCrashBehaviour, err error) {
	method := "VM.get_actions_after_crash"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumOnCrashBehaviour(method+" -> ", result)
	return
}

// GetActionsAfterReboot: Get the actions/after_reboot field of the given VM.
// Version: rio
func (vm) GetActionsAfterReboot(session *Session, self VMRef) (retval OnNormalExit, err error) {
	method := "VM.get_actions_after_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumOnNormalExit(method+" -> ", result)
	return
}

// GetActionsAfterReboot2: Get the actions/after_reboot field of the given VM.
// Version: rio
func (vm) GetActionsAfterReboot2(session *Session, self VMRef) (retval OnNormalExit, err error) {
	method := "VM.get_actions_after_reboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumOnNormalExit(method+" -> ", result)
	return
}

// GetActionsAfterShutdown: Get the actions/after_shutdown field of the given VM.
// Version: rio
func (vm) GetActionsAfterShutdown(session *Session, self VMRef) (retval OnNormalExit, err error) {
	method := "VM.get_actions_after_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumOnNormalExit(method+" -> ", result)
	return
}

// GetActionsAfterShutdown2: Get the actions/after_shutdown field of the given VM.
// Version: rio
func (vm) GetActionsAfterShutdown2(session *Session, self VMRef) (retval OnNormalExit, err error) {
	method := "VM.get_actions_after_shutdown"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumOnNormalExit(method+" -> ", result)
	return
}

// GetActionsAfterSoftreboot: Get the actions/after_softreboot field of the given VM.
// Version: rio
func (vm) GetActionsAfterSoftreboot(session *Session, self VMRef) (retval OnSoftrebootBehavior, err error) {
	method := "VM.get_actions_after_softreboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumOnSoftrebootBehavior(method+" -> ", result)
	return
}

// GetActionsAfterSoftreboot2: Get the actions/after_softreboot field of the given VM.
// Version: rio
func (vm) GetActionsAfterSoftreboot2(session *Session, self VMRef) (retval OnSoftrebootBehavior, err error) {
	method := "VM.get_actions_after_softreboot"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumOnSoftrebootBehavior(method+" -> ", result)
	return
}

// GetVCPUsAtStartup: Get the VCPUs/at_startup field of the given VM.
// Version: rio
func (vm) GetVCPUsAtStartup(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_VCPUs_at_startup"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVCPUsAtStartup2: Get the VCPUs/at_startup field of the given VM.
// Version: rio
func (vm) GetVCPUsAtStartup2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_VCPUs_at_startup"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVCPUsMax: Get the VCPUs/max field of the given VM.
// Version: rio
func (vm) GetVCPUsMax(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_VCPUs_max"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVCPUsMax2: Get the VCPUs/max field of the given VM.
// Version: rio
func (vm) GetVCPUsMax2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_VCPUs_max"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVCPUsParams: Get the VCPUs/params field of the given VM.
// Version: rio
func (vm) GetVCPUsParams(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVCPUsParams2: Get the VCPUs/params field of the given VM.
// Version: rio
func (vm) GetVCPUsParams2(session *Session, self VMRef) (retval map[string]string, err error) {
	method := "VM.get_VCPUs_params"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryStaticMin: Get the memory/static_min field of the given VM.
// Version: rio
func (vm) GetMemoryStaticMin(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_memory_static_min"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryStaticMin2: Get the memory/static_min field of the given VM.
// Version: rio
func (vm) GetMemoryStaticMin2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_memory_static_min"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryDynamicMin: Get the memory/dynamic_min field of the given VM.
// Version: rio
func (vm) GetMemoryDynamicMin(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_memory_dynamic_min"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryDynamicMin2: Get the memory/dynamic_min field of the given VM.
// Version: rio
func (vm) GetMemoryDynamicMin2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_memory_dynamic_min"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryDynamicMax: Get the memory/dynamic_max field of the given VM.
// Version: rio
func (vm) GetMemoryDynamicMax(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_memory_dynamic_max"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryDynamicMax2: Get the memory/dynamic_max field of the given VM.
// Version: rio
func (vm) GetMemoryDynamicMax2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_memory_dynamic_max"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryStaticMax: Get the memory/static_max field of the given VM.
// Version: rio
func (vm) GetMemoryStaticMax(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_memory_static_max"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryStaticMax2: Get the memory/static_max field of the given VM.
// Version: rio
func (vm) GetMemoryStaticMax2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_memory_static_max"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryTarget: Get the memory/target field of the given VM.
// Version: rio
func (vm) GetMemoryTarget(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_memory_target"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryTarget2: Get the memory/target field of the given VM.
// Version: rio
func (vm) GetMemoryTarget2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_memory_target"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryOverhead: Get the memory/overhead field of the given VM.
// Version: rio
func (vm) GetMemoryOverhead(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_memory_overhead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetMemoryOverhead2: Get the memory/overhead field of the given VM.
// Version: rio
func (vm) GetMemoryOverhead2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_memory_overhead"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAffinity: Get the affinity field of the given VM.
// Version: rio
func (vm) GetAffinity(session *Session, self VMRef) (retval HostRef, err error) {
	method := "VM.get_affinity"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAffinity2: Get the affinity field of the given VM.
// Version: rio
func (vm) GetAffinity2(session *Session, self VMRef) (retval HostRef, err error) {
	method := "VM.get_affinity"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetScheduledToBeResidentOn: Get the scheduled_to_be_resident_on field of the given VM.
// Version: rio
func (vm) GetScheduledToBeResidentOn(session *Session, self VMRef) (retval HostRef, err error) {
	method := "VM.get_scheduled_to_be_resident_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetScheduledToBeResidentOn2: Get the scheduled_to_be_resident_on field of the given VM.
// Version: rio
func (vm) GetScheduledToBeResidentOn2(session *Session, self VMRef) (retval HostRef, err error) {
	method := "VM.get_scheduled_to_be_resident_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetResidentOn: Get the resident_on field of the given VM.
// Version: rio
func (vm) GetResidentOn(session *Session, self VMRef) (retval HostRef, err error) {
	method := "VM.get_resident_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetResidentOn2: Get the resident_on field of the given VM.
// Version: rio
func (vm) GetResidentOn2(session *Session, self VMRef) (retval HostRef, err error) {
	method := "VM.get_resident_on"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSuspendVDI: Get the suspend_VDI field of the given VM.
// Version: rio
func (vm) GetSuspendVDI(session *Session, self VMRef) (retval VDIRef, err error) {
	method := "VM.get_suspend_VDI"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetSuspendVDI2: Get the suspend_VDI field of the given VM.
// Version: rio
func (vm) GetSuspendVDI2(session *Session, self VMRef) (retval VDIRef, err error) {
	method := "VM.get_suspend_VDI"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsDefaultTemplate: Get the is_default_template field of the given VM.
// Version: rio
func (vm) GetIsDefaultTemplate(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_is_default_template"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsDefaultTemplate2: Get the is_default_template field of the given VM.
// Version: rio
func (vm) GetIsDefaultTemplate2(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_is_default_template"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsATemplate: Get the is_a_template field of the given VM.
// Version: rio
func (vm) GetIsATemplate(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_is_a_template"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsATemplate2: Get the is_a_template field of the given VM.
// Version: rio
func (vm) GetIsATemplate2(session *Session, self VMRef) (retval bool, err error) {
	method := "VM.get_is_a_template"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUserVersion: Get the user_version field of the given VM.
// Version: rio
func (vm) GetUserVersion(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_user_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUserVersion2: Get the user_version field of the given VM.
// Version: rio
func (vm) GetUserVersion2(session *Session, self VMRef) (retval int, err error) {
	method := "VM.get_user_version"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetPowerState: Get the power_state field of the given VM.
// Version: rio
func (vm) GetPowerState(session *Session, self VMRef) (retval VMPowerState, err error) {
	method := "VM.get_power_state"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVMPowerState(method+" -> ", result)
	return
}

// GetPowerState2: Get the power_state field of the given VM.
// Version: rio
func (vm) GetPowerState2(session *Session, self VMRef) (retval VMPowerState, err error) {
	method := "VM.get_power_state"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVMPowerState(method+" -> ", result)
	return
}

// GetNameDescription: Get the name/description field of the given VM.
// Version: rio
func (vm) GetNameDescription(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription2: Get the name/description field of the given VM.
// Version: rio
func (vm) GetNameDescription2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel: Get the name/label field of the given VM.
// Version: rio
func (vm) GetNameLabel(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel2: Get the name/label field of the given VM.
// Version: rio
func (vm) GetNameLabel2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetCurrentOperations: Get the current_operations field of the given VM.
// Version: rio
func (vm) GetCurrentOperations(session *Session, self VMRef) (retval map[string]VMOperations, err error) {
	method := "VM.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumVMOperationsMap(method+" -> ", result)
	return
}

// GetCurrentOperations2: Get the current_operations field of the given VM.
// Version: rio
func (vm) GetCurrentOperations2(session *Session, self VMRef) (retval map[string]VMOperations, err error) {
	method := "VM.get_current_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringToEnumVMOperationsMap(method+" -> ", result)
	return
}

// GetAllowedOperations: Get the allowed_operations field of the given VM.
// Version: rio
func (vm) GetAllowedOperations(session *Session, self VMRef) (retval []VMOperations, err error) {
	method := "VM.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVMOperationsSet(method+" -> ", result)
	return
}

// GetAllowedOperations2: Get the allowed_operations field of the given VM.
// Version: rio
func (vm) GetAllowedOperations2(session *Session, self VMRef) (retval []VMOperations, err error) {
	method := "VM.get_allowed_operations"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVMOperationsSet(method+" -> ", result)
	return
}

// GetUUID: Get the uuid field of the given VM.
// Version: rio
func (vm) GetUUID(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given VM.
// Version: rio
func (vm) GetUUID2(session *Session, self VMRef) (retval string, err error) {
	method := "VM.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByNameLabel: Get all the VM instances with the given label.
// Version: rio
func (vm) GetByNameLabel(session *Session, label string) (retval []VMRef, err error) {
	method := "VM.get_by_name_label"
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
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// GetByNameLabel2: Get all the VM instances with the given label.
// Version: rio
func (vm) GetByNameLabel2(session *Session, label string) (retval []VMRef, err error) {
	method := "VM.get_by_name_label"
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
	retval, err = deserializeVMRefSet(method+" -> ", result)
	return
}

// Destroy: Destroy the specified VM.  The VM is completely removed from the system.  This function can only be called when the VM is in the Halted State.
// Version: rio
func (vm) Destroy(session *Session, self VMRef) (err error) {
	method := "VM.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy the specified VM.  The VM is completely removed from the system.  This function can only be called when the VM is in the Halted State.
// Version: rio
func (vm) AsyncDestroy(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy the specified VM.  The VM is completely removed from the system.  This function can only be called when the VM is in the Halted State.
// Version: rio
func (vm) Destroy2(session *Session, self VMRef) (err error) {
	method := "VM.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy the specified VM.  The VM is completely removed from the system.  This function can only be called when the VM is in the Halted State.
// Version: rio
func (vm) AsyncDestroy2(session *Session, self VMRef) (retval TaskRef, err error) {
	method := "Async.VM.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: NOT RECOMMENDED! VM.clone or VM.copy (or VM.import) is a better choice in almost all situations. The standard way to obtain a new VM is to call VM.clone on a template VM, then call VM.provision on the new clone. Caution: if VM.create is used and then the new VM is attached to a virtual disc that has an operating system already installed, then there is no guarantee that the operating system will boot and run. Any software that calls VM.create on a future version of this API may fail or give unexpected results. For example this could happen if an additional parameter were added to VM.create. VM.create is intended only for use in the automatic creation of the system VM templates. It creates a new VM instance, and returns its handle. The constructor args are: name_label, name_description, power_state, user_version*, is_a_template*, suspend_VDI, affinity*, memory_target, memory_static_max*, memory_dynamic_max*, memory_dynamic_min*, memory_static_min*, VCPUs_params*, VCPUs_max*, VCPUs_at_startup*, actions_after_softreboot, actions_after_shutdown*, actions_after_reboot*, actions_after_crash*, PV_bootloader*, PV_kernel*, PV_ramdisk*, PV_args*, PV_bootloader_args*, PV_legacy_args*, HVM_boot_policy*, HVM_boot_params*, HVM_shadow_multiplier, platform*, PCI_bus*, other_config*, last_boot_CPU_flags, last_booted_record, recommendations*, xenstore_data, ha_always_run, ha_restart_priority, tags, blocked_operations, protection_policy, is_snapshot_from_vmpp, snapshot_schedule, is_vmss_snapshot, appliance, start_delay, shutdown_delay, order, suspend_SR, version, generation_id, hardware_platform_version, has_vendor_device, reference_label, domain_type, NVRAM (* = non-optional).
// Version: rio
func (vm) Create(session *Session, args VMRecord) (retval VMRef, err error) {
	method := "VM.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVMRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncCreate: NOT RECOMMENDED! VM.clone or VM.copy (or VM.import) is a better choice in almost all situations. The standard way to obtain a new VM is to call VM.clone on a template VM, then call VM.provision on the new clone. Caution: if VM.create is used and then the new VM is attached to a virtual disc that has an operating system already installed, then there is no guarantee that the operating system will boot and run. Any software that calls VM.create on a future version of this API may fail or give unexpected results. For example this could happen if an additional parameter were added to VM.create. VM.create is intended only for use in the automatic creation of the system VM templates. It creates a new VM instance, and returns its handle. The constructor args are: name_label, name_description, power_state, user_version*, is_a_template*, suspend_VDI, affinity*, memory_target, memory_static_max*, memory_dynamic_max*, memory_dynamic_min*, memory_static_min*, VCPUs_params*, VCPUs_max*, VCPUs_at_startup*, actions_after_softreboot, actions_after_shutdown*, actions_after_reboot*, actions_after_crash*, PV_bootloader*, PV_kernel*, PV_ramdisk*, PV_args*, PV_bootloader_args*, PV_legacy_args*, HVM_boot_policy*, HVM_boot_params*, HVM_shadow_multiplier, platform*, PCI_bus*, other_config*, last_boot_CPU_flags, last_booted_record, recommendations*, xenstore_data, ha_always_run, ha_restart_priority, tags, blocked_operations, protection_policy, is_snapshot_from_vmpp, snapshot_schedule, is_vmss_snapshot, appliance, start_delay, shutdown_delay, order, suspend_SR, version, generation_id, hardware_platform_version, has_vendor_device, reference_label, domain_type, NVRAM (* = non-optional).
// Version: rio
func (vm) AsyncCreate(session *Session, args VMRecord) (retval TaskRef, err error) {
	method := "Async.VM.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVMRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
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

// Create2: NOT RECOMMENDED! VM.clone or VM.copy (or VM.import) is a better choice in almost all situations. The standard way to obtain a new VM is to call VM.clone on a template VM, then call VM.provision on the new clone. Caution: if VM.create is used and then the new VM is attached to a virtual disc that has an operating system already installed, then there is no guarantee that the operating system will boot and run. Any software that calls VM.create on a future version of this API may fail or give unexpected results. For example this could happen if an additional parameter were added to VM.create. VM.create is intended only for use in the automatic creation of the system VM templates. It creates a new VM instance, and returns its handle. The constructor args are: name_label, name_description, power_state, user_version*, is_a_template*, suspend_VDI, affinity*, memory_target, memory_static_max*, memory_dynamic_max*, memory_dynamic_min*, memory_static_min*, VCPUs_params*, VCPUs_max*, VCPUs_at_startup*, actions_after_softreboot, actions_after_shutdown*, actions_after_reboot*, actions_after_crash*, PV_bootloader*, PV_kernel*, PV_ramdisk*, PV_args*, PV_bootloader_args*, PV_legacy_args*, HVM_boot_policy*, HVM_boot_params*, HVM_shadow_multiplier, platform*, PCI_bus*, other_config*, last_boot_CPU_flags, last_booted_record, recommendations*, xenstore_data, ha_always_run, ha_restart_priority, tags, blocked_operations, protection_policy, is_snapshot_from_vmpp, snapshot_schedule, is_vmss_snapshot, appliance, start_delay, shutdown_delay, order, suspend_SR, version, generation_id, hardware_platform_version, has_vendor_device, reference_label, domain_type, NVRAM (* = non-optional).
// Version: rio
func (vm) Create2(session *Session, args VMRecord) (retval VMRef, err error) {
	method := "VM.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVMRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// AsyncCreate2: NOT RECOMMENDED! VM.clone or VM.copy (or VM.import) is a better choice in almost all situations. The standard way to obtain a new VM is to call VM.clone on a template VM, then call VM.provision on the new clone. Caution: if VM.create is used and then the new VM is attached to a virtual disc that has an operating system already installed, then there is no guarantee that the operating system will boot and run. Any software that calls VM.create on a future version of this API may fail or give unexpected results. For example this could happen if an additional parameter were added to VM.create. VM.create is intended only for use in the automatic creation of the system VM templates. It creates a new VM instance, and returns its handle. The constructor args are: name_label, name_description, power_state, user_version*, is_a_template*, suspend_VDI, affinity*, memory_target, memory_static_max*, memory_dynamic_max*, memory_dynamic_min*, memory_static_min*, VCPUs_params*, VCPUs_max*, VCPUs_at_startup*, actions_after_softreboot, actions_after_shutdown*, actions_after_reboot*, actions_after_crash*, PV_bootloader*, PV_kernel*, PV_ramdisk*, PV_args*, PV_bootloader_args*, PV_legacy_args*, HVM_boot_policy*, HVM_boot_params*, HVM_shadow_multiplier, platform*, PCI_bus*, other_config*, last_boot_CPU_flags, last_booted_record, recommendations*, xenstore_data, ha_always_run, ha_restart_priority, tags, blocked_operations, protection_policy, is_snapshot_from_vmpp, snapshot_schedule, is_vmss_snapshot, appliance, start_delay, shutdown_delay, order, suspend_SR, version, generation_id, hardware_platform_version, has_vendor_device, reference_label, domain_type, NVRAM (* = non-optional).
// Version: rio
func (vm) AsyncCreate2(session *Session, args VMRecord) (retval TaskRef, err error) {
	method := "Async.VM.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVMRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
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

// GetByUUID: Get a reference to the VM instance with the specified UUID.
// Version: rio
func (vm) GetByUUID(session *Session, uuid string) (retval VMRef, err error) {
	method := "VM.get_by_uuid"
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
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the VM instance with the specified UUID.
// Version: rio
func (vm) GetByUUID2(session *Session, uuid string) (retval VMRef, err error) {
	method := "VM.get_by_uuid"
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
	retval, err = deserializeVMRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given VM.
// Version: rio
func (vm) GetRecord(session *Session, self VMRef) (retval VMRecord, err error) {
	method := "VM.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given VM.
// Version: rio
func (vm) GetRecord2(session *Session, self VMRef) (retval VMRecord, err error) {
	method := "VM.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMRecord(method+" -> ", result)
	return
}
