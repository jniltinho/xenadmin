<script setup lang="ts">
// The GPU tab of a pool or server (GpuView: placement policy, GPUs by settings, allowed vGPU types and their editor)
// and the USB tab of a server (UsbView: the devices and Enable / Disable pass-through); read again on tick.
import { computed, onMounted, reactive, ref, watch } from "vue";
import { api, type GPUTab, type USBTab } from "./api/workspace";
import DetailsTable from "./DetailsTable.vue";
import { serial } from "./serial";

const props = defineProps<{ server: string; kind: string; uuid: string; tab: "gpu" | "usb"; tick?: number }>();
const emit = defineEmits<{ (e: "status", text: string, error: boolean): void }>();

const gpu = ref<GPUTab | null>(null);
const usb = ref<USBTab | null>(null);
const picked = ref("");
const busy = ref(false);
const load = serial(async () => {
  try {
    if (props.tab === "gpu") gpu.value = await api.gpu(props.server, props.kind, props.uuid);
    else {
      usb.value = await api.usb(props.server, props.uuid);
      if (!usb.value.table.keys.includes(picked.value)) picked.value = usb.value.table.keys[0] ?? "";
    }
  } catch (e) { emit("status", String(e), true); }
});
onMounted(load);
watch(() => props.tick, load);

async function run(label: string, f: () => Promise<void>) {
  busy.value = true;
  emit("status", `${label}...`, false);
  try { await f(); emit("status", "", false); } catch (e) { emit("status", `${label}: ${String(e)}`, true); }
  await load();
  busy.value = false;
}

// PoolGpuEditPage
const placement = ref<string | null>(null);
function savePlacement() {
  const chosen = placement.value;
  placement.value = null;
  if (chosen && chosen !== gpu.value?.policy_value) run("Setting the GPU placement policy", () => api.setGpuPlacement(props.server, chosen));
}

// GpuConfiguration: tick the allowed types; the ones in use cannot be cleared; only changed boxes are applied
const editing = ref<GPUTab["groups"][number] | null>(null);
const ticks = reactive<Record<string, boolean>>(Object.create(null));
function editTypes(row: GPUTab["groups"][number]) {
  for (const k of Object.keys(ticks)) delete ticks[k];
  for (const t of row.types) ticks[t.uuid] = t.enabled;
  editing.value = row;
}
function saveTypes() {
  const row = editing.value;
  editing.value = null;
  if (!row) return;
  const changes: Record<string, boolean> = {};
  for (const t of row.types) if (!t.in_use && ticks[t.uuid] !== t.enabled) changes[t.uuid] = ticks[t.uuid];
  if (Object.keys(changes).length) run("Configuring GPUs", () => api.setVgpuTypes(props.server, row.pgpus, changes));
}

// UsbPage: the question, then the change if the device is still as it was shown
const device = computed(() => usb.value?.devices.find((d) => d.uuid === picked.value) ?? null);
const ask = ref<{ uuid: string; was: boolean; text: string; yes: string } | null>(null);
function answer(ok: boolean) {
  const a = ask.value;
  ask.value = null;
  if (ok && a) run("Changing USB pass-through", () => api.setUsbPassthrough(props.server, a.uuid, a.was));
}
</script>

<template>
  <div class="flex min-h-0 flex-1 flex-col overflow-auto">
    <div v-if="tab === 'gpu' && gpu" class="max-w-[1100px] space-y-3 p-3">
      <div v-if="gpu.policy" class="flex items-center gap-3 border border-[#DCDCDC] bg-white p-2.5">
        <span class="flex-1">{{ gpu.policy }}</span>
        <button class="btn" :disabled="busy" @click="placement = gpu.policy_value">Edit...</button>
      </div>
      <div v-for="(g, i) in gpu.groups" :key="i" class="space-y-1.5 border border-[#DCDCDC] bg-white p-2.5">
        <template v-for="(h, j) in g.hosts" :key="j">
          <p class="font-semibold">{{ h.label }}</p>
          <p v-for="n in h.gpus" :key="n" class="whitespace-pre pl-6">{{ n }}</p>
        </template>
        <table class="mt-1.5">
          <thead><tr><th /><th class="px-2 text-left font-semibold">Allowed vGPU types</th><th class="px-2 text-left font-semibold">Max per GPU</th><th class="px-2 text-left font-semibold">Video RAM</th></tr></thead>
          <tbody>
            <tr v-for="t in g.types" :key="t.uuid">
              <td><img :src="`/icons/${t.enabled ? '000_Tick_h32bit_16' : '000_Abort_h32bit_16'}.png`" alt="" class="h-4 w-4" /></td>
              <td class="px-2">{{ t.name }}</td><td class="px-2">{{ t.max }}</td><td class="px-2">{{ t.vram }}</td>
            </tr>
          </tbody>
        </table>
        <button v-if="g.can_edit" class="btn mt-1.5" :disabled="busy" @click="editTypes(g)">{{ g.edit }}</button>
      </div>
      <p v-if="!gpu.groups.length">{{ gpu.no_gpu }}</p>
    </div>

    <template v-if="tab === 'usb' && usb">
      <DetailsTable :table="usb.table" :selected="picked" @select="picked = $event" />
      <div class="flex gap-2 px-3">
        <button class="btn min-w-36" :disabled="busy || !device?.can_toggle"
                @click="device && (ask = { uuid: device.uuid, was: device.passthrough, text: device.question, yes: device.yes })">
          {{ device?.button ?? "Enable pass-through" }}
        </button>
      </div>
    </template>

    <Teleport to="body">
      <div v-if="placement !== null && gpu" class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
        <form class="w-[32rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="savePlacement">
          <div class="bg-band-dark px-3 py-1 font-semibold text-white">GPU Placement Policy</div>
          <div class="space-y-2 p-3">
            <p>Choose how VMs with vGPUs are placed on the GPUs of this pool:</p>
            <label class="flex items-center gap-2"><input v-model="placement" type="radio" value="depth_first" />{{ gpu.density }}</label>
            <label class="flex items-center gap-2"><input v-model="placement" type="radio" value="breadth_first" />{{ gpu.performance }}</label>
            <div class="flex justify-end gap-2 pt-1">
              <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5">OK</button>
              <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="placement = null">Cancel</button>
            </div>
          </div>
        </form>
      </div>
      <div v-if="editing" class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
        <form class="w-[35rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="saveTypes">
          <div class="bg-band-dark px-3 py-1 font-semibold text-white">GPU Configuration</div>
          <div class="space-y-2 p-3">
            <p class="font-semibold">{{ editing.title }}</p>
            <p>{{ editing.rubric }}</p>
            <div class="max-h-[300px] space-y-1 overflow-auto">
              <label v-for="t in editing.types" :key="t.uuid" class="flex items-center gap-2">
                <input v-model="ticks[t.uuid]" type="checkbox" :disabled="t.in_use" />{{ t.in_use ? `${t.name}  (in use)` : t.name }}
              </label>
            </div>
            <div class="flex justify-end gap-2 pt-1">
              <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5">OK</button>
              <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="editing = null">Cancel</button>
            </div>
          </div>
        </form>
      </div>
      <div v-if="ask" class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
        <div class="w-[32rem] border border-neutral-500 bg-white shadow-lg">
          <div class="bg-band-dark px-3 py-1 font-semibold text-white">USB</div>
          <div class="space-y-3 p-3">
            <p class="whitespace-pre-line">{{ ask.text }}</p>
            <div class="flex justify-end gap-2">
              <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="answer(true)">{{ ask.yes }}</button>
              <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" autofocus @click="answer(false)">Cancel</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
