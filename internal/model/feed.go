package model

// FeedSource maps an RSS feed to its default topic(s).
type FeedSource struct {
	Name   string  `json:"name"`
	URL    string  `json:"url"`
	Topics []Topic `json:"topics"`
}
