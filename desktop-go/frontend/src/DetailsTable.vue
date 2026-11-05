<script setup lang="ts">
// One of Details' tables (a band with its title, then a grid); a row can be selected by its key.
import type { Table } from "./api/workspace";

defineProps<{ table: Table; selected?: string; plain?: boolean }>();
const emit = defineEmits<{ (e: "select", key: string): void }>();
</script>

<template>
  <div class="bg-gradient-to-r from-band-dark to-band px-3 py-1.5 font-semibold text-white">{{ table.title }}</div>
  <div class="overflow-auto p-3">
    <table class="w-full border-collapse border border-line">
      <thead v-if="!plain" class="bg-pane text-left">
        <tr><th v-for="h in table.headers" :key="h" class="border-b border-line px-2 py-1 font-normal whitespace-nowrap">{{ h }}</th></tr>
      </thead>
      <tbody>
        <tr v-for="(r, i) in table.rows" :key="table.keys[i] ?? i" class="cursor-default"
            :class="selected && table.keys[i] === selected ? 'bg-select text-white' : 'hover:bg-neutral-100'" @click="table.keys[i] && emit('select', table.keys[i])">
          <td v-for="(c, j) in r" :key="j" class="border-b border-line px-2 py-1 select-text" :class="{ 'w-64 text-neutral-600': plain && j === 0 }">{{ c }}</td>
        </tr>
      </tbody>
    </table>
    <p v-if="!table.rows.length" class="pt-2 text-neutral-500">None.</p>
  </div>
</template>
