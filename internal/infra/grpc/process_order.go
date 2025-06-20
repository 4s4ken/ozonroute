package grpc

import (
	"context"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// ProcessOrders обрабатывает gRPC-запрос на обработку заказа.
func (s *OrdersServer) ProcessOrders(ctx context.Context, req *pkg.ProcessOrdersRequest) (*pkg.ProcessResult, error) {
	input := mapProcessOrdersRequestToInput(req)

	resp, err := s.ordersHandler.Process(ctx, input)
	if err != nil {
		return nil, err
	}
	return mapProcessResultToProto(resp), nil
}
