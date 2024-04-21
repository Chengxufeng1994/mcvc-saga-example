package event

import (
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageType int

const (
	TRX_MSG MessageType = iota
	RESULT_MSG
)

type Event struct {
	Topic       string
	MessageType MessageType
	Message     *message.Message
}

var (
	StepUpdateProductInventory = "UPDATE_PRODUCT_INVENTORY"
	StepCreateOrder            = "CREATE_ORDER"
	StepCreatePayment          = "CREATE_PAYMENT"

	StatusExecute        = "STATUS_EXUCUTE"
	StatusSucess         = "STATUS_SUCCESS"
	StatusFailed         = "STATUS_FAILED"
	StatusRollbacked     = "STATUS_ROLLBACKED"
	StatusRollbackFailed = "STATUS_ROLLBACK_FAIL"
)

// PurchaseResult event
type PurchaseResultEvent struct {
	UserID     uint64
	PurchaseID uint64
	Step       string
	Status     string
	Timestamp  time.Time
}

func NewPurchaseResultEvent(userID, purchaseID uint64, step, status string) *PurchaseResultEvent {
	return &PurchaseResultEvent{
		UserID:     userID,
		PurchaseID: purchaseID,
		Step:       step,
		Status:     status,
		Timestamp:  time.Now(),
	}
}
