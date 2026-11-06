<script setup lang="ts">
// The Performance tab (PerformancePage): the graphs of a VM or a server (the ones saved in the pool by either app, else
// the defaults), the Zoom menu, the lifecycle events drawn on the time axis; refreshed as the RRD step goes (5 s for
// 10 minutes). Actions has New, Edit and Delete Graph and Restore Default Graphs, Move Up/Down reorder the selected
// graph (also in the graphs' context menu); each change saves the layout at once (GraphList.SaveGraphs).
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { api, type DesignedGraph, type Performance } from "./api/workspace";
import GraphDetailsDialog from "./GraphDetailsDialog.vue";
import MenuList, { type MenuItem } from "./MenuList.vue";

// area: Tools > Options' "Area graph" (FillAreaUnderGraphs), else lines
const props = defineProps<{ server: string; uuid: string; kind: "vms" | "hosts"; area?: boolean }>();
const zooms = [["10m", "10 Minutes"], ["1h", "1 Hour"], ["1d", "1 Day"], ["1w", "1 Week"], ["1mo", "1 Month"], ["1y", "1 Year"]] as const;
const zoom = ref<string>("10m");
const data = ref<Performance | null>(null);
const error = ref("");
let timer = 0;
let seq = 0, loading = 0;
const W = 760, H = 150, L = 56, B = 18, R = 16; // R keeps the last time label inside

async function load() {
  const z = zoom.value, n = ++seq; // only the latest read is shown: an older one may end after a save's
  loading++;
  try {
    const p = await api.performance(props.server, props.kind, props.uuid, z);
    if (n === seq && z === zoom.value) { data.value = p; error.value = ""; if (sel.value >= p.graphs.length) sel.value = -1; }
  } catch (e) { if (n === seq) error.value = String(e); }
  finally { loading--; }
}
const poll = () => { if (!loading) load(); }; // a slow read is never overtaken by the next poll
function setZoom(z: string) { zoom.value = z; data.value = null; load(); }

// the graph editor (PerformancePage): the selected graph, the open dialog or question, one save at a time
const sel = ref(-1);
// the dialog's graph and the layout it was opened on (base: as saved, [] for the defaults)
const details = ref<{ title: string; graph: DesignedGraph; index: number; base: DesignedGraph[]; gs: DesignedGraph[] } | null>(null);
const ask = ref<{ title: string; text: string; run: () => void } | null>(null);
const menu = ref<{ x: number; y: number } | null>(null);
const actionsOpen = ref(false);
const pending = ref(false);
const graphs = computed(() => data.value?.graphs ?? []);
const layout = (): DesignedGraph[] => graphs.value.map((g) => ({ name: g.title, sources: [...g.sources] }));
// the colours shown, by data source (the memory totals have no line of their own)
const shown = computed(() => {
  const out: Record<string, string> = {};
  for (const g of graphs.value)
    g.sources.filter((src) => src !== "memory" && src !== "memory_total_kib").forEach((src, i) => { if (g.series[i]) out[src] = g.series[i].colour; });
  return out;
});
// GraphList.GetGraphsDataSources: the colours of the graphs' sources are saved with the layout
const coloursOf = (gs: DesignedGraph[]) =>
  Object.fromEntries(gs.flatMap((g) => g.sources).filter((src) => shown.value[src]).map((src) => [src, shown.value[src]]));

// the layout shown: base is what the server must still have for a change to be saved (SaveGraphs refuses otherwise)
const snapshot = () => ({ gs: layout(), base: data.value?.default ? [] : layout() });
// true when saved; else the error stays shown (after the refresh, which would clear it)
async function save(base: DesignedGraph[], gs: DesignedGraph[], colours: Record<string, string>, select: number) {
  if (pending.value) return false;
  pending.value = true;
  error.value = "";
  try {
    await api.saveGraphs(props.server, props.kind, props.uuid, base, gs, colours);
    sel.value = select;
    await load();
    return true;
  } catch (e) { await load(); error.value = String(e); return false; }
  finally { pending.value = false; }
}
const picked = computed(() => sel.value >= 0 && sel.value < graphs.value.length);
const canEdit = computed(() => picked.value && !data.value?.edit_blocker);
let opening = 0; // only the last New or Edit Graph opens the dialog, and never over an open one
async function open(index: number) {
  if (!data.value || pending.value || details.value) return;
  const n = ++opening, snap = snapshot(), g = index < 0 ? { name: data.value.new_name, sources: [] } : snap.gs[index];
  const t = await api.graphTexts(g.name);
  if (n === opening && !details.value) details.value = { title: index < 0 ? t.new_title : t.edit_title, graph: g, index, ...snap };
}
function newGraph() { if (data.value && !data.value.new_blocker) open(-1); }
function editGraph() { if (canEdit.value) open(sel.value); }
// the dialog stays open, with its draft, until the layout is saved
async function saveDetails(g: DesignedGraph, colours: Record<string, string>) {
  const d = details.value!, gs = d.gs.map((x) => ({ ...x }));
  if (d.index < 0) gs.push(g); else gs[d.index] = g;
  if (await save(d.base, gs, { ...coloursOf(gs), ...colours }, d.index < 0 ? gs.length - 1 : d.index)) details.value = null;
}
async function deleteGraph() {
  if (!picked.value || graphs.value.length < 2 || pending.value) return;
  const i = sel.value, { base, gs } = snapshot(), t = await api.graphTexts(gs[i].name);
  ask.value = { title: "Delete Graph", text: t.delete, run: () => {
    gs.splice(i, 1);
    save(base, gs, coloursOf(gs), Math.min(i, gs.length - 1)); // DeleteGraph: the next graph, the previous one for the last
  } };
}
async function restoreDefaults() {
  if (!data.value || data.value.default || pending.value) return;
  const { base } = snapshot(), t = await api.graphTexts("");
  ask.value = { title: "Restore Default Graphs", text: t.restore, run: () => save(base, [], {}, -1) };
}
function move(by: number) {
  const i = sel.value, j = i + by, { base, gs } = snapshot();
  if (!picked.value || j < 0 || j >= gs.length) return;
  [gs[i], gs[j]] = [gs[j], gs[i]];
  save(base, gs, coloursOf(gs), j);
}
function close() { actionsOpen.value = false; menu.value = null; }
const item = (label: string, enabled: boolean, run: () => void): MenuItem => ({ label, enabled: enabled && !pending.value, run: () => { close(); run(); } });
// PerformancePage.UpdateActionButtons and UpdateMoveButtons
const actions = computed<MenuItem[]>(() => [
  item("New Graph...", !!data.value && !data.value.new_blocker, newGraph),
  item("Edit Graph...", canEdit.value, editGraph),
  item("Delete Graph...", picked.value && graphs.value.length > 1, deleteGraph),
  item("Restore Default Graphs", !!data.value && !data.value.default, restoreDefaults),
]);
const contextItems = computed<MenuItem[]>(() => [...actions.value, { label: "", sep: true },
  item("Move Up", picked.value && sel.value > 0, () => move(-1)), item("Move Down", picked.value && sel.value < graphs.value.length - 1, () => move(1))]);
