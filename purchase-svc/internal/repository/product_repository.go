package repository

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/domain"
)

// ProductRepository is the product repository interface
type ProductRepository interface {
	CheckProducts(ctx context.Context, cartItems []*domain.CartItem) ([]*domain.ProductStatus, error)
}
