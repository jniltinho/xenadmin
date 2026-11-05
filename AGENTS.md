# Agent guide — XCP-ng Center desktop port

The Windows app (WinForms, `XenAdmin/`) is upstream. The work here is the cross-platform desktop app:
`XenAdmin.Desktop.Core` (UI-free logic, no Avalonia, no OS APIs) and `XenAdmin.Desktop` (Avalonia). Design and
task list: `openspec/changes/add-linux-desktop-app/` (validate with `openspec validate add-linux-desktop-app --strict`).

## Build and test (containers only)

```bash
make test        # unit tests; integration tests when XCP_* are set
make test-win    # the core still builds for net481 (Windows app)
make dist        # self-contained Linux build in dist/
make appimage && make smoke
```

`tools/linux/container.sh test` forwards `XCP_HOST`, `XCP_RO_HOST`, `XCP_USER`, `XCP_PASS`, `XCP_ALLOW_WRITE`,
`XCP_LAB_HOSTS`, `XCP_LAB_GUEST_ISO`; `XCP_DOCKER_NETWORK=xcp-kvm-lab` runs it on the KVM lab's network.

## Test labs: KVM and Docker

Pools, storage, live migration and anything that needs running VMs are tested on the **KVM lab**, never on real
servers. The VirtualBox lab cannot run guests on AMD (no nested NPT; XCP-ng 8.3's Xen has no shadow paging), and
KVM cannot run while any VirtualBox VM is up.

```bash
# two nested XCP-ng 8.3 hosts, installed unattended (10.77.58.11 / .12, Docker network xcp-kvm-lab)
XCP_LAB_BACKEND=kvm XCP_LAB_PASS=... tools/vagrant/xcp-ng-lab/lab.sh up        # also: status | halt | start | destroy
XCP_LAB_BACKEND=kvm XCP_LAB_PASS=... tools/vagrant/xcp-ng-lab/lab.sh storage-up  # NFS 10.77.58.1:/export, SMB \\10.77.58.1\sr, iSCSI 10.77.58.20

# lab integration tests (pool test needs standalone hosts; storage and migration need the pool + a shared SR)
XCP_DOCKER_NETWORK=xcp-kvm-lab XCP_LAB_HOSTS=10.77.58.11 XCP_ALLOW_WRITE=1 \
XCP_LAB_GUEST_ISO=/src/_output/xcp-lab/nfs/iso/alpine-virt-3.22.7-x86_64.iso XCP_USER=root XCP_PASS=... make test
```

The Windows reference (official XCP-ng Center v25.04, to compare screens, dialogs and behaviour) runs on KVM too,
with [dockur/windows](https://github.com/dockur/windows) on the same Docker network:

```bash
tools/vagrant/win11-reference/kvm.sh up            # first run installs Windows 11 (~30 min); http://127.0.0.1:8006
tools/vagrant/win11-reference/kvm.sh shot NAME     # screenshot to _output/reference/NAME.png
tools/vagrant/win11-reference/kvm.sh ps 'command'  # PowerShell in the desktop session (ui.ps1 clicks and keys)
```

The Linux app is driven the same way in an Xvfb container (screenshots with `xwd`, input with `xdotool`) on the
`xcp-kvm-lab` network. Compare every new screen or dialog with the Windows one before calling it done.

## Rules

- Credentials live in `.env` or in variables on the command line; never in files, tests, logs or commit messages.
- `_output/` (lab disks, ISOs, answerfiles, screenshots, the Windows password) and `dist/` are never committed.
- Real servers: `192.168.137.236` is read-only; `192.168.137.114` is view-only (no input to its VMs);
  `192.168.137.110` allows writes on throwaway VMs only (`xcpng-center-linux-test*`), deleted afterwards.
- Logic goes to `XenAdmin.Desktop.Core`, ported from the WinForms code it mirrors (cite it); the UI only shows it.
  Reuse XenModel actions and `Messages` strings instead of rewriting them.
- Each milestone: tests (unit + lab), a Codex review (`codex exec -s read-only ...`), then commit and push.
- Releases: the `create-release` skill (`.claude/skills/create-release`).
