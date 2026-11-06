<script setup lang="ts">
// The VM's Properties dialog (PropertiesDialog): General (name, description, folder, tags), Custom Fields, CPU, Boot
// Options, Start Options, Alerts (PerfmonAlertsPage), Home Server, GPU (GpuEditPage), USB (UsbEditPage) and (HVM)
// Advanced Options pages on the left, OK and Cancel. The GPU and USB pages are shown when they say so.
import { computed, onMounted, reactive, ref } from "vue";
import { api, type ObjectProps, type VMProperties } from "./api/workspace";
import FolderTags from "./FolderTags.vue";
import GpuEditPage from "./GpuEditPage.vue";
import PerfmonAlertsPage from "./PerfmonAlertsPage.vue";
import UsbEditPage from "./UsbEditPage.vue";

const props = defineProps<{ server: string; vm: string }>();
const emit = defineEmits<{ (e: "close", saved: boolean): void }>();

type Page = "general" | "fields" | "cpu" | "boot" | "start" | "alerts" | "home" | "gpu" | "usb" | "advanced";
const alerts = ref<InstanceType<typeof PerfmonAlertsPage> | null>(null);
const gpu = ref<InstanceType<typeof GpuEditPage> | null>(null);
const gpuShown = ref(false);
// VMHAEditPage.RefillPrioritiesComboBox: no Restart priority for a VM with vGPUs on the GPU page (saved or not)
const haChoices = computed(() => (orig.value?.ha_priorities ?? []).filter((p) => !(p.value === "restart" && (gpu.value?.count ?? 0) > 0)));
const usbShown = ref(false);
const page = ref<Page>("general");
const orig = ref<VMProperties | null>(null);
const f = reactive({ name: "", description: "", vcpus: 1, atStartup: 1, cores: 1, priority: 4, home: "",
  autoPowerOn: false, devices: [] as { code: string; on: boolean }[], picked: -1, pvArgs: "", pvFromCd: false, shadow: "1",
  shadowMode: "general" as "general" | "apps" | "manual", order: 0, delay: 0, haPriority: "" });
const error = ref("");
const busy = ref(false);
// GeneralEditPage's Folder and Tags (the same pages as a server's or an SR's, ObjectPropertiesDialog)
const org = ref<ObjectProps | null>(null);
const ft = reactive({ folder: "", tags: [] as string[] });

