package app

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type OrdersListOutput struct {
	Orders []pvz.Order
	Total  int32
}

func (s *Service) ListOrders(ctx context.Context, input pvz.ListOptions) (*OrdersListOutput, error) {
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
