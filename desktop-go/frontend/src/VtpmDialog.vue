<script setup lang="ts">
// VM > Manage vTPMs (VtpmManagementDialog): the VM's vTPMs on the left, the selected one's properties (Protected,
// Unique) with Remove vTPM, and Add vTPM; the Go side checks the server's allowed operations.
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { api, type Vtpms } from "./api/workspace";

const props = defineProps<{ server: string; vm: string }>();
const emit = defineEmits<{ (e: "close"): void }>();
const data = ref<Vtpms | null>(null);
const picked = ref("");
const error = ref("");
const busy = ref(false);
const ask = ref("");
const current = computed(() => data.value?.vtpms.find((t) => t.uuid === picked.value) ?? data.value?.vtpms[0] ?? null);

async function read() {
  try { data.value = await api.vtpms(props.server, props.vm); } catch (e) { error.value = String(e); }
}
// read again every two seconds: the VM may change elsewhere (VtpmManagementDialog follows the cache)
let timer: ReturnType<typeof setInterval> | undefined;
onMounted(() => { read(); timer = setInterval(() => { if (!busy.value) read(); }, 2000); });
onBeforeUnmount(() => clearInterval(timer));
async function run(f: () => Promise<void>) {
  if (busy.value) return;
  busy.value = true;
  error.value = "";
  try { await f(); } catch (e) { error.value = String(e); }
  await read();
  busy.value = false;
}
async function remove() { ask.value = await api.vtpmRemoveWarning(); }
function confirmRemove() { const t = current.value; ask.value = ""; if (t) run(() => api.removeVtpm(props.server, props.vm, t.uuid)); }
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <div class="flex h-[24rem] w-[40rem] flex-col border border-neutral-500 bg-white shadow-lg">
      <div class="truncate bg-band-dark px-3 py-1 font-semibold text-white">Trusted Platform Module (TPM) Manager - {{ data?.name ?? "" }}</div>
      <p v-if="data?.blocker" class="p-3">{{ data.blocker }}</p>
      <div v-else class="flex min-h-0 flex-1">
        <div class="flex w-48 flex-col border-r border-line bg-pane">
          <div class="min-h-0 flex-1 py-1">
            <div v-for="(t, i) in data?.vtpms ?? []" :key="t.uuid" class="cursor-default px-2 py-1.5" :class="current?.uuid === t.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100'"
                 @click="picked = t.uuid">{{ data!.vtpms.length === 1 ? "vTPM" : `vTPM ${i + 1}` }}</div>
          </div>
          <div class="p-2"><button class="btn w-full" :disabled="busy || !data || !!data.add_blocker" :title="data?.add_blocker" @click="run(() => api.addVtpm(server, vm))">Add vTPM</button></div>
        </div>
        <div class="min-w-0 flex-1 p-4">
          <template v-if="current">
            <p class="mb-2 font-bold">Properties</p>
            <div class="grid grid-cols-[7rem_1fr] gap-1"><span>Protected:</span><span>{{ current.protected ? "Yes" : "No" }}</span><span>Unique:</span><span>{{ current.unique ? "Yes" : "No" }}</span></div>
            <button class="btn mt-4" :disabled="busy || !!current.remove_blocker" :title="current.remove_blocker" @click="remove">Remove vTPM</button>
          </template>
          <p v-else-if="data" class="pt-16 text-center">There are no vTPMs attached to this VM.</p>
        </div>
      </div>
      <p v-if="error" class="px-3 text-red-600">{{ error }}</p>
      <div class="flex justify-end border-t border-line p-2"><button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Close</button></div>
    </div>
    <div v-if="ask" class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
      <div class="w-[28rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Remove vTPM</div>
        <div class="space-y-3 p-3"><p>{{ ask }}</p>
          <div class="flex justify-end gap-2">
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="confirmRemove">Yes</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="ask = ''">No</button>
          </div></div>
      </div>
    </div>
  </div>
</template>
