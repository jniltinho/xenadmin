<script setup lang="ts">
// The Delete VM dialog (ConfirmVMDeleteDialog): the disks only this VM uses, checked when it owns them, and its
// snapshots, unchecked; Delete runs VMDestroyAction on the Go side (DeleteVM), which checks the choices again.
import { onMounted, ref } from "vue";
import { api, type VMDelete } from "./api/workspace";

const props = defineProps<{ server: string; vm: string }>();
const emit = defineEmits<{ (e: "close", deleted: boolean): void }>();

const info = ref<VMDelete | null>(null);
const checked = ref(new Set<string>());
const busy = ref(false);
const error = ref("");
const gb = (b?: number) => `${+((b ?? 0) / 1024 ** 3).toFixed(2)} GB`;

onMounted(async () => {
  try {
    info.value = await api.vmDeleteChoices(props.server, props.vm);
    checked.value = new Set(info.value.disks.filter((d) => d.checked).map((d) => d.uuid));
  } catch (e) { error.value = String(e); }
});
const all = () => [...(info.value?.disks ?? []), ...(info.value?.snapshots ?? [])];
function toggle(uuid: string) { const s = new Set(checked.value); if (!s.delete(uuid)) s.add(uuid); checked.value = s; }

async function remove() {
  if (!info.value) return;
  busy.value = true; error.value = "";
  try {
    await api.deleteVm(props.server, props.vm, info.value.disks.filter((d) => checked.value.has(d.uuid)).map((d) => d.uuid),
      info.value.snapshots.filter((d) => checked.value.has(d.uuid)).map((d) => d.uuid));
    emit("close", true);
  } catch (e) { error.value = String(e); }
  busy.value = false;
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <div class="flex w-[34rem] flex-col border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ info?.title ?? "Delete VM" }}</div>
      <div v-if="info" class="p-4">
        <div class="mb-3 flex gap-3">
          <img src="/icons/000_WarningAlert_h32bit_32.png" alt="" class="h-8 w-8 shrink-0" />
          <div><p>{{ info.text }}</p><p class="font-semibold">Deleted VMs and Templates cannot be recovered.</p></div>
        </div>
        <template v-if="all().length">
          <p class="mb-1">The selected virtual disks and snapshots below will also be deleted:</p>
          <div class="h-40 overflow-auto border border-neutral-400">
            <template v-for="[title, items] in [['Attached virtual disks', info.disks], ['Snapshots', info.snapshots]] as const" :key="title">
              <div v-if="items.length" class="bg-pane px-2 py-0.5 font-semibold text-select">{{ title }}</div>
              <label v-for="d in items" :key="d.uuid" class="flex items-center gap-2 px-2 py-0.5 hover:bg-neutral-100">
                <input type="checkbox" :checked="checked.has(d.uuid)" @change="toggle(d.uuid)" />
                <span class="flex-1 truncate">{{ d.name }}</span><span v-if="d.size" class="text-neutral-600">{{ gb(d.size) }}</span></label>
            </template>
          </div>
          <div class="mt-2 flex gap-2">
            <button class="btn" @click="checked = new Set(all().map((d) => d.uuid))">Select All</button>
            <button class="btn" @click="checked = new Set()">Clear All</button>
          </div>
        </template>
      </div>
      <div v-else class="p-4 text-neutral-500">{{ error || "Loading..." }}</div>
      <div class="flex items-center justify-end gap-2 border-t border-line p-2">
        <span v-if="error && info" class="mr-auto text-red-600">{{ error }}</span>
        <button class="btn min-w-20" :disabled="!info || busy" @click="remove">{{ busy ? "Deleting..." : "Delete" }}</button>
        <button class="btn min-w-20" :disabled="busy" @click="emit('close', false)">Cancel</button>
      </div>
    </div>
  </div>
</template>
