package app

import (
	"context"
	"github.com/shopspring/decimal"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"strconv"
)

// accept — внутренняя логика добавления заказа с проверками.
func (s *Storage) accept(ctx context.Context, order pvz.OrderCreate) error {
	now := s.nowFn().UTC()

	if order.ExpiresAt.Before(now) {
		return storageErrors.EntityValidationError("Order", strconv.FormatUint(uint64(order.OrderID), 10))
	}

	if _, ok := s.store[order.UserID]; !ok {
		s.store[order.UserID] = make(map[pvz.OrderID]*pvz.Order)
	}

	if _, exists := s.store[order.UserID][order.OrderID]; exists {
		return storageErrors.EntityAlreadyExistsError("Order", strconv.FormatUint(uint64(order.OrderID), 10))
	}

	price := decimal.NewFromFloat(float64(order.Price))

	if order.Package != nil {
		switch *order.Package {
		case pvz.PackageTypeBag:
			if order.Weight > 10 {
				return storageErrors.EntityWeightError("Bag", strconv.FormatUint(uint64(order.OrderID), 10))
			}
			price = price.Add(decimal.NewFromInt(5))
		case pvz.PackageTypeBox:
			if order.Weight > 30 {
				return storageErrors.EntityWeightError("Box", strconv.FormatUint(uint64(order.OrderID), 10))
			}
			price = price.Add(decimal.NewFromInt(20))
		case pvz.PackageTypeTape:
			price = price.Add(decimal.NewFromInt(1))
		case pvz.PackageTypeBagTape:
			if order.Weight > 10 {
				return storageErrors.EntityWeightError("Bag+Tape", strconv.FormatUint(uint64(order.OrderID), 10))
			}
			price = price.Add(decimal.NewFromInt(6))
		case pvz.PackageTypeBoxTape:
			if order.Weight > 30 {
				return storageErrors.EntityWeightError("Box+Tape", strconv.FormatUint(uint64(order.OrderID), 10))
			}
			price = price.Add(decimal.NewFromInt(21))
		default:
			return storageErrors.EntityPackageError("Package Type", strconv.FormatUint(uint64(order.OrderID), 10))
		}
	}

	s.store[order.UserID][order.OrderID] = &pvz.Order{
		OrderCreate: order,
		Status:      pvz.OrderStatusAccepted,
	}

	err := s.writeDataToFile()
	return err
}

// Accept — публичный метод, вызывающий внутренний accept.
func (s *Storage) Accept(ctx context.Context, order pvz.OrderCreate) error {
	return s.accept(ctx, order)
}
