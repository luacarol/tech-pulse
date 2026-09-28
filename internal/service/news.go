package service

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/luacarol/tech-pulse/internal/model"
	"github.com/luacarol/tech-pulse/internal/repository"
)

const defaultPageSize = 20

// NewsService is the application layer for news queries.
type NewsService struct {
	repo repository.NewsRepository
}

func NewNewsService(repo repository.NewsRepository) *NewsService {
	return &NewsService{repo: repo}
}

// List returns a page of news, validating/normalizing the filter.
func (s *NewsService) List(ctx context.Context, f model.NewsFilter) (model.PaginatedNews, error) {
	if f.Limit <= 0 {
		f.Limit = defaultPageSize
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	if f.Order != "asc" {
		f.Order = "desc"
	}
	if f.Sort != "source" {
		f.Sort = "date"
	}

	items, total, err := s.repo.List(ctx, f)
	if err != nil {
		return model.PaginatedNews{}, fmt.Errorf("list news: %w", err)
	}

	page := (f.Offset / f.Limit) + 1
	totalPages := int(math.Ceil(float64(total) / float64(f.Limit)))

	return model.PaginatedNews{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   f.Limit,
		TotalPages: totalPages,
	}, nil
}

// ParseTopics splits a comma-separated topic filter into valid Topic values.
func ParseTopics(raw string) []model.Topic {
	if raw == "" {
		return nil
	}
	valid := map[string]model.Topic{
		"ai/ml":                model.TopicAI,
		"cloud":                model.TopicCloud,
		"software engineering": model.TopicSE,
		"devops/sre":           model.TopicDevOps,
		"security":             model.TopicSecurity,
		"career/market":        model.TopicCareer,
	}
	seen := map[model.Topic]bool{}
	var out []model.Topic
	for _, part := range strings.Split(raw, ",") {
		t, ok := valid[strings.ToLower(strings.TrimSpace(part))]
		if !ok || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

// ParseSources splits a comma-separated source filter.
func ParseSources(raw string) []string {
	if raw == "" {
		return nil
	}
	var out []string
	for _, part := range strings.Split(raw, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}
