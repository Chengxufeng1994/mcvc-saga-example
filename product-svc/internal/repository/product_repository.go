package repository

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/valueobject"
)

// ProductRepository is the product repository interface
type ProductRepository interface {
	CheckProduct(ctx context.Context, productID uint64) (*entity.ProductStatus, error)
	CreateProduct(ctx context.Context, product *entity.Product) (uint64, error)
	ListProducts(ctx context.Context, offset, size int) (*[]entity.Product, error)
	GetProduct(ctx context.Context, productID uint64) (*entity.Product, error)
	GetProductDetail(ctx context.Context, productID uint64) (*valueobject.ProductDetail, error)
	GetProductInventory(ctx context.Context, productID uint64) (int64, error)
	// saga pattern
	UpdateProductInventory(ctx context.Context, idempotencyKey uint64, purchasedItems *[]valueobject.PurchasedItem) error
	RollbackProductInventory(ctx context.Context, idempotencyKey uint64) (bool, *[]entity.Idempotency, error)
}
