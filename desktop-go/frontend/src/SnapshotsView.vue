<script setup lang="ts">
// The Snapshots tab (SnapshotsPage / SnapshotsView): the VM's snapshots, Take Snapshot (VmSnapshotDialog),
// Revert To (with "take a snapshot first") and Delete, with the C# questions from the Go side.
import { computed, onMounted, reactive, ref, watch } from "vue";
import { serial } from "./serial";
import { api, type Snapshot, type VMSnapshots } from "./api/workspace";

// tick: the server's objects changed (xcp:changed), read the snapshots again
// askTake: VM > Take a Snapshot... opens the dialog (answered with "asked")
const props = defineProps<{ server: string; vm: string; tick?: number; askTake?: boolean }>();
const emit = defineEmits<{ (e: "status", text: string, error: boolean): void; (e: "changed"): void; (e: "asked"): void }>();

const data = ref<VMSnapshots | null>(null);
const picked = ref<string>("");
const busy = ref(false);
const take = reactive({ open: false, name: "", description: "", memory: false, error: "" });
const ask = ref<{ text: string; option?: string; checked: boolean; resolve: (r: { ok: boolean; checked: boolean }) => void } | null>(null);
const selected = computed(() => data.value?.snapshots.find((s) => s.uuid === picked.value) ?? null);

// Util.DiskSizeString: binary units with one decimal
function size(bytes: number) {
  const units = ["B", "KB", "MB", "GB", "TB"];
  let i = 0, v = bytes;
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++; }
  return `${i ? v.toFixed(1) : v} ${units[i]}`;
}
const when = (t: string) => new Date(t).toLocaleString(undefined, { dateStyle: "medium", timeStyle: "medium" });

// one read at a time (serial): changes coming faster than a slow read make one more read; an answer for the VM shown
// before is dropped
const load = serial(async () => {
  const vm = props.vm, server = props.server;
  try {
    const got = await api.snapshots(server, vm);
    if (vm !== props.vm || server !== props.server) return;
    data.value = got;
    if (!data.value.snapshots.some((s) => s.uuid === picked.value)) picked.value = "";
  } catch (e) { if (vm === props.vm && server === props.server) emit("status", String(e), true); }
});

async function run(label: string, f: () => Promise<void>) {
  busy.value = true;
  emit("status", `${label}...`, false);
  try { await f(); emit("status", "", false); } catch (e) { emit("status", `${label}: ${String(e)}`, true); }
  busy.value = false;
  await load();
  emit("changed");
}

function openTake() {
  Object.assign(take, { open: true, name: "", description: "", memory: false, error: "" });
}
watch(() => props.askTake, (v) => { if (v) { openTake(); emit("asked"); } }, { immediate: true });

async function submitTake(): Promise<boolean> {
  if (!take.name.trim()) { take.error = "Enter a name for the snapshot."; return false; }
  take.open = false;
  let ok = true;
  await run(`Taking snapshot '${take.name.trim()}'`, async () => {
    try { await api.takeSnapshot(props.server, props.vm, take.name.trim(), take.description, take.memory); } catch (e) { ok = false; throw e; }
  });
  return ok;
}

function question(text: string, option?: string) {
  return new Promise<{ ok: boolean; checked: boolean }>((resolve) => (ask.value = { text, option, checked: false, resolve }));
}
function answer(ok: boolean) { ask.value?.resolve({ ok, checked: ask.value.checked }); ask.value = null; }

let afterTake: ((ok: boolean) => void) | null = null;

async function revert(s: Snapshot) {
  const q = await api.snapshotQuestion("revert", s.name, !!data.value?.snapshot_allowed);
  const a = await question(q.text ?? "", q.option);
  if (!a.ok) return;
  if (a.checked) { // RevertToSnapshotCommand: take a snapshot of the current state, then revert
    openTake();
    const taken = await new Promise<boolean>((resolve) => (afterTake = resolve));
    if (!taken) return;
  }
  await run(`Reverting to '${s.name}'`, () => api.revertSnapshot(props.server, s.uuid));
}

async function del(s: Snapshot) {
  const q = await api.snapshotQuestion("delete", s.name, false);
  if (!(await question(q.text ?? "")).ok) return;
  await run(`Deleting snapshot '${s.name}'`, () => api.deleteSnapshot(props.server, s.uuid));
}

