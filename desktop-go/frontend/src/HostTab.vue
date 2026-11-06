<script setup lang="ts">
// The Memory, Storage, Networking and NICs tabs of a server or pool, and the Storage tab of an SR (Details' tables).
// shortcut: read-only; their buttons (New Network, bonds, disks of an SR) come later.
import { onMounted, ref, watch } from "vue";
import { api, type Table } from "./api/workspace";
import DetailsTable from "./DetailsTable.vue";
import { serial } from "./serial";

// tick: the server's objects changed (xcp:changed), read the tab again
const props = defineProps<{ server: string; kind: string; uuid: string; tab: string; tick?: number }>();
const emit = defineEmits<{ (e: "status", text: string, error: boolean): void }>();

const table = ref<Table | null>(null);
const picked = ref("");
const load = serial(async () => {
  try { table.value = await api.hostTab(props.server, props.kind, props.uuid, props.tab); } catch (e) { emit("status", String(e), true); }
});
onMounted(load);
watch(() => props.tick, load);
</script>

<template>
  <div class="flex min-h-0 flex-1 flex-col overflow-auto">
    <DetailsTable v-if="table" :table="table" :selected="picked" :plain="tab === 'memory'" @select="picked = $event" />
  </div>
</template>
