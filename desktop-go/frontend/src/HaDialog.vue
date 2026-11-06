<script setup lang="ts">
// Pool > High Availability (MainWindow.ConfigureHa and DisableHa, XenAdmin.Desktop.Core/HaRules.cs):
// - Configure with HA off: the HA wizard (HaWizard.cs): Prerequisites, Heartbeat SR (GetHeartbeatSRsAction), HA Plan
//   (HaPlanEditor.cs: each VM's restart priority, start order and delay; the failures to tolerate, whose maximum the
//   pool computes for the priorities) and Finish (HAFinishPage), then EnableHAAction;
// - Configure with HA on: the plan alone (EditVmHaPrioritiesDialog), then SetHaPrioritiesAction;
// - Disable: HADisableCommand's question, then DisableHAAction.
// The Go side checks everything again (ops.EnableHA, SetHAPriorities, DisableHA).
import { computed, onMounted, reactive, ref, watch } from "vue";
import { api, type HAConfig, type HeartbeatSR, type HostQuestion } from "./api/workspace";
import { serial } from "./serial";

const props = defineProps<{ server: string; disable?: boolean }>();
const emit = defineEmits<{ (e: "close"): void; (e: "run", label: string, f: () => Promise<void>): void }>();

const cfg = ref<HAConfig | null>(null);
const question = ref<HostQuestion | null>(null);
const error = ref("");
type Row = { uuid: string; name: string; priority: string; order: number; delay: number; notAgile: string };
const rows = reactive<Row[]>([]);
const ntol = ref(0);
let touched = false; // HaNtolControl: a value the user chose (or the pool's, HA on) is never changed for them
const max = ref(-1);
const info = ref("");
const pending = ref(false);

onMounted(async () => {
  try {
    if (props.disable) { question.value = await api.haDisableQuestion(props.server); return; }
    const c = await api.haConfig(props.server);
    cfg.value = c;
    if (c.blocker) return;
    rows.push(...(c.vms ?? []).map((v) => ({ uuid: v.uuid, name: v.name, priority: v.priority, order: v.order, delay: v.start_delay, notAgile: v.not_agile ?? "" })));
    ntol.value = c.ha_enabled ? c.ntol : 0;
    touched = c.ha_enabled;
    recompute();
  } catch (e) { error.value = String(e); }
});

const priorities = computed(() => cfg.value?.priorities ?? []);
const priorityText = (p: string) => priorities.value.find((x) => x.value === p)?.text ?? p;
// HaRules.Allowed: a VM that is not agile is only restarted if possible, or not at all
const choices = (r: Row) => priorities.value.filter((p) => !r.notAgile || p.value !== "restart");

// HaNtolControl.updateNtol: the maximum for the current priorities; the last change wins
const recompute = serial(async () => {
  max.value = -1;
  info.value = "Computing the number of server failures the pool can tolerate...";
  const m = await api.haMaxFailures(props.server, Object.fromEntries(rows.map((r) => [r.uuid, r.priority])));
  max.value = m.max;
  info.value = m.info;
  if (m.max >= 0 && !touched) ntol.value = m.max; // a new configuration starts at the maximum
});
watch(() => rows.map((r) => r.priority).join(), () => recompute().catch((e) => { error.value = String(e); }));
function ntolChanged() { touched = true; }

// AssignPriorities.updateButtons (HaRules.PlanError)
const planError = computed(() => {
  if (max.value < 0) return "The number of server failures that can be tolerated is still being computed.";
  if (!Number.isInteger(ntol.value) || ntol.value < 0) return "Enter the number of server failures to tolerate.";
  if (ntol.value > max.value) return cfg.value!.overcommit;
  if (!rows.every((r) => Number.isInteger(r.order) && r.order >= 0 && Number.isInteger(r.delay) && r.delay >= 0)) return "Enter a whole number of 0 or more.";
  return "";
});
const plan = () => rows.map((r) => ({ vm: r.uuid, priority: r.priority, order: r.order, start_delay: r.delay }));

// the wizard (HA off)
const pages = ["Prerequisites", "Heartbeat SR", "HA Plan", "Finish"];
const titles = ["HA configuration prerequisites", "Choose a heartbeat SR",
  "Configure the HA restart priority, restart order and delay interval for the VMs in this pool", "Review configuration and activate HA"];
