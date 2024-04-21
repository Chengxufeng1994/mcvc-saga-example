package usecase

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/product-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
)

// PaymentService interface
type PaymentUseCase interface {
	// command
	// query
	GetPayment(ctx context.Context, userID, paymentID uint64) (*dto.Payment, error)
}

// SagaPaymentService interface
type SagaPaymentUseCase interface {
	ExecuteCreatePayment(ctx context.Context, payment *entity.Payment) error
	RollbackCreatePayment(ctx context.Context, paymentID uint64) error
}
