package valueobject

// DetailedPurchasedItem value object
type DetailedPurchasedItem struct {
	ProductID   uint64
	Name        string
	Description string
	BrandName   string
	Price       int64
	Amount      int64
}
