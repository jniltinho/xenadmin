<script setup lang="ts">
// A network's Properties: General (name, description) and Network Settings (EditNetworkPage: auto-add and MTU, with
// the reason the MTU is fixed). shortcut: no NIC/VLAN change, bond mode or tags/folder/custom fields pages; add them
// with the full Properties dialog of networks.
import { computed, onMounted, reactive, ref } from "vue";
import { api, type NetworkEdit } from "./api/workspace";
import { mtuCheck } from "./netrules";

const props = defineProps<{ server: string; uuid: string }>();
const emit = defineEmits<{ (e: "close"): void; (e: "run", label: string, f: () => Promise<void>): void }>();

const edit = ref<NetworkEdit | null>(null);
const f = reactive({ name: "", description: "", auto: true, mtu: "" });
const error = ref("");
onMounted(async () => {
  try {
    const e = (edit.value = await api.networkEdit(props.server, props.uuid));
    Object.assign(f, { name: e.name, description: e.description, auto: e.auto, mtu: String(e.mtu) });
  } catch (e) { error.value = String(e); }
});
const m = computed(() => edit.value?.mtu_state);
const mtu = computed(() => (m.value?.enabled ? mtuCheck(f.mtu, m.value.min, m.value.max) : { ok: true, text: "" }));
const why = computed(() => (!f.name.trim() ? "Enter a name for the network." : mtu.value.ok ? "" : mtu.value.text));

function ok() {
  if (why.value || !edit.value) return;
  const ch = { name: f.name.trim(), description: f.description, auto: f.auto, mtu: f.mtu.trim() };
  emit("run", `Saving the properties of '${ch.name}'`, () => api.saveNetwork(props.server, props.uuid, ch));
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <form novalidate class="w-[36rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="ok">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">'{{ edit?.name ?? "" }}' Properties</div>
      <div v-if="edit" class="space-y-3 p-3">
        <p class="font-semibold">General</p>
        <label class="grid grid-cols-[7rem_1fr] items-center gap-2">Name:<input v-model="f.name" class="border border-line px-1 py-0.5" /></label>
        <label class="grid grid-cols-[7rem_1fr] items-start gap-2">Description:<textarea v-model="f.description" rows="3" class="border border-line px-1 py-0.5" /></label>
        <p class="pt-2 font-semibold">Network Settings</p>
        <div v-if="m?.shown" class="flex flex-wrap items-center gap-2">
          <label for="netpmtu" class="w-[7rem]">MTU:</label>
          <input id="netpmtu" v-model="f.mtu" class="w-24 border border-line px-1 py-0.5" :disabled="!m.enabled" />
          <span v-if="m.enabled" :class="mtu.ok ? 'text-neutral-600' : 'text-red-600'">{{ mtu.text }}</span>
          <p v-if="m.note" class="w-full text-neutral-600">{{ m.note }}</p>
          <!-- EditNetworkPage.WillDisrupt: a new MTU unplugs and plugs the network's interfaces -->
          <p v-if="m.enabled && mtu.ok && Number(f.mtu) !== edit.mtu" class="w-full text-amber-700">Warning: Your network will be temporarily disrupted when you apply these settings</p>
        </div>
        <label class="flex items-center gap-2"><input v-model="f.auto" type="checkbox" :disabled="!edit.auto_enabled" />Automatically add this network to new virtual machines</label>
      </div>
      <p v-if="error || why" class="mx-3 text-red-600">{{ error || why }}</p>
      <div class="flex justify-end gap-2 p-3">
        <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="!!why || !edit">OK</button>
        <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Cancel</button>
      </div>
    </form>
  </div>
</template>
