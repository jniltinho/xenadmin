// The field checks of the New Network wizard, BondDetails and a network's Properties (NetworkRules.ValidateVlan,
// ValidateMtu, BondRules.MaxMtu), as rules/bond.go checks them again on the Go side.
export const MTU_MIN = 1500, MTU_MAX = 9216;

const int = (text: string) => (/^\s*-?\d+\s*$/.test(text) ? Number(text) : NaN);

export function mtuCheck(text: string, lo: number, hi: number): { ok: boolean; text: string } {
  const m = int(text);
  if (Number.isNaN(m)) return { ok: false, text: "Invalid number" };
  return { ok: m >= lo && m <= hi, text: lo === hi ? `Allowed MTU value: ${lo}` : `Allowed MTU range: ${lo} to ${hi}` };
}

export function vlanCheck(text: string, min: number, inUse: number[]): { ok: boolean; text: string } {
  const v = int(text);
  if (Number.isNaN(v)) return { ok: false, text: "Invalid number" };
  if (v < min || v > 4094) return { ok: false, text: `Allowed VLAN range: ${min} to 4094` };
  if (inUse.includes(v)) return { ok: false, text: "This VLAN number is already in use" };
  return { ok: true, text: v === 0 ? "VLAN 0 will receive all traffic not on any other VLAN" : "" };
}

// NetworkRules.VlanFor: the VLAN if free on the NIC, else the first free one from 1
export function vlanFor(current: number, inUse: number[]): number {
  if (!inUse.includes(current)) return current;
  for (let i = 1; i <= 4094; i++) if (!inUse.includes(i)) return i;
  return current;
}

// the NumericUpDown keeps its value inside lo..hi; a fixed range resets it to the only value
export const clampMtu = (text: string, hi: number) => String(Math.min(Math.max(Number.isNaN(int(text)) ? MTU_MIN : int(text), MTU_MIN), hi));
