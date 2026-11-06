<script setup lang="ts">
// Export (ExportApplianceWizard in one dialog): the name, location and format (ExportAppliancePage: an OVF/OVA package,
// or for one VM an XVA file), the VMs (ExportSelectVMsPage: the halted or suspended VMs of the pool, the selected VM or
// the vApp's VMs ticked), and for OVF/OVA the EULA files (ExportEulaPage), the options (ExportOptionsPage: manifest,
// OVA, compression, signing with a .pfx certificate and its private key's password) and Verify (ExportFinishPage); then
// ExportApplianceAction or ExportVmAction. The password stays in this dialog and the export's call; it is not kept.
// shortcut: no encryption (hidden in the C# wizard too); the certificate is checked on Export, not by a Validate button.
import { computed, onMounted, ref, watch } from "vue";
import { api, type ExportChoices } from "./api/workspace";

const props = defineProps<{ server: string; vm: string; vapp: string; name: string }>();
const emit = defineEmits<{ (e: "close"): void; (e: "run", label: string, f: () => Promise<string>): void }>();

const choices = ref<ExportChoices | null>(null);
const name = ref(props.name);
const dir = ref("");
const xva = ref(false);
const ticked = ref<string[]>([]);
const eulas = ref<string[]>([]);
const eula = ref("");
const manifest = ref(false);
const ova = ref(false);
const compress = ref(false);
const verify = ref(false);
const sign = ref(false);
const certFile = ref("");
const certPassword = ref("");
const error = ref("");
const pending = ref(false);
// ExportOptionsPage: signing needs the manifest, no manifest no signing
watch(sign, (v) => { if (v) manifest.value = true; });
watch(manifest, (v) => { if (!v) sign.value = false; });

onMounted(async () => {
  try {
    dir.value = await api.exportFolder();
    choices.value = await api.exportChoices(props.server, props.vm, props.vapp);
    ticked.value = choices.value.vms.filter((v) => v.selected).map((v) => v.uuid);
  } catch (e) { error.value = String(e); }
});
const vms = computed(() => choices.value?.vms ?? []);
const counter = computed(() => (ticked.value.length === 1 ? "1 VM selected" : `${ticked.value.length} VMs selected`));
const vtpm = computed(() => !xva.value && vms.value.some((v) => v.vtpm && ticked.value.includes(v.uuid)));

async function browse() {
  try { const d = await api.pickExportFolder(dir.value); if (d) dir.value = d; } catch (e) { error.value = String(e); }
}
async function addEula() {
  try {
    const f = await api.pickEulaFile();
    if (f && !eulas.value.includes(f) && eulas.value.length < 25) eulas.value.push(f);
  } catch (e) { error.value = String(e); }
}
async function browseCert() {
  try { const f = await api.pickCertificateFile(); if (f) { certFile.value = f; sign.value = true; } } catch (e) { error.value = String(e); }
}

async function ok() {
  if (pending.value || !choices.value) return;
  pending.value = true;
  const n = name.value.trim(), d = dir.value.trim(), picked = vms.value.filter((v) => ticked.value.includes(v.uuid));
  error.value = "";
  if (!picked.length) error.value = "Select the VM(s) you want to export.";
  else if (xva.value && picked.length > 1)
    error.value = "Only one VM may be selected for export in XVA format. To export multiple VMs, select OVF/OVA Package as the export format.";
  else {
    try {
      error.value = await api.applianceTargetError(d, n, xva.value);
      if (!error.value && !xva.value && sign.value) error.value = await api.signingCertificateError(certFile.value, certPassword.value);
    } catch (e) { error.value = String(e); }
  }
  if (error.value) { pending.value = false; return; }
  if (xva.value) {
    const file = `${d}/${n.endsWith(".xva") ? n : n + ".xva"}`, check = verify.value;
    emit("run", `Exporting '${picked[0].name}'`, async () => { await api.exportVm(props.server, picked[0].uuid, file, check); return file; });
  } else {
    const e = { dir: d, name: n, vms: picked.map((v) => v.uuid), eulas: [...eulas.value], manifest: manifest.value, ova: ova.value, compress: compress.value, verify: verify.value,
      sign: sign.value, certificate: sign.value ? certFile.value : "", password: sign.value ? certPassword.value : "" };
    emit("run", `Exporting '${n}'`, () => api.exportAppliance(props.server, e));
  }
  emit("close");
}
</script>

<template>
  <div class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
    <form novalidate class="flex max-h-[95vh] w-[44rem] flex-col border border-neutral-500 bg-white shadow-lg" @submit.prevent="ok">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ xva ? "Export to File" : "Export OVF/OVA Package" }}</div>
      <fieldset :disabled="pending || !choices" class="min-h-0 space-y-2 overflow-auto p-3">
        <p class="whitespace-pre-line">vApps and multiple VMs can be exported as an OVF/OVA package.
