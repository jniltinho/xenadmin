# Development

- Go 1.27; everything runs in `golang:1.27.3` through the Makefile (`make check`, `make cli`).
- Tests and `go vet` run with `-tags nogui` (no CGO, no GTK); the window builds with `make linux`, which runs
  `make frontend` (Vue/Vite into `internal/gui/dist`, embedded) and links GTK 4 + WebKitGTK 6.0 (`build/linux.Dockerfile`).
- The window's API is `internal/gui.Workspace` (Wails service); the SPA calls it in `frontend/src/api/workspace.ts`.
  Frontend alone: `cd frontend && npm install && npm run build`.
- Running the window under Xvfb in a container needs WebKit's sandbox (bubblewrap user namespaces); try it on a host
  with GTK 4 and WebKitGTK 6.0 instead, with its own `HOME` so your settings are untouched.
- Lab tests (`internal/app/lab_test.go`) run only with `XCP_LAB_HOST`, `XCP_PASS` and `XCP_TRUST_FINGERPRINT` (the lab
  certificate's SHA-256); the ones that write also need `XCP_ALLOW_WRITE=1` and touch only the VM named in `XCP_LAB_VM`:
  `XCP_LAB_HOST=10.77.58.11 XCP_TRUST_FINGERPRINT=... XCP_PASS=... XCP_LAB_VM=... XCP_ALLOW_WRITE=1 go test -tags nogui -run TestLab ./internal/app/`.
  The SR tests also need the lab storage (`lab.sh storage-up`): `XCP_LAB_NFS=10.77.58.1:/export` (an NFS SR and an
  NFS ISO library on `/export/iso`, removed afterwards) and `XCP_LAB_SMB='\\10.77.58.1\sr'` (an SMB SR, user `xcp`);
  `TestLabReattachSR` and `TestLabSRCommands` detach, forget, reattach, repair and destroy their own SRs only;
  New Pool is checked by hand on the lab: move the VMs off `xcp-lab-2`, `pool.eject` it (it restarts standalone),
  add it in the window and create the pool again with `xcp-lab-1` as coordinator.
  `XCP_LAB_ISCSI=10.77.58.20` formats the lab target's LUN, reattaches the SR and destroys it (`TestLabISCSI`).
- Tests use a fake JSON-RPC server over TLS (`internal/xapi/client_test.go`) and an in-memory MCP client
  (`internal/mcp/server_test.go`); none needs a real server.
- Against a server: `XCP_HOST`, `XCP_USER`, `XCP_PASS` in the environment (never in files); the lab servers and rules
  are in the repository's [AGENTS.md](../../AGENTS.md) (read-only `.236`, view-only `.114`).
- A new SDK: `tools/sdk-update.sh XenServer-SDK-<version>.zip` replaces `sdk/xenapi` and re-applies the
  `XCP-NG PATCH` changes listed in `sdk/xenapi/PATCHES.md`.
