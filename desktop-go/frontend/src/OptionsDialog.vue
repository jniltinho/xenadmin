<script setup lang="ts">
// Tools > Options (OptionsDialog): Security (certificates, the pool secret reminder), Connection (proxy server and
// timeout), Confirmations (dismissing alerts) and Display (graph type, the last tab of each resource).
// shortcut: no "Warn me when a certificate changes": a changed certificate is always refused until the user removes
// the old pin (the app's trust rule).
import { onMounted, reactive, ref } from "vue";
import { api, type AppOptions } from "./api/workspace";

const emit = defineEmits<{ (e: "close", saved: boolean): void }>();
type Page = "security" | "connection" | "confirmations" | "display";
const pages: { id: Page; label: string; sub: string; icon: string }[] = [
  { id: "security", label: "Security", sub: "Certificate warnings", icon: "padlock" },
  { id: "connection", label: "Connection", sub: "Proxy server and timeout", icon: "000_Network_h32bit_16" },
  { id: "confirmations", label: "Confirmations", sub: "Dismissing alerts", icon: "075_TickRound_h32bit_16" },
  { id: "display", label: "Display", sub: "Graphs and tabs", icon: "001_PerformanceGraph_h32bit_16" },
];
const page = ref<Page>("security");
const f = reactive({ newCert: true, remind: true, proxy: "direct", address: "", port: "80", bypass: false, timeout: 20, dismiss: true, area: false, rememberTab: true });
const error = ref("");
const pending = ref(false);

onMounted(async () => {
  try {
    const o = await api.options();
    Object.assign(f, { newCert: !o.no_new_certificate_warning, remind: !o.no_rotate_secret_reminder, proxy: o.proxy || "direct", address: o.proxy_address ?? "",
      port: o.proxy_port ? String(o.proxy_port) : "80", // AppOptions.ProxyPort's default bypass: !!o.bypass_proxy_for_servers, timeout: o.connection_timeout || 20,
      dismiss: !o.no_dismiss_confirmation, area: !!o.area_graphs, rememberTab: !o.no_remember_last_tab });
  } catch (e) { error.value = String(e); }
});

async function ok() {
  if (pending.value) return;
  pending.value = true;
  try {
    const o: AppOptions = { no_new_certificate_warning: !f.newCert, no_rotate_secret_reminder: !f.remind, proxy: f.proxy, proxy_address: f.address.trim(),
      proxy_port: Number(f.port) || 0, bypass_proxy_for_servers: f.bypass, connection_timeout: Math.trunc(Number(f.timeout)) || 0,
      no_dismiss_confirmation: !f.dismiss, area_graphs: f.area, no_remember_last_tab: !f.rememberTab };
    await api.setOptions(o);
    emit("close", true);
  } catch (e) { error.value = String(e); } finally { pending.value = false; }
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <form novalidate class="flex h-[28rem] w-[44rem] flex-col border border-neutral-500 bg-white shadow-lg" @submit.prevent="ok">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">Options</div>
      <div class="flex min-h-0 flex-1">
        <div class="w-48 border-r border-line bg-pane py-1">
          <div v-for="p in pages" :key="p.id" class="flex cursor-default items-start gap-2 px-2 py-1.5" :class="page === p.id ? 'bg-select text-white' : 'hover:bg-neutral-100'" @click="page = p.id">
            <img :src="`/icons/${p.icon}.png`" alt="" class="mt-0.5 h-4 w-4" /><span>{{ p.label }}<span class="block text-xs opacity-75">{{ p.sub }}</span></span>
          </div>
        </div>
        <fieldset :disabled="pending" class="min-w-0 flex-1 space-y-3 overflow-auto p-4">
          <template v-if="page === 'security'">
            <p>You can configure XCP-ng Center to display a warning whenever a new or modified certificate is found on a managed server.</p>
            <label class="flex items-center gap-2"><input v-model="f.newCert" type="checkbox" />Warn me when a new certificate is found</label>
            <p>You can configure XCP-ng Center to display a reminder for changing the server root password whenever you rotate the pool secret.</p>
            <label class="flex items-center gap-2"><input v-model="f.remind" type="checkbox" />Remind me to change the server root password when I rotate the pool secret</label>
          </template>
          <template v-else-if="page === 'connection'">
            <p>These settings allow you to configure a proxy server or to customize your server connections for high latency environments.</p>
            <p class="font-semibold">Proxy server</p>
            <label class="flex items-center gap-2"><input v-model="f.proxy" type="radio" value="direct" />Don't use a proxy server</label>
            <label class="flex items-center gap-2"><input v-model="f.proxy" type="radio" value="system" />Use your operating system's proxy server settings</label>
            <label class="flex items-center gap-2"><input v-model="f.proxy" type="radio" value="specified" />Use this proxy server:</label>
            <div class="ml-6 flex items-center gap-2" :class="{ 'text-neutral-400': f.proxy !== 'specified' }">
              <label for="paddr">Address:</label><input id="paddr" v-model="f.address" :disabled="f.proxy !== 'specified'" class="w-56 border border-line px-1 py-0.5" />
              <label for="pport">Port:</label><input id="pport" v-model="f.port" :disabled="f.proxy !== 'specified'" class="w-20 border border-line px-1 py-0.5" />
            </div>
            <label class="ml-6 flex items-center gap-2" :class="{ 'text-neutral-400': f.proxy !== 'specified' }">
              <input v-model="f.bypass" type="checkbox" :disabled="f.proxy !== 'specified'" />Bypass proxy server for all server connections</label>
            <p class="font-semibold">Connection timeout</p>
            <p>Specify how long to wait when the connection to a managed server is disrupted.</p>
            <div class="flex items-center gap-2"><label for="ptime">Connection timeout:</label>
              <input id="ptime" v-model.number="f.timeout" type="number" min="1" max="2147483" class="w-28 border border-line px-1 py-0.5" /><span>seconds</span></div>
          </template>
          <template v-else-if="page === 'confirmations'">
            <p>You can configure XCP-ng Center to display a warning whenever you dismiss one or more alerts on the 'Notifications' view.</p>
            <label class="flex items-center gap-2"><input type="checkbox" :checked="!f.dismiss" @change="f.dismiss = !($event.target as HTMLInputElement).checked" />Don't show dismiss alerts confirmation dialog</label>
          </template>
          <template v-else>
            <p class="font-semibold">Graph type</p>
            <p>Configure the appearance of the performance graphs displayed on the Performance tab.</p>
            <label class="flex items-center gap-2"><input v-model="f.area" type="radio" :value="true" />Area graph</label>
            <label class="flex items-center gap-2"><input v-model="f.area" type="radio" :value="false" />Line graph</label>
            <p class="font-semibold">Tab options</p>
            <label class="flex items-center gap-2"><input v-model="f.rememberTab" type="checkbox" />Remember the last selected Tab for a resource</label>
          </template>
        </fieldset>
      </div>
      <div class="flex items-center justify-end gap-2 border-t border-line p-2">
        <span v-if="error" class="mr-auto text-red-600">{{ error }}</span>
        <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" :disabled="pending">OK</button>
        <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close', false)">Cancel</button>
      </div>
    </form>
  </div>
</template>
