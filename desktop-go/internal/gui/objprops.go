package gui

import (
	"xcpng-center/internal/app"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/xapi"
)

// ObjectProperties reads a server, pool, SR, network (or VM) for the General and Custom Fields pages of Properties.
func (w *Workspace) ObjectProperties(key, kind, uuid string) (out app.ObjectProperties, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadObjectProperties(c, kind, uuid); return err })
	return out, err
}

// SaveObjectProperties applies those pages' changes.
func (w *Workspace) SaveObjectProperties(key, kind, uuid string, ch ops.ObjectChanges) error {
	err := w.with(key, func(c *xapi.Connection) error { return ops.SaveObject(c, kind, uuid, ch) })
	if err == nil {
		w.changed(key)
	}
	return err
}
