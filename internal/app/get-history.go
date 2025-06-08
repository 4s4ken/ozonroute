package app

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

func (s *Storage) OrderHistory(ctx context.Context) (*pvz.OrderHistoryList, error) {
	history := make([]pvz.OrderHistory, 0)

	for _, userOrders := range s.store {
		for _, order := range userOrders {
			history = append(history, pvz.OrderHistory{
				OrderID:   uint64(order.OrderID),
				Status:    order.Status,
				CreatedAt: order.ExpiresAt,
			})
		}
	}

	return &pvz.OrderHistoryList{History: history}, nil
}
