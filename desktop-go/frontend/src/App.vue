<script setup lang="ts">
// The main window (MainWindow): the menus and toolbar of the C# app, the navigation pane (Infrastructure, Objects,
// Organization Views, Saved Searches, Notifications) with the tree and its search box, the tabs of the selected object
// (Home and Search for the root), Add New Server and the certificate question, a status bar.
import { computed, nextTick, onMounted, reactive, ref, watch } from "vue";
import { Events } from "@wailsio/runtime";
import { api, type ViewInfo, type ViewNode, type Alert, type Confirm, type ServerRow, type XenObject } from "./api/workspace";
import ConsoleView from "./ConsoleView.vue";
import { serial } from "./serial";
import SnapshotsView from "./SnapshotsView.vue";
import VmTab from "./VmTab.vue";
import PropertiesDialog from "./PropertiesDialog.vue";
import PerformanceView from "./PerformanceView.vue";
import NewVmWizard from "./NewVmWizard.vue";
import NewSrWizard from "./NewSrWizard.vue";
import MenuList, { type MenuItem } from "./MenuList.vue";
import CopyVmDialog from "./CopyVmDialog.vue";
import NewPoolDialog from "./NewPoolDialog.vue";
import DeleteVmDialog from "./DeleteVmDialog.vue";

interface Node { id: string; label: string; icon: string; object?: XenObject; server?: string; children: Node[]; disconnected?: boolean }

const servers = ref<ServerRow[]>([]);
const data = reactive<Record<string, Record<string, XenObject[]>>>({});
const selected = ref<Node | null>(null);
const expanded = reactive(new Set<string>(["root"]));
const status = ref("Ready");
const failure = ref(""); // the last error, kept visible in the status bar until the next success

const addOpen = ref(false);
const form = reactive({ address: "", user: "root", password: "", error: "", busy: false });
// The VM power commands (StartVMCommand ... ForceVMRebootCommand): label, icon, the allowed_operations entry they need.
const vmCommands = [
  { op: "start", label: "Start", icon: "001_PowerOn", needs: "start", busy: "Starting" },
  { op: "shutdown", label: "Shut Down", icon: "001_ShutDown", needs: "clean_shutdown", busy: "Shutting down" },
  { op: "reboot", label: "Reboot", icon: "001_Reboot", needs: "clean_reboot", busy: "Rebooting" },
  { op: "suspend", label: "Suspend", icon: "000_Paused", menuIcon: "000_paused", needs: "suspend", busy: "Suspending" },
  { op: "resume", label: "Resume", icon: "000_Resumed", needs: "resume", busy: "Resuming" },
  { op: "force_shutdown", label: "Force Shut Down", icon: "001_ForceShutDown", needs: "hard_shutdown", busy: "Forcing shut down of" },
  { op: "force_reboot", label: "Force Reboot", icon: "001_ForceReboot", needs: "hard_reboot", busy: "Forcing reboot of" },
];
type VmCommand = (typeof vmCommands)[number];
// the toolbar's Pause button (PauseVMCommand), Unpause on a paused VM (UnPauseVMCommand); neither asks first
const pauseCommand = computed<VmCommand>(() => (selectedVm.value?.object?.power_state === "Paused"
  ? { op: "unpause", label: "Unpause", icon: "000_Resumed", needs: "unpause", busy: "Unpausing" }
  : { op: "pause", label: "Pause", icon: "000_Paused", needs: "pause", busy: "Pausing" }) as VmCommand);
const selectedVm = computed(() => (selected.value?.object?.kind === "vms" ? selected.value : null));
const canRun = (c: VmCommand, n: Node | null) => !!n?.object?.allowed?.includes(c.needs);
const menu = ref<{ x: number; y: number; items: MenuItem[] } | null>(null);
const openMenu = ref("");
// the message box: a title, a text and its buttons; resolves with the index of the button clicked (-1: closed)
const ask = ref<(Confirm & { buttons: string[]; resolve: (i: number) => void }) | null>(null);
const choose = (q: Confirm, buttons: string[]) => new Promise<number>((resolve) => (ask.value = { ...q, buttons, resolve }));

async function runVm(c: VmCommand, n: Node | null) {
  closeMenus();
  if (!n?.object || !n.server || !canRun(c, n)) return;
  const q = await api.vmConfirmation(c.op, !!n.object.ha_protected);
  if (q.title && !(await confirm(q))) return;
  status.value = `${c.busy} VM '${n.label}'...`;
  let error = "";
  try {
    await api.vmPower(n.server, n.object.uuid, c.op);
  } catch (e) {
    error = `${c.label} '${n.label}': ${String(e)}`;
  }
  await load(n.server);
  if (error) failure.value = error; // the refresh clears old errors, never this one
  reselect(n.id);
}

// keep the selection on the object after a reload (a VM moves under its host when it starts): expand its parents
function reselect(id: string): boolean {
  const path = (n: Node): Node[] | null => {
    if (n.id === id) return [n];
    for (const c of n.children) { const p = path(c); if (p) return [n, ...p]; }
    return null;
  };
  const p = path(currentTree.value);
  if (!p) return false; // gone (or in another view): the selection stays
  p.slice(0, -1).forEach((a) => expanded.add(a.id));
  selected.value = p[p.length - 1];
  return true;
}

// Notifications > Alerts (AlertSummaryPage): every connected server's alerts, newest first; Dismiss asks first.
type View = "infrastructure" | "objects" | "organization" | "saved" | "alerts";
const view = ref<View>("infrastructure");
// NavigationPane's buttons, in the C# order
const navButtons: { v: View; label: string; icon: string }[] = [
  { v: "infrastructure", label: "Infrastructure", icon: "infra_view_24" },
  { v: "objects", label: "Objects", icon: "objects_24" },
  { v: "organization", label: "Organization Views", icon: "org_view_24" },
  { v: "saved", label: "Saved Searches", icon: "saved_searches_24" },
  { v: "alerts", label: "Notifications", icon: "notif_none_24" },
];
type AlertRow = Alert & { server: string };
const alerts = ref<AlertRow[]>([]);
const picked = reactive(new Set<string>());
const confirm = async (q: Confirm) => (await choose(q, ["Yes", "No"])) === 0;

// one alerts read at a time (serial): reloads asked by several servers never finish out of order
const loadAlerts = serial(async () => {
  const rows: AlertRow[] = [];
  for (const s of servers.value.filter((x) => x.connected)) {
    try {
      rows.push(...((await api.alerts(s.key)) ?? []).map((a) => ({ ...a, server: s.key })));
    } catch (e) {
      failure.value = `Could not load the alerts of ${s.key}: ${String(e)}`;
    }
  }
  alerts.value = rows.sort((a, b) => b.time.localeCompare(a.time));
  for (const uuid of [...picked]) if (!rows.some((a) => a.uuid === uuid)) picked.delete(uuid); // the rest stays picked
});

async function dismiss(rows: AlertRow[], all = false) {
  if (!rows.length || !(await confirm(await api.dismissConfirmation(rows.length, all)))) return;
  for (const key of new Set(rows.map((r) => r.server))) {
    try {
      await api.dismissAlerts(key, rows.filter((r) => r.server === key).map((r) => r.uuid));
    } catch (e) {
      failure.value = `Dismiss: ${String(e)}`;
    }
  }
  await loadAlerts();
}

