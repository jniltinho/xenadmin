<script setup lang="ts">
// The GPU page of a VM's Properties (GpuEditPage, its Add virtual GPU dialog AddVGPUDialog): the VM's vGPUs, Add and
// Delete, and the warnings (rules.GPUPage, read from Go after each change). The dialog's OK calls changes(); the
// dialog adds the page when "shown" says the VM can have a GPU and the pool has one (PropertiesDialog).
import { computed, onMounted, ref, watch } from "vue";
import { api, type GPUPageState, type VGPUPicks, type VMGPUPage, type VMGPURow } from "./api/workspace";

const props = defineProps<{ server: string; vm: string; haPriority: string; disabled?: boolean }>();
const emit = defineEmits<{ (e: "shown", shown: boolean): void }>();

const page = ref<VMGPUPage | null>(null);
const rows = ref<VMGPURow[]>([]);
const picked = ref(-1);
const st = ref<GPUPageState | null>(null);
const error = ref("");
const add = ref<(VGPUPicks & { pick: number }) | null>(null);

// GpuEditPage.HasChanged: VGpus.SequenceEqual(currentGpus) (a new row has no uuid)
const changed = () => JSON.stringify(rows.value.map((r) => r.uuid)) !== JSON.stringify((page.value?.rows ?? []).map((r) => r.uuid));
async function refresh() {
  const p = page.value;
  if (!p) return;
  try {
    st.value = await api.vmGpuState({ ...p.input, ha_priority: props.haPriority, rows: rows.value.map((r) => r.facts), changed: changed(), selected: picked.value >= 0 });
  } catch (e) { error.value = String(e); }
}
async function load() {
  page.value = await api.vmGpu(props.server, props.vm);
  rows.value = [...page.value.rows];
  picked.value = -1;
  await refresh();
}
onMounted(async () => {
  try { await load(); emit("shown", page.value!.shown); } catch (e) { error.value = String(e); emit("shown", false); }
});
// after a failed save: the page shows what the server has, so OK never makes a vGPU twice
async function reload() { try { await load(); } catch (e) { page.value = null; error.value = String(e); } }
watch(() => props.haPriority, refresh);

async function openAdd() {
  try {
    const c = await api.vgpuChoices(props.server, props.vm, rows.value.map((r) => r.type));
    add.value = { ...c, pick: c.picked };
  } catch (e) { error.value = String(e); }
}
function doAdd() {
  const item = add.value?.items[add.value.pick];
  add.value = null;
  if (!item || item.header || !item.enabled) return;
  rows.value.push(item.row);
  refresh();
}
function remove() {
  if (picked.value < 0) return;
  rows.value.splice(picked.value, 1);
  picked.value = -1;
  refresh();
}
function pick(i: number) { picked.value = i; refresh(); }

// null: nothing changed; else the vGPUs read that were removed and the ones to create (GpuAssignAction)
function changes(): null | { removed: string[]; add: { group: string; type: string }[] } {
  if (!page.value || !changed()) return null;
  const kept = new Set(rows.value.map((r) => r.uuid));
  return { removed: page.value.rows.filter((r) => !kept.has(r.uuid)).map((r) => r.uuid),
    add: rows.value.filter((r) => !r.uuid).map((r) => ({ group: r.group, type: r.type })) };
}
// VMHAEditPage.RefillPrioritiesComboBox: a VM with vGPUs (saved or not) is not offered the Restart priority
const count = computed(() => rows.value.length);
defineExpose({ changes, count, reload });

// GpuEditPage.resx, top to bottom (warningsTable)
const lines = (s: GPUPageState) => [
  [s.experimental, "000_Alert2_h32bit_16", "This type of vGPU is in Tech Preview, and is currently not suitable for production use."],
  [s.rdp, "000_Alert2_h32bit_16", "Before making this change, you must make sure that you can connect to your VM over Remote Desktop. The Default (VNC) Desktop will not work with a dedicated GPU, so you will not be able to connect to your VM except over Remote Desktop."],
  [s.need_gpu, "000_Alert2_h32bit_16", "It is essential that the appropriate graphics card drivers are installed on the VM once it has been booted."],
  [s.need_gpu, "000_Info3_h32bit_16", "If there is no available virtual GPU of the correct type when the VM is next booted, it will not be able to start."],
  [s.stop_vm, "000_Info3_h32bit_16", "The virtual GPUs can only be changed when the VM is shut down."],
  [s.ha, "000_Info3_h32bit_16", "The virtual GPUs cannot be changed because the VM is currently protected by HA."],
  [!!s.multiple, "000_Info3_h32bit_16", s.multiple],
].filter((l) => l[0]) as [boolean, string, string][];
</script>

<template>
  <div>
    <p class="mb-3">You can improve graphics performance by assigning one or more virtual graphics processing units to this VM.</p>
    <p v-if="error" class="mb-2 text-red-600">{{ error }}</p>
    <template v-if="page">
      <div class="flex gap-2">
        <div class="h-36 flex-1 overflow-auto border border-line">
          <table class="w-full">
            <thead><tr class="bg-pane text-left">
              <th v-if="page.rows.length" class="px-1 font-normal">PCI slot number</th><th class="px-1 font-normal">Virtual GPU type</th>
              <th class="px-1 font-normal">Virtual GPUs per GPU</th><th class="px-1 font-normal">Video RAM</th></tr></thead>
            <tbody>
              <tr v-for="(r, i) in rows" :key="i" class="cursor-default" :class="picked === i ? 'bg-select text-white' : 'hover:bg-neutral-100'" @click="pick(i)">
                <td v-if="page.rows.length" class="px-1">{{ r.slot }}</td><td class="px-1">{{ r.name }}</td><td class="px-1">{{ r.per_gpu }}</td><td class="px-1">{{ r.vram }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="flex flex-col gap-1">
          <button class="btn min-w-20" :disabled="disabled || !st?.can_add" @click="openAdd">Add...</button>
          <button class="btn min-w-20" :disabled="disabled || !st?.can_delete" @click="remove">Delete</button>
        </div>
      </div>
      <div v-if="st" class="mt-3 space-y-2">
        <p v-for="l in lines(st)" :key="l[2]" class="flex items-start gap-2"><img :src="`/icons/${l[1]}.png`" alt="" class="mt-0.5" />{{ l[2] }}</p>
      </div>
    </template>
    <div v-if="add" class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
      <div class="w-[28rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Add virtual GPU</div>
        <div class="space-y-3 p-3">
          <label for="vgpu-type">Select virtual GPU type:</label>
          <div id="vgpu-type" role="listbox" class="h-48 overflow-auto border border-line">
            <div v-for="(it, i) in add.items" :key="i" role="option" :aria-selected="add.pick === i" :aria-disabled="it.header || !it.enabled"
                 class="cursor-default px-1" :class="[it.header ? 'font-semibold' : !it.enabled ? 'text-neutral-400' : add.pick === i ? 'bg-select text-white' : 'hover:bg-neutral-100', it.sub ? 'pl-6' : '']"
                 @click="!it.header && it.enabled && (add.pick = i)">{{ it.text }}</div>
          </div>
          <div class="flex justify-end gap-2">
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="add.pick < 0" @click="doAdd">Add</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="add = null">Cancel</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
