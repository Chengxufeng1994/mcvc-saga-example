package grpc

import (
	libconfig "github.com/Chengxufeng1994/go-saga-example/common/config"
	libgrpc "github.com/Chengxufeng1994/go-saga-example/common/grpc"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/config"
	"google.golang.org/grpc"
)

// ProductClientConn grpc connection
var ProductClientConn *ProductConn

// ProductConn is a wrapper for product grpc connection
type ProductConn struct {
	conn *grpc.ClientConn
}

func NewProductConn(logger *config.Logger, libconfig *libconfig.ApplicationConfig) *ProductConn {
	logger.ContextLogger.Infoln("starting connecting to product auth service...")
	conn := libgrpc.MustGRPCConn(libconfig.RpcEndpoints.ProductServiceHost)
	ProductClientConn = &ProductConn{conn: conn}
	return ProductClientConn
}

func (c ProductConn) Conn() *grpc.ClientConn {
	return c.conn
}
