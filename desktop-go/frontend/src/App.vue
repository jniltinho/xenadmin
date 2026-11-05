<script setup lang="ts">
// The main window (first stage): menu and toolbar, the Infrastructure tree (server -> pool -> hosts -> VMs, and the
// SRs), the General tab of the selected object, Add New Server and the certificate question, a status bar.
import { computed, onMounted, reactive, ref } from "vue";
import { api, type Alert, type Confirm, type ServerRow, type XenObject } from "./api/workspace";

interface Node { id: string; label: string; icon: string; object?: XenObject; server?: string; children: Node[] }

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
const selectedVm = computed(() => (selected.value?.object?.kind === "vms" ? selected.value : null));
const canRun = (c: VmCommand, n: Node | null) => !!n?.object?.allowed?.includes(c.needs);
const menu = ref<{ x: number; y: number; node: Node } | null>(null);
const vmMenuOpen = ref(false);
// the message box: a title, a text and its buttons; resolves with the index of the button clicked (-1: closed)
const ask = ref<(Confirm & { buttons: string[]; resolve: (i: number) => void }) | null>(null);
const choose = (q: Confirm, buttons: string[]) => new Promise<number>((resolve) => (ask.value = { ...q, buttons, resolve }));

async function runVm(c: VmCommand, n: Node | null) {
  menu.value = null; vmMenuOpen.value = false;
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
function reselect(id: string) {
  const path = (n: Node): Node[] | null => {
    if (n.id === id) return [n];
    for (const c of n.children) { const p = path(c); if (p) return [n, ...p]; }
    return null;
  };
  const p = path(tree.value);
  p?.slice(0, -1).forEach((a) => expanded.add(a.id));
  selected.value = p?.[p.length - 1] ?? null;
}

// Notifications > Alerts (AlertSummaryPage): every connected server's alerts, newest first; Dismiss asks first.
const view = ref<"infrastructure" | "alerts">("infrastructure");
type AlertRow = Alert & { server: string };
const alerts = ref<AlertRow[]>([]);
const picked = reactive(new Set<string>());
const confirm = async (q: Confirm) => (await choose(q, ["Yes", "No"])) === 0;

async function loadAlerts() {
  const rows: AlertRow[] = [];
  for (const s of servers.value.filter((x) => x.connected)) {
    try {
      rows.push(...((await api.alerts(s.key)) ?? []).map((a) => ({ ...a, server: s.key })));
    } catch (e) {
      failure.value = `Could not load the alerts of ${s.key}: ${String(e)}`;
    }
  }
  alerts.value = rows.sort((a, b) => b.time.localeCompare(a.time));
  picked.clear();
}

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

async function showView(v: "infrastructure" | "alerts") {
  view.value = v;
  if (v === "alerts") await loadAlerts();
}

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
const hostCanRun = (c: HostCommand, n: Node | null) =>
  !!n?.object?.live && (c.op === "enter_maintenance" ? !n.object.maintenance : c.op === "exit_maintenance" ? !!n.object.maintenance : true);
const serverMenuOpen = ref(false);

async function runHost(c: HostCommand, n: Node | null) {
  menu.value = null; vmMenuOpen.value = false; serverMenuOpen.value = false;
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

const vmIcon = (o: XenObject) => (o.power_state === "Running" ? "000_StartVM_h32bit_16" : o.power_state === "Suspended" ? "000_SuspendVM_h32bit_16" : "000_StoppedVM_h32bit_16");

const tree = computed<Node>(() => ({
  id: "root", label: "XCP-ng Center", icon: "Logo", children: servers.value.map((s) => {
    const d = data[s.key];
    if (!s.connected || !d) return { id: s.key, label: s.name || s.host, icon: "000_ServerDisconnected_h32bit_16", server: s.key, children: [] };
    const pool = d.pools?.[0];
    const hosts = (d.hosts ?? []).map<Node>((h) => ({
      id: h.uuid, label: h.name, icon: !h.live ? "000_ServerDisconnected_h32bit_16" : h.maintenance ? "000_ServerMaintenance_h32bit_16" : "000_TreeConnected_h32bit_16", object: h, server: s.key,
      children: (d.vms ?? []).filter((v) => v.host_uuid === h.uuid).map((v) => ({ id: v.uuid, label: v.name, icon: vmIcon(v), object: v, server: s.key, children: [] })),
    }));
    const loose = (d.vms ?? []).filter((v) => !v.host_uuid).map<Node>((v) => ({ id: v.uuid, label: v.name, icon: vmIcon(v), object: v, server: s.key, children: [] }));
    const srs = (d.srs ?? []).map<Node>((r) => ({ id: r.uuid, label: r.name, icon: "000_Storage_h32bit_16", object: r, server: s.key, children: [] }));
    // Pool.Name: an unnamed pool is called after its coordinator
    const coordinator = (d.hosts ?? []).find((h) => h.coordinator);
    return { id: s.key, label: pool?.name || coordinator?.name || s.name || s.host, icon: "000_PoolConnected_h32bit_16", object: pool, server: s.key, children: [...hosts, ...loose, ...srs] };
  }),
}));

const flat = computed(() => {
  const rows: { node: Node; depth: number }[] = [];
  const walk = (n: Node, depth: number) => {
    rows.push({ node: n, depth });
    if (expanded.has(n.id)) n.children.forEach((c) => walk(c, depth + 1));
  };
  walk(tree.value, 0);
  return rows;
});

async function refreshServers() { servers.value = await api.servers(); }

async function load(key: string) {
  status.value = `Loading ${key}...`;
  try {
    const result: Record<string, XenObject[]> = {};
    for (const k of ["pools", "hosts", "vms", "srs"]) result[k] = (await api.objects(key, k)) ?? [];
    data[key] = result;
    expanded.add(key);
    status.value = `Connected to ${key}`;
    failure.value = "";
  } catch (e) {
    status.value = "Ready";
    failure.value = `Could not load ${key}: ${String(e)}`;
  }
}

async function submitAdd() {
  form.busy = true; form.error = "";
  try {
    const r = await api.connect(form.address, form.user, form.password);
    if (r.fingerprint) { cert.value = { key: r.key!, fingerprint: r.fingerprint, previous: r.previous }; return; }
    if (r.error && !r.connected) { form.error = r.error; return; }
    addOpen.value = false; form.password = "";
    await refreshServers(); await load(r.key!);
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
  if (!n.object) return [["Name", n.label], ["Address", n.server ?? ""], ["State", "Disconnected"]];
  const o = n.object;
  const rows: [string, string][] = [["Name", o.name], ["Description", o.description || "<None>"], ["UUID", o.uuid]];
  if (o.kind === "vms") rows.push(["Power state", o.power_state ?? ""], ["Home server", o.host || "<None>"]);
  return rows;
}

onMounted(async () => { await refreshServers(); });
</script>

<template>
  <div class="flex h-full flex-col" @click="menu = null">
    <nav class="flex gap-4 border-b border-line bg-pane px-2 py-1">
      <span class="cursor-default hover:bg-neutral-200 px-1" @click="addOpen = true">File</span>
      <span v-for="m in ['View', 'Pool']" :key="m" class="cursor-default px-1 text-neutral-400">{{ m }}</span>
      <span class="relative cursor-default px-1 hover:bg-neutral-200" @click.stop="serverMenuOpen = !serverMenuOpen; vmMenuOpen = false">Server
        <div v-if="serverMenuOpen" class="absolute left-0 top-full z-10 w-56 border border-neutral-400 bg-white py-1 shadow">
          <div v-for="c in hostCommands" :key="c.op" class="flex items-center gap-2 px-2 py-0.5"
               :class="hostCanRun(c, selectedHost) ? 'hover:bg-select hover:text-white' : 'text-neutral-400'" @click.stop="runHost(c, selectedHost)">
            <img :src="`/icons/${c.icon}_h32bit_16.png`" alt="" class="h-4 w-4" :class="{ 'opacity-40 grayscale': !hostCanRun(c, selectedHost) }" />{{ c.label }}
          </div>
        </div>
      </span>
      <span class="relative cursor-default px-1 hover:bg-neutral-200" @click.stop="vmMenuOpen = !vmMenuOpen; serverMenuOpen = false">VM
        <div v-if="vmMenuOpen" class="absolute left-0 top-full z-10 w-48 border border-neutral-400 bg-white py-1 shadow">
          <div v-for="c in vmCommands" :key="c.op" class="flex items-center gap-2 px-2 py-0.5"
               :class="canRun(c, selectedVm) ? 'hover:bg-select hover:text-white' : 'text-neutral-400'" @click.stop="runVm(c, selectedVm)">
            <img :src="`/icons/${c.menuIcon ?? c.icon}_h32bit_16.png`" alt="" class="h-4 w-4" :class="{ 'opacity-40 grayscale': !canRun(c, selectedVm) }" />{{ c.label }}
          </div>
        </div>
      </span>
      <span v-for="m in ['Storage', 'Templates', 'Tools', 'Help']" :key="m" class="cursor-default px-1 text-neutral-400">{{ m }}</span>
    </nav>
    <div class="flex items-center gap-1 border-b border-line bg-pane px-2 py-1">
      <button class="flex items-center gap-1 px-2 py-0.5 hover:bg-neutral-200" @click="addOpen = true">
        <img src="/icons/000_AddApplicationServer_h32bit_16.png" alt="" /> Add New Server
      </button>
      <template v-for="c in hostCommands.slice(0, 2)" :key="c.op">
        <button v-if="hostCanRun(c, selectedHost)" class="flex items-center gap-1 px-2 py-0.5 hover:bg-neutral-200" @click="runHost(c, selectedHost)">
          <img :src="`/icons/${c.icon}_h32bit_24.png`" alt="" class="h-6 w-6" /> {{ c.label }}
        </button>
      </template>
      <template v-for="c in vmCommands" :key="c.op">
        <button v-if="canRun(c, selectedVm)" class="flex items-center gap-1 px-2 py-0.5 hover:bg-neutral-200" @click="runVm(c, selectedVm)">
          <img :src="`/icons/${c.icon}_h32bit_24.png`" alt="" class="h-6 w-6" /> {{ c.label }}
        </button>
      </template>
    </div>
    <div class="flex min-h-0 flex-1">
      <aside class="flex w-64 flex-col border-r border-line">
       <div class="min-h-0 flex-1 overflow-auto">
        <template v-if="view === 'infrastructure'">
        <div v-for="{ node, depth } in flat" :key="node.id"
             class="flex cursor-default items-center gap-1 whitespace-nowrap px-1 py-px"
             :class="selected?.id === node.id ? 'bg-select text-white' : 'hover:bg-neutral-100'"
             :style="{ paddingLeft: `${0.25 + depth * 1.1}rem` }"
             @click="selected = node" @dblclick="toggle(node)"
             @contextmenu.prevent.stop="selected = node; menu = node.object?.kind === 'vms' || node.object?.kind === 'hosts' ? { x: $event.clientX, y: $event.clientY, node } : null">
          <span class="w-3 text-center text-[10px]" @click.stop="toggle(node)">{{ node.children.length ? (expanded.has(node.id) ? "−" : "+") : "" }}</span>
          <img :src="`/icons/${node.icon}.png`" alt="" class="h-4 w-4" />
          <span>{{ node.label }}</span>
        </div>
        </template>
        <div v-else class="flex items-center gap-1 bg-select px-2 py-px text-white">
          <img src="/icons/notif_alerts_16.png" alt="" class="h-4 w-4" /> Alerts ({{ alerts.length }})
        </div>
       </div>
        <div class="border-t border-line bg-pane">
          <div v-for="b in [{ v: 'infrastructure', label: 'Infrastructure', icon: 'infra_view_24' }, { v: 'alerts', label: 'Notifications', icon: 'notif_none_24' }] as const"
               :key="b.v" class="flex cursor-default items-center gap-2 border-b border-line px-2 py-1.5 font-semibold"
               :class="view === b.v ? 'bg-neutral-200' : 'hover:bg-neutral-100'" @click="showView(b.v)">
            <img :src="`/icons/${b.icon}.png`" alt="" class="h-6 w-6" />{{ b.label }}
          </div>
        </div>
      </aside>
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
          <img :src="`/icons/${selected?.icon ?? 'Logo'}.png`" alt="" class="h-4 w-4" />{{ selected?.label ?? "XCP-ng Center" }}
        </div>
        <div class="flex border-b border-line bg-pane px-2 pt-1">
          <span class="border border-b-0 border-line bg-white px-3 py-0.5">General</span>
        </div>
        <div class="bg-gradient-to-r from-band-dark to-band px-3 py-1.5 font-semibold text-white">General Properties</div>
        <div class="overflow-auto p-3">
          <table v-if="selected" class="w-full max-w-3xl border border-line">
            <tr v-for="[k, v] in details(selected)" :key="k" class="border-b border-line">
              <td class="w-48 px-2 py-1 text-neutral-600">{{ k }}:</td>
              <td class="px-2 py-1 select-text">{{ v }}</td>
            </tr>
          </table>
          <p v-else class="text-neutral-500">Select a server, host, VM or SR. Use Add New Server to connect.</p>
        </div>
      </main>
    </div>
    <footer class="flex gap-4 border-t border-line bg-pane px-2 py-0.5 text-xs">
      <span>{{ status }}</span><span v-if="failure" class="text-red-700">{{ failure }}</span>
    </footer>

    <div v-if="menu" class="fixed z-20 w-60 border border-neutral-400 bg-white py-1 shadow" :style="{ left: `${menu.x}px`, top: `${menu.y}px` }" @click.stop>
      <template v-if="menu.node.object?.kind === 'hosts'">
        <div v-for="c in hostCommands" :key="c.op" class="flex items-center gap-2 px-2 py-0.5"
             :class="hostCanRun(c, menu.node) ? 'hover:bg-select hover:text-white' : 'text-neutral-400'" @click="runHost(c, menu.node)">
          <img :src="`/icons/${c.icon}_h32bit_16.png`" alt="" class="h-4 w-4" :class="{ 'opacity-40 grayscale': !hostCanRun(c, menu.node) }" />{{ c.label }}
        </div>
      </template>
      <div v-for="c in menu.node.object?.kind === 'vms' ? vmCommands : []" :key="c.op" class="flex items-center gap-2 px-2 py-0.5"
           :class="canRun(c, menu.node) ? 'hover:bg-select hover:text-white' : 'text-neutral-400'" @click="runVm(c, menu.node)">
        <img :src="`/icons/${c.menuIcon ?? c.icon}_h32bit_16.png`" alt="" class="h-4 w-4" :class="{ 'opacity-40 grayscale': !canRun(c, menu.node) }" />{{ c.label }}
      </div>
    </div>

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

    <div v-if="addOpen" class="fixed inset-0 flex items-center justify-center bg-black/20">
      <form class="w-[28rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="submitAdd">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Add New Server</div>
        <div class="space-y-2 p-3">
          <p>Enter the host name or IP address of the server you want to add and your user login credentials for that server.</p>
          <label class="grid grid-cols-[6rem_1fr] items-center gap-2">Server:<input v-model="form.address" class="border border-line px-1 py-0.5" autofocus /></label>
          <label class="grid grid-cols-[6rem_1fr] items-center gap-2">User name:<input v-model="form.user" class="border border-line px-1 py-0.5" /></label>
          <label class="grid grid-cols-[6rem_1fr] items-center gap-2">Password:<input v-model="form.password" type="password" class="border border-line px-1 py-0.5" /></label>
          <p v-if="form.error" class="text-red-600">{{ form.error }}</p>
          <div class="flex justify-end gap-2 pt-1">
            <button type="submit" :disabled="form.busy" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5">Add</button>
            <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="addOpen = false">Cancel</button>
          </div>
        </div>
      </form>
    </div>

    <div v-if="cert" class="fixed inset-0 flex items-center justify-center bg-black/20">
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
