package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRedisCmdable struct {
	getValue string
	getErr   error
	setErr   error
	delErr   error
}

func (f *fakeRedisCmdable) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	if f.getErr != nil {
		cmd.SetErr(f.getErr)
		return cmd
	}
	cmd.SetVal(f.getValue)
	return cmd
}

func (f *fakeRedisCmdable) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	if f.setErr != nil {
		cmd.SetErr(f.setErr)
		return cmd
	}
	cmd.SetVal("OK")
	return cmd
}

func (f *fakeRedisCmdable) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	if f.delErr != nil {
		cmd.SetErr(f.delErr)
		return cmd
	}
	cmd.SetVal(1)
	return cmd
}

func TestRedisClientGetReturnsHit(t *testing.T) {
	client := &RedisClient{client: &fakeRedisCmdable{getValue: "cached"}}
	value, found, err := client.Get(context.Background(), "k")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "cached", value)
}

func TestRedisClientGetReturnsMissOnNil(t *testing.T) {
	client := &RedisClient{client: &fakeRedisCmdable{getErr: redis.Nil}}
	value, found, err := client.Get(context.Background(), "k")
	require.NoError(t, err)
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestRedisClientGetReturnsError(t *testing.T) {
	client := &RedisClient{client: &fakeRedisCmdable{getErr: errors.New("boom")}}
	_, _, err := client.Get(context.Background(), "k")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}

func TestRedisClientSetAndDelete(t *testing.T) {
	client := &RedisClient{client: &fakeRedisCmdable{}}
	require.NoError(t, client.Set(context.Background(), "k", "v", time.Hour))
	require.NoError(t, client.Delete(context.Background(), "k"))
}
