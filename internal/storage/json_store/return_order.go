package json_store

import (
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"time"
)

func (s *Storage) returnOrder(orderID pvz.OrderID) error {
	for _, orders := range s.store {
		if order, ok := orders[orderID]; ok {
			if order.Status == 2 {
				return storageErrors.EntityValidationError("Order", fmt.Sprint(orderID))

			}
			if order.Status == 1 && order.Time.After(time.Now().UTC()) {
				return storageErrors.EntityExpiredError("Order", fmt.Sprint(orderID))
			}

			order.Status = 3
			now := time.Now().UTC()
			order.ReturnedAt = now
			order.StatusChanged = true
			return s.writeDataToFile()
		}
	}
	return storageErrors.EntityNotFoundError("Order", fmt.Sprint(orderID))
}

func (s *Storage) Return(orderID uint64) error {
	return s.returnOrder(pvz.OrderID(orderID))
}
