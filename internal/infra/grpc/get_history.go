package grpc

import (
	"context"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetHistory — обрабатывает gRPC-запрос на отправку истории заказов
func (s *OrdersServer) GetHistory(ctx context.Context, req *pkg.GetHistoryRequest) (*pkg.OrderHistoryList, error) {
	domainResp, err := s.ordersHandler.OrderHistory(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}
	protoHistory := make([]*pkg.OrderHistory, len(domainResp.History))
	for i, his := range domainResp.History {
		protoHistory[i] = mapDomainHistoryToProto(his)
	}

	return &pkg.OrderHistoryList{
		History: protoHistory,
	}, nil
}
