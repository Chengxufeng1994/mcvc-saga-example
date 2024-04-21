package grpc

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/common/pb"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/domain"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/repository"
	servergrpc "github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/server/grpc"
)

type GrpcProductRepository struct {
	productConn *servergrpc.ProductConn
}

func NewGrpcProductRepository(productConn *servergrpc.ProductConn) repository.ProductRepository {
	return &GrpcProductRepository{
		productConn: productConn,
	}
}

// CheckProducts implements repository.ProductRepository.
func (r *GrpcProductRepository) CheckProducts(ctx context.Context, cartItems []*domain.CartItem) ([]*domain.ProductStatus, error) {
	var pbCartItems []*pb.CartItem
	for _, cartItem := range cartItems {
		pbCartItems = append(pbCartItems, &pb.CartItem{
			ProductId: cartItem.ProductID,
			Amount:    cartItem.Amount,
		})
	}

	cli := pb.NewProductServiceClient(r.productConn.Conn())
	resp, err := cli.CheckProducts(ctx, &pb.CheckProductsRequest{
		CartItems: pbCartItems,
	})
	if err != nil {
		return nil, err
	}

	var productProductStates []*domain.ProductStatus
	for _, productStatus := range resp.ProductStatuses {
		productProductStates = append(productProductStates, &domain.ProductStatus{
			ProductID: productStatus.ProductId,
			Price:     productStatus.Price,
			Status:    getProductStatus(productStatus.Status),
		})

	}

	return productProductStates, nil
}

func getProductStatus(status pb.Status) domain.Status {
	switch status {
	case pb.Status_STATUS_OK:
		return domain.ProductOk
	case pb.Status_STATUS_NOT_FOUND:
		return domain.ProductNotFound
	}
	return -1
}
