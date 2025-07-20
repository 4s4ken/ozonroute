package app

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	tracerpackage "gitlab.ozon.dev/gleb_evseev/homework/pkg/tracer"
)

type OrdersListOutput struct {
	Orders []pvz.Order
	Total  int32
}

// ListOrders — возвращает список заказов c фильтрами [IN-PVZ, Page & limit, LastN].
func (s *Service) ListOrders(ctx context.Context, input pvz.ListOptions) (*OrdersListOutput, error) {
	ctx, span := tracerpackage.Start(ctx, "List_orders")
	defer span.End()

	orders, total, err := s.repo.ListOrders(ctx, input)
	if err != nil {
		return nil, err
	}

	result := make([]pvz.Order, 0, len(orders))
	for _, o := range orders {
		result = append(result, *o)
	}

	return &OrdersListOutput{
		Orders: result,
		Total:  total,
	}, nil
}
