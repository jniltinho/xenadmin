<script setup lang="ts">
// The New Storage wizard (NewSRWizard): Type, Name and Location pages for NFS and SMB storage and the NFS and SMB ISO
// libraries; the checks and the creation come from the Go side (SrRules, SrCreateAction).
// NFS and SMB storage can also reattach an SR a Scan finds on the share (SrReattachAction / SrIntroduceAction).
// iSCSI (LVMoISCSI): Scan Target Host lists the IQNs, then the LUNs; Finish probes the LUN and asks to reattach the SR
// found on it or to format it (LVMoIsciWarningDialog).
// shortcut: no HBA, FCoE or local storage.
import { computed, onMounted, reactive, ref, watch } from "vue";
import { api, type IQN, type LUN, type NewSRKind } from "./api/workspace";

const props = defineProps<{ server: string; pool: string }>();
const emit = defineEmits<{ (e: "close", created: string): void }>();

const pages = ["Type", "Name", "Location"] as const;
const page = ref(0);
const kinds = ref<NewSRKind[] | null>(null);
const busy = ref(false);
const error = ref("");
const f = reactive({ kind: "nfs", name: "", nameFor: "", autoDescription: true, description: "", path: "", options: "", nfsVersion: "",
  otherUser: false, user: "", password: "", reattach: false, existing: "",
  target: "", port: "3260", chap: false, chapUser: "", chapPassword: "", iqn: -1, lun: -1 });
const iqns = ref<IQN[]>([]);
const luns = ref<LUN[]>([]);
const found = ref<{ uuid: string; in_use?: string }[]>([]);
const scanStatus = ref("");
// a question of the wizard (the NFS reattach warning, iSCSI's Reattach/Format): its text and buttons
const question = ref<{ text: string; buttons: string[]; resolve: (i: number) => void } | null>(null);
const ask = (text: string, buttons: string[]) => new Promise<number>((resolve) => (question.value = { text, buttons, resolve }));
function answer(i: number) { const q = question.value; question.value = null; q?.resolve(i); }
let scanGeneration = 0;
const iso = computed(() => f.kind === "nfs_iso" || f.kind === "smb_iso");
const iscsi = computed(() => f.kind === "iscsi");

onMounted(async () => {
  try { kinds.value = await api.newSrKinds(props.server); } catch (e) { error.value = String(e); }
});
const kind = computed(() => kinds.value?.find((k) => k.kind === f.kind));
const nfs = computed(() => f.kind === "nfs" || f.kind === "nfs_iso");
const titles = computed(() => ["Choose the type of new storage", "What do you want to call this Storage Repository?",
  iscsi.value ? "Enter a path for your iSCSI storage" : nfs.value ? "Enter a path for your NFS storage" : "Enter a path for your SMB storage"]);
const groups = [{ title: "Virtual disk storage", sub: [
  { title: "Block based storage", items: [["iscsi", "iscsi"], ["Hardware HBA", ""], ["Software FCoE", ""]] },
  { title: "File based storage", items: [["nfs", "nfs"], ["smb", "smb"], ["Local Storage", ""]] }] },
{ title: "ISO library", sub: [{ title: "", items: [["smb_iso", "smb_iso"], ["nfs_iso", "nfs_iso"]] }] }];
const label = (text: string, k: string) => (k ? kinds.value?.find((x) => x.kind === k)?.name : text) ?? text;

// a new share, other options or credentials: the SRs found before belong to the old ones
function resetScan() { scanGeneration++; found.value = []; f.existing = ""; f.reattach = false; scanStatus.value = ""; }
watch(() => [f.kind, f.path, f.options, f.user, f.password], resetScan);
// another target or CHAP setting: what was scanned belongs to the old one (LVMoISCSI.ResetAll)
function resetTarget() { scanGeneration++; iqns.value = []; luns.value = []; f.iqn = -1; f.lun = -1; scanStatus.value = ""; }
watch(() => [f.target, f.port, f.chap, f.chapUser, f.chapPassword], resetTarget);
const chap = () => (f.chap ? [f.chapUser, f.chapPassword] : ["", ""]);

