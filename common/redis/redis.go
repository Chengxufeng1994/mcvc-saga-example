package redis

import "context"

type RedisCache interface {
	Get(ctx context.Context, key string, dst interface{}) (bool, error)
	Set(ctx context.Context, key string, val interface{}, ttl int) error
	Del(ctx context.Context, key string) error
	Ping() error
	Close() error
}
