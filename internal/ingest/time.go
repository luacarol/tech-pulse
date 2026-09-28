package ingest

import "time"

// rfc822 and rfc3339 are the two most common date formats in RSS/Atom feeds.
var timeLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	"2006-01-02",
}

// parseTime attempts a best-effort parse across common feed date layouts.
func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	for _, layout := range timeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}
