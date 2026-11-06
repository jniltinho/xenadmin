<script setup lang="ts">
// A network's Properties: General (name, description) and Network Settings (EditNetworkPage: the NIC and VLAN of a
// private or VLAN network, the MTU, a bond network's mode and auto-add, with why they are fixed and the disruption
// warning); its folder, tags and custom fields open in ObjectPropertiesDialog.
import { computed, onMounted, reactive, ref } from "vue";
import { api, type NetworkEdit } from "./api/workspace";
import { mtuCheck, vlanCheck } from "./netrules";

const props = defineProps<{ server: string; uuid: string }>();
const emit = defineEmits<{ (e: "close"): void; (e: "run", label: string, f: () => Promise<void>): void; (e: "more"): void }>();

const edit = ref<NetworkEdit | null>(null);
// bond: balance-slb, active-backup, lacp:src_mac or lacp:tcpudp_ports (as BondEditor); "lacp" without the vSwitch
const f = reactive({ name: "", description: "", auto: true, mtu: "", nic: "", vlan: "", bond: "" });
const error = ref("");
onMounted(async () => {
  try {
    const e = (edit.value = await api.networkEdit(props.server, props.uuid));
    // Repopulate: an LACP bond of unknown hashing shows source MAC
    const bond = e.mode === "lacp" && e.lacp ? `lacp:${e.hashing || "src_mac"}` : e.mode;
    Object.assign(f, { name: e.name, description: e.description, auto: e.auto, mtu: String(e.mtu), nic: e.nic, vlan: String(e.vlan), bond });
  } catch (e) { error.value = String(e); }
});
const internal = computed(() => f.nic === ""); // SelectedIsInternal
const nic = computed(() => edit.value?.nics.find((n) => n.uuid === f.nic));
// NetworkAction checks a new VLAN only (the network's own VLANs are in use); VLAN 0 shows labelVLAN0Info
const vlan = computed(() => {
  const e = edit.value;
  if (!e?.nic_state.shown || internal.value) return { ok: true, text: "" };
  if (f.nic === e.nic && f.vlan.trim() !== "" && Number(f.vlan) === e.vlan) return vlanCheck(f.vlan, 0, []);
  return vlanCheck(f.vlan, e.min_vlan, nic.value?.vlans ?? []);
});
// InternalToExternal, ExternalToInternal, ExternalChangedDeviceOrVlan
const nicChanged = computed(() => {
  const e = edit.value;
  return !!e?.nic_state.shown && (f.nic !== e.nic || (!internal.value && Number(f.vlan) !== e.vlan));
});
// SetMTUControlEnablement for the picked NIC: fixed on <internal>, a new VLAN's as UpdateNetwork checks it
const m = computed(() => {
  const e = edit.value;
  if (!e || !nicChanged.value) return e?.mtu_state;
  return internal.value ? { ...e.mtu_state, enabled: false, note: "" } : e.vlan_mtu_state;
});
const mtu = computed(() => (m.value?.enabled ? mtuCheck(f.mtu, m.value.min, m.value.max) : { ok: true, text: "" }));
const mode = computed(() => f.bond.split(":")[0]);
const hashing = computed(() => (mode.value === "lacp" && edit.value?.lacp ? (f.bond.split(":")[1] ?? "") : ""));
// WillDisrupt: MtuHasChanged, BondModeHasChanged, HashingAlgorithmHasChanged or new PIFs
const disrupt = computed(() => {
  const e = edit.value;
  if (!e) return false;
  const mtuChanged = !!m.value?.enabled && mtu.value.ok && Number(f.mtu) !== e.mtu;
  const bondChanged = e.bond && (mode.value !== e.mode || (hashing.value !== "" && hashing.value !== e.hashing));
  return mtuChanged || bondChanged || (nicChanged.value && vlan.value.ok);
});
const why = computed(() =>
  !f.name.trim() ? "Enter a name for the network." : !vlan.value.ok ? vlan.value.text : mtu.value.ok ? "" : mtu.value.text);

