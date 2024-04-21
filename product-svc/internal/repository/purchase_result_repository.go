package repository

import "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/event"

type PurchaseResultRepository interface {
	PublishPurchaseResult(correlationID string, evt *event.PurchaseResultEvent) error
}
