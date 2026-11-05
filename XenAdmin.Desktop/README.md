# XCP-ng Center desktop (preview: Linux, macOS and Windows)

Avalonia (.NET 10) desktop app that reuses the XenModel/XenAPI core of XCP-ng Center.
Design and scope: [`openspec/changes/add-linux-desktop-app`](../openspec/changes/add-linux-desktop-app/proposal.md).

## What works

- Add/connect/disconnect servers and pools (redirects to the pool coordinator); double-click connects.
- Trust-on-first-use TLS: SHA-256 fingerprint confirmation, alert if it changes.
- Live tree of pools, hosts, VMs and storage; General, Storage, Snapshots, Search and Console tabs.
- VM power operations, New VM wizard, Copy and Delete VM, snapshots, disks and DVD, Properties
  (General, CPU, Boot Options).
- Pools and storage: New Pool, New Storage Repository for NFS, SMB and iSCSI (create, reattach, introduce) and NFS/SMB
  ISO libraries, live migration
  between the servers of a pool (VM > Migrate to Server, with the reason a server cannot take the VM).
- Servers: Shut Down, Reboot and Enter/Exit Maintenance Mode (evacuates the VMs and offers to bring them back).
- Performance tab: the default graphs of hosts and VMs (CPU, memory, network, disk), zoom from 10 minutes to 1 year,
  VM lifecycle events.
- Notifications: the servers' alerts (severity, text, server/pool, date) with severity and server filters, Dismiss and
  Dismiss All, and the alert count on the Notifications button and the status bar; the Events list of this session.
- Server list and fingerprints in the platform settings directory (below). Passwords are never stored.

The console's SSH and Remote Desktop buttons start the desktop's own clients: a terminal emulator (`$TERMINAL`,
`x-terminal-emulator`, GNOME Terminal, Konsole, xterm...) running OpenSSH, and FreeRDP (in a terminal, so its
certificate question can be answered) or Remmina. As in XenCenter, SSH is offered once the guest agent reports an
address, and only for guests it does not take for Windows.

Not yet: the graph designer of the Performance tab, Objects/Organization Views/Saved Searches, the date filter and
Export of alerts, HA/WLB, updates, import/export, the other SR types and wizards, NICs/GPU/USB tabs, OS keychain.

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
