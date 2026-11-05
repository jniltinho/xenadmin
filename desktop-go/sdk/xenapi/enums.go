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

type VusbOperations string

const (
	// Attempting to attach this VUSB to a VM
	VusbOperationsAttach VusbOperations = "attach"
	// Attempting to plug this VUSB into a VM
	VusbOperationsPlug VusbOperations = "plug"
	// Attempting to hot unplug this VUSB
	VusbOperationsUnplug VusbOperations = "unplug"
	// The value does not belong to this enumeration
	VusbOperationsUnrecognized VusbOperations = "unrecognized"
)

type VtpmOperations string

const (
	// Destroy a VTPM
	VtpmOperationsDestroy VtpmOperations = "destroy"
	// The value does not belong to this enumeration
	VtpmOperationsUnrecognized VtpmOperations = "unrecognized"
)

type VmssType string

const (
	// The snapshot is a disk snapshot
	VmssTypeSnapshot VmssType = "snapshot"
	// The snapshot is a checkpoint
	VmssTypeCheckpoint VmssType = "checkpoint"
	// Support for VSS has been removed.
	VmssTypeSnapshotWithQuiesce VmssType = "snapshot_with_quiesce"
	// The value does not belong to this enumeration
	VmssTypeUnrecognized VmssType = "unrecognized"
)

type VmssFrequency string

const (
	// Hourly snapshots
	VmssFrequencyHourly VmssFrequency = "hourly"
	// Daily snapshots
	VmssFrequencyDaily VmssFrequency = "daily"
	// Weekly snapshots
	VmssFrequencyWeekly VmssFrequency = "weekly"
	// The value does not belong to this enumeration
	VmssFrequencyUnrecognized VmssFrequency = "unrecognized"
)

type VmppBackupType string

const (
	// The backup is a snapshot
	VmppBackupTypeSnapshot VmppBackupType = "snapshot"
	// The backup is a checkpoint
	VmppBackupTypeCheckpoint VmppBackupType = "checkpoint"
	// The value does not belong to this enumeration
	VmppBackupTypeUnrecognized VmppBackupType = "unrecognized"
)

type VmppBackupFrequency string

const (
	// Hourly backups
	VmppBackupFrequencyHourly VmppBackupFrequency = "hourly"
	// Daily backups
	VmppBackupFrequencyDaily VmppBackupFrequency = "daily"
	// Weekly backups
	VmppBackupFrequencyWeekly VmppBackupFrequency = "weekly"
	// The value does not belong to this enumeration
	VmppBackupFrequencyUnrecognized VmppBackupFrequency = "unrecognized"
)

type VmppArchiveTargetType string

const (
	// No target config
	VmppArchiveTargetTypeNone VmppArchiveTargetType = "none"
	// CIFS target config
	VmppArchiveTargetTypeCifs VmppArchiveTargetType = "cifs"
	// NFS target config
	VmppArchiveTargetTypeNfs VmppArchiveTargetType = "nfs"
	// The value does not belong to this enumeration
	VmppArchiveTargetTypeUnrecognized VmppArchiveTargetType = "unrecognized"
)

type VmppArchiveFrequency string

const (
	// Never archive
	VmppArchiveFrequencyNever VmppArchiveFrequency = "never"
	// Archive after backup
	VmppArchiveFrequencyAlwaysAfterBackup VmppArchiveFrequency = "always_after_backup"
	// Daily archives
	VmppArchiveFrequencyDaily VmppArchiveFrequency = "daily"
	// Weekly backups
	VmppArchiveFrequencyWeekly VmppArchiveFrequency = "weekly"
	// The value does not belong to this enumeration
	VmppArchiveFrequencyUnrecognized VmppArchiveFrequency = "unrecognized"
)

type VifOperations string

const (
	// Attempting to attach this VIF to a VM
	VifOperationsAttach VifOperations = "attach"
	// Attempting to hotplug this VIF
	VifOperationsPlug VifOperations = "plug"
	// Attempting to hot unplug this VIF
	VifOperationsUnplug VifOperations = "unplug"
	// The value does not belong to this enumeration
	VifOperationsUnrecognized VifOperations = "unrecognized"
)

type VifLockingMode string

const (
	// No specific configuration set - default network policy applies
	VifLockingModeNetworkDefault VifLockingMode = "network_default"
	// Only traffic to a specific MAC and a list of IPv4 or IPv6 addresses is permitted
	VifLockingModeLocked VifLockingMode = "locked"
	// All traffic is permitted
	VifLockingModeUnlocked VifLockingMode = "unlocked"
	// No traffic is permitted
	VifLockingModeDisabled VifLockingMode = "disabled"
	// The value does not belong to this enumeration
	VifLockingModeUnrecognized VifLockingMode = "unrecognized"
)

type VifIpv6ConfigurationMode string

const (
	// Follow the default IPv6 configuration of the guest (this is guest-dependent)
	VifIpv6ConfigurationModeNone VifIpv6ConfigurationMode = "None"
	// Static IPv6 address configuration
	VifIpv6ConfigurationModeStatic VifIpv6ConfigurationMode = "Static"
	// Acquire an IPv6 address automatically
	VifIpv6ConfigurationModeAutoconf VifIpv6ConfigurationMode = "Autoconf"
	// The value does not belong to this enumeration
	VifIpv6ConfigurationModeUnrecognized VifIpv6ConfigurationMode = "unrecognized"
)

type VifIpv4ConfigurationMode string

const (
	// Follow the default IPv4 configuration of the guest (this is guest-dependent)
	VifIpv4ConfigurationModeNone VifIpv4ConfigurationMode = "None"
	// Static IPv4 address configuration
	VifIpv4ConfigurationModeStatic VifIpv4ConfigurationMode = "Static"
	// Acquire an IP address by DHCP
	VifIpv4ConfigurationModeDHCP VifIpv4ConfigurationMode = "DHCP"
	// The value does not belong to this enumeration
	VifIpv4ConfigurationModeUnrecognized VifIpv4ConfigurationMode = "unrecognized"
)

type VgpuTypeImplementation string

const (
	// Pass through an entire physical GPU to a guest
	VgpuTypeImplementationPassthrough VgpuTypeImplementation = "passthrough"
	// vGPU using NVIDIA hardware
	VgpuTypeImplementationNvidia VgpuTypeImplementation = "nvidia"
	// vGPU using NVIDIA hardware with SR-IOV
	VgpuTypeImplementationNvidiaSriov VgpuTypeImplementation = "nvidia_sriov"
	// vGPU using Intel GVT-g
	VgpuTypeImplementationGvtG VgpuTypeImplementation = "gvt_g"
	// vGPU using AMD MxGPU
	VgpuTypeImplementationMxgpu VgpuTypeImplementation = "mxgpu"
	// The value does not belong to this enumeration
	VgpuTypeImplementationUnrecognized VgpuTypeImplementation = "unrecognized"
)

