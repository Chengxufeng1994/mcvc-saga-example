package redis

import (
	"context"
	"strings"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/common/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient redis.UniversalClient
)

const (
	minRetryBackoff = 8 * time.Millisecond
	maxRetryBackoff = 512 * time.Millisecond
	dialTimeout     = 5 * time.Second
	readTimeout     = 3 * time.Second
	writeTimeout    = 3 * time.Second
	delimiter       = ","
)

func NewRedisClusterClient(libconfig *config.ApplicationConfig) redis.UniversalClient {
	poolTimeout := time.Duration(libconfig.RedisConfig.PoolTimeout) * time.Second
	RedisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           getServerAddrs(libconfig.RedisConfig.Addrs),
		Password:        libconfig.RedisConfig.Password,
		PoolSize:        libconfig.RedisConfig.PoolSize,
		PoolTimeout:     poolTimeout,
		MaxRetries:      libconfig.RedisConfig.MaxRetries,
		MinRetryBackoff: minRetryBackoff,
		MaxRetryBackoff: maxRetryBackoff,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		// To route commands by latency or randomly, enable one of the following.
		ReadOnly: libconfig.RedisConfig.ReadOnly,
		// RouteRandomly:  true,
		// RouteByLatency: true,
	})
	ctx := context.Background()
	err := RedisClient.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err == redis.Nil || err != nil {
		panic(err)
	}

	err = RedisClient.ForEachSlave(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err == redis.Nil || err != nil {
		panic(err)
	}

	redisotel.InstrumentTracing(RedisClient)

	return RedisClient
}

func getServerAddrs(addrs string) []string {
	return strings.Split(addrs, ",")
}
