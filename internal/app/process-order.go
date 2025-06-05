package app

import (
	"context"

	service "gitlab.ozon.dev/gleb_evseev/homework/internal/service/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// ProcessOrder обрабатывает gRPC-запрос на обработку заказа.
func (s *OrdersServer) ProcessOrder(ctx context.Context, req *pvz.ProcessOrdersRequest) (*pvz.ProcessResult, error) {
	resp, err := s.ordersService.ProcessOrders(ctx, req)
	if err != nil {
		return nil, service.MapBusinessErrorToGRPC(err)
	}
	return resp, nil
}
