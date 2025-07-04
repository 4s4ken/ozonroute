package app

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"sync"
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
		mu        sync.Mutex
		processed []uint64
		failed    []uint64
		now       = s.nowFn().UTC()
	)

	type res struct {
		id  uint64
		err error
	}

	errCh := make(chan res, len(input.OrderIDs))

	for _, id := range input.OrderIDs {
		err := s.wp.Submit(func(ctx context.Context) {
			order, err := s.repo.GetOrder(ctx, pvz.OrderID(id))
			if err != nil || order == nil {
				errCh <- res{id: id, err: fmt.Errorf("order not found or error")}
				return
			}

			switch input.Action {
			case pvz.ActionTypeIssue:
				if order.OrderCreate.ExpiresAt.Before(now) {
					errCh <- res{id: id, err: fmt.Errorf("expired")}
					return
				}
				order.Status = pvz.OrderStatusIssued

			case pvz.ActionTypeReturn:
				if order.IsExpiredForReturn(now) {
					errCh <- res{id: id, err: fmt.Errorf("expired for return")}
					return
				}
				order.Status = pvz.OrderStatusReturned
				order.OrderCreate.ExpiresAt = now

			default:
				errCh <- res{id: id, err: fmt.Errorf("unsupported action")}
				return
			}

			err = s.repo.UpdateOrder(ctx, order)
			errCh <- res{id: id, err: err}
		})
		if err != nil {
			return ProcessOrdersOutput{}, err
		}
	}

	for i := 0; i < len(input.OrderIDs); i++ {
		r := <-errCh
		mu.Lock()
		if r.err != nil {
			failed = append(failed, r.id)
		} else {
			processed = append(processed, r.id)
		}
		mu.Unlock()
	}

	return ProcessOrdersOutput{
		Processed: processed,
		Errors:    failed,
	}, nil
}
