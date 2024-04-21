package event

const (
	// PurchaseTopic is the topic to which we publish new purchase
	PurchaseTopic = "purchase"
	// PurchaseResultTopic is the subscribed topic for purchase result
	PurchaseResultTopic = "purchase.result"
	// ReplyTopic is saga step reply topic
	ReplyTopic = "reply"
	// UpdateProductInventoryTopic topic
	UpdateProductInventoryTopic = "product_update_inventory"
	// RollbackProductInventoryTopic topic
	RollbackProductInventoryTopic = "product_rollback_inventory"
	// Create Order Topic
	CreateOrderTopic = "order_create"
	// Rollback Order Topic
	RollbackOrderTopic = "order_rollback"
	// Payment Order Topic
	CreatePaymentTopic = "payment_create"
	// Rollback Order Topic
	RollbackPaymentTopic = "payment_rollback"
)
