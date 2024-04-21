//go:build wireinject
// +build wireinject

package di

import (
	"github.com/Chengxufeng1994/go-saga-example/common/bootstrap"
	"github.com/Chengxufeng1994/go-saga-example/common/config"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/db"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/adapter/broker"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/adapter/repository"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/application"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure"
	infrabroker "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/broker"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/grpc/client"
	infragrpcproduct "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/grpc/product"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/middleware"
	httporder "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/order"
	httppayment "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/payment"
	httpproduct "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/product"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/observe"
	"github.com/google/wire"
)

func InitApplicationConfig(path string) *config.ApplicationConfig {
	wire.Build(
		config.LoadApplicationConfig,
	)
	return &config.ApplicationConfig{}
}

func InitBootstrapConfig(path string) *bootstrap.BootstrapConfig {
	wire.Build(
		bootstrap.LoadBootstrapConfig,
	)
	return &bootstrap.BootstrapConfig{}
}

func InitializeMigrator(app string, appCfg *config.ApplicationConfig) (*db.Migrator, error) {
	wire.Build(
		db.NewDatabase,
		db.NewMigrator,
	)
	return &db.Migrator{}, nil
}

func InitializeProductServer(appCfg *config.ApplicationConfig, bootCfg *bootstrap.BootstrapConfig) *infrastructure.ProductServer {
	wire.Build(
		db.NewDatabase,

		observe.NewTracer,

		infrabroker.InitializeRouter,
		infrabroker.NewNATSPublisher,
		infrabroker.NewNATSSubscriber,
		infragrpcproduct.NewGrpcProductServer,
		broker.NewSagaProductController,
		broker.NewProductEventRouter,

		repository.NewGormProductRepository,

		client.NewAuthConn,
		application.NewAuthService,
		application.NewProductService,
		application.NewSagaProductService,
		application.NewProductApplication,

		middleware.NewJwtAuthenticator,
		httpproduct.NewGinEngine,
		httpproduct.NewRouter,
		httpproduct.New,

		infrastructure.NewProductServer,
	)

	return &infrastructure.ProductServer{}
}

func InitializeOrderServer(appCfg *config.ApplicationConfig, bootCfg *bootstrap.BootstrapConfig) *infrastructure.OrderServer {
	wire.Build(
		db.NewDatabase,

		observe.NewTracer,

		infrabroker.InitializeRouter,
		infrabroker.NewNATSPublisher,
		infrabroker.NewNATSSubscriber,
		broker.NewSagaOrderController,
		broker.NewOrderEventRouter,

		repository.NewGormOrderRepository,

		client.NewAuthConn,
		client.NewProductConn,
		application.NewAuthService,
		application.NewOrderService,
		application.NewSagaOrderService,
		application.NewOrderApplication,

		middleware.NewJwtAuthenticator,
		httporder.NewGinEngine,
		httporder.NewRouter,
		httporder.New,

		infrastructure.NewOrderServer,
	)

	return &infrastructure.OrderServer{}
}

func InitializePaymentServer(appCfg *config.ApplicationConfig, bootCfg *bootstrap.BootstrapConfig) *infrastructure.PaymentServer {
	wire.Build(
		db.NewDatabase,

		observe.NewTracer,

		infrabroker.InitializeRouter,
		infrabroker.NewNATSPublisher,
		infrabroker.NewNATSSubscriber,
		broker.NewSagaPaymentController,
		broker.NewPaymentEventRouter,

		repository.NewGormPaymentRepository,

		client.NewAuthConn,
		application.NewAuthService,
		application.NewPaymentService,
		application.NewSagaPaymentService,
		application.NewPaymentApplication,

		middleware.NewJwtAuthenticator,
		httppayment.NewGinEngine,
		httppayment.NewRouter,
		httppayment.New,

		infrastructure.NewPaymentServer,
	)

	return &infrastructure.PaymentServer{}
}

func InitializeOrchestratorServer(appCfg *config.ApplicationConfig, bootCfg *bootstrap.BootstrapConfig) *infrastructure.OrchestratorServer {
	wire.Build(
		observe.NewTracer,

		infrabroker.InitializeRouter,
		infrabroker.NewNATSPublisher,
		infrabroker.NewNATSSubscriber,
		infrabroker.NewRedisPublisher,
		application.NewOrchestratorService,
		broker.NewPurchaseResultPublisher,
		broker.NewSagaOrchestratorController,
		broker.NewOrchestratorEventRouter,

		infrastructure.NewOrchestratorServer,
	)

	return &infrastructure.OrchestratorServer{}
}
