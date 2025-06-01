package storage

import (
	"context"
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

func (s *Storage) returnOrder(ctx context.Context, orderID pvz.OrderID) error {
	for _, orders := range s.store {
		if order, ok := orders[orderID]; ok {
			if order.Status == pvz.StatusReturned {
				return storageErrors.EntityValidationError("Order", fmt.Sprint(orderID))

			}
			if order.Status == pvz.StatusAccepted && order.Time.After(s.nowFn().UTC()) {
				return storageErrors.EntityExpiredError("Order", fmt.Sprint(orderID))
			}

			order.Status = pvz.StatusReturned
			now := s.nowFn().UTC()
			order.ReturnedAt = now
			order.StatusChanged = true
			return s.writeDataToFile()
		}
	}
	return storageErrors.EntityNotFoundError("Order", fmt.Sprint(orderID))
}

func (s *Storage) Return(ctx context.Context, orderID pvz.OrderID) error {
	return s.returnOrder(ctx, orderID)
}
