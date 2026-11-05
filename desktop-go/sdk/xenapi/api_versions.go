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
	"errors"
	"fmt"
)

type APIVersion int

const (
	// XenServer 4.0 (rio)
	APIVersion1_1 APIVersion = iota + 1
	// XenServer 4.1 (miami)
	APIVersion1_2
	// XenServer 5.0 (orlando)
	APIVersion1_3
	// Unreleased ()
	APIVersion1_4
	// XenServer 5.0 update 3 ()
	APIVersion1_5
	// XenServer 5.5 (george)
	APIVersion1_6
	// XenServer 5.6 (midnight-ride)
	APIVersion1_7
	// XenServer 5.6 FP1 (cowley)
	APIVersion1_8
	// XenServer 6.0 (boston)
	APIVersion1_9
	// XenServer 6.1 (tampa)
	APIVersion1_10
	// XenServer 6.2 (clearwater)
	APIVersion2_0
	// XenServer 6.2 SP1 (vgpu-productisation)
	APIVersion2_1
	// XenServer 6.2 SP1 Hotfix 4 (clearwater-felton)
	APIVersion2_2
	// XenServer 6.5 (creedence)
	APIVersion2_3
	// XenServer 6.5 SP1 (cream)
	APIVersion2_4
	// XenServer 7.0 (dundee)
	APIVersion2_5
	// XenServer 7.1 (ely)
	APIVersion2_6
	// XenServer 7.2 (falcon)
	APIVersion2_7
	// XenServer 7.3 (inverness)
	APIVersion2_8
	// XenServer 7.4 (jura)
	APIVersion2_9
	// XenServer 7.5 (kolkata)
	APIVersion2_10
	// XenServer 7.6 (lima)
	APIVersion2_11
	// Citrix Hypervisor 8.0 (naples)
	APIVersion2_12
	// Unreleased (oslo)
	APIVersion2_13
	// Citrix Hypervisor 8.1 (quebec)
	APIVersion2_14
	// Citrix Hypervisor 8.2 (stockholm)
	APIVersion2_15
	// XenServer 8 Preview (nile-preview)
	APIVersion2_20
	// XenServer 8 (nile)
	APIVersion2_21
	// XenServer 9 (orinoco)
	APIVersion2_23
	APIVersionLatest  APIVersion = 29
	APIVersionUnknown APIVersion = 99
)

func (v APIVersion) String() string {
	switch v {
	case APIVersion1_1:
		return "1.1"
	case APIVersion1_2:
		return "1.2"
	case APIVersion1_3:
		return "1.3"
	case APIVersion1_4:
		return "1.4"
	case APIVersion1_5:
		return "1.5"
	case APIVersion1_6:
		return "1.6"
	case APIVersion1_7:
		return "1.7"
	case APIVersion1_8:
		return "1.8"
	case APIVersion1_9:
		return "1.9"
	case APIVersion1_10:
		return "1.10"
	case APIVersion2_0:
		return "2.0"
	case APIVersion2_1:
		return "2.1"
	case APIVersion2_2:
		return "2.2"
	case APIVersion2_3:
		return "2.3"
	case APIVersion2_4:
		return "2.4"
	case APIVersion2_5:
		return "2.5"
	case APIVersion2_6:
		return "2.6"
	case APIVersion2_7:
		return "2.7"
	case APIVersion2_8:
		return "2.8"
	case APIVersion2_9:
		return "2.9"
	case APIVersion2_10:
		return "2.10"
	case APIVersion2_11:
		return "2.11"
	case APIVersion2_12:
		return "2.12"
	case APIVersion2_13:
		return "2.13"
	case APIVersion2_14:
		return "2.14"
	case APIVersion2_15:
		return "2.15"
	case APIVersion2_20:
		return "2.20"
	case APIVersion2_21:
		return "2.21"
	case APIVersion2_23:
		return "2.23"
	case APIVersionUnknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

var APIVersionMap = map[string]APIVersion{
	//
	"APIVersion1_1": APIVersion1_1,
	//
	"APIVersion1_2": APIVersion1_2,
	//
	"APIVersion1_3": APIVersion1_3,
	//
	"APIVersion1_4": APIVersion1_4,
	//
	"APIVersion1_5": APIVersion1_5,
	//
	"APIVersion1_6": APIVersion1_6,
	//
	"APIVersion1_7": APIVersion1_7,
	//
	"APIVersion1_8": APIVersion1_8,
	//
	"APIVersion1_9": APIVersion1_9,
	//
	"APIVersion1_10": APIVersion1_10,
	//
	"APIVersion2_0": APIVersion2_0,
	//
	"APIVersion2_1": APIVersion2_1,
	//
	"APIVersion2_2": APIVersion2_2,
	//
	"APIVersion2_3": APIVersion2_3,
	//
	"APIVersion2_4": APIVersion2_4,
	//
	"APIVersion2_5": APIVersion2_5,
	//
	"APIVersion2_6": APIVersion2_6,
	//
	"APIVersion2_7": APIVersion2_7,
	//
	"APIVersion2_8": APIVersion2_8,
	//
	"APIVersion2_9": APIVersion2_9,
	//
	"APIVersion2_10": APIVersion2_10,
	//
	"APIVersion2_11": APIVersion2_11,
	//
	"APIVersion2_12": APIVersion2_12,
	//
	"APIVersion2_13": APIVersion2_13,
	//
	"APIVersion2_14": APIVersion2_14,
	//
	"APIVersion2_15": APIVersion2_15,
	//
	"APIVersion2_20": APIVersion2_20,
	//
	"APIVersion2_21": APIVersion2_21,
	//
	"APIVersion2_23": APIVersion2_23,
	//
	"APIVersionLatest": APIVersionLatest,
	//
	"APIVersionUnknown": APIVersionUnknown,
}

func GetAPIVersion(major int, minor int) APIVersion {
	versionName := fmt.Sprintf("APIVersion%d_%d", major, minor)
	apiVersion, ok := APIVersionMap[versionName]
	if !ok {
		apiVersion = APIVersionUnknown
	}

	return apiVersion
}

func getPoolMaster(session *Session) (HostRef, error) {
	var master HostRef
	poolRefs, err := Pool.GetAll(session)
	if err != nil {
		return master, err
	}
	if len(poolRefs) > 0 {
		poolRecord, err := Pool.GetRecord(session, poolRefs[0])
		if err != nil {
			return master, err
		}
		return poolRecord.Master, nil
	}
	return master, errors.New("pool master not found")
}

func setSessionDetails(session *Session) error {
	err := setAPIVersion(session)
	if err != nil {
		return err
	}
	err = setXAPIVersion(session)
	if err != nil {
		return err
	}
	return nil
}

func setAPIVersion(session *Session) error {
	session.APIVersion = APIVersionUnknown
	masterRef, err := getPoolMaster(session)
	if err != nil {
		return err
	}
	hostRecord, err := Host.GetRecord(session, masterRef)
	if err != nil {
		return err
	}
	session.APIVersion = GetAPIVersion(hostRecord.APIVersionMajor, hostRecord.APIVersionMinor)
	return nil
}

func setXAPIVersion(session *Session) error {
	masterRef, err := getPoolMaster(session)
	if err != nil {
		return err
	}
	hostRecord, err := Host.GetRecord(session, masterRef)
	if err != nil {
		return err
	}
	version, ok := hostRecord.SoftwareVersion["xapi"]
	if !ok {
		return errors.New("xapi version not found")
	}
	session.XAPIVersion = version
	return nil
}