A single VM can be exported as an OVF/OVA package or as an XVA file.</p>
        <div class="grid grid-cols-[6rem_1fr_auto] items-center gap-x-2 gap-y-1">
          <label for="expname">Name:</label><input id="expname" v-model="name" class="col-span-2 border border-line px-1 py-0.5" />
          <label for="expdir">Location:</label><input id="expdir" v-model="dir" class="border border-line px-1 py-0.5" />
          <button type="button" class="btn" @click="browse">Browse...</button>
          <label for="expformat">Format:</label>
          <select id="expformat" v-model="xva" class="col-span-2 border border-line bg-white px-1 py-0.5">
            <option :value="false">OVF/OVA Package (*.ovf, *.ova)</option>
            <option v-if="!vapp" :value="true">XVA File (*.xva)</option>
          </select>
        </div>
        <p>{{ choices?.is_pool ? `VMs in pool '${choices?.pool}':` : `VMs in server '${choices?.pool ?? ""}':` }}</p>
        <div class="h-36 overflow-auto border border-line">
          <table class="w-full border-collapse">
            <thead class="sticky top-0 bg-pane text-left"><tr><th class="w-6"></th><th class="px-1 font-normal">Name</th><th class="px-1 font-normal">Description</th>
              <th class="px-1 font-normal">Size</th><th class="px-1 font-normal">vApp</th></tr></thead>
            <tbody><tr v-for="v in vms" :key="v.uuid" class="hover:bg-neutral-100">
              <td class="text-center"><input v-model="ticked" type="checkbox" :value="v.uuid" :aria-label="v.name" /></td>
              <td class="px-1">{{ v.name }}</td><td class="px-1">{{ v.description }}</td><td class="px-1 whitespace-nowrap">{{ v.size }}</td>
              <td class="px-1">{{ v.vapp || "None" }}</td></tr></tbody>
          </table>
        </div>
        <div class="flex items-center gap-2">
          <span class="flex-1">{{ counter }}</span>
          <button type="button" class="btn" :disabled="ticked.length === vms.length" @click="ticked = vms.map((v) => v.uuid)">Select All</button>
          <button type="button" class="btn" :disabled="!ticked.length" @click="ticked = []">Clear All</button>
        </div>
        <p v-if="vtpm" class="text-yellow-700">The OVF/OVA format does not preserve vTPM data. Some of the VMs' security features may be deactivated and the VMs may not start. Use the XVA format to preserve vTPM data.</p>
        <template v-if="!xva">
          <span class="block font-semibold">EULAs</span>
          <p>You can include end user license agreement (EULA) documents in the exported OVF/OVA package.</p>
          <div class="flex gap-2">
            <select v-model="eula" size="3" aria-label="EULA files" class="min-w-0 flex-1 border border-line">
              <option v-for="f in eulas" :key="f" :value="f">{{ f }}</option>
            </select>
            <div class="flex flex-col gap-1">
              <button type="button" class="btn" :disabled="eulas.length >= 25" @click="addEula">Add...</button>
              <button type="button" class="btn" :disabled="!eula" @click="eulas = eulas.filter((f) => f !== eula); eula = ''">Remove</button>
            </div>
          </div>
          <span class="block font-semibold">Advanced Options</span>
          <label class="flex items-center gap-2"><input v-model="manifest" type="checkbox" />Create a manifest</label>
          <label class="flex items-center gap-2"><input v-model="sign" type="checkbox" />Sign the OVF package</label>
          <div class="grid grid-cols-[10rem_1fr_auto] items-center gap-x-2 gap-y-1 pl-6">
            <label for="expcert">Certificate path:</label>
            <input id="expcert" v-model="certFile" class="border border-line px-1 py-0.5" @input="sign = true" />
            <button type="button" class="btn" @click="browseCert">Browse...</button>
            <label for="expcertpw">Private key password:</label>
            <input id="expcertpw" v-model="certPassword" type="password" autocomplete="off" class="col-span-2 border border-line px-1 py-0.5" @input="sign = true" />
          </div>
          <label class="flex items-center gap-2"><input v-model="ova" type="checkbox" />Create OVA package (single .ova file)</label>
          <label class="flex items-center gap-2"><input v-model="compress" type="checkbox" />Compress package</label>
          <p v-if="compress" class="pl-6 text-neutral-600">{{ ova ? "The OVA package will be compressed to an .ova.gz file." : "Each of the disks of the OVF package will be compressed to a .vhd.gz file." }}</p>
        </template>
        <label class="flex items-center gap-2"><input v-model="verify" type="checkbox" />Verify export on completion</label>
        <p v-if="verify && xva" class="pl-6 text-yellow-700">Verify export can take up to 5 hours per TB.</p>
        <p v-if="error" class="text-red-600">{{ error }}</p>
      </fieldset>
      <div class="flex justify-end gap-2 border-t border-line p-2">
        <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="pending || !choices">Export</button>
        <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Cancel</button>
      </div>
    </form>
  </div>
</template>
