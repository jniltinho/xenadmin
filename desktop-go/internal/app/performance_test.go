package app

import (
	"math"
	"strings"
	"testing"
)

func TestRRDUpdatesAndAxes(t *testing.T) {
	xport := `<xport><meta><legend><entry>AVERAGE:vm:u:cpu0</entry><entry>AVERAGE:vm:u:memory</entry></legend></meta>
<data><row><t>200</t><v>0.5</v><v>NaN</v></row><row><t>100</t><v>0.25</v><v>1024</v></row></data></xport>`
	u, err := parseUpdates(strings.NewReader(xport))
	if err != nil || len(u.times) != 2 || u.times[0] != 100 {
		t.Fatalf("%+v %v", u, err)
	}
	cpu := scaled(u, "vm:u:cpu0", map[string]string{"cpu0": "(fraction)"})
	if cpu[0] != 25 || cpu[1] != 50 { // rows sorted by time, a fraction shown as %
		t.Errorf("cpu %v", cpu)
	}
	if mem := scaled(u, "vm:u:memory", nil); mem[0] != 1024 || !math.IsNaN(mem[1]) {
		t.Errorf("memory %v", mem)
	}
	if missing := scaled(u, "vm:u:absent", nil); !math.IsNaN(missing[0]) {
		t.Error("a series the server lacks is a gap")
	}
	if top, _, label := axis(unitPercent, 0, 37); top != 100 || label != "%" {
		t.Errorf("percent axis %v %v", top, label)
	}
	if top, div, label := axis(unitBytesPerSecond, 0, 3000); top != 10*1024 || div != 1024 || label != "kB/s" {
		t.Errorf("bytes axis %v %v %v", top, div, label)
	}
	if top, _, label := axis(unitBytes, 512<<20, 1); top != 512<<20 || label != "MB" {
		t.Errorf("fixed memory axis %v %v", top, label)
	}
}
