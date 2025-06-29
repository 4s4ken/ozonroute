package app

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
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

func (s *Service) Process(ctx context.Context, input ProcessOrdersInput) (ProcessOrdersOutput, error) {

	if input.Action != pvz.ActionTypeIssue && input.Action != pvz.ActionTypeReturn {
		return ProcessOrdersOutput{}, errors.EntityValidationError("Action", "unsupported action")
	}

	var (
		processed []uint64
		failed    []uint64
		now       = s.nowFn().UTC()
	)

	for _, id := range input.OrderIDs {
		order, err := s.repo.GetOrder(ctx, pvz.OrderID(id))
		if err != nil {
			failed = append(failed, id)
			continue
		}
		if order == nil {
			failed = append(failed, id)
			continue
		}
		
		switch input.Action {
		case pvz.ActionTypeIssue:
			if order.OrderCreate.ExpiresAt.Before(now) {
				failed = append(failed, id)
				continue
			}
			order.Status = pvz.OrderStatusIssued

		case pvz.ActionTypeReturn:
			if order.IsExpiredForReturn(now) {
				failed = append(failed, id)
				continue
			}
			order.Status = pvz.OrderStatusReturned
			order.OrderCreate.ExpiresAt = now

		default:
			return ProcessOrdersOutput{}, errors.EntityValidationError("Action", fmt.Sprint(input.Action))
		}

		err = s.repo.UpdateOrder(ctx, order)
		if err != nil {
			failed = append(failed, id)
			continue
		}

		processed = append(processed, id)
	}

	return ProcessOrdersOutput{
		Processed: processed,
		Errors:    failed,
	}, nil
}