type VdiType string

const (
	// a disk that may be replaced on upgrade
	VdiTypeSystem VdiType = "system"
	// a disk that is always preserved on upgrade
	VdiTypeUser VdiType = "user"
	// a disk that may be reformatted on upgrade
	VdiTypeEphemeral VdiType = "ephemeral"
	// a disk that stores a suspend image
	VdiTypeSuspend VdiType = "suspend"
	// a disk that stores VM crashdump information
	VdiTypeCrashdump VdiType = "crashdump"
	// a disk used for HA storage heartbeating
	VdiTypeHaStatefile VdiType = "ha_statefile"
	// a disk used for HA Pool metadata
	VdiTypeMetadata VdiType = "metadata"
	// a disk used for a general metadata redo-log
	VdiTypeRedoLog VdiType = "redo_log"
	// a disk that stores SR-level RRDs
	VdiTypeRrd VdiType = "rrd"
	// a disk that stores PVS cache data
	VdiTypePvsCache VdiType = "pvs_cache"
	// Metadata about a snapshot VDI that has been deleted: the set of blocks that changed between some previous version of the disk and the version tracked by the snapshot.
	VdiTypeCbtMetadata VdiType = "cbt_metadata"
	// The value does not belong to this enumeration
	VdiTypeUnrecognized VdiType = "unrecognized"
)

type VdiOperations string

const (
	// Cloning the VDI
	VdiOperationsClone VdiOperations = "clone"
	// Copying the VDI
	VdiOperationsCopy VdiOperations = "copy"
	// Resizing the VDI
	VdiOperationsResize VdiOperations = "resize"
	// Resizing the VDI which may or may not be online
	VdiOperationsResizeOnline VdiOperations = "resize_online"
	// Snapshotting the VDI
	VdiOperationsSnapshot VdiOperations = "snapshot"
	// Mirroring the VDI
	VdiOperationsMirror VdiOperations = "mirror"
	// Destroying the VDI
	VdiOperationsDestroy VdiOperations = "destroy"
	// Forget about the VDI
	VdiOperationsForget VdiOperations = "forget"
	// Refreshing the fields of the VDI
	VdiOperationsUpdate VdiOperations = "update"
	// Forcibly unlocking the VDI
	VdiOperationsForceUnlock VdiOperations = "force_unlock"
	// Generating static configuration
	VdiOperationsGenerateConfig VdiOperations = "generate_config"
	// Enabling changed block tracking for a VDI
	VdiOperationsEnableCbt VdiOperations = "enable_cbt"
	// Disabling changed block tracking for a VDI
	VdiOperationsDisableCbt VdiOperations = "disable_cbt"
	// Deleting the data of the VDI
	VdiOperationsDataDestroy VdiOperations = "data_destroy"
	// Exporting a bitmap that shows the changed blocks between two VDIs
	VdiOperationsListChangedBlocks VdiOperations = "list_changed_blocks"
	// Setting the on_boot field of the VDI
	VdiOperationsSetOnBoot VdiOperations = "set_on_boot"
	// Operations on this VDI are temporarily blocked
	VdiOperationsBlocked VdiOperations = "blocked"
	// Reverting a VDI to a clone of this snapshot
	VdiOperationsRevertTo VdiOperations = "revert_to"
	// Reverting this VDI to a clone of a snapshot
	VdiOperationsRevertFrom VdiOperations = "revert_from"
	// The value does not belong to this enumeration
	VdiOperationsUnrecognized VdiOperations = "unrecognized"
)

type VbdType string

const (
	// VBD will appear to guest as CD
	VbdTypeCD VbdType = "CD"
	// VBD will appear to guest as disk
	VbdTypeDisk VbdType = "Disk"
	// VBD will appear as a floppy
	VbdTypeFloppy VbdType = "Floppy"
	// The value does not belong to this enumeration
	VbdTypeUnrecognized VbdType = "unrecognized"
)

type VbdOperations string

const (
	// Attempting to attach this VBD to a VM
	VbdOperationsAttach VbdOperations = "attach"
	// Attempting to eject the media from this VBD
	VbdOperationsEject VbdOperations = "eject"
	// Attempting to insert new media into this VBD
	VbdOperationsInsert VbdOperations = "insert"
	// Attempting to hotplug this VBD
	VbdOperationsPlug VbdOperations = "plug"
	// Attempting to hot unplug this VBD
	VbdOperationsUnplug VbdOperations = "unplug"
	// Attempting to forcibly unplug this VBD
	VbdOperationsUnplugForce VbdOperations = "unplug_force"
	// Attempting to pause a block device backend
	VbdOperationsPause VbdOperations = "pause"
	// Attempting to unpause a block device backend
	VbdOperationsUnpause VbdOperations = "unpause"
	// The value does not belong to this enumeration
	VbdOperationsUnrecognized VbdOperations = "unrecognized"
)

type VbdMode string

const (
	// only read-only access will be allowed
	VbdModeRO VbdMode = "RO"
	// read-write access will be allowed
	VbdModeRW VbdMode = "RW"
	// The value does not belong to this enumeration
	VbdModeUnrecognized VbdMode = "unrecognized"
)

type VMUefiMode string

const (
	// clears a VM&apos;s EFI variables related to Secure Boot and places it into Setup Mode
	VMUefiModeSetup VMUefiMode = "setup"
	// resets a VM&apos;s EFI variables related to Secure Boot to the defaults, placing it into User Mode
	VMUefiModeUser VMUefiMode = "user"
	// The value does not belong to this enumeration
	VMUefiModeUnrecognized VMUefiMode = "unrecognized"
)

type VMSecurebootReadiness string

const (
	// VM&apos;s firmware is not UEFI
	VMSecurebootReadinessNotSupported VMSecurebootReadiness = "not_supported"
	// Secureboot is disabled on this VM
	VMSecurebootReadinessDisabled VMSecurebootReadiness = "disabled"
	// Secured boot is enabled on this VM and its NVRAM.EFI-variables are empty
	VMSecurebootReadinessFirstBoot VMSecurebootReadiness = "first_boot"
	// Secured boot is enabled on this VM and PK, KEK, db and dbx are defined in its EFI variables
	VMSecurebootReadinessReady VMSecurebootReadiness = "ready"
	// Secured boot is enabled on this VM and PK, KEK, db but not dbx are defined in its EFI variables
	VMSecurebootReadinessReadyNoDbx VMSecurebootReadiness = "ready_no_dbx"
	// Secured boot is enabled on this VM and PK is not defined in its EFI variables
	VMSecurebootReadinessSetupMode VMSecurebootReadiness = "setup_mode"
	// Secured boot is enabled on this VM and the certificates defined in its EFI variables are incomplete
	VMSecurebootReadinessCertsIncomplete VMSecurebootReadiness = "certs_incomplete"
	// The value does not belong to this enumeration
	VMSecurebootReadinessUnrecognized VMSecurebootReadiness = "unrecognized"
)

