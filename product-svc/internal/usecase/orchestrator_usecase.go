package usecase

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
	"github.com/ThreeDotsLabs/watermill/message"
)

type OrchestratorUseCase interface {
	HandleTrx(ctx context.Context, purchase *entity.Purchase, correlationID string) error
	HandleReply(ctx context.Context, msg *message.Message, correlationID string) error
}
