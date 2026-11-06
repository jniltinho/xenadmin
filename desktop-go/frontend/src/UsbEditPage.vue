<script setup lang="ts">
// The USB page of a VM's Properties (USBEditPage, its Attach USB dialog AttachUsbDialog): the VM's USB devices, Attach
// and Detach. As in the C# app these act at once (CreateVUSBAction, DeleteVUSBAction), not on the dialog's OK. The
// dialog adds the page when "shown" says so (PropertiesDialog: HVM, not a template, a USB device in the pool).
import { onMounted, ref, watch } from "vue";
import { api, type USBAttachChoices, type VMUSBPage } from "./api/workspace";

const props = defineProps<{ server: string; vm: string; haPriority: string; disabled?: boolean }>();
const emit = defineEmits<{ (e: "shown", shown: boolean): void }>();

const page = ref<VMUSBPage | null>(null);
const picked = ref("");
const error = ref("");
const busy = ref(false);
const attach = ref<(USBAttachChoices & { pick: string }) | null>(null);
const confirm = ref(false);

async function load() {
  try {
    page.value = await api.vmUsb(props.server, props.vm, props.haPriority);
    if (!page.value.table.keys.includes(picked.value)) picked.value = "";
  } catch (e) { error.value = String(e); }
}
onMounted(async () => { await load(); emit("shown", page.value?.shown ?? false); });
watch(() => props.haPriority, load); // USBEditPage.ShowHideWarnings with the Start Options page's priority

async function act(f: () => Promise<void>) {
  busy.value = true;
  error.value = "";
  try { await f(); } catch (e) { error.value = String(e); }
  await load();
  busy.value = false;
}
async function openAttach() {
  try { attach.value = { ...(await api.usbAttachChoices(props.server, props.vm)), pick: "" }; } catch (e) { error.value = String(e); }
}
function doAttach() {
  const pusb = attach.value?.pick ?? "";
  attach.value = null;
  if (pusb) act(() => api.attachUsb(props.server, props.vm, pusb));
}
function doDetach() {
  const vusb = picked.value;
  confirm.value = false;
  if (vusb) act(() => api.detachUsb(props.server, props.vm, vusb));
}
</script>

<template>
  <div>
    <p v-if="error" class="mb-2 text-red-600">{{ error }}</p>
    <template v-if="page">
      <p class="mb-2 font-semibold">{{ page.title }}</p>
      <div class="flex gap-2">
        <div class="h-48 flex-1 overflow-auto border border-line">
          <table class="w-full">
            <thead><tr class="bg-pane text-left"><th v-for="h in page.table.headers" :key="h" class="px-1 font-normal">{{ h }}</th></tr></thead>
            <tbody>
              <tr v-for="(r, i) in page.table.rows" :key="page.table.keys[i]" class="cursor-default"
                  :class="picked === page.table.keys[i] ? 'bg-select text-white' : 'hover:bg-neutral-100'" @click="picked = page.table.keys[i]">
                <td v-for="(c, j) in r" :key="j" class="px-1">{{ c }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="flex flex-col gap-1">
          <button class="btn min-w-20" :disabled="disabled || busy || !page.can_attach" @click="openAttach">Attach</button>
          <button class="btn min-w-20" :disabled="disabled || busy || !picked" @click="confirm = true">Detach</button>
        </div>
      </div>
      <p v-if="!page.can_attach" class="mt-3 flex items-center gap-2"><img src="/icons/000_Info3_h32bit_16.png" alt="" />{{ page.warning }}</p>
    </template>
    <div v-if="confirm && page" class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
      <div class="w-[26rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ page.confirm_title }}</div>
        <div class="space-y-3 p-3">
          <p>{{ page.confirm }}</p>
          <div class="flex justify-end gap-2">
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="doDetach">Yes</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" autofocus @click="confirm = false">No</button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="attach" class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
      <div class="w-[30rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Attach USB</div>
        <div class="space-y-3 p-3">
          <p>Select a device to be passed through to this virtual machine.</p>
          <div role="listbox" aria-label="USB devices" class="h-40 overflow-auto border border-line">
            <p v-if="!attach.hosts.length" class="px-1 text-neutral-500">{{ attach.empty }}</p>
            <template v-for="h in attach.hosts" :key="h.name">
              <div class="px-1 font-semibold">{{ h.name }}</div>
              <div v-for="d in h.devices" :key="d.uuid" role="option" :aria-selected="attach.pick === d.uuid" class="cursor-default pl-6 pr-1"
                   :class="attach.pick === d.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100'" @click="attach.pick = d.uuid">{{ d.name }}</div>
            </template>
          </div>
          <div class="space-y-1 text-neutral-700">
            <p class="flex items-start gap-2"><img src="/icons/000_Alert2_h32bit_16.png" alt="" class="mt-0.5" />Plugging in untrustworthy USB devices to your computer may put your computer at risk.</p>
            <p class="ml-6">USB devices with modifiable behaviour should only be assigned to trustworthy guest VMs.</p>
            <p v-if="attach.upgrade" class="ml-6">A major component update of VM will be performed the first time when a USB device is passed through to it. This update is irreversible.</p>
          </div>
          <div class="flex justify-end gap-2">
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="!attach.pick" @click="doAttach">Attach</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="attach = null">Cancel</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
