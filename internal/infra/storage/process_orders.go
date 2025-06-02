package storage

import (
	"context"
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// Process внутренняя логика выдачи/возврата заказа клиенту.
func (s *Storage) process(ctx context.Context, userID pvz.UserID, act string, orderIDs []pvz.OrderID) error {
	orders, ok := s.store[userID]
	if !ok {
		return storageErrors.EntityNotFoundError("User", fmt.Sprint(orders))
	}

	for _, orderID := range orderIDs {
		order, ok := orders[orderID]
		if !ok {
			return storageErrors.EntityNotFoundError("Order", fmt.Sprint(orderID))
		}

		switch act {
		case "issue":
			if order.IsExpired(s.nowFn().UTC()) {
				return storageErrors.EntityExpiredError("Order", fmt.Sprint(orderID))
			}
			order.Status = pvz.StatusAccepted
			order.StatusChanged = true

		case "return":
			if order.IsExpiredForReturn(s.nowFn().UTC()) {
				return storageErrors.EntityExpiredError("Order", fmt.Sprint(orderID))
			}
			order.Status = pvz.StatusIssued
			order.StatusChanged = true
			order.ReturnedAt = s.nowFn()
		default:
			return storageErrors.EntityValidationError("Order", fmt.Sprint(orderID))
		}
	}
	return s.writeDataToFile()
}

func (s *Storage) Process(ctx context.Context, userID pvz.UserID, act string, orderIDs []pvz.OrderID) error {
	pvzOrderIDs := make([]pvz.OrderID, len(orderIDs))
	for i, id := range orderIDs {
		pvzOrderIDs[i] = id
	}
	return s.process(ctx, userID, act, pvzOrderIDs)
}
