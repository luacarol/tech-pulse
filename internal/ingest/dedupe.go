package ingest

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
)

var htmlTag = regexp.MustCompile(`<[^>]*>`)
var ws = regexp.MustCompile(`\s+`)

// HashURL computes a deterministic dedupe key from a canonicalized URL.
func HashURL(url string) string {
	u := strings.TrimSpace(url)
	// strip common tracking/utm query params for dedupe stability
	u = stripTracking(u)
	sum := sha256.Sum256([]byte(u))
	return hex.EncodeToString(sum[:])
}

func stripTracking(u string) string {
	parts := strings.SplitN(u, "?", 2)
	if len(parts) != 2 {
		return u
	}
	q := parts[1]
	var keep []string
	for _, pair := range strings.Split(q, "&") {
		key := strings.ToLower(strings.SplitN(pair, "=", 2)[0])
		if strings.HasPrefix(key, "utm_") || strings.HasPrefix(key, "fbclid") || strings.HasPrefix(key, "gclid") {
			continue
		}
		keep = append(keep, pair)
	}
	if len(keep) == 0 {
		return parts[0]
	}
	return parts[0] + "?" + strings.Join(keep, "&")
}

func stripHTML(s string) string {
	return strings.TrimSpace(ws.ReplaceAllString(htmlTag.ReplaceAllString(s, " "), " "))
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return strings.TrimSpace(string(r[:n])) + "…"
}
