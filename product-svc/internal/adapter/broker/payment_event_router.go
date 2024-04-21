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

type sagaPaymentController struct {
	paymentService usecase.SagaPaymentUseCase
}

func NewSagaPaymentController(paymentService usecase.SagaPaymentUseCase) *sagaPaymentController {
	return &sagaPaymentController{
		paymentService: paymentService,
	}
}

func (c *sagaPaymentController) HandleExecuteCreatePayment(msg *message.Message) ([]*message.Message, error) {
	log.Println("handleExecuteCreateOrder received message", msg.UUID)
	carrier := make(propagation.HeaderCarrier)
	carrier.Set(broker.TraceparentHeader, msg.Metadata.Get(string(constant.CtxSpanKey)))
	parentCtx := broker.TraceContext.Extract(context.Background(), carrier)
	tr := otel.Tracer("executeCreatePayment")
	ctx, span := tr.Start(parentCtx, "event.ExecuteCreatePayment")
	defer span.End()

	purchase, pbPurchase, err := broker.DecodeCreatePurchaseCommand(msg.Payload)
	if err != nil {
		return nil, err
	}

	reply := pb.CreatePurchaseResponse{
		PurchaseId: purchase.ID,
		Purchase:   pbPurchase,
	}
	err = c.paymentService.ExecuteCreatePayment(ctx, purchase.Payment)
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
	replyMsg.Metadata.Set(constant.HandlerHeader, constant.CreatePaymentHandler)
	broker.SetSpanContext(ctx, replyMsg)
	replyMsgs = append(replyMsgs, replyMsg)
	return replyMsgs, nil
}

func (c *sagaPaymentController) HandleRollbackCreatePayment(msg *message.Message) ([]*message.Message, error) {
	log.Println("handleRollbackCreatePayment received message", msg.UUID)
	carrier := make(propagation.HeaderCarrier)
	carrier.Set(broker.TraceparentHeader, msg.Metadata.Get(string(constant.CtxSpanKey)))
	parentCtx := broker.TraceContext.Extract(context.Background(), carrier)
	tr := otel.Tracer("rollbackCreatePayment")
	ctx, span := tr.Start(parentCtx, "event.RollbackCreatePayment")
	defer span.End()

	var cmd pb.RollbackCommand
	if err := json.Unmarshal(msg.Payload, &cmd); err != nil {
		return nil, err
	}

	reply := pb.RollbackResponse{
		UserId:     cmd.UserId,
		PurchaseId: cmd.PurchaseId,
	}
	err := c.paymentService.RollbackCreatePayment(parentCtx, cmd.PurchaseId)
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
	replyMsg.Metadata.Set(constant.HandlerHeader, constant.RollbackPaymentHandler)
	broker.SetSpanContext(ctx, replyMsg)
	replyMsgs = append(replyMsgs, replyMsg)
	return replyMsgs, nil
}

type PaymentEventRouter struct {
	router     *message.Router
	publisher  broker.NatsPublisher
	subscriber broker.NatsSubscriber
	controller *sagaPaymentController
}

func NewPaymentEventRouter(
	router *message.Router,
	publisher broker.NatsPublisher,
	subscriber broker.NatsSubscriber,
	controller *sagaPaymentController) broker.EventRouter {
	return &PaymentEventRouter{
		router:     router,
		publisher:  publisher,
		subscriber: subscriber,
		controller: controller,
	}
}

func (r *PaymentEventRouter) RegisterHandlers() {

	r.router.AddHandler(
		"saga_payment_create_payment_handler",
		event.CreatePaymentTopic,
		r.subscriber,
		event.ReplyTopic,
		r.publisher,
		r.controller.HandleExecuteCreatePayment,
	)

	r.router.AddHandler(
		"saga_payment_rollback_payment_handler",
		event.RollbackPaymentTopic,
		r.subscriber,
		event.ReplyTopic,
		r.publisher,
		r.controller.HandleRollbackCreatePayment,
	)

}

func (r *PaymentEventRouter) Run() error {
	r.RegisterHandlers()
	return r.router.Run(context.Background())
}

func (r *PaymentEventRouter) GracefulShutdown() error {
	return r.router.Close()
}
