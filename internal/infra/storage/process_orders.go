package storage

import (
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"time"
)

func (s *Storage) process(userID pvz.UserID, act string, orderIDs []pvz.OrderID) error {
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
			if order.IsExpired(time.Now().UTC()) {
				return storageErrors.EntityExpiredError("Order", fmt.Sprint(orderID))
			}
			order.Status = 2
			order.StatusChanged = true

		case "return":
			if order.IsExpiredForReturn(time.Now().UTC()) {
				return storageErrors.EntityExpiredError("Order", fmt.Sprint(orderID))
			}
			order.Status = 3
			order.StatusChanged = true
			order.ReturnedAt = time.Now()
		default:
			return storageErrors.EntityValidationError("Order", fmt.Sprint(orderID))
		}
	}
	orders = nil
	return s.writeDataToFile()
}

func (s *Storage) Process(userID uint64, act string, orderIDs []pvz.OrderID) error {
	pvzOrderIDs := make([]pvz.OrderID, len(orderIDs))
	for i, id := range orderIDs {
		pvzOrderIDs[i] = id
	}
	return s.process(pvz.UserID(userID), act, pvzOrderIDs)
}
