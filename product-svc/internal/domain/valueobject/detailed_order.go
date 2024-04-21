package valueobject

// DetailedOrder value object
type DetailedOrder struct {
	ID                     uint64
	UserID                 uint64
	DetailedPurchasedItems *[]DetailedPurchasedItem
}
