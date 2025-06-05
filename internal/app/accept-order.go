package app

import (
	"context"

	service "gitlab.ozon.dev/gleb_evseev/homework/internal/service/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// AcceptOrder обрабатывает gRPC-запрос на приём заказа.
func (s *OrdersServer) AcceptOrder(ctx context.Context, req *pvz.AcceptOrderRequest) (*pvz.OrderResponse, error) {
	resp, err := s.ordersService.AcceptOrder(ctx, req)
	if err != nil {
		return nil, service.MapBusinessErrorToGRPC(err)
	}
	return resp, nil
}
