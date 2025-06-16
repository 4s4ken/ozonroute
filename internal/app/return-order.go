package app

import (
	"context"
	"fmt"

	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// returnOrder возврата заказа курьеру.
func (s *Storage) returnOrder(ctx context.Context, order pvz.OrderCreate) error {
	for _, orders := range s.store {
		if order, ok := orders[order.OrderID]; ok {
			if order.Status == pvz.OrderStatusReturned {
				return storageErrors.EntityValidationError("Order", fmt.Sprint(order.OrderID))
			}
			if order.Status == pvz.OrderStatusAccepted && order.ExpiresAt.Before(s.nowFn().UTC()) {
				return storageErrors.EntityExpiredError("Order", fmt.Sprint(order.OrderID))
			}
			order.Status = pvz.OrderStatusReturned
			order.ExpiresAt = s.nowFn().UTC()

			return s.writeDataToFile()
		}
	}
	return storageErrors.EntityNotFoundError("Order", fmt.Sprint(order.OrderID))
}

func (s *Storage) Return(ctx context.Context, input pvz.OrderCreate) error {
	return s.returnOrder(ctx, input)
}
