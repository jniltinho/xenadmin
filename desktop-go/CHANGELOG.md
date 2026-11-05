# Changelog — XCP-ng Center (Go)

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions are the `go-v*` tags.

## [Unreleased]

### Added
- Server Shut Down, Reboot, Enter Maintenance Mode (saves where the VMs are, evacuates, hands the coordinator role
  over) and Exit Maintenance Mode (restores the VMs), with the C# questions; XAPI calls run as polled tasks.
- MCP `list_alerts`, and with `--allow-write`: `host_power`, `host_maintenance`, `dismiss_alerts`.
- Notifications > Alerts: the alerts of every connected server with the C# texts (FriendlyNames, failed logins),
  Dismiss, Dismiss Selected, Dismiss All with the C# confirmations; `tools/resx2go` generates `internal/messages`.
- `make windows` (WebView2, no CGO) and `make macos` (a universal `.app` zip, cross-built with `-tags purego`).
- VM power operations (Start, Shut Down, Reboot, Suspend, Resume, Force Shut Down, Force Reboot) in the toolbar, the
  VM menu and the context menu, with the C# confirmations and HA start check (`internal/rules`, `internal/ops`);
  MCP `vm_power` behind `--allow-write` and `confirm: true`.
- The desktop window (Wails v3, Vue 3, Tailwind 4): Add New Server with the certificate question, the tree with the
  C# app's icons, the General tab, a status bar; `make frontend` and `make linux` build it in containers.
- One binary (`xcp-ng-center-go`) with Cobra: `server add|list`, `pool|host|vm|template|sr|network list` (table or
  `--json`), `version`, and `mcp`.
- Connection core: server addresses (IPv6; 80/8080 refused), TLS 1.2+ without redirects, trust on first use with
  SHA-256 pins and `--trust-fingerprint`, the pool coordinator redirect, error messages without secrets.
- MCP server on the official Go SDK: `list_servers`, `connect`, `disconnect`, `list_objects` (read-only), with the
  C# app's tool names.
- The XenServer Go SDK 26.16.1 vendored in `sdk/xenapi` with documented patches (the app's HTTP client, no
  `log.Fatal`, typed `APIError`, `Session.Ref`).
