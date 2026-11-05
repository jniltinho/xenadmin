# Changelog — XCP-ng Center (Go)

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions are the `go-v*` tags.

## [Unreleased]

### Added
- iSCSI in the New Storage wizard (LVMoISCSI): Scan Target Host lists the IQNs and then the LUNs (the SM's refused
  SR.create, as the C# app), CHAP optional; Finish probes the LUN and offers Reattach of the SR found on it or Format,
  never on a LUN an SR of a connected pool uses.
- The window follows the servers live, as the C# XenObjectDownloader: an `event.from` loop per connection (its own
  cancellable session copy, stopped before logout) tells the window, coalesced, to read the server again; the tree,
  the open tab and the alerts follow, and the selection stays.
- VM menu and context menu: Copy VM (fast clone or full copy to an SR, CopyVMDialog), Delete VM (the disks only it
  uses and its snapshots, ConfirmVMDeleteDialog / VMDestroyAction, every choice checked again on the server) and the
  Resume on, Migrate to and Start on Server submenus (the home server first, each server with the reason it cannot
  take the VM, from VM.assert_can_boot_here; a live VM.pool_migrate).
- XAPI failures read as in the Windows app (FriendlyErrorNames: "Cannot see required storage", not the code).
- New Storage wizard (Type, Name, Location) for NFS and SMB storage and the NFS and SMB ISO libraries, as
  SrCreateAction: the password in a secret only for SR.create, an unplugged PBD plugged or the SR forgotten, the
  pool's first shared SR made its default. Scan finds the SRs already on an NFS or SMB share and reattaches one
  (SrReattachAction / SrIntroduceAction, with the C# warning; never one another connected pool uses).
- Storage menu and SR context menu: Repair, Set as Default, Detach, Forget and Destroy, enabled as the C# app and
  with its questions; the tree shows the C# SR icons (default, broken or detached, hidden).
- New VM wizard (Template, Name, Installation Media, Home Server, CPU & Memory, Storage, Networking, Finish) with the
  C# defaults; the VM is created as CreateVMAction does and removed again if a step fails.
- The navigation pane's splitter is wider, with a grip.
- Performance tab of VMs and servers: the C# default graphs, zoom from 10 minutes to 1 year, lifecycle events, in
  the server's time; RRDs read from the member that holds them, through the coordinator's certificate list.
- The XCP-ng window icon (written as _NET_WM_ICON on X11; the .desktop file on Wayland), the C# tree expanders, a
  draggable splitter for the navigation pane.
- The main window laid out as the C# app: menus, the toolbar with its power-button rules, Back/Forward, the tree's
  search box, the navigation pane (Objects view), Home and Search tabs of the root, saved servers shown disconnected
  with Connect (password focused), Disconnect and Remove.
- VM Properties: General, CPU (vCPUs, topology, priority, hotplug and its licence) and Home Server.
- `make appimage`; every build is amd64 only.
- Storage tab: Add Virtual Disk (unique name, the SRs the VM's server can use, the default SR), Detach and Delete with
  the C# rules and questions (StorageRules.CanDetach/CanDelete, DetachVirtualDiskAction, DestroyDiskAction).
- Memory, Storage and Networking tabs of a VM (the C# Details tables): Activate/Deactivate a disk, create a DVD drive
  (with the C# notice when it needs a restart), choose or eject its ISO; lab tests in `internal/app/lab_test.go`.
- Snapshots tab: Take Snapshot (disks, or disks and memory, with the reason when not available), Revert To (optionally
  taking a snapshot first; the VM is started or resumed again when it was running), Delete (snapshots only, with the
  disks they own); MCP `list_snapshots`, `vm_snapshot`, `snapshot_revert`, `snapshot_delete`.
- Console tab: the graphical console of a running VM or of a server (its control domain) with noVNC, Send
  Ctrl+Alt+Del and Scale; `internal/console` dials the console host with the trust policy (pool members vouched by
  the coordinator) and serves it on 127.0.0.1 to the webview's origin with a one-time ticket.
- Tabs with an icon and a bold label, the selected one joined to its page with an accent line on top.
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
