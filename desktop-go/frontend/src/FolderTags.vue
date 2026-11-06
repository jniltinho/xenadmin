<script setup lang="ts">
// The General page's Folder and Tags rows (GeneralEditPage's FolderEditor and TagsEditor) with the Change Folder and
// Edit Tags dialogs (Dialogs.ChangeFolder / EditTags); folders and allTags are the connection's.
import { ref } from "vue";

const props = defineProps<{ folders: string[]; allTags: string[] }>();
const folder = defineModel<string>("folder", { required: true });
const tags = defineModel<string[]>("tags", { required: true });

// Change Folder: none, an existing folder, or a new path
const folderEdit = ref<{ none: boolean; path: string } | null>(null);
function chooseFolder() {
  const fe = folderEdit.value!;
  folder.value = fe.none ? "" : fe.path.trim();
  folderEdit.value = null;
}
// Edit Tags: the connection's tags ticked for this object, and new ones
const tagEdit = ref<{ all: string[]; on: Set<string>; added: string } | null>(null);
function openTags() { tagEdit.value = { all: [...new Set([...props.allTags, ...tags.value])].sort(), on: new Set(tags.value), added: "" }; }
function addTag() {
  const t = tagEdit.value!, name = t.added.trim();
  if (name && !t.all.includes(name)) t.all = [...t.all, name].sort();
  if (name) t.on.add(name);
  t.added = "";
}
function chooseTags() { tags.value = [...tagEdit.value!.on].sort(); tagEdit.value = null; }
</script>

<template>
  <span class="pt-0.5">Folder:</span>
  <div class="flex items-center gap-6 border border-line px-1 py-0.5">{{ folder || "<None>" }}
    <a class="cursor-pointer underline" @click="folderEdit = { none: !folder, path: folder }">Change...</a></div>
  <span class="pt-0.5">Tags:</span>
  <div class="min-h-20 border border-line px-1 py-1">
    <span v-for="t in tags" :key="t" class="mr-1 mb-1 inline-block border border-neutral-300 bg-neutral-100 px-1">{{ t }}</span>
    <span v-if="!tags.length">&lt;None&gt;</span>
    <div class="mt-2 flex items-center gap-1"><img src="/icons/000_Tag_h32bit_16.png" alt="" /><a class="cursor-pointer underline" @click="openTags">Edit tags...</a></div>
  </div>

  <Teleport to="body">
    <div v-if="folderEdit" class="fixed inset-0 z-50 flex items-center justify-center bg-black/20">
      <form novalidate class="w-[28rem] border border-neutral-500 bg-white shadow-lg" @submit.prevent="chooseFolder">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Change Folder</div>
        <div class="space-y-2 p-3">
          <label class="flex items-center gap-2"><input v-model="folderEdit.none" type="radio" :value="true" />Do not place this resource in any folder</label>
          <label class="flex items-center gap-2"><input v-model="folderEdit.none" type="radio" :value="false" />Place this resource in this folder:</label>
          <div class="ml-5 h-40 overflow-auto border border-line">
            <div v-for="x in folders" :key="x" class="cursor-default px-1 py-0.5" :class="!folderEdit.none && folderEdit.path === x ? 'bg-select text-white' : 'hover:bg-neutral-100'"
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
    <div v-if="tagEdit" class="fixed inset-0 z-50 flex items-center justify-center bg-black/20">
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
  </Teleport>
</template>
