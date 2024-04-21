package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/common/config"
	"github.com/redis/go-redis/v9"
)

type ClusterClient struct {
	client *redis.ClusterClient
}

const (
	minRetryBackoff = 8 * time.Millisecond
	maxRetryBackoff = 512 * time.Millisecond
	dialTimeout     = 5 * time.Second
	readTimeout     = 3 * time.Second
	writeTimeout    = 3 * time.Second
	delimiter       = ","
)

func NewClusterClient(libconfig *config.ApplicationConfig) (RedisCache, error) {
	cc := &ClusterClient{}

	poolTimeout := time.Duration(libconfig.RedisConfig.PoolTimeout) * time.Second
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
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

	cc.client = rdb
	err := cc.Ping()
	if err != nil {
		return nil, err
	}

	return cc, nil
}

// Get gets the value for the given key.
func (c *ClusterClient) Get(ctx context.Context, key string, dst interface{}) (bool, error) {
	fmt.Println(key)
	result := c.client.Get(ctx, key)
	val, err := result.Result()
	fmt.Println(val)
	switch {
	case err == redis.Nil:
		return false, nil
	case err != nil:
		return false, err
	default:
		err := json.Unmarshal([]byte(val), dst)
		fmt.Println(err)
	}

	return true, nil
}

// Set stores the given value for the given key along with a
func (c *ClusterClient) Set(ctx context.Context, key string, val interface{}, ttl int) error {
	dat, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if err := c.client.Set(ctx, key, dat, time.Duration(ttl)*time.Second).Err(); err != nil {
		return err
	}
	return nil
}

// Del deletes the value for the given key.
func (c *ClusterClient) Del(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

// Ping check redis connection
func (c *ClusterClient) Ping() error {
	ctx := context.Background()
	err := c.client.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err == redis.Nil || err != nil {
		return err
	}

	err = c.client.ForEachSlave(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err == redis.Nil || err != nil {
		return err
	}

	return nil
}

// Close implements RedisCache.
func (c *ClusterClient) Close() error {
	return c.client.Close()
}

func getServerAddrs(addrs string) []string {
	return strings.Split(addrs, delimiter)
}