async function confirmTake() { const ok = await submitTake(); afterTake?.(ok); afterTake = null; }
function cancelTake() { take.open = false; afterTake?.(false); afterTake = null; }

watch(() => [props.server, props.vm, props.tick], load);
onMounted(load);
</script>

<template>
  <div class="flex min-h-0 flex-1 flex-col">
    <div class="bg-gradient-to-r from-band-dark to-band px-3 py-1.5 font-semibold text-white">Snapshots</div>
    <div class="flex items-center gap-2 border-b border-line bg-pane px-2 py-1.5">
      <button class="btn" :disabled="busy || !(data?.modes.disk || data?.modes.memory)" @click="openTake">Take Snapshot...</button>
      <button class="btn" :disabled="busy || !selected" @click="selected && revert(selected)">Revert To...</button>
      <button class="btn" :disabled="busy || !selected" @click="selected && del(selected)">Delete</button>
    </div>
    <div class="min-h-0 flex-1 overflow-auto">
      <table v-if="data?.snapshots.length" class="w-full border-collapse">
        <thead class="sticky top-0 bg-pane text-left">
          <tr><th class="border-b border-line px-2 py-1 font-normal">Name</th><th class="w-48 border-b border-line px-2 py-1 font-normal">Created</th>
            <th class="w-36 border-b border-line px-2 py-1 font-normal">Type</th><th class="w-24 border-b border-line px-2 py-1 font-normal">Size</th>
            <th class="border-b border-line px-2 py-1 font-normal">Description</th></tr>
        </thead>
        <tbody>
          <tr v-for="s in data.snapshots" :key="s.uuid" class="cursor-default" :class="picked === s.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100'"
              @click="picked = s.uuid" @dblclick="revert(s)">
            <td class="border-b border-line px-2 py-1">
              <span class="flex items-center gap-1.5"><img :src="`/icons/${s.type === 'Disks and memory' ? '000_VMSnapshotDiskMemory' : '000_VMSnapShotDiskOnly'}_h32bit_16.png`" alt="" class="h-4 w-4" />{{ s.name }}</span>
            </td>
            <td class="border-b border-line px-2 py-1">{{ when(s.time) }}</td>
            <td class="border-b border-line px-2 py-1">{{ s.type }}</td>
            <td class="border-b border-line px-2 py-1">{{ size(s.size) }}</td>
            <td class="border-b border-line px-2 py-1 select-text">{{ s.description }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else-if="data" class="p-3 text-neutral-500">There are currently no snapshots of this VM. To create a snapshot, click Take Snapshot.</p>
    </div>

    <div v-if="take.open" class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
      <form class="w-[30rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="confirmTake">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Take Snapshot</div>
        <div class="space-y-2 p-3">
          <label class="grid grid-cols-[6rem_1fr] items-center gap-2">Name:<input v-model="take.name" class="border border-line px-1 py-0.5" autofocus /></label>
          <label class="grid grid-cols-[6rem_1fr] items-start gap-2">Description:<textarea v-model="take.description" rows="3" class="border border-line px-1 py-0.5" /></label>
          <fieldset class="space-y-1 border border-line p-2">
            <legend class="px-1">Snapshot mode</legend>
            <label class="flex items-center gap-2"><input v-model="take.memory" type="radio" :value="false" :disabled="!data?.modes.disk" />Snapshot the virtual machine's disks</label>
            <label class="flex items-center gap-2"><input v-model="take.memory" type="radio" :value="true" :disabled="!data?.modes.memory" />Snapshot the virtual machine's disks and memory</label>
            <p v-if="!data?.modes.memory" class="pl-6 text-xs text-neutral-500">{{ data?.modes.memory_info }}</p>
          </fieldset>
          <p v-if="take.error" class="text-red-600">{{ take.error }}</p>
          <div class="flex justify-end gap-2 pt-1">
            <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5">OK</button>
            <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="cancelTake">Cancel</button>
          </div>
        </div>
      </form>
    </div>

    <div v-if="ask" class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
      <div class="w-[32rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">XCP-ng Center</div>
        <div class="space-y-3 p-3">
          <p class="whitespace-pre-line">{{ ask.text }}</p>
          <label v-if="ask.option" class="flex items-center gap-2"><input v-model="ask.checked" type="checkbox" />{{ ask.option }}</label>
          <div class="flex justify-end gap-2">
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="answer(true)">Yes</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="answer(false)">No</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
