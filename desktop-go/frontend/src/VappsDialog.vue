<script setup lang="ts">
// Manage vApps (VMAppliancesDialog / XenAdmin.Desktop VappsDialog): the vApps of a pool with New vApp, Delete,
// Properties, Start and Shut down, and the VMs of the selected one; read again every second while open.
// Export opens the Export wizard for the vApp (the window's ExportApplianceDialog, OVF/OVA only).
// shortcut: Import of a vApp (OVF/OVA) is not here; File > Import... imports an appliance.
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { api, type Vapps } from "./api/workspace";
import VappEditor from "./VappEditor.vue";

const props = defineProps<{ server: string }>();
const emit = defineEmits<{ (e: "close"): void; (e: "export", uuid: string, name: string): void }>();

const data = ref<Vapps | null>(null);
const picked = ref("");
const error = ref("");
const busy = ref(false);
const editing = ref<{ uuid: string } | null>(null); // "" for a new vApp
const ask = ref<{ text: string; resolve: (ok: boolean) => void } | null>(null);
const current = computed(() => data.value?.vapps.find((a) => a.uuid === picked.value) ?? null);
let reading = false;
async function read() {
  if (reading) return;
  reading = true;
  try { data.value = await api.vapps(props.server); } catch (e) { error.value = String(e); } finally { reading = false; }
}
let timer: ReturnType<typeof setInterval> | undefined;
onMounted(() => { read(); timer = setInterval(read, 1000); });
onBeforeUnmount(() => clearInterval(timer));

const question = (text: string) => new Promise<boolean>((resolve) => (ask.value = { text, resolve }));
function answer(ok: boolean) { ask.value?.resolve(ok); ask.value = null; }
async function run(f: () => Promise<void>) {
  busy.value = true;
  error.value = "";
  try { await f(); } catch (e) { error.value = String(e); }
  await read();
  busy.value = false;
}
async function del() {
  const a = current.value;
  if (a && (await question(await api.vappDeleteQuestion(a.name, a.vms.length)))) run(() => api.deleteVapp(props.server, a.uuid));
}
async function shutDown() {
  const a = current.value;
  if (a && (await question("Are you sure you want to shut down the selected vApp?"))) run(() => api.vappPower(props.server, a.uuid, "shutdown"));
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <div class="flex h-[32rem] w-[48rem] flex-col border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">Manage vApps - pool {{ data?.pool ?? "" }}</div>
      <div class="flex min-h-0 flex-1 flex-col gap-2 p-3">
        <p>{{ data?.is_pool ? `vApps defined in pool '${data?.pool}': ` : `vApps defined in server '${data?.pool ?? ""}': ` }}</p>
        <div class="flex gap-2">
          <button class="btn" :disabled="busy || !data" @click="editing = { uuid: '' }">New vApp...</button>
          <button class="btn" :disabled="busy || !current" @click="del">Delete</button>
          <button class="btn" :disabled="busy || !current" @click="editing = { uuid: current!.uuid }">Properties</button>
          <button class="btn" :disabled="busy || !current?.can_start" @click="run(() => api.vappPower(server, current!.uuid, 'start'))">Start</button>
          <button class="btn" :disabled="busy || !current?.can_shutdown" @click="shutDown">Shut down</button>
          <button class="btn" disabled title="Importing a vApp (OVF/OVA) is not available yet.">Import</button>
          <button class="btn" :disabled="busy || !current" @click="emit('export', current!.uuid, current!.name)">Export</button>
        </div>
        <div class="min-h-0 flex-1 overflow-auto border border-line">
          <table class="w-full border-collapse">
            <thead class="sticky top-0 bg-pane text-left"><tr><th class="px-2 py-1 font-normal">Name</th><th class="px-2 py-1 font-normal">Description</th><th class="w-20 px-2 py-1 font-normal">#VMs</th></tr></thead>
            <tbody><tr v-for="a in data?.vapps ?? []" :key="a.uuid" class="cursor-default" :class="picked === a.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100'"
                       @click="picked = a.uuid" @dblclick="editing = { uuid: a.uuid }">
              <td class="px-2 py-1">{{ a.name }}</td><td class="px-2 py-1">{{ a.description }}</td><td class="px-2 py-1">{{ a.vms.length }}</td></tr></tbody>
          </table>
        </div>
        <div><p class="font-semibold">VMs in this vApp:</p>
          <p>{{ !current ? "" : current.vms.length ? current.vms.map((v) => `${v.name} (${v.power_state})`).join(", ") : "None" }}</p></div>
        <p v-if="error" class="text-red-600">{{ error }}</p>
      </div>
      <div class="flex justify-end border-t border-line p-2"><button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Close</button></div>
    </div>
    <VappEditor v-if="editing && data" :server="server" :data="data" :vapp="data.vapps.find((a) => a.uuid === editing!.uuid) ?? null"
                @close="editing = null; read()" />
    <div v-if="ask" class="fixed inset-0 z-50 flex items-center justify-center bg-black/20">
      <div class="w-[30rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Manage vApps</div>
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
