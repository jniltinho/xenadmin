package gui

import (
	"xcpng-center/internal/app"
	"xcpng-center/internal/xapi"
)

// MemoryBars are the Memory tab's bars of a VM, server or pool.
func (w *Workspace) MemoryBars(key, kind, uuid string) (out app.MemoryBars, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadMemoryBars(c, kind, uuid); return err })
	return out, err
}
