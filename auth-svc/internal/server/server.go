package server

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/server/grpc"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/server/http"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/server/observe"
)

type Server struct {
	HttpSrv *http.HttpServer
	GrpcSrv *grpc.GrpcServer
}

func New(httpSrv *http.HttpServer, grpcSrv *grpc.GrpcServer) *Server {
	return &Server{
		HttpSrv: httpSrv,
		GrpcSrv: grpcSrv,
	}
}

func (srv *Server) Run() error {
	config.ContextLogger.Infoln("server.Run")
	go func() {
		if err := srv.HttpSrv.Run(); err != nil {
			config.ContextLogger.Fatal(err)
		}
	}()
	go func() {
		err := srv.GrpcSrv.Run()
		if err != nil {
			config.ContextLogger.Fatal(err)
		}
	}()
	return nil
}

func (srv *Server) GracefulShutdown(ctx context.Context) {
	config.ContextLogger.Infoln("server.GracefulShutdown")
	srv.HttpSrv.GracefulShutdown(ctx)
	srv.GrpcSrv.GracefulShutdown(ctx)
	if observe.TracerProvider != nil {
		err := observe.TracerProvider.Shutdown(ctx)
		if err != nil {
			config.ContextLogger.WithError(err).Error("server.GracefulShutdown trace provider shutdown")
		}
	}
}