type VMSecurebootCertificatesState string

const (
	// The VM&apos;s certificates do not need to be updated (including the case where Secure Boot does not apply to this VM, e.g. BIOS VM).
	VMSecurebootCertificatesStateOk VMSecurebootCertificatesState = "ok"
	// The Secure Boot certificates are due to expire or have already expired.
	VMSecurebootCertificatesStateUpdateAvailable VMSecurebootCertificatesState = "update_available"
	// An update of the certificates will be triggered whenever the VM boots. This includes VM.start, VM.reboot and a guest-triggered reboot.
	VMSecurebootCertificatesStateUpdateOnBoot VMSecurebootCertificatesState = "update_on_boot"
	// The value does not belong to this enumeration
	VMSecurebootCertificatesStateUnrecognized VMSecurebootCertificatesState = "unrecognized"
)

type VMPowerState string

const (
	// VM is offline and not using any resources
	VMPowerStateHalted VMPowerState = "Halted"
	// All resources have been allocated but the VM itself is paused and its vCPUs are not running
	VMPowerStatePaused VMPowerState = "Paused"
	// Running
	VMPowerStateRunning VMPowerState = "Running"
	// VM state has been saved to disk and it is nolonger running. Note that disks remain in-use while the VM is suspended.
	VMPowerStateSuspended VMPowerState = "Suspended"
	// The value does not belong to this enumeration
	VMPowerStateUnrecognized VMPowerState = "unrecognized"
)

type VMOperations string

const (
	// refers to the operation &quot;snapshot&quot;
	VMOperationsSnapshot VMOperations = "snapshot"
	// refers to the operation &quot;clone&quot;
	VMOperationsClone VMOperations = "clone"
	// refers to the operation &quot;copy&quot;
	VMOperationsCopy VMOperations = "copy"
	// refers to the operation &quot;create_template&quot;
	VMOperationsCreateTemplate VMOperations = "create_template"
	// refers to the operation &quot;revert&quot;
	VMOperationsRevert VMOperations = "revert"
	// refers to the operation &quot;checkpoint&quot;
	VMOperationsCheckpoint VMOperations = "checkpoint"
	// refers to the operation &quot;snapshot_with_quiesce&quot;
	VMOperationsSnapshotWithQuiesce VMOperations = "snapshot_with_quiesce"
	// refers to the operation &quot;provision&quot;
	VMOperationsProvision VMOperations = "provision"
	// refers to the operation &quot;start&quot;
	VMOperationsStart VMOperations = "start"
	// refers to the operation &quot;start_on&quot;
	VMOperationsStartOn VMOperations = "start_on"
	// refers to the operation &quot;pause&quot;
	VMOperationsPause VMOperations = "pause"
	// refers to the operation &quot;unpause&quot;
	VMOperationsUnpause VMOperations = "unpause"
	// refers to the operation &quot;clean_shutdown&quot;
	VMOperationsCleanShutdown VMOperations = "clean_shutdown"
	// refers to the operation &quot;clean_reboot&quot;
	VMOperationsCleanReboot VMOperations = "clean_reboot"
	// refers to the operation &quot;hard_shutdown&quot;
	VMOperationsHardShutdown VMOperations = "hard_shutdown"
	// refers to the operation &quot;power_state_reset&quot;
	VMOperationsPowerStateReset VMOperations = "power_state_reset"
	// refers to the operation &quot;hard_reboot&quot;
	VMOperationsHardReboot VMOperations = "hard_reboot"
	// refers to the operation &quot;suspend&quot;
	VMOperationsSuspend VMOperations = "suspend"
	// refers to the operation &quot;csvm&quot;
	VMOperationsCsvm VMOperations = "csvm"
	// refers to the operation &quot;resume&quot;
	VMOperationsResume VMOperations = "resume"
	// refers to the operation &quot;resume_on&quot;
	VMOperationsResumeOn VMOperations = "resume_on"
	// refers to the operation &quot;pool_migrate&quot;
	VMOperationsPoolMigrate VMOperations = "pool_migrate"
	// refers to the operation &quot;migrate_send&quot;
	VMOperationsMigrateSend VMOperations = "migrate_send"
	// refers to the operation &quot;get_boot_record&quot;
	VMOperationsGetBootRecord VMOperations = "get_boot_record"
	// refers to the operation &quot;send_sysrq&quot;
	VMOperationsSendSysrq VMOperations = "send_sysrq"
	// refers to the operation &quot;send_trigger&quot;
	VMOperationsSendTrigger VMOperations = "send_trigger"
	// refers to the operation &quot;query_services&quot;
	VMOperationsQueryServices VMOperations = "query_services"
	// refers to the operation &quot;shutdown&quot;
	VMOperationsShutdown VMOperations = "shutdown"
	// refers to the operation &quot;call_plugin&quot;
	VMOperationsCallPlugin VMOperations = "call_plugin"
	// Changing the memory settings
	VMOperationsChangingMemoryLive VMOperations = "changing_memory_live"
	// Waiting for the memory settings to change
	VMOperationsAwaitingMemoryLive VMOperations = "awaiting_memory_live"
	// Changing the memory dynamic range
	VMOperationsChangingDynamicRange VMOperations = "changing_dynamic_range"
	// Changing the memory static range
	VMOperationsChangingStaticRange VMOperations = "changing_static_range"
	// Changing the memory limits
	VMOperationsChangingMemoryLimits VMOperations = "changing_memory_limits"
	// Changing the shadow memory for a halted VM.
	VMOperationsChangingShadowMemory VMOperations = "changing_shadow_memory"
	// Changing the shadow memory for a running VM.
	VMOperationsChangingShadowMemoryLive VMOperations = "changing_shadow_memory_live"
	// Changing VCPU settings for a halted VM.
	VMOperationsChangingVCPUs VMOperations = "changing_VCPUs"
	// Changing VCPU settings for a running VM.
	VMOperationsChangingVCPUsLive VMOperations = "changing_VCPUs_live"
	// Changing NVRAM for a halted VM.
	VMOperationsChangingNVRAM VMOperations = "changing_NVRAM"
	//
	VMOperationsAssertOperationValid VMOperations = "assert_operation_valid"
	// Add, remove, query or list data sources
	VMOperationsDataSourceOp VMOperations = "data_source_op"
	//
	VMOperationsUpdateAllowedOperations VMOperations = "update_allowed_operations"
	// Turning this VM into a template
	VMOperationsMakeIntoTemplate VMOperations = "make_into_template"
	// importing a VM from a network stream
	VMOperationsImport VMOperations = "import"
	// exporting a VM to a network stream
	VMOperationsExport VMOperations = "export"
	// exporting VM metadata to a network stream
	VMOperationsMetadataExport VMOperations = "metadata_export"
	// Reverting the VM to a previous snapshotted state
	VMOperationsReverting VMOperations = "reverting"
	// refers to the act of uninstalling the VM
	VMOperationsDestroy VMOperations = "destroy"
	// Creating and adding a VTPM to this VM
	VMOperationsCreateVtpm VMOperations = "create_vtpm"
	// Performing a Windows sysprep on this VM
	VMOperationsSysprep VMOperations = "sysprep"
	// The value does not belong to this enumeration
	VMOperationsUnrecognized VMOperations = "unrecognized"
)

