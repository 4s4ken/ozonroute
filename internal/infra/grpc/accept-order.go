package grpc

import (
	"context"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AcceptOrder обрабатывает gRPC-запрос на приём заказа.
func (s *OrdersServer) AcceptOrder(ctx context.Context, req *pkg.AcceptOrderRequest) (*pkg.OrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	input := mapAcceptOrderRequestToInput(req)

	err := s.ordersHandler.Accept(ctx, input)
	if err != nil {
		return nil, MapBusinessErrorToGRPC(err)
	}

	return &pkg.OrderResponse{
		Status:  pkg.OrderStatus_ORDER_STATUS_ACCEPTED,
		OrderId: input.OrderID,
	}, nil
}
