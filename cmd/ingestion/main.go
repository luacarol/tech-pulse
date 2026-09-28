package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luacarol/tech-pulse/internal/config"
	"github.com/luacarol/tech-pulse/internal/ingest"
	"github.com/luacarol/tech-pulse/internal/repository/postgres"
	"github.com/luacarol/tech-pulse/internal/store"
)

func main() {
	once := flag.Bool("once", false, "run a single ingestion pass and exit")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := store.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := postgres.NewNewsRepo(pool)
	runner, err := ingest.NewRunner(cfg.FeedsPath, repo, logger, cfg.FetchTimeout)
	if err != nil {
		logger.Error("init runner", "error", err)
		os.Exit(1)
	}

	if *once {
		fetched, stored, ferr := runner.RunOnce(ctx)
		logger.Info("ingestion pass complete", "fetched", fetched, "stored", stored, "error", ferr)
		return
	}

	logger.Info("ingestion service started", "interval", cfg.Interval.String())
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	run := func() {
		fetched, stored, ferr := runner.RunOnce(ctx)
		if ferr != nil {
			logger.Warn("ingestion pass completed with errors", "fetched", fetched, "stored", stored, "error", ferr)
			return
		}
		logger.Info("ingestion pass complete", "fetched", fetched, "stored", stored)
	}

	run() // run immediately on start
	for {
		select {
		case <-ctx.Done():
			logger.Info("ingestion service shutting down")
			return
		case <-ticker.C:
			run()
		}
	}
}