type VMApplianceOperation string

const (
	// Start
	VMApplianceOperationStart VMApplianceOperation = "start"
	// Clean shutdown
	VMApplianceOperationCleanShutdown VMApplianceOperation = "clean_shutdown"
	// Hard shutdown
	VMApplianceOperationHardShutdown VMApplianceOperation = "hard_shutdown"
	// Shutdown
	VMApplianceOperationShutdown VMApplianceOperation = "shutdown"
	// The value does not belong to this enumeration
	VMApplianceOperationUnrecognized VMApplianceOperation = "unrecognized"
)

type UpdateSyncFrequency string

const (
	// The update synchronizations happen every day
	UpdateSyncFrequencyDaily UpdateSyncFrequency = "daily"
	// The update synchronizations happen every week on the chosen day
	UpdateSyncFrequencyWeekly UpdateSyncFrequency = "weekly"
	// The value does not belong to this enumeration
	UpdateSyncFrequencyUnrecognized UpdateSyncFrequency = "unrecognized"
)

type UpdateGuidances string

const (
	// Indicates the updated host should reboot as soon as possible
	UpdateGuidancesRebootHost UpdateGuidances = "reboot_host"
	// Indicates the updated host should reboot as soon as possible since one or more livepatch(es) failed to be applied.
	UpdateGuidancesRebootHostOnLivepatchFailure UpdateGuidances = "reboot_host_on_livepatch_failure"
	// Indicates the updated host should reboot as soon as possible since one or more kernel livepatch(es) failed to be applied.
	UpdateGuidancesRebootHostOnKernelLivepatchFailure UpdateGuidances = "reboot_host_on_kernel_livepatch_failure"
	// Indicates the updated host should reboot as soon as possible since one or more xen livepatch(es) failed to be applied.
	UpdateGuidancesRebootHostOnXenLivepatchFailure UpdateGuidances = "reboot_host_on_xen_livepatch_failure"
	// Indicates the Toolstack running on the updated host should restart as soon as possible
	UpdateGuidancesRestartToolstack UpdateGuidances = "restart_toolstack"
	// Indicates the device model of a running VM should restart as soon as possible
	UpdateGuidancesRestartDeviceModel UpdateGuidances = "restart_device_model"
	// Indicates the VM should restart as soon as possible
	UpdateGuidancesRestartVM UpdateGuidances = "restart_vm"
	// The value does not belong to this enumeration
	UpdateGuidancesUnrecognized UpdateGuidances = "unrecognized"
)

type UpdateAfterApplyGuidance string

const (
	// This update requires HVM guests to be restarted once applied.
	UpdateAfterApplyGuidanceRestartHVM UpdateAfterApplyGuidance = "restartHVM"
	// This update requires PV guests to be restarted once applied.
	UpdateAfterApplyGuidanceRestartPV UpdateAfterApplyGuidance = "restartPV"
	// This update requires the host to be restarted once applied.
	UpdateAfterApplyGuidanceRestartHost UpdateAfterApplyGuidance = "restartHost"
	// This update requires XAPI to be restarted once applied.
	UpdateAfterApplyGuidanceRestartXAPI UpdateAfterApplyGuidance = "restartXAPI"
	// The value does not belong to this enumeration
	UpdateAfterApplyGuidanceUnrecognized UpdateAfterApplyGuidance = "unrecognized"
)

type TunnelProtocol string

const (
	// GRE protocol
	TunnelProtocolGre TunnelProtocol = "gre"
	// VxLAN Protocol
	TunnelProtocolVxlan TunnelProtocol = "vxlan"
	// The value does not belong to this enumeration
	TunnelProtocolUnrecognized TunnelProtocol = "unrecognized"
)

type TristateType string

const (
	// Known to be true
	TristateTypeYes TristateType = "yes"
	// Known to be false
	TristateTypeNo TristateType = "no"
	// Unknown or unspecified
	TristateTypeUnspecified TristateType = "unspecified"
	// The value does not belong to this enumeration
	TristateTypeUnrecognized TristateType = "unrecognized"
)

type TelemetryFrequency string

const (
	// Run telemetry task daily
	TelemetryFrequencyDaily TelemetryFrequency = "daily"
	// Run telemetry task weekly
	TelemetryFrequencyWeekly TelemetryFrequency = "weekly"
	// Run telemetry task monthly
	TelemetryFrequencyMonthly TelemetryFrequency = "monthly"
	// The value does not belong to this enumeration
	TelemetryFrequencyUnrecognized TelemetryFrequency = "unrecognized"
)

type TaskStatusType string

const (
	// task is in progress
	TaskStatusTypePending TaskStatusType = "pending"
	// task was completed successfully
	TaskStatusTypeSuccess TaskStatusType = "success"
	// task has failed
	TaskStatusTypeFailure TaskStatusType = "failure"
	// task is being cancelled
	TaskStatusTypeCancelling TaskStatusType = "cancelling"
	// task has been cancelled
	TaskStatusTypeCancelled TaskStatusType = "cancelled"
	// The value does not belong to this enumeration
	TaskStatusTypeUnrecognized TaskStatusType = "unrecognized"
)

type TaskAllowedOperations string

const (
	// refers to the operation &quot;cancel&quot;
	TaskAllowedOperationsCancel TaskAllowedOperations = "cancel"
	// refers to the operation &quot;destroy&quot;
	TaskAllowedOperationsDestroy TaskAllowedOperations = "destroy"
	// The value does not belong to this enumeration
	TaskAllowedOperationsUnrecognized TaskAllowedOperations = "unrecognized"
)

type StorageOperations string

