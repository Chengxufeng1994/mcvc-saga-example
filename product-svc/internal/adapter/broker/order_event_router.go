package broker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/common/constant"
	"github.com/Chengxufeng1994/go-saga-example/common/event"
	"github.com/Chengxufeng1994/go-saga-example/common/pb"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/broker"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/usecase"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sagaOrderController struct {
	orderService usecase.SagaOrderUseCase
}

func NewSagaOrderController(orderService usecase.SagaOrderUseCase) *sagaOrderController {
	return &sagaOrderController{
		orderService: orderService,
	}
}

func (c *sagaOrderController) HandleExecuteCreateOrder(msg *message.Message) ([]*message.Message, error) {
	log.Println("handleExecuteCreateOrder received message", msg.UUID)
	carrier := make(propagation.HeaderCarrier)
	carrier.Set(broker.TraceparentHeader, msg.Metadata.Get(string(constant.CtxSpanKey)))
	parentCtx := broker.TraceContext.Extract(context.Background(), carrier)
	tr := otel.Tracer("executeCreateOrder")
	ctx, span := tr.Start(parentCtx, "event.ExecuteCreateOrder")
	defer span.End()

	purchase, pbPurchase, err := broker.DecodeCreatePurchaseCommand(msg.Payload)
	if err != nil {
		return nil, err
	}

	reply := pb.CreatePurchaseResponse{
		PurchaseId: purchase.ID,
		Purchase:   pbPurchase,
	}
	err = c.orderService.ExecuteCreateOrder(ctx, purchase.Order)
	if err != nil {
		reply.Success = false
		reply.Error = err.Error()
	} else {
		reply.Success = true
		reply.Error = ""
	}
	reply.Timestamp = timestamppb.New(time.Now())

	payload, err := json.Marshal(&reply)
	if err != nil {
		return nil, err
	}

	var replyMsgs []*message.Message
	replyMsg := message.NewMessage(watermill.NewUUID(), payload)
	replyMsg.Metadata.Set(constant.HandlerHeader, constant.CreateOrderHandler)
	broker.SetSpanContext(ctx, replyMsg)
	replyMsgs = append(replyMsgs, replyMsg)
	return replyMsgs, nil
}

func (c *sagaOrderController) HandleRollbackCreateOrder(msg *message.Message) ([]*message.Message, error) {
	log.Println("handleRollbackCreateOrder received message", msg.UUID)
	carrier := make(propagation.HeaderCarrier)
	carrier.Set(broker.TraceparentHeader, msg.Metadata.Get(string(constant.CtxSpanKey)))
	parentCtx := broker.TraceContext.Extract(context.Background(), carrier)
	tr := otel.Tracer("rollbackCreateOrder")
	ctx, span := tr.Start(parentCtx, "event.RollbackCreateOrder")
	defer span.End()

	var cmd pb.RollbackCommand
	if err := json.Unmarshal(msg.Payload, &cmd); err != nil {
		return nil, err
	}

	reply := pb.RollbackResponse{
		UserId:     cmd.UserId,
		PurchaseId: cmd.PurchaseId,
	}
	err := c.orderService.RollbackCreateOrder(parentCtx, cmd.PurchaseId)
	if err != nil {
		reply.Success = false
		reply.Error = err.Error()
	} else {
		reply.Success = true
		reply.Error = ""
	}
	reply.Timestamp = timestamppb.New(time.Now())

	payload, err := json.Marshal(&reply)
	if err != nil {
		return nil, err
	}

	var replyMsgs []*message.Message
	replyMsg := message.NewMessage(watermill.NewUUID(), payload)
	replyMsg.Metadata.Set(constant.HandlerHeader, constant.RollbackOrderHandler)
	broker.SetSpanContext(ctx, replyMsg)
	replyMsgs = append(replyMsgs, replyMsg)
	return replyMsgs, nil
}

type OrderEventRouter struct {
	router     *message.Router
	publisher  broker.NatsPublisher
	subscriber broker.NatsSubscriber
	controller *sagaOrderController
}

func NewOrderEventRouter(router *message.Router, publisher broker.NatsPublisher, subscriber broker.NatsSubscriber, controller *sagaOrderController) broker.EventRouter {
	return &OrderEventRouter{
		router:     router,
		publisher:  publisher,
		subscriber: subscriber,
		controller: controller,
	}
}

// RegisterHandlers implements broker.EventRouter.
func (r *OrderEventRouter) RegisterHandlers() {

	r.router.AddHandler(
		"saga_order_create_order_handler",
		event.CreateOrderTopic,
		r.subscriber,
		event.ReplyTopic,
		r.publisher,
		r.controller.HandleExecuteCreateOrder,
	)

	r.router.AddHandler(
		"saga_order_rollback_order_handler",
		event.RollbackOrderTopic,
		r.subscriber,
		event.ReplyTopic,
		r.publisher,
		r.controller.HandleRollbackCreateOrder,
	)
}

// Run implements broker.EventRouter.
func (r *OrderEventRouter) Run() error {
	r.RegisterHandlers()
	return r.router.Run(context.Background())
}

// GracefulShutdown implements broker.EventRouter.
func (r *OrderEventRouter) GracefulShutdown() error {
	return r.router.Close()
}
