// The window's API: internal/gui.Workspace, bound by Wails. A web build will implement the same interface over HTTP.
import { Call } from "@wailsio/runtime";

const svc = "xcpng-center/internal/gui.Workspace.";
const call = <T>(method: string, ...args: unknown[]) => Call.ByName(svc + method, ...args) as Promise<T>;

export interface ServerRow { key: string; host: string; port: number; user: string; name: string; connected: boolean }
export interface ConnectResult { connected: boolean; key?: string; fingerprint?: string; previous?: string; error?: string }
export interface XenObject { kind: string; uuid: string; name: string; description?: string; power_state?: string; host?: string; host_uuid?: string; allowed?: string[]; ha_protected?: boolean;
  maintenance?: boolean; coordinator?: boolean; live?: boolean }
export interface Confirm { title: string; text: string }
export interface HostQuestion { error?: string; title?: string; text?: string }
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
  vmConfirmation: (op: string, haProtected: boolean) => call<Confirm>("VMConfirmation", op, haProtected),
  vmPower: (key: string, uuid: string, op: string) => call<void>("VMPower", key, uuid, op),
};
