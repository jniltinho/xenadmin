<script setup lang="ts">
// The New Pool dialog (NewPoolDialog): name, description, the coordinator and the additional members among the
// connected servers, each with the reason it cannot join (PoolJoinRules); Create Pool runs CreatePoolAction on the Go
// side (the pool's name, then pool.join from each member).
import { onMounted, reactive, ref } from "vue";
import { api, type PoolCandidate } from "./api/workspace";

// close: the coordinator of the pool created ("" for none) and an error to show in the window (a member that failed)
const emit = defineEmits<{ (e: "close", created: string, error: string): void; (e: "add-server"): void }>();
const f = reactive({ name: "New Pool (1)", description: "", coordinator: "" }); // the C# dialog's default name
const candidates = ref<PoolCandidate[]>([]);
const members = ref(new Set<string>());
const busy = ref(false);
const error = ref("");

// each read is numbered: an answer for a coordinator chosen before is dropped, and Create waits for the last one
let reads = 0;
const loading = ref(false);
async function load() {
  const n = ++reads;
  loading.value = true;
  try {
    const got = await api.poolCandidates(f.coordinator);
    if (n !== reads) return;
    candidates.value = got;
    f.coordinator = got.find((c) => c.is_coordinator)?.key ?? "";
    members.value = new Set([...members.value].filter((k) => got.some((c) => c.key === k && c.can_join)));
  } catch (e) { if (n === reads) error.value = String(e); } finally { if (n === reads) loading.value = false; }
}
const question = ref("");
let answer: (yes: boolean) => void = () => {};
const ask = (text: string) => new Promise<boolean>((resolve) => { question.value = text; answer = (yes) => { question.value = ""; resolve(yes); }; });
onMounted(load);
defineExpose({ load }); // the window reloads it after Add New Server

function toggle(key: string) { const s = new Set(members.value); if (!s.delete(key)) s.add(key); members.value = s; }

async function create() {
  if (!f.name.trim()) { error.value = "Pool name cannot be empty"; return; }
  if (!f.coordinator) { error.value = "There is no server that can be the coordinator of a new pool."; return; }
  busy.value = true; error.value = "";
  // one snapshot for the warning, the creation and the close (the fields are locked while busy)
  const coordinator = f.coordinator, chosen = [...members.value], name = f.name.trim(), description = f.description;
  try {
    const warning = await api.poolWarning(coordinator, chosen); // PoolRules.Warnings: CPU levelling
    if (warning && !(await ask(warning))) return;
    const r = await api.createPool(coordinator, chosen, name, description);
    if (r.created) { emit("close", coordinator, r.error ?? ""); return; } // a member that failed: the pool exists all the same
    error.value = r.error ?? "";
    await load();
  } catch (e) { error.value = String(e); } finally { busy.value = false; }
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <div class="flex w-[38rem] flex-col border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">Create New Pool</div>
      <fieldset class="p-4" :disabled="busy">
        <p class="mb-3">To create your pool provide a name and select which servers you would like to be added to the pool.</p>
        <label class="mb-2 grid grid-cols-[7rem_1fr] items-center gap-2">Name:<input v-model="f.name" class="border border-line px-1 py-0.5" /></label>
        <label class="mb-2 grid grid-cols-[7rem_1fr] items-center gap-2">Description:<input v-model="f.description" placeholder="(optional)" class="border border-line px-1 py-0.5" /></label>
        <label class="mb-3 grid grid-cols-[7rem_1fr] items-center gap-2">Coordinator:
          <select v-model="f.coordinator" class="border border-line bg-white px-1 py-0.5" @change="load">
            <option v-for="c in candidates.filter((x) => x.can_coordinate)" :key="c.key" :value="c.key">{{ c.name || c.key }}</option>
          </select></label>
        <p class="mb-1">Additional members:</p>
        <div class="h-40 overflow-auto border border-neutral-400">
          <label v-for="c in candidates.filter((x) => x.key !== f.coordinator)" :key="c.key" class="flex items-center gap-2 px-2 py-0.5"
                 :class="c.can_join ? 'hover:bg-neutral-100' : 'text-neutral-400'">
            <input type="checkbox" :disabled="!c.can_join" :checked="members.has(c.key)" @change="toggle(c.key)" />
            <span class="flex-1">{{ c.name || c.key }}</span><span class="text-xs">{{ c.reason }}</span></label>
          <p v-if="candidates.length <= 1" class="p-2 text-neutral-500">Connect to the servers that will join the pool with Add New Server.</p>
        </div>
        <button class="btn mt-2" :disabled="busy" @click="emit('add-server')">Add New Server</button>
      </fieldset>
      <div class="flex items-center justify-end gap-2 border-t border-line p-2">
        <span v-if="error" class="mr-auto text-red-600">{{ error }}</span>
        <button class="btn min-w-24" :disabled="busy || loading || !f.coordinator" @click="create">{{ busy ? "Creating..." : "Create Pool" }}</button>
        <button class="btn min-w-20" :disabled="busy" @click="emit('close', '', '')">Cancel</button>
      </div>
    </div>
    <div v-if="question" class="fixed inset-0 z-40 flex items-center justify-center bg-black/10">
      <div class="w-[34rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Create New Pool</div>
        <p class="whitespace-pre-line p-4">{{ question }}</p>
        <div class="flex justify-end gap-2 border-t border-line p-2">
          <button class="btn min-w-20" @click="answer(true)">Yes</button><button class="btn min-w-20" @click="answer(false)">No</button>
        </div>
      </div>
    </div>
  </div>
</template>