function pick(a: AlertRow, e: MouseEvent) {
  if (!e.ctrlKey && !e.metaKey) picked.clear();
  picked.has(a.uuid) ? picked.delete(a.uuid) : picked.add(a.uuid);
}

const when = (t: string) => new Date(t).toLocaleString(undefined, { dateStyle: "medium", timeStyle: "short" });

async function showView(v: View) {
  view.value = v;
  if (v === "alerts") await loadAlerts();
  // the selection belongs to the tree shown: the same object if it is there, else the root
  if (v !== "alerts" && (!selected.value || !reselect(selected.value.id))) selected.value = currentTree.value;
}

// Organization Views and Saved Searches (OrganizationalView, Search): their button and the View menu list the views;
// the chosen one is read from every connected server and shown as the tree (MainWindowTreeBuilder.CreateNewRootNode)
const views = ref<{ organization: ViewInfo[]; searches: ViewInfo[] }>({ organization: [], searches: [] });
const chosen = ref<ViewInfo | null>(null);
const viewNodes = ref<ViewNode[]>([]);
const viewMenu = ref<{ x: number; y: number; items: MenuItem[] } | null>(null);
async function viewItems(kind: "organization" | "saved"): Promise<MenuItem[]> {
  const [organization, searches] = await api.views();
  views.value = { organization, searches };
  return (kind === "organization" ? organization : searches).map((v) => ({
    label: v.name, icon: v.icon, enabled: v.supported, run: () => { closeMenus(); viewMenu.value = null; openView(kind, v); },
  }));
}
async function navClick(v: View, e: MouseEvent) {
  if (v !== "organization" && v !== "saved") { showView(v); return; }
  const r = (e.currentTarget as HTMLElement).getBoundingClientRect();
  viewMenu.value = { x: r.left + 8, y: r.top, items: [] };
  try { const items = await viewItems(v); if (viewMenu.value) viewMenu.value = { x: r.left + 8, y: r.top - 4 - items.length * 26, items }; }
  catch (e2) { viewMenu.value = null; failure.value = String(e2); }
}
async function openView(kind: "organization" | "saved", v: ViewInfo) {
  chosen.value = v;
  viewNodes.value = [];
  view.value = kind;
  await readView(); // it puts the selection on the view (its root, for an object or group it lacks)
}
// one view read at a time (serial): changes during a read read it once more; a failed read empties the view and says
// why (never the previous view's objects under the new name)
const readView = serial(async () => {
  if (!chosen.value) return;
  const id = chosen.value.id;
  try {
    const nodes = (await api.view(id)) ?? [];
    if (chosen.value?.id === id) {
      viewNodes.value = nodes;
      // the selection follows the reread view: the same node if it is still there, else the view's root
      if ((view.value === "organization" || view.value === "saved") && (!selected.value || !reselect(selected.value.id))) selected.value = viewTree.value;
    }
    // VM.GetVirtualizationStatus: "Unknown" lasts a VM's first two minutes; read again when it may have ended
    if (id === "vms_without_tools" && nodes.some((n) => !n.object && n.name === "Unknown")) setTimeout(() => { if (chosen.value?.id === id) void readView(); }, 125_000);
  } catch (e) {
    if (chosen.value?.id === id) {
      viewNodes.value = [];
      failure.value = `Could not read ${chosen.value.name}: ${String(e)}`;
      if (view.value === "organization" || view.value === "saved") selected.value = viewTree.value; // nothing of the old tree stays shown
    }
  }
});
const viewTree = computed<Node>(() => {
  // ids from identities, not positions: a group is its path of values (or object keys), an object its server, uuid and
  // place (one VM under two tags is two nodes), so selection and expansion survive a reread
  const seg = (s: string) => encodeURIComponent(s);
  const nodeOf = (v: ViewNode, path: string): Node => {
    if (v.object && v.server) {
      const o = v.object;
      const id = `o:${path}|${seg(v.server)}/${o.uuid}`;
      return { id, label: v.name, icon: objectIcon(o), object: o, server: v.server, children: v.children.map((c) => nodeOf(c, `${path}|${o.uuid}`)) };
    }
    const here = `${path}/${seg(v.key || v.name)}`;
    return { id: `g:${here}`, label: v.name, icon: v.icon || "000_Folder_open_h32bit_16", children: v.children.map((c) => nodeOf(c, here)) };
  };
  return { id: "root", label: chosen.value?.name ?? "", icon: chosen.value?.icon ?? "objects_16", children: viewNodes.value.map((v) => nodeOf(v, "")) };
});
const currentTree = computed(() => (view.value === "objects" ? objectsTree.value : view.value === "organization" || view.value === "saved" ? viewTree.value : tree.value));

function answer(i: number) { ask.value?.resolve(i); ask.value = null; }

// The server commands (ShutDownHostCommand, RebootHostCommand, HostMaintenanceModeCommand) with the C# questions.
const hostCommands = [
  { op: "shutdown", label: "Shut Down", icon: "001_ShutDown" },
  { op: "reboot", label: "Reboot", icon: "001_Reboot" },
  { op: "enter_maintenance", label: "Enter Maintenance Mode...", icon: "000_ServerMaintenance" },
  { op: "exit_maintenance", label: "Exit Maintenance Mode", icon: "000_ServerMaintenance" },
];
type HostCommand = (typeof hostCommands)[number];
const selectedHost = computed(() => (selected.value?.object?.kind === "hosts" ? selected.value : null));
const selectedSr = computed(() => (selected.value?.object?.kind === "srs" ? selected.value : null));

// the SR commands of the Storage menu (RepairSRCommand, SetAsDefaultSRCommand, DetachSRCommand, ForgetSRCommand,
// DestroySRCommand); shortcut: Repair runs without the C# RepairSRDialog's per-host list
const srCommands = [{ op: "repair", label: "Repair...", icon: "000_StorageBroken_h32bit_16" }, { op: "set_default", label: "Set as Default", icon: "000_StorageDefault_h32bit_16" }, { op: "detach", label: "Detach..." },
  { op: "forget", label: "Forget" }, { op: "destroy", label: "Destroy..." }];
async function runSr(op: string, label: string, n: Node | null) {
  closeMenus();
  if (!n?.object || !n.server || !n.object.allowed?.includes(op)) return;
  const q = await api.srQuestion(op, n.label);
  if (q.title && (await choose({ title: q.title, text: q.text ?? "" }, [q.option || "Yes", "No"])) !== 0) return;
  label = label.replace("...", "");
  status.value = `${label}: '${n.label}'...`;
  let error = "";
  try { await api.srCommand(n.server, n.object.uuid, op); } catch (e) { error = `${label} '${n.label}': ${String(e)}`; }
  await load(n.server);
  if (error) failure.value = error;
  if (!reselect(n.id)) reselect(n.server); // a forgotten or destroyed SR: its pool
}
const hostCanRun = (c: HostCommand, n: Node | null) =>
  !!n?.object?.live && (c.op === "enter_maintenance" ? !n.object.maintenance : c.op === "exit_maintenance" ? !!n.object.maintenance : true);
function closeMenus() { menu.value = null; openMenu.value = ""; viewMenu.value = null; }

