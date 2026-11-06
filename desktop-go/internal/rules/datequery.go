package rules

import "time"

// The date search (XenSearch DatePropertyQuery: start_time, the only date property).

// ParseDateQuery is DatePropertyQuery's "query" attribute: yyyyMMdd (an invalid one is an error, as ParseExact), or an
// older ISO 8601 date and time; else the zero time (DateTime.MinValue). Only its calendar date is used.
func ParseDateQuery(s string) (time.Time, error) {
	if len(s) == 8 {
		return time.Parse("20060102", s)
	}
	for _, layout := range []string{"20060102T15:04:05Z", "2006-01-02T15:04:05Z"} { // Util.TryParseIso8601DateTime
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, nil
}

// MatchDate is DatePropertyQuery.MatchDate: the value (UTC) in local time against the day queried; "this week" is
// today and the six days before, "last week" the seven before those, before and after include the day itself.
func MatchDate(value, query time.Time, kind string, now time.Time) bool {
	v := value.In(now.Location())
	day := func(t time.Time) time.Time { return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC) }
	diff := int(day(now).Sub(day(v)).Hours() / 24) // whole days back from today
	q := day(query)                                // query.Date: the date as written, no time zone change
	switch kind {
	case "today":
		return diff == 0
	case "yesterday":
		return diff == 1
	case "thisweek":
		return v.Year() >= now.Year()-1 && diff >= 0 && diff < 7
	case "lastweek":
		return v.Year() >= now.Year()-1 && diff >= 7 && diff < 14
	case "before":
		return !day(v).After(q)
	case "after":
		return !day(v).Before(q)
	case "exact":
		return day(v).Equal(q)
	}
	return false
}
