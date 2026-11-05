# Changelog — XCP-ng Center desktop (Linux)

All notable changes to the cross-platform desktop app (`XenAdmin.Desktop`). Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
versions are the `linux-v*` git tags. The Windows (WinForms) app keeps its own `CHANGELOG.md` at the repository root.

## [Unreleased]

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
- Nested XCP-ng lab in VirtualBox (`tools/vagrant/xcp-ng-lab`) for pool and write tests.

### Changed
- Projects reorganised for future macOS, FreeBSD and Windows (non-WinForms) ports: `XenAdmin.Desktop.Core` (UI-free logic,
  platform seam) and `XenAdmin.Desktop` (Avalonia).
- Settings directory follows the platform (XDG, `~/Library/Application Support`, `%APPDATA%`).
- One TLS trust policy for XenAPI and the console tunnel; pool members trusted through the coordinator.

### Fixed
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
