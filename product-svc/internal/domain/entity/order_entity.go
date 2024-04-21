package entity

import "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/valueobject"

// Order entity
type Order struct {
	ID             uint64
	UserID         uint64
	PurchasedItems *[]valueobject.PurchasedItem
}
