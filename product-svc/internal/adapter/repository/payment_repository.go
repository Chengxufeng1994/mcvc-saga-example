package repository

import (
	"context"
	"errors"
	"strconv"

	libcommon "github.com/Chengxufeng1994/go-saga-example/common/model"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/db/model"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/repository"
	"gorm.io/gorm"
)

// PaymentRepositoryImpl implementation
type GormPaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository factory
func NewGormPaymentRepository(db *gorm.DB) repository.PaymentRepository {
	return &GormPaymentRepository{
		db: db,
	}
}

// GetPayment get an payment
func (repo *GormPaymentRepository) GetPayment(ctx context.Context, paymentID uint64) (*entity.Payment, error) {
	var payment model.Payment
	if err := repo.db.Model(&model.Payment{}).Select("id", "customer_id", "currency_code", "amount").Where("id = ?", paymentID).First(&payment).WithContext(ctx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewErrNotFound("payment", strconv.Itoa(int(paymentID)))
		}
		return nil, err
	}
	return &entity.Payment{
		ID:           payment.ID,
		UserID:       payment.UserID,
		CurrencyCode: payment.CurrencyCode,
		Amount:       payment.Amount,
	}, nil
}

// CreatePayment creates a payment
func (repo *GormPaymentRepository) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	if err := repo.db.WithContext(ctx).Create(&model.Payment{
		BaseModel: libcommon.BaseModel{
			ID: payment.ID,
		},
		UserID:       payment.UserID,
		CurrencyCode: payment.CurrencyCode,
		Amount:       payment.Amount,
	}).Error; err != nil {
		return err
	}
	return nil
}

// DeletePayment deletes an payment
func (repo *GormPaymentRepository) DeletePayment(ctx context.Context, paymentID uint64) error {
	if err := repo.db.Exec("DELETE FROM payments where id = ?", paymentID).Error; err != nil {
		return err
	}
	return nil
}
