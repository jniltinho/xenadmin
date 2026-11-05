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

type VMPPRecord struct {
	// Unique identifier/object reference
	UUID string `json:"uuid,omitempty"`
	// a human-readable name
	NameLabel string `json:"name_label,omitempty"`
	// a notes field containing human-readable description
	NameDescription string `json:"name_description,omitempty"`
	// enable or disable this policy
	IsPolicyEnabled bool `json:"is_policy_enabled,omitempty"`
	// type of the backup sub-policy
	BackupType VmppBackupType `json:"backup_type,omitempty"`
	// maximum number of backups that should be stored at any time
	BackupRetentionValue int `json:"backup_retention_value,omitempty"`
	// frequency of the backup schedule
	BackupFrequency VmppBackupFrequency `json:"backup_frequency,omitempty"`
	// schedule of the backup containing &apos;hour&apos;, &apos;min&apos;, &apos;days&apos;. Date/time-related information is in Local Timezone
	BackupSchedule map[string]string `json:"backup_schedule,omitempty"`
	// true if this protection policy&apos;s backup is running
	IsBackupRunning bool `json:"is_backup_running,omitempty"`
	// time of the last backup
	BackupLastRunTime time.Time `json:"backup_last_run_time,omitempty"`
	// type of the archive target config
	ArchiveTargetType VmppArchiveTargetType `json:"archive_target_type,omitempty"`
	// configuration for the archive, including its &apos;location&apos;, &apos;username&apos;, &apos;password&apos;
	ArchiveTargetConfig map[string]string `json:"archive_target_config,omitempty"`
	// frequency of the archive schedule
	ArchiveFrequency VmppArchiveFrequency `json:"archive_frequency,omitempty"`
	// schedule of the archive containing &apos;hour&apos;, &apos;min&apos;, &apos;days&apos;. Date/time-related information is in Local Timezone
	ArchiveSchedule map[string]string `json:"archive_schedule,omitempty"`
	// true if this protection policy&apos;s archive is running
	IsArchiveRunning bool `json:"is_archive_running,omitempty"`
	// time of the last archive
	ArchiveLastRunTime time.Time `json:"archive_last_run_time,omitempty"`
	// all VMs attached to this protection policy
	VMs []VMRef `json:"VMs,omitempty"`
	// true if alarm is enabled for this policy
	IsAlarmEnabled bool `json:"is_alarm_enabled,omitempty"`
	// configuration for the alarm
	AlarmConfig map[string]string `json:"alarm_config,omitempty"`
	// recent alerts
	RecentAlerts []string `json:"recent_alerts,omitempty"`
}

type VMPPRef string

// VM Protection Policy
type vmpp struct{}

var VMPP vmpp

// GetAllRecords: Return a map of VMPP references to VMPP records for all VMPPs known to the system.
// Version: cowley
func (vmpp) GetAllRecords(session *Session) (retval map[VMPPRef]VMPPRecord, err error) {
	method := "VMPP.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMPPRefToVMPPRecordMap(method+" -> ", result)
	return
}

// GetAllRecords1: Return a map of VMPP references to VMPP records for all VMPPs known to the system.
// Version: cowley
func (vmpp) GetAllRecords1(session *Session) (retval map[VMPPRef]VMPPRecord, err error) {
	method := "VMPP.get_all_records"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMPPRefToVMPPRecordMap(method+" -> ", result)
	return
}

// GetAll: Return a list of all the VMPPs known to the system.
// Version: cowley
func (vmpp) GetAll(session *Session) (retval []VMPPRef, err error) {
	method := "VMPP.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMPPRefSet(method+" -> ", result)
	return
}

// GetAll1: Return a list of all the VMPPs known to the system.
// Version: cowley
func (vmpp) GetAll1(session *Session) (retval []VMPPRef, err error) {
	method := "VMPP.get_all"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMPPRefSet(method+" -> ", result)
	return
}

