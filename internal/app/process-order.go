package app

import (
	"context"
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// ProcessOrders — обработка заказов покупателей (выдача/возврат).
func (s *Storage) ProcessOrders(ctx context.Context, orders pvz.ProcessOrdersInput) error {
	ordersByUser, ok := s.store[orders.UserID]
	if !ok {
		return storageErrors.EntityNotFoundError("User", fmt.Sprint(orders.UserID))
	}

	for _, orderID := range orders.OrderIDs {
		order, ok := ordersByUser[orderID]
		if !ok {
			return storageErrors.EntityNotFoundError("Order", fmt.Sprint(orderID))
		}

		switch orders.Action {
		case pvz.ActionTypeIssue:
			if order.IsExpired(s.nowFn().UTC()) {
				return storageErrors.EntityExpiredError("Order", fmt.Sprint(orderID))
			}
			order.Status = pvz.OrderStatusIssued

		case pvz.ActionTypeReturn:
			if order.IsExpiredForReturn(s.nowFn().UTC()) {
				return storageErrors.EntityExpiredError("Order", fmt.Sprint(orderID))
			}
			order.Status = pvz.OrderStatusReturned
			order.ExpiresAt = s.nowFn()

		default:
			return storageErrors.EntityValidationError("Action", fmt.Sprint(orders.Action))
		}
	}

	return s.writeDataToFile()
}

func (s *Storage) Process(ctx context.Context, orders pvz.ProcessOrdersInput) (pvz.ProcessResult, error) {
	err := s.ProcessOrders(ctx, orders)
	if err != nil {
		return pvz.ProcessResult{}, err
	}

	processed := make([]uint64, 0, len(orders.OrderIDs))
	for _, id := range orders.OrderIDs {
		processed = append(processed, uint64(id))
	}

	return pvz.ProcessResult{
		Processed: processed,
	}, nil
}
