<script setup lang="ts">
// New Search / Edit Search (XenAdmin.Desktop SearchEditor: SearchPage with QueryPanel): the name, the object types, a
// text filter and the grouping, then Apply (show it), Save to a pool (it joins Saved Searches) or, for a saved search,
// Delete. The Go side builds the XML (SearchRules.Build) and keeps the parts of an edited search not shown here.
import { computed, onMounted, ref } from "vue";
import { api, type SearchEditorForm } from "./api/workspace";

export type SearchResult = { kind: "apply" | "save" | "delete"; xml: string; name: string; uuid?: string };
// xml: the edited search ("" for New Search); uuid: set when it is saved in a pool (Delete, Save replaces it)
const props = defineProps<{ xml: string; uuid: string; server: string; isNew: boolean }>();
const emit = defineEmits<{ (e: "close", r: SearchResult | null): void }>();

const form = ref<SearchEditorForm | null>(null);
const name = ref(""), text = ref(""), grouping = ref(0), saveTo = ref("");
const checked = ref(new Set<string>());
const groupingChanged = ref(false); // an unchanged choice keeps the search's own grouping (it may be nested or not listed)
const error = ref("");
const pending = ref(false);
const askDelete = ref(false);
const title = computed(() => (props.isNew ? "New Search" : `Edit Search '${form.value?.name ?? ""}'`));

onMounted(async () => {
  try {
    const f = await api.editSearch(props.xml);
    form.value = f;
    name.value = f.name;
    checked.value = new Set(f.types ?? []);
    grouping.value = f.grouping;
    saveTo.value = f.pools.find((p) => p.key === props.server)?.key ?? f.pools[0]?.key ?? "";
  } catch (e) { error.value = String(e); }
});
function toggle(id: string) {
  const s = new Set(checked.value);
  if (!s.delete(id)) s.add(id);
  checked.value = s;
}
const build = () => api.buildSearch({ name: name.value, types: [...checked.value], text: text.value, grouping: grouping.value,
  grouping_changed: groupingChanged.value, edited: props.xml, saved: !!props.uuid });

async function run(f: () => Promise<SearchResult>) {
  if (pending.value || !form.value) return;
  pending.value = true;
  error.value = "";
  try { emit("close", await f()); } catch (e) { error.value = String(e); } finally { pending.value = false; }
}
const apply = () => run(async () => ({ kind: "apply", xml: await build(), name: name.value.trim() }));
function save() {
  if (!saveTo.value) { error.value = "Choose the pool to save the search to."; return; }
  run(async () => { const xml = await build(); return { kind: "save", xml, name: name.value.trim(), uuid: await api.saveSearch(saveTo.value, xml) }; });
}
function remove() { askDelete.value = false; run(async () => { await api.deleteSearch(props.uuid); return { kind: "delete", xml: "", name: "" }; }); }
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <form novalidate class="w-[42.5rem] border border-neutral-500 bg-white shadow-lg" autocomplete="off" @submit.prevent="apply">
      <div class="truncate bg-band-dark px-3 py-1 font-semibold text-white">{{ title }}</div>
      <fieldset :disabled="pending || !form" class="space-y-2.5 p-3">
        <div class="grid grid-cols-[110px_1fr] items-center"><label for="sename">Name:</label>
          <input id="sename" v-model="name" class="border border-line px-1 py-0.5" /></div>
        <p class="font-semibold">Search for:</p>
        <div class="grid grid-cols-3 gap-x-3.5 gap-y-1">
          <label v-for="t in form?.type_list ?? []" :key="t.id" class="flex items-center gap-1.5">
            <input type="checkbox" :checked="checked.has(t.id)" @change="toggle(t.id)" />{{ t.name }}</label>
        </div>
        <!-- shortcut: no IP address query (RunSearch does not read addresses), so the hint leaves addresses out -->
        <div class="grid grid-cols-[110px_1fr] items-center"><label for="setext">Filter:</label>
          <input id="setext" v-model="text" class="border border-line px-1 py-0.5"
                 :placeholder="form?.has_filter ? 'Add a text filter to the current conditions' : 'Name or description contains...'" /></div>
        <p v-if="form?.has_filter" class="text-neutral-500">The search's current conditions are kept.</p>
        <div class="grid grid-cols-[110px_1fr] items-center"><label for="segroup">Group by:</label>
          <select id="segroup" v-model.number="grouping" class="w-56 border border-line px-1 py-0.5" @change="groupingChanged = true">
            <option v-for="(g, i) in form?.groupings ?? []" :key="g.id" :value="i">{{ g.name }}</option>
          </select></div>
        <div class="grid grid-cols-[110px_1fr] items-center"><label for="sepool">Save to pool:</label>
          <select id="sepool" v-model="saveTo" class="w-56 border border-line px-1 py-0.5" :disabled="!form?.pools.length">
            <option v-for="p in form?.pools ?? []" :key="p.key" :value="p.key">{{ p.name }}</option>
          </select></div>
      </fieldset>
      <!-- outside the fieldset: Cancel works when the editor could not be filled -->
      <div class="space-y-2.5 px-3 pb-3">
        <p v-if="error" class="text-red-600">{{ error }}</p>
        <div class="flex gap-2">
          <button v-if="uuid" type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="pending || !form" @click="askDelete = true">Delete</button>
          <span class="flex-1" />
          <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="pending || !form">Apply</button>
          <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="pending || !form?.pools.length" @click="save">Save</button>
          <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="pending" @click="emit('close', null)">Cancel</button>
        </div>
      </div>
    </form>
    <div v-if="askDelete" class="fixed inset-0 z-40 flex items-center justify-center bg-black/20">
      <div class="w-[28rem] border border-neutral-500 bg-white shadow-lg">
        <div class="bg-band-dark px-3 py-1 font-semibold text-white">Delete Search</div>
        <div class="space-y-3 p-3"><p>Delete the saved search '{{ form?.name }}'?</p>
          <div class="flex justify-end gap-2">
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="remove">Yes</button>
            <button class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="askDelete = false">No</button>
          </div></div>
      </div>
    </div>
  </div>
</template>
