// The window's API: internal/gui.Workspace, bound by Wails. A web build will implement the same interface over HTTP.
import { Call } from "@wailsio/runtime";

const svc = "xcpng-center/internal/gui.Workspace.";
// a failed call rejects with the Go error's text alone (String(e) shows it without Wails' "RuntimeError: ")
const call = <T>(method: string, ...args: unknown[]) =>
  (Call.ByName(svc + method, ...args) as Promise<T>).catch((e: unknown) => { throw e instanceof Error ? e.message : e; });

export interface ServerRow { key: string; host: string; port: number; user: string; name: string; connected: boolean }
export interface ConnectResult { connected: boolean; key?: string; fingerprint?: string; previous?: string; error?: string }
export interface XenObject { kind: string; uuid: string; name: string; description?: string; power_state?: string; host?: string; host_uuid?: string; allowed?: string[]; ha_protected?: boolean;
  maintenance?: boolean; coordinator?: boolean; live?: boolean; control_domain?: string; icon?: string; resident?: number }
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
export interface MemoryEdit { name: string; can_edit: boolean; advanced: boolean; dynamic: boolean; unavailable: string; rubric: string;
  in_mb: boolean; static_min: number; static_max: number; dynamic_min: number; dynamic_max: number; max: number }
export interface Performance { times: number[]; start: number; end: number;
  graphs: { title: string; max: number; divisor: number; label: string; series: { name: string; values: (number | null)[] }[] }[];
  events: { time: string; name: string; text: string }[] }
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
export interface PoolCandidate { key: string; name: string; reason: string; can_coordinate: boolean; can_join: boolean; is_coordinator: boolean }
export interface ViewInfo { id: string; name: string; icon: string; supported: boolean }
export interface ViewNode { name: string; key?: string; icon?: string; server?: string; object?: XenObject; children: ViewNode[] }
export interface Alert { uuid: string; priority: number; title: string; description: string; applies_to: string; time: string }

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
  pickExportFile: (name: string) => call<string>("PickExportFile", name),
  pickImportFile: () => call<string>("PickImportFile"),
  exportBlocker: (key: string, vm: string) => call<string>("ExportBlocker", key, vm),
  exportVm: (key: string, vm: string, file: string, verify: boolean) => call<void>("ExportVM", key, vm, file, verify),
  importOptions: (key: string, host: string) => call<{ servers: { uuid: string; name: string; live: boolean }[]; srs: { uuid: string; text: string }[] | null; default_sr: string }>("ImportOptions", key, host),
  importFileError: (file: string) => call<string>("ImportFileError", file),
  importVm: (key: string, file: string, host: string, sr: string, start: boolean) => call<string>("ImportVM", key, file, host, sr, start),
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
  srQuestion: (command: string, name: string) => call<HostQuestion>("SRQuestion", command, name),
  srCommand: (key: string, uuid: string, command: string) => call<void>("SRCommand", key, uuid, command),
  srAttachWarning: (kind: string, uuid: string) => call<string>("SRAttachWarning", kind, uuid),
  openConsole: (key: string, vm: string) => call<string>("OpenConsole", key, vm),
  vmConfirmation: (op: string, haProtected: boolean) => call<Confirm>("VMConfirmation", op, haProtected),
  vmPower: (key: string, uuid: string, op: string) => call<void>("VMPower", key, uuid, op),
};
