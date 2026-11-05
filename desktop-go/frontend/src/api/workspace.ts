// The window's API: internal/gui.Workspace, bound by Wails. A web build will implement the same interface over HTTP.
import { Call } from "@wailsio/runtime";

const svc = "xcpng-center/internal/gui.Workspace.";
const call = <T>(method: string, ...args: unknown[]) => Call.ByName(svc + method, ...args) as Promise<T>;

export interface ServerRow { key: string; host: string; port: number; user: string; name: string; connected: boolean }
export interface ConnectResult { connected: boolean; key?: string; fingerprint?: string; previous?: string; error?: string }
export interface XenObject { kind: string; uuid: string; name: string; description?: string; power_state?: string; host?: string; host_uuid?: string; allowed?: string[]; ha_protected?: boolean;
  maintenance?: boolean; coordinator?: boolean; live?: boolean; control_domain?: string }
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
  home_server: string; servers: { uuid: string; name: string; live: boolean }[]; can_have_none: boolean }
export interface Performance { times: number[]; start: number; end: number;
  graphs: { title: string; max: number; divisor: number; label: string; series: { name: string; values: (number | null)[] }[] }[];
  events: { time: string; name: string; text: string }[] }
export interface Snapshot { uuid: string; name: string; description?: string; time: string; type: string; size: number }
export interface VMSnapshots { snapshots: Snapshot[]; modes: { disk: boolean; memory: boolean; memory_info: string }; snapshot_allowed: boolean }
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
  forget: (key: string) => call<void>("Forget", key),
  storage: (key: string, vm: string) => call<VMStorage>("Storage", key, vm),
  networking: (key: string, vm: string) => call<Table>("Networking", key, vm),
  memory: (key: string, vm: string) => call<Table>("Memory", key, vm),
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
  openConsole: (key: string, vm: string) => call<string>("OpenConsole", key, vm),
  vmConfirmation: (op: string, haProtected: boolean) => call<Confirm>("VMConfirmation", op, haProtected),
  vmPower: (key: string, uuid: string, op: string) => call<void>("VMPower", key, uuid, op),
};
