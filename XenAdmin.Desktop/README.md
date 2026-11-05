# XCP-ng Center desktop (preview, Linux first)

Avalonia (.NET 10) desktop app that reuses the XenModel/XenAPI core of XCP-ng Center.
Design and scope: [`openspec/changes/add-linux-desktop-app`](../openspec/changes/add-linux-desktop-app/proposal.md).

## What works

- Add/connect/disconnect servers and pools (redirects to the pool coordinator); double-click connects.
- Trust-on-first-use TLS: SHA-256 fingerprint confirmation, alert if it changes.
- Live tree of pools, hosts, VMs and storage; General, Storage, Snapshots, Search and Console tabs.
- VM power operations, New VM wizard, Copy and Delete VM, snapshots, disks and DVD, Properties
  (General, CPU, Boot Options).
- Server list and fingerprints in the platform settings directory (below). Passwords are never stored.

The console's SSH and Remote Desktop buttons start the desktop's own clients: a terminal emulator (`$TERMINAL`,
`x-terminal-emulator`, GNOME Terminal, Konsole, xterm...) running OpenSSH, and FreeRDP (in a terminal, so its
certificate question can be answered) or Remmina. As in XenCenter, SSH is offered once the guest agent reports an
address, and only for guests it does not take for Windows.

Not yet: the graph designer of the Performance tab, Notifications, HA/WLB, updates, host actions, the other wizards, OS keychain.

## Layout (ready for macOS, FreeBSD and Windows without WinForms)

| Project | Role |
|---|---|
| `XenAdmin.Desktop.Core` | UI-free logic: tree, details, command rules, wizard defaults, certificates, settings. No Avalonia, no OS APIs (`CA1416` is an error). |
| `XenAdmin.Desktop.Core/Platform` | The OS seam contract (`IPlatformServices`): settings directory and password store. |
| `XenAdmin.Desktop` | Avalonia windows, dialogs and the VNC console (`xcp-ng-center`). |
| `XenAdmin.Desktop/Platform` | The OS adapters (XDG, macOS, Windows), picked at startup. |
| `XenAdmin.Desktop/packaging` | Per-OS packaging (AppImage today). |

Settings directory: `$XDG_CONFIG_HOME/xcp-ng-center` (Linux, FreeBSD), `~/Library/Application Support/xcp-ng-center`
(macOS), `%APPDATA%\xcp-ng-center` (Windows). macOS and Windows ports add a `-r <rid>` publish (`osx-arm64`, `win-x64`) and their package; the code stays.
FreeBSD also needs the runtime from the community port and FreeBSD builds of the SkiaSharp/HarfBuzz natives
that Avalonia loads, since the NuGet packages ship none.

## Build (only docker + make needed)

```bash
make dist        # dist/xcp-ng-center-<ver>-linux-x64/ and .tar.gz
make appimage    # dist/xcp-ng-center-<ver>-x86_64.AppImage
make smoke       # runs the AppImage under Xvfb in Ubuntu 22.04
make test        # unit tests; integration tests too when .env defines XCP_*
make test-win    # checks the core still builds for net481 (Windows app)
```

`XCP_LAB_HOSTS=192.168.56.11,192.168.56.12` (with `XCP_ALLOW_WRITE=1`) runs the pool test against the nested lab of
`tools/vagrant/xcp-ng-lab/lab.sh` (two XCP-ng hosts in VirtualBox, installed unattended); it forms a pool and ejects
the member again. Never point it at real servers.

`XCP_LAB_HOSTS=192.168.56.11` with the lab already pooled runs the storage tests instead: `lab.sh storage-up` serves
NFS (`192.168.56.1:/export`) and SMB (`\\192.168.56.1\sr`, user `xcp`) from a container bound to the host-only address,
and the tests create, scan, detach, reattach and destroy a shared SR of each kind. The VirtualBox lab cannot run guests
on AMD (no nested NPT, and XCP-ng 8.3's Xen has no shadow paging), so the migration test is inconclusive there.

Integration tests read `.env` (git-ignored, see `.env.example`). `XCP_ALLOW_WRITE=1` lets them create,
power-cycle and destroy a throwaway diskless VM named `xcpng-center-linux-test` on `XCP_HOST`.

## Run

```bash
./dist/xcp-ng-center-*-x86_64.AppImage
```

Needs an X11 or XWayland desktop with fontconfig (present on standard Ubuntu/Fedora desktops).

## Release

Push a tag `linux-v<version>`; `.github/workflows/linux.yml` builds, tests and publishes the tarball,
AppImage and `SHA256SUMS` as a GitHub pre-release.
