package rules

import (
	"testing"
	"time"
)

func TestMatchDate(t *testing.T) {
	now := time.Date(2026, 1, 3, 10, 0, 0, 0, time.Local)
	at := func(y int, m time.Month, d int) time.Time { return time.Date(y, m, d, 12, 0, 0, 0, time.Local).UTC() }
	q, _ := ParseDateQuery("20251231")
	for _, c := range []struct {
		value time.Time
		kind  string
		want  bool
	}{
		{at(2026, 1, 3), "today", true},
		{at(2026, 1, 2), "today", false},
		{at(2026, 1, 2), "yesterday", true},
		{at(2025, 12, 28), "thisweek", true},
		{at(2025, 12, 27), "thisweek", false},
		{at(2025, 12, 27), "lastweek", true},
		{at(2025, 12, 21), "lastweek", true},
		{at(2025, 12, 20), "lastweek", false},
		{at(2025, 12, 31), "exact", true},
		{at(2025, 12, 31), "before", true},
		{at(2026, 1, 1), "before", false},
		{at(2025, 12, 31), "after", true},
		{at(2025, 12, 30), "after", false},
	} {
		if got := MatchDate(c.value, q, c.kind, now); got != c.want {
			t.Errorf("%v %s = %v", c.value, c.kind, got)
		}
	}
	if q, err := ParseDateQuery("bad"); err != nil || !q.IsZero() {
		t.Error("an unreadable date is not MinValue")
	}
	if _, err := ParseDateQuery("20260230"); err == nil {
		t.Error("an invalid yyyyMMdd date parsed")
	}
	if q, _ := ParseDateQuery("2026-01-01T00:00:00Z"); !MatchDate(at(2026, 1, 1), q, "exact", now) {
		t.Error("an ISO date moved to another day")
	}
}
