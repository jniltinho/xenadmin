package gui

import (
	"errors"
	"sync"

	"xcpng-center/internal/app"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// Pool > High Availability (MainWindow.ConfigureHa and DisableHa): the HA wizard, the restart priorities, Disable.

// haBusy holds the connections running an HA action (MainWindow._haBusy): one at a time per pool.
var haBusy sync.Map

// runHA runs an HA action on the connection, refusing a second one while it runs.
func (w *Workspace) runHA(key string, fn func(*xapi.Connection) error) error {
	if _, running := haBusy.LoadOrStore(key, true); running {
		return errors.New(rules.HAConfigureBlocker(true, false, true, ""))
	}
	defer haBusy.Delete(key)
	err := w.with(key, fn)
	w.changed(key)
	return err
}

// HAConfig is the pool's HA state with, when Configure can run, its plan (the VMs, their priorities and agility).
func (w *Workspace) HAConfig(key string) (out app.HAConfig, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadHAConfig(c, true); return err })
	if _, running := haBusy.Load(key); err == nil && running && out.Blocker == "" {
		out = app.HAConfig{Pool: out.Pool, HAEnabled: out.HAEnabled, Blocker: rules.HAConfigureBlocker(true, false, true, "")}
	}
	return out, err
}

// HAMax is HaNtolControl's maximum for a plan (-1: it could not be computed) and the text shown with it.
type HAMax struct {
	Max  int    `json:"max"`
	Info string `json:"info"`
}

// HAMaxFailures is the most server failures the pool tolerates with these priorities (VM uuid to priority); a
// failure to compute it reads as unknown, as HaNtolControl.
func (w *Workspace) HAMaxFailures(key string, priorities map[string]string) (HAMax, error) {
	n := -1
	err := w.with(key, func(c *xapi.Connection) error {
		if m, err := ops.HAMaxFailures(c, priorities); err == nil {
			n = m
		}
		return nil
	})
	return HAMax{Max: n, Info: rules.HANtolInfo(n)}, err
}

// HeartbeatSRs are the HA wizard's heartbeat SRs, each usable or with the reason it is not.
func (w *Workspace) HeartbeatSRs(key string) (out []ops.HeartbeatSR, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = ops.HeartbeatSRs(c); return err })
	return out, err
}

// EnableHA turns HA on with the heartbeat SR (uuid), the plan and the failures to tolerate.
func (w *Workspace) EnableHA(key, heartbeatSR string, plan []ops.HAStartup, ntol int) error {
	return w.runHA(key, func(c *xapi.Connection) error { return ops.EnableHA(c, heartbeatSR, plan, ntol) })
}

// SetHAPriorities changes the restart priorities, start order and delay and the failures to tolerate, HA on.
func (w *Workspace) SetHAPriorities(key string, plan []ops.HAStartup, ntol int) error {
	return w.runHA(key, func(c *xapi.Connection) error { return ops.SetHAPriorities(c, plan, ntol) })
}

// HADisableQuestion is Disable's question (HADisableCommand), or why it cannot run (Error).
func (w *Workspace) HADisableQuestion(key string) (q HostQuestion, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		if why := ops.DisableHABlocker(c); why != "" {
			q.Error = why
			return nil
		}
		cfg, err := app.ReadHAConfig(c, false)
		if err == nil {
			q.Title, q.Text = rules.HADisableConfirmation(ops.Ellipsise(cfg.Pool, 30))
		}
		return err
	})
	return q, err
}

// DisableHA turns HA off (after HADisableQuestion).
func (w *Workspace) DisableHA(key string) error {
	return w.runHA(key, ops.DisableHA)
}
