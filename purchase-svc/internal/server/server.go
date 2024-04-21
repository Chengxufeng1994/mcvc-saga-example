package server

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/server/grpc"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/server/http"
	log "github.com/sirupsen/logrus"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Server struct {
	HttpSrv        *http.HttpServer
	TracerProvider *sdktrace.TracerProvider
}

func New(httpSrv *http.HttpServer, tracerProvider *sdktrace.TracerProvider) *Server {
	return &Server{
		HttpSrv: httpSrv,
	}
}

func (srv *Server) Run() error {
	log.Infoln("server.Run")
	go func() {
		if err := srv.HttpSrv.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}

func (srv *Server) GracefulShutdown(ctx context.Context) {
	log.Infoln("server.GracefulShutdown")
	srv.HttpSrv.GracefulShutdown(ctx)

	if srv.TracerProvider != nil {
		err := srv.TracerProvider.Shutdown(ctx)
		if err != nil {
			log.WithError(err).Error("server.GracefulShutdown trace provider shutdown")
		}
	}

	if grpc.AuthClientConn != nil {
		if err := grpc.AuthClientConn.Conn().Close(); err != nil {
			log.WithError(err).Error("server.GracefulShutdown grpc auth client close")
		}
	}
}
