package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient(env EnvDevType) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     env.REDIS_ADDR,
		Password: env.REDIS_PASS,
		DB:       env.REDIS_DB,
	})
	return rdb
}
