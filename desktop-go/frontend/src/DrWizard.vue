<script setup lang="ts">
// Pool > Disaster Recovery > Disaster Recovery Wizard (DRFailoverWizard, XenAdmin.Desktop/DrWindows.cs DrWizard):
// Before You Start (failover, failback or test failover), Storage (the SRs with metadata), VMs and vApps (from each
// pool's metadata, and what to do with them after recovery), Prechecks, Summary. Finish recovers them; a test failover
// then offers to remove what it recovered. The Go side checks everything again (gui.DRRecover runs the prechecks).
import { computed, ref, watch } from "vue";
import { api, type DRCheck, type DRChoice, type DRPool, type DRSR } from "./api/workspace";

const props = defineProps<{ server: string; pool: string }>();
const emit = defineEmits<{ (e: "close"): void; (e: "run", label: string, f: () => Promise<void>): void }>();

type Type = "failover" | "failback" | "test";
const type = ref<Type>("failover");
const verb = computed(() => (type.value === "failback" ? "failback" : type.value === "test" ? "test failover" : "failover"));
const pages = ["Before You Start", "Storage", "VMs and vApps", "Prechecks", "Summary"];
const titles = computed(() => ["Disaster Recovery Wizard", "Locate the mirrored storage", "Select the vApps and VMs to recover",
  "Check that the selected vApps and VMs can be recovered", `Ready to start the ${verb.value}`]);
const page = ref(0);
const error = ref("");
const pending = ref(false);

// Storage: the SRs that can hold metadata, those with metadata ticked
const srs = ref<DRSR[] | null>(null);
const chosenSrs = ref<string[]>([]);
async function loadSrs() {
  try {
    const list = (await api.drSrs(props.server)).srs;
    srs.value = list;
    chosenSrs.value = list.filter((s) => s.has_metadata).map((s) => s.uuid);
  } catch (e) { srs.value = []; error.value = String(e); }
}

// VMs and vApps: the metadata of the chosen SRs, read again when they change; a newer read makes an older one stale
const metadata = ref<DRPool[] | null>(null);
let metadataFor = "";
let metadataGeneration = 0;
const chosen = ref<string[]>([]); // `${pool}/${uuid}` of the vApps and VMs
async function loadMetadata() {
  const key = [...chosenSrs.value].sort().join(",");
  if (metadata.value && metadataFor === key) return;
  metadata.value = null;
  chosen.value = [];
  const generation = ++metadataGeneration;
  let result: DRPool[] = [];
  try { result = await api.drMetadata(props.server, [...chosenSrs.value]); } catch (e) { error.value = String(e); }
  if (generation !== metadataGeneration) return;
  metadata.value = result;
  metadataFor = key;
}
// DRFailoverWizardAppliancesPage: a test failover never boots the recovered guests
const modes = computed(() => (type.value === "test"
  ? [{ value: "paused", text: "Start them paused" }, { value: "none", text: "Do not start them" }]
  : [{ value: "start", text: "Start the recovered vApps and VMs" }, { value: "paused", text: "Start them paused" }, { value: "none", text: "Do not start them" }]));
const start = ref("start");
watch(modes, (m) => { if (!m.some((x) => x.value === start.value)) start.value = m[0].value; }, { immediate: true });

const choices = computed<DRChoice[]>(() => (metadata.value ?? []).filter((p) => !p.current).map((p) => ({
  vdi: p.vdi,
  vapps: p.vapps.filter((a) => chosen.value.includes(`${p.uuid}/${a.uuid}`)).map((a) => a.uuid),
  vms: p.vms.filter((v) => chosen.value.includes(`${p.uuid}/${v.uuid}`)).map((v) => v.uuid),
})).filter((c) => c.vapps.length + c.vms.length > 0));
const chosenCount = computed(() => choices.value.reduce((n, c) => n + c.vapps.length + c.vms.length, 0));

// Prechecks: run again each time the page is shown (the user may have fixed what failed)
const checks = ref<DRCheck[] | null>(null);
let checksGeneration = 0;
async function runChecks() {
  checks.value = null;
  const generation = ++checksGeneration;
  let result: DRCheck[];
  try { result = await api.drPrechecks(props.server, type.value, choices.value); } catch (e) { result = [{ text: String(e) || "The prechecks failed.", error: true }]; }
  if (generation === checksGeneration) checks.value = result;
}

watch(page, (p) => {
  error.value = "";
  if (p === 1 && srs.value === null) loadSrs();
  if (p === 2) loadMetadata();
  if (p === 3) runChecks();
});

