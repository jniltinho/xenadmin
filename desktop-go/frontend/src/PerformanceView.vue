<script setup lang="ts">
// The Performance tab (PerformancePage): the default graphs of a VM or a server, the Zoom menu, the lifecycle events
// drawn on the time axis; refreshed as the RRD step goes (5 s for 10 minutes).
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { api, type Performance } from "./api/workspace";

// area: Tools > Options' "Area graph" (FillAreaUnderGraphs), else lines
const props = defineProps<{ server: string; uuid: string; kind: "vms" | "hosts"; area?: boolean }>();
const zooms = [["10m", "10 Minutes"], ["1h", "1 Hour"], ["1d", "1 Day"], ["1w", "1 Week"], ["1mo", "1 Month"], ["1y", "1 Year"]] as const;
const zoom = ref<string>("10m");
const data = ref<Performance | null>(null);
const error = ref("");
let timer = 0;
// DataPlot's palette (the first colours of Palette.GetColour)
const colours = ["#1f77b4", "#d62728", "#2ca02c", "#ff7f0e", "#9467bd", "#8c564b", "#e377c2", "#17becf", "#bcbd22", "#7f7f7f"];
const W = 760, H = 150, L = 56, B = 18, R = 16; // R keeps the last time label inside

async function load() {
  const z = zoom.value;
  try {
    const p = await api.performance(props.server, props.kind, props.uuid, z);
    if (z === zoom.value) { data.value = p; error.value = ""; }
  } catch (e) { error.value = String(e); }
}
function setZoom(z: string) { zoom.value = z; data.value = null; load(); }

const x = (t: number) => L + ((t - data.value!.start) / Math.max(1, data.value!.end - data.value!.start)) * (W - L - R);
const y = (v: number, max: number) => H - B - (Math.min(v, max) / (max || 1)) * (H - B - 6);
// one polyline per run of values (a null is a gap)
function paths(values: (number | null)[], max: number) {
  const out: string[] = [];
  let cur: string[] = [];
  values.forEach((v, i) => {
    if (v === null) { if (cur.length) out.push(cur.join(" ")); cur = []; return; }
    cur.push(`${x(data.value!.times[i]).toFixed(1)},${y(v, max).toFixed(1)}`);
  });
  if (cur.length) out.push(cur.join(" "));
  return out;
}
const ticks = computed(() => {
  if (!data.value) return [];
  const { start, end } = data.value, n = 6;
  return Array.from({ length: n + 1 }, (_, i) => start + ((end - start) * i) / n);
});
const timeText = (t: number) => {
  const d = new Date(t * 1000);
  return zoom.value === "10m" || zoom.value === "1h" || zoom.value === "1d"
    ? d.toLocaleTimeString(undefined, { hour: "2-digit", minute: "2-digit" }) : d.toLocaleDateString(undefined, { day: "2-digit", month: "short" });
};
const eventIcon: Record<string, string> = { VM_STARTED: "001_PowerOn_h32bit_16", VM_SHUTDOWN: "001_ShutDown_h32bit_16", VM_REBOOTED: "001_Reboot_h32bit_16",
  VM_SUSPENDED: "000_paused_h32bit_16", VM_RESUMED: "000_Resumed_h32bit_16" };

onMounted(() => { load(); timer = window.setInterval(load, 10000); });
onBeforeUnmount(() => window.clearInterval(timer));
</script>

<template>
  <div class="min-h-0 flex-1 overflow-auto">
    <div class="bg-gradient-to-r from-band-dark to-band px-3 py-1.5 font-semibold text-white">Performance Graphs</div>
    <div class="flex items-center gap-2 border-b border-line bg-pane px-3 py-1.5">
      <span>Zoom:</span>
      <button v-for="[z, label] in zooms" :key="z" class="btn" :class="{ 'bg-select text-white hover:bg-select': zoom === z }" @click="setZoom(z)">{{ label }}</button>
      <span v-if="error" class="text-red-700">{{ error }}</span>
    </div>
    <p v-if="!data && !error" class="p-3 text-neutral-500">Loading...</p>
    <div v-for="(g, gi) in data?.graphs ?? []" :key="g.title" class="mx-3 my-3 border border-line">
      <div class="border-b border-line bg-pane px-2 py-1 font-semibold">{{ g.title }}</div>
      <div class="flex gap-3 p-2">
        <svg :viewBox="`0 0 ${W} ${H}`" class="h-40 min-w-0 flex-1" preserveAspectRatio="none" role="img" :aria-label="g.title">
          <g v-for="i in 5" :key="i">
            <line :x1="L" :x2="W - R" :y1="y((g.max * (i - 1)) / 4, g.max)" :y2="y((g.max * (i - 1)) / 4, g.max)" stroke="#e5e5e5" />
            <text :x="L - 4" :y="y((g.max * (i - 1)) / 4, g.max) + 3" text-anchor="end" font-size="10" fill="#555">
              {{ +((g.max * (i - 1)) / 4 / g.divisor).toFixed(1) }} {{ g.label }}</text>
          </g>
          <g v-for="t in ticks" :key="t">
            <line :x1="x(t)" :x2="x(t)" :y1="6" :y2="H - B" stroke="#f0f0f0" />
            <text :x="x(t)" :y="H - 4" text-anchor="middle" font-size="10" fill="#555">{{ timeText(t) }}</text>
          </g>
          <g v-if="gi === 0">
            <line v-for="ev in data!.events" :key="ev.time + ev.name" :x1="x(Date.parse(ev.time) / 1000)" :x2="x(Date.parse(ev.time) / 1000)" :y1="6" :y2="H - B" stroke="#999" stroke-dasharray="3 3" />
          </g>
          <g v-for="(s, si) in g.series" :key="s.name">
            <template v-for="(pts, pi) in paths(s.values, g.max)" :key="pi">
              <polygon v-if="area" :points="`${pts.split(' ')[0].split(',')[0]},${H - B} ${pts} ${pts.split(' ').at(-1)!.split(',')[0]},${H - B}`"
                       :fill="colours[si % colours.length]" fill-opacity="0.35" stroke="none" />
              <polyline :points="pts" fill="none" :stroke="colours[si % colours.length]" stroke-width="1.5" vector-effect="non-scaling-stroke" />
            </template>
          </g>
        </svg>
        <ul class="w-56 shrink-0 space-y-0.5 text-xs">
          <li v-for="(s, si) in g.series" :key="s.name" class="flex items-center gap-1.5">
            <span class="inline-block h-0.5 w-4" :style="{ background: colours[si % colours.length] }"></span>{{ s.name }}</li>
          <li v-if="!g.series.length" class="text-neutral-500">No data sources.</li>
        </ul>
      </div>
    </div>
    <div v-if="data?.events.length" class="mx-3 mb-3">
      <div class="mb-1 font-semibold">Events</div>
      <div v-for="ev in data.events" :key="ev.time + ev.name" class="flex items-center gap-2 text-xs">
        <img :src="`/icons/${eventIcon[ev.name] ?? 'alert6_16'}.png`" alt="" class="h-4 w-4" />{{ new Date(ev.time).toLocaleString() }} — {{ ev.text || ev.name }}
      </div>
    </div>
  </div>
</template>
