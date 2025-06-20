package json

import (
	"context"
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// UpdateOrder обновляет заказ и сохраняет данные в файл.
func (s *Storage) UpdateOrder(ctx context.Context, order *pvz.Order) error {
	orders, ok := s.store[order.UserID]
	if !ok {
		return storageErrors.EntityNotFoundError("User", fmt.Sprint(order.UserID))
	}
	_, ok = orders[order.OrderID]
	if !ok {
		return storageErrors.EntityNotFoundError("Order", fmt.Sprint(order.OrderID))
	}
	orders[order.OrderID] = order

	if err := s.writeDataToFile(); err != nil {
		return err
	}

	return nil
}
