<script setup lang="ts">
// The Memory, Storage, Networking and NICs tabs of a server or pool, and the Storage tab of an SR (Details' tables);
// the Networking tab's Add Network... / Properties / Remove Network (HostNetworkView), the NICs tab's Create Bond... /
// Delete Bond / Rescan (NicsView), an SR's Storage tab's Rescan / Add... / Move... / Delete / Properties (SrStorageView)
// and a server's or pool's Storage tab's New SR... / Reclaim freed space / Properties (PhysicalStorageView).
import { computed, onMounted, ref, watch } from "vue";
import { api, type HostNetworking, type HostNics, type NewDisk, type PhysicalStorage, type SRStorage, type Table } from "./api/workspace";
import DetailsTable from "./DetailsTable.vue";
import MemoryBars from "./MemoryBars.vue";
import MoveDiskDialog from "./MoveDiskDialog.vue";
import NetworkPropertiesDialog from "./NetworkPropertiesDialog.vue";
import NewDiskDialog from "./NewDiskDialog.vue";
import NewSrWizard from "./NewSrWizard.vue";
import ObjectPropertiesDialog from "./ObjectPropertiesDialog.vue";
import NewNetworkWizard from "./NewNetworkWizard.vue";
import { serial } from "./serial";

// tick: the server's objects changed (xcp:changed), read the tab again
const props = defineProps<{ server: string; kind: string; uuid: string; tab: string; tick?: number }>();
const emit = defineEmits<{ (e: "status", text: string, error: boolean): void }>();

const table = ref<Table | null>(null);
const nets = ref<HostNetworking | null>(null);
const nics = ref<HostNics | null>(null);
const srStorage = ref<SRStorage | null>(null);
const physical = ref<PhysicalStorage | null>(null);
const picked = ref("");
const busy = ref(false);
const load = serial(async () => {
  try {
    if (props.tab === "networking") table.value = (nets.value = await api.hostNetworking(props.server, props.kind, props.uuid)).table;
    else if (props.tab === "nics") table.value = (nics.value = await api.hostNics(props.server, props.uuid)).table;
    else if (props.tab === "storage" && props.kind === "srs") table.value = (srStorage.value = await api.srStorage(props.server, props.uuid)).table;
    else if (props.tab === "storage") table.value = (physical.value = await api.physicalStorage(props.server, props.kind, props.uuid)).table;
    else table.value = await api.hostTab(props.server, props.kind, props.uuid, props.tab);
    if (props.tab !== "memory" && !table.value.keys.includes(picked.value)) picked.value = table.value.keys[0] ?? "";
  } catch (e) { emit("status", String(e), true); }
});
onMounted(load);
watch(() => props.tick, load);

async function run(label: string, f: () => Promise<void>) {
  wizard.value = props_.value = newDisk.value = moveOf.value = null;
  busy.value = true;
  emit("status", `${label}...`, false);
  try { await f(); emit("status", "", false); } catch (e) { emit("status", `${label}: ${String(e)}`, true); }
  await load(); // the buttons stay disabled until the tab shows the new state
  busy.value = false;
}

const network = computed(() => nets.value?.networks.find((n) => n.uuid === picked.value) ?? null);
const nic = computed(() => nics.value?.nics.find((n) => n.uuid === picked.value) ?? null);
const disk = computed(() => srStorage.value?.disks.find((d) => d.uuid === picked.value) ?? null);
const sr = computed(() => physical.value?.srs.find((s) => s.uuid === picked.value) ?? null);
const wizard = ref<{ bond: boolean } | null>(null);
const props_ = ref<string | null>(null); // the network whose Properties are open
const more = ref<{ kind: string; uuid: string } | null>(null); // the object whose folder, tags and custom fields are open

// the questions: ok "" is an error to read (a blocked bond change)
const ask = ref<{ title: string; text: string; ok: string; cancel: string; resolve: (ok: boolean) => void } | null>(null);
const question = (title: string, text: string, ok: string, cancel: string) =>
  new Promise<boolean>((resolve) => (ask.value = { title, text, ok, cancel, resolve }));
function answer(ok: boolean) { ask.value?.resolve(ok); ask.value = null; }

// DestroyBondCommand: the question (or the error under HA), then DestroyBondAction
async function deleteBond(networkUUID: string) {
  try {
    const w = await api.deleteBondQuestion(props.server, networkUUID);
    if (w.block) { await question(w.title, w.text, "", "OK"); return; }
    if (await question(w.title, w.text, w.ok, "Cancel")) run("Deleting bond", () => api.deleteBond(props.server, networkUUID));
  } catch (e) { emit("status", String(e), true); }
}
async function removeNetwork() {
  const n = network.value;
  if (!n?.can_remove) return;
  if (n.bond) return deleteBond(n.uuid);
  const q = await api.removeNetworkQuestion();
  if (await question(q.title, q.text, "Yes", "No")) run("Removing network", () => api.removeNetwork(props.server, n.uuid));
}
function edit() { if (network.value?.can_edit) props_.value = network.value.uuid; }

