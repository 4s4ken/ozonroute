package json

import (
	"context"
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"time"
)

// SaveOrder — сохраняет заказ и записывает изменение статуса в историю.
func (s *Storage) SaveOrder(ctx context.Context, order *pvz.Order) error {
	if s.store[order.UserID] == nil {
		s.store[order.UserID] = make(map[pvz.OrderID]*pvz.Order)
	}

	if _, exists := s.store[order.UserID][order.OrderID]; exists {
		return storageErrors.EntityAlreadyExistsError("Order", fmt.Sprint(order.OrderID))
	}

	order.LastUpdate = time.Now().UTC()
	s.store[order.UserID][order.OrderID] = order

	s.history = append(s.history, &pvz.StatusChange{
		OrderID: order.OrderID,
		Status:  order.Status,
		Time:    time.Now().UTC(),
	})

	if err := s.writeDataToFile(); err != nil {
		return err
	}

	return nil
}
