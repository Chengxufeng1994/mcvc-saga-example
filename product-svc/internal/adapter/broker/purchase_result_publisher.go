package broker

import (
	"encoding/json"

	"github.com/Chengxufeng1994/go-saga-example/common/event"
	domainevent "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/event"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/broker"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/repository"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type PurchaseResultPublisher struct {
	topic     string
	publisher *broker.RedisPublisher
}

func NewPurchaseResultPublisher(publisher *broker.RedisPublisher) repository.PurchaseResultRepository {
	return &PurchaseResultPublisher{
		topic:     event.PurchaseResultTopic,
		publisher: publisher,
	}
}

// PublishPurchaseResult implements repository.PurchaseResultPublisher.
func (p *PurchaseResultPublisher) PublishPurchaseResult(correlationID string, evt *domainevent.PurchaseResultEvent) error {
	encoded := EncodeDomainPurchaseResult(evt)
	payload, err := json.Marshal(encoded)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	middleware.SetCorrelationID(correlationID, msg)

	if err := p.publisher.GetPublisher().Publish(p.topic, msg); err != nil {
		return err
	}

	return nil
}
