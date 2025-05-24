package storage

import "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"

func (s *Storage) ListOrders(options pvz.ListOrdersOptions) ([]*pvz.OrderInfo, error) {
	ordersMap, ok := s.store[pvz.UserID(options.UserID)]
	if !ok || len(ordersMap) == 0 {
		return nil, nil
	}

	listOrd := make([]*pvz.OrderInfo, 0, len(ordersMap))
	for _, order := range ordersMap {
		if options.InPvz && order.Status != pvz.StatusAccepted {
			continue
		}
		listOrd = append(listOrd, order)
	}

	if options.Page != nil && options.Limit != nil {
		start := int(*options.Page) * int(*options.Limit)
		end := start + int(*options.Limit)
		if start < len(listOrd) {
			if end > len(listOrd) {
				end = len(listOrd)
			}
			listOrd = listOrd[start:end]
		} else {
			listOrd = []*pvz.OrderInfo{}
		}
	}

	if options.Last != nil {
		n := int(*options.Last)
		if n < len(listOrd) {
			listOrd = listOrd[len(listOrd)-n:]
		}
	}

	return listOrd, nil
}
