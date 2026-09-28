# Tech Pulse

Tech-news aggregator: ingests RSS feeds from ~10 tech sources, normalizes and
deduplicates them into PostgreSQL, and serves a filterable/searchable dashboard.

## Architecture

```
┌──────────────┐   RSS   ┌──────────────┐   SQL   ┌────────────┐
│  RSS sources  │ ──────▶ │  ingestion    │ ──────▶ │ PostgreSQL │
└──────────────┘         │ (Go, cmd/)    │         └─────┬──────┘
                         └──────────────┘               │
                                                        ▼
┌──────────────┐   REST   ┌──────────────┐         ┌────────────┐
│  Next.js UI   │ ◀────── │  Go API       │ ──────▶ │   Redis     │
│  (dashboard)  │         │  (cmd/api)    │  cache  │  (optional) │
└──────────────┘         └──────────────┘         └────────────┘
```

- **Ingestion** (`cmd/ingestion`): fetches feeds from `configs/feeds.json`,
  normalizes into a common schema, dedupes by URL-hash, stores in Postgres.
- **API** (`cmd/api`): REST endpoints with clean architecture
  (handlers → services → repositories).
- **Dashboard** (`web/`): Next.js 15 App Router + Tailwind, dark/light toggle,
  PT/EN/ES language switcher, topic/source filters, search, sort.

## Stack

- Go 1.26, `net/http` routing, `slog` structured logging
- PostgreSQL 16 + `pgx/v5`, migrations with `golang-migrate`
- Redis 7 (response cache, graceful noop fallback)
- Next.js 15, React 19, Tailwind CSS, TypeScript
- Docker Compose for local dev

## Repository layout

```
tech-pulse/
├── cmd/
│   ├── api/          # REST API entrypoint
│   ├── ingestion/    # RSS ingestion entrypoint
│   └── migrate/      # golang-migrate runner
├── internal/
│   ├── api/          # HTTP handlers
│   ├── cache/        # Redis cache (noop fallback)
│   ├── config/       # env config
│   ├── ingest/       # fetch/normalize/dedupe
│   ├── model/        # shared domain types
│   ├── repository/   # interfaces + postgres impl
│   ├── service/      # application logic
│   └── store/        # pgx pool bootstrap
├── migrations/       # SQL migrations
├── configs/feeds.json
├── web/              # Next.js dashboard
└── build/            # Dockerfiles
```

## Quick start (Docker)

```bash
cp .env.example .env
make dev          # postgres+redis, then api + web
```

- Dashboard: http://localhost:3000
- API: http://localhost:8080/api/v1/news
- Health: http://localhost:8080/healthz

Run migrations and one ingestion pass (also possible on host):

```bash
make migrate      # apply DB migrations (needs local DATABASE_URL)
make ingest       # run a single ingestion pass
```

## Local development (without Docker for Go)

```bash
cp .env.example .env
docker compose up -d postgres redis   # infra only
make migrate
make ingest
go run ./cmd/api                      # API on :8080

cd web
npm install
npm run dev                           # dashboard on :3000
```

## API

### `GET /api/v1/news`

Query parameters:

| Param     | Description                                | Example                     |
|-----------|--------------------------------------------|-----------------------------|
| `topics`  | Comma-separated topic filter               | `topics=security,cloud`     |
| `sources` | Comma-separated source filter              | `sources=InfoQ`             |
| `q`       | Keyword search (title + summary, ILIKE)    | `q=kubernetes`              |
| `sort`    | `date` (default) or `source`               | `sort=source`               |
| `order`   | `desc` (default) or `asc`                  | `order=asc`                 |
| `limit`   | Page size (1–100, default 20)              | `limit=10`                  |
| `offset`  | Offset (default 0)                         | `offset=20`                 |

Response:

```json
{
  "items": [
    {
      "id": 1,
      "url": "https://…",
      "title": "…",
      "summary": "…",
      "source": "InfoQ",
      "topics": ["DevOps/SRE"],
      "published_at": "2026-09-27T12:00:00Z",
      "language": "en"
    }
  ],
  "total": 42,
  "page": 1,
  "page_size": 20,
  "total_pages": 3
}
```

## Topics

`AI/ML`, `Cloud`, `Software Engineering`, `DevOps/SRE`, `Security`, `Career/Market`.

For the MVP, titles/summaries are served in the original language; translations
are a later milestone.

## Make targets

| Target     | Description                                   |
|------------|-----------------------------------------------|
| `make dev` | Start infra + build/run API and web (Docker)  |
| `make build` | Build Go binaries + Next.js app             |
| `make test` | `go test ./...` + frontend lint              |
| `make lint` | `go vet`/golangci-lint + frontend lint       |
| `make migrate` | Apply migrations (host-local DB)          |
| `make ingest`  | One ingestion pass (host-local DB)        |
| `make up`/`down` | Start/stop docker-compose stack          |

## Conventions

- Conventional commits (`feat:`, `fix:`, `docs:`, `chore:`).
- See `AGENTS.md` for AI-tool guidance.
- No real secrets in the repo; `.env` is gitignored, use `.env.example` as template.
