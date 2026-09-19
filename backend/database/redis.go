package database

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	var opts *redis.Options

	// Hosted Redis (Upstash etc.) provides a full URL, including TLS and password.
	if url := os.Getenv("REDIS_URL"); url != "" {
		parsed, err := redis.ParseURL(url)
		if err != nil {
			return nil, err
		}
		opts = parsed
	} else {
		// Local development fallback.
		opts = &redis.Options{Addr: "localhost:6379"}
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}