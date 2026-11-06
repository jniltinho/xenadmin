<script setup lang="ts">
// Import VM (the import wizard's XVA path in one dialog, as XenAdmin.Desktop's Dialogs.ImportVm): the file, the pool or
// server to import to (a pool's servers as the home server), the SR and "Start VM after import". An OVF/OVA appliance
// goes to its own dialog (ImportApplianceDialog), as the C# app's Dialogs.ImportAppliance.
import { onMounted, ref, watch } from "vue";
import { api, type ServerRow } from "./api/workspace";

const emit = defineEmits<{ (e: "close"): void; (e: "run", label: string, key: string, f: () => Promise<string>): void; (e: "appliance", file: string): void }>();

type Target = { key: string; host: string; text: string };
const file = ref("");
const targets = ref<Target[]>([]);
const target = ref(0);
const srs = ref<{ uuid: string; text: string }[]>([]);
const sr = ref("");
const start = ref(false);
const error = ref("");
const pending = ref(false);

onMounted(async () => {
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
// the SRs of the chosen target, the pool's default selected
let asked = 0; // a slow answer for a target no longer chosen is dropped
watch([targets, target], async () => {
  const t = targets.value[target.value], id = ++asked;
  srs.value = []; sr.value = "";
  if (!t) return;
  try {
    const o = await api.importOptions(t.key, t.host);
    if (id !== asked) return;
    srs.value = o.srs ?? [];
    sr.value = o.default_sr && srs.value.some((x) => x.uuid === o.default_sr) ? o.default_sr : srs.value[0]?.uuid ?? "";
  } catch (e) { error.value = String(e); }
});

async function browse() {
  try { const p = await api.pickImportFile(); if (p) file.value = p; } catch (e) { error.value = String(e); }
}
async function ok() {
  if (pending.value) return;
  pending.value = true; // the form is frozen: the values checked are the values imported
  const t = targets.value[target.value], path = file.value.trim(), srUuid = sr.value, run = start.value;
  try {
    if (await api.isAppliance(path)) { emit("appliance", path); emit("close"); return; } // the appliance path has its own pages
    error.value = await api.importFileError(path);
    if (!error.value && !t) error.value = "Select the pool or server to import to.";
    if (!error.value && !srUuid) error.value = "Select a storage repository.";
  } catch (e) { error.value = String(e); }
  if (error.value) { pending.value = false; return; }
  emit("run", `Importing '${path.split("/").pop()}'`, t.key, () => api.importVm(t.key, path, t.host, srUuid, run));
  emit("close");
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <form novalidate class="w-[36rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="ok">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">Import VM</div>
      <fieldset :disabled="pending" class="space-y-2 p-3">
        <p>Import a VM from an XVA file, or an appliance from an OVF/OVA package, to a pool or server.</p>
        <div class="flex gap-2">
          <input v-model="file" placeholder="File to import (.xva, .xva.gz, .ovf, .ova, .ova.gz, .ova.bz2)" aria-label="File to import" class="min-w-0 flex-1 border border-line px-1 py-0.5" />
          <button type="button" class="btn" @click="browse">Browse...</button>
        </div>
        <label for="imptarget" class="block">Import to:</label>
        <select id="imptarget" v-model.number="target" class="w-full border border-line bg-white px-1 py-0.5 whitespace-pre">
          <option v-for="(t, i) in targets" :key="t.key + t.host" :value="i">{{ t.text }}</option>
        </select>
        <span class="block">Storage:</span>
        <div class="h-36 overflow-auto border border-line">
          <div v-for="s in srs" :key="s.uuid" class="cursor-default whitespace-pre px-1 py-0.5" :class="sr === s.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100'"
               @click="sr = s.uuid">{{ s.text }}</div>
        </div>
        <label class="flex items-center gap-2"><input v-model="start" type="checkbox" />Start VM after import</label>
        <p v-if="error" class="text-red-600">{{ error }}</p>
        <div class="flex justify-end gap-2 pt-1">
          <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5">Import</button>
          <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Cancel</button>
        </div>
      </fieldset>
    </form>
  </div>
</template>
