# Changelog — XCP-ng Center (Go)

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions are the `go-v*` tags.

## [Unreleased]

### Added

- A VLAN on an SR-IOV network (NetWDetails' checkBoxSriov, NewNetworkWizard.CreateNonBonded): on a NIC with an
  SR-IOV network the External network's Interface page offers "Create the VLAN on the SR-IOV network", which checks
  the VLAN against those already on that network and makes it on its logical PIF. Unit-tested; not run on hardware.
- OVF/OVA import of hosted sparse VMDKs (monolithicSparse, what `qemu-img convert -O vmdk` and VMware Workstation
  write) and monolithic flat VMDKs (a descriptor and its flat extent), converted to raw as they are read, as DiscUtils
  does in the Windows app. A hosted sparse disk needs a folder package or an uncompressed OVA.
- Signed OVF/OVA packages: the Export wizard's "Sign the OVF package" with a .pfx certificate and its private key's
  password (never kept), writing the `.cert` file as the Windows app does; the import's "Verify digital signature"
  with View Certificate (subject, issuer, validity, thumbprint), which refuses an untrusted certificate or a bad
  signature before the manifest is checked.

## [0.3.0] — 2026-11-06

**The last parity gaps with the Windows app: OVF/OVA export, the import's per-VM, disk and interface choices, a VM's
GPU and USB pages, SR-IOV networks and searches on related objects.**

### Added

- Searches on a related object's properties (XenSearch RecursiveXMOPropertyQuery and RecursiveXMOListPropertyQuery:
  the VMs on an SR or a server, the servers running a VM, the VMs of a disk or a pool), run on every object's facts.
- A VM's GPU page in its Properties (GpuEditPage): its vGPUs with their PCI slot, type, vGPUs per GPU and video RAM,
  Add (the types each GPU group allows beside the VM's other vGPUs, whole-GPU pass-through only when vGPUs are
  restricted) and Delete, with the C# warnings (Remote Desktop, drivers, a VM that is running or protected by HA,
  multiple vGPU support, Tech Preview types); shown when the VM can have a GPU and the pool has one. Saved on OK.
- A VM's USB page in its Properties (USBEditPage): its USB devices, Attach (the free pass-through devices of the
  servers the VM can start on; the first one moves the VM to the upstream QEMU) and Detach; for an HVM VM when a
  server has a USB device and the licence allows it.
- SR-IOV networks: the New Network wizard's SR-IOV Network (a NIC that can enable SR-IOV, the disruption question,
  then SR-IOV on that device of every server, the coordinator first), offered only with an SR-IOV capable NIC; the
  NICs tab's SR-IOV Capable column (virtual functions left, a reboot needed, an unplugged logical PIF). The vGPU, USB
  and SR-IOV writes were not run on hardware (the lab has none).
- OVF/OVA export as the Windows app's Export wizard (ExportApplianceAction): VM > Export... offers an OVF/OVA package
  besides the XVA file, and Manage vApps' Export exports a vApp; the name, location and format, the halted or
  suspended VMs of the pool (the selected VM or the vApp's VMs ticked, the vTPM warning), EULA files, a manifest, a
  single .ova file, compression (the disks as .vhd.gz, or the OVA as .ova.gz) and Verify. The disks come as VHD from
  export_raw_vdi; the descriptor is XenOvf's (the Xen settings, the start order and delays) and the appliance import
  reads it back. The package is made beside the target and only renamed into place once complete.
- The import wizard's per-VM choices when importing an OVF/OVA appliance (ImportSelectHostPage, ImportSelectStoragePage,
  ImportSelectNetworkPage, VmMapping): the server of each VM (to put its disks on that server's local storage), all
  the disks on one SR or each disk on its own, and the network of each interface; the defaults keep it one click. The
  choices are checked against the package and the pool when the dialog closes and again by the import (an SR the
  server cannot see, a full one, one too small for the disks put on it, a network the server does not see).


## [0.2.0] — 2026-11-06

**Close to the Windows app's feature set: VM, server, pool, SR and network Properties and tabs, HA, disaster
recovery, vApps, XVA and OVF/OVA import, the search editor and its filters, the Performance graph editor, alert
filters and export, and Tools > Options.**

### Added

- Searches by address (XenSearch IPAddressQuery: a VM's guest addresses, a server's interfaces, full or partial such
  as 10.77.58.*), also from the search editor's text filter.
- Searches by a VM's start time (XenSearch DatePropertyQuery: today, yesterday, this or last week, before, after or
  on a day).
- Searches by custom field (XenSearch CustomFieldQuery for text fields, CustomFieldDateQuery for date fields).
- Searches by the pool, server, vApp or network an object belongs to (XenSearch XenModelObjectPropertyQuery and
  XenModelObjectListContains(Name)Query).

- The Alerts page of a VM's, server's and SR's Properties (CPU, network, disk, free memory, control domain memory and
  storage throughput alerts with their repeat interval, in the perfmon other_config key), then perfmon's refresh.

