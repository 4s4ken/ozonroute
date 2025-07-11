package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	buierr "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/outbox"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
	"strconv"
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
		return ProcessOrdersOutput{}, buierr.EntityValidationError("Action", "unsupported action")
	}

	wp, ok := workerpool.WorkerPoolFromContext(ctx)
	if !ok {
		return ProcessOrdersOutput{}, errors.New("worker pool not found in context")
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
		idCopy := id
		err := wp.Submit(func(ctx context.Context) {
			order, err := s.repo.GetOrder(ctx, pvz.OrderID(idCopy))
			if err != nil || order == nil {
				errCh <- res{id: idCopy, err: fmt.Errorf("order not found or error")}

				return
			}

			if err := ProcessActionOrder(order, input.Action, now); err != nil {
				errCh <- res{id: idCopy, err: err}
				return
			}

			var eventType outbox.EventType
			var actorType outbox.ActorType

			switch input.Action {
			case pvz.ActionTypeIssue:
				eventType = outbox.OrderIssued
				actorType = outbox.ActorClient
			case pvz.ActionTypeReturn:
				eventType = outbox.OrderReturnedByClient
				actorType = outbox.ActorClient
			default:
			}

			event := outbox.Event{
				EventID:   uuid.New(),
				EventType: eventType,
				Timestamp: now,
				Actor: outbox.Actor{
					Type: actorType,
					ID:   int(input.UserID),
				},
				Order: outbox.Order{
					ID:     int(order.OrderID),
					UserID: strconv.FormatUint(uint64(order.UserID), 10),
					Status: order.Status.String(),
				},
				Source: "pvz-api",
			}

			if err := s.outboxRepo.SaveEvent(ctx, event); err != nil {
				errCh <- res{id: idCopy, err: err}
				return
			}

			err = s.repo.UpdateOrder(ctx, order)
			errCh <- res{id: idCopy, err: err}
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
