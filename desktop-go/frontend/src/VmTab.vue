<script setup lang="ts">
// The Storage, Networking and Memory tabs of a VM (VmStorageView, VmNetworkView, MemoryView): the tables from the
// Go side, plus the DVD drive and Activate/Deactivate of a disk on Storage.
import { computed, onMounted, ref } from "vue";
import { api, type Table, type VMStorage } from "./api/workspace";
import DetailsTable from "./DetailsTable.vue";

const props = defineProps<{ server: string; vm: string; kind: "storage" | "networking" | "memory" }>();
const emit = defineEmits<{ (e: "status", text: string, error: boolean): void }>();

const table = ref<Table | null>(null);
const storage = ref<VMStorage | null>(null);
const picked = ref("");
const busy = ref(false);
const disk = computed(() => storage.value?.disks.find((d) => d.vbd === picked.value) ?? null);

async function load() {
  try {
    if (props.kind === "storage") { storage.value = await api.storage(props.server, props.vm); table.value = storage.value.table; }
    else table.value = await (props.kind === "networking" ? api.networking : api.memory)(props.server, props.vm);
  } catch (e) { emit("status", String(e), true); }
}

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
function createDvd() {
  run("Creating a DVD drive", async () => { notice.value = await api.createDvdDrive(props.server, props.vm); });
}

onMounted(load);
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
    <div v-if="kind === 'storage'" class="flex gap-2 px-3">
      <button class="btn" :disabled="busy || !disk || !(disk.attached ? disk.can_deactivate : disk.can_activate)" @click="toggle">
        {{ disk && !disk.attached ? "Activate" : "Deactivate" }}
      </button>
    </div>
  </div>
</template>
