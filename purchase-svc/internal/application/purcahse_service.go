package application

import (
	"context"
	"errors"

	libmodel "github.com/Chengxufeng1994/go-saga-example/common/model"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/domain"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/repository"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/usecase"
	log "github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"
)

var (
	// ErrInvalidCartItemAmount is invalid cart item amount error
	ErrInvalidCartItemAmount = errors.New("invalid cart item amount")
	// ErrProductNotfound is product not found error
	ErrProductNotfound = errors.New("product not found")
	// ErrUnkownProductStatus unkown product status error
	ErrUnkownProductStatus = errors.New("unknown product status")
)

type PurchaseService struct {
	logger               *log.Entry
	productRepository    repository.ProductRepository
	purchasingRepository repository.PurchasingRepository
	sf                   *sonyflake.Sonyflake
}

func NewPurchaseService(
	logger *config.Logger,
	productRepository repository.ProductRepository,
	purchasingRepository repository.PurchasingRepository) usecase.PurchaseUseCase {

	var st sonyflake.Settings
	sf := sonyflake.NewSonyflake(st)

	return &PurchaseService{
		logger:               logger.ContextLogger.WithFields(log.Fields{"type": "service:PurchaseService"}),
		productRepository:    productRepository,
		purchasingRepository: purchasingRepository,
		sf:                   sf,
	}
}

// CheckProducts implements usecase.PurchaseUseCase.
func (svc *PurchaseService) CheckProducts(ctx context.Context, req *dto.CheckProductRequest) (*dto.CheckProductResponse, error) {
	for _, cartItem := range req.CartItems {
		if cartItem.Amount <= 0 {
			return nil, ErrInvalidCartItemAmount
		}
	}

	productStatuses, err := svc.productRepository.CheckProducts(ctx, svc.toCartItemDomain(req.CartItems))
	if err != nil {
		return nil, err
	}

	for _, productStatus := range productStatuses {
		switch productStatus.Status {
		case domain.ProductOk:
			continue
		case domain.ProductNotFound:
			return nil, ErrProductNotfound
		default:
			return nil, ErrUnkownProductStatus
		}
	}

	return &dto.CheckProductResponse{
		ProductStatues: svc.toProductStatusDto(productStatuses),
	}, nil
}

func (svc *PurchaseService) toCartItemDomain(dtos []*dto.CartItem) []*domain.CartItem {
	var domains []*domain.CartItem
	for _, dto := range dtos {
		domains = append(domains, &domain.CartItem{
			ProductID: dto.ProductID,
			Amount:    dto.Amount,
		})
	}

	return domains
}

func (svc *PurchaseService) toProductStatusDto(domains []*domain.ProductStatus) []*dto.ProductStatus {
	var dtos []*dto.ProductStatus
	for _, ps := range domains {
		dtos = append(dtos, &dto.ProductStatus{
			ProductID: ps.ProductID,
			Price:     ps.Price,
			Status:    ps.Status,
		})
	}

	return dtos
}

// CreatePurchase implements usecase.PurchaseUseCase.
func (svc *PurchaseService) CreatePurchase(ctx context.Context, userID uint64, req *dto.PurchaseCreationRequest) (*dto.PurchaseCreationResponse, error) {
	cartItems := req.CartItems
	resp, err := svc.CheckProducts(ctx, &dto.CheckProductRequest{
		CartItems: cartItems,
	})
	if err != nil {
		return nil, libmodel.NewAppError("CreatePurchase", "app.purchase.create.error", nil, "")
	}

	var amount int64 = 0
	for i, productStatus := range resp.ProductStatues {
		amount += cartItems[i].Amount * productStatus.Price
	}

	var cis []domain.CartItem
	for _, item := range cartItems {
		cis = append(cis, domain.CartItem{
			ProductID: item.ProductID,
			Amount:    item.Amount,
		})
	}

	purchaseID, err := svc.sf.NextID()
	if err != nil {
		return nil, libmodel.NewAppError("CreatePurchase", "app.purchase.gen_id.error", nil, "")
	}

	if err := svc.purchasingRepository.CreatePurchase(ctx, &domain.Purchase{
		ID: purchaseID,
		Order: &domain.Order{
			UserID:    userID,
			CartItems: &cis,
		},
		Payment: &domain.Payment{
			CurrencyCode: req.Payment.CurrencyCode,
			Amount:       amount,
		},
	}); err != nil {
		return nil, libmodel.NewAppError("CreatePurchase", "app.purchase.create.error", nil, "")
	}

	return &dto.PurchaseCreationResponse{
		PurchaseID: purchaseID,
	}, nil
}
