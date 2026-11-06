<script setup lang="ts">
// The Storage, Networking and Memory tabs of a VM (VmStorageView, VmNetworkView, MemoryView): the tables from the
// Go side, plus the DVD drive and Activate/Deactivate of a disk on Storage.
import { computed, onMounted, ref, watch } from "vue";
import { api, type Confirm, type NewDisk, type Table, type VIF, type VIFSettings, type VMNetworking, type VMStorage } from "./api/workspace";
import DetailsTable from "./DetailsTable.vue";
import MemoryDialog from "./MemoryDialog.vue";
import NewDiskDialog from "./NewDiskDialog.vue";
import VifDialog from "./VifDialog.vue";
import { serial } from "./serial";

// tick: the server's objects changed (xcp:changed), read the tab again
const props = defineProps<{ server: string; vm: string; kind: "storage" | "networking" | "memory"; tick?: number }>();
const emit = defineEmits<{ (e: "status", text: string, error: boolean): void }>();

const table = ref<Table | null>(null);
const storage = ref<VMStorage | null>(null);
const picked = ref("");
const busy = ref(false);
const disk = computed(() => storage.value?.disks.find((d) => d.vbd === picked.value) ?? null);
const network = ref<VMNetworking | null>(null);
const vif = computed(() => network.value?.vifs.find((v) => v.uuid === picked.value) ?? null);
let editing: { device: string; uuid: string; done?: boolean } | null = null; // an edit in flight (UpdateVIFAction deletes and creates the VIF)

// one read at a time (serial): changes coming faster than a slow read make one more read, never a stale table
const load = serial(async () => {
  try {
    if (props.kind === "storage") { storage.value = await api.storage(props.server, props.vm); table.value = storage.value.table; }
    else if (props.kind === "networking") {
      network.value = await api.networking(props.server, props.vm);
      table.value = network.value.table;
      // NetworkRules.Selection / EditReplaced: while an edit recreates a VIF, its device stays selected (none while it
      // is missing) until a new VIF holds it; otherwise the selected VIF, else the first
      const vifs = network.value.vifs;
      if (editing) {
        const fresh = vifs.find((v) => v.device === editing!.device && v.uuid !== editing!.uuid);
        if (fresh) { picked.value = fresh.uuid; editing = null; }
        else {
          picked.value = vifs.find((v) => v.uuid === editing!.uuid)?.uuid ?? "";
          if (editing.done && picked.value) editing = null; // an edit that changed nothing kept its VIF
        }
      } else if (!table.value.keys.includes(picked.value)) picked.value = table.value.keys[0] ?? "";
    }
    else table.value = await api.memory(props.server, props.vm);
  } catch (e) { emit("status", String(e), true); }
});

async function run(label: string, f: () => Promise<void>) {
  busy.value = true;
  emit("status", `${label}...`, false);
  try { await f(); emit("status", "", false); } catch (e) { emit("status", `${label}: ${String(e)}`, true); }
  await load(); // the buttons stay disabled until the tab shows the new state
  busy.value = false;
}

function toggle() {
  const d = disk.value;
  if (d) run(d.attached ? "Deactivating the disk" : "Activating the disk", () => api.setDiskActive(props.server, d.vbd, !d.attached));
}

function changeIso(e: Event) {
  const iso = (e.target as HTMLSelectElement).value, st = storage.value;
  if (st?.dvd && iso !== st.iso) run(iso ? "Inserting the ISO" : "Ejecting the ISO", () => api.changeIso(props.server, st.dvd!, iso));
}

const notice = ref("");
const editMemory = ref(false); // the Memory tab's Edit... (VMMemoryControlsNoEdit)

// the Networking tab's buttons (NetworkList): Add Interface, Properties, Remove Interface, Activate/Deactivate
const vifEdit = ref<{ vif: VIF | null } | null>(null);
function addVif() {
  if (network.value?.at_max) notice.value = network.value.at_max;
  else vifEdit.value = { vif: null };
}
function saveVif(set: VIFSettings) {
  const v = vifEdit.value?.vif ?? null;
  vifEdit.value = null;
  if (v) editing = { device: v.device, uuid: v.uuid };
  run(v ? "Updating the virtual network interface" : "Creating a virtual network interface", async () => {
    try { notice.value = await api.saveVif(props.server, props.vm, v?.uuid ?? "", set); if (editing) editing.done = true; }
    catch (e) { editing = null; throw e; } // a failed edit: the selection follows the tab again
  });
}
async function removeVif() {
  const v = vif.value;
  if (!v) return;
  const q = await api.vifDeleteQuestion();
  if (await new Promise<boolean>((resolve) => (ask.value = { ...q, resolve }))) run("Deleting the virtual network interface", () => api.deleteVif(props.server, v.uuid));
}
function toggleVif() {
  const v = vif.value;
  if (v) run(v.attached ? "Unplugging the virtual network interface" : "Plugging the virtual network interface", () => api.setVifActive(props.server, v.uuid, !v.attached));
}

// Add Virtual Disk (NewDiskDialog): name, description, size in GB or MB, the SR
const add = ref<NewDisk | null>(null);
async function openAdd() {
  try { add.value = await api.newDiskOptions(props.server, props.vm); } catch (e) { emit("status", String(e), true); }
}
function submitAdd(name: string, description: string, sr: string, bytes: number) {
  add.value = null;
  run(`Adding virtual disk '${name}'`, async () => { notice.value = await api.addDisk(props.server, props.vm, name, description, sr, bytes); });
}

