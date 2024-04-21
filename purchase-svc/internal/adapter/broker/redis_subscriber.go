package broker

import (
	"github.com/Chengxufeng1994/go-saga-example/common/bootstrap"
	"github.com/Chengxufeng1994/go-saga-example/common/config"

	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/ThreeDotsLabs/watermill/message"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

var (
	Subscriber message.Subscriber
)

func NewRedisSubscriber(bootstrapCfg *bootstrap.BootstrapConfig, appCfg *config.ApplicationConfig, subClient redis.UniversalClient) message.Subscriber {
	var err error
	Subscriber, err = redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:        subClient,
			Unmarshaller:  redisstream.DefaultMarshallerUnmarshaller{},
			Consumer:      appCfg.RedisConfig.Subscriber.ConsumerID,
			ConsumerGroup: appCfg.RedisConfig.Subscriber.ConsumerGroup,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	registry, ok := prom.DefaultRegisterer.(*prom.Registry)
	if !ok {
		panic(err)
	}

	metricsBuilder := metrics.NewPrometheusMetricsBuilder(registry, bootstrapCfg.Application, "pubsub")
	Subscriber, err = metricsBuilder.DecorateSubscriber(Subscriber)
	if err != nil {
		panic(err)
	}

	return Subscriber
}
