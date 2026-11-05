# Architecture

```
main.go → cmd (Cobra: GUI when no subcommand, CLI, mcp)
             │
             ├── internal/mcp     MCP server (modelcontextprotocol/go-sdk), the C# tool names
             ├── internal/rules   UI-free decisions ported from XenAdmin.Desktop.Core (VmOperations)
             ├── internal/ops     operations on a server (VM and host power, maintenance), XAPI tasks
             ├── internal/console the VM console: TLS + CONNECT to the host, a loopback WebSocket for noVNC
             ├── internal/messages the C# texts, generated from the .resx files by tools/resx2go
             ├── internal/gui     Wails v3 window (Workspace service) + Vue/Tailwind SPA in frontend/
             │
             └── internal/app     use cases shared by every front end
                    ├── internal/xapi     addresses, HTTP client with the trust policy, login, coordinator redirect
                    ├── internal/trust    TOFU pins, changed certificates, operator pins (CertificateTrust.cs)
                    └── internal/config   settings.toml (servers, fingerprints; never passwords)
sdk/xenapi   XenServer Go SDK 26.16.1, its own module, with the patches in PATCHES.md
```

The core (`internal/{xapi,trust,config,app}`) imports no front end; `internal/archtest` fails the build if it does.
That keeps one core for the desktop app, the CLI, the MCP server and a future web server (`internal/server`).
The C# sources it ports are cited in each package (`ServerAddress.cs`, `CertificateTrust.cs`, `XcpTools.cs`).

## Live updates

The C# app keeps a cache per connection fed by `event.from` (XenObjectDownloader). Here `app.Watch` runs that loop
for each connected server, on a copy of the session made with `Session.WithContext` (an SDK patch), so that
`Disconnect` can cut a waiting call and wait for the loop before logging the session out. The loop only says that
something changed: the Workspace coalesces the batches of one operation (700 ms) into one Wails event `xcp:changed`
with the server's key, and the window reads that server's objects again (one read at a time; a read of a server
disconnected meanwhile is dropped), then the open tab and the alerts. A cache per object, with diffs, is task 2.4.
