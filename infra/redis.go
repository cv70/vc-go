package infra

import (
	"context"
	"vc-go/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context, c *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Host,
		Username: c.User,
		Password: c.Password,
		DB:       c.DB,
	})
	return client, nil
}
