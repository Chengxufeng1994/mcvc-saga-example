package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/common/constant"
	"github.com/Chengxufeng1994/go-saga-example/common/event"
	"github.com/Chengxufeng1994/go-saga-example/common/pb"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
	domainevent "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/event"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/broker"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/repository"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/usecase"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrchestratorService struct {
	logger                   *logrus.Entry
	natsPublisher            broker.NatsPublisher
	purchaseResultRepository repository.PurchaseResultRepository
}

func NewOrchestratorService(natsPublisher broker.NatsPublisher, purchaseResultRepository repository.PurchaseResultRepository) usecase.OrchestratorUseCase {
	return &OrchestratorService{
		logger:                   config.ContextLogger.WithFields(logrus.Fields{"type": "service:OrchestratorService"}),
		natsPublisher:            natsPublisher,
		purchaseResultRepository: purchaseResultRepository,
	}
}

// HandleTrx implements usecase.OrchestratorUseCase.
// start transaction starts the first transaction, which is UpdateProductInventory
func (svc *OrchestratorService) HandleTrx(parentCtx context.Context, purchase *entity.Purchase, correlationID string) error {
	tr := otel.Tracer("startTransaction")
	ctx, span := tr.Start(parentCtx, "event.StartTransaction")
	defer span.End()

	cmd := EncodeDomainPurchase(purchase)
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	middleware.SetCorrelationID(correlationID, msg)
	svc.logger.Infof("update product inventory %v", purchase.ID)
	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID,
		domainevent.NewPurchaseResultEvent(purchase.Order.UserID, cmd.PurchaseId, domainevent.StepUpdateProductInventory, domainevent.StatusExecute),
	)

	updateProductInventoryTopicEvt := &domainevent.Event{
		Topic:       event.UpdateProductInventoryTopic,
		MessageType: domainevent.TRX_MSG,
		Message:     msg,
	}
	return svc.publishEvent(ctx, updateProductInventoryTopicEvt)
}

// HandleReply implements usecase.OrchestratorUseCase.
func (svc *OrchestratorService) HandleReply(parentCtx context.Context, msg *message.Message, correlationID string) error {
	tr := otel.Tracer("handleReply")
	ctx, span := tr.Start(parentCtx, "event.HandleReply")
	defer span.End()

	handler := msg.Metadata.Get(constant.HandlerHeader)
	switch handler {
	case constant.UpdateProductInventoryHandler:
		resp, err := DecodeCreatePurchaseResponse(msg.Payload)
		if err != nil {
			return err
		}
		if resp.Success {
			return svc.createOrder(ctx, resp.Purchase, correlationID)
		}
		svc.logger.WithError(err).Error(resp.Error)
		return svc.rollbackProductInventory(ctx, resp.Purchase.Order.UserID, resp.Purchase.ID, correlationID)
	case constant.RollbackProductInventoryHandler:
		resp, err := DecodeRollbackResponse(msg.Payload)
		if err != nil {
			return err
		}
		return svc.purchaseResultRepository.PublishPurchaseResult(
			correlationID, domainevent.NewPurchaseResultEvent(resp.UserID, resp.PurchaseID, domainevent.StepUpdateProductInventory, domainevent.StatusRollbackFailed))
	case constant.CreateOrderHandler:
		resp, err := DecodeCreatePurchaseResponse(msg.Payload)
		if err != nil {
			return err
		}
		if resp.Success {
			return svc.createPayment(ctx, resp.Purchase, correlationID)
		}
		svc.logger.WithError(err).Error(resp.Error)
		return svc.rollbackFromOrder(ctx, resp.Purchase.Order.UserID, resp.Purchase.ID, correlationID)
	case constant.RollbackOrderHandler:
		resp, err := DecodeRollbackResponse(msg.Payload)
		if err != nil {
			return err
		}
		return svc.purchaseResultRepository.PublishPurchaseResult(
			correlationID, domainevent.NewPurchaseResultEvent(resp.UserID, resp.PurchaseID, domainevent.StepCreateOrder, domainevent.StatusRollbackFailed))
	case constant.CreatePaymentHandler:
		resp, err := DecodeCreatePurchaseResponse(msg.Payload)
		if err != nil {
			return err
		}
		if resp.Success {
			return svc.purchaseResultRepository.PublishPurchaseResult(
				correlationID, domainevent.NewPurchaseResultEvent(resp.Purchase.Order.UserID, resp.Purchase.ID, domainevent.StepCreatePayment, domainevent.StatusSucess))
		}
		svc.logger.WithError(err).Error(resp.Error)
		return svc.rollbackFromPayment(ctx, resp.Purchase.Order.UserID, resp.Purchase.ID, correlationID)
	case constant.RollbackPaymentHandler:
		resp, err := DecodeRollbackResponse(msg.Payload)
		if err != nil {
			return err
		}
		return svc.purchaseResultRepository.PublishPurchaseResult(
			correlationID, domainevent.NewPurchaseResultEvent(resp.UserID, resp.PurchaseID, domainevent.StepCreatePayment, domainevent.StatusRollbackFailed))
	default:
		return nil
	}
}

