package grpc

import (
	"context"

	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// ReturnOrder обрабатывает gRPC-запрос на возвращение заказа.
func (s *OrdersServer) ReturnOrder(ctx context.Context, req *pkg.ReturnOrderRequest) (*pkg.OrderResponse, error) {
	input := mapReturnOrderRequestToInput(req)

	err := s.ordersHandler.Return(ctx, input)
	if err != nil {
		return nil, MapBusinessErrorToGRPC(err)
	}

	return &pkg.OrderResponse{
		Status:  pkg.OrderStatus_ORDER_STATUS_RETURNED,
		OrderId: uint64(input.OrderID),
	}, nil
}
