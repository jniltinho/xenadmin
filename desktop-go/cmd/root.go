// Package cmd is the command line of the single xcp-ng-center-go binary: with no subcommand it opens the desktop
// app; the subcommands are the CLI and the MCP server (mcp).
package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"xcpng-center/internal/config"
	"xcpng-center/internal/trust"
	"xcpng-center/internal/xapi"
)

var flags struct {
	json             bool
	config           string
	host             string
	user             string
	trustFingerprint string
}

var rootCmd = &cobra.Command{
	Use:           "xcp-ng-center-go",
	Short:         "XCP-ng Center: desktop app, command line and MCP server",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGUI(cmd.Context())
	},
}

// Execute runs the command line.
func Execute() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "error:", redact(err.Error()))
		os.Exit(1)
	}
}

func init() {
	p := rootCmd.PersistentFlags()
	p.BoolVar(&flags.json, "json", false, "print JSON instead of tables")
	p.StringVar(&flags.config, "config", "", "settings file (default: the user config dir)")
	p.StringVar(&flags.host, "host", os.Getenv("XCP_HOST"), "server address (host[:port]); default $XCP_HOST")
	p.StringVar(&flags.user, "user", envOr("XCP_USER", "root"), "user name; default $XCP_USER or root")
	p.StringVar(&flags.trustFingerprint, "trust-fingerprint", os.Getenv("XCP_TRUST_FINGERPRINT"), "accept only this SHA-256 certificate fingerprint; default $XCP_TRUST_FINGERPRINT")
}

func envOr(name, def string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return def
}

func loadSettings() (*config.Settings, error) {
	path := flags.config
	if path == "" {
		var err error
		if path, err = config.DefaultPath(); err != nil {
			return nil, err
		}
	}
	return config.Load(path)
}

// policy: the stored pins, plus the operator's fingerprint for any host when given.
func policy(s *config.Settings, requirePin bool) *trust.Policy {
	p := &trust.Policy{Store: s, OperatorPins: map[string]string{}, RequirePin: requirePin}
	if flags.trustFingerprint != "" {
		p.OperatorPins["*"] = flags.trustFingerprint
	}
	return p
}

// typedPassword is kept only to redact it from what the CLI prints.
var typedPassword string

// redact removes the password and session references from any text the CLI prints.
func redact(text string) string { return xapi.Redact(text, os.Getenv("XCP_PASS"), typedPassword) }

// password: $XCP_PASS, else a prompt without echo on a terminal. Never a flag (it would show in ps and history).
func password() (string, error) {
	if v := os.Getenv("XCP_PASS"); v != "" {
		return v, nil
	}
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return "", errors.New("set XCP_PASS or run in a terminal to type the password")
	}
	fmt.Fprint(os.Stderr, "Password: ")
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	typedPassword = string(b)
	return typedPassword, err
}

// connect logs in to --host; an untrusted certificate is offered for acceptance on a terminal.
func connect(ctx context.Context, s *config.Settings) (*xapi.Connection, error) {
	if flags.host == "" {
		return nil, errors.New("no server: use --host or XCP_HOST")
	}
	host, port, err := xapi.ParseAddress(flags.host)
	if err != nil {
		return nil, err
	}
	pass, err := password()
	if err != nil {
		return nil, err
	}
	p := policy(s, false)
	conn, err := xapi.Connect(ctx, p, host, port, flags.user, pass)
	var untrusted *trust.UntrustedError
	// a new certificate may be accepted on a terminal, never when the operator gave the fingerprint to expect
	if errors.As(err, &untrusted) && untrusted.Previous == "" && flags.trustFingerprint == "" && term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintf(os.Stderr, "%s\nTrust it? [y/N] ", untrusted.Error())
		var answer string
		fmt.Fscanln(os.Stdin, &answer)
		if strings.EqualFold(answer, "y") || strings.EqualFold(answer, "yes") {
			if err := s.Pin(untrusted.HostKey, untrusted.Fingerprint); err != nil {
				return nil, err
			}
			conn, err = xapi.Connect(ctx, p, host, port, flags.user, pass)
		}
	}
	return conn, err
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
