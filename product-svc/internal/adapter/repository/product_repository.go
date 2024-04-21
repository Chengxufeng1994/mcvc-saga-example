package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	libmodel "github.com/Chengxufeng1994/go-saga-example/common/model"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/db/model"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/valueobject"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) repository.ProductRepository {
	return &GormProductRepository{
		db: db,
	}
}

// CreateProduct implements repository.ProductRepository.
func (g *GormProductRepository) CreateProduct(ctx context.Context, product *entity.Product) (uint64, error) {
	newRow := model.Product{
		Name:        product.Detail.Name,
		Description: product.Detail.Description,
		BrandName:   product.Detail.BrandName,
		Inventory:   product.Inventory,
		Price:       product.Detail.Price,
	}

	if err := g.db.WithContext(ctx).Clauses(clause.Returning{}).Model(&model.Product{}).Create(&newRow).Error; err != nil {
		return 0, err
	}

	return newRow.ID, nil
}

// ListProducts implements repository.ProductRepository.
func (g *GormProductRepository) ListProducts(ctx context.Context, offset int, size int) (*[]entity.Product, error) {
	var rows []model.Product
	if err := g.db.WithContext(ctx).Model(&model.Product{}).Find(&rows).Error; err != nil {
		return nil, err
	}

	n := len(rows)
	entities := make([]entity.Product, 0, n)
	for i := 0; i < n; i++ {
		row := rows[i]
		entities = append(entities, entity.Product{
			ID:        row.ID,
			Detail:    valueobject.NewProductDetail(row.Name, row.Description, row.BrandName, row.Price),
			Inventory: row.Inventory,
		})
	}

	return &entities, nil
}

// GetProduct implements repository.ProductRepository.
func (g *GormProductRepository) GetProduct(ctx context.Context, productID uint64) (*entity.Product, error) {
	var row model.Product
	if err := g.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", productID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewErrNotFound("product", strconv.Itoa(int(productID)))
		}
		return nil, err
	}

	return &entity.Product{
		ID:        row.ID,
		Detail:    valueobject.NewProductDetail(row.Name, row.Description, row.BrandName, row.Price),
		Inventory: row.Inventory,
	}, nil
}

// GetProductDetail implements repository.ProductRepository.
func (g *GormProductRepository) GetProductDetail(ctx context.Context, productID uint64) (*valueobject.ProductDetail, error) {
	var row valueobject.ProductDetail
	if err := g.db.WithContext(ctx).Model(&model.Product{}).
		Select("name", "description", "brand_name", "price").Where("id = ?", productID).First(&row).Error; err != nil {
		return nil, err
	}

	return &row, nil
}

// GetProductInventory implements repository.ProductRepository.
func (g *GormProductRepository) GetProductInventory(ctx context.Context, productID uint64) (int64, error) {
	var inventory int64
	if err := g.db.WithContext(ctx).Model(&model.Product{}).
		Select("inventory").Where("id = ?", productID).First(&inventory).Error; err != nil {
		return 0, err
	}

	return inventory, nil
}

// CheckProduct implements repository.ProductRepository.
func (g *GormProductRepository) CheckProduct(ctx context.Context, productID uint64) (*entity.ProductStatus, error) {
	var row model.Product
	if err := g.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", productID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entity.ProductStatus{
				ProductID: productID,
				Price:     0,
				Existed:   false,
			}, nil
		}
		return nil, err
	}

	return &entity.ProductStatus{
		ProductID: productID,
		Price:     row.Price,
		Existed:   true,
	}, nil
}

// UpdateProductInventory implements repository.ProductRepository.
func (g *GormProductRepository) UpdateProductInventory(ctx context.Context, idempotencyKey uint64, purchasedItems *[]valueobject.PurchasedItem) error {
	var err error
	var row model.Idempotency
	err = g.db.WithContext(ctx).Model(&model.Idempotency{}).Where("id = ?", idempotencyKey).First(&row).Error
	if err == nil {
		return repository.ErrInvalidIdempotency
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tx := g.db.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, purchasedItem := range *purchasedItems {
		var inventory int64
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&model.Product{}).Select("inventory").Where("id = ?", purchasedItem.ProductID).First(&inventory).Error; err != nil {
			tx.Rollback()
			return err
		}
		if inventory < purchasedItem.Amount {
			tx.Rollback()
			return repository.ErrInsuffientInventory
		}

		if err := tx.Model(&model.Product{}).Where("id = ?", purchasedItem.ProductID).Update("inventory", gorm.Expr("inventory - ?", purchasedItem.Amount)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	var idempotencies []model.Idempotency
	for _, purchasedItem := range *purchasedItems {
		idempotencies = append(idempotencies, model.Idempotency{
			BaseModel: libmodel.BaseModel{
				ID: idempotencyKey,
			},
			ProductID:  purchasedItem.ProductID,
			Amount:     purchasedItem.Amount,
			Rollbacked: false,
		})
	}

	if err := tx.Model(&model.Idempotency{}).Create(&idempotencies).WithContext(ctx).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// RollbackProductInventory implements repository.ProductRepository.
func (g *GormProductRepository) RollbackProductInventory(ctx context.Context, idempotencyKey uint64) (bool, *[]entity.Idempotency, error) {
	var idempotencies []model.Idempotency
	if err := g.db.WithContext(ctx).Model(&model.Idempotency{}).Select("product_id", "amount", "rollbacked").Where("id = ?", idempotencyKey).Order("product_id").Find(&idempotencies).Error; err != nil {
		return false, nil, err
	}
	if len(idempotencies) == 0 {
		return false, nil, fmt.Errorf("idempotency key not found: %v", idempotencyKey)
	}
	if idempotencies[0].Rollbacked {
		return true, nil, nil
	}

	tx := g.db.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, nil, err
	}

	for _, idempotency := range idempotencies {
		var inventory int64
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&model.Product{}).Select("inventory").Where("id = ?", idempotency.ProductID).First(&inventory).Error; err != nil {
			tx.Rollback()
			return false, nil, err
		}
		if err := tx.Model(&model.Product{}).Where("id = ?", idempotency.ProductID).Update("inventory", gorm.Expr("inventory + ?", idempotency.Amount)).Error; err != nil {
			tx.Rollback()
			return false, nil, err
		}
	}
	if err := tx.Model(&model.Idempotency{}).Where("id = ?", idempotencyKey).Update("rollbacked", true).WithContext(ctx).Error; err != nil {
		tx.Rollback()
		return false, nil, err
	}
	var domainIdempotencies []entity.Idempotency
	for _, idempotency := range idempotencies {
		domainIdempotencies = append(domainIdempotencies, entity.Idempotency{
			ID:        idempotencyKey,
			ProductID: idempotency.ProductID,
			Amount:    idempotency.Amount,
		})
	}
	return false, &domainIdempotencies, tx.Commit().Error
}