const page = ref(0);
const srs = ref<HeartbeatSR[] | null>(null);
const heartbeat = ref("");
const scanning = ref(false);
// ChooseSR: a new scan decides again which SRs are usable; nothing stays selected from the previous one
async function scan() {
  srs.value = null; heartbeat.value = ""; scanning.value = true; error.value = "";
  try {
    const list = await api.heartbeatSrs(props.server);
    srs.value = list;
    heartbeat.value = list.find((s) => !s.reason)?.uuid ?? "";
    if (!list.length) error.value = "There are no shared storage repositories in this pool.";
  } catch (e) { error.value = String(e); } finally { scanning.value = false; }
}
watch(page, (p) => { if (p === 1 && srs.value === null && !scanning.value) scan(); });
const why = computed(() => (page.value === 1 && !heartbeat.value ? "Select a heartbeat SR." : page.value === 2 ? planError.value : ""));
const heartbeatName = computed(() => srs.value?.find((s) => s.uuid === heartbeat.value)?.name ?? "");
const count = (restart: boolean | null) => rows.filter((r) => (restart === null ? r.priority === "best-effort" : restart ? r.priority === "restart" : r.priority === "")).length;

function finish() {
  if (pending.value || planError.value || !heartbeat.value) return;
  pending.value = true;
  const sr = heartbeat.value, p = plan(), n = ntol.value;
  emit("run", `Enabling HA on pool '${cfg.value!.pool}'`, () => api.enableHa(props.server, sr, p, n));
  emit("close");
}

// the restart priorities (HA on): no failure tolerated asks first, No being the default
const askZero = ref(false);
function save(confirmed = false) {
  if (pending.value || planError.value) return;
  if (ntol.value === 0 && !confirmed) { askZero.value = true; return; }
  pending.value = true;
  const p = plan(), n = ntol.value;
  emit("run", "Setting HA protection levels", () => api.setHaPriorities(props.server, p, n));
  emit("close");
}

