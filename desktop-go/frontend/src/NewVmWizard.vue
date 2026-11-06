<script setup lang="ts">
// The New VM wizard (NewVMWizard): Template, Name, Installation Media, Home Server, CPU & Memory, Storage, Networking
// and Finish pages; the defaults and the creation come from the Go side (NewVmDefaults, CreateVMAction).
// shortcut: no boot mode page for secure boot choices beyond BIOS/UEFI, no vGPU, vTPM or cloud-config pages.
import { computed, onMounted, reactive, ref } from "vue";
import { api, type NewVMDefaults, type NewVMOptions } from "./api/workspace";

// template: Create VM From Selection opens on that template (its uuid)
const props = defineProps<{ server: string; template?: string }>();
const emit = defineEmits<{ (e: "close", created: string, startError: string): void }>();

const pages = ["Template", "Name", "Installation Media", "Home Server", "CPU & Memory", "Storage", "Networking", "Finish"] as const;
const page = ref(0);
const opts = ref<NewVMOptions | null>(null);
const filter = ref("");
const busy = ref(false);
const error = ref("");
const f = reactive({ template: "", name: "", description: "", install: "cd", iso: "", bootMode: "bios", home: "", vcpus: 1, cores: 1,
  memoryMb: 1024, disks: [] as NewVMDefaults["disks"], networks: [] as NewVMDefaults["networks"], start: true, hvm: true, noDvd: false,
  eliloader: false, url: "" });

onMounted(async () => {
  try {
    opts.value = await api.newVmOptions(props.server);
    f.template = opts.value.templates.find((t) => t.uuid === props.template)?.uuid ?? opts.value.templates[0]?.uuid ?? "";
  } catch (e) { error.value = String(e); }
});
const templates = computed(() => (opts.value?.templates ?? []).filter((t) => t.name.toLowerCase().includes(filter.value.toLowerCase())));
const srName = (uuid: string) => opts.value?.srs.find((s) => s.uuid === uuid)?.text.split("   ")[0] ?? "(no storage)";
const netName = (uuid: string) => opts.value?.networks.find((n) => n.vdi === uuid)?.text ?? uuid;
const gb = (b: number) => +(b / 1024 ** 3).toFixed(2);

// leaving the Template page loads that template's defaults (NewVmDefaults.FromTemplate / DisksFor / VifsFor); leaving
// the Home Server page only places the disks again (the SRs the server can use), keeping every other choice
async function loadDefaults(disksOnly: boolean) {
  const d = await api.newVmDefaults(props.server, f.template, f.home);
  if (disksOnly) { f.disks = d.disks; return; }
  Object.assign(f, { hvm: d.hvm, bootMode: d.boot_mode, install: d.install, vcpus: d.vcpus, cores: d.cores, memoryMb: Math.round(d.memory / 1024 ** 2),
    disks: d.disks, networks: d.networks, noDvd: d.no_dvd, eliloader: d.eliloader, url: "", name: d.name, description: d.description });
}

async function next() {
  error.value = "";
  try {
    if (page.value === 0) { if (!f.template) { error.value = "Select a template."; return; } await loadDefaults(false); }
    if (page.value === 2 && f.eliloader && f.install === "network" && !f.url.trim()) { error.value = "Enter the URL of the installation repository."; return; }
    if (page.value === 1 && !f.name.trim()) { error.value = "Enter a name for the VM."; return; }
    if (page.value === 3) await loadDefaults(true); // the disks' SRs depend on the home server
    if (page.value === 5 && f.disks.some((d) => !d.sr)) { error.value = "Every disk needs a storage repository."; return; }
    page.value++;
  } catch (e) { error.value = String(e); }
}

function addDisk() {
  const used = new Set(f.disks.map((d) => d.position));
  const position = Array.from({ length: 16 }, (_, i) => String(i)).find((p) => !used.has(p) && p !== "3") ?? "1"; // 3: the DVD drive
  f.disks.push({ position, name: `${f.name} ${position}`, description: "", size: 1024 ** 3, sr: opts.value?.default_sr || opts.value?.srs[0]?.uuid || "", bootable: false, inherited: false });
}