// The tabs of the selected object: General, and Console for a running VM or a live server (its control domain).
type Tab = "home" | "search" | "general" | "memory" | "storage" | "networking" | "console" | "performance" | "snapshots";
const tab = ref<Tab>("general");
const consoleVm = computed(() => {
  const o = selected.value?.object;
  if (o?.kind === "vms" && o.power_state === "Running") return o.uuid;
  if (o?.kind === "hosts" && o.live) return o.control_domain ?? "";
  return "";
});
const tabInfo = {
  home: { label: "Home", icon: "infra_view_16_textured" },
  search: { label: "Search", icon: "000_Search_h32bit_16" },
  general: { label: "General", icon: "000_Info3_h32bit_16" },
  memory: { label: "Memory", icon: "000_GetMemoryInfo_h32bit_16" },
  storage: { label: "Storage", icon: "000_VirtualStorage_h32bit_16" },
  networking: { label: "Networking", icon: "000_Network_h32bit_16" },
  console: { label: "Console", icon: "console_16" },
  performance: { label: "Performance", icon: "001_PerformanceGraph_h32bit_16" },
  snapshots: { label: "Snapshots", icon: "000_VMSnapShotDiskOnly_h32bit_16" },
} as const;
// VMTabs: General, Console, Snapshots for a VM; General, Console for a server
const tabs = computed<Tab[]>(() => {
  const n = selected.value;
  if (!n || n.id === "root") return view.value === "organization" || view.value === "saved" ? ["search"] : ["home", "search"]; // MainWindow: the root shows Home and Search; a view's root, its Search
  if (n.id.startsWith("g:")) return ["search"]; // a group of a view: the Search tab of its objects
  if (n.disconnected) return ["general"]; // GeneralTabPage.GenerateDisconnectedHostBox
  // the C# order: General, Memory, Storage, Networking, Console, Snapshots
  const t: Tab[] = ["general"];
  if (selected.value?.object?.kind === "vms") t.push("memory", "storage", "networking");
  if (consoleVm.value) t.push("console");
  if (selected.value?.object?.kind === "vms" || selected.value?.object?.kind === "hosts") t.push("performance");
  if (selected.value?.object?.kind === "vms") t.push("snapshots");
  return t;
});
function onTabStatus(text: string, error: boolean) {
  if (error) failure.value = text; else { status.value = text || "Ready"; if (!text) failure.value = ""; }
}
watch(tabs, (t) => { if (!t.includes(tab.value)) tab.value = t[0]; }, { immediate: true });

async function runHost(c: HostCommand, n: Node | null) {
  closeMenus();
  if (!n?.object || !n.server || !hostCanRun(c, n)) return;
  const q = await api.hostQuestion(n.server, n.object.uuid, c.op);
  let restore = false;
  if (q.error) { await choose({ title: "XCP-ng Center", text: q.error }, ["OK"]); return; }
  if (c.op === "exit_maintenance") {
    if (q.text) {
      const i = await choose({ title: q.title ?? "", text: q.text }, ["Restore VMs", "Skip", "Cancel"]);
      if (i !== 0 && i !== 1) return; // Cancel leaves the server in maintenance mode
      restore = i === 0;
    }
  } else {
    const enter = c.op === "enter_maintenance";
    if ((await choose({ title: q.title ?? "", text: q.text ?? "" }, enter ? ["Enter Maintenance Mode", "Cancel"] : ["Yes", "No"])) !== 0) return;
  }
  const label = c.label.replace("...", "");
  status.value = `${label}: '${n.label}'...`;
  let error = "";
  try { await api.hostCommand(n.server, n.object.uuid, c.op, restore); } catch (e) { error = `${label} '${n.label}': ${String(e)}`; }
  await load(n.server);
  if (error) failure.value = error;
  reselect(n.id);
}

const cert = ref<{ key: string; fingerprint: string; previous?: string } | null>(null);

const vmIcon = (o: XenObject) => (o.power_state === "Running" ? "000_StartVM_h32bit_16" : o.power_state === "Suspended" ? "000_SuspendVM_h32bit_16"
  : o.power_state === "Paused" ? "000_paused_h32bit_16" : "000_StoppedVM_h32bit_16");
const hostIcon = (h: XenObject) => (!h.live ? "000_ServerDisconnected_h32bit_16" : h.maintenance ? "000_ServerMaintenance_h32bit_16" : "000_TreeConnected_h32bit_16");
// Images.GetIconFor by kind, for the objects of the Organization Views and searches
const objectIcon = (o: XenObject) => ({ vms: vmIcon(o), hosts: hostIcon(o), pools: "000_PoolConnected_h32bit_16", srs: o.icon || "000_Storage_h32bit_16",
  templates: "000_VMTemplate_h32bit_16", networks: "000_Network_h32bit_16", snapshots: "000_VMSnapShotDiskOnly_h32bit_16",
  vdis: "000_VirtualStorage_h32bit_16" } as Record<string, string>)[o.kind] ?? "000_Folder_open_h32bit_16";

const tree = computed<Node>(() => ({
  id: "root", label: "XCP-ng Center", icon: "infra_view_16_textured", children: servers.value.map((s) => {
    const d = data[s.key];
    if (!s.connected || !d) return { id: s.key, label: s.name || s.host, icon: "000_ServerDisconnected_h32bit_16", server: s.key, children: [], disconnected: true };
    const pool = d.pools?.[0];
    const hosts = (d.hosts ?? []).map<Node>((h) => ({
      id: `${s.key}/${h.uuid}`, label: h.name, icon: !h.live ? "000_ServerDisconnected_h32bit_16" : h.maintenance ? "000_ServerMaintenance_h32bit_16" : "000_TreeConnected_h32bit_16", object: h, server: s.key,
      children: (d.vms ?? []).filter((v) => v.host_uuid === h.uuid).map((v) => ({ id: `${s.key}/${v.uuid}`, label: v.name, icon: vmIcon(v), object: v, server: s.key, children: [] })),
    }));
    const loose = (d.vms ?? []).filter((v) => !v.host_uuid).map<Node>((v) => ({ id: `${s.key}/${v.uuid}`, label: v.name, icon: vmIcon(v), object: v, server: s.key, children: [] }));
    const srs = (d.srs ?? []).map<Node>((r) => ({ id: `${s.key}/${r.uuid}`, label: r.name, icon: r.icon || "000_Storage_h32bit_16", object: r, server: s.key, children: [] }));
    // Pool.Name: an unnamed pool is called after its coordinator
    const coordinator = (d.hosts ?? []).find((h) => h.coordinator);
    return { id: s.key, label: pool?.name || coordinator?.name || s.name || s.host, icon: "000_PoolConnected_h32bit_16", object: pool, server: s.key, children: [...hosts, ...loose, ...srs] };
  }),
}));