function ok() {
  if (why.value || !edit.value) return;
  const ch = { name: f.name.trim(), description: f.description, auto: f.auto, mtu: f.mtu.trim(), nic: f.nic, vlan: f.vlan.trim(), mode: mode.value, hashing: hashing.value };
  emit("run", `Saving the properties of '${ch.name}'`, () => api.saveNetwork(props.server, props.uuid, ch));
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <form novalidate class="max-h-full w-[38rem] overflow-auto border border-neutral-500 bg-white shadow-lg" @submit.prevent="ok">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">'{{ edit?.name ?? "" }}' Properties</div>
      <div v-if="edit" class="space-y-3 p-3">
        <p class="font-semibold">General</p>
        <label class="grid grid-cols-[7rem_1fr] items-center gap-2">Name:<input v-model="f.name" class="border border-line px-1 py-0.5" /></label>
        <label class="grid grid-cols-[7rem_1fr] items-start gap-2">Description:<textarea v-model="f.description" rows="3" class="border border-line px-1 py-0.5" /></label>
        <p class="pt-2 font-semibold">Network Settings</p>
        <p class="text-neutral-600">Use the controls below to configure advanced settings for your network. If your changes will result in temporary disruption of the network you will be informed below. Changes will apply to all hosts in a pool.</p>
        <template v-if="edit.nic_state.shown">
          <div class="grid grid-cols-[7rem_1fr] items-center gap-2">
            <label for="netpnic">NIC:</label>
            <select id="netpnic" v-model="f.nic" class="w-56 border border-line px-1 py-0.5" :disabled="!edit.nic_state.enabled">
              <option value="">&lt;internal&gt;</option>
              <option v-for="n in edit.nics" :key="n.uuid" :value="n.uuid">{{ n.name }}</option>
            </select>
            <span />
            <span class="text-neutral-600">Network interface card. Can also be an &lt;internal&gt; Virtual Network.</span>
            <label for="netpvlan">VLAN:</label>
            <span class="flex items-center gap-2">
              <input id="netpvlan" v-model="f.vlan" class="w-24 border border-line px-1 py-0.5" :disabled="!edit.nic_state.enabled || internal" />
              <span v-if="!internal && edit.nic_state.enabled" :class="vlan.ok ? 'text-neutral-600' : 'text-red-600'">{{ vlan.text }}</span>
            </span>
          </div>
          <p v-if="edit.nic_state.note" class="text-neutral-600">{{ edit.nic_state.note }}</p>
        </template>
        <div v-if="m?.shown" class="flex flex-wrap items-center gap-2">
          <label for="netpmtu" class="w-[7rem]">MTU:</label>
          <input id="netpmtu" v-model="f.mtu" class="w-24 border border-line px-1 py-0.5" :disabled="!m.enabled" />
          <span v-if="m.enabled" :class="mtu.ok ? 'text-neutral-600' : 'text-red-600'">{{ mtu.text }}</span>
          <p v-if="m.note" class="w-full text-neutral-600">{{ m.note }}</p>
        </div>
        <fieldset v-if="edit.bond" class="space-y-1 border border-line px-2 pb-2">
          <legend class="px-1">Bond mode</legend>
          <label class="flex items-center gap-2"><input v-model="f.bond" type="radio" value="balance-slb" />Active-active</label>
          <label class="flex items-center gap-2"><input v-model="f.bond" type="radio" value="active-backup" />Active-passive</label>
          <template v-if="edit.lacp">
            <label class="flex items-center gap-2"><input v-model="f.bond" type="radio" value="lacp:src_mac" />LACP with load balancing based on source MAC address</label>
            <label class="flex items-center gap-2"><input v-model="f.bond" type="radio" value="lacp:tcpudp_ports" />LACP with load balancing based on IP and port of source and destination</label>
          </template>
          <p v-if="f.bond.startsWith('lacp:')" class="text-amber-700">LACP must also be configured on the switch ports</p>
        </fieldset>
        <p v-if="disrupt" class="text-amber-700">Warning: Your network will be temporarily disrupted when you apply these settings</p>
        <label class="flex items-center gap-2"><input v-model="f.auto" type="checkbox" :disabled="!edit.auto_enabled" />Automatically add this network to new virtual machines</label>
        <!-- the folder, tags and custom fields pages (ObjectPropertiesDialog) -->
        <p class="pt-1"><a class="cursor-pointer underline" @click="emit('more')">Folder, tags and custom fields...</a></p>
      </div>
      <p v-if="error || why" class="mx-3 text-red-600">{{ error || why }}</p>
      <div class="flex justify-end gap-2 p-3">
        <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="!!why || !edit">OK</button>
        <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Cancel</button>
      </div>
    </form>
  </div>
</template>
