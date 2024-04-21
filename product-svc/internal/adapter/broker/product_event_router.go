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

type sagaProductController struct {
	productService usecase.SagaProductUseCase
}

func NewSagaProductController(productService usecase.SagaProductUseCase) *sagaProductController {
	return &sagaProductController{
		productService: productService,
	}
}

func (c *sagaProductController) HandleUpdateProductInventory(msg *message.Message) ([]*message.Message, error) {
	log.Println("handleUpdateProductInventory received message", msg.UUID)
	carrier := make(propagation.HeaderCarrier)
	carrier.Set(broker.TraceparentHeader, msg.Metadata.Get(string(constant.CtxSpanKey)))
	parentCtx := broker.TraceContext.Extract(context.Background(), carrier)
	tr := otel.Tracer("updateProductInventory")
	ctx, span := tr.Start(parentCtx, "event.UpdateProductInventory")
	defer span.End()

	purchase, pbPurchase, err := broker.DecodeCreatePurchaseCommand(msg.Payload)
	if err != nil {
		return nil, err
	}

	reply := pb.CreatePurchaseResponse{
		PurchaseId: purchase.ID,
		Purchase:   pbPurchase,
	}
	err = c.productService.UpdateProductInventory(ctx, purchase.ID, purchase.Order.PurchasedItems)
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
	replyMsg.Metadata.Set(constant.HandlerHeader, constant.UpdateProductInventoryHandler)
	broker.SetSpanContext(ctx, replyMsg)
	replyMsgs = append(replyMsgs, replyMsg)
	return replyMsgs, nil
}

func (c *sagaProductController) HandleRollbackProductInventory(msg *message.Message) ([]*message.Message, error) {
	log.Println("handleRollbackProductInventory received message", msg.UUID)
	carrier := make(propagation.HeaderCarrier)
	carrier.Set(broker.TraceparentHeader, msg.Metadata.Get(string(constant.CtxSpanKey)))
	parentCtx := broker.TraceContext.Extract(context.Background(), carrier)
	tr := otel.Tracer("updateProductInventory")
	ctx, span := tr.Start(parentCtx, "event.UpdateProductInventory")
	defer span.End()

	var cmd pb.RollbackCommand
	if err := json.Unmarshal(msg.Payload, &cmd); err != nil {
		return nil, err
	}

	reply := pb.RollbackResponse{
		UserId:     cmd.UserId,
		PurchaseId: cmd.PurchaseId,
	}
	err := c.productService.RollbackProductInventory(ctx, cmd.PurchaseId)
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
	replyMsg.Metadata.Set(constant.HandlerHeader, constant.RollbackProductInventoryHandler)
	broker.SetSpanContext(ctx, replyMsg)
	replyMsgs = append(replyMsgs, replyMsg)
	return replyMsgs, nil
}

type ProductEventRouter struct {
	router     *message.Router
	publisher  broker.NatsPublisher
	subscriber broker.NatsSubscriber
	controller *sagaProductController
}

func NewProductEventRouter(
	router *message.Router,
	publisher broker.NatsPublisher,
	subscriber broker.NatsSubscriber,
	controller *sagaProductController) broker.EventRouter {
	return &ProductEventRouter{
		router:     router,
		publisher:  publisher,
		subscriber: subscriber,
		controller: controller,
	}
}

func (r *ProductEventRouter) RegisterHandlers() {

	r.router.AddHandler(
		"saga_product_update_product_inventory_handler",
		event.UpdateProductInventoryTopic,
		r.subscriber,
		event.ReplyTopic,
		r.publisher,
		r.controller.HandleUpdateProductInventory,
	)

	r.router.AddHandler(
		"saga_product_rollback_product_inventory_handler",
		event.RollbackProductInventoryTopic,
		r.subscriber,
		event.ReplyTopic,
		r.publisher,
		r.controller.HandleRollbackProductInventory,
	)
}

func (r *ProductEventRouter) Run() error {
	r.RegisterHandlers()
	return r.router.Run(context.Background())
}

func (r *ProductEventRouter) GracefulShutdown() error {
	return r.router.Close()
}
