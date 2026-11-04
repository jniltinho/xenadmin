# XCP-ng Center for Linux (preview)

Avalonia (.NET 10) desktop app that reuses the XenModel/XenAPI core of XCP-ng Center.
Design and scope: [`openspec/changes/add-linux-desktop-app`](../openspec/changes/add-linux-desktop-app/proposal.md).

## What works (MVP)

- Add/connect/disconnect servers and pools (redirects to the pool coordinator).
- Trust-on-first-use TLS: SHA-256 fingerprint confirmation, alert if it changes.
- Live tree of pools, hosts, VMs and storage, with search and a properties panel.
- VM power operations (start, shutdown, reboot, suspend, resume, force variants) with a task list.
- Server list and fingerprints in `~/.config/xcp-ng-center/config.json`. Passwords are never stored.

Not yet: VM console (VNC), wizards, HA/WLB, updates, Secret Service password storage.

## Build (only docker + make needed)

```bash
make dist        # dist/xcp-ng-center-<ver>-linux-x64/ and .tar.gz
make appimage    # dist/xcp-ng-center-<ver>-x86_64.AppImage
make smoke       # runs the AppImage under Xvfb in Ubuntu 22.04
make test        # unit tests; integration tests too when .env defines XCP_*
make test-win    # checks the core still builds for net481 (Windows app)
```

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
