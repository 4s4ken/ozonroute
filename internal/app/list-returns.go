package app

import (
	"context"
	"fmt"
	"sort"

	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// ListReturns возвращает все возвращённые заказы, отсортированные по убыванию срока хранения.
func (s *Storage) ListReturns(ctx context.Context, options pvz.ListReturnsInput) (*pvz.ReturnsList, error) {
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

	if options.Pagination != nil {
		page := options.Pagination.Page
		limit := options.Pagination.CountOnPage
		if limit <= 0 {
			return nil, storageErrors.EntityValidationError("Limit", "must be > 0")
		}

		start := page * limit
		end := start + limit

		if start >= uint64(len(rawReturns)) {
			return nil, storageErrors.EntityValidationError("Pagination", fmt.Sprintf("page %d out of bounds", page))
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

	return &pvz.ReturnsList{
		Returns: returns,
	}, nil
}
