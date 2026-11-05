<script setup lang="ts">
// The Copy VM dialog (CopyVMDialog): name, description, and a fast clone or a full copy to an SR; the Go side runs
// VMCloneAction or VMCopyAction (CopyVM).
import { onMounted, reactive, ref } from "vue";
import { api, type VMCopy } from "./api/workspace";

const props = defineProps<{ server: string; vm: string }>();
const emit = defineEmits<{ (e: "close", copied: string): void }>();

const info = ref<VMCopy | null>(null);
const f = reactive({ name: "", description: "", full: false, sr: "" });
const busy = ref(false);
const error = ref("");

onMounted(async () => {
  try {
    const c = (info.value = await api.vmCopyChoices(props.server, props.vm));
    Object.assign(f, { name: c.name, description: c.description, full: !c.clone, sr: c.default_sr || c.srs.find((s) => !s.reason)?.uuid || "" });
  } catch (e) { error.value = String(e); }
});

async function copy() {
  if (!f.name.trim()) { error.value = "Enter a name for the copy."; return; }
  if (f.full && !f.sr) { error.value = "Select the storage repository for the copy."; return; }
  busy.value = true; error.value = "";
  try { emit("close", await api.copyVm(props.server, props.vm, f.name.trim(), f.description, f.full, f.full ? f.sr : "")); } catch (e) { error.value = String(e); }
  busy.value = false;
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <div class="flex w-[34rem] flex-col border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">Copy Virtual Machine</div>
      <div v-if="info" class="p-4">
        <label class="mb-2 grid grid-cols-[6.5rem_1fr] items-center gap-2">Name:<input v-model="f.name" class="border border-line px-1 py-0.5" /></label>
        <label class="mb-3 grid grid-cols-[6.5rem_1fr] items-start gap-2">Description:<textarea v-model="f.description" rows="2" class="border border-line px-1 py-0.5" /></label>
        <fieldset class="border border-line px-3 pb-3 pt-1">
          <legend class="px-1">Copy mode</legend>
          <label class="flex items-center gap-2" :class="info.clone ? '' : 'text-neutral-400'"><input v-model="f.full" type="radio" :value="false" :disabled="!info.clone" />Fast clone</label>
          <p class="mb-2 ml-6 text-neutral-600">Clone the existing VM, using a storage-level fast disk clone operation</p>
          <label class="flex items-center gap-2" :class="info.copy ? '' : 'text-neutral-400'"><input v-model="f.full" type="radio" :value="true" :disabled="!info.copy" />Full copy</label>
          <p class="mb-1 ml-6 text-neutral-600">Create a full copy of the existing VM on this storage repository:</p>
          <div class="ml-6 h-28 overflow-auto border border-neutral-400" :class="f.full ? '' : 'bg-neutral-50 text-neutral-400'">
            <div v-for="s in info.srs" :key="s.uuid" class="cursor-default whitespace-pre px-2" :title="s.reason"
                 :class="f.full && f.sr === s.uuid ? 'bg-select text-white' : f.full && !s.reason ? 'hover:bg-neutral-100' : 'text-neutral-400'"
                 @click="f.full && !s.reason && (f.sr = s.uuid)">{{ s.text }}<template v-if="s.reason"> - {{ s.reason }}</template></div>
          </div>
        </fieldset>
      </div>
      <div v-else class="p-4 text-neutral-500">{{ error || "Loading..." }}</div>
      <div class="flex items-center justify-end gap-2 border-t border-line p-2">
        <span v-if="error && info" class="mr-auto text-red-600">{{ error }}</span>
        <button class="btn min-w-20" :disabled="!info || busy" @click="copy">{{ busy ? "Copying..." : "Copy" }}</button>
        <button class="btn min-w-20" :disabled="busy" @click="emit('close', '')">Cancel</button>
      </div>
    </div>
  </div>
</template>
