package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Chengxufeng1994/go-saga-example/common/constant"
	"github.com/Chengxufeng1994/go-saga-example/common/pb"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/valueobject"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	logger    = watermill.NewStdLogger(true, false)
	marshaler = &nats.GobMarshaler{}

	TraceContext        = propagation.TraceContext{}
	TraceparentHeader   = TraceContext.Fields()[0]
	W3CSupportedVersion = 0
)

// decode createPurchasedCommand to entity.Purchase
func DecodeCreatePurchaseCommand(payload message.Payload) (*entity.Purchase, *pb.Purchase, error) {
	var cmd pb.CreatePurchaseCommand
	if err := json.Unmarshal(payload, &cmd); err != nil {
		return nil, nil, err
	}

	var purchasedItems []valueobject.PurchasedItem
	for _, item := range cmd.Purchase.Order.PurchasedItems {
		purchasedItems = append(purchasedItems, valueobject.PurchasedItem{
			ProductID: item.ProductId,
			Amount:    item.Amount,
		})
	}

	purchaseID := cmd.PurchaseId
	ent := entity.Purchase{
		ID: purchaseID,
		Order: &entity.Order{
			ID:             purchaseID,
			UserID:         cmd.Purchase.Order.UserId,
			PurchasedItems: &purchasedItems,
		},
		Payment: &entity.Payment{
			ID:           purchaseID,
			UserID:       cmd.Purchase.Order.UserId,
			CurrencyCode: cmd.Purchase.Payment.CurrencyCode,
			Amount:       cmd.Purchase.Payment.Amount,
		},
	}

	return &ent, cmd.Purchase, nil
}

// SetSpanContext set span context to the message
func SetSpanContext(ctx context.Context, msg *message.Message) {
	msg.Metadata.Set(string(constant.CtxSpanKey), spanContextToW3C(ctx))
}

func spanContextToW3C(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return ""
	}
	// Clear all flags other than the trace-context supported sampling bit.
	flags := sc.TraceFlags() & trace.FlagsSampled
	return fmt.Sprintf("%.2x-%s-%s-%s",
		W3CSupportedVersion,
		sc.TraceID(),
		sc.SpanID(),
		flags)
}
