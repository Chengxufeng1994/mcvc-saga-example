package usecase

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/dto"
)

// PurchasingUseCase is the interface of purchasing service
type PurchaseUseCase interface {
	CheckProducts(ctx context.Context, req *dto.CheckProductRequest) (*dto.CheckProductResponse, error)
	CreatePurchase(ctx context.Context, userID uint64, req *dto.PurchaseCreationRequest) (*dto.PurchaseCreationResponse, error)
}
