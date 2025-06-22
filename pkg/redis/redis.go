package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"simple-product-api/pkg/config"
)

func NewRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress,
	})
}

func PingRedis(rdb *redis.Client) error {
	return rdb.Ping(context.Background()).Err()
}
