//go:build wireinject
// +build wireinject

package di

import (
	"github.com/Chengxufeng1994/go-saga-example/common/bootstrap"
	libconfig "github.com/Chengxufeng1994/go-saga-example/common/config"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/adapter/broker"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/adapter/grpc"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/application"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/server"
	infrabroker "github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/server/broker"
	srvgrpc "github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/server/grpc"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/server/http"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/server/observe"
	"github.com/google/wire"
)

func InitApplicationConfig(path string) *libconfig.ApplicationConfig {
	wire.Build(
		libconfig.LoadApplicationConfig,
	)
	return &libconfig.ApplicationConfig{}
}

func InitBootstrapConfig(path string) *bootstrap.BootstrapConfig {
	wire.Build(
		bootstrap.LoadBootstrapConfig,
	)
	return &bootstrap.BootstrapConfig{}
}

func InitializeServer(appCfg *libconfig.ApplicationConfig, bootCfg *bootstrap.BootstrapConfig) *server.Server {
	wire.Build(
		config.InitLogger,
		// init tracer
		observe.NewTracer,
		//init redis
		// redis.NewRedisClusterClient,
		// init broker
		// broker.NewRedisSubscriber,
		infrabroker.NewNatsPublisher,

		// grpc client
		srvgrpc.NewProductConn,
		srvgrpc.NewAuthConn,

		// internal services
		grpc.NewGrpcAuthRepository,
		grpc.NewGrpcProductRepository,
		broker.NewNatsNatsPurchasePublisher,
		application.NewAuthService,
		application.NewPurchaseService,
		application.New,

		// init server
		http.NewGinEngine,
		http.NewRouter,
		http.New,
		server.New,
	)

	return &server.Server{}
}
