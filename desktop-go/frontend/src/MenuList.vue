<script setup lang="ts">
// A drop-down menu (the menu bar's and the context menu's): items, separators, and submenus whose items are loaded
// when the pointer reaches them, as the C# app fills "Migrate to Server" when it opens.
import { ref, watch } from "vue";

export interface MenuItem { label: string; icon?: string; enabled?: boolean; run?: () => void; sep?: boolean; sub?: () => Promise<MenuItem[]> }

const props = defineProps<{ items: MenuItem[]; hint?: string }>();
const open = ref(-1);
const children = ref<MenuItem[] | null>(null);
let generation = 0;
// another object's menu in the same place: its submenus must never run the previous object's items
watch(() => props.items, () => { open.value = -1; children.value = null; generation++; });

async function enter(i: number, it: MenuItem) {
  if (open.value === i) return;
  open.value = i;
  children.value = null;
  if (!it.sub || !it.enabled) return;
  const g = ++generation;
  const items = await it.sub().catch((e) => [{ label: String(e), enabled: false }]);
  if (g === generation) children.value = items;
}
</script>

<template>
  <div class="min-w-56 border border-neutral-400 bg-white py-1 shadow" @click.stop>
    <template v-for="(it, i) in items" :key="i">
      <div v-if="it.sep" class="my-1 border-t border-line"></div>
      <div v-else class="relative flex items-center gap-2 whitespace-nowrap px-2 py-0.5" :title="it.enabled ? '' : hint"
           :class="it.enabled ? 'hover:bg-select hover:text-white' : 'text-neutral-400'" @mouseenter="enter(i, it)"
           @click="it.enabled && !it.sub && it.run?.()">
        <img v-if="it.icon" :src="`/icons/${it.icon}.png`" alt="" class="h-4 w-4" :class="{ 'opacity-40 grayscale': !it.enabled }" /><span v-else class="w-4"></span>
        <span class="flex-1">{{ it.label }}</span><span v-if="it.sub" class="pl-4 text-xs">&#9656;</span>
        <div v-if="it.sub && open === i && it.enabled" class="absolute left-full top-0 z-30 text-black">
          <MenuList v-if="children" :items="children" />
          <div v-else class="min-w-40 border border-neutral-400 bg-white px-2 py-1 text-neutral-500 shadow">Loading...</div>
        </div>
      </div>
    </template>
  </div>
</template>
