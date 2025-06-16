package grpc

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/mw"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// ListOrders обрабатывает gRPC-запрос на отправку списка заказов.
func (s *OrdersServer) ListOrders(ctx context.Context, req *pkg.ListOrdersRequest) (*pkg.OrdersList, error) {
	var lastN *uint64
	if req.LastN != nil {
		tmp := uint64(*req.LastN)
		lastN = &tmp
	}

	page := uint64(req.Pagination.Page)
	limit := uint64(req.Pagination.CountOnPage)

	options := app.ListOrdersInput{
		UserID: req.UserId,
		InPVZ:  req.InPvz,
		LastN:  lastN,
		Pagination: &pvz.Pagination{
			Page:  &page,
			Limit: &limit,
		},
	}

	resp, err := s.ordersHandler.ListOrders(ctx, options)
	if err != nil {
		return nil, mw.MapBusinessErrorToGRPC(err)
	}

	ordersProto := mapDomainOrdersToProto(resp.Orders)

	return &pkg.OrdersList{
		Orders: ordersProto,
		Total:  resp.Total,
	}, nil
}