// NavigationViews.Objects: the objects grouped by type under the root
const objectsTree = computed<Node>(() => {
  const group = (id: string, label: string, icon: string, kind: string, iconOf: (o: XenObject) => string): Node => ({
    id, label, icon, children: Object.entries(data).flatMap(([key, d]) => (d[kind] ?? []).map((o) => ({ id: `${key}/${o.uuid}`, label: o.name, icon: iconOf(o), object: o, server: key, children: [] })))
      .sort((a, b) => a.label.localeCompare(b.label)),
  });
  const groups = [
    group("g:pools", "Pools", "000_PoolConnected_h32bit_16", "pools", () => "000_PoolConnected_h32bit_16"), // pool ids are fixed below
    group("g:hosts", "Servers", "000_TreeConnected_h32bit_16", "hosts", (h) => (!h.live ? "000_ServerDisconnected_h32bit_16" : h.maintenance ? "000_ServerMaintenance_h32bit_16" : "000_TreeConnected_h32bit_16")),
    group("g:vms", "Virtual Machines", "000_StartVM_h32bit_16", "vms", vmIcon),
    group("g:srs", "Storage", "000_Storage_h32bit_16", "srs", (r) => r.icon || "000_Storage_h32bit_16"),
  ];
  for (const n of groups[0].children) n.id = n.server!; // a pool is its server's node, as in the Infrastructure tree
  return { id: "root", label: "Objects by Type", icon: "objects_16", children: groups.filter((g) => g.children.length) };
});

// the tree's search box (MainWindow's "Search..."): the nodes whose name matches, with their parents
const filter = ref("");
const flat = computed(() => {
  const rows: { node: Node; depth: number }[] = [];
  const q = filter.value.trim().toLowerCase();
  const matches = (n: Node): boolean => !q || n.label.toLowerCase().includes(q) || n.children.some(matches);
  const walk = (n: Node, depth: number) => {
    if (depth > 0 && !matches(n)) return;
    rows.push({ node: n, depth });
    if (expanded.has(n.id) || q) n.children.forEach((c) => walk(c, depth + 1));
  };
  walk(currentTree.value, 0);
  return rows;
});

// Back and Forward (MainWindow's history of selections)
const history = ref<string[]>(["root"]);
const historyAt = ref(0);
let navigating = false;
watch(selected, (n) => {
  if (navigating || !n) return;
  if (history.value[historyAt.value] === n.id) return;
  history.value = [...history.value.slice(0, historyAt.value + 1), n.id];
  historyAt.value = history.value.length - 1;
});
function navigate(step: number) {
  const i = historyAt.value + step;
  if (i < 0 || i >= history.value.length) return;
  navigating = true;
  if (reselect(history.value[i])) historyAt.value = i;
  navigating = false;
}

// the toolbar's power buttons (MainWindow.UpdateContent): one of Start / Shut Down / Force Shut Down, Reboot or Force
// Reboot, Suspend or Resume; disabled when the VM does not allow them; a selected server's Shut Down and Reboot are the
// server commands
const cmd = (op: string) => vmCommands.find((c) => c.op === op)!;
const powerButtons = computed(() => {
  const vm = selectedVm.value, host = selectedHost.value;
  const allowed = (op: string) => canRun(cmd(op), vm);
  const startVisible = !!vm && vm.object?.power_state === "Halted";
  const shutdownVisible = !startVisible && (!vm || allowed("shutdown") || !allowed("force_shutdown"));
  const forceRebootVisible = !!vm && !allowed("reboot") && allowed("force_reboot");
  const resumeVisible = !!vm && vm.object?.power_state === "Suspended";
  const list: { op: string; visible: boolean }[] = [
    { op: "start", visible: startVisible }, { op: "shutdown", visible: shutdownVisible }, { op: "force_shutdown", visible: !startVisible && !shutdownVisible },
    { op: "reboot", visible: !forceRebootVisible }, { op: "force_reboot", visible: forceRebootVisible },
    { op: "suspend", visible: !resumeVisible }, { op: "resume", visible: resumeVisible },
  ];
  return list.filter((b) => b.visible).map((b) => {
    const c = cmd(b.op);
    const hostCmd = host && (b.op === "shutdown" || b.op === "reboot") ? hostCommands.find((h) => h.op === b.op)! : null;
    return { ...c, enabled: hostCmd ? hostCanRun(hostCmd, host) : allowed(b.op), run: () => (hostCmd ? runHost(hostCmd, host) : runVm(c, vm)) };
  });
});

// Connect, Disconnect and Forget of a server (the tree's server node)
// AddServerDialog: a known server opens with its address and user, the cursor in Password
const serverField = ref<HTMLInputElement>();
const passwordField = ref<HTMLInputElement>();
function connectTo(n: Node) {
  const s = servers.value.find((x) => x.key === n.server);
  Object.assign(form, { address: s?.key ?? n.label, user: s?.user || "root", password: "", error: "" });
  addOpen.value = true;
  nextTick(() => passwordField.value?.focus());
}
function openAdd() {
  Object.assign(form, { address: "", user: "root", password: "", error: "" });
  addOpen.value = true;
  nextTick(() => serverField.value?.focus());
}
async function disconnect(n: Node) {
  if (!n.server) return;
  await api.disconnect(n.server);
  retire(n.server);
  delete data[n.server];
  if (view.value === "organization" || view.value === "saved") void readView(); // its objects leave the view
  await refreshServers();
  if (!reselect(n.server)) selected.value = currentTree.value; // gone from this view
}
async function forget(n: Node) {
  if (!n.server || !(await confirm({ title: "Remove Server", text: `Remove '${n.label}' from the list of servers?` }))) return;
  try { await api.forget(n.server); retire(n.server); delete data[n.server]; if (view.value === "organization" || view.value === "saved") void readView(); await refreshServers(); selected.value = tree.value; } catch (e) { failure.value = String(e); }
}

// the splitter's position (MainWindow's SplitContainer), kept for the next start
const sideWidth = ref(260);
try { sideWidth.value = Math.min(800, Math.max(160, Number(localStorage.getItem("sideWidth")) || 260)); } catch { /* no storage */ }
function startResize(e: MouseEvent) {
  const x0 = e.clientX, w0 = sideWidth.value;
  const move = (m: MouseEvent) => { sideWidth.value = Math.min(800, Math.max(160, w0 + m.clientX - x0)); };
  const up = () => {
    window.removeEventListener("mousemove", move);
    window.removeEventListener("mouseup", up);
    try { localStorage.setItem("sideWidth", String(sideWidth.value)); } catch { /* no storage */ }
  };
  window.addEventListener("mousemove", move);
  window.addEventListener("mouseup", up);
}

// New VM (NewVMCommand): on the selected object's server, else the first connected one
const newVmOn = ref("");
const targetServer = computed(() => selected.value?.server && data[selected.value.server] ? selected.value.server : servers.value.find((s) => s.connected)?.key ?? "");
function openNewVm() { closeMenus(); if (targetServer.value) newVmOn.value = targetServer.value; }
async function newVmClosed(uuid: string, startError: string) {
  const key = newVmOn.value;
  newVmOn.value = "";
  if (uuid) { await load(key); reselect(`${key}/${uuid}`); }
  if (startError) failure.value = `The VM was created, but it did not start: ${startError}`;
}

// New Storage (NewSRCommand): on the same server as New VM
const newSrOn = ref("");
const poolLabel = (key: string) => tree.value.children.find((n) => n.id === key)?.label ?? key;
function openNewSr() { closeMenus(); if (targetServer.value) newSrOn.value = targetServer.value; }
async function newSrClosed(uuid: string) {
  const key = newSrOn.value;
  newSrOn.value = "";
  if (uuid) { await load(key); reselect(`${key}/${uuid}`); }
}

