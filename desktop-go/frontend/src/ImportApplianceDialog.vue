<script setup lang="ts">
// Import Appliance (XenAdmin.Desktop's Dialogs.ImportAppliance: the import wizard's OVF/OVA path in one dialog): the
// package's VMs, its EULA to accept, the pool or server, then the wizard's VmMappings: the server of each VM
// (ImportSelectHostPage, only to put its disks on local storage), all the disks on one SR or each on its own
// (ImportSelectStoragePage), the network of each interface (ImportSelectNetworkPage); manifest verification and Start,
// then ImportApplianceAction. The defaults (the target's server, its default SR, the first network) keep it one click.
// shortcut: no OS fixup page (ImportOptionsPage): see rules/ovfmap.go.
import { computed, onMounted, ref, watch } from "vue";
import { api, type Appliance, type ApplianceMapping, type ServerRow } from "./api/workspace";

const props = defineProps<{ file: string }>();
const emit = defineEmits<{ (e: "close"): void; (e: "run", label: string, key: string, f: () => Promise<string>): void }>();

type Target = { key: string; host: string; text: string };
type Options = { srs: { uuid: string; text: string }[]; default_sr: string; networks: { uuid: string; name: string }[] };
const info = ref<Appliance | null>(null);
const loadError = ref("");
const targets = ref<Target[]>([]);
const target = ref(0);
const servers = ref<Record<string, { uuid: string; name: string }[]>>({}); // per connection: the pool, then its servers
const hosts = ref<string[]>([]); // the server of each VM ("" for the pool)
const opts = ref<Record<string, Options>>({}); // the SRs and networks seen from each server of the chosen connection
const sameSr = ref(true);
const sr = ref("");
const diskSrs = ref<string[][]>([]);
const nets = ref<string[][]>([]);
const accept = ref(false);
const verify = ref(false);
const start = ref(false);
const error = ref("");
const pending = ref(false);

const vms = computed(() => info.value?.vms ?? []);
const conn = computed(() => targets.value[target.value]);
const poolServers = computed(() => (conn.value ? servers.value[conn.value.key] ?? [] : []));
const eulas = computed(() => info.value?.eulas ?? []);
const summary = computed(() => {
  const a = info.value;
  if (!a) return "";
  return `Appliance '${a.name}': ${vms.value.length} VM(s) (${vms.value.map((v) => v.name).join(", ")}), ${a.disks} disk(s), ${a.interfaces} network interface(s).`;
});
// "all on this SR": the SRs every chosen server sees (SelectVMStorageWithMultipleVirtualDisksPage's commonSRs)
const commonSrs = computed(() => {
  const all = [...new Set(hosts.value)].map((h) => opts.value[h]);
  if (all.length === 0 || all.some((o) => !o)) return [];
  return all[0].srs.filter((s) => all.every((o) => o.srs.some((x) => x.uuid === s.uuid)));
});
const size = (b: number) => {
  const units = ["B", "KB", "MB", "GB", "TB"];
  let i = 0;
  while (b >= 1024 && i < units.length - 1) { b /= 1024; i++; }
  return `${b.toFixed(i ? 1 : 0)} ${units[i]}`;
};