async function finish() {
  busy.value = true; error.value = "";
  try {
    const r = await api.createVm(props.server, { template: f.template, name: f.name.trim(), description: f.description, install: f.install,
      iso: f.install === "cd" ? f.iso : "", url: f.url.trim(), boot_mode: f.bootMode, home_server: f.home, vcpus: f.vcpus, cores: f.cores, memory: f.memoryMb * 1024 ** 2,
      disks: f.disks.map((d) => ({ position: d.position, name: d.name, description: d.description, size: d.size, sr: d.sr })), networks: f.networks, start: f.start });
    emit("close", r.uuid, r.start_error ?? ""); // created: never offered again, also when starting it failed
  } catch (e) { error.value = String(e); }
  busy.value = false;
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <div class="flex h-[34rem] w-[54rem] flex-col border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">New VM</div>
      <div class="flex min-h-0 flex-1">
        <ol class="w-48 shrink-0 border-r border-line bg-pane py-2">
          <li v-for="(p, i) in pages" :key="p" class="px-3 py-1" :class="i === page ? 'font-semibold text-select' : i < page ? 'text-neutral-700' : 'text-neutral-400'">{{ p }}</li>
        </ol>
        <div v-if="opts" class="min-w-0 flex-1 overflow-auto p-4">
          <template v-if="page === 0">
            <p class="mb-2">Select a VM template:</p>
            <input v-model="filter" placeholder="Search..." class="mb-2 w-full border border-line px-1 py-0.5" />
            <div class="h-80 overflow-auto border border-line">
              <div v-for="t in templates" :key="t.uuid" class="flex cursor-default items-center gap-2 px-2 py-0.5"
                   :class="f.template === t.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100'" @click="f.template = t.uuid" @dblclick="next">
                {{ t.name }}<span v-if="t.custom" class="text-xs opacity-70">(custom)</span></div>
            </div>
          </template>
          <template v-else-if="page === 1">
            <p class="mb-3">Enter a name that will help you to identify the virtual machine later.</p>
            <label class="mb-2 grid grid-cols-[7rem_1fr] items-center gap-2">Name:<input v-model="f.name" class="border border-line px-1 py-0.5" /></label>
            <label class="grid grid-cols-[7rem_1fr] items-start gap-2">Description:<textarea v-model="f.description" rows="4" class="border border-line px-1 py-0.5" /></label>
          </template>
          <template v-else-if="page === 2">
            <p class="mb-3">Select the installation method, operating system software you want to install on the new VM.</p>
            <p v-if="f.noDvd" class="mb-3 text-neutral-600">This template has no DVD drive: the new VM gets none.</p>
            <label class="mb-2 flex items-center gap-2"><input v-model="f.install" type="radio" value="cd" :disabled="f.noDvd" />Install from ISO library or DVD drive:</label>
            <select v-model="f.iso" :disabled="f.install !== 'cd'" class="mb-4 ml-6 w-[30rem] border border-line bg-white">
              <option value="">&lt;empty&gt;</option>
              <option v-for="i in opts.isos" :key="i.vdi" :value="i.vdi">{{ i.text }}</option>
            </select>
            <label class="flex items-center gap-2"><input v-model="f.install" type="radio" value="network" :disabled="f.noDvd" />{{ f.eliloader ? "Install from URL:" : "Boot from network" }}</label>
            <input v-if="f.eliloader" v-model="f.url" :disabled="f.install !== 'network'" placeholder="http://, ftp:// or nfs://" class="ml-6 mt-1 w-[30rem] border border-line px-1 py-0.5" />
            <template v-if="f.hvm">
              <p class="mt-5 mb-1">Boot mode:</p>
              <label class="mr-4"><input v-model="f.bootMode" type="radio" value="bios" /> BIOS Boot</label>
              <label class="mr-4"><input v-model="f.bootMode" type="radio" value="uefi" /> UEFI Boot</label>
              <label><input v-model="f.bootMode" type="radio" value="secureuefi" /> UEFI Secure Boot</label>
            </template>
          </template>
          <template v-else-if="page === 3">
            <p class="mb-3">Select a home server. The VM will run on it when it is available.</p>
            <label class="mb-2 flex items-center gap-2"><input v-model="f.home" type="radio" value="" />Don't assign this VM a home server</label>
            <div class="h-56 overflow-auto border border-line">
              <div v-for="s in opts.servers" :key="s.uuid" class="flex cursor-default items-center gap-2 px-1 py-0.5"
                   :class="[f.home === s.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100', { 'text-neutral-400': !s.live }]" @click="s.live && (f.home = s.uuid)">
                <img :src="`/icons/${s.live ? '000_TreeConnected_h32bit_16' : '000_ServerDisconnected_h32bit_16'}.png`" alt="" />{{ s.name }}</div>
            </div>
          </template>
          <template v-else-if="page === 4">
            <p class="mb-3">Specify the number of virtual CPUs, their topology and the amount of memory for the new VM.</p>
            <div class="grid grid-cols-[11rem_auto] items-center gap-2">
              <span>Number of vCPUs:</span><input v-model.number="f.vcpus" type="number" min="1" max="64" class="w-24 border border-line px-1" />
              <span>Cores per socket:</span><input v-model.number="f.cores" type="number" min="1" class="w-24 border border-line px-1" />
              <span>Memory (MB):</span><input v-model.number="f.memoryMb" type="number" min="128" step="128" class="w-32 border border-line px-1" />
            </div>
            <p v-if="f.vcpus % f.cores" class="mt-3 text-red-800">The number of vCPUs must be a multiple of the number of cores per socket.</p>
          </template>
          <template v-else-if="page === 5">
            <p class="mb-3">The virtual disks of the new VM. Their storage repository and size can be changed.</p>
            <table class="w-full border border-line">
              <thead class="bg-pane text-left"><tr><th class="px-2 py-1 font-normal">Position</th><th class="px-2 py-1 font-normal">Name</th>
                <th class="px-2 py-1 font-normal">Size (GB)</th><th class="px-2 py-1 font-normal">Location</th><th></th></tr></thead>
              <tbody>
                <tr v-for="(d, i) in f.disks" :key="i" class="border-t border-line">
                  <td class="px-2 py-1">{{ d.position }}</td>
                  <td class="px-2 py-1"><input v-model="d.name" class="w-40 border border-line px-1" /></td>
                  <td class="px-2 py-1"><input :value="gb(d.size)" type="number" min="0.1" step="any" :disabled="d.inherited" class="w-20 border border-line px-1" @input="d.size = Math.round(Number(($event.target as HTMLInputElement).value) * 1024 ** 3)" /></td>
                  <td class="px-2 py-1"><select v-model="d.sr" class="w-56 border border-line bg-white"><option v-for="s in opts.srs" :key="s.uuid" :value="s.uuid">{{ s.text }}</option></select></td>
                  <td class="px-2 py-1"><button class="btn" :disabled="d.inherited" :title="d.inherited ? 'A disk of the template: its storage can be changed' : ''" @click="f.disks.splice(i, 1)">Delete</button></td>
                </tr>
              </tbody>
            </table>
            <button class="btn mt-2" @click="addDisk">Add...</button>
          </template>
          <template v-else-if="page === 6">
            <p class="mb-3">The virtual network interfaces of the new VM (one per network).</p>
            <div v-for="(n, i) in f.networks" :key="i" class="mb-1 flex items-center gap-2">
              <span class="w-24 whitespace-nowrap">Interface {{ i }}</span>
              <select v-model="n.network" class="w-72 border border-line bg-white"><option v-for="o in opts.networks" :key="o.vdi" :value="o.vdi">{{ o.text }}</option></select>
              <button class="btn" @click="f.networks.splice(i, 1)">Delete</button>
            </div>
            <button class="btn mt-2" :disabled="f.networks.length >= 4 || !opts.networks.length" @click="f.networks.push({ network: opts.networks[0].vdi })">Add</button>
          </template>
          <template v-else>
            <p class="mb-3">All the necessary information has been collected and the wizard is ready to create the new VM.</p>
            <table class="w-full border border-line">
              <tbody>
                <tr v-for="[k, v] in [['Template', opts.templates.find((t) => t.uuid === f.template)?.name ?? ''], ['Name', f.name],
                  ['Install method', f.install === 'cd' ? (opts.isos.find((i) => i.vdi === f.iso)?.text ?? 'DVD drive (empty)') : f.install === 'network' ? (f.eliloader ? f.url : 'Network boot') : 'None'],
                  ['Home server', opts.servers.find((s) => s.uuid === f.home)?.name ?? 'None'], ['vCPUs', String(f.vcpus)], ['Memory', `${f.memoryMb} MB`],
                  ['Disks', f.disks.map((d) => `${d.name} (${gb(d.size)} GB on ${srName(d.sr)})`).join(', ') || 'None'],
                  ['Networks', f.networks.map((n) => netName(n.network)).join(', ') || 'None']]" :key="k" class="border-b border-line">
                  <td class="w-40 px-2 py-1 text-neutral-600">{{ k }}:</td><td class="px-2 py-1">{{ v }}</td></tr>
              </tbody>
            </table>
            <label class="mt-3 flex items-center gap-2"><input v-model="f.start" type="checkbox" />Start the new VM automatically</label>
          </template>
        </div>
        <div v-else class="p-4 text-neutral-500">{{ error || "Loading..." }}</div>
      </div>
      <div class="flex items-center justify-end gap-2 border-t border-line p-2">
        <span v-if="error && opts" class="mr-auto text-red-600">{{ error }}</span>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="page === 0 || busy" @click="page--">&lt; Previous</button>
        <button v-if="page < pages.length - 1" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="!opts" @click="next">Next &gt;</button>
        <button v-else class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="busy" @click="finish">{{ busy ? "Creating..." : "Create Now" }}</button>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="busy" @click="emit('close', '', '')">Cancel</button>
      </div>
    </div>
  </div>
</template>
