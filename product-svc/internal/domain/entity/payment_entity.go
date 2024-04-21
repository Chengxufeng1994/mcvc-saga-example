package entity

// Payment entity
type Payment struct {
	ID           uint64
	UserID       uint64
	CurrencyCode string
	Amount       int64
}
