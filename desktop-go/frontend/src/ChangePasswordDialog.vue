<script setup lang="ts">
// Change Server Password (ChangeServerPasswordDialog): the current root password, the new one twice; the Go side
// checks them against the connection (PoolRules.ChangePasswordError) and runs session.change_password.
import { onMounted, ref } from "vue";
import { api } from "./api/workspace";

const props = defineProps<{ server: string }>();
const emit = defineEmits<{ (e: "close", done: boolean): void }>();
const current = ref(""), next = ref(""), confirm = ref("");
const title = ref("Change Server Password");
const error = ref("");
const pending = ref(false);
onMounted(async () => { try { title.value = await api.changePasswordTitle(props.server); } catch { /* the default title */ } });

async function ok() {
  if (pending.value || !current.value || !next.value || !confirm.value) return;
  pending.value = true;
  error.value = "";
  try { await api.changePassword(props.server, current.value, next.value, confirm.value); emit("close", true); }
  catch (e) { error.value = String(e); } finally { pending.value = false; }
}
</script>

<template>
  <div class="fixed inset-0 z-30 flex items-center justify-center bg-black/20">
    <form novalidate class="w-[30rem] border border-neutral-500 bg-white shadow-lg" autocomplete="off" @submit.prevent="ok">
      <div class="bg-band-dark px-3 py-1 font-semibold text-white">{{ title }}</div>
      <fieldset :disabled="pending" class="space-y-3 p-3">
        <div class="grid grid-cols-[11rem_1fr] items-center gap-2">
          <label for="cpcur">Current root password:</label><input id="cpcur" v-model="current" type="password" class="border border-line px-1 py-0.5" />
          <label for="cpnew">New root password:</label><input id="cpnew" v-model="next" type="password" class="border border-line px-1 py-0.5" />
          <label for="cpre">Re-enter new password:</label><input id="cpre" v-model="confirm" type="password" class="border border-line px-1 py-0.5" />
        </div>
        <p class="text-neutral-600">Passwords are case sensitive. If you lose or forget the password, it cannot be recovered.</p>
        <p>To reduce the risk of unauthorized access to your system, it is recommended that you also rotate the pool secret.</p>
        <p v-if="error" class="text-red-600">{{ error }}</p>
        <div class="flex justify-end gap-2">
          <button type="submit" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5 disabled:text-neutral-400" :disabled="!current || !next || !confirm">OK</button>
          <button type="button" class="min-w-20 border border-neutral-400 bg-pane px-3 py-0.5" @click="emit('close', false)">Cancel</button>
        </div>
      </fieldset>
    </form>
  </div>
</template>