const (
	// Scanning backends for new or deleted VDIs
	StorageOperationsScan StorageOperations = "scan"
	// Destroying the SR
	StorageOperationsDestroy StorageOperations = "destroy"
	// Forgetting about SR
	StorageOperationsForget StorageOperations = "forget"
	// Plugging a PBD into this SR
	StorageOperationsPlug StorageOperations = "plug"
	// Unplugging a PBD from this SR
	StorageOperationsUnplug StorageOperations = "unplug"
	// Refresh the fields on the SR
	StorageOperationsUpdate StorageOperations = "update"
	// Creating a new VDI
	StorageOperationsVdiCreate StorageOperations = "vdi_create"
	// Introducing a new VDI
	StorageOperationsVdiIntroduce StorageOperations = "vdi_introduce"
	// Destroying a VDI
	StorageOperationsVdiDestroy StorageOperations = "vdi_destroy"
	// Resizing a VDI
	StorageOperationsVdiResize StorageOperations = "vdi_resize"
	// Cloneing a VDI
	StorageOperationsVdiClone StorageOperations = "vdi_clone"
	// Snapshotting a VDI
	StorageOperationsVdiSnapshot StorageOperations = "vdi_snapshot"
	// Mirroring a VDI
	StorageOperationsVdiMirror StorageOperations = "vdi_mirror"
	// Enabling changed block tracking for a VDI
	StorageOperationsVdiEnableCbt StorageOperations = "vdi_enable_cbt"
	// Disabling changed block tracking for a VDI
	StorageOperationsVdiDisableCbt StorageOperations = "vdi_disable_cbt"
	// Deleting the data of the VDI
	StorageOperationsVdiDataDestroy StorageOperations = "vdi_data_destroy"
	// Exporting a bitmap that shows the changed blocks between two VDIs
	StorageOperationsVdiListChangedBlocks StorageOperations = "vdi_list_changed_blocks"
	// Setting the on_boot field of the VDI
	StorageOperationsVdiSetOnBoot StorageOperations = "vdi_set_on_boot"
	// Blocking other operations for a VDI
	StorageOperationsVdiBlocked StorageOperations = "vdi_blocked"
	// Copying the VDI
	StorageOperationsVdiCopy StorageOperations = "vdi_copy"
	// Forcefully unlocking the VDI
	StorageOperationsVdiForceUnlock StorageOperations = "vdi_force_unlock"
	// Forgetting about the VDI
	StorageOperationsVdiForget StorageOperations = "vdi_forget"
	// Generating the configuration of the VDI
	StorageOperationsVdiGenerateConfig StorageOperations = "vdi_generate_config"
	// Resizing the VDI online
	StorageOperationsVdiResizeOnline StorageOperations = "vdi_resize_online"
	// Refreshing the fields on the VDI
	StorageOperationsVdiUpdate StorageOperations = "vdi_update"
	// Reverting a VDI to the snapshot
	StorageOperationsVdiRevert StorageOperations = "vdi_revert"
	// Creating a PBD for this SR
	StorageOperationsPbdCreate StorageOperations = "pbd_create"
	// Destroying one of this SR&apos;s PBDs
	StorageOperationsPbdDestroy StorageOperations = "pbd_destroy"
	// The value does not belong to this enumeration
	StorageOperationsUnrecognized StorageOperations = "unrecognized"
)

type SriovConfigurationMode string

const (
	// Configure network sriov by sysfs, do not need reboot
	SriovConfigurationModeSysfs SriovConfigurationMode = "sysfs"
	// Configure network sriov by modprobe, need reboot
	SriovConfigurationModeModprobe SriovConfigurationMode = "modprobe"
	// Configure network sriov manually
	SriovConfigurationModeManual SriovConfigurationMode = "manual"
	// Unknown mode
	SriovConfigurationModeUnknown SriovConfigurationMode = "unknown"
	// The value does not belong to this enumeration
	SriovConfigurationModeUnrecognized SriovConfigurationMode = "unrecognized"
)

type SrHealth string

const (
	// Storage is fully available
	SrHealthHealthy SrHealth = "healthy"
	// Storage is busy recovering, for example, rebuilding mirrors.
	SrHealthRecovering SrHealth = "recovering"
	// Storage is unreachable but may be recoverable with admin intervention
	SrHealthUnreachable SrHealth = "unreachable"
	// Storage is unavailable, a host reboot will be required
	SrHealthUnavailable SrHealth = "unavailable"
	// The value does not belong to this enumeration
	SrHealthUnrecognized SrHealth = "unrecognized"
)

type SdnControllerProtocol string

const (
	// Active ssl connection
	SdnControllerProtocolSsl SdnControllerProtocol = "ssl"
	// Passive ssl connection
	SdnControllerProtocolPssl SdnControllerProtocol = "pssl"
	// The value does not belong to this enumeration
	SdnControllerProtocolUnrecognized SdnControllerProtocol = "unrecognized"
)

type PvsProxyStatus string

const (
	// The proxy is not currently running
	PvsProxyStatusStopped PvsProxyStatus = "stopped"
	// The proxy is setup but has not yet cached anything
	PvsProxyStatusInitialised PvsProxyStatus = "initialised"
	// The proxy is currently caching data
	PvsProxyStatusCaching PvsProxyStatus = "caching"
	// The PVS device is configured to use an incompatible write-cache mode
	PvsProxyStatusIncompatibleWriteCacheMode PvsProxyStatus = "incompatible_write_cache_mode"
	// The PVS protocol in use is not compatible with the PVS proxy
	PvsProxyStatusIncompatibleProtocolVersion PvsProxyStatus = "incompatible_protocol_version"
	// The value does not belong to this enumeration
	PvsProxyStatusUnrecognized PvsProxyStatus = "unrecognized"
)

type PrimaryAddressType string

const (
	// Primary address is the IPv4 address
	PrimaryAddressTypeIPv4 PrimaryAddressType = "IPv4"
	// Primary address is the IPv6 address
	PrimaryAddressTypeIPv6 PrimaryAddressType = "IPv6"
	// The value does not belong to this enumeration
	PrimaryAddressTypeUnrecognized PrimaryAddressType = "unrecognized"
)

type PoolGuestSecurebootReadiness string

const (
	// Pool is ready for SecureBoot, all auth files are present
	PoolGuestSecurebootReadinessReady PoolGuestSecurebootReadiness = "ready"
	// Pool is ready for SecureBoot, but there is no dbx auth file
	PoolGuestSecurebootReadinessReadyNoDbx PoolGuestSecurebootReadiness = "ready_no_dbx"
	// Pool is not ready for SecureBoot, mandatory auth files are missing
	PoolGuestSecurebootReadinessNotReady PoolGuestSecurebootReadiness = "not_ready"
	// The value does not belong to this enumeration
	PoolGuestSecurebootReadinessUnrecognized PoolGuestSecurebootReadiness = "unrecognized"
)

type PoolAllowedOperations string

