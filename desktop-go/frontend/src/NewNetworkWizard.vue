<script setup lang="ts">
// The New Network wizard (XenAdmin.Desktop/NewNetworkWizard.cs): Select Type, Details (NetWName), Interface
// (NetWDetails) for External and Single-Server Private networks, or Select Type, Bond Members (NetWBondDetails) for a
// bonded network; with bond set it is Create Bond (BondProperties: the BondDetails control alone). An SR-IOV Network
// (NetWTypeSelect, NetWSriovDetails) is Select Type, Details and its NIC, then a question (CreateSriovAction). An
// External network on a NIC with an SR-IOV network can be a VLAN on that network (NetWDetails' checkBoxSriov).
import { computed, onMounted, reactive, ref, watch } from "vue";
import { api, type BondSpec, type BondWarning, type NewNetworkChoices } from "./api/workspace";
import BondEditor from "./BondEditor.vue";
import { clampMtu, MTU_MIN, mtuCheck, vlanCheck, vlanFor } from "./netrules";

const props = defineProps<{ server: string; kind: string; uuid: string; bond?: boolean }>();
const emit = defineEmits<{ (e: "close"): void; (e: "run", label: string, f: () => Promise<void>): void }>();

const opts = ref<NewNetworkChoices | null>(null);
const error = ref("");
const type = ref<"external" | "private" | "bond" | "sriov">(props.bond ? "bond" : "external");
const page = ref(props.bond ? 1 : 0);
const f = reactive({ name: "", description: "", auto: true, nic: "", vlan: "0", mtu: String(MTU_MIN), sriovNic: "", onSriov: false });
const bondSpec = ref<BondSpec>({ members: [], mode: "balance-slb", hashing: "", mtu: MTU_MIN, auto: true });
const bondValid = ref(false);

onMounted(async () => {
  try {
    opts.value = await api.newNetworkOptions(props.server, props.kind, props.uuid);
    f.name = opts.value.external_name;
    f.vlan = String(opts.value.min_vlan);
    pickNic(opts.value.nics[0]?.uuid ?? "");
    f.sriovNic = opts.value.sriov_nics[0]?.uuid ?? "";
  } catch (e) { error.value = String(e); }
});
// NetWName.PopulatePage: "New Network" or "New Private Network", made unique
watch(type, (t) => {
  if (opts.value && t !== "bond") f.name = t === "external" ? opts.value.external_name : t === "sriov" ? opts.value.sriov_name : opts.value.private_name;
});

const nic = computed(() => opts.value?.nics.find((n) => n.uuid === f.nic) ?? null);
// NetWDetails.CreateVlanOnSriovNetwork (checkBoxSriov shown and ticked) and GetVLANList: the VLANs of the NIC's SR-IOV
// network, or those of the others; ticking it checks the VLAN again without moving it
const onSriov = computed(() => !!nic.value?.sriov && f.onSriov);
const inUse = computed(() => (nic.value ? (onSriov.value ? nic.value.sriov_vlans : nic.value.vlans) : []));
// NetWDetails: a VLAN already used on the NIC moves to the first free one; the MTU stays inside the NIC's range
function pickNic(uuid: string) {
  f.nic = uuid;
  const n = nic.value;
  if (!n) return;
  f.vlan = String(vlanFor(Number(f.vlan) || opts.value!.min_vlan, inUse.value));
  f.mtu = n.max_mtu > MTU_MIN ? clampMtu(f.mtu, n.max_mtu) : String(n.max_mtu);
}
const vlan = computed(() => (nic.value ? vlanCheck(f.vlan, opts.value!.min_vlan, inUse.value) : { ok: false, text: "" }));
const mtu = computed(() => mtuCheck(f.mtu, MTU_MIN, nic.value?.max_mtu ?? MTU_MIN));

