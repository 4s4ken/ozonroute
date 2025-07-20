package json

import (
	"context"
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"sort"
)

// ListOrders — возвращает список заказов с фильтрацией, пагинацией и общее количество.
func (s *Storage) ListOrders(ctx context.Context, options pvz.ListOptions) ([]*pvz.Order, int32, error) {
	userID := pvz.UserID(options.UserID)
	userOrders, ok := s.store[userID]
	if !ok {
		return nil, 0, storageErrors.EntityNotFoundError("User", fmt.Sprint(userID))
	}

	filtered := make([]*pvz.Order, 0, len(userOrders))
	for _, order := range userOrders {
		if options.InPVZ && order.Status != pvz.OrderStatusAccepted {
			continue
		}
		filtered = append(filtered, order)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].OrderID < filtered[j].OrderID
	})

	if options.LastN != nil {
		n := int(*options.LastN)
		if n <= 0 {
			return nil, 0, storageErrors.EntityValidationError("LastN", "must be > 0")
		}
		if n < len(filtered) {
			filtered = filtered[len(filtered)-n:]
		}
	}

	if options.Pagination != nil && options.Pagination.Page != nil && options.Pagination.Limit != nil {
		page := int(*options.Pagination.Page)
		limit := int(*options.Pagination.Limit)

		if page < 0 || limit <= 0 {
			return nil, 0, storageErrors.EntityValidationError("Pagination", "page must be >= 0 and limit > 0")
		}

		start := page * limit
		if start >= len(filtered) {
			return []*pvz.Order{}, int32(len(userOrders)), nil
		}

		end := start + limit
		if end > len(filtered) {
			end = len(filtered)
		}
		filtered = filtered[start:end]
	}

	return filtered, int32(len(userOrders)), nil
}

// ListReturns — возвращает список возвращённых заказов с пагинацией.
func (s *Storage) ListReturns(ctx context.Context, pagination *pvz.Pagination) ([]*pvz.Order, error) {
	var returns []*pvz.Order
	for _, userOrders := range s.store {
		for _, order := range userOrders {
			if order.Status == pvz.OrderStatusReturned {
				returns = append(returns, order)
			}
		}
	}

	sort.Slice(returns, func(i, j int) bool {
		return returns[i].ExpiresAt.After(returns[j].ExpiresAt)
	})

	//total := int32(len(returns))

	if pagination != nil && pagination.Page != nil && pagination.Limit != nil {
		page := int(*pagination.Page)
		limit := int(*pagination.Limit)
		if limit <= 0 {
			return nil, storageErrors.EntityValidationError("Limit", "must be > 0")
		}

		start := page * limit
		if start >= len(returns) {
			return []*pvz.Order{}, nil
		}

		end := start + limit
		if end > len(returns) {
			end = len(returns)
		}
		returns = returns[start:end]
	}

	return returns, nil
}

// GetHistory — возвращает историю статусов всех заказов, отсортированную по времени изменения.
func (s *Storage) GetHistory(ctx context.Context) ([]*pvz.StatusChange, error) {
	var result []*pvz.StatusChange
	for _, userOrders := range s.store {
		for _, order := range userOrders {
			result = append(result, &pvz.StatusChange{
				OrderID: order.OrderID,
				Status:  order.Status,
				Time:    order.LastUpdate,
			})
		}
	}
	return result, nil
}

// GetStatusHistory - заглушка для json репозитория из-за неэффективности метода.
func (s *Storage) GetStatusHistory(ctx context.Context, orderID pvz.OrderID) ([]*pvz.StatusChange, error) {
	return nil, nil
}
