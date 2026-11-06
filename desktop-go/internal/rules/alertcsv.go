package rules

import (
	"strconv"
	"strings"
	"time"

	"xcpng-center/internal/messages"
)

// AlertCSV is an alert row as the CSV export writes it.
type AlertCSV struct {
	Title       string    `json:"title"`
	Priority    int       `json:"priority"` // 1..5, 6 unknown
	Description string    `json:"description"`
	AppliesTo   string    `json:"applies_to"`
	Time        time.Time `json:"time"`
}

// AlertsCSV is AlertSummaryPage's export (NotificationRules.CsvHeader / CsvLine, AlertExtensions.GetAlertDetailsCSVQuotes):
// a header, then every field quoted with its quotes doubled, the time in local time as "MMM d, yyyy h:mm tt".
func AlertsCSV(rows []AlertCSV) string {
	q := func(s string) string { return `"` + strings.ReplaceAll(s, `"`, `""`) + `"` }
	var b strings.Builder
	b.WriteString(strings.Join([]string{messages.Get("TITLE"), messages.Get("SEVERITY"), messages.Get("DESCRIPTION"),
		messages.Get("APPLIES_TO"), messages.Get("TIMESTAMP")}, ",") + "\n")
	for _, r := range rows {
		prio := "Unknown" // AlertPriority.GetString
		if r.Priority >= 1 && r.Priority <= 5 {
			prio = strconv.Itoa(r.Priority)
		}
		b.WriteString(strings.Join([]string{q(r.Title), q(prio), q(r.Description), q(r.AppliesTo), q(r.Time.Local().Format("Jan 2, 2006 3:04 PM"))}, ",") + "\n")
	}
	return b.String()
}
