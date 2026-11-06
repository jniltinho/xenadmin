<script setup lang="ts">
// The Memory, Storage, Networking and NICs tabs of a server or pool, and the Storage tab of an SR (Details' tables);
// the Networking tab's Add Network... / Properties / Remove Network (HostNetworkView) and the NICs tab's Create Bond... /
// Delete Bond / Rescan (NicsView). shortcut: the Storage tabs have no buttons yet.
import { computed, onMounted, ref, watch } from "vue";
import { api, type HostNetworking, type HostNics, type Table } from "./api/workspace";
import DetailsTable from "./DetailsTable.vue";
import NetworkPropertiesDialog from "./NetworkPropertiesDialog.vue";
import NewNetworkWizard from "./NewNetworkWizard.vue";
import { serial } from "./serial";

// tick: the server's objects changed (xcp:changed), read the tab again
const props = defineProps<{ server: string; kind: string; uuid: string; tab: string; tick?: number }>();
const emit = defineEmits<{ (e: "status", text: string, error: boolean): void }>();

const table = ref<Table | null>(null);
const nets = ref<HostNetworking | null>(null);
const nics = ref<HostNics | null>(null);
const picked = ref("");
const busy = ref(false);
const load = serial(async () => {
  try {
    if (props.tab === "networking") table.value = (nets.value = await api.hostNetworking(props.server, props.kind, props.uuid)).table;
    else if (props.tab === "nics") table.value = (nics.value = await api.hostNics(props.server, props.uuid)).table;
    else table.value = await api.hostTab(props.server, props.kind, props.uuid, props.tab);
    if (props.tab !== "memory" && !table.value.keys.includes(picked.value)) picked.value = table.value.keys[0] ?? "";
  } catch (e) { emit("status", String(e), true); }
});
onMounted(load);
watch(() => props.tick, load);

async function run(label: string, f: () => Promise<void>) {
  wizard.value = props_.value = null;
  busy.value = true;
  emit("status", `${label}...`, false);
  try { await f(); emit("status", "", false); } catch (e) { emit("status", `${label}: ${String(e)}`, true); }
  await load(); // the buttons stay disabled until the tab shows the new state
  busy.value = false;
}

const network = computed(() => nets.value?.networks.find((n) => n.uuid === picked.value) ?? null);
const nic = computed(() => nics.value?.nics.find((n) => n.uuid === picked.value) ?? null);
const wizard = ref<{ bond: boolean } | null>(null);
const props_ = ref<string | null>(null); // the network whose Properties are open

// the questions: ok "" is an error to read (a blocked bond change)
const ask = ref<{ title: string; text: string; ok: string; cancel: string; resolve: (ok: boolean) => void } | null>(null);
const question = (title: string, text: string, ok: string, cancel: string) =>
  new Promise<boolean>((resolve) => (ask.value = { title, text, ok, cancel, resolve }));
function answer(ok: boolean) { ask.value?.resolve(ok); ask.value = null; }

// DestroyBondCommand: the question (or the error under HA), then DestroyBondAction
async function deleteBond(networkUUID: string) {
  try {
    const w = await api.deleteBondQuestion(props.server, networkUUID);
    if (w.block) { await question(w.title, w.text, "", "OK"); return; }
    if (await question(w.title, w.text, w.ok, "Cancel")) run("Deleting bond", () => api.deleteBond(props.server, networkUUID));
  } catch (e) { emit("status", String(e), true); }
}
async function removeNetwork() {
  const n = network.value;
  if (!n?.can_remove) return;
  if (n.bond) return deleteBond(n.uuid);
  const q = await api.removeNetworkQuestion();
  if (await question(q.title, q.text, "Yes", "No")) run("Removing network", () => api.removeNetwork(props.server, n.uuid));
}
function edit() { if (network.value?.can_edit) props_.value = network.value.uuid; }
</script>

<template>
  <div class="flex min-h-0 flex-1 flex-col overflow-auto">
    <div v-if="table" @dblclick="tab === 'networking' && !busy && edit()"><!-- NetworkList: a double click is Properties -->
      <DetailsTable :table="table" :selected="picked" :plain="tab === 'memory'" @select="picked = $event" />
    </div>
    <div v-if="tab === 'networking' && nets" class="flex gap-2 px-3">
      <button class="btn" :disabled="busy || !nets.can_add" @click="wizard = { bond: false }">Add Network...</button>
      <button class="btn" :disabled="busy || !network?.can_edit" @click="edit">Properties</button>
      <button class="btn" :disabled="busy || !network?.can_remove" :title="network?.remove_tip" @click="removeNetwork">Remove Network</button>
    </div>
    <div v-if="tab === 'nics' && nics" class="flex gap-2 px-3">
      <button class="btn" :disabled="busy || !nics.can_bond" @click="wizard = { bond: true }">Create Bond...</button>
      <button class="btn" :disabled="busy || !nic?.bond_network" @click="nic?.bond_network && deleteBond(nic.bond_network)">Delete Bond</button>
      <button class="btn" :disabled="busy || !nics.can_bond" @click="run('Scanning for NICs', () => api.rescanNics(server, uuid))">Rescan</button>
    </div>
    <NewNetworkWizard v-if="wizard" :server="server" :kind="kind" :uuid="uuid" :bond="wizard.bond" @close="wizard = null" @run="run" />
    <NetworkPropertiesDialog v-if="props_" :server="server" :uuid="props_" @close="props_ = null" @run="run" />
    <div v-if="ask" class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
      <div class="w-[32rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ ask.title }}</div>
        <div class="space-y-3 p-3">
          <p class="whitespace-pre-line">{{ ask.text }}</p>
          <div class="flex justify-end gap-2">
            <button v-if="ask.ok" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="answer(true)">{{ ask.ok }}</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" autofocus @click="answer(false)">{{ ask.cancel }}</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