const (
	// Indicates this pool is in the process of enabling HA
	PoolAllowedOperationsHaEnable PoolAllowedOperations = "ha_enable"
	// Indicates this pool is in the process of disabling HA
	PoolAllowedOperationsHaDisable PoolAllowedOperations = "ha_disable"
	// Indicates this pool is in the process of creating a cluster
	PoolAllowedOperationsClusterCreate PoolAllowedOperations = "cluster_create"
	// Indicates this pool is in the process of changing master
	PoolAllowedOperationsDesignateNewMaster PoolAllowedOperations = "designate_new_master"
	// Indicates this pool is in the process of configuring repositories
	PoolAllowedOperationsConfigureRepositories PoolAllowedOperations = "configure_repositories"
	// Indicates this pool is in the process of syncing updates
	PoolAllowedOperationsSyncUpdates PoolAllowedOperations = "sync_updates"
	// Indicates this pool is in the process of syncing bundle
	PoolAllowedOperationsSyncBundle PoolAllowedOperations = "sync_bundle"
	// Indicates this pool is in the process of getting updates
	PoolAllowedOperationsGetUpdates PoolAllowedOperations = "get_updates"
	// Indicates this pool is in the process of applying updates
	PoolAllowedOperationsApplyUpdates PoolAllowedOperations = "apply_updates"
	// Indicates this pool is in the process of enabling TLS verification
	PoolAllowedOperationsTLSVerificationEnable PoolAllowedOperations = "tls_verification_enable"
	// A certificate refresh and distribution is in progress
	PoolAllowedOperationsCertRefresh PoolAllowedOperations = "cert_refresh"
	// Indicates this pool is exchanging internal certificates with a new joiner
	PoolAllowedOperationsExchangeCertificatesOnJoin PoolAllowedOperations = "exchange_certificates_on_join"
	// Indicates this pool is exchanging ca certificates with a new joiner
	PoolAllowedOperationsExchangeCaCertificatesOnJoin PoolAllowedOperations = "exchange_ca_certificates_on_join"
	// Indicates the primary host is sending its certificates to another host
	PoolAllowedOperationsCopyPrimaryHostCerts PoolAllowedOperations = "copy_primary_host_certs"
	// Ejection of a host from the pool is under way
	PoolAllowedOperationsEject PoolAllowedOperations = "eject"
	// Indicates this pool is exchanging CRLs with a new joiner
	PoolAllowedOperationsExchangeCrlsOnJoin PoolAllowedOperations = "exchange_crls_on_join"
	// The value does not belong to this enumeration
	PoolAllowedOperationsUnrecognized PoolAllowedOperations = "unrecognized"
)

type PlacementPolicy string

const (
	// Anti-affinity placement policy
	PlacementPolicyAntiAffinity PlacementPolicy = "anti_affinity"
	// Default placement policy
	PlacementPolicyNormal PlacementPolicy = "normal"
	// The value does not belong to this enumeration
	PlacementPolicyUnrecognized PlacementPolicy = "unrecognized"
)

type PifIgmpStatus string

const (
	// IGMP Snooping is enabled in the corresponding backend bridge.&apos;
	PifIgmpStatusEnabled PifIgmpStatus = "enabled"
	// IGMP Snooping is disabled in the corresponding backend bridge.&apos;
	PifIgmpStatusDisabled PifIgmpStatus = "disabled"
	// IGMP snooping status is unknown. If this is a VLAN master, then please consult the underlying VLAN slave PIF.
	PifIgmpStatusUnknown PifIgmpStatus = "unknown"
	// The value does not belong to this enumeration
	PifIgmpStatusUnrecognized PifIgmpStatus = "unrecognized"
)

type PersistenceBackend string

const (
	// This VTPM is persisted in XAPI&apos;s DB
	PersistenceBackendXapi PersistenceBackend = "xapi"
	// The value does not belong to this enumeration
	PersistenceBackendUnrecognized PersistenceBackend = "unrecognized"
)

type PciDom0Access string

const (
	// dom0 can access this device as normal
	PciDom0AccessEnabled PciDom0Access = "enabled"
	// On host reboot dom0 will be blocked from accessing this device
	PciDom0AccessDisableOnReboot PciDom0Access = "disable_on_reboot"
	// dom0 cannot access this device
	PciDom0AccessDisabled PciDom0Access = "disabled"
	// On host reboot dom0 will be allowed to access this device
	PciDom0AccessEnableOnReboot PciDom0Access = "enable_on_reboot"
	// The value does not belong to this enumeration
	PciDom0AccessUnrecognized PciDom0Access = "unrecognized"
)

type Origin string

const (
	// The origin of the repository is a remote one
	OriginRemote Origin = "remote"
	// The origin of the repository is a local bundle file
	OriginBundle Origin = "bundle"
	// The origin of the repository is a remote pool
	OriginRemotePool Origin = "remote_pool"
	// The value does not belong to this enumeration
	OriginUnrecognized Origin = "unrecognized"
)

type OnSoftrebootBehavior string

const (
	// perform soft-reboot
	OnSoftrebootBehaviorSoftReboot OnSoftrebootBehavior = "soft_reboot"
	// destroy the VM state
	OnSoftrebootBehaviorDestroy OnSoftrebootBehavior = "destroy"
	// restart the VM
	OnSoftrebootBehaviorRestart OnSoftrebootBehavior = "restart"
	// leave the VM paused
	OnSoftrebootBehaviorPreserve OnSoftrebootBehavior = "preserve"
	// The value does not belong to this enumeration
	OnSoftrebootBehaviorUnrecognized OnSoftrebootBehavior = "unrecognized"
)

type OnNormalExit string

const (
	// destroy the VM state
	OnNormalExitDestroy OnNormalExit = "destroy"
	// restart the VM
	OnNormalExitRestart OnNormalExit = "restart"
	// The value does not belong to this enumeration
	OnNormalExitUnrecognized OnNormalExit = "unrecognized"
)

type OnCrashBehaviour string

const (
	// destroy the VM state
	OnCrashBehaviourDestroy OnCrashBehaviour = "destroy"
	// record a coredump and then destroy the VM state
	OnCrashBehaviourCoredumpAndDestroy OnCrashBehaviour = "coredump_and_destroy"
	// restart the VM
	OnCrashBehaviourRestart OnCrashBehaviour = "restart"
	// record a coredump and then restart the VM
	OnCrashBehaviourCoredumpAndRestart OnCrashBehaviour = "coredump_and_restart"
	// leave the crashed VM paused
	OnCrashBehaviourPreserve OnCrashBehaviour = "preserve"
	// rename the crashed VM and start a new copy
	OnCrashBehaviourRenameRestart OnCrashBehaviour = "rename_restart"
	// The value does not belong to this enumeration
	OnCrashBehaviourUnrecognized OnCrashBehaviour = "unrecognized"
)

type OnBoot string

const (
	// When a VM containing this VDI is started, the contents of the VDI are reset to the state they were in when this flag was last set.
	OnBootReset OnBoot = "reset"
	// Standard behaviour.
	OnBootPersist OnBoot = "persist"
	// The value does not belong to this enumeration
	OnBootUnrecognized OnBoot = "unrecognized"
)

type NetworkPurpose string

