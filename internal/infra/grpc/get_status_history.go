package grpc

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetStatusHistory — обрабатывает gRPC-запрос на отправку истории изменении статуса заказа.
func (s *OrdersServer) GetStatusHistory(ctx context.Context, req *pkg.GetStatusHistoryRequest) (*pkg.StatusHistoryList, error) {
	history, err := s.ordersHandler.StatusHistory(ctx, pvz.OrderID(req.OrderId))
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	var protoHistory []*pkg.StatusHistory
	for _, h := range history {
		protoHistory = append(protoHistory, &pkg.StatusHistory{
			OrderId:   req.OrderId,
			Status:    pkg.OrderStatus(h.Status),
			ChangesAt: timestamppb.New(h.Time),
		})
	}

	return &pkg.StatusHistoryList{History: protoHistory}, nil
}
