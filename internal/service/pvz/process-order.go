package pvz

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (s *OrdersService) ProcessOrders(ctx context.Context, req *pvz.ProcessOrdersRequest) (*pvz.ProcessResult, error) {
	if req == nil {
		return nil, businesserror.EntityValidationError("ProcessOrdersRequest", "")
	}

	now := time.Now()
	processed := make([]uint64, 0, len(req.OrderIds))
	errorsList := make([]uint64, 0)

	for _, orderID := range req.OrderIds {
		order, exists := s.storage[orderID]
		if !exists || order.UserId != req.UserId {
			errorsList = append(errorsList, orderID)
			continue
		}

		switch req.Action {
		case pvz.ActionType_ACTION_TYPE_ISSUE:
			if order.Status != pvz.OrderStatus_ORDER_STATUS_ACCEPTED {
				errorsList = append(errorsList, orderID)
				continue
			}
			if order.ExpiresAt == nil || order.ExpiresAt.AsTime().Before(now) {
				errorsList = append(errorsList, orderID)
				continue
			}
			expireAt := now.Add(48 * time.Hour)
			order.ExpiresAt = timestamppb.New(expireAt)
			order.Status = pvz.OrderStatus_ORDER_STATUS_EXPECTS
			processed = append(processed, orderID)

		case pvz.ActionType_ACTION_TYPE_RETURN:
			if order.Status != pvz.OrderStatus_ORDER_STATUS_EXPECTS {
				errorsList = append(errorsList, orderID)
				continue
			}
			if order.ExpiresAt == nil || order.ExpiresAt.AsTime().Before(now) {
				errorsList = append(errorsList, orderID)
				continue
			}
			order.Status = pvz.OrderStatus_ORDER_STATUS_RETURNED
			processed = append(processed, orderID)

		default:
			errorsList = append(errorsList, orderID)
		}
	}

	return &pvz.ProcessResult{
		Processed: processed,
		Errors:    errorsList,
	}, nil
}
