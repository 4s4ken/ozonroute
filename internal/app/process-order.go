package app

import (
	"context"
	"fmt"

	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type ProcessOrdersInput struct {
	UserID   uint64
	Action   pvz.ActionType
	OrderIDs []uint64
}

type ProcessOrdersOutput struct {
	Processed []uint64
	Errors    []uint64
}

// ProcessOrders — обработка заказов покупателей (выдача/возврат).
func (s *Storage) ProcessOrders(ctx context.Context, orders ProcessOrdersInput) (ProcessOrdersOutput, error) {
	ordersByUser, ok := s.store[pvz.UserID(orders.UserID)]
	if !ok {
		return ProcessOrdersOutput{}, storageErrors.EntityNotFoundError("User", fmt.Sprint(orders.UserID))
	}

	var processed []uint64
	var errors []uint64

	now := s.nowFn().UTC()

	for _, orderID := range orders.OrderIDs {
		order, ok := ordersByUser[pvz.OrderID(orderID)]
		if !ok {
			errors = append(errors, orderID)
			continue
		}

		switch orders.Action {
		case pvz.ActionTypeIssue:
			if order.IsExpired(now) {
				errors = append(errors, orderID)
				continue
			}
			order.Status = pvz.OrderStatusIssued

		case pvz.ActionTypeReturn:
			if order.IsExpiredForReturn(now) {
				errors = append(errors, orderID)
				continue
			}
			order.Status = pvz.OrderStatusReturned
			order.ExpiresAt = s.nowFn()

		default:
			return ProcessOrdersOutput{}, storageErrors.EntityValidationError("Action", fmt.Sprint(orders.Action))
		}

		processed = append(processed, orderID)
	}

	if err := s.writeDataToFile(); err != nil {
		return ProcessOrdersOutput{}, err
	}

	return ProcessOrdersOutput{
		Processed: processed,
		Errors:    errors,
	}, nil
}

func (s *Storage) Process(ctx context.Context, input ProcessOrdersInput) (ProcessOrdersOutput, error) {
	return s.ProcessOrders(ctx, input)
}
