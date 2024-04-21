package constant

type key string

const (
	CtxUserKey key = "ctx_user_key"
	CtxSpanKey key = "span_ctx_key"

	// HandlerHeader identifies a handler in the ReplyTopic
	HandlerHeader = "Handler"
	// UpdateProductInventoryHandler identifier
	UpdateProductInventoryHandler = "update_product_inventory_handler"
	// RollbackProductInventoryHandler identifier
	RollbackProductInventoryHandler = "rollback_product_inventory_handler"
	// CreateOrderHandler identifier
	CreateOrderHandler = "create_order_handler"
	// RollbackOrderHandler identifier
	RollbackOrderHandler = "rollback_order_handler"
	// CreatePaymentHandler identifier
	CreatePaymentHandler = "create_payment_handler"
	// RollbackPaymentHandler identifier
	RollbackPaymentHandler = "rollback_payment_handler"

	//
	JaegerHeader = "Uber-Trace-Id"
)