// Custom Fields (CustomFieldsPage): the pool's fields with the VM's values; a Date in the local time of a
// datetime-local input, stored as RFC 3339 universal time
type Field = { name: string; type: string; value: string };
const fields = ref<Field[]>([]);
// null-prototype dictionaries: a field may be named "constructor" or "__proto__"
const fieldEdits = reactive<Record<string, string>>(Object.create(null));
const fieldTypes: Record<string, string> = Object.create(null);
const toLocal = (rfc: string) => {
  if (!rfc) return "";
  const d = new Date(rfc), p = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}T${p(d.getHours())}:${p(d.getMinutes())}`;
};
const shown = (f: Field) => (f.type === "Date" ? toLocal(f.value) : f.value);
async function loadFields() {
  fields.value = (await api.customFields(props.server, props.vm)) ?? [];
  // edits survive a reload; a deleted field (or one made again with another type) starts from the VM's value
  for (const name of Object.keys(fieldEdits))
    if (!fields.value.some((f) => f.name === name && f.type === fieldTypes[name])) { delete fieldEdits[name]; delete fieldTypes[name]; }
  for (const f of fields.value) if (!Object.hasOwn(fieldEdits, f.name)) { fieldEdits[f.name] = shown(f); fieldTypes[f.name] = f.type; }
}
function fieldChanges(): Record<string, string> | null {
  const out: Record<string, string> = Object.create(null);
  for (const f of fields.value) {
    const v = fieldEdits[f.name] ?? "";
    if (v === shown(f)) continue;
    if (f.type === "Date" && v && Number.isNaN(new Date(v).getTime())) return null;
    out[f.name] = f.type === "Date" && v ? new Date(v).toISOString().replace(/\.\d{3}Z$/, "Z") : v;
  }
  return out;
}
// CustomFieldsDialog: add a field or delete one, for every object of the pool (at once, as the C# actions)
const editFields = ref<{ name: string; type: string; error: string; confirm: string } | null>(null);
async function fieldAction(f: () => Promise<void>) {
  const ed = editFields.value!; // the dialog may be closed before the call ends: the page is read again anyway
  try { await f(); ed.error = ""; } catch (e) { ed.error = String(e); }
  try { await loadFields(); } catch (e) { error.value = String(e); }
}

onMounted(async () => {
  try {
    const p = await api.properties(props.server, props.vm);
    orig.value = p;
    Object.assign(f, { name: p.name, description: p.description, vcpus: p.hotplug ? p.vcpus_max : p.vcpus_at_startup, atStartup: p.vcpus_at_startup,
      cores: p.cores_per_socket, priority: p.priority, home: p.home_server, autoPowerOn: p.auto_power_on, pvArgs: p.pv_args,
      pvFromCd: p.pv_boot_from_cd && p.has_cd, shadow: String(p.shadow_multiplier), order: p.order, delay: p.start_delay, haPriority: p.ha_restart_priority,
      shadowMode: p.shadow_multiplier === 1 ? "general" : p.shadow_multiplier === 4 ? "apps" : "manual",
      // BootOptionsEditPage: the boot order first, checked, then the devices it leaves out
      devices: [...p.boot_order.split(""), ..."CDN".split("").filter((c) => !p.boot_order.includes(c))].map((code) => ({ code, on: p.boot_order.includes(code) })) });
    await loadFields();
    org.value = await api.objectProperties(props.server, "vms", props.vm);
    Object.assign(ft, { folder: org.value.folder, tags: [...org.value.tags] });
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
const pages = computed<{ id: Page; label: string; icon: string }[]>(() => [
  { id: "general", label: "General", icon: "000_Info3_h32bit_16" },
  { id: "fields", label: "Custom Fields", icon: "000_Fields_h32bit_16" },
  { id: "cpu", label: "CPU", icon: "000_CPU_h32bit_16" },
  { id: "boot", label: "Boot Options", icon: "001_PowerOn_h32bit_16" },
  { id: "start", label: "Start Options", icon: "000_RebootVM_h32bit_16" },
  { id: "alerts", label: "Alerts", icon: "000_Alert2_h32bit_16" },
  { id: "home", label: "Home Server", icon: "000_ServerHome_h32bit_16" },
  ...(gpuShown.value ? [{ id: "gpu" as const, label: "GPU", icon: "000_GetMemoryInfo_h32bit_16" }] : []),
  ...(usbShown.value ? [{ id: "usb" as const, label: "USB", icon: "usb_16" }] : []),
  ...(orig.value?.hvm ? [{ id: "advanced" as const, label: "Advanced", icon: "002_Configure_h32bit_16" }] : []), // VMAdvancedEditPage: HVM only
]);
const bootOrder = () => f.devices.filter((d) => d.on).map((d) => d.code).join("");
function move(delta: number) {
  const i = f.picked, j = i + delta;
  if (i < 0 || j < 0 || j >= f.devices.length) return;
  [f.devices[i], f.devices[j]] = [f.devices[j], f.devices[i]];
  f.picked = j;
}
// VMAdvancedEditPage: general use (1), Virtual Apps (4) or a manual multiplier of at least 1; the mode is chosen by
// its radio button, the manual value typed (never changed by the mode)
const shadowValue = () => (f.shadowMode === "general" ? 1 : f.shadowMode === "apps" ? 4 : Number(f.shadow));

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
  if (f.autoPowerOn !== p.auto_power_on) ch.auto_power_on = f.autoPowerOn;
  if (p.hvm && bootOrder() !== p.boot_order) ch.boot_order = bootOrder();
  if (!p.hvm && f.pvArgs !== p.pv_args) ch.pv_args = f.pvArgs;
  if (!p.hvm && f.pvFromCd !== (p.pv_boot_from_cd && p.has_cd)) ch.pv_boot_from_cd = f.pvFromCd;
  const shadow = shadowValue();
  if (p.hvm && shadow !== p.shadow_multiplier) {
    if (!(Number.isFinite(shadow) && shadow >= 1)) { error.value = "Value should be a number greater than or equal to 1.0"; page.value = "advanced"; return; }
    ch.shadow_multiplier = shadow;
  }
  // StartOptionsPage: with HA on, saved with the restart priority (SetHaPrioritiesAction, the pool's failures to tolerate kept)
  const haVm = p.ha_priorities.length > 0; // a VM in the HA plan (not a template)
  if (f.order !== p.order || f.delay !== p.start_delay || haVm && f.haPriority !== p.ha_restart_priority) {
    if (![f.order, f.delay].every((n) => Number.isInteger(n) && n >= 0)) { error.value = "Enter a whole number of 0 or more."; page.value = "start"; return; }
    ch.order = f.order; ch.start_delay = f.delay;
    if (haVm) ch.ha_restart_priority = f.haPriority;
  }
  if (!f.name.trim()) { error.value = "Enter a name for the VM."; page.value = "general"; return; }
  const fc = fieldChanges();
  if (!fc) { error.value = "Enter a valid date."; page.value = "fields"; return; }
  const og: Record<string, unknown> = {};
  if (org.value && ft.folder !== org.value.folder) og.folder = ft.folder;
  if (org.value && ft.tags.join("\n") !== [...org.value.tags].sort().join("\n")) { og.tags = ft.tags; og.tags_set = true; }
  const pm = alerts.value?.changes() ?? null;
  if (typeof pm === "string") { error.value = pm; page.value = "alerts"; return; }
  const gc = gpu.value?.changes() ?? null;
  if (haVm && f.haPriority === "restart" && (gpu.value?.count ?? 0) > 0 && (gc || f.haPriority !== p.ha_restart_priority)) {
    error.value = "A VM with virtual GPUs cannot have the Restart priority."; page.value = "start"; return;
  }
  if (!Object.keys(ch).length && !Object.keys(fc).length && !Object.keys(og).length && !pm && !gc) { emit("close", false); return; }
  busy.value = true;
  try {
    if (Object.keys(ch).length) await api.saveProperties(props.server, props.vm, ch);
    if (Object.keys(fc).length) await api.saveCustomFields(props.server, props.vm, fc);
    if (Object.keys(og).length) await api.saveObjectProperties(props.server, "vms", props.vm, og);
    if (pm) await api.savePerfmon(props.server, "vms", props.vm, pm.on, pm.interval);
    // after the HA restart priority (saveProperties), which decides whether the vGPUs can change
    if (gc) await api.assignVgpus(props.server, props.vm, gc.removed, gc.add).catch(async (e) => { await gpu.value?.reload(); throw e; });
    emit("close", true);
  } catch (e) { error.value = String(e); }
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
          <PerfmonAlertsPage v-show="page === 'alerts'" ref="alerts" :server="server" kind="vms" :uuid="vm" :disabled="busy" />
          <GpuEditPage v-show="page === 'gpu'" ref="gpu" :server="server" :vm="vm" :ha-priority="f.haPriority" :disabled="busy" @shown="gpuShown = $event" />
          <UsbEditPage v-show="page === 'usb'" :server="server" :vm="vm" :ha-priority="f.haPriority" :disabled="busy" @shown="usbShown = $event" />
          <template v-if="page === 'alerts' || page === 'gpu' || page === 'usb'" />
          <template v-else-if="page === 'general'">
            <label class="mb-2 grid grid-cols-[7rem_1fr] items-center gap-2">Name:<input v-model="f.name" class="border border-line px-1 py-0.5" /></label>
            <label class="grid grid-cols-[7rem_1fr] items-start gap-2">Description:<textarea v-model="f.description" rows="4" class="border border-line px-1 py-0.5" /></label>
            <div v-if="org" class="mt-2 grid grid-cols-[7rem_1fr] items-start gap-2">
              <FolderTags v-model:folder="ft.folder" v-model:tags="ft.tags" :folders="org.folders" :all-tags="org.all_tags" />
            </div>
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
          <template v-else-if="page === 'fields'">
            <p class="mb-3">Custom fields hold your own information about this object, such as its owner or its expiry date.</p>
            <div class="mb-3 grid grid-cols-[10rem_1fr] items-center gap-2">
              <template v-for="fd in fields" :key="fd.name">
                <label :for="'cf-' + fd.name" class="truncate">{{ fd.name }}:</label>
                <input :id="'cf-' + fd.name" v-model="fieldEdits[fd.name]" :type="fd.type === 'Date' ? 'datetime-local' : 'text'" class="border border-line px-1 py-0.5" />
              </template>
            </div>
            <p v-if="!fields.length" class="mb-3 text-neutral-500">There are no custom fields. Use Edit Custom Fields to add one.</p>
            <button class="btn" @click="editFields = { name: '', type: 'String', error: '', confirm: '' }">Edit Custom Fields...</button>
          </template>
          <template v-else-if="page === 'start'">
            <p class="mb-3">{{ orig.ha_priorities.length ? "HA is enabled on this pool: set this VM's restart priority, its start order and the delay before the next VM starts."
              : "Set the order in which this VM starts in a vApp or when HA restarts it, and the delay before the next VM starts." }}</p>
            <div class="grid grid-cols-[auto_14rem_auto] items-center gap-x-3 gap-y-2">
              <label for="order">Start order:</label><input id="order" v-model.number="f.order" type="number" min="0" step="1" class="w-32 border border-line px-1 py-0.5" /><span />
              <label for="delay">Attempt to start next VM after:</label><input id="delay" v-model.number="f.delay" type="number" min="0" step="1" class="w-32 border border-line px-1 py-0.5" /><span>seconds</span>
              <template v-if="orig.ha_priorities.length">
                <label for="harestart">HA restart priority:</label>
                <select id="harestart" v-model="f.haPriority" class="border border-line bg-white px-1 py-0.5">
                  <option v-for="p in haChoices" :key="p.value" :value="p.value">{{ p.text }}</option>
                </select><span />
              </template>
            </div>
          </template>
          <template v-else-if="page === 'boot'">
            <p class="mb-3">{{ orig.hvm ? "Choose the devices from which you would like to boot the VM’s operating system, and configure the boot device order. The first item in the list will be scanned first when booting; if boot media is not found, the second device will then be scanned, and so on down the list of selected devices."
              : "Choose the device from which you would like to boot the VM’s operating system, and, optionally, specify any additional boot parameters." }}</p>
            <fieldset class="mb-4 border border-line px-3 pb-3 pt-1"><legend class="px-1">Autostart</legend>
              <label class="flex items-center gap-2"><input v-model="f.autoPowerOn" type="checkbox" />Autostart on server boot</label></fieldset>
            <div v-if="orig.hvm" class="flex items-start gap-3">
              <span class="w-32">Boot order:</span>
              <div class="h-24 w-40 overflow-auto border border-neutral-400">
                <!-- a row is selected by a click on it; only its check box turns the device on or off -->
                <div v-for="(d, i) in f.devices" :key="d.code" class="flex cursor-default items-center gap-2 px-1"
                     :class="f.picked === i ? 'bg-select text-white' : ''" @click="f.picked = i">
                  <input v-model="d.on" type="checkbox" @click.stop="f.picked = i" />{{ orig.boot_devices[d.code] }}</div>
              </div>
              <div class="flex flex-col gap-1">
                <button class="btn min-w-24" :disabled="f.picked <= 0" @click="move(-1)">Move Up</button>
                <button class="btn min-w-24" :disabled="f.picked < 0 || f.picked >= f.devices.length - 1" @click="move(1)">Move Down</button>
              </div>
            </div>
            <div v-else class="grid grid-cols-[10rem_1fr] items-center gap-2">
              <span>Boot device:</span>
              <select v-model="f.pvFromCd" class="w-48 border border-line bg-white">
                <option :value="false">{{ orig.boot_devices.C }}</option><option v-if="orig.has_cd" :value="true">DVD drive</option>
              </select>
              <span>OS boot parameters:</span><input v-model="f.pvArgs" class="border border-line px-1 py-0.5" />
            </div>
          </template>
          <template v-else-if="page === 'advanced'">
            <p class="mb-3">Use these options to optimize this VM for increased performance in different environments.</p>
            <label class="mb-1 flex items-center gap-2"><input v-model="f.shadowMode" type="radio" value="general" :disabled="orig.frozen" />Optimize for general use</label>
            <label class="mb-1 flex items-center gap-2"><input v-model="f.shadowMode" type="radio" value="apps" :disabled="orig.frozen" />Optimize for Virtual Apps</label>
            <label class="mb-1 flex items-center gap-2"><input v-model="f.shadowMode" type="radio" value="manual" :disabled="orig.frozen" />Optimize manually (advanced use only)</label>
            <label class="ml-6 flex items-center gap-2">Shadow memory multiplier:<input v-model="f.shadow" :disabled="orig.frozen || f.shadowMode !== 'manual'" class="w-24 border border-line px-1 py-0.5" /></label>
            <p v-if="orig.frozen" class="mt-3 text-amber-700">These settings cannot be adjusted when the VM is suspended</p>
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
    <div v-if="editFields" class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
      <div class="w-[30rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Edit Custom Fields</div>
        <div v-if="editFields.confirm" class="space-y-3 p-3">
          <p>Delete the custom field '{{ editFields.confirm }}' and its values on every object?</p>
          <div class="flex justify-end gap-2">
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="fieldAction(() => api.removeCustomField(server, editFields!.confirm)); editFields!.confirm = ''">Yes</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="editFields!.confirm = ''">No</button>
          </div>
        </div>
        <div v-else class="space-y-3 p-3">
          <div class="h-40 overflow-auto border border-line">
            <div v-for="fd in fields" :key="fd.name" class="flex items-center justify-between px-2 py-1 hover:bg-neutral-100">
              <span>{{ fd.name }}&nbsp;&nbsp;&nbsp;({{ fd.type === "Date" ? "Date and time" : "Text" }})</span>
              <button class="btn" @click="editFields!.confirm = fd.name">Delete</button>
            </div>
          </div>
          <form novalidate class="flex items-center gap-2" @submit.prevent="fieldAction(() => api.addCustomField(server, editFields!.name, editFields!.type)); editFields!.name = ''">
            <input v-model="editFields.name" placeholder="Name of the new field" aria-label="Name of the new field" class="w-56 border border-line px-1 py-0.5" />
            <select v-model="editFields.type" aria-label="Type" class="border border-line bg-white px-1 py-0.5"><option value="String">Text</option><option value="Date">Date and time</option></select>
            <button type="submit" class="btn" :disabled="!editFields.name.trim()">Add</button>
          </form>
          <p v-if="editFields.error" class="text-red-600">{{ editFields.error }}</p>
          <div class="flex justify-end"><button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="editFields = null">Close</button></div>
        </div>
      </div>
    </div>
  </div>
</template>
