<script setup lang="ts">
// Graph Details (GraphDetailsDialog): the graph's name and the data sources it shows, with their type, recording state
// and colour; "Show data sources" adds the hidden and the disabled ones, Enable turns recording on for one. Save hands
// the graph and the colours back (the Performance tab saves the layout).
import { computed, onMounted, ref, watch } from "vue";
import { api, type DataSource, type DesignedGraph } from "./api/workspace";

// failed: why the last Save did not save; busy: a Save runs
const props = defineProps<{ server: string; kind: string; uuid: string; title: string; graph: DesignedGraph; colours: Record<string, string>;
  failed?: string; busy?: boolean }>();
const emit = defineEmits<{ (e: "close"): void; (e: "save", g: DesignedGraph, colours: Record<string, string>): void }>();
// Palette.DefaultColors
const palette = ["#0066ff", "#339900", "#ff9900", "#660099", "#990033", "#995e00", "#ffff00", "#00ffb4", "#009982", "#bfa4a6"];
const name = ref(props.graph.name);
const chosen = ref<string[]>([...props.graph.sources]); // DesignedGraph.DataSourceItems: kept in order, new ones last
const sources = ref<DataSource[] | null>(null);
const colour = ref<Record<string, string>>({});
const changed = new Set<string>();
const error = ref("");
const search = ref("");
const showHidden = ref(false), showDisabled = ref(false), showMenu = ref(false);
const picked = ref("");
const enabling = ref(false);

async function read() {
  try {
    const list = await api.graphDataSources(props.server, props.kind, props.uuid);
    list.forEach((d, i) => { colour.value[d.name] ??= d.colour || props.colours[d.name] || palette[i % palette.length]; });
    sources.value = list;
  } catch (e) { error.value = String(e); }
}
onMounted(read);

// PopulateDataGridView: the checked rows first, then by name (natural order), sorted again when the filter changes
const collator = new Intl.Collator(undefined, { numeric: true });
const first = ref(new Set(props.graph.sources));
watch([search, showHidden, showDisabled], () => { first.value = new Set(chosen.value); });
const rows = computed(() => (sources.value ?? [])
  .filter((d) => (showHidden.value || !d.hidden) && (showDisabled.value || d.enabled) && d.text.toLowerCase().includes(search.value.trim().toLowerCase()))
  .sort((a, b) => Number(first.value.has(b.name)) - Number(first.value.has(a.name)) || collator.compare(a.text, b.text)));
const checkedRows = computed(() => rows.value.filter((d) => chosen.value.includes(d.name)).length);
const current = computed(() => rows.value.find((d) => d.name === picked.value));

function toggle(d: DataSource) {
  chosen.value = chosen.value.includes(d.name) ? chosen.value.filter((n) => n !== d.name) : [...chosen.value, d.name];
}
function clearAll() { const shown = new Set(rows.value.map((d) => d.name)); chosen.value = chosen.value.filter((n) => !shown.has(n)); }
function setColour(d: DataSource, c: string) { colour.value[d.name] = c; changed.add(d.name); }

async function enable() {
  const d = current.value;
  if (!d || d.enabled || enabling.value) return;
  enabling.value = true;
  error.value = "";
  try { await api.enableDataSource(props.server, props.kind, props.uuid, d.name, d.text); await read(); }
  catch (e) { error.value = String(e); }
  finally { enabling.value = false; }
}

// SaveButton_Click: the colours of the graph's sources (GraphList.UpdateDataSources keeps those without one) and the
// ones changed here
function save() {
  if (!checkedRows.value || props.busy) return;
  const out: Record<string, string> = {};
  for (const n of chosen.value) if (colour.value[n]) out[n] = colour.value[n];
  for (const n of changed) out[n] = colour.value[n];
  emit("save", { name: name.value, sources: chosen.value }, out);
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20" @click="showMenu = false">
    <form class="flex h-[32rem] w-[48rem] flex-col border border-neutral-500 bg-white shadow-lg" @submit.prevent="save">
      <div class="truncate bg-band-dark px-3 py-1 font-semibold text-white">{{ title }}</div>
      <div class="flex min-h-0 flex-1 flex-col gap-2 p-3">
        <p>Specify a name for the graph and select the data sources you want to include in it.</p>
        <label class="flex items-center gap-2">Name:<input v-model="name" class="flex-1 border border-neutral-400 px-1" autofocus /></label>
        <div class="flex items-center gap-2">
          <input v-model="search" class="w-56 border border-neutral-400 px-1" placeholder="Search" :disabled="!sources" aria-label="Search" />
          <div class="relative">
            <button type="button" class="btn" @click.stop="showMenu = !showMenu">Show data sources &#9662;</button>
            <div v-if="showMenu" class="absolute left-0 top-full z-10 min-w-40 border border-neutral-400 bg-white py-1 shadow" @click.stop>
              <label class="flex items-center gap-2 px-2 py-0.5 hover:bg-neutral-100"><input v-model="showHidden" type="checkbox" />Hidden</label>
              <label class="flex items-center gap-2 px-2 py-0.5 hover:bg-neutral-100"><input v-model="showDisabled" type="checkbox" />Disabled</label>
            </div>
          </div>
        </div>
        <p v-if="!sources && !error" class="text-neutral-500">Retrieving data sources...</p>
        <div v-else class="min-h-0 flex-1 overflow-auto border border-line">
          <table class="w-full text-left text-xs">
            <thead class="sticky top-0 bg-pane"><tr><th class="w-6"></th><th class="px-1">Data Source</th><th class="px-1">Description</th>
              <th class="px-1">Type</th><th class="px-1">Enabled</th><th class="w-16 px-1">Color</th></tr></thead>
            <tbody>
              <tr v-for="d in rows" :key="d.name" :class="[picked === d.name ? 'bg-select text-white' : '', d.hidden ? 'text-neutral-400' : '']" @click="picked = d.name">
                <td class="px-1"><input type="checkbox" :checked="chosen.includes(d.name)" :aria-label="d.text" @change="toggle(d)" /></td>
                <td class="px-1">{{ d.text }}</td><td class="px-1">{{ d.description }}</td><td class="px-1">{{ d.type }}</td>
                <td class="px-1">{{ d.enabled ? "Yes" : "No" }}</td>
                <td class="px-1"><input type="color" class="h-4 w-8" :value="colour[d.name]" :aria-label="`${d.text} color`" @input="setColour(d, ($event.target as HTMLInputElement).value)" /></td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="flex gap-2">
          <button type="button" class="btn" :disabled="!checkedRows" @click="clearAll">Clear All</button>
          <button type="button" class="btn" :disabled="!current || current.enabled || enabling" @click="enable">Enable</button>
        </div>
        <p v-if="error || failed" class="whitespace-pre-line text-red-600">{{ error || failed }}</p>
      </div>
      <div class="flex justify-end gap-2 border-t border-line p-2">
        <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="!checkedRows || busy">Save</button>
        <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Cancel</button>
      </div>
    </form>
  </div>
</template>
