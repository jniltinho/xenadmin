package ops

import (
	"errors"
	"maps"
	"slices"
	"strings"

	xenapi "go/xenapi"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/rules"
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

// poolOf reads the connection's pool.
func poolOf(s *xenapi.Session) (xenapi.PoolRef, xenapi.PoolRecord, int, error) {
	pools, err := xenapi.Pool.GetAllRecords(s)
	if err != nil {
		return "", xenapi.PoolRecord{}, 0, xapi.Sanitize(err)
	}
	hosts, err := xenapi.Host.GetAll(s)
	if err != nil {
		return "", xenapi.PoolRecord{}, 0, xapi.Sanitize(err)
	}
	for r, p := range pools {
		return r, p, len(hosts), nil
	}
	return "", xenapi.PoolRecord{}, 0, errors.New("the server has no pool")
}

// RotateSecretBlocker is PoolRules.CanRotateSecret and RotateSecretRestriction: why the pool secret cannot be rotated
// ("" when it can): HA on, a rolling pool upgrade, the licence, or a role without pool.rotate_secret.
func RotateSecretBlocker(c *xapi.Connection) string {
	s := c.Session
	_, p, _, err := poolOf(s)
	if err != nil {
		return err.Error()
	}
	hosts, err := xenapi.Host.GetAllRecords(s)
	if err != nil {
		return xapi.Sanitize(err).Error()
	}
	coordinator := hosts[p.Master]
	switch {
	case rules.CompareVersions(coordinator.SoftwareVersion["platform_version"], "3.1.50") < 0: // Helpers.StockholmOrGreater
		return "The pool secret can only be rotated from XCP-ng 8.2 on."
	case p.HaEnabled:
		return messages.Get("ROTATE_POOL_SECRET_HA")
	case hasKey(p.OtherConfig, "rolling_upgrade_in_progress"): // Pool.RollingUpgrade
		return messages.Get("ROTATE_POOL_SECRET_RPU")
	case featureForbidden(hosts, "restrict_pool_secret_rotation"):
		return "The pool secret rotation is not available due to license restrictions on the server."
	}
	if err := permitted(s, "pool.rotate_secret"); err != nil {
		return err.Error()
	}
	return ""
}

// RotateSecret is RotatePoolSecretCommand: pool.rotate_secret.
func RotateSecret(c *xapi.Connection) error {
	if why := RotateSecretBlocker(c); why != "" {
		return errors.New(why)
	}
	s := c.Session
	return xapi.Sanitize(wait(s)(xenapi.Pool.AsyncRotateSecret(s)))
}

// MakeStandaloneBlocker is PoolRules.CanMakeStandalone and MakeStandaloneBlocker: a visible pool of one server, HA off.
func MakeStandaloneBlocker(c *xapi.Connection) string {
	_, p, hosts, err := poolOf(c.Session)
	switch {
	case err != nil:
		return err.Error()
	case p.NameLabel == "" && hosts < 2: // Pool.IsVisible
		return "The server is not in a pool."
	case slices.Contains(slices.Collect(maps.Values(p.CurrentOperations)), xenapi.PoolAllowedOperationsHaEnable):
		return messages.Format(messages.Get("POOL_DELETE_HA_ENABLING"), p.NameLabel)
	case p.HaEnabled:
		return messages.Format(messages.Get("POOL_DELETE_HA_ENABLED"), p.NameLabel)
	case hosts > 1:
		return messages.Get("MESSAGEBOX_POOL_MEMBERS_EJECT")
	}
	return ""
}

// MakeStandalone is DestroyPoolAction ("Make into standalone server"): the pool of one server loses its name and
// description.
func MakeStandalone(c *xapi.Connection) error {
	if why := MakeStandaloneBlocker(c); why != "" {
		return errors.New(why)
	}
	s := c.Session
	ref, _, _, err := poolOf(s)
	if err != nil {
		return err
	}
	if err := xenapi.Pool.SetNameLabel(s, ref, ""); err != nil {
		return xapi.Sanitize(err)
	}
	return xapi.Sanitize(xenapi.Pool.SetNameDescription(s, ref, ""))
}

func hasKey(m map[string]string, k string) bool { _, ok := m[k]; return ok }

// ClearNonSharedSrs is the pool actions' ClearNonSharedSrs: a default SR that is not shared cannot serve the members,
// so the pool keeps no default, suspend or crash dump SR on local storage.
func ClearNonSharedSrs(c *xapi.Connection) error {
	s := c.Session
	ref, p, _, err := poolOf(s)
	if err != nil {
		return err
	}
	if p.DefaultSR == "" || p.DefaultSR == nullRef {
		return nil
	}
	shared, err := xenapi.SR.GetShared(s, p.DefaultSR)
	if err != nil || shared {
		return xapi.Sanitize(err)
	}
	for _, set := range []func(*xenapi.Session, xenapi.PoolRef, xenapi.SRRef) error{xenapi.Pool.SetDefaultSR, xenapi.Pool.SetSuspendImageSR, xenapi.Pool.SetCrashDumpSR} {
		if err := set(s, ref, nullRef); err != nil {
			return xapi.Sanitize(err)
		}
	}
	return nil
}

// CoordinatorAddress is the pool coordinator's management address (pool.join is given it, as AddHostToPoolAction).
func CoordinatorAddress(c *xapi.Connection) (string, error) {
	_, p, _, err := poolOf(c.Session)
	if err != nil {
		return "", err
	}
	a, err := xenapi.Host.GetAddress(c.Session, p.Master)
	if err != nil || a == "" {
		return c.Host, nil // the address the window connects to
	}
	return a, nil
}
