package app

import (
	"context"

	service "gitlab.ozon.dev/gleb_evseev/homework/internal/service/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// ReturnOrder обрабатывает gRPC-запрос на возвращение заказа.
func (s *OrdersServer) ReturnOrder(ctx context.Context, req *pvz.OrderIdRequest) (*pvz.OrderResponse, error) {
	resp, err := s.ordersService.ReturnOrder(ctx, req)
	if err != nil {
		return nil, service.MapBusinessErrorToGRPC(err)
	}
	return resp, nil
}
