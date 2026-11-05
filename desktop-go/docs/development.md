# Development

- Go 1.27; everything runs in `golang:1.27.3` through the Makefile (`make check`, `make cli`).
- Tests and `go vet` run with `-tags nogui` (no CGO, no GTK); the window builds with `make linux`, which runs
  `make frontend` (Vue/Vite into `internal/gui/dist`, embedded) and links GTK 4 + WebKitGTK 6.0 (`build/linux.Dockerfile`).
- The window's API is `internal/gui.Workspace` (Wails service); the SPA calls it in `frontend/src/api/workspace.ts`.
  Frontend alone: `cd frontend && npm install && npm run build`.
- Running the window under Xvfb in a container needs WebKit's sandbox (bubblewrap user namespaces); try it on a host
  with GTK 4 and WebKitGTK 6.0 instead, with its own `HOME` so your settings are untouched.
- Tests use a fake JSON-RPC server over TLS (`internal/xapi/client_test.go`) and an in-memory MCP client
  (`internal/mcp/server_test.go`); none needs a real server.
- Against a server: `XCP_HOST`, `XCP_USER`, `XCP_PASS` in the environment (never in files); the lab servers and rules
  are in the repository's [AGENTS.md](../../AGENTS.md) (read-only `.236`, view-only `.114`).
- A new SDK: `tools/sdk-update.sh XenServer-SDK-<version>.zip` replaces `sdk/xenapi` and re-applies the
  `XCP-NG PATCH` changes listed in `sdk/xenapi/PATCHES.md`.
