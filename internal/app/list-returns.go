package app

import (
	"context"
	"sort"

	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type ListReturnsInput struct {
	Pagination pvz.Pagination
}

type ReturnsListOutput struct {
	Returns []pvz.Order
}

// ListReturns возвращает все возвращённые заказы, отсортированные по убыванию срока хранения.
func (s *Storage) ListReturns(ctx context.Context, options ListReturnsInput) (*ReturnsListOutput, error) {
	var rawReturns []*pvz.Order

	for _, userOrders := range s.store {
		for _, order := range userOrders {
			if order.Status == pvz.OrderStatusReturned {
				rawReturns = append(rawReturns, order)
			}
		}
	}

	sort.Slice(rawReturns, func(i, j int) bool {
		return rawReturns[i].ExpiresAt.After(rawReturns[j].ExpiresAt)
	})

	if options.Pagination.Page != nil && options.Pagination.Limit != nil {
		page := *options.Pagination.Page
		limit := *options.Pagination.Limit

		if limit == 0 {
			return nil, storageErrors.EntityValidationError("Limit", "must be > 0")
		}

		start := page * limit
		end := start + limit

		if start >= uint64(len(rawReturns)) {
			return &ReturnsListOutput{
				Returns: []pvz.Order{},
			}, nil
		}
		if end > uint64(len(rawReturns)) {
			end = uint64(len(rawReturns))
		}
		rawReturns = rawReturns[start:end]
	}

	returns := make([]pvz.Order, 0, len(rawReturns))
	for _, ptr := range rawReturns {
		returns = append(returns, *ptr)
	}

	return &ReturnsListOutput{
		Returns: returns,
	}, nil
}
