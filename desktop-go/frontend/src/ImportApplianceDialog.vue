<script setup lang="ts">
// Import Appliance (XenAdmin.Desktop's Dialogs.ImportAppliance: the import wizard's OVF/OVA path in one dialog): the
// package's VMs, its EULA to accept, the pool or server, the SR for every disk, the network for every interface,
// manifest verification and Start, then ImportApplianceAction.
// shortcut: one SR and one network for the whole appliance, as OvfRules (the wizard has a choice per disk and interface).
import { computed, onMounted, ref, watch } from "vue";
import { api, type Appliance, type ServerRow } from "./api/workspace";

const props = defineProps<{ file: string }>();
const emit = defineEmits<{ (e: "close"): void; (e: "run", label: string, key: string, f: () => Promise<string>): void }>();

type Target = { key: string; host: string; text: string };
const info = ref<Appliance | null>(null);
const loadError = ref("");
const targets = ref<Target[]>([]);
const target = ref(0);
const srs = ref<{ uuid: string; text: string }[]>([]);
const sr = ref("");
const networks = ref<{ uuid: string; name: string }[]>([]);
const network = ref("");
const accept = ref(false);
const verify = ref(false);
const start = ref(false);
const error = ref("");
const pending = ref(false);

const eulas = computed(() => info.value?.eulas ?? []);
const summary = computed(() => {
  const a = info.value;
  if (!a) return "";
  const vms = a.vms ?? [];
  return `Appliance '${a.name}': ${vms.length} VM(s) (${vms.map((v) => v.name).join(", ")}), ${a.disks} disk(s), ${a.interfaces} network interface(s).`;
});

onMounted(async () => {
  try { info.value = await api.applianceInfo(props.file); } catch (e) { loadError.value = String(e); return; } // ImportSourcePage: read once
  try {
    const servers: ServerRow[] = (await api.servers()).filter((s) => s.connected);
    const out: Target[] = [];
    for (const s of servers) {
      const o = await api.importOptions(s.key, "");
      o.servers.forEach((h, i) => out.push({ key: s.key, host: h.uuid, text: i === 0 ? h.name : "    " + h.name }));
    }
    targets.value = out;
  } catch (e) { error.value = String(e); }
});
// the SRs and networks of the chosen target, the pool's default SR and the first network selected
let asked = 0; // a slow answer for a target no longer chosen is dropped
watch([targets, target], async () => {
  const t = targets.value[target.value], id = ++asked;
  srs.value = []; sr.value = ""; networks.value = []; network.value = "";
  if (!t) return;
  try {
    const o = await api.importOptions(t.key, t.host);
    if (id !== asked) return;
    srs.value = o.srs ?? [];
    sr.value = o.default_sr && srs.value.some((x) => x.uuid === o.default_sr) ? o.default_sr : srs.value[0]?.uuid ?? "";
    networks.value = o.networks ?? [];
    network.value = networks.value[0]?.uuid ?? "";
  } catch (e) { error.value = String(e); }
});

async function ok() {
  const a = info.value;
  if (pending.value || !a) return;
  pending.value = true; // the form is frozen: the values checked are the values imported
  const t = targets.value[target.value], srUuid = sr.value, net = network.value, check = verify.value && a.has_manifest, run = start.value;
  error.value = "";
  if (eulas.value.length > 0 && !accept.value) error.value = "You must accept the EULA(s) to import this appliance.";
  else if (!t) error.value = "Select the pool or server to import to.";
  else if (!srUuid) error.value = "Select a storage repository.";
  else if (a.interfaces > 0 && !net) error.value = "Select a network.";
  if (!error.value) {
    try { error.value = await api.applianceSpaceError(t.key, srUuid, a.required_space); } catch (e) { error.value = String(e); }
  }
  if (error.value) { pending.value = false; return; }
  emit("run", `Importing '${a.name}'`, t.key, () => api.importAppliance(t.key, props.file, srUuid, net, check, run));
  emit("close");
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <form novalidate class="w-[38rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="ok">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">Import Appliance</div>
      <div v-if="loadError" class="space-y-2 p-3">
        <p class="whitespace-pre-wrap">{{ loadError }}</p>
        <div class="flex justify-end pt-1">
          <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">OK</button>
        </div>
      </div>
      <p v-else-if="!info" class="p-3">Reading '{{ file.split("/").pop() }}'...</p>
      <fieldset v-else :disabled="pending" class="space-y-2 p-3">
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
        <span class="block">Storage for the virtual disks:</span>
        <div class="h-28 overflow-auto border border-line">
          <div v-for="s in srs" :key="s.uuid" class="cursor-default whitespace-pre px-1 py-0.5" :class="sr === s.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100'"
               @click="sr = s.uuid">{{ s.text }}</div>
        </div>
        <label for="appnet" class="block">Network for the network interfaces:</label>
        <select id="appnet" v-model="network" class="w-full border border-line bg-white px-1 py-0.5">
          <option v-for="n in networks" :key="n.uuid" :value="n.uuid">{{ n.name }}</option>
        </select>
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
