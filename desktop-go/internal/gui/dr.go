package gui

import (
	"errors"
	"sync"

	"xcpng-center/internal/app"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// Pool > Disaster Recovery (MainWindow.ConfigureDr and the DrWizard): Configure, and the wizard's storage, metadata,
// prechecks, recovery and test failover clean-up.

var (
	drBusy sync.Map // the connections running a recovery or a clean-up: one at a time per pool
	// shortcut: the guard is per connection key; a pool added twice (by name and by address) could run two at once. Key
	// it by the pool uuid if the server list ever allows that.
	drTests sync.Map // key -> *ops.DRRecovered: what the last test failover recovered, for its clean-up only
)

// DRSRs are the SRs that can hold the pool's metadata (DR on, metadata found) and the Configure dialog's title.
func (w *Workspace) DRSRs(key string) (out app.DRSRs, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadDRSRs(c); return err })
	return out, err
}

// DRWarning is DRConfigureDialog's warning for a choice and whether OK is allowed.
type DRWarning struct {
	Text    string `json:"text"`
	CanSave bool   `json:"can_save"`
}

// DRConfigureWarning is the Configure dialog's warning (originally: DR was on somewhere; checked: SRs ticked now).
func (w *Workspace) DRConfigureWarning(originally bool, checked int) DRWarning {
	text, ok := rules.DRConfigureWarning(originally, checked)
	return DRWarning{text, ok}
}

// ConfigureDR leaves database replication on exactly these SRs (uuids).
func (w *Workspace) ConfigureDR(key string, srs []string) error {
	err := w.with(key, func(c *xapi.Connection) error { return ops.ConfigureDR(c, srs) })
	w.changed(key)
	return err
}

// DRMetadata is the wizard's VMs and vApps page: the pools whose metadata is on these SRs (uuids).
func (w *Workspace) DRMetadata(key string, srs []string) (out []app.DRPool, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.LoadDRMetadata(c, srs); return err })
	return out, err
}

// DRPrechecks are the wizard's prechecks; a failover or failback also looks at the other pools connected here for
// VMs still running in their source pool (RunningVmCheck). A pool that does not answer is left out, as the C# app
// checks only the connections it can read: a failed primary site must not block its own failover.
func (w *Workspace) DRPrechecks(key, typ string, choices []ops.DRChoice) ([]rules.DRCheck, error) {
	sources := map[string]app.DRSource{}
	if typ != rules.DRTest {
		w.mu.Lock()
		var keys []string
		for k := range w.conns {
			if k != key {
				keys = append(keys, k)
			}
		}
		w.mu.Unlock()
		for _, k := range keys {
			_ = w.with(k, func(c *xapi.Connection) error {
				uuid, running, err := app.ReadDRSource(c)
				if err == nil {
					sources[uuid] = running
				}
				return err
			})
		}
	}
	var out []rules.DRCheck
	err := w.with(key, func(c *xapi.Connection) (err error) { out, err = app.DRPrechecks(c, typ, choices, sources); return err })
	return out, err
}

// runDR runs a recovery or a clean-up (prechecks and its record of a test included), refusing a second one on the
// connection while it runs.
func (w *Workspace) runDR(key string, fn func() error) error {
	if _, running := drBusy.LoadOrStore(key, true); running {
		return errors.New("a disaster recovery is running on this pool: wait for it to end")
	}
	defer drBusy.Delete(key)
	defer w.changed(key)
	return fn()
}

// DRRecover runs the prechecks again (the pool may have changed since the wizard showed them) and, when none fails,
// recovers the choices; a test failover keeps what it recovered for DRTestCleanup.
func (w *Workspace) DRRecover(key, typ string, choices []ops.DRChoice, start string) error {
	return w.runDR(key, func() error {
		drTests.Delete(key) // a test not cleaned up (its clean-up was declined) is never offered for this run
		checks, err := w.DRPrechecks(key, typ, choices)
		if err != nil {
			return err
		}
		if rules.DRBlocked(checks) {
			return errors.New("the prechecks failed: run the wizard again")
		}
		return w.with(key, func(c *xapi.Connection) error {
			rec, err := ops.RecoverDR(c, typ, choices, start)
			if typ == rules.DRTest && !rec.Empty() {
				drTests.Store(key, &rec)
			}
			return err
		})
	})
}

// DRTestPending is how many vApps and VMs the last test failover recovered and did not remove yet.
func (w *Workspace) DRTestPending(key string) int {
	if v, ok := drTests.Load(key); ok {
		rec := v.(*ops.DRRecovered)
		return len(rec.VMs) + len(rec.VApps) + len(rec.Members)
	}
	return 0
}

// DRTestCleanup removes what the last test failover recovered (DrRules.TestCleanupAction).
func (w *Workspace) DRTestCleanup(key string) error {
	return w.runDR(key, func() error {
		v, ok := drTests.Load(key)
		if !ok {
			return errors.New("there is no test failover to clean up")
		}
		err := w.with(key, func(c *xapi.Connection) error { return ops.DRTestCleanup(c, *v.(*ops.DRRecovered)) })
		if err == nil {
			drTests.Delete(key)
		}
		return err
	})
}
