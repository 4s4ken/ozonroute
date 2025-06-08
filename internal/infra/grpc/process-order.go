package grpc

import (
	"context"

	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// ProcessOrder обрабатывает gRPC-запрос на обработку заказа.
func (s *OrdersServer) ProcessOrder(ctx context.Context, req *pkg.ProcessOrdersRequest) (*pkg.ProcessResult, error) {
	input := mapProcessOrdersRequestToInput(req)

	resp, err := s.ordersHandler.Process(ctx, input)
	if err != nil {
		return nil, MapBusinessErrorToGRPC(err)
	}
	return mapProcessResultToProto(resp), nil
}
