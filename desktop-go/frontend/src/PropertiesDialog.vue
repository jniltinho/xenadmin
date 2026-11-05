<script setup lang="ts">
// The VM's Properties dialog (PropertiesDialog): General, CPU and Home Server pages on the left, OK and Cancel.
// shortcut: General has name and description only; no Boot Options, Start Options, Advanced or Custom Fields pages.
import { computed, onMounted, reactive, ref } from "vue";
import { api, type VMProperties } from "./api/workspace";

const props = defineProps<{ server: string; vm: string }>();
const emit = defineEmits<{ (e: "close", saved: boolean): void }>();

const page = ref<"general" | "cpu" | "home">("general");
const orig = ref<VMProperties | null>(null);
const f = reactive({ name: "", description: "", vcpus: 1, atStartup: 1, cores: 1, priority: 4, home: "" });
const error = ref("");
const busy = ref(false);

onMounted(async () => {
  try {
    const p = await api.properties(props.server, props.vm);
    orig.value = p;
    Object.assign(f, { name: p.name, description: p.description, vcpus: p.hotplug ? p.vcpus_max : p.vcpus_at_startup, atStartup: p.vcpus_at_startup,
      cores: p.cores_per_socket, priority: p.priority, home: p.home_server });
  } catch (e) { error.value = String(e); }
});

const choice = computed(() => orig.value?.vcpu_choices.find((c) => c.vcpus === f.vcpus));
// CpuPage.FillTopology: keep the cores per socket when still valid, else the first topology
function vcpusChanged() {
  const t = choice.value?.topologies ?? [];
  if (!t.some((x) => x.cores === f.cores)) f.cores = t[0]?.cores ?? f.cores;
  if (!orig.value?.hotplug || f.atStartup > f.vcpus) f.atStartup = f.vcpus;
}
const warnings = computed(() => {
  const p = orig.value, w: string[] = [];
  if (!p) return w;
  if (f.vcpus < p.min_vcpus || f.atStartup < p.min_vcpus) w.push(`The VM needs at least ${p.min_vcpus} vCPUs.`);
  if (p.max_physical > 0 && p.max_physical < f.vcpus) w.push("More vCPUs than the physical CPUs of any server can slow the VM down.");
  if (f.vcpus % f.cores !== 0) w.push("The number of vCPUs must be a multiple of the number of cores per socket.");
  return w;
});
const pages = [
  { id: "general", label: "General", icon: "000_Info3_h32bit_16" },
  { id: "cpu", label: "CPU", icon: "000_CPU_h32bit_16" },
  { id: "home", label: "Home Server", icon: "000_ServerHome_h32bit_16" },
] as const;

