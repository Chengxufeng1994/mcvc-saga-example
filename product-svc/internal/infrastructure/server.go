package infrastructure

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/product-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/broker"
	infragrpcproduct "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/grpc/product"
	httporder "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/order"
	httppayment "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/payment"
	httpproduct "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/product"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type ProductServer struct {
	HttpSrv        *httpproduct.HttpServer
	GrpcSrv        *infragrpcproduct.GrpcProductServer
	EventRouter    broker.EventRouter
	TracerProvider *sdktrace.TracerProvider
}

type OrderServer struct {
	HttpSrv        *httporder.HttpServer
	EventRouter    broker.EventRouter
	TracerProvider *sdktrace.TracerProvider
}

type PaymentServer struct {
	HttpSrv        *httppayment.HttpServer
	EventRouter    broker.EventRouter
	TracerProvider *sdktrace.TracerProvider
}

// OrchestratorServer wrapper
type OrchestratorServer struct {
	EventRouter    broker.EventRouter
	TracerProvider *sdktrace.TracerProvider
}

func NewProductServer(
	httpSrv *httpproduct.HttpServer,
	grpcSrv *infragrpcproduct.GrpcProductServer,
	eventRouter broker.EventRouter,
	tracerProvider *sdktrace.TracerProvider) *ProductServer {
	return &ProductServer{
		HttpSrv:     httpSrv,
		GrpcSrv:     grpcSrv,
		EventRouter: eventRouter,
	}
}

func (srv *ProductServer) Run() error {
	config.ContextLogger.Infoln("server.Run")
	go func() {
		if err := srv.HttpSrv.Run(); err != nil {
			config.ContextLogger.Fatal(err)
		}
	}()

	go func() {
		if err := srv.GrpcSrv.Run(); err != nil {
			config.ContextLogger.Fatal(err)
		}
	}()

	go func() {
		err := srv.EventRouter.Run()
		if err != nil {
			config.ContextLogger.Fatal(err)
		}
	}()
	return nil
}

func (srv *ProductServer) GracefulShutdown(ctx context.Context) {
	config.ContextLogger.Infoln("server.GracefulShutdown")
	srv.HttpSrv.GracefulShutdown(ctx)
	srv.GrpcSrv.GracefulShutdown(ctx)

	if err := srv.EventRouter.GracefulShutdown(); err != nil {
		config.ContextLogger.WithError(err).Error("server.GracefulShutdown event router shutdown")
	}

	if srv.TracerProvider != nil {
		err := srv.TracerProvider.Shutdown(ctx)
		if err != nil {
			config.ContextLogger.WithError(err).Error("server.GracefulShutdown trace provider shutdown")
		}
	}
}

func NewOrderServer(
	httpSrv *httporder.HttpServer,
	eventRouter broker.EventRouter,
	tracerProvider *sdktrace.TracerProvider) *OrderServer {
	return &OrderServer{
		HttpSrv:     httpSrv,
		EventRouter: eventRouter,
	}
}

func (srv *OrderServer) Run() error {
	config.ContextLogger.Infoln("server.Run")
	go func() {
		if err := srv.HttpSrv.Run(); err != nil {
			config.ContextLogger.Fatal(err)
		}
	}()

	go func() {
		err := srv.EventRouter.Run()
		if err != nil {
			config.ContextLogger.Fatal(err)
		}
	}()

	return nil
}

func (srv *OrderServer) GracefulShutdown(ctx context.Context) {
	config.ContextLogger.Infoln("server.GracefulShutdown")
	srv.HttpSrv.GracefulShutdown(ctx)

	if err := srv.EventRouter.GracefulShutdown(); err != nil {
		config.ContextLogger.WithError(err).Error("server.GracefulShutdown event router shutdown")
	}

	if srv.TracerProvider != nil {
		err := srv.TracerProvider.Shutdown(ctx)
		if err != nil {
			config.ContextLogger.WithError(err).Error("server.GracefulShutdown trace provider shutdown")
		}
	}
}

func NewPaymentServer(
	httpSrv *httppayment.HttpServer,
	eventRouter broker.EventRouter,
	tracerProvider *sdktrace.TracerProvider) *PaymentServer {
	return &PaymentServer{
		HttpSrv:     httpSrv,
		EventRouter: eventRouter,
	}
}

func (srv *PaymentServer) Run() error {
	config.ContextLogger.Infoln("server.Run")
	go func() {
		if err := srv.HttpSrv.Run(); err != nil {
			config.ContextLogger.Fatal(err)
		}
	}()

	go func() {
		err := srv.EventRouter.Run()
		if err != nil {
			config.ContextLogger.Fatal(err)
		}
	}()

	return nil
}

func (srv *PaymentServer) GracefulShutdown(ctx context.Context) {
	config.ContextLogger.Infoln("server.GracefulShutdown")
	srv.HttpSrv.GracefulShutdown(ctx)

	if err := srv.EventRouter.GracefulShutdown(); err != nil {
		config.ContextLogger.WithError(err).Error("server.GracefulShutdown event router shutdown")
	}

	if srv.TracerProvider != nil {
		err := srv.TracerProvider.Shutdown(ctx)
		if err != nil {
			config.ContextLogger.WithError(err).Error("server.GracefulShutdown trace provider shutdown")
		}
	}
}

func NewOrchestratorServer(
	eventRouter broker.EventRouter,
	tracerProvider *sdktrace.TracerProvider) *OrchestratorServer {
	return &OrchestratorServer{
		EventRouter: eventRouter,
	}
}

func (srv *OrchestratorServer) Run() error {
	config.ContextLogger.Infoln("server.Run")

	go func() {
		err := srv.EventRouter.Run()
		if err != nil {
			config.ContextLogger.Fatal(err)
		}
	}()

	return nil
}

func (srv *OrchestratorServer) GracefulShutdown(ctx context.Context) {
	config.ContextLogger.Infoln("server.GracefulShutdown")

	if err := srv.EventRouter.GracefulShutdown(); err != nil {
		config.ContextLogger.WithError(err).Error("server.GracefulShutdown event router shutdown")
	}

	if srv.TracerProvider != nil {
		err := srv.TracerProvider.Shutdown(ctx)
		if err != nil {
			config.ContextLogger.WithError(err).Error("server.GracefulShutdown trace provider shutdown")
		}
	}
}
