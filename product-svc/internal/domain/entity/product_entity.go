package entity

import "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/valueobject"

// Product entity
type Product struct {
	ID        uint64
	Detail    *valueobject.ProductDetail
	Inventory int64
}

// Idempotency entity
type Idempotency struct {
	ID        uint64
	ProductID uint64
	Amount    int64
}

// ProductStatus entity
type ProductStatus struct {
	ProductID uint64
	Price     int64
	Existed   bool
}
