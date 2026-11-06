<script setup lang="ts">
// The General tab (Details.General shown by XenAdmin.Desktop's Views.Sections, after the WinForms PDSection): the band,
// Properties and Expand all / Collapse all, then the collapsible sections with a selectable row; read again on each
// tick (xcp:changed). A section keeps its open state across reads.
import { onMounted, reactive, ref, watch } from "vue";
import { api } from "./api/workspace";
import { serial } from "./serial";

const props = defineProps<{ server: string; kind: string; uuid: string; tick?: number }>();
const emit = defineEmits<{ (e: "properties"): void; (e: "status", text: string, error: boolean): void }>();

type Section = Awaited<ReturnType<typeof api.general>>[number];
const sections = ref<Section[]>([]);
const open = reactive(new Map<string, boolean>());
const picked = ref(""); // "title/row" of the selected row; the first row starts selected
// Details.TabTitle
const title = ({ vms: "VM", templates: "VM", hosts: "Server", pools: "Pool", srs: "Storage" } as Record<string, string>)[props.kind];

const load = serial(async () => {
  try {
    sections.value = await api.general(props.server, props.kind, props.uuid);
    for (const s of sections.value) if (!open.has(s.title)) open.set(s.title, s.expanded);
    if (!picked.value && sections.value.length) picked.value = `${sections.value[0].title}/0`;
  } catch (e) { emit("status", String(e), true); }
});
onMounted(load);
watch(() => props.tick, load);

function all(on: boolean) { for (const s of sections.value) open.set(s.title, on); }
</script>

<template>
  <div class="bg-gradient-to-r from-band-dark to-band px-3 py-1.5 font-semibold text-white">{{ title ? `${title} General Properties` : "General Properties" }}</div>
  <div class="min-h-0 flex-1 overflow-auto">
    <div class="flex max-w-[924px] flex-col gap-2.5 px-3 pb-3 pt-2.5">
      <div class="flex items-center">
        <button class="btn" @click="emit('properties')">Properties</button>
        <span class="flex-1"></span>
        <a class="mr-4 cursor-pointer text-blue-700 underline" @click="all(true)">Expand all</a>
        <a class="cursor-pointer text-blue-700 underline" @click="all(false)">Collapse all</a>
      </div>
      <section v-for="s in sections" :key="s.title" class="border border-[#dcdcdc] bg-white">
        <div class="flex h-8 cursor-pointer items-center bg-pane pl-2 pr-1.5" @click="open.set(s.title, !open.get(s.title))">
          <span class="flex-1 text-[15px]">{{ s.title }}</span>
          <span class="mr-2 h-5 w-px bg-[#dcdcdc]"></span>
          <img :src="`/icons/PDChevron${open.get(s.title) ? 'Up' : 'Down'}.png`" alt="" class="h-5 w-5" />
        </div>
        <div v-if="open.get(s.title)">
          <div v-for="([k, v], i) in s.rows" :key="i" class="grid grid-cols-[150px_1fr]"
               :class="picked === `${s.title}/${i}` ? 'bg-select text-white' : ''" @click="picked = `${s.title}/${i}`">
            <span class="break-words px-2 py-1.5">{{ k }}:</span>
            <span class="whitespace-pre-line break-words py-1.5 pr-2 select-text">{{ v }}</span>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>
