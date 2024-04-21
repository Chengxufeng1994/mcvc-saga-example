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

type OrderService struct {
	logger          *logrus.Entry
	orderRepository repository.OrderRepository
}

func NewOrderService(orderRepository repository.OrderRepository) usecase.OrderUseCase {
	return &OrderService{
		logger:          config.ContextLogger.WithFields(logrus.Fields{"type": "service:OrderService"}),
		orderRepository: orderRepository,
	}
}

// GetDetailedOrder implements usecase.OrderUseCase.
func (svc *OrderService) GetDetailedOrder(ctx context.Context, userID uint64, orderID uint64) (*dto.GetDetailedOrderResponse, error) {
	order, err := svc.orderRepository.GetOrder(ctx, orderID)
	if err != nil {
		var nfErr *repository.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetOrder", "app.order.get_by_id.error", nil, "").Wrap(err)
		default:
			return nil, model.NewAppError("GetOrder", "app.order.get_by_id.error", nil, "").Wrap(err)
		}
	}

	if userID != order.UserID {
		return nil, model.NewAppError("GetOrder", "app.order.user_id.error", nil, "").Wrap(err)
	}

	detailedPurchasedItems, err := svc.orderRepository.GetDetailedPurchasedItems(ctx, order.PurchasedItems)
	if err != nil {
		svc.logger.Error(err.Error())
		return nil, model.NewAppError("GetOrder", "app.order.get_detailed_purchased_items.error", nil, "").Wrap(err)
	}

	var purchasedItems []dto.PurchasedItem
	for _, detailedPurchasedItem := range *detailedPurchasedItems {
		purchasedItems = append(purchasedItems, dto.PurchasedItem{
			ProductID:   detailedPurchasedItem.ProductID,
			Name:        detailedPurchasedItem.Name,
			Description: detailedPurchasedItem.Description,
			BrandName:   detailedPurchasedItem.BrandName,
			Price:       detailedPurchasedItem.Price,
			Amount:      detailedPurchasedItem.Amount,
		})
	}

	return &dto.GetDetailedOrderResponse{
		ID:             orderID,
		PurchasedItems: purchasedItems,
	}, nil
}

type SagaOrderService struct {
	logger          *logrus.Entry
	orderRepository repository.OrderRepository
}

func NewSagaOrderService(orderRepository repository.OrderRepository) usecase.SagaOrderUseCase {
	return &SagaOrderService{
		logger:          config.ContextLogger.WithFields(logrus.Fields{"type": "service:SagaOrderService"}),
		orderRepository: orderRepository,
	}
}

// ExecuteCreateOrder implements usecase.SagaOrderUseCase.
func (svc *SagaOrderService) ExecuteCreateOrder(ctx context.Context, order *entity.Order) error {
	if err := svc.orderRepository.CreateOrder(ctx, order); err != nil {
		svc.logger.WithError(err).Error(err.Error())
		return model.NewAppError("ExecuteCreateOrder", "app.order.create_order.error", nil, "").Wrap(err)
	}

	return nil
}

// RollbackCreateOrder implements usecase.SagaOrderUseCase.
func (svc *SagaOrderService) RollbackCreateOrder(ctx context.Context, orderID uint64) error {
	if err := svc.orderRepository.DeleteOrder(ctx, orderID); err != nil {
		svc.logger.WithError(err).Error(err.Error())
		return model.NewAppError("RollbackCreateOrder", "app.order.delete_order.error", nil, "").Wrap(err)
	}

	return nil
}
