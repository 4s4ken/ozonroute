package app

import (
	"context"

	service "gitlab.ozon.dev/gleb_evseev/homework/internal/service/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// ListReturns обрабатывает gRPC-запрос на отправку списка вернувшихся заказов.
func (s *OrdersServer) ListReturns(ctx context.Context, req *pvz.ListReturnsRequest) (*pvz.ReturnsList, error) {
	resp, err := s.ordersService.ListReturns(ctx, req)
	if err != nil {
		return nil, service.MapBusinessErrorToGRPC(err)
	}
	return resp, nil
}
