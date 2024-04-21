package repository

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/valueobject"
)

// OrderRepository interface
type OrderRepository interface {
	GetOrder(ctx context.Context, orderID uint64) (*entity.Order, error)
	GetDetailedPurchasedItems(ctx context.Context, purchasedItems *[]valueobject.PurchasedItem) (*[]valueobject.DetailedPurchasedItem, error)
	CreateOrder(ctx context.Context, order *entity.Order) error
	DeleteOrder(ctx context.Context, orderID uint64) error
}
