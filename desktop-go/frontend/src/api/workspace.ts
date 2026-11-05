// The window's API: internal/gui.Workspace, bound by Wails. A web build will implement the same interface over HTTP.
import { Call } from "@wailsio/runtime";

const svc = "xcpng-center/internal/gui.Workspace.";
const call = <T>(method: string, ...args: unknown[]) => Call.ByName(svc + method, ...args) as Promise<T>;

export interface ServerRow { key: string; host: string; port: number; user: string; name: string; connected: boolean }
export interface ConnectResult { connected: boolean; key?: string; fingerprint?: string; previous?: string; error?: string }
export interface XenObject { kind: string; uuid: string; name: string; description?: string; power_state?: string; host?: string; host_uuid?: string }

export const api = {
  servers: () => call<ServerRow[]>("Servers"),
  connect: (address: string, user: string, password: string) => call<ConnectResult>("Connect", address, user, password),
  trust: (key: string, fingerprint: string) => call<void>("Trust", key, fingerprint),
  disconnect: (key: string) => call<void>("Disconnect", key),
  objects: (key: string, kind: string) => call<XenObject[] | null>("Objects", key, kind),
};
