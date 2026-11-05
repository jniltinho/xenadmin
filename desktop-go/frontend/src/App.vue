<script setup lang="ts">
// The main window (first stage): menu and toolbar, the Infrastructure tree (server -> pool -> hosts -> VMs, and the
// SRs), the General tab of the selected object, Add New Server and the certificate question, a status bar.
import { computed, onMounted, reactive, ref } from "vue";
import { api, type ServerRow, type XenObject } from "./api/workspace";

interface Node { id: string; label: string; icon: string; object?: XenObject; server?: string; children: Node[] }

const servers = ref<ServerRow[]>([]);
const data = reactive<Record<string, Record<string, XenObject[]>>>({});
const selected = ref<Node | null>(null);
const expanded = reactive(new Set<string>(["root"]));
const status = ref("Ready");
const failure = ref(""); // the last error, kept visible in the status bar until the next success

const addOpen = ref(false);
const form = reactive({ address: "", user: "root", password: "", error: "", busy: false });
const cert = ref<{ key: string; fingerprint: string; previous?: string } | null>(null);

const vmIcon = (o: XenObject) => (o.power_state === "Running" ? "000_StartVM_h32bit_16" : o.power_state === "Suspended" ? "000_SuspendVM_h32bit_16" : "000_StoppedVM_h32bit_16");

const tree = computed<Node>(() => ({
  id: "root", label: "XCP-ng Center", icon: "Logo", children: servers.value.map((s) => {
    const d = data[s.key];
    if (!s.connected || !d) return { id: s.key, label: s.name || s.host, icon: "000_ServerDisconnected_h32bit_16", server: s.key, children: [] };
    const pool = d.pools?.[0];
    const hosts = (d.hosts ?? []).map<Node>((h) => ({
      id: h.uuid, label: h.name, icon: "000_TreeConnected_h32bit_16", object: h, server: s.key,
      children: (d.vms ?? []).filter((v) => v.host_uuid === h.uuid).map((v) => ({ id: v.uuid, label: v.name, icon: vmIcon(v), object: v, server: s.key, children: [] })),
    }));
    const loose = (d.vms ?? []).filter((v) => !v.host_uuid).map<Node>((v) => ({ id: v.uuid, label: v.name, icon: vmIcon(v), object: v, server: s.key, children: [] }));
    const srs = (d.srs ?? []).map<Node>((r) => ({ id: r.uuid, label: r.name, icon: "000_Storage_h32bit_16", object: r, server: s.key, children: [] }));
    return { id: s.key, label: pool?.name || s.name || s.host, icon: "000_PoolConnected_h32bit_16", object: pool, server: s.key, children: [...hosts, ...loose, ...srs] };
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
  <div class="flex h-full flex-col">
    <nav class="flex gap-4 border-b border-line bg-pane px-2 py-1">
      <span class="cursor-default hover:bg-neutral-200 px-1" @click="addOpen = true">File</span>
      <span v-for="m in ['View', 'Pool', 'Server', 'VM', 'Storage', 'Templates', 'Tools', 'Help']" :key="m" class="cursor-default px-1 text-neutral-400">{{ m }}</span>
    </nav>
    <div class="flex items-center gap-1 border-b border-line bg-pane px-2 py-1">
      <button class="flex items-center gap-1 px-2 py-0.5 hover:bg-neutral-200" @click="addOpen = true">
        <img src="/icons/000_AddApplicationServer_h32bit_16.png" alt="" /> Add New Server
      </button>
    </div>
    <div class="flex min-h-0 flex-1">
      <aside class="w-64 overflow-auto border-r border-line">
        <div v-for="{ node, depth } in flat" :key="node.id"
             class="flex cursor-default items-center gap-1 whitespace-nowrap px-1 py-px"
             :class="selected?.id === node.id ? 'bg-select text-white' : 'hover:bg-neutral-100'"
             :style="{ paddingLeft: `${0.25 + depth * 1.1}rem` }"
             @click="selected = node" @dblclick="toggle(node)">
          <span class="w-3 text-center text-[10px]" @click.stop="toggle(node)">{{ node.children.length ? (expanded.has(node.id) ? "−" : "+") : "" }}</span>
          <img :src="`/icons/${node.icon}.png`" alt="" class="h-4 w-4" />
          <span>{{ node.label }}</span>
        </div>
      </aside>
      <main class="flex min-w-0 flex-1 flex-col">
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