// CHAP needs a user name (LVMoISCSI.validChap): without it the scans and the SR would go unauthenticated
const chapMissing = () => f.chap && !f.chapUser.trim();
async function scanTarget() {
  if (chapMissing()) { error.value = "Enter the CHAP user name."; return; }
  const port = Number(f.port);
  if (!f.target.trim() || !Number.isInteger(port) || port < 1 || port > 65535) { error.value = "Enter the target host and a port number."; return; }
  resetTarget();
  const generation = scanGeneration;
  error.value = ""; busy.value = true;
  scanStatus.value = `Scanning for IQNs on iSCSI filer ${f.target.trim()}`;
  try {
    const found = await api.scanIqns(props.server, f.target.trim(), port, ...(chap() as [string, string]));
    if (generation !== scanGeneration) return;
    iqns.value = found;
    scanStatus.value = "";
    if (found.length === 1) { f.iqn = 0; await scanLuns(); }
  } catch (e) { if (generation === scanGeneration) scanStatus.value = String(e); } finally { busy.value = false; }
}
async function scanLuns() {
  luns.value = []; f.lun = -1;
  const iqn = iqns.value[f.iqn];
  if (!iqn) return;
  const generation = scanGeneration;
  busy.value = true;
  scanStatus.value = `Scanning for LUNs on iSCSI filer ${iqn.address}`;
  try {
    const found = await api.scanLuns(props.server, iqn, ...(chap() as [string, string]));
    if (generation !== scanGeneration || iqns.value[f.iqn] !== iqn) return;
    luns.value = found;
    if (found.length === 1) f.lun = 0;
    scanStatus.value = "";
  } catch (e) { if (generation === scanGeneration) scanStatus.value = String(e); } finally { busy.value = false; }
}

function spec() {
  const smbIso = f.kind === "smb_iso" && !f.otherUser; // an ISO library's credentials only with "Use different user name"
  return { kind: f.kind, name: f.name.trim(), description: f.description, auto_description: f.autoDescription, path: f.path,
    options: f.kind === "nfs" ? f.options : "", nfs_version: nfs.value ? f.nfsVersion : "", user: smbIso || nfs.value ? "" : f.user,
    password: smbIso || nfs.value ? "" : f.password, existing: f.reattach ? f.existing : "",
    iqn: iqns.value[f.iqn], lun: luns.value[f.lun], chap_user: chap()[0], chap_password: chap()[1] };
}

async function scan() {
  const generation = ++scanGeneration;
  error.value = ""; busy.value = true;
  scanStatus.value = `Performing SR scan from '${props.pool}'`;
  try {
    const srs = (await api.probeSr(props.server, spec())) ?? [];
    if (generation !== scanGeneration) return; // the location changed while scanning
    found.value = srs;
    scanStatus.value = srs.length ? "" : "<No SRs were found>";
  } catch (e) { if (generation === scanGeneration) scanStatus.value = String(e); } finally { busy.value = false; }
}

async function next() {
  error.value = "";
  if (page.value === 0 && f.nameFor !== f.kind) { f.name = kind.value?.default_name ?? ""; f.nameFor = f.kind; } // SrRules.DefaultName
  if (page.value === 1 && !f.name.trim()) { error.value = "Enter a name for the new SR."; return; }
  page.value++;
}

async function finish() {
  error.value = "";
  if (iscsi.value) return finishIscsi();
  if (f.reattach && !f.existing) { error.value = "Select the SR to reattach."; return; }
  const inUse = found.value.find((x) => x.uuid === f.existing)?.in_use;
  if (f.reattach && inUse) { error.value = `This SR is already in use by ${inUse}.`; return; }
  busy.value = true;
  try {
    const problem = await api.checkSrLocation(props.server, f.kind, f.path);
    if (problem) { error.value = problem; return; }
    const warning = f.reattach ? await api.srAttachWarning(f.kind, f.existing) : "";
    if (warning && (await ask(warning, ["Yes", "No"])) !== 0) return;
    emit("close", await api.createSr(props.server, spec()));
  } catch (e) { error.value = String(e); } finally { busy.value = false; }
}

