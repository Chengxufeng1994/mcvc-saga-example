package product

import (
	"context"
	"fmt"
	"net"

	"github.com/Chengxufeng1994/go-saga-example/common/bootstrap"
	"github.com/Chengxufeng1994/go-saga-example/common/pb"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/valueobject"
	infragrpc "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/grpc"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/usecase"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GrpcProductServer struct {
	application     string
	bootstrapConfig *bootstrap.BootstrapConfig
	productService  usecase.ProductUseCase
	srv             *grpc.Server
	pb.UnimplementedProductServiceServer
}

func NewGrpcProductServer(bootCfg *bootstrap.BootstrapConfig, productService usecase.ProductUseCase) *GrpcProductServer {
	grpcProductSrv := &GrpcProductServer{
		application:     bootCfg.Application,
		bootstrapConfig: bootCfg,
		productService:  productService,
	}

	grpcProductSrv.srv = infragrpc.InitializeServer(config.ContextLogger)
	pb.RegisterProductServiceServer(grpcProductSrv.srv, grpcProductSrv)

	grpc_prometheus.Register(grpcProductSrv.srv)
	reflection.Register(grpcProductSrv.srv)
	return grpcProductSrv
}

// CheckProducts implements pb.ProductServiceServer.
func (s *GrpcProductServer) CheckProducts(ctx context.Context, req *pb.CheckProductsRequest) (*pb.CheckProductsResponse, error) {
	var cartItems []*valueobject.CartItem
	pbCartItems := req.CartItems
	for _, pbCartItem := range pbCartItems {
		cartItems = append(cartItems, &valueobject.CartItem{
			ProductID: pbCartItem.ProductId,
			Amount:    pbCartItem.Amount,
		})
	}

	resp, err := s.productService.CheckProduct(ctx, &dto.ProductCheckRequest{CartItems: cartItems})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("internal error: %v", err),
		)
	}

	productStatues := resp.ProductStatus
	pbStatues := make([]*pb.ProductStatus, 0, len(productStatues))
	for _, status := range productStatues {
		pbStatues = append(pbStatues, &pb.ProductStatus{
			ProductId: status.ProductID,
			Price:     status.Price,
			Status:    getPbProductStatus(status.Status),
		})
	}

	return &pb.CheckProductsResponse{
		ProductStatuses: pbStatues,
	}, nil
}

// GetProducts implements pb.ProductServiceServer.
func (s *GrpcProductServer) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	productIDs := req.ProductIds
	result, err := s.productService.GetProducts(ctx, productIDs)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("internal error: %v", err),
		)
	}
	var products []*pb.Product
	for _, product := range *result {
		products = append(products, &pb.Product{
			ProductId:   product.ID,
			ProductName: product.Name,
			Description: product.Description,
			BrandName:   product.BrandName,
			Inventory:   product.Inventory,
			Price:       product.Price,
		})
	}
	return &pb.GetProductsResponse{
		Products: products,
	}, nil
}

func (s *GrpcProductServer) Run() error {
	addr := fmt.Sprintf("%s:%d", s.bootstrapConfig.Grpc.Host, s.bootstrapConfig.Grpc.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	config.ContextLogger.Infoln("grpc.Run listening on", s.bootstrapConfig.Grpc.Port)
	if err := s.srv.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *GrpcProductServer) GracefulShutdown(ctx context.Context) {
	config.ContextLogger.Infoln("grpc.GracefulShutdown")
	s.srv.GracefulStop()
}

func getPbProductStatus(status valueobject.Status) pb.Status {
	switch status {
	case valueobject.ProductOk:
		return pb.Status_STATUS_OK
	case valueobject.ProductNotFound:
		return pb.Status_STATUS_NOT_FOUND
	}
	return pb.Status_STATUS_NOT_FOUND
}
