<script setup lang="ts">
// Add Virtual Interface / Virtual Interface Properties (VIFDialog): the network, an auto-generated or typed MAC
// address and an optional QoS limit; NetworkRules.Validate on the Go side, and the duplicate MAC question first.
import { computed, onMounted, reactive, ref } from "vue";
import { api, type Confirm, type VIF, type VIFSettings } from "./api/workspace";

const props = defineProps<{ server: string; vm: string; vif: VIF | null }>();
const emit = defineEmits<{ (e: "close"): void; (e: "save", set: VIFSettings): void }>();

const networks = ref<{ uuid: string; name: string }[] | null>(null);
const f = reactive({ network: "", manual: false, mac: "", qos: false, limit: "" });
const error = ref("");
const ask = ref<Confirm | null>(null);

onMounted(async () => {
  try {
    networks.value = (await api.vifNetworks(props.server, props.vm)) ?? [];
    const v = props.vif;
    f.network = networks.value.find((n) => n.uuid === v?.network)?.uuid ?? networks.value[0]?.uuid ?? "";
    if (v) Object.assign(f, { manual: v.mac !== "", mac: v.mac, qos: v.qos, limit: v.limit });
  } catch (e) { error.value = String(e); }
});

// NetworkRules.Validate: the reason, or "" with OK disabled when there is nothing to say
const why = computed(() => {
  if (!f.network) return { text: "There are no networks available for this network interface", ok: false };
  if (f.manual) {
    if (!f.mac.trim()) return { text: "", ok: false };
    if (!/^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$/.test(f.mac.trim())) return { text: "Invalid MAC address.", ok: false };
  }
  if (f.qos) {
    if (!f.limit.trim()) return { text: "", ok: false };
    if (!/^\d+$/.test(f.limit.trim()) || Number(f.limit) <= 0) return { text: "The QoS limit value should be an integer greater than zero", ok: false };
  }
  return { text: "", ok: true };
});
const settings = (): VIFSettings => ({ network: f.network, mac: f.manual ? f.mac.trim() : "", qos: f.qos, limit: f.limit.trim() });

// one OK at a time, the form frozen meanwhile: the values checked for a duplicate MAC are the ones saved
const pending = ref(false);
const checked = ref<VIFSettings | null>(null);
async function ok() {
  if (!why.value.ok || pending.value) return;
  pending.value = true;
  const set = settings();
  try {
    const q = await api.vifQuestion(props.vif?.uuid ?? "", set.mac);
    if (q.title) { checked.value = set; ask.value = q; return; }
    emit("save", set);
  } catch (e) { error.value = String(e); } finally { pending.value = false; }
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <form v-if="!ask" novalidate class="w-[34rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="ok">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ vif ? "Virtual Interface Properties" : "Add Virtual Interface" }}</div>
      <fieldset :disabled="pending" class="space-y-3 p-3">
        <p>Select your network and MAC address for this virtual interface and, optionally, specify a QoS limit.</p>
        <div class="grid grid-cols-[auto_1fr] items-start gap-x-3 gap-y-3">
          <label for="vifnet" class="pt-0.5">Network:</label>
          <select id="vifnet" v-model="f.network" class="border border-line bg-white px-1 py-0.5" :disabled="!networks?.length">
            <option v-if="networks && !networks.length" value="">None</option>
            <option v-for="n in networks ?? []" :key="n.uuid" :value="n.uuid">{{ n.name }}</option>
          </select>
          <span class="pt-0.5">MAC address:</span>
          <div class="space-y-1">
            <label class="flex items-center gap-2"><input v-model="f.manual" type="radio" :value="false" />Auto-generate a MAC address</label>
            <div class="flex items-center gap-2">
              <label class="flex items-center gap-2"><input v-model="f.manual" type="radio" :value="true" />Use this MAC address:</label>
              <input v-model="f.mac" placeholder="aa:bb:cc:dd:ee:ff" aria-label="MAC address" class="w-48 border border-line px-1 py-0.5" @focus="f.manual = true" />
            </div>
          </div>
          <span class="pt-0.5">QoS settings:</span>
          <div class="flex items-center gap-2">
            <label class="flex items-center gap-2"><input v-model="f.qos" type="checkbox" />Enable a QoS limit of:</label>
            <input v-model="f.limit" aria-label="QoS limit" class="w-24 border border-line px-1 py-0.5" @focus="f.qos = true" /><span>Kbytes/s</span>
          </div>
        </div>
        <p v-if="error || why.text" class="text-red-600">{{ error || why.text }}</p>
        <div class="flex justify-end gap-2 pt-1">
          <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="!why.ok || pending">{{ vif ? "OK" : "Add" }}</button>
          <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Cancel</button>
        </div>
      </fieldset>
    </form>
    <div v-else class="w-[32rem] border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ ask.title }}</div>
      <div class="space-y-3 p-3">
        <p class="whitespace-pre-line">{{ ask.text }}</p>
        <div class="flex justify-end gap-2">
          <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('save', checked!)">Yes</button>
          <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" autofocus @click="ask = null">No</button>
        </div>
      </div>
    </div>
  </div>
</template>
