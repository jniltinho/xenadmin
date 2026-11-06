package gui

import (
	"errors"
	"strings"
	"time"

	"xcpng-center/internal/config"
	"xcpng-center/internal/xapi"
)

// Options are Tools > Options.
func (w *Workspace) Options() config.Options { return w.settings.GetOptions() }

// SetOptions checks and saves the options (OptionsDialog's OK) and applies the network ones to the next connections.
func (w *Workspace) SetOptions(o config.Options) error {
	o.ProxyAddress = strings.TrimSpace(o.ProxyAddress)
	switch o.Proxy {
	case "", "direct", "system":
	case "specified":
		if o.ProxyAddress == "" || strings.ContainsAny(o.ProxyAddress, " /") {
			return errors.New("enter the address of the proxy server")
		}
		if o.ProxyPort < 1 || o.ProxyPort > 65535 {
			return errors.New("enter a port between 1 and 65535")
		}
	default:
		return errors.New("choose how to connect to the servers")
	}
	if o.ConnectionTimeout < 1 || o.ConnectionTimeout > 2147483 { // AppOptions.MinTimeout / MaxTimeout
		return errors.New("the connection timeout must be between 1 and 2147483 seconds")
	}
	if err := w.settings.SetOptions(o); err != nil {
		return err
	}
	applyNetwork(o)
	return nil
}

// applyNetwork hands the proxy and timeout options to the connections.
func applyNetwork(o config.Options) {
	timeout := o.ConnectionTimeout
	if timeout <= 0 {
		timeout = 20
	}
	proxy := o.Proxy
	if proxy == "" {
		proxy = "direct" // AppOptions.ProxyStyle's default
	}
	xapi.SetNetwork(xapi.Network{Proxy: proxy, Address: o.ProxyAddress, Port: o.ProxyPort, Bypass: o.BypassProxyForServers,
		Timeout: time.Duration(timeout) * time.Second})
}
