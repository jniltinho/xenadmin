package gui

import (
	"errors"

	"xcpng-center/internal/app"
	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/xapi"
)

// CanChangePassword is PoolRules.CanChangePassword: the connection is logged in as the local root.
func (w *Workspace) CanChangePassword(key string) (out bool, err error) {
	err = w.with(key, func(c *xapi.Connection) error {
		su, err := c.Session.GetIsLocalSuperuser(c.Session.Ref())
		out = err == nil && su
		return nil
	})
	return out, err
}

// ChangePassword is Change Server Password (ChangeServerPasswordDialog, ChangeHostPasswordAction): the current
// password must be the connection's, the new one typed twice; session.change_password, then the connection keeps the
// new password. Passwords never reach a log or an error text.
func (w *Workspace) ChangePassword(key, current, next, confirm string) error {
	if current == "" || next == "" || confirm == "" {
		return errors.New("enter the current password and the new one twice")
	}
	w.mu.Lock()
	known, conn := w.passwords[key], w.conns[key]
	w.mu.Unlock()
	if current != known { // PoolRules.ChangePasswordError
		return errors.New(messages.Get("PASSWORD_INCORRECT"))
	}
	if next != confirm {
		return errors.New(messages.Get("PASSWORDS_DONT_MATCH"))
	}
	err := w.with(key, func(c *xapi.Connection) error {
		if su, err := c.Session.GetIsLocalSuperuser(c.Session.Ref()); err != nil || !su {
			return errors.New("only the local root account can change the server password")
		}
		return xapi.Sanitize(c.Session.ChangePassword(current, next), current, next)
	})
	if err != nil {
		return err
	}
	w.mu.Lock() // only for the connection that changed it (not after a Disconnect, Close or another login)
	if !w.closed && conn != nil && w.conns[key] == conn {
		w.passwords[key] = next
	}
	w.mu.Unlock()
	return nil
}

// ChangePasswordTitle is the dialog's title for the connection (CHANGEPASS_DIALOG_TITLE).
func (w *Workspace) ChangePasswordTitle(key string) string {
	name := key
	_ = w.with(key, func(c *xapi.Connection) error { name = ops.Ellipsise(app.PoolName(c), 30); return nil })
	return messages.Format(messages.Get("CHANGEPASS_DIALOG_TITLE"), name)
}
