package client

import (
	libconfig "github.com/Chengxufeng1994/go-saga-example/common/config"
	libgrpc "github.com/Chengxufeng1994/go-saga-example/common/grpc"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/config"
	"google.golang.org/grpc"
)

type ProductConn struct {
	conn *grpc.ClientConn
}

func (cli ProductConn) Conn() *grpc.ClientConn {
	return cli.conn
}

func NewProductConn(libconfig *libconfig.ApplicationConfig) *ProductConn {
	config.ContextLogger.Infoln("starting connecting to grpc product service...")
	conn := libgrpc.MustGRPCConn(libconfig.RpcEndpoints.ProductServiceHost)
	ProductClient := &ProductConn{conn: conn}
	return ProductClient
}
