package app

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"sort"
	"time"
)

type OrderHistory struct {
	OrderID   uint64
	Status    pvz.OrderStatus
	CreatedAt time.Time
}
type OrderHistoryList struct {
	History []OrderHistory
}

func (s *Storage) OrderHistory(ctx context.Context) (*OrderHistoryList, error) {
	history := make([]OrderHistory, 0)

	for _, userOrders := range s.store {
		for _, order := range userOrders {
			history = append(history, OrderHistory{
				OrderID:   uint64(order.OrderID),
				Status:    order.Status,
				CreatedAt: order.LastUpdate,
			})
		}
	}

	sort.Slice(history, func(i, j int) bool {
		return history[i].CreatedAt.After(history[j].CreatedAt)
	})

	return &OrderHistoryList{History: history}, nil
}
