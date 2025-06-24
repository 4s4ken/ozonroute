package grpc

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListOrders обрабатывает gRPC-запрос на отправку списка заказов.
func (s *OrdersServer) ListOrders(ctx context.Context, req *pkg.ListOrdersRequest) (*pkg.OrdersList, error) {
	var lastN *uint64
	if req.LastN != nil {
		tmp := uint64(*req.LastN)
		lastN = &tmp
	}

	var pagination *pvz.Pagination
	if req.Pagination != nil {
		page := uint64(req.Pagination.Page)
		limit := uint64(req.Pagination.CountOnPage)
		pagination = &pvz.Pagination{
			Page:  &page,
			Limit: &limit,
		}
	}

	options := pvz.ListOptions{
		UserID:     req.UserId,
		InPVZ:      req.InPvz,
		LastN:      lastN,
		Pagination: pagination,
	}

	resp, err := s.ordersHandler.ListOrders(ctx, options)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	ordersProto := mapDomainOrdersToProto(resp.Orders)

	return &pkg.OrdersList{
		Orders: ordersProto,
		Total:  resp.Total,
	}, nil
}
