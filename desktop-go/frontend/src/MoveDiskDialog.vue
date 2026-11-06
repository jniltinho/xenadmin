<script setup lang="ts">
// Move Virtual Disk (Dialogs.MoveDisk: MoveVirtualDiskDialog / MigrateVirtualDiskDialog): the SRs with their free space
// or the reason they cannot take the disk (greyed, not selectable), Rescan (SrPicker.ScanSRs), and Move.
import { onMounted, ref } from "vue";
import { api, type DiskTargets } from "./api/workspace";

const props = defineProps<{ server: string; vdi: string }>();
const emit = defineEmits<{ (e: "close"): void; (e: "move", sr: string, migrate: boolean): void }>();

const targets = ref<DiskTargets | null>(null);
const sr = ref("");
const error = ref("");
const pending = ref(false);
async function fill() {
  targets.value = await api.moveDiskTargets(props.server, props.vdi);
  if (!targets.value.srs.some((s) => s.uuid === sr.value && !s.reason)) sr.value = targets.value.srs.find((s) => !s.reason)?.uuid ?? "";
}
onMounted(() => fill().catch((e) => (error.value = String(e))));
async function rescan() {
  if (pending.value || !targets.value) return;
  pending.value = true;
  error.value = "";
  try { await api.rescanSrs(props.server, targets.value.srs.map((s) => s.uuid)); } catch (e) { error.value = String(e); }
  try { await fill(); } catch (e) { error.value ||= String(e); }
  pending.value = false;
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <div class="w-[35rem] border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">Move Virtual Disk</div>
      <div class="space-y-2 p-3">
        <p>Select the storage repository that you would like to move the virtual disk to:</p>
        <div class="h-56 overflow-auto border border-line">
          <div v-for="s in targets?.srs ?? []" :key="s.uuid" class="cursor-default whitespace-pre px-1 py-0.5"
               :class="s.reason ? 'text-neutral-400' : sr === s.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100'"
               @click="!s.reason && (sr = s.uuid)">{{ s.text }}</div>
        </div>
        <p v-if="error" class="text-red-600">{{ error }}</p>
        <div class="flex gap-2 pt-1">
          <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="pending || !targets" @click="rescan">Rescan</button>
          <span class="flex-1" />
          <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="pending || !sr"
                  @click="emit('move', sr, !!targets?.migrate)">Move</button>
          <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>
