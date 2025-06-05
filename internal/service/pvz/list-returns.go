package pvz

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"sort"
)

func (s *OrdersService) ListReturns(ctx context.Context, req *pvz.ListReturnsRequest) (*pvz.ReturnsList, error) {
	page := uint32(1)
	count := uint32(0)

	if req.Pagination != nil {
		page = req.Pagination.Page
		count = req.Pagination.CountOnPage
	}

	var returns []*pvz.Order
	for _, order := range s.storage {
		if order.Status == pvz.OrderStatus_ORDER_STATUS_RETURNED {
			returns = append(returns, order)
		}
	}

	sort.Slice(returns, func(i, j int) bool {
		return returns[i].ExpiresAt.AsTime().After(returns[j].ExpiresAt.AsTime())
	})

	if req.Pagination == nil || count == 0 {
		return &pvz.ReturnsList{
			Returns: returns,
		}, nil
	}

	start := (page - 1) * count
	if start >= uint32(len(returns)) {
		return &pvz.ReturnsList{Returns: []*pvz.Order{}}, nil
	}
	end := start + count
	if end > uint32(len(returns)) {
		end = uint32(len(returns))
	}

	return &pvz.ReturnsList{
		Returns: returns[start:end],
	}, nil
}
