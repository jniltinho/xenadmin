package gui

import (
	"xcpng-center/internal/app"
	"xcpng-center/internal/xapi"
)

// General reads the General tab of a VM, template, server, pool, SR or network (Details.General).
func (w *Workspace) General(key, kind, uuid string) (out []app.Section, err error) {
	err = w.with(key, func(c *xapi.Connection) error { out, err = app.ReadGeneral(c, kind, uuid); return err })
	return out, err
}
