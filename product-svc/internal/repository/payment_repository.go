package repository

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
)

type PaymentRepository interface {
	GetPayment(ctx context.Context, paymentID uint64) (*entity.Payment, error)
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	DeletePayment(ctx context.Context, paymentID uint64) error
}