onMounted(async () => {
  try { info.value = await api.applianceInfo(props.file); } catch (e) { loadError.value = String(e); return; } // ImportSourcePage: read once
  try {
    const rows: ServerRow[] = (await api.servers()).filter((s) => s.connected);
    const out: Target[] = [];
    for (const s of rows) {
      const o = await api.importOptions(s.key, "");
      servers.value[s.key] = o.servers;
      o.servers.forEach((h, i) => out.push({ key: s.key, host: h.uuid, text: i === 0 ? h.name : "    " + h.name }));
    }
    targets.value = out;
  } catch (e) { error.value = String(e); }
});
// a new target: every VM on its server, the options read again
watch([targets, target], () => { hosts.value = vms.value.map(() => conn.value?.host ?? ""); opts.value = {}; });
// the options of each VM's server, then each row keeps its choice when still offered, else the default
let asked = 0; // a slow answer for a target or servers no longer chosen is dropped
watch([conn, () => hosts.value.join("\n")], async () => {
  const t = conn.value, id = ++asked;
  if (!t) return;
  for (const h of new Set(hosts.value)) {
    if (opts.value[h]) continue;
    try {
      const o = await api.importOptions(t.key, h);
      if (id !== asked) return;
      opts.value = { ...opts.value, [h]: { srs: o.srs ?? [], default_sr: o.default_sr, networks: o.networks ?? [] } };
    } catch (e) { error.value = String(e); return; }
  }
  vms.value.forEach((v, i) => {
    const o = opts.value[hosts.value[i]];
    const def = o?.srs.some((s) => s.uuid === o.default_sr) ? o.default_sr : o?.srs[0]?.uuid ?? "";
    diskSrs.value[i] = (v.disks ?? []).map((_, j) => { const c = diskSrs.value[i]?.[j]; return c && o?.srs.some((s) => s.uuid === c) ? c : def; });
    nets.value[i] = (v.macs ?? []).map((_, j) => { const c = nets.value[i]?.[j]; return c && o?.networks.some((n) => n.uuid === c) ? c : o?.networks[0]?.uuid ?? ""; });
  });
  if (!commonSrs.value.some((s) => s.uuid === sr.value)) {
    const d = opts.value[hosts.value[0]]?.default_sr;
    sr.value = commonSrs.value.some((s) => s.uuid === d) ? d! : commonSrs.value[0]?.uuid ?? "";
  }
});

