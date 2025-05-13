package cache

import "context"

type ICache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, data interface{}, second int) error
	Lock(ctx context.Context, key string, ttl int64, proses func(ctx context.Context) error) error
}
