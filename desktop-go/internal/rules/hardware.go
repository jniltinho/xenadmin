package rules

import "sort"

// The GPU and USB tabs (XenAdmin.Desktop.Core/HardwareRules.cs: GpuPage, GpuConfiguration, UsbPage).

// VGPUType is what the GPU tab shows of a vGPU type.
type VGPUType struct {
	Ref          string
	Model        string // model_name; "passthrough" for a whole GPU (VGPU_type.IsPassthrough)
	Capacity     int64  // vGPUs per GPU (VGPU_type.Capacity)
	Resolution   int64  // max_resolution_x * max_resolution_y
	Framebuffer  int64
	Compatible   []string // compatible_types_in_vm: the types another vGPU of the VM may have beside this one
	Experimental bool
}

// Passthrough is VGPU_type.IsPassthrough.
func (t VGPUType) Passthrough() bool { return t.Model == "passthrough" }

// SortVGPUTypes is GpuConfiguration's OrderByDescending(VGPU_type.CompareTo): whole-GPU pass-through first, then the
// types with fewer vGPUs per GPU, then the larger resolutions, then by name (natural order, descending).
func SortVGPUTypes(types []VGPUType) {
	sort.SliceStable(types, func(i, j int) bool {
		a, b := types[i], types[j]
		if a.Passthrough() != b.Passthrough() {
			return a.Passthrough()
		}
		if a.Capacity != b.Capacity {
			return a.Capacity < b.Capacity
		}
		if a.Resolution != b.Resolution {
			return a.Resolution > b.Resolution
		}
		return naturalLess(b.Model, a.Model)
	})
}

// ApplyVGPUChanges is HardwareRules.ConfigureAction for one GPU: only the boxes the user changed (type ref, now
// allowed) apply to its current allowed types, so changes made elsewhere meanwhile stay.
func ApplyVGPUChanges(enabled []string, changes map[string]bool) []string {
	out, has := []string{}, map[string]bool{}
	for _, r := range enabled {
		if allowed, ok := changes[r]; !ok || allowed {
			out, has[r] = append(out, r), true
		}
	}
	var added []string
	for r, allowed := range changes {
		if allowed && !has[r] {
			added = append(added, r)
		}
	}
	sort.Strings(added) // map order is random
	return append(out, added...)
}

// HostLess is Host.CompareTo (the names in natural order), with the ref as a tie-break so a server's GPUs stay
// together when names repeat.
func HostLess(nameA, refA, nameB, refB string) bool {
	if nameA != nameB {
		return naturalLess(nameA, nameB)
	}
	return refA < refB
}
