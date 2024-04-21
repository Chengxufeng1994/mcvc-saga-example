package dto

import "github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/domain"

// CartItem is the JSON request that represents an order
type CartItem struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	Amount    int64  `json:"amount" binding:"required,number,min=1"`
}

type CheckProductRequest struct {
	CartItems []*CartItem `json:"cart_items"`
}

// CartItem is the JSON request that represents an order
type ProductStatus struct {
	ProductID uint64        `json:"product_id"`
	Price     int64         `json:"price""`
	Status    domain.Status `json:"status"`
}

type CheckProductResponse struct {
	ProductStatues []*ProductStatus `json:"product_statues"`
}