// the tooltip of New and Edit Graph on a VM that is not running
const hint = computed(() => data.value?.new_blocker ?? "");
function context(e: MouseEvent, gi: number) { sel.value = gi; actionsOpen.value = false; menu.value = { x: e.clientX, y: e.clientY }; }

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

onMounted(() => { load(); timer = window.setInterval(poll, 10000); window.addEventListener("click", close); });
onBeforeUnmount(() => { window.clearInterval(timer); window.removeEventListener("click", close); });
</script>

<template>
  <div class="min-h-0 flex-1 overflow-auto">
    <div class="bg-gradient-to-r from-band-dark to-band px-3 py-1.5 font-semibold text-white">Performance Graphs</div>
    <div class="flex items-center gap-2 border-b border-line bg-pane px-3 py-1.5">
      <span>Zoom:</span>
      <button v-for="[z, label] in zooms" :key="z" class="btn" :class="{ 'bg-select text-white hover:bg-select': zoom === z }" @click="setZoom(z)">{{ label }}</button>
      <div class="relative ml-4">
        <button class="btn" @click.stop="menu = null; actionsOpen = !actionsOpen">Actions &#9662;</button>
        <MenuList v-if="actionsOpen" class="absolute left-0 top-full z-20" :items="actions" :hint="hint" />
      </div>
      <button class="btn" :disabled="!picked || sel === 0 || pending" @click="move(-1)">Move Up</button>
      <button class="btn" :disabled="!picked || sel === graphs.length - 1 || pending" @click="move(1)">Move Down</button>
      <span v-if="error" class="text-red-700">{{ error }}</span>
    </div>
    <p v-if="!data && !error" class="p-3 text-neutral-500">Loading...</p>
    <div v-for="(g, gi) in graphs" :key="gi" class="mx-3 my-3 border" :class="sel === gi ? 'border-select bg-sky-100' : 'border-line'"
         @click="sel = gi" @dblclick="sel = gi; editGraph()" @contextmenu.prevent.stop="context($event, gi)">
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
          <g v-for="(s, si) in g.series" :key="si">
            <template v-for="(pts, pi) in paths(s.values, g.max)" :key="pi">
              <polygon v-if="area" :points="`${pts.split(' ')[0].split(',')[0]},${H - B} ${pts} ${pts.split(' ').at(-1)!.split(',')[0]},${H - B}`"
                       :fill="s.colour" fill-opacity="0.35" stroke="none" />
              <polyline :points="pts" fill="none" :stroke="s.colour" stroke-width="1.5" vector-effect="non-scaling-stroke" />
            </template>
          </g>
        </svg>
        <ul class="w-56 shrink-0 space-y-0.5 text-xs">
          <li v-for="(s, si) in g.series" :key="si" class="flex items-center gap-1.5">
            <span class="inline-block h-0.5 w-4" :style="{ background: s.colour }"></span>{{ s.name }}</li>
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
    <MenuList v-if="menu" class="fixed z-20" :style="{ left: `${menu.x}px`, top: `${menu.y}px` }" :items="contextItems" :hint="hint" />
    <GraphDetailsDialog v-if="details" :server="server" :kind="kind" :uuid="uuid" :title="details.title" :graph="details.graph" :colours="shown"
                        :failed="error" :busy="pending" @close="details = null" @save="saveDetails" />
    <div v-if="ask" class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
      <div class="w-[28rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ ask.title }}</div>
        <div class="space-y-3 p-3"><p class="whitespace-pre-line">{{ ask.text }}</p>
          <div class="flex justify-end gap-2">
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="pending" @click="ask.run(); ask = null">Yes</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="ask = null">No</button>
          </div></div>
      </div>
    </div>
  </div>
</template>