// New Pool (NewPoolCommand); Add New Server from it connects another server, then the dialog reads the servers again
const newPool = ref(false);
const poolDialog = ref<InstanceType<typeof NewPoolDialog> | null>(null);
async function newPoolClosed(coordinator: string, error: string) {
  newPool.value = false;
  await refreshServers();
  // the members left the server list (they are in the pool now): their objects go too
  for (const key of Object.keys(data)) if (!servers.value.some((s) => s.key === key && s.connected)) { retire(key); delete data[key]; }
  if (coordinator) { await load(coordinator); reselect(coordinator); }
  if (error) failure.value = error;
}

// the Copy VM and Delete VM dialogs
const copyOf = ref<Node | null>(null);
const deleteOf = ref<Node | null>(null);
async function copyClosed(uuid: string) {
  const n = copyOf.value;
  copyOf.value = null;
  if (uuid && n?.server) { await load(n.server); reselect(`${n.server}/${uuid}`); }
}
async function deleteClosed(deleted: boolean) {
  const n = deleteOf.value;
  deleteOf.value = null;
  if (deleted && n?.server) { await load(n.server); reselect(n.server); }
}

// the Properties dialog of a VM
const propsOf = ref<Node | null>(null);
async function propertiesClosed(saved: boolean) {
  const n = propsOf.value;
  propsOf.value = null;
  if (saved && n?.server) { await load(n.server); reselect(n.id); }
}

const soon = "Not available in this version yet";
const hostItems = (n: Node | null): MenuItem[] => hostCommands.map((c) => ({ label: c.label, icon: `${c.icon}_h32bit_16`, enabled: hostCanRun(c, n), run: () => runHost(c, n) }));
const srItems = (n: Node | null): MenuItem[] => srCommands.flatMap((c, i) => [...(i === 2 ? [{ label: "", sep: true }] : []),
  { label: c.label, icon: c.icon, enabled: !!n?.object?.allowed?.includes(c.op), run: () => runSr(c.op, c.label, n) }]);
// CopyVMCommand: a shut down VM that can be cloned or copied; DeleteVMCommand: destroy allowed
const canCopy = (n: Node | null) => n?.object?.power_state === "Halted" && !!n.object.allowed?.some((a) => a === "clone" || a === "copy");
// Resume on, Migrate to and Start on Server (VMOperationHostCommand): a submenu of the servers, the home server first
const hostOps = [{ op: "resume_on", label: "Resume on Server", busy: "Resuming" }, { op: "migrate", label: "Migrate to Server", busy: "Migrating" },
  { op: "start_on", label: "Start on Server", busy: "Starting" }];
async function runOnHost(n: Node, op: (typeof hostOps)[number], host: string, hostName: string) {
  closeMenus();
  if (!n.server || !n.object) return;
  status.value = `${op.busy} VM '${n.label}' on '${hostName}'...`;
  let error = "";
  try { await api.vmOnHost(n.server, n.object.uuid, host, op.op); } catch (e) { error = `${op.label} '${n.label}': ${String(e)}`; }
  await load(n.server);
  if (error) failure.value = error; else status.value = "Ready";
  reselect(n.id);
}
const onHostItems = (n: Node | null): MenuItem[] => hostOps.map((o) => ({
  label: o.label, enabled: !!n?.object?.allowed?.includes(o.op === "migrate" ? "pool_migrate" : o.op),
  sub: async () => (await api.vmTargets(n!.server!, n!.object!.uuid)).flatMap((t, i) => [
    { label: t.text, icon: t.home ? "000_ServerHome_h32bit_16" : "000_TreeConnected_h32bit_16", enabled: t.enabled, run: () => runOnHost(n!, o, t.uuid, t.text) },
    ...(i === 0 ? [{ label: "", sep: true }] : [])]),
}));
const vmItems = (n: Node | null): MenuItem[] => [
  ...vmCommands.map((c) => ({ label: c.label, icon: `${c.menuIcon ?? c.icon}_h32bit_16`, enabled: canRun(c, n), run: () => runVm(c, n) })),
  ...onHostItems(n),
  { label: "", sep: true },
  { label: "Copy VM...", enabled: canCopy(n), run: () => { closeMenus(); copyOf.value = n; } },
  { label: "", sep: true },
  { label: "Delete VM...", enabled: !!n?.object?.allowed?.includes("destroy"), run: () => { closeMenus(); deleteOf.value = n; } },
  { label: "", sep: true },
  { label: "Properties", icon: "000_Info3_h32bit_16", enabled: !!n?.object, run: () => { closeMenus(); propsOf.value = n; } },
];
const serverNode = computed(() => (selected.value?.server && selected.value.id === selected.value.server ? selected.value : null));
// the menu bar of the C# app; what is not ported yet is shown disabled
const menus = computed<{ label: string; items: MenuItem[] }[]>(() => [
  { label: "File", items: [{ label: "Add New Server...", icon: "000_AddApplicationServer_h32bit_16", enabled: true, run: () => { closeMenus(); openAdd(); } }] },
  { label: "View", items: navButtons.map((b) => b.v === "organization" || b.v === "saved"
    ? { label: b.label, enabled: true, sub: () => viewItems(b.v as "organization" | "saved") }
    : { label: b.label, enabled: true, run: () => { closeMenus(); showView(b.v); } }) },
  { label: "Pool", items: [{ label: "New Pool...", icon: "000_PoolNew_h32bit_16", enabled: servers.value.some((x) => x.connected), run: () => { closeMenus(); newPool.value = true; } }] },
  { label: "Server", items: [
    { label: "Add...", icon: "000_AddApplicationServer_h32bit_16", enabled: true, run: () => { closeMenus(); openAdd(); } },
    { label: "Connect", enabled: !!serverNode.value?.disconnected, run: () => { closeMenus(); connectTo(serverNode.value!); } },
    { label: "Disconnect", enabled: !!serverNode.value && !serverNode.value.disconnected, run: () => { closeMenus(); disconnect(serverNode.value!); } },
    { label: "", sep: true }, ...hostItems(selectedHost.value)] },
  { label: "VM", items: [{ label: "New VM...", icon: "001_CreateVM_h32bit_16", enabled: !!targetServer.value, run: openNewVm }, { label: "", sep: true }, ...vmItems(selectedVm.value)] },
  { label: "Storage", items: [{ label: "New SR...", icon: "000_Storage_h32bit_16", enabled: !!targetServer.value, run: openNewSr }, { label: "", sep: true }, ...srItems(selectedSr.value)] },
  { label: "Templates", items: [{ label: "New VM from Template...", enabled: false }] },
  { label: "Tools", items: [{ label: "Options...", enabled: false }] },
  { label: "Help", items: [{ label: "About XCP-ng Center", enabled: true, run: () => { closeMenus(); choose({ title: "About XCP-ng Center", text: "XCP-ng Center (Go)\nThe XCP-ng Center desktop app in Go, Wails and Vue." }, ["OK"]); } }] },
]);
// the tree's context menu by object
function contextItems(n: Node): MenuItem[] {
  if (n.id === "root") return [{ label: "Add New Server...", icon: "000_AddApplicationServer_h32bit_16", enabled: true, run: () => { closeMenus(); openAdd(); } }];
  if (n.disconnected) return [{ label: "Connect", enabled: true, run: () => { closeMenus(); connectTo(n); } }, { label: "Remove", enabled: true, run: () => { closeMenus(); forget(n); } }];
  if (n.server && n.id === n.server) return [{ label: "Disconnect", enabled: true, run: () => { closeMenus(); disconnect(n); } }];
  if (n.object?.kind === "hosts") return hostItems(n);
  if (n.object?.kind === "vms") return vmItems(n);
  if (n.object?.kind === "srs") return srItems(n);
  return [];
}
function activate(n: Node) { if (n.disconnected) connectTo(n); else toggle(n); }

