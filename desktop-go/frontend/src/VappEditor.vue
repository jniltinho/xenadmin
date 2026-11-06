<script setup lang="ts">
// The New vApp wizard (NewVMApplianceWizard: Name, Virtual Machines, VM Startup Sequence, Finish) and a vApp's
// Properties (General, Virtual Machines, VM Startup Sequence), with the question when VMs leave another vApp.
import { computed, reactive, ref } from "vue";
import { api, type Vapp, type Vapps } from "./api/workspace";

const props = defineProps<{ server: string; data: Vapps; vapp: Vapp | null; preselect?: string }>();
// the candidates as they were when the editor opened (Manage vApps reads them again every second)
const data: Vapps = JSON.parse(JSON.stringify(props.data));
const emit = defineEmits<{ (e: "close", saved: boolean): void }>();

const pages = computed(() => (props.vapp ? ["General", "Virtual Machines", "VM Startup Sequence"] : ["Name", "Virtual Machines", "VM Startup Sequence", "Finish"]));
const page = ref(0);
const f = reactive({ name: props.vapp?.name ?? "", description: props.vapp?.description ?? "",
  vms: new Set<string>(props.vapp ? props.vapp.vms.map((v) => v.uuid) : props.preselect ? [props.preselect] : []) });
const startup = reactive<Record<string, { order: number; delay: number }>>(Object.create(null));
for (const v of data.candidates) startup[v.uuid] = { order: v.order, delay: v.delay };
const error = ref("");
const pending = ref(false);
const ask = ref("");
const summary = ref("");
const chosen = computed(() => data.candidates.filter((v) => f.vms.has(v.uuid)));
const vappName = (uuid: string) => data.vapps.find((a) => a.uuid === uuid)?.name ?? "";

async function go(i: number) {
  error.value = "";
  if (!props.vapp && page.value === 0 && i > 0 && !f.name.trim()) { error.value = "Enter a name for the vApp."; return; }
  page.value = i;
  if (pages.value[i] === "Finish") summary.value = await api.vappSummary(f.name.trim(), chosen.value.map((v) => v.name));
}
async function finish(confirmed = false) {
  if (pending.value) return;
  if (!f.name.trim()) { error.value = "Enter a name for the vApp."; page.value = 0; return; }
  pending.value = true; // before the first await: a second click never creates a second vApp
  if (!confirmed) {
    const moves = chosen.value.filter((v) => v.vapp && v.vapp !== props.vapp?.uuid).map((v) => ({ vm: v.name, from: vappName(v.vapp) }));
    const q = await api.vappChangeWarning(moves, f.name.trim()).catch(() => "");
    if (q) { ask.value = q; pending.value = false; return; }
  }
  pending.value = false;
  ask.value = "";
  const changed: Record<string, { order: number; delay: number }> = Object.create(null);
  for (const v of chosen.value) {
    const s = startup[v.uuid];
    if (!Number.isInteger(s.order) || s.order < 0 || !Number.isInteger(s.delay) || s.delay < 0) { error.value = "Enter whole numbers of 0 or more."; page.value = 2; return; }
    if (s.order !== v.order || s.delay !== v.delay) changed[v.uuid] = s;
  }
  pending.value = true;
  const spec = { name: f.name.trim(), description: f.description, vms: [...f.vms], startup: changed };
  try {
    if (props.vapp) await api.saveVapp(props.server, props.vapp.uuid, spec);
    else await api.createVapp(props.server, spec);
    emit("close", true);
  } catch (e) { error.value = String(e); } finally { pending.value = false; }
}
function toggle(uuid: string, on: boolean) { if (on) f.vms.add(uuid); else f.vms.delete(uuid); }
</script>

