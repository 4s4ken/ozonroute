package grpc

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReturnOrder — обрабатывает gRPC-запрос на возвращение заказа.
func (s *OrdersServer) ReturnOrder(ctx context.Context, req *pkg.ReturnOrderRequest) (*pkg.OrderResponse, error) {
	input := app.ReturnOrderInput{
		OrderID: pvz.OrderID(req.GetOrderId()),
	}
	err := s.ordersHandler.Return(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	return &pkg.OrderResponse{
		Status:  pkg.OrderStatus_ORDER_STATUS_RETURNED,
		OrderId: uint64(input.OrderID),
	}, nil
}
