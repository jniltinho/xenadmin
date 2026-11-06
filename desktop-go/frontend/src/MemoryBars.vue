<script setup lang="ts">
// The Memory tab's bars (MemoryView): a VM's memory in use inside its dynamic range (VMShinyBar), or each server's
// memory split into Xen and the control domain, each VM and free memory (HostShinyBar); read again on tick.
import { onMounted, ref, watch } from "vue";
import { api, type MemoryBars } from "./api/workspace";
import { serial } from "./serial";

const props = defineProps<{ server: string; kind: string; uuid: string; tick?: number }>();
const bars = ref<MemoryBars | null>(null);
const colours = ["#3C78B4", "#5A9BD5", "#7FB3E0", "#A5C9EA"];
const read = serial(async () => { try { bars.value = await api.memoryBars(props.server, props.kind, props.uuid); } catch { bars.value = null; } });
onMounted(read);
watch(() => props.tick, read);
const total = (segs: { bytes: number }[]) => Math.max(1, segs.reduce((a, s) => a + s.bytes, 0));
</script>

<template>
  <div v-if="bars" class="space-y-2 px-3 pt-3">
    <template v-if="kind === 'vms' || kind === 'templates'">
      <div class="relative h-[26px] w-[600px] border border-[#828790] bg-[#F2F2F2]" :title="bars.used_text">
        <div class="absolute inset-y-px bg-[#C9E7B0]" :style="{ left: `${bars.dynamic_min * 100}%`, width: `${Math.max(0, bars.dynamic_max - bars.dynamic_min) * 100}%` }" />
        <div class="absolute inset-y-px left-0 bg-[#4E9A06]" :style="{ width: `${bars.used * 100}%` }" />
      </div>
    </template>
    <template v-else>
      <div v-for="h in bars.hosts" :key="h.name">
        <p v-if="kind === 'pools'" class="font-semibold">{{ h.name }}</p>
        <div class="flex h-[26px] w-[600px] border border-[#828790]">
          <div v-for="(s, i) in h.segments" :key="i" :title="s.text" class="h-full"
               :style="{ width: `${Math.max(0.2, (s.bytes / total(h.segments)) * 100)}%`, background: i === 0 ? '#2E5E8A' : i === h.segments.length - 1 ? '#F2F2F2' : colours[(i - 1) % colours.length] }" />
        </div>
      </div>
    </template>
  </div>
</template>
