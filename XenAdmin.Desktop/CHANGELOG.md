# Changelog — XCP-ng Center desktop (Linux, macOS, Windows)

All notable changes to the cross-platform desktop app (`XenAdmin.Desktop`). Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
versions are the `linux-v*` git tags. The Windows (WinForms) app keeps its own `CHANGELOG.md` at the repository root.

## [Unreleased]

### Added
- VM Networking tab: Add Interface, Properties (network, MAC, QoS limit), Remove Interface and Activate/Deactivate,
  with the rules, tooltips and confirmations of the Windows app, and the Limit column.
- SR Storage tab: Rescan, Add (a disk on that SR), Delete and Properties with the Windows rules and confirmations;
  the Storage Volume and Changed Block Tracking columns when they apply.
- Host and pool Storage tab: New SR, Reclaim freed space (with confirmation) and Properties.
- The storage and network grids select rows with the keyboard (Up/Down, Enter) and open Properties on double click.
- Host and pool Networking tab: Add Network (New Network wizard: External network on a VLAN and Single-Server
  Private network; Bonded and SR-IOV listed, not yet available), Properties and Remove Network, with the Windows
  columns (NIC, VLAN, Auto, Link Status, MAC, MTU, SR-IOV).
- Host NICs tab: the physical interfaces (link, speed, duplex, vendor, device, PCI bus path, FCoE and SR-IOV) and Rescan.
- Server menu: Restart Toolstack and Remove From Pool (also in the server's context menu); VM menu: Take a Snapshot
  and Convert to Template, with the Windows rules and confirmations.
- Notifications > Alerts: Filter by Date (last 24 hours, 3, 7 or 30 days, or a custom start and end) and Export to a
  CSV file (all alerts or only the filtered ones).
- Move Virtual Disk from the VM and SR Storage tabs: live storage migration for attached disks, an offline move
  otherwise, with the reason each storage repository cannot take the disk, and Rescan.
- Bonded networks: the New Network wizard's Bonded Network (members, Active-active / Active-passive / LACP, MTU),
  Create Bond and Delete Bond in the host's NICs tab, and Remove Network of a bond, with the warnings for the
  management interfaces.
- Lab: every KVM lab host gets two spare NICs on an isolated bridge (bonds and VLANs), and `lab.sh halt` shuts the
  hosts down before stopping them.
- New VM: the templates show their OS logos and categories and follow the Windows order (custom and Windows templates
  first, newest version first); the list no longer starts with an empty band.
- Console: Undock (Alt+Shift+U) into its own window, Redock, and Fullscreen (Ctrl+Enter); while undocked the tab
  offers Find Console and Reattach Console.
- Memory tab: the memory bar of a VM (in use, dynamic range) or of a server (Xen and control domain, each VM, free),
  and Edit: the Memory Settings dialog (a fixed size or a range, or the advanced dynamic minimum, maximum and static
  maximum), with the server's limits and the restart warning.
- Pool menu: Add Server (the connected standalone servers that may join), Remove Server, Make into standalone server,
  Reconnect As, Disconnect, Change Server Password and Rotate Pool Secret; Server menu: Change Server
  Password (current password checked, the passwords kept in memory only).
- Templates menu: Create VM From Selection (New VM wizard on the selected template), Copy, Delete Template and
  Properties.

### Fixed
- The Networking tab buttons follow the Windows order (Add, Properties, Remove, Activate/Deactivate).
- Delete Virtual Disk also requires the server to allow `destroy` on the disk, as the Windows app does.

## [0.2.0] — 2026-11-05

**XCP-ng Center layout and tabs, wizards, storage (NFS, SMB, iSCSI, ISO libraries), migration, maintenance mode,
Performance and Notifications; universal macOS and Windows x64 builds; an MCP server for AI agents.**

### Added
- XCP-ng Center layout: menus, toolbar, navigation tree with the original icons and events (double-click connects or starts).
- VNC console tab, VM General/Storage/Snapshots/Search tabs, tab remembered per object.
- Properties dialog (General, CPU, Boot Options) with WinForms-style controls.
- VM Storage actions: new, attach, activate, deactivate, detach and delete disks; DVD drive and ISO.
- New VM wizard, Copy VM (fast clone / full copy) and Delete VM.
- Create New Pool dialog (coordinator, members with the join rules of the Windows app).
- New Storage Repository wizard for NFS and SMB (Scan, create, reattach or introduce an existing SR).
- VM > Migrate to Server inside a pool, with the reason each server cannot take the VM.
- Server commands: Shut Down, Reboot and Enter/Exit Maintenance Mode (evacuates the VMs and offers to bring them
  back), from the toolbar, the Server menu and the tree's context menu.
- Performance tab: the default graphs of hosts and VMs (CPU, used memory, network, disk) from the RRD archives,
  zoom from 10 minutes to 1 year, and the VM lifecycle events.
- Console: Launch SSH Console (the desktop's terminal with OpenSSH) and Turn on / Switch to Remote Desktop (the
  desktop's RDP client: FreeRDP in a terminal or Remmina on Linux, mstsc on Windows, the RDP app on macOS).
- Notifications > Alerts: the servers' alerts with severity, expandable text, server/pool and date; filters by
  severity and server; Dismiss, Dismiss All and Dismiss Selected; the count on the Notifications button and in the
  status bar.
- Properties > Home Server: no home server (with shared storage) or one of the live servers of the pool.
- macOS: one universal app for Apple Silicon and Intel (`make macos`, `xcp-ng-center-<ver>_macos_universal.zip`),
  built in the Linux container and attached to the releases.
- Windows: the same Avalonia app for Windows 10/11 x64 as one `xcp-ng-center.exe` (`make windows`,
  `xcp-ng-center-<ver>-win-x64.zip`), next to the WinForms XCP-ng Center.
- New Storage Repository: iSCSI (target host, CHAP, Scan Target Host, target IQN and LUN; Format, or Reattach the SR
  found on the LUN) and the NFS ISO and SMB ISO libraries.
- Lab: an iSCSI target (tgt) next to the NFS and SMB servers (`lab.sh storage-up`).
- MCP server: `xcp-ng-center mcp [--allow-write]` lets AI agents (Claude Code, Codex...) list and operate pools, hosts,
  VMs, snapshots, storage and alerts with the app's rules; servers and credentials come from the operator's
  environment, certificates must already be trusted, and writes need `--allow-write` and `confirm: true`.
- Nested XCP-ng lab in VirtualBox (`tools/vagrant/xcp-ng-lab`) for pool and write tests.

### Changed
- Projects reorganised for future macOS, FreeBSD and Windows (non-WinForms) ports: `XenAdmin.Desktop.Core` (UI-free logic,
  platform seam) and `XenAdmin.Desktop` (Avalonia).
- Settings directory follows the platform (XDG, `~/Library/Application Support`, `%APPDATA%`).
- One TLS trust policy for XenAPI and the console tunnel; pool members trusted through the coordinator.

### Fixed
- New VM: the Storage page can add and delete disks (a template without disks no longer gives a diskless VM).
- Shut Down, Reboot and Suspend of a VM ask for confirmation, and Enter Maintenance Mode lists the VMs it migrates,
  as in the Windows app.
- Storage repositories: Set as Default, Detach, Forget and Destroy in the context menu, with the Windows rules and
  confirmations.
- Console: the DVD band shows the VM's drive and its ISO picker (or creates a drive).
- Boot Options: clicking a device selects it for Move Up/Down; Esc closes Properties and the other dialogs; the Events
  list lays out correctly; a VM's IPv4 address is listed once; no doubled separators in context menus.
- Alpine guests count as Linux (XenModel), so the console offers Launch SSH Console for them.
- XenAPI.HTTP refuses a redirect from HTTPS to plain HTTP.
- Connections on ports 80/8080 (plain HTTP) are refused, so passwords never cross the network unencrypted; IPv6
  addresses (bare or `[addr]:port`) are parsed correctly.
- The console key log (`XCP_CONSOLE_KEYLOG`) exists only in debug builds.
- The DVD drive's Eject button follows an ISO inserted elsewhere.
- The main window is maximized on screens smaller than its default size.
- Leaving the Console tab no longer drops the connection to the pool, and coming back to it reconnects
  (it stayed at "Connecting..." with the old picture).
- Console keyboard: HVM guests get scancodes, so symbols such as % and layouts such as ABNT2 work.
- Full copy no longer falls back to a fast clone when no storage is selected.
- Hidden storage repositories are never picked as the default for new disks.

## [0.1.0] — 2026-11-04

**First preview: connect to servers and pools, live resource tree and VM power operations.**

### Added
- Add, connect and disconnect servers and pools; trust-on-first-use certificates.
- Live tree of pools, hosts, VMs and storage with a details panel.
- VM start, shutdown, reboot, suspend and resume, with their force variants.
- Self-contained Linux build, AppImage and GitHub release workflow.
