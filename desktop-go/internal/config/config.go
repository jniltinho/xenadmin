// Package config is the app's settings file: the server list and the accepted certificate fingerprints, in
// settings.toml under os.UserConfigDir()/xcp-ng-center-go, readable by the user only. Passwords are never stored.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

// Server is one entry of the server list.
type Server struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	User string `toml:"user"`
	Name string `toml:"name,omitempty"`
}

// Settings is settings.toml.
type Settings struct {
	Servers      []Server          `toml:"servers"`
	Fingerprints map[string]string `toml:"fingerprints"`

	path string
	mu   sync.Mutex
}

// DefaultPath is os.UserConfigDir()/xcp-ng-center-go/settings.toml.
func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "xcp-ng-center-go", "settings.toml"), nil
}

// Load reads the file (a missing file is an empty configuration).
func Load(path string) (*Settings, error) {
	s := &Settings{path: path, Fingerprints: map[string]string{}}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if _, err := toml.Decode(string(data), s); err != nil {
		return nil, err
	}
	if s.Fingerprints == nil {
		s.Fingerprints = map[string]string{}
	}
	return s, nil
}

// Save writes the file atomically with mode 0600 (its directory 0700).
func (s *Settings) Save() error {
	return s.update(func(*Settings) {})
}

// update applies a change under a lock shared by every process (the CLI and the app may run at once): it reloads
// the file, applies the change to the fresh copy, and writes it through a unique temporary file, so no process loses
// another's pins or servers.
func (s *Settings) update(change func(*Settings)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	unlock, err := lockFile(s.path + ".lock")
	if err != nil {
		return err
	}
	defer unlock()
	fresh, err := Load(s.path)
	if err != nil {
		return err
	}
	change(fresh)
	s.Servers, s.Fingerprints = fresh.Servers, fresh.Fingerprints // the in-memory copy follows the file
	f, err := os.CreateTemp(filepath.Dir(s.path), ".settings-*.toml")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if err := f.Chmod(0o600); err != nil {
		f.Close()
		return err
	}
	if err := toml.NewEncoder(f).Encode(fresh); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// lockFile takes an exclusive lock file (created with O_EXCL), waiting up to 5 s. A lock left by a crashed process
// is never taken over automatically (two processes could both take it): the error says how to release it.
func lockFile(path string) (func(), error) {
	deadline := time.Now().Add(5 * time.Second)
	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			f.Close()
			return func() { os.Remove(path) }, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("the settings file is locked by another process (if none is running, delete %s)", path)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// AddServer adds or updates an entry (by host and port) and saves.
func (s *Settings) AddServer(srv Server) error {
	return s.update(func(t *Settings) {
		for i, e := range t.Servers {
			if strings.EqualFold(e.Host, srv.Host) && e.Port == srv.Port {
				t.Servers[i].User = srv.User
				if srv.Name != "" {
					t.Servers[i].Name = srv.Name
				}
				return
			}
		}
		t.Servers = append(t.Servers, srv)
	})
}

// Path is the settings file.
func (s *Settings) Path() string { return s.path }

// RemoveServer takes a server (and its pinned certificate) off the list: the C# app's Forget/Remove server.
func (s *Settings) RemoveServer(host string, port int) error {
	return s.update(func(t *Settings) {
		kept := t.Servers[:0]
		for _, e := range t.Servers {
			if !(strings.EqualFold(e.Host, host) && e.Port == port) {
				kept = append(kept, e)
			}
		}
		t.Servers = kept
	})
}

// ServerList is a copy of the servers (the window reads it while the CLI or MCP may add one).
func (s *Settings) ServerList() []Server {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Server(nil), s.Servers...)
}

// Fingerprint implements trust.Store.
func (s *Settings) Fingerprint(hostKey string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fp, ok := s.Fingerprints[strings.ToLower(hostKey)]
	return fp, ok
}

// Pin implements trust.Store and saves.
func (s *Settings) Pin(hostKey, fingerprint string) error {
	return s.update(func(t *Settings) { t.Fingerprints[strings.ToLower(hostKey)] = fingerprint })
}

// PinNew pins a certificate only when the host has none in the file (checked under the lock): accepting a
// certificate in the window never replaces a pin, even one another process saved meanwhile.
func (s *Settings) PinNew(hostKey, fingerprint string) error {
	key, pinned := strings.ToLower(hostKey), false
	err := s.update(func(t *Settings) {
		if _, ok := t.Fingerprints[key]; ok {
			pinned = true
			return
		}
		t.Fingerprints[key] = fingerprint
	})
	if err == nil && pinned {
		return ErrPinned
	}
	return err
}

// ErrPinned: the host already has a pinned certificate.
var ErrPinned = errors.New("this server's certificate is already pinned: a changed certificate is not replaced from the window")
