<script setup lang="ts">
// BondDetails (XenAdmin/Controls/BondDetails.cs, XenAdmin.Desktop/BondEditor.cs): the NICs to bond (ticked up to the
// pool's limit; NICs already in a bond are dimmed), the bond mode, the MTU and "add to new VMs". Used by Create Bond
// and the New Network wizard's Bonded Network.
import { computed, ref, watch, watchEffect } from "vue";
import type { BondChoices, BondSpec } from "./api/workspace";
import { clampMtu, MTU_MAX, MTU_MIN, mtuCheck } from "./netrules";

const props = defineProps<{ choices: BondChoices; disabled?: boolean }>();
const spec = defineModel<BondSpec>({ required: true });
const valid = defineModel<boolean>("valid", { default: false });

const mtu = ref(String(spec.value.mtu || MTU_MIN));
const ticked = computed(() => props.choices.nics.filter((n) => spec.value.members.includes(n.uuid)));
// BondRules.MaxMtu: up to the smallest MTU of the members
const maxMtu = computed(() => Math.min(MTU_MAX, ...ticked.value.map((n) => n.mtu)));
const check = computed(() => mtuCheck(mtu.value, MTU_MIN, maxMtu.value));
// BondRules.CanTick: ticked rows stay; the others while the limit allows
const canTick = (n: BondChoices["nics"][number]) => spec.value.members.includes(n.uuid) || (!n.in_bond && spec.value.members.length < props.choices.limit);
function tick(uuid: string, on: boolean) {
  const m = spec.value.members.filter((u) => u !== uuid);
  spec.value.members = on ? [...m, uuid] : m;
}
watch(maxMtu, (hi) => { mtu.value = hi > MTU_MIN ? clampMtu(mtu.value, hi) : String(hi); }, { immediate: true });
// active-active is balance-slb; the LACP modes carry their hashing algorithm
const mode = computed({
  get: () => (spec.value.mode === "lacp" ? `lacp:${spec.value.hashing}` : spec.value.mode),
  set: (v: string) => { const [m, h] = v.split(":"); spec.value.mode = m; spec.value.hashing = h ?? ""; },
});
watchEffect(() => {
  spec.value.mtu = Number(mtu.value);
  valid.value = spec.value.members.length >= 2 && spec.value.members.length <= props.choices.limit && check.value.ok; // BondRules.ValidSelection
});
</script>

<template>
  <fieldset :disabled="disabled" class="space-y-3">
    <div class="max-h-40 overflow-auto border border-neutral-500">
      <table class="w-full border-collapse">
        <thead><tr class="text-left"><th class="border-r border-b border-line px-2 py-0.5" />
          <th v-for="h in ['NIC', 'MAC', 'Link Status', 'Speed', 'Vendor']" :key="h" class="border-r border-b border-line px-2 py-0.5 font-normal">{{ h }}</th></tr></thead>
        <tbody>
          <tr v-for="n in choices.nics" :key="n.uuid" :class="canTick(n) ? '' : 'text-neutral-400'">
            <td class="px-2"><input type="checkbox" :aria-label="n.name" :checked="spec.members.includes(n.uuid)" :disabled="!canTick(n)"
                                    @change="tick(n.uuid, ($event.target as HTMLInputElement).checked)" /></td>
            <td class="px-2">{{ n.name }}</td><td class="px-2">{{ n.mac }}</td><td class="px-2">{{ n.link }}</td>
            <td class="px-2">{{ n.speed }}</td><td class="px-2">{{ n.vendor }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="space-y-1">
      <p class="font-semibold">Bond mode</p>
      <label class="flex items-center gap-2"><input v-model="mode" type="radio" value="balance-slb" />Active-active</label>
      <label class="flex items-center gap-2"><input v-model="mode" type="radio" value="active-backup" />Active-passive</label>
      <template v-if="choices.lacp">
        <label class="flex items-center gap-2"><input v-model="mode" type="radio" value="lacp:src_mac" />LACP with load balancing based on source MAC address</label>
        <label class="flex items-center gap-2"><input v-model="mode" type="radio" value="lacp:tcpudp_ports" />LACP with load balancing based on IP and port of source and destination</label>
      </template>
      <p v-if="spec.mode === 'lacp'" class="text-amber-700">LACP must also be configured on the switch ports</p>
    </div>
    <div class="flex items-center gap-2">
      <label for="bondmtu">MTU:</label>
      <input id="bondmtu" v-model="mtu" class="w-24 border border-line px-1 py-0.5" :disabled="maxMtu <= MTU_MIN" />
      <span :class="check.ok ? 'text-neutral-600' : 'text-red-600'">{{ check.text }}</span>
    </div>
    <label class="flex items-center gap-2"><input v-model="spec.auto" type="checkbox" />Automatically add this network to new virtual machines</label>
  </fieldset>
</template>
