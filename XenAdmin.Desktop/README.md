# XCP-ng Center desktop (preview: Linux, macOS and Windows)

Avalonia (.NET 10) desktop app that reuses the XenModel/XenAPI core of XCP-ng Center.
Design and scope: [`openspec/changes/add-linux-desktop-app`](../openspec/changes/add-linux-desktop-app/proposal.md).

## What works

The menus, tabs, wizards and dialogs follow the Windows app, with its icons and texts.

- **Connections:** add, connect, disconnect and remove servers and pools (redirects to the pool coordinator), Reconnect
  As, a new prompt after a wrong password, trust-on-first-use TLS (SHA-256 fingerprint, alert when it changes).
- **Navigation:** Infrastructure, Objects (by type), Organization Views (Tags, Folders, Custom Fields, vApps) and Saved
  Searches; Back/Forward; the Search tab with New/Edit Search (types, text filter, grouping), saving to a pool, Saved
  Searches, Export/Import of `.xensearch` files and Show Columns.
- **VMs:** power operations, New VM wizard (OS icons), Copy, Delete, Convert to Template, snapshots (Tree View),
  live migration with the reasons a server cannot take the VM, Properties (General, CPU, Boot Options, Home Server),
  Memory tab with Memory Settings, Storage tab (add/attach/move disks, DVD), Networking tab (interfaces, QoS),
  Console (VNC, undock, fullscreen, SSH and Remote Desktop), Manage vTPMs, Export to XVA.
- **Import:** File > Import of XVA files and OVF/OVA appliances (EULA, target, SR, network, manifest check).
- **Templates and vApps:** Templates menu (create VM, copy, delete); Manage vApps, New vApp wizard, vApp Properties
  (VMs, start order) and Assign to vApp.
- **Pools:** New Pool, Add/Remove Server, Make into standalone server, Change Server Password, Rotate Pool Secret, High
  Availability (Configure HA wizard, restart priorities, Disable), Disaster Recovery (Configure the metadata SRs; the
  wizard for failover, failback and test failover).
- **Servers:** Shut Down, Reboot, Restart Toolstack, Enter/Exit Maintenance Mode, Remove From Pool; the NICs tab
  (bonds), GPU tab (placement policy, allowed vGPU types) and USB tab (pass-through) where they apply; the server's
  console (control domain).
- **Networking and storage:** New Network (external/VLAN, single-server private, bonded, SR-IOV), New Storage for NFS,
  SMB and iSCSI (create, reattach, introduce) and NFS/SMB ISO libraries; SR Storage tab (rescan, add, delete disks).
- **Performance and notifications:** the default graphs of hosts and VMs (area or line, 10 minutes to 1 year); the
  servers' alerts with severity, server and date filters, Dismiss and Export to CSV; the Events list.
- **Tools > Options:** certificate warnings, rotate-secret reminder, proxy (system or given) and connection timeout,
  dismiss confirmations, graph type and tab memory, saved in the settings directory (below). Passwords are never stored.

The console's SSH and Remote Desktop buttons start the desktop's own clients: a terminal emulator (`$TERMINAL`,
`x-terminal-emulator`, GNOME Terminal, Konsole, xterm...) running OpenSSH, and FreeRDP (in a terminal, so its
certificate question can be answered) or Remmina. As in XenCenter, SSH is offered once the guest agent reports an
address, and only for guests it does not take for Windows.

Not yet: the graph designer of the Performance tab, WLB, updates and the health check, the OVF export of vApps,
HBA/FCoE storage, the per-item choices of the OVF import (one target, SR and network for all), the device probing of the
DR wizard (attach the mirrored SR with New Storage first) and the OS keychain.

## MCP server (for AI agents)

