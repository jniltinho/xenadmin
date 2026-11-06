package rules

import (
	"reflect"
	"testing"
)

func TestSortVGPUTypes(t *testing.T) {
	types := []VGPUType{{Model: "K120Q", Capacity: 32}, {Model: "passthrough"}, {Model: "K160Q", Capacity: 8, Resolution: 4},
		{Model: "K180Q", Capacity: 8, Resolution: 9}, {Model: "K140Q", Capacity: 16}, {Model: "K2Q", Capacity: 16}, {Model: "K10Q", Capacity: 16}}
	SortVGPUTypes(types)
	var got []string
	for _, x := range types {
		got = append(got, x.Model)
	}
	if want := []string{"passthrough", "K180Q", "K160Q", "K140Q", "K10Q", "K2Q", "K120Q"}; !reflect.DeepEqual(got, want) {
		t.Fatal(got)
	}
}

func TestApplyVGPUChanges(t *testing.T) {
	got := ApplyVGPUChanges([]string{"a", "b", "c"}, map[string]bool{"b": false, "e": true, "d": true, "a": true})
	if want := []string{"a", "c", "d", "e"}; !reflect.DeepEqual(got, want) {
		t.Fatal(got)
	}
	if got := ApplyVGPUChanges(nil, nil); len(got) != 0 {
		t.Fatal(got)
	}
}
