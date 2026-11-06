<script setup lang="ts">
// Add Virtual Disk (Dialogs.NewDisk, NewDiskDialog): name, description, size in GB or MB and the SR, for a VM's
// Storage tab or an SR's (the SR preselected).
import { reactive } from "vue";
import type { NewDisk } from "./api/workspace";

const props = defineProps<{ options: NewDisk }>();
const emit = defineEmits<{ (e: "close"): void; (e: "add", name: string, description: string, sr: string, bytes: number): void }>();

const f = reactive({ name: props.options.name, description: "", size: 10, unit: "GB" as "GB" | "MB",
  sr: props.options.default_sr || props.options.srs[0]?.uuid || "", error: "" });
function submit() {
  const bytes = Math.round(f.size * (f.unit === "GB" ? 1024 ** 3 : 1024 ** 2));
  if (!f.name.trim()) { f.error = "Enter a name for the disk."; return; }
  if (!f.sr) { f.error = "Select a storage repository."; return; }
  if (!(bytes > 0)) { f.error = "Enter a disk size."; return; }
  emit("add", f.name.trim(), f.description, f.sr, bytes);
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <form novalidate class="w-[36rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="submit">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">Add Virtual Disk</div>
      <div class="space-y-2 p-3">
        <p>Enter a name, description and size for the new virtual disk, and select the storage repository where it will be stored.</p>
        <label class="grid grid-cols-[6rem_1fr] items-center gap-2">Name:<input v-model="f.name" class="border border-line px-1 py-0.5" /></label>
        <label class="grid grid-cols-[6rem_1fr] items-center gap-2">Description:<input v-model="f.description" class="border border-line px-1 py-0.5" /></label>
        <div class="grid grid-cols-[6rem_1fr] items-center gap-2">Size:
          <span class="flex gap-2"><input v-model.number="f.size" type="number" min="0.01" step="any" class="w-28 border border-line px-1 py-0.5" />
            <select v-model="f.unit" class="border border-line bg-white px-1"><option>GB</option><option>MB</option></select></span>
        </div>
        <div class="grid grid-cols-[6rem_1fr] gap-2">Location:
          <div class="h-36 overflow-auto border border-line">
            <div v-for="s in options.srs" :key="s.uuid" class="cursor-default whitespace-pre px-1 py-0.5" :class="f.sr === s.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100'" @click="f.sr = s.uuid">{{ s.text }}</div>
          </div>
        </div>
        <p v-if="f.error" class="text-red-600">{{ f.error }}</p>
        <div class="flex justify-end gap-2 pt-1">
          <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5">Add</button>
          <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Cancel</button>
        </div>
      </div>
    </form>
  </div>
</template>
