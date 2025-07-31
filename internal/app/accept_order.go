package app

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	buierr "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/outbox"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/metrics"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/tracer"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
	"go.opentelemetry.io/otel/attribute"
	"strconv"
)

// Accept — контроль конкурентности через каналы.
func (s *Service) Accept(ctx context.Context, order pvz.OrderCreate) error {
	wp, ok := workerpool.WorkerPoolFromContext(ctx)
	if !ok {
		return errors.New("worker pool not found in context")
	}

	errCh := make(chan error, 1)

	err := wp.Submit(func(ctx context.Context) {
		errCh <- s.accept(ctx, order)
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

// Accept — бизнес-логика принятия и валидации заказа.
func (s *Service) accept(ctx context.Context, order pvz.OrderCreate) error {

	ctx, span := tracer.Start(ctx, "Accept_order")
	defer span.End()

	span.SetAttributes(
		attribute.Int("order.id", int(order.OrderID)),
		attribute.Int("order.user_id", int(order.UserID)),
	)

	now := s.nowFn().UTC()
	if order.ExpiresAt.Before(now) {
		return buierr.EntityValidationError("Order", strconv.FormatUint(uint64(order.OrderID), 10))
	}

	existing, err := s.GetOrder(ctx, order.OrderID)
	if err != nil {
		return err
	}
	if existing != nil {
		return buierr.EntityAlreadyExistsError("Order", strconv.FormatUint(uint64(order.OrderID), 10))
	}

	finalPrice := order.Price
	if order.Package != nil {
		switch *order.Package {
		case pvz.PackageTypeBag:
			if order.Weight > 10 {
				return buierr.EntityWeightError("Bag", strconv.FormatUint(uint64(order.OrderID), 10))
			}
			finalPrice = finalPrice.Add(decimal.NewFromInt(5))
		case pvz.PackageTypeBox:
			if order.Weight > 30 {
				return buierr.EntityWeightError("Box", strconv.FormatUint(uint64(order.OrderID), 10))
			}
			finalPrice = finalPrice.Add(decimal.NewFromInt(20))
		case pvz.PackageTypeTape:
			finalPrice = finalPrice.Add(decimal.NewFromInt(1))
		case pvz.PackageTypeBagTape:
			if order.Weight > 10 {
				return buierr.EntityWeightError("Bag+Tape", strconv.FormatUint(uint64(order.OrderID), 10))
			}
			finalPrice = finalPrice.Add(decimal.NewFromInt(6))
		case pvz.PackageTypeBoxTape:
			if order.Weight > 30 {
				return buierr.EntityWeightError("Box+Tape", strconv.FormatUint(uint64(order.OrderID), 10))
			}
			finalPrice = finalPrice.Add(decimal.NewFromInt(21))
		default:
			return buierr.EntityPackageError("Package Type", strconv.FormatUint(uint64(order.OrderID), 10))
		}
	}

	newOrder := &pvz.Order{
		OrderCreate: order,
		Status:      pvz.OrderStatusAccepted,
	}
	newOrder.Price = finalPrice
	event := outbox.Event{
		EventID:   uuid.New(),
		EventType: outbox.OrderAccepted,
		Timestamp: now,
		Actor: outbox.Actor{
			Type: outbox.ActorClient,
			ID:   int(order.UserID),
		},
		Order: outbox.Order{
			ID:     int(order.OrderID),
			UserID: strconv.FormatUint(uint64(order.UserID), 10),
			Status: pvz.OrderStatusAccepted.String(),
		},
		Source: "pvz-api",
	}

	if err := s.outboxRepo.SaveEvent(ctx, event); err != nil {
		return err
	}
	if err := s.repo.SaveOrder(ctx, newOrder); err != nil {
		return err
	}
	metrics.OrdersAccepted.Inc()
	s.orderCache.Delete(strconv.FormatUint(uint64(order.OrderID), 10))

	return nil
}
