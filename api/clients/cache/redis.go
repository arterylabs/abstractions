package cache

import (
	"abstractions/config"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCmdable interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

// RedisClient is a redis-backed implementation of Client.
type RedisClient struct {
	client redisCmdable
}

func NewRedisClient(redisAddr string) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr: redisAddr,
		}),
	}
}

func NewRedisClientFromEnvironment() (*RedisClient, error) {
	if config.Environment == nil {
		return nil, fmt.Errorf("cache: redis: environment is not loaded")
	}
	if config.Environment.RedisHost == "" || config.Environment.RedisPort == 0 {
		return nil, fmt.Errorf("cache: redis: host/port not configured")
	}
	return NewRedisClient(config.Environment.RedisHost + ":" + strconv.Itoa(config.Environment.RedisPort)), nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, bool, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

func (r *RedisClient) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
