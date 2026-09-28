package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luacarol/tech-pulse/internal/api"
	"github.com/luacarol/tech-pulse/internal/cache"
	"github.com/luacarol/tech-pulse/internal/config"
	"github.com/luacarol/tech-pulse/internal/repository/postgres"
	"github.com/luacarol/tech-pulse/internal/service"
	"github.com/luacarol/tech-pulse/internal/store"
)

func main() {
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

	_ = cache.New(cfg.RedisAddr) // best-effort; cache wiring is a later milestone

	repo := postgres.NewNewsRepo(pool)
	svc := service.NewNewsService(repo)
	handler := api.NewHandler(svc, logger, cfg.CORSOrigins)

	srv := &http.Server{
		Addr:    cfg.APIAddr,
		Handler: handler.Routes(),
	}

	go func() {
		logger.Info("api listening", "addr", cfg.APIAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api server", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("api shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