// the Search tab of the root (SearchPage's overview): every object of the connected servers; of a view's root or group,
// the objects under it
const typeName: Record<string, string> = { pools: "Pool", hosts: "Server", vms: "VM", srs: "Storage", templates: "Template", networks: "Network",
  snapshots: "Snapshot", vdis: "Virtual Disk" };
const overview = computed(() => {
  const n = selected.value;
  const viewing = view.value === "organization" || view.value === "saved";
  if (n && (n.id.startsWith("g:") || (viewing && n.id === "root"))) {
    const rows: { o: XenObject; type: string; key: string; icon: string; id: string }[] = [];
    const seen = new Set<string>();
    const walk = (x: Node) => x.children.forEach((c) => {
      const uid = c.object && c.server ? `${c.server}/${c.object.uuid}` : "";
      if (c.object && c.server && !seen.has(uid)) { seen.add(uid); rows.push({ o: c.object, type: typeName[c.object.kind] ?? c.object.kind, key: c.server, icon: c.icon, id: c.id }); }
      walk(c);
    });
    walk(n);
    return rows;
  }
  return everything.value;
});
// the tree node a row of the Search tab opens: a view's own node, else the Infrastructure tree's
const rowId = (r: { key: string; type: string; o: XenObject; id?: string }) => r.id ?? (r.type === "Pool" ? r.key : `${r.key}/${r.o.uuid}`);
const everything = computed(() => Object.entries(data).flatMap(([key, d]) => [
  ...(d.pools ?? []).map((o) => ({ o: { ...o, name: o.name || (d.hosts ?? []).find((h) => h.coordinator)?.name || key }, type: "Pool", key, icon: "000_PoolConnected_h32bit_16" })),
  ...(d.hosts ?? []).map((o) => ({ o, type: "Server", key, icon: "000_TreeConnected_h32bit_16" })),
  ...(d.vms ?? []).map((o) => ({ o, type: "VM", key, icon: vmIcon(o) })), ...(d.srs ?? []).map((o) => ({ o, type: "Storage", key, icon: o.icon || "000_Storage_h32bit_16" })),
]));

async function refreshServers() { servers.value = await api.servers(); }

async function load(key: string) {
  status.value = `Loading ${key}...`;
  try {
    await reader(key)();
    expanded.add(key);
    status.value = `Connected to ${key}`;
    failure.value = "";
  } catch (e) {
    status.value = "Ready";
    failure.value = `Could not load ${key}: ${String(e)}`;
  }
}
async function read(key: string) {
  const result: Record<string, XenObject[]> = {};
  for (const k of ["pools", "hosts", "vms", "srs"]) result[k] = (await api.objects(key, k)) ?? [];
  return result;
}

// One read of a server's objects at a time (serial), shared by the window's own loads and the live updates: a read
// that ends never overwrites a newer one, and a read of a server disconnected meanwhile (its generation moved) is
// dropped.
const readers = new Map<string, () => Promise<void>>();
const generations = new Map<string, number>();
const generation = (key: string) => generations.get(key) ?? 0;
function reader(key: string) {
  let r = readers.get(key);
  if (!r) {
    r = serial(async () => {
      const gen = generation(key);
      const result = await read(key);
      if (gen === generation(key)) { data[key] = result; followed(key); }
    });
    readers.set(key, r);
  }
  return r;
}
function retire(key: string) { generations.set(key, generation(key) + 1); readers.delete(key); }

// XenObjectDownloader: the server's objects changed (app.Watch, coalesced on the Go side). After each read the tree,
// the open tab and the alerts follow without touching the status bar; the selection stays (or goes to its server
// when its object is gone). A change during the first load joins that load (it reads once more).
const tick = ref(0);
function followed(key: string) {
  if (selected.value?.server === key) {
    if (!reselect(selected.value.id)) reselect(key);
    tick.value++;
  }
  if (view.value === "alerts") void loadAlerts();
  if (view.value === "organization" || view.value === "saved") void readView();
}
async function refresh(key: string) {
  if (!readers.has(key)) return; // never loaded: not shown
  await reader(key)().catch(() => {}); // a failed background read: the next change tries again
}
Events.On("xcp:changed", (ev: { data: unknown }) => refresh(String(Array.isArray(ev.data) ? ev.data[0] : ev.data)));

async function submitAdd() {
  form.busy = true; form.error = "";
  try {
    const r = await api.connect(form.address, form.user, form.password);
    if (r.fingerprint) { cert.value = { key: r.key!, fingerprint: r.fingerprint, previous: r.previous }; return; }
    if (r.error && !r.connected) { form.error = r.error; return; }
    addOpen.value = false; form.password = "";
    await refreshServers(); await load(r.key!);
    void poolDialog.value?.load(); // New Pool is open: the new server is a candidate
    if (selected.value?.id === r.key) reselect(r.key!); // the selected node was the disconnected one
    if (r.error) failure.value = `Connected, but the server list was not saved: ${r.error}`;
  } catch (e) { form.error = String(e); } finally { form.busy = false; }
}

async function acceptCert() {
  const c = cert.value!; cert.value = null;
  try { await api.trust(c.key, c.fingerprint); await submitAdd(); } catch (e) { form.error = String(e); }
}

function toggle(n: Node) { expanded.has(n.id) ? expanded.delete(n.id) : expanded.add(n.id); }

function details(n: Node | null): [string, string][] {
  if (!n) return [];
  if (!n.object) return [["Name", n.label], ["Address", n.server ?? ""]];
  const o = n.object;
  const rows: [string, string][] = [["Name", o.kind === "pools" ? n.label : o.name], ["Description", o.description || "<None>"], ["UUID", o.uuid]];
  if (o.kind === "vms") rows.push(["Power state", o.power_state ?? ""], ["Home server", o.host || "<None>"]);
  return rows;
}

onMounted(async () => { await refreshServers(); selected.value = tree.value; });
</script>

