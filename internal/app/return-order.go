package app

import (
	"context"
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// returnOrder возврата заказа курьеру.
func (s *Storage) returnOrder(ctx context.Context, input pvz.ReturnOrderInput) error {
	for _, orders := range s.store {
		if order, ok := orders[input.OrderID]; ok {
			if order.Status == pvz.OrderStatusReturned {
				return storageErrors.EntityValidationError("Order", fmt.Sprint(input.OrderID))
			}
			if order.Status == pvz.OrderStatusAccepted && order.ExpiresAt.Before(s.nowFn().UTC()) {
				return storageErrors.EntityExpiredError("Order", fmt.Sprint(input.OrderID))
			}
			order.Status = pvz.OrderStatusReturned
			order.ExpiresAt = s.nowFn().UTC()

			return s.writeDataToFile()
		}
	}
	return storageErrors.EntityNotFoundError("Order", fmt.Sprint(input.OrderID))
}

func (s *Storage) Return(ctx context.Context, input pvz.ReturnOrderInput) error {
	return s.returnOrder(ctx, input)
}