function disable() {
  if (pending.value) return;
  pending.value = true;
  emit("run", "Disabling HA", () => api.disableHa(props.server));
  emit("close");
}
const btn = "min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400";
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <!-- a message: Disable's question or why it cannot run, Configure's blocker, a failed read -->
    <div v-if="question || error && !cfg || cfg?.blocker || askZero" class="w-[30rem] border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ askZero ? "High Availability" : question ? question.title || "Disable HA" : "Configure HA" }}</div>
      <div class="space-y-3 p-3">
        <p class="whitespace-pre-line" :class="{ 'text-red-600': error && !cfg }">{{ askZero ? cfg!.ntol_zero_query : question ? question.error || question.text : cfg?.blocker || error }}</p>
        <div class="flex justify-end gap-2">
          <template v-if="askZero">
            <button :class="btn" :disabled="pending" @click="save(true)">Yes</button>
            <button :class="btn" autofocus @click="askZero = false">No</button>
          </template>
          <template v-else-if="question && !question.error">
            <button :class="btn" :disabled="pending" @click="disable">Yes</button>
            <button :class="btn" autofocus @click="emit('close')">No</button>
          </template>
          <button v-else :class="btn" @click="emit('close')">OK</button>
        </div>
      </div>
    </div>
    <div v-else-if="!cfg" class="w-[30rem] border border-neutral-500 bg-white p-3 shadow-lg">
      <p>Loading...</p>
      <div class="flex justify-end"><button :class="btn" @click="emit('close')">Cancel</button></div>
    </div>
    <div v-else class="flex w-[52rem] flex-col border border-neutral-500 bg-white shadow-lg" :class="cfg.ha_enabled ? '' : 'h-[34rem]'">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ cfg.ha_enabled ? "HA restart priorities" : `Configure HA on '${cfg.pool}'` }}</div>
      <div v-if="!cfg.ha_enabled" class="border-b border-line px-4 py-2 text-base font-semibold">{{ titles[page] }}</div>
      <div class="flex min-h-0 flex-1">
        <ol v-if="!cfg.ha_enabled" class="w-40 shrink-0 border-r border-line bg-pane py-2">
          <li v-for="(p, i) in pages" :key="p" class="px-3 py-1" :class="i === page ? 'font-semibold text-select' : i < page ? 'text-neutral-700' : 'text-neutral-400'">{{ p }}</li>
        </ol>
        <div class="min-w-0 flex-1 overflow-auto p-4">
          <div v-if="!cfg.ha_enabled && page === 0" class="space-y-3">
            <p>This wizard helps you configure HA (high availability) on a resource pool.</p>
            <p>Before you begin, ensure that the following requirements are satisfied for all servers and virtual machines in the pool:</p>
            <ul class="ml-6 list-disc space-y-2">
              <li>Shared storage must be available, including at least one iSCSI, NFS or Fibre Channel LUN. This LUN will be used for the heartbeat SR.</li>
              <li>All the virtual machines you want to protect with HA must be agile.</li>
            </ul>
            <p>We strongly recommend the use of a bonded management interface on the servers in the pool if HA is enabled, and multipathed storage for the heartbeat SR.</p>
            <p>To continue, click Next.</p>
          </div>
          <div v-else-if="!cfg.ha_enabled && page === 1" class="space-y-3">
            <p>Select a heartbeat SR that will be used to monitor the availability and health of servers in the pool. Heartbeat SRs must use shared storage.</p>
            <div class="h-60 overflow-auto border border-neutral-400">
              <p v-if="srs === null" class="p-1 text-neutral-600">Scanning for HA heartbeat SRs on {{ cfg.pool }}</p>
              <label v-for="s in srs ?? []" :key="s.uuid" class="flex items-center gap-2 px-1 py-0.5" :class="s.reason ? 'text-neutral-400' : ''">
                <input v-model="heartbeat" type="radio" :value="s.uuid" :disabled="!!s.reason" />{{ s.name }}<span v-if="s.reason">({{ s.reason }})</span>
              </label>
            </div>
            <button :class="btn" :disabled="scanning" @click="scan">Rescan</button>
          </div>
          <div v-else-if="cfg.ha_enabled || page === 2" class="space-y-3">
            <div class="max-h-60 overflow-auto border border-neutral-400">
              <table v-if="rows.length" class="w-full">
                <thead><tr class="bg-pane text-left"><th class="px-2 py-1">VM</th><th class="px-2">HA restart priority</th><th class="px-2">Start order</th><th class="px-2">Delay (s)</th><th class="px-2">Agile</th></tr></thead>
                <tbody>
                  <tr v-for="r in rows" :key="r.uuid" class="border-t border-line">
                    <td class="px-2 py-0.5">{{ r.name }}</td>
                    <td class="px-2"><select v-model="r.priority" :aria-label="`HA restart priority of ${r.name}`" class="w-44 border border-line bg-white">
                      <option v-for="p in choices(r)" :key="p.value" :value="p.value">{{ p.text }}</option></select></td>
                    <td class="px-2"><input v-model.number="r.order" type="number" min="0" step="1" :aria-label="`Start order of ${r.name}`" class="w-24 border border-line px-1" /></td>
                    <td class="px-2"><input v-model.number="r.delay" type="number" min="0" step="1" :aria-label="`Delay of ${r.name}`" class="w-24 border border-line px-1" /></td>
                    <td class="px-2" :title="r.notAgile">{{ r.notAgile ? "No" : "Yes" }}</td>
                  </tr>
                </tbody>
              </table>
              <p v-else class="p-1">There are no VMs in this pool that HA can protect.</p>
            </div>
            <p>You can specify the number of server failures you need to be able to tolerate in the pool.</p>
            <label class="flex items-center gap-2">Failures tolerated:
              <input v-model.number="ntol" type="number" min="0" :max="cfg.hosts" step="1" class="w-28 border border-line px-1 py-0.5" @input="ntolChanged" /></label>
            <p class="whitespace-pre-line text-neutral-700">{{ info }}</p>
          </div>
          <div v-else class="space-y-3">
            <p>Review the final HA configuration below. You can press Finish to enable HA with this configuration, or Previous to alter the configuration.</p>
            <div class="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1.5 border border-line p-3">
              <span>Heartbeat SR:</span><span>{{ heartbeatName }}</span>
              <span>Failure Tolerance:</span><span>{{ ntol }}</span>
              <span>{{ priorityText("restart") }}:</span><span>{{ count(true) }} VMs</span>
              <span>{{ priorityText("best-effort") }}:</span><span>{{ count(null) }} VMs</span>
              <span>{{ priorityText("") }}:</span><span>{{ count(false) }} VMs</span>
            </div>
            <p v-if="ntol === 0" class="text-yellow-700">Under this configuration, the failover tolerance of your pool would be 0 (zero) and it may not be possible to restart your virtual machines in the event of a server failure.</p>
          </div>
        </div>
      </div>
      <p v-if="error || (cfg.ha_enabled || page === 2) && planError" class="mx-4 mb-2 whitespace-pre-line text-red-600">{{ error || planError }}</p>
      <div class="flex justify-end gap-2 border-t border-line px-4 py-2">
        <template v-if="cfg.ha_enabled">
          <button :class="btn" :disabled="!!planError || pending" @click="save()">OK</button>
        </template>
        <template v-else>
          <button :class="btn" :disabled="page === 0 || pending" @click="page--">&lt; Previous</button>
          <button :class="btn" :disabled="page === 3 || !!why" :title="why" @click="page++">Next &gt;</button>
          <button :class="btn" :disabled="page !== 3 || !!planError || !heartbeat || pending" @click="finish">Finish</button>
        </template>
        <button :class="btn" @click="emit('close')">Cancel</button>
      </div>
    </div>
  </div>
</template>
