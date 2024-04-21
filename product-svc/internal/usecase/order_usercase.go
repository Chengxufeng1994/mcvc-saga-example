package usecase

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/product-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
)

// OrderUseCase interface
type OrderUseCase interface {
	GetDetailedOrder(ctx context.Context, userID, orderID uint64) (*dto.GetDetailedOrderResponse, error)
}

// SagaOrderUseCase interface
type SagaOrderUseCase interface {
	ExecuteCreateOrder(ctx context.Context, order *entity.Order) error
	RollbackCreateOrder(ctx context.Context, orderID uint64) error
}