// Attach Disk (AttachDiskDialog): an unattached disk, optionally read-only
const attach = ref<{ disks: { vdi: string; text: string }[]; vdi: string; readOnly: boolean } | null>(null);
async function openAttach() {
  try { attach.value = { disks: (await api.attachableDisks(props.server, props.vm)) ?? [], vdi: "", readOnly: false }; } catch (e) { emit("status", String(e), true); }
}
function submitAttach() {
  const a = attach.value!;
  attach.value = null;
  run("Attaching the disk", async () => { notice.value = await api.attachDisk(props.server, props.vm, a.vdi, a.readOnly); });
}

// Detach and Delete ask the C# questions first
const ask = ref<(Confirm & { resolve: (ok: boolean) => void }) | null>(null);
async function confirm(del: boolean, system: boolean) {
  const q = await api.diskQuestion(del, system);
  return new Promise<boolean>((resolve) => (ask.value = { ...q, resolve }));
}
function answer(ok: boolean) { ask.value?.resolve(ok); ask.value = null; }
async function detach() {
  const d = disk.value;
  if (d && (await confirm(false, d.system))) run(`Detaching '${d.name}'`, () => api.detachDisk(props.server, d.vbd));
}
async function del() {
  const d = disk.value;
  if (d && (await confirm(true, d.system))) run(`Deleting '${d.name}'`, () => api.deleteDisk(props.server, props.vm, d.vbd));
}
function createDvd() {
  run("Creating a DVD drive", async () => { notice.value = await api.createDvdDrive(props.server, props.vm); });
}

onMounted(load);
watch(() => props.tick, load);
</script>

<template>
  <div class="flex min-h-0 flex-1 flex-col overflow-auto">
    <div v-if="kind === 'storage' && storage?.dvd" class="flex items-center gap-2 px-3 pt-3">
      <span>{{ storage.dvd_name }}:</span>
      <select class="w-96 border border-line bg-white px-1 py-0.5" :value="storage.iso" :disabled="busy" @change="changeIso">
        <option v-for="i in storage.isos" :key="i.vdi" :value="i.vdi">{{ i.text }}</option>
      </select>
      <button class="btn" :disabled="busy || !storage.iso" @click="run('Ejecting the ISO', () => api.changeIso(server, storage!.dvd!, ''))">Eject</button>
    </div>
    <p v-if="kind === 'storage' && storage?.can_add_dvd" class="pt-3 text-center">
      <a class="cursor-pointer font-semibold underline" @click="!busy && createDvd()">Click here to create a DVD drive</a>
    </p>
    <p v-if="notice" class="mx-3 mt-2 border border-amber-300 bg-amber-50 px-2 py-1">{{ notice }}</p>
    <DetailsTable v-if="table" :table="table" :selected="picked" :plain="kind === 'memory'" @select="picked = $event" />
    <div v-if="kind === 'networking' && network" class="flex gap-2 px-3">
      <button class="btn" :disabled="busy || !network.can_add" @click="addVif">Add Interface...</button>
      <button class="btn" :disabled="busy || !vif?.can_edit" :title="vif?.edit_tip" @click="vifEdit = { vif }">Properties</button>
      <button class="btn" :disabled="busy || !vif?.can_edit" :title="vif?.remove_tip" @click="removeVif">Remove Interface</button>
      <button class="btn" :disabled="busy || !vif?.can_toggle" :title="vif?.toggle_tip" @click="toggleVif">{{ vif?.toggle_label ?? "Deactivate" }}</button>
    </div>
    <VifDialog v-if="vifEdit" :server="server" :vm="vm" :vif="vifEdit.vif" @close="vifEdit = null" @save="saveVif" />
    <div v-if="kind === 'memory'" class="flex gap-2 px-3">
      <button class="btn" :disabled="busy" @click="editMemory = true">Edit...</button>
    </div>
    <MemoryDialog v-if="editMemory" :server="server" :vm="vm" @close="editMemory = false" @run="run" />
    <div v-if="kind === 'storage'" class="flex gap-2 px-3">
      <button class="btn" :disabled="busy" @click="openAdd">Add...</button>
      <button class="btn" :disabled="busy" @click="openAttach">Attach Disk...</button>
      <button class="btn" :disabled="busy || !disk || !(disk.attached ? disk.can_deactivate : disk.can_activate)" @click="toggle">
        {{ disk && !disk.attached ? "Activate" : "Deactivate" }}
      </button>
      <button class="btn" :disabled="busy || !disk?.can_detach" @click="detach">Detach</button>
      <button class="btn" :disabled="busy || !disk?.can_delete" @click="del">Delete</button>
    </div>

    <NewDiskDialog v-if="add" :options="add" @close="add = null" @add="submitAdd" />

    <div v-if="attach" class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
      <div class="w-[34rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Attach Disk</div>
        <div class="space-y-2 p-3">
          <p>{{ attach.disks.length ? "Select the virtual disk you want to attach to this VM." : "There are no unattached virtual disks available to this VM." }}</p>
          <div class="h-56 overflow-auto border border-line">
            <div v-for="d in attach.disks" :key="d.vdi" class="cursor-default px-1 py-0.5" :class="attach.vdi === d.vdi ? 'bg-select text-white' : 'hover:bg-neutral-100'"
                 @click="attach!.vdi = d.vdi">{{ d.text }}</div>
          </div>
          <label class="flex items-center gap-2"><input v-model="attach.readOnly" type="checkbox" />Attach as read-only</label>
          <div class="flex justify-end gap-2 pt-1">
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="!attach.vdi" @click="submitAttach">Attach</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="attach = null">Cancel</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="ask" class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
      <div class="w-[32rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ ask.title }}</div>
        <div class="space-y-3 p-3">
          <p class="whitespace-pre-line">{{ ask.text }}</p>
          <div class="flex justify-end gap-2">
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="answer(true)">Yes</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="answer(false)">No</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
