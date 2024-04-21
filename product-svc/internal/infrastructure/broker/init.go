package broker

import (
	"time"

	"github.com/Chengxufeng1994/go-saga-example/common/bootstrap"
	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	prom "github.com/prometheus/client_golang/prometheus"
)

// InitializeRouter factory
func InitializeRouter(bootCfg *bootstrap.BootstrapConfig) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	// SignalsHandler will gracefully shutdown Router when SIGTERM is received.
	// You can also close the router by just calling `r.Close()`.
	router.AddPlugin(plugin.SignalsHandler)

	registry, ok := prom.DefaultRegisterer.(*prom.Registry)
	if !ok {
		panic("prometheus type casting error")
	}
	metricsBuilder := metrics.NewPrometheusMetricsBuilder(registry, bootCfg.Application, "pubsub")
	metricsBuilder.AddPrometheusRouterMetrics(router)

	// Router level middleware are executed for every message sent to the router
	router.AddMiddleware(
		// CorrelationID will copy the correlation id from the incoming message's metadata to the produced messages
		middleware.CorrelationID,
		// Timeout makes the handler cancel the incoming message's context after a specified time
		middleware.Timeout(time.Second*15),
		// The handler function is retried if it returns an error.
		// After MaxRetries, the message is Nacked and it's up to the PubSub to resend it.
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,
		// Recoverer handles panics from handlers.
		// In this case, it passes them as errors to the Retry middleware.
		middleware.Recoverer,
	)
	return router
}
