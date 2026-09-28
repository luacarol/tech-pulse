package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config holds environment-derived settings shared across binaries.
type Config struct {
	DatabaseURL  string
	RedisAddr    string
	APIAddr      string
	FeedsPath    string
	Interval     time.Duration
	FetchTimeout time.Duration
	CORSOrigins  []string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() Config {
	return Config{
		DatabaseURL:  getenv("DATABASE_URL", "postgres://techpulse:techpulse@localhost:5432/techpulse?sslmode=disable"),
		RedisAddr:    getenv("REDIS_ADDR", "localhost:6379"),
		APIAddr:      ":" + getenv("API_PORT", "8080"),
		FeedsPath:    getenv("INGESTION_FEEDS", "configs/feeds.json"),
		Interval:     getDuration("INGESTION_INTERVAL", 30*time.Minute),
		FetchTimeout: 15 * time.Second,
		CORSOrigins:  splitCSV(getenv("CORS_ORIGINS", "http://localhost:3000")),
	}
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		panic(fmt.Sprintf("invalid duration for %s: %q", key, v))
	}
	return d
}
