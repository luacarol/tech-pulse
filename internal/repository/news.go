package repository

import (
	"context"

	"github.com/luacarol/tech-pulse/internal/model"
)

// NewsRepository abstracts persistence of news items.
type NewsRepository interface {
	// Upsert inserts a news item, deduplicating by URL hash.
	Upsert(ctx context.Context, n model.News) error
	// Exists reports whether a URL hash is already stored.
	Exists(ctx context.Context, urlHash string) (bool, error)
	// List returns a page of news matching the filter.
	List(ctx context.Context, f model.NewsFilter) ([]model.News, int64, error)
}