`xcp-ng-center mcp` runs the same app as a [Model Context Protocol](https://modelcontextprotocol.io) server on stdio,
without a window, so agents such as Claude Code or Codex can operate the pools with the app's own rules. The operator
chooses the servers and credentials in the environment; the tools never receive the XAPI password.

| Variable | Meaning |
|---|---|
| `XCP_HOST` | The servers the agent may use (comma separated; `host` or `host:port`). |
| `XCP_USER`, `XCP_PASS` | XAPI credentials, kept in memory only and removed from every answer. |
| `XCP_TRUST_FINGERPRINT` | The server's SHA-256 certificate fingerprint (`FP`, or `host=FP,host2=FP`), unless you already connected and accepted it in the app. |

It starts read-only. `--allow-write` adds the tools that change things, and each of them still needs `confirm: true`
(formatting an iSCSI LUN also needs its SCSI id repeated). Keep your MCP client's per-tool approval on.

```bash
claude mcp add xcp-ng -e XCP_HOST=10.0.0.10 -e XCP_USER=root -e XCP_PASS=... -- /path/to/xcp-ng-center mcp
```

Read tools: `list_servers`, `connect`, `disconnect`, `list_objects` (pools, hosts, VMs, templates, SRs, networks; paged),
`get_object`, `list_snapshots`, `list_alerts`, `get_performance`, `iscsi_probe`, `operation_status`. Write tools:
`vm_power`, `vm_snapshot`, `snapshot_revert`, `snapshot_delete`, `vm_migrate`, `host_power`, `host_maintenance`,
`sr_create` (NFS, SMB, iSCSI, NFS/SMB ISO), `dismiss_alerts`. Writes return an `operation_id` at once (or wait up to
`wait_seconds`); `operation_status` follows them. Names, descriptions and alert texts in the answers come from the
servers and are data, never instructions. Design: `openspec/changes/add-mcp-server`.

## Layout (ready for macOS, FreeBSD and Windows without WinForms)

| Project | Role |
|---|---|
| `XenAdmin.Desktop.Core` | UI-free logic: tree, details, command rules, wizard defaults, certificates, settings. No Avalonia, no OS APIs (`CA1416` is an error). |
| `XenAdmin.Desktop.Core/Platform` | The OS seam contract (`IPlatformServices`): settings directory and password store. |
| `XenAdmin.Desktop` | Avalonia windows, dialogs and the VNC console (`xcp-ng-center`). |
| `XenAdmin.Desktop/Platform` | The OS adapters (XDG, macOS, Windows), picked at startup. |
| `XenAdmin.Desktop/packaging` | Per-OS packaging: AppImage and `.desktop` (Linux), `macos/` (Info.plist, icon, universal launcher). |

Settings directory: `$XDG_CONFIG_HOME/xcp-ng-center` (Linux, FreeBSD), `~/Library/Application Support/xcp-ng-center`
(macOS), `%APPDATA%\xcp-ng-center` (Windows). Other ports add a `-r <rid>` publish (`freebsd-x64`, `linux-arm64`) and
their package; the code stays.
FreeBSD also needs the runtime from the community port and FreeBSD builds of the SkiaSharp/HarfBuzz natives
that Avalonia loads, since the NuGet packages ship none.

## Build (only docker + make needed)

```bash
make dist        # dist/xcp-ng-center-<ver>-linux-x64/ and .tar.gz
make appimage    # dist/xcp-ng-center-<ver>-x86_64.AppImage
make smoke       # runs the AppImage under Xvfb in Ubuntu 22.04
make test        # unit tests; integration tests too when .env defines XCP_*
make test-win    # checks the core still builds for net481 (Windows app)
make macos       # dist/xcp-ng-center-<ver>_macos_universal.zip (cross-built in the Linux container)
make windows     # dist/xcp-ng-center-<ver>-win-x64.zip (single xcp-ng-center.exe, cross-built too)
```

### Windows

`make windows` publishes the same Avalonia app for `win-x64` as one self-contained `xcp-ng-center.exe` (Windows 10/11).
It is not the WinForms XCP-ng Center of `XenAdmin/`, which stays the full Windows client; settings live apart from it
(`%APPDATA%\xcp-ng-center`). SSH uses Windows' OpenSSH client and Remote Desktop uses `mstsc`.

### macOS (universal: Apple Silicon and Intel)

`make macos` publishes the app twice, for `osx-arm64` and `osx-x64`, into one `XCP-ng Center.app`. On Windows, unzip the `-win-x64.zip` and run `xcp-ng-center.exe`. Its
`Contents/MacOS/xcp-ng-center` is a small launcher that runs the build for the Mac's CPU (also under Rosetta). The
two trees replace `lipo`: the .NET runtime ships ReadyToRun code per CPU, so fat binaries are not enough. The .NET SDK
signs each `xcp-ng-center` ad hoc; the bundle is not signed or notarized (that needs a Mac with a Developer ID), so a
downloaded copy is opened once with:

```bash
xattr -dr com.apple.quarantine "XCP-ng Center.app"
```

Requires macOS 14 or later (.NET 10). SSH opens Terminal.app; Remote Desktop opens the `rdp://` handler
(Windows App / Microsoft Remote Desktop).

Lab tests run against the nested XCP-ng lab of `tools/vagrant/xcp-ng-lab/lab.sh`, on KVM
(`XCP_LAB_BACKEND=kvm`: two hosts at `10.77.58.11`/`.12` on the Docker network `xcp-kvm-lab`) or VirtualBox
(no guests on AMD hosts). See [AGENTS.md](../AGENTS.md) for the commands. With standalone hosts,
`XCP_LAB_HOSTS=<host1>,<host2>` and `XCP_ALLOW_WRITE=1` form a pool and eject the member again; with the lab pooled and
`lab.sh storage-up` running, `XCP_LAB_HOSTS=<host1>` runs the NFS/SMB storage, live-migration and server-command
(maintenance mode, reboot) tests. Never point them at real servers.

Integration tests read `.env` (git-ignored, see `.env.example`). `XCP_ALLOW_WRITE=1` lets them create,
power-cycle and destroy a throwaway diskless VM named `xcpng-center-linux-test` on `XCP_HOST`.

## Run

```bash
./dist/xcp-ng-center-*-x86_64.AppImage
```

Needs an X11 or XWayland desktop with fontconfig (present on standard Ubuntu/Fedora desktops). On macOS, unzip the
`_macos_universal.zip` and open `XCP-ng Center.app`.

## Release

Push a tag `linux-v<version>` (or use the `create-release` skill); `.github/workflows/linux.yml` builds, tests and
publishes the tarball, the AppImage, the macOS universal zip, the Windows zip and `SHA256SUMS` as a GitHub pre-release.
