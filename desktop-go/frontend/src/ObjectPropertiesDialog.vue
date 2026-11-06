<script setup lang="ts">
// Properties of a server, pool, SR or network (PropertiesDialog.PagesFor): General (name, description, folder, tags; a
// server's iSCSI IQN) and Custom Fields, with the Change Folder and Edit Tags dialogs (Dialogs.ChangeFolder/EditTags).
// shortcut: the Alerts, Multipathing, Log Destination and Network Settings pages are not ported (as in the C# app).
import { computed, onMounted, reactive, ref } from "vue";
import { api, type ObjectProps } from "./api/workspace";
import FolderTags from "./FolderTags.vue";

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
              <FolderTags v-model:folder="f.folder" v-model:tags="f.tags" :folders="orig.folders" :all-tags="orig.all_tags" />
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

  </div>
</template>