const (
	// Network Block Device service using TLS
	NetworkPurposeNbd NetworkPurpose = "nbd"
	// Network Block Device service without integrity or confidentiality: NOT RECOMMENDED
	NetworkPurposeInsecureNbd NetworkPurpose = "insecure_nbd"
	// The value does not belong to this enumeration
	NetworkPurposeUnrecognized NetworkPurpose = "unrecognized"
)

type NetworkOperations string

const (
	// Indicates this network is attaching to a VIF or PIF
	NetworkOperationsAttaching NetworkOperations = "attaching"
	// The value does not belong to this enumeration
	NetworkOperationsUnrecognized NetworkOperations = "unrecognized"
)

type NetworkDefaultLockingMode string

const (
	// Treat all VIFs on this network with locking_mode = &apos;default&apos; as if they have locking_mode = &apos;unlocked&apos;
	NetworkDefaultLockingModeUnlocked NetworkDefaultLockingMode = "unlocked"
	// Treat all VIFs on this network with locking_mode = &apos;default&apos; as if they have locking_mode = &apos;disabled&apos;
	NetworkDefaultLockingModeDisabled NetworkDefaultLockingMode = "disabled"
	// The value does not belong to this enumeration
	NetworkDefaultLockingModeUnrecognized NetworkDefaultLockingMode = "unrecognized"
)

type LivepatchStatus string

const (
	// An applicable live patch exists for every required component
	LivepatchStatusOkLivepatchComplete LivepatchStatus = "ok_livepatch_complete"
	// An applicable live patch exists but it is not sufficient
	LivepatchStatusOkLivepatchIncomplete LivepatchStatus = "ok_livepatch_incomplete"
	// There is no applicable live patch
	LivepatchStatusOk LivepatchStatus = "ok"
	// The value does not belong to this enumeration
	LivepatchStatusUnrecognized LivepatchStatus = "unrecognized"
)

type LatestSyncedUpdatesAppliedState string

const (
	// The host is up to date with the latest updates synced from remote CDN
	LatestSyncedUpdatesAppliedStateYes LatestSyncedUpdatesAppliedState = "yes"
	// The host is outdated with the latest updates synced from remote CDN
	LatestSyncedUpdatesAppliedStateNo LatestSyncedUpdatesAppliedState = "no"
	// If the host is up to date with the latest updates synced from remote CDN is unknown
	LatestSyncedUpdatesAppliedStateUnknown LatestSyncedUpdatesAppliedState = "unknown"
	// The value does not belong to this enumeration
	LatestSyncedUpdatesAppliedStateUnrecognized LatestSyncedUpdatesAppliedState = "unrecognized"
)

type Ipv6ConfigurationMode string

const (
	// Do not acquire an IPv6 address
	Ipv6ConfigurationModeNone Ipv6ConfigurationMode = "None"
	// Acquire an IPv6 address by DHCP
	Ipv6ConfigurationModeDHCP Ipv6ConfigurationMode = "DHCP"
	// Static IPv6 address configuration
	Ipv6ConfigurationModeStatic Ipv6ConfigurationMode = "Static"
	// Router assigned prefix delegation IPv6 allocation
	Ipv6ConfigurationModeAutoconf Ipv6ConfigurationMode = "Autoconf"
	// The value does not belong to this enumeration
	Ipv6ConfigurationModeUnrecognized Ipv6ConfigurationMode = "unrecognized"
)

type ImageFormatType string

const (
	// Plain disk image
	ImageFormatTypeRaw ImageFormatType = "raw"
	// Virtual Hard Disk
	ImageFormatTypeVhd ImageFormatType = "vhd"
	// Qemu Copy-On-Write version 2
	ImageFormatTypeQcow2 ImageFormatType = "qcow2"
	// The value does not belong to this enumeration
	ImageFormatTypeUnrecognized ImageFormatType = "unrecognized"
)

type IPConfigurationMode string

const (
	// Do not acquire an IP address
	IPConfigurationModeNone IPConfigurationMode = "None"
	// Acquire an IP address by DHCP
	IPConfigurationModeDHCP IPConfigurationMode = "DHCP"
	// Static IP address configuration
	IPConfigurationModeStatic IPConfigurationMode = "Static"
	// The value does not belong to this enumeration
	IPConfigurationModeUnrecognized IPConfigurationMode = "unrecognized"
)

type HostSchedGran string

const (
	// core scheduling
	HostSchedGranCore HostSchedGran = "core"
	// CPU scheduling
	HostSchedGranCPU HostSchedGran = "cpu"
	// socket scheduling
	HostSchedGranSocket HostSchedGran = "socket"
	// The value does not belong to this enumeration
	HostSchedGranUnrecognized HostSchedGran = "unrecognized"
)

type HostNumaAffinityPolicy string

const (
	// VMs are spread across all available NUMA nodes
	HostNumaAffinityPolicyAny HostNumaAffinityPolicy = "any"
	// VMs are placed on the smallest number of NUMA nodes that they fit using soft-pinning, but the policy doesn&apos;t guarantee a balanced placement, falling back to the &apos;any&apos; policy.
	HostNumaAffinityPolicyBestEffort HostNumaAffinityPolicy = "best_effort"
	// Use the NUMA affinity policy that is the default for the current version
	HostNumaAffinityPolicyDefaultPolicy HostNumaAffinityPolicy = "default_policy"
	// The value does not belong to this enumeration
	HostNumaAffinityPolicyUnrecognized HostNumaAffinityPolicy = "unrecognized"
)

type HostNtpMode string

const (
	// Using NTP servers assigned by DHCP to sync time
	HostNtpModeDHCP HostNtpMode = "DHCP"
	// Using custom NTP servers configured by user to sync time
	HostNtpModeCustom HostNtpMode = "Custom"
	// Using built-in NTP servers to sync time
	HostNtpModeFactory HostNtpMode = "Factory"
	// NTP is disabled on the host
	HostNtpModeDisabled HostNtpMode = "Disabled"
	// The value does not belong to this enumeration
	HostNtpModeUnrecognized HostNtpMode = "unrecognized"
)

type HostDisplay string

const (
	// This host is outputting its console to a physical display device
	HostDisplayEnabled HostDisplay = "enabled"
	// The host will stop outputting its console to a physical display device on next boot
	HostDisplayDisableOnReboot HostDisplay = "disable_on_reboot"
	// This host is not outputting its console to a physical display device
	HostDisplayDisabled HostDisplay = "disabled"
	// The host will start outputting its console to a physical display device on next boot
	HostDisplayEnableOnReboot HostDisplay = "enable_on_reboot"
	// The value does not belong to this enumeration
	HostDisplayUnrecognized HostDisplay = "unrecognized"
)

type HostAllowedOperations string