// LVMoISCSI.PageLeave: the LUN must not hold an SR of a connected pool; an SR found on it is reattached, or the LUN
// formatted, as the user answers
async function finishIscsi() {
  const lun = luns.value[f.lun];
  if (!iqns.value[f.iqn] || !lun) { error.value = "Scan the target host and select the target IQN and LUN."; return; }
  if (chapMissing()) { error.value = "Enter the CHAP user name."; return; }
  busy.value = true;
  try {
    const used = await api.lunInUse(lun.scsi_id);
    if (used) { error.value = used; return; }
    const srs = (await api.probeSr(props.server, { ...spec(), existing: "" })) ?? [];
    const attached = srs.find((x) => x.in_use);
    if (attached) { error.value = `This SR is already in use by ${attached.in_use}.`; return; } // attached meanwhile somewhere
    const text = await api.lunQuestion(srs.length > 0);
    const i = await ask(text, srs.length ? ["Reattach", "Format", "Cancel"] : ["Format", "Cancel"]);
    if (srs.length ? i > 1 || i < 0 : i !== 0) return;
    emit("close", await api.createSr(props.server, { ...spec(), existing: srs.length && i === 0 ? srs[0].uuid : "" }));
  } catch (e) { error.value = String(e); } finally { busy.value = false; }
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <div class="flex h-[32rem] w-[50rem] flex-col border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">New Storage Repository - {{ pool }}</div>
      <div class="border-b border-line px-4 py-2 text-base font-semibold">{{ titles[page] }}</div>
      <div class="flex min-h-0 flex-1">
        <ol class="w-40 shrink-0 border-r border-line bg-pane py-2">
          <li v-for="(p, i) in pages" :key="p" class="px-3 py-1" :class="i === page ? 'font-semibold text-select' : i < page ? 'text-neutral-700' : 'text-neutral-400'">{{ p }}</li>
        </ol>
        <div v-if="kinds" class="min-w-0 flex-1 overflow-auto p-4">
          <div v-if="page === 0" class="grid grid-cols-[15rem_1fr] gap-4">
            <div>
              <template v-for="g in groups" :key="g.title">
                <p class="mt-2 mb-1 font-semibold">{{ g.title }}</p>
                <template v-for="s in g.sub" :key="s.title">
                  <p v-if="s.title" class="mt-1 mb-1 ml-3 font-semibold">{{ s.title }}</p>
                  <label v-for="[text, k] in s.items" :key="text" class="ml-6 flex items-center gap-2 py-0.5" :class="k ? '' : 'text-neutral-400'">
                    <input v-model="f.kind" type="radio" :value="k" :disabled="!k" />{{ label(text, k) }}</label>
                </template>
              </template>
            </div>
            <div class="mt-3"><p class="mb-3 italic">{{ kind?.name }}</p><p class="whitespace-pre-line">{{ kind?.blurb }}</p></div>
          </div>
          <template v-else-if="page === 1">
            <p class="mb-4">Provide a name and a description (optional) for your SR.</p>
            <label class="mb-4 grid grid-cols-[6rem_1fr] items-center gap-2">Name:<input v-model="f.name" class="border border-line px-1 py-0.5" /></label>
            <label class="mb-2 flex items-center gap-2"><input v-model="f.autoDescription" type="checkbox" />Autogenerate description based on SR settings (e.g., IP address, LUN etc.)</label>
            <label class="grid grid-cols-[6rem_1fr] items-start gap-2" :class="f.autoDescription ? 'text-neutral-400' : ''">Description:
              <textarea v-model="f.description" rows="4" :disabled="f.autoDescription" class="border border-line px-1 py-0.5" /></label>
          </template>
          <template v-else-if="iscsi">
            <p class="mb-4">Provide a target host for your iSCSI storage, indicating your target IQN and your target LUN before proceeding.</p>
            <div class="grid grid-cols-[8.5rem_1fr] items-center gap-x-2 gap-y-2">
              <span>Target Host:</span>
              <div class="flex items-center gap-2"><input v-model="f.target" class="min-w-0 flex-1 border border-line px-1 py-0.5" />
                <span>:</span><input v-model="f.port" class="w-16 border border-line px-1 py-0.5" /></div>
              <span /><label class="flex items-center gap-2"><input v-model="f.chap" type="checkbox" />Use CHAP authentication</label>
              <span :class="f.chap ? '' : 'text-neutral-400'">User:</span><input v-model="f.chapUser" :disabled="!f.chap" class="border border-line px-1 py-0.5" />
              <span :class="f.chap ? '' : 'text-neutral-400'">Password:</span><input v-model="f.chapPassword" type="password" :disabled="!f.chap" class="border border-line px-1 py-0.5" />
              <span /><div><button class="border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="busy || !f.target.trim()" @click="scanTarget">Scan Target Host</button></div>
              <span :class="iqns.length ? '' : 'text-neutral-400'">Target IQN:</span>
              <select v-model.number="f.iqn" :disabled="!iqns.length || busy" class="border border-line bg-white px-1 py-0.5" @change="scanLuns">
                <option :value="-1">Select target IQN</option>
                <option v-for="(q, i) in iqns" :key="q.text" :value="i">{{ q.text }}</option>
              </select>
              <span :class="luns.length ? '' : 'text-neutral-400'">Target LUN:</span>
              <select v-model.number="f.lun" :disabled="!luns.length || busy" class="border border-line bg-white px-1 py-0.5">
                <option :value="-1">Select target LUN</option>
                <option v-for="(l, i) in luns" :key="l.scsi_id + l.serial" :value="i">{{ l.text }}</option>
              </select>
            </div>
            <p v-if="scanStatus" class="mt-3 text-neutral-600">{{ scanStatus }}</p>
          </template>
          <template v-else>
            <p class="mb-4">{{ f.kind === "nfs_iso" ? "Provide the name of the share where your NFS ISO storage is located."
              : f.kind === "smb_iso" ? "Provide the name of the share where your SR is located. You can optionally specify alternative credentials by setting the server options."
              : nfs ? "Provide the name of the share where your SR is located, optionally specifying advanced options."
              : "Provide the name of the share where your SR is located, optionally specifying the credentials for it." }}</p>
            <div class="grid grid-cols-[8.5rem_1fr] items-center gap-x-2 gap-y-2">
              <span>Share Name:</span>
              <div class="flex gap-2"><input v-model="f.path" class="min-w-0 flex-1 border border-line px-1 py-0.5" />
                <button v-if="!iso" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="busy || !f.path.trim()" @click="scan">Scan</button></div>
              <span /><span class="-mt-1 text-neutral-600">{{ nfs ? "Example: server:/path" : "Example: \\\\server\\share" }}</span>
              <template v-if="nfs">
                <span class="self-start">NFS Version:</span>
                <div class="flex flex-col gap-1">
                  <label v-for="[t, v] in f.kind === 'nfs' ? [['NFSv3', ''], ['NFSv4', '4'], ['NFSv4.1', '4.1']] : [['NFSv3', ''], ['NFSv4', '4']]" :key="v"
                         class="flex items-center gap-2"><input v-model="f.nfsVersion" type="radio" :value="v" />{{ t }}</label>
                </div>
                <template v-if="f.kind === 'nfs'"><span>Advanced Options:</span><input v-model="f.options" class="border border-line px-1 py-0.5" /></template>
              </template>
              <template v-else>
                <template v-if="f.kind === 'smb_iso'"><span /><label class="flex items-center gap-2"><input v-model="f.otherUser" type="checkbox" />Use different user name</label></template>
                <span>Username:</span><input v-model="f.user" :disabled="f.kind === 'smb_iso' && !f.otherUser" class="border border-line px-1 py-0.5" />
                <span>Password:</span><input v-model="f.password" type="password" :disabled="f.kind === 'smb_iso' && !f.otherUser" class="border border-line px-1 py-0.5" />
              </template>
            </div>
            <p v-if="scanStatus" class="mt-3 text-neutral-600">{{ scanStatus }}</p>
            <template v-if="!iso">
              <label class="mt-3 flex items-center gap-2"><input v-model="f.reattach" type="radio" :value="false" />Create a new SR</label>
              <label class="mt-1 flex items-center gap-2" :class="found.length ? '' : 'text-neutral-400'">
                <input v-model="f.reattach" type="radio" :value="true" :disabled="!found.length" />Reattach an existing SR:</label>
              <div class="ml-6 mt-1 h-20 overflow-auto border border-neutral-400" :class="found.length ? '' : 'bg-neutral-50'">
                <div v-for="x in found" :key="x.uuid" class="cursor-default px-2" :class="f.existing === x.uuid ? 'bg-select text-white' : 'hover:bg-neutral-100'"
                     @click="f.existing = x.uuid; f.reattach = true">{{ x.uuid }}<span v-if="x.in_use"> (in use by {{ x.in_use }})</span></div>
              </div>
            </template>
          </template>
        </div>
        <div v-else class="p-4 text-neutral-500">{{ error || "Loading..." }}</div>
      </div>
      <div class="flex items-center justify-end gap-2 border-t border-line p-2">
        <span v-if="error && kinds" class="mr-auto text-red-600">{{ error }}</span>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="page === 0 || busy" @click="page--">&lt; Previous</button>
        <button v-if="page < pages.length - 1" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="!kinds" @click="next">Next &gt;</button>
        <button v-else class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="busy || (iscsi ? f.lun < 0 : !f.path.trim())" @click="finish()">{{ busy ? "Working..." : "Finish" }}</button>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="busy" @click="emit('close', '')">Cancel</button>
      </div>
    </div>
    <div v-if="question" class="fixed inset-0 z-40 flex items-center justify-center bg-black/10">
      <div class="w-[34rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">New Storage Repository - {{ pool }}</div>
        <p class="whitespace-pre-line p-4">{{ question.text }}</p>
        <div class="flex justify-end gap-2 border-t border-line p-2">
          <button v-for="(b, i) in question.buttons" :key="b" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="answer(i)">{{ b }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
