package dto

// Purchase is the  request of creating new purchase
type PurchaseCreationRequest struct {
	Payment   *Payment    `json:"payment"`
	CartItems []*CartItem `json:"purchase_items" binding:"min=1"`
}

// PurchaseCreation response payload
type PurchaseCreationResponse struct {
	PurchaseID uint64 `json:"purchase_id"`
}

// Payment is the JSON request that represents a payment
type Payment struct {
	CurrencyCode string `json:"currency_code" binding:"required,oneof=NT US"`
}

// PurchaseResult is the HTTP JSON response of purchase result
type PurchaseResult struct {
	PurchaseID uint64 `json:"purchase_id"`
	Step       string `json:"step"`
	Status     string `json:"status"`
	Timestamp  int64  `json:"timestamp"`
}
