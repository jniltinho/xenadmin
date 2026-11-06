package rules

import (
	"strings"
	"testing"
	"time"
)

func TestAlertsCSV(t *testing.T) {
	out := AlertsCSV([]AlertCSV{{Title: `a "b"`, Priority: 2, Description: "d", AppliesTo: "pool", Time: time.Date(2026, 1, 2, 15, 4, 0, 0, time.Local)},
		{Title: "x", Priority: 6}})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if lines[0] != "Title,Severity,Description,Applies To,TimeStamp" || lines[1] != `"a ""b""","2","d","pool","Jan 2, 2026 3:04 PM"` ||
		!strings.Contains(lines[2], `"Unknown"`) {
		t.Fatal(out)
	}
}
