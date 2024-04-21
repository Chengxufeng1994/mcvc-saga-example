package client

import (
	libconfig "github.com/Chengxufeng1994/go-saga-example/common/config"
	libgrpc "github.com/Chengxufeng1994/go-saga-example/common/grpc"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/config"
	"google.golang.org/grpc"
)

type AuthConn struct {
	conn *grpc.ClientConn
}

func (cli AuthConn) Conn() *grpc.ClientConn {
	return cli.conn
}

func NewAuthConn(libconfig *libconfig.ApplicationConfig) *AuthConn {
	config.ContextLogger.Infoln("starting connecting to grpc auth service...")
	conn := libgrpc.MustGRPCConn(libconfig.RpcEndpoints.AuthServiceHost)
	AuthClient := &AuthConn{conn: conn}
	return AuthClient
}
