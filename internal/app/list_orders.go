package app

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"go.opentelemetry.io/otel"
)

type OrdersListOutput struct {
	Orders []pvz.Order
	Total  int32
}

func (s *Service) ListOrders(ctx context.Context, input pvz.ListOptions) (*OrdersListOutput, error) {
	tr := otel.Tracer("app.Service")
	ctx, span := tr.Start(ctx, "List_orders")
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