func (svc *OrchestratorService) createOrder(ctx context.Context, purchase *entity.Purchase, correlationID string) error {
	svc.logger.Infof("create order %v", purchase.ID)
	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID, domainevent.NewPurchaseResultEvent(purchase.Order.UserID, purchase.ID, domainevent.StepUpdateProductInventory, domainevent.StatusSucess))

	cmd := EncodeDomainPurchase(purchase)
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	middleware.SetCorrelationID(correlationID, msg)

	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID, domainevent.NewPurchaseResultEvent(purchase.Order.UserID, purchase.ID, domainevent.StepCreateOrder, domainevent.StatusExecute))

	createOrderTopicEvt := &domainevent.Event{
		Topic:       event.CreateOrderTopic,
		MessageType: domainevent.TRX_MSG,
		Message:     msg,
	}
	return svc.publishEvent(ctx, createOrderTopicEvt)
}

func (svc *OrchestratorService) createPayment(ctx context.Context, purchase *entity.Purchase, correlationID string) error {
	svc.logger.Infof("create payment %v", purchase.ID)
	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID, domainevent.NewPurchaseResultEvent(purchase.Order.UserID, purchase.ID, domainevent.StepCreateOrder, domainevent.StatusSucess))

	cmd := EncodeDomainPurchase(purchase)
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	middleware.SetCorrelationID(correlationID, msg)

	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID, domainevent.NewPurchaseResultEvent(purchase.Order.UserID, purchase.ID, domainevent.StepCreatePayment, domainevent.StatusExecute))

	createPaymentTopicEvt := &domainevent.Event{
		Topic:       event.CreatePaymentTopic,
		MessageType: domainevent.TRX_MSG,
		Message:     msg,
	}
	return svc.publishEvent(ctx, createPaymentTopicEvt)
}

func (svc *OrchestratorService) rollbackProductInventory(ctx context.Context, userID uint64, purchaseID uint64, correlationID string) error {
	svc.logger.Infof("rollback product inventory %v", purchaseID)
	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID,
		domainevent.NewPurchaseResultEvent(userID, purchaseID, domainevent.StepUpdateProductInventory, domainevent.StatusFailed),
	)

	cmd := &pb.RollbackCommand{
		PurchaseId: purchaseID,
		Timestamp:  timestamppb.New(time.Now()),
	}
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	middleware.SetCorrelationID(correlationID, msg)
	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID,
		domainevent.NewPurchaseResultEvent(userID, purchaseID, domainevent.StepUpdateProductInventory, domainevent.StatusRollbackFailed),
	)

	rollbackProductInventoryTopicEvt := &domainevent.Event{
		Topic:       event.RollbackProductInventoryTopic,
		MessageType: domainevent.TRX_MSG,
		Message:     msg,
	}
	return svc.publishEvent(ctx, rollbackProductInventoryTopicEvt)
}

