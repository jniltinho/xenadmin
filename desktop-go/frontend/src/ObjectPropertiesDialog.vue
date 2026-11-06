<script setup lang="ts">
// Properties of a server, pool, SR or network (PropertiesDialog.PagesFor): General (name, description, folder, tags; a
// server's iSCSI IQN) and Custom Fields, with the Change Folder and Edit Tags dialogs (Dialogs.ChangeFolder/EditTags).
// shortcut: the Alerts, Multipathing, Log Destination and Network Settings pages are not ported (as in the C# app).
import { computed, onMounted, reactive, ref } from "vue";
import { api, type ObjectProps } from "./api/workspace";

const props = defineProps<{ server: string; kind: string; uuid: string }>();
const emit = defineEmits<{ (e: "close", saved: boolean): void }>();

const orig = ref<ObjectProps | null>(null);
const page = ref<"general" | "fields">("general");
const f = reactive({ name: "", description: "", folder: "", tags: [] as string[], iqn: "" });
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
    Object.assign(f, { name: o.name, description: o.description, folder: o.folder, tags: [...o.tags], iqn: o.iqn ?? "" });
    for (const x of o.fields) fieldEdits[x.name] = x.type === "Date" ? toLocal(x.value) : x.value;
  } catch (e) { error.value = String(e); }
});
const title = computed(() => `'${orig.value?.name ?? ""}' Properties`);

// Change Folder: none, an existing folder, or a new path
const folderEdit = ref<{ none: boolean; path: string } | null>(null);
function chooseFolder() {
  const fe = folderEdit.value!;
  f.folder = fe.none ? "" : fe.path.trim();
  folderEdit.value = null;
}
// Edit Tags: the connection's tags ticked for this object, and new ones
const tagEdit = ref<{ all: string[]; on: Set<string>; added: string } | null>(null);
function openTags() { tagEdit.value = { all: [...new Set([...(orig.value?.all_tags ?? []), ...f.tags])].sort(), on: new Set(f.tags), added: "" }; }
function addTag() {
  const t = tagEdit.value!, name = t.added.trim();
  if (name && !t.all.includes(name)) t.all = [...t.all, name].sort();
  if (name) t.on.add(name);
  t.added = "";
}
function chooseTags() { f.tags = [...tagEdit.value!.on].sort(); tagEdit.value = null; }

async function save() {
  const o = orig.value!;
  if (pending.value) return;
  const ch: Record<string, unknown> = {};
  if (f.name.trim() !== o.name) ch.name = f.name.trim();
  if (f.description !== o.description) ch.description = f.description;
  if (f.folder !== o.folder) ch.folder = f.folder;
  if (f.tags.join("\n") !== [...o.tags].sort().join("\n")) { ch.tags = f.tags; ch.tags_set = true; }
  if (o.iqn != null && f.iqn.trim() !== o.iqn) ch.iqn = f.iqn.trim();
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
          <div v-for="p in [{ id: 'general', label: 'General', icon: '000_Info3_h32bit_16' }, { id: 'fields', label: 'Custom Fields', icon: '000_Fields_h32bit_16' }]" :key="p.id"
               class="flex cursor-default items-center gap-2 px-2 py-1.5" :class="page === p.id ? 'bg-select text-white' : 'hover:bg-neutral-100'" @click="page = p.id as 'general' | 'fields'">
            <img :src="`/icons/${p.icon}.png`" alt="" class="h-4 w-4" />{{ p.label }}
          </div>
        </div>
        <fieldset v-if="orig" :disabled="pending" class="min-w-0 flex-1 overflow-auto p-4">
          <template v-if="page === 'general'">
            <p class="mb-3">XCP-ng Center provides several methods for labelling and organizing your resources. By using features such as folders and tags you can quickly group and manage resources across your connected servers.</p>
            <div class="grid grid-cols-[7rem_1fr] items-start gap-2">
              <label for="opname" class="pt-0.5">Name:</label><input id="opname" v-model="f.name" class="border border-line px-1 py-0.5" />
              <label for="opdesc" class="pt-0.5">Description:</label><textarea id="opdesc" v-model="f.description" rows="3" class="border border-line px-1 py-0.5" />
              <span class="pt-0.5">Folder:</span>
              <div class="flex items-center gap-6 border border-line px-1 py-0.5">{{ f.folder || "<None>" }}
                <a class="cursor-pointer underline" @click="folderEdit = { none: !f.folder, path: f.folder }">Change...</a></div>
              <span class="pt-0.5">Tags:</span>
              <div class="min-h-20 border border-line px-1 py-1">
                <span v-for="t in f.tags" :key="t" class="mr-1 mb-1 inline-block border border-neutral-300 bg-neutral-100 px-1">{{ t }}</span>
                <span v-if="!f.tags.length">&lt;None&gt;</span>
                <div class="mt-2 flex items-center gap-1"><img src="/icons/000_Tag_h32bit_16.png" alt="" /><a class="cursor-pointer underline" @click="openTags">Edit tags...</a></div>
              </div>
              <template v-if="orig.iqn != null"><label for="opiqn" class="pt-0.5">iSCSI IQN:</label><input id="opiqn" v-model="f.iqn" class="border border-line px-1 py-0.5" /></template>
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

    <div v-if="folderEdit" class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
      <form novalidate class="w-[28rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="chooseFolder">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Change Folder</div>
        <div class="space-y-2 p-3">
          <label class="flex items-center gap-2"><input v-model="folderEdit.none" type="radio" :value="true" />Do not place this resource in any folder</label>
          <label class="flex items-center gap-2"><input v-model="folderEdit.none" type="radio" :value="false" />Place this resource in this folder:</label>
          <div class="ml-5 h-40 overflow-auto border border-line">
            <div v-for="x in orig?.folders ?? []" :key="x" class="cursor-default px-1 py-0.5" :class="!folderEdit.none && folderEdit.path === x ? 'bg-select text-white' : 'hover:bg-neutral-100'"
                 @click="folderEdit!.path = x; folderEdit!.none = false">{{ x }}</div>
          </div>
          <div class="ml-5 flex items-center gap-2"><label for="newfolder" class="w-12">New...</label>
            <input id="newfolder" v-model="folderEdit.path" placeholder="/Folder/Subfolder" class="flex-1 border border-line px-1 py-0.5" @input="folderEdit!.none = false" /></div>
          <div class="flex justify-end gap-2 pt-1">
            <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5">OK</button>
            <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="folderEdit = null">Cancel</button>
          </div>
        </div>
      </form>
    </div>

    <div v-if="tagEdit" class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
      <div class="w-[26rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Edit Tags</div>
        <div class="space-y-2 p-3">
          <form novalidate class="flex gap-2" @submit.prevent="addTag">
            <input v-model="tagEdit.added" placeholder="New tag" aria-label="New tag" class="flex-1 border border-line px-1 py-0.5" />
            <button type="submit" class="btn" :disabled="!tagEdit.added.trim()">Create</button>
          </form>
          <div class="h-48 overflow-auto border border-line">
            <label v-for="t in tagEdit.all" :key="t" class="flex items-center gap-2 px-1 py-0.5 hover:bg-neutral-100">
              <input type="checkbox" :checked="tagEdit.on.has(t)" @change="($event.target as HTMLInputElement).checked ? tagEdit!.on.add(t) : tagEdit!.on.delete(t)" />{{ t }}</label>
          </div>
          <div class="flex justify-end gap-2 pt-1">
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="chooseTags">OK</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="tagEdit = null">Cancel</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
