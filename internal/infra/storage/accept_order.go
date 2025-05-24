package storage

import (
	"context"
	storageerrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"strconv"
	"time"
)

func (s *Storage) accept(ctx context.Context, orderID pvz.OrderID, userID pvz.UserID, expires time.Time) (err error) {
	if expires.Before(time.Now().UTC()) {
		return storageerrors.EntityValidationError("Order", strconv.FormatUint(uint64(orderID), 10))
	}
	if _, ok := s.store[userID]; !ok {
		s.store[userID] = make(map[pvz.OrderID]*pvz.OrderInfo)
	}

	if _, exists := s.store[userID][orderID]; exists {
		return storageerrors.EntityAlreadyExistsError("Order", strconv.FormatUint(uint64(orderID), 10))
	}

	s.store[userID][orderID] = &pvz.OrderInfo{
		Userid:        userID,
		Orderid:       orderID,
		Time:          expires,
		Status:        1,
		StatusChanged: true,
	}
	return s.writeDataToFile()
}

func (s *Storage) Accept(ctx context.Context, orderID pvz.OrderID, userID pvz.UserID, exp time.Time) error {
	return s.accept(ctx, orderID, userID, exp)
}
