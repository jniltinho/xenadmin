// The window's API: internal/gui.Workspace, bound by Wails. A web build will implement the same interface over HTTP.
import { Call } from "@wailsio/runtime";

const svc = "xcpng-center/internal/gui.Workspace.";
// a failed call rejects with the Go error's text alone (String(e) shows it without Wails' "RuntimeError: ")
const call = <T>(method: string, ...args: unknown[]) =>
  (Call.ByName(svc + method, ...args) as Promise<T>).catch((e: unknown) => { throw e instanceof Error ? e.message : e; });

export interface ServerRow { key: string; host: string; port: number; user: string; name: string; connected: boolean }
export interface ConnectResult { connected: boolean; key?: string; fingerprint?: string; previous?: string; error?: string }
export interface XenObject { kind: string; uuid: string; name: string; description?: string; power_state?: string; host?: string; host_uuid?: string; allowed?: string[]; ha_protected?: boolean;
  maintenance?: boolean; coordinator?: boolean; live?: boolean; control_domain?: string; icon?: string; resident?: number; user_template?: boolean;
  gpu_tab?: boolean; usb_tab?: boolean }
export interface Confirm { title: string; text: string }
export interface HostQuestion { error?: string; title?: string; text?: string; option?: string }
export interface Table { title: string; headers: string[]; rows: string[][]; keys: string[] }
export interface Disk { vbd: string; attached: boolean; can_activate: boolean; can_deactivate: boolean; vdi: string; name: string;
  system: boolean; can_detach: boolean; can_delete: boolean }
export interface NewDisk { name: string; srs: { uuid: string; text: string }[]; default_sr: string }
export interface VMStorage { table: Table; disks: Disk[];
  dvd?: string; dvd_name?: string; can_add_dvd: boolean; iso: string; isos: { vdi: string; text: string }[] }
export interface VMProperties { name: string; description: string; halted: boolean; running: boolean; hotplug: boolean; vcpus_max: number;
  vcpus_at_startup: number; cores_per_socket: number; priority: number; min_vcpus: number; max_physical: number;
  vcpu_choices: { vcpus: number; topologies: { sockets: number; cores: number; text: string }[] }[];
  home_server: string; servers: { uuid: string; name: string; live: boolean }[]; can_have_none: boolean;
  hvm: boolean; auto_power_on: boolean; boot_order: string; boot_devices: Record<string, string>; pv_args: string; has_cd: boolean;
  pv_boot_from_cd: boolean; shadow_multiplier: number; frozen: boolean; order: number; start_delay: number; ha_enabled: boolean }
export interface VIF { uuid: string; device: string; network: string; mac: string; qos: boolean; limit: string; attached: boolean; can_edit: boolean;
  can_toggle: boolean; toggle_label: string; edit_tip: string; remove_tip: string; toggle_tip: string }
export interface VMNetworking { table: Table; vifs: VIF[]; can_add: boolean; at_max: string }
export interface VIFSettings { network: string; mac: string; qos: boolean; limit: string }
export interface OverviewRow { cpu: string; cpu_bar: number; memory: string; memory_bar: number; disks: string; network: string; address: string; uptime: string }
export interface AppOptions { no_new_certificate_warning?: boolean; no_rotate_secret_reminder?: boolean; proxy?: string; proxy_address?: string; proxy_port?: number;
  bypass_proxy_for_servers?: boolean; connection_timeout?: number; no_dismiss_confirmation?: boolean; area_graphs?: boolean; no_remember_last_tab?: boolean }
export interface ObjectProps { name: string; description: string; folder: string; tags: string[]; all_tags: string[]; folders: string[]; iqn?: string;
  log_destination?: string; multipath?: boolean; multipath_editable?: boolean;
  fields: { name: string; type: string; value: string }[] }
export interface VappVM { uuid: string; name: string; power_state: string; vapp: string; order: number; delay: number }
export interface Vapp { uuid: string; name: string; description: string; vms: VappVM[]; can_start: boolean; can_shutdown: boolean }
export interface Vapps { pool: string; is_pool: boolean; vapps: Vapp[]; candidates: VappVM[] }
export interface Vtpms { name: string; blocker: string; add_blocker: string; vtpms: { uuid: string; protected: boolean; unique: boolean; remove_blocker: string }[] }
export interface GPUTab { policy: string; policy_value: string; no_gpu: string; density: string; performance: string;
  groups: { hosts: { label: string; gpus: string[] }[]; types: { uuid: string; name: string; max: string; vram: string; enabled: boolean; in_use: boolean }[];
    pgpus: string[]; can_edit: boolean; edit: string; title: string; rubric: string }[] }
