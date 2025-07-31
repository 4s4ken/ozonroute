package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	buierr "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/outbox"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/tracer"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
	"strconv"
)

type ReturnOrderInput struct {
	UserID  pvz.UserID
	OrderID pvz.OrderID
}

// Return — контроль конкурентности через каналы.
func (s *Service) Return(ctx context.Context, input ReturnOrderInput) error {
	wp, ok := workerpool.WorkerPoolFromContext(ctx)
	if !ok {
		return errors.New("worker pool not found in context")
	}

	errCh := make(chan error, 1)

	err := wp.Submit(func(ctx context.Context) {
		errCh <- s.returnOrder(ctx, input)
	})
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

// Return — бизнес-логика возврата заказа курьеру по истечению срока хранения.
func (s *Service) returnOrder(ctx context.Context, input ReturnOrderInput) error {
	ctx, span := tracer.Start(ctx, "Return_order")
	defer span.End()

	now := s.nowFn().UTC()

	order, err := s.GetOrder(ctx, input.OrderID)
	if err != nil {
		return err
	}
	if order == nil {
		return buierr.EntityNotFoundError("Order", strconv.FormatUint(uint64(input.OrderID), 10))
	}
	if order.Status == pvz.OrderStatusReturned {
		return buierr.EntityValidationError("Order", strconv.FormatUint(uint64(order.OrderID), 10))
	}
	if order.Status == pvz.OrderStatusAccepted && order.ExpiresAt.Before(now) {
		return buierr.EntityExpiredError("Order", strconv.FormatUint(uint64(order.OrderID), 10))
	}

	order.Status = pvz.OrderStatusReturned
	order.ExpiresAt = now

	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	s.orderCache.Delete(strconv.FormatUint(uint64(order.OrderID), 10))

	event := outbox.Event{
		EventID:   uuid.New(),
		EventType: outbox.OrderReturnedToCourier,
		Timestamp: now,
		Actor: outbox.Actor{
			Type: outbox.ActorCourier,
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
		return err
	}
	return nil
}
