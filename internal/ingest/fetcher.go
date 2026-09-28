package ingest

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/luacarol/tech-pulse/internal/model"
)

// rssFeed is the minimal RSS/Atom subset we care about.
type rssFeed struct {
	Channel struct {
		Items []rssItem `xml:"item"`
	} `xml:"channel"`
	// Atom-style entries are also mapped into Items by the atom types below.
	Entries []rssItem `xml:"entry"`
}

type rssItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	GUID    string `xml:"guid"`
	PubDate string `xml:"pubDate"`
	Updated string `xml:"updated"`
	Summary string `xml:"description"`
	// Atom link alternative: <link href="..."/>
	LinkHref string `xml:"-"`
}

// Fetcher downloads and parses a feed URL into normalized News items.
type Fetcher struct {
	client *http.Client
}

func NewFetcher(timeout time.Duration) *Fetcher {
	return &Fetcher{client: &http.Client{Timeout: timeout}}
}

// Fetch retrieves the feed at feedURL and returns normalized items for the source.
func (f *Fetcher) Fetch(ctx context.Context, feedURL string, src model.FeedSource) ([]model.News, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "tech-pulse/0.1 (+https://github.com/luacarol/tech-pulse)")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", feedURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: unexpected status %d", feedURL, resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20)) // 5 MiB cap
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", feedURL, err)
	}

	var feed rssFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("parse %s: %w", feedURL, err)
	}

	items := feed.Channel.Items
	if len(items) == 0 {
		items = feed.Entries
	}

	out := make([]model.News, 0, len(items))
	for _, it := range items {
		n, ok := normalizeItem(it, src)
		if !ok {
			continue
		}
		out = append(out, n)
	}
	return out, nil
}

// normalizeItem converts a parsed RSS/Atom item into the common schema.
// It returns ok=false when the item lacks a usable link or title.
func normalizeItem(it rssItem, src model.FeedSource) (model.News, bool) {
	link := it.Link
	if link == "" {
		link = it.LinkHref
	}
	if link == "" {
		link = it.GUID
	}
	if link == "" || it.Title == "" {
		return model.News{}, false
	}

	pub := parseTime(it.PubDate)
	if pub.IsZero() {
		pub = parseTime(it.Updated)
	}
	if pub.IsZero() {
		pub = time.Now().UTC()
	}

	return model.News{
		URL:         link,
		URLHash:     HashURL(link),
		Title:       it.Title,
		Summary:     truncate(stripHTML(it.Summary), 280),
		Source:      src.Name,
		Topics:      src.Topics,
		PublishedAt: pub,
		Language:    "en", // original language; translations come later
	}, true
}
