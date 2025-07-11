package app

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/outbox"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/store"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
	"time"
)

type Service struct {
	repo       store.OrderRepository
	nowFn      func() time.Time
	wp         workerpool.WorkerPool
	outboxRepo outbox.OutboxRepository
}

func NewService(repo store.OrderRepository, outboxRepo outbox.OutboxRepository, nowFn func() time.Time, wp workerpool.WorkerPool) *Service {
	return &Service{
		repo:       repo,
		nowFn:      nowFn,
		wp:         wp,
		outboxRepo: outboxRepo,
	}
}

func (s *Service) ResizeWorkerPool(_ context.Context, n int) error {
	s.wp.Resize(n)
	return nil
}

func ProcessActionOrder(order *pvz.Order, action pvz.ActionType, now time.Time) error {
	switch action {
	case pvz.ActionTypeIssue:
		if order.OrderCreate.ExpiresAt.Before(now) {
			return fmt.Errorf("expired")
		}
		order.Status = pvz.OrderStatusIssued

	case pvz.ActionTypeReturn:
		if order.IsExpiredForReturn(now) {
			return fmt.Errorf("expired for return")
		}
		order.Status = pvz.OrderStatusReturned
		order.OrderCreate.ExpiresAt = now

	default:
		return errors.EntityValidationError("Action", "unsupported action")
	}
	return nil
}
