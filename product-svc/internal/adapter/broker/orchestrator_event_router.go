package broker

import (
	"context"
	"log"

	"github.com/Chengxufeng1994/go-saga-example/common/constant"
	"github.com/Chengxufeng1994/go-saga-example/common/event"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/broker"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/usecase"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"go.opentelemetry.io/otel/propagation"
)

type sagaOrchestratorController struct {
	svc usecase.OrchestratorUseCase
}

func NewSagaOrchestratorController(svc usecase.OrchestratorUseCase) *sagaOrchestratorController {
	return &sagaOrchestratorController{
		svc: svc,
	}
}

// handleTransaction start the purchase transaction
func (ctrl sagaOrchestratorController) HandleTrx(msg *message.Message) error {
	log.Println("handleTransaction received message", msg.UUID)
	req, _, err := broker.DecodeCreatePurchaseCommand(msg.Payload)
	if err != nil {
		return err
	}
	correlationID := msg.Metadata.Get(middleware.CorrelationIDMetadataKey)
	carrier := make(propagation.HeaderCarrier)
	carrier.Set(broker.TraceparentHeader, msg.Metadata.Get(string(constant.CtxSpanKey)))
	parentCtx := broker.TraceContext.Extract(context.Background(), carrier)
	return ctrl.svc.HandleTrx(parentCtx, req, correlationID)
}

func (ctrl sagaOrchestratorController) HandleReply(msg *message.Message) error {
	correlationID := msg.Metadata.Get(middleware.CorrelationIDMetadataKey)
	carrier := make(propagation.HeaderCarrier)
	carrier.Set(broker.TraceparentHeader, msg.Metadata.Get(string(constant.CtxSpanKey)))
	parentCtx := broker.TraceContext.Extract(context.Background(), carrier)
	return ctrl.svc.HandleReply(parentCtx, msg, correlationID)
}

type OrchestratorEventRouter struct {
	router     *message.Router
	publisher  broker.NatsPublisher
	subscriber broker.NatsSubscriber
	controller *sagaOrchestratorController
}

func NewOrchestratorEventRouter(
	router *message.Router,
	publisher broker.NatsPublisher,
	subscriber broker.NatsSubscriber,
	controller *sagaOrchestratorController) broker.EventRouter {
	return &OrchestratorEventRouter{
		router:     router,
		publisher:  publisher,
		subscriber: subscriber,
		controller: controller,
	}
}

// RegisterHandlers implements broker.EventRouter.
func (r *OrchestratorEventRouter) RegisterHandlers() {

	r.router.AddNoPublisherHandler(
		"saga_orchestrator_handle_transaction_handler",
		event.PurchaseTopic,
		r.subscriber,
		r.controller.HandleTrx,
	)

	r.router.AddNoPublisherHandler(
		"saga_orchestrator_handle_reply_handler",
		event.ReplyTopic,
		r.subscriber,
		r.controller.HandleReply,
	)
}

// Run implements broker.EventRouter.
func (r *OrchestratorEventRouter) Run() error {
	r.RegisterHandlers()
	return r.router.Run(context.Background())
}

// GracefulShutdown implements broker.EventRouter.
func (r *OrchestratorEventRouter) GracefulShutdown() error {
	return r.router.Close()
}
