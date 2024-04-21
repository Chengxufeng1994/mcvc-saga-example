package broker

import (
	"fmt"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/common/config"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"

	nc "github.com/nats-io/nats.go"
)

var (
	logger    = watermill.NewStdLogger(false, false)
	marshaler = &nats.GobMarshaler{}
)

// NewNatsPublisher returns a NATS publisher for event streaming
func NewNatsPublisher(appCfg *config.ApplicationConfig) message.Publisher {
	var err error
	natsUrl := fmt.Sprintf("nats://%s:%d", appCfg.NatsConfig.Host, appCfg.NatsConfig.Port)
	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
	}
	jsConfig := nats.JetStreamConfig{
		Disabled:       false,
		AutoProvision:  true,
		ConnectOptions: nil,
		PublishOptions: nil,
		TrackMsgId:     false,
		AckAsync:       false,
		DurablePrefix:  appCfg.NatsConfig.NatsSubscriber.DurableName,
	}
	pub, err := nats.NewPublisher(
		nats.PublisherConfig{
			URL:         natsUrl,
			NatsOptions: options,
			Marshaler:   marshaler,
			JetStream:   jsConfig,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	// registry, ok := prom.DefaultRegisterer.(*prom.Registry)
	// if !ok {
	// 	return nil, fmt.Errorf("prometheus type casting error")
	// }
	// metricsBuilder := metrics.NewPrometheusMetricsBuilder(registry, config.App, "pubsub")
	// TxPublisher, err = metricsBuilder.DecoratePublisher(TxPublisher)
	// if err != nil {
	// 	return nil, err
	// }

	return pub
}
