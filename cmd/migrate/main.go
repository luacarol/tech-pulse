package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/luacarol/tech-pulse/internal/config"
)

func main() {
	path := flag.String("path", "migrations", "path to migration files")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := config.Load()

	action := "up"
	if flag.NArg() > 0 {
		action = flag.Arg(0)
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", *path),
		cfg.DatabaseURL,
	)
	if err != nil {
		logger.Error("init migrate", "error", err)
		os.Exit(1)
	}

	switch action {
	case "up":
		err = m.Up()
	case "down":
		err = m.Steps(-1)
	case "drop":
		err = m.Drop()
	default:
		err = fmt.Errorf("unknown action %q (use up|down|drop)", action)
	}

	if err != nil && err.Error() != "no change" && err.Error() != "file does not exist" {
		logger.Error("migration failed", "action", action, "error", err)
		os.Exit(1)
	}
	logger.Info("migration applied", "action", action)
}