<template>
  <div class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
    <div class="flex h-[30rem] w-[46rem] flex-col border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ vapp ? `'${vapp.name}' Properties` : `Create vApp - ${data.pool}` }}</div>
      <div class="flex min-h-0 flex-1">
        <div class="w-48 border-r border-line bg-pane py-1">
          <div v-for="(p, i) in pages" :key="p" class="cursor-default px-2 py-1.5" :class="page === i ? 'bg-select text-white' : vapp ? 'hover:bg-neutral-100' : 'text-neutral-500'"
               @click="vapp && go(i)">{{ p }}</div>
        </div>
        <fieldset :disabled="pending" class="min-w-0 flex-1 overflow-auto p-4">
          <template v-if="page === 0">
            <p class="mb-3">Enter a name and an optional description for the {{ vapp ? "" : "new " }}vApp.</p>
            <label class="mb-2 grid grid-cols-[7rem_1fr] items-center gap-2">Name:<input v-model="f.name" class="border border-line px-1 py-0.5" /></label>
            <label class="grid grid-cols-[7rem_1fr] items-start gap-2">Description:<textarea v-model="f.description" rows="3" class="border border-line px-1 py-0.5" /></label>
          </template>
          <template v-else-if="page === 1">
            <p class="mb-2">Select the VMs to include in this vApp.</p>
            <div class="h-72 overflow-auto border border-line">
              <label v-for="v in data.candidates" :key="v.uuid" class="flex items-center gap-2 px-1 py-0.5 hover:bg-neutral-100">
                <input type="checkbox" :checked="f.vms.has(v.uuid)" @change="toggle(v.uuid, ($event.target as HTMLInputElement).checked)" />
                <span class="flex-1">{{ v.name }}</span><span class="w-40 truncate text-neutral-500">{{ v.vapp && v.vapp !== vapp?.uuid ? vappName(v.vapp) : "" }}</span></label>
            </div>
          </template>
          <template v-else-if="page === 2">
            <p class="mb-2">Specify the startup sequence of the VMs in the vApp: the start order, and the delay before the next VM starts.</p>
            <table v-if="chosen.length" class="w-full border-collapse">
              <thead class="bg-pane text-left"><tr><th class="px-2 py-1 font-normal">VM</th><th class="px-2 py-1 font-normal">Start order</th><th class="px-2 py-1 font-normal">Attempt to start next VM after (seconds)</th></tr></thead>
              <tbody><tr v-for="v in chosen" :key="v.uuid" class="border-b border-line">
                <td class="px-2 py-1">{{ v.name }}</td>
                <td class="px-2 py-1"><input v-model.number="startup[v.uuid].order" type="number" min="0" :aria-label="`Start order of ${v.name}`" class="w-20 border border-line px-1" /></td>
                <td class="px-2 py-1"><input v-model.number="startup[v.uuid].delay" type="number" min="0" :aria-label="`Delay of ${v.name}`" class="w-20 border border-line px-1" /></td>
              </tr></tbody>
            </table>
            <p v-else class="text-neutral-500">No VMs are selected.</p>
          </template>
          <template v-else>
            <p class="mb-2">The new vApp will be created with the following settings:</p>
            <pre class="border border-line bg-pane p-2 whitespace-pre-wrap">{{ summary }}</pre>
          </template>
        </fieldset>
      </div>
      <div class="flex items-center justify-end gap-2 border-t border-line p-2">
        <span v-if="error" class="mr-auto text-red-600">{{ error }}</span>
        <template v-if="!vapp">
          <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="page === 0" @click="go(page - 1)">&lt; Previous</button>
          <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="page === pages.length - 1" @click="go(page + 1)">Next &gt;</button>
        </template>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="pending || (!vapp && page !== pages.length - 1)" @click="finish()">{{ vapp ? "OK" : "Finish" }}</button>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close', false)">Cancel</button>
      </div>
    </div>
    <div v-if="ask" class="fixed inset-0 z-50 flex items-center justify-center bg-black/20">
      <div class="w-[30rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Change vApp</div>
        <div class="space-y-3 p-3">
          <p class="whitespace-pre-line">{{ ask }}</p>
          <div class="flex justify-end gap-2">
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="finish(true)">Yes</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="ask = ''">No</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
