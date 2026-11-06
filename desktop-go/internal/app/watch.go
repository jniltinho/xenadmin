package app

import (
	"context"
	"errors"
	"time"

	xenapi "go/xenapi"
)

// WatchClasses are the objects whose changes the window shows (the tree, the tabs and the alerts).
// shortcut: no vm_metrics/host_metrics, which change every few seconds (memory figures); a host going down is seen
// through its VMs and PBDs, not its "live" flag, until the reload is per object rather than per server
var WatchClasses = []string{"pool", "host", "vm", "sr", "pbd", "vbd", "vdi", "network", "vif", "message", "vm_guest_metrics", "vm_appliance",
	"pgpu", "gpu_group", "vgpu", "pusb", "usb_group", "vusb"}

// watchTimeout is how long one event.from waits for changes (under the HTTP client's 60 s); watchRetry is the pause
// after a failed call.
var (
	watchTimeout = 30.0
	watchRetry   = 5 * time.Second
)

// Watch is XenObjectDownloader's event loop: event.from until ctx ends, calling changed after the first batch (the
// baseline: whatever changed while the window read the objects is in it) and after each batch with events.
// EVENTS_LOST starts over from a new baseline; a failed call is retried after a pause. s must be a session of its
// own (Session.WithContext with ctx), so that cancelling cuts a waiting call before the session is logged out.
func Watch(ctx context.Context, s *xenapi.Session, changed func()) {
	token := ""
	for ctx.Err() == nil {
		batch, err := xenapi.Event.From(s, WatchClasses, token, watchTimeout)
		if ctx.Err() != nil {
			return
		}
		var api *xenapi.APIError
		switch {
		case errors.As(err, &api) && api.Message == "EVENTS_LOST":
			token = ""
			continue
		case err != nil:
			select {
			case <-ctx.Done():
				return
			case <-time.After(watchRetry):
			}
			continue
		}
		if token == "" || len(batch.Events) > 0 {
			changed()
		}
		token = batch.Token
	}
}