export interface USBTab { title: string; table: Table;
  devices: { uuid: string; passthrough: boolean; button: string; can_toggle: boolean; question: string; yes: string }[] }
export interface VGPURowFacts { known: boolean; passthrough: boolean; multiple: boolean; experimental: boolean }
export interface VMGPURow { uuid: string; group: string; type: string; slot: string; name: string; per_gpu: string; vram: string; facts: VGPURowFacts }
export interface GPUPageInput { halted: boolean; ha_enabled: boolean; ha_priority: string; gpus_available: boolean; rows: VGPURowFacts[];
  original: number; changed: boolean; selected: boolean }
export interface GPUPageState { can_add: boolean; can_delete: boolean; experimental: boolean; rdp: boolean; need_gpu: boolean; stop_vm: boolean; ha: boolean; multiple: string }
export interface VMGPUPage { shown: boolean; rows: VMGPURow[]; input: GPUPageInput }
export interface VGPUPicks { picked: number;
  items: { group: string; type: string; text: string; header: boolean; sub: boolean; enabled: boolean; row: VMGPURow }[] }
export interface VMUSBPage { shown: boolean; title: string; table: Table; can_attach: boolean; warning: string; confirm: string; confirm_title: string }
export interface USBAttachChoices { hosts: { name: string; devices: { uuid: string; name: string }[] }[]; empty: string; upgrade: boolean }
export interface PerfmonPage { interval: number;
  alerts: { name: string; enabled: boolean; level: number; period: number; min: number; max: number; sub: string }[] }
export interface MemoryBars { used: number; dynamic_min: number; dynamic_max: number; used_text: string;
  hosts: { name: string; segments: { name: string; bytes: number; text: string }[] }[] }
export interface MemoryEdit { name: string; can_edit: boolean; advanced: boolean; dynamic: boolean; unavailable: string; rubric: string;
  in_mb: boolean; static_min: number; static_max: number; dynamic_min: number; dynamic_max: number; max: number }
export interface Performance { times: number[]; start: number; end: number;
  graphs: { title: string; max: number; divisor: number; label: string; sources: string[]; series: { name: string; colour: string; values: (number | null)[] }[] }[];
  events: { time: string; name: string; text: string }[]; default: boolean; new_blocker?: string; edit_blocker?: string; new_name: string }
export interface DesignedGraph { name: string; sources: string[] }
export interface DataSource { name: string; text: string; description: string; type: string; enabled: boolean; standard: boolean; hidden: boolean; colour?: string }
export interface NewVMOptions { templates: { uuid: string; name: string; custom: boolean }[]; isos: { vdi: string; text: string }[];
  servers: { uuid: string; name: string; live: boolean }[]; networks: { vdi: string; text: string }[]; srs: { uuid: string; text: string }[]; default_sr: string }
export interface NewVMDefaults { name: string; description: string; hvm: boolean; boot_mode: string; install: string; vcpus: number; cores: number;
  memory: number; disks: { position: string; name: string; description: string; size: number; sr: string; bootable: boolean; inherited: boolean }[];
  networks: { network: string; device?: string; qos_type?: string; qos_params?: Record<string, string> }[]; no_dvd: boolean; eliloader: boolean }
export interface Snapshot { uuid: string; name: string; description?: string; time: string; type: string; size: number }
export interface VMSnapshots { snapshots: Snapshot[]; modes: { disk: boolean; memory: boolean; memory_info: string }; snapshot_allowed: boolean }
export interface NewSRKind { kind: string; name: string; blurb: string; default_name: string; iso: boolean }
export interface DeleteItem { uuid: string; name: string; size?: number; checked: boolean }
export interface VMDelete { title: string; text: string; disks: DeleteItem[]; snapshots: DeleteItem[] }
export interface VMCopy { name: string; description: string; clone: boolean; copy: boolean; srs: { uuid: string; text: string; reason?: string }[]; default_sr: string }
export interface IQN { index: number; iqn: string; address: string; port: number; text: string }
export interface LUN { scsi_id: string; lun_id: number; vendor: string; serial: string; size: number; text: string }
export interface Appliance { name: string; warnings: string[] | null; disks: number; interfaces: number;
  vms: { sys_id: string; name: string; disks: { name: string; capacity: number }[] | null; macs: string[] | null }[] | null;
  eulas: string[] | null; has_manifest: boolean; required_space: number; has_signature: boolean;
  certificate?: { subject: string; issuer: string; valid_from: string; valid_to: string; thumbprint: string }; certificate_error?: string }