async function ok() {
  const a = info.value;
  if (pending.value || !a) return;
  pending.value = true; // the form is frozen: the values checked are the values imported
  const t = conn.value, check = verify.value && a.has_manifest, run = start.value;
  const maps: ApplianceMapping[] = vms.value.map((v, i) => ({ sys_id: v.sys_id, host: hosts.value[i] ?? "",
    srs: (v.disks ?? []).map((_, j) => (sameSr.value ? sr.value : diskSrs.value[i]?.[j] ?? "")),
    networks: (v.macs ?? []).map((_, j) => nets.value[i]?.[j] ?? "") }));
  error.value = "";
  if (eulas.value.length > 0 && !accept.value) error.value = "You must accept the EULA(s) to import this appliance.";
  else if (!t) error.value = "Select the pool or server to import to.";
  else if (maps.some((m) => m.srs.includes(""))) error.value = "Select a storage repository.";
  else if (maps.some((m) => m.networks.includes(""))) error.value = "Select a network.";
  if (!error.value) {
    try { error.value = await api.applianceMappingsError(t.key, a, maps); } catch (e) { error.value = String(e); }
  }
  if (error.value) { pending.value = false; return; }
  emit("run", `Importing '${a.name}'`, t.key, () => api.importAppliance(t.key, props.file, maps, check, run));
  emit("close");
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <form novalidate class="w-[42rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="ok">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">Import Appliance</div>
      <div v-if="loadError" class="space-y-2 p-3">
        <p class="whitespace-pre-wrap">{{ loadError }}</p>
        <div class="flex justify-end pt-1">
          <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">OK</button>
        </div>
      </div>
      <p v-else-if="!info" class="p-3">Reading '{{ file.split("/").pop() }}'...</p>
      <fieldset v-else :disabled="pending" class="max-h-[85vh] space-y-2 overflow-auto p-3">
        <p>{{ summary }}</p>
        <p v-if="info.warnings?.length" class="whitespace-pre-wrap text-yellow-700">{{ info.warnings.join("\n") }}</p>
        <template v-if="eulas.length">
          <span class="block font-semibold">End User License Agreement:</span>
          <div class="max-h-36 overflow-auto whitespace-pre-wrap border border-line p-1.5">{{ eulas.join("\n\n") }}</div>
          <label class="flex items-center gap-2"><input v-model="accept" type="checkbox" />I accept the EULA(s)</label>
        </template>
        <label for="apptarget" class="block">Import to:</label>
        <select id="apptarget" v-model.number="target" class="w-full border border-line bg-white px-1 py-0.5 whitespace-pre">
          <option v-for="(t, i) in targets" :key="t.key + t.host" :value="i">{{ t.text }}</option>
        </select>
        <template v-if="poolServers.length > 1">
          <span class="block">If you want to place the disks of a VM on local storage, specify the server below.</span>
          <div class="max-h-28 overflow-auto border border-line">
            <div v-for="(v, i) in vms" :key="v.sys_id + i" class="flex items-center gap-2 px-1 py-0.5">
              <span class="w-1/2 truncate">{{ v.name }}</span>
              <select v-model="hosts[i]" :aria-label="`Server for ${v.name}`" class="w-1/2 border border-line bg-white px-1">
                <option v-for="h in poolServers" :key="h.uuid" :value="h.uuid">{{ h.name }}</option>
              </select>
            </div>
          </div>
        </template>
        <label class="flex items-center gap-2"><input v-model="sameSr" type="radio" :value="true" />Place all imported virtual disks on this target SR:</label>
        <div v-if="sameSr" class="h-24 overflow-auto border border-line">
          <div v-for="s in commonSrs" :key="s.uuid" class="cursor-default whitespace-pre px-1 py-0.5" :class="sr === s.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100'"
               @click="sr = s.uuid">{{ s.text }}</div>
        </div>
        <label class="flex items-center gap-2"><input v-model="sameSr" type="radio" :value="false" />Place imported virtual disks onto specified target SRs:</label>
        <div v-if="!sameSr" class="max-h-32 overflow-auto border border-line">
          <template v-for="(v, i) in vms" :key="v.sys_id + i">
            <div v-for="(d, j) in diskSrs[i] ? v.disks ?? [] : []" :key="j" class="flex items-center gap-2 px-1 py-0.5">
              <span class="w-1/2 truncate">{{ v.name }} - {{ d.name }} ({{ size(d.capacity) }})</span>
              <select v-model="diskSrs[i][j]" :aria-label="`Storage for ${v.name} - ${d.name}`" class="w-1/2 border border-line bg-white px-1 whitespace-pre">
                <option v-for="s in opts[hosts[i]]?.srs ?? []" :key="s.uuid" :value="s.uuid">{{ s.text }}</option>
              </select>
            </div>
          </template>
        </div>
        <template v-if="info.interfaces > 0">
          <span class="block">Virtual network interfaces in imported VMs:</span>
          <div class="max-h-28 overflow-auto border border-line">
            <template v-for="(v, i) in vms" :key="v.sys_id + i">
              <div v-for="(mac, j) in nets[i] ? v.macs ?? [] : []" :key="j" class="flex items-center gap-2 px-1 py-0.5">
                <span class="w-1/2 truncate">{{ v.name }} - Interface {{ j }}{{ mac ? ` (${mac})` : "" }}</span>
                <select v-model="nets[i][j]" :aria-label="`Network for ${v.name} interface ${j}`" class="w-1/2 border border-line bg-white px-1">
                  <option v-for="n in opts[hosts[i]]?.networks ?? []" :key="n.uuid" :value="n.uuid">{{ n.name }}</option>
                </select>
              </div>
            </template>
          </div>
        </template>
        <label class="flex items-center gap-2" :class="{ 'text-neutral-400': !info.has_manifest }">
          <input v-model="verify" type="checkbox" :disabled="!info.has_manifest" />Verify manifest content</label>
        <label class="flex items-center gap-2"><input v-model="start" type="checkbox" />Start VM(s) after import</label>
        <p v-if="error" class="text-red-600">{{ error }}</p>
        <div class="flex justify-end gap-2 pt-1">
          <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5">Import</button>
          <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Cancel</button>
        </div>
      </fieldset>
    </form>
  </div>
</template>
