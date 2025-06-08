package grpc

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"

	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// ListOrders обрабатывает gRPC-запрос на отправку списка заказов.
func (s *OrdersServer) ListOrders(ctx context.Context, req *pkg.ListOrdersRequest) (*pkg.OrdersList, error) {
	var lastN *uint64
	if req.LastN != nil {
		tmp := uint64(*req.LastN)
		lastN = &tmp
	}

	input := pvz.ListOrdersInput{
		UserID: pvz.UserID(req.UserId),
		InPVZ:  req.InPvz,
		LastN:  lastN,
		Pagination: &pvz.PaginationInput{
			Page:        uint64(req.Pagination.Page),
			CountOnPage: uint64(req.Pagination.CountOnPage),
		},
	}

	resp, err := s.ordersHandler.ListOrders(ctx, input)
	if err != nil {
		return nil, MapBusinessErrorToGRPC(err)
	}

	ordersProto := mapDomainOrdersToProto(resp.Orders)

	return &pkg.OrdersList{
		Orders: ordersProto,
		Total:  int32(resp.Total),
	}, nil
}
