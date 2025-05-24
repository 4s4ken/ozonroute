package storage

import (
	"fmt"
	storageErrors "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/businesserror"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	"sort"
)

func (s *Storage) ScrollOrders(userID uint64, limit uint, lastID *pvz.OrderID) ([]*pvz.OrderInfo, *pvz.OrderID, error) {
	ordersMap, ok := s.store[pvz.UserID(userID)]
	if !ok || len(ordersMap) == 0 {
		return nil, nil, storageErrors.EntityNotFoundError("Order", fmt.Sprint(ordersMap))
	}

	scrollListOrders := make([]*pvz.OrderInfo, 0, len(ordersMap))
	for _, order := range ordersMap {
		scrollListOrders = append(scrollListOrders, order)
	}

	sort.Slice(scrollListOrders, func(i, j int) bool {
		return scrollListOrders[i].Orderid < scrollListOrders[j].Orderid
	})

	startIndex := 0
	if lastID != nil {
		found := false
		for i, order := range scrollListOrders {
			if order.Orderid == *lastID {
				startIndex = i + 1
				found = true
				break
			}
		}
		if !found {
			return nil, nil, fmt.Errorf("lastID %v not found", *lastID)
		}
	}

	endIndex := startIndex + int(limit)
	if endIndex > len(scrollListOrders) {
		endIndex = len(scrollListOrders)
	}
	result := scrollListOrders[startIndex:endIndex]

	var nextLastID *pvz.OrderID
	if endIndex < len(scrollListOrders) {
		nextLastID = &result[len(result)-1].Orderid
	}

	return result, nextLastID, nil
}