- Network Properties, Network Settings page complete: the NIC and VLAN of a VLAN network, a private network made
  external on a NIC and back (the network's PIFs destroyed and its VLAN made again on every server), and a bond
  network's mode (active-active, active-passive, LACP) and LACP load balancing; fixed, with the C# reason, on a
  management interface or with VMs attached, and the disruption warning shown first.
- Performance tab: the graph editor of the Windows app — Actions > New Graph, Edit Graph, Delete Graph, Restore Default
  Graphs, Move Up / Move Down (also in the graphs' context menu; a double click edits) and Graph Details (the graph's
  name; its data sources with their type, recording state and colour, the hidden and disabled ones on demand, Enable).
  The layout and colours are kept in the pool's gui_config where the Windows app keeps them, so both apps show the
  same graphs for a VM or a server; a non-standard data source put on a graph is recorded. Each change checks the
  role first and is refused when the graphs changed meanwhile.
- Server Properties: the Multipathing page (editable in maintenance mode; the SRs are unplugged and plugged again
  around the change) and the Log Destination page (a remote syslog server, then syslog_reconfigure).

- Memory tab bars: a VM's memory in use inside its dynamic range, and each server's memory split into Xen and the
  control domain, each VM and free memory (MemoryView).
- GPU tab of pools and servers (placement policy, GPUs by settings, allowed vGPU types and their editor) and USB tab of
  servers with USB devices (Enable / Disable pass-through), shown where the C# app shows them.
- Search tab: New Search and Edit Search (name, the object types, a name/description text filter and a grouping, as
  the C# Desktop editor; an edited search keeps its conditions, other types, nested grouping, columns and sorting),
  Apply, Save to a pool (gui_config `xensearch.Search-<uuid>`, it joins Saved Searches), Delete (from every pool), and
  Export... / Import... of `.xensearch` files in the SearchMarshalling format the Windows app loads.
- File > Import... of OVF/OVA appliances (.ovf, .ova, .ova.gz, .ova.bz2) as the C# app (OvfRules,
  ImportApplianceAction): the package's VMs, disks and interfaces, its EULA to accept, the pool or server, one SR for
  every disk, one network for every interface, manifest verification and Start; a vApp named after the package with
  each VM's vCPUs, memory, Xen settings, disks, DVD drive and interfaces. Raw, VHD and stream-optimized VMDK disks
  (converted to raw as they are read), streamed from the archive without extracting it; a failed import removes what
  it made. The package is checked as untrusted input (file names inside it, sizes, no DOCTYPE).
- General tab of VMs, templates, servers, pools, SRs and networks as in the C# app (Details.General): collapsible
  sections with Expand all / Collapse all and the Properties button; General (tags, folder, OS, virtualization mode and
  state, uptime, home server, pool coordinator, enabled, autoboot, certificate verification, iSCSI IQN, log
  destination, type, size, SCSI ID), Boot Options, CPUs, Read Caching, Certificates, Version Details, Management
  Interfaces, Memory, Status and Multipathing, with the C# labels and formats.
- SR Storage tab: Rescan, Add... (a disk on the SR attached to no VM), Move... (live storage migration of a disk in
  use, else an offline move of a disk of halted VMs, each SR with its free space or the reason it cannot take the
  disk, and Rescan), Delete (the C# question; a snapshot disk deletes its snapshots) and Properties (name,
  description, folder, tags, custom fields; also a double click). Servers and pools, Storage tab: New SR... (the New
  Storage wizard), Reclaim freed space (with its question) and Properties. Each change checks the role first.
- Servers and pools, Networking tab: Add Network... (the New Network wizard: External network on a VLAN of a NIC,
  Single-Server Private network, Bonded network), Properties (name, description, auto-add, MTU with the C# reasons it
  is fixed; a new MTU unplugs and plugs the network's PIFs and VIFs) and Remove Network (its VLAN PIFs first; a bond
  network deletes the bond).
- Servers, NICs tab: Create Bond... (the NICs up to the pool's limit, active-active, active-passive or LACP, MTU,
  auto-add; the bond on every server, the secondary management interfaces moved onto it and back on a failure),
  Delete Bond (on every server) and Rescan, with the C# warnings (management interfaces, HA) and an RBAC check first.
- Searches saved in pools run, for the XenSearch filters and groupings they use most (GroupQuery, string, enum,
  boolean, integer, tag and null queries; type, power state, OS, tags, tools, folder, custom field, pool, server,
  network and vApp groupings), with the C# semantics; a search with another filter is still listed, not run.
- VM Properties: Boot Options (autostart, the HVM boot order, the PV boot device and parameters) and, for HVM VMs,
  Advanced (the shadow memory multiplier, changed live on a running VM).
- Pool > High Availability: Configure... (with HA off the HA wizard: the heartbeat SR, each VM's restart priority,
  start order and delay, a VM that is not agile only restarted if possible, and the failures to tolerate up to the
  maximum the pool computes for the priorities, then EnableHAAction; with HA on the restart priorities and the
  failures to tolerate, SetHaPrioritiesAction) and Disable, with the C# questions and an RBAC check first.
- Pool > Disaster Recovery: Configure... (the SRs that can hold the pool's VM metadata, at most 8, with the C#
  warnings; database replication on or off where the choice changed) and the Disaster Recovery Wizard (Failover,
  Failback and Test Failover: the SRs with metadata, the pools, vApps and VMs found in it with the current pool
  marked, the prechecks, the recovery and its start mode, and the clean-up of a test failover), an RBAC check first.
  The wizard reads SRs already attached to the pool, as the C# Desktop app (no iSCSI/HBA probe of mirrored LUNs).
- VM Properties: Start Options (the start order and the delay before the next VM starts; on an HA pool the restart
  priority too, saved with the pool's failures to tolerate).
- VM Properties: Custom Fields (the pool's fields with the VM's values, dates as the C# app stores them) and Edit
  Custom Fields (add a text or date field to the pool, delete one with its values on every object).
- Memory tab: Edit... opens Memory Settings (a fixed memory or a range, or the dynamic minimum, maximum and static
  maximum with ballooning), with the C# limits and questions; a new static maximum restarts a running VM on its server.
- Servers and pools: the Memory, Storage, Networking and NICs (a server) tabs; SRs: the Storage tab (its disks and
  the VMs using them).
- Pool menu: Add Server and Remove Server (submenus), Disconnect, Rotate Pool Secret and Make into standalone server,
  with the C# questions and warnings (supplemental packs and CPU levelling, also in New Pool).
- VM Properties, General: the folder and tags (Change Folder, Edit Tags), as for the other objects.
- Properties of servers, pools, SRs and networks (General: name, description, folder, tags, a server's iSCSI IQN;
  Custom Fields), with Change Folder and Edit Tags; a folder left empty is kept (EMPTY_FOLDERS).
- VM > Manage vTPMs (VtpmManagementDialog): add and remove a UEFI VM's vTPMs, with Protected and Unique.
- Notifications > Alerts: Filter by Severity, Server and Date (the presets), "Filters are ON/OFF", and Export to a CSV
  file (all, or only the alerts shown).
- Pool and Server menus: Change Server Password (the local root; the connection keeps the new password) and Pool >
  Reconnect As.
- vApps: Pool > Manage vApps... (New vApp wizard, Delete, Properties with the VMs and their startup sequence, Start,
  Shut down) and VM > Assign to vApp, with the C# questions. Its Import button is not ported (File > Import...
  imports an appliance).
- The user templates in the tree and the Templates menu: Create VM From Selection, Copy, Delete Template, Properties;
  Server > Remove (the server leaves the list).
- Tools > Options: certificate warnings, the pool secret reminder, the proxy server and connection timeout, the
  dismiss alerts confirmation, area or line graphs, and the last tab of each resource.
- File > Import... and VM > Export...: a VM as an XVA file (ExportVmAction, ImportVmAction), checked as it downloads
  (the block checksums, SHA-1 or xxHash64) and optionally read back; the import picks the pool or a server and the SR,
  and can start the VM.
- VM menu: Take a Snapshot... and Convert to Template... (ConvertVMToTemplateCommand); Server menu: Restart Toolstack
  and Remove From Pool (a member running only its control domain), with the C# questions.
- The Search tab of a VM, server or pool: the overview of its connection ("<pool>: Overview").
- The Search overview follows the C# SearchView: the pool, its servers and their VMs (folding with +/-), with the live
  CPU and memory (with bars), disks, network, address and uptime columns, read every ten seconds (MetricUpdater).
- Networking tab: Add Interface, Properties, Remove Interface and Activate/Deactivate (network, generated or typed MAC
  with the duplicate MAC question, QoS limit), plugged live when the server allows it.

### Fixed
- Errors show the server's text without the "RuntimeError:" prefix.

### Changed
- The Organization Views and searches follow the C# visibility rules (VM.Show, SR.Show, Network.Show).

## [0.1.0] — 2026-11-05

**The first release of the Go app: the desktop window for Linux, Windows and macOS (amd64), the CLI and the MCP
server, following the Windows app's screens, checked on XCP-ng 8.3.**

### Added
- macOS: `build/darwin/build-on-mac.sh` builds the app natively on a Mac (arm64 or amd64), signs it with the hardened
  runtime and notarizes it when the Apple credentials are given; `make icons` copies the C# icons without a container.
- The toolbar's Pause button pauses a running VM and becomes Unpause on a paused one (PauseVMCommand / UnPauseVMCommand,
  no question); a paused VM has the C# icon. The MCP server keeps the seven C# operations.
- Organization Views (Tags, Folders with the empty ones, Custom Fields, vApps) and Saved Searches (the seven C#
  default searches) from their navigation buttons and the View menu, over every connected server; a group shows the
  Search tab of its objects, and the views follow the servers live. Searches saved in pools are listed, not run yet.
- New Pool (NewPoolDialog / CreatePoolAction): the coordinator and the members among the connected servers, each with
  the reason it cannot join (PoolJoinRules: a pool already, CPU vendor, server version, shared storage, running VMs,
  HA, clustering, TLS verification), the CPU levelling question, then the pool's name and pool.join from each member,
  which leaves the server list. The window keeps each connection's password in memory for that join, as the C# app.
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
