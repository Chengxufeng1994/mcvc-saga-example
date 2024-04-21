package application

import (
	"context"
	"errors"

	"github.com/Chengxufeng1994/go-saga-example/common/model"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/repository"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/usecase"
	"github.com/sirupsen/logrus"
)

type PaymentService struct {
	logger            *logrus.Entry
	paymentRepository repository.PaymentRepository
}

func NewPaymentService(paymentRepository repository.PaymentRepository) usecase.PaymentUseCase {
	return &PaymentService{
		logger:            config.ContextLogger.WithFields(logrus.Fields{"type": "service:PaymentService"}),
		paymentRepository: paymentRepository,
	}
}

// GetPayment implements usecase.PaymentUseCase.
func (svc *PaymentService) GetPayment(ctx context.Context, userID uint64, paymentID uint64) (*dto.Payment, error) {
	payment, err := svc.paymentRepository.GetPayment(ctx, paymentID)
	if err != nil {
		var nfErr *repository.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetPayment", "app.payment.get_by_id.error", nil, "").Wrap(err)
		default:
			return nil, model.NewAppError("GetPayment", "app.payment.get_by_id.error", nil, "").Wrap(err)
		}
	}

	if userID != payment.UserID {
		return nil, model.NewAppError("GetPayment", "app.payment.user_id.error", nil, "").Wrap(err)
	}

	return &dto.Payment{
		ID:           paymentID,
		UserID:       userID,
		CurrencyCode: payment.CurrencyCode,
		Amount:       payment.Amount,
	}, nil
}

type SagaPaymentService struct {
	logger            *logrus.Entry
	paymentRepository repository.PaymentRepository
}

func NewSagaPaymentService(paymentRepository repository.PaymentRepository) usecase.SagaPaymentUseCase {
	return &SagaPaymentService{
		logger:            config.ContextLogger.WithFields(logrus.Fields{"type": "service:SagaPaymentService"}),
		paymentRepository: paymentRepository,
	}
}

// ExecuteCreatePayment implements usecase.SagaPaymentUseCase.
func (svc *SagaPaymentService) ExecuteCreatePayment(ctx context.Context, payment *entity.Payment) error {
	if err := svc.paymentRepository.CreatePayment(ctx, payment); err != nil {
		svc.logger.WithError(err).Error(err.Error())
		return model.NewAppError("ExecuteCreatePayment", "app.payment.create_payment.error", nil, "").Wrap(err)
	}

	return nil
}

// RollbackCreatePayment implements usecase.SagaPaymentUseCase.
func (svc *SagaPaymentService) RollbackCreatePayment(ctx context.Context, paymentID uint64) error {
	if err := svc.paymentRepository.DeletePayment(ctx, paymentID); err != nil {
		svc.logger.WithError(err).Error(err.Error())
		return model.NewAppError("RollbackCreatePayment", "app.payment.delete_payment.error", nil, "").Wrap(err)
	}

	return nil
}