func (svc *OrchestratorService) rollbackCreateOrder(ctx context.Context, userID uint64, purchaseID uint64, correlationID string) error {
	svc.logger.Infof("rollback create order %v", purchaseID)
	cmd := &pb.RollbackCommand{
		PurchaseId: purchaseID,
		Timestamp:  timestamppb.New(time.Now()),
	}
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	middleware.SetCorrelationID(correlationID, msg)

	rollbackCreateOrderTopicEvt := &domainevent.Event{
		Topic:       event.RollbackOrderTopic,
		MessageType: domainevent.TRX_MSG,
		Message:     msg,
	}
	return svc.publishEvent(ctx, rollbackCreateOrderTopicEvt)
}

func (svc *OrchestratorService) rollbackCreatePayment(ctx context.Context, userID uint64, purchaseID uint64, correlationID string) error {
	svc.logger.Infof("rollback create payment %v", purchaseID)
	cmd := &pb.RollbackCommand{
		PurchaseId: purchaseID,
		Timestamp:  timestamppb.New(time.Now()),
	}
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	middleware.SetCorrelationID(correlationID, msg)

	rollbackPaymentTopicEvt := &domainevent.Event{
		Topic:       event.RollbackPaymentTopic,
		MessageType: domainevent.TRX_MSG,
		Message:     msg,
	}
	return svc.publishEvent(ctx, rollbackPaymentTopicEvt)
}

func (svc *OrchestratorService) rollbackFromOrder(ctx context.Context, userID uint64, purchaseID uint64, correlationID string) error {
	var err error
	svc.logger.Infof("rollback from order %v", purchaseID)
	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID,
		domainevent.NewPurchaseResultEvent(userID, purchaseID, domainevent.StepCreateOrder, domainevent.StatusFailed),
	)

	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID,
		domainevent.NewPurchaseResultEvent(userID, purchaseID, domainevent.StepCreateOrder, domainevent.StatusRollbacked),
	)
	if err = svc.rollbackCreateOrder(ctx, userID, purchaseID, correlationID); err != nil {
		svc.logger.WithError(err).Error(err.Error())
	}

	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID,
		domainevent.NewPurchaseResultEvent(userID, purchaseID, domainevent.StepUpdateProductInventory, domainevent.StatusRollbacked),
	)
	if err = svc.rollbackCreatePayment(ctx, userID, purchaseID, correlationID); err != nil {
		svc.logger.WithError(err).Error(err.Error())
	}

	return err
}

func (svc *OrchestratorService) rollbackFromPayment(ctx context.Context, userID uint64, purchaseID uint64, correlationID string) error {
	var err error
	svc.logger.Infof("rollback from payment %v", purchaseID)
	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID,
		domainevent.NewPurchaseResultEvent(userID, purchaseID, domainevent.StepCreatePayment, domainevent.StatusFailed),
	)

	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID,
		domainevent.NewPurchaseResultEvent(userID, purchaseID, domainevent.StepCreatePayment, domainevent.StatusRollbacked),
	)
	if err = svc.rollbackCreatePayment(ctx, userID, purchaseID, correlationID); err != nil {
		svc.logger.WithError(err).Error(err.Error())
	}

	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID,
		domainevent.NewPurchaseResultEvent(userID, purchaseID, domainevent.StepCreateOrder, domainevent.StatusRollbacked),
	)
	if err = svc.rollbackCreateOrder(ctx, userID, purchaseID, correlationID); err != nil {
		svc.logger.WithError(err).Error(err.Error())
	}

	svc.purchaseResultRepository.PublishPurchaseResult(
		correlationID,
		domainevent.NewPurchaseResultEvent(userID, purchaseID, domainevent.StepUpdateProductInventory, domainevent.StatusRollbacked),
	)
	if err = svc.rollbackCreatePayment(ctx, userID, purchaseID, correlationID); err != nil {
		svc.logger.WithError(err).Error(err.Error())
	}

	return err
}

func (svc *OrchestratorService) publishEvent(ctx context.Context, evt *domainevent.Event) error {
	broker.SetSpanContext(ctx, evt.Message)
	switch evt.MessageType {
	case domainevent.TRX_MSG:
		return svc.natsPublisher.Publish(evt.Topic, evt.Message)
	default:
		return fmt.Errorf("unknown message type: %v", evt.MessageType)
	}
}
