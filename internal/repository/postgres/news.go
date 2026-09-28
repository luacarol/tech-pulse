package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/luacarol/tech-pulse/internal/model"
	"github.com/luacarol/tech-pulse/internal/repository"
)

type NewsRepo struct {
	pool *pgxpool.Pool
}

func NewNewsRepo(pool *pgxpool.Pool) *NewsRepo {
	return &NewsRepo{pool: pool}
}

var _ repository.NewsRepository = (*NewsRepo)(nil)

func (r *NewsRepo) Upsert(ctx context.Context, n model.News) error {
	const q = `
INSERT INTO news (url_hash, url, title, summary, source, topics, published_at, language)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (url_hash) DO NOTHING`
	_, err := r.pool.Exec(ctx, q,
		n.URLHash, n.URL, n.Title, n.Summary, n.Source,
		topicsToStrings(n.Topics), n.PublishedAt, n.Language,
	)
	if err != nil {
		return fmt.Errorf("upsert news %q: %w", n.URL, err)
	}
	return nil
}

func (r *NewsRepo) Exists(ctx context.Context, urlHash string) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM news WHERE url_hash = $1)`
	var exists bool
	if err := r.pool.QueryRow(ctx, q, urlHash).Scan(&exists); err != nil {
		return false, fmt.Errorf("check exists %q: %w", urlHash, err)
	}
	return exists, nil
}

func (r *NewsRepo) List(ctx context.Context, f model.NewsFilter) ([]model.News, int64, error) {
	where, args := buildWhere(f)

	countQ := "SELECT count(*) FROM news" + where
	var total int64
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count news: %w", err)
	}

	order := "published_at DESC"
	if f.Sort == "source" {
		order = "source ASC, published_at DESC"
	}
	if f.Order == "asc" && f.Sort == "date" {
		order = "published_at ASC"
	}

	limit, offset := f.Limit, f.Offset
	if limit <= 0 {
		limit = 20
	}

	selectQ := fmt.Sprintf(`
SELECT id, url_hash, url, title, summary, source, topics, published_at, language, created_at
FROM news%s ORDER BY %s LIMIT $%d OFFSET $%d`,
		where, order, len(args)+1, len(args)+2)

	rows, err := r.pool.Query(ctx, selectQ, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, fmt.Errorf("list news: %w", err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, scanNews)
	if err != nil {
		return nil, 0, fmt.Errorf("collect news: %w", err)
	}
	return items, total, nil
}

func scanNews(row pgx.CollectableRow) (model.News, error) {
	var n model.News
	var topics []string
	if err := row.Scan(
		&n.ID, &n.URLHash, &n.URL, &n.Title, &n.Summary, &n.Source,
		&topics, &n.PublishedAt, &n.Language, &n.CreatedAt,
	); err != nil {
		return model.News{}, err
	}
	n.Topics = stringsToTopics(topics)
	return n, nil
}

func buildWhere(f model.NewsFilter) (string, []any) {
	var conds []string
	var args []any

	if len(f.Topics) > 0 {
		args = append(args, topicsToStrings(f.Topics))
		conds = append(conds, fmt.Sprintf("topics && $%d::text[]", len(args)))
	}
	if len(f.Sources) > 0 {
		args = append(args, f.Sources)
		conds = append(conds, fmt.Sprintf("source = ANY($%d::text[])", len(args)))
	}
	if q := strings.TrimSpace(f.Query); q != "" {
		args = append(args, "%"+q+"%")
		conds = append(conds, fmt.Sprintf("(title ILIKE $%d OR summary ILIKE $%d)", len(args), len(args)))
	}

	if len(conds) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}

func topicsToStrings(topics []model.Topic) []string {
	out := make([]string, len(topics))
	for i, t := range topics {
		out[i] = string(t)
	}
	return out
}

func stringsToTopics(s []string) []model.Topic {
	out := make([]model.Topic, len(s))
	for i, v := range s {
		out[i] = model.Topic(v)
	}
	return out
}
