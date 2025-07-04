package json

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// GetOrder возвращает заказ по OrderID.
func (s *Storage) GetOrder(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
	for _, userOrders := range s.store {
		if order, ok := userOrders[orderID]; ok {
			return order, nil
		}
	}
	return nil, nil
}
