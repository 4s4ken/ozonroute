package app

import (
	"context"
	"strconv"
	"sync"

	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

func (s *Service) Accept(ctx context.Context, order pvz.OrderCreate) error {
	errCh := make(chan error, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)

	err := s.wp.Submit(func(ctx context.Context) {
		defer wg.Done()
		errCh <- s.accept(ctx, order)
	})
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		wg.Wait()
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (s *Service) accept(ctx context.Context, order pvz.OrderCreate) error {
	now := s.nowFn().UTC()
	if order.ExpiresAt.Before(now) {
		return errors.EntityValidationError("Order", strconv.FormatUint(uint64(order.OrderID), 10))
	}

	existing, err := s.repo.GetOrder(ctx, order.OrderID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.EntityAlreadyExistsError("Order", strconv.FormatUint(uint64(order.OrderID), 10))
	}

	finalPrice := order.Price
	if order.Package != nil {
		switch *order.Package {
		case pvz.PackageTypeBag:
			if order.Weight > 10 {
				return errors.EntityWeightError("Bag", strconv.FormatUint(uint64(order.OrderID), 10))
			}
			finalPrice = finalPrice.Add(decimal.NewFromInt(5))
		case pvz.PackageTypeBox:
			if order.Weight > 30 {
				return errors.EntityWeightError("Box", strconv.FormatUint(uint64(order.OrderID), 10))
			}
			finalPrice = finalPrice.Add(decimal.NewFromInt(20))
		case pvz.PackageTypeTape:
			finalPrice = finalPrice.Add(decimal.NewFromInt(1))
		case pvz.PackageTypeBagTape:
			if order.Weight > 10 {
				return errors.EntityWeightError("Bag+Tape", strconv.FormatUint(uint64(order.OrderID), 10))
			}
			finalPrice = finalPrice.Add(decimal.NewFromInt(6))
		case pvz.PackageTypeBoxTape:
			if order.Weight > 30 {
				return errors.EntityWeightError("Box+Tape", strconv.FormatUint(uint64(order.OrderID), 10))
			}
			finalPrice = finalPrice.Add(decimal.NewFromInt(21))
		default:
			return errors.EntityPackageError("Package Type", strconv.FormatUint(uint64(order.OrderID), 10))
		}
	}

	newOrder := &pvz.Order{
		OrderCreate: order,
		Status:      pvz.OrderStatusAccepted,
	}
	newOrder.Price = finalPrice

	return s.repo.SaveOrder(ctx, newOrder)
}
