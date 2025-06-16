package grpc

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/app"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/mw"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// ListReturns обрабатывает gRPC-запрос на отправку списка вернувшихся заказов.
func (s *OrdersServer) ListReturns(ctx context.Context, req *pkg.ListReturnsRequest) (*pkg.ReturnsList, error) {

	page := uint64(req.Pagination.Page)
	limit := uint64(req.Pagination.CountOnPage)

	options := app.ListReturnsInput{
		Pagination: pvz.Pagination{
			Page:  &page,
			Limit: &limit,
		},
	}

	resp, err := s.ordersHandler.ListReturns(ctx, options)
	if err != nil {
		return nil, mw.MapBusinessErrorToGRPC(err)
	}

	return mapDomainReturnsListToProto(resp), nil
}
