package ports

import (
	"context"
	"errors"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Cacher interface {
	GenKey(ctx context.Context, s string) (string, error)
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}
