package valueobject

// CartItem value object
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
