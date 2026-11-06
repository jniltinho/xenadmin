package gui

import (
	"xcpng-center/internal/app"
	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// The Networking and NICs tabs' buttons of a server or pool (HostNetworkView, NicsView): New Network wizard, network
// Properties, Remove Network, Create Bond, Delete Bond and Rescan.

// HostNetworking is the Networking tab of a server ("hosts") or pool ("pools") with its buttons.
func (w *Workspace) HostNetworking(key, kind, uuid string) (out app.HostNetworking, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.HostNetworkingTab(c, kind, uuid); return err })
	return out, err
}

// HostNics is the NICs tab of a server with its buttons.
func (w *Workspace) HostNics(key, uuid string) (out app.HostNics, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.HostNicsTab(c, uuid); return err })
	return out, err
}

// NewNetworkOptions are the New Network wizard's (and Create Bond's) choices for a server or pool.
func (w *Workspace) NewNetworkOptions(key, kind, uuid string) (out app.NewNetworkChoices, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.NewNetworkOptions(c, kind, uuid); return err })
	return out, err
}

// CreateNetwork makes an External (VLAN) or Single-Server Private network.
func (w *Workspace) CreateNetwork(key string, spec ops.NewNetwork) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.CreateNetwork(c, spec) })
}

// BondQuestion is the question Create Bond asks first (nil: none; Block: an error, the bond cannot be made).
func (w *Workspace) BondQuestion(key string, members []string) (out *rules.Warning, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = ops.BondCreateWarning(c, members); return err })
	return out, err
}

// CreateBond bonds the NICs (after BondQuestion).
func (w *Workspace) CreateBond(key string, spec ops.BondSpec) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.CreateBond(c, spec) })
}

// DeleteBondQuestion is Delete Bond's question for a bond network (Block: an error, the bond cannot be deleted).
func (w *Workspace) DeleteBondQuestion(key, networkUUID string) (out rules.Warning, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = ops.BondDeleteWarning(c, networkUUID); return err })
	return out, err
}

// DeleteBond deletes the bond of a bond network on every server (after DeleteBondQuestion).
func (w *Workspace) DeleteBond(key, networkUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.DestroyBond(c, networkUUID) })
}

// RescanNics looks for new NICs on a server.
func (w *Workspace) RescanNics(key, hostUUID string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.RescanPIFs(c, hostUUID) })
}

// NetworkEdit reads a network for its Properties.
func (w *Workspace) NetworkEdit(key, uuid string) (out app.NetworkEdit, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadNetworkEdit(c, uuid); return err })
	return out, err
}

// SaveNetwork saves a network's Properties.
func (w *Workspace) SaveNetwork(key, uuid string, ch ops.NetworkChanges) error {
	return w.with(key, func(c *xapi.Connection) error { return app.SaveNetwork(c, uuid, ch) })
}

// RemoveNetworkQuestion is Remove Network's question (MESSAGEBOX_NETWORK_DELETE); a bond network asks DeleteBondQuestion.
func (w *Workspace) RemoveNetworkQuestion() Confirm {
	return Confirm{Title: messages.Get("MESSAGEBOX_NETWORK_DELETE_TITLE"), Text: messages.Get("MESSAGEBOX_NETWORK_DELETE")}
}

// RemoveNetwork removes a network that is not a bond.
func (w *Workspace) RemoveNetwork(key, uuid string) error {
	return w.with(key, func(c *xapi.Connection) error { return ops.RemoveNetwork(c, uuid) })
}
