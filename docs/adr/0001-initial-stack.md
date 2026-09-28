# ADR 0001 — Initial stack choices

- Status: accepted
- Date: 2026-09-27

## Context

Tech Pulse is a tech-news aggregator MVP: ingest ~10 RSS feeds, normalize and
deduplicate them, store in PostgreSQL, and serve a filterable dashboard. We need
to pick a stack that is fast to build, easy to run locally, and cheap to operate
as a solo/small project, while leaving room to grow.

## Decision

| Concern          | Choice                          | Rationale |
|------------------|---------------------------------|-----------|
| Ingestion + API  | Go 1.26, `net/http`             | Single binary, low ops cost, strong stdlib, `slog` for structured logs. |
| Database         | PostgreSQL 16 + `pgx/v5`        | Relational fits the normalized schema; `pgx` is the de-facto high-perf driver. |
| Migrations       | `golang-migrate/v4`             | Versioned up/down SQL, CLI + library, works well with Docker entrypoints. |
| Cache            | Redis 7 (`go-redis/v9`)         | Fast response cache; wrapped behind an interface with a noop fallback so the API runs without it. |
| Dashboard        | Next.js 15 App Router + React 19 | RSC for data fetching, client components only for interactivity. |
| Styling          | Tailwind CSS                    | Fast iteration, dark mode via `class` strategy. |
| Orchestration    | Docker Compose                  | Reproducible local env for Postgres/Redis/API/web. |

## Architecture

Clean architecture in the Go service:

```
handlers (net/http) → services (business logic) → repositories (interfaces)
                                                      └── postgres impl
```

Dependencies point inward; handlers depend on service interfaces, services on
repository interfaces. PostgreSQL and Redis implementations sit behind interfaces,
which keeps the application layer testable without real infra.

## Consequences

- **Pros**: low operational footprint, reproducible local dev, clear separation of
  concerns, easy to test the application layer with fake repositories.
- **Cons**: two languages (Go + TypeScript) to maintain; RSC/client split requires
  discipline about where interactivity lives.

## Alternatives considered

- **Python/FastAPI for ingestion**: quicker to prototype, but weaker concurrency
  story and larger runtime footprint for a long-running fetcher.
- **Node.js for the API**: one language across the stack, but Go's single-binary
  deployment and `slog`/`pgx` ergonomics were preferred for the data path.
- **MongoDB**: schema flexibility, but relational constraints (unique URL hash,
  array topic filters with GIN index) are a natural fit for Postgres.
