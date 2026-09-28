package model

import "time"

// Topic is a classification tag for a news item.
type Topic string

const (
	TopicAI         Topic = "AI/ML"
	TopicCloud      Topic = "Cloud"
	TopicSE         Topic = "Software Engineering"
	TopicDevOps     Topic = "DevOps/SRE"
	TopicSecurity   Topic = "Security"
	TopicCareer     Topic = "Career/Market"
)

// News is the normalized, common schema for every ingested item.
type News struct {
	ID          int64     `json:"id"`
	URLHash     string    `json:"-"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	Source      string    `json:"source"`
	Topics      []Topic   `json:"topics"`
	PublishedAt time.Time `json:"published_at"`
	Language    string    `json:"language"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewsFilter carries the query parameters for listing news.
type NewsFilter struct {
	Topics  []Topic
	Sources []string
	Query   string
	Sort    string // "date" (default) | "source"
	Order   string // "desc" (default) | "asc"
	Limit   int
	Offset  int
}

// PaginatedNews is a page of results with total count.
type PaginatedNews struct {
	Items      []News `json:"items"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}