// a VM's server ("" for the pool), the SR of each of its disks and the network of each of its interfaces
export interface ApplianceMapping { sys_id: string; host: string; srs: string[]; networks: string[] }
export interface ExportChoices { pool: string; is_pool: boolean;
  vms: { uuid: string; name: string; description: string; size: string; vapp: string; selected: boolean; vtpm: boolean }[] }
export interface ApplianceExport { dir: string; name: string; vms: string[]; eulas: string[]; manifest: boolean; ova: boolean; compress: boolean; verify: boolean;
  sign: boolean; certificate: string; password: string }
export interface PoolCandidate { key: string; name: string; reason: string; can_coordinate: boolean; can_join: boolean; is_coordinator: boolean }
export interface ViewInfo { id: string; name: string; icon: string; supported: boolean }
export interface ViewNode { name: string; key?: string; icon?: string; server?: string; object?: XenObject; children: ViewNode[] }
// the Search tab's editor (internal/gui/search.go)
export interface SearchEditorForm { name: string; types: string[] | null; grouping: number; has_filter: boolean; type_list: { id: string; name: string }[];
  groupings: { id: string; name: string }[]; pools: { key: string; name: string }[] }
export interface SearchEdit { name: string; types: string[]; text: string; grouping: number; grouping_changed: boolean; edited: string; saved: boolean }
export interface Alert { uuid: string; priority: number; title: string; description: string; applies_to: string; time: string; host: string }
// server networks and bonds (internal/gui/network.go)
export interface HostNetworking { table: Table; networks: { uuid: string; can_edit: boolean; can_remove: boolean; remove_tip?: string; bond: boolean }[]; can_add: boolean }
export interface HostNics { table: Table; nics: { uuid: string; bond_network?: string }[]; can_bond: boolean }
export interface BondNIC { uuid: string; name: string; in_bond: boolean; mac: string; link: string; speed: string; vendor: string; mtu: number }
export interface BondChoices { nics: BondNIC[]; limit: number; lacp: boolean }
export interface BondSpec { members: string[]; mode: string; hashing: string; mtu: number; auto: boolean }
export interface WizardNIC { uuid: string; name: string; max_mtu: number; vlans: number[] }
export interface NewNetworkChoices { title: string; host: string; external_name: string; private_name: string; nics: WizardNIC[]; min_vlan: number; bond: BondChoices;
  sriov: boolean; sriov_why: string; sriov_name: string; sriov_nics: WizardNIC[]; sriov_warning: string; sriov_ok: string }
export interface NewSriovSpec { name: string; description: string; auto: boolean; host: string; nic: string }
export interface NewNetworkSpec { external: boolean; name: string; description: string; auto: boolean; host: string; nic: string; vlan: string; mtu: string }
export interface BondWarning { block: boolean; title: string; text: string; ok: string }
export interface MTUState { shown: boolean; enabled: boolean; min: number; max: number; note: string }
export interface NICState { shown: boolean; enabled: boolean; note: string }
export interface SettingsNIC { uuid: string; name: string; vlans: number[] }
export interface NetworkEdit {
  name: string; description: string; auto: boolean; auto_enabled: boolean; mtu: number; mtu_state: MTUState; vlan_mtu_state: MTUState
  nic_state: NICState; nics: SettingsNIC[]; nic: string; vlan: number; min_vlan: number; bond: boolean; lacp: boolean; mode: string; hashing: string
}
export interface NetworkChanges { name: string; description: string; auto: boolean; mtu: string; nic: string; vlan: string; mode: string; hashing: string }

