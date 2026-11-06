# XCP-ng Center (Go)

One binary, `xcp-ng-center-go`, for the desktop app, the command line and an MCP server for AI agents. It follows the
proposal in [`openspec/changes/add-go-desktop-app`](../openspec/changes/add-go-desktop-app/proposal.md): Go, the
XenServer Go SDK, Cobra, the official MCP Go SDK, and Wails v3 with Vue and Tailwind for the window, built for
Linux, Windows and macOS and ready to become a web app.

The window follows the Windows app (XCP-ng Center 25.04) screen by screen: its menus, toolbar, tree, tabs, dialogs
and wizards, with its texts (generated from the C# `.resx` files), its icons and its rules, ported from the C# code
they mirror. Every feature below was checked against XCP-ng 8.3 on a nested KVM lab.

## Download

The [releases](https://github.com/jniltinho/xenadmin/releases) (`go-v*`) hold, for amd64:

| File | What |
| --- | --- |
| `xcp-ng-center-go-X.Y.Z-x86_64.AppImage` | Linux desktop app; needs the system's GTK 4 and WebKitGTK 6.0 (tested on Ubuntu 24.04); `chmod +x` and run |
| `xcp-ng-center-go-X.Y.Z-windows-amd64-desktop.exe` | Windows 10/11 desktop app (WebView2) |
| `xcp-ng-center-go-X.Y.Z_macos_amd64.zip` | macOS app, cross-built and unsigned (see [macOS](#macos)) |
| `xcp-ng-center-go-X.Y.Z-{linux,windows,darwin}-amd64` | the CLI and MCP server alone (no window), static |
| `SHA256SUMS` | the checksums of the files above |

## What the window does

**Connection.** Add New Server with `host[:port]` or IPv6, TLS 1.2+, the certificate question on first use (SHA-256
pins; a changed certificate is refused), the pool coordinator redirect, saved servers shown disconnected (Connect,
Disconnect, Remove). Passwords stay in memory, never in files, logs or errors. The window follows the servers live
(`event.from`): what another client changes shows up at once.

**Main window.** The C# menus, toolbar (Back, Forward, Add New Server, New Pool, New Storage, New VM, the power
buttons, Pause/Unpause) and context menus, enabled by the same rules; the tree with its search box and the C#
icons; the navigation pane (Infrastructure, Objects, Organization Views, Saved Searches, Notifications), resizable;
the Home and Search tabs (the overview with live CPU, memory, disk, network, address and uptime columns); XAPI errors in the Windows app's words (`FriendlyErrorNames`).

**VMs.** Start, Shut Down, Reboot, Suspend, Resume, Pause, Unpause and the forced ones, with the C# questions;
Export to an XVA file and Import (File > Import...) of an XVA file or an OVF/OVA appliance (.ovf, .ova, .ova.gz, .ova.bz2; raw, VHD and stream-optimized VMDK disks; its EULA, manifest verification, the server of each VM, the SR of each disk and the network of each interface); Resume on, Migrate to and Start on Server (live migration, each server with the reason it cannot take the VM);
User templates in the tree (Create VM From Selection, Delete Template); New VM wizard (template, installation media, home server, CPU and memory, disks, networks); Convert to Template; Copy VM (fast clone or
full copy to an SR that can hold it); Delete VM (with the disks only it uses and its snapshots); Properties
(General, Custom Fields, CPU, Boot Options, Start Options, Alerts, Home Server, Advanced); the General (the C# sections: General, Boot
Options, CPUs, Read Caching; for servers, pools and SRs Certificates, Version Details, Management Interfaces, Memory,
CPUs, Status, Multipathing), Memory (Memory Settings: fixed, range or static
maximum), Storage (add, attach, activate, detach and delete disks; the DVD drive and its ISO), Networking (add,
change, remove, activate and deactivate interfaces), Console (noVNC), Performance and Snapshots (take, revert,
delete) tabs.

**Servers and pools.** Shut Down, Reboot, Enter/Exit Maintenance Mode (evacuation, restoring the VMs), Restart
Toolstack, Remove From Pool, the Pool menu (Add Server, Remove Server, Rotate Pool Secret, Make into standalone
server, Manage vApps, High Availability: the HA wizard with the heartbeat SR, the VMs' restart priorities, start
order and delay and the failures to tolerate, the restart priorities of a pool with HA on, and Disable; Disaster
Recovery: Configure (the SRs holding the pool's metadata) and the wizard (failover, failback and test failover from the
metadata on SRs attached to the pool, with the prechecks and the clean-up of a test)), the Memory
(with the bars of each server's memory), Storage, Networking, NICs, GPU (placement policy and allowed vGPU types), USB
(pass-through), Console and Performance tabs; on Networking the New Network wizard
(External/VLAN, Single-Server Private and Bonded networks), network Properties (name, description, auto-add, MTU, NIC and VLAN, bond mode and LACP hashing) and
Remove Network; on NICs Create Bond (active-active, active-passive, LACP), Delete Bond and Rescan, with the C#
management interface and HA warnings; New Pool (the coordinator and the members with the reason each cannot join, the CPU levelling
question, `pool.join`); a server's Properties with Alerts (the perfmon thresholds), Multipathing and Log
Destination, an SR's with Alerts.

**Storage.** The New Storage wizard: NFS, SMB, iSCSI (Scan Target Host, IQN, LUN, Reattach or Format), NFS and SMB
ISO libraries, Scan and Reattach of an existing SR; the SR's Storage tab (its disks: Rescan, Add, Move to another SR
live or offline, Delete, Properties); a server's or pool's Storage tab (New SR, Reclaim freed space, Properties);
Repair, Set as Default, Detach, Forget and Destroy, with the C# icons of a default, broken or detached SR.

**Views and notifications.** Organization Views (Tags, Folders, Custom Fields, vApps), the default Saved Searches
and the searches saved in the pools (the XenSearch filters and groupings they use most) over every connected server;
the Search tab's New Search, Edit Search (types, a text filter, a grouping; Save to a pool, Delete), Export... and
Import... of `.xensearch` files the Windows app reads; searches by name, type, power state, OS, tags, folder, custom
fields, dates, IP addresses and the pool, server, vApp or network of an object; Notifications > Alerts with the
severity, server and date filters, Export (CSV), Dismiss, Dismiss Selected and Dismiss All.

Not there yet: a VM's GPU and USB Properties pages, SR-IOV networks, the OS fixup,
encrypted packages and VMDKs that are not stream-optimized when importing an OVF/OVA appliance, OVF/OVA export, and saved searches that filter on a related object's
properties (XenSearch Recursive queries; listed, not run). The Windows app has those.

## Command line and MCP

```bash
export XCP_HOST=192.168.0.10 XCP_USER=root     # or --host / --user
read -rs XCP_PASS && export XCP_PASS           # or type it when asked (never a flag)
xcp-ng-center-go server add                    # asks to trust the certificate on first use
xcp-ng-center-go vm list --json                # pool|host|vm|template|sr|network list, table or --json
xcp-ng-center-go mcp                           # MCP server on stdio, read-only
xcp-ng-center-go mcp --allow-write             # also the tools that change the servers (each needs confirm: true)
```

The MCP tools have the names and arguments of the C# app's MCP server: `list_servers`, `connect`, `disconnect`,
`list_objects`, `list_alerts`, `list_snapshots`, and with `--allow-write` `vm_power`, `vm_snapshot`,
`snapshot_revert`, `snapshot_delete`, `host_power`, `host_maintenance`, `dismiss_alerts`. The password comes from
`XCP_PASS`; the agent never sees it.

Settings (servers, certificate fingerprints and Tools > Options; never passwords) live in `settings.toml` under the user config
directory (`~/.config/xcp-ng-center-go` on Linux), mode 0600, shared by the window and the CLI.

## Build (only docker and make)

```bash
make check     # gofmt, go vet (app and SDK), tests
make cli       # static CLI/MCP binaries for Linux, macOS and Windows in dist/
make linux     # the desktop app for Linux (node:24 + Ubuntu 24.04 GTK/WebKit images)
make appimage  # the desktop app as an x86_64 AppImage
make windows   # the desktop app for Windows (WebView2; no CGO)
make macos     # "XCP-ng Center (Go).app" for macOS, cross-built, zipped
make dist      # all of the above
```

Every build is amd64. The version comes from the `go-v*` git tag.

### macOS

`make macos` cross-builds the app on Linux with Wails' CGO-free backend (`-tags purego`, experimental in Wails 3.1).
It is unsigned: after downloading, `xattr -dr com.apple.quarantine "XCP-ng Center (Go).app"`; Apple Silicon Macs run
it with Rosetta 2.

For a Mac release, build natively on a Mac (macOS 14+, Go 1.27, Node 24, Xcode command line tools):

```bash
build/darwin/build-on-mac.sh                                   # this Mac's architecture (arm64 or amd64)
ARCH=amd64 build/darwin/build-on-mac.sh                        # the other one
SIGN_IDENTITY="Developer ID Application: ..." NOTARY_PROFILE=xcp-ng build/darwin/build-on-mac.sh
```

It builds the frontend, the app with the native (CGO) backend, the same `.app` bundle (`build/darwin/Info.plist`,
the XCP-ng icon) and its zip in `dist/`; with `SIGN_IDENTITY` it signs with the hardened runtime, and with
`NOTARY_PROFILE` (`xcrun notarytool store-credentials`) it notarizes and staples.

## More

- [docs/architecture.md](docs/architecture.md): the packages, the core/UI split, live updates.
- [docs/development.md](docs/development.md): the KVM lab, the lab tests and their variables, the UI checks.
- [CHANGELOG.md](CHANGELOG.md).
