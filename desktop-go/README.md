# XCP-ng Center (Go)

One binary, `xcp-ng-center-go`, for the desktop app, the command line and an MCP server for AI agents. It follows the
proposal in [`openspec/changes/add-go-desktop-app`](../openspec/changes/add-go-desktop-app/proposal.md): Go, the
XenServer Go SDK, Cobra, the official MCP Go SDK, and Wails v3 with Vue and Tailwind for the window, built for
Linux, macOS and Windows and ready to become a web app.

**Status: first stage.** The connection core, the CLI and a read-only MCP server work against real servers; the desktop
window (`xcp-ng-center-go` without a subcommand) connects and shows the tree and the General tab. The C# app in
[`XenAdmin.Desktop`](../XenAdmin.Desktop/README.md) remains the complete cross-platform client.

## What works

- Connection: `host[:port]` and IPv6 addresses (80 and 8080 refused: plain HTTP), TLS 1.2+, no HTTP redirects,
  certificate trust on first use with SHA-256 pins (a changed certificate is refused), `--trust-fingerprint` for
  unattended use, the pool coordinator redirect (`HOST_IS_SLAVE`). Passwords never reach logs, errors or files.
- Window (Wails v3, GTK 4 + WebKitGTK 6.0 on Linux): Add New Server with the certificate question, the tree (pool,
  hosts, VMs, SRs, with the C# app's icons), the General tab, square corners like the Windows app; VM Start, Shut Down,
  Reboot, Suspend, Resume and the forced ones from the toolbar, the VM menu and the context menu, enabled by the VM's
  allowed operations and with the C# confirmations; Notifications > Alerts with the C# texts, Dismiss, Dismiss
  Selected and Dismiss All; Server Shut Down, Reboot and Enter/Exit Maintenance Mode (evacuation, restoring the VMs);
  the Console tab (noVNC) for VMs and servers, through a loopback WebSocket with one-time tickets (the session and the
  certificates stay in the Go side).
- CLI: `server add|list`, `pool|host|vm|template|sr|network list` (table or `--json`), `version`.
- MCP (`mcp`, stdio): `list_servers`, `connect`, `disconnect`, `list_objects` — the names and arguments of the C# app's
  MCP server, plus `list_alerts`; `mcp --allow-write` adds `vm_power`, `host_power`, `host_maintenance` and
  `dismiss_alerts`, which run only with `confirm: true`. The operator's
  password comes from `XCP_PASS`; the agent never sees it.

## Use

```bash
export XCP_HOST=192.168.0.10 XCP_USER=root     # or --host / --user
read -rs XCP_PASS && export XCP_PASS           # or type it when asked (never a flag)
xcp-ng-center-go server add                    # asks to trust the certificate on first use
xcp-ng-center-go vm list --json
xcp-ng-center-go mcp                           # for an MCP client (stdio)
```

Settings (servers and certificate fingerprints, no passwords) live in `settings.toml` under the user config directory
(`~/.config/xcp-ng-center-go` on Linux), mode 0600.

## Build (only docker and make)

```bash
make check     # gofmt, go vet (app and SDK), tests
make cli       # static CLI/MCP binaries for Linux, macOS (arm64, amd64) and Windows in dist/
make linux     # the desktop app for Linux (builds the frontend first, node:24 + Ubuntu 24.04 GTK/WebKit images)
make windows   # the desktop app for Windows (WebView2; no CGO)
make macos     # "XCP-ng Center (Go).app" for Apple Silicon and Intel, zipped (cross-built, see below)
```

### macOS

`make macos` cross-builds the app on Linux with Wails' CGO-free backend (`-tags purego`, experimental in Wails 3.1):
one binary per CPU inside the bundle and a launcher that picks it, like the C# app's macOS zip. It is unsigned and not
yet tried on a Mac: after downloading, `xattr -dr com.apple.quarantine "XCP-ng Center (Go).app"`. On a Mac, the
native build is `cd frontend && npm ci && npm run build && cd .. && CGO_ENABLED=1 go build .` (Xcode command line
tools), then sign and notarize before a release that claims macOS support.

See [docs/architecture.md](docs/architecture.md) and [docs/development.md](docs/development.md).
