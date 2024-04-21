package entity

// CreatePurchaseResponse value object
type CreatePurchaseResponse struct {
	Purchase *Purchase
	Success  bool
	Error    string
}

// RollbackResponse value object
type RollbackResponse struct {
	UserID     uint64
	PurchaseID uint64
	Success    bool
	Error      string
}
