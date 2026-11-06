<script setup lang="ts">
// Properties of a server, pool, SR or network (PropertiesDialog.PagesFor): General (name, description, folder, tags; a
// server's iSCSI IQN) and Custom Fields, with the Change Folder and Edit Tags dialogs (Dialogs.ChangeFolder/EditTags).
// A server also has Multipathing (HostMultipathPage) and Log Destination (LogDestinationEditPage).
// shortcut: the Alerts and Network Settings pages are not ported (as in the C# app).
import { computed, onMounted, reactive, ref } from "vue";
import { api, type ObjectProps } from "./api/workspace";
import FolderTags from "./FolderTags.vue";

const props = defineProps<{ server: string; kind: string; uuid: string }>();
const emit = defineEmits<{ (e: "close", saved: boolean): void }>();

const orig = ref<ObjectProps | null>(null);
type Page = "general" | "fields" | "multipath" | "log";
const page = ref<Page>("general");
const f = reactive({ name: "", description: "", folder: "", tags: [] as string[], iqn: "", multipath: false, remote: false, server: "" });
// the pages, with a server's own two (PropertiesDialog.PagesFor)
const pages = computed(() => [{ id: "general", label: "General", icon: "000_Info3_h32bit_16" }, { id: "fields", label: "Custom Fields", icon: "000_Fields_h32bit_16" },
  ...(orig.value?.multipath != null ? [{ id: "multipath", label: "Multipathing", icon: "000_Storage_h32bit_16" }] : []),
  ...(orig.value?.log_destination != null ? [{ id: "log", label: "Log Destination", icon: "log_destination_16" }] : [])]);
