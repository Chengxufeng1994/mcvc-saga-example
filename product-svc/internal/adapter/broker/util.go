package broker

import (
	"time"

	"github.com/Chengxufeng1994/go-saga-example/common/pb"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/event"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EncodeDomainPurchaseResult(purchaseResult *event.PurchaseResultEvent) *pb.PurchaseResult {
	step := getPbPurchaseStep(purchaseResult.Step)
	status := getPbPurchaseStatus(purchaseResult.Status)
	return &pb.PurchaseResult{
		UserId:     purchaseResult.UserID,
		PurchaseId: purchaseResult.PurchaseID,
		Step:       step,
		Status:     status,
		Timestamp:  timestamppb.New(time.Now()),
	}
}

func getPbPurchaseStep(step string) pb.PurchaseStep {
	switch step {
	case event.StepUpdateProductInventory:
		return pb.PurchaseStep_STEP_UPDATE_PRODUCT_INVENTORY
	case event.StepCreateOrder:
		return pb.PurchaseStep_STEP_CREATE_ORDER
	case event.StepCreatePayment:
		return pb.PurchaseStep_STEP_CREATE_PAYMENT
	}
	return -1
}

func getPbPurchaseStatus(status string) pb.PurchaseStatus {
	switch status {
	case event.StatusExecute:
		return pb.PurchaseStatus_STATUS_EXUCUTE
	case event.StatusSucess:
		return pb.PurchaseStatus_STATUS_SUCCESS
	case event.StatusFailed:
		return pb.PurchaseStatus_STATUS_FAILED
	case event.StatusRollbacked:
		return pb.PurchaseStatus_STATUS_ROLLBACKED
	case event.StatusRollbackFailed:
		return pb.PurchaseStatus_STATUS_ROLLBACK_FAIL
	}
	return -1
}
