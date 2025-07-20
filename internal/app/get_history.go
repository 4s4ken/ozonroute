package app

import (
	"context"
	"fmt"
	tracerpackage "gitlab.ozon.dev/gleb_evseev/homework/pkg/tracer"
	"go.opentelemetry.io/otel"
	"sort"
	"time"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type OrderHistory struct {
	OrderID uint64
	Status  pvz.OrderStatus
	Time    time.Time
}

type OrderHistoryList struct {
	History []OrderHistory
}

// OrderHistory — возвращает список изменений статусов заказов.
func (s *Service) OrderHistory(ctx context.Context) (*OrderHistoryList, error) {
	ctx, span := tracerpackage.Start(ctx, "Orders_history")
	defer span.End()

	rawHistory, err := s.repo.GetHistory(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]OrderHistory, 0, len(rawHistory))
	for _, item := range rawHistory {
		result = append(result, OrderHistory{
			OrderID: uint64(item.OrderID),
			Status:  item.Status,
			Time:    item.Time,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.After(result[j].Time)
	})

	return &OrderHistoryList{History: result}, nil
}

// StatusHistory — возвращает список изменений статуса одного заказа.
func (s *Service) StatusHistory(ctx context.Context, orderID pvz.OrderID) ([]*pvz.StatusChange, error) {
	tr := otel.Tracer("app.Service")
	ctx, span := tr.Start(ctx, "Order_status_history")
	defer span.End()

	rawHistory, err := s.repo.GetStatusHistory(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if len(rawHistory) == 0 {
		return nil, fmt.Errorf("order %d not found", orderID)
	}

	return rawHistory, nil
}
