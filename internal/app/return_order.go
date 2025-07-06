package app

import (
	"context"
	"errors"
	"fmt"
	buierr "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
	"strconv"
)

type ReturnOrderInput struct {
	UserID  pvz.UserID
	OrderID pvz.OrderID
}

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

func (s *Service) returnOrder(ctx context.Context, input ReturnOrderInput) error {
	order, err := s.repo.GetOrder(ctx, input.OrderID)
	if err != nil {
		return err
	}

	if order == nil {
		return buierr.EntityNotFoundError("Order", strconv.FormatUint(uint64(input.OrderID), 10))
	}

	if order.Status == pvz.OrderStatusReturned {
		return buierr.EntityValidationError("Order", strconv.FormatUint(uint64(order.OrderID), 10))
	}

	if order.Status == pvz.OrderStatusAccepted && order.ExpiresAt.Before(s.nowFn().UTC()) {
		return buierr.EntityExpiredError("Order", strconv.FormatUint(uint64(order.OrderID), 10))
	}

	order.Status = pvz.OrderStatusReturned
	order.ExpiresAt = s.nowFn().UTC()

	err = s.repo.UpdateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	return nil
}
