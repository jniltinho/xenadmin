<script setup lang="ts">
// The Search tab's overview (SearchView / QueryPanel): the objects with their depth (pool, its servers, their VMs), a
// +/- to fold a node, and the live columns of OverviewColumns (CPU and memory with a bar, disks, network, address,
// uptime), read again every ten seconds while shown (MetricUpdater).
// shortcut: no New/Edit Search, Saved Searches, Export/Import or Show Columns buttons yet.
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from "vue";
import { api, type OverviewRow, type XenObject } from "./api/workspace";

export type SearchRow = { o: XenObject; type: string; key: string; icon: string; id?: string; depth?: number };
const props = defineProps<{ title: string; rows: SearchRow[] }>();
const emit = defineEmits<{ (e: "open", r: SearchRow): void }>();

const columns = ["CPU Usage", "Used Memory", "Disks\n(avg / max KBs)", "Network\n(avg / max KBs)", "Address", "Uptime"];
const live = reactive<Record<string, Record<string, OverviewRow>>>({}); // server key -> uuid -> columns
const folded = ref(new Set<string>());

const keys = computed(() => [...new Set(props.rows.map((r) => r.key))]);
let timer: ReturnType<typeof setInterval> | undefined;
let reading = false;
async function read() {
  if (reading) return;
  reading = true;
  try {
    await Promise.all(keys.value.map(async (k) => { try { live[k] = (await api.overview(k)) ?? {}; } catch { live[k] = {}; } }));
  } finally { reading = false; }
}
onMounted(() => { read(); timer = setInterval(read, 10000); });
onBeforeUnmount(() => clearInterval(timer));
watch(keys, read);

// the rows shown: those under a folded node are hidden; a node has children when the next row is deeper
const shown = computed(() => {
  const out: (SearchRow & { children: boolean; uid: string })[] = [];
  let skip = -1;
  props.rows.forEach((r, i) => {
    const depth = r.depth ?? 0;
    if (skip >= 0 && depth > skip) return;
    skip = -1;
    const uid = `${r.key}/${r.o.uuid}`, children = (props.rows[i + 1]?.depth ?? 0) > depth;
    if (children && folded.value.has(uid)) skip = depth;
    out.push({ ...r, children, uid });
  });
  return out;
});
function fold(uid: string) {
  const s = new Set(folded.value);
  if (!s.delete(uid)) s.add(uid);
  folded.value = s;
}
const none = (r: SearchRow): OverviewRow => ({ cpu: "-", cpu_bar: -1, memory: "-", memory_bar: -1, disks: "-", network: "-",
  address: r.type === "Pool" ? "-" : "", uptime: r.type === "Pool" ? "-" : "" });
const cells = (r: SearchRow) => {
  const v = live[r.key]?.[r.o.uuid] ?? none(r);
  return [[v.cpu, v.cpu_bar], [v.memory, v.memory_bar], [v.disks, -1], [v.network, -1], [v.address, -1], [v.uptime, -1]] as [string, number][];
};
</script>

<template>
  <div class="min-h-0 flex-1 overflow-auto">
    <div class="bg-gradient-to-r from-band-dark to-band px-3 py-1.5 font-semibold text-white">{{ title }}</div>
    <table class="m-3 border-collapse text-sm" style="table-layout: fixed">
      <colgroup><col style="width: 250px" /><col style="width: 115px" /><col style="width: 125px" /><col style="width: 120px" />
        <col style="width: 120px" /><col style="width: 120px" /><col style="width: 170px" /></colgroup>
      <thead class="sticky top-0 bg-white">
        <tr><th class="border-b border-r border-[#E3E3E3] px-1 py-1 font-semibold">Name</th>
          <th v-for="c in columns" :key="c" class="border-b border-r border-[#E3E3E3] px-1 py-1 font-semibold whitespace-pre-line">{{ c }}</th></tr>
      </thead>
      <tbody>
        <tr v-for="r in shown" :key="r.uid" class="cursor-default hover:bg-neutral-100" @dblclick="emit('open', r)">
          <td class="py-1">
            <span class="flex items-center gap-1" :style="{ paddingLeft: `${4 + (r.depth ?? 0) * 25}px` }">
              <img v-if="r.children" :src="`/icons/${folded.has(r.uid) ? 'tree_plus' : 'tree_minus'}.png`" alt="" class="mr-0.5 h-[9px] w-[9px] cursor-pointer" @click="fold(r.uid)" />
              <span v-else class="mr-0.5 inline-block w-[9px]" />
              <img :src="`/icons/${r.icon}.png`" alt="" class="h-4 w-4" />
              <span class="min-w-0"><span class="block truncate">{{ r.o.name }}</span>
                <span v-if="r.o.description" class="block truncate text-xs text-neutral-500">{{ r.o.description }}</span></span>
            </span>
          </td>
          <td v-for="([text, bar], i) in cells(r)" :key="i" class="px-1 py-1 text-center">
            <div v-if="bar >= 0" class="mx-auto mb-0.5 h-2 w-24 border border-neutral-400 bg-white"><div class="h-full bg-[#4A9E2B]" :style="{ width: `${Math.min(100, bar)}%` }" /></div>
            <span class="block truncate">{{ text }}</span>
          </td>
        </tr>
      </tbody>
    </table>
    <p v-if="!rows.length" class="p-3 text-neutral-500">No connected server. Use Add New Server, or double-click a server of the tree to connect.</p>
  </div>
</template>
