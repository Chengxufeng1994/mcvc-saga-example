package domain

// Order entity
type Order struct {
	UserID    uint64
	CartItems *[]CartItem
}

// CartItem entity
type CartItem struct {
	ProductID uint64
	Amount    int64
}

func NewCartItem(productId uint64, amount int64) *CartItem {
	return &CartItem{
		ProductID: productId,
		Amount:    amount,
	}
}
