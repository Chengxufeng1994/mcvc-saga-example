package application

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/common/model"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/valueobject"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/repository"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/usecase"
	"github.com/sirupsen/logrus"
)

type ProductService struct {
	logger            *logrus.Entry
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) usecase.ProductUseCase {
	return &ProductService{
		logger: config.ContextLogger.WithFields(logrus.Fields{
			"type": "service:ProductService",
		}),
		productRepository: productRepository,
	}
}

// CreateProduct implements usecase.ProductUseCase.
func (p *ProductService) CreateProduct(ctx context.Context, req *dto.ProductCreationRequest) (*dto.ProductCreationResponse, error) {
	entity := &entity.Product{
		Detail:    valueobject.NewProductDetail(req.Name, req.Description, req.BrandName, req.Price),
		Inventory: req.Inventory,
	}

	id, err := p.productRepository.CreateProduct(ctx, entity)
	if err != nil {
		return nil, model.NewAppError("CreateProduct", "app.product.create.error", nil, "").Wrap(err)
	}

	return &dto.ProductCreationResponse{
		ID: id,
	}, nil
}

// ListProducts implements usecase.ProductUseCase.
func (p *ProductService) ListProducts(ctx context.Context) (*[]dto.Product, error) {
	entities, err := p.productRepository.ListProducts(ctx, 0, 0)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.Product, 0, len(*entities))
	for _, entity := range *entities {
		dtos = append(dtos, dto.Product{
			ID:          entity.ID,
			Name:        entity.Detail.Name,
			Description: entity.Detail.Description,
			BrandName:   entity.Detail.BrandName,
			Price:       entity.Detail.Price,
			Inventory:   entity.Inventory,
		})
	}

	return &dtos, nil
}

// GetProduct implements usecase.ProductUseCase.
func (p *ProductService) GetProduct(ctx context.Context, id uint64) (*dto.Product, error) {
	entity, err := p.productRepository.GetProduct(ctx, id)
	if err != nil {
		return nil, model.NewAppError("GetPorduct", "app.product.not_found.error", nil, "").Wrap(err)
	}

	return &dto.Product{
		ID:          entity.ID,
		Name:        entity.Detail.Name,
		Description: entity.Detail.Description,
		BrandName:   entity.Detail.BrandName,
		Price:       entity.Detail.Price,
		Inventory:   entity.Inventory,
	}, nil
}

// GetProducts implements usecase.ProductUseCase.
func (p *ProductService) GetProducts(ctx context.Context, ids []uint64) (*[]dto.Product, error) {
	var products []dto.Product
	for _, id := range ids {
		productDetail, err := p.productRepository.GetProductDetail(ctx, id)
		if err != nil {
			p.logger.WithError(err).Error("GetProducts")
			return nil, model.NewAppError("GetProducts", "app.product.not_found.error", nil, "").Wrap(err)
		}

		inventory, err := p.productRepository.GetProductInventory(ctx, id)
		if err != nil {
			p.logger.WithError(err).Error("GetProducts")
			return nil, model.NewAppError("GetProducts", "app.product.not_found.error", nil, "").Wrap(err)
		}

		products = append(products, dto.Product{
			ID:          id,
			Name:        productDetail.Name,
			Description: productDetail.Description,
			BrandName:   productDetail.BrandName,
			Price:       productDetail.Price,
			Inventory:   inventory,
		})
	}

	return &products, nil
}

// CheckProduct implements usecase.ProductUseCase.
func (p *ProductService) CheckProduct(ctx context.Context, req *dto.ProductCheckRequest) (*dto.ProductCheckResponse, error) {
	cartItems := req.CartItems
	productStatues := make([]*valueobject.ProductStatus, 0)
	for _, ci := range cartItems {
		entity, err := p.productRepository.CheckProduct(ctx, ci.ProductID)
		if err != nil {
			return nil, err
		}

		productStatues = append(productStatues, valueobject.NewProductStatus(entity.ProductID, entity.Price, entity.Existed))
	}

	return &dto.ProductCheckResponse{
		ProductStatus: productStatues,
	}, nil
}

type SagaProductService struct {
	logger            *logrus.Entry
	productRepository repository.ProductRepository
}

func NewSagaProductService(productRepository repository.ProductRepository) usecase.SagaProductUseCase {
	return &SagaProductService{
		logger: config.ContextLogger.WithFields(logrus.Fields{
			"type": "service:SagaProductService",
		}),
		productRepository: productRepository,
	}
}

// UpdateProductInventory implements usecase.SagaProductUseCase.
func (svc *SagaProductService) UpdateProductInventory(ctx context.Context, idempotencyKey uint64, purchasedItems *[]valueobject.PurchasedItem) error {
	err := svc.productRepository.UpdateProductInventory(ctx, idempotencyKey, purchasedItems)
	if err != nil {
		svc.logger.WithError(err).Error(err.Error())
		switch err {
		case repository.ErrInsuffientInventory:
			return model.NewAppError("UpdateProductInventory", "app.product.insuffient_inventory.error", nil, "insufficient inventory")
		case repository.ErrInvalidIdempotency:
			return model.NewAppError("UpdateProductInventory", "app.product.insuffient_inventory.error", nil, "invalid dempotency")
		default:
			return model.NewAppError("UpdateProductInventory", "app.product.update_product_inventory.error", nil, "unknown error")
		}
	}
	return nil
}

// RollbackProductInventory implements usecase.SagaProductUseCase.
func (svc *SagaProductService) RollbackProductInventory(ctx context.Context, idempotencyKey uint64) error {
	_, _, err := svc.productRepository.RollbackProductInventory(ctx, idempotencyKey)
	if err != nil {
		svc.logger.WithError(err).Error(err.Error())
		return model.NewAppError("RollbackProductInventory", "app.product.rollback_product_inventory.error", nil, "")
	}
	return nil

}
