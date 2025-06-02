package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

// ImportOrders внутренняя логика импорта файла с заказами.
func (s *Storage) ImportOrders(ctx context.Context, filePath string) (int, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("cannot read file %q: %w", filePath, err)
	}

	temp := make(map[string]map[string]*pvz.OrderInfo)
	if err := json.Unmarshal(data, &temp); err != nil {
		return 0, storageErrors.EntityValidationError("JSON file", "invalid format")
	}

	count := 0
	for userStr, ordersMap := range temp {
		uid, err := strconv.ParseUint(userStr, 10, 64)
		if err != nil {
			return 0, storageErrors.EntityValidationError("User", userStr)
		}
		userID := pvz.UserID(uid)

		if _, exists := s.store[userID]; !exists {
			s.store[userID] = make(map[pvz.OrderID]*pvz.OrderInfo)
		}

		for orderStr, order := range ordersMap {
			oid, err := strconv.ParseUint(orderStr, 10, 64)
			if err != nil {
				return 0, storageErrors.EntityValidationError("Order", orderStr)
			}
			orderID := pvz.OrderID(oid)
			s.store[userID][orderID] = order
			count++
		}
	}

	if err := s.writeDataToFile(); err != nil {
		return 0, storageErrors.EntityInternalError("Order storage", "write failure")
	}

	return count, nil
}
