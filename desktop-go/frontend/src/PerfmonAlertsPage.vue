<script setup lang="ts">
// The Alerts page of a VM's, server's or SR's Properties (PerfmonAlertEditPage): an alert per resource with its
// threshold and duration, and the repeat interval shared by every alert. The dialog's OK calls changes().
import { onMounted, reactive, ref } from "vue";
import { api, type PerfmonPage } from "./api/workspace";

const props = defineProps<{ server: string; kind: string; uuid: string; disabled?: boolean }>();

// the WinForms page's labels (PerfmonAlertEditPage.resx)
const labels: Record<string, { box: string; level: string; unit: string }> = {
  cpu_usage: { box: "Generate CPU usage alerts", level: "When CPU usage exceeds:", unit: "%" },
  network_usage: { box: "Generate network usage alerts", level: "When network usage exceeds:", unit: "KB/s" },
  disk_usage: { box: "Generate disk usage alerts", level: "When disk usage exceeds:", unit: "KB/s" },
  memory_free_kib: { box: "Generate memory usage alerts", level: "When free memory falls below:", unit: "MB" },
  sr_io_throughput_total_per_host: { box: "Generate storage throughput alerts", level: "When SR throughput for any host exceeds:", unit: "KB/s" },
  mem_usage: { box: "Generate control domain memory usage alerts", level: "When control domain memory usage exceeds:", unit: "%" },
};

const orig = ref<PerfmonPage | null>(null);
const f = reactive({ alerts: [] as PerfmonPage["alerts"], interval: 60 });
const error = ref("");
onMounted(async () => {
  try {
    const p = (orig.value = await api.perfmon(props.server, props.kind, props.uuid));
    f.alerts = p.alerts.map((a) => ({ ...a }));
    f.interval = p.interval;
  } catch (e) { error.value = String(e); }
});
const anyOn = () => f.alerts.some((a) => a.enabled);

// null: nothing changed (AlertGroup.HasChanged); a string: the reason it cannot be saved; else what to save
function changes(): null | string | { on: { name: string; level: number; period: number }[]; interval: number } {
  const o = orig.value;
  if (!o) return null;
  const same = o.interval === f.interval && o.alerts.every((a, i) => {
    const b = f.alerts[i];
    return a.enabled === b.enabled && (!b.enabled || a.level === b.level && a.period === b.period);
  });
  if (same) return null;
  if (!Number.isInteger(f.interval) || f.interval < 5 || f.interval > 86400 || f.interval % 5 !== 0) return "Interval must be a multiple of 5";
  for (const a of f.alerts.filter((x) => x.enabled)) {
    if (!(a.level >= a.min && a.level <= a.max)) return `${labels[a.name]?.level ?? a.name} enter a value from ${a.min} to ${a.max}.`;
    if (!(a.period >= 1 && a.period <= 60)) return "Enter a duration from 1 to 60 minutes.";
  }
  return { on: f.alerts.filter((a) => a.enabled).map((a) => ({ name: a.name, level: a.level, period: a.period })), interval: f.interval };
}
defineExpose({ changes });
</script>

<template>
  <div>
    <p class="mb-3">The server can notify you when your resource usage crosses certain thresholds. These notifications are generated as system alerts and can be forwarded to an email server the same as any other System Alert.</p>
    <p v-if="error" class="text-red-600">{{ error }}</p>
    <fieldset v-if="orig" :disabled="disabled" class="space-y-3">
      <div v-for="a in f.alerts" :key="a.name" class="border border-line p-2">
        <label class="flex items-center gap-2 font-semibold"><input v-model="a.enabled" type="checkbox" />{{ labels[a.name]?.box ?? a.name }}</label>
        <div class="mt-2 grid grid-cols-[17rem_7rem_auto] items-center gap-x-2 gap-y-1.5 pl-6" :class="a.enabled ? '' : 'text-neutral-400'">
          <label :for="`pm-${a.name}`">{{ labels[a.name]?.level }}</label>
          <input :id="`pm-${a.name}`" v-model.number="a.level" type="number" :min="a.min" :max="a.max" :step="a.name === 'cpu_usage' || a.name === 'mem_usage' ? 5 : 1"
                 :disabled="!a.enabled" class="border border-line px-1 py-0.5" />
          <span>{{ labels[a.name]?.unit }}</span>
          <label :for="`pmp-${a.name}`">For longer than:</label>
          <input :id="`pmp-${a.name}`" v-model.number="a.period" type="number" min="1" max="60" :disabled="!a.enabled" class="border border-line px-1 py-0.5" />
          <span>minutes</span>
        </div>
      </div>
      <div class="flex items-center gap-2" :class="anyOn() ? '' : 'text-neutral-400'">
        <label for="pm-interval">Alert repeat interval:</label>
        <input id="pm-interval" v-model.number="f.interval" type="number" min="5" max="86400" step="5" :disabled="!anyOn()" class="w-28 border border-line px-1 py-0.5" />
        <span>minutes</span>
      </div>
    </fieldset>
  </div>
</template>