// DrWizard.Validate
const why = computed(() => {
  switch (page.value) {
    case 1: return srs.value === null ? "Looking for the SRs..." : !chosenSrs.value.length ? "Select the SRs holding the metadata of the pool to recover." : "";
    case 2: return metadata.value === null ? "Loading the metadata..." : !chosenCount.value ? "Select the vApps and VMs to recover." : "";
    case 3: return checks.value === null ? "Running the prechecks..." : checks.value.some((c) => c.error) ? "Fix the errors above before you continue." : "";
  }
  return "";
});

const summary = computed(() => (metadata.value ?? []).filter((p) => !p.current).map((p) => ({
  name: p.name,
  items: [...p.vapps.filter((a) => chosen.value.includes(`${p.uuid}/${a.uuid}`)).map((a) => `vApp '${a.name}'`),
    ...p.vms.filter((v) => chosen.value.includes(`${p.uuid}/${v.uuid}`)).map((v) => v.name)],
})).filter((p) => p.items.length));

// Finish: a failover or failback runs in the background; a test failover waits here to offer its clean-up
const testDone = ref<{ error: string; pending: number } | null>(null);
async function finish() {
  if (pending.value || !chosenCount.value) return;
  pending.value = true;
  const t = type.value, c = choices.value, s = start.value;
  if (t !== "test") {
    emit("run", "Disaster recovery", () => api.drRecover(props.server, t, c, s));
    emit("close");
    return;
  }
  let failed = "";
  try { await api.drRecover(props.server, t, c, s); } catch (e) { failed = String(e); }
  let left = 0;
  try { left = await api.drTestPending(props.server); } catch (e) { failed ||= String(e); }
  testDone.value = { error: failed, pending: left };
}
function cleanup() {
  emit("run", "Clean up the test failover", () => api.drTestCleanup(props.server));
  emit("close");
}
const btn = "min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400";
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <!-- the test failover is over: DRFailoverWizard.FinishWizard's question -->
    <div v-if="testDone" class="w-[30rem] border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">Test Failover</div>
      <div class="space-y-3 p-3">
        <p v-if="testDone.error" class="whitespace-pre-line text-red-600">{{ testDone.error }}</p>
        <p v-if="testDone.pending">The test failover is done. Remove the vApps and VMs it recovered?</p>
        <div class="flex justify-end gap-2">
          <template v-if="testDone.pending">
            <button :class="btn" @click="cleanup">Yes</button>
            <button :class="btn" autofocus @click="emit('close')">No</button>
          </template>
          <button v-else :class="btn" @click="emit('close')">OK</button>
        </div>
      </div>
    </div>
    <div v-else class="flex h-[34rem] w-[52rem] flex-col border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">Disaster Recovery - {{ pool }}</div>
      <div class="border-b border-line px-4 py-2 text-base font-semibold">{{ titles[page] }}</div>
      <div class="flex min-h-0 flex-1">
        <ol class="w-40 shrink-0 border-r border-line bg-pane py-2">
          <li v-for="(p, i) in pages" :key="p" class="px-3 py-1" :class="i === page ? 'font-semibold text-select' : i < page ? 'text-neutral-700' : 'text-neutral-400'">{{ p }}</li>
        </ol>
        <div class="min-w-0 flex-1 overflow-auto p-4">
          <div v-if="page === 0" class="space-y-2">
            <p>This wizard will take you through the steps needed to fail over your critical vApps and VMs in the event of a site/datacenter failure at your primary data site, to fail them back, or to test a failover.</p>
            <p class="text-neutral-700">Configuration information for the vApps and VMs to be recovered will be retrieved from remote storage mirrors. Before the recovery begins, the wizard will run a number of checks to ensure that the selected vApps and VMs can be brought back on this pool.</p>
            <label class="mt-3 flex items-center gap-2"><input v-model="type" type="radio" value="failover" :disabled="pending" />Failover</label>
            <p class="ml-6 text-neutral-700">Recover the vApps and VMs of a failed primary site on this (secondary) pool.</p>
            <label class="flex items-center gap-2"><input v-model="type" type="radio" value="failback" :disabled="pending" />Failback</label>
            <p class="ml-6 text-neutral-700">Bring the vApps and VMs back to the primary site once it has been restored.</p>
            <label class="flex items-center gap-2"><input v-model="type" type="radio" value="test" :disabled="pending" />Test Failover</label>
            <p class="ml-6 text-neutral-700">Recover the vApps and VMs here to check the failover works; they are removed when the wizard ends.</p>
            <p class="mt-3 font-semibold">IMPORTANT</p>
            <p>Before failover can take place, storage mirroring - remote replication of your VM data - must be halted. The mirrored SRs must be attached to this pool.</p>
          </div>
          <div v-else-if="page === 1" class="space-y-3">
            <p>Select the SRs holding the VM configuration data (metadata) of the pool you want to recover:</p>
            <div class="h-64 overflow-auto border border-neutral-400 p-1">
              <p v-if="srs === null" class="text-neutral-600">Looking for the SRs...</p>
              <p v-else-if="!srs.length">No SR attached to this pool can hold VM metadata. Attach the mirrored SR with New Storage first.</p>
              <label v-for="s in srs ?? []" :key="s.uuid" class="flex items-center gap-2 px-1 py-0.5">
                <input v-model="chosenSrs" type="checkbox" :value="s.uuid" />{{ s.name }}<span class="text-neutral-600">({{ s.type }})</span>
                <span v-if="s.has_metadata" class="text-neutral-600">- metadata detected</span>
              </label>
            </div>
          </div>
          <div v-else-if="page === 2" class="space-y-3">
            <p>Select the vApps and VMs to recover. A vApp brings its own VMs.</p>
            <div class="h-56 overflow-auto border border-neutral-400 p-1">
              <p v-if="metadata === null" class="text-neutral-600">Loading the metadata...</p>
              <p v-else-if="!metadata.length">No pool metadata was found on the selected SRs.</p>
              <template v-for="p in metadata ?? []" :key="p.uuid">
                <p class="mt-1 font-semibold">{{ p.name }}<span v-if="p.note" class="ml-2 font-normal text-neutral-600">({{ p.note }})</span></p>
                <label v-for="a in p.vapps" :key="a.uuid" class="ml-5 flex items-center gap-2" :class="p.current ? 'text-neutral-400' : ''">
                  <input v-model="chosen" type="checkbox" :value="`${p.uuid}/${a.uuid}`" :disabled="p.current" />vApp '{{ a.name }}' ({{ a.vms ?? 0 }} VMs)
                </label>
                <label v-for="v in p.vms" :key="v.uuid" class="ml-5 flex items-center gap-2" :class="p.current ? 'text-neutral-400' : ''">
                  <input v-model="chosen" type="checkbox" :value="`${p.uuid}/${v.uuid}`" :disabled="p.current" />{{ v.name }}
                </label>
              </template>
            </div>
            <label class="flex items-center gap-2">After recovery:
              <select v-model="start" class="w-72 border border-line bg-white"><option v-for="m in modes" :key="m.value" :value="m.value">{{ m.text }}</option></select>
            </label>
          </div>
          <div v-else-if="page === 3" class="space-y-3">
            <p>Checking that the {{ verb }} can take place:</p>
            <div class="h-64 overflow-auto border border-neutral-400 p-1">
              <p v-if="checks === null" class="text-neutral-600">Running the prechecks...</p>
              <p v-else-if="!checks.length" class="flex items-center gap-2"><img src="/icons/000_Tick_h32bit_16.png" alt="" />All checks passed.</p>
              <p v-for="(c, i) in checks ?? []" :key="i" class="flex items-start gap-2">
                <img :src="c.error ? '/icons/000_Abort_h32bit_16.png' : '/icons/000_Alert2_h32bit_16.png'" :alt="c.error ? 'Error' : 'Warning'" class="mt-0.5" />
                <span class="whitespace-pre-line">{{ c.text }}</span>
              </p>
            </div>
          </div>
          <div v-else class="space-y-2">
            <p>Press Finish to start the {{ verb }}. These will be recovered on '{{ pool }}':</p>
            <template v-for="p in summary" :key="p.name">
              <p class="font-semibold">From {{ p.name }}</p>
              <p v-for="item in p.items" :key="item" class="ml-6">{{ item }}</p>
            </template>
            <p class="text-neutral-700">{{ start === "start" ? "They will be started after recovery." : start === "paused" ? "They will be started paused after recovery." : "They will not be started." }}</p>
            <p v-if="type === 'test'" class="text-neutral-700">This is a test: when it is done, you will be asked to remove what it recovered.</p>
            <p v-if="pending" class="text-neutral-700">Recovering...</p>
          </div>
        </div>
      </div>
      <p v-if="error || why && page !== 0" class="mx-4 mb-2 whitespace-pre-line" :class="error || page === 3 && checks?.some((c) => c.error) ? 'text-red-600' : 'text-neutral-600'">{{ error || why }}</p>
      <div class="flex justify-end gap-2 border-t border-line px-4 py-2">
        <button :class="btn" :disabled="page === 0 || pending" @click="page--">&lt; Previous</button>
        <button :class="btn" :disabled="page === 4 || !!why || pending" :title="why" @click="page++">Next &gt;</button>
        <button :class="btn" :disabled="page !== 4 || !chosenCount || pending" @click="finish">Finish</button>
        <button :class="btn" :disabled="pending" @click="emit('close')">Cancel</button>
      </div>
    </div>
  </div>
</template>
