package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/Chateaubriand-g/bili/common/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	newclient := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MaxRetries:   3,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := newclient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("init client failed: %w", err)
	}

	return newclient, nil
}
