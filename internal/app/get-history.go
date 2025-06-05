package app

import (
	"context"

	service "gitlab.ozon.dev/gleb_evseev/homework/internal/service/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// GetHistory обрабатывает gRPC-запрос на отправку истории заказов.
func (s *OrdersServer) GetHistory(ctx context.Context, req *pvz.GetHistoryRequest) (*pvz.OrderHistoryList, error) {
	resp, err := s.ordersService.GetHistory(ctx, req)
	if err != nil {
		return nil, service.MapBusinessErrorToGRPC(err)
	}
	return resp, nil
}
