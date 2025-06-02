package storage

import "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"

// OrderHistory внутренняя логика вывода истории изменения заказов.
func (s *Storage) OrderHistory() ([]*pvz.OrderInfo, error) {
	changedOrders := make([]*pvz.OrderInfo, 0)

	for _, userOrders := range s.store {
		for _, order := range userOrders {
			if order.StatusChanged {
				changedOrders = append(changedOrders, order)
			}
		}
	}
	return changedOrders, nil
}