// an SR's Storage tab (SrStorageView): Add makes a disk attached to no VM, Move migrates or moves it, Delete asks the
// DeleteConfirmation question; the new or moved disk stays selected
const newDisk = ref<NewDisk | null>(null);
const moveOf = ref<string | null>(null);
async function openAdd() {
  try { newDisk.value = await api.newSrDiskOptions(props.server, props.uuid); } catch (e) { emit("status", String(e), true); }
}
function addDisk(name: string, description: string, srUUID: string, bytes: number) {
  run(`Adding virtual disk '${name}'`, async () => { picked.value = await api.createSrDisk(props.server, srUUID, name, description, bytes); });
}
function moveDisk(target: string, migrate: boolean) {
  const vdi = moveOf.value!;
  run(migrate ? "Migrating virtual disk" : "Moving virtual disk", async () => { picked.value = await api.moveSrDisk(props.server, vdi, target); });
}
async function deleteDisk() {
  const d = disk.value;
  if (d?.can_delete && (await question(d.delete_title, d.delete_text, "Yes", "No"))) run("Deleting virtual disk", () => api.deleteSrDisk(props.server, d.uuid));
}
function diskProperties() { if (disk.value?.can_edit) more.value = { kind: "vdis", uuid: disk.value.uuid }; }

// a server's or pool's Storage tab (PhysicalStorageView): New SR (NewSRCommand), Reclaim freed space (TrimSRCommand)
const newSr = ref(false);
async function trim() {
  const s = sr.value;
  if (!s?.can_trim) return;
  const q = await api.trimQuestion();
  if (await question(q.title, q.text, "Yes", "No")) run("Reclaiming freed space", () => api.trimSr(props.server, s.uuid));
}
function srProperties() { if (sr.value) more.value = { kind: "srs", uuid: sr.value.uuid }; }
function newSrClosed(created: string) {
  newSr.value = false;
  if (created) picked.value = created;
  load();
}
// NetworkList, SrStoragePage and PhysicalStoragePage: a double click is Properties
function activate() {
  if (busy.value) return;
  if (props.tab === "networking") edit();
  else if (props.tab === "storage" && props.kind === "srs") diskProperties();
  else if (props.tab === "storage") srProperties();
}
</script>

<template>
  <div class="flex min-h-0 flex-1 flex-col overflow-auto">
    <MemoryBars v-if="tab === 'memory'" :server="server" :kind="kind" :uuid="uuid" :tick="tick" />
    <div v-if="table" @dblclick="activate">
      <DetailsTable :table="table" :selected="picked" :plain="tab === 'memory'" @select="picked = $event" />
    </div>
    <div v-if="tab === 'networking' && nets" class="flex gap-2 px-3">
      <button class="btn" :disabled="busy || !nets.can_add" @click="wizard = { bond: false }">Add Network...</button>
      <button class="btn" :disabled="busy || !network?.can_edit" @click="edit">Properties</button>
      <button class="btn" :disabled="busy || !network?.can_remove" :title="network?.remove_tip" @click="removeNetwork">Remove Network</button>
    </div>
    <div v-if="tab === 'nics' && nics" class="flex gap-2 px-3">
      <button class="btn" :disabled="busy || !nics.can_bond" @click="wizard = { bond: true }">Create Bond...</button>
      <button class="btn" :disabled="busy || !nic?.bond_network" @click="nic?.bond_network && deleteBond(nic.bond_network)">Delete Bond</button>
      <button class="btn" :disabled="busy || !nics.can_bond" @click="run('Scanning for NICs', () => api.rescanNics(server, uuid))">Rescan</button>
    </div>
    <div v-if="tab === 'storage' && kind === 'srs' && srStorage" class="flex gap-2 px-3">
      <button class="btn" :disabled="busy || !srStorage.can_rescan" :title="srStorage.rescan_tip" @click="run('Rescanning', () => api.rescanSrs(server, [uuid]))">Rescan</button>
      <button class="btn" :disabled="busy || !srStorage.can_add" @click="openAdd">Add...</button>
      <button class="btn" :disabled="busy || !disk?.can_move" @click="moveOf = disk!.uuid">Move...</button>
      <button class="btn" :disabled="busy || !disk?.can_delete" @click="deleteDisk">Delete</button>
      <button class="btn" :disabled="busy || !disk?.can_edit" @click="diskProperties">Properties</button>
    </div>
    <div v-if="tab === 'storage' && kind !== 'srs' && physical" class="flex gap-2 px-3">
      <button class="btn" :disabled="busy || !physical.can_new" @click="newSr = true">New SR...</button>
      <button class="btn" :disabled="busy || !sr?.can_trim" :title="sr?.trim_tip" @click="trim">Reclaim freed space</button>
      <button class="btn" :disabled="busy || !sr" @click="srProperties">Properties</button>
    </div>
    <NewNetworkWizard v-if="wizard" :server="server" :kind="kind" :uuid="uuid" :bond="wizard.bond" @close="wizard = null" @run="run" />
    <NetworkPropertiesDialog v-if="props_" :server="server" :uuid="props_" @close="props_ = null" @run="run" @more="more = { kind: 'networks', uuid: props_! }; props_ = null" />
    <ObjectPropertiesDialog v-if="more" :server="server" :kind="more.kind" :uuid="more.uuid" @close="more = null; load()" />
    <NewDiskDialog v-if="newDisk" :options="newDisk" @close="newDisk = null" @add="addDisk" />
    <MoveDiskDialog v-if="moveOf" :server="server" :vdi="moveOf" @close="moveOf = null" @move="moveDisk" />
    <NewSrWizard v-if="newSr && physical" :server="server" :pool="physical.pool" @close="newSrClosed" />
    <div v-if="ask" class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
      <div class="w-[32rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ ask.title }}</div>
        <div class="space-y-3 p-3">
          <p class="whitespace-pre-line">{{ ask.text }}</p>
          <div class="flex justify-end gap-2">
            <button v-if="ask.ok" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="answer(true)">{{ ask.ok }}</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" autofocus @click="answer(false)">{{ ask.cancel }}</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