const (
	// Indicates this host is able to provision another VM
	HostAllowedOperationsProvision HostAllowedOperations = "provision"
	// Indicates this host is evacuating
	HostAllowedOperationsEvacuate HostAllowedOperations = "evacuate"
	// Indicates this host is in the process of shutting itself down
	HostAllowedOperationsShutdown HostAllowedOperations = "shutdown"
	// Indicates this host is in the process of rebooting
	HostAllowedOperationsReboot HostAllowedOperations = "reboot"
	// Indicates this host is in the process of being powered on
	HostAllowedOperationsPowerOn HostAllowedOperations = "power_on"
	// This host is starting a VM
	HostAllowedOperationsVMStart HostAllowedOperations = "vm_start"
	// This host is resuming a VM
	HostAllowedOperationsVMResume HostAllowedOperations = "vm_resume"
	// This host is the migration target of a VM
	HostAllowedOperationsVMMigrate HostAllowedOperations = "vm_migrate"
	// Indicates this host is being updated
	HostAllowedOperationsApplyUpdates HostAllowedOperations = "apply_updates"
	// Indicates this host is in the process of enabling
	HostAllowedOperationsEnable HostAllowedOperations = "enable"
	// The value does not belong to this enumeration
	HostAllowedOperationsUnrecognized HostAllowedOperations = "unrecognized"
)

type EventOperation string

const (
	// An object has been created
	EventOperationAdd EventOperation = "add"
	// An object has been deleted
	EventOperationDel EventOperation = "del"
	// An object has been modified
	EventOperationMod EventOperation = "mod"
	// The value does not belong to this enumeration
	EventOperationUnrecognized EventOperation = "unrecognized"
)

type DomainType string

const (
	// HVM; Fully Virtualised
	DomainTypeHvm DomainType = "hvm"
	// PV: Paravirtualised
	DomainTypePv DomainType = "pv"
	// PV inside a PVH container
	DomainTypePvInPvh DomainType = "pv_in_pvh"
	// PVH
	DomainTypePvh DomainType = "pvh"
	// Not specified or unknown domain type
	DomainTypeUnspecified DomainType = "unspecified"
	// The value does not belong to this enumeration
	DomainTypeUnrecognized DomainType = "unrecognized"
)

type ConsoleProtocol string

const (
	// VT100 terminal
	ConsoleProtocolVt100 ConsoleProtocol = "vt100"
	// Remote FrameBuffer protocol (as used in VNC)
	ConsoleProtocolRfb ConsoleProtocol = "rfb"
	// Remote Desktop Protocol
	ConsoleProtocolRdp ConsoleProtocol = "rdp"
	// The value does not belong to this enumeration
	ConsoleProtocolUnrecognized ConsoleProtocol = "unrecognized"
)

type ClusterOperation string

const (
	// adding a new member to the cluster
	ClusterOperationAdd ClusterOperation = "add"
	// removing a member from the cluster
	ClusterOperationRemove ClusterOperation = "remove"
	// enabling any cluster member
	ClusterOperationEnable ClusterOperation = "enable"
	// disabling any cluster member
	ClusterOperationDisable ClusterOperation = "disable"
	// completely destroying a cluster
	ClusterOperationDestroy ClusterOperation = "destroy"
	// The value does not belong to this enumeration
	ClusterOperationUnrecognized ClusterOperation = "unrecognized"
)

type ClusterHostOperation string

const (
	// enabling cluster membership on a particular host
	ClusterHostOperationEnable ClusterHostOperation = "enable"
	// disabling cluster membership on a particular host
	ClusterHostOperationDisable ClusterHostOperation = "disable"
	// completely destroying a cluster host
	ClusterHostOperationDestroy ClusterHostOperation = "destroy"
	// The value does not belong to this enumeration
	ClusterHostOperationUnrecognized ClusterHostOperation = "unrecognized"
)

type Cls string

const (
	// VM
	ClsVM Cls = "VM"
	// Host
	ClsHost Cls = "Host"
	// SR
	ClsSR Cls = "SR"
	// Pool
	ClsPool Cls = "Pool"
	// VMPP
	ClsVMPP Cls = "VMPP"
	// VMSS
	ClsVMSS Cls = "VMSS"
	// PVS_proxy
	ClsPVSProxy Cls = "PVS_proxy"
	// VDI
	ClsVDI Cls = "VDI"
	// Certificate
	ClsCertificate Cls = "Certificate"
	// The value does not belong to this enumeration
	ClsUnrecognized Cls = "unrecognized"
)

type CertificateType string

const (
	// Certificate that is trusted by the whole pool
	CertificateTypeCa CertificateType = "ca"
	// Certificate that identifies a single host to entities outside the pool
	CertificateTypeHost CertificateType = "host"
	// Certificate that identifies a single host to other pool members
	CertificateTypeHostInternal CertificateType = "host_internal"
	// Pinned leaf certificate that is trusted by the whole pool
	CertificateTypePinned CertificateType = "pinned"
	// The value does not belong to this enumeration
	CertificateTypeUnrecognized CertificateType = "unrecognized"
)

type CertificatePurpose string

const (
	// Trusted certificates that are for licensing purpose.
	CertificatePurposeLicensing CertificatePurpose = "licensing"
	// Trusted certificates that are for LDAPS connections to external authentication services.
	CertificatePurposeLdaps CertificatePurpose = "ldaps"
	// The value does not belong to this enumeration
	CertificatePurposeUnrecognized CertificatePurpose = "unrecognized"
)

type BondMode string

const (
	// Source-level balancing
	BondModeBalanceSlb BondMode = "balance-slb"
	// Active/passive bonding: only one NIC is carrying traffic
	BondModeActiveBackup BondMode = "active-backup"
	// Link aggregation control protocol
	BondModeLacp BondMode = "lacp"
	// The value does not belong to this enumeration
	BondModeUnrecognized BondMode = "unrecognized"
)

type AllocationAlgorithm string

const (
	// vGPUs of a given type are allocated evenly across supporting pGPUs.
	AllocationAlgorithmBreadthFirst AllocationAlgorithm = "breadth_first"
	// vGPUs of a given type are allocated on supporting pGPUs until they are full.
	AllocationAlgorithmDepthFirst AllocationAlgorithm = "depth_first"
	// The value does not belong to this enumeration
	AllocationAlgorithmUnrecognized AllocationAlgorithm = "unrecognized"
)

type AfterApplyGuidance string

const (
	// This patch requires HVM guests to be restarted once applied.
	AfterApplyGuidanceRestartHVM AfterApplyGuidance = "restartHVM"
	// This patch requires PV guests to be restarted once applied.
	AfterApplyGuidanceRestartPV AfterApplyGuidance = "restartPV"
	// This patch requires the host to be restarted once applied.
	AfterApplyGuidanceRestartHost AfterApplyGuidance = "restartHost"
	// This patch requires XAPI to be restarted once applied.
	AfterApplyGuidanceRestartXAPI AfterApplyGuidance = "restartXAPI"
	// The value does not belong to this enumeration
	AfterApplyGuidanceUnrecognized AfterApplyGuidance = "unrecognized"
)
