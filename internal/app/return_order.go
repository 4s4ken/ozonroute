package app

import (
	"context"
	"fmt"
	"strconv"

	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

type ReturnOrderInput struct {
	UserID  pvz.UserID
	OrderID pvz.OrderID
}

func (s *Service) Return(ctx context.Context, input ReturnOrderInput) error {
	order, err := s.repo.GetOrder(ctx, input.OrderID)
	if err != nil {
		return err
	}

	if order.Status == pvz.OrderStatusReturned {
		return errors.EntityValidationError("Order", strconv.FormatUint(uint64(order.OrderID), 10))
	}

	if order.Status == pvz.OrderStatusAccepted && order.ExpiresAt.Before(s.nowFn().UTC()) {
		return errors.EntityExpiredError("Order", strconv.FormatUint(uint64(order.OrderID), 10))
	}

	order.Status = pvz.OrderStatusReturned
	order.ExpiresAt = s.nowFn().UTC()

	err = s.repo.UpdateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	return nil
}
