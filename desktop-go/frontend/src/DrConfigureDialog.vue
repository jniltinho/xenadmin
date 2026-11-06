<script setup lang="ts">
// Pool > Disaster Recovery > Configure (DRConfigureCommand, XenAdmin.Desktop/DrWindows.cs DrConfigureDialog): the SRs
// that can hold the pool's VM metadata (found in the background), ticked where it is stored already; OK enables or
// disables database replication where the choice changed (ops.ConfigureDR checks it all again).
import { computed, onMounted, ref, watch } from "vue";
import { api, type DRSRs } from "./api/workspace";
import { serial } from "./serial";

const props = defineProps<{ server: string }>();
const emit = defineEmits<{ (e: "close"): void; (e: "run", label: string, f: () => Promise<void>): void }>();

const info = ref<DRSRs | null>(null);
const checked = ref<string[]>([]);
const error = ref("");
const warning = ref("");
const canSave = ref(false);
const pending = ref(false);
const originally = computed(() => !!info.value?.srs.some((s) => s.enabled));

onMounted(async () => {
  try {
    info.value = await api.drSrs(props.server);
    checked.value = info.value.srs.filter((s) => s.enabled).map((s) => s.uuid);
    await update();
  } catch (e) { error.value = String(e); }
});

// DRConfigureDialog.ToggleWarningsVisibleState (DrRules.ConfigureWarning); the last change wins
const update = serial(async () => {
  const w = await api.drConfigureWarning(originally.value, checked.value.length);
  warning.value = w.text;
  canSave.value = w.can_save;
});
watch(() => checked.value.length, () => update().catch((e) => { error.value = String(e); }));

function ok() {
  if (pending.value || !info.value || !canSave.value) return;
  const want = [...checked.value];
  const before = info.value.srs.filter((s) => s.enabled).map((s) => s.uuid);
  if (want.length === before.length && want.every((u) => before.includes(u))) { emit("close"); return; } // nothing changed
  pending.value = true;
  emit("run", info.value.action, () => api.configureDr(props.server, want));
  emit("close");
}
const btn = "min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400";
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <div class="w-[36rem] border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ info?.title ?? "Configure DR" }}</div>
      <div class="space-y-3 p-3">
        <p>The configuration data needed to recover VMs and vApps in the event of a disaster must be stored on one or more SRs in this pool.</p>
        <p>Select the SRs where the VM configuration data will be stored:</p>
        <div class="h-48 overflow-auto border border-neutral-400 p-1">
          <p v-if="error" class="text-red-600">{{ error }}</p>
          <p v-else-if="!info" class="text-neutral-600">Looking for the SRs that can store the VM configuration data...</p>
          <p v-else-if="!info.srs.length">There are no SRs in this pool that can store the VM configuration data.</p>
          <label v-for="s in info?.srs ?? []" :key="s.uuid" class="flex items-center gap-2 px-1 py-0.5">
            <input v-model="checked" type="checkbox" :value="s.uuid" />{{ s.name }}<span class="text-neutral-600">({{ s.type }}, {{ s.free }} free)</span>
          </label>
        </div>
        <p v-if="warning" class="text-yellow-700">{{ warning }}</p>
        <div class="flex justify-end gap-2">
          <button :class="btn" :disabled="!info || !canSave || pending" @click="ok">OK</button>
          <button :class="btn" @click="emit('close')">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>
