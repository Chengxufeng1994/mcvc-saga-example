package repository

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/domain"
)

// PurchasingRepository is the repository interface of purchase aggregate
type PurchasingRepository interface {
	CreatePurchase(ctx context.Context, purchase *domain.Purchase) error
}
