package app

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/outbox"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/store"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/cache"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
	"strconv"
	"time"
)

// Service — структура бизнес-логики работы с заказами.
type Service struct {
	repo       store.OrderRepository
	nowFn      func() time.Time
	wp         workerpool.WorkerPool
	outboxRepo outbox.OutboxRepository
	orderCache cache.Cache[string, *pvz.Order]
}

func NewService(repo store.OrderRepository, outboxRepo outbox.OutboxRepository, nowFn func() time.Time, wp workerpool.WorkerPool, orderCache cache.Cache[string, *pvz.Order]) *Service {
	return &Service{
		repo:       repo,
		nowFn:      nowFn,
		wp:         wp,
		outboxRepo: outboxRepo,
		orderCache: orderCache,
	}
}

// ResizeWorkerPool — изменяет количество воркеров в пуле.
func (s *Service) ResizeWorkerPool(_ context.Context, n int) error {
	s.wp.Resize(n)
	return nil
}

// ProcessActionOrder — обработка действия issue/return к заказу.
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

// GetOrder — проверяет, находится ли заказ в кеше, в ином случае загружает из репозитория.
func (s *Service) GetOrder(ctx context.Context, orderID pvz.OrderID) (*pvz.Order, error) {
	key := strconv.FormatUint(uint64(orderID), 10)

	if cached, found := s.orderCache.Get(key); found {
		return cached, nil
	}

	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order != nil {
		s.orderCache.Set(key, order)
	}

	return order, nil
}