const pages = computed(() => (props.bond ? [] : type.value === "bond" ? ["Select Type", "Bond Members"] : ["Select Type", "Details", "Interface"]));
const titles: Record<string, string> = {
  "Select Type": "Choose the type of network to create", "Bond Members": "Select members for the new bonded network",
  Details: "Enter a name and description for the new network", Interface: "Configure the new network",
};
const step = computed(() => (props.bond ? "Bond Members" : pages.value[page.value]));
const last = computed(() => props.bond || page.value === pages.value.length - 1);
// NewNetworkWizard.Validate
const why = computed(() => {
  if (step.value === "Bond Members") return bondValid.value ? "" : "Select two or more NICs that are not already bonded.";
  if (step.value === "Details" && !f.name.trim()) return "Enter a name for the new network.";
  if (step.value === "Interface" && type.value === "sriov" && !f.sriovNic) return "Select a NIC.";
  if (step.value === "Interface" && type.value === "external") {
    if (!nic.value) return "Select a NIC.";
    if (!vlan.value.ok) return vlan.value.text;
    if (!mtu.value.ok) return mtu.value.text;
  }
  return "";
});

const pending = ref(false);
const ask = ref<BondWarning | null>(null);
async function finish() {
  if (why.value || pending.value || !opts.value) return;
  if (type.value === "sriov") { // NetWSriovDetails.PageLeaveCore: the servers' connections are disturbed
    ask.value = { title: "XCP-ng Center", text: opts.value.sriov_warning, ok: opts.value.sriov_ok, block: false };
    return;
  }
  if (type.value !== "bond") {
    const o = opts.value, spec = { external: type.value === "external", name: f.name.trim(), description: f.description, auto: f.auto, host: o.host, nic: f.nic, vlan: f.vlan, mtu: f.mtu, sriov: type.value === "external" && onSriov.value };
    emit("run", `Creating network '${spec.name}'`, () => api.createNetwork(props.server, spec));
    return;
  }
  // BondDetails.CanCreateBond: an error stops, a warning asks
  pending.value = true;
  try {
    const q = await api.bondQuestion(props.server, bondSpec.value.members);
    if (q?.block) error.value = q.text;
    else if (q) ask.value = q;
    else createBond();
  } catch (e) { error.value = String(e); } finally { pending.value = false; }
}
function confirmed() {
  if (type.value !== "sriov") return createBond();
  const spec = { name: f.name.trim(), description: f.description, auto: f.auto, host: opts.value!.host, nic: f.sriovNic };
  ask.value = null;
  emit("run", `Creating network '${spec.name}'`, () => api.createSriov(props.server, spec));
}
function createBond() {
  const spec = { ...bondSpec.value, members: [...bondSpec.value.members] };
  ask.value = null;
  emit("run", "Creating bond", () => api.createBond(props.server, spec));
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <div v-if="!ask" class="flex w-[50rem] flex-col border border-neutral-500 bg-white shadow-lg" :class="bond ? '' : 'h-[32rem]'">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ bond ? "Create Bond" : opts?.title ?? "New Network" }}</div>
      <div v-if="!bond" class="border-b border-line px-4 py-2 text-base font-semibold">{{ titles[step] }}</div>
      <div class="flex min-h-0 flex-1">
        <ol v-if="!bond" class="w-40 shrink-0 border-r border-line bg-pane py-2">
          <li v-for="(p, i) in pages" :key="p" class="px-3 py-1" :class="i === page ? 'font-semibold text-select' : i < page ? 'text-neutral-700' : 'text-neutral-400'">{{ p }}</li>
        </ol>
        <div v-if="opts" class="min-w-0 flex-1 overflow-auto p-4">
          <div v-if="step === 'Select Type'" class="space-y-1">
            <p class="mb-2">Select the type of new network you would like to create:</p>
            <label class="mt-2 flex items-center gap-2"><input v-model="type" type="radio" value="external" />External Network</label>
            <p class="ml-6 text-neutral-600">Create a network that passes traffic over one of your VLANs.</p>
            <label class="mt-2 flex items-center gap-2"><input v-model="type" type="radio" value="private" />Single-Server Private Network</label>
            <p class="ml-6 text-neutral-600">Create a network that does not leave each server. This can be used as a private connection between VMs on the same host.</p>
            <label class="mt-2 flex items-center gap-2"><input v-model="type" type="radio" value="bond" />Bonded Network</label>
            <p class="ml-6 text-neutral-600">Create a network that bonds together two or more of your NICs. This will create a single higher performing channel.</p>
            <label class="mt-2 flex items-center gap-2" :class="opts.sriov ? '' : 'text-neutral-400'"><input v-model="type" type="radio" value="sriov" :disabled="!opts.sriov" />SR-IOV Network</label>
            <p class="ml-6" :class="opts.sriov ? 'text-neutral-600' : 'text-neutral-400'">Enable SR-IOV on a NIC and create an SR-IOV network on that NIC.</p>
            <p v-if="!opts.sriov" class="ml-6 flex items-center gap-2"><img src="/icons/000_Info3_h32bit_16.png" alt="" />{{ opts.sriov_why }}</p>
          </div>
          <BondEditor v-else-if="step === 'Bond Members'" v-model="bondSpec" v-model:valid="bondValid" :choices="opts.bond" :disabled="pending" />
          <div v-else-if="step === 'Details'" class="space-y-3">
            <p>Provide a name and optional description for the new network.</p>
            <label class="grid grid-cols-[6rem_1fr] items-center gap-2">Name:<input v-model="f.name" class="border border-line px-1 py-0.5" /></label>
            <label class="grid grid-cols-[6rem_1fr] items-start gap-2">Description:<textarea v-model="f.description" rows="3" class="border border-line px-1 py-0.5" /></label>
          </div>
          <div v-else-if="type === 'sriov'" class="space-y-3">
            <p>Select a NIC to enable SR-IOV.</p>
            <label class="flex items-center gap-3">NIC:<select v-model="f.sriovNic" class="w-56 border border-line bg-white px-1 py-0.5">
              <option v-for="n in opts.sriov_nics" :key="n.uuid" :value="n.uuid">{{ n.name }}</option></select></label>
            <label class="flex items-center gap-2"><input v-model="f.auto" type="checkbox" />Automatically add this network to new virtual machines.</label>
          </div>
          <div v-else-if="type === 'private'" class="space-y-3">
            <p>Select whether you would like to automatically add this network to new VMs.</p>
            <label class="flex items-center gap-2"><input v-model="f.auto" type="checkbox" />Automatically add this network to new virtual machines</label>
          </div>
          <div v-else class="space-y-3">
            <p>Your new network will be mapped to an existing physical network interface and assigned a VLAN number to use on that interface. You can select the physical interface you would like to use below.</p>
            <div class="grid grid-cols-[4rem_auto_1fr] items-center gap-x-3 gap-y-2">
              <label for="netnic">NIC:</label>
              <select id="netnic" class="w-56 border border-line bg-white px-1 py-0.5" :value="f.nic" @change="pickNic(($event.target as HTMLSelectElement).value)">
                <option v-for="n in opts.nics" :key="n.uuid" :value="n.uuid">{{ n.name }}</option>
              </select><span />
              <label for="netvlan">VLAN:</label><input id="netvlan" v-model="f.vlan" class="w-24 border border-line px-1 py-0.5" />
              <span :class="vlan.ok ? '' : 'text-red-600'">{{ vlan.text }}</span>
              <label for="netmtu">MTU:</label><input id="netmtu" v-model="f.mtu" class="w-24 border border-line px-1 py-0.5" :disabled="(nic?.max_mtu ?? MTU_MIN) <= MTU_MIN" />
              <span :class="mtu.ok ? 'text-neutral-600' : 'text-red-600'">{{ mtu.text }}</span>
            </div>
            <label v-if="nic?.sriov" class="flex items-center gap-2"><input v-model="f.onSriov" type="checkbox" />Create the VLAN on the SR-IOV network</label>
            <label class="flex items-center gap-2"><input v-model="f.auto" type="checkbox" />Automatically add this network to new virtual machines</label>
          </div>
        </div>
      </div>
      <p v-if="error" class="mx-4 mb-2 whitespace-pre-line text-red-600">{{ error }}</p>
      <div class="flex justify-end gap-2 border-t border-line px-4 py-2">
        <template v-if="!bond">
          <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="page === 0 || pending" @click="page--; error = ''">&lt; Previous</button>
          <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="last || !!why || !opts" :title="why" @click="page++">Next &gt;</button>
        </template>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="!last || !!why || pending || !opts" :title="why" @click="finish">
          {{ bond ? "Create" : "Finish" }}</button>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close')">Cancel</button>
      </div>
    </div>
    <div v-else class="w-[32rem] border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ ask.title }}</div>
      <div class="space-y-3 p-3">
        <p class="whitespace-pre-line">{{ ask.text }}</p>
        <div class="flex justify-end gap-2">
          <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="confirmed">{{ ask.ok }}</button>
          <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" autofocus @click="ask = null">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>
