package grpc

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
)

// ListReturns обрабатывает gRPC-запрос на отправку списка вернувшихся заказов.
func (s *OrdersServer) ListReturns(ctx context.Context, req *pkg.ListReturnsRequest) (*pkg.ReturnsList, error) {
	input := pvz.ListReturnsInput{
		Pagination: &pvz.PaginationInput{
			Page:        uint64(req.Pagination.Page),
			CountOnPage: uint64(req.Pagination.CountOnPage),
		},
	}
	resp, err := s.ordersHandler.ListReturns(ctx, input)
	if err != nil {
		return nil, MapBusinessErrorToGRPC(err)
	}

	return mapDomainReturnsListToProto(resp), nil
}
