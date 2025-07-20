package json

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// GetOrder — получает заказ из БД(json) по orderID.
func (s *Storage) GetOrder(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
	for _, userOrders := range s.store {
		if order, ok := userOrders[orderID]; ok {
			return order, nil
		}
	}
	return nil, nil
}
