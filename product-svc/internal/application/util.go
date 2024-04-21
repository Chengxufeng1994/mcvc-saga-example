package application

import (
	"encoding/json"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/common/pb"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/entity"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/domain/valueobject"
	"github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EncodeDomainPurchase(purchase *entity.Purchase) *pb.CreatePurchaseCommand {
	var pbPurchasedItems []*pb.PurchasedItem
	for _, purchasedItem := range *purchase.Order.PurchasedItems {
		pbPurchasedItems = append(pbPurchasedItems, &pb.PurchasedItem{
			ProductId: purchasedItem.ProductID,
			Amount:    purchasedItem.Amount,
		})
	}
	cmd := &pb.CreatePurchaseCommand{
		PurchaseId: purchase.ID,
		Purchase: &pb.Purchase{
			Order: &pb.Order{
				UserId:         purchase.Order.UserID,
				PurchasedItems: pbPurchasedItems,
			},
			Payment: &pb.Payment{
				CurrencyCode: purchase.Payment.CurrencyCode,
				Amount:       purchase.Payment.Amount,
			},
		},
		Timestamp: timestamppb.New(time.Now()),
	}
	return cmd
}

func DecodeCreatePurchaseResponse(payload message.Payload) (*entity.CreatePurchaseResponse, error) {
	var resp pb.CreatePurchaseResponse
	if err := json.Unmarshal(payload, &resp); err != nil {
		return nil, err
	}
	purchaseID := resp.PurchaseId
	pbPurchasedItems := resp.Purchase.Order.PurchasedItems
	var purchasedItems []valueobject.PurchasedItem
	for _, pbPurchasedItem := range pbPurchasedItems {
		purchasedItems = append(purchasedItems, valueobject.PurchasedItem{
			ProductID: pbPurchasedItem.ProductId,
			Amount:    pbPurchasedItem.Amount,
		})
	}

	return &entity.CreatePurchaseResponse{
		Purchase: &entity.Purchase{
			ID: purchaseID,
			Order: &entity.Order{
				ID:             purchaseID,
				UserID:         resp.Purchase.Order.UserId,
				PurchasedItems: &purchasedItems,
			},
			Payment: &entity.Payment{
				ID:           purchaseID,
				CurrencyCode: resp.Purchase.Payment.CurrencyCode,
				Amount:       resp.Purchase.Payment.Amount,
			},
		},
		Success: resp.Success,
		Error:   resp.Error,
	}, nil
}

func DecodeRollbackResponse(payload message.Payload) (*entity.RollbackResponse, error) {
	var resp pb.RollbackResponse
	if err := json.Unmarshal(payload, &resp); err != nil {
		return nil, err
	}

	return &entity.RollbackResponse{
		UserID:     resp.UserId,
		PurchaseID: resp.PurchaseId,
		Success:    resp.Success,
		Error:      resp.Error,
	}, nil
}
