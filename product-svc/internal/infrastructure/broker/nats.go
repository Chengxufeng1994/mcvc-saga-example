package broker

import (
	"fmt"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/common/bootstrap"
	"github.com/Chengxufeng1994/go-saga-example/common/config"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/ThreeDotsLabs/watermill/message"
	nc "github.com/nats-io/nats.go"
	prom "github.com/prometheus/client_golang/prometheus"
)

type NatsPublisher message.Publisher
type NatsSubscriber message.Subscriber

var (
	TxPublisher  NatsPublisher
	TxSubscriber NatsSubscriber
)

// NewNATSPublisher returns a NATS publisher for event streaming
func NewNATSPublisher(bootCfg *bootstrap.BootstrapConfig, appCfg *config.ApplicationConfig) NatsPublisher {
	var err error
	natsUrl := fmt.Sprintf("nats://%s:%d", appCfg.NatsConfig.Host, appCfg.NatsConfig.Port)
	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
	}

	jsConfig := nats.JetStreamConfig{
		AutoProvision: true,
		AckAsync:      true,
		TrackMsgId:    true,
		Disabled:      false,
	}

	TxPublisher, err = nats.NewPublisher(
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

	registry, ok := prom.DefaultRegisterer.(*prom.Registry)
	if !ok {
		panic("prometheus type casting error")
	}
	metricsBuilder := metrics.NewPrometheusMetricsBuilder(registry, bootCfg.Application, "pubsub")
	TxPublisher, err = metricsBuilder.DecoratePublisher(TxPublisher)
	if err != nil {
		panic("prometheus decorate publisher")
	}

	return TxPublisher
}

// NewNATSSubscriber returns a NATS subscriber for event streaming
func NewNATSSubscriber(bootCfg *bootstrap.BootstrapConfig, appCfg *config.ApplicationConfig) NatsSubscriber {
	var err error
	natsUrl := fmt.Sprintf("nats://%s:%d", appCfg.NatsConfig.Host, appCfg.NatsConfig.Port)
	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
	}

	subOpts := []nc.SubOpt{
		// nc.DeliverAll(),
		nc.DeliverNew(),
		nc.AckExplicit(),
	}

	jsConfig := nats.JetStreamConfig{
		AutoProvision:    true,
		AckAsync:         true,
		TrackMsgId:       true,
		Disabled:         false,
		SubscribeOptions: subOpts,
		DurablePrefix:    appCfg.NatsConfig.NatsSubscriber.DurableName,
	}

	TxSubscriber, err = nats.NewSubscriber(
		nats.SubscriberConfig{
			URL:              natsUrl,
			NatsOptions:      options,
			Unmarshaler:      marshaler,
			JetStream:        jsConfig,
			QueueGroupPrefix: appCfg.NatsConfig.NatsSubscriber.QueueGroup,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	registry, ok := prom.DefaultRegisterer.(*prom.Registry)
	if !ok {
		panic("prometheus type casting error")
	}
	metricsBuilder := metrics.NewPrometheusMetricsBuilder(registry, bootCfg.Application, "pubsub")
	TxSubscriber, err = metricsBuilder.DecorateSubscriber(TxSubscriber)
	if err != nil {
		panic("prometheus decorate publisher")
	}

	return TxSubscriber
}
