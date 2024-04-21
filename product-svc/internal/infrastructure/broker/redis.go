package broker

import (
	"context"
	"strings"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/common/bootstrap"
	"github.com/Chengxufeng1994/go-saga-example/common/config"
	"github.com/Chengxufeng1994/go-saga-example/common/event"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

const (
	minRetryBackoff = 8 * time.Millisecond
	maxRetryBackoff = 512 * time.Millisecond
	dialTimeout     = 5 * time.Second
	readTimeout     = 3 * time.Second
	writeTimeout    = 3 * time.Second
	delimiter       = ","
)

type RedisPublisher struct {
	redisClient    redis.UniversalClient
	redisPublisher message.Publisher
}

func (rp RedisPublisher) GetPublisher() message.Publisher {
	return rp.redisPublisher
}

func (rp RedisPublisher) GetRedisClient() redis.UniversalClient {
	return rp.redisClient
}

func NewRedisPublisher(bootCfg *bootstrap.BootstrapConfig, libconfig *config.ApplicationConfig) *RedisPublisher {
	rp := &RedisPublisher{}

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

	rp.redisClient = RedisClient
	redisotel.InstrumentTracing(rp.redisClient)

	// TODO: purchaseResultTopic mas lens in configuration
	publisherConfig := redisstream.PublisherConfig{
		Client:     rp.redisClient,
		Marshaller: &redisstream.DefaultMarshallerUnmarshaller{},
		Maxlens: map[string]int64{
			event.PurchaseResultTopic: 5000,
		},
	}
	rp.redisPublisher, err = redisstream.NewPublisher(publisherConfig, logger)
	if err != nil {
		panic("failed to new redis publisher")
	}

	// registry, ok := prom.DefaultRegisterer.(*prom.Registry)
	// if !ok {
	// 	panic("failed to prometheus type casting")
	// }
	// metricsBuilder := metrics.NewPrometheusMetricsBuilder(registry, bootCfg.Application, "pubsub")
	// pub, err := metricsBuilder.DecoratePublisher(rp.redisPublisher)
	// if err != nil {
	// 	panic("failed to decorate publisher")
	// }
	// rp.redisPublisher = pub

	return rp
}

func getServerAddrs(addrs string) []string {
	return strings.Split(addrs, ",")
}
