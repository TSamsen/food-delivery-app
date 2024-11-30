package adapters

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisClient(address string) *RedisAdapter {
	return &RedisAdapter{
		client: redis.NewClient(&redis.Options{Addr: address}),
	}
}

func (r *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisAdapter) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}
