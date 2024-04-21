package broker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/common/event"
	"github.com/Chengxufeng1994/go-saga-example/common/pb"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/domain"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/repository"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NatsPurchasePublisher struct {
	publisher message.Publisher
}

func NewNatsNatsPurchasePublisher(publisher message.Publisher) repository.PurchasingRepository {
	return &NatsPurchasePublisher{
		publisher: publisher,
	}
}

// CreatePurchase implements repository.PurchasingRepository.
func (n *NatsPurchasePublisher) CreatePurchase(ctx context.Context, purchase *domain.Purchase) error {
	order := *purchase.Order
	var purchasedItems []*pb.PurchasedItem
	for _, item := range *order.CartItems {
		purchasedItems = append(purchasedItems, &pb.PurchasedItem{
			ProductId: item.ProductID,
			Amount:    item.Amount,
		})
	}
	payment := *purchase.Payment

	createPurchaseDto := &pb.CreatePurchaseCommand{
		PurchaseId: purchase.ID,
		Purchase: &pb.Purchase{
			Order: &pb.Order{
				UserId:         order.UserID,
				PurchasedItems: purchasedItems,
			},
			Payment: &pb.Payment{
				CurrencyCode: payment.CurrencyCode,
				Amount:       payment.Amount,
			},
		},
		Timestamp: timestamppb.New(time.Now()),
	}
	payload, err := json.Marshal(createPurchaseDto)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	middleware.SetCorrelationID(watermill.NewUUID(), msg)

	if err := n.publisher.Publish(event.PurchaseTopic, msg); err != nil {
		return err
	}

	return nil
}
