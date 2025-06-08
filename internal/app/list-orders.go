package app

import (
	"context"
	"fmt"

	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// ListOrders возвращает заказы пользователя с фильтрами (в ПВЗ, страница, лимит, последние N).
func (s *Storage) ListOrders(ctx context.Context, options pvz.ListOrdersInput) (*pvz.OrdersList, error) {
	userID := pvz.UserID(options.UserID)
	ordersMap, ok := s.store[userID]
	if !ok {
		return nil, storageErrors.EntityNotFoundError("User", fmt.Sprint(userID))
	}

	var filtered []*pvz.Order
	for _, order := range ordersMap {
		if options.InPVZ && order.Status != pvz.OrderStatusAccepted {
			continue
		}
		filtered = append(filtered, order)
	}

	if options.LastN != nil {
		n := int(*options.LastN)
		if n <= 0 {
			return nil, storageErrors.EntityValidationError("Last", "must be > 0")
		}
		if n < len(filtered) {
			filtered = filtered[len(filtered)-n:]
		}
	}

	if options.Pagination != nil {
		page := int(options.Pagination.Page)
		limit := int(options.Pagination.CountOnPage)
		start := page * limit
		end := start + limit

		if start >= len(filtered) {
			return nil, storageErrors.EntityValidationError("Pagination", "page out of bounds")
		}
		if end > len(filtered) {
			end = len(filtered)
		}
		filtered = filtered[start:end]
	}

	result := make([]pvz.Order, 0, len(filtered))
	for _, o := range filtered {
		result = append(result, *o)
	}

	return &pvz.OrdersList{
		Orders: result,
		Total:  int32(len(ordersMap)),
	}, nil
}
