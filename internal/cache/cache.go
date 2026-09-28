package cache

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

// Cache is a minimal key/value abstraction over Redis with a graceful noop fallback.
type Cache struct {
	client *redis.Client
	logger *slog.Logger
}

// New connects to Redis at addr. It never fails hard: on error it returns a
// Cache with a nil client whose methods become no-ops, so the API keeps working
// without Redis.
func New(addr string) *Cache {
	logger := slog.Default()
	c := &Cache{logger: logger}
	if addr == "" {
		return c
	}
	client := redis.NewClient(&redis.Options{Addr: addr})
	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Warn("redis unavailable, cache disabled", "addr", addr, "error", err)
		return c
	}
	c.client = client
	logger.Info("redis cache connected", "addr", addr)
	return c
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	if c.client == nil {
		return "", nil
	}
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (c *Cache) Set(ctx context.Context, key, val string, ttl int) error {
	if c.client == nil {
		return nil
	}
	return c.client.Set(ctx, key, val, duration(ttl)).Err()
}
