<script setup lang="ts">
// Memory Settings (BallooningDialog): a fixed memory or a range (VMMemoryControlsBasic), or the dynamic minimum,
// maximum and static maximum (VMMemoryControlsAdvanced); the Go side validates (MemoryRules.Validate) and asks first
// when a running VM must restart (MemoryRules.Confirmation).
import { computed, onMounted, reactive, ref, watch } from "vue";
import { api, type MemoryEdit } from "./api/workspace";
import { serial } from "./serial";

const props = defineProps<{ server: string; vm: string }>();
const emit = defineEmits<{ (e: "close"): void; (e: "run", label: string, f: () => Promise<void>): void }>();

const m = ref<MemoryEdit | null>(null);
const f = reactive({ dmin: 0, dmax: 0, smax: 0, range: false });
const shown = { dmin: 0, dmax: 0, smax: 0 }; // a field left as shown keeps the server's exact bytes
const error = ref("");
const question = ref("");
const unit = computed(() => (m.value?.in_mb ? 1024 ** 2 : 1024 ** 3));
const show = (bytes: number) => (m.value?.in_mb ? Math.round(bytes / unit.value) : Math.round((bytes / unit.value) * 10) / 10);

onMounted(async () => {
  try {
    const e = await api.memoryEdit(props.server, props.vm);
    m.value = e;
    Object.assign(shown, { dmin: show(e.dynamic_min), dmax: show(e.dynamic_max), smax: show(e.static_max) });
    Object.assign(f, shown, { range: e.dynamic && e.dynamic_min !== e.static_max });
  } catch (e) { error.value = String(e); }
});

function bytes(k: "dmin" | "dmax" | "smax") {
  const e = m.value!;
  if (f[k] === shown[k]) return { dmin: e.dynamic_min, dmax: e.dynamic_max, smax: e.static_max }[k];
  return Math.round(f[k] * unit.value);
}
const values = (): [number, number, number] => m.value!.advanced ? [bytes("dmin"), bytes("dmax"), bytes("smax")]
  : f.range ? [bytes("dmin"), bytes("dmax"), bytes("dmax")] : [bytes("smax"), bytes("smax"), bytes("smax")];

// the fields the mode uses (an unused field may be left blank)
const used = (): number[] => m.value!.advanced ? [f.dmin, f.dmax, f.smax] : f.range ? [f.dmin, f.dmax] : [f.smax];
const check = serial(async () => {
  if (!m.value) return;
  if (!used().every(Number.isFinite)) { error.value = "Enter a number."; return; }
  try { question.value = await api.memoryQuestion(props.server, props.vm, ...values()); error.value = ""; }
  catch (e) { error.value = String(e); }
});
watch(() => [f.dmin, f.dmax, f.smax, f.range, m.value], () => check());

const ask = ref(false);
const pending = ref(false); // one OK or Yes at a time: a second click never stands for the Yes
// OK asks first when a running VM must restart; Yes (confirmed) checks again and applies
async function ok(confirmed = false) {
  if (pending.value) return;
  pending.value = true;
  try {
    await check();
    if (error.value) { ask.value = false; return; } // back to the values, with the reason
    if (question.value && !confirmed) { ask.value = true; return; }
    const [a, b, c] = values(), name = m.value!.name;
    emit("run", `Changing memory settings for '${name}'`, () => api.changeMemory(props.server, props.vm, a, b, c));
    emit("close");
  } finally { pending.value = false; }
}
const u = computed(() => (m.value?.in_mb ? "MB" : "GB"));
const step = computed(() => (m.value?.in_mb ? 64 : 0.1));
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <form novalidate class="w-[30rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="ok(ask)">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">Memory Settings</div>
      <div v-if="m && ask" class="space-y-3 p-3">
        <p class="whitespace-pre-line">{{ question }}</p>
        <div class="flex justify-end gap-2">
          <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="pending">Yes</button>
          <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="pending" @click="ask = false">No</button>
        </div>
      </div>
      <div v-else-if="m" class="space-y-3 p-3">
        <p>{{ m.rubric }}</p>
        <div v-if="m.advanced" class="grid grid-cols-[auto_8rem_auto] items-center gap-x-3 gap-y-2">
          <label for="dmin">Minimum memory:</label><input id="dmin" v-model.number="f.dmin" type="number" min="0" :step="step" class="border border-line px-1 py-0.5" /><span>{{ u }}</span>
          <label for="dmax">Maximum memory:</label><input id="dmax" v-model.number="f.dmax" type="number" min="0" :step="step" class="border border-line px-1 py-0.5" /><span>{{ u }}</span>
          <label for="smax">Static maximum:</label><input id="smax" v-model.number="f.smax" type="number" min="0" :step="step" class="border border-line px-1 py-0.5" /><span>{{ u }}</span>
        </div>
        <template v-else>
          <div class="flex items-center gap-2">
            <label class="flex items-center gap-2"><input v-model="f.range" type="radio" :value="false" />Set a fixed memory of</label>
            <input v-model.number="f.smax" type="number" min="0" :step="step" aria-label="Fixed memory" class="w-32 border border-line px-1 py-0.5" @focus="f.range = false" /><span>{{ u }}</span>
          </div>
          <label class="flex items-center gap-2" :class="{ 'text-neutral-400': !m.dynamic }">
            <input v-model="f.range" type="radio" :value="true" :disabled="!m.dynamic" />Automatically allocate memory within this range:</label>
          <div class="ml-6 grid grid-cols-[auto_8rem_auto] items-center gap-x-3 gap-y-2" :class="{ 'text-neutral-400': !m.dynamic }">
            <label for="bmin">Minimum memory:</label><input id="bmin" v-model.number="f.dmin" type="number" min="0" :step="step" :disabled="!m.dynamic" class="border border-line px-1 py-0.5" @focus="f.range = true" /><span>{{ u }}</span>
            <label for="bmax">Maximum memory:</label><input id="bmax" v-model.number="f.dmax" type="number" min="0" :step="step" :disabled="!m.dynamic" class="border border-line px-1 py-0.5" @focus="f.range = true" /><span>{{ u }}</span>
          </div>
          <p v-if="m.unavailable" class="text-neutral-600">{{ m.unavailable }}</p>
        </template>
        <p v-if="error" class="text-red-600">{{ error }}</p>
        <div class="flex justify-end gap-2 pt-1">
          <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="!!error || !m.can_edit || pending">OK</button>
          <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Cancel</button>
        </div>
      </div>
      <div v-else class="space-y-3 p-3">
        <p :class="{ 'text-red-600': error }">{{ error || "Loading..." }}</p>
        <div class="flex justify-end"><button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Cancel</button></div>
      </div>
    </form>
  </div>
</template>