<template>
  <div class="flex h-full flex-col" @click="closeMenus()">
    <!-- the menu bar -->
    <nav class="flex gap-1 border-b border-line bg-pane px-1 py-0.5">
      <span v-for="m in menus" :key="m.label" class="relative cursor-default px-2 py-0.5"
            :class="openMenu === m.label ? 'bg-neutral-200' : 'hover:bg-neutral-200'"
            @click.stop="openMenu = openMenu === m.label ? '' : m.label" @mouseenter="openMenu && (openMenu = m.label)">{{ m.label }}
        <MenuList v-if="openMenu === m.label" class="absolute left-0 top-full z-20" :items="m.items" :hint="soon" />
      </span>
    </nav>
    <!-- the toolbar (MainWindow): Back, Forward | Add New Server | New Pool, New Storage, New VM | the power buttons -->
    <div class="flex items-center gap-0.5 border-b border-line bg-pane px-1 py-0.5">
      <button class="tool" :disabled="historyAt === 0" @click="navigate(-1)"><img src="/icons/001_Back_h32bit_24.png" alt="" />Back</button>
      <button class="tool" :disabled="historyAt >= history.length - 1" @click="navigate(1)"><img src="/icons/001_Forward_h32bit_24.png" alt="" />Forward</button>
      <span class="toolsep"></span>
      <button class="tool" @click="openAdd()"><img src="/icons/000_AddApplicationServer_h32bit_24.png" alt="" />Add New Server</button>
      <span class="toolsep"></span>
      <button class="tool" :disabled="!servers.some((x) => x.connected)" @click="newPool = true"><img src="/icons/000_PoolNew_h32bit_24.png" alt="" />New Pool</button>
      <button class="tool" :disabled="!targetServer" @click="openNewSr"><img src="/icons/000_NewStorage_h32bit_24.png" alt="" />New Storage</button>
      <button class="tool" :disabled="!targetServer" @click="openNewVm"><img src="/icons/000_CreateVM_h32bit_24.png" alt="" />New VM</button>
      <span class="toolsep"></span>
      <button v-for="b in powerButtons" :key="b.op" class="tool" :disabled="!b.enabled" @click="b.run()"><img :src="`/icons/${b.icon}_h32bit_24.png`" alt="" />{{ b.label }}</button>
      <button class="tool" :disabled="!canRun(pauseCommand, selectedVm)" @click="runVm(pauseCommand, selectedVm)">
        <img :src="`/icons/${pauseCommand.icon}_h32bit_24_green.png`" alt="" />{{ pauseCommand.label }}</button>
    </div>
    <div class="flex min-h-0 flex-1">
      <aside class="flex shrink-0 flex-col bg-pane" :style="{ width: `${sideWidth}px` }">
        <div v-if="view !== 'alerts'" class="m-1 flex items-center border border-line bg-white">
          <input v-model="filter" placeholder="Search..." class="min-w-0 flex-1 px-1 py-0.5 outline-none" />
          <img src="/icons/000_Search_h32bit_16.png" alt="" class="mx-1 h-4 w-4" />
        </div>
        <div class="mx-1 min-h-0 flex-1 overflow-auto border border-line bg-white">
          <template v-if="view !== 'alerts'">
            <div v-for="{ node, depth } in flat" :key="node.id"
                 class="flex cursor-default items-center gap-1 whitespace-nowrap px-1 py-px"
                 :class="selected?.id === node.id ? 'bg-select text-white' : 'hover:bg-neutral-100'"
                 :style="{ paddingLeft: `${0.25 + depth * 1.1}rem` }"
                 @click="selected = node" @dblclick="activate(node)"
                 @contextmenu.prevent.stop="selected = node; menu = contextItems(node).length ? { x: $event.clientX, y: $event.clientY, items: contextItems(node) } : null">
              <!-- the C# tree's boxed expander -->
              <span class="flex h-4 w-4 shrink-0 items-center justify-center" @click.stop="toggle(node)">
                <img v-if="node.children.length" :src="`/icons/${expanded.has(node.id) ? 'tree_minus' : 'tree_plus'}.png`" alt="" />
              </span>
              <img :src="`/icons/${node.icon}.png`" alt="" class="h-4 w-4 shrink-0" />
              <span>{{ node.label }}</span>
            </div>
          </template>
          <div v-else-if="view === 'alerts'" class="flex items-center gap-1 bg-select px-2 py-px text-white">
            <img src="/icons/notif_alerts_16.png" alt="" class="h-4 w-4" /> Alerts ({{ alerts.length }})
          </div>
        </div>
        <!-- NavigationPane -->
        <div class="mt-1 border-t border-line">
          <div v-for="b in navButtons" :key="b.v" class="flex cursor-default items-center gap-2 px-2 py-1"
               :class="view === b.v ? 'navactive' : 'hover:bg-neutral-200'" @click.stop="navClick(b.v, $event)">
            <img :src="`/icons/${b.icon}.png`" alt="" class="h-6 w-6" />{{ b.label }}
          </div>
        </div>
      </aside>
      <!-- the splitter between the navigation pane and the tabs: drag it to resize -->
      <div class="splitter" title="Drag to resize" @mousedown.prevent="startResize"></div>
      <main v-if="view === 'alerts'" class="flex min-w-0 flex-1 flex-col">
        <div class="flex items-center gap-2 border-b border-line bg-pane px-2 py-1 font-semibold">
          <img src="/icons/notif_alerts_16.png" alt="" class="h-4 w-4" />Alerts
        </div>
        <div class="flex gap-1 border-b border-line bg-pane px-2 py-1">
          <button class="px-2 py-0.5 hover:bg-neutral-200" @click="loadAlerts">Refresh</button>
          <button class="px-2 py-0.5 hover:bg-neutral-200 disabled:text-neutral-400" :disabled="!alerts.length" @click="dismiss(alerts, true)">Dismiss All</button>
          <button class="px-2 py-0.5 hover:bg-neutral-200 disabled:text-neutral-400" :disabled="!picked.size" @click="dismiss(alerts.filter((a) => picked.has(a.uuid)))">Dismiss Selected</button>
        </div>
        <div class="min-h-0 flex-1 overflow-auto">
          <table class="w-full border-collapse">
            <thead class="sticky top-0 bg-pane text-left">
              <tr><th class="w-20 border-b border-line px-2 py-1 font-normal">Severity</th><th class="border-b border-line px-2 py-1 font-normal">Message</th>
                <th class="w-48 border-b border-line px-2 py-1 font-normal">Location</th><th class="w-44 border-b border-line px-2 py-1 font-normal">Date</th><th class="w-24 border-b border-line"></th></tr>
            </thead>
            <tbody>
              <tr v-for="a in alerts" :key="a.uuid" class="cursor-default align-top" :class="picked.has(a.uuid) ? 'bg-select text-white' : 'hover:bg-neutral-100'" @click="pick(a, $event)">
                <td class="border-b border-line px-2 py-1"><img :src="`/icons/alert${a.priority}_16.png`" alt="" class="h-4 w-4" /></td>
                <td class="border-b border-line px-2 py-1 select-text"><div class="font-semibold">{{ a.title }}</div><div class="whitespace-pre-line">{{ a.description }}</div></td>
                <td class="border-b border-line px-2 py-1">{{ a.applies_to }}</td>
                <td class="border-b border-line px-2 py-1">{{ when(a.time) }}</td>
                <td class="border-b border-line px-2 py-1"><button class="underline" @click.stop="dismiss([a])">Dismiss</button></td>
              </tr>
            </tbody>
          </table>
          <p v-if="!alerts.length" class="p-3 text-neutral-500">There are no alerts.</p>
        </div>
      </main>
      <main v-else class="flex min-w-0 flex-1 flex-col">
        <div class="flex items-center gap-2 border-b border-line bg-pane px-2 py-1 font-semibold">
          <img :src="`/icons/${!selected || selected.id === 'root' ? 'Logo' : selected.icon}.png`" alt="" class="h-4 w-4" />{{ selected?.label ?? "XCP-ng Center" }}
        </div>
        <div class="tabstrip" role="tablist">
          <button v-for="t in tabs" :key="t" role="tab" :aria-selected="tab === t" class="tab" :class="{ active: tab === t }" @click="tab = t">
            <img :src="`/icons/${tabInfo[t].icon}.png`" alt="" class="h-4 w-4" />{{ tabInfo[t].label }}
          </button>
        </div>
        <!-- Views.HomePage: the logo and "Add a Server" -->
        <div v-if="tab === 'home'" class="flex min-h-0 flex-1 flex-col items-center gap-11 overflow-auto bg-white pt-12">
          <img src="/icons/home-logo.png" alt="" class="h-72 w-72" />
          <a class="cursor-pointer font-serif text-[15px] text-blue-700 underline" @click="openAdd()">Add a Server</a>
        </div>
        <div v-else-if="tab === 'search'" class="min-h-0 flex-1 overflow-auto">
          <div class="bg-gradient-to-r from-band-dark to-band px-3 py-1.5 font-semibold text-white">Overview</div>
          <table class="w-full border-collapse">
            <thead class="sticky top-0 bg-pane text-left"><tr>
              <th class="border-b border-line px-2 py-1 font-normal">Name</th><th class="border-b border-line px-2 py-1 font-normal">Type</th>
              <th class="border-b border-line px-2 py-1 font-normal">Description</th><th class="border-b border-line px-2 py-1 font-normal">Server</th></tr></thead>
            <tbody>
              <tr v-for="r in overview" :key="`${r.key}/${r.o.uuid}`" class="cursor-default hover:bg-neutral-100" @dblclick="reselect(rowId(r))">
                <td class="border-b border-line px-2 py-1"><span class="flex items-center gap-1.5"><img :src="`/icons/${r.icon}.png`" alt="" class="h-4 w-4" />{{ r.o.name }}</span></td><td class="border-b border-line px-2 py-1">{{ r.type }}</td>
                <td class="border-b border-line px-2 py-1">{{ r.o.description }}</td><td class="border-b border-line px-2 py-1">{{ r.key }}</td></tr>
            </tbody>
          </table>
          <p v-if="!overview.length" class="p-3 text-neutral-500">No connected server. Use Add New Server, or double-click a server of the tree to connect.</p>
        </div>
        <ConsoleView v-else-if="tab === 'console' && consoleVm && selected?.server" :server="selected.server" :vm="consoleVm" />
        <PerformanceView v-else-if="tab === 'performance' && selected?.object && selected.server" :key="`perf/${selected.id}`"
                         :server="selected.server" :uuid="selected.object.uuid" :kind="selected.object.kind as 'vms' | 'hosts'" />
        <VmTab v-else-if="(tab === 'memory' || tab === 'storage' || tab === 'networking') && selected?.object && selected.server"
               :key="tab + '/' + selected.server + '/' + selected.object.uuid" :server="selected.server" :vm="selected.object.uuid" :kind="tab" :tick="tick"
               @status="onTabStatus" />
        <!-- keyed: another VM gets a new tab, never the previous VM's rows or selection -->
        <SnapshotsView v-else-if="tab === 'snapshots' && selected?.object && selected.server" :key="selected.server + '/' + selected.object.uuid"
                       :server="selected.server" :vm="selected.object.uuid" :tick="tick"
                       @status="onTabStatus" @changed="selected?.server && load(selected.server).then(() => reselect(selected!.id))" />
        <template v-else>
        <div class="bg-gradient-to-r from-band-dark to-band px-3 py-1.5 font-semibold text-white">{{ selected?.disconnected ? "Connection General Properties" : "General Properties" }}</div>
        <div class="overflow-auto p-3">
          <table v-if="selected" class="w-full max-w-3xl border border-line">
            <tr v-for="[k, v] in details(selected)" :key="k" class="border-b border-line">
              <td class="w-48 px-2 py-1 text-neutral-600">{{ k }}:</td>
              <td class="px-2 py-1 select-text">{{ v }}</td>
            </tr>
          </table>
          <p v-if="selected?.disconnected" class="mt-3"><button class="btn" @click="connectTo(selected!)">Connect</button></p>
        </div>
        </template>
      </main>
    </div>
    <footer class="flex gap-4 border-t border-line bg-pane px-2 py-0.5 text-xs">
      <span>{{ status }}</span><span v-if="failure" class="text-red-700">{{ failure }}</span>
    </footer>

    <MenuList v-if="viewMenu && viewMenu.items.length" class="fixed z-20" :style="{ left: `${viewMenu.x}px`, top: `${viewMenu.y}px` }" :items="viewMenu.items" />
    <MenuList v-if="menu" class="fixed z-20" :style="{ left: `${menu.x}px`, top: `${menu.y}px` }" :items="menu.items" />

    <NewVmWizard v-if="newVmOn" :server="newVmOn" @close="newVmClosed" />
    <NewPoolDialog v-if="newPool" ref="poolDialog" @close="newPoolClosed" @add-server="openAdd()" />
    <CopyVmDialog v-if="copyOf?.server && copyOf.object" :server="copyOf.server" :vm="copyOf.object.uuid" @close="copyClosed" />
    <DeleteVmDialog v-if="deleteOf?.server && deleteOf.object" :server="deleteOf.server" :vm="deleteOf.object.uuid" @close="deleteClosed" />
    <NewSrWizard v-if="newSrOn" :server="newSrOn" :pool="poolLabel(newSrOn)" @close="newSrClosed" />
    <PropertiesDialog v-if="propsOf?.object && propsOf.server" :server="propsOf.server" :vm="propsOf.object.uuid" @close="propertiesClosed" />
    <div v-if="ask" class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
      <div class="w-[30rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ ask.title }}</div>
        <div class="space-y-3 p-3">
          <p class="whitespace-pre-line">{{ ask.text }}</p>
          <div class="flex justify-end gap-2">
            <button v-for="(b, i) in ask.buttons" :key="b" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="answer(i)">{{ b }}</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="addOpen" class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
      <form class="w-[28rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="submitAdd">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Add New Server</div>
        <div class="space-y-2 p-3">
          <p>Enter the host name or IP address of the server you want to add and your user login credentials for that server.</p>
          <label class="grid grid-cols-[6rem_1fr] items-center gap-2">Server:<input ref="serverField" v-model="form.address" class="border border-line px-1 py-0.5" /></label>
          <label class="grid grid-cols-[6rem_1fr] items-center gap-2">User name:<input v-model="form.user" class="border border-line px-1 py-0.5" /></label>
          <label class="grid grid-cols-[6rem_1fr] items-center gap-2">Password:<input ref="passwordField" v-model="form.password" type="password" class="border border-line px-1 py-0.5" /></label>
          <p v-if="form.error" class="text-red-600">{{ form.error }}</p>
          <div class="flex justify-end gap-2 pt-1">
            <button type="submit" :disabled="form.busy" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5">Add</button>
            <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="addOpen = false">Cancel</button>
          </div>
        </div>
      </form>
    </div>

    <div v-if="cert" class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
      <div class="w-[34rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ cert.previous ? "Certificate Changed" : "Unknown Certificate" }}</div>
        <div class="space-y-2 p-3">
          <p v-if="!cert.previous">The certificate of {{ cert.key }} is not trusted yet. Accept it only if this SHA-256 fingerprint is the server's:</p>
          <p v-else class="text-red-700">The certificate of {{ cert.key }} changed. It is refused: verify the server before trusting it.</p>
          <code class="block break-all bg-pane p-2 select-text">{{ cert.fingerprint.match(/../g)?.join(":") }}</code>
          <div class="flex justify-end gap-2 pt-1">
            <button v-if="!cert.previous" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="acceptCert">Accept</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="cert = null; form.busy = false">Cancel</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