// LogDestinationEditPage.regex: a host name or address
const syslogOk = (v: string) => /^[a-zA-Z0-9]([-a-zA-Z0-9]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([-a-zA-Z0-9]{0,61}[a-zA-Z0-9])?)*$/.test(v);
const fieldEdits = reactive<Record<string, string>>(Object.create(null));
const error = ref("");
const pending = ref(false);
const toLocal = (rfc: string) => {
  if (!rfc) return "";
  const d = new Date(rfc), p = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}T${p(d.getHours())}:${p(d.getMinutes())}`;
};

onMounted(async () => {
  try {
    const o = await api.objectProperties(props.server, props.kind, props.uuid);
    orig.value = o;
    Object.assign(f, { name: o.name, description: o.description, folder: o.folder, tags: [...o.tags], iqn: o.iqn ?? "",
      multipath: !!o.multipath, remote: !!o.log_destination, server: o.log_destination ?? "" });
    for (const x of o.fields) fieldEdits[x.name] = x.type === "Date" ? toLocal(x.value) : x.value;
  } catch (e) { error.value = String(e); }
});
const title = computed(() => `'${orig.value?.name ?? ""}' Properties`);

async function save() {
  const o = orig.value!;
  if (pending.value) return;
  const ch: Record<string, unknown> = {};
  if (f.name.trim() !== o.name) ch.name = f.name.trim();
  if (f.description !== o.description) ch.description = f.description;
  if (f.folder !== o.folder) ch.folder = f.folder;
  if (f.tags.join("\n") !== [...o.tags].sort().join("\n")) { ch.tags = f.tags; ch.tags_set = true; }
  if (o.iqn != null && f.iqn.trim() !== o.iqn) ch.iqn = f.iqn.trim();
  if (o.multipath != null && f.multipath !== o.multipath) ch.multipath = f.multipath;
  if (o.log_destination != null) {
    const to = f.remote ? f.server.trim() : "";
    if (f.remote && !syslogOk(to)) { error.value = "Please enter a valid IP address or server name."; page.value = "log"; return; }
    if (to !== o.log_destination) ch.log_destination = to;
  }
  const fields: Record<string, string> = Object.create(null);
  for (const x of o.fields) {
    const v = fieldEdits[x.name] ?? "", was = x.type === "Date" ? toLocal(x.value) : x.value;
    if (v === was) continue;
    if (x.type === "Date" && v && Number.isNaN(new Date(v).getTime())) { error.value = "Enter a valid date."; page.value = "fields"; return; }
    fields[x.name] = x.type === "Date" && v ? new Date(v).toISOString().replace(/\.\d{3}Z$/, "Z") : v;
  }
  if (Object.keys(fields).length) ch.fields = fields;
  if (!f.name.trim()) { error.value = "Enter a name."; page.value = "general"; return; }
  if (!Object.keys(ch).length) { emit("close", false); return; }
  pending.value = true;
  try { await api.saveObjectProperties(props.server, props.kind, props.uuid, ch); emit("close", true); }
  catch (e) { error.value = String(e); } finally { pending.value = false; }
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <div class="flex h-[30rem] w-[46rem] flex-col border border-neutral-500 bg-white shadow-lg">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ title }}</div>
      <div class="flex min-h-0 flex-1">
        <div class="w-48 border-r border-line bg-pane py-1">
          <div v-for="p in pages" :key="p.id"
               class="flex cursor-default items-center gap-2 px-2 py-1.5" :class="page === p.id ? 'bg-select text-white' : 'hover:bg-neutral-100'" @click="page = p.id as Page">
            <img :src="`/icons/${p.icon}.png`" alt="" class="h-4 w-4" />{{ p.label }}
          </div>
        </div>
        <fieldset v-if="orig" :disabled="pending" class="min-w-0 flex-1 overflow-auto p-4">
          <template v-if="page === 'general'">
            <p class="mb-3">XCP-ng Center provides several methods for labelling and organizing your resources. By using features such as folders and tags you can quickly group and manage resources across your connected servers.</p>
            <div class="grid grid-cols-[7rem_1fr] items-start gap-2">
              <label for="opname" class="pt-0.5">Name:</label><input id="opname" v-model="f.name" class="border border-line px-1 py-0.5" />
              <label for="opdesc" class="pt-0.5">Description:</label><textarea id="opdesc" v-model="f.description" rows="3" class="border border-line px-1 py-0.5" />
              <FolderTags v-model:folder="f.folder" v-model:tags="f.tags" :folders="orig.folders" :all-tags="orig.all_tags" />
              <template v-if="orig.iqn != null"><label for="opiqn" class="pt-0.5">iSCSI IQN:</label><input id="opiqn" v-model="f.iqn" class="border border-line px-1 py-0.5" /></template>
            </div>
          </template>
          <template v-else-if="page === 'multipath'">
            <p class="mb-3 whitespace-pre-line">Dynamic multipathing support is available for some types of storage repository.

The server must be in maintenance mode before you can change its multipathing setting. This ensures that any running virtual machines with virtual disks in the affected storage repository are migrated before the changes are made.</p>
            <label class="flex items-center gap-2"><input v-model="f.multipath" type="checkbox" :disabled="!orig.multipath_editable" />Enable multipathing on this server</label>
            <p v-if="!orig.multipath_editable" class="mt-3 flex items-center gap-2"><img src="/icons/000_WarningAlert_h32bit_32.png" alt="" class="h-4 w-4" />The server must be in maintenance mode to edit these settings.</p>
          </template>
          <template v-else-if="page === 'log'">
            <p class="mb-3 font-semibold">Server system log destination</p>
            <p class="mb-3">XCP-ng stores its system logs locally on the server. If you want the system logs to be directed to a remote server in addition to being stored locally, you may specify a remote log destination using the settings below.</p>
            <label class="flex items-center gap-2"><input v-model="f.remote" type="checkbox" />Also store the system logs on a remote server</label>
            <div class="mt-2 flex items-center gap-2 pl-6">
              <label for="oplog">Server:</label><input id="oplog" v-model="f.server" class="w-64 border border-line px-1 py-0.5" @focus="f.remote = true" />
            </div>
          </template>
          <template v-else>
            <p class="mb-3">Custom fields hold your own information about this object, such as its owner or its expiry date.</p>
            <div class="grid grid-cols-[10rem_1fr] items-center gap-2">
              <template v-for="x in orig.fields" :key="x.name">
                <label :for="'ocf-' + x.name" class="truncate">{{ x.name }}:</label>
                <input :id="'ocf-' + x.name" v-model="fieldEdits[x.name]" :type="x.type === 'Date' ? 'datetime-local' : 'text'" class="border border-line px-1 py-0.5" />
              </template>
            </div>
            <p v-if="!orig.fields.length" class="text-neutral-500">There are no custom fields. Use Edit Custom Fields in a VM's Properties to add one.</p>
          </template>
        </fieldset>
        <div v-else class="p-4 text-neutral-500">{{ error || "Loading..." }}</div>
      </div>
      <div class="flex items-center justify-end gap-2 border-t border-line p-2">
        <span v-if="error && orig" class="mr-auto text-red-600">{{ error }}</span>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="!orig || pending" @click="save">OK</button>
        <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close', false)">Cancel</button>
      </div>
    </div>

  </div>
</template>
