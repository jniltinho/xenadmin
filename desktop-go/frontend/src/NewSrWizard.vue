<script setup lang="ts">
// The New Storage wizard (NewSRWizard): Type, Name and Location pages for NFS and SMB storage and the NFS and SMB ISO
// libraries; the checks and the creation come from the Go side (SrRules, SrCreateAction).
// shortcut: no Scan/Reattach of an existing SR and no iSCSI, HBA or local storage yet.
import { computed, onMounted, reactive, ref } from "vue";
import { api, type NewSRKind } from "./api/workspace";

const props = defineProps<{ server: string; pool: string }>();
const emit = defineEmits<{ (e: "close", created: string): void }>();

const pages = ["Type", "Name", "Location"] as const;
const page = ref(0);
const kinds = ref<NewSRKind[] | null>(null);
const busy = ref(false);
const error = ref("");
const f = reactive({ kind: "nfs", name: "", nameFor: "", autoDescription: true, description: "", path: "", options: "", nfsVersion: "",
  otherUser: false, user: "", password: "" });

onMounted(async () => {
  try { kinds.value = await api.newSrKinds(props.server); } catch (e) { error.value = String(e); }
});
const kind = computed(() => kinds.value?.find((k) => k.kind === f.kind));
const nfs = computed(() => f.kind === "nfs" || f.kind === "nfs_iso");
const titles = computed(() => ["Choose the type of new storage", "What do you want to call this Storage Repository?",
  nfs.value ? "Enter a path for your NFS storage" : "Enter a path for your SMB storage"]);
const groups = [{ title: "Virtual disk storage", sub: [
  { title: "Block based storage", items: [["iSCSI", ""], ["Hardware HBA", ""], ["Software FCoE", ""]] },
  { title: "File based storage", items: [["nfs", "nfs"], ["smb", "smb"], ["Local Storage", ""]] }] },
{ title: "ISO library", sub: [{ title: "", items: [["smb_iso", "smb_iso"], ["nfs_iso", "nfs_iso"]] }] }];
const label = (text: string, k: string) => (k ? kinds.value?.find((x) => x.kind === k)?.name : text) ?? text;

async function next() {
  error.value = "";
  if (page.value === 0 && f.nameFor !== f.kind) { f.name = kind.value?.default_name ?? ""; f.nameFor = f.kind; } // SrRules.DefaultName
  if (page.value === 1 && !f.name.trim()) { error.value = "Enter a name for the new SR."; return; }
  page.value++;
}

async function finish() {
  busy.value = true; error.value = "";
  try {
    const problem = await api.checkSrLocation(props.server, f.kind, f.path);
    if (problem) { error.value = problem; busy.value = false; return; }
    const iso = f.kind === "smb_iso" && !f.otherUser; // an ISO library's credentials only with "Use different user name"
    const uuid = await api.createSr(props.server, { kind: f.kind, name: f.name.trim(), description: f.description, auto_description: f.autoDescription,
      path: f.path, options: f.kind === "nfs" ? f.options : "", nfs_version: nfs.value ? f.nfsVersion : "", user: iso || nfs.value ? "" : f.user,
      password: iso || nfs.value ? "" : f.password });
    emit("close", uuid);
  } catch (e) { error.value = String(e); }
  busy.value = false;
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
          <template v-else>
            <p class="mb-4">{{ f.kind === "nfs_iso" ? "Provide the name of the share where your NFS ISO storage is located."
              : f.kind === "smb_iso" ? "Provide the name of the share where your SR is located. You can optionally specify alternative credentials by setting the server options."
              : nfs ? "Provide the name of the share where your SR is located, optionally specifying advanced options."
              : "Provide the name of the share where your SR is located, optionally specifying the credentials for it." }}</p>
            <div class="grid grid-cols-[8.5rem_1fr] items-center gap-x-2 gap-y-2">
              <span>Share Name:</span><input v-model="f.path" class="border border-line px-1 py-0.5" />
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
          </template>
        </div>
        <div v-else class="p-4 text-neutral-500">{{ error || "Loading..." }}</div>
      </div>
      <div class="flex items-center justify-end gap-2 border-t border-line p-2">
        <span v-if="error && kinds" class="mr-auto text-red-600">{{ error }}</span>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="page === 0 || busy" @click="page--">&lt; Previous</button>
        <button v-if="page < pages.length - 1" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="!kinds" @click="next">Next &gt;</button>
        <button v-else class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="busy || !f.path.trim()" @click="finish">{{ busy ? "Creating..." : "Finish" }}</button>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="busy" @click="emit('close', '')">Cancel</button>
      </div>
    </div>
  </div>
</template>
