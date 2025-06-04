package storage

import (
	"context"
	"github.com/shopspring/decimal"
	storageerrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"strconv"
	"time"
)

// Accept внутренняя логика добавления заказа с проверками.
func (s *Storage) accept(ctx context.Context, orderID pvz.OrderID, userID pvz.UserID, expires time.Time, weight pvz.OrderWeight, price pvz.OrderPrice, packageType pvz.PackageType) (pvz.OrderPrice, error) {
	if expires.Before(s.nowFn().UTC()) {
		return price, storageerrors.EntityValidationError("Order", strconv.FormatUint(uint64(orderID), 10))
	}
	if _, ok := s.store[userID]; !ok {
		s.store[userID] = make(map[pvz.OrderID]*pvz.OrderInfo)
	}

	if _, exists := s.store[userID][orderID]; exists {
		return price, storageerrors.EntityAlreadyExistsError("Order", strconv.FormatUint(uint64(orderID), 10))
	}

	switch packageType {
	case "":
	case pvz.Bag:
		if weight >= 10 {
			return price, storageerrors.EntityWeightError("Bag", strconv.FormatUint(uint64(orderID), 10))
		}
		price = price.Add(decimal.NewFromInt(5))
	case pvz.Box:
		if weight >= 30 {
			return price, storageerrors.EntityWeightError("Box", string(packageType))
		}
		price = price.Add(decimal.NewFromInt(20))
	case pvz.Film:
		price = price.Add(decimal.NewFromInt(1))
	case pvz.BagAndFilm:
		if weight >= 10 {
			return price, storageerrors.EntityWeightError("Bag+film", string(packageType))
		}
		price = price.Add(decimal.NewFromInt(6))
	case pvz.BoxAndFilm:
		if weight >= 30 {
			return price, storageerrors.EntityWeightError("Box+film", string(packageType))
		}
		price = price.Add(decimal.NewFromInt(21))
	default:
		return price, storageerrors.EntityPackageError("Package Type", string(packageType))
	}

	s.store[userID][orderID] = &pvz.OrderInfo{
		Userid:        userID,
		Orderid:       orderID,
		Time:          expires,
		Status:        pvz.StatusAccepted,
		StatusChanged: true,
		Weight:        weight,
		Price:         price,
		Package:       &packageType,
	}
	err := s.writeDataToFile()
	return price, err
}

func (s *Storage) Accept(ctx context.Context, orderID pvz.OrderID, userID pvz.UserID, exp time.Time, weight pvz.OrderWeight, price pvz.OrderPrice, packageType pvz.PackageType) (pvz.OrderPrice, error) {
	return s.accept(ctx, orderID, userID, exp, weight, price, packageType)
}
