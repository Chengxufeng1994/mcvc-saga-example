package usecase

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/product-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/valueobject"
)

type ProductUseCase interface {
	// command
	CreateProduct(ctx context.Context, req *dto.ProductCreationRequest) (*dto.ProductCreationResponse, error)
	// query
	ListProducts(ctx context.Context) (*[]dto.Product, error)
	GetProduct(ctx context.Context, id uint64) (*dto.Product, error)
	GetProducts(ctx context.Context, ids []uint64) (*[]dto.Product, error)
	CheckProduct(ctx context.Context, req *dto.ProductCheckRequest) (*dto.ProductCheckResponse, error)
}

type SagaProductUseCase interface {
	UpdateProductInventory(ctx context.Context, idempotencyKey uint64, purchasedItems *[]valueobject.PurchasedItem) error
	RollbackProductInventory(ctx context.Context, idempotencyKey uint64) error
}
