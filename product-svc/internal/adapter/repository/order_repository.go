package repository

import (
	"context"
	"errors"
	"strconv"

	libcommon "github.com/Chengxufeng1994/go-saga-example/common/model"
	"github.com/Chengxufeng1994/go-saga-example/common/pb"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/db/model"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/valueobject"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/grpc/client"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/repository"
	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db            *gorm.DB
	productClient *client.ProductConn
}

func NewGormOrderRepository(db *gorm.DB, productClient *client.ProductConn) repository.OrderRepository {
	return &GormOrderRepository{
		db:            db,
		productClient: productClient,
	}
}

// CreateOrder implements repository.OrderRepository.
func (g *GormOrderRepository) CreateOrder(ctx context.Context, order *entity.Order) error {
	id := order.ID
	userID := order.UserID
	var entries []model.Order
	for _, purchasedItem := range *order.PurchasedItems {
		entries = append(entries, model.Order{
			BaseModel: libcommon.BaseModel{
				ID: id,
			},
			ProductID: purchasedItem.ProductID,
			Amount:    purchasedItem.Amount,
			UserID:    userID,
		})
	}

	if err := g.db.WithContext(ctx).Model(&model.Order{}).Create(&entries).Error; err != nil {
		return err
	}

	return nil
}

// GetOrder implements repository.OrderRepository.
func (g *GormOrderRepository) GetOrder(ctx context.Context, orderID uint64) (*entity.Order, error) {
	var orders []model.Order
	if err := g.db.WithContext(ctx).Model(&model.Order{}).Select("id", "product_id", "amount", "user_id").Where("id = ?", orderID).Order("product_id").Find(&orders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewErrNotFound("order", strconv.Itoa(int(orderID)))
		}
		return nil, err
	}
	if len(orders) == 0 {
		return nil, repository.NewErrNotFound("order", strconv.Itoa(int(orderID)))
	}

	var purchasedItems []valueobject.PurchasedItem
	for _, order := range orders {
		purchasedItems = append(purchasedItems, valueobject.PurchasedItem{
			ProductID: order.ProductID,
			Amount:    order.Amount,
		})
	}

	return &entity.Order{
		ID:             orders[0].ID,
		UserID:         orders[0].UserID,
		PurchasedItems: &purchasedItems,
	}, nil
}

// GetDetailedPurchasedItems implements repository.OrderRepository.
func (g *GormOrderRepository) GetDetailedPurchasedItems(ctx context.Context, purchasedItems *[]valueobject.PurchasedItem) (*[]valueobject.DetailedPurchasedItem, error) {
	var productIDs []uint64
	for _, purchasedItem := range *purchasedItems {
		productIDs = append(productIDs, purchasedItem.ProductID)
	}

	cli := pb.NewProductServiceClient(g.productClient.Conn())
	res, err := cli.GetProducts(ctx, &pb.GetProductsRequest{
		ProductIds: productIDs,
	})
	if err != nil {
		return nil, err
	}

	var detailedPurchasedItems []valueobject.DetailedPurchasedItem
	for i, prod := range res.Products {
		detailedPurchasedItems = append(detailedPurchasedItems, valueobject.DetailedPurchasedItem{
			ProductID:   prod.ProductId,
			Name:        prod.ProductName,
			Description: prod.Description,
			BrandName:   prod.BrandName,
			Price:       prod.Price,
			Amount:      (*purchasedItems)[i].Amount,
		})
	}

	return &detailedPurchasedItems, nil
}

// DeleteOrder implements repository.OrderRepository.
func (g *GormOrderRepository) DeleteOrder(ctx context.Context, orderID uint64) error {
	if err := g.db.Exec("DELETE FROM orders where id = ?", orderID).Error; err != nil {
		return err
	}
	return nil
}