export const api = {
  servers: () => call<ServerRow[]>("Servers"),
  connect: (address: string, user: string, password: string) => call<ConnectResult>("Connect", address, user, password),
  trust: (key: string, fingerprint: string) => call<void>("Trust", key, fingerprint),
  disconnect: (key: string) => call<void>("Disconnect", key),
  objects: (key: string, kind: string) => call<XenObject[] | null>("Objects", key, kind),
  alerts: (key: string) => call<Alert[] | null>("Alerts", key),
  dismissAlerts: (key: string, uuids: string[]) => call<void>("DismissAlerts", key, uuids),
  dismissConfirmation: (count: number, all: boolean) => call<Confirm>("DismissConfirmation", count, all),
  hostQuestion: (key: string, uuid: string, op: string) => call<HostQuestion>("HostQuestion", key, uuid, op),
  hostCommand: (key: string, uuid: string, op: string, restore: boolean) => call<void>("HostCommand", key, uuid, op, restore),
  properties: (key: string, vm: string) => call<VMProperties>("Properties", key, vm),
  saveProperties: (key: string, vm: string, changes: Record<string, unknown>) => call<void>("SaveProperties", key, vm, changes),
  customFields: (key: string, vm: string) => call<{ name: string; type: string; value: string }[]>("CustomFields", key, vm),
  saveCustomFields: (key: string, vm: string, values: Record<string, string>) => call<void>("SaveCustomFields", key, vm, values),
  addCustomField: (key: string, name: string, kind: string) => call<void>("AddCustomField", key, name, kind),
  removeCustomField: (key: string, name: string) => call<void>("RemoveCustomField", key, name),
  forget: (key: string) => call<void>("Forget", key),
  storage: (key: string, vm: string) => call<VMStorage>("Storage", key, vm),
  networking: (key: string, vm: string) => call<VMNetworking>("Networking", key, vm),
  vifNetworks: (key: string, vm: string) => call<{ uuid: string; name: string }[]>("VIFNetworks", key, vm),
  vifQuestion: (vif: string, mac: string) => call<Confirm>("VIFQuestion", vif, mac),
  vifDeleteQuestion: () => call<Confirm>("VIFDeleteQuestion"),
  saveVif: (key: string, vm: string, vif: string, set: VIFSettings) => call<string>("SaveVIF", key, vm, vif, set),
  deleteVif: (key: string, vif: string) => call<void>("DeleteVIF", key, vif),
  setVifActive: (key: string, vif: string, active: boolean) => call<void>("SetVIFActive", key, vif, active),
  memory: (key: string, vm: string) => call<Table>("Memory", key, vm),
  hostTab: (key: string, kind: string, uuid: string, tab: string) => call<Table>("HostTab", key, kind, uuid, tab),
  overview: (key: string) => call<Record<string, OverviewRow> | null>("Overview", key),
  convertQuestion: (name: string) => call<Confirm>("ConvertQuestion", name),
  convertToTemplate: (key: string, vm: string) => call<void>("ConvertToTemplate", key, vm),
  joinCandidates: (pool: string) => call<{ key: string; name: string; reason: string; can_join: boolean }[]>("JoinCandidates", pool),
  addToPoolQuestions: (pool: string, member: string) => call<string[]>("AddToPoolQuestions", pool, member),
  addToPool: (pool: string, member: string) => call<void>("AddToPool", pool, member),
  rotateSecretBlocker: (key: string) => call<string>("RotateSecretBlocker", key),
  rotateSecret: (key: string) => call<string>("RotateSecret", key),
  makeStandaloneBlocker: (key: string) => call<string>("MakeStandaloneBlocker", key),
  makeStandalone: (key: string) => call<void>("MakeStandalone", key),
  options: () => call<AppOptions>("Options"),
  setOptions: (o: AppOptions) => call<void>("SetOptions", o),
  objectProperties: (key: string, kind: string, uuid: string) => call<ObjectProps>("ObjectProperties", key, kind, uuid),
  saveObjectProperties: (key: string, kind: string, uuid: string, ch: Record<string, unknown>) => call<void>("SaveObjectProperties", key, kind, uuid, ch),
  vapps: (key: string) => call<Vapps>("Vapps", key),
  vappDeleteQuestion: (name: string, vms: number) => call<string>("VappDeleteQuestion", name, vms),
  vappChangeWarning: (moves: { vm: string; from: string }[], target: string) => call<string>("VappChangeWarning", moves, target),
  vappSummary: (name: string, vms: string[]) => call<string>("VappSummary", name, vms),
  createVapp: (key: string, spec: Record<string, unknown>) => call<string>("CreateVapp", key, spec),
  saveVapp: (key: string, uuid: string, spec: Record<string, unknown>) => call<void>("SaveVapp", key, uuid, spec),
  deleteVapp: (key: string, uuid: string) => call<void>("DeleteVapp", key, uuid),
  vappPower: (key: string, uuid: string, op: string) => call<void>("VappPower", key, uuid, op),
  assignVapp: (key: string, vapp: string, vm: string) => call<void>("AssignVapp", key, vapp, vm),
  canChangePassword: (key: string) => call<boolean>("CanChangePassword", key),
  changePasswordTitle: (key: string) => call<string>("ChangePasswordTitle", key),
  changePassword: (key: string, current: string, next: string, confirm: string) => call<void>("ChangePassword", key, current, next, confirm),
  vtpms: (key: string, vm: string) => call<Vtpms>("Vtpms", key, vm),
  addVtpm: (key: string, vm: string) => call<void>("AddVtpm", key, vm),
  vtpmRemoveWarning: () => call<string>("VtpmRemoveWarning"),
  removeVtpm: (key: string, vm: string, vtpm: string) => call<void>("RemoveVtpm", key, vm, vtpm),
  exportAlerts: (rows: { title: string; priority: number; description: string; applies_to: string; time: string }[]) => call<string>("ExportAlerts", rows),
  memoryBars: (key: string, kind: string, uuid: string) => call<MemoryBars>("MemoryBars", key, kind, uuid),
  // the Alerts page of the Properties (internal/gui/hardware.go)
  perfmon: (key: string, kind: string, uuid: string) => call<PerfmonPage>("Perfmon", key, kind, uuid),
  savePerfmon: (key: string, kind: string, uuid: string, on: { name: string; level: number; period: number }[], interval: number) =>
    call<void>("SavePerfmon", key, kind, uuid, on, interval),
  // the GPU and USB tabs (internal/gui/hardware.go)
  gpu: (key: string, kind: string, uuid: string) => call<GPUTab>("GPU", key, kind, uuid),
  setGpuPlacement: (key: string, algorithm: string) => call<void>("SetGPUPlacement", key, algorithm),
  setVgpuTypes: (key: string, pgpus: string[], changes: Record<string, boolean>) => call<void>("SetVGPUTypes", key, pgpus, changes),
  usb: (key: string, uuid: string) => call<USBTab>("USB", key, uuid),
  setUsbPassthrough: (key: string, pusb: string, wasEnabled: boolean) => call<void>("SetUSBPassthrough", key, pusb, wasEnabled),
  // a VM's GPU and USB pages (internal/gui/vmhardware.go)
  vmGpu: (key: string, vm: string) => call<VMGPUPage>("VMGPU", key, vm),
  vmGpuState: (input: GPUPageInput) => call<GPUPageState>("VMGPUState", input),
  vgpuChoices: (key: string, vm: string, existing: string[]) => call<VGPUPicks>("VGPUChoices", key, vm, existing),
  assignVgpus: (key: string, vm: string, removed: string[], add: { group: string; type: string }[]) => call<void>("AssignVGPUs", key, vm, removed, add),
  vmUsb: (key: string, vm: string, haPriority: string) => call<VMUSBPage>("VMUSB", key, vm, haPriority),
  usbAttachChoices: (key: string, vm: string) => call<USBAttachChoices>("USBAttachChoices", key, vm),
  attachUsb: (key: string, vm: string, pusb: string) => call<void>("AttachUSB", key, vm, pusb),
  detachUsb: (key: string, vm: string, vusb: string) => call<void>("DetachUSB", key, vm, vusb),
  pickImportFile: () => call<string>("PickImportFile"),
  exportBlocker: (key: string, vm: string) => call<string>("ExportBlocker", key, vm),
  exportVm: (key: string, vm: string, file: string, verify: boolean) => call<void>("ExportVM", key, vm, file, verify),
  importOptions: (key: string, host: string) => call<{ servers: { uuid: string; name: string; live: boolean }[]; srs: { uuid: string; text: string }[] | null; default_sr: string;
    networks: { uuid: string; name: string }[] | null }>("ImportOptions", key, host),
  importFileError: (file: string) => call<string>("ImportFileError", file),
  importVm: (key: string, file: string, host: string, sr: string, start: boolean) => call<string>("ImportVM", key, file, host, sr, start),
  isAppliance: (file: string) => call<boolean>("IsAppliance", file),
  applianceInfo: (file: string) => call<Appliance>("ApplianceInfo", file),
  applianceMappingsError: (key: string, a: Appliance, maps: ApplianceMapping[]) => call<string>("ApplianceMappingsError", key, a, maps),
  importAppliance: (key: string, file: string, maps: ApplianceMapping[], verify: boolean, start: boolean) =>
    call<string>("ImportAppliance", key, file, maps, verify, start),
  // the Export wizard's OVF/OVA path (internal/gui/ovfexport.go)
  exportChoices: (key: string, vm: string, vapp: string) => call<ExportChoices>("ExportChoices", key, vm, vapp),
  exportFolder: () => call<string>("ExportFolder"),
  pickExportFolder: (dir: string) => call<string>("PickExportFolder", dir),
  pickEulaFile: () => call<string>("PickEulaFile"),
  pickCertificateFile: () => call<string>("PickCertificateFile"),
  signingCertificateError: (file: string, password: string) => call<string>("SigningCertificateError", file, password),
  applianceTargetError: (dir: string, name: string, xva: boolean) => call<string>("ApplianceTargetError", dir, name, xva),
  exportAppliance: (key: string, e: ApplianceExport) => call<string>("ExportAppliance", key, e),
  memoryEdit: (key: string, vm: string) => call<MemoryEdit>("MemoryEdit", key, vm),
  memoryQuestion: (key: string, vm: string, dynMin: number, dynMax: number, staticMax: number) => call<string>("MemoryQuestion", key, vm, dynMin, dynMax, staticMax),
  changeMemory: (key: string, vm: string, dynMin: number, dynMax: number, staticMax: number) => call<void>("ChangeMemory", key, vm, dynMin, dynMax, staticMax),
  setDiskActive: (key: string, vbd: string, active: boolean) => call<void>("SetDiskActive", key, vbd, active),
  newDiskOptions: (key: string, vm: string) => call<NewDisk>("NewDiskOptions", key, vm),
  addDisk: (key: string, vm: string, name: string, description: string, sr: string, size: number) => call<string>("AddDisk", key, vm, name, description, sr, size),
  attachableDisks: (key: string, vm: string) => call<{ vdi: string; text: string }[] | null>("AttachableDisks", key, vm),
  attachDisk: (key: string, vm: string, vdi: string, readOnly: boolean) => call<string>("AttachDisk", key, vm, vdi, readOnly),
  diskQuestion: (del: boolean, system: boolean) => call<Confirm>("DiskQuestion", del, system),
  detachDisk: (key: string, vbd: string) => call<void>("DetachDisk", key, vbd),
  deleteDisk: (key: string, vm: string, vbd: string) => call<void>("DeleteDisk", key, vm, vbd),
  createDvdDrive: (key: string, vm: string) => call<string>("CreateDVDDrive", key, vm),
  changeIso: (key: string, dvd: string, iso: string) => call<void>("ChangeISO", key, dvd, iso),
  snapshots: (key: string, vm: string) => call<VMSnapshots>("Snapshots", key, vm),
  snapshotQuestion: (op: string, name: string, snapshotAllowed: boolean) => call<HostQuestion>("SnapshotQuestion", op, name, snapshotAllowed),
  takeSnapshot: (key: string, vm: string, name: string, description: string, memory: boolean) => call<void>("TakeSnapshot", key, vm, name, description, memory),
  revertSnapshot: (key: string, snapshot: string) => call<void>("RevertSnapshot", key, snapshot),
  deleteSnapshot: (key: string, snapshot: string) => call<void>("DeleteSnapshot", key, snapshot),
  performance: (key: string, kind: string, uuid: string, zoom: string) => call<Performance>("Performance", key, kind, uuid, zoom),
  graphDataSources: (key: string, kind: string, uuid: string) => call<DataSource[]>("GraphDataSources", key, kind, uuid),
  saveGraphs: (key: string, kind: string, uuid: string, base: DesignedGraph[], graphs: DesignedGraph[], colours: Record<string, string>) =>
    call<void>("SaveGraphs", key, kind, uuid, base, graphs, colours),
  enableDataSource: (key: string, kind: string, uuid: string, name: string, text: string) => call<void>("EnableDataSource", key, kind, uuid, name, text),
  graphTexts: (name: string) => call<{ new_title: string; edit_title: string; delete: string; restore: string }>("GraphTexts", name),
  newVmOptions: (key: string) => call<NewVMOptions>("NewVMOptions", key),
  newVmDefaults: (key: string, template: string, home: string) => call<NewVMDefaults>("NewVMDefaults", key, template, home),
  createVm: (key: string, spec: Record<string, unknown>) => call<{ uuid: string; start_error?: string }>("CreateVM", key, spec),
  newSrKinds: (key: string) => call<NewSRKind[]>("NewSRKinds", key),
  checkSrLocation: (key: string, kind: string, path: string) => call<string>("CheckSRLocation", key, kind, path),
  createSr: (key: string, spec: Record<string, unknown>) => call<string>("CreateSR", key, spec),
  probeSr: (key: string, spec: Record<string, unknown>) => call<{ uuid: string; in_use?: string }[] | null>("ProbeSR", key, spec),
  vmDeleteChoices: (key: string, uuid: string) => call<VMDelete>("VMDeleteChoices", key, uuid),
  deleteVm: (key: string, uuid: string, vbds: string[], snapshots: string[]) => call<void>("DeleteVM", key, uuid, vbds, snapshots),
  vmCopyChoices: (key: string, uuid: string) => call<VMCopy>("VMCopyChoices", key, uuid),
  copyVm: (key: string, uuid: string, name: string, description: string, full: boolean, sr: string) => call<string>("CopyVM", key, uuid, name, description, full, sr),
  vmTargets: (key: string, uuid: string) => call<{ uuid: string; text: string; enabled: boolean; home?: boolean }[]>("VMTargets", key, uuid),
  vmOnHost: (key: string, uuid: string, host: string, op: string) => call<void>("VMOnHost", key, uuid, host, op),
  scanIqns: (key: string, target: string, port: number, user: string, password: string) => call<IQN[]>("ScanIQNs", key, target, port, user, password),
  scanLuns: (key: string, iqn: IQN, user: string, password: string) => call<LUN[]>("ScanLUNs", key, iqn, user, password),
  lunInUse: (scsiId: string) => call<string>("LUNInUse", scsiId),
  lunQuestion: (found: boolean) => call<string>("LUNQuestion", found),
  poolCandidates: (coordinator: string) => call<PoolCandidate[]>("PoolCandidates", coordinator),
  poolWarning: (coordinator: string, members: string[]) => call<string>("PoolWarning", coordinator, members),
  createPool: (coordinator: string, members: string[], name: string, description: string) =>
    call<{ created: boolean; error?: string }>("CreatePool", coordinator, members, name, description),
  views: () => call<[ViewInfo[], ViewInfo[]]>("Views"),
  view: (id: string) => call<ViewNode[] | null>("View", id),
  editSearch: (xml: string) => call<SearchEditorForm>("EditSearch", xml),
  buildSearch: (e: SearchEdit) => call<string>("BuildSearch", e),
  overviewSearch: () => call<string>("OverviewSearch"),
  searchDefinition: (uuid: string) => call<{ xml: string; server: string }>("SearchDefinition", uuid),
  saveSearch: (key: string, xml: string) => call<string>("SaveSearch", key, xml),
  deleteSearch: (uuid: string) => call<void>("DeleteSearch", uuid),
  exportSearch: (xml: string) => call<string>("ExportSearch", xml),
  importSearch: () => call<{ xml: string; name: string }>("ImportSearch"),
  runSearch: (xml: string) => call<ViewNode[] | null>("RunSearch", xml),
  srQuestion: (command: string, name: string) => call<HostQuestion>("SRQuestion", command, name),
  srCommand: (key: string, uuid: string, command: string) => call<void>("SRCommand", key, uuid, command),
  srAttachWarning: (kind: string, uuid: string) => call<string>("SRAttachWarning", kind, uuid),
  openConsole: (key: string, vm: string) => call<string>("OpenConsole", key, vm),
  vmConfirmation: (op: string, haProtected: boolean) => call<Confirm>("VMConfirmation", op, haProtected),
  vmPower: (key: string, uuid: string, op: string) => call<void>("VMPower", key, uuid, op),
  // server networks and bonds (internal/gui/network.go)
  hostNetworking: (key: string, kind: string, uuid: string) => call<HostNetworking>("HostNetworking", key, kind, uuid),
  hostNics: (key: string, uuid: string) => call<HostNics>("HostNics", key, uuid),
  newNetworkOptions: (key: string, kind: string, uuid: string) => call<NewNetworkChoices>("NewNetworkOptions", key, kind, uuid),
  createNetwork: (key: string, spec: NewNetworkSpec) => call<void>("CreateNetwork", key, spec),
  createSriov: (key: string, spec: NewSriovSpec) => call<void>("CreateSriov", key, spec),
  bondQuestion: (key: string, members: string[]) => call<BondWarning | null>("BondQuestion", key, members),
  createBond: (key: string, spec: BondSpec) => call<void>("CreateBond", key, spec),
  deleteBondQuestion: (key: string, network: string) => call<BondWarning>("DeleteBondQuestion", key, network),
  deleteBond: (key: string, network: string) => call<void>("DeleteBond", key, network),
  rescanNics: (key: string, host: string) => call<void>("RescanNics", key, host),
  networkEdit: (key: string, uuid: string) => call<NetworkEdit>("NetworkEdit", key, uuid),
  saveNetwork: (key: string, uuid: string, changes: NetworkChanges) => call<void>("SaveNetwork", key, uuid, changes),
  removeNetworkQuestion: () => call<Confirm>("RemoveNetworkQuestion"),
  removeNetwork: (key: string, uuid: string) => call<void>("RemoveNetwork", key, uuid),
  // Pool > High Availability (internal/gui/ha.go)
  haConfig: (key: string) => call<HAConfig>("HAConfig", key),
  haMaxFailures: (key: string, priorities: Record<string, string>) => call<{ max: number; info: string }>("HAMaxFailures", key, priorities),
  heartbeatSrs: (key: string) => call<HeartbeatSR[]>("HeartbeatSRs", key),
  enableHa: (key: string, sr: string, plan: HAStartup[], ntol: number) => call<void>("EnableHA", key, sr, plan, ntol),
  setHaPriorities: (key: string, plan: HAStartup[], ntol: number) => call<void>("SetHAPriorities", key, plan, ntol),
  haDisableQuestion: (key: string) => call<HostQuestion>("HADisableQuestion", key),
  disableHa: (key: string) => call<void>("DisableHA", key),
  // the General tab (internal/gui/general.go): Details.General's sections
  general: (key: string, kind: string, uuid: string) => call<{ title: string; rows: [string, string][]; expanded: boolean }[]>("General", key, kind, uuid),
  // the Storage tabs' buttons of an SR and of a server or pool (internal/gui/srdisks.go)
  srStorage: (key: string, sr: string) => call<SRStorage>("SRStorage", key, sr),
  physicalStorage: (key: string, kind: string, uuid: string) => call<PhysicalStorage>("PhysicalStorage", key, kind, uuid),
  newSrDiskOptions: (key: string, sr: string) => call<NewDisk>("NewSRDiskOptions", key, sr),
  createSrDisk: (key: string, sr: string, name: string, description: string, size: number) => call<string>("CreateSRDisk", key, sr, name, description, size),
  rescanSrs: (key: string, srs: string[]) => call<void>("RescanSRs", key, srs),
  moveDiskTargets: (key: string, vdi: string) => call<DiskTargets>("MoveDiskTargets", key, vdi),
  moveSrDisk: (key: string, vdi: string, sr: string) => call<string>("MoveSRDisk", key, vdi, sr),
  deleteSrDisk: (key: string, vdi: string) => call<void>("DeleteSRDisk", key, vdi),
  trimQuestion: () => call<Confirm>("TrimQuestion"),
  trimSr: (key: string, sr: string) => call<void>("TrimSR", key, sr),
  // Pool > Disaster Recovery (internal/gui/dr.go)
  drSrs: (key: string) => call<DRSRs>("DRSRs", key),
  drConfigureWarning: (originally: boolean, checked: number) => call<{ text: string; can_save: boolean }>("DRConfigureWarning", originally, checked),
  configureDr: (key: string, srs: string[]) => call<void>("ConfigureDR", key, srs),
  drMetadata: (key: string, srs: string[]) => call<DRPool[]>("DRMetadata", key, srs),
  drPrechecks: (key: string, type: string, choices: DRChoice[]) => call<DRCheck[]>("DRPrechecks", key, type, choices),
  drRecover: (key: string, type: string, choices: DRChoice[], start: string) => call<void>("DRRecover", key, type, choices, start),
  drTestPending: (key: string) => call<number>("DRTestPending", key),
  drTestCleanup: (key: string) => call<void>("DRTestCleanup", key),
};
export interface HAPriority { value: string; text: string }
export interface HAConfig { pool: string; ha_enabled: boolean; blocker?: string; ntol: number; hosts: number; priorities: HAPriority[] | null;
  vms: { uuid: string; name: string; priority: string; order: number; start_delay: number; not_agile?: string }[] | null;
  overcommit: string; ntol_zero_query: string }
