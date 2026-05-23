package cache

import (
	"context"
	"time"
)

// Client is a simple string-key/string-value cache abstraction.
type Client interface {
	Get(ctx context.Context, key string) (value string, found bool, err error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
