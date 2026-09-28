# AGENTS.md — Tech Pulse

Guidance for AI coding tools (Claude Code, Copilot, etc.) working in this repo.

## What this is

Tech Pulse is a tech-news aggregator monorepo:

- `cmd/ingestion` — Go RSS ingestion service (fetch → normalize → dedupe → store).
- `cmd/api` — Go REST API (clean architecture: handlers → services → repositories).
- `cmd/migrate` — golang-migrate runner.
- `web/` — Next.js dashboard (App Router, Tailwind).
- `migrations/` — SQL migrations (golang-migrate, up/down pairs).
- `configs/feeds.json` — RSS source + topic mapping.

## Stack decisions

- Go 1.26, `net/http` stdlib routing (Go 1.22+ method patterns), `slog` for logging.
- PostgreSQL via `pgx/v5` (pgxpool). Migrations via `golang-migrate/v4`.
- Redis via `go-redis/v9` (response cache with graceful noop fallback).
- Next.js 15 App Router, Tailwind, TypeScript. Server components for data fetching;
  client components ONLY for interactive elements (filters, search, theme, language).

## Conventions

- **Clean architecture**: handlers depend on service interfaces; services depend on
  repository interfaces; PostgreSQL/Redis implementations live behind interfaces.
- **Logging**: always `slog` (never `fmt.Println` / `log`). Structured key/value pairs.
- **Config**: read from env vars only (see `.env.example`). No hardcoded credentials.
- **Errors**: wrap with context (`fmt.Errorf("fetch %s: %w", url, err)`), handle at
  boundaries, return proper HTTP status codes.
- **Migrations**: every schema change is a new `NNNNNN_name.up.sql` / `.down.sql` pair.
- **Commits**: conventional commits (`feat:`, `fix:`, `docs:`, `chore:`, `refactor:`).
- **No real secrets** in the repo; `.env` is gitignored, `.env.example` is the template.

## Common commands

- `make dev` — start postgres+redis, then API + web (Docker).
- `make build` — build Go binaries and the Next.js app.
- `make test` — `go test ./...` + frontend checks.
- `make migrate` — run DB migrations locally.
- `make ingest` — run one ingestion pass locally.

## Gotchas

- `DATABASE_URL` in `.env` points to `localhost` (host-local Go tools). Inside
  docker-compose, the URL is constructed with host `postgres` — do not pass the
  host-local URL into containers.
- Go module path is `github.com/luacarol/tech-pulse`; internal packages live under
  `internal/` and are not importable outside this module.
- Topic strings are fixed constants in `internal/model` — keep them in sync with
  `configs/feeds.json`.
