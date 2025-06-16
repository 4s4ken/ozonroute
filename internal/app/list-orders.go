package app

import (
	"context"
	"fmt"

	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type ListOrdersInput struct {
	UserID uint64
	InPVZ  bool
	LastN  *uint64
	*pvz.Pagination
}

type OrdersListOutput struct {
	Orders []pvz.Order
	Total  int32
}

// ListOrders возвращает заказы пользователя с фильтрами (в ПВЗ, страница, лимит, последние N).
func (s *Storage) ListOrders(ctx context.Context, options ListOrdersInput) (*OrdersListOutput, error) {
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

	if options.Page != nil && options.Limit != nil {
		page := int(*options.Page)
		limit := int(*options.Limit)
		if page < 0 || limit <= 0 {
			return nil, storageErrors.EntityValidationError("Pagination", "page must be >= 0 and limit > 0")
		}

		start := page * limit
		if start >= len(filtered) {
			return nil, storageErrors.EntityValidationError("Pagination", "page out of bounds")
		}

		end := start + limit
		if end > len(filtered) {
			end = len(filtered)
		}

		filtered = filtered[start:end]
	}

	result := make([]pvz.Order, 0, len(filtered))
	for _, o := range filtered {
		result = append(result, *o)
	}

	return &OrdersListOutput{
		Orders: result,
		Total:  int32(len(ordersMap)),
	}, nil
}