// SetArchiveLastRunTime:
// Version: cowley
func (vmpp) SetArchiveLastRunTime(session *Session, self VMPPRef, value time.Time) (err error) {
	method := "VMPP.set_archive_last_run_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeTime(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetArchiveLastRunTime3:
// Version: cowley
func (vmpp) SetArchiveLastRunTime3(session *Session, self VMPPRef, value time.Time) (err error) {
	method := "VMPP.set_archive_last_run_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeTime(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetBackupLastRunTime:
// Version: cowley
func (vmpp) SetBackupLastRunTime(session *Session, self VMPPRef, value time.Time) (err error) {
	method := "VMPP.set_backup_last_run_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeTime(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetBackupLastRunTime3:
// Version: cowley
func (vmpp) SetBackupLastRunTime3(session *Session, self VMPPRef, value time.Time) (err error) {
	method := "VMPP.set_backup_last_run_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeTime(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// RemoveFromAlarmConfig:
// Version: cowley
func (vmpp) RemoveFromAlarmConfig(session *Session, self VMPPRef, key string) (err error) {
	method := "VMPP.remove_from_alarm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromAlarmConfig3:
// Version: cowley
func (vmpp) RemoveFromAlarmConfig3(session *Session, self VMPPRef, key string) (err error) {
	method := "VMPP.remove_from_alarm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromArchiveSchedule:
// Version: cowley
func (vmpp) RemoveFromArchiveSchedule(session *Session, self VMPPRef, key string) (err error) {
	method := "VMPP.remove_from_archive_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromArchiveSchedule3:
// Version: cowley
func (vmpp) RemoveFromArchiveSchedule3(session *Session, self VMPPRef, key string) (err error) {
	method := "VMPP.remove_from_archive_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromArchiveTargetConfig:
// Version: cowley
func (vmpp) RemoveFromArchiveTargetConfig(session *Session, self VMPPRef, key string) (err error) {
	method := "VMPP.remove_from_archive_target_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromArchiveTargetConfig3:
// Version: cowley
func (vmpp) RemoveFromArchiveTargetConfig3(session *Session, self VMPPRef, key string) (err error) {
	method := "VMPP.remove_from_archive_target_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromBackupSchedule:
// Version: cowley
func (vmpp) RemoveFromBackupSchedule(session *Session, self VMPPRef, key string) (err error) {
	method := "VMPP.remove_from_backup_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// RemoveFromBackupSchedule3:
// Version: cowley
func (vmpp) RemoveFromBackupSchedule3(session *Session, self VMPPRef, key string) (err error) {
	method := "VMPP.remove_from_backup_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToAlarmConfig:
// Version: cowley
func (vmpp) AddToAlarmConfig(session *Session, self VMPPRef, key string, value string) (err error) {
	method := "VMPP.add_to_alarm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToAlarmConfig4:
// Version: cowley
func (vmpp) AddToAlarmConfig4(session *Session, self VMPPRef, key string, value string) (err error) {
	method := "VMPP.add_to_alarm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToArchiveSchedule:
// Version: cowley
func (vmpp) AddToArchiveSchedule(session *Session, self VMPPRef, key string, value string) (err error) {
	method := "VMPP.add_to_archive_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToArchiveSchedule4:
// Version: cowley
func (vmpp) AddToArchiveSchedule4(session *Session, self VMPPRef, key string, value string) (err error) {
	method := "VMPP.add_to_archive_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToArchiveTargetConfig:
// Version: cowley
func (vmpp) AddToArchiveTargetConfig(session *Session, self VMPPRef, key string, value string) (err error) {
	method := "VMPP.add_to_archive_target_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToArchiveTargetConfig4:
// Version: cowley
func (vmpp) AddToArchiveTargetConfig4(session *Session, self VMPPRef, key string, value string) (err error) {
	method := "VMPP.add_to_archive_target_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToBackupSchedule:
// Version: cowley
func (vmpp) AddToBackupSchedule(session *Session, self VMPPRef, key string, value string) (err error) {
	method := "VMPP.add_to_backup_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// AddToBackupSchedule4:
// Version: cowley
func (vmpp) AddToBackupSchedule4(session *Session, self VMPPRef, key string, value string) (err error) {
	method := "VMPP.add_to_backup_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetAlarmConfig:
// Version: cowley
func (vmpp) SetAlarmConfig(session *Session, self VMPPRef, value map[string]string) (err error) {
	method := "VMPP.set_alarm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetAlarmConfig3:
// Version: cowley
func (vmpp) SetAlarmConfig3(session *Session, self VMPPRef, value map[string]string) (err error) {
	method := "VMPP.set_alarm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetIsAlarmEnabled: Set the value of the is_alarm_enabled field
// Version: cowley
func (vmpp) SetIsAlarmEnabled(session *Session, self VMPPRef, value bool) (err error) {
	method := "VMPP.set_is_alarm_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetIsAlarmEnabled3: Set the value of the is_alarm_enabled field
// Version: cowley
func (vmpp) SetIsAlarmEnabled3(session *Session, self VMPPRef, value bool) (err error) {
	method := "VMPP.set_is_alarm_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetArchiveTargetConfig:
// Version: cowley
func (vmpp) SetArchiveTargetConfig(session *Session, self VMPPRef, value map[string]string) (err error) {
	method := "VMPP.set_archive_target_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetArchiveTargetConfig3:
// Version: cowley
func (vmpp) SetArchiveTargetConfig3(session *Session, self VMPPRef, value map[string]string) (err error) {
	method := "VMPP.set_archive_target_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetArchiveTargetType: Set the value of the archive_target_config_type field
// Version: cowley
func (vmpp) SetArchiveTargetType(session *Session, self VMPPRef, value VmppArchiveTargetType) (err error) {
	method := "VMPP.set_archive_target_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumVmppArchiveTargetType(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetArchiveTargetType3: Set the value of the archive_target_config_type field
// Version: cowley
func (vmpp) SetArchiveTargetType3(session *Session, self VMPPRef, value VmppArchiveTargetType) (err error) {
	method := "VMPP.set_archive_target_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumVmppArchiveTargetType(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetArchiveSchedule:
// Version: cowley
func (vmpp) SetArchiveSchedule(session *Session, self VMPPRef, value map[string]string) (err error) {
	method := "VMPP.set_archive_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetArchiveSchedule3:
// Version: cowley
func (vmpp) SetArchiveSchedule3(session *Session, self VMPPRef, value map[string]string) (err error) {
	method := "VMPP.set_archive_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetArchiveFrequency: Set the value of the archive_frequency field
// Version: cowley
func (vmpp) SetArchiveFrequency(session *Session, self VMPPRef, value VmppArchiveFrequency) (err error) {
	method := "VMPP.set_archive_frequency"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumVmppArchiveFrequency(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetArchiveFrequency3: Set the value of the archive_frequency field
// Version: cowley
func (vmpp) SetArchiveFrequency3(session *Session, self VMPPRef, value VmppArchiveFrequency) (err error) {
	method := "VMPP.set_archive_frequency"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumVmppArchiveFrequency(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetBackupSchedule:
// Version: cowley
func (vmpp) SetBackupSchedule(session *Session, self VMPPRef, value map[string]string) (err error) {
	method := "VMPP.set_backup_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetBackupSchedule3:
// Version: cowley
func (vmpp) SetBackupSchedule3(session *Session, self VMPPRef, value map[string]string) (err error) {
	method := "VMPP.set_backup_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetBackupFrequency: Set the value of the backup_frequency field
// Version: cowley
func (vmpp) SetBackupFrequency(session *Session, self VMPPRef, value VmppBackupFrequency) (err error) {
	method := "VMPP.set_backup_frequency"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumVmppBackupFrequency(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetBackupFrequency3: Set the value of the backup_frequency field
// Version: cowley
func (vmpp) SetBackupFrequency3(session *Session, self VMPPRef, value VmppBackupFrequency) (err error) {
	method := "VMPP.set_backup_frequency"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumVmppBackupFrequency(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetBackupRetentionValue:
// Version: cowley
func (vmpp) SetBackupRetentionValue(session *Session, self VMPPRef, value int) (err error) {
	method := "VMPP.set_backup_retention_value"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetBackupRetentionValue3:
// Version: cowley
func (vmpp) SetBackupRetentionValue3(session *Session, self VMPPRef, value int) (err error) {
	method := "VMPP.set_backup_retention_value"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAlerts: This call fetches a history of alerts for a given protection policy
// Version: cowley
func (vmpp) GetAlerts(session *Session, vmpp VMPPRef, hoursFromNow int) (retval []string, err error) {
	method := "VMPP.get_alerts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmppArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "vmpp"), vmpp)
	if err != nil {
		return
	}
	hoursFromNowArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "hours_from_now"), hoursFromNow)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmppArg, hoursFromNowArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// GetAlerts3: This call fetches a history of alerts for a given protection policy
// Version: cowley
func (vmpp) GetAlerts3(session *Session, vmpp VMPPRef, hoursFromNow int) (retval []string, err error) {
	method := "VMPP.get_alerts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmppArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "vmpp"), vmpp)
	if err != nil {
		return
	}
	hoursFromNowArg, err := serializeInt(fmt.Sprintf("%s(%s)", method, "hours_from_now"), hoursFromNow)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmppArg, hoursFromNowArg)
	if err != nil {
		return
	}
	retval, err = deserializeStringSet(method+" -> ", result)
	return
}

// ArchiveNow: This call archives the snapshot provided as a parameter
// Version: cowley
func (vmpp) ArchiveNow(session *Session, snapshot VMRef) (retval string, err error) {
	method := "VMPP.archive_now"
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
	retval, err = deserializeString(method+" -> ", result)
	return
}

// ArchiveNow2: This call archives the snapshot provided as a parameter
// Version: cowley
func (vmpp) ArchiveNow2(session *Session, snapshot VMRef) (retval string, err error) {
	method := "VMPP.archive_now"
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
	retval, err = deserializeString(method+" -> ", result)
	return
}

// ProtectNow: This call executes the protection policy immediately
// Version: cowley
func (vmpp) ProtectNow(session *Session, vmpp VMPPRef) (retval string, err error) {
	method := "VMPP.protect_now"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmppArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "vmpp"), vmpp)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmppArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// ProtectNow2: This call executes the protection policy immediately
// Version: cowley
func (vmpp) ProtectNow2(session *Session, vmpp VMPPRef) (retval string, err error) {
	method := "VMPP.protect_now"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	vmppArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "vmpp"), vmpp)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, vmppArg)
	if err != nil {
		return
	}
	retval, err = deserializeString(method+" -> ", result)
	return
}

// SetBackupType: Set the backup_type field of the given VMPP.
// Version: cowley
func (vmpp) SetBackupType(session *Session, self VMPPRef, value VmppBackupType) (err error) {
	method := "VMPP.set_backup_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumVmppBackupType(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetBackupType3: Set the backup_type field of the given VMPP.
// Version: cowley
func (vmpp) SetBackupType3(session *Session, self VMPPRef, value VmppBackupType) (err error) {
	method := "VMPP.set_backup_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	valueArg, err := serializeEnumVmppBackupType(fmt.Sprintf("%s(%s)", method, "value"), value)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg, valueArg)
	return
}

// SetIsPolicyEnabled: Set the is_policy_enabled field of the given VMPP.
// Version: cowley
func (vmpp) SetIsPolicyEnabled(session *Session, self VMPPRef, value bool) (err error) {
	method := "VMPP.set_is_policy_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetIsPolicyEnabled3: Set the is_policy_enabled field of the given VMPP.
// Version: cowley
func (vmpp) SetIsPolicyEnabled3(session *Session, self VMPPRef, value bool) (err error) {
	method := "VMPP.set_is_policy_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameDescription: Set the name/description field of the given VMPP.
// Version: cowley
func (vmpp) SetNameDescription(session *Session, self VMPPRef, value string) (err error) {
	method := "VMPP.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameDescription3: Set the name/description field of the given VMPP.
// Version: cowley
func (vmpp) SetNameDescription3(session *Session, self VMPPRef, value string) (err error) {
	method := "VMPP.set_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameLabel: Set the name/label field of the given VMPP.
// Version: cowley
func (vmpp) SetNameLabel(session *Session, self VMPPRef, value string) (err error) {
	method := "VMPP.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// SetNameLabel3: Set the name/label field of the given VMPP.
// Version: cowley
func (vmpp) SetNameLabel3(session *Session, self VMPPRef, value string) (err error) {
	method := "VMPP.set_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRecentAlerts: Get the recent_alerts field of the given VMPP.
// Version: cowley
func (vmpp) GetRecentAlerts(session *Session, self VMPPRef) (retval []string, err error) {
	method := "VMPP.get_recent_alerts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetRecentAlerts2: Get the recent_alerts field of the given VMPP.
// Version: cowley
func (vmpp) GetRecentAlerts2(session *Session, self VMPPRef) (retval []string, err error) {
	method := "VMPP.get_recent_alerts"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAlarmConfig: Get the alarm_config field of the given VMPP.
// Version: cowley
func (vmpp) GetAlarmConfig(session *Session, self VMPPRef) (retval map[string]string, err error) {
	method := "VMPP.get_alarm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetAlarmConfig2: Get the alarm_config field of the given VMPP.
// Version: cowley
func (vmpp) GetAlarmConfig2(session *Session, self VMPPRef) (retval map[string]string, err error) {
	method := "VMPP.get_alarm_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsAlarmEnabled: Get the is_alarm_enabled field of the given VMPP.
// Version: cowley
func (vmpp) GetIsAlarmEnabled(session *Session, self VMPPRef) (retval bool, err error) {
	method := "VMPP.get_is_alarm_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsAlarmEnabled2: Get the is_alarm_enabled field of the given VMPP.
// Version: cowley
func (vmpp) GetIsAlarmEnabled2(session *Session, self VMPPRef) (retval bool, err error) {
	method := "VMPP.get_is_alarm_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVMs: Get the VMs field of the given VMPP.
// Version: cowley
func (vmpp) GetVMs(session *Session, self VMPPRef) (retval []VMRef, err error) {
	method := "VMPP.get_VMs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetVMs2: Get the VMs field of the given VMPP.
// Version: cowley
func (vmpp) GetVMs2(session *Session, self VMPPRef) (retval []VMRef, err error) {
	method := "VMPP.get_VMs"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetArchiveLastRunTime: Get the archive_last_run_time field of the given VMPP.
// Version: cowley
func (vmpp) GetArchiveLastRunTime(session *Session, self VMPPRef) (retval time.Time, err error) {
	method := "VMPP.get_archive_last_run_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetArchiveLastRunTime2: Get the archive_last_run_time field of the given VMPP.
// Version: cowley
func (vmpp) GetArchiveLastRunTime2(session *Session, self VMPPRef) (retval time.Time, err error) {
	method := "VMPP.get_archive_last_run_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsArchiveRunning: Get the is_archive_running field of the given VMPP.
// Version: cowley
func (vmpp) GetIsArchiveRunning(session *Session, self VMPPRef) (retval bool, err error) {
	method := "VMPP.get_is_archive_running"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsArchiveRunning2: Get the is_archive_running field of the given VMPP.
// Version: cowley
func (vmpp) GetIsArchiveRunning2(session *Session, self VMPPRef) (retval bool, err error) {
	method := "VMPP.get_is_archive_running"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetArchiveSchedule: Get the archive_schedule field of the given VMPP.
// Version: cowley
func (vmpp) GetArchiveSchedule(session *Session, self VMPPRef) (retval map[string]string, err error) {
	method := "VMPP.get_archive_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetArchiveSchedule2: Get the archive_schedule field of the given VMPP.
// Version: cowley
func (vmpp) GetArchiveSchedule2(session *Session, self VMPPRef) (retval map[string]string, err error) {
	method := "VMPP.get_archive_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetArchiveFrequency: Get the archive_frequency field of the given VMPP.
// Version: cowley
func (vmpp) GetArchiveFrequency(session *Session, self VMPPRef) (retval VmppArchiveFrequency, err error) {
	method := "VMPP.get_archive_frequency"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVmppArchiveFrequency(method+" -> ", result)
	return
}

// GetArchiveFrequency2: Get the archive_frequency field of the given VMPP.
// Version: cowley
func (vmpp) GetArchiveFrequency2(session *Session, self VMPPRef) (retval VmppArchiveFrequency, err error) {
	method := "VMPP.get_archive_frequency"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVmppArchiveFrequency(method+" -> ", result)
	return
}

// GetArchiveTargetConfig: Get the archive_target_config field of the given VMPP.
// Version: cowley
func (vmpp) GetArchiveTargetConfig(session *Session, self VMPPRef) (retval map[string]string, err error) {
	method := "VMPP.get_archive_target_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetArchiveTargetConfig2: Get the archive_target_config field of the given VMPP.
// Version: cowley
func (vmpp) GetArchiveTargetConfig2(session *Session, self VMPPRef) (retval map[string]string, err error) {
	method := "VMPP.get_archive_target_config"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetArchiveTargetType: Get the archive_target_type field of the given VMPP.
// Version: cowley
func (vmpp) GetArchiveTargetType(session *Session, self VMPPRef) (retval VmppArchiveTargetType, err error) {
	method := "VMPP.get_archive_target_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVmppArchiveTargetType(method+" -> ", result)
	return
}

// GetArchiveTargetType2: Get the archive_target_type field of the given VMPP.
// Version: cowley
func (vmpp) GetArchiveTargetType2(session *Session, self VMPPRef) (retval VmppArchiveTargetType, err error) {
	method := "VMPP.get_archive_target_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVmppArchiveTargetType(method+" -> ", result)
	return
}

// GetBackupLastRunTime: Get the backup_last_run_time field of the given VMPP.
// Version: cowley
func (vmpp) GetBackupLastRunTime(session *Session, self VMPPRef) (retval time.Time, err error) {
	method := "VMPP.get_backup_last_run_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBackupLastRunTime2: Get the backup_last_run_time field of the given VMPP.
// Version: cowley
func (vmpp) GetBackupLastRunTime2(session *Session, self VMPPRef) (retval time.Time, err error) {
	method := "VMPP.get_backup_last_run_time"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsBackupRunning: Get the is_backup_running field of the given VMPP.
// Version: cowley
func (vmpp) GetIsBackupRunning(session *Session, self VMPPRef) (retval bool, err error) {
	method := "VMPP.get_is_backup_running"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsBackupRunning2: Get the is_backup_running field of the given VMPP.
// Version: cowley
func (vmpp) GetIsBackupRunning2(session *Session, self VMPPRef) (retval bool, err error) {
	method := "VMPP.get_is_backup_running"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBackupSchedule: Get the backup_schedule field of the given VMPP.
// Version: cowley
func (vmpp) GetBackupSchedule(session *Session, self VMPPRef) (retval map[string]string, err error) {
	method := "VMPP.get_backup_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBackupSchedule2: Get the backup_schedule field of the given VMPP.
// Version: cowley
func (vmpp) GetBackupSchedule2(session *Session, self VMPPRef) (retval map[string]string, err error) {
	method := "VMPP.get_backup_schedule"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBackupFrequency: Get the backup_frequency field of the given VMPP.
// Version: cowley
func (vmpp) GetBackupFrequency(session *Session, self VMPPRef) (retval VmppBackupFrequency, err error) {
	method := "VMPP.get_backup_frequency"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVmppBackupFrequency(method+" -> ", result)
	return
}

// GetBackupFrequency2: Get the backup_frequency field of the given VMPP.
// Version: cowley
func (vmpp) GetBackupFrequency2(session *Session, self VMPPRef) (retval VmppBackupFrequency, err error) {
	method := "VMPP.get_backup_frequency"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVmppBackupFrequency(method+" -> ", result)
	return
}

// GetBackupRetentionValue: Get the backup_retention_value field of the given VMPP.
// Version: cowley
func (vmpp) GetBackupRetentionValue(session *Session, self VMPPRef) (retval int, err error) {
	method := "VMPP.get_backup_retention_value"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBackupRetentionValue2: Get the backup_retention_value field of the given VMPP.
// Version: cowley
func (vmpp) GetBackupRetentionValue2(session *Session, self VMPPRef) (retval int, err error) {
	method := "VMPP.get_backup_retention_value"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetBackupType: Get the backup_type field of the given VMPP.
// Version: cowley
func (vmpp) GetBackupType(session *Session, self VMPPRef) (retval VmppBackupType, err error) {
	method := "VMPP.get_backup_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVmppBackupType(method+" -> ", result)
	return
}

// GetBackupType2: Get the backup_type field of the given VMPP.
// Version: cowley
func (vmpp) GetBackupType2(session *Session, self VMPPRef) (retval VmppBackupType, err error) {
	method := "VMPP.get_backup_type"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeEnumVmppBackupType(method+" -> ", result)
	return
}

// GetIsPolicyEnabled: Get the is_policy_enabled field of the given VMPP.
// Version: cowley
func (vmpp) GetIsPolicyEnabled(session *Session, self VMPPRef) (retval bool, err error) {
	method := "VMPP.get_is_policy_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetIsPolicyEnabled2: Get the is_policy_enabled field of the given VMPP.
// Version: cowley
func (vmpp) GetIsPolicyEnabled2(session *Session, self VMPPRef) (retval bool, err error) {
	method := "VMPP.get_is_policy_enabled"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription: Get the name/description field of the given VMPP.
// Version: cowley
func (vmpp) GetNameDescription(session *Session, self VMPPRef) (retval string, err error) {
	method := "VMPP.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameDescription2: Get the name/description field of the given VMPP.
// Version: cowley
func (vmpp) GetNameDescription2(session *Session, self VMPPRef) (retval string, err error) {
	method := "VMPP.get_name_description"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel: Get the name/label field of the given VMPP.
// Version: cowley
func (vmpp) GetNameLabel(session *Session, self VMPPRef) (retval string, err error) {
	method := "VMPP.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetNameLabel2: Get the name/label field of the given VMPP.
// Version: cowley
func (vmpp) GetNameLabel2(session *Session, self VMPPRef) (retval string, err error) {
	method := "VMPP.get_name_label"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID: Get the uuid field of the given VMPP.
// Version: cowley
func (vmpp) GetUUID(session *Session, self VMPPRef) (retval string, err error) {
	method := "VMPP.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetUUID2: Get the uuid field of the given VMPP.
// Version: cowley
func (vmpp) GetUUID2(session *Session, self VMPPRef) (retval string, err error) {
	method := "VMPP.get_uuid"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// GetByNameLabel: Get all the VMPP instances with the given label.
// Version: cowley
func (vmpp) GetByNameLabel(session *Session, label string) (retval []VMPPRef, err error) {
	method := "VMPP.get_by_name_label"
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
	retval, err = deserializeVMPPRefSet(method+" -> ", result)
	return
}

// GetByNameLabel2: Get all the VMPP instances with the given label.
// Version: cowley
func (vmpp) GetByNameLabel2(session *Session, label string) (retval []VMPPRef, err error) {
	method := "VMPP.get_by_name_label"
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
	retval, err = deserializeVMPPRefSet(method+" -> ", result)
	return
}

// Destroy: Destroy the specified VMPP instance.
// Version: cowley
func (vmpp) Destroy(session *Session, self VMPPRef) (err error) {
	method := "VMPP.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy: Destroy the specified VMPP instance.
// Version: cowley
func (vmpp) AsyncDestroy(session *Session, self VMPPRef) (retval TaskRef, err error) {
	method := "Async.VMPP.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Destroy2: Destroy the specified VMPP instance.
// Version: cowley
func (vmpp) Destroy2(session *Session, self VMPPRef) (err error) {
	method := "VMPP.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	_, err = session.client.sendCall(method, sessionIDArg, selfArg)
	return
}

// AsyncDestroy2: Destroy the specified VMPP instance.
// Version: cowley
func (vmpp) AsyncDestroy2(session *Session, self VMPPRef) (retval TaskRef, err error) {
	method := "Async.VMPP.destroy"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
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

// Create: Create a new VMPP instance, and return its handle. The constructor args are: name_label, name_description, is_policy_enabled, backup_type, backup_retention_value, backup_frequency, backup_schedule, archive_target_type, archive_target_config, archive_frequency, archive_schedule, is_alarm_enabled, alarm_config (* = non-optional).
// Version: cowley
func (vmpp) Create(session *Session, args VMPPRecord) (retval VMPPRef, err error) {
	method := "VMPP.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVMPPRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMPPRef(method+" -> ", result)
	return
}

// AsyncCreate: Create a new VMPP instance, and return its handle. The constructor args are: name_label, name_description, is_policy_enabled, backup_type, backup_retention_value, backup_frequency, backup_schedule, archive_target_type, archive_target_config, archive_frequency, archive_schedule, is_alarm_enabled, alarm_config (* = non-optional).
// Version: cowley
func (vmpp) AsyncCreate(session *Session, args VMPPRecord) (retval TaskRef, err error) {
	method := "Async.VMPP.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVMPPRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
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

// Create2: Create a new VMPP instance, and return its handle. The constructor args are: name_label, name_description, is_policy_enabled, backup_type, backup_retention_value, backup_frequency, backup_schedule, archive_target_type, archive_target_config, archive_frequency, archive_schedule, is_alarm_enabled, alarm_config (* = non-optional).
// Version: cowley
func (vmpp) Create2(session *Session, args VMPPRecord) (retval VMPPRef, err error) {
	method := "VMPP.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVMPPRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, argsArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMPPRef(method+" -> ", result)
	return
}

// AsyncCreate2: Create a new VMPP instance, and return its handle. The constructor args are: name_label, name_description, is_policy_enabled, backup_type, backup_retention_value, backup_frequency, backup_schedule, archive_target_type, archive_target_config, archive_frequency, archive_schedule, is_alarm_enabled, alarm_config (* = non-optional).
// Version: cowley
func (vmpp) AsyncCreate2(session *Session, args VMPPRecord) (retval TaskRef, err error) {
	method := "Async.VMPP.create"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	argsArg, err := serializeVMPPRecord(fmt.Sprintf("%s(%s)", method, "args"), args)
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

// GetByUUID: Get a reference to the VMPP instance with the specified UUID.
// Version: cowley
func (vmpp) GetByUUID(session *Session, uuid string) (retval VMPPRef, err error) {
	method := "VMPP.get_by_uuid"
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
	retval, err = deserializeVMPPRef(method+" -> ", result)
	return
}

// GetByUUID2: Get a reference to the VMPP instance with the specified UUID.
// Version: cowley
func (vmpp) GetByUUID2(session *Session, uuid string) (retval VMPPRef, err error) {
	method := "VMPP.get_by_uuid"
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
	retval, err = deserializeVMPPRef(method+" -> ", result)
	return
}

// GetRecord: Get a record containing the current state of the given VMPP.
// Version: cowley
func (vmpp) GetRecord(session *Session, self VMPPRef) (retval VMPPRecord, err error) {
	method := "VMPP.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMPPRecord(method+" -> ", result)
	return
}

// GetRecord2: Get a record containing the current state of the given VMPP.
// Version: cowley
func (vmpp) GetRecord2(session *Session, self VMPPRef) (retval VMPPRecord, err error) {
	method := "VMPP.get_record"
	sessionIDArg, err := serializeSessionRef(fmt.Sprintf("%s(%s)", method, "session_id"), session.ref)
	if err != nil {
		return
	}
	selfArg, err := serializeVMPPRef(fmt.Sprintf("%s(%s)", method, "self"), self)
	if err != nil {
		return
	}
	result, err := session.client.sendCall(method, sessionIDArg, selfArg)
	if err != nil {
		return
	}
	retval, err = deserializeVMPPRecord(method+" -> ", result)
	return
}
