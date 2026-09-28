package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/luacarol/tech-pulse/internal/model"
	"github.com/luacarol/tech-pulse/internal/repository"
)

// Runner orchestrates fetching all configured feeds and storing new items.
type Runner struct {
	feeds  []model.FeedSource
	fetch  *Fetcher
	repo   repository.NewsRepository
	logger *slog.Logger
}

func NewRunner(feedsPath string, repo repository.NewsRepository, logger *slog.Logger, timeout time.Duration) (*Runner, error) {
	feeds, err := loadFeeds(feedsPath)
	if err != nil {
		return nil, err
	}
	return &Runner{
		feeds:  feeds,
		fetch:  NewFetcher(timeout),
		repo:   repo,
		logger: logger,
	}, nil
}

func loadFeeds(path string) ([]model.FeedSource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read feeds config: %w", err)
	}
	var feeds []model.FeedSource
	if err := json.Unmarshal(data, &feeds); err != nil {
		return nil, fmt.Errorf("parse feeds config: %w", err)
	}
	return feeds, nil
}

// RunOnce performs a single ingestion pass and returns counts.
func (r *Runner) RunOnce(ctx context.Context) (fetched, stored int, err error) {
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		firstErr error
	)

	for _, src := range r.feeds {
		wg.Add(1)
		go func(src model.FeedSource) {
			defer wg.Done()

			items, ferr := r.fetch.Fetch(ctx, src.URL, src)
			if ferr != nil {
				r.logger.Warn("feed fetch failed", "source", src.Name, "error", ferr)
				mu.Lock()
				if firstErr == nil {
					firstErr = ferr
				}
				mu.Unlock()
				return
			}

			var localStored int
			for _, n := range items {
				exists, eerr := r.repo.Exists(ctx, n.URLHash)
				if eerr != nil {
					r.logger.Warn("dedupe check failed", "url", n.URL, "error", eerr)
					continue
				}
				if exists {
					continue
				}
				if uerr := r.repo.Upsert(ctx, n); uerr != nil {
					r.logger.Warn("store failed", "url", n.URL, "error", uerr)
					continue
				}
				localStored++
			}

			mu.Lock()
			fetched += len(items)
			stored += localStored
			mu.Unlock()

			r.logger.Info("feed ingested",
				"source", src.Name,
				"fetched", len(items),
				"stored", localStored,
			)
		}(src)
	}
	wg.Wait()

	return fetched, stored, firstErr
}