async function save() {
  const p = orig.value!;
  const ch: Record<string, unknown> = {};
  if (f.name.trim() !== p.name) ch.name = f.name.trim();
  if (f.description !== p.description) ch.description = f.description;
  if (p.halted && f.priority !== p.priority) ch.priority = f.priority;
  if (p.halted && f.cores !== p.cores_per_socket) ch.cores_per_socket = f.cores;
  const origVcpus = p.hotplug ? p.vcpus_max : p.vcpus_at_startup;
  if (p.halted && f.vcpus !== origVcpus) { ch.vcpus_max = f.vcpus; ch.vcpus_at_startup = p.hotplug ? f.atStartup : f.vcpus; }
  else if (p.hotplug && f.atStartup !== p.vcpus_at_startup) ch.vcpus_at_startup = f.atStartup;
  if (f.home !== p.home_server) ch.home_server = f.home;
  if (!f.name.trim()) { error.value = "Enter a name for the VM."; page.value = "general"; return; }
  if (!Object.keys(ch).length) { emit("close", false); return; }
  busy.value = true;
  try { await api.saveProperties(props.server, props.vm, ch); emit("close", true); } catch (e) { error.value = String(e); }
  busy.value = false;
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <div class="flex h-[30rem] w-[46rem] flex-col border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">'{{ orig?.name ?? "" }}' Properties</div>
      <div class="flex min-h-0 flex-1">
        <div class="w-48 border-r border-line bg-pane py-1">
          <div v-for="p in pages" :key="p.id" class="flex cursor-default items-center gap-2 px-2 py-1.5"
               :class="page === p.id ? 'bg-select text-white' : 'hover:bg-neutral-100'" @click="page = p.id">
            <img :src="`/icons/${p.icon}.png`" alt="" class="h-4 w-4" />{{ p.label }}
          </div>
        </div>
        <div v-if="orig" class="min-w-0 flex-1 overflow-auto p-4">
          <template v-if="page === 'general'">
            <label class="mb-2 grid grid-cols-[7rem_1fr] items-center gap-2">Name:<input v-model="f.name" class="border border-line px-1 py-0.5" /></label>
            <label class="grid grid-cols-[7rem_1fr] items-start gap-2">Description:<textarea v-model="f.description" rows="4" class="border border-line px-1 py-0.5" /></label>
          </template>
          <template v-else-if="page === 'cpu'">
            <p class="mb-3">Specify the number of vCPUs, their topology, and the priority to assign them over other vCPUs.</p>
            <p v-if="!orig.halted" class="mb-3 flex items-center gap-2"><img src="/icons/000_Info3_h32bit_16.png" alt="" />
              {{ orig.hotplug ? "The maximum number of vCPUs can only be changed when the VM is shut down." : "The vCPUs can only be changed when the VM is shut down." }}</p>
            <div class="grid grid-cols-[11rem_auto] items-center gap-2">
              <span>{{ orig.hotplug ? "Maximum number of vCPUs:" : "Number of vCPUs:" }}</span>
              <select v-model.number="f.vcpus" :disabled="!orig.halted" class="w-24 border border-line bg-white" @change="vcpusChanged">
                <option v-for="c in orig.vcpu_choices" :key="c.vcpus" :value="c.vcpus">{{ c.vcpus }}</option>
              </select>
              <span>Topology:</span>
              <select v-model.number="f.cores" :disabled="!orig.halted" class="w-80 border border-line bg-white">
                <option v-for="t in choice?.topologies ?? []" :key="t.cores" :value="t.cores">{{ t.text }}</option>
              </select>
              <template v-if="orig.hotplug">
                <span>{{ orig.halted ? "Initial number of vCPUs:" : "Current number of vCPUs:" }}</span>
                <select v-model.number="f.atStartup" :disabled="!(orig.halted || orig.running)" class="w-24 border border-line bg-white">
                  <option v-for="n in f.vcpus - (orig.halted ? 0 : orig.vcpus_at_startup - 1)" :key="n" :value="n + (orig.halted ? 0 : orig.vcpus_at_startup - 1)">
                    {{ n + (orig.halted ? 0 : orig.vcpus_at_startup - 1) }}</option>
                </select>
              </template>
            </div>
            <p class="mt-5 mb-1">vCPU priority for this virtual machine:</p>
            <input v-model.number="f.priority" type="range" min="0" max="8" step="1" :disabled="!orig.halted" class="w-[30rem]" />
            <div class="flex w-[30rem] justify-between text-xs text-neutral-500"><span>Lowest</span><span>Normal</span><span>Highest</span></div>
            <p v-for="w in warnings" :key="w" class="mt-3 text-red-800">{{ w }}</p>
          </template>
          <template v-else>
            <p class="mb-3">You can nominate a server that will always run this VM, if it is available. If it is not, the VM is started on another server of the pool.</p>
            <label class="mb-2 flex items-center gap-2" :class="{ 'text-neutral-400': !orig.can_have_none }">
              <input v-model="f.home" type="radio" value="" :disabled="!orig.can_have_none" />Don't assign this VM a home server</label>
            <label class="mb-1 flex items-center gap-2"><input type="radio" :checked="f.home !== ''" @change="f.home = orig.servers.find((s) => s.live)?.uuid ?? ''" />Place the VM on this server:</label>
            <div class="ml-6 h-40 overflow-auto border border-line">
              <div v-for="s in orig.servers" :key="s.uuid" class="flex cursor-default items-center gap-2 px-1 py-0.5"
                   :class="[f.home === s.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100', { 'text-neutral-400': !s.live }]" @click="s.live && (f.home = s.uuid)">
                <img :src="`/icons/${s.live ? '000_TreeConnected_h32bit_16' : '000_ServerDisconnected_h32bit_16'}.png`" alt="" />{{ s.name }}</div>
            </div>
          </template>
        </div>
        <div v-else class="p-4 text-neutral-500">{{ error || "Loading..." }}</div>
      </div>
      <div class="flex items-center justify-end gap-2 border-t border-line p-2">
        <span v-if="error && orig" class="mr-auto text-red-600">{{ error }}</span>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="busy || !orig" @click="save">OK</button>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close', false)">Cancel</button>
      </div>
    </div>
  </div>
</template>
