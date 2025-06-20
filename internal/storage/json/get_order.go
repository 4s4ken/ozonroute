package json

import (
	"context"
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// GetOrder возвращает заказ по OrderID.
func (s *Storage) GetOrder(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
	for _, userOrders := range s.store {
		if order, ok := userOrders[orderID]; ok {
			return order, nil
		}
	}
	return nil, storageErrors.EntityNotFoundError("Order", fmt.Sprint(orderID))
}
