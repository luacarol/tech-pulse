.PHONY: dev build test lint migrate ingest up down clean

# Load local env for host-run Go tools (migrate/ingestion/api).
# docker-compose reads .env automatically.
-include .env
export

dev: ## Start postgres+redis, then build & run API + web
	docker compose up -d postgres redis
	docker compose up --build api web

build: ## Build Go binaries and the Next.js app
	go build -o bin/api ./cmd/api
	go build -o bin/ingestion ./cmd/ingestion
	go build -o bin/migrate ./cmd/migrate
	cd web && npm run build

test: ## Run Go tests and frontend type/lint checks
	go test ./...
	cd web && npm run lint

lint: ## Lint Go and frontend
	( command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... ) || go vet ./...
	cd web && npm run lint

migrate: ## Apply DB migrations (host-local DATABASE_URL)
	go run ./cmd/migrate -path ./migrations up

ingest: ## Run a single ingestion pass (host-local DATABASE_URL)
	go run ./cmd/ingestion -once

up: ## Start the full stack via docker-compose
	docker compose up -d

down: ## Stop the stack
	docker compose down

clean:
	rm -rf bin web/.next web/node_modules
