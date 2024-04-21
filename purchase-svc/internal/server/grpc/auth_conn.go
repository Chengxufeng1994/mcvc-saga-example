package grpc

import (
	libconfig "github.com/Chengxufeng1994/go-saga-example/common/config"
	libgrpc "github.com/Chengxufeng1994/go-saga-example/common/grpc"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/config"
	"google.golang.org/grpc"
)

// AuthClientConn grpc connection
var AuthClientConn *AuthConn

// AuthConn is a wrapper for Auth grpc connection
type AuthConn struct {
	conn *grpc.ClientConn
}

func NewAuthConn(logger *config.Logger, libconfig *libconfig.ApplicationConfig) *AuthConn {
	logger.ContextLogger.Infoln("starting connecting to grpc auth service...")
	conn := libgrpc.MustGRPCConn(libconfig.RpcEndpoints.AuthServiceHost)
	AuthClientConn = &AuthConn{conn: conn}
	return AuthClientConn
}

func (c AuthConn) Conn() *grpc.ClientConn {
	return c.conn
}
