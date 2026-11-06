package gui

import (
	"fmt"
	"sort"
	"strings"

	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// The Pool menu's Add Server, Rotate Pool Secret and Make into standalone server (MainWindow's Pool menu).

// JoinCandidates is PoolRules.JoinCandidates: the other connected servers that are not in a pool, each with why it
// cannot join the pool of poolKey ("" when it can).
func (w *Workspace) JoinCandidates(poolKey string) ([]PoolCandidate, error) {
	facts, err := w.poolFacts()
	if err != nil {
		return nil, err
	}
	pf, ok := facts[poolKey]
	if !ok {
		return nil, errNotConnected
	}
	out := []PoolCandidate{}
	for k, f := range facts {
		if k == poolKey || f.PoolName != "" || f.Hosts > 1 { // Helpers.GetPool(c) == null
			continue
		}
		r := rules.CanJoinPool(f, pf)
		out = append(out, PoolCandidate{Key: k, Name: f.Name, Reason: rules.JoinReasonText(r), reason: r, CanJoin: r == rules.JoinAllowed})
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	return out, nil
}

// AddToPoolQuestions are AddHostToPoolCommand's questions: moving the server, then the CPU warning (each asked in turn).
func (w *Workspace) AddToPoolQuestions(poolKey, memberKey string) ([]string, error) {
	facts, err := w.poolFacts()
	if err != nil {
		return nil, err
	}
	pf, mf, err := joinFacts(facts, poolKey, memberKey)
	if err != nil {
		return nil, err
	}
	out := []string{messages.Format(messages.Get("MAINWINDOW_CONFIRM_MOVE_TO_POOL"), ops.Ellipsise(mf.Name, 500), ops.Ellipsise(pf.PoolName, 500))}
	for _, warn := range []string{rules.PackWarning([]rules.PoolFacts{mf}, pf), rules.CPUWarning([]rules.PoolFacts{mf}, pf)} {
		if warn != "" {
			out = append(out, warn)
		}
	}
	return out, nil
}

// AddToPool is AddHostToPoolAction: the server joins the pool (pool.join with the pool's credentials) and leaves the
// server list.
func (w *Workspace) AddToPool(poolKey, memberKey string) error {
	facts, err := w.poolFacts()
	if err != nil {
		return err
	}
	pf, mf, err := joinFacts(facts, poolKey, memberKey)
	if err != nil {
		return err
	}
	if r := rules.CanJoinPool(mf, pf); r != rules.JoinAllowed {
		return fmt.Errorf("%s: %s", mf.Name, rules.JoinReasonText(r))
	}
	w.mu.Lock()
	coord, password := w.conns[poolKey], w.passwords[poolKey]
	w.mu.Unlock()
	if coord == nil {
		return errNotConnected
	}
	var address string
	if err := w.with(poolKey, func(c *xapi.Connection) (err error) { address, err = ops.CoordinatorAddress(c); return err }); err != nil {
		return err
	}
	if err := w.with(memberKey, func(c *xapi.Connection) error { return ops.JoinPool(c, address, coord.User, password) }); err != nil {
		return err
	}
	_ = w.Forget(memberKey) // a member now: the pool's connection shows it
	err = w.with(poolKey, func(c *xapi.Connection) error { return ops.ClearNonSharedSrs(c) })
	w.changed(poolKey)
	return err
}

// joinFacts are the pool's and the member's facts, both read (a server not read may hide a warning).
func joinFacts(facts map[string]rules.PoolFacts, poolKey, memberKey string) (rules.PoolFacts, rules.PoolFacts, error) {
	pf, ok1 := facts[poolKey]
	mf, ok2 := facts[memberKey]
	switch {
	case !ok1 || !ok2:
		return pf, mf, errNotConnected
	case pf.ReadError != "":
		return pf, mf, fmt.Errorf("%s could not be read: try again", poolKey)
	case mf.ReadError != "":
		return pf, mf, fmt.Errorf("%s could not be read: try again", memberKey)
	}
	return pf, mf, nil
}

// RotateSecretBlocker is why the pool secret cannot be rotated ("" when it can).
func (w *Workspace) RotateSecretBlocker(key string) (out string, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out = ops.RotateSecretBlocker(c); return nil })
	return out, err
}

// RotateSecret rotates the pool secret; the reminder to change the root password comes back to be shown after.
func (w *Workspace) RotateSecret(key string) (string, error) {
	err := w.with(key, func(c *xapi.Connection) error { return ops.RotateSecret(c) })
	if err != nil || w.settings.GetOptions().NoRotateSecretReminder {
		return "", err
	}
	return messages.Get("ROTATE_POOL_SECRET_REMIND_CHANGE_PASSWORD"), nil
}

// MakeStandaloneBlocker is why the pool cannot be made into a standalone server ("" when it can).
func (w *Workspace) MakeStandaloneBlocker(key string) (out string, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out = ops.MakeStandaloneBlocker(c); return nil })
	return out, err
}

// MakeStandalone makes a pool of one server into a standalone server.
func (w *Workspace) MakeStandalone(key string) error {
	err := w.with(key, func(c *xapi.Connection) error { return ops.MakeStandalone(c) })
	if err == nil {
		w.changed(key)
	}
	return err
}
