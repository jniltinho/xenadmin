<script setup lang="ts">
// The Console tab (VNCTabView / XSVNCScreen): the VM's graphical console through noVNC. The window asks the Go side
// for a one-time loopback ticket (OpenConsole); the session and the server's certificate never reach the page.
import { onBeforeUnmount, onMounted, ref, watch } from "vue";
import * as novnc from "@novnc/novnc/lib/rfb";
import { api } from "./api/workspace";

// noVNC 1.6 ships CommonJS: its default export may come wrapped once more
const RFB = ((novnc as any).default?.default ?? (novnc as any).default) as any;

const props = defineProps<{ server: string; vm: string }>();
const screen = ref<HTMLDivElement>();
const state = ref("");
const scale = ref(true);
let rfb: any = null;
let generation = 0; // a newer connect (another VM, unmount) makes the older answers stale

async function connect() {
  disconnect();
  const mine = generation;
  state.value = "Connecting...";
  try {
    const url = await api.openConsole(props.server, props.vm);
    if (mine !== generation) return; // the ticket expires unused
    const r = new RFB(screen.value!, url, { wsProtocols: ["binary"] });
    rfb = r;
    r.scaleViewport = scale.value;
    r.focusOnClick = true;
    r.addEventListener("connect", () => { if (rfb === r) { state.value = ""; r.focus(); } });
    r.addEventListener("disconnect", (e: any) => {
      if (rfb !== r) return;
      state.value = e.detail?.clean ? "Disconnected." : "The console connection was lost.";
      rfb = null;
    });
  } catch (e) {
    if (mine === generation) state.value = String(e);
  }
}

function disconnect() { generation++; rfb?.disconnect(); rfb = null; }

watch(() => [props.server, props.vm], connect);
watch(scale, (v) => { if (rfb) rfb.scaleViewport = v; });
onMounted(connect);
onBeforeUnmount(disconnect);
</script>

<template>
  <div class="flex min-h-0 flex-1 flex-col">
    <div class="flex items-center gap-3 border-b border-line bg-pane px-2 py-1">
      <button class="border border-neutral-400 bg-white px-2 py-0.5 disabled:text-neutral-400" :disabled="!!state" @click="rfb?.sendCtrlAltDel()">Send Ctrl+Alt+Del</button>
      <label class="flex items-center gap-1"><input v-model="scale" type="checkbox" />Scale</label>
      <button v-if="state && state !== 'Connecting...'" class="border border-neutral-400 bg-white px-2 py-0.5" @click="connect">Reconnect</button>
      <span class="text-neutral-600">{{ state }}</span>
    </div>
    <div ref="screen" class="min-h-0 flex-1 bg-black"></div>
  </div>
</template>
