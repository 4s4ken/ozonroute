package app

import (
	"context"

	service "gitlab.ozon.dev/gleb_evseev/homework/internal/service/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// ListOrders обрабатывает gRPC-запрос на отправку списка заказов.
func (s *OrdersServer) ListOrders(ctx context.Context, req *pvz.ListOrdersRequest) (*pvz.OrdersList, error) {
	resp, err := s.ordersService.ListOrders(ctx, req)
	if err != nil {
		return nil, service.MapBusinessErrorToGRPC(err)
	}
	return resp, nil
}
