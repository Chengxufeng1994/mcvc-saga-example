package dto

import "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/valueobject"

type Product struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BrandName   string `json:"brand_name"`
	Price       int64  `json:"price"`
	Inventory   int64  `json:"inventory"`
}

type ProductCreationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BrandName   string `json:"brand_name"`
	Price       int64  `json:"price"`
	Inventory   int64  `json:"inventory"`
}

type ProductCreationResponse struct {
	ID uint64 `json:"id"`
}

type ProductCheckRequest struct {
	CartItems []*valueobject.CartItem `json:"cart_items"`
}

type ProductCheckResponse struct {
	ProductStatus []*valueobject.ProductStatus `json:"product_status"`
}
