package ops

import (
	"errors"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/xapi"
)

// NamePool is the first part of CreatePoolAction on the coordinator: the pool's name and description, then
// ClearNonSharedSrs (a default SR that is not shared cannot serve the members: no default, suspend or crash dump SR).
// named tells whether the pool got its name (the server is a pool from then on, even when a later step fails).
func NamePool(c *xapi.Connection, name, description string) (named bool, err error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return false, errors.New("enter a name for the pool")
	}
	s := c.Session
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return false, xapi.Sanitize(err)
	}
	for ref, p := range pools {
		if err := xenapi.Pool.SetNameLabel(s, ref, name); err != nil {
			return false, xapi.Sanitize(err)
		}
		named = true
		if err := xenapi.Pool.SetNameDescription(s, ref, description); err != nil {
			return named, xapi.Sanitize(err)
		}
		if p.DefaultSR != "" && p.DefaultSR != nullRef {
			shared, err := xenapi.SR.GetShared(s, p.DefaultSR)
			if err != nil {
				return named, xapi.Sanitize(err) // never join members to a pool whose defaults may be local storage
			}
			if !shared {
				for _, set := range []func(*xenapi.Session, xenapi.PoolRef, xenapi.SRRef) error{xenapi.Pool.SetDefaultSR, xenapi.Pool.SetSuspendImageSR, xenapi.Pool.SetCrashDumpSR} {
					if err := set(s, ref, nullRef); err != nil {
						return named, xapi.Sanitize(err)
					}
				}
			}
		}
	}
	return named, nil
}

// JoinPool is CreatePoolAction's pool.join from a member's session to the coordinator (its address and credentials,
// which never leave this call but for XAPI).
func JoinPool(member *xapi.Connection, coordinator, user, password string) error {
	s := member.Session
	return xapi.Sanitize(wait(s)(xenapi.Pool.AsyncJoin(s, coordinator, user, password)), password)
}
