package storage

import (
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
)

func (s *Storage) ListReturns(options pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error) {

	returnedOrders := make([]*pvz.OrderInfo, 0)
	for _, userOrders := range s.store {
		for _, order := range userOrders {
			if order.Status == 3 {
				returnedOrders = append(returnedOrders, order)
			}
		}
	}

	if options.Page != nil && options.Limit != nil {
		start := int(*options.Page) * int(*options.Limit)
		end := start + int(*options.Limit)

		if start < len(returnedOrders) {
			if end > len(returnedOrders) {
				end = len(returnedOrders)
			}
			return returnedOrders[start:end], nil
		}
		return []*pvz.OrderInfo{}, nil
	}

	return returnedOrders, nil
}