export interface HeartbeatSR { uuid: string; name: string; reason?: string }
export interface HAStartup { vm: string; priority: string; order: number; start_delay: number }
// a pool's HA state (app.Object) and the Start Options page's restart priority (app.VMProperties)
export interface XenObject { ha_enabled?: boolean }
export interface VMProperties { ha_restart_priority: string; ha_priorities: HAPriority[] }
// the Storage tabs of an SR and of a server or pool (app.SRStorage, app.PhysicalStorage, app.DiskTargets)
export interface SRStorage { table: Table; can_add: boolean; can_rescan: boolean; rescan_tip?: string;
  disks: { uuid: string; can_edit: boolean; can_delete: boolean; can_move: boolean; delete_title: string; delete_text: string }[] }
export interface PhysicalStorage { table: Table; srs: { uuid: string; can_trim: boolean; trim_tip?: string }[]; can_new: boolean; pool: string }
export interface DiskTargets { migrate: boolean; srs: { uuid: string; text: string; reason?: string }[] }
// Pool > Disaster Recovery (app.DRSRs, app.DRPool, ops.DRChoice, rules.DRCheck)
export interface DRSR { uuid: string; name: string; type: string; free: string; enabled: boolean; has_metadata: boolean }
export interface DRSRs { pool: string; title: string; action: string; max_srs: number; srs: DRSR[] }
export interface DRObject { uuid: string; name: string; vms?: number }
export interface DRPool { uuid: string; name: string; vdi: string; current: boolean; note?: string; vapps: DRObject[]; vms: DRObject[] }
export interface DRChoice { vdi: string; vapps: string[]; vms: string[] }
export interface DRCheck { text: string; error: boolean }
